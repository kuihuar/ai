# 项目结构

## Kratos 项目标准结构

```
sre/
├── api/                    # API 定义层
│   └── helloworld/
│       └── v1/            # API 版本
│           ├── greeter.proto
│           └── error_reason.proto
├── cmd/                    # 应用入口
│   └── sre/
│       ├── main.go        # 主函数
│       ├── wire.go        # Wire 定义
│       └── wire_gen.go    # Wire 生成代码
├── configs/                # 配置文件
│   └── config.yaml
├── internal/               # 内部代码
│   ├── biz/               # 业务逻辑层
│   │   ├── greeter.go
│   │   └── biz.go
│   ├── data/              # 数据访问层
│   │   ├── greeter.go
│   │   └── data.go
│   ├── service/           # 服务层
│   │   ├── greeter.go
│   │   └── service.go
│   ├── server/            # 服务器配置
│   │   ├── http.go
│   │   ├── grpc.go
│   │   └── server.go
│   └── conf/              # 配置定义
│       └── conf.proto
├── docs/                   # 文档目录
│   ├── architecture/
│   ├── code-standards/
│   ├── development/
│   ├── operations/
│   └── project/
├── third_party/            # 第三方 Protobuf 定义
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 目录说明

### api/
存放 Protobuf 接口定义文件，按服务名和版本组织。

### cmd/
应用入口目录，每个应用一个子目录。

**多应用支持**：Kratos 支持在同一项目中管理多个应用，每个应用在 `cmd/` 目录下有独立的子目录。每个应用可以：
- 共享 `internal/` 目录下的业务逻辑
- 使用独立的配置文件和启动入口
- 通过 gRPC 或 HTTP 进行服务间通信

详细说明请参考 [多应用支持文档](../architecture/multi-app.md)。

### configs/
配置文件目录，支持多环境配置。

### internal/
内部代码目录，不对外暴露。

- **biz/**：业务逻辑层，核心业务代码
- **data/**：数据访问层，数据库和外部服务调用
- **service/**：服务层，实现 API 接口
- **server/**：服务器配置，HTTP/gRPC 服务器
- **conf/**：配置结构定义

### docs/
文档目录，存放项目文档和最佳实践。

### third_party/
第三方 Protobuf 定义，如 Google API、OpenAPI 等。

## 文件命名规范

- **Go 文件**：小写字母，下划线分隔（如 `greeter.go`）
- **Proto 文件**：小写字母，下划线分隔（如 `greeter.proto`）
- **配置文件**：小写字母，下划线或连字符（如 `config.yaml`）

## 最佳实践

1. **保持结构清晰**：按职责组织代码
2. **避免循环依赖**：遵循依赖方向
3. **文档同步**：代码变更时更新文档
4. **版本管理**：API 变更时使用版本号

