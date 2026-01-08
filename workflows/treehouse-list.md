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

1. **Group nooks by parent** - Create a map of parent -> children
2. **Start from base branch** (main/master/dev)
3. **Recursively render children** with proper indentation

Tree characters:
- `├──` for middle items (has siblings after)
- `└──` for last items
- `│   ` for continuing vertical lines
- `    ` for empty space (after last item)

Status icons:
- `●` Active (current nook or has worktree)
- `○` Inactive (tracked but no worktree)
- `✗` Orphan (worktree exists but not tracked)

Mark current nook (from `current_nook` field) as active.

Example output:
```
main
├── ● a1b2-auth-spike
│   └── ○ c3d4-jwt-variant
├── ○ e5f6-redis-cache
└── ○ g7h8-refactor-api
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
 [s] sleep        [r] refresh     [q] quit
```

### Step 6: Handle Menu Selection

Wait for user input and handle:

- **[n] Navigate**: Ask for nook ID, then `cd` to its worktree path
- **[d] Delete**: Ask for nook ID, confirm, then run `th remove {nook-id}`
- **[p] Prune**: Run `th prune` to clean up orphan worktrees
- **[s] Sleep**: Ask for nook ID, remove worktree but keep branch tracking
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
