package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
)

// hatExitCode stores the exit code for the hat command
var hatExitCode int

// hatCmd represents the hat command (parent for subcommands)
var hatCmd = &cobra.Command{
	Use:           "hat",
	Short:         "Manage treehouse hats",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		output.PrintError("HAT_NO_SUBCOMMAND", "Subcommand required: th hat add <name> '<config>'")
		hatExitCode = 2
		return nil
	},
}

// hatAddCmd represents the hat add subcommand
var hatAddCmd = &cobra.Command{
	Use:           "add <name> <json-config>",
	Short:         "Create a new hat from JSON config",
	Long:          "Create a new hat as a simple .md file",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		hatExitCode = runHatAdd(args[0], args[1])
		return nil
	},
}

func init() {
	hatCmd.AddCommand(hatAddCmd)
	rootCmd.AddCommand(hatCmd)
}

// HatConfig represents the JSON config for creating a hat
type HatConfig struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

// HatAddResult contains the data returned on successful hat add
type HatAddResult struct {
	Hat     string `json:"hat"`
	File    string `json:"file"`
	Command string `json:"command"`
}

// runHatAdd executes the hat add command logic and returns the exit code
func runHatAdd(name string, configJSON string) int {
	// 1. Validate name format (lowercase, no special chars except hyphens)
	sanitizedName := sanitizeHatName(name)
	if sanitizedName == "" {
		output.PrintError("HAT_INVALID_NAME", "Invalid hat name: must be lowercase alphanumeric with hyphens only")
		return 2
	}

	// 2. Parse JSON config
	var config HatConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		output.PrintError("HAT_INVALID_CONFIG", fmt.Sprintf("Invalid JSON config: %s", err.Error()))
		return 2
	}

	// 3. Validate required fields
	if err := validateHatConfig(&config); err != nil {
		output.PrintError("HAT_INVALID_CONFIG", err.Error())
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

	// 5. Check if hat already exists (as .md file)
	hatPath := filepath.Join(thInfo.TreehousePath, "hats", sanitizedName+".md")
	if _, err := os.Stat(hatPath); err == nil {
		output.PrintError("HAT_ALREADY_EXISTS", fmt.Sprintf("Hat '%s' already exists", sanitizedName))
		return 1
	}

	// 6. Create hat .md file
	if err := createHatFile(thInfo.TreehousePath, sanitizedName, &config); err != nil {
		output.PrintError("HAT_CREATE_FAILED", err.Error())
		return 1
	}

	// 7. Create Claude command stub
	repoRoot := filepath.Dir(thInfo.TreehousePath) // .treehouse parent is repo root
	commandPath, err := createHatCommandStub(repoRoot, sanitizedName, config.Title, config.Icon)
	if err != nil {
		// Cleanup hat file on failure
		_ = os.Remove(hatPath)
		output.PrintError("HAT_COMMAND_FAILED", err.Error())
		return 1
	}

	// 8. Return success
	output.PrintSuccess(HatAddResult{
		Hat:     sanitizedName,
		File:    hatPath,
		Command: commandPath,
	})
	return 0
}

// sanitizeHatName validates and sanitizes the hat name
func sanitizeHatName(name string) string {
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

// validateHatConfig checks all required fields are present
func validateHatConfig(config *HatConfig) error {
	if config.Title == "" {
		return fmt.Errorf("missing required field: title")
	}
	return nil
}

// createHatFile creates the hat .md file directly in .treehouse/hats/
func createHatFile(treehousePath, name string, config *HatConfig) error {
	hatsDir := filepath.Join(treehousePath, "hats")

	// Ensure hats directory exists
	if err := os.MkdirAll(hatsDir, 0755); err != nil {
		return err
	}

	titleCaser := cases.Title(language.English)
	displayName := titleCaser.String(strings.ReplaceAll(config.Name, "-", " "))
	if config.Name == "" {
		displayName = "Hat"
	}

	hatMD := fmt.Sprintf(`# %s %s

> %s

---

## Domain Knowledge

(Add domain expertise that applies across all nooks)

---

Last Updated: (auto-updated on save)
`, config.Icon, displayName, config.Title)

	hatPath := filepath.Join(hatsDir, name+".md")
	return os.WriteFile(hatPath, []byte(hatMD), 0644)
}

// createHatCommandStub creates the Claude command stub with inline hat-loader
func createHatCommandStub(repoRoot, name, title, icon string) (string, error) {
	// Create directory structure
	commandDir := filepath.Join(repoRoot, ".claude", "commands", "th", "hats")
	if err := os.MkdirAll(commandDir, 0755); err != nil {
		return "", err
	}

	titleCaser := cases.Title(language.English)
	displayName := titleCaser.String(strings.ReplaceAll(name, "-", " "))

	commandMD := fmt.Sprintf(`# %s %s

%s

TREEHOUSE_BASE_WORKSPACE=%s

<hat-loader>
1. Load {BASE_PATH}/.treehouse/hats/%s.md for domain knowledge
</hat-loader>
`, icon, displayName, title, repoRoot, name)

	commandPath := filepath.Join(commandDir, name+".md")
	if err := os.WriteFile(commandPath, []byte(commandMD), 0644); err != nil {
		return "", err
	}

	return commandPath, nil
}
