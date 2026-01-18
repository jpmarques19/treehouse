# Treehouse

> Fork easily. Stay focused. Sync what matters.

---

**Early Development** - Tested primarily with Claude Code + Opus 4.5 on Linux.

---

## Why Treehouse?

When projects grow, you need to explore multiple directions at once—investigating issues, trying new approaches, running experiments. Each exploration benefits from isolation, but creates a hidden cost: your documents, configurations, and decisions start drifting apart.

Treehouse solves this with **focused nook environments** (git worktrees), **persistent AI hats** for domain knowledge, and **boards** that capture learnings per nook.

**The key insight**: Each nook is isolated for focused work, but hats provide cross-cutting domain knowledge. Boards capture nook-specific learnings that persist across sessions.

## How It Works

```
Your Repository
│
├── .treehouse/                       ← Treehouse configuration
│   ├── hats/                         ← AI hats with project specific domain knowledge
│   │   ├── reviewer.md               ← Simple .md files
│   │   └── architect.md
│   ├── nooks/                        ← Isolated worktrees for focused work
│   │   ├── boards.yaml               ← Pins (learnings) per nook
│   │   ├── a1b2-feature-auth/
│   │   └── 7f3e-spike-perf/
│   └── workflows/                    ← Workflow definitions
│
├── .claude/commands/th/              ← Claude Code integration
│   ├── workflows/                    ← Workflow stubs
│   └── hats/                         ← Hat command stubs
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
- Claude Code workflow stubs

### 2. Create a Hat

```
/th:workflows:hat-add
```

Hats contain domain knowledge that applies across all nooks—coding standards, architecture decisions, project conventions.

### 3. Create a Nook for Focused Work

```
/th:workflows:nook-fork
```

Creates an isolated git worktree for focused exploration. Work freely without affecting your main branch.

### 4. Pin Learnings to the Board

```bash
th pin "Decided to use YAML for config files"
```

Or use the workflow:
```
/th:workflows:pin
```

Pins capture nook-specific learnings—decisions made, context discovered, notes for future sessions. View them with:

```bash
th board
```

### 5. Merge When Ready

When your work is complete, merge changes back and clean up the worktree.

## CLI Commands

| Command | Description |
|---------|-------------|
| `th init` | Initialize treehouse in a repository |
| `th fork <name>` | Create a new nook (worktree) |
| `th list` | List all nooks |
| `th pin <content>` | Save a learning to the current nook's board |
| `th board` | View all pins for the current nook |
| `th hat add <name> <config>` | Create a new hat |

## Claude Workflows

All treehouse operations also work through Claude workflows:

| Command | Description |
|---------|-------------|
| `/th:hats:{name}` | Load a hat |
| `/th:workflows:nook-fork` | Create a focused worktree |
| `/th:workflows:pin` | Save learnings to the board |
| `/th:workflows:treehouse-view` | View all nooks |
| `/th:workflows:hat-add` | Create a new hat |

## Hats System

Hats are lightweight domain knowledge containers:

- **Knowledge base** (`knowledge.md`) - Domain expertise that applies across all nooks
- Create specialized hats for different concerns—code review standards, architecture patterns, domain terminology

## Boards System

Boards capture nook-specific learnings:

- **Pins** - Timestamped notes saved to `boards.yaml`
- Each nook has its own section in the board
- Pins persist across sessions, helping you resume context
- Use `th pin` from within a nook to save learnings
- Use `th board` to view all pins for the current nook

```yaml
# .treehouse/nooks/boards.yaml
a1b2-feature-auth:
  - ts: "2026-01-18T10:30:00Z"
    content: |
      Decided to use JWT for auth tokens.
      Session duration: 24 hours.
```

## Path Resolution

Treehouse hardcodes the base workspace path into Claude stubs during `th init`. This ensures workflows work correctly from any nook without environment variable configuration.

## Why "Treehouse"?

Your repo is the tree. The `.treehouse/` is your elevated workspace built on top of it. **Nooks** are cozy corners for focused work. **Hats** provide domain expertise. **Boards** remember what you've learned.

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
