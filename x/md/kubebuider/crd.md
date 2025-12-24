1. 什么是 Kubebuilder？它的核心功能是什么？
答案：
Kubebuilder 是一个基于 Go 语言的 SDK，用于简化 Kubernetes 自定义资源（CRD） 和 控制器（Controller） 的开发。它提供脚手架工具和代码库，帮助开发者快速构建符合 Kubernetes API 规范的扩展组件。

核心功能：

生成 CRD 的 YAML 定义和 Go 代码框架。

自动化控制器的开发，集成调谐循环（Reconcile Loop）。

支持 Webhook（验证、默认值设置、转换等）。

提供与 Kubernetes API 交互的客户端工具。

2. CRD 是什么？其核心字段有哪些？
答案：
CRD（Custom Resource Definition）允许用户在 Kubernetes 中定义新的资源类型，扩展集群的功能。

核心字段：

Spec：用户定义的期望状态（如配置参数）。

Status：由控制器维护的实际状态（如运行状态、错误信息）。

Metadata：元数据（名称、命名空间、标签等）。

3. 如何用 Kubebuilder 创建一个 CRD？
步骤：

初始化项目：

kubebuilder init --domain my.domain
创建 API（CRD + Controller）：

kubebuilder create api --group demo --version v1 --kind Example1
定义 Spec 和 Status：
修改 api/v1/example1_types.go：

4. 控制器（Controller）的核心逻辑是什么？
答案：
控制器通过 调谐循环（Reconcile Loop） 监听资源变化，确保实际状态（Status）与期望状态（Spec）一致。

关键逻辑：

监听资源（CRD）的创建、更新、删除事件。

触发 Reconcile 方法，执行状态同步逻辑。

处理错误并决定重试策略。

5. 如何实现 Webhook 的验证和默认值设置？
答案：
Webhook 用于在资源操作（创建、更新、删除）时执行自定义逻辑。

6. Finalizer 的作用及实现方式
答案：
Finalizer 用于在资源删除前执行清理逻辑（如释放外部资源）。
