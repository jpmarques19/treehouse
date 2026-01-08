package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/jpmarques19/treehouse/internal/output"
)

// Version is set by ldflags at build time
var Version = "0.4.0"

var versionFlag bool

// SetOutput sets the output writer for JSON responses (used in testing)
func SetOutput(w io.Writer) {
	output.SetWriter(w)
}

// rootCmd is the base command for the th CLI
var rootCmd = &cobra.Command{
	Use:           "th",
	Short:         "Treehouse CLI for managing isolated development environments",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			output.PrintSuccess(map[string]string{"version": Version})
			return nil
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			// Already handled in PersistentPreRunE
			return nil
		}
		output.PrintError("NO_COMMAND", "No command specified")
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")
	// Disable default help flag - we output JSON for everything
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		output.PrintError("NO_COMMAND", "No command specified. Use a subcommand.")
	})
}

// Execute runs the root command with the provided arguments
// Returns the exit code (0=success, 1=general error, 2=invalid args, 3=git error)
func Execute(args []string) int {
	// Reset version flag and exit codes for each execution (important for tests)
	versionFlag = false
	initExitCode = 0
	listExitCode = 0

	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	if err != nil {
		// Handle unknown flags and other parsing errors
		output.PrintError("INVALID_ARGS", err.Error())
		return 2
	}

	// If version was requested, return success
	if versionFlag {
		return 0
	}

	// If no args provided (no subcommand), return exit code 2
	if len(args) == 0 {
		return 2
	}

	// Check if init command was executed and return its exit code
	if len(args) > 0 && args[0] == "init" {
		return initExitCode
	}

	// Check if list command was executed and return its exit code
	if len(args) > 0 && args[0] == "list" {
		return listExitCode
	}

	return 0
}
