# Canal 缓存一致性配置说明

## Canal Server 配置

### 1. 基础配置

**conf/canal.properties**：
```properties
# Canal Server 端口
canal.port = 11111

# Canal Server 实例配置目录
canal.instance.global.spring.xml = classpath:spring/default-instance.xml

# Canal Server 数据存储目录
canal.instance.data.dir = ${canal.file.data.dir:../conf}
```

### 2. 实例配置

**conf/example/instance.properties**：
```properties
# MySQL 连接配置
canal.instance.master.address=127.0.0.1:3306
canal.instance.master.journal.name=
canal.instance.master.position=
canal.instance.master.timestamp=
canal.instance.master.gtid=

# MySQL 账号配置
canal.instance.dbUsername=canal
canal.instance.dbPassword=canal_password
canal.instance.connectionCharset=UTF-8

# 订阅配置
canal.instance.filter.regex=test\\..*
canal.instance.filter.black.regex=

# 批量处理配置
canal.instance.batchSize=1000
canal.instance.batchTimeout=1000
```

### 3. 过滤规则配置

**表过滤**：
```properties
# 只订阅 test 数据库下的 users 和 orders 表
canal.instance.filter.regex=test\\.(users|orders)

# 排除 test 数据库下的 logs 表
canal.instance.filter.black.regex=test\\.logs
```

**字段过滤**（需要自定义）：
- 默认 Canal 不支持字段过滤
- 需要在 Canal Client 中实现字段过滤逻辑

## Canal Client 配置

### 1. 应用配置

**configs/config.yaml**：
```yaml
data:
  # ... 其他配置
  
  canal:
    # 是否启用 Canal
    enable: true
    
    # Canal Server 地址
    server_addr: 127.0.0.1:11111
    
    # Canal 实例名称
    instance: example
    
    # MySQL 连接配置（用于 Canal Client 直连，可选）
    mysql:
      addr: 127.0.0.1:3306
      user: canal
      password: canal_password
      database: test
    
    # 订阅配置
    subscribe:
      # 订阅的数据库列表
      databases:
        - test
      
      # 订阅的表列表（为空则订阅所有表）
      tables:
        - users
        - orders
        - products
    
    # 缓存同步配置
    cache_sync:
      # 批量删除大小
      batch_size: 100
      
      # 批量删除超时时间（毫秒）
      batch_timeout: 1000
      
      # 重试配置
      retry:
        max_retries: 3
        retry_interval: 1s
      
      # 异步处理
      async: true
      async_workers: 10
```

### 2. 配置结构定义

**internal/conf/conf.proto**：
```protobuf
message Data {
  // ... 其他配置
  
  message Canal {
    // 是否启用 Canal
    bool enable = 1;
    
    // Canal Server 地址
    string server_addr = 2;
    
    // Canal 实例名称
    string instance = 3;
    
    // MySQL 连接配置（可选）
    message MySQL {
      string addr = 1;
      string user = 2;
      string password = 3;
      string database = 4;
    }
    MySQL mysql = 4;
    
    // 订阅配置
    message Subscribe {
      repeated string databases = 1;
      repeated string tables = 2;
    }
    Subscribe subscribe = 5;
    
    // 缓存同步配置
    message CacheSync {
      int32 batch_size = 1;
      int32 batch_timeout = 2;
      
      message Retry {
        int32 max_retries = 1;
        google.protobuf.Duration retry_interval = 2;
      }
      Retry retry = 3;
      
      bool async = 4;
      int32 async_workers = 5;
    }
    CacheSync cache_sync = 6;
  }
  
  Canal canal = 5;
}
```

## 缓存键映射配置

### 1. 表到缓存键映射

**internal/data/canal/cache_mapping.go**：
```go
package canal

// CacheKeyMapping 表到缓存键的映射配置
var CacheKeyMapping = map[string]CacheKeyRule{
    "users": {
        PrimaryKey: "id",
        Keys: []KeyTemplate{
            {Template: "user:%d", Fields: []string{"id"}},
            {Template: "user:username:%s", Fields: []string{"username"}},
            {Template: "user:email:%s", Fields: []string{"email"}},
        },
        RelatedKeys: []RelatedKeyRule{
            {
                Condition: "status",
                Keys: []string{
                    "user:list:*",
                    "user:list:active:*",
                },
            },
        },
    },
    "orders": {
        PrimaryKey: "id",
        Keys: []KeyTemplate{
            {Template: "order:%d", Fields: []string{"id"}},
        },
        RelatedKeys: []RelatedKeyRule{
            {
                Condition: "user_id",
                Keys: []string{
                    "user:%d:orders",
                    "user:%d:order:stats",
                },
            },
        },
    },
}

// CacheKeyRule 缓存键规则
type CacheKeyRule struct {
    PrimaryKey string
    Keys       []KeyTemplate
    RelatedKeys []RelatedKeyRule
}

// KeyTemplate 缓存键模板
type KeyTemplate struct {
    Template string
    Fields   []string
}

// RelatedKeyRule 关联缓存键规则
type RelatedKeyRule struct {
    Condition string
    Keys      []string
}
```

## 环境变量配置

### 1. 开发环境

**.env.development**：
```bash
CANAL_ENABLE=true
CANAL_SERVER_ADDR=127.0.0.1:11111
CANAL_INSTANCE=example
CANAL_MYSQL_ADDR=127.0.0.1:3306
CANAL_MYSQL_USER=canal
CANAL_MYSQL_PASSWORD=canal_password
```

### 2. 生产环境

**.env.production**：
```bash
CANAL_ENABLE=true
CANAL_SERVER_ADDR=canal-server:11111
CANAL_INSTANCE=production
CANAL_MYSQL_ADDR=mysql-master:3306
CANAL_MYSQL_USER=canal
CANAL_MYSQL_PASSWORD=${CANAL_PASSWORD}
```

## 配置验证

### 1. 配置检查脚本

**scripts/check-canal-config.sh**：
```bash
#!/bin/bash

# 检查 Canal Server 配置
echo "Checking Canal Server configuration..."

# 检查 MySQL 连接
mysql -h 127.0.0.1 -u canal -p${CANAL_PASSWORD} -e "SELECT 1" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ MySQL connection OK"
else
    echo "✗ MySQL connection failed"
    exit 1
fi

# 检查 Binlog 是否启用
mysql -h 127.0.0.1 -u canal -p${CANAL_PASSWORD} -e "SHOW VARIABLES LIKE 'log_bin'" | grep -q "ON"
if [ $? -eq 0 ]; then
    echo "✓ Binlog enabled"
else
    echo "✗ Binlog not enabled"
    exit 1
fi

# 检查 Canal Server 是否运行
curl -s http://127.0.0.1:11111/ > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ Canal Server running"
else
    echo "✗ Canal Server not running"
    exit 1
fi

echo "All checks passed!"
```

### 2. 配置测试

**测试 Canal Client 连接**：
```go
// internal/data/canal/manager_test.go
func TestCanalManager_Connect(t *testing.T) {
    cfg := &conf.Data_Canal{
        Enable:     true,
        ServerAddr: "127.0.0.1:11111",
        Instance:   "example",
    }
    
    manager, err := NewManager(cfg, nil, log.NewStdLogger(os.Stdout))
    assert.NoError(t, err)
    
    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err = manager.Start(ctx)
    assert.NoError(t, err)
}
```

## 配置最佳实践

### 1. 生产环境配置

**建议配置**：
- 启用 Canal 集群模式
- 配置 Canal Server 高可用
- 使用独立的 MySQL 从库
- 配置监控和告警

### 2. 性能优化配置

**批量处理**：
```yaml
cache_sync:
  batch_size: 1000      # 增大批量大小
  batch_timeout: 5000   # 增大超时时间
  async: true
  async_workers: 50     # 增加异步工作线程
```

### 3. 安全配置

**权限控制**：
- Canal 账号使用最小权限（只读 + REPLICATION）
- Redis 账号使用独立账号，限制删除操作
- 配置网络访问控制（防火墙规则）

## 相关文档

- [Canal 方案概述](./canal-cache-consistency-overview.md) - 方案概述
- [Canal 架构设计](./canal-cache-consistency-architecture.md) - 详细架构设计
- [Canal 实现指南](./canal-cache-consistency-implementation.md) - 具体实现步骤
- [Canal 监控运维](./canal-cache-consistency-monitoring.md) - 监控和运维指南

