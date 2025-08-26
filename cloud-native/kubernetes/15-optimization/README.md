# 性能优化

## 📈 性能优化概述

Kubernetes 性能优化涉及集群、应用、网络、存储等多个层面。通过系统性的优化，可以提高资源利用率、应用性能和用户体验。

## 🎯 优化目标

### 1. 资源利用率
- 提高 CPU 和内存利用率
- 减少资源浪费
- 优化存储使用

### 2. 应用性能
- 减少响应时间
- 提高吞吐量
- 降低延迟

### 3. 集群效率
- 提高调度效率
- 优化网络性能
- 减少运维开销

## 🏗️ 集群优化

### 1. 节点优化

#### 系统参数调优
```bash
# 内核参数优化
cat >> /etc/sysctl.conf << EOF
# 网络优化
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_max_tw_buckets = 5000

# 内存优化
vm.swappiness = 0
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5

# 文件系统优化
fs.file-max = 1000000
fs.inotify.max_user_watches = 1048576
EOF

# 应用配置
sysctl -p
```

#### 容器运行时优化
```yaml
# containerd 配置优化
cat > /etc/containerd/config.toml << EOF
version = 2
root = "/var/lib/containerd"
state = "/run/containerd"

[grpc]
  address = "/run/containerd/containerd.sock"
  max_recv_message_size = 16777216
  max_send_message_size = 16777216

[plugins]
  [plugins."io.containerd.grpc.v1.cri"]
    sandbox_image = "registry.k8s.io/pause:3.6"
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    enable_selinux = false
    enable_tls_streaming = false
    max_container_log_line_size = 16384
    disable_cgroup = false
    disable_apparmor = false
    restrict_oom_score_adj = false
    max_concurrent_downloads = 3
    disable_proc_mount = false
    unset_seccomp_profile = ""
    tolerate_missing_hugetlb_controller = true
    ignore_image_defined_volumes = false

[metrics]
  address = ""
  grpc_histogram = false
EOF
```

### 2. 调度器优化

#### 调度器配置
```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
clientConnection:
  kubeconfig: /etc/kubernetes/scheduler.conf
profiles:
- schedulerName: default-scheduler
  plugins:
    score:
      enabled:
      - name: NodeResourcesBalancedAllocation
      - name: NodePreferAvoidPods
      - name: NodeAffinity
      - name: TaintToleration
      - name: ImageLocality
      - name: InterPodAffinity
      - name: NodeResourcesFit
      disabled:
      - name: "*"
    filter:
      enabled:
      - name: "*"
      disabled:
      - name: "*"
    preFilter:
      enabled:
      - name: "*"
      disabled:
      - name: "*"
    preScore:
      enabled:
      - name: "*"
      disabled:
      - name: "*"
    reserve:
      enabled:
      - name: "*"
      disabled:
      - name: "*"
    permit:
      enabled:
      - name: "*"
      disabled:
      - name: "*"
    bind:
      enabled:
      - name: "*"
      disabled:
      - name: "*"
```

#### 节点亲和性配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: optimized-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: optimized-app
  template:
    metadata:
      labels:
        app: optimized-app
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: node-type
                operator: In
                values:
                - high-performance
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - optimized-app
              topologyKey: kubernetes.io/hostname
      containers:
      - name: app
        image: optimized-app:latest
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
```

## 🚀 应用优化

### 1. 资源管理

#### 资源请求和限制
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: resource-optimized-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: resource-optimized-app
  template:
    metadata:
      labels:
        app: resource-optimized-app
    spec:
      containers:
      - name: app
        image: resource-optimized-app:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        env:
        - name: JAVA_OPTS
          value: "-Xms256m -Xmx512m -XX:+UseG1GC"
```

#### 自动扩缩容
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: optimized-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: optimized-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

### 2. 应用配置优化

#### JVM 优化
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jvm-optimized-app
spec:
  template:
    spec:
      containers:
      - name: app
        image: jvm-optimized-app:latest
        env:
        - name: JAVA_OPTS
          value: |
            -server
            -Xms1g
            -Xmx2g
            -XX:+UseG1GC
            -XX:MaxGCPauseMillis=200
            -XX:+UnlockExperimentalVMOptions
            -XX:+UseStringDeduplication
            -XX:+UseCompressedOops
            -XX:+UseCompressedClassPointers
            -XX:+UseContainerSupport
            -XX:MaxRAMPercentage=75.0
```

#### Node.js 优化
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nodejs-optimized-app
spec:
  template:
    spec:
      containers:
      - name: app
        image: nodejs-optimized-app:latest
        env:
        - name: NODE_ENV
          value: "production"
        - name: NODE_OPTIONS
          value: "--max-old-space-size=1024 --optimize-for-size"
```

## 🌐 网络优化

### 1. 网络策略优化
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: optimized-network-policy
spec:
  podSelector:
    matchLabels:
      app: optimized-app
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: database
    ports:
    - protocol: TCP
      port: 5432
  - to: []
    ports:
    - protocol: UDP
      port: 53
    - protocol: TCP
      port: 53
```

### 2. 服务优化
```yaml
apiVersion: v1
kind: Service
metadata:
  name: optimized-service
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
    service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
spec:
  selector:
    app: optimized-app
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
  type: LoadBalancer
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
```

## 💾 存储优化

### 1. 存储类优化
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
  iops: "3000"
  throughput: "125"
  encrypted: "true"
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
```

### 2. 持久化卷优化
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: optimized-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
  storageClassName: fast-ssd
```

## 📊 监控和调优

### 1. 性能监控
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      evaluation_interval: 15s
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

### 2. 性能指标
```promql
# CPU 使用率
rate(container_cpu_usage_seconds_total{container!=""}[5m])

# 内存使用率
container_memory_usage_bytes{container!=""} / container_spec_memory_limit_bytes{container!=""}

# 网络 I/O
rate(container_network_receive_bytes_total{container!=""}[5m])
rate(container_network_transmit_bytes_total{container!=""}[5m])

# 磁盘 I/O
rate(container_fs_reads_bytes_total{container!=""}[5m])
rate(container_fs_writes_bytes_total{container!=""}[5m])
```

## 🎯 最佳实践

### 1. 资源规划
- 合理设置资源请求和限制
- 监控资源使用趋势
- 定期调整资源配置

### 2. 应用优化
- 使用合适的容器镜像
- 优化应用配置
- 实现健康检查

### 3. 网络优化
- 配置网络策略
- 优化服务发现
- 使用合适的负载均衡

### 4. 存储优化
- 选择合适的存储类型
- 优化 I/O 性能
- 实现数据备份

## 🛠️ 实践练习

### 练习 1：资源优化
1. 分析应用资源使用
2. 调整资源配置
3. 测试性能提升

### 练习 2：网络优化
1. 配置网络策略
2. 优化服务配置
3. 测试网络性能

### 练习 3：存储优化
1. 配置高性能存储
2. 优化 I/O 配置
3. 测试存储性能

## 📚 扩展阅读

- [Kubernetes 性能优化指南](https://kubernetes.io/docs/concepts/cluster-administration/)
- [应用性能优化](https://kubernetes.io/docs/tasks/configure-pod-container/)
- [网络性能优化](https://kubernetes.io/docs/concepts/services-networking/)

## 🎯 下一步

掌握性能优化后，继续学习：
- [生产环境最佳实践](./projects/)
- [高级运维技巧](./projects/) 