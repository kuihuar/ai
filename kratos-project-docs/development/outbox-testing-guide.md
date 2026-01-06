# Outbox 模式端到端测试指南

本文档说明如何验证 Outbox 模式的完整流程：从订单创建到事件投递。

## 前置条件

1. **数据库已配置并运行**
   - MySQL 数据库已启动
   - 已执行数据库迁移（包含 `outbox_events` 表）

2. **应用已编译**
   ```bash
   make build
   ```

3. **配置已就绪**
   - `configs/config.yaml` 中数据库配置正确

## 测试步骤

### 步骤 1：启动 cron-worker

启动 cron-worker，它会每 10 秒扫描一次 Outbox 表：

```bash
./bin/cron-worker -conf ./configs/config.yaml
```

你应该看到类似日志：
```
INFO msg="starting cron manager with 2 jobs"
INFO msg="registered cron job: sync-user, spec: 0 0 2 * * *"
INFO msg="registered cron job: outbox-dispatcher, spec: */10 * * * * *"
```

### 步骤 2：创建订单（触发 Outbox 事件写入）

有两种方式创建订单：

#### 方式 A：通过 gRPC API（如果已实现）

```bash
grpcurl -plaintext -d '{
  "user_id": 1,
  "amount": 10000,
  "currency": "CNY",
  "description": "测试订单"
}' localhost:9000 api.order.v1.OrderService/CreateOrder
```

#### 方式 B：通过 HTTP API（如果已实现）

```bash
curl -X POST http://localhost:8000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 10000,
    "currency": "CNY",
    "description": "测试订单"
  }'
```

#### 方式 C：直接调用代码（用于开发测试）

创建一个简单的测试脚本 `scripts/test-outbox.go`：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"sre/internal/biz"
	"sre/internal/data"
	// ... 其他导入
)

func main() {
	// 初始化依赖（简化版，实际应该使用 Wire）
	// ...
	
	orderUsecase := biz.NewOrderUsecase(orderRepo, orderItemRepo, productRepo, logger)
	
	ctx := context.Background()
	order, err := orderUsecase.CreateOrder(ctx, 1, 10000, "CNY", "测试订单", nil)
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	
	fmt.Printf("订单创建成功: ID=%d, OrderNo=%s\n", order.ID, order.OrderNo)
}
```

### 步骤 3：验证数据库中的记录

#### 3.1 检查订单表

```sql
SELECT * FROM orders ORDER BY id DESC LIMIT 1;
```

应该看到刚创建的订单记录。

#### 3.2 检查 Outbox 表

```sql
SELECT * FROM outbox_events ORDER BY id DESC LIMIT 1;
```

应该看到：
- `aggregate_type` = "order"
- `aggregate_id` = 订单ID（字符串格式）
- `event_type` = "OrderCreated"
- `status` = 1（PENDING）
- `payload` 包含订单信息的 JSON

### 步骤 4：观察 cron-worker 日志

等待最多 10 秒（因为 cron 表达式是 `*/10 * * * * *`），你应该看到类似日志：

```
INFO msg="starting outbox dispatch cycle"
INFO msg="found 1 pending outbox events"
INFO msg="dispatching outbox event: id=1, event_id=order:created:ORD1234567890, aggregate=order/1, type=OrderCreated"
INFO msg="outbox dispatch cycle completed"
```

### 步骤 5：验证事件状态已更新

再次查询 Outbox 表：

```sql
SELECT * FROM outbox_events WHERE id = <刚才的ID>;
```

应该看到：
- `status` = 2（SENT）
- `updated_at` 已更新

## 验证清单

- [ ] cron-worker 成功启动，OutboxDispatchJob 已注册
- [ ] 创建订单后，`outbox_events` 表中有新记录
- [ ] Outbox 记录的 `status` = 1（PENDING）
- [ ] Outbox 记录的 `payload` 包含正确的订单信息
- [ ] cron-worker 日志显示扫描到待处理事件
- [ ] cron-worker 日志显示事件已处理
- [ ] Outbox 记录的 `status` 已更新为 2（SENT）

## 故障排查

### 问题 1：cron-worker 启动失败

**可能原因**：
- 数据库连接失败
- 配置错误

**解决方法**：
- 检查 `configs/config.yaml` 中的数据库配置
- 检查数据库是否运行
- 查看日志中的错误信息

### 问题 2：创建订单后，Outbox 表中没有记录

**可能原因**：
- 事务回滚
- `SaveWithEvent` 方法未正确调用

**解决方法**：
- 检查订单创建日志，确认是否调用了 `SaveWithEvent`
- 检查数据库事务日志
- 在 `repo_order.go` 的 `SaveWithEvent` 中添加更多日志

### 问题 3：cron-worker 没有处理事件

**可能原因**：
- OutboxDispatchJob 未注册
- cron 表达式配置错误
- 查询条件不匹配

**解决方法**：
- 检查 cron-worker 启动日志，确认 OutboxDispatchJob 已注册
- 检查 `outbox_dispatcher.go` 中的查询逻辑
- 手动执行一次 `ListPending` 查询，验证是否能查到数据

### 问题 4：事件状态未更新为 SENT

**可能原因**：
- `MarkSent` 方法执行失败
- 数据库更新失败

**解决方法**：
- 检查 cron-worker 日志中的错误信息
- 检查 `outbox_repo.go` 的 `MarkSent` 实现
- 验证数据库连接是否正常

## 下一步

完成基本验证后，可以：

1. **集成真实的 MQ**：在 `OutboxDispatchJob` 中接入 Kafka/RabbitMQ
2. **添加更多事件类型**：为订单支付、取消等操作添加 Outbox 事件
3. **实现 Saga 模式**：基于 Outbox 实现跨服务的分布式事务编排
4. **添加监控和告警**：监控 Outbox 事件的处理延迟和失败率

## 相关文档

- [分布式事件一致性架构设计](../architecture/distributed-event-consistency.md)
- [Job 任务调度实现方案](./job-implementation.md)

