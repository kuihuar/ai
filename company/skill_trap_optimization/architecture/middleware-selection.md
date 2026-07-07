# Web 框架与中间件选型

## 一、Web 框架选型

### 1.1 候选框架概览

| 框架 | Stars | 底层 HTTP | 路由算法 | 性能 | 社区生态 |
|------|-------|-----------|----------|------|----------|
| **Gin** | ~77k | net/http | Radix Tree | 中 | 最成熟，插件最多 |
| **Chi** | ~18k | net/http | Radix Tree | 中 | 轻量，标准库风格 |
| **Fiber** | ~33k | fasthttp | Radix Tree | 高 | 增长快，受 Express.js 启发 |
| **Echo** | ~29k | net/http | Radix Tree | 中偏高 | 成熟，功能全面 |
| **net/http** | 标准库 | net/http | - | 中 | 零依赖，Go 1.22+ 路由加强 |

### 1.2 逐项对比

#### Gin

```go
r := gin.Default()
r.GET("/users/:id", func(c *gin.Context) { ... })
r.POST("/users", func(c *gin.Context) { ... })
r.Run(":8080")
```

| 优点 | 缺点 |
|------|------|
| 社区最大，中文资料最丰富 | Context 对象大量使用 `interface{}`（Any/MustBind），失去类型安全 |
| 中间件生态最全（gin-contrib） | 并发设置不当（高并发）可能有边界问题 |
| 验证器内置（binding tag），开箱即用 | Handler 签名是 `gin.Context`，耦合框架 |
| 最成熟稳定，生产验证案例最多 | 性能不如 Fiber/fasthttp 栈 |

**推荐场景**：最通用的选择，尤其适合中文社区、团队新手较多的情况。

#### Chi

```go
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) { ... })
http.ListenAndServe(":8080", r)
```

| 优点 | 缺点 |
|------|------|
| 完全兼容 `net/http` Handler 签名，`func(w, r)` | 生态不如 Gin 丰富 |
| 中间件兼容所有 net/http 中间件 | 无内置验证器，需要额外引入 |
| 轻量、简洁、Go 惯用 | 学习资料相对少 |
| 路由支持嵌套、分组、中间件作用域 | |

**推荐场景**：偏好标准库风格、不希望被框架绑定的团队。

#### Fiber

```go
app := fiber.New()
app.Get("/users/:id", func(c *fiber.Ctx) error { ... })
app.Listen(":8080")
```

| 优点 | 缺点 |
|------|------|
| 性能最高（基于 fasthttp，比 gin 快 5-10 倍） | **不兼容 net/http**，不能用 gin/chi 中间件 |
| API 类似 Express.js，前端全栈友好 | fasthttp 有自身生态限制 |
| 内置了大量实用功能 | 遇到 net/http 中间件需要用 adaptor 转换 |

**注意**：Fiber 底层使用 fasthttp，而非 Go 标准 net/http。这意味着：
- 标准库的 `http.Handler` 中间件不能直接用，需要用 `adaptor`
- 社区生态与 net/http 体系不互通

**推荐场景**：对延迟和吞吐有极致要求、不依赖 net/http 生态的中间件。

#### Echo

```go
e := echo.New()
e.GET("/users/:id", func(c echo.Context) error { ... })
e.Start(":8080")
```

| 优点 | 缺点 |
|------|------|
| 功能全面（内置 render/validator/binder/middleware） | 社区和生态不如 Gin |
| 集中化错误处理（`e.HTTPErrorHandler`） | 更新和维护节奏不稳定 |
| 支持 HTTP/2、自动 TLS | 中文资料较少 |

**推荐场景**：想要 Gin 级别的功能但希望框架更专注和高性能。

#### 标准库 net/http（Go 1.22+）

Go 1.22 加强了默认路由功能，支持路径参数和 HTTP method：
```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users/{id}", handler)
```

| 优点 | 缺点 |
|------|------|
| 零依赖 | 无内置中间件链 |
| 完全掌控 | 无参数验证、render 等功能 |
| 永远不会被废弃 | 需要用三方库补充中间件 |

**推荐场景**：极简服务、不需要复杂路由、对依赖有严格限制的项目。

### 1.3 框架选型决策树

```
需要高 QPS / 低延迟？
    ├── 是 → 是否依赖 net/http 生态？
    │        ├── 是 → Gin / Echo
    │        └── 否 → Fiber
    └── 否 → 偏好标准库风格？
             ├── 是 → Chi（Go < 1.22）或 net/http（Go 1.22+）
             └── 否 → 团队有 Node.js 背景？
                      ├── 是 → Fiber
                      └── 否 → Gin（最稳妥）
```

---

## 二、中间件详解与选型

中间件执行顺序 = 注册顺序，符合**洋葱模型**：

```
请求 → Recovery → Logger → CORS → Auth → RateLimit → Tracing → Handler
响应 ← Recovery ← Logger ← CORS ← Auth ← RateLimit ← Tracing ← Handler
```

### 2.1 Recovery（panic 恢复）

**作用**：捕获 handler 中的 panic，返回 500 错误，避免整个进程崩溃。

```go
// Gin 已内置
r := gin.Default() // 包含 Logger + Recovery

// 自定义 recovery，加上通知
func CustomRecovery() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, err any) {
        // 发送告警通知
        notifyPanic(c, err)
        c.AbortWithStatusJSON(500, gin.H{"code": 50000, "msg": "internal error"})
    })
}
```

**关键点**：
- 必须放在中间件链第一位（最外层），否则 panic 跳过其他中间件
- goroutine 中发生的 panic **不会被 recovery 捕获**——goroutine 内部需要自己的 recover
- 生产环境不要返回 panic 的原始错误给客户端，暴露内部信息

### 2.2 Logger（请求日志）

| 方案 | 适用 | 特点 |
|------|------|------|
| Gin 默认 Logger | 开发环境 | 彩色输出，信息简单 |
| 自定义 slog 中间件 | 生产环境 | 结构化 JSON，对接日志平台 |
| 第三方（zap/zerolog） | 生产环境 | 高性能，结构化 |

```go
// 使用 slog 的结构化日志中间件
func SlogLogger(logger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        logger.Info("request",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "latency_ms", time.Since(start).Milliseconds(),
            "client_ip", c.ClientIP(),
        )
    }
}
```

**选型**：
- 新项目直接选用 `log/slog`（Go 1.21+ 标准库），零依赖且够用
- 对日志分配有极致要求的场景选 `zap` 或 `zerolog`
- 不要新项目上 `logrus`——性能差，功能不如 slog

### 2.3 CORS（跨域资源共享）

**仅在浏览器调用后端时**才需要 CORS。内部服务间调用不需要。

```go
// Gin 官方中间件
import "github.com/gin-contrib/cors"

r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))

// Chi
import "github.com/go-chi/cors"
r.Use(cors.Handler(cors.Options{...}))
```

**关键点**：
- `AllowOrigins` 不要写成 `"*"`，生产环境要列明白名单域名
- `AllowCredentials: true` 时不能同时 `AllowOrigins: ["*"]`（浏览器会拒绝）
- 预检请求（OPTIONS）的 `MaxAge` 不要设 0，合理设置可减少预检请求次数

### 2.4 Auth（认证）

#### JWT 方案

```go
import "github.com/golang-jwt/jwt/v5"

func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
            c.AbortWithStatusJSON(401, gin.H{"msg": "missing token"})
            return
        }
        tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims,
            func(token *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"msg": "invalid token"})
            return
        }
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

#### 认证方案对比

| 方案 | 适用 | 优缺点 |
|------|------|--------|
| JWT (Access + Refresh) | 无状态 API，前后端分离 | 无状态好扩展；token 泄露无法失效（除非黑名单） |
| Session + Cookie | 传统 Web 应用 | 服务端可控强；不易水平扩展 |
| OAuth2 / OIDC | 三方登录、SSO | 标准协议，安全；实现复杂 |
| API Key | 机器间调用、Open API | 简单；安全粒度粗 |

**推荐 JWT 双 Token 模式**：
- Access Token：短期（15min-2h），放在内存或 header
- Refresh Token：长期（7-30天），放在 httpOnly Cookie 或安全存储
- Access Token 过期用 Refresh Token 换新的，Refresh Token 过期需重新登录

```go
type Claims struct {
    jwt.RegisteredClaims
    UserID int64  `json:"uid"`
    Role   string `json:"role"`
}
```

#### 可选增强

| 增强功能 | 方案 |
|----------|------|
| Token 黑名单（登出即失效） | Redis 存黑名单 Access Token（TTL = 过期时间） |
| Casbin RBAC/ABAC | 基于角色的权限控制（role-permission） |
| 多租户隔离 | 从 JWT Claims 取 tenant_id，注入 Context |

### 2.5 RateLimit（限流）

#### 方案对比

| 方案 | 算法 | 适用 | 分布式的支持 |
|------|------|------|-------------|
| `golang.org/x/time/rate` | 令牌桶 | 单实例限流 | 需自己实现 |
| `github.com/ulule/limiter` | 固定窗口/滑动窗口 | 单实例或 Redis 共享 | 内建 Redis 支持 |
| API Gateway (Kong/APISIX) | 多种 | 网关层统一限流 | 原生支持 |
| K8s Ingress Controller | 连接数/IP | K8s 部署场景 | 原生支持 |

```go
// ulule/limiter + Gin 中间件（单实例，内存存储）
import "github.com/ulule/limiter/v3"

rate := limiter.Rate{Period: 1 * time.Minute, Limit: 100}
store := memory.NewStore()
instance := limiter.New(store, rate)
middleware := ginlimiter.NewMiddleware(instance)
r.Use(middleware)

// Redis 共享（多实例）
store, _ := redis.NewStore(client)
```

#### 限流策略

| 级别 | 目标 | 示例 |
|------|------|------|
| 全局 | 保护整个服务 | 1000 req/s |
| 路由 | 保护高消耗接口 | `/api/export` 10 req/min |
| 用户/IP | 防止单用户滥用 | 每个 IP 100 req/min |
| 租户 | SaaS 多租户配额 | 租户A 1000 req/min |

**选型**：
- 单实例 → `x/time/rate` 最轻量
- 多实例需统一计数 → `ulule/limiter` + Redis
- 最推荐：在 API Gateway 层统一做限流（服务本身不做）

### 2.6 Tracing（链路追踪）

**统一标准**：OpenTelemetry（CNCF 项目），各语言的 SDK、自动插桩、多 exporter（Jaeger/Zipkin/OTLP）。

```go
import (
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// 初始化 tracer
func initTracer() (*sdktrace.TracerProvider, error) {
    exp, _ := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("jaeger:4317"),
        otlptracegrpc.WithInsecure(),
    )
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
    )
    otel.SetTracerProvider(tp)
    return tp, nil
}

// Gin 中间件自动为每个请求创建 span
r.Use(otelgin.Middleware("my-service"))
```

**选型**：一律用 OpenTelemetry 作为 SDK，后端（exporter）根据基础设施选：
- **Jaeger**：最常用，自建或使用 Jaeger all-in-one
- **Grafana Tempo**：已用 Grafana 生态就选它
- **Datadog / New Relic**：商业 APM 平台
- **OTLP gRPC**：标准协议，各后端通吃

### 2.7 中间件完整注册

```go
func setupRouter() *gin.Engine {
    r := gin.New()

    // 1. Recovery：最外层，捕获一切 panic
    r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
        c.AbortWithStatusJSON(500, gin.H{"code": 50000, "msg": "internal error"})
    }))

    // 2. Tracing：尽早注入，追踪全链路
    r.Use(otelgin.Middleware("my-service"))

    // 3. Logger：记录请求信息（用 tracing span id 关联）
    r.Use(slogMiddleware())

    // 4. CORS（浏览器才需要）
    r.Use(cors.Default())

    // 5. Auth（按需保护特定路由）
    // 注册顺序靠后，可以在前面中间件拿到 tracing/slog context

    // 6. RateLimit 只应用到高消耗路由
    return r
}
```
