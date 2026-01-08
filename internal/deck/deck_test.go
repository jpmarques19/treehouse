package deck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDecks_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create empty decks.yaml
	decksContent := "decks: {}\n"
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte(decksContent), 0644); err != nil {
		t.Fatalf("Failed to write decks.yaml: %v", err)
	}

	decks, err := LoadDecks(treehousePath)
	if err != nil {
		t.Fatalf("LoadDecks() error = %v", err)
	}

	if decks.Decks == nil {
		t.Error("Decks should not be nil")
	}

	if len(decks.Decks) != 0 {
		t.Errorf("Expected 0 decks, got %d", len(decks.Decks))
	}
}

func TestLoadDecks_WithNooks(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create decks.yaml with data
	decksContent := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-auth-spike:
        parent: main
        created: "2026-01-08"
      c3d4-jwt-variant:
        parent: a1b2-auth-spike
        created: "2026-01-09"
`
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte(decksContent), 0644); err != nil {
		t.Fatalf("Failed to write decks.yaml: %v", err)
	}

	decks, err := LoadDecks(treehousePath)
	if err != nil {
		t.Fatalf("LoadDecks() error = %v", err)
	}

	if len(decks.Decks) != 1 {
		t.Errorf("Expected 1 deck, got %d", len(decks.Decks))
	}

	deck, ok := decks.Decks["dk-a1b2"]
	if !ok {
		t.Fatal("Expected deck dk-a1b2 to exist")
	}

	if len(deck.Nooks) != 2 {
		t.Errorf("Expected 2 nooks, got %d", len(deck.Nooks))
	}

	nook, ok := deck.Nooks["a1b2-auth-spike"]
	if !ok {
		t.Fatal("Expected nook a1b2-auth-spike to exist")
	}

	if nook.Parent != "main" {
		t.Errorf("Nook parent = %q, want %q", nook.Parent, "main")
	}
}

func TestLoadDecks_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	_, err := LoadDecks(treehousePath)
	if err == nil {
		t.Fatal("Expected error for missing decks.yaml")
	}

	deckErr, ok := err.(*DeckError)
	if !ok {
		t.Fatalf("Expected DeckError, got %T", err)
	}

	if deckErr.Code != "DECK_FILE_NOT_FOUND" {
		t.Errorf("Error code = %q, want %q", deckErr.Code, "DECK_FILE_NOT_FOUND")
	}
}

func TestNookExists(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create decks.yaml with data
	decksContent := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-auth-spike:
        parent: main
        created: "2026-01-08"
`
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte(decksContent), 0644); err != nil {
		t.Fatalf("Failed to write decks.yaml: %v", err)
	}

	tests := []struct {
		name     string
		nookID   string
		expected bool
	}{
		{
			name:     "existing nook",
			nookID:   "a1b2-auth-spike",
			expected: true,
		},
		{
			name:     "non-existing nook",
			nookID:   "x1y2-nonexistent",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := NookExists(treehousePath, tt.nookID)
			if err != nil {
				t.Fatalf("NookExists() error = %v", err)
			}

			if exists != tt.expected {
				t.Errorf("NookExists(%q) = %v, want %v", tt.nookID, exists, tt.expected)
			}
		})
	}
}

func TestNookExists_NoDecksFile(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// No decks.yaml file
	exists, err := NookExists(treehousePath, "a1b2-test")
	if err != nil {
		t.Fatalf("NookExists() error = %v", err)
	}

	if exists {
		t.Error("Expected false when decks.yaml doesn't exist")
	}
}

func TestGetNook(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create decks.yaml with data
	decksContent := `decks:
  dk-a1b2:
    created: "2026-01-08"
    nooks:
      a1b2-auth-spike:
        parent: main
        created: "2026-01-08"
`
	decksPath := filepath.Join(treehousePath, "decks.yaml")
	if err := os.WriteFile(decksPath, []byte(decksContent), 0644); err != nil {
		t.Fatalf("Failed to write decks.yaml: %v", err)
	}

	// Test finding existing nook
	nook, deckID, err := GetNook(treehousePath, "a1b2-auth-spike")
	if err != nil {
		t.Fatalf("GetNook() error = %v", err)
	}

	if deckID != "dk-a1b2" {
		t.Errorf("Deck ID = %q, want %q", deckID, "dk-a1b2")
	}

	if nook.Parent != "main" {
		t.Errorf("Nook parent = %q, want %q", nook.Parent, "main")
	}

	// Test non-existing nook
	_, _, err = GetNook(treehousePath, "x1y2-nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existing nook")
	}

	deckErr, ok := err.(*DeckError)
	if !ok {
		t.Fatalf("Expected DeckError, got %T", err)
	}

	if deckErr.Code != "NOOK_NOT_FOUND" {
		t.Errorf("Error code = %q, want %q", deckErr.Code, "NOOK_NOT_FOUND")
	}
}

func TestSaveDecks(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	// Create a DecksFile to save
	decks := &DecksFile{
		Decks: map[string]*Deck{
			"dk-a1b2": {
				Created: "2026-01-08",
				Nooks: map[string]*Nook{
					"a1b2-auth-spike": {
						Parent:  "main",
						Created: "2026-01-08",
					},
				},
			},
		},
	}

	// Save the decks
	if err := SaveDecks(treehousePath, decks); err != nil {
		t.Fatalf("SaveDecks() error = %v", err)
	}

	// Load it back and verify
	loaded, err := LoadDecks(treehousePath)
	if err != nil {
		t.Fatalf("LoadDecks() error = %v", err)
	}

	if len(loaded.Decks) != 1 {
		t.Errorf("Expected 1 deck, got %d", len(loaded.Decks))
	}

	deck, ok := loaded.Decks["dk-a1b2"]
	if !ok {
		t.Fatal("Expected deck dk-a1b2 to exist")
	}

	if deck.Created != "2026-01-08" {
		t.Errorf("Deck created = %q, want %q", deck.Created, "2026-01-08")
	}

	nook, ok := deck.Nooks["a1b2-auth-spike"]
	if !ok {
		t.Fatal("Expected nook a1b2-auth-spike to exist")
	}

	if nook.Parent != "main" {
		t.Errorf("Nook parent = %q, want %q", nook.Parent, "main")
	}
}

func TestSaveDecks_NoTempFileLeftOnSuccess(t *testing.T) {
	dir := t.TempDir()
	treehousePath := filepath.Join(dir, ".treehouse")
	if err := os.MkdirAll(treehousePath, 0755); err != nil {
		t.Fatalf("Failed to create .treehouse: %v", err)
	}

	decks := &DecksFile{Decks: map[string]*Deck{}}

	if err := SaveDecks(treehousePath, decks); err != nil {
		t.Fatalf("SaveDecks() error = %v", err)
	}

	// Check that no temp files remain
	entries, err := os.ReadDir(treehousePath)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	for _, entry := range entries {
		if entry.Name() != "decks.yaml" {
			t.Errorf("Unexpected file in .treehouse: %s", entry.Name())
		}
	}
}

func TestGenerateDeckID(t *testing.T) {
	tests := []struct {
		hash     string
		expected string
	}{
		{"a1b2", "dk-a1b2"},
		{"c3d4", "dk-c3d4"},
		{"ABCD", "dk-ABCD"},
	}

	for _, tt := range tests {
		result := GenerateDeckID(tt.hash)
		if result != tt.expected {
			t.Errorf("GenerateDeckID(%q) = %q, want %q", tt.hash, result, tt.expected)
		}
	}
}

func TestCreateDeck(t *testing.T) {
	deck := CreateDeck("dk-a1b2", "2026-01-08")

	if deck.Created != "2026-01-08" {
		t.Errorf("Deck created = %q, want %q", deck.Created, "2026-01-08")
	}

	if deck.Nooks == nil {
		t.Error("Deck nooks should not be nil")
	}

	if len(deck.Nooks) != 0 {
		t.Errorf("Expected 0 nooks, got %d", len(deck.Nooks))
	}
}

func TestAddNookToDeck_NewDeck(t *testing.T) {
	decks := &DecksFile{Decks: map[string]*Deck{}}

	deckID := AddNookToDeck(decks, "dk-a1b2", "a1b2-auth-spike", "main", "2026-01-08")

	if deckID != "dk-a1b2" {
		t.Errorf("Deck ID = %q, want %q", deckID, "dk-a1b2")
	}

	// Verify deck was created
	deck, ok := decks.Decks["dk-a1b2"]
	if !ok {
		t.Fatal("Expected deck dk-a1b2 to exist")
	}

	if deck.Created != "2026-01-08" {
		t.Errorf("Deck created = %q, want %q", deck.Created, "2026-01-08")
	}

	// Verify nook was added
	nook, ok := deck.Nooks["a1b2-auth-spike"]
	if !ok {
		t.Fatal("Expected nook a1b2-auth-spike to exist")
	}

	if nook.Parent != "main" {
		t.Errorf("Nook parent = %q, want %q", nook.Parent, "main")
	}
}

func TestAddNookToDeck_ExistingDeck(t *testing.T) {
	decks := &DecksFile{
		Decks: map[string]*Deck{
			"dk-a1b2": {
				Created: "2026-01-08",
				Nooks: map[string]*Nook{
					"a1b2-auth-spike": {
						Parent:  "main",
						Created: "2026-01-08",
					},
				},
			},
		},
	}

	// Add a child nook to existing deck
	AddNookToDeck(decks, "dk-a1b2", "c3d4-jwt-variant", "a1b2-auth-spike", "2026-01-09")

	// Verify the original nook is still there
	_, ok := decks.Decks["dk-a1b2"].Nooks["a1b2-auth-spike"]
	if !ok {
		t.Fatal("Original nook should still exist")
	}

	// Verify new nook was added
	nook, ok := decks.Decks["dk-a1b2"].Nooks["c3d4-jwt-variant"]
	if !ok {
		t.Fatal("Expected nook c3d4-jwt-variant to exist")
	}

	if nook.Parent != "a1b2-auth-spike" {
		t.Errorf("Nook parent = %q, want %q", nook.Parent, "a1b2-auth-spike")
	}
}

func TestRemoveNook(t *testing.T) {
	decks := &DecksFile{
		Decks: map[string]*Deck{
			"dk-a1b2": {
				Created: "2026-01-08",
				Nooks: map[string]*Nook{
					"a1b2-auth-spike": {
						Parent:  "main",
						Created: "2026-01-08",
					},
					"c3d4-jwt-variant": {
						Parent:  "a1b2-auth-spike",
						Created: "2026-01-09",
					},
				},
			},
		},
	}

	// Remove one nook (should not remove deck)
	deckEmpty, err := RemoveNook(decks, "dk-a1b2", "c3d4-jwt-variant")
	if err != nil {
		t.Fatalf("RemoveNook() error = %v", err)
	}

	if deckEmpty {
		t.Error("Deck should not be empty after removing one of two nooks")
	}

	if _, exists := decks.Decks["dk-a1b2"].Nooks["c3d4-jwt-variant"]; exists {
		t.Error("Removed nook should not exist")
	}

	// Remove last nook (should remove deck)
	deckEmpty, err = RemoveNook(decks, "dk-a1b2", "a1b2-auth-spike")
	if err != nil {
		t.Fatalf("RemoveNook() error = %v", err)
	}

	if !deckEmpty {
		t.Error("Deck should be empty after removing last nook")
	}

	if _, exists := decks.Decks["dk-a1b2"]; exists {
		t.Error("Empty deck should be removed")
	}
}

func TestGetDeckForNook(t *testing.T) {
	decks := &DecksFile{
		Decks: map[string]*Deck{
			"dk-a1b2": {
				Created: "2026-01-08",
				Nooks: map[string]*Nook{
					"a1b2-auth-spike": {Parent: "main", Created: "2026-01-08"},
				},
			},
			"dk-e5f6": {
				Created: "2026-01-09",
				Nooks: map[string]*Nook{
					"e5f6-feature": {Parent: "main", Created: "2026-01-09"},
				},
			},
		},
	}

	// Find nook in first deck
	deckID, found := GetDeckForNook(decks, "a1b2-auth-spike")
	if !found {
		t.Fatal("Expected to find nook")
	}
	if deckID != "dk-a1b2" {
		t.Errorf("Deck ID = %q, want %q", deckID, "dk-a1b2")
	}

	// Find nook in second deck
	deckID, found = GetDeckForNook(decks, "e5f6-feature")
	if !found {
		t.Fatal("Expected to find nook")
	}
	if deckID != "dk-e5f6" {
		t.Errorf("Deck ID = %q, want %q", deckID, "dk-e5f6")
	}

	// Non-existing nook
	_, found = GetDeckForNook(decks, "nonexistent")
	if found {
		t.Error("Should not find non-existing nook")
	}
}

func TestGetChildNooks(t *testing.T) {
	decks := &DecksFile{
		Decks: map[string]*Deck{
			"dk-a1b2": {
				Created: "2026-01-08",
				Nooks: map[string]*Nook{
					"a1b2-auth-spike":  {Parent: "main", Created: "2026-01-08"},
					"c3d4-jwt-variant": {Parent: "a1b2-auth-spike", Created: "2026-01-09"},
					"e5f6-refresh":     {Parent: "a1b2-auth-spike", Created: "2026-01-10"},
				},
			},
		},
	}

	// Get children of a1b2-auth-spike
	children := GetChildNooks(decks, "a1b2-auth-spike")
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}

	// Get children of main (should be one - a1b2-auth-spike)
	children = GetChildNooks(decks, "main")
	if len(children) != 1 {
		t.Errorf("Expected 1 child for main, got %d", len(children))
	}

	// Get children of leaf nook (should be none)
	children = GetChildNooks(decks, "c3d4-jwt-variant")
	if len(children) != 0 {
		t.Errorf("Expected 0 children for leaf, got %d", len(children))
	}
}
