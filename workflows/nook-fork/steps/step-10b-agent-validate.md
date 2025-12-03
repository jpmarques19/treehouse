---
name: 'step-10b-agent-validate'
description: 'Validate Full wizard agent source against BMAD quality checklist'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-10b-agent-validate.md'
nextStepFile: '{workflow_path}/steps/step-11-create-context.md'

# Validation Reference
validationChecklist: '{project-root}/.bmad/bmb/workflows/create-agent/data/agent-validation-checklist.md'
---

# Step 10b: Validate Full Wizard Agent

## STEP GOAL:

Validate agent source created by Full wizard mode against BMAD checklist. Fix any issues before proceeding.

**Note:** This step only runs for Full wizard mode. YOLO agents are validated in step-10 before deployment.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- 🎯 Only runs for agent_mode == 'full'
- 🛑 HALT if critical issues found - offer to fix or abort

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Quality gate - don't let broken agents through
- ✅ Be helpful - offer fixes, don't just report problems

## EXECUTION SEQUENCE:

### 1. Check if Validation Needed

<check if="agent_mode != 'full'">
  <action>This step should only be reached for Full wizard mode</action>
  <action>If reached with agent_mode == 'none' or 'yolo', auto-proceed to step-11</action>
</check>

<check if="agent_mode == 'full'">
  <action>Display: "Validating Full wizard agent source..."</action>
  <action>Continue with validation</action>
</check>

### 2. Load Agent Source File

<action>Read the generated agent YAML from nook:</action>
```bash
cat "{{worktree_path}}/bmad/agents/{{custom_agent}}/{{custom_agent}}.agent.yaml"
```

<action>Store content for validation</action>

### 3. Run Validation Checks

<action>Run validation against {validationChecklist}:</action>

**YAML Structure Checks:**
- [ ] `agent.metadata.name` exists and is non-empty
- [ ] `agent.metadata.title` exists and is non-empty
- [ ] `agent.metadata.icon` exists (emoji)
- [ ] `agent.metadata.type` exists (simple/expert/module)
- [ ] `agent.persona` exists with role, identity, communication_style, principles

**Persona Field Separation (CRITICAL - #1 Quality Issue):**
- [ ] communication_style is 1-2 sentences MAX
- [ ] communication_style does NOT contain behavioral words: "ensures", "makes sure", "always", "never"
- [ ] communication_style does NOT contain identity words: "experienced", "expert who", "senior", "seasoned"
- [ ] communication_style does NOT contain philosophy words: "believes in", "focused on", "committed to"
- [ ] communication_style does NOT contain behavioral descriptions: "who does X", "that does Y"
- [ ] communication_style describes only HOW they talk (verbal patterns, word choice, quirks)

**Menu Validation:**
- [ ] `agent.menu` exists with at least one item
- [ ] Menu triggers do NOT start with `*` (auto-added by compiler)
- [ ] Each item has `description` field
- [ ] Each item has handler: `action`, `workflow`, `exec`, `tmpl`, or `data`
- [ ] If `action="#prompt-id"`, verify prompt with that ID exists in `agent.prompts`

**Prompts Validation (if present):**
- [ ] Each prompt has `id` field
- [ ] Each prompt has `content` field
- [ ] Prompt IDs are unique
- [ ] All referenced prompt IDs exist

**Quality Checks:**
- [ ] No placeholder text ({{NAME}}, {ROLE}, TODO, etc.)
- [ ] Filename is kebab-case ending in `.agent.yaml`
- [ ] Agent purpose is clear from persona

### 4. Report Validation Results

<check if="no issues found">
  <action>Display:</action>
  ```
  ✅ AGENT VALIDATION PASSED

  All checks passed:
  ✓ YAML structure valid
  ✓ Persona fields properly separated
  ✓ Menu items valid
  ✓ Prompts valid (if any)
  ✓ Quality checks passed

  Agent '{{custom_agent}}' is ready for compilation.
  ```
  <action>Auto-proceed to next step</action>
</check>

<check if="issues found">
  <action>Display:</action>
  ```
  ❌ AGENT VALIDATION FAILED

  Issues found in {{worktree_path}}/bmad/agents/{{custom_agent}}/

  {{list_issues_with_details}}

  Options:
  [F] Fix issues - I'll correct the agent source file
  [S] Skip agent - Remove agent and continue without it
  [A] Abort - Stop workflow and investigate manually

  Choose [F/S/A]:
  ```
  <action>HALT and wait for user input</action>
</check>

### 5. Handle User Choice

<check if="user selects F (Fix)">
  <action>Fix each issue in the agent source file:</action>

  **Common Fixes:**
  - Remove `*` prefix from menu triggers
  - Shorten communication_style to 1-2 sentences of pure verbal patterns
  - Move behavioral content from communication_style to principles
  - Move identity content from communication_style to identity
  - Add missing metadata fields (name, title, icon, type)
  - Fix broken prompt references in menu actions

  <action>Write corrected file:</action>
  ```bash
  # Write corrected YAML to:
  # {{worktree_path}}/bmad/agents/{{custom_agent}}/{{custom_agent}}.agent.yaml
  ```

  <action>Display: "Issues fixed. Re-validating..."</action>
  <action>Return to step 3 (re-validate)</action>
</check>

<check if="user selects S (Skip)">
  <action>Remove agent from nook:</action>
  ```bash
  rm -rf "{{worktree_path}}/bmad/agents/{{custom_agent}}"
  ```
  <action>Set {{custom_agent}} = "none"</action>
  <action>Set {{agent_mode}} = "none"</action>
  <action>Display: "Agent removed. Continuing without custom agent."</action>
  <action>Proceed to next step</action>
</check>

<check if="user selects A (Abort)">
  <action>Display:</action>
  ```
  Workflow aborted. The nook was created but agent validation failed.

  To investigate:
    cd {{worktree_path}}
    cat bmad/agents/{{custom_agent}}/{{custom_agent}}.agent.yaml

  To retry agent creation:
    /bmad:bmb:workflows:create-agent

  To remove the nook:
    git worktree remove {{worktree_path}}
  ```
  <action>HALT workflow</action>
</check>

## MENU OPTIONS:

- If validation passes: Auto-proceed to `{nextStepFile}`
- If issues found: Wait for [F/S/A] selection, then proceed or abort

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Validation runs for Full wizard agents
- All critical checks pass (or issues fixed)
- Agent source meets BMAD quality standards
- Proceeded to next step with valid agent

### ❌ SYSTEM FAILURE:
- Not validating Full wizard output
- Proceeding with known critical issues
- Not offering fix option for fixable issues
- Not re-validating after fixes
