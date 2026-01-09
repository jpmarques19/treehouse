# Treehouse

> Fork easily. Stay focused. Sync what matters.

---

**Early Development** - Tested primarily with Claude Code + Opus 4.5 on Linux.

---

## Why Treehouse?

When projects grow, you need to explore multiple directions at once—investigating issues, trying new approaches, running experiments. Each exploration benefits from isolation, but creates a hidden cost: your documents, configurations, and decisions start drifting apart.

Treehouse solves this with **focused nook environments** (git worktrees) and **persistent AI crew members** that maintain context across sessions.

**The key insight**: Each nook is isolated for focused work, but crew members persist. Your AI assistants remember decisions, accumulate knowledge, and stay consistent regardless of which nook you're working in.

## How It Works

```
Your Repository
│
├── .treehouse/                       ← Treehouse configuration
│   ├── crew/                         ← AI crew members with persistent memory
│   │   └── oak/
│   │       ├── oak.agent.yaml        ← Agent definition
│   │       ├── knowledge.md          ← Accumulated knowledge
│   │       ├── memories/             ← Session memories
│   │       └── sessions/             ← Session logs
│   ├── nooks/                        ← Isolated worktrees for focused work
│   │   ├── a1b2-feature-auth/
│   │   └── 7f3e-spike-perf/
│   └── workflows/                    ← Workflow definitions
│
├── .claude/commands/th/              ← Claude Code integration
│   ├── workflows/                    ← Workflow stubs
│   └── crew/                         ← Crew member stubs
│
└── your-code/
```

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/jpmarques19/treehouse/main/install.sh | bash
```

## Quick Start

### 1. Initialize Treehouse

```bash
th init
```

This bootstraps your repository with:
- `.treehouse/` folder structure
- Default Oak assistant (your first crew member)
- Claude Code workflow stubs

### 2. Load Your Assistant

```
/th:crew:oak
```

Oak is ready to help. Crew members persist across sessions and accumulate knowledge as you work.

### 3. Create a Nook for Focused Work

```
/th:workflows:nook-fork
```

Creates an isolated git worktree for focused exploration. Work freely without affecting your main branch.

### 4. Work, Checkpoint, Merge

```
/th:workflows:checkpoint    # Save context before switching
```

When ready, merge changes back and clean up the worktree.

## Claude Workflows

All treehouse operations happen through Claude workflows:

| Command | Description |
|---------|-------------|
| `/th:crew:{name}` | Load a crew member |
| `/th:workflows:nook-fork` | Create a focused worktree |
| `/th:workflows:checkpoint` | Save current context |
| `/th:workflows:treehouse-view` | View all nooks |
| `/th:workflows:crew-add` | Create a new crew member |
| `/th:workflows:huddle` | Multi-agent discussion |

## Crew System

Crew members are AI agents with persistent identity and memory:

- **Agent definition** (`{name}.agent.yaml`) - Role, persona, principles
- **Knowledge base** (`knowledge.md`) - Accumulated context
- **Memories** - Session-to-session continuity
- **Sessions** - Interaction logs

Create specialized crew members for different tasks—code reviewers, architects, domain experts—each maintaining their own perspective and knowledge.

## Path Resolution

Treehouse hardcodes the base workspace path into Claude stubs during `th init`. This ensures workflows work correctly from any nook without environment variable configuration.

## Why "Treehouse"?

Your repo is the tree. The `.treehouse/` is your elevated workspace built on top of it. **Nooks** are cozy corners for focused work. **Crew** members live there with you, remembering what matters.

---

## Development

For contributors building from source:

```bash
git clone https://github.com/jpmarques19/treehouse.git
cd treehouse
go build -o th ./cmd/th
sudo mv th /usr/local/bin/
```

---

Feedback welcome!
