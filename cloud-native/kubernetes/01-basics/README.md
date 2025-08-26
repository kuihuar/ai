# Kubernetes 基础概念

## 📖 什么是 Kubernetes？

Kubernetes (K8s) 是一个开源的容器编排平台，用于自动化部署、扩展和管理容器化应用程序。它提供了一个可移植、可扩展的开源平台，用于管理容器化的工作负载和服务。

## 🎯 核心概念

### 1. 容器 (Container)
容器是轻量级、可移植的软件包，包含运行应用程序所需的所有依赖项。

**特点：**
- 轻量级：比虚拟机更小、更快
- 可移植：在任何支持Docker的环境中运行
- 隔离性：容器间相互隔离
- 一致性：开发、测试、生产环境一致

### 2. Pod
Pod 是 Kubernetes 中最小的可部署单元，包含一个或多个容器。

**Pod 特点：**
- 一个 Pod 可以包含多个容器
- Pod 内的容器共享网络命名空间
- Pod 内的容器可以通过 localhost 通信
- Pod 是临时的，可以被创建、删除、替换

**示例 Pod 定义：**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
```

### 3. Node (节点)
Node 是 Kubernetes 集群中的工作机器，可以是物理机或虚拟机。

**Node 类型：**
- **Master Node**: 控制平面节点，运行集群控制组件
- **Worker Node**: 工作节点，运行应用程序容器

**Node 组件：**
- **kubelet**: 节点代理，管理容器生命周期
- **kube-proxy**: 网络代理，实现服务抽象
- **Container Runtime**: 容器运行时（如 Docker、containerd）

### 4. Namespace (命名空间)
Namespace 提供了一种在单个集群内隔离资源组的机制。

**默认命名空间：**
- `default`: 默认命名空间
- `kube-system`: 系统组件
- `kube-public`: 公开访问的资源
- `kube-node-lease`: 节点租约信息

**创建命名空间：**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
```

## 🔧 基本操作

### 1. 查看集群信息
```bash
# 查看集群信息
kubectl cluster-info

# 查看节点
kubectl get nodes

# 查看命名空间
kubectl get namespaces
```

### 2. 创建和管理 Pod
```bash
# 创建 Pod
kubectl apply -f pod.yaml

# 查看 Pod
kubectl get pods

# 查看 Pod 详细信息
kubectl describe pod <pod-name>

# 删除 Pod
kubectl delete pod <pod-name>
```

### 3. 进入 Pod 容器
```bash
# 进入容器执行命令
kubectl exec -it <pod-name> -- /bin/bash

# 查看容器日志
kubectl logs <pod-name>
```

## 📊 资源对象层次结构

```
Cluster
├── Namespace
│   ├── Pod
│   │   ├── Container
│   │   └── Container
│   ├── Service
│   ├── ConfigMap
│   └── Secret
└── Node
    ├── Pod
    └── Pod
```

## 🎯 学习要点

### 1. 理解容器化
- 容器 vs 虚拟机
- Docker 基础操作
- 容器镜像构建

### 2. 掌握 Pod 概念
- Pod 生命周期
- Pod 网络模型
- Pod 资源限制

### 3. 熟悉基本命令
- kubectl 常用命令
- YAML 配置文件格式
- 资源创建和管理

## 🛠️ 实践练习

### 练习 1：创建第一个 Pod
1. 创建一个简单的 nginx Pod
2. 验证 Pod 运行状态
3. 访问 Pod 中的服务

### 练习 2：多容器 Pod
1. 创建一个包含多个容器的 Pod
2. 观察容器间通信
3. 理解 Pod 网络模型

### 练习 3：命名空间管理
1. 创建自定义命名空间
2. 在命名空间中部署应用
3. 理解资源隔离

## 📚 扩展阅读

- [Kubernetes 官方概念文档](https://kubernetes.io/docs/concepts/)
- [Docker 容器基础](https://docs.docker.com/get-started/)
- [YAML 语法指南](https://yaml.org/spec/)

## 🎯 下一步

掌握基础概念后，继续学习：
- [K8s架构](./02-architecture/README.md)
- [K8s安装部署](./03-installation/README.md) 