# Apple Notes CLI 產品需求文件 (PRD)

## 1. 專案概述

### 1.1 目標
使用 Golang 建立一個命令列工具 (CLI)，透過 AppleScript 與 macOS Notes.app 互動，實現筆記管理功能。

### 1.2 目標用戶
- macOS 用戶
- 喜歡使用命令列工具管理筆記的開發者
- 需要自動化筆記操作的工作者

### 1.3 技術棧
- **語言**: Go (Golang)
- **CLI 框架**: Cobra
- **系統整合**: AppleScript via `osascript`
- **目標平台**: macOS

---

## 2. 功能需求

### 2.1 筆記操作

#### 2.1.1 建立筆記 (create)

**語法**：
```bash
notes create <title> [options]
```

**參數**：

| 參數 | 簡寫 | 說明 | 範例 |
|------|------|------|------|
| `title` | - | 筆記標題（必填） | `"會議記錄"` |
| `--body` | `-b` | 純文字內容（自動轉 HTML） | `-b "內容..."` |
| `--html` | - | HTML 格式內容 | `--html "<p>內容</p>"` |
| `--folder` | `-f` | 指定資料夾 | `-f "工作"` |
| `--file` | - | 從檔案讀取內容 | `--file ./note.html` |
| `--stdin` | - | 從標準輸入讀取內容 | `cat note.txt \| notes create "標題" --stdin` |

**使用範例**：
```bash
# 基本建立
$ notes create "會議記錄"

# 帶純文字內容（自動轉 HTML）
$ notes create "會議記錄" -b "這是會議內容"

# 帶 HTML 格式內容
$ notes create "會議記錄" --html "<h1>會議</h1><p>內容...</p>"

# 指定資料夾
$ notes create "會議記錄" -f "工作"

# 從檔案讀取內容
$ notes create "會議記錄" --file ./content.html

# 從 stdin 讀取內容
$ echo "會議內容" | notes create "會議記錄" --stdin
```

**AppleScript 對應**：
```applescript
tell application "Notes"
    make new note with properties {name:"標題", body:"<p>內容</p>"} at folder "資料夾"
end tell
```

**技術細節**：
- `body` 屬性必須是 HTML 格式，CLI 會將純文字自動轉換
- 唯讀屬性（`id`, `creation date`, `modification date`）由系統自動產生
- 若指定資料夾不存在，回傳錯誤訊息

---

#### 2.1.2 列出筆記 (list)

**語法**：
```bash
notes list [options]
```

**參數**：

| 參數 | 簡寫 | 說明 |
|------|------|------|
| `--folder` | `-f` | 篩選特定資料夾 |
| `--long` | `-l` | 顯示完整資訊（含完整日期、內容預覽） |
| `--preview` | `-p` | 顯示內容預覽 |
| `--json` | - | JSON 格式輸出 |

**顯示欄位**：

| 欄位 | 預設 | --long | 說明 |
|------|:----:|:------:|------|
| # | ✅ | ✅ | 列表序號 |
| Name | ✅ | ✅ | 筆記標題 |
| Folder | ✅ | ✅ | 所在資料夾 |
| Created | ✅ | ✅ | 建立日期 |
| Modified | ✅ | ✅ | 修改日期 |
| Attachments | ❌ | ✅ | 附件數量及名稱 |
| Preview | ❌ | ✅ | 內容預覽（前 50 字元） |

**輸出範例**：
```bash
# 預設輸出
$ notes list
#  Name              Folder        Created      Modified
1  會議記錄          工作          03/01        03/04
2  購物清單          個人          03/02        03/03

# 詳細輸出
$ notes list --long
#  Name              Folder        Created              Modified             Attachments    Preview
1  會議記錄          工作          2026/03/01 10:30     2026/03/04 15:20     📎 2           今天討論了專案進度...

# JSON 輸出
$ notes list --json
[
  {
    "id": "x-coredata://...",
    "name": "會議記錄",
    "folder": "工作",
    "created": "2026-03-01T10:30:00+08:00",
    "modified": "2026-03-04T15:20:00+08:00",
    "preview": "今天討論了..."
  }
]
```

---

#### 2.1.3 顯示筆記 (show)

**語法**：
```bash
notes show <id|name>
```

**識別方式**：
- 列表序號：`notes show 1`
- 完整名稱：`notes show "會議記錄"`
- 部分名稱：`notes show 會議`（模糊比對）

**輸出內容**：
```
標題：會議記錄
資料夾：工作
建立日期：2026/03/01 10:30
修改日期：2026/03/04 15:20

內容：
今天討論了專案進度...

附件 (2)：
  1. 報告.pdf (1.2 MB)
  2. 簽名.png (45 KB)
```

---

#### 2.1.4 搜尋筆記 (search)

**語法**：
```bash
notes search <keyword> [options]
```

**參數**：

| 參數 | 簡寫 | 說明 |
|------|------|------|
| `--title-only` | `-t` | 只搜尋標題 |
| `--body-only` | `-b` | 只搜尋內容 |
| `--folder` | `-f` | 限定資料夾 |
| `--since` | - | 搜尋某日期之後修改的筆記 |
| `--before` | - | 搜尋某日期之前修改的筆記 |

**日期格式**：
- 絕對日期：`YYYY-MM-DD`
- 相對日期：`today`, `yesterday`, `last-week`, `last-month`

**輸出範例**：
```bash
$ notes search "會議" --folder "工作" --since "2026-01-01"
Found 3 notes:

#  Name              Folder    Modified
1  會議記錄-0304      工作      2026/03/04
2  週會筆記          工作      2026/03/01
3  會議待辦事項      工作      2026/02/28
```

---

#### 2.1.5 刪除筆記 (delete)

**語法**：
```bash
notes delete <id|name> [options]
```

**參數**：

| 參數 | 簡寫 | 說明 |
|------|------|------|
| `--force` | `-f` | 略過確認提示 |
| `--permanent` | - | 永久刪除（不經垃圾桶） |

**預設行為**：
- 刪除後移至「最近刪除」資料夾
- 30 天後自動清除
- 可從 Notes.app 復原

**輸出範例**：
```bash
# 依名稱刪除
$ notes delete "會議記錄"
確定要刪除筆記「會議記會議記錄」嗎？: y
已將筆記「會議記錄」移到「最近刪除」

# 依 ID 刪除
$ notes delete "x-coredata://..."
已將筆記移到「最近刪除」: x-coredata://...


$ notes delete "會議記錄" --permanent
⚠️  警告：永久刪除無法復原！
確定要永久刪除筆記「會議記錄」嗎？: y
已永久刪除筆記「會議記錄」
```

---

#### 2.1.6 匯出筆記 (export)

**語法**：
```bash
notes export <id|name> [options]
```

**識別方式**：
- 筆記名稱：`notes export "會議記錄" --format md -o note.md`
- 筆記 ID：`notes export "x-coredata://..." --format html -o note.html`
- 當可能存在同名筆記時，應優先使用筆記 ID

**目前已實作參數**：

| 參數 | 簡寫 | 說明 |
|------|------|------|
| `--format` | `-f` | 輸出格式：`html` 或 `md` |
| `--output` | `-o` | 輸出檔案路徑 |

**規則**：
- 若未指定 `--format`，會依輸出副檔名推斷格式
- 若輸出到 stdout 且未指定格式，預設使用 `md`
- 目前尚未支援 `txt`、`json`、`--folder`、`--all`、附件匯出

**格式說明**：

| 格式 | 說明 | 附件處理 |
|------|------|----------|
| `html` | 保留 Apple Notes 回傳的 HTML 內容 | 目前未支援附件匯出 |
| `md` | 將常見 HTML 標籤轉為 Markdown | 目前未支援附件匯出 |

---

### 2.2 資料夾操作

#### 2.2.1 列出資料夾
```bash
notes folder list [--json]
```

#### 2.2.2 建立資料夾
```bash
notes folder create <name>
```

---

### 2.3 筆記識別機制

採用**混合識別**方式：

| 識別方式 | 範例 | 說明 |
|----------|------|------|
| 名稱 | `notes show "會議記錄"` | 完整比對筆記名稱 |
| 列表序號 | `notes show 1` | CLI 內部對應到實際筆記 |
| 部分名稱 | `notes show 會議` | 模糊比對（多個結果時提示用戶選擇） |

---

## 3. 附件處理

### 3.1 支援的附件類型

| 類型 | 格式 |
|------|------|
| 圖片 | PNG, JPG, GIF, HEIC |
| 文件 | PDF, DOC, XLS |
| 音訊 | M4A, MP3 |
| 影片 | MP4, MOV |

### 3.2 附件相關功能

| 功能 | 支援程度 | 說明 |
|------|:--------:|------|
| 列出附件 | ✅ 完整 | 可取得名稱、ID |
| 讀取附件 | ✅ 完整 | 透過 `save` 指令存檔 |
| 匯出附件 | ✅ 完整 | 匯出筆記時一併匯出 |
| 新增附件 | ⚠️ 有限 | 需用 UI Scripting（第二階段） |
| 刪除附件 | ⚠️ 有限 | 需用 UI Scripting（第二階段） |

### 3.3 AppleScript 附件操作

```applescript
-- 取得筆記的所有附件
attachments of note "會議記錄"

-- 取得附件屬性
name of attachment 1 of note "會議記錄"

-- 儲存附件到檔案
save attachment 1 of note "會議記錄" in file "/path/to/output.pdf"
```

---

## 4. 非功能需求

### 4.1 效能
- 大量筆記操作時顯示進度提示
- 快取筆記列表以加速序號對應

### 4.2 錯誤處理
- 友善的錯誤訊息
- 筆記不存在時的提示
- 權限不足時的引導
- 匯出時磁碟空間不足的警告

### 4.4 測試策略
- 單元測試必須可在不連線 Notes.app 的情況下執行
- AppleScript/Notes 互動改以整合測試驗證，避免一般 CI 或本機 `go test ./...` 誤失敗
- 整合測試僅在 macOS 且可使用 `osascript` 時執行
- 測試至少覆蓋字串跳脫、日期解析、HTML 轉換、建立/查詢/搜尋/移動/刪除筆記與資料夾列舉

### 4.3 使用者體驗
- 刪除前的確認提示
- 操作成功的回饋訊息
- 支援 `--help` 查看指令說明
- 永久刪除時的警告訊息

---

## 5. 專案結構

```
auto_notes/
├── cmd/
│   └── notes/
│       └── main.go           # CLI 入口點
├── internal/
│   ├── apple/
│   │   └── notes.go          # AppleScript 執行邏輯
│   └── cli/
│       ├── root.go           # 根指令
│       ├── create.go         # create 指令
│       ├── list.go           # list 指令
│       ├── show.go           # show 指令
│       ├── search.go         # search 指令
│       ├── delete.go         # delete 指令
│       ├── export.go         # export 指令
│       └── folder.go         # folder 相關指令
├── go.mod
├── go.sum
├── prd.md
└── README.md
```

---

## 6. 實作階段

### 階段一：基礎功能
- [x] 初始化 Go module 專案結構
- [x] 安裝 cobra 套件
- [x] 實作 AppleScript 執行模組
- [x] 實作 create, list, show, delete 指令

### 階段二：進階功能
- [x] 實作 search 指令
- [x] 實作 export 指令（`md` / `html`）
- [x] 實作 folder 相關指令
- [x] 實作 move 指令
- [ ] 支援 `stdin` 建立筆記
- [ ] 支援 JSON 輸出
- [ ] 補齊附件匯出能力

### 階段三：優化與文件
- [x] 測試並優化錯誤處理
- [x] 編寫 README 文件
- [x] 建立單元測試與整合測試分層
- [x] 補齊 Apple Notes 整合測試案例
- [ ] 補強 AppleScript 注入防護
- [ ] 補齊安裝發佈流程（`go install` / Homebrew）

### 6.1 目前測試案例
- [x] `escapeAppleScriptString()` 特殊字元跳脫
- [x] `parseAppleDate()` Apple 日期解析
- [x] `textToHTML()` 純文字轉 HTML
- [x] `NewNotesClient()` 建立 client
- [x] `ListFolders()` 列出資料夾
- [x] `CreateNote()` 建立筆記
- [x] `DeleteNote()` 刪除筆記
- [x] `ShowNote()` 顯示筆記
- [x] `ExportNote()` 依 ID 匯出筆記
- [x] `SearchNotes()` 搜尋筆記
- [x] `FindNotesByName()` 依名稱查找筆記
- [x] `MoveNote()` 移動筆記

---

## 7. AppleScript 技術參考

### 7.1 筆記屬性

| 屬性 | 類型 | 可讀 | 可寫 | 說明 |
|------|------|:----:|:----:|------|
| `id` | text | ✅ | ❌ | 系統唯一 ID |
| `name` | text | ✅ | ✅ | 筆記標題 |
| `body` | text | ✅ | ✅ | 筆記內容（HTML 格式） |
| `creation date` | date | ✅ | ❌ | 建立日期 |
| `modification date` | date | ✅ | ❌ | 修改日期 |
| `container` | folder | ✅ | ✅ | 所在資料夾 |

### 7.2 資料夾屬性

| 屬性 | 類型 | 可讀 | 可寫 | 說明 |
|------|------|:----:|:----:|------|
| `id` | text | ✅ | ❌ | 系統唯一 ID |
| `name` | text | ✅ | ✅ | 資料夾名稱 |
