# 配置中心接入指南

## 概述

配置中心（Config Center）是用于集中管理和动态更新分布式系统中应用配置的工具。通过配置中心，可以在不重启应用的情况下实时调整配置参数，提升系统的灵活性和可维护性。

### 为什么需要配置中心？

在微服务架构中，配置管理面临以下挑战：

- **配置分散**：配置分布在各个服务的配置文件中，难以统一管理
- **环境差异**：开发、测试、生产环境配置不同，容易出错
- **配置变更困难**：修改配置需要重新部署，影响服务可用性
- **敏感信息管理**：密码、密钥等敏感信息需要安全存储
- **配置版本管理**：需要记录配置变更历史，支持回滚

**配置中心的优势：**
- ✅ 集中管理：所有配置统一存储在配置中心
- ✅ 动态更新：配置变更实时生效，无需重启服务
- ✅ 环境隔离：通过命名空间或环境标识区分不同环境
- ✅ 版本管理：记录配置变更历史，支持回滚
- ✅ 权限控制：细粒度的配置访问权限管理
- ✅ 配置加密：敏感信息加密存储

## 架构原理

### 配置加载流程

```
┌─────────────┐
│  应用启动   │
└──────┬──────┘
       │ 1. 连接配置中心
       │ 2. 获取配置（应用名、环境、命名空间）
       ▼
┌─────────────┐
│  配置中心   │
│ (Config     │
│  Center)    │
└──────┬──────┘
       │ 3. 返回配置数据
       ▼
┌─────────────┐
│  应用加载   │
│  配置到内存 │
└──────┬──────┘
       │ 4. 监听配置变更
       │ 5. 配置变更时自动更新
       ▼
┌─────────────┐
│  配置生效   │
└─────────────┘
```

### 配置更新流程

```
┌─────────────┐
│  管理员     │
│  修改配置   │
└──────┬──────┘
       │ 1. 在配置中心修改配置
       ▼
┌─────────────┐
│  配置中心   │
│  推送变更   │
└──────┬──────┘
       │ 2. 通知所有订阅的应用
       ▼
┌─────────────┐
│  应用接收   │
│  配置变更   │
└──────┬──────┘
       │ 3. 验证配置有效性
       │ 4. 更新内存中的配置
       │ 5. 触发回调函数
       ▼
┌─────────────┐
│  配置生效   │
│  (无需重启) │
└─────────────┘
```

## Kratos 框架支持

Kratos 框架通过 `config.Source` 接口支持多种配置中心，可以轻松接入远程配置。

### 支持的配置中心

Kratos 官方支持以下配置中心：

1. **Apollo** - 携程开源的配置管理平台
2. **Nacos** - 阿里巴巴的服务发现和配置管理平台
3. **Consul** - HashiCorp 的分布式配置管理系统
4. **etcd** - 分布式键值存储系统
5. **Kubernetes ConfigMap** - Kubernetes 原生配置管理

### 核心接口

```go
// 配置源接口
type Source interface {
    Load() ([]*KeyValue, error)
    Watch() (Watcher, error)
}

// 配置键值对
type KeyValue struct {
    Key    string
    Value  []byte
    Format string  // yaml, json, toml 等
}
```

## 实现步骤

### 步骤 1: 添加依赖

根据选择的配置中心，添加对应的依赖：

#### Apollo

```bash
go get github.com/go-kratos/kratos/contrib/config/apollo/v2
```

#### Nacos

```bash
go get github.com/go-kratos/kratos/contrib/config/nacos/v2
go get github.com/nacos-group/nacos-sdk-go/v2
```

#### Consul

```bash
go get github.com/go-kratos/kratos/contrib/config/consul/v2
go get github.com/hashicorp/consul/api
```

#### etcd

```bash
go get github.com/go-kratos/kratos/contrib/config/etcd/v2
go get go.etcd.io/etcd/client/v3
```

### 步骤 2: 更新配置定义

在 `internal/conf/conf.proto` 中添加配置中心配置：

```protobuf
syntax = "proto3";
package kratos.api;

option go_package = "sre/internal/conf;conf";

import "google/protobuf/duration.proto";

message Bootstrap {
  Server server = 1;
  Data data = 2;
  Registry registry = 3;
  ConfigCenter config_center = 4;  // 新增：配置中心配置
}

message ConfigCenter {
  message Apollo {
    string app_id = 1;           // 应用 ID
    string cluster = 2;          // 集群名称
    string namespace = 3;        // 命名空间
    string ip = 4;               // Apollo 配置中心地址
    string release_key = 5;      // 发布密钥（可选）
  }
  message Nacos {
    repeated string endpoints = 1;  // Nacos 地址列表
    string namespace = 2;           // 命名空间
    string group = 3;               // 配置分组
    string data_id = 4;             // 配置 ID
    string username = 5;            // 用户名
    string password = 6;             // 密码
  }
  message Consul {
    string address = 1;          // Consul 地址
    string scheme = 2;            // 协议（http/https）
    string datacenter = 3;        // 数据中心
    string prefix = 4;            // 配置前缀
  }
  message Etcd {
    repeated string endpoints = 1;  // etcd 地址列表
    string prefix = 2;              // 配置前缀
    int64 timeout = 3;              // 超时时间（秒）
  }
  
  Apollo apollo = 1;
  Nacos nacos = 2;
  Consul consul = 3;
  Etcd etcd = 4;
}
```

重新生成配置代码：

```bash
make config
```

### 步骤 3: 创建配置中心客户端

创建 `internal/config/center.go`：

```go
package config

import (
	"sre/internal/conf"

	"github.com/go-kratos/kratos/v2/config"
	apollo "github.com/go-kratos/kratos/contrib/config/apollo/v2"
	nacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	consul "github.com/go-kratos/kratos/contrib/config/consul/v2"
	etcd "github.com/go-kratos/kratos/contrib/config/etcd/v2"

	consulAPI "github.com/hashicorp/consul/api"
	nacosClients "github.com/nacos-group/nacos-sdk-go/v2/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosVo "github.com/nacos-group/nacos-sdk-go/v2/vo"
	etcdClient "go.etcd.io/etcd/client/v3"
)

// NewConfigCenterSource 根据配置创建配置中心源
func NewConfigCenterSource(cc *conf.ConfigCenter) (config.Source, error) {
	if cc == nil {
		return nil, nil
	}

	// Apollo
	if cc.Apollo != nil && cc.Apollo.Ip != "" {
		return apollo.NewSource(
			apollo.WithAppID(cc.Apollo.AppId),
			apollo.WithCluster(cc.Apollo.Cluster),
			apollo.WithNamespace(cc.Apollo.Namespace),
			apollo.WithIP(cc.Apollo.Ip),
		)
	}

	// Nacos
	if cc.Nacos != nil && len(cc.Nacos.Endpoints) > 0 {
		sc := []nacosConstant.ServerConfig{
			*nacosConstant.NewServerConfig(cc.Nacos.Endpoints[0], 8848),
		}
		cc := nacosConstant.ClientConfig{
			NamespaceId:         cc.Nacos.Namespace,
			Username:            cc.Nacos.Username,
			Password:            cc.Nacos.Password,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
		}
		nacosClient, err := nacosClients.NewConfigClient(
			nacosVo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
		if err != nil {
			return nil, err
		}
		return nacos.NewConfigSource(nacosClient, nacos.WithGroup(cc.Nacos.Group), nacos.WithDataID(cc.Nacos.DataId)), nil
	}

	// Consul
	if cc.Consul != nil && cc.Consul.Address != "" {
		consulConfig := consulAPI.DefaultConfig()
		consulConfig.Address = cc.Consul.Address
		if cc.Consul.Scheme != "" {
			consulConfig.Scheme = cc.Consul.Scheme
		}
		if cc.Consul.Datacenter != "" {
			consulConfig.Datacenter = cc.Consul.Datacenter
		}
		client, err := consulAPI.NewClient(consulConfig)
		if err != nil {
			return nil, err
		}
		return consul.New(client, consul.WithPrefix(cc.Consul.Prefix)), nil
	}

	// etcd
	if cc.Etcd != nil && len(cc.Etcd.Endpoints) > 0 {
		etcdClient, err := etcdClient.New(etcdClient.Config{
			Endpoints: cc.Etcd.Endpoints,
		})
		if err != nil {
			return nil, err
		}
		return etcd.New(etcdClient, etcd.WithPrefix(cc.Etcd.Prefix)), nil
	}

	return nil, nil
}
```

### 步骤 4: 更新 main.go

修改 `cmd/sre/main.go`，添加配置中心支持：

```go
package main

import (
	"flag"
	"os"

	"sre/internal/conf"
	"sre/internal/config"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	kratosConfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// ... existing code ...

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

	// ============================================
	// 步骤 1: 先加载本地配置文件（获取元配置）
	// ============================================
	// 配置中心地址必须放在本地配置文件或环境变量中
	// 这是"元配置"，用于连接配置中心
	tempConfig := kratosConfig.New(
		kratosConfig.WithSource(file.NewSource(flagconf)),
	)
	defer tempConfig.Close()
	
	if err := tempConfig.Load(); err != nil {
		panic(err)
	}
	
	var bootstrap conf.Bootstrap
	if err := tempConfig.Scan(&bootstrap); err != nil {
		panic(err)
	}

	// ============================================
	// 步骤 2: 创建配置源列表（本地文件 + 配置中心）
	// ============================================
	sources := []kratosConfig.Source{
		file.NewSource(flagconf), // 本地配置文件（兜底配置）
	}

	// 如果配置了配置中心，添加配置中心源
	// 配置中心源放在前面，优先级更高（会覆盖本地配置）
	if bootstrap.ConfigCenter != nil {
		configCenterSource, err := config.NewConfigCenterSource(bootstrap.ConfigCenter)
		if err != nil {
			logger.Log(log.LevelWarn, "failed to connect config center, using local config only", "err", err)
		} else if configCenterSource != nil {
			// 配置中心源放在前面，优先级更高
			sources = append([]kratosConfig.Source{configCenterSource}, sources...)
			logger.Log(log.LevelInfo, "config center connected", "type", getConfigCenterType(bootstrap.ConfigCenter))
		}
	}

	// ============================================
	// 步骤 3: 创建最终配置对象（合并本地和远程配置）
	// ============================================
	c := kratosConfig.New(
		kratosConfig.WithSource(sources...),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	// 重新扫描配置（合并后的配置：配置中心 + 本地文件）
	if err := c.Scan(&bootstrap); err != nil {
		panic(err)
	}

	// ============================================
	// 步骤 4: 监听配置变更（热更新）
	// ============================================
	watcher, err := c.Watch()
	if err == nil {
		go func() {
			for {
				values, err := watcher.Next()
				if err != nil {
					logger.Log(log.LevelError, "watch config error", "err", err)
					break
				}
				logger.Log(log.LevelInfo, "config changed", "values", values)
				// 重新扫描配置
				if err := c.Scan(&bootstrap); err != nil {
					logger.Log(log.LevelError, "scan config error", "err", err)
				}
			}
		}()
	}

	// ============================================
	// 步骤 5: 启动应用
	// ============================================
	app, cleanup, err := wireApp(bootstrap.Server, bootstrap.Data, bootstrap.Registry, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

// getConfigCenterType 获取配置中心类型（用于日志）
func getConfigCenterType(cc *conf.ConfigCenter) string {
	if cc == nil {
		return "none"
	}
	if cc.Apollo != nil {
		return "apollo"
	}
	if cc.Nacos != nil {
		return "nacos"
	}
	if cc.Consul != nil {
		return "consul"
	}
	if cc.Etcd != nil {
		return "etcd"
	}
	return "unknown"
}
```

### 步骤 5: 配置文件

在 `configs/config.yaml` 中添加配置中心配置：

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s

data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s

registry:
  etcd:
    endpoints:
      - "127.0.0.1:2379"
    timeout: 5

# 配置中心配置（示例：使用 Nacos）
config_center:
  nacos:
    endpoints:
      - "127.0.0.1:8848"
    namespace: "public"
    group: "DEFAULT_GROUP"
    data_id: "sre-config.yaml"
    username: "nacos"
    password: "nacos"
```

### 步骤 6: 更新 wire.go

在 `cmd/sre/wire.go` 中添加配置中心依赖（如果需要）：

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

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		registry.ProviderSet,
		newApp,
	))
}
```

运行 Wire 生成代码：

```bash
make wire
```

## 各配置中心详细配置

### Apollo 配置

#### 1. 安装 Apollo

参考 [Apollo 官方文档](https://www.apolloconfig.com/) 安装 Apollo 配置中心。

#### 2. 配置示例

```yaml
config_center:
  apollo:
    app_id: "sre"                    # 应用 ID
    cluster: "default"               # 集群名称
    namespace: "application"         # 命名空间
    ip: "http://127.0.0.1:8080"     # Apollo 配置中心地址
```

#### 3. Apollo 配置格式

在 Apollo 控制台创建配置，格式为 YAML：

```yaml
server:
  http:
    addr: 0.0.0.0:8000
data:
  database:
    source: "user:pass@tcp(host:3306)/db"
```

### Nacos 配置

#### 1. 安装 Nacos

```bash
# 下载 Nacos
wget https://github.com/alibaba/nacos/releases/download/2.3.0/nacos-server-2.3.0.tar.gz
tar -xzf nacos-server-2.3.0.tar.gz
cd nacos/bin

# 启动 Nacos（单机模式）
sh startup.sh -m standalone
```

#### 2. 配置示例

```yaml
config_center:
  nacos:
    endpoints:
      - "127.0.0.1:8848"
    namespace: "public"              # 命名空间 ID
    group: "DEFAULT_GROUP"          # 配置分组
    data_id: "sre-config.yaml"      # 配置 ID（Data ID）
    username: "nacos"                # 用户名
    password: "nacos"                # 密码
```

#### 3. Nacos 配置格式

在 Nacos 控制台（`http://127.0.0.1:8848/nacos`）创建配置：

- **Data ID**: `sre-config.yaml`
- **Group**: `DEFAULT_GROUP`
- **配置格式**: `YAML`
- **配置内容**:

```yaml
server:
  http:
    addr: 0.0.0.0:8000
data:
  database:
    source: "user:pass@tcp(host:3306)/db"
```

### Consul 配置

#### 1. 安装 Consul

```bash
# 下载 Consul
wget https://releases.hashicorp.com/consul/1.17.0/consul_1.17.0_linux_amd64.zip
unzip consul_1.17.0_linux_amd64.zip

# 启动 Consul（开发模式）
./consul agent -dev
```

#### 2. 配置示例

```yaml
config_center:
  consul:
    address: "127.0.0.1:8500"       # Consul 地址
    scheme: "http"                   # 协议
    datacenter: "dc1"                # 数据中心
    prefix: "sre/config"             # 配置前缀（KV 路径）
```

#### 3. Consul 配置格式

使用 Consul KV API 或 Web UI 存储配置：

```bash
# 使用 Consul CLI 设置配置
consul kv put sre/config/server.http.addr "0.0.0.0:8000"
consul kv put sre/config/data.database.source "user:pass@tcp(host:3306)/db"

# 或者使用 JSON 格式存储整个配置
consul kv put sre/config @config.json
```

### etcd 配置

#### 1. 安装 etcd

```bash
# 下载 etcd
wget https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-linux-amd64.tar.gz
tar -xzf etcd-v3.5.9-linux-amd64.tar.gz
cd etcd-v3.5.9-linux-amd64

# 启动 etcd
./etcd
```

#### 2. 配置示例

```yaml
config_center:
  etcd:
    endpoints:
      - "127.0.0.1:2379"
    prefix: "sre/config"             # 配置前缀
    timeout: 5                        # 超时时间（秒）
```

#### 3. etcd 配置格式

使用 etcdctl 设置配置：

```bash
# 设置单个配置项
etcdctl put sre/config/server.http.addr "0.0.0.0:8000"
etcdctl put sre/config/data.database.source "user:pass@tcp(host:3306)/db"

# 或者使用 YAML 格式存储整个配置
etcdctl put sre/config @config.yaml
```

## 配置中心地址获取与配置分层

### 核心问题：配置中心地址从哪里来？

在微服务架构中，存在一个"鸡生蛋"的问题：
- **问题**：需要配置中心地址才能连接配置中心，但配置中心地址本身也是配置
- **解决方案**：采用**配置分层策略**，将配置分为"元配置"和"业务配置"

### 配置分层策略

#### 分层原则

```
┌─────────────────────────────────────┐
│  第一层：元配置（Bootstrap Config） │
│  - 配置中心地址                     │
│  - 服务注册中心地址                 │
│  - 应用基础信息（名称、版本等）      │
│  存储位置：本地文件 / 环境变量       │
└─────────────────────────────────────┘
              │
              │ 使用元配置连接
              ▼
┌─────────────────────────────────────┐
│  第二层：业务配置（Business Config） │
│  - 数据库连接信息                    │
│  - Redis 配置                       │
│  - 业务参数                         │
│  - 第三方服务配置                   │
│  存储位置：配置中心                  │
└─────────────────────────────────────┘
```

#### 配置分类

| 配置类型 | 存储位置 | 变更频率 | 示例 |
|---------|---------|---------|------|
| **元配置** | 本地文件/环境变量 | 低 | 配置中心地址、注册中心地址 |
| **基础配置** | 本地文件 | 低 | 服务端口、日志级别 |
| **环境配置** | 配置中心 | 中 | 数据库连接、Redis 地址 |
| **业务配置** | 配置中心 | 高 | 业务参数、开关配置 |
| **敏感配置** | 配置中心（加密） | 低 | 密码、密钥、Token |

### 实现方式

#### 方式 1: 本地配置文件（推荐）

**配置中心地址放在本地配置文件**，这是最简单和常用的方式。

**配置文件结构** (`configs/config.yaml`):

```yaml
# 元配置：配置中心地址（必须放在本地文件）
config_center:
  nacos:
    endpoints:
      - "127.0.0.1:8848"    # 配置中心地址
    namespace: "public"
    group: "DEFAULT_GROUP"
    data_id: "sre-config.yaml"
    username: "nacos"
    password: "nacos"

# 基础配置：服务基础信息（也可以放在本地）
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s

# 注册中心地址（元配置）
registry:
  etcd:
    endpoints:
      - "127.0.0.1:2379"
    timeout: 5
```

**加载流程**:

```go
// 1. 先加载本地配置文件（包含配置中心地址）
c := config.New(
    config.WithSource(
        file.NewSource(flagconf), // 本地文件
    ),
)
defer c.Close()

if err := c.Load(); err != nil {
    panic(err)
}

// 2. 读取配置中心地址
var bc conf.Bootstrap
if err := c.Scan(&bc); err != nil {
    panic(err)
}

// 3. 如果配置了配置中心，添加配置中心源
sources := []kratosConfig.Source{
    file.NewSource(flagconf), // 本地文件（兜底）
}

if bc.ConfigCenter != nil {
    configCenterSource, err := config.NewConfigCenterSource(bc.ConfigCenter)
    if err == nil && configCenterSource != nil {
        // 配置中心源放在前面，优先级更高
        sources = append([]kratosConfig.Source{configCenterSource}, sources...)
    }
}

// 4. 重新创建配置对象（本地文件 + 配置中心）
c = config.New(
    config.WithSource(sources...),
)
if err := c.Load(); err != nil {
    panic(err)
}

// 5. 重新扫描配置（合并本地和远程配置）
if err := c.Scan(&bc); err != nil {
    panic(err)
}
```

#### 方式 2: 环境变量（推荐用于容器化部署）

**配置中心地址通过环境变量传递**，适合容器化部署场景。

**环境变量设置**:

```bash
# 配置中心地址
export CONFIG_CENTER_NACOS_ENDPOINTS="127.0.0.1:8848"
export CONFIG_CENTER_NACOS_NAMESPACE="public"
export CONFIG_CENTER_NACOS_GROUP="DEFAULT_GROUP"
export CONFIG_CENTER_NACOS_DATA_ID="sre-config.yaml"
export CONFIG_CENTER_NACOS_USERNAME="nacos"
export CONFIG_CENTER_NACOS_PASSWORD="nacos"
```

**代码实现**:

```go
// 从环境变量读取配置中心地址
func getConfigCenterFromEnv() *conf.ConfigCenter {
    endpoints := os.Getenv("CONFIG_CENTER_NACOS_ENDPOINTS")
    if endpoints == "" {
        return nil
    }
    
    return &conf.ConfigCenter{
        Nacos: &conf.ConfigCenter_Nacos{
            Endpoints: []string{endpoints},
            Namespace: os.Getenv("CONFIG_CENTER_NACOS_NAMESPACE"),
            Group:     os.Getenv("CONFIG_CENTER_NACOS_GROUP"),
            DataId:    os.Getenv("CONFIG_CENTER_NACOS_DATA_ID"),
            Username:  os.Getenv("CONFIG_CENTER_NACOS_USERNAME"),
            Password:  os.Getenv("CONFIG_CENTER_NACOS_PASSWORD"),
        },
    }
}

func main() {
    // 1. 从环境变量获取配置中心地址
    configCenter := getConfigCenterFromEnv()
    
    // 2. 如果环境变量没有，则从本地文件读取
    if configCenter == nil {
        c := config.New(
            config.WithSource(file.NewSource(flagconf)),
        )
        if err := c.Load(); err == nil {
            var bc conf.Bootstrap
            if err := c.Scan(&bc); err == nil {
                configCenter = bc.ConfigCenter
            }
        }
        c.Close()
    }
    
    // 3. 创建配置源列表
    sources := []kratosConfig.Source{
        file.NewSource(flagconf), // 本地文件（兜底）
    }
    
    // 4. 添加配置中心源
    if configCenter != nil {
        configCenterSource, err := config.NewConfigCenterSource(configCenter)
        if err == nil && configCenterSource != nil {
            sources = append([]kratosConfig.Source{configCenterSource}, sources...)
        }
    }
    
    // 5. 加载配置
    c := config.New(
        config.WithSource(sources...),
    )
    // ...
}
```

#### 方式 3: 命令行参数

**配置中心地址通过命令行参数传递**。

```go
var (
    flagconf          string
    flagConfigCenter  string  // 配置中心地址
)

func init() {
    flag.StringVar(&flagconf, "conf", "../../configs", "config path")
    flag.StringVar(&flagConfigCenter, "config-center", "", "config center address, eg: nacos://127.0.0.1:8848")
}

func main() {
    flag.Parse()
    
    // 解析配置中心地址
    var configCenter *conf.ConfigCenter
    if flagConfigCenter != "" {
        configCenter = parseConfigCenterAddress(flagConfigCenter)
    }
    
    // ... 后续流程
}
```

#### 方式 4: 服务发现（高级）

**通过服务注册中心发现配置中心地址**，适用于大规模微服务架构。

```go
// 1. 先连接服务注册中心（地址在本地配置或环境变量）
registryClient := connectRegistry(localConfig.Registry)

// 2. 从服务注册中心发现配置中心服务
configCenterInstances := registryClient.Discover("config-center")

// 3. 选择可用的配置中心实例
configCenterAddr := selectInstance(configCenterInstances)

// 4. 连接配置中心
configCenterSource := connectConfigCenter(configCenterAddr)
```

### 最佳实践总结

#### ✅ 推荐做法

1. **配置中心地址放在本地配置文件**
   - 简单直接，易于理解
   - 适合大多数场景
   - 便于版本控制和部署

2. **容器化部署使用环境变量**
   - 不同环境使用不同的环境变量
   - 避免将敏感信息提交到代码仓库
   - 符合 12-Factor App 原则

3. **配置分层明确**
   - 元配置（配置中心地址）→ 本地文件/环境变量
   - 业务配置（数据库、Redis 等）→ 配置中心
   - 基础配置（服务端口等）→ 本地文件

4. **提供兜底机制**
   - 配置中心不可用时，使用本地配置文件
   - 本地配置文件包含完整的默认配置

#### ❌ 不推荐做法

1. **配置中心地址放在配置中心**
   - 会导致循环依赖问题
   - 无法启动应用

2. **所有配置都放在配置中心**
   - 配置中心地址等元配置必须本地化
   - 基础配置（如服务端口）建议本地化

3. **硬编码配置中心地址**
   - 不利于多环境部署
   - 难以维护

### 配置文件示例

**完整的配置文件结构** (`configs/config.yaml`):

```yaml
# ============================================
# 元配置：必须放在本地文件
# ============================================

# 配置中心地址（元配置）
config_center:
  nacos:
    endpoints:
      - "127.0.0.1:8848"
    namespace: "public"
    group: "DEFAULT_GROUP"
    data_id: "sre-config.yaml"
    username: "nacos"
    password: "nacos"

# 服务注册中心地址（元配置）
registry:
  etcd:
    endpoints:
      - "127.0.0.1:2379"
    timeout: 5

# ============================================
# 基础配置：建议放在本地文件
# ============================================

# 服务配置（基础配置）
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s

# 日志配置（基础配置）
log:
  level: "info"
  format: "json"

# ============================================
# 以下配置建议放在配置中心
# ============================================

# 数据源配置（环境配置，建议放在配置中心）
# data:
#   database:
#     driver: mysql
#     source: "user:pass@tcp(host:3306)/db"
#   redis:
#     addr: "127.0.0.1:6379"

# 业务配置（建议放在配置中心）
# business:
#   feature_flags:
#     enable_new_feature: true
#   rate_limit:
#     qps: 1000
```

**配置中心中的配置** (`sre-config.yaml` in Nacos):

```yaml
# 数据源配置（环境配置）
data:
  database:
    driver: mysql
    source: "user:pass@tcp(host:3306)/db"
  redis:
    addr: "127.0.0.1:6379"
    read_timeout: 0.2s
    write_timeout: 0.2s
  kafka:
    brokers:
      - "127.0.0.1:9092"
    topic: "user-create-events"
    group_id: "user-consumer-group"

# 业务配置
business:
  feature_flags:
    enable_new_feature: true
  rate_limit:
    qps: 1000
  timeout:
    api_timeout: 5s
```

### 配置加载流程图

```
应用启动
    │
    ├─► 1. 加载本地配置文件
    │   └─► 获取元配置（配置中心地址）
    │
    ├─► 2. 连接配置中心
    │   └─► 使用元配置中的地址
    │
    ├─► 3. 从配置中心获取业务配置
    │   └─► 数据库、Redis、业务参数等
    │
    ├─► 4. 合并配置
    │   ├─► 配置中心配置（优先级高）
    │   └─► 本地配置文件（兜底）
    │
    └─► 5. 应用启动完成
        └─► 监听配置变更
```

## 配置优先级

配置加载的优先级（从高到低）：

1. **配置中心** - 远程配置，优先级最高
2. **本地配置文件** - 作为默认配置和兜底方案
3. **环境变量** - 如果使用 Viper 配置系统

**配置合并规则：**
- 配置中心的值会覆盖本地配置文件中的相同配置项
- 未在配置中心配置的项，使用本地配置文件的值
- 支持部分配置在配置中心，部分配置在本地文件

## 配置热更新

配置中心支持配置热更新，无需重启服务即可生效。

### 实现方式

在 `main.go` 中已经添加了配置监听：

```go
// 监听配置变更
watcher, err := c.Watch()
if err == nil {
    go func() {
        for {
            values, err := watcher.Next()
            if err != nil {
                logger.Log(log.LevelError, "watch config error", "err", err)
                break
            }
            logger.Log(log.LevelInfo, "config changed", "values", values)
            // 重新扫描配置
            if err := c.Scan(&bc); err != nil {
                logger.Log(log.LevelError, "scan config error", "err", err)
            }
        }
    }()
}
```

### 配置更新回调

如果需要根据配置变更执行特定操作，可以添加回调函数：

```go
watcher, err := c.Watch()
if err == nil {
    go func() {
        for {
            values, err := watcher.Next()
            if err != nil {
                break
            }
            
            // 重新扫描配置
            var newBc conf.Bootstrap
            if err := c.Scan(&newBc); err != nil {
                continue
            }
            
            // 检查特定配置是否变更
            if bc.Server.Http.Addr != newBc.Server.Http.Addr {
                logger.Log(log.LevelInfo, "server address changed", 
                    "old", bc.Server.Http.Addr, 
                    "new", newBc.Server.Http.Addr)
                // 执行相应的更新操作
            }
            
            // 更新配置
            bc = newBc
        }
    }()
}
```

## 最佳实践

### 1. 配置分层

- **基础配置**：放在本地配置文件（如服务端口、日志级别）
- **环境配置**：放在配置中心（如数据库连接、Redis 地址）
- **敏感配置**：放在配置中心，并启用加密（如密码、密钥）

### 2. 配置命名规范

- 使用有意义的配置键名
- 使用点分隔符组织层级（如 `server.http.addr`）
- 统一配置前缀（如 `sre/config`）

### 3. 配置验证

- 启动时验证配置完整性
- 配置变更时验证配置有效性
- 提供清晰的错误提示

### 4. 配置版本管理

- 记录配置变更历史
- 支持配置回滚
- 配置变更通知相关人员

### 5. 高可用性

- 配置中心应具备高可用性（集群部署）
- 配置中心不可用时，使用本地配置文件作为兜底
- 实现配置缓存，减少对配置中心的依赖

### 6. 安全性

- 敏感信息加密存储
- 配置访问权限控制
- 配置传输使用 HTTPS/TLS

### 7. 监控和告警

- 监控配置中心连接状态
- 监控配置变更频率
- 配置异常时及时告警

## 常见问题

### Q1: 配置中心不可用时怎么办？

**A:** 配置加载时会先尝试从配置中心获取配置，如果失败，会使用本地配置文件作为兜底。建议：
- 本地配置文件包含完整的默认配置
- 配置中心恢复后，配置会自动更新

### Q2: 如何实现配置的灰度发布？

**A:** 可以通过以下方式实现：
- **Apollo**: 使用灰度发布功能
- **Nacos**: 使用配置分组和命名空间
- **Consul/etcd**: 通过不同的配置前缀区分环境

### Q3: 配置变更会影响正在处理的请求吗？

**A:** 配置变更通常不会影响正在处理的请求，但需要注意：
- 数据库连接池等资源需要重新初始化
- 某些配置变更可能需要重启服务才能生效（如服务端口）

### Q4: 如何管理多个环境的配置？

**A:** 推荐方式：
- **方式 1**: 使用不同的命名空间（如 `dev`、`test`、`prod`）
- **方式 2**: 使用不同的配置分组（如 `DEFAULT_GROUP`、`PROD_GROUP`）
- **方式 3**: 使用不同的配置前缀（如 `sre/config/dev`、`sre/config/prod`）

### Q5: 配置中心和服务注册中心可以共用吗？

**A:** 可以。Nacos、Consul、etcd 都同时支持服务注册发现和配置管理：
- **Nacos**: 同时支持服务注册发现和配置管理
- **Consul**: 同时支持服务注册发现和 KV 存储（配置）
- **etcd**: 同时支持服务注册发现和配置存储

## 参考资源

- [Kratos 配置文档](https://go-kratos.dev/docs/component/config)
- [Apollo 官方文档](https://www.apolloconfig.com/)
- [Nacos 官方文档](https://nacos.io/docs)
- [Consul 官方文档](https://www.consul.io/docs)
- [etcd 官方文档](https://etcd.io/docs)

