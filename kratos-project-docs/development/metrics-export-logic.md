# Metrics 导出逻辑详细说明

## 概述

本文档详细说明当前系统中 Metrics（指标）的收集、记录和导出逻辑，包括数据流、导出时机、数据格式等。

## 快速回答

**Q: 定时导出后，会清理已导出的数据吗？**  
A: **不会**。当前系统使用 `CumulativeTemporality`（累积时间性），导出后数据继续在内存中累积，不会清理。

**Q: 会不会重复导出？**  
A: **不会重复导出相同数据**。虽然每次导出都包含数据，但：
- **Counter 类型**：每次导出的是累积值（不断增长），例如第 1 次导出 10，第 2 次导出 25，第 3 次导出 40
- **Histogram 类型**：每次导出的是累积的统计信息（样本数、总和等都在增长）
- **Gauge 类型**：每次导出的是当前最新值

因此，虽然每次导出都有数据，但**值是不同的**（累积值在增长），不是重复的相同数据。

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                     请求处理流程                              │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  HTTP/gRPC 请求                                              │
│  └─> Metrics 中间件 (internal/metrics/middleware.go)        │
│      ├─> 记录请求开始时间                                     │
│      ├─> 增加活跃请求计数                                     │
│      ├─> 执行业务逻辑                                         │
│      ├─> 计算请求耗时                                         │
│      └─> 记录指标到 OpenTelemetry Meter                      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  OpenTelemetry MeterProvider                                 │
│  └─> 收集指标数据                                             │
│      └─> 使用 PeriodicReader 定期导出                        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  导出器 (Exporter)                                           │
│  ├─> Prometheus Exporter (HTTP 端点)                        │
│  ├─> OTLP Exporter (gRPC 推送)                              │
│  └─> JSON File Exporter (文件写入) ← 当前使用                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  输出目标                                                     │
│  └─> logs/metrics.jsonl (JSON Lines 格式)                   │
└─────────────────────────────────────────────────────────────┘
```

## 详细流程

### 1. 初始化阶段

#### 1.1 MeterProvider 初始化

**位置**：`cmd/sre/main.go`

```go
// 初始化 OpenTelemetry MeterProvider
if metricsCleanup, err := metrics.InitMeterProvider(ctx, bootstrap.Metrics, Name, Version, logger); err == nil && metricsCleanup != nil {
    defer metricsCleanup()
    // 初始化 metrics 中间件
    if err := metrics.InitMetricsMiddleware(); err != nil {
        log.NewHelper(logger).Warnf("failed to init metrics middleware: %v", err)
    }
}
```

**流程**：
1. 读取配置（`bootstrap.Metrics`）
2. 根据配置选择导出器（Prometheus/OTLP/JSON File）
3. 创建 `MeterProvider` 并设置为全局 MeterProvider
4. 创建 `PeriodicReader`，配置导出间隔（默认 10s）和超时（默认 5s）

**关键代码**：`internal/metrics/provider.go`

```go
// 创建 PeriodicReader（定期导出）
reader = metric.NewPeriodicReader(
    jsonExporter,
    metric.WithInterval(cfg.ExportInterval),  // 默认 10s
    metric.WithTimeout(cfg.ExportTimeout),      // 默认 5s
)

// 创建 MeterProvider
mp := metric.NewMeterProvider(
    metric.WithResource(res),
    metric.WithReader(reader),
)

// 设置为全局 MeterProvider
otel.SetMeterProvider(mp)
```

#### 1.2 Metrics 中间件初始化

**位置**：`internal/metrics/middleware.go`

```go
func InitMetricsMiddleware() error {
    meter := otel.Meter("kratos-server")
    
    // 创建 HTTP 指标
    httpRequestCounter, _ = meter.Int64Counter("http_server_requests_total", ...)
    httpRequestDuration, _ = meter.Float64Histogram("http_server_request_duration_seconds", ...)
    httpRequestSize, _ = meter.Int64Histogram("http_server_request_size_bytes", ...)
    httpResponseSize, _ = meter.Int64Histogram("http_server_response_size_bytes", ...)
    httpActiveRequests, _ = meter.Int64UpDownCounter("http_server_active_requests", ...)
    
    // 创建 gRPC 指标（类似）
    ...
    
    initialized = true
    return nil
}
```

**创建的指标类型**：

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `http_server_requests_total` | Counter | HTTP 请求总数 |
| `http_server_request_duration_seconds` | Histogram | HTTP 请求耗时分布 |
| `http_server_request_size_bytes` | Histogram | HTTP 请求大小 |
| `http_server_response_size_bytes` | Histogram | HTTP 响应大小 |
| `http_server_active_requests` | UpDownCounter | 当前活跃请求数 |
| `grpc_server_requests_total` | Counter | gRPC 请求总数 |
| `grpc_server_request_duration_seconds` | Histogram | gRPC 请求耗时分布 |
| `grpc_server_active_requests` | UpDownCounter | 当前活跃 gRPC 请求数 |

### 2. 请求处理阶段

#### 2.1 HTTP 请求处理

**位置**：`internal/metrics/middleware.go::handleHTTP`

**流程**：

```go
func handleHTTP(ctx context.Context, req interface{}, handler middleware.Handler) (interface{}, error) {
    // 1. 获取请求信息
    request := httpTr.Request()
    method := request.Method
    path := request.URL.Path
    
    // 2. 增加活跃请求数（请求开始时）
    httpActiveRequests.Add(ctx, 1)
    defer httpActiveRequests.Add(ctx, -1)  // 请求结束时减少
    
    // 3. 记录请求大小（如果有）
    if requestSize > 0 {
        httpRequestSize.Record(ctx, requestSize, ...)
    }
    
    // 4. 记录开始时间
    start := time.Now()
    
    // 5. 执行业务逻辑
    resp, err := handler(ctx, req)
    
    // 6. 计算耗时
    duration := time.Since(start).Seconds()
    
    // 7. 构建属性（标签）
    attrs := []attribute.KeyValue{
        attribute.String("http.method", method),
        attribute.String("http.route", path),
        attribute.Int("http.status_code", statusCode),
        attribute.String("status", status),  // "success" 或 "error"
    }
    
    // 8. 记录指标
    httpRequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
    httpRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
    
    return resp, err
}
```

**关键点**：
- 指标记录是**异步的**，不会阻塞请求处理
- 每个指标都带有**属性（标签）**，用于区分不同的维度
- 活跃请求数使用 `UpDownCounter`，在请求开始和结束时自动增减

#### 2.2 gRPC 请求处理

**位置**：`internal/metrics/middleware.go::handleGRPC`

**流程**：与 HTTP 类似，但使用 gRPC 特定的属性：
- `rpc.service`：服务名
- `rpc.method`：方法名
- `rpc.status_code`：gRPC 状态码

### 3. 数据收集阶段

#### 3.1 OpenTelemetry 内部收集

**机制**：
- OpenTelemetry SDK 使用**内存缓冲区**收集指标数据
- 指标数据在内存中**累积**，不会立即导出
- 使用**聚合**机制（Aggregation）处理相同标签的指标值

**聚合类型**：
- **Counter**：累加（Sum）
- **Histogram**：分桶统计（ExplicitBucketHistogram）
- **Gauge**：最新值（LastValue）
- **UpDownCounter**：累加（Sum）

#### 3.2 导出触发

**触发时机**：
1. **定期导出**：每 `export_interval`（默认 10 秒）自动触发一次
2. **强制刷新**：调用 `ForceFlush()` 时立即导出
3. **关闭时**：调用 `Shutdown()` 时导出所有剩余数据

**代码位置**：`internal/metrics/provider.go`

```go
reader = metric.NewPeriodicReader(
    jsonExporter,
    metric.WithInterval(cfg.ExportInterval),  // 每 10 秒触发一次
    metric.WithTimeout(cfg.ExportTimeout),      // 导出超时 5 秒
)
```

### 4. 数据导出阶段

#### 4.1 JSON File Exporter

**位置**：`internal/metrics/exporter/json_file_exporter.go`

**导出流程**：

```go
func (e *JSONFileExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
    // 1. 加锁（确保线程安全）
    e.mu.Lock()
    defer e.mu.Unlock()
    
    // 2. 提取服务名
    serviceName := e.extractServiceName(rm.Resource)
    
    // 3. 遍历所有 ScopeMetrics（作用域指标）
    for _, sm := range rm.ScopeMetrics {
        // 4. 遍历所有 Metrics（指标）
        for _, m := range sm.Metrics {
            // 5. 转换指标为 JSON 格式
            metricDataList := e.convertMetric(m, serviceName, rm.Resource)
            
            // 6. 写入文件（JSON Lines 格式）
            for _, metricData := range metricDataList {
                e.encoder.Encode(metricData)  // 每行一个 JSON 对象
            }
        }
    }
    
    // 7. 同步到磁盘
    e.file.Sync()
    
    return nil
}
```

#### 4.2 数据转换

**指标类型转换**：

1. **Counter (Sum)**
   ```json
   {
     "timestamp": "2025-12-16T09:30:17.594+0800",
     "service_name": "sre",
     "metric_name": "http_server_requests_total",
     "metric_type": "counter",
     "value": 1234,
     "unit": "1",
     "description": "Total number of HTTP requests",
     "attributes": {
       "http.method": "POST",
       "http.route": "/api/v1/orders",
       "http.status_code": 200,
       "status": "success"
     }
   }
   ```

2. **Histogram**
   ```json
   {
     "timestamp": "2025-12-16T09:30:17.594+0800",
     "service_name": "sre",
     "metric_name": "http_server_request_duration_seconds",
     "metric_type": "histogram",
     "value": {
       "count": 100,
       "sum": 5.234,
       "min": 0.001,
       "max": 0.5,
       "buckets": [
         {"upper_bound": 0.005, "count": 10},
         {"upper_bound": 0.01, "count": 20},
         {"upper_bound": 0.025, "count": 30},
         ...
       ]
     },
     "attributes": {
       "http.method": "POST",
       "http.route": "/api/v1/orders",
       "http.status_code": 200
     }
   }
   ```

3. **Gauge (UpDownCounter)**
   ```json
   {
     "timestamp": "2025-12-16T09:30:17.594+0800",
     "service_name": "sre",
     "metric_name": "http_server_active_requests",
     "metric_type": "gauge",
     "value": 5,
     "unit": "1",
     "description": "Number of active HTTP requests"
   }
   ```

#### 4.3 文件写入

**文件格式**：JSON Lines（每行一个 JSON 对象）

**文件路径**：`./logs/metrics.jsonl`（配置文件中指定）

**写入模式**：
- **追加模式**：`os.O_CREATE|os.O_WRONLY|os.O_APPEND`
- **自动创建目录**：如果目录不存在，自动创建
- **同步写入**：每次导出后调用 `file.Sync()` 确保数据写入磁盘

**线程安全**：使用 `sync.Mutex` 确保并发安全

### 5. 导出时机详解

#### 5.1 定期导出

**默认间隔**：10 秒

**配置位置**：`configs/config.yaml`

```yaml
metrics:
  export_interval: 10s  # 每 10 秒导出一次
  export_timeout: 5s    # 导出超时时间
```

**工作原理**：
1. `PeriodicReader` 内部维护一个定时器
2. 每 10 秒触发一次 `Export()` 调用
3. 收集当前内存中的所有指标数据
4. 调用 `Exporter.Export()` 导出数据
5. 根据 Temporality 类型决定是否重置数据

#### 5.2 Temporality（时间性）机制

**当前配置**：`CumulativeTemporality`（累积时间性）

**位置**：`internal/metrics/exporter/json_file_exporter.go`

```go
func (e *JSONFileExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
    return metricdata.CumulativeTemporality  // 所有指标类型都使用累积时间性
}
```

**CumulativeTemporality 的特点**：

1. **累积值**：每次导出的是**从应用启动到现在的累积值**
   - Counter：总请求数（从启动到现在的总和）
   - Histogram：总样本数、总和、最小值、最大值（从启动到现在）
   - Gauge：当前值（最新值）

2. **不会重置**：导出后**不会清理或重置**数据
   - **"不清理"的含义**：
     - 导出后，内存中的数据**不会被删除**
     - 数据**继续保留在内存中**
     - 下次导出时，会**继续累积**新的数据
   - 类比：就像一个**计数器**，每次导出时记录当前总数，但计数器本身不会清零，继续计数
   - 数据在内存中继续累积
   - 下次导出时，值会更大（对于 Counter）或更新（对于 Gauge）

3. **不会重复导出**：
   - 虽然每次导出都包含数据，但**值是不同的**（累积值在增长）
   - 例如：
     - 第 1 次导出：`http_server_requests_total = 10`
     - 第 2 次导出：`http_server_requests_total = 25`（包含前 10 个 + 新增 15 个）
     - 第 3 次导出：`http_server_requests_total = 40`（包含前 25 个 + 新增 15 个）

**示例时间线**：

```
时间    请求数    导出值（累积）    说明
00:00   0         -               应用启动
00:10   10        10               第 1 次导出：10 个请求
00:20   25        25               第 2 次导出：25 个请求（10 + 15）
00:30   40        40               第 3 次导出：40 个请求（25 + 15）
```

**实际文件内容示例**：

```jsonl
{"timestamp":"2025-12-16T09:30:10.000+0800","metric_name":"http_server_requests_total","value":10,"attributes":{"http.route":"/api/v1/orders"}}
{"timestamp":"2025-12-16T09:30:20.000+0800","metric_name":"http_server_requests_total","value":25,"attributes":{"http.route":"/api/v1/orders"}}
{"timestamp":"2025-12-16T09:30:30.000+0800","metric_name":"http_server_requests_total","value":40,"attributes":{"http.route":"/api/v1/orders"}}
```

**关键点**：
- 每次导出的值都在增长（10 → 25 → 40）
- 不是重复的相同数据（10 → 10 → 10）
- 每次导出都包含从启动到现在的所有数据（累积值）

**"数据清理"的详细说明**：

**Cumulative 模式（当前使用）- "不清理"**：

```
内存中的数据状态：
┌─────────────────────────────────────┐
│ 时间 00:00 - 应用启动               │
│ 内存：http_requests_total = 0      │
└─────────────────────────────────────┘
           │
           ▼ 处理 10 个请求
┌─────────────────────────────────────┐
│ 时间 00:10 - 第 1 次导出             │
│ 内存：http_requests_total = 10      │
│ 导出：10                            │
│ 导出后：http_requests_total = 10    │ ← 数据保留，不清理
└─────────────────────────────────────┘
           │
           ▼ 处理 15 个请求
┌─────────────────────────────────────┐
│ 时间 00:20 - 第 2 次导出             │
│ 内存：http_requests_total = 25      │ ← 继续累积（10 + 15）
│ 导出：25                            │
│ 导出后：http_requests_total = 25    │ ← 数据保留，不清理
└─────────────────────────────────────┘
```

**Delta 模式（未使用）- "清理"**：

```
内存中的数据状态：
┌─────────────────────────────────────┐
│ 时间 00:00 - 应用启动               │
│ 内存：http_requests_total = 0       │
└─────────────────────────────────────┘
           │
           ▼ 处理 10 个请求
┌─────────────────────────────────────┐
│ 时间 00:10 - 第 1 次导出             │
│ 内存：http_requests_total = 10      │
│ 导出：10（增量）                    │
│ 导出后：http_requests_total = 0     │ ← 数据被清理/重置
└─────────────────────────────────────┘
           │
           ▼ 处理 15 个请求
┌─────────────────────────────────────┐
│ 时间 00:20 - 第 2 次导出             │
│ 内存：http_requests_total = 15      │ ← 重新开始计数
│ 导出：15（增量）                    │
│ 导出后：http_requests_total = 0     │ ← 数据被清理/重置
└─────────────────────────────────────┘
```

**与 Delta Temporality 的对比**：

| 特性 | Cumulative（当前使用） | Delta（未使用） |
|------|----------------------|----------------|
| 导出值 | 累积值（从启动到现在） | 增量值（自上次导出） |
| 导出后 | 不重置，继续累积 | 重置为 0，重新开始 |
| 数据清理 | **不清理**：内存中的数据保留，继续累积 | **清理**：内存中的数据被重置为 0 |
| 内存占用 | 持续增长（累积值） | 保持稳定（只记录增量） |
| 重复导出 | 不会（值不同） | 不会（已清理） |
| 适用场景 | 长期监控、趋势分析 | 短期统计、实时监控 |

**"不清理"的具体含义**：

1. **内存中的数据不会被删除**：
   - 导出后，内存中的指标值**保持不变**
   - 例如：导出时 `http_requests_total = 10`，导出后内存中仍然是 `10`
   - 不会变成 `0` 或 `null`

2. **数据继续累积**：
   - 新的请求到来时，在现有值基础上**继续累加**
   - 例如：导出后内存中是 `10`，来了 5 个新请求，变成 `15`

3. **下次导出包含所有历史数据**：
   - 下次导出时，值会包含从启动到现在的**所有数据**
   - 例如：第 1 次导出 `10`，第 2 次导出 `25`（包含之前的 10 + 新增的 15）

**类比理解**：

- **Cumulative（不清理）**：像一个**总里程表**
  - 每次查看时显示总里程（从买车到现在）
  - 查看后里程表不会清零，继续累积
  - 下次查看时显示更大的数值

- **Delta（清理）**：像一个**单次行程表**
  - 每次查看时显示本次行程的里程
  - 查看后行程表清零，重新开始
  - 下次查看时只显示新的行程里程

**注意**：
- 指标数据在内存中**累积**，不会在导出后清理
- 每次导出的是**累积值**，不是增量值
- 虽然每次导出都包含数据，但**值在增长**，不是重复的相同数据
- 如果应用在导出间隔内退出，可能丢失部分数据
- 建议在应用关闭时调用 `cleanup()` 函数，确保数据导出

#### 5.2 导出超时

**默认超时**：5 秒

**作用**：如果导出操作超过 5 秒未完成，会取消操作并记录错误

**场景**：
- 文件写入速度慢
- 磁盘空间不足
- 网络问题（OTLP 导出器）

#### 5.3 关闭时导出

**位置**：`cmd/sre/main.go`

```go
defer metricsCleanup()  // 应用退出时调用
```

**实现**：`internal/metrics/provider.go`

```go
cleanup := func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := mp.Shutdown(ctx); err != nil {
        fmt.Printf("failed to shutdown MeterProvider: %v\n", err)
    }
}
```

**作用**：
- 调用 `MeterProvider.Shutdown()` 会触发最后一次导出
- 确保所有内存中的指标数据都被导出
- 关闭导出器，释放资源

### 6. 数据格式说明

#### 6.1 JSON Lines 格式

**特点**：
- 每行一个完整的 JSON 对象
- 使用换行符 `\n` 分隔
- 便于流式处理和大文件分析

**示例**：

```jsonl
{"timestamp":"2025-12-16T09:30:17.594+0800","service_name":"sre","metric_name":"http_server_requests_total","metric_type":"counter","value":1,"attributes":{"http.method":"POST","http.route":"/api/v1/orders","http.status_code":200,"status":"success"}}
{"timestamp":"2025-12-16T09:30:17.604+0800","service_name":"sre","metric_name":"http_server_request_duration_seconds","metric_type":"histogram","value":{"count":1,"sum":0.123,"min":0.123,"max":0.123,"buckets":[...]},"attributes":{"http.method":"POST","http.route":"/api/v1/orders"}}
{"timestamp":"2025-12-16T09:30:27.594+0800","service_name":"sre","metric_name":"http_server_requests_total","metric_type":"counter","value":5,"attributes":{"http.method":"GET","http.route":"/api/v1/orders","http.status_code":200,"status":"success"}}
```

#### 6.2 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `timestamp` | string | 指标记录时间（RFC3339Nano 格式） |
| `service_name` | string | 服务名称（从 Resource 中提取） |
| `metric_name` | string | 指标名称 |
| `metric_type` | string | 指标类型：`counter`、`gauge`、`histogram` |
| `value` | interface{} | 指标值（根据类型不同而不同） |
| `unit` | string | 单位（如 "1"、"s"、"By"） |
| `description` | string | 指标描述 |
| `attributes` | object | 标签（用于区分不同维度） |
| `resource` | object | 资源属性（服务元数据） |

#### 6.3 指标值格式

**Counter/Gauge**：
```json
"value": 1234
```

**Histogram**：
```json
"value": {
  "count": 100,      // 样本总数
  "sum": 5.234,      // 总和
  "min": 0.001,      // 最小值
  "max": 0.5,        // 最大值
  "buckets": [       // 分桶统计
    {"upper_bound": 0.005, "count": 10},
    {"upper_bound": 0.01, "count": 20},
    ...
  ]
}
```

### 7. 性能考虑

#### 7.1 异步处理

- 指标记录是**异步的**，不会阻塞业务逻辑
- 使用内存缓冲区，批量导出
- 导出操作在后台线程执行

#### 7.2 内存使用

- 指标数据在内存中累积，直到导出
- 内存占用取决于：
  - 指标数量
  - 标签组合数量（高基数标签会增加内存）
  - 导出间隔（间隔越长，累积数据越多）

#### 7.3 文件 I/O

- 使用**追加模式**写入，性能较好
- 每次导出后调用 `Sync()`，确保数据持久化
- 使用互斥锁，避免并发写入冲突

### 8. 配置说明

#### 8.1 完整配置示例

```yaml
metrics:
  service_name: "sre"                    # 服务名称
  service_version: "v1.0.0"             # 服务版本
  environment: "dev"                      # 环境
  json_file_path: "./logs/metrics.jsonl" # JSON 文件路径
  export_interval: 10s                   # 导出间隔
  export_timeout: 5s                     # 导出超时
```

#### 8.2 配置参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `service_name` | string | 应用名称 | 用于标识服务的指标 |
| `service_version` | string | 应用版本 | 服务版本号 |
| `environment` | string | "dev" | 环境标识 |
| `json_file_path` | string | - | JSON 文件路径（相对或绝对） |
| `export_interval` | duration | 10s | 导出间隔，建议 5-60 秒 |
| `export_timeout` | duration | 5s | 导出超时时间 |

#### 8.3 导出间隔建议

- **开发环境**：10-30 秒（便于调试）
- **测试环境**：10 秒（平衡性能和实时性）
- **生产环境**：30-60 秒（减少 I/O 开销）

### 9. 故障排查

#### 9.1 没有生成 metrics.jsonl 文件

**可能原因**：
1. 配置未加载：检查 `bootstrap.Metrics` 是否为 nil
2. 导出器未初始化：检查日志中是否有 "OpenTelemetry MeterProvider initialized"
3. 文件路径错误：检查路径是否正确，是否有写入权限
4. 没有请求：确保有 HTTP/gRPC 请求经过中间件

**排查步骤**：
```bash
# 1. 检查配置
grep -A 5 "=== DEBUG: Metrics config ===" logs/app.log

# 2. 检查初始化日志
grep "MeterProvider initialized" logs/app.log

# 3. 检查文件
ls -lh logs/metrics.jsonl

# 4. 检查文件内容
tail -f logs/metrics.jsonl
```

#### 9.2 文件内容为空或更新不及时

**可能原因**：
1. 导出间隔太长：检查 `export_interval` 配置
2. 没有请求：确保有请求经过中间件
3. 中间件未启用：检查 `initialized` 标志

**排查步骤**：
```bash
# 1. 检查导出间隔配置
grep "export_interval" configs/config.yaml

# 2. 等待一个导出间隔（默认 10 秒）后检查文件
sleep 11 && ls -lh logs/metrics.jsonl

# 3. 检查中间件是否初始化
grep "InitMetricsMiddleware" logs/app.log
```

#### 9.3 文件写入失败

**可能原因**：
1. 磁盘空间不足
2. 权限问题
3. 目录不存在（应该会自动创建）

**排查步骤**：
```bash
# 1. 检查磁盘空间
df -h

# 2. 检查文件权限
ls -l logs/metrics.jsonl

# 3. 检查目录权限
ls -ld logs/
```

### 10. 最佳实践

#### 10.1 标签使用

**✅ 好的实践**：
```go
// 使用低基数的标签
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("http.method", "POST"),
    attribute.String("http.route", "/api/v1/orders"),
    attribute.Int("http.status_code", 200),
))
```

**❌ 避免的做法**：
```go
// 不要使用高基数的标签（如用户ID、请求ID）
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.Int64("user_id", userID),  // 高基数！
    attribute.String("request_id", reqID), // 高基数！
))
```

#### 10.2 指标命名

**遵循规范**：
- 使用小写字母和下划线
- 包含单位后缀（`_total`、`_seconds`、`_bytes`）
- 使用描述性名称

**示例**：
- ✅ `http_server_requests_total`
- ✅ `http_server_request_duration_seconds`
- ❌ `httpRequests`（驼峰命名）
- ❌ `http_requests`（缺少类型后缀）

#### 10.3 导出间隔调优

**根据场景选择**：
- **高频服务**：较短的间隔（5-10 秒），及时发现问题
- **低频服务**：较长的间隔（30-60 秒），减少 I/O 开销
- **开发环境**：较短的间隔（5-10 秒），便于调试

### 11. 数据查看和分析

#### 11.1 实时查看

```bash
# 实时查看最新写入的指标
tail -f logs/metrics.jsonl

# 查看最后 10 条记录
tail -n 10 logs/metrics.jsonl
```

#### 11.2 统计分析

```bash
# 统计总请求数
jq -r 'select(.metric_name=="http_server_requests_total") | .value' logs/metrics.jsonl | awk '{sum+=$1} END {print sum}'

# 查看请求耗时分布
jq -r 'select(.metric_name=="http_server_request_duration_seconds") | .value.sum' logs/metrics.jsonl

# 按路由统计请求数
jq -r 'select(.metric_name=="http_server_requests_total") | "\(.attributes["http.route"]) \(.value)"' logs/metrics.jsonl | sort | uniq -c
```

#### 11.3 转换为 Prometheus 格式

可以使用工具将 JSON Lines 格式转换为 Prometheus 格式，或直接使用 Prometheus 导出器。

## 总结

当前系统的 Metrics 导出逻辑：

1. **收集**：Metrics 中间件自动记录所有 HTTP/gRPC 请求的指标
2. **累积**：指标数据在内存中累积，使用 OpenTelemetry SDK 的聚合机制
3. **导出**：每 10 秒（可配置）自动导出一次到 JSON 文件
4. **格式**：JSON Lines 格式，每行一个指标记录
5. **持久化**：每次导出后同步到磁盘，确保数据不丢失
6. **时间性**：使用 Cumulative（累积）模式，导出累积值，不清理数据

**关键特性**：
- ✅ 自动收集，无需手动编写代码
- ✅ 异步处理，不影响业务性能
- ✅ 批量导出，减少 I/O 开销
- ✅ 结构化数据，便于分析和查询
- ✅ 累积模式，不会重复导出相同数据（每次值不同）

**关于数据清理和重复导出**：

1. **不会清理数据**：使用 Cumulative Temporality，导出后数据继续在内存中累积
2. **不会重复导出**：虽然每次导出都包含数据，但值是累积的（不断增长），不是重复的相同数据
3. **内存管理**：数据在内存中累积，直到应用重启或关闭
4. **文件增长**：JSON 文件会持续增长，每次导出都会追加新行（累积值）

**如果需要增量导出**：

如果需要每次导出增量值（而不是累积值），可以修改 `Temporality()` 方法返回 `metricdata.DeltaTemporality`，这样：
- 每次导出的是自上次导出以来的增量
- 导出后数据会被重置
- 可以避免文件快速增长

