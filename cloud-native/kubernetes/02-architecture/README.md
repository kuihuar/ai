# Kubernetes 架构详解

## 🏗️ 整体架构

Kubernetes 采用主从架构（Master-Worker），由控制平面（Control Plane）和数据平面（Data Plane）组成。

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
├─────────────────────────────────────────────────────────────┤
│                    Control Plane (Master)                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   API       │ │  Scheduler  │ │ Controller  │            │
│  │  Server     │ │             │ │   Manager   │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
│  ┌─────────────┐ ┌─────────────┐                            │
│  │   etcd      │ │  Cloud      │                            │
│  │             │ │ Controller  │                            │
│  └─────────────┘ └─────────────┘                            │
├─────────────────────────────────────────────────────────────┤
│                    Data Plane (Worker)                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   Worker    │ │   Worker    │ │   Worker    │            │
│  │   Node 1    │ │   Node 2    │ │   Node N    │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

## 🎛️ 控制平面组件

### 1. API Server
API Server 是 Kubernetes 集群的统一入口，所有组件都通过 API Server 进行通信。

**功能：**
- 提供 RESTful API 接口
- 验证和授权请求
- 数据持久化到 etcd
- 提供集群状态查询

**特点：**
- 无状态设计，可水平扩展
- 支持多种认证方式
- 提供审计日志

### 2. etcd
etcd 是一个分布式键值存储系统，存储 Kubernetes 集群的所有数据。

**存储内容：**
- 集群配置信息
- 资源对象状态
- 集群元数据

**特点：**
- 高可用性（通常部署3个或5个节点）
- 强一致性
- 支持事务操作

### 3. Scheduler
Scheduler 负责将 Pod 调度到合适的 Node 上运行。

**调度过程：**
1. **过滤阶段（Predicates）**：过滤不满足条件的节点
2. **评分阶段（Priorities）**：对满足条件的节点进行评分
3. **选择阶段**：选择得分最高的节点

**调度策略：**
- 资源需求匹配
- 节点亲和性/反亲和性
- Pod 亲和性/反亲和性
- 污点和容忍

### 4. Controller Manager
Controller Manager 运行各种控制器，确保集群状态符合期望。

**主要控制器：**
- **Node Controller**：监控节点状态
- **Replication Controller**：确保 Pod 副本数量
- **Endpoints Controller**：维护 Service 端点
- **Service Account & Token Controller**：管理服务账户
- **Namespace Controller**：管理命名空间生命周期

### 5. Cloud Controller Manager
Cloud Controller Manager 与云服务商集成，管理云资源。

**功能：**
- 节点管理（创建、删除）
- 路由管理
- 负载均衡器管理
- 存储卷管理

## 🔧 数据平面组件

### 1. kubelet
kubelet 是每个节点上的主要代理，管理该节点上的容器。

**职责：**
- 管理 Pod 生命周期
- 监控容器健康状态
- 执行容器探针
- 挂载存储卷
- 下载容器镜像

**工作流程：**
1. 从 API Server 获取 Pod 清单
2. 创建和管理容器
3. 定期向 API Server 报告状态

### 2. kube-proxy
kube-proxy 是网络代理，实现 Service 抽象。

**功能：**
- 维护网络规则
- 实现负载均衡
- 支持多种代理模式

**代理模式：**
- **userspace**：用户空间代理
- **iptables**：基于 iptables 的代理（默认）
- **ipvs**：基于 IPVS 的代理

### 3. Container Runtime
Container Runtime 负责运行容器，如 Docker、containerd、CRI-O。

**功能：**
- 拉取容器镜像
- 启动和停止容器
- 管理容器资源

## 🌐 网络架构

### 1. Pod 网络模型
每个 Pod 都有唯一的 IP 地址，Pod 间可以直接通信。

**网络特点：**
- Pod 网络是扁平的
- Pod IP 在集群内唯一
- 支持多种网络插件

### 2. Service 网络
Service 为 Pod 提供稳定的网络端点。

**Service 类型：**
- **ClusterIP**：集群内部访问（默认）
- **NodePort**：通过节点端口访问
- **LoadBalancer**：外部负载均衡器
- **ExternalName**：外部服务别名

### 3. 网络插件
常见的网络插件包括：
- **Flannel**：简单易用
- **Calico**：企业级功能丰富
- **Weave Net**：自动发现
- **Cilium**：基于 eBPF

## 🔐 安全架构

### 1. 认证（Authentication）
支持多种认证方式：
- **证书认证**：基于 TLS 证书
- **Token 认证**：基于 Bearer Token
- **Basic 认证**：用户名密码
- **OpenID Connect**：OAuth2 集成

### 2. 授权（Authorization）
基于 RBAC（Role-Based Access Control）：
- **Role**：命名空间级别的权限
- **ClusterRole**：集群级别的权限
- **RoleBinding**：绑定角色到用户
- **ClusterRoleBinding**：绑定集群角色

### 3. 准入控制（Admission Control）
在请求持久化前进行验证和修改：
- **ValidatingAdmissionWebhook**：验证请求
- **MutatingAdmissionWebhook**：修改请求

## 📊 高可用架构

### 1. 控制平面高可用
- **多 Master 节点**：通常部署3个或5个
- **负载均衡器**：分发 API 请求
- **etcd 集群**：数据高可用

### 2. 工作节点高可用
- **多 Worker 节点**：避免单点故障
- **Pod 反亲和性**：分散部署
- **节点故障转移**：自动重新调度

## 🛠️ 实践练习

### 练习 1：查看集群组件
```bash
# 查看所有命名空间的 Pod
kubectl get pods --all-namespaces

# 查看系统组件
kubectl get pods -n kube-system

# 查看节点信息
kubectl describe nodes
```

### 练习 2：理解调度过程
1. 创建带有资源限制的 Pod
2. 观察调度过程
3. 分析调度决策

### 练习 3：网络连通性测试
1. 创建多个 Pod
2. 测试 Pod 间通信
3. 创建 Service 并测试访问

## 📚 扩展阅读

- [Kubernetes 架构设计](https://kubernetes.io/docs/concepts/architecture/)
- [etcd 官方文档](https://etcd.io/docs/)
- [Kubernetes 网络模型](https://kubernetes.io/docs/concepts/services-networking/)

## 🎯 下一步

理解架构后，继续学习：
- [K8s安装部署](./03-installation/README.md)
- [Pod详解](./04-pod/README.md) 