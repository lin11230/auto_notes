package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var exportOutput string

var exportCmd = &cobra.Command{
	Use:   "export <筆記名稱或ID>",
	Short: "匯出筆記",
	Long: `將筆記內容匯出為檔案。預設輸出到標準輸出，可使用 -o 指定檔案。

範例：
  notes export "我的筆記"
  notes export "我的筆記" -o note.txt
  notes export "我的筆記" -o note.html`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		name, body, err := client.ExportNote(identifier)
		if err != nil {
			fmt.Fprintf(os.Stderr, "錯誤：%v\n", err)
			os.Exit(1)
		}

		content := body
		if !strings.HasSuffix(exportOutput, ".html") {
			content = stripHTML(body)
		}

		if exportOutput != "" {
			err := ioutil.WriteFile(exportOutput, []byte(content), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "錯誤：無法寫入檔案: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ 已匯出筆記「%s」到 %s\n", name, exportOutput)
		} else {
			fmt.Println(content)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "輸出檔案路徑")
}
