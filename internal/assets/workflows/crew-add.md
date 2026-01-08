# Crew Add Workflow

Create a new crew member from template.

## Overview

This workflow guides the creation of a new treehouse agent with proper folder structure and YAML definition. It collects basic info, creates the agent files, and confirms success.

## Execution

### Step 1: Check Treehouse

1. Run `th list` to verify treehouse is initialized
2. If error with `INIT_NOT_FOUND`, show: `✗ Treehouse not initialized. Run /treehouse-init first` and exit

### Step 2: Collect Agent Info

Prompt for agent details:

```
Create new crew member

1. Name (lowercase, no spaces):
2. Title (e.g., "Code Reviewer"):
3. Icon (single emoji, default ◇):
4. Role (one sentence - what they do):
5. Identity (2-3 sentences - who they are, their approach):
```

Validate:
- Name must be lowercase, alphanumeric, no spaces (convert spaces to hyphens)
- Name must not already exist in `.treehouse/agents/`
- Icon defaults to ◇ if not provided

### Step 3: Create Folder Structure

Create agent folder:
```
.treehouse/agents/{name}/
├── {name}.agent.yaml
├── knowledge.md
├── memories/
└── sessions/
```

### Step 4: Generate Agent YAML

Create `{name}.agent.yaml` with treehouse agent schema:

```yaml
agent:
  metadata:
    name: "{Name}"
    title: "{title}"
    icon: "{icon}"

  persona:
    role: "{role}"

    identity: |
      {identity}

    communication_style: |
      Direct and helpful. Focuses on clear explanations and
      practical solutions.

    principles:
      - "Context window awareness - save state before refresh"
      - "Complete tasks fully regardless of context remaining"

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

### Step 5: Create Knowledge Template

Create `knowledge.md`:

```markdown
# {Name} - Long-term Knowledge

> Cross-nook persistent memory

---

## Learnings

(Add global learnings that apply across all nooks)

---

Last Updated: {date}
```

### Step 6: Confirm Creation

Display success:

```
✓ Created crew member: {name}

  {icon} {Name} - {title}

  Folder: .treehouse/agents/{name}/
  Files:
    - {name}.agent.yaml
    - knowledge.md
    - memories/
    - sessions/

Activate with: /th:agents:{name}
```

## Dependencies

- Treehouse must be initialized (`th init`)
- Write access to `.treehouse/agents/`
