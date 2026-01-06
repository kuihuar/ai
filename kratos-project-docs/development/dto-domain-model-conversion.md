# DTO 与 Domain Model 转换流程

## 概述

在 Kratos 分层架构中，数据在不同层之间传递时需要转换：
- **API 层**：使用 Protobuf 生成的 DTO（Data Transfer Object）
- **Biz 层**：使用 Domain Model（领域模型）
- **Data 层**：使用 Ent Entity（数据实体）

本文档详细说明这些转换在何时、何地、如何进行。

## 数据流向

```
请求流程：
Client → API (v1.OrderInfo) → Service → Biz (biz.Order) → Data (ent.Order) → Database

响应流程：
Database → Data (ent.Order) → Biz (biz.Order) → Service → API (v1.OrderInfo) → Client
```

## 转换位置：Service 层

**关键原则**：所有 DTO 与 Domain Model 的转换都在 **Service 层** 进行。

### 为什么在 Service 层？

1. **职责清晰**：Service 层负责协议转换，不包含业务逻辑
2. **解耦**：Biz 层不依赖 API 层的具体类型
3. **灵活性**：可以支持多种协议（gRPC、HTTP、GraphQL）而不影响业务逻辑

## 转换流程详解

### 1. API DTO → Domain Model（请求转换）

**位置**：`internal/service/order.go` 的各个方法中

**方式**：直接提取 DTO 字段，传递给 Biz 层方法

#### 示例 1：CreateOrder

```go
// internal/service/order.go
func (s *OrderService) CreateOrder(ctx context.Context, in *v1.CreateOrderRequest) (*v1.CreateOrderReply, error) {
    // 直接从 DTO 提取字段，传递给 Biz 层
    // 注意：Biz 层方法接收的是基本类型，不是 DTO
    order, err := s.uc.CreateOrder(ctx, 
        in.UserId,      // int64
        in.Amount,      // int64
        in.Currency,    // string
        in.Description, // string
    )
    // ...
}
```

**DTO 结构**（`api/order/v1/order.proto`）：
```protobuf
message CreateOrderRequest {
  int64 user_id = 1;
  int64 amount = 2;
  string currency = 3;
  string description = 4;
}
```

**Biz 层方法签名**（`internal/biz/order.go`）：
```go
func (uc *OrderUsecase) CreateOrder(
    ctx context.Context, 
    userID int64, 
    amount int64, 
    currency, description string,
) (*Order, error)
```

#### 示例 2：UpdateOrder（需要处理可选字段）

```go
// internal/service/order.go
func (s *OrderService) UpdateOrder(ctx context.Context, in *v1.UpdateOrderRequest) (*v1.UpdateOrderReply, error) {
    // 处理可选字段：将枚举转换为 int32
    var status *int32
    var description *string

    if in.Status != v1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
        statusVal := int32(in.Status)
        status = &statusVal
    }
    if in.Description != "" {
        description = &in.Description
    }

    // 传递给 Biz 层
    order, err := s.uc.UpdateOrder(ctx, in.Id, status, description)
    // ...
}
```

**关键点**：
- ✅ 直接提取 DTO 字段，不创建中间对象
- ✅ 处理可选字段（使用指针类型）
- ✅ 处理枚举类型转换
- ✅ 参数校验（如分页参数）

### 2. Domain Model → API DTO（响应转换）

**位置**：`internal/service/order.go` 的 `toOrderInfo` 方法

**方式**：使用转换方法将 Domain Model 转换为 DTO

#### 转换方法实现

```go
// internal/service/order.go
// toOrderInfo converts biz.Order to v1.OrderInfo.
func (s *OrderService) toOrderInfo(order *biz.Order) *v1.OrderInfo {
    orderInfo := &v1.OrderInfo{
        Id:          order.ID,
        UserId:      order.UserID,
        OrderNo:     order.OrderNo,
        Status:      v1.OrderStatus(order.Status),  // int32 → enum
        Amount:      order.Amount,
        Currency:    order.Currency,
        Description: order.Description,
        CreatedAt:   order.CreatedAt.Unix(),        // time.Time → int64
        UpdatedAt:   order.UpdatedAt.Unix(),        // time.Time → int64
    }

    // 处理可选字段
    if order.PaidAt != nil {
        orderInfo.PaidAt = order.PaidAt.Unix()
    }
    if order.CancelledAt != nil {
        orderInfo.CancelledAt = order.CancelledAt.Unix()
    }

    return orderInfo
}
```

#### 使用转换方法

```go
// CreateOrder 返回响应
func (s *OrderService) CreateOrder(ctx context.Context, in *v1.CreateOrderRequest) (*v1.CreateOrderReply, error) {
    order, err := s.uc.CreateOrder(ctx, in.UserId, in.Amount, in.Currency, in.Description)
    if err != nil {
        return nil, err
    }

    return &v1.CreateOrderReply{
        Order: s.toOrderInfo(order),  // Domain Model → DTO
    }, nil
}
```

#### 批量转换（ListOrders）

```go
func (s *OrderService) ListOrders(ctx context.Context, in *v1.ListOrdersRequest) (*v1.ListOrdersReply, error) {
    orders, total, err := s.uc.ListOrders(ctx, page, pageSize, userID, status, in.Keyword)
    if err != nil {
        return nil, err
    }

    // 批量转换
    orderInfos := make([]*v1.OrderInfo, 0, len(orders))
    for _, order := range orders {
        orderInfos = append(orderInfos, s.toOrderInfo(order))
    }

    return &v1.ListOrdersReply{
        Orders:   orderInfos,
        Total:    int32(total),
        Page:     int32(page),
        PageSize: int32(pageSize),
    }, nil
}
```

## 类型映射表

### 基本类型映射

| Biz 层类型 | API 层类型 | 转换方式 |
|-----------|----------|---------|
| `int64` | `int64` | 直接赋值 |
| `int32` | `int32` | 直接赋值 |
| `string` | `string` | 直接赋值 |
| `time.Time` | `int64` | `.Unix()` 转换为时间戳 |
| `*time.Time` | `int64` | 先判断 nil，再 `.Unix()` |
| `int32` (状态) | `enum` | `v1.OrderStatus(status)` |

### 复杂类型映射

| Biz 层类型 | API 层类型 | 转换方式 |
|-----------|----------|---------|
| `*int32` (可选) | `enum` | 判断是否为默认值，转换为指针 |
| `[]*Order` | `[]*OrderInfo` | 循环调用转换方法 |

## 完整转换示例

### 请求流程（CreateOrder）

```go
// 1. API 层接收请求（自动由 Kratos 处理）
// Client 发送 HTTP POST /api/v1/orders
// Kratos 自动将 JSON 转换为 v1.CreateOrderRequest

// 2. Service 层：DTO → Domain Model
func (s *OrderService) CreateOrder(ctx context.Context, in *v1.CreateOrderRequest) (*v1.CreateOrderReply, error) {
    // 提取 DTO 字段，直接传递给 Biz 层
    order, err := s.uc.CreateOrder(ctx, 
        in.UserId,      // DTO 字段
        in.Amount,      // DTO 字段
        in.Currency,    // DTO 字段
        in.Description, // DTO 字段
    )
    // ...
}

// 3. Biz 层：使用 Domain Model
func (uc *OrderUsecase) CreateOrder(ctx context.Context, userID int64, amount int64, currency, description string) (*Order, error) {
    // 创建 Domain Model
    order := &Order{
        UserID:      userID,
        OrderNo:     uc.generateOrderNo(),
        Status:      1,
        Amount:      amount,
        Currency:    currency,
        Description: description,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    // 调用 Data 层保存
    return uc.repo.Save(ctx, order)
}

// 4. Data 层：Domain Model → Ent Entity（在 Repository 中转换）
func (r *orderRepo) Save(ctx context.Context, order *biz.Order) (*biz.Order, error) {
    // 转换为 Ent Entity
    entOrder, err := r.client.Order.Create().
        SetUserID(order.UserID).
        SetOrderNo(order.OrderNo).
        SetStatus(order.Status).
        SetAmount(order.Amount).
        SetCurrency(order.Currency).
        SetDescription(order.Description).
        SetCreatedAt(order.CreatedAt.Unix()).
        SetUpdatedAt(order.UpdatedAt.Unix()).
        Save(ctx)
    // ...
}
```

### 响应流程（GetOrder）

```go
// 1. Data 层：Ent Entity → Domain Model
func (r *orderRepo) FindByID(ctx context.Context, id int64) (*biz.Order, error) {
    entOrder, err := r.client.Order.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 转换为 Domain Model
    order := &biz.Order{
        ID:          entOrder.ID,
        UserID:      entOrder.UserID,
        OrderNo:     entOrder.OrderNo,
        Status:      entOrder.Status,
        Amount:      entOrder.Amount,
        Currency:    entOrder.Currency,
        Description: entOrder.Description,
        CreatedAt:   time.Unix(entOrder.CreatedAt, 0),
        UpdatedAt:   time.Unix(entOrder.UpdatedAt, 0),
    }
    // ...
    return order, nil
}

// 2. Biz 层：直接返回 Domain Model
func (uc *OrderUsecase) GetOrder(ctx context.Context, id int64) (*Order, error) {
    return uc.repo.FindByID(ctx, id)
}

// 3. Service 层：Domain Model → DTO
func (s *OrderService) GetOrder(ctx context.Context, in *v1.GetOrderRequest) (*v1.GetOrderReply, error) {
    order, err := s.uc.GetOrder(ctx, in.Id)
    if err != nil {
        return nil, err
    }

    return &v1.GetOrderReply{
        Order: s.toOrderInfo(order),  // 转换方法
    }, nil
}

// 4. API 层：自动序列化为 JSON（由 Kratos 处理）
// 返回给 Client
```

## 最佳实践

### 1. 转换方法命名

使用统一的命名规范：
- `toOrderInfo` - Domain Model → DTO
- `fromOrderRequest` - DTO → Domain Model（如果需要）

### 2. 处理可选字段

```go
// ✅ 正确：使用指针类型
var status *int32
if in.Status != v1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
    statusVal := int32(in.Status)
    status = &statusVal
}

// ❌ 错误：直接使用零值
status := int32(in.Status)  // 无法区分"未设置"和"设置为0"
```

### 3. 处理时间转换

```go
// ✅ 正确：处理 nil 指针
if order.PaidAt != nil {
    orderInfo.PaidAt = order.PaidAt.Unix()
}

// ❌ 错误：直接调用会 panic
orderInfo.PaidAt = order.PaidAt.Unix()  // 如果 PaidAt 为 nil 会 panic
```

### 4. 批量转换优化

```go
// ✅ 正确：预分配容量
orderInfos := make([]*v1.OrderInfo, 0, len(orders))
for _, order := range orders {
    orderInfos = append(orderInfos, s.toOrderInfo(order))
}

// ❌ 错误：未预分配，会多次扩容
var orderInfos []*v1.OrderInfo
for _, order := range orders {
    orderInfos = append(orderInfos, s.toOrderInfo(order))
}
```

### 5. 错误处理

```go
// ✅ 正确：错误直接传播，不转换
order, err := s.uc.CreateOrder(ctx, ...)
if err != nil {
    return nil, err  // Kratos 会自动处理错误转换
}
```

## 总结

### 转换位置总结

| 转换类型 | 位置 | 方法 |
|---------|------|------|
| **DTO → Domain Model** | Service 层 | 直接提取字段，传递给 Biz 层 |
| **Domain Model → DTO** | Service 层 | `toOrderInfo()` 等转换方法 |
| **Domain Model → Ent Entity** | Data 层 Repository | 在 Repository 方法中转换 |
| **Ent Entity → Domain Model** | Data 层 Repository | 在 Repository 方法中转换 |

### 关键原则

1. ✅ **Service 层负责协议转换**：所有 DTO 与 Domain Model 的转换都在 Service 层
2. ✅ **Biz 层不依赖 API 类型**：Biz 层只使用 Domain Model 和基本类型
3. ✅ **Data 层负责数据转换**：Domain Model 与 Ent Entity 的转换在 Data 层
4. ✅ **保持转换方法集中**：使用 `toXxxInfo` 等方法集中管理转换逻辑

### 转换流程图

```
┌─────────────────────────────────────────────────────────────┐
│                     请求流程                                 │
├─────────────────────────────────────────────────────────────┤
│ Client                                                    │
│   ↓                                                         │
│ API Layer (v1.CreateOrderRequest)                         │
│   ↓                                                         │
│ Service Layer                                             │
│   • 提取 DTO 字段                                          │
│   • 传递给 Biz 层（基本类型）                               │
│   ↓                                                         │
│ Biz Layer (biz.Order)                                     │
│   • 创建 Domain Model                                      │
│   • 业务逻辑处理                                           │
│   ↓                                                         │
│ Data Layer (ent.Order)                                    │
│   • Domain Model → Ent Entity                              │
│   • 保存到数据库                                           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     响应流程                                 │
├─────────────────────────────────────────────────────────────┤
│ Database                                                  │
│   ↓                                                         │
│ Data Layer (ent.Order)                                    │
│   • Ent Entity → Domain Model                             │
│   ↓                                                         │
│ Biz Layer (biz.Order)                                     │
│   • 返回 Domain Model                                      │
│   ↓                                                         │
│ Service Layer                                             │
│   • Domain Model → DTO (toOrderInfo)                      │
│   ↓                                                         │
│ API Layer (v1.OrderInfo)                                  │
│   ↓                                                         │
│ Client                                                    │
└─────────────────────────────────────────────────────────────┘
```

