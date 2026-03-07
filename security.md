# 🔒 資安風險分析報告 - auto_notes 專案

**分析日期**: 2026-03-04  
**分析工具**: 程式碼審查

---

## 🚨 高風險問題

### 1. AppleScript 注入漏洞 (Critical)

**嚴重程度**: 🔴 高

**位置**: `internal/apple/notes.go`

**問題描述**: 
`escapeAppleScriptString()` 函數的跳脫機制不足，用戶輸入直接嵌入 AppleScript 字串中。

```go
func escapeAppleScriptString(s string) string {
    s = strings.ReplaceAll(s, "\\", "\\\\")
    s = strings.ReplaceAll(s, "\"", "\\\"")
    s = strings.ReplaceAll(s, "'", "\\'")
    return s
}
```

**風險**: 
攻擊者可以透過特製輸入注入惡意 AppleScript 指令，例如：
- 使用 `\"` 組合繞過跳脫
- 注入換行符號 (`\n`) 在某些情況下可能執行額外指令
- 利用 Unicode 字元編碼繞過過濾

**攻擊範例**:
```bash
notes create -t "test" -b $'hello\"\ndo shell script \"rm -rf ~\"\n--'
```

**影響範圍**:
- `CreateNote()` - title, body, folder 參數
- `ShowNote()` - identifier 參數
- `DeleteNote()` - identifier 參數
- `SearchNotes()` - keyword, folder 參數
- `CreateFolder()` - name 參數

---

### 2. 命令執行風險

**嚴重程度**: 🔴 高

**位置**: `internal/apple/notes.go` 的 `runAppleScript()`

```go
func runAppleScript(script string) (string, error) {
    cmd := exec.Command("osascript", "-e", script)
    // ...
}
```

**風險**: 
雖然使用 `-e` 參數，但 script 內容包含用戶輸入，若跳脫不完全，可能執行任意系統命令。AppleScript 的 `do shell script` 指令可以執行任意 shell 命令。

**潛在影響**:
- 資料外洩
- 系統被入侵
- 惡意程式執行

---

## ⚠️ 中風險問題

### 3. 敏感資料權限問題

**嚴重程度**: 🟡 中

**位置**: `cmd/export.go`

```go
err := ioutil.WriteFile(exportOutput, []byte(content), 0644)
```

**風險**: 
- 匯出檔案使用 `0644` 權限，所有使用者皆可讀取
- 筆記可能包含密碼、API keys、個人敏感資訊
- 應考慮使用更嚴格的權限 (如 `0600`)

---

### 4. 輸入驗證不足

**嚴重程度**: 🟡 中

**問題描述**: 
多個功能缺乏嚴格的輸入驗證：
- 資料夾名稱未驗證長度和特殊字元
- 筆記標題未限制長度
- 搜尋關鍵字未過濾特殊字元

**風險**:
- 可能導致阻斷服務 (DoS)
- 可能觸發未預期的程式行為

---

### 5. 錯誤訊息洩漏

**嚴重程度**: 🟡 中

**位置**: 多個 `cmd/*.go` 檔案

```go
fmt.Fprintf(os.Stderr, "錯誤：無法讀取檔案 %s: %v\n", createFile, err)
```

**風險**: 
錯誤訊息可能洩漏系統路徑、檔案結構等敏感資訊，方便攻擊者進一步攻擊。

---

## 📋 低風險問題

### 6. 使用已棄用的 API

**嚴重程度**: 🟢 低

**位置**: `cmd/create.go`, `cmd/export.go`

```go
content, err := ioutil.ReadFile(createFile)  // 已棄用
```

**建議**: 應使用 `os.ReadFile()` 和 `os.WriteFile()` (Go 1.16+)

---

### 7. 缺乏日誌和稽核

**嚴重程度**: 🟢 低

**問題**: 沒有任何操作日誌記錄，難以追蹤異常行為或進行取證分析。

---

## 📊 風險評估摘要

| 風險等級 | 問題數量 | 說明 |
|---------|---------|------|
| 🔴 高 | 2 | AppleScript 注入、命令執行 |
| 🟡 中 | 3 | 權限、輸入驗證、錯誤洩漏 |
| 🟢 低 | 2 | 已棄用 API、缺乏日誌 |

---

## 💡 修復建議

### 1. AppleScript 注入修復

**優先級**: 高

**建議方案**:
- 不要使用字串拼接構建 AppleScript
- 改用參數化方式傳遞資料
- 或使用更嚴格的白名單輸入驗證

**範例修復**:
```go
// 使用白名單驗證
func validateInput(s string) error {
    // 只允許安全字元
    validPattern := regexp.MustCompile(`^[a-zA-Z0-9\u4e00-\u9fff\s\-_,.!?]+$`)
    if !validPattern.MatchString(s) {
        return errors.New("輸入包含不允許的字元")
    }
    return nil
}
```

### 2. 命令執行防護

**優先級**: 高

**建議**:
- 限制 AppleScript 可執行的指令
- 考慮使用沙箱機制
- 實作權限分級

### 3. 檔案權限修正

**優先級**: 中

```go
// 修改前
err := ioutil.WriteFile(exportOutput, []byte(content), 0644)

// 修改後
err := os.WriteFile(exportOutput, []byte(content), 0600)
```

### 4. 輸入驗證加強

**優先級**: 中

- 對所有用戶輸入進行白名單驗證
- 限制輸入長度
- 過濾特殊字元

### 5. 錯誤處理改善

**優先級**: 中

- 使用通用錯誤訊息對使用者顯示
- 詳細錯誤僅記錄於日誌

### 6. 安全測試納入自動化

**優先級**: 中

- 單元測試持續覆蓋 `escapeAppleScriptString()` 的特殊字元處理
- 後續應新增惡意 payload 測試案例，包含換行、跳脫字元與長字串輸入
- Apple Notes 整合測試應與一般單元測試分離，避免環境問題掩蓋真正的安全缺陷
- 建議固定使用：

```bash
GOCACHE=$(pwd)/.gocache go test ./...
GOCACHE=$(pwd)/.gocache go test -tags=integration ./internal/apple
```

---

## 🔐 安全開發建議

1. **防禦深度**: 不要依賴單一防禦機制
2. **最小權限原則**: 只給予必要的權限
3. **輸入驗證**: 永遠不信任用戶輸入
4. **安全編碼規範**: 遵循 OWASP 安全編碼規範
5. **定期審查**: 定期進行安全程式碼審查

---

## 📚 參考資料

- [OWASP Top 10](https://owasp.org/Top10/)
- [OWASP Command Injection](https://owasp.org/www-community/attacks/Command_Injection)
- [AppleScript Security Considerations](https://developer.apple.com/library/archive/documentation/AppleScript/Conceptual/AppleScriptLangGuide/reference/ASLR_keywords.html)
