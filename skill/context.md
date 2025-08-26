
## 概述

`context.Context` 是 Go 语言中用于跨 API 边界和进程间传递截止时间、取消信号和其他请求范围值的标准方式。它不是同步原语，而是一种协调机制，用于在多个 goroutine 之间传递控制信息。

## 核心价值

1. **取消控制** - 支持优雅取消和超时控制
2. **值传递** - 在调用链中传递请求范围的值
3. **截止时间** - 设置操作的截止时间
4. **请求追踪** - 支持分布式追踪和日志记录

## Context 接口

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

### 方法说明

- `Deadline()` - 返回 Context 的截止时间，如果没有设置则返回 `ok=false`
- `Done()` - 返回一个只读的 channel，当 Context 被取消或超时时会关闭
- `Err()` - 返回 Context 被取消的原因，如果 Context 未被取消则返回 `nil`
- `Value(key)` - 返回与 key 关联的值，如果 key 不存在则返回 `nil`

## 基础 Context 类型

### 1. 空 Context

#### context.Background()
- 根 Context，永不取消
- 通常用作最顶层的 Context
- 适用于 main 函数、初始化、测试等场景

#### context.TODO()
- 当不确定使用哪个 Context 时使用
- 与 Background() 功能相同，但语义更明确
- 适用于临时占位或重构过程中的过渡

### 2. 可取消 Context

#### context.WithCancel(parent Context)
- 基于父 Context 创建可取消的子 Context
- 返回 Context 和 cancel 函数
- cancel 函数可以手动取消 Context

**使用场景：**
- 需要手动控制取消的长时间操作
- 用户主动取消请求
- 资源清理

### 3. 带超时 Context

#### context.WithTimeout(parent Context, timeout time.Duration)
- 创建带超时的 Context
- 超时时间到达后自动取消
- 内部使用 `WithDeadline` 实现

**使用场景：**
- API 调用超时控制
- 数据库查询超时
- 网络请求超时

### 4. 带截止时间 Context

#### context.WithDeadline(parent Context, deadline time.Time)
- 创建带截止时间的 Context
- 到达截止时间后自动取消
- 比 WithTimeout 更精确的时间控制

**使用场景：**
- 精确的时间控制
- 基于绝对时间的操作
- 定时任务

### 5. 带值 Context

#### context.WithValue(parent Context, key, val interface{})
- 在 Context 中存储键值对
- 子 Context 会继承父 Context 的值
- 同名的 key 会被覆盖

**使用场景：**
- 传递请求 ID
- 传递用户信息
- 传递追踪信息

## Context 链式传递

### 特点

1. **继承性** - 子 Context 继承父 Context 的所有特性
2. **覆盖性** - 子 Context 可以覆盖父 Context 的某些特性
3. **独立性** - 子 Context 的取消不会影响父 Context

### 示例

```go
// 创建 Context 链
rootCtx := context.Background()
timeoutCtx, cancel1 := context.WithTimeout(rootCtx, 5*time.Second)
valueCtx := context.WithValue(timeoutCtx, "user_id", "12345")
cancelCtx, cancel2 := context.WithCancel(valueCtx)
```

## 实际应用场景

### 1. HTTP 请求处理

**场景描述：**
- 处理 HTTP 请求时需要超时控制
- 在请求处理过程中传递用户信息
- 支持请求取消

**实现要点：**
- 使用 `WithTimeout` 设置请求超时
- 使用 `WithValue` 传递用户信息
- 在数据库查询和 API 调用中检查 Context

### 2. 并发任务控制

**场景描述：**
- 启动多个并发任务
- 需要统一取消所有任务
- 收集任务执行结果

**实现要点：**
- 使用 `WithCancel` 创建可取消的 Context
- 在任务中定期检查 Context 状态
- 使用 channel 收集错误信息

### 3. 资源池管理

**场景描述：**
- 管理有限数量的资源
- 支持资源获取超时
- 支持资源池关闭

**实现要点：**
- 使用 `WithTimeout` 控制资源获取超时
- 使用 `WithCancel` 支持资源池关闭
- 在资源操作中检查 Context 状态

## 最佳实践

### 1. Context 传递原则

#### 总是将 Context 作为第一个参数
```go
// 正确
func Process(ctx context.Context, data []byte) error

// 错误
func Process(data []byte, ctx context.Context) error
```

#### 不要将 Context 存储在结构体中
```go
// 错误示例
type BadStruct struct {
    ctx context.Context
}

// 正确做法
func (s *GoodStruct) Process(ctx context.Context) error
```

#### 总是调用 cancel 函数
```go
ctx, cancel := context.WithCancel(parentCtx)
defer cancel() // 确保资源被释放
```

### 2. Context 值传递最佳实践

#### 使用类型安全的键
```go
type contextKey string

const (
    UserIDKey    contextKey = "user_id"
    RequestIDKey contextKey = "request_id"
)
```

#### 安全地获取值
```go
if userID, ok := ctx.Value(UserIDKey).(string); ok {
    // 使用 userID
}
```

#### 避免传递可变数据
```go
// 错误：传递可变数据
ctx = context.WithValue(ctx, "shared_map", make(map[string]interface{}))

// 正确：传递不可变数据
ctx = context.WithValue(ctx, "user_id", "12345")
```

### 3. Context 取消最佳实践

#### 检查取消原因
```go
select {
case <-ctx.Done():
    switch ctx.Err() {
    case context.Canceled:
        // 手动取消
    case context.DeadlineExceeded:
        // 超时取消
    }
}
```

#### 在长时间操作中定期检查
```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // 执行工作
        time.Sleep(time.Millisecond * 100)
    }
}
```

## 与其他同步原语的结合

### 1. Context 与 Channel

**结合方式：**
- 使用 Context 控制 channel 操作
- 在 select 语句中同时监听 Context 和 channel

**示例：**
```go
select {
case <-ctx.Done():
    return ctx.Err()
case data := <-dataChan:
    // 处理数据
}
```

### 2. Context 与 WaitGroup

**结合方式：**
- 使用 Context 控制 goroutine 的生命周期
- 使用 WaitGroup 等待所有 goroutine 完成

**示例：**
```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            select {
            case <-ctx.Done():
                return
            default:
                // 执行工作
            }
        }
    }()
}
wg.Wait()
```

## 性能考虑

### 1. Context 链的性能影响

**问题：**
- Context 链越长，值查找越慢
- 每次 `Value()` 调用都需要遍历整个链

**优化建议：**
- 避免创建过长的 Context 链
- 将频繁访问的值放在较浅的位置
- 考虑使用其他方式传递大量数据

### 2. 避免 Context 值滥用

**问题：**
- 在 Context 中存储大量数据
- 频繁创建和销毁 Context

**优化建议：**
- 只存储必要的标识符
- 重用 Context 对象
- 使用对象池管理 Context

## 常见错误和陷阱

### 1. 忘记调用 cancel 函数

**错误示例：**
```go
ctx, cancel := context.WithCancel(parentCtx)
// 忘记调用 cancel()
```

**正确做法：**
```go
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()
```

### 2. 在 Context 中存储可变数据

**错误示例：**
```go
ctx = context.WithValue(ctx, "map", make(map[string]interface{}))
```

**正确做法：**
```go
ctx = context.WithValue(ctx, "user_id", "12345")
```

### 3. 忽略 Context 取消

**错误示例：**
```go
func process(ctx context.Context) {
    // 长时间操作，不检查 Context
    time.Sleep(time.Hour)
}
```

**正确做法：**
```go
func process(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 执行工作
        }
    }
}
```

## 总结

Context 是 Go 并发编程中非常重要的工具，它提供了一种优雅的方式来处理取消、超时和值传递。正确使用 Context 可以：

1. **提高代码的可维护性** - 统一的取消机制
2. **增强系统的健壮性** - 避免资源泄漏
3. **改善用户体验** - 支持请求取消和超时
4. **简化错误处理** - 统一的错误传播机制

掌握 Context 的使用对于编写高质量的 Go 并发程序至关重要。
# Go Context 详解

