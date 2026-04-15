# 虚拟机操作全集（KubeVirt 面试版）

本文整理“虚拟机除了创建之外的常见操作”，覆盖面试高频动作，并给出每个操作的 1-N 步骤、关键命令、失败信号与排障方向。

适用范围：Kubernetes + KubeVirt 场景。

---

## 1. 操作总览（按生命周期分组）

## 1.1 运行控制类

- 启动（Start）
- 停止（Stop）
- 重启（Restart）
- 强制停止（Force Stop）
- 暂停（Pause）
- 恢复（Unpause）

## 1.2 访问连接类

- 控制台登录（VNC / Serial Console）
- SSH 连接

## 1.3 数据保护类

- 克隆（Clone）
- 快照（Snapshot）
- 恢复（Restore）

## 1.4 可用性类

- 热迁移（Live Migration）

## 1.5 配置变更类

- CPU/内存规格调整
- 磁盘扩容（系统盘/数据盘）
- 网络调整（附加网卡/NAD 变更）

## 1.6 回收类

- 删除与清理（Delete + Finalizer 回收）

---

## 2. 运行控制类

## 2.1 启动（Start）

**用途**：将已停止 VM 启动为运行态。

**关键对象/组件**：`VirtualMachine`、`VirtualMachineInstance`、`virt-controller`、`virt-launcher`。

**1-N 步骤**
1. 确认 VM 存在且处于停止态。
2. 发起启动请求（声明 `running: true` 或 `virtctl start`）。
3. `virt-controller` Reconcile，生成或恢复 VMI。
4. scheduler 绑定节点，kubelet 拉起 virt-launcher。
5. 网络/存储接入完成，状态收敛到 Running/Ready。

**关键命令**
```bash
virtctl start <vm-name> -n <namespace>
kubectl get vm,vmi -n <namespace>
kubectl get pod -n <namespace> | rg virt-launcher
```

## 2.2 停止（Stop）

**用途**：正常关闭 VM，释放运行资源。

**关键对象/组件**：`VM`、`VMI`、`virt-controller`。

**1-N 步骤**
1. 确认目标 VM 当前 Running。
2. 发起停止请求（`running: false` 或 `virtctl stop`）。
3. 控制器推进 VMI 退出并回收运行单元。
4. 释放计算资源，保留 VM 声明对象。
5. 状态收敛到 Stopped。

**关键命令**
```bash
virtctl stop <vm-name> -n <namespace>
kubectl get vm,vmi -n <namespace>
```

## 2.3 重启（Restart）

**用途**：在不删除 VM 声明的前提下重建运行实例。

**关键对象/组件**：`VM`、`VMI`、`virt-controller`。

**1-N 步骤**
1. 发起重启请求。
2. 控制器终止当前 VMI。
3. 重新创建 VMI 并调度。
4. 重新拉起 virt-launcher。
5. 状态恢复 Running。

**关键命令**
```bash
virtctl restart <vm-name> -n <namespace>
kubectl get vmi -n <namespace> -w
```

## 2.4 强制停止（Force Stop）

**用途**：常规停止卡住时强制终止运行实例。

**关键对象/组件**：`VMI`、`virt-handler`、节点 kubelet。

**1-N 步骤**
1. 判定常规 stop 超时或失效。
2. 执行强制停止动作。
3. 强制终止运行实例并清理挂起状态。
4. 检查卷与网络资源是否已释放。
5. 确认 VM 状态一致。

**关键命令**
```bash
virtctl stop <vm-name> -n <namespace> --force
kubectl describe vmi <vmi-name> -n <namespace>
```

## 2.5 暂停（Pause）与恢复（Unpause）

**用途**：临时冻结/恢复虚机执行状态（不等同停止）。

**关键对象/组件**：`VMI`、`virt-handler`。

**1-N 步骤**
1. 对运行中的 VMI 发起 pause。
2. 计算执行冻结，内存状态保持。
3. 业务验证暂停效果。
4. 发起 unpause 恢复执行。
5. 确认恢复后状态正常。

**关键命令**
```bash
virtctl pause vm <vm-name> -n <namespace>
virtctl unpause vm <vm-name> -n <namespace>
kubectl get vmi -n <namespace>
```

---

## 3. 访问连接类

## 3.1 控制台登录（VNC / Serial）

**用途**：进入虚机控制台进行运维或故障诊断。

**关键对象/组件**：`virt-api`、`virtctl`、`VMI`。

**1-N 步骤**
1. 确认 VMI Running。
2. 选择连接方式（VNC 或 serial）。
3. 建立到 virt-api 的控制台会话。
4. 进入客体系统执行诊断。
5. 退出并记录结果。

**关键命令**
```bash
virtctl console <vm-name> -n <namespace>
virtctl vnc <vm-name> -n <namespace>
```

## 3.2 SSH 连接

**用途**：通过网络远程登录虚机系统。

**关键对象/组件**：VM 网络（Multus/CNI）、Service/LB（可选）、SSH 服务。

**1-N 步骤**
1. 确认虚机网络可达（IP/路由/安全策略）。
2. 确认客体系统 SSH 服务已开启。
3. 准备密钥或密码认证。
4. 建立 SSH 连接。
5. 验证业务与系统状态。

**关键命令**
```bash
ssh <user>@<vm-ip>
kubectl get vmi -n <namespace> -o wide
```

---

## 4. 数据保护类

## 4.1 克隆（Clone）

**用途**：快速复制一台 VM（模板化交付、环境复制）。

**关键对象/组件**：`VirtualMachineClone`（如启用）、源 VM、目标 VM、底层卷。

**1-N 步骤**
1. 选择源 VM 并校验状态（通常建议停机或快照一致性）。
2. 创建 Clone 请求对象。
3. 控制器复制 VM 配置与磁盘数据引用/副本。
4. 生成目标 VM。
5. 校验目标 VM 可启动并可用。

**关键命令**
```bash
kubectl get vm -n <namespace>
kubectl get vmclone -n <namespace>
kubectl describe vmclone <clone-name> -n <namespace>
```

## 4.2 快照（Snapshot）

**用途**：记录某一时点状态，用于回滚与备份。

**关键对象/组件**：`VirtualMachineSnapshot`、存储快照能力（CSI/Longhorn）。

**1-N 步骤**
1. 确认虚机与卷支持快照能力。
2. 创建 VM Snapshot。
3. 控制器协调磁盘快照。
4. 等待 snapshot Ready。
5. 记录恢复点信息。

**关键命令**
```bash
kubectl get vmsnapshot -n <namespace>
kubectl describe vmsnapshot <snapshot-name> -n <namespace>
```

## 4.3 恢复（Restore）

**用途**：从快照恢复虚机状态。

**关键对象/组件**：`VirtualMachineRestore`、目标 VM、底层卷恢复能力。

**1-N 步骤**
1. 选定 snapshot 作为恢复源。
2. 创建 Restore 请求。
3. 控制器执行卷与配置回滚。
4. 恢复完成后重建/拉起 VM。
5. 验证业务数据一致性。

**关键命令**
```bash
kubectl get vmrestore -n <namespace>
kubectl describe vmrestore <restore-name> -n <namespace>
```

---

## 5. 可用性类

## 5.1 热迁移（Live Migration）

**用途**：不停机把 VM 从一个节点迁移到另一个节点。

**关键对象/组件**：`VirtualMachineInstanceMigration`、`virt-controller`、`virt-handler`。

**1-N 步骤**
1. 校验 VMI 支持迁移（网络/存储/配置符合条件）。
2. 创建 Migration 请求。
3. 控制器协调源节点与目标节点迁移。
4. 迁移切换流量与执行上下文。
5. 迁移完成并回收源侧资源。

**关键命令**
```bash
virtctl migrate <vm-name> -n <namespace>
kubectl get vmim -n <namespace>
kubectl describe vmim <migration-name> -n <namespace>
```

---

## 6. 配置变更类

## 6.1 CPU/内存规格调整

**用途**：按负载变更计算规格。

**关键对象/组件**：`VM spec.domain.cpu/memory`、`VMI`、`virt-controller`。

**1-N 步骤**
1. 修改 VM 规格声明（CPU/内存）。
2. 判断是否支持热更新，不支持则重启生效。
3. 控制器推进配置到运行实例。
4. 观察 VMI 与业务指标变化。
5. 回写并确认 Ready。

**关键命令**
```bash
kubectl edit vm <vm-name> -n <namespace>
kubectl get vm,vmi -n <namespace>
```

## 6.2 磁盘扩容（系统盘/数据盘）

**用途**：扩大系统盘或数据盘容量。

**关键对象/组件**：`PVC`、`DataVolume`、`Longhorn`、`CSI`。

**1-N 步骤**
1. 修改 PVC/DataVolume 目标容量。
2. Longhorn 扩容底层卷。
3. CSI 在节点侧完成卷扩展。
4. 客体系统执行文件系统扩容。
5. 验证新容量生效。

**关键命令**
```bash
kubectl get pvc,dv -n <namespace>
kubectl describe pvc <pvc-name> -n <namespace>
```

## 6.3 网络调整（附加网卡/NAD 变更）

**用途**：新增、删除或修改虚机附加网络。

**关键对象/组件**：`VM networks/interfaces`、`NAD`、`Multus`。

**1-N 步骤**
1. 调整 VM 网络声明或 NAD。
2. 控制器检查网络依赖。
3. 按能力决定热更新或重建实例。
4. 多网络重新附加并校验连通。
5. 回写状态与事件。

**关键命令**
```bash
kubectl get network-attachment-definitions -n <namespace>
kubectl describe vm <vm-name> -n <namespace>
kubectl describe pod <virt-launcher-pod> -n <namespace>
```

---

## 7. 回收类

## 7.1 删除与清理（Delete + Finalizer）

**用途**：删除 VM 并按策略回收运行与依赖资源。

**关键对象/组件**：`VM`、`VMI`、Finalizer、PVC/数据卷策略。

**1-N 步骤**
1. 发起删除请求。
2. Finalizer 阶段执行依赖清理逻辑。
3. 删除 VMI 与运行单元。
4. 按保留策略处理卷与快照。
5. Finalizer 移除后对象真正删除。

**关键命令**
```bash
kubectl delete vm <vm-name> -n <namespace>
kubectl get vm,vmi,pvc -n <namespace>
kubectl describe vm <vm-name> -n <namespace>
```

---

## 8. 高频失败信号与排障入口

- **VM 有，VMI 没有**：先查控制器推进与依赖状态
- **VMI Pending**：先查调度（资源/亲和性/污点）
- **virt-launcher 启动失败**：先查 CNI/CSI 与镜像事件
- **盘未挂载**：先查 PVC/PV/DataVolume 与 Longhorn 组件
- **能启动但连不上**：先查 Multus/NAD、策略、CoreDNS
- **删除卡住**：先查 Finalizer 与依赖清理失败

---

## 9. 面试口述模板（30-40 秒）

「创建之外的 VM 操作我一般分六类：运行控制、访问连接、数据保护、可用性、配置变更、回收清理。每个操作都遵循同一逻辑：先提交声明，控制器 Reconcile，依赖就绪后推进运行态，再回写 phase/conditions。比如启动是 VM 到 VMI 再到 virt-launcher；克隆/快照/恢复走数据一致性链路；迁移走 VMI 迁移对象；扩容和网络调整分别由存储和网络组件协同完成；删除则由 Finalizer 做安全清理。排障我按 VM/VMI、virt-launcher、网络、存储、事件分层定位。」 

---

## 10. 开发实现视角（逐操作细化）

这一节回答三个问题：

1. 每个步骤涉及哪些组件、各做什么动作  
2. 每个步骤的知识点/原理是什么  
3. 开发里通常落到哪些代码位置

> 说明：下面的“代码落点”按常见 Operator/KubeVirt 项目结构给出，你可以映射到自己项目中的具体文件名。

## 10.1 启动（Start）

### 步骤 1：提交启动意图
- **组件动作**：客户端（`virtctl`/API）更新 VM 期望状态（如 `spec.running=true`）。
- **知识点/原理**：声明式控制；用户表达“目标状态”，不是直接操作进程。
- **代码落点**：API Handler 或 CLI 调用层；CR 更新逻辑（`client.Update`）。

### 步骤 2：控制器 Reconcile
- **组件动作**：`virt-controller` 监听 VM 变更，判断当前无运行实例，进入创建分支。
- **知识点/原理**：List/Watch + Reconcile 幂等收敛。
- **代码落点**：`Reconcile()` 主流程，`ensureVMI()`/`createVMIIfNeeded()`。

### 步骤 3：VMI 生成与调度
- **组件动作**：创建 VMI，scheduler 绑定节点。
- **知识点/原理**：控制器只负责声明，调度器负责放置。
- **代码落点**：VMI 构造器、调度策略注入（亲和性/污点容忍）。

### 步骤 4：节点落地与状态回写
- **组件动作**：kubelet 拉起 `virt-launcher`，CNI/CSI 接网挂盘，控制器回写状态。
- **知识点/原理**：状态机（`Pending -> Running`）+ conditions 驱动观测。
- **代码落点**：状态聚合函数（`syncStatus()`、`setCondition()`）。

## 10.2 停止（Stop）/重启（Restart）/强制停止（Force Stop）

### 步骤 1：提交动作
- **组件动作**：停止或重启动作下发到 VM。
- **知识点/原理**：停止是期望“无实例”；重启是“先删后建”；强停是异常兜底。
- **代码落点**：动作路由层（stop/restart API）、VM 规格更新逻辑。

### 步骤 2：控制器推进实例终止
- **组件动作**：识别当前 VMI，执行优雅终止或强制终止路径。
- **知识点/原理**：优雅停机优先；强制路径用于超时/卡死状态。
- **代码落点**：`reconcileStop()`、`forceDeleteVMI()`、超时处理。

### 步骤 3：资源回收与一致性校验
- **组件动作**：回收运行单元，校验 VM/VMI 状态一致。
- **知识点/原理**：最终一致性；避免“VM 显示停机但 VMI 残留”。
- **代码落点**：清理函数、状态对账函数。

## 10.3 暂停（Pause）/恢复（Unpause）

### 步骤 1：执行冻结/恢复动作
- **组件动作**：向运行实例发起 pause/unpause。
- **知识点/原理**：暂停保留内存态，不等于停机。
- **代码落点**：VMI subresource 调用层（pause/unpause client）。

### 步骤 2：状态同步
- **组件动作**：控制器更新条件与事件，暴露当前执行态。
- **知识点/原理**：控制面与运行面状态一致性。
- **代码落点**：condition 更新、event recorder。

## 10.4 控制台登录（VNC/Serial）与 SSH

### 步骤 1：会话建立
- **组件动作**：`virtctl` 通过 `virt-api` 建立控制台通道；SSH 通过网络直连 VM。
- **知识点/原理**：控制台是管理通道，SSH 是业务网络通道。
- **代码落点**：console proxy/API route；网络服务配置代码。

### 步骤 2：连接鉴权与通道维护
- **组件动作**：校验权限，建立 websocket/stream 或 SSH 会话。
- **知识点/原理**：RBAC + 双通道模型（管理面/数据面）。
- **代码落点**：RBAC 配置、API server route、连接超时处理。

## 10.5 克隆（Clone）

### 步骤 1：提交 Clone 对象
- **组件动作**：创建 `VirtualMachineClone`，指定源 VM 和目标 VM。
- **知识点/原理**：控制器编排复制流程，不直接复制进程态。
- **代码落点**：Clone CR 构建器、校验 webhook（可选）。

### 步骤 2：配置与存储复制
- **组件动作**：复制 VM 配置与卷数据引用/副本。
- **知识点/原理**：配置复制 + 数据一致性策略。
- **代码落点**：`reconcileClone()`、卷复制子流程。

### 步骤 3：目标 VM 落地
- **组件动作**：创建目标 VM 并更新 Clone 状态。
- **知识点/原理**：长流程任务状态机（`Progressing/Ready/Failed`）。
- **代码落点**：status 条件更新、阶段事件上报。

## 10.6 快照（Snapshot）与恢复（Restore）

### 步骤 1：创建快照/恢复请求
- **组件动作**：创建 `VirtualMachineSnapshot` 或 `VirtualMachineRestore`。
- **知识点/原理**：面向恢复点的声明式对象。
- **代码落点**：snapshot/restore controller 入口。

### 步骤 2：协调卷层操作
- **组件动作**：调用 CSI/存储后端执行快照或回滚。
- **知识点/原理**：控制平面编排 + 存储平面执行。
- **代码落点**：卷操作适配层、存储状态轮询逻辑。

### 步骤 3：状态收敛
- **组件动作**：更新 Ready/Failed，必要时重建实例。
- **知识点/原理**：异步任务必须有清晰可观测状态。
- **代码落点**：`updateSnapshotStatus()`、`updateRestoreStatus()`。

## 10.7 热迁移（Live Migration）

### 步骤 1：迁移请求与前置校验
- **组件动作**：创建 `VirtualMachineInstanceMigration`，校验迁移条件。
- **知识点/原理**：可迁移性取决于网络、存储和实例配置。
- **代码落点**：migration controller 校验器。

### 步骤 2：源/目标节点协同迁移
- **组件动作**：控制器协调源目标节点转移运行上下文。
- **知识点/原理**：不停机迁移核心是内存与执行状态搬迁。
- **代码落点**：迁移阶段状态机、超时与回退逻辑。

### 步骤 3：切换与回收
- **组件动作**：目标接管运行，源侧资源回收。
- **知识点/原理**：切换点一致性与资源清理一致性。
- **代码落点**：迁移完成处理、源侧清理函数。

## 10.8 CPU/内存规格调整

### 步骤 1：规格变更提交
- **组件动作**：更新 VM 规格字段。
- **知识点/原理**：声明式变更，是否热更新取决于能力矩阵。
- **代码落点**：规格 patch/update 层、变更审计日志。

### 步骤 2：运行态应用
- **组件动作**：控制器判定热生效还是重建实例。
- **知识点/原理**：控制面能力与运行时限制的折中。
- **代码落点**：`reconcileSpecChange()`、重启策略分支。

### 步骤 3：结果验证
- **组件动作**：更新状态并验证业务指标。
- **知识点/原理**：容量变更后的性能回归检查。
- **代码落点**：状态更新与指标埋点。

## 10.9 磁盘扩容（系统盘/数据盘）

### 步骤 1：提交容量变更
- **组件动作**：更新 PVC/DataVolume 容量。
- **知识点/原理**：控制面扩容与文件系统扩容是两段流程。
- **代码落点**：存储变更控制器、PVC patch 逻辑。

### 步骤 2：后端卷扩容
- **组件动作**：Longhorn 扩容卷，CSI 完成节点侧扩展。
- **知识点/原理**：块设备扩容成功后仍需客体系统扩文件系统。
- **代码落点**：Longhorn 控制器交互、CSI 事件监听处理。

### 步骤 3：客体系统确认
- **组件动作**：系统内识别并扩展文件系统。
- **知识点/原理**：云原生卷扩容“控制平面完成”不代表“业务立即可用”。
- **代码落点**：运维自动化脚本（cloud-init/agent hooks）。

## 10.10 网络调整（附加网卡/NAD 变更）

### 步骤 1：网络声明变更
- **组件动作**：修改 VM 网络字段或 NAD 配置。
- **知识点/原理**：默认网络与附加网络分工不同。
- **代码落点**：网络 spec patch、NAD 管理器。

### 步骤 2：依赖重检与实例更新
- **组件动作**：控制器检查 NMState/NAD 条件，决定热更新或重建。
- **知识点/原理**：网络变更通常依赖底层节点能力和 CNI 插件能力。
- **代码落点**：`reconcileNetworkChange()`。

### 步骤 3：连通性验证
- **组件动作**：验证网卡、路由、DNS、业务连通。
- **知识点/原理**：网络变更后的验收应分层（L2/L3/L7）。
- **代码落点**：诊断任务/健康检查流程代码。

## 10.11 删除与清理（Delete + Finalizer）

### 步骤 1：删除请求进入
- **组件动作**：对象进入 `deletionTimestamp` 阶段。
- **知识点/原理**：带 Finalizer 的对象不会立即消失。
- **代码落点**：删除分支入口判断（`if !DeletionTimestamp.IsZero()`）。

### 步骤 2：依赖清理
- **组件动作**：按策略清理 VMI、网络对象、卷快照等。
- **知识点/原理**：先清理外部依赖，再移除 Finalizer。
- **代码落点**：`reconcileDelete()`、`cleanupResources()`。

### 步骤 3：终态收敛
- **组件动作**：Finalizer 移除，对象真正删除。
- **知识点/原理**：避免僵尸资源、保证幂等删除。
- **代码落点**：`controllerutil.RemoveFinalizer(...)` + `client.Update(...)`。

---

## 11. 项目代码导航版（vmoperator 实际落点）

这一节把上面的通用步骤映射到你项目里的真实代码。  
项目路径：`/Users/jianfenliu/Workspace/vmoperator`。

## 11.1 主调度入口（所有操作共用）

- **控制器入口**：`internal/controller/wukong_controller.go`
  - `WukongReconciler.Reconcile()`
  - `reconcileNetworks()`
  - `reconcileDisks()`
  - `reconcileVirtualMachine()`
  - `syncNetworkStatusFromVMI()`
  - `updateConditions()`
  - `reconcileDelete()`
- **核心状态字段**：`api/v1alpha1/wukong_types.go`
  - `WukongSpec`（CPU/Memory/Networks/Disks/StartStrategy）
  - `WukongStatus.Phase`（Pending/Creating/Running/Stopped/Error）
  - `WukongStatus.Conditions`（Ready/NetworksConfigured/VolumesBound）

## 11.2 计算面（KubeVirt）真实代码

- **VM 构建与更新**：`pkg/kubevirt/vm.go`
  - `ReconcileVirtualMachine()`
  - `buildVirtualMachine()` / `buildVMSpec()`
  - `buildNetworks()` / `buildInterfaces()` / `buildVolumes()`
  - `updateVMSpec()`
  - `GetVMStatus()`
- **启动相关实现状态**：`已实现（声明驱动）`
  - 通过 `StartStrategy.AutoStart` 与 VM `Running` 字段驱动运行态
  - 当前更偏“声明式启动/停止”，非完整 `virtctl start/stop` 指令编排控制器

## 11.3 网络依赖真实代码（Multus/NMState）

- **网络编排**：`pkg/network/multus.go`
  - `ReconcileNetworks()`（自动创建/复用 NAD）
  - `checkMultusCRDExists()`
  - `buildCNIConfig()`（bridge/ovs 分支 + IPAM 构造）
- **控制器调用点**：`internal/controller/wukong_controller.go:reconcileNetworks()`
- **实现状态**
  - Multus NAD 编排：`已实现`
  - NMState：控制器里有 `network.ReconcileNMState(...)` 调用，当前注释为占位/轻实现（需按你环境确认）

## 11.4 存储依赖真实代码（PVC/DataVolume/扩容）

- **存储主流程**：`pkg/storage/reconcile.go`
  - `ReconcileDisks()`（有镜像走 DataVolume，无镜像走 PVC）
- **DataVolume 流程**：`pkg/storage/datavolume.go`
  - `ReconcileDataVolume()`
  - `CheckDataVolumeStatus()`
  - `DeleteDataVolume()`
- **扩容流程**：`pkg/storage/expand.go`
  - `ExpandPVC()`
  - `CheckPVCExpansionStatus()`
  - `ReconcileDiskExpansion()`
- **控制器调用点**：`internal/controller/wukong_controller.go`
  - `reconcileDisks()`
  - `storage.ReconcileDiskExpansion(...)`

## 11.5 快照相关真实代码

- **快照控制器**：`internal/controller/snapshot_controller.go`
  - `WukongSnapshotReconciler.Reconcile()`
  - 通过 `snapshot.kubevirt.io/VirtualMachineSnapshot`（Unstructured）创建并轮询 `status.readyToUse`
- **对应 CR 定义**：`api/v1alpha1/snapshot_types.go`
- **实现状态**：`已实现（创建 + 状态轮询）`

## 11.6 删除与清理真实代码

- **删除编排主入口**：`internal/controller/wukong_controller.go:reconcileDelete()`
  - 删除 KubeVirt `VirtualMachine`
  - 删除 DataVolume（并间接清理 PVC）
  - 无 ownerReference 的 PVC 兜底删除
  - 检查残留资源 -> Requeue -> 最终移除 Finalizer
- **实现状态**：`已实现（含重试与等待）`

## 11.7 “操作全集”在当前项目中的落地成熟度

- **已有明确代码闭环**
  - 创建/更新（Wukong -> VM）
  - 删除清理
  - 存储创建（PVC/DataVolume）与扩容
  - 网络 NAD 编排
  - 快照创建与状态跟踪
- **部分具备/依赖 KubeVirt 原生能力**
  - 启动/停止（主要通过声明字段驱动）
  - 控制台/SSH（更多依赖 virtctl 与环境配置）
- **文档规划为主（需确认是否已有业务控制器实现）**
  - 克隆、恢复、热迁移、暂停/恢复、强制停止等高级运维动作
  - 这些可先落在 `docs/interview/13-*` 方案中，再逐步沉淀到 controller 代码

## 11.8 开发时如何把“通用步骤”映射到项目代码

1. 先从 `WukongReconciler.Reconcile()` 定位所在阶段（网络/存储/VM/状态）。
2. 再进入对应子模块：
   - 网络看 `pkg/network/multus.go`
   - 存储看 `pkg/storage/*.go`
   - KubeVirt 组装看 `pkg/kubevirt/vm.go`
3. 最后回到 `updateConditions()` 看状态暴露是否符合预期。

这样就能把“面试步骤”直接对齐成“代码执行链路”。
