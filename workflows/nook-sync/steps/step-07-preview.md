---
name: 'step-07-preview'
description: 'Show sync preview'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-07-preview.md'
nextStepFile: '{workflow_path}/steps/step-08-confirm.md'
---

# Step 7: Show Sync Preview

## STEP GOAL:

Show the user exactly what will be synced to the tracking folder.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Check which paths exist in current directory
- Display clear preview of sync operation

### Role Reinforcement:

- You are a workspace automation assistant
- User needs visibility before confirming
- Be thorough in path checking

## EXECUTION SEQUENCE:

### 1. Check Source Paths

<action>For each path in {{sync_paths}}, check if it exists in current directory:</action>

```bash
for path in {{sync_paths as space-separated}}; do
  if [ -e "$path" ]; then
    if [ -d "$path" ]; then
      count=$(find "$path" -type f 2>/dev/null | wc -l)
      echo "EXISTS_DIR:$path:$count"
    else
      echo "EXISTS_FILE:$path"
    fi
  else
    echo "MISSING:$path"
  fi
done
```

<action>Store results as {{source_status}} for each path</action>

### 2. Display Sync Preview

<action>Display preview:</action>

```
============================================================
                     SYNC PREVIEW
============================================================

Source:  {{current_path}} ({{current_branch}})
Target:  {{tracking_path}}
Hash:    {{current_hash}}

------------------------------------------------------------
Paths to sync:
------------------------------------------------------------

{{for each path in sync_paths}}
{{if path exists as directory}}
  [EXISTS] {{path}}/ ({{file_count}} files)
           -> {{tracking_path}}/{{path}}/

{{else if path exists as file}}
  [EXISTS] {{path}}
           -> {{tracking_path}}/{{path}}

{{else}}
  [SKIP]   {{path}} (not found in current directory)

{{end if}}
{{end for}}

------------------------------------------------------------
WARNING: Existing content in tracking will be REPLACED (not merged)
------------------------------------------------------------
```

### 3. Handle Nothing to Sync

<check if="no paths exist to sync">
  <action>Display warning:</action>
  ```
  WARNING - NOTHING TO SYNC

  None of the configured sync_paths exist in the current directory:
  {{for each path in sync_paths}}
    - {{path}}
  {{end for}}

  This is expected if:
  - This is a fresh nook that hasn't created any artifacts yet
  - You're in a worktree where sync artifacts are gitignored

  The context.yaml will still be updated with branch info.
  ```

  <action>Store {{nothing_to_sync}} = true</action>
</check>

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After preview:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All sync paths checked for existence
- Clear preview displayed
- Warning shown if nothing to sync
- User informed of overwrite behavior
- Proceeded to step 8

### SYSTEM FAILURE:
- Not checking all paths
- Unclear preview
- Not warning about overwrites
