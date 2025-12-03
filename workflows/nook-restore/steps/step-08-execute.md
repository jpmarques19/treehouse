---
name: 'step-08-execute'
description: 'Execute the restore operations'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'

# File References
thisStepFile: '{workflow_path}/steps/step-08-execute.md'
nextStepFile: '{workflow_path}/steps/step-09-report.md'
---

# Step 8: Execute Restore Operations

## STEP GOAL:

Execute the actual file copy/sync operations to restore context from tracking to local.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Track success/failure for each path

### Role Reinforcement:

- You are a workspace automation assistant
- Execute operations carefully and track results
- Report any errors immediately

## EXECUTION SEQUENCE:

### 1. Initialize Results Tracking

<action>Initialize {{restore_results}} array to track each operation</action>

### 2. Process Each Path in Restore Plan

<action>For each path in {{restore_plan}} where action is RESTORE or OVERWRITE:</action>

#### 2a. Determine Path Type

```bash
tracking_loc="{{tracking_path}}/{{path}}"
if [ -d "$tracking_loc" ]; then
  echo "TYPE:directory"
elif [ -f "$tracking_loc" ]; then
  echo "TYPE:file"
fi
```

#### 2b. Restore Directory

<check if="path is a directory">
  <action>Create target directory if needed and sync:</action>
  ```bash
  mkdir -p "{{current_path}}/{{path}}"
  rsync -av --delete "{{tracking_path}}/{{path}}/" "{{current_path}}/{{path}}/"
  ```

  <action>Count restored files:</action>
  ```bash
  find "{{current_path}}/{{path}}" -type f | wc -l
  ```

  <action>Record result:</action>
  - Path: {{path}}
  - Type: directory
  - Status: SUCCESS or FAILED
  - File count: {{count}}
</check>

#### 2c. Restore File

<check if="path is a file">
  <action>Create parent directory if needed and copy:</action>
  ```bash
  mkdir -p "$(dirname "{{current_path}}/{{path}}")"
  cp "{{tracking_path}}/{{path}}" "{{current_path}}/{{path}}"
  ```

  <action>Verify copy:</action>
  ```bash
  test -f "{{current_path}}/{{path}}" && echo "SUCCESS" || echo "FAILED"
  ```

  <action>Record result:</action>
  - Path: {{path}}
  - Type: file
  - Status: SUCCESS or FAILED
</check>

### 3. Handle Errors

<check if="any operation failed">
  <action>Display error summary:</action>
  ```
  ⚠️  SOME OPERATIONS FAILED

  Failed paths:
  {{for each failed result}}
    - {{path}}: {{error_message}}
  {{end for}}

  Successfully restored:
  {{for each successful result}}
    - {{path}}
  {{end for}}
  ```

  <action>Continue to report step despite failures</action>
</check>

### 4. Store Results and Proceed

<action>Store {{restore_results}} for report step</action>
<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After execution:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All restore operations executed
- Results tracked for each path
- Errors reported immediately
- Proceeded to step 9

### SYSTEM FAILURE:
- Not executing all planned restores
- Not tracking operation results
- Silent failures without reporting
- Stopping on first error instead of continuing
