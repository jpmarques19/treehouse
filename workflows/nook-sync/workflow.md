---
name: "nook-sync"
description: "Save working context TO base workspace tracking folder - the single source of truth"
version: "3.0.0"
web_bundle: true
---

# Nook Sync - Save Context to Base Workspace

**Goal:** Save working context (docs, configs, etc.) TO the base workspace's tracking folder, maintaining lineage and generating AI summaries for context analysis.

**Your Role:** You are a workspace automation assistant helping save development context. This workflow copies artifacts from the current worktree to the centralized tracking folder, updates lineage information, and generates semantic summaries for the nook-context-analyst.

---

## WORKFLOW ARCHITECTURE

This uses **step-file architecture** for disciplined execution:

### Core Principles

- **Micro-file Design**: Each step is a self contained instruction file that must be followed exactly
- **Just-In-Time Loading**: Only the current step file is in memory - never load future step files until told to do so
- **Sequential Enforcement**: Sequence within the step files must be completed in order, no skipping or optimization allowed
- **Action Workflow**: This workflow performs file operations, no document output - state tracked via context.yaml

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

- All context is saved TO `base_workspace_path/.bmad-tracking/{branch}/` (git submodule)
- Can be run from ANY worktree (base or nook)
- Sync artifacts from current working directory to base workspace
- Updates or creates context.yaml with lineage information
- Generates AI summary for nook-context-analyst
- Optional auto-commit: commits submodule changes + updates parent pointer

### Submodule Commit Pattern

With the submodule architecture, committing tracking changes requires TWO commits:
1. **Inside .bmad-tracking/ submodule** - commits the actual file changes
2. **In parent repo** - commits the updated submodule pointer

### Nook Summary Generation

Each sync generates an AI-analyzed summary:
- **purpose**: Why this nook exists (derived from branch name + commits)
- **insights**: Key findings/learnings discovered during work
- **decisions**: Important choices made with rationale
- **status**: Current state of work (done/in-progress)
- **blockers**: Known issues or obstacles
- **next_steps**: Recommended actions for context resumption

---

## INITIALIZATION SEQUENCE

### 1. Configuration Loading

Load and read full config from `{project-root}/.bmad/th/config.yaml` and resolve:

- `base_workspace_path`, `sync_paths`, `auto_commit_tracking`

### 2. First Step EXECUTION

Load, read the full file and then execute `{project-root}/.bmad/th/workflows/nook-sync/steps/step-01-check-init.md` to begin the workflow.
