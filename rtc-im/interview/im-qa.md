# IM 面试问答

即时通讯（Instant Messaging）相关面试题和答案。

## 目录

- [基础概念](#基础概念)
- [消息系统](#消息系统)
- [在线状态](#在线状态)
- [推送技术](#推送技术)
- [存储设计](#存储设计)
- [系统设计](#系统设计)

---

## 基础概念

### Q1: IM 系统的核心功能有哪些？

**答案：**

1. **消息收发**
   - 单聊、群聊
   - 消息类型（文本、图片、文件、语音、视频）
   - 消息状态（发送中、已发送、已读）

2. **在线状态**
   - 在线/离线
   - 最后在线时间
   - 正在输入

3. **消息存储**
   - 消息历史
   - 离线消息
   - 消息搜索

4. **推送通知**
   - 离线推送
   - 消息提醒
   - 免打扰

5. **其他功能**
   - 好友管理
   - 群组管理
   - 文件传输

### Q2: IM 系统的主要技术挑战是什么？

**答案：**

1. **消息可靠性**
   - 保证消息不丢失
   - 消息重传机制
   - 消息确认机制

2. **消息顺序性**
   - 保证消息顺序
   - 处理并发消息
   - 处理网络乱序

3. **高并发**
   - 大量用户同时在线
   - 消息高并发
   - 连接管理

4. **低延迟**
   - 消息实时到达
   - 减少网络延迟
   - 优化处理流程

5. **存储和扩展**
   - 海量消息存储
   - 历史消息查询
   - 数据分片

---

## 消息系统

### Q3: 如何保证消息的可靠性？

**答案：**

**多级确认机制：**

1. **发送端确认**
   ```go
   发送消息 → 等待 ACK → 标记已发送
   如果超时未收到 ACK → 重传
   ```

2. **服务端确认**
   - 消息持久化后发送 ACK
   - 保证消息不丢失

3. **接收端确认**
   - 收到消息后发送 ACK
   - 标记已读状态

**实现策略：**
- 消息 ID 唯一标识
- 重传机制（指数退避）
- 消息去重
- 幂等性保证

### Q4: 如何保证消息的顺序性？

**答案：**

**方案1：单连接单线程**
- 每个用户一个连接
- 服务端单线程处理
- 简单但性能受限

**方案2：消息序列号**
- 每条消息分配序列号
- 客户端按序列号排序
- 处理乱序消息

**方案3：分区有序**
- 按会话分区
- 每个分区保证顺序
- 不同分区可并行

**实现细节：**
```go
type Message struct {
    ID        int64
    SeqID     int64  // 序列号
    SessionID string
    Content   string
    Timestamp int64
}

// 客户端排序
func sortMessages(messages []Message) {
    sort.Slice(messages, func(i, j int) bool {
        return messages[i].SeqID < messages[j].SeqID
    })
}
```

### Q5: 如何处理离线消息？

**答案：**

**方案1：服务端存储**
- 用户离线时，消息存储在服务端
- 用户上线时，拉取离线消息
- 优点：可靠
- 缺点：存储压力大

**方案2：消息队列**
- 使用消息队列暂存
- 用户上线后消费
- 优点：解耦
- 缺点：需要额外组件

**实现流程：**
```
用户A发送消息给用户B
  ↓
用户B在线？
  ├─ 是 → 直接推送
  └─ 否 → 存储到离线消息表
         ↓
      用户B上线
         ↓
      拉取离线消息
         ↓
      标记为已读
```

### Q6: 群聊消息如何实现？

**答案：**

**方案1：写扩散（Write Fan-out）**
- 发送者写入所有成员的消息箱
- 优点：读简单
- 缺点：写压力大

**方案2：读扩散（Read Fan-out）**
- 消息只存储一份
- 成员读取时聚合
- 优点：写简单
- 缺点：读复杂

**方案3：混合方案**
- 小群：写扩散
- 大群：读扩散
- 根据群大小选择策略

**实现示例：**
```go
// 写扩散
func sendGroupMessage(groupID string, message Message) {
    members := getGroupMembers(groupID)
    for _, member := range members {
        saveMessage(member.UserID, message)
    }
}

// 读扩散
func getGroupMessages(groupID string, userID string) []Message {
    // 从群消息表读取
    // 过滤掉用户已删除的消息
    return queryGroupMessages(groupID, userID)
}
```

---

## 在线状态

### Q7: 如何实现用户在线状态？

**答案：**

**方案1：心跳机制**
```go
// 客户端定期发送心跳
func heartbeat() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        sendHeartbeat()
    }
}

// 服务端检测
func checkOnlineStatus() {
    // 如果 60 秒内没有心跳，标记为离线
    if time.Since(lastHeartbeat) > 60*time.Second {
        setOffline(userID)
    }
}
```

**方案2：连接状态**
- WebSocket 连接断开 → 离线
- 连接建立 → 在线
- 简单但不够准确

**方案3：混合方案**
- 连接状态 + 心跳
- 更准确的状态判断

### Q8: 如何实现"正在输入"功能？

**答案：**

**实现方式：**
1. 客户端检测输入事件
2. 发送"正在输入"状态
3. 服务端转发给接收方
4. 接收方显示状态
5. 3 秒后自动取消

**优化：**
- 防抖：避免频繁发送
- 节流：限制发送频率
- 状态缓存：避免重复发送

```go
var typingTimer *time.Timer

func onInput() {
    if typingTimer != nil {
        typingTimer.Stop()
    }
    
    sendTypingStatus(true)
    
    typingTimer = time.AfterFunc(3*time.Second, func() {
        sendTypingStatus(false)
    })
}
```

---

## 推送技术

### Q9: 离线推送的实现方案？

**答案：**

**推送渠道：**

1. **APNs (Apple Push Notification Service)**
   - iOS 官方推送
   - 需要证书配置
   - 可靠性高

2. **FCM (Firebase Cloud Messaging)**
   - Android 官方推送
   - 免费使用
   - 支持多平台

3. **第三方推送**
   - 极光推送
   - 个推
   - 小米推送

**推送策略：**
- 消息优先级
- 免打扰时间
- 推送内容摘要
- 推送去重

### Q10: 如何实现消息推送的优先级？

**答案：**

**优先级分类：**
1. **高优先级**：@ 消息、系统通知
2. **中优先级**：普通消息
3. **低优先级**：群消息（非 @）

**实现：**
```go
type MessagePriority int

const (
    PriorityHigh MessagePriority = iota
    PriorityMedium
    PriorityLow
)

func pushMessage(message Message) {
    switch message.Priority {
    case PriorityHigh:
        // 立即推送，带声音
        pushImmediately(message, sound: true)
    case PriorityMedium:
        // 正常推送
        pushNormal(message)
    case PriorityLow:
        // 延迟推送，静音
        pushDelayed(message, sound: false)
    }
}
```

---

## 存储设计

### Q11: 消息存储的数据库设计？

**答案：**

**表设计：**

```sql
-- 消息表
CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT,
    group_id BIGINT,
    content TEXT NOT NULL,
    msg_type INT NOT NULL,  -- 文本、图片、文件等
    seq_id BIGINT NOT NULL,
    status INT NOT NULL,    -- 发送中、已发送、已读
    created_at TIMESTAMP NOT NULL,
    INDEX idx_to_user (to_user_id, created_at),
    INDEX idx_group (group_id, created_at),
    INDEX idx_seq (seq_id)
);

-- 离线消息表
CREATE TABLE offline_messages (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    INDEX idx_user (user_id, created_at)
);

-- 已读状态表
CREATE TABLE read_status (
    user_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    read_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, message_id)
);
```

### Q12: 如何处理海量消息存储？

**答案：**

**方案1：分库分表**
- 按用户 ID 分片
- 按时间分表
- 冷热数据分离

**方案2：消息归档**
- 近期消息：热存储（MySQL）
- 历史消息：冷存储（对象存储）
- 定期归档

**方案3：读写分离**
- 写：主库
- 读：从库
- 减少主库压力

**实现示例：**
```go
// 分片策略
func getShardKey(userID int64) int {
    return int(userID % 16) // 16 个分片
}

// 归档策略
func archiveOldMessages() {
    // 归档 3 个月前的消息
    cutoffTime := time.Now().AddDate(0, -3, 0)
    archiveMessagesBefore(cutoffTime)
}
```

---

## 系统设计

### Q13: 设计一个支持千万级用户的 IM 系统？

**答案：**

**架构设计：**

1. **接入层**
   - WebSocket 网关
   - 负载均衡
   - 连接管理

2. **业务层**
   - 消息服务
   - 用户服务
   - 群组服务

3. **存储层**
   - 消息存储（分库分表）
   - 用户信息（Redis + MySQL）
   - 在线状态（Redis）

4. **基础设施**
   - 消息队列（Kafka）
   - 缓存（Redis）
   - 推送服务

**关键技术：**
- 水平扩展
- 消息队列解耦
- 缓存加速
- 数据库分片

### Q14: 如何实现消息的实时性？

**答案：**

**优化策略：**

1. **长连接**
   - WebSocket 保持连接
   - 减少连接建立开销

2. **消息队列**
   - 异步处理
   - 提高吞吐量

3. **缓存优化**
   - 热点数据缓存
   - 减少数据库查询

4. **就近部署**
   - CDN 加速
   - 边缘节点

5. **协议优化**
   - 二进制协议
   - 消息压缩

---

## 扩展问题

### Q15: 如何处理消息的撤回和删除？

**答案：**

**撤回：**
- 标记消息为已撤回
- 通知接收方
- 保留消息记录（审计）

**删除：**
- 软删除：标记删除
- 硬删除：物理删除
- 根据业务需求选择

### Q16: 如何实现消息的搜索功能？

**答案：**

**方案1：数据库全文索引**
- MySQL Full-Text Index
- 简单但性能有限

**方案2：搜索引擎**
- Elasticsearch
- 高性能搜索
- 支持复杂查询

**方案3：混合方案**
- 近期消息：数据库
- 历史消息：搜索引擎

---

## 参考资料

- [IM 系统设计实践](https://github.com/)
- [消息队列最佳实践](https://kafka.apache.org/)
- [WebSocket 协议](https://tools.ietf.org/html/rfc6455)

