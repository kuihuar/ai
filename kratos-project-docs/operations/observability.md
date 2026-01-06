# 监控与可观测性

## 可观测性三大支柱

### 1. Metrics（指标）
- **用途**：量化系统状态
- **特点**：数值型数据，适合监控和告警
- **示例**：请求数、错误率、响应时间

### 2. Traces（追踪）
- **用途**：跟踪请求在系统中的流转
- **特点**：分布式追踪，了解请求路径
- **示例**：请求从入口到数据库的完整路径

### 3. Logs（日志）
- **用途**：记录系统事件
- **特点**：文本型数据，包含详细信息
- **示例**：错误日志、访问日志、业务日志

## Kratos 可观测性

### Metrics

Kratos 支持 Prometheus 指标：

```go
import "github.com/go-kratos/kratos/v2/metrics"

// 注册指标
counter := metrics.NewCounter("requests_total")
histogram := metrics.NewHistogram("request_duration_seconds")

// 记录指标
counter.Inc()
histogram.Observe(duration.Seconds())
```

### Traces

Kratos 支持 OpenTelemetry 追踪：

```go
import "go.opentelemetry.io/otel"

// 创建 tracer
tracer := otel.Tracer("service-name")

// 创建 span
ctx, span := tracer.Start(ctx, "operation-name")
defer span.End()

// 添加属性
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("request.size", size),
)
```

### Logs

使用结构化日志：

```go
import "github.com/go-kratos/kratos/v2/log"

logger := log.NewHelper(log.FromContext(ctx))
logger.Infow("request_processed",
    "method", "GET",
    "path", "/api/users",
    "status", 200,
    "duration_ms", 150,
)
```

## 监控指标

### 业务指标
- 用户注册数
- 订单创建数
- API 调用次数

### 技术指标
- 请求 QPS
- 错误率
- 响应时间（P50, P95, P99）
- 数据库连接数
- 缓存命中率

## 告警策略

### 告警规则
1. **错误率告警**：错误率 > 1%
2. **响应时间告警**：P95 响应时间 > 1s
3. **可用性告警**：服务可用性 < 99.9%
4. **资源告警**：CPU/内存使用率 > 80%

### 告警级别
- **Critical**：服务不可用，立即处理
- **Warning**：性能下降，需要关注
- **Info**：信息通知，了解即可

## 最佳实践

1. **指标命名**：使用统一的命名规范
2. **采样策略**：合理设置采样率，平衡性能和可观测性
3. **日志级别**：合理使用日志级别，避免日志过多
4. **告警收敛**：避免告警风暴，合理设置告警规则
5. **可视化**：使用 Grafana 等工具可视化指标

## 工具链

- **Prometheus**：指标收集和存储
- **Grafana**：指标可视化
- **Jaeger/Zipkin**：分布式追踪
- **ELK/Loki**：日志聚合和分析

