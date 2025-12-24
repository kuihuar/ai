#### 四层目录结构

my-project/
├── cmd/                 # 第1层：应用入口（Main Entrypoints）
│   └── myapp/
│       └── main.go      # 主程序入口（HTTP、gRPC、CLI等）
├── internal/            # 第2层：内部实现（私有业务逻辑）
│   ├── handler/         # 接口层（HTTP路由、gRPC服务）
│   ├── service/         # 业务逻辑层（核心功能实现）
│   ├── repository/      # 数据访问层（数据库、缓存操作）
│   └── model/           # 数据模型（结构体定义、DTOs）
├── pkg/                 # 第3层：公共库（可复用的模块）
│   ├── utils/           # 工具函数（日志、加密等）
│   └── config/          # 配置管理（解析YAML/Env）
└── api/                 # 第4层：API协议（接口定义）
    ├── proto/           # Protobuf文件（gRPC协议）
    └── openapi/         # OpenAPI/Swagger文档（REST协议）
└── scripts/
└── docs/     


#### DDD 目录结构示例

my-ddd-project/
├── cmd/                     # 应用入口（启动 HTTP、CLI、定时任务等）
│   └── api/
│       └── main.go          # 主程序入口
├── internal/                # 内部代码（不对外暴露）
│   ├── user/                # 用户领域模块（按业务划分）
│   │   ├── domain/          # 领域层（核心）
│   │   │   ├── entity/      # 领域实体（如 User）
│   │   │   ├── valueobject/ # 值对象（如 Email）
│   │   │   ├── service/     # 领域服务（纯业务逻辑）
│   │   │   └── event/       # 领域事件（如 UserRegistered）
│   │   ├── application/     # 应用层（用例编排）
│   │   │   └── service/     # 应用服务（如 UserAppService）
│   │   └── interfaces/      # 接口层（适配外部）
│   │       ├── http/        # HTTP 控制器（如 Gin 路由）
│   │       ├── grpc/        # gRPC 服务端
│   │       └── repository/  # 仓储接口定义（依赖倒置）
│   ├── order/               # 订单领域模块（结构同上）
│   └── shared/              # 跨领域共享代码
│       └── kernel/          # 通用领域基础设施（如 ID 生成器）
├── pkg/                     # 可复用的公共库
│   ├── errors/              # 自定义错误类型
│   └── utils/               # 工具函数（如加密、日期处理）
└── infra/                   # 基础设施层（具体实现）
    ├── persistence/         # 数据持久化
    │   ├── mysql/           # MySQL 仓储实现
    │   └── redis/           # Redis 缓存实现
    ├── mq/                  # 消息队列实现（如 Kafka）
    └── config/              # 配置加载（Env、YAML 等）



#### 基于特性的组织方式 (Feature-Based Organization):

这种结构按功能而不是按技术层对代码进行分组。
重点： 围绕应用程序的特定功能或用例组织代码。
优点：
改进的代码局部性：与功能相关的所有代码都位于一个位置，从而更容易理解和维护。
更容易添加或删除功能：可以添加或删除功能，而不会影响应用程序的其他部分。
降低耦合性：功能之间的耦合性更低，因为它们仅依赖于共享的 pkg/ 目录。
缺点：
代码重复的潜在性：如果多个功能需要类似的功能，则代码可能会在功能目录中重复。 将通用功能仔细提取到 pkg/ 目录中至关重要。
在大型项目中可能会变得复杂：随着功能数量的增长，根目录可能会变得混乱。
myproject/
├── feature1/
│   ├── api/
│   │   └── handler.go
│   ├── domain/
│   │   └── model.go
│   ├── service/
│   │   └── service.go
│   ├── repository/
│   │   └── repository.go
│   └── ...
├── feature2/
│   ├── api/
│   │   └── handler.go
│   ├── domain/
│   │   └── model.go
│   ├── service/
│   │   └── service.go
│   ├── repository/
│   │   └── repository.go
│   └── ...
├── cmd/
│   └── myapp/
│       └── main.go
├── pkg/
│   └── ...
└── vendor/
    └── ...    

#### 模块化单体结构 (Modular Monolith Structure):

这种结构适用于作为单个单元部署但内部组织成模块的较大型应用程序。
重点： 将应用程序划分为具有明确边界的独立模块。
优点：
改进的代码组织：模块提供了清晰的关注点分离，使应用程序更易于理解和维护。
增加的可重用性：模块可以在应用程序的不同部分重复使用。
更容易扩展：如果需要，可以独立扩展模块（尽管作为单个单元部署）。
缺点：
需要仔细规划：定义模块边界可能具有挑战性。
潜在的紧密耦合：如果模块设计不当，可能会变得紧密耦合。
myproject/
├── module1/
│   ├── internal/
│   │   └── ... (DDD-like structure)
│   ├── api/
│   │   └── ...
│   └── ...
├── module2/
│   ├── internal/
│   │   └── ...
│   ├── api/
│   │   └── ...
│   └── ...
├── cmd/
│   └── myapp/
│       └── main.go
├── pkg/
│   └── ...
└── vendor/
    └── ...


#### 干净架构（洋葱架构）(Clean Architecture / Onion Architecture):

这种结构强调依赖倒置和可测试性。
重点： 使核心业务逻辑独立于框架、数据库和外部服务。
优点：
高可测试性：核心应用程序逻辑独立于基础设施细节，易于测试。
灵活性：应用程序可以轻松适应不同的框架、数据库或外部服务。
可维护性：清晰的关注点分离使应用程序更易于理解和维护。
缺点：
可能很复杂：分层架构可能会增加复杂性，尤其是在小型项目中。
需要仔细规划：定义层之间的边界可能具有挑战性。
myproject/
├── app/
│   ├── entities/
│   ├── usecases/
│   ├── interfaces/
│   │   ├── controllers/
│   │   ├── presenters/
│   │   └── repositories/
│   └── ...
├── infrastructure/
│   ├── persistence/
│   ├── messaging/
│   ├── api/
│   └── ...
├── cmd/
│   └── myapp/
│       └── main.go
├── pkg/
│   └── ...
└── vendor/
    └── ...

### 简单分层架构 (Simple Layered Architecture):

这是分层架构的简化版本，适用于较小的项目。
重点： 提供基本的关注点分离，且复杂性最低。
优点：
简单易懂。
提供基本的关注点分离。
缺点：
在大型项目中可能难以维护。
不如更复杂的架构灵活。
myproject/
├── api/
│   └── handlers.go
├── service/
│   └── services.go
├── repository/
│   └── repositories.go
├── model/
│   └── models.go
├── cmd/
│   └── myapp/
│       └── main.go
├── pkg/
│   └── ...
└── vendor/
    └── ...


#### 六边形架构（端口和适配器）(Hexagonal Architecture / Ports and Adapters):

这种架构通过定义应用程序核心和外部系统之间的清晰边界来强调可测试性和解耦。
重点： 通过使用端口和适配器将应用程序核心与外部系统解耦。
优点：
高可测试性：应用程序核心独立于外部系统，易于测试。
灵活性：可以通过更换适配器轻松地将应用程序适应不同的技术。
可维护性：清晰的关注点分离使应用程序更易于理解和维护。
缺点：
可能很复杂：定义端口和适配器可能具有挑战性。
需要仔细规划：该架构需要仔细规划，以确保核心和适配器之间的边界定义明确。
myproject/
├── core/
│   ├── 领域/
│   ├── 用例/
│   ├── 端口/
│   │   ├── 输入/
│   │   └── 输出/
│   └── ...
├── adapters/
│   ├── api/
│   ├── 持久化/
│   ├── 消息传递/
│   └── ...
├── cmd/
│   └── myapp/
│       └── main.go
├── pkg/
│   └── ...
└── vendor/
    └── ...    






1. 四层目录结构示例 (不严格遵循 DDD):

假设我们正在构建一个简单的任务管理应用程序，允许用户创建、更新和删除任务。

重点:

提供一种简单的、分层的代码组织方式。
将应用程序划分为 API 层、服务层、数据访问层和模型层。
易于理解和实现，尤其是在小型项目中。
优点:

简单易懂: 结构简单，容易上手，适合快速开发小型项目。
基本的分层: 提供了一定的关注点分离，有助于代码的组织和维护。
开发速度快: 由于结构简单，可以快速搭建项目框架并开始开发。
缺点:

缺乏领域模型: 领域模型通常比较简单，缺乏丰富的行为和业务逻辑。
业务逻辑集中在 Service 层: 导致 Service 层过于臃肿，难以维护。
可测试性较差: 由于层之间的依赖关系比较紧密，难以进行单元测试。
可扩展性差: 随着项目规模的增大，结构会变得难以维护和扩展。
与业务脱节: 代码结构可能与实际业务领域不匹配，导致理解和维护困难。
不适合复杂项目: 难以应对复杂的业务逻辑和需求变化。

复制
taskmanager/
├── cmd/
│   └── taskmanager/
│       └── main.go  // 应用程序入口点
├── api/
│   └── handlers.go  // HTTP API 处理程序
├── service/
│   └── task_service.go  // 任务管理业务逻辑
├── repository/
│   └── task_repository.go  // 数据访问逻辑
├── model/
│   └── task.go  // 任务数据模型
├── config/
│   └── config.go  // 应用程序配置
├── pkg/
│   └── utils.go  // 通用工具函数
├── vendor/
│   └── (依赖)
├── Dockerfile
├── Makefile
└── README.md    

2. DDD 目录结构示例 (更接近 DDD):

现在，让我们将相同的任务管理应用程序使用 DDD 原则进行组织。

重点:

将代码与业务领域紧密结合。
使用领域模型来表示业务概念和规则。
通过命令、查询和事件来驱动应用程序的行为。
强调关注点分离和可测试性。
优点:

与业务对齐: 代码结构反映了业务领域，更容易理解和维护。
丰富的领域模型: 领域模型包含丰富的行为和业务逻辑，更贴近实际业务。
高可测试性: 领域模型可以独立于基础设施细节进行测试。
高可维护性: 分层架构和关注点分离使应用程序更易于修改和扩展。
更好的可扩展性: 能够更好地应对复杂的业务逻辑和需求变化。
更强的内聚性： 每个层都有明确的职责，并且层之间的依赖关系清晰。
缺点:

复杂性高: 需要理解 DDD 的核心概念，并进行领域建模。
学习曲线陡峭: 需要开发人员掌握 DDD 的相关知识和技能。
开发速度慢: 需要花费更多的时间进行领域建模和代码设计。
过度设计风险: 对于简单的应用程序，DDD 可能过度复杂。
需要领域专家参与: 需要领域专家的参与，才能构建准确的领域模型。
初始成本高: 需要投入更多的时间和精力进行项目初始化和架构设计。

复制
taskmanager/
├── cmd/
│   └── taskmanager/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── task/
│   │   │   ├── task.go  // 任务实体 (聚合根)
│   │   │   ├── task_id.go  // 任务 ID 值对象
│   │   │   ├── task_status.go  // 任务状态值对象
│   │   │   ├── ...
│   │   ├── events/
│   │   │   ├── task_created_event.go
│   │   │   ├── task_updated_event.go
│   │   │   ├── ...
│   │   ├── services/
│   │   │   ├── task_domain_service.go  // 任务领域服务
│   │   │   ├── ...
│   ├── application/
│   │   ├── commands/
│   │   │   ├── create_task_command.go
│   │   │   ├── update_task_command.go
│   │   │   ├── ...
│   │   ├── queries/
│   │   │   ├── get_task_query.go
│   │   │   ├── list_tasks_query.go
│   │   │   ├── ...
│   │   ├── services/
│   │   │   ├── task_application_service.go  // 任务应用服务
│   │   │   ├── ...
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── task_repository.go  // 任务仓库接口
│   │   │   ├── task_repository_impl.go  // 任务仓库实现
│   │   │   ├── ...
│   │   ├── api/
│   │   │   ├── task_handler.go  // HTTP API 处理程序
│   │   │   ├── ...
├── pkg/
│   └── utils.go
├── vendor/
│   └── (依赖)
├── Dockerfile
├── Makefile
└── README.md

|架构类型|适用场景|核心优势|主要挑战|
|---|---|---|---|
|四层目录结构|小型项目、快速原型|简单易用|业务复杂后耦合度高|
|DDD 目录结构|复杂业务系统、微服务|业务与技术解耦|领域建模成本高|
|基于特性的组织方式|功能明确的中型应用|模块独立开发|重复代码风险|
|模块化单体结构|渐进式迁移微服务|平衡灵活性与简单性|依赖管理复杂|
|干净架构/洋葱架构|高可维护性核心业务系统|技术无关性|适配器冗余|
|简单分层架构（三层架构）|CRUD 应用、入门项目|结构清晰|贫血模型、扩展性差|
|六边形架构|多外部依赖的中大型系统|高度解耦与扩展性|设计复杂度高|

选型建议
- 初创项目：优先选择 简单分层 或 基于特性，快速验证业务。
- 复杂业务：采用 DDD 或 干净架构，确保领域模型清晰。
- 技术多样性：六边形架构 适配多外部服务场景。
- 平滑演进：模块化单体 作为微服务过渡方案。


干净架构和DDD

特性	干净架构	DDD
核心目标	技术解耦，依赖倒置	业务建模，统一语言
目录划分依据	功能分层（实体、用例等）	业务模块（user/order）
适用场景	技术栈复杂或需频繁替换	业务复杂且需持续演进
代码示例重点	分层接口与适配器	领域模型与限界上下文
学习曲线	中等	高（需掌握DDD方法论



---

### 一、干净架构（Clean Architecture）示例
#### 目录结构
bash
my-clean-project/
├── cmd/                     # 应用入口（启动HTTP/CLI）
│   └── main.go
├── internal/                # 核心业务逻辑（不对外暴露）
│   ├── entity/              # 实体层（纯业务模型）
│   │   └── user.go          # 用户实体定义
│   ├── usecase/             # 用例层（业务逻辑编排）
│   │   └── user_usecase.go  # 用户相关用例
│   └── repository/          # 仓储接口（抽象数据访问）
│       └── user_repository.go
├── pkg/                     # 公共库（可复用）
│   └── database/            # 数据库连接池等工具
└── infra/                   # 基础设施层（技术实现）
    ├── http/                # HTTP服务实现
    │   └── handler.go       # HTTP路由和控制器
    └── persistence/         # 数据持久化实现
        └── mysql/           # MySQL仓储实现
            └── user_repository.go
#### 关键代码示例

1. 实体层（Entity）

```go
// internal/entity/user.go
package entity

type User struct {
    ID    string
    Name  string
    Email string
}

// 业务规则：邮箱格式校验
func (u *User) ValidateEmail() bool {
    return strings.Contains(u.Email, "@")
}
```
2. 用例层（Usecase）

```go
// internal/usecase/user_usecase.go
package usecase

type UserUsecase struct {
    repo repository.UserRepository
}

// 创建用户的业务逻辑
func (uc *UserUsecase) CreateUser(name, email string) error {
    user := entity.User{Name: name, Email: email}
    if !user.ValidateEmail() {
        return errors.New("invalid email")
    }
    return uc.repo.Save(&user)
}
```
3. 仓储接口（Repository Interface）

```go
// internal/repository/user_repository.go
package repository

type UserRepository interface {
    Save(user *entity.User) error
    FindByID(id string) (*entity.User, error)
}
```
4. 基础设施层（MySQL仓储实现）

```go
// infra/persistence/mysql/user_repository.go
package mysql

type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Save(user *entity.User) error {
    return r.db.Create(user).Error
}
```
5. HTTP控制器（适配器层）

```go
// infra/http/handler.go
package http

type UserHandler struct {
    userUsecase *usecase.UserUsecase
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
```
### 二、领域驱动设计（DDD）示例
#### 目录结构
bash
my-ddd-project/
├── cmd/
│   └── main.go
├── internal/
│   ├── user/                # 用户领域模块
│   │   ├── domain/          # 领域层（核心）
│   │   │   ├── entity/      # 用户实体
│   │   │   ├── service/     # 领域服务
│   │   │   └── event/       # 领域事件（如UserRegistered）
│   │   ├── application/     # 应用层（用例编排）
│   │   │   └── service/     # 应用服务
│   │   └── interfaces/      # 接口层
│   │       ├── http/        # HTTP控制器
│   │       └── repository/  # 仓储接口
│   └── order/               # 订单领域模块（结构同上）
├── pkg/                     # 公共库
│   └── eventbus/            # 事件总线实现
└── infra/
    ├── http/                # HTTP框架（如Gin）
    ├── persistence/         # 数据持久化
    └── mq/                  # 消息队列实现

#### 关键代码示例
1. 领域实体（Entity）

```go
// internal/user/domain/entity/user.go
package entity

type User struct {
    ID    string
    Name  string
    Email string
}

func (u *User) ChangeName(newName string) {
    u.Name = newName
}
```
2. 领域服务（Domain Service）

```go
// internal/user/domain/service/registration.go
package service

type RegistrationService struct {
    repo        repository.UserRepository
    eventPub    eventbus.Publisher
}

func (s *RegistrationService) Register(user *entity.User) error {
    if exists, _ := s.repo.Exists(user.Email); exists {
        return errors.New("email already exists")
    }
    if err := s.repo.Save(user); err != nil {
        return err
    }
    s.eventPub.Publish(event.UserRegistered{UserID: user.ID})
    return nil
}
```
3. 应用服务（Application Service）

```go
// internal/user/application/service/user_app_service.go
package service

type UserAppService struct {
    regService *domain.RegistrationService
}

func (s *UserAppService) CreateUser(name, email string) error {
    user := entity.User{Name: name, Email: email}
    return s.regService.Register(&user)
}
```
3. 仓储接口（Repository Interface）

```go
// internal/user/interfaces/repository/user_repository.go
package repository

type UserRepository interface {
    Save(user *entity.User) error
    Exists(email string) (bool, error)
}
```
4. MySQL仓储实现

```go
// infra/persistence/mysql/user_repository.go
package mysql

type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Save(user *entity.User) error {
    return r.db.Create(user).Error
}
```