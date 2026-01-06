# 文档编写

## 文档类型

### 1. API 文档
- **内容**：接口定义、参数说明、示例
- **工具**：Swagger/OpenAPI
- **位置**：`openapi.yaml` 或自动生成

### 2. 架构文档
- **内容**：系统架构、组件设计、数据流
- **格式**：Markdown + 图表
- **位置**：`docs/architecture/`

### 3. 开发文档
- **内容**：开发流程、代码规范、最佳实践
- **格式**：Markdown
- **位置**：`docs/development/`

### 4. 运维文档
- **内容**：部署流程、配置说明、故障处理
- **格式**：Markdown
- **位置**：`docs/operations/`

### 5. 代码注释
- **内容**：函数说明、参数说明、示例
- **格式**：Go 注释
- **位置**：代码文件中

## 文档编写原则

### 1. 清晰明确
- 使用简洁的语言
- 避免歧义
- 提供具体示例

### 2. 结构完整
- 有清晰的目录结构
- 章节划分合理
- 内容层次分明

### 3. 及时更新
- 代码变更时同步更新文档
- 定期审查文档准确性
- 标记过时内容

### 4. 易于查找
- 使用清晰的标题
- 提供索引和链接
- 支持搜索

## Markdown 编写规范

### 标题层级
```markdown
# 一级标题（文档标题）
## 二级标题（主要章节）
### 三级标题（子章节）
#### 四级标题（细节说明）
```

### 代码块
```markdown
```go
// Go 代码示例
func example() {
    // ...
}
```

```bash
# Shell 命令示例
go run main.go
```
```

### 列表
```markdown
- 无序列表项
- 另一个列表项

1. 有序列表项
2. 另一个列表项
```

### 链接和引用
```markdown
[链接文本](URL)
![图片描述](图片URL)
> 引用内容
```

## 代码注释规范

### 包注释
```go
// Package example provides utilities for example operations.
package example
```

### 函数注释
```go
// GetUser retrieves a user by ID.
// It returns an error if the user is not found.
func GetUser(ctx context.Context, id int64) (*User, error) {
    // ...
}
```

### 复杂逻辑注释
```go
// 计算用户积分：
// 1. 基础积分 = 注册天数 * 10
// 2. 奖励积分 = 完成任务数 * 5
// 3. 总积分 = 基础积分 + 奖励积分
func calculatePoints(user *User) int {
    // ...
}
```

## 最佳实践

1. **从 README 开始**：每个项目都应该有清晰的 README
2. **示例驱动**：提供丰富的代码示例
3. **图表辅助**：使用图表说明复杂概念
4. **版本管理**：文档随代码一起版本管理
5. **定期审查**：定期审查和更新文档

