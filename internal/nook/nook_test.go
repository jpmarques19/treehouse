package nook

import (
	"testing"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "auth-spike",
			expected: "auth-spike",
		},
		{
			name:     "uppercase conversion",
			input:    "Auth-Spike",
			expected: "auth-spike",
		},
		{
			name:     "spaces to hyphens",
			input:    "my feature",
			expected: "my-feature",
		},
		{
			name:     "special characters removed",
			input:    "test@feature!",
			expected: "testfeature",
		},
		{
			name:     "mixed case with special chars",
			input:    "My Feature @123!",
			expected: "my-feature-123",
		},
		{
			name:     "multiple spaces",
			input:    "my   feature",
			expected: "my-feature",
		},
		{
			name:     "leading trailing special chars",
			input:    "@auth-spike!",
			expected: "auth-spike",
		},
		{
			name:     "all special characters",
			input:    "!@#$%",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "numbers preserved",
			input:    "feature123",
			expected: "feature123",
		},
		{
			name:     "complex mixed input",
			input:    "  @My--Feature!!  123  ",
			expected: "my-feature-123",
		},
		{
			name:     "underscores removed",
			input:    "my_feature",
			expected: "myfeature",
		},
		{
			name:     "dots removed",
			input:    "my.feature",
			expected: "myfeature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name       string
		nookName   string
		commitSHA  string
		expectedID string
		wantErr    bool
		errCode    string
	}{
		{
			name:       "standard generation",
			nookName:   "auth-spike",
			commitSHA:  "a1b2c3d4e5f6",
			expectedID: "a1b2-auth-spike",
			wantErr:    false,
		},
		{
			name:       "uppercase SHA converted",
			nookName:   "auth-spike",
			commitSHA:  "A1B2C3D4E5F6",
			expectedID: "a1b2-auth-spike",
			wantErr:    false,
		},
		{
			name:       "name with spaces",
			nookName:   "my feature",
			commitSHA:  "c3d4e5f6g7h8",
			expectedID: "c3d4-my-feature",
			wantErr:    false,
		},
		{
			name:       "name with special chars",
			nookName:   "test@feature!",
			commitSHA:  "e5f6g7h8i9j0",
			expectedID: "e5f6-testfeature",
			wantErr:    false,
		},
		{
			name:       "short SHA (exactly 4 chars)",
			nookName:   "spike",
			commitSHA:  "abcd",
			expectedID: "abcd-spike",
			wantErr:    false,
		},
		{
			name:      "SHA too short",
			nookName:  "spike",
			commitSHA: "abc",
			wantErr:   true,
			errCode:   "NOOK_INVALID_SHA",
		},
		{
			name:      "empty name",
			nookName:  "",
			commitSHA: "a1b2c3d4",
			wantErr:   true,
			errCode:   "NOOK_NAME_INVALID",
		},
		{
			name:      "name with only special chars",
			nookName:  "!@#$%",
			commitSHA: "a1b2c3d4",
			wantErr:   true,
			errCode:   "NOOK_NAME_INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateID(tt.nookName, tt.commitSHA)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateID(%q, %q) expected error, got nil", tt.nookName, tt.commitSHA)
					return
				}
				nookErr, ok := err.(*NookError)
				if !ok {
					t.Errorf("Expected NookError, got %T", err)
					return
				}
				if nookErr.Code != tt.errCode {
					t.Errorf("Error code = %q, want %q", nookErr.Code, tt.errCode)
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateID(%q, %q) unexpected error: %v", tt.nookName, tt.commitSHA, err)
				return
			}

			if result != tt.expectedID {
				t.Errorf("GenerateID(%q, %q) = %q, want %q", tt.nookName, tt.commitSHA, result, tt.expectedID)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid simple name",
			input:   "auth-spike",
			wantErr: false,
		},
		{
			name:    "valid name with spaces",
			input:   "my feature",
			wantErr: false,
		},
		{
			name:    "valid name with special chars",
			input:   "test@feature!",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errCode: "NOOK_NAME_INVALID",
		},
		{
			name:    "only special chars",
			input:   "!@#$%",
			wantErr: true,
			errCode: "NOOK_NAME_INVALID",
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: true,
			errCode: "NOOK_NAME_INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateName(%q) expected error, got nil", tt.input)
					return
				}
				nookErr, ok := err.(*NookError)
				if !ok {
					t.Errorf("Expected NookError, got %T", err)
					return
				}
				if nookErr.Code != tt.errCode {
					t.Errorf("Error code = %q, want %q", nookErr.Code, tt.errCode)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateName(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestParseNookID(t *testing.T) {
	tests := []struct {
		name         string
		nookID       string
		expectedHash string
		expectedName string
		wantErr      bool
		errCode      string
	}{
		{
			name:         "valid nook ID",
			nookID:       "a1b2-auth-spike",
			expectedHash: "a1b2",
			expectedName: "auth-spike",
			wantErr:      false,
		},
		{
			name:         "valid nook ID with numbers",
			nookID:       "c3d4-feature-123",
			expectedHash: "c3d4",
			expectedName: "feature-123",
			wantErr:      false,
		},
		{
			name:         "minimum valid ID",
			nookID:       "abcd-x",
			expectedHash: "abcd",
			expectedName: "x",
			wantErr:      false,
		},
		{
			name:    "too short",
			nookID:  "a1b2",
			wantErr: true,
			errCode: "NOOK_ID_INVALID",
		},
		{
			name:    "missing hyphen",
			nookID:  "a1b2authspike",
			wantErr: true,
			errCode: "NOOK_ID_INVALID",
		},
		{
			name:    "empty name part",
			nookID:  "a1b2-",
			wantErr: true,
			errCode: "NOOK_ID_INVALID",
		},
		{
			name:    "empty string",
			nookID:  "",
			wantErr: true,
			errCode: "NOOK_ID_INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, name, err := ParseNookID(tt.nookID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseNookID(%q) expected error, got nil", tt.nookID)
					return
				}
				nookErr, ok := err.(*NookError)
				if !ok {
					t.Errorf("Expected NookError, got %T", err)
					return
				}
				if nookErr.Code != tt.errCode {
					t.Errorf("Error code = %q, want %q", nookErr.Code, tt.errCode)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseNookID(%q) unexpected error: %v", tt.nookID, err)
				return
			}

			if hash != tt.expectedHash {
				t.Errorf("ParseNookID(%q) hash = %q, want %q", tt.nookID, hash, tt.expectedHash)
			}
			if name != tt.expectedName {
				t.Errorf("ParseNookID(%q) name = %q, want %q", tt.nookID, name, tt.expectedName)
			}
		})
	}
}
