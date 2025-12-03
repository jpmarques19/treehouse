---
name: 'step-10-generate-summary'
description: 'Generate nook summary for context analyst'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-sync'

# File References
thisStepFile: '{workflow_path}/steps/step-10-generate-summary.md'
nextStepFile: '{workflow_path}/steps/step-11-build-lineage.md'
---

# Step 10: Generate Nook Summary

## STEP GOAL:

Generate an AI-analyzed summary of the nook's work for the nook-context-analyst.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- Read the complete step file before taking any action
- Analyze actual work done - NO placeholder text
- Generate meaningful, specific summary

### Role Reinforcement:

- You are a workspace automation assistant
- Summary quality is critical for context analyst
- Be thorough in analysis, accurate in output

## EXECUTION SEQUENCE:

### 1. Gather Information

<action>Get git commits (last 20 or since fork):</action>
```bash
git log --oneline -20 --no-merges
```

<action>Get list of files changed:</action>
```bash
git diff --name-only HEAD~10 2>/dev/null || git diff --name-only --cached
```

<action>If docs/ was synced, scan doc titles/headers:</action>
```bash
head -20 docs/*.md 2>/dev/null || true
```

<action>Check for existing summary in context (if updating):</action>

### 2. Analyze Gathered Information

<critical>
Generate a MEANINGFUL summary based on actual work done in this nook.
Do NOT use placeholder text. Analyze:
- Commit messages reveal PURPOSE and DECISIONS
- Changed files reveal SCOPE of work
- Doc content reveals INSIGHTS and FINDINGS
- Unfinished work or TODOs reveal NEXT STEPS
- Error mentions or "fix" commits may reveal BLOCKERS
</critical>

### 3. Construct Summary Object

<action>Build {{nook_summary}} with these fields:</action>

```yaml
summary:
  purpose: |
    [1-2 sentences: WHY this nook exists, derived from branch name + commits]

  insights:
    [List key findings/learnings from docs and commits, or empty if none]

  decisions:
    [List of {decision, rationale} pairs from commits/docs, or empty if none]

  status: |
    [Current state: what's done, what's in progress]

  blockers:
    [Known issues or blockers discovered, or empty if none]

  next_steps:
    [Recommended actions based on TODOs, incomplete work, or logical next steps]
```

### 4. Handle Minimal Content

<note>
If this is a fresh nook with minimal commits, generate a minimal but accurate summary.
It's better to have "purpose: Exploring X" with empty insights than fabricated content.
</note>

<check if="very few commits or minimal work">
  <action>Generate minimal but accurate summary:</action>
  ```yaml
  summary:
    purpose: |
      [Derive from branch name, e.g., "Exploring step-engine migration for TH workflows"]
    insights: []
    decisions: []
    status: |
      Fresh nook - work just beginning
    blockers: []
    next_steps:
      - Begin implementation
  ```
</check>

### 5. Store and Proceed

<action>Store {{nook_summary}} for manifest generation</action>
<action>Auto-proceed to next step</action>

## MENU OPTIONS:

This is an auto-proceed step. After generation:
- Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### SUCCESS:
- Git history analyzed
- Meaningful summary generated (not placeholder)
- All summary fields populated appropriately
- Empty arrays for genuinely empty sections
- Proceeded to step 11

### SYSTEM FAILURE:
- Using placeholder text
- Fabricating content not supported by evidence
- Not analyzing available information
- Skipping summary generation
