---
name: 'step-07-update-config'
description: 'Update th config with base workspace info'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-07-update-config.md'
nextStepFile: '{workflow_path}/steps/step-08-commit-changes.md'
---

# Step 7: Update Config

## STEP GOAL:

Update the th module config with base workspace path and branch.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Update both base_workspace_path and base_branch
- Preserve rest of config file

### Role Reinforcement:

- You are a workspace automation assistant
- Config is source of truth for workspace location
- Be precise with file editing

## EXECUTION SEQUENCE:

### 1. Update base_workspace_path

<action>Edit .bmad/th/config.yaml:</action>
- Replace `base_workspace_path: ""` with `base_workspace_path: "{{workspace_path}}"`

### 2. Update base_branch

<action>Edit .bmad/th/config.yaml:</action>
- Replace `base_branch: ""` with `base_branch: "{{current_branch}}"`

### 3. Verify Updates

```bash
grep -E "^base_workspace_path:|^base_branch:" .bmad/th/config.yaml
```

<action>Display:</action>
```
Config updated:
  base_workspace_path: {{workspace_path}}
  base_branch: {{current_branch}}
```

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After update:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- base_workspace_path set correctly
- base_branch set correctly
- Config file valid after edit
- Proceeded to step 8

### SYSTEM FAILURE:
- Corrupting config file
- Missing either field
- Incorrect values
