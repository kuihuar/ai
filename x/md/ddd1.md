在电商系统中，处理订单支付涉及多个领域对象交互，这里我们用领域服务来封装这一复杂业务逻辑。

package main

import (
	"errors"
	"fmt"
)

// Order 订单实体
type Order struct {
	OrderID    int
	TotalPrice float64
	IsPaid     bool
}

// MarkAsPaid 将订单标记为已支付
func (o *Order) MarkAsPaid() {
	o.IsPaid = true
}

// Payment 支付实体
type Payment struct {
	PaymentID int
	Amount    float64
}

// ProcessPayment 处理支付
func (p *Payment) ProcessPayment() bool {
	fmt.Printf("Processing payment of %.2f\n", p.Amount)
	return true
}

// PaymentDomainService 领域服务，处理订单支付
type PaymentDomainService struct{}

// PayOrder 处理订单支付逻辑
func (pds *PaymentDomainService) PayOrder(order *Order, payment *Payment) error {
	if order.TotalPrice != payment.Amount {
		return errors.New("payment amount does not match order total")
	}
	if payment.ProcessPayment() {
		order.MarkAsPaid()
		fmt.Printf("Order %d has been paid.\n", order.OrderID)
		return nil
	}
	return errors.New("payment failed")
}

可以使用以下方式调用上述代码进行测试：

func main() {
	order := &Order{
		OrderID:    1,
		TotalPrice: 100.0,
	}
	payment := &Payment{
		PaymentID: 101,
		Amount:    100.0,
	}
	paymentDomainService := &PaymentDomainService{}
	err := paymentDomainService.PayOrder(order, payment)
	if err != nil {
		fmt.Println(err)
	}
}
应用服务示例
创建订单的业务流程需要协调多个操作，我们使用应用服务来完成这个任务。

package main

import (
	"fmt"
)

// 假设这是基础设施层的订单仓储接口
type OrderRepository interface {
	SaveOrder(order *Order)
}

// InMemoryOrderRepository 内存中的订单仓储实现
type InMemoryOrderRepository struct{}

// SaveOrder 保存订单到“数据库”（这里只是模拟）
func (ior *InMemoryOrderRepository) SaveOrder(order *Order) {
	fmt.Printf("Saving order %d to database.\n", order.OrderID)
}

// OrderApplicationService 应用服务，处理订单创建
type OrderApplicationService struct {
	orderRepository OrderRepository
}

// NewOrderApplicationService 创建订单应用服务实例
func NewOrderApplicationService(repo OrderRepository) *OrderApplicationService {
	return &OrderApplicationService{
		orderRepository: repo,
	}
}

// CreateOrder 创建订单的业务流程
func (oas *OrderApplicationService) CreateOrder(totalPrice float64) {
	order := &Order{
		OrderID:    2,
		TotalPrice: totalPrice,
	}
	oas.orderRepository.SaveOrder(order)
	fmt.Printf("Order %d created successfully.\n", order.OrderID)
}

可以使用以下方式调用上述代码进行测试：

func main() {
	orderRepository := &InMemoryOrderRepository{}
	orderAppService := NewOrderApplicationService(orderRepository)
	orderAppService.CreateOrder(200.0)
}
代码解释
领域服务：PaymentDomainService 负责处理订单支付的核心业务逻辑，它协调了 Order 和 Payment 两个领域对象的操作，确保支付金额和订单总价匹配，并处理支付流程和订单状态更新。
应用服务：OrderApplicationService 负责协调订单创建的业务流程，它调用领域层的 Order 实体创建订单，并使用基础设施层的 OrderRepository 保存订单，完成一个完整的业务用例。
通过这种分层设计，不同层的职责清晰，提高了代码的可维护性和可扩展性。

代码目录是啥样子的？

WeTab AI-4o mini
根据上述代码示例，结合 Go 语言的项目结构最佳实践，一个典型的目录结构可以如下设计，以更好地组织领域服务、应用服务以及相关实体和仓储等代码：

.
├── cmd
│   └── main.go
├── internal
│   ├── domain
│   │   ├── order.go
│   │   ├── payment.go
│   │   └── payment_domain_service.go
│   ├── application
│   │   └── order_application_service.go
│   └── infrastructure
│       └── order_repository.go
└── go.mod
下面对各部分目录和文件进行详细解释：

1. cmd 目录
用途：存放项目的入口文件。通常每个可执行程序对应一个子目录，这里简单地将主程序入口文件 main.go 放在该目录下。
main.go：程序的入口点，负责初始化应用服务、领域服务以及仓储等组件，并调用相应的服务方法来启动业务流程。示例代码如下：
package main

import (
    "github.com/yourproject/internal/application"
    "github.com/yourproject/internal/infrastructure"
)

func main() {
    orderRepository := infrastructure.NewInMemoryOrderRepository()
    orderAppService := application.NewOrderApplicationService(orderRepository)
    orderAppService.CreateOrder(200.0)

    // 测试领域服务
    // ...
}
2. internal 目录
用途：存放项目内部使用的代码，这些代码不对外暴露，只能在项目内部被引用。按照不同的层次进行划分，包括领域层、应用层和基础设施层。
2.1 domain 目录
用途：包含领域层的代码，如实体和领域服务。
order.go：定义 Order 实体及其相关方法。
package domain

type Order struct {
    OrderID    int
    TotalPrice float64
    IsPaid     bool
}

func (o *Order) MarkAsPaid() {
    o.IsPaid = true
}
payment.go：定义 Payment 实体及其相关方法。
package domain

import "fmt"

type Payment struct {
    PaymentID int
    Amount    float64
}

func (p *Payment) ProcessPayment() bool {
    fmt.Printf("Processing payment of %.2f\n", p.Amount)
    return true
}
payment_domain_service.go：实现订单支付的领域服务。
package domain

import (
    "errors"
    "fmt"
)

type PaymentDomainService struct{}

func (pds *PaymentDomainService) PayOrder(order *Order, payment *Payment) error {
    if order.TotalPrice != payment.Amount {
        return errors.New("payment amount does not match order total")
    }
    if payment.ProcessPayment() {
        order.MarkAsPaid()
        fmt.Printf("Order %d has been paid.\n", order.OrderID)
        return nil
    }
    return errors.New("payment failed")
}
2.2 application 目录
用途：存放应用层的代码，主要是应用服务。
order_application_service.go：实现订单创建的应用服务。
package application

import (
    "github.com/yourproject/internal/domain"
    "github.com/yourproject/internal/infrastructure"
)

type OrderApplicationService struct {
    orderRepository infrastructure.OrderRepository
}

func NewOrderApplicationService(repo infrastructure.OrderRepository) *OrderApplicationService {
    return &OrderApplicationService{
        orderRepository: repo,
    }
}

func (oas *OrderApplicationService) CreateOrder(totalPrice float64) {
    order := &domain.Order{
        OrderID:    2,
        TotalPrice: totalPrice,
    }
    oas.orderRepository.SaveOrder(order)
    fmt.Printf("Order %d created successfully.\n", order.OrderID)
}
2.3 infrastructure 目录
用途：包含基础设施层的代码，如数据库访问、文件操作等。这里主要实现订单仓储接口。
order_repository.go：定义订单仓储接口并提供一个简单的内存实现。
package infrastructure

import (
    "fmt"
    "github.com/yourproject/internal/domain"
)

type OrderRepository interface {
    SaveOrder(order *domain.Order)
}

type InMemoryOrderRepository struct{}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{}
}

func (ior *InMemoryOrderRepository) SaveOrder(order *domain.Order) {
    fmt.Printf("Saving order %d to database.\n", order.OrderID)
}
3. go.mod 文件
用途：用于管理项目的依赖关系，记录项目所依赖的外部包及其版本信息。可以使用 go mod init yourproject 命令来初始化该文件。
通过这样的目录结构，各个层次的代码分离清晰，便于维护和扩展。不同层之间的依赖关系也更加明确，遵循了领域驱动设计的原则。