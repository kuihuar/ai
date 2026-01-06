# 使用 gcache 重构缓存

## 为什么使用 gcache？

GoFrame 提供了 `gcache` 包，这是一个统一的缓存抽象层，具有以下优势：

### 1. 统一的缓存接口
- ✅ 支持多种后端：内存缓存、Redis缓存
- ✅ 可以轻松切换缓存后端，无需修改业务代码
- ✅ 提供统一的 API，使用简单

### 2. 高级功能
- ✅ `GetOrSet` - 获取或设置（原子操作）
- ✅ `GetOrSetFunc` - 获取或执行函数设置（防止缓存击穿）
- ✅ `GetOrSetFuncLock` - 带锁的获取或执行函数设置（防止缓存穿透）

### 3. 自动序列化
- ✅ 自动处理对象的序列化和反序列化
- ✅ 支持任意类型的值

## 实现方式

### 1. 初始化缓存（使用Redis适配器）

```go
// internal/service/cache.go
func InitCache() error {
	ctx := gctx.New()

	// 获取Redis配置
	redisConfig, ok := gredis.GetConfig("default")
	if !ok {
		return fmt.Errorf("获取Redis配置失败: 配置不存在")
	}

	// 创建Redis客户端
	redis, err := gredis.New(redisConfig)
	if err != nil {
		return fmt.Errorf("创建Redis客户端失败: %w", err)
	}

	// 创建Redis适配器
	adapter := gcache.NewAdapterRedis(redis)

	// 创建缓存实例
	cacheInstance = gcache.NewWithAdapter(adapter)

	return nil
}
```

### 2. 在启动时初始化

```go
// internal/cmd/cmd.go
Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
	// 初始化缓存（使用Redis适配器）
	if err = service.InitCache(); err != nil {
		return fmt.Errorf("初始化缓存失败: %w", err)
	}
	
	// ... 其他初始化
}
```

## 使用示例

### 基础操作

```go
// 设置缓存
err := service.Cache.Set(ctx, "key", "value", 1*time.Hour)

// 获取缓存
val, err := service.Cache.Get(ctx, "key")

// 获取对象
var user entity.User
err := service.Cache.GetObject(ctx, "user:1", &user)

// 删除缓存
err := service.Cache.Delete(ctx, "key")
```

### 高级功能（gcache 原生支持）

```go
import "github.com/gogf/gf/v2/os/gcache"

cache := service.getCache() // 获取内部缓存实例

// GetOrSet - 获取或设置（原子操作）
val, err := cache.GetOrSet(ctx, "key", "default_value", 1*time.Hour)

// GetOrSetFunc - 获取或执行函数设置（防止缓存击穿）
val, err := cache.GetOrSetFunc(ctx, "user:1", func(ctx context.Context) (interface{}, error) {
	// 如果缓存不存在，执行这个函数
	user, err := dao.User.GetById(ctx, 1)
	if err != nil {
		return nil, err
	}
	return user, nil
}, 1*time.Hour)

// GetOrSetFuncLock - 带锁的获取或执行函数设置（防止缓存穿透）
val, err := cache.GetOrSetFuncLock(ctx, "user:1", func(ctx context.Context) (interface{}, error) {
	// 带锁，防止并发时重复执行
	user, err := dao.User.GetById(ctx, 1)
	if err != nil {
		return nil, err
	}
	return user, nil
}, 1*time.Hour)
```

## 与直接使用 Redis 的对比

### 直接使用 g.Redis()

```go
// 需要手动处理序列化
data, _ := json.Marshal(user)
g.Redis().Set(ctx, "user:1", data, expiration)

// 需要手动处理反序列化
val, _ := g.Redis().Get(ctx, "user:1")
json.Unmarshal([]byte(val.String()), &user)
```

### 使用 gcache

```go
// 自动处理序列化
cache.Set(ctx, "user:1", user, expiration)

// 自动处理反序列化
cache.Get(ctx, "user:1", &user)
```

## 优势总结

### ✅ 代码更简洁
- 自动序列化/反序列化
- 统一的 API

### ✅ 功能更强大
- `GetOrSet` - 原子操作
- `GetOrSetFunc` - 防止缓存击穿
- `GetOrSetFuncLock` - 防止缓存穿透

### ✅ 更灵活
- 可以切换缓存后端（内存/Redis）
- 统一的接口，易于测试

### ✅ 更安全
- 防止缓存击穿（GetOrSetFunc）
- 防止缓存穿透（GetOrSetFuncLock）

## 注意事项

### 1. 哈希操作
`gcache` 不直接支持 Redis 的哈希操作（HSet、HGet等），我们使用复合key来模拟：

```go
// HSet 使用复合key
hashKey := fmt.Sprintf("%s:%s", key, field)
cache.Set(ctx, hashKey, value, 0)

// HGet 使用复合key
hashKey := fmt.Sprintf("%s:%s", key, field)
val, err := cache.Get(ctx, hashKey)
```

### 2. 计数器操作
`gcache` 不直接支持 Incr/Decr，我们手动实现：

```go
// 获取当前值
val, err := cache.Get(ctx, key)
var count int64
if err != nil || val == nil {
	count = 1  // 初始化为1
} else {
	count = val.Int64() + 1  // 加1
}
// 设置新值
cache.Set(ctx, key, count, 0)
```

### 3. 模式删除
`gcache` 不直接支持模式删除，我们通过遍历所有key来实现：

```go
// 获取所有key
keys, err := cache.Keys(ctx)
// 过滤匹配的key
matchedKeys := filterKeys(keys, pattern)
// 批量删除
cache.Removes(ctx, matchedKeys)
```

## 参考文档

- [GoFrame gcache 文档](https://pkg.go.dev/github.com/gogf/gf/v2/os/gcache)
- [GoFrame Redis 文档](https://pkg.go.dev/github.com/gogf/gf/v2/database/gredis)

