package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var moveTargetFolder string

var moveCmd = &cobra.Command{
	Use:   "move <note-title-or-id>... -t <target-folder>",
	Short: "Move notes to another folder",
	Long: `Move one or more notes to a target folder.

Examples:
  notes move "Meeting Notes" -t "Work"
  notes move "Note1" "Note2" "Note3" -t "Work"
  notes move "x-coredata://..." -t "Personal"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if moveTargetFolder == "" {
			exitWithError("please specify a target folder (-t, --to)", nil)
		}

		client := apple.NewNotesClient()

		// Step 1: 驗證目標資料夾是否存在
		folders, err := client.ListFolders()
		if err != nil {
			exitWithError("unable to list folders", err)
		}

		folderExists := false
		for _, folder := range folders {
			if folder.Name == moveTargetFolder {
				folderExists = true
				break
			}
		}

		if !folderExists {
			fmt.Fprintf(os.Stderr, "Error: target folder %q does not exist\n", moveTargetFolder)
			fmt.Fprintln(os.Stderr, "\nAvailable folders:")
			for _, folder := range folders {
				fmt.Fprintf(os.Stderr, "  - %s\n", folder.Name)
			}
			os.Exit(1)
		}

		// Step 2: 處理每個筆記
		successCount := 0
		failCount := 0
		hasAmbiguousNote := false

		for _, identifier := range args {
			// 如果是 ID，直接移動
			if strings.HasPrefix(identifier, "x-coredata://") {
				sourceFolder, targetFolder, err := client.MoveNote(identifier, moveTargetFolder)
				if err != nil {
					if debugMode {
						fmt.Fprintf(os.Stderr, "Failed to move note %s\nDetails: %v\n", identifier, err)
					} else {
						fmt.Fprintf(os.Stderr, "Failed to move note %s\n", identifier)
					}
					failCount++
					continue
				}
				fmt.Printf("Moved note from %q to %q\n", sourceFolder, targetFolder)
				successCount++
				continue
			}

			// 如果是名稱，檢查是否有多個同名筆記
			notes, err := client.FindNotesByName(identifier)
			if err != nil {
				if debugMode {
					fmt.Fprintf(os.Stderr, "Failed to look up note %q\nDetails: %v\n", identifier, err)
				} else {
					fmt.Fprintf(os.Stderr, "Failed to look up note %q\n", identifier)
				}
				failCount++
				continue
			}

			if len(notes) == 0 {
				fmt.Fprintf(os.Stderr, "Note %q not found\n", identifier)
				failCount++
				continue
			}

			if len(notes) > 1 {
				// 有多個同名筆記，列出並要求使用 ID
				fmt.Fprintf(os.Stderr, "\nError: found %d notes named %q. Please use an ID instead:\n\n", len(notes), identifier)

				w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "#\tTitle\tFolder\tID")
				fmt.Fprintln(w, "-\t-----\t------\t--")
				for i, note := range notes {
					fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i+1, note.Name, note.Container, note.ID)
				}
				w.Flush()

				fmt.Fprintf(os.Stderr, "\nUse this command instead:\n")
				fmt.Fprintf(os.Stderr, "  notes move \"<ID>\" -t \"%s\"\n\n", moveTargetFolder)

				failCount++
				hasAmbiguousNote = true
				continue
			}

			// 唯一筆記，直接移動
			note := notes[0]
			sourceFolder, targetFolder, err := client.MoveNote(note.ID, moveTargetFolder)
			if err != nil {
				if debugMode {
					fmt.Fprintf(os.Stderr, "Failed to move note %q\nDetails: %v\n", identifier, err)
				} else {
					fmt.Fprintf(os.Stderr, "Failed to move note %q\n", identifier)
				}
				failCount++
				continue
			}
			fmt.Printf("Moved note %q from %q to %q\n", note.Name, sourceFolder, targetFolder)
			successCount++
		}

		// Step 3: 顯示統計結果
		fmt.Println()
		if failCount == 0 {
			if successCount == 1 {
				fmt.Println("Moved 1 note successfully")
			} else {
				fmt.Printf("Moved %d notes successfully\n", successCount)
			}
		} else {
			fmt.Printf("Moved %d notes, failed to move %d\n", successCount, failCount)
			if hasAmbiguousNote {
				os.Exit(1)
			}
		}

		if failCount > 0 && !hasAmbiguousNote {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(&moveTargetFolder, "to", "t", "", "Target folder name (required)")
	moveCmd.MarkFlagRequired("to")
}
