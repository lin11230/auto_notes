package apple

import (
	"testing"
)

// TestEscapeAppleScriptString tests the escaping function
func TestEscapeAppleScriptString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No special characters",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Backslash",
			input:    `C:\Users\test`,
			expected: `C:\\Users\\test`,
		},
		{
			name:     "Double quotes",
			input:    `He said "Hello"`,
			expected: `He said \"Hello\"`,
		},
		{
			name:     "Single quotes",
			input:    `It's a test`,
			expected: `It\'s a test`,
		},
		{
			name:     "Mixed special characters",
			input:    `Path: "C:\test" and it's working`,
			expected: `Path: \"C:\\test\" and it\'s working`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeAppleScriptString(tt.input)
			if result != tt.expected {
				t.Errorf("escapeAppleScriptString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseAppleDate tests the date parsing function
func TestParseAppleDate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid date with AM",
			input:       "Wednesday, March 4, 2026 at 12:00:00 PM",
			expectError: false,
		},
		{
			name:        "Valid date with PM",
			input:       "Monday, January 1, 2024 at 3:45:30 PM",
			expectError: false,
		},
		{
			name:        "Invalid date format",
			input:       "2026-03-04",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAppleDate(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("parseAppleDate(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseAppleDate(%q) unexpected error: %v", tt.input, err)
				}
				if result.IsZero() {
					t.Errorf("parseAppleDate(%q) returned zero time", tt.input)
				}
			}
		})
	}
}

// TestTextToHTML tests the text to HTML conversion
func TestTextToHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Hello World",
			expected: "<p>Hello World</p>",
		},
		{
			name:     "Text with newline",
			input:    "Line 1\nLine 2",
			expected: "<p>Line 1<br>Line 2</p>",
		},
		{
			name:     "Multiple newlines",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "<p>Line 1<br>Line 2<br>Line 3</p>",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "<p></p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textToHTML(tt.input)
			if result != tt.expected {
				t.Errorf("textToHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNewNotesClient tests the client creation
func TestNewNotesClient(t *testing.T) {
	client := NewNotesClient()
	if client == nil {
		t.Error("NewNotesClient() returned nil")
	}
}
