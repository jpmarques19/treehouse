---
name: 'step-06-agent-decision'
description: 'Agent wizard decision point - gather intent and choose YOLO or Full mode'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-06-agent-decision.md'
nextStepFile: '{workflow_path}/steps/step-07-validate-folder.md'

# Inputs (may be provided via workflow inputs)
# skip_wizard: optional - if true, skip agent creation entirely

# Agent References
nookContextAnalyst: 'nook-context-analyst'
createAgentWorkflow: '{project-root}/.bmad/bmb/workflows/create-agent/workflow.md'
---

# Step 6: Agent Decision

## STEP GOAL:

Offer custom agent creation with two tracks: YOLO (fast auto-generation) or Full (interactive wizard). If YOLO selected, spawn subagent immediately for parallel execution.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- ⏸️ HALT and wait for user decisions at each choice point
- 🎯 For YOLO: spawn subagent and continue WITHOUT waiting

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Explain the two tracks clearly
- ✅ YOLO runs in parallel - do not block

## EXECUTION SEQUENCE:

### 1. Check for Skip Parameter

<check if="skip_wizard input is true">
  <action>Display: "Agent Wizard skipped via input parameter..."</action>
  <action>Set {{agent_mode}} = "none"</action>
  <action>Proceed to next step</action>
</check>

### 2. Ask About Custom Agent

```
AGENT WIZARD (optional)

Would you like a custom agent for this nook?

A custom agent will have:
- Task-specific expertise and persona
- Pre-loaded context from your lineage
- Relevant workflows in its menu

[Y] Yes - Create a custom agent
[N] No  - Skip (clean nook)

Choose [Y/N]:
```

<action>HALT and wait for user input</action>

<check if="user selects N">
  <action>Display: "Skipping Agent Wizard - creating clean nook..."</action>
  <action>Set {{agent_mode}} = "none"</action>
  <action>Proceed to next step</action>
</check>

### 3. Choose Agent Creation Mode (BEFORE task description)

<check if="user selects Y">
  <action>Present mode selection:</action>

```
How should we create the agent?

[Y] YOLO - Fast auto-generation
    You provide a brief task description
    Agent created in background while nook is set up
    Quick and efficient

[F] Full - Interactive wizard
    Step-by-step agent customization
    More control over persona and settings
    Runs after nook is created

Choose [Y/F]:
```

  <action>HALT and wait for user input</action>
</check>

### 4. Handle YOLO Mode

<check if="user selects Y (YOLO)">
  <action>Set {{agent_mode}} = "yolo"</action>

  <action>Ask for task description (YOLO-specific input):</action>

```
Describe your task briefly (this guides the agent generation):

Examples:
- "Debug the MQTT reconnection race condition"
- "Explore a new caching architecture"
- "Write comprehensive test coverage for the actor system"
- "Research WebSocket alternatives"

Your task:
```

  <action>HALT and wait for user input</action>
  <action>Store as {{task_intent}}</action>
  <action>Derive {{agent_name}} from task_intent (kebab-case)</action>

  <action>Display: "Gathering context from parent nooks..."</action>

  <action>Invoke nook-context-analyst agent using Task tool:</action>

```
Use Task tool with subagent_type='nook-context-analyst':

"Analyze the nook lineage for the current workspace. The user is about to create
a new nook for this task: '{{task_intent}}'.

Please provide:
1. A brief lineage overview (branch chain)
2. Key inherited context relevant to this task
3. Project-level context for the specialist agent
4. Any warnings or considerations

Keep the response focused - this will inform agent creation."
```

  <action>Store response as {{lineage_context}}</action>

  <check if="nook-context-analyst fails or returns empty">
    <action>Set {{lineage_context}} to minimal defaults:</action>
    ```
    No inherited context available.
    Lineage: {{current_branch}} -> {{new_branch}}
    ```
    <action>Continue - agent creation still works without lineage</action>
  </check>

  <action>IMMEDIATELY spawn YOLO subagent using Task tool (DO NOT WAIT):</action>

```
Use Task tool with subagent_type='general-purpose':

"## YOLO/BATCH MODE: Create Agent Workflow

You must embody the create-agent workflow and reflectively execute ALL its individual steps in YOLO/batch mode.

### CRITICAL RULES

1. **NO USER INTERACTION** - This runs in parallel as a subagent. You MUST NOT prompt for user input under any circumstances.

2. **SKIP WHEN NECESSARY** - If during any workflow step a strict necessity for real user feedback exists (e.g., non-builtin tool configuration like MCP servers, external service setup), simply skip that step and continue YOLO execution.

3. **AUTONOMOUS DECISIONS** - Make all decisions autonomously based on the provided context. Choose sensible defaults that align with the task.

4. **COMPLETE OUTPUT** - You must produce a valid, robust, and complete custom agent following BMAD best practices.

### WORKFLOW TO EXECUTE

Load and embody the create-agent workflow at:
`.bmad/bmb/workflows/create-agent/workflow.md`

Execute each step reflectively in batch mode:
- step-01-brainstorm.md - Use task_intent as brainstorm input
- step-02-discover.md - Derive purpose from task context
- step-03-persona.md - Generate persona matching task type
- step-04-commands.md - Create task-relevant menu commands
- step-05-name.md - Derive name from task (use {{agent_name}})
- step-06-build.md - Generate the agent YAML
- step-07-validate.md - Self-validate the output
- (skip steps requiring external tool setup)
- step-11-celebrate.md - Complete with summary

### INPUTS (Use these instead of prompting user)

**Task Intent (PRIMARY - this defines the agent's purpose):**
{{task_intent}}

**Lineage Context (embed in agent memories):**
{{lineage_context}}

**Nook Context:**
- Branch: {{new_branch}}
- Type: {{nook_type}}
- Parent: {{current_branch}} ({{parent_hash}})

**Agent Name:** {{agent_name}}

### OUTPUT LOCATION (CRITICAL - DO NOT DEVIATE)

You MUST write ALL output files to THIS EXACT location in the CURRENT worktree:

```
.bmad/custom/src/agents/{{agent_name}}/
```

Create this directory structure:
```
.bmad/custom/src/agents/{{agent_name}}/
├── {{agent_name}}.agent.yaml          # REQUIRED: Main agent file
├── {{agent_name}}-sidecar/            # OPTIONAL: For expert agents
│   ├── memories.md
│   ├── instructions.md
│   └── knowledge/
│       └── [knowledge files].md
└── info-and-installation-guide.md     # OPTIONAL: Installation notes
```

⚠️ DO NOT write to `.bmad/_cfg/agents/` - that is for compiled output only
⚠️ DO NOT write `.md` files directly - agent source must be `.agent.yaml`
⚠️ The folder name MUST match the agent filename ({{agent_name}})

### REQUIRED OUTPUT

1. Create folder: `.bmad/custom/src/agents/{{agent_name}}/`
2. Write main file: `.bmad/custom/src/agents/{{agent_name}}/{{agent_name}}.agent.yaml`
3. For expert agents with substantial context, also create sidecar folder

### AGENT REQUIREMENTS

The generated agent MUST include:
- nook-context-analyst in critical_actions (for lineage access)
- Task intent and lineage embedded in memories
- Menu items relevant to {{nook_type}} work
- Persona tailored to the specific task

### COMPLETION

When done, return brief summary:
- Agent name created
- Exact file paths written (list each file)
- Ready for: npx bmad-method@alpha agent-install
"
```

  <action>Display: "YOLO agent creation started in background..."</action>
  <action>Continue immediately to next step (DO NOT WAIT for subagent)</action>
</check>

### 5. Handle Full Mode

<check if="user selects F (Full)">
  <action>Set {{agent_mode}} = "full"</action>
  <action>Display: "Full wizard will run after nook is created..."</action>
  <action>DO NOT ask for task description - create-agent workflow handles that</action>
  <action>DO NOT gather lineage context - create-agent workflow handles that</action>
  <action>Proceed to next step</action>
</check>

## MENU OPTIONS:

After mode selection: Load and execute `{nextStepFile}`

**CRITICAL**: For YOLO mode, spawn the subagent and continue immediately. Do NOT wait for completion - the subagent runs in parallel while the main workflow continues.

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- User decision captured (skip/YOLO/full)
- Mode selection comes BEFORE task description
- Task intent gathered only for YOLO mode (Full gathers it later)
- Lineage context gathered via nook-context-analyst (YOLO only)
- YOLO subagent spawned in parallel (if selected)
- {{agent_mode}} stored; {{task_intent}}, {{lineage_context}}, {{agent_name}} stored for YOLO
- Proceeded to next step without blocking on YOLO

### ❌ SYSTEM FAILURE:
- Asking for task description before mode selection
- Asking for task description when Full mode selected
- Gathering lineage context when Full mode selected
- Waiting for YOLO subagent to complete (should be parallel)
- Not gathering lineage context for YOLO
- Blocking the workflow on agent creation
