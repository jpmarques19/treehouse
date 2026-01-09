---
name: 'treehouse-init'
description: 'Initialize base workspace - establishes the official starting point from which all nooks branch, ensuring clean configuration isolation'
---

TREEHOUSE_BASE_WORKSPACE={{TREEHOUSE_BASE_WORKSPACE}}

<workflow-loader>
1. Use `TREEHOUSE_BASE_WORKSPACE` defined above as `BASE_PATH`
2. LOAD and EXECUTE the workflow from `{BASE_PATH}/.treehouse/workflows/treehouse-init.md`
</workflow-loader>
