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

func TestInitCommand_Success(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Parse JSON response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v\nOutput: %s", err, buf.String())
	}

	if !resp.Success {
		t.Errorf("expected success=true, got success=false")
	}

	// Verify .treehouse folder structure was created
	treehousePath := filepath.Join(dir, ".treehouse")
	if _, err := os.Stat(treehousePath); os.IsNotExist(err) {
		t.Error(".treehouse folder was not created")
	}

	// Verify decks.yaml was created
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if _, err := os.Stat(decksPath); os.IsNotExist(err) {
		t.Error("decks.yaml was not created")
	}

	// Verify hats folder was created
	hatsPath := filepath.Join(treehousePath, "hats")
	if _, err := os.Stat(hatsPath); os.IsNotExist(err) {
		t.Error("hats/ folder was not created")
	}

	// Verify workflows folder was created
	workflowsPath := filepath.Join(treehousePath, "workflows")
	if _, err := os.Stat(workflowsPath); os.IsNotExist(err) {
		t.Error("workflows folder was not created")
	}
}

func TestInitCommand_AlreadyInitialized(t *testing.T) {
	dir := testutil.SetupGitRepo(t)

	// Create .treehouse folder to simulate already initialized
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("failed to create .treehouse: %v", err)
	}

	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	// Parse JSON response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v\nOutput: %s", err, buf.String())
	}

	if resp.Success {
		t.Error("expected success=false for already initialized")
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Code != "INIT_ALREADY_EXISTS" {
		t.Errorf("expected error code INIT_ALREADY_EXISTS, got %s", resp.Error.Code)
	}
}

func TestInitCommand_NotGitRepo(t *testing.T) {
	// Create a non-git directory
	dir := t.TempDir()
	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})

	if exitCode != 3 {
		t.Errorf("expected exit code 3, got %d", exitCode)
	}

	// Parse JSON response
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v\nOutput: %s", err, buf.String())
	}

	if resp.Success {
		t.Error("expected success=false for non-git repo")
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Code != "INIT_NOT_GIT_REPO" {
		t.Errorf("expected error code INIT_NOT_GIT_REPO, got %s", resp.Error.Code)
	}
}

func TestDecksYamlContent(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("init failed with exit code %d", exitCode)
	}

	// Read and verify decks.yaml content
	decksPath := filepath.Join(dir, ".treehouse", "decks.yaml")
	content, err := os.ReadFile(decksPath)
	if err != nil {
		t.Fatalf("failed to read decks.yaml: %v", err)
	}

	expected := "decks: {}\n"
	if string(content) != expected {
		t.Errorf("decks.yaml content = %q, want %q", string(content), expected)
	}
}


func TestWorkflowInstallation(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("init failed with exit code %d", exitCode)
	}

	// Verify workflows are installed
	workflowsPath := filepath.Join(dir, ".treehouse", "workflows")
	expectedWorkflows := []string{"nook-fork.md", "treehouse-init.md", "treehouse-view.md"}

	for _, wf := range expectedWorkflows {
		wfPath := filepath.Join(workflowsPath, wf)
		if _, err := os.Stat(wfPath); os.IsNotExist(err) {
			t.Errorf("workflow %s was not installed", wf)
		}
	}
}

func TestClaudeCommandInstallation(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("init failed with exit code %d", exitCode)
	}

	// Verify Claude commands are installed
	claudePath := filepath.Join(dir, ".claude", "commands", "th", "workflows")
	expectedCommands := []string{"nook-fork.md", "treehouse-init.md", "treehouse-view.md"}

	for _, cmd := range expectedCommands {
		cmdPath := filepath.Join(claudePath, cmd)
		if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
			t.Errorf("Claude command %s was not installed", cmd)
		}
	}

	// Verify stub content has loader directive
	nookForkPath := filepath.Join(claudePath, "nook-fork.md")
	content, err := os.ReadFile(nookForkPath)
	if err != nil {
		t.Fatalf("failed to read nook-fork.md: %v", err)
	}

	if !bytes.Contains(content, []byte("<workflow-loader>")) {
		t.Error("nook-fork.md missing workflow-loader directive")
	}

	if !bytes.Contains(content, []byte(".treehouse/workflows/nook-fork.md")) {
		t.Error("nook-fork.md has incorrect workflow path")
	}
}


func TestGitignoreCreation(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("init failed with exit code %d", exitCode)
	}

	// Verify .gitignore was created with treehouse entries
	gitignorePath := filepath.Join(dir, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	if !bytes.Contains(content, []byte(".treehouse/")) {
		t.Error(".gitignore missing .treehouse/ entry")
	}

	if !bytes.Contains(content, []byte(".claude/commands/th")) {
		t.Error(".gitignore missing .claude/commands/th entry")
	}

	if !bytes.Contains(content, []byte("# treehouse managed")) {
		t.Error(".gitignore missing treehouse managed marker")
	}
}

func TestGitignoreIdempotent(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Create existing .gitignore with treehouse entries
	existingContent := `# my project
node_modules/

# treehouse managed
.treehouse/
.claude/commands/th
`
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("init failed with exit code %d", exitCode)
	}

	// Verify .gitignore was not duplicated
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	// Count occurrences of treehouse marker
	count := bytes.Count(content, []byte("# treehouse managed"))
	if count != 1 {
		t.Errorf("expected 1 treehouse managed marker, found %d", count)
	}

	// Parse response to verify gitignore_updated is false
	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		if updated, ok := data["gitignore_updated"].(bool); ok && updated {
			t.Error("expected gitignore_updated to be false when entries already exist")
		}
	}
}

func TestGitignoreAppendToExisting(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Create existing .gitignore without treehouse entries
	existingContent := `# my project
node_modules/
*.log`
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Run init command
	exitCode := Execute([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("init failed with exit code %d", exitCode)
	}

	// Verify .gitignore was appended
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	// Check existing content is preserved
	if !bytes.Contains(content, []byte("node_modules/")) {
		t.Error("existing .gitignore content was lost")
	}

	// Check treehouse entries were added
	if !bytes.Contains(content, []byte(".treehouse/")) {
		t.Error(".gitignore missing .treehouse/ entry after append")
	}

	if !bytes.Contains(content, []byte(".claude/commands/th")) {
		t.Error(".gitignore missing .claude/commands/th entry after append")
	}
}
