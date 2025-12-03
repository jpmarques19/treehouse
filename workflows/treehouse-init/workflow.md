---
name: "treehouse-init"
description: "Initialize base workspace - sets up gitignore for sync artifacts, creates tracking folder, establishes source of truth for all nooks"
version: "3.0.0"
web_bundle: true
---

# Treehouse Init - Initialize Base Workspace

**Goal:** Initialize this worktree as the base workspace - the source of truth for all workspace tracking.

**Your Role:** You are a workspace automation assistant helping establish the base workspace. This workflow sets up the foundation for the treehouse architecture, enabling isolated development nooks with shared context tracking.

---

## WORKFLOW ARCHITECTURE

This uses **step-file architecture** for disciplined execution:

### Core Principles

- **Micro-file Design**: Each step is a self contained instruction file that must be followed exactly
- **Just-In-Time Loading**: Only the current step file is in memory - never load future step files until told to do so
- **Sequential Enforcement**: Sequence within the step files must be completed in order, no skipping or optimization allowed
- **Action Workflow**: This workflow performs git operations, no document output - state tracked via config and submodule

### Step Processing Rules

1. **READ COMPLETELY**: Always read the entire step file before taking any action
2. **FOLLOW SEQUENCE**: Execute all numbered sections in order, never deviate
3. **WAIT FOR INPUT**: If a menu is presented, halt and wait for user selection
4. **EXECUTE COMMANDS**: Run bash commands exactly as specified
5. **LOAD NEXT**: When directed, load, read entire file, then execute the next step file

### Critical Rules (NO EXCEPTIONS)

- **NEVER** load multiple step files simultaneously
- **ALWAYS** read entire step file before execution
- **NEVER** skip steps or optimize the sequence
- **ALWAYS** follow the exact instructions in the step file
- **ALWAYS** halt at menus and wait for user input
- **NEVER** create mental todo lists from future steps

---

## KEY ARCHITECTURE POINTS

This workflow establishes the three-tier artifact tracking system:

| Type | Paths | Behavior in Nooks |
|------|-------|-------------------|
| **Gitignored** | `docs/` | Don't exist - must restore |
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git, local changes ignored |
| **Submodule** | `.bmad-tracking/`, `.gitmodules` | Skip-worktree applied, then removed |

### Why This Architecture

- **Gitignored paths**: Large directories that would bloat git history. Simpler mental model.
- **Skip-worktree paths**: Small config files that should exist immediately in nooks.
- **Submodule**: Tracking folder isolated from nooks. No sparse-checkout complexity.

---

## INITIALIZATION SEQUENCE

### 1. Configuration Loading

Load and read full config from `{project-root}/.bmad/th/config.yaml` and resolve:

- `gitignore_paths`, `skip_worktree_paths`, `sync_paths`, `tracking_folder`

### 2. First Step EXECUTION

Load, read the full file and then execute `{project-root}/.bmad/th/workflows/treehouse-init/steps/step-01-check-git-clean.md` to begin the workflow.
