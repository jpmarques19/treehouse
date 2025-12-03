---
name: 'step-06-build-tree'
description: 'Build workspace tree structure with status and stale detection'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-06-build-tree.md'
nextStepFile: '{workflow_path}/steps/step-07-display-tree.md'
---

# Step 6: Build Workspace Tree

## STEP GOAL:

Build the hierarchical tree structure from workspace data, determine status for each workspace, and detect stale/orphan conditions.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Process all data silently
- Do NOT display anything to user in this step

### Role Reinforcement:

- You are a workspace navigation assistant
- Build tree data structure silently
- Auto-proceed when done

## EXECUTION SEQUENCE:

### 1. Build Tree Hierarchy

<action>For each workspace in `{{workspaces}}`:</action>

Determine parent-child relationships:
- `parent` = the `forked_from` field value
- `children` = workspaces where `forked_from` matches this branch
- `depth` = levels from root (workspaces with no parent or parent not in tracked list are depth 0)

### 2. Check Git Branch Existence

<action>For each tracked workspace, verify git branch exists:</action>

```bash
git show-ref --verify --quiet "refs/heads/{{branch}}" && echo "EXISTS" || echo "MISSING"
```

Store result as `branch_exists` for each workspace.

### 3. Calculate Staleness

<action>For each workspace, calculate days since last sync:</action>

```bash
# Extract synced_at from context.yaml
synced_at=$(grep "synced_at:" "{{base_path}}/.bmad-tracking/{{tracking_path}}/context.yaml" | cut -d: -f2- | xargs)

# Calculate days since sync
if [ -n "$synced_at" ]; then
  synced_epoch=$(date -d "$synced_at" +%s 2>/dev/null || echo 0)
  now_epoch=$(date +%s)
  days_old=$(( (now_epoch - synced_epoch) / 86400 ))
else
  days_old=0
fi
echo $days_old
```

Store `days_old` for each workspace.

### 4. Determine Status for Each Workspace

<action>Apply status logic:</action>

| Condition | Status |
|-----------|--------|
| context.yaml exists + branch exists + worktree exists | **active** |
| context.yaml exists + branch exists + NO worktree | **inactive** |
| context.yaml exists + branch does NOT exist | **orphan-tracking** |
| days_old > 14 | add **stale** flag |
| branch == current_branch | add **current** flag |

### 5. Identify Orphan Worktrees

<action>Find worktrees with no matching context.yaml:</action>

For each worktree in `{{worktrees}}`:
- Check if any workspace in `{{workspaces}}` has matching branch
- If no match: add to `{{orphan_worktrees}}` list

### 6. Calculate Summary Counts

<action>Count workspaces by status:</action>
- `{{total}}` = total tracked workspaces
- `{{active_count}}` = workspaces with status "active"
- `{{inactive_count}}` = workspaces with status "inactive"
- `{{stale_count}}` = workspaces with stale flag
- `{{orphan_tracking_count}}` = workspaces with status "orphan-tracking"
- `{{orphan_wt_count}}` = count of orphan worktrees

### 7. Store Tree Structure

Store complete tree as `{{workspace_tree}}` containing:
- All workspaces with their status, depth, parent, children
- Orphan worktrees list
- Summary counts

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. No user interaction required.

Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Tree hierarchy built correctly
- All statuses determined
- Stale detection working
- Orphan detection working
- Summary counts calculated
- Proceeded to step 7

### SYSTEM FAILURE:
- Displaying intermediate output to user
- Not checking branch existence
- Not calculating staleness
- Missing orphan detection
