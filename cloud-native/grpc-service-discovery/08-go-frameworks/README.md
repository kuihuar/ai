# Go 框架服务注册与发现集成

## 📖 概述

本文档介绍如何将 gRPC 服务注册与发现集成到各种 Go Web 框架中，专注于核心的服务发现功能。

## �� 支持的框架

- **Web 框架**: Gin, Echo, Fiber, Chi
- **微服务框架**: Go-kit, Go-micro, Kratos

## 🚀 核心组件

### 1. 服务发现接口

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
    // Discover 发现服务实例
    Discover(serviceName string) ([]string, error)
    // Register 注册服务实例
    Register(serviceName, address string) error
    // Deregister 注销服务实例
    Deregister(serviceName, address string) error
    // Watch 监听服务变化
    Watch(serviceName string) (<-chan []string, error)
}

// ServiceInstance 服务实例
type ServiceInstance struct {
    ID      string            `json:"id"`
    Name    string            `json:"name"`
    Address string            `json:"address"`
    Port    int               `json:"port"`
    Tags    []string          `json:"tags"`
    Meta    map[string]string `json:"meta"`
}

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
    // Register 注册服务
    Register(instance *ServiceInstance) error
    // Deregister 注销服务
    Deregister(instanceID string) error
    // Discover 发现服务
    Discover(serviceName string) ([]*ServiceInstance, error)
    // Watch 监听服务变化
    Watch(serviceName string) (<-chan []*ServiceInstance, error)
}
```

### 2. Consul 服务发现实现

```go
package main

import (
    "fmt"
    "time"

    "github.com/hashicorp/consul/api"
)

// ConsulServiceDiscovery Consul 服务发现
type ConsulServiceDiscovery struct {
    client *api.Client
}

// NewConsulServiceDiscovery 创建 Consul 服务发现
func NewConsulServiceDiscovery(address string) (*ConsulServiceDiscovery, error) {
    config := api.DefaultConfig()
    config.Address = address
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ConsulServiceDiscovery{client: client}, nil
}

// Discover 发现服务实例
func (csd *ConsulServiceDiscovery) Discover(serviceName string) ([]string, error) {
    services, _, err := csd.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, service := range services {
        address := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
        addresses = append(addresses, address)
    }
    
    return addresses, nil
}

// Register 注册服务实例
func (csd *ConsulServiceDiscovery) Register(serviceName, address string) error {
    registration := &api.AgentServiceRegistration{
        Name: serviceName,
        Address: address,
        Port: 8080, // 默认端口
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s/health", address),
            Interval:                       "10s",
            Timeout:                        "3s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    return csd.client.Agent().ServiceRegister(registration)
}

// Deregister 注销服务实例
func (csd *ConsulServiceDiscovery) Deregister(serviceName, address string) error {
    return csd.client.Agent().ServiceDeregister(serviceName)
}

// Watch 监听服务变化
func (csd *ConsulServiceDiscovery) Watch(serviceName string) (<-chan []string, error) {
    ch := make(chan []string, 1)
    
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                addresses, err := csd.Discover(serviceName)
                if err == nil {
                    ch <- addresses
                }
            }
        }
    }()
    
    return ch, nil
}
```

### 3. Etcd 服务发现实现

```go
package main

import (
    "context"
    "fmt"
    "strings"
    "time"

    "go.etcd.io/etcd/clientv3"
)

// EtcdServiceDiscovery Etcd 服务发现
type EtcdServiceDiscovery struct {
    client *clientv3.Client
}

// NewEtcdServiceDiscovery 创建 Etcd 服务发现
func NewEtcdServiceDiscovery(endpoints []string) (*EtcdServiceDiscovery, error) {
    client, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        return nil, err
    }
    
    return &EtcdServiceDiscovery{client: client}, nil
}

// Discover 发现服务实例
func (esd *EtcdServiceDiscovery) Discover(serviceName string) ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    resp, err := esd.client.Get(ctx, fmt.Sprintf("/services/%s/", serviceName), clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, kv := range resp.Kvs {
        key := string(kv.Key)
        if strings.HasSuffix(key, "/address") {
            address := string(kv.Value)
            addresses = append(addresses, address)
        }
    }
    
    return addresses, nil
}

// Register 注册服务实例
func (esd *EtcdServiceDiscovery) Register(serviceName, address string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    key := fmt.Sprintf("/services/%s/%s/address", serviceName, address)
    _, err := esd.client.Put(ctx, key, address)
    return err
}

// Deregister 注销服务实例
func (esd *EtcdServiceDiscovery) Deregister(serviceName, address string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    key := fmt.Sprintf("/services/%s/%s/", serviceName, address)
    _, err := esd.client.Delete(ctx, key, clientv3.WithPrefix())
    return err
}

// Watch 监听服务变化
func (esd *EtcdServiceDiscovery) Watch(serviceName string) (<-chan []string, error) {
    ch := make(chan []string, 1)
    
    go func() {
        watchCh := esd.client.Watch(context.Background(), fmt.Sprintf("/services/%s/", serviceName), clientv3.WithPrefix())
        
        for watchResp := range watchCh {
            if watchResp.Err() != nil {
                continue
            }
            
            addresses, err := esd.Discover(serviceName)
            if err == nil {
                ch <- addresses
            }
        }
    }()
    
    return ch, nil
}
```

## 🔧 框架集成

### 1. Gin 框架集成

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    pb "your-project/proto"
)

// GinServiceDiscovery Gin 服务发现
type GinServiceDiscovery struct {
    discovery ServiceDiscovery
    clients   map[string]*grpc.ClientConn
}

// NewGinServiceDiscovery 创建 Gin 服务发现
func NewGinServiceDiscovery(discovery ServiceDiscovery) *GinServiceDiscovery {
    return &GinServiceDiscovery{
        discovery: discovery,
        clients:   make(map[string]*grpc.ClientConn),
    }
}

// GetServiceClient 获取服务客户端
func (gsd *GinServiceDiscovery) GetServiceClient(serviceName string) (pb.YourServiceClient, error) {
    if client, ok := gsd.clients[serviceName]; ok {
        return pb.NewYourServiceClient(client), nil
    }
    
    addresses, err := gsd.discovery.Discover(serviceName)
    if err != nil {
        return nil, err
    }
    
    if len(addresses) == 0 {
        return nil, fmt.Errorf("no instances found for service %s", serviceName)
    }
    
    // 选择第一个地址（实际应用中可以使用负载均衡）
    conn, err := grpc.Dial(addresses[0], grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    gsd.clients[serviceName] = conn
    return pb.NewYourServiceClient(conn), nil
}

// RegisterService 注册服务
func (gsd *GinServiceDiscovery) RegisterService(serviceName, address string) error {
    return gsd.discovery.Register(serviceName, address)
}

// DeregisterService 注销服务
func (gsd *GinServiceDiscovery) DeregisterService(serviceName, address string) error {
    return gsd.discovery.Deregister(serviceName, address)
}

// 使用示例
func main() {
    // 创建服务发现
    discovery, err := NewConsulServiceDiscovery("localhost:8500")
    if err != nil {
        panic(err)
    }
    
    // 创建 Gin 服务发现
    ginSD := NewGinServiceDiscovery(discovery)
    
    // 注册当前服务
    err = ginSD.RegisterService("user-service", "localhost:8080")
    if err != nil {
        panic(err)
    }
    
    // 创建 Gin 路由
    r := gin.Default()
    
    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    
    // 服务发现测试
    r.GET("/discover/:service", func(c *gin.Context) {
        serviceName := c.Param("service")
        addresses, err := discovery.Discover(serviceName)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"service": serviceName, "addresses": addresses})
    })
    
    // 启动服务器
    r.Run(":8080")
}
```

### 2. Echo 框架集成

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "google.golang.org/grpc"
    pb "your-project/proto"
)

// EchoServiceDiscovery Echo 服务发现
type EchoServiceDiscovery struct {
    discovery ServiceDiscovery
    clients   map[string]*grpc.ClientConn
}

// NewEchoServiceDiscovery 创建 Echo 服务发现
func NewEchoServiceDiscovery(discovery ServiceDiscovery) *EchoServiceDiscovery {
    return &EchoServiceDiscovery{
        discovery: discovery,
        clients:   make(map[string]*grpc.ClientConn),
    }
}

// GetServiceClient 获取服务客户端
func (esd *EchoServiceDiscovery) GetServiceClient(serviceName string) (pb.YourServiceClient, error) {
    if client, ok := esd.clients[serviceName]; ok {
        return pb.NewYourServiceClient(client), nil
    }
    
    addresses, err := esd.discovery.Discover(serviceName)
    if err != nil {
        return nil, err
    }
    
    if len(addresses) == 0 {
        return nil, fmt.Errorf("no instances found for service %s", serviceName)
    }
    
    conn, err := grpc.Dial(addresses[0], grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    esd.clients[serviceName] = conn
    return pb.NewYourServiceClient(conn), nil
}

// RegisterService 注册服务
func (esd *EchoServiceDiscovery) RegisterService(serviceName, address string) error {
    return esd.discovery.Register(serviceName, address)
}

// DeregisterService 注销服务
func (esd *EchoServiceDiscovery) DeregisterService(serviceName, address string) error {
    return esd.discovery.Deregister(serviceName, address)
}

// 使用示例
func main() {
    // 创建服务发现
    discovery, err := NewConsulServiceDiscovery("localhost:8500")
    if err != nil {
        panic(err)
    }
    
    // 创建 Echo 服务发现
    echoSD := NewEchoServiceDiscovery(discovery)
    
    // 注册当前服务
    err = echoSD.RegisterService("user-service", "localhost:8080")
    if err != nil {
        panic(err)
    }
    
    // 创建 Echo 实例
    e := echo.New()
    
    // 添加中间件
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    // 健康检查
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
    })
    
    // 服务发现测试
    e.GET("/discover/:service", func(c echo.Context) error {
        serviceName := c.Param("service")
        addresses, err := discovery.Discover(serviceName)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
        }
        return c.JSON(http.StatusOK, map[string]interface{}{
            "service":   serviceName,
            "addresses": addresses,
        })
    })
    
    // 启动服务器
    e.Logger.Fatal(e.Start(":8080"))
}
```

### 3. Fiber 框架集成

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "google.golang.org/grpc"
    pb "your-project/proto"
)

// FiberServiceDiscovery Fiber 服务发现
type FiberServiceDiscovery struct {
    discovery ServiceDiscovery
    clients   map[string]*grpc.ClientConn
}

// NewFiberServiceDiscovery 创建 Fiber 服务发现
func NewFiberServiceDiscovery(discovery ServiceDiscovery) *FiberServiceDiscovery {
    return &FiberServiceDiscovery{
        discovery: discovery,
        clients:   make(map[string]*grpc.ClientConn),
    }
}

// GetServiceClient 获取服务客户端
func (fsd *FiberServiceDiscovery) GetServiceClient(serviceName string) (pb.YourServiceClient, error) {
    if client, ok := fsd.clients[serviceName]; ok {
        return pb.NewYourServiceClient(client), nil
    }
    
    addresses, err := fsd.discovery.Discover(serviceName)
    if err != nil {
        return nil, err
    }
    
    if len(addresses) == 0 {
        return nil, fmt.Errorf("no instances found for service %s", serviceName)
    }
    
    conn, err := grpc.Dial(addresses[0], grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    fsd.clients[serviceName] = conn
    return pb.NewYourServiceClient(conn), nil
}

// RegisterService 注册服务
func (fsd *FiberServiceDiscovery) RegisterService(serviceName, address string) error {
    return fsd.discovery.Register(serviceName, address)
}

// DeregisterService 注销服务
func (fsd *FiberServiceDiscovery) DeregisterService(serviceName, address string) error {
    return fsd.discovery.Deregister(serviceName, address)
}

// 使用示例
func main() {
    // 创建服务发现
    discovery, err := NewConsulServiceDiscovery("localhost:8500")
    if err != nil {
        panic(err)
    }
    
    // 创建 Fiber 服务发现
    fiberSD := NewFiberServiceDiscovery(discovery)
    
    // 注册当前服务
    err = fiberSD.RegisterService("user-service", "localhost:8080")
    if err != nil {
        panic(err)
    }
    
    // 创建 Fiber 应用
    app := fiber.New()
    
    // 添加中间件
    app.Use(logger.New())
    app.Use(recover.New())
    
    // 健康检查
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })
    
    // 服务发现测试
    app.Get("/discover/:service", func(c *fiber.Ctx) error {
        serviceName := c.Params("service")
        addresses, err := discovery.Discover(serviceName)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{
            "service":   serviceName,
            "addresses": addresses,
        })
    })
    
    // 启动服务器
    app.Listen(":8080")
}
```

## 📚 使用说明

1. **选择服务发现后端**: 支持 Consul、Etcd 等
2. **集成到框架**: 根据使用的 Web 框架选择对应的集成方案
3. **注册服务**: 在服务启动时注册到服务发现
4. **发现服务**: 在需要调用其他服务时通过服务发现获取地址
5. **注销服务**: 在服务关闭时注销服务

## 🔗 相关资源

- [Consul 文档](https://www.consul.io/docs)
- [Etcd 文档](https://etcd.io/docs/)
- [gRPC Go 文档](https://grpc.io/docs/languages/go/)
