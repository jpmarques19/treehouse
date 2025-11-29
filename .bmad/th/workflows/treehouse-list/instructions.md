# Treehouse List - View Workspace Lineage Tree

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/treehouse-list/workflow.yaml</critical>

## Purpose

Display a visual tree of all workspaces showing:
- Nook lineage and hierarchy
- Current location in the tree
- Workspace health status (active, inactive, stale, orphan)
- Quick access to navigation and cleanup operations

**Key Architecture Points:**
- All tracking data is read from `base_workspace_path/.bmad-tracking/`
- Worktree status is determined by cross-referencing `git worktree list`
- Stale detection based on `synced_at` timestamp (configurable threshold)
- Can be run from any worktree

<workflow>

<step n="1" goal="Check if workspace is initialized and get base path">
  <action>Read {module_config} and check base_workspace_path:</action>
  ```bash
  grep "^base_workspace_path:" .bmad/th/config.yaml
  ```

  <check if="base_workspace_path is empty or not set">
    <action>Display:</action>
    ```
    🌳 Workspace Lineage

    ⚠ not initialized

    run /bmad:th:workflows:treehouse-init first

    [i] initialize  [q] quit
    ```

    <action>Wait for user input: {{init_choice}}</action>

    <check if="init_choice == I">
      <action>Execute workflow: /bmad:th:workflows:treehouse-init</action>
      <action>Exit workflow</action>
    </check>

    <check if="init_choice == Q">
      <action>Exit workflow</action>
    </check>
  </check>

  <action>Store base_workspace_path as {{base_path}}</action>
  <action>Store base_branch from config as {{base_branch}}</action>
</step>

<step n="2" goal="Get current branch info">
  <action>Get current branch: `git branch --show-current`</action>
  <action>Store as {{current_branch}}</action>
  <action>Get current hash: `git rev-parse --short=4 HEAD`</action>
  <action>Store as {{current_hash}}</action>
  <action>Get current directory: `pwd`</action>
  <action>Store as {{current_path}}</action>
</step>

<step n="3" goal="Verify base workspace tracking folder exists">
  <action>Check if tracking folder exists at base workspace:
  ```bash
  test -d "{{base_path}}/.bmad-tracking"
  ```
  </action>

  <check if="tracking folder does NOT exist">
    <action>Display:</action>
    ```
    🌳 Workspace Lineage

    ⚠ tracking folder not found

    base     {{base_path}}
    expected {{base_path}}/.bmad-tracking/

    re-run /bmad:th:workflows:treehouse-init from base
    ```
    <action>HALT workflow</action>
  </check>
</step>

<step n="4" goal="Scan tracking folder for all workspaces">
  <action>Find all context.yaml files in base workspace tracking:
  ```bash
  find "{{base_path}}/.bmad-tracking" -name "context.yaml" -type f 2>/dev/null
  ```
  </action>

  <check if="no context.yaml files found">
    <action>Display:</action>
    ```
    🌳 Workspace Lineage

    no workspaces tracked yet

    base {{base_path}}
    here {{current_branch}} · {{current_hash}}

    run /bmad:th:workflows:nook-sync to sync current context
    ```
    <action>Exit workflow</action>
  </check>

  <action>For each context.yaml found:
    - Load the file
    - Extract: branch, hash, type, forked_from, lineage, synced_at
    - Derive path from context.yaml location (relative to .bmad-tracking/)
    - Store in {{workspaces}} array
  </action>
</step>

<step n="5" goal="Scan git worktrees">
  <action>List all worktrees:
  ```bash
  git worktree list --porcelain
  ```
  </action>

  <action>Parse output to get:
    - worktree path
    - branch name
    - HEAD commit
  </action>
  <action>Store as {{worktrees}} array</action>
</step>

<step n="6" goal="Build workspace tree structure with stale detection">
  <action>Identify root workspaces (no forked_from or forked_from not in tracked workspaces)</action>

  <action>For each workspace, determine:
    - parent: from forked_from field
    - children: workspaces that have this branch as forked_from
    - depth: levels from root
  </action>

  <action>Calculate staleness for each workspace:
  ```bash
  # Extract synced_at from context.yaml
  synced_at=$(grep "synced_at:" "{{tracking_path}}/context.yaml" | cut -d: -f2- | xargs)

  # Calculate days since sync
  if [ -n "$synced_at" ]; then
    synced_epoch=$(date -d "$synced_at" +%s 2>/dev/null || echo 0)
    now_epoch=$(date +%s)
    days_old=$(( (now_epoch - synced_epoch) / 86400 ))
  else
    days_old=0
  fi
  ```
  </action>

  <action>For each tracked workspace, check if git branch exists:
  ```bash
  git show-ref --verify --quiet "refs/heads/{{branch}}" && echo "EXISTS" || echo "MISSING"
  ```
  </action>

  <action>Determine workspace status for each:
    - **active**: context.yaml exists + branch exists + worktree exists
    - **inactive**: context.yaml exists + branch exists + NO worktree
    - **orphan-tracking**: context.yaml exists + branch does NOT exist in git
    - **stale**: sync > 14 days old (applies to active and inactive)
    - **current**: is the current branch we're on (marked with ←)
  </action>

  <action>Also identify orphan worktrees:
    - Worktrees that exist but have no matching context.yaml
  </action>

  <action>Store tree structure as {{workspace_tree}}</action>
</step>

<step n="7" goal="Display workspace tree and menu">
  <critical>Do NOT show internal processing steps to the user. Only display the final tree output.</critical>

  <action>Build tree silently, then display ONLY this output:</action>
  ```
  🌳 Workspace Lineage

  base {{base_path}}
  here {{current_branch}} · {{current_hash}}

  {{tree output - see format below}}

  {{if orphan_worktrees}}
  ? orphan worktrees
    {{for each orphan}}
    {{orphan.path}} → {{orphan.branch}}
    {{end for}}
  {{end if}}

  {{total}} tracked  ● {{active}} active  ○ {{inactive}} inactive
  ! {{stale}} stale  ✗ {{orphan_tracking}} orphan  ? {{orphan_wt}} orphan-wt
  ─────────────────────────────────────────────────
  n navigate    s sleep    d delete
  p prune       r refresh  q quit
  ```

  <action>Status icons:</action>
  ```
  ● = active (has worktree + branch exists)
  ○ = inactive (no worktree, branch exists)
  ✗ = orphan-tracking (context exists, but git branch deleted)
  ! = stale (>14 days since sync)
  ? = orphan-worktree (worktree exists, no tracking)
  ← = you are here
  ```

  <action>Example complete output:</action>
  ```
  🌳 Workspace Lineage

  base /home/joao/Documents/amora-os
  here discovery/peek-a-box-mvp · 7bee

  ● main a1b2                                      2 days ago
  │
  ├─ ● discovery/peek-a-box-mvp 7bee  ←            today
  │  │
  │  ├─ ○ explore/7bee-provision-race c3d4         5 days ago
  │  │  │
  │  │  └─ ! spike/c3d4-mutex-fix e5f6             18 days ago
  │  │
  │  └─ ○ bugfix/7bee-mqtt-reconnect g7h8          7 days ago
  │
  └─ ○ feature/a1b2-auth-system i9j0               7 days ago

  5 tracked  ● 2 active  ○ 3 inactive
  ! 1 stale  ✗ 0 orphan  ? 0 orphan-wt
  ─────────────────────────────────────────────────
  n navigate    s sleep    d delete
  p prune       r refresh  q quit
  ```

  <action>Wait for user input: {{menu_choice}}</action>
</step>

<step n="8" goal="Handle menu selection">
  <check if="menu_choice == N or menu_choice == navigate">
    <action>Show numbered list of branches:</action>
    ```
    select branch:
    {{for each workspace with index}}
    {{index}}. {{status}} {{workspace.branch}}
    {{end for}}

    enter number or q to cancel:
    ```
    <action>Wait for selection</action>

    <check if="selected branch has active worktree">
      <action>Display:</action>
      ```
      → cd {{worktree_path}}
      ```
    </check>

    <check if="selected branch is synced only (no worktree)">
      <action>Display:</action>
      ```
      ○ {{branch_name}} has no worktree

      to reactivate:
        git worktree add "{{suggested_path}}" "{{branch_name}}"
        cd "{{suggested_path}}"
        /bmad:th:workflows:nook-restore
      ```
    </check>
  </check>

  <check if="menu_choice == S or menu_choice == sleep">
    <action>Show numbered list of ACTIVE branches only:</action>
    ```
    select branch to sleep (removes worktree, keeps tracking + branch):
    {{for each ACTIVE workspace with index}}
    {{index}}. ● {{workspace.branch}}
    {{end for}}

    enter number or q to cancel:
    ```
    <action>Wait for selection: {{sleep_index}}</action>
    <check if="valid selection">
      <check if="selected branch is current branch">
        <action>Display:</action>
        ```
        ⚠ cannot sleep current branch

        switch to another worktree first
        ```
      </check>
      <check if="selected branch is NOT current branch">
        <action>Display confirmation:</action>
        ```
        sleep {{selected_branch}}?

        this will:
          ✓ remove worktree folder
          ✓ keep git branch
          ✓ keep tracking data

        to reactivate later:
          git worktree add "{{worktree_path}}" "{{branch}}"
          /bmad:th:workflows:nook-restore

        proceed? [y/n]
        ```
        <check if="confirmed">
          <action>Remove worktree: `git worktree remove "{{worktree_path}}"`</action>
          <action>Display:</action>
          ```
          ○ {{selected_branch}} is now inactive
          ```
          <action>Go back to Step 4 (rescan) - display tree again</action>
        </check>
      </check>
    </check>
  </check>

  <check if="menu_choice == D or menu_choice == delete">
    <action>Show numbered list of branches:</action>
    ```
    select branch to delete:
    {{for each workspace with index}}
    {{index}}. {{status}} {{workspace.branch}}
    {{end for}}

    enter number or q to cancel:
    ```
    <action>Wait for selection: {{delete_index}}</action>
    <check if="valid selection">
      <action>Display confirmation:</action>
      ```
      delete {{selected_branch}}? [y/n]
      ```
      <check if="confirmed">
        <action>Execute task: delete-nook with target={{selected_branch}}</action>
      </check>
    </check>
  </check>

  <check if="menu_choice == P or menu_choice == prune">
    <check if="orphan_count == 0 AND stale_count == 0">
      <action>Display:</action>
      ```
      nothing to prune
      ```
    </check>
    <check if="orphan_count > 0 OR stale_count > 0">
      <action>Display:</action>
      ```
      prune will remove:
      {{if orphan_count > 0}}
      ? orphan worktrees:
      {{for each orphan}}
        {{orphan.path}} → {{orphan.branch}}
      {{end for}}
      {{end if}}
      {{if stale_count > 0}}
      ! stale workspaces (>14 days):
      {{for each stale}}
        {{stale.branch}} ({{stale.days_old}} days)
      {{end for}}
      {{end if}}

      proceed? [y/n]
      ```
      <action>Wait for confirmation</action>
      <check if="confirmed">
        <action>Execute task: prune-nooks</action>
      </check>
    </check>
  </check>

  <check if="menu_choice == R or menu_choice == refresh">
    <action>Go back to Step 4 (rescan) - display tree again</action>
  </check>

  <check if="menu_choice == Q or menu_choice == quit">
    <action>Exit workflow</action>
  </check>
</step>

</workflow>

## Notes

<notes>
**Status Icons:**

| Icon | Status | Meaning |
|------|--------|---------|
| ● | active | has worktree + git branch exists |
| ○ | inactive | no worktree, but git branch exists |
| ✗ | orphan-tracking | context.yaml exists, but git branch DELETED |
| ! | stale | sync > 14 days old |
| ? | orphan-worktree | worktree exists, no tracking |
| ← | here | current branch |

**Critical Status Check:**
For each tracked workspace, verify the git branch actually exists:
```bash
git show-ref --verify --quiet "refs/heads/{{branch}}"
```
If branch doesn't exist, it's **orphan-tracking** (✗), not synced (○).

**Stale Detection:**
- Workspaces synced > 14 days ago show ! instead of ● or ○
- Time shown as "X days ago"

**How Stale Detection Works:**
1. Read `synced_at` from each context.yaml
2. Calculate days since sync
3. If > 14 days: mark as stale (!)

**Tracking Location:**
All tracking data lives in `base_workspace_path/.bmad-tracking/`

**Menu Actions:**
- `sleep`: soft-delete - removes worktree only, keeps branch + tracking (● → ○)
- `delete`: hard-delete - removes worktree + branch + tracking
- `prune`: removes orphan worktrees and stale tracking

**Reactivating Inactive Workspaces:**
```bash
git worktree add "{{worktrees_folder}}/{{name}}" "{{branch}}"
cd "{{worktrees_folder}}/{{name}}"
/bmad:th:workflows:nook-restore
```
</notes>
