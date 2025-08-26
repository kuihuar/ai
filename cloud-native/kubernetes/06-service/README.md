# Service 与网络

## 📖 什么是 Service？

Service 是 Kubernetes 中用于为一组 Pod 提供统一访问入口的资源对象。它通过标签选择器将请求路由到后端 Pod，实现了服务发现和负载均衡。

## 🎯 Service 特点

### 1. 服务发现
- 为 Pod 提供稳定的网络端点
- 自动发现后端 Pod
- 支持动态扩缩容

### 2. 负载均衡
- 自动分发请求到后端 Pod
- 支持多种负载均衡算法
- 健康检查确保流量分发到健康 Pod

### 3. 网络抽象
- 隐藏 Pod 的 IP 变化
- 提供稳定的服务名称
- 支持集群内外访问

## 🌐 Service 类型

### 1. ClusterIP（默认）
仅在集群内部可访问，为 Pod 提供内部服务发现。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  type: ClusterIP
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
```

### 2. NodePort
通过节点端口暴露服务，可以从集群外部访问。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  type: NodePort
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080
    protocol: TCP
```

### 3. LoadBalancer
使用云服务商的负载均衡器暴露服务。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  type: LoadBalancer
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
```

### 4. ExternalName
将服务映射到外部域名。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: external-service
spec:
  type: ExternalName
  externalName: api.example.com
```

## 📝 Service 配置详解

### 基础 Service 配置
```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-service
  labels:
    app: web
spec:
  selector:
    app: web
    tier: frontend
  ports:
  - name: http
    port: 80
    targetPort: 8080
    protocol: TCP
  - name: https
    port: 443
    targetPort: 8443
    protocol: TCP
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
```

### 多端口 Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: multi-port-service
spec:
  selector:
    app: multi-port-app
  ports:
  - name: http
    port: 80
    targetPort: 8080
  - name: https
    port: 443
    targetPort: 8443
  - name: metrics
    port: 9090
    targetPort: 9090
```

### 无选择器 Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: external-service
spec:
  ports:
  - port: 80
    targetPort: 8080
---
apiVersion: v1
kind: Endpoints
metadata:
  name: external-service
subsets:
- addresses:
  - ip: 192.168.1.10
  - ip: 192.168.1.11
  ports:
  - port: 8080
```

## 🔄 负载均衡策略

### 1. 轮询（Round Robin）
默认策略，依次将请求分发到后端 Pod。

### 2. 会话亲和性（Session Affinity）
基于客户端 IP 的会话亲和性。

```yaml
spec:
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
```

### 3. 自定义负载均衡
通过 EndpointSlice 实现自定义负载均衡。

## 🛠️ 常用操作

### 1. 创建 Service
```bash
# 从 YAML 文件创建
kubectl apply -f service.yaml

# 为 Deployment 创建 Service
kubectl expose deployment nginx --port=80 --target-port=80
```

### 2. 查看 Service
```bash
# 查看所有 Service
kubectl get services

# 查看详细信息
kubectl describe service <service-name>

# 查看 Endpoints
kubectl get endpoints <service-name>
```

### 3. 测试 Service
```bash
# 在集群内测试
kubectl run test-pod --image=busybox --rm -it --restart=Never -- nslookup nginx-service

# 端口转发
kubectl port-forward service/nginx-service 8080:80
```

### 4. 删除 Service
```bash
# 删除 Service
kubectl delete service <service-name>
```

## 🌐 网络策略

### 1. NetworkPolicy
控制 Pod 间的网络通信。

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

### 2. 允许特定流量
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-web-traffic
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 80
```

## 🔍 服务发现

### 1. DNS 服务发现
Kubernetes 自动为 Service 创建 DNS 记录。

```bash
# 在 Pod 中解析服务
nslookup nginx-service
nslookup nginx-service.default.svc.cluster.local
```

### 2. 环境变量
Pod 启动时自动注入 Service 环境变量。

```bash
# 查看环境变量
env | grep SERVICE
```

## 📊 监控和调试

### 1. 查看 Service 状态
```bash
# 查看 Service 详情
kubectl describe service nginx-service

# 查看 Endpoints
kubectl get endpoints nginx-service

# 查看 Service 事件
kubectl get events --field-selector involvedObject.name=nginx-service
```

### 2. 网络连通性测试
```bash
# 测试 Service 连通性
kubectl run test-pod --image=busybox --rm -it --restart=Never -- wget -O- nginx-service

# 测试端口连通性
kubectl run test-pod --image=busybox --rm -it --restart=Never -- nc -zv nginx-service 80
```

## 🎯 最佳实践

### 1. 命名规范
- 使用有意义的服务名称
- 遵循命名空间约定
- 使用标签进行分组

### 2. 端口管理
- 使用标准端口号
- 避免端口冲突
- 文档化端口用途

### 3. 安全配置
- 使用 NetworkPolicy 限制访问
- 配置适当的会话亲和性
- 监控异常流量

### 4. 性能优化
- 合理配置负载均衡策略
- 监控服务性能
- 优化网络配置

## 🛠️ 实践练习

### 练习 1：基础 Service
1. 创建 Deployment
2. 创建 ClusterIP Service
3. 测试服务发现

### 练习 2：外部访问
1. 创建 NodePort Service
2. 配置 LoadBalancer
3. 测试外部访问

### 练习 3：网络策略
1. 创建 NetworkPolicy
2. 测试网络隔离
3. 配置允许规则

## 📚 扩展阅读

- [Kubernetes Service 官方文档](https://kubernetes.io/docs/concepts/services-networking/service/)
- [网络策略](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [服务发现](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)

## 🎯 下一步

掌握 Service 后，继续学习：
- [ConfigMap与Secret](./07-config/README.md)
- [存储管理](./08-storage/README.md) 