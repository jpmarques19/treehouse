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

func setupInitializedRepo(t *testing.T) string {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Initialize treehouse
	var buf bytes.Buffer
	output.SetWriter(&buf)
	exitCode := Execute([]string{"init"})
	output.SetWriter(os.Stdout)

	if exitCode != 0 {
		t.Fatalf("failed to initialize treehouse: %s", buf.String())
	}

	return dir
}

func TestHatAdd_Success(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Create hat
	config := `{"name":"reviewer","title":"Code Reviewer","icon":"🔍"}`

	exitCode := Execute([]string{"hat", "add", "reviewer", config})

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

	// Verify hat .md file was created
	hatPath := filepath.Join(dir, ".treehouse", "hats", "reviewer.md")
	if _, err := os.Stat(hatPath); os.IsNotExist(err) {
		t.Error("hat .md file was not created")
	}

	// Verify Claude command stub was created
	commandPath := filepath.Join(dir, ".claude", "commands", "th", "hats", "reviewer.md")
	if _, err := os.Stat(commandPath); os.IsNotExist(err) {
		t.Error("Claude command stub was not created")
	}
}

func TestHatAdd_AlreadyExists(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Create hat first
	var buf bytes.Buffer
	output.SetWriter(&buf)

	config := `{"name":"reviewer","title":"Code Reviewer","icon":"🔍"}`
	Execute([]string{"hat", "add", "reviewer", config})

	// Try to create again
	buf.Reset()
	exitCode := Execute([]string{"hat", "add", "reviewer", config})
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
		t.Error("expected success=false for already existing hat")
	}

	if resp.Error == nil || resp.Error.Code != "HAT_ALREADY_EXISTS" {
		t.Errorf("expected error code HAT_ALREADY_EXISTS, got %v", resp.Error)
	}

	// Verify hat .md file still exists
	hatPath := filepath.Join(dir, ".treehouse", "hats", "reviewer.md")
	if _, err := os.Stat(hatPath); os.IsNotExist(err) {
		t.Error("original hat .md file should still exist")
	}
}

func TestHatAdd_InvalidJSON(t *testing.T) {
	setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"hat", "add", "reviewer", "not valid json"})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}

	// Parse JSON response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false for invalid JSON")
	}

	if resp.Error == nil || resp.Error.Code != "HAT_INVALID_CONFIG" {
		t.Errorf("expected error code HAT_INVALID_CONFIG, got %v", resp.Error)
	}
}

func TestHatAdd_MissingRequiredFields(t *testing.T) {
	setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Missing title
	config := `{"name":"reviewer"}`
	exitCode := Execute([]string{"hat", "add", "reviewer", config})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != "HAT_INVALID_CONFIG" {
		t.Errorf("expected error code HAT_INVALID_CONFIG, got %v", resp.Error)
	}
}

func TestHatAdd_NotInitialized(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Don't initialize treehouse

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	config := `{"name":"reviewer","title":"Code Reviewer","icon":"🔍"}`
	exitCode := Execute([]string{"hat", "add", "reviewer", config})

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

func TestHatAdd_DefaultIcon(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// No icon in config
	config := `{"name":"reviewer","title":"Code Reviewer"}`
	exitCode := Execute([]string{"hat", "add", "reviewer", config})

	if exitCode != 0 {
		t.Fatalf("hat add failed: %s", buf.String())
	}

	// Read Claude command stub and check for default icon
	commandPath := filepath.Join(dir, ".claude", "commands", "th", "hats", "reviewer.md")
	content, err := os.ReadFile(commandPath)
	if err != nil {
		t.Fatalf("failed to read Claude command stub: %v", err)
	}

	if !bytes.Contains(content, []byte("◇")) {
		t.Error("Claude command stub should have default icon ◇")
	}
}

func TestHatAdd_NameSanitization(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Name with spaces and special chars
	config := `{"name":"Code Reviewer!","title":"Code Reviewer"}`
	exitCode := Execute([]string{"hat", "add", "Code Reviewer!", config})

	if exitCode != 0 {
		t.Fatalf("hat add failed: %s", buf.String())
	}

	// Should be sanitized to "code-reviewer.md"
	hatPath := filepath.Join(dir, ".treehouse", "hats", "code-reviewer.md")
	if _, err := os.Stat(hatPath); os.IsNotExist(err) {
		t.Error("hat .md file was not created with sanitized name")
	}
}

func TestHatAdd_ClaudeCommandContent(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	config := `{"name":"reviewer","title":"Code Reviewer","icon":"🔍"}`
	exitCode := Execute([]string{"hat", "add", "reviewer", config})

	if exitCode != 0 {
		t.Fatalf("hat add failed: %s", buf.String())
	}

	// Read Claude command stub
	commandPath := filepath.Join(dir, ".claude", "commands", "th", "hats", "reviewer.md")
	content, err := os.ReadFile(commandPath)
	if err != nil {
		t.Fatalf("failed to read Claude command stub: %v", err)
	}

	// Check content - should have hat-loader block referencing .md file
	requiredContent := []string{
		"🔍",
		"Reviewer",
		"Code Reviewer",
		"<hat-loader>",
		"reviewer.md",
	}

	for _, s := range requiredContent {
		if !bytes.Contains(content, []byte(s)) {
			t.Errorf("Claude command stub missing: %s", s)
		}
	}

	// Verify it does NOT contain agent-loader (old crew approach)
	if bytes.Contains(content, []byte("agent-loader")) {
		t.Error("Claude command stub should not contain 'agent-loader' reference")
	}
}

func TestHatCommand_NoSubcommand(t *testing.T) {
	setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"hat"})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != "HAT_NO_SUBCOMMAND" {
		t.Errorf("expected error code HAT_NO_SUBCOMMAND, got %v", resp.Error)
	}
}

func TestSanitizeHatName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"reviewer", "reviewer"},
		{"Code Reviewer", "code-reviewer"},
		{"test@hat!", "testhat"},
		{"My--Hat", "my-hat"},
		{"-leading-trailing-", "leading-trailing"},
		{"UPPERCASE", "uppercase"},
		{"with123numbers", "with123numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeHatName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeHatName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
