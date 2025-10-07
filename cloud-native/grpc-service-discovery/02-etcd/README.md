# Etcd 服务注册与发现

## 📖 概述

Etcd 是 CoreOS 开源的高可用键值存储系统，被广泛用作 Kubernetes 的存储后端。它提供了强一致性、高可用性和高性能的服务注册与发现功能。

## 🎯 核心特性

### 1. 强一致性
- 基于 Raft 算法保证强一致性
- 支持线性化读取
- 分布式锁和选举

### 2. 高可用性
- 支持集群部署
- 自动故障转移
- 数据持久化

### 3. 高性能
- 低延迟读写
- 高吞吐量
- 内存优化

## 🚀 快速开始

### 1. 安装 Etcd

```bash
# 下载 Etcd
wget https://github.com/etcd-io/etcd/releases/download/v3.5.7/etcd-v3.5.7-linux-amd64.tar.gz
tar -xzf etcd-v3.5.7-linux-amd64.tar.gz
cd etcd-v3.5.7-linux-amd64
sudo cp etcd etcdctl /usr/local/bin/

# 验证安装
etcd --version
```

### 2. 启动 Etcd 服务器

```bash
# 单节点模式
etcd --name node1 --data-dir /tmp/etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379

# 集群模式
etcd --name node1 --data-dir /tmp/etcd1 --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --initial-cluster node1=http://0.0.0.0:2380,node2=http://0.0.0.0:2480,node3=http://0.0.0.0:2580 --initial-cluster-token etcd-cluster-1 --initial-advertise-peer-urls http://0.0.0.0:2380 --listen-peer-urls http://0.0.0.0:2380
```

### 3. Go 客户端实现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "path"
    "strings"
    "time"

    "go.etcd.io/etcd/clientv3"
    "google.golang.org/grpc"
)

// EtcdRegistry Etcd 服务注册器
type EtcdRegistry struct {
    client *clientv3.Client
    prefix string
}

// NewEtcdRegistry 创建 Etcd 服务注册器
func NewEtcdRegistry(endpoints []string, prefix string) (*EtcdRegistry, error) {
    config := clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    }
    
    client, err := clientv3.New(config)
    if err != nil {
        return nil, err
    }
    
    return &EtcdRegistry{
        client: client,
        prefix: prefix,
    }, nil
}

// Register 注册服务
func (er *EtcdRegistry) Register(serviceName, serviceID, address string, port int, ttl int64) error {
    key := path.Join(er.prefix, serviceName, serviceID)
    value := fmt.Sprintf("%s:%d", address, port)
    
    // 创建租约
    lease, err := er.client.Grant(context.Background(), ttl)
    if err != nil {
        return err
    }
    
    // 注册服务
    _, err = er.client.Put(context.Background(), key, value, clientv3.WithLease(lease.ID))
    if err != nil {
        return err
    }
    
    // 保持租约活跃
    ch, kaerr := er.client.KeepAlive(context.Background(), lease.ID)
    if kaerr != nil {
        return kaerr
    }
    
    // 处理租约续期
    go func() {
        for ka := range ch {
            log.Printf("KeepAlive response: %v", ka)
        }
    }()
    
    return nil
}

// Deregister 注销服务
func (er *EtcdRegistry) Deregister(serviceName, serviceID string) error {
    key := path.Join(er.prefix, serviceName, serviceID)
    _, err := er.client.Delete(context.Background(), key)
    return err
}

// Discover 发现服务
func (er *EtcdRegistry) Discover(serviceName string) ([]string, error) {
    key := path.Join(er.prefix, serviceName)
    resp, err := er.client.Get(context.Background(), key, clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, kv := range resp.Kvs {
        addresses = append(addresses, string(kv.Value))
    }
    
    return addresses, nil
}

// Watch 监听服务变化
func (er *EtcdRegistry) Watch(serviceName string, callback func([]string)) {
    key := path.Join(er.prefix, serviceName)
    
    go func() {
        rch := er.client.Watch(context.Background(), key, clientv3.WithPrefix())
        for wresp := range rch {
            for _, ev := range wresp.Events {
                log.Printf("Watch event: %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
                
                // 重新获取所有服务实例
                addresses, err := er.Discover(serviceName)
                if err != nil {
                    log.Printf("Error discovering services: %v", err)
                    continue
                }
                
                callback(addresses)
            }
        }
    }()
}

// Close 关闭客户端
func (er *EtcdRegistry) Close() error {
    return er.client.Close()
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
    registry, err := NewEtcdRegistry([]string{"localhost:2379"}, "/services")
    if err != nil {
        log.Fatalf("failed to create registry: %v", err)
    }
    defer registry.Close()
    
    // 注册服务
    serviceID := fmt.Sprintf("%s-%s", "your-service", getLocalIP())
    err = registry.Register("your-service", serviceID, getLocalIP(), 8080, 30)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
    
    // 优雅关闭
    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c
        
        // 注销服务
        registry.Deregister("your-service", serviceID)
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
    registry, err := NewEtcdRegistry([]string{"localhost:2379"}, "/services")
    if err != nil {
        log.Fatalf("failed to create registry: %v", err)
    }
    defer registry.Close()
    
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

### 1. 集群配置

```go
// 多节点 Etcd 集群
endpoints := []string{
    "http://etcd1:2379",
    "http://etcd2:2379",
    "http://etcd3:2379",
}

config := clientv3.Config{
    Endpoints:   endpoints,
    DialTimeout: 5 * time.Second,
    Username:    "root",
    Password:    "password",
}
```

### 2. TLS 配置

```go
config := clientv3.Config{
    Endpoints:   []string{"https://etcd.example.com:2379"},
    DialTimeout: 5 * time.Second,
    TLS: &tls.Config{
        InsecureSkipVerify: false,
        ServerName:         "etcd.example.com",
        Certificates:       []tls.Certificate{cert},
        RootCAs:            rootCAs,
    },
}
```

### 3. 连接池配置

```go
config := clientv3.Config{
    Endpoints:   endpoints,
    DialTimeout: 5 * time.Second,
    DialKeepAliveTime:    30 * time.Second,
    DialKeepAliveTimeout: 5 * time.Second,
    MaxCallSendMsgSize:   2 * 1024 * 1024, // 2MB
    MaxCallRecvMsgSize:   4 * 1024 * 1024, // 4MB
}
```

## 📊 性能优化

### 1. 批量操作

```go
// 批量注册服务
func (er *EtcdRegistry) BatchRegister(services []ServiceInfo) error {
    var ops []clientv3.Op
    for _, service := range services {
        key := path.Join(er.prefix, service.Name, service.ID)
        value := fmt.Sprintf("%s:%d", service.Address, service.Port)
        ops = append(ops, clientv3.OpPut(key, value))
    }
    
    _, err := er.client.Txn(context.Background()).Then(ops...).Commit()
    return err
}
```

### 2. 缓存优化

```go
// 使用缓存减少 Etcd 查询
type CachedRegistry struct {
    *EtcdRegistry
    cache map[string][]string
    mutex sync.RWMutex
}

func (cr *CachedRegistry) Discover(serviceName string) ([]string, error) {
    cr.mutex.RLock()
    if addresses, ok := cr.cache[serviceName]; ok {
        cr.mutex.RUnlock()
        return addresses, nil
    }
    cr.mutex.RUnlock()
    
    addresses, err := cr.EtcdRegistry.Discover(serviceName)
    if err != nil {
        return nil, err
    }
    
    cr.mutex.Lock()
    cr.cache[serviceName] = addresses
    cr.mutex.Unlock()
    
    return addresses, nil
}
```

## ��️ 安全配置

### 1. 认证配置

```go
config := clientv3.Config{
    Endpoints:   endpoints,
    DialTimeout: 5 * time.Second,
    Username:    "root",
    Password:    "password",
}
```

### 2. 权限控制

```go
// 创建用户
client, err := clientv3.New(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 创建用户
_, err = client.Auth.UserAdd(context.Background(), "user1", "password1")
if err != nil {
    log.Fatal(err)
}

// 创建角色
_, err = client.Auth.RoleAdd(context.Background(), "role1")
if err != nil {
    log.Fatal(err)
}

// 分配权限
_, err = client.Auth.RoleGrantPermission(context.Background(), "role1", "/services/", "", clientv3.PermissionTypePrefix, clientv3.PermissionReadWrite)
if err != nil {
    log.Fatal(err)
}

// 分配角色给用户
_, err = client.Auth.UserGrantRole(context.Background(), "user1", "role1")
if err != nil {
    log.Fatal(err)
}
```

## 🔍 监控和调试

### 1. 健康检查

```go
// 检查 Etcd 集群健康状态
func (er *EtcdRegistry) HealthCheck() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := er.client.Status(ctx, er.client.Endpoints()[0])
    return err
}
```

### 2. 指标监控

```go
// 使用 Prometheus 监控
import "github.com/prometheus/client_golang/prometheus"

var (
    etcdOperations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "etcd_operations_total",
            Help: "Total number of etcd operations",
        },
        []string{"operation", "status"},
    )
)

func init() {
    prometheus.MustRegister(etcdOperations)
}
```

## 📚 最佳实践

1. **集群部署**: 使用奇数个节点部署 Etcd 集群
2. **数据备份**: 定期备份 Etcd 数据
3. **监控告警**: 监控集群健康状态和性能指标
4. **资源限制**: 设置合理的内存和磁盘限制
5. **网络优化**: 优化网络配置和连接池

## 🔗 相关资源

- [Etcd 官方文档](https://etcd.io/docs/)
- [Etcd Go 客户端文档](https://pkg.go.dev/go.etcd.io/etcd/clientv3)
- [Etcd 最佳实践](https://etcd.io/docs/latest/op-guide/)
