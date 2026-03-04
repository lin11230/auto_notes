package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var appVersion = "dev"

var rootCmd = &cobra.Command{
	Use:   "notes",
	Short: "CLI 工具用於管理 Apple Notes",
	Long: `notes 是一個命令列工具，透過 AppleScript 管理 macOS 本機的 Notes 應用程式。

支援的功能：
  • 建立新筆記
  • 列出筆記
  • 查看筆記內容
  • 搜尋筆記
  • 刪除筆記
  • 匯出筆記
  • 管理資料夾`,
}

func SetVersion(v string) {
	appVersion = v
	rootCmd.Version = appVersion
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "顯示說明訊息")
	rootCmd.SetVersionTemplate(fmt.Sprintf("notes version %s\n", appVersion))
}
