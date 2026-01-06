# API 响应定义说明

## API 响应结构

在 GoFrame 项目中，API 响应通过 **API 层的 Res 结构体**定义，框架会自动处理响应格式。

## 响应定义位置

### API 层定义响应结构

```go
// api/user/v1/user.go

// CreateRes 创建用户响应
type CreateRes struct {
	Id uint `json:"id" dc:"用户ID"`
}

// GetByIdRes 根据ID获取用户响应
type GetByIdRes struct {
	*entity.User  // 嵌入实体对象
}

// GetListRes 获取用户列表响应
type GetListRes struct {
	List  []*entity.User `json:"list"  dc:"用户列表"`
	Total int            `json:"total" dc:"总数"`
	Page  int            `json:"page"  dc:"当前页码"`
	Size  int            `json:"size"  dc:"每页数量"`
}
```

## 响应格式

### 1. 统一响应格式（框架自动包装）

GoFrame 框架会自动将 Controller 返回的 Res 包装成统一格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 这里是 Res 结构体的内容
  }
}
```

### 2. 实际响应示例

#### 创建用户响应

**定义**:
```go
type CreateRes struct {
	Id uint `json:"id" dc:"用户ID"`
}
```

**Controller返回**:
```go
return &v1.CreateRes{Id: id}, nil
```

**实际HTTP响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

#### 根据ID获取用户响应

**定义**:
```go
type GetByIdRes struct {
	*entity.User
}
```

**Controller返回**:
```go
return &v1.GetByIdRes{User: user}, nil
```

**实际HTTP响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "张三",
    "email": "zhang@example.com",
    "phone": "13800138000",
    "createdAt": "2024-01-01 10:00:00",
    "updatedAt": "2024-01-01 10:00:00"
  }
}
```

#### 获取用户列表响应

**定义**:
```go
type GetListRes struct {
	List  []*entity.User `json:"list"  dc:"用户列表"`
	Total int            `json:"total" dc:"总数"`
	Page  int            `json:"page"  dc:"当前页码"`
	Size  int            `json:"size"  dc:"每页数量"`
}
```

**Controller返回**:
```go
return &v1.GetListRes{
	List:  users,
	Total: total,
	Page:  req.Page,
	Size:  req.PageSize,
}, nil
```

**实际HTTP响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "张三",
        "email": "zhang@example.com",
        "phone": "13800138000",
        "createdAt": "2024-01-01 10:00:00",
        "updatedAt": "2024-01-01 10:00:00"
      },
      {
        "id": 2,
        "name": "李四",
        "email": "lisi@example.com",
        "phone": "13900139000",
        "createdAt": "2024-01-02 10:00:00",
        "updatedAt": "2024-01-02 10:00:00"
      }
    ],
    "total": 2,
    "page": 1,
    "size": 10
  }
}
```

## 响应定义方式

### 方式1: 直接定义字段

```go
type CreateRes struct {
	Id uint `json:"id" dc:"用户ID"`
}
```

**适用场景**: 简单响应，只有几个字段

### 方式2: 嵌入实体对象

```go
type GetByIdRes struct {
	*entity.User
}
```

**适用场景**: 返回完整的实体对象

**注意**: 使用指针 `*entity.User`，如果 user 为 nil，响应中 data 也为 null

### 方式3: 组合多个字段

```go
type GetListRes struct {
	List  []*entity.User `json:"list"  dc:"用户列表"`
	Total int            `json:"total" dc:"总数"`
	Page  int            `json:"page"  dc:"当前页码"`
	Size  int            `json:"size"  dc:"每页数量"`
}
```

**适用场景**: 需要返回多个字段（如分页信息）

### 方式4: 自定义响应结构

```go
type CustomRes struct {
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
	Data    interface{} `json:"data" dc:"数据"`
}
```

**适用场景**: 需要自定义响应格式

## 标签说明

### json 标签

```go
Id uint `json:"id"`
```

**作用**: 定义 JSON 序列化时的字段名

**示例**:
- `json:"id"` → JSON 中字段名为 "id"
- `json:"user_id"` → JSON 中字段名为 "user_id"
- `json:"-"` → 不序列化该字段

### dc 标签

```go
Id uint `json:"id" dc:"用户ID"`
```

**作用**: 定义字段的描述（用于 Swagger 文档）

**示例**: Swagger UI 中会显示 "用户ID" 作为字段说明

## Controller 层返回响应

### 1. 成功响应

```go
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err  // 返回错误
	}
	return &v1.CreateRes{Id: id}, nil  // 返回成功响应
}
```

### 2. 错误响应

```go
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err  // 返回错误，框架自动处理
	}
	return &v1.CreateRes{Id: id}, nil
}
```

**框架自动处理**:
```json
{
  "code": 1,
  "message": "错误信息",
  "data": null
}
```

## 统一响应格式

### 中间件处理

在 `cmd.go` 中注册的中间件会自动处理响应格式：

```go
s.Group("/", func(group *ghttp.RouterGroup) {
	group.Middleware(ghttp.MiddlewareHandlerResponse)  // 统一响应处理
	group.Bind(
		user.NewV1(),
	)
})
```

### 响应格式说明

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // Res 结构体的内容
  }
}
```

#### 错误响应

```json
{
  "code": 1,
  "message": "错误信息",
  "data": null
}
```

## 完整示例

### 示例1: 创建用户

**API定义**:
```go
type CreateRes struct {
	Id uint `json:"id" dc:"用户ID"`
}
```

**Controller实现**:
```go
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	data := &do.User{
		Name:  &req.Name,
		Email: &req.Email,
		Phone: &req.Phone,
	}
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
```

**HTTP响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

### 示例2: 获取用户列表

**API定义**:
```go
type GetListRes struct {
	List  []*entity.User `json:"list"  dc:"用户列表"`
	Total int            `json:"total" dc:"总数"`
	Page  int            `json:"page"  dc:"当前页码"`
	Size  int            `json:"size"  dc:"每页数量"`
}
```

**Controller实现**:
```go
func (c *ControllerV1) GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error) {
	users, total, err := service.UserEnhanced.GetList(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	return &v1.GetListRes{
		List:  users,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}
```

**HTTP响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "张三",
        "email": "zhang@example.com",
        "phone": "13800138000",
        "createdAt": "2024-01-01 10:00:00",
        "updatedAt": "2024-01-01 10:00:00"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

### 示例3: 简单响应

**API定义**:
```go
type UpdateRes struct {
	Success bool `json:"success" dc:"是否成功"`
}
```

**Controller实现**:
```go
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = service.UserEnhanced.UpdateWithValidation(ctx, req.Id, data)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{Success: true}, nil
}
```

**HTTP响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true
  }
}
```

## 响应字段类型

### 基本类型

```go
type SimpleRes struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Price float64 `json:"price"`
	Active bool  `json:"active"`
}
```

### 数组类型

```go
type ListRes struct {
	Items []string `json:"items"`
}
```

### 对象类型

```go
type ObjectRes struct {
	User *entity.User `json:"user"`
}
```

### 对象数组

```go
type ListRes struct {
	List []*entity.User `json:"list"`
}
```

### 嵌套对象

```go
type NestedRes struct {
	User struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}
```

## 最佳实践

### 1. 使用 Entity 对象

```go
// ✅ 推荐：使用实体对象
type GetByIdRes struct {
	*entity.User
}

// ❌ 不推荐：重复定义字段
type GetByIdRes struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	// ... 重复定义所有字段
}
```

### 2. 分页响应

```go
// ✅ 推荐：包含分页信息
type GetListRes struct {
	List  []*entity.User `json:"list"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}
```

### 3. 错误处理

```go
// ✅ 推荐：直接返回错误
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := service.UserEnhanced.CreateWithValidation(ctx, data)
	if err != nil {
		return nil, err  // 框架自动处理错误响应
	}
	return &v1.CreateRes{Id: id}, nil
}
```

### 4. 使用 dc 标签

```go
// ✅ 推荐：添加描述
type CreateRes struct {
	Id uint `json:"id" dc:"用户ID"`
}

// ❌ 不推荐：缺少描述
type CreateRes struct {
	Id uint `json:"id"`
}
```

## Swagger 文档

响应定义会自动生成 Swagger 文档：

- 访问: http://localhost:8000/swagger
- 查看每个 API 的响应结构
- `dc` 标签会显示为字段说明

## 总结

### API 响应定义步骤

1. **在 API 层定义 Res 结构体**
   ```go
   type CreateRes struct {
       Id uint `json:"id" dc:"用户ID"`
   }
   ```

2. **在 Controller 层返回 Res**
   ```go
   return &v1.CreateRes{Id: id}, nil
   ```

3. **框架自动处理**
   - 自动包装成统一格式
   - 自动序列化为 JSON
   - 自动生成 Swagger 文档

### 响应格式

- **成功**: `{"code":0, "message":"success", "data":{...}}`
- **错误**: `{"code":1, "message":"错误信息", "data":null}`

### 关键点

1. ✅ 响应定义在 **API 层**（Res 结构体）
2. ✅ Controller 返回 **Res 结构体实例**
3. ✅ 框架自动**包装和序列化**
4. ✅ 使用 **json 标签**定义字段名
5. ✅ 使用 **dc 标签**添加描述（Swagger）

**记住**: API 响应在 API 层定义，Controller 层只需返回 Res 结构体实例！

