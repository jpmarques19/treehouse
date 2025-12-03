---
name: 'step-05-update-gitignore'
description: 'Update .gitignore with gitignore_paths only'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'

# File References
thisStepFile: '{workflow_path}/steps/step-05-update-gitignore.md'
nextStepFile: '{workflow_path}/steps/step-06-remove-from-cache.md'
---

# Step 5: Update .gitignore

## STEP GOAL:

Add gitignore_paths to .gitignore (NOT skip_worktree_paths).

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Only add gitignore_paths, NOT skip_worktree_paths
- Check for existing section before adding

### Role Reinforcement:

- You are a workspace automation assistant
- Precise gitignore management is important
- Avoid duplicate entries

## EXECUTION SEQUENCE:

### 1. Check .gitignore Exists

```bash
test -f .gitignore && echo "EXISTS" || echo "NOT_FOUND"
```

<action if=".gitignore does not exist">Create empty .gitignore</action>

### 2. Check for Existing Section

```bash
grep -q "# Treehouse" .gitignore 2>/dev/null && echo "EXISTS" || echo "NOT_FOUND"
```

### 3. Add Section If Missing

<check if="Treehouse section does not exist">
  <action>Append to .gitignore (ONLY gitignore_paths):</action>

```gitignore

# Treehouse (th) - Tracked artifacts
# These paths are per-workspace and synced via th workflows
# Do not remove - managed by /bmad:th:workflows:treehouse-init
docs/
```

  <action>Note: Add trailing `/` for directories</action>
  <action>Note: {{skip_worktree_paths}} are NOT gitignored - they use skip-worktree in nooks</action>
  <action>Note: .bmad-tracking/ is NOT gitignored - it's a submodule</action>

  <action>Display:</action>
  ```
  Added to .gitignore:
  {{for each path in gitignore_paths}}
    - {{path}}/
  {{end for}}
  ```
</check>

<check if="Treehouse section already exists">
  <action>Display: ".gitignore already has Treehouse section - skipping"</action>
</check>

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After update:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- .gitignore checked/created
- Existing section detected
- gitignore_paths added (not skip_worktree_paths)
- Proceeded to step 6

### SYSTEM FAILURE:
- Adding skip_worktree_paths to .gitignore
- Creating duplicate sections
- Not checking for existing section
