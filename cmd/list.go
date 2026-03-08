package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var (
	listFolder string
	listLimit  int
	listOffset int
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List notes",
	Long: `List all Apple Notes or notes in a specific folder.

Examples:
  notes list
  notes list --folder "Work"
  notes list --limit 10 --offset 20
  notes ls`,
	Run: func(cmd *cobra.Command, args []string) {
		client := apple.NewNotesClient()

		var notes []apple.Note
		var err error

		if listLimit > 0 {
			notes, err = client.ListNotesPaginated(listFolder, listLimit, listOffset)
		} else {
			notes, err = client.ListNotes(listFolder)
		}

		if err != nil {
			exitWithError("unable to list notes", err)
		}

		if len(notes) == 0 {
			fmt.Println("No notes found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTitle\tFolder\tModified")
		fmt.Fprintln(w, "--\t-----\t------\t--------")
		for _, note := range notes {
			modTime := note.ModificationDate.Format("2006-01-02 15:04")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", note.ID, truncate(note.Name, 30), note.Container, modTime)
		}
		w.Flush()

		if listLimit > 0 {
			fmt.Printf("\nShowing %d notes (starting from %d)\n", len(notes), listOffset+1)
		} else {
			fmt.Printf("\nTotal notes: %d\n", len(notes))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listFolder, "folder", "F", "", "Folder name")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 0, "Maximum number of notes to return")
	listCmd.Flags().IntVarP(&listOffset, "offset", "o", 0, "Number of notes to skip")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
