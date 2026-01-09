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

func TestCrewAdd_Success(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Create crew member
	config := `{"name":"sparrow","title":"Code Reviewer","icon":"🔍","persona":{"role":"Review code for quality and correctness","identity":"Meticulous code reviewer with an eye for detail","communication_style":"Direct and constructive feedback","principles":["Code quality over speed","Clear explanations"]}}`

	exitCode := Execute([]string{"crew", "add", "sparrow", config})

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

	// Verify agent folder was created
	agentPath := filepath.Join(dir, ".treehouse", "agents", "sparrow")
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		t.Error("agent folder was not created")
	}

	// Verify agent YAML was created
	agentYAML := filepath.Join(agentPath, "sparrow.agent.yaml")
	if _, err := os.Stat(agentYAML); os.IsNotExist(err) {
		t.Error("agent YAML was not created")
	}

	// Verify knowledge.md was created
	knowledgePath := filepath.Join(agentPath, "knowledge.md")
	if _, err := os.Stat(knowledgePath); os.IsNotExist(err) {
		t.Error("knowledge.md was not created")
	}

	// Verify memories/ folder was created
	memoriesPath := filepath.Join(agentPath, "memories")
	if _, err := os.Stat(memoriesPath); os.IsNotExist(err) {
		t.Error("memories/ folder was not created")
	}

	// Verify sessions/ folder was created
	sessionsPath := filepath.Join(agentPath, "sessions")
	if _, err := os.Stat(sessionsPath); os.IsNotExist(err) {
		t.Error("sessions/ folder was not created")
	}

	// Verify Claude command stub was created
	commandPath := filepath.Join(dir, ".claude", "commands", "th", "crew", "sparrow.md")
	if _, err := os.Stat(commandPath); os.IsNotExist(err) {
		t.Error("Claude command stub was not created")
	}
}

func TestCrewAdd_AgentYAMLContent(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	config := `{"name":"reviewer","title":"PR Reviewer","icon":"👁","persona":{"role":"Review pull requests","identity":"Careful reviewer","communication_style":"Constructive","principles":["Quality first"]}}`

	exitCode := Execute([]string{"crew", "add", "reviewer", config})
	if exitCode != 0 {
		t.Fatalf("crew add failed: %s", buf.String())
	}

	// Read agent YAML
	agentYAML := filepath.Join(dir, ".treehouse", "agents", "reviewer", "reviewer.agent.yaml")
	content, err := os.ReadFile(agentYAML)
	if err != nil {
		t.Fatalf("failed to read agent YAML: %v", err)
	}

	// Check required content
	requiredFields := []string{
		"name:",
		"Reviewer",
		"title:",
		"PR Reviewer",
		"icon:",
		"👁",
		"role:",
		"Review pull requests",
		"identity:",
		"Careful reviewer",
		"communication_style:",
		"Constructive",
		"principles:",
		"Quality first",
		"critical_actions:",
	}

	for _, field := range requiredFields {
		if !bytes.Contains(content, []byte(field)) {
			t.Errorf("agent YAML missing required field: %s", field)
		}
	}
}

func TestCrewAdd_AlreadyExists(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Create crew member first
	var buf bytes.Buffer
	output.SetWriter(&buf)

	config := `{"name":"sparrow","title":"Code Reviewer","icon":"🔍","persona":{"role":"Review code","identity":"Reviewer","communication_style":"Direct","principles":["Quality"]}}`
	Execute([]string{"crew", "add", "sparrow", config})

	// Try to create again
	buf.Reset()
	exitCode := Execute([]string{"crew", "add", "sparrow", config})
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
		t.Error("expected success=false for already existing crew")
	}

	if resp.Error == nil || resp.Error.Code != "CREW_ALREADY_EXISTS" {
		t.Errorf("expected error code CREW_ALREADY_EXISTS, got %v", resp.Error)
	}

	// Verify only one agent folder exists
	agentPath := filepath.Join(dir, ".treehouse", "agents", "sparrow")
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		t.Error("original agent folder should still exist")
	}
}

func TestCrewAdd_InvalidJSON(t *testing.T) {
	setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"crew", "add", "sparrow", "not valid json"})

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

	if resp.Error == nil || resp.Error.Code != "CREW_INVALID_CONFIG" {
		t.Errorf("expected error code CREW_INVALID_CONFIG, got %v", resp.Error)
	}
}

func TestCrewAdd_MissingRequiredFields(t *testing.T) {
	setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Missing title
	config := `{"name":"sparrow","persona":{"role":"Review code","identity":"Reviewer","communication_style":"Direct","principles":["Quality"]}}`
	exitCode := Execute([]string{"crew", "add", "sparrow", config})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != "CREW_INVALID_CONFIG" {
		t.Errorf("expected error code CREW_INVALID_CONFIG, got %v", resp.Error)
	}
}

func TestCrewAdd_NotInitialized(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	// Don't initialize treehouse

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	config := `{"name":"sparrow","title":"Code Reviewer","icon":"🔍","persona":{"role":"Review code","identity":"Reviewer","communication_style":"Direct","principles":["Quality"]}}`
	exitCode := Execute([]string{"crew", "add", "sparrow", config})

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

func TestCrewAdd_DefaultIcon(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// No icon in config
	config := `{"name":"sparrow","title":"Code Reviewer","persona":{"role":"Review code","identity":"Reviewer","communication_style":"Direct","principles":["Quality"]}}`
	exitCode := Execute([]string{"crew", "add", "sparrow", config})

	if exitCode != 0 {
		t.Fatalf("crew add failed: %s", buf.String())
	}

	// Read agent YAML and check for default icon
	agentYAML := filepath.Join(dir, ".treehouse", "agents", "sparrow", "sparrow.agent.yaml")
	content, err := os.ReadFile(agentYAML)
	if err != nil {
		t.Fatalf("failed to read agent YAML: %v", err)
	}

	if !bytes.Contains(content, []byte("◇")) {
		t.Error("agent YAML should have default icon ◇")
	}
}

func TestCrewAdd_NameSanitization(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	// Name with spaces and special chars
	config := `{"name":"Code Reviewer!","title":"Code Reviewer","persona":{"role":"Review code","identity":"Reviewer","communication_style":"Direct","principles":["Quality"]}}`
	exitCode := Execute([]string{"crew", "add", "Code Reviewer!", config})

	if exitCode != 0 {
		t.Fatalf("crew add failed: %s", buf.String())
	}

	// Should be sanitized to "code-reviewer"
	agentPath := filepath.Join(dir, ".treehouse", "agents", "code-reviewer")
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		t.Error("agent folder was not created with sanitized name")
	}
}

func TestCrewAdd_ClaudeCommandContent(t *testing.T) {
	dir := setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	config := `{"name":"sparrow","title":"Code Reviewer","icon":"🔍","persona":{"role":"Review code","identity":"Reviewer","communication_style":"Direct","principles":["Quality"]}}`
	exitCode := Execute([]string{"crew", "add", "sparrow", config})

	if exitCode != 0 {
		t.Fatalf("crew add failed: %s", buf.String())
	}

	// Read Claude command stub
	commandPath := filepath.Join(dir, ".claude", "commands", "th", "crew", "sparrow.md")
	content, err := os.ReadFile(commandPath)
	if err != nil {
		t.Fatalf("failed to read Claude command stub: %v", err)
	}

	// Check content
	requiredContent := []string{
		"🔍",
		"Sparrow",
		"Code Reviewer",
		"agent-loader",
		"sparrow",
	}

	for _, s := range requiredContent {
		if !bytes.Contains(content, []byte(s)) {
			t.Errorf("Claude command stub missing: %s", s)
		}
	}
}

func TestCrewCommand_NoSubcommand(t *testing.T) {
	setupInitializedRepo(t)

	// Capture output
	var buf bytes.Buffer
	output.SetWriter(&buf)
	t.Cleanup(func() { output.SetWriter(os.Stdout) })

	exitCode := Execute([]string{"crew"})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != "CREW_NO_SUBCOMMAND" {
		t.Errorf("expected error code CREW_NO_SUBCOMMAND, got %v", resp.Error)
	}
}

func TestSanitizeCrewName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sparrow", "sparrow"},
		{"Code Reviewer", "code-reviewer"},
		{"test@agent!", "testagent"},
		{"My--Agent", "my-agent"},
		{"-leading-trailing-", "leading-trailing"},
		{"UPPERCASE", "uppercase"},
		{"with123numbers", "with123numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeCrewName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeCrewName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
