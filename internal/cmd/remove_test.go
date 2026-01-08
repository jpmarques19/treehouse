package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRemoveCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		setupFunc      func(t *testing.T, tmpDir string)
		expectedCode   int
		expectedError  string
		checkSuccess   bool
		checkRemoved   []string
		checkDeckGone  bool
	}{
		{
			name:          "missing nook ID",
			args:          []string{"remove"},
			expectedCode:  2,
			expectedError: "NOOK_ID_REQUIRED",
		},
		{
			name:          "nook not found",
			args:          []string{"remove", "nonexistent-nook"},
			expectedCode:  4,
			expectedError: "NOOK_NOT_FOUND",
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with empty decks
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				os.MkdirAll(treehousePath, 0755)
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte("decks: {}\n"), 0644)
			},
		},
		{
			name:         "remove single nook without children",
			args:         []string{"remove", "a1b2-test-nook"},
			expectedCode: 0,
			checkSuccess: true,
			checkRemoved: []string{"a1b2-test-nook"},
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with one nook
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				os.MkdirAll(nooksPath, 0755)

				// Create decks.yaml with one nook
				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-test-nook:
        parent: main
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)

				// Create a mock worktree directory
				os.MkdirAll(filepath.Join(nooksPath, "a1b2-test-nook"), 0755)
			},
		},
		{
			name:          "remove nook with children (cascade)",
			args:          []string{"remove", "a1b2-parent"},
			expectedCode:  0,
			checkSuccess:  true,
			checkRemoved:  []string{"c3d4-child", "a1b2-parent"},
			checkDeckGone: true,
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with parent and child nook
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				os.MkdirAll(nooksPath, 0755)

				// Create decks.yaml with parent and child
				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-parent:
        parent: main
        created: "2026-01-08"
      c3d4-child:
        parent: a1b2-parent
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)

				// Create mock worktree directories
				os.MkdirAll(filepath.Join(nooksPath, "a1b2-parent"), 0755)
				os.MkdirAll(filepath.Join(nooksPath, "c3d4-child"), 0755)
			},
		},
		{
			name:         "cleanup crew memory files",
			args:         []string{"remove", "a1b2-test-nook"},
			expectedCode: 0,
			checkSuccess: true,
			setupFunc: func(t *testing.T, tmpDir string) {
				// Create initialized treehouse with crew memory files
				treehousePath := filepath.Join(tmpDir, ".treehouse")
				nooksPath := filepath.Join(treehousePath, "nooks")
				crewPath := filepath.Join(treehousePath, "crew", "spruce")
				os.MkdirAll(nooksPath, 0755)
				os.MkdirAll(filepath.Join(crewPath, "memories"), 0755)
				os.MkdirAll(filepath.Join(crewPath, "sessions"), 0755)

				// Create decks.yaml
				decksYaml := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-test-nook:
        parent: main
        created: "2026-01-08"
`
				os.WriteFile(filepath.Join(treehousePath, "decks.yaml"), []byte(decksYaml), 0644)
				os.MkdirAll(filepath.Join(nooksPath, "a1b2-test-nook"), 0755)

				// Create memory files that should be deleted
				os.WriteFile(filepath.Join(crewPath, "memories", "a1b2-test-nook.md"), []byte("# Memory"), 0644)
				os.WriteFile(filepath.Join(crewPath, "sessions", "a1b2-test-nook.md"), []byte("# Session"), 0644)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Initialize as real git repo (required for worktree detection)
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

				if tt.checkRemoved != nil {
					data, ok := resp["data"].(map[string]interface{})
					if !ok {
						t.Fatalf("Expected data object in response")
					}
					removedRaw, ok := data["removed"].([]interface{})
					if !ok {
						t.Fatalf("Expected removed array in data")
					}

					var removed []string
					for _, r := range removedRaw {
						removed = append(removed, r.(string))
					}

					if len(removed) != len(tt.checkRemoved) {
						t.Errorf("Expected %d removed nooks, got %d: %v", len(tt.checkRemoved), len(removed), removed)
					}
					for i, expected := range tt.checkRemoved {
						if i < len(removed) && removed[i] != expected {
							t.Errorf("Expected removed[%d] = %q, got %q", i, expected, removed[i])
						}
					}

					if tt.checkDeckGone {
						if data["deck_removed"] != true {
							t.Errorf("Expected deck_removed=true, got %v", data["deck_removed"])
						}
					}
				}
			}

			// Verify crew memory files were deleted
			if tt.name == "cleanup crew memory files" {
				memFile := filepath.Join(tmpDir, ".treehouse", "crew", "spruce", "memories", "a1b2-test-nook.md")
				sessFile := filepath.Join(tmpDir, ".treehouse", "crew", "spruce", "sessions", "a1b2-test-nook.md")

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

// initGitRepo initializes a real git repo in the given directory
func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	cmd.Run()

	cmd = exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
}
