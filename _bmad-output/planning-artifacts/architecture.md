---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
status: 'complete'
completedAt: '2026-01-08'
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/product-brief-treehouse-v0.4.0-2026-01-07.md
  - _bmad-output/planning-artifacts/ux-design-specification.md
  - docs/index.md
  - docs/architecture.md
  - docs/development-guide.md
  - docs/project-overview.md
  - docs/source-tree-analysis.md
workflowType: 'architecture'
project_name: '8aa9-treehouse-planning'
user_name: 'Joao'
date: '2026-01-08'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

32 functional requirements across 6 domains:

| Domain | FRs | Key Capabilities |
|--------|-----|------------------|
| Workspace Init | FR1-4 | Initialize `.treehouse/`, create `decks.yaml`, validate git repo |
| Nook Management | FR5-10 | Create/list/remove nooks, auto-generate IDs, prune orphans |
| Deck & Lineage | FR11-15 | Auto-create decks, nest sub-forks, track parent-child, `dk-<id>` format |
| Crew Memory | FR16-22 | Long-term/short-term/session memory, nook detection, `/handoff` workflow |
| Crew Management | FR23-25 | Create crew from template, store in `.treehouse/crew/`, huddles |
| CLI Operations | FR26-32 | JSON output, exit codes, idempotent ops, single-command install |

**Non-Functional Requirements:**

| Category | Requirement | Target |
|----------|-------------|--------|
| Performance | CLI command response | < 1 second |
| Performance | Nook creation | < 5 seconds |
| Performance | Memory file operations | < 100ms |
| Reliability | CLI success rate | 100% (zero crashes) |
| Reliability | Workflow execution | Zero failures in core path |
| Reliability | Data integrity | No corruption of `decks.yaml` or memory files |
| Integration | Git compatibility | Standard git (2.5+ for worktrees) |
| Integration | Claude Code | JSON output parseable by workflows |

**Scale & Complexity:**

- Primary domain: CLI tool (Go) + Agentic workflows
- Complexity level: Low
- Estimated architectural components: ~5 (CLI, Deck Tracker, Crew Memory, Workflow Layer, File Storage)

### Technical Constraints & Dependencies

| Constraint | Source | Impact |
|------------|--------|--------|
| Go language for CLI | PRD decision | Enables < 1 second response, single binary distribution |
| JSON-only output | PRD requirement | No human-readable flags, workflows parse all output |
| Zero-config design | Product vision | No config file, sensible defaults only |
| Git worktree dependency | Core functionality | Requires git 2.5+, standard POSIX filesystem |
| Linux/macOS only (MVP) | Scope boundary | No Windows support initially |
| Single-user focus | Scope boundary | No team/multi-user features |

### Cross-Cutting Concerns Identified

| Concern | Affected Components | Resolution Approach |
|---------|---------------------|---------------------|
| Nook Detection | CLI, Crew, Workflows | Convention: detect from pwd against `decks.yaml` |
| Path Conventions | All | `.treehouse/` as root, predictable subdirectory structure |
| CLI-Workflow Contract | CLI, Workflows | JSON schema for all commands, documented error codes |
| Memory File Format | Crew, Workflows | Markdown files in predictable locations |
| Error Handling | CLI, Workflows | JSON error objects with code + message |

## Starter Template Evaluation

### Primary Technology Domain

Go CLI Tool with Agentic Workflow integration - not a typical web app starter scenario.

### Starter Options Considered

| Option | Decision |
|--------|----------|
| Go CLI Framework | Cobra (lightweight usage) |
| Project Structure | Standard Go layout with `cmd/` and `internal/` |
| Agentic Workflows | No starter - build on v0.3.x patterns |

### Selected Approach: Go + Cobra (Minimal)

**Rationale:**
- Cobra provides clean subcommand organization
- Skip features we don't need (Viper, auto-help)
- Custom JSON output for workflow consumption
- Standard Go project structure

### Development Setup

**Initialization:**

```bash
go mod init github.com/yourusername/th
go get -u github.com/spf13/cobra@latest
```

**Project Structure:**

```
th/
├── cmd/th/main.go          # Entry point, minimal code
├── internal/
│   ├── cmd/                # Cobra command definitions
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── fork.go
│   │   ├── list.go
│   │   ├── prune.go
│   │   └── remove.go
│   ├── deck/               # Deck/lineage logic
│   │   └── deck.go
│   ├── nook/               # Nook management
│   │   └── nook.go
│   └── output/             # JSON formatting
│       └── json.go
├── go.mod
└── go.sum
```

### End-User Distribution Strategy

**Installation Methods (per PRD FR30-32):**

| Method | Command | Target User |
|--------|---------|-------------|
| curl \| sh | `curl -sSL https://get.treehouse.dev/install.sh \| sh` | Most users |
| go install | `go install github.com/user/th@latest` | Go developers |
| Homebrew | `brew install th` (future) | macOS users |

**Cross-Compilation Targets:**

| OS | Architecture | Binary Name |
|----|--------------|-------------|
| Linux | amd64 | th-linux-amd64 |
| Linux | arm64 | th-linux-arm64 |
| macOS | amd64 (Intel) | th-darwin-amd64 |
| macOS | arm64 (Apple Silicon) | th-darwin-arm64 |

**Release Pipeline:**

- **Release Please** - Automated semantic versioning & changelog generation
- **GoReleaser** - Automates cross-compilation and GitHub releases
- **GitHub Actions** - CI/CD triggered on version tags
- **Install script** - Detects OS/arch, downloads correct binary, installs to PATH

**Release Flow:**
```
Conventional Commits → Release Please PR → Merge → GoReleaser → GitHub Release
```

**Release Artifacts:**

```
releases/v0.4.0/
├── th-linux-amd64
├── th-linux-arm64
├── th-darwin-amd64
├── th-darwin-arm64
├── checksums.txt
└── install.sh
```

### Architectural Decisions Established

**Language & Runtime:**
- Go 1.21+ with Go modules

**CLI Framework:**
- Cobra for subcommand structure
- No Viper config (zero-config design)
- Custom JSON-only output

**Testing:**
- Standard `go test`
- Table-driven test patterns

**Build & Distribution:**
- Release Please for semantic versioning & changelogs
- GoReleaser for cross-compilation
- GitHub Releases for distribution
- Single binary, zero runtime dependencies

**Note:** CLI initialization and release pipeline setup should be early implementation stories.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- `.treehouse/` folder structure
- `decks.yaml` schema
- Nook ID format
- Crew memory architecture
- CLI JSON output contract

**Important Decisions (Shape Architecture):**
- Exit code conventions
- Error code taxonomy
- Worktree path conventions

**Deferred Decisions (Post-MVP):**
- TUI dashboard architecture
- Agent edit workflow patterns
- Homebrew formula structure

### File Storage Architecture

**`.treehouse/` Folder Structure:**

```
.treehouse/
├── decks.yaml                        # Deck/nook lineage tracking
└── crew/                           # Agent definitions and memory
    └── {agent-name}/
        ├── {agent-name}.agent.yaml   # Agent definition
        ├── knowledge.md              # Long-term memory (global, cross-nook)
        ├── memories/                 # Short-term memory (per-nook)
        │   └── {nook-id}.md
        └── sessions/                 # Session state (per-nook)
            └── {nook-id}.md
```

**Rationale:** Centralized storage in `.treehouse/` ensures memory survives nook deletion and enables cross-nook context access. One file per nook keeps structure flat and simple.

### Nook Identification

**Format:** `{4-char-hash}-{user-name}`

**Examples:**
- `a1b2-auth-spike`
- `c3d4-jwt-experiment`
- `f7e8-refactor-deck-logic`

**Hash Source:** First 4 characters of the git commit SHA at fork time.

**Rationale:** Hash provides uniqueness and ties nook to its origin point. User-provided name gives human-readable context.

### Deck & Lineage Schema

**`decks.yaml` Structure:**

```yaml
decks:
  dk-a1b2:                        # Deck ID = dk-{first-nook-hash}
    created: 2026-01-08
    nooks:
      a1b2-auth-spike:
        parent: main
        created: 2026-01-08
      c3d4-jwt-variant:
        parent: a1b2-auth-spike
        created: 2026-01-09
```

**Conventions:**
- Deck ID: `dk-{first-nook-4-char-hash}`
- Worktree paths derived by convention: `{worktrees_folder}/{nook-id}/`
- No redundant path storage in YAML

### Crew Memory Architecture

**Three-Tier Memory Model:**

| Tier | File | Scope | Lifecycle |
|------|------|-------|-----------|
| Long-term | `knowledge.md` | Global (cross-nook) | Persistent, updated via learning |
| Short-term | `memories/{nook-id}.md` | Per-nook | Updated via `/handoff`, current work context |
| Session | `sessions/{nook-id}.md` | Per-nook | Updated via `/handoff`, for session restoration |

**Memory Location:** Centralized in `.treehouse/crew/{agent}/`

**Rationale:** Centralized storage ensures:
- Memory survives nook deletion
- Cross-nook access for long-term knowledge
- Single source of truth for crew state

### Crew Communication

**Huddle Workflow:**
- Purpose: Interactive inter-agent context sharing
- Status: Existing workflow needs review and update for v0.4.0
- Location: Will need to work with centralized memory in `.treehouse/crew/`

**Huddle Integration Points:**
- Reads from: `knowledge.md` (long-term), `memories/{nook-id}.md` (short-term)
- Writes to: Can update `knowledge.md` with shared insights
- Participants: Multiple crew members in same nook context

**Review Needed:**
- Verify huddle workflow compatibility with new `.treehouse/` structure
- Update file paths for centralized memory location
- Ensure huddle outputs can be persisted to appropriate memory tier

### CLI-Workflow Contract

**JSON Output Schema:**

**Success Response:**
```json
{
  "success": true,
  "data": {
    "nook_id": "a1b2-auth-spike",
    "deck_id": "dk-a1b2",
    "worktree": "/path/to/worktrees/a1b2-auth-spike"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "NOOK_NOT_FOUND",
    "message": "Nook 'a1b2-auth-spike' does not exist"
  }
}
```

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Git operation failed |
| 4 | Nook/deck not found |

**Error Code Categories:**
- `INIT_*` - Initialization errors (e.g., `INIT_NOT_GIT_REPO`)
- `NOOK_*` - Nook operation errors (e.g., `NOOK_NOT_FOUND`, `NOOK_ALREADY_EXISTS`)
- `DECK_*` - Deck tracking errors (e.g., `DECK_CORRUPT`)
- `GIT_*` - Git operation errors (e.g., `GIT_WORKTREE_FAILED`)

### Worktree Conventions

**Path Convention:** `{worktrees_folder}/{nook-id}/`

**Default worktrees_folder:** `../worktrees/` (sibling to base repo)

**Example:**
```
~/projects/
├── my-project/               # Base repo with .treehouse/
└── worktrees/
    ├── a1b2-auth-spike/      # Nook worktree
    └── c3d4-jwt-variant/     # Nook worktree
```

**Rationale:** Sibling folder keeps worktrees out of the main repo, avoiding gitignore complexity.

### Decision Impact Analysis

**Implementation Sequence:**
1. CLI scaffolding with Cobra
2. `decks.yaml` read/write operations
3. `th init` command
4. `th fork` command (nook creation + deck tracking)
5. `th list` command (JSON tree output)
6. `th remove` and `th prune` commands
7. Crew memory file operations
8. Agentic workflow integration (`/handoff`, `/th::{agent}`)
9. Huddle workflow review and update

**Cross-Component Dependencies:**

| Decision | Affects |
|----------|---------|
| Nook ID format | CLI, decks.yaml, memory file names, worktree paths |
| Memory location | Crew workflows, `/handoff`, agent activation, huddles |
| JSON schema | All CLI commands, all workflow parsing |
| Worktree convention | `th fork`, `th list`, nook detection |

## Implementation Patterns & Consistency Rules

### Agent Definition Schema (Intent-Based)

Agents are intent-based, not prescriptive. No menus, no rigid prompts - agents understand their role and act autonomously.

```yaml
agent:
  metadata:
    name: "AgentName"
    title: "Role Title"
    icon: "◇"

  persona:
    role: "Primary role description"

    identity: |
      Multi-line detailed expertise and domain knowledge.
      What this agent specializes in and how it approaches problems.

    communication_style: |
      How agent communicates with users.

    principles:
      - "Principle 1 - what this agent believes"
      - "Principle 2 - how this agent operates"

  critical_actions:
    # Nook detection
    - "Detect current nook from git worktree folder name"
    - "Construct paths: AGENT_FOLDER, MEMORY_FILE, SESSION_FILE"

    # State loading
    - "Load knowledge.md for global context"
    - "Load memories/{nook-id}.md for nook-specific context"
    - "Load sessions/{nook-id}.md for session restoration"

    # Context window awareness (Claude 4 best practice)
    - |
      Your context window will be automatically compacted as it approaches
      its limit, allowing you to continue working indefinitely from where
      you left off. Therefore, do not stop tasks early due to token budget
      concerns. As you approach your token budget limit, save your current
      progress and state to memory before the context window refreshes.
      Always be as persistent and autonomous as possible and complete tasks
      fully, even if the end of your budget is approaching. Never artificially
      stop any task early regardless of the context remaining.
```

### Memory File Patterns

**Long-term Memory (`knowledge.md`):**

```markdown
# {Agent Name} - Project Knowledge

> Global context accumulated across all nooks

---

## Project Understanding
{What the agent has learned about the project}

## Key Patterns
{Patterns the agent has identified}

## Technical Debt
{Items to address later}

## Key Files
{Important file paths and their purposes}
```

**Short-term Memory (`memories/{nook-id}.md`):**

```markdown
# {Agent Name} - {nook-id} Nook Memory

> Nook-specific context for this implementation branch

---

## YYYY-MM-DD: {Activity Title}

### Context
{Why this happened}

### Key Information
- **Item 1**: {value}
- **Item 2**: {value}

### Decisions Made
- {Decision 1}
- {Decision 2}
```

**Session State (`sessions/{nook-id}.md`):**

```markdown
# {Agent Name} - {nook-id} Session State

> Last Updated: YYYY-MM-DD

---

## Current State

### Session Focus
{What we're working on} - {STATUS}

### Recent Commits (this session)
```
{commit list}
```

### Build Status
✅ **BUILD PASSES** / ❌ **BUILD FAILS**

---

## Next Actions

### Priority 1: {Action}
- {Details}

---

## Resume Context

{Summary of where to pick up - what was happening, what's next}
```

### Naming Conventions

**Files:**

| Type | Convention | Example |
|------|------------|---------|
| Agent definition | `{name}.agent.yaml` | `prism.agent.yaml` |
| Nook memory | `{nook-id}.md` | `a1b2-auth-spike.md` |
| Session state | `{nook-id}.md` | `a1b2-auth-spike.md` |
| Go source files | `snake_case.go` | `deck_tracker.go` |

**Go Code:**

| Type | Convention | Example |
|------|------------|---------|
| Exported functions | `PascalCase` | `CreateNook`, `GetDeckByID` |
| Unexported functions | `camelCase` | `parseYAML`, `validateInput` |
| Variables | `camelCase` | `nookID`, `deckPath` |
| Constants | `PascalCase` | `DefaultWorktreesFolder` |

**YAML Fields:**

| Convention | Example |
|------------|---------|
| `snake_case` | `nook_id`, `created_at`, `parent_nook` |

### Date/Time Formats

| Context | Format | Example |
|---------|--------|---------|
| YAML fields | `YYYY-MM-DD` | `2026-01-08` |
| JSON output | ISO 8601 | `2026-01-08T10:30:00Z` |
| Memory files | `YYYY-MM-DD` | `## 2026-01-08: Activity` |
| Session header | `YYYY-MM-DD` | `> Last Updated: 2026-01-08` |

### Workflow File Patterns

**Treehouse v0.4.0 uses standalone workflow format** - single `.md` files without step-file architecture.

| Location | Purpose |
|----------|---------|
| `{dev-nook}/workflows/{name}.md` | Development (in planning nook) |
| `.treehouse/workflows/{name}.md` | Installation (in target repo) |

**Workflow Structure:**

```markdown
# {Workflow Name}

{Brief description of what the workflow does}

## Execution

### Step 1: {Action}
{Instructions for the agent}

### Step 2: {Action}
{Instructions for the agent}

## JSON Response Structure
{Document expected CLI output format}

## Error Handling
{Document error codes and user-friendly messages}
```

**Design Principles:**
- Single file per workflow (no step-file architecture)
- Agent interprets intent, not rigid step sequences
- Fully decoupled from BMAD workflow.xml engine
- Claude commands reference workflow files directly

**Workflow List (v0.4.0):**

| Workflow | Purpose | Epic |
|----------|---------|------|
| `treehouse-init.md` | Initialize workspace | Epic 1 |
| `nook-fork.md` | Create isolated nook | Epic 2 |
| `treehouse-list.md` | View/manage workspace tree | Epic 3 |
| `handoff.md` | Save agent context | Epic 5 |
| `crew-add.md` | Create new crew member | Epic 6 |
| `huddle.md` | Multi-agent context sharing | Epic 6 |

**Note:** Step-file architecture (`steps/step-{NN}-{name}.md`) is deprecated for Treehouse workflows

### Error Handling Patterns

**CLI errors (JSON):**

```json
{
  "success": false,
  "error": {
    "code": "CATEGORY_SPECIFIC_ERROR",
    "message": "Human-readable description",
    "details": { }
  }
}
```

**Workflow error handling:**
1. Check CLI exit code first
2. Parse JSON error for user-friendly message
3. Log raw JSON for debugging

### All AI Agents MUST

1. **Detect nook automatically** from working directory against `decks.yaml`
2. **Load context on activation** - knowledge.md, memories/{nook}.md, sessions/{nook}.md
3. **Save state before context refresh** - persist to memory/session files
4. **Complete tasks fully** - never stop early due to token concerns
5. **Be autonomous** - interpret intent, don't wait for prescriptive commands
6. **Use JSON for CLI interaction** - parse all `th` command output as JSON

### Anti-Patterns to Avoid

| Anti-Pattern | Correct Pattern |
|--------------|-----------------|
| Stopping early due to context limits | Save state, continue after refresh |
| Waiting for explicit commands | Interpret intent, act autonomously |
| Hardcoding paths | Derive from nook ID and conventions |
| Human-readable CLI output | JSON-only for workflow parsing |
| Prescriptive menus in agents | Intent-based persona |

## Project Structure & Boundaries

### Complete Project Directory Structure

Treehouse v0.4.0 consists of three distinct components:

**1. Go CLI (`th`) Repository:**

```
th/
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── .goreleaser.yaml              # Cross-compilation config
├── .github/
│   └── workflows/
│       ├── ci.yml                # Build + test on PR
│       └── release.yml           # GoReleaser on tag
├── cmd/
│   └── th/
│       └── main.go               # Entry point
├── internal/
│   ├── cmd/                      # Cobra command definitions
│   │   ├── root.go               # Root command, JSON output setup
│   │   ├── init.go               # th init
│   │   ├── fork.go               # th fork <name>
│   │   ├── list.go               # th list
│   │   ├── remove.go             # th remove <nook-id>
│   │   └── prune.go              # th prune
│   ├── deck/                     # Deck/lineage operations
│   │   ├── deck.go               # Deck struct, CRUD
│   │   ├── yaml.go               # decks.yaml read/write
│   │   └── deck_test.go
│   ├── nook/                     # Nook operations
│   │   ├── nook.go               # Nook struct, ID generation
│   │   ├── worktree.go           # Git worktree operations
│   │   └── nook_test.go
│   ├── output/                   # JSON output formatting
│   │   ├── json.go               # Success/error response builders
│   │   └── json_test.go
│   └── git/                      # Git operations wrapper
│       ├── git.go                # Git command execution
│       └── git_test.go
├── scripts/
│   └── install.sh                # curl | sh installer
└── testdata/                     # Test fixtures
    ├── valid_decks.yaml
    └── invalid_decks.yaml
```

**2. Treehouse System (`.treehouse/` in user's repo):**

```
{user-project}/
├── .treehouse/                   # Treehouse system root
│   ├── decks.yaml                # Deck/nook lineage tracking
│   ├── nooks/                    # Nook worktrees (git worktrees stored here)
│   │   └── {nook-id}/            # Individual nook worktree
│   ├── crew/                     # Agent definitions and memory
│   │   └── {agent-name}/
│   │       ├── {agent-name}.agent.yaml
│   │       ├── knowledge.md      # Long-term memory
│   │       ├── memories/         # Short-term (per-nook)
│   │       │   └── {nook-id}.md
│   │       └── sessions/         # Session state (per-nook)
│   │           └── {nook-id}.md
│   └── workflows/                # Treehouse workflows (standalone files)
│       ├── treehouse-init.md     # /treehouse-init
│       ├── nook-fork.md          # /nook-fork
│       ├── treehouse-list.md     # /treehouse-list
│       ├── handoff.md            # /handoff
│       ├── crew-add.md           # /crew-add
│       └── huddle.md             # /huddle
└── ...                           # User's project files
```

**Note:** Nook worktrees are stored inside `.treehouse/nooks/` (not as sibling folder).
Workflows are standalone single files (no folder wrappers or step-file architecture).

### User-Facing Commands (Claude Code Integration)

Users interact with Treehouse via Claude Code slash commands, not directly with the CLI.

**Command Architecture:**

Treehouse commands are registered as BMAD skills. Each skill points to a Claude command file that references the workflow.

```
.claude/commands/th/workflows/
├── treehouse-init.md            # /treehouse-init command
├── nook-fork.md                 # /nook-fork command
├── treehouse-list.md            # /treehouse-list command
└── ...
```

**Command File Format (Standalone):**

```markdown
# {Command Name}

{Brief description}

## Instructions

Execute the workflow by following the instructions in:
`workflows/{name}.md`

Read that workflow file and execute all steps.
```

**Command Mapping:**

| Skill | Workflow | CLI Called | Purpose |
|-------|----------|------------|---------|
| `th:workflows:treehouse-init` | `workflows/treehouse-init.md` | `th init` | Initialize workspace |
| `th:workflows:nook-fork` | `workflows/nook-fork.md` | `th fork <name>` | Create new nook |
| `th:workflows:treehouse-list` | `workflows/treehouse-list.md` | `th list` | Show deck/nook tree |
| `th:workflows:handoff` | `workflows/handoff.md` | None | Save agent context |
| `th:workflows:crew-add` | `workflows/crew-add.md` | None | Add crew member |
| `th:workflows:huddle` | `workflows/huddle.md` | None | Multi-agent sync |

**Development vs Installation:**

| Context | Workflow Location |
|---------|-------------------|
| Development (in nook) | `{nook}/workflows/` |
| Installation (target repo) | `.treehouse/workflows/` |

**Command Flow Example (`/nook-fork`):**

```
User: /nook-fork

┌─────────────────────────────────────────┐
│ workflows/nook-fork.md                  │
├─────────────────────────────────────────┤
│ 1. Detect current nook from pwd         │
│ 2. Ask user for nook name               │
│ 3. Call: th fork <name>                 │
│ 4. Parse JSON response                  │
│ 5. Display result with UX formatting:   │
│                                         │
│    ● Created nook: a1b2-auth-spike      │
│    ├── Deck: dk-a1b2                    │
│    ├── Parent: main                     │
│    └── Path: ../worktrees/a1b2-auth-... │
│                                         │
│ 6. Offer: cd to new worktree?           │
└─────────────────────────────────────────┘
```

**Command Flow Example (`/treehouse-list`):**

```
User: /treehouse-list

┌─────────────────────────────────────────┐
│ treehouse-list/workflow.md              │
├─────────────────────────────────────────┤
│ 1. Call: th list                        │
│ 2. Parse JSON response                  │
│ 3. Render tree with Unicode:            │
│                                         │
│    Treehouse: my-project                │
│    ─────────────────────                │
│    dk-a1b2                              │
│    ├── ● a1b2-auth-spike (← you)        │
│    │   └── ○ c3d4-jwt-variant           │
│    └── ○ e5f6-redis-cache               │
│                                         │
│    dk-g7h8                              │
│    └── ○ g7h8-refactor-api              │
│                                         │
│    Legend: ● current  ○ other           │
└─────────────────────────────────────────┘
```

### Architectural Boundaries

**CLI Boundaries:**

| Boundary | In | Out |
|----------|----|----|
| Input | Command args, pwd detection | No stdin interaction |
| Output | JSON to stdout | No human-readable text |
| Exit codes | 0-4 as defined | No other codes |
| Filesystem | `.treehouse/`, `../worktrees/` | No other locations |

**Workflow Boundaries:**

| Boundary | In | Out |
|----------|----|----|
| CLI interaction | Parse JSON output | Call `th` commands |
| Memory access | Read/write `.treehouse/crew/` | No direct decks.yaml writes |
| User interaction | Claude Code conversation | No external APIs |
| Display | Unicode tree rendering | UX spec formatting |

**Data Boundaries:**

| Data | Owner | Readers |
|------|-------|---------|
| `decks.yaml` | CLI (`th`) | Workflows (read-only) |
| `knowledge.md` | Agent workflows | Other agents (huddle) |
| `memories/*.md` | Agent workflows | Same agent only |
| `sessions/*.md` | Agent workflows | Same agent only |

### Requirements to Structure Mapping

**FR1-4 (Workspace Init) →**
- CLI: `internal/cmd/init.go`
- Workflow: `.treehouse/workflows/treehouse-init/`
- Creates: `.treehouse/decks.yaml`

**FR5-10 (Nook Management) →**
- CLI: `internal/cmd/fork.go`, `remove.go`, `prune.go`, `list.go`
- CLI: `internal/nook/nook.go`, `worktree.go`
- Workflow: `.treehouse/workflows/nook-fork/`
- Workflow: `.treehouse/workflows/treehouse-list/`

**FR11-15 (Deck & Lineage) →**
- CLI: `internal/deck/deck.go`, `yaml.go`
- Manages: `.treehouse/decks.yaml`

**FR16-22 (Crew Memory) →**
- Workflow: `.treehouse/workflows/handoff/`
- Storage: `.treehouse/crew/{agent}/`
- Not in CLI - workflow responsibility

**FR23-25 (Crew Management) →**
- Workflow: `.treehouse/workflows/crew-add/`
- Workflow: `.treehouse/workflows/huddle/`
- Storage: `.treehouse/crew/{agent}/`
- Not in CLI - workflow responsibility

**FR26-32 (CLI Operations) →**
- CLI: `internal/output/json.go`
- CLI: `internal/cmd/root.go`
- Release: `.goreleaser.yaml`
- Install: `scripts/install.sh`

### Integration Points

**CLI ↔ Workflows:**

```
┌─────────────┐         JSON stdout          ┌─────────────────┐
│   th CLI    │ ──────────────────────────▶  │ Agentic Workflow│
│             │                              │                 │
│  - init     │ ◀────────────────────────── │  /treehouse-init│
│  - fork     │    Parse & render UX         │  /nook-fork     │
│  - list     │                              │  /treehouse-list│
│  - remove   │                              │  /handoff       │
│  - prune    │                              │  /huddle        │
└─────────────┘                              └─────────────────┘
```

**Agent ↔ Memory:**

```
┌─────────────┐    load on activation    ┌─────────────────┐
│   Agent     │ ◀─────────────────────── │  knowledge.md   │
│  (workflow) │                          │  (global)       │
│             │ ◀─────────────────────── │  memories/*.md  │
│             │    load nook context     │  (per-nook)     │
│             │                          │                 │
│             │ ──────────────────────▶  │  sessions/*.md  │
│             │    save before refresh   │  (per-nook)     │
└─────────────┘                          └─────────────────┘
```

### File Organization Patterns

**Go CLI Source Organization:**
- `cmd/` - Entry point only, minimal code
- `internal/cmd/` - Cobra command implementations
- `internal/{domain}/` - Domain logic (deck, nook, git, output)
- Tests co-located: `*_test.go` alongside source

**Workflow Organization:**
- `workflow.md` - Entry point (single file for simple workflows)
- `steps/step-{NN}-{name}.md` - Sequential steps (for complex workflows)

**Memory File Organization:**
- Agent definition: `{name}.agent.yaml`
- Global memory: `knowledge.md`
- Per-nook: `memories/{nook-id}.md`, `sessions/{nook-id}.md`

### Development Workflow Integration

**Go CLI Development:**
```bash
# Development
go run ./cmd/th init
go run ./cmd/th fork my-feature

# Testing
go test ./...

# Build
go build -o th ./cmd/th

# Release (via CI)
git tag v0.4.0 && git push --tags
# GoReleaser handles cross-compilation
```

**Workflow Development:**
- Edit `.treehouse/workflows/` in base repo
- Changes propagate to worktrees via git
- Test by invoking `/nook-fork`, `/treehouse-list` in Claude Code

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:**
All technology choices (Go, Cobra, JSON output) work together without conflicts. The CLI and workflow boundaries are clean and non-overlapping.

**Pattern Consistency:**
Naming conventions are consistent across Go (snake_case files, PascalCase exports), YAML (snake_case fields), and all date formats use ISO 8601.

**Structure Alignment:**
The project structure fully supports all architectural decisions. CLI in `th/`, workflows in `.treehouse/workflows/`, memory in `.treehouse/crew/`.

### Requirements Coverage Validation ✅

**Functional Requirements Coverage:**

| FR Group | Coverage | Implementation |
|----------|----------|----------------|
| FR1-4 (Workspace Init) | ✅ | `th init` + `/treehouse-init` workflow |
| FR5-10 (Nook Management) | ✅ | `th fork/list/remove/prune` + `/nook-fork`, `/treehouse-list` |
| FR11-15 (Deck & Lineage) | ✅ | `internal/deck/` package + decks.yaml schema |
| FR16-22 (Crew Memory) | ✅ | Three-tier memory model + `/handoff` workflow |
| FR23-25 (Crew Management) | ✅ | `/crew-add` + `/huddle` workflows |
| FR26-32 (CLI Operations) | ✅ | JSON output, exit codes, GoReleaser distribution |

**Non-Functional Requirements Coverage:**

| NFR | Architectural Support | Status |
|-----|----------------------|--------|
| CLI < 1 second | Go binary, no heavy deps | ✅ |
| Nook creation < 5 seconds | Git worktree (native) | ✅ |
| Memory ops < 100ms | File-based, local | ✅ |
| Zero config | No config file in design | ✅ |
| Git 2.5+ compatibility | Standard worktree API | ✅ |

### Implementation Readiness Validation ✅

**Decision Completeness:**
All critical decisions documented with versions. JSON schemas, exit codes, and error taxonomy fully specified.

**Structure Completeness:**
Complete directory trees for CLI (`th/`) and treehouse system (`.treehouse/`). All files and integration points defined.

**Pattern Completeness:**
Agent schema, memory templates, session templates all provided. Claude 4 best practices (context window awareness) integrated into agent critical_actions.

### Gap Analysis Results

**Critical Gaps:** None

**Important Gaps:** None - all operations covered

`/treehouse-list` workflow handles:
- List all decks and nooks (via `th list`)
- Remove nook with cascade (via `th remove`)
- Prune orphaned worktrees (via `th prune`)

**Nook Remove Cascade Behavior:**

When a nook is removed via `th remove`:
1. Delete the nook's git worktree
2. Delete all child nooks (recursively)
3. For each deleted nook, clean up crew memory:
   - `.treehouse/crew/{agent}/memories/{nook-id}.md` (all agents)
   - `.treehouse/crew/{agent}/sessions/{nook-id}.md` (all agents)
4. Update `decks.yaml` to remove nook entries
5. If deck becomes empty, remove deck entry

**Nice-to-Have (future enhancement):**
- Mermaid diagrams for visual documentation
- Example fixtures (decks.yaml samples)

### Architecture Completeness Checklist

**✅ Requirements Analysis**
- [x] Project context thoroughly analyzed
- [x] Scale and complexity assessed (Low)
- [x] Technical constraints identified (Go, JSON-only, zero-config)
- [x] Cross-cutting concerns mapped (nook detection, path conventions)

**✅ Architectural Decisions**
- [x] Critical decisions documented with versions
- [x] Technology stack fully specified (Go + Cobra + workflows)
- [x] Integration patterns defined (CLI ↔ Workflow via JSON)
- [x] Performance considerations addressed (< 1s, < 5s targets)

**✅ Implementation Patterns**
- [x] Naming conventions established (Go, YAML, files)
- [x] Structure patterns defined (memory tiers, workflow organization)
- [x] Communication patterns specified (JSON contract)
- [x] Process patterns documented (error handling, context refresh)

**✅ Project Structure**
- [x] Complete directory structure defined
- [x] Component boundaries established (CLI vs workflows)
- [x] Integration points mapped (CLI ↔ workflows diagram)
- [x] Requirements to structure mapping complete

### Architecture Readiness Assessment

**Overall Status:** READY FOR IMPLEMENTATION

**Confidence Level:** High

**Key Strengths:**
- Clean separation between CLI (deterministic) and workflows (agentic)
- Intent-based agent design aligned with Claude 4 best practices
- Comprehensive JSON contract for CLI-workflow communication
- Three-tier memory model with clear file organization
- Cascade delete ensures no orphaned memory files

**Areas for Future Enhancement:**
- TUI dashboard (deferred post-MVP)
- Windows support (deferred)
- Team/multi-user features (deferred)

### Implementation Handoff

**AI Agent Guidelines:**
- Follow all architectural decisions exactly as documented
- Use implementation patterns consistently across all components
- Respect project structure and boundaries
- Refer to this document for all architectural questions

**First Implementation Priority:**
1. Go CLI scaffolding with Cobra
2. `th init` command + `/treehouse-init` workflow
3. `th fork` command + `/nook-fork` workflow
4. `th list` command + `/treehouse-list` workflow (includes remove/prune)

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED ✅
**Total Steps Completed:** 8
**Date Completed:** 2026-01-08
**Document Location:** `_bmad-output/planning-artifacts/architecture.md`

### Final Architecture Deliverables

**Complete Architecture Document**
- All architectural decisions documented with specific versions
- Implementation patterns ensuring AI agent consistency
- Complete project structure with all files and directories
- Requirements to architecture mapping
- Validation confirming coherence and completeness

**Implementation Ready Foundation**
- 15+ architectural decisions made
- Intent-based agent design patterns
- 3 distinct components (Go CLI, Treehouse System, Worktrees)
- 32 functional requirements fully supported

**AI Agent Implementation Guide**
- Technology stack: Go 1.21+ with Cobra
- Claude 4 best practices integrated
- JSON contract for CLI-workflow communication
- Three-tier memory model

### Implementation Handoff

**For AI Agents:**
This architecture document is your complete guide for implementing Treehouse v0.4.0. Follow all decisions, patterns, and structures exactly as documented.

**First Implementation Priority:**
```bash
# 1. Create Go CLI repository
mkdir th && cd th
go mod init github.com/yourusername/th
go get -u github.com/spf13/cobra@latest

# 2. Scaffold Cobra commands
# 3. Implement th init + /treehouse-init workflow
```

**Development Sequence:**
1. Go CLI scaffolding with Cobra
2. `th init` command implementation
3. `th fork` command + `/nook-fork` workflow
4. `th list` command + `/treehouse-list` workflow
5. Crew memory system + `/handoff` workflow
6. `/crew-add` and `/huddle` workflows

### Quality Assurance Checklist

**✅ Architecture Coherence**
- [x] All decisions work together without conflicts
- [x] Technology choices are compatible
- [x] Patterns support the architectural decisions
- [x] Structure aligns with all choices

**✅ Requirements Coverage**
- [x] All 32 functional requirements are supported
- [x] All non-functional requirements are addressed
- [x] Cross-cutting concerns are handled
- [x] Integration points are defined

**✅ Implementation Readiness**
- [x] Decisions are specific and actionable
- [x] Patterns prevent agent conflicts
- [x] Structure is complete and unambiguous
- [x] Claude 4 best practices integrated

---

**Architecture Status:** READY FOR IMPLEMENTATION ✅

**Next Phase:** Begin implementation using the architectural decisions and patterns documented herein.

**Document Maintenance:** Update this architecture when major technical decisions are made during implementation.

