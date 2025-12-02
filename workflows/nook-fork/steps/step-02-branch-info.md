---
name: 'step-02-branch-info'
description: 'Get current branch info and validate tracking folder status'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-02-branch-info.md'
nextStepFile: '{workflow_path}/steps/step-03-context-sync.md'
---

# Step 2: Get Branch Info and Validate Tracking

## STEP GOAL:

Gather current branch information for lineage tracking and validate that the tracking folder has no uncommitted changes.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 Execute bash commands exactly as specified
- ⏸️ HALT and wait for user input if tracking folder is dirty

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Warn user about uncommitted tracking changes
- ✅ Offer options to handle dirty state

## EXECUTION SEQUENCE:

### 1. Get Current Branch

```bash
git branch --show-current
```

Store result as `{{current_branch}}`

### 2. Get Parent Hash (4 chars)

```bash
git rev-parse --short=4 HEAD
```

Store result as `{{parent_hash}}`

### 2b. Get Current Worktree Path

```bash
pwd
```

Store result as `{{parent_worktree_path}}` (the path of the current/parent worktree where YOLO agent will run)

### 3. Check Tracking Folder Status

```bash
git -C "{{base_path}}" status --porcelain .bmad-tracking/
```

### 4. Handle Tracking Status

<check if="tracking folder has modifications (non-empty output)">
  <action>Display warning with file list:</action>

```
WARNING - TRACKING FOLDER HAS MODIFICATIONS

The tracking folder at {{base_path}}/.bmad-tracking/
has uncommitted changes. This may indicate:
- A sync operation was interrupted
- Manual edits to context files
- Work in progress that should be committed

Uncommitted changes:
{{list first 10 files from git status}}

RECOMMENDATION: Commit these changes before creating nook
to ensure a clean baseline for the new nook.

Options:
[C] Commit changes now - Stage and commit tracking folder
[P] Proceed anyway    - Create nook without committing (not recommended)
[A] Abort             - Exit and handle manually

Choose [C/P/A]:
```

  <action>HALT and wait for user input</action>

  <check if="user_choice == C or commit">
    <action>Stage and commit tracking changes:</action>
    ```bash
    cd "{{base_path}}"
    git add .bmad-tracking/
    git commit -m "chore(th): commit tracking changes before nook-fork

Generated with Claude Code"
    ```
    <action>Display: "Tracking changes committed. Continuing..."</action>
    <action>Proceed to next step</action>
  </check>

  <check if="user_choice == P or proceed">
    <action>Display: "Proceeding with uncommitted tracking changes..."</action>
    <action>Proceed to next step</action>
  </check>

  <check if="user_choice == A or abort">
    <action>Display: "Aborted. Please handle tracking changes and try again."</action>
    <action>HALT workflow</action>
  </check>
</check>

<check if="tracking folder is clean (empty output)">
  <action>Display: "Tracking folder is clean"</action>
  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

- If tracking clean: Auto-proceed to `{nextStepFile}`
- If tracking dirty: Wait for [C/P/A] selection

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Current branch identified
- Parent hash captured (4 chars)
- Tracking status checked
- User informed of dirty state (if applicable)
- Appropriate action taken based on user choice

### ❌ SYSTEM FAILURE:
- Not checking tracking folder status
- Proceeding without user consent when dirty
- Not storing branch info for lineage
