---
name: 'step-14-report'
description: 'Report results and remind about commit'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-14-report.md'
# No nextStepFile - this is the final step
---

# Step 14: Report Completion

## STEP GOAL:

Report the sync operation results, show summary, and provide commit instructions if needed.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Provide clear summary of what was done
- Include commit instructions if manual commit needed

### Role Reinforcement:

- You are a workspace automation assistant
- Provide helpful, actionable completion report
- End on a positive, forward-looking note

## EXECUTION SEQUENCE:

### 1. Display Completion Header

<action>Display completion summary:</action>

```
============================================================
                    SYNC COMPLETE
============================================================

Source: {{current_path}}
Branch: {{current_branch}} ({{current_hash}})
Target: {{tracking_path}}

------------------------------------------------------------
Synced Paths:
------------------------------------------------------------

{{for each result in sync_results}}
{{if status is SUCCESS}}
  [DONE] {{path}}
         {{if type is directory}}{{file_count}} files synced{{end if}}
         {{if type is file}}File synced{{end if}}

{{else if status is FAILED}}
  [FAIL] {{path}}
         Error: {{error_message}}

{{else if status is SKIPPED}}
  [SKIP] {{path}}
         Not found in source

{{end if}}
{{end for}}
```

### 2. Show Lineage

<action>Display lineage chain:</action>

```
------------------------------------------------------------
Lineage ({{full_lineage | length}} levels):
------------------------------------------------------------

{{for each entry in full_lineage}}
  {{index}}. {{entry.branch}} ({{entry.hash}})
{{end for}}
```

### 3. Show Summary Preview

<action>Display summary preview:</action>

```
------------------------------------------------------------
Summary Generated:
------------------------------------------------------------

┌─────────────────────────────────────────────────────────
│ Purpose: {{nook_summary.purpose | truncate(80)}}
│ Status:  {{nook_summary.status | truncate(80)}}
│ Insights: {{nook_summary.insights | length}} | Decisions: {{nook_summary.decisions | length}}
│ Blockers: {{nook_summary.blockers | length}} | Next Steps: {{nook_summary.next_steps | length}}
└─────────────────────────────────────────────────────────

Manifest: {{tracking_path}}/context.yaml
```

### 4. Show Commit Status

<check if="auto_commit was enabled and succeeded">
  <action>Display:</action>
  ```
  ------------------------------------------------------------
  Auto-commit: ENABLED
  ------------------------------------------------------------

    - Submodule changes committed
    - Parent pointer updated

  Context is ready to use in other nooks!

  To restore this context in another worktree:
    /bmad:th:workflows:nook-restore
  ```
</check>

<check if="manual_commit_needed == true">
  <action>Display:</action>
  ```
  ------------------------------------------------------------
  IMPORTANT: Manual commit required
  ------------------------------------------------------------

  The tracking folder is a git submodule. To commit changes:

  1. Go to base workspace:
     cd {{base_path}}

  2. Commit inside submodule:
     cd .bmad-tracking
     git add .
     git commit -m "sync: {{current_branch}} context"
     cd ..

  3. Update parent pointer:
     git add .bmad-tracking
     git commit -m "chore(th): update tracking pointer"

  To enable auto-commit, set in .bmad/th/config.yaml:
    auto_commit_tracking: true

  To restore this context in another worktree:
    /bmad:th:workflows:nook-restore
  ```
</check>

### 5. Workflow Complete

<action>Display:</action>
```
============================================================
```

<action>HALT - workflow finished</action>

## MENU OPTIONS:

This is the final step. After report:
- Workflow complete - no further steps

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Clear completion summary displayed
- All results reported with status
- Lineage chain shown
- Summary preview displayed
- Commit instructions provided if needed
- Workflow completed cleanly

### SYSTEM FAILURE:
- Incomplete or unclear report
- Missing result information
- No commit instructions when needed
- Workflow not properly terminated
