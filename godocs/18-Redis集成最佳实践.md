# Redis 集成最佳实践

## GoFrame Redis 支持

GoFrame 框架内置了 Redis 支持，可以直接使用 `g.Redis()` 访问 Redis。

## 配置方式

### 1. 配置文件

在 `manifest/config/config.yaml` 中添加 Redis 配置：

```yaml
# Redis配置
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

### 2. 使用方式

GoFrame 提供了两种使用方式：

#### 方式1: 直接使用 g.Redis()

```go
import "github.com/gogf/gf/v2/frame/g"

// 获取Redis实例
redis := g.Redis()

// 设置值
redis.Set(ctx, "key", "value")

// 获取值
val, err := redis.Get(ctx, "key")
```

#### 方式2: 封装为Service层（推荐）

```go
// internal/service/cache.go
package service

import (
	"context"
	"time"
	"github.com/gogf/gf/v2/frame/g"
)

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
```

## 完整实现示例

### 1. 配置文件

```yaml
# manifest/config/config.yaml
redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    pass: ""
    minIdle: 10
    maxIdle: 100
    maxActive: 200
    idleTimeout: "10s"
    maxConnLifetime: "30s"
```

### 2. Cache Service层

```go
// internal/service/cache.go
package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type ICache interface {
	// Set 设置缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get 获取缓存（字符串）
	Get(ctx context.Context, key string) (string, error)
	// GetObject 获取缓存（对象）
	GetObject(ctx context.Context, key string, dest interface{}) error
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	// Exists 检查key是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error
	// TTL 获取剩余过期时间
	TTL(ctx context.Context, key string) (time.Duration, error)
}

type cacheImpl struct{}

var Cache = cacheImpl{}

// Set 设置缓存
func (s *cacheImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		// 序列化为JSON
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		val = string(data)
	}
	
	_, err := g.Redis().Set(ctx, key, val, expiration)
	return err
}

// Get 获取缓存（字符串）
func (s *cacheImpl) Get(ctx context.Context, key string) (string, error) {
	val, err := g.Redis().Get(ctx, key)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

// GetObject 获取缓存（对象）
func (s *cacheImpl) GetObject(ctx context.Context, key string, dest interface{}) error {
	val, err := g.Redis().Get(ctx, key)
	if err != nil {
		return err
	}
	
	// 反序列化JSON
	return json.Unmarshal([]byte(val.String()), dest)
}

// Delete 删除缓存
func (s *cacheImpl) Delete(ctx context.Context, key string) error {
	_, err := g.Redis().Del(ctx, key)
	return err
}

// Exists 检查key是否存在
func (s *cacheImpl) Exists(ctx context.Context, key string) (bool, error) {
	count, err := g.Redis().Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Expire 设置过期时间
func (s *cacheImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	_, err := g.Redis().Expire(ctx, key, expiration)
	return err
}

// TTL 获取剩余过期时间
func (s *cacheImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := g.Redis().TTL(ctx, key)
	if err != nil {
		return 0, err
	}
	return ttl, nil
}
```

### 3. 在Service层使用Cache

```go
// internal/service/user.go (增强版)
func (s *userImpl) GetById(ctx context.Context, id uint) (user *entity.User, err error) {
	// 1. 先查缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	cached, err := service.Cache.GetObject(ctx, cacheKey, &user)
	if err == nil && user != nil {
		return user, nil
	}
	
	// 2. 缓存未命中，查数据库
	user, err = dao.User.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// 3. 写入缓存
	if user != nil {
		service.Cache.Set(ctx, cacheKey, user, 1*time.Hour)
	}
	
	return user, nil
}
```

## 常用场景示例

### 1. 缓存用户信息

```go
// 设置缓存
err := service.Cache.Set(ctx, "user:1", user, 1*time.Hour)

// 获取缓存
var user entity.User
err := service.Cache.GetObject(ctx, "user:1", &user)
```

### 2. 缓存列表数据

```go
// 设置缓存
err := service.Cache.Set(ctx, "users:list:page:1", users, 10*time.Minute)

// 获取缓存
var users []*entity.User
err := service.Cache.GetObject(ctx, "users:list:page:1", &users)
```

### 3. 分布式锁

```go
// 使用Redis实现分布式锁
func (s *cacheImpl) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	result, err := g.Redis().DoVar(ctx, "SET", key, "1", "NX", "EX", int(expiration.Seconds()))
	if err != nil {
		return false, err
	}
	return result.String() == "OK", nil
}

func (s *cacheImpl) Unlock(ctx context.Context, key string) error {
	_, err := g.Redis().Del(ctx, key)
	return err
}
```

### 4. 计数器

```go
// 增加计数
count, err := g.Redis().Incr(ctx, "counter:key")

// 减少计数
count, err := g.Redis().Decr(ctx, "counter:key")
```

## 最佳实践

### 1. 缓存键命名规范

```go
// 推荐：使用冒号分隔，层次清晰
"user:1"              // 用户信息
"user:1:profile"      // 用户资料
"users:list:page:1"   // 用户列表
"order:123:status"    // 订单状态
```

### 2. 缓存过期时间

```go
// 根据数据特性设置过期时间
1 * time.Hour         // 用户信息：1小时
10 * time.Minute      // 列表数据：10分钟
5 * time.Minute       // 统计数据：5分钟
24 * time.Hour        // 配置数据：24小时
```

### 3. 缓存更新策略

```go
// 策略1: 先更新数据库，再删除缓存
func (s *userImpl) Update(ctx context.Context, id uint, data *do.User) error {
	// 1. 更新数据库
	err := dao.User.Update(ctx, id, data)
	if err != nil {
		return err
	}
	
	// 2. 删除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	service.Cache.Delete(ctx, cacheKey)
	
	return nil
}

// 策略2: 先更新数据库，再更新缓存
func (s *userImpl) Update(ctx context.Context, id uint, data *do.User) error {
	// 1. 更新数据库
	err := dao.User.Update(ctx, id, data)
	if err != nil {
		return err
	}
	
	// 2. 重新查询并更新缓存
	user, err := dao.User.GetById(ctx, id)
	if err == nil && user != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		service.Cache.Set(ctx, cacheKey, user, 1*time.Hour)
	}
	
	return nil
}
```

## 总结

### GoFrame Redis使用

1. **配置** → 在config.yaml中配置
2. **使用** → 通过 `g.Redis()` 访问
3. **封装** → 在Service层封装，便于使用

### 关键点

- ✅ 配置简单，框架内置支持
- ✅ 直接使用 `g.Redis()` 或封装为Service
- ✅ 支持连接池管理
- ✅ 支持所有Redis命令

