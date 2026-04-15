# AI 组件对比：训练与 MLOps 编排

## 对比范围

- Kubeflow (Training Operator / KFP / Katib)
- Argo Workflows
- Ray (Train + Serve 生态)

## 一句话定位

- **Kubeflow**：端到端 ML 平台，组件齐全
- **Argo Workflows**：通用工作流编排引擎，灵活可控
- **Ray 生态**：分布式计算与服务一体，偏 Python/LLM 工程化

## 核心能力对比

| 维度 | Kubeflow | Argo Workflows | Ray |
| :--- | :--- | :--- | :--- |
| 训练任务 CRD 生态 | 强 | 中 | 中高 |
| 流水线管理 | 强（KFP） | 强 | 中 |
| 自动调参 | 强（Katib） | 需自建 | 中 |
| 学习与运维成本 | 高 | 中 | 中 |
| 适合团队 | 平台化团队 | DevOps 团队 | 算法/工程混合团队 |

## 选型建议

- **要完整 MLOps 平台**：Kubeflow
- **要轻量流水线与通用编排**：Argo
- **要分布式计算 + 在线服务一体**：Ray

## 落地建议

1. 先统一实验、模型、数据元信息规范
2. 先固化最小流水线（训练->评估->发布）
3. 自动调参与复杂编排放第二阶段
