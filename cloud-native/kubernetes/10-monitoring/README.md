# 监控与日志

## 📊 监控概述

Kubernetes 集群的监控是确保应用稳定运行的关键。监控包括基础设施监控、应用监控、日志收集和分析等多个方面。

## 🎯 监控架构

### 1. 监控层次
- **基础设施层**: 节点、网络、存储监控
- **容器层**: Pod、容器资源使用监控
- **应用层**: 应用性能、业务指标监控
- **用户层**: 用户体验、SLA 监控

### 2. 监控组件
```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   Metrics   │ │   Logs      │ │   Traces    │            │
│  │   Server    │ │   Collector │ │   Collector │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │  Prometheus │ │   Grafana   │ │   Alert     │            │
│  │             │ │             │ │  Manager    │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

## 📈 Prometheus 监控

### 1. Prometheus 简介
Prometheus 是一个开源的监控和告警系统，特别适合 Kubernetes 环境。

**特点：**
- 多维度数据模型
- 强大的查询语言 PromQL
- 支持服务发现
- 告警管理

### 2. Prometheus 部署
```yaml
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
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:latest
        ports:
        - containerPort: 9090
        volumeMounts:
        - name: config
          mountPath: /etc/prometheus
      volumes:
      - name: config
        configMap:
          name: prometheus-config
```

### 3. 关键指标
```promql
# CPU 使用率
rate(container_cpu_usage_seconds_total{container!=""}[5m])

# 内存使用率
container_memory_usage_bytes{container!=""} / container_spec_memory_limit_bytes{container!=""}

# Pod 重启次数
increase(kube_pod_container_status_restarts_total[1h])

# 节点状态
kube_node_status_condition{condition="Ready"}
```

## 📊 Grafana 可视化

### 1. Grafana 部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - name: grafana
        image: grafana/grafana:latest
        ports:
        - containerPort: 3000
        env:
        - name: GF_SECURITY_ADMIN_PASSWORD
          value: "admin"
        volumeMounts:
        - name: grafana-storage
          mountPath: /var/lib/grafana
      volumes:
      - name: grafana-storage
        persistentVolumeClaim:
          claimName: grafana-pvc
```

### 2. 常用仪表板
- **Kubernetes 集群概览**: 节点状态、Pod 数量、资源使用
- **应用性能监控**: 响应时间、错误率、吞吐量
- **基础设施监控**: CPU、内存、磁盘、网络

## 📝 日志收集

### 1. ELK Stack
Elasticsearch + Logstash + Kibana 的日志收集方案。

#### Elasticsearch
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: elasticsearch
spec:
  serviceName: elasticsearch
  replicas: 3
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      containers:
      - name: elasticsearch
        image: docker.elastic.co/elasticsearch/elasticsearch:7.17.0
        env:
        - name: cluster.name
          value: "k8s-logs"
        - name: node.name
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: discovery.seed_hosts
          value: "elasticsearch-0.elasticsearch,elasticsearch-1.elasticsearch,elasticsearch-2.elasticsearch"
        - name: cluster.initial_master_nodes
          value: "elasticsearch-0,elasticsearch-1,elasticsearch-2"
        ports:
        - containerPort: 9200
          name: http
        - containerPort: 9300
          name: transport
```

#### Logstash
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logstash
spec:
  replicas: 2
  selector:
    matchLabels:
      app: logstash
  template:
    metadata:
      labels:
        app: logstash
    spec:
      containers:
      - name: logstash
        image: docker.elastic.co/logstash/logstash:7.17.0
        ports:
        - containerPort: 5044
        volumeMounts:
        - name: logstash-config
          mountPath: /usr/share/logstash/pipeline
      volumes:
      - name: logstash-config
        configMap:
          name: logstash-config
```

#### Kibana
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kibana
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kibana
  template:
    metadata:
      labels:
        app: kibana
    spec:
      containers:
      - name: kibana
        image: docker.elastic.co/kibana/kibana:7.17.0
        env:
        - name: ELASTICSEARCH_HOSTS
          value: "http://elasticsearch:9200"
        ports:
        - containerPort: 5601
```

### 2. Fluentd
轻量级的日志收集器。

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
spec:
  selector:
    matchLabels:
      name: fluentd
  template:
    metadata:
      labels:
        name: fluentd
    spec:
      serviceAccount: fluentd
      containers:
      - name: fluentd
        image: fluent/fluentd-kubernetes-daemonset:v1.14-debian-elasticsearch7-1
        env:
        - name: FLUENT_ELASTICSEARCH_HOST
          value: "elasticsearch"
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

## 🚨 告警管理

### 1. Alertmanager
Prometheus 的告警管理器。

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: alertmanager-config
data:
  alertmanager.yml: |
    global:
      smtp_smarthost: 'localhost:587'
      smtp_from: 'alertmanager@example.com'
    route:
      group_by: ['alertname']
      group_wait: 10s
      group_interval: 10s
      repeat_interval: 1h
      receiver: 'web.hook'
    receivers:
    - name: 'web.hook'
      webhook_configs:
      - url: 'http://127.0.0.1:5001/'
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alertmanager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: alertmanager
  template:
    metadata:
      labels:
        app: alertmanager
    spec:
      containers:
      - name: alertmanager
        image: prom/alertmanager:latest
        ports:
        - containerPort: 9093
        volumeMounts:
        - name: config
          mountPath: /etc/alertmanager
      volumes:
      - name: config
        configMap:
          name: alertmanager-config
```

### 2. 告警规则
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-rules
data:
  rules.yml: |
    groups:
    - name: kubernetes
      rules:
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage on {{ $labels.instance }}"
          description: "CPU usage is above 80% for 5 minutes"
      
      - alert: HighMemoryUsage
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage on {{ $labels.instance }}"
          description: "Memory usage is above 80% for 5 minutes"
```

## 🔍 分布式追踪

### 1. Jaeger
分布式追踪系统。

```yaml
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
```

### 2. 应用集成
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        env:
        - name: JAEGER_AGENT_HOST
          value: "jaeger"
        - name: JAEGER_AGENT_PORT
          value: "6831"
```

## 🛠️ 监控工具

### 1. kube-state-metrics
Kubernetes 资源状态指标。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-state-metrics
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kube-state-metrics
  template:
    metadata:
      labels:
        app: kube-state-metrics
    spec:
      containers:
      - name: kube-state-metrics
        image: k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.5.0
        ports:
        - containerPort: 8080
```

### 2. node-exporter
节点指标收集器。

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-exporter
spec:
  selector:
    matchLabels:
      app: node-exporter
  template:
    metadata:
      labels:
        app: node-exporter
    spec:
      containers:
      - name: node-exporter
        image: prom/node-exporter:latest
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

## 🎯 最佳实践

### 1. 监控策略
- 监控关键指标
- 设置合理的告警阈值
- 定期审查监控配置

### 2. 日志管理
- 统一日志格式
- 设置日志保留策略
- 监控日志量

### 3. 性能优化
- 合理配置资源
- 优化查询性能
- 监控系统开销

### 4. 安全考虑
- 保护监控数据
- 控制访问权限
- 加密敏感信息

## 🛠️ 实践练习

### 练习 1：基础监控
1. 部署 Prometheus 和 Grafana
2. 配置数据源
3. 创建仪表板

### 练习 2：日志收集
1. 部署 ELK Stack
2. 配置日志收集
3. 创建日志分析

### 练习 3：告警配置
1. 配置告警规则
2. 设置通知渠道
3. 测试告警功能

## 📚 扩展阅读

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [ELK Stack 官方文档](https://www.elastic.co/guide/)

## 🎯 下一步

掌握监控与日志后，继续学习：
- [Helm包管理](./11-helm/README.md)
- [微服务部署](./12-microservices/README.md) 