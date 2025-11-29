[![BMAD Extension](https://img.shields.io/badge/BMAD-Extension-blue)](https://github.com/bmad-code-org/BMAD-METHOD)

# 🌳 Treehouse (BMad Workspace) v0.1

>Fork easily. Stay focused. Sync what matters.

Ever had one document say one thing, another artifact evolve differently, and references pointing to a hybrid that exists nowhere?

That's **context drift** — and it gets worse as you explore multiple directions in parallel.

Treehouse enables multiple work environments (nooks) while keeping your **contextual artifacts** synchronized through a central tracking hub.

```
Base Workspace
└── .bmad-tracking/          ← single source of truth
    ├── main/docs/
    ├── spike/new-approach/docs/
    └── explore/research-thread/docs/
```

## Workflows

- `treehouse-init` → Set up base workspace
- `nook-fork` → Create focused, clean context environment
- `nook-sync` → Save context to tracking
- `nook-restore` → Load context from tracking
- `treehouse-list` → View workspace lineage tree

---

⚠️ **EXPERIMENTAL** — Very much WIP, narrowly tested with Claude Code + Opus 4.5.
Tightly coupled with git internals (worktrees, skip-worktree). Use at your own risk!

📄 README: `.bmad/th/README.md`

Feedback welcome! 🙏

---

Built for [BMAD-METHOD](https://github.com/bmad-code-org/BMAD-METHOD)
