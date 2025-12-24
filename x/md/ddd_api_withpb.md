以下是在 API 服务中引入 Google Protocol Buffers 后，重新生成目录结构和相关代码示例的详细步骤。

1. 项目目录结构
首先，我们来规划一下项目的目录结构。引入 Protocol Buffers 后，会新增 proto 目录用于存放 .proto 文件，生成的 .pb.go 文件也会放在相关位置。以下是一个典型的项目目录结构：

your_api_project/
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── application
│   │   └── usecase
│   │       └── user_usecase.go
│   ├── domain
│   │   └── model
│   │       └── user.go
│   └── infrastructure
│       └── http
│           └── handler
│               └── user_handler.go
├── proto
│   └── user.proto
|   └── user.pb.go
|   └── user_grpc.pb.go
├── go.mod
├── go.sum
2. 详细代码示例
2.1 proto/user.proto
syntax = "proto3";

package proto;

// 用户消息定义
message User {
  string id = 1;
  string name = 2;
  string email = 3;
}

// 获取用户请求消息
message GetUserRequest {
  string id = 1;
}

// 获取用户响应消息
message GetUserResponse {
  User user = 1;
}

// 用户服务定义
service UserService {
  // 获取用户信息的 RPC 方法
  rpc GetUser (GetUserRequest) returns (GetUserResponse);
}
2.2 生成 Go 代码
在项目根目录下执行以下命令，使用 protoc 编译器和 Go 插件生成对应的 Go 代码：

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto
执行后，会在 proto 目录下生成 user.pb.go 和 user_grpc.pb.go 文件。

2.3 internal/domain/model/user.go
package model

import (
    "github.com/your_api_project/proto"
)

// 直接使用 proto 消息作为模型
type User = proto.User
2.4 internal/application/usecase/user_usecase.go
package usecase

import (
    "context"
    "github.com/your_api_project/internal/domain/model"
    "github.com/your_api_project/proto"
)

// UserUsecase 定义用户用例接口
type UserUsecase interface {
    GetUser(ctx context.Context, id string) (*model.User, error)
}

// userUsecase 实现用户用例接口
type userUsecase struct{}

// NewUserUsecase 创建用户用例实例
func NewUserUsecase() UserUsecase {
    return &userUsecase{}
}

// GetUser 实现获取用户信息的方法
func (u *userUsecase) GetUser(ctx context.Context, id string) (*model.User, error) {
    // 模拟获取用户信息
    user := &proto.User{
        Id:    id,
        Name:  "John Doe",
        Email: "johndoe@example.com",
    }
    return user, nil
}
2.5 internal/infrastructure/http/handler/user_handler.go
package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/your_api_project/internal/application/usecase"
    "github.com/your_api_project/proto"
)

// UserHandler 定义用户处理程序结构体
type UserHandler struct {
    userUsecase usecase.UserUsecase
}

// NewUserHandler 创建用户处理程序实例
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{
        userUsecase: userUsecase,
    }
}

// GetUser 处理获取用户信息的 HTTP 请求
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    user, err := h.userUsecase.GetUser(context.Background(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 将 proto 消息转换为 JSON 响应
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
2.6 cmd/api/main.go
package main

import (
    "log"
    "net/http"
    "github.com/your_api_project/internal/application/usecase"
    "github.com/your_api_project/internal/infrastructure/http/handler"
    "github.com/gorilla/mux"
)

func main() {
    // 创建用户用例实例
    userUsecase := usecase.NewUserUsecase()
    // 创建用户处理程序实例
    userHandler := handler.NewUserHandler(userUsecase)

    // 创建路由
    r := mux.NewRouter()
    // 注册获取用户信息的路由
    r.HandleFunc("/users", userHandler.GetUser).Methods("GET")

    log.Println("Starting API server on port 8080...")
    // 启动 HTTP 服务器
    log.Fatal(http.ListenAndServe(":8080", r))
}
3. 测试 API
在项目根目录下，使用以下命令启动 API 服务：

go run cmd/api/main.go
然后使用 curl 工具测试 API：

curl http://localhost:8080/users?id=123
你应该会看到类似以下的 JSON 响应：

{
    "id": "123",
    "name": "John Doe",
    "email": "johndoe@example.com"
}
通过以上步骤，你就完成了在 API 服务中引入 Google Protocol Buffers 并重新生成目录和代码的过程。