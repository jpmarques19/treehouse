# Nook Fork - Create Isolated Development Nook

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/nook-fork/workflow.yaml</critical>

## Purpose

Create an **isolated development nook** - a new git worktree for focused work:
- Nook has config files from git (skip-worktree applied - local changes ignored)
- Nook does NOT have docs/ (gitignored - must restore if needed)
- Create focused workspace for bugs, spikes, or exploration
- Lineage tracked via branch name prefix (hash) and context.yaml

**Key Architecture Points:**
- Source of truth: `base_workspace_path/.bmad-tracking/` (git submodule)
- Three-tier artifact handling:
  - **skip_worktree_paths** (`.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml`): Exist from git, skip-worktree applied
  - **submodule_paths** (`.bmad-tracking/`, `.gitmodules`): Skip-worktree applied, then removed from nook
  - **gitignore_paths** (`docs/`): Don't exist, must restore if needed
- User decides when to sync (no forced auto-sync)
- To get docs in nook: run `nook-restore` after creating nook

## Branch Naming Scheme

Uses **hash-based lineage** to avoid git ref conflicts while maintaining traceability:

```
{type}/{parent-hash-4-chars}-{short-name}
```

**Examples:**
```
peek-a-box-mvp                    <- Base workspace (hash: a1b2)
explore/a1b2-provision-race       <- Explore nook from a1b2 (hash: 7f3e)
spike/7f3e-mutex-solution         <- Spike nook from 7f3e (hash: c9d1)
bugfix/c9d1-deadlock-retry        <- Bugfix nook from c9d1 (hash: jf2b)
```

<workflow>

<step n="1" goal="Check if workspace is initialized">
  <action>Read {module_config} and check base_workspace_path:</action>
  ```bash
  grep "^base_workspace_path:" .bmad/th/config.yaml
  ```

  <check if="base_workspace_path is empty or not set">
    <action>Display critical halt:</action>
    ```
    CANNOT CREATE NOOK - WORKSPACE NOT INITIALIZED

    This workspace hasn't been initialized as a base workspace.
    The base workspace holds the tracking folder (.bmad-tracking/)
    which is the source of truth for all nooks.

    Please run first:
      /bmad:th:workflows:treehouse-init

    This will:
    1. Add sync artifacts to .gitignore
    2. Create the tracking folder
    3. Establish this worktree as the base workspace
    ```
    <action>HALT workflow - do not proceed</action>
  </check>

  <action>Store base_workspace_path as {{base_path}}</action>
  <action>Store base_branch from config as {{base_branch}}</action>
  <action>Display: "Base workspace: {{base_path}} ({{base_branch}})"</action>
</step>

<step n="2" goal="Get current branch info and validate tracking">
  <action>Get current branch: `git branch --show-current`</action>
  <action>Store as {{current_branch}}</action>

  <action>Get short hash of current HEAD (4 chars): `git rev-parse --short=4 HEAD`</action>
  <action>Store as {{parent_hash}}</action>

  <action>Check if tracking folder has uncommitted changes:</action>
  ```bash
  git -C "{{base_path}}" status --porcelain .bmad-tracking/
  ```

  <check if="tracking folder has modifications (non-empty output)">
    <action>Display warning with options:</action>
    ```
    WARNING - TRACKING FOLDER HAS MODIFICATIONS

    The tracking folder at {{base_path}}/.bmad-tracking/
    has uncommitted changes. This may indicate:
    - A sync operation was interrupted
    - Manual edits to context files
    - Work in progress that should be committed

    Uncommitted changes:
    {{list first 10 files from git status}}

    RECOMMENDATION: Commit these changes before creating nook
    to ensure a clean baseline for the new nook.

    Options:
    [C] Commit changes now - Stage and commit tracking folder
    [P] Proceed anyway    - Create nook without committing (not recommended)
    [A] Abort             - Exit and handle manually

    Choose [C/P/A]:
    ```

    <check if="user_choice == C or user_choice == commit">
      <action>Stage and commit tracking changes:</action>
      ```bash
      cd "{{base_path}}"
      git add .bmad-tracking/
      git commit -m "chore(th): commit tracking changes before nook

      Generated with Claude Code"
      ```
      <action>Display: "Tracking changes committed. Continuing..."</action>
    </check>

    <check if="user_choice == P or user_choice == proceed">
      <action>Display: "Proceeding with uncommitted tracking changes..."</action>
    </check>

    <check if="user_choice == A or user_choice == abort">
      <action>Display: "Aborted. Please handle tracking changes and try again."</action>
      <action>HALT workflow</action>
    </check>
  </check>

  <check if="tracking folder is clean">
    <action>Display: "Tracking folder is clean"</action>
  </check>
</step>

<step n="3" goal="Check if current branch has context synced">
  <action>Check if context.yaml exists for current branch:</action>
  ```bash
  test -f "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml"
  ```

  <check if="context.yaml does NOT exist for current branch">
    <action>Display informational notice:</action>
    ```
    NOTE - NO SYNCED CONTEXT FOR CURRENT BRANCH

    Branch '{{current_branch}}' hasn't been synced yet.
    This means the nook won't inherit any context lineage.

    This is fine if:
    - This is a new workspace just starting out
    - You haven't made changes worth preserving yet

    If you want to save context first:
      /bmad:th:workflows:nook-sync

    Continuing with nook creation...
    ```
  </check>

  <check if="context.yaml exists">
    <action>Load existing context for lineage:</action>
    ```bash
    cat "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml"
    ```
    <action>Extract {{existing_lineage}} array from context</action>
  </check>
</step>

<step n="4" goal="Select nook type">
  <check if="nook_type provided as input">
    <action>Validate nook_type is one of: explore, spike, bugfix, discovery, feature, experiment, hotfix</action>
    <action>Store as {{nook_type}}</action>
  </check>

  <check if="nook_type NOT provided">
    <ask>What kind of nook do you want to create?

    1. **explore/** - Deep exploration or investigation
    2. **spike/** - Quick prototype or proof of concept
    3. **bugfix/** - Focused bug investigation and fix
    4. **discovery/** - Requirements or architecture discovery
    5. **feature/** - New feature branch
    6. **experiment/** - Wild experimentation
    7. **hotfix/** - Production hotfix
    8. **custom** - Custom branch prefix

    Choose [1-8]:</ask>

    <action>Map choice to type:
      1 -> explore
      2 -> spike
      3 -> bugfix
      4 -> discovery
      5 -> feature
      6 -> experiment
      7 -> hotfix
      8 -> ask for custom type (validate: lowercase, no slashes)
    </action>
    <action>Store as {{nook_type}}</action>
  </check>
</step>

<step n="5" goal="Define nook name">
  <check if="nook_name provided as input">
    <action>Validate nook_name is kebab-case</action>
    <action>Store as {{nook_name}}</action>
  </check>

  <check if="nook_name NOT provided">
    <ask>What should this nook be called? (kebab-case, e.g., "race-condition", "mqtt-reconnect")

    The full branch will be: {{nook_type}}/{{parent_hash}}-<your-name></ask>

    <action>Validate name:
      - kebab-case only
      - no spaces or special chars
      - no slashes
    </action>
    <action>Store as {{nook_name}}</action>
  </check>

  <action>Construct full branch name: {{nook_type}}/{{parent_hash}}-{{nook_name}}</action>
  <action>Store as {{new_branch}}</action>

  <action>Display preview:</action>
  ```
  Nook Preview

  Parent branch:  {{current_branch}} ({{parent_hash}})
  New branch:     {{new_branch}}
  Nook type:      {{nook_type}}
  ```
</step>

<step n="6" goal="Validate worktrees folder">
  <action>Read worktrees_folder from {module_config}</action>

  <check if="worktrees_folder not configured or doesn't exist">
    <ask>Where should worktrees be created? (e.g., /home/user/my-workspace-worktrees)</ask>
    <action>Validate path is writable</action>
    <action>Create directory if needed: `mkdir -p {{worktrees_folder}}`</action>
    <action>Update {module_config} with new worktrees_folder value</action>
  </check>

  <action>Store as {{worktrees_folder}}</action>
</step>

<step n="7" goal="Create the worktree">
  <action>Derive worktree folder name: {{parent_hash}}-{{nook_name}}</action>
  <action>Calculate full worktree path: {{worktrees_folder}}/{{parent_hash}}-{{nook_name}}</action>
  <action>Store as {{worktree_path}}</action>

  <action>Check path doesn't already exist:</action>
  ```bash
  test -d "{{worktree_path}}"
  ```

  <check if="path exists">
    <action>Display error:</action>
    ```
    ERROR - WORKTREE PATH EXISTS

    Path already exists: {{worktree_path}}

    Options:
    - Choose a different nook name
    - Remove the existing worktree first:
        git worktree remove "{{worktree_path}}"
    ```
    <action>HALT workflow</action>
  </check>

  <action>Display what will happen:</action>
  ```
  Creating Nook

  New branch:    {{new_branch}}
  Base:          {{current_branch}} ({{parent_hash}})
  Worktree path: {{worktree_path}}
  ```

  <ask>Proceed with nook creation? [y/n]</ask>

  <check if="user confirms">
    <action>Execute worktree creation:</action>
    ```bash
    git worktree add "{{worktree_path}}" -b "{{new_branch}}" "{{current_branch}}"
    ```
  </check>

  <check if="user declines">
    <action>Display: "Nook creation cancelled."</action>
    <action>HALT workflow</action>
  </check>
</step>

<step n="8" goal="Apply skip-worktree and remove submodule from nook">
  <action>Display:</action>
  ```
  Configuring nook isolation...

  Applying skip-worktree to:
  - Config files (local changes won't show in git status)
  - Submodule pointer and .gitmodules (will be removed from nook)
  ```

  <action>Read skip_worktree_paths from {module_config}:</action>
  ```yaml
  skip_worktree_paths:
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml
  ```

  <action>Apply skip-worktree to config files:</action>
  ```bash
  cd "{{worktree_path}}"

  # For directories - apply to all files recursively
  git ls-files .bmad/_cfg/agents/ | xargs -r git update-index --skip-worktree

  # For single files
  git update-index --skip-worktree .bmad/bmm/config.yaml 2>/dev/null || true
  ```

  <action>Apply skip-worktree to submodule pointer and .gitmodules:</action>
  ```bash
  cd "{{worktree_path}}"

  # Skip-worktree on submodule pointer (tracked as gitlink in index)
  git update-index --skip-worktree .bmad-tracking 2>/dev/null || true

  # Skip-worktree on .gitmodules file
  git update-index --skip-worktree .gitmodules 2>/dev/null || true
  ```

  <action>Remove submodule placeholder and .gitmodules from nook:</action>
  ```bash
  cd "{{worktree_path}}"

  # Remove empty submodule placeholder folder
  rm -rf .bmad-tracking/

  # Remove .gitmodules file
  rm -f .gitmodules
  ```

  <action>Verify isolation and clean git status:</action>
  ```bash
  cd "{{worktree_path}}"
  echo "Skip-worktree files: $(git ls-files -v | grep "^S" | wc -l)"
  git status --short
  ```

  <check if="git status is clean (no output from status --short)">
    <action>Display:</action>
    ```
    Nook isolation configured:
      - Config files: skip-worktree applied (local changes ignored)
      - Submodule (.bmad-tracking/): removed from nook
      - .gitmodules: removed from nook
      - Git status: clean
    ```
  </check>

  <check if="git status shows changes">
    <action>Display warning:</action>
    ```
    Nook isolation incomplete - git status shows uncommitted changes.
    This may indicate skip-worktree didn't apply correctly.

    Try manually:
      git update-index --skip-worktree .bmad-tracking
      git update-index --skip-worktree .gitmodules
    ```
  </check>
</step>

<step n="9" goal="Create context.yaml for new nook in BASE workspace">
  <action>Build lineage for new nook:</action>

  <action>If parent has existing_lineage, use it as base</action>
  <action>Add parent entry if not already in lineage:</action>
  ```yaml
  parent_entry:
    branch: {{current_branch}}
    hash: {{parent_hash}}
    type: (extract from branch name or "base")
    name: (derive human name from branch)
  ```

  <action>Create tracking directory for new nook in base workspace:</action>
  ```bash
  mkdir -p "{{base_path}}/.bmad-tracking/{{new_branch}}"
  ```

  <action>Create context.yaml at {{base_path}}/.bmad-tracking/{{new_branch}}/context.yaml:</action>
  ```yaml
  # Context Manifest - Nook
  # Auto-generated by nook-fork
  branch: {{new_branch}}
  type: {{nook_type}}
  forked_from: {{current_branch}}
  forked_at: {{date}}
  forked_by: {user_name}

  # Parent reference
  parent:
    branch: {{current_branch}}
    hash: {{parent_hash}}

  # Full lineage chain (oldest to newest, ending with this nook)
  lineage:
    {{for each entry in existing_lineage}}
    - branch: {{entry.branch}}
      hash: {{entry.hash}}
      type: {{entry.type}}
      name: "{{entry.name}}"
    {{end for}}
    - branch: {{new_branch}}
      hash: pending
      type: {{nook_type}}
      name: "{{nook_name}}"
      parent_hash: {{parent_hash}}
      created: {{date}}

  # Artifacts (none synced yet - nook starts clean)
  artifacts:
    docs:
      synced: false
    bmad_cfg_agents:
      synced: false
    bmad_bmm_config:
      synced: false

  notes: |
    Nook created from {{current_branch}}.
    Nook starts clean - no artifacts inherited.
    Run /bmad:th:workflows:nook-restore to load artifacts from parent.
    Run /bmad:th:workflows:nook-sync after making changes to save.
  ```

  <action>Write the context.yaml file</action>
</step>

<step n="10" goal="Report completion and next steps">
  <action>Calculate lineage depth</action>

  <action>Display completion summary:</action>

  ```
  Nook Created Successfully!

  Parent Workspace
     Branch: {{current_branch}}
     Hash:   {{parent_hash}}
     Path:   (current directory or base_path)

  New Nook
     Branch:   {{new_branch}}
     Type:     {{nook_type}}
     Location: {{worktree_path}}

  Lineage ({{lineage_count}} levels):
     {{for each entry in lineage}}
     {{index}}. {{entry.branch}} ({{entry.hash}}) - {{entry.type}}
     {{end for}}
     -> {{new_branch}} (new) - {{nook_type}} <- YOU ARE HERE

  What's in your nook:
     Config files exist (.bmad/_cfg/agents/, .bmad/bmm/config.yaml)
       - From git, with skip-worktree (local changes won't show in status)
     docs/ does NOT exist (gitignored)
       - Run nook-restore if you need docs
     .bmad-tracking/ does NOT exist (submodule removed)
     .gitmodules does NOT exist (removed)

  Next Steps:

  1. Switch to the new worktree:
     cd {{worktree_path}}

  2. Config files are ready to use!
     They exist from git with skip-worktree applied.
     You can modify them freely - changes stay local.

  3. To load docs from parent context:
     /bmad:th:workflows:nook-restore

  4. When done with the nook:
     - Commit your code changes
     - Run /bmad:th:workflows:nook-sync to save artifacts
     - Merge back: git checkout {{current_branch}} && git merge {{new_branch}}
     - Remove worktree: git worktree remove "{{worktree_path}}"

  5. To return to parent workspace:
     cd {{base_path}}

  Start a new AI session in the nook for maximum context clarity!
  ```
</step>

</workflow>

## Branch Naming Reference

<notes>
**Type Folders:**
| Type | Purpose |
|------|---------|
| `discovery/` | New workspace/epic discovery |
| `explore/` | Deep investigation |
| `spike/` | Quick prototype/POC |
| `bugfix/` | Bug investigation & fix |
| `feature/` | Feature development |
| `experiment/` | Wild experimentation |
| `hotfix/` | Production hotfix |

**Why Hash-Based Naming:**
- Git refs are files in `.git/refs/heads/`
- A path cannot be both a file and directory
- `explore/` is a directory, `explore/a1b2-foo` is a file (branch)
- No conflicts, infinite depth via hash chaining

**Lineage Tracing:**
- 4-char hash prefix links to parent commit
- Full chain stored in context.yaml in base_workspace_path/.bmad-tracking/
- Can trace any branch back to its root

**Three-Tier Nook Architecture:**

| Type | Paths | Behavior in Nook |
|------|-------|------------------|
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git, local changes ignored |
| **Submodule (skip-worktree + removed)** | `.bmad-tracking/`, `.gitmodules` | Skip-worktree applied, then removed |
| **Gitignored** | `docs/` | Don't exist, must restore |

**Why skip-worktree for configs:**
- Files exist immediately in nook (no restore needed)
- `git update-index --skip-worktree` marks files as "assume unchanged"
- Local changes don't show in git status
- Each nook can customize configs independently
- Changes are saved via nook-sync, not git commit

**Why submodule for tracking folder:**
- Submodules don't auto-initialize in worktrees - nooks get empty placeholder
- Skip-worktree on submodule pointer hides deletion from git status
- No sparse-checkout complexity needed
- `.gitmodules` also gets skip-worktree + removal for clean nooks
- All tracking goes through base workspace via absolute paths

**Why gitignore for docs:**
- Large directories would bloat git history
- User explicitly restores if needed
- Single source of truth in base workspace tracking folder
</notes>
