package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadBoards_NotExists(t *testing.T) {
	dir := t.TempDir()
	nooksPath := filepath.Join(dir, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks dir: %v", err)
	}

	// No boards.yaml file - should return empty BoardsFile, not error
	boards, err := LoadBoards(nooksPath)
	if err != nil {
		t.Fatalf("LoadBoards() error = %v", err)
	}

	if boards.Boards == nil {
		t.Error("Boards map should not be nil")
	}

	if len(boards.Boards) != 0 {
		t.Errorf("Expected 0 boards, got %d", len(boards.Boards))
	}
}

func TestLoadBoards_Valid(t *testing.T) {
	dir := t.TempDir()
	nooksPath := filepath.Join(dir, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks dir: %v", err)
	}

	// Create boards.yaml with data
	boardsContent := `a8ca-test-nook:
  - ts: "2026-01-18T10:30:00Z"
    content: |
      First pin content
  - ts: "2026-01-18T11:00:00Z"
    content: |
      Second pin content
b7c3-another-nook:
  - ts: "2026-01-18T12:00:00Z"
    content: |
      Another nook's pin
`
	boardsPath := filepath.Join(nooksPath, "boards.yaml")
	if err := os.WriteFile(boardsPath, []byte(boardsContent), 0644); err != nil {
		t.Fatalf("Failed to write boards.yaml: %v", err)
	}

	boards, err := LoadBoards(nooksPath)
	if err != nil {
		t.Fatalf("LoadBoards() error = %v", err)
	}

	if len(boards.Boards) != 2 {
		t.Errorf("Expected 2 nooks with boards, got %d", len(boards.Boards))
	}

	pins, ok := boards.Boards["a8ca-test-nook"]
	if !ok {
		t.Fatal("Expected a8ca-test-nook to exist")
	}

	if len(pins) != 2 {
		t.Errorf("Expected 2 pins, got %d", len(pins))
	}

	if pins[0].Ts != "2026-01-18T10:30:00Z" {
		t.Errorf("Pin ts = %q, want %q", pins[0].Ts, "2026-01-18T10:30:00Z")
	}

	if !strings.Contains(pins[0].Content, "First pin content") {
		t.Errorf("Pin content = %q, expected to contain 'First pin content'", pins[0].Content)
	}
}

func TestLoadBoards_Invalid(t *testing.T) {
	dir := t.TempDir()
	nooksPath := filepath.Join(dir, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks dir: %v", err)
	}

	// Create invalid boards.yaml
	invalidContent := `this is not: [valid yaml: for boards`
	boardsPath := filepath.Join(nooksPath, "boards.yaml")
	if err := os.WriteFile(boardsPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to write boards.yaml: %v", err)
	}

	_, err := LoadBoards(nooksPath)
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}

	boardErr, ok := err.(*BoardError)
	if !ok {
		t.Fatalf("Expected BoardError, got %T", err)
	}

	if boardErr.Code != "BOARD_PARSE_FAILED" {
		t.Errorf("Error code = %q, want %q", boardErr.Code, "BOARD_PARSE_FAILED")
	}
}

func TestSaveBoards_New(t *testing.T) {
	dir := t.TempDir()
	nooksPath := filepath.Join(dir, "nooks")
	// Don't create nooks dir - SaveBoards should create it

	boards := &BoardsFile{
		Boards: map[string][]Pin{
			"a8ca-test-nook": {
				{Ts: "2026-01-18T10:30:00Z", Content: "Test content"},
			},
		},
	}

	if err := SaveBoards(nooksPath, boards); err != nil {
		t.Fatalf("SaveBoards() error = %v", err)
	}

	// Verify file was created
	boardsPath := filepath.Join(nooksPath, "boards.yaml")
	if _, err := os.Stat(boardsPath); os.IsNotExist(err) {
		t.Fatal("boards.yaml was not created")
	}

	// Load it back and verify
	loaded, err := LoadBoards(nooksPath)
	if err != nil {
		t.Fatalf("LoadBoards() error = %v", err)
	}

	if len(loaded.Boards) != 1 {
		t.Errorf("Expected 1 board, got %d", len(loaded.Boards))
	}

	pins, ok := loaded.Boards["a8ca-test-nook"]
	if !ok {
		t.Fatal("Expected a8ca-test-nook to exist")
	}

	if len(pins) != 1 {
		t.Errorf("Expected 1 pin, got %d", len(pins))
	}

	if pins[0].Content != "Test content" {
		t.Errorf("Pin content = %q, want %q", pins[0].Content, "Test content")
	}
}

func TestSaveBoards_Existing(t *testing.T) {
	dir := t.TempDir()
	nooksPath := filepath.Join(dir, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks dir: %v", err)
	}

	// Create initial boards.yaml
	initialContent := `a8ca-test-nook:
  - ts: "2026-01-18T10:30:00Z"
    content: Initial content
`
	boardsPath := filepath.Join(nooksPath, "boards.yaml")
	if err := os.WriteFile(boardsPath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write initial boards.yaml: %v", err)
	}

	// Save updated boards
	boards := &BoardsFile{
		Boards: map[string][]Pin{
			"a8ca-test-nook": {
				{Ts: "2026-01-18T10:30:00Z", Content: "Initial content"},
				{Ts: "2026-01-18T11:00:00Z", Content: "New content"},
			},
		},
	}

	if err := SaveBoards(nooksPath, boards); err != nil {
		t.Fatalf("SaveBoards() error = %v", err)
	}

	// Load and verify
	loaded, err := LoadBoards(nooksPath)
	if err != nil {
		t.Fatalf("LoadBoards() error = %v", err)
	}

	pins := loaded.Boards["a8ca-test-nook"]
	if len(pins) != 2 {
		t.Errorf("Expected 2 pins, got %d", len(pins))
	}
}

func TestSaveBoards_NoTempFileLeftOnSuccess(t *testing.T) {
	dir := t.TempDir()
	nooksPath := filepath.Join(dir, "nooks")
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		t.Fatalf("Failed to create nooks dir: %v", err)
	}

	boards := &BoardsFile{Boards: map[string][]Pin{}}

	if err := SaveBoards(nooksPath, boards); err != nil {
		t.Fatalf("SaveBoards() error = %v", err)
	}

	// Check that no temp files remain
	entries, err := os.ReadDir(nooksPath)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	for _, entry := range entries {
		if entry.Name() != "boards.yaml" {
			t.Errorf("Unexpected file in nooks: %s", entry.Name())
		}
	}
}

func TestAddPin(t *testing.T) {
	boards := &BoardsFile{Boards: map[string][]Pin{}}

	beforeTime := time.Now().UTC().Truncate(time.Second)
	ts := AddPin(boards, "a8ca-test-nook", "Test pin content")
	afterTime := time.Now().UTC().Add(time.Second).Truncate(time.Second)

	// Verify pin was added
	pins, ok := boards.Boards["a8ca-test-nook"]
	if !ok {
		t.Fatal("Expected a8ca-test-nook to exist")
	}

	if len(pins) != 1 {
		t.Errorf("Expected 1 pin, got %d", len(pins))
	}

	if pins[0].Content != "Test pin content" {
		t.Errorf("Pin content = %q, want %q", pins[0].Content, "Test pin content")
	}

	// Verify timestamp is RFC3339 and within expected range
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Errorf("Timestamp %q is not RFC3339: %v", ts, err)
	}

	// RFC3339 truncates to seconds, so we compare truncated times
	if parsedTime.Before(beforeTime) || parsedTime.After(afterTime) {
		t.Errorf("Timestamp %v not in expected range [%v, %v]", parsedTime, beforeTime, afterTime)
	}
}

func TestAddPin_AppendsToExisting(t *testing.T) {
	boards := &BoardsFile{
		Boards: map[string][]Pin{
			"a8ca-test-nook": {
				{Ts: "2026-01-18T10:30:00Z", Content: "First pin"},
			},
		},
	}

	AddPin(boards, "a8ca-test-nook", "Second pin")

	pins := boards.Boards["a8ca-test-nook"]
	if len(pins) != 2 {
		t.Errorf("Expected 2 pins, got %d", len(pins))
	}

	if pins[0].Content != "First pin" {
		t.Errorf("First pin content = %q, want %q", pins[0].Content, "First pin")
	}

	if pins[1].Content != "Second pin" {
		t.Errorf("Second pin content = %q, want %q", pins[1].Content, "Second pin")
	}
}

func TestAddPin_NilBoards(t *testing.T) {
	boards := &BoardsFile{Boards: nil}

	AddPin(boards, "a8ca-test-nook", "Test content")

	if boards.Boards == nil {
		t.Fatal("Boards map should be initialized")
	}

	if len(boards.Boards["a8ca-test-nook"]) != 1 {
		t.Error("Pin should be added even when boards map was nil")
	}
}

func TestGetBoard_Exists(t *testing.T) {
	boards := &BoardsFile{
		Boards: map[string][]Pin{
			"a8ca-test-nook": {
				{Ts: "2026-01-18T10:30:00Z", Content: "First pin"},
				{Ts: "2026-01-18T11:00:00Z", Content: "Second pin"},
			},
			"b7c3-another-nook": {
				{Ts: "2026-01-18T12:00:00Z", Content: "Other nook pin"},
			},
		},
	}

	pins := GetBoard(boards, "a8ca-test-nook")
	if len(pins) != 2 {
		t.Errorf("Expected 2 pins, got %d", len(pins))
	}

	if pins[0].Content != "First pin" {
		t.Errorf("First pin content = %q, want %q", pins[0].Content, "First pin")
	}
}

func TestGetBoard_NotExists(t *testing.T) {
	boards := &BoardsFile{
		Boards: map[string][]Pin{
			"a8ca-test-nook": {
				{Ts: "2026-01-18T10:30:00Z", Content: "A pin"},
			},
		},
	}

	pins := GetBoard(boards, "nonexistent-nook")
	if pins == nil {
		t.Error("Pins should not be nil, expected empty slice")
	}

	if len(pins) != 0 {
		t.Errorf("Expected 0 pins, got %d", len(pins))
	}
}

func TestGetBoard_NilBoards(t *testing.T) {
	boards := &BoardsFile{Boards: nil}

	pins := GetBoard(boards, "any-nook")
	if pins == nil {
		t.Error("Pins should not be nil, expected empty slice")
	}

	if len(pins) != 0 {
		t.Errorf("Expected 0 pins, got %d", len(pins))
	}
}
