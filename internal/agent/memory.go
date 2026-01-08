package agent

import (
	"os"
	"path/filepath"

	"github.com/jpmarques19/treehouse/internal/nook"
)

// DeleteMemoryFiles removes memory and session files for a nook across all agents
// Directory structure: .treehouse/agents/*/memories/{nook-id}.md and .treehouse/agents/*/sessions/{nook-id}.md
func DeleteMemoryFiles(treehousePath, nookID string) error {
	// Validate nookID to prevent path traversal
	if err := nook.ValidateNookID(nookID); err != nil {
		return err
	}

	agentsPath := filepath.Join(treehousePath, "agents")

	// Check if agents directory exists
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		return nil // Nothing to clean up
	}

	// Pattern: .treehouse/agents/*/memories/{nook-id}.md
	memPattern := filepath.Join(agentsPath, "*", "memories", nookID+".md")
	memFiles, _ := filepath.Glob(memPattern)
	for _, f := range memFiles {
		_ = os.Remove(f) // Best-effort removal
	}

	// Pattern: .treehouse/agents/*/sessions/{nook-id}.md
	sessPattern := filepath.Join(agentsPath, "*", "sessions", nookID+".md")
	sessFiles, _ := filepath.Glob(sessPattern)
	for _, f := range sessFiles {
		_ = os.Remove(f) // Best-effort removal
	}

	return nil
}
