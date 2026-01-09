package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
)

// crewExitCode stores the exit code for the crew command
var crewExitCode int

// crewCmd represents the crew command (parent for subcommands)
var crewCmd = &cobra.Command{
	Use:           "crew",
	Short:         "Manage treehouse crew members",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		output.PrintError("CREW_NO_SUBCOMMAND", "Subcommand required: th crew --add <name> '<config>'")
		crewExitCode = 2
		return nil
	},
}

// crewAddCmd represents the crew --add subcommand
var crewAddCmd = &cobra.Command{
	Use:           "add <name> <json-config>",
	Short:         "Create a new crew member from JSON config",
	Long:          "Create a new crew member with agent definition files and Claude command stub",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		crewExitCode = runCrewAdd(args[0], args[1])
		return nil
	},
}

func init() {
	crewCmd.AddCommand(crewAddCmd)
	rootCmd.AddCommand(crewCmd)
}

// CrewConfig represents the JSON config for creating a crew member
type CrewConfig struct {
	Name    string       `json:"name"`
	Title   string       `json:"title"`
	Icon    string       `json:"icon"`
	Persona CrewPersona  `json:"persona"`
}

// CrewPersona represents the persona section of the agent config
type CrewPersona struct {
	Role               string   `json:"role"`
	Identity           string   `json:"identity"`
	CommunicationStyle string   `json:"communication_style"`
	Principles         []string `json:"principles"`
}

// CrewAddResult contains the data returned on successful crew add
type CrewAddResult struct {
	Agent   string `json:"agent"`
	Folder  string `json:"folder"`
	Command string `json:"command"`
}

// runCrewAdd executes the crew add command logic and returns the exit code
func runCrewAdd(name string, configJSON string) int {
	// 1. Validate name format (lowercase, no special chars except hyphens)
	sanitizedName := sanitizeCrewName(name)
	if sanitizedName == "" {
		output.PrintError("CREW_INVALID_NAME", "Invalid crew name: must be lowercase alphanumeric with hyphens only")
		return 2
	}

	// 2. Parse JSON config
	var config CrewConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		output.PrintError("CREW_INVALID_CONFIG", fmt.Sprintf("Invalid JSON config: %s", err.Error()))
		return 2
	}

	// 3. Validate required fields
	if err := validateCrewConfig(&config); err != nil {
		output.PrintError("CREW_INVALID_CONFIG", err.Error())
		return 2
	}

	// Use sanitized name
	config.Name = sanitizedName

	// Set default icon if not provided
	if config.Icon == "" {
		config.Icon = "◇"
	}

	// 4. Find treehouse (validates we're in an initialized repo)
	thInfo, err := treehouse.FindTreehouse(".")
	if err != nil {
		thErr, ok := err.(*treehouse.TreehouseError)
		if ok {
			output.PrintError(thErr.Code, thErr.Message)
		} else {
			output.PrintError("INIT_NOT_FOUND", err.Error())
		}
		return 1
	}

	// 5. Check if agent already exists
	agentPath := filepath.Join(thInfo.TreehousePath, "agents", sanitizedName)
	if _, err := os.Stat(agentPath); err == nil {
		output.PrintError("CREW_ALREADY_EXISTS", fmt.Sprintf("Crew member '%s' already exists", sanitizedName))
		return 1
	}

	// 6. Create agent folder structure
	if err := createCrewStructure(thInfo.TreehousePath, sanitizedName, &config); err != nil {
		output.PrintError("CREW_CREATE_FAILED", err.Error())
		return 1
	}

	// 7. Create Claude command stub
	repoRoot := filepath.Dir(thInfo.TreehousePath) // .treehouse parent is repo root
	commandPath, err := createCrewCommandStub(repoRoot, sanitizedName, config.Title, config.Icon)
	if err != nil {
		// Cleanup agent folder on failure
		_ = os.RemoveAll(agentPath)
		output.PrintError("CREW_COMMAND_FAILED", err.Error())
		return 1
	}

	// 8. Return success
	output.PrintSuccess(CrewAddResult{
		Agent:   sanitizedName,
		Folder:  agentPath,
		Command: commandPath,
	})
	return 0
}

// sanitizeCrewName validates and sanitizes the crew name
func sanitizeCrewName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any character that isn't alphanumeric or hyphen
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Collapse multiple hyphens and trim
	sanitized := result.String()
	for strings.Contains(sanitized, "--") {
		sanitized = strings.ReplaceAll(sanitized, "--", "-")
	}
	sanitized = strings.Trim(sanitized, "-")

	return sanitized
}

// validateCrewConfig checks all required fields are present
func validateCrewConfig(config *CrewConfig) error {
	if config.Title == "" {
		return fmt.Errorf("missing required field: title")
	}
	if config.Persona.Role == "" {
		return fmt.Errorf("missing required field: persona.role")
	}
	if config.Persona.Identity == "" {
		return fmt.Errorf("missing required field: persona.identity")
	}
	if config.Persona.CommunicationStyle == "" {
		return fmt.Errorf("missing required field: persona.communication_style")
	}
	if len(config.Persona.Principles) == 0 {
		return fmt.Errorf("missing required field: persona.principles (must have at least one)")
	}
	return nil
}

// createCrewStructure creates the agent folder and files
func createCrewStructure(treehousePath, name string, config *CrewConfig) error {
	agentPath := filepath.Join(treehousePath, "agents", name)

	// Create directories
	dirs := []string{
		agentPath,
		filepath.Join(agentPath, "memories"),
		filepath.Join(agentPath, "sessions"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			_ = os.RemoveAll(agentPath)
			return err
		}
	}

	// Create agent YAML file
	if err := createAgentYAML(agentPath, name, config); err != nil {
		_ = os.RemoveAll(agentPath)
		return err
	}

	// Create knowledge.md
	if err := createKnowledgeMD(agentPath, config); err != nil {
		_ = os.RemoveAll(agentPath)
		return err
	}

	return nil
}

// createAgentYAML creates the {name}.agent.yaml file
func createAgentYAML(agentPath, name string, config *CrewConfig) error {
	// Build principles YAML
	var principlesYAML strings.Builder
	for _, p := range config.Persona.Principles {
		principlesYAML.WriteString(fmt.Sprintf("      - %q\n", p))
	}

	// Capitalize name for display
	displayName := strings.Title(name)

	agentYAML := fmt.Sprintf(`agent:
  metadata:
    name: %q
    title: %q
    icon: %q

  persona:
    role: %q

    identity: |
      %s

    communication_style: |
      %s

    principles:
%s
  critical_actions:
    - "Detect current nook from git worktree folder name"
    - "Construct paths: AGENT_FOLDER, MEMORY_FILE, SESSION_FILE"
    - "Load knowledge.md for global cross-nook context"
    - "Load memories/{nook-id}.md for nook-specific work context"
    - "Load sessions/{nook-id}.md for session restoration"
    - |
      Your context window will be automatically compacted as it approaches
      its limit, allowing you to continue working indefinitely from where
      you left off. Therefore, do not stop tasks early due to token budget
      concerns. As you approach your token budget limit, save your current
      progress and state to memory before the context window refreshes.
      Always be as persistent and autonomous as possible and complete tasks
      fully, even if the end of your budget is approaching. Never artificially
      stop any task early regardless of the context remaining.
`,
		displayName,
		config.Title,
		config.Icon,
		config.Persona.Role,
		indentMultiline(config.Persona.Identity, 6),
		indentMultiline(config.Persona.CommunicationStyle, 6),
		principlesYAML.String(),
	)

	agentFile := filepath.Join(agentPath, name+".agent.yaml")
	return os.WriteFile(agentFile, []byte(agentYAML), 0644)
}

// createKnowledgeMD creates the knowledge.md template
func createKnowledgeMD(agentPath string, config *CrewConfig) error {
	displayName := strings.Title(strings.ReplaceAll(config.Name, "-", " "))
	if config.Name == "" {
		displayName = "Agent"
	}

	knowledgeMD := fmt.Sprintf(`# %s - Long-term Knowledge

> Cross-nook persistent memory

---

## Learnings

(Add global learnings that apply across all nooks)

---

Last Updated: (auto-updated on save)
`, displayName)

	knowledgePath := filepath.Join(agentPath, "knowledge.md")
	return os.WriteFile(knowledgePath, []byte(knowledgeMD), 0644)
}

// createCrewCommandStub creates the Claude command stub
func createCrewCommandStub(repoRoot, name, title, icon string) (string, error) {
	// Create directory structure
	commandDir := filepath.Join(repoRoot, ".claude", "commands", "th", "crew")
	if err := os.MkdirAll(commandDir, 0755); err != nil {
		return "", err
	}

	displayName := strings.Title(strings.ReplaceAll(name, "-", " "))

	commandMD := fmt.Sprintf(`# %s %s

%s

Load and activate the %s agent.

## Activation

`+"```"+`
/th:workflows:agent-loader %s
`+"```"+`
`, icon, displayName, title, displayName, name)

	commandPath := filepath.Join(commandDir, name+".md")
	if err := os.WriteFile(commandPath, []byte(commandMD), 0644); err != nil {
		return "", err
	}

	return commandPath, nil
}

// indentMultiline indents each line of a multiline string
func indentMultiline(s string, spaces int) string {
	indent := strings.Repeat(" ", spaces)
	lines := strings.Split(strings.TrimSpace(s), "\n")

	// First line doesn't need indent (it comes after the |)
	if len(lines) == 0 {
		return ""
	}
	if len(lines) == 1 {
		return lines[0]
	}

	var result strings.Builder
	result.WriteString(lines[0])
	for i := 1; i < len(lines); i++ {
		result.WriteString("\n")
		result.WriteString(indent)
		result.WriteString(lines[i])
	}
	return result.String()
}
