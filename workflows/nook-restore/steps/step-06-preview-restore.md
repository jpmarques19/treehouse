---
name: 'step-06-preview-restore'
description: 'Show restore preview and warn about overwrites'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-restore'

# File References
thisStepFile: '{workflow_path}/steps/step-06-preview-restore.md'
nextStepFile: '{workflow_path}/steps/step-07-confirm.md'
---

# Step 6: Preview Restore Operation

## STEP GOAL:

Show the user exactly what will happen during restore, with clear warnings about overwrites.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Display clear, formatted preview
- Highlight any destructive operations prominently

### Role Reinforcement:

- You are a workspace automation assistant
- User needs full visibility before confirming
- Make overwrite warnings unmissable

## EXECUTION SEQUENCE:

### 1. Display Restore Summary

<action>Display restore preview header:</action>

```
============================================================
                    RESTORE PREVIEW
============================================================

Source:  {{tracking_path}}
Target:  {{current_path}}
Branch:  {{restore_branch}} -> {{current_branch}}

{{if synced_at available}}
Context synced at: {{synced_at}}
{{end if}}

------------------------------------------------------------
```

### 2. Show Restore Plan

<action>Display each path with its action:</action>

```
Paths to restore:

{{for each entry in restore_plan}}
{{if action is RESTORE and local is MISSING}}
  [RESTORE]   {{path}}
              Source: {{tracking_path}}/{{path}} ({{tracking_count}} files)
              Target: {{path}} (will be created)

{{else if action is OVERWRITE}}
  [OVERWRITE] {{path}}
              Source: {{tracking_path}}/{{path}} ({{tracking_count}} files)
              Target: {{path}} ({{local_count}} files will be REPLACED)

{{else if action is SKIP}}
  [SKIP]      {{path}}
              Not in saved context - will not be touched

{{end if}}
{{end for}}
```

### 3. Show Overwrite Warning (if applicable)

<check if="any paths have action OVERWRITE">
  <action>Display prominent warning:</action>

```
============================================================
                    ⚠️  WARNING  ⚠️
============================================================

The following paths have EXISTING CONTENT that will be REPLACED:

{{for each entry in restore_plan where action is OVERWRITE}}
  {{path}}
    Local:    {{local_count}} files
    Tracking: {{tracking_count}} files
    Result:   Local files DELETED, replaced with tracking version

{{end for}}

This operation is DESTRUCTIVE. Local changes will be LOST.
============================================================
```
</check>

### 4. Show Context Metadata

<check if="context metadata available">
  <action>Display context info:</action>

```
Context Information:
  Synced at:    {{synced_at}}
  {{if synced_by}}Synced by:    {{synced_by}}{{end if}}
  {{if forked_from}}Forked from:  {{forked_from}}{{end if}}
```
</check>

### 5. Explain Artifact Visibility

<action>Display artifact visibility note:</action>

```
------------------------------------------------------------
NOTE: Restored files won't show in git status

  docs/                    - gitignored (doesn't travel with git)
  .bmad/_cfg/agents/       - skip-worktree (tracked but local ignored)
  .bmad/bmm/config.yaml    - skip-worktree (tracked but local ignored)

Each worktree manages its own copy of these artifacts.
------------------------------------------------------------
```

### 6. Auto-Proceed to Confirmation

<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After preview:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Clear preview displayed with all paths
- Overwrite warnings prominently shown
- Context metadata displayed if available
- Artifact visibility explained
- Proceeded to step 7

### SYSTEM FAILURE:
- Not showing all affected paths
- Hiding or minimizing overwrite warnings
- Unclear about what will happen
