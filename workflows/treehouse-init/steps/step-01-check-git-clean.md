---
name: 'step-01-check-git-clean'
description: 'Check for clean git state before initialization'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'

# File References
thisStepFile: '{workflow_path}/steps/step-01-check-git-clean.md'
nextStepFile: '{workflow_path}/steps/step-02-check-initialized.md'
---

# Step 1: Check Git State (CRITICAL)

## STEP GOAL:

Verify the git working tree is clean before proceeding with initialization.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- HALT immediately if uncommitted changes exist

### Role Reinforcement:

- You are a workspace automation assistant
- Clean git state is CRITICAL for initialization
- Do not proceed with uncommitted changes

## EXECUTION SEQUENCE:

### 1. Check Git Status

```bash
git status --porcelain
```

### 2. Handle Uncommitted Changes

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

### 3. Continue If Clean

<check if="git status is clean (no output)">
  <action>Display: "Git working tree is clean"</action>
  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

This is an auto-proceed step. After check:
- If clean: Load and execute `{nextStepFile}`
- If dirty: HALT with instructions

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Git status checked
- Working tree confirmed clean
- Proceeded to step 2

### SYSTEM FAILURE:
- Proceeding with uncommitted changes
- Not checking git status
- Ignoring dirty working tree
