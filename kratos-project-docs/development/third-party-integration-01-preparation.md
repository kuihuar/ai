# 第三方服务集成指南 - 第一步：准备工作

## 概述

在集成第三方服务之前，需要完成以下准备工作：
1. 扩展配置定义（如果需要）
2. 创建目录结构
3. 了解服务类型和集成方式

## 步骤 1: 扩展配置定义

### 1.1 更新 conf.proto

如果需要在配置文件中管理第三方服务的连接信息，需要在 `internal/conf/conf.proto` 中添加配置定义。

#### 添加 gRPC 客户端配置

```protobuf
// internal/conf/conf.proto
message Data {
  // ... 现有配置 ...
  
  // gRPC 客户端配置
  message GRPC {
    map<string, string> clients = 1;  // name -> endpoint
  }
  
  // HTTP 客户端配置
  message HTTP {
    map<string, ClientConfig> clients = 1;  // name -> config
  }
  
  message ClientConfig {
    string endpoint = 1;                    // 服务地址
    google.protobuf.Duration timeout = 2;   // 超时时间
    map<string, string> headers = 3;        // 默认请求头
  }
  
  // ... 现有字段 ...
  GRPC grpc = 5;
  HTTP http = 6;
}
```

#### 重新生成配置代码

```bash
make config
# 或
protoc --proto_path=. \
       --proto_path=./third_party \
       --go_out=paths=source_relative:. \
       internal/conf/conf.proto
```

### 1.2 更新配置文件

在 `configs/config.yaml` 中添加第三方服务配置：

```yaml
data:
  # ... 现有配置 ...
  
  # gRPC 客户端配置
  grpc:
    clients:
      user-service: 127.0.0.1:9001
      order-service: 127.0.0.1:9002
      payment-service: payment.example.com:443
  
  # HTTP 客户端配置
  http:
    clients:
      payment-api:
        endpoint: https://api.payment.com
        timeout: 5s
        headers:
          X-API-Key: your-api-key
      notification-api:
        endpoint: https://api.notification.com
        timeout: 3s
```

## 步骤 2: 创建目录结构

根据第三方服务的类型，创建相应的目录结构。

### 2.1 gRPC 服务目录结构

对于 gRPC 服务，在 `api/external/` 下创建服务目录：

```
api/
└── external/
    └── {service-name}/        # 服务名称，如 user-service
        └── v1/
            ├── {service}.proto
            ├── {service}.pb.go
            └── {service}_grpc.pb.go
```

**示例：**
```
api/
└── external/
    └── user-service/
        └── v1/
            ├── user.proto
            ├── user.pb.go
            └── user_grpc.pb.go
```

### 2.2 HTTP REST API 目录结构

对于 HTTP REST API，在 `internal/data/external/` 下创建服务目录：

```
internal/data/
└── external/
    └── {service-name}/        # 服务名称，如 payment
        ├── types.go           # 请求/响应类型定义
        ├── client.go          # HTTP 客户端实现
        └── {service}.go       # 业务封装（可选）
```

**示例：**
```
internal/data/
└── external/
    └── payment/
        ├── types.go
        ├── client.go
        └── payment.go
```

## 步骤 3: 确定服务类型

在开始集成前，需要明确第三方服务的类型：

| 服务类型 | 协议 | 定义位置 | 客户端位置 |
|---------|------|---------|-----------|
| 内部 gRPC 服务 | gRPC | `api/external/{service}/v1/` | `internal/data/clients/grpc.go` |
| 外部 gRPC 服务 | gRPC | `api/external/{service}/v1/` 或 `third_party/` | `internal/data/clients/grpc.go` |
| HTTP REST API | HTTP | `internal/data/external/{service}/types.go` | `internal/data/external/{service}/client.go` |
| GraphQL API | GraphQL | `internal/data/external/{service}/types.go` | `internal/data/external/{service}/client.go` |

## 步骤 4: 安装依赖

根据服务类型安装相应的依赖包。

### gRPC 服务依赖

```bash
# gRPC 核心库
go get google.golang.org/grpc
go get google.golang.org/protobuf

# Kratos gRPC 客户端
go get github.com/go-kratos/kratos/v2/transport/grpc
```

### HTTP 客户端依赖

```bash
# HTTP 客户端（推荐使用标准库或 resty）
go get github.com/go-resty/resty/v2

# 或使用 Kratos HTTP 客户端
go get github.com/go-kratos/kratos/v2/transport/http
```

## 下一步

完成准备工作后，根据服务类型选择相应的集成步骤：

- **gRPC 服务**：参考 [第二步：gRPC 服务集成](./third-party-integration-02-grpc.md)
- **HTTP REST API**：参考 [第三步：HTTP REST API 集成](./third-party-integration-03-http.md)

## 注意事项

1. **配置管理**：如果服务地址是动态的或通过服务发现获取，可能不需要在配置文件中硬编码
2. **版本管理**：对于 gRPC 服务，建议使用版本号目录（如 `v1/`, `v2/`）来管理 API 版本
3. **目录命名**：使用小写字母和连字符，与服务名称保持一致
4. **依赖注入**：确保在 `internal/data/data.go` 的 `ProviderSet` 中包含新的客户端提供者

