# Hat Add Workflow

Create a new treehouse hat.

## Overview

This workflow guides the creation of a new treehouse hat - a simple .md file containing domain knowledge that applies across all nooks.

## Execution

### Step 1: Check Treehouse

1. Run `th list` to verify treehouse is initialized
2. If error with `INIT_NOT_FOUND`, show: `Treehouse not initialized. Run /th:workflows:treehouse-init first` and exit

### Step 2: Collect Hat Info

Prompt for hat details:

```
Create new hat

1. Name (lowercase, no spaces):
2. Title (e.g., "Code Reviewer"):
3. Icon (single emoji, default diamond):
```

Validate:
- Name must be lowercase, alphanumeric, no spaces (convert spaces to hyphens)
- Name must not already exist in `.treehouse/hats/`
- Icon defaults to diamond if not provided

### Step 3: Create Hat

Run the CLI command:

```bash
th hat add {name} '{"title":"{title}","icon":"{icon}"}'
```

### Step 4: Confirm Creation

Display success:

```
Created hat: {name}

  {icon} {Name} - {title}

  File: .treehouse/hats/{name}.md

Activate with: /th:hats:{name}
```

## Dependencies

- Treehouse must be initialized (`th init`)
- Write access to `.treehouse/hats/`
