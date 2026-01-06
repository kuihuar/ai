# HTTP 服务入口 Recover 机制详解

## 一、核心问题

### 1.1 问题

**如果是提供 HTTP 服务，入口处肯定有 recover 吗？**

**答案**: **是的，几乎所有的 Web 框架都会在入口处提供 recover 机制**，但实现方式不同：

1. **有些框架内置了 recover**（如 GoFrame）
2. **有些框架通过中间件提供 recover**（如 Gin、Fiber）
3. **有些框架需要手动添加 recover 中间件**（如标准库）

## 二、各框架的 Recover 实现

### 2.1 GoFrame - 内置 Recover

**GoFrame 在框架层面内置了 recover**，不需要手动添加：

```go
// net/ghttp/ghttp_server_handler.go
func (s *Server) handleRequest(r *Request) {
    // ...
    defer func() {
        if exception := recover(); exception != nil {
            // 自动捕获 panic，转换为错误响应
            r.Response.WriteStatus(http.StatusInternalServerError, exception)
        }
    }()
    // 处理请求
}
```

**特点**:
- ✅ **内置 recover**，开箱即用
- ✅ 自动捕获 panic，返回 500 错误
- ✅ 不需要手动添加中间件

**使用示例**:

```go
func main() {
    s := g.Server()
    s.Group("/", func(group *ghttp.RouterGroup) {
        group.GET("/panic", func(r *ghttp.Request) {
            panic("something went wrong")  // ← 会被自动捕获
        })
    })
    s.Run()
}
```

**结果**: panic 被自动捕获，返回 500 错误，程序继续运行。

### 2.2 Gin - 通过中间件提供 Recover

**Gin 通过 `gin.Default()` 默认包含 Recovery 中间件**：

```go
// gin/gin.go
func Default() *Engine {
    debugPrintWARNINGDefault()
    engine := New()
    engine.Use(Logger(), Recovery())  // ← 默认包含 Recovery
    return engine
}

func New() *Engine {
    // 不包含任何中间件，需要手动添加
}
```

**特点**:
- ✅ `gin.Default()` **默认包含 Recovery**
- ⚠️ `gin.New()` **不包含 Recovery**，需要手动添加
- ✅ 可以自定义 Recovery 行为

**使用示例**:

```go
// ✅ 方式 1：使用 Default（推荐）
func main() {
    r := gin.Default()  // ← 默认包含 Recovery
    r.GET("/panic", func(c *gin.Context) {
        panic("something went wrong")  // ← 会被捕获
    })
    r.Run()
}

// ⚠️ 方式 2：使用 New（需要手动添加）
func main() {
    r := gin.New()
    r.Use(gin.Recovery())  // ← 需要手动添加
    r.GET("/panic", func(c *gin.Context) {
        panic("something went wrong")
    })
    r.Run()
}

// ❌ 方式 3：使用 New 但不添加 Recovery（危险！）
func main() {
    r := gin.New()  // ← 没有 Recovery
    r.GET("/panic", func(c *gin.Context) {
        panic("something went wrong")  // ← 会导致程序崩溃！
    })
    r.Run()
}
```

### 2.3 Fiber - 通过中间件提供 Recover

**Fiber 需要手动添加 recover 中间件**：

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
        return c.Next()
    }
}
```

**特点**:
- ⚠️ **需要手动添加** recover 中间件
- ✅ 可以自定义配置
- ✅ 支持堆栈跟踪

**使用示例**:

```go
// ✅ 正确方式：手动添加 recover
func main() {
    app := fiber.New()
    app.Use(recover.New())  // ← 需要手动添加
    app.Get("/panic", func(c fiber.Ctx) error {
        panic("something went wrong")  // ← 会被捕获
        return nil
    })
    app.Listen(":3000")
}

// ❌ 错误方式：没有添加 recover
func main() {
    app := fiber.New()  // ← 没有 recover
    app.Get("/panic", func(c fiber.Ctx) error {
        panic("something went wrong")  // ← 会导致程序崩溃！
        return nil
    })
    app.Listen(":3000")
}
```

### 2.4 Kratos - 通过中间件提供 Recover

**Kratos 提供了 Recovery 中间件，但需要手动添加**：

```go
// middleware/recovery/recovery.go
func Recovery(opts ...Option) middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req any) (reply any, err error) {
            defer func() {
                if rerr := recover(); rerr != nil {
                    // 记录堆栈跟踪
                    buf := make([]byte, 64<<10)
                    n := runtime.Stack(buf, false)
                    log.Context(ctx).Errorf("%v: %+v\n%s\n", rerr, req, buf)
                    // 调用自定义处理函数
                    err = op.handler(ctx, req, rerr)
                }
            }()
            return handler(ctx, req)
        }
    }
}
```

**特点**:
- ⚠️ **需要手动添加** Recovery 中间件
- ✅ 可以自定义错误处理
- ✅ 自动记录堆栈跟踪
- ✅ 支持设置延迟信息

**使用示例**:

```go
// ✅ 正确方式：手动添加 Recovery
func main() {
    httpSrv := http.NewServer(
        http.Address(":8000"),
        http.Middleware(
            recovery.Recovery(  // ← 需要手动添加
                recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
                    // 自定义错误处理
                    return errors.InternalServer("INTERNAL_ERROR", "internal server error")
                }),
            ),
        ),
    )
    
    // 注册路由
    httpSrv.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
        panic("something went wrong")  // ← 会被捕获
    })
    
    app := kratos.New(
        kratos.Server(httpSrv),
    )
    app.Run()
}

// ❌ 错误方式：没有添加 Recovery
func main() {
    httpSrv := http.NewServer(
        http.Address(":8000"),
        // 没有添加 Recovery 中间件
    )
    // panic 会导致程序崩溃！
}
```

### 2.5 标准库 - 需要手动实现

**标准库 `net/http` 没有内置 recover**，需要手动实现：

```go
// ❌ 没有 recover（危险！）
func handler(w http.ResponseWriter, r *http.Request) {
    panic("something went wrong")  // ← 会导致程序崩溃！
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

**需要手动添加 recover**:

```go
// ✅ 手动添加 recover
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

func handler(w http.ResponseWriter, r *http.Request) {
    panic("something went wrong")  // ← 会被捕获
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)
    
    handler := recoveryMiddleware(mux)  // ← 包装 recover
    http.ListenAndServe(":8080", handler)
}
```

## 三、框架对比总结

### 3.1 Recover 机制对比

| 框架 | Recover 方式 | 默认包含 | 是否需要手动添加 |
|------|-------------|---------|----------------|
| **GoFrame** | 内置在框架层面 | ✅ 是 | ❌ 不需要 |
| **Gin (Default)** | 中间件 | ✅ 是 | ❌ 不需要 |
| **Gin (New)** | 中间件 | ❌ 否 | ✅ 需要 |
| **Fiber** | 中间件 | ❌ 否 | ✅ 需要 |
| **Kratos** | 中间件 | ❌ 否 | ✅ 需要 |
| **标准库** | 无 | ❌ 否 | ✅ 需要手动实现 |

### 3.2 最佳实践

#### GoFrame

```go
// ✅ 开箱即用，不需要任何配置
s := g.Server()
s.GET("/panic", func(r *ghttp.Request) {
    panic("error")  // 自动捕获
})
s.Run()
```

#### Gin

```go
// ✅ 推荐：使用 Default
r := gin.Default()  // 默认包含 Recovery
r.GET("/panic", handler)

// ⚠️ 如果使用 New，必须手动添加
r := gin.New()
r.Use(gin.Recovery())  // 必须添加
r.GET("/panic", handler)
```

#### Fiber

```go
// ✅ 必须手动添加
app := fiber.New()
app.Use(recover.New())  // 必须添加
app.Get("/panic", handler)
```

#### Kratos

```go
// ✅ 必须手动添加
httpSrv := http.NewServer(
    http.Address(":8000"),
    http.Middleware(
        recovery.Recovery(  // 必须添加
            recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
                return errors.InternalServer("INTERNAL_ERROR", "internal server error")
            }),
        ),
    ),
)
```

#### 标准库

```go
// ✅ 必须手动实现
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## 四、为什么入口处必须有 Recover？

### 4.1 防止程序崩溃

```go
// 没有 recover
请求1 (goroutine) → panic → 整个程序崩溃 ❌
请求2 (goroutine) → 被中断 ❌
请求3 (goroutine) → 被中断 ❌

// 有 recover
请求1 (goroutine) → panic → 被捕获 → 返回 500 ✅
请求2 (goroutine) → 正常处理 ✅
请求3 (goroutine) → 正常处理 ✅
```

### 4.2 提高可用性

```go
// 没有 recover
服务器可用性: 0% (一个 panic 就崩溃)

// 有 recover
服务器可用性: 99.9% (单个请求失败不影响其他请求)
```

### 4.3 用户体验

```go
// 没有 recover
用户看到: 连接错误 / 服务器无响应

// 有 recover
用户看到: 500 Internal Server Error (友好的错误响应)
```

## 五、如何确认框架是否有 Recover？

### 5.1 查看框架文档

- **GoFrame**: 文档明确说明内置 recover
- **Gin**: 文档说明 `Default()` 包含 Recovery
- **Fiber**: 文档说明需要手动添加 recover 中间件

### 5.2 查看源码

```go
// 查看框架的入口处理函数
// 是否有 defer recover() 或 Recovery 中间件
```

### 5.3 实际测试

```go
// 创建一个会 panic 的 handler
func panicHandler(c *gin.Context) {
    panic("test panic")
}

// 如果程序没有崩溃，说明有 recover
// 如果程序崩溃了，说明没有 recover
```

## 六、总结

### 6.1 核心答案

**是的，HTTP 服务的入口处应该有 recover**，但不同框架的实现方式不同：

1. ✅ **GoFrame**: 内置 recover，开箱即用
2. ✅ **Gin (Default)**: 默认包含 Recovery 中间件
3. ⚠️ **Gin (New)**: 需要手动添加 Recovery
4. ⚠️ **Fiber**: 需要手动添加 recover 中间件
5. ⚠️ **标准库**: 需要手动实现 recover

### 6.2 关键要点

- ✅ **总是使用有 recover 的框架或中间件**
- ✅ **如果使用 `gin.New()` 或 `fiber.New()`，必须手动添加 recover**
- ✅ **在生产环境中，recover 是必需的，不是可选的**
- ✅ **定期检查框架是否启用了 recover**

### 6.3 检查清单

在启动 HTTP 服务前，确认：

- [ ] 框架是否内置了 recover？
- [ ] 如果使用中间件，是否添加了 Recovery 中间件？
- [ ] Recovery 是否配置正确？
- [ ] 是否记录了 panic 日志？
- [ ] 是否返回了友好的错误响应？

**记住**: 在生产环境中，**recover 是必需的**，没有 recover 的 HTTP 服务是**不安全的**！

