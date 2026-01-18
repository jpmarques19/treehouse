package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jpmarques19/treehouse/internal/output"
	"github.com/jpmarques19/treehouse/internal/testutil"
)

func TestBoard_Success(t *testing.T) {
	_, nookID := setupInitializedRepoWithNook(t)

	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// First add a pin
	exitCode := Execute([]string{"pin", "Test pin content"})
	if exitCode != 0 {
		t.Fatalf("pin failed: %s", buf.String())
	}

	// Now get the board
	buf.Reset()
	exitCode = Execute([]string{"board"})

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

	pins, ok := data["pins"].([]interface{})
	if !ok {
		t.Fatalf("pins should be an array, got %T", data["pins"])
	}

	if len(pins) != 1 {
		t.Errorf("expected 1 pin, got %d", len(pins))
	}

	// Verify pin structure
	pin := pins[0].(map[string]interface{})
	if pin["content"] != "Test pin content" {
		t.Errorf("pin content = %v, want %v", pin["content"], "Test pin content")
	}

	if pin["ts"] == nil || pin["ts"] == "" {
		t.Error("pin ts should not be empty")
	}
}

func TestBoard_Empty(t *testing.T) {
	_, nookID := setupInitializedRepoWithNook(t)

	// Capture output - no pins added, just get board
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"board"})

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
	data := resp.Data.(map[string]interface{})

	if data["nook"] != nookID {
		t.Errorf("nook = %v, want %v", data["nook"], nookID)
	}

	pins, ok := data["pins"].([]interface{})
	if !ok {
		t.Fatalf("pins should be an array, got %T", data["pins"])
	}

	// Empty pins array, not error
	if len(pins) != 0 {
		t.Errorf("expected 0 pins, got %d", len(pins))
	}
}

func TestBoard_NotInNook(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Initialize treehouse
	var buf bytes.Buffer
	output.SetWriter(&buf)
	Execute([]string{"init"})

	// Try to view board from base repo (not in nook)
	buf.Reset()
	exitCode := Execute([]string{"board"})
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

	if resp.Error == nil || resp.Error.Code != "BOARD_NOT_IN_NOOK" {
		t.Errorf("expected error code BOARD_NOT_IN_NOOK, got %v", resp.Error)
	}
}

func TestBoard_NotInitialized(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Don't initialize treehouse

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"board"})

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

func TestBoard_OnlyShowsCurrentNookPins(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Initialize treehouse
	var buf bytes.Buffer
	output.SetWriter(&buf)
	Execute([]string{"init"})

	// Create first nook and add pin
	buf.Reset()
	Execute([]string{"fork", "first-nook"})
	var resp output.Response
	json.Unmarshal(buf.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	firstNookID := data["nook_id"].(string)

	treehousePath := filepath.Join(dir, ".treehouse")
	firstNookPath := filepath.Join(treehousePath, "nooks", firstNookID)

	if err := os.Chdir(firstNookPath); err != nil {
		t.Fatalf("failed to chdir to first nook: %v", err)
	}

	buf.Reset()
	Execute([]string{"pin", "First nook pin"})

	// Go back to base and create second nook
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir back to dir: %v", err)
	}

	buf.Reset()
	Execute([]string{"fork", "second-nook"})
	json.Unmarshal(buf.Bytes(), &resp)
	data = resp.Data.(map[string]interface{})
	secondNookID := data["nook_id"].(string)

	secondNookPath := filepath.Join(treehousePath, "nooks", secondNookID)

	if err := os.Chdir(secondNookPath); err != nil {
		t.Fatalf("failed to chdir to second nook: %v", err)
	}

	buf.Reset()
	Execute([]string{"pin", "Second nook pin"})

	// Now get board in second nook - should only have its own pin
	buf.Reset()
	exitCode := Execute([]string{"board"})
	output.SetWriter(os.Stdout)

	if exitCode != 0 {
		t.Fatalf("board failed: %s", buf.String())
	}

	json.Unmarshal(buf.Bytes(), &resp)
	data = resp.Data.(map[string]interface{})

	if data["nook"] != secondNookID {
		t.Errorf("nook = %v, want %v", data["nook"], secondNookID)
	}

	pins := data["pins"].([]interface{})
	if len(pins) != 1 {
		t.Errorf("expected 1 pin in second nook, got %d", len(pins))
	}

	pin := pins[0].(map[string]interface{})
	if pin["content"] != "Second nook pin" {
		t.Errorf("expected second nook's pin, got %v", pin["content"])
	}
}
