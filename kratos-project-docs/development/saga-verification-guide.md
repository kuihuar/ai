# Saga 流程验证指南

本文档提供详细的步骤来验证 Saga 订单创建流程是否正常工作，包括成功场景和失败补偿场景。

## 目录

- [前置准备](#前置准备)
- [验证成功流程](#验证成功流程)
- [验证失败补偿流程](#验证失败补偿流程)
- [检查数据库状态](#检查数据库状态)
- [查看日志](#查看日志)
- [常见问题排查](#常见问题排查)

---

## 前置准备

### 1. 确保数据库已创建表

运行数据库迁移，确保以下表已创建：
- `saga_instances`
- `saga_steps`
- `orders`
- `order_items`
- `outbox_events`

```bash
# 如果使用 Ent 迁移
cd /Users/jianfenliu/Workspace/sre
make ent-generate
# 然后运行数据库迁移（根据你的迁移工具）
```

### 2. 启动服务

```bash
# 编译
go build ./cmd/sre

# 启动服务（根据你的启动方式）
./sre -conf configs/config.yaml
# 或者
make run
```

服务启动后，应该看到：
- gRPC 服务监听在 `:8989`
- HTTP 服务监听在 `:8000`

### 3. 准备测试数据

确保数据库中有：
- 至少一个用户（用于创建订单）
- 至少一个产品（用于订单项）

---

## 验证成功流程

### 方法 1: 使用 curl 调用 HTTP API

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

**预期响应**:
```json
{
  "order": {
    "id": 123,
    "user_id": 1,
    "order_no": "SAGA-1234567890-1",
    "status": 1,
    "amount": 10000,
    "currency": "CNY",
    "description": "测试订单（Saga模式）",
    "items": [
      {
        "id": 1,
        "product_id": 1,
        "quantity": 2,
        "price": 5000
      }
    ],
    "created_at": 1234567890,
    "updated_at": 1234567890
  },
  "saga_id": "SAGA-1234567890-1",
  "compensated": false
}
```

### 方法 2: 使用 grpcurl 调用 gRPC API

```bash
# 安装 grpcurl（如果未安装）
# brew install grpcurl

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

### 方法 3: 使用项目自带的客户端工具

如果项目有 `sre-client` 工具：

```bash
cd /Users/jianfenliu/Workspace/sre
go run ./cmd/sre-client --help
# 查看如何使用客户端工具调用 CreateOrderSaga
```

### 验证步骤

1. **检查 API 响应**
   - ✅ 返回 `200 OK` 或成功状态
   - ✅ `compensated = false`
   - ✅ 返回了订单信息

2. **检查日志输出**
   应该看到类似以下日志：
   ```
   INFO msg="CreateOrderSaga: userID=1, currency=CNY, itemsCount=1"
   INFO msg="Saga started: id=SAGA-xxx, type=order.create, steps=3"
   INFO msg="Saga step execute: saga_id=SAGA-xxx, step=create-order"
   INFO msg="Saga Step Execute: create-order, saga_id=SAGA-xxx, user_id=1"
   INFO msg="Saga step execute: saga_id=SAGA-xxx, step=reserve-inventory"
   INFO msg="Saga Step Execute: reserve-inventory, saga_id=SAGA-xxx, order_id=123"
   INFO msg="ReserveInventory: order_id=123, order_no=SAGA-xxx, items_count=1"
   INFO msg="Inventory reserved successfully: reserve_id=RESERVE-xxx"
   INFO msg="Saga step execute: saga_id=SAGA-xxx, step=freeze-payment"
   INFO msg="Saga Step Execute: freeze-payment, saga_id=SAGA-xxx, order_id=123"
   INFO msg="FreezePayment: order_id=123, order_no=SAGA-xxx, user_id=1, amount=10000"
   INFO msg="Payment frozen successfully: freeze_id=FREEZE-xxx"
   INFO msg="Saga completed: id=SAGA-xxx, type=order.create"
   ```

3. **检查数据库状态**
   见 [检查数据库状态](#检查数据库状态) 部分

---

## 验证失败补偿流程

### 场景 1: 模拟支付服务失败

由于当前是 Mock 实现，我们需要修改代码来模拟失败。创建一个测试版本：

#### 方法 A: 修改 Mock 客户端返回错误

临时修改 `internal/data/external/payment/client.go`:

```go
func (c *Client) FreezePayment(ctx context.Context, req *FreezePaymentRequest) (*FreezePaymentResponse, error) {
    // 临时添加：模拟失败
    if req.Amount > 100000 { // 金额超过 1000 元时失败
        return nil, fmt.Errorf("insufficient balance")
    }
    
    // ... 原有代码
}
```

然后调用 API，传入大金额：

```bash
curl -X POST http://localhost:8000/api/v1/orders/saga \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "测试订单（大金额，触发失败）",
    "items": [
      {
        "product_id": 1,
        "quantity": 100,
        "price": 2000
      }
    ]
  }'
```

**预期响应**:
```json
{
  "order": {
    "id": 124,
    ...
  },
  "saga_id": "SAGA-xxx",
  "compensated": true
}
```

**预期错误**: HTTP 500 或 gRPC 错误，但 `compensated = true`

#### 方法 B: 使用测试工具直接调用

创建一个测试脚本 `scripts/test-saga-failure.sh`:

```bash
#!/bin/bash

# 测试 Saga 失败场景
echo "测试 Saga 失败场景..."

# 调用 API（假设支付服务会失败）
curl -X POST http://localhost:8000/api/v1/orders/saga \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "测试失败场景",
    "items": [
      {
        "product_id": 1,
        "quantity": 1,
        "price": 1
      }
    ]
  }' 2>&1

echo ""
echo "检查日志和数据库状态..."
```

### 场景 2: 模拟库存服务失败

类似地，可以修改 `internal/data/external/inventory/client.go` 来模拟库存不足：

```go
func (c *Client) ReserveInventory(ctx context.Context, req *ReserveInventoryRequest) (*ReserveInventoryResponse, error) {
    // 临时添加：模拟失败
    if len(req.Items) > 10 { // 商品数量超过 10 时失败
        return &ReserveInventoryResponse{
            Success: false,
            Message: "insufficient inventory",
        }, nil
    }
    
    // ... 原有代码
}
```

### 验证补偿流程

1. **检查日志输出**
   应该看到补偿相关的日志：
   ```
   ERROR msg="Saga step failed: saga_id=SAGA-xxx, step=freeze-payment, err=..."
   WARN msg="Saga step compensate: saga_id=SAGA-xxx, step=reserve-inventory"
   INFO msg="Saga Step Compensate: reserve-inventory, saga_id=SAGA-xxx"
   INFO msg="ReleaseInventory: reserve_id=RESERVE-xxx, order_id=123"
   INFO msg="Inventory released successfully: reserve_id=RESERVE-xxx"
   WARN msg="Saga step compensate: saga_id=SAGA-xxx, step=create-order"
   INFO msg="Saga Step Compensate: create-order, saga_id=SAGA-xxx"
   INFO msg="Order cancelled successfully: order_id=123"
   WARN msg="Saga completed with compensation: id=SAGA-xxx"
   ```

2. **检查数据库状态**
   - Saga 实例状态应该是 `COMPENSED (5)`
   - 步骤状态：
     - Step1: `COMPENSATED (5)`
     - Step2: `COMPENSATED (5)`
     - Step3: `EXECUTE_FAILED (3)`

3. **检查订单状态**
   - 订单应该被取消（`status = CANCELLED (5)`）

---

## 检查数据库状态

### 1. 查询 Saga 实例

```sql
-- 查看最近的 Saga 实例
SELECT 
    id,
    saga_id,
    saga_type,
    status,
    metadata,
    started_at,
    completed_at,
    error_message,
    created_at,
    updated_at
FROM saga_instances
ORDER BY created_at DESC
LIMIT 10;
```

**状态值说明**:
- `1`: PENDING（待开始）
- `2`: RUNNING（执行中）
- `3`: COMPLETED（已完成）
- `4`: FAILED（失败待补偿）
- `5`: COMPENSED（已补偿完成）

### 2. 查询 Saga 步骤

```sql
-- 查看某个 Saga 的所有步骤
SELECT 
    id,
    saga_id,
    step_name,
    step_order,
    status,
    executed_at,
    compensated_at,
    error_message,
    created_at,
    updated_at
FROM saga_steps
WHERE saga_id = 'SAGA-1234567890-1'
ORDER BY step_order;
```

**状态值说明**:
- `0`: NOT_EXECUTED（未执行）
- `1`: EXECUTING（执行中）
- `2`: EXECUTED（执行成功）
- `3`: EXECUTE_FAILED（执行失败）
- `4`: COMPENSATING（补偿中）
- `5`: COMPENSATED（补偿成功）
- `6`: COMPENSATE_FAILED（补偿失败）

### 3. 查询订单状态

```sql
-- 查看订单是否被创建/取消
SELECT 
    id,
    user_id,
    order_no,
    status,
    amount,
    description,
    created_at,
    cancelled_at
FROM orders
WHERE order_no LIKE 'SAGA-%'
ORDER BY created_at DESC
LIMIT 10;
```

### 4. 查询 Outbox 事件

```sql
-- 查看订单创建事件
SELECT 
    id,
    event_id,
    aggregate_type,
    aggregate_id,
    event_type,
    status,
    created_at
FROM outbox_events
WHERE aggregate_type = 'order'
  AND aggregate_id = '123'
ORDER BY created_at DESC;
```

### 5. 完整查询脚本

创建一个 SQL 文件 `scripts/check-saga-status.sql`:

```sql
-- 检查最近的 Saga 实例
SELECT 
    si.saga_id,
    si.saga_type,
    CASE si.status
        WHEN 1 THEN 'PENDING'
        WHEN 2 THEN 'RUNNING'
        WHEN 3 THEN 'COMPLETED'
        WHEN 4 THEN 'FAILED'
        WHEN 5 THEN 'COMPENSED'
    END AS status_name,
    si.started_at,
    si.completed_at,
    si.error_message,
    COUNT(ss.id) AS steps_count,
    SUM(CASE WHEN ss.status = 2 THEN 1 ELSE 0 END) AS executed_count,
    SUM(CASE WHEN ss.status = 5 THEN 1 ELSE 0 END) AS compensated_count
FROM saga_instances si
LEFT JOIN saga_steps ss ON si.saga_id = ss.saga_id
GROUP BY si.id
ORDER BY si.created_at DESC
LIMIT 5;

-- 查看步骤详情
SELECT 
    ss.saga_id,
    ss.step_name,
    ss.step_order,
    CASE ss.status
        WHEN 0 THEN 'NOT_EXECUTED'
        WHEN 1 THEN 'EXECUTING'
        WHEN 2 THEN 'EXECUTED'
        WHEN 3 THEN 'EXECUTE_FAILED'
        WHEN 4 THEN 'COMPENSATING'
        WHEN 5 THEN 'COMPENSATED'
        WHEN 6 THEN 'COMPENSATE_FAILED'
    END AS status_name,
    ss.executed_at,
    ss.compensated_at,
    ss.error_message
FROM saga_steps ss
WHERE ss.saga_id IN (
    SELECT saga_id FROM saga_instances ORDER BY created_at DESC LIMIT 1
)
ORDER BY ss.step_order;
```

---

## 查看日志

### 1. 实时查看日志

如果使用文件日志：

```bash
# 查看实时日志
tail -f logs/sre.log | grep -i saga

# 或者查看所有日志
tail -f logs/sre.log
```

### 2. 过滤关键日志

```bash
# 只看 Saga 相关日志
tail -f logs/sre.log | grep -E "(Saga|saga|SAGA)"

# 只看错误日志
tail -f logs/sre.log | grep -E "(ERROR|WARN)"

# 只看补偿相关日志
tail -f logs/sre.log | grep -E "(compensate|Compensate|COMPENSATE)"
```

### 3. 日志关键字段

关注以下日志字段：
- `saga_id`: Saga 实例ID
- `step`: 步骤名称
- `order_id`: 订单ID
- `reserve_id`: 库存预留ID
- `freeze_id`: 支付冻结ID

### 4. 使用 jq 格式化 JSON 日志

如果日志是 JSON 格式：

```bash
tail -f logs/sre.log | jq 'select(.msg | contains("Saga"))'
```

---

## 常见问题排查

### 问题 1: API 调用返回 500 错误

**可能原因**:
- 数据库连接失败
- 表未创建
- 配置错误

**排查步骤**:
1. 检查数据库连接配置
2. 检查表是否存在
3. 查看错误日志详情

### 问题 2: Saga 执行成功但订单未创建

**可能原因**:
- 订单创建失败但 Saga 未捕获错误
- 事务未提交

**排查步骤**:
1. 检查订单表是否有记录
2. 检查 `outbox_events` 表是否有事件
3. 查看订单创建相关的日志

### 问题 3: 补偿未执行

**可能原因**:
- 补偿逻辑有错误
- 外部服务调用失败

**排查步骤**:
1. 检查 `saga_steps` 表的 `compensated_at` 字段
2. 查看补偿相关的日志
3. 检查外部服务是否可访问

### 问题 4: 状态持久化失败

**可能原因**:
- Ent 客户端未正确初始化
- 数据库事务问题

**排查步骤**:
1. 检查 `saga_instances` 和 `saga_steps` 表是否有记录
2. 查看数据库连接日志
3. 检查 Ent 客户端初始化代码

---

## 自动化测试脚本

创建一个测试脚本 `scripts/test-saga.sh`:

```bash
#!/bin/bash

set -e

echo "=== Saga 流程验证测试 ==="

BASE_URL="http://localhost:8000"

# 测试 1: 成功场景
echo ""
echo "测试 1: 成功创建订单（Saga）"
RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/orders/saga" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency": "CNY",
    "description": "自动化测试订单",
    "items": [
      {
        "product_id": 1,
        "quantity": 1,
        "price": 1000
      }
    ]
  }')

echo "响应: $RESPONSE"

# 提取 saga_id
SAGA_ID=$(echo $RESPONSE | jq -r '.saga_id')
if [ "$SAGA_ID" != "null" ] && [ -n "$SAGA_ID" ]; then
    echo "✅ 成功场景测试通过，Saga ID: $SAGA_ID"
else
    echo "❌ 成功场景测试失败"
    exit 1
fi

# 测试 2: 检查数据库状态
echo ""
echo "测试 2: 检查数据库状态"
# 这里可以添加数据库查询命令

echo ""
echo "=== 所有测试完成 ==="
```

运行测试：

```bash
chmod +x scripts/test-saga.sh
./scripts/test-saga.sh
```

---

## 性能测试

### 并发测试

使用 Apache Bench (ab) 或 wrk:

```bash
# 安装 wrk
# brew install wrk

# 并发测试
wrk -t4 -c100 -d30s --script=scripts/saga-test.lua http://localhost:8000/api/v1/orders/saga
```

创建 `scripts/saga-test.lua`:

```lua
wrk.method = "POST"
wrk.body   = '{"user_id":1,"currency":"CNY","description":"压力测试","items":[{"product_id":1,"quantity":1,"price":1000}]}'
wrk.headers["Content-Type"] = "application/json"
```

### 监控指标

关注以下指标：
- API 响应时间
- Saga 执行时间
- 数据库查询时间
- 外部服务调用时间
- 成功率/失败率

---

## 总结

验证 Saga 流程的关键步骤：

1. ✅ **API 调用成功** - 返回正确的响应
2. ✅ **日志输出正确** - 看到所有步骤的执行日志
3. ✅ **数据库状态正确** - Saga 实例和步骤状态正确
4. ✅ **补偿机制工作** - 失败时能正确补偿
5. ✅ **订单状态正确** - 成功时订单创建，失败时订单取消

如果所有验证都通过，说明 Saga 流程实现正确！

