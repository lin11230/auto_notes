package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var exportOutput string
var exportFormat string

var exportCmd = &cobra.Command{
	Use:   "export <筆記名稱或ID>",
	Short: "匯出筆記",
	Long: `將筆記內容匯出為檔案。預設輸出到標準輸出，可使用 -o 指定檔案。

範例：
  notes export "我的筆記"
  notes export "我的筆記" --format md -o note.md
  notes export "我的筆記" --format html -o note.html`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		client := apple.NewNotesClient()
		name, body, err := client.ExportNote(identifier)
		if err != nil {
			exitWithError("無法匯出筆記", err)
		}

		format := resolveExportFormat(exportFormat, exportOutput)
		if format == "" {
			exitWithError("匯出格式只支援 html 或 md", nil)
		}

		content, err := renderExportContent(body, format)
		if err != nil {
			exitWithError("無法轉換匯出內容", err)
		}

		if exportOutput != "" {
			err := os.WriteFile(exportOutput, []byte(content), 0600)
			if err != nil {
				exitWithError("無法寫入輸出檔案", err)
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
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "", "匯出格式：html 或 md")
}

func resolveExportFormat(flagValue, outputPath string) string {
	if flagValue != "" {
		return strings.ToLower(flagValue)
	}

	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".html", ".htm":
		return "html"
	case ".md", ".markdown":
		return "md"
	default:
		if outputPath == "" {
			return "md"
		}
		return ""
	}
}

func renderExportContent(body, format string) (string, error) {
	switch format {
	case "html":
		return body, nil
	case "md":
		return htmlToMarkdown(body), nil
	default:
		return "", fmt.Errorf("不支援的匯出格式: %s", format)
	}
}

func htmlToMarkdown(html string) string {
	replacer := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n\n",
		"<p>", "",
		"</div>", "\n",
		"<div>", "",
		"<strong>", "**",
		"</strong>", "**",
		"<b>", "**",
		"</b>", "**",
		"<em>", "*",
		"</em>", "*",
		"<i>", "*",
		"</i>", "*",
		"&nbsp;", " ",
	)

	result := replacer.Replace(html)

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

	lines := strings.Split(clean.String(), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
