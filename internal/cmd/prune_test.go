package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPruneCommand(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		setupFunc         func(t *testing.T, tmpDir string)
		expectedCode      int
		expectedError     string
		checkSuccess      bool
		checkPruned       []string
		checkDecksRemoved []string
	}{
		{
			name:          "not initialized",
			args:          []string{"prune"},
			expectedCode:  1,
			expectedError: "INIT_NOT_FOUND",
		},
		{
			name:         "no orphans",
			args:         []string{"prune"},
			expectedCode: 0,
			checkSuccess: true,
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with valid nook
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				os.MkdirAll(nooksPath, 0755)

				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-valid-nook:
        parent: main
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)
				// Create the worktree folder so it's not orphan
				os.MkdirAll(filepath.Join(nooksPath, "a1b2-valid-nook"), 0755)
			},
		},
		{
			name:         "prune single orphan",
			args:         []string{"prune"},
			expectedCode: 0,
			checkSuccess: true,
			checkPruned:  []string{"a1b2-orphan-nook"},
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with orphan nook (no worktree folder)
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				os.MkdirAll(nooksPath, 0755)

				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-orphan-nook:
        parent: main
        created: "2026-01-08"
      c3d4-valid-nook:
        parent: main
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)
				// Only create folder for valid nook, orphan has no folder
				os.MkdirAll(filepath.Join(nooksPath, "c3d4-valid-nook"), 0755)
			},
		},
		{
			name:              "prune orphan and remove empty deck",
			args:              []string{"prune"},
			expectedCode:      0,
			checkSuccess:      true,
			checkPruned:       []string{"a1b2-orphan-nook"},
			checkDecksRemoved: []string{"dk-a1b2"},
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with only orphan nook
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				os.MkdirAll(nooksPath, 0755)

				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-orphan-nook:
        parent: main
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)
				// No worktree folder created - it's an orphan
			},
		},
		{
			name:         "cleanup crew memory files on prune",
			args:         []string{"prune"},
			expectedCode: 0,
			checkSuccess: true,
			setupFunc: func(t *testing.T, tmpDir string) {
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				crewPath := filepath.Join(treehousePath, "crew", "spruce")
				os.MkdirAll(nooksPath, 0755)
				os.MkdirAll(filepath.Join(crewPath, "memories"), 0755)
				os.MkdirAll(filepath.Join(crewPath, "sessions"), 0755)

				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-orphan:
        parent: main
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)

				// Create memory files that should be deleted on prune
				os.WriteFile(filepath.Join(crewPath, "memories", "a1b2-orphan.md"), []byte("# Memory"), 0644)
				os.WriteFile(filepath.Join(crewPath, "sessions", "a1b2-orphan.md"), []byte("# Session"), 0644)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Initialize as real git repo
			initGitRepo(t, tmpDir)

			// Run setup function if provided
			if tt.setupFunc != nil {
				tt.setupFunc(t, tmpDir)
			}

			// Change to temp directory
			oldWd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldWd)

			// Capture output
			var buf bytes.Buffer
			SetOutput(&buf)

			// Execute command
			exitCode := Execute(tt.args)

			// Check exit code
			if exitCode != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d. Output: %s", tt.expectedCode, exitCode, buf.String())
			}

			// Parse response
			var resp map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse JSON response: %v. Raw: %s", err, buf.String())
			}

			// Check error code if expected
			if tt.expectedError != "" {
				if resp["success"] != false {
					t.Errorf("Expected success=false, got %v", resp["success"])
				}
				errorObj, ok := resp["error"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected error object in response")
				}
				if errorObj["code"] != tt.expectedError {
					t.Errorf("Expected error code %q, got %q", tt.expectedError, errorObj["code"])
				}
			}

			// Check success response
			if tt.checkSuccess {
				if resp["success"] != true {
					t.Errorf("Expected success=true, got %v. Response: %s", resp["success"], buf.String())
				}

				data, ok := resp["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected data object in response")
				}

				if tt.checkPruned != nil {
					prunedRaw, ok := data["pruned"].([]interface{})
					if !ok {
						t.Fatalf("Expected pruned array in data")
					}

					var pruned []string
					for _, p := range prunedRaw {
						pruned = append(pruned, p.(string))
					}

					if len(pruned) != len(tt.checkPruned) {
						t.Errorf("Expected %d pruned nooks, got %d: %v", len(tt.checkPruned), len(pruned), pruned)
					}
				}

				if tt.checkDecksRemoved != nil {
					decksRemovedRaw, ok := data["decks_removed"].([]interface{})
					if !ok {
						t.Fatalf("Expected decks_removed array in data")
					}

					var decksRemoved []string
					for _, d := range decksRemovedRaw {
						decksRemoved = append(decksRemoved, d.(string))
					}

					if len(decksRemoved) != len(tt.checkDecksRemoved) {
						t.Errorf("Expected %d decks removed, got %d: %v", len(tt.checkDecksRemoved), len(decksRemoved), decksRemoved)
					}
				}
			}

			// Verify crew memory files were deleted on prune
			if tt.name == "cleanup crew memory files on prune" {
				memFile := filepath.Join(tmpDir, ".treehouse", "crew", "spruce", "memories", "a1b2-orphan.md")
				sessFile := filepath.Join(tmpDir, ".treehouse", "crew", "spruce", "sessions", "a1b2-orphan.md")

				if _, err := os.Stat(memFile); !os.IsNotExist(err) {
					t.Errorf("Memory file should have been deleted: %s", memFile)
				}
				if _, err := os.Stat(sessFile); !os.IsNotExist(err) {
					t.Errorf("Session file should have been deleted: %s", sessFile)
				}
			}
		})
	}
}
