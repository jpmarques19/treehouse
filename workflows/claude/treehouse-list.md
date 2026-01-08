---
name: 'treehouse-list'
description: 'View workspace lineage tree - displays all workspaces with visual hierarchy, marks current location, and provides options for cleanup operations'
---

<workflow-loader>
1. Get `TREEHOUSE_BASE_WORKSPACE` from `env` - this is the base treehouse installation path
2. If set, use that as `BASE_PATH`. If not set, use current working directory as `BASE_PATH`
3. LOAD and EXECUTE the workflow from `{BASE_PATH}/.treehouse/workflows/treehouse-list.md`
</workflow-loader>
