# Kubernetes 安装部署

## 🎯 安装方式概览

Kubernetes 有多种安装方式，根据使用场景选择合适的方案：

### 1. 本地开发环境
- **Minikube**: 单节点集群，适合学习和开发
- **Docker Desktop**: 内置 K8s，简单易用
- **Kind**: 使用 Docker 容器运行 K8s 集群

### 2. 生产环境
- **云服务商**: AWS EKS、Azure AKS、GCP GKE
- **自建集群**: kubeadm、kops、Rancher

## 🛠️ 本地开发环境安装

### Minikube 安装

**优点：**
- 轻量级，资源占用少
- 支持多种驱动（Docker、VirtualBox、KVM）
- 适合学习和测试

**安装步骤：**
```bash
# 下载 Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 启动集群
minikube start --driver=docker

# 验证安装
kubectl cluster-info
```

**常用命令：**
```bash
# 启动集群
minikube start

# 停止集群
minikube stop

# 删除集群
minikube delete

# 查看状态
minikube status

# 打开仪表板
minikube dashboard
```

### Docker Desktop

**优点：**
- 一键安装，配置简单
- 与 Docker 集成良好
- 支持 Windows、macOS、Linux

**安装步骤：**
1. 下载并安装 Docker Desktop
2. 在设置中启用 Kubernetes
3. 等待集群启动完成

### Kind (Kubernetes in Docker)

**优点：**
- 使用 Docker 容器运行 K8s
- 支持多节点集群
- 适合 CI/CD 环境

**安装步骤：**
```bash
# 安装 Kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# 创建集群
kind create cluster

# 验证安装
kubectl cluster-info
```

## ☁️ 云服务商部署

### AWS EKS (Elastic Kubernetes Service)

**特点：**
- 托管服务，无需管理控制平面
- 与 AWS 服务深度集成
- 自动扩缩容和更新

**部署步骤：**
```bash
# 安装 eksctl
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# 创建集群
eksctl create cluster --name my-cluster --region us-west-2 --nodegroup-name workers --node-type t3.medium --nodes 3 --nodes-min 1 --nodes-max 4
```

### Azure AKS (Azure Kubernetes Service)

**特点：**
- 完全托管的 K8s 服务
- 与 Azure 服务集成
- 支持 Windows 容器

**部署步骤：**
```bash
# 创建资源组
az group create --name myResourceGroup --location eastus

# 创建 AKS 集群
az aks create --resource-group myResourceGroup --name myAKSCluster --node-count 3 --enable-addons monitoring --generate-ssh-keys

# 获取凭据
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster
```

### GCP GKE (Google Kubernetes Engine)

**特点：**
- Google 原生 K8s 服务
- 自动扩缩容和升级
- 与 Google Cloud 服务集成

**部署步骤：**
```bash
# 创建集群
gcloud container clusters create my-cluster --zone us-central1-a --num-nodes 3

# 获取凭据
gcloud container clusters get-credentials my-cluster --zone us-central1-a
```

## 🏗️ 自建集群部署

### kubeadm 部署

**适用场景：**
- 生产环境
- 需要完全控制
- 学习 K8s 内部机制

**前置要求：**
- 至少 2GB RAM
- 2 个 CPU 核心
- 网络连接
- 禁用 swap

**部署步骤：**

1. **安装 Docker**
```bash
# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

2. **安装 kubeadm、kubelet、kubectl**
```bash
# 添加 Kubernetes 仓库
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# 安装组件
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

3. **初始化 Master 节点**
```bash
# 初始化集群
sudo kubeadm init --pod-network-cidr=10.244.0.0/16

# 配置 kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

4. **安装网络插件**
```bash
# 安装 Flannel
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

5. **添加 Worker 节点**
```bash
# 在 Master 节点获取 join 命令
kubeadm token create --print-join-command

# 在 Worker 节点执行 join 命令
sudo kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>
```

### Rancher 部署

**特点：**
- 图形化管理界面
- 多集群管理
- 应用商店

**部署步骤：**
```bash
# 使用 Docker 运行 Rancher
docker run -d --restart=unless-stopped \
  -p 80:80 -p 443:443 \
  --privileged \
  rancher/rancher:latest
```

## 🔧 安装后配置

### 1. 配置 kubectl
```bash
# 设置别名
echo 'alias k=kubectl' >> ~/.bashrc
source ~/.bashrc

# 启用自动补全
echo 'source <(kubectl completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### 2. 安装常用工具
```bash
# 安装 Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# 安装 kubectx 和 kubens
sudo git clone https://github.com/ahmetb/kubectx /opt/kubectx
sudo ln -s /opt/kubectx/kubectx /usr/local/bin/kubectx
sudo ln -s /opt/kubectx/kubens /usr/local/bin/kubens
```

### 3. 验证安装
```bash
# 检查集群状态
kubectl cluster-info

# 查看节点
kubectl get nodes

# 查看系统 Pod
kubectl get pods --all-namespaces
```

## 🎯 选择建议

### 学习阶段
- **推荐**: Minikube 或 Docker Desktop
- **原因**: 简单易用，资源占用少

### 开发测试
- **推荐**: Kind 或本地 kubeadm
- **原因**: 更接近生产环境

### 生产环境
- **推荐**: 云服务商托管服务
- **原因**: 高可用、自动维护、成本效益

## 🛠️ 实践练习

### 练习 1：Minikube 环境搭建
1. 安装 Minikube
2. 启动集群
3. 部署示例应用
4. 访问应用

### 练习 2：多节点集群
1. 使用 kubeadm 创建集群
2. 添加 Worker 节点
3. 部署应用并测试

### 练习 3：云环境部署
1. 在云服务商创建集群
2. 配置 kubectl
3. 部署应用

## 📚 扩展阅读

- [Kubernetes 官方安装指南](https://kubernetes.io/docs/setup/)
- [Minikube 官方文档](https://minikube.sigs.k8s.io/)
- [kubeadm 官方文档](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)

## 🎯 下一步

完成安装后，继续学习：
- [Pod详解](./04-pod/README.md)
- [ReplicaSet与Deployment](./05-deployment/README.md) 