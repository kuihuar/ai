# 故障排查

## 🔍 故障排查概述

Kubernetes 故障排查是运维工作中的重要技能。掌握系统性的排查方法，能够快速定位和解决问题，确保集群和应用稳定运行。

## 🎯 排查方法论

### 1. 排查步骤
1. **收集信息**: 了解问题现象和影响范围
2. **分析日志**: 查看相关组件的日志
3. **检查状态**: 验证资源状态和配置
4. **测试连通性**: 确认网络和通信正常
5. **对比正常**: 与正常状态进行对比
6. **逐步排查**: 从简单到复杂逐步排查

### 2. 排查工具
- **kubectl**: 基础排查命令
- **日志工具**: 日志收集和分析
- **监控工具**: 性能指标监控
- **网络工具**: 网络连通性测试

## 🛠️ 基础排查命令

### 1. 集群状态检查
```bash
# 检查集群状态
kubectl cluster-info

# 检查节点状态
kubectl get nodes
kubectl describe nodes

# 检查组件状态
kubectl get componentstatuses
kubectl get cs

# 检查 API 服务器
kubectl get --raw='/readyz?verbose'
```

### 2. 资源状态检查
```bash
# 查看所有资源
kubectl get all --all-namespaces

# 查看特定命名空间
kubectl get all -n default

# 查看资源详情
kubectl describe pod <pod-name>
kubectl describe service <service-name>
kubectl describe deployment <deployment-name>

# 查看资源 YAML
kubectl get pod <pod-name> -o yaml
```

### 3. 日志查看
```bash
# 查看 Pod 日志
kubectl logs <pod-name>
kubectl logs <pod-name> -f
kubectl logs <pod-name> --previous

# 查看多个容器
kubectl logs <pod-name> -c <container-name>

# 查看事件
kubectl get events --sort-by='.lastTimestamp'
kubectl get events -n <namespace>
```

## 🚨 常见问题排查

### 1. Pod 启动失败

#### 问题现象
```bash
# Pod 状态为 Pending 或 Failed
kubectl get pods
NAME                     READY   STATUS    RESTARTS   AGE
myapp-pod               0/1     Pending   0          5m
```

#### 排查步骤
```bash
# 1. 查看 Pod 详情
kubectl describe pod myapp-pod

# 2. 检查事件
kubectl get events --field-selector involvedObject.name=myapp-pod

# 3. 检查节点资源
kubectl describe nodes

# 4. 检查镜像
kubectl get pod myapp-pod -o yaml | grep image
```

#### 常见原因
- **资源不足**: CPU 或内存不足
- **镜像拉取失败**: 镜像不存在或网络问题
- **节点污点**: 节点有污点，Pod 无法调度
- **存储问题**: PVC 绑定失败

### 2. 服务无法访问

#### 问题现象
```bash
# 服务无法访问
curl http://service-name
# 连接超时或错误
```

#### 排查步骤
```bash
# 1. 检查服务状态
kubectl get svc
kubectl describe svc service-name

# 2. 检查 Endpoints
kubectl get endpoints service-name

# 3. 检查 Pod 状态
kubectl get pods -l app=myapp

# 4. 测试 Pod 连通性
kubectl exec -it <pod-name> -- curl localhost:8080
```

#### 常见原因
- **Pod 未运行**: 后端 Pod 未启动
- **端口不匹配**: 服务端口与 Pod 端口不匹配
- **标签选择器错误**: 服务无法找到后端 Pod
- **网络策略**: 网络策略阻止访问

### 3. 应用性能问题

#### 问题现象
```bash
# 应用响应慢或超时
# 资源使用率高
```

#### 排查步骤
```bash
# 1. 检查资源使用
kubectl top pods
kubectl top nodes

# 2. 查看资源限制
kubectl describe pod <pod-name> | grep -A 5 Resources

# 3. 检查日志
kubectl logs <pod-name> --tail=100

# 4. 检查网络
kubectl exec -it <pod-name> -- netstat -tulpn
```

#### 常见原因
- **资源限制**: CPU 或内存限制过低
- **网络延迟**: 网络连接问题
- **应用问题**: 应用本身性能问题
- **存储 I/O**: 存储性能问题

## 🔧 高级排查技巧

### 1. 调试容器
```bash
# 进入容器调试
kubectl exec -it <pod-name> -- /bin/bash
kubectl exec -it <pod-name> -- /bin/sh

# 在容器中执行命令
kubectl exec <pod-name> -- ps aux
kubectl exec <pod-name> -- netstat -tulpn
kubectl exec <pod-name> -- df -h
```

### 2. 端口转发
```bash
# 端口转发到本地
kubectl port-forward pod/<pod-name> 8080:80
kubectl port-forward svc/<service-name> 8080:80

# 访问本地端口
curl http://localhost:8080
```

### 3. 临时 Pod 调试
```bash
# 创建调试 Pod
kubectl run debug-pod --image=busybox --rm -it --restart=Never -- sh

# 使用 kubectl debug
kubectl debug <pod-name> -it --image=busybox --target=<container-name>
```

### 4. 网络调试
```bash
# 测试 DNS 解析
kubectl run test-dns --image=busybox --rm -it --restart=Never -- nslookup kubernetes.default

# 测试网络连通性
kubectl run test-connectivity --image=busybox --rm -it --restart=Never -- wget -O- http://service-name

# 检查网络策略
kubectl get networkpolicies --all-namespaces
```

## 📊 监控和告警

### 1. 资源监控
```bash
# 查看资源使用情况
kubectl top pods --all-namespaces
kubectl top nodes

# 查看资源配额
kubectl describe resourcequota --all-namespaces
kubectl describe limitrange --all-namespaces
```

### 2. 健康检查
```bash
# 检查 Pod 健康状态
kubectl get pods -o wide
kubectl describe pod <pod-name> | grep -A 10 "Events:"

# 检查服务健康状态
kubectl get endpoints
kubectl describe endpoints <service-name>
```

### 3. 日志分析
```bash
# 实时查看日志
kubectl logs -f deployment/<deployment-name>

# 查看错误日志
kubectl logs <pod-name> | grep ERROR
kubectl logs <pod-name> | grep -i error

# 日志聚合
kubectl logs --all-containers=true -l app=myapp
```

## 🛡️ 故障预防

### 1. 健康检查配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
```

### 2. 资源限制
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

### 3. 自动扩缩容
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## 🛠️ 排查工具

### 1. kubectl 插件
```bash
# 安装 kubectl 插件
kubectl krew install access-matrix
kubectl krew install resource-capacity
kubectl krew install view-secret

# 使用插件
kubectl access-matrix
kubectl resource-capacity
kubectl view-secret <secret-name>
```

### 2. 第三方工具
```bash
# Lens - Kubernetes IDE
# 图形化界面，便于排查

# K9s - 终端 UI
# 实时监控和操作

# Popeye - 集群健康检查
# 自动检查集群问题
```

## 🎯 最佳实践

### 1. 排查流程
- 建立标准排查流程
- 记录排查步骤和结果
- 建立知识库和文档

### 2. 监控告警
- 设置合理的告警阈值
- 配置多级告警
- 建立告警升级机制

### 3. 日志管理
- 统一日志格式
- 配置日志轮转
- 建立日志分析流程

### 4. 备份恢复
- 定期备份配置和数据
- 测试恢复流程
- 建立灾难恢复计划

## 🛠️ 实践练习

### 练习 1：Pod 故障排查
1. 创建有问题的 Pod
2. 使用排查命令诊断
3. 修复问题并验证

### 练习 2：服务故障排查
1. 创建服务访问问题
2. 排查网络连通性
3. 修复服务配置

### 练习 3：性能问题排查
1. 模拟性能问题
2. 使用监控工具分析
3. 优化资源配置

## 📚 扩展阅读

- [Kubernetes 故障排查官方文档](https://kubernetes.io/docs/tasks/debug/)
- [调试应用](https://kubernetes.io/docs/tasks/debug-application-cluster/)
- [调试集群](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-cluster/)

## 🎯 下一步

掌握故障排查后，继续学习：
- [性能优化](./15-optimization/README.md)
- [生产环境最佳实践](./projects/) 