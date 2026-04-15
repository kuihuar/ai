# Longhorn：虚拟机创建依赖（存储链路）

本文聚焦虚拟机创建流程中的 Longhorn 存储依赖，重点说明：

1. Longhorn 在 VM 创建中负责什么
2. 系统盘和普通盘各自如何落地
3. 卷卡住时怎么排障

---

## 1. 一句话定位

Longhorn 在虚拟机创建中承担分布式块存储角色，负责卷供应、卷副本维护、attach/mount，并通过 CSI 把系统盘/数据盘挂载到 `virt-launcher`，为 VM 启动与数据读写提供底层存储。

---

## 2. 关键对象关系（存储视角）

```text
StorageClass(Longhorn)
  -> PVC / DataVolume
    -> PV (Longhorn volume)
      -> virt-launcher Pod 挂载
        -> VM 识别系统盘/数据盘
```

关键对象：

- `StorageClass`：指定 Longhorn 作为供应后端
- `PersistentVolumeClaim`（PVC）：声明需要多大卷
- `DataVolume`（DV，可选）：常用于镜像导入场景
- `PersistentVolume`（PV）：实际供应的卷

---

## 3. 系统盘 vs 普通盘（面试高频）

## 3.1 系统盘

- 常见来源：镜像导入（DataVolume）或模板克隆
- 关注点：导入速度、镜像可达性、启动盘可引导
- 风险点：镜像导入失败会导致 VM 无法启动

## 3.2 普通盘（数据盘）

- 常见来源：空卷（PVC）或已有卷复用
- 关注点：容量、性能、副本数、读写模式
- 风险点：卷未绑定或未挂载时，VM 可能启动但业务盘不可用

---

## 4. 从 1 到 N 的存储主线（面试版 8 步）

1. **声明磁盘需求**
   - 在 VM/上层 CR 声明系统盘、数据盘、容量与存储类。
   - 关键词：声明式期望。

2. **创建 PVC 或 DataVolume**
   - 系统盘常见 DV（导入镜像），数据盘常见 PVC（空盘或复用）。
   - 关键词：对象先存在，数据后就绪。

3. **Longhorn 动态供应卷**
   - Longhorn 根据 StorageClass 创建底层卷并分配副本。
   - 关键词：供应（provision）。

4. **PVC/PV 绑定完成**
   - PVC 状态变为 `Bound`，进入可挂载阶段。
   - 关键词：绑定是启动前门槛。

5. **VMI 调度后执行 attach**
   - VMI 被调度到节点后，卷开始 attach 到目标节点。
   - 关键词：先调度再 attach。

6. **kubelet/CSI 挂载到 virt-launcher**
   - 节点侧 CSI 完成卷挂载，virt-launcher 可访问块设备。
   - 关键词：attach != mount，二者都要成功。

7. **虚机识别磁盘**
   - VM 启动后识别系统盘并挂载数据盘。
   - 关键词：客体系统可见盘才算可用。

8. **状态回写与收敛**
   - 控制器回写 `VolumesReady/Ready`，最终 `Running` 或 `Error`。
   - 关键词：最终一致性收敛。

---

## 5. 与总流程的映射（介入时机）

- 在 `crd-create.md` 的“准备依赖”阶段，Longhorn 是存储主责任组件
- 在 VM -> VMI 推进前，PVC/DV 通常要达到可用状态
- 在节点落地阶段，CSI 挂载成功与否直接决定 VM 是否能正常读写磁盘

常见失败信号：

- PVC 一直 Pending：供应链路有问题（StorageClass/容量/后端状态）
- VMI 卡在启动：卷 attach/mount 异常
- VM Running 但业务异常：数据盘未识别或文件系统未就绪

---

## 6. 高频排障（现象 -> 检查 -> 结论）

## 6.1 PVC 长时间 Pending

检查：
```bash
kubectl get pvc,pv -A
kubectl describe pvc <pvc-name> -n <namespace>
kubectl get sc
```

结论方向：
- StorageClass 配置不对
- 后端容量不足或供应失败

## 6.2 VMI 启动时报卷挂载错误

检查：
```bash
kubectl describe vmi <vmi-name> -n <namespace>
kubectl describe pod <virt-launcher-pod> -n <namespace>
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
```

结论方向：
- 卷 attach 失败
- CSI mount 失败

## 6.3 Longhorn 后端异常

检查：
```bash
kubectl get pods -n longhorn-system
kubectl logs -n longhorn-system -l app=longhorn-manager --tail=200
kubectl get volumes.longhorn.io -n longhorn-system
```

结论方向：
- 副本不健康
- manager/engine 组件异常

---

## 7. 常用命令清单（面试/排障）

```bash
# 声明与绑定
kubectl get sc
kubectl get pvc,pv,dv -A
kubectl describe pvc <pvc-name> -n <namespace>

# 运行态关联
kubectl get vm,vmi -A
kubectl get pod -A | grep virt-launcher
kubectl describe pod <virt-launcher-pod> -n <namespace>

# Longhorn 组件与卷状态
kubectl get pods -n longhorn-system
kubectl get volumes.longhorn.io -n longhorn-system
kubectl logs -n longhorn-system -l app=longhorn-manager --tail=200

# 事件
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
```

---

## 8. 面试口述模板（40 秒）

「Longhorn 在虚拟机创建里负责存储依赖。流程是先在 VM 里声明系统盘和数据盘，然后创建 PVC 或 DataVolume，Longhorn 动态供应卷并完成 PVC/PV 绑定。VMI 调度到节点后，CSI 把卷 attach/mount 到 virt-launcher，虚机启动后识别系统盘和业务盘。最后控制器回写 VolumesReady/Ready 收敛状态。排障我会先看 PVC 是否 Bound，再看 virt-launcher 的挂载事件，最后看 longhorn-system 里的 manager 和卷副本健康。」 
