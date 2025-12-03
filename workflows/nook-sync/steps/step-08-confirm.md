---
name: 'step-08-confirm'
description: 'Confirm sync operation'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-08-confirm.md'
nextStepFile: '{workflow_path}/steps/step-09-execute.md'
---

# Step 8: Confirm Sync Operation

## STEP GOAL:

Get explicit user confirmation before executing the sync operation.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- WAIT for user input - do not auto-proceed
- Only proceed with explicit confirmation

### Role Reinforcement:

- You are a workspace automation assistant
- This is a user-input step - must wait for response
- Respect user's decision

## EXECUTION SEQUENCE:

### 1. Present Confirmation Prompt

<action>Display confirmation prompt:</action>

```
------------------------------------------------------------
                      CONFIRM SYNC
------------------------------------------------------------

You are about to sync context:

  From: {{current_path}}
  To:   {{tracking_path}}

  Existing content in tracking will be REPLACED.

------------------------------------------------------------

Options:
  [Y] Yes, proceed with sync
  [N] No, cancel sync
  [D] Show details again

Choose:
```

### 2. Handle User Response

<check if="user chooses Y (yes)">
  <action>Display: "Proceeding with sync..."</action>
  <action>Load and execute `{nextStepFile}`</action>
</check>

<check if="user chooses N (no)">
  <action>Display:</action>
  ```
  Sync cancelled by user.

  No changes were made.

  To sync later, run:
    /bmad:th:workflows:nook-sync
  ```
  <action>HALT workflow</action>
</check>

<check if="user chooses D (details)">
  <action>Re-display the preview from step 7</action>
  <action>Return to confirmation prompt</action>
</check>

## MENU OPTIONS:

This is a user-input step. WAIT for user response:
- [Y] Proceed to `{nextStepFile}`
- [N] Cancel and HALT
- [D] Show details and re-prompt

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Confirmation prompt displayed clearly
- Waited for user input
- Handled response appropriately
- Proceeded only with explicit Y confirmation

### SYSTEM FAILURE:
- Auto-proceeding without confirmation
- Not offering cancel option
- Ignoring user's cancel choice
