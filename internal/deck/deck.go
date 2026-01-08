package deck

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DecksFile represents the structure of decks.yaml
type DecksFile struct {
	Decks map[string]*Deck `yaml:"decks"`
}

// Deck represents a deck containing nooks
type Deck struct {
	Created string           `yaml:"created"`
	Nooks   map[string]*Nook `yaml:"nooks"`
}

// Nook represents a nook entry in a deck
type Nook struct {
	Parent  string `yaml:"parent"`
	Created string `yaml:"created"`
}

// DeckError represents a deck operation error
type DeckError struct {
	Code    string
	Message string
}

func (e *DeckError) Error() string {
	return e.Message
}

// LoadDecks reads and parses the decks.yaml file
func LoadDecks(treehousePath string) (*DecksFile, error) {
	decksPath := filepath.Join(treehousePath, "decks.yaml")

	data, err := os.ReadFile(decksPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &DeckError{
				Code:    "DECK_FILE_NOT_FOUND",
				Message: fmt.Sprintf("decks.yaml not found at %s", decksPath),
			}
		}
		return nil, &DeckError{
			Code:    "DECK_READ_FAILED",
			Message: fmt.Sprintf("Failed to read decks.yaml: %v", err),
		}
	}

	var decksFile DecksFile
	if err := yaml.Unmarshal(data, &decksFile); err != nil {
		return nil, &DeckError{
			Code:    "DECK_PARSE_FAILED",
			Message: fmt.Sprintf("Failed to parse decks.yaml: %v", err),
		}
	}

	// Initialize empty maps if nil
	if decksFile.Decks == nil {
		decksFile.Decks = make(map[string]*Deck)
	}

	return &decksFile, nil
}

// NookExists checks if a nook with the given ID exists in any deck
func NookExists(treehousePath string, nookID string) (bool, error) {
	decks, err := LoadDecks(treehousePath)
	if err != nil {
		// If decks.yaml doesn't exist, no nooks exist
		if deckErr, ok := err.(*DeckError); ok && deckErr.Code == "DECK_FILE_NOT_FOUND" {
			return false, nil
		}
		return false, err
	}

	for _, deck := range decks.Decks {
		if deck.Nooks != nil {
			if _, exists := deck.Nooks[nookID]; exists {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetNook returns a nook by ID, searching all decks
func GetNook(treehousePath string, nookID string) (*Nook, string, error) {
	decks, err := LoadDecks(treehousePath)
	if err != nil {
		return nil, "", err
	}

	for deckID, deck := range decks.Decks {
		if deck.Nooks != nil {
			if nook, exists := deck.Nooks[nookID]; exists {
				return nook, deckID, nil
			}
		}
	}

	return nil, "", &DeckError{
		Code:    "NOOK_NOT_FOUND",
		Message: fmt.Sprintf("Nook '%s' not found", nookID),
	}
}

// SaveDecks writes the decks file atomically (temp file + rename)
func SaveDecks(treehousePath string, decks *DecksFile) error {
	decksPath := filepath.Join(treehousePath, "decks.yaml")

	// Marshal to YAML
	data, err := yaml.Marshal(decks)
	if err != nil {
		return &DeckError{
			Code:    "DECK_MARSHAL_FAILED",
			Message: fmt.Sprintf("Failed to marshal decks: %v", err),
		}
	}

	// Create temp file in same directory for atomic rename
	tempFile, err := os.CreateTemp(treehousePath, "decks.yaml.tmp.*")
	if err != nil {
		return &DeckError{
			Code:    "DECK_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to create temp file: %v", err),
		}
	}
	tempPath := tempFile.Name()

	// Ensure temp file cleanup on error
	defer func() {
		if tempPath != "" {
			os.Remove(tempPath)
		}
	}()

	// Write data to temp file
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return &DeckError{
			Code:    "DECK_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to write temp file: %v", err),
		}
	}

	// Sync to disk
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return &DeckError{
			Code:    "DECK_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to sync temp file: %v", err),
		}
	}

	// Close before rename
	if err := tempFile.Close(); err != nil {
		return &DeckError{
			Code:    "DECK_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to close temp file: %v", err),
		}
	}

	// Atomic rename
	if err := os.Rename(tempPath, decksPath); err != nil {
		return &DeckError{
			Code:    "DECK_WRITE_FAILED",
			Message: fmt.Sprintf("Failed to rename temp file: %v", err),
		}
	}

	// Clear tempPath to prevent cleanup of the renamed file
	tempPath = ""

	return nil
}

// GenerateDeckID creates a deck ID from a nook's 4-char hash
// Format: dk-{hash}
func GenerateDeckID(nookHash string) string {
	return fmt.Sprintf("dk-%s", nookHash)
}

// CreateDeck creates a new deck with the given ID and creation date
func CreateDeck(deckID string, created string) *Deck {
	return &Deck{
		Created: created,
		Nooks:   make(map[string]*Nook),
	}
}

// AddNookToDeck adds a nook to a deck, creating the deck if it doesn't exist
// Returns the deck ID (which may be new if a deck was created)
func AddNookToDeck(decks *DecksFile, deckID string, nookID string, parent string, created string) string {
	// Initialize decks map if nil
	if decks.Decks == nil {
		decks.Decks = make(map[string]*Deck)
	}

	// Create deck if it doesn't exist
	deck, exists := decks.Decks[deckID]
	if !exists {
		deck = CreateDeck(deckID, created)
		decks.Decks[deckID] = deck
	}

	// Initialize nooks map if nil
	if deck.Nooks == nil {
		deck.Nooks = make(map[string]*Nook)
	}

	// Add the nook
	deck.Nooks[nookID] = &Nook{
		Parent:  parent,
		Created: created,
	}

	return deckID
}

// RemoveNook removes a nook from a deck
// Returns true if the deck is now empty (and should be removed)
func RemoveNook(decks *DecksFile, deckID string, nookID string) (bool, error) {
	deck, exists := decks.Decks[deckID]
	if !exists {
		return false, &DeckError{
			Code:    "DECK_NOT_FOUND",
			Message: fmt.Sprintf("Deck '%s' not found", deckID),
		}
	}

	if deck.Nooks == nil {
		return true, nil
	}

	delete(deck.Nooks, nookID)

	// Check if deck is now empty
	if len(deck.Nooks) == 0 {
		delete(decks.Decks, deckID)
		return true, nil
	}

	return false, nil
}

// GetDeckForNook finds the deck ID for a given nook
func GetDeckForNook(decks *DecksFile, nookID string) (string, bool) {
	for deckID, deck := range decks.Decks {
		if deck.Nooks != nil {
			if _, exists := deck.Nooks[nookID]; exists {
				return deckID, true
			}
		}
	}
	return "", false
}

// GetChildNooks returns all nooks that have the given nook as their parent
func GetChildNooks(decks *DecksFile, parentNookID string) []string {
	var children []string
	for _, deck := range decks.Decks {
		if deck.Nooks != nil {
			for nookID, nook := range deck.Nooks {
				if nook.Parent == parentNookID {
					children = append(children, nookID)
				}
			}
		}
	}
	return children
}
