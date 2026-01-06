# 微服务注册与发现

## 概述

服务注册与发现（Service Registry and Discovery）是微服务架构中的核心组件，用于解决服务之间的动态发现和通信问题。

### 为什么需要服务注册与发现？

在微服务架构中，服务实例是动态变化的：
- 服务实例会启动、停止、重启
- 服务实例的 IP 和端口可能变化
- 服务实例可能分布在不同的机器上
- 需要实现负载均衡和故障转移

**没有服务注册与发现的问题：**
- 硬编码服务地址，难以维护
- 服务实例变化时需要手动更新配置
- 无法实现动态负载均衡
- 难以实现服务健康检查

**有了服务注册与发现的好处：**
- 服务自动注册和注销
- 客户端自动发现可用服务实例
- 支持负载均衡和故障转移
- 支持服务健康检查

## 架构原理

### 服务注册流程

```
┌─────────────┐
│  服务实例   │
│  (启动时)   │
└──────┬──────┘
       │ 1. 注册服务信息
       │    (服务名、IP、端口、元数据)
       ▼
┌─────────────┐
│  注册中心   │
│ (Registry)  │
└──────┬──────┘
       │ 2. 存储服务信息
       │ 3. 定期健康检查
       │ 4. 服务下线时删除
```

### 服务发现流程

```
┌─────────────┐
│  客户端     │
│ (Consumer)  │
└──────┬──────┘
       │ 1. 查询服务列表
       │    (通过服务名)
       ▼
┌─────────────┐
│  注册中心   │
│ (Registry)  │
└──────┬──────┘
       │ 2. 返回可用实例列表
       │    (IP、端口、健康状态)
       ▼
┌─────────────┐
│  客户端     │
│             │
│ 3. 负载均衡 │
│ 4. 调用服务 │
└─────────────┘
```

## Kratos 框架支持

Kratos 框架内置了服务注册与发现的支持，通过 `registry` 和 `discovery` 接口实现。

### 支持的注册中心

Kratos 支持多种注册中心：

1. **Consul** - HashiCorp 的分布式服务发现和配置管理系统
2. **etcd** - 分布式键值存储系统
3. **Nacos** - 阿里巴巴的服务发现和配置管理平台
4. **Kubernetes** - 基于 Kubernetes 的服务发现
5. **Zookeeper** - Apache 的分布式协调服务

### 核心接口

```go
// 注册中心接口
type Registrar interface {
    Register(ctx context.Context, service *ServiceInstance) error
    Deregister(ctx context.Context, service *ServiceInstance) error
}

// 服务发现接口
type Discovery interface {
    GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
    Watch(ctx context.Context, serviceName string) (Watcher, error)
}
```

## 实现步骤

### 步骤 1: 添加依赖

根据选择的注册中心，添加对应的依赖：

#### Consul

```bash
go get github.com/go-kratos/kratos/contrib/registry/consul/v2
go get github.com/hashicorp/consul/api
```

#### etcd

```bash
go get github.com/go-kratos/kratos/contrib/registry/etcd/v2
go get go.etcd.io/etcd/client/v3
```

#### Nacos

```bash
go get github.com/go-kratos/kratos/contrib/registry/nacos/v2
go get github.com/nacos-group/nacos-sdk-go/v2
```

### 步骤 2: 更新配置定义

在 `internal/conf/conf.proto` 中添加注册中心配置：

```protobuf
syntax = "proto3";
package kratos.api;

option go_package = "sre/internal/conf;conf";

import "google/protobuf/duration.proto";

message Bootstrap {
  Server server = 1;
  Data data = 2;
  Registry registry = 3;  // 新增：注册中心配置
}

message Registry {
  message Consul {
    string address = 1;  // Consul 地址，如 "127.0.0.1:8500"
    string scheme = 2;   // 协议，如 "http" 或 "https"
  }
  message Etcd {
    repeated string endpoints = 1;  // etcd 地址列表，如 ["127.0.0.1:2379"]
    int64 timeout = 2;              // 超时时间（秒）
  }
  message Nacos {
    repeated string endpoints = 1;  // Nacos 地址列表
    string namespace = 2;           // 命名空间
    string username = 3;           // 用户名
    string password = 4;             // 密码
  }
  
  Consul consul = 1;
  Etcd etcd = 2;
  Nacos nacos = 3;
}
```

重新生成配置代码：

```bash
make config
```

### 步骤 3: 配置文件

在 `configs/config.yaml` 中添加注册中心配置：

```yaml
server:
  http:
    network: tcp
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    network: tcp
    addr: 0.0.0.0:9000
    timeout: 1s

data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    network: tcp
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s

# 注册中心配置
registry:
  consul:
    address: "127.0.0.1:8500"
    scheme: "http"
  # 或者使用 etcd
  # etcd:
  #   endpoints:
  #     - "127.0.0.1:2379"
  #   timeout: 5
  # 或者使用 Nacos
  # nacos:
  #   endpoints:
  #     - "127.0.0.1:8848"
  #   namespace: "public"
  #   username: "nacos"
  #   password: "nacos"
```

### 步骤 4: 创建注册中心客户端

创建 `internal/registry/registry.go`：

```go
package registry

import (
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/registry"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	nacos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	
	consulAPI "github.com/hashicorp/consul/api"
	etcdClient "go.etcd.io/etcd/client/v3"
	nacosClient "github.com/nacos-group/nacos-sdk-go/v2/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosVo "github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NewRegistry 根据配置创建注册中心客户端
func NewRegistry(c *conf.Registry) (registry.Registrar, error) {
	// Consul
	if c.Consul != nil && c.Consul.Address != "" {
		consulConfig := consulAPI.DefaultConfig()
		consulConfig.Address = c.Consul.Address
		consulConfig.Scheme = c.Consul.Scheme
		client, err := consulAPI.NewClient(consulConfig)
		if err != nil {
			return nil, err
		}
		return consul.New(client), nil
	}
	
	// etcd
	if c.Etcd != nil && len(c.Etcd.Endpoints) > 0 {
		etcdClient, err := etcdClient.New(etcdClient.Config{
			Endpoints: c.Etcd.Endpoints,
		})
		if err != nil {
			return nil, err
		}
		return etcd.New(etcdClient), nil
	}
	
	// Nacos
	if c.Nacos != nil && len(c.Nacos.Endpoints) > 0 {
		sc := []nacosConstant.ServerConfig{
			{
				IpAddr: c.Nacos.Endpoints[0],
				Port:   8848,
			},
		}
		cc := nacosConstant.ClientConfig{
			NamespaceId:         c.Nacos.Namespace,
			Username:            c.Nacos.Username,
			Password:            c.Nacos.Password,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
		}
		nacosClient, err := nacosClient.NewNamingClient(
			nacosVo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
		if err != nil {
			return nil, err
		}
		return nacos.New(nacosClient), nil
	}
	
	return nil, nil
}
```

### 步骤 5: 创建服务发现客户端

创建 `internal/registry/discovery.go`：

```go
package registry

import (
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/registry"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	nacos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	
	consulAPI "github.com/hashicorp/consul/api"
	etcdClient "go.etcd.io/etcd/client/v3"
	nacosClient "github.com/nacos-group/nacos-sdk-go/v2/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosVo "github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	// Consul
	if c.Consul != nil && c.Consul.Address != "" {
		consulConfig := consulAPI.DefaultConfig()
		consulConfig.Address = c.Consul.Address
		consulConfig.Scheme = c.Consul.Scheme
		client, err := consulAPI.NewClient(consulConfig)
		if err != nil {
			return nil, err
		}
		return consul.New(client), nil
	}
	
	// etcd
	if c.Etcd != nil && len(c.Etcd.Endpoints) > 0 {
		etcdClient, err := etcdClient.New(etcdClient.Config{
			Endpoints: c.Etcd.Endpoints,
		})
		if err != nil {
			return nil, err
		}
		return etcd.New(etcdClient), nil
	}
	
	// Nacos
	if c.Nacos != nil && len(c.Nacos.Endpoints) > 0 {
		sc := []nacosConstant.ServerConfig{
			{
				IpAddr: c.Nacos.Endpoints[0],
				Port:   8848,
			},
		}
		cc := nacosConstant.ClientConfig{
			NamespaceId:         c.Nacos.Namespace,
			Username:            c.Nacos.Username,
			Password:            c.Nacos.Password,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
		}
		nacosClient, err := nacosClient.NewNamingClient(
			nacosVo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
		if err != nil {
			return nil, err
		}
		return nacos.New(nacosClient), nil
	}
	
	return nil, nil
}
```

### 步骤 6: 更新 Wire 配置

更新 `cmd/sre/wire.go`：

```go
//go:build wireinject
// +build wireinject

package main

import (
	"sre/internal/biz"
	"sre/internal/conf"
	"sre/internal/data"
	"sre/internal/registry"
	"sre/internal/server"
	"sre/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is wire providers.
var ProviderSet = wire.NewSet(
	server.ProviderSet,
	data.ProviderSet,
	biz.ProviderSet,
	service.ProviderSet,
	registry.NewRegistry,  // 新增
	registry.NewDiscovery, // 新增
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(ProviderSet, newApp))
}
```

### 步骤 7: 更新 main.go

更新 `cmd/sre/main.go`，在创建应用时注册服务：

```go
package main

import (
	"flag"
	"os"
	
	"sre/internal/conf"
	"sre/internal/registry"
	
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	
	_ "go.uber.org/automaxprocs"
)

var (
	Name     string
	Version  string
	flagconf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(
	logger log.Logger,
	gs *grpc.Server,
	hs *http.Server,
	registrar registry.Registrar,  // 新增：注册中心
) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
		kratos.Registrar(registrar),  // 新增：注册服务
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Registry, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

### 步骤 8: 更新 wire_gen.go

运行 Wire 生成代码：

```bash
make wire
```

### 步骤 9: 设置服务名称和版本 ⚠️ 重要

**服务名称（Name）和版本（Version）必须在编译时通过 `-ldflags` 参数设置**，否则服务无法正确注册到注册中心。

#### 方法 1: 使用 Makefile（推荐）

项目已配置 Makefile，默认服务名称为 `sre`：

```bash
# 使用默认服务名称 sre
make build

# 或者指定自定义服务名称
make build NAME=your-service-name
```

Makefile 会自动设置：
- `-X main.Name=$(NAME)`（默认为 `sre`）
- `-X main.Version=$(VERSION)`（从 git tag 获取）

#### 方法 2: 直接使用 go build

```bash
# 设置服务名称和版本
go build -ldflags "-X main.Name=sre -X main.Version=v1.0.0" -o ./bin/sre ./cmd/sre
```

#### 方法 3: 在代码中设置默认值（不推荐）

如果未通过 `-ldflags` 设置，服务名称将为空，导致：
- etcd 中的注册键为 `/microservices//{ID}`（中间有两个斜杠）
- 服务无法被正确发现
- 多个服务实例无法区分

#### 验证服务名称设置

启动服务后，检查日志中是否包含服务名称：

```
service.id=JianfendeMacBook-Pro.local service.name=sre service.version=v1.0.0
```

如果 `service.name` 为空，说明服务名称未正确设置。

#### 使用检查工具验证注册

```bash
# 运行检查工具查看注册信息
go run ./cmd/check-registry/main.go

# 或指定 etcd 地址
go run ./cmd/check-registry/main.go 127.0.0.1:2379
```

**正确的注册键格式**：
- ✅ `/microservices/sre/grpc/JianfendeMacBook-Pro.local`
- ✅ `/microservices/sre/http/JianfendeMacBook-Pro.local`

**错误的注册键格式**（服务名称为空）：
- ❌ `/microservices//JianfendeMacBook-Pro.local`（中间有两个斜杠）

## 客户端服务发现

### 创建 gRPC 客户端

当服务 A 需要调用服务 B 时，使用服务发现：

```go
package client

import (
	"sre/internal/conf"
	"sre/internal/registry"
	
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/registry/discovery"
)

// NewGreeterClient 创建带服务发现的 gRPC 客户端
func NewGreeterClient(
	registry *conf.Registry,
	discovery registry.Discovery,
) (v1.GreeterClient, error) {
	// 创建 gRPC 连接，使用服务发现
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///greeter"),  // 使用 discovery:// 协议
		grpc.WithDiscovery(discovery),
	)
	if err != nil {
		return nil, err
	}
	
	return v1.NewGreeterClient(conn), nil
}
```

### 使用客户端

```go
// 在业务代码中使用
greeterClient, err := client.NewGreeterClient(registry, discovery)
if err != nil {
	return err
}

reply, err := greeterClient.SayHello(ctx, &v1.HelloRequest{
	Name: "world",
})
```

## 配置说明

### Consul 配置

```yaml
registry:
  consul:
    address: "127.0.0.1:8500"  # Consul 服务器地址
    scheme: "http"              # 协议：http 或 https
```

### etcd 配置

```yaml
registry:
  etcd:
    endpoints:
      - "127.0.0.1:2379"       # etcd 服务器地址列表
      - "127.0.0.1:2380"
    timeout: 5                  # 超时时间（秒）
```

### Nacos 配置

```yaml
registry:
  nacos:
    endpoints:
      - "127.0.0.1:8848"       # Nacos 服务器地址
    namespace: "public"         # 命名空间
    username: "nacos"           # 用户名
    password: "nacos"           # 密码
```

## 服务元数据

可以在注册服务时添加元数据：

```go
func newApp(
	logger log.Logger,
	gs *grpc.Server,
	hs *http.Server,
	registrar registry.Registrar,
) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			"env":      "production",
			"region":   "us-west-1",
			"zone":     "zone-a",
			"weight":   "100",
		}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
		kratos.Registrar(registrar),
	)
}
```

## 服务注册信息

当服务启动并注册到注册中心时，Kratos 框架会自动收集并注册以下信息：

### 1. 基本信息

这些信息通过 `newApp` 函数中的 `kratos.Option` 配置：

```37:50:cmd/sre/main.go
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, registrar kratosRegistry.Registrar) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	}
	// 如果配置了注册中心，则注册服务
	if registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}
	return kratos.New(opts...)
}
```

#### 服务实例 ID (`ID`)
- **来源**: `os.Hostname()`，通常是运行服务的主机名
- **用途**: 唯一标识服务实例
- **示例**: `hostname-001`、`server-01`

#### 服务名称 (`Name`)
- **来源**: 编译时通过 `-ldflags "-X main.Name=xxx"` 传入
- **用途**: 标识服务类型，用于服务发现
- **示例**: `sre`、`user-service`、`order-service`

#### 服务版本 (`Version`)
- **来源**: 编译时通过 `-ldflags "-X main.Version=x.y.z"` 传入
- **用途**: 标识服务版本，支持多版本共存
- **示例**: `v1.0.0`、`v2.1.3`

#### 元数据 (`Metadata`)
- **来源**: 自定义的键值对映射
- **用途**: 存储额外的服务信息，如环境、区域、权重等
- **示例**: 
  ```go
  map[string]string{
      "env":    "production",
      "region": "us-west-1",
      "zone":   "zone-a",
      "weight": "100",
  }
  ```

### 2. 服务端点信息

Kratos 会自动从已注册的服务器中提取端点信息：

#### gRPC 端点
- **来源**: `configs/config.yaml` 中的 `server.grpc.addr` 配置
- **格式**: `grpc://<实际IP>:<端口>`
- **示例**: `grpc://192.168.1.100:9000`
- **说明**: 如果配置为 `0.0.0.0:9000`，Kratos 会自动解析为实际绑定的 IP 地址

#### HTTP 端点
- **来源**: `configs/config.yaml` 中的 `server.http.addr` 配置
- **格式**: `http://<实际IP>:<端口>`
- **示例**: `http://192.168.1.100:8000`
- **说明**: 如果配置为 `0.0.0.0:8000`，Kratos 会自动解析为实际绑定的 IP 地址

### 3. 在注册中心中的存储格式

根据不同的注册中心，服务信息会以不同格式存储：

#### Nacos 注册格式

注册到 Nacos 时，每个服务端点会创建一个独立的服务实例：

**gRPC 服务实例:**
- **服务名称**: `{Name}.grpc`（例如：`sre.grpc`）
- **IP 地址**: 从 gRPC 端点解析的实际 IP
- **端口**: gRPC 配置的端口（如 `9000`）
- **协议类型**: 存储在 metadata 的 `kind` 字段，值为 `grpc`
- **版本信息**: 存储在 metadata 的 `version` 字段
- **权重**: 默认 `100`（可在 Nacos 控制台或配置中调整）
- **集群名称**: 默认 `DEFAULT`（可在配置中调整）
- **分组名称**: 默认 `DEFAULT_GROUP`（可在配置中调整）
- **健康状态**: `Healthy: true`
- **临时实例**: `Ephemeral: true`（服务下线时自动删除）

**HTTP 服务实例:**
- **服务名称**: `{Name}.http`（例如：`sre.http`）
- **IP 地址**: 从 HTTP 端点解析的实际 IP
- **端口**: HTTP 配置的端口（如 `8000`）
- **协议类型**: 存储在 metadata 的 `kind` 字段，值为 `http`
- **版本信息**: 存储在 metadata 的 `version` 字段
- **其他字段**: 与 gRPC 服务实例相同

#### Consul 注册格式

注册到 Consul 时，服务信息存储在 Consul 的 Service Catalog 中：

- **Service ID**: `{Name}-{ID}-{scheme}`（例如：`sre-hostname-001-grpc`）
- **Service Name**: `{Name}.{scheme}`（例如：`sre.grpc`）
- **Address**: 服务 IP 地址
- **Port**: 服务端口
- **Tags**: 包含版本、协议类型等信息
- **Meta**: 包含元数据信息
- **Check**: 健康检查配置

#### etcd 注册格式

注册到 etcd 时，服务信息存储在 etcd 的键值对中：

- **Key**: `/microservices/{Name}/{scheme}/{ID}`
  - 如果服务名称（Name）为空，键格式为：`/microservices//{ID}`（注意中间有两个斜杠）
  - **重要**: 必须通过编译参数设置服务名称，否则服务无法正确注册和发现
- **Value**: JSON 格式的服务实例信息，包含：
  - `id`: 服务实例 ID（主机名）
  - `name`: 服务名称（必须设置）
  - `version`: 服务版本
  - `endpoints`: 服务端点列表（gRPC 和 HTTP）
  - `metadata`: 元数据信息
- **TTL**: 带 TTL 的租约，定期续约保持服务在线

**⚠️ 常见问题**:
- 如果发现注册键为 `/microservices//{ID}`（中间有两个斜杠），说明服务名称未设置
- 解决方法：使用 `-ldflags "-X main.Name=your-service-name"` 编译参数设置服务名称

### 4. 实际注册示例

假设服务配置如下：

**编译参数:**
```bash
go build -ldflags "-X main.Name=sre -X main.Version=v1.0.0"
```

**配置文件 (`configs/config.yaml`):**
```yaml
server:
  http:
    addr: 0.0.0.0:8000
  grpc:
    addr: 0.0.0.0:9000
```

**运行环境:**
- 主机名: `server-01`
- 实际 IP: `192.168.1.100`

**注册结果:**

在 Nacos 中会注册两个服务实例：

1. **gRPC 服务实例**:
   - 服务名: `sre.grpc`
   - IP: `192.168.1.100`
   - 端口: `9000`
   - 实例 ID: `server-01`
   - Metadata: 
     ```json
     {
       "kind": "grpc",
       "version": "v1.0.0"
     }
     ```

2. **HTTP 服务实例**:
   - 服务名: `sre.http`
   - IP: `192.168.1.100`
   - 端口: `8000`
   - 实例 ID: `server-01`
   - Metadata:
     ```json
     {
       "kind": "http",
       "version": "v1.0.0"
     }
     ```

### 5. 查看注册信息

#### Nacos 控制台
1. 访问 Nacos 控制台（默认地址：`http://127.0.0.1:8848/nacos`）
2. 进入「服务管理」→「服务列表」
3. 搜索服务名称（如 `sre.grpc` 或 `sre.http`）
4. 点击服务名称查看实例详情

#### Consul 控制台
1. 访问 Consul UI（默认地址：`http://127.0.0.1:8500/ui`）
2. 进入「Services」页面
3. 查看已注册的服务列表和实例详情

#### etcd
```bash
# 查看所有注册的服务
etcdctl get --prefix /kratos/

# 查看特定服务的实例
etcdctl get --prefix /kratos/sre/
```

### 6. 注意事项

1. **服务名称规范**: 建议使用小写字母和连字符，避免特殊字符
2. **IP 地址解析**: 如果配置为 `0.0.0.0`，Kratos 会自动解析为实际绑定的 IP
3. **多端点注册**: 如果同时配置了 gRPC 和 HTTP，会注册两个独立的服务实例
4. **版本管理**: 通过版本号可以支持多版本共存，实现灰度发布
5. **元数据扩展**: 可以在 `Metadata` 中添加自定义信息，如环境、区域、权重等

## 健康检查

Kratos 会自动处理健康检查：

- **gRPC 服务**：使用 gRPC 健康检查协议
- **HTTP 服务**：提供 `/health` 端点

注册中心会定期检查服务健康状态，不健康的实例会被自动移除。

## 最佳实践

### 1. 服务命名规范

- 使用小写字母和连字符：`user-service`、`order-service`
- 避免使用下划线和特殊字符
- 保持命名简洁且有意义

### 2. 服务版本管理

- 在服务元数据中记录版本号
- 支持多版本共存（灰度发布）
- 使用语义化版本号

### 3. 负载均衡

- 使用客户端负载均衡（Kratos 内置支持）
- 支持多种负载均衡策略：轮询、随机、加权轮询等

### 4. 故障转移

- 自动剔除不健康的服务实例
- 实现重试机制
- 使用熔断器防止雪崩

### 5. 监控和日志

- 记录服务注册和注销事件
- 监控服务发现延迟
- 记录服务调用失败情况

## 常见问题

### Q1: 服务注册失败怎么办？

**A:** 检查以下几点：
- 注册中心是否正常运行
- 网络连接是否正常
- 配置是否正确
- 查看日志中的错误信息

### Q2: 服务发现不到实例？

**A:** 检查以下几点：
- 服务是否已成功注册
- 服务名称是否匹配
- 服务是否健康（健康检查是否通过）
- 注册中心连接是否正常

### Q3: 如何实现服务灰度发布？

**A:** 使用服务版本和元数据：
- 在服务元数据中标记版本
- 客户端根据版本选择实例
- 逐步切换流量到新版本

### Q4: 服务下线时如何优雅关闭？

**A:** Kratos 会自动处理：
- 收到停止信号时，先停止接收新请求
- 等待现有请求处理完成
- 从注册中心注销服务
- 关闭服务

### Q5: 多个注册中心如何选择？

**A:** 根据场景选择：
- **Consul**：功能全面，适合复杂场景
- **etcd**：轻量级，适合 Kubernetes 环境
- **Nacos**：功能丰富，适合 Java/Go 混合技术栈
- **Kubernetes**：如果使用 K8s，可以直接使用其服务发现

## 总结

服务注册与发现是微服务架构的基础设施，Kratos 框架提供了完整的支持：

1. ✅ **服务注册**：服务启动时自动注册
2. ✅ **服务发现**：客户端自动发现可用实例
3. ✅ **健康检查**：自动检测服务健康状态
4. ✅ **负载均衡**：内置多种负载均衡策略
5. ✅ **多注册中心**：支持 Consul、etcd、Nacos 等

通过服务注册与发现，可以实现：
- 动态服务管理
- 自动负载均衡
- 故障自动转移
- 服务治理能力

## 参考资源

- [Kratos Registry 文档](https://go-kratos.dev/docs/component/registry)
- [Consul 官方文档](https://www.consul.io/docs)
- [etcd 官方文档](https://etcd.io/docs)
- [Nacos 官方文档](https://nacos.io/docs)
