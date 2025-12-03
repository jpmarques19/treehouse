---
name: 'step-02-branch-info'
description: 'Get current branch, hash, and path information'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-02-branch-info.md'
nextStepFile: '{workflow_path}/steps/step-03-verify-tracking.md'
---

# Step 2: Get Branch Info

## STEP GOAL:

Gather current branch information for display in the tree header.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Store all results for later display

### Role Reinforcement:

- You are a workspace navigation assistant
- Gather info silently, no user interaction needed
- Auto-proceed when done

## EXECUTION SEQUENCE:

### 1. Get Current Branch

```bash
git branch --show-current
```

Store result as `{{current_branch}}`

### 2. Get Current Hash (4 chars)

```bash
git rev-parse --short=4 HEAD
```

Store result as `{{current_hash}}`

### 3. Get Current Directory

```bash
pwd
```

Store result as `{{current_path}}`

### 4. Continue

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. No user interaction required.

Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Current branch identified
- Current hash captured (4 chars)
- Current path stored
- Proceeded to step 3

### SYSTEM FAILURE:
- Not capturing all three values
- Halting for user input unnecessarily
