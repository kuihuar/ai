### 干净架构（Clean Architecture）
干净架构由 Robert C. Martin（Uncle Bob）提出，是一种以 业务逻辑为核心、技术实现为外围 的软件设计模式。其核心思想是通过分层和依赖规则，隔离业务逻辑与技术细节，使系统具备 高可维护性、可测试性 和 技术无关性。

#### 一、核心分层与依赖规则
干净架构通过 同心圆分层 实现依赖倒置，内层定义业务规则，外层实现技术细节。各层从内到外依次为：

1. 实体层（Entities）

    - 职责：定义核心业务模型（数据 + 业务规则）。

    - 示例：用户实体、订单实体及其验证逻辑。

    - 特点：完全独立于框架、数据库和外部服务。

2. 用例层（Use Cases）

    - 职责：编排实体完成具体业务场景（如“创建订单”、“支付流程”）。

    - 示例：调用仓储接口保存数据，发布领域事件。

    - 特点：通过接口与外部交互，不依赖具体实现。

3. 接口适配器层（Interface Adapters）

    - 职责：转换数据格式，适配外部系统（如数据库、UI、第三方API）。

    - 示例：HTTP 控制器、数据库仓储实现、消息队列生产者。

    - 特点：将外部数据转换为用例层所需格式，或反向转换。

4. 框架与驱动层（Frameworks & Drivers）

    - 职责：实现具体技术细节（如 Web 框架、数据库驱动）。

    - 示例：Gin 路由配置、MySQL 连接池、Kafka 客户端。

    - 特点：最外层，可随时替换不影响核心业务。

依赖规则
1. 单向依赖：外层可依赖内层，内层 绝不依赖 外层。
2. 依赖倒置：通过接口（如仓储接口）实现外层对內层的依赖。

#### 二、实现示例
1. 实体层（Entities）
定义核心业务模型和规则，无任何外部依赖。

```go
// internal/entity/user.go
package entity

type User struct {
    ID    string
    Name  string
    Email string
}

// 业务规则：邮箱格式校验
func (u *User) ValidateEmail() error {
    if !strings.Contains(u.Email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}
```
2. 用例层（Use Cases）
通过接口调用外部服务，实现业务逻辑。

```go
// internal/usecase/user_usecase.go
package usecase

type UserRepository interface {
    Save(user *entity.User) error
}

type UserUsecase struct {
    repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
    return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) CreateUser(name, email string) error {
    user := &entity.User{Name: name, Email: email}
    if err := user.ValidateEmail(); err != nil {
        return err
    }
    return uc.repo.Save(user) // 依赖倒置：通过接口调用仓储
}
```

3. 接口适配器层（Interface Adapters）
实现用例层定义的接口，适配具体技术。

```go
// infra/persistence/mysql/user_repository.go
package mysql

import (
    "myapp/internal/entity"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// 实现用例层的 UserRepository 接口
func (r *UserRepository) Save(user *entity.User) error {
    return r.db.Create(user).Error
}
```
4. 框架与驱动层（Frameworks & Drivers）
配置技术组件，启动服务。

```go
// infra/http/server.go
package http

import (
    "myapp/internal/usecase"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userUsecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUsecase: uc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req struct { Name string; Email string }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "bad request"})
        return
    }
    if err := h.userUsecase.CreateUser(req.Name, req.Email); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"status": "ok"})
}

// 主函数中启动服务
func main() {
    db := ConnectMySQL() // 初始化数据库
    userRepo := mysql.NewUserRepository(db)
    userUsecase := usecase.NewUserUsecase(userRepo)
    handler := NewUserHandler(userUsecase)

    router := gin.Default()
    router.POST("/users", handler.CreateUser)
    router.Run(":8080")
}
```

#### 三、干净架构的核心优势
1. 技术无关性
    - 替换技术栈：更换数据库（如 MySQL → PostgreSQL）或 Web 框架（如 Gin → Echo）时，只需修改适配器层，无需改动业务逻辑。

    - 示例：替换仓储实现：

```go
// infra/persistence/mongodb/user_repository.go
func NewMongoUserRepository(client *mongo.Client) *MongoUserRepository {
    return &MongoUserRepository{col: client.Database("test").Collection("users")}
}
```
2. 高可测试性
    - 独立测试业务逻辑：通过 Mock 适配器隔离外部依赖。

```go
// 测试用例
func TestCreateUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Save", mock.Anything).Return(nil)
    
    uc := usecase.NewUserUsecase(mockRepo)
    err := uc.CreateUser("Alice", "alice@example.com")
    assert.NoError(t, err)
}
```
3. 清晰的代码边界
    - 防止代码腐化：严格的分层规则避免业务逻辑与技术代码混杂。
    - 团队协作：开发者按层分工（如领域专家负责实体层，后端工程师负责适配器层）。

#### 四、适用场景与挑战
适用场景
- 长期维护的项目：业务逻辑稳定，技术栈可能频繁变更。

- 复杂业务系统：需清晰隔离核心业务与外部依赖。

- 微服务架构：每个服务独立采用干净架构，提升整体系统灵活性。

挑战与解决方案
|挑战|解决方案|
|----|----|
|适配器代码冗余|通过代码生成工具（如 Protobuf）自动生成适配器代码。|
|学习曲线陡峭|提供分层规范文档，结合代码审查确保团队理解架构规则。|
|初期开发成本高|在复杂项目中逐步引入，优先核心模块使用干净架构，其他模块采用简单分层。|

#### 五、与其他架构对比
|架构|核心差异|
|----|----|
|传统分层架构|业务逻辑分散在各层，技术实现与业务耦合度高。|
|六边形架构|强调端口与适配器的对称性，干净架构是六边形架构的一种变体。|
|DDD|更关注领域建模和统一语言，干净架构提供技术解耦的实现框架。|

#### 六、总结
干净架构通过 分层设计 和 依赖倒置，将业务逻辑置于系统核心，使其免受技术细节变更的影响。在 Go 语言中，结合接口和依赖注入，可高效实现这一架构。尽管初期需要一定的设计成本，但其带来的 可维护性、灵活性 和 可测试性 优势，使其成为中大型项目的理想选择。

#### 七、示例项目结构
myapp/
├── cmd/                     # 应用入口
│   └── main.go              # 主程序入口（启动HTTP服务）
├── internal/                # 核心业务逻辑
│   └── user/                # 用户领域模块
│       ├── domain/          # 领域层（核心业务模型）
│       │   ├── entity/      # 用户实体
│       │   │   └── user.go
│       │   ├── valueobject/ # 值对象（如邮箱）
│       │   │   └── email.go
│       │   └── service/     # 领域服务（纯业务逻辑）
│       │       └── registration_service.go
│       ├── application/     # 应用层（用例编排）
│       │   └── service/     # 应用服务
│       │       └── user_app_service.go
│       └── interfaces/      # 接口层（适配外部系统）
│           ├── http/        # HTTP控制器
│           │   └── user_controller.go
│           └── repository/  # 仓储接口（抽象数据访问）
│               └── user_repository.go
├── pkg/                     # 公共库（可复用）
│   ├── eventbus/            # 事件总线（发布领域事件）
│   └── utils/               # 工具函数（如加密）
└── infra/                   # 基础设施层（技术实现）
    ├── config/         # 配置管理模块
    │   ├── config.go   # 配置加载逻辑
    │   └── env/        # 环境变量解析
    ├── persistence/         # 数据持久化
    │   └── mysql/           # MySQL仓储实现
    │       └── user_repository.go
    │   └── redis/           # Redis仓储实现
    │       └── user_repository.go
    ├── http/                # HTTP框架（如Gin）
    │   └── server.go
    └── event/               # 事件处理器（如发送邮件）
        └── user_registered_handler.go



        


└── infra/              # 基础设施层（框架与驱动层）
    ├── config/         # 配置管理模块
    │   ├── config.go   # 配置加载逻辑
    │   └── env/        # 环境变量解析
    ├── persistence/    # 数据持久化实现
    │   ├── mysql/      # MySQL 配置与连接池
    │   └── redis/      # Redis 配置与客户端
    └── http/           # Web框架配置（如Gin）



干净架构可以替换框架
