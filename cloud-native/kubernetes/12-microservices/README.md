# 微服务部署

## 📖 微服务概述

微服务架构是一种将应用程序拆分为小型、独立服务的架构模式。在 Kubernetes 上部署微服务需要考虑服务发现、负载均衡、配置管理、监控等多个方面。

## 🎯 微服务特点

### 1. 服务拆分
- 按业务功能拆分
- 独立开发部署
- 技术栈多样化

### 2. 服务通信
- 同步通信（HTTP/gRPC）
- 异步通信（消息队列）
- 服务网格（Service Mesh）

### 3. 数据管理
- 数据库 per 服务
- 分布式事务
- 数据一致性

## 🏗️ 微服务架构

### 1. 典型架构
```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway                              │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   User      │ │   Order     │ │   Payment   │            │
│  │  Service    │ │  Service    │ │  Service    │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   User      │ │   Order     │ │   Payment   │            │
│  │  Database   │ │  Database   │ │  Database   │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### 2. 服务组件
- **API Gateway**: 统一入口
- **Service Registry**: 服务注册
- **Load Balancer**: 负载均衡
- **Circuit Breaker**: 熔断器
- **Distributed Tracing**: 分布式追踪

## 📝 服务定义

### 1. User Service
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: user-service-config
              key: db_host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: user-service-secret
              key: db_password
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
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

### 2. Order Service
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  labels:
    app: order-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
      - name: order-service
        image: order-service:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: USER_SERVICE_URL
          value: "http://user-service"
        - name: PAYMENT_SERVICE_URL
          value: "http://payment-service"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: order-service-config
              key: db_host
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: order-service
spec:
  selector:
    app: order-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

### 3. Payment Service
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-service
  labels:
    app: payment-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: payment-service
  template:
    metadata:
      labels:
        app: payment-service
    spec:
      containers:
      - name: payment-service
        image: payment-service:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: payment-service-config
              key: db_host
        - name: API_KEY
          valueFrom:
            secretKeyRef:
              name: payment-service-secret
              key: api_key
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: payment-service
spec:
  selector:
    app: payment-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

## 🌐 API Gateway

### 1. Kong Gateway
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kong-gateway
spec:
  replicas: 2
  selector:
    matchLabels:
      app: kong-gateway
  template:
    metadata:
      labels:
        app: kong-gateway
    spec:
      containers:
      - name: kong
        image: kong:2.8
        env:
        - name: KONG_DATABASE
          value: "off"
        - name: KONG_PROXY_ACCESS_LOG
          value: "/dev/stdout"
        - name: KONG_ADMIN_ACCESS_LOG
          value: "/dev/stdout"
        - name: KONG_PROXY_ERROR_LOG
          value: "/dev/stderr"
        - name: KONG_ADMIN_ERROR_LOG
          value: "/dev/stderr"
        - name: KONG_ADMIN_LISTEN
          value: "0.0.0.0:8001"
        - name: KONG_ADMIN_GUI_URL
          value: "http://localhost:8002"
        ports:
        - containerPort: 8000
          name: proxy
        - containerPort: 8001
          name: admin
        - containerPort: 8443
          name: proxy-ssl
        - containerPort: 8444
          name: admin-ssl
---
apiVersion: v1
kind: Service
metadata:
  name: kong-gateway
spec:
  selector:
    app: kong-gateway
  ports:
  - port: 80
    targetPort: 8000
    name: proxy
  - port: 8001
    targetPort: 8001
    name: admin
  type: LoadBalancer
```

### 2. Ingress Controller
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: microservices-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /users
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 80
      - path: /orders
        pathType: Prefix
        backend:
          service:
            name: order-service
            port:
              number: 80
      - path: /payments
        pathType: Prefix
        backend:
          service:
            name: payment-service
            port:
              number: 80
```

## 🔄 服务通信

### 1. 同步通信
```yaml
# 服务间 HTTP 调用
apiVersion: v1
kind: ConfigMap
metadata:
  name: service-config
data:
  user_service_url: "http://user-service"
  order_service_url: "http://order-service"
  payment_service_url: "http://payment-service"
```

### 2. 异步通信
```yaml
# Redis 消息队列
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:6.2
        ports:
        - containerPort: 6379
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: redis
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: 6379
  type: ClusterIP
```

## 🔍 服务发现

### 1. DNS 服务发现
```yaml
# 使用 Kubernetes DNS
apiVersion: v1
kind: ConfigMap
metadata:
  name: service-discovery
data:
  user_service: "user-service.default.svc.cluster.local"
  order_service: "order-service.default.svc.cluster.local"
  payment_service: "payment-service.default.svc.cluster.local"
```

### 2. 服务注册
```yaml
# 使用 Consul 服务注册
apiVersion: apps/v1
kind: Deployment
metadata:
  name: consul
spec:
  replicas: 3
  selector:
    matchLabels:
      app: consul
  template:
    metadata:
      labels:
        app: consul
    spec:
      containers:
      - name: consul
        image: consul:1.12
        ports:
        - containerPort: 8500
        - containerPort: 8600
        command:
        - consul
        - agent
        - -server
        - -bootstrap-expect=3
        - -ui
        - -client=0.0.0.0
---
apiVersion: v1
kind: Service
metadata:
  name: consul
spec:
  selector:
    app: consul
  ports:
  - port: 8500
    targetPort: 8500
    name: http
  - port: 8600
    targetPort: 8600
    name: dns
  type: ClusterIP
```

## 🛡️ 熔断器

### 1. Hystrix 配置
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hystrix-config
data:
  hystrix.command.default.circuitBreaker.enabled: "true"
  hystrix.command.default.circuitBreaker.requestVolumeThreshold: "20"
  hystrix.command.default.circuitBreaker.errorThresholdPercentage: "50"
  hystrix.command.default.circuitBreaker.sleepWindowInMilliseconds: "5000"
  hystrix.command.default.execution.isolation.thread.timeoutInMilliseconds: "3000"
```

### 2. 重试机制
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: retry-config
data:
  max_retries: "3"
  retry_delay: "1000"
  backoff_multiplier: "2"
```

## 📊 监控和追踪

### 1. 分布式追踪
```yaml
# Jaeger 追踪
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jaeger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jaeger
  template:
    metadata:
      labels:
        app: jaeger
    spec:
      containers:
      - name: jaeger
        image: jaegertracing/all-in-one:latest
        ports:
        - containerPort: 16686
          name: ui
        - containerPort: 14268
          name: collector
        env:
        - name: COLLECTOR_OTLP_ENABLED
          value: "true"
---
apiVersion: v1
kind: Service
metadata:
  name: jaeger
spec:
  selector:
    app: jaeger
  ports:
  - port: 16686
    targetPort: 16686
    name: ui
  - port: 14268
    targetPort: 14268
    name: collector
  type: ClusterIP
```

### 2. 服务监控
```yaml
# Prometheus 监控
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: 'kubernetes-pods'
      kubernetes_sd_configs:
      - role: pod
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
```

## 🎯 最佳实践

### 1. 服务设计
- 单一职责原则
- 松耦合设计
- 接口标准化

### 2. 部署策略
- 蓝绿部署
- 金丝雀发布
- 滚动更新

### 3. 数据管理
- 数据库设计
- 事务处理
- 数据一致性

### 4. 安全考虑
- 服务间认证
- 数据加密
- 访问控制

## 🛠️ 实践练习

### 练习 1：基础微服务
1. 部署用户服务
2. 部署订单服务
3. 配置服务间通信

### 练习 2：API Gateway
1. 部署 Kong Gateway
2. 配置路由规则
3. 测试服务访问

### 练习 3：监控追踪
1. 部署 Jaeger
2. 配置分布式追踪
3. 分析调用链路

## 📚 扩展阅读

- [微服务架构模式](https://microservices.io/patterns/)
- [Kubernetes 微服务部署](https://kubernetes.io/docs/concepts/services-networking/)
- [服务网格](https://istio.io/docs/)

## 🎯 下一步

掌握微服务部署后，继续学习：
- [CI/CD流水线](./13-cicd/README.md)
- [故障排查](./14-troubleshooting/README.md) 