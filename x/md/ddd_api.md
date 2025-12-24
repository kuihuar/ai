当开发提供 API 服务的 Go 项目时，采用合理的目录结构有助于提高代码的可维护性、可扩展性和可读性。以下是一种常见且规范的目录结构示例，结合了分层架构和模块化设计的思想：

.
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── domain
│   │   ├── model
│   │   │   ├── user.go
│   │   │   └── product.go
│   │   └── service
│   │       ├── user_service.go
│   │       └── product_service.go
│   ├── application
│   │   ├── dto
│   │   │   ├── user_dto.go
│   │   │   └── product_dto.go
│   │   └── usecase
│   │       ├── user_usecase.go
│   │       └── product_usecase.go
│   ├── infrastructure
│   │   ├── database
│   │   │   ├── db_connection.go
│   │   │   ├── user_repository.go
│   │   │   └── product_repository.go
│   │   └── http
│   │       ├── middleware
│   │       │   ├── auth_middleware.go
│   │       │   └── logging_middleware.go
│   │       └── handler
│   │           ├── user_handler.go
│   │           └── product_handler.go
├── pkg
│   ├── utils
│   │   ├── validation.go
│   │   └── response.go
├── config
│   └── config.yaml
├── go.mod
├── go.sum
各目录和文件的详细说明
1. cmd 目录
用途：存放项目的可执行文件入口。每个可执行程序通常对应一个子目录。
cmd/api/main.go：API 服务的入口文件，负责初始化依赖、配置路由和启动服务器。示例代码如下：
package main

import (
    "log"
    "net/http"
    "github.com/yourproject/internal/infrastructure/http/handler"
    "github.com/yourproject/internal/infrastructure/http/middleware"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.Use(middleware.LoggingMiddleware)
    r.Use(middleware.AuthMiddleware)

    // 注册处理程序
    handler.RegisterUserHandlers(r)
    handler.RegisterProductHandlers(r)

    log.Println("Starting API server on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
2. internal 目录
用途：存放项目内部使用的代码，这些代码不对外暴露，只能在项目内部被引用。按照不同的层次进行划分，包括领域层、应用层和基础设施层。
2.1 domain 目录
用途：包含领域层的代码，如实体和领域服务。
domain/model：存放领域模型，即业务实体的定义。
user.go：定义 User 实体及其相关方法。
product.go：定义 Product 实体及其相关方法。
domain/service：实现领域服务，处理复杂的业务逻辑。
user_service.go：提供用户相关的领域服务。
product_service.go：提供产品相关的领域服务。
2.2 application 目录
用途：存放应用层的代码，主要是用例和数据传输对象（DTO）。
application/dto：定义数据传输对象，用于在不同层之间传递数据。
user_dto.go：定义用户相关的 DTO。
product_dto.go：定义产品相关的 DTO。
application/usecase：实现用例，协调领域服务和基础设施层，处理业务流程。
user_usecase.go：处理用户相关的业务用例。
product_usecase.go：处理产品相关的业务用例。
2.3 infrastructure 目录
用途：包含基础设施层的代码，如数据库访问、HTTP 处理等。
infrastructure/database：负责数据库相关操作。
db_connection.go：建立数据库连接。
user_repository.go：实现用户数据的持久化操作。
product_repository.go：实现产品数据的持久化操作。
infrastructure/http：处理 HTTP 请求和响应。
middleware：存放中间件，如身份验证、日志记录等。
auth_middleware.go：实现身份验证中间件。
logging_middleware.go：实现日志记录中间件。
handler：定义 HTTP 处理程序，负责接收请求、调用用例和返回响应。
user_handler.go：处理用户相关的 HTTP 请求。
product_handler.go：处理产品相关的 HTTP 请求。
3. pkg 目录
用途：存放可复用的工具函数和库，这些代码可以被项目内的其他模块或外部项目使用。
pkg/utils：包含一些通用的工具函数。
validation.go：实现数据验证功能。
response.go：封装 HTTP 响应的通用方法。
4. config 目录
用途：存放项目的配置文件，如数据库连接信息、服务器端口等。
config.yaml：以 YAML 格式存储配置信息。示例内容如下：
database:
  host: localhost
  port: 5432
  user: your_user
  password: your_password
  dbname: your_dbname
server:
  port: 8080
5. go.mod 和 go.sum 文件
go.mod：用于管理项目的依赖关系，记录项目所依赖的外部包及其版本信息。
go.sum：用于确保依赖包的版本和哈希值的一致性，防止依赖包被篡改。
通过这种目录结构，各个层次的代码职责明确，不同层之间的依赖关系清晰，便于团队协作和项目的长期维护。