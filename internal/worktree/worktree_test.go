package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupGitRepo creates a temporary git repository for testing
func setupGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to config user.email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to config user.name: %v", err)
	}

	// Create an initial commit
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git commit: %v", err)
	}

	return dir
}

func TestGetParentInfo_BaseRepo(t *testing.T) {
	repoDir := setupGitRepo(t)

	// Change to repo directory
	originalDir, _ := os.Getwd()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo dir: %v", err)
	}
	defer os.Chdir(originalDir)

	treehousePath := filepath.Join(repoDir, ".treehouse")

	info, err := GetParentInfo(treehousePath)
	if err != nil {
		t.Fatalf("GetParentInfo() error = %v", err)
	}

	if info.IsNook {
		t.Error("Expected IsNook to be false in base repo")
	}

	// Default branch could be 'main' or 'master' depending on git config
	if info.ParentID == "" {
		t.Error("Expected ParentID to be non-empty")
	}

	if info.CurrentNook != "" {
		t.Errorf("Expected CurrentNook to be empty, got %q", info.CurrentNook)
	}
}

func TestCreate(t *testing.T) {
	repoDir := setupGitRepo(t)

	// Change to repo directory for git commands
	originalDir, _ := os.Getwd()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo dir: %v", err)
	}
	defer os.Chdir(originalDir)

	// Get current branch
	branchOut, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	parentBranch := string(branchOut)
	parentBranch = parentBranch[:len(parentBranch)-1] // trim newline

	// Create worktrees directory at sibling level
	worktreesPath := filepath.Join(repoDir, "..", "worktrees")

	worktreePath, err := Create(worktreesPath, "a1b2-test-nook", parentBranch)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Errorf("Worktree path does not exist: %s", worktreePath)
	}

	// Verify it's a valid git repo
	cmd := exec.Command("git", "status")
	cmd.Dir = worktreePath
	if err := cmd.Run(); err != nil {
		t.Errorf("Worktree is not a valid git repo: %v", err)
	}

	// Clean up
	_ = exec.Command("git", "worktree", "remove", worktreePath, "--force").Run()
	_ = exec.Command("git", "branch", "-D", "a1b2-test-nook").Run()
}

func TestCreate_AlreadyExists(t *testing.T) {
	repoDir := setupGitRepo(t)

	// Change to repo directory
	originalDir, _ := os.Getwd()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo dir: %v", err)
	}
	defer os.Chdir(originalDir)

	// Get current branch
	branchOut, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	parentBranch := string(branchOut)
	parentBranch = parentBranch[:len(parentBranch)-1]

	worktreesPath := filepath.Join(repoDir, "..", "worktrees")

	// Create first worktree
	worktreePath, err := Create(worktreesPath, "a1b2-duplicate", parentBranch)
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	// Try to create again - should fail
	_, err = Create(worktreesPath, "a1b2-duplicate", parentBranch)
	if err == nil {
		t.Fatal("Expected error when creating duplicate worktree")
	}

	wtErr, ok := err.(*WorktreeError)
	if !ok {
		t.Fatalf("Expected WorktreeError, got %T", err)
	}

	if wtErr.Code != "NOOK_ALREADY_EXISTS" {
		t.Errorf("Error code = %q, want %q", wtErr.Code, "NOOK_ALREADY_EXISTS")
	}

	// Clean up
	_ = exec.Command("git", "worktree", "remove", worktreePath, "--force").Run()
	_ = exec.Command("git", "branch", "-D", "a1b2-duplicate").Run()
}

func TestRemove(t *testing.T) {
	repoDir := setupGitRepo(t)

	// Change to repo directory
	originalDir, _ := os.Getwd()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo dir: %v", err)
	}
	defer os.Chdir(originalDir)

	// Get current branch
	branchOut, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	parentBranch := string(branchOut)
	parentBranch = parentBranch[:len(parentBranch)-1]

	worktreesPath := filepath.Join(repoDir, "..", "worktrees")

	// Create a worktree first
	worktreePath, err := Create(worktreesPath, "a1b2-to-remove", parentBranch)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify it exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Fatal("Worktree should exist after creation")
	}

	// Remove it
	if err := Remove(worktreesPath, "a1b2-to-remove"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// Verify it's gone
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Error("Worktree should not exist after removal")
	}
}

func TestRemove_NotFound(t *testing.T) {
	dir := t.TempDir()
	worktreesPath := filepath.Join(dir, "worktrees")

	err := Remove(worktreesPath, "nonexistent")
	if err == nil {
		t.Fatal("Expected error when removing non-existent worktree")
	}

	wtErr, ok := err.(*WorktreeError)
	if !ok {
		t.Fatalf("Expected WorktreeError, got %T", err)
	}

	if wtErr.Code != "NOOK_NOT_FOUND" {
		t.Errorf("Error code = %q, want %q", wtErr.Code, "NOOK_NOT_FOUND")
	}
}
