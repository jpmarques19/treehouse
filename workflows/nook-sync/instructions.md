# Nook Sync - Save Context to Base Workspace

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/nook-sync/workflow.yaml</critical>

## Purpose

Save working context (docs, configs, etc.) **TO the base workspace's tracking folder**.
This is the single source of truth for all workspace context.

**Key Architecture Points:**
- All context is saved to `base_workspace_path/.bmad-tracking/{branch}/` (git submodule)
- Can be run from ANY worktree (base or nook)
- Sync artifacts from current working directory to base workspace
- Updates or creates context.yaml with lineage information
- Optional auto-commit: commits submodule changes + updates parent pointer

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

<step n="9a" goal="Generate nook summary for context analyst">
  <action>Gather information to generate a meaningful summary:</action>

  <action>1. Get git commits since nook was forked (or last 20 if many):</action>
  ```bash
  # Find fork point from parent branch (if known from existing context)
  git log --oneline -20 --no-merges
  ```

  <action>2. Get list of files changed in this nook:</action>
  ```bash
  git diff --name-only HEAD~10 2>/dev/null || git diff --name-only --cached
  ```

  <action>3. If docs/ was synced, scan doc titles/headers for context:</action>
  ```bash
  # Read first few lines of each doc for context
  head -20 docs/*.md 2>/dev/null || true
  ```

  <action>4. Check for any existing summary in parent's context.yaml (if this is an update)</action>

  <action>5. Analyze the gathered information and generate summary:</action>

  <critical>
  Generate a MEANINGFUL summary based on actual work done in this nook.
  Do NOT use placeholder text. Analyze:
  - Commit messages reveal PURPOSE and DECISIONS
  - Changed files reveal SCOPE of work
  - Doc content reveals INSIGHTS and FINDINGS
  - Unfinished work or TODOs reveal NEXT STEPS
  - Error mentions or "fix" commits may reveal BLOCKERS
  </critical>

  <action>Construct {{nook_summary}} object:</action>
  ```yaml
  summary:
    purpose: |
      [1-2 sentences: WHY this nook exists, derived from branch name + commits]

    insights:
      [List key findings/learnings from docs and commits, or empty if none]

    decisions:
      [List of {decision, rationale} pairs from commits/docs, or empty if none]

    status: |
      [Current state: what's done, what's in progress]

    blockers:
      [Known issues or blockers discovered, or empty if none]

    next_steps:
      [Recommended actions based on TODOs, incomplete work, or logical next steps]
  ```

  <note>
  If this is a fresh nook with minimal commits, generate a minimal but accurate summary.
  It's better to have "purpose: Exploring X" with empty insights than fabricated content.
  </note>
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

<step n="11" goal="Generate context manifest with lineage and summary">
  <action>Create or update `{{tracking_path}}/context.yaml` with:</action>

  ```yaml
  # Context Manifest - Auto-generated by nook-sync
  branch: {{current_branch}}
  hash: {{current_hash}}
  type: {{branch_type}}
  synced_at: {{date}}
  synced_by: {user_name}
  synced_from: {{current_path}}

  # AI-generated nook summary (for nook-context-analyst)
  # This section provides semantic context about the nook's work
  summary:
    purpose: |
      {{nook_summary.purpose}}

    insights:
      {{for each insight in nook_summary.insights}}
      - {{insight}}
      {{end for}}

    decisions:
      {{for each decision in nook_summary.decisions}}
      - decision: {{decision.decision}}
        rationale: {{decision.rationale}}
      {{end for}}

    status: |
      {{nook_summary.status}}

    blockers:
      {{for each blocker in nook_summary.blockers}}
      - {{blocker}}
      {{end for}}

    next_steps:
      {{for each step in nook_summary.next_steps}}
      - {{step}}
      {{end for}}

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

<step n="11a" goal="Optional auto-commit for submodule">
  <action>Read auto_commit_tracking from {module_config}:</action>
  ```bash
  grep "^auto_commit_tracking:" .bmad/th/config.yaml | awk '{print $2}'
  ```
  <action>Store as {{auto_commit}} (default: false)</action>

  <check if="auto_commit == true">
    <action>Display:</action>
    ```
    Auto-commit enabled - committing tracking changes...

    With submodule architecture, this requires TWO commits:
    1. Inside .bmad-tracking/ submodule (the actual file changes)
    2. In parent repo (update submodule pointer)
    ```

    <action>Commit changes inside tracking submodule:</action>
    ```bash
    cd "{{base_path}}/.bmad-tracking"
    git add .
    git commit -m "sync: {{current_branch}} context

    Synced from: {{current_path}}
    Hash: {{current_hash}}

    Generated with Claude Code" || true
    ```

    <action>Update parent repo's submodule pointer:</action>
    ```bash
    cd "{{base_path}}"
    git add .bmad-tracking
    git commit -m "chore(th): update tracking pointer for {{current_branch}}

    Generated with Claude Code" || true
    ```

    <action>Display: "Auto-commit complete (submodule + pointer)"</action>
  </check>

  <check if="auto_commit != true">
    <action>Store {{manual_commit_needed}} = true</action>
  </check>
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

  Summary Generated:
  ┌─────────────────────────────────────────────────────────
  │ Purpose: {{nook_summary.purpose | truncate(80)}}
  │ Status:  {{nook_summary.status | truncate(80)}}
  │ Insights: {{nook_summary.insights | length}} | Decisions: {{nook_summary.decisions | length}}
  │ Blockers: {{nook_summary.blockers | length}} | Next Steps: {{nook_summary.next_steps | length}}
  └─────────────────────────────────────────────────────────

  Manifest: {{tracking_path}}/context.yaml
  ```

  <check if="auto_commit == true">
    <action>Display:</action>
    ```
    Auto-commit: ENABLED
    - Submodule changes committed
    - Parent pointer updated

    Context is ready to use in other nooks!

    To restore this context in another worktree:
      /bmad:th:workflows:nook-restore
    ```
  </check>

  <check if="manual_commit_needed == true">
    <action>Display:</action>
    ```
    IMPORTANT: Manual commit required (auto_commit_tracking: false)

    The tracking folder is a git submodule. To commit changes:

    1. Go to base workspace:
       cd {{base_path}}

    2. Commit inside submodule:
       cd .bmad-tracking
       git add .
       git commit -m "sync: {{current_branch}} context"
       cd ..

    3. Update parent pointer:
       git add .bmad-tracking
       git commit -m "chore(th): update tracking pointer"

    To enable auto-commit, set in .bmad/th/config.yaml:
      auto_commit_tracking: true

    To restore this context in another worktree:
      /bmad:th:workflows:nook-restore
    ```
  </check>
</step>

</workflow>

## Notes

<notes>
**Architecture:**
- All context is stored in base_workspace_path/.bmad-tracking/ (git submodule)
- This workflow copies FROM current directory TO base workspace
- Can be run from any worktree (base or nook)
- The tracking folder only exists in the base workspace (initialized submodule)

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

**Submodule Commit Pattern:**
With the submodule architecture, committing tracking changes requires TWO commits:
1. **Inside .bmad-tracking/ submodule** - commits the actual file changes
2. **In parent repo** - commits the updated submodule pointer

This is standard git submodule behavior and provides:
- Better version control of tracking data
- Clear separation between code and tracking commits
- Ability to roll back tracking changes independently

**Auto-Commit Option:**
- Controlled by `auto_commit_tracking` in .bmad/th/config.yaml
- Default: false (manual commits required)
- When true: automatically commits both submodule and pointer
- When false: shows detailed manual commit instructions

**Commit Responsibility:**
- If auto_commit_tracking: true - workflow handles both commits
- If auto_commit_tracking: false - user must commit manually
- Manual mode allows batching multiple syncs before committing

**Nook Summary Generation:**
Each sync generates an AI-analyzed summary for the nook-context-analyst:
- **purpose**: Why this nook exists (derived from branch name + commits)
- **insights**: Key findings/learnings discovered during work
- **decisions**: Important choices made with rationale
- **status**: Current state of work (done/in-progress)
- **blockers**: Known issues or obstacles
- **next_steps**: Recommended actions for context resumption

The summary is generated by analyzing:
- Git commit messages (reveal purpose, decisions, fixes)
- Changed files (reveal scope of work)
- Doc content (reveal insights and findings)
- TODOs and incomplete work (reveal next steps)

This enables the nook-context-analyst to provide meaningful context
without reading all artifacts - useful for quick lineage understanding.
</notes>
