---
name: 'step-02-branch-detect'
description: 'Detect current branch and determine restore source'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'

# File References
thisStepFile: '{workflow_path}/steps/step-02-branch-detect.md'
nextStepFile: '{workflow_path}/steps/step-03-validate-tracking.md'
---

# Step 2: Detect Branch and Determine Source

## STEP GOAL:

Detect current branch and determine which branch's context to restore from.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Wait for user input when source selection is needed

### Role Reinforcement:

- You are a workspace automation assistant
- Help user select appropriate restore source
- Be clear about options and implications

## EXECUTION SEQUENCE:

### 1. Detect Current Branch

```bash
git branch --show-current
```

<action>Store result as {{current_branch}}</action>

```bash
pwd
```

<action>Store result as {{current_path}}</action>

<action>Display: "Current branch: {{current_branch}}"</action>

### 2. Check for Provided Source Branch

<check if="source_branch was provided as workflow input">
  <action>Use provided source_branch as {{restore_branch}}</action>
  <action>Display: "Using provided source: {{restore_branch}}"</action>
  <action>Auto-proceed to next step</action>
</check>

### 3. Check If Current Branch Has Context

<check if="source_branch NOT provided">
  <action>Check if current branch has context in tracking:</action>
  ```bash
  test -d "{{base_path}}/.bmad-tracking/{{current_branch}}" && echo "EXISTS" || echo "NOT_FOUND"
  ```

  <check if="context exists for current branch">
    <action>Default to current branch: {{restore_branch}} = {{current_branch}}</action>
    <action>Display: "Context found for current branch"</action>
    <action>Auto-proceed to next step</action>
  </check>
</check>

### 4. Handle No Context for Current Branch

<check if="context does NOT exist for current branch">
  <action>Check if this is a nook by looking for parent in context.yaml:</action>
  ```bash
  # Check all context.yaml files for one that lists current branch as forked_from
  grep -l "forked_from: {{current_branch}}" {{base_path}}/.bmad-tracking/*/context.yaml 2>/dev/null || echo "NO_CHILDREN"
  ```

  <action>Also check if any context has current branch in lineage - indicating parent:</action>
  ```bash
  # Look for context that might be our parent
  ls -1 {{base_path}}/.bmad-tracking/ 2>/dev/null
  ```

  <action>Display source selection menu:</action>
  ```
  SELECT RESTORE SOURCE

  Current branch '{{current_branch}}' has no synced context.

  Available contexts in tracking:
  {{list directories from base_path/.bmad-tracking/}}

  Options:
  [L] List available contexts with details
  [E] Enter branch name manually
  [C] Cancel restore

  Choose:
  ```

  <check if="user chooses L (list)">
    <action>Show detailed listing:</action>
    ```bash
    for dir in {{base_path}}/.bmad-tracking/*/; do
      branch=$(basename "$dir")
      if [ -f "$dir/context.yaml" ]; then
        synced=$(grep "synced_at:" "$dir/context.yaml" | cut -d: -f2-)
        echo "$branch (synced:$synced)"
      else
        echo "$branch (no context.yaml)"
      fi
    done
    ```
    <ask>Enter branch name to restore from:</ask>
    <action>Set {{restore_branch}} from user input</action>
  </check>

  <check if="user chooses E (enter)">
    <ask>Enter branch name to restore from:</ask>
    <action>Set {{restore_branch}} from user input</action>
  </check>

  <check if="user chooses C (cancel)">
    <action>Display: "Restore cancelled."</action>
    <action>HALT workflow</action>
  </check>
</check>

## MENU OPTIONS:

- If source branch provided or context exists: Auto-proceed to `{nextStepFile}`
- If no context found: Present selection menu and wait for user input

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Current branch detected
- Restore source determined (from input, current branch, or user selection)
- {{restore_branch}} variable set
- Proceeded to step 3

### SYSTEM FAILURE:
- Not detecting current branch
- Proceeding without determining restore source
- Not offering user selection when no automatic source found
