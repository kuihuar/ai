# 订单 API 列表

本文档列出了所有订单相关的 API 接口。

## API 概览

| 序号 | API 名称 | gRPC 方法 | HTTP 方法 | HTTP 路径 | 功能描述 |
|------|---------|-----------|-----------|-----------|----------|
| 1 | 创建订单 | `CreateOrder` | POST | `/api/v1/orders` | 创建新订单 |
| 2 | 获取订单 | `GetOrder` | GET | `/api/v1/orders/{id}` | 根据ID获取订单信息 |
| 3 | 更新订单 | `UpdateOrder` | PUT | `/api/v1/orders/{id}` | 更新订单信息（状态、描述） |
| 4 | 删除订单 | `DeleteOrder` | DELETE | `/api/v1/orders/{id}` | 删除订单 |
| 5 | 列出订单 | `ListOrders` | GET | `/api/v1/orders` | 分页列出订单，支持多条件筛选 |
| 6 | 取消订单 | `CancelOrder` | POST | `/api/v1/orders/{id}/cancel` | 取消订单 |

## 详细说明

### 1. 创建订单

**接口信息**
- **gRPC**: `order.v1.Order.CreateOrder`
- **HTTP**: `POST /api/v1/orders`

**请求参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| user_id | int64 | 是 | 用户ID（必须 > 0） |
| currency | string | 否 | 货币类型（默认 CNY） |
| description | string | 否 | 订单描述 |
| items | OrderItemRequest[] | 是 | 订单项列表（至少包含一个订单项） |

**OrderItemRequest 结构**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| product_id | int64 | 是 | 产品ID（必须 > 0，产品必须存在） |
| quantity | int32 | 是 | 商品数量（必须 > 0） |
| price | int64 | 否 | 单价（分），下单时的价格快照。如果未提供或 <= 0，则使用产品当前价格 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| order | OrderInfo | 订单信息对象 |

**业务规则**

1. **订单号自动生成**：系统自动生成唯一订单号（格式：`ORD{timestamp}{nanosecond}`）
2. **初始状态**：新订单默认状态为 `ORDER_STATUS_PENDING`（待支付）
3. **用户验证**：用户ID必须 > 0
4. **订单项验证**：
   - 订单项列表不能为空（至少包含一个订单项）
   - 每个订单项的产品ID必须 > 0，且产品必须存在
   - 每个订单项的数量必须 > 0
   - 每个订单项的价格必须 >= 0（如果未提供或 <= 0，则使用产品当前价格）
5. **金额自动计算**：
   - 订单总金额由系统根据订单项自动计算：`总金额 = Σ(订单项单价 × 数量)`
   - 不需要手动传入 `amount` 字段
6. **订单项创建**：创建订单成功后，系统会自动批量创建对应的订单项
7. **事务处理**：如果订单项创建失败，订单创建也会失败（保证数据一致性）

**错误码**

| 错误码 | 说明 |
|--------|------|
| ORDER_INVALID_USER_ID | 用户ID无效 |
| ORDER_INVALID_AMOUNT | 订单金额无效（计算后金额 <= 0） |
| ORDER_ITEMS_REQUIRED | 订单项列表为空（必须至少包含一个订单项） |
| PRODUCT_NOT_FOUND | 产品不存在（订单项中的 product_id 无效） |
| INVALID_ORDER_ITEM_QUANTITY | 订单项数量无效（必须 > 0） |
| INVALID_ORDER_ITEM_PRICE | 订单项价格无效（必须 >= 0） |

---

### 2. 获取订单

**接口信息**
- **gRPC**: `order.v1.Order.GetOrder`
- **HTTP**: `GET /api/v1/orders/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 订单ID |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| order | OrderInfo | 订单信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| ORDER_NOT_FOUND | 订单不存在 |

---

### 3. 更新订单

**接口信息**
- **gRPC**: `order.v1.Order.UpdateOrder`
- **HTTP**: `PUT /api/v1/orders/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 订单ID |
| status | OrderStatus | 否 | 请求体 | 订单状态（可选） |
| description | string | 否 | 请求体 | 订单描述（可选） |

**订单状态说明**

| 状态值 | 枚举 | 说明 |
|--------|------|------|
| 0 | ORDER_STATUS_UNSPECIFIED | 未指定 |
| 1 | ORDER_STATUS_PENDING | 待支付 |
| 2 | ORDER_STATUS_PAID | 已支付 |
| 3 | ORDER_STATUS_SHIPPED | 已发货 |
| 4 | ORDER_STATUS_COMPLETED | 已完成 |
| 5 | ORDER_STATUS_CANCELLED | 已取消 |
| 6 | ORDER_STATUS_REFUNDED | 已退款 |

**业务规则**

1. **状态变更**：
   - 状态变为 `ORDER_STATUS_PAID`（2）时，自动设置 `paid_at` 时间
   - 状态变为 `ORDER_STATUS_CANCELLED`（5）时，自动设置 `cancelled_at` 时间
2. **状态验证**：状态值必须在 1-6 之间

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| order | OrderInfo | 订单信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| ORDER_NOT_FOUND | 订单不存在 |
| ORDER_INVALID_STATUS | 订单状态无效 |

---

### 4. 删除订单

**接口信息**
- **gRPC**: `order.v1.Order.DeleteOrder`
- **HTTP**: `DELETE /api/v1/orders/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 订单ID |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | bool | 是否成功 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| ORDER_NOT_FOUND | 订单不存在 |

---

### 5. 列出订单

**接口信息**
- **gRPC**: `order.v1.Order.ListOrders`
- **HTTP**: `GET /api/v1/orders`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 默认值 | 说明 |
|--------|------|------|------|--------|------|
| page | int32 | 否 | 查询参数 | 1 | 页码（从1开始） |
| page_size | int32 | 否 | 查询参数 | 10 | 每页数量（最大100） |
| user_id | int64 | 否 | 查询参数 | - | 用户ID（可选，筛选特定用户的订单） |
| status | OrderStatus | 否 | 查询参数 | - | 订单状态（可选，筛选特定状态的订单） |
| keyword | string | 否 | 查询参数 | - | 搜索关键词（可选，搜索订单号或描述） |
| start_time | int64 | 否 | 查询参数 | - | 开始时间（Unix时间戳，筛选创建时间 >= start_time） |
| end_time | int64 | 否 | 查询参数 | - | 结束时间（Unix时间戳，筛选创建时间 <= end_time） |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| orders | OrderInfo[] | 订单列表 |
| total | int32 | 总数量 |
| page | int32 | 当前页码 |
| page_size | int32 | 每页数量 |

**示例**

```bash
# 基本查询
GET /api/v1/orders?page=1&page_size=10

# 按用户筛选
GET /api/v1/orders?user_id=1

# 按状态筛选
GET /api/v1/orders?status=2

# 时间范围筛选
GET /api/v1/orders?start_time=1704067200&end_time=1704153600

# 关键词搜索
GET /api/v1/orders?keyword=ORD123456

# 组合查询
GET /api/v1/orders?page=1&page_size=20&user_id=1&status=2&start_time=1704067200&end_time=1704153600
```

---

### 6. 取消订单

**接口信息**
- **gRPC**: `order.v1.Order.CancelOrder`
- **HTTP**: `POST /api/v1/orders/{id}/cancel`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 订单ID |
| reason | string | 否 | 请求体 | 取消原因（可选） |

**业务规则**

1. **状态检查**：
   - 已取消的订单（状态 5）不能再次取消
   - 已完成的订单（状态 4）不能取消
2. **自动设置**：
   - 订单状态自动更新为 `ORDER_STATUS_CANCELLED`（5）
   - 自动设置 `cancelled_at` 时间
   - 如果提供了取消原因，会追加到订单描述中

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| order | OrderInfo | 订单信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| ORDER_NOT_FOUND | 订单不存在 |
| ORDER_INVALID_STATUS | 订单状态无效（已取消或已完成） |

---

## 数据模型

### OrderInfo

订单信息对象，用于返回订单的基本信息。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 订单ID |
| user_id | int64 | 用户ID |
| order_no | string | 订单号（唯一） |
| status | OrderStatus | 订单状态 |
| amount | int64 | 订单金额（分） |
| currency | string | 货币类型（默认 CNY） |
| description | string | 订单描述 |
| created_at | int64 | 创建时间（Unix时间戳） |
| updated_at | int64 | 更新时间（Unix时间戳） |
| paid_at | int64 | 支付时间（Unix时间戳，0表示未支付） |
| cancelled_at | int64 | 取消时间（Unix时间戳，0表示未取消） |

### OrderStatus 枚举

| 值 | 枚举名称 | 说明 |
|---|---------|------|
| 0 | ORDER_STATUS_UNSPECIFIED | 未指定 |
| 1 | ORDER_STATUS_PENDING | 待支付 |
| 2 | ORDER_STATUS_PAID | 已支付 |
| 3 | ORDER_STATUS_SHIPPED | 已发货 |
| 4 | ORDER_STATUS_COMPLETED | 已完成 |
| 5 | ORDER_STATUS_CANCELLED | 已取消 |
| 6 | ORDER_STATUS_REFUNDED | 已退款 |

---

## 错误码

所有 API 的错误码定义在 `api/order/v1/error_reason.proto` 中。

常见错误码：

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| ORDER_NOT_FOUND | 404 | 订单不存在 |
| ORDER_ALREADY_EXISTS | 409 | 订单已存在 |
| ORDER_INVALID_STATUS | 400 | 订单状态无效 |
| ORDER_INVALID_AMOUNT | 400 | 订单金额无效 |
| ORDER_INVALID_USER_ID | 400 | 用户ID无效 |
| ORDER_SAVE_FAILED | 500 | 保存订单失败 |
| ORDER_UPDATE_FAILED | 500 | 更新订单失败 |
| ORDER_QUERY_FAILED | 500 | 查询订单失败 |
| ORDER_DELETE_FAILED | 500 | 删除订单失败 |
| ORDER_CANCEL_FAILED | 500 | 取消订单失败 |

---

## 使用示例

### HTTP 请求示例

#### 创建订单（单个产品）
```bash
curl -X POST http://localhost:8000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "购买 iPhone 15",
    "items": [
      {
        "product_id": 1,
        "quantity": 1,
        "price": 599900
      }
    ]
  }'
```

#### 创建订单（多个产品）
```bash
curl -X POST http://localhost:8000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "购买 iPhone 15 和 AirPods",
    "items": [
      {
        "product_id": 1,
        "quantity": 1,
        "price": 599900
      },
      {
        "product_id": 2,
        "quantity": 1,
        "price": 129900
      }
    ]
  }'
```

**响应示例：**
```json
{
  "order": {
    "id": 1,
    "user_id": 1,
    "order_no": "ORD1704067200123456",
    "status": 1,
    "amount": 729800,
    "currency": "CNY",
    "description": "购买 iPhone 15 和 AirPods",
    "created_at": 1704067200,
    "updated_at": 1704067200,
    "paid_at": 0,
    "cancelled_at": 0
  }
}
```

**说明**：
- `amount` 字段（729800）由系统根据订单项自动计算：599900 + 129900 = 729800
- 订单项已自动创建并关联到该订单

#### 获取订单
```bash
curl http://localhost:8000/api/v1/orders/1
```

#### 更新订单状态
```bash
curl -X PUT http://localhost:8000/api/v1/orders/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": 2,
    "description": "订单已支付"
  }'
```

#### 取消订单
```bash
curl -X POST http://localhost:8000/api/v1/orders/1/cancel \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "用户主动取消"
  }'
```

#### 列出订单
```bash
# 基本列表
curl "http://localhost:8000/api/v1/orders?page=1&page_size=10"

# 按用户筛选
curl "http://localhost:8000/api/v1/orders?user_id=1"

# 按状态筛选
curl "http://localhost:8000/api/v1/orders?status=2"

# 时间范围筛选
curl "http://localhost:8000/api/v1/orders?start_time=1704067200&end_time=1704153600"

# 关键词搜索
curl "http://localhost:8000/api/v1/orders?keyword=ORD123456"

# 组合查询
curl "http://localhost:8000/api/v1/orders?page=1&page_size=20&user_id=1&status=2&start_time=1704067200&end_time=1704153600&keyword=ORD"
```

#### 删除订单
```bash
curl -X DELETE http://localhost:8000/api/v1/orders/1
```

### gRPC 请求示例

```go
import (
    "context"
    "sre/api/order/v1"
    "google.golang.org/grpc"
)

conn, _ := grpc.Dial("localhost:8989", grpc.WithInsecure())
client := v1.NewOrderClient(conn)

// 创建订单（单个产品）
order, err := client.CreateOrder(context.Background(), &v1.CreateOrderRequest{
    UserId:      1,
    Currency:    "CNY",
    Description: "购买 iPhone 15",
    Items: []*v1.OrderItemRequest{
        {
            ProductId: 1,
            Quantity:  1,
            Price:     599900,
        },
    },
})

// 创建订单（多个产品）
order, err := client.CreateOrder(context.Background(), &v1.CreateOrderRequest{
    UserId:      1,
    Currency:    "CNY",
    Description: "购买 iPhone 15 和 AirPods",
    Items: []*v1.OrderItemRequest{
        {
            ProductId: 1,
            Quantity:  1,
            Price:     599900,
        },
        {
            ProductId: 2,
            Quantity:  1,
            Price:     129900,
        },
    },
})

// 获取订单
order, err := client.GetOrder(context.Background(), &v1.GetOrderRequest{
    Id: 1,
})

// 更新订单状态
order, err := client.UpdateOrder(context.Background(), &v1.UpdateOrderRequest{
    Id:          1,
    Status:      v1.OrderStatus_ORDER_STATUS_PAID,
    Description: "订单已支付",
})

// 取消订单
order, err := client.CancelOrder(context.Background(), &v1.CancelOrderRequest{
    Id:     1,
    Reason: "用户主动取消",
})

// 列出订单
orders, err := client.ListOrders(context.Background(), &v1.ListOrdersRequest{
    Page:      1,
    PageSize:  10,
    UserId:    1,
    Status:    v1.OrderStatus_ORDER_STATUS_PENDING,
    Keyword:   "ORD",
    StartTime: 1704067200,
    EndTime:   1704153600,
})
```

---

## 业务规则

### 订单创建规则

1. **订单号生成**：系统自动生成唯一订单号，格式为 `ORD{timestamp}{nanosecond}`
2. **初始状态**：新订单默认状态为 `ORDER_STATUS_PENDING`（待支付）
3. **用户验证**：用户ID必须 > 0
4. **订单项要求**：
   - 订单项列表不能为空（至少包含一个订单项）
   - 每个订单项的产品必须存在
   - 每个订单项的数量必须 > 0
5. **金额计算**：
   - 订单总金额由系统根据订单项自动计算：`总金额 = Σ(订单项单价 × 数量)`
   - 如果订单项未提供价格或价格 <= 0，则使用产品当前价格
   - 计算后的总金额必须 > 0
6. **货币类型**：默认货币为 `CNY`，如果未提供则使用默认值
7. **订单项创建**：创建订单成功后，系统会自动批量创建对应的订单项

### 订单状态流转规则

```
待支付 (1) 
  ↓ [支付]
已支付 (2)
  ↓ [发货]
已发货 (3)
  ↓ [确认收货]
已完成 (4)

待支付 (1) 
  ↓ [取消]
已取消 (5)

已支付 (2)
  ↓ [退款]
已退款 (6)
```

**状态变更限制**：
- 已取消的订单不能再次取消
- 已完成的订单不能取消
- 状态变更时会自动设置相应的时间戳（`paid_at`、`cancelled_at`）

### 订单查询规则

1. **分页**：默认每页 10 条，最大 100 条
2. **关键词搜索**：支持搜索订单号和描述
3. **时间范围筛选**：支持按创建时间范围筛选
4. **排序**：按 ID 降序排列（最新创建的在前）

---

## 订单与订单项的关系

### 当前实现

当前版本的订单 API 中，创建订单时**支持直接传入订单项（OrderItem）**：

- **订单（Order）**：包含订单基本信息（用户、金额、状态等）
- **订单项（OrderItem）**：包含订单中的具体产品信息（产品ID、数量、单价等）

### 创建包含产品的订单

创建订单时，可以直接在请求中传入订单项列表：

```bash
curl -X POST http://localhost:8000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "购买 iPhone 15 和 AirPods",
    "items": [
      {
        "product_id": 1,
        "quantity": 1,
        "price": 599900
      },
      {
        "product_id": 2,
        "quantity": 1,
        "price": 129900
      }
    ]
  }'
```

**系统处理流程：**

1. **验证订单项**：
   - 验证产品是否存在
   - 验证数量是否 > 0
   - 如果价格未提供或 <= 0，使用产品当前价格

2. **计算总金额**：
   - 根据订单项自动计算：`总金额 = Σ(订单项单价 × 数量)`
   - 示例：599900 × 1 + 129900 × 1 = 729800

3. **创建订单**：
   - 生成订单号
   - 保存订单基本信息

4. **创建订单项**：
   - 批量创建订单项并关联到订单
   - 如果订单项创建失败，订单创建也会失败（保证数据一致性）

### 订单项数据结构

每个订单项包含以下信息：

- **product_id**：产品ID（必填，产品必须存在）
- **quantity**：商品数量（必填，必须 > 0）
- **price**：单价（可选，如果未提供或 <= 0，则使用产品当前价格）

订单项的价格是下单时的价格快照，即使后续产品价格发生变化，订单项的价格也不会改变。

---

## 认证和授权

### 当前配置

- **认证要求**：`/api/v1/orders` 及其子路径**不需要认证**（公开接口）
- **限流策略**：公开接口限流（100 请求/分钟）

### 修改认证配置

如需为订单接口添加认证，请修改 `internal/server/http.go` 中的 `isAuthenticatedRoute` 函数和 `createAuthMiddlewareWithSkip` 函数。

---

## 相关文档

- [API 设计规范](../../code-standards/api-design.md)
- [用户 API 列表](./user-api-list.md)
- [产品 API 列表](./product-api-list.md)
- [订单业务逻辑](../../internal/biz/order.go)
- [订单项业务逻辑](../../internal/biz/order_item.go)
- [用户订单产品查询设计](../../development/user-order-product-query-design.md)
