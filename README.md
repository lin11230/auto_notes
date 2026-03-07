# Apple Notes CLI

一個使用 Golang 打造的 CLI 工具，透過 AppleScript 管理 macOS 本機的 Notes 應用程式。

## 安裝

### 從原始碼建置

```bash
# 複製儲存庫
git clone https://github.com/kclin/auto_notes.git
cd auto_notes

# 編譯
go build -o notes .

# 移動到 PATH（選用）
sudo mv notes /usr/local/bin/
```

## 使用說明

### 顯示說明

```bash
./notes --help
./notes <command> --help
```

### 列出筆記

```bash
# 列出所有筆記
./notes list
./notes ls

# 列出特定資料夾的筆記
./notes list -f "資料夾名稱"
```

### 建立筆記

```bash
# 建立新筆記
./notes create -t "筆記標題" -b "筆記內容"

# 在特定資料夾建立筆記
./notes create -t "筆記標題" -b "筆記內容" -f "資料夾名稱"
```

### 顯示筆記內容

```bash
# 依名稱顯示筆記
./notes show "筆記名稱"

# 依 ID 顯示筆記
./notes show "x-coredata://..."
```

### 搜尋筆記

```bash
# 搜尋包含關鍵字的筆記
./notes search "關鍵字"

# 在特定資料夾搜尋
./notes search "關鍵字" -f "資料夾名稱"
```

### 刪除筆記

```bash
# 依名稱刪除筆記（移到「最近刪除」）
./notes delete "筆記名稱"

# 依 ID 刪除筆記
./notes delete "x-coredata://..."
```

### 移動筆記

```bash
# 移動單個筆記到指定資料夾
./notes move "筆記名稱" -t "目標資料夾"

# 批次移動多個筆記
./notes move "筆記1" "筆記2" "筆記3" -t "工作"

# 使用 ID 移動（避免同名衝突）
./notes move "x-coredata://..." -t "個人"
```

### 匯出筆記

```bash
# 匯出到檔案
./notes export "筆記名稱" -o output.txt

# 匯出到 stdout
./notes export "筆記名稱"
```

### 資料夾管理

```bash
# 列出所有資料夾
./notes folder list

# 建立新資料夾
./notes folder create "新資料夾名稱"
```

## 指令一覽

| 指令 | 說明 |
|------|------|
| `notes list` | 列出筆記 |
| `notes create` | 建立筆記 |
| `notes show` | 顯示筆記內容 |
| `notes search` | 搜尋筆記 |
| `notes delete` | 刪除筆記 |
| `notes move` | 移動筆記到指定資料夾 |
| `notes export` | 匯出筆記 |
| `notes folder list` | 列出資料夾 |
| `notes folder create` | 建立資料夾 |

## 系統需求

- macOS（需要 Apple Notes 應用程式）
- Go 1.21+（僅建置時需要）

## 測試

專案目前將測試分為兩層：

- 單元測試：不依賴 Apple Notes，可在一般 `go test` 流程執行
- 整合測試：需要 macOS、`osascript` 與可存取的 Notes.app

```bash
# 執行全部單元測試
GOCACHE=$(pwd)/.gocache go test ./...

# 執行 Apple Notes 整合測試
GOCACHE=$(pwd)/.gocache go test -tags=integration ./internal/apple
```

目前已涵蓋的測試案例：

- `escapeAppleScriptString()` 特殊字元跳脫
- `parseAppleDate()` Apple 日期格式解析
- `textToHTML()` 純文字轉 HTML
- `NewNotesClient()` client 建立
- `ListFolders()` 整合測試
- `CreateNote()` / `DeleteNote()` 整合測試
- `ShowNote()` 整合測試
- `SearchNotes()` 整合測試
- `FindNotesByName()` 整合測試
- `MoveNote()` 整合測試

## 技術細節

本工具使用 AppleScript 與 macOS Notes 應用程式通訊，支援：

- 筆記的 CRUD 操作
- 筆記移動功能
- 資料夾管理
- 搜尋功能
- 匯出功能

## 授權

MIT License
