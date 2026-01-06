# Cursor IDE 配置指南

## Cursor 配置文件位置

### 1. 项目级别配置

#### `.cursorrules` 文件（推荐）

在项目根目录创建 `.cursorrules` 文件，用于定义项目特定的 AI 规则和上下文。这个文件会被 Git 跟踪，团队成员可以共享。

**位置：** `项目根目录/.cursorrules`

**示例内容：**

```
# 项目上下文
这是一个基于 Kratos 框架的 Go 微服务项目。

# 架构原则
- 遵循 Clean Architecture 分层架构
- 使用依赖注入（Wire）
- 数据层负责外部依赖（数据库、Redis、Kafka、第三方服务等）
- 业务层不依赖数据层实现细节

# 代码规范
- 遵循 Go 官方代码规范
- 使用 Kratos 框架的最佳实践
- 所有公开函数必须有文档注释
- 错误处理使用 Kratos 的错误定义

# 项目结构
- api/: API 定义（proto 文件）
- internal/biz/: 业务逻辑层
- internal/data/: 数据访问层
- internal/service/: 服务层
- cmd/: 应用入口

# 第三方服务接口定义
- gRPC 服务：api/external/{service}/v1/*.proto
- HTTP REST API：internal/data/external/{service}/types.go
```

#### `.vscode/settings.json`（工作区设置）

Cursor 基于 VS Code，可以使用 `.vscode/settings.json` 配置工作区设置。

**位置：** `项目根目录/.vscode/settings.json`

**示例内容：**

```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "[go]": {
    "editor.defaultFormatter": "golangci-lint",
    "editor.formatOnSave": true
  },
  "files.exclude": {
    "**/.git": true,
    "**/.DS_Store": true,
    "**/bin": true
  }
}
```

### 2. 用户级别配置

Cursor 的用户配置存储在用户目录下，不同操作系统位置不同：

- **macOS:** `~/Library/Application Support/Cursor/User/settings.json`
- **Windows:** `%APPDATA%\Cursor\User\settings.json`
- **Linux:** `~/.config/Cursor/User/settings.json`

这些配置是用户个人的，不会被 Git 跟踪。

### 3. Cursor 特定配置

Cursor 还有一些特定的配置目录：

- **扩展配置：** `~/.cursor/extensions/`
- **AI 模型缓存：** `~/.cursor/models/`（如果使用本地模型）
- **工作区状态：** `~/.cursor/workspaceStorage/`

## 在其他电脑上继续操作

### 方法 1：使用 Git 同步（推荐）

项目级别的 Cursor 配置可以通过 Git 同步：

1. **提交配置文件到 Git：**
   ```bash
   # 创建 .cursorrules 文件（如果还没有）
   touch .cursorrules
   
   # 创建 .vscode 目录和设置文件
   mkdir -p .vscode
   # 编辑 .vscode/settings.json
   
   # 提交到 Git
   git add .cursorrules .vscode/
   git commit -m "Add Cursor IDE configuration"
   git push
   ```

2. **在其他电脑上克隆项目：**
   ```bash
   git clone <repository-url>
   cd sre
   ```

3. **打开项目：**
   ```bash
   cursor .
   ```

   Cursor 会自动读取项目中的 `.cursorrules` 和 `.vscode/settings.json` 配置。

### 方法 2：手动同步配置

如果需要同步用户级别的配置：

1. **导出配置：**
   - 复制 `~/.cursor/User/settings.json`
   - 记录已安装的扩展列表

2. **在新电脑上：**
   - 安装 Cursor
   - 复制 `settings.json` 到对应位置
   - 安装相同的扩展

### 方法 3：使用 Cursor 账户同步（如果支持）

如果 Cursor 支持账户同步功能（类似 VS Code 的 Settings Sync），可以：
1. 登录 Cursor 账户
2. 启用设置同步
3. 在其他电脑上登录相同账户

## 推荐的项目配置

### `.cursorrules` 文件

在项目根目录创建 `.cursorrules` 文件，定义项目特定的 AI 上下文：

```markdown
# SRE 项目 - Cursor AI 规则

## 项目概述
这是一个基于 Kratos 框架的 Go 微服务研究项目，专注于软件工程最佳实践。

## 架构原则
- Clean Architecture 分层架构
- 依赖注入（Wire）
- 数据层管理所有外部依赖
- 业务层不依赖数据层实现

## 代码组织
- api/: Protobuf API 定义
- internal/biz/: 业务逻辑（不依赖外部）
- internal/data/: 数据访问层（数据库、Redis、Kafka、第三方服务）
- internal/service/: 服务层（gRPC/HTTP 处理）
- cmd/: 应用入口

## 第三方服务接口定义
- gRPC: api/external/{service}/v1/*.proto
- HTTP REST: internal/data/external/{service}/types.go

## 代码规范
- 遵循 Go 官方规范
- 所有公开函数必须有文档注释
- 使用 Kratos 错误定义
- 参考 docs/code-standards/ 目录下的规范文档

## 文档优先
- 重要决策和最佳实践都记录在 docs/ 目录
- 编写代码前先查看相关文档
```

### `.vscode/settings.json`

```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnSave": false,
  "go.useLanguageServer": true,
  
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit"
  },
  
  "[go]": {
    "editor.defaultFormatter": "golangci-lint",
    "editor.formatOnSave": true,
    "editor.snippetSuggestions": "none"
  },
  
  "[go.mod]": {
    "editor.defaultFormatter": "golang.go"
  },
  
  "[proto]": {
    "editor.defaultFormatter": "bufbuild.buf"
  },
  
  "files.exclude": {
    "**/.git": true,
    "**/.DS_Store": true,
    "**/bin": true,
    "**/*.pb.go": false
  },
  
  "search.exclude": {
    "**/vendor": true,
    "**/node_modules": true,
    "**/*.pb.go": true
  },
  
  "files.watcherExclude": {
    "**/.git/objects/**": true,
    "**/.git/subtree-cache/**": true,
    "**/node_modules/**": true,
    "**/vendor/**": true
  }
}
```

### `.vscode/extensions.json`（推荐扩展）

```json
{
  "recommendations": [
    "golang.go",
    "bufbuild.buf",
    "golangci.golangci-lint",
    "ms-vscode.vscode-json"
  ]
}
```

## 更新 .gitignore

确保 `.gitignore` 包含以下内容（如果需要排除某些 Cursor 文件）：

```gitignore
# Cursor IDE
.cursor/
.cursor-workspace

# 但保留项目配置
!.cursorrules
!.vscode/
```

## 团队协作

### 共享配置

以下文件应该提交到 Git，以便团队共享：
- ✅ `.cursorrules` - 项目 AI 规则
- ✅ `.vscode/settings.json` - 工作区设置
- ✅ `.vscode/extensions.json` - 推荐扩展

### 个人配置

以下配置不应该提交到 Git（个人偏好）：
- ❌ `~/.cursor/User/settings.json` - 用户个人设置
- ❌ `~/.cursor/User/keybindings.json` - 个人快捷键
- ❌ 扩展配置（除非是团队必需的）

## 快速设置脚本

创建 `scripts/setup-cursor.sh` 用于快速设置：

```bash
#!/bin/bash

# 创建 .cursorrules 文件
if [ ! -f .cursorrules ]; then
  cat > .cursorrules << 'EOF'
# SRE 项目 - Cursor AI 规则
# 项目概述、架构原则、代码规范等
EOF
  echo "Created .cursorrules"
fi

# 创建 .vscode 目录
mkdir -p .vscode

# 创建 settings.json
if [ ! -f .vscode/settings.json ]; then
  cat > .vscode/settings.json << 'EOF'
{
  "go.formatTool": "goimports",
  "editor.formatOnSave": true
}
EOF
  echo "Created .vscode/settings.json"
fi

echo "Cursor configuration setup complete!"
```

## 总结

### 配置文件位置

| 配置类型 | 位置 | 是否提交 Git |
|---------|------|-------------|
| 项目 AI 规则 | `.cursorrules` | ✅ 是 |
| 工作区设置 | `.vscode/settings.json` | ✅ 是 |
| 推荐扩展 | `.vscode/extensions.json` | ✅ 是 |
| 用户设置 | `~/.cursor/User/settings.json` | ❌ 否 |

### 在其他电脑上继续操作

1. **克隆项目：** `git clone <repo-url>`
2. **打开项目：** `cursor .`
3. **安装推荐扩展：** Cursor 会提示安装 `.vscode/extensions.json` 中推荐的扩展
4. **配置完成：** `.cursorrules` 和 `.vscode/settings.json` 会自动生效

通过这种方式，团队成员可以在任何电脑上获得一致的开发体验。

