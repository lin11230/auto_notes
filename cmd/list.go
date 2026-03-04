package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var listFolder string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出所有筆記",
	Long: `列出所有 Apple Notes 或指定資料夾中的筆記。

範例：
  notes list
  notes list --folder "工作"
  notes ls`,
	Run: func(cmd *cobra.Command, args []string) {
		client := apple.NewNotesClient()
		notes, err := client.ListNotes(listFolder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "錯誤：無法列出筆記: %v\n", err)
			os.Exit(1)
		}

		if len(notes) == 0 {
			fmt.Println("沒有找到任何筆記")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\t標題\t資料夾\t修改時間")
		fmt.Fprintln(w, "──\t────\t──────\t────────")
		for _, note := range notes {
			modTime := note.ModificationDate.Format("2006-01-02 15:04")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", note.ID, truncate(note.Name, 30), note.Container, modTime)
		}
		w.Flush()
		fmt.Printf("\n共 %d 則筆記\n", len(notes))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listFolder, "folder", "F", "", "指定資料夾")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
