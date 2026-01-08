package treehouse

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
// If we're in a git worktree (nook), it finds the base repo's .treehouse.
func FindTreehouse(startPath string) (*Info, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return nil, &TreehouseError{
			Code:    "INIT_PATH_ERROR",
			Message: fmt.Sprintf("Failed to resolve path: %v", err),
		}
	}

	// Check if we're in a git worktree (nook) - if so, find the base repo
	baseRepoPath := findBaseRepo(absPath)
	if baseRepoPath != "" {
		// We're in a worktree, use the base repo path
		treehousePath := filepath.Join(baseRepoPath, ".treehouse")
		info, err := os.Stat(treehousePath)
		if err == nil && info.IsDir() {
			return &Info{
				TreehousePath: treehousePath,
				RepoRoot:      baseRepoPath,
				WorktreesPath: filepath.Join(treehousePath, "nooks"),
			}, nil
		}
	}

	// Walk up the directory tree from current path
	currentPath := absPath
	for {
		treehousePath := filepath.Join(currentPath, ".treehouse")
		info, err := os.Stat(treehousePath)
		if err == nil && info.IsDir() {
			// Check if this .treehouse has a decks.yaml (indicates it's a real treehouse, not a nook's local folder)
			decksPath := filepath.Join(treehousePath, "decks.yaml")
			if _, err := os.Stat(decksPath); err == nil {
				return &Info{
					TreehousePath: treehousePath,
					RepoRoot:      currentPath,
					WorktreesPath: filepath.Join(treehousePath, "nooks"),
				}, nil
			}
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

// findBaseRepo checks if we're in a git worktree and returns the base repo path
func findBaseRepo(path string) string {
	// Run git rev-parse --git-dir to check if we're in a worktree
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	gitDir := strings.TrimSpace(string(output))

	// If git-dir contains "worktrees", we're in a worktree
	// Pattern: /path/to/base/.git/worktrees/nook-name
	if strings.Contains(gitDir, string(filepath.Separator)+"worktrees"+string(filepath.Separator)) {
		// Extract base repo: remove .git/worktrees/nook-name
		parts := strings.Split(gitDir, string(filepath.Separator)+".git"+string(filepath.Separator)+"worktrees")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	return ""
}

// Exists checks if a treehouse workspace exists at or above the given path
func Exists(startPath string) bool {
	_, err := FindTreehouse(startPath)
	return err == nil
}
