# K8s 中 AI 相关组件及作用

## 一、K8s 原生 AI 基础能力

### 1. Device Plugin 机制
- 作用：让 Kubernetes 识别并调度 GPU、NPU、TPU、昇腾等异构加速设备
- 原理：
  - 设备厂商实现 Device Plugin 接口，以 DaemonSet 运行在每个节点
  - 向 kubelet 注册扩展资源，如 `nvidia.com/gpu`
  - 调度器根据扩展资源进行 Pod 调度
- 常见实现：
  - NVIDIA Device Plugin
  - 华为昇腾 Device Plugin
  - AMD Device Plugin

### 2. Dynamic Resource Allocation (DRA)
- 作用：K8s 1.31+ 正式支持的细粒度资源分配机制
- 解决：GPU 整卡分配利用率低的问题
- 能力：
  - 支持按显存、算力、MIG 实例动态分配
  - 支持 GPU 共享、超配、弹性切分

### 3. 扩展调度机制
- 支持自定义调度器，实现 GPU 亲和性、拓扑感知、任务队列、优先级调度

---

## 二、GPU 高级管理与共享组件

### 1. NVIDIA GPU Operator
- 作用：GPU 节点全生命周期自动化管理
- 功能：
  - 自动安装驱动、CUDA、container-toolkit
  - 配置 MIG 硬件切片
  - 部署 DCGM Exporter 监控
  - 管理 Device Plugin

### 2. MIG (Multi-Instance GPU)
- 作用：硬件级 GPU 切片，单卡划分为多个独立小 GPU
- 适用：A100、H100、L40
- 优势：强隔离、性能几乎无损失

### 3. GPU 软件共享方案
- HAMi
  - 显存+算力分时复用，多 Pod 共享单卡
  - 支持算力限制、优先级
- vGPU / GPU-Shared Manager
  - 时间片轮转共享
  - 适合推理、低优先级任务

### 4. DCGM Exporter
- 作用：GPU 深度监控指标暴露
- 监控内容：
  - GPU 利用率
  - 显存使用
  - 温度、功耗、NVLink 状态
  - 提供 Prometheus 指标

---

## 三、AI 训练任务调度组件

### 1. Volcano（CNCF，AI 训练调度标准）
- 作用：分布式训练/批量作业增强调度
- 核心能力：
  - Gang Scheduling：要么全部启动，要么不启动，避免分布式训练死锁
  - 任务队列、优先级、抢占、公平调度
  - 支持 TensorFlow、PyTorch、MPI 分布式训练
  - 多租户、资源池管理

### 2. Kueue（K8s 官方批量调度）
- 作用：轻量级训练任务队列管理
- 结构：LocalQueue → ClusterQueue → ResourceFlavor
- 支持配额、优先级、抢占、多集群

### 3. Cluster Autoscaler / Karpenter
- 作用：根据 pending Pod 自动扩缩容 GPU 节点
- 实现算力弹性，降低成本

---

## 四、模型训练相关组件

### 1. Kubeflow（端到端 ML 平台）
核心组件：
- Training Operator
  - 支持 TFJob、PyTorchJob、MPIJob、MXNetJob
  - 简化分布式训练编排
- Kubeflow Pipelines (KFP)
  - DAG 工作流、实验管理、流水线可视化
- Katib
  - 自动超参搜索、AutoML、NAS
- Notebook Controller
  - 一键启动 Jupyter，支持 GPU/PVC
- Model Registry
  - 模型版本管理、模型元数据存储

### 2. 训练专用 Operator
- TF Operator：TensorFlow 分布式训练
- PyTorch Operator：PyTorch DDP/FSDP 训练
- MPI Operator：高性能分布式训练
- Spark Operator：数据预处理 + AI 流水线

---

## 五、模型推理服务组件

### 1. KServe（原 KFServing）
- 作用：云原生模型推理服务标准
- 功能：
  - 支持 TensorFlow、PyTorch、ONNX、SKLearn、XGBoost
  - 自动扩缩容、蓝绿发布、金丝雀
  - 预测、批预测、解释器
  - 监控、日志、追踪

### 2. Triton Inference Server
- 作用：NVIDIA 高性能推理引擎
- 特点：
  - 多框架支持
  - 动态批处理、并发
  - 模型仓库、版本管理
  - GPU/CPU 推理优化

### 3. TorchServe / TensorFlow Serving
- 专用框架推理服务
- 轻量、适合单一框架场景

---

## 六、AI 数据与流水线组件

### 1. Argo Workflows
- 作用：K8s 原生 DAG 工作流引擎
- 用于：数据预处理 → 训练 → 验证 → 部署 全链路

### 2. JuiceFS / Fluid
- 作用：AI 数据集缓存、加速
- 解决：训练读取对象存储慢的问题

### 3. MinIO
- 作用：S3 兼容对象存储
- 存储数据集、模型、检查点

---

## 七、监控与可观测性

### 1. Prometheus + Grafana
- 监控：GPU、节点、Pod、训练任务指标
- 看板：GPU 利用率、训练进度、异常告警

### 2. ELK / Loki
- 日志：训练日志、推理日志、错误排查

### 3. Jaeger / Tempo
- 追踪：分布式训练、推理链路追踪