# DAO 层说明

## 什么是 DAO

**DAO (Data Access Object)** 是**数据访问对象**，专门负责**数据库操作**。

## DAO 层的职责

### 1. 封装数据库操作
- 封装所有 SQL 操作
- 提供统一的数据访问接口
- 隐藏数据库实现细节

### 2. 数据持久化
- 执行 INSERT（创建）
- 执行 UPDATE（更新）
- 执行 DELETE（删除）
- 执行 SELECT（查询）

### 3. 数据转换
- 数据库结果 → Entity
- DO → 数据库记录

## DAO 层的特点

### 1. 只关注数据访问
- ❌ 不包含业务逻辑
- ❌ 不处理业务规则
- ✅ 只负责数据的增删改查

### 2. 使用 GoFrame ORM
- 使用 `g.DB()` 进行数据库操作
- 使用链式调用构建查询
- 自动处理 SQL 注入防护

### 3. 返回 Entity
- 查询操作返回 `*entity.User`
- 创建操作返回生成的 ID
- 更新/删除操作返回影响行数

## DAO 层代码示例

### 完整的 User DAO

```go
package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"hz/internal/model/do"
	"hz/internal/model/entity"
)

type userDao struct{}

var User = userDao{}

// GetById 根据ID获取用户
func (dao *userDao) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
	err = g.DB().Model("users").Ctx(ctx).Where("id", id).Scan(&user)
	return
}

// GetList 获取用户列表（支持分页）
func (dao *userDao) GetList(ctx context.Context, page, pageSize int) (users []*entity.User, total int, err error) {
	m := g.DB().Model("users").Ctx(ctx)
	
	// 获取总数
	total, err = m.Count()
	if err != nil {
		return
	}
	
	// 分页查询
	err = m.Page(page, pageSize).OrderDesc("id").Scan(&users)
	return
}

// Create 创建用户
func (dao *userDao) Create(ctx context.Context, data *do.User) (id uint, err error) {
	result, err := g.DB().Model("users").Ctx(ctx).Data(data).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return uint(result), nil
}

// Update 更新用户
func (dao *userDao) Update(ctx context.Context, id uint, data *do.User) (err error) {
	updateData := g.Map{}
	if data.Name != nil {
		updateData["name"] = *data.Name
	}
	if data.Email != nil {
		updateData["email"] = *data.Email
	}
	if data.Phone != nil {
		updateData["phone"] = *data.Phone
	}
	if len(updateData) == 0 {
		return nil
	}
	_, err = g.DB().Model("users").Ctx(ctx).Where("id", id).Data(updateData).Update()
	return
}

// Delete 删除用户
func (dao *userDao) Delete(ctx context.Context, id uint) (err error) {
	_, err = g.DB().Model("users").Ctx(ctx).Where("id", id).Delete()
	return
}
```

## DAO 层的方法说明

### 1. GetById - 根据ID查询

```go
func (dao *userDao) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
	err = g.DB().Model("users").Ctx(ctx).Where("id", id).Scan(&user)
	return
}
```

**作用**: 
- 根据主键ID查询单条记录
- 返回完整的 Entity 对象

**SQL等价**:
```sql
SELECT * FROM users WHERE id = ?
```

### 2. GetList - 列表查询（分页）

```go
func (dao *userDao) GetList(ctx context.Context, page, pageSize int) (users []*entity.User, total int, err error) {
	m := g.DB().Model("users").Ctx(ctx)
	total, err = m.Count()  // 获取总数
	err = m.Page(page, pageSize).OrderDesc("id").Scan(&users)  // 分页查询
	return
}
```

**作用**:
- 支持分页查询
- 返回总数和列表
- 按ID倒序排列

**SQL等价**:
```sql
-- 获取总数
SELECT COUNT(*) FROM users

-- 分页查询
SELECT * FROM users ORDER BY id DESC LIMIT ? OFFSET ?
```

### 3. Create - 创建记录

```go
func (dao *userDao) Create(ctx context.Context, data *do.User) (id uint, err error) {
	result, err := g.DB().Model("users").Ctx(ctx).Data(data).InsertAndGetId()
	return uint(result), nil
}
```

**作用**:
- 插入新记录
- 返回自动生成的ID
- 使用 DO 对象作为数据源

**SQL等价**:
```sql
INSERT INTO users (name, email, phone) VALUES (?, ?, ?)
-- 返回自增ID
```

### 4. Update - 更新记录

```go
func (dao *userDao) Update(ctx context.Context, id uint, data *do.User) (err error) {
	updateData := g.Map{}
	if data.Name != nil {
		updateData["name"] = *data.Name
	}
	if data.Email != nil {
		updateData["email"] = *data.Email
	}
	// 只更新非nil的字段
	_, err = g.DB().Model("users").Ctx(ctx).Where("id", id).Data(updateData).Update()
	return
}
```

**作用**:
- 部分更新（只更新设置的字段）
- 通过 DO 的指针判断哪些字段要更新
- 支持部分字段更新

**SQL等价**:
```sql
UPDATE users SET name=?, email=? WHERE id=?
```

### 5. Delete - 删除记录

```go
func (dao *userDao) Delete(ctx context.Context, id uint) (err error) {
	_, err = g.DB().Model("users").Ctx(ctx).Where("id", id).Delete()
	return
}
```

**作用**:
- 根据ID删除记录
- 物理删除（从数据库删除）

**SQL等价**:
```sql
DELETE FROM users WHERE id=?
```

## DAO 层的设计原则

### 1. 单一职责
- 只负责数据访问
- 不包含业务逻辑
- 不处理业务规则

### 2. 接口统一
- 所有方法接收 `context.Context`
- 统一的错误处理
- 统一的返回格式

### 3. 可复用性
- 可被多个 Service 调用
- 可被 Logic 层调用（用于验证）
- 方法粒度适中

## DAO 层在架构中的位置

```
┌─────────────────┐
│  Controller     │  处理HTTP请求
└────────┬────────┘
         ↓
┌─────────────────┐
│  Service        │  业务逻辑编排
└────────┬────────┘
    ↓         ↓
┌────────┐ ┌────────┐
│ Logic  │ │  DAO   │  业务逻辑 | 数据访问
└────────┘ └───┬────┘
               ↓
         ┌─────────┐
         │ Database│  数据库
         └─────────┘
```

## DAO 层的调用关系

### 1. 被 Service 层调用

```go
// Service层
func (s *userImpl) Create(ctx context.Context, data *do.User) (id uint, err error) {
	return dao.User.Create(ctx, data)  // 调用DAO层
}
```

### 2. 被 Logic 层调用（用于验证）

```go
// Logic层
func (l *userLogicImpl) PrepareUserForUpdate(ctx context.Context, id uint, data *do.User) error {
	// 调用DAO验证用户是否存在
	user, err := dao.User.GetById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}
	// ... 其他业务逻辑
}
```

## DAO 层 vs Service 层 vs Logic 层

| 特性 | DAO层 | Service层 | Logic层 |
|------|-------|-----------|---------|
| **职责** | 数据访问 | 业务编排 | 业务逻辑 |
| **操作** | CRUD操作 | 调用Logic和DAO | 业务规则验证 |
| **依赖** | Model层 | Logic层、DAO层 | DAO层（用于验证） |
| **返回** | Entity | Entity | 无返回值或bool |
| **SQL** | ✅ 直接操作 | ❌ 不操作 | ❌ 不操作 |

## 实际使用示例

### 场景1: Service 调用 DAO

```go
// Service层
func (s *userImpl) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
	// 直接调用DAO，不做任何业务处理
	return dao.User.GetById(ctx, id)
}
```

### 场景2: Logic 调用 DAO（验证）

```go
// Logic层
func (l *userLogicImpl) ValidateEmailUnique(ctx context.Context, email string) (bool, error) {
	// 调用DAO查询数据库验证
	var user *entity.User
	err := g.DB().Model("users").Ctx(ctx).Where("email", email).Scan(&user)
	if user != nil {
		return false, errors.New("邮箱已被使用")
	}
	return true, nil
}
```

### 场景3: 复杂查询（在DAO中）

```go
// DAO层可以添加复杂查询方法
func (dao *userDao) GetByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	err = g.DB().Model("users").Ctx(ctx).Where("email", email).Scan(&user)
	return
}

func (dao *userDao) GetByPhone(ctx context.Context, phone string) (user *entity.User, err error) {
	err = g.DB().Model("users").Ctx(ctx).Where("phone", phone).Scan(&user)
	return
}

func (dao *userDao) Search(ctx context.Context, keyword string) (users []*entity.User, err error) {
	err = g.DB().Model("users").Ctx(ctx).
		Where("name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Scan(&users)
	return
}
```

## DAO 层的优势

### 1. 数据访问集中化
- 所有数据库操作都在DAO层
- 便于统一管理和维护
- 便于性能优化

### 2. 解耦合
- Service层不直接操作数据库
- 业务逻辑与数据访问分离
- 便于测试和替换

### 3. 可复用性
- 一个DAO可被多个Service使用
- 可被Logic层使用（用于验证）
- 减少重复代码

### 4. 易于维护
- 数据库结构变化只需修改DAO
- SQL优化只需修改DAO
- 不影响业务逻辑

## 最佳实践

### 1. 方法命名规范
- `GetById` - 根据ID获取
- `GetList` - 获取列表
- `Create` - 创建
- `Update` - 更新
- `Delete` - 删除
- `GetByXxx` - 根据某个字段获取

### 2. 参数和返回值
- 所有方法接收 `context.Context`
- 查询方法返回 `*entity.User` 或 `[]*entity.User`
- 创建方法返回生成的ID
- 更新/删除方法返回 `error`

### 3. 错误处理
- 数据库错误直接返回
- 不在这里处理业务错误
- 让上层（Service/Logic）处理业务逻辑

### 4. 性能优化
- 使用索引字段查询
- 合理使用分页
- 避免N+1查询问题

## 总结

**DAO层是数据访问的抽象层**，它：

1. ✅ **封装数据库操作** - 所有SQL操作都在这里
2. ✅ **提供统一接口** - 便于Service和Logic调用
3. ✅ **隐藏实现细节** - 上层不需要知道SQL细节
4. ✅ **易于维护** - 数据库变化只需修改DAO
5. ✅ **可复用** - 可被多个Service和Logic使用

**记住**: DAO层只负责数据访问，不包含业务逻辑！

