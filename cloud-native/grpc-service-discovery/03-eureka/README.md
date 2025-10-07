# Eureka 服务注册与发现

## 📖 概述

Eureka 是 Netflix 开源的服务发现组件，是 Spring Cloud 生态系统的核心组件之一。虽然 Netflix 已经停止维护，但在很多遗留系统中仍在使用。

## 🎯 核心特性

### 1. 服务注册与发现
- 基于 REST API 的服务注册
- 客户端服务发现
- 自我保护机制

### 2. 高可用性
- 支持集群部署
- 自动故障转移
- 数据同步

### 3. 简单易用
- RESTful API
- 配置简单
- 集成方便

## 🚀 快速开始

### 1. 启动 Eureka 服务器

```bash
# 下载 Eureka 服务器 JAR
wget https://repo1.maven.org/maven2/com/netflix/eureka/eureka-server/1.10.17/eureka-server-1.10.17.war

# 启动 Eureka 服务器
java -jar eureka-server-1.10.17.war --server.port=8761
```

### 2. Go 客户端实现

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// EurekaRegistry Eureka 服务注册器
type EurekaRegistry struct {
    eurekaURL string
    client    *http.Client
}

// NewEurekaRegistry 创建 Eureka 服务注册器
func NewEurekaRegistry(eurekaURL string) *EurekaRegistry {
    return &EurekaRegistry{
        eurekaURL: eurekaURL,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// InstanceInfo 实例信息
type InstanceInfo struct {
    InstanceID string            `json:"instanceId"`
    HostName   string            `json:"hostName"`
    App        string            `json:"app"`
    IPAddr     string            `json:"ipAddr"`
    Status     string            `json:"status"`
    Port       Port              `json:"port"`
    DataCenter DataCenter        `json:"dataCenterInfo"`
    Lease      Lease             `json:"leaseInfo"`
    Metadata   map[string]string `json:"metadata"`
}

type Port struct {
    Port    int  `json:"$"`
    Enabled bool `json:"@enabled"`
}

type DataCenter struct {
    Name  string `json:"name"`
    Class string `json:"@class"`
}

type Lease struct {
    RenewalIntervalInSecs int `json:"renewalIntervalInSecs"`
    DurationInSecs        int `json:"durationInSecs"`
}

// Register 注册服务
func (er *EurekaRegistry) Register(instanceInfo *InstanceInfo) error {
    url := fmt.Sprintf("%s/apps/%s", er.eurekaURL, instanceInfo.App)
    
    data := map[string]interface{}{
        "instance": instanceInfo,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := er.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to register service: %d", resp.StatusCode)
    }
    
    return nil
}

// Deregister 注销服务
func (er *EurekaRegistry) Deregister(app, instanceID string) error {
    url := fmt.Sprintf("%s/apps/%s/%s", er.eurekaURL, app, instanceID)
    
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }
    
    resp, err := er.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to deregister service: %d", resp.StatusCode)
    }
    
    return nil
}

// Discover 发现服务
func (er *EurekaRegistry) Discover(app string) ([]*InstanceInfo, error) {
    url := fmt.Sprintf("%s/apps/%s", er.eurekaURL, app)
    
    resp, err := er.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to discover service: %d", resp.StatusCode)
    }
    
    var result struct {
        Application struct {
            Instance []*InstanceInfo `json:"instance"`
        } `json:"application"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Application.Instance, nil
}

// Heartbeat 发送心跳
func (er *EurekaRegistry) Heartbeat(app, instanceID string) error {
    url := fmt.Sprintf("%s/apps/%s/%s", er.eurekaURL, app, instanceID)
    
    req, err := http.NewRequest("PUT", url, nil)
    if err != nil {
        return err
    }
    
    resp, err := er.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send heartbeat: %d", resp.StatusCode)
    }
    
    return nil
}
```

## 📝 使用示例

### 1. 服务端注册

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

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
    registry := NewEurekaRegistry("http://localhost:8761/eureka")
    
    // 注册服务
    instanceInfo := &InstanceInfo{
        InstanceID: fmt.Sprintf("%s:%s:%d", "your-service", getLocalIP(), 8080),
        HostName:   getLocalIP(),
        App:        "your-service",
        IPAddr:     getLocalIP(),
        Status:     "UP",
        Port: Port{
            Port:    8080,
            Enabled: true,
        },
        DataCenter: DataCenter{
            Name:  "MyOwn",
            Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
        },
        Lease: Lease{
            RenewalIntervalInSecs: 30,
            DurationInSecs:        90,
        },
        Metadata: map[string]string{
            "version": "1.0.0",
        },
    }
    
    err = registry.Register(instanceInfo)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
    
    // 定期发送心跳
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if err := registry.Heartbeat("your-service", instanceInfo.InstanceID); err != nil {
                    log.Printf("failed to send heartbeat: %v", err)
                }
            }
        }
    }()
    
    // 优雅关闭
    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c
        
        // 注销服务
        registry.Deregister("your-service", instanceInfo.InstanceID)
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
    registry := NewEurekaRegistry("http://localhost:8761/eureka")
    
    // 发现服务
    instances, err := registry.Discover("your-service")
    if err != nil {
        log.Fatalf("failed to discover service: %v", err)
    }
    
    if len(instances) == 0 {
        log.Fatalf("no service instances found")
    }
    
    // 选择健康的实例
    var healthyInstance *InstanceInfo
    for _, instance := range instances {
        if instance.Status == "UP" {
            healthyInstance = instance
            break
        }
    }
    
    if healthyInstance == nil {
        log.Fatalf("no healthy service instances found")
    }
    
    // 连接服务
    address := fmt.Sprintf("%s:%d", healthyInstance.IPAddr, healthyInstance.Port.Port)
    conn, err := grpc.Dial(address, grpc.WithInsecure())
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

### 1. 集群配置

```go
// 多节点 Eureka 集群
eurekaURLs := []string{
    "http://eureka1:8761/eureka",
    "http://eureka2:8761/eureka",
    "http://eureka3:8761/eureka",
}

// 轮询选择 Eureka 服务器
func (er *EurekaRegistry) selectEurekaServer() string {
    // 简单的轮询实现
    return eurekaURLs[time.Now().Unix()%int64(len(eurekaURLs))]
}
```

### 2. 负载均衡

```go
// 简单的负载均衡
func (er *EurekaRegistry) DiscoverWithLoadBalance(app string) (*InstanceInfo, error) {
    instances, err := er.Discover(app)
    if err != nil {
        return nil, err
    }
    
    if len(instances) == 0 {
        return nil, fmt.Errorf("no instances found")
    }
    
    // 简单的轮询负载均衡
    index := time.Now().Unix() % int64(len(instances))
    return instances[index], nil
}
```

### 3. 健康检查

```go
// 健康检查
func (er *EurekaRegistry) HealthCheck(instance *InstanceInfo) bool {
    url := fmt.Sprintf("http://%s:%d/health", instance.IPAddr, instance.Port.Port)
    
    resp, err := er.client.Get(url)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == http.StatusOK
}
```

## 📊 性能优化

### 1. 连接池配置

```go
// 配置 HTTP 客户端连接池
func NewEurekaRegistry(eurekaURL string) *EurekaRegistry {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    }
    
    return &EurekaRegistry{
        eurekaURL: eurekaURL,
        client: &http.Client{
            Transport: transport,
            Timeout:   10 * time.Second,
        },
    }
}
```

### 2. 缓存优化

```go
// 使用缓存减少 API 调用
type CachedEurekaRegistry struct {
    *EurekaRegistry
    cache map[string][]*InstanceInfo
    mutex sync.RWMutex
}

func (cer *CachedEurekaRegistry) Discover(app string) ([]*InstanceInfo, error) {
    cer.mutex.RLock()
    if instances, ok := cer.cache[app]; ok {
        cer.mutex.RUnlock()
        return instances, nil
    }
    cer.mutex.RUnlock()
    
    instances, err := cer.EurekaRegistry.Discover(app)
    if err != nil {
        return nil, err
    }
    
    cer.mutex.Lock()
    cer.cache[app] = instances
    cer.mutex.Unlock()
    
    return instances, nil
}
```

## 🛡️ 安全配置

### 1. 认证配置

```go
// 添加基本认证
func (er *EurekaRegistry) RegisterWithAuth(instanceInfo *InstanceInfo, username, password string) error {
    url := fmt.Sprintf("%s/apps/%s", er.eurekaURL, instanceInfo.App)
    
    data := map[string]interface{}{
        "instance": instanceInfo,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.SetBasicAuth(username, password)
    
    resp, err := er.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to register service: %d", resp.StatusCode)
    }
    
    return nil
}
```

## 🔍 监控和调试

### 1. 日志配置

```go
// 添加详细日志
func (er *EurekaRegistry) Register(instanceInfo *InstanceInfo) error {
    log.Printf("Registering service: %s", instanceInfo.App)
    
    url := fmt.Sprintf("%s/apps/%s", er.eurekaURL, instanceInfo.App)
    
    data := map[string]interface{}{
        "instance": instanceInfo,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Printf("Failed to marshal instance info: %v", err)
        return err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Failed to create request: %v", err)
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := er.client.Do(req)
    if err != nil {
        log.Printf("Failed to send request: %v", err)
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent {
        log.Printf("Failed to register service: %d", resp.StatusCode)
        return fmt.Errorf("failed to register service: %d", resp.StatusCode)
    }
    
    log.Printf("Successfully registered service: %s", instanceInfo.App)
    return nil
}
```

## 📚 最佳实践

1. **服务命名**: 使用有意义的服务名称
2. **心跳间隔**: 设置合理的心跳间隔
3. **自我保护**: 了解 Eureka 的自我保护机制
4. **监控告警**: 监控服务注册状态
5. **优雅关闭**: 确保服务关闭时正确注销

## ⚠️ 注意事项

- Eureka 已经停止维护，建议迁移到其他方案
- 仅适用于遗留系统维护
- 新项目不建议使用 Eureka

## 🔗 相关资源

- [Eureka 官方文档](https://github.com/Netflix/eureka)
- [Spring Cloud Eureka](https://spring.io/projects/spring-cloud-netflix)
- [Eureka 迁移指南](https://spring.io/blog/2018/12/12/spring-cloud-greenwich-rc1-available-now)
