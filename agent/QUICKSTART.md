# 🚀 快速开始指南

## 1. 环境准备

### 1.1 克隆项目
```bash
git clone <your-repo-url>
cd agent
```

### 1.2 安装依赖
```bash
pip install -r requirements.txt
```

### 1.3 配置环境变量
```bash
# 复制环境变量模板
cp env.example .env

# 编辑 .env 文件，填入你的 OpenAI API 密钥
# OPENAI_API_KEY=your_actual_api_key_here
```

## 2. 快速测试

### 2.1 使用启动脚本
```bash
python start.py
```
选择选项2进行测试（无需API密钥）

### 2.2 直接运行示例
```bash
# 基础对话测试（无需API密钥）
python examples/basic_agent.py

# 或者导入使用
python -c "
from examples.basic_agent import test_basic_conversation
test_basic_conversation()
"
```

## 3. 基本使用

### 3.1 创建简单的Agent
```python
from src.agent.core import Agent
from src.llm.openai_client import OpenAIClient

# 创建LLM客户端
llm_client = OpenAIClient(
    model_name="gpt-3.5-turbo",
    temperature=0.7
)

# 创建Agent
agent = Agent(
    name="助手",
    description="一个有用的AI助手",
    llm_client=llm_client
)

# 开始对话
response = agent.chat("你好，请介绍一下你自己")
print(response)
```

### 3.2 使用工具
```python
from src.tools.base import Tool
from src.agent.core import Agent

# 定义自定义工具
class CalculatorTool(Tool):
    def __init__(self):
        super().__init__(
            name="calculator",
            description="进行数学计算",
            parameters={
                "expression": {
                    "type": "string",
                    "description": "要计算的数学表达式"
                }
            }
        )
    
    def execute(self, expression: str) -> str:
        try:
            result = eval(expression)
            return f"计算结果: {result}"
        except Exception as e:
            return f"计算错误: {e}"

# 创建带工具的Agent
agent = Agent(
    name="计算助手",
    description="可以进行数学计算的AI助手",
    llm_client=llm_client,
    tools=[CalculatorTool()]
)

# 使用工具
response = agent.chat("请计算 2 + 3 * 4")
print(response)
```

## 4. 项目结构

```
agent/
├── src/                    # 源代码
│   ├── agent/             # Agent核心模块
│   ├── tools/             # 工具模块
│   ├── llm/               # LLM客户端
│   ├── knowledge/         # 知识库模块
│   └── utils/             # 工具函数
├── examples/              # 示例代码
├── tests/                 # 测试文件
├── data/                  # 数据存储
├── docs/                  # 文档
├── logs/                  # 日志文件
├── requirements.txt       # 依赖包
├── env.example           # 环境变量模板
├── start.py              # 启动脚本
└── README.md             # 项目说明
```

## 5. 常见问题

### 5.1 API密钥错误
- 确保在 `.env` 文件中正确设置了 `OPENAI_API_KEY`
- 检查API密钥是否有效且有足够的余额

### 5.2 依赖安装失败
```bash
# 升级pip
pip install --upgrade pip

# 重新安装依赖
pip install -r requirements.txt --force-reinstall
```

### 5.3 模块导入错误
```bash
# 确保在项目根目录运行
cd agent

# 或者设置PYTHONPATH
export PYTHONPATH="${PYTHONPATH}:$(pwd)"
```

## 6. 下一步

- 查看 `README.md` 了解详细功能
- 探索 `examples/` 目录中的更多示例
- 阅读源代码了解实现细节
- 根据需要扩展工具和功能

## 7. 获取帮助

- 查看项目文档
- 检查日志文件 `logs/agent.log`
- 提交Issue或Pull Request

---

🎉 **恭喜！你已经成功搭建了AI Agent项目的基础框架！** 