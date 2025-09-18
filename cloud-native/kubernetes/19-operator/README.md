# Kubernetes Operator 模式

## 📚 学习目标

通过本模块学习，您将掌握：
- Operator 模式的核心概念和原理
- 如何设计和实现自定义 Operator
- 使用 Operator SDK 开发 Operator
- 生产环境中的 Operator 最佳实践
- 常见 Operator 的使用和管理

## 🎯 核心概念

### 1. 什么是 Operator？

Operator 是 Kubernetes 的扩展，它使用自定义资源来管理应用程序及其组件。Operator 遵循 Kubernetes 的理念，特别是控制器模式。

**Operator 的核心思想：**
- **声明式 API**: 描述期望状态
- **控制器模式**: 持续协调实际状态与期望状态
- **领域知识**: 将运维知识编码到软件中
- **自动化**: 减少人工干预

### 2. Operator 模式的优势

```
传统运维                    Operator 模式
┌─────────────┐            ┌─────────────────┐
│ 手动部署    │            │ 声明式部署      │
├─────────────┤            ├─────────────────┤
│ 手动配置    │            │ 自动化配置      │
├─────────────┤            ├─────────────────┤
│ 手动升级    │            │ 自动化升级      │
├─────────────┤            ├─────────────────┤
│ 手动备份    │            │ 自动化备份      │
├─────────────┤            ├─────────────────┤
│ 手动恢复    │            │ 自动化恢复      │
└─────────────┘            └─────────────────┘
```

### 3. Operator 架构

```
┌─────────────────────────────────────────────────┐
│                Kubernetes API                   │
├─────────────────────────────────────────────────┤
│  Custom Resource Definition (CRD)               │
│  ┌─────────────────────────────────────────┐   │
│  │  Custom Resource (CR)                   │   │
│  └─────────────────────────────────────────┘   │
├─────────────────────────────────────────────────┤
│                Operator                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────┐   │
│  │ Controller  │ │ Reconciler  │ │ Watcher │   │
│  └─────────────┘ └─────────────┘ └─────────┘   │
├─────────────────────────────────────────────────┤
│                Managed Resources                │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────┐   │
│  │ Pods    │ │ Services│ │ ConfigMaps│ PVCs│   │
│  └─────────┘ └─────────┘ └─────────┘ └─────┘   │
└─────────────────────────────────────────────────┘
```

## 🔧 Operator 开发

### 1. 使用 Operator SDK

```bash
# 安装 Operator SDK
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.28.0/operator-sdk_darwin_amd64
chmod +x operator-sdk_darwin_amd64
sudo mv operator-sdk_darwin_amd64 /usr/local/bin/operator-sdk

# 创建新项目
operator-sdk init --domain example.com --repo github.com/example/memcached-operator
cd memcached-operator

# 创建 API
operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller
```

### 2. 自定义资源定义 (CRD)

```yaml
# config/crd/bases/cache.example.com_memcacheds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: memcacheds.cache.example.com
spec:
  group: cache.example.com
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              size:
                type: integer
                minimum: 1
                maximum: 10
              image:
                type: string
                default: "memcached:1.6"
          status:
            type: object
            properties:
              nodes:
                type: array
                items:
                  type: string
  scope: Namespaced
  names:
    plural: memcacheds
    singular: memcached
    kind: Memcached
```

### 3. 控制器实现

```go
// controllers/memcached_controller.go
package controllers

import (
    "context"
    "fmt"
    "time"

    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MemcachedReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // 获取 Memcached 实例
    memcached := &cachev1alpha1.Memcached{}
    err := r.Get(ctx, req.NamespacedName, memcached)
    if err != nil {
        if errors.IsNotFound(err) {
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }

    // 检查 Deployment 是否存在
    deployment := &appsv1.Deployment{}
    err = r.Get(ctx, client.ObjectKey{
        Namespace: memcached.Namespace,
        Name:      memcached.Name,
    }, deployment)

    if err != nil && errors.IsNotFound(err) {
        // 创建 Deployment
        dep := r.deploymentForMemcached(memcached)
        if err := r.Create(ctx, dep); err != nil {
            log.Error(err, "Failed to create Deployment")
            return ctrl.Result{}, err
        }
        return ctrl.Result{Requeue: true}, nil
    } else if err != nil {
        log.Error(err, "Failed to get Deployment")
        return ctrl.Result{}, err
    }

    // 更新 Deployment 副本数
    size := memcached.Spec.Size
    if *deployment.Spec.Replicas != size {
        deployment.Spec.Replicas = &size
        if err := r.Update(ctx, deployment); err != nil {
            log.Error(err, "Failed to update Deployment")
            return ctrl.Result{}, err
        }
    }

    // 更新状态
    memcached.Status.Nodes = []string{fmt.Sprintf("%s-pod", memcached.Name)}
    if err := r.Status().Update(ctx, memcached); err != nil {
        log.Error(err, "Failed to update Memcached status")
        return ctrl.Result{}, err
    }

    return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *MemcachedReconciler) deploymentForMemcached(m *cachev1alpha1.Memcached) *appsv1.Deployment {
    ls := labelsForMemcached(m.Name)
    replicas := m.Spec.Size

    dep := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      m.Name,
            Namespace: m.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: ls,
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: ls,
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{{
                        Image: m.Spec.Image,
                        Name:  "memcached",
                        Ports: []corev1.ContainerPort{{
                            ContainerPort: 11211,
                            Name:          "memcached",
                        }},
                    }},
                },
            },
        },
    }
    ctrl.SetControllerReference(m, dep, r.Scheme)
    return dep
}

func labelsForMemcached(name string) map[string]string {
    return map[string]string{"app": "memcached", "memcached_cr": name}
}

func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&cachev1alpha1.Memcached{}).
        Owns(&appsv1.Deployment{}).
        Complete(r)
}
```

## 🚀 部署和测试

### 1. 构建和部署 Operator

```bash
# 构建镜像
make docker-build docker-push IMG=example.com/memcached-operator:v0.0.1

# 部署 CRD
make install

# 部署 Operator
make deploy IMG=example.com/memcached-operator:v0.0.1
```

### 2. 创建自定义资源

```yaml
# config/samples/cache_v1alpha1_memcached.yaml
apiVersion: cache.example.com/v1alpha1
kind: Memcached
metadata:
  name: memcached-sample
spec:
  size: 3
  image: "memcached:1.6"
```

### 3. 测试 Operator

```bash
# 应用示例
kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml

# 检查状态
kubectl get memcacheds
kubectl describe memcached memcached-sample

# 检查创建的 Deployment
kubectl get deployments
kubectl get pods
```

## 🌟 生产环境最佳实践

### 1. 错误处理和重试

```go
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 设置超时
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // 指数退避重试
    backoff := wait.Backoff{
        Duration: 1 * time.Second,
        Factor:   2.0,
        Steps:    5,
    }

    return wait.ExponentialBackoff(backoff, func() (bool, error) {
        // 执行协调逻辑
        return true, nil
    })
}
```

### 2. 事件记录

```go
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 记录事件
    r.Recorder.Event(memcached, corev1.EventTypeNormal, "Reconciled", "Memcached reconciled successfully")
    
    // 记录错误事件
    if err != nil {
        r.Recorder.Event(memcached, corev1.EventTypeWarning, "Error", err.Error())
    }
}
```

### 3. 健康检查

```go
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&cachev1alpha1.Memcached{}).
        WithOptions(controller.Options{
            MaxConcurrentReconciles: 1,
        }).
        Complete(r)
}
```

## 📦 常见 Operator 示例

### 1. Prometheus Operator

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
spec:
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      team: frontend
  ruleSelector:
    matchLabels:
      prometheus: k8s
      role: alert-rules
  resources:
    requests:
      memory: 400Mi
```

### 2. Elasticsearch Operator

```yaml
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: quickstart
spec:
  version: 8.8.0
  nodeSets:
  - name: default
    count: 1
    config:
      node.roles: ["master", "data", "ingest"]
```

### 3. MySQL Operator

```yaml
apiVersion: mysql.oracle.com/v2
kind: InnoDBCluster
metadata:
  name: mycluster
spec:
  secretName: mycluster-secret
  tlsUseSelfSigned: true
  instances: 3
  router:
    instances: 2
```

## 🛠️ 实践练习

### 练习1: 创建简单的 Web 应用 Operator

```yaml
apiVersion: apps.example.com/v1alpha1
kind: WebApp
metadata:
  name: my-webapp
spec:
  replicas: 3
  image: nginx:alpine
  port: 80
  domain: example.com
```

### 练习2: 实现数据库 Operator

```yaml
apiVersion: database.example.com/v1alpha1
kind: Database
metadata:
  name: my-db
spec:
  type: postgresql
  version: "13"
  storage:
    size: 10Gi
  backup:
    enabled: true
    schedule: "0 2 * * *"
```

## 📚 相关资源

### 官方文档
- [Operator SDK 文档](https://sdk.operatorframework.io/)
- [Kubernetes Operator 模式](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

### 学习资源
- [Operator 最佳实践](https://sdk.operatorframework.io/docs/best-practices/)
- [Operator 生命周期管理](https://olm.operatorframework.io/)

### 工具推荐
- **Operator SDK**: Operator 开发框架
- **Operator Lifecycle Manager**: Operator 生命周期管理
- **Kubebuilder**: 另一个 Operator 开发框架
- **Helm Operator**: 基于 Helm 的 Operator

---

**掌握 Operator 模式，实现 Kubernetes 的无限扩展！** 🚀
