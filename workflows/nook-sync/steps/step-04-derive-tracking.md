---
name: 'step-04-derive-tracking'
description: 'Derive tracking path in base workspace'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-04-derive-tracking.md'
nextStepFile: '{workflow_path}/steps/step-05-load-context.md'
---

# Step 4: Derive Tracking Path

## STEP GOAL:

Calculate the tracking path in the base workspace where context will be saved.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Branch hierarchy is preserved in tracking path
- Store tracking path for subsequent steps

### Role Reinforcement:

- You are a workspace automation assistant
- Path derivation is straightforward
- Be precise with path construction

## EXECUTION SEQUENCE:

### 1. Calculate Tracking Path

<action>Calculate tracking path based on branch name:</action>

```
{{tracking_path}} = {{base_path}}/.bmad-tracking/{{current_branch}}
```

**Note:** The branch hierarchy is preserved. For example:
- Branch `discovery/peek-a-box-mvp` → `.bmad-tracking/discovery/peek-a-box-mvp/`
- Branch `main` → `.bmad-tracking/main/`
- Branch `th/feature/agent-wizard` → `.bmad-tracking/th/feature/agent-wizard/`

### 2. Display Tracking Path

<action>Display:</action>
```
Tracking path: {{tracking_path}}
```

### 3. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After derivation:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Tracking path calculated correctly
- Branch hierarchy preserved
- Path stored for subsequent steps
- Proceeded to step 5

### SYSTEM FAILURE:
- Incorrect path construction
- Not preserving branch hierarchy
- Missing path variable
