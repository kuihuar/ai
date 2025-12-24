

#### prometheus
1. 主动暴露 /metrics
2. 指标转换器exporter
3. 存储使用的自己TSDB
4. 应用 Pod（/metrics） → Prometheus（拉取） → TSDB 存储 → Grafana（PromQL 查询）
#### loki

```mermaid
flowchart LR
    A[日志流] --> B(Distributor)
    B --> C[Ingester: 内存缓冲 + 压缩]
    C --> D[对象存储: S3/GCS]
    D --> E[Querier: 查询处理]
    
```
- 应用 Pod（stdout/stderr） → Promtail（采集+标签） → Loki 存储 → Grafana（LogQL 查询）
### 架构

```mermaid
flowchart TB
    subgraph 业务应用
        A[业务代码] -->|自定义指标| B(Prometheus Client)
        A -->|业务日志文件| C[(/app/logs)]
        A -->|stdout/stderr| D[容器日志]
    end

    subgraph K8s集群
        B -->|暴露 /metrics| E(Prometheus)
        C -->|挂载卷| F[Fluent Bit]
        D --> G[Promtail]
        F -->|日志流| H(Loki)
        G --> H
        E -->|指标存储| I[Prometheus TSDB]
    end

    H & I -->|查询| J[Grafana]

```    


#### 综合trace后

```mermaid
flowchart TB
    %% 业务应用层
    subgraph 业务应用
        A[Service] -->|自动埋点| B(OTel Instrumentation)
        A -->|手动埋点| C[metrics端点]
        A -->|写入文件| D[app/logs]
        A -->|stdout/stderr| E[容器日志]
    end

    %% 数据采集层
    subgraph 数据采集层
        B -->|OTLP 指标/日志/追踪| F[OTel Collector]
        C -->|被拉取| G[Prometheus]
        D -->|采集| H[Fluent Bit]
        E -->|采集| I[Promtail]
        H -->|日志流| J[Loki]
        I -->|日志流| J
        F -->|指标导出| G
        F -->|日志导出| J
        F -->|追踪导出| K[Tempo]
    end

    %% 存储层
    subgraph 存储层
        G -->|指标存储| L[Prometheus TSDB]
        J -->|日志存储| M[Loki Storage]
        K -->|追踪存储| N[Tempo Backend]
    end

    %% 可视化层
    subgraph 可视化层
        L & M & N -->|查询| O[Grafana]
    end

    %% 样式调整
    style A fill:#f9f,stroke:#333
    style B fill:#bbf,stroke:#333
    style C fill:#f96,stroke:#333
    style D fill:#f96,stroke:#333
    style E fill:#f96,stroke:#333
    style F fill:#6bf,stroke:#333
    style G fill:#f66,stroke:#333
    style H fill:#6f6,stroke:#333
    style I fill:#6f6,stroke:#333
    style J fill:#9f9,stroke:#333
    style K fill:#99f,stroke:#333
    style L fill:#f66,stroke:#333
    style M fill:#9f9,stroke:#333
    style N fill:#99f,stroke:#333
    style O fill:#ff9,stroke:#333
```

```mermaid
sequenceDiagram
    participant App as 业务应用
    participant OTel as OTel Collector
    participant Prom as Prometheus
    participant Fluent as Fluent Bit
    participant Loki as Loki
    participant Tempo as Tempo
    participant Grafana as Grafana

    App->>OTel: OTLP 指标/日志/追踪
    App->>Prom: /metrics (Prometheus 格式)
    App->>Fluent: 业务日志文件
    App->>Loki: 容器 stdout/stderr

    OTel->>Prom: 转发指标（Prometheus Remote Write）
    OTel->>Loki: 转发日志（Loki Push API）
    OTel->>Tempo: 转发追踪（OTLP/gRPC）

    Prom-->>Grafana: 提供指标查询
    Loki-->>Grafana: 提供日志查询
    Tempo-->>Grafana: 提供追踪查询

    Note right of Grafana: 用户通过 Grafana 统一界面<br>关联分析指标、日志、追踪
```


```mermaid

flowchart TB
    subgraph Go微服务
        A[业务代码] -->|埋点| B(OTel SDK)
        A -->|日志| C[logrus/zap]
        A -->|暴露/metrics| D[Prometheus Client]
    end

    subgraph 采集与存储
        B -->|OTLP 数据| E[OTel Collector]
        C -->|日志输出| F[Promtail/Fluent Bit]
        D -->|拉取| G[Prometheus]
        F -->|日志流| H[Loki]
        E -->|指标| G
        E -->|日志| H
        E -->|追踪| I[Tempo]
    end

    subgraph 可视化
        G & H & I --> J[Grafana]
    end

```  

### 指标、日志、追踪的统一管理


#### 架构

```mermaid

flowchart TB
    subgraph Go微服务
        A[业务代码] -->|埋点| B(OTel SDK)
        A -->|日志| C[logrus/zap]
        A -->|暴露/metrics| D[Prometheus Client]
    end

    subgraph 采集与存储
        B -->|OTLP 数据| E[OTel Collector]
        C -->|日志输出| F[Promtail/Fluent Bit]
        D -->|拉取| G[Prometheus]
        F -->|日志流| H[Loki]
        E -->|指标| G
        E -->|日志| H
        E -->|追踪| I[Tempo]
    end

    subgraph 可视化
        G & H & I --> J[Grafana]
    end
```
#### 集成方案

1. (1) 指标监控（Prometheus）
   - prometheus client：暴露 Prometheus 指标  /metrics
   - K8s 配置（ServiceMonitor） /metrics
2. (2) 日志收集（Loki） 
   - 结构化日志（logrus）
   - Promtail 配置收集
3. (3) 分布式追踪（Tempo + OpenTelemetry）
   - OTel SDK 埋点（otel.Tracer）
4. (4) 统一采集（OpenTelemetry Collector）
``` yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

exporters:
  prometheusremotewrite:
    endpoint: "http://prometheus:9090/api/v1/write"
  loki:
    endpoint: "http://loki:3100/loki/api/v1/push"
  otlp:
    endpoint: "tempo:4317"

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp]
    metrics:
      receivers: [otlp]
      exporters: [prometheusremotewrite]
    logs:
      receivers: [otlp]
      exporters: [loki]
      processors: [batch]
```      
#### 部署编排
1. (1)示例      
``` yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-service
spec:
  template:
    spec:
      containers:
      - name: go-service
        image: my-go-app:latest
        ports:
        - containerPort: 8080
        env:
        - name: OTEL_SERVICE_NAME
          value: "go-service"
        - name: OTEL_EXPORTER_OTLP_ENDPOINT
          value: "http://otel-collector:4317"
        volumeMounts:
        - name: logs
          mountPath: /var/log/go-service
      volumes:
      - name: logs
        emptyDir: {}
```
2. (2) Sidecar 模式（日志采集）


#### 数据关联与 Grafana 展示
(1) 标签一致性
指标：http_requests_total{service="go-service"}

日志：{job="go-service", trace_id="abc123"}

追踪：tempo_query(trace_id="abc123")

(2) Grafana 仪表盘
指标面板：PromQL 查询 rate(http_requests_total[5m])

日志面板：LogQL 查询 {job="go-service"} |= "error"

追踪面板：直接关联 TraceID


```mermaid
flowchart TB
    subgraph OTel生态
        A[OTel API/SDK] -->|埋点| B(应用程序)
        B -->|生成数据| C[OTel Collector]
        C -->|导出指标| D[Prometheus]
        C -->|导出日志| E[Loki]
        C -->|导出追踪| F[Tempo/Jaeger]
        C -->|导出到云服务| G[AWS CloudWatch/Datadog]
    end

```


zap/zerolog/logrus 日志生成库（应用层）
OTel Collector	日志收集与路由中枢
Loki	日志存储与查询引擎
Promtail/Fluent Bit	日志采集代理

```mermaid
flowchart TB
    subgraph 应用代码
        A[zap/zerolog/logrus] -->|打印JSON日志| B[(日志文件)]
        A -->|直接发送| C[OTel SDK]
    end

    subgraph 采集层
        B -->|文件日志| D[Fluent Bit]
        C -->|OTLP 日志| E[OTel Collector]
        D -->|转发| E
    end

    subgraph 存储层
        E -->|路由和导出| F[Loki]
    end

    subgraph 查询层
        F -->|LogQL| G[Grafana]
    end
```

组件选型建议
|需求|	推荐组合|
|--|--|
|高性能微服务日志|	zap/zerolog + OTel Collector + Loki|
|Kubernetes 环境|	logrus + Promtail + Loki|
|混合环境（VM+K8s）|	zerolog + Fluent Bit + OTel Collector + Loki|
|需要复杂日志分析|	任何日志库 + Elasticsearch（替代 Loki）|


维度	高性能微服务日志	Kubernetes 环境日志
核心目标	极致性能（低延迟、高吞吐）	自动化管理（动态发现、弹性伸缩）
日志来源	业务逻辑高频日志（如每请求多次打印）	容器标准输出（stdout/stderr）+ 少量文件日志
采集挑战	避免日志库成为性能瓶颈	处理 Pod 动态创建销毁的日志流
典型组件	zap/zerolog + OTel Collector + Kafka + Loki	logrus + Promtail + Loki
优化重点	减少内存分配、异步写入	自动标签注入、日志路由

需求	高性能微服务方案	K8s 通用方案
日志分级	代码控制（Debug 日志采样）	采集器过滤（丢弃 Info 级）
多租户隔离	业务标签区分（如 tenant_id）	Loki 的 X-Scope-OrgID
故障排查	依赖 TraceID 跨服务关联	依赖 Pod 名称和事件时间



#### 若服务 既是高性能微服务又运行在 K8s，可融合两者优势

```mermaid
flowchart TB
    subgraph 应用层
        A[Go服务] -->|zap JSON日志| B[(EmptyDir Volume)]
    end
    subgraph K8s层
        B -->|挂载| C[Fluent Bit]
        C -->|添加k8s标签| D[Kafka]
        D --> E[OTel Collector]
        E --> F[Loki]
    end
```