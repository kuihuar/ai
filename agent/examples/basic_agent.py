"""
基础Agent示例
演示如何使用Agent进行简单对话
"""

import os
import sys
from pathlib import Path

# 添加项目根目录到Python路径
project_root = Path(__file__).parent.parent
sys.path.append(str(project_root))

from src.agent.core import Agent
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()


def main():
    """主函数"""
    print("🤖 AI Agent 基础示例")
    print("=" * 50)
    
    # 获取API密钥
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        print("❌ 错误: 请设置 OPENAI_API_KEY 环境变量")
        print("在 .env 文件中添加: OPENAI_API_KEY=your_api_key_here")
        return
    
    # 创建Agent
    print("🔄 正在初始化Agent...")
    agent = Agent(
        model_name="gpt-3.5-turbo",  # 使用更便宜的模型进行测试
        api_key=api_key
    )
    
    print("✅ Agent初始化完成!")
    print("\n💬 开始对话 (输入 'quit' 退出)")
    print("-" * 50)
    
    # 对话循环
    while True:
        try:
            # 获取用户输入
            user_input = input("\n👤 你: ").strip()
            
            # 检查退出命令
            if user_input.lower() in ['quit', 'exit', '退出']:
                print("\n👋 再见!")
                break
            
            if not user_input:
                continue
            
            # 处理用户输入
            print("🤖 Agent: 正在思考...")
            response = agent.chat(user_input)
            print(f"🤖 Agent: {response}")
            
        except KeyboardInterrupt:
            print("\n\n👋 再见!")
            break
        except Exception as e:
            print(f"❌ 错误: {e}")


def test_basic_conversation():
    """测试基础对话功能"""
    print("🧪 测试基础对话功能")
    print("=" * 30)
    
    # 创建Agent（不连接API，用于测试）
    agent = Agent(model_name="gpt-3.5-turbo")
    
    # 测试对话
    test_messages = [
        "你好，请介绍一下你自己",
        "你能帮我做什么？",
        "今天天气怎么样？",
        "请给我讲个笑话"
    ]
    
    for message in test_messages:
        print(f"\n👤 用户: {message}")
        print("🤖 Agent: [模拟回复] 这是一个测试回复，实际使用时会调用真实的AI模型。")
    
    print("\n✅ 基础对话测试完成")


if __name__ == "__main__":
    # 检查是否有API密钥
    if os.getenv("OPENAI_API_KEY"):
        main()
    else:
        print("⚠️  未找到API密钥，运行测试模式")
        test_basic_conversation() 