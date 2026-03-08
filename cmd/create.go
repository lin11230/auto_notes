package cmd

import (
	"fmt"
	"os"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var (
	createTitle  string
	createBody   string
	createFile   string
	createFolder string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	Long: `Create a new Apple Note.

Examples:
  notes create -t "My Note" -b "This is the note body"
  notes create -t "My Note" -f content.txt
  notes create -t "My Note" -b "Body" --folder "Work"`,
	Run: func(cmd *cobra.Command, args []string) {
		if createTitle == "" {
			exitWithError("please provide a note title (-t, --title)", nil)
		}

		body := createBody
		if createFile != "" {
			content, err := os.ReadFile(createFile)
			if err != nil {
				exitWithError("unable to read input file", err)
			}
			body = string(content)
		}

		client := apple.NewNotesClient()
		note, err := client.CreateNote(createTitle, body, createFolder)
		if err != nil {
			exitWithError("unable to create note", err)
		}

		fmt.Printf("Created note successfully\n")
		fmt.Printf("  ID: %s\n", note.ID)
		fmt.Printf("  Title: %s\n", note.Name)
		if createFolder != "" {
			fmt.Printf("  Folder: %s\n", createFolder)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&createTitle, "title", "t", "", "Note title (required)")
	createCmd.Flags().StringVarP(&createBody, "body", "b", "", "Note body")
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "Read content from file")
	createCmd.Flags().StringVarP(&createFolder, "folder", "F", "", "Target folder")
	createCmd.MarkFlagRequired("title")
}
