# NMState + Multus + CoreDNS：虚拟机“准备依赖”详解

本文用于补充虚拟机创建流程中的“准备依赖”阶段，重点回答：

1. 为什么 VM 经常卡在依赖阶段
2. NMState、Multus、CoreDNS 分别负责什么
3. 面试时如何按“时序 + 组件 + 排障”讲清楚

---

## 1. 一句话结论

在 VM 真正启动前，需要先保证**节点网络能力（NMState）**、**多网络接入能力（Multus + NAD）**、**服务发现能力（CoreDNS）**可用；否则即使 VM/VMI 对象创建成功，也可能无法联网或无法解析服务域名。

---

## 2. 组件职责分工

## 2.1 NMState（节点网络先决条件）

- 负责把节点网卡配置声明式下发到主机层（桥接、VLAN、Bond 等）
- 典型用途：为 KubeVirt 桥接网络、SR-IOV/二层网络准备底座
- 关键点：NMState 成功不代表 VM 已联网，只代表“节点网络底座就绪”

常见对象：
- `NodeNetworkConfigurationPolicy`（NNCP）
- `NodeNetworkState`（NNS）

## 2.2 Multus（多网络编排入口）

- CNI 元插件，本身不做具体网络转发
- 作用是把“默认网络 + 附加网络”组合注入到 Pod/virt-launcher
- 通过 `NetworkAttachmentDefinition`（NAD）声明附加网络

关键点：
- Multus 只负责“附加网卡编排”
- 具体网卡创建仍由底层 CNI 插件完成（如 macvlan/bridge/whereabouts 等）

## 2.3 CoreDNS（服务发现）

- 负责集群内 DNS 解析
- 让 VM/Pod 通过服务名访问，例如 `svc.ns.svc.cluster.local`
- 在虚机场景中，常用于访问集群内 API、数据库、消息队列等服务

关键点：
- 网络通了但 DNS 不通，业务仍然“不可用”
- CoreDNS 通常不参与网卡创建，但直接影响业务连通性体验

---

## 3. 准备依赖的推荐时序（面试主线）

1. **先做 NMState**：确认节点桥接/VLAN/Bond 等网络策略已生效  
2. **再做 Multus/NAD**：准备附加网络定义，确保 VM 需要的网段可附加  
3. **检查 CoreDNS**：确保虚机启动后能解析集群服务名  
4. **最后启动 VM/VMI**：由 KubeVirt 推进实例运行并接入网络

> 口诀：**先节点，后网卡，再解析，最后启动。**

---

## 4. 与 VM 创建流程的映射关系

- VM `spec` 声明网络需求（默认网络 + 附加网络）
- 控制器 Reconcile 时先检查 NMState/NAD 等依赖
- 依赖未就绪则 Requeue，不进入 VMI 落地
- 依赖就绪后，virt-launcher 启动并接网
- 运行后通过 CoreDNS 完成服务发现

---

## 5. 典型 YAML 关注点（面试可说）

## 5.1 VM 侧网络声明

- `spec.template.spec.networks`：声明默认/附加网络
- `spec.template.spec.domain.devices.interfaces`：声明网卡模型与绑定方式

## 5.2 NAD 侧网络声明

- `metadata.name` + `metadata.namespace`：被 VM 引用时必须可见
- `spec.config`：底层 CNI 类型、IPAM、网段等配置

## 5.3 NMState 侧声明

- NNCP 描述目标节点的网络期望
- NNS 展示节点当前网络状态（是否已收敛）

---

## 6. 高频故障与排障命令

## 6.1 NMState 未生效

现象：
- VM 附加网络失败
- 节点缺少桥接/目标接口

命令：
```bash
kubectl get nncp,nns -A
kubectl describe nncp <nncp-name>
kubectl describe nns <node-name>
```

## 6.2 Multus/NAD 问题

现象：
- virt-launcher 创建后报 CNI 错误
- 提示 NAD 不存在或网络配置不合法

命令：
```bash
kubectl get network-attachment-definitions -A
kubectl describe pod <virt-launcher-pod> -n <namespace>
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
```

## 6.3 CoreDNS 解析异常

现象：
- 能 ping IP 但服务名无法解析
- 应用报 `no such host` / DNS timeout

命令：
```bash
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns --tail=200
kubectl exec -it <pod-name> -n <namespace> -- nslookup kubernetes.default.svc.cluster.local
```

---

## 7. 面试回答模板（40 秒）

「VM 创建前我会先做依赖准备，顺序是 NMState、Multus、CoreDNS。NMState 负责节点网络底座，比如桥接或 VLAN；Multus 通过 NAD 把附加网卡编排给 virt-launcher；CoreDNS 负责服务发现，保证虚机启动后能用服务名访问集群服务。控制器在 Reconcile 时会先检查这些依赖，未就绪就 Requeue，全部就绪后才推进到 VMI 运行。排障我会按 NNCP/NNS、NAD、CoreDNS 三层去查。」 

---

## 8. 你关心的三个细化问题

## 8.1 VM 创建时，网络组件到底有多少步骤？每步做什么？

下面给出一个可用于口述的 8 步网络链路（从声明到网卡可用）：

1. **声明网络意图（控制面）**
   - 在 VM 中声明 `networks` 与 `interfaces`，标明默认网卡和附加网卡需求。
   - 关键对象：`VirtualMachine.spec.template.spec.networks/interfaces`。

2. **准备节点网络底座（NMState）**
   - 下发 NNCP，把桥接/VLAN/Bond 等主机网络配置到目标节点。
   - 判定信号：NNCP 已应用，NNS 显示节点网络状态收敛。

3. **准备附加网络定义（NAD）**
   - 创建或复用 `NetworkAttachmentDefinition`，定义附加网络类型与 IPAM。
   - 判定信号：NAD 存在且 namespace 正确。

4. **控制器依赖检查（Reconcile）**
   - Operator/KubeVirt 检查 NMState/NAD 是否就绪；未就绪则 Requeue。
   - 判定信号：事件中可看到等待网络依赖的重试记录。

5. **VMI/virt-launcher 进入创建**
   - VM 推进到 VMI，调度后在节点创建 virt-launcher Pod。
   - 判定信号：出现 virt-launcher Pod，状态由 Pending 向 Running 变化。

6. **CNI + Multus 注入网卡**
   - 默认 CNI 先接默认网卡，Multus 再根据 NAD 注入附加网卡。
   - 判定信号：Pod 注解或 CNI 日志出现多网络 attach 结果。

7. **网络策略与连通性生效**
   - 若有 NetworkPolicy/安全策略，开始对流量放行或限制。
   - 判定信号：IP 能否互通、端口是否可达。

8. **业务网络可用验收**
   - 验证 VM 内能访问目标网段、网关、服务地址。
   - 判定信号：ping/curl/tcp 连通成功，接口配置符合预期。

> 面试简化版：**声明 -> NMState -> NAD -> Reconcile 检查 -> VMI 落地 -> CNI/Multus 接网 -> 策略生效 -> 连通性验收**。

## 8.2 涉及系统盘/普通盘时，存储组件多少步骤？分别是什么？

按“系统盘 + 数据盘”统一看，可以拆成 8 步：

1. **声明磁盘需求**
   - 在 VM/上层 CR 里声明系统盘、数据盘容量、访问模式、存储类。
   - 关键对象：`dataVolumeTemplates`、`volumes`、`disks`。

2. **选择存储后端（Longhorn）**
   - 通过 StorageClass 指向 Longhorn，确定卷供应策略。
   - 判定信号：PVC 的 `storageClassName` 正确。

3. **创建 PVC/DataVolume**
   - 系统盘常见 DataVolume（可从镜像导入）；普通盘可用 PVC/DataVolume。
   - 判定信号：PVC/DataVolume 对象创建成功。

4. **卷供应与绑定（Provision + Bind）**
   - Longhorn 控制器创建底层卷，PVC 进入 Bound。
   - 判定信号：`kubectl get pvc` 显示 `Bound`。

5. **调度与 Attach 前置**
   - VMI 调度到节点后，卷进入 attach/mount 准备阶段。
   - 判定信号：事件出现 attach/mount 正常推进日志。

6. **CSI 挂载到 virt-launcher**
   - kubelet 调用 CSI，把系统盘/数据盘挂载到 virt-launcher。
   - 判定信号：Pod 未出现 MountVolume/AttachVolume 错误。

7. **虚机启动并识别磁盘**
   - VM 启动后识别系统盘并挂载数据盘。
   - 判定信号：VMI Running，客体系统看到对应磁盘设备。

8. **状态回写与持续监控**
   - 控制器回写 VolumesReady/Ready，Longhorn 继续维护副本健康。
   - 判定信号：conditions Ready，Longhorn 卷副本健康。

> 面试简化版：**声明磁盘 -> 创建 PVC/DV -> Longhorn 供应绑定 -> CSI 挂载 -> VM 识别磁盘 -> 状态收敛**。

## 8.3 CoreDNS 什么时候介入？什么时候不介入？

CoreDNS 的介入点很关键：

- **不介入的阶段**
  - VM 对象创建、VMI 调度、网卡创建、卷绑定/挂载这些基础资源准备阶段。
  - 这些阶段主要由 API Server、控制器、CNI、CSI、Longhorn 完成。

- **开始介入的阶段**
  - VM 网络已经可用后，业务开始通过“服务名”访问集群服务时。
  - 例如访问 `kubernetes.default.svc.cluster.local` 或内部业务 Service 域名。

- **典型场景**
  1. 仅用 IP 访问：即使 CoreDNS 异常，可能仍可直连 IP。
  2. 用服务名访问：CoreDNS 异常会直接导致 `no such host` 或 DNS 超时。
  3. 业务启动依赖域名：会表现为“虚机已 Running，但应用不可用”。

一句话判断：

**CoreDNS 不决定 VM 能不能“启动”，但决定 VM 能不能“方便且稳定地通过服务名通信”。**

---

## 9. 扩展排障命令（网络/存储/DNS 分层）

```bash
# 网络分层：NMState + NAD + virt-launcher
kubectl get nncp,nns -A
kubectl get network-attachment-definitions -A
kubectl describe pod <virt-launcher-pod> -n <namespace>

# 存储分层：PVC/DV + Longhorn + 事件
kubectl get pvc,pv,dv -A
kubectl get pods -n longhorn-system
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp

# DNS 分层：CoreDNS + 域名解析
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns --tail=200
kubectl exec -it <pod-or-vm-proxy-pod> -n <namespace> -- nslookup kubernetes.default.svc.cluster.local
```
