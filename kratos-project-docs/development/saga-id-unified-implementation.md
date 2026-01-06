# Saga ID 统一实现说明

## 实现概述

已按照推荐方案实现：**统一使用订单号作为 Saga ID**，确保 Saga ID 和订单号始终一致。

## 实现内容

### 1. 公开订单号生成方法

**文件**: `internal/biz/order.go`

```go
// GenerateOrderNo generates a unique order number.
// 公开方法，允许外部（如 Saga）在创建订单前生成订单号。
func (uc *OrderUsecase) GenerateOrderNo() string {
    return fmt.Sprintf("ORD%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}
```

**变更**:
- ✅ 将 `generateOrderNo()` 改为公开方法 `GenerateOrderNo()`
- ✅ 保留内部方法 `generateOrderNo()` 以保持向后兼容

### 2. 添加带订单号的创建方法

**文件**: `internal/biz/order.go`

```go
// CreateOrderWithOrderNo 使用指定的订单号创建订单。
// 用于 Saga 模式，确保 Saga ID 和订单号一致。
func (uc *OrderUsecase) CreateOrderWithOrderNo(ctx context.Context, orderNo string, userID int64, amount int64, currency, description string, items []*OrderItem) (*Order, error)
```

**功能**:
- ✅ 接受订单号作为参数
- ✅ 使用指定的订单号创建订单
- ✅ 其他逻辑与 `CreateOrder()` 相同

### 3. 修改 Saga 使用订单号作为 ID

**文件**: `internal/biz/order_saga.go`

**变更前**:
```go
orderNo := fmt.Sprintf("SAGA-%d-%d", time.Now().Unix(), userID)
sagaCtx := &SagaContext{
    ID: orderNo,  // Saga ID = "SAGA-xxx"
    ...
}
```

**变更后**:
```go
// 在 Saga 开始前生成订单号，使用订单号作为 Saga ID
orderNo := s.uc.GenerateOrderNo()

sagaCtx := &SagaContext{
    ID: orderNo,  // Saga ID = 订单号（如 "ORD1234567890123456"）
    ...
}
```

### 4. 修改订单创建步骤使用 Saga ID

**文件**: `internal/biz/order_saga.go`

**变更**:
- ✅ `orderCreateStep` 结构体添加 `orderNo` 字段
- ✅ `Execute()` 方法使用 `CreateOrderWithOrderNo()` 创建订单
- ✅ 验证订单号与 Saga ID 一致

```go
type orderCreateStep struct {
    uc          *OrderUsecase
    orderNo     string // 订单号（Saga ID）
    ...
}

func (s *orderCreateStep) Execute(ctx context.Context, sagaCtx *SagaContext) error {
    // 使用 CreateOrderWithOrderNo 确保订单号与 Saga ID 一致
    order, err := s.uc.CreateOrderWithOrderNo(ctx, s.orderNo, ...)
    
    // 验证订单号是否与 Saga ID 一致
    if order.OrderNo != sagaCtx.ID {
        return fmt.Errorf("order_no mismatch: expected=%s, got=%s", sagaCtx.ID, order.OrderNo)
    }
    ...
}
```

### 5. 修改 Service 层返回正确的 Saga ID

**文件**: `internal/service/order.go`

**变更**:
- ✅ 成功时：返回订单号作为 `saga_id`
- ✅ 失败时：如果订单已创建，返回订单号；否则返回空字符串

```go
// Saga 成功
return &v1.CreateOrderSagaReply{
    Order:       s.toOrderInfo(order),
    SagaId:      order.OrderNo, // Saga ID = 订单号
    Compensated: false,
}, nil

// Saga 失败
if order != nil {
    sagaID = order.OrderNo  // 使用订单号
} else {
    sagaID = ""  // 订单未创建，无法获取 Saga ID
}
```

## 实现效果

### 1. Saga ID 和订单号一致

**之前**:
```
Saga ID: "SAGA-1234567890-1"
订单号: "ORD1234567890123456"
返回的 saga_id: "ORD1234567890123456"  ❌ 无法查询 saga_instances
```

**现在**:
```
Saga ID: "ORD1234567890123456"
订单号: "ORD1234567890123456"
返回的 saga_id: "ORD1234567890123456"  ✅ 可以查询 saga_instances
```

### 2. 可以通过订单号查询 Saga 状态

```sql
-- 通过订单号查询 Saga 实例
SELECT * FROM saga_instances WHERE saga_id = 'ORD1234567890123456';

-- 通过订单号查询 Saga 步骤
SELECT * FROM saga_steps WHERE saga_id = 'ORD1234567890123456';

-- 通过订单号查询订单
SELECT * FROM orders WHERE order_no = 'ORD1234567890123456';
```

### 3. 用户体验提升

- ✅ 用户只需要记住订单号
- ✅ 通过订单号可以追踪整个 Saga 流程
- ✅ API 返回的 `saga_id` 可以直接用于查询

## 数据流

### 成功场景

```
1. Service 层调用 OrderCreateSaga.Run()
   ↓
2. OrderCreateSaga.Run() 生成订单号: "ORD1234567890123456"
   ↓
3. 创建 SagaContext，ID = "ORD1234567890123456"
   ↓
4. 持久化到 saga_instances 表，saga_id = "ORD1234567890123456"
   ↓
5. orderCreateStep.Execute() 使用订单号创建订单
   ↓
6. 订单创建成功，order_no = "ORD1234567890123456"
   ↓
7. 验证：order.OrderNo == sagaCtx.ID ✅
   ↓
8. Service 层返回：saga_id = "ORD1234567890123456"
```

### 失败场景

```
1. Service 层调用 OrderCreateSaga.Run()
   ↓
2. OrderCreateSaga.Run() 生成订单号: "ORD1234567890123456"
   ↓
3. 创建 SagaContext，ID = "ORD1234567890123456"
   ↓
4. 持久化到 saga_instances 表，saga_id = "ORD1234567890123456"
   ↓
5. orderCreateStep.Execute() 使用订单号创建订单
   ↓
6. 订单创建成功，order_no = "ORD1234567890123456"
   ↓
7. 后续步骤失败，执行补偿
   ↓
8. Service 层返回：saga_id = "ORD1234567890123456"（即使失败也能查询）
```

## 验证方法

### 1. 验证 Saga ID 和订单号一致

```sql
SELECT 
    si.saga_id,
    o.order_no,
    CASE 
        WHEN si.saga_id = o.order_no THEN '✅ 一致'
        ELSE '❌ 不一致'
    END AS validation
FROM saga_instances si
LEFT JOIN orders o ON si.saga_id = o.order_no
WHERE si.saga_type = 'order.create'
ORDER BY si.created_at DESC
LIMIT 10;
```

### 2. 通过订单号查询 Saga 状态

```bash
# 调用 API 获取订单号和 Saga ID
curl -X POST http://localhost:8000/api/v1/orders/saga \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "测试订单",
    "items": [{"product_id": 1, "quantity": 1, "price": 1000}]
  }'

# 响应中的 saga_id 就是订单号
# 使用这个订单号查询 Saga 状态
```

### 3. 验证 API 响应

**成功响应**:
```json
{
  "order": {
    "order_no": "ORD1234567890123456"
  },
  "saga_id": "ORD1234567890123456",  // ✅ 与订单号一致
  "compensated": false
}
```

**失败响应**:
```json
{
  "order": {
    "order_no": "ORD1234567890123456"
  },
  "saga_id": "ORD1234567890123456",  // ✅ 即使失败也能返回
  "compensated": true
}
```

## 注意事项

### 1. 订单号生成时机

- ✅ 在 Saga 开始前生成订单号
- ✅ 确保订单号在创建订单前就确定
- ✅ 避免订单号和 Saga ID 不一致

### 2. 失败场景处理

- ✅ 如果订单已创建，返回订单号作为 Saga ID
- ✅ 如果订单未创建，返回空字符串（无法确定 Saga ID）
- ⚠️ 极端情况下，如果 Saga 在创建实例前失败，无法获取 Saga ID

### 3. 向后兼容

- ✅ `CreateOrder()` 方法保持不变
- ✅ `generateOrderNo()` 内部方法保留
- ✅ 不影响现有的订单创建流程

## 总结

通过统一使用订单号作为 Saga ID，实现了：

1. ✅ **语义清晰**: Saga ID = 订单号，一一对应
2. ✅ **查询便利**: 通过订单号可以查询所有相关信息
3. ✅ **用户体验**: 用户只需要记住订单号
4. ✅ **追踪能力**: 完整的 Saga 执行历史可以通过订单号查询

现在，用户可以通过订单号直接查询 Saga 状态，无需额外的 ID 映射！

