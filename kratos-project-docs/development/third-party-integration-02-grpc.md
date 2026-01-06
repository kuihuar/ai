# 第三方服务集成指南 - 第二步：gRPC 服务集成

## 概述

本文档介绍如何集成 gRPC 第三方服务，包括内部微服务和外部 gRPC 服务。

## 步骤 1: 定义 Proto 文件

### 1.1 创建目录结构

```bash
mkdir -p api/external/{service-name}/v1
```

**示例：**
```bash
mkdir -p api/external/user-service/v1
```

### 1.2 编写 Proto 文件

在 `api/external/{service-name}/v1/` 目录下创建 `.proto` 文件。

**示例：`api/external/user-service/v1/user.proto`**

```protobuf
syntax = "proto3";

package user.service.v1;

option go_package = "sre/api/external/user-service/v1;v1";

// 用户服务定义
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

// 请求消息
message GetUserRequest {
  int64 id = 1;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

// 响应消息
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int64 created_at = 4;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}
```

### 1.3 生成 Go 代码

```bash
# 在项目根目录执行
protoc --proto_path=. \
       --proto_path=./third_party \
       --go_out=paths=source_relative:. \
       --go-grpc_out=paths=source_relative:. \
       api/external/user-service/v1/user.proto
```

生成的文件：
- `user.pb.go` - 消息类型
- `user_grpc.pb.go` - 服务客户端和服务器代码

## 步骤 2: 创建 gRPC 客户端管理器

### 2.1 创建客户端管理器

在 `internal/data/clients/grpc.go` 中创建 gRPC 客户端管理器：

```go
package clients

import (
	"context"
	"fmt"
	"sync"

	"sre/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClients 管理所有 gRPC 客户端连接
type GRPCClients struct {
	clients map[string]*grpc.ClientConn
	mu      sync.RWMutex
	log     *log.Helper
}

// NewGRPCClients 创建 gRPC 客户端管理器
func NewGRPCClients(c *conf.Data, logger log.Logger) (*GRPCClients, error) {
	if c.Grpc == nil || len(c.Grpc.Clients) == 0 {
		return &GRPCClients{
			clients: make(map[string]*grpc.ClientConn),
			log:     log.NewHelper(logger),
		}, nil
	}

	clients := &GRPCClients{
		clients: make(map[string]*grpc.ClientConn),
		log:     log.NewHelper(logger),
	}

	// 初始化所有配置的客户端
	for name, endpoint := range c.Grpc.Clients {
		conn, err := grpc.Dial(
			endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// 可以添加更多选项，如超时、重试等
		)
		if err != nil {
			clients.log.Warnf("failed to connect to gRPC service %s at %s: %v", name, endpoint, err)
			continue
		}
		clients.clients[name] = conn
		clients.log.Infof("gRPC client connected: %s -> %s", name, endpoint)
	}

	return clients, nil
}

// GetClient 获取指定名称的 gRPC 客户端连接
func (c *GRPCClients) GetClient(name string) (*grpc.ClientConn, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, ok := c.clients[name]
	if !ok {
		return nil, fmt.Errorf("gRPC client %s not found", name)
	}
	return conn, nil
}

// Close 关闭所有客户端连接
func (c *GRPCClients) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error
	for name, conn := range c.clients {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close %s: %w", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing gRPC clients: %v", errs)
	}
	return nil
}
```

### 2.2 更新 Data 结构体

在 `internal/data/data.go` 中添加 gRPC 客户端管理器：

```go
package data

import (
	"sre/internal/conf"
	"sre/internal/data/clients"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	clients.NewGRPCClients,  // 添加 gRPC 客户端提供者
	NewUserRepo,
)

// Data 统一管理所有数据访问依赖
type Data struct {
	db          *gorm.DB
	grpcClients *clients.GRPCClients  // 添加 gRPC 客户端
}

// NewData 创建 Data 实例
func NewData(
	db *gorm.DB,
	grpcClients *clients.GRPCClients,
	logger log.Logger,
) (*Data, func(), error) {
	cleanup := func() {
		logHelper := log.NewHelper(logger)
		logHelper.Info("closing the data resources")
		
		// 关闭数据库
		if db != nil {
			if sqlDB, err := db.DB(); err == nil {
				sqlDB.Close()
			}
		}
		
		// 关闭 gRPC 客户端
		if grpcClients != nil {
			if err := grpcClients.Close(); err != nil {
				logHelper.Warnf("failed to close gRPC clients: %v", err)
			}
		}
	}
	
	return &Data{
		db:          db,
		grpcClients: grpcClients,
	}, cleanup, nil
}
```

## 步骤 3: 在 Repository 中使用 gRPC 客户端

### 3.1 创建业务封装（可选）

在 `internal/data/external/{service}/` 目录下创建业务封装：

**示例：`internal/data/external/user-service/user.go`**

```go
package userservice

import (
	"context"

	v1 "sre/api/external/user-service/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

// Client 用户服务客户端封装
type Client struct {
	client v1.UserServiceClient
	log    *log.Helper
}

// NewClient 创建用户服务客户端
func NewClient(conn *grpc.ClientConn, logger log.Logger) *Client {
	return &Client{
		client: v1.NewUserServiceClient(conn),
		log:    log.NewHelper(logger),
	}
}

// GetUser 获取用户信息
func (c *Client) GetUser(ctx context.Context, id int64) (*v1.User, error) {
	resp, err := c.client.GetUser(ctx, &v1.GetUserRequest{Id: id})
	if err != nil {
		c.log.Errorf("failed to get user %d: %v", id, err)
		return nil, err
	}
	return resp, nil
}

// CreateUser 创建用户
func (c *Client) CreateUser(ctx context.Context, name, email string) (*v1.User, error) {
	resp, err := c.client.CreateUser(ctx, &v1.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		c.log.Errorf("failed to create user: %v", err)
		return nil, err
	}
	return resp, nil
}
```

### 3.2 在 Repository 中使用

在 Repository 中通过 Data 获取 gRPC 客户端：

```go
package data

import (
	"context"
	"sre/internal/biz"
	"sre/internal/data/external/userservice"
	v1 "sre/api/external/user-service/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data   *Data
	log    *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) GetExternalUser(ctx context.Context, id int64) (*biz.User, error) {
	// 获取 gRPC 客户端连接
	conn, err := r.data.grpcClients.GetClient("user-service")
	if err != nil {
		return nil, err
	}

	// 创建客户端封装
	client := userservice.NewClient(conn, r.log)
	
	// 调用服务
	user, err := client.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为业务类型
	return &biz.User{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
```

## 步骤 4: 更新 Wire 配置

确保在 `cmd/{app}/wire.go` 中包含 gRPC 客户端提供者：

```go
//go:build wireinject
// +build wireinject

package main

import (
	// ... 其他导入 ...
	"sre/internal/data/clients"
)

func initApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// ... 其他提供者 ...
		clients.NewGRPCClients,
		data.ProviderSet,
		// ...
	))
}
```

运行 Wire 生成代码：

```bash
cd cmd/{app}
go generate
```

## 步骤 5: 配置服务地址

在 `configs/config.yaml` 中配置服务地址：

```yaml
data:
  grpc:
    clients:
      user-service: 127.0.0.1:9001
      order-service: 127.0.0.1:9002
```

## 高级配置

### 使用 TLS 加密

```go
import (
	"crypto/tls"
	"google.golang.org/grpc/credentials"
)

conn, err := grpc.Dial(
	endpoint,
	grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false, // 生产环境应设为 false
	})),
)
```

### 使用服务发现

如果使用服务发现（如 Consul、Etcd），可以通过 resolver 自动发现服务：

```go
import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
)

// 在创建连接时使用服务发现
conn, err := grpc.Dial(
	fmt.Sprintf("discovery:///%s", serviceName),
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithResolvers(discovery.NewBuilder(r)), // r 是 registry.Registry
)
```

## 错误处理

建议定义服务特定的错误类型，参考 `api/helloworld/v1/error_reason.proto` 的模式。

## 测试

创建单元测试时，可以使用 mock 来模拟 gRPC 客户端：

```go
//go:build !integration

package data_test

import (
	"testing"
	// 使用 mock 库，如 github.com/golang/mock
)
```

## 下一步

完成 gRPC 服务集成后，可以参考：
- [第四步：在 Repository 中使用第三方服务](./third-party-integration-04-usage.md)
- [HTTP REST API 集成指南](./third-party-integration-03-http.md)

