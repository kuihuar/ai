# Kubernetes CSI 面试题（含参考答案）

## 基础题

### 1）什么是 CSI？为什么 Kubernetes 要引入 CSI？
**答：**
- CSI 是容器编排系统与存储系统的标准接口。
- 目的：存储插件解耦、独立演进、跨平台复用，避免 in-tree 插件和 K8s 强绑定。

### 2）PVC、PV、StorageClass 分别是什么？
**答：**
- `PVC`：应用对存储的申请。
- `PV`：集群中可用的实际存储卷对象。
- `StorageClass`：定义如何供应卷（驱动、参数、回收策略、绑定策略）。

### 3）动态供给流程是怎样的？
**答：**
1. 创建 PVC。
2. provisioner 监听到 PVC，调用 `CreateVolume`。
3. 后端创建卷，K8s 生成 PV 并绑定 PVC。
4. Pod 使用 PVC，后续进入 attach/mount。

### 4）`WaitForFirstConsumer` 的价值是什么？
**答：**
- 延迟卷创建到 Pod 确定调度节点后，避免提前在错误可用区创建卷。
- 对多可用区集群尤其关键。

### 5）`Delete` 和 `Retain` 的区别？
**答：**
- `Delete`：PVC 删除时联动删除底层卷。
- `Retain`：保留底层卷，适合关键数据防误删，需要人工回收。

## 进阶题

### 6）CSI Controller 和 Node 侧分别负责什么？
**答：**
- Controller 侧负责控制平面操作：建卷、删卷、attach/detach 等。
- Node 侧负责节点挂载生命周期：`NodeStageVolume` / `NodePublishVolume`。

### 7）为什么 Pod 会卡在 `ContainerCreating`？
**答：**
- 常见是卷未成功挂载：设备不存在、权限错误、文件系统异常、Node 插件异常。
- 需要看 Pod 事件、kubelet 日志和 CSI Node 插件日志。

### 8）PVC 一直 Pending，怎么排查？
**答：**
1. 看 PVC 事件（是否找不到 StorageClass/provisioner）。
2. 检查 StorageClass 的 `provisioner` 字段。
3. 检查 provisioner sidecar 日志（`CreateVolume` 报错）。
4. 检查后端存储配额、认证、网络连通性。

### 9）如何理解 attach 与 mount 的区别？
**答：**
- attach：把卷“连接”到节点（云盘映射到主机）。
- mount：把节点上的设备挂载到 Pod 可见目录。
- 某些存储类型可跳过 attach，直接 mount（取决于驱动能力）。

### 10）CSI 如何支持卷扩容？
**答：**
- 前提是 `StorageClass` 开启 `allowVolumeExpansion` 且驱动支持。
- 扩容分为控制面扩容与节点文件系统扩容两部分，二者都成功才算完成。

## 场景题

### 11）线上出现跨可用区挂载失败，你如何处理？
**答题思路：**
- 先确认报错是否拓扑不匹配（zone/region 标签）。
- 校验 StorageClass 是否使用 `WaitForFirstConsumer`。
- 检查 Pod 的 nodeSelector / affinity 是否把工作负载调度到错误区域。
- 回看 CSI 驱动拓扑能力声明和节点标签。

### 12）如何设计“高可靠 + 可恢复”的有状态服务存储方案？
**答题思路：**
- 数据卷：高可用存储 + 合理 `reclaimPolicy`（关键业务建议 `Retain`）。
- 备份：快照策略（定时 + 异地/跨区复制）。
- 恢复：演练从 `VolumeSnapshot` 或备份恢复 PVC。
- 发布：灰度升级并验证数据一致性与恢复时间目标（RTO/RPO）。

## 高频追问（可直接背）

### 追问 1：in-tree 到 CSI 的迁移风险有哪些？
- 驱动兼容性与版本矩阵。
- 旧 PV 元数据与新驱动识别关系。
- 迁移窗口中的业务中断风险与回滚预案。

### 追问 2：什么时候不建议用共享文件存储？
- 高并发随机写、低延迟强依赖场景。
- 元数据操作密集导致性能抖动场景。
- 需要严格块设备语义的数据库场景。

### 追问 3：如何给面试官展示排障能力？
- 先讲“对象状态”再讲“组件日志”再讲“后端存储”。
- 给出可执行命令路径（events -> sidecar -> node plugin）。
- 明确根因、修复动作、预防措施三段式结论。

## Longhorn / Hwameistor 专项题

### 13）Longhorn 和 Hwameistor 的定位分别是什么？
**答：**
- **Longhorn**：云原生分布式块存储，主打易用性、可视化运维和副本高可用，适合通用有状态业务。
- **Hwameistor**：以本地盘为核心的容器存储方案，强调本地卷性能与 Kubernetes 深度集成，适合对性能和本地数据路径敏感场景。

### 14）Longhorn 核心架构怎么讲？
**答：**
- 每个卷通常有多个副本（Replica）分布在不同节点。
- 通过 Engine 负责 I/O 聚合和副本一致性维护。
- Manager 负责调度、健康状态、重建、副本迁移等控制逻辑。
- 具备快照、备份到对象存储、卷恢复等能力。

### 15）Hwameistor 的核心能力与组件思路是什么？
**答：**
- 以节点本地磁盘池化（HDD/SSD/NVMe）为基础。
- 通过 CSI 驱动提供动态供给、挂载、扩容等能力。
- 提供本地卷高可用能力（依赖多副本或数据同步机制，具体取决于部署模式）。
- 结合调度策略让 Pod 尽量与卷在拓扑上匹配，降低跨节点数据访问开销。

### 16）Longhorn 为什么常被认为“上手快”？
**答：**
- 安装和接入相对简单，文档和可视化控制面较完善。
- 备份/恢复、快照、卷健康状态可视化能力成熟。
- 对中小规模集群和通用场景成本较低、运维门槛更友好。

### 17）Hwameistor 适合哪些工作负载？
**答：**
- 对本地 I/O 性能敏感的数据库、中间件、日志类工作负载。
- 希望充分利用本地 NVMe/SSD 且可接受本地盘运维复杂度的场景。
- 对“数据本地性 + 调度协同”有明确要求的场景。

### 18）Longhorn 常见故障排查路径？
**答题思路：**
1. 先看卷状态（Degraded/Faulted）与副本健康。
2. 检查副本重建是否受节点资源、网络抖动影响。
3. 查看 Manager / Engine / Instance Manager 日志。
4. 检查底层磁盘空间、inode、IO 延迟和网络连通性。

### 19）Hwameistor 常见故障排查路径？
**答题思路：**
1. 看本地磁盘池状态与可分配容量是否充足。
2. 看 PVC/PV 绑定事件与 CSI provision/mount 日志。
3. 检查卷副本或同步状态（若启用高可用）。
4. 检查节点标签、亲和策略与 Pod 调度是否匹配。

### 20）面试中怎么做 Longhorn vs Hwameistor 选型表达？
**答：**
- 若目标是“快速落地 + 易运维 + 完整可视化能力”，优先 Longhorn。
- 若目标是“极致本地盘性能 + 数据本地性”，倾向 Hwameistor。
- 最终要结合：性能目标、可用性要求、团队运维能力、备份恢复链路、故障演练结果。

## 速记版（1 分钟）

- CSI = K8s 存储标准接口，解耦存储插件。
- 抽象主线：`StorageClass -> PVC -> PV -> Pod`。
- 关键流程：`CreateVolume -> (Attach) -> NodePublishVolume`。
- 关键策略：`WaitForFirstConsumer`、`Retain/Delete`、`allowVolumeExpansion`。
- 常见故障：PVC Pending、挂载失败、拓扑不匹配；排查按“对象 -> 日志 -> 后端”逐层推进。




## CSI 存储 面试题 + 答案
1. 什么是 CSI？作用？
答：CSI 容器存储接口，是 K8s 标准化存储插件接口；把存储能力从 K8s 源码解耦，第三方存储只需实现 CSI 即可对接，支持动态供给、挂载、扩容、快照。
2. CSI 架构组成
答：CSI Controller：控制面，创建 / 删除 / 扩容 / 快照卷；CSI Node：节点面，挂载、卸载卷；官方 Sidecar：provisioner、attacher、resizer、snapshotter 等控制器。
3. PV、PVC、StorageClass 关系
答：PVC 用户申请存储；PV 实际后端存储卷；StorageClass 存储模板，实现动态供给，不用手动建 PV。
4. 访问模式 RWO、ROX、RWOP
答：RWO：只能单节点读写（块存储常用）ROX：多节点只读RWOP：多节点读写
5. 块存储和文件存储区别
答：块存储：裸磁盘设备，/dev/sda，单节点独占，适合数据库，volumeMode:Block；文件存储：共享目录，多节点可同时挂载，适合日志配置，volumeMode:Filesystem。
6. StorageClass 立即绑定 和 延迟绑定
答：立即绑定：PVC 创建就建 PV，适合普通云盘；延迟绑定：等 Pod 调度到节点再建 PV，适合本地盘、有拓扑约束存储。


## 本地开源存储 
Longhorn 特点
答：Rancher 开源，基于本地盘做分布式块存储；支持副本、快照、备份、在线扩容，部署简单、有 UI，运维友好。




## 核心组件
- External Components (K8s 侧托管)：

    - External Provisioner：监听 PVC，负责在云端创建/删除实际的存储卷。

    - External Attacher：负责将存储卷挂载/卸载到对应的Node 节点上（逻辑挂载）。

    - External Resizer：处理卷扩容请求。

- CSI Driver (厂商实现)：

    - Identity Service：身份验证，告知 K8s 插件信息。

    - Controller Service：负责创建、删除、挂载、卸载卷（对应云端操作）。

    - Node Service：运行在每个 Node 上，负责最后一步：将块设备格式化并挂载到 Pod 的目录里。