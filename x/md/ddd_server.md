在 DDD 架构中，cmd 目录通常用来存放应用程序的入口文件，为了遵循单一职责原则，让代码结构更清晰，我们可以将服务器的启动逻辑从 cmd 中的入口文件分离出来。以下是具体的操作步骤和示例代码：

1. 目录结构调整
假设原有的目录结构如下：

your_project/
├── cmd/
│   ├── grpc_server/
│   │   └── main.go
│   └── http_server/
│       └── main.go
├── internal/
│   ├── application/
│   ├── domain/
│   └── infrastructure/
│       ├── persistence/
│       ├── grpc/
│       └── http/
└── go.mod
我们可以在 internal 目录下创建一个新的 server 目录，用于存放服务器的启动逻辑，调整后的目录结构如下：

your_project/
├── cmd/
│   ├── grpc_server/
│   │   └── main.go
│   └── http_server/
│       └── main.go
├── internal/
│   ├── application/
│   ├── domain/
│   ├── infrastructure/
│   │   ├── persistence/
│   │   ├── grpc/
│   │   └── http/
│   └── server/
│       ├── grpc/
│       │   └── server.go
│       └── http/
│           └── server.go
└── go.mod
2. 分离 gRPC 服务器逻辑
internal/server/grpc/server.go
package grpc

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "your_project/internal/application"
    "your_project/internal/infrastructure/grpc"
)

// StartGRPCServer 启动 gRPC 服务器
func StartGRPCServer(ctx context.Context, app *application.YourServiceApp, addr string) error {
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }

    s := grpc.NewServer()
    grpc.RegisterYourServiceServer(s, grpc.NewYourServiceServer(app))

    log.Printf("Starting gRPC server on %s", addr)
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Printf("Failed to serve gRPC: %v", err)
        }
    }()

    // 监听上下文取消信号，优雅关闭服务器
    <-ctx.Done()
    log.Println("Shutting down gRPC server...")
    s.GracefulStop()
    return nil
}
cmd/grpc_server/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "your_project/internal/application"
    "your_project/internal/server/grpc"
)

func main() {
    app := application.NewYourServiceApp()
    addr := ":50051"

    ctx, cancel := context.WithCancel(context.Background())

    // 监听系统信号，用于优雅关闭服务器
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    if err := grpc.StartGRPCServer(ctx, app, addr); err != nil {
        log.Fatalf("Failed to start gRPC server: %v", err)
    }
}
3. 分离 HTTP 服务器逻辑
internal/server/http/server.go
package http

import (
    "context"
    "log"
    "net/http"
    "time"

    "your_project/internal/application"
    "your_project/internal/infrastructure/http"
)

// StartHTTPServer 启动 HTTP 服务器
func StartHTTPServer(ctx context.Context, app *application.YourServiceApp, addr string) error {
    httpServer := http.NewYourServiceHTTPServer(app)
    mux := http.NewServeMux()
    mux.HandleFunc("/sayhello", httpServer.SayHelloHandler)

    srv := &http.Server{
        Addr:    addr,
        Handler: mux,
    }

    log.Printf("Starting HTTP server on %s", addr)
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("Failed to serve HTTP: %v", err)
        }
    }()

    // 监听上下文取消信号，优雅关闭服务器
    <-ctx.Done()
    log.Println("Shutting down HTTP server...")
    ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return srv.Shutdown(ctxShutDown)
}
cmd/http_server/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "your_project/internal/application"
    "your_project/internal/server/http"
)

func main() {
    app := application.NewYourServiceApp()
    addr := ":8080"

    ctx, cancel := context.WithCancel(context.Background())

    // 监听系统信号，用于优雅关闭服务器
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    if err := http.StartHTTPServer(ctx, app, addr); err != nil {
        log.Fatalf("Failed to start HTTP server: %v", err)
    }
}
总结
通过上述操作，我们将服务器的启动逻辑从 cmd 目录下的入口文件分离到了 internal/server 目录中，使得 cmd 目录下的入口文件只负责简单的初始化和启动操作，而服务器的具体启动逻辑则封装在 internal/server 目录中，提高了代码的可维护性和可测试性。同时，我们还实现了服务器的优雅关闭功能，确保在接收到系统信号时能够正常关闭服务器。