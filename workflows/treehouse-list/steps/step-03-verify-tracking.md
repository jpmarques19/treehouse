---
name: 'step-03-verify-tracking'
description: 'Verify base workspace tracking folder exists'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-03-verify-tracking.md'
nextStepFile: '{workflow_path}/steps/step-04-scan-workspaces.md'
---

# Step 3: Verify Tracking Folder

## STEP GOAL:

Verify that the tracking folder exists at the base workspace path.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Execute bash commands exactly as specified
- HALT if tracking folder is missing

### Role Reinforcement:

- You are a workspace navigation assistant
- Report clear error if tracking folder missing
- Auto-proceed if folder exists

## EXECUTION SEQUENCE:

### 1. Check Tracking Folder Exists

```bash
test -d "{{base_path}}/.bmad-tracking" && echo "EXISTS" || echo "MISSING"
```

### 2. Handle Result

<check if="tracking folder does NOT exist (MISSING)">
  <action>Display error:</action>

```
Workspace Lineage

tracking folder not found

base     {{base_path}}
expected {{base_path}}/.bmad-tracking/

re-run /bmad:th:workflows:treehouse-init from base
```

  <action>HALT workflow</action>
</check>

<check if="tracking folder EXISTS">
  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

- If exists: Auto-proceed to `{nextStepFile}`
- If missing: HALT with error message

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Tracking folder verified to exist
- Proceeded to step 4

### SYSTEM FAILURE:
- Proceeding without verifying folder exists
- Not providing clear instructions when missing
