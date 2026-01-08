package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/deck"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// removeExitCode stores the exit code for the remove command
var removeExitCode int

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:           "remove <nook-id>",
	Short:         "Remove a nook and its children",
	Long:          "Remove a specific nook (git worktree) and all its child nooks recursively",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			output.PrintError("NOOK_ID_REQUIRED", "Nook ID required: th remove <nook-id>")
			removeExitCode = 2
			return nil
		}
		removeExitCode = runRemove(args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

// RemoveResult contains the data returned on successful remove
type RemoveResult struct {
	Removed     []string `json:"removed"`
	DeckRemoved bool     `json:"deck_removed"`
}

// runRemove executes the remove command logic and returns the exit code
func runRemove(nookID string) int {
	// 1. Find treehouse
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

	// 2. Check if trying to remove current nook
	parentInfo, err := worktree.GetParentInfo(thInfo.TreehousePath)
	if err != nil {
		wtErr, ok := err.(*worktree.WorktreeError)
		if ok {
			output.PrintError(wtErr.Code, wtErr.Message)
			return 1
		}
		output.PrintError("WORKTREE_ERROR", err.Error())
		return 1
	}

	if parentInfo.IsNook && parentInfo.CurrentNook == nookID {
		output.PrintError("NOOK_IS_CURRENT", "Cannot remove current nook. Switch to base repo first.")
		return 1
	}

	// 3. Load decks and verify nook exists
	decks, err := deck.LoadDecks(thInfo.TreehousePath)
	if err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok {
			output.PrintError(deckErr.Code, deckErr.Message)
		} else {
			output.PrintError("DECK_ERROR", err.Error())
		}
		return 1
	}

	deckID, found := deck.GetDeckForNook(decks, nookID)
	if !found {
		output.PrintError("NOOK_NOT_FOUND", "Nook '"+nookID+"' not found")
		return 4
	}

	// 4. Collect all nooks to remove (recursive children + target)
	nooksToRemove := collectNooksToRemove(decks, nookID)

	// 5. Remove each nook (children first, then parent)
	var removed []string
	for _, id := range nooksToRemove {
		// Remove worktree (ignore error if already missing)
		_ = worktree.Remove(thInfo.WorktreesPath, id)

		// Clean up crew memory files
		deleteCrewMemoryFiles(thInfo.TreehousePath, id)

		removed = append(removed, id)
	}

	// 6. Remove nooks from decks.yaml (in reverse order to handle parent-child)
	var deckRemoved bool
	for i := len(nooksToRemove) - 1; i >= 0; i-- {
		id := nooksToRemove[i]
		nookDeckID, _ := deck.GetDeckForNook(decks, id)
		if nookDeckID != "" {
			empty, _ := deck.RemoveNook(decks, nookDeckID, id)
			if empty && nookDeckID == deckID {
				deckRemoved = true
			}
		}
	}

	// 7. Save updated decks
	if err := deck.SaveDecks(thInfo.TreehousePath, decks); err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok {
			output.PrintError(deckErr.Code, deckErr.Message)
		} else {
			output.PrintError("DECK_WRITE_FAILED", err.Error())
		}
		return 1
	}

	// 8. Return success
	output.PrintSuccess(RemoveResult{
		Removed:     removed,
		DeckRemoved: deckRemoved,
	})
	return 0
}

// collectNooksToRemove returns all nooks that need to be removed in order (children first)
func collectNooksToRemove(decks *deck.DecksFile, nookID string) []string {
	var result []string

	// Get direct children
	children := deck.GetChildNooks(decks, nookID)

	// Recursively collect all descendants (depth-first)
	for _, child := range children {
		result = append(result, collectNooksToRemove(decks, child)...)
	}

	// Add this nook last (after all children)
	result = append(result, nookID)

	return result
}

// deleteCrewMemoryFiles removes memory and session files for a nook across all agents
func deleteCrewMemoryFiles(treehousePath, nookID string) {
	crewPath := filepath.Join(treehousePath, "crew")

	// Check if crew directory exists
	if _, err := os.Stat(crewPath); os.IsNotExist(err) {
		return
	}

	// Pattern: .treehouse/crew/*/memories/{nook-id}.md
	memPattern := filepath.Join(crewPath, "*", "memories", nookID+".md")
	memFiles, _ := filepath.Glob(memPattern)
	for _, f := range memFiles {
		os.Remove(f) // Ignore errors - file may not exist
	}

	// Pattern: .treehouse/crew/*/sessions/{nook-id}.md
	sessPattern := filepath.Join(crewPath, "*", "sessions", nookID+".md")
	sessFiles, _ := filepath.Glob(sessPattern)
	for _, f := range sessFiles {
		os.Remove(f) // Ignore errors - file may not exist
	}
}
