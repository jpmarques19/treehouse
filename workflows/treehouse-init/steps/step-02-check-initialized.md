---
name: 'step-02-check-initialized'
description: 'Check if workspace is already initialized'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-02-check-initialized.md'
nextStepFile: '{workflow_path}/steps/step-03-get-branch-path.md'
---

# Step 2: Check If Already Initialized

## STEP GOAL:

Check if this workspace was already initialized and handle accordingly.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Check config for existing initialization
- WAIT for user input if already initialized

### Role Reinforcement:

- You are a workspace automation assistant
- Offer reinitialize option if already set up
- Respect user's choice

## EXECUTION SEQUENCE:

### 1. Check Config for Existing Initialization

```bash
grep "^base_workspace_path:" .bmad/th/config.yaml | grep -v '""'
```

<action>Also check base_branch:</action>
```bash
grep "^base_branch:" .bmad/th/config.yaml | grep -v '""'
```

### 2. Handle Already Initialized

<check if="base_workspace_path is not empty">
  <action>Extract current values:</action>
  - {{existing_base_path}} = base_workspace_path value
  - {{existing_base_branch}} = base_branch value

  <action>Display:</action>
  ```
  WORKSPACE ALREADY INITIALIZED

  This workspace tree was already initialized.
  Base workspace: {{existing_base_path}}
  Base branch:    {{existing_base_branch}}

  Options:
  [R] Reinitialize - Reset as new base workspace (use if base was deleted)
  [C] Cancel       - Exit without changes

  Choose [R/C]:
  ```

  <action>WAIT for user input</action>

  <check if="user chooses C (cancel)">
    <action>Display: "Cancelled. No changes made."</action>
    <action>HALT workflow</action>
  </check>

  <check if="user chooses R (reinitialize)">
    <action>Display: "Proceeding with reinitialization..."</action>
    <action>Auto-proceed to next step</action>
  </check>
</check>

### 3. Handle New Initialization

<check if="base_workspace_path is empty or not set">
  <action>Display: "Initializing new base workspace..."</action>
  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

- If not initialized: Auto-proceed to `{nextStepFile}`
- If initialized: Present [R/C] menu and wait for input

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Config checked for existing initialization
- User given choice if already initialized
- Proceeded only with user consent or new init

### SYSTEM FAILURE:
- Not checking for existing initialization
- Proceeding without user consent on reinit
- Not offering cancel option
