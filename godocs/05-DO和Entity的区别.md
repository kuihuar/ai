# DO 和 Entity 的区别

## 什么是 DO (Data Object)

**DO (Data Object)** 是**数据对象**，专门用于**数据库操作**。

## DO 的特点

### 1. 字段类型
- 使用**指针类型** (`*string`, `*uint` 等)
- 指针为 `nil` 表示字段未设置
- 便于区分"零值"和"未设置"

### 2. ORM 标签
- 使用 `orm` 标签定义数据库字段映射
- 使用 `g.Meta` 定义表名和 DO 标识

### 3. 使用场景
- **写入操作**: Create, Update
- **数据库交互**: 与数据库表结构对应
- **部分更新**: 通过指针判断哪些字段需要更新

## 什么是 Entity (实体)

**Entity** 是**实体对象**，用于**业务逻辑**和**数据展示**。

## Entity 的特点

### 1. 字段类型
- 使用**值类型** (`string`, `uint` 等)
- 包含完整的业务字段
- 包含时间戳等业务字段

### 2. JSON 标签
- 使用 `json` 标签用于序列化
- 用于 API 响应

### 3. 使用场景
- **读取操作**: GetById, GetList
- **业务逻辑**: Logic 层使用
- **API 响应**: Controller 返回给客户端

## 对比示例

### DO (Data Object)

```go
// internal/model/do/user.go
type User struct {
    g.Meta `orm:"table:users, do:true"`  // 定义表名和DO标识
    Id     *uint   `orm:"id,primary"`    // 指针类型
    Name   *string `orm:"name"`           // 指针类型
    Email  *string `orm:"email"`          // 指针类型
    Phone  *string `orm:"phone"`         // 指针类型
}
```

**特点**:
- ✅ 字段都是指针类型
- ✅ 使用 `orm` 标签
- ✅ 只包含需要操作的字段
- ✅ 不包含时间戳等自动字段

### Entity (实体)

```go
// internal/model/entity/user.go
type User struct {
    Id        uint   `json:"id"`         // 值类型
    Name      string `json:"name"`       // 值类型
    Email     string `json:"email"`      // 值类型
    Phone     string `json:"phone"`      // 值类型
    CreatedAt string `json:"createdAt"`  // 业务字段
    UpdatedAt string `json:"updatedAt"` // 业务字段
}
```

**特点**:
- ✅ 字段都是值类型
- ✅ 使用 `json` 标签
- ✅ 包含完整的业务字段
- ✅ 包含时间戳等业务字段

## 使用场景对比

### 场景1: 创建用户

**使用 DO**:
```go
// Controller层
data := &do.User{
    Name:  &req.Name,   // 指针类型
    Email: &req.Email,  // 指针类型
    Phone: &req.Phone,  // 指针类型
}
service.User.Create(ctx, data)
```

**为什么用 DO**:
- 只需要设置要插入的字段
- 指针类型便于判断字段是否设置
- 与数据库表结构对应

### 场景2: 查询用户

**使用 Entity**:
```go
// DAO层
func (dao *userDao) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
    err = g.DB().Model("users").Ctx(ctx).Where("id", id).Scan(&user)
    return
}
```

**为什么用 Entity**:
- 返回完整的业务对象
- 包含所有字段（包括时间戳）
- 直接用于 API 响应

### 场景3: 部分更新

**使用 DO**:
```go
// Controller层
data := &do.User{}
if req.Name != "" {
    data.Name = &req.Name  // 只设置要更新的字段
}
if req.Email != "" {
    data.Email = &req.Email
}
// Phone 不设置，表示不更新该字段
service.User.Update(ctx, id, data)
```

**为什么用 DO**:
- 指针为 `nil` 表示不更新该字段
- 指针不为 `nil` 表示要更新该字段
- 便于实现部分更新

## 数据流转

### 创建用户流程

```
1. HTTP请求 (JSON)
   {"name":"张三","email":"zhang@example.com"}
   ↓
2. Controller层
   do.User{Name:&"张三", Email:&"zhang@example.com"}
   ↓ (使用 DO)
3. Service层 → Logic层 → DAO层
   都使用 do.User
   ↓
4. 数据库
   INSERT INTO users (name, email) VALUES ...
   ↓
5. 返回 Entity
   entity.User{Id:1, Name:"张三", Email:"zhang@example.com", ...}
   ↓
6. HTTP响应 (JSON)
   {"id":1,"name":"张三","email":"zhang@example.com",...}
```

### 查询用户流程

```
1. HTTP请求
   GET /user/1
   ↓
2. DAO层查询
   使用 do.User 作为查询条件（如果需要）
   ↓
3. 数据库查询
   SELECT * FROM users WHERE id=1
   ↓
4. 返回 Entity
   entity.User{Id:1, Name:"张三", ...}
   ↓
5. HTTP响应
   {"id":1,"name":"张三",...}
```

## 为什么需要两种类型？

### 1. 职责分离

- **DO**: 专注于数据库操作
- **Entity**: 专注于业务逻辑和展示

### 2. 灵活性

- **DO**: 指针类型支持部分更新
- **Entity**: 值类型更符合业务语义

### 3. 类型安全

- **DO**: 明确标识用于数据库操作
- **Entity**: 明确标识用于业务逻辑

## 实际代码示例

### 创建用户（使用 DO）

```go
// Controller层
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
    // 使用 DO 构建数据
    data := &do.User{
        Name:  &req.Name,   // 指针类型
        Email: &req.Email,
        Phone: &req.Phone,
    }
    id, err := service.User.Create(ctx, data)
    return &v1.CreateRes{Id: id}, nil
}
```

### 查询用户（返回 Entity）

```go
// DAO层
func (dao *userDao) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
    // 查询返回 Entity
    err = g.DB().Model("users").Ctx(ctx).Where("id", id).Scan(&user)
    return
}
```

### 部分更新（使用 DO）

```go
// Controller层
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
    data := &do.User{}
    // 只设置要更新的字段
    if req.Name != "" {
        data.Name = &req.Name  // 设置指针，表示要更新
    }
    // Phone 不设置（指针为 nil），表示不更新
    err = service.User.Update(ctx, req.Id, data)
    return &v1.UpdateRes{Success: true}, nil
}
```

## 总结

| 特性 | DO (Data Object) | Entity (实体) |
|------|-----------------|---------------|
| **用途** | 数据库操作 | 业务逻辑和展示 |
| **字段类型** | 指针类型 (`*string`) | 值类型 (`string`) |
| **标签** | `orm` 标签 | `json` 标签 |
| **使用场景** | Create, Update | GetById, GetList, API响应 |
| **字段完整性** | 只包含操作字段 | 包含完整业务字段 |
| **部分更新** | ✅ 支持（通过指针） | ❌ 不支持 |
| **时间戳** | ❌ 不包含 | ✅ 包含 |

## 最佳实践

1. **写入操作** (Create, Update) → 使用 **DO**
2. **读取操作** (GetById, GetList) → 返回 **Entity**
3. **业务逻辑** → 使用 **Entity**
4. **API 响应** → 使用 **Entity**

这样的设计使得代码更加清晰，职责更加明确！

