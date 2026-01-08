package main

import (
	"os"

	"github.com/treehouse-cli/th/internal/cmd"
)

func main() {
	exitCode := cmd.Execute(os.Args[1:])
	os.Exit(exitCode)
}
