# OpenTelemetry Metrics 集成指南

## 概述

本文档说明如何在 Kratos 框架中集成 OpenTelemetry Metrics，实现应用程序的指标监控。

> **📌 配置位置说明**：OpenTelemetry Metrics 配置定义在 **`internal/conf/conf.proto`** 中，这是项目的统一配置定义文件。配置通过 `internal/config` 包加载，支持 Viper 和 Kratos 两种配置系统。

### 为什么需要 OpenTelemetry Metrics？

- **系统监控**：量化系统状态，如请求数、错误率、响应时间
- **性能分析**：识别系统瓶颈和性能问题
- **告警支持**：基于指标设置告警规则
- **容量规划**：通过历史指标数据规划系统容量

### 当前项目状态

项目中已经：
- ✅ 引入了 OpenTelemetry Metrics 基础依赖
- ✅ 实现了 MeterProvider 初始化（`internal/metrics/provider.go`）
- ✅ 支持多种导出器（Prometheus、OTLP、JSON File）
- ✅ 在主程序中自动初始化 MeterProvider

## 前置条件

### 1. 依赖检查

确保项目中已引入必要的依赖：

```go
// go.mod
require (
    go.opentelemetry.io/otel v1.39.0
    go.opentelemetry.io/otel/sdk/metric v1.39.0
    go.opentelemetry.io/otel/exporters/prometheus v0.61.0
    go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.39.0
)
```

依赖已通过 `go mod tidy` 自动添加。

### 2. 配置检查

在 `configs/config.yaml` 中配置 Metrics：

```yaml
metrics:
  service_name: "sre"
  service_version: "v1.0.0"
  environment: "dev"
  json_file_path: "./metrics/metrics.jsonl"
  export_interval: 10s
  export_timeout: 5s
```

## 快速开始

### 1. 使用中间件（推荐，已自动配置）

Metrics 中间件已自动配置，无需手动编写代码即可自动记录所有 HTTP 和 gRPC 请求的指标。

#### 自动记录的指标

**HTTP 指标：**
- `http_server_requests_total` - HTTP 请求总数
- `http_server_request_duration_seconds` - HTTP 请求耗时
- `http_server_request_size_bytes` - HTTP 请求大小
- `http_server_response_size_bytes` - HTTP 响应大小
- `http_server_active_requests` - 当前活跃请求数

**gRPC 指标：**
- `grpc_server_requests_total` - gRPC 请求总数
- `grpc_server_request_duration_seconds` - gRPC 请求耗时
- `grpc_server_request_size_bytes` - gRPC 请求大小
- `grpc_server_response_size_bytes` - gRPC 响应大小
- `grpc_server_active_requests` - 当前活跃请求数

#### 中间件配置

中间件已在以下位置自动配置：

1. **初始化**（`cmd/sre/main.go`）：
```go
if metricsCleanup, err := metrics.InitMeterProvider(ctx, bootstrap.Metrics, Name, Version, logger); err == nil && metricsCleanup != nil {
    defer metricsCleanup()
    // 初始化 metrics 中间件
    if err := metrics.InitMetricsMiddleware(); err != nil {
        log.NewHelper(logger).Warnf("failed to init metrics middleware: %v", err)
    }
}
```

2. **HTTP 服务器**（`internal/server/http.go`）：
```go
globalChain.Add(tracing.Server())
globalChain.Add(metrics.Server()) // Metrics 中间件
```

3. **gRPC 服务器**（`internal/server/grpc.go`）：
```go
grpc.Middleware(
    recovery.Recovery(),
    tracing.Server(),
    metrics.Server(), // Metrics 中间件
)
```

**无需额外配置，中间件会自动记录所有请求的指标！**

#### 自定义：选择性地记录某些路由

如果需要排除某些路由（如健康检查、监控端点），可以使用 `ServerWithConfig`：

```go
import (
    "context"
    "sre/internal/metrics"
    "github.com/go-kratos/kratos/v2/transport"
    "github.com/go-kratos/kratos/v2/transport/http"
)

// 在 internal/server/http.go 中
globalChain.Add(metrics.ServerWithConfig(metrics.MetricsConfig{
    SkipFunc: func(ctx context.Context) bool {
        tr, ok := transport.FromServerContext(ctx)
        if !ok {
            return false
        }
        if httpTr, ok := tr.(*http.Transport); ok {
            path := httpTr.Request().URL.Path
            // 排除健康检查和监控端点
            skipPaths := []string{
                "/health",
                "/metrics",
                "/ready",
            }
            for _, skipPath := range skipPaths {
                if path == skipPath {
                    return true // 跳过记录
                }
            }
        }
        return false // 记录指标
    },
}))
```

**gRPC 示例：**

```go
// 排除某些 gRPC 方法
globalChain.Add(metrics.ServerWithConfig(metrics.MetricsConfig{
    SkipFunc: func(ctx context.Context) bool {
        tr, ok := transport.FromServerContext(ctx)
        if !ok {
            return false
        }
        if grpcTr, ok := tr.(*grpc.Transport); ok {
            method := grpcTr.Operation()
            // 排除健康检查方法
            if strings.Contains(method, "Health") {
                return true
            }
        }
        return false
    },
}))
```

### 2. 手动记录业务指标

如果需要记录业务特定的指标，可以手动创建和记录：

```go
package main

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

func main() {
    // MeterProvider 已在 main.go 中自动初始化
    meter := otel.Meter("my-service")
    
    // 创建 Counter
    counter, _ := meter.Int64Counter("requests_total")
    
    // 记录指标
    ctx := context.Background()
    counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("method", "GET"),
    ))
}
```

### 3. 在 HTTP 处理器中使用

```go
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    meter := otel.Meter("user-service")
    requestCounter, _ := meter.Int64Counter("user_requests_total")
    requestDuration, _ := meter.Float64Histogram("user_request_duration_seconds")
    
    start := time.Now()
    user, err := h.repo.GetUser(ctx, req.Id)
    duration := time.Since(start)
    
    status := "success"
    if err != nil {
        status = "error"
    }
    
    requestCounter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("method", "GetUser"),
        attribute.String("status", status),
    ))
    requestDuration.Record(ctx, duration.Seconds())
    
    return user, err
}
```

## 指标类型

### Counter（计数器）

用于累计值，只能增加：

```go
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("method", "GET"),
    attribute.String("status", "200"),
))
```

### Gauge（仪表盘）

用于当前值，可以增加或减少：

```go
gauge, _ := meter.Int64ObservableGauge("active_connections")

_, _ = meter.RegisterCallback(
    func(ctx context.Context, o metric.Observer) error {
        count := getActiveConnectionCount()
        o.ObserveInt64(gauge, count)
        return nil
    },
    gauge,
)
```

### Histogram（直方图）

用于记录值的分布：

```go
histogram, _ := meter.Float64Histogram("request_duration_seconds")
histogram.Record(ctx, duration.Seconds(), metric.WithAttributes(
    attribute.String("method", "GET"),
))
```

## 导出器配置

### Prometheus 导出器

```yaml
metrics:
  prometheus_endpoint: ":9090/metrics"
```

**注意**：Prometheus 导出器需要框架提供 HTTP 端点暴露 `/metrics`，通常由 Kratos 的 metrics 中间件处理。

### OTLP 导出器

```yaml
metrics:
  otlp_endpoint: "localhost:4317"
```

适用于需要发送到 OpenTelemetry Collector 的场景。

### JSON File 导出器

```yaml
metrics:
  json_file_path: "./metrics/metrics.jsonl"
```

适用于开发、调试或本地存储的场景。

## 最佳实践

### 1. 指标命名

遵循 OpenTelemetry 和 Prometheus 的命名规范：

- 使用小写字母和下划线：`http_requests_total`
- 包含单位：`request_duration_seconds`、`memory_usage_bytes`
- 使用后缀表示类型：
  - Counter: `_total`、`_count`
  - Gauge: 无后缀
  - Histogram: `_seconds`、`_bytes` 等

### 2. 标签使用

- **使用标签区分维度**：通过标签区分不同的维度，而不是创建大量指标
- **避免高基数标签**：不要使用用户ID、请求ID等高基数值作为标签
- **常用标签**：
  - HTTP: `method`、`path`、`status`
  - 数据库: `operation`、`table`、`status`
  - 业务: `service`、`operation`、`status`

### 3. 性能考虑

- **异步记录**：指标记录是异步的，不会阻塞业务逻辑
- **批量导出**：使用 PeriodicReader 批量导出，减少网络开销
- **采样**：对于高频指标，考虑采样

### 4. 错误处理

```go
counter, err := meter.Int64Counter("requests_total")
if err != nil {
    log.Errorf("failed to create counter: %v", err)
    return
}
```

## 与 Tracing 的配合

Metrics 和 Traces 可以配合使用，提供完整的可观测性：

```go
func handleRequest(ctx context.Context, req *Request) error {
    // Tracing: 创建 span
    tracer := otel.Tracer("service")
    ctx, span := tracer.Start(ctx, "handleRequest")
    defer span.End()
    
    // Metrics: 记录指标
    meter := otel.Meter("service")
    counter, _ := meter.Int64Counter("requests_total")
    duration, _ := meter.Float64Histogram("request_duration_seconds")
    
    start := time.Now()
    err := processRequest(ctx, req)
    elapsed := time.Since(start)
    
    status := "success"
    if err != nil {
        status = "error"
        span.RecordError(err) // Tracing: 记录错误
    }
    counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("status", status),
    ))
    duration.Record(ctx, elapsed.Seconds())
    
    return err
}
```

## 参考资源

- [OpenTelemetry Metrics 文档](https://opentelemetry.io/docs/specs/otel/metrics/)
- [Prometheus 最佳实践](https://prometheus.io/docs/practices/naming/)
- [Kratos Metrics 中间件](https://github.com/go-kratos/kratos/tree/main/contrib/metrics)
- [项目 Metrics README](../internal/metrics/README.md)

