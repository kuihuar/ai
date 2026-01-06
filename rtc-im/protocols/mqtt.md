# MQTT 协议

MQTT (Message Queuing Telemetry Transport) 协议详解。

## 目录

- [协议概述](#协议概述)
- [QoS 级别](#qos-级别)
- [主题和订阅](#主题和订阅)
- [实现示例](#实现示例)
- [应用场景](#应用场景)

---

## 协议概述

### 什么是 MQTT

MQTT 是一个轻量级的消息传输协议，专为低带宽、高延迟或不稳定的网络环境设计。

**特点：**
- 轻量级
- 发布/订阅模式
- 低带宽消耗
- 支持 QoS

### 协议架构

```
发布者 → MQTT Broker → 订阅者
```

**组件：**
- **Client**: 发布者或订阅者
- **Broker**: 消息代理服务器
- **Topic**: 主题
- **Message**: 消息

---

## QoS 级别

### QoS 0 - 最多一次

**特点：**
- 不保证送达
- 不重复
- 最低延迟

**适用场景：**
- 传感器数据
- 实时性要求高
- 可容忍丢失

### QoS 1 - 至少一次

**特点：**
- 保证送达
- 可能重复
- 需要确认

**流程：**
```
发布者 → PUBLISH → Broker → PUBACK
```

**适用场景：**
- 重要数据
- 需要保证送达
- 可容忍重复

### QoS 2 - 恰好一次

**特点：**
- 保证送达
- 不重复
- 最高可靠性

**流程：**
```
发布者 → PUBLISH → Broker → PUBREC
发布者 ← PUBREL ← Broker
发布者 → PUBCOMP → Broker
```

**适用场景：**
- 关键数据
- 不能丢失
- 不能重复

---

## 主题和订阅

### 主题 (Topic)

**主题格式：**
```
sensor/temperature/room1
sensor/humidity/room1
device/+/status
device/#/error
```

**通配符：**
- `+`: 单级通配符
- `#`: 多级通配符

### 订阅示例

```go
// 订阅单个主题
client.Subscribe("sensor/temperature", 1, nil)

// 订阅多个主题
topics := map[string]byte{
    "sensor/temperature": 1,
    "sensor/humidity":    1,
}
client.SubscribeMultiple(topics, nil)

// 使用通配符
client.Subscribe("device/+/status", 1, nil)
```

---

## 实现示例

### Go MQTT 客户端

```go
import "github.com/eclipse/paho.mqtt.golang"

func connectMQTT() (mqtt.Client, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("client1")
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    
    return client, nil
}

func publish(client mqtt.Client, topic string, payload []byte) error {
    token := client.Publish(topic, 1, false, payload)
    token.Wait()
    return token.Error()
}

func subscribe(client mqtt.Client, topic string, handler mqtt.MessageHandler) error {
    token := client.Subscribe(topic, 1, handler)
    token.Wait()
    return token.Error()
}
```

---

## 应用场景

### 1. IoT 设备

- 传感器数据采集
- 设备状态监控
- 远程控制

### 2. 移动应用

- 推送通知
- 实时消息
- 状态同步

### 3. IM 系统

- 消息推送
- 在线状态
- 系统通知

---

## 参考资料

- [MQTT 官方文档](https://mqtt.org/)
- [Eclipse Paho](https://www.eclipse.org/paho/)
- [MQTT 规范](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)

