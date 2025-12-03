---
name: 'step-06-remove-from-cache'
description: 'Remove gitignore_paths from git cache (skip_worktree_paths stay tracked)'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'

# File References
thisStepFile: '{workflow_path}/steps/step-06-remove-from-cache.md'
nextStepFile: '{workflow_path}/steps/step-07-update-config.md'
---

# Step 6: Remove from Git Cache

## STEP GOAL:

Remove gitignore_paths from git cache so they're truly untracked. Skip_worktree_paths MUST stay tracked.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- ONLY remove gitignore_paths from cache
- NEVER remove skip_worktree_paths from cache

### Role Reinforcement:

- You are a workspace automation assistant
- Cache management is critical for correct behavior
- Preserve local files while removing from index

## EXECUTION SEQUENCE:

### 1. Explain Operation

<action>Display:</action>
```
Removing gitignored paths from git cache...

Adding to .gitignore alone doesn't stop tracking files that are
already in the git index. We must remove them from cache.

Note: skip_worktree_paths (.bmad/_cfg/agents, .bmad/bmm/config.yaml)
stay tracked - they use skip-worktree in nooks instead.
```

### 2. Remove gitignore_paths from Cache

<action>For each path in {{gitignore_paths}} ONLY, remove from git cache:</action>

```bash
# Only remove gitignore_paths (docs), NOT skip_worktree_paths
git rm -r --cached docs/ 2>/dev/null || true
```

<action>Note: `|| true` ensures command doesn't fail if path isn't tracked</action>
<action>Note: `--cached` removes from index only, preserving local files</action>
<action>CRITICAL: Do NOT remove skip_worktree_paths from cache!</action>

### 3. Report Results

```bash
# Check if anything was unstaged
git status --porcelain | grep "^D " | head -10
```

<check if="files were removed from cache">
  <action>Display:</action>
  ```
  Removed from git tracking (files preserved locally):
  {{list removed files}}
  ```
</check>

<check if="no files were removed">
  <action>Display: "No tracked files to remove (already clean)"</action>
</check>

### 4. Auto-Proceed

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After removal:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- gitignore_paths removed from cache
- skip_worktree_paths NOT touched
- Local files preserved
- Proceeded to step 7

### SYSTEM FAILURE:
- Removing skip_worktree_paths from cache
- Deleting local files (not using --cached)
- Not reporting what was removed
