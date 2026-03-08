package cmd

import (
	"fmt"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var deletePermanent bool

var deleteCmd = &cobra.Command{
	Use:   "delete <note-title-or-id>",
	Short: "Delete a note",
	Long: `Delete the specified note. The note will be moved to Recently Deleted.

Examples:
  notes delete "My Note"
  notes delete "x-coredata://..."`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		err := client.DeleteNote(identifier, false)
		if err != nil {
			exitWithError("unable to delete note", err)
		}

		fmt.Printf("Moved note to Recently Deleted: %s\n", identifier)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deletePermanent, "permanent", "p", false, "Permanently delete instead of moving to Recently Deleted")
}
