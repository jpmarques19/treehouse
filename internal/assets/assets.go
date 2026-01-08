// Package assets provides embedded files for treehouse installation
package assets

import "embed"

// Workflows contains all workflow files to be installed by th init
// - workflows/*.md -> .treehouse/workflows/
// - workflows/claude/*.md -> .claude/commands/th/workflows/
//
//go:embed workflows/*.md
//go:embed workflows/claude/*.md
var Workflows embed.FS
