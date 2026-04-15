# AI 组件对比：调度与资源编排

## 对比范围

- Volcano
- Kueue
- YuniKorn

## 一句话定位

- **Volcano**：AI/HPC 训练导向，Gang Scheduling 能力强
- **Kueue**：K8s 官方队列治理，和原生生态融合更自然
- **YuniKorn**：通用多租户队列调度器，偏企业资源治理

## 核心能力对比

| 维度 | Volcano | Kueue | YuniKorn |
| :--- | :--- | :--- | :--- |
| Gang Scheduling | 强 | 通过队列策略间接支持 | 有 |
| 官方生态整合 | 中 | 高 | 中 |
| 多租户配额 | 高 | 高 | 高 |
| AI 训练场景成熟度 | 高 | 中高 | 中 |
| 上手复杂度 | 中 | 中 | 中高 |

## 选型建议

- **分布式训练优先**：优先 Volcano
- **希望官方路线、简化运维**：优先 Kueue
- **企业统一资源池治理**：优先 YuniKorn

## 落地建议

1. 不建议一开始多调度器并存，先单栈稳定
2. 先做队列、优先级、配额模型，再上抢占策略
3. 结合 Cluster Autoscaler/Karpenter 做容量闭环
