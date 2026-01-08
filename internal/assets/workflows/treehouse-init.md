# Treehouse Init Workflow

Initialize a treehouse workspace in the current git repository.

## Execution

### Step 1: Check Prerequisites

Verify we're in a git repository by running:
```bash
git rev-parse --git-dir
```

If this fails, show: "Not a git repository. Initialize git first with: git init"

### Step 2: Run Init Command

```bash
th init
```

### Step 3: Display Result

Parse JSON response and display appropriate feedback.

**On success** (JSON has `"success": true`):
```
Treehouse initialized

Created:
  .treehouse/
  ├── decks.yaml
  ├── nooks/
  ├── crew/oak/
  └── workflows/

You're ready to fork your first nook with /nook-fork
```

**On error** (JSON has `"success": false`):

| Error Code | Message |
|------------|---------|
| `INIT_ALREADY_EXISTS` | Treehouse already initialized in this repository |
| `INIT_NOT_GIT_REPO` | Not a git repository. Initialize git first with: git init |
| `GIT_VERSION_UNSUPPORTED` | Git 2.5+ required for worktree support. Update git first. |
| Other | Failed to initialize: {message} |

## JSON Response Structure

```json
{
  "success": true,
  "data": {
    "path": "/path/to/.treehouse",
    "created": ["decks.yaml", "crew/oak/", "workflows/"]
  }
}
```

## Notes

- Requires `th` CLI to be installed and in PATH
- If `th` command is not found, show: "th CLI not installed. Install it first."
- Creates Oak agent as default crew member
- All file system operations are handled by the CLI, not this workflow
