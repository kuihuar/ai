# Kubeflow 安装与部署

## 部署前准备

- 已有可用 Kubernetes 集群（建议 v1.26+）
- 可用的动态存储类（StorageClass）
- Ingress 或网关方案
- 域名与证书（生产环境推荐）

## 常见部署方式

### 1) 本地学习环境

- 使用 Kind / Minikube 快速体验
- 重点熟悉组件关系和 Pipeline 使用流程

### 2) 云上托管 Kubernetes

- 在 EKS/AKS/GKE 中部署
- 结合云厂商对象存储与负载均衡

## 基础安装步骤（通用）

1. 准备命名空间和基础依赖
2. 部署 Kubeflow manifests
3. 等待核心组件就绪
4. 配置访问入口（Ingress / Gateway）
5. 创建用户与权限策略

## 验证安装

```bash
kubectl get pods -A | rg "kubeflow|kserve|katib"
kubectl get svc -n kubeflow
```

## 生产实践建议

- 将控制面组件与训练负载分节点池
- 启用镜像仓库加速和镜像签名校验
- 配置资源配额、限制范围与审计日志
- 对 Pipeline 元数据和模型资产做备份

