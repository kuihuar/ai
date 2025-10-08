# 脚本使用指南

本文档详细介绍了项目中提供的各种脚本工具的使用方法。

## 📁 脚本目录结构

```
scripts/
├── setup.sh        # 环境设置脚本
├── dev.sh          # 开发环境脚本
├── test.sh         # 测试脚本
└── run_example.sh  # 快速启动脚本
```

## 🚀 快速开始

### 首次使用
```bash
# 1. 克隆项目
git clone <repository-url>
cd go-validator-optimization

# 2. 运行环境设置脚本
./scripts/setup.sh

# 3. 开始使用
./scripts/run_example.sh
```

## 🔧 脚本详细说明

### 1. `scripts/setup.sh` - 环境设置脚本

**用途**: 首次使用项目时的环境设置

**功能**:
- 检查Go环境
- 安装项目依赖
- 构建项目
- 创建必要目录
- 生成默认配置文件

**使用方法**:
```bash
./scripts/setup.sh
```

**输出示例**:
```
🔧 设置多数据库一致性验证工具环境
=====================================
✅ Go环境检查通过: go version go1.23.8 darwin/amd64
📦 安装Go依赖...
✅ 依赖安装完成
🔨 构建项目...
✅ 项目构建完成
📁 创建必要目录...
✅ 目录创建完成
📝 创建默认配置文件...
✅ 默认配置文件已创建

🎉 环境设置完成！

📖 下一步:
1. 编辑配置文件: vim configs/config.yaml
2. 运行验证: ./bin/validator-optimization validate --dry-run
3. 查看帮助: ./bin/validator-optimization --help
```

### 2. `scripts/dev.sh` - 开发环境脚本

**用途**: 开发环境快速启动

**功能**:
- 构建项目
- 运行测试
- 格式化代码
- 使用开发配置验证

**使用方法**:
```bash
./scripts/dev.sh
```

**输出示例**:
```
🚀 开发环境快速启动
===================
🔨 构建项目...
✅ 项目构建完成
🧪 运行测试...
✅ 测试通过
🎨 格式化代码...
✅ 代码格式化完成
🔍 使用开发配置运行验证...
✅ 开发环境启动完成！

💡 常用开发命令:
  make build        # 构建项目
  make test         # 运行测试
  make fmt          # 格式化代码
  make lint         # 代码检查
  make clean        # 清理构建产物
```

### 3. `scripts/test.sh` - 测试脚本

**用途**: 运行完整的测试套件

**功能**:
- 运行单元测试
- 运行基准测试
- 生成覆盖率报告
- 竞态检测

**使用方法**:
```bash
./scripts/test.sh
```

**输出示例**:
```
🧪 运行测试套件
===============
📋 运行单元测试...
✅ 单元测试通过
📊 运行基准测试...
✅ 基准测试完成
📈 检查测试覆盖率...
✅ 覆盖率报告已生成: coverage.html
🔍 运行竞态检测...
✅ 竞态检测通过

✅ 测试完成！

📊 测试结果:
  - 单元测试: 通过
  - 基准测试: 完成
  - 覆盖率报告: coverage.html
  - 竞态检测: 通过
```

### 4. `scripts/run_example.sh` - 快速启动脚本

**用途**: 快速启动和示例演示

**功能**:
- 检查环境
- 构建项目
- 创建配置文件
- 提供使用示例

**使用方法**:
```bash
./scripts/run_example.sh
```

**输出示例**:
```
🚀 多数据库一致性验证工具 - Cobra + Viper 优化版本
==================================================
✅ Go环境检查通过: go version go1.23.8 darwin/amd64

📖 可用命令:
  ./bin/validator-optimization --help              # 显示帮助
  ./bin/validator-optimization init --help         # 配置文件创建帮助
  ./bin/validator-optimization validate --help     # 验证命令帮助

💡 使用示例:

1. 编辑配置文件:
   vim configs/config.yaml

2. 执行验证:
   ./bin/validator-optimization validate

3. 试运行模式:
   ./bin/validator-optimization validate --dry-run

4. 设置并发数:
   ./bin/validator-optimization validate --workers 5

5. 详细输出:
   ./bin/validator-optimization validate --verbose

是否要立即执行验证? (y/N):
```

## 🛠️ Makefile 命令

除了脚本外，项目还提供了Makefile命令：

### 基本构建
```bash
make build              # 构建二进制文件
make clean              # 清理构建产物
make run                # 构建并运行
```

### 开发工具
```bash
make test               # 运行测试
make fmt                # 格式化代码
make lint               # 代码检查
make deps               # 安装依赖
```

### 跨平台构建
```bash
make build-all          # 构建所有平台版本
```

### 配置管理
```bash
make init-config        # 创建默认配置
make validate-dry       # 试运行验证
make validate-dev       # 使用开发配置验证
```

### 信息查看
```bash
make help               # 显示帮助
make version            # 显示版本信息
```

## 🔄 工作流程

### 开发工作流
```bash
# 1. 环境设置（首次）
./scripts/setup.sh

# 2. 日常开发
./scripts/dev.sh

# 3. 运行测试
./scripts/test.sh

# 4. 提交代码前
make fmt
make lint
make test
```

### 部署工作流
```bash
# 1. 构建所有平台版本
make build-all

# 2. 运行完整测试
./scripts/test.sh

# 3. 生产环境验证
./bin/validator-optimization validate --config configs/prod.yaml --dry-run
```

## 🐛 故障排除

### 脚本执行权限问题
```bash
# 添加执行权限
chmod +x scripts/*.sh
```

### Go环境问题
```bash
# 检查Go版本
go version

# 安装依赖
go mod tidy
```

### 构建问题
```bash
# 清理并重新构建
make clean
make build
```

## 📝 自定义脚本

你可以根据需要创建自定义脚本：

```bash
# 创建自定义脚本
touch scripts/my-script.sh
chmod +x scripts/my-script.sh

# 脚本模板
#!/bin/bash
set -e

echo "🚀 我的自定义脚本"
echo "================="

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 你的脚本逻辑
echo "✅ 脚本执行完成"
```

## 📞 支持

如有脚本使用问题，请：
1. 检查脚本执行权限
2. 确认Go环境正确
3. 查看错误日志
4. 提交Issue或联系维护者
