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

### 2.1 Generate Role

Based on the intent summary, generate a concise role statement (1-2 sentences) that defines:
- Expertise domain and knowledge areas
- Capabilities and specializations
- Functional definition (NOT personality)

### 2.2 Display & Review

```
ROLE:
{generated_role}

  [c] Continue
  [e] Edit
```

**If [e]:** Ask what to change, regenerate, redisplay
**If [c]:** Proceed to Step 3

---

## Step 3: Identity Development

### 3.1 Generate Identity

Based on the intent and role, generate an identity (2-3 sentences) that defines:
- Personality type and attitude
- Emotional intelligence and worldview
- Character traits (NOT job description)

### 3.2 Display & Review

```
IDENTITY:
{generated_identity}

  [c] Continue
  [e] Edit
```

**If [e]:** Ask what to change, regenerate, redisplay
**If [c]:** Proceed to Step 4

---

## Step 4: Communication Style

### 4.1 Generate Communication Style

Based on the role and identity, generate a communication style that defines:
- Tone and formality level
- Verbosity and linguistic preferences
- Voice characteristics

### 4.2 Display & Review

```
COMMUNICATION STYLE:
{generated_communication_style}

  [c] Continue
  [e] Edit
  [p] Use preset instead
```

**If [e]:** Ask what to change, regenerate, redisplay
**If [p]:** Show presets:
```
Presets:
  [1] Direct & Technical - Clear explanations, practical solutions
  [2] Warm & Supportive - Friendly, uses examples and analogies
  [3] Concise & Precise - Minimal verbosity, bullets over prose
```
Apply selected preset and redisplay
**If [c]:** Proceed to Step 5

---

## Step 5: Principles Crafting

### 5.1 Generate Principles

Based on the intent, role, identity, and communication style, generate 5-7 principles:

**Generation guidelines:**
- First principle: Expert activator (core mission statement derived from role)
- Principles 2-4: Decision framework values (derived from identity)
- Principles 5-6: Behavioral constraints (derived from communication style)
- Always include Treehouse defaults:
  - "Context window awareness - save state before refresh"
  - "Complete tasks fully regardless of context remaining"

### 5.2 Display & Review

```
PRINCIPLES:
{generated_principles_list}

  [c] Continue
  [e] Edit
```

**If [e]:** Ask what to change, regenerate, redisplay
**If [c]:** Proceed to Step 6

---

## Step 6: Finalize Identity (Name, Title, Icon)

### 6.1 Generate Suggestions

Based on the complete persona, generate:
- Name: lowercase, hyphen-separated, memorable
- Title: 2-4 words, professional
- Icon: single emoji representing expertise

### 6.2 Display & Review

```
  Name: {suggested_name}
  Title: {suggested_title}
  Icon: {suggested_icon}

  Command: /th:agents:{suggested_name}

  [c] Continue
  [e] Edit
  [r] Regenerate
```

**If [e]:** Ask which field to change, update, redisplay
**If [r]:** Generate new suggestions, redisplay
**If [c]:** Validate and proceed to Step 7

**Validation:** Name must be unique in `.treehouse/agents/`

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
