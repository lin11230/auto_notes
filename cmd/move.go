package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/kclin/auto_notes/internal/apple"
	"github.com/spf13/cobra"
)

var moveTargetFolder string

var moveCmd = &cobra.Command{
	Use:   "move <筆記名稱或ID>... -t <目標資料夾>",
	Short: "移動筆記到指定資料夾",
	Long: `移動一個或多個筆記到指定的資料夾。

範例：
  # 移動單個筆記
  notes move "會議記錄" -t "工作"
  
  # 批次移動多個筆記
  notes move "筆記1" "筆記2" "筆記3" -t "工作"
  
  # 使用 ID 移動（避免同名衝突）
  notes move "x-coredata://..." -t "個人"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if moveTargetFolder == "" {
			fmt.Fprintln(os.Stderr, "錯誤：請指定目標資料夾 (-t, --to)")
			os.Exit(1)
		}

		client := apple.NewNotesClient()

		// Step 1: 驗證目標資料夾是否存在
		folders, err := client.ListFolders()
		if err != nil {
			fmt.Fprintf(os.Stderr, "錯誤：無法列出資料夾: %v\n", err)
			os.Exit(1)
		}

		folderExists := false
		for _, folder := range folders {
			if folder.Name == moveTargetFolder {
				folderExists = true
				break
			}
		}

		if !folderExists {
			fmt.Fprintf(os.Stderr, "錯誤：目標資料夾「%s」不存在\n", moveTargetFolder)
			fmt.Fprintln(os.Stderr, "\n可用的資料夾：")
			for _, folder := range folders {
				fmt.Fprintf(os.Stderr, "  - %s\n", folder.Name)
			}
			os.Exit(1)
		}

		// Step 2: 處理每個筆記
		successCount := 0
		failCount := 0
		hasAmbiguousNote := false

		for _, identifier := range args {
			// 如果是 ID，直接移動
			if strings.HasPrefix(identifier, "x-coredata://") {
				sourceFolder, targetFolder, err := client.MoveNote(identifier, moveTargetFolder)
				if err != nil {
					fmt.Fprintf(os.Stderr, "✗ 錯誤：無法移動筆記 %s: %v\n", identifier, err)
					failCount++
					continue
				}
				fmt.Printf("✓ 已移動筆記從「%s」到「%s」\n", sourceFolder, targetFolder)
				successCount++
				continue
			}

			// 如果是名稱，檢查是否有多個同名筆記
			notes, err := client.FindNotesByName(identifier)
			if err != nil {
				fmt.Fprintf(os.Stderr, "✗ 錯誤：無法查找筆記「%s」: %v\n", identifier, err)
				failCount++
				continue
			}

			if len(notes) == 0 {
				fmt.Fprintf(os.Stderr, "✗ 錯誤：找不到筆記「%s」\n", identifier)
				failCount++
				continue
			}

			if len(notes) > 1 {
				// 有多個同名筆記，列出並要求使用 ID
				fmt.Fprintf(os.Stderr, "\n錯誤：找到 %d 個名為「%s」的筆記，請使用 ID 指定：\n\n", len(notes), identifier)
				
				w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "#\t標題\t資料夾\tID")
				fmt.Fprintln(w, "─\t────\t────\t──")
				for i, note := range notes {
					fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i+1, note.Name, note.Container, note.ID)
				}
				w.Flush()
				
				fmt.Fprintf(os.Stderr, "\n請使用以下指令：\n")
				fmt.Fprintf(os.Stderr, "  notes move \"<ID>\" -t \"%s\"\n\n", moveTargetFolder)
				
				failCount++
				hasAmbiguousNote = true
				continue
			}

			// 唯一筆記，直接移動
			note := notes[0]
			sourceFolder, targetFolder, err := client.MoveNote(note.ID, moveTargetFolder)
			if err != nil {
				fmt.Fprintf(os.Stderr, "✗ 錯誤：無法移動筆記「%s」: %v\n", identifier, err)
				failCount++
				continue
			}
			fmt.Printf("✓ 已移動筆記「%s」從「%s」到「%s」\n", note.Name, sourceFolder, targetFolder)
			successCount++
		}

		// Step 3: 顯示統計結果
		fmt.Println()
		if failCount == 0 {
			if successCount == 1 {
				fmt.Println("成功移動 1 個筆記")
			} else {
				fmt.Printf("成功移動 %d 個筆記\n", successCount)
			}
		} else {
			fmt.Printf("成功移動 %d 個筆記，失敗 %d 個\n", successCount, failCount)
			if hasAmbiguousNote {
				os.Exit(1)
			}
		}

		if failCount > 0 && !hasAmbiguousNote {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(&moveTargetFolder, "to", "t", "", "目標資料夾名稱 (必填)")
	moveCmd.MarkFlagRequired("to")
}
