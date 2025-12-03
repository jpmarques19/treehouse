---
name: 'step-09-execute'
description: 'Execute sync operations'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-09-execute.md'
nextStepFile: '{workflow_path}/steps/step-10-generate-summary.md'
---

# Step 9: Execute Sync Operations

## STEP GOAL:

Execute the actual file copy/sync operations from current directory to tracking folder.

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

<action>Initialize {{sync_results}} array to track each operation</action>

### 2. Process Each Path

<action>For each path in {{sync_paths}} that exists in source:</action>

#### 2a. Sync Directory

<check if="path is a directory">
  <action>Sync directory with rsync:</action>
  ```bash
  rsync -av --delete "{{current_path}}/{{path}}/" "{{tracking_path}}/{{path}}/"
  ```

  <action>Count synced files:</action>
  ```bash
  find "{{tracking_path}}/{{path}}" -type f | wc -l
  ```

  <action>Record result:</action>
  - Path: {{path}}
  - Type: directory
  - Status: SUCCESS or FAILED
  - File count: {{count}}
</check>

#### 2b. Sync File

<check if="path is a file">
  <action>Copy file:</action>
  ```bash
  mkdir -p "$(dirname "{{tracking_path}}/{{path}}")"
  cp "{{current_path}}/{{path}}" "{{tracking_path}}/{{path}}"
  ```

  <action>Verify copy:</action>
  ```bash
  test -f "{{tracking_path}}/{{path}}" && echo "SUCCESS" || echo "FAILED"
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
  WARNING: SOME OPERATIONS FAILED

  Failed paths:
  {{for each failed result}}
    - {{path}}: {{error_message}}
  {{end for}}

  Successfully synced:
  {{for each successful result}}
    - {{path}}
  {{end for}}
  ```

  <action>Continue to next step despite failures</action>
</check>

### 4. Store Results and Proceed

<action>Store {{sync_results}} for report step</action>
<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After execution:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All sync operations executed
- Results tracked for each path
- Errors reported immediately
- Proceeded to step 10

### SYSTEM FAILURE:
- Not executing all planned syncs
- Not tracking operation results
- Silent failures without reporting
- Stopping on first error instead of continuing
