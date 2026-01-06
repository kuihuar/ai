# Go 代码规范

## 命名规范

### 包名
- 小写字母，简短有意义
- 避免下划线和混合大小写
- 单数形式，避免复数

### 变量名
- 驼峰命名（camelCase）
- 首字母大写表示导出，小写表示私有
- 简短但有意义

### 函数名
- 驼峰命名
- 导出函数首字母大写
- Getter/Setter 遵循约定

### 接口名
- 通常以 `-er` 结尾（如 `Reader`, `Writer`）
- 或使用描述性名称

## 代码组织

### 文件结构
```go
package example

import (
    // 标准库
    "fmt"
    "context"
    
    // 第三方库
    "github.com/go-kratos/kratos/v2"
    
    // 项目内部
    "sre/internal/biz"
)
```

### 函数长度
- 单个函数不超过 50 行
- 复杂逻辑拆分为多个小函数
- 保持函数职责单一

### 错误处理
- 总是检查错误
- 错误信息要清晰
- 使用 `fmt.Errorf` 包装错误，添加上下文

```go
if err != nil {
    return nil, fmt.Errorf("failed to save user: %w", err)
}
```

## 注释规范

### 包注释
每个包都应该有包注释，说明包的用途。

```go
// Package example provides utilities for example operations.
package example
```

### 导出函数注释
所有导出的函数、类型、变量都应该有注释。

```go
// GetUser retrieves a user by ID.
// It returns an error if the user is not found.
func GetUser(id int64) (*User, error) {
    // ...
}
```

## 最佳实践

### 1. 使用 context
- 所有可能长时间运行的函数都应该接受 `context.Context`
- 使用 context 传递请求范围的值
- 使用 context 控制超时和取消

### 2. 避免全局变量
- 使用依赖注入传递依赖
- 避免使用包级别的全局变量

### 3. 接口设计
- 接口应该小而专注
- 接口定义在使用方，而非实现方

### 4. 错误处理
- 使用 `errors.Is` 和 `errors.As` 检查错误
- 使用 `fmt.Errorf` 和 `%w` 包装错误
- 提供有意义的错误信息

### 5. 并发安全
- 明确并发安全的要求
- 使用适当的同步原语
- 避免数据竞争

