package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var searchFolder string

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search notes",
	Long: `Search notes by keyword in title or body.

Examples:
  notes search "important"
  notes search "meeting" --folder "Work"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		client := apple.NewNotesClient()
		notes, err := client.SearchNotes(keyword, searchFolder)
		if err != nil {
			exitWithError("unable to search notes", err)
		}

		if len(notes) == 0 {
			fmt.Println("No matching notes found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Title\tFolder\tModified")
		fmt.Fprintln(w, "-----\t------\t--------")
		for _, note := range notes {
			modTime := note.ModificationDate.Format("2006-01-02 15:04")
			fmt.Fprintf(w, "%s\t%s\t%s\n", truncate(note.Name, 30), note.Container, modTime)
		}
		w.Flush()
		fmt.Printf("\nFound %d matching notes\n", len(notes))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&searchFolder, "folder", "F", "", "Search within a folder")
}
