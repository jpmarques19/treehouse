---
name: "nook-restore"
description: "Restore saved context FROM base workspace tracking folder to current worktree"
version: "3.0.0"
web_bundle: true
---

# Nook Restore - Load Context from Base Workspace

**Goal:** Restore saved context (docs, configs, etc.) FROM the base workspace's tracking folder to current worktree.

**Your Role:** You are a workspace automation assistant helping restore development context. This workflow copies saved artifacts from the centralized tracking folder to the current worktree, enabling context resumption after nook creation or workspace switching.

---

## WORKFLOW ARCHITECTURE

This uses **step-file architecture** for disciplined execution:

### Core Principles

- **Micro-file Design**: Each step is a self contained instruction file that must be followed exactly
- **Just-In-Time Loading**: Only the current step file is in memory - never load future step files until told to do so
- **Sequential Enforcement**: Sequence within the step files must be completed in order, no skipping or optimization allowed
- **Action Workflow**: This workflow performs file operations, no document output - state tracked via restored files

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

- All context is loaded FROM `base_workspace_path/.bmad-tracking/{branch}/`
- Can restore context from ANY branch (not just current)
- Common use: nook starts clean, then restore parent's context
- Overwrites local content with saved version

### Two-Tier Artifact System

| Type | Paths | Nook State | Restore Behavior |
|------|-------|------------|------------------|
| **Gitignored** | `docs/` | Don't exist | Restored from tracking |
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git | Overwritten with tracked version |

---

## INITIALIZATION SEQUENCE

### 1. Configuration Loading

Load and read full config from `{project-root}/.bmad/th/config.yaml` and resolve:

- `base_workspace_path`, `sync_paths`

### 2. First Step EXECUTION

Load, read the full file and then execute `{project-root}/.bmad/th/workflows/nook-restore/steps/step-01-check-init.md` to begin the workflow.
