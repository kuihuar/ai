行为型模式

策略模式（Strategy Pattern）
作用：定义算法族，使其可互换，客户端根据场景选择不同算法。

Go 实现：通过接口定义策略，具体策略实现接口。
```go
// PaymentStrategy 是一个接口，它定义了一个 Pay 方法，该方法接受一个 float64 类型的参数 amount，表示支付的金额，并返回一个 string 类型的结果，用于描述支付的信息。这个接口是策略模式的核心，所有具体的支付策略都需要实现这个接口。
type PaymentStrategy interface {
    Pay(amount float64) string
}

// 具体的支付策略，它们都实现了 PaymentStrategy 接口
type CreditCard struct{}
func (c CreditCard) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f via Credit Card", amount)
}

// 具体的支付策略，它们都实现了 PaymentStrategy 接口
type PayPal struct{}
func (p PayPal) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f via PayPal", amount)
}

// PaymentContext 是一个上下文结构体，它持有一个 PaymentStrategy 类型的字段 strategy，用于存储当前使用的支付策略
type PaymentContext struct {
    strategy PaymentStrategy
}
// 实现上下文
// 方法用于设置当前使用的支付策略。通过这个方法，我们可以在运行时动态地改变支付策略。
// Execute 方法用于执行支付操作。它调用当前支付策略的 Pay 方法，并返回支付结果。
func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}
func (p *PaymentContext) Execute(amount float64) string {
    return p.strategy.Pay(amount)
}

// 使用
ctx := &PaymentContext{}
ctx.SetStrategy(CreditCard{})
fmt.Println(ctx.Execute(100.5)) // 输出: Paid $100.50 via Credit Card
```