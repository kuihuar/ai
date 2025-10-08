#!/bin/bash

# 多数据库一致性验证工具 - 快速启动脚本
# Cobra + Viper 优化版本

set -e

echo "🚀 多数据库一致性验证工具 - Cobra + Viper 优化版本"
echo "=================================================="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go 1.19+"
    exit 1
fi

echo "✅ Go环境检查通过: $(go version)"

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 检查可执行文件
if [ ! -f "./bin/validator-optimization" ]; then
    echo "📦 构建可执行文件..."
    make build
    echo "✅ 构建完成"
fi

# 显示帮助信息
echo ""
echo "📖 可用命令:"
echo "  ./bin/validator-optimization --help              # 显示帮助"
echo "  ./bin/validator-optimization init --help         # 配置文件创建帮助"
echo "  ./bin/validator-optimization validate --help     # 验证命令帮助"
echo ""

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo "📝 创建默认配置文件..."
    ./bin/validator-optimization init --format yaml --output configs/config.yaml
    echo "✅ 配置文件已创建: configs/config.yaml"
    echo "⚠️  请编辑配置文件设置正确的数据库连接信息"
    echo ""
fi

# 显示配置文件内容
if [ -f "configs/config.yaml" ]; then
    echo "📋 当前配置文件内容:"
    echo "----------------------------------------"
    head -20 configs/config.yaml
    echo "----------------------------------------"
    echo ""
fi

# 提供使用示例
echo "💡 使用示例:"
echo ""
echo "1. 编辑配置文件:"
echo "   vim configs/config.yaml"
echo ""
echo "2. 执行验证:"
echo "   ./bin/validator-optimization validate"
echo ""
echo "3. 试运行模式:"
echo "   ./bin/validator-optimization validate --dry-run"
echo ""
echo "4. 设置并发数:"
echo "   ./bin/validator-optimization validate --workers 5"
echo ""
echo "5. 详细输出:"
echo "   ./bin/validator-optimization validate --verbose"
echo ""

# 检查是否要立即执行验证
read -p "是否要立即执行验证? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🔍 开始执行验证..."
    ./bin/validator-optimization validate --verbose
else
    echo "👋 退出。请编辑配置文件后手动执行验证。"
fi
