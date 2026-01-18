package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/assets"
	"github.com/jpmarques19/treehouse/internal/git"
	"github.com/jpmarques19/treehouse/internal/output"
)

// initExitCode stores the exit code for the init command
var initExitCode int

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:           "init",
	Short:         "Initialize treehouse workspace",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		initExitCode = runInit()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// runInit executes the init command logic and returns the exit code
func runInit() int {
	// 1. Validate git repository
	repoInfo, err := git.DetectRepo()
	if err != nil {
		gitErr, ok := err.(*git.GitError)
		if ok {
			output.PrintError(gitErr.Code, gitErr.Message)
			return 3
		}
		output.PrintError("GIT_ERROR", err.Error())
		return 3
	}

	// 2. Check git version
	if err := git.CheckVersion(); err != nil {
		gitErr, ok := err.(*git.GitError)
		if ok {
			output.PrintError(gitErr.Code, gitErr.Message)
			return 3
		}
		output.PrintError("GIT_ERROR", err.Error())
		return 3
	}

	// 3. Check if .treehouse already exists
	treehousePath := filepath.Join(repoInfo.Root, ".treehouse")
	if _, err := os.Stat(treehousePath); err == nil {
		output.PrintError("INIT_ALREADY_EXISTS", "Treehouse already initialized")
		return 1
	}

	// 4. Update .gitignore before creating structure
	gitignoreUpdated, err := updateGitignore(repoInfo.Root)
	if err != nil {
		output.PrintError("INIT_GITIGNORE_FAILED", err.Error())
		return 1
	}

	// 5. Create folder structure
	created, err := createTreehouseStructure(treehousePath)
	if err != nil {
		output.PrintError("INIT_CREATE_FAILED", err.Error())
		return 1
	}

	// 6. Return success
	output.PrintSuccess(map[string]interface{}{
		"path":             treehousePath,
		"created":          created,
		"gitignore_updated": gitignoreUpdated,
	})
	return 0
}

// createTreehouseStructure creates the .treehouse folder structure
// If any step fails, it cleans up partially created files/directories
func createTreehouseStructure(basePath string) ([]string, error) {
	created := []string{}

	// Create directories
	dirs := []string{
		basePath,
		filepath.Join(basePath, "nooks"),
		filepath.Join(basePath, "hats"),
		filepath.Join(basePath, "workflows"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			// Cleanup on failure
			_ = os.RemoveAll(basePath)
			return nil, err
		}
	}

	// Create decks.yaml
	decksPath := filepath.Join(basePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte("decks: {}\n"), 0644); err != nil {
		_ = os.RemoveAll(basePath)
		return nil, err
	}
	created = append(created, "decks.yaml")
	created = append(created, "nooks/")
	created = append(created, "hats/")

	// Install workflows from embedded assets
	workflowsPath := filepath.Join(basePath, "workflows")
	installedWorkflows, err := installWorkflows(workflowsPath)
	if err != nil {
		_ = os.RemoveAll(basePath)
		return nil, err
	}
	for _, w := range installedWorkflows {
		created = append(created, "workflows/"+w)
	}

	// Install Claude commands
	repoRoot := filepath.Dir(basePath) // basePath is .treehouse, parent is repo root
	claudeCommandsPath := filepath.Join(repoRoot, ".claude", "commands", "th", "workflows")
	installedCommands, err := installClaudeCommands(claudeCommandsPath, repoRoot)
	if err != nil {
		_ = os.RemoveAll(basePath)
		return nil, err
	}
	for _, c := range installedCommands {
		created = append(created, ".claude/commands/th/workflows/"+c)
	}

	return created, nil
}

// installWorkflows copies embedded workflow files to .treehouse/workflows/
func installWorkflows(destPath string) ([]string, error) {
	installed := []string{}

	// Read workflow files from embedded assets
	entries, err := assets.Workflows.ReadDir("workflows")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip claude/ subdirectory, handled separately
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		content, err := assets.Workflows.ReadFile("workflows/" + entry.Name())
		if err != nil {
			return nil, err
		}

		destFile := filepath.Join(destPath, entry.Name())
		if err := os.WriteFile(destFile, content, 0644); err != nil {
			return nil, err
		}
		installed = append(installed, entry.Name())
	}

	return installed, nil
}

// installClaudeCommands copies embedded Claude command stubs to .claude/commands/th/workflows/
// and replaces the {{TREEHOUSE_BASE_WORKSPACE}} placeholder with the actual repoRoot path
func installClaudeCommands(destPath string, repoRoot string) ([]string, error) {
	installed := []string{}

	// Create directory structure
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return nil, err
	}

	// Read Claude command files from embedded assets
	entries, err := assets.Workflows.ReadDir("workflows/claude")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		content, err := assets.Workflows.ReadFile("workflows/claude/" + entry.Name())
		if err != nil {
			return nil, err
		}

		// Replace placeholder with actual repoRoot path
		contentStr := string(content)
		contentStr = strings.Replace(contentStr, "{{TREEHOUSE_BASE_WORKSPACE}}", repoRoot, -1)

		destFile := filepath.Join(destPath, entry.Name())
		if err := os.WriteFile(destFile, []byte(contentStr), 0644); err != nil {
			return nil, err
		}
		installed = append(installed, entry.Name())
	}

	return installed, nil
}

// treehouseGitignoreMarker is the marker used to identify the treehouse managed block
const treehouseGitignoreMarker = "# treehouse managed"

// treehouseGitignoreBlock is the content added to .gitignore
const treehouseGitignoreBlock = `# treehouse managed
.treehouse/
.claude/commands/th
`

// updateGitignore adds treehouse entries to .gitignore if not already present
// Returns true if the file was updated, false if entries already existed
func updateGitignore(repoRoot string) (bool, error) {
	gitignorePath := filepath.Join(repoRoot, ".gitignore")

	// Check if .gitignore exists
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new .gitignore with our block
			if err := os.WriteFile(gitignorePath, []byte(treehouseGitignoreBlock), 0644); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	// Check if treehouse entries already exist
	if containsTreehouseEntries(string(content)) {
		return false, nil // Already present, no update needed
	}

	// Append treehouse block to existing content
	newContent := string(content)
	if len(newContent) > 0 && newContent[len(newContent)-1] != '\n' {
		newContent += "\n"
	}
	newContent += "\n" + treehouseGitignoreBlock

	if err := os.WriteFile(gitignorePath, []byte(newContent), 0644); err != nil {
		return false, err
	}

	return true, nil
}

// containsTreehouseEntries checks if the gitignore content already has treehouse entries
func containsTreehouseEntries(content string) bool {
	// Check for the managed marker or the actual entries
	if strings.Contains(content, treehouseGitignoreMarker) {
		return true
	}
	// Also check for the individual entries in case added manually
	if strings.Contains(content, ".treehouse/") && strings.Contains(content, ".claude/commands/th") {
		return true
	}
	return false
}
