# 可观察性知识点整理（Prometheus 方向）

本文从可观察性角度整理 Prometheus 相关知识点，并结合当前云原生虚拟化场景（KubeVirt/Longhorn/Multus）给出排障思路。

---

## 1. 可观察性三支柱

## 1.1 Metrics（指标）

- 关注“趋势”和“聚合”，适合看容量、性能、错误率
- 典型工具：Prometheus + Alertmanager + Grafana
- 关键词：采样、时序数据、标签（label）、聚合查询

## 1.2 Logs（日志）

- 关注“事件细节”，适合定位具体报错与上下文
- 常见栈：Loki/ELK
- 关键词：结构化日志、TraceID 关联、日志采集与索引

## 1.3 Traces（链路追踪）

- 关注“请求路径与耗时分解”
- 常见栈：Jaeger/Tempo + OpenTelemetry
- 关键词：span、parent-child、采样率

---

## 2. Prometheus 核心原理（面试高频）

1. **拉模型（Pull）**
   - Prometheus 主动从目标的 `/metrics` 抓取数据
2. **时序存储（TSDB）**
   - 指标按 `metric + labels + timestamp + value` 存储
3. **PromQL 查询**
   - 提供聚合、分组、时间窗口函数（`rate`、`sum by`、`histogram_quantile`）
4. **告警链路**
   - Prometheus 触发规则 -> Alertmanager 分组/抑制/路由 -> 通知渠道

---

## 3. Prometheus 关键组件

- `Prometheus Server`：抓取、存储、查询、规则评估
- `Alertmanager`：告警去重、分组、抑制、路由
- `node-exporter`：主机资源指标
- `kube-state-metrics`：K8s 对象状态指标
- `Grafana`：可视化与看板

在 Kubernetes 里常见部署方式：

- kube-prometheus-stack（Helm）
- Prometheus Operator（通过 `ServiceMonitor/PodMonitor` 管理抓取目标）

---

## 3.1 Prometheus 组件知识点（展开版）

## 3.1.1 Prometheus Server 内部模块

- `scrape manager`：管理抓取任务调度与并发抓取
- `service discovery`：自动发现目标（K8s/Consul/static）
- `TSDB`：本地时序存储（WAL + block compaction）
- `rule manager`：评估 recording rules 和 alert rules
- `query engine`：执行 PromQL 即时查询

知识点：

- 写入路径：`scrape -> relabel -> ingest -> WAL -> memory -> block`
- 查询路径：`PromQL -> 查询引擎 -> TSDB block -> 结果聚合`
- 性能瓶颈常见在：高基数标签、宽查询窗口、大量 regex 匹配

## 3.1.2 Operator 生态组件（K8s 常用）

- `Prometheus Operator`：管理 Prometheus/Alertmanager CR
- `ServiceMonitor`：按 Service 发现并抓取目标
- `PodMonitor`：按 Pod 发现并抓取目标
- `PrometheusRule`：托管告警/录制规则
- `AlertmanagerConfig`：细粒度告警路由

知识点：

- 业务团队通常只需要维护 `ServiceMonitor + PrometheusRule`
- 平台团队维护 Prometheus 生命周期和容量规划

## 3.1.3 数据接入与长期存储

- `remote_write`：把指标写到远端存储（Thanos/Cortex/Mimir/VictoriaMetrics）
- `remote_read`：查询远端历史数据
- `federation`：多 Prometheus 分层聚合查询

知识点：

- 短期高频查询用本地 TSDB
- 长期保留与跨集群查询用远端时序后端

## 3.1.4 Exporter 体系

- 基础设施：`node-exporter`
- K8s 资源：`kube-state-metrics`
- 容器层：`cAdvisor`（常由 kubelet 暴露）
- 应用层：业务自定义 `/metrics`

知识点：

- Exporter 负责“暴露”，Prometheus 负责“抓取与存储”
- 优先标准化指标，减少自定义采集成本

---

## 3.2 日志组件全景（Logs）

## 3.2.1 采集层

- `Fluent Bit` / `Promtail` / `Vector`
- 负责从容器 stdout、文件、systemd 读取日志并打标签

## 3.2.2 存储与检索层

- `Loki`：标签索引 + 对象存储日志内容（成本友好）
- `Elasticsearch/OpenSearch`：全文检索强，运维成本更高

## 3.2.3 查询与分析层

- `Grafana`（Loki 数据源）
- `Kibana`（Elasticsearch 数据源）

## 3.2.4 告警与关联

- 基于日志关键字或频次告警（如 5xx 激增）
- 与 Metrics 联动：指标告警触发后自动跳转日志查询

日志体系关键知识点：

- 统一日志字段（service、namespace、pod、trace_id）
- 结构化 JSON 日志优先
- 日志保留策略按冷热分层

---

## 3.3 链路追踪组件全景（Traces）

## 3.3.1 采集与注入

- `OpenTelemetry SDK`：应用埋点与上下文传播
- `OpenTelemetry Collector`：统一接收、处理、转发 traces/metrics/logs

## 3.3.2 存储与查询

- `Jaeger`：经典链路追踪平台
- `Tempo`：与 Grafana 结合紧密，成本可控

## 3.3.3 关联能力

- Trace -> Logs（通过 trace_id）
- Trace -> Metrics（查看慢请求的指标背景）

链路追踪关键知识点：

- 采样策略（head/tail sampling）直接影响成本和可见性
- 跨服务上下文传播（W3C Trace Context）是追踪成功前提

---

## 3.4 三类组件如何联动（实战）

1. Prometheus 发现 `reconcile_errors_total` 突增并触发告警  
2. 跳转 Loki 查询同时间窗 `reconcile_id`、`vm_name` 相关错误日志  
3. 再看 Trace（若接入 OTel）定位慢点发生在网络、存储还是 API 调用  
4. 最后回到代码与对象状态确认根因

这就是“Metrics 定位范围、Logs 给细节、Trace 给路径”的闭环。

---

## 4. 指标设计最佳实践

## 4.1 指标类型

- Counter：单调递增（请求总数、错误总数）
- Gauge：可增可减（并发数、队列长度、资源使用）
- Histogram/Summary：分位数与延迟分布

## 4.2 Label 设计

- 保持低基数（避免 label 爆炸）
- 避免用户 ID、请求 ID 这类高基数字段
- 用 `namespace/workload/component` 等稳定维度

## 4.3 命名建议

- 格式：`<domain>_<subsystem>_<metric>_<unit>`
- 示例：
  - `vmoperator_reconcile_duration_seconds`
  - `vmoperator_reconcile_total`
  - `vmoperator_reconcile_errors_total`

---

## 5. PromQL 常用查询模板

## 5.1 QPS/错误率

```promql
sum(rate(http_requests_total[5m]))
sum(rate(http_requests_total{code=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

## 5.2 Reconcile 观测（Operator 场景）

```promql
sum(rate(controller_runtime_reconcile_total[5m])) by (controller)
sum(rate(controller_runtime_reconcile_errors_total[5m])) by (controller)
histogram_quantile(0.95, sum(rate(controller_runtime_reconcile_time_seconds_bucket[5m])) by (le, controller))
```

## 5.3 资源与容量

```promql
sum(container_memory_working_set_bytes{namespace="kubevirt"})
sum(rate(container_cpu_usage_seconds_total{namespace="kubevirt"}[5m]))
```

---

## 6. 告警规则设计思路

## 6.1 四层告警

1. **SLO 告警**：错误率、延迟、可用性
2. **平台告警**：节点资源、API Server、etcd、CoreDNS
3. **业务告警**：Reconcile error、队列堆积、VM Pending 超时
4. **依赖告警**：Longhorn 卷异常、Multus/NAD 失败、KubeVirt 组件异常

## 6.2 示例规则（思路）

- Reconcile 错误率持续升高（5m）
- VMI 长时间 Pending（>10m）
- PVC 未绑定持续超时
- virt-launcher 启动失败次数异常

---

## 7. 虚拟化场景下的可观察性重点

## 7.1 KubeVirt

重点看：

- VM/VMI phase 分布（Running/Pending/Error）
- virt-controller/virt-handler 组件健康
- 迁移成功率与迁移耗时

## 7.2 Longhorn（存储）

重点看：

- 卷健康状态、副本健康、attach/mount 失败
- PVC 绑定时长、扩容耗时

## 7.3 Multus + NMState（网络）

重点看：

- NAD 变更失败次数
- CNI attach 错误事件
- 节点网络策略（NNCP/NNS）收敛状态

---

## 8. 从可观察性做排障（标准路径）

1. **先看告警面板**：确定是容量、错误率还是可用性告警
2. **再看指标趋势**：定位峰值时间窗（CPU/内存/Reconcile errors）
3. **关联日志**：按时间窗和对象名（VM/VMI/PVC）搜日志
4. **落到对象状态**：`kubectl describe` 看 events + conditions
5. **回到代码路径**：定位对应 reconcile 分支与依赖调用

---

## 9. 与当前项目可结合的观测点（建议）

结合 `vmoperator`，建议补充这些业务指标：

- `wukong_reconcile_total{result}`
- `wukong_reconcile_duration_seconds{phase}`
- `wukong_network_reconcile_errors_total`
- `wukong_storage_reconcile_errors_total`
- `wukong_vm_phase_total{phase}`
- `wukong_volume_bound_latency_seconds`

日志建议统一字段：

- `wukong_name`
- `namespace`
- `vm_name`
- `phase`
- `reconcile_id`

---

## 10. 面试回答模板（30 秒）

「可观察性我会按指标、日志、链路三层讲。Prometheus 负责拉取时序指标和规则告警，Grafana做可视化，Alertmanager做告警路由。落地时我会先定义 SLO 指标，再补平台与依赖指标，比如 KubeVirt 的 VMI 状态、Longhorn 卷健康、Operator Reconcile 错误率。排障流程是先看告警和趋势，再关联日志与对象状态，最后定位到控制器代码分支，形成闭环。」 
