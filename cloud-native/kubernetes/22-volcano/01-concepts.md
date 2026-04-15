# Volcano 核心概念与调度模型

## 关键对象

- **Queue**：队列，承载租户/团队的资源份额
- **Job**：批任务抽象，支持多角色 Task
- **PodGroup**：任务组，支持 Gang Scheduling
- **PriorityClass**：优先级，用于抢占和排序

## 与 kube-scheduler 的关系

- Volcano 可作为批任务调度器，与默认调度器并存
- 通过插件链实现队列、公平性、抢占、回填等策略

## 典型调度能力

1. **Gang Scheduling**：任务组成员必须同时满足最小可运行数
2. **资源队列与配额**：限制团队资源占用，避免互相干扰
3. **抢占策略**：高优先级任务可抢占低优先级任务资源
4. **回填（Backfill）**：利用空闲碎片资源提高集群利用率

## 适配场景

- 分布式训练（MPI、PyTorch、TensorFlow）
- 大规模 ETL / Spark 批任务
- HPC 队列化作业

