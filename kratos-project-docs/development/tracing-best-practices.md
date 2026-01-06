# OpenTelemetry Tracing 最佳实践

## 概述

本文档介绍在项目中为数据库操作、第三方服务调用等关键操作添加 OpenTelemetry tracing 的最佳实践。通过添加 tracing，我们可以：

- **追踪性能瓶颈**：了解每个操作的耗时
- **分布式追踪**：将数据库、Redis、第三方服务调用纳入完整的请求追踪链路
- **问题排查**：快速定位是哪个操作出现问题
- **性能分析**：分析系统性能瓶颈

## 目录

1. [数据库操作 Tracing（GORM）](#数据库操作-tracinggorm)
2. [Redis 操作 Tracing](#redis-操作-tracing)
3. [第三方服务调用 Tracing](#第三方服务调用-tracing)
4. [通用工具函数](#通用工具函数)
5. [最佳实践总结](#最佳实践总结)

---

## 数据库操作 Tracing（GORM）

### 基本用法

在 Repository 方法中为数据库操作创建 span：

```go
package data

import (
    "context"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    
    "gorm.io/gorm"
)

type userRepo struct {
    data   *Data
    log    *log.Helper
    tracer trace.Tracer // 添加 tracer
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
    return &userRepo{
        data:   data,
        log:    log.NewHelper(logger),
        tracer: otel.Tracer("user-repo"), // 初始化 tracer
    }
}

// Save 保存用户
func (r *userRepo) Save(ctx context.Context, user *biz.User) (*biz.User, error) {
    // 创建 span
    ctx, span := r.tracer.Start(ctx, "user-repo.Save",
        trace.WithAttributes(
            attribute.String("db.operation", "INSERT"),
            attribute.String("db.table", "users"),
            attribute.String("db.user.username", user.Username),
            attribute.String("db.user.email", user.Email),
        ),
    )
    defer span.End()
    
    start := time.Now()
    
    if err := r.checkDB(); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "database not available")
        return nil, err
    }
    
    u := userFromBiz(user)
    
    // 执行数据库操作（使用包含 span 的 context）
    if err := r.data.db.WithContext(ctx).Create(u).Error; err != nil {
        duration := time.Since(start)
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("db.duration_ms", duration.Milliseconds()),
            attribute.String("db.error", err.Error()),
        )
        r.log.WithContext(ctx).Errorf("failed to save user: %v", err)
        return nil, errors.InternalServer(v1.ErrorReason_USER_SAVE_FAILED.String(), "failed to save user")
    }
    
    duration := time.Since(start)
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int64("db.duration_ms", duration.Milliseconds()),
        attribute.Int64("db.user.id", u.ID),
    )
    
    return userToBiz(u), nil
}

// FindByID 根据 ID 查询用户
func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
    ctx, span := r.tracer.Start(ctx, "user-repo.FindByID",
        trace.WithAttributes(
            attribute.String("db.operation", "SELECT"),
            attribute.String("db.table", "users"),
            attribute.Int64("db.user.id", id),
        ),
    )
    defer span.End()
    
    start := time.Now()
    
    if err := r.checkDB(); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "database not available")
        return nil, err
    }
    
    var u model.User
    if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
        duration := time.Since(start)
        if err == gorm.ErrRecordNotFound {
            span.SetStatus(codes.Ok, "not found") // 未找到不算错误
            span.SetAttributes(
                attribute.Int64("db.duration_ms", duration.Milliseconds()),
                attribute.Bool("db.not_found", true),
            )
            return nil, errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
        }
        
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("db.duration_ms", duration.Milliseconds()),
            attribute.String("db.error", err.Error()),
        )
        r.log.WithContext(ctx).Errorf("failed to find user by id: %v", err)
        return nil, errors.InternalServer(v1.ErrorReason_USER_QUERY_FAILED.String(), "failed to query user")
    }
    
    duration := time.Since(start)
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int64("db.duration_ms", duration.Milliseconds()),
    )
    
    return userToBiz(&u), nil
}

// List 分页查询用户列表
func (r *userRepo) List(ctx context.Context, page, pageSize int64, keyword string) ([]*biz.User, int64, error) {
    ctx, span := r.tracer.Start(ctx, "user-repo.List",
        trace.WithAttributes(
            attribute.String("db.operation", "SELECT"),
            attribute.String("db.table", "users"),
            attribute.Int64("db.query.page", page),
            attribute.Int64("db.query.page_size", pageSize),
            attribute.String("db.query.keyword", keyword),
        ),
    )
    defer span.End()
    
    start := time.Now()
    
    if err := r.checkDB(); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "database not available")
        return nil, 0, err
    }
    
    var users []model.User
    var total int64
    query := r.data.db.WithContext(ctx).Model(&model.User{})
    
    if keyword != "" {
        query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
            "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
    }
    
    // 计算总数
    if err := query.Count(&total).Error; err != nil {
        duration := time.Since(start)
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("db.duration_ms", duration.Milliseconds()),
        )
        return nil, 0, errors.InternalServer(v1.ErrorReason_USER_QUERY_FAILED.String(), "failed to count users")
    }
    
    // 分页查询
    offset := (page - 1) * pageSize
    if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("id DESC").Find(&users).Error; err != nil {
        duration := time.Since(start)
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("db.duration_ms", duration.Milliseconds()),
        )
        return nil, 0, errors.InternalServer(v1.ErrorReason_USER_QUERY_FAILED.String(), "failed to list users")
    }
    
    duration := time.Since(start)
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int64("db.duration_ms", duration.Milliseconds()),
        attribute.Int64("db.result.count", int64(len(users))),
        attribute.Int64("db.result.total", total),
    )
    
    result := make([]*biz.User, 0, len(users))
    for i := range users {
        result = append(result, userToBiz(&users[i]))
    }
    
    return result, total, nil
}
```

### 事务操作

对于事务操作，可以创建嵌套的 span：

```go
func (r *userRepo) CreateWithTransaction(ctx context.Context, user *biz.User) (*biz.User, error) {
    ctx, span := r.tracer.Start(ctx, "user-repo.CreateWithTransaction",
        trace.WithAttributes(
            attribute.String("db.operation", "TRANSACTION"),
            attribute.String("db.table", "users"),
        ),
    )
    defer span.End()
    
    // 在事务中执行操作
    err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 创建子 span
        ctx, txSpan := r.tracer.Start(ctx, "user-repo.CreateInTransaction",
            trace.WithAttributes(
                attribute.String("db.operation", "INSERT"),
                attribute.String("db.table", "users"),
            ),
        )
        defer txSpan.End()
        
        u := userFromBiz(user)
        if err := tx.WithContext(ctx).Create(u).Error; err != nil {
            txSpan.RecordError(err)
            txSpan.SetStatus(codes.Error, err.Error())
            return err
        }
        
        txSpan.SetStatus(codes.Ok, "success")
        return nil
    })
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    span.SetStatus(codes.Ok, "success")
    return user, nil
}
```

### 数据库操作属性规范

| 属性名 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `db.operation` | string | 操作类型 | `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `TRANSACTION` |
| `db.table` | string | 表名 | `users`, `orders` |
| `db.duration_ms` | int64 | 操作耗时（毫秒） | `150` |
| `db.user.id` | int64 | 用户ID（示例） | `123` |
| `db.query.page` | int64 | 分页页码 | `1` |
| `db.query.page_size` | int64 | 分页大小 | `20` |
| `db.result.count` | int64 | 返回结果数量 | `10` |
| `db.result.total` | int64 | 总记录数 | `100` |
| `db.error` | string | 错误信息 | `connection timeout` |
| `db.not_found` | bool | 是否未找到记录 | `true` |

---

## Redis 操作 Tracing

### 基本用法

为 Redis 操作创建 span：

```go
package data

import (
    "context"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    
    "github.com/redis/go-redis/v9"
)

type cacheRepo struct {
    redis  *redis.Client
    log    *log.Helper
    tracer trace.Tracer
}

func NewCacheRepo(redisClient *redis.Client, logger log.Logger) biz.CacheRepo {
    if redisClient == nil {
        return nil
    }
    return &cacheRepo{
        redis:  redisClient,
        log:    log.NewHelper(logger),
        tracer: otel.Tracer("cache-repo"),
    }
}

// Get 获取缓存
func (r *cacheRepo) Get(ctx context.Context, key string) (string, error) {
    ctx, span := r.tracer.Start(ctx, "cache-repo.Get",
        trace.WithAttributes(
            attribute.String("cache.operation", "GET"),
            attribute.String("cache.key", key),
        ),
    )
    defer span.End()
    
    start := time.Now()
    
    val, err := r.redis.Get(ctx, key).Result()
    duration := time.Since(start)
    
    if err != nil {
        if err == redis.Nil {
            span.SetStatus(codes.Ok, "not found") // 缓存未命中不算错误
            span.SetAttributes(
                attribute.Int64("cache.duration_ms", duration.Milliseconds()),
                attribute.Bool("cache.miss", true),
            )
            return "", nil
        }
        
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("cache.duration_ms", duration.Milliseconds()),
            attribute.String("cache.error", err.Error()),
        )
        return "", err
    }
    
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int64("cache.duration_ms", duration.Milliseconds()),
        attribute.Bool("cache.hit", true),
        attribute.Int("cache.value_size", len(val)),
    )
    
    return val, nil
}

// Set 设置缓存
func (r *cacheRepo) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
    ctx, span := r.tracer.Start(ctx, "cache-repo.Set",
        trace.WithAttributes(
            attribute.String("cache.operation", "SET"),
            attribute.String("cache.key", key),
            attribute.Int64("cache.expiration_seconds", int64(expiration.Seconds())),
            attribute.Int("cache.value_size", len(value)),
        ),
    )
    defer span.End()
    
    start := time.Now()
    
    err := r.redis.Set(ctx, key, value, expiration).Err()
    duration := time.Since(start)
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("cache.duration_ms", duration.Milliseconds()),
            attribute.String("cache.error", err.Error()),
        )
        return err
    }
    
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int64("cache.duration_ms", duration.Milliseconds()),
    )
    
    return nil
}

// Delete 删除缓存
func (r *cacheRepo) Delete(ctx context.Context, key string) error {
    ctx, span := r.tracer.Start(ctx, "cache-repo.Delete",
        trace.WithAttributes(
            attribute.String("cache.operation", "DELETE"),
            attribute.String("cache.key", key),
        ),
    )
    defer span.End()
    
    start := time.Now()
    
    err := r.redis.Del(ctx, key).Err()
    duration := time.Since(start)
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("cache.duration_ms", duration.Milliseconds()),
        )
        return err
    }
    
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int64("cache.duration_ms", duration.Milliseconds()),
    )
    
    return nil
}
```

### Redis 操作属性规范

| 属性名 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `cache.operation` | string | 操作类型 | `GET`, `SET`, `DELETE`, `EXISTS`, `INCR` |
| `cache.key` | string | 缓存键 | `user:123` |
| `cache.duration_ms` | int64 | 操作耗时（毫秒） | `5` |
| `cache.hit` | bool | 缓存命中 | `true` |
| `cache.miss` | bool | 缓存未命中 | `true` |
| `cache.value_size` | int | 值大小（字节） | `1024` |
| `cache.expiration_seconds` | int64 | 过期时间（秒） | `3600` |
| `cache.error` | string | 错误信息 | `connection timeout` |

---

## 第三方服务调用 Tracing

### HTTP 客户端（已有示例）

参考 `internal/data/external/dingtalk/client.go` 中的实现，或查看文档：
- [第三方服务调用 Tracing 文档](./opentelemetry-tracing-third-party.md)

### gRPC 客户端

```go
package payment

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    
    "google.golang.org/grpc"
    pb "api/external/payment/v1"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.PaymentServiceClient
    tracer trace.Tracer
}

func NewClient(conn *grpc.ClientConn, logger log.Logger) *Client {
    return &Client{
        conn:   conn,
        client: pb.NewPaymentServiceClient(conn),
        tracer: otel.Tracer("payment-client"),
    }
}

func (c *Client) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    ctx, span := c.tracer.Start(ctx, "payment.CreateOrder",
        trace.WithAttributes(
            attribute.String("rpc.service", "payment.PaymentService"),
            attribute.String("rpc.method", "CreateOrder"),
            attribute.String("payment.user_id", req.UserId),
            attribute.Float64("payment.amount", float64(req.Amount)),
        ),
    )
    defer span.End()
    
    resp, err := c.client.CreateOrder(ctx, req)
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.String("payment.order_id", resp.OrderId),
        attribute.String("rpc.status_code", "OK"),
    )
    
    return resp, nil
}
```

### gRPC 操作属性规范

| 属性名 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `rpc.service` | string | gRPC 服务名 | `payment.PaymentService` |
| `rpc.method` | string | gRPC 方法名 | `CreateOrder` |
| `rpc.status_code` | string | gRPC 状态码 | `OK`, `NOT_FOUND`, `INTERNAL` |
| `payment.user_id` | string | 业务相关属性（示例） | `user123` |
| `payment.amount` | float64 | 业务相关属性（示例） | `99.99` |

---

## 通用工具函数

为了减少重复代码，可以创建通用的 tracing 工具函数：

```go
package tracing

import (
    "context"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

// TraceDBOperation 追踪数据库操作
func TraceDBOperation(ctx context.Context, operation, table string, fn func(context.Context) error) error {
    tracer := otel.Tracer("db")
    ctx, span := tracer.Start(ctx, "db."+operation,
        trace.WithAttributes(
            attribute.String("db.operation", operation),
            attribute.String("db.table", table),
        ),
    )
    defer span.End()
    
    start := time.Now()
    err := fn(ctx)
    duration := time.Since(start)
    
    span.SetAttributes(
        attribute.Int64("db.duration_ms", duration.Milliseconds()),
    )
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    } else {
        span.SetStatus(codes.Ok, "success")
    }
    
    return err
}

// TraceCacheOperation 追踪缓存操作
func TraceCacheOperation(ctx context.Context, operation, key string, fn func(context.Context) error) error {
    tracer := otel.Tracer("cache")
    ctx, span := tracer.Start(ctx, "cache."+operation,
        trace.WithAttributes(
            attribute.String("cache.operation", operation),
            attribute.String("cache.key", key),
        ),
    )
    defer span.End()
    
    start := time.Now()
    err := fn(ctx)
    duration := time.Since(start)
    
    span.SetAttributes(
        attribute.Int64("cache.duration_ms", duration.Milliseconds()),
    )
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    } else {
        span.SetStatus(codes.Ok, "success")
    }
    
    return err
}

// 使用示例
func (r *userRepo) Save(ctx context.Context, user *biz.User) (*biz.User, error) {
    var u *model.User
    err := TraceDBOperation(ctx, "INSERT", "users", func(ctx context.Context) error {
        u = userFromBiz(user)
        return r.data.db.WithContext(ctx).Create(u).Error
    })
    
    if err != nil {
        return nil, err
    }
    
    return userToBiz(u), nil
}
```

---

## 最佳实践总结

### 1. Span 命名规范

- **格式**：`{service-name}.{operation-name}`
- **示例**：
  - `user-repo.Save`
  - `cache-repo.Get`
  - `dingtalk.GetUserInfo`
  - `payment.CreateOrder`

### 2. 属性命名规范

使用统一的前缀：
- **数据库**：`db.*`
- **缓存**：`cache.*`
- **HTTP**：`http.*`
- **gRPC**：`rpc.*`
- **业务相关**：`{service}.*`（如 `dingtalk.*`, `payment.*`）

### 3. 错误处理

始终记录错误：
```go
if err != nil {
    span.RecordError(err)           // 记录错误堆栈
    span.SetStatus(codes.Error, err.Error()) // 设置错误状态
    return nil, err
}
```

### 4. 性能指标

记录关键性能指标：
```go
start := time.Now()
// ... 执行操作
duration := time.Since(start)
span.SetAttributes(
    attribute.Int64("db.duration_ms", duration.Milliseconds()),
)
```

### 5. Context 传播

**重要**：确保将包含 span 的 context 传递给所有子操作：
- ✅ `db.WithContext(ctx)`
- ✅ `redis.Get(ctx, key)`
- ✅ `httpClient.R().SetContext(ctx)`
- ✅ `grpcClient.Method(ctx, req)`

### 6. 何时创建 Span

**应该创建 span 的操作**：
- ✅ 数据库操作（CRUD）
- ✅ 缓存操作（Get/Set/Delete）
- ✅ 第三方服务调用（HTTP/gRPC）
- ✅ 消息队列操作（发送/接收）
- ✅ 文件操作（读写）
- ✅ 复杂业务逻辑（包含多个步骤）

**不需要创建 span 的操作**：
- ❌ 简单的数据转换
- ❌ 简单的计算
- ❌ 内存操作

### 7. 嵌套 Span

对于复杂操作，使用嵌套 span：
```go
// 父 span：整个操作
ctx, span := tracer.Start(ctx, "user-repo.CreateWithTransaction")
defer span.End()

// 子 span：事务中的操作
err := db.Transaction(func(tx *gorm.DB) error {
    ctx, txSpan := tracer.Start(ctx, "user-repo.CreateInTransaction")
    defer txSpan.End()
    // ...
})
```

### 8. 避免过度追踪

不要为每个小操作都创建 span，只在关键操作上使用。过多的 span 会导致：
- 性能开销
- 追踪数据过多，难以分析
- 存储成本增加

---

## 相关文档

- [OpenTelemetry Tracing 集成文档](./kratos-opentelemetry-integration.md)
- [第三方服务调用 Tracing](./opentelemetry-tracing-third-party.md)
- [OpenTelemetry Go 官方文档](https://opentelemetry.io/docs/instrumentation/go/)

