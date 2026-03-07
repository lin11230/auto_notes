# Apple Notes CLI 開發狀態

## 專案概述
使用 Golang 打造 CLI 工具，透過 AppleScript 管理 macOS 本機的 Notes 應用程式。

## 專案結構
```
auto_notes/
├── main.go                 # 程式進入點
├── go.mod                  # Go module 定義
├── go.sum                  # 依賴版本鎖定
├── prd.md                  # 產品需求文件
├── .gitignore              # Git 忽略設定
├── cmd/                    # CLI 指令目錄
│   ├── root.go            # 根指令
│   ├── create.go          # 建立筆記指令
│   ├── list.go            # 列出筆記指令
│   ├── show.go            # 顯示筆記指令
│   ├── search.go          # 搜尋筆記指令
│   ├── delete.go          # 刪除筆記指令
│   ├── move.go            # 移動筆記指令
│   ├── export.go          # 匯出筆記指令
│   └── folder.go          # 資料夾管理指令
└── internal/
    └── apple/
        ├── notes.go                 # AppleScript 執行模組
        ├── notes_test.go            # 單元測試
        └── notes_integration_test.go # Apple Notes 整合測試
```

## 已完成項目

### 基礎建設
- [x] 初始化 Git repository 和 .gitignore
- [x] 初始化 Go module (github.com/kclin/auto_notes)
- [x] 安裝 spf13/cobra CLI 框架
- [x] 建立專案目錄結構

### AppleScript 模組 (internal/apple/notes.go)
- [x] `runAppleScript()` - 執行 AppleScript 的核心函數
- [x] `ListNotes()` - 列出筆記
- [x] `CreateNote()` - 建立新筆記
- [x] `ShowNote()` - 顯示筆記內容
- [x] `DeleteNote()` - 刪除筆記（支援永久刪除）
- [x] `SearchNotes()` - 搜尋筆記
- [x] `ListFolders()` - 列出資料夾
- [x] `CreateFolder()` - 建立資料夾
- [x] `ExportNote()` - 匯出筆記
- [x] `FindNotesByName()` - 依名稱查找所有同名筆記
- [x] `MoveNote()` - 移動筆記到指定資料夾
- [x] `parseAppleDate()` - 解析 Apple 日期格式
- [x] `escapeAppleScriptString()` - 跳脫特殊字元
- [x] `textToHTML()` - 純文字轉 HTML

### CLI 指令 (cmd/)
- [x] root.go - 根指令與說明
- [x] create.go - `notes create -t "標題" -b "內容"`
- [x] list.go - `notes list` / `notes ls`
- [x] show.go - `notes show <名稱或ID>`
- [x] search.go - `notes search <關鍵字>`
- [x] delete.go - `notes delete <名稱或ID>`
- [x] move.go - `notes move <名稱或ID> -t <資料夾>`
- [x] export.go - `notes export <名稱或ID> --format <md|html> -o file`
- [x] folder.go - `notes folder list` / `notes folder create`

### 編譯
- [x] 成功編譯產出 `notes` 執行檔
- [x] `./notes --help` 正確顯示說明

### 測試與驗證 (已完成)
- [x] 測試 `./notes list` 指令 ✓ (功能正常，可列出 1170 則筆記)
- [x] 測試 `./notes folder list` 指令 ✓
- [x] 測試 `./notes create` 建立筆記 ✓
- [x] 測試 `./notes show` 顯示筆記 ✓
- [x] 測試 `./notes search` 搜尋筆記 ✓
- [x] 測試 `./notes delete` 刪除筆記 ✓ (支援依名稱或 **ID** 刪除，並移到「最近刪除」)
- [x] 測試 `./notes move` 移動筆記 ✓ (支援單個/批次移動，同名筆記檢測)
- [x] 測試 `./notes export` 匯出筆記 ✓
- [x] 新增 `internal/apple/notes_test.go` 單元測試
- [x] 新增 `internal/apple/notes_integration_test.go` 整合測試
- [x] 將整合測試改為 `integration && darwin` build tag，避免在一般 `go test` 中誤跑
- [x] 整合測試加入 `osascript` 與 Notes 可用性檢查

### 測試案例覆蓋
- [x] `escapeAppleScriptString()` 特殊字元跳脫
- [x] `parseAppleDate()` 日期字串解析
- [x] `textToHTML()` 文字轉 HTML
- [x] `NewNotesClient()` client 建立
- [x] `ListFolders()` 整合驗證
- [x] `CreateNote()` / `DeleteNote()` 整合驗證
- [x] `ShowNote()` 整合驗證
- [x] `ExportNote()` 依筆記 ID 匯出整合驗證
- [x] `SearchNotes()` 整合驗證
- [x] `FindNotesByName()` 整合驗證
- [x] `MoveNote()` 整合驗證
- [x] `export` 格式判定與 HTML/Markdown 轉換

### 測試執行方式
```bash
# 單元測試
GOCACHE=$(pwd)/.gocache go test ./...

# Apple Notes 整合測試
GOCACHE=$(pwd)/.gocache go test -tags=integration ./internal/apple
```

### 已解決問題
1. ~~AppleScript `container` 屬性取得失敗~~ - 已改用 `folder` 屬性並加入 try/error 處理
2. ~~`move to folder "Recently Deleted"` 失敗~~ - 已改用 `delete` 命令，直接移到垃圾桶

### 已完成功能
- [x] 編寫 README.md 使用說明文件
- [x] 新增版本資訊 (notes --version)
- [x] 新增 Makefile
- [x] 新增筆記移動功能 (notes move)

### 待開發功能 (未來版本)
- [ ] 新增筆記編輯功能 (notes edit)
- [ ] 新增筆記複製功能 (notes copy)
- [ ] 支援從 stdin 讀取內容
- [ ] 新增 JSON 輸出格式選項
- [ ] 支援 `go install` 安裝
- [ ] 新增 Homebrew formula (選用)

## 技術筆記

### AppleScript 日期格式
AppleScript 回傳的日期格式為：`Wednesday, March 4, 2026 at 12:00:00 PM`

### Apple Notes 筆記 ID 格式
筆記 ID 格式為：`x-coredata://...`

### 特殊字元處理
AppleScript 字串需要跳脫：`\`, `"`, `'`

## 編譯與執行
```bash
# 編譯
go build -o notes .

# 執行
./notes --help
./notes list
./notes create -t "測試筆記" -b "這是內容"
./notes export "x-coredata://..." --format html -o note.html
./notes export "測試筆記" --format md -o note.md
```

## 下一個 Agent 接續事項
1. 先執行 `GOCACHE=$(pwd)/.gocache go test ./...` 確認單元測試
2. 在 macOS Notes 可用環境執行 `GOCACHE=$(pwd)/.gocache go test -tags=integration ./internal/apple`
3. 如整合測試失敗，優先檢查 AppleScript 與 Notes 權限/可用性
4. 考慮補強 AppleScript 注入防護與更多進階功能
