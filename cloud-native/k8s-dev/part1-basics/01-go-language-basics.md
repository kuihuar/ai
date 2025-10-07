# 1. Go 语言基础

## 📖 概述

Go 语言是 Kubernetes 生态系统的核心开发语言。掌握 Go 语言对于进行 Kubernetes 二次开发至关重要。

## 🎯 学习目标

- 掌握 Go 语言基本语法和特性
- 理解 Go 的并发模型（Goroutines 和 Channels）
- 熟悉 Go 的包管理和模块系统
- 掌握 Go 的测试框架
- 了解 Go 的性能优化技巧

## 📚 核心内容

### 1.1 基础语法

#### 变量和常量
```go
// 变量声明
var name string = "Kubernetes"
var age int = 10
var isActive bool = true

// 短变量声明
name := "Kubernetes"
age := 10

// 常量
const (
    APIVersion = "v1"
    Kind       = "Pod"
)
```

#### 数据类型
```go
// 基本类型
var i int = 42
var f float64 = 3.14
var s string = "Hello"
var b bool = true

// 复合类型
type PodSpec struct {
    Containers []Container `json:"containers"`
    RestartPolicy string   `json:"restartPolicy"`
}

// 接口
type Object interface {
    GetName() string
    GetNamespace() string
}
```

### 1.2 函数和方法

```go
// 函数
func CreatePod(name, namespace string) (*Pod, error) {
    return &Pod{
        Name:      name,
        Namespace: namespace,
    }, nil
}

// 方法
func (p *Pod) GetFullName() string {
    return fmt.Sprintf("%s/%s", p.Namespace, p.Name)
}

// 可变参数
func LogInfo(format string, args ...interface{}) {
    log.Printf(format, args...)
}
```

### 1.3 并发编程

#### Goroutines
```go
func main() {
    // 启动 goroutine
    go processPod("pod-1")
    go processPod("pod-2")
    
    // 等待完成
    time.Sleep(2 * time.Second)
}

func processPod(name string) {
    fmt.Printf("Processing pod: %s\n", name)
    time.Sleep(1 * time.Second)
}
```

#### Channels
```go
func main() {
    // 创建 channel
    ch := make(chan string, 2)
    
    // 发送数据
    go func() {
        ch <- "pod-1"
        ch <- "pod-2"
        close(ch)
    }()
    
    // 接收数据
    for pod := range ch {
        fmt.Printf("Received: %s\n", pod)
    }
}
```

### 1.4 错误处理

```go
func GetPod(client kubernetes.Interface, name, namespace string) (*v1.Pod, error) {
    pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get pod %s/%s: %w", namespace, name, err)
    }
    return pod, nil
}

// 自定义错误类型
type PodNotFoundError struct {
    Name      string
    Namespace string
}

func (e *PodNotFoundError) Error() string {
    return fmt.Sprintf("pod %s not found in namespace %s", e.Name, e.Namespace)
}
```

### 1.5 包管理

#### Go Modules
```go
// go.mod
module github.com/example/k8s-operator

go 1.21

require (
    k8s.io/api v0.28.0
    k8s.io/client-go v0.28.0
    sigs.k8s.io/controller-runtime v0.16.0
)
```

#### 导入包
```go
import (
    "context"
    "fmt"
    "time"
    
    "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)
```

### 1.6 测试

#### 单元测试
```go
func TestCreatePod(t *testing.T) {
    tests := []struct {
        name      string
        namespace string
        want      string
        wantErr   bool
    }{
        {
            name:      "valid pod",
            namespace: "default",
            want:      "default/pod-1",
            wantErr:   false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pod, err := CreatePod(tt.name, tt.namespace)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreatePod() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if pod.GetFullName() != tt.want {
                t.Errorf("CreatePod() = %v, want %v", pod.GetFullName(), tt.want)
            }
        })
    }
}
```

#### 集成测试
```go
func TestPodIntegration(t *testing.T) {
    // 使用 envtest 进行集成测试
    testEnv := &envtest.Environment{
        CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
    }
    
    cfg, err := testEnv.Start()
    if err != nil {
        t.Fatal(err)
    }
    defer testEnv.Stop()
    
    // 测试逻辑...
}
```

## 🛠️ 开发工具

### 必需工具
```bash
# 安装 Go
brew install go  # macOS
# 或下载官方安装包

# 验证安装
go version

# 设置环境变量
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### 推荐工具
```bash
# 代码格式化
go install golang.org/x/tools/cmd/goimports@latest

# 静态分析
go install honnef.co/go/tools/cmd/staticcheck@latest

# 测试覆盖率
go test -cover ./...

# 性能分析
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## 📖 学习资源

### 官方文档
- [Go 官方文档](https://golang.org/doc/)
- [Go 语言规范](https://golang.org/ref/spec)
- [Effective Go](https://golang.org/doc/effective_go.html)

### 推荐书籍
- 《Go 语言圣经》
- 《Go 并发编程实战》
- 《Go 语言高级编程》

### 在线资源
- [Go by Example](https://gobyexample.com/)
- [Go Playground](https://play.golang.org/)
- [Go 语言中文网](https://studygolang.com/)

## 🎯 实践练习

### 练习 1：基础语法
创建一个简单的 Pod 管理程序，包含：
- Pod 结构体定义
- 创建、删除、查询 Pod 的方法
- 错误处理

### 练习 2：并发编程
实现一个并发处理多个 Pod 的程序：
- 使用 goroutines 并发处理
- 使用 channels 进行通信
- 实现优雅关闭

### 练习 3：测试
为你的 Pod 管理程序编写：
- 单元测试
- 集成测试
- 性能测试

## 🔗 下一步

掌握 Go 语言基础后，建议继续学习：
- [Kubernetes 核心架构](./02-kubernetes-core-architecture.md)
- [开发环境搭建](./03-development-environment.md)
