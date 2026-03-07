# auto_notes 資安檢視報告

**檢視日期**: 2026-03-07  
**檢視方式**: 原始碼審查  
**檢視範圍**: `cmd/*.go`, `internal/apple/notes.go`

## 風險摘要

本專案目前最大的資安風險仍集中在 AppleScript 指令組裝與本機檔案處理。整體來說，這不是遠端服務型系統，攻擊面偏向本機使用者輸入、惡意筆記內容、以及導出的敏感資料暴露；但因為程式會直接呼叫 `osascript` 控制 Notes.app，一旦輸入處理失當，風險等級仍然偏高。

本次審查的結論：

- 高風險：1 項
- 中風險：3 項
- 低風險：2 項

## 主要發現

### 1. 高風險：`ListNotes(folder)` 仍存在未跳脫的 AppleScript 注入點

**位置**: [`internal/apple/notes.go`](/Users/kclin/Documents/auto_notes/internal/apple/notes.go#L51)

`ListNotes(folder)` 在 `folder != ""` 的分支中，直接把使用者提供的 `folder` 插入 AppleScript：

```go
script = fmt.Sprintf(`
    tell application "Notes"
        ...
        repeat with eachNote in notes of folder "%s"
```

這裡和同檔案其他函式不同，沒有經過 `escapeAppleScriptString()`。  
代表只要 `folder` 來源可控，就可能破壞 AppleScript 字串結構，最終插入額外指令。由於 AppleScript 可透過 `do shell script` 執行 shell 指令，這個問題應視為高風險。

**改善建議**:

1. 立即將 `ListNotes(folder)` 改為使用 `escapeAppleScriptString(folder)`。
2. 不要只依賴跳脫；應補上輸入白名單或長度限制，例如限制資料夾名稱長度與控制字元。
3. 中長期應避免以字串拼接方式建 AppleScript，改成固定模板加嚴格驗證過的參數。
4. 新增惡意 payload 測試，例如包含 `"`, `\`, 換行與 AppleScript 關鍵字的 folder 名稱。

---

### 2. 中風險：匯出檔案權限過寬，敏感筆記可能被其他本機使用者讀取

**位置**: [`cmd/export.go`](/Users/kclin/Documents/auto_notes/cmd/export.go#L47)

目前匯出使用：

```go
os.WriteFile(exportOutput, []byte(content), 0644)
```

`0644` 代表同機其他使用者可讀。對筆記內容這類可能包含帳密、API keys、個資、商業資訊的資料來說，權限過寬。

**改善建議**:

1. 將匯出檔案權限改為 `0600`。
2. 若輸出檔已存在，考慮提示是否覆蓋，避免敏感內容被不預期地寫到共享位置。
3. 文件中應明確提醒使用者匯出筆記可能含敏感資訊。

---

### 3. 中風險：HTML 筆記內容可被原樣匯出，存在本機內容注入與主動內容風險

**位置**: [`cmd/export.go`](/Users/kclin/Documents/auto_notes/cmd/export.go#L85), [`internal/apple/notes.go`](/Users/kclin/Documents/auto_notes/internal/apple/notes.go#L231)

`CreateNote()` 目前不會對內容做 HTML 安全處理，只是用 `textToHTML()` 包裝；而 `export --format html` 會直接把 Notes 回傳的 body 原樣寫出。

若筆記內容本身含有：

- `<script>`
- `<meta http-equiv=refresh>`
- 外部資源載入
- 惡意連結或釣魚內容

則使用者一旦在瀏覽器開啟匯出的 HTML，可能觸發主動內容。

這比較像「本機內容安全」與「儲存型內容注入」風險，而不是典型 Web XSS，但仍值得管控。

**改善建議**:

1. 預設匯出格式維持 Markdown 或純文字較安全。
2. HTML 匯出可加入 `--unsafe-html` 明確確認，或在輸出前加警告。
3. 若產品定位允許，可考慮在 HTML 匯出時移除 `<script>`、`iframe`、`object`、`embed>`、危險事件屬性。
4. 若未來支援從檔案匯入 HTML，應明確區分「純文字建立」與「信任 HTML 建立」模式。

---

### 4. 中風險：錯誤訊息直接回傳底層 `osascript` 與檔案路徑資訊，可能洩漏系統細節

**位置**: [`internal/apple/notes.go`](/Users/kclin/Documents/auto_notes/internal/apple/notes.go#L42), [`cmd/create.go`](/Users/kclin/Documents/auto_notes/cmd/create.go#L36), [`cmd/export.go`](/Users/kclin/Documents/auto_notes/cmd/export.go#L49)

目前錯誤訊息會直接包含：

- AppleScript 原始錯誤輸出
- 本機檔案路徑
- 可能的系統環境資訊

這在本機 CLI 場景不算最嚴重，但如果未來這個工具被包進自動化流程、CI log、共享終端或聊天機器人代理，就可能造成資訊外洩。

**改善建議**:

1. 對使用者顯示摘要錯誤，例如「無法執行 Notes 操作」。
2. 詳細錯誤改由 `--debug` 開啟時才顯示。
3. 匯出或讀檔失敗訊息避免直接暴露完整敏感路徑。

---

### 5. 低風險：輸入缺乏長度與字元限制，容易造成異常行為或可預期外的腳本失敗

**位置**: [`cmd/create.go`](/Users/kclin/Documents/auto_notes/cmd/create.go#L28), [`cmd/export.go`](/Users/kclin/Documents/auto_notes/cmd/export.go#L26), [`internal/apple/notes.go`](/Users/kclin/Documents/auto_notes/internal/apple/notes.go#L35)

目前大多數參數都直接進入 AppleScript 或檔案輸出流程，缺乏：

- 最大長度限制
- 控制字元限制
- 路徑合理性檢查

這會增加腳本解析錯誤、記憶體使用膨脹、以及例外輸出情況。

**改善建議**:

1. 對 title、folder、keyword、identifier 設定最大長度。
2. 明確拒絕 ASCII 控制字元與不必要的不可見字元。
3. 對輸出路徑做基本檢查，例如禁止空白檔名或目錄覆寫。

---

### 6. 低風險：仍使用 `ioutil.ReadFile`，不影響安全本質，但表示輸入處理路徑尚未整理

**位置**: [`cmd/create.go`](/Users/kclin/Documents/auto_notes/cmd/create.go#L36)

這不是直接漏洞，但通常代表檔案 I/O 路徑還未進行完整整理與標準化，連帶也常與權限、大小限制、例外處理不一致一起出現。

**改善建議**:

1. 改為 `os.ReadFile`。
2. 讀檔前先檢查檔案大小上限。
3. 若未來支援 `stdin`，需統一檔案/標準輸入的大小限制與驗證流程。

## 建議修正優先順序

### 第一優先

1. 修補 `ListNotes(folder)` 的未跳脫插值。
2. 將匯出檔案權限調整為 `0600`。
3. 為 AppleScript 相關輸入補上安全測試。

### 第二優先

1. 建立統一的輸入驗證函式，處理 title/folder/identifier/keyword。
2. 將詳細錯誤輸出收斂到 `--debug` 模式。
3. 補上 HTML 匯出風險提示或安全模式。

### 第三優先

1. 逐步淘汰 `fmt.Sprintf` 拼接 AppleScript 的寫法。
2. 整理檔案 I/O API 與大小限制。
3. 規劃發佈前的安全檢查清單。

## 建議新增的安全測試

### 單元測試

- `escapeAppleScriptString()` 加入換行、tab、控制字元、長字串案例
- 輸出格式與路徑檢查的異常案例
- HTML 匯出內容中含 `<script>`、`iframe`、`onerror=` 的案例

### 整合測試

- 以惡意 folder 名稱驗證 `ListNotes()` 不會破壞 AppleScript
- 以同名筆記與 ID 查找驗證不會誤操作錯誤目標
- 匯出後檔案權限檢查

## 發佈前最低安全基線

在下一次版本發佈前，至少應完成下列項目：

- 修補 `ListNotes()` 的 AppleScript 注入點
- 將匯出檔案權限改為 `0600`
- 對 AppleScript 參數建立一致的驗證規則
- 建立最少一組注入測試與一組檔案權限測試

## 結論

這個專案目前最值得優先處理的不是網路攻擊面，而是本機腳本注入與敏感資料保護。  
若只做一件事，應先修掉 [`internal/apple/notes.go`](/Users/kclin/Documents/auto_notes/internal/apple/notes.go#L51) 的未跳脫 `folder` 插值；若做第二件事，應把 [`cmd/export.go`](/Users/kclin/Documents/auto_notes/cmd/export.go#L48) 的輸出權限改為 `0600`。
