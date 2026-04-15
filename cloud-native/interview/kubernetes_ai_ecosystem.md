# Kubernetes AI 架构选型指南 (Cloud Native AI)

本文不是“组件罗列”，而是“选型决策文档”。目标是回答三件事：

1. 不同业务场景该选哪套 AI on K8s 架构  
2. 组件如何组合，避免“全家桶”  
3. 如何分阶段落地，控制风险与成本

---

## 1. 选型方法（先场景，后组件）

建议按以下顺序决策：

1. **业务类型**：训练优先 / 推理优先 / 训练+推理一体  
2. **规模与并发**：单团队、中型、多租户平台  
3. **硬件形态**：NVIDIA-only / 异构（AMD/Ascend）  
4. **成本目标**：性能优先 / 成本优先 / 平衡型  
5. **治理要求**：是否需要强隔离、审计、灰度、SLA

> 原则：先确定“最小可用架构（MVA）”，再增量补组件。

---

## 2. 场景化架构选型（核心）

## 2.1 场景 A：中小团队，推理优先（快速上线）

**目标**：尽快提供稳定 API 推理能力，低运维复杂度。

**推荐组件组合**
- 资源层：`NVIDIA Device Plugin` + `GPU Operator`
- 推理层：`KServe` 或 `vLLM`（二选一起步）
- 调度层：默认 scheduler + `Cluster Autoscaler/Karpenter`
- 存储层：`MinIO`（模型仓库）+ 基础 CSI
- 可观测性：`Prometheus + Grafana + Alertmanager` + `Loki`
- 治理层：`RBAC + Namespace 配额`

**不建议起步就上**
- Volcano/Kueue 双栈并存
- Service Mesh 全套
- 复杂数据编排缓存链路

## 2.2 场景 B：分布式训练优先（吞吐与队列治理）

**目标**：提升训练集群利用率与作业成功率。

**推荐组件组合**
- 资源层：`GPU Operator` + `NFD` + `MIG`（可选）
- 调度层：`Volcano` 或 `Kueue`（二选一为主）
- 训练层：`Kubeflow Training Operator` + `KFP` + `Katib`
- 数据层：`MinIO` + `Fluid/Alluxio` 或 `JuiceFS`
- 可观测性：`Prometheus + DCGM Exporter + Loki/Tempo`
- 治理层：`OPA/Kyverno` + 配额策略

**关键价值**
- Gang Scheduling 避免分布式任务“半启动死锁”
- 队列与配额机制提升多团队公平性

## 2.3 场景 C：企业级多租户 AI 平台（训练+推理一体）

**目标**：平台化、可审计、可运营、可持续扩展。

**推荐组件组合**
- 资源层：Device Plugin + MIG + DRA（按版本成熟度评估）
- 调度层：`Kueue`（队列治理）+ `Volcano`（重训练场景，可按租户隔离）
- MLOps 层：`Kubeflow`（Training/KFP/Katib/Notebook/Registry）
- 推理层：`KServe + Triton/vLLM + Ray Serve`（按模型类型组合）
- 数据层：`MinIO + 数据版本治理（LakeFS/DVC）`
- 可观测性：`Prometheus + Alertmanager + Loki + Tempo/Jaeger + OTel`
- 安全治理：`RBAC + NetworkPolicy/Cilium + Vault + OPA/Kyverno + Service Mesh`

**关键价值**
- 形成“开发 -> 训练 -> 发布 -> 推理 -> 观测 -> 回滚”闭环
- 满足多租户隔离、审计与合规要求

---

## 3. 关键组件选型矩阵（简版）

| 选型维度 | 选项A | 选项B | 选择建议 |
| :--- | :--- | :--- | :--- |
| 批量调度 | Volcano | Kueue | 训练密集且要 Gang 选 Volcano；偏官方队列治理选 Kueue |
| 推理服务 | KServe | vLLM/Triton | 多模型治理与标准化选 KServe；LLM 高吞吐选 vLLM；GPU 推理优化选 Triton |
| 数据加速 | Fluid/Alluxio | JuiceFS | 大规模远端缓存选 Fluid；POSIX 强需求选 JuiceFS |
| 日志方案 | Loki | ELK/OpenSearch | 成本与简单运维优先 Loki；复杂检索分析优先 ELK |
| 服务治理 | 无 Mesh | Istio/Linkerd | 初期可不引入；多租户+灰度+mTLS 强需求再上 |

---

## 4. 分层参考架构（推荐）

## 4.1 基础层（必须）
- Kubernetes 集群
- GPU Operator / Device Plugin / NFD
- MinIO（对象存储）
- Prometheus + Grafana + Alertmanager

## 4.2 能力层（按场景）
- 训练：Volcano/Kueue + Kubeflow Training/KFP/Katib
- 推理：KServe 或 vLLM/Triton
- 数据：Fluid/JuiceFS

## 4.3 治理层（平台化必备）
- RBAC、配额、策略准入（OPA/Kyverno）
- NetworkPolicy/Cilium
- Secret 治理（Vault/External Secrets）
- Trace + 日志体系（Tempo/Jaeger + Loki/ELK）

---

## 5. 落地路线图（避免一步到位失败）

## 阶段 1：最小可用（2~4 周）
- 先交付单场景（训练或推理）
- 上线基础监控与告警
- 建立最小权限与配额

## 阶段 2：性能与成本优化（4~8 周）
- 加入批量调度和弹性扩缩容
- 推进 GPU 共享/切片策略（MIG/MPS）
- 建立容量看板与成本分析

## 阶段 3：平台化治理（8 周+）
- 引入策略准入、密钥治理、网络隔离
- 完成日志与链路追踪闭环
- 建立模型版本治理与灰度发布流程

---

## 6. 风险清单与规避策略

- **风险：组件过多导致运维复杂度失控**  
  对策：每阶段只引入 1~2 个增量能力，先 POC 后生产。

- **风险：高基数指标拖垮监控系统**  
  对策：指标标签治理、录制规则、分层存储。

- **风险：GPU 利用率低，成本不可控**  
  对策：调度队列治理 + 自动扩缩容 + MIG/MPS 策略。

- **风险：多租户安全边界不清**  
  对策：RBAC/配额/准入策略/网络策略组合落地。

---

## 7. 当前推荐基线（给你现阶段）

如果你现在处于“面试与项目并行推进阶段”，推荐先采用：

- 推理：`KServe 或 vLLM`（二选一）
- 训练：`Volcano 或 Kueue`（二选一）
- 数据：`MinIO +（按需）JuiceFS`
- 观测：`Prometheus + Alertmanager + Loki`
- 治理：`RBAC + Namespace 配额 + OPA/Kyverno（最小策略）`

这样可以在保证可用性的同时，保留后续扩展空间。

---

## 8. 附录：AI on K8s 组件分层清单（快速查阅）

## 8.1 基础设施与资源层 (Hardware Enablement)
| 组件名称 | 核心作用 | 场景示例 |
| :--- | :--- | :--- |
| **NVIDIA Device Plugin** | GPU 注册与资源暴露 | `limits: nvidia.com/gpu: 1` |
| **Node Feature Discovery (NFD)** | 节点特征打标 | 调度到特定型号 GPU 节点 |
| **GPU Operator** | GPU 节点全生命周期管理 | 自动化 GPU 集群运维 |
| **MIG (Multi-Instance GPU)** | 硬件切片与隔离 | 多租户推理/训练 |
| **厂商 Device Plugin（AMD/Ascend 等）** | 异构算力接入 | 混合算力集群 |

## 8.2 调度与性能优化层 (High-Density Scheduling)
| 组件名称 | 核心作用 | 场景示例 |
| :--- | :--- | :--- |
| **Volcano / YuniKorn** | 批量调度、Gang Scheduling | 分布式训练 |
| **Kueue** | 官方队列调度与配额治理 | 多租户训练排队 |
| **GPU Sharing / MPS** | GPU 共享与复用 | 推理密度优化 |
| **DRA (Dynamic Resource Allocation)** | 动态资源分配 | 异构资源细粒度调度 |
| **Cluster Autoscaler / Karpenter** | GPU 节点弹性伸缩 | 成本优化 |

## 8.3 工作流与任务编排层 (MLOps)
| 组件名称 | 核心作用 | 场景示例 |
| :--- | :--- | :--- |
| **Kubeflow Training Operator** | 分布式训练管理 | 大模型预训练 |
| **KServe (原 KFServing)** | 模型推理服务治理 | 生产 API 推理 |
| **Katib** | 自动调参 | 模型调优 |
| **Kubeflow Pipelines (KFP)** | 训练流水线编排 | 持续训练发布 |
| **Argo Workflows** | 通用 DAG 工作流 | 数据处理到部署 |
| **Notebook Controller** | Notebook 研发环境 | 数据科学实验 |
| **Model Registry** | 模型版本治理 | 版本追溯与灰度 |
| **Triton Inference Server** | 高性能推理运行时 | GPU 高吞吐推理 |
| **vLLM** | LLM 推理引擎 | OpenAI 兼容服务 |
| **Ray Serve** | 分布式服务编排 | 复杂推理拓扑 |

## 8.4 存储与数据加速层 (Data Orchestration)
| 组件名称 | 核心作用 | 场景示例 |
| :--- | :--- | :--- |
| **Fluid / Alluxio** | 数据缓存与编排 | 降低远端读取延迟 |
| **JuiceFS / CSI-S3** | 共享存储接入 | 并行训练数据访问 |
| **MinIO** | S3 兼容对象存储 | 数据集/模型/Checkpoint |
| **LakeFS / DVC（可选）** | 数据版本治理 | 实验可复现 |

## 8.5 可观测性与运维层 (Observability & Ops)
| 组件名称 | 核心作用 | 场景示例 |
| :--- | :--- | :--- |
| **Prometheus + Grafana** | 指标监控与看板 | GPU 利用率与失败率 |
| **DCGM Exporter** | GPU 深度指标暴露 | 显存/温度/功耗告警 |
| **Alertmanager** | 告警治理 | 告警分级与路由 |
| **Loki / ELK** | 日志观测 | 错误定位与回放 |
| **Jaeger / Tempo + OpenTelemetry** | 链路追踪 | 延迟瓶颈定位 |

## 8.6 安全与多租户治理层 (Security & Governance)
| 组件名称 | 核心作用 | 场景示例 |
| :--- | :--- | :--- |
| **RBAC + Namespace 配额** | 租户隔离与资源边界 | 多团队共享集群 |
| **OPA Gatekeeper / Kyverno** | 策略准入治理 | 强制安全与资源规范 |
| **NetworkPolicy / Cilium** | 网络隔离 | 最小暴露面 |
| **External Secrets / Vault** | 密钥治理 | 外部依赖安全访问 |
| **Service Mesh (Istio/Linkerd)** | 流量治理与零信任 | 灰度、mTLS、限流 |
