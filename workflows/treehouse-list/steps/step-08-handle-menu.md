---
name: 'step-08-handle-menu'
description: 'Handle menu selection and execute corresponding action'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-08-handle-menu.md'
refreshStepFile: '{workflow_path}/steps/step-04-scan-workspaces.md'
---

# Step 8: Handle Menu Selection

## STEP GOAL:

Process the user's menu selection and execute the corresponding action.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Execute the action matching user's `{{menu_choice}}`
- Loop back to refresh step after actions that modify state

### Role Reinforcement:

- You are a workspace navigation assistant
- Execute menu actions efficiently
- Provide clear feedback

## EXECUTION SEQUENCE:

### Handle: Navigate (n)

<check if="menu_choice == N or menu_choice == n or menu_choice == navigate">
  <action>Display numbered list of branches:</action>

```
select branch:
{{for each workspace with index}}
{{index}}. {{status_icon}} {{workspace.branch}}
{{end for}}

enter number or q to cancel:
```

  <action>Wait for selection: `{{nav_selection}}`</action>

  <check if="nav_selection == q">
    <action>Return to tree display (load step-07)</action>
  </check>

  <check if="valid number selected">
    <check if="selected branch has active worktree">
      <action>Display:</action>
      ```
      -> cd {{worktree_path}}
      ```
    </check>

    <check if="selected branch is inactive (no worktree)">
      <action>Display:</action>
      ```
      o {{branch_name}} has no worktree

      to reactivate:
        git worktree add "{{suggested_path}}" "{{branch_name}}"
        cd "{{suggested_path}}"
        /bmad:th:workflows:nook-restore
      ```
    </check>

    <action>Return to tree display (load step-07)</action>
  </check>
</check>

---

### Handle: Sleep (s)

<check if="menu_choice == S or menu_choice == s or menu_choice == sleep">
  <action>Display numbered list of ACTIVE branches only:</action>

```
select branch to sleep (removes worktree, keeps tracking + branch):
{{for each ACTIVE workspace with index}}
{{index}}. . {{workspace.branch}}
{{end for}}

enter number or q to cancel:
```

  <action>Wait for selection: `{{sleep_selection}}`</action>

  <check if="sleep_selection == q">
    <action>Return to tree display (load step-07)</action>
  </check>

  <check if="valid number selected">
    <check if="selected branch is current branch">
      <action>Display:</action>
      ```
      cannot sleep current branch

      switch to another worktree first
      ```
      <action>Return to tree display (load step-07)</action>
    </check>

    <check if="selected branch is NOT current branch">
      <action>Display confirmation:</action>
      ```
      sleep {{selected_branch}}?

      this will:
        remove worktree folder
        keep git branch
        keep tracking data

      to reactivate later:
        git worktree add "{{worktree_path}}" "{{branch}}"
        /bmad:th:workflows:nook-restore

      proceed? [y/n]
      ```

      <action>Wait for confirmation: `{{sleep_confirm}}`</action>

      <check if="sleep_confirm == y or sleep_confirm == Y">
        <action>Remove worktree:</action>
        ```bash
        git worktree remove "{{worktree_path}}"
        ```

        <action>Display:</action>
        ```
        o {{selected_branch}} is now inactive
        ```

        <action>Rescan and redisplay (load `{refreshStepFile}`)</action>
      </check>

      <check if="sleep_confirm == n or sleep_confirm == N">
        <action>Return to tree display (load step-07)</action>
      </check>
    </check>
  </check>
</check>

---

### Handle: Delete (d)

<check if="menu_choice == D or menu_choice == d or menu_choice == delete">
  <action>Display numbered list of ALL branches:</action>

```
select branch to delete:
{{for each workspace with index}}
{{index}}. {{status_icon}} {{workspace.branch}}
{{end for}}

enter number or q to cancel:
```

  <action>Wait for selection: `{{delete_selection}}`</action>

  <check if="delete_selection == q">
    <action>Return to tree display (load step-07)</action>
  </check>

  <check if="valid number selected">
    <check if="selected branch is current branch">
      <action>Display:</action>
      ```
      cannot delete current branch

      switch to another worktree first
      ```
      <action>Return to tree display (load step-07)</action>
    </check>

    <check if="selected branch is NOT current branch">
      <action>Display confirmation:</action>
      ```
      DELETE {{selected_branch}}? [y/n]

      this will PERMANENTLY remove:
        - worktree folder (if exists)
        - git branch
        - tracking data
      ```

      <action>Wait for confirmation: `{{delete_confirm}}`</action>

      <check if="delete_confirm == y or delete_confirm == Y">
        <action>Remove worktree if exists:</action>
        ```bash
        git worktree remove "{{worktree_path}}" 2>/dev/null || true
        ```

        <action>Delete git branch:</action>
        ```bash
        git branch -D "{{selected_branch}}"
        ```

        <action>Remove tracking folder:</action>
        ```bash
        rm -rf "{{base_path}}/.bmad-tracking/{{tracking_path}}"
        ```

        <action>Commit tracking removal:</action>
        ```bash
        cd "{{base_path}}"
        git add .bmad-tracking/
        git commit -m "chore(th): delete nook {{selected_branch}}

Generated with Claude Code"
        ```

        <action>Display:</action>
        ```
        {{selected_branch}} deleted
        ```

        <action>Rescan and redisplay (load `{refreshStepFile}`)</action>
      </check>

      <check if="delete_confirm == n or delete_confirm == N">
        <action>Return to tree display (load step-07)</action>
      </check>
    </check>
  </check>
</check>

---

### Handle: Prune (p)

<check if="menu_choice == P or menu_choice == p or menu_choice == prune">
  <check if="orphan_wt_count == 0 AND stale_count == 0">
    <action>Display:</action>
    ```
    nothing to prune
    ```
    <action>Return to tree display (load step-07)</action>
  </check>

  <check if="orphan_wt_count > 0 OR stale_count > 0">
    <action>Display:</action>
    ```
    prune will remove:
    {{if orphan_wt_count > 0}}
    ? orphan worktrees:
    {{for each orphan_worktree}}
      {{orphan.path}} -> {{orphan.branch}}
    {{end for}}
    {{end if}}
    {{if stale_count > 0}}
    ! stale workspaces (>14 days):
    {{for each stale_workspace}}
      {{stale.branch}} ({{stale.days_old}} days)
    {{end for}}
    {{end if}}

    proceed? [y/n]
    ```

    <action>Wait for confirmation: `{{prune_confirm}}`</action>

    <check if="prune_confirm == y or prune_confirm == Y">
      <action>For each orphan worktree:</action>
      ```bash
      git worktree remove "{{orphan.path}}" --force
      ```

      <action>For each stale workspace:</action>
      ```bash
      # Remove worktree if exists
      git worktree remove "{{worktree_path}}" 2>/dev/null || true
      # Delete branch
      git branch -D "{{stale.branch}}" 2>/dev/null || true
      # Remove tracking
      rm -rf "{{base_path}}/.bmad-tracking/{{tracking_path}}"
      ```

      <action>Commit tracking changes:</action>
      ```bash
      cd "{{base_path}}"
      git add .bmad-tracking/
      git commit -m "chore(th): prune stale nooks and orphan worktrees

Generated with Claude Code"
      ```

      <action>Display:</action>
      ```
      pruned {{orphan_wt_count}} orphan worktrees and {{stale_count}} stale workspaces
      ```

      <action>Rescan and redisplay (load `{refreshStepFile}`)</action>
    </check>

    <check if="prune_confirm == n or prune_confirm == N">
      <action>Return to tree display (load step-07)</action>
    </check>
  </check>
</check>

---

### Handle: Refresh (r)

<check if="menu_choice == R or menu_choice == r or menu_choice == refresh">
  <action>Rescan and redisplay (load `{refreshStepFile}`)</action>
</check>

---

### Handle: Quit (q)

<check if="menu_choice == Q or menu_choice == q or menu_choice == quit">
  <action>Exit workflow</action>
</check>

---

### Handle: Unknown Input

<check if="menu_choice does not match any option">
  <action>Display:</action>
  ```
  unknown option: {{menu_choice}}
  ```
  <action>Return to tree display (load step-07)</action>
</check>

## MENU OPTIONS:

After action completion:
- Most actions loop back to `{refreshStepFile}` to rescan and redisplay
- Quit exits the workflow
- Invalid input returns to step-07 display

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Correct action executed for menu choice
- Proper confirmations before destructive actions
- State-modifying actions trigger rescan
- Clean exit on quit

### SYSTEM FAILURE:
- Executing wrong action for menu choice
- Not confirming before delete/prune
- Not rescanning after state changes
- Not handling unknown input gracefully
