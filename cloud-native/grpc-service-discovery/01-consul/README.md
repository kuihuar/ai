# Consul 服务注册与发现

## 📖 概述

Consul 是 HashiCorp 开源的服务发现和配置管理工具，提供了完整的服务注册与发现功能，支持多数据中心、健康检查、键值存储等特性。

## 🎯 核心特性

### 1. 服务注册与发现
- 自动服务注册和注销
- 基于 DNS 和 HTTP 的服务发现
- 支持多数据中心部署

### 2. 健康检查
- 支持多种健康检查方式
- 自动故障检测和恢复
- 可配置的检查间隔和超时

### 3. 键值存储
- 分布式键值存储
- 支持事务操作
- 事件通知机制

## 🚀 快速开始

### 1. 安装 Consul

```bash
# 下载 Consul
wget https://releases.hashicorp.com/consul/1.15.2/consul_1.15.2_linux_amd64.zip
unzip consul_1.15.2_linux_amd64.zip
sudo mv consul /usr/local/bin/

# 验证安装
consul version
```

### 2. 启动 Consul 服务器

```bash
# 开发模式启动
consul agent -dev -ui -client=0.0.0.0

# 生产模式启动
consul agent -server -bootstrap-expect=3 -data-dir=/tmp/consul -node=consul-1 -bind=0.0.0.0 -client=0.0.0.0 -ui
```

### 3. Go 客户端实现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"

    "github.com/hashicorp/consul/api"
    "google.golang.org/grpc"
)

// ServiceRegistry Consul 服务注册器
type ServiceRegistry struct {
    client *api.Client
}

// NewServiceRegistry 创建服务注册器
func NewServiceRegistry(consulAddr string) (*ServiceRegistry, error) {
    config := api.DefaultConfig()
    config.Address = consulAddr
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    return &ServiceRegistry{client: client}, nil
}

// Register 注册服务
func (sr *ServiceRegistry) Register(serviceName, serviceID, address string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    serviceName,
        Address: address,
        Port:    port,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
            Timeout:                        "3s",
            Interval:                       "10s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    return sr.client.Agent().ServiceRegister(registration)
}

// Deregister 注销服务
func (sr *ServiceRegistry) Deregister(serviceID string) error {
    return sr.client.Agent().ServiceDeregister(serviceID)
}

// Discover 发现服务
func (sr *ServiceRegistry) Discover(serviceName string) ([]string, error) {
    services, _, err := sr.client.Health().Service(serviceName, "", false, nil)
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

// Watch 监听服务变化
func (sr *ServiceRegistry) Watch(serviceName string, callback func([]string)) {
    go func() {
        for {
            services, _, err := sr.client.Health().Service(serviceName, "", false, nil)
            if err != nil {
                log.Printf("Error watching service %s: %v", serviceName, err)
                time.Sleep(5 * time.Second)
                continue
            }
            
            var addresses []string
            for _, service := range services {
                address := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
                addresses = append(addresses, address)
            }
            
            callback(addresses)
            time.Sleep(10 * time.Second)
        }
    }()
}
```

## 📝 使用示例

### 1. 服务端注册

```go
package main

import (
    "context"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"

    "google.golang.org/grpc"
    pb "your-project/proto"
)

func main() {
    // 创建 gRPC 服务器
    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    pb.RegisterYourServiceServer(s, &server{})
    
    // 创建服务注册器
    registry, err := NewServiceRegistry("localhost:8500")
    if err != nil {
        log.Fatalf("failed to create registry: %v", err)
    }
    
    // 注册服务
    serviceID := fmt.Sprintf("%s-%s", "your-service", getLocalIP())
    err = registry.Register("your-service", serviceID, getLocalIP(), 8080)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
    
    // 优雅关闭
    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c
        
        // 注销服务
        registry.Deregister(serviceID)
        s.GracefulStop()
    }()
    
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

func getLocalIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "127.0.0.1"
    }
    defer conn.Close()
    return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
```

### 2. 客户端发现

```go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "your-project/proto"
)

func main() {
    // 创建服务注册器
    registry, err := NewServiceRegistry("localhost:8500")
    if err != nil {
        log.Fatalf("failed to create registry: %v", err)
    }
    
    // 发现服务
    addresses, err := registry.Discover("your-service")
    if err != nil {
        log.Fatalf("failed to discover service: %v", err)
    }
    
    if len(addresses) == 0 {
        log.Fatalf("no service instances found")
    }
    
    // 连接服务
    conn, err := grpc.Dial(addresses[0], grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()
    
    client := pb.NewYourServiceClient(conn)
    
    // 调用服务
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    resp, err := client.YourMethod(ctx, &pb.YourRequest{})
    if err != nil {
        log.Fatalf("failed to call service: %v", err)
    }
    
    log.Printf("Response: %v", resp)
}
```

## 🔧 高级配置

### 1. 健康检查配置

```go
// 自定义健康检查
registration := &api.AgentServiceRegistration{
    ID:      serviceID,
    Name:    serviceName,
    Address: address,
    Port:    port,
    Check: &api.AgentServiceCheck{
        // HTTP 健康检查
        HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
        Timeout:                        "3s",
        Interval:                       "10s",
        DeregisterCriticalServiceAfter: "30s",
        
        // 或者使用 gRPC 健康检查
        // GRPC:                          fmt.Sprintf("%s:%d", address, port),
        // GRPCUseTLS:                    false,
    },
}
```

### 2. 多数据中心配置

```go
// 多数据中心客户端
config := api.DefaultConfig()
config.Address = "consul-dc1:8500"
config.Datacenter = "dc1"

client, err := api.NewClient(config)
if err != nil {
    log.Fatalf("failed to create consul client: %v", err)
}
```

### 3. 服务标签和元数据

```go
registration := &api.AgentServiceRegistration{
    ID:      serviceID,
    Name:    serviceName,
    Address: address,
    Port:    port,
    Tags:    []string{"v1", "production", "web"},
    Meta: map[string]string{
        "version": "1.0.0",
        "region":  "us-west-2",
        "env":     "production",
    },
    Check: &api.AgentServiceCheck{
        HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
        Timeout:                        "3s",
        Interval:                       "10s",
        DeregisterCriticalServiceAfter: "30s",
    },
}
```

## 📊 性能优化

### 1. 连接池配置

```go
config := api.DefaultConfig()
config.Address = "localhost:8500"
config.Transport = &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
}
```

### 2. 缓存配置

```go
// 使用缓存减少 API 调用
services, _, err := sr.client.Health().Service(serviceName, "", false, &api.QueryOptions{
    UseCache: true,
    MaxAge:   30 * time.Second,
})
```

## 🛡️ 安全配置

### 1. ACL 配置

```go
config := api.DefaultConfig()
config.Address = "localhost:8500"
config.Token = "your-acl-token"
```

### 2. TLS 配置

```go
config := api.DefaultConfig()
config.Address = "https://consul.example.com:8501"
config.TLSConfig = api.TLSConfig{
    Address: "consul.example.com",
    CAFile:  "/path/to/ca.pem",
    CertFile: "/path/to/cert.pem",
    KeyFile:  "/path/to/key.pem",
}
```

## 🔍 监控和调试

### 1. 日志配置

```go
config := api.DefaultConfig()
config.Address = "localhost:8500"
config.LogLevel = "DEBUG"
```

### 2. 指标监控

```go
// 使用 Prometheus 监控
import "github.com/prometheus/client_golang/prometheus"

var (
    serviceRegistrations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "consul_service_registrations_total",
            Help: "Total number of service registrations",
        },
        []string{"service_name", "status"},
    )
)

func init() {
    prometheus.MustRegister(serviceRegistrations)
}
```

## 📚 最佳实践

1. **服务命名**: 使用有意义的服务名称，遵循命名规范
2. **健康检查**: 实现快速响应的健康检查端点
3. **优雅关闭**: 确保服务关闭时正确注销
4. **监控告警**: 监控服务注册状态和健康检查结果
5. **多数据中心**: 在多个数据中心部署 Consul 集群

## 🔗 相关资源

- [Consul 官方文档](https://www.consul.io/docs)
- [Consul Go API 文档](https://pkg.go.dev/github.com/hashicorp/consul/api)
- [Consul 最佳实践](https://www.consul.io/docs/guides)
