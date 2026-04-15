# KubeVirt：虚拟机创建依赖（计算控制面）

本文聚焦虚拟机创建流程中的 KubeVirt 组件，重点说明：

1. KubeVirt 在 VM 创建中负责什么
2. 从声明到运行的 1 到 N 步骤
3. 出问题后如何快速定位

---

## 1. 一句话定位

KubeVirt 是 Kubernetes 上的虚拟化控制面，负责把 `VirtualMachine` 声明推进为 `VirtualMachineInstance` 运行实例，并通过 `virt-launcher` 在节点上真正拉起虚机。

---

## 2. 关键组件与对象关系

## 2.1 控制面组件

- `virt-api`：提供 KubeVirt API 能力
- `virt-controller`：推进 VM/VMI 生命周期（核心控制器）
- `virt-handler`：节点侧代理，管理虚机运行时交互

## 2.2 数据面运行单元

- `virt-launcher`：承载 VM 运行实例（qemu/libvirt）

## 2.3 对象关系

```text
VirtualMachine (声明)
  -> VirtualMachineInstance (运行实例)
    -> virt-launcher Pod (节点落地)
```

---

## 3. 从 1 到 N 的创建主线（面试版 8 步）

1. **提交声明**
   - 提交 `VirtualMachine`（或上层业务 CR 触发 VM 创建）。
   - 关键词：声明式期望状态。

2. **准入与持久化**
   - API Server 做认证/授权/准入，写入 etcd。
   - 关键词：对象存在 != 已运行。

3. **控制器 Reconcile**
   - `virt-controller` Watch 到变更，进入 Reconcile。
   - 关键词：幂等循环、持续收敛。

4. **依赖检查**
   - 检查网络与存储依赖（NAD、PVC/DataVolume 等）是否可用。
   - 关键词：未就绪会 Requeue。

5. **生成 VMI**
   - VM 被推进为 VMI，进入可调度状态。
   - 关键词：VM 是声明，VMI 是运行态。

6. **调度绑定**
   - kube-scheduler 为 VMI 选择节点并绑定。
   - 关键词：Filter/Score/Bind。

7. **节点落地运行**
   - 节点 kubelet 拉起 `virt-launcher`，完成网卡与卷接入。
   - 关键词：真正“跑起来”发生在节点。

8. **状态回写收敛**
   - KubeVirt/业务控制器回写 `phase/conditions`，达到 `Running/Ready` 或 `Error`。
   - 关键词：最终一致性。

---

## 4. 与总流程的映射（介入时机）

- `crd-create.md` 中“提交声明/准入持久化”主要由 API Server 完成
- 到“创建运行实例”阶段，KubeVirt 开始成为主角色
- `virt-controller` 负责 VM -> VMI 推进
- `virt-handler + virt-launcher` 负责节点运行态

常见失败信号：

- VM 存在但无 VMI：多为依赖未就绪或控制器未推进
- VMI Pending：多为调度约束/资源不足
- virt-launcher 起不来：多为 CNI/CSI/镜像或权限问题

---

## 5. 高频排障（现象 -> 检查 -> 结论）

## 5.1 VM 创建了但没有 VMI

检查：
```bash
kubectl get vm,vmi -A
kubectl describe vm <vm-name> -n <namespace>
kubectl get pods -n kubevirt
```

结论方向：
- 控制器未推进
- 前置依赖未就绪（网络/存储）

## 5.2 VMI 一直 Pending

检查：
```bash
kubectl describe vmi <vmi-name> -n <namespace>
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
kubectl describe node <node-name>
```

结论方向：
- 调度失败（资源、亲和性、污点）

## 5.3 virt-launcher 启动失败

检查：
```bash
kubectl get pod -A | grep virt-launcher
kubectl describe pod <virt-launcher-pod> -n <namespace>
kubectl logs -n kubevirt -l kubevirt.io=virt-handler --tail=200
```

结论方向：
- 网络附加失败
- 卷挂载失败
- 容器运行时错误

---

## 6. 常用命令清单（面试/排障）

```bash
# 主对象
kubectl get vm,vmi -A
kubectl describe vm <vm-name> -n <namespace>
kubectl describe vmi <vmi-name> -n <namespace>

# KubeVirt 组件
kubectl get pods -n kubevirt
kubectl logs -n kubevirt -l kubevirt.io=virt-controller --tail=200
kubectl logs -n kubevirt -l kubevirt.io=virt-handler --tail=200

# 运行单元与事件
kubectl get pod -A | grep virt-launcher
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
```

---

## 7. 面试口述模板（40 秒）

「KubeVirt 在虚拟机创建里负责计算控制面。流程是：先提交 VM 声明，API Server 持久化后，virt-controller 进入 Reconcile，检查网络和存储依赖，依赖就绪后把 VM 推进成 VMI。随后 scheduler 选节点，kubelet 在节点拉起 virt-launcher，完成网卡和卷接入。最后 KubeVirt 持续回写 phase/conditions，收敛到 Running/Ready；如果失败，就通过 VM/VMI、virt-launcher 和 kubevirt 命名空间日志分层排查。」 
