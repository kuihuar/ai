#!/bin/bash
# run_example.sh
# Go语言多数据库验证工具使用示例

echo "=== Go语言多数据库一致性验证工具 ==="
echo ""

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go环境，请先安装Go"
    exit 1
fi

echo "1. 初始化Go模块..."
go mod tidy

echo ""
echo "2. 创建默认配置文件..."
echo "   创建JSON格式配置文件..."
go run . init config.json

echo "   创建YAML格式配置文件..."
go run . init config.yaml

echo ""
echo "3. 显示帮助信息..."
go run . help

echo ""
echo "4. 配置文件已创建，请编辑config.json或config.yaml文件设置正确的数据库连接信息"
echo "   然后运行: go run ."
echo ""

echo "=== 使用说明 ==="
echo "1. 编辑config.json或config.yaml文件，设置Azure和AWS数据库连接信息"
echo "2. 运行验证: go run ."
echo "3. 查看结果: consistency_report.json"
echo ""
echo "支持的配置文件格式:"
echo "  - JSON: config.json"
echo "  - YAML: config.yaml 或 config.yml"
echo "4. 查看日志: 控制台输出和日志文件"
echo ""

echo "=== 项目结构 ==="
echo "go-validator/"
echo "├── go.mod              # Go模块文件"
echo "├── main.go             # 主程序入口"
echo "├── types.go            # 数据结构定义"
echo "├── validator.go        # 验证器核心逻辑"
echo "├── config.go           # 配置文件处理"
echo "├── config.json         # 配置文件（需要编辑）"
echo "├── config.json.example # 配置文件示例"
echo "└── README.md           # 说明文档"
echo ""

echo "=== 编译和运行 ==="
echo "# 直接运行"
echo "go run ."
echo ""
echo "# 编译后运行"
echo "go build -o validator"
echo "./validator"
echo ""

echo "=== 验证结果 ==="
echo "验证完成后会生成以下文件："
echo "- consistency_report.json: 详细的验证报告"
echo "- 控制台输出: 实时验证进度和摘要"
echo ""

echo "准备就绪！请编辑config.json文件后运行验证。"
