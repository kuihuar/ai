# Job 和 CronJob 详解

## 🚀 什么是 Job？

Job 是 Kubernetes 中用于运行一次性任务的控制器。Job 会创建一个或多个 Pod 来执行任务，当任务成功完成后，Job 会标记为完成状态。如果任务失败，Job 会根据配置的重试策略进行重试。

## 🎯 Job 特点

- **一次性执行**：任务完成后不会重新运行
- **并行执行**：可以配置多个 Pod 并行执行任务
- **重试机制**：支持配置重试次数和策略
- **完成状态**：任务完成后保持完成状态
- **清理策略**：支持自动清理完成的 Job

## 📝 Job 配置

### 基础 Job 配置

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl:5.34
        command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

### 并行 Job 配置

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: parallel-job
spec:
  parallelism: 3
  completions: 6
  template:
    spec:
      containers:
      - name: worker
        image: busybox
        command: ["sh", "-c", "echo Processing item $ITEM && sleep 10"]
        env:
        - name: ITEM
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
      restartPolicy: Never
  backoffLimit: 3
```

### 带资源限制的 Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: resource-limited-job
spec:
  template:
    spec:
      containers:
      - name: worker
        image: nginx:1.21
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        command: ["sh", "-c", "echo 'Job completed successfully'"]
      restartPolicy: Never
  backoffLimit: 2
```

## 🚀 什么是 CronJob？

CronJob 是 Kubernetes 中用于创建定时任务的控制器。它基于 Job 构建，可以按照 Cron 表达式定期创建 Job 来执行任务。

## 🎯 CronJob 特点

- **定时执行**：按照 Cron 表达式定期执行
- **基于 Job**：每次执行创建一个新的 Job
- **并发控制**：可以控制同时运行的 Job 数量
- **历史记录**：保留成功和失败的 Job 历史
- **时区支持**：支持配置时区

## 📝 CronJob 配置

### 基础 CronJob 配置

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox:1.28
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```

### 带并发控制的 CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-job
spec:
  schedule: "0 2 * * *"
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:13
            command:
            - /bin/bash
            - -c
            - |
              echo "Starting backup at $(date)"
              pg_dump -h postgres-service -U admin mydb > /backup/backup-$(date +%Y%m%d).sql
              echo "Backup completed at $(date)"
            env:
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: password
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-pvc
          restartPolicy: OnFailure
```

### 带时区配置的 CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: timezone-job
spec:
  schedule: "0 9 * * 1-5"
  timeZone: "Asia/Shanghai"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: timezone-worker
            image: busybox:1.28
            command:
            - /bin/sh
            - -c
            - echo "Current time: $(date)"
          restartPolicy: OnFailure
```

## 🛠️ Job 和 CronJob 操作

### 1. 创建 Job

```bash
# 使用 YAML 文件创建
kubectl apply -f job.yaml

# 使用命令行创建
kubectl create job pi --image=perl:5.34 -- perl -Mbignum=bpi -wle 'print bpi(2000)'
```

### 2. 查看 Job 状态

```bash
# 查看所有 Job
kubectl get jobs

# 查看 Job 详情
kubectl describe job <job-name>

# 查看 Job 的 Pod
kubectl get pods -l job-name=<job-name>
```

### 3. 创建 CronJob

```bash
# 使用 YAML 文件创建
kubectl apply -f cronjob.yaml

# 使用命令行创建
kubectl create cronjob hello --image=busybox:1.28 --schedule="*/1 * * * *" -- /bin/sh -c 'date; echo Hello from the Kubernetes cluster'
```

### 4. 查看 CronJob 状态

```bash
# 查看所有 CronJob
kubectl get cronjobs

# 查看 CronJob 详情
kubectl describe cronjob <cronjob-name>

# 查看 CronJob 创建的 Job
kubectl get jobs -l cronjob=<cronjob-name>
```

### 5. 手动触发 CronJob

```bash
# 立即触发 CronJob
kubectl create job --from=cronjob/<cronjob-name> <job-name>
```

### 6. 删除 Job 和 CronJob

```bash
# 删除 Job
kubectl delete job <job-name>

# 删除 CronJob
kubectl delete cronjob <cronjob-name>

# 删除 CronJob 但保留已创建的 Job
kubectl delete cronjob <cronjob-name> --cascade=orphan
```

## 🔧 实际应用场景

### 1. 数据备份任务

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: database-backup
spec:
  schedule: "0 2 * * *"
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:13
            command:
            - /bin/bash
            - -c
            - |
              echo "Starting database backup at $(date)"
              pg_dump -h postgres-service -U admin mydb | gzip > /backup/backup-$(date +%Y%m%d_%H%M%S).sql.gz
              echo "Backup completed at $(date)"
            env:
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: password
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-pvc
          restartPolicy: OnFailure
```

### 2. 日志清理任务

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: log-cleanup
spec:
  schedule: "0 0 * * 0"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: cleanup
            image: busybox:1.28
            command:
            - /bin/sh
            - -c
            - |
              echo "Starting log cleanup at $(date)"
              find /var/log -name "*.log" -mtime +7 -delete
              echo "Log cleanup completed at $(date)"
            volumeMounts:
            - name: log-storage
              mountPath: /var/log
          volumes:
          - name: log-storage
            hostPath:
              path: /var/log
          restartPolicy: OnFailure
```

### 3. 健康检查任务

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: health-check
spec:
  schedule: "*/5 * * * *"
  concurrencyPolicy: Replace
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: health-check
            image: curlimages/curl:latest
            command:
            - /bin/sh
            - -c
            - |
              echo "Starting health check at $(date)"
              curl -f http://web-service:80/health || exit 1
              echo "Health check passed at $(date)"
          restartPolicy: OnFailure
```

### 4. 批处理数据处理任务

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: data-processing
spec:
  parallelism: 5
  completions: 10
  template:
    spec:
      containers:
      - name: processor
        image: python:3.9
        command:
        - python
        - -c
        - |
          import time
          import os
          import random
          print(f"Processing item {os.environ.get('JOB_COMPLETION_INDEX', 'unknown')}")
          time.sleep(random.randint(10, 30))
          print("Processing completed")
        env:
        - name: JOB_COMPLETION_INDEX
          valueFrom:
            fieldRef:
              fieldPath: metadata.annotations['batch.kubernetes.io/job-completion-index']
      restartPolicy: Never
  backoffLimit: 3
```

## 🎯 练习

### 练习 1：基础 Job
1. 创建一个计算 π 值的 Job
2. 查看 Job 执行状态
3. 验证任务完成

### 练习 2：并行 Job
1. 创建一个并行处理数据的 Job
2. 配置多个 Pod 并行执行
3. 验证所有任务完成

### 练习 3：定时任务 CronJob
1. 创建一个每分钟执行的 CronJob
2. 查看定时执行情况
3. 测试手动触发

### 练习 4：数据备份 CronJob
1. 创建一个每日备份的 CronJob
2. 配置持久化存储
3. 验证备份功能

## 🔍 故障排查

### 常见问题

1. **Job 执行失败**
   ```bash
   # 查看 Job 状态和事件
   kubectl describe job <job-name>
   kubectl get pods -l job-name=<job-name>
   kubectl logs <pod-name>
   ```

2. **CronJob 不执行**
   ```bash
   # 检查 CronJob 状态
   kubectl describe cronjob <cronjob-name>
   kubectl get cronjobs
   ```

3. **资源不足**
   ```bash
   # 检查节点资源
   kubectl top nodes
   kubectl describe nodes
   ```

4. **权限问题**
   ```bash
   # 检查 ServiceAccount 和 RBAC
   kubectl get serviceaccount
   kubectl describe clusterrolebinding
   ```

## 📚 相关资源

- [Kubernetes Job 官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
- [Kubernetes CronJob 官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
- [Cron 表达式参考](https://en.wikipedia.org/wiki/Cron)

## 🎯 下一步学习

掌握 Job 和 CronJob 后，继续学习：
- [Storage](./08-storage/README.md) - 存储管理
- [ConfigMap 和 Secret](./07-config/README.md) - 配置管理
- [Service](./06-service/README.md) - 服务发现和负载均衡
