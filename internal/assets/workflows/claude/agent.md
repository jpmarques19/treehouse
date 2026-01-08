---
name: 'agent'
description: 'Load and execute a Treehouse agent'
---

<agent-loader>
1. Extract AGENT_NAME from the command invocation (e.g., "th:agents:spruce" → "spruce")
2. Get `TREEHOUSE_BASE_WORKSPACE` from `env` - this is the base treehouse installation path
3. If set, use that as `BASE_PATH`. If not set, use current working directory as `BASE_PATH`
4. LOAD and EXECUTE the agent from `{BASE_PATH}/.treehouse/agents/{AGENT_NAME}/{AGENT_NAME}.agent.yaml`
</agent-loader>
