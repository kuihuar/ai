# RabbitMQ 集成说明

## 概述

本文档说明如何在项目中使用 RabbitMQ 作为消息队列，包括 Producer 和 Consumer 的实现。

## 前置条件

### 1. 安装 RabbitMQ 客户端库

```bash
go get github.com/rabbitmq/amqp091-go
```

### 2. 重新生成 Protobuf 代码

由于在 `conf.proto` 中添加了 `Rabbitmq` 配置，需要重新生成 protobuf 代码：

```bash
# 生成 conf.pb.go
make api
# 或者
protoc --proto_path=. --proto_path=third_party --go_out=paths=source_relative:. internal/conf/conf.proto
```

## 配置

### 1. 配置文件

在 `configs/config.yaml` 中添加 RabbitMQ 配置：

```yaml
data:
  rabbitmq:
    enable: true              # 是否启用 RabbitMQ
    url: amqp://guest:guest@127.0.0.1:5672/  # RabbitMQ 连接 URL
    queue: order-events       # 队列名称（用于 consumer）
    exchange: order-exchange  # 交换机名称（用于 producer，可选）
```

### 2. 配置字段说明

- `enable`: 是否启用 RabbitMQ（`true`/`false`）
- `url`: RabbitMQ 连接 URL，格式：`amqp://user:password@host:port/vhost`
- `queue`: 队列名称（用于 Consumer）
- `exchange`: 交换机名称（用于 Producer，可选）

## 使用方式

### 1. RabbitMQ Producer

**创建 Producer**：
```go
producer, err := data.NewRabbitMQProducer(confData, logger)
if err != nil {
    // 处理错误
}
defer producer.Close()
```

**发布消息**：
```go
err := producer.Publish(ctx, "exchange-name", "routing-key", []byte("message body"))
if err != nil {
    // 处理错误
}
```

**使用 Event Publisher**（与 Kafka 接口一致）：
```go
publisher := data.NewRabbitMQEventPublisher(producer, logger)
err := publisher.Publish(ctx, "exchange-name", "routing-key", []byte("event payload"))
```

### 2. RabbitMQ Consumer

**创建 Consumer**：
```go
consumer, err := data.NewRabbitMQConsumer(confData, logger)
if err != nil {
    // 处理错误
}
defer consumer.Close()
```

**消费消息**：
```go
handler := func(ctx context.Context, body []byte) error {
    // 处理消息
    log.Info("received message: %s", string(body))
    return nil
}

err := consumer.Consume(ctx, handler)
if err != nil {
    // 处理错误
}
```

## 与 Kafka 的对比

| 特性 | Kafka | RabbitMQ |
|------|-------|----------|
| **消息模型** | Topic/Partition | Exchange/Queue |
| **路由方式** | 基于 Partition | 基于 Routing Key |
| **消息持久化** | 支持 | 支持 |
| **消息确认** | Offset | ACK/NACK |
| **适用场景** | 高吞吐量日志流 | 复杂路由规则 |

## 在 Outbox 模式中使用

### 使用 RabbitMQ Producer

```go
// 在 OutboxDispatchJob 中使用 RabbitMQ
rabbitmqProducer, err := data.NewRabbitMQProducer(confData, logger)
if err != nil {
    // 处理错误
}
rabbitmqPublisher := data.NewRabbitMQEventPublisher(rabbitmqProducer, logger)

// 发布事件
err := rabbitmqPublisher.Publish(ctx, "order-exchange", "order.created", event.Payload)
```

### 使用 RabbitMQ Consumer

```go
// 创建 Consumer
consumer, err := data.NewRabbitMQConsumer(confData, logger)
if err != nil {
    // 处理错误
}

// 消费消息
handler := func(ctx context.Context, body []byte) error {
    // 解析事件
    var event OrderCreatedEvent
    if err := json.Unmarshal(body, &event); err != nil {
        return err
    }
    
    // 处理事件
    return processOrderCreatedEvent(ctx, &event)
}

err := consumer.Consume(ctx, handler)
```

## 注意事项

### 1. 连接管理

- Producer 和 Consumer 都会创建独立的连接和 Channel
- 使用完毕后务必调用 `Close()` 方法关闭连接
- 连接断开后需要重新创建

### 2. 消息确认

- Consumer 使用手动确认模式（`autoAck: false`）
- 处理成功后调用 `msg.Ack(false)` 确认消息
- 处理失败时调用 `msg.Nack(false, true)` 拒绝消息并重新入队

### 3. 队列声明

- Consumer 会自动声明队列（如果不存在）
- 队列设置为持久化（`durable: true`）
- 建议在生产环境中预先创建队列

### 4. 错误处理

- 连接失败时返回错误
- 消息发布失败时返回错误
- 消息处理失败时拒绝消息并重新入队

## 部署建议

### 1. 开发环境

```yaml
rabbitmq:
  enable: true
  url: amqp://guest:guest@127.0.0.1:5672/
  queue: order-events
  exchange: order-exchange
```

### 2. 生产环境

```yaml
rabbitmq:
  enable: true
  url: amqp://user:password@rabbitmq.example.com:5672/vhost
  queue: order-events-prod
  exchange: order-exchange-prod
```

### 3. 高可用配置

- 使用 RabbitMQ 集群
- 配置镜像队列
- 使用连接池
- 实现连接重试机制

## 相关文档

- [Kafka 集成说明](./kafka-integration.md) - Kafka 集成文档（待创建）
- [Outbox 模式](./outbox-testing-guide.md) - Outbox 模式测试指南
- [消息队列对比](./mq-comparison.md) - 消息队列对比文档（待创建）

