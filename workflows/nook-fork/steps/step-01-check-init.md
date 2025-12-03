---
name: 'step-01-check-init'
description: 'Validate that the workspace has been initialized with treehouse-init'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-01-check-init.md'
nextStepFile: '{workflow_path}/steps/step-02-branch-info.md'
---

# Step 1: Check Workspace Initialized

## STEP GOAL:

Validate that the base workspace has been initialized with `treehouse-init` before proceeding with nook creation.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 Execute bash commands exactly as specified
- 🛑 HALT immediately if workspace is not initialized

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ This is an action workflow - execute commands, check results
- ✅ Be direct and efficient

## EXECUTION SEQUENCE:

### 1. Read Configuration

Read `{module_config}` and extract:
- `base_workspace_path`
- `base_branch`

```bash
cat .bmad/th/config.yaml
```

### 2. Check Initialization Status

<check if="base_workspace_path is empty or not set">
  <action>Display critical halt message:</action>

```
CANNOT CREATE NOOK - WORKSPACE NOT INITIALIZED

This workspace hasn't been initialized as a base workspace.
The base workspace holds the tracking folder (.bmad-tracking/)
which is the source of truth for all nooks.

Please run first:
  /th:workflows:treehouse-init

This will:
1. Add sync artifacts to .gitignore
2. Create the tracking folder
3. Establish this worktree as the base workspace
```

  <action>HALT workflow - do not proceed</action>
</check>

### 3. Store Variables and Continue

<check if="base_workspace_path is configured">
  <action>Store variables:</action>
  - `{{base_path}}` = base_workspace_path value
  - `{{base_branch}}` = base_branch value

  <action>Display confirmation:</action>
  ```
  Base workspace: {{base_path}} ({{base_branch}})
  ```

  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

This is an auto-proceed step. After validation:
- If initialized: Load and execute `{nextStepFile}`
- If not initialized: HALT with instructions

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Configuration read successfully
- base_workspace_path is set
- Variables stored for subsequent steps
- Proceeded to step 2

### ❌ SYSTEM FAILURE:
- Proceeding without checking config
- Not halting when workspace is not initialized
- Skipping to later steps
