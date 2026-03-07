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
	Use:   "search <關鍵字>",
	Short: "搜尋筆記",
	Long: `搜尋包含關鍵字的筆記（標題或內容）。

範例：
  notes search "重要"
  notes search "會議" --folder "工作"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		client := apple.NewNotesClient()
		notes, err := client.SearchNotes(keyword, searchFolder)
		if err != nil {
			exitWithError("無法搜尋筆記", err)
		}

		if len(notes) == 0 {
			fmt.Println("沒有找到匹配的筆記")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "標題\t資料夾\t修改時間")
		fmt.Fprintln(w, "────\t──────\t────────")
		for _, note := range notes {
			modTime := note.ModificationDate.Format("2006-01-02 15:04")
			fmt.Fprintf(w, "%s\t%s\t%s\n", truncate(note.Name, 30), note.Container, modTime)
		}
		w.Flush()
		fmt.Printf("\n找到 %d 則匹配的筆記\n", len(notes))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&searchFolder, "folder", "F", "", "指定搜尋的資料夾")
}
