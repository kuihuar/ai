# AI 组件对比：可观测性与安全治理

## 对比范围

- 可观测性：Prometheus/Grafana、Loki/ELK、Jaeger/Tempo(+OTel)
- 安全治理：OPA/Kyverno、RBAC/Quota、NetworkPolicy/Cilium、Vault

## 可观测性选型

| 维度 | 方案A | 方案B | 建议 |
| :--- | :--- | :--- | :--- |
| 指标监控 | Prometheus + Grafana | 云厂商托管监控 | 自建可控性强；托管省运维 |
| 日志平台 | Loki | ELK/OpenSearch | 成本优先选 Loki；检索复杂选 ELK |
| 链路追踪 | Tempo | Jaeger | Grafana 一体化优先 Tempo；传统生态选 Jaeger |

## 安全治理选型

| 维度 | 方案A | 方案B | 建议 |
| :--- | :--- | :--- | :--- |
| 策略准入 | Kyverno | OPA Gatekeeper | YAML 友好选 Kyverno；复杂策略选 OPA |
| 网络隔离 | 原生 NetworkPolicy | Cilium | 简单隔离用原生；高级可视化/策略用 Cilium |
| 密钥管理 | K8s Secret | Vault/External Secrets | 生产优先 Vault/External Secrets |

## 联动闭环建议

1. Prometheus 告警定位时间窗
2. Loki/ELK 找错误上下文
3. Tempo/Jaeger 定位链路慢点
4. 回到配置与代码修复并验证

## 落地建议

- 起步：Prometheus + Alertmanager + Loki + RBAC/Quota
- 进阶：OTel + Trace + OPA/Kyverno + Vault
