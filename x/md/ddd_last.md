领域驱动设计（DDD）分层架构通常包含用户界面层（Presentation）、应用层（Application）、领域层（Domain）和基础设施层（Infrastructure）。下面为你展示如何结合之前的 gRPC 和 HTTP 服务示例，生成符合 DDD 分层架构的目录结构，并简单说明各层的代码组织。

1. 目录结构
your_project/
├── api/
│   └── pb/
│       ├── your_service.proto
│       ├── your_service.pb.go
│       ├── your_service_grpc.pb.go
│       └── your_service.pb.gw.go
├── cmd/
│   ├── grpc_server/
│   │   └── main.go
│   └── http_server/
│       └── main.go
├── internal/
│   ├── application/
│   │   └── your_service_app.go
│   ├── domain/
│   │   ├── model/
│   │   │   └── your_service_model.go
│   │   └── service/
│   │       └── your_service_domain.go
│   └── infrastructure/
│       ├── persistence/
│       │   └── your_service_persistence.go
│       └── grpc/
│           └── your_service_grpc_server.go
└── go.mod
2. 各层说明及代码示例
api 目录
此目录存放 .proto 文件以及生成的代码。

api/
└── pb/
    ├── your_service.proto
    ├── your_service.pb.go
    ├── your_service_grpc.pb.go
    └── your_service.pb.gw.go
这些文件是通过 protoc 工具根据 .proto 文件生成的，用于定义 gRPC 服务和消息结构。

cmd 目录
该目录包含可执行程序的入口文件。

cmd/
├── grpc_server/
│   └── main.go
└── http_server/
    └── main.go
grpc_server/main.go：启动 gRPC 服务的入口文件。
package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    "your_project/internal/infrastructure/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    grpc.RegisterYourServiceServer(grpcServer, grpc.NewYourServiceServer())

    log.Println("Starting gRPC server on port 50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC: %v", err)
    }
}
http_server/main.go：启动 HTTP 服务的入口文件。
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "your_project/api/pb"
)

func main() {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    err := pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
    if err != nil {
        log.Fatalf("Failed to register HTTP handler: %v", err)
    }

    log.Println("Starting HTTP server on port 8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("Failed to serve HTTP: %v", err)
    }
}
internal 目录
此目录包含项目的核心业务逻辑，按照 DDD 分层架构进行组织。

application 目录
负责协调领域层和基础设施层，处理业务流程。

internal/
└── application/
    └── your_service_app.go
package application

import (
    "your_project/internal/domain/service"
)

type YourServiceApp struct {
    yourServiceDomain service.YourServiceDomain
}

func NewYourServiceApp() *YourServiceApp {
    return &YourServiceApp{
        yourServiceDomain: service.NewYourServiceDomain(),
    }
}

func (a *YourServiceApp) SayHello(message string) string {
    return a.yourServiceDomain.SayHello(message)
}
domain 目录
包含领域模型和领域服务，实现核心业务逻辑。

internal/
└── domain/
    ├── model/
    │   └── your_service_model.go
    └── service/
        └── your_service_domain.go
model/your_service_model.go：定义领域模型。
package model

type YourServiceModel struct {
    Message string
}
service/your_service_domain.go：实现领域服务。
package service

import (
    "your_project/internal/domain/model"
)

type YourServiceDomain struct{}

func NewYourServiceDomain() *YourServiceDomain {
    return &YourServiceDomain{}
}

func (s *YourServiceDomain) SayHello(message string) string {
    model := model.YourServiceModel{Message: message}
    return "Hello, " + model.Message
}
infrastructure 目录
提供底层的基础设施服务，如数据库访问、网络通信等。

internal/
└── infrastructure/
    ├── persistence/
    │   └── your_service_persistence.go
    └── grpc/
        └── your_service_grpc_server.go
persistence/your_service_persistence.go：处理数据持久化。
package persistence

import (
    "your_project/internal/domain/model"
)

type YourServicePersistence struct{}

func NewYourServicePersistence() *YourServicePersistence {
    return &YourServicePersistence{}
}

func (p *YourServicePersistence) Save(model model.YourServiceModel) error {
    // 实现数据保存逻辑
    return nil
}
grpc/your_service_grpc_server.go：实现 gRPC 服务。
package grpc

import (
    "context"

    "google.golang.org/grpc"
    "your_project/api/pb"
    "your_project/internal/application"
)

type YourServiceServer struct {
    pb.UnimplementedYourServiceServer
    app *application.YourServiceApp
}

func NewYourServiceServer() *YourServiceServer {
    return &YourServiceServer{
        app: application.NewYourServiceApp(),
    }
}

func (s *YourServiceServer) SayHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {
    message := s.app.SayHello(in.Message)
    return &pb.Response{Message: message}, nil
}
3. 总结
通过上述目录结构和代码示例，你可以将 gRPC 和 HTTP 服务整合到 DDD 分层架构中，实现代码的清晰组织和可维护性。各层之间通过依赖注入的方式进行解耦，便于后续的扩展和维护。