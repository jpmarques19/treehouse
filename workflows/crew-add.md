# Crew Add Workflow

Progressive agent creation through iterative persona development.

## Overview

This workflow guides you through creating a thoughtful, well-crafted treehouse agent through intent-driven discovery. We'll progressively build the persona from your needs, then derive the agent's identity naturally from what we've created.

**Philosophy:** Quality agents emerge from understanding intent, not filling forms.

**Flow:**
1. **Discover Intent** - What problem are you solving?
2. **Build Persona** - Role → Identity → Communication → Principles
3. **Derive Identity** - Name, title, icon emerge from the persona
4. **Create Agent** - Generate files and celebrate

---

## Step 1: Intent Discovery

### 1.1 Verify Treehouse

Run `th list` to verify treehouse is initialized.

If error with `INIT_NOT_FOUND`:
- Show: `✗ Treehouse not initialized. Run /treehouse-init first`
- Exit workflow

### 1.2 Discover User Intent

Show:
```
Create New Crew Member

Let's discover what agent you need through conversation.

What challenge or task do you want this agent to help with?

Describe the problem you're trying to solve, the work you need help with,
or the expertise you're looking for:
```

Wait for user response.

### 1.3 Clarify & Expand Intent

Based on user's initial response, ask 2-3 clarifying questions to understand:
- **Context:** Where/when will this agent be used?
- **Scope:** Narrow specialist or broad generalist?
- **Value:** What makes this agent valuable vs. doing it yourself?

Synthesize the conversation into a clear intent statement:
```
Intent Summary:
{synthesized_intent_paragraph}

[c] Continue to build this agent
[r] Restart - I want a different agent
[q] Quit
```

**If [r]:** Return to 1.2
**If [c]:** Proceed to Step 2

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

### 2.2 Generate Role from Intent

Using the intent summary from Step 1, generate the role field that defines:
- Expertise domain and knowledge areas
- Capabilities and what they're expert in
- Functional definition (NOT personality)

**Generation Prompt (internal):**
"Based on this intent: {intent_summary}, generate a concise role statement that defines what this agent is expert in and what they do professionally."

### 2.3 Display & Review

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

## Step 6: Finalize Identity (Name, Title, Icon)

**Purpose:** Derive the agent's external identity from the persona we've built.

### 6.1 Generate Suggestions

Based on the complete persona (role + identity + communication + principles), generate suggestions:

**Generation Prompt (internal):**
"Based on this agent's complete persona, suggest:
1. A name (lowercase, hyphen-separated, memorable, reflects expertise)
2. A title (2-4 words, professional, clear role indication)
3. An icon (single emoji that represents their expertise domain)"

### 6.2 Present & Refine

```
Final Identity

Based on your agent's persona, here are suggestions:

  Name: {suggested_name}
  Title: {suggested_title}
  Icon: {suggested_icon}

The name will be used for: /th:agents:{suggested_name}

Options:
  [c] Use these suggestions
  [e] Edit name, title, or icon
  [r] Regenerate different suggestions
```

### 6.3 Collect Changes (if needed)

**If [e]:** Ask which field to change, get new value, redisplay
**If [r]:** Generate new suggestions, redisplay
**If [c]:** Validate and proceed to Step 7

**Validation:**
- Name: lowercase, alphanumeric with hyphens only
- Name: must not exist in `.treehouse/agents/`
- Icon: single emoji (default ◇ if invalid)
- Title: non-empty string

---

## Step 7: Build & Celebrate

### 7.1 Generate Complete YAML

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

### 7.2 Create Folder Structure

Create:
```
.treehouse/agents/{name}/
├── {name}.agent.yaml
├── knowledge.md
├── memories/
└── sessions/
```

### 7.3 Create Knowledge Template

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

### 7.4 Display Success

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
