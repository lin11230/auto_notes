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
# 刪除筆記（移到「最近刪除」）
./notes delete "筆記名稱"
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
| `notes export` | 匯出筆記 |
| `notes folder list` | 列出資料夾 |
| `notes folder create` | 建立資料夾 |

## 系統需求

- macOS（需要 Apple Notes 應用程式）
- Go 1.21+（僅建置時需要）

## 技術細節

本工具使用 AppleScript 與 macOS Notes 應用程式通訊，支援：

- 筆記的 CRUD 操作
- 資料夾管理
- 搜尋功能
- 匯出功能

## 授權

MIT License