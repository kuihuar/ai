1. 接口基础
- 接口定义：接口是一组方法签名的集合，用于定义行为而非数据。
- 隐式实现：类型无需显式声明实现接口，只需实现接口的全部方法。
- 空接口（interface{}): 空接口可表示任意类型，常用于泛型占位 TODO

2. 接口底层实现
- iface：非空接口的底层结构，包含动态类型（_type）和动态值（data指针）。
- eface：空接口的底层结构，仅包含动态类型和值。
- 方法调用机制：通过接口调用方法时，Go 在运行时查找动态类型的方法表，实现多态。

3. 接口高级特性
- 类型断言（Type Assertion）检查接口值的实际类型并转换
```go
value, ok := any.(int)  // 安全断言
if ok {
    fmt.Println(value)
}
```
- 类型选择（Type Switch）
```go
switch v := any.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
default:
    fmt.Println("unknown type")
}
```
4. 接口设计模式
- 小接口原则, 最佳实践：定义单一职责的小接口 
- 接口组合, 扩展行为：通过嵌入接口组合功能。

5. 常见问题
Q1：接口的 nil 值问题
```go
var w Writer  // 接口变量 w 的初始值为 nil
var f *File   // f 是一个 nil 指针
w = f         // w 的动态类型是 *File，动态值为 nil
fmt.Println(w == nil) // 输出 false
```
解释：接口是否为 nil 取决于其动态类型和值是否均为 nil。

Q2：接口与性能
方法调用开销：接口方法调用需通过动态派发，比直接调用略慢，但差异通常可忽略。

Q3：接口的零值
默认值：未初始化的接口变量为 nil，调用方法会触发 panic。

6. 实际应用场景

标准库示例
sort.Interface：需实现 Len(), Less(i, j int), Swap(i, j int) 以支持自定义排序。

error 接口：

```go
type error interface {
    Error() string
}
```

依赖注入
解耦组件：通过接口定义服务，实现模块间松耦合。

```go
type Database interface {
    Query(query string) (Result, error)
}
func NewService(db Database) *Service { /*...*/ }
```

依赖注入详细例子
依赖注入的不同方式
- 构造函数注入
- 方法注入
- 属性注入

在Go中通常推荐使用构造函数注入，因为清晰且易于管理。同时，解释为什么接口在这里是关键，因为它定义了契约，允许不同的实现替换。

依赖注入带来的好处，如提高代码的可维护性、可测试性，以及如何促进模块化开发。可能还需要提到一些相关的设计模式，如工厂模式，或者Go中的最佳实践，比如避免全局状态，使用接口进行抽象

依赖注入（Dependency Injection，DI）和解耦组件是现代软件设计的核心思想，尤其在 Go 中通过接口实现非常直观。

实现步骤： 

第一步：定义接口（解耦的关键）
通过接口定义存储层的契约，而不是依赖具体实现
```go
// user_repository.go
package user

// 定义存储接口
type Repository interface {
    CreateUser(user *User) error
    FindUserByID(id string) (*User, error)
}

// User 实体定义
type User struct {
    ID    string
    Name  string
    Email string
}
```
第二步：实现具体存储（如 MySQL）

```go
// mysql_repository.go
package user

import "database/sql"

// MySQL 实现接口
type MySQLRepository struct {
    db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
    return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateUser(user *User) error {
    // 具体 SQL 操作
    _, err := r.db.Exec("INSERT INTO users (...) VALUES (...)", user.ID, user.Name, user.Email)
    return err
}

func (r *MySQLRepository) FindUserByID(id string) (*User, error) {
    // 查询逻辑...
    return &User{ID: id, Name: "John"}, nil
}
```
第三步：业务服务层（依赖接口）


```go
// user_service.go
package user

// 业务服务层，依赖接口而非具体实现
type Service struct {
    repo Repository // 关键：这里引用的是接口
}

// 通过构造函数注入依赖（依赖注入的典型方式）
func NewUserService(repo Repository) *Service {
    return &Service{repo: repo}
}

// 业务方法
func (s *Service) RegisterUser(name, email string) (*User, error) {
    user := &User{
        ID:    generateID(),
        Name:  name,
        Email: email,
    }
    if err := s.repo.CreateUser(user); err != nil {
        return nil, err
    }
    return user, nil
}
```
解耦的实际效果
1. 更换存储实现
假设需要切换到 MongoDB，只需新增一个实现相同接口的类：
```go

// mongo_repository.go
package user

type MongoRepository struct {
    collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
    return &MongoRepository{collection: collection}
}

func (r *MongoRepository) CreateUser(user *User) error {
    // MongoDB 插入操作
    _, err := r.collection.InsertOne(context.TODO(), user)
    return err
}
```
业务层无需修改任何代码：
```go
// 初始化时注入 MongoDB 实现
mongoRepo := user.NewMongoRepository(mongoCollection)
userService := user.NewUserService(mongoRepo) // 同一接口，无缝切换
```
2. 单元测试
通过 Mock 存储层，无需连接真实数据库：

```go
// user_service_test.go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// 定义一个 Mock 存储实现
type MockRepository struct {
    users map[string]*User
}

func NewMockRepository() *MockRepository {
    return &MockRepository{users: make(map[string]*User)}
}

func (m *MockRepository) CreateUser(user *User) error {
    m.users[user.ID] = user
    return nil
}

func TestRegisterUser(t *testing.T) {
    mockRepo := NewMockRepository()
    service := NewUserService(mockRepo) // 注入 Mock

    user, err := service.RegisterUser("Alice", "alice@example.com")
    assert.NoError(t, err)
    assert.Equal(t, "Alice", user.Name)

    // 验证 Mock 存储中是否保存了用户
    savedUser, ok := mockRepo.users[user.ID]
    assert.True(t, ok)
    assert.Equal(t, "alice@example.com", savedUser.Email)
}
```


依赖注入的多种方式
Go 中常见的依赖注入方式：

方式	示例代码	适用场景
构造函数注入	func NewService(repo Repository) *Service	主要依赖，推荐优先使用
方法注入	func (s *Service) SetRepository(repo Repository)	可选依赖或动态替换依赖
全局容器	使用依赖注入框架（如 wire）	大型项目，依赖关系复杂
设计优势
可测试性：通过 Mock 轻松隔离测试业务逻辑。

可扩展性：新增存储实现无需修改业务代码。

可维护性：组件职责清晰，符合单一职责原则（SRP）。

团队协作：接口先行定义，不同开发者可并行开发存储层和业务层。

总结
通过 接口定义契约 + 依赖注入实现，Go 可以优雅地实现组件解耦。这种模式在以下场景尤为重要：

需要支持多数据源（如 MySQL、MongoDB、缓存）。

编写单元测试时隔离外部依赖。

构建插件化架构（如可替换的算法模块）。
