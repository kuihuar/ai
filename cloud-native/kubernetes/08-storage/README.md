# 存储管理

## 📖 存储概述

Kubernetes 提供了多种存储解决方案来满足不同应用的需求。从临时存储到持久化存储，从本地存储到分布式存储，Kubernetes 支持各种存储类型。

## 🎯 存储类型

### 1. 临时存储
- **emptyDir**: Pod 生命周期内的临时存储
- **hostPath**: 挂载主机文件系统
- **tmpfs**: 内存文件系统

### 2. 持久化存储
- **PersistentVolume (PV)**: 集群级别的存储资源
- **PersistentVolumeClaim (PVC)**: 用户对存储的请求
- **StorageClass**: 动态供应存储

### 3. 特殊存储
- **ConfigMap**: 配置文件存储
- **Secret**: 敏感数据存储
- **Downward API**: 容器信息存储

## 💾 临时存储

### 1. emptyDir
Pod 创建时创建，Pod 删除时销毁的临时存储。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    volumeMounts:
    - name: cache-volume
      mountPath: /cache
  volumes:
  - name: cache-volume
    emptyDir: {}
```

### 2. hostPath
挂载主机文件系统，数据持久化在主机上。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    volumeMounts:
    - name: host-volume
      mountPath: /host-data
  volumes:
  - name: host-volume
    hostPath:
      path: /data
      type: Directory
```

### 3. tmpfs
内存文件系统，数据存储在内存中。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    volumeMounts:
    - name: tmp-volume
      mountPath: /tmp
  volumes:
  - name: tmp-volume
    emptyDir:
      medium: Memory
      sizeLimit: "100Mi"
```

## 🔗 持久化存储

### 1. PersistentVolume (PV)
集群级别的存储资源，由管理员创建。

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-example
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
    - ReadOnlyMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: fast
  hostPath:
    path: /data
```

### 2. PersistentVolumeClaim (PVC)
用户对存储的请求，类似于 Pod 对 Node 的请求。

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-example
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: fast
```

### 3. 在 Pod 中使用 PVC
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    volumeMounts:
    - name: data-volume
      mountPath: /data
  volumes:
  - name: data-volume
    persistentVolumeClaim:
      claimName: pvc-example
```

## 🏭 StorageClass

StorageClass 用于动态供应存储，支持多种存储后端。

### 1. 本地存储 StorageClass
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```

### 2. NFS StorageClass
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nfs-storage
provisioner: example.com/nfs
parameters:
  server: nfs-server.example.com
  path: /exports
```

### 3. 云存储 StorageClass
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
  fsType: ext4
```

## 🔄 动态供应

### 1. 自动创建 PV
当创建 PVC 时，StorageClass 会自动创建对应的 PV。

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: dynamic-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: fast-ssd  # 指定 StorageClass
```

### 2. 默认 StorageClass
设置默认 StorageClass，PVC 可以不指定 storageClassName。

```bash
# 设置默认 StorageClass
kubectl patch storageclass fast-ssd -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

## 📊 访问模式

### 1. ReadWriteOnce (RWO)
- 单节点读写
- 只能被一个节点挂载
- 适合单实例应用

### 2. ReadOnlyMany (ROX)
- 多节点只读
- 可以被多个节点同时挂载
- 适合共享配置文件

### 3. ReadWriteMany (RWM)
- 多节点读写
- 可以被多个节点同时读写
- 需要支持分布式文件系统

## 🛠️ 常用操作

### 1. 创建和管理存储
```bash
# 创建 PV
kubectl apply -f pv.yaml

# 创建 PVC
kubectl apply -f pvc.yaml

# 查看 PV 和 PVC
kubectl get pv
kubectl get pvc

# 查看详细信息
kubectl describe pv pv-example
kubectl describe pvc pvc-example
```

### 2. 存储类管理
```bash
# 查看 StorageClass
kubectl get storageclass

# 创建 StorageClass
kubectl apply -f storageclass.yaml

# 设置默认 StorageClass
kubectl patch storageclass fast-ssd -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

### 3. 存储清理
```bash
# 删除 PVC
kubectl delete pvc pvc-example

# 删除 PV
kubectl delete pv pv-example

# 清理 StorageClass
kubectl delete storageclass fast-ssd
```

## 🎯 存储最佳实践

### 1. 存储选择
- 根据应用需求选择合适的存储类型
- 考虑数据持久性要求
- 评估性能和成本

### 2. 容量规划
- 合理规划存储容量
- 监控存储使用情况
- 设置存储配额

### 3. 备份策略
- 定期备份重要数据
- 测试恢复流程
- 文档化备份策略

### 4. 安全考虑
- 控制存储访问权限
- 加密敏感数据
- 监控存储访问

## 🛠️ 实践练习

### 练习 1：基础存储
1. 创建 PV 和 PVC
2. 在 Pod 中使用存储
3. 测试数据持久性

### 练习 2：动态供应
1. 创建 StorageClass
2. 使用动态供应创建存储
3. 测试自动创建 PV

### 练习 3：存储迁移
1. 创建不同存储类型
2. 迁移数据
3. 测试数据完整性

## 📚 扩展阅读

- [Kubernetes 存储官方文档](https://kubernetes.io/docs/concepts/storage/)
- [持久化存储](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
- [存储类](https://kubernetes.io/docs/concepts/storage/storage-classes/)

## 🎯 下一步

掌握存储管理后，继续学习：
- [安全机制](./09-security/README.md)
- [监控与日志](./10-monitoring/README.md) 