---
name: 'step-10-agent-deploy'
description: 'Deploy agent - move YOLO files from parent or run full wizard'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-10-agent-deploy.md'
nextStepFile: '{workflow_path}/steps/step-11-create-context.md'

# Agent References
createAgentWorkflow: '{project-root}/.bmad/bmb/workflows/create-agent/workflow.md'
---

# Step 10: Agent Deploy

## STEP GOAL:

Complete agent deployment based on agent_mode: move YOLO-created files from parent, run full wizard, or skip.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 For YOLO: wait for subagent if still running, then MOVE files
- 🎯 For Full: invoke create-agent workflow interactively

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Handle each agent_mode correctly
- ✅ Verify agent files after deployment

## EXECUTION SEQUENCE:

### 1. Check Agent Mode

<check if="agent_mode == 'none'">
  <action>Display: "No custom agent requested - skipping agent deploy"</action>
  <action>Set {{custom_agent}} = "none"</action>
  <action>Auto-proceed to next step</action>
</check>

### 2. Handle YOLO Mode

<check if="agent_mode == 'yolo'">
  <action>Display: "Deploying YOLO-generated agent..."</action>

  <action>Check if YOLO subagent has completed:</action>
  - If still running: Wait for completion (check Task tool status)
  - If completed: Proceed with file move

  <action>Verify agent was created in parent worktree:</action>
  ```bash
  test -d "{{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}" && echo "FOUND" || echo "NOT_FOUND"
  ```

  <check if="agent NOT_FOUND in parent worktree">
    <action>Display error:</action>
    ```
    WARNING: YOLO agent not found in parent worktree

    Expected: {{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}/

    The YOLO subagent may have:
    - Failed to create the agent
    - Used a different name
    - Written to a different location (e.g., _cfg/agents instead of custom/src/agents)

    Checking for any agents created during this session...
    ```

    <action>List agents in parent worktree (correct location):</action>
    ```bash
    ls -la "{{parent_worktree_path}}/.bmad/custom/src/agents/" 2>/dev/null || echo "No custom/src/agents folder"
    ```

    <action>Also check wrong location (_cfg/agents) for diagnostic:</action>
    ```bash
    ls -la "{{parent_worktree_path}}/.bmad/_cfg/agents/" 2>/dev/null | grep -i "{{agent_name}}" || echo "Not in _cfg either"
    ```

    <action>If agent found with different name in correct location, update {{agent_name}} and continue</action>
    <action>If agent found in wrong location (_cfg), display error explaining the YOLO subagent wrote to wrong path</action>
    <action>If no agent found anywhere, display warning and continue without agent</action>
  </check>

  <check if="agent FOUND in parent worktree">
    <action>Create target directory in nook:</action>
    ```bash
    mkdir -p "{{worktree_path}}/.bmad/custom/src/agents"
    ```

    <action>MOVE agent folder from parent worktree to nook:</action>
    ```bash
    mv "{{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}" \
       "{{worktree_path}}/.bmad/custom/src/agents/"
    ```

    <action>Verify move succeeded:</action>
    ```bash
    test -d "{{worktree_path}}/.bmad/custom/src/agents/{{agent_name}}" && echo "SUCCESS" || echo "FAILED"
    ```

    <check if="move FAILED">
      <action>Display error and continue without agent</action>
    </check>

    <check if="move SUCCESS">
      <action>Display: "Agent '{{agent_name}}' deployed to nook"</action>
      <action>Set {{custom_agent}} = {{agent_name}}</action>
    </check>
  </check>

  <action>Proceed to next step</action>
</check>

### 3. Handle Full Wizard Mode

<check if="agent_mode == 'full'">
  <action>Display:</action>
  ```
  RUNNING FULL AGENT WIZARD

  You'll now go through the interactive agent creation process.
  The agent will be created directly in this nook.

  Task context: {{task_intent}}
  Lineage context available for reference.

  Starting wizard...
  ```

  <action>Change to nook directory:</action>
  ```bash
  cd "{{worktree_path}}"
  ```

  <action>Invoke create-agent workflow:</action>
  ```
  Load and execute {createAgentWorkflow}

  Pass context:
  - task_intent: {{task_intent}}
  - lineage_context: {{lineage_context}}
  - output_location: {{worktree_path}}/.bmad/custom/src/agents/
  ```

  <action>After wizard completes, capture {{custom_agent}} from workflow output</action>
  <action>Proceed to next step</action>
</check>

## MENU OPTIONS:

After agent deployment: Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Agent mode handled correctly (none/yolo/full)
- YOLO: Files moved from parent to nook
- Full: Wizard completed and agent created
- {{custom_agent}} set to agent name or "none"
- Agent files verified in nook

### ❌ SYSTEM FAILURE:
- Not waiting for YOLO subagent to complete
- Copying instead of moving files (leaves duplicates)
- Not verifying file move succeeded
- Full wizard not receiving context
