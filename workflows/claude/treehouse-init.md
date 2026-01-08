---
name: 'treehouse-init'
description: 'Initialize base workspace - establishes the official starting point from which all nooks branch, ensuring clean configuration isolation'
---

<workflow-loader>
1. Read `.claude/settings.json` in the current working directory
2. Get `TREEHOUSE_BASE_WORKSPACE` from `env` - this is the base treehouse installation path
3. If set, use that as `BASE_PATH`. If not set, use current working directory as `BASE_PATH`
4. LOAD and EXECUTE the workflow from `{BASE_PATH}/.treehouse/workflows/treehouse-init.md`
</workflow-loader>
