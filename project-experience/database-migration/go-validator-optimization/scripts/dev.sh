#!/bin/bash

# 开发环境快速启动脚本
# Development Environment Quick Start Script

set -e

echo "🚀 开发环境快速启动"
echo "==================="

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go 1.19+"
    exit 1
fi

# 构建项目
echo "🔨 构建项目..."
make build

# 运行测试
echo "🧪 运行测试..."
make test

# 代码格式化
echo "🎨 格式化代码..."
make fmt

# 使用开发配置运行验证
echo "🔍 使用开发配置运行验证..."
./bin/validator-optimization validate --config configs/dev.yaml --dry-run --verbose

echo ""
echo "✅ 开发环境启动完成！"
echo ""
echo "💡 常用开发命令:"
echo "  make build        # 构建项目"
echo "  make test         # 运行测试"
echo "  make fmt          # 格式化代码"
echo "  make lint         # 代码检查"
echo "  make clean        # 清理构建产物"
echo ""
