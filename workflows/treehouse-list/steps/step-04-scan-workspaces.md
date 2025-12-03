---
name: 'step-04-scan-workspaces'
description: 'Scan tracking folder for all tracked workspaces'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-list'

# File References
thisStepFile: '{workflow_path}/steps/step-04-scan-workspaces.md'
nextStepFile: '{workflow_path}/steps/step-05-scan-worktrees.md'
---

# Step 4: Scan Tracked Workspaces

## STEP GOAL:

Find all context.yaml files in the tracking folder and extract workspace metadata.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- CRITICAL: Read the complete step file before taking any action
- Execute bash commands exactly as specified
- Handle empty tracking folder gracefully

### Role Reinforcement:

- You are a workspace navigation assistant
- Scan silently, no user interaction
- Auto-proceed when done

## EXECUTION SEQUENCE:

### 1. Find All Context Files

```bash
find "{{base_path}}/.bmad-tracking" -name "context.yaml" -type f 2>/dev/null
```

### 2. Handle Empty Result

<check if="no context.yaml files found">
  <action>Display:</action>

```
🌳 Workspace Lineage

no workspaces tracked yet

base {{base_path}}
here {{current_branch}} · {{current_hash}}

run /bmad:th:workflows:nook-sync to sync current context
```

  <action>Exit workflow</action>
</check>

### 3. Parse Each Context File

<check if="context.yaml files found">
  <action>For each context.yaml found:</action>

  Read each file and extract:
  - `branch` - the branch name
  - `hash` - the commit hash (4 chars)
  - `type` - nook type (feature, bugfix, spike, etc.)
  - `forked_from` - parent branch name
  - `lineage` - ancestry chain
  - `synced_at` - last sync timestamp

  Derive tracking path from context.yaml location:
  - Path relative to `.bmad-tracking/`
  - Example: `.bmad-tracking/feature/abc1-my-feature/context.yaml` -> `feature/abc1-my-feature`

  Store all workspace data in `{{workspaces}}` array.

  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

- If no workspaces: Exit with message
- If workspaces found: Auto-proceed to `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- All context.yaml files discovered
- Metadata extracted from each
- Workspaces array populated
- Proceeded to step 5

### SYSTEM FAILURE:
- Not handling empty tracking folder
- Not extracting all required fields
- Displaying intermediate output to user
