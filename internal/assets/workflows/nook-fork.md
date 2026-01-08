# Nook Fork Workflow

Create a new nook (isolated workspace) from the current branch for exploration.

## Execution

### Step 1: Verify Treehouse Initialized

Run `th list` to verify treehouse is initialized.

If error code `INIT_NOT_FOUND`:
- Show: "Treehouse not initialized. Run /treehouse-init first"
- Exit workflow

### Step 2: Detect Current Context

Run `git rev-parse --git-dir` to check if we're in a worktree (nook) or base repo.

If output contains `.git/worktrees/`, we're in a nook - show context:
```
Forking from: {nook-folder-name}
```

Otherwise, we're in base repo - no extra context needed.

### Step 3: Get Nook Name

Ask user: "What would you like to name this nook?"

The name should be descriptive of what they're exploring (e.g., "auth-spike", "jwt-experiment", "fix-login-bug").

### Step 4: Create Nook

Run:
```bash
th fork {user-provided-name}
```

### Step 5: Display Result

Parse JSON response and display appropriate feedback.

**On success** (JSON has `"success": true`):
```
Created nook: {nook_id}

  Deck:   {deck_id}
  Parent: {parent}
  Path:   {worktree}

Switch to new nook?
  cd {worktree}
```

**On error** (JSON has `"success": false`):

| Error Code | Message |
|------------|---------|
| `INIT_NOT_FOUND` | Treehouse not initialized. Run /treehouse-init first |
| `NOOK_ALREADY_EXISTS` | Nook '{id}' already exists. Choose a different name. |
| `NOOK_NAME_REQUIRED` | Nook name required. Please provide a name for your nook. |
| `NOOK_NAME_INVALID` | Invalid nook name. Use letters, numbers, and hyphens only. |
| `GIT_WORKTREE_FAILED` | Failed to create worktree: {message} |
| Other | Failed to create nook: {message} |

## JSON Response Structure

```json
{
  "success": true,
  "data": {
    "nook_id": "a1b2-auth-spike",
    "deck_id": "dk-a1b2",
    "parent": "main",
    "worktree": "/path/to/repo/.treehouse/nooks/a1b2-auth-spike"
  }
}
```

## Notes

- Requires `th` CLI to be installed and in PATH
- If `th` command is not found, show: "th CLI not installed. Install it first."
- Nook ID format: `{4-char-hash}-{sanitized-name}` where hash comes from current commit SHA
- If forking from another nook, new nook is added to same deck as parent
- If forking from base branch (main/dev), a new deck is created
- Worktrees stored in `.treehouse/nooks/{nook-id}/`
