#!/bin/bash

# 测试脚本
# Test Script

set -e

echo "🧪 运行测试套件"
echo "==============="

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go 1.19+"
    exit 1
fi

# 运行单元测试
echo "📋 运行单元测试..."
go test -v ./...

# 运行基准测试
echo "📊 运行基准测试..."
go test -bench=. ./...

# 检查测试覆盖率
echo "📈 检查测试覆盖率..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo "✅ 覆盖率报告已生成: coverage.html"

# 运行竞态检测
echo "🔍 运行竞态检测..."
go test -race ./...

echo ""
echo "✅ 测试完成！"
echo ""
echo "📊 测试结果:"
echo "  - 单元测试: 通过"
echo "  - 基准测试: 完成"
echo "  - 覆盖率报告: coverage.html"
echo "  - 竞态检测: 通过"
echo ""
