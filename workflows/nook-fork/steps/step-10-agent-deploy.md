---
name: 'step-10-agent-deploy'
description: 'Deploy agent - validate YOLO source, fix issues, then move to nook; or run full wizard'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-10-agent-deploy.md'
nextStepFile: '{workflow_path}/steps/step-10b-agent-validate.md'

# Agent References
createAgentWorkflow: '{project-root}/.bmad/bmb/workflows/create-agent/workflow.md'
validationChecklist: '{project-root}/.bmad/bmb/workflows/create-agent/data/agent-validation-checklist.md'
---

# Step 10: Agent Deploy

## STEP GOAL:

Deploy agent based on agent_mode. For YOLO: validate source BEFORE moving (fix issues first). For Full: run wizard interactively.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 For YOLO: VALIDATE before moving - fix issues in parent worktree first
- 🎯 For Full: invoke create-agent workflow interactively

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Quality gate - don't deploy broken agents
- ✅ Fix issues before deployment, not after

## EXECUTION SEQUENCE:

### 1. Check Agent Mode

<check if="agent_mode == 'none'">
  <action>Display: "No custom agent requested - skipping agent deploy"</action>
  <action>Set {{custom_agent}} = "none"</action>
  <action>Auto-proceed to next step (skip step-10b validation)</action>
  <action>Load step-11-create-context.md directly</action>
</check>

### 2. Handle YOLO Mode

<check if="agent_mode == 'yolo'">
  <action>Display: "Processing YOLO-generated agent..."</action>

  <action>Check if YOLO subagent has completed:</action>
  - If still running: Wait for completion (check Task tool status)
  - If completed: Proceed with validation

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
    - Written to a different location (e.g., .bmad/custom/agents instead of .bmad/custom/src/agents)

    Checking for any agents created during this session...
    ```

    <action>List agents in parent worktree (correct location):</action>
    ```bash
    ls -la "{{parent_worktree_path}}/.bmad/custom/src/agents/" 2>/dev/null || echo "No .bmad/custom/src/agents folder"
    ```

    <action>Also check wrong location (.bmad/custom/agents) for diagnostic:</action>
    ```bash
    ls -la "{{parent_worktree_path}}/.bmad/custom/agents/" 2>/dev/null | grep -i "{{agent_name}}" || echo "Not in .bmad/custom/agents either"
    ```

    <action>If agent found with different name in correct location, update {{agent_name}} and continue</action>
    <action>If agent found in wrong location (.bmad/custom/agents), display error explaining the YOLO subagent wrote to compiled output path instead of source path</action>
    <action>If no agent found anywhere, set {{custom_agent}} = "none" and continue without agent</action>
  </check>

  <check if="agent FOUND in parent worktree">

    ### 2a. VALIDATE YOLO Agent BEFORE Moving (CRITICAL)

    <action>Display: "Validating YOLO-generated agent source..."</action>

    <action>Read agent source from parent worktree:</action>
    ```bash
    cat "{{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}/{{agent_name}}.agent.yaml"
    ```

    <action>Run validation against {validationChecklist}:</action>

    **YAML Structure Checks:**
    - [ ] `agent.metadata.name` exists and is non-empty
    - [ ] `agent.metadata.title` exists and is non-empty
    - [ ] `agent.metadata.icon` exists (emoji)
    - [ ] `agent.metadata.type` exists (simple/expert/module)
    - [ ] `agent.persona` exists with role, identity, communication_style, principles

    **Persona Field Separation (CRITICAL):**
    - [ ] communication_style is 1-2 sentences MAX
    - [ ] communication_style does NOT contain: "ensures", "makes sure", "always", "never"
    - [ ] communication_style does NOT contain: "experienced", "expert who", "senior"
    - [ ] communication_style does NOT contain: "believes in", "focused on", "committed to"
    - [ ] communication_style describes only HOW they talk (verbal patterns)

    **Menu Validation:**
    - [ ] `agent.menu` exists with at least one item
    - [ ] Menu triggers do NOT start with `*` (auto-added by compiler)
    - [ ] Each item has description and handler (action, workflow, exec, etc.)
    - [ ] If `action="#id"`, verify prompt with that ID exists

    **Quality Checks:**
    - [ ] No placeholder text ({{NAME}}, TODO, etc.)
    - [ ] Filename is kebab-case ending in `.agent.yaml`

    ### 2b. Handle Validation Results

    <check if="validation PASSED">
      <action>Display: "✅ YOLO agent validation passed"</action>
      <action>Continue to move agent to nook</action>
    </check>

    <check if="validation FAILED with issues">
      <action>Display:</action>
      ```
      ❌ YOLO AGENT VALIDATION FAILED

      Issues found in {{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}/

      {{list_issues_with_details}}

      Options:
      [F] Fix issues - I'll correct the agent source file in parent worktree
      [S] Skip agent - Continue without agent
      [A] Abort - Stop workflow and investigate manually

      Choose [F/S/A]:
      ```
      <action>HALT and wait for user input</action>
    </check>

    <check if="user selects F (Fix)">
      <action>Fix each issue in the agent source file (in parent worktree):</action>

      **Common Fixes:**
      - Remove `*` prefix from menu triggers
      - Shorten communication_style to 1-2 sentences of pure verbal patterns
      - Move behavioral content from communication_style to principles
      - Move identity content from communication_style to identity
      - Add missing metadata fields
      - Fix broken prompt references

      <action>Write corrected file to parent worktree:</action>
      ```bash
      # Write corrected YAML to:
      # {{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}/{{agent_name}}.agent.yaml
      ```

      <action>Display: "Issues fixed. Re-validating..."</action>
      <action>Return to step 2a (re-validate)</action>
    </check>

    <check if="user selects S (Skip)">
      <action>Remove agent from parent worktree:</action>
      ```bash
      rm -rf "{{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}"
      ```
      <action>Set {{custom_agent}} = "none"</action>
      <action>Set {{agent_mode}} = "none"</action>
      <action>Display: "Agent removed. Continuing without custom agent."</action>
      <action>Skip to step-11 (no validation needed)</action>
    </check>

    <check if="user selects A (Abort)">
      <action>Display:</action>
      ```
      Workflow aborted. The nook worktree was created but agent deployment failed.

      To investigate the YOLO agent:
        cat {{parent_worktree_path}}/.bmad/custom/src/agents/{{agent_name}}/{{agent_name}}.agent.yaml

      To remove the nook:
        git worktree remove {{worktree_path}}
      ```
      <action>HALT workflow</action>
    </check>

    ### 2c. Move Validated Agent to Nook

    <check if="validation passed (or fixed and re-validated)">
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
        <action>Display error and set {{custom_agent}} = "none"</action>
      </check>

      <check if="move SUCCESS">
        <action>Display: "✅ Agent '{{agent_name}}' deployed to nook (validated)"</action>
        <action>Set {{custom_agent}} = {{agent_name}}</action>
      </check>
    </check>
  </check>

  <action>Skip step-10b (already validated) - proceed to step-11</action>
</check>

### 3. Handle Full Wizard Mode

<check if="agent_mode == 'full'">
  <action>Display:</action>
  ```
  RUNNING FULL AGENT WIZARD

  You'll now go through the interactive agent creation process.
  The agent will be created directly in this nook.

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
  - output_location: {{worktree_path}}/.bmad/custom/src/agents/
  ```

  <action>After wizard completes, capture {{custom_agent}} from workflow output</action>
  <action>Proceed to step-10b for final validation</action>
</check>

## MENU OPTIONS:

- For agent_mode == 'none': Skip to step-11-create-context.md
- For agent_mode == 'yolo': Validate → Fix → Move → Skip to step-11-create-context.md
- For agent_mode == 'full': Run wizard → Proceed to step-10b-agent-validate.md

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Agent mode handled correctly (none/yolo/full)
- YOLO: Validated BEFORE moving, issues fixed if found
- YOLO: Only valid agents deployed to nook
- Full: Wizard completed, agent passed to step-10b for validation
- {{custom_agent}} set to agent name or "none"

### ❌ SYSTEM FAILURE:
- Moving YOLO agent without validating first
- Deploying agent with known validation issues
- Not offering fix option for fixable issues
- Not re-validating after fixes
