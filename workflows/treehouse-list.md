# Treehouse List Workflow

Display workspace lineage as a visual tree with status icons and action menu.

## Overview

This workflow calls `th list` to get JSON data about all decks and nooks, then renders a visual tree representation with status indicators and an interactive action menu.

## Execution

### Step 1: Get Workspace Data

Run the list command:
```bash
th list
```

Parse the JSON response. If `success` is false, handle errors:

- **INIT_NOT_FOUND**: Show "Treehouse not initialized. Run /treehouse-init first"
- **Other errors**: Show the error message

### Step 2: Render Tree Header

```
Workspace Lineage

base {base.path}
here {base.branch} · {base.commit (first 7 chars)}
```

### Step 3: Build Tree Visualization

If `decks` array is empty:
```
No nooks yet. Create your first with /nook-fork
```

Otherwise, build the tree using this algorithm:

1. **Collect all nooks from all decks** - Create a unified map of all nooks
2. **Group nooks by parent** - Create a map of parent -> children across all decks
3. **Always start from base branch** (main/master/dev from `base.branch`) - This is the trunk/root
4. **Recursively render children** with proper indentation, following all parent-child chains
5. **Mark current nook** - Place `●` icon on the nook matching `current_nook` field, regardless of tree depth

Tree characters:
- `├──` for middle items (has siblings after)
- `└──` for last items
- `│   ` for continuing vertical lines
- `    ` for empty space (after last item)

Status icons:
- `●` Current nook (only one, marked from `current_nook` field)
- `○` Inactive nook (has worktree but not current)
- `✗` Orphan (tracked in decks.yaml but no worktree)

**Important:** The tree ALWAYS starts from the base branch trunk, showing the full dependency graph from root to all leaves, regardless of which nook is current.

Example output (with current nook at depth 2):
```
main
├── ○ 8aa9-test-nook
│   └── ○ 8aa9-sub-test-nook
└── ● 8aa9-treehouse-planning (current)
    └── ○ 6c81-auth-spike2
```

### Step 4: Render Summary Bar

Count nooks by status and display:
```
─────────────────────────────────────────────────
 ● N active       ○ N inactive    ✗ N orphan
─────────────────────────────────────────────────
```

### Step 5: Show Action Menu

```
 [n] navigate     [d] delete      [p] prune
 [r] refresh      [q] quit
```

### Step 6: Handle Menu Selection

Wait for user input and handle:

- **[n] Navigate**: Ask for nook ID, then `cd` to its worktree path
- **[d] Delete**: Show delete submenu (see Step 6a)
- **[p] Prune**: Show prune confirmation (see Step 6b)
- **[r] Refresh**: Re-run `th list` and re-render tree
- **[q] Quit**: Exit workflow

#### Step 6a: Delete Submenu

When `[d]` is selected, show numbered list of deletable nooks (excludes current nook):

```
Delete nook:

1. ○ c3d4-jwt-variant
2. ○ e5f6-redis-cache
3. ✗ x1y2-old-experiment

enter number or q to cancel:
```

When user enters a number, show confirmation with consequences:

```
Delete {nook-id}?

this will:
  - remove git worktree
  - delete crew memory files
  - update decks.yaml

[y/n]
```

If confirmed with `y`:
1. Run `th remove {nook-id}`
2. Show: `✓ Removed {nook-id}`
3. Refresh tree and return to main menu

If declined with `n`:
1. Show: `Cancelled`
2. Return to main menu

#### Step 6b: Prune Confirmation

When `[p]` is selected:

First, find orphans by checking nooks with status `"orphan"` from the list data.

If no orphans exist:
```
No orphaned nooks to prune.
```
Return to main menu.

If orphans exist, show confirmation:
```
Prune orphaned nooks?

found {N} orphans:
  - x1y2-old-experiment
  - z3w4-abandoned

this will:
  - remove entries from decks.yaml
  - delete associated crew memory files

[y/n]
```

If confirmed with `y`:
1. Run `th prune`
2. Show: `✓ Pruned {N} orphaned nooks`
3. Refresh tree and return to main menu

If declined with `n`:
1. Show: `Cancelled`
2. Return to main menu

## JSON Response Structure

```json
{
  "success": true,
  "data": {
    "base": {
      "path": "/path/to/repo",
      "branch": "main",
      "commit": "abc1234def5678"
    },
    "current_nook": "a1b2-auth-spike",
    "decks": [
      {
        "id": "dk-a1b2",
        "created": "2026-01-08",
        "nooks": [
          {
            "id": "a1b2-auth-spike",
            "parent": "main",
            "created": "2026-01-08",
            "worktree": "/path/to/repo/.treehouse/nooks/a1b2-auth-spike",
            "status": "active"
          }
        ]
      }
    ]
  }
}
```

## Dependencies

- `th` CLI must be installed and in PATH
- Treehouse must be initialized (`th init`)
