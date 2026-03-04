package cmd

import (
	"fmt"
	"io/ioutil"
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
	Short: "建立新筆記",
	Long: `建立一個新的 Apple Note。

範例：
  notes create -t "我的筆記" -b "這是筆記內容"
  notes create -t "我的筆記" -f content.txt
  notes create -t "我的筆記" -b "內容" --folder "工作"`,
	Run: func(cmd *cobra.Command, args []string) {
		if createTitle == "" {
			fmt.Fprintln(os.Stderr, "錯誤：請提供筆記標題 (-t, --title)")
			os.Exit(1)
		}

		body := createBody
		if createFile != "" {
			content, err := ioutil.ReadFile(createFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "錯誤：無法讀取檔案 %s: %v\n", createFile, err)
				os.Exit(1)
			}
			body = string(content)
		}

		client := apple.NewNotesClient()
		note, err := client.CreateNote(createTitle, body, createFolder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "錯誤：無法建立筆記: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ 已建立筆記\n")
		fmt.Printf("  ID: %s\n", note.ID)
		fmt.Printf("  標題: %s\n", note.Name)
		if createFolder != "" {
			fmt.Printf("  資料夾: %s\n", createFolder)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&createTitle, "title", "t", "", "筆記標題 (必填)")
	createCmd.Flags().StringVarP(&createBody, "body", "b", "", "筆記內容")
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "從檔案讀取內容")
	createCmd.Flags().StringVarP(&createFolder, "folder", "F", "", "指定資料夾")
	createCmd.MarkFlagRequired("title")
}
