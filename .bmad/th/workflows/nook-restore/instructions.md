# Nook Restore - Load Context from Base Workspace

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/nook-restore/workflow.yaml</critical>

## Purpose

Restore saved context (docs, configs, etc.) **FROM the base workspace's tracking folder** to current worktree.

**Key Architecture Points:**
- All context is loaded from `base_workspace_path/.bmad-tracking/{branch}/`
- Can restore context from ANY branch (not just current)
- Common use: nook starts clean, then restore parent's context
- Overwrites local content with saved version

<workflow>

<step n="1" goal="Check if workspace is initialized and get base path">
  <action>Read {module_config} and check base_workspace_path:</action>
  ```bash
  grep "^base_workspace_path:" .bmad/th/config.yaml
  ```

  <check if="base_workspace_path is empty or not set">
    <action>Display critical halt:</action>
    ```
    CANNOT RESTORE - WORKSPACE NOT INITIALIZED

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

<step n="2" goal="Detect current branch and determine source">
  <action>Execute: `git branch --show-current`</action>
  <action>Store result as {{current_branch}}</action>

  <action>Get current working directory: `pwd`</action>
  <action>Store as {{current_path}}</action>

  <check if="source_branch provided as input">
    <action>Use provided source_branch as {{restore_branch}}</action>
  </check>

  <check if="source_branch NOT provided">
    <action>Check if current branch has context in tracking:</action>
    ```bash
    test -d "{{base_path}}/.bmad-tracking/{{current_branch}}"
    ```

    <check if="context exists for current branch">
      <action>Default to current branch: {{restore_branch}} = {{current_branch}}</action>
    </check>

    <check if="context does NOT exist for current branch">
      <action>Check if this is a nook by looking for parent in context.yaml:</action>
      ```bash
      cat "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml" 2>/dev/null | grep "forked_from:"
      ```

      <check if="forked_from exists">
        <action>Extract parent branch as potential restore source</action>
        <action>Store as {{parent_branch}}</action>
      </check>

      <action>Display source selection:</action>
      ```
      SELECT RESTORE SOURCE

      Current branch '{{current_branch}}' has no synced context.

      {{if parent_branch exists}}
      This branch was forked from: {{parent_branch}}

      Options:
      [P] Restore from parent ({{parent_branch}})
      [L] List all available contexts
      [C] Cancel
      {{else}}
      Options:
      [L] List all available contexts
      [C] Cancel
      {{end if}}

      Choose:
      ```

      <check if="user chooses P (parent)">
        <action>Set {{restore_branch}} = {{parent_branch}}</action>
      </check>

      <check if="user chooses L (list)">
        <action>List all available contexts:</action>
        ```bash
        ls -1 "{{base_path}}/.bmad-tracking/"
        ```
        <ask>Enter branch name to restore from:</ask>
        <action>Set {{restore_branch}} from user input</action>
      </check>

      <check if="user chooses C (cancel)">
        <action>Display: "Restore cancelled."</action>
        <action>HALT workflow</action>
      </check>
    </check>
  </check>
</step>

<step n="3" goal="Validate tracking path exists">
  <action>Calculate tracking path: `{{base_path}}/.bmad-tracking/{{restore_branch}}`</action>
  <action>Store as {{tracking_path}}</action>

  <action>Check if {{tracking_path}} exists:</action>
  ```bash
  test -d "{{tracking_path}}"
  ```

  <check if="tracking path does not exist">
    <action>Display error:</action>
    ```
    NO CONTEXT FOUND

    No saved context for branch: {{restore_branch}}

    Expected location: {{tracking_path}}

    Available contexts:
    {{list directories in base_path/.bmad-tracking/}}

    Options:
    - Run nook-sync from a worktree with content
    - Check if branch name is correct
    - Start fresh without restoring
    ```
    <action>HALT workflow</action>
  </check>

  <action>Load context.yaml if exists:</action>
  ```bash
  cat "{{tracking_path}}/context.yaml"
  ```
  <action>Extract and display: synced_at, synced_by, artifacts info</action>
</step>

<step n="4" goal="Load sync_paths from config">
  <action>Read sync_paths from {module_config}:</action>
  ```yaml
  sync_paths:
    - docs
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml
  ```
  <action>Store as {{sync_paths}} array</action>
</step>

<step n="5" goal="Check existing local content">
  <action>For each path in {{sync_paths}}:</action>

  <action>Check if path exists locally:
    - Check if path exists
    - If directory: count files
    - If file: check if exists
    - Flag as "has content" or "empty/missing"
  </action>

  <action>Store findings as {{existing_content_report}}</action>
</step>

<step n="6" goal="Show restore preview and warn about overwrites">
  <action>Display restore preview:</action>

  ```
  Restore Preview

  Source:  {{tracking_path}}
  Target:  {{current_path}}
  Branch:  {{restore_branch}} -> {{current_branch}}

  Paths to restore:
  {{for each path in sync_paths}}
  {{if path exists in tracking}}
  [RESTORE] {{tracking_path}}/{{path}} -> {{path}}
  {{else}}
  [SKIP]    {{path}} (not in saved context)
  {{end if}}
  {{end for}}
  ```

  <check if="any target path has existing content">
    <action>Display warning:</action>
    ```
    WARNING: Existing content will be REPLACED!

    Paths with existing content:
    {{list paths with details}}

    This operation will DELETE existing files and replace with saved versions.
    ```
  </check>

  <action>Show context metadata if available (synced_at, synced_by)</action>
</step>

<step n="7" goal="Confirm operation">
  <ask>Proceed with restore? This will REPLACE existing content. [y/n]</ask>
  <action if="user declines">Abort workflow with message: "Restore cancelled by user"</action>
</step>

<step n="8" goal="Execute restore operations">
  <action>For each path in {{sync_paths}} that exists in tracking:</action>

  <action>Determine if path is file or directory</action>

  <action>For directories (if exists in tracking):
  ```bash
  mkdir -p "{{current_path}}/{{path}}"
  rsync -av --delete "{{tracking_path}}/{{path}}/" "{{current_path}}/{{path}}/"
  ```
  </action>

  <action>For files (if exists in tracking):
  ```bash
  mkdir -p "$(dirname "{{current_path}}/{{path}}")"
  cp "{{tracking_path}}/{{path}}" "{{current_path}}/{{path}}"
  ```
  </action>

  <action>Record which paths were restored and file counts</action>
</step>

<step n="9" goal="Report results">
  <action>Display completion summary:</action>

  ```
  Restore Complete!

  Source:  {{tracking_path}}
  Target:  {{current_path}}
  Branch:  {{restore_branch}} -> {{current_branch}}

  Restored:
  {{for each path in sync_paths}}
  {{if restored}}
  [DONE] {{path}} ({{file_count}} files)
  {{else}}
  [SKIP] {{path}} (not in saved context)
  {{end if}}
  {{end for}}

  Context from {{synced_at}} has been restored.

  REMEMBER: Restored files won't show in git status

  - docs/ - gitignored (doesn't travel with git)
  - .bmad/_cfg/agents/, .bmad/bmm/config.yaml - skip-worktree
    (tracked in git but local changes ignored)

  Each worktree manages its own copy of these artifacts.

  Next Steps:
  - Continue working with the restored context
  - Run /bmad:th:workflows:nook-sync to save changes back
  ```
</step>

</workflow>

## Notes

<notes>
**Architecture:**
- All context is loaded FROM base_workspace_path/.bmad-tracking/
- This workflow copies FROM base workspace TO current directory
- Can be run from any worktree (typically nooks)
- Supports restoring from any branch, not just current

**Two-Tier Artifact System:**

| Type | Paths | Nook State | Restore Behavior |
|------|-------|------------|------------------|
| **Gitignored** | `docs/` | Don't exist | Restored from tracking |
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git | Overwritten with tracked version |

**Common Use Cases:**
1. Nook needs docs (most common):
   - Nook already has config files (skip-worktree)
   - Run nook-restore to get docs from parent
   - Work with full context

2. Reset configs to tracked version:
   - Local config changes gone wrong?
   - nook-restore overwrites with saved version
   - Skip-worktree keeps new version from showing in status

3. Return to base workspace after nook:
   - cd to base workspace
   - Run nook-restore to get latest synced context

4. Restore specific branch's context:
   - Provide source_branch parameter
   - Or select from list when prompted

**Sync Paths:**
- Configured in .bmad/th/config.yaml under sync_paths
- Includes both gitignore_paths and skip_worktree_paths
- Directories are restored with rsync --delete (mirror copy)
- Files are copied directly
- Non-existent paths in tracking are skipped

**Artifact Visibility After Restore:**
- `docs/` - gitignored, won't show in git status
- `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` - skip-worktree, won't show in git status
- Both types are "invisible" to git but for different reasons
</notes>
