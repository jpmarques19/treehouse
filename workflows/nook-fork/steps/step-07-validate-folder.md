---
name: 'step-07-validate-folder'
description: 'Validate that the worktrees folder exists and is configured'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-07-validate-folder.md'
nextStepFile: '{workflow_path}/steps/step-08-create-worktree.md'
---

# Step 7: Validate Worktrees Folder

## STEP GOAL:

Ensure the worktrees destination folder exists and is configured in th/config.yaml.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 Create folder if it doesn't exist
- ⏸️ Ask user for path if not configured

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Handle missing configuration gracefully
- ✅ Validate folder is writable

## EXECUTION SEQUENCE:

### 1. Read Worktrees Folder Config

<action>Read worktrees_folder from {module_config}</action>

### 2. Check Configuration

<check if="worktrees_folder is NOT configured or empty">
  <action>Ask user for path:</action>

```
WORKTREES FOLDER NOT CONFIGURED

Where should nook worktrees be created?

Common locations:
- ../    (sibling folder to current workspace)
- ~/dev/ (home directory dev folder)
- /path/to/worktrees (custom location)

Enter path:
```

  <action>HALT and wait for user input</action>
  <action>Store user path as {{worktrees_folder}}</action>

  <action>Update config file:</action>
  ```
  Add/update worktrees_folder in .bmad/th/config.yaml
  ```
</check>

<check if="worktrees_folder is configured">
  <action>Store as {{worktrees_folder}}</action>
</check>

### 3. Validate Folder Exists

```bash
test -d "{{worktrees_folder}}" && echo "EXISTS" || echo "NOT_FOUND"
```

<check if="folder does NOT exist">
  <action>Create it:</action>
  ```bash
  mkdir -p "{{worktrees_folder}}"
  ```
  <action>Display: "Created worktrees folder: {{worktrees_folder}}"</action>
</check>

<check if="folder exists">
  <action>Display: "Worktrees folder: {{worktrees_folder}}"</action>
</check>

### 4. Validate Writable

```bash
test -w "{{worktrees_folder}}" && echo "WRITABLE" || echo "NOT_WRITABLE"
```

<check if="folder is NOT writable">
  <action>Display error:</action>
  ```
  ERROR: Cannot write to worktrees folder
  Path: {{worktrees_folder}}

  Please check permissions and try again.
  ```
  <action>HALT workflow</action>
</check>

### 5. Continue

<action>Display: "Worktrees folder validated: {{worktrees_folder}}"</action>
<action>Auto-proceed to next step</action>

## MENU OPTIONS:

After validation: Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Worktrees folder path resolved (from config or user)
- Folder exists (created if needed)
- Folder is writable
- Config updated if user provided new path
- {{worktrees_folder}} stored

### ❌ SYSTEM FAILURE:
- Proceeding with non-existent folder
- Not validating write permissions
- Not asking user when config is missing
