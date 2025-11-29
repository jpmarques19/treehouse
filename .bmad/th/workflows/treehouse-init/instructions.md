# Treehouse Init - Initialize Base Workspace

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/treehouse-init/workflow.yaml</critical>

## Purpose

Initialize this worktree as the **base workspace** - the source of truth for all workspace tracking.
This workflow:
1. Requires a clean git state (no uncommitted changes)
2. Adds sync artifacts to .gitignore
3. Commits the .gitignore changes
4. Creates the tracking folder as a git submodule with initial context.yaml
5. Updates th config with base workspace path

<workflow>

<step n="1" goal="Check for clean git state (CRITICAL)">
  <action>Check git status for uncommitted changes:
  ```bash
  git status --porcelain
  ```
  </action>

  <check if="git status has output (uncommitted changes exist)">
    <action>Display critical halt:</action>
    ```
    CANNOT INITIALIZE - UNCOMMITTED CHANGES

    treehouse-init requires a clean git working tree.

    You have uncommitted changes:
    {{list first 10 files from git status}}

    Please commit or stash your changes first:

      git add -A && git commit -m "your message"
      # OR
      git stash

    Then run /bmad:th:workflows:treehouse-init again.
    ```
    <action>HALT workflow - do not proceed</action>
  </check>

  <check if="git status is clean">
    <action>Display: "Git working tree is clean"</action>
    <action>Continue to next step</action>
  </check>
</step>

<step n="2" goal="Check if already initialized">
  <action>Read th config to check if base_workspace_path is set:
  ```bash
  grep "^base_workspace_path:" .bmad/th/config.yaml | grep -v '""'
  ```
  </action>

  <check if="base_workspace_path is not empty">
    <action>Display:</action>
    ```
    WORKSPACE ALREADY INITIALIZED

    This workspace tree was already initialized.
    Base workspace: {{base_workspace_path}}
    Base branch:  {{base_branch}}

    Options:
    [R] Reinitialize - Reset as new base workspace (use if base was deleted)
    [C] Cancel       - Exit without changes

    Choose [R/C]:
    ```

    <action>Wait for user input: {{user_choice}}</action>

    <check if="user_choice == C or user_choice == cancel">
      <action>Display: "Cancelled. No changes made."</action>
      <action>Exit workflow</action>
    </check>

    <check if="user_choice == R or user_choice == reinitialize">
      <action>Display: "Proceeding with reinitialization..."</action>
      <action>Continue to next step</action>
    </check>
  </check>

  <check if="base_workspace_path is empty">
    <action>Display: "Initializing new base workspace..."</action>
    <action>Continue to next step</action>
  </check>
</step>

<step n="3" goal="Get current branch and workspace path">
  <action>Get current branch: `git branch --show-current`</action>
  <action>Store as {{current_branch}}</action>

  <action>Get absolute path to current directory: `pwd`</action>
  <action>Store as {{workspace_path}}</action>

  <action>Display:</action>
  ```
  Base Workspace Setup
  Path:   {{workspace_path}}
  Branch: {{current_branch}}
  ```
</step>

<step n="4" goal="Load path configurations from config">
  <action>Read path configs from {module_config}:</action>
  ```yaml
  gitignore_paths:
    - docs

  skip_worktree_paths:
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml

  sync_paths:
    - docs
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml
  ```
  <action>Store {{gitignore_paths}} - paths to add to .gitignore</action>
  <action>Store {{skip_worktree_paths}} - paths to mark with skip-worktree in nooks</action>
  <action>Store {{sync_paths}} - all paths for sync operations</action>
  <action>Also store tracking_folder (default: .bmad-tracking)</action>
</step>

<step n="5" goal="Update .gitignore with gitignore_paths only">
  <action>Check if .gitignore exists, create if not</action>

  <action>Check if Treehouse section already exists in .gitignore:
  ```bash
  grep -q "# Treehouse" .gitignore 2>/dev/null
  ```
  </action>

  <check if="Treehouse section does not exist">
    <action>Append to .gitignore (ONLY gitignore_paths, NOT skip_worktree_paths):</action>
    ```gitignore

    # Treehouse (th) - Tracked artifacts
    # These paths are per-workspace and synced via th workflows
    # Do not remove - managed by /bmad:th:workflows:treehouse-init
    docs/
    ```

    <action>Note: Only add paths from {{gitignore_paths}}, adding trailing / for directories</action>
    <action>Note: {{skip_worktree_paths}} are NOT gitignored - they use skip-worktree in nooks instead</action>
    <action>Note: .bmad-tracking/ is NOT gitignored - it's tracked and sparse-checkout excludes it from nooks</action>
  </check>

  <check if="Treehouse section already exists">
    <action>Display: ".gitignore already has Treehouse section"</action>
  </check>

  <action>Display what was added to .gitignore</action>
</step>

<step n="5a" goal="Remove gitignore_paths from git cache (skip_worktree_paths stay tracked)">
  <action>Display:</action>
  ```
  Removing gitignored paths from git cache...

  Adding to .gitignore alone doesn't stop tracking files that are
  already in the git index. We must remove them from cache.

  Note: skip_worktree_paths (.bmad/_cfg/agents, .bmad/bmm/config.yaml)
  stay tracked - they use skip-worktree in nooks instead.
  ```

  <action>For each gitignore_path ONLY, remove from git cache if tracked:</action>
  ```bash
  # Only remove gitignore_paths (docs), NOT skip_worktree_paths
  git rm -r --cached docs/ 2>/dev/null || true
  ```

  <action>Note: The `|| true` ensures the command doesn't fail if the path isn't tracked</action>
  <action>Note: `--cached` removes from index only, preserving local files</action>
  <action>IMPORTANT: Do NOT remove skip_worktree_paths from cache - they must stay tracked!</action>

  <action>Report what was removed:</action>
  ```bash
  # Check if anything was unstaged
  git status --porcelain | grep "^D " | head -10
  ```

  <check if="files were removed from cache">
    <action>Display:</action>
    ```
    Removed from git tracking (files preserved locally):
    {{list removed files}}
    ```
  </check>

  <check if="no files were removed">
    <action>Display: "No tracked files to remove (already clean)"</action>
  </check>
</step>

<step n="6" goal="Update th config with base workspace info">
  <action>Update .bmad/th/config.yaml:</action>

  <action>Set base_workspace_path to {{workspace_path}}</action>
  <action>Set base_branch to {{current_branch}}</action>

  <action>Edit the config file to update both fields:
  - Replace `base_workspace_path: ""` with `base_workspace_path: "{{workspace_path}}"`
  - Replace `base_branch: ""` with `base_branch: "{{current_branch}}"`
  </action>
</step>

<step n="7" goal="Commit .gitignore, config, and cache removal">
  <action>Stage all changes (including the cache removals):
  ```bash
  git add -A .gitignore .bmad/th/config.yaml
  ```
  </action>

  <action>Note: `git add -A` will also stage the deletions from `git rm --cached`</action>

  <action>Commit:
  ```bash
  git commit -m "chore(th): initialize base workspace tracking

  - Add gitignore_paths to .gitignore (docs/)
  - Remove gitignore_paths from git tracking (git rm --cached)
  - Set base_workspace_path in th config
  - Base branch: {{current_branch}}

  Two-tier artifact tracking:
  - Gitignored (don't travel with git): docs/
  - Skip-worktree in nooks (tracked, local changes ignored):
    .bmad/_cfg/agents/, .bmad/bmm/config.yaml

  Generated with Claude Code"
  ```
  </action>

  <action>Display: "Changes committed"</action>
</step>

<step n="8" goal="Create tracking submodule and initial context">
  <action>Display:</action>
  ```
  Creating tracking folder as git submodule...

  The tracking folder will be a standalone git repository registered
  as a submodule. This approach:
  - Keeps tracking data isolated from nooks (submodules don't auto-initialize in worktrees)
  - No sparse-checkout or complex git hacks needed
  - Local-only (no remote repository required)
  ```

  <action>Create .bmad-tracking as standalone git repository:</action>
  ```bash
  mkdir -p .bmad-tracking/{{current_branch}}
  cd .bmad-tracking
  git init
  ```

  <action>Create initial context.yaml at .bmad-tracking/{{current_branch}}/context.yaml:</action>
  ```yaml
  # Context Manifest - Base Workspace
  # Auto-generated by treehouse-init
  branch: {{current_branch}}
  type: base
  initialized_at: {{date}}
  initialized_by: {user_name}

  # This is the base workspace - source of truth for all nooks
  is_base: true

  # Lineage starts here
  lineage:
    - branch: {{current_branch}}
      type: base
      name: "{{current_branch}}"
      created: {{date}}

  # Artifacts (none synced yet)
  artifacts:
    docs:
      synced: false
    bmad_cfg_agents:
      synced: false
    bmad_bmm_config:
      synced: false

  notes: |
    Base workspace initialized.
    Run /bmad:th:workflows:nook-sync to save current artifacts.
    Run /bmad:th:workflows:nook-fork to create isolated development nooks.
  ```

  <action>Write the context.yaml file</action>

  <action>Commit tracking repo initial state:</action>
  ```bash
  cd .bmad-tracking
  git add .
  git commit -m "init: tracking repository for {{current_branch}}

  Base workspace initialized.
  Generated with Claude Code"
  cd ..
  ```

  <action>Register as submodule (local path, no remote):</action>
  ```bash
  git submodule add ./.bmad-tracking .bmad-tracking
  ```

  <action>Commit submodule registration:</action>
  ```bash
  git commit -m "chore(th): add tracking submodule

  - .bmad-tracking/ is now a git submodule
  - Submodule only initializes in base workspace
  - Nooks get empty placeholder (removed via skip-worktree)
  - No sparse-checkout needed for isolation

  Generated with Claude Code"
  ```

  <action>Display: "Tracking submodule created and registered"</action>
</step>

<step n="9" goal="Report completion">
  <action>Display completion summary:</action>

  ```
  Base Workspace Initialized!

  Base Workspace
     Path:   {{workspace_path}}
     Branch: {{current_branch}}

  Gitignored (per-nook, synced via th workflows):
  {{for each path in gitignore_paths}}
     - {{path}}
  {{end for}}

  Skip-worktree in nooks (tracked, local changes ignored):
  {{for each path in skip_worktree_paths}}
     - {{path}}
  {{end for}}

  Tracking: .bmad-tracking/ (git submodule)
     - Submodule initialized here with content
     - .gitmodules tracked in parent repo
     - Nooks: submodule not initialized (empty placeholder removed)

  Next Steps:

  1. Save your current artifacts:
     /bmad:th:workflows:nook-sync

  2. Create a focused development nook:
     /bmad:th:workflows:nook-fork

  3. View all workspaces:
     /bmad:th:workflows:treehouse-list

  Architecture:
     - This worktree is the SOURCE OF TRUTH
     - .bmad-tracking/ is a git submodule (initialized here only)
     - Nooks get skip_worktree_paths from git (files exist!)
     - Nooks don't get gitignore_paths (must restore)
     - Nooks don't get .bmad-tracking/ or .gitmodules (skip-worktree + removed)
     - nook-sync saves TO here (optionally auto-commits submodule)
     - nook-restore loads FROM here
  ```
</step>

</workflow>

## Notes

<notes>
**Three-tier artifact tracking architecture:**

| Type | Paths | Behavior in Nooks |
|------|-------|-------------------|
| **Gitignored** | `docs/` | Don't exist - must restore |
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git, local changes ignored |
| **Submodule** | `.bmad-tracking/`, `.gitmodules` | Skip-worktree applied, then removed |

**Why gitignore for docs/:**
- Large directories that would bloat git history
- Simpler mental model - gitignored files don't travel with git
- User explicitly syncs when ready

**Why skip-worktree for config files:**
- Small files that should exist immediately in nooks
- `git update-index --skip-worktree` marks files as "assume unchanged"
- Files exist in worktree (from git) but local changes don't show in git status
- Each nook can customize configs without committing changes
- nook-fork applies skip-worktree automatically

**Why git rm --cached is required for gitignore_paths:**
- Adding to .gitignore only affects NEW files
- Already-tracked files continue to be tracked even if in .gitignore
- `git rm --cached` removes files from the index WITHOUT deleting them locally
- After this, the files are truly untracked and won't appear in git status

**Why tracking folder is a git submodule:**
- Submodules don't auto-initialize in worktrees - nooks get empty placeholder
- No sparse-checkout complexity needed
- Skip-worktree on submodule pointer hides deletion from git status
- Local-only (no remote repository required)
- Two commits for sync: submodule content + parent pointer update
- Standard git behavior - easy to understand and maintain

**Why .gitmodules gets skip-worktree in nooks:**
- .gitmodules is a tracked file that references the submodule
- Nooks don't need it (they don't use the submodule)
- Skip-worktree + removal keeps nooks clean

**Base workspace responsibilities:**
- Holds .bmad-tracking/ submodule (source of truth)
- All workspace context is archived here
- Should not be deleted while nooks exist

**Re-initialization:**
- Safe to reinitialize if base was accidentally deleted
- Will reset base_workspace_path to current location
- Existing tracking data (if any) is preserved
</notes>
