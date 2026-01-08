---
name: 'crew-add'
description: 'Create a new crew member from template'
---

<workflow-loader>
1. Get `TREEHOUSE_BASE_WORKSPACE` from `env` - this is the base treehouse installation path
2. If set, use that as `BASE_PATH`. If not set, use current working directory as `BASE_PATH`
3. LOAD and EXECUTE the workflow from `{BASE_PATH}/.treehouse/workflows/crew-add.md`
</workflow-loader>
