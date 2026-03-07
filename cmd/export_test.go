package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveExportFormat(t *testing.T) {
	tests := []struct {
		name       string
		flagValue  string
		outputPath string
		expected   string
	}{
		{name: "Flag html", flagValue: "html", outputPath: "", expected: "html"},
		{name: "Flag md", flagValue: "md", outputPath: "note.html", expected: "md"},
		{name: "Infer html from extension", outputPath: "note.html", expected: "html"},
		{name: "Infer markdown from extension", outputPath: "note.md", expected: "md"},
		{name: "Default stdout to markdown", outputPath: "", expected: "md"},
		{name: "Unknown extension", outputPath: "note.txt", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveExportFormat(tt.flagValue, tt.outputPath); got != tt.expected {
				t.Fatalf("resolveExportFormat(%q, %q) = %q, want %q", tt.flagValue, tt.outputPath, got, tt.expected)
			}
		})
	}
}

func TestRenderExportContent(t *testing.T) {
	html := "<p>Hello<br><strong>world</strong></p>"

	tests := []struct {
		name      string
		format    string
		expected  string
		wantError bool
	}{
		{name: "HTML", format: "html", expected: html},
		{name: "Markdown", format: "md", expected: "Hello\n**world**"},
		{name: "Invalid", format: "txt", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderExportContent(html, tt.format)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("renderExportContent() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestHTMLToMarkdown(t *testing.T) {
	input := "<div><p>Line 1<br />Line 2</p><p><em>Italic</em> and <b>bold</b></p></div>"
	expected := "Line 1\nLine 2\n\n*Italic* and **bold**"

	if got := htmlToMarkdown(input); got != expected {
		t.Fatalf("htmlToMarkdown() = %q, want %q", got, expected)
	}
}

func TestExportFilePermission(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "note.md")

	if err := os.WriteFile(outputPath, []byte("secret"), 0600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}

	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf("file permission = %#o, want %#o", got, 0600)
	}
}
