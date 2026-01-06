# Kratos 多应用支持

## 概述

Kratos **支持多应用**架构。在 Kratos 项目中，`cmd/` 目录就是为多应用设计的，每个应用可以有独立的入口、配置和服务。

## 多应用架构

### 项目结构

```
sre/
├── api/                    # 共享的 API 定义
│   ├── user/
│   │   └── v1/
│   └── order/
│       └── v1/
├── cmd/                    # 应用入口目录
│   ├── api-server/        # API 网关服务
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── user-service/      # 用户服务
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── order-service/     # 订单服务
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── internal/              # 共享的内部代码
│   ├── biz/
│   ├── data/
│   └── service/
├── configs/               # 配置文件
│   ├── api-server.yaml
│   ├── user-service.yaml
│   └── order-service.yaml
└── go.mod
```

## 多应用实现方式

### 方式一：共享代码库（Monorepo）

所有应用共享 `internal/` 目录下的代码，但每个应用有独立的入口和配置。

#### 优点
- 代码复用：共享业务逻辑和数据访问层
- 统一管理：所有服务在一个仓库中
- 类型安全：共享类型定义，避免不一致

#### 缺点
- 耦合风险：需要严格控制依赖关系
- 构建复杂：需要为每个应用单独构建

#### 示例：创建新应用

```bash
# 创建新的应用目录
mkdir -p cmd/user-service

# 创建 main.go
cat > cmd/user-service/main.go << 'EOF'
package main

import (
    "flag"
    "os"
    "user-service/internal/conf"
    // ... 其他导入
)

var (
    Name    string
    Version string
    flagconf string
    id, _ = os.Hostname()
)

func init() {
    flag.StringVar(&flagconf, "conf", "../../configs", "config path")
}

func main() {
    // 应用启动逻辑
}
EOF
```

### 方式二：独立服务（推荐用于大型项目）

每个服务作为独立的模块，通过 API 通信。

#### 优点
- 独立部署：每个服务可以独立部署和扩展
- 技术栈灵活：不同服务可以使用不同技术
- 团队独立：不同团队可以独立开发

#### 缺点
- 代码重复：可能需要在多个服务中重复实现
- 通信开销：服务间需要网络通信
- 分布式复杂性：需要处理分布式系统的问题

## 多应用配置管理

### 独立配置文件

每个应用使用独立的配置文件：

```yaml
# configs/user-service.yaml
server:
  http:
    addr: 0.0.0.0:8001
  grpc:
    addr: 0.0.0.0:9001
data:
  database:
    source: "user_db_connection_string"

# configs/order-service.yaml
server:
  http:
    addr: 0.0.0.0:8002
  grpc:
    addr: 0.0.0.0:9002
data:
  database:
    source: "order_db_connection_string"
```

### 共享配置

使用配置继承或环境变量：

```yaml
# configs/base.yaml - 共享配置
server:
  http:
    timeout: 1s
  grpc:
    timeout: 1s

# configs/user-service.yaml - 应用特定配置
server:
  http:
    addr: 0.0.0.0:8001
  # 继承 base.yaml 的其他配置
```

## Wire 依赖注入

每个应用有独立的 Wire 配置：

```go
// cmd/user-service/wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "user-service/internal/biz"
    "user-service/internal/data"
    "user-service/internal/server"
    "user-service/internal/service"
    
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        newApp,
    ))
}
```

## 服务间通信

### gRPC 调用

服务间通过 gRPC 进行通信：

```go
// internal/data/user_client.go
type userClient struct {
    conn *grpc.ClientConn
}

func NewUserClient(conn *grpc.ClientConn) *userClient {
    return &userClient{conn: conn}
}

func (c *userClient) GetUser(ctx context.Context, id int64) (*User, error) {
    client := v1.NewUserServiceClient(c.conn)
    resp, err := client.GetUser(ctx, &v1.GetUserRequest{Id: id})
    if err != nil {
        return nil, err
    }
    return toUser(resp), nil
}
```

### HTTP 调用

也可以使用 HTTP 进行服务间通信：

```go
import "github.com/go-kratos/kratos/v2/transport/http"

func NewUserHTTPClient(endpoint string) *userHTTPClient {
    conn, _ := http.NewClient(context.Background(),
        http.WithEndpoint(endpoint),
    )
    return &userHTTPClient{client: conn}
}
```

## 构建和部署

### 构建单个应用

```bash
# 构建 user-service
go build -o bin/user-service ./cmd/user-service

# 构建 order-service
go build -o bin/order-service ./cmd/order-service
```

### Makefile 支持

```makefile
# Makefile
.PHONY: build-user build-order build-all

build-user:
	go build -o bin/user-service ./cmd/user-service

build-order:
	go build -o bin/order-service ./cmd/order-service

build-all: build-user build-order
	@echo "All services built successfully"
```

### Docker 多阶段构建

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download

# 构建 user-service
RUN go build -o /bin/user-service ./cmd/user-service

# 构建 order-service
RUN go build -o /bin/order-service ./cmd/order-service

FROM alpine:latest
COPY --from=builder /bin/user-service /bin/user-service
COPY --from=builder /bin/order-service /bin/order-service
```

## 最佳实践

### 1. 代码组织
- **共享代码**：将通用代码放在 `internal/` 目录
- **应用特定代码**：应用特定的代码放在各自的 `cmd/` 目录
- **API 定义**：共享的 API 定义放在 `api/` 目录

### 2. 依赖管理
- **避免循环依赖**：确保依赖方向清晰
- **接口抽象**：使用接口解耦服务间依赖
- **版本控制**：API 变更时使用版本号

### 3. 配置管理
- **环境分离**：不同环境使用不同配置
- **敏感信息**：敏感信息使用环境变量或密钥管理
- **配置验证**：启动时验证配置完整性

### 4. 服务发现
- **服务注册**：使用服务注册中心（如 Consul、etcd）
- **负载均衡**：使用客户端或服务端负载均衡
- **健康检查**：实现健康检查接口

> 📖 **详细文档**：关于服务注册与发现的完整实现指南，请参考 [服务注册与发现文档](./service-registry-discovery.md)

## 常见场景

### 场景一：API 网关 + 微服务
```
cmd/
├── gateway/        # API 网关，统一入口
├── user-service/   # 用户服务
├── order-service/  # 订单服务
└── payment-service/ # 支付服务
```

### 场景二：管理后台 + 业务服务
```
cmd/
├── admin-api/      # 管理后台 API
├── user-api/       # 用户端 API
└── worker/         # 后台任务服务
```

### 场景三：多环境部署
```
cmd/
├── prod-service/   # 生产环境服务
├── staging-service/ # 预发布环境服务
└── dev-service/    # 开发环境服务
```

## 独立部署场景的最佳实践

当多个应用需要**独立部署**时（每个应用部署到不同的服务器或容器），需要特别注意代码结构的组织。以下是针对独立部署场景的最佳实践：

### 1. 代码结构优化

#### 使用 pkg/ 目录存放共享代码

对于需要独立部署的应用，建议将共享代码放在 `pkg/` 目录，而不是 `internal/`：

```
sre/
├── api/                    # 共享的 API 定义
│   ├── user/v1/
│   └── order/v1/
├── pkg/                    # 可共享的公共代码（推荐用于独立部署）
│   ├── errors/            # 错误定义
│   ├── middleware/        # 中间件
│   ├── utils/             # 工具函数
│   ├── logger/            # 日志封装
│   └── validator/         # 验证器
├── internal/              # 应用特定的内部代码（可选）
│   ├── biz/               # 如果业务逻辑不共享，放在各自应用目录
│   ├── data/
│   └── service/
├── cmd/
│   ├── user-service/
│   │   ├── internal/      # 应用特定的内部代码
│   │   │   ├── biz/
│   │   │   ├── data/
│   │   │   └── service/
│   │   ├── main.go
│   │   └── wire.go
│   └── order-service/
│       ├── internal/      # 应用特定的内部代码
│       ├── main.go
│       └── wire.go
└── go.mod
```

**关键区别**：
- `pkg/`：可被外部导入的公共代码，适合独立部署场景
- `internal/`：项目内部代码，Go 编译器会阻止外部包导入

#### 应用特定的 internal/ 目录

每个应用可以有自己独立的 `internal/` 目录：

```go
// cmd/user-service/internal/biz/user.go
package biz

import (
    "sre/pkg/errors"  // 使用共享的错误定义
    "sre/api/user/v1" // 使用共享的 API 定义
)

type UserUseCase struct {
    // 业务逻辑
}
```

### 2. Go Workspace 支持（Go 1.18+）

对于大型项目，可以使用 Go Workspace 管理多个模块：

```
sre/
├── go.work                # Workspace 配置文件
├── pkg/                   # 共享包模块
│   ├── go.mod
│   └── errors/
├── api/                   # API 定义模块
│   ├── go.mod
│   └── user/
├── cmd/
│   ├── user-service/      # 用户服务模块
│   │   ├── go.mod
│   │   └── main.go
│   └── order-service/     # 订单服务模块
│       ├── go.mod
│       └── main.go
```

**go.work 配置**：

```go
// go.work
go 1.21

use (
    ./pkg
    ./api
    ./cmd/user-service
    ./cmd/order-service
)
```

**优势**：
- 每个服务可以独立管理依赖版本
- 支持本地开发时的模块替换
- 构建时可以选择性构建特定模块

### 3. 构建优化

#### 独立构建配置

为每个应用创建独立的构建脚本：

```makefile
# Makefile
.PHONY: build-user build-order

# 构建用户服务
build-user:
	@echo "Building user-service..."
	@mkdir -p bin
	@go build -ldflags "-X main.Version=$(VERSION) -X main.Name=user-service" \
		-o bin/user-service ./cmd/user-service

# 构建订单服务
build-order:
	@echo "Building order-service..."
	@mkdir -p bin
	@go build -ldflags "-X main.Version=$(VERSION) -X main.Name=order-service" \
		-o bin/order-service ./cmd/order-service

# 构建所有服务
build-all: build-user build-order
	@echo "All services built successfully"
```

#### 构建标签（Build Tags）

使用构建标签控制编译内容，减少构建体积：

```go
// cmd/user-service/main.go
//go:build !order_service
// +build !order_service

package main

// 用户服务的代码
```

```go
// cmd/order-service/main.go
//go:build !user_service
// +build !user_service

package main

// 订单服务的代码
```

#### 最小化依赖

每个应用只导入需要的依赖：

```go
// cmd/user-service/main.go
import (
    // 只导入用户服务需要的包
    "sre/pkg/errors"
    "sre/api/user/v1"
    // 不导入 order 相关的包
)
```

### 4. 依赖管理策略

#### 版本锁定

为每个应用独立管理依赖版本：

```go
// cmd/user-service/go.mod
module sre/cmd/user-service

go 1.21

require (
    sre/pkg v0.0.0
    sre/api v0.0.0
    github.com/go-kratos/kratos/v2 v2.8.0
)

replace (
    sre/pkg => ../../pkg
    sre/api => ../../api
)
```

#### 依赖分离

将共享依赖和特定依赖分离：

```go
// pkg/go.mod - 共享包的最小依赖
module sre/pkg

go 1.21

require (
    github.com/go-kratos/kratos/v2 v2.8.0
    // 只包含共享包需要的依赖
)

// cmd/user-service/go.mod - 应用特定依赖
module sre/cmd/user-service

require (
    sre/pkg v0.0.0
    github.com/go-redis/redis/v8 v8.11.5  // 用户服务特有的依赖
)
```

### 5. 独立 Dockerfile

为每个应用创建独立的 Dockerfile：

```dockerfile
# cmd/user-service/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY api/ ./api/
COPY cmd/user-service/ ./cmd/user-service/

# 构建应用
RUN cd cmd/user-service && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/cmd/user-service/user-service .
COPY --from=builder /app/configs/user-service.yaml ./configs/

EXPOSE 8001 9001

CMD ["./user-service", "-conf", "./configs"]
```

#### Docker Compose 多服务

```yaml
# docker-compose.yml
version: '3.8'

services:
  user-service:
    build:
      context: .
      dockerfile: cmd/user-service/Dockerfile
    ports:
      - "8001:8001"
      - "9001:9001"
    volumes:
      - ./configs/user-service.yaml:/app/configs/config.yaml

  order-service:
    build:
      context: .
      dockerfile: cmd/order-service/Dockerfile
    ports:
      - "8002:8002"
      - "9002:9002"
    volumes:
      - ./configs/order-service.yaml:/app/configs/config.yaml
```

### 6. CI/CD 优化

#### 独立构建流水线

为每个应用创建独立的 CI/CD 流水线：

```yaml
# .github/workflows/user-service.yml
name: Build User Service

on:
  push:
    paths:
      - 'cmd/user-service/**'
      - 'pkg/**'
      - 'api/user/**'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: |
          cd cmd/user-service
          go build -o user-service .
      - name: Docker Build
        run: |
          docker build -f cmd/user-service/Dockerfile -t user-service:${{ github.sha }} .
```

#### 构建缓存优化

使用构建缓存加速构建：

```dockerfile
# 优化后的 Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY api/go.mod api/go.sum ./api/
COPY cmd/user-service/go.mod cmd/user-service/go.sum ./cmd/user-service/

# 下载依赖（如果依赖文件没变，这步会被缓存）
RUN go mod download

# 再复制源代码
COPY . .

# 构建
RUN cd cmd/user-service && go build -o user-service .
```

### 7. 配置管理分离

#### 应用特定配置

每个应用有独立的配置目录：

```
configs/
├── user-service/
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
└── order-service/
    ├── config.yaml
    ├── config.dev.yaml
    └── config.prod.yaml
```

#### 环境变量注入

使用环境变量覆盖配置：

```yaml
# configs/user-service/config.yaml
server:
  http:
    addr: ${HTTP_ADDR:0.0.0.0:8001}
  grpc:
    addr: ${GRPC_ADDR:0.0.0.0:9001}
```

### 8. 版本管理策略

#### 独立版本号

每个应用可以有独立的版本号：

```go
// cmd/user-service/main.go
var (
    Name    = "user-service"
    Version = "1.2.3"  // 用户服务版本
)

// cmd/order-service/main.go
var (
    Name    = "order-service"
    Version = "2.1.0"  // 订单服务版本
)
```

#### 共享包版本

共享包使用语义化版本：

```go
// pkg/go.mod
module sre/pkg

go 1.21

// 版本号：v1.0.0, v1.1.0, v2.0.0 等
```

### 9. 代码复用策略

#### 共享工具包

将通用工具放在 `pkg/` 目录：

```
pkg/
├── errors/          # 错误定义
│   └── errors.go
├── logger/          # 日志封装
│   └── logger.go
├── middleware/      # 中间件
│   └── auth.go
└── utils/           # 工具函数
    ├── crypto.go
    └── validator.go
```

#### 接口抽象

使用接口解耦服务依赖：

```go
// pkg/interfaces/user.go
package interfaces

type UserRepository interface {
    GetUser(ctx context.Context, id int64) (*User, error)
}

// cmd/user-service/internal/data/user.go
package data

import "sre/pkg/interfaces"

type userRepo struct {
    // 实现
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*User, error) {
    // 实现
}
```

### 10. 测试策略

#### 独立测试

每个应用有独立的测试：

```bash
# 测试用户服务
go test ./cmd/user-service/...

# 测试订单服务
go test ./cmd/order-service/...

# 测试共享包
go test ./pkg/...
```

#### 集成测试

为服务间通信编写集成测试：

```go
// cmd/user-service/internal/integration/user_test.go
package integration

func TestUserService_GetUser(t *testing.T) {
    // 测试用户服务
}
```

## 注意事项

1. **代码共享 vs 独立**：根据项目规模决定是否共享代码
2. **版本管理**：多应用时注意 API 版本管理
3. **测试策略**：每个应用需要独立的测试
4. **监控和日志**：为每个应用配置独立的监控和日志
5. **依赖管理**：独立部署时注意依赖版本的一致性
6. **构建优化**：使用构建缓存和标签减少构建时间
7. **部署隔离**：确保应用间不会相互影响

