# Nook Sync - Save Context to Base Workspace

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/nook-sync/workflow.yaml</critical>

## Purpose

Save working context (docs, configs, etc.) **TO the base workspace's tracking folder**.
This is the single source of truth for all workspace context.

**Key Architecture Points:**
- All context is saved to `base_workspace_path/.bmad-tracking/{branch}/`
- Can be run from ANY worktree (base or nook)
- Sync artifacts from current working directory to base workspace
- Updates or creates context.yaml with lineage information

<workflow>

<step n="1" goal="Check if workspace is initialized and get base path">
  <action>Read {module_config} and check base_workspace_path:</action>
  ```bash
  grep "^base_workspace_path:" .bmad/th/config.yaml
  ```

  <check if="base_workspace_path is empty or not set">
    <action>Display critical halt:</action>
    ```
    CANNOT SYNC - WORKSPACE NOT INITIALIZED

    No base workspace has been configured.
    The base workspace holds the tracking folder (.bmad-tracking/)
    which is the single source of truth for all context.

    Please run first (from the base worktree):
      /bmad:th:workflows:treehouse-init

    This will establish the tracking folder location.
    ```
    <action>HALT workflow - do not proceed</action>
  </check>

  <action>Store base_workspace_path as {{base_path}}</action>
  <action>Verify base_path exists and is accessible:</action>
  ```bash
  test -d "{{base_path}}"
  ```

  <check if="base_path does not exist">
    <action>Display error:</action>
    ```
    ERROR - BASE WORKSPACE NOT FOUND

    Configured base_workspace_path does not exist:
      {{base_path}}

    The base workspace may have been moved or deleted.

    Options:
    - Ensure the base worktree exists at that path
    - Re-run /bmad:th:workflows:treehouse-init from the base
    - Update base_workspace_path in .bmad/th/config.yaml
    ```
    <action>HALT workflow</action>
  </check>

  <action>Display: "Base workspace: {{base_path}}"</action>
</step>

<step n="2" goal="Detect current git branch and hash">
  <action>Execute: `git branch --show-current`</action>
  <action>Store result as {{current_branch}}</action>
  <action>If not in a git repository or no branch detected, halt with error</action>

  <action>Get short hash of current HEAD (4 chars): `git rev-parse --short=4 HEAD`</action>
  <action>Store as {{current_hash}}</action>

  <action>Get current working directory: `pwd`</action>
  <action>Store as {{current_path}}</action>
</step>

<step n="3" goal="Load sync_paths from config">
  <action>Read sync_paths from {module_config}:</action>
  ```yaml
  sync_paths:
    - docs
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml
  ```
  <action>Store as {{sync_paths}} array</action>
</step>

<step n="4" goal="Derive tracking path in base workspace">
  <action>Calculate tracking path: `{{base_path}}/.bmad-tracking/{{current_branch}}`</action>
  <action>The branch hierarchy is preserved (e.g., `discovery/peek-a-box-mvp` -> `.bmad-tracking/discovery/peek-a-box-mvp/`)</action>
  <action>Store as {{tracking_path}}</action>
</step>

<step n="5" goal="Load existing context if available">
  <action>Check if {{tracking_path}}/context.yaml exists</action>

  <check if="context.yaml exists">
    <action>Load existing context.yaml</action>
    <action>Extract lineage array if present</action>
    <action>Store as {{existing_lineage}}</action>
    <action>Store existing artifacts info as {{existing_artifacts}}</action>
  </check>

  <check if="context.yaml does not exist">
    <action>Initialize empty lineage array</action>
    <action>Extract type from branch name (e.g., "discovery/foo" -> type: "discovery")</action>
    <action>If branch has no type prefix, use "base" as type</action>
  </check>
</step>

<step n="6" goal="Create tracking directories if missing">
  <action>Check if {{tracking_path}} exists</action>
  <action>Create directory structure: `mkdir -p "{{tracking_path}}"`</action>

  <action>For each path in {{sync_paths}}, create parent directory in tracking:
  ```bash
  # For each sync_path, ensure directory structure exists
  # e.g., docs -> {{tracking_path}}/docs/
  # e.g., .bmad/_cfg/agents -> {{tracking_path}}/.bmad/_cfg/agents/
  # e.g., .bmad/bmm/config.yaml -> {{tracking_path}}/.bmad/bmm/
  ```
  </action>
</step>

<step n="7" goal="Show sync preview">
  <action>Check which source paths exist and have content in current directory</action>

  <action>Display what will be synced:</action>

  ```
  Sync Preview

  Source:  {{current_path}} ({{current_branch}})
  Target:  {{tracking_path}}
  Hash:    {{current_hash}}

  Paths to sync:
  {{for each path in sync_paths}}
  {{if path exists}}
  [EXISTS] {{path}} -> {{tracking_path}}/{{path}}
  {{else}}
  [SKIP]   {{path}} (not found in current directory)
  {{end if}}
  {{end for}}

  WARNING: Existing content in tracking will be REPLACED (not merged)
  ```

  <check if="no paths exist to sync">
    <action>Display warning:</action>
    ```
    WARNING - NOTHING TO SYNC

    None of the configured sync_paths exist in the current directory:
    {{list sync_paths}}

    This is expected if:
    - This is a fresh nook that hasn't created any artifacts yet
    - You're in a worktree where sync artifacts are gitignored

    The context.yaml will still be updated with branch info.
    ```
  </check>
</step>

<step n="8" goal="Confirm operation">
  <ask>Proceed with sync to base workspace? [y/n]</ask>
  <action if="user declines">Abort workflow with message: "Sync cancelled by user"</action>
</step>

<step n="9" goal="Execute sync operations">
  <action>For each path in {{sync_paths}} that exists:</action>

  <action>Determine if path is file or directory</action>

  <action>For directories:
  ```bash
  rsync -av --delete "{{current_path}}/{{path}}/" "{{tracking_path}}/{{path}}/"
  ```
  </action>

  <action>For files:
  ```bash
  mkdir -p "$(dirname "{{tracking_path}}/{{path}}")"
  cp "{{current_path}}/{{path}}" "{{tracking_path}}/{{path}}"
  ```
  </action>

  <action>Record which paths were synced and file counts</action>
</step>

<step n="10" goal="Build and update lineage">
  <action>Build current branch's lineage entry:</action>
  ```yaml
  current_entry:
    branch: {{current_branch}}
    hash: {{current_hash}}
    type: (extracted from branch prefix, or "base" if none)
    name: (derive human-readable name from branch)
    created: (use existing created date if updating, else {{date}})
    updated: {{date}}
  ```

  <action>Check if current branch already in lineage:</action>
  <check if="branch already in lineage">
    <action>Update existing entry with new hash and updated timestamp</action>
  </check>
  <check if="branch not in lineage">
    <action>Append current_entry to lineage</action>
  </check>

  <action>Store final lineage as {{full_lineage}}</action>
</step>

<step n="11" goal="Generate context manifest with lineage">
  <action>Create or update `{{tracking_path}}/context.yaml` with:</action>

  ```yaml
  # Context Manifest - Auto-generated by nook-sync
  branch: {{current_branch}}
  hash: {{current_hash}}
  type: {{branch_type}}
  synced_at: {{date}}
  synced_by: {user_name}
  synced_from: {{current_path}}

  # Full lineage chain (oldest to newest)
  lineage:
    {{for each entry in full_lineage}}
    - branch: {{entry.branch}}
      hash: {{entry.hash}}
      type: {{entry.type}}
      name: "{{entry.name}}"
      created: {{entry.created}}
      updated: {{entry.updated}}
      {{if entry.parent_hash}}parent_hash: {{entry.parent_hash}}{{end if}}
    {{end for}}

  # Artifacts synced (based on sync_paths config)
  artifacts:
    {{for each path in sync_paths}}
    {{path_key}}:
      synced: true/false
      files: [count]
      source: "{{current_path}}/{{path}}"
    {{end for}}

  notes: |
    Context synced from {{current_path}}.
    Use /bmad:th:workflows:nook-restore to restore this context.
    Tracking folder is in base workspace: {{base_path}}
  ```
</step>

<step n="12" goal="Report results and remind about commit">
  <action>Display completion summary:</action>

  ```
  Sync Complete!

  Source: {{current_path}}
  Branch: {{current_branch}} ({{current_hash}})
  Target: {{tracking_path}}

  Synced:
  {{for each path in sync_paths}}
  {{if synced}}
  [DONE] {{path}} ({{file_count}} files)
  {{else}}
  [SKIP] {{path}} (not found)
  {{end if}}
  {{end for}}

  Lineage ({{lineage_count}} levels):
  {{for each entry in full_lineage}}
     {{index}}. {{entry.branch}} ({{entry.hash}})
  {{end for}}

  Manifest: {{tracking_path}}/context.yaml


  IMPORTANT: Changes saved to BASE WORKSPACE

  The tracking folder is in: {{base_path}}/.bmad-tracking/

  To make this context available for nooks:
  1. Go to base workspace: cd {{base_path}}
  2. Stage changes: git add .bmad-tracking/
  3. Commit: git commit -m "chore(th): sync context for {{current_branch}}"

  To restore this context in another worktree:
    /bmad:th:workflows:nook-restore
  ```
</step>

</workflow>

## Notes

<notes>
**Architecture:**
- All context is stored in base_workspace_path/.bmad-tracking/
- This workflow copies FROM current directory TO base workspace
- Can be run from any worktree (base or nook)
- The tracking folder only exists in the base workspace

**Sync Paths:**
- Configured in .bmad/th/config.yaml under sync_paths
- Directories are synced with rsync --delete (mirror copy)
- Files are copied directly
- Non-existent paths are skipped (logged in context.yaml)

**Lineage Tracking:**
- Each branch gets a lineage entry
- Hash is updated on each sync
- Created date is preserved, updated date changes
- Full chain allows tracing nook history

**Commit Responsibility:**
- This workflow does NOT auto-commit
- User must commit the tracking folder in base workspace
- This is intentional - allows batching multiple syncs
</notes>
