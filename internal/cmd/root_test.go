package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/jpmarques19/treehouse/internal/output"
)

func TestRootCommand_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	exitCode := Execute([]string{})

	if exitCode != 2 {
		t.Errorf("Execute() exit code = %d, want 2", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Output is not valid JSON: %v, got: %s", err, buf.String())
	}

	if resp.Success {
		t.Error("Expected success=false for no args")
	}

	if resp.Error == nil {
		t.Fatal("Expected error in response")
	}

	if resp.Error.Code != "NO_COMMAND" {
		t.Errorf("Error code = %s, want NO_COMMAND", resp.Error.Code)
	}

	if resp.Error.Message != "No command specified" {
		t.Errorf("Error message = %s, want 'No command specified'", resp.Error.Message)
	}
}

func TestRootCommand_Version(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	exitCode := Execute([]string{"--version"})

	if exitCode != 0 {
		t.Errorf("Execute() exit code = %d, want 0", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Output is not valid JSON: %v, got: %s", err, buf.String())
	}

	if !resp.Success {
		t.Error("Expected success=true for --version")
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected Data to be a map, got %T", resp.Data)
	}

	version, ok := data["version"].(string)
	if !ok {
		t.Fatal("Expected version to be a string")
	}

	if version != "0.4.0" {
		t.Errorf("Version = %s, want 0.4.0", version)
	}
}

func TestRootCommand_ShortVersion(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	exitCode := Execute([]string{"-v"})

	if exitCode != 0 {
		t.Errorf("Execute() exit code = %d, want 0", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Output is not valid JSON: %v, got: %s", err, buf.String())
	}

	if !resp.Success {
		t.Error("Expected success=true for -v")
	}
}

func TestRootCommand_UnknownFlag(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	exitCode := Execute([]string{"--unknown"})

	if exitCode != 2 {
		t.Errorf("Execute() exit code = %d, want 2", exitCode)
	}

	var resp output.Response
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("Output is not valid JSON: %v, got: %s", err, buf.String())
	}

	if resp.Success {
		t.Error("Expected success=false for unknown flag")
	}
}
