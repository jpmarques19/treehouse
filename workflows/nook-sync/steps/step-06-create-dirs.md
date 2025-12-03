---
name: 'step-06-create-dirs'
description: 'Create tracking directories if missing'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-06-create-dirs.md'
nextStepFile: '{workflow_path}/steps/step-07-preview.md'
---

# Step 6: Create Tracking Directories

## STEP GOAL:

Ensure all required tracking directories exist before sync operations.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Create parent directories as needed

### Role Reinforcement:

- You are a workspace automation assistant
- Directory creation is prerequisite for sync
- Be thorough with directory structure

## EXECUTION SEQUENCE:

### 1. Create Main Tracking Directory

```bash
mkdir -p "{{tracking_path}}"
```

### 2. Create Subdirectories for Sync Paths

<action>For each path in {{sync_paths}}, create parent directory in tracking:</action>

```bash
# For directories: create the directory itself
# For files: create the parent directory

# Examples:
# docs -> {{tracking_path}}/docs/
# .bmad/_cfg/agents -> {{tracking_path}}/.bmad/_cfg/agents/
# .bmad/bmm/config.yaml -> {{tracking_path}}/.bmad/bmm/

for path in {{sync_paths as space-separated}}; do
  if [[ "$path" == *.* ]]; then
    # It's a file - create parent directory
    mkdir -p "$(dirname "{{tracking_path}}/$path")"
  else
    # It's a directory - create it
    mkdir -p "{{tracking_path}}/$path"
  fi
done
```

### 3. Verify Creation

```bash
ls -la "{{tracking_path}}"
```

### 4. Display Status

<action>Display:</action>
```
Tracking directories ready:
  {{tracking_path}}
```

### 5. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After creation:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Main tracking directory created
- Subdirectories for all sync paths created
- Directory structure verified
- Proceeded to step 7

### SYSTEM FAILURE:
- Directory creation failed
- Missing parent directories
- Not handling file vs directory paths
