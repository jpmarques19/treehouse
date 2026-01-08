package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/deck"
	"github.com/jpmarques19/treehouse/internal/git"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// listExitCode stores the exit code for the list command
var listExitCode int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:           "list",
	Short:         "List all decks and nooks",
	Long:          "Display all decks and nooks in JSON format for visualization",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		listExitCode = runList()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// ListResult contains the data returned by list command
type ListResult struct {
	Base        BaseInfo     `json:"base"`
	CurrentNook *string      `json:"current_nook"`
	Decks       []DeckResult `json:"decks"`
}

// BaseInfo contains information about the base repository
type BaseInfo struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

// DeckResult contains deck information for list output
type DeckResult struct {
	ID      string       `json:"id"`
	Created string       `json:"created"`
	Nooks   []NookResult `json:"nooks"`
}

// NookResult contains nook information for list output
type NookResult struct {
	ID       string `json:"id"`
	Parent   string `json:"parent"`
	Created  string `json:"created"`
	Worktree string `json:"worktree"`
	Status   string `json:"status"`
}

// runList executes the list command logic and returns the exit code
func runList() int {
	// 1. Find treehouse (validates we're in an initialized repo)
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

	// 2. Load decks
	decksData, err := deck.LoadDecks(thInfo.TreehousePath)
	if err != nil {
		output.PrintError("DECK_LOAD_ERROR", err.Error())
		return 1
	}

	// 3. Get base repository information
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

	baseInfo := BaseInfo{
		Path:   thInfo.RepoRoot,
		Branch: repoInfo.Branch,
		Commit: repoInfo.CommitSHA, // Already short form (7 chars)
	}

	// 4. Determine current nook
	parentInfo, err := worktree.GetParentInfo(thInfo.TreehousePath)
	if err != nil {
		// If we can't detect parent info, assume we're in base repo
		parentInfo = &worktree.ParentInfo{
			IsNook:      false,
			CurrentNook: "",
		}
	}

	var currentNook *string
	if parentInfo.IsNook && parentInfo.CurrentNook != "" {
		currentNook = &parentInfo.CurrentNook
	}

	// 5. Build deck results with nook status
	deckResults := buildDeckResults(decksData, thInfo, currentNook)

	// 6. Output result
	result := ListResult{
		Base:        baseInfo,
		CurrentNook: currentNook,
		Decks:       deckResults,
	}

	output.PrintSuccess(result)
	return 0
}

// buildDeckResults converts deck data to output format with nook status
func buildDeckResults(decksData *deck.DecksFile, thInfo *treehouse.Info, currentNook *string) []DeckResult {
	if decksData == nil || len(decksData.Decks) == 0 {
		return []DeckResult{}
	}

	results := make([]DeckResult, 0, len(decksData.Decks))

	for deckID, deckData := range decksData.Decks {
		deckResult := DeckResult{
			ID:      deckID,
			Created: deckData.Created,
			Nooks:   make([]NookResult, 0, len(deckData.Nooks)),
		}

		for nookID, nookData := range deckData.Nooks {
			// Build worktree path
			worktreePath := filepath.Join(thInfo.WorktreesPath, nookID)

			// Determine status
			status := determineNookStatus(nookID, worktreePath, currentNook)

			nookResult := NookResult{
				ID:       nookID,
				Parent:   nookData.Parent,
				Created:  nookData.Created,
				Worktree: worktreePath,
				Status:   status,
			}

			deckResult.Nooks = append(deckResult.Nooks, nookResult)
		}

		results = append(results, deckResult)
	}

	return results
}

// determineNookStatus determines the status of a nook (active, inactive, orphan)
func determineNookStatus(nookID string, worktreePath string, currentNook *string) string {
	// Check if this is the current nook
	if currentNook != nil && *currentNook == nookID {
		return "active"
	}

	// Check if worktree exists
	exists, err := worktree.Exists(worktreePath)
	if err != nil || !exists {
		return "orphan"
	}

	return "inactive"
}
