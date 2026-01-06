# 在线状态管理

即时通讯中的用户在线状态管理实现。

## 目录

- [状态类型](#状态类型)
- [状态管理](#状态管理)
- [心跳机制](#心跳机制)
- [状态同步](#状态同步)

---

## 状态类型

### 基本状态

```go
type PresenceStatus int

const (
    StatusOffline PresenceStatus = iota  // 离线
    StatusOnline                          // 在线
    StatusAway                           // 离开
    StatusBusy                           // 忙碌
    StatusInvisible                      // 隐身
)
```

### 扩展状态

- 正在输入
- 最后在线时间
- 设备信息（手机、电脑、Web）

---

## 状态管理

### 状态存储

**Redis 存储：**
```go
// 设置在线状态
func setOnlineStatus(userID int64) error {
    key := fmt.Sprintf("presence:%d", userID)
    return redis.Set(key, "online", 60*time.Second)
}

// 获取在线状态
func getOnlineStatus(userID int64) (string, error) {
    key := fmt.Sprintf("presence:%d", userID)
    return redis.Get(key)
}
```

### 状态更新

```go
func updatePresence(userID int64, status PresenceStatus) error {
    // 更新 Redis
    key := fmt.Sprintf("presence:%d", userID)
    redis.Set(key, status, 60*time.Second)
    
    // 通知好友
    friends := getFriends(userID)
    for _, friendID := range friends {
        notifyPresenceChange(friendID, userID, status)
    }
    
    return nil
}
```

---

## 心跳机制

### 客户端心跳

```go
func heartbeat(client *Client) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := client.SendHeartbeat(); err != nil {
                // 连接断开
                return
            }
        case <-client.Done():
            return
        }
    }
}
```

### 服务端检测

```go
func checkPresence() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // 检查所有在线用户
        onlineUsers := redis.Keys("presence:*")
        for _, key := range onlineUsers {
            // 检查最后心跳时间
            lastHeartbeat := redis.Get(key + ":heartbeat")
            if time.Since(lastHeartbeat) > 60*time.Second {
                // 标记为离线
                setOffline(key)
            }
        }
    }
}
```

---

## 边界情况处理

### 问题场景分析

#### 1. 客户端进程被强制杀死

**场景：**
- 用户强制关闭应用
- 系统杀死进程（内存不足）
- 应用崩溃

**问题：**
- 无法发送心跳
- 无法发送断开连接消息
- 服务端认为用户仍在线

**解决方案：**

**方案1：连接状态检测（推荐）**
```go
// 服务端检测连接状态
func handleConnection(client *Client) {
    // 连接建立时设置在线状态
    setOnlineStatus(client.UserID)
    
    // 监听连接断开事件
    go func() {
        <-client.Conn.CloseNotify()  // 检测连接关闭
        // 连接断开，立即标记离线
        setOfflineStatus(client.UserID)
        cleanupConnection(client.UserID)
    }()
}

// WebSocket 连接关闭检测
func (c *WebSocketClient) CloseNotify() <-chan bool {
    ch := make(chan bool, 1)
    go func() {
        c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            ch <- true  // 连接已断开
        }
    }()
    return ch
}
```

**方案2：TCP Keepalive**
```go
// 设置 TCP Keepalive
func setKeepAlive(conn net.Conn) error {
    tcpConn := conn.(*net.TCPConn)
    if err := tcpConn.SetKeepAlive(true); err != nil {
        return err
    }
    if err := tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
        return err
    }
    return nil
}
```

**方案3：应用层心跳 + 超时检测**
```go
// 服务端心跳检测
type ConnectionManager struct {
    connections map[int64]*Connection
    mu          sync.RWMutex
}

type Connection struct {
    UserID         int64
    LastHeartbeat  time.Time
    Conn           net.Conn
    CloseChan      chan struct{}
}

func (cm *ConnectionManager) checkConnections() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        cm.mu.RLock()
        for userID, conn := range cm.connections {
            // 检查心跳超时
            if time.Since(conn.LastHeartbeat) > 90*time.Second {
                // 尝试发送探测包
                if !cm.probeConnection(conn) {
                    // 连接已断开，清理
                    go cm.cleanupConnection(userID)
                }
            }
        }
        cm.mu.RUnlock()
    }
}

// 探测连接是否存活
func (cm *ConnectionManager) probeConnection(conn *Connection) bool {
    // 发送 Ping 帧（WebSocket）
    if err := conn.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
        return false
    }
    
    // 尝试写入 Ping
    if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
        return false
    }
    
    // 等待 Pong 响应
    conn.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    _, _, err := conn.Conn.ReadMessage()
    return err == nil
}
```

#### 2. 网络突然断开

**场景：**
- WiFi 断开
- 移动网络切换
- 网络信号丢失

**问题：**
- 客户端无法感知断开
- 服务端无法立即检测

**解决方案：**

**多级检测机制：**
```go
type ConnectionHealthChecker struct {
    // 心跳超时：30秒
    HeartbeatTimeout time.Duration
    
    // 连接探测超时：60秒
    ProbeTimeout time.Duration
    
    // 最终超时：90秒
    FinalTimeout time.Duration
}

func (c *ConnectionHealthChecker) checkHealth(conn *Connection) {
    now := time.Now()
    
    // 第一级：心跳超时
    if now.Sub(conn.LastHeartbeat) > c.HeartbeatTimeout {
        // 发送探测包
        if !c.sendProbe(conn) {
            // 连接可能断开，进入第二级检测
            c.startProbeMode(conn)
        }
    }
}

func (c *ConnectionHealthChecker) startProbeMode(conn *Connection) {
    // 第二级：频繁探测
    probeTicker := time.NewTicker(5 * time.Second)
    defer probeTicker.Stop()
    
    probeCount := 0
    for range probeTicker.C {
        probeCount++
        if c.sendProbe(conn) {
            // 连接恢复
            conn.LastHeartbeat = time.Now()
            return
        }
        
        // 第三级：最终超时
        if time.Since(conn.LastHeartbeat) > c.FinalTimeout {
            // 确认断开，清理连接
            c.cleanupConnection(conn)
            return
        }
        
        // 限制探测次数
        if probeCount > 6 {
            c.cleanupConnection(conn)
            return
        }
    }
}
```

#### 3. 客户端休眠/锁屏

**场景：**
- 手机锁屏
- 应用进入后台
- 系统休眠

**问题：**
- 心跳可能停止
- 但应用未退出

**解决方案：**

**客户端处理：**
```go
// iOS/Android 后台心跳策略
func handleAppStateChange(state AppState) {
    switch state {
    case AppStateForeground:
        // 前台：正常心跳（30秒）
        startHeartbeat(30 * time.Second)
    case AppStateBackground:
        // 后台：降低心跳频率（5分钟）
        startHeartbeat(5 * time.Minute)
    case AppStateInactive:
        // 非活跃：发送最后心跳，标记为离开
        sendLastHeartbeat()
        setStatus(StatusAway)
    }
}
```

**服务端处理：**
```go
// 区分正常离线和临时离线
func handleHeartbeatTimeout(userID int64, lastHeartbeat time.Time) {
    timeout := time.Since(lastHeartbeat)
    
    if timeout > 5*time.Minute && timeout < 30*time.Minute {
        // 可能是休眠，标记为离开
        setStatus(userID, StatusAway)
    } else if timeout > 30*time.Minute {
        // 确认离线
        setOffline(userID)
    }
}
```

#### 4. 服务端重启/崩溃

**场景：**
- 服务端重启
- 服务端崩溃
- 网络分区

**问题：**
- 连接状态丢失
- 需要客户端重连

**解决方案：**

**客户端重连机制：**
```go
func reconnectWithBackoff(client *Client) {
    maxRetries := 10
    baseDelay := 1 * time.Second
    
    for i := 0; i < maxRetries; i++ {
        delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
        time.Sleep(delay)
        
        if err := client.Connect(); err == nil {
            // 重连成功，恢复状态
            client.RestoreState()
            return
        }
    }
    
    // 重连失败，提示用户
    notifyReconnectFailed()
}
```

**服务端状态恢复：**
```go
// 服务端启动时恢复连接状态
func restoreConnections() {
    // 从持久化存储恢复用户状态
    // 但标记为"待确认"状态
    users := loadPersistedUsers()
    for _, userID := range users {
        setStatus(userID, StatusPending)
    }
    
    // 等待客户端重连确认
    // 如果 5 分钟内未重连，标记为离线
    go func() {
        time.Sleep(5 * time.Minute)
        for userID := range pendingUsers {
            if !isConnected(userID) {
                setOffline(userID)
            }
        }
    }()
}
```

#### 5. 网络延迟/抖动

**场景：**
- 网络延迟突然增加
- 心跳包丢失
- 网络抖动

**问题：**
- 误判为离线
- 频繁状态切换

**解决方案：**

**容错机制：**
```go
type HeartbeatManager struct {
    // 允许连续丢失的心跳次数
    MaxMissedHeartbeats int
    
    // 心跳历史记录
    heartbeatHistory []time.Time
}

func (hm *HeartbeatManager) updateHeartbeat(timestamp time.Time) {
    hm.heartbeatHistory = append(hm.heartbeatHistory, timestamp)
    
    // 只保留最近 10 次
    if len(hm.heartbeatHistory) > 10 {
        hm.heartbeatHistory = hm.heartbeatHistory[len(hm.heartbeatHistory)-10:]
    }
}

func (hm *HeartbeatManager) shouldMarkOffline() bool {
    if len(hm.heartbeatHistory) == 0 {
        return true
    }
    
    lastHeartbeat := hm.heartbeatHistory[len(hm.heartbeatHistory)-1]
    missedDuration := time.Since(lastHeartbeat)
    
    // 允许连续丢失 3 次心跳（90秒）
    expectedHeartbeats := int(missedDuration / (30 * time.Second))
    return expectedHeartbeats > hm.MaxMissedHeartbeats
}
```

---

## 完整的心跳检测方案

### 多维度检测

```go
type PresenceDetector struct {
    // 连接状态检测
    connectionCheck func(conn *Connection) bool
    
    // 心跳超时检测
    heartbeatCheck func(userID int64) bool
    
    // 应用层探测
    probeCheck func(conn *Connection) bool
    
    // TCP Keepalive
    keepaliveEnabled bool
}

func (pd *PresenceDetector) detect(userID int64, conn *Connection) PresenceStatus {
    // 1. 首先检查连接状态（最快）
    if !pd.connectionCheck(conn) {
        return StatusOffline
    }
    
    // 2. 检查心跳超时
    if !pd.heartbeatCheck(userID) {
        // 心跳超时，但连接可能还在，进行探测
        if !pd.probeCheck(conn) {
            return StatusOffline
        }
        // 探测成功，可能是网络延迟
        return StatusAway
    }
    
    return StatusOnline
}
```

### 状态机设计

```go
type PresenceStateMachine struct {
    state     PresenceStatus
    lastCheck time.Time
    checkCount int
}

func (sm *PresenceStateMachine) transition(connAlive, heartbeatOK, probeOK bool) {
    switch sm.state {
    case StatusOnline:
        if !connAlive {
            sm.state = StatusOffline
        } else if !heartbeatOK {
            if probeOK {
                sm.state = StatusAway  // 临时离开
            } else {
                sm.state = StatusOffline
            }
        }
        
    case StatusAway:
        if !connAlive {
            sm.state = StatusOffline
        } else if heartbeatOK {
            sm.state = StatusOnline  // 恢复在线
        } else if !probeOK {
            sm.checkCount++
            if sm.checkCount > 3 {
                sm.state = StatusOffline
            }
        }
        
    case StatusOffline:
        if connAlive && heartbeatOK {
            sm.state = StatusOnline  // 重新上线
            sm.checkCount = 0
        }
    }
}
```

---

## 最佳实践总结

### 1. 多层检测机制

```
连接层检测（最快） → 心跳检测 → 应用层探测 → 最终确认
```

### 2. 容错设计

- 允许短暂的心跳丢失
- 区分临时离线和永久离线
- 避免频繁状态切换

### 3. 客户端配合

- 应用状态变化时主动通知
- 重连时恢复状态
- 优雅关闭时发送离线消息

### 4. 监控和告警

- 监控异常断开率
- 监控心跳延迟
- 告警异常模式

---

## 参考资料

- [WebSocket 连接管理](https://tools.ietf.org/html/rfc6455)
- [TCP Keepalive](https://tools.ietf.org/html/rfc1122)
- [心跳机制最佳实践](https://github.com/)

---

## 状态同步

### 好友状态同步

```go
func syncFriendPresence(userID int64) {
    friends := getFriends(userID)
    for _, friendID := range friends {
        status := getPresence(friendID)
        sendPresenceUpdate(userID, friendID, status)
    }
}
```

### 状态变更通知

```go
func notifyPresenceChange(userID int64, friendID int64, status PresenceStatus) {
    // 检查用户是否在线
    if !isOnline(userID) {
        return
    }
    
    // 发送状态更新
    message := &PresenceMessage{
        UserID: friendID,
        Status: status,
        Timestamp: time.Now(),
    }
    
    sendToUser(userID, message)
}
```

---

---

## 大规模客户端检测策略

### 问题：大量客户端的检测挑战

**场景：**
- 百万级在线用户
- 每秒数万心跳
- 服务端检测压力大

**挑战：**
- 检测效率
- 资源消耗
- 实时性
- 准确性

### 检测策略对比

| 策略 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **纯心跳检测** | 简单 | 服务端压力大 | 小规模（< 1万） |
| **连接状态检测** | 实时、准确 | 需要维护连接状态 | 中大规模（1-10万） |
| **混合检测** | 平衡 | 实现复杂 | 大规模（> 10万） |
| **事件驱动** | 高效 | 需要事件机制 | 超大规模（> 百万） |

---

## 大规模检测方案

### 方案1：连接状态优先（推荐）

**核心思想：**
- 不依赖心跳，优先检测连接状态
- 心跳作为辅助验证

**实现：**

```go
type ConnectionRegistry struct {
    // 连接映射：userID -> connection
    connections sync.Map
    
    // 心跳时间：userID -> lastHeartbeat
    heartbeats sync.Map
}

// 连接建立时注册
func (cr *ConnectionRegistry) Register(userID int64, conn net.Conn) {
    cr.connections.Store(userID, conn)
    
    // 监听连接关闭
    go cr.monitorConnection(userID, conn)
}

// 监控连接状态
func (cr *ConnectionRegistry) monitorConnection(userID int64, conn net.Conn) {
    // 方法1：使用 ReadDeadline 检测
    for {
        conn.SetReadDeadline(time.Now().Add(1 * time.Second))
        buf := make([]byte, 1)
        _, err := conn.Read(buf)
        if err != nil {
            // 连接断开
            cr.handleDisconnect(userID)
            return
        }
    }
}

// 方法2：使用 TCP Keepalive
func (cr *ConnectionRegistry) setupKeepAlive(conn net.Conn) error {
    tcpConn := conn.(*net.TCPConn)
    
    // 启用 Keepalive
    if err := tcpConn.SetKeepAlive(true); err != nil {
        return err
    }
    
    // 设置 Keepalive 参数
    // Linux: 30秒后开始探测，每10秒探测一次，最多3次
    if err := tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
        return err
    }
    
    return nil
}

// 处理断开
func (cr *ConnectionRegistry) handleDisconnect(userID int64) {
    cr.connections.Delete(userID)
    cr.heartbeats.Delete(userID)
    setOfflineStatus(userID)
}
```

**优势：**
- ✅ 实时检测（连接断开立即知道）
- ✅ 不依赖心跳（减少服务端压力）
- ✅ 准确性高（TCP 层检测）

**适用：**
- WebSocket 连接
- 长连接场景
- 大规模系统

### 方案2：事件驱动检测

**核心思想：**
- 连接断开时触发事件
- 批量处理状态更新

**实现：**

```go
type PresenceService struct {
    // 连接管理器
    connMgr *ConnectionManager
    
    // 状态更新队列
    statusQueue chan *StatusUpdate
    
    // 批量处理协程
    batchProcessor *BatchProcessor
}

// 连接断开事件
func (ps *PresenceService) onConnectionClose(userID int64) {
    // 发送到队列，异步处理
    ps.statusQueue <- &StatusUpdate{
        UserID: userID,
        Status: StatusOffline,
        Timestamp: time.Now(),
    }
}

// 批量处理状态更新
func (ps *PresenceService) processStatusUpdates() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    batch := make([]*StatusUpdate, 0, 1000)
    
    for {
        select {
        case update := <-ps.statusQueue:
            batch = append(batch, update)
            
            // 批量达到阈值或超时，处理一批
            if len(batch) >= 1000 {
                ps.batchUpdateStatus(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            // 定时处理
            if len(batch) > 0 {
                ps.batchUpdateStatus(batch)
                batch = batch[:0]
            }
        }
    }
}

// 批量更新状态
func (ps *PresenceService) batchUpdateStatus(updates []*StatusUpdate) {
    // 使用 Pipeline 批量更新 Redis
    pipe := redis.Pipeline()
    for _, update := range updates {
        key := fmt.Sprintf("presence:%d", update.UserID)
        if update.Status == StatusOffline {
            pipe.Del(key)
        } else {
            pipe.Set(key, update.Status, 60*time.Second)
        }
    }
    pipe.Exec()
}
```

### 方案3：分层检测机制

**核心思想：**
- 不同规模使用不同策略
- 动态调整检测频率

**实现：**

```go
type TieredPresenceDetector struct {
    // 活跃用户：频繁检测
    activeUsers map[int64]*ActiveUser
    
    // 普通用户：正常检测
    normalUsers map[int64]*NormalUser
    
    // 不活跃用户：低频检测
    inactiveUsers map[int64]*InactiveUser
}

type ActiveUser struct {
    UserID        int64
    LastActivity  time.Time
    CheckInterval time.Duration  // 10秒
}

type NormalUser struct {
    UserID        int64
    LastActivity  time.Time
    CheckInterval time.Duration  // 30秒
}

type InactiveUser struct {
    UserID        int64
    LastActivity  time.Time
    CheckInterval time.Duration  // 5分钟
}

func (tpd *TieredPresenceDetector) checkUsers() {
    // 分层检测
    go tpd.checkActiveUsers()   // 每10秒
    go tpd.checkNormalUsers()   // 每30秒
    go tpd.checkInactiveUsers() // 每5分钟
}

func (tpd *TieredPresenceDetector) checkActiveUsers() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        for userID, user := range tpd.activeUsers {
            if !tpd.isAlive(userID) {
                tpd.moveToOffline(userID)
            }
        }
    }
}
```

### 方案4：分布式检测

**核心思想：**
- 多服务器分布式检测
- 按用户 ID 分片

**实现：**

```go
type DistributedPresenceDetector struct {
    // 当前节点负责的用户范围
    shardID    int
    totalShards int
    
    // 用户分片
    userShards map[int][]int64
}

// 判断用户是否由本节点负责
func (dpd *DistributedPresenceDetector) isResponsible(userID int64) bool {
    shard := int(userID) % dpd.totalShards
    return shard == dpd.shardID
}

// 只检测本节点负责的用户
func (dpd *DistributedPresenceDetector) checkLocalUsers() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // 只获取本节点负责的用户
        localUsers := dpd.getLocalUsers()
        
        for _, userID := range localUsers {
            if !dpd.isAlive(userID) {
                dpd.setOffline(userID)
            }
        }
    }
}

// 使用一致性哈希分片
func (dpd *DistributedPresenceDetector) getShard(userID int64) int {
    // 一致性哈希
    hash := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%d", userID)))
    return int(hash) % dpd.totalShards
}
```

---

## 性能优化策略

### 1. 减少检测频率

**策略：**
- 连接状态检测：实时（事件驱动）
- 心跳检测：低频（30-60秒）
- 批量检测：批量处理

```go
// 优化前：每秒检测所有用户
func checkAllUsers() {
    for range time.Tick(1 * time.Second) {
        for userID := range allUsers {
            checkUser(userID)  // 百万用户，压力大
        }
    }
}

// 优化后：分片检测
func checkUsersSharded() {
    shards := 100  // 100个分片
    for i := 0; i < shards; i++ {
        go func(shardID int) {
            for range time.Tick(10 * time.Second) {
                // 只检测本分片的用户
                users := getUsersByShard(shardID)
                for _, userID := range users {
                    checkUser(userID)
                }
            }
        }(i)
    }
}
```

### 2. 使用 Redis 过期机制（TTL）

#### TTL 是什么？

**TTL (Time To Live)**：键的生存时间，过期后 Redis 会自动删除该键。

**Redis 中保存的数据：**

```go
// Redis 中保存的键值对
Key:   "presence:12345"           // 用户ID为12345的在线状态键
Value: "online"                   // 状态值（online/offline/away）
TTL:   90秒                        // 90秒后自动过期删除
```

**数据结构：**
```
Redis 内存中：
┌─────────────────┬──────────┬─────────┐
│ Key             │ Value    │ TTL     │
├─────────────────┼──────────┼─────────┤
│ presence:12345  │ "online" │ 90秒    │
│ presence:67890  │ "online" │ 85秒    │
│ presence:11111  │ "away"   │ 60秒    │
└─────────────────┴──────────┴─────────┘
```

#### TTL 工作机制

**1. 设置键值对时指定 TTL**
```go
// 用户上线时
func setOnlineWithTTL(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    // 设置键值对，并指定 90 秒后过期
    redis.Set(key, "online", 90*time.Second)
    
    // Redis 内部：
    // - 保存键值对：presence:12345 = "online"
    // - 启动倒计时：90秒
    // - 90秒后自动删除这个键
}
```

**2. TTL 倒计时过程**
```
时间轴：
0秒    → 设置键，TTL = 90秒
       → Redis: presence:12345 = "online" (剩余90秒)
30秒   → 收到心跳，刷新 TTL
       → Redis: presence:12345 = "online" (剩余90秒，重新计时)
60秒   → 收到心跳，刷新 TTL
       → Redis: presence:12345 = "online" (剩余90秒，重新计时)
90秒   → 如果没收到心跳，TTL 到期
       → Redis: 自动删除 presence:12345
       → 键不存在 = 用户离线
```

**3. 心跳时刷新 TTL**
```go
// 收到心跳时，刷新过期时间
func refreshHeartbeat(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    
    // 刷新 TTL，重新设置为 90 秒
    redis.Expire(key, 90*time.Second)
    
    // 或者使用 SET 命令，同时更新值和 TTL
    redis.Set(key, "online", 90*time.Second)
}
```

**4. 检测在线状态**
```go
// 检测用户是否在线
func isOnline(userID int64) bool {
    key := fmt.Sprintf("presence:%d", userID)
    
    // 检查键是否存在
    exists, _ := redis.Exists(key)
    
    // 如果键存在 = 在线（TTL 还没过期）
    // 如果键不存在 = 离线（TTL 已过期或从未设置）
    return exists
}
```

#### TTL 的优势

**1. 自动清理**
- ✅ Redis 自动删除过期键
- ✅ 无需手动检测和清理
- ✅ 减少服务端压力

**2. 精确计时**
- ✅ Redis 内部精确计时
- ✅ 不依赖外部定时器
- ✅ 性能高效

**3. 原子操作**
- ✅ TTL 设置和检查是原子操作
- ✅ 线程安全
- ✅ 无竞态条件

#### 实际应用示例

**完整流程：**
```go
// 1. 用户上线
func userOnline(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    // 设置在线状态，90秒后自动过期
    redis.Set(key, "online", 90*time.Second)
    
    // Redis 中：
    // Key: "presence:12345"
    // Value: "online"
    // TTL: 90秒（倒计时开始）
}

// 2. 收到心跳
func onHeartbeat(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    // 刷新 TTL，重新设置为 90 秒
    redis.Expire(key, 90*time.Second)
    
    // Redis 中：
    // Key: "presence:12345"
    // Value: "online"（不变）
    // TTL: 90秒（重新开始倒计时）
}

// 3. 检查在线状态
func checkOnline(userID int64) bool {
    key := fmt.Sprintf("presence:%d", userID)
    exists, _ := redis.Exists(key)
    
    if exists {
        // 键存在 = 在线
        return true
    } else {
        // 键不存在 = 离线（TTL 已过期）
        return false
    }
}

// 4. 用户离线（手动或自动）
func userOffline(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    // 手动删除键
    redis.Del(key)
    
    // 或者等待 TTL 自动过期（90秒后）
}
```

#### TTL 命令示例

**Redis 命令：**
```bash
# 设置键值对，带 TTL
SET presence:12345 "online" EX 90

# 查看剩余 TTL（秒）
TTL presence:12345
# 返回：85（表示还有85秒过期）

# 刷新 TTL
EXPIRE presence:12345 90

# 检查键是否存在
EXISTS presence:12345
# 返回：1（存在）或 0（不存在）

# 查看键的值
GET presence:12345
# 返回："online"
```

#### TTL 状态说明

**TTL 返回值：**
- **正数**：剩余秒数（如 85 表示还有 85 秒过期）
- **-1**：键存在但没有设置过期时间
- **-2**：键不存在（已过期或被删除）

```go
// 检查 TTL
func checkTTL(userID int64) int {
    key := fmt.Sprintf("presence:%d", userID)
    ttl, _ := redis.TTL(key)
    
    switch ttl {
    case -2:
        // 键不存在（已过期）
        return -2
    case -1:
        // 键存在但没有过期时间（不应该出现）
        return -1
    default:
        // 剩余秒数
        return ttl
    }
}
```

#### 为什么使用 TTL？

**对比：不使用 TTL**
```go
// ❌ 需要主动检测
func checkAllUsers() {
    for userID := range allUsers {
        lastHeartbeat := getLastHeartbeat(userID)
        if time.Since(lastHeartbeat) > 90*time.Second {
            setOffline(userID)  // 需要手动删除
        }
    }
}
// 问题：需要遍历所有用户，压力大
```

**使用 TTL：**
```go
// ✅ 自动过期，无需主动检测
func isOnline(userID int64) bool {
    key := fmt.Sprintf("presence:%d", userID)
    exists, _ := redis.Exists(key)
    return exists
}
// 优势：Redis 自动处理，无需遍历
```

#### 注意事项

**1. TTL 精度**
- Redis TTL 精度为秒级
- 对于毫秒级需求，需要额外处理

**2. 内存管理**
- 过期键会占用内存直到被删除
- 大量过期键可能影响性能

**3. 持久化**
- TTL 信息不持久化到 RDB
- 重启后需要重新设置

---

## Redis 中保存的数据结构

### 在线状态数据

```go
// Redis 中的数据结构
type PresenceData struct {
    Key   string  // "presence:12345"
    Value string  // "online" | "offline" | "away"
    TTL   int     // 剩余秒数（90秒）
}
```

### 连接信息数据（可选）

```go
// 也可以保存连接信息
type ConnectionData struct {
    Key   string  // "connection:12345"
    Value string  // 连接ID或服务器ID
    TTL   int     // 连接超时时间
}
```

### 心跳时间数据（可选）

```go
// 保存最后心跳时间
type HeartbeatData struct {
    Key   string  // "heartbeat:12345"
    Value int64   // Unix 时间戳
    TTL   int     // 心跳超时时间
}
```

### 完整的数据模型

```go
// 用户上线时，设置多个键
func setUserOnline(userID int64, connID string) {
    // 1. 在线状态（主要）
    redis.Set(fmt.Sprintf("presence:%d", userID), "online", 90*time.Second)
    
    // 2. 连接信息（可选）
    redis.Set(fmt.Sprintf("connection:%d", userID), connID, 90*time.Second)
    
    // 3. 心跳时间（可选）
    redis.Set(fmt.Sprintf("heartbeat:%d", userID), time.Now().Unix(), 90*time.Second)
}

// 检查在线状态
func isUserOnline(userID int64) bool {
    // 只需要检查 presence 键
    key := fmt.Sprintf("presence:%d", userID)
    exists, _ := redis.Exists(key)
    return exists
}
```

---

## TTL 机制总结

### 核心概念

1. **TTL = Time To Live**：键的生存时间
2. **自动过期**：TTL 到期后 Redis 自动删除键
3. **无需轮询**：不需要主动检测，Redis 自动处理

### Redis 中保存什么？

**主要数据：**
- **Key**: `presence:{userID}` - 用户在线状态键
- **Value**: `"online"` - 状态值
- **TTL**: `90秒` - 过期时间

**工作原理：**
```
用户上线 → SET presence:12345 "online" EX 90
    ↓
收到心跳 → EXPIRE presence:12345 90（刷新 TTL）
    ↓
90秒无心跳 → Redis 自动删除键
    ↓
检查在线 → EXISTS presence:12345 → 返回 0（离线）
```

### 优势

- ✅ **自动清理**：无需手动删除
- ✅ **性能高效**：Redis 内部处理
- ✅ **减少压力**：不需要轮询检测
- ✅ **精确计时**：Redis 精确管理时间

### 3. 采样检测

**策略：**
- 不检测所有用户
- 采样检测 + 统计推断

```go
// 采样检测
func sampleCheck(sampleRate float64) {
    allUsers := getAllOnlineUsers()
    sampleSize := int(float64(len(allUsers)) * sampleRate)
    
    // 随机采样
    sampledUsers := randomSample(allUsers, sampleSize)
    
    for _, userID := range sampledUsers {
        if !isAlive(userID) {
            // 发现离线用户，标记
            setOffline(userID)
        }
    }
}
```

### 4. 异步批量处理

```go
type BatchPresenceChecker struct {
    checkQueue chan int64
    workers    int
}

func (bpc *BatchPresenceChecker) Start() {
    // 启动多个工作协程
    for i := 0; i < bpc.workers; i++ {
        go bpc.worker()
    }
}

func (bpc *BatchPresenceChecker) worker() {
    batch := make([]int64, 0, 100)
    ticker := time.NewTicker(1 * time.Second)
    
    for {
        select {
        case userID := <-bpc.checkQueue:
            batch = append(batch, userID)
            if len(batch) >= 100 {
                bpc.batchCheck(batch)
                batch = batch[:0]
            }
        case <-ticker.C:
            if len(batch) > 0 {
                bpc.batchCheck(batch)
                batch = batch[:0]
            }
        }
    }
}

func (bpc *BatchPresenceChecker) batchCheck(userIDs []int64) {
    // 批量检查连接状态
    pipe := redis.Pipeline()
    for _, userID := range userIDs {
        key := fmt.Sprintf("connection:%d", userID)
        pipe.Exists(key)
    }
    results, _ := pipe.Exec()
    
    // 批量更新状态
    for i, result := range results {
        if !result.(*redis.BoolCmd).Val() {
            // 连接不存在，标记离线
            setOffline(userIDs[i])
        }
    }
}
```

---

## 完整的大规模检测方案

### 架构设计

```
┌─────────────────────────────────────────┐
│         客户端层                          │
│  百万级客户端                             │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         接入层 (Gateway)                  │
│  ┌──────────┐  ┌──────────┐            │
│  │ Gateway1 │  │ Gateway2 │  ...       │
│  └──────────┘  └──────────┘            │
│  连接状态检测（事件驱动）                  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         业务层                            │
│  ┌──────────┐  ┌──────────┐            │
│  │ Presence │  │ Presence │  ...       │
│  │ Service1 │  │ Service2 │            │
│  └──────────┘  └──────────┘            │
│  分片检测（按用户ID）                      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         存储层                            │
│  Redis Cluster (分片存储)                 │
│  利用 TTL 自动过期                        │
└─────────────────────────────────────────┘
```

### 实现要点

**1. 连接状态优先**
```go
// 不依赖心跳，优先检测连接
func detectPresence(userID int64, conn net.Conn) {
    // 1. 连接状态检测（最快、最准确）
    if !isConnectionAlive(conn) {
        setOffline(userID)
        return
    }
    
    // 2. 心跳作为辅助验证（低频）
    if time.Since(lastHeartbeat) > 90*time.Second {
        // 心跳超时，但连接可能还在，探测
        if !probeConnection(conn) {
            setOffline(userID)
        }
    }
}
```

**2. 事件驱动更新**
```go
// 连接断开时立即更新
func onConnectionClose(userID int64) {
    // 立即标记离线，不等待心跳超时
    setOffline(userID)
    notifyFriends(userID, StatusOffline)
}
```

**3. 分片检测**
```go
// 按用户ID分片，每个节点只检测部分用户
func getShard(userID int64) int {
    return int(userID) % totalShards
}

// 只检测本节点负责的用户
func checkLocalUsers() {
    localUsers := getUsersByShard(myShardID)
    for _, userID := range localUsers {
        checkUser(userID)
    }
}
```

**4. Redis TTL 机制**
```go
// 利用 Redis 自动过期
func setOnline(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    // 90秒后自动过期（相当于离线）
    redis.Set(key, "online", 90*time.Second)
}

// 心跳时刷新
func heartbeat(userID int64) {
    key := fmt.Sprintf("presence:%d", userID)
    redis.Expire(key, 90*time.Second)
}
```

---

## 性能对比

### 方案对比

| 方案 | 检测方式 | 服务端压力 | 实时性 | 适用规模 |
|------|---------|-----------|--------|----------|
| **纯心跳** | 被动等待心跳 | 高（每秒数万心跳） | 差（延迟高） | < 1万 |
| **连接检测** | 主动检测连接 | 中（事件驱动） | 好（实时） | 1-10万 |
| **混合检测** | 连接+心跳 | 低（优化后） | 好 | 10-100万 |
| **分布式+TTL** | 分片+自动过期 | 极低 | 好 | > 100万 |

### 推荐方案（大规模）

**组合策略：**
1. **连接状态检测**（主要）
   - 事件驱动，实时更新
   - 不依赖心跳

2. **Redis TTL**（辅助）
   - 自动过期机制
   - 减少主动检测

3. **分片检测**（补充）
   - 分布式检测
   - 负载均衡

4. **心跳验证**（可选）
   - 低频验证
   - 网络质量检测

---

## 总结：大规模客户端检测核心要点

### 回答你的问题

**Q: 如果有很多客户端，服务端是如何检测的？依靠客户端的心跳吗？**

**A: 不完全依赖心跳，推荐使用连接状态检测为主，心跳为辅。**

### 核心策略

**1. 连接状态检测（主要方式）**
- ✅ **不依赖心跳**：通过 TCP/WebSocket 连接状态检测
- ✅ **事件驱动**：连接断开时立即触发事件
- ✅ **实时准确**：连接断开立即知道，无需等待心跳超时
- ✅ **服务端压力小**：不需要轮询所有用户

**2. Redis TTL 机制（兜底）**
- ✅ **自动过期**：设置 90 秒 TTL，过期自动删除
- ✅ **心跳刷新**：收到心跳时刷新 TTL
- ✅ **无需主动检测**：Redis 自动处理过期

**3. 分布式分片（扩展性）**
- ✅ **按用户分片**：每个节点只检测部分用户
- ✅ **负载均衡**：分散检测压力
- ✅ **水平扩展**：可以增加节点

**4. 心跳作为辅助（可选）**
- ⚠️ **低频验证**：30-60 秒一次，不是主要检测方式
- ⚠️ **网络质量**：用于检测网络质量，不是检测在线状态

### 实际工作流程

```
客户端连接建立
    ↓
服务端注册连接 + 设置 Redis TTL(90秒)
    ↓
监听连接关闭事件（TCP/WebSocket 层）
    ↓
[如果连接断开] → 立即标记离线（事件驱动，< 1秒）
    ↓
[如果连接正常] → 心跳刷新 TTL（30秒一次，辅助）
    ↓
[如果心跳停止但连接还在] → 探测连接 → 确认离线
    ↓
[如果 Redis TTL 过期] → 自动标记离线（兜底，90秒）
```

### 关键优势

| 特性 | 纯心跳方案 | 连接检测方案 |
|------|-----------|-------------|
| **服务端压力** | 高（每秒数万心跳） | 低（事件驱动） |
| **检测延迟** | 高（等待心跳超时） | 低（< 1秒） |
| **准确性** | 中（可能误判） | 高（TCP 层检测） |
| **扩展性** | 差 | 好（分布式） |

### 最佳实践

1. **优先使用连接状态检测**，不依赖心跳
2. **利用 Redis TTL** 作为兜底机制
3. **分布式分片** 处理大规模用户
4. **心跳作为辅助**，用于网络质量检测

---

## 状态同步

### 好友状态同步

```go
func syncFriendPresence(userID int64) {
    friends := getFriends(userID)
    for _, friendID := range friends {
        status := getPresence(friendID)
        sendPresenceUpdate(userID, friendID, status)
    }
}
```

### 状态变更通知

```go
func notifyPresenceChange(userID int64, friendID int64, status PresenceStatus) {
    // 检查用户是否在线
    if !isOnline(userID) {
        return
    }
    
    // 发送状态更新
    message := &PresenceMessage{
        UserID: friendID,
        Status: status,
        Timestamp: time.Now(),
    }
    
    sendToUser(userID, message)
}
```

---

## 参考资料

- [WebSocket 连接管理](https://tools.ietf.org/html/rfc6455)
- [TCP Keepalive](https://tools.ietf.org/html/rfc1122)
- [心跳机制最佳实践](https://github.com/)
- [大规模系统设计](https://github.com/)

