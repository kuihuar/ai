# 云服务提供商详解

## 📚 学习目标

通过本模块学习，您将掌握：
- 主流云服务提供商的核心服务
- AWS、Azure 云原生服务架构
- 容器化、无服务器、存储服务
- 监控、安全、最佳实践
- 多云架构和迁移策略

## 🎯 云服务生态对比

### 1. 服务映射表

| 功能 | AWS | Azure | 说明 |
|------|-----|-------|------|
| **容器编排** | EKS | AKS | Kubernetes 托管服务 |
| **容器运行** | ECS | Container Instances | 容器即服务 |
| **无服务器容器** | Fargate | Container Instances | 无基础设施管理 |
| **无服务器函数** | Lambda | Functions | 事件驱动计算 |
| **对象存储** | S3 | Blob Storage | 可扩展对象存储 |
| **关系数据库** | RDS | SQL Database | 托管关系数据库 |
| **NoSQL 数据库** | DynamoDB | Cosmos DB | 全球分布式数据库 |
| **缓存服务** | ElastiCache | Redis Cache | 内存缓存服务 |
| **CDN** | CloudFront | CDN | 内容分发网络 |
| **负载均衡** | ALB/NLB | Load Balancer | 应用负载均衡 |
| **DNS** | Route 53 | DNS | 域名解析服务 |
| **监控** | CloudWatch | Monitor | 应用监控服务 |
| **日志** | CloudWatch Logs | Log Analytics | 日志聚合分析 |
| **追踪** | X-Ray | Application Insights | 分布式追踪 |
| **安全** | IAM | AAD | 身份和访问管理 |
| **密钥管理** | KMS | Key Vault | 密钥和证书管理 |

### 2. 架构对比

#### AWS 架构
```
AWS 云原生架构
├── 计算层
│   ├── EC2 (虚拟机)
│   ├── EKS (Kubernetes)
│   ├── ECS (容器编排)
│   ├── Fargate (无服务器容器)
│   └── Lambda (无服务器函数)
├── 存储层
│   ├── S3 (对象存储)
│   ├── EBS (块存储)
│   ├── EFS (文件存储)
│   └── FSx (托管文件系统)
├── 数据层
│   ├── RDS (关系数据库)
│   ├── DynamoDB (NoSQL)
│   ├── ElastiCache (缓存)
│   └── Redshift (数据仓库)
└── 网络层
    ├── VPC (虚拟网络)
    ├── ALB/NLB (负载均衡)
    ├── CloudFront (CDN)
    └── Route 53 (DNS)
```

#### Azure 架构
```
Azure 云原生架构
├── 计算层
│   ├── Virtual Machines
│   ├── AKS (Kubernetes)
│   ├── Container Instances
│   ├── App Service
│   └── Functions
├── 存储层
│   ├── Blob Storage
│   ├── Managed Disks
│   ├── Files
│   └── Data Lake Storage
├── 数据层
│   ├── SQL Database
│   ├── Cosmos DB
│   ├── Redis Cache
│   └── Synapse Analytics
└── 网络层
    ├── Virtual Network
    ├── Load Balancer
    ├── CDN
    └── DNS
```

## 🚀 快速开始

### 1. AWS 快速开始

```bash
# 安装 AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# 配置 AWS 凭据
aws configure

# 创建 EKS 集群
eksctl create cluster --name my-cluster --region us-west-2

# 部署应用
kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/application/nginx-app.yaml
```

### 2. Azure 快速开始

```bash
# 安装 Azure CLI
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

# 登录 Azure
az login

# 创建资源组
az group create --name myResourceGroup --location eastus

# 创建 AKS 集群
az aks create --resource-group myResourceGroup --name myAKSCluster --node-count 3

# 获取凭据
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster
```

## 📖 学习路径

### 1. 基础阶段
- [AWS 基础服务](./01-aws/README.md)
- [Azure 基础服务](./02-azure/README.md)
- 云服务概念和架构
- 基础服务配置和部署

### 2. 进阶阶段
- 容器化服务深入
- 无服务器架构设计
- 存储和数据库优化
- 监控和日志分析

### 3. 高级阶段
- 多云架构设计
- 服务迁移策略
- 安全最佳实践
- 成本优化策略

## 🛠️ 实践项目

### 项目1: 多云容器部署
- 在 AWS EKS 和 Azure AKS 上部署相同应用
- 比较性能和成本
- 实现跨云负载均衡

### 项目2: 无服务器 API
- 使用 AWS Lambda 和 Azure Functions 构建 API
- 实现统一的前端接口
- 监控和日志聚合

### 项目3: 数据迁移
- 在 AWS RDS 和 Azure SQL Database 间迁移数据
- 实现数据同步和备份
- 性能对比和优化

## 📚 相关资源

### 官方文档
- [AWS 官方文档](https://docs.aws.amazon.com/)
- [Azure 官方文档](https://docs.microsoft.com/azure/)

### 学习资源
- [AWS 架构中心](https://aws.amazon.com/architecture/)
- [Azure 架构中心](https://docs.microsoft.com/azure/architecture/)

### 工具推荐
- **Terraform**: 多云基础设施管理
- **Ansible**: 配置管理
- **Kubernetes**: 容器编排
- **Helm**: 包管理

---

**掌握多云服务，构建灵活的云原生应用！** 🚀
