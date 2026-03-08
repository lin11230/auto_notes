package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Manage folders",
	Long: `Manage Apple Notes folders.

Subcommands:
  list    List all folders
  create  Create a new folder`,
}

var folderListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all folders",
	Run: func(cmd *cobra.Command, args []string) {
		client := apple.NewNotesClient()
		folders, err := client.ListFolders()
		if err != nil {
			exitWithError("unable to list folders", err)
		}

		if len(folders) == 0 {
			fmt.Println("No folders found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Name")
		fmt.Fprintln(w, "----")
		for _, folder := range folders {
			fmt.Fprintf(w, "%s\n", folder.Name)
		}
		w.Flush()
		fmt.Printf("\nTotal folders: %d\n", len(folders))
	},
}

var folderCreateName string

var folderCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new folder",
	Long: `Create a new folder.

Examples:
  notes folder create -n "Work"
  notes folder create --name "Personal"`,
	Run: func(cmd *cobra.Command, args []string) {
		if folderCreateName == "" {
			exitWithError("please provide a folder name (-n, --name)", nil)
		}

		client := apple.NewNotesClient()
		folder, err := client.CreateFolder(folderCreateName)
		if err != nil {
			exitWithError("unable to create folder", err)
		}

		fmt.Printf("Created folder successfully\n")
		fmt.Printf("  ID: %s\n", folder.ID)
		fmt.Printf("  Name: %s\n", folder.Name)
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
	folderCmd.AddCommand(folderListCmd)
	folderCmd.AddCommand(folderCreateCmd)

	folderCreateCmd.Flags().StringVarP(&folderCreateName, "name", "n", "", "Folder name (required)")
	folderCreateCmd.MarkFlagRequired("name")
}
