package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorktreeError represents a worktree operation error
type WorktreeError struct {
	Code    string
	Message string
}

func (e *WorktreeError) Error() string {
	return e.Message
}

// ParentInfo contains information about the current context's parent
type ParentInfo struct {
	IsNook       bool   // True if currently in a nook worktree
	ParentID     string // Nook ID if in a nook, branch name if in base repo
	CurrentNook  string // Current nook ID if in a nook (empty if base repo)
	WorktreeRoot string // Root of the current worktree/repo
}

// GetParentInfo detects the current context (base repo or nook) and returns parent info
// This is used to determine what the parent of a new nook should be
func GetParentInfo(treehousePath string) (*ParentInfo, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, &WorktreeError{
			Code:    "WORKTREE_CWD_ERROR",
			Message: fmt.Sprintf("Failed to get current directory: %v", err),
		}
	}

	// Check if we're in a git worktree
	// git rev-parse --show-toplevel gives us the root
	// git rev-parse --git-dir gives us the git dir (if worktree, it's .git/worktrees/<name>)
	gitDirOut, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return nil, &WorktreeError{
			Code:    "GIT_DIR_ERROR",
			Message: fmt.Sprintf("Failed to get git directory: %v", err),
		}
	}

	gitDir := strings.TrimSpace(string(gitDirOut))

	// Get worktree root
	rootOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, &WorktreeError{
			Code:    "GIT_ROOT_ERROR",
			Message: fmt.Sprintf("Failed to get git root: %v", err),
		}
	}
	worktreeRoot := strings.TrimSpace(string(rootOut))

	// Check if gitDir indicates we're in a worktree
	// If we're in a worktree, gitDir will be something like /path/to/repo/.git/worktrees/<nook-id>
	// If we're in base repo, gitDir will be .git or /path/to/repo/.git
	isWorktree := strings.Contains(gitDir, filepath.Join(".git", "worktrees"))

	if isWorktree {
		// Extract nook ID from the worktree directory name
		// The worktree folder name should be the nook ID
		worktreeBaseName := filepath.Base(worktreeRoot)

		// The worktree name should be the nook ID
		return &ParentInfo{
			IsNook:       true,
			ParentID:     worktreeBaseName, // Parent for sub-fork is current nook
			CurrentNook:  worktreeBaseName,
			WorktreeRoot: worktreeRoot,
		}, nil
	}

	// We're in the base repository
	// Get current branch name
	branchOut, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return nil, &WorktreeError{
			Code:    "GIT_BRANCH_ERROR",
			Message: fmt.Sprintf("Failed to get current branch: %v", err),
		}
	}

	branchName := strings.TrimSpace(string(branchOut))

	// Check for detached HEAD
	if branchName == "HEAD" {
		return nil, &WorktreeError{
			Code:    "GIT_DETACHED_HEAD",
			Message: "Cannot fork from detached HEAD. Checkout a branch first",
		}
	}

	return &ParentInfo{
		IsNook:       false,
		ParentID:     branchName, // Parent is the branch name (main, dev, etc.)
		CurrentNook:  "",
		WorktreeRoot: cwd, // In base repo, this is just cwd
	}, nil
}

// Create creates a new git worktree with the given nook ID
// Returns the absolute path to the created worktree
func Create(worktreesPath string, nookID string, parentBranch string) (string, error) {
	// Ensure worktrees directory exists
	if err := os.MkdirAll(worktreesPath, 0755); err != nil {
		return "", &WorktreeError{
			Code:    "GIT_WORKTREE_FAILED",
			Message: fmt.Sprintf("Failed to create worktrees directory: %v", err),
		}
	}

	// Calculate worktree path
	worktreePath := filepath.Join(worktreesPath, nookID)

	// Check if worktree path already exists
	if _, err := os.Stat(worktreePath); err == nil {
		return "", &WorktreeError{
			Code:    "NOOK_ALREADY_EXISTS",
			Message: fmt.Sprintf("Worktree path already exists: %s", worktreePath),
		}
	}

	// Create the worktree with a new branch
	// git worktree add -b <branch-name> <path> [<commit-ish>]
	// We create a new branch named after the nook ID, based on parent branch
	cmd := exec.Command("git", "worktree", "add", "-b", nookID, worktreePath, parentBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &WorktreeError{
			Code:    "GIT_WORKTREE_FAILED",
			Message: fmt.Sprintf("Failed to create worktree: %s", strings.TrimSpace(string(output))),
		}
	}

	// Return absolute path
	absPath, err := filepath.Abs(worktreePath)
	if err != nil {
		return worktreePath, nil
	}

	return absPath, nil
}

// Remove removes a git worktree by nook ID
func Remove(worktreesPath string, nookID string) error {
	worktreePath := filepath.Join(worktreesPath, nookID)

	// Check if worktree exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return &WorktreeError{
			Code:    "NOOK_NOT_FOUND",
			Message: fmt.Sprintf("Worktree not found: %s", worktreePath),
		}
	}

	// Remove the worktree
	cmd := exec.Command("git", "worktree", "remove", worktreePath, "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &WorktreeError{
			Code:    "GIT_WORKTREE_FAILED",
			Message: fmt.Sprintf("Failed to remove worktree: %s", strings.TrimSpace(string(output))),
		}
	}

	// Delete the branch
	cmd = exec.Command("git", "branch", "-D", nookID)
	// Ignore branch deletion errors - branch might not exist
	_ = cmd.Run()

	return nil
}
