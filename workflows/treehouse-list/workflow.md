---
name: "treehouse-list"
description: "View workspace lineage tree - displays all tracked workspaces from base, with status and cleanup options"
version: "2.0.0"
web_bundle: true
---

# Treehouse List - View Workspace Lineage Tree

**Goal:** Display a visual tree of all workspaces showing nook lineage, hierarchy, health status, and provide quick access to navigation and cleanup operations.

**Your Role:** In addition to your name, communication_style, and persona, you are also a workspace navigation assistant helping a developer view and manage their workspace lineage. This is a utility workflow - display information clearly and handle menu actions efficiently.

---

## WORKFLOW ARCHITECTURE

This uses **step-file architecture** for disciplined execution:

### Core Principles

- **Micro-file Design**: Each step is a self contained instruction file that must be followed exactly
- **Just-In-Time Loading**: Only the current step file is in memory - never load future step files until told to do so
- **Sequential Enforcement**: Sequence within the step files must be completed in order, no skipping or optimization allowed
- **Action Workflow**: This workflow displays status and performs operations, no document output

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

## INITIALIZATION SEQUENCE

### 1. Configuration Loading

Load and read full config from `{project-root}/.bmad/th/config.yaml` and resolve:

- `base_workspace_path`, `base_branch`, `worktrees_folder`

### 2. First Step EXECUTION

Load, read the full file and then execute `{project-root}/.bmad/th/workflows/treehouse-list/steps/step-01-check-init.md` to begin the workflow.
