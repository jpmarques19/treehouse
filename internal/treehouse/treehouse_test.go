package treehouse

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindTreehouse_InCurrentDirectory(t *testing.T) {
	dir := t.TempDir()

	// Create .treehouse directory
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	info, err := FindTreehouse(dir)
	if err != nil {
		t.Fatalf("FindTreehouse() error = %v", err)
	}

	if info.TreehousePath != treehousePath {
		t.Errorf("TreehousePath = %q, want %q", info.TreehousePath, treehousePath)
	}

	if info.RepoRoot != dir {
		t.Errorf("RepoRoot = %q, want %q", info.RepoRoot, dir)
	}
}

func TestFindTreehouse_InParentDirectory(t *testing.T) {
	dir := t.TempDir()

	// Create .treehouse in root
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create a subdirectory
	subdir := filepath.Join(dir, "subdir", "nested")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Search from subdirectory
	info, err := FindTreehouse(subdir)
	if err != nil {
		t.Fatalf("FindTreehouse() error = %v", err)
	}

	if info.TreehousePath != treehousePath {
		t.Errorf("TreehousePath = %q, want %q", info.TreehousePath, treehousePath)
	}

	if info.RepoRoot != dir {
		t.Errorf("RepoRoot = %q, want %q", info.RepoRoot, dir)
	}
}

func TestFindTreehouse_NotFound(t *testing.T) {
	dir := t.TempDir()

	// No .treehouse directory
	_, err := FindTreehouse(dir)
	if err == nil {
		t.Fatal("Expected error when .treehouse not found")
	}

	thErr, ok := err.(*TreehouseError)
	if !ok {
		t.Fatalf("Expected TreehouseError, got %T", err)
	}

	if thErr.Code != "INIT_NOT_FOUND" {
		t.Errorf("Error code = %q, want %q", thErr.Code, "INIT_NOT_FOUND")
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()

	// Not initialized
	if Exists(dir) {
		t.Error("Expected Exists to return false when .treehouse doesn't exist")
	}

	// Create .treehouse
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Now initialized
	if !Exists(dir) {
		t.Error("Expected Exists to return true when .treehouse exists")
	}
}
