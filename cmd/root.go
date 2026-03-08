package cmd

import (
	"fmt"
	"os"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var appVersion = "dev"
var debugMode bool

var rootCmd = &cobra.Command{
	Use:   "notes",
	Short: "CLI tool for managing Apple Notes",
	Long: `notes is a command-line tool for managing the macOS Notes app through AppleScript.

Supported features:
  • Create notes
  • List notes
  • Show note details
  • Search notes
  • Delete notes
  • Export notes
  • Manage folders`,
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
	rootCmd.Flags().BoolP("help", "h", false, "Show help")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Show detailed error output")
	rootCmd.SetVersionTemplate(fmt.Sprintf("notes version %s\n", appVersion))
	cobra.OnInitialize(func() {
		apple.SetDebug(debugMode)
	})
}
