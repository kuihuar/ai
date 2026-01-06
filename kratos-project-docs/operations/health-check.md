# 健康检查端点实现说明

## 概述

项目已完整实现健康检查、就绪探针和存活探针功能，支持 HTTP 和 gRPC 两种协议。

## 已实现的端点

### HTTP 端点

1. **`GET /health`** - 健康检查
   - 检查所有依赖组件的健康状态（数据库、Redis 等）
   - 返回详细的状态信息

2. **`GET /ready`** - 就绪探针（Readiness Probe）
   - 检查服务是否准备好接收流量
   - 只检查关键依赖（如数据库）
   - 用于 Kubernetes 就绪探针

3. **`GET /live`** - 存活探针（Liveness Probe）
   - 检查服务进程是否存活
   - 不检查依赖，只确认进程运行
   - 用于 Kubernetes 存活探针

### gRPC 端点

1. **`api.health.v1.Health/Check`** - 健康检查
2. **`api.health.v1.Health/Readiness`** - 就绪检查
3. **`api.health.v1.Health/Liveness`** - 存活检查

同时注册了标准的 gRPC 健康检查服务（`grpc.health.v1.Health`），用于兼容性。

## 实现细节

### 1. Proto 定义

**文件**：`api/health/v1/health.proto`

```protobuf
service Health {
  rpc Check(HealthCheckRequest) returns (HealthCheckResponse) {
    option (google.api.http) = { get: "/health" };
  }
  
  rpc Readiness(ReadinessCheckRequest) returns (ReadinessCheckResponse) {
    option (google.api.http) = { get: "/ready" };
  }
  
  rpc Liveness(LivenessCheckRequest) returns (LivenessCheckResponse) {
    option (google.api.http) = { get: "/live" };
  }
}
```

### 2. 服务实现

**文件**：`internal/service/health.go`

```go
type HealthService struct {
    healthv1.UnimplementedHealthServer
    data   *data.Data
    logger log.Logger
}

// Check - 健康检查
func (s *HealthService) Check(ctx context.Context, req *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error)

// Readiness - 就绪检查
func (s *HealthService) Readiness(ctx context.Context, req *healthv1.ReadinessCheckRequest) (*healthv1.ReadinessCheckResponse, error)

// Liveness - 存活检查
func (s *HealthService) Liveness(ctx context.Context, req *healthv1.LivenessCheckRequest) (*healthv1.LivenessCheckResponse, error)
```

### 3. 数据层健康检查逻辑

**文件**：`internal/data/health.go`

#### HealthCheck（健康检查）

- 检查数据库连接（如果配置）
- 检查 Redis 连接（如果配置）
- 返回所有依赖组件的状态
- 包含连接池统计信息

#### ReadinessCheck（就绪检查）

- **关键依赖**：数据库必须可用
- **非关键依赖**：Redis 不可用不影响就绪状态
- 用于判断服务是否可以接收流量

#### LivenessCheck（存活检查）

- 只检查进程是否运行
- 不检查任何依赖
- 用于判断服务是否需要重启

## 响应格式

### HTTP 响应示例

#### `/health` 响应

```json
{
  "status": "SERVING",
  "message": "healthy",
  "details": {
    "database": "ok",
    "database_open_conns": "5",
    "database_idle_conns": "3",
    "redis": "ok"
  }
}
```

#### `/ready` 响应

```json
{
  "ready": true,
  "message": "ready",
  "details": {
    "database": "ready",
    "redis": "ready"
  }
}
```

#### `/live` 响应

```json
{
  "alive": true,
  "message": "alive"
}
```

### 状态码

- **HTTP 200**：服务健康/就绪/存活
- **HTTP 503**：服务不健康/未就绪（仅健康检查和就绪检查）

## 中间件配置

健康检查端点已自动排除以下中间件，避免产生不必要的追踪和指标：

- **Tracing 中间件**：排除 `/health`、`/ready`、`/live`
- **Metrics 中间件**：排除 `/health`、`/ready`、`/live`

**配置位置**：`internal/server/http.go`

```go
globalChain.Add(tracing.ServerWithConfig(tracing.TracingConfig{
    SkipFunc: func(ctx context.Context) bool {
        // 排除健康检查端点
        return path == "/health" || path == "/ready" || path == "/live"
    },
}))
```

## Kubernetes 配置示例

### Deployment 配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sre
spec:
  template:
    spec:
      containers:
      - name: sre
        image: sre:latest
        ports:
        - containerPort: 8000
          name: http
        # 存活探针
        livenessProbe:
          httpGet:
            path: /live
            port: 8000
          initialDelaySeconds: 30  # 启动后 30 秒开始检查
          periodSeconds: 10         # 每 10 秒检查一次
          timeoutSeconds: 3         # 超时 3 秒
          failureThreshold: 3       # 失败 3 次后重启
        # 就绪探针
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 10  # 启动后 10 秒开始检查
          periodSeconds: 5          # 每 5 秒检查一次
          timeoutSeconds: 2         # 超时 2 秒
          failureThreshold: 3       # 失败 3 次后标记为未就绪
```

### gRPC 健康检查（可选）

```yaml
livenessProbe:
  grpc:
    port: 8989
    service: api.health.v1.Health
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  grpc:
    port: 8989
    service: api.health.v1.Health
  initialDelaySeconds: 10
  periodSeconds: 5
```

## 测试

### 使用 curl 测试

```bash
# 健康检查
curl http://localhost:8000/health

# 就绪检查
curl http://localhost:8000/ready

# 存活检查
curl http://localhost:8000/live
```

### 使用 gRPC 测试

```bash
# 使用 grpcurl 工具
grpcurl -plaintext localhost:8989 api.health.v1.Health/Check
grpcurl -plaintext localhost:8989 api.health.v1.Health/Readiness
grpcurl -plaintext localhost:8989 api.health.v1.Health/Liveness
```

## 监控和告警

### Prometheus 监控

健康检查端点可以用于 Prometheus 监控：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'sre-health'
    metrics_path: '/health'
    static_configs:
      - targets: ['sre:8000']
```

### 告警规则

```yaml
# prometheus/alerts.yml
groups:
  - name: sre_health
    rules:
      - alert: ServiceUnhealthy
        expr: up{job="sre-health"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is unhealthy"
```

## 最佳实践

### 1. 探针配置建议

- **存活探针**：
  - `initialDelaySeconds`: 30-60 秒（给服务启动时间）
  - `periodSeconds`: 10 秒
  - `timeoutSeconds`: 3 秒
  - `failureThreshold`: 3（失败 3 次后重启）

- **就绪探针**：
  - `initialDelaySeconds`: 10-30 秒
  - `periodSeconds`: 5 秒
  - `timeoutSeconds`: 2 秒
  - `failureThreshold`: 3（失败 3 次后从负载均衡移除）

### 2. 健康检查逻辑

- **健康检查**：检查所有依赖，用于监控和诊断
- **就绪检查**：只检查关键依赖，用于流量控制
- **存活检查**：不检查依赖，只确认进程运行

### 3. 性能考虑

- 健康检查端点已排除 Tracing 和 Metrics 中间件
- 使用超时上下文（2 秒）避免长时间阻塞
- 连接池检查使用轻量级的 Ping 操作

## 故障排查

### 问题 1：健康检查返回 503

**可能原因**：
- 数据库连接失败
- Redis 连接失败
- 网络问题

**排查步骤**：
1. 检查数据库连接配置
2. 检查 Redis 连接配置
3. 查看服务日志
4. 检查网络连接

### 问题 2：就绪探针一直失败

**可能原因**：
- 数据库未启动
- 数据库连接配置错误
- 数据库连接池耗尽

**排查步骤**：
1. 检查数据库服务状态
2. 验证数据库连接配置
3. 检查连接池配置
4. 查看数据库日志

### 问题 3：存活探针失败

**可能原因**：
- 服务进程崩溃
- 服务未启动
- 端口被占用

**排查步骤**：
1. 检查服务进程状态
2. 查看服务日志
3. 检查端口占用情况

## 参考资源

- [Kubernetes 健康检查文档](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [gRPC 健康检查协议](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
- [Kratos 健康检查示例](https://github.com/go-kratos/kratos/tree/main/examples)

