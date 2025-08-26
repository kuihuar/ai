"""
工具基类模块
定义所有工具的基础接口和通用功能
"""

from abc import ABC, abstractmethod
from typing import Any, Dict, Optional
from pydantic import BaseModel
from loguru import logger


class ToolSchema(BaseModel):
    """工具参数模式定义"""
    name: str
    description: str
    parameters: Dict[str, Any]


class Tool(ABC):
    """
    工具基类
    所有工具都应该继承此类并实现必要的方法
    """
    
    def __init__(self, name: str, description: str):
        """
        初始化工具
        
        Args:
            name: 工具名称
            description: 工具描述
        """
        self.name = name
        self.description = description
        self.schema = self._get_schema()
        
        logger.info(f"Initialized tool: {name}")
    
    @abstractmethod
    def execute(self, **kwargs) -> str:
        """
        执行工具功能
        
        Args:
            **kwargs: 工具参数
            
        Returns:
            执行结果字符串
        """
        pass
    
    def _get_schema(self) -> ToolSchema:
        """
        获取工具的参数模式
        
        Returns:
            工具模式定义
        """
        return ToolSchema(
            name=self.name,
            description=self.description,
            parameters=self._get_parameters_schema()
        )
    
    def _get_parameters_schema(self) -> Dict[str, Any]:
        """
        获取参数模式定义
        
        Returns:
            参数模式字典
        """
        return {}
    
    def validate_parameters(self, **kwargs) -> bool:
        """
        验证输入参数
        
        Args:
            **kwargs: 输入参数
            
        Returns:
            验证是否通过
        """
        try:
            # 这里可以实现具体的参数验证逻辑
            return True
        except Exception as e:
            logger.error(f"Parameter validation failed: {e}")
            return False
    
    def get_help(self) -> str:
        """
        获取工具使用帮助
        
        Returns:
            帮助信息
        """
        return f"""
工具名称: {self.name}
描述: {self.description}
参数: {self._get_parameters_schema()}
        """.strip()


class WebSearchTool(Tool):
    """
    网络搜索工具
    提供网络搜索功能
    """
    
    def __init__(self):
        super().__init__(
            name="web_search",
            description="搜索网络获取最新信息"
        )
    
    def execute(self, query: str) -> str:
        """
        执行网络搜索
        
        Args:
            query: 搜索查询
            
        Returns:
            搜索结果
        """
        try:
            # 这里应该实现实际的网络搜索逻辑
            # 可以使用Google Search API或其他搜索服务
            logger.info(f"Searching for: {query}")
            
            # 模拟搜索结果
            return f"搜索结果: 关于'{query}'的信息..."
            
        except Exception as e:
            logger.error(f"Web search failed: {e}")
            return f"搜索失败: {str(e)}"
    
    def _get_parameters_schema(self) -> Dict[str, Any]:
        return {
            "query": {
                "type": "string",
                "description": "搜索查询",
                "required": True
            }
        }


class CalculatorTool(Tool):
    """
    计算器工具
    提供数学计算功能
    """
    
    def __init__(self):
        super().__init__(
            name="calculator",
            description="执行数学计算"
        )
    
    def execute(self, expression: str) -> str:
        """
        执行数学计算
        
        Args:
            expression: 数学表达式
            
        Returns:
            计算结果
        """
        try:
            # 安全地执行数学表达式
            # 这里应该实现安全的表达式解析和计算
            logger.info(f"Calculating: {expression}")
            
            # 简单的实现，实际应该使用更安全的计算库
            result = eval(expression)
            return f"计算结果: {expression} = {result}"
            
        except Exception as e:
            logger.error(f"Calculation failed: {e}")
            return f"计算失败: {str(e)}"
    
    def _get_parameters_schema(self) -> Dict[str, Any]:
        return {
            "expression": {
                "type": "string",
                "description": "数学表达式",
                "required": True
            }
        }


class WeatherTool(Tool):
    """
    天气查询工具
    提供天气信息查询功能
    """
    
    def __init__(self):
        super().__init__(
            name="weather",
            description="查询指定城市的天气信息"
        )
    
    def execute(self, city: str) -> str:
        """
        查询天气信息
        
        Args:
            city: 城市名称
            
        Returns:
            天气信息
        """
        try:
            # 这里应该实现实际的天气API调用
            logger.info(f"Querying weather for: {city}")
            
            # 模拟天气数据
            return f"{city}的天气: 晴天，温度25°C，湿度60%"
            
        except Exception as e:
            logger.error(f"Weather query failed: {e}")
            return f"天气查询失败: {str(e)}"
    
    def _get_parameters_schema(self) -> Dict[str, Any]:
        return {
            "city": {
                "type": "string",
                "description": "城市名称",
                "required": True
            }
        }


class FileOperationTool(Tool):
    """
    文件操作工具
    提供文件读写功能
    """
    
    def __init__(self):
        super().__init__(
            name="file_ops",
            description="执行文件操作（读取、写入、创建等）"
        )
    
    def execute(self, operation: str, file_path: str, content: Optional[str] = None) -> str:
        """
        执行文件操作
        
        Args:
            operation: 操作类型（read, write, create, delete）
            file_path: 文件路径
            content: 文件内容（写操作时需要）
            
        Returns:
            操作结果
        """
        try:
            logger.info(f"File operation: {operation} on {file_path}")
            
            if operation == "read":
                with open(file_path, 'r', encoding='utf-8') as f:
                    return f.read()
            
            elif operation == "write":
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content or "")
                return f"文件写入成功: {file_path}"
            
            elif operation == "create":
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content or "")
                return f"文件创建成功: {file_path}"
            
            elif operation == "delete":
                import os
                os.remove(file_path)
                return f"文件删除成功: {file_path}"
            
            else:
                return f"不支持的操作: {operation}"
                
        except Exception as e:
            logger.error(f"File operation failed: {e}")
            return f"文件操作失败: {str(e)}"
    
    def _get_parameters_schema(self) -> Dict[str, Any]:
        return {
            "operation": {
                "type": "string",
                "description": "操作类型 (read, write, create, delete)",
                "required": True
            },
            "file_path": {
                "type": "string",
                "description": "文件路径",
                "required": True
            },
            "content": {
                "type": "string",
                "description": "文件内容（写操作时需要）",
                "required": False
            }
        } 