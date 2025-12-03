---
name: 'step-05-check-content'
description: 'Check existing local content that may be overwritten'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'

# File References
thisStepFile: '{workflow_path}/steps/step-05-check-content.md'
nextStepFile: '{workflow_path}/steps/step-06-preview-restore.md'
---

# Step 5: Check Existing Local Content

## STEP GOAL:

Check what local content exists that will be affected by the restore operation.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Document all existing content for user awareness

### Role Reinforcement:

- You are a workspace automation assistant
- Help user understand what will be affected
- Be thorough in content discovery

## EXECUTION SEQUENCE:

### 1. Check Each Sync Path Locally

<action>For each path in {{sync_paths}}, check local status:</action>

```bash
# For each path, determine:
# 1. Does it exist?
# 2. Is it a file or directory?
# 3. If directory, how many files?

for path in {{sync_paths as space-separated list}}; do
  if [ -d "$path" ]; then
    count=$(find "$path" -type f 2>/dev/null | wc -l)
    echo "DIR:$path:$count files"
  elif [ -f "$path" ]; then
    echo "FILE:$path:exists"
  else
    echo "MISSING:$path"
  fi
done
```

### 2. Build Content Report

<action>Parse results and build {{existing_content_report}}:</action>

For each path, record:
- **Status**: EXISTS_DIR, EXISTS_FILE, or MISSING
- **Details**: File count for directories, "exists" for files
- **Will be affected**: Yes if local content exists

### 3. Check What Exists in Tracking

<action>Also check what exists in the tracking folder:</action>

```bash
for path in {{sync_paths as space-separated list}}; do
  tracking_loc="{{tracking_path}}/$path"
  if [ -d "$tracking_loc" ]; then
    count=$(find "$tracking_loc" -type f 2>/dev/null | wc -l)
    echo "TRACKING_DIR:$path:$count files"
  elif [ -f "$tracking_loc" ]; then
    echo "TRACKING_FILE:$path:exists"
  else
    echo "TRACKING_MISSING:$path"
  fi
done
```

### 4. Build Restore Plan

<action>Create {{restore_plan}} combining local and tracking info:</action>

| Path | Local Status | Tracking Status | Action |
|------|-------------|-----------------|--------|
| For each path... | | | |

Actions:
- **RESTORE**: Tracking has content, will copy to local
- **SKIP**: Tracking missing this path
- **OVERWRITE**: Both have content, tracking will replace local

### 5. Store Results and Proceed

<action>Store {{existing_content_report}} and {{restore_plan}} for next step</action>
<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After checking:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All sync paths checked locally
- All sync paths checked in tracking
- Content report built with accurate counts
- Restore plan created
- Proceeded to step 6

### SYSTEM FAILURE:
- Not checking all paths
- Inaccurate content counts
- Not identifying overwrite situations
