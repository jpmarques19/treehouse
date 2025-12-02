---
name: 'step-12-report-completion'
description: 'Display success summary with next steps including agent compile instructions'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-12-report-completion.md'
# No nextStepFile - this is the final step
---

# Step 12: Report Completion

## STEP GOAL:

Display comprehensive success summary with clear next steps, including agent compilation instructions if an agent was created.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 Provide clear, actionable next steps
- ✅ Include compile instructions if agent created

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Make success visible and celebration-worthy
- ✅ Ensure user knows exactly what to do next

## EXECUTION SEQUENCE:

### 1. Set Workflow Outputs

<action>Set final output variables:</action>
- `new_branch` = {{new_branch}}
- `worktree_path` = {{worktree_path}}
- `custom_agent` = {{custom_agent}}

### 2. Display Success Summary

<check if="agent was created (custom_agent != 'none')">
  <action>Display full summary with agent:</action>

```
╔═══════════════════════════════════════════════════════════════════╗
║                    NOOK CREATED SUCCESSFULLY!                     ║
╠═══════════════════════════════════════════════════════════════════╣
║  PARENT WORKSPACE                                                 ║
║  Branch: {{current_branch}}                                       ║
║  Path:   {{base_path}}                                            ║
╠═══════════════════════════════════════════════════════════════════╣
║  NEW NOOK                                                         ║
║  Branch: {{new_branch}}                                           ║
║  Type:   {{nook_type}}                                            ║
║  Path:   {{worktree_path}}                                        ║
╠═══════════════════════════════════════════════════════════════════╣
║  CUSTOM AGENT READY TO COMPILE                                    ║
╠═══════════════════════════════════════════════════════════════════╣
║                                                                   ║
║  Agent source: {{custom_agent}}                                   ║
║  Location: .bmad/custom/src/agents/{{custom_agent}}/              ║
║  Mode: {{agent_mode}}                                             ║
║                                                                   ║
║  To activate your agent, run these commands in the nook:          ║
║                                                                   ║
║  1. Switch to nook:                                               ║
║     cd {{worktree_path}}                                          ║
║                                                                   ║
║  2. Compile the agent:                                            ║
║     npx bmad-method@alpha agent-install                           ║
║                                                                   ║
║  3. Invoke your agent:                                            ║
║     /bmad:custom:agents:{{custom_agent}}                          ║
║                                                                   ║
╠═══════════════════════════════════════════════════════════════════╣
║  LINEAGE                                                          ║
║  {{lineage_display}}                                              ║
╠═══════════════════════════════════════════════════════════════════╣
║  WHAT'S IN THE NOOK                                               ║
║  ✓ Config files (skip-worktree - local changes ignored)          ║
║  ✓ Agent source (needs compilation)                               ║
║  ✗ docs/ folder (run nook-restore if needed)                     ║
╠═══════════════════════════════════════════════════════════════════╣
║  NEXT STEPS                                                       ║
║  1. cd {{worktree_path}}                                          ║
║  2. npx bmad-method@alpha agent-install                           ║
║  3. /bmad:custom:agents:{{custom_agent}}                          ║
║  4. (optional) /th:workflows:nook-restore for docs                ║
║  5. When done: /th:workflows:nook-sync to save context            ║
║  6. Merge back: git checkout {{current_branch}} && git merge      ║
║  7. Cleanup: git worktree remove {{worktree_path}}                ║
╚═══════════════════════════════════════════════════════════════════╝
```
</check>

<check if="no agent created (custom_agent == 'none')">
  <action>Display summary without agent:</action>

```
╔═══════════════════════════════════════════════════════════════════╗
║                    NOOK CREATED SUCCESSFULLY!                     ║
╠═══════════════════════════════════════════════════════════════════╣
║  PARENT WORKSPACE                                                 ║
║  Branch: {{current_branch}}                                       ║
║  Path:   {{base_path}}                                            ║
╠═══════════════════════════════════════════════════════════════════╣
║  NEW NOOK                                                         ║
║  Branch: {{new_branch}}                                           ║
║  Type:   {{nook_type}}                                            ║
║  Path:   {{worktree_path}}                                        ║
╠═══════════════════════════════════════════════════════════════════╣
║  LINEAGE                                                          ║
║  {{lineage_display}}                                              ║
╠═══════════════════════════════════════════════════════════════════╣
║  WHAT'S IN THE NOOK                                               ║
║  ✓ Config files (skip-worktree - local changes ignored)          ║
║  ✗ Custom agent (not requested)                                   ║
║  ✗ docs/ folder (run nook-restore if needed)                     ║
╠═══════════════════════════════════════════════════════════════════╣
║  NEXT STEPS                                                       ║
║  1. cd {{worktree_path}}                                          ║
║  2. (optional) /th:workflows:nook-restore for docs                ║
║  3. When done: /th:workflows:nook-sync to save context            ║
║  4. Merge back: git checkout {{current_branch}} && git merge      ║
║  5. Cleanup: git worktree remove {{worktree_path}}                ║
╚═══════════════════════════════════════════════════════════════════╝
```
</check>

### 3. Workflow Complete

<action>This is the final step - workflow complete</action>
<action>No next step to load</action>

## MENU OPTIONS:

This is the final step. Workflow is complete.

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- All outputs set correctly
- Clear summary displayed
- Agent compile instructions included (if applicable)
- Next steps are actionable
- Lineage chain displayed

### ❌ SYSTEM FAILURE:
- Missing agent compile instructions when agent created
- Unclear or incomplete next steps
- Not displaying lineage
- Missing key information in summary
