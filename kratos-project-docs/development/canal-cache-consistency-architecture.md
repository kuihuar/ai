# Canal 架构缓存一致性详细设计

## 架构组件

### 1. Canal Server

**职责**：
- 订阅 MySQL Binlog
- 解析 Binlog 事件
- 将变更事件发送给 Canal Client

**部署方式**：
- 独立部署（推荐）
- 或集成到应用服务中

**关键配置**：
- MySQL 连接信息
- 订阅的表和数据库
- 事件过滤规则

### 2. Canal Client

**职责**：
- 连接 Canal Server
- 接收变更事件
- 解析事件并删除缓存
- 处理错误和重试

**部署方式**：
- 集成到应用服务中（推荐）
- 或独立部署为缓存同步服务

### 3. 缓存层（Redis）

**职责**：
- 存储业务数据缓存
- 提供缓存删除接口
- 支持批量删除

## 数据流设计

### 1. 变更事件流

```
MySQL Binlog
    │
    │ (Row Change Event)
    ▼
Canal Server
    │
    │ (Canal Entry)
    ▼
Canal Client
    │
    │ (解析为业务事件)
    ▼
缓存同步处理器
    │
    │ (删除/更新缓存)
    ▼
Redis
```

### 2. 事件格式

**Canal Entry 结构**：
```json
{
  "header": {
    "logfileName": "mysql-bin.000001",
    "logfileOffset": 12345,
    "executeTime": 1640000000000,
    "schemaName": "test",
    "tableName": "users",
    "eventType": "UPDATE"
  },
  "entryType": "ROWDATA",
  "storeValue": "..."
}
```

**解析后的业务事件**：
```go
type CacheSyncEvent struct {
    Database    string   // 数据库名
    Table       string   // 表名
    EventType   string   // INSERT/UPDATE/DELETE
    PrimaryKey  string   // 主键值
    RowData     map[string]interface{} // 行数据
    Timestamp   int64    // 时间戳
}
```

## 缓存同步策略

### 1. 单表缓存同步

**场景**：单表数据变更，只影响该表的缓存

**策略**：
```go
// 用户表更新
if table == "users" {
    userID := rowData["id"]
    cacheKey := fmt.Sprintf("user:%d", userID)
    redis.Del(cacheKey)
    
    // 同时删除相关缓存
    redis.Del(fmt.Sprintf("user:username:%s", rowData["username"]))
    redis.Del(fmt.Sprintf("user:email:%s", rowData["email"]))
}
```

### 2. 关联表缓存同步

**场景**：关联表数据变更，影响多个表的缓存

**策略**：
```go
// 订单表更新
if table == "orders" {
    orderID := rowData["id"]
    userID := rowData["user_id"]
    
    // 删除订单缓存
    redis.Del(fmt.Sprintf("order:%d", orderID))
    
    // 删除用户订单列表缓存
    redis.Del(fmt.Sprintf("user:%d:orders", userID))
    
    // 删除用户订单统计缓存
    redis.Del(fmt.Sprintf("user:%d:order:stats", userID))
}
```

### 3. 批量缓存同步

**场景**：批量更新操作，需要删除多个缓存

**策略**：
```go
// 批量删除用户缓存
if table == "users" && eventType == "UPDATE" {
    // 如果更新了影响列表的字段（如状态）
    if rowData["status"] != oldRowData["status"] {
        // 删除所有用户列表缓存
        redis.Del("user:list:*")
        redis.Del("user:list:active:*")
        redis.Del("user:list:inactive:*")
    }
}
```

## 错误处理

### 1. Canal Server 连接失败

**处理策略**：
- 自动重连（指数退避）
- 记录告警日志
- 降级到 Cache-Aside 模式

### 2. 事件解析失败

**处理策略**：
- 记录错误日志
- 跳过该事件（不影响其他事件）
- 发送告警通知

### 3. 缓存删除失败

**处理策略**：
- 重试删除（最多3次）
- 记录失败事件
- 异步补偿删除

## 性能优化

### 1. 批量处理

**策略**：累积多个事件后批量删除缓存

```go
// 批量删除缓存
func (h *CacheSyncHandler) BatchDeleteCache(events []CacheSyncEvent) {
    keys := make([]string, 0, len(events)*2)
    for _, event := range events {
        keys = append(keys, h.getCacheKeys(event)...)
    }
    redis.Del(keys...)
}
```

### 2. 异步处理

**策略**：缓存删除异步执行，不阻塞事件处理

```go
// 异步删除缓存
go func() {
    if err := redis.Del(cacheKey); err != nil {
        log.Errorf("failed to delete cache: %v", err)
    }
}()
```

### 3. 过滤规则

**策略**：只处理需要同步的表和字段

```go
// 配置过滤规则
filters := []FilterRule{
    {Database: "test", Table: "users", Fields: []string{"id", "username", "email"}},
    {Database: "test", Table: "orders", Fields: []string{"id", "user_id", "status"}},
}
```

## 一致性保证

### 1. 最终一致性

**说明**：
- Canal 基于 Binlog，保证所有数据变更都会被捕获
- 缓存删除可能延迟（通常 < 1 秒）
- 保证最终一致性，不保证强一致性

### 2. 顺序保证

**说明**：
- Canal 保证同一表的变更事件顺序
- 不同表的变更事件可能乱序
- 需要根据业务需求处理顺序问题

### 3. 幂等性

**说明**：
- 缓存删除操作天然幂等
- 重复删除不影响结果
- 可以安全重试

## 与 Cache-Aside 的配合

### 1. 混合使用策略

**高频热点数据**：
- 使用 Canal 自动同步
- 业务代码无需删除缓存

**低频数据**：
- 继续使用 Cache-Aside
- 业务代码显式删除缓存

### 2. 降级策略

**Canal 故障时**：
- 自动降级到 Cache-Aside
- 业务代码显式删除缓存
- 记录告警日志

## 部署架构

### 1. 单机部署

```
┌─────────────────┐
│  应用服务        │
│  ┌───────────┐  │
│  │Canal Client│  │
│  └───────────┘  │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  Canal Server   │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│     MySQL       │
└─────────────────┘
```

### 2. 集群部署

```
┌─────────────────┐     ┌─────────────────┐
│  应用服务 A      │     │  应用服务 B      │
│  ┌───────────┐  │     │  ┌───────────┐  │
│  │Canal Client│  │     │  │Canal Client│  │
│  └───────────┘  │     │  └───────────┘  │
└────────┬────────┘     └────────┬────────┘
         │                        │
         └──────────┬────────────┘
                    ▼
         ┌─────────────────┐
         │  Canal Server   │
         │   (集群)        │
         └─────────────────┘
                    │
                    ▼
         ┌─────────────────┐
         │  MySQL (主从)   │
         └─────────────────┘
```

## 安全考虑

### 1. 权限控制

**Canal Server**：
- 使用只读账号连接 MySQL
- 限制访问的数据库和表

**Canal Client**：
- 使用独立的 Redis 账号
- 限制删除操作的权限

### 2. 数据脱敏

**敏感数据**：
- 不在日志中记录敏感字段
- 缓存删除不记录完整数据

## 相关文档

- [Canal 方案概述](./canal-cache-consistency-overview.md) - 方案概述
- [Canal 实现指南](./canal-cache-consistency-implementation.md) - 具体实现步骤
- [Canal 配置说明](./canal-cache-consistency-config.md) - 配置参数说明
- [Canal 监控运维](./canal-cache-consistency-monitoring.md) - 监控和运维指南

