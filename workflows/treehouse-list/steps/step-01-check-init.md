---
name: 'step-01-check-init'
description: 'Check if workspace is initialized and get base path'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-01-check-init.md'
nextStepFile: '{workflow_path}/steps/step-02-branch-info.md'
---

# Step 1: Check Workspace Initialized

## STEP GOAL:

Check if the base workspace has been initialized with `treehouse-init` and retrieve the base path configuration.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Execute bash commands exactly as specified
- HALT immediately if workspace is not initialized

### Role Reinforcement:

- You are a workspace navigation assistant
- This is a utility workflow - check prerequisites
- Be direct and efficient

## EXECUTION SEQUENCE:

### 1. Read Configuration

```bash
grep "^base_workspace_path:" .bmad/th/config.yaml
```

### 2. Check Initialization Status

<check if="base_workspace_path is empty or not set">
  <action>Display initialization prompt:</action>

```
🌳 Workspace Lineage

⚠ not initialized

run /bmad:th:workflows:treehouse-init first

[i] initialize  [q] quit
```

  <action>Wait for user input: {{init_choice}}</action>

  <check if="init_choice == I or init_choice == i or init_choice == initialize">
    <action>Display: "Launching treehouse-init..."</action>
    <action>Execute workflow: /bmad:th:workflows:treehouse-init</action>
    <action>Exit this workflow</action>
  </check>

  <check if="init_choice == Q or init_choice == q or init_choice == quit">
    <action>Exit workflow</action>
  </check>
</check>

### 3. Store Variables and Continue

<check if="base_workspace_path is configured">
  <action>Store variables:</action>
  - `{{base_path}}` = base_workspace_path value
  - `{{base_branch}}` = base_branch value from config

  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

- If initialized: Auto-proceed to `{nextStepFile}`
- If not initialized: Wait for [i/q] selection

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Configuration read successfully
- base_workspace_path is set
- Variables stored for subsequent steps
- Proceeded to step 2

### SYSTEM FAILURE:
- Proceeding without checking config
- Not offering init option when not initialized
- Skipping to later steps
