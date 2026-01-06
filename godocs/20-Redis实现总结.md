# Redis 实现总结

## 已实现的Redis功能

### 1. 配置文件
- ✅ 在 `manifest/config/config.yaml` 中添加Redis配置

### 2. Cache Service层
- ✅ `internal/service/cache.go` - 完整的缓存服务封装

### 3. 用户缓存服务
- ✅ `internal/service/user_cache.go` - 带缓存的用户服务示例

## 配置文件

```yaml
# manifest/config/config.yaml
redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    pass: ""  # 密码（如果有）
    minIdle: 10
    maxIdle: 100
    maxActive: 200
    idleTimeout: "10s"
    maxConnLifetime: "30s"
```

## Cache Service功能

### 基础操作
- ✅ `Set()` - 设置缓存（支持过期时间）
- ✅ `Get()` - 获取缓存（字符串）
- ✅ `GetObject()` - 获取缓存（对象，自动JSON序列化）
- ✅ `Delete()` - 删除缓存
- ✅ `DeletePattern()` - 按模式删除缓存
- ✅ `Exists()` - 检查key是否存在

### 高级操作
- ✅ `Expire()` - 设置过期时间
- ✅ `TTL()` - 获取剩余过期时间
- ✅ `Incr()` - 增加计数
- ✅ `Decr()` - 减少计数
- ✅ `HSet()` - 设置哈希字段
- ✅ `HGet()` - 获取哈希字段
- ✅ `HGetAll()` - 获取所有哈希字段

## 使用示例

### 1. 基础缓存操作

```go
// 设置缓存
err := service.Cache.Set(ctx, "key", "value", 1*time.Hour)

// 获取缓存（字符串）
val, err := service.Cache.Get(ctx, "key")

// 获取缓存（对象）
var user entity.User
err := service.Cache.GetObject(ctx, "user:1", &user)

// 删除缓存
err := service.Cache.Delete(ctx, "key")
```

### 2. 带缓存的用户查询

```go
// 获取用户（自动缓存）
user, err := service.UserCache.GetByIdWithCache(ctx, 1)
// 第一次：查数据库，写入缓存
// 后续：直接从缓存读取

// 获取用户列表（自动缓存）
users, total, err := service.UserCache.GetListWithCache(ctx, 1, 10)
```

### 3. 缓存失效

```go
// 更新用户后，使缓存失效
service.UserCache.InvalidateUserCache(ctx, id)
// 清除：user:1 和 users:list:*
```

### 4. 计数器

```go
// 增加计数
count, err := service.Cache.Incr(ctx, "user:1:view_count")

// 减少计数
count, err := service.Cache.Decr(ctx, "counter:key")
```

### 5. 哈希操作

```go
// 设置哈希字段
err := service.Cache.HSet(ctx, "user:1", "name", "张三")

// 获取哈希字段
name, err := service.Cache.HGet(ctx, "user:1", "name")

// 获取所有哈希字段
all, err := service.Cache.HGetAll(ctx, "user:1")
```

## 缓存键命名规范

```go
// 单个对象
"user:1"                    // 用户信息
"user:1:profile"           // 用户资料
"order:123"                // 订单信息

// 列表数据
"users:list:page:1:size:10"  // 用户列表
"orders:list:status:paid"    // 订单列表

// 计数器
"user:1:view_count"        // 用户浏览量
"order:123:status"         // 订单状态
```

## 缓存策略

### Cache-Aside（旁路缓存）

```go
// 读取：先查缓存，未命中查数据库，再写入缓存
func GetByIdWithCache(ctx context.Context, id uint) {
	// 1. 查缓存
	user := cache.Get("user:1")
	if user != nil {
		return user
	}
	
	// 2. 查数据库
	user = db.GetById(id)
	
	// 3. 写入缓存
	cache.Set("user:1", user, 1*time.Hour)
	
	return user
}
```

## 在Service层集成缓存

### 示例：用户服务集成缓存

```go
// 获取用户（带缓存）
func (s *userImpl) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
	return service.UserCache.GetByIdWithCache(ctx, id)
}

// 更新用户（清除缓存）
func (s *userImpl) Update(ctx context.Context, id uint, data *do.User) error {
	// 1. 更新数据库
	err := dao.User.Update(ctx, id, data)
	if err != nil {
		return err
	}
	
	// 2. 使缓存失效
	service.UserCache.InvalidateUserCache(ctx, id)
	
	return nil
}
```

## 总结

### ✅ 已实现的功能

1. **Redis配置** - 在config.yaml中配置
2. **Cache Service** - 完整的Redis操作封装
3. **用户缓存服务** - 带缓存的用户查询示例
4. **缓存失效** - 更新时自动清除缓存

### 📝 使用方式

```go
// 基础操作
service.Cache.Set(ctx, "key", "value", 1*time.Hour)
service.Cache.Get(ctx, "key")

// 带缓存的用户查询
user, err := service.UserCache.GetByIdWithCache(ctx, 1)

// 缓存失效
service.UserCache.InvalidateUserCache(ctx, 1)
```

### 🎯 关键特性

- ✅ **自动序列化** - 对象自动JSON序列化/反序列化
- ✅ **过期时间** - 支持设置过期时间
- ✅ **模式删除** - 支持按模式批量删除
- ✅ **哈希操作** - 支持Redis哈希操作
- ✅ **计数器** - 支持计数操作

所有代码已实现并通过lint检查，可以直接使用！

