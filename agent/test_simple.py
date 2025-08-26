#!/usr/bin/env python3
"""
简单的项目测试脚本
不依赖外部包，用于验证项目结构
"""

import os
import sys
from pathlib import Path

def test_project_structure():
    """测试项目结构"""
    print("🔍 检查项目结构...")
    
    # 检查必要目录
    required_dirs = [
        "src",
        "src/agent", 
        "src/tools",
        "src/llm",
        "src/knowledge",
        "src/utils",
        "examples",
        "tests",
        "data",
        "docs",
        "logs"
    ]
    
    missing_dirs = []
    for dir_path in required_dirs:
        if not Path(dir_path).exists():
            missing_dirs.append(dir_path)
    
    if missing_dirs:
        print("❌ 缺少以下目录:")
        for dir_path in missing_dirs:
            print(f"   - {dir_path}")
        return False
    else:
        print("✅ 项目目录结构完整")
        return True

def test_core_files():
    """测试核心文件"""
    print("\n📁 检查核心文件...")
    
    required_files = [
        "src/agent/core.py",
        "src/tools/base.py", 
        "src/llm/openai_client.py",
        "examples/basic_agent.py",
        "requirements.txt",
        "README.md",
        "env.example",
        "start.py"
    ]
    
    missing_files = []
    for file_path in required_files:
        if not Path(file_path).exists():
            missing_files.append(file_path)
    
    if missing_files:
        print("❌ 缺少以下文件:")
        for file_path in missing_files:
            print(f"   - {file_path}")
        return False
    else:
        print("✅ 核心文件完整")
        return True

def test_python_imports():
    """测试Python模块导入（不依赖外部包）"""
    print("\n🐍 检查Python模块...")
    
    # 检查__init__.py文件
    init_files = [
        "src/__init__.py",
        "src/agent/__init__.py",
        "src/tools/__init__.py", 
        "src/llm/__init__.py",
        "src/knowledge/__init__.py",
        "src/utils/__init__.py",
        "tests/__init__.py"
    ]
    
    missing_inits = []
    for init_file in init_files:
        if not Path(init_file).exists():
            missing_inits.append(init_file)
    
    if missing_inits:
        print("❌ 缺少以下__init__.py文件:")
        for init_file in missing_inits:
            print(f"   - {init_file}")
        return False
    else:
        print("✅ Python模块结构正确")
        return True

def test_syntax():
    """测试Python语法（不执行）"""
    print("\n🔤 检查Python语法...")
    
    python_files = [
        "src/agent/core.py",
        "src/tools/base.py",
        "src/llm/openai_client.py", 
        "examples/basic_agent.py",
        "start.py"
    ]
    
    syntax_errors = []
    for file_path in python_files:
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            compile(content, file_path, 'exec')
        except SyntaxError as e:
            syntax_errors.append(f"{file_path}: {e}")
        except Exception as e:
            syntax_errors.append(f"{file_path}: {e}")
    
    if syntax_errors:
        print("❌ 发现语法错误:")
        for error in syntax_errors:
            print(f"   - {error}")
        return False
    else:
        print("✅ Python语法检查通过")
        return True

def show_next_steps():
    """显示下一步操作"""
    print("\n🎯 下一步操作:")
    print("1. 安装依赖: pip install -r requirements.txt")
    print("2. 配置环境变量: cp env.example .env")
    print("3. 编辑 .env 文件，填入 OpenAI API 密钥")
    print("4. 运行测试: python start.py")
    print("5. 开始使用: python examples/basic_agent.py")

def main():
    """主函数"""
    print("🤖 AI Agent 项目测试")
    print("=" * 40)
    
    # 运行所有测试
    tests = [
        test_project_structure,
        test_core_files,
        test_python_imports,
        test_syntax
    ]
    
    all_passed = True
    for test in tests:
        if not test():
            all_passed = False
    
    print("\n" + "=" * 40)
    if all_passed:
        print("🎉 所有测试通过！项目结构完整。")
        show_next_steps()
    else:
        print("❌ 部分测试失败，请检查项目结构。")

if __name__ == "__main__":
    main() 