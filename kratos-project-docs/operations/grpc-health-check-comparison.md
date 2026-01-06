# 自定义 vs 标准 gRPC 健康检查服务对比

## 概述

gRPC 健康检查服务有两种实现方式：
1. **标准健康检查服务**：`grpc.health.v1.Health`（gRPC 官方定义）
2. **自定义健康检查服务**：`api.health.v1.Health`（项目自定义）

## 详细对比

### 1. 服务定义

#### 标准健康检查服务 (`grpc.health.v1.Health`)

**Proto 定义**（gRPC 官方提供）：

```protobuf
syntax = "proto3";

package grpc.health.v1;

service Health {
  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);
  rpc Watch(HealthCheckRequest) returns (stream HealthCheckResponse);
}

message HealthCheckRequest {
  string service = 1;
}

message HealthCheckResponse {
  enum ServingStatus {
    UNKNOWN = 0;
    SERVING = 1;
    NOT_SERVING = 2;
    SERVICE_UNKNOWN = 3;
  }
  ServingStatus status = 1;
}
```

**特点**：
- ✅ 服务名固定：`grpc.health.v1.Health`
- ✅ 只有 `Check` 和 `Watch` 两个方法
- ✅ 响应格式简单：只有 `status` 字段
- ✅ 不支持自定义消息和详细信息

#### 自定义健康检查服务 (`api.health.v1.Health`)

**Proto 定义**（项目自定义）：

```protobuf
syntax = "proto3";

package api.health.v1;

service Health {
  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);
  rpc Readiness(ReadinessCheckRequest) returns (ReadinessCheckResponse);
  rpc Liveness(LivenessCheckRequest) returns (LivenessCheckResponse);
}

message HealthCheckResponse {
  enum ServingStatus {
    UNKNOWN = 0;
    SERVING = 1;
    NOT_SERVING = 2;
    SERVICE_UNKNOWN = 3;
  }
  ServingStatus status = 1;
  string message = 2;                    // 自定义：状态消息
  map<string, string> details = 3;       // 自定义：详细信息
}

message ReadinessCheckResponse {
  bool ready = 1;
  string message = 2;
  map<string, string> details = 3;       // 各依赖组件的状态
}

message LivenessCheckResponse {
  bool alive = 1;
  string message = 2;
}
```

**特点**：
- ✅ 服务名自定义：`api.health.v1.Health`
- ✅ 三个独立方法：`Check`、`Readiness`、`Liveness`
- ✅ 响应格式丰富：包含 `message` 和 `details`
- ✅ 支持 HTTP 端点（通过 `google.api.http` 注解）

### 2. 功能对比

| 特性 | 标准服务 | 自定义服务 |
|------|---------|-----------|
| **服务名** | `grpc.health.v1.Health` | `api.health.v1.Health` |
| **方法数量** | 2 个（Check、Watch） | 3 个（Check、Readiness、Liveness） |
| **健康检查** | ✅ Check | ✅ Check |
| **就绪检查** | ❌ 无 | ✅ Readiness |
| **存活检查** | ❌ 无 | ✅ Liveness |
| **流式监控** | ✅ Watch | ❌ 无 |
| **状态消息** | ❌ 无 | ✅ message |
| **详细信息** | ❌ 无 | ✅ details |
| **HTTP 支持** | ❌ 无 | ✅ 支持（通过注解） |
| **Kubernetes 兼容** | ✅ 支持 | ✅ 支持 |

### 3. 使用场景

#### 标准服务适用场景

1. **简单健康检查**
   - 只需要检查服务是否运行
   - 不需要详细的诊断信息

2. **标准化要求**
   - 需要遵循 gRPC 官方标准
   - 需要与标准工具兼容

3. **流式监控**
   - 需要实时监控服务状态变化
   - 使用 `Watch` 方法进行流式检查

4. **负载均衡器集成**
   - 负载均衡器默认支持标准服务
   - 无需额外配置

#### 自定义服务适用场景

1. **Kubernetes 探针**
   - 需要独立的就绪探针（Readiness）
   - 需要独立的存活探针（Liveness）
   - 标准服务只有一个 `Check` 方法

2. **详细诊断信息**
   - 需要返回详细的健康状态信息
   - 需要包含各依赖组件的状态

3. **HTTP 和 gRPC 双协议**
   - 需要同时支持 HTTP 和 gRPC
   - 通过 `google.api.http` 注解自动生成 HTTP 端点

4. **业务特定检查**
   - 需要检查业务特定的健康状态
   - 需要自定义检查逻辑

### 4. 代码实现对比

#### 标准服务实现

```go
import (
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
)

// 创建标准健康检查服务
grpcHealthServer := health.NewServer()
grpc_health_v1.RegisterHealthServer(srv, grpcHealthServer)

// 设置服务状态
grpcHealthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
```

**特点**：
- 使用 gRPC 官方库
- 需要手动设置服务状态
- 状态管理简单

#### 自定义服务实现

```go
import (
    healthv1 "sre/api/health/v1"
    "sre/internal/data"
)

// 创建自定义健康检查服务
type HealthService struct {
    healthv1.UnimplementedHealthServer
    data   *data.Data
    logger log.Logger
}

// 实现 Check 方法
func (s *HealthService) Check(ctx context.Context, req *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
    status := s.data.HealthCheck(ctx)
    // 返回详细的状态信息
    return &healthv1.HealthCheckResponse{
        Status:  healthv1.HealthCheckResponse_SERVING,
        Message: status.Message,
        Details: status.Details,
    }, nil
}

// 实现 Readiness 方法
func (s *HealthService) Readiness(ctx context.Context, req *healthv1.ReadinessCheckRequest) (*healthv1.ReadinessCheckResponse, error) {
    // 检查服务是否就绪
}

// 实现 Liveness 方法
func (s *HealthService) Liveness(ctx context.Context, req *healthv1.LivenessCheckRequest) (*healthv1.LivenessCheckResponse, error) {
    // 检查服务是否存活
}
```

**特点**：
- 完全自定义实现
- 可以检查数据库、Redis 等依赖
- 返回详细的状态信息

### 5. Kubernetes 配置对比

#### 标准服务配置

```yaml
livenessProbe:
  grpc:
    port: 8989
    service: grpc.health.v1.Health  # 标准服务名
  initialDelaySeconds: 30

readinessProbe:
  grpc:
    port: 8989
    service: grpc.health.v1.Health  # 使用同一个服务
  initialDelaySeconds: 10
```

**问题**：
- 存活和就绪探针使用同一个 `Check` 方法
- 无法区分存活检查和就绪检查的逻辑
- 需要手动管理服务状态

#### 自定义服务配置

```yaml
livenessProbe:
  grpc:
    port: 8989
    service: api.health.v1.Health  # 自定义服务名
  initialDelaySeconds: 30

readinessProbe:
  grpc:
    port: 8989
    service: api.health.v1.Health  # 自定义服务名
  initialDelaySeconds: 10
```

**或者使用 HTTP 端点**：

```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 8000
  initialDelaySeconds: 30

readinessProbe:
  httpGet:
    path: /ready
    port: 8000
  initialDelaySeconds: 10
```

**优势**：
- 存活和就绪探针使用不同的方法
- 可以有不同的检查逻辑
- 支持 HTTP 和 gRPC 两种方式

### 6. 兼容性对比

#### 标准服务兼容性

| 工具/框架 | 支持情况 | 说明 |
|----------|---------|------|
| gRPC 官方工具 | ✅ 完全支持 | 默认支持标准服务 |
| Kubernetes | ✅ 支持 | 支持 gRPC 健康检查 |
| 负载均衡器 | ✅ 广泛支持 | 大多数负载均衡器支持 |
| HTTP 客户端 | ❌ 不支持 | 只支持 gRPC |

#### 自定义服务兼容性

| 工具/框架 | 支持情况 | 说明 |
|----------|---------|------|
| gRPC 官方工具 | ⚠️ 需要配置 | 需要指定服务名 |
| Kubernetes | ✅ 支持 | 支持自定义服务名 |
| 负载均衡器 | ⚠️ 需要配置 | 需要指定服务名 |
| HTTP 客户端 | ✅ 支持 | 通过 HTTP 端点 |

### 7. 性能对比

| 指标 | 标准服务 | 自定义服务 |
|------|---------|-----------|
| **响应时间** | 快（简单检查） | 稍慢（详细检查） |
| **资源占用** | 低 | 中等 |
| **网络开销** | 小 | 中等（包含详细信息） |

### 8. 维护成本

#### 标准服务

- ✅ 使用官方库，维护成本低
- ✅ 接口稳定，很少变化
- ❌ 功能有限，扩展困难

#### 自定义服务

- ⚠️ 需要自己实现和维护
- ✅ 功能灵活，易于扩展
- ⚠️ 需要维护 Proto 定义和实现代码

## 当前项目的选择

### 为什么选择自定义服务？

1. **Kubernetes 需求**
   - 需要独立的就绪探针和存活探针
   - 标准服务只有一个 `Check` 方法，无法区分

2. **详细诊断信息**
   - 需要返回数据库、Redis 等依赖的状态
   - 标准服务只返回简单的状态码

3. **HTTP 支持**
   - 需要同时支持 HTTP 和 gRPC
   - 通过 `google.api.http` 注解自动生成 HTTP 端点

4. **业务特定检查**
   - 需要检查业务特定的健康状态
   - 可以自定义检查逻辑

### 是否可以同时使用？

**不建议同时注册两个服务**，原因：

1. **服务名冲突**
   - 如果自定义服务也使用 `grpc.health.v1.Health` 作为服务名，会导致冲突
   - 错误：`duplicate service registration for "grpc.health.v1.Health"`

2. **功能重复**
   - 两个服务功能重叠
   - 增加维护成本

3. **资源浪费**
   - 两个服务都会占用资源
   - 增加代码复杂度

### 如何选择？

#### 选择标准服务，如果：

- ✅ 只需要简单的健康检查
- ✅ 需要与标准工具兼容
- ✅ 需要流式监控（Watch）
- ✅ 不需要详细的诊断信息

#### 选择自定义服务，如果：

- ✅ 需要 Kubernetes 就绪/存活探针
- ✅ 需要详细的诊断信息
- ✅ 需要 HTTP 和 gRPC 双协议支持
- ✅ 需要业务特定的检查逻辑

## 最佳实践建议

### 1. 单一服务原则

只注册一个健康检查服务，避免冲突和混淆。

### 2. 当前项目推荐

**使用自定义服务**（`api.health.v1.Health`），因为：
- ✅ 支持 Kubernetes 探针
- ✅ 提供详细的诊断信息
- ✅ 支持 HTTP 和 gRPC
- ✅ 功能更丰富

### 3. 如果需要标准服务

如果确实需要标准服务（例如与特定工具集成），可以：

1. **移除自定义服务**，只使用标准服务
2. **或者实现适配器**，让自定义服务同时实现标准接口

### 4. 混合方案（不推荐）

如果必须同时支持，可以：
- 使用不同的服务名（自定义：`api.health.v1.Health`，标准：`grpc.health.v1.Health`）
- 但会增加维护成本和复杂度

## 总结

| 对比项 | 标准服务 | 自定义服务 | 推荐 |
|--------|---------|-----------|------|
| **功能丰富度** | ⭐⭐ | ⭐⭐⭐⭐⭐ | 自定义 |
| **兼容性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | 标准 |
| **灵活性** | ⭐⭐ | ⭐⭐⭐⭐⭐ | 自定义 |
| **维护成本** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | 标准 |
| **Kubernetes 支持** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 自定义 |

**当前项目选择**：自定义服务（`api.health.v1.Health`）

**原因**：
- 需要 Kubernetes 就绪/存活探针
- 需要详细的诊断信息
- 需要 HTTP 和 gRPC 双协议支持

