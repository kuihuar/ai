# Kubernetes 实战项目

## 📚 项目概述

本目录包含 Kubernetes 实战项目，通过实际案例来巩固和应用所学知识。每个项目都包含完整的部署配置、最佳实践和故障排查指南。

## 🎯 项目列表

### 1. 基础项目
- [Web 应用部署](./01-web-app/) - 简单的 Web 应用部署
- [数据库应用](./02-database/) - MySQL/PostgreSQL 数据库部署
- [缓存服务](./03-cache/) - Redis 缓存服务部署

### 2. 中级项目
- [微服务架构](./04-microservices/) - 完整的微服务应用
- [API 网关](./05-api-gateway/) - Kong/Envoy API 网关
- [监控栈](./06-monitoring/) - Prometheus + Grafana 监控

### 3. 高级项目
- [GitOps 部署](./07-gitops/) - ArgoCD GitOps 实践
- [服务网格](./08-service-mesh/) - Istio 服务网格
- [机器学习平台](./09-ml-platform/) - Kubeflow 机器学习平台

## 🛠️ 项目结构

每个项目都包含以下内容：

```
project-name/
├── README.md              # 项目说明
├── k8s/                   # Kubernetes 配置
│   ├── namespace.yaml     # 命名空间
│   ├── configmap.yaml     # 配置
│   ├── secret.yaml        # 密钥
│   ├── deployment.yaml    # 部署
│   ├── service.yaml       # 服务
│   ├── ingress.yaml       # 入口
│   └── pvc.yaml          # 存储
├── helm/                  # Helm Charts
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
├── scripts/               # 脚本文件
│   ├── deploy.sh         # 部署脚本
│   ├── test.sh           # 测试脚本
│   └── cleanup.sh        # 清理脚本
├── docs/                  # 文档
│   ├── architecture.md   # 架构说明
│   ├── troubleshooting.md # 故障排查
│   └── best-practices.md # 最佳实践
└── examples/              # 示例代码
    ├── dockerfile        # Dockerfile
    ├── docker-compose.yml # Docker Compose
    └── app/              # 应用代码
```

## 🚀 快速开始

### 1. 环境准备
```bash
# 确保有可用的 Kubernetes 集群
kubectl cluster-info

# 安装必要工具
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### 2. 项目部署
```bash
# 进入项目目录
cd projects/01-web-app

# 部署项目
./scripts/deploy.sh

# 验证部署
./scripts/test.sh
```

### 3. 项目清理
```bash
# 清理项目资源
./scripts/cleanup.sh
```

## 📖 学习路径

### 初学者
1. 从 [Web 应用部署](./01-web-app/) 开始
2. 学习 [数据库应用](./02-database/)
3. 实践 [缓存服务](./03-cache/)

### 进阶者
1. 深入 [微服务架构](./04-microservices/)
2. 学习 [API 网关](./05-api-gateway/)
3. 配置 [监控栈](./06-monitoring/)

### 高级用户
1. 实践 [GitOps 部署](./07-gitops/)
2. 探索 [服务网格](./08-service-mesh/)
3. 构建 [机器学习平台](./09-ml-platform/)

## 🎯 项目目标

通过实战项目，您将掌握：

1. **实际部署技能** - 真实应用的部署和管理
2. **问题解决能力** - 常见问题的诊断和解决
3. **最佳实践** - 生产环境的最佳实践
4. **架构设计** - 复杂系统的架构设计
5. **运维技能** - 日常运维和故障处理

## 📚 扩展资源

- [Kubernetes 官方示例](https://github.com/kubernetes/examples)
- [Helm Charts](https://github.com/helm/charts)
- [Kubernetes 最佳实践](https://kubernetes.io/docs/concepts/)

## 🤝 贡献指南

欢迎提交项目改进建议和新的实战项目！

1. Fork 本仓库
2. 创建特性分支
3. 提交更改
4. 发起 Pull Request

---

**开始您的 Kubernetes 实战之旅吧！** 🎉 