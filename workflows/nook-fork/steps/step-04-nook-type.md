---
name: 'step-04-nook-type'
description: 'Select the type of nook to create (explore, spike, bugfix, etc.)'

# Path Definitions
workflow_path: '{project-root}/.bmad/th/workflows/nook-fork'

# File References
thisStepFile: '{workflow_path}/steps/step-04-nook-type.md'
nextStepFile: '{workflow_path}/steps/step-05-nook-name.md'

# Inputs (may be provided via workflow inputs)
# nook_type: optional - if provided, skip interactive selection
---

# Step 4: Select Nook Type

## STEP GOAL:

Determine the type of nook to create, which becomes the branch folder prefix.

## MANDATORY EXECUTION RULES (READ FIRST):

### Universal Rules:

- 📖 CRITICAL: Read the complete step file before taking any action
- ⏸️ HALT and wait for user selection (unless input provided)
- 🎯 Validate input if provided via workflow parameters

### Role Reinforcement:

- ✅ You are a workspace automation assistant
- ✅ Present clear options with descriptions
- ✅ Handle both interactive and parameter-driven input

## EXECUTION SEQUENCE:

### 1. Check for Input Parameter

<check if="nook_type input was provided">
  <action>Validate nook_type is one of: explore, spike, bugfix, discovery, feature, experiment, hotfix</action>
  <action>If valid: Store as {{nook_type}} and proceed to next step</action>
  <action>If invalid: Display error and show interactive menu</action>
</check>

### 2. Interactive Selection (if no input)

<check if="nook_type NOT provided">
  <action>Display selection menu:</action>

```
What kind of nook do you want to create?

1. explore/    - Deep exploration or investigation
2. spike/      - Quick prototype or proof of concept
3. bugfix/     - Focused bug investigation and fix
4. discovery/  - Requirements or architecture discovery
5. feature/    - New feature branch
6. experiment/ - Wild experimentation
7. hotfix/     - Production hotfix
8. custom      - Custom branch prefix

Choose [1-8]:
```

  <action>HALT and wait for user input</action>
</check>

### 3. Process Selection

<action>Map choice to type:</action>
- 1 → explore
- 2 → spike
- 3 → bugfix
- 4 → discovery
- 5 → feature
- 6 → experiment
- 7 → hotfix
- 8 → Ask for custom type

<check if="user selects 8 (custom)">
  <action>Ask: "Enter custom branch prefix (lowercase, no slashes):"</action>
  <action>Validate: lowercase only, no slashes or special characters</action>
</check>

<action>Store as {{nook_type}}</action>
<action>Display: "Selected nook type: {{nook_type}}"</action>

## MENU OPTIONS:

After selection: Load and execute `{nextStepFile}`

---

## SUCCESS/FAILURE METRICS

### ✅ SUCCESS:
- Nook type selected (interactive or via input)
- Type validated against allowed values
- Custom type validated if used
- {{nook_type}} stored for branch naming

### ❌ SYSTEM FAILURE:
- Accepting invalid nook types
- Not validating custom prefix format
- Proceeding without user selection
