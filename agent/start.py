#!/usr/bin/env python3
"""
AI Agent 项目启动脚本
提供快速启动和测试功能
"""

import os
import sys
from pathlib import Path

# 添加项目根目录到Python路径
project_root = Path(__file__).parent
sys.path.append(str(project_root))

from examples.basic_agent import main as basic_agent_main
from examples.basic_agent import test_basic_conversation


def show_menu():
    """显示主菜单"""
    print("🤖 AI Agent 项目启动器")
    print("=" * 40)
    print("1. 启动基础Agent (需要API密钥)")
    print("2. 运行测试模式 (无需API密钥)")
    print("3. 安装依赖")
    print("4. 查看项目状态")
    print("5. 退出")
    print("-" * 40)


def install_dependencies():
    """安装项目依赖"""
    print("📦 正在安装项目依赖...")
    
    try:
        import subprocess
        result = subprocess.run([
            sys.executable, "-m", "pip", "install", "-r", "requirements.txt"
        ], capture_output=True, text=True)
        
        if result.returncode == 0:
            print("✅ 依赖安装成功!")
        else:
            print("❌ 依赖安装失败:")
            print(result.stderr)
            
    except Exception as e:
        print(f"❌ 安装过程中出现错误: {e}")


def check_project_status():
    """检查项目状态"""
    print("🔍 检查项目状态...")
    
    # 检查必要文件
    required_files = [
        "requirements.txt",
        "src/agent/core.py",
        "src/tools/base.py",
        "src/llm/openai_client.py",
        "examples/basic_agent.py"
    ]
    
    missing_files = []
    for file_path in required_files:
        if not Path(file_path).exists():
            missing_files.append(file_path)
    
    if missing_files:
        print("❌ 缺少以下文件:")
        for file_path in missing_files:
            print(f"   - {file_path}")
    else:
        print("✅ 项目文件完整")
    
    # 检查环境变量
    api_key = os.getenv("OPENAI_API_KEY")
    if api_key:
        print("✅ OpenAI API密钥已设置")
    else:
        print("⚠️  OpenAI API密钥未设置 (复制 env.example 为 .env 并填入密钥)")
    
    # 检查Python包
    required_packages = [
        "openai",
        "langchain",
        "chromadb",
        "fastapi",
        "pydantic"
    ]
    
    missing_packages = []
    for package in required_packages:
        try:
            __import__(package)
        except ImportError:
            missing_packages.append(package)
    
    if missing_packages:
        print(f"❌ 缺少以下Python包: {', '.join(missing_packages)}")
        print("   运行选项3安装依赖")
    else:
        print("✅ 所有依赖包已安装")


def main():
    """主函数"""
    while True:
        show_menu()
        
        try:
            choice = input("请选择操作 (1-5): ").strip()
            
            if choice == "1":
                print("\n🚀 启动基础Agent...")
                basic_agent_main()
                
            elif choice == "2":
                print("\n🧪 运行测试模式...")
                test_basic_conversation()
                
            elif choice == "3":
                install_dependencies()
                
            elif choice == "4":
                check_project_status()
                
            elif choice == "5":
                print("👋 再见!")
                break
                
            else:
                print("❌ 无效选择，请输入1-5")
                
        except KeyboardInterrupt:
            print("\n\n👋 再见!")
            break
        except Exception as e:
            print(f"❌ 发生错误: {e}")
        
        input("\n按回车键继续...")


if __name__ == "__main__":
    main() 