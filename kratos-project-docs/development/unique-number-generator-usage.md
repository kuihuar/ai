# 唯一编号生成器使用示例

## 概述

本文档展示如何在不同的业务场景中使用唯一编号生成器，通过不同的前缀区分不同的业务类型。

## 基本使用

### 1. 订单号生成（前缀：ORD）

```go
// 在 OrderUsecase 中
func (uc *OrderUsecase) CreateOrder(ctx context.Context, ...) (*Order, error) {
    // 生成订单号
    orderNo, err := uc.GenerateOrderNo(ctx)  // 内部使用 "ORD" 前缀
    if err != nil {
        return nil, err
    }
    
    order := &Order{
        OrderNo: orderNo,  // 例如: ORD20241216000001
        // ...
    }
    // ...
}
```

### 2. 支付单号生成（前缀：PAY）

```go
// 在 PaymentUsecase 中
type PaymentUsecase struct {
    numberGenerator number.Generator
    // ...
}

func (uc *PaymentUsecase) CreatePayment(ctx context.Context, ...) (*Payment, error) {
    // 生成支付单号
    paymentNo, err := uc.numberGenerator.Generate(ctx, "PAY")
    if err != nil {
        return nil, err
    }
    
    payment := &Payment{
        PaymentNo: paymentNo,  // 例如: PAY20241216000001
        // ...
    }
    // ...
}
```

### 3. 退款单号生成（前缀：REF）

```go
// 在 RefundUsecase 中
type RefundUsecase struct {
    numberGenerator number.Generator
    // ...
}

func (uc *RefundUsecase) CreateRefund(ctx context.Context, ...) (*Refund, error) {
    // 生成退款单号
    refundNo, err := uc.numberGenerator.Generate(ctx, "REF")
    if err != nil {
        return nil, err
    }
    
    refund := &Refund{
        RefundNo: refundNo,  // 例如: REF20241216000001
        // ...
    }
    // ...
}
```

## 业务封装示例

### 1. 为每个业务创建专用的生成方法

```go
// internal/biz/payment.go
type PaymentUsecase struct {
    repo            PaymentRepo
    numberGenerator number.Generator
    log             *log.Helper
}

func NewPaymentUsecase(
    repo PaymentRepo,
    numberGenerator number.Generator,
    logger log.Logger,
) *PaymentUsecase {
    return &PaymentUsecase{
        repo:            repo,
        numberGenerator: numberGenerator,
        log:             log.NewHelper(logger),
    }
}

// GeneratePaymentNo 生成支付单号
func (uc *PaymentUsecase) GeneratePaymentNo(ctx context.Context) (string, error) {
    return uc.numberGenerator.Generate(ctx, "PAY")
}
```

### 2. 在业务逻辑中使用

```go
func (uc *PaymentUsecase) CreatePayment(ctx context.Context, orderID int64, amount int64) (*Payment, error) {
    // 生成支付单号
    paymentNo, err := uc.GeneratePaymentNo(ctx)
    if err != nil {
        return nil, err
    }
    
    payment := &Payment{
        PaymentNo: paymentNo,
        OrderID:   orderID,
        Amount:    amount,
        Status:    PaymentStatusPending,
        CreatedAt: time.Now(),
    }
    
    return uc.repo.Save(ctx, payment)
}
```

## 前缀规范建议

### 1. 前缀命名规范

- **长度**: 建议 3-5 个字符
- **格式**: 使用大写字母
- **语义**: 应该清晰表达业务含义

### 2. 常用前缀列表

| 前缀 | 业务类型 | 说明 | 示例 |
|------|---------|------|------|
| ORD  | 订单 | Order | ORD20241216000001 |
| PAY  | 支付单 | Payment | PAY20241216000001 |
| REF  | 退款单 | Refund | REF20241216000001 |
| INV  | 发票 | Invoice | INV20241216000001 |
| SHI  | 发货单 | Shipment | SHI20241216000001 |
| RET  | 退货单 | Return | RET20241216000001 |
| COU  | 优惠券 | Coupon | COU20241216000001 |
| ACT  | 活动 | Activity | ACT20241216000001 |
| TIC  | 票据 | Ticket | TIC20241216000001 |
| MEM  | 会员 | Member | MEM20241216000001 |

### 3. 前缀管理

建议在项目文档中维护一个前缀列表，避免冲突：

```markdown
# 业务前缀列表

- ORD: 订单
- PAY: 支付单
- REF: 退款单
- INV: 发票
- ...
```

## 完整示例：支付业务

```go
package biz

import (
    "context"
	"sre/internal/pkg/number"
    "github.com/go-kratos/kratos/v2/log"
)

// PaymentUsecase 支付业务用例
type PaymentUsecase struct {
    repo            PaymentRepo
    numberGenerator number.Generator
    log             *log.Helper
}

// NewPaymentUsecase 创建支付业务用例
func NewPaymentUsecase(
    repo PaymentRepo,
    numberGenerator number.Generator,
    logger log.Logger,
) *PaymentUsecase {
    return &PaymentUsecase{
        repo:            repo,
        numberGenerator: numberGenerator,
        log:             log.NewHelper(logger),
    }
}

// GeneratePaymentNo 生成支付单号
func (uc *PaymentUsecase) GeneratePaymentNo(ctx context.Context) (string, error) {
    return uc.numberGenerator.Generate(ctx, "PAY")
}

// CreatePayment 创建支付单
func (uc *PaymentUsecase) CreatePayment(ctx context.Context, orderID int64, amount int64) (*Payment, error) {
    // 生成支付单号
    paymentNo, err := uc.GeneratePaymentNo(ctx)
    if err != nil {
        uc.log.WithContext(ctx).Errorf("Failed to generate payment number: %v", err)
        return nil, err
    }
    
    payment := &Payment{
        PaymentNo: paymentNo,  // 例如: PAY20241216000001
        OrderID:   orderID,
        Amount:    amount,
        Status:    PaymentStatusPending,
        CreatedAt: time.Now(),
    }
    
    return uc.repo.Save(ctx, payment)
}
```

## 依赖注入配置

### Wire 配置

```go
// cmd/sre/wire.go
func wireApp(...) (*kratos.App, func(), error) {
    panic(wire.Build(
        // ...
        data.NewNumberGenerator,  // 生成器
        biz.NewOrderUsecase,           // 订单业务（使用 "ORD" 前缀）
        biz.NewPaymentUsecase,         // 支付业务（使用 "PAY" 前缀）
        biz.NewRefundUsecase,          // 退款业务（使用 "REF" 前缀）
        // ...
    ))
}
```

## 测试示例

```go
func TestPaymentUsecase_GeneratePaymentNo(t *testing.T) {
    // 创建 mock 生成器（需要实现 number.Generator 接口）
    // 可以使用 mock 库或者实现一个测试用的生成器
    generator := createMockNumberGenerator()
    
    uc := NewPaymentUsecase(
        mockPaymentRepo,
        generator,
        log.NewStdLogger(),
    )
    
    // 生成支付单号
    paymentNo, err := uc.GeneratePaymentNo(context.Background())
    assert.NoError(t, err)
    assert.True(t, strings.HasPrefix(paymentNo, "PAY"))
}
```

## 注意事项

1. **前缀唯一性**: 确保不同业务使用不同的前缀，避免冲突
2. **前缀长度**: 前缀不能超过10个字符
3. **前缀验证**: 生成器会自动验证前缀的有效性
4. **错误处理**: 生成失败时返回错误，调用方需要处理
5. **并发安全**: 生成器是并发安全的，可以在多个 goroutine 中使用

## 迁移说明

如果之前使用的是固定前缀 "ORD" 的订单号生成器，现在可以：

1. **继续使用订单号生成**: `OrderUsecase.GenerateOrderNo()` 仍然使用 "ORD" 前缀
2. **扩展其他业务**: 直接使用 `generator.Generate(ctx, "PAY")` 等生成其他业务的编号
3. **无需修改现有代码**: 订单号生成的接口保持不变，向后兼容

