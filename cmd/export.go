package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var exportOutput string
var exportFormat string

var exportCmd = &cobra.Command{
	Use:   "export <note-title-or-id>",
	Short: "Export a note",
	Long: `Export note content to a file. By default, content is printed to stdout unless -o is provided.

Examples:
  notes export "My Note"
  notes export "My Note" --format md -o note.md
  notes export "My Note" --format html -o note.html`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		name, body, err := client.ExportNote(identifier)
		if err != nil {
			exitWithError("unable to export note", err)
		}

		format := resolveExportFormat(exportFormat, exportOutput)
		if format == "" {
			exitWithError("export format only supports html or md", nil)
		}

		content, err := renderExportContent(body, format)
		if err != nil {
			exitWithError("unable to render export content", err)
		}

		if exportOutput != "" {
			err := os.WriteFile(exportOutput, []byte(content), 0600)
			if err != nil {
				exitWithError("unable to write output file", err)
			}
			fmt.Printf("Exported note %q to %s\n", name, exportOutput)
		} else {
			fmt.Println(content)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "", "Export format: html or md")
}

func resolveExportFormat(flagValue, outputPath string) string {
	if flagValue != "" {
		return strings.ToLower(flagValue)
	}

	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".html", ".htm":
		return "html"
	case ".md", ".markdown":
		return "md"
	default:
		if outputPath == "" {
			return "md"
		}
		return ""
	}
}

func renderExportContent(body, format string) (string, error) {
	switch format {
	case "html":
		return body, nil
	case "md":
		return htmlToMarkdown(body), nil
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}

func htmlToMarkdown(html string) string {
	replacer := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n\n",
		"<p>", "",
		"</div>", "\n",
		"<div>", "",
		"<strong>", "**",
		"</strong>", "**",
		"<b>", "**",
		"</b>", "**",
		"<em>", "*",
		"</em>", "*",
		"<i>", "*",
		"</i>", "*",
		"&nbsp;", " ",
	)

	result := replacer.Replace(html)

	inTag := false
	var clean strings.Builder
	for _, r := range result {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			clean.WriteRune(r)
		}
	}

	lines := strings.Split(clean.String(), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
