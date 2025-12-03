---
name: 'step-03-get-branch-path'
description: 'Get current branch and workspace path'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'

# File References
thisStepFile: '{workflow_path}/steps/step-03-get-branch-path.md'
nextStepFile: '{workflow_path}/steps/step-04-load-config.md'
---

# Step 3: Get Branch and Path

## STEP GOAL:

Capture the current git branch and absolute workspace path.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Store variables for subsequent steps

### Role Reinforcement:

- You are a workspace automation assistant
- Capture accurate workspace information
- Be precise with path handling

## EXECUTION SEQUENCE:

### 1. Get Current Branch

```bash
git branch --show-current
```

<action>Store result as {{current_branch}}</action>

### 2. Get Workspace Path

```bash
pwd
```

<action>Store result as {{workspace_path}}</action>

### 3. Display Setup Info

<action>Display:</action>
```
Base Workspace Setup
Path:   {{workspace_path}}
Branch: {{current_branch}}
```

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After capturing:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Current branch captured
- Absolute workspace path captured
- Variables stored for subsequent steps
- Proceeded to step 4

### SYSTEM FAILURE:
- Not capturing branch
- Using relative path
- Missing variable storage
