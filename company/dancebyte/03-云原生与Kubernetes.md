# 云原生与 Kubernetes 面试题（含参考答案）

## 1. Pod、Deployment、StatefulSet 区别？
**答：**Pod 是最小调度单元；Deployment 适合无状态副本；StatefulSet 适合有状态、稳定标识与存储绑定场景。

## 2. 为什么需要 Service？
**答：**为动态 Pod 提供稳定访问入口与负载分发能力，解耦实例变化。

## 3. Ingress 和 Service 区别？
**答：**Service 解决 L4 访问与服务抽象；Ingress 解决 L7 路由（Host/Path/TLS）。

## 4. K8s 调度核心流程？
**答：**Filter 过滤可选节点 -> Score 打分 -> Bind 绑定节点。

## 5. 常见工作负载 Pending 原因？
**答：**资源不足、亲和/反亲和不满足、污点未容忍、PVC 未绑定、镜像拉取失败。

## 6. CRD + Controller 价值？
**答：**把业务对象声明化，通过 Reconcile 持续收敛，实现平台自动化运维。

## 7. Operator 的关键能力？
**答：**声明式接口、状态机、幂等 Reconcile、依赖编排、故障重试、状态回写。

## 8. 为什么要用 Helm？
**答：**统一应用打包、参数化部署、版本回滚，降低多环境发布复杂度。

## 9. K8s 网络三件套 CRI/CNI/CSI 是什么？
**答：**CRI 管容器运行时，CNI 管网络接入，CSI 管存储卷供应与挂载。

## 10. 如何做 Kubernetes 故障排查？
**答：**先看事件与对象状态，再看 Pod 日志，再查节点与资源，再回到配置与代码。

## 高频追问（含回答要点）

### 追问 1：Deployment 滚动发布卡住了你先看哪里？
**要点：**
- `kubectl rollout status/history`
- Pod 就绪探针和启动日志
- 资源配额与调度事件

### 追问 2：CRD 控制器为什么强调“幂等”？
**要点：**
- Reconcile 会被反复触发
- 任意中断后要可重入恢复
- 幂等是最终一致性前提

### 追问 3：如何区分网络问题还是应用问题？
**要点：**
- 先 DNS（CoreDNS）再 Service/Endpoint
- 再看 Pod 内连通性（curl/nslookup）
- 最后回应用日志与探针状态



### 追问 4：k8s 资源创建流程？
**要点：**
- 资源声明
- API Server 接收
- Controller 处理
- 事件触发
- 资源状态更新
- 通知其他组件 

####  控制面流程
1. 客户端提交资源声明（kubectl or API）
2. API Server 接收并校验
   - 认证
   - 授权
   - 校验
   - 准入
3. 持久化
4. 调度器决策 scheduler
   - Filter(预筛选)
   - Score(打分)
   - Bind(绑定节点)
5. 更新ETCD（API SERVER 将nodeName 更新到资源声明中，写入ETCD中）
####  工作节点流程
1. Kubelet 监听 到API server 变化,发瑞POD 创建事件
2. 创建 容器沙箱（Pause Container）
    - kubetl 调用CRI API 启动 CRI 运行时
3. 调用 CNI API 配置网络
4. 调用 CSI API 挂载存储卷
5. 拉取应用镜像
6. 启动应用容器
7. 更新 Pod 状态为 Running
####  Pod 生命周期
- Pending
- Running
- Succeeded
- Failed
- Unknown
