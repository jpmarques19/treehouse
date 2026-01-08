---
name: 'nook-fork'
description: 'Create a new nook for focused development - creates isolated worktree, syncs current context, and initializes clean-slate workspace branch'
---

<workflow-loader>
1. Read `.claude/settings.json` in the current working directory
2. Get `TREEHOUSE_BASE_WORKSPACE` from `env` - this is the base treehouse installation path
3. If set, use that as `BASE_PATH`. If not set, use current working directory as `BASE_PATH`
4. LOAD and EXECUTE the workflow from `{BASE_PATH}/.treehouse/workflows/nook-fork.md`
</workflow-loader>
