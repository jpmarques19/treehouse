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
- **[d] Delete**: Ask for nook ID, confirm, then run `th remove {nook-id}`
- **[p] Prune**: Run `th prune` to clean up orphan worktrees
- **[r] Refresh**: Re-run `th list` and re-render tree
- **[q] Quit**: Exit workflow

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
