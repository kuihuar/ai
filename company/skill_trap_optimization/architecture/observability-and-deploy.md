# 可观测性、部署与微服务演进

## 一、可观测性

### 1.1 日志库选型

| 库 | 性能（ns/op） | 分配（allocs/op） | 结构化 | 特点 |
|----|-------------|-------------------|--------|------|
| **slog** (Go 1.21+) | ~1200 | ~0 | 支持 | 标准库，无需三方依赖 |
| **zap** (Uber) | ~200 | ~0 | 支持 | 极致性能，生产实战最多 |
| **zerolog** | ~200 | ~0 | JSON 原生 | 极致性能，链式 API |
| **logrus** | ~3000 | ~30+ | 支持 | 不推荐新项目 |

**结论**：新项目一律用 slog，极致性能场景选 zap，不要新上 logrus。

```go
// slog 使用（Go 1.21+标准库）
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
slog.SetDefault(logger)

// 结构化记录
slog.Info("user created",
    "user_id", 123,
    "duration_ms", 45,
)

// 与 Context 集成，自动关联 trace_id
slog.InfoContext(ctx, "processing order", "order_id", 456)
```

### 1.2 指标（Metrics）

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "http_requests_total"},
        []string{"method", "path", "status"},
    )
    httpDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        status := strconv.Itoa(c.Writer.Status())
        httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        httpDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(time.Since(start).Seconds())
    }
}

// 暴露 /metrics 端点
r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

**黄金指标（RED Method）**：
- **Rate**：请求速率（counter）
- **Errors**：错误率（counter，status >= 500）
- **Duration**：请求耗时（histogram）

**USE Method**（资源维度）：CPU 使用率、内存使用、连接数、goroutine 数量。

### 1.3 链路追踪（Tracing）

**一律用 OpenTelemetry**作为 SDK。选型的唯一变量是后端 exporter：

| Exporter | 后端 | 适用场景 |
|----------|------|----------|
| OTLP gRPC | Jaeger / Grafana Tempo / 任意兼容 OTLP | 标准协议，推荐 |
| Jaeger Thrift | Jaeger | Jaeger 老版本 |
| Zipkin | Zipkin | 已有 Zipkin 设施 |
| stdout | 控制台 | 本地调试 |

```go
// OpenTelemetry 初始化（Gin）
import (
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func initTracer(serviceName, otelEndpoint string) (*sdktrace.TracerProvider, error) {
    ctx := context.Background()
    exp, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(otelEndpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.ServiceNameKey.String(serviceName),
        )),
        sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // 10% 采样
    )
    otel.SetTracerProvider(tp)
    return tp, nil
}

// 在中间件链中注入
r.Use(otelgin.Middleware("my-service"))
```

**采样策略**：
- 开发环境：AlwaysSample（100%）
- 生产环境：`TraceIDRatioBased(0.1)`（10%），或 `ParentBased`（有上游 trace 的必采）

---

## 二、API 设计

### 2.1 REST vs gRPC vs GraphQL

| 维度 | REST | gRPC | GraphQL |
|------|------|------|---------|
| 协议 | HTTP/1.1 + JSON | HTTP/2 + Protobuf | HTTP + JSON |
| 性能 | 中 | 高（二进制、多路复用） | 中 |
| 类型安全 | 无（需 OpenAPI/Swagger 补充） | 编译时（proto 生成） | Schema 运行时校验 |
| 浏览器直接调用 | 可以 | 需要 grpc-web | 可以 |
| 学习成本 | 低 | 中（protobuf 语法） | 中 |
| 工具链 | Postman、curl | grpcurl、grpcui | GraphiQL |
| 适用 | 对外 API、前后端分离 | 内部微服务 | 前端灵活查询 |
| 文件上传 | 原生支持 | 需额外处理 | 需额外处理 |
| 流式传输 | SSE/WebSocket | 原生双向流 | Subscription |

### 2.2 RESTful 设计规范

```
GET    /api/v1/users              # 列表（支持 ?page=1&per_page=20&sort=created_at:desc）
GET    /api/v1/users/:id          # 详情
POST   /api/v1/users              # 创建
PUT    /api/v1/users/:id          # 全量更新
PATCH  /api/v1/users/:id          # 部分更新
DELETE /api/v1/users/:id          # 删除
```

### 2.3 统一响应结构

```go
// 成功
{
    "code": 0,
    "message": "ok",
    "data": { ... }
}

// 列表（带分页）
{
    "code": 0,
    "message": "ok",
    "data": [ ... ],
    "total": 100,
    "page": 1,
    "per_page": 20
}

// 错误
{
    "code": 40401,
    "message": "user not found"
}
```

**错误码设计**：
- `0` — 成功
- `4xxxx` — 客户端错误（40001 参数错误、40101 未登录、40301 无权限、40401 资源不存在）
- `5xxxx` — 服务端错误（50001 内部错误、50002 数据库错误、50003 第三方调用失败）

### 2.4 gRPC 设计要点

```protobuf
syntax = "proto3";
package user.v1;
option go_package = "github.com/myproject/api/user/v1;userv1";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

message GetUserRequest {
    int64 user_id = 1;
}
```

- proto 文件是 API 的单一事实来源（SSOT）
- 用 `buf` 管理 proto lint、breaking change 检测、代码生成
- 字段号一旦分配不可更改（protobuf 向后兼容的基础）

---

## 三、测试策略

### 3.1 测试金字塔与工具

| 层级 | 工具 | 数量占比 |
|------|------|----------|
| 单元测试 | `testing` + table-driven | ~70% |
| 集成测试 | `testcontainers` （真实 DB/Redis） | ~20% |
| E2E | 看具体场景（很少，手测为主） | ~10% |

### 3.2 Mock 方案对比

| 方案 | 方式 | 适用 |
|------|------|------|
| **手动 Mock** | 实现接口的结构体 | 接口少、团队小 |
| **gomock** (Google) | 反射生成 Mock | 大型项目，需要精确控制调用次数 |
| **mockery** | 代码生成，基于接口 | 最常用，平衡效率和灵活性 |

```go
// mockery 使用
//go:generate mockery --name UserRepo --output ./mocks

// 生成的 mock 在测试中使用
mockRepo := mocks.NewUserRepo(t)
mockRepo.On("GetByID", ctx, int64(1)).Return(&model.User{ID: 1}, nil)

svc := service.NewUserService(mockRepo)
user, err := svc.GetUser(ctx, 1)
assert.NoError(t, err)
assert.Equal(t, int64(1), user.ID)
mockRepo.AssertExpectations(t)
```

### 3.3 集成测试：testcontainers-go

```go
func TestUserRepo_Create(t *testing.T) {
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:16-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env:          map[string]string{"POSTGRES_PASSWORD": "test", "POSTGRES_DB": "testdb"},
        WaitingFor:   wait.ForListeningPort("5432/tcp"),
    }
    container, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req, Started: true,
    })
    defer container.Terminate(ctx)

    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")
    dsn := fmt.Sprintf("postgres://postgres:test@%s:%s/testdb?sslmode=disable", host, port.Port())
    db, _ := sql.Open("postgres", dsn)

    // 执行迁移
    // ... golang-migrate 跑 up

    // 真实测试
    repo := NewUserRepo(db)
    err := repo.Create(ctx, &model.User{Name: "test"})
    assert.NoError(t, err)
}
```

**testcontainers 的优点**：
- 真实数据库，不是 mock，测试结果可信
- 用完自动销毁，不污染环境
- 支持 Postgres/MySQL/Redis/Kafka/ES 等

---

## 四、构建与部署

### 4.1 Docker 镜像方案对比

| 方案 | 镜像大小 | 安全性 | 适用 |
|------|----------|--------|------|
| `scratch`（空镜像） | ~5MB | 最高（无 shell/工具） | 静态编译（CGO_ENABLED=0） |
| `alpine` | ~8MB | 高 | 需要 ca-certificates/tzdata |
| `distroless` (Google) | ~10MB | 最高 | 生产环境推荐 |
| `ubuntu/debian-slim` | ~80MB | 中 | 需要使用系统库 |

```dockerfile
# 推荐：distroless（Google 维护，生产友好）
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /server /server
COPY --from=builder /app/configs /configs
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
```

**选择**：
- 追求小 → `scratch`
- 需要证书+时区 → `alpine`
- 追求安全+生产 → `distroless`

### 4.2 CI/CD 关键阶段

```yaml
# .github/workflows/ci.yml 关键步骤
steps:
  - lint:        golangci-lint run
  - test:        go test -race -coverprofile=coverage.out ./...
  - build:       go build -o bin/server ./cmd/server
  - docker:      docker build -t myapp:${{ github.sha }} .
  - push:        docker push registry/myapp:${{ github.sha }}
  - deploy:      kubectl set image deployment/myapp myapp=registry/myapp:${{ github.sha }}
```

### 4.3 K8s 健康检查

```yaml
livenessProbe:       # 容器是否存活（失败 → 重启容器）
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 15
readinessProbe:      # 容器是否就绪（失败 → 从 Service 摘除）
  httpGet:
    path: /readyz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

- **/healthz**：检查进程存活，通常只返回 200
- **/readyz**：检查是否能处理请求（DB/Redis 是否连通），就绪前不接收流量
- `initialDelaySeconds` 要大于应用启动时间

---

## 五、微服务演进

### 5.1 何时拆 vs 不拆

| 保持单体 | 开始拆分 |
|----------|----------|
| 团队 < 10 人 | 多团队独立迭代，单团队 > 8 人 |
| 业务模式未定型，变更频繁 | 子领域边界清晰稳定 |
| 数据强耦合，事务跨模块 | 各模块数据独立，最终一致性可接受 |
| 部署频率低 | 需要独立部署、独立扩缩容 |

**原则：先走单体，运行 6-12 个月后根据痛点再拆。** 不要一开始就微服务。

### 5.2 服务通信选型

| 方式 | 同步/异步 | 性能 | 适用 |
|------|-----------|------|------|
| **HTTP REST** | 同步 | 中 | 对外 API、跨语言 |
| **gRPC** | 同步 | 高 | 内部服务间调用 |
| **Kafka** | 异步 | 高吞吐 | 事件驱动、日志/消息流 |
| **NATS / NATS JetStream** | 同步+异步 | 极高 | 轻量消息、云原生 |
| **RabbitMQ** | 异步 | 中 | 复杂路由、传统企业 |

```
内部服务间调用  ──→  gRPC（性能+类型安全）
外部 API         ──→  REST（通用性强）
异步解耦         ──→  Kafka（大吞吐）/ NATS（轻量）
```

### 5.3 服务治理组件选型

| 组件 | 方案 | 适用 |
|------|------|------|
| **服务发现** | K8s Service（DNS）/ Consul / Etcd | K8s 部署直接用 K8s Service |
| **负载均衡** | K8s Service（iptables）/ gRPC client-side LB | 内部服务 L4 LB 用 K8s，gRPC 用客户端 LB |
| **限流** | `golang.org/x/time/rate` / sentinel-golang | 单点限流用 token bucket，分布式用 Sentinel |
| **熔断** | `sony/gobreaker` / `sentinel-golang` | 模式简单用 gobreaker，要阿里生态用 Sentinel |
| **降级** | 业务代码中返回缓存/默认值 | 非核心功能降级 |

### 5.4 分布式事务

| 方案 | 一致性 | 复杂度 | 适用 |
|------|--------|--------|------|
| **Saga（编排/编制）** | 最终一致 | 中 | 长事务，需要补偿逻辑 |
| **Outbox Pattern** | 最终一致 | 低 | 确保 DB 写入和消息发出的原子性 |
| **2PC / TCC** | 强一致 | 高 | 对一致性要求极高（支付） |
| **Eventual Consistency** | 最终一致 | 低 | 大多数场景接受 |

**大部分场景用 Outbox + 最终一致性就够了**：写 DB 时同时写 outbox 表，一个异步 worker 扫描 outbox 发消息。

---

## 六、安全要点

### 6.1 JWT 双 Token 模式

```
Access Token   — 短期（15min-2h），放内存或 Authorization header
Refresh Token  — 长期（7-30day），放 httpOnly Secure Cookie
```

```go
// 登录时返回双 token
func Login(c *gin.Context) {
    accessToken, _ := generateJWT(userID, 15*time.Minute)
    refreshToken, _ := generateJWT(userID, 7*24*time.Hour)

    // Refresh Token 放 httpOnly Cookie
    c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "", true, true)
    c.JSON(200, gin.H{"access_token": accessToken})
}

// /refresh 端点：用 refresh token 换新 access token
func Refresh(c *gin.Context) {
    refreshToken, _ := c.Cookie("refresh_token")
    // 验证 refresh token → 生成新 access token
}
```

### 6.2 安全检查清单

| 检查项 | 做法 |
|--------|------|
| SQL 注入 | 参数化查询，不用字符串拼接 SQL |
| XSS | 输出转义，设置 Content-Type |
| CSRF | 用 SameSite Cookie + CSRF Token（SPA 用 token header） |
| HTTPS | 生产环境强制 HTTPS，HSTS header |
| 敏感信息 | 密码 bcrypt 哈希；密钥从 env/Secret Manager 读取 |
| 依赖漏洞 | `govulncheck ./...`（Go 官方漏洞扫描） |
| 静态编译 | `CGO_ENABLED=0` 减少攻击面 |

### 6.3 错误信息安全

```go
// 错误时不要给客户端返回内部细节
// ❌ 错误做法
c.JSON(500, gin.H{"error": err.Error()})
// → 可能暴露 DB schema、内部 IP、堆栈信息

// ✅ 正确做法：生产环境只返回通用错误，日志中记录详情
slog.Error("db error", "error", err, "user_id", userID)
c.JSON(500, gin.H{"code": 50000, "msg": "internal server error"})
```
