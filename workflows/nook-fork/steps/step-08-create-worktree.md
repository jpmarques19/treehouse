---
name: 'step-08-create-worktree'
description: 'Create the git worktree for the new nook'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-08-create-worktree.md'
nextStepFile: '{workflow_path}/steps/step-09-configure-isolation.md'
---

# Step 8: Create Worktree

## STEP GOAL:

Create the git worktree for the new nook, establishing the isolated development environment.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- ⏸️ HALT if target path already exists
- ⏸️ Get user confirmation before creating

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Check for conflicts before creating
- ✅ Display clear summary before execution

## EXECUTION SEQUENCE:

### 1. Calculate Worktree Path

<action>Derive folder name from branch:</action>
```
Folder: {{parent_hash}}-{{nook_name}}
Full path: {{worktrees_folder}}/{{parent_hash}}-{{nook_name}}
```

<action>Store as {{worktree_path}}</action>

### 2. Check Path Doesn't Exist

```bash
test -e "{{worktree_path}}" && echo "EXISTS" || echo "AVAILABLE"
```

<check if="path already EXISTS">
  <action>Display error:</action>

```
ERROR - PATH ALREADY EXISTS

The target worktree path already exists:
{{worktree_path}}

This could mean:
- A nook with this name was created before
- There's a stale worktree that wasn't cleaned up

Options:
- Choose a different nook name
- Remove the existing folder: rm -rf "{{worktree_path}}"
- Check existing worktrees: git worktree list
```

  <action>HALT workflow</action>
</check>

### 3. Display Summary and Confirm

```
READY TO CREATE NOOK

Summary:
  Parent branch:    {{current_branch}} ({{parent_hash}})
  New branch:       {{new_branch}}
  Worktree path:    {{worktree_path}}
  Agent mode:       {{agent_mode}}

This will:
1. Create new git branch: {{new_branch}}
2. Create worktree at: {{worktree_path}}
3. Apply isolation configuration

[Y] Create nook
[N] Cancel

Proceed? [Y/N]:
```

<action>HALT and wait for user confirmation</action>

<check if="user selects N">
  <action>Display: "Cancelled. No changes made."</action>
  <action>HALT workflow</action>
</check>

### 4. Create Worktree

<check if="user selects Y">
  <action>Display: "Creating worktree..."</action>

  <action>Execute git worktree add:</action>
  ```bash
  git worktree add "{{worktree_path}}" -b "{{new_branch}}" "{{current_branch}}"
  ```

  <action>Verify success:</action>
  ```bash
  git worktree list | grep "{{worktree_path}}"
  ```

  <check if="worktree creation failed">
    <action>Display error with git output</action>
    <action>HALT workflow</action>
  </check>

  <action>Display: "Worktree created successfully at {{worktree_path}}"</action>
  <action>Proceed to next step</action>
</check>

## MENU OPTIONS:

After confirmation:
- If Y: Create worktree and proceed to `{nextStepFile}`
- If N: HALT workflow

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Target path verified as available
- User confirmed creation
- Git worktree created successfully
- New branch created from current branch
- {{worktree_path}} stored and verified

### ❌ SYSTEM FAILURE:
- Creating worktree without confirmation
- Proceeding when target path exists
- Not checking git command result
- Not verifying worktree in list
