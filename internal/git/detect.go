package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compiled regex for version parsing (compiled once at package init)
var versionRegex = regexp.MustCompile(`git version (\d+)\.(\d+)(?:\.(\d+))?`)

// RepoInfo contains information about a git repository
type RepoInfo struct {
	Root      string // Absolute path to repository root
	Branch    string // Current branch name
	CommitSHA string // First 7 characters of HEAD commit
}

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
}

// GitError represents a git-related error with a code
type GitError struct {
	Code    string
	Message string
}

func (e *GitError) Error() string {
	return e.Message
}

// DetectRepo detects the git repository from the current working directory
// Returns RepoInfo with root path, branch name, and commit SHA
// Returns GitError with code "INIT_NOT_GIT_REPO" if not in a git repository
func DetectRepo() (*RepoInfo, error) {
	// Check if inside git repo and get root
	rootOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, &GitError{
			Code:    "INIT_NOT_GIT_REPO",
			Message: "Not a git repository",
		}
	}

	// Get current branch
	branchOut, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return nil, &GitError{
			Code:    "GIT_BRANCH_ERROR",
			Message: fmt.Sprintf("Failed to get branch: %v", err),
		}
	}

	// Get commit SHA (first 7 characters)
	shaOut, err := exec.Command("git", "rev-parse", "--short=7", "HEAD").Output()
	if err != nil {
		return nil, &GitError{
			Code:    "GIT_SHA_ERROR",
			Message: fmt.Sprintf("Failed to get commit SHA: %v", err),
		}
	}

	return &RepoInfo{
		Root:      strings.TrimSpace(string(rootOut)),
		Branch:    strings.TrimSpace(string(branchOut)),
		CommitSHA: strings.TrimSpace(string(shaOut)),
	}, nil
}

// CheckVersion verifies that the installed git version supports worktrees (2.5+)
// Returns GitError with code "GIT_VERSION_UNSUPPORTED" if git version is below 2.5
// Returns GitError with code "GIT_NOT_FOUND" if git is not installed
func CheckVersion() error {
	output, err := exec.Command("git", "--version").Output()
	if err != nil {
		return &GitError{
			Code:    "GIT_NOT_FOUND",
			Message: "Git not installed",
		}
	}

	version, err := parseVersion(string(output))
	if err != nil {
		return &GitError{
			Code:    "GIT_VERSION_PARSE_ERROR",
			Message: fmt.Sprintf("Failed to parse git version: %v", err),
		}
	}

	return checkVersionRequirement(version)
}

// parseVersion extracts version numbers from "git version X.Y.Z" output
// Handles formats: "git version 2.34.1", "git version 2.34", "git version 2.34.1 (Apple Git-146)"
func parseVersion(versionStr string) (Version, error) {
	// Use pre-compiled regex with optional patch version
	matches := versionRegex.FindStringSubmatch(versionStr)

	if len(matches) < 3 {
		return Version{}, fmt.Errorf("invalid git version format: %s", versionStr)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid major version: %v", err)
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor version: %v", err)
	}

	// Patch is optional (some git versions report only major.minor)
	var patch int
	if len(matches) > 3 && matches[3] != "" {
		patch, err = strconv.Atoi(matches[3])
		if err != nil {
			return Version{}, fmt.Errorf("invalid patch version: %v", err)
		}
	}

	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// checkVersionRequirement verifies git version is at least 2.5.0 (worktree support)
func checkVersionRequirement(v Version) error {
	// Git 2.5+ is required for worktree support
	if v.Major < 2 || (v.Major == 2 && v.Minor < 5) {
		return &GitError{
			Code:    "GIT_VERSION_UNSUPPORTED",
			Message: "Git 2.5+ required for worktree support",
		}
	}
	return nil
}
