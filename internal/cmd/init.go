package cmd

import (
	"os"
	"path/filepath"

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

	// 4. Create folder structure
	created, err := createTreehouseStructure(treehousePath)
	if err != nil {
		output.PrintError("INIT_CREATE_FAILED", err.Error())
		return 1
	}

	// 5. Return success
	output.PrintSuccess(map[string]interface{}{
		"path":    treehousePath,
		"created": created,
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
		filepath.Join(basePath, "crew"),
		filepath.Join(basePath, "crew", "oak"),
		filepath.Join(basePath, "crew", "oak", "memories"),
		filepath.Join(basePath, "crew", "oak", "sessions"),
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

	// Create Oak agent files
	if err := createOakAgent(filepath.Join(basePath, "crew", "oak")); err != nil {
		_ = os.RemoveAll(basePath)
		return nil, err
	}
	created = append(created, "crew/oak/")

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
	installedCommands, err := installClaudeCommands(claudeCommandsPath)
	if err != nil {
		_ = os.RemoveAll(basePath)
		return nil, err
	}
	for _, c := range installedCommands {
		created = append(created, ".claude/commands/th/workflows/"+c)
	}

	return created, nil
}

// createOakAgent creates the default Oak agent files
func createOakAgent(oakPath string) error {
	// Create oak.agent.yaml
	agentYAML := `agent:
  metadata:
    name: "Oak"
    title: "Treehouse Assistant"
    icon: "🌳"

  persona:
    role: "Navigate treehouse workspace and manage context"

    identity: |
      I am Oak, your treehouse assistant. I help you navigate between nooks,
      manage your exploration context, and ensure nothing gets lost when
      switching between branches of investigation.

      I understand the deck/nook structure and can help you organize your
      explorations effectively.

    communication_style: |
      Friendly and supportive. I provide clear guidance on treehouse operations
      and help you maintain context across your explorations.

    principles:
      - "Always detect current nook context before acting"
      - "Preserve context when switching between nooks"
      - "Help organize explorations into logical decks"

  critical_actions:
    - "Detect current nook from git worktree folder name"
    - "Construct paths: AGENT_FOLDER, MEMORY_FILE, SESSION_FILE"
    - "Load knowledge.md for global context"
    - "Load memories/{nook-id}.md for nook-specific context"
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
`
	agentPath := filepath.Join(oakPath, "oak.agent.yaml")
	if err := os.WriteFile(agentPath, []byte(agentYAML), 0644); err != nil {
		return err
	}

	// Create knowledge.md
	knowledgeMD := `# Oak - Project Knowledge

> Global context accumulated across all nooks

---

## Project Understanding

{What Oak has learned about the project}

## Key Patterns

{Patterns Oak has identified}

## Technical Debt

{Items to address later}

## Key Files

{Important file paths and their purposes}
`
	knowledgePath := filepath.Join(oakPath, "knowledge.md")
	if err := os.WriteFile(knowledgePath, []byte(knowledgeMD), 0644); err != nil {
		return err
	}

	return nil
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
func installClaudeCommands(destPath string) ([]string, error) {
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

		destFile := filepath.Join(destPath, entry.Name())
		if err := os.WriteFile(destFile, content, 0644); err != nil {
			return nil, err
		}
		installed = append(installed, entry.Name())
	}

	return installed, nil
}
