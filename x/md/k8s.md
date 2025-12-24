一、基础概念
什么是 Kubernetes？它的核心功能是什么？

答：Kubernetes 是容器编排平台，核心功能包括自动化部署、扩缩容、服务发现、负载均衡、自愈（如重启故障容器）等。

Pod 和容器的区别是什么？

答：Pod 是 Kubernetes 的最小调度单元，包含一个或多个共享网络/存储的容器。容器是实际运行应用的进程。

Deployment、StatefulSet、DaemonSet 的区别？

答：

Deployment：管理无状态应用，支持滚动更新和回滚。

StatefulSet：管理有状态应用（如数据库），提供稳定的网络标识和持久存储。

DaemonSet：确保每个节点运行一个 Pod（如日志收集组件）。

Service 的作用是什么？有哪些类型？

答：Service 提供 Pod 的稳定访问入口，类型包括：

ClusterIP（集群内部访问）、NodePort（通过节点端口暴露）、LoadBalancer（云平台负载均衡器）、ExternalName（映射到外部服务）。

二、核心组件
Kubernetes 主节点（Master）包含哪些组件？

答：

API Server：集群操作的入口。

Scheduler：负责 Pod 调度到节点。

Controller Manager：运行控制器（如 Deployment 控制器）。

etcd：分布式键值存储，保存集群状态。

kubelet 和 kube-proxy 的作用是什么？

答：

kubelet：在节点上管理 Pod 生命周期（如启动/停止容器）。

kube-proxy：维护节点网络规则，实现 Service 的流量转发。

Ingress 和 Ingress Controller 的区别？

答：

Ingress：定义 HTTP/HTTPS 路由规则的资源对象。

Ingress Controller：实际处理流量的组件（如 Nginx、Traefik）。

三、网络与存储
Kubernetes 的网络模型是什么？Pod 之间如何通信？

答：要求所有 Pod 无需 NAT 即可直接通信，通常通过 CNI 插件（如 Calico、Flannel）实现 overlay 网络。

Persistent Volume (PV) 和 Persistent Volume Claim (PVC) 的作用？

答：

PV：集群级别的存储资源（如云盘）。

PVC：用户对存储资源的请求，绑定到 PV。

ConfigMap 和 Secret 的区别？

答：两者均用于配置注入，但 Secret 存储敏感数据（如密码），内容需 Base64 编码。

四、运维与故障排查
如何查看 Pod 的日志？

答：kubectl logs <pod-name>，若 Pod 有多个容器则需 -c <container-name>。

Pod 一直处于 Pending 状态的可能原因？

答：资源不足（CPU/内存）、节点选择器（nodeSelector）不匹配、未绑定持久卷（PVC未找到PV）等。

如何实现滚动更新和回滚？

答：

更新：修改 Deployment 的镜像版本，触发滚动更新。

回滚：kubectl rollout undo deployment/<name>。

五、安全与权限
什么是 RBAC？如何限制用户权限？

答：基于角色的访问控制，通过定义 Role（权限集合）和 RoleBinding（将角色绑定到用户/组）。

Service Account 的作用是什么？

答：为 Pod 中的进程提供身份认证，用于访问 Kubernetes API。

六、高级特性
什么是 Operator 模式？

答：通过自定义控制器（Custom Controller）扩展 Kubernetes，管理复杂有状态应用（如数据库）。

如何实现 HPA（Horizontal Pod Autoscaler）？

答：基于 CPU/内存或自定义指标自动扩缩 Pod 数量，需 Metrics Server 提供指标。

Helm 的作用是什么？

答：Kubernetes 包管理工具，通过 Charts 定义、安装和升级复杂应用。

七、场景题
如何在一个 Pod 中运行多个容器？它们的通信方式？

答：在 Pod 的 containers 字段定义多个容器，通过 localhost 直接通信（共享网络命名空间）。

如何调试一个启动失败的 Pod？

答：

kubectl describe pod <name> 查看事件。

kubectl logs <pod-name> 查看日志。

检查镜像是否存在、资源配额、存储卷挂载等。

八、生态工具
Prometheus 在 Kubernetes 中的作用？

答：监控集群资源、应用性能，通过 Alertmanager 实现告警。

Istio 的核心功能？

答：服务网格，提供流量管理、熔断、链路追踪、安全（mTLS）等。

总结
面试官可能结合实际场景提问（如设计高可用架构、优化资源利用率）。建议深入理解 Kubernetes 核心机制（如控制器循环、声明式 API），并熟悉常见运维操作和故障排查流程。

首先，我需要从基础概念开始，比如什么是Kubernetes，它的架构组件，然后是核心概念如Pod、Deployment、Service这些。然后进阶问题，比如控制器、网络、存储、安全等。还有故障排查和运维相关的问题。

要确保覆盖基础、核心组件、网络、存储、安全、运维、监控、故障排查、扩展性、CI/CD集成等方面。同时，可能需要给出一些具体的问题例子，比如Pod和容器的区别，Service的作用，ConfigMap和Secret的使用场景。

另外，可能用户需要了解不同层次的题目，从初级到高级，所以可以分层次列举。还要注意一些常见的实际操作题，比如如何排查Pod启动失败的问题，或者如何扩展Deployment的副本数。

还要考虑最新版本的特性和最佳实践，比如StatefulSet的应用场景，Helm的作用，以及云原生生态中的工具如Prometheus、Istio等。

需要确保回答结构清晰，分点明确，问题覆盖全面，同时给出参考答案或思路，这样用户才能更好地准备。可能还需要提醒用户不仅要记住答案，还要理解背后的原理和实际应用场景。

检查是否有遗漏的重要主题，比如RBAC、Ingress、持久卷声明、滚动更新策略、资源限制等。确保每个部分的问题都有代表性，能够帮助用户全面复习Kubernetes的知识点。
好的，用户之前问过Kubernetes的面试题，现在要求再深入讲解一下。我需要分析用户的需求，可能他们正在准备面试，需要更详细的知识点或更深入的解释。之前的回答已经覆盖了基础、核心组件、网络存储、运维、安全、高级特性、场景题和生态工具，但用户可能希望更深入每个部分，或者了解更多的实际应用例子。

首先，用户可能希望每个问题有更详细的答案，比如在基础概念中，Pod和容器的区别可能需要更技术性的解释，比如共享网络命名空间和存储卷。另外，关于Deployment、StatefulSet、DaemonSet的区别，可能需要更具体的用例或配置示例。

核心组件部分，用户可能需要了解每个组件如何协同工作，比如API Server的请求处理流程，或者etcd的高可用配置。kubelet和kube-proxy的具体工作机制，比如kubelet如何与容器运行时交互，kube-proxy如何实现iptables或IPVS模式。

网络与存储方面，用户可能想了解不同CNI插件的比较，或者PV/PVC的动态供应流程。ConfigMap和Secret的使用场景和最佳实践，比如如何更新ConfigMap而不重启Pod。

运维与故障排查可能需要更详细的步骤，比如如何分析Pod的Pending状态，具体命令的使用，或者如何查看事件日志。滚动更新和回滚的具体策略，比如maxUnavailable和maxSurge的配置。

安全方面，RBAC的具体配置例子，如何创建Role和RoleBinding，Service Account的令牌生成和使用方式。高级特性中的Operator模式可能需要具体框架如Operator SDK的说明，或者HPA的自定义指标实现。

场景题部分可能需要更多的实际案例，比如如何设计一个高可用应用，或者如何处理节点故障。生态工具如Prometheus和Istio的集成细节，如何在K8s中部署和配置它们。

此外，用户可能希望了解最新的K8s版本的变化，比如v1.25以后的特性，或者常见面试中的陷阱问题，如如何调试服务发现失败，如何处理持久卷的权限问题等。需要确保内容结构清晰，分点详细，同时结合实际操作和命令示例，帮助用户更好地理解和应用知识。

好的！以下是对 Kubernetes 核心概念和常见问题的深入解析，结合实际场景和底层机制，帮助你更全面理解：

一、基础概念深入
1. Pod 的底层实现
共享命名空间：Pod 中的多个容器共享网络命名空间（同一 IP）、IPC 命名空间（进程间通信）和 UTS 命名空间（主机名）。

生命周期：

Pending：调度中或镜像下载中。

Running：至少一个容器启动。

Succeeded/Failed：所有容器终止且状态明确。

Init 容器：在应用容器前运行，用于环境准备（如下载依赖），必须成功退出才会启动主容器。

2. Deployment vs. StatefulSet
Deployment：

无状态设计：Pod 名称随机（如 web-5d89b4d6f-abcx1），无固定顺序。

滚动更新策略：通过 spec.strategy 控制 maxUnavailable（最大不可用 Pod 数）和 maxSurge（超出期望的 Pod 数）。

StatefulSet：

稳定标识：Pod 名称固定（如 redis-0, redis-1），重启后保持不变。

有序部署/扩展：按顺序启停 Pod（如数据库主从）。

持久化存储：通过 volumeClaimTemplates 为每个 Pod 动态创建 PVC，绑定到不同 PV。

3. Service 的流量转发机制
kube-proxy 模式：

iptables：默认模式，通过规则链实现负载均衡（随机选择后端 Pod）。

IPVS：高性能模式，支持更多负载均衡算法（如轮询、最少连接）。

Endpoint 对象：Service 的后端 Pod 列表由 Endpoint 动态维护，当 Pod 变化时自动更新。

二、网络模型深度解析
1. CNI 插件的工作原理
Flannel：

使用 overlay 网络（如 VXLAN），每个节点分配子网，跨节点通信封装在 UDP 包中。

简单易用，但性能损耗较高。

Calico：

基于 BGP 协议的路由方案，无 overlay，性能接近原生网络。

支持网络策略（NetworkPolicy），实现 Pod 级别的流量控制。

2. DNS 服务发现
CoreDNS：

解析规则：<service-name>.<namespace>.svc.cluster.local → Service 的 ClusterIP。

Pod 的 DNS 配置：/etc/resolv.conf 中定义搜索域和 DNS 服务器（如 10.96.0.10）。

3. Ingress 高级配置
多路径路由：

yaml
复制
spec:
  rules:
  - host: demo.com
    http:
      paths:
      - path: /api
        backend:
          serviceName: api-service
          servicePort: 80
      - path: /static
        backend:
          serviceName: static-service
          servicePort: 80
TLS 终止：通过 spec.tls 配置证书，Ingress Controller 负责 HTTPS 解密。

三、存储机制深入
1. PV/PVC 的动态供应（Dynamic Provisioning）
StorageClass：

定义存储类型（如 AWS EBS、GCE PD）。

示例：

yaml
复制
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
2. ConfigMap/Secret 的热更新
ConfigMap 更新后：

如果通过环境变量注入，需重启 Pod 生效。

如果通过 Volume 挂载，kubelet 会定期同步更新文件（可能需要应用支持文件热加载）。

Secret 安全性：

避免在日志或环境中明文暴露。

推荐使用加密存储（如 Kubernetes 的 EncryptionConfiguration）。

四、运维与故障排查进阶
1. Pod 常见故障场景
CrashLoopBackOff：

检查日志：kubectl logs -p <pod-name>（查看上一次崩溃的日志）。

可能原因：应用启动失败、依赖缺失、资源不足（如 OOMKilled）。

ImagePullBackOff：

检查镜像名称、权限（私有仓库需配置 imagePullSecrets）。

2. 节点故障处理
节点 NotReady：

检查 kubelet 状态：systemctl status kubelet。

驱逐 Pod：kubectl drain <node-name> --ignore-daemonsets。

3. 资源配额（Resource Quotas）
限制命名空间资源：

yaml
复制
apiVersion: v1
kind: ResourceQuota
metadata:
  name: dev-quota
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
五、安全机制深入
1. RBAC 配置示例
创建 Role：

yaml
复制
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: dev
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
绑定 Role 到用户：

yaml
复制
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: dev
subjects:
- kind: User
  name: alice
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
2. Pod 安全策略（PSP）
控制 Pod 的权限（如禁止特权容器）：

yaml
复制
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: restricted
spec:
  privileged: false  # 禁止特权模式
  runAsUser:
    rule: MustRunAsNonRoot  # 必须非 root 用户运行
  volumes:
  - configMap
  - secret
六、高级特性与设计模式
1. Operator 模式实现
Custom Resource Definition (CRD)：

定义自定义资源（如 RedisCluster）。

控制器逻辑：

Watch 资源变化，协调实际状态与期望状态（如自动扩缩 Redis 节点）。

2. HPA 自定义指标
基于 Prometheus 的指标扩缩：

yaml
复制
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: custom-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  metrics:
  - type: Pods
    pods:
      metric:
        name: http_requests_per_second
      target:
        type: AverageValue
        averageValue: 100
七、典型场景解决方案
1. 有状态应用部署（如 MySQL 集群）
StatefulSet 配置：

每个 Pod 绑定独立的 PV，持久化数据目录（如 /var/lib/mysql）。

使用 Headless Service（clusterIP: None）提供稳定的 DNS 记录。

2. 金丝雀发布（Canary Release）
通过 Deployment 分批次更新：

先创建新版本 Deployment（副本数较少），逐步调整流量比例。

或结合 Ingress 的流量切分（如 Nginx Annotations）。

八、性能优化与监控
1. 资源请求与限制（Requests/Limits）
Requests：调度依据（如 requests.cpu: 0.5）。

Limits：容器资源上限（如 limits.memory: 1Gi）。

避免资源争抢：合理设置 Limits，防止单个 Pod 耗尽节点资源。

2. 使用 Prometheus + Grafana 监控
监控指标：

集群级别：节点 CPU/内存、Pod 状态、存储使用量。

应用级别：HTTP 请求延迟、错误率。

告警规则示例：

yaml
复制
- alert: HighPodRestartRate
  expr: rate(kube_pod_container_status_restarts_total[5m]) > 0.1
  for: 10m
  labels:
    severity: critical
总结
Kubernetes 的深度理解需结合实践与底层原理：

核心机制：控制器循环（Control Loop）、声明式 API、事件驱动。

故障排查：从事件（kubectl describe）、日志、资源状态层层深入。

扩展性：通过 CRD、Operator、CNI/CSI 插件定制化集群能力。

建议通过实际项目（如搭建高可用集群、部署复杂应用）巩固知识，同时关注社区动态（如 KEPs - Kubernetes Enhancement Proposals）。
https://github.com/edisonwong520/program_term_pronounciation
好的，我现在需要理解用户的问题：“如何理解KubeVirt 基于 QEMU/KVM 或 libvirt 实现虚拟化？”这个问题是在用户之前询问KubeVirt的基本概念之后提出的，显然用户希望更深入地了解KubeVirt如何与底层的虚拟化技术如QEMU/KVM和libvirt集成。我需要确保我的回答既准确又易于理解，同时覆盖用户可能关心的技术细节。

首先，我应该回顾KubeVirt的基本架构，特别是它如何利用现有的虚拟化技术。KubeVirt的目标是在Kubernetes中管理虚拟机，而Kubernetes本身是设计用于容器的，因此必须依赖现有的虚拟化工具来创建和管理虚拟机。QEMU/KVM和libvirt都是常见的开源虚拟化解决方案，我需要解释它们各自的作用以及KubeVirt是如何与它们集成的。

接下来，我需要拆解问题中的关键词：“基于 QEMU/KVM 或 libvirt 实现虚拟化”。这里用户可能想知道KubeVirt是如何具体使用这些工具的，例如是否直接调用它们的API，或者是否有中间层。我需要解释QEMU和KVM的关系，以及libvirt的作用。QEMU是一个机器模拟器，而KVM是Linux内核模块，提供硬件虚拟化支持，两者结合使用可以实现高效的虚拟化。Libvirt则是一个管理虚拟化平台的开源API，守护进程和工具，常用于管理QEMU/KVM虚拟机。

然后，我需要描述KubeVirt组件如何与这些技术交互。例如，KubeVirt的virt-launcher Pod负责启动和管理虚拟机实例，这个Pod内部可能会运行QEMU进程，利用KVM加速。Libvirt可能用于配置和管理这些虚拟机，但需要确认KubeVirt是否直接使用libvirt，还是绕过它直接与QEMU/KVM交互。

用户可能还想了解为什么KubeVirt选择这些技术，而不是其他虚拟化方案，比如VirtualBox或Hyper-V。这里需要提到性能、开源生态、与Linux内核的集成等因素。此外，KubeVirt需要与Kubernetes的调度和资源管理机制集成，QEMU/KVM和libvirt在这些方面的兼容性可能更好。

另一个方面是安全性和隔离性。容器和虚拟机的隔离机制不同，QEMU/KVM提供硬件级别的隔离，而容器依赖于命名空间和cgroups。用户可能想知道KubeVirt如何利用QEMU/KVM来确保虚拟机之间的隔离，以及这与容器隔离的区别。

还需要考虑存储和网络的配置。KubeVirt如何利用Kubernetes的存储卷（如PVC）作为虚拟机的磁盘，这可能涉及到QEMU的存储驱动和libvirt的存储池管理。网络方面，KubeVirt可能使用CNI插件来为虚拟机提供网络连接，类似于容器，但需要处理虚拟机的网络设备配置，这里QEMU的网络后端可能被使用。

此外，用户可能关心性能开销。虽然KVM提供了硬件虚拟化支持，但运行虚拟机仍然会有一定的资源消耗，与容器相比如何？KubeVirt在资源管理方面如何优化，例如CPU和内存的分配，是否支持资源限制和请求，类似容器的requests和limits。

还有，故障排查和监控。使用QEMU/KVM的情况下，如何监控虚拟机的性能，KubeVirt是否提供了与Kubernetes集成的监控工具，或者需要依赖传统的虚拟化监控手段。

最后，我需要总结KubeVirt基于这些技术的原因，以及这种设计带来的优势和潜在的挑战，比如复杂性增加、依赖底层虚拟化工具的稳定性等。

在组织回答时，应该分点说明，先介绍QEMU/KVM和libvirt的基本概念，然后解释KubeVirt如何集成它们，接着讨论这种集成的优势、工作原理、组件交互，最后结合实际用例或性能考量。确保语言通俗易懂，避免过于专业的术语，或在必要时进行解释。

KubeVirt 基于 QEMU/KVM 或 libvirt 实现虚拟化，本质上是将传统虚拟化技术无缝集成到 Kubernetes 生态中，利用 Kubernetes 的调度和管理能力运行虚拟机（VM）。以下是这一设计的关键点解析：

1. QEMU/KVM 与 libvirt 的角色
(1) QEMU/KVM：虚拟化的执行引擎
QEMU：开源的机器模拟器（全称 Quick Emulator），负责虚拟化 CPU、设备（如磁盘、网卡）等硬件资源，支持多种架构（如 x86、ARM）。

KVM（Kernel-based Virtual Machine）：Linux 内核模块，提供硬件辅助虚拟化（如 Intel VT-x/AMD-V），将 QEMU 的指令翻译直接映射到物理 CPU，大幅提升虚拟机性能。

QEMU + KVM 是黄金组合：QEMU 负责设备模拟，KVM 负责 CPU/内存的硬件加速。

(2) libvirt：虚拟化管理工具
功能：提供统一的 API 和工具（如 virsh）管理多种虚拟化技术（如 QEMU/KVM、VMware、Hyper-V）。

作用：简化虚拟机的生命周期管理（创建、启动、停止）、网络/存储配置等。

KubeVirt 的选择：早期 KubeVirt 依赖 libvirt，但新版本逐渐转向直接调用 QEMU/KVM，减少依赖层级。

2. KubeVirt 如何集成 QEMU/KVM？
(1) 架构设计
虚拟机运行在 Pod 中：
KubeVirt 为每个虚拟机创建一个 Pod（名为 virt-launcher），Pod 内运行 qemu 进程，直接调用 KVM 加速。

优势：

利用 Kubernetes 的调度能力（如节点选择、资源限制）。

共享 Kubernetes 的网络（CNI）和存储（CSI）插件。

示例流程：

用户创建 VirtualMachine 资源（CRD）。

virt-controller 生成 VirtualMachineInstance（VMI）。

virt-handler 在节点上启动 virt-launcher Pod，内部调用 qemu-system-x86_64 命令启动虚拟机。

QEMU 进程通过 /dev/kvm 设备文件使用 KVM 加速。

(2) 关键配置
QEMU 命令行参数：
KubeVirt 通过模板生成 QEMU 启动命令，例如：

bash
复制
qemu-system-x86_64 \
  -enable-kvm \                  # 启用 KVM 加速
  -cpu host \                    # 透传物理 CPU 特性
  -m 4096M \                     # 内存分配
  -drive file=/var/lib/kubevirt/disks/disk.img,format=qcow2 \  # 磁盘
  -netdev tap,id=net0 \          # 网络设备
  -device virtio-net-pci,netdev=net0
资源隔离：

通过 Kubernetes 的 limits 和 requests 限制虚拟机的 CPU/内存。

使用 cgroups 和 namespaces 隔离资源（与容器机制一致）。

3. 为什么选择 QEMU/KVM 而不是其他虚拟化技术？
(1) 开源与生态兼容性
完全开源：QEMU/KVM 是 Linux 生态的主流虚拟化方案，与 Kubernetes 的开源理念一致。

社区支持：广泛的开发者社区和成熟的工具链（如 libvirt、virsh）。

(2) 性能与硬件支持
接近原生性能：KVM 的硬件虚拟化支持（Intel VT-x/AMD-V）使虚拟机性能损失极低（通常 <5%）。

设备直通（Passthrough）：支持 GPU、NVMe 磁盘等硬件直接挂载到虚拟机。

(3) 与 Kubernetes 的无缝集成
轻量化：QEMU 进程作为 Pod 中的容器运行，复用 Kubernetes 的存储、网络和安全机制。

动态调度：Kubernetes 根据节点资源情况自动调度虚拟机（如通过 kube-scheduler）。

4. KubeVirt 与 libvirt 的关系演进
早期依赖：KubeVirt 最初通过 libvirt 管理 QEMU/KVM 虚拟机，利用其成熟的 API 简化开发。

去 libvirt 化趋势：

原因：减少依赖层级，提升灵活性和性能。

现状：KubeVirt 逐步转向直接调用 QEMU 命令，仅保留部分 libvirt 工具（如 virsh 用于调试）。

5. 实际工作流程示例
场景：在 Kubernetes 中启动一个 CentOS 虚拟机
定义虚拟机 CRD：

yaml
复制
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: centos-vm
spec:
  running: true
  template:
    spec:
      domain:
        devices:
          disks:
          - name: rootdisk
            disk: {}
        resources:
          requests:
            memory: 2Gi
            cpu: 1
      volumes:
      - name: rootdisk
        persistentVolumeClaim:
          claimName: centos-pvc  # 使用 PVC 作为系统盘
底层执行过程：

virt-launcher Pod 启动后，调用 QEMU 命令：

bash
复制
qemu-system-x86_64 -enable-kvm -cpu host -m 2048M \
  -drive file=/var/lib/kubevirt-disks/centos.qcow2,format=qcow2 \
  -netdev tap,id=net0 -device virtio-net-pci,netdev=net0
QEMU 通过 /dev/kvm 接口使用 KVM 加速。

资源监控：

使用 kubectl top pod 查看 virt-launcher Pod 的资源占用。

通过 virsh 或 virt-top 工具监控虚拟机内部状态。

6. 性能优化与注意事项
(1) 性能调优
启用 CPU 透传（-cpu host）：让虚拟机直接使用物理 CPU 指令集。

使用 Virtio 驱动：优化磁盘和网络 I/O（如 virtio-blk、virtio-net）。

大页内存（HugePages）：减少内存分页开销，提升内存敏感型应用性能。

(2) 安全与隔离
限制特权操作：避免虚拟机获得宿主机 root 权限（通过 Pod 的 SecurityContext 配置）。

使用 Kata Containers：结合轻量级虚拟机实现容器强隔离（可选方案）。

总结
KubeVirt 基于 QEMU/KVM 实现虚拟化的核心逻辑是：

复用 Kubernetes 基础设施（如调度、存储、网络）。

通过 Pod 封装 QEMU 进程，直接调用 KVM 加速。

平衡性能与灵活性：既支持遗留虚拟机负载，又能与容器共存。

这种设计使 Kubernetes 真正成为 混合工作负载（容器 + 虚拟机）的统一平台，适用于传统应用迁移、GPU 直通、边缘计算等复杂场景。
