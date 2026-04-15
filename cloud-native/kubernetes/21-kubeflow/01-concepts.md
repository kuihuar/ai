# Kubeflow 核心概念与架构

## 核心组件

- **Kubeflow Pipelines (KFP)**：定义与编排机器学习流水线
- **Notebook**：提供 Jupyter 开发环境
- **Katib**：超参数搜索与 AutoML
- **Training Operator**：统一管理分布式训练任务（如 PyTorch、TFJob）
- **KServe**：模型推理服务部署与灰度发布

## 控制面与数据面

1. 控制面负责声明式配置、任务编排与权限控制  
2. 数据面负责实际训练、推理、数据读写

在生产环境中，通常会把元数据存储、对象存储、模型仓库与 Kubeflow 集成起来。

## 典型工作流

1. 在 Notebook 中进行数据探索和原型开发
2. 把训练流程拆成 Pipeline 组件
3. 使用 Katib 做参数搜索
4. 训练完成后登记模型并发布到 KServe
5. 配置监控与回滚策略

## 与原生 Kubernetes 的关系

- 运行时依赖 Kubernetes 的调度、存储、网络与 RBAC
- 可以结合 `Namespace` 实现团队隔离
- 可与 Argo Workflows、Istio、Prometheus 等生态组件配合

