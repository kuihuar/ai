# Saga API 使用指南

本文档说明如何使用 Saga 模式的订单创建 API。

## API 概述

`CreateOrderSaga` 是一个使用 Saga 模式创建订单的 API，支持跨服务事务编排和自动补偿。

### 与普通 CreateOrder 的区别

| 特性 | CreateOrder | CreateOrderSaga |
|------|-------------|-----------------|
| **事务模式** | 本地事务 + Outbox | Saga 编排（跨服务） |
| **补偿机制** | ❌ 无 | ✅ 自动补偿 |
| **适用场景** | 单服务内操作 | 跨多个服务的复杂流程 |
| **性能** | 更快 | 稍慢（需要多次调用） |

## API 定义

### gRPC

```protobuf
rpc CreateOrderSaga (CreateOrderSagaRequest) returns (CreateOrderSagaReply)
```

### HTTP REST

```http
POST /api/v1/orders/saga
Content-Type: application/json

{
  "user_id": 1,
  "currency": "CNY",
  "description": "测试订单（Saga模式）",
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 5000
    }
  ]
}
```

## 请求参数

### CreateOrderSagaRequest

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `user_id` | int64 | ✅ | 用户ID |
| `currency` | string | ❌ | 货币类型（默认：CNY） |
| `description` | string | ❌ | 订单描述 |
| `items` | OrderItemRequest[] | ✅ | 订单项列表 |

### OrderItemRequest

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `product_id` | int64 | ✅ | 产品ID |
| `quantity` | int32 | ✅ | 数量 |
| `price` | int64 | ✅ | 单价（分） |

## 响应参数

### CreateOrderSagaReply

| 字段 | 类型 | 说明 |
|------|------|------|
| `order` | OrderInfo | 订单信息（如果 Saga 失败且已补偿，可能为 null） |
| `saga_id` | string | Saga 实例ID（成功时等于订单号） |
| `compensated` | bool | 是否执行了补偿（true 表示流程失败但已回滚） |

## 使用示例

### 1. 使用 curl 调用 HTTP API

```bash
curl -X POST http://localhost:8000/api/v1/orders/saga \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "测试订单（Saga模式）",
    "items": [
      {
        "product_id": 1,
        "quantity": 2,
        "price": 5000
      }
    ]
  }'
```

### 2. 使用 grpcurl 调用 gRPC API

```bash
grpcurl -plaintext -d '{
  "user_id": 1,
  "currency": "CNY",
  "description": "测试订单（Saga模式）",
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 5000
    }
  ]
}' localhost:8989 order.v1.Order/CreateOrderSaga
```

### 3. Go 客户端示例

```go
package main

import (
    "context"
    "sre/api/order/v1"
    "google.golang.org/grpc"
)

func main() {
    conn, _ := grpc.Dial("localhost:8989", grpc.WithInsecure())
    defer conn.Close()
    
    client := v1.NewOrderClient(conn)
    
    req := &v1.CreateOrderSagaRequest{
        UserId:      1,
        Currency:    "CNY",
        Description: "测试订单（Saga模式）",
        Items: []*v1.OrderItemRequest{
            {
                ProductId: 1,
                Quantity:  2,
                Price:     5000,
            },
        },
    }
    
    reply, err := client.CreateOrderSaga(context.Background(), req)
    if err != nil {
        // 处理错误
        // 如果 compensated = true，说明已执行补偿
        return
    }
    
    // 使用 reply.Order 和 reply.SagaId
}
```

## Saga 执行流程

### 成功流程

```
1. 调用 CreateOrderSaga API
   ↓
2. Step1: 创建订单（本地服务，带 Outbox 事件）
   ↓
3. Step2: 预留库存（占位步骤，当前仅打日志）
   ↓
4. Step3: 冻结支付（占位步骤，当前仅打日志）
   ↓
5. 返回成功响应（order + saga_id）
```

### 失败 + 补偿流程

```
1. 调用 CreateOrderSaga API
   ↓
2. Step1: 创建订单 ✅ 成功
   ↓
3. Step2: 预留库存 ❌ 失败
   ↓
4. 触发补偿：
   - Compensate Step2（释放库存，占位）
   - Compensate Step1（取消订单，占位）
   ↓
5. 返回错误响应（compensated = true）
```

## 日志示例

### 成功场景

```
INFO msg="CreateOrderSaga: userID=1, currency=CNY, itemsCount=1"
INFO msg="Saga started: id=SAGA-1234567890-1, type=order.create"
INFO msg="Saga Step Execute: create-order, saga_id=SAGA-1234567890-1, user_id=1"
INFO msg="Saga Step Execute: reserve-inventory (no-op), saga_id=SAGA-1234567890-1"
INFO msg="Saga Step Execute: freeze-payment (no-op), saga_id=SAGA-1234567890-1"
INFO msg="Saga completed successfully: id=SAGA-1234567890-1"
```

### 失败 + 补偿场景

```
INFO msg="CreateOrderSaga: userID=1, currency=CNY, itemsCount=1"
INFO msg="Saga started: id=SAGA-1234567890-1, type=order.create"
INFO msg="Saga Step Execute: create-order, saga_id=SAGA-1234567890-1, user_id=1"
INFO msg="Saga Step Execute: reserve-inventory (no-op), saga_id=SAGA-1234567890-1"
ERROR msg="Saga step failed: name=freeze-payment, error=..."
WARN msg="Saga Step Compensate: reserve-inventory (no-op), saga_id=SAGA-1234567890-1"
WARN msg="Saga Step Compensate: create-order, saga_id=SAGA-1234567890-1"
WARN msg="Saga completed with compensation: id=SAGA-1234567890-1"
ERROR msg="CreateOrderSaga failed: ..."
```

## 错误处理

### 常见错误

1. **ORDER_ITEMS_REQUIRED**
   - 原因：未提供订单项
   - 处理：确保 `items` 字段不为空

2. **Saga 执行失败**
   - 原因：某个步骤执行失败
   - 处理：
     - 检查 `compensated` 字段
     - 如果 `compensated = true`，说明已自动回滚，可以重试
     - 如果 `compensated = false`，说明补偿失败，需要人工介入

### 补偿状态说明

- **compensated = false**：Saga 成功，或失败但未执行补偿
- **compensated = true**：Saga 失败，但已成功执行补偿（所有步骤已回滚）

## 当前实现状态

### 已实现

✅ **Saga 编排框架**：通用的 Saga Orchestrator  
✅ **订单创建 Saga**：订单创建流程的 Saga 实现  
✅ **API 接口**：CreateOrderSaga RPC/HTTP 接口  
✅ **自动补偿**：失败时自动执行补偿流程  

### 待实现（占位步骤）

⏳ **库存服务集成**：Step2 目前是占位步骤，需要接入真实的库存服务  
⏳ **支付服务集成**：Step3 目前是占位步骤，需要接入真实的支付服务  
⏳ **补偿逻辑完善**：Step1 的补偿目前仅打日志，需要实现真正的订单取消逻辑  
⏳ **Saga 状态持久化**：当前 Saga 状态仅在内存中，需要持久化到数据库  

## 扩展建议

### 1. 集成真实的库存服务

在 `order_saga.go` 中，将 `noOpStep{name: "reserve-inventory"}` 替换为：

```go
&inventoryReserveStep{
    inventoryClient: inventoryClient, // gRPC 客户端
    logger: s.log,
}
```

### 2. 集成真实的支付服务

类似地，将 `noOpStep{name: "freeze-payment"}` 替换为：

```go
&paymentFreezeStep{
    paymentClient: paymentClient, // gRPC 客户端
    logger: s.log,
}
```

### 3. 实现真正的补偿逻辑

在 `orderCreateStep.Compensate` 中：

```go
func (s *orderCreateStep) Compensate(ctx context.Context, sagaCtx *SagaContext) error {
    if idStr, ok := sagaCtx.Metadata["order_id"]; ok {
        orderID, _ := strconv.ParseInt(idStr, 10, 64)
        // 调用 CancelOrder
        return s.uc.CancelOrder(ctx, orderID, "Saga compensation")
    }
    return nil
}
```

## 相关文档

- [分布式事件一致性架构设计](../architecture/distributed-event-consistency.md)
- [Outbox 模式端到端测试指南](./outbox-testing-guide.md)

