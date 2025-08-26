"""
AI Agent 核心模块
实现基础的Agent功能和架构
"""

import json
from typing import List, Dict, Any, Optional
from datetime import datetime
from loguru import logger

from ..llm.openai_client import OpenAIClient
from ..agent.memory import MemoryManager
from ..agent.planner import TaskPlanner
from ..agent.executor import TaskExecutor
from ..tools.base import Tool


class Agent:
    """
    AI Agent 核心类
    负责协调各个组件，处理用户输入并生成回复
    """
    
    def __init__(
        self,
        model_name: str = "gpt-4",
        api_key: Optional[str] = None,
        memory_path: str = "./data/memory",
        max_memory: int = 1000
    ):
        """
        初始化Agent
        
        Args:
            model_name: 使用的语言模型名称
            api_key: OpenAI API密钥
            memory_path: 记忆存储路径
            max_memory: 最大记忆数量
        """
        self.model_name = model_name
        self.llm = OpenAIClient(model_name=model_name, api_key=api_key)
        self.memory = MemoryManager(memory_path, max_memory)
        self.planner = TaskPlanner(self.llm)
        self.executor = TaskExecutor(self)
        
        # 工具列表
        self.tools: List[Tool] = []
        
        # 对话历史
        self.conversation_history: List[Dict[str, str]] = []
        
        # Agent状态
        self.status = "ready"
        
        logger.info(f"Agent initialized with model: {model_name}")
    
    def add_tool(self, tool: Tool) -> None:
        """
        添加工具到Agent
        
        Args:
            tool: 要添加的工具
        """
        self.tools.append(tool)
        logger.info(f"Added tool: {tool.name}")
    
    def get_tool_descriptions(self) -> str:
        """
        获取所有工具的描述信息
        
        Returns:
            工具描述字符串
        """
        descriptions = []
        for tool in self.tools:
            descriptions.append(f"- {tool.name}: {tool.description}")
        return "\n".join(descriptions)
    
    def chat(self, message: str) -> str:
        """
        处理用户消息并生成回复
        
        Args:
            message: 用户输入的消息
            
        Returns:
            Agent的回复
        """
        try:
            self.status = "processing"
            logger.info(f"Processing message: {message[:50]}...")
            
            # 更新对话历史
            self.conversation_history.append({
                "role": "user",
                "content": message,
                "timestamp": datetime.now().isoformat()
            })
            
            # 分析用户意图
            intent = self._analyze_intent(message)
            
            # 根据意图选择处理方式
            if intent.get("requires_tool"):
                response = self._handle_tool_request(message, intent)
            elif intent.get("is_task"):
                response = self._handle_task(message, intent)
            else:
                response = self._handle_conversation(message)
            
            # 更新对话历史
            self.conversation_history.append({
                "role": "assistant",
                "content": response,
                "timestamp": datetime.now().isoformat()
            })
            
            # 保存到记忆
            self.memory.add_memory(
                f"User: {message}\nAssistant: {response}",
                "conversation"
            )
            
            self.status = "ready"
            return response
            
        except Exception as e:
            logger.error(f"Error processing message: {e}")
            self.status = "error"
            return f"抱歉，处理您的请求时出现了错误: {str(e)}"
    
    def _analyze_intent(self, message: str) -> Dict[str, Any]:
        """
        分析用户意图
        
        Args:
            message: 用户消息
            
        Returns:
            意图分析结果
        """
        prompt = f"""
        分析以下用户消息的意图，返回JSON格式的结果：
        
        消息: {message}
        
        可用工具:
        {self.get_tool_descriptions()}
        
        请分析：
        1. 是否需要使用工具
        2. 是否是复杂任务
        3. 用户的具体需求
        
        返回格式：
        {{
            "requires_tool": true/false,
            "is_task": true/false,
            "tool_name": "工具名称或null",
            "task_type": "任务类型或null",
            "intent": "用户意图描述"
        }}
        """
        
        try:
            response = self.llm.generate(prompt)
            return json.loads(response)
        except:
            # 如果解析失败，返回默认分析
            return {
                "requires_tool": False,
                "is_task": False,
                "tool_name": None,
                "task_type": None,
                "intent": "general_conversation"
            }
    
    def _handle_conversation(self, message: str) -> str:
        """
        处理一般对话
        
        Args:
            message: 用户消息
            
        Returns:
            回复内容
        """
        # 获取相关记忆
        relevant_memories = self.memory.retrieve_relevant(message, k=3)
        
        # 构建上下文
        context = self._build_context(relevant_memories)
        
        prompt = f"""
        {context}
        
        基于以上上下文，请回复用户的消息。回复要自然、有帮助且符合上下文。
        
        用户消息: {message}
        
        回复:
        """
        
        return self.llm.generate(prompt)
    
    def _handle_tool_request(self, message: str, intent: Dict[str, Any]) -> str:
        """
        处理工具请求
        
        Args:
            message: 用户消息
            intent: 意图分析结果
            
        Returns:
            工具执行结果
        """
        tool_name = intent.get("tool_name")
        
        # 查找对应的工具
        tool = None
        for t in self.tools:
            if t.name == tool_name:
                tool = t
                break
        
        if not tool:
            return "抱歉，我没有找到合适的工具来处理您的请求。"
        
        try:
            # 提取工具参数
            params = self._extract_tool_params(message, tool)
            
            # 执行工具
            result = tool.execute(**params)
            
            return f"工具执行结果:\n{result}"
            
        except Exception as e:
            logger.error(f"Tool execution error: {e}")
            return f"工具执行失败: {str(e)}"
    
    def _handle_task(self, message: str, intent: Dict[str, Any]) -> str:
        """
        处理复杂任务
        
        Args:
            message: 用户消息
            intent: 意图分析结果
            
        Returns:
            任务执行结果
        """
        return self.executor.execute_task(message)
    
    def _build_context(self, relevant_memories: List[str]) -> str:
        """
        构建对话上下文
        
        Args:
            relevant_memories: 相关记忆
            
        Returns:
            上下文字符串
        """
        context_parts = []
        
        # 添加相关记忆
        if relevant_memories:
            context_parts.append("相关记忆:")
            for memory in relevant_memories:
                context_parts.append(f"- {memory}")
        
        # 添加最近的对话历史
        recent_history = self.conversation_history[-6:]  # 最近3轮对话
        if recent_history:
            context_parts.append("\n最近的对话:")
            for entry in recent_history:
                role = "用户" if entry["role"] == "user" else "助手"
                context_parts.append(f"{role}: {entry['content']}")
        
        return "\n".join(context_parts)
    
    def _extract_tool_params(self, message: str, tool: Tool) -> Dict[str, Any]:
        """
        从用户消息中提取工具参数
        
        Args:
            message: 用户消息
            tool: 要使用的工具
            
        Returns:
            工具参数字典
        """
        # 这里可以根据具体工具类型实现参数提取逻辑
        # 简化实现，返回空字典
        return {}
    
    def get_status(self) -> Dict[str, Any]:
        """
        获取Agent状态信息
        
        Returns:
            状态信息字典
        """
        return {
            "status": self.status,
            "model": self.model_name,
            "tools_count": len(self.tools),
            "memory_size": len(self.conversation_history),
            "available_tools": [tool.name for tool in self.tools]
        }
    
    def reset(self) -> None:
        """
        重置Agent状态
        """
        self.conversation_history = []
        self.status = "ready"
        logger.info("Agent reset completed") 