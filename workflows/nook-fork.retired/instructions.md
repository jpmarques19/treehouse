# Nook Fork - Create Isolated Development Nook

<critical>The workflow execution engine is governed by: {project-root}/.bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/.bmad/th/workflows/nook-fork/workflow.yaml</critical>

## Purpose

Create an **isolated development nook** - a new git worktree for focused work:
- Nook has config files from git (skip-worktree applied - local changes ignored)
- Nook does NOT have docs/ (gitignored - must restore if needed)
- Create focused workspace for bugs, spikes, or exploration
- Lineage tracked via branch name prefix (hash) and context.yaml

**Key Architecture Points:**
- Source of truth: `base_workspace_path/.bmad-tracking/` (git submodule)
- Three-tier artifact handling:
  - **skip_worktree_paths** (`.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml`): Exist from git, skip-worktree applied
  - **submodule_paths** (`.bmad-tracking/`, `.gitmodules`): Skip-worktree applied, then removed from nook
  - **gitignore_paths** (`docs/`): Don't exist, must restore if needed
- User decides when to sync (no forced auto-sync)
- To get docs in nook: run `nook-restore` after creating nook

## Branch Naming Scheme

Uses **hash-based lineage** to avoid git ref conflicts while maintaining traceability:

```
{type}/{parent-hash-4-chars}-{short-name}
```

**Examples:**
```
peek-a-box-mvp                    <- Base workspace (hash: a1b2)
explore/a1b2-provision-race       <- Explore nook from a1b2 (hash: 7f3e)
spike/7f3e-mutex-solution         <- Spike nook from 7f3e (hash: c9d1)
bugfix/c9d1-deadlock-retry        <- Bugfix nook from c9d1 (hash: jf2b)
```

<workflow>

<step n="1" goal="Check if workspace is initialized">
  <action>Read {module_config} and check base_workspace_path:</action>
  ```bash
  grep "^base_workspace_path:" .bmad/th/config.yaml
  ```

  <check if="base_workspace_path is empty or not set">
    <action>Display critical halt:</action>
    ```
    CANNOT CREATE NOOK - WORKSPACE NOT INITIALIZED

    This workspace hasn't been initialized as a base workspace.
    The base workspace holds the tracking folder (.bmad-tracking/)
    which is the source of truth for all nooks.

    Please run first:
      /bmad:th:workflows:treehouse-init

    This will:
    1. Add sync artifacts to .gitignore
    2. Create the tracking folder
    3. Establish this worktree as the base workspace
    ```
    <action>HALT workflow - do not proceed</action>
  </check>

  <action>Store base_workspace_path as {{base_path}}</action>
  <action>Store base_branch from config as {{base_branch}}</action>
  <action>Display: "Base workspace: {{base_path}} ({{base_branch}})"</action>
</step>

<step n="2" goal="Get current branch info and validate tracking">
  <action>Get current branch: `git branch --show-current`</action>
  <action>Store as {{current_branch}}</action>

  <action>Get short hash of current HEAD (4 chars): `git rev-parse --short=4 HEAD`</action>
  <action>Store as {{parent_hash}}</action>

  <action>Check if tracking folder has uncommitted changes:</action>
  ```bash
  git -C "{{base_path}}" status --porcelain .bmad-tracking/
  ```

  <check if="tracking folder has modifications (non-empty output)">
    <action>Display warning with options:</action>
    ```
    WARNING - TRACKING FOLDER HAS MODIFICATIONS

    The tracking folder at {{base_path}}/.bmad-tracking/
    has uncommitted changes. This may indicate:
    - A sync operation was interrupted
    - Manual edits to context files
    - Work in progress that should be committed

    Uncommitted changes:
    {{list first 10 files from git status}}

    RECOMMENDATION: Commit these changes before creating nook
    to ensure a clean baseline for the new nook.

    Options:
    [C] Commit changes now - Stage and commit tracking folder
    [P] Proceed anyway    - Create nook without committing (not recommended)
    [A] Abort             - Exit and handle manually

    Choose [C/P/A]:
    ```

    <check if="user_choice == C or user_choice == commit">
      <action>Stage and commit tracking changes:</action>
      ```bash
      cd "{{base_path}}"
      git add .bmad-tracking/
      git commit -m "chore(th): commit tracking changes before nook

      Generated with Claude Code"
      ```
      <action>Display: "Tracking changes committed. Continuing..."</action>
    </check>

    <check if="user_choice == P or user_choice == proceed">
      <action>Display: "Proceeding with uncommitted tracking changes..."</action>
    </check>

    <check if="user_choice == A or user_choice == abort">
      <action>Display: "Aborted. Please handle tracking changes and try again."</action>
      <action>HALT workflow</action>
    </check>
  </check>

  <check if="tracking folder is clean">
    <action>Display: "Tracking folder is clean"</action>
  </check>
</step>

<step n="3" goal="Check if current branch has context synced">
  <action>Check if context.yaml exists for current branch:</action>
  ```bash
  test -f "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml"
  ```

  <check if="context.yaml does NOT exist for current branch">
    <action>Display informational notice:</action>
    ```
    NOTE - NO SYNCED CONTEXT FOR CURRENT BRANCH

    Branch '{{current_branch}}' hasn't been synced yet.
    This means the nook won't inherit any context lineage.

    This is fine if:
    - This is a new workspace just starting out
    - You haven't made changes worth preserving yet

    If you want to save context first:
      /bmad:th:workflows:nook-sync

    Continuing with nook creation...
    ```
  </check>

  <check if="context.yaml exists">
    <action>Load existing context for lineage:</action>
    ```bash
    cat "{{base_path}}/.bmad-tracking/{{current_branch}}/context.yaml"
    ```
    <action>Extract {{existing_lineage}} array from context</action>
  </check>
</step>

<step n="4" goal="Select nook type">
  <check if="nook_type provided as input">
    <action>Validate nook_type is one of: explore, spike, bugfix, discovery, feature, experiment, hotfix</action>
    <action>Store as {{nook_type}}</action>
  </check>

  <check if="nook_type NOT provided">
    <ask>What kind of nook do you want to create?

    1. **explore/** - Deep exploration or investigation
    2. **spike/** - Quick prototype or proof of concept
    3. **bugfix/** - Focused bug investigation and fix
    4. **discovery/** - Requirements or architecture discovery
    5. **feature/** - New feature branch
    6. **experiment/** - Wild experimentation
    7. **hotfix/** - Production hotfix
    8. **custom** - Custom branch prefix

    Choose [1-8]:</ask>

    <action>Map choice to type:
      1 -> explore
      2 -> spike
      3 -> bugfix
      4 -> discovery
      5 -> feature
      6 -> experiment
      7 -> hotfix
      8 -> ask for custom type (validate: lowercase, no slashes)
    </action>
    <action>Store as {{nook_type}}</action>
  </check>
</step>

<step n="5" goal="Define nook name">
  <check if="nook_name provided as input">
    <action>Validate nook_name is kebab-case</action>
    <action>Store as {{nook_name}}</action>
  </check>

  <check if="nook_name NOT provided">
    <ask>What should this nook be called? (kebab-case, e.g., "race-condition", "mqtt-reconnect")

    The full branch will be: {{nook_type}}/{{parent_hash}}-<your-name></ask>

    <action>Validate name:
      - kebab-case only
      - no spaces or special chars
      - no slashes
    </action>
    <action>Store as {{nook_name}}</action>
  </check>

  <action>Construct full branch name: {{nook_type}}/{{parent_hash}}-{{nook_name}}</action>
  <action>Store as {{new_branch}}</action>

  <action>Display preview:</action>
  ```
  Nook Preview

  Parent branch:  {{current_branch}} ({{parent_hash}})
  New branch:     {{new_branch}}
  Nook type:      {{nook_type}}
  ```
</step>

<step n="5a" goal="Agent Wizard - Gather Task Intent">
  <check if="skip_wizard input is true">
    <action>Display: "Agent Wizard skipped via input parameter..."</action>
    <action>Skip to step 6</action>
  </check>

  <ask>
  AGENT WIZARD (optional)

  What are you trying to achieve in this nook?

  Describe your task and I'll create a custom specialist agent tailored
  specifically for this work. The agent will have:
  - A unique name and personality
  - Task-specific expertise and principles
  - Pre-loaded context about your goal
  - Relevant workflows in its menu
  - Custom prompts for quick actions

  Examples:
  - "Debug the MQTT reconnection race condition"
  - "Explore a new caching architecture"
  - "Write comprehensive test coverage for the actor system"
  - "Research WebSocket alternatives"
  - "Refactor the playlist management code"

  Press Enter to skip (creates clean fork with default configs)
  </ask>

  <check if="user provides description (non-empty input)">
    <action>Store as {{task_intent}}</action>
    <action>Proceed to step 5a-context</action>
  </check>

  <check if="user presses Enter / provides empty input">
    <action>Display: "Skipping Agent Wizard - creating clean nook..."</action>
    <action>Set {{wizard_skipped}} = true</action>
    <goto step="6" reason="User skipped Agent Wizard"/>
  </check>
</step>

<step n="5a-context" goal="Agent Wizard - Gather Big Picture Context">
  <action>Display:</action>
  ```
  Gathering context from parent nooks...
  This will help create a well-informed specialist agent.
  ```

  <action>Invoke nook-context-analyst agent to gather lineage context:</action>
  <critical>
  Use the Task tool with subagent_type='nook-context-analyst' to:
  1. Trace the lineage of the current nook/branch
  2. Read summaries from all parent context.yaml files
  3. Gather inherited decisions, constraints, and project context
  4. Focus on context relevant to {{task_intent}}

  Prompt for the agent:
  "Analyze the nook lineage for the current workspace. The user is about to create
  a new nook for this task: '{{task_intent}}'.

  Please provide:
  1. A brief lineage overview (branch chain)
  2. Key inherited context relevant to this task (decisions, constraints, architecture)
  3. Project-level context that should be embedded in a specialist agent's memory
  4. Any warnings or considerations for this type of work

  Keep the response focused and actionable - this will inform agent creation."
  </critical>

  <action>Store the context analyst's response as {{lineage_context}}</action>

  <check if="nook-context-analyst fails, times out, or returns empty/error">
    <action>Display:</action>
    ```
    Note: Could not gather full lineage context.
    This may happen if:
    - No parent context.yaml files exist yet
    - The lineage is shallow (close to base workspace)
    - The context analyst encountered an error

    Proceeding with task intent only - your agent will still be customized
    for "{{task_intent}}" but without inherited context.
    ```
    <action>Set {{lineage_context}} to minimal defaults:</action>
    ```yaml
    lineage_context:
      brief_summary: "No inherited context available"
      inherited: []
      constraints: []
      parent_findings: null
      branch_chain: "{{current_branch}} -> {{new_branch}}"
      memory_count: 0
    ```
    <action>Continue to step 5b (agent selection still works without lineage)</action>
  </check>

  <check if="nook-context-analyst succeeds">
    <action>Extract key elements:</action>
    - {{inherited_decisions}} - decisions from parent nooks
    - {{project_constraints}} - constraints that apply
    - {{architecture_context}} - relevant architecture info
    - {{parent_purpose}} - what parent nooks were working on
  </check>

  <action>Display summary to user:</action>
  ```
  Lineage Context Gathered

  {{lineage_context.brief_summary}}

  Relevant Inherited Context:
  {{for each item in inherited_context, max 5}}
  - {{item}}
  {{end for}}

  This context will be embedded in your specialist agent.
  ```

  <action>Proceed to step 5b</action>
</step>

<step n="5b" goal="Agent Wizard - Select Base Agent">
  <action>Load list of available base agents from .bmad/_cfg/agents/*.customize.yaml</action>

  <action>Analyze {{task_intent}} and match to best base agent using this mapping:</action>

  <agent_selection_matrix>
  | Task Pattern Keywords | Best Base Agent | Agent Title | Reasoning |
  |----------------------|-----------------|-------------|-----------|
  | debug, fix, bug, error, race, crash, broken, issue | bmm-dev | Developer Agent | Senior engineer for implementation and debugging |
  | explore, investigate, research, understand, analyze, dig | bmm-analyst | Business Analyst | Discovery and analysis expert |
  | architecture, design, system, scale, infrastructure, api | bmm-architect | Architect | System design specialist |
  | test, coverage, QA, validation, verify, quality | bmm-tea | Test Engineer Agent | Testing and quality expert |
  | docs, documentation, explain, write, readme, guide | bmm-tech-writer | Technical Writer | Documentation specialist |
  | ui, ux, user experience, interface, design, wireframe | bmm-ux-designer | UX Designer | Design specialist |
  | brainstorm, ideas, creative, innovate, think | cis-brainstorming-coach | Brainstorming Coach | Creative facilitation |
  | problem, solve, challenge, stuck, complex | cis-creative-problem-solver | Problem Solver | Systematic problem-solving |
  | story, narrative, presentation, pitch, communicate | cis-storyteller | Storyteller | Communication specialist |
  | sprint, agile, tickets, stories, backlog, jira | bmm-sm | Scrum Master | Sprint and agile management |
  | requirements, features, specs, product, roadmap | bmm-pm | Product Manager | Product management |
  | framework, library, code patterns, refactor, clean | bmm-frame-expert | Framework Expert | Framework and patterns expertise |
  </agent_selection_matrix>

  <action>If task matches multiple patterns equally, or no clear match:</action>
  <ask>I found multiple agents that could help. Which best fits your needs?

  1. {{agent_option_1}} - {{agent_1_reasoning}}
  2. {{agent_option_2}} - {{agent_2_reasoning}}
  3. {{agent_option_3}} - {{agent_3_reasoning}}

  Choose [1-3]:
  </ask>

  <action>Store selected agent as {{base_agent}}</action>
  <action>Store agent title as {{base_agent_title}}</action>
  <action>Store agent file path as {{base_agent_path}} (.bmad/_cfg/agents/{{base_agent}}.customize.yaml)</action>

  <action>Display: "Selected base agent: {{base_agent_title}}"</action>
</step>

<step n="5c" goal="Agent Wizard - Generate Custom Agent Profile">
  <action>Based on {{task_intent}} and {{base_agent}}, generate customization:</action>

  <critical>
  The template below is a STRUCTURAL GUIDE for generating YAML content.
  The {{for each}} and {{if}} blocks are pseudo-code indicating iteration/conditionals.
  Generate actual YAML by:
  1. Reading the structure
  2. Replacing placeholders with real generated values
  3. Expanding loops with actual items
  4. Evaluating conditionals based on available data
  </critical>

  <customization_generation>
    Generate the following customization YAML based on the task intent:

    1. **agent.metadata.name** - Create a fitting human name:
       - MUST be a proper human first name (not a job title or descriptor)
       - Choose a name that resonates with and reads well alongside the role
       - Name must NOT repeat words from the role - avoid redundancy
       - Examples: "Aria", "Felix", "Mira", "Quinn", "Sage", "Nova"

    2. **persona.role** - Specialize the role:
       - Start with base agent's role
       - Add task-specific specialization
       - Example: "Concurrency Bug Hunter + Thread Safety Expert"

    3. **persona.identity** - Add specific expertise:
       - Include relevant technical areas from task description
       - Reference specific technologies, patterns, or domains mentioned
       - Example: "Expert in Python threading, asyncio, and actor model concurrency"

    4. **persona.communication_style** - Tailor to task type:
       - Debugging tasks: More methodical, precise, investigative
       - Research tasks: More exploratory, questioning, curious
       - Creative tasks: More enthusiastic, building on ideas
       - Example: "Thinks out loud through timing sequences. Uses diagrams to visualize thread interleaving."

    5. **persona.principles** - Add task-specific principles (3-5):
       - Ground rules specific to this type of work
       - Example for debugging: "Always assume timing issue until proven otherwise"
       - Example for research: "Question every assumption"

    6. **critical_actions** (2-4 items):
       - Task-specific actions the agent must always do
       - Reference project-specific patterns if relevant
       - Example: "Check BaseMPDConnection._execute_with_lock() for thread safety"

       <critical>
       ALWAYS include this critical action:
       - "For broader project context or to understand inherited constraints, invoke the nook-context-analyst agent (Task tool with subagent_type='nook-context-analyst'). This agent can trace lineage and gather context from parent nooks."
       </critical>

    7. **memories** (5-10 items):
       - Pre-load context about the task
       - Include the task description itself
       - Add relevant project-specific knowledge
       - Example: "Currently investigating MQTT reconnection timing issues"

       <critical>
       EMBED LINEAGE CONTEXT IN MEMORIES:
       Use {{lineage_context}} gathered in step 5a-context to add:
       - Key inherited decisions from parent nooks
       - Project-level constraints that affect this work
       - Architecture patterns relevant to the task
       - Parent nook's purpose/findings if related
       - Any warnings or considerations identified

       Example memories from lineage context:
       - "INHERITED: Project uses Pykka actor architecture - all actors must be supervised"
       - "INHERITED: MPD connections require _execute_with_lock() for thread safety"
       - "PARENT CONTEXT: Boot flow validation discovered PipeWire audio instability"
       - "CONSTRAINT: Docker deployment is mandatory for edge device"
       - "LINEAGE: This work traces back to peek-a-box-mvp (9d4d) epic"
       </critical>

    8. **menu** (1-3 items):
       ONLY use existing workflows from these modules:
       - bmm: dev-story, code-review, story-context, create-story, research, brainstorm-project, architecture, create-diagram, create-dataflow, create-flowchart, create-wireframe, document-project, prd
       - core: brainstorming, party-mode
       - cis: problem-solving, design-thinking, storytelling, innovation-strategy

       Match workflows to task:
       - Debug tasks: code-review, story-context
       - Research tasks: research, brainstorm-project
       - Architecture: architecture, create-diagram, create-dataflow
       - Testing: dev-story (for writing tests)
       - Documentation: document-project

    9. **prompts** (1-2 custom inline prompts):
       - Create task-specific quick actions
       - Example for debugging: prompt to analyze recent git changes
       - Example for research: prompt to summarize findings
       - Each prompt needs: id (kebab-case), content (the instructions)
  </customization_generation>

  <action>Generate YAML content matching the customize.yaml schema:</action>
  ```yaml
  # Agent Customization - Generated by Agent Wizard
  # Nook: {{new_branch}}
  # Task: {{task_intent}}
  # Lineage: {{lineage_context.branch_chain}}

  agent:
    metadata:
      name: "{{generated_custom_name}}"

  persona:
    role: "{{generated_role}}"
    identity: "{{generated_identity}}"
    communication_style: "{{generated_style}}"
    principles:
      {{for each principle}}
      - "{{principle}}"
      {{end for}}

  critical_actions:
    # REQUIRED: Context analyst availability
    - "For broader project context or inherited constraints, use the nook-context-analyst agent (Task tool with subagent_type='nook-context-analyst'). It traces lineage and gathers context from parent nooks."
    # Task-specific critical actions
    {{for each action}}
    - "{{action}}"
    {{end for}}

  memories:
    # Task intent
    - "TASK: {{task_intent}}"
    # Lineage context (from nook-context-analyst)
    {{for each inherited_item in lineage_context.inherited}}
    - "INHERITED: {{inherited_item}}"
    {{end for}}
    {{for each constraint in lineage_context.constraints}}
    - "CONSTRAINT: {{constraint}}"
    {{end for}}
    {{if lineage_context.parent_findings}}
    - "PARENT CONTEXT: {{lineage_context.parent_findings}}"
    {{end if}}
    - "LINEAGE: {{lineage_context.branch_chain}}"
    # Task-specific memories
    {{for each memory}}
    - "{{memory}}"
    {{end for}}

  menu:
    {{for each menu_item}}
    - trigger: {{menu_item.trigger}}
      workflow: "{{menu_item.workflow_path}}"
      description: "{{menu_item.description}}"
    {{end for}}

  prompts:
    {{for each prompt}}
    - id: {{prompt.id}}
      content: |
        {{prompt.content}}
    {{end for}}
  ```

  <action>Store generated YAML as {{customization_yaml}}</action>
  <action>Store custom name as {{custom_name}}</action>
</step>

<step n="5d" goal="Agent Wizard - User Approval">
  <action>Display customization preview:</action>

  ```
  ╔═══════════════════════════════════════════════════════════════════╗
  ║                     NOOK AGENT PREVIEW                            ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║  Base Agent: {{base_agent_title}}                                 ║
  ║  Task: {{task_intent_short}}                                      ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║                                                                   ║
  ║  Your Custom Specialist:                                          ║
  ║  ┌───────────────────────────────────────────────────────────┐   ║
  ║  │  {{custom_name}}                                          │   ║
  ║  │  "{{custom_identity_short}}"                              │   ║
  ║  └───────────────────────────────────────────────────────────┘   ║
  ║                                                                   ║
  ║  Role: {{custom_role}}                                            ║
  ║  Style: {{custom_style_preview}}                                  ║
  ║                                                                   ║
  ║  Principles:                                                      ║
  ║  {{for each principle, max 3}}                                    ║
  ║    - {{principle}}                                                ║
  ║  {{end for}}                                                      ║
  ║                                                                   ║
  ╠───────────────────────────────────────────────────────────────────╣
  ║  INHERITED CONTEXT (from lineage):                                ║
  ║  {{for each item in lineage_context.inherited, max 3}}            ║
  ║    ◆ {{item}}                                                     ║
  ║  {{end for}}                                                      ║
  ╠───────────────────────────────────────────────────────────────────╣
  ║                                                                   ║
  ║  Pre-loaded Memories:                                             ║
  ║  {{for each memory, max 3}}                                       ║
  ║    - {{memory_short}}                                             ║
  ║  {{end for}}                                                      ║
  ║  + {{lineage_context.memory_count}} inherited context items       ║
  ║                                                                   ║
  ║  Critical Actions:                                                ║
  ║    ⚡ nook-context-analyst available for big picture              ║
  ║  {{for each action, max 2}}                                       ║
  ║    ⚡ {{action_short}}                                            ║
  ║  {{end for}}                                                      ║
  ║                                                                   ║
  ║  Added Menu Commands:                                             ║
  ║  {{for each menu_item}}                                           ║
  ║    *{{menu_item.trigger}} - {{menu_item.description}}             ║
  ║  {{end for}}                                                      ║
  ║                                                                   ║
  ║  Custom Prompts:                                                  ║
  ║  {{for each prompt}}                                              ║
  ║    #{{prompt.id}} - {{prompt.preview}}                            ║
  ║  {{end for}}                                                      ║
  ║                                                                   ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║  Lineage: {{lineage_context.branch_chain}}                        ║
  ║  This will modify: .bmad/_cfg/agents/{{base_agent}}.customize.yaml║
  ║  (Skip-worktree applied - changes stay local to this nook)        ║
  ╚═══════════════════════════════════════════════════════════════════╝
  ```

  <ask>
  Options:
  [A] Accept  - Create nook with this custom agent
  [M] Modify  - Tell me what to change
  [S] Skip    - Create clean nook without customization
  [C] Cancel  - Abort nook creation

  Choose [A/M/S/C]:
  </ask>

  <check if="user_choice == A or Accept or accept">
    <action>Display: "Agent customization approved!"</action>
    <action>Keep {{customization_yaml}} for application in Step 8a</action>
    <action>Continue to Step 6</action>
  </check>

  <check if="user_choice == M or Modify or modify">
    <ask>What would you like to change?

    You can say things like:
    - "Make the name more serious"
    - "Add a memory about XYZ"
    - "Remove the code-review workflow"
    - "Change the communication style to be more casual"
    </ask>
    <action>Update {{customization_yaml}} based on user feedback</action>
    <action>Return to display preview (loop until Accept/Skip/Cancel)</action>
  </check>

  <check if="user_choice == S or Skip or skip">
    <action>Clear {{customization_yaml}}</action>
    <action>Set {{wizard_skipped}} = true</action>
    <action>Display: "Creating clean nook without agent customization..."</action>
    <action>Continue to Step 6</action>
  </check>

  <check if="user_choice == C or Cancel or cancel">
    <action>Display: "Nook creation cancelled."</action>
    <action>HALT workflow</action>
  </check>
</step>

<step n="6" goal="Validate worktrees folder">
  <action>Read worktrees_folder from {module_config}</action>

  <check if="worktrees_folder not configured or doesn't exist">
    <ask>Where should worktrees be created? (e.g., /home/user/my-workspace-worktrees)</ask>
    <action>Validate path is writable</action>
    <action>Create directory if needed: `mkdir -p {{worktrees_folder}}`</action>
    <action>Update {module_config} with new worktrees_folder value</action>
  </check>

  <action>Store as {{worktrees_folder}}</action>
</step>

<step n="7" goal="Create the worktree">
  <action>Derive worktree folder name: {{parent_hash}}-{{nook_name}}</action>
  <action>Calculate full worktree path: {{worktrees_folder}}/{{parent_hash}}-{{nook_name}}</action>
  <action>Store as {{worktree_path}}</action>

  <action>Check path doesn't already exist:</action>
  ```bash
  test -d "{{worktree_path}}"
  ```

  <check if="path exists">
    <action>Display error:</action>
    ```
    ERROR - WORKTREE PATH EXISTS

    Path already exists: {{worktree_path}}

    Options:
    - Choose a different nook name
    - Remove the existing worktree first:
        git worktree remove "{{worktree_path}}"
    ```
    <action>HALT workflow</action>
  </check>

  <action>Display what will happen:</action>
  ```
  Creating Nook

  New branch:    {{new_branch}}
  Base:          {{current_branch}} ({{parent_hash}})
  Worktree path: {{worktree_path}}
  ```

  <ask>Proceed with nook creation? [y/n]</ask>

  <check if="user confirms">
    <action>Execute worktree creation:</action>
    ```bash
    git worktree add "{{worktree_path}}" -b "{{new_branch}}" "{{current_branch}}"
    ```
  </check>

  <check if="user declines">
    <action>Display: "Nook creation cancelled."</action>
    <action>HALT workflow</action>
  </check>
</step>

<step n="8" goal="Apply skip-worktree and remove submodule from nook">
  <action>Display:</action>
  ```
  Configuring nook isolation...

  Applying skip-worktree to:
  - Config files (local changes won't show in git status)
  - Submodule pointer and .gitmodules (will be removed from nook)
  ```

  <action>Read skip_worktree_paths from {module_config}:</action>
  ```yaml
  skip_worktree_paths:
    - .bmad/_cfg/agents
    - .bmad/bmm/config.yaml
  ```

  <action>Apply skip-worktree to config files:</action>
  ```bash
  cd "{{worktree_path}}"

  # For directories - apply to all files recursively
  git ls-files .bmad/_cfg/agents/ | xargs -r git update-index --skip-worktree

  # For single files
  git update-index --skip-worktree .bmad/bmm/config.yaml 2>/dev/null || true
  ```

  <action>Apply skip-worktree to submodule pointer and .gitmodules:</action>
  ```bash
  cd "{{worktree_path}}"

  # Skip-worktree on submodule pointer (tracked as gitlink in index)
  git update-index --skip-worktree .bmad-tracking 2>/dev/null || true

  # Skip-worktree on .gitmodules file
  git update-index --skip-worktree .gitmodules 2>/dev/null || true
  ```

  <action>Remove submodule placeholder and .gitmodules from nook:</action>
  ```bash
  cd "{{worktree_path}}"

  # Remove empty submodule placeholder folder
  rm -rf .bmad-tracking/

  # Remove .gitmodules file
  rm -f .gitmodules
  ```

  <action>Verify isolation and clean git status:</action>
  ```bash
  cd "{{worktree_path}}"
  echo "Skip-worktree files: $(git ls-files -v | grep "^S" | wc -l)"
  git status --short
  ```

  <check if="git status is clean (no output from status --short)">
    <action>Display:</action>
    ```
    Nook isolation configured:
      - Config files: skip-worktree applied (local changes ignored)
      - Submodule (.bmad-tracking/): removed from nook
      - .gitmodules: removed from nook
      - Git status: clean
    ```
  </check>

  <check if="git status shows changes">
    <action>Display warning:</action>
    ```
    Nook isolation incomplete - git status shows uncommitted changes.
    This may indicate skip-worktree didn't apply correctly.

    Try manually:
      git update-index --skip-worktree .bmad-tracking
      git update-index --skip-worktree .gitmodules
    ```
  </check>
</step>

<step n="8a" goal="Apply Agent Customization (if wizard was used)">
  <check if="{{customization_yaml}} exists AND {{wizard_skipped}} is NOT true">
    <action>Display:</action>
    ```
    Applying agent customization...
    ```

    <action>Write customization to nook's customize.yaml:</action>
    ```bash
    # The file path in the NEW nook
    CUSTOMIZE_FILE="{{worktree_path}}/.bmad/_cfg/agents/{{base_agent}}.customize.yaml"
    ```

    <action>Write the generated {{customization_yaml}} content to the file:</action>
    <action>The file already exists (from git with skip-worktree applied)</action>
    <action>Overwrite with the generated customization YAML</action>

    <action>Verify the file was written:</action>
    ```bash
    head -5 "{{worktree_path}}/.bmad/_cfg/agents/{{base_agent}}.customize.yaml"
    ```

    <action>Display success:</action>
    ```
    Custom agent "{{custom_name}}" ready!

    Customization saved to:
      {{worktree_path}}/.bmad/_cfg/agents/{{base_agent}}.customize.yaml

    To activate after switching to the nook:
      /bmad:bmm:agents:{{base_agent_cmd}}

    (Customization is loaded at runtime when agent activates)
    (Skip-worktree applied - changes stay local to this nook)
    ```

    <action>Store {{custom_agent_configured}} = true</action>
  </check>

  <check if="{{customization_yaml}} does NOT exist OR {{wizard_skipped}} is true">
    <action>Display: "Using default agent configurations..."</action>
    <action>Store {{custom_agent_configured}} = false</action>
  </check>
</step>

<step n="9" goal="Create context.yaml for new nook in BASE workspace">
  <action>Build lineage for new nook:</action>

  <action>If parent has existing_lineage, use it as base</action>
  <action>Add parent entry if not already in lineage:</action>
  ```yaml
  parent_entry:
    branch: {{current_branch}}
    hash: {{parent_hash}}
    type: (extract from branch name or "base")
    name: (derive human name from branch)
  ```

  <action>Create tracking directory for new nook in base workspace:</action>
  ```bash
  mkdir -p "{{base_path}}/.bmad-tracking/{{new_branch}}"
  ```

  <action>Create context.yaml at {{base_path}}/.bmad-tracking/{{new_branch}}/context.yaml:</action>
  ```yaml
  # Context Manifest - Nook
  # Auto-generated by nook-fork
  branch: {{new_branch}}
  type: {{nook_type}}
  forked_from: {{current_branch}}
  forked_at: {{date}}
  forked_by: {user_name}

  # Parent reference
  parent:
    branch: {{current_branch}}
    hash: {{parent_hash}}

  # Full lineage chain (oldest to newest, ending with this nook)
  lineage:
    {{for each entry in existing_lineage}}
    - branch: {{entry.branch}}
      hash: {{entry.hash}}
      type: {{entry.type}}
      name: "{{entry.name}}"
    {{end for}}
    - branch: {{new_branch}}
      hash: pending
      type: {{nook_type}}
      name: "{{nook_name}}"
      parent_hash: {{parent_hash}}
      created: {{date}}

  # Artifacts (none synced yet - nook starts clean)
  artifacts:
    docs:
      synced: false
    bmad_cfg_agents:
      synced: false
    bmad_bmm_config:
      synced: false

  # Agent Wizard Customization (if used)
  agent_customization:
    enabled: {{custom_agent_configured or false}}
    {{if custom_agent_configured}}
    base_agent: {{base_agent}}
    custom_name: "{{custom_name}}"
    task_intent: |
      {{task_intent}}
    activate_command: "/bmad:bmm:agents:{{base_agent_cmd}}"
    {{else}}
    base_agent: none
    custom_name: none
    task_intent: "No task specified - clean fork"
    {{end if}}

  notes: |
    Nook created from {{current_branch}}.
    {{if custom_agent_configured}}
    Custom agent "{{custom_name}}" configured for: {{task_intent_short}}
    Activate with: /bmad:bmm:agents:{{base_agent_cmd}}
    {{else}}
    Nook starts clean - no artifacts inherited.
    {{end if}}
    Run /bmad:th:workflows:nook-restore to load artifacts from parent.
    Run /bmad:th:workflows:nook-sync after making changes to save.
  ```

  <action>Write the context.yaml file</action>
</step>

<step n="10" goal="Report completion and next steps">
  <action>Calculate lineage depth</action>

  <action>Set workflow outputs (as defined in workflow.yaml):</action>
  ```
  outputs:
    new_branch: {{new_branch}}
    worktree_path: {{worktree_path}}
    custom_agent: {{custom_name}} (or "none" if wizard was skipped)
    lineage_context: {{lineage_context.brief_summary}} (or "none" if not gathered)
  ```

  <check if="{{custom_agent_configured}} is true">
    <action>Display completion summary WITH custom agent:</action>

    ```
    ╔═══════════════════════════════════════════════════════════════════╗
    ║              NOOK CREATED SUCCESSFULLY!                           ║
    ╠═══════════════════════════════════════════════════════════════════╣
    ║                                                                   ║
    ║  Parent Workspace                                                 ║
    ║     Branch: {{current_branch}}                                    ║
    ║     Hash:   {{parent_hash}}                                       ║
    ║     Path:   (current directory or base_path)                      ║
    ║                                                                   ║
    ║  New Nook                                                         ║
    ║     Branch:   {{new_branch}}                                      ║
    ║     Type:     {{nook_type}}                                       ║
    ║     Location: {{worktree_path}}                                   ║
    ║                                                                   ║
    ╠═══════════════════════════════════════════════════════════════════╣
    ║                    CUSTOM AGENT READY!                            ║
    ╠═══════════════════════════════════════════════════════════════════╣
    ║                                                                   ║
    ║  Your Specialist: "{{custom_name}}"                               ║
    ║  Base Agent:      {{base_agent_title}}                            ║
    ║  Task:            {{task_intent_short}}                           ║
    ║                                                                   ║
    ║  To activate your custom agent:                                   ║
    ║     /bmad:bmm:agents:{{base_agent_cmd}}                           ║
    ║                                                                   ║
    ╠═══════════════════════════════════════════════════════════════════╣
    ║  Lineage ({{lineage_count}} levels):                              ║
    ║     {{for each entry in lineage}}                                 ║
    ║     {{index}}. {{entry.branch}} ({{entry.hash}}) - {{entry.type}} ║
    ║     {{end for}}                                                   ║
    ║     -> {{new_branch}} (new) - {{nook_type}} <- YOU ARE HERE       ║
    ╚═══════════════════════════════════════════════════════════════════╝

    What's in your nook:
       Config files exist (.bmad/_cfg/agents/, .bmad/bmm/config.yaml)
         - Custom agent "{{custom_name}}" pre-configured
         - Skip-worktree applied (changes stay local)
       docs/ does NOT exist (gitignored)
         - Run nook-restore if you need docs

    Next Steps:

    1. Switch to the new worktree:
       cd {{worktree_path}}

    2. Activate your custom specialist:
       /bmad:bmm:agents:{{base_agent_cmd}}

    3. Your agent has pre-loaded context about:
       "{{task_intent_short}}"

    4. To load docs from parent context:
       /bmad:th:workflows:nook-restore

    5. When done with the nook:
       - Commit your code changes
       - Run /bmad:th:workflows:nook-sync to save artifacts
       - Merge back: git checkout {{current_branch}} && git merge {{new_branch}}
       - Remove worktree: git worktree remove "{{worktree_path}}"

    Start a new AI session in the nook and summon your specialist!
    ```
  </check>

  <check if="{{custom_agent_configured}} is NOT true">
    <action>Display completion summary WITHOUT custom agent:</action>

    ```
    Nook Created Successfully!

    Parent Workspace
       Branch: {{current_branch}}
       Hash:   {{parent_hash}}
       Path:   (current directory or base_path)

    New Nook
       Branch:   {{new_branch}}
       Type:     {{nook_type}}
       Location: {{worktree_path}}

    Lineage ({{lineage_count}} levels):
       {{for each entry in lineage}}
       {{index}}. {{entry.branch}} ({{entry.hash}}) - {{entry.type}}
       {{end for}}
       -> {{new_branch}} (new) - {{nook_type}} <- YOU ARE HERE

    What's in your nook:
       Config files exist (.bmad/_cfg/agents/, .bmad/bmm/config.yaml)
         - From git, with skip-worktree (local changes won't show in status)
       docs/ does NOT exist (gitignored)
         - Run nook-restore if you need docs
       .bmad-tracking/ does NOT exist (submodule removed)
       .gitmodules does NOT exist (removed)

    Next Steps:

    1. Switch to the new worktree:
       cd {{worktree_path}}

    2. Config files are ready to use!
       They exist from git with skip-worktree applied.
       You can modify them freely - changes stay local.

    3. To load docs from parent context:
       /bmad:th:workflows:nook-restore

    4. When done with the nook:
       - Commit your code changes
       - Run /bmad:th:workflows:nook-sync to save artifacts
       - Merge back: git checkout {{current_branch}} && git merge {{new_branch}}
       - Remove worktree: git worktree remove "{{worktree_path}}"

    5. To return to parent workspace:
       cd {{base_path}}

    Start a new AI session in the nook for maximum context clarity!
    ```
  </check>
</step>

</workflow>

## Branch Naming Reference

<notes>
**Type Folders:**
| Type | Purpose |
|------|---------|
| `discovery/` | New workspace/epic discovery |
| `explore/` | Deep investigation |
| `spike/` | Quick prototype/POC |
| `bugfix/` | Bug investigation & fix |
| `feature/` | Feature development |
| `experiment/` | Wild experimentation |
| `hotfix/` | Production hotfix |

**Why Hash-Based Naming:**
- Git refs are files in `.git/refs/heads/`
- A path cannot be both a file and directory
- `explore/` is a directory, `explore/a1b2-foo` is a file (branch)
- No conflicts, infinite depth via hash chaining

**Lineage Tracing:**
- 4-char hash prefix links to parent commit
- Full chain stored in context.yaml in base_workspace_path/.bmad-tracking/
- Can trace any branch back to its root

**Three-Tier Nook Architecture:**

| Type | Paths | Behavior in Nook |
|------|-------|------------------|
| **Skip-worktree** | `.bmad/_cfg/agents/`, `.bmad/bmm/config.yaml` | Exist from git, local changes ignored |
| **Submodule (skip-worktree + removed)** | `.bmad-tracking/`, `.gitmodules` | Skip-worktree applied, then removed |
| **Gitignored** | `docs/` | Don't exist, must restore |

**Why skip-worktree for configs:**
- Files exist immediately in nook (no restore needed)
- `git update-index --skip-worktree` marks files as "assume unchanged"
- Local changes don't show in git status
- Each nook can customize configs independently
- Changes are saved via nook-sync, not git commit

**Why submodule for tracking folder:**
- Submodules don't auto-initialize in worktrees - nooks get empty placeholder
- Skip-worktree on submodule pointer hides deletion from git status
- No sparse-checkout complexity needed
- `.gitmodules` also gets skip-worktree + removal for clean nooks
- All tracking goes through base workspace via absolute paths

**Why gitignore for docs:**
- Large directories would bloat git history
- User explicitly restores if needed
- Single source of truth in base workspace tracking folder

**Agent Wizard Context Integration:**
The Agent Wizard uses nook-context-analyst to gather broader context before creating custom agents:

1. **Pre-Creation Context Gathering** (step 5a-context):
   - Invokes nook-context-analyst agent before agent selection
   - Traces lineage and reads parent context.yaml summaries
   - Extracts inherited decisions, constraints, architecture patterns
   - Filters context relevant to the user's task intent

2. **Context Embedding in Memories**:
   - Task intent is the first memory entry
   - Inherited decisions prefixed with "INHERITED:"
   - Constraints prefixed with "CONSTRAINT:"
   - Parent findings prefixed with "PARENT CONTEXT:"
   - Full lineage chain included for traceability

3. **Context Analyst Availability**:
   - Every custom agent includes a critical action pointing to nook-context-analyst
   - Agents can invoke it for deeper big-picture understanding
   - Enables agents to self-orient when facing unfamiliar constraints

4. **Context Flow**:
   ```
   Parent Nooks → nook-sync → context.yaml (with summary)
                      ↓
   nook-fork → nook-context-analyst → lineage_context
                      ↓
   Agent Wizard → embeds in memories + critical_actions
                      ↓
   Custom Agent → pre-loaded with inherited context
                  + knows to call analyst for more
   ```

This ensures specialist agents understand not just their immediate task,
but also the broader project context they're operating within.
</notes>
