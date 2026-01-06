# Redis 完整实现示例

## 已实现的Redis功能

### 1. 配置文件
- ✅ 在 `manifest/config/config.yaml` 中添加Redis配置

### 2. Cache Service层
- ✅ `internal/service/cache.go` - 缓存服务封装

### 3. 用户缓存服务
- ✅ `internal/service/user_cache.go` - 带缓存的用户服务

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

## Cache Service层

### 基础操作

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

// 检查是否存在
exists, err := service.Cache.Exists(ctx, "key")
```

### 高级操作

```go
// 按模式删除
err := service.Cache.DeletePattern(ctx, "user:*")

// 设置过期时间
err := service.Cache.Expire(ctx, "key", 1*time.Hour)

// 获取剩余过期时间
ttl, err := service.Cache.TTL(ctx, "key")

// 计数器
count, err := service.Cache.Incr(ctx, "counter:key")
count, err := service.Cache.Decr(ctx, "counter:key")

// 哈希操作
err := service.Cache.HSet(ctx, "user:1", "name", "张三")
name, err := service.Cache.HGet(ctx, "user:1", "name")
all, err := service.Cache.HGetAll(ctx, "user:1")
```

## 使用示例

### 1. 缓存用户信息

```go
// 在Service层使用
func (s *userImpl) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
	// 使用带缓存的服务
	return service.UserCache.GetByIdWithCache(ctx, id)
}
```

### 2. 缓存失效

```go
// 更新用户后，使缓存失效
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

### 3. 缓存列表数据

```go
// 获取用户列表（带缓存）
users, total, err := service.UserCache.GetListWithCache(ctx, 1, 10)
```

## 缓存键命名规范

### 推荐命名方式

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

// 哈希
"user:1"                   // 用户哈希表
  - "name" -> "张三"
  - "email" -> "zhang@example.com"
```

## 缓存策略

### 1. Cache-Aside（旁路缓存）

```go
// 读取：先查缓存，未命中查数据库，再写入缓存
func GetById(ctx context.Context, id uint) {
	// 1. 查缓存
	user := cache.Get("user:1")
	if user != nil {
		return user
	}
	
	// 2. 查数据库
	user = db.GetById(id)
	
	// 3. 写入缓存
	cache.Set("user:1", user)
	
	return user
}
```

### 2. Write-Through（写穿透）

```go
// 写入：同时更新数据库和缓存
func Update(ctx context.Context, id uint, data *do.User) {
	// 1. 更新数据库
	db.Update(id, data)
	
	// 2. 更新缓存
	cache.Set("user:1", data)
}
```

### 3. Write-Back（写回）

```go
// 写入：先写缓存，异步写数据库
func Update(ctx context.Context, id uint, data *do.User) {
	// 1. 更新缓存
	cache.Set("user:1", data)
	
	// 2. 标记为脏数据
	cache.MarkDirty("user:1")
	
	// 3. 异步批量写入数据库
	go batchWriteToDB()
}
```

## 实际业务场景

### 场景1: 用户信息缓存

```go
// 获取用户（带缓存）
user, err := service.UserCache.GetByIdWithCache(ctx, 1)
// 第一次：查数据库，写入缓存
// 后续：直接从缓存读取
```

### 场景2: 列表数据缓存

```go
// 获取用户列表（带缓存）
users, total, err := service.UserCache.GetListWithCache(ctx, 1, 10)
// 缓存10分钟，减少数据库查询
```

### 场景3: 缓存失效

```go
// 更新用户后，自动清除相关缓存
service.UserCache.InvalidateUserCache(ctx, id)
// 清除：user:1 和 users:list:*
```

## 性能优化建议

### 1. 缓存预热

```go
// 应用启动时预热热点数据
func warmupCache(ctx context.Context) {
	// 预热热门用户
	hotUserIds := []uint{1, 2, 3}
	for _, id := range hotUserIds {
		service.UserCache.GetByIdWithCache(ctx, id)
	}
}
```

### 2. 缓存穿透防护

```go
// 使用布隆过滤器或缓存空值
func GetById(ctx context.Context, id uint) {
	// 1. 检查是否在黑名单（不存在的ID）
	if isBlacklisted(id) {
		return nil, errors.New("用户不存在")
	}
	
	// 2. 查缓存
	user := cache.Get("user:1")
	if user != nil {
		return user, nil
	}
	
	// 3. 查数据库
	user = db.GetById(id)
	if user == nil {
		// 缓存空值，防止穿透
		cache.Set("user:1:null", "", 5*time.Minute)
		return nil, errors.New("用户不存在")
	}
	
	// 4. 写入缓存
	cache.Set("user:1", user)
	return user, nil
}
```

### 3. 缓存雪崩防护

```go
// 设置随机过期时间
func SetWithRandomExpire(ctx context.Context, key string, value interface{}, baseExpire time.Duration) {
	// 基础过期时间 + 随机0-10分钟
	randomOffset := time.Duration(rand.Intn(600)) * time.Second
	expire := baseExpire + randomOffset
	Cache.Set(ctx, key, value, expire)
}
```

## 总结

### ✅ 已实现的功能

1. **Redis配置** - 在config.yaml中配置
2. **Cache Service** - 封装Redis操作
3. **用户缓存服务** - 带缓存的用户查询
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

所有代码已实现，可以直接使用！

