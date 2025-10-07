# Nacos 服务注册与发现

## 📖 概述

Nacos 是阿里巴巴开源的服务发现和配置管理平台，支持服务注册与发现、配置管理、服务管理等功能，是云原生应用的重要基础设施。

## 🎯 核心特性

### 1. 服务注册与发现
- 支持多种服务类型
- 健康检查机制
- 服务元数据管理

### 2. 配置管理
- 动态配置推送
- 配置版本管理
- 配置监听

### 3. 服务管理
- 服务健康检查
- 服务权重管理
- 服务分组管理

## 🚀 快速开始

### 1. 安装 Nacos

```bash
# 下载 Nacos
wget https://github.com/alibaba/nacos/releases/download/2.2.3/nacos-server-2.2.3.tar.gz
tar -xzf nacos-server-2.2.3.tar.gz
cd nacos/bin

# 启动 Nacos (单机模式)
sh startup.sh -m standalone

# 启动 Nacos (集群模式)
sh startup.sh
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

// NacosRegistry Nacos 服务注册器
type NacosRegistry struct {
    serverAddr string
    client     *http.Client
}

// NewNacosRegistry 创建 Nacos 服务注册器
func NewNacosRegistry(serverAddr string) *NacosRegistry {
    return &NacosRegistry{
        serverAddr: serverAddr,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// Instance 实例信息
type Instance struct {
    IP          string            `json:"ip"`
    Port        int               `json:"port"`
    Weight      float64           `json:"weight"`
    Enabled     bool              `json:"enabled"`
    Healthy     bool              `json:"healthy"`
    Metadata    map[string]string `json:"metadata"`
    ClusterName string            `json:"clusterName"`
    ServiceName string            `json:"serviceName"`
    GroupName   string            `json:"groupName"`
}

// Register 注册服务
func (nr *NacosRegistry) Register(instance *Instance) error {
    url := fmt.Sprintf("http://%s/nacos/v1/ns/instance", nr.serverAddr)
    
    params := map[string]string{
        "serviceName": instance.ServiceName,
        "ip":          instance.IP,
        "port":        fmt.Sprintf("%d", instance.Port),
        "weight":      fmt.Sprintf("%.2f", instance.Weight),
        "enabled":     fmt.Sprintf("%t", instance.Enabled),
        "healthy":     fmt.Sprintf("%t", instance.Healthy),
        "metadata":    encodeMetadata(instance.Metadata),
        "clusterName": instance.ClusterName,
        "groupName":   instance.GroupName,
    }
    
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return err
    }
    
    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to register service: %d, %s", resp.StatusCode, string(body))
    }
    
    return nil
}

// Deregister 注销服务
func (nr *NacosRegistry) Deregister(serviceName, ip string, port int, groupName string) error {
    url := fmt.Sprintf("http://%s/nacos/v1/ns/instance", nr.serverAddr)
    
    params := map[string]string{
        "serviceName": serviceName,
        "ip":          ip,
        "port":        fmt.Sprintf("%d", port),
        "groupName":   groupName,
    }
    
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }
    
    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to deregister service: %d, %s", resp.StatusCode, string(body))
    }
    
    return nil
}

// Discover 发现服务
func (nr *NacosRegistry) Discover(serviceName, groupName string) ([]*Instance, error) {
    url := fmt.Sprintf("http://%s/nacos/v1/ns/instance/list", nr.serverAddr)
    
    params := map[string]string{
        "serviceName": serviceName,
        "groupName":   groupName,
    }
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("failed to discover service: %d, %s", resp.StatusCode, string(body))
    }
    
    var result struct {
        Hosts []*Instance `json:"hosts"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Hosts, nil
}

// Heartbeat 发送心跳
func (nr *NacosRegistry) Heartbeat(serviceName, ip string, port int, groupName string) error {
    url := fmt.Sprintf("http://%s/nacos/v1/ns/instance/beat", nr.serverAddr)
    
    beatInfo := map[string]interface{}{
        "ip":          ip,
        "port":        port,
        "serviceName": serviceName,
        "groupName":   groupName,
    }
    
    jsonData, err := json.Marshal(beatInfo)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to send heartbeat: %d, %s", resp.StatusCode, string(body))
    }
    
    return nil
}

// 编码元数据
func encodeMetadata(metadata map[string]string) string {
    if len(metadata) == 0 {
        return ""
    }
    
    jsonData, err := json.Marshal(metadata)
    if err != nil {
        return ""
    }
    
    return string(jsonData)
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
    registry := NewNacosRegistry("localhost:8848")
    
    // 注册服务
    instance := &Instance{
        IP:          getLocalIP(),
        Port:        8080,
        Weight:      1.0,
        Enabled:     true,
        Healthy:     true,
        Metadata: map[string]string{
            "version": "1.0.0",
            "region":  "us-west-2",
        },
        ClusterName: "DEFAULT",
        ServiceName: "your-service",
        GroupName:   "DEFAULT_GROUP",
    }
    
    err = registry.Register(instance)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
    
    // 定期发送心跳
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if err := registry.Heartbeat("your-service", instance.IP, instance.Port, "DEFAULT_GROUP"); err != nil {
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
        registry.Deregister("your-service", instance.IP, instance.Port, "DEFAULT_GROUP")
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
    registry := NewNacosRegistry("localhost:8848")
    
    // 发现服务
    instances, err := registry.Discover("your-service", "DEFAULT_GROUP")
    if err != nil {
        log.Fatalf("failed to discover service: %v", err)
    }
    
    if len(instances) == 0 {
        log.Fatalf("no service instances found")
    }
    
    // 选择健康的实例
    var healthyInstance *Instance
    for _, instance := range instances {
        if instance.Healthy && instance.Enabled {
            healthyInstance = instance
            break
        }
    }
    
    if healthyInstance == nil {
        log.Fatalf("no healthy service instances found")
    }
    
    // 连接服务
    address := fmt.Sprintf("%s:%d", healthyInstance.IP, healthyInstance.Port)
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
// 多节点 Nacos 集群
serverAddrs := []string{
    "nacos1:8848",
    "nacos2:8848",
    "nacos3:8848",
}

// 轮询选择 Nacos 服务器
func (nr *NacosRegistry) selectNacosServer() string {
    return serverAddrs[time.Now().Unix()%int64(len(serverAddrs))]
}
```

### 2. 负载均衡

```go
// 基于权重的负载均衡
func (nr *NacosRegistry) DiscoverWithLoadBalance(serviceName, groupName string) (*Instance, error) {
    instances, err := nr.Discover(serviceName, groupName)
    if err != nil {
        return nil, err
    }
    
    if len(instances) == 0 {
        return nil, fmt.Errorf("no instances found")
    }
    
    // 过滤健康的实例
    var healthyInstances []*Instance
    for _, instance := range instances {
        if instance.Healthy && instance.Enabled {
            healthyInstances = append(healthyInstances, instance)
        }
    }
    
    if len(healthyInstances) == 0 {
        return nil, fmt.Errorf("no healthy instances found")
    }
    
    // 基于权重的负载均衡
    totalWeight := 0.0
    for _, instance := range healthyInstances {
        totalWeight += instance.Weight
    }
    
    if totalWeight == 0 {
        return healthyInstances[0], nil
    }
    
    random := time.Now().UnixNano() % int64(totalWeight*100)
    currentWeight := 0.0
    
    for _, instance := range healthyInstances {
        currentWeight += instance.Weight
        if float64(random) < currentWeight*100 {
            return instance, nil
        }
    }
    
    return healthyInstances[0], nil
}
```

### 3. 配置管理

```go
// 配置管理
func (nr *NacosRegistry) GetConfig(dataId, group string) (string, error) {
    url := fmt.Sprintf("http://%s/nacos/v1/cs/configs", nr.serverAddr)
    
    params := map[string]string{
        "dataId": dataId,
        "group":  group,
    }
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }
    
    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to get config: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    return string(body), nil
}

// 监听配置变化
func (nr *NacosRegistry) ListenConfig(dataId, group string, callback func(string)) {
    go func() {
        for {
            config, err := nr.GetConfig(dataId, group)
            if err != nil {
                log.Printf("Failed to get config: %v", err)
                time.Sleep(5 * time.Second)
                continue
            }
            
            callback(config)
            time.Sleep(10 * time.Second)
        }
    }()
}
```

## 📊 性能优化

### 1. 连接池配置

```go
// 配置 HTTP 客户端连接池
func NewNacosRegistry(serverAddr string) *NacosRegistry {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    }
    
    return &NacosRegistry{
        serverAddr: serverAddr,
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
type CachedNacosRegistry struct {
    *NacosRegistry
    cache map[string][]*Instance
    mutex sync.RWMutex
}

func (cnr *CachedNacosRegistry) Discover(serviceName, groupName string) ([]*Instance, error) {
    key := fmt.Sprintf("%s:%s", serviceName, groupName)
    
    cnr.mutex.RLock()
    if instances, ok := cnr.cache[key]; ok {
        cnr.mutex.RUnlock()
        return instances, nil
    }
    cnr.mutex.RUnlock()
    
    instances, err := cnr.NacosRegistry.Discover(serviceName, groupName)
    if err != nil {
        return nil, err
    }
    
    cnr.mutex.Lock()
    cnr.cache[key] = instances
    cnr.mutex.Unlock()
    
    return instances, nil
}
```

## 🛡️ 安全配置

### 1. 认证配置

```go
// 添加认证
func (nr *NacosRegistry) RegisterWithAuth(instance *Instance, username, password string) error {
    // 先获取访问令牌
    token, err := nr.getAccessToken(username, password)
    if err != nil {
        return err
    }
    
    url := fmt.Sprintf("http://%s/nacos/v1/ns/instance", nr.serverAddr)
    
    params := map[string]string{
        "serviceName": instance.ServiceName,
        "ip":          instance.IP,
        "port":        fmt.Sprintf("%d", instance.Port),
        "weight":      fmt.Sprintf("%.2f", instance.Weight),
        "enabled":     fmt.Sprintf("%t", instance.Enabled),
        "healthy":     fmt.Sprintf("%t", instance.Healthy),
        "metadata":    encodeMetadata(instance.Metadata),
        "clusterName": instance.ClusterName,
        "groupName":   instance.GroupName,
        "accessToken": token,
    }
    
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return err
    }
    
    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to register service: %d, %s", resp.StatusCode, string(body))
    }
    
    return nil
}

// 获取访问令牌
func (nr *NacosRegistry) getAccessToken(username, password string) (string, error) {
    url := fmt.Sprintf("http://%s/nacos/v1/auth/login", nr.serverAddr)
    
    params := map[string]string{
        "username": username,
        "password": password,
    }
    
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return "", err
    }
    
    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
    
    resp, err := nr.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to get access token: %d", resp.StatusCode)
    }
    
    var result struct {
        AccessToken string `json:"accessToken"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    
    return result.AccessToken, nil
}
```

## 🔍 监控和调试

### 1. 健康检查

```go
// 检查 Nacos 服务器健康状态
func (nr *NacosRegistry) HealthCheck() error {
    url := fmt.Sprintf("http://%s/nacos/v1/ns/operator/metrics", nr.serverAddr)
    
    resp, err := nr.client.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("nacos server is not healthy: %d", resp.StatusCode)
    }
    
    return nil
}
```

### 2. 指标监控

```go
// 使用 Prometheus 监控
import "github.com/prometheus/client_golang/prometheus"

var (
    nacosOperations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "nacos_operations_total",
            Help: "Total number of nacos operations",
        },
        []string{"operation", "status"},
    )
)

func init() {
    prometheus.MustRegister(nacosOperations)
}
```

## 📚 最佳实践

1. **服务分组**: 使用合理的服务分组策略
2. **权重管理**: 根据实例性能设置合适的权重
3. **健康检查**: 实现快速响应的健康检查
4. **配置管理**: 利用 Nacos 的配置管理功能
5. **监控告警**: 监控服务注册状态和配置变化

## �� 相关资源

- [Nacos 官方文档](https://nacos.io/zh-cn/docs/)
- [Nacos Go SDK](https://github.com/nacos-group/nacos-sdk-go)
- [Nacos 最佳实践](https://nacos.io/zh-cn/docs/best-practices.html)
