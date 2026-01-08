package treehouse

import (
	"fmt"
	"os"
	"path/filepath"
)

// TreehouseError represents a treehouse operation error
type TreehouseError struct {
	Code    string
	Message string
}

func (e *TreehouseError) Error() string {
	return e.Message
}

// Info contains information about a treehouse workspace
type Info struct {
	TreehousePath string // Path to .treehouse directory
	RepoRoot      string // Path to the repository root
	WorktreesPath string // Path to worktrees directory
}

// FindTreehouse searches for a .treehouse directory starting from the given
// path and walking up the directory tree. Returns the treehouse info if found.
func FindTreehouse(startPath string) (*Info, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return nil, &TreehouseError{
			Code:    "INIT_PATH_ERROR",
			Message: fmt.Sprintf("Failed to resolve path: %v", err),
		}
	}

	currentPath := absPath

	// Walk up the directory tree
	for {
		treehousePath := filepath.Join(currentPath, ".treehouse")
		info, err := os.Stat(treehousePath)
		if err == nil && info.IsDir() {
			// Found it!
			return &Info{
				TreehousePath: treehousePath,
				RepoRoot:      currentPath,
				WorktreesPath: filepath.Join(currentPath, "..", "worktrees"),
			}, nil
		}

		// Move to parent directory
		parentPath := filepath.Dir(currentPath)

		// Check if we've reached the root
		if parentPath == currentPath {
			return nil, &TreehouseError{
				Code:    "INIT_NOT_FOUND",
				Message: "Treehouse not initialized. Run 'th init' first",
			}
		}

		currentPath = parentPath
	}
}

// Exists checks if a treehouse workspace exists at or above the given path
func Exists(startPath string) bool {
	_, err := FindTreehouse(startPath)
	return err == nil
}
