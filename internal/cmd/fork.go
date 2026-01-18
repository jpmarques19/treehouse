package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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

	// 3b. Validate commit SHA is present
	if repoInfo.CommitSHA == "" {
		output.PrintError("GIT_NO_COMMITS", "No commits found. Make at least one commit first")
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

	// 5. Load decks.yaml (once, reuse throughout)
	decks, err := deck.LoadDecks(thInfo.TreehousePath)
	if err != nil {
		deckErr, ok := err.(*deck.DeckError)
		if ok && deckErr.Code == "DECK_FILE_NOT_FOUND" {
			// Create empty decks structure - file will be created on save
			decks = &deck.DecksFile{Decks: make(map[string]*deck.Deck)}
		} else {
			if ok {
				output.PrintError(deckErr.Code, deckErr.Message)
			} else {
				output.PrintError("DECK_ERROR", err.Error())
			}
			return 1
		}
	}

	// 6. Check if nook already exists
	for _, d := range decks.Decks {
		if _, exists := d.Nooks[nookID]; exists {
			output.PrintError("NOOK_ALREADY_EXISTS", "Nook '"+nookID+"' already exists")
			return 1
		}
	}

	// 7. Get parent info (are we in base repo or a nook?)
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

	// 8. Use worktrees path from treehouse info (.treehouse/nooks/)
	worktreesPath := thInfo.WorktreesPath

	// 9. Create the git worktree
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

	// 9b. Set up Claude commands symlink (if th commands exist in base repo)
	if err := setupClaudeCommandsSymlink(thInfo.RepoRoot, worktreePath); err != nil {
		// Log warning but don't fail - th commands are optional
		// Could add debug logging here if needed
	}

	// 10. Determine deck ID
	// If forking from base branch (main/dev), create new deck
	// If forking from existing nook, use same deck
	var deckID string
	nookHash := nookID[:4] // First 4 chars of nook ID

	if parentInfo.IsNook {
		// Sub-fork: find parent's deck
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

	created := time.Now().Format("2006-01-02")
	deck.AddNookToDeck(decks, deckID, nookID, parentInfo.ParentID, created)

	if err := deck.SaveDecks(thInfo.TreehousePath, decks); err != nil {
		// Rollback worktree - log if rollback also fails
		if rollbackErr := worktree.Remove(worktreesPath, nookID); rollbackErr != nil {
			// Worktree cleanup failed - user may need to manually remove
			output.PrintError("DECK_WRITE_FAILED", err.Error()+" (warning: worktree cleanup also failed, manual removal may be needed)")
			return 1
		}
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

// setupClaudeCommandsSymlink creates a symlink in the nook for th Claude commands
// pointing back to the base repo's .claude/commands/th directory.
// Returns nil if source doesn't exist (th commands are optional).
func setupClaudeCommandsSymlink(repoRoot, nookPath string) error {
	// Check if source exists
	sourcePath := filepath.Join(repoRoot, ".claude", "commands", "th")
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return nil // th commands don't exist in base repo, skip
	}

	// Create .claude/commands/ directory in nook
	nookCommandsDir := filepath.Join(nookPath, ".claude", "commands")
	if err := os.MkdirAll(nookCommandsDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude/commands: %w", err)
	}

	// Create relative symlink
	// From: {nook}/.claude/commands/th
	// To:   {repo}/.claude/commands/th
	// Relative: ../../../../../.claude/commands/th
	symlinkPath := filepath.Join(nookCommandsDir, "th")
	relTarget := filepath.Join("..", "..", "..", "..", "..", ".claude", "commands", "th")

	if err := os.Symlink(relTarget, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}
