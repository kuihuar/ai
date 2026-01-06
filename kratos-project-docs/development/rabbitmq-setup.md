# RabbitMQ 集成 - 设置步骤

## 已完成的工作

1. ✅ 实现了 `NewRabbitMQProducer` - RabbitMQ 消息生产者
2. ✅ 实现了 `NewRabbitMQConsumer` - RabbitMQ 消息消费者
3. ✅ 实现了 `NewRabbitMQEventPublisher` - RabbitMQ 事件发布器（与 KafkaEventPublisher 接口一致）
4. ✅ 添加到 `internal/data/data.go` 的 `ProviderSet` 中
5. ✅ 配置文件已包含 RabbitMQ 配置（`configs/config.yaml`）
6. ✅ Protobuf 定义已添加（`internal/conf/conf.proto`）

## 需要手动完成的操作

### 1. 添加 RabbitMQ 依赖

```bash
go get github.com/rabbitmq/amqp091-go
go mod tidy
```

### 2. 重新生成 Protobuf 代码

```bash
# 方式1：使用 Makefile（推荐）
make api

# 方式2：直接使用 protoc
protoc --proto_path=. --proto_path=third_party --go_out=paths=source_relative:. internal/conf/conf.proto
```

**重要：** 如果不重新生成 protobuf 代码，会出现编译错误：
```
c.Rabbitmq undefined (type *conf.Data has no field or method Rabbitmq)
```

### 3. 重新生成 Wire 代码（如果需要）

如果修改了 `wire.go` 文件，需要重新生成 Wire 代码：

```bash
# 在相应的 cmd 目录下运行
cd cmd/sre
go generate ./...

# 或者使用 Wire 工具
go run github.com/google/wire/cmd/wire ./...
```

## 代码位置

- **Producer/Consumer 实现**: `internal/data/client_rabbitmq.go`
- **Provider 注册**: `internal/data/data.go` (ProviderSet)
- **配置定义**: `internal/conf/conf.proto` (message Rabbitmq)
- **配置文件**: `configs/config.yaml` (data.rabbitmq)

## 使用示例

### Producer 示例

```go
// 通过 Wire 依赖注入
producer, err := data.NewRabbitMQProducer(confData, logger)
if err != nil {
    // 处理错误
}
defer producer.Close()

// 发布消息
err := producer.Publish(ctx, "exchange-name", "routing-key", []byte("message body"))
```

### Consumer 示例

```go
// 通过 Wire 依赖注入
consumer, err := data.NewRabbitMQConsumer(confData, logger)
if err != nil {
    // 处理错误
}
defer consumer.Close()

// 消费消息
handler := func(ctx context.Context, body []byte) error {
    // 处理消息
    log.Info("received message: %s", string(body))
    return nil
}

err := consumer.Consume(ctx, handler)
```

### Event Publisher 示例（与 Kafka 接口一致）

```go
// 通过 Wire 依赖注入
producer, _ := data.NewRabbitMQProducer(confData, logger)
publisher := data.NewRabbitMQEventPublisher(producer, logger)

// 发布事件
err := publisher.Publish(ctx, "exchange-name", "routing-key", eventPayload)
```

## 验证步骤

1. 添加依赖：`go get github.com/rabbitmq/amqp091-go`
2. 生成 protobuf：`make api`
3. 编译项目：`make build`
4. 启动 RabbitMQ 服务器（如果未运行）
5. 更新 `configs/config.yaml` 中的 RabbitMQ 配置
6. 运行应用并测试

## 相关文档

- [RabbitMQ 集成说明](./rabbitmq-integration.md) - 详细的使用文档
- [Kafka 集成说明](./kafka-integration.md) - Kafka 集成文档（待创建）
- [Outbox 模式](./outbox-testing-guide.md) - Outbox 模式测试指南

