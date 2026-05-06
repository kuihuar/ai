
### 控制面 (Control Plane) —— 集群的大脑
1. kube-apiserver	整个集群的网关，API 入口，权限控制
2. etcd		键值存储，存储集群状态
3. controller-manager		控制器，负责资源的创建、删除、更新等操作
4. scheduler-scheduler		调度器，负责将 Pod 分配到节点上
5. cloud-controller-manager		云控制器，负责与云平台交互，管理云资源
6. CSI Controller

### 工作节点 (Worker Node/ Node / Data Plan) —— 集群的执行单元

1. kubelet
2. kube-proxy
3. (CRI)Container Runtime
4. Add-on Components
  - CNI(Clico)
  - CoreDNS
  - Metrics Server
  - Dashboard
5. CSI Node  


CSI Driver = CSI Controller + CSI Node + (一些辅助组件)



1. kubectl 提交 yaml 到 apiserver
2. apiserver 认证、授权、准入控制
3. 写入 etcd 存储
4. 调度器为 Pod 选节点，绑定 nodeName
5. 目标节点 kubelet 感知到 Pod
6. 调用 CRI 拉镜像、创建启动容器
7. CNI 配网络、CSI 挂载存储
8. 上报 Pod 状态到 apiserver
