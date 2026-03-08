package cmd

import (
	"fmt"
	"strings"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <note-title-or-id>",
	Short: "Show note details",
	Long: `Show the full details of a note.

Examples:
  notes show "My Note"
  notes show "x-coredata://..."`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		note, err := client.ShowNote(identifier)
		if err != nil {
			exitWithError("unable to show note", err)
		}

		fmt.Printf("Title: %s\n", note.Name)
		fmt.Printf("ID: %s\n", note.ID)
		fmt.Printf("Folder: %s\n", note.Container)
		fmt.Printf("Created: %s\n", note.CreationDate.Format("2006-01-02 15:04:05"))
		fmt.Printf("Modified: %s\n", note.ModificationDate.Format("2006-01-02 15:04:05"))
		fmt.Println(strings.Repeat("-", 40))
		fmt.Println(stripHTML(note.Body))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func stripHTML(html string) string {
	// Simple HTML tag removal
	result := html
	result = strings.ReplaceAll(result, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")
	result = strings.ReplaceAll(result, "</p>", "\n")
	result = strings.ReplaceAll(result, "</div>", "\n")

	// Remove remaining HTML tags
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

	return strings.TrimSpace(clean.String())
}
