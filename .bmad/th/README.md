# Treehouse (BMad Workspace)

> Fork easily. Stay focused. Sync what matters.

---

**EXPERIMENTAL - USE AT YOUR OWN RISK**

This module is a **work in progress** with very narrow testing:
- Tested exclusively with **Claude Code + Opus 4.5**
- **Tightly coupled with git internals** (worktrees, skip-worktree, sparse-checkout)
- May leave orphaned worktrees or unexpected git state if workflows are interrupted
- Not battle-tested across different environments or git versions

**If you're not comfortable with low level git operations, proceed with caution.**

---

## The Problem: Context Drift

When projects grow, you need to explore multiple directions at once—investigating issues, trying new approaches, running experiments. Each exploration benefits from isolation, but creates a hidden cost: **your documents, configurations, and decisions start drifting apart.**

A week later, you find yourself with:
- One document describing a certain approach
- Another artifact that evolved differently
- References pointing to a hybrid that exists nowhere
- Configs with conflicting assumptions

Reconciling this drift is painful. The longer parallel explorations run, the harder it becomes to maintain coherence across your workspace.

## The Solution: Focused Nooks with Shared Context

Treehouse creates **separate, focused nook environments** (via git worktrees) while maintaining **conceptual continuity** across all of them through a centralized tracking system.

```
Your Repository
│
├── main-worktree/                    ← Base workspace (source of truth)
│   ├── .bmad-tracking/               ← Centralized context storage
│   │   ├── main/                     ← Context for main branch
│   │   ├── explore/mvp-scope/        ← Context for exploration nook
│   │   └── spike/a1b2-new-approach/  ← Context for spike nook
│   ├── docs/                         ← Current working docs
│   └── content/                      ← Your files
│
├── worktrees/
│   ├── a1b2-new-approach/            ← Focused spike environment
│   │   ├── content/                  ← Can diverge freely
│   │   └── docs/                     ← Restored from tracking
│   │
│   └── 7f3e-research-thread/         ← Focused research environment
│       ├── content/                  ← Dedicated investigation
│       └── docs/                     ← Restored from tracking
```

**Key insight**: Content can diverge freely in each worktree, but documentation and configuration artifacts flow through a single tracking hub. You explicitly choose when to sync changes back, maintaining conscious control over what context persists.

## Core Concepts

### Git Worktrees
Treehouse uses [git worktrees](https://git-scm.com/docs/git-worktree) to create multiple working directories from a single repository. Each worktree has its own branch and working state, but shares the same git history.

### The Tracking Folder
The `.bmad-tracking/` folder in your base workspace is the **single source of truth** for all context artifacts. It stores:
- Documentation snapshots (`docs/`)
- Configuration files (`.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml`)
- Context manifests with lineage information (`context.yaml`)

Each branch gets its own subdirectory, preserving context independently.

### Lineage Tracking
Every nook records its parent, creating a traceable chain:

```
my-project                        ← Base workspace (hash: a1b2)
  └── explore/a1b2-new-direction  ← Nook from a1b2 (hash: 7f3e)
        └── spike/7f3e-deep-dive  ← Nook from 7f3e (hash: c9d1)
```

Branch names encode lineage: `{type}/{parent-hash-4chars}-{name}`

### Two-Tier Artifact System

| Type | Paths | In Nook | Git Behavior |
|------|-------|---------|--------------|
| **Gitignored** | `docs/` | Don't exist until restored | Not tracked by git |
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git | Tracked but local changes hidden |

- **Gitignored paths**: Must be explicitly restored via `nook-restore`
- **Skip-worktree paths**: Exist immediately in nooks, but modifications don't show in `git status`

## Workflows

### 1. Initialize Base Workspace

```
/bmad:th:workflows:treehouse-init
```

Establishes the current worktree as the base workspace:
- Creates `.bmad-tracking/` folder
- Adds `gitignore_paths` (e.g., `docs/`) to `.gitignore`
- Sets `base_workspace_path` in module config

**Run this once** from your main worktree before using other workflows.

### 2. Create a Nook

```
/bmad:th:workflows:nook-fork
```

Creates a focused worktree for dedicated work:
- Prompts for nook type (`explore/`, `spike/`, `bugfix/`, `discovery/`, `feature/`, `experiment/`, `hotfix/`)
- Creates new git worktree with hash-based branch name
- Applies skip-worktree to config files
- Records lineage in base workspace's tracking folder

**Nook starts clean**—run `nook-restore` if you need to replicate docs from any other workspace.

### 3. Sync Context to Base

```
/bmad:th:workflows:nook-sync
```

Saves current context **to** the base workspace's tracking folder:
- Copies configured `sync_paths` to `.bmad-tracking/{branch}/`
- Updates `context.yaml` with lineage and artifact info
- Works from any worktree (base or nook)

**Does not auto-commit**—you must commit `.bmad-tracking/` changes in the base workspace.

### 4. Restore Context from Base

```
/bmad:th:workflows:nook-restore
```

Loads saved context **from** the base workspace's tracking folder:
- Restores docs and configs to current worktree
- Can restore from any branch (current, parent, or selected)
- Overwrites local content with saved version

**Common use**: Nook needs docs from parent branch.

### 5. View Workspace Lineage

```
/bmad:th:workflows:treehouse-list
```

Displays the workspace lineage tree showing all tracked branches and their relationships.

## Typical Workflow

```bash
# 1. Initialize your base workspace (once)
/bmad:th:workflows:treehouse-init

# 2. Work on your workspace, sync periodically
/bmad:th:workflows:nook-sync

# 3. Commit tracking changes (FROM BASE WORKSPACE ONLY!)
cd /path/to/base-workspace
git add .bmad-tracking/ && git commit -m "chore(th): sync context"

# 4. Need to explore something? Create a nook!
/bmad:th:workflows:nook-fork
# Choose: explore, spike, experiment, etc.
# → Creates worktree at configured location

# 5. Switch to the nook
cd /path/to/worktrees/a1b2-my-exploration

# 6. Need the docs? Restore from parent
/bmad:th:workflows:nook-restore

# 7. Work, iterate, document...
# When ready, sync your context
/bmad:th:workflows:nook-sync

# 8. Merge changes back (from base workspace)
cd /path/to/base-workspace
git merge explore/a1b2-my-exploration

# 9. Clean up the worktree
git worktree remove /path/to/worktrees/a1b2-my-exploration
```

## Configuration

Module configuration lives in `.bmad/th/config.yaml`:

```yaml
# Base workspace (set by treehouse-init)
base_workspace_path: "/path/to/your/repo"
base_branch: "main"

# Where nook worktrees are created
worktrees_folder: "/path/to/worktrees"

# Paths synced between worktrees
sync_paths:
  - docs
  - .bmad/_cfg/agents
  - .bmad/bmm/config.yaml

# Paths added to .gitignore (don't travel with git)
gitignore_paths:
  - docs
```

## Known Limitations

1. **Git expertise required**: Understanding worktrees, skip-worktree, and gitignore is essential
2. **Manual cleanup**: Orphaned worktrees must be removed manually (`git worktree remove`)
3. **No auto-commit**: You must commit `.bmad-tracking/` changes yourself
4. **Single base workspace**: All nooks reference one base—no nested base workspaces
5. **Narrowly tested**: Only validated with Claude Code + Opus 4.5 on Linux
6. **Interrupted workflows**: May leave git in unexpected state if aborted mid-workflow

## Why "Treehouse"?

The name reflects the core architecture: your workspace is like a **treehouse** - an elevated, cozy space built on the tree (git). **Nooks** are the individual spaces where focused work happens, each sitting on top of a git branch. The metaphor is intuitive: trees have branches, nooks are cozy isolated spaces, and the whole system provides a safe place to build and experiment.

