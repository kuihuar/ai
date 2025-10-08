#!/bin/bash

# Python版本多数据库一致性验证工具运行示例

echo "=== Python版本多数据库一致性验证工具 ==="
echo

# 检查Python是否安装
if ! command -v python3 &> /dev/null; then
    echo "错误: 未找到python3，请先安装Python 3"
    exit 1
fi

# 检查是否安装了pymysql
if ! python3 -c "import pymysql" 2>/dev/null; then
    echo "警告: 未安装pymysql，正在安装..."
    if [ -f "requirements.txt" ]; then
        pip3 install -r requirements.txt
    else
        pip3 install pymysql
    fi
fi

echo "1. 创建默认配置文件..."
python3 multi_database_validator.py --init

echo
echo "2. 显示帮助信息..."
python3 multi_database_validator.py --help

echo
echo "3. 使用默认配置运行验证（会失败，因为连接信息是示例）..."
echo "注意: 实际使用时请修改config.json中的连接信息"
python3 multi_database_validator.py --config config.json

echo
echo "4. 查看生成的报告文件..."
if [ -f "consistency_report.json" ]; then
    echo "报告文件已生成: consistency_report.json"
    echo "报告内容预览:"
    head -20 consistency_report.json
else
    echo "未找到报告文件"
fi

echo
echo "=== 使用说明 ==="
echo "1. 修改config.json中的数据库连接信息"
echo "2. 运行: python3 multi_database_validator.py"
echo "3. 查看生成的consistency_report.json报告"
echo
echo "=== 完成 ==="
