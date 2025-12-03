---
name: 'step-11-build-lineage'
description: 'Build and update lineage chain'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-11-build-lineage.md'
nextStepFile: '{workflow_path}/steps/step-12-update-manifest.md'
---

# Step 11: Build and Update Lineage

## STEP GOAL:

Build the current branch's lineage entry and update the full lineage chain.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Preserve existing lineage entries
- Update or add current branch entry

### Role Reinforcement:

- You are a workspace automation assistant
- Lineage tracking is critical for context analyst
- Be precise with timestamps and hashes

## EXECUTION SEQUENCE:

### 1. Extract Branch Type and Name

<action>Parse branch name for type and human-readable name:</action>

Examples:
- `discovery/peek-a-box-mvp` → type: "discovery", name: "peek-a-box-mvp"
- `th/feature/agent-wizard` → type: "th", name: "feature/agent-wizard"
- `main` → type: "base", name: "main"
- `feature/auth-system` → type: "feature", name: "auth-system"

<action>Store as {{branch_type}} and {{branch_name}}</action>

### 2. Build Current Entry

<action>Build current branch's lineage entry:</action>

```yaml
current_entry:
  branch: {{current_branch}}
  hash: {{current_hash}}
  type: {{branch_type}}
  name: "{{branch_name}}"
  created: {{existing_created or current_date}}
  updated: {{current_date}}
```

**Note:** If updating existing entry, preserve original `created` date.

### 3. Update Lineage Chain

<check if="current branch already in existing_lineage">
  <action>Update existing entry:</action>
  - Keep original created date
  - Update hash to current
  - Update updated timestamp
</check>

<check if="current branch NOT in existing_lineage">
  <action>Append current_entry to lineage:</action>
  - Add as new entry at end of chain
  - Set created = current_date
</check>

### 4. Store Full Lineage

<action>Store final lineage as {{full_lineage}}</action>

<action>Display lineage summary:</action>
```
Lineage chain ({{full_lineage | length}} entries):
{{for each entry in full_lineage}}
  {{index}}. {{entry.branch}} ({{entry.hash}}) - {{entry.type}}
{{end for}}
```

### 5. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After building:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Branch type and name extracted correctly
- Current entry built with accurate data
- Existing lineage preserved
- Created dates preserved on updates
- Proceeded to step 12

### SYSTEM FAILURE:
- Losing existing lineage entries
- Overwriting created dates
- Incorrect type extraction
- Duplicate entries in lineage
