---
name: 'step-04-load-config'
description: 'Load path configurations from config'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-04-load-config.md'
nextStepFile: '{workflow_path}/steps/step-05-update-gitignore.md'
---

# Step 4: Load Path Configurations

## STEP GOAL:

Load the path configurations that define which artifacts to track and how.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Parse all path configuration sections
- Store arrays for subsequent steps

### Role Reinforcement:

- You are a workspace automation assistant
- Configuration loading is critical for correct setup
- Be precise with path parsing

## EXECUTION SEQUENCE:

### 1. Read Module Configuration

```bash
cat .bmad/th/config.yaml
```

### 2. Extract Path Configurations

<action>Parse and store the following from config:</action>

**gitignore_paths** - Paths to add to .gitignore:
```yaml
gitignore_paths:
  - docs
```

**skip_worktree_paths** - Paths to mark with skip-worktree in nooks:
```yaml
skip_worktree_paths:
  - .bmad/_cfg/agents
  - .bmad/bmm/config.yaml
```

**sync_paths** - All paths for sync operations:
```yaml
sync_paths:
  - docs
  - .bmad/_cfg/agents
  - .bmad/bmm/config.yaml
```

<action>Store as {{gitignore_paths}}, {{skip_worktree_paths}}, {{sync_paths}}</action>

### 3. Extract Tracking Folder

<action>Get tracking_folder (default: .bmad-tracking):</action>
```bash
grep "^tracking_folder:" .bmad/th/config.yaml | awk '{print $2}' || echo ".bmad-tracking"
```

<action>Store as {{tracking_folder}}</action>

### 4. Display Configuration

<action>Display:</action>
```
Path Configuration Loaded:

Gitignore paths (added to .gitignore):
{{for each path in gitignore_paths}}
  - {{path}}
{{end for}}

Skip-worktree paths (tracked, local changes ignored in nooks):
{{for each path in skip_worktree_paths}}
  - {{path}}
{{end for}}

Tracking folder: {{tracking_folder}}
```

### 5. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After loading:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All path configs extracted
- Arrays stored correctly
- Tracking folder determined
- Proceeded to step 5

### SYSTEM FAILURE:
- Missing path configurations
- Not distinguishing gitignore from skip_worktree paths
- Incorrect array parsing
