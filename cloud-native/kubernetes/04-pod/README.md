# Pod 详解

## 📖 什么是 Pod？

Pod 是 Kubernetes 中最小的可部署单元，包含一个或多个容器。Pod 内的容器共享网络命名空间和存储卷，可以通过 localhost 相互通信。

## 🎯 Pod 特点

### 1. 生命周期
- **Pending**: Pod 已被调度，但容器镜像还在下载或容器还在启动
- **Running**: Pod 已绑定到节点，所有容器都已创建
- **Succeeded**: Pod 中所有容器都已成功终止
- **Failed**: Pod 中至少有一个容器异常终止
- **Unknown**: 无法获取 Pod 状态

### 2. 网络模型
- Pod 内的容器共享同一个 IP 地址
- 容器间可以通过 localhost 通信
- 每个 Pod 在集群内有唯一的 IP

### 3. 存储模型
- Pod 内的容器可以共享存储卷
- 支持多种存储类型（emptyDir、hostPath、PVC等）

## 📝 Pod 配置详解

### 基础 Pod 配置
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
    tier: frontend
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
      protocol: TCP
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
```

### 多容器 Pod
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-app
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
  - name: log-collector
    image: busybox
    command: ['sh', '-c', 'while true; do echo "$(date) - Log entry"; sleep 10; done']
```

## 🔧 Pod 生命周期

### 1. 启动阶段
1. **调度**: Scheduler 将 Pod 分配到节点
2. **镜像拉取**: 下载容器镜像
3. **容器启动**: 启动容器进程
4. **就绪检查**: 执行就绪探针
5. **服务就绪**: Pod 可以接收流量

### 2. 运行阶段
- **健康检查**: 定期执行存活探针
- **资源监控**: 监控 CPU、内存使用
- **日志收集**: 收集容器日志

### 3. 终止阶段
1. **优雅终止**: 发送 SIGTERM 信号
2. **强制终止**: 发送 SIGKILL 信号
3. **清理资源**: 清理网络、存储等资源

## 🏥 健康检查

### 1. 存活探针 (Liveness Probe)
检测容器是否正常运行，失败时重启容器。

```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    livenessProbe:
      httpGet:
        path: /health
        port: 80
      initialDelaySeconds: 30
      periodSeconds: 10
      timeoutSeconds: 5
      failureThreshold: 3
```

### 2. 就绪探针 (Readiness Probe)
检测容器是否准备好接收流量。

```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    readinessProbe:
      httpGet:
        path: /ready
        port: 80
      initialDelaySeconds: 5
      periodSeconds: 5
```

### 3. 启动探针 (Startup Probe)
检测容器是否完成启动。

```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    startupProbe:
      httpGet:
        path: /startup
        port: 80
      failureThreshold: 30
      periodSeconds: 10
```

## 💾 存储配置

### 1. emptyDir
临时存储，Pod 删除时数据丢失。

```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    volumeMounts:
    - name: cache-volume
      mountPath: /cache
  volumes:
  - name: cache-volume
    emptyDir: {}
```

### 2. hostPath
挂载主机文件系统。

```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    volumeMounts:
    - name: host-volume
      mountPath: /host-data
  volumes:
  - name: host-volume
    hostPath:
      path: /data
      type: Directory
```

### 3. ConfigMap
挂载配置文件。

```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    volumeMounts:
    - name: config-volume
      mountPath: /etc/nginx/conf.d
  volumes:
  - name: config-volume
    configMap:
      name: nginx-config
```

## 🌐 网络配置

### 1. 端口配置
```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - name: http
      containerPort: 80
      protocol: TCP
    - name: https
      containerPort: 443
      protocol: TCP
```

### 2. 环境变量
```yaml
spec:
  containers:
  - name: nginx
    image: nginx:latest
    env:
    - name: NGINX_PORT
      value: "80"
    - name: NGINX_HOST
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
```

## 🎯 Pod 调度

### 1. 节点选择器 (Node Selector)
```yaml
spec:
  nodeSelector:
    disk: ssd
    environment: production
  containers:
  - name: nginx
    image: nginx:latest
```

### 2. 节点亲和性 (Node Affinity)
```yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/e2e-az-name
            operator: In
            values:
            - e2e-az1
            - e2e-az2
  containers:
  - name: nginx
    image: nginx:latest
```

### 3. Pod 亲和性 (Pod Affinity)
```yaml
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - web
        topologyKey: kubernetes.io/hostname
  containers:
  - name: nginx
    image: nginx:latest
```

## 🛠️ 常用操作

### 1. 创建 Pod
```bash
# 从 YAML 文件创建
kubectl apply -f pod.yaml

# 直接创建
kubectl run nginx --image=nginx:latest
```

### 2. 查看 Pod
```bash
# 查看所有 Pod
kubectl get pods

# 查看详细信息
kubectl describe pod <pod-name>

# 查看日志
kubectl logs <pod-name>
```

### 3. 进入 Pod
```bash
# 进入容器
kubectl exec -it <pod-name> -- /bin/bash

# 在容器中执行命令
kubectl exec <pod-name> -- ls /app
```

### 4. 删除 Pod
```bash
# 删除 Pod
kubectl delete pod <pod-name>

# 强制删除
kubectl delete pod <pod-name> --grace-period=0 --force
```

## 🎯 最佳实践

### 1. 资源管理
- 设置资源请求和限制
- 监控资源使用情况
- 避免资源竞争

### 2. 健康检查
- 配置合适的探针
- 设置合理的超时时间
- 避免过于频繁的检查

### 3. 存储管理
- 选择合适的存储类型
- 注意数据持久性
- 管理存储容量

### 4. 网络配置
- 合理配置端口
- 使用服务发现
- 配置网络策略

## 🛠️ 实践练习

### 练习 1：基础 Pod
1. 创建简单的 nginx Pod
2. 配置端口映射
3. 测试访问

### 练习 2：多容器 Pod
1. 创建包含多个容器的 Pod
2. 配置容器间通信
3. 观察日志输出

### 练习 3：健康检查
1. 配置存活探针
2. 配置就绪探针
3. 测试故障恢复

## 📚 扩展阅读

- [Kubernetes Pod 官方文档](https://kubernetes.io/docs/concepts/workloads/pods/)
- [Pod 生命周期](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
- [Pod 安全策略](https://kubernetes.io/docs/concepts/security/pod-security-policy/)

## 🎯 下一步

掌握 Pod 后，继续学习：
- [ReplicaSet与Deployment](./05-deployment/README.md)
- [Service与网络](./06-service/README.md) 