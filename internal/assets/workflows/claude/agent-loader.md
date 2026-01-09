---
name: 'agent-loader'
description: 'Load and execute a Treehouse agent'
---

TREEHOUSE_BASE_WORKSPACE={{TREEHOUSE_BASE_WORKSPACE}}

<agent-loader>
1. Extract AGENT_NAME from the command invocation (e.g., "th:crew:spruce" → "spruce")
2. Use `TREEHOUSE_BASE_WORKSPACE` defined above as `BASE_PATH`
3. LOAD and EXECUTE the agent from `{BASE_PATH}/.treehouse/crew/{AGENT_NAME}/{AGENT_NAME}.agent.yaml`
</agent-loader>
