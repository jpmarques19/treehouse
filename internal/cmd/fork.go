package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/deck"
	"github.com/jpmarques19/treehouse/internal/git"
	"github.com/jpmarques19/treehouse/internal/nook"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// forkExitCode stores the exit code for the fork command
var forkExitCode int

// forkCmd represents the fork command
var forkCmd = &cobra.Command{
	Use:           "fork <name>",
	Short:         "Create a new nook from current branch",
	Long:          "Create an isolated nook (git worktree) from the current branch for exploration",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			output.PrintError("NOOK_NAME_REQUIRED", "Nook name required: th fork <name>")
			forkExitCode = 2
			return nil
		}
		forkExitCode = runFork(args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(forkCmd)
}

// ForkResult contains the data returned on successful fork
type ForkResult struct {
	NookID   string `json:"nook_id"`
	DeckID   string `json:"deck_id"`
	Parent   string `json:"parent"`
	Worktree string `json:"worktree"`
}

// runFork executes the fork command logic and returns the exit code
func runFork(name string) int {
	// 1. Validate nook name
	if err := nook.ValidateName(name); err != nil {
		nookErr, ok := err.(*nook.NookError)
		if ok {
			output.PrintError(nookErr.Code, nookErr.Message)
		} else {
			output.PrintError("NOOK_NAME_INVALID", err.Error())
		}
		return 2
	}

	// 2. Find treehouse (this validates we're in an initialized repo)
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

	// 3. Get git repo info for commit SHA
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

	// 4. Generate nook ID
	nookID, err := nook.GenerateID(name, repoInfo.CommitSHA)
	if err != nil {
		nookErr, ok := err.(*nook.NookError)
		if ok {
			output.PrintError(nookErr.Code, nookErr.Message)
		} else {
			output.PrintError("NOOK_ERROR", err.Error())
		}
		return 1
	}

	// 5. Check if nook already exists in decks.yaml
	exists, err := deck.NookExists(thInfo.TreehousePath, nookID)
	if err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok && deckErr.Code != "DECK_FILE_NOT_FOUND" {
			output.PrintError(deckErr.Code, deckErr.Message)
			return 1
		}
		// If deck file not found, that's OK - we'll create it
	}
	if exists {
		output.PrintError("NOOK_ALREADY_EXISTS", "Nook '"+nookID+"' already exists")
		return 1
	}

	// 6. Get parent info (are we in base repo or a nook?)
	parentInfo, err := worktree.GetParentInfo(thInfo.TreehousePath)
	if err != nil {
		wtErr, ok := err.(*worktree.WorktreeError)
		if ok {
			output.PrintError(wtErr.Code, wtErr.Message)
			return 3
		}
		output.PrintError("WORKTREE_ERROR", err.Error())
		return 3
	}

	// 7. Use worktrees path from treehouse info (.treehouse/nooks/)
	worktreesPath := thInfo.WorktreesPath

	// 8. Create the git worktree
	worktreePath, err := worktree.Create(worktreesPath, nookID, parentInfo.ParentID)
	if err != nil {
		wtErr, ok := err.(*worktree.WorktreeError)
		if ok {
			output.PrintError(wtErr.Code, wtErr.Message)
			return 3
		}
		output.PrintError("GIT_WORKTREE_FAILED", err.Error())
		return 3
	}

	// 9. Determine deck ID
	// If forking from base branch (main/dev), create new deck
	// If forking from existing nook, use same deck
	var deckID string
	nookHash := nookID[:4] // First 4 chars of nook ID

	if parentInfo.IsNook {
		// Sub-fork: find parent's deck
		decks, err := deck.LoadDecks(thInfo.TreehousePath)
		if err != nil {
			// Worktree created but can't find parent deck - rollback
			_ = worktree.Remove(worktreesPath, nookID)
			deckErr, ok := err.(*deck.DeckError)
			if ok {
				output.PrintError(deckErr.Code, deckErr.Message)
			} else {
				output.PrintError("DECK_ERROR", err.Error())
			}
			return 1
		}

		parentDeckID, found := deck.GetDeckForNook(decks, parentInfo.ParentID)
		if !found {
			// Parent nook not in any deck - strange state, create new deck anyway
			deckID = deck.GenerateDeckID(nookHash)
		} else {
			deckID = parentDeckID
		}
	} else {
		// Fork from base branch: create new deck
		deckID = deck.GenerateDeckID(nookHash)
	}

	// 10. Update decks.yaml
	decks, err := deck.LoadDecks(thInfo.TreehousePath)
	if err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok && deckErr.Code == "DECK_FILE_NOT_FOUND" {
			// Create empty decks file
			decks = &deck.DecksFile{Decks: make(map[string]*deck.Deck)}
		} else {
			// Rollback worktree
			_ = worktree.Remove(worktreesPath, nookID)
			if ok {
				output.PrintError(deckErr.Code, deckErr.Message)
			} else {
				output.PrintError("DECK_ERROR", err.Error())
			}
			return 1
		}
	}

	created := time.Now().Format("2006-01-02")
	deck.AddNookToDeck(decks, deckID, nookID, parentInfo.ParentID, created)

	if err := deck.SaveDecks(thInfo.TreehousePath, decks); err != nil {
		// Rollback worktree
		_ = worktree.Remove(worktreesPath, nookID)
		deckErr, ok := err.(*deck.DeckError)
		if ok {
			output.PrintError(deckErr.Code, deckErr.Message)
		} else {
			output.PrintError("DECK_WRITE_FAILED", err.Error())
		}
		return 1
	}

	// 11. Return success
	output.PrintSuccess(ForkResult{
		NookID:   nookID,
		DeckID:   deckID,
		Parent:   parentInfo.ParentID,
		Worktree: worktreePath,
	})
	return 0
}
