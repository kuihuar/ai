# 推送技术

即时通讯中的推送通知技术实现。

## 目录

- [推送概述](#推送概述)
- [iOS 推送 (APNs)](#ios-推送-apns)
- [Android 推送 (FCM)](#android-推送-fcm)
- [Web 推送](#web-推送)
- [推送策略](#推送策略)

---

## 推送概述

### 推送的必要性

**场景：**
- 用户应用在后台
- 用户离线
- 需要及时通知

**挑战：**
- 不同平台机制不同
- 需要保持连接
- 电池消耗
- 可靠性保证

---

## iOS 推送 (APNs)

### APNs 工作原理

```
应用服务器 → APNs → iOS 设备
```

**流程：**
1. 应用注册推送
2. 获取 Device Token
3. 服务器发送到 APNs
4. APNs 推送到设备

### 实现示例

**服务端发送：**
```go
type APNsClient struct {
    client *apns2.Client
}

func (c *APNsClient) Send(deviceToken string, payload map[string]interface{}) error {
    notification := &apns2.Notification{
        DeviceToken: deviceToken,
        Topic:       "com.example.app",
        Payload:     payload,
    }
    
    res, err := c.client.Push(notification)
    if err != nil {
        return err
    }
    
    if !res.Sent() {
        return fmt.Errorf("failed to send: %s", res.Reason)
    }
    
    return nil
}
```

### 推送格式

```json
{
    "aps": {
        "alert": {
            "title": "新消息",
            "body": "您有一条新消息"
        },
        "sound": "default",
        "badge": 1
    },
    "custom": {
        "message_id": "123",
        "type": "text"
    }
}
```

---

## Android 推送 (FCM)

### FCM 工作原理

```
应用服务器 → FCM → Android 设备
```

**流程：**
1. 应用注册 FCM
2. 获取 FCM Token
3. 服务器发送到 FCM
4. FCM 推送到设备

### 实现示例

```go
import "firebase.google.com/go/messaging"

func sendFCMNotification(client *messaging.Client, token string, title, body string) error {
    message := &messaging.Message{
        Token: token,
        Notification: &messaging.Notification{
            Title: title,
            Body:  body,
        },
        Data: map[string]string{
            "message_id": "123",
            "type":       "text",
        },
    }
    
    _, err := client.Send(context.Background(), message)
    return err
}
```

---

## Web 推送

### Web Push API

**特点：**
- 浏览器原生支持
- 需要 Service Worker
- 跨平台

**实现：**
```javascript
// 注册 Service Worker
navigator.serviceWorker.register('/sw.js');

// 订阅推送
const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: urlBase64ToUint8Array(publicKey)
});

// 发送订阅信息到服务器
await fetch('/api/push/subscribe', {
    method: 'POST',
    body: JSON.stringify(subscription)
});
```

---

## 推送策略

### 1. 推送优先级

```go
type PushPriority int

const (
    PriorityHigh PushPriority = iota
    PriorityNormal
    PriorityLow
)

func pushWithPriority(message *Message, priority PushPriority) {
    switch priority {
    case PriorityHigh:
        // 立即推送，带声音
        pushImmediately(message, sound: true)
    case PriorityNormal:
        // 正常推送
        pushNormal(message)
    case PriorityLow:
        // 延迟推送，静音
        pushDelayed(message, sound: false)
    }
}
```

### 2. 免打扰策略

```go
func shouldPush(userID int64, message *Message) bool {
    // 检查免打扰时间
    if isQuietHours(userID) {
        return false
    }
    
    // 检查免打扰设置
    if isMuted(userID, message.FromUser) {
        return false
    }
    
    // 检查消息优先级
    if message.Priority < getMinPriority(userID) {
        return false
    }
    
    return true
}
```

### 3. 推送去重

```go
func pushWithDeduplication(userID int64, message *Message) {
    key := fmt.Sprintf("push:%d:%s", userID, message.ID)
    
    // 检查是否已推送
    if exists, _ := cache.Exists(key); exists {
        return
    }
    
    // 推送
    push(message)
    
    // 标记已推送（5分钟过期）
    cache.Set(key, "1", 5*time.Minute)
}
```

### 4. 批量推送

```go
func batchPush(messages []*Message) {
    // 按用户分组
    userMessages := make(map[int64][]*Message)
    for _, msg := range messages {
        userMessages[msg.ToUser] = append(userMessages[msg.ToUser], msg)
    }
    
    // 批量推送
    for userID, msgs := range userMessages {
        go pushToUser(userID, msgs)
    }
}
```

---

## 推送优化

### 1. 连接池

```go
type PushClientPool struct {
    clients chan *PushClient
    maxSize int
}

func (p *PushClientPool) Get() *PushClient {
    select {
    case client := <-p.clients:
        return client
    default:
        return NewPushClient()
    }
}

func (p *PushClientPool) Put(client *PushClient) {
    select {
    case p.clients <- client:
    default:
        client.Close()
    }
}
```

### 2. 异步推送

```go
func pushAsync(message *Message) {
    go func() {
        if err := push(message); err != nil {
            // 重试或记录日志
            retryPush(message)
        }
    }()
}
```

### 3. 推送统计

```go
func recordPushStats(message *Message, success bool) {
    stats := &PushStats{
        UserID:    message.ToUser,
        MessageID: message.ID,
        Success:   success,
        Timestamp: time.Now(),
    }
    
    // 记录到数据库或监控系统
    statsService.Record(stats)
}
```

---

## 参考资料

- [APNs 官方文档](https://developer.apple.com/documentation/usernotifications)
- [FCM 官方文档](https://firebase.google.com/docs/cloud-messaging)
- [Web Push API](https://developer.mozilla.org/en-US/docs/Web/API/Push_API)

