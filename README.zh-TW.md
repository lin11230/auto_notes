# Apple Notes CLI

一個使用 Go 撰寫的命令列工具，透過 AppleScript 管理 macOS 的 Notes.app。

## 安裝

### 從原始碼建置

```bash
git clone https://github.com/kclin/auto_notes.git
cd auto_notes
go build -o notes .
sudo mv notes /usr/local/bin/
```

## 使用方式

### 顯示說明

```bash
./notes --help
./notes <command> --help
```

### 列出筆記

```bash
./notes list
./notes ls
./notes list -f "工作"
```

### 建立筆記

```bash
./notes create -t "筆記標題" -b "筆記內容"
./notes create -t "筆記標題" -b "筆記內容" -f "工作"
```

### 顯示筆記

```bash
./notes show "會議記錄"
./notes show "x-coredata://..."
```

### 搜尋筆記

```bash
./notes search "關鍵字"
./notes search "關鍵字" -f "工作"
```

### 刪除筆記

```bash
./notes delete "會議記錄"
./notes delete "x-coredata://..."
```

### 移動筆記

```bash
./notes move "會議記錄" -t "封存"
./notes move "筆記1" "筆記2" "筆記3" -t "工作"
./notes move "x-coredata://..." -t "個人"
```

### 匯出筆記

```bash
./notes export "x-coredata://..." --format md -o output.md
./notes export "會議記錄" --format md -o output.md
./notes export "會議記錄" --format html -o output.html
./notes export "會議記錄"
```

### 資料夾管理

```bash
./notes folder list
./notes folder create "新資料夾"
```

## 系統需求

- macOS 與 Apple Notes app
- Go 1.21+（建置原始碼時需要）

## 測試

```bash
GOCACHE=$(pwd)/.gocache go test ./...
GOCACHE=$(pwd)/.gocache go test -tags=integration ./internal/apple
```

## 授權

MIT
