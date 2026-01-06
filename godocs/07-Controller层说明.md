# Controller 层说明

## 什么是 Controller

**Controller (控制器)** 是**控制器层**，专门负责**处理 HTTP 请求**。

## Controller 层的职责

### 1. 接收 HTTP 请求
- 接收客户端发送的 HTTP 请求
- 解析请求参数（路径参数、查询参数、请求体）
- 框架自动进行参数验证

### 2. 参数转换
- API 层结构（Req）→ DO 层结构（do.User）
- 将 HTTP 请求数据转换为业务层可用的格式

### 3. 调用 Service 层
- 调用 Service 层处理业务逻辑
- 不直接操作数据库
- 不包含业务逻辑

### 4. 返回 HTTP 响应
- 将 Service 层返回的数据转换为 API 响应结构（Res）
- 返回给客户端

## Controller 层的特点

### 1. 薄层设计
- 代码简洁，逻辑简单
- 主要做参数转换和调用 Service
- 不包含复杂业务逻辑

### 2. 依赖 API 层和 Service 层
- 依赖 API 层：使用 Req/Res 结构
- 依赖 Service 层：调用业务逻辑
- 依赖 Model/DO 层：用于参数转换

### 3. 框架自动路由
- 通过 `g.Meta` 标签自动注册路由
- 框架自动处理参数验证
- 框架自动处理响应格式

## Controller 层代码结构

### 1. 控制器定义

```go
// internal/controller/user/user_new.go
package user

import (
	"hz/api/user"
)

type ControllerV1 struct{}

func NewV1() user.IUserV1 {
	return &ControllerV1{}
}
```

**作用**:
- 定义控制器结构体
- 实现 API 层定义的接口
- 提供构造函数

### 2. 控制器方法

每个 HTTP 请求对应一个控制器方法：

```go
// internal/controller/user/user_v1_create.go
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	// 1. 参数转换
	data := &do.User{
		Name:  &req.Name,
		Email: &req.Email,
		Phone: &req.Phone,
	}
	
	// 2. 调用Service层
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err
	}
	
	// 3. 返回响应
	return &v1.CreateRes{Id: id}, nil
}
```

## Controller 层完整示例

### 1. Create - 创建用户

```go
// internal/controller/user/user_v1_create.go
package user

import (
	"context"
	v1 "hz/api/user/v1"
	"hz/internal/model/do"
	"hz/internal/service"
)

// Create 创建用户
// 完整调用链: Controller -> Service -> Logic -> DAO
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	// 1. Controller层：参数转换（API层 -> DO层）
	data := &do.User{
		Name:  &req.Name,
		Email: &req.Email,
		Phone: &req.Phone,
	}
	
	// 2. 调用Service层（带业务逻辑验证的增强版）
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err
	}
	
	// 3. 返回响应（DO层数据 -> API层结构）
	return &v1.CreateRes{Id: id}, nil
}
```

**流程**:
1. 接收 `v1.CreateReq`（来自 HTTP 请求体）
2. 转换为 `do.User`（用于业务层）
3. 调用 `service.UserEnhanced.CreateWithValidation()`
4. 返回 `v1.CreateRes`（HTTP 响应）

### 2. GetById - 根据ID获取用户

```go
// internal/controller/user/user_v1_get_by_id.go
package user

import (
	"context"
	v1 "hz/api/user/v1"
	"hz/internal/service"
)

func (c *ControllerV1) GetById(ctx context.Context, req *v1.GetByIdReq) (res *v1.GetByIdRes, err error) {
	// 直接调用Service层，无需参数转换
	user, err := service.UserEnhanced.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	// 返回Entity（自动转换为API响应）
	return &v1.GetByIdRes{User: user}, nil
}
```

**流程**:
1. 接收路径参数 `req.Id`
2. 直接调用 Service 层
3. 返回 Entity（自动序列化为 JSON）

### 3. GetList - 获取用户列表

```go
// internal/controller/user/user_v1_get_list.go
package user

import (
	"context"
	v1 "hz/api/user/v1"
	"hz/internal/service"
)

func (c *ControllerV1) GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error) {
	// 调用Service层
	users, total, err := service.UserEnhanced.GetList(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	
	// 构建响应
	return &v1.GetListRes{
		List:  users,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}
```

**流程**:
1. 接收查询参数（page, pageSize）
2. 调用 Service 层获取列表
3. 构建分页响应

### 4. Update - 更新用户

```go
// internal/controller/user/user_v1_update.go
package user

import (
	"context"
	v1 "hz/api/user/v1"
	"hz/internal/model/do"
	"hz/internal/service"
)

// Update 更新用户
// 完整调用链: Controller -> Service -> Logic -> DAO
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	// 1. Controller层：参数转换（API层 -> DO层）
	data := &do.User{}
	
	// 只设置要更新的字段（部分更新）
	if req.Name != "" {
		data.Name = &req.Name
	}
	if req.Email != "" {
		data.Email = &req.Email
	}
	if req.Phone != "" {
		data.Phone = &req.Phone
	}
	
	// 2. 调用Service层（带业务逻辑验证的增强版）
	err = service.UserEnhanced.UpdateWithValidation(ctx, req.Id, data)
	if err != nil {
		return nil, err
	}
	
	// 3. 返回响应
	return &v1.UpdateRes{Success: true}, nil
}
```

**流程**:
1. 接收路径参数（id）和请求体（name, email, phone）
2. 转换为 DO 对象（只设置要更新的字段）
3. 调用 Service 层
4. 返回成功响应

### 5. Delete - 删除用户

```go
// internal/controller/user/user_v1_delete.go
package user

import (
	"context"
	v1 "hz/api/user/v1"
	"hz/internal/service"
)

// Delete 删除用户
// 完整调用链: Controller -> Service -> Logic -> DAO
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	// 1. Controller层：直接传递参数
	
	// 2. 调用Service层（带业务逻辑验证的增强版）
	err = service.UserEnhanced.DeleteWithValidation(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	// 3. 返回响应
	return &v1.DeleteRes{Success: true}, nil
}
```

**流程**:
1. 接收路径参数（id）
2. 直接调用 Service 层
3. 返回成功响应

## Controller 层的路由注册

### 1. API 层定义路由

```go
// api/user/v1/user.go
type CreateReq struct {
	g.Meta `path:"/user" method:"post" tags:"User" summary:"创建用户"`
	Name   string `json:"name" v:"required"`
	Email  string `json:"email" v:"required|email"`
}
```

**路由信息**:
- `path:"/user"` - 路由路径
- `method:"post"` - HTTP 方法
- `tags:"User"` - Swagger 标签
- `summary:"创建用户"` - API 摘要

### 2. Controller 实现接口

```go
// internal/controller/user/user_new.go
func NewV1() user.IUserV1 {
	return &ControllerV1{}
}
```

### 3. 在 cmd.go 中注册

```go
// internal/cmd/cmd.go
s.Group("/", func(group *ghttp.RouterGroup) {
	group.Middleware(ghttp.MiddlewareHandlerResponse)
	group.Bind(
		hello.NewV1(),
		user.NewV1(),  // 注册用户控制器
	)
})
```

**框架自动**:
- 根据 `g.Meta` 标签注册路由
- 自动参数验证
- 自动处理响应

## Controller 层的数据流转

### 创建用户完整流程

```
1. HTTP请求
   POST /user
   Content-Type: application/json
   Body: {"name":"张三","email":"zhang@example.com","phone":"13800138000"}
   ↓
2. GoFrame框架
   - 解析请求体
   - 参数验证（根据v标签）
   - 创建 v1.CreateReq 对象
   ↓
3. Controller层
   user_v1_create.go: Create()
   - 接收: v1.CreateReq{Name:"张三", Email:"zhang@example.com", Phone:"13800138000"}
   - 转换: do.User{Name:&"张三", Email:&"zhang@example.com", Phone:&"13800138000"}
   - 调用: service.UserEnhanced.CreateWithValidation(ctx, data)
   ↓
4. Service层
   - 业务逻辑处理
   - 调用DAO层
   ↓
5. 返回路径
   Service返回: id=1
   Controller返回: v1.CreateRes{Id:1}
   框架序列化: {"id":1}
   HTTP响应: 200 OK {"id":1}
```

## Controller 层的设计原则

### 1. 薄层设计
- 代码简洁，逻辑简单
- 主要做参数转换
- 不包含业务逻辑

### 2. 单一职责
- 只负责 HTTP 请求处理
- 不直接操作数据库
- 不包含复杂业务逻辑

### 3. 依赖注入
- 依赖 Service 层
- 通过接口调用，便于测试

### 4. 错误处理
- 直接返回 Service 层的错误
- 框架自动处理错误响应

## Controller 层 vs Service 层 vs Logic 层

| 特性 | Controller层 | Service层 | Logic层 |
|------|--------------|-----------|---------|
| **职责** | HTTP请求处理 | 业务编排 | 业务逻辑 |
| **操作** | 参数转换 | 调用Logic和DAO | 业务规则验证 |
| **依赖** | API层、Service层 | Logic层、DAO层 | DAO层（验证） |
| **返回** | API响应结构 | Entity | 无或bool |
| **HTTP** | ✅ 处理HTTP | ❌ 不处理 | ❌ 不处理 |

## Controller 层的优势

### 1. 职责清晰
- 只负责 HTTP 请求处理
- 不包含业务逻辑
- 易于理解和维护

### 2. 易于测试
- 可以 Mock Service 层
- 独立测试 HTTP 处理逻辑
- 不依赖数据库

### 3. 易于扩展
- 新增接口只需添加方法
- 不影响业务逻辑
- 框架自动处理路由

### 4. 统一响应格式
- 框架自动处理响应格式
- 统一的错误处理
- 自动生成 Swagger 文档

## 最佳实践

### 1. 方法命名
- 与 API 层接口方法名一致
- 清晰表达功能

### 2. 参数转换
- API 层结构 → DO 层结构
- 只转换必要的字段
- 保持数据一致性

### 3. 错误处理
- 直接返回 Service 层的错误
- 不在这里处理业务错误
- 让框架处理 HTTP 错误响应

### 4. 代码简洁
- 保持方法简短
- 只做参数转换和调用 Service
- 不包含复杂逻辑

## 常见模式

### 模式1: 简单查询（无需转换）

```go
func (c *ControllerV1) GetById(ctx context.Context, req *v1.GetByIdReq) (res *v1.GetByIdRes, err error) {
	user, err := service.UserEnhanced.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetByIdRes{User: user}, nil
}
```

### 模式2: 创建/更新（需要转换）

```go
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	data := &do.User{
		Name:  &req.Name,
		Email: &req.Email,
	}
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
```

### 模式3: 部分更新（条件转换）

```go
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	data := &do.User{}
	if req.Name != "" {
		data.Name = &req.Name
	}
	if req.Email != "" {
		data.Email = &req.Email
	}
	err = service.UserEnhanced.UpdateWithValidation(ctx, req.Id, data)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{Success: true}, nil
}
```

## 总结

**Controller层是HTTP请求的入口**，它：

1. ✅ **接收HTTP请求** - 解析请求参数
2. ✅ **参数转换** - API层结构 → DO层结构
3. ✅ **调用Service层** - 处理业务逻辑
4. ✅ **返回HTTP响应** - DO层数据 → API层结构
5. ✅ **薄层设计** - 代码简洁，职责单一

**记住**: Controller层只负责HTTP请求处理，不包含业务逻辑！

