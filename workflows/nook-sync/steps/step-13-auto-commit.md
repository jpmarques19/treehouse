---
name: 'step-13-auto-commit'
description: 'Optional auto-commit for submodule'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-13-auto-commit.md'
nextStepFile: '{workflow_path}/steps/step-14-report.md'
---

# Step 13: Auto-Commit (Optional)

## STEP GOAL:

If auto_commit_tracking is enabled, commit changes to both the submodule and parent repo.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Check auto_commit setting before executing
- Handle both submodule and parent commits

### Role Reinforcement:

- You are a workspace automation assistant
- Auto-commit is optional based on config
- Follow git submodule commit pattern

## EXECUTION SEQUENCE:

### 1. Check Auto-Commit Setting

<check if="auto_commit == true">
  <action>Display:</action>
  ```
  Auto-commit enabled - committing tracking changes...

  With submodule architecture, this requires TWO commits:
  1. Inside .bmad-tracking/ submodule (the actual file changes)
  2. In parent repo (update submodule pointer)
  ```

  <goto step="2">Execute commits</goto>
</check>

<check if="auto_commit != true or not set">
  <action>Store {{manual_commit_needed}} = true</action>
  <action>Display: "Auto-commit disabled - manual commit will be required"</action>
  <action>Auto-proceed to next step</action>
</check>

### 2. Commit Inside Submodule

<action>Navigate to tracking submodule and commit:</action>
```bash
cd "{{base_path}}/.bmad-tracking"
git add .
git commit -m "sync: {{current_branch}} context

Synced from: {{current_path}}
Hash: {{current_hash}}

Generated with Claude Code" || true
```

### 3. Update Parent Pointer

<action>Navigate to parent repo and update submodule pointer:</action>
```bash
cd "{{base_path}}"
git add .bmad-tracking
git commit -m "chore(th): update tracking pointer for {{current_branch}}

Generated with Claude Code" || true
```

### 4. Confirm Commits

<action>Display:</action>
```
Auto-commit complete:
  - Submodule changes committed
  - Parent pointer updated
```

<action>Store {{manual_commit_needed}} = false</action>

### 5. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After commit handling:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Auto-commit setting checked
- If enabled: both commits executed
- If disabled: manual_commit_needed flag set
- Proceeded to step 14

### SYSTEM FAILURE:
- Not checking auto_commit setting
- Partial commits (submodule without parent)
- Silent commit failures
