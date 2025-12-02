---
name: 'step-09-configure-isolation'
description: 'Apply skip-worktree flags and remove submodule artifacts for nook isolation'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'
module_config: '{project-root}/.bmad/th/config.yaml'

# File References
thisStepFile: '{workflow_path}/steps/step-09-configure-isolation.md'
nextStepFile: '{workflow_path}/steps/step-10-agent-deploy.md'
---

# Step 9: Configure Nook Isolation

## STEP GOAL:

Apply skip-worktree flags to config files and remove submodule artifacts to establish nook isolation.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 Execute all isolation steps even if individual ones fail
- ✅ This is an auto-proceed step

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Errors on individual files are warnings, not failures
- ✅ Verify final state after all operations

## EXECUTION SEQUENCE:

### 1. Change to Worktree Directory

```bash
cd "{{worktree_path}}"
```

### 2. Apply Skip-Worktree to Agent Config Files

<action>Read skip_worktree_paths from {module_config}</action>

<action>Apply skip-worktree to configured paths:</action>
```bash
# Apply to agent config folder files
git ls-files .bmad/_cfg/agents/ 2>/dev/null | xargs -r git update-index --skip-worktree

# Apply to bmm config (if it exists)
git update-index --skip-worktree .bmad/bmm/config.yaml 2>/dev/null || true

# Apply any additional configured paths
# (iterate over skip_worktree_paths from config)
```

<action>Display: "Applied skip-worktree to config files"</action>

### 3. Apply Skip-Worktree to Submodule Artifacts

```bash
# Mark submodule files as skip-worktree before removing
git update-index --skip-worktree .bmad-tracking 2>/dev/null || true
git update-index --skip-worktree .gitmodules 2>/dev/null || true
```

<action>Display: "Applied skip-worktree to submodule artifacts"</action>

### 4. Remove Submodule Artifacts from Nook

```bash
# Remove tracking folder (it lives in base workspace only)
rm -rf .bmad-tracking/

# Remove gitmodules file
rm -f .gitmodules
```

<action>Display: "Removed submodule artifacts from nook"</action>

### 5. Verify Clean Git Status

```bash
git status --porcelain
```

<check if="git status shows changes">
  <action>Display warning (but continue):</action>
  ```
  Note: Git status shows some changes after isolation.
  This is usually fine - skip-worktree files may show initially.
  ```
</check>

<check if="git status is clean">
  <action>Display: "Nook has clean git status"</action>
</check>

### 6. Display Isolation Summary

```
ISOLATION CONFIGURED

Applied skip-worktree to:
  - .bmad/_cfg/agents/* (agent config files)
  - .bmad/bmm/config.yaml (if exists)
  - .bmad-tracking (submodule)
  - .gitmodules (submodule config)

Removed from nook:
  - .bmad-tracking/ folder
  - .gitmodules file

Local changes to skip-worktree files won't show in git status.
```

### 7. Continue

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step.
After isolation: Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Skip-worktree applied to config files
- Skip-worktree applied to submodule artifacts
- Submodule artifacts removed from nook
- Git status verified (clean or explained)
- Isolation summary displayed

### ❌ SYSTEM FAILURE:
- Not applying skip-worktree before removing files
- Stopping on individual file errors (should continue)
- Not verifying final state
