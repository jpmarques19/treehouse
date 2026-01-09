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

### 7.1 Compile JSON Config

Compile all session data into JSON config:

```json
{
  "name": "{name}",
  "title": "{title}",
  "icon": "{icon}",
  "persona": {
    "role": "{generated_role}",
    "identity": "{generated_identity}",
    "communication_style": "{generated_communication_style}",
    "principles": ["{principle_1}", "{principle_2}", ...]
  }
}
```

### 7.2 Create Agent via CLI

Run: `th crew --add {name} '{json_config}'`

Parse JSON response.

**If error `CREW_ALREADY_EXISTS`:**
Show `✗ Crew member '{name}' already exists` and offer to return to Step 6 to choose different name.

**If error `CREW_INVALID_CONFIG`:**
Show error details and offer to retry.

**If success:** Proceed to celebration.

### 7.3 Display Success

```
✓ Created crew member: {name}

  {icon} {Name} - {title}

  Files created:
    .treehouse/agents/{name}/
    .treehouse/.claude/commands/th/crew/{name}.md

Activate with: /th:crew:{name}
```

---

## Dependencies

- Treehouse initialized (`th init`)
- `th crew --add` CLI command available

## Design Notes

**Single-File Architecture:**
This workflow implements progressive development within a single file, maintaining Treehouse's architectural principle while achieving iterative refinement.

**CLI Delegation:**
File creation is delegated to `th crew --add` command, following the pattern of other workflows (treehouse-init → th init, nook-fork → th fork).

**Simple Review Pattern:**
Uses [c] Continue / [e] Edit for simplicity, avoiding complex menu systems. Treehouse workflows are focused and minimal.
