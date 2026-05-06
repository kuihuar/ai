# Kubernetes CSI 核心知识点

## 1. CSI 是什么，解决了什么问题

- CSI（Container Storage Interface）是容器编排系统（如 Kubernetes）与存储系统之间的标准接口。
- 目标是把存储能力从 Kubernetes 核心代码中解耦，存储厂商只需实现 CSI 驱动即可接入。
- 相比早期 in-tree 插件，CSI 具备更快迭代、独立发布、跨平台复用的优势。

一句话：**K8s 定义规范，存储厂商实现驱动，集群按统一方式使用存储。**

## 2. CSI 架构与组件

Kubernetes 中通常会看到 3 类组件：

- **CSI Driver（厂商实现）**
  - 实现 CSI 标准 RPC，负责创建卷、删除卷、挂载卷等实际能力。
  - 常见分为 Controller 插件和 Node 插件。
- **CSI Sidecar（社区通用组件）**
  - 如 `external-provisioner`、`external-attacher`、`external-resizer`、`external-snapshotter`。
  - 负责监听 K8s 对象变化并调用对应 CSI RPC。
- **Kubelet + API Server（K8s 自身）**
  - Kubelet 负责节点侧卷挂载生命周期。
  - API Server 持久化 PVC/PV/VolumeAttachment 等对象状态。

## 3. 关键对象与关系

- `StorageClass`：定义存储“类型”和参数（provisioner、reclaimPolicy、bindingMode、参数等）。
- `PVC`（PersistentVolumeClaim）：应用发起的存储申请（容量、访问模式、存储类）。
- `PV`（PersistentVolume）：实际被供应出来的存储卷资源。
- `Pod`：通过 `volumes.persistentVolumeClaim` 使用 PVC。

关系可简化为：

`Pod -> PVC -> PV -> 后端存储系统`

## 4. 从声明到挂载的完整流程（动态供给）

1. 用户创建 PVC（指定 StorageClass）。
2. `external-provisioner` 监听到 PVC，调用 CSI `CreateVolume`。
3. 后端存储创建卷，K8s 生成并绑定 PV 与 PVC。
4. Pod 调度到某节点后，若需要附加卷：
   - `external-attacher` 调用 `ControllerPublishVolume`（Attach）。
5. 目标节点 kubelet 调用 Node 侧能力：
   - `NodeStageVolume`（可选，设备级准备）
   - `NodePublishVolume`（把卷挂载到 Pod 目录）
6. Pod 启动后可读写卷数据。

## 5. 常见能力点（面试高频）

- **动态供给（Dynamic Provisioning）**：PVC 自动触发创建 PV。
- **静态供给（Static Provisioning）**：管理员预先创建 PV，再由 PVC 绑定。
- **卷扩容（Expansion）**：修改 PVC 容量，驱动与文件系统共同完成扩容。
- **快照（Snapshot）**：基于 `VolumeSnapshot` 实现数据时间点保护。
- **克隆（Clone）**：从已有 PVC 快速生成新卷（同存储类常见）。
- **拓扑感知（Topology-aware）**：按可用区/节点拓扑创建与调度存储。

## 6. 重要配置项与设计取舍

- `volumeBindingMode`：
  - `Immediate`：PVC 创建时立刻供给。
  - `WaitForFirstConsumer`：等 Pod 调度时再决策，减少跨可用区错误绑定。
- `reclaimPolicy`：
  - `Delete`：PVC 删除后，底层卷也删除。
  - `Retain`：保留底层数据，适合高价值数据防误删。
- `allowVolumeExpansion`：
  - 是否允许在线/离线扩容（取决于驱动与文件系统支持）。

## 7. 常见故障与排查思路

- PVC 一直 `Pending`
  - 看 `StorageClass` 是否存在，provisioner 名称是否匹配。
  - 看 CSI Controller sidecar 日志是否有 `CreateVolume` 错误。
- Pod 卡在 `ContainerCreating`
  - 看节点 kubelet 事件与 Node 插件日志（挂载、权限、设备识别）。
- 多可用区挂载失败
  - 检查 `volumeBindingMode` 是否为 `WaitForFirstConsumer`，并确认拓扑标签配置。
- 卷删除异常或残留
  - 检查 `reclaimPolicy` 与后端存储回收策略是否一致。

## 8. 与 CNI / CRI 的关系

- `CRI`：管理容器运行时（创建/启动容器）。
- `CNI`：为 Pod 提供网络连接。
- `CSI`：为 Pod 提供持久化存储。

三者分别解决“算力运行、网络连通、数据持久化”三类核心问题。

## 9. 面试表达模板（30 秒）

“CSI 是 K8s 的存储标准接口，核心价值是把存储插件从 K8s 主仓解耦。业务侧主要关注 StorageClass、PVC、PV 三层抽象。PVC 触发 dynamic provisioning 后，provisioner 调用 CSI CreateVolume 创建后端卷，Pod 调度后再由 attacher 和 kubelet 完成 attach/mount。线上问题通常集中在 PVC Pending、挂载失败和拓扑不匹配，排查时要按 K8s 对象事件到 CSI sidecar 日志再到存储后端逐层定位。”

## 10. 延伸阅读

- 面试题与追问见：`csi-interview.md`
