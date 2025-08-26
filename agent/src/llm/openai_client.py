"""
OpenAI API 客户端
封装OpenAI API调用功能
"""

import os
from typing import Optional, Dict, Any
from openai import OpenAI
from loguru import logger


class OpenAIClient:
    """
    OpenAI API 客户端
    提供与OpenAI API的交互功能
    """
    
    def __init__(
        self,
        model_name: str = "gpt-4",
        api_key: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 1000
    ):
        """
        初始化OpenAI客户端
        
        Args:
            model_name: 使用的模型名称
            api_key: OpenAI API密钥
            temperature: 生成温度
            max_tokens: 最大token数
        """
        self.model_name = model_name
        self.temperature = temperature
        self.max_tokens = max_tokens
        
        # 获取API密钥
        self.api_key = api_key or os.getenv("OPENAI_API_KEY")
        if not self.api_key:
            raise ValueError("OpenAI API密钥未设置")
        
        # 初始化客户端
        self.client = OpenAI(api_key=self.api_key)
        
        logger.info(f"OpenAI client initialized with model: {model_name}")
    
    def generate(self, prompt: str, **kwargs) -> str:
        """
        生成文本回复
        
        Args:
            prompt: 输入提示
            **kwargs: 其他参数
            
        Returns:
            生成的文本
        """
        try:
            # 设置默认参数
            params = {
                "model": self.model_name,
                "messages": [{"role": "user", "content": prompt}],
                "temperature": kwargs.get("temperature", self.temperature),
                "max_tokens": kwargs.get("max_tokens", self.max_tokens)
            }
            
            # 调用API
            response = self.client.chat.completions.create(**params)
            
            # 提取回复内容
            content = response.choices[0].message.content
            
            logger.debug(f"Generated response: {content[:100]}...")
            return content
            
        except Exception as e:
            logger.error(f"OpenAI API call failed: {e}")
            raise
    
    def generate_with_messages(self, messages: list, **kwargs) -> str:
        """
        使用消息列表生成回复
        
        Args:
            messages: 消息列表
            **kwargs: 其他参数
            
        Returns:
            生成的文本
        """
        try:
            params = {
                "model": self.model_name,
                "messages": messages,
                "temperature": kwargs.get("temperature", self.temperature),
                "max_tokens": kwargs.get("max_tokens", self.max_tokens)
            }
            
            response = self.client.chat.completions.create(**params)
            content = response.choices[0].message.content
            
            logger.debug(f"Generated response with messages: {content[:100]}...")
            return content
            
        except Exception as e:
            logger.error(f"OpenAI API call with messages failed: {e}")
            raise
    
    def generate_stream(self, prompt: str, **kwargs):
        """
        流式生成文本
        
        Args:
            prompt: 输入提示
            **kwargs: 其他参数
            
        Yields:
            生成的文本片段
        """
        try:
            params = {
                "model": self.model_name,
                "messages": [{"role": "user", "content": prompt}],
                "temperature": kwargs.get("temperature", self.temperature),
                "max_tokens": kwargs.get("max_tokens", self.max_tokens),
                "stream": True
            }
            
            stream = self.client.chat.completions.create(**params)
            
            for chunk in stream:
                if chunk.choices[0].delta.content is not None:
                    yield chunk.choices[0].delta.content
                    
        except Exception as e:
            logger.error(f"OpenAI streaming failed: {e}")
            raise
    
    def get_embeddings(self, text: str) -> list:
        """
        获取文本嵌入向量
        
        Args:
            text: 输入文本
            
        Returns:
            嵌入向量
        """
        try:
            response = self.client.embeddings.create(
                model="text-embedding-ada-002",
                input=text
            )
            
            return response.data[0].embedding
            
        except Exception as e:
            logger.error(f"OpenAI embeddings failed: {e}")
            raise
    
    def get_models(self) -> list:
        """
        获取可用的模型列表
        
        Returns:
            模型列表
        """
        try:
            response = self.client.models.list()
            return [model.id for model in response.data]
            
        except Exception as e:
            logger.error(f"Failed to get models: {e}")
            return []
    
    def get_usage(self) -> Dict[str, Any]:
        """
        获取API使用情况
        
        Returns:
            使用情况信息
        """
        try:
            # 注意：这个功能需要OpenAI Plus账户
            response = self.client.usage.list()
            return response.dict()
            
        except Exception as e:
            logger.error(f"Failed to get usage: {e}")
            return {} 