package nook

import (
	"fmt"
	"regexp"
	"strings"
)

// Pre-compiled regex patterns for performance
var (
	// Match any character that is NOT alphanumeric or hyphen
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9-]`)
	// Match multiple consecutive hyphens
	multipleHyphensRegex = regexp.MustCompile(`-+`)
)

// NookError represents a nook operation error
type NookError struct {
	Code    string
	Message string
}

func (e *NookError) Error() string {
	return e.Message
}

// GenerateID creates a nook ID from commit SHA and user-provided name
// Format: {4-char-hash}-{sanitized-name}
// Example: "a1b2-auth-spike"
func GenerateID(name string, commitSHA string) (string, error) {
	if len(commitSHA) < 4 {
		return "", &NookError{
			Code:    "NOOK_INVALID_SHA",
			Message: "Commit SHA must be at least 4 characters",
		}
	}

	sanitized := SanitizeName(name)
	if sanitized == "" {
		return "", &NookError{
			Code:    "NOOK_NAME_INVALID",
			Message: "Nook name cannot be empty or contain only special characters",
		}
	}

	hash := strings.ToLower(commitSHA[:4])
	return fmt.Sprintf("%s-%s", hash, sanitized), nil
}

// SanitizeName cleans user input for nook naming
// Rules:
// 1. Convert to lowercase
// 2. Replace spaces with hyphens
// 3. Remove non-alphanumeric characters (except hyphens)
// 4. Collapse multiple consecutive hyphens to single hyphen
// 5. Trim leading/trailing hyphens
func SanitizeName(name string) string {
	// Step 1: Convert to lowercase
	result := strings.ToLower(name)

	// Step 2: Replace spaces with hyphens
	result = strings.ReplaceAll(result, " ", "-")

	// Step 3: Remove non-alphanumeric characters (except hyphens)
	result = nonAlphanumericRegex.ReplaceAllString(result, "")

	// Step 4: Collapse multiple consecutive hyphens
	result = multipleHyphensRegex.ReplaceAllString(result, "-")

	// Step 5: Trim leading/trailing hyphens
	result = strings.Trim(result, "-")

	return result
}

// ValidateName checks if a name would produce a valid nook ID
func ValidateName(name string) error {
	sanitized := SanitizeName(name)
	if sanitized == "" {
		return &NookError{
			Code:    "NOOK_NAME_INVALID",
			Message: "Nook name cannot be empty or contain only special characters",
		}
	}
	return nil
}

// ParseNookID extracts the hash and name parts from a nook ID
// Returns hash, name, and any error
func ParseNookID(nookID string) (hash string, name string, err error) {
	if len(nookID) < 6 { // minimum: 4-char hash + hyphen + 1-char name
		return "", "", &NookError{
			Code:    "NOOK_ID_INVALID",
			Message: fmt.Sprintf("Invalid nook ID format: %s", nookID),
		}
	}

	// Find the first hyphen after the 4-char hash
	if nookID[4] != '-' {
		return "", "", &NookError{
			Code:    "NOOK_ID_INVALID",
			Message: fmt.Sprintf("Invalid nook ID format: %s", nookID),
		}
	}

	hash = nookID[:4]
	name = nookID[5:]

	if name == "" {
		return "", "", &NookError{
			Code:    "NOOK_ID_INVALID",
			Message: fmt.Sprintf("Invalid nook ID format: %s", nookID),
		}
	}

	return hash, name, nil
}
