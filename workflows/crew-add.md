# Crew Add Workflow

Progressive agent creation through iterative persona development.

## Overview

This workflow guides you through creating a thoughtful, well-crafted treehouse agent by developing their persona section-by-section. Each section is generated from your intent and reviewed before proceeding.

**Philosophy:** Quality agents come from progressive refinement, not form-filling.

---

## Step 1: Initialize & Collect Basics

### 1.1 Verify Treehouse

Run `th list` to verify treehouse is initialized.

If error with `INIT_NOT_FOUND`:
- Show: `✗ Treehouse not initialized. Run /treehouse-init first`
- Exit workflow

### 1.2 Collect Basic Information

Prompt for essential details:

```
Create new crew member

Name (lowercase, no spaces): _______
Title (e.g., "Code Reviewer"): _______
Icon (single emoji, default ◇): _______
```

**Validation:**
- Name: lowercase, alphanumeric, convert spaces to hyphens
- Name: must not exist in `.treehouse/agents/`
- Icon: defaults to ◇ if empty

**Once collected, confirm:**
```
Creating: {icon} {Name} - {title}

[c] Continue to persona development
[q] Quit
```

---

## Step 2: Role Development

**Field Purpose:** Define WHAT the agent does - their expertise domain and capabilities.

### 2.1 Explain Role Field

Show user:
```
ROLE - What This Agent Does

The role field defines the agent's professional identity:
  • Expertise domain and knowledge areas
  • Capabilities and what they're expert in
  • Functional definition (NOT personality)

Quality Check: Can you describe their job without mentioning personality?

Example: "Expert Go developer specializing in CLI tools and testing patterns"
```

### 2.2 Gather Intent

Ask:
```
What is this agent's primary expertise and purpose?

Describe what they do, what they're expert in, and their core capabilities:
```

Wait for user response.

### 2.3 Generate Role

From user's description, generate a concise role statement (1-2 sentences):
- Focus on expertise and capabilities
- Keep functional, not personal
- Be specific about domain knowledge

### 2.4 Display & Review

```
ROLE:
{generated_role}

Quality Check:
  ✓ Functional definition (job description)
  ✓ No personality traits mentioned
  ✓ Clear expertise domain

Review this role:
  [c] Continue to identity
  [e] Edit role
```

**If [e]:** Ask what to change, regenerate, redisplay, re-prompt
**If [c]:** Save role, proceed to Step 3

---

## Step 3: Identity Development

**Field Purpose:** Define WHO the agent is - their character and personality.

### 3.1 Explain Identity Field

Show user:
```
IDENTITY - Who This Agent Is

The identity field defines the agent's character:
  • Personality type and attitude
  • Emotional intelligence and worldview
  • Character definition (NOT job description)

Quality Check: Can you describe their character without their job title?

Example: "Thoughtful and detail-oriented. Approaches problems systematically
with a focus on maintainability. Values clarity and believes good code tells a story."
```

### 3.2 Gather Intent

Ask:
```
What is this agent's personality and character?

Describe their attitude, emotional approach, and worldview:
```

Wait for user response.

### 3.3 Generate Identity

From user's description, generate identity (2-3 sentences):
- Focus on personality and character
- Keep personal, not functional
- Establish emotional intelligence

### 3.4 Display & Review

```
IDENTITY:
{generated_identity}

Quality Check:
  ✓ Personality and character traits
  ✓ No job description elements
  ✓ Distinct from role

Review this identity:
  [c] Continue to communication style
  [e] Edit identity
```

**If [e]:** Ask what to change, regenerate, redisplay, re-prompt
**If [c]:** Save identity, proceed to Step 4

---

## Step 4: Communication Style

**Field Purpose:** Define HOW the agent speaks - tone, voice, and language patterns.

### 4.1 Explain Communication Style

Show user:
```
COMMUNICATION STYLE - How This Agent Speaks

The communication style defines speech patterns:
  • Tone and formality level
  • Verbosity and linguistic preferences
  • Voice characteristics (NOT expertise or personality)

Quality Check: Could a voice actor use this for direction?

Example: "Direct and technical. Prefers code examples over lengthy explanations.
Thinks in terms of interfaces and contracts."
```

### 4.2 Offer Presets

Show common presets:
```
Communication Style Presets:

1. Direct & Technical
   "Direct and helpful. Focuses on clear explanations and practical solutions."

2. Warm & Supportive
   "Friendly and encouraging. Uses examples and analogies to clarify concepts."

3. Concise & Precise
   "Minimal verbosity. Prefers bullet points and code over prose."

4. Custom
   Describe your own communication style

Select preset [1-3] or [4] for custom:
```

### 4.3 Generate Communication Style

If preset selected: Use preset text
If custom: Generate from user description

### 4.4 Display & Review

```
COMMUNICATION STYLE:
{generated_communication_style}

Quality Check:
  ✓ Language patterns and tone
  ✓ No expertise or personality overlap
  ✓ Actionable voice direction

Review this communication style:
  [c] Continue to principles
  [e] Edit communication style
```

**If [e]:** Ask what to change, regenerate, redisplay, re-prompt
**If [c]:** Save communication style, proceed to Step 5

---

## Step 5: Principles Crafting

**Field Purpose:** Define WHY the agent acts - decision-making framework and values.

### 5.1 Explain Principles

Show user:
```
PRINCIPLES - Why This Agent Acts

Principles define the agent's decision-making framework:
  • First principle: Expert Activator (core mission)
  • Principles 2-5: Decision framework (values that guide choices)
  • Principles 6+: Behavioral constraints (operational boundaries)

Quality Check: Would following these principles produce the desired behavior?

Example:
  - "Context window awareness - save state before refresh"
  - "Complete tasks fully regardless of context remaining"
  - "Test-driven development - write tests first"
```

### 5.2 Gather Intent

Ask:
```
What are this agent's core values and operating principles?

Describe their decision-making framework and behavioral guidelines:
```

Wait for user response.

### 5.3 Generate Principles

From user's description, generate 5-7 principles:
- First principle: Expert activator (core mission statement)
- Middle principles: Decision framework values
- Later principles: Behavioral constraints

Include Treehouse defaults:
- "Context window awareness - save state before refresh"
- "Complete tasks fully regardless of context remaining"

### 5.4 Display & Review

```
PRINCIPLES:
{generated_principles_list}

Quality Check:
  ✓ First principle activates expertise
  ✓ Creates decision-making clarity
  ✓ Includes behavioral constraints

Review these principles:
  [c] Continue to build
  [e] Edit principles
```

**If [e]:** Ask what to change, regenerate, redisplay, re-prompt
**If [c]:** Save principles, proceed to Step 6

---

## Step 6: Build & Celebrate

### 6.1 Generate Complete YAML

Create `{name}.agent.yaml`:

```yaml
agent:
  metadata:
    name: "{Name}"
    title: "{title}"
    icon: "{icon}"

  persona:
    role: "{generated_role}"

    identity: |
      {generated_identity}

    communication_style: |
      {generated_communication_style}

    principles:
{generated_principles_as_yaml_list}

  critical_actions:
    - "Detect current nook from git worktree folder name"
    - "Construct paths: AGENT_FOLDER, MEMORY_FILE, SESSION_FILE"
    - "Load knowledge.md for global cross-nook context"
    - "Load memories/{nook-id}.md for nook-specific work context"
    - "Load sessions/{nook-id}.md for session restoration"
    - |
      Your context window will be automatically compacted as it approaches
      its limit. Save progress to memory before context refresh. Complete
      tasks fully regardless of context remaining.
```

### 6.2 Create Folder Structure

Create:
```
.treehouse/agents/{name}/
├── {name}.agent.yaml
├── knowledge.md
├── memories/
└── sessions/
```

### 6.3 Create Knowledge Template

Create `knowledge.md`:

```markdown
# {Name} - Long-term Knowledge

> Cross-nook persistent memory

---

## Learnings

(Add global learnings that apply across all nooks)

---

Last Updated: {current_date}
```

### 6.4 Display Success

```
✓ Created crew member: {name}

  {icon} {Name} - {title}

  Folder: .treehouse/agents/{name}/
  Files:
    - {name}.agent.yaml
    - knowledge.md
    - memories/
    - sessions/

Persona Summary:
  Role: {generated_role}
  Identity: {generated_identity_first_sentence}
  Communication: {generated_communication_first_sentence}
  Principles: {count} principles defined

Activate with: /th:agents:{name}
```

---

## Dependencies

- Treehouse initialized (`th init`)
- Write access to `.treehouse/agents/`

## Design Notes

**Single-File Architecture:**
This workflow implements progressive development within a single file, maintaining Treehouse's architectural principle while achieving iterative refinement.

**Simple Review Pattern:**
Uses [c] Continue / [e] Edit for simplicity, avoiding complex menu systems. Treehouse workflows are focused and minimal.

**Quality Through Iteration:**
Each persona field is developed independently, reviewed, and refined before proceeding. This ensures thoughtful agent design.
