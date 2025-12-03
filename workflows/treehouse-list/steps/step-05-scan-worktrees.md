---
name: 'step-05-scan-worktrees'
description: 'Scan git worktrees to determine active/inactive status'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-05-scan-worktrees.md'
nextStepFile: '{workflow_path}/steps/step-06-build-tree.md'
---

# Step 5: Scan Git Worktrees

## STEP GOAL:

Get list of all git worktrees to determine which tracked workspaces are active (have a worktree) vs inactive (no worktree).

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Store results for tree building

### Role Reinforcement:

- You are a workspace navigation assistant
- Scan silently, no user interaction
- Auto-proceed when done

## EXECUTION SEQUENCE:

### 1. List All Worktrees

```bash
git worktree list --porcelain
```

### 2. Parse Worktree Output

The porcelain output format is:

```
worktree /path/to/worktree
HEAD abc123def456
branch refs/heads/branch-name

worktree /path/to/another
HEAD def456abc123
branch refs/heads/other-branch
```

<action>Parse output to extract for each worktree:</action>
- `path` - the worktree filesystem path
- `head` - the HEAD commit hash
- `branch` - the branch name (strip `refs/heads/` prefix)

Store as `{{worktrees}}` array.

### 3. Continue

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. No user interaction required.

Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All worktrees discovered
- Path, head, and branch extracted for each
- Worktrees array populated
- Proceeded to step 6

### SYSTEM FAILURE:
- Not parsing porcelain format correctly
- Missing worktree entries
- Displaying intermediate output to user
