# Kratos OpenTelemetry 集成指南

## 概述

本文档说明如何在 Kratos 框架中集成 OpenTelemetry 分布式追踪，实现完整的请求追踪链路。

> **📌 配置位置说明**：OpenTelemetry 配置定义在 **`internal/conf/conf.proto`** 中，这是项目的统一配置定义文件。配置通过 `internal/config` 包加载，支持 Viper 和 Kratos 两种配置系统。详见 [配置管理](#配置管理) 章节。

### 为什么需要 OpenTelemetry？

- **分布式追踪**：跟踪请求在微服务架构中的完整流转路径
- **性能分析**：识别系统瓶颈和慢请求
- **问题排查**：快速定位错误发生的服务和方法
- **服务依赖关系**：可视化服务间的调用关系

### Kratos 对 OpenTelemetry 的支持

Kratos 框架原生支持 OpenTelemetry，提供了 tracing 中间件：

- **服务端中间件**：`github.com/go-kratos/kratos/v2/middleware/tracing`
  - 自动为每个请求创建 span
  - 自动传播 trace context
  - 自动记录请求元数据（方法、路径、状态码等）

- **客户端中间件**：支持 HTTP 和 gRPC 客户端
  - 自动注入 trace context 到请求头
  - 自动创建客户端 span

## 前置条件

### 1. 依赖检查

确保项目中已引入必要的依赖：

```go
// go.mod
require (
    go.opentelemetry.io/otel v1.34.0
    go.opentelemetry.io/otel/trace v1.34.0
    go.opentelemetry.io/otel/sdk v1.34.0
    go.opentelemetry.io/otel/exporters/jaeger v1.17.0  // 如果使用 Jaeger
    go.opentelemetry.io/otel/exporters/zipkin v1.17.0  // 如果使用 Zipkin
    go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.34.0  // 如果使用 OTLP
    github.com/go-kratos/kratos/v2 v2.9.0
)
```

### 2. 当前项目状态

项目中已经：
- ✅ 引入了 OpenTelemetry 基础依赖
- ✅ 在第三方服务调用中使用了 OpenTelemetry（见 `docs/development/opentelemetry-tracing-third-party.md`）
- ✅ 日志系统集成了 OpenTelemetry trace 信息提取
- ❌ **尚未在 Kratos HTTP/gRPC 服务器中使用 tracing 中间件**

## 实现步骤

### 步骤 1: 初始化 TracerProvider

首先需要初始化 OpenTelemetry TracerProvider，这是追踪系统的核心组件。

#### 1.1 创建 TracerProvider 初始化函数

创建 `internal/tracing/provider.go`：

```go
package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Config TracerProvider 配置
type Config struct {
	ServiceName    string // 服务名称
	ServiceVersion string // 服务版本
	Environment    string // 环境（dev, staging, prod）
	
	// 导出器配置（三选一）
	JaegerEndpoint string // Jaeger 端点，如: http://localhost:14268/api/traces
	ZipkinEndpoint string  // Zipkin 端点，如: http://localhost:9411/api/v2/spans
	OTLPEndpoint   string  // OTLP 端点，如: localhost:4317
	
	// 采样配置
	SamplingRatio float64 // 采样率，0.0-1.0，1.0 表示采样所有请求
}

// InitTracerProvider 初始化 OpenTelemetry TracerProvider
func InitTracerProvider(ctx context.Context, cfg Config) (trace.TracerProvider, func(), error) {
	// 创建 Resource（描述服务的元数据）
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建导出器（Exporter）
	var exporter sdktrace.SpanExporter
	var exporterName string
	
	switch {
	case cfg.JaegerEndpoint != "":
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
		exporterName = "jaeger"
	case cfg.ZipkinEndpoint != "":
		exporter, err = zipkin.New(cfg.ZipkinEndpoint)
		exporterName = "zipkin"
	case cfg.OTLPEndpoint != "":
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
			otlptracegrpc.WithInsecure(), // 生产环境应使用 TLS
		)
		exporterName = "otlp"
	default:
		return nil, nil, fmt.Errorf("no exporter endpoint configured")
	}
	
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create %s exporter: %w", exporterName, err)
	}

	// 配置采样率
	samplingRatio := cfg.SamplingRatio
	if samplingRatio <= 0 {
		samplingRatio = 1.0 // 默认采样所有请求
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // 批量导出
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(samplingRatio)), // 基于采样率的采样器
	)

	// 设置全局 TracerProvider 和 TextMapPropagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C Trace Context
		propagation.Baggage{},      // W3C Baggage
	))

	// 返回清理函数
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			// 记录错误，但不影响程序退出
			fmt.Printf("failed to shutdown TracerProvider: %v\n", err)
		}
	}

	return tp, cleanup, nil
}
```

#### 1.2 在 main.go 中初始化

在 `cmd/sre/main.go` 中初始化 TracerProvider：

```go
package main

import (
	"context"
	"flag"
	"os"

	"sre/internal/config"
	"sre/internal/logger"
	"sre/internal/tracing" // 新增

	_ "go.uber.org/automaxprocs"
)

func main() {
	flag.Parse()

	// 加载配置
	bootstrap, err := config.LoadBootstrapWithViper(flagconf)
	if err != nil {
		panic(err)
	}

	// 初始化 logger
	logger := logger.NewZapLoggerWithConfig(
		bootstrap.Log.Level,
		bootstrap.Log.Format,
		bootstrap.Log.OutputPaths,
		Name,
		id,
		Version,
	)

	// ============================================
	// 初始化 OpenTelemetry TracerProvider
	// ============================================
	ctx := context.Background()
	
	// 从配置中读取 Tracing 配置
	var tracingConfig tracing.Config
	if bootstrap.Tracing != nil {
		tracingConfig = tracing.Config{
			ServiceName:    bootstrap.Tracing.ServiceName,
			ServiceVersion: bootstrap.Tracing.ServiceVersion,
			Environment:    bootstrap.Tracing.Environment,
			JaegerEndpoint: bootstrap.Tracing.JaegerEndpoint,
			ZipkinEndpoint: bootstrap.Tracing.ZipkinEndpoint,
			OTLPEndpoint:   bootstrap.Tracing.OtlpEndpoint,
			SamplingRatio:  bootstrap.Tracing.SamplingRatio,
		}
	} else {
		// 如果没有配置，使用默认值
		tracingConfig = tracing.Config{
			ServiceName:    Name,
			ServiceVersion: Version,
			Environment:    "dev",
			SamplingRatio:  1.0,
		}
	}
	
	tp, tracingCleanup, err := tracing.InitTracerProvider(ctx, tracingConfig)
	if err != nil {
		log.NewHelper(logger).Warnf("failed to initialize TracerProvider: %v", err)
		// 不中断程序启动，使用 NoOp TracerProvider
	} else {
		log.NewHelper(logger).Info("OpenTelemetry TracerProvider initialized")
	}
	defer tracingCleanup()

	// 启动应用
	app, cleanup, err := wireApp(bootstrap.Server, bootstrap.Data, bootstrap.Registry, bootstrap.Service, bootstrap.Worker, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

### 步骤 2: 在 HTTP 服务器中使用 Tracing 中间件

修改 `internal/server/http.go`，添加 tracing 中间件：

```go
package server

import (
	// ... 其他导入
	"github.com/go-kratos/kratos/v2/middleware/tracing" // 新增
)

func NewHTTPServer(
	c *conf.Server,
	user *service.UserService,
	cache biz.CacheRepo,
	serviceConf *conf.Service,
	logger log.Logger,
) *http.Server {
	globalChain := middleware.NewChain(logger)
	
	// 1. Recovery 中间件（最外层，必须）
	globalChain.Add(recovery.Recovery())
	
	// 2. Tracing 中间件（在 Recovery 之后，其他中间件之前）
	// 这样即使发生 panic，也能记录到 trace 中
	globalChain.Add(tracing.Server())
	
	// 3. 其他中间件（限流、认证等）
	// ... 现有代码 ...
	
	var opts = []http.ServerOption{
		http.Middleware(globalChain.ToSlice()...),
	}
	// ... 其他配置 ...
	
	srv := http.NewServer(opts...)
	userv1.RegisterUserHTTPServer(srv, user)
	return srv
}
```

### 步骤 3: 在 gRPC 服务器中使用 Tracing 中间件

修改 `internal/server/grpc.go`，添加 tracing 中间件：

```go
package server

import (
	// ... 其他导入
	"github.com/go-kratos/kratos/v2/middleware/tracing" // 新增
)

func NewGRPCServer(c *conf.Server, user *service.UserService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(), // 新增：添加 tracing 中间件
		),
	}
	// ... 其他配置 ...
	
	srv := grpc.NewServer(opts...)
	userv1.RegisterUserServer(srv, user)
	return srv
}
```

### 步骤 4: 在 HTTP 客户端中使用 Tracing 中间件

如果项目中有 HTTP 客户端调用其他服务，也需要添加 tracing 中间件：

```go
import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

// 创建 HTTP 客户端
client := http.NewClient(
	http.WithMiddleware(
		tracing.Client(), // 自动注入 trace context
	),
	http.WithEndpoint("http://other-service:8080"),
)
```

### 步骤 5: 在 gRPC 客户端中使用 Tracing 中间件

```go
import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

// 创建 gRPC 客户端
conn, err := grpc.DialInsecure(
	context.Background(),
	grpc.WithEndpoint("other-service:9000"),
	grpc.WithMiddleware(
		tracing.Client(), // 自动注入 trace context
	),
)
```

## 配置管理

### 配置位置

OpenTelemetry 配置定义在 **`internal/conf/conf.proto`** 中，这是项目的统一配置定义文件。

### 配置定义

在 `internal/conf/conf.proto` 中已经定义了 Tracing 配置：

```protobuf
message Bootstrap {
  // ... 其他配置
  Tracing tracing = 8;  // OpenTelemetry 追踪配置
}

message Tracing {
  string service_name = 1;      // 服务名称（用于标识追踪中的服务）
  string service_version = 2;    // 服务版本
  string environment = 3;       // 环境（dev, staging, prod）
  
  // 导出器配置（三选一，优先级：jaeger > zipkin > otlp）
  string jaeger_endpoint = 4;   // Jaeger 端点，如: http://localhost:14268/api/traces
  string zipkin_endpoint = 5;   // Zipkin 端点，如: http://localhost:9411/api/v2/spans
  string otlp_endpoint = 6;     // OTLP 端点，如: localhost:4317
  
  // 采样配置
  double sampling_ratio = 7;    // 采样率，0.0-1.0，1.0 表示采样所有请求（默认 1.0）
}
```

### 配置文件

在 `configs/config.yaml` 中添加配置：

```yaml
tracing:
  service_name: "sre"
  service_version: "v1.0.0"
  environment: "dev"
  jaeger_endpoint: "http://localhost:14268/api/traces"
  sampling_ratio: 1.0
```

### 配置加载

配置通过 `internal/config` 包加载，该包支持：
- ✅ Viper 配置系统（推荐）：支持环境变量、多配置文件、配置热更新
- ✅ Kratos 配置系统：基于 Protobuf 的配置加载

配置会自动从 `configs/config.yaml` 加载，并转换为 `conf.Bootstrap` 结构。Tracing 配置位于 `bootstrap.Tracing` 字段中。

## 导出器选择

### Jaeger（推荐用于开发环境）

Jaeger 是最流行的分布式追踪系统之一，易于部署和使用。

**优点**：
- 提供完整的 UI 界面
- 支持多种存储后端（内存、Cassandra、Elasticsearch）
- 社区活跃，文档完善

**部署**：
```bash
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  jaegertracing/all-in-one:latest
```

访问 UI：http://localhost:16686

### Zipkin

Zipkin 是另一个流行的分布式追踪系统。

**优点**：
- 轻量级
- 支持多种存储后端
- 提供简洁的 UI

**部署**：
```bash
docker run -d --name zipkin \
  -p 9411:9411 \
  openzipkin/zipkin:latest
```

访问 UI：http://localhost:9411

### OTLP（推荐用于生产环境）

OTLP（OpenTelemetry Protocol）是 OpenTelemetry 的标准协议，可以导出到任何支持 OTLP 的后端。

**优点**：
- 标准协议，兼容性好
- 可以导出到多种后端（Jaeger、Zipkin、Tempo、Datadog 等）
- 支持 gRPC 和 HTTP 两种传输方式

**使用 OTLP Collector**：
```bash
docker run -d --name otel-collector \
  -p 4317:4317 \
  -p 4318:4318 \
  otel/opentelemetry-collector:latest
```

## 最佳实践

### 1. 中间件顺序

正确的中间件顺序很重要：

```go
globalChain := middleware.NewChain(logger)

// 1. Recovery（最外层，捕获 panic）
globalChain.Add(recovery.Recovery())

// 2. Tracing（在 Recovery 之后，记录所有请求）
globalChain.Add(tracing.Server())

// 3. 日志中间件（记录请求日志）
globalChain.Add(logging.Server(logger))

// 4. 限流中间件
globalChain.Add(ratelimit.Server(...))

// 5. 认证中间件
globalChain.Add(auth.Server(...))

// 6. 业务逻辑
```

### 2. 采样策略

**开发环境**：
- 采样率：1.0（采样所有请求）
- 便于调试和开发

**生产环境**：
- 采样率：0.1-0.5（采样 10%-50% 的请求）
- 平衡性能和可观测性
- 对于错误请求，可以配置采样率 1.0

### 3. 服务名称规范

使用清晰的服务名称：
- ✅ `sre-api`、`sre-worker`、`sre-cron`
- ❌ `service1`、`app`、`backend`

### 4. Span 命名规范

Kratos tracing 中间件会自动创建 span，命名格式为：
- HTTP: `GET /api/v1/users`
- gRPC: `user.v1.User/GetUser`

在业务代码中创建自定义 span 时，遵循以下规范：
- ✅ `user.CreateUser`、`order.ProcessPayment`
- ❌ `create_user`、`process_payment`、`doSomething`

### 5. 属性（Attributes）使用

在业务代码中添加有意义的属性：

```go
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("order.amount", amount),
    attribute.String("payment.method", "credit_card"),
)
```

### 6. 错误处理

始终记录错误到 span：

```go
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return err
}
```

### 7. Context 传播

确保在所有异步操作中传播 context：

```go
// ✅ 正确：传播 context
go func(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "async.operation")
    defer span.End()
    // ... 业务逻辑
}(ctx)

// ❌ 错误：使用新的 context
go func() {
    ctx := context.Background()
    // ... 业务逻辑
}()
```

## 与现有代码集成

### 与第三方服务调用集成

项目中已经在第三方服务调用中使用了 OpenTelemetry（见 `docs/development/opentelemetry-tracing-third-party.md`）。添加 Kratos tracing 中间件后，这些 span 会自动关联到同一个 trace 中。

**示例**：
```
HTTP Request (Kratos tracing)
  └── dingtalk.GetUserInfo (第三方服务调用)
      └── HTTP Request (Resty 客户端)
```

### 与日志系统集成

日志系统已经集成了 OpenTelemetry trace 信息提取（见 `internal/logger/zap.go`）。添加 tracing 中间件后，日志会自动包含 trace ID 和 span ID。

**日志示例**：
```json
{
  "level": "info",
  "msg": "user created",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "user_id": "12345"
}
```

## 验证和测试

### 1. 检查 TracerProvider 是否初始化

在应用启动后，检查日志中是否有：
```
OpenTelemetry TracerProvider initialized
```

### 2. 发送测试请求

```bash
curl http://localhost:8000/api/v1/users/123
```

### 3. 查看 Jaeger UI

访问 http://localhost:16686，应该能看到：
- Service: `sre`
- Operation: `GET /api/v1/users/:id`
- Trace ID: 完整的追踪链路

### 4. 检查日志中的 Trace ID

查看应用日志，应该能看到每条日志都包含 `trace_id` 和 `span_id`。

## 故障排查

### 问题 1: 没有看到 trace

**可能原因**：
1. TracerProvider 未初始化
2. 导出器配置错误
3. 网络连接问题

**解决方案**：
1. 检查日志中是否有初始化错误
2. 验证导出器端点是否可访问
3. 检查采样率配置

### 问题 2: Trace 不完整

**可能原因**：
1. Context 未正确传播
2. 异步操作未传播 context

**解决方案**：
1. 确保所有函数调用都传递 context
2. 在异步操作中使用 `context.WithValue` 或 `context.WithTimeout`

### 问题 3: 性能影响

**可能原因**：
1. 采样率过高
2. 导出器性能问题

**解决方案**：
1. 降低采样率（生产环境建议 0.1-0.5）
2. 使用批量导出器（已默认启用）
3. 考虑使用异步导出

## 参考资源

- [Kratos Tracing 中间件文档](https://github.com/go-kratos/kratos/tree/main/middleware/tracing)
- [OpenTelemetry Go SDK 文档](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger 文档](https://www.jaegertracing.io/docs/)
- [Zipkin 文档](https://zipkin.io/)
- [OTLP 协议规范](https://opentelemetry.io/docs/specs/otlp/)

## 总结

通过以上步骤，我们实现了：

1. ✅ 初始化 OpenTelemetry TracerProvider
2. ✅ 在 HTTP/gRPC 服务器中使用 tracing 中间件
3. ✅ 在客户端中使用 tracing 中间件
4. ✅ 配置导出器（Jaeger/Zipkin/OTLP）
5. ✅ 与现有代码集成（第三方服务调用、日志系统）

现在整个系统已经具备完整的分布式追踪能力，可以：
- 追踪请求在系统中的完整路径
- 关联日志和追踪信息
- 快速定位性能瓶颈和错误
- 可视化服务依赖关系

