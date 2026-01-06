# Go 中 Panic 与协程的关系详解

## 一、核心问题

### 1.1 问题澄清

**问题**: 每个请求是在独立的协程里吗？一个协程的 panic 如果不 recover 会导致整个服务崩溃吗？

**答案**:
1. ✅ **是的**，每个 HTTP 请求确实在独立的 goroutine 中处理
2. ✅ **是的**，如果一个 goroutine 发生 panic 且没有被 recover，**会导致整个程序崩溃**

### 1.2 Go 语言的重要特性

在 Go 语言中，**任何 goroutine 的未恢复 panic 都会导致整个程序崩溃**，这是 Go 语言的设计决策。

## 二、HTTP 服务器的请求处理机制

### 2.1 标准库的实现

```go
// net/http/server.go
func (srv *Server) Serve(l net.Listener) error {
    for {
        // 1. 接受连接
        rw, err := l.Accept()
        if err != nil {
            // ...
        }
        
        // 2. 为每个连接创建 goroutine
        go c.serve(connCtx)  // ← 每个请求在独立的 goroutine 中
    }
}
```

**关键点**:
- `go c.serve(connCtx)` - 每个连接在独立的 goroutine 中处理
- 这意味着每个 HTTP 请求都在独立的 goroutine 中处理

### 2.2 请求处理流程

```
客户端请求
  ↓
服务器接受连接
  ↓
创建新的 goroutine (go c.serve)
  ↓
在 goroutine 中处理请求
  ↓
返回响应
  ↓
goroutine 结束
```

## 三、Panic 对程序的影响

### 3.1 未恢复的 Panic 会导致程序崩溃

```go
package main

import (
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // 这个 panic 没有被 recover
    panic("something went wrong")  // ← 会导致整个程序崩溃！
}

func main() {
    http.HandleFunc("/", handler)
    go http.ListenAndServe(":8080", nil)
    
    // 即使主 goroutine 在运行，如果 handler goroutine panic
    // 整个程序也会崩溃
    select {}  // 阻塞主 goroutine
}
```

**运行结果**:
```
panic: something went wrong

goroutine 6 [running]:
main.handler(...)
    /path/to/main.go:10
net/http.HandlerFunc.ServeHTTP(...)
    /usr/local/go/src/net/http/server.go:2136
net/http.serverHandler.ServeHTTP(...)
    /usr/local/go/src/net/http/server.go:2968
net/http.(*conn).serve(...)
    /usr/local/go/src/net/http/server.go:1967
created by net/http.(*Server).Serve
    /usr/local/go/src/net/http/server.go:3109

Process finished with exit code 2  ← 程序崩溃退出
```

### 3.2 为什么会导致整个程序崩溃？

这是 **Go 语言的设计决策**：

1. **Panic 是严重错误**: Go 认为 panic 是程序无法恢复的严重错误
2. **避免静默失败**: 如果 panic 被忽略，可能导致数据不一致等问题
3. **快速失败原则**: 让问题立即暴露，而不是隐藏

### 3.3 实际测试

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func handler1(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handler 1: processing...")
    time.Sleep(1 * time.Second)
    w.Write([]byte("Handler 1: OK"))
}

func handler2(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handler 2: processing...")
    panic("Handler 2 panicked!")  // ← 这个 panic 会导致整个程序崩溃
    w.Write([]byte("Handler 2: OK"))
}

func main() {
    http.HandleFunc("/handler1", handler1)
    http.HandleFunc("/handler2", handler2)
    
    go func() {
        if err := http.ListenAndServe(":8080", nil); err != nil {
            fmt.Println("Server error:", err)
        }
    }()
    
    // 等待服务器启动
    time.Sleep(100 * time.Millisecond)
    
    // 测试：先访问 handler1（正常）
    go func() {
        time.Sleep(200 * time.Millisecond)
        http.Get("http://localhost:8080/handler1")
    }()
    
    // 测试：然后访问 handler2（会 panic）
    go func() {
        time.Sleep(500 * time.Millisecond)
        http.Get("http://localhost:8080/handler2")
    }()
    
    // 主 goroutine 等待
    time.Sleep(2 * time.Second)
    fmt.Println("Main goroutine still running...")
}
```

**运行结果**:
```
Handler 1: processing...
Handler 2: processing...
panic: Handler 2 panicked!

goroutine 8 [running]:
main.handler2(...)
    /path/to/main.go:20
...

Process finished with exit code 2  ← 整个程序崩溃，handler1 的请求也被中断
```

**关键观察**:
- Handler 1 正在处理请求
- Handler 2 发生 panic
- **整个程序崩溃**，Handler 1 的请求也被中断

## 四、为什么需要 Recovery 中间件？

### 4.1 没有 Recovery 的问题

```go
// ❌ 没有 Recovery
func handler(w http.ResponseWriter, r *http.Request) {
    panic("error")  // 导致整个程序崩溃
}

// 问题：
// 1. 一个请求的 panic 导致整个服务器崩溃
// 2. 所有正在处理的请求都被中断
// 3. 服务不可用
```

### 4.2 有 Recovery 的解决方案

```go
// ✅ 有 Recovery
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // 恢复 panic，返回错误响应
                http.Error(w, "Internal Server Error", 500)
                // 记录日志
                log.Printf("Panic recovered: %v", err)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    panic("error")  // 被 recoveryMiddleware 捕获，不会导致程序崩溃
}

// 优势：
// 1. 单个请求的 panic 不会影响其他请求
// 2. 服务器保持运行
// 3. 可以返回友好的错误响应
```

### 4.3 实际测试（有 Recovery）

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                http.Error(w, "Internal Server Error", 500)
                log.Printf("Panic recovered: %v", err)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func handler1(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handler 1: processing...")
    time.Sleep(1 * time.Second)
    w.Write([]byte("Handler 1: OK"))
    fmt.Println("Handler 1: completed")
}

func handler2(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handler 2: processing...")
    panic("Handler 2 panicked!")  // ← 被 recoveryMiddleware 捕获
    w.Write([]byte("Handler 2: OK"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/handler1", handler1)
    mux.HandleFunc("/handler2", handler2)
    
    handler := recoveryMiddleware(mux)
    
    go func() {
        if err := http.ListenAndServe(":8080", handler); err != nil {
            fmt.Println("Server error:", err)
        }
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    // 测试：先访问 handler1（正常）
    go func() {
        time.Sleep(200 * time.Millisecond)
        http.Get("http://localhost:8080/handler1")
    }()
    
    // 测试：然后访问 handler2（会 panic，但被恢复）
    go func() {
        time.Sleep(500 * time.Millisecond)
        http.Get("http://localhost:8080/handler2")
    }()
    
    time.Sleep(2 * time.Second)
    fmt.Println("Main goroutine still running...")  // ← 程序继续运行
}
```

**运行结果**:
```
Handler 1: processing...
Handler 2: processing...
2024/01/01 12:00:00 Panic recovered: Handler 2 panicked!
Handler 1: completed  ← Handler 1 正常完成
Main goroutine still running...  ← 程序继续运行
```

**关键观察**:
- Handler 2 发生 panic，但被 recoveryMiddleware 捕获
- Handler 1 正常完成，不受影响
- **程序继续运行**，没有崩溃

## 五、各框架的 Recovery 实现

### 5.1 Gin 的 Recovery

```go
// recovery.go
func Recovery() HandlerFunc {
    return func(c *Context) {
        defer func() {
            if err := recover(); err != nil {
                // 恢复 panic，返回错误响应
                c.AbortWithStatus(http.StatusInternalServerError)
            }
        }()
        c.Next()  // 继续执行下一个中间件/处理器
    }
}

// 使用
r := gin.Default()  // 默认包含 Recovery 中间件
// 或
r := gin.New()
r.Use(gin.Recovery())
```

### 5.2 GoFrame 的 Recovery

```go
// net/ghttp/ghttp_server.go
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if err := recover(); err != nil {
            s.handlePanic(w, r, err)  // 处理 panic
        }
    }()
    // 处理请求
}

// 内置了 Recovery，不需要手动添加
```

### 5.3 Fiber 的 Recovery

```go
// middleware/recover/recover.go
func New(config ...Config) fiber.Handler {
    return func(c fiber.Ctx) (err error) {
        defer func() {
            if r := recover(); r != nil {
                // 转换为 error
                if e, ok := r.(error); ok {
                    err = e
                } else {
                    err = fmt.Errorf("%v", r)
                }
            }
        }()
        return c.Next()  // 继续执行
    }
}

// 使用
app := fiber.New()
app.Use(recover.New())
```

### 5.4 Kratos 的 Recovery

```go
// middleware/recovery/recovery.go
func Recovery(opts ...Option) middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
            defer func() {
                if rerr := recover(); rerr != nil {
                    // 记录堆栈
                    buf := make([]byte, 64<<10)
                    n := runtime.Stack(buf, false)
                    log.Context(ctx).Errorf("%v: %+v\n%s\n", rerr, req, buf)
                    // 返回错误
                    err = errors.InternalServer("INTERNAL_ERROR", "internal server error")
                }
            }()
            return handler(ctx, req)
        }
    }
}
```

## 六、为什么需要 Recovery？

### 6.1 隔离错误

```go
// 没有 Recovery：一个请求的 panic 影响所有请求
请求1 (goroutine 1) → panic → 整个程序崩溃
请求2 (goroutine 2) → 被中断
请求3 (goroutine 3) → 被中断

// 有 Recovery：一个请求的 panic 只影响该请求
请求1 (goroutine 1) → panic → 被捕获，返回 500
请求2 (goroutine 2) → 正常处理
请求3 (goroutine 3) → 正常处理
```

### 6.2 提高可用性

```go
// 没有 Recovery
服务器可用性: 0% (一个 panic 就崩溃)

// 有 Recovery
服务器可用性: 99.9% (单个请求失败不影响其他请求)
```

### 6.3 用户体验

```go
// 没有 Recovery
用户看到: 连接错误 / 服务器无响应

// 有 Recovery
用户看到: 500 Internal Server Error (友好的错误响应)
```

## 七、最佳实践

### 7.1 总是使用 Recovery

```go
// ✅ 好的方式：总是使用 Recovery
r := gin.Default()  // 默认包含 Recovery

// 或
r := gin.New()
r.Use(gin.Recovery())
```

### 7.2 记录详细日志

```go
// ✅ 好的方式：记录详细日志
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // 记录堆栈跟踪
                buf := make([]byte, 64<<10)
                n := runtime.Stack(buf, false)
                log.Printf("Panic recovered: %v\n%s", err, buf[:n])
                
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### 7.3 区分错误类型

```go
// ✅ 好的方式：区分错误类型
defer func() {
    if err := recover(); err != nil {
        switch e := err.(type) {
        case *MyCustomError:
            // 自定义错误，返回友好消息
            c.JSON(400, gin.H{"error": e.Message})
        case error:
            // 标准错误，记录日志
            log.Error(e)
            c.JSON(500, gin.H{"error": "Internal Server Error"})
        default:
            // 其他类型，记录详细日志
            log.Errorf("panic: %v", err)
            c.JSON(500, gin.H{"error": "Internal Server Error"})
        }
    }
}()
```

## 八、总结

### 8.1 核心要点

1. ✅ **每个 HTTP 请求在独立的 goroutine 中处理**
2. ✅ **未恢复的 panic 会导致整个程序崩溃**
3. ✅ **必须使用 Recovery 中间件来隔离错误**
4. ✅ **Recovery 可以捕获 panic，返回错误响应，而不导致程序崩溃**

### 8.2 关键理解

```
没有 Recovery:
  请求1 (goroutine) → panic → 整个程序崩溃 ❌

有 Recovery:
  请求1 (goroutine) → panic → 被捕获 → 返回 500 → 程序继续运行 ✅
  请求2 (goroutine) → 正常处理 ✅
  请求3 (goroutine) → 正常处理 ✅
```

### 8.3 实际应用

- ✅ **总是使用 Recovery**: 这是 Web 框架的必备功能
- ✅ **记录详细日志**: 便于调试和排查问题
- ✅ **区分错误类型**: 提供更友好的错误响应
- ✅ **监控和告警**: 监控 panic 的发生频率

**记住**: 在 Go 中，任何 goroutine 的未恢复 panic 都会导致整个程序崩溃，所以 Recovery 中间件是**必需的**，而不是可选的！

