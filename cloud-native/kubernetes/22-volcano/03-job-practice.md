# Volcano 批任务与队列实战

## 实战目标

创建一个带队列与任务组约束的批处理作业，理解 Volcano 如何做整体调度。

## 示例流程

1. 创建 Queue 并设置资源权重
2. 提交 Volcano Job
3. 观察 PodGroup 是否满足最小可运行副本
4. 查看任务完成状态与重试行为

## Queue 示例（简化）

```yaml
apiVersion: scheduling.volcano.sh/v1beta1
kind: Queue
metadata:
  name: ml-queue
spec:
  weight: 1
  reclaimable: true
```

## Job 示例（简化）

```yaml
apiVersion: batch.volcano.sh/v1alpha1
kind: Job
metadata:
  name: demo-batch
spec:
  minAvailable: 2
  schedulerName: volcano
  tasks:
    - replicas: 2
      name: worker
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: worker
              image: busybox
              command: ["sh", "-c", "echo hello && sleep 30"]
```

## 排查建议

- Job 长时间 Pending：检查 Queue 配额和节点空闲资源
- 仅部分 Pod 启动：检查 `minAvailable` 与 Gang Scheduling 约束
- 高频失败重试：检查容器命令、镜像、依赖服务可达性

