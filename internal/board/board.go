package board

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// BoardsFile represents the structure of boards.yaml
// Maps nook-id to a list of pins
type BoardsFile struct {
	Boards map[string][]Pin `yaml:",inline"`
}

// Pin represents a single pin entry on a nook's board
type Pin struct {
	Ts      string `yaml:"ts" json:"ts"`
	Content string `yaml:"content" json:"content"`
}

// BoardError represents a board operation error
type BoardError struct {
	Code    string
	Message string
}

func (e *BoardError) Error() string {
	return e.Message
}

// LoadBoards reads and parses the boards.yaml file
// Returns empty BoardsFile if file doesn't exist (not an error)
func LoadBoards(nooksPath string) (*BoardsFile, error) {
	boardsPath := filepath.Join(nooksPath, "boards.yaml")

	data, err := os.ReadFile(boardsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty BoardsFile, not an error
			return &BoardsFile{
				Boards: make(map[string][]Pin),
			}, nil
		}
		return nil, &BoardError{
			Code:    "BOARD_READ_FAILED",
			Message: fmt.Sprintf("Failed to read boards.yaml: %v", err),
		}
	}

	var boards map[string][]Pin
	if err := yaml.Unmarshal(data, &boards); err != nil {
		return nil, &BoardError{
			Code:    "BOARD_PARSE_FAILED",
			Message: fmt.Sprintf("Failed to parse boards.yaml: %v", err),
		}
	}

	// Initialize empty map if nil
	if boards == nil {
		boards = make(map[string][]Pin)
	}

	return &BoardsFile{Boards: boards}, nil
}

// SaveBoards writes the boards file atomically (temp file + rename)
func SaveBoards(nooksPath string, boards *BoardsFile) error {
	boardsPath := filepath.Join(nooksPath, "boards.yaml")

	// Marshal to YAML (just the map, not the wrapper struct)
	data, err := yaml.Marshal(boards.Boards)
	if err != nil {
		return &BoardError{
			Code:    "BOARD_MARSHAL_FAILED",
			Message: fmt.Sprintf("Failed to marshal boards: %v", err),
		}
	}

	// Ensure nooks directory exists
	if err := os.MkdirAll(nooksPath, 0755); err != nil {
		return &BoardError{
			Code:    "BOARD_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to create nooks directory: %v", err),
		}
	}

	// Create temp file in same directory for atomic rename
	tempFile, err := os.CreateTemp(nooksPath, "boards.yaml.tmp.*")
	if err != nil {
		return &BoardError{
			Code:    "BOARD_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to create temp file: %v", err),
		}
	}
	tempPath := tempFile.Name()

	// Ensure temp file cleanup on error
	defer func() {
		if tempPath != "" {
			os.Remove(tempPath)
		}
	}()

	// Write data to temp file
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return &BoardError{
			Code:    "BOARD_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to write temp file: %v", err),
		}
	}

	// Sync to disk
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return &BoardError{
			Code:    "BOARD_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to sync temp file: %v", err),
		}
	}

	// Close before rename
	if err := tempFile.Close(); err != nil {
		return &BoardError{
			Code:    "BOARD_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to close temp file: %v", err),
		}
	}

	// Atomic rename
	if err := os.Rename(tempPath, boardsPath); err != nil {
		return &BoardError{
			Code:    "BOARD_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to rename temp file: %v", err),
		}
	}

	// Clear tempPath to prevent cleanup of the renamed file
	tempPath = ""

	return nil
}

// AddPin appends a new pin to the specified nook's board
// Returns the timestamp of the created pin
func AddPin(boards *BoardsFile, nookID string, content string) string {
	// Initialize boards map if nil
	if boards.Boards == nil {
		boards.Boards = make(map[string][]Pin)
	}

	// Generate RFC3339 timestamp
	ts := time.Now().UTC().Format(time.RFC3339)

	// Append pin to nook's board
	boards.Boards[nookID] = append(boards.Boards[nookID], Pin{
		Ts:      ts,
		Content: content,
	})

	return ts
}

// GetBoard returns all pins for a specific nook
// Returns empty slice if nook has no pins
func GetBoard(boards *BoardsFile, nookID string) []Pin {
	if boards.Boards == nil {
		return []Pin{}
	}

	pins, exists := boards.Boards[nookID]
	if !exists {
		return []Pin{}
	}

	return pins
}
