package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/deck"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// pruneExitCode stores the exit code for the prune command
var pruneExitCode int

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:           "prune",
	Short:         "Remove orphaned nooks from decks.yaml",
	Long:          "Remove all orphaned nook entries (nooks where the worktree folder no longer exists)",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pruneExitCode = runPrune()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}

// PruneResult contains the data returned on successful prune
type PruneResult struct {
	Pruned       []string `json:"pruned"`
	DecksRemoved []string `json:"decks_removed"`
}

// runPrune executes the prune command logic and returns the exit code
func runPrune() int {
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

	// 2. Load decks
	decks, err := deck.LoadDecks(thInfo.TreehousePath)
	if err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok {
			// If no decks file, nothing to prune
			if deckErr.Code == "DECK_FILE_NOT_FOUND" {
				output.PrintSuccess(PruneResult{
					Pruned:       []string{},
					DecksRemoved: []string{},
				})
				return 0
			}
			output.PrintError(deckErr.Code, deckErr.Message)
		} else {
			output.PrintError("DECK_ERROR", err.Error())
		}
		return 1
	}

	// 3. Find orphaned nooks (entries in decks.yaml with missing worktrees)
	var orphans []string
	for _, d := range decks.Decks {
		if d.Nooks == nil {
			continue
		}
		for nookID := range d.Nooks {
			worktreePath := filepath.Join(thInfo.WorktreesPath, nookID)
			exists, _ := worktree.Exists(worktreePath)
			if !exists {
				orphans = append(orphans, nookID)
			}
		}
	}

	// 4. If no orphans, return empty result
	if len(orphans) == 0 {
		output.PrintSuccess(PruneResult{
			Pruned:       []string{},
			DecksRemoved: []string{},
		})
		return 0
	}

	// 5. Remove each orphan
	var pruned []string
	decksRemovedMap := make(map[string]bool)

	for _, nookID := range orphans {
		// Clean up crew memory files
		deleteCrewMemoryFiles(thInfo.TreehousePath, nookID)

		// Remove from decks.yaml
		deckID, found := deck.GetDeckForNook(decks, nookID)
		if found {
			empty, _ := deck.RemoveNook(decks, deckID, nookID)
			if empty {
				decksRemovedMap[deckID] = true
			}
		}

		pruned = append(pruned, nookID)
	}

	// 6. Save updated decks
	if err := deck.SaveDecks(thInfo.TreehousePath, decks); err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok {
			output.PrintError(deckErr.Code, deckErr.Message)
		} else {
			output.PrintError("DECK_WRITE_FAILED", err.Error())
		}
		return 1
	}

	// 7. Convert decks removed map to slice
	var decksRemoved []string
	for deckID := range decksRemovedMap {
		decksRemoved = append(decksRemoved, deckID)
	}

	// 8. Return success
	output.PrintSuccess(PruneResult{
		Pruned:       pruned,
		DecksRemoved: decksRemoved,
	})
	return 0
}
