# StatefulSet 详解

## 🚀 什么是 StatefulSet？

StatefulSet 是 Kubernetes 中用于管理有状态应用的工作负载控制器。与 Deployment 不同，StatefulSet 为每个 Pod 维护一个稳定的、唯一的标识符，即使 Pod 被重新调度到其他节点，其标识符和存储也会保持不变。

## 🎯 StatefulSet 特点

- **稳定网络标识**：每个 Pod 有固定的主机名和 DNS 记录
- **稳定存储**：每个 Pod 有独立的持久化存储
- **有序部署**：Pod 按顺序创建、更新和删除
- **有序扩缩容**：按顺序进行扩缩容操作
- **有序删除**：删除时按逆序进行

## 📝 StatefulSet 配置

### 基础 StatefulSet 配置

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
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
        image: nginx:1.21
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: www
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

### 带 Headless Service 的 StatefulSet

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
  - port: 80
    name: web
  clusterIP: None
  selector:
    app: nginx
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
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
        image: nginx:1.21
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: www
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

### 数据库 StatefulSet 示例

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
  - port: 3306
    name: mysql
  clusterIP: None
  selector:
    app: mysql
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: "mysql"
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "password"
        - name: MYSQL_DATABASE
          value: "testdb"
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysql-storage
          mountPath: /var/lib/mysql
  volumeClaimTemplates:
  - metadata:
      name: mysql-storage
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 10Gi
```

## 🛠️ StatefulSet 操作

### 1. 创建 StatefulSet

```bash
# 使用 YAML 文件创建
kubectl apply -f statefulset.yaml

# 查看 StatefulSet 状态
kubectl get statefulsets
kubectl describe statefulset web
```

### 2. 查看 Pod 状态

```bash
# 查看 StatefulSet 管理的 Pod
kubectl get pods -l app=nginx

# 查看 Pod 详细信息
kubectl describe pod web-0
kubectl describe pod web-1
kubectl describe pod web-2
```

### 3. 扩缩容 StatefulSet

```bash
# 扩容到 5 个副本
kubectl scale statefulset web --replicas=5

# 缩容到 2 个副本
kubectl scale statefulset web --replicas=2

# 查看扩缩容状态
kubectl get pods -l app=nginx -w
```

### 4. 更新 StatefulSet

```bash
# 更新镜像
kubectl set image statefulset/web nginx=nginx:1.22

# 查看更新状态
kubectl rollout status statefulset/web

# 查看更新历史
kubectl rollout history statefulset/web
```

### 5. 删除 StatefulSet

```bash
# 删除 StatefulSet（会删除所有 Pod）
kubectl delete statefulset web

# 删除 StatefulSet 但保留 Pod
kubectl delete statefulset web --cascade=orphan
```

## 🔧 实际应用场景

### 1. 数据库集群 - MySQL 主从复制

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql-config
data:
  master.cnf: |
    [mysqld]
    log-bin
    skip-name-resolve
  slave.cnf: |
    [mysqld]
    super-read-only
    skip-name-resolve
---
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
  - port: 3306
    name: mysql
  clusterIP: None
  selector:
    app: mysql
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: "mysql"
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      initContainers:
      - name: init-mysql
        image: mysql:8.0
        command:
        - bash
        - "-c"
        - |
          set -ex
          [[ $HOSTNAME =~ -([0-9]+)$ ]] || exit 1
          ordinal=${BASH_REMATCH[1]}
          [[ $ordinal -eq 0 ]] && echo "server-id=1" > /mnt/conf.d/server-id.cnf
          [[ $ordinal -ne 0 ]] && echo "server-id=$((100 + $ordinal))" > /mnt/conf.d/server-id.cnf
          [[ $ordinal -ne 0 ]] && echo "log-slave-updates=1" >> /mnt/conf.d/server-id.cnf
          cp /mnt/config-map/master.cnf /mnt/conf.d/
          [[ $ordinal -ne 0 ]] && cp /mnt/config-map/slave.cnf /mnt/conf.d/
        volumeMounts:
        - name: conf
          mountPath: /mnt/conf.d
        - name: config-map
          mountPath: /mnt/config-map
      containers:
      - name: mysql
        image: mysql:8.0
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "password"
        - name: MYSQL_REPLICATION_MODE
          value: "master"
        - name: MYSQL_REPLICATION_USER
          value: "replicator"
        - name: MYSQL_REPLICATION_PASSWORD
          value: "replpass"
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
        - name: conf
          mountPath: /etc/mysql/conf.d
      volumes:
      - name: conf
        emptyDir: {}
      - name: config-map
        configMap:
          name: mysql-config
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 10Gi
```

### 2. 消息队列 - RabbitMQ 集群

```yaml
apiVersion: v1
kind: Service
metadata:
  name: rabbitmq
  labels:
    app: rabbitmq
spec:
  ports:
  - port: 5672
    name: amqp
  - port: 15672
    name: management
  clusterIP: None
  selector:
    app: rabbitmq
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: rabbitmq
spec:
  serviceName: "rabbitmq"
  replicas: 3
  selector:
    matchLabels:
      app: rabbitmq
  template:
    metadata:
      labels:
        app: rabbitmq
    spec:
      containers:
      - name: rabbitmq
        image: rabbitmq:3.11-management
        env:
        - name: RABBITMQ_ERLANG_COOKIE
          value: "SWQOKODSQALRPCLNMEQG"
        - name: RABBITMQ_NODENAME
          value: "rabbit@$(HOSTNAME).rabbitmq.default.svc.cluster.local"
        ports:
        - containerPort: 5672
          name: amqp
        - containerPort: 15672
          name: management
        volumeMounts:
        - name: rabbitmq-storage
          mountPath: /var/lib/rabbitmq
        command:
        - bash
        - -c
        - |
          set -euo pipefail
          if [ -f /var/lib/rabbitmq/.erlang.cookie ]; then
            echo "Cookie already exists"
          else
            echo "SWQOKODSQALRPCLNMEQG" > /var/lib/rabbitmq/.erlang.cookie
            chmod 600 /var/lib/rabbitmq/.erlang.cookie
          fi
          exec docker-entrypoint.sh rabbitmq-server
  volumeClaimTemplates:
  - metadata:
      name: rabbitmq-storage
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 5Gi
```

### 3. 缓存集群 - Redis 集群

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-cluster
data:
  redis.conf: |
    port 6379
    cluster-enabled yes
    cluster-config-file nodes.conf
    cluster-node-timeout 5000
    appendonly yes
    protected-mode no
---
apiVersion: v1
kind: Service
metadata:
  name: redis
  labels:
    app: redis
spec:
  ports:
  - port: 6379
    name: redis
  - port: 16379
    name: cluster
  clusterIP: None
  selector:
    app: redis
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
spec:
  serviceName: "redis"
  replicas: 6
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
        image: redis:7-alpine
        command:
        - redis-server
        - /etc/redis/redis.conf
        ports:
        - containerPort: 6379
          name: redis
        - containerPort: 16379
          name: cluster
        volumeMounts:
        - name: redis-storage
          mountPath: /data
        - name: redis-config
          mountPath: /etc/redis
      volumes:
      - name: redis-config
        configMap:
          name: redis-cluster
  volumeClaimTemplates:
  - metadata:
      name: redis-storage
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

## 🎯 练习

### 练习 1：基础 StatefulSet
1. 创建一个 nginx StatefulSet
2. 查看 Pod 的命名规则
3. 验证持久化存储
4. 测试扩缩容

### 练习 2：数据库 StatefulSet
1. 创建 MySQL StatefulSet
2. 配置持久化存储
3. 验证数据持久性
4. 测试故障恢复

### 练习 3：集群应用 StatefulSet
1. 创建 Redis 集群 StatefulSet
2. 配置集群发现
3. 验证集群功能
4. 测试节点故障恢复

## 🔍 故障排查

### 常见问题

1. **Pod 启动失败**
   ```bash
   # 查看 Pod 事件和日志
   kubectl describe pod <pod-name>
   kubectl logs <pod-name>
   ```

2. **存储问题**
   ```bash
   # 检查 PVC 状态
   kubectl get pvc
   kubectl describe pvc <pvc-name>
   ```

3. **网络问题**
   ```bash
   # 检查 Service 和 DNS
   kubectl get svc
   nslookup <service-name>.<namespace>.svc.cluster.local
   ```

4. **有序性问题**
   ```bash
   # 查看 StatefulSet 状态
   kubectl get statefulset
   kubectl describe statefulset <statefulset-name>
   ```

## 📚 相关资源

- [Kubernetes StatefulSet 官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)
- [StatefulSet 最佳实践](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#limitations)

## 🎯 下一步学习

掌握 StatefulSet 后，继续学习：
- [Job 和 CronJob](./18-job-cronjob/README.md) - 批处理任务
- [Storage](./08-storage/README.md) - 存储管理
- [Service](./06-service/README.md) - 服务发现和负载均衡
