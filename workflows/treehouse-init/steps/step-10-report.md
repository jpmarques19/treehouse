---
name: 'step-10-report'
description: 'Report completion with next steps'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/treehouse-init'

# File References
thisStepFile: '{workflow_path}/steps/step-10-report.md'
# No nextStepFile - this is the final step
---

# Step 10: Report Completion

## STEP GOAL:

Report the initialization results and guide user on next steps.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Provide clear summary of what was set up
- Include actionable next steps

### Role Reinforcement:

- You are a workspace automation assistant
- Help user understand the new architecture
- End on a positive, forward-looking note

## EXECUTION SEQUENCE:

### 1. Display Completion Summary

<action>Display:</action>

```
============================================================
              BASE WORKSPACE INITIALIZED
============================================================

Base Workspace
   Path:   {{workspace_path}}
   Branch: {{current_branch}}

------------------------------------------------------------
Artifact Tracking Configuration:
------------------------------------------------------------

Gitignored (per-nook, synced via th workflows):
{{for each path in gitignore_paths}}
   - {{path}}/
{{end for}}

Skip-worktree in nooks (tracked, local changes ignored):
{{for each path in skip_worktree_paths}}
   - {{path}}
{{end for}}

Tracking: .bmad-tracking/ (git submodule)
   - Submodule initialized here with content
   - .gitmodules tracked in parent repo
   - Nooks: submodule not initialized (empty placeholder removed)

------------------------------------------------------------
Next Steps:
------------------------------------------------------------

1. Save your current artifacts:
   /bmad:th:workflows:nook-sync

2. Create a focused development nook:
   /bmad:th:workflows:nook-fork

3. View all workspaces:
   /bmad:th:workflows:treehouse-list

------------------------------------------------------------
Architecture Overview:
------------------------------------------------------------

   - This worktree is the SOURCE OF TRUTH
   - .bmad-tracking/ is a git submodule (initialized here only)
   - Nooks get skip_worktree_paths from git (files exist!)
   - Nooks don't get gitignore_paths (must restore)
   - Nooks don't get .bmad-tracking/ or .gitmodules (skip-worktree + removed)
   - nook-sync saves TO here (optionally auto-commits submodule)
   - nook-restore loads FROM here

============================================================
```

### 2. Workflow Complete

<action>HALT - workflow finished</action>

## MENU OPTIONS:

This is the final step. After report:
- Workflow complete - no further steps

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Clear completion summary displayed
- All configuration shown
- Next steps provided
- Architecture explained
- Workflow completed cleanly

### SYSTEM FAILURE:
- Incomplete summary
- Missing next steps
- Unclear architecture explanation
- Workflow not properly terminated
