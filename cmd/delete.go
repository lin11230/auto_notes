package cmd

import (
	"fmt"
	"os"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var deletePermanent bool

var deleteCmd = &cobra.Command{
	Use:   "delete <筆記名稱或ID>",
	Short: "刪除筆記",
	Long: `刪除指定的筆記。筆記會被移到「最近刪除」資料夾。

範例：
  notes delete "我的筆記"
  notes delete "x-coredata://..."`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		err := client.DeleteNote(identifier, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "錯誤：%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ 已將筆記移到「最近刪除」: %s\n", identifier)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deletePermanent, "permanent", "p", false, "永久刪除（不移到垃圾桶）")
}
