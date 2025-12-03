---
name: 'step-09-report'
description: 'Report restore completion and next steps'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'

# File References
thisStepFile: '{workflow_path}/steps/step-09-report.md'
# No nextStepFile - this is the final step
---

# Step 9: Report Completion

## STEP GOAL:

Report the restore operation results and guide user on next steps.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Provide clear summary of what was done
- Guide user on what to do next

### Role Reinforcement:

- You are a workspace automation assistant
- Provide helpful, actionable completion report
- End on a positive, forward-looking note

## EXECUTION SEQUENCE:

### 1. Display Completion Header

<action>Display completion summary:</action>

```
============================================================
                   RESTORE COMPLETE
============================================================

Source:  {{tracking_path}}
Target:  {{current_path}}
Branch:  {{restore_branch}} -> {{current_branch}}

------------------------------------------------------------
```

### 2. Show Restore Results

<action>Display results for each path:</action>

```
Restored Paths:

{{for each result in restore_results}}
{{if status is SUCCESS}}
  [DONE] {{path}}
         {{if type is directory}}{{file_count}} files restored{{end if}}
         {{if type is file}}File restored{{end if}}

{{else if status is FAILED}}
  [FAIL] {{path}}
         Error: {{error_message}}

{{else if status is SKIPPED}}
  [SKIP] {{path}}
         Not in saved context

{{end if}}
{{end for}}
```

### 3. Show Statistics

<action>Calculate and display statistics:</action>

```
------------------------------------------------------------
Summary:
  Total paths processed: {{total_count}}
  Successfully restored: {{success_count}}
  Failed:               {{failed_count}}
  Skipped:              {{skipped_count}}

{{if synced_at available}}
Context from {{synced_at}} has been restored.
{{end if}}
------------------------------------------------------------
```

### 4. Explain Artifact Visibility

<action>Display visibility reminder:</action>

```
REMEMBER: Restored files won't show in git status

  docs/                    - gitignored (doesn't travel with git)
  .bmad/_cfg/agents/       - skip-worktree (tracked but local ignored)
  .bmad/bmm/config.yaml    - skip-worktree (tracked but local ignored)

Each worktree manages its own copy of these artifacts.
```

### 5. Provide Next Steps

<action>Display next steps based on context:</action>

```
------------------------------------------------------------
                      NEXT STEPS
------------------------------------------------------------

Your workspace now has the restored context. You can:

  1. Continue working with the restored content
     - docs/ contains your project documentation
     - Configs are set to the tracked versions

  2. Save changes back after making updates
     /bmad:th:workflows:nook-sync

  3. Check what was restored
     ls -la docs/
     cat .bmad/th/config.yaml

{{if current_branch != restore_branch}}
  NOTE: You restored from '{{restore_branch}}' to '{{current_branch}}'.
  When you sync, changes will be saved under '{{current_branch}}'.
{{end if}}

============================================================
```

### 6. Workflow Complete

<action>Display: "Restore workflow complete."</action>
<action>HALT - workflow finished</action>

## MENU OPTIONS:

This is the final step. After report:
- Workflow complete - no further steps

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Clear completion summary displayed
- All results reported with status
- Statistics calculated correctly
- Next steps provided
- Workflow completed cleanly

### SYSTEM FAILURE:
- Incomplete or unclear report
- Missing result information
- No guidance on next steps
- Workflow not properly terminated
