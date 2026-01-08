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

// runGitCommand runs a git command in the specified directory
func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Git command failed: git %v: %v", args, err)
	}
}

func TestListCommand_NotInitialized(t *testing.T) {
	// Create temp directory without .treehouse folder
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Capture output
	var buf bytes.Buffer
	SetOutput(&buf)
	defer SetOutput(os.Stdout)

	// Run list command
	exitCode := Execute([]string{"list"})

	// Verify exit code
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	// Verify JSON error output
	var response output.Response
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if response.Success {
		t.Error("Expected success=false")
	}

	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}

	if response.Error.Code != "INIT_NOT_FOUND" {
		t.Errorf("Expected error code INIT_NOT_FOUND, got %s", response.Error.Code)
	}
}

func TestListCommand_EmptyDecks(t *testing.T) {
	// Create temp directory with .treehouse but no decks
	tempDir := t.TempDir()
	treehousePath := filepath.Join(tempDir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create empty decks.yaml
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte("decks: {}\n"), 0644); err != nil {
		t.Fatalf("Failed to create decks.yaml: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize git repo (required for GetCurrentBranch and GetCurrentCommit)
	runGitCommand(t, tempDir, "init")
	runGitCommand(t, tempDir, "config", "user.name", "Test User")
	runGitCommand(t, tempDir, "config", "user.email", "test@example.com")
	runGitCommand(t, tempDir, "commit", "--allow-empty", "-m", "Initial commit")

	// Capture output
	var buf bytes.Buffer
	SetOutput(&buf)
	defer SetOutput(os.Stdout)

	// Run list command
	exitCode := Execute([]string{"list"})

	// Verify exit code
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	// Verify JSON output
	var response output.Response
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, buf.String())
	}

	if !response.Success {
		t.Errorf("Expected success=true\nError: %v", response.Error)
	}

	// Parse data
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	// Check base info exists
	if _, exists := data["base"]; !exists {
		t.Error("Expected base field")
	}

	// Check current_nook is null (we're in base repo)
	if data["current_nook"] != nil {
		t.Errorf("Expected current_nook=null, got %v", data["current_nook"])
	}

	// Check decks is empty array
	decks, ok := data["decks"].([]interface{})
	if !ok {
		t.Fatal("Expected decks to be an array")
	}
	if len(decks) != 0 {
		t.Errorf("Expected empty decks array, got %d decks", len(decks))
	}
}

func TestListCommand_WithDecks(t *testing.T) {
	// Create temp directory with .treehouse and decks
	tempDir := t.TempDir()
	treehousePath := filepath.Join(tempDir, ".treehouse")
	nooksPath := filepath.Join(treehousePath, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks directory: %v", err)
	}

	// Create decks.yaml with one deck and one nook
	decksContent := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-test-nook:
        parent: main
        created: "2026-01-08"
`
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte(decksContent), 0644); err != nil {
		t.Fatalf("Failed to create decks.yaml: %v", err)
	}

	// Create the nook worktree folder (so it's not orphan)
	nookPath := filepath.Join(nooksPath, "a1b2-test-nook")
	if err := os.MkdirAll(nookPath, 0755); err != nil {
		t.Fatalf("Failed to create nook folder: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize git repo
	runGitCommand(t, tempDir, "init")
	runGitCommand(t, tempDir, "config", "user.name", "Test User")
	runGitCommand(t, tempDir, "config", "user.email", "test@example.com")
	runGitCommand(t, tempDir, "commit", "--allow-empty", "-m", "Initial commit")

	// Capture output
	var buf bytes.Buffer
	SetOutput(&buf)
	defer SetOutput(os.Stdout)

	// Run list command
	exitCode := Execute([]string{"list"})

	// Verify exit code
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	// Verify JSON output
	var response output.Response
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, buf.String())
	}

	if !response.Success {
		t.Errorf("Expected success=true\nError: %v", response.Error)
	}

	// Parse data
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	// Check decks exists and has one deck
	decks, ok := data["decks"].([]interface{})
	if !ok {
		t.Fatal("Expected decks to be an array")
	}
	if len(decks) != 1 {
		t.Errorf("Expected 1 deck, got %d", len(decks))
	}

	// Verify deck structure
	deck0, ok := decks[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected deck to be a map")
	}

	if deck0["id"] != "dk-a1b2" {
		t.Errorf("Expected deck id=dk-a1b2, got %v", deck0["id"])
	}

	// Verify nook
	nooks, ok := deck0["nooks"].([]interface{})
	if !ok {
		t.Fatal("Expected nooks to be an array")
	}
	if len(nooks) != 1 {
		t.Errorf("Expected 1 nook, got %d", len(nooks))
	}

	nook0, ok := nooks[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected nook to be a map")
	}

	if nook0["id"] != "a1b2-test-nook" {
		t.Errorf("Expected nook id=a1b2-test-nook, got %v", nook0["id"])
	}
	if nook0["status"] != "inactive" {
		t.Errorf("Expected status=inactive, got %v", nook0["status"])
	}
}

func TestListCommand_OrphanNook(t *testing.T) {
	// Create temp directory with .treehouse and deck with orphan nook
	tempDir := t.TempDir()
	treehousePath := filepath.Join(tempDir, ".treehouse")
	nooksPath := filepath.Join(treehousePath, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks directory: %v", err)
	}

	// Create decks.yaml with a nook that doesn't have a worktree folder
	decksContent := `decks:
  dk-c3d4:
    created: "2026-01-08"
    nooks:
      c3d4-orphan-nook:
        parent: main
        created: "2026-01-08"
`
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte(decksContent), 0644); err != nil {
		t.Fatalf("Failed to create decks.yaml: %v", err)
	}

	// NOTE: We do NOT create the nook folder - making it orphan

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize git repo
	runGitCommand(t, tempDir, "init")
	runGitCommand(t, tempDir, "config", "user.name", "Test User")
	runGitCommand(t, tempDir, "config", "user.email", "test@example.com")
	runGitCommand(t, tempDir, "commit", "--allow-empty", "-m", "Initial commit")

	// Capture output
	var buf bytes.Buffer
	SetOutput(&buf)
	defer SetOutput(os.Stdout)

	// Run list command
	exitCode := Execute([]string{"list"})

	// Verify exit code - should still succeed
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	// Verify JSON output
	var response output.Response
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, buf.String())
	}

	if !response.Success {
		t.Errorf("Expected success=true\nError: %v", response.Error)
	}

	// Parse data
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	// Check decks
	decks, ok := data["decks"].([]interface{})
	if !ok {
		t.Fatal("Expected decks to be an array")
	}
	if len(decks) != 1 {
		t.Errorf("Expected 1 deck, got %d", len(decks))
	}

	// Verify nook is marked as orphan
	deck0, _ := decks[0].(map[string]interface{})
	nooks, _ := deck0["nooks"].([]interface{})
	nook0, _ := nooks[0].(map[string]interface{})

	if nook0["status"] != "orphan" {
		t.Errorf("Expected status=orphan, got %v", nook0["status"])
	}
}
