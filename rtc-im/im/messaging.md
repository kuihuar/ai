# 消息系统设计

即时通讯消息系统的设计和实现。

## 目录

- [消息模型](#消息模型)
- [消息流程](#消息流程)
- [可靠性保证](#可靠性保证)
- [顺序性保证](#顺序性保证)
- [存储设计](#存储设计)

---

## 消息模型

### 消息类型

**文本消息：**
```go
type TextMessage struct {
    ID        int64
    FromUser  string
    ToUser    string
    Content   string
    Timestamp int64
}
```

**图片消息：**
```go
type ImageMessage struct {
    ID        int64
    FromUser  string
    ToUser    string
    ImageURL  string
    ThumbURL  string
    Width     int
    Height    int
    Timestamp int64
}
```

**文件消息：**
```go
type FileMessage struct {
    ID        int64
    FromUser  string
    ToUser    string
    FileName  string
    FileURL   string
    FileSize  int64
    FileType  string
    Timestamp int64
}
```

### 消息状态

```go
type MessageStatus int

const (
    StatusSending MessageStatus = iota  // 发送中
    StatusSent                         // 已发送
    StatusDelivered                    // 已送达
    StatusRead                         // 已读
    StatusFailed                       // 发送失败
)
```

---

## 消息流程

### 单聊消息流程

```
发送方                   服务端                   接收方
  │                        │                        │
  │── 发送消息 ────────────>│                        │
  │                        │── 持久化 ──────────────>│
  │                        │                        │
  │<── ACK ────────────────│                        │
  │                        │                        │
  │                        │── 推送消息 ────────────>│
  │                        │<── ACK ────────────────│
  │                        │                        │
  │<── 已读回执 ───────────│<── 已读回执 ───────────│
```

### 群聊消息流程

**写扩散模式：**
```
发送方 → 服务端 → 写入所有成员消息箱 → 推送所有成员
```

**读扩散模式：**
```
发送方 → 服务端 → 写入群消息表 → 成员读取时聚合
```

---

## 可靠性保证

### 三级确认机制

**1. 发送端确认**
```go
func sendMessage(msg *Message) error {
    // 发送消息
    err := client.Send(msg)
    if err != nil {
        return err
    }
    
    // 等待 ACK
    ack := <-ackChannel
    if ack.Success {
        msg.Status = StatusSent
    } else {
        // 重传
        return retrySend(msg)
    }
    return nil
}
```

**2. 服务端确认**
```go
func handleMessage(msg *Message) error {
    // 持久化
    if err := db.Save(msg); err != nil {
        return err
    }
    
    // 发送 ACK
    sendACK(msg.ID, true)
    return nil
}
```

**3. 接收端确认**
```go
func receiveMessage(msg *Message) {
    // 处理消息
    processMessage(msg)
    
    // 发送 ACK
    sendACK(msg.ID, true)
    
    // 标记已读
    markAsRead(msg.ID)
}
```

### 重传机制

```go
func retrySend(msg *Message, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := sendMessage(msg)
        if err == nil {
            return nil
        }
        
        // 指数退避
        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    return errors.New("max retries exceeded")
}
```

---

## 顺序性保证

### 消息序列号

```go
type Message struct {
    ID        int64
    SeqID     int64  // 序列号
    SessionID string // 会话 ID
    Content   string
    Timestamp int64
}

// 服务端分配序列号
func generateSeqID(sessionID string) int64 {
    // 原子操作保证唯一性
    return atomic.AddInt64(&seqCounter, 1)
}
```

### 客户端排序

```go
// 客户端维护消息队列
type MessageQueue struct {
    messages map[int64]*Message
    expectedSeq int64
    mu sync.Mutex
}

func (q *MessageQueue) Add(msg *Message) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    q.messages[msg.SeqID] = msg
    
    // 处理连续的消息
    for {
        if msg, ok := q.messages[q.expectedSeq]; ok {
            processMessage(msg)
            delete(q.messages, q.expectedSeq)
            q.expectedSeq++
        } else {
            break
        }
    }
}
```

---

## 存储设计

### 数据库设计

**消息表：**
```sql
CREATE TABLE messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT,
    group_id BIGINT,
    content TEXT NOT NULL,
    msg_type INT NOT NULL,
    seq_id BIGINT NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    INDEX idx_to_user (to_user_id, created_at),
    INDEX idx_group (group_id, created_at),
    INDEX idx_seq (seq_id)
);
```

**离线消息表：**
```sql
CREATE TABLE offline_messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    INDEX idx_user (user_id, created_at)
);
```

### 分库分表策略

**按用户 ID 分片：**
```go
func getShardKey(userID int64) int {
    return int(userID % 16) // 16 个分片
}

func getTableName(userID int64, timestamp int64) string {
    shard := getShardKey(userID)
    month := time.Unix(timestamp, 0).Format("200601")
    return fmt.Sprintf("messages_%d_%s", shard, month)
}
```

### 冷热数据分离

```go
// 热数据：最近 3 个月
func getHotData(userID int64) []*Message {
    cutoff := time.Now().AddDate(0, -3, 0)
    return queryMessages(userID, cutoff)
}

// 冷数据：3 个月以前
func getColdData(userID int64, startTime time.Time) []*Message {
    // 从对象存储读取
    return readFromObjectStorage(userID, startTime)
}
```

---

## 性能优化

### 1. 批量操作

```go
// 批量发送消息
func batchSendMessages(messages []*Message) error {
    batch := make([]*Message, 0, 100)
    for _, msg := range messages {
        batch = append(batch, msg)
        if len(batch) >= 100 {
            if err := db.BatchInsert(batch); err != nil {
                return err
            }
            batch = batch[:0]
        }
    }
    if len(batch) > 0 {
        return db.BatchInsert(batch)
    }
    return nil
}
```

### 2. 缓存优化

```go
// 缓存最近消息
func getRecentMessages(userID int64) []*Message {
    key := fmt.Sprintf("messages:%d", userID)
    
    // 从缓存读取
    if msgs, err := cache.Get(key); err == nil {
        return msgs
    }
    
    // 从数据库读取
    msgs := db.QueryRecentMessages(userID)
    
    // 写入缓存
    cache.Set(key, msgs, 5*time.Minute)
    return msgs
}
```

### 3. 异步处理

```go
// 异步处理消息
func handleMessageAsync(msg *Message) {
    go func() {
        // 持久化
        db.Save(msg)
        
        // 推送
        pushService.Push(msg)
        
        // 统计
        statsService.Record(msg)
    }()
}
```

---

## 参考资料

- [消息队列最佳实践](https://kafka.apache.org/)
- [分布式系统设计](https://en.wikipedia.org/wiki/Distributed_computing)

