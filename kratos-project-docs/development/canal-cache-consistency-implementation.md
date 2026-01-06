# Canal 缓存一致性实现指南

## 前置条件

### 1. MySQL 配置

**启用 Binlog**：
```ini
[mysqld]
log-bin=mysql-bin
binlog-format=ROW
binlog-row-image=FULL
server-id=1
```

**创建 Canal 账号**：
```sql
CREATE USER 'canal'@'%' IDENTIFIED BY 'canal_password';
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';
FLUSH PRIVILEGES;
```

### 2. Canal Server 部署

**下载 Canal**：
```bash
wget https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.deployer-1.1.7.tar.gz
tar -xzf canal.deployer-1.1.7.tar.gz
cd canal.deployer-1.1.7
```

**配置 Canal Server**：
```properties
# conf/canal.properties
canal.instance.master.address=127.0.0.1:3306
canal.instance.dbUsername=canal
canal.instance.dbPassword=canal_password
canal.instance.filter.regex=test\\..*
```

**启动 Canal Server**：
```bash
./bin/startup.sh
```

## 代码实现

### 1. 添加依赖

**go.mod**：
```go
require (
    github.com/go-mysql-org/go-mysql v1.7.0
    github.com/siddontang/go-log v0.0.0-20190221022429-1e957dd83bed
)
```

### 2. Canal Client 实现

**创建 Canal Client**：
```go
// internal/data/canal/client.go
package canal

import (
    "context"
    "fmt"
    "time"

    "github.com/go-mysql-org/go-mysql/canal"
    "github.com/go-mysql-org/go-mysql/mysql"
    "github.com/go-mysql-org/go-mysql/replication"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/redis/go-redis/v9"
)

// CacheSyncHandler 缓存同步处理器
type CacheSyncHandler struct {
    redis  *redis.Client
    logger *log.Helper
}

// NewCacheSyncHandler 创建缓存同步处理器
func NewCacheSyncHandler(redis *redis.Client, logger log.Logger) *CacheSyncHandler {
    return &CacheSyncHandler{
        redis:  redis,
        logger: log.NewHelper(logger),
    }
}

// OnRow 处理行变更事件
func (h *CacheSyncHandler) OnRow(e *canal.RowsEvent) error {
    ctx := context.Background()
    
    // 解析事件
    event := &CacheSyncEvent{
        Database:  e.Table.Schema,
        Table:     e.Table.Name,
        EventType: e.Action,
        Timestamp: time.Now().Unix(),
    }
    
    // 根据表名选择处理策略
    switch event.Table {
    case "users":
        return h.handleUserEvent(ctx, event, e)
    case "orders":
        return h.handleOrderEvent(ctx, event, e)
    default:
        h.logger.Debugf("ignoring table: %s", event.Table)
        return nil
    }
}

// handleUserEvent 处理用户表变更
func (h *CacheSyncHandler) handleUserEvent(ctx context.Context, event *CacheSyncEvent, e *canal.RowsEvent) error {
    keys := make([]string, 0)
    
    // 解析主键和字段
    for _, row := range e.Rows {
        // 获取主键值（假设主键是 id）
        userID := row[0].(int64)
        
        // 构建缓存键
        keys = append(keys, fmt.Sprintf("user:%d", userID))
        
        // 如果是 UPDATE，获取旧值
        if e.Action == "update" && len(e.Rows) > 1 {
            oldRow := e.Rows[0]
            newRow := e.Rows[1]
            
            // 如果用户名或邮箱变更，删除相关缓存
            if oldRow[1] != newRow[1] { // username
                keys = append(keys, fmt.Sprintf("user:username:%v", oldRow[1]))
            }
            if oldRow[2] != newRow[2] { // email
                keys = append(keys, fmt.Sprintf("user:email:%v", oldRow[2]))
            }
        }
    }
    
    // 批量删除缓存
    if len(keys) > 0 {
        if err := h.redis.Del(ctx, keys...).Err(); err != nil {
            h.logger.Errorf("failed to delete cache: keys=%v, error=%v", keys, err)
            return err
        }
        h.logger.Infof("cache deleted: keys=%v", keys)
    }
    
    return nil
}

// handleOrderEvent 处理订单表变更
func (h *CacheSyncHandler) handleOrderEvent(ctx context.Context, event *CacheSyncEvent, e *canal.RowsEvent) error {
    keys := make([]string, 0)
    userIDs := make(map[int64]bool)
    
    // 解析订单数据
    for _, row := range e.Rows {
        orderID := row[0].(int64)
        userID := row[1].(int64) // user_id
        
        // 删除订单缓存
        keys = append(keys, fmt.Sprintf("order:%d", orderID))
        
        // 记录用户ID，用于删除用户订单列表缓存
        userIDs[userID] = true
    }
    
    // 删除用户订单列表缓存
    for userID := range userIDs {
        keys = append(keys, fmt.Sprintf("user:%d:orders", userID))
        keys = append(keys, fmt.Sprintf("user:%d:order:stats", userID))
    }
    
    // 批量删除缓存
    if len(keys) > 0 {
        if err := h.redis.Del(ctx, keys...).Err(); err != nil {
            h.logger.Errorf("failed to delete cache: keys=%v, error=%v", keys, err)
            return err
        }
        h.logger.Infof("cache deleted: keys=%v", keys)
    }
    
    return nil
}

// CacheSyncEvent 缓存同步事件
type CacheSyncEvent struct {
    Database  string
    Table     string
    EventType string
    Timestamp int64
}
```

### 3. Canal Client 启动

**启动 Canal Client**：
```go
// internal/data/canal/manager.go
package canal

import (
    "context"
    "fmt"

    "github.com/go-mysql-org/go-mysql/canal"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/redis/go-redis/v9"
)

// Manager Canal Client 管理器
type Manager struct {
    canal   *canal.Canal
    handler *CacheSyncHandler
    logger  *log.Helper
}

// NewManager 创建 Canal Manager
func NewManager(redis *redis.Client, logger log.Logger) (*Manager, error) {
    cfg := canal.NewDefaultConfig()
    cfg.Addr = "127.0.0.1:3306"
    cfg.User = "canal"
    cfg.Password = "canal_password"
    cfg.Dump.ExecutionPath = ""
    
    c, err := canal.NewCanal(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create canal: %w", err)
    }
    
    handler := NewCacheSyncHandler(redis, logger)
    
    // 注册事件处理器
    c.SetEventHandler(handler)
    
    return &Manager{
        canal:   c,
        handler: handler,
        logger:  log.NewHelper(logger),
    }, nil
}

// Start 启动 Canal Client
func (m *Manager) Start(ctx context.Context) error {
    m.logger.Info("starting canal client")
    
    // 启动 Canal
    go func() {
        if err := m.canal.Run(); err != nil {
            m.logger.Errorf("canal client error: %v", err)
        }
    }()
    
    // 等待上下文取消
    <-ctx.Done()
    m.logger.Info("stopping canal client")
    
    return m.canal.Close()
}
```

### 4. 集成到应用

**在 data.go 中集成**：
```go
// internal/data/data.go
package data

import (
    "sre/internal/data/canal"
    // ...
)

// Data 结构体
type Data struct {
    db      *gorm.DB
    ent     *ent.Client
    redis   *redis.Client
    locker  DistributedLocker
    canal   *canal.Manager  // 新增
}

// NewData 创建 Data
func NewData(
    db *gorm.DB,
    entClient *ent.Client,
    redisClient *redis.Client,
    locker DistributedLocker,
    canalManager *canal.Manager,  // 新增
    logger log.Logger,
) (*Data, func(), error) {
    // ...
    
    data := &Data{
        db:     db,
        ent:    entClient,
        redis:  redisClient,
        locker: locker,
        canal:  canalManager,  // 新增
    }
    
    // 启动 Canal Client
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        if err := canalManager.Start(ctx); err != nil {
            logHelper.Errorf("canal client error: %v", err)
        }
    }()
    
    cleanup := func() {
        cancel()
        // ... 其他清理逻辑
    }
    
    return data, cleanup, nil
}
```

## 配置管理

### 1. 配置文件

**config.yaml**：
```yaml
data:
  canal:
    enable: true
    addr: 127.0.0.1:3306
    user: canal
    password: canal_password
    databases:
      - test
    tables:
      - users
      - orders
      - products
```

### 2. 配置结构

**conf.proto**：
```protobuf
message Data {
  // ... 其他配置
  
  message Canal {
    bool enable = 1;
    string addr = 2;
    string user = 3;
    string password = 4;
    repeated string databases = 5;
    repeated string tables = 6;
  }
  
  Canal canal = 5;
}
```

## 测试验证

### 1. 单元测试

```go
// internal/data/canal/client_test.go
func TestCacheSyncHandler_OnRow(t *testing.T) {
    // 创建 Mock Redis
    redis := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    handler := NewCacheSyncHandler(redis, log.NewStdLogger(os.Stdout))
    
    // 创建测试事件
    event := &canal.RowsEvent{
        Table: &canal.Table{
            Schema: "test",
            Name:   "users",
        },
        Action: "update",
        Rows: [][]interface{}{
            {1, "old_username", "old_email"},
            {1, "new_username", "new_email"},
        },
    }
    
    // 执行处理
    err := handler.OnRow(event)
    assert.NoError(t, err)
    
    // 验证缓存已删除
    val, err := redis.Get(context.Background(), "user:1").Result()
    assert.Equal(t, redis.Nil, err)
    assert.Empty(t, val)
}
```

### 2. 集成测试

```go
// 测试步骤：
// 1. 启动 Canal Server
// 2. 启动应用服务
// 3. 更新数据库
// 4. 验证缓存已删除
```

## 部署步骤

### 1. 部署 Canal Server

```bash
# 下载并解压
wget https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.deployer-1.1.7.tar.gz
tar -xzf canal.deployer-1.1.7.tar.gz
cd canal.deployer-1.1.7

# 配置
vim conf/canal.properties

# 启动
./bin/startup.sh

# 查看日志
tail -f logs/canal/canal.log
```

### 2. 部署应用服务

```bash
# 构建应用
make build

# 启动应用
./bin/sre -conf ./configs

# 验证 Canal Client 连接
# 查看应用日志，确认 Canal Client 已启动
```

## 故障排查

### 1. Canal Server 连接失败

**检查项**：
- MySQL Binlog 是否启用
- Canal 账号权限是否正确
- 网络连接是否正常

**解决方案**：
```bash
# 检查 MySQL Binlog
mysql> SHOW VARIABLES LIKE 'log_bin';

# 检查 Canal 账号权限
mysql> SHOW GRANTS FOR 'canal'@'%';

# 测试连接
mysql -h 127.0.0.1 -u canal -p
```

### 2. 缓存未删除

**检查项**：
- Canal Client 是否正常运行
- 事件处理器是否注册
- Redis 连接是否正常

**解决方案**：
```bash
# 查看应用日志
tail -f logs/app.log | grep canal

# 测试 Redis 连接
redis-cli ping

# 手动触发缓存删除测试
```

## 相关文档

- [Canal 方案概述](./canal-cache-consistency-overview.md) - 方案概述
- [Canal 架构设计](./canal-cache-consistency-architecture.md) - 详细架构设计
- [Canal 配置说明](./canal-cache-consistency-config.md) - 配置参数说明
- [Canal 监控运维](./canal-cache-consistency-monitoring.md) - 监控和运维指南

