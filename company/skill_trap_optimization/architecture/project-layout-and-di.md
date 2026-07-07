# 目录结构与分层架构 + 依赖注入

## 一、项目目录结构

### 1.1 标准布局详解

```
project/
├── cmd/                    # 可执行程序入口，每个子目录一个 main.go
│   └── server/
│       └── main.go
├── internal/               # Go 编译器强制隔离的私有包
│   ├── handler/            # HTTP/gRPC 接口层：参数校验、序列化、响应
│   ├── service/            # 业务逻辑层：核心逻辑、事务编排
│   ├── repository/         # 数据访问层：DB/Cache/外部 API 调用
│   ├── model/              # 数据模型/领域对象（贫血模型，只有数据+tag）
│   ├── middleware/         # 中间件：Auth/Logger/RateLimit 等
│   └── router/             # 路由注册
├── pkg/                    # 可被外部项目导入的公共库
│   └── errcode/            # 统一错误码
├── api/                    # API 定义文件（proto、openapi/swagger）
│   └── proto/
│       └── user/v1/
│           └── user.proto
├── configs/                # 配置文件模板（yaml/toml/json）
├── scripts/                # 构建、部署、代码生成脚本
├── migrations/             # 数据库迁移 SQL 文件
├── deployments/            # Dockerfile、docker-compose、k8s 清单
├── go.mod
├── Makefile
└── README.md
```

### 1.2 为什么用 internal/ 而不是全都放 pkg/？

`internal/` 是 Go 编译器层面的强制保护——**放在 internal 下的包，外部项目 import 时会编译报错**。这比约定或 code review 更可靠。

```go
// 外部项目 import "github.com/your/project/internal/handler"
// 编译错误：use of internal package not allowed
```

**选择原则**：
- `internal/`：业务逻辑、当前项目独有代码，不希望被外部依赖
- `pkg/`：可复用的工具库、通用错误码、与业务无关的基础组件

### 1.3 是否一定要用这个布局？

**不用**。Go 官方没有强制标准布局。以下是不同阶段的建议：

| 阶段 | 推荐布局 | 理由 |
|------|----------|------|
| 原型/雏形 | `main.go` + 几个平级文件 | 不要过早分层 |
| 中等规模 | `cmd/` + `internal/` + 分层 | 结构清晰，编译隔离 |
| 大型单体/微服务 | 标准布局 + 按领域分包 | 多人协作，分而治之 |
| 单体仓库 monorepo | `pkg/` 共享库 + 多 `cmd/` | 代码复用，统一发版 |

### 1.4 常见变体：按领域分包

当业务复杂时，可以按领域（domain）拆分，避免单个包文件过多：

```
internal/
├── user/           # 用户领域
│   ├── handler.go
│   ├── service.go
│   └── repository.go
├── order/          # 订单领域
│   ├── handler.go
│   ├── service.go
│   └── repository.go
└── middleware/     # 跨领域的公共中间件
```

**分层 vs 分领域？** 团队小、模块少用分层；模块多、团队多按领域分，参考 DDD（领域驱动设计）思想。

---

## 二、分层架构

### 2.1 经典三层

```
┌─────────────────────────────┐
│  Handler (接口层)            │  ← 解析请求、参数校验、调用 Service、序列化响应
├─────────────────────────────┤
│  Service (业务层)            │  ← 核心逻辑、编排多个 Repository、事务边界
├─────────────────────────────┤
│  Repository (数据层)         │  ← DB SQL、缓存读写、外部 API 调用
├─────────────────────────────┤
│  Model (数据模型)            │  ← 结构体定义、DB tag，不包含方法
└─────────────────────────────┘
```

**依赖规则**：上层依赖下层，下层**定义接口**给上层用，不感知上层。

### 2.2 为什么 Service 要依赖 Repository 接口而不是具体实现？

这是**依赖反转原则（DIP）**——高层模块不依赖低层模块，两者都依赖抽象。

```go
// ===== 不推荐的写法（Service 依赖具体实现）=====
type UserService struct {
    repo *mysqlUserRepo  // 直接依赖 MySQL 实现
}
// 问题：换成 PostgreSQL 要改 Service；测试必须连 MySQL

// ===== 推荐的写法（依赖接口）=====
type UserRepo interface {
    GetByID(ctx context.Context, id int64) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
}

type UserService struct {
    repo UserRepo  // 依赖接口，不关心实现
}
// 好处：测试时传 mock；换数据库只改 NewXxxRepo 的实现，不改 Service
```

### 2.3 更多层的场景

当项目进一步复杂化，可以考虑加层：

| 额外层 | 作用 | 何时引入 |
|--------|------|----------|
| DTO (Data Transfer Object) | Handler 与 Service 之间的数据结构 | 接口返回结构与 DB 模型不同时 |
| Facade | 组合多个 Service | 复杂业务流程跨多个领域 |
| Adapter | 对接第三方服务 | 外部 API 多，需要统一封装 |

但原则是：**不要过早加层**。三层先跑起来，复杂度上来了再拆。

### 2.4 Clean Architecture（整洁架构）

如果团队推崇 Clean Architecture，可以进一步抽象：

```
Handler → UseCase → Repository → Entity
   ↑ 外层      ↑ 内层     ↑ 外层
```

- UseCase 是纯粹的业务用例，不依赖任何框架和数据库
- 所有依赖指向内层（Domain）
- 优点是极高的可测试性；代价是更多的接口和文件

**选型建议**：绝大多数 Go 项目用经典三层就够了，Clean Architecture 适合对测试覆盖率要求极高或业务规则极复杂的场景。

---

## 三、依赖注入（DI）

### 3.1 方案对比总览

| 维度 | 手动注入 | Wire (Google) | Uber Fx |
|------|----------|---------------|---------|
| 原理 | main.go 中手动 `NewXxx` | 代码生成，编译时检查 | 运行时反射+容器 |
| 复杂度 | 低，代码直接 | 中，需学习 Provider/Injector 概念 | 高，概念多（Provide/Invoke/Lifecycle）|
| 性能 | 无额外开销 | 无运行时开销 | 有运行时开销（反射） |
| 报错时机 | 编译时（类型不匹配） | 编译时（wire 生成报错） | 运行时（启动时才报错） |
| 依赖可视化 | 无 | `wire` 命令展示依赖图 | 无内建 |
| Go 惯用程度 | 最符合 Go 哲学 | 较符合 | 较重的框架风格 |
| 社区/维护 | - | Google 官方，中等活跃 | Uber 官方，活跃 |

### 3.2 手动构造函数注入

**适用场景**：依赖 < 20 个，小型项目，追求简单直接的团队。

```go
func main() {
    // 1. 基础设施
    db, _ := sql.Open("postgres", dsn)
    redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    // 2. Repository 层
    userRepo := repository.NewUserRepo(db)
    orderRepo := repository.NewOrderRepo(db)
    cacheRepo := repository.NewCacheRepo(redisClient)

    // 3. Service 层（可依赖多个 repo，也可依赖其他 service）
    userSvc := service.NewUserService(userRepo, cacheRepo)
    orderSvc := service.NewOrderService(orderRepo, userRepo)

    // 4. Handler 层
    userHandler := handler.NewUserHandler(userSvc)
    orderHandler := handler.NewOrderHandler(orderSvc)

    // 5. 路由注册
    r := gin.Default()
    userHandler.Register(r)
    orderHandler.Register(r)
    r.Run(":8080")
}
```

**优点**：
- 完全显式，new 一个看一个，不出意外
- 零学习成本，不需要理解 DI 框架
- 不引入额外的依赖包

**缺点**：
- 依赖超过 20 个后 main.go 会很长
- 依赖调整后需手动维护初始化顺序
- 单例管理、生命周期需要自己控制

### 3.3 Wire（推荐中型以上项目）

**适用场景**：依赖 > 20 个，需要编译时依赖检查，不希望引入运行时反射。

核心概念：
- **Provider**：一个返回某类型实例的构造函数，如 `func NewUserRepo(db *sql.DB) *UserRepo`
- **Injector**：Wire 生成的函数，调用所有 Provider 完成注入

```go
// wire.go —— 你写的声明文件
//go:build wireinject
// +build wireinject

func InitializeApp(dsn string) (*App, error) {
    wire.Build(
        // 基础设施
        sql.Open,
        // Repository
        repository.NewUserRepo,
        repository.NewOrderRepo,
        // Service
        service.NewUserService,
        service.NewOrderService,
        // Handler
        handler.NewUserHandler,
        handler.NewOrderHandler,
        // App
        NewApp,
    )
    return &App{}, nil
}

// 运行 wire 命令自动生成 wire_gen.go（不要手动编辑）
```

```bash
# 安装 wire
go install github.com/google/wire/cmd/wire@latest

# 生成依赖注入代码
wire

# 查看依赖图
wire ./...
```

**优点**：
- 编译时检查所有依赖是否满足，缺失或循环依赖在生成阶段就报错
- 生成的是纯 Go 代码，无运行时开销
- `wire` 命令可以展示依赖关系图

**缺点**：
- 依赖变更后需要重新运行 `wire`
- Provider 命名和组合有一定学习成本
- 不支持运行时动态注册（不过大多数项目不需要）

### 3.4 Uber Fx

**适用场景**：大型微服务框架、需要生命周期管理、动态注册（如按需加载插件）。

```go
func main() {
    fx.New(
        // Provide：注册构造函数
        fx.Provide(
            NewDB,
            NewRedis,
            repository.NewUserRepo,
            service.NewUserService,
            handler.NewUserHandler,
        ),
        // Invoke：启动时立即执行的函数（如注册路由）
        fx.Invoke(RegisterRoutes),
        // Lifecycle hooks
        fx.Invoke(func(lc fx.Lifecycle, db *sql.DB) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    return db.Ping()
                },
                OnStop: func(ctx context.Context) error {
                    return db.Close()
                },
            })
        }),
        // 日志
        fx.WithLogger(func() fxevent.Logger {
            return &fxevent.ConsoleLogger{W: os.Stdout}
        }),
    ).Run()
}
```

**优点**：
- 自动解析依赖顺序，不需要手动排列
- 内建生命周期管理（OnStart/OnStop）
- 支持 `fx.Module` 模块化组织
- 支持 annotated 类型（同类型多实例）

**缺点**：
- 运行时 DI，错误在启动时才暴露（不像 Wire 编译时就报错）
- 引入反射，有少量性能开销
- 概念较多（Provide/Invoke/Annotate/Lifecycle/Module），团队学习成本高
- "魔法"较多，排查问题不如手动注入直观

### 3.5 选型推荐

```
依赖 < 20 个  ──→  手动注入（最简单）
依赖 20~50 个 ──→  Wire（编译时安全）
依赖 > 50 个  ──→  Wire 或 Fx（Fx 更省心但需团队接受）

微服务框架    ──→  Fx（生命周期管理是刚需）
对"魔法"敏感  ──→  手动或 Wire（明确可控）
```
