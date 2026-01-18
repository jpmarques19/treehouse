# Pin Workflow

Save learnings to the nook's board.

## Overview

This workflow helps agents save key learnings and context to the nook's board so they persist across sessions. Pins are lightweight notes that capture decisions, discoveries, and context worth remembering.

## Execution

### Step 1: Detect Current Nook

Check if currently in a nook worktree:

1. Run `th list` to get current context
2. Check `current_nook` field in response
3. If `current_nook` is null, show: `⚠ Not in a nook. Nothing to pin.` and exit

### Step 2: Generate Pin Content

Prompt the agent to summarize what should be saved:

```
What should be pinned to this nook's board?

Consider:
- Key decisions made
- Important context discovered
- Learnings for future sessions
```

Wait for agent to generate content.

### Step 3: Save Pin

Run the pin command with the generated content:

```bash
th pin "<content>"
```

### Step 4: Confirm

After saving, display confirmation:

```
✓ Pinned to {nook-id}'s board.
```

## Dependencies

- Treehouse must be initialized
- Must be in a nook worktree (not base repo)
