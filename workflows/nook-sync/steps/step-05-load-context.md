---
name: 'step-05-load-context'
description: 'Load existing context if available'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-05-load-context.md'
nextStepFile: '{workflow_path}/steps/step-06-create-dirs.md'
---

# Step 5: Load Existing Context

## STEP GOAL:

Check for and load existing context.yaml if available, extracting lineage for preservation.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Preserve existing lineage if found
- Initialize empty lineage if new

### Role Reinforcement:

- You are a workspace automation assistant
- Context preservation is important
- Handle both new and existing contexts

## EXECUTION SEQUENCE:

### 1. Check for Existing Context

```bash
test -f "{{tracking_path}}/context.yaml" && echo "EXISTS" || echo "NOT_FOUND"
```

### 2. Handle Existing Context

<check if="context.yaml exists">
  <action>Load existing context.yaml:</action>
  ```bash
  cat "{{tracking_path}}/context.yaml"
  ```

  <action>Extract and store:</action>
  - {{existing_lineage}} = lineage array from context
  - {{existing_artifacts}} = artifacts section
  - {{existing_summary}} = summary section (if present)

  <action>Display:</action>
  ```
  Existing context found:
  - Lineage entries: {{existing_lineage | length}}
  - Last synced: {{synced_at from context}}
  ```
</check>

### 3. Handle New Context

<check if="context.yaml does not exist">
  <action>Initialize empty structures:</action>
  - {{existing_lineage}} = empty array
  - {{existing_artifacts}} = empty
  - {{existing_summary}} = empty

  <action>Extract type from branch name:</action>
  - If branch has prefix (e.g., "discovery/foo") → type = "discovery"
  - If branch has no prefix → type = "base"

  <action>Store as {{branch_type}}</action>

  <action>Display:</action>
  ```
  No existing context - this will be a new sync.
  Branch type: {{branch_type}}
  ```
</check>

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After loading:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Existing context checked
- Lineage extracted if available
- Empty structures initialized if new
- Branch type determined
- Proceeded to step 6

### SYSTEM FAILURE:
- Not checking for existing context
- Losing existing lineage
- Not initializing new context properly
