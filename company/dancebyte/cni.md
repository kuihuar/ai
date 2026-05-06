# Kubernetes CNI 核心知识点

## 1. CNI 是什么，解决什么问题

- CNI（Container Network Interface）是容器网络插件标准。
- Kubernetes 通过 CNI 为 Pod 提供网络连接能力，包括分配 IP、配置路由、连通集群网络。
- 目标是网络能力标准化：K8s 定义调用方式，网络厂商实现插件。

一句话：**CNI 负责让 Pod “拿到网络身份并能通信”。**

## 2. CNI 在 K8s 中的位置

- **kubelet** 在 Pod sandbox 创建阶段调用 CNI。
- **CRI 运行时**（containerd / CRI-O）负责触发网络配置流程。
- **CNI 插件**完成具体动作：
  - 创建一对虚拟网卡 veth pair（内核虚拟以太网卡对）
  - 将 veth pair 一端移入 Pod 网络命名空间，留在 Pod 里作为 eth0
  - 将 veth pair 另一端接入宿主机网络命名空间，作为 eth1
  - 为 Pod 分配唯一 Pod IP 地址、子网、网关（通过 IPAM）
  - 在 Pod 内部网卡 eth0 配置 IP、子网掩码、网关
  - 配置 Pod 内部静态路由：默认路由指向宿主机网关，实现跨节点 / 外网访问
  - 宿主机侧配置网桥 / OVS 端口、二层转发规则
  - 配置路由、转发规则、策略等
  - 配置 iptables /nftables 规则：
    - 容器南北向流量 SNAT/DNAT
    - 集群内 Service 负载均衡转发
    - 网络策略流量放行 / 拒绝规则
  - 加载并应用 NetworkPolicy 网络策略，做 Pod 间访问隔离与准入控制
  - 配置内核网络参数：ip 转发开启、ARP 转发、桥接参数等
  - 做链路就绪检测，保证 Pod 网络通后才就绪


## 3. 核心概念

- **Pod 网络模型**：每个 Pod 有独立 IP；同 Pod 内容器共享网络命名空间。
- **Node-to-Node 通信**：不同节点 Pod 之间可达（依赖 Overlay 或 Underlay）。
- **IPAM**：IP 地址管理，负责分配/回收 Pod IP。
- **NetworkPolicy**：基于标签的 L3/L4 访问控制（由支持策略的 CNI 实现）。

## 4. Pod 网络建立流程（简版）

1. Pod 被调度到某节点。
2. kubelet 让运行时创建 sandbox。
3. 运行时调用 CNI `ADD`。
4. CNI 插件创建 veth、分配 IP、写路由和规则。
5. Pod 启动后可进行网络通信。
6. Pod 删除时调用 CNI `DEL`，回收网络资源。

## 5. 常见网络实现模式

- **Overlay**（如 VXLAN）  
  跨节点通过隧道封装，部署简单，性能有额外开销。
- **BGP / 路由直连**（如 Calico BGP 模式）  
  网络路径更直接，性能通常更好，但路由规划复杂度更高。
- **主机本地桥接 + 路由**  
  依赖底层网络能力和路由可达性设计。

## 6. CNI 常见能力点

- Pod IP 自动分配与回收
- 跨节点连通
- NetworkPolicy（白名单/黑名单控制）
- Egress/NAT（集群出网）
- 多网卡多网络（通过 Multus 等多 CNI 机制）

## 7. Calico / Flannel / Canal / Multus 快速定位

- **Calico**：网络 + 网络策略能力强，支持 BGP 和 Overlay 多模式。
- **Flannel**：专注基础网络连通，架构简单、上手快，不以策略能力见长。
- **Canal**：Flannel + Calico 组合方案（Flannel 做连通，Calico 提供策略）。
- **Multus**：元插件，可给 Pod 挂多个网络接口（默认网 + SR-IOV/Macvlan 等附加网）。

## 8. 高频故障排查思路

- **Pod 间不通**
  - 看 Pod IP、路由、节点转发、隧道状态是否正常。
- **DNS 解析失败**
  - 先看 CoreDNS，再看 Pod 到 DNS Service 的链路与 NetworkPolicy。
- **策略误杀**
  - 检查命名空间/Pod 标签与 ingress/egress 规则是否匹配。
- **跨节点偶发丢包**
  - 检查 MTU、封装模式、底层网络质量和 conntrack 压力。

## 9. 常用排障命令（面试可口述）

- `kubectl get pod -A -o wide`：确认 Pod IP 与节点分布
- `kubectl get netpol -A`：查看策略是否限制流量
- `kubectl describe pod <pod>`：查看事件与 CNI 报错
- `kubectl logs -n kube-system <cni-pod>`：查看插件日志
- 节点侧：`ip a`、`ip route`、`iptables -S`、`ss -ant`

## 10. 面试表达模板（30 秒）

“CNI 是 Kubernetes 的容器网络标准接口。Pod 创建时 kubelet/CRI 调用 CNI 插件完成 veth 创建、IPAM 分配和路由配置，从而实现 Pod 间可达。选型上 Flannel 更轻量，Calico 策略更完整，Canal 是两者组合，Multus 用于多网卡场景。排障一般按 Pod IP 与路由、策略规则、CNI 日志、底层网络质量逐层定位。”

## 11. 延伸阅读

- 面试题与专项题见：`cni-interview.md`


|模式|实现技术|性能|配置难度|适用场景|
|--|--|--|--|--|
|VXLAN (Overlay)|UDP 封装|中,低（即插即用） |跨网段、对物理网无控制权、混合云|
|BGP (L3)|路由广播|高|高（需配置网络设备） |自建机房、对延迟极度敏感的应用|
|Host-GW|静态路由|高|中|节点都在同一个二层交换机下|
|VPC 模式|云厂商 API|极高|极低（云平台自动处理） |纯公有云环境|
