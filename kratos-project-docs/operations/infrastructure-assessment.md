# 基础设施评估与改进建议

## 概述

本文档对当前项目的基础设施进行全面评估，识别存在的问题和改进机会，并提供具体的改进建议。

## 当前基础设施状态

### ✅ 已实现的功能

1. **可观测性**
   - ✅ OpenTelemetry Tracing（支持 Jaeger、Zipkin、OTLP、JSON File）
   - ✅ OpenTelemetry Metrics（支持 Prometheus、OTLP、JSON File）
   - ✅ 结构化日志（Zap + Kratos Logger）
   - ✅ 日志轮转和清理策略（lumberjack）

2. **数据存储**
   - ✅ MySQL 数据库（GORM + Ent）
   - ✅ Redis（连接池已配置）
   - ✅ Kafka（已集成，但默认禁用）
   - ✅ 数据库连接池配置（可配置化）

3. **服务框架**
   - ✅ Kratos 微服务框架
   - ✅ HTTP/gRPC 双协议支持
   - ✅ 中间件链（Recovery、Tracing、Metrics、Auth、RateLimit）

4. **配置管理**
   - ✅ Viper 配置系统
   - ✅ 配置中心支持（Nacos、Apollo、Consul、Etcd）
   - ✅ 环境变量支持

5. **容器化**
   - ✅ Dockerfile（多阶段构建优化）
   - ✅ Docker Compose（Jaeger + MySQL）

6. **依赖注入**
   - ✅ Wire 依赖注入
   - ✅ 优雅关闭（Cleanup 函数）

## ⚠️ 需要改进的方面

### 1. 健康检查（Health Check）✅ 已完成

**状态**：
- ✅ 已实现健康检查端点（`/health`）
- ✅ 已实现就绪探针（`/ready`）
- ✅ 已实现存活探针（`/live`）
- ✅ 支持 HTTP 和 gRPC 两种协议

**实现位置**：
- `api/health/v1/health.proto` - Proto 定义
- `internal/service/health.go` - 服务实现
- `internal/data/health.go` - 健康检查逻辑
- `internal/server/http.go` - HTTP 端点注册
- `internal/server/grpc.go` - gRPC 端点注册

**文档**：`docs/operations/health-check.md`

**建议**（已实现）：

#### 1.1 实现健康检查服务

```go
// internal/service/health.go
package service

import (
    "context"
    "sre/api/health/v1"
    "sre/internal/data"
)

type HealthService struct {
    v1.UnimplementedHealthServer
    data *data.Data
}

func (s *HealthService) Check(ctx context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
    status := v1.HealthCheckResponse_SERVING
    
    // 检查数据库连接
    if s.data.DB() != nil {
        sqlDB, err := s.data.DB().DB()
        if err != nil || sqlDB.PingContext(ctx) != nil {
            status = v1.HealthCheckResponse_NOT_SERVING
        }
    }
    
    // 检查 Redis 连接
    if s.data.Redis() != nil {
        if err := s.data.Redis().Ping(ctx).Err(); err != nil {
            status = v1.HealthCheckResponse_NOT_SERVING
        }
    }
    
    return &v1.HealthCheckResponse{
        Status: status,
    }, nil
}
```

#### 1.2 添加 HTTP 健康检查端点

```go
// internal/server/http.go
import "github.com/go-kratos/kratos/v2/transport/http"

// 注册健康检查路由
httpSrv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    // 执行健康检查
    status := healthService.Check(r.Context(), &v1.HealthCheckRequest{})
    if status.Status == v1.HealthCheckResponse_SERVING {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("NOT_SERVING"))
    }
})

// 就绪探针
httpSrv.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
    // 检查服务是否就绪（数据库、Redis 等是否连接）
    // ...
})

// 存活探针
httpSrv.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
    // 检查服务是否存活（进程是否运行）
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
```

#### 1.3 使用 gRPC 健康检查协议

```go
// 使用标准的 gRPC 健康检查
import "google.golang.org/grpc/health"
import "google.golang.org/grpc/health/grpc_health_v1"

healthServer := health.NewServer()
grpc_health_v1.RegisterHealthServer(grpcSrv, healthServer)

// 设置服务状态
healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
```

### 2. 数据库连接池配置 ⭐ 高优先级

**现状**：
- ❌ 没有配置数据库连接池参数
- ❌ 使用 GORM 默认连接池设置
- ❌ 可能导致连接泄漏或性能问题

**建议**：

```go
// internal/data/data.go
func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
    // ... 现有代码 ...
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })
    if err != nil {
        return nil, err
    }
    
    // 配置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // 设置连接池参数
    sqlDB.SetMaxOpenConns(25)        // 最大打开连接数
    sqlDB.SetMaxIdleConns(10)        // 最大空闲连接数
    sqlDB.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生存时间
    sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大生存时间
    
    // 测试连接
    if err := sqlDB.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

**配置化**：

```yaml
# configs/config.yaml
data:
  database:
    pool:
      max_open_conns: 25
      max_idle_conns: 10
      conn_max_lifetime: 5m
      conn_max_idle_time: 10m
```

### 3. 环境变量管理 ⭐ 中优先级

**现状**：
- ❌ 没有 `.env.example` 文件
- ❌ 敏感信息（如数据库密码）硬编码在配置文件中
- ❌ 缺少环境隔离配置

**建议**：

#### 3.1 创建 `.env.example`

```bash
# .env.example
# 数据库配置
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=test

# Redis 配置
REDIS_ADDR=127.0.0.1:6379

# 服务配置
SERVICE_NAME=sre
SERVICE_VERSION=v1.0.0
ENVIRONMENT=dev

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json

# Tracing 配置
TRACING_JSON_FILE_PATH=./logs/traces.jsonl
TRACING_SAMPLING_RATIO=1.0

# Metrics 配置
METRICS_JSON_FILE_PATH=./logs/metrics.jsonl
METRICS_EXPORT_INTERVAL=10s
```

#### 3.2 使用环境变量覆盖配置

```go
// internal/config/kratos.go
func LoadBootstrapFromViper(v *viper.Viper) (*conf.Bootstrap, error) {
    // 环境变量优先级最高
    v.SetEnvPrefix("SRE")
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // 例如：SRE_DATA_DATABASE_SOURCE 会覆盖 data.database.source
    // ...
}
```

#### 3.3 敏感信息加密存储

```yaml
# configs/config.yaml
data:
  database:
    # 使用环境变量或加密存储
    source: ${DB_DSN}  # 从环境变量读取
    # 或使用加密的 DSN
    encrypted_dsn: ${ENCRYPTED_DB_DSN}
    decrypt_key: ${DB_DECRYPT_KEY}
```

### 4. Docker 镜像优化 ✅ 已完成

**状态**：
- ✅ 已实现多阶段构建
- ✅ 镜像体积优化（从 ~500MB 减少到 ~20MB）
- ✅ 包含健康检查
- ✅ 安全性优化（不包含构建工具）

**实现位置**：
- `Dockerfile` - 多阶段构建配置

**优化效果**：
- 镜像体积：从 ~500MB 减少到 ~20MB
- 构建时间：更快（分离构建和运行环境）
- 安全性：更好（不包含构建工具）
- 功能：包含健康检查、时区配置、日志目录

**建议**（已实现）：

### 5. Kubernetes 部署配置 ⭐ 高优先级

**现状**：
- ❌ 没有 Kubernetes 部署配置
- ❌ 没有 Service、Deployment、ConfigMap 等资源定义
- ❌ 无法在 K8s 环境中运行

**建议**：

#### 5.1 创建 Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sre
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sre
  template:
    metadata:
      labels:
        app: sre
    spec:
      containers:
      - name: sre
        image: sre:latest
        ports:
        - containerPort: 8000
          name: http
        - containerPort: 8989
          name: grpc
        env:
        - name: ENVIRONMENT
          value: "prod"
        - name: LOG_LEVEL
          value: "info"
        volumeMounts:
        - name: config
          mountPath: /app/configs
        livenessProbe:
          httpGet:
            path: /live
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: config
        configMap:
          name: sre-config
```

#### 5.2 创建 Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: sre
spec:
  selector:
    app: sre
  ports:
  - name: http
    port: 8000
    targetPort: 8000
  - name: grpc
    port: 8989
    targetPort: 8989
  type: ClusterIP
```

#### 5.3 创建 ConfigMap

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sre-config
data:
  config.yaml: |
    server:
      http:
        addr: 0.0.0.0:8000
      grpc:
        addr: 0.0.0.0:8989
    # ... 其他配置 ...
```

### 6. Prometheus Metrics 端点 ⭐ 中优先级

**现状**：
- ✅ 已实现 OpenTelemetry Metrics
- ❌ 没有暴露 Prometheus `/metrics` 端点
- ❌ 无法被 Prometheus 抓取

**建议**：

```go
// internal/server/http.go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// 注册 Prometheus metrics 端点
httpSrv.Handle("/metrics", promhttp.Handler())
```

**配置 Prometheus**：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'sre'
    static_configs:
      - targets: ['sre:8000']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### 7. 日志轮转和清理 ✅ 已完成

**状态**：
- ✅ 已实现日志轮转功能（使用 lumberjack）
- ✅ 已实现日志清理策略
- ✅ 支持配置化调整轮转参数
- ✅ 支持自动压缩旧日志文件

**实现位置**：
- `internal/conf/conf.proto` - 配置定义
- `configs/config.yaml` - 配置文件示例
- `internal/config/kratos.go` - 配置加载
- `internal/logger/zap.go` - 日志轮转实现
- `internal/logger/provider.go` - 配置解析

**配置参数**：
- `enable`: 是否启用日志轮转（默认 false）
- `max_size`: 每个日志文件最大大小（MB，默认 100）
- `max_backups`: 保留的备份文件数量（默认 10）
- `max_age`: 保留天数（默认 30）
- `compress`: 是否压缩旧日志文件（默认 true）
- `local_time`: 使用本地时间而非 UTC（默认 true）

**文档**：`docs/operations/log-rotation.md`

**建议**（已实现）：

### 8. CI/CD 流水线 ⭐ 中优先级

**现状**：
- ❌ 没有 CI/CD 配置
- ❌ 没有自动化测试
- ❌ 没有自动化构建和部署

**建议**：

#### 8.1 GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - name: Run tests
      run: go test ./...
    - name: Build
      run: make build
    - name: Lint
      run: golangci-lint run
```

#### 8.2 GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

test:
  stage: test
  script:
    - go test ./...
    - make build

build:
  stage: build
  script:
    - docker build -t sre:$CI_COMMIT_SHA .
    - docker push sre:$CI_COMMIT_SHA
```

### 9. 配置中心集成 ⭐ 低优先级

**现状**：
- ✅ 已支持多种配置中心（Nacos、Apollo、Consul、Etcd）
- ❌ 配置中心配置被注释，未启用
- ❌ 缺少配置中心使用文档

**建议**：

#### 9.1 启用配置中心

```yaml
# configs/config.yaml
config_center:
  nacos:
    endpoints:
      - "127.0.0.1:8848"
    namespace: "public"
    group: "DEFAULT_GROUP"
    data_id: "sre-config.yaml"
    username: "nacos"
    password: "nacos"
```

#### 9.2 环境隔离

```yaml
# 开发环境
config_center:
  nacos:
    namespace: "dev"
    data_id: "sre-config-dev.yaml"

# 生产环境
config_center:
  nacos:
    namespace: "prod"
    data_id: "sre-config-prod.yaml"
```

### 10. 监控告警 ⭐ 中优先级

**现状**：
- ✅ 已实现 Metrics 和 Tracing
- ❌ 没有告警规则
- ❌ 没有监控面板（Grafana）

**建议**：

#### 10.1 Prometheus 告警规则

```yaml
# prometheus/alerts.yml
groups:
  - name: sre
    rules:
      - alert: HighErrorRate
        expr: rate(http_server_requests_total{status="error"}[5m]) > 0.01
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          
      - alert: HighLatency
        expr: histogram_quantile(0.95, http_server_request_duration_seconds_bucket) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
```

#### 10.2 Grafana 仪表板

创建 Grafana 仪表板配置文件，监控：
- QPS（每秒请求数）
- 错误率
- 响应时间（P50、P95、P99）
- 数据库连接数
- Redis 连接数

## 优先级总结

### ✅ 已完成

1. **健康检查端点** ✅ - 已实现 HTTP 和 gRPC 健康检查
2. **数据库连接池配置** ✅ - 已实现配置化连接池管理
3. **Docker 镜像优化** ✅ - 已实现多阶段构建优化
4. **日志轮转和清理** ✅ - 已实现日志轮转和自动清理

### 🔴 高优先级（必须实现）

1. **Kubernetes 部署配置** - 影响生产环境部署

### 🟡 中优先级（建议实现）

2. **环境变量管理** - 提高配置灵活性
3. **Prometheus Metrics 端点** - 完善监控
4. **CI/CD 流水线** - 提高开发效率
5. **监控告警** - 及时发现问题

### 🟢 低优先级（可选）

6. **配置中心集成** - 根据实际需求决定

## 实施计划

### ✅ 第一阶段（已完成）

1. ✅ 实现健康检查端点
2. ✅ 配置数据库连接池
3. ✅ 优化 Docker 镜像
4. ✅ 实现日志轮转

### 🔄 第二阶段（进行中）

1. 创建 Kubernetes 部署配置
2. 添加 Prometheus Metrics 端点
3. 环境变量管理优化

### 📋 第三阶段（待实施）

1. 设置 CI/CD 流水线
2. 配置监控告警
3. 完善文档

## 参考资源

- [Kubernetes 健康检查](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [Prometheus 最佳实践](https://prometheus.io/docs/practices/)
- [Docker 多阶段构建](https://docs.docker.com/build/building/multi-stage/)
- [Grafana 仪表板](https://grafana.com/docs/grafana/latest/dashboards/)

