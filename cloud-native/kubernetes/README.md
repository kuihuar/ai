# Kubernetes (K8s) 学习指南

## 📚 学习路径

### 第一阶段：基础概念
- [K8s基础概念](./01-basics/README.md) - 容器、Pod、Node等核心概念
- [K8s架构](./02-architecture/README.md) - Master节点、Worker节点、组件详解
- [K8s安装部署](./03-installation/README.md) - 本地环境搭建、生产环境部署

### 第二阶段：核心资源
- [Pod详解](./04-pod/README.md) - Pod生命周期、配置、调度
- [ReplicaSet与Deployment](./05-deployment/README.md) - 应用部署和扩缩容
- [Service与网络](./06-service/README.md) - 服务发现、负载均衡、网络策略
- [ConfigMap与Secret](./07-config/README.md) - 配置管理和敏感信息

### 第三阶段：高级特性
- [存储管理](./08-storage/README.md) - PV、PVC、StorageClass
- [安全机制](./09-security/README.md) - RBAC、NetworkPolicy、PodSecurityPolicy
- [监控与日志](./10-monitoring/README.md) - Prometheus、Grafana、ELK Stack
- [Helm包管理](./11-helm/README.md) - Chart、Release、Repository

### 第四阶段：实战应用
- [微服务部署](./12-microservices/README.md) - 微服务架构在K8s上的实践
- [CI/CD流水线](./13-cicd/README.md) - GitOps、ArgoCD、Jenkins集成
- [故障排查](./14-troubleshooting/README.md) - 常见问题诊断和解决方案
- [性能优化](./15-optimization/README.md) - 资源优化、性能调优

## 🎯 学习目标

通过本学习路径，您将掌握：

1. **基础概念**：理解容器编排、K8s核心概念和架构
2. **资源管理**：熟练使用K8s各种资源对象
3. **网络配置**：掌握K8s网络模型和服务发现
4. **存储管理**：了解持久化存储和动态供应
5. **安全实践**：掌握K8s安全最佳实践
6. **运维技能**：具备K8s集群运维和故障排查能力
7. **实战经验**：通过实际项目积累生产环境经验

## 🛠️ 学习环境

### 本地开发环境
- **Minikube**: 单节点K8s集群，适合本地开发
- **Docker Desktop**: 内置K8s，简单易用
- **Kind**: 使用Docker容器运行K8s集群

### 生产环境
- **云服务商**: AWS EKS、Azure AKS、GCP GKE
- **自建集群**: 使用kubeadm、kops等工具

## 📖 推荐资源

### 官方文档
- [Kubernetes官方文档](https://kubernetes.io/docs/)
- [Kubernetes中文文档](https://kubernetes.io/zh/docs/)

### 在线课程
- [Kubernetes官方教程](https://kubernetes.io/docs/tutorials/)
- [CKA认证课程](https://www.cncf.io/certification/cka/)

### 实践项目
- [Kubernetes示例应用](https://github.com/kubernetes/examples)
- [Kubernetes实战项目](./projects/)

## 🚀 快速开始

1. 安装本地K8s环境（推荐Minikube）
2. 学习基础概念和架构
3. 动手实践Pod、Deployment等资源
4. 逐步深入高级特性和实战应用

---

**开始您的Kubernetes学习之旅吧！** 🎉 