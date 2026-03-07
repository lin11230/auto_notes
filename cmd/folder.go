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
	Short: "管理資料夾",
	Long: `管理 Apple Notes 資料夾。

子指令：
  list    列出所有資料夾
  create  建立新資料夾`,
}

var folderListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出所有資料夾",
	Run: func(cmd *cobra.Command, args []string) {
		client := apple.NewNotesClient()
		folders, err := client.ListFolders()
		if err != nil {
			exitWithError("無法列出資料夾", err)
		}

		if len(folders) == 0 {
			fmt.Println("沒有找到任何資料夾")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "名稱")
		fmt.Fprintln(w, "────")
		for _, folder := range folders {
			fmt.Fprintf(w, "%s\n", folder.Name)
		}
		w.Flush()
		fmt.Printf("\n共 %d 個資料夾\n", len(folders))
	},
}

var folderCreateName string

var folderCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "建立新資料夾",
	Long: `建立新的資料夾。

範例：
  notes folder create -n "工作"
  notes folder create --name "個人"`,
	Run: func(cmd *cobra.Command, args []string) {
		if folderCreateName == "" {
			exitWithError("請提供資料夾名稱 (-n, --name)", nil)
		}

		client := apple.NewNotesClient()
		folder, err := client.CreateFolder(folderCreateName)
		if err != nil {
			exitWithError("無法建立資料夾", err)
		}

		fmt.Printf("✓ 已建立資料夾\n")
		fmt.Printf("  ID: %s\n", folder.ID)
		fmt.Printf("  名稱: %s\n", folder.Name)
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
	folderCmd.AddCommand(folderListCmd)
	folderCmd.AddCommand(folderCreateCmd)

	folderCreateCmd.Flags().StringVarP(&folderCreateName, "name", "n", "", "資料夾名稱 (必填)")
	folderCreateCmd.MarkFlagRequired("name")
}
