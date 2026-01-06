# Panic Recovery 最佳实践

## 概述

本文档说明在业务代码中，哪些位置需要手动添加 Recovery 捕获 panic，哪些地方不需要。

## 基本原则

### 1. **框架层已提供 Recovery，业务代码通常不需要手动处理**

Kratos 框架在 HTTP/gRPC 服务器层已经提供了 Recovery 中间件，会自动捕获请求处理过程中的 panic，并返回 500 错误给客户端。

### 2. **Panic 用于不可恢复的错误，应该让程序崩溃**

在大多数情况下，panic 表示程序遇到了不可恢复的错误（如编程错误、内存溢出等），应该让程序崩溃并记录日志，而不是捕获并继续运行。

### 3. **手动 Recovery 的适用场景**

只在以下特定场景需要手动添加 Recovery：
- **独立的 goroutine**（不会被框架 Recovery 覆盖）
- **后台任务**（如定时任务、消息队列消费）
- **关键业务逻辑**（需要优雅降级）

## 需要手动添加 Recovery 的场景

### 1. ✅ 独立的 Goroutine（后台任务）

**场景**：在独立的 goroutine 中运行的后台任务，不受框架 Recovery 保护。

**位置**：
- `cmd/daemon-worker/app.go` - Daemon Job 启动
- `cmd/cron-worker/app.go` - Cron Job 启动
- `internal/app/worker/manager.go` - Worker 启动

**示例**：

```go
// ✅ 需要 Recovery
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Errorf("daemon job %s panic: %v, stack: %s", job.Name(), r, debug.Stack())
            // 可以发送告警、记录指标等
        }
    }()
    
    if err := job.Run(ctx); err != nil {
        log.Errorf("daemon job %s stopped with error: %v", job.Name(), err)
    }
}()
```

**原因**：独立的 goroutine 不受框架 Recovery 保护，如果 panic 会导致整个 goroutine 退出，但不会影响主程序。

### 2. ✅ 消息队列消费者循环

**场景**：长时间运行的消息消费循环，需要处理单个消息的 panic，避免整个消费者停止。

**位置**：
- `internal/data/client_rabbitmq.go` - RabbitMQ 消费者
- 类似的消息队列消费逻辑

**示例**：

```go
// ✅ 需要 Recovery（在消息处理循环中）
func (c *RabbitMQConsumer) Consume(ctx context.Context, handler func(ctx context.Context, body []byte) error) error {
    msgs, err := c.ch.Consume(c.queue, "", false, false, false, false, nil)
    if err != nil {
        return err
    }
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case msg, ok := <-msgs:
            if !ok {
                return nil
            }
            
            // ✅ 在每个消息处理中捕获 panic
            func() {
                defer func() {
                    if r := recover(); r != nil {
                        c.logger.Errorf("panic processing message: %v, stack: %s", r, debug.Stack())
                        // 拒绝消息并重新入队
                        msg.Nack(false, true)
                    }
                }()
                
                if err := handler(ctx, msg.Body); err != nil {
                    msg.Nack(false, true)
                    return
                }
                msg.Ack(false)
            }()
        }
    }
}
```

**原因**：单个消息处理失败不应该导致整个消费者停止，需要捕获 panic 并拒绝消息。

### 3. ✅ 定时任务（Cron Jobs）

**场景**：定时执行的任务，单个任务失败不应该影响后续任务。

**位置**：
- `internal/biz/cron/jobs/outbox_dispatcher.go` - Outbox 分发任务
- 其他 Cron Job 实现

**示例**：

```go
// ✅ 需要在 Run 方法中添加 Recovery
func (j *OutboxDispatchJob) Run(ctx context.Context) error {
    defer func() {
        if r := recover(); r != nil {
            j.logger.Errorf("outbox dispatch job panic: %v, stack: %s", r, debug.Stack())
            // 可以发送告警
        }
    }()
    
    // 业务逻辑
    events, err := j.outboxRepo.ListPending(ctx, batchSize)
    // ...
}
```

**原因**：定时任务通常是独立的，单个任务 panic 不应该影响后续任务的执行。

### 4. ✅ 异步操作（Fire-and-Forget）

**场景**：启动异步任务但不等待结果，需要确保 panic 不会影响主流程。

**位置**：
- `internal/data/repo_user.go` - 延迟删除缓存

**示例**：

```go
// ✅ 需要 Recovery
go func() {
    defer func() {
        if r := recover(); r != nil {
            r.log.Errorf("panic in delayed cache delete: %v, stack: %s", r, debug.Stack())
        }
    }()
    
    time.Sleep(userCacheDeleteDelay)
    r.cacheRepo.Delete(context.Background(), cacheKey)
}()
```

**原因**：异步操作的 panic 不会影响主流程，但需要记录日志避免静默失败。

## 不需要手动添加 Recovery 的场景

### 1. ❌ HTTP/gRPC 请求处理

**位置**：
- `internal/service/*.go` - Service 层
- `internal/biz/*.go` - Business 层
- `internal/data/*.go` - Data 层（除了独立的 goroutine）

**原因**：Kratos 框架已经在服务器层提供了 Recovery 中间件（`recovery.Recovery()`），会自动捕获 panic 并返回 500 错误。

**代码位置**：
```go
// internal/server/http.go
globalChain.Add(recovery.Recovery())  // ✅ 框架已提供

// internal/server/grpc.go
grpc.Middleware(recovery.Recovery())  // ✅ 框架已提供
```

### 2. ❌ 业务逻辑函数

**位置**：
- `internal/biz/*.go` - Usecase 方法
- `internal/data/repo_*.go` - Repository 方法（除了独立的 goroutine）

**原因**：业务逻辑函数的 panic 应该被框架 Recovery 捕获，不需要手动处理。

**示例**：

```go
// ❌ 不需要 Recovery
func (uc *UserUsecase) GetUser(ctx context.Context, userID int64) (*User, error) {
    // 不需要 defer recover()
    // 如果 panic，会被框架 Recovery 捕获
    return uc.repo.FindByID(ctx, userID)
}
```

### 3. ❌ 主函数（main）

**位置**：
- `cmd/*/main.go`

**原因**：主函数的 panic 应该让程序崩溃，便于发现启动配置错误。

**示例**：

```go
// ❌ 不需要 Recovery（应该让程序崩溃）
func main() {
    app, cleanup, err := wireApp(...)
    if err != nil {
        panic(err)  // ✅ 正确：启动失败应该 panic
    }
    
    if err := app.Run(); err != nil {
        panic(err)  // ✅ 正确：运行时错误应该 panic
    }
}
```

### 4. ❌ 初始化代码（init、构造函数）

**位置**：
- `New*` 构造函数
- `init()` 函数

**原因**：初始化失败应该让程序无法启动，便于及早发现问题。

**示例**：

```go
// ❌ 不需要 Recovery
func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
    // 如果初始化失败，返回错误即可
    // 不需要 defer recover()
    return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

### 5. ⚠️ 资源清理（特殊场景）

**注意**：`internal/data/repo_order.go` 中的以下代码是**特殊的资源清理模式**：

```go
// ⚠️ 资源清理模式（不是 Recovery）
defer func() {
    if p := recover(); p != nil {
        _ = tx.Rollback()  // 确保事务回滚
        panic(p)           // 重新 panic，让框架 Recovery 处理
    }
}()
```

**原因**：
- 这不是真正的 Recovery，而是**资源清理**模式
- 目的是确保发生 panic 时，事务能够正确回滚
- 然后重新 panic，让框架的 Recovery 中间件处理
- 这是合理的，因为需要保证事务的完整性

**类似场景**：
- 事务回滚
- 文件关闭
- 连接关闭
- 锁释放

**正确做法**：
```go
// ✅ 资源清理 + 重新 panic（合理）
defer func() {
    if p := recover(); p != nil {
        _ = tx.Rollback()  // 清理资源
        panic(p)           // 重新 panic
    }
}()
```

**不推荐的做法**：
```go
// ❌ 不应该静默恢复
defer func() {
    if p := recover(); p != nil {
        _ = tx.Rollback()
        // 不要在这里 return，应该 panic
    }
}()
```

## 最佳实践总结

### ✅ 需要 Recovery 的场景

| 场景 | 位置 | 原因 |
|------|------|------|
| 独立的 goroutine | `go func()` | 不受框架 Recovery 保护 |
| 消息队列消费者 | Consumer 循环 | 单个消息失败不应停止整个消费者 |
| 定时任务 | Cron Job Run 方法 | 单个任务失败不应影响后续任务 |
| 异步操作 | Fire-and-Forget | 不影响主流程但需要记录日志 |

### ❌ 不需要 Recovery 的场景

| 场景 | 位置 | 原因 |
|------|------|------|
| HTTP/gRPC 请求处理 | Service/Biz/Data 层 | 框架已提供 Recovery |
| 业务逻辑函数 | Usecase/Repository 方法 | 框架已提供 Recovery |
| 主函数 | main() | 应该让程序崩溃 |
| 初始化代码 | New* 构造函数 | 失败应该让程序无法启动 |
| 资源清理（事务等） | defer func() { recover(); panic() } | 确保资源清理后重新 panic |

## 实现建议

### 1. 通用 Recovery Helper

可以创建一个通用的 Recovery Helper 函数：

```go
// internal/pkg/recovery/recovery.go
package recovery

import (
    "runtime/debug"
    "github.com/go-kratos/kratos/v2/log"
)

// Recover 捕获 panic 并记录日志
func Recover(logger log.Logger, context string) {
    if r := recover(); r != nil {
        log.NewHelper(logger).Errorf("%s panic: %v, stack: %s", context, r, string(debug.Stack()))
        // 可以发送告警、记录指标等
    }
}

// RecoverWithHandler 捕获 panic 并执行处理函数
func RecoverWithHandler(logger log.Logger, context string, handler func(r interface{})) {
    if r := recover(); r != nil {
        log.NewHelper(logger).Errorf("%s panic: %v, stack: %s", context, r, string(debug.Stack()))
        if handler != nil {
            handler(r)
        }
    }
}
```

### 2. 使用示例

```go
// 在 goroutine 中使用
go func() {
    defer recovery.Recover(logger, "daemon job: "+job.Name())
    job.Run(ctx)
}()

// 在消息处理中使用
func() {
    defer recovery.RecoverWithHandler(logger, "message processing", func(r interface{}) {
        msg.Nack(false, true)  // 拒绝消息
    })
    handler(ctx, msg.Body)
}()
```

### 3. 监控和告警

在 Recovery 中应该：
- ✅ 记录详细的错误日志（包含 stack trace）
- ✅ 发送告警（如发送到监控系统）
- ✅ 记录指标（如 panic 计数器）
- ✅ 优雅降级（如拒绝消息、跳过任务等）

## 参考

- [Kratos Recovery Middleware](https://github.com/go-kratos/kratos/tree/main/middleware/recovery)
- [Go Panic Best Practices](https://go.dev/blog/defer-panic-and-recover)
- 项目中的 Recovery 使用：
  - `internal/server/http.go` - HTTP 服务器 Recovery
  - `internal/server/grpc.go` - gRPC 服务器 Recovery

