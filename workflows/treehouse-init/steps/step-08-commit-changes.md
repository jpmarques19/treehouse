---
name: 'step-08-commit-changes'
description: 'Commit .gitignore, config, and cache removal'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'

# File References
thisStepFile: '{workflow_path}/steps/step-08-commit-changes.md'
nextStepFile: '{workflow_path}/steps/step-09-create-submodule.md'
---

# Step 8: Commit Changes

## STEP GOAL:

Commit the .gitignore changes, config updates, and cache removals.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Stage all relevant changes
- Create descriptive commit message

### Role Reinforcement:

- You are a workspace automation assistant
- Clean commit history is important
- Include all context in commit message

## EXECUTION SEQUENCE:

### 1. Stage Changes

```bash
git add -A .gitignore .bmad/th/config.yaml
```

<action>Note: `git add -A` will also stage the deletions from `git rm --cached`</action>

### 2. Create Commit

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

### 3. Verify Commit

```bash
git log -1 --oneline
```

<action>Display: "Changes committed: {{commit_hash}}"</action>

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After commit:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All changes staged
- Commit created with descriptive message
- Commit verified
- Proceeded to step 9

### SYSTEM FAILURE:
- Commit failed
- Missing files in commit
- No commit message
