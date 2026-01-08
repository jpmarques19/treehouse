# Handoff Workflow

Save current context before switching nooks or ending a session.

## Overview

This workflow helps agents save their current work context to memory files so they can resume later. It prompts the agent to summarize accomplishments, current state, and next actions, then persists this to the appropriate memory files.

## Execution

### Step 1: Detect Current Nook

Check if currently in a nook worktree:

1. Run `th list` to get current context
2. Check `current_nook` field in response
3. If `current_nook` is null, show: `⚠ Not in a nook. Nothing to hand off.` and exit

### Step 2: Generate Session Summary

Prompt the agent to summarize the current session:

```
Handoff: Please summarize this session for future context.

What was accomplished?
Current state (build status, blockers)?
What are the next actions?
```

Wait for agent to generate summary.

### Step 3: Update Memory Files

**Short-term Memory** (`memories/{nook-id}.md`):
- Add date-stamped entry: `## YYYY-MM-DD: {Activity Title}`
- Include key decisions made
- Include important context discovered
- Keep entries concise but complete

**Session State** (`sessions/{nook-id}.md`):
- Update `## Current State` section
- Update `### Session Focus` with current work and status
- Update `### Recent Commits` from git log
- Update `### Build Status` (pass/fail)
- Update `## Next Actions` with prioritized tasks
- Update `## Resume Context` with summary for next session

### Step 4: Confirm Save

After saving, display confirmation:

```
✓ Context saved for {nook-id}

Saved:
  - Session state (what you were doing)
  - Short-term memory (decisions & context)

Safe to switch nooks or close session.
```

### Step 5: Prompt for Long-term Knowledge (Optional)

If the agent identified learnings that apply globally:

```
Would you like to add any of this to long-term knowledge?
(This persists across all nooks)

[y/n]
```

If confirmed, update `knowledge.md` with the global learnings.

## Memory File Locations

```
.treehouse/agents/{agent}/
├── knowledge.md              # Long-term (global)
├── memories/{nook-id}.md     # Short-term (per-nook)
└── sessions/{nook-id}.md     # Session state (per-nook)
```

## Dependencies

- Agent must be loaded with nook context
- Treehouse must be initialized
- Must be in a nook worktree (not base repo)
