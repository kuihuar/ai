# Volcano 安装与部署

## 部署前提

- Kubernetes 集群可用（建议 v1.25+）
- 具备 `cluster-admin` 权限
- 集群中已有默认 CNI 和 DNS 组件

## 安装步骤（通用）

1. 安装 Volcano CRD 与控制器组件
2. 部署 Volcano 调度器
3. 校验 `volcano-system` 命名空间组件状态

## 基础验证命令

```bash
kubectl get ns volcano-system
kubectl get pods -n volcano-system
kubectl get crd | rg "volcano.sh"
```

## 与业务命名空间结合

- 为团队创建独立 Queue
- 定义资源上限（CPU/内存/GPU）
- 配置优先级类用于紧急任务抢占

## 生产建议

- 将 Volcano 调度器高可用部署
- 配置监控指标（调度时延、任务等待时长）
- 对关键队列启用容量告警

