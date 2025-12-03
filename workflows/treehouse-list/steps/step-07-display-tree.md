---
name: 'step-07-display-tree'
description: 'Display the workspace tree and menu options'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-07-display-tree.md'
nextStepFile: '{workflow_path}/steps/step-08-handle-menu.md'
---

# Step 7: Display Workspace Tree

## STEP GOAL:

Display the complete workspace tree with status indicators and present the menu for user interaction.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Display ONLY the formatted output below
- Do NOT show any processing steps or intermediate data
- HALT and wait for user menu selection

### Role Reinforcement:

- You are a workspace navigation assistant
- Present information cleanly and clearly
- Wait for user input after displaying

## EXECUTION SEQUENCE:

### 1. Display Tree Output

<critical>Do NOT show internal processing. Display ONLY this formatted output:</critical>

```
🌳 Workspace Lineage

base {{base_path}}
here {{current_branch}} · {{current_hash}}

{{tree_output}}

{{if orphan_worktrees exist}}
? orphan worktrees
{{for each orphan}}
  {{orphan.path}} → {{orphan.branch}}
{{end for}}
{{end if}}

─────────────────────────────────────────────────
 {{total}} tracked   ● {{active_count}} active   ○ {{inactive_count}} inactive
 ! {{stale_count}} stale      ✗ {{orphan_tracking_count}} orphan    ? {{orphan_wt_count}} orphan-wt
─────────────────────────────────────────────────
 [n] navigate   [s] sleep    [d] delete
 [p] prune      [r] refresh  [q] quit
```

### 2. Tree Output Format

**Status Icons:**
- `●` = active (has worktree + branch exists)
- `○` = inactive (no worktree, branch exists)
- `✗` = orphan-tracking (context exists, but git branch deleted)
- `!` = stale (>14 days since sync) - replaces ● or ○
- `?` = orphan-worktree (worktree exists, no tracking)
- `←` = you are here (current branch marker)

**Tree Structure:**

```
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
```

**Time Display:**
- "today" for 0 days
- "yesterday" for 1 day
- "X days ago" for 2-30 days
- "X weeks ago" for 31-60 days
- "X months ago" for 61+ days

**Summary Grid Format (CRITICAL - must be aligned):**

The summary section MUST use fixed-width columns for grid alignment:

```
─────────────────────────────────────────────────
 COL1 (width 14)  COL2 (width 13)  COL3 (width 14)
─────────────────────────────────────────────────
```

Row 1 format (pad numbers to 2 chars):
- Column 1: `{total:>2} tracked` padded to 14 chars
- Column 2: `● {active:>2} active` padded to 13 chars
- Column 3: `○ {inactive:>2} inactive` padded to 14 chars

Row 2 format:
- Column 1: `! {stale:>2} stale` padded to 14 chars
- Column 2: `✗ {orphan:>2} orphan` padded to 13 chars
- Column 3: `? {orphan_wt:>2} orphan-wt` padded to 14 chars

Menu row format (same columns):
- Column 1: `[n] navigate` padded to 14 chars
- Column 2: `[s] sleep` padded to 13 chars
- Column 3: `[d] delete` padded to 14 chars

Example with actual alignment:
```
─────────────────────────────────────────────────
 11 tracked     ● 11 active    ○  0 inactive
 !  0 stale     ✗  0 orphan    ?  2 orphan-wt
─────────────────────────────────────────────────
 [n] navigate   [s] sleep      [d] delete
 [p] prune      [r] refresh    [q] quit
```

### 3. Wait for Menu Selection

<action>HALT and wait for user input: `{{menu_choice}}`</action>

User can enter:
- `n` or `navigate` - Navigate to a workspace
- `s` or `sleep` - Sleep an active workspace
- `d` or `delete` - Delete a workspace
- `p` or `prune` - Prune orphans and stale
- `r` or `refresh` - Refresh the tree display
- `q` or `quit` - Exit the workflow

## MENU OPTIONS:

After user input, load and execute `{nextStepFile}` with `{{menu_choice}}` value.

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Tree displayed cleanly without processing noise
- All status indicators correct
- Summary counts displayed
- Menu presented
- Waited for user input

### SYSTEM FAILURE:
- Showing intermediate processing output
- Not displaying tree structure correctly
- Not waiting for user input
- Missing status indicators or counts
