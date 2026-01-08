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

	// Verify crew/oak folder was created
	oakPath := filepath.Join(treehousePath, "crew", "oak")
	if _, err := os.Stat(oakPath); os.IsNotExist(err) {
		t.Error("crew/oak folder was not created")
	}

	// Verify oak.agent.yaml was created
	oakAgentPath := filepath.Join(oakPath, "oak.agent.yaml")
	if _, err := os.Stat(oakAgentPath); os.IsNotExist(err) {
		t.Error("oak.agent.yaml was not created")
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

func TestOakAgentCreation(t *testing.T) {
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

	oakPath := filepath.Join(dir, ".treehouse", "crew", "oak")

	// Verify oak.agent.yaml exists and contains required fields
	agentPath := filepath.Join(oakPath, "oak.agent.yaml")
	content, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("failed to read oak.agent.yaml: %v", err)
	}

	// Check for required content
	requiredFields := []string{"name:", "Oak", "title:", "Treehouse Assistant", "icon:"}
	for _, field := range requiredFields {
		if !bytes.Contains(content, []byte(field)) {
			t.Errorf("oak.agent.yaml missing required field: %s", field)
		}
	}

	// Verify knowledge.md exists
	knowledgePath := filepath.Join(oakPath, "knowledge.md")
	if _, err := os.Stat(knowledgePath); os.IsNotExist(err) {
		t.Error("knowledge.md was not created")
	}

	// Verify memories/ directory exists
	memoriesPath := filepath.Join(oakPath, "memories")
	if _, err := os.Stat(memoriesPath); os.IsNotExist(err) {
		t.Error("memories/ directory was not created")
	}

	// Verify sessions/ directory exists
	sessionsPath := filepath.Join(oakPath, "sessions")
	if _, err := os.Stat(sessionsPath); os.IsNotExist(err) {
		t.Error("sessions/ directory was not created")
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
	expectedWorkflows := []string{"huddle.md", "nook-fork.md", "treehouse-init.md", "treehouse-list.md"}

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
	expectedCommands := []string{"huddle.md", "nook-fork.md", "treehouse-init.md", "treehouse-list.md"}

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
