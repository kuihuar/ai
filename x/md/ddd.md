在领域驱动设计（DDD）分层架构中，通常包含用户界面层（Presentation Layer）、应用层（Application Layer）、领域层（Domain Layer）和基础设施层（Infrastructure Layer）。一般是基础设施层依赖数据库（DB），下面详细解释各层的职责以及与数据库的依赖关系：

各层职责及与数据库的依赖关系
1. 用户界面层（Presentation Layer）
职责：负责与用户进行交互，将用户的请求传递给应用层，同时将应用层返回的结果展示给用户。它可以是 Web 界面、移动应用界面、命令行界面等。
与数据库的依赖关系：该层不直接依赖数据库。它的核心在于处理用户交互，不涉及数据库操作的细节。例如，在一个电商系统的 Web 界面中，用户点击“查看商品详情”按钮，用户界面层会将这个请求传递给应用层进行处理，而不是直接去数据库查询商品信息。
2. 应用层（Application Layer）
职责：协调领域层的领域对象来完成具体的业务流程，负责处理业务用例和事务管理。它不包含业务规则，只是调用领域层的服务和方法来实现业务逻辑。
与数据库的依赖关系：应用层不直接依赖数据库。它主要负责调用领域层的服务，将不同的领域操作组合起来完成一个完整的业务流程。例如，在创建订单的业务流程中，应用层会调用领域层的订单服务和商品服务，确保订单创建过程中的业务规则得到正确执行。
3. 领域层（Domain Layer）
职责：包含领域模型（实体、值对象、聚合根等）、领域服务和领域事件等，是业务逻辑的核心所在。它定义了业务规则和业务行为，负责处理业务的核心逻辑。
与数据库的依赖关系：领域层本身不直接依赖数据库。它专注于业务逻辑的实现，而不关心数据的持久化。例如，在一个图书管理系统中，领域层会定义图书实体和借阅服务，处理图书的借阅和归还等业务规则，但不会直接与数据库交互。
4. 基础设施层（Infrastructure Layer）
职责：为其他层提供通用的技术支持，包括数据库访问、消息队列、缓存等。它实现了领域层和应用层所需要的各种基础设施服务。
与数据库的依赖关系：基础设施层直接依赖数据库。它负责封装数据库的操作，将领域层的实体对象持久化到数据库中，以及从数据库中检索数据并转换为领域对象。例如，在基础设施层中会有具体的仓储（Repository）实现，用于对数据库中的数据进行增删改查操作。
示例代码说明 DDD 分层架构的依赖关系

```go
// 领域层（Domain Layer）
package domain

// Book 图书实体
type Book struct {
    ID    int
    Title string
}

// BookRepository 图书仓储接口
type BookRepository interface {
    GetBookByID(id int) (*Book, error)
    SaveBook(book *Book) error
}

// 应用层（Application Layer）
package application

import (
    "your_project/domain"
)

// BookService 图书应用服务
type BookService struct {
    bookRepository domain.BookRepository
}

// NewBookService 创建图书应用服务实例
func NewBookService(bookRepository domain.BookRepository) *BookService {
    return &BookService{
       bookRepository: bookRepository,
    }
}

// GetBookByID 根据图书 ID 获取图书信息
func (b *BookService) GetBookByID(id int) (*domain.Book, error) {
    return b.bookRepository.GetBookByID(id)
}

// 基础设施层（Infrastructure Layer）
package infrastructure

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "your_project/domain"
)

// MySQLBookRepository MySQL 图书仓储实现
type MySQLBookRepository struct {
    db *sql.DB
}

// NewMySQLBookRepository 创建 MySQL 图书仓储实例
func NewMySQLBookRepository(db *sql.DB) *MySQLBookRepository {
    return &MySQLBookRepository{
       db: db,
    }
}

// GetBookByID 根据图书 ID 从数据库中获取图书信息
func (m *MySQLBookRepository) GetBookByID(id int) (*domain.Book, error) {
    var book domain.Book
    err := m.db.QueryRow("SELECT id, title FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title)
    if err != nil {
       if err == sql.ErrNoRows {
          return nil, nil
       }
       return nil, err
    }
    return &book, nil
}

// SaveBook 将图书信息保存到数据库中
func (m *MySQLBookRepository) SaveBook(book *domain.Book) error {
    _, err := m.db.Exec("INSERT INTO books (id, title) VALUES (?, ?)", book.ID, book.Title)
    return err
}

// 用户界面层（Presentation Layer）
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "your_project/application"
    "your_project/infrastructure"
)

func main() {
    // 建立数据库连接
    db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/test")
    if err != nil {
       fmt.Println("数据库连接失败:", err)
       return
    }
    defer db.Close()

    // 创建基础设施层的图书仓储实例
    bookRepository := infrastructure.NewMySQLBookRepository(db)
    // 创建应用层的图书服务实例
    bookService := application.NewBookService(bookRepository)

    // 调用应用层的服务方法
    book, err := bookService.GetBookByID(1)
    if err != nil {
       fmt.Println("获取图书信息失败:", err)
       return
    }
    if book != nil {
       fmt.Printf("图书 ID: %d, 书名: %s\n", book.ID, book.Title)
    } else {
       fmt.Println("未找到该图书")
    }
}
```


在这个示例中，基础设施层的 MySQLBookRepository 直接依赖数据库进行数据的持久化和查询操作，而应用层和领域层通过接口与基础设施层进行交互，不直接依赖数据库。这样的分层架构使得各层职责清晰，提高了代码的可维护性和可扩展性。