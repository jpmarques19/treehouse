package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jpmarques19/treehouse/internal/output"
)

// setupForkTestRepo creates a git repo with treehouse initialized for testing
func setupForkTestRepo(t *testing.T) string {
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

	// Create initial commit
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

	// Initialize treehouse
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create empty decks.yaml
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte("decks: {}\n"), 0644); err != nil {
		t.Fatalf("Failed to create decks.yaml: %v", err)
	}

	return dir
}

func TestForkCommand_MissingName(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	defer output.SetWriter(os.Stdout)

	// Execute fork without name
	rootCmd.SetArgs([]string{"fork"})
	_ = rootCmd.Execute()

	// Parse response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}

	if resp.Error == nil {
		t.Fatal("Expected error to be present")
	}

	if resp.Error.Code != "NOOK_NAME_REQUIRED" {
		t.Errorf("Error code = %q, want %q", resp.Error.Code, "NOOK_NAME_REQUIRED")
	}
}

func TestForkCommand_InvalidName(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	defer output.SetWriter(os.Stdout)

	// Execute fork with invalid name (only special chars)
	rootCmd.SetArgs([]string{"fork", "!@#$%"})
	_ = rootCmd.Execute()

	// Parse response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}

	if resp.Error == nil {
		t.Fatal("Expected error to be present")
	}

	if resp.Error.Code != "NOOK_NAME_INVALID" {
		t.Errorf("Error code = %q, want %q", resp.Error.Code, "NOOK_NAME_INVALID")
	}
}

func TestForkCommand_NotInitialized(t *testing.T) {
	dir := t.TempDir()

	// Initialize git but not treehouse
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = dir
	_ = cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	_ = cmd.Run()

	// Create initial commit
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dir
	_ = cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	_ = cmd.Run()

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	defer output.SetWriter(os.Stdout)

	// Execute fork
	rootCmd.SetArgs([]string{"fork", "test-nook"})
	_ = rootCmd.Execute()

	// Parse response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, buf.String())
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}

	if resp.Error == nil {
		t.Fatal("Expected error to be present")
	}

	if resp.Error.Code != "INIT_NOT_FOUND" {
		t.Errorf("Error code = %q, want %q", resp.Error.Code, "INIT_NOT_FOUND")
	}
}

func TestForkCommand_Success(t *testing.T) {
	repoDir := setupForkTestRepo(t)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	defer output.SetWriter(os.Stdout)

	// Execute fork
	rootCmd.SetArgs([]string{"fork", "auth-spike"})
	_ = rootCmd.Execute()

	// Parse response
	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			NookID   string `json:"nook_id"`
			DeckID   string `json:"deck_id"`
			Parent   string `json:"parent"`
			Worktree string `json:"worktree"`
		} `json:"data"`
		Error *output.ErrorInfo `json:"error"`
	}

	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, buf.String())
	}

	if !resp.Success {
		if resp.Error != nil {
			t.Fatalf("Fork failed with error: %s - %s", resp.Error.Code, resp.Error.Message)
		}
		t.Fatal("Fork failed without error message")
	}

	// Verify nook ID format (4-char hash + hyphen + sanitized name)
	if len(resp.Data.NookID) < 6 {
		t.Errorf("NookID too short: %s", resp.Data.NookID)
	}

	// Verify deck ID format
	if len(resp.Data.DeckID) < 3 || resp.Data.DeckID[:3] != "dk-" {
		t.Errorf("DeckID format invalid: %s", resp.Data.DeckID)
	}

	// Verify worktree was created
	if _, err := os.Stat(resp.Data.Worktree); os.IsNotExist(err) {
		t.Errorf("Worktree path does not exist: %s", resp.Data.Worktree)
	}

	// Clean up worktree
	_ = exec.Command("git", "worktree", "remove", resp.Data.Worktree, "--force").Run()
	_ = exec.Command("git", "branch", "-D", resp.Data.NookID).Run()
}
