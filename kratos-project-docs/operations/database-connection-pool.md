# 数据库连接池配置说明

## 概述

数据库连接池配置已添加到配置文件中，可以通过配置文件灵活调整连接池参数，无需修改代码。

## 配置位置

### 配置文件

**位置**：`configs/config.yaml`

```yaml
data:
  database:
    enable: true
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/test?timeout=15s&charset=utf8mb4&parseTime=True
    decrypt_key: "your_decrypt_key"
    pool:
      max_open_conns: 25       # 最大打开连接数
      max_idle_conns: 10       # 最大空闲连接数
      conn_max_lifetime: 5m    # 连接最大生存时间
      conn_max_idle_time: 10m  # 空闲连接最大生存时间
```

## 配置参数说明

### 1. max_open_conns（最大打开连接数）

**类型**：`int32`  
**默认值**：`25`  
**说明**：数据库连接池中同时打开的最大连接数。

**建议值**：
- **开发环境**：10-25
- **测试环境**：25-50
- **生产环境**：根据数据库服务器配置和并发请求量调整
  - 低并发：25-50
  - 中并发：50-100
  - 高并发：100-200

**注意事项**：
- 不要超过数据库服务器的 `max_connections` 配置
- 考虑应用实例数量（如果有多个实例，总连接数 = max_open_conns × 实例数）

### 2. max_idle_conns（最大空闲连接数）

**类型**：`int32`  
**默认值**：`10`  
**说明**：连接池中保持的最大空闲连接数。

**建议值**：
- 通常设置为 `max_open_conns` 的 40-50%
- 例如：`max_open_conns = 25`，则 `max_idle_conns = 10`

**注意事项**：
- 空闲连接会占用数据库服务器资源
- 设置过大会浪费资源，设置过小会导致频繁创建连接

### 3. conn_max_lifetime（连接最大生存时间）

**类型**：`duration`（如 `5m`、`1h`）  
**默认值**：`5m`  
**说明**：连接的最大生存时间，超过此时间的连接会被关闭并重新创建。

**建议值**：
- **开发/测试环境**：5-10 分钟
- **生产环境**：根据数据库服务器配置调整
  - MySQL：通常设置为小于 `wait_timeout`（默认 8 小时）
  - 建议：1-4 小时

**注意事项**：
- 防止使用过期的连接（数据库服务器可能关闭空闲连接）
- 不要设置过长，避免使用有问题的连接

### 4. conn_max_idle_time（空闲连接最大生存时间）

**类型**：`duration`（如 `10m`、`30m`）  
**默认值**：`10m`  
**说明**：空闲连接的最大生存时间，超过此时间的空闲连接会被关闭。

**建议值**：
- **开发/测试环境**：10-30 分钟
- **生产环境**：30 分钟 - 1 小时

**注意事项**：
- 帮助回收长时间未使用的连接
- 减少数据库服务器的连接数

## 配置示例

### 开发环境配置

```yaml
data:
  database:
    pool:
      max_open_conns: 10
      max_idle_conns: 5
      conn_max_lifetime: 5m
      conn_max_idle_time: 10m
```

### 测试环境配置

```yaml
data:
  database:
    pool:
      max_open_conns: 25
      max_idle_conns: 10
      conn_max_lifetime: 10m
      conn_max_idle_time: 30m
```

### 生产环境配置（低并发）

```yaml
data:
  database:
    pool:
      max_open_conns: 25
      max_idle_conns: 10
      conn_max_lifetime: 1h
      conn_max_idle_time: 30m
```

### 生产环境配置（高并发）

```yaml
data:
  database:
    pool:
      max_open_conns: 100
      max_idle_conns: 50
      conn_max_lifetime: 2h
      conn_max_idle_time: 1h
```

## 配置加载逻辑

### 默认值

如果配置文件中没有设置连接池参数，将使用以下默认值：

```go
maxOpenConns := 25
maxIdleConns := 10
connMaxLifetime := 5 * time.Minute
connMaxIdleTime := 10 * time.Minute
```

### 配置优先级

1. **配置文件中的值**（最高优先级）
2. **默认值**（如果配置文件中没有设置）

### 配置验证

- `max_open_conns > 0`：如果设置为 0 或负数，使用默认值 25
- `max_idle_conns > 0`：如果设置为 0 或负数，使用默认值 10
- `conn_max_lifetime`：如果未设置，使用默认值 5 分钟
- `conn_max_idle_time`：如果未设置，使用默认值 10 分钟

## 实现位置

### 1. Proto 定义

**文件**：`internal/conf/conf.proto`

```protobuf
message Data {
  message Database {
    // ... 其他字段 ...
    Pool pool = 5;
  }
  
  message Pool {
    int32 max_open_conns = 1;
    int32 max_idle_conns = 2;
    google.protobuf.Duration conn_max_lifetime = 3;
    google.protobuf.Duration conn_max_idle_time = 4;
  }
}
```

### 2. 配置加载

**文件**：`internal/config/kratos.go`

```go
// 连接池配置
if v.IsSet("data.database.pool") {
    pool := &conf.Data_Pool{}
    if v.IsSet("data.database.pool.max_open_conns") {
        pool.MaxOpenConns = int32(v.GetInt("data.database.pool.max_open_conns"))
    }
    // ... 其他字段 ...
    database.Pool = pool
}
```

### 3. 连接池应用

**文件**：`internal/data/data.go`

```go
// 配置连接池
sqlDB, err := db.DB()
if err != nil {
    return nil, nil
}

// 使用配置的连接池参数
maxOpenConns := 25  // 默认值
if c.Database.Pool != nil && c.Database.Pool.MaxOpenConns > 0 {
    maxOpenConns = int(c.Database.Pool.MaxOpenConns)
}

sqlDB.SetMaxOpenConns(maxOpenConns)
sqlDB.SetMaxIdleConns(maxIdleConns)
sqlDB.SetConnMaxLifetime(connMaxLifetime)
sqlDB.SetConnMaxIdleTime(connMaxIdleTime)
```

## 监控和调优

### 1. 查看连接池状态

连接池状态可以通过健康检查端点查看：

```bash
curl http://localhost:8000/health
```

响应中包含连接池统计信息：

```json
{
  "status": "SERVING",
  "details": {
    "database": "ok",
    "database_open_conns": "5",
    "database_idle_conns": "3"
  }
}
```

### 2. 连接池指标

可以通过以下方式监控连接池：

- **健康检查**：查看 `database_open_conns` 和 `database_idle_conns`
- **日志**：启动时会输出连接池配置信息
- **数据库监控**：查看数据库服务器的连接数

### 3. 调优建议

#### 如果看到 "too many connections" 错误

**原因**：连接数超过数据库服务器限制

**解决方案**：
1. 减少 `max_open_conns`
2. 增加数据库服务器的 `max_connections`
3. 减少应用实例数量

#### 如果看到连接创建频繁

**原因**：`conn_max_lifetime` 或 `conn_max_idle_time` 设置过短

**解决方案**：
1. 增加 `conn_max_lifetime`
2. 增加 `conn_max_idle_time`
3. 增加 `max_idle_conns`

#### 如果看到连接等待时间长

**原因**：`max_open_conns` 设置过小

**解决方案**：
1. 增加 `max_open_conns`
2. 检查是否有慢查询导致连接占用时间过长

## 最佳实践

### 1. 根据环境调整

- **开发环境**：使用较小的连接数，快速发现问题
- **测试环境**：使用中等连接数，模拟生产环境
- **生产环境**：根据实际负载调整

### 2. 考虑应用实例数

如果有多个应用实例，总连接数 = `max_open_conns` × 实例数

例如：
- 3 个应用实例
- 每个实例 `max_open_conns = 25`
- 总连接数 = 75

确保数据库服务器的 `max_connections` > 75

### 3. 监控和告警

建议监控：
- 当前打开的连接数
- 当前空闲的连接数
- 连接等待时间
- 连接创建频率

### 4. 定期检查

定期检查连接池配置是否合理：
- 查看健康检查中的连接池统计
- 查看数据库服务器的连接数
- 查看应用日志中的连接池配置信息

## 故障排查

### 问题 1：连接数不足

**症状**：
- 请求超时
- 日志中出现 "connection pool exhausted"

**排查**：
1. 检查 `max_open_conns` 配置
2. 检查应用实例数量
3. 检查是否有慢查询

**解决**：
1. 增加 `max_open_conns`
2. 优化慢查询
3. 增加数据库服务器连接数限制

### 问题 2：连接泄漏

**症状**：
- 连接数持续增长
- 数据库服务器连接数接近上限

**排查**：
1. 检查是否有未关闭的连接
2. 检查事务是否正确提交/回滚
3. 检查是否有长时间运行的查询

**解决**：
1. 确保所有数据库操作都使用 context
2. 确保事务正确关闭
3. 设置合理的超时时间

### 问题 3：连接创建频繁

**症状**：
- 日志中频繁出现连接创建信息
- 性能下降

**排查**：
1. 检查 `conn_max_lifetime` 和 `conn_max_idle_time` 配置
2. 检查数据库服务器的 `wait_timeout` 配置

**解决**：
1. 增加 `conn_max_lifetime`
2. 增加 `conn_max_idle_time`
3. 调整数据库服务器超时配置

## 参考资源

- [Go database/sql 文档](https://pkg.go.dev/database/sql#DB.SetMaxOpenConns)
- [MySQL 连接管理](https://dev.mysql.com/doc/refman/8.0/en/connection-management.html)
- [连接池最佳实践](https://www.alexedwards.net/blog/configuring-sqldb)

