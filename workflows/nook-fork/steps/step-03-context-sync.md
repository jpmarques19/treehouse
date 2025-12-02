---
name: 'step-03-context-sync'
description: 'Check if current branch has synced context for lineage inheritance'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-03-context-sync.md'
nextStepFile: '{workflow_path}/steps/step-04-nook-type.md'
---

# Step 3: Check Context Sync Status

## STEP GOAL:

Check if the current branch has a context.yaml file for lineage inheritance. This is informational - not blocking.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 This is an informational check, not a blocker
- ✅ Extract existing lineage if available

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Inform user about lineage inheritance status
- ✅ Don't block - just inform

## EXECUTION SEQUENCE:

### 1. Check for Existing Context

```bash
test -f "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml" && echo "EXISTS" || echo "NOT_FOUND"
```

### 2. Handle Context Status

<check if="context.yaml does NOT exist for current branch">
  <action>Display informational notice:</action>

```
NOTE - NO SYNCED CONTEXT FOR CURRENT BRANCH

Branch '{{current_branch}}' hasn't been synced yet.
This means the new nook won't inherit any context lineage.

This is fine if:
- This is a new workspace just starting out
- You haven't made changes worth preserving yet

If you want to save context first:
  /th:workflows:nook-sync

Continuing with nook creation...
```

  <action>Set {{existing_lineage}} = empty array</action>
  <action>Auto-proceed to next step</action>
</check>

<check if="context.yaml exists">
  <action>Load existing context:</action>
  ```bash
  cat "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml"
  ```

  <action>Extract {{existing_lineage}} array from the lineage section</action>
  <action>Display: "Lineage context found - will be inherited by new nook"</action>
  <action>Auto-proceed to next step</action>
</check>

## MENU OPTIONS:

This is an auto-proceed step (informational only).
After checking: Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Context file existence checked
- User informed of lineage status
- Existing lineage extracted if available
- Proceeded to step 4

### ❌ SYSTEM FAILURE:
- Blocking on missing context (should be informational)
- Not extracting lineage when available
- Not informing user of inheritance status
