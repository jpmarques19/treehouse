# Task: prune-nooks

Clean up orphan worktrees and tracking data across the repository.

**Architecture Note:** All tracking data lives in `base_workspace_path/.bmad-tracking/`

## What Gets Pruned

| Type | Description | Action |
|------|-------------|--------|
| **Broken worktrees** | Git worktree references that point to missing paths | `git worktree prune` |
| **Orphan worktrees** | Worktree exists but no context.yaml in base tracking | Offer to delete |
| **Orphan tracking** | context.yaml exists but branch deleted | Delete tracking |

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

tracking_folder="$base_path/.bmad-tracking"
if [ ! -d "$tracking_folder" ]; then
  echo "ERROR: Tracking folder not found: $tracking_folder"
  exit 1
fi
```

Store: `{{base_path}}`, `{{tracking_folder}}`

### Step 1: Scan for issues

**Find broken worktrees:**
```bash
git worktree list | grep -E '\(error|broken\)'
```

**Find orphan worktrees (no context.yaml in base tracking):**
```bash
# List all worktrees
git worktree list --porcelain | grep "^worktree " | cut -d' ' -f2- > /tmp/worktrees.txt

# For each worktree, check if context.yaml exists in base tracking
while read wt_path; do
  branch=$(git -C "$wt_path" branch --show-current 2>/dev/null)
  if [ -n "$branch" ] && [ ! -f "{{tracking_folder}}/$branch/context.yaml" ]; then
    echo "ORPHAN_WORKTREE: $wt_path ($branch)"
  fi
done < /tmp/worktrees.txt
```

**Find orphan tracking (branch no longer exists):**
```bash
# List all tracking entries
find "{{tracking_folder}}" -name "context.yaml" -type f | while read ctx; do
  branch=$(dirname "$ctx" | sed "s|{{tracking_folder}}/||")
  # Check if branch exists
  if ! git show-ref --verify --quiet "refs/heads/$branch"; then
    echo "ORPHAN_TRACKING: $branch (branch deleted)"
  fi
done
```

### Step 2: Display findings

```
Nook Prune Scan

Base workspace: {{base_path}}
Tracking:       {{tracking_folder}}

Broken Worktrees (git references to missing paths):
{{if broken_worktrees}}
{{for each broken in broken_worktrees}}
  [X] {{broken.path}} -> MISSING
{{end for}}
  -> Will be cleaned with `git worktree prune`
{{else}}
  (none found)
{{end if}}

Orphan Worktrees (no tracking in base workspace):
{{if orphan_worktrees}}
{{for each orphan in orphan_worktrees}}
  [?] {{orphan.path}} ({{orphan.branch}})
      No context.yaml at: {{tracking_folder}}/{{orphan.branch}}/
{{end for}}
{{else}}
  (none found)
{{end if}}

Orphan Tracking (branch deleted):
{{if orphan_tracking}}
{{for each orphan in orphan_tracking}}
  [X] {{tracking_folder}}/{{orphan.branch}}/
      Branch no longer exists - safe to delete tracking
{{end for}}
{{else}}
  (none found)
{{end if}}
```

### Step 3: Prune broken worktrees (automatic)

```bash
git worktree prune
```

This is safe and always runs - just cleans git's internal references.

### Step 4: Handle orphan worktrees (interactive)

For each orphan worktree, first check if it's currently open:

**Check if worktree is open (cannot delete):**
```bash
current_dir=$(pwd)
current_branch=$(git branch --show-current)

# Check 1: Current branch matches orphan branch
if [ "$current_branch" = "{{orphan.branch}}" ]; then
  echo "[SKIP] {{orphan.path}} ({{orphan.branch}}) - CURRENTLY CHECKED OUT"
  continue  # Skip to next orphan
fi

# Check 2: Current directory is inside orphan worktree
case "$current_dir" in
  "{{orphan.path}}"*)
    echo "[SKIP] {{orphan.path}} ({{orphan.branch}}) - YOU ARE INSIDE THIS WORKTREE"
    continue  # Skip to next orphan
    ;;
esac
```

**If not open, prompt for action:**
```
Orphan worktree: {{path}} ({{branch}})

This worktree has no tracking data (never synced to base workspace).

What would you like to do?
[D] Delete     - Remove the worktree entirely
[S] Skip       - Leave it alone for now

Choice [D/S]:
```

- **D**: Run `git worktree remove "{{path}}"`
- **S**: Continue to next

**If worktree is open, display skip message:**
```
[SKIP] {{orphan.path}} ({{orphan.branch}})
       Cannot delete - this worktree is currently open
       {{if current_branch == orphan.branch}}
       Reason: You are checked out on this branch
       {{else}}
       Reason: Your working directory is inside this worktree
       {{end if}}
       Navigate elsewhere first, then re-run prune.
```

### Step 5: Handle orphan tracking (batch)

```
Found {{count}} orphan tracking entries (branches deleted).

These context files have no matching branch and can be safely deleted:
{{for each orphan}}
  - {{tracking_folder}}/{{orphan.branch}}/
{{end for}}

Delete all orphan tracking? [y/N]:
```

If confirmed:
```bash
{{for each orphan}}
rm -rf "{{tracking_folder}}/{{orphan.branch}}"
{{end for}}

# Remind to commit
echo ""
echo "Orphan tracking deleted. Commit changes in base workspace:"
echo "  cd {{base_path}}"
echo "  git add .bmad-tracking/"
echo "  git commit -m 'chore(th): prune orphan tracking'"
```

### Step 6: Report summary

```
Prune Complete

- Broken worktree refs cleaned: {{broken_count}}
- Orphan worktrees deleted: {{deleted_wt_count}}
- Orphan tracking deleted: {{orphan_tracking_deleted}}

{{if orphan_tracking_deleted > 0}}
IMPORTANT: Commit tracking changes in base workspace:
  cd {{base_path}}
  git add .bmad-tracking/
  git commit -m "chore(th): prune orphan tracking"
{{end if}}
```

## Notes

- **Broken worktrees** are always pruned (safe git cleanup)
- **Orphan worktrees** require user decision (may have uncommitted work)
- **Open worktrees cannot be deleted** - auto-skipped if:
  - You are checked out on that branch
  - Your working directory is inside that worktree
- **Orphan tracking** is batch-deleted with confirmation
- **Tracking changes** must be committed separately in base workspace
