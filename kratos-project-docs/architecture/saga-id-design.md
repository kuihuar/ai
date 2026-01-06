# Saga ID 设计说明

## 当前实现问题

### 问题 1: Saga ID 和订单号不一致

**当前实现**:
```go
// OrderCreateSaga.Run() - 生成 Saga ID
orderNo := fmt.Sprintf("SAGA-%d-%d", time.Now().Unix(), userID)
sagaCtx := &SagaContext{
    ID: orderNo,  // Saga ID = "SAGA-1234567890-1"
    ...
}

// OrderUsecase.CreateOrder() - 生成订单号
orderNo := uc.generateOrderNo()  // 格式: "ORD1234567890123456"
order := &Order{
    OrderNo: orderNo,  // 订单号 = "ORD1234567890123456"
    ...
}

// Service 层返回
return &v1.CreateOrderSagaReply{
    SagaId: order.OrderNo,  // 返回的是订单号，不是 Saga ID！
    ...
}
```

**问题**:
- Saga ID (`SAGA-1234567890-1`) 和订单号 (`ORD1234567890123456`) **不一致**
- 返回的 `saga_id` 实际上是订单号，不是真正的 Saga ID
- 无法通过返回的 `saga_id` 查询 `saga_instances` 表

### 问题 2: Saga 失败时无法获取订单号

**当前实现**:
```go
if err != nil {
    return &v1.CreateOrderSagaReply{
        SagaId: "",  // Saga 失败时无法获取 ID
        ...
    }, err
}
```

**问题**:
- 如果 Saga 在 Step1（创建订单）之前失败，没有订单号
- 如果 Saga 在 Step1 之后失败，订单可能已被创建，但无法获取订单号来查询 Saga 状态

---

## 使用订单号作为 Saga ID 的含义

### 方案 A: Saga ID = 订单号（推荐）

**设计思路**:
- 在创建 Saga 之前，先生成订单号
- 使用订单号作为 Saga ID
- 这样 Saga ID 和订单号保持一致

**优点**:
1. ✅ **业务语义清晰**: Saga ID 就是订单号，用户可以直接用订单号查询 Saga 状态
2. ✅ **便于追踪**: 通过订单号可以同时查询订单和 Saga 状态
3. ✅ **简化查询**: 不需要维护两套 ID 映射关系
4. ✅ **用户体验好**: 用户只需要记住订单号，就可以追踪整个流程

**缺点**:
- ⚠️ 订单号必须在 Saga 开始前生成（当前实现是在 Step1 中生成）

### 方案 B: Saga ID ≠ 订单号（当前实现）

**设计思路**:
- Saga ID 和订单号独立生成
- 通过 `saga_instances.metadata` 关联

**优点**:
- ✅ Saga ID 可以在订单创建前就确定
- ✅ 支持更复杂的 Saga 类型（不一定是订单相关）

**缺点**:
- ❌ 需要额外的查询来关联 Saga 和订单
- ❌ 用户需要记住两个 ID
- ❌ 当前实现中返回的 `saga_id` 实际上是订单号，造成混淆

---

## 推荐方案：统一使用订单号作为 Saga ID

### 实现方式

#### 1. 修改 OrderCreateSaga.Run()

```go
func (s *OrderCreateSaga) Run(ctx context.Context, userID int64, amount int64, currency, description string, items []*OrderItem) (*Order, error) {
    // 先生成订单号（在 Saga 开始前）
    orderNo := s.uc.GenerateOrderNo()  // 需要将 generateOrderNo 改为公开方法
    
    sagaCtx := &SagaContext{
        ID:        orderNo,  // 使用订单号作为 Saga ID
        Type:      "order.create",
        Metadata:  map[string]string{"user_id": fmt.Sprintf("%d", userID)},
        StartedAt: time.Now(),
    }
    
    // ... 后续逻辑
}
```

#### 2. 修改 orderCreateStep.Execute()

```go
func (s *orderCreateStep) Execute(ctx context.Context, sagaCtx *SagaContext) error {
    // 从 SagaContext 获取订单号（而不是生成新的）
    orderNo := sagaCtx.ID  // Saga ID 就是订单号
    
    // 创建订单时使用这个订单号
    order, err := s.uc.CreateOrderWithOrderNo(ctx, orderNo, ...)
    // ...
}
```

#### 3. Service 层返回

```go
return &v1.CreateOrderSagaReply{
    Order:       s.toOrderInfo(order),
    SagaId:      order.OrderNo,  // 现在 Saga ID 和订单号一致了
    Compensated: false,
}, nil
```

---

## 使用订单号作为 Saga ID 的含义

### 1. **业务语义**

**含义**: Saga ID = 订单号，表示这个 Saga 实例就是为创建这个订单而存在的。

**好处**:
- 用户只需要记住订单号
- 通过订单号可以直接查询：
  - 订单信息 (`orders` 表)
  - Saga 执行状态 (`saga_instances` 表)
  - Saga 步骤详情 (`saga_steps` 表)

### 2. **追踪和审计**

**含义**: 订单号和 Saga ID 一一对应，便于追踪整个订单创建流程。

**查询示例**:
```sql
-- 通过订单号查询 Saga 状态
SELECT * FROM saga_instances WHERE saga_id = 'ORD1234567890123456';

-- 通过订单号查询 Saga 步骤
SELECT * FROM saga_steps WHERE saga_id = 'ORD1234567890123456';

-- 通过订单号查询订单
SELECT * FROM orders WHERE order_no = 'ORD1234567890123456';
```

### 3. **API 响应语义**

**当前问题**:
```json
{
  "order": {
    "order_no": "ORD1234567890123456"
  },
  "saga_id": "ORD1234567890123456"  // 实际上是订单号，不是真正的 Saga ID
}
```

**改进后**:
```json
{
  "order": {
    "order_no": "ORD1234567890123456"
  },
  "saga_id": "ORD1234567890123456"  // 明确：Saga ID = 订单号
}
```

**含义**: 
- `saga_id` 字段明确表示：这个 Saga 实例的 ID 就是这个订单号
- 用户可以用这个 ID 查询 Saga 的完整执行历史

### 4. **日志和监控**

**含义**: 所有日志和监控指标都可以使用订单号作为唯一标识。

**好处**:
- 日志中 `saga_id` 和 `order_no` 一致，便于关联
- 监控指标可以直接用订单号分组
- 分布式追踪中，订单号可以作为 trace ID 的一部分

---

## 当前实现的问题总结

### 问题 1: ID 不一致

```
Saga ID (saga_instances.saga_id): "SAGA-1234567890-1"
订单号 (orders.order_no): "ORD1234567890123456"
返回的 saga_id: "ORD1234567890123456"  ❌ 无法查询 saga_instances
```

### 问题 2: 失败时无法获取

```go
if err != nil {
    SagaId: "",  // ❌ 无法获取，用户无法追踪失败的 Saga
}
```

### 问题 3: 语义混淆

- 返回的 `saga_id` 实际上是订单号
- 但真正的 Saga ID 存储在数据库中，无法通过返回的 `saga_id` 查询

---

## 修复建议

### 方案 1: 统一使用订单号（推荐）

**步骤**:
1. 在 `OrderCreateSaga.Run()` 中，先调用 `GenerateOrderNo()` 生成订单号
2. 使用订单号作为 Saga ID
3. 在 `orderCreateStep.Execute()` 中使用这个订单号创建订单
4. Service 层返回时，`saga_id` 就是订单号

**优点**: 
- ✅ 语义清晰
- ✅ 便于查询
- ✅ 用户体验好

### 方案 2: 返回真正的 Saga ID

**步骤**:
1. 保持当前实现（Saga ID 和订单号独立）
2. 在 Service 层，从 `sagaCtx.ID` 获取真正的 Saga ID
3. 返回真正的 Saga ID

**优点**:
- ✅ 保持当前架构
- ✅ 支持更复杂的 Saga 类型

**缺点**:
- ❌ 需要额外的查询来关联 Saga 和订单
- ❌ 用户需要记住两个 ID

---

## 总结

**使用订单号作为 Saga ID 的含义**:

1. **业务语义**: Saga 实例和订单一一对应
2. **追踪能力**: 通过订单号可以追踪整个 Saga 流程
3. **用户体验**: 用户只需要记住订单号
4. **查询便利**: 不需要额外的 ID 映射

**当前实现的问题**:
- Saga ID 和订单号不一致
- 返回的 `saga_id` 无法查询 `saga_instances` 表
- 失败时无法获取 Saga ID

**推荐方案**:
- 统一使用订单号作为 Saga ID
- 在 Saga 开始前生成订单号
- 确保 Saga ID 和订单号始终一致

