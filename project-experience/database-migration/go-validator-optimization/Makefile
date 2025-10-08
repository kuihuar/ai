# Makefile for multi-database-validator-optimization

# 项目信息
PROJECT_NAME := multi-database-validator-optimization
BINARY_NAME := validator-optimization
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.goVersion=$(GO_VERSION)"

# 目录
BIN_DIR := bin
DIST_DIR := dist

# 默认目标
.PHONY: all
all: build

# 构建二进制文件
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) .
	@echo "✅ Build completed: $(BIN_DIR)/$(BINARY_NAME)"

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
	@echo "✅ Clean completed"

# 运行程序
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BIN_DIR)/$(BINARY_NAME)

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# 跨平台构建
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ Cross-platform build completed"

# 创建配置文件
.PHONY: init-config
init-config:
	@echo "Creating default configuration..."
	./$(BIN_DIR)/$(BINARY_NAME) init --format yaml --output configs/config.yaml

# 运行验证（试运行模式）
.PHONY: validate-dry
validate-dry: build
	@echo "Running validation in dry-run mode..."
	./$(BIN_DIR)/$(BINARY_NAME) validate --dry-run

# 运行验证（使用开发配置）
.PHONY: validate-dev
validate-dev: build
	@echo "Running validation with dev config..."
	./$(BIN_DIR)/$(BINARY_NAME) validate --config configs/dev.yaml --dry-run

# 显示帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the program"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  deps         - Install dependencies"
	@echo "  build-all    - Build for all platforms"
	@echo "  init-config  - Create default configuration"
	@echo "  validate-dry - Run validation in dry-run mode"
	@echo "  validate-dev - Run validation with dev config"
	@echo "  help         - Show this help"

# 显示版本信息
.PHONY: version
version:
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"
