package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/board"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// boardExitCode stores the exit code for the board command
var boardExitCode int

// boardCmd represents the board command
var boardCmd = &cobra.Command{
	Use:           "board",
	Short:         "View the current nook's board",
	Long:          "Display all pins on the current nook's board",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		boardExitCode = runBoard()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(boardCmd)
}

// BoardResult contains the data returned by board command
type BoardResult struct {
	Nook string       `json:"nook"`
	Pins []board.Pin  `json:"pins"`
}

// runBoard executes the board command logic and returns the exit code
func runBoard() int {
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

	// 2. Get current nook context
	parentInfo, err := worktree.GetParentInfo(thInfo.TreehousePath)
	if err != nil {
		output.PrintError("BOARD_NOT_IN_NOOK", "Must be in a nook to view board")
		return 1
	}

	// 3. Verify we're in a nook
	if !parentInfo.IsNook {
		output.PrintError("BOARD_NOT_IN_NOOK", "Must be in a nook to view board")
		return 1
	}

	nookID := parentInfo.CurrentNook

	// 4. Build nooks path and load boards
	nooksPath := filepath.Join(thInfo.TreehousePath, "nooks")
	boards, err := board.LoadBoards(nooksPath)
	if err != nil {
		boardErr, ok := err.(*board.BoardError)
		if ok {
			output.PrintError(boardErr.Code, boardErr.Message)
		} else {
			output.PrintError("BOARD_LOAD_FAILED", err.Error())
		}
		return 1
	}

	// 5. Get board for current nook
	pins := board.GetBoard(boards, nookID)

	// 6. Return success (empty pins array is valid, not an error)
	result := BoardResult{
		Nook: nookID,
		Pins: pins,
	}

	output.PrintSuccess(result)
	return 0
}
