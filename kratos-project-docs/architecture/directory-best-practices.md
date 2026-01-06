# Kratos 项目目录结构最佳实践

## 概述

本文档基于 Kratos 框架和 Clean Architecture 原则，详细说明项目目录结构的最佳实践，帮助开发者正确组织代码，避免架构混乱。

## 目录结构总览

```
sre/
├── api/                    # API 定义层（Protobuf）
│   ├── user/
│   │   └── v1/
│   └── external/           # 第三方服务 API 定义
│       └── {service}/
│           └── v1/
├── cmd/                    # 应用入口目录
│   ├── sre/                # 主服务应用
│   ├── cron-worker/        # Cron 任务应用
│   ├── daemon-worker/      # Daemon 任务应用
│   └── sre-client/         # 命令行客户端
├── configs/                # 配置文件
├── internal/               # 内部代码（不对外暴露）
│   ├── app/                # 应用协调层（应用生命周期管理）
│   ├── biz/                # 业务逻辑层（核心业务代码）
│   ├── data/               # 数据访问层（数据库、外部服务）
│   ├── service/            # 服务层（gRPC/HTTP 处理）
│   ├── server/             # 服务器配置
│   ├── conf/               # 配置结构定义
│   ├── config/             # 配置加载器
│   ├── logger/             # 日志封装
│   ├── registry/           # 服务注册发现
│   └── pkd/                # 公共工具包
├── third_party/            # 第三方 Protobuf 定义
├── docs/                   # 文档目录
├── go.mod
└── README.md
```

## 核心目录详解

### 1. `api/` - API 定义层

**职责**：存放 Protobuf 接口定义文件

**结构**：
```
api/
├── user/                    # 业务服务 API
│   └── v1/
│       ├── user.proto
│       └── error_reason.proto
└── external/                # 第三方服务 API 定义
    ├── dingtalk/
    │   └── v1/
    └── wps/
        └── v1/
```

**最佳实践**：
- ✅ 按服务名和版本组织目录
- ✅ 业务服务 API 放在 `api/{service}/v1/`
- ✅ 第三方服务 API 定义放在 `api/external/{service}/v1/`
- ✅ 使用版本号管理 API 变更（v1, v2, ...）
- ❌ 不要在 `api/` 目录下放置实现代码

**示例**：
```protobuf
// api/user/v1/user.proto
syntax = "proto3";
package api.user.v1;

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserReply);
}
```

### 2. `cmd/` - 应用入口目录

**职责**：存放应用入口点，每个应用一个子目录

**结构**：
```
cmd/
├── sre/                    # 主服务应用
│   ├── main.go
│   ├── wire.go
│   └── wire_gen.go
├── cron-worker/            # Cron 任务应用
├── daemon-worker/          # Daemon 任务应用
└── sre-client/             # 命令行客户端
```

**最佳实践**：
- ✅ 每个应用有独立的目录和入口
- ✅ 使用 Wire 进行依赖注入
- ✅ 每个应用可以有独立的配置文件
- ✅ 应用之间可以共享 `internal/` 目录下的代码
- ❌ 不要在 `cmd/` 目录下放置业务逻辑

**示例**：
```go
// cmd/sre/main.go
func main() {
    flag.Parse()
    bootstrap, err := config.LoadBootstrapWithViper(flagconf)
    // ...
    app, cleanup, err := wireApp(...)
    defer cleanup()
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

### 3. `internal/app/` - 应用协调层

**职责**：管理应用级别的组件和生命周期

**适用场景**：
- 应用级别的协调组件（如 Worker Manager）
- 需要管理生命周期的服务（如 DingTalk Event Service）
- 跨业务模块的协调逻辑
- 应用启动/停止时的初始化逻辑

**结构**：
```
internal/app/
├── worker/
│   └── manager.go          # Worker 管理器
└── dingtalk/
    └── event.go            # 钉钉事件服务（生命周期管理）
```

**最佳实践**：
- ✅ 用于应用级别的协调和管理
- ✅ 可以依赖 `biz` 层和 `data` 层
- ✅ 管理组件的启动和停止
- ✅ 协调多个业务模块
- ❌ 不要包含具体的业务逻辑（业务逻辑应该在 `biz` 层）

**示例**：
```go
// internal/app/worker/manager.go
type Manager struct {
    daemonJobSet *DaemonJobSet
    cronManager  *cron.Manager
}

func (m *Manager) Start(ctx context.Context) error {
    // 启动所有 Worker
}
```

### 4. `internal/biz/` - 业务逻辑层

**职责**：核心业务逻辑，不依赖外部实现

**设计原则**：
- **不依赖外部**：不直接依赖数据库、Redis、HTTP 客户端等
- **定义接口**：定义数据访问接口，由 `data` 层实现
- **纯业务逻辑**：只包含业务规则、流程编排、验证等

**结构**：
```
internal/biz/
├── user.go                 # 用户业务逻辑
├── dingtalk_event.go       # 钉钉事件业务逻辑
├── cron/                   # Cron 任务接口和业务逻辑
│   ├── job.go
│   ├── manager.go
│   └── jobs/
└── daemon/                 # Daemon 任务接口和业务逻辑
    ├── job.go
    └── table_consumer.go
```

**最佳实践**：
- ✅ 定义 Repository 接口（如 `UserRepo`）
- ✅ 实现 UseCase（业务用例）
- ✅ 包含业务规则和验证逻辑
- ✅ 可以定义任务接口（如 `CronJob`、`DaemonJob`）
- ❌ 不要直接访问数据库或外部服务
- ❌ 不要包含应用级别的协调逻辑

**示例**：
```go
// internal/biz/user.go
type UserRepo interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id int64) (*User, error)
}

type UserUsecase struct {
    repo UserRepo
    log  *log.Helper
}

func (uc *UserUsecase) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // 业务逻辑：验证、规则检查、流程编排
    user := &User{...}
    return uc.repo.Save(ctx, user)
}
```

### 5. `internal/data/` - 数据访问层

**职责**：所有外部依赖的实现（数据库、Redis、Kafka、第三方服务等）

**设计原则**：
- **实现接口**：实现 `biz` 层定义的接口
- **管理外部依赖**：统一管理所有外部依赖
- **数据转换**：将外部数据格式转换为业务对象

**结构**：
```
internal/data/
├── data.go                 # 数据层初始化
├── user_repo.go            # 用户 Repository 实现
├── model/                  # 数据模型
│   └── user.go
├── external/               # 第三方服务客户端
│   ├── dingtalk/
│   │   ├── client.go
│   │   └── types.go
│   └── wps/
│       ├── client.go
│       └── types.go
├── kafka.go                # Kafka 客户端
├── redis.go                # Redis 客户端
└── sql/                    # SQL 脚本
```

**最佳实践**：
- ✅ 实现 `biz` 层定义的接口
- ✅ 管理所有外部依赖（数据库、Redis、Kafka、第三方服务）
- ✅ 第三方服务客户端放在 `external/{service}/`
- ✅ 第三方服务的类型定义放在 `external/{service}/types.go`
- ✅ gRPC 服务的 API 定义放在 `api/external/{service}/v1/`
- ✅ HTTP REST API 的类型定义放在 `internal/data/external/{service}/types.go`
- ❌ 不要包含业务逻辑
- ❌ 不要直接暴露给 `service` 层

**示例**：
```go
// internal/data/user_repo.go
type userRepo struct {
    data *Data
    log  *log.Helper
}

func (r *userRepo) Save(ctx context.Context, user *biz.User) error {
    // 数据库操作
    return r.data.db.Create(user).Error
}

// internal/data/external/dingtalk/client.go
type Client struct {
    httpClient *http.Client
}

func (c *Client) GetUser(ctx context.Context, userID string) (*types.User, error) {
    // HTTP 调用第三方服务
}
```

### 6. `internal/service/` - 服务层

**职责**：实现 gRPC/HTTP 接口，处理协议转换

**结构**：
```
internal/service/
├── service.go
└── user.go                 # 用户服务实现
```

**最佳实践**：
- ✅ 实现 `api/` 中定义的接口
- ✅ 参数校验和协议转换
- ✅ 调用 `biz` 层的 UseCase
- ✅ 错误处理和响应转换
- ❌ 不要包含业务逻辑
- ❌ 不要直接访问 `data` 层

**示例**：
```go
// internal/service/user.go
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
    // 参数校验
    if req.Id == 0 {
        return nil, errors.BadRequest("INVALID_ID", "id is required")
    }
    
    // 调用业务层
    user, err := s.uc.GetUser(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    
    // 协议转换
    return &pb.GetUserReply{
        User: &pb.User{
            Id:   user.ID,
            Name: user.Name,
        },
    }, nil
}
```

### 7. `internal/server/` - 服务器配置

**职责**：HTTP/gRPC 服务器配置

**结构**：
```
internal/server/
├── server.go
├── http.go
└── grpc.go
```

**最佳实践**：
- ✅ 配置 HTTP/gRPC 服务器
- ✅ 注册中间件
- ✅ 注册路由和服务
- ❌ 不要包含业务逻辑

### 8. `internal/conf/` - 配置结构定义

**职责**：配置结构定义（从 Protobuf 生成）

**结构**：
```
internal/conf/
├── conf.proto
└── conf.pb.go
```

**最佳实践**：
- ✅ 使用 Protobuf 定义配置结构
- ✅ 通过 `make config` 生成 Go 代码
- ✅ 支持多环境配置

### 9. `internal/config/` - 配置加载器

**职责**：配置加载和解析

**结构**：
```
internal/config/
├── loader.go               # 配置加载器
├── viper.go                # Viper 实现
├── kratos.go               # Kratos 配置
└── center.go               # 配置中心支持
```

**最佳实践**：
- ✅ 支持多种配置源（文件、配置中心）
- ✅ 支持配置热更新
- ✅ 统一的配置加载接口

## 关键区别：`internal/app` vs `internal/biz`

### `internal/app/` - 应用协调层

**用途**：应用级别的协调和管理

**特点**：
- ✅ 可以依赖 `biz` 层和 `data` 层
- ✅ 管理应用生命周期（启动、停止）
- ✅ 协调多个业务模块
- ✅ 应用级别的组件管理

**示例场景**：
- Worker Manager：管理所有 Worker 的启动和停止
- DingTalk Event Service：管理钉钉事件服务的生命周期
- 应用级别的初始化逻辑

### `internal/biz/` - 业务逻辑层

**用途**：核心业务逻辑

**特点**：
- ✅ 不依赖外部实现（只依赖接口）
- ✅ 定义数据访问接口
- ✅ 纯业务逻辑和规则
- ✅ 可以被多个应用复用

**示例场景**：
- UserUsecase：用户相关的业务逻辑
- DingTalkEventUsecase：钉钉事件处理的业务逻辑
- CronJob/DaemonJob 接口定义：任务接口定义

## 目录选择决策树

### 我应该把代码放在哪里？

```
开始
  │
  ├─ 是 Protobuf 定义？
  │   └─→ api/
  │
  ├─ 是应用入口？
  │   └─→ cmd/{app}/
  │
  ├─ 是业务逻辑？
  │   ├─ 依赖外部实现（数据库、HTTP 等）？
  │   │   └─→ internal/data/  （实现接口）
  │   │
  │   └─ 纯业务逻辑（不依赖外部）？
  │       └─→ internal/biz/  （定义接口和业务逻辑）
  │
  ├─ 是应用级别的协调？
  │   └─→ internal/app/  （管理生命周期、协调组件）
  │
  ├─ 是 gRPC/HTTP 接口实现？
  │   └─→ internal/service/
  │
  ├─ 是服务器配置？
  │   └─→ internal/server/
  │
  ├─ 是配置相关？
  │   ├─ 配置结构定义？
  │   │   └─→ internal/conf/
  │   │
  │   └─ 配置加载器？
  │       └─→ internal/config/
  │
  └─ 是工具函数？
      └─→ internal/pkd/
```

## 常见场景示例

### 场景 1：新增业务功能

**需求**：实现订单创建功能

**目录结构**：
```
api/order/v1/
  └── order.proto           # 定义订单 API

internal/biz/
  └── order.go              # 订单业务逻辑（定义 OrderRepo 接口）

internal/data/
  └── order_repo.go         # 实现 OrderRepo 接口

internal/service/
  └── order.go              # 实现 gRPC/HTTP 接口
```

### 场景 2：集成第三方服务

**需求**：集成支付服务（HTTP REST API）

**目录结构**：
```
api/external/payment/v1/
  └── payment.proto         # 如果支付服务提供 gRPC

internal/data/external/payment/
  ├── client.go             # HTTP 客户端实现
  └── types.go              # 请求/响应类型定义

internal/biz/
  └── payment.go            # 支付业务逻辑（定义 PaymentClient 接口）

internal/data/
  └── payment_client.go     # 实现 PaymentClient 接口
```

### 场景 3：新增后台任务

**需求**：实现定时同步用户数据

**目录结构**：
```
internal/biz/cron/
  ├── job.go                # CronJob 接口定义
  └── jobs/
      └── sync_user.go      # 同步用户任务（实现 CronJob）

internal/app/worker/
  └── manager.go            # Worker Manager（管理任务启动）

cmd/cron-worker/
  └── main.go               # Cron Worker 应用入口
```

### 场景 4：应用级别的服务

**需求**：实现消息队列消费者服务

**目录结构**：
```
internal/biz/
  └── message_consumer.go   # 消息消费的业务逻辑

internal/app/message/
  └── consumer_service.go   # 消息消费者服务（管理生命周期）

cmd/sre/
  └── main.go               # 在主应用中注册服务
```

## 最佳实践总结

### ✅ 应该做的

1. **严格遵循分层架构**：Service → Biz → Data
2. **接口定义在 Biz 层**：Biz 层定义接口，Data 层实现
3. **业务逻辑在 Biz 层**：纯业务逻辑放在 `internal/biz/`
4. **应用协调在 App 层**：应用级别的协调放在 `internal/app/`
5. **外部依赖在 Data 层**：所有外部依赖统一在 `internal/data/` 管理
6. **多应用共享代码**：多个应用共享 `internal/` 目录下的代码
7. **版本管理 API**：使用版本号管理 API 变更

### ❌ 不应该做的

1. **不要跨层调用**：Service 层不要直接调用 Data 层
2. **不要在 Biz 层依赖外部**：Biz 层只依赖接口，不依赖具体实现
3. **不要在 App 层写业务逻辑**：业务逻辑应该在 Biz 层
4. **不要在 Data 层写业务逻辑**：Data 层只负责数据访问
5. **不要循环依赖**：严格遵循依赖方向
6. **不要在 cmd/ 写业务逻辑**：cmd/ 只放应用入口

## 参考文档

- [分层架构设计](./layered-architecture.md)
- [多应用支持](./multi-app.md)
- [依赖注入](./dependency-injection.md)
- [第三方服务集成](./third-party-api-definitions.md)
- [项目结构](../project/structure.md)

