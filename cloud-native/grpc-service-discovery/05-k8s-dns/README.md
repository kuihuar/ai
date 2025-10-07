# Kubernetes DNS 服务发现

## 📖 概述

Kubernetes 内置了强大的服务发现机制，通过 DNS 和 Service 资源实现服务注册与发现。这是云原生环境中最常用的服务发现方案。

## 🎯 核心特性

### 1. 内置服务发现
- 基于 DNS 的服务发现
- 自动服务注册
- 负载均衡

### 2. 服务类型
- ClusterIP: 集群内部访问
- NodePort: 节点端口访问
- LoadBalancer: 云负载均衡器
- ExternalName: 外部服务映射

### 3. 健康检查
- 内置健康检查
- 自动故障转移
- 服务端点管理

## 🚀 快速开始

### 1. 创建 Service 资源

```yaml
# service.yaml
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

### 2. 创建 Deployment

```yaml
# deployment.yaml
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
```

### 3. Go 客户端实现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/resolver"
    pb "your-project/proto"
)

// K8sResolver Kubernetes DNS 解析器
type K8sResolver struct {
    serviceName string
    namespace   string
    port        string
}

// NewK8sResolver 创建 K8s 解析器
func NewK8sResolver(serviceName, namespace, port string) *K8sResolver {
    return &K8sResolver{
        serviceName: serviceName,
        namespace:   namespace,
        port:        port,
    }
}

// Build 构建解析器
func (r *K8sResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    // 构建服务地址
    serviceAddr := fmt.Sprintf("%s.%s.svc.cluster.local:%s", r.serviceName, r.namespace, r.port)
    
    // 解析服务地址
    addresses, err := r.resolveService(serviceAddr)
    if err != nil {
        return nil, err
    }
    
    // 更新客户端连接
    cc.UpdateState(resolver.State{
        Addresses: addresses,
    })
    
    return r, nil
}

// Scheme 返回解析器方案
func (r *K8sResolver) Scheme() string {
    return "k8s"
}

// Close 关闭解析器
func (r *K8sResolver) Close() {
    // 清理资源
}

// ResolveNow 立即解析
func (r *K8sResolver) ResolveNow(opts resolver.ResolveNowOptions) {
    // 重新解析服务
}

// resolveService 解析服务地址
func (r *K8sResolver) resolveService(serviceAddr string) ([]resolver.Address, error) {
    // 解析 DNS
    addrs, err := net.LookupHost(serviceAddr)
    if err != nil {
        return nil, err
    }
    
    var addresses []resolver.Address
    for _, addr := range addrs {
        addresses = append(addresses, resolver.Address{
            Addr: addr,
        })
    }
    
    return addresses, nil
}

// K8sServiceDiscovery Kubernetes 服务发现
type K8sServiceDiscovery struct {
    serviceName string
    namespace   string
    port        string
}

// NewK8sServiceDiscovery 创建服务发现
func NewK8sServiceDiscovery(serviceName, namespace, port string) *K8sServiceDiscovery {
    return &K8sServiceDiscovery{
        serviceName: serviceName,
        namespace:   namespace,
        port:        port,
    }
}

// Discover 发现服务
func (ksd *K8sServiceDiscovery) Discover() ([]string, error) {
    serviceAddr := fmt.Sprintf("%s.%s.svc.cluster.local:%s", ksd.serviceName, ksd.namespace, ksd.port)
    
    addrs, err := net.LookupHost(serviceAddr)
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, addr := range addrs {
        addresses = append(addresses, fmt.Sprintf("%s:%s", addr, ksd.port))
    }
    
    return addresses, nil
}

// GetServiceAddress 获取服务地址
func (ksd *K8sServiceDiscovery) GetServiceAddress() string {
    return fmt.Sprintf("%s.%s.svc.cluster.local:%s", ksd.serviceName, ksd.namespace, ksd.port)
}
```

## 📝 使用示例

### 1. 服务端部署

```yaml
# k8s-deployment.yaml
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
        - name: SERVICE_NAME
          value: "your-service"
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
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
```

### 2. 客户端使用

```go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/resolver"
    pb "your-project/proto"
)

func main() {
    // 注册解析器
    resolver.Register(&K8sResolver{
        serviceName: "your-service",
        namespace:   "default",
        port:        "80",
    })
    
    // 连接服务
    conn, err := grpc.Dial("k8s://your-service.default.svc.cluster.local:80", grpc.WithInsecure())
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

### 1. Headless Service

```yaml
# headless-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: your-service-headless
  namespace: default
spec:
  selector:
    app: your-service
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  clusterIP: None  # Headless Service
```

### 2. ExternalName Service

```yaml
# external-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: external-service
  namespace: default
spec:
  type: ExternalName
  externalName: external-service.example.com
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
```

### 3. LoadBalancer Service

```yaml
# loadbalancer-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: your-service-lb
  namespace: default
spec:
  selector:
    app: your-service
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

## 📊 性能优化

### 1. 连接池配置

```go
// 配置 gRPC 连接池
func createGRPCConnection(serviceAddr string) (*grpc.ClientConn, error) {
    return grpc.Dial(serviceAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
        Time:                10 * time.Second,
        Timeout:             3 * time.Second,
        PermitWithoutStream: true,
    }))
}
```

### 2. 负载均衡

```go
// 使用 gRPC 负载均衡
func createGRPCConnectionWithLB(serviceAddr string) (*grpc.ClientConn, error) {
    return grpc.Dial(serviceAddr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
}
```

### 3. 服务网格集成

```go
// 使用 Istio 服务网格
func createGRPCConnectionWithIstio(serviceAddr string) (*grpc.ClientConn, error) {
    return grpc.Dial(serviceAddr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"), grpc.WithUnaryInterceptor(istioInterceptor))
}

func istioInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    // 添加 Istio 相关的头部信息
    ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", generateRequestID())
    return invoker(ctx, method, req, reply, cc, opts...)
}
```

## 🛡️ 安全配置

### 1. mTLS 配置

```yaml
# mTLS 配置
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
  name: your-service
  namespace: default
spec:
  host: your-service
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
```

### 2. 网络策略

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: your-service-netpol
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: your-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: your-client
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: your-service
    ports:
    - protocol: TCP
      port: 8080
```

## 🔍 监控和调试

### 1. 服务监控

```yaml
# service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: your-service-monitor
  namespace: default
spec:
  selector:
    matchLabels:
      app: your-service
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
```

### 2. 日志配置

```yaml
# fluentd-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
  namespace: kube-system
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*your-service*.log
      pos_file /var/log/fluentd-your-service.log.pos
      tag kubernetes.*
      format json
    </source>
    <match kubernetes.**>
      @type elasticsearch
      host elasticsearch.logging.svc.cluster.local
      port 9200
      index_name your-service
    </match>
```

## 📚 最佳实践

1. **服务命名**: 使用有意义的服务名称
2. **命名空间**: 合理使用命名空间隔离服务
3. **健康检查**: 实现快速响应的健康检查
4. **资源限制**: 设置合理的资源限制
5. **监控告警**: 监控服务状态和性能指标

## 🔗 相关资源

- [Kubernetes 服务文档](https://kubernetes.io/docs/concepts/services-networking/service/)
- [Kubernetes DNS 文档](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)
- [gRPC 负载均衡](https://grpc.io/docs/guides/load-balancing/)
