---
name: 'step-02-detect-branch'
description: 'Detect current git branch and hash'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-02-detect-branch.md'
nextStepFile: '{workflow_path}/steps/step-03-load-config.md'
---

# Step 2: Detect Current Branch and Hash

## STEP GOAL:

Detect the current git branch, commit hash, and working directory for sync operations.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- HALT if not in a git repository

### Role Reinforcement:

- You are a workspace automation assistant
- Gather accurate git information
- Be precise with variable storage

## EXECUTION SEQUENCE:

### 1. Get Current Branch

```bash
git branch --show-current
```

<action>Store result as {{current_branch}}</action>

<check if="no branch detected or not in git repo">
  <action>Display error:</action>
  ```
  ERROR - NOT IN GIT REPOSITORY

  Cannot detect current branch.
  This workflow must be run from within a git worktree.
  ```
  <action>HALT workflow</action>
</check>

### 2. Get Current Commit Hash

```bash
git rev-parse --short=4 HEAD
```

<action>Store result as {{current_hash}}</action>

### 3. Get Current Working Directory

```bash
pwd
```

<action>Store result as {{current_path}}</action>

### 4. Display Branch Info

<action>Display:</action>
```
Current branch: {{current_branch}}
Commit hash:    {{current_hash}}
Working dir:    {{current_path}}
```

### 5. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After detection:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Current branch detected
- Commit hash captured (4 char short hash)
- Working directory stored
- All variables available for subsequent steps
- Proceeded to step 3

### SYSTEM FAILURE:
- Not detecting branch
- Not capturing hash
- Proceeding without git context
