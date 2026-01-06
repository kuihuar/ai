# Emitter 消息网关

## 概述

Emitter 是一个基于 MQTT 协议的高性能分布式消息网关，专为实时消息传递设计。它提供了低延迟、高吞吐量的消息传递能力，特别适合物联网（IoT）、实时 Web 应用、游戏和移动应用等场景。

## 核心特性

### 1. 高性能

- **低延迟**：毫秒级消息传递延迟
- **高吞吐量**：支持百万级消息/秒的吞吐量
- **分布式架构**：支持水平扩展，可部署多个节点

### 2. MQTT 协议支持

- **MQTT 3.1.1 和 5.0**：完整支持 MQTT 协议标准
- **QoS 级别**：支持 QoS 0、1、2
- **保留消息**：支持消息持久化和保留
- **遗嘱消息**：支持客户端断开时的遗嘱消息

### 3. 安全特性

- **TLS/SSL 加密**：支持安全连接
- **密钥认证**：基于密钥的访问控制
- **频道权限**：细粒度的频道访问控制

### 4. 消息存储

- **消息持久化**：支持消息的持久化存储
- **消息历史**：支持查询历史消息
- **离线消息**：支持客户端离线时的消息缓存

## 适用场景

### ✅ 适合使用 Emitter

- **物联网（IoT）**：设备数据采集和控制
- **实时 Web 应用**：实时仪表板、数据可视化
- **在线游戏**：实时游戏状态同步
- **移动应用**：实时推送和通知
- **聊天系统**：实时消息传递
- **股票行情**：实时数据推送

### ❌ 不适合使用 Emitter

- **批量数据处理**：适合使用 Kafka
- **事务性消息**：适合使用 RabbitMQ
- **复杂路由**：适合使用 RabbitMQ 或 RocketMQ
- **消息持久化要求极高**：适合使用 Kafka

## 与其他消息中间件对比

| 特性 | Emitter | Kafka | RabbitMQ | MQTT Broker |
|------|---------|-------|----------|-------------|
| **协议** | MQTT | 自定义 | AMQP | MQTT |
| **延迟** | 极低（毫秒级） | 低 | 低 | 低 |
| **吞吐量** | 极高 | 极高 | 高 | 中 |
| **消息持久化** | ✅ | ✅ | ✅ | 部分支持 |
| **实时性** | ✅ 优秀 | ⚠️ 一般 | ⚠️ 一般 | ✅ 优秀 |
| **IoT 支持** | ✅ 优秀 | ❌ | ⚠️ 一般 | ✅ 优秀 |
| **Web 支持** | ✅ 优秀 | ❌ | ⚠️ 一般 | ⚠️ 一般 |
| **学习曲线** | 低 | 中 | 中 | 低 |
| **适用场景** | 实时推送、IoT | 大数据流 | 消息队列 | IoT |

## 架构设计

```
┌─────────────┐
│   Client    │
│  (Go/JS)    │
└──────┬──────┘
       │
       │ MQTT Protocol
       │
┌──────▼─────────────────────────────────┐
│        Emitter Gateway                  │
│  - MQTT Broker                         │
│  - Message Router                      │
│  - Key Manager                         │
│  - Storage Engine                      │
└──────┬─────────────────────────────────┘
       │
       ├─────────────┬─────────────┐
       │             │             │
┌──────▼─────┐  ┌────▼──────┐  ┌──▼──────┐
│  Node 1    │  │  Node 2   │  │ Node N  │
│ (Cluster)  │  │ (Cluster)  │  │(Cluster)│
└────────────┘  └───────────┘  └─────────┘
```

## 核心概念

### 1. Channel（频道）

Channel 是消息传递的主题，类似于 MQTT 的 Topic。Channel 使用层级结构命名：

```
device/12345/sensor/temperature
device/12345/sensor/humidity
user/67890/notifications
```

### 2. Key（密钥）

Key 用于控制对 Channel 的访问权限。每个 Key 可以配置：
- **读取权限**：允许订阅（Subscribe）
- **写入权限**：允许发布（Publish）
- **频道范围**：限制可访问的频道

### 3. Message（消息）

消息是传递的数据单元，支持二进制和文本格式。

### 4. QoS（服务质量）

支持 MQTT 的三种 QoS 级别：
- **QoS 0**：最多一次（At most once）
- **QoS 1**：至少一次（At least once）
- **QoS 2**：恰好一次（Exactly once）

## 详细文档

### 1. 安装和部署

参见 [Emitter 安装部署文档](./emitter-installation.md)

### 2. Go SDK 使用

参见 [Emitter Go SDK 使用文档](./emitter-go-sdk.md)

### 3. 与项目集成

参见 [Emitter 项目集成文档](./emitter-integration.md)

### 4. 最佳实践

参见 [Emitter 最佳实践文档](./emitter-best-practices.md)

## 快速开始

### 1. 安装 Emitter Server

```bash
# 使用 Docker
docker run -d -p 8080:8080 -p 443:443 emitter/server
```

### 2. 安装 Go SDK

```bash
go get github.com/emitter-io/emitter
```

### 3. 基本使用

```go
package main

import (
    "context"
    "log"
    
    "github.com/emitter-io/emitter"
)

func main() {
    // 连接到 Emitter
    client, err := emitter.Connect("tcp://localhost:8080", "your-key")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // 订阅频道
    client.Subscribe("device/+/sensor/+", func(msg *emitter.Message) {
        log.Printf("Received: %s", msg.Payload())
    })
    
    // 发布消息
    client.Publish("device/12345/sensor/temperature", []byte("25.5"))
}
```

## 与项目现有架构的集成

### 当前项目使用 Kafka 和 RabbitMQ

当前项目已经集成了：
- **Kafka**：用于事件发布（Outbox 模式）
- **RabbitMQ**：用于消息队列

### Emitter 的补充作用

Emitter 可以作为补充，用于：
- **实时推送**：向 Web 客户端推送实时数据
- **IoT 设备通信**：与物联网设备进行消息传递
- **实时通知**：实时通知和提醒

### 集成建议

1. **保持现有架构**：Kafka 和 RabbitMQ 继续用于后端服务间通信
2. **添加 Emitter**：用于前端实时通信和 IoT 设备
3. **统一事件总线**：可以考虑将 Emitter 作为实时事件总线的补充

## 性能指标

| 指标 | 数值 |
|------|------|
| **消息延迟** | < 1ms（本地网络） |
| **吞吐量** | 100万+ 消息/秒 |
| **并发连接** | 100万+ |
| **消息大小** | 支持任意大小（建议 < 256KB） |

## 参考资源

- [Emitter 官方网站](https://emitter.io/)
- [Emitter GitHub](https://github.com/emitter-io/emitter)
- [Emitter Go SDK](https://github.com/emitter-io/emitter)
- [MQTT 协议规范](https://mqtt.org/)

## 总结

Emitter 是一个高性能的 MQTT 消息网关，特别适合实时消息传递场景。它的主要优势：

- ✅ 极低的延迟和极高的吞吐量
- ✅ 完整的 MQTT 协议支持
- ✅ 适合 IoT 和实时 Web 应用
- ✅ 简单易用，学习曲线低

如果需要实时消息推送、IoT 设备通信或实时 Web 应用，Emitter 是一个很好的选择。

