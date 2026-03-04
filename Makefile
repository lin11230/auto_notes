# Apple Notes CLI Makefile

# 版本號
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# 建置相關變數
BINARY_NAME = notes
GO = go
GOFLAGS = -v

# 預設目標
.PHONY: all
all: build

# 建置執行檔
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) .

# 安裝到系統
.PHONY: install
install: build
	@cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✓ 已安裝 $(BINARY_NAME) 到 /usr/local/bin/"

# 清理建置產物
.PHONY: clean
clean:
	@rm -f $(BINARY_NAME)
	@echo "✓ 已清理建置產物"

# 執行測試
.PHONY: test
test:
	$(GO) test -v ./...

# 格式化程式碼
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# 檢查程式碼
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "請先安裝 golangci-lint" && exit 1)
	golangci-lint run ./...

# 顯示說明
.PHONY: help
help:
	@echo "Apple Notes CLI - 可用指令："
	@echo ""
	@echo "  make build      - 建置執行檔"
	@echo "  make install    - 安裝到 /usr/local/bin/"
	@echo "  make clean      - 清理建置產物"
	@echo "  make test       - 執行測試"
	@echo "  make fmt        - 格式化程式碼"
	@echo "  make lint       - 檢查程式碼品質"
	@echo ""
	@echo "建置時可指定版本："
	@echo "  make build VERSION=v1.0.0"