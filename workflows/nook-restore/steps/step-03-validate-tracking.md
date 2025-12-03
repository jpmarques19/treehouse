---
name: 'step-03-validate-tracking'
description: 'Validate that tracking path exists for selected branch'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'

# File References
thisStepFile: '{workflow_path}/steps/step-03-validate-tracking.md'
nextStepFile: '{workflow_path}/steps/step-04-load-config.md'
---

# Step 3: Validate Tracking Path

## STEP GOAL:

Validate that the tracking folder exists for the selected restore branch and load context metadata.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- HALT if tracking path does not exist

### Role Reinforcement:

- You are a workspace automation assistant
- Validate before proceeding
- Provide clear error messages with actionable options

## EXECUTION SEQUENCE:

### 1. Calculate Tracking Path

<action>Calculate tracking path:</action>
```
{{tracking_path}} = {{base_path}}/.bmad-tracking/{{restore_branch}}
```

<action>Display: "Tracking path: {{tracking_path}}"</action>

### 2. Check Tracking Path Exists

```bash
test -d "{{tracking_path}}" && echo "EXISTS" || echo "NOT_FOUND"
```

<check if="tracking path does not exist">
  <action>List available contexts:</action>
  ```bash
  ls -1 "{{base_path}}/.bmad-tracking/" 2>/dev/null || echo "(none)"
  ```

  <action>Display error:</action>
  ```
  NO CONTEXT FOUND

  No saved context for branch: {{restore_branch}}

  Expected location: {{tracking_path}}

  Available contexts:
  {{list from above command}}

  Options:
  - Run nook-sync from a worktree with content
  - Check if branch name is correct
  - Start fresh without restoring
  ```
  <action>HALT workflow</action>
</check>

### 3. Load Context Metadata

<check if="tracking path exists">
  <action>Check for context.yaml:</action>
  ```bash
  test -f "{{tracking_path}}/context.yaml" && echo "EXISTS" || echo "NOT_FOUND"
  ```

  <check if="context.yaml exists">
    <action>Load and display context metadata:</action>
    ```bash
    cat "{{tracking_path}}/context.yaml"
    ```

    <action>Extract and store key fields:</action>
    - {{synced_at}} = synced_at value
    - {{synced_by}} = synced_by value (if present)
    - {{forked_from}} = forked_from value (if present)

    <action>Display context summary:</action>
    ```
    Context Metadata:
    - Synced at: {{synced_at}}
    - Synced by: {{synced_by}}
    - Forked from: {{forked_from}}
    ```
  </check>

  <check if="context.yaml does not exist">
    <action>Display warning:</action>
    ```
    NOTE: No context.yaml found in tracking folder.
    Files will be restored but metadata is unavailable.
    ```
  </check>

  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

This is an auto-proceed step. After validation:
- If tracking path exists: Load and execute `{nextStepFile}`
- If tracking path missing: HALT with error

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Tracking path calculated correctly
- Path existence verified
- Context metadata loaded (if available)
- Key variables stored
- Proceeded to step 4

### SYSTEM FAILURE:
- Not verifying tracking path exists
- Proceeding when path is missing
- Not loading available context metadata
