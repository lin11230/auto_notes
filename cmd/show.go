package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <筆記名稱或ID>",
	Short: "顯示筆記內容",
	Long: `顯示指定筆記的詳細內容。

範例：
  notes show "我的筆記"
  notes show "x-coredata://..."`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		note, err := client.ShowNote(identifier)
		if err != nil {
			fmt.Fprintf(os.Stderr, "錯誤：%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("標題: %s\n", note.Name)
		fmt.Printf("ID: %s\n", note.ID)
		fmt.Printf("資料夾: %s\n", note.Container)
		fmt.Printf("建立時間: %s\n", note.CreationDate.Format("2006-01-02 15:04:05"))
		fmt.Printf("修改時間: %s\n", note.ModificationDate.Format("2006-01-02 15:04:05"))
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
