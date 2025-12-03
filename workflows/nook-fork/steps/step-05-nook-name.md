---
name: 'step-05-nook-name'
description: 'Define the nook name and construct the full branch name'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-05-nook-name.md'
nextStepFile: '{workflow_path}/steps/step-06-agent-decision.md'

# Inputs (may be provided via workflow inputs)
# nook_name: optional - if provided, skip interactive input
---

# Step 5: Define Nook Name

## STEP GOAL:

Get the nook name from user and construct the full branch name using the hash-based lineage scheme.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- ⏸️ HALT and wait for user input (unless parameter provided)
- 🎯 Validate kebab-case naming

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Show preview of full branch name
- ✅ Validate naming conventions strictly

## EXECUTION SEQUENCE:

### 1. Check for Input Parameter

<check if="nook_name input was provided">
  <action>Validate nook_name is kebab-case (lowercase, hyphens only, no spaces/slashes)</action>
  <action>If valid: Store as {{nook_name}} and construct branch name</action>
  <action>If invalid: Display error and ask for input</action>
</check>

### 2. Interactive Input (if no parameter)

<check if="nook_name NOT provided">
  <action>Ask for name:</action>

```
What should this nook be called? (kebab-case, e.g., "race-condition", "mqtt-reconnect")

The full branch will be: {{nook_type}}/{{parent_hash}}-<your-name>
```

  <action>HALT and wait for user input</action>
</check>

### 3. Validate Name

<action>Validate the provided name:</action>
- Must be kebab-case (lowercase letters and hyphens only)
- No spaces or special characters
- No slashes
- Not empty

<check if="name is invalid">
  <action>Display error and ask again:</action>
  ```
  Invalid name. Please use kebab-case:
  - Lowercase letters and hyphens only
  - No spaces, slashes, or special characters
  - Example: "my-feature", "bug-fix", "quick-test"
  ```
</check>

### 4. Construct Branch Name

<action>Construct full branch name:</action>
```
{{nook_type}}/{{parent_hash}}-{{nook_name}}
```

<action>Store as {{new_branch}}</action>

### 5. Display Preview

```
Nook Preview

Parent branch:  {{current_branch}} ({{parent_hash}})
New branch:     {{new_branch}}
Nook type:      {{nook_type}}
```

## MENU OPTIONS:

After preview: Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Nook name collected (interactive or via input)
- Name validated as kebab-case
- Full branch name constructed
- Preview displayed to user
- {{nook_name}} and {{new_branch}} stored

### ❌ SYSTEM FAILURE:
- Accepting invalid name formats
- Not showing preview
- Constructing branch name incorrectly
