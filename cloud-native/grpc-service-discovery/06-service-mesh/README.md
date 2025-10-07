# 服务网格服务发现

## 📖 概述

服务网格（Service Mesh）是一种处理服务间通信的基础设施层，提供了服务发现、负载均衡、故障恢复、指标收集等功能。Istio 是最流行的服务网格实现。

## 🎯 核心特性

### 1. 服务发现
- 自动服务注册和发现
- 多集群服务发现
- 服务拓扑感知

### 2. 流量管理
- 智能路由
- 负载均衡
- 故障注入

### 3. 安全
- mTLS 加密
- 身份认证
- 授权策略

### 4. 可观测性
- 分布式追踪
- 指标收集
- 日志聚合

## 🚀 快速开始

### 1. 安装 Istio

```bash
# 下载 Istio
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.19.0
export PATH=$PWD/bin:$PATH

# 安装 Istio
istioctl install --set values.defaultRevision=default
kubectl label namespace default istio-injection=enabled
```

### 2. 部署服务

```yaml
# service-mesh-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: your-service
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: your-service
  template:
    metadata:
      labels:
        app: your-service
    spec:
      containers:
      - name: your-service
        image: your-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
---
apiVersion: v1
kind: Service
metadata:
  name: your-service
  namespace: default
spec:
  selector:
    app: your-service
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

### 3. Go 客户端实现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "your-project/proto"
)

// ServiceMeshDiscovery 服务网格服务发现
type ServiceMeshDiscovery struct {
    serviceName string
    namespace   string
    port        string
}

// NewServiceMeshDiscovery 创建服务网格服务发现
func NewServiceMeshDiscovery(serviceName, namespace, port string) *ServiceMeshDiscovery {
    return &ServiceMeshDiscovery{
        serviceName: serviceName,
        namespace:   namespace,
        port:        port,
    }
}

// GetServiceAddress 获取服务地址
func (smd *ServiceMeshDiscovery) GetServiceAddress() string {
    return fmt.Sprintf("%s.%s.svc.cluster.local:%s", smd.serviceName, smd.namespace, smd.port)
}

// CreateConnection 创建 gRPC 连接
func (smd *ServiceMeshDiscovery) CreateConnection() (*grpc.ClientConn, error) {
    serviceAddr := smd.GetServiceAddress()
    
    // 使用服务网格的负载均衡和故障恢复
    return grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// CreateConnectionWithOptions 创建带选项的连接
func (smd *ServiceMeshDiscovery) CreateConnectionWithOptions(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
    serviceAddr := smd.GetServiceAddress()
    
    // 添加默认选项
    defaultOpts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                10 * time.Second,
            Timeout:             3 * time.Second,
            PermitWithoutStream: true,
        }),
    }
    
    // 合并选项
    allOpts := append(defaultOpts, opts...)
    
    return grpc.Dial(serviceAddr, allOpts...)
}
```

## 📝 使用示例

### 1. 服务端部署

```yaml
# istio-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: your-service
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: your-service
  template:
    metadata:
      labels:
        app: your-service
    spec:
      containers:
      - name: your-service
        image: your-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: your-service
  namespace: default
spec:
  selector:
    app: your-service
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: your-service
  namespace: default
spec:
  hosts:
  - your-service
  http:
  - route:
    - destination:
        host: your-service
        port:
          number: 80
---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: your-service
  namespace: default
spec:
  host: your-service
  trafficPolicy:
    loadBalancer:
      simple: ROUND_ROBIN
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 10
        maxRequestsPerConnection: 2
    circuitBreaker:
      consecutiveErrors: 3
      interval: 30s
      baseEjectionTime: 30s
```

### 2. 客户端使用

```go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "your-project/proto"
)

func main() {
    // 创建服务网格服务发现
    discovery := NewServiceMeshDiscovery("your-service", "default", "80")
    
    // 创建连接
    conn, err := discovery.CreateConnection()
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()
    
    client := pb.NewYourServiceClient(conn)
    
    // 调用服务
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    resp, err := client.YourMethod(ctx, &pb.YourRequest{})
    if err != nil {
        log.Fatalf("failed to call service: %v", err)
    }
    
    log.Printf("Response: %v", resp)
}
```

## 🔧 高级配置

### 1. 流量管理

```yaml
# traffic-management.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: your-service-traffic
  namespace: default
spec:
  hosts:
  - your-service
  http:
  - match:
    - headers:
        version:
          exact: v1
    route:
    - destination:
        host: your-service
        subset: v1
        port:
          number: 80
  - match:
    - headers:
        version:
          exact: v2
    route:
    - destination:
        host: your-service
        subset: v2
        port:
          number: 80
  - route:
    - destination:
        host: your-service
        subset: v1
        port:
          number: 80
      weight: 90
    - destination:
        host: your-service
        subset: v2
        port:
          number: 80
      weight: 10
---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: your-service-subsets
  namespace: default
spec:
  host: your-service
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
```

### 2. 故障注入

```yaml
# fault-injection.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: your-service-fault
  namespace: default
spec:
  hosts:
  - your-service
  http:
  - fault:
      delay:
        percentage:
          value: 0.1
        fixedDelay: 5s
      abort:
        percentage:
          value: 0.1
        httpStatus: 500
    route:
    - destination:
        host: your-service
        port:
          number: 80
```

### 3. 安全配置

```yaml
# security-config.yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: default
spec:
  mtls:
    mode: STRICT
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: your-service-auth
  namespace: default
spec:
  selector:
    matchLabels:
      app: your-service
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/default/sa/your-client"]
    to:
    - operation:
        methods: ["GET", "POST"]
```

## 📊 性能优化

### 1. 连接池配置

```yaml
# connection-pool.yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: your-service-pool
  namespace: default
spec:
  host: your-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
        connectTimeout: 30ms
        tcpKeepalive:
          time: 7200s
          interval: 75s
          probes: 9
      http:
        http1MaxPendingRequests: 10
        http2MaxRequests: 100
        maxRequestsPerConnection: 2
        maxRetries: 3
        consecutiveGatewayErrors: 5
        interval: 30s
        baseEjectionTime: 30s
        maxEjectionPercent: 50
```

### 2. 负载均衡配置

```yaml
# load-balancing.yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: your-service-lb
  namespace: default
spec:
  host: your-service
  trafficPolicy:
    loadBalancer:
      simple: LEAST_CONN
      consistentHash:
        httpHeaderName: x-user-id
```

### 3. 熔断器配置

```yaml
# circuit-breaker.yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: your-service-cb
  namespace: default
spec:
  host: your-service
  trafficPolicy:
    circuitBreaker:
      consecutiveErrors: 3
      interval: 30s
      baseEjectionTime: 30s
      maxEjectionPercent: 50
      minHealthPercent: 30
```

## 🛡️ 安全配置

### 1. mTLS 配置

```yaml
# mtls-config.yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: default
spec:
  mtls:
    mode: STRICT
---
apiVersion: security.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: your-service-mtls
  namespace: default
spec:
  host: your-service
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
```

### 2. 授权策略

```yaml
# authorization-policy.yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: your-service-auth
  namespace: default
spec:
  selector:
    matchLabels:
      app: your-service
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/default/sa/your-client"]
    to:
    - operation:
        methods: ["GET", "POST"]
        paths: ["/api/v1/*"]
```

## 🔍 监控和调试

### 1. 指标监控

```yaml
# metrics-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: istio-config
  namespace: istio-system
data:
  mesh: |
    defaultConfig:
      proxyStatsMatcher:
        inclusionRegexps:
        - ".*circuit_breakers.*"
        - ".*upstream_rq_retries.*"
        - ".*upstream_rq_timeout.*"
        - ".*upstream_rq_pending.*"
        - ".*upstream_rq_active.*"
        - ".*upstream_rq_total.*"
        - ".*upstream_rq_retry.*"
        - ".*upstream_rq_retry_overflow.*"
        - ".*upstream_rq_retry_success.*"
        - ".*upstream_rq_retry_cancelled.*"
        - ".*upstream_rq_retry_direct_reset.*"
        - ".*upstream_rq_retry_remote_reset.*"
        - ".*upstream_rq_retry_local_reset.*"
        - ".*upstream_rq_retry_connect_failure.*"
        - ".*upstream_rq_retry_connect_timeout.*"
        - ".*upstream_rq_retry_remote_origin.*"
        - ".*upstream_rq_retry_remote_origin_timeout.*"
        - ".*upstream_rq_retry_remote_origin_connect_failure.*"
        - ".*upstream_rq_retry_remote_origin_connect_timeout.*"
        - ".*upstream_rq_retry_remote_origin_direct_reset.*"
        - ".*upstream_rq_retry_remote_origin_remote_reset.*"
        - ".*upstream_rq_retry_remote_origin_local_reset.*"
        - ".*upstream_rq_retry_remote_origin_overflow.*"
        - ".*upstream_rq_retry_remote_origin_success.*"
        - ".*upstream_rq_retry_remote_origin_cancelled.*"
        - ".*upstream_rq_retry_remote_origin_direct_reset.*"
        - ".*upstream_rq_retry_remote_origin_remote_reset.*"
        - ".*upstream_rq_retry_remote_origin_local_reset.*"
        - ".*upstream_rq_retry_remote_origin_overflow.*"
        - ".*upstream_rq_retry_remote_origin_success.*"
        - ".*upstream_rq_retry_remote_origin_cancelled.*"
```

### 2. 分布式追踪

```yaml
# tracing-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tracing-config
  namespace: istio-system
data:
  tracing: |
    sampling: 100
    zipkin:
      address: zipkin.istio-system.svc.cluster.local:9411
```

## 📚 最佳实践

1. **渐进式部署**: 逐步启用服务网格功能
2. **监控告警**: 监控服务网格指标
3. **安全策略**: 实施严格的安全策略
4. **性能优化**: 优化连接池和负载均衡
5. **故障恢复**: 配置合适的熔断器

## 🔗 相关资源

- [Istio 官方文档](https://istio.io/latest/docs/)
- [Istio 服务发现](https://istio.io/latest/docs/ops/configuration/traffic-management/service-discovery/)
- [Istio 流量管理](https://istio.io/latest/docs/concepts/traffic-management/)
