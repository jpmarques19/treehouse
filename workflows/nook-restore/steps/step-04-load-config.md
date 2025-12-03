---
name: 'step-04-load-config'
description: 'Load sync_paths configuration from module config'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-04-load-config.md'
nextStepFile: '{workflow_path}/steps/step-05-check-content.md'
---

# Step 4: Load Sync Paths Configuration

## STEP GOAL:

Load the sync_paths configuration that defines which paths should be restored.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Ensure sync_paths is properly loaded before proceeding

### Role Reinforcement:

- You are a workspace automation assistant
- Configuration loading is critical for correct restore
- Be precise with path handling

## EXECUTION SEQUENCE:

### 1. Load Module Configuration

<action>Read full module config:</action>
```bash
cat .bmad/th/config.yaml
```

### 2. Extract Sync Paths

<action>Parse the sync_paths section from config. Expected format:</action>
```yaml
sync_paths:
  - docs
  - .bmad/_cfg/agents
  - .bmad/bmm/config.yaml
  # ... additional paths
```

<action>Store as {{sync_paths}} array</action>

<check if="sync_paths is empty or not defined">
  <action>Display error:</action>
  ```
  ERROR: No sync_paths defined in config

  The sync_paths configuration is required to know what to restore.

  Please add sync_paths to .bmad/th/config.yaml:

  sync_paths:
    - docs
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml
  ```
  <action>HALT workflow</action>
</check>

### 3. Display Configured Paths

<action>Display the paths that will be considered for restore:</action>
```
Configured sync paths:
{{for each path in sync_paths}}
  - {{path}}
{{end for}}
```

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After loading:
- If sync_paths loaded: Load and execute `{nextStepFile}`
- If sync_paths missing: HALT with error

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Module config loaded
- sync_paths extracted and stored as array
- Paths displayed to user
- Proceeded to step 5

### SYSTEM FAILURE:
- Not loading sync_paths properly
- Proceeding with empty sync_paths
- Not handling missing configuration
