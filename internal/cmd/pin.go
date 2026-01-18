package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/board"
	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/treehouse"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// pinExitCode stores the exit code for the pin command
var pinExitCode int

// pinCmd represents the pin command
var pinCmd = &cobra.Command{
	Use:           "pin <content>",
	Short:         "Add a pin to the current nook's board",
	Long:          "Save a note/learning to the current nook's board for future reference",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pinExitCode = runPin(args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pinCmd)
}

// PinResult contains the data returned by pin command
type PinResult struct {
	Nook      string `json:"nook"`
	PinsCount int    `json:"pins_count"`
	Ts        string `json:"ts"`
}

// runPin executes the pin command logic and returns the exit code
func runPin(content string) int {
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
		output.PrintError("PIN_NOT_IN_NOOK", "Must be in a nook to pin")
		return 1
	}

	// 3. Verify we're in a nook
	if !parentInfo.IsNook {
		output.PrintError("PIN_NOT_IN_NOOK", "Must be in a nook to pin")
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

	// 5. Add pin
	ts := board.AddPin(boards, nookID, content)

	// 6. Save boards atomically
	if err := board.SaveBoards(nooksPath, boards); err != nil {
		boardErr, ok := err.(*board.BoardError)
		if ok {
			output.PrintError(boardErr.Code, boardErr.Message)
		} else {
			output.PrintError("BOARD_SAVE_FAILED", err.Error())
		}
		return 1
	}

	// 7. Return success
	result := PinResult{
		Nook:      nookID,
		PinsCount: len(boards.Boards[nookID]),
		Ts:        ts,
	}

	output.PrintSuccess(result)
	return 0
}
