package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/testutil"
	"github.com/jpmarques19/treehouse/internal/worktree"
)

// setupInitializedRepoWithNook creates a repo with treehouse initialized and a nook worktree.
// Changes directory into the nook worktree.
// Returns the base repo path and nook ID.
func setupInitializedRepoWithNook(t *testing.T) (string, string) {
	t.Helper()
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Initialize treehouse
	var buf bytes.Buffer
	output.SetWriter(&buf)
	exitCode := Execute([]string{"init"})

	if exitCode != 0 {
		output.SetWriter(os.Stdout)
		t.Fatalf("failed to initialize treehouse: %s", buf.String())
	}

	// Create a nook worktree using fork command
	buf.Reset()
	exitCode = Execute([]string{"fork", "test-nook"})

	if exitCode != 0 {
		output.SetWriter(os.Stdout)
		t.Fatalf("failed to fork nook: %s", buf.String())
	}

	// Parse fork response to get nook ID
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		output.SetWriter(os.Stdout)
		t.Fatalf("failed to parse fork response: %v", err)
	}

	// Get nook ID from response
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		output.SetWriter(os.Stdout)
		t.Fatalf("unexpected response format: %v", resp.Data)
	}
	nookID := data["nook_id"].(string)

	// Get the nook path
	treehousePath := filepath.Join(dir, ".treehouse")
	worktreesPath := filepath.Join(treehousePath, "nooks")
	nookPath := filepath.Join(worktreesPath, nookID)

	output.SetWriter(os.Stdout)

	// Change to nook directory
	if err := os.Chdir(nookPath); err != nil {
		t.Fatalf("failed to chdir to nook: %v", err)
	}

	return dir, nookID
}

func TestPin_Success(t *testing.T) {
	dir, nookID := setupInitializedRepoWithNook(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"pin", "Test pin content"})

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	// Parse JSON response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v\nOutput: %s", err, buf.String())
	}

	if !resp.Success {
		t.Errorf("expected success=true, got success=false")
	}

	// Verify response data
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected data format: %v", resp.Data)
	}

	if data["nook"] != nookID {
		t.Errorf("nook = %v, want %v", data["nook"], nookID)
	}

	if data["pins_count"].(float64) != 1 {
		t.Errorf("pins_count = %v, want 1", data["pins_count"])
	}

	if data["ts"] == nil || data["ts"] == "" {
		t.Error("ts should not be empty")
	}

	// Verify boards.yaml was created
	boardsPath := filepath.Join(dir, ".treehouse", "nooks", "boards.yaml")
	if _, err := os.Stat(boardsPath); os.IsNotExist(err) {
		t.Error("boards.yaml was not created")
	}
}

func TestPin_NotInNook(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Initialize treehouse
	var buf bytes.Buffer
	output.SetWriter(&buf)
	Execute([]string{"init"})

	// Try to pin from base repo (not in nook)
	buf.Reset()
	exitCode := Execute([]string{"pin", "Test content"})
	output.SetWriter(os.Stdout)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	// Parse JSON response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false for base repo")
	}

	if resp.Error == nil || resp.Error.Code != "PIN_NOT_IN_NOOK" {
		t.Errorf("expected error code PIN_NOT_IN_NOOK, got %v", resp.Error)
	}
}

func TestPin_NotInitialized(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Don't initialize treehouse

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"pin", "Test content"})

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false for uninitialized treehouse")
	}
}

func TestPin_NoArgs(t *testing.T) {
	_, _ = setupInitializedRepoWithNook(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Pin with no content argument
	exitCode := Execute([]string{"pin"})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}
}

func TestPin_MultipleInSameNook(t *testing.T) {
	dir, nookID := setupInitializedRepoWithNook(t)

	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Add first pin
	exitCode := Execute([]string{"pin", "First pin"})
	if exitCode != 0 {
		t.Fatalf("first pin failed: %s", buf.String())
	}

	// Add second pin
	buf.Reset()
	exitCode = Execute([]string{"pin", "Second pin"})
	if exitCode != 0 {
		t.Fatalf("second pin failed: %s", buf.String())
	}

	// Verify response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	data := resp.Data.(map[string]interface{})
	if data["pins_count"].(float64) != 2 {
		t.Errorf("pins_count = %v, want 2", data["pins_count"])
	}

	// Verify boards.yaml has both pins
	boardsPath := filepath.Join(dir, ".treehouse", "nooks", "boards.yaml")
	content, err := os.ReadFile(boardsPath)
	if err != nil {
		t.Fatalf("failed to read boards.yaml: %v", err)
	}

	if !bytes.Contains(content, []byte(nookID)) {
		t.Error("boards.yaml should contain nook ID")
	}

	if !bytes.Contains(content, []byte("First pin")) {
		t.Error("boards.yaml should contain first pin content")
	}

	if !bytes.Contains(content, []byte("Second pin")) {
		t.Error("boards.yaml should contain second pin content")
	}
}

// Helper to verify GetParentInfo detects nook correctly
func TestWorktreeDetection(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Initialize treehouse
	var buf bytes.Buffer
	output.SetWriter(&buf)
	Execute([]string{"init"})
	output.SetWriter(os.Stdout)

	treehousePath := filepath.Join(dir, ".treehouse")

	// In base repo
	info, err := worktree.GetParentInfo(treehousePath)
	if err != nil {
		t.Fatalf("GetParentInfo() error = %v", err)
	}
	if info.IsNook {
		t.Error("expected IsNook=false in base repo")
	}

	// Create a nook
	buf.Reset()
	output.SetWriter(&buf)
	Execute([]string{"fork", "test"})

	var resp output.Response
	json.Unmarshal(buf.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	nookID := data["nook_id"].(string)
	output.SetWriter(os.Stdout)

	// Change to nook
	nookPath := filepath.Join(treehousePath, "nooks", nookID)
	if err := os.Chdir(nookPath); err != nil {
		t.Fatalf("failed to chdir to nook: %v", err)
	}

	// Should detect as nook
	info, err = worktree.GetParentInfo(treehousePath)
	if err != nil {
		t.Fatalf("GetParentInfo() error = %v", err)
	}
	if !info.IsNook {
		t.Error("expected IsNook=true in nook worktree")
	}
	if info.CurrentNook != nookID {
		t.Errorf("CurrentNook = %v, want %v", info.CurrentNook, nookID)
	}
}
