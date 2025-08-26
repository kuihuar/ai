# ReplicaSet 与 Deployment

## 📖 什么是 ReplicaSet？

ReplicaSet 是 Kubernetes 中用于确保指定数量的 Pod 副本始终运行的一种控制器。它通过标签选择器来管理 Pod，当 Pod 数量不足时会自动创建新的 Pod，当数量过多时会删除多余的 Pod。

## 🎯 ReplicaSet 特点

### 1. 自动扩缩容
- 根据配置的副本数自动调整 Pod 数量
- 支持手动扩缩容和自动扩缩容（HPA）

### 2. 故障恢复
- 当 Pod 故障时自动创建新的 Pod
- 确保应用的高可用性

### 3. 标签管理
- 通过标签选择器管理 Pod
- 支持复杂的标签匹配规则

## 📝 ReplicaSet 配置

### 基础 ReplicaSet 配置
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-replicaset
  labels:
    app: nginx
    tier: frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
```

### 高级选择器配置
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-replicaset
spec:
  replicas: 3
  selector:
    matchExpressions:
    - key: app
      operator: In
      values:
      - nginx
      - web
    - key: environment
      operator: NotIn
      values:
      - test
  template:
    metadata:
      labels:
        app: nginx
        environment: production
    spec:
      containers:
      - name: nginx
        image: nginx:latest
```

## 🚀 什么是 Deployment？

Deployment 是 Kubernetes 中用于管理应用程序部署的高级控制器，它基于 ReplicaSet 构建，提供了声明式更新、回滚、暂停和恢复等功能。

## 🎯 Deployment 特点

### 1. 声明式更新
- 支持滚动更新和重新创建更新
- 自动管理更新过程

### 2. 回滚功能
- 支持快速回滚到之前的版本
- 保留更新历史记录

### 3. 暂停和恢复
- 可以暂停更新过程
- 支持分阶段更新

### 4. 扩缩容
- 支持手动和自动扩缩容
- 集成 HPA（Horizontal Pod Autoscaler）

## 📝 Deployment 配置

### 基础 Deployment 配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.19
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

### 滚动更新配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.20
        ports:
        - containerPort: 80
```

### 重新创建更新配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
```

## 🔄 更新策略

### 1. RollingUpdate（滚动更新）
- 逐步替换 Pod，确保服务不中断
- 可以配置最大可用和最大超出数量
- 适合大多数应用场景

### 2. Recreate（重新创建）
- 先删除所有旧 Pod，再创建新 Pod
- 更新过程中服务会短暂中断
- 适合不支持多版本并存的应用

## 🛠️ 常用操作

### 1. 创建 Deployment
```bash
# 从 YAML 文件创建
kubectl apply -f deployment.yaml

# 直接创建
kubectl create deployment nginx --image=nginx:latest
```

### 2. 查看 Deployment
```bash
# 查看所有 Deployment
kubectl get deployments

# 查看详细信息
kubectl describe deployment <deployment-name>

# 查看 ReplicaSet
kubectl get replicasets
```

### 3. 扩缩容
```bash
# 手动扩缩容
kubectl scale deployment nginx --replicas=5

# 自动扩缩容
kubectl autoscale deployment nginx --min=2 --max=10 --cpu-percent=80
```

### 4. 更新镜像
```bash
# 更新镜像版本
kubectl set image deployment/nginx nginx=nginx:1.21

# 查看更新状态
kubectl rollout status deployment/nginx
```

### 5. 回滚操作
```bash
# 查看更新历史
kubectl rollout history deployment/nginx

# 回滚到上一个版本
kubectl rollout undo deployment/nginx

# 回滚到指定版本
kubectl rollout undo deployment/nginx --to-revision=2
```

### 6. 暂停和恢复
```bash
# 暂停更新
kubectl rollout pause deployment/nginx

# 恢复更新
kubectl rollout resume deployment/nginx
```

## 📊 状态监控

### 1. 查看更新状态
```bash
# 查看更新进度
kubectl rollout status deployment/nginx

# 查看 Pod 状态
kubectl get pods -l app=nginx
```

### 2. 查看事件
```bash
# 查看 Deployment 事件
kubectl describe deployment nginx

# 查看 Pod 事件
kubectl describe pods -l app=nginx
```

## 🎯 最佳实践

### 1. 标签管理
- 使用有意义的标签
- 保持标签的一致性
- 避免标签冲突

### 2. 资源管理
- 设置合理的资源请求和限制
- 监控资源使用情况
- 配置 HPA 实现自动扩缩容

### 3. 更新策略
- 选择合适的更新策略
- 配置合理的更新参数
- 测试更新过程

### 4. 健康检查
- 配置存活探针和就绪探针
- 设置合理的超时时间
- 监控应用健康状态

## 🛠️ 实践练习

### 练习 1：基础 Deployment
1. 创建 nginx Deployment
2. 配置 3 个副本
3. 测试扩缩容功能

### 练习 2：滚动更新
1. 创建 Deployment
2. 执行滚动更新
3. 观察更新过程
4. 测试回滚功能

### 练习 3：自动扩缩容
1. 配置 HPA
2. 模拟负载增加
3. 观察自动扩缩容

## 📚 扩展阅读

- [Kubernetes Deployment 官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [ReplicaSet 官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/)
- [滚动更新最佳实践](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/)

## 🎯 下一步

掌握 Deployment 后，继续学习：
- [Service与网络](./06-service/README.md)
- [ConfigMap与Secret](./07-config/README.md) 