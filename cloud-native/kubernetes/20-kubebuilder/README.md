# Kubebuilder 开发框架

## 📚 学习目标

通过本模块学习，您将掌握：
- Kubebuilder 框架的核心概念和使用方法
- 如何快速创建和开发 Kubernetes Controller
- 自定义资源定义 (CRD) 的设计和实现
- 测试和调试 Controller 的最佳实践
- 与 Operator SDK 的区别和选择

## 🎯 核心概念

### 1. 什么是 Kubebuilder？

Kubebuilder 是一个用于构建 Kubernetes API 的 SDK 框架，它简化了自定义资源和控制器的开发过程。

**主要特性：**
- **代码生成**: 自动生成 CRD、Client、DeepCopy 等代码
- **测试框架**: 内置测试工具和模拟环境
- **最佳实践**: 遵循 Kubernetes 社区最佳实践
- **简单易用**: 命令行工具简化开发流程

### 2. Kubebuilder vs Operator SDK

| 特性 | Kubebuilder | Operator SDK |
|------|-------------|--------------|
| 开发语言 | Go | Go/Ansible/Helm |
| 代码生成 | 强大 | 基础 |
| 测试支持 | 优秀 | 良好 |
| 学习曲线 | 中等 | 简单 |
| 社区支持 | 官方 | CNCF |

### 3. 项目结构

```
my-operator/
├── api/
│   └── v1/
│       ├── groupversion_info.go
│       ├── memcached_types.go
│       └── zz_generated.deepcopy.go
├── config/
│   ├── crd/
│   ├── rbac/
│   ├── manager/
│   └── samples/
├── controllers/
│   ├── memcached_controller.go
│   └── suite_test.go
├── main.go
├── Makefile
└── PROJECT
```

## 🚀 快速开始

### 1. 安装 Kubebuilder

```bash
# 下载并安装
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder
sudo mv kubebuilder /usr/local/bin/

# 验证安装
kubebuilder version
```

### 2. 创建新项目

```bash
# 初始化项目
kubebuilder init --domain example.com --repo github.com/example/memcached-operator

# 创建 API
kubebuilder create api --group cache --version v1 --kind Memcached --resource --controller

# 生成代码
make generate
```

### 3. 定义自定义资源

```go
// api/v1/memcached_types.go
package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MemcachedSpec 定义期望状态
type MemcachedSpec struct {
    // 副本数
    Size int32 `json:"size"`
    
    // 镜像
    Image string `json:"image,omitempty"`
    
    // 端口
    Port int32 `json:"port,omitempty"`
}

// MemcachedStatus 定义观察状态
type MemcachedStatus struct {
    // 节点列表
    Nodes []string `json:"nodes,omitempty"`
    
    // 就绪副本数
    ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced

// Memcached 是 Memcached 资源的 Schema
type Memcached struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   MemcachedSpec   `json:"spec,omitempty"`
    Status MemcachedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MemcachedList 包含 Memcached 项目列表
type MemcachedList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Memcached `json:"items"`
}

func init() {
    SchemeBuilder.Register(&Memcached{}, &MemcachedList{})
}
```

### 4. 实现控制器

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

    cachev1 "github.com/example/memcached-operator/api/v1"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MemcachedReconciler 协调 Memcached 对象
type MemcachedReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile 是主要的协调循环
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // 获取 Memcached 实例
    memcached := &cachev1.Memcached{}
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
    memcached.Status.ReadyReplicas = *deployment.Spec.Replicas
    memcached.Status.Nodes = []string{fmt.Sprintf("%s-pod", memcached.Name)}
    if err := r.Status().Update(ctx, memcached); err != nil {
        log.Error(err, "Failed to update Memcached status")
        return ctrl.Result{}, err
    }

    return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// deploymentForMemcached 创建 Memcached 的 Deployment
func (r *MemcachedReconciler) deploymentForMemcached(m *cachev1.Memcached) *appsv1.Deployment {
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
                            ContainerPort: m.Spec.Port,
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

// labelsForMemcached 返回 Memcached 的标签
func labelsForMemcached(name string) map[string]string {
    return map[string]string{"app": "memcached", "memcached_cr": name}
}

// SetupWithManager 设置控制器管理器
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&cachev1.Memcached{}).
        Owns(&appsv1.Deployment{}).
        Complete(r)
}
```

## 🧪 测试

### 1. 单元测试

```go
// controllers/suite_test.go
package controllers

import (
    "context"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "k8s.io/client-go/kubernetes/scheme"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/envtest"

    cachev1 "github.com/example/memcached-operator/api/v1"
)

var k8sClient client.Client
var testEnv *envtest.Environment

func TestControllers(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
    By("bootstrapping test environment")
    testEnv = &envtest.Environment{}

    cfg, err := testEnv.Start()
    Expect(err).NotTo(HaveOccurred())
    Expect(cfg).NotTo(BeNil())

    err = cachev1.AddToScheme(scheme.Scheme)
    Expect(err).NotTo(HaveOccurred())

    k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
    Expect(err).NotTo(HaveOccurred())
    Expect(k8sClient).NotTo(BeNil())
}, 60)

var _ = AfterSuite(func() {
    By("tearing down the test environment")
    err := testEnv.Stop()
    Expect(err).NotTo(HaveOccurred())
})
```

### 2. 集成测试

```go
// controllers/memcached_controller_test.go
package controllers

import (
    "context"
    "time"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    appsv1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"

    cachev1 "github.com/example/memcached-operator/api/v1"
)

var _ = Describe("Memcached Controller", func() {
    const (
        MemcachedName      = "test-memcached"
        MemcachedNamespace = "default"
        MemcachedSize      = 3
    )

    Context("When creating a Memcached resource", func() {
        It("Should create a Deployment with the correct size", func() {
            By("Creating a new Memcached resource")
            memcached := &cachev1.Memcached{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      MemcachedName,
                    Namespace: MemcachedNamespace,
                },
                Spec: cachev1.MemcachedSpec{
                    Size:  MemcachedSize,
                    Image: "memcached:1.6",
                    Port:  11211,
                },
            }
            Expect(k8sClient.Create(context.TODO(), memcached)).Should(Succeed())

            By("Checking if Deployment was created")
            deployment := &appsv1.Deployment{}
            Eventually(func() bool {
                err := k8sClient.Get(context.TODO(), types.NamespacedName{
                    Name:      MemcachedName,
                    Namespace: MemcachedNamespace,
                }, deployment)
                return err == nil
            }, time.Second*10, time.Millisecond*250).Should(BeTrue())

            By("Checking Deployment spec")
            Expect(*deployment.Spec.Replicas).Should(Equal(int32(MemcachedSize)))
        })
    })
})
```

## 🚀 部署和运行

### 1. 构建和部署

```bash
# 生成代码
make generate

# 生成 manifests
make manifests

# 构建镜像
make docker-build IMG=example.com/memcached-operator:v0.0.1

# 推送镜像
make docker-push IMG=example.com/memcached-operator:v0.0.1

# 部署 CRD
make install

# 部署控制器
make deploy IMG=example.com/memcached-operator:v0.0.1
```

### 2. 本地运行

```bash
# 本地运行控制器
make run

# 在另一个终端测试
kubectl apply -f config/samples/cache_v1_memcached.yaml
```

### 3. 清理

```bash
# 删除控制器
make undeploy

# 删除 CRD
make uninstall
```

## 🛠️ 实践练习

### 练习1: 创建 Web 应用控制器

```go
// 定义 WebApp 资源
type WebAppSpec struct {
    Replicas int32  `json:"replicas"`
    Image    string `json:"image"`
    Port     int32  `json:"port"`
    Domain   string `json:"domain,omitempty"`
}

type WebAppStatus struct {
    ReadyReplicas int32    `json:"readyReplicas,omitempty"`
    URL           string   `json:"url,omitempty"`
}
```

### 练习2: 实现数据库控制器

```go
// 定义 Database 资源
type DatabaseSpec struct {
    Type    string `json:"type"`
    Version string `json:"version"`
    Storage struct {
        Size string `json:"size"`
    } `json:"storage"`
    Backup struct {
        Enabled  bool   `json:"enabled"`
        Schedule string `json:"schedule,omitempty"`
    } `json:"backup"`
}
```

## 📚 相关资源

### 官方文档
- [Kubebuilder 官方文档](https://book.kubebuilder.io/)
- [Kubernetes API 约定](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)

### 学习资源
- [Kubebuilder 教程](https://book.kubebuilder.io/quick-start.html)
- [Controller Runtime 文档](https://pkg.go.dev/sigs.k8s.io/controller-runtime)

### 工具推荐
- **Kubebuilder**: 开发框架
- **envtest**: 测试环境
- **controller-gen**: 代码生成工具
- **kustomize**: 配置管理

---

**使用 Kubebuilder 快速构建强大的 Kubernetes 控制器！** 🚀
