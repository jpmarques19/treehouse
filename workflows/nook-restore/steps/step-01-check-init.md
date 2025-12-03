---
name: 'step-01-check-init'
description: 'Validate that the workspace has been initialized with treehouse-init'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-01-check-init.md'
nextStepFile: '{workflow_path}/steps/step-02-branch-detect.md'
---

# Step 1: Check Workspace Initialized

## STEP GOAL:

Validate that the base workspace has been initialized with `treehouse-init` before proceeding with context restore.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- HALT immediately if workspace is not initialized

### Role Reinforcement:

- You are a workspace automation assistant
- This is an action workflow - execute commands, check results
- Be direct and efficient

## EXECUTION SEQUENCE:

### 1. Read Configuration

Read `{module_config}` and extract:
- `base_workspace_path`

```bash
cat .bmad/th/config.yaml
```

### 2. Check Initialization Status

<check if="base_workspace_path is empty or not set">
  <action>Display critical halt message:</action>

```
CANNOT RESTORE - WORKSPACE NOT INITIALIZED

No base workspace has been configured.
The base workspace holds the tracking folder (.bmad-tracking/)
which is the single source of truth for all context.

Please run first (from the base worktree):
  /bmad:th:workflows:treehouse-init

This will establish the tracking folder location.
```

  <action>HALT workflow - do not proceed</action>
</check>

### 3. Verify Base Path Exists

<check if="base_workspace_path is configured">
  <action>Store variable:</action>
  - `{{base_path}}` = base_workspace_path value

  <action>Verify base_path exists and is accessible:</action>
  ```bash
  test -d "{{base_path}}" && echo "EXISTS" || echo "NOT_FOUND"
  ```

  <check if="base_path does not exist">
    <action>Display error:</action>
    ```
    ERROR - BASE WORKSPACE NOT FOUND

    Configured base_workspace_path does not exist:
      {{base_path}}

    The base workspace may have been moved or deleted.

    Options:
    - Ensure the base worktree exists at that path
    - Re-run /bmad:th:workflows:treehouse-init from the base
    - Update base_workspace_path in .bmad/th/config.yaml
    ```
    <action>HALT workflow</action>
  </check>

  <action>Display confirmation:</action>
  ```
  Base workspace: {{base_path}}
  ```

  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

This is an auto-proceed step. After validation:
- If initialized and path exists: Load and execute `{nextStepFile}`
- If not initialized or path missing: HALT with instructions

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Configuration read successfully
- base_workspace_path is set
- Base path exists and is accessible
- Variables stored for subsequent steps
- Proceeded to step 2

### SYSTEM FAILURE:
- Proceeding without checking config
- Not halting when workspace is not initialized
- Not verifying base path exists
- Skipping to later steps
