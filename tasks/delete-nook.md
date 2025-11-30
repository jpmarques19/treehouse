# Task: delete-nook

Delete a specific nook by removing its worktree and/or tracking data.

**Architecture Note:** Tracking data lives in `base_workspace_path/.bmad-tracking/`

## Input

- `target`: Branch name to delete (required)
- `mode`: What to delete - `all` (default), `worktree-only`, `tracking-only`

## Execution

### Step 0: Get base workspace path

```bash
# Read from config
base_path=$(grep "^base_workspace_path:" .bmad/th/config.yaml | sed 's/base_workspace_path: *//' | tr -d '"')

if [ -z "$base_path" ]; then
  echo "ERROR: Workspace not initialized. Run /bmad:th:workflows:treehouse-init first."
  exit 1
fi

# Verify base path exists
if [ ! -d "$base_path" ]; then
  echo "ERROR: Base workspace not found at: $base_path"
  exit 1
fi
```

Store: `{{base_path}}`

### Step 1: Validate target exists

```bash
# Check if branch exists
git branch --list "{{target}}"

# Check if worktree exists
git worktree list | grep "{{target}}"

# Check if tracking exists (in base workspace)
ls -la "{{base_path}}/.bmad-tracking/{{target}}/context.yaml" 2>/dev/null
```

Store findings:
- `{{has_worktree}}`: true/false
- `{{worktree_path}}`: path if exists
- `{{has_tracking}}`: true/false
- `{{tracking_path}}`: `{{base_path}}/.bmad-tracking/{{target}}`

### Step 2: Prevent deleting active/open nook

This step checks TWO conditions that would block deletion:

**Check 1: Current branch matches target**
```bash
current_branch=$(git branch --show-current)
if [ "$current_branch" = "{{target}}" ]; then
  echo "ERROR: Cannot delete - you are checked out on this branch!"
  echo ""
  echo "You are currently on branch '{{target}}'."
  echo "Switch to a different worktree before deleting."
  exit 1
fi
```

**Check 2: Current working directory is inside target worktree**
```bash
current_dir=$(pwd)
if [ "{{has_worktree}}" = "true" ]; then
  # Check if current directory is inside the target worktree
  case "$current_dir" in
    "{{worktree_path}}"*)
      echo "ERROR: Cannot delete - you are inside this worktree!"
      echo ""
      echo "Current directory: $current_dir"
      echo "Target worktree:   {{worktree_path}}"
      echo ""
      echo "Navigate outside this worktree before deleting:"
      echo "  cd {{base_path}}"
      exit 1
      ;;
  esac
fi
```

**Display halt message if either check fails:**
```
CANNOT DELETE - NOOK IS OPEN

You cannot delete a nook that is currently "open" in any form:
- Checked out as current branch
- Working directory inside the worktree

Current location: {{current_dir}}
Target branch:    {{target}}
{{if has_worktree}}
Target worktree:  {{worktree_path}}
{{end if}}

To delete this nook:
1. Navigate to a different worktree: cd {{base_path}}
2. Then retry the delete operation
```

### Step 3: Show what will be deleted

```
Delete Nook: {{target}}

Base workspace: {{base_path}}

Will delete:
{{if has_worktree}}
- Worktree: {{worktree_path}}
  All local changes in this worktree will be LOST
{{end if}}
{{if has_tracking}}
- Tracking: {{tracking_path}}
  Synced context (docs, configs) will be LOST
{{end if}}

{{if has_children}}
WARNING: This nook has child nooks:
{{for each child in children}}
   - {{child.branch}}
{{end for}}
Deleting the parent won't delete children, but lineage will be broken.
{{end if}}

This action cannot be undone.
```

### Step 4: Confirm deletion

```
Type the branch name to confirm deletion:
```

Wait for input. Must match `{{target}}` exactly.

### Step 5: Execute deletion

**Delete worktree (if exists and mode allows):**
```bash
# Check for uncommitted changes
cd "{{worktree_path}}"
if [ -n "$(git status --porcelain)" ]; then
  echo "WARNING: Worktree has uncommitted changes!"
  read -p "Force delete anyway? [y/N]: " force
  if [ "$force" != "y" ]; then
    exit 1
  fi
fi
cd -

# Remove worktree
git worktree remove "{{worktree_path}}" --force
```

**Delete tracking (if exists and mode allows):**
```bash
# Remove from base workspace tracking
rm -rf "{{base_path}}/.bmad-tracking/{{target}}"

# Note: User should commit this change in base workspace
echo "Tracking removed. Remember to commit in base workspace:"
echo "  cd {{base_path}}"
echo "  git add .bmad-tracking/"
echo "  git commit -m 'chore(th): delete tracking for {{target}}'"
```

**Delete branch (optional, ask user):**
```bash
# Only if worktree was removed
git branch -d "{{target}}"  # or -D to force
```

### Step 6: Report completion

```
Nook Deleted

{{if deleted_worktree}}
- Worktree removed: {{worktree_path}}
{{end if}}
{{if deleted_tracking}}
- Tracking removed: {{tracking_path}}
{{end if}}
{{if deleted_branch}}
- Branch deleted: {{target}}
{{end if}}

{{if deleted_tracking}}
IMPORTANT: Commit tracking changes in base workspace:
  cd {{base_path}}
  git add .bmad-tracking/
  git commit -m "chore(th): delete tracking for {{target}}"
{{end if}}

{{if kept_branch}}
Note: Branch "{{target}}" was kept. Delete manually with:
  git branch -d {{target}}
{{end if}}
```

## Safety Checks

1. **Cannot delete if checked out** - branch is your current branch
2. **Cannot delete if inside worktree** - working directory is inside target worktree path
3. **Uncommitted changes warning** - requires force flag
4. **Confirmation by typing name** - prevents accidental deletion
5. **Child nook warning** - informs about broken lineage
6. **Tracking in base workspace** - reminds to commit changes
