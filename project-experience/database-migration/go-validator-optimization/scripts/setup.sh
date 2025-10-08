#!/bin/bash

# 项目环境设置脚本
# Project Setup Script

set -e

echo "🔧 设置多数据库一致性验证工具环境"
echo "====================================="

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go 1.19+"
    exit 1
fi

echo "✅ Go环境检查通过: $(go version)"

# 安装依赖
echo "📦 安装Go依赖..."
go mod tidy
go mod download
echo "✅ 依赖安装完成"

# 构建项目
echo "🔨 构建项目..."
make build
echo "✅ 项目构建完成"

# 创建必要的目录
echo "📁 创建必要目录..."
mkdir -p configs examples output/{reports,logs,temp}
echo "✅ 目录创建完成"

# 创建默认配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo "📝 创建默认配置文件..."
    ./bin/validator-optimization init --format yaml --output configs/config.yaml
    echo "✅ 默认配置文件已创建"
fi

echo ""
echo "🎉 环境设置完成！"
echo ""
echo "📖 下一步:"
echo "1. 编辑配置文件: vim configs/config.yaml"
echo "2. 运行验证: ./bin/validator-optimization validate --dry-run"
echo "3. 查看帮助: ./bin/validator-optimization --help"
echo ""
