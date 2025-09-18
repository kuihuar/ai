# DaemonSet 详解

## 🚀 什么是 DaemonSet？

DaemonSet 是 Kubernetes 中用于确保集群中每个节点（或指定节点）都运行一个 Pod 副本的工作负载控制器。当有新节点加入集群时，DaemonSet 会自动在新节点上创建 Pod；当节点从集群中移除时，相应的 Pod 也会被删除。

## 🎯 DaemonSet 特点

- **节点级部署**：确保每个节点运行一个 Pod 副本
- **自动扩缩容**：节点加入/移除时自动创建/删除 Pod
- **系统级服务**：常用于运行系统级守护进程
- **资源监控**：每个节点运行监控、日志收集等服务
- **网络代理**：如 kube-proxy 在每个节点运行

## 📝 DaemonSet 配置

### 基础 DaemonSet 配置

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nginx-daemonset
  labels:
    app: nginx
spec:
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

### 带节点选择器的 DaemonSet

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: monitoring-daemonset
spec:
  selector:
    matchLabels:
      app: monitoring
  template:
    metadata:
      labels:
        app: monitoring
    spec:
      nodeSelector:
        kubernetes.io/os: linux
      containers:
      - name: monitoring
        image: prom/node-exporter:latest
        ports:
        - containerPort: 9100
```

### 带容忍度的 DaemonSet（在 Master 节点运行）

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: master-monitoring
spec:
  selector:
    matchLabels:
      app: master-monitoring
  template:
    metadata:
      labels:
        app: master-monitoring
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: NoSchedule
      - key: node-role.kubernetes.io/control-plane
        operator: Exists
        effect: NoSchedule
      containers:
      - name: monitoring
        image: prom/node-exporter:latest
        ports:
        - containerPort: 9100
```

## 🛠️ DaemonSet 操作

### 1. 创建 DaemonSet

```bash
# 使用 YAML 文件创建
kubectl apply -f daemonset.yaml

# 使用命令行创建
kubectl create daemonset nginx --image=nginx:latest
```

### 2. 查看 DaemonSet

```bash
# 查看所有 DaemonSet
kubectl get daemonsets

# 查看特定 DaemonSet 详情
kubectl describe daemonset <daemonset-name>

# 查看 DaemonSet 管理的 Pod
kubectl get pods -l app=nginx
```

### 3. 更新 DaemonSet

```bash
# 更新镜像
kubectl set image daemonset/nginx nginx=nginx:1.22

# 查看更新状态
kubectl rollout status daemonset/nginx

# 查看更新历史
kubectl rollout history daemonset/nginx
```

### 4. 回滚 DaemonSet

```bash
# 回滚到上一个版本
kubectl rollout undo daemonset/nginx

# 回滚到指定版本
kubectl rollout undo daemonset/nginx --to-revision=2
```

### 5. 删除 DaemonSet

```bash
# 删除 DaemonSet（会删除所有相关 Pod）
kubectl delete daemonset nginx

# 删除 DaemonSet 但保留 Pod
kubectl delete daemonset nginx --cascade=orphan
```

## 🔧 实际应用场景

### 1. 日志收集 - Fluentd

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: kube-system
spec:
  selector:
    matchLabels:
      name: fluentd
  template:
    metadata:
      labels:
        name: fluentd
    spec:
      serviceAccountName: fluentd
      containers:
      - name: fluentd
        image: fluent/fluentd-kubernetes-daemonset:v1-debian-elasticsearch
        env:
        - name: FLUENT_ELASTICSEARCH_HOST
          value: "elasticsearch.logging.svc.cluster.local"
        - name: FLUENT_ELASTICSEARCH_PORT
          value: "9200"
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
```

### 2. 监控代理 - Node Exporter

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-exporter
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: node-exporter
  template:
    metadata:
      labels:
        app: node-exporter
    spec:
      hostNetwork: true
      hostPID: true
      containers:
      - name: node-exporter
        image: prom/node-exporter:latest
        args:
        - --path.procfs=/host/proc
        - --path.sysfs=/host/sys
        - --collector.filesystem.ignored-mount-points
        - ^/(sys|proc|dev|host|etc)($|/)
        ports:
        - containerPort: 9100
        volumeMounts:
        - name: proc
          mountPath: /host/proc
          readOnly: true
        - name: sys
          mountPath: /host/sys
          readOnly: true
      volumes:
      - name: proc
        hostPath:
          path: /proc
      - name: sys
        hostPath:
          path: /sys
```

## 🎯 练习

### 练习 1：基础 DaemonSet
1. 创建一个 nginx DaemonSet
2. 查看 Pod 分布情况
3. 更新镜像版本
4. 验证更新结果

### 练习 2：日志收集 DaemonSet
1. 创建 Fluentd DaemonSet 用于日志收集
2. 配置挂载宿主机日志目录
3. 验证日志收集功能

### 练习 3：监控 DaemonSet
1. 创建 Node Exporter DaemonSet
2. 配置监控数据收集
3. 验证监控数据可用性

## 🔍 故障排查

### 常见问题

1. **Pod 无法调度到某些节点**
   ```bash
   # 检查节点标签和选择器
   kubectl get nodes --show-labels
   kubectl describe daemonset <daemonset-name>
   ```

2. **Pod 启动失败**
   ```bash
   # 查看 Pod 事件和日志
   kubectl describe pod <pod-name>
   kubectl logs <pod-name>
   ```

3. **权限问题**
   ```bash
   # 检查 ServiceAccount 和 RBAC
   kubectl get serviceaccount
   kubectl describe clusterrolebinding
   ```

## 📚 相关资源

- [Kubernetes DaemonSet 官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
- [DaemonSet 最佳实践](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/#writing-a-daemonset-spec)

## 🎯 下一步学习

掌握 DaemonSet 后，继续学习：
- [StatefulSet](./17-statefulset/README.md) - 有状态应用管理
- [Job 和 CronJob](./18-job-cronjob/README.md) - 批处理任务
- [Service](./06-service/README.md) - 服务发现和负载均衡
