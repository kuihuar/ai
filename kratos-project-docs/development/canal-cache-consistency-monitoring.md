# Canal 缓存一致性监控运维指南

## 监控指标

### 1. Canal Server 监控

**关键指标**：
- Canal Server 运行状态
- Binlog 订阅延迟
- 事件处理速度
- 错误率

**监控实现**：
```go
// internal/metrics/canal_metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    // Canal Server 连接状态
    CanalServerConnected = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "canal_server_connected",
            Help: "Canal Server connection status",
        },
        []string{"instance"},
    )
    
    // Binlog 延迟（秒）
    CanalBinlogDelay = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "canal_binlog_delay_seconds",
            Help: "Canal binlog delay in seconds",
        },
        []string{"instance", "database", "table"},
    )
    
    // 事件处理速度（事件/秒）
    CanalEventRate = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "canal_events_total",
            Help: "Total number of canal events processed",
        },
        []string{"instance", "database", "table", "event_type"},
    )
    
    // 事件处理错误数
    CanalEventErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "canal_event_errors_total",
            Help: "Total number of canal event errors",
        },
        []string{"instance", "database", "table", "error_type"},
    )
    
    // 缓存删除操作数
    CanalCacheDeletes = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "canal_cache_deletes_total",
            Help: "Total number of cache delete operations",
        },
        []string{"database", "table"},
    )
    
    // 缓存删除失败数
    CanalCacheDeleteErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "canal_cache_delete_errors_total",
            Help: "Total number of cache delete errors",
        },
        []string{"database", "table"},
    )
)

func init() {
    prometheus.MustRegister(
        CanalServerConnected,
        CanalBinlogDelay,
        CanalEventRate,
        CanalEventErrors,
        CanalCacheDeletes,
        CanalCacheDeleteErrors,
    )
}
```

### 2. 缓存同步监控

**关键指标**：
- 缓存删除成功率
- 缓存删除延迟
- 缓存命中率变化

**监控实现**：
```go
// 在缓存同步处理器中记录指标
func (h *CacheSyncHandler) handleUserEvent(ctx context.Context, event *CacheSyncEvent, e *canal.RowsEvent) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        metrics.CanalCacheDeleteLatency.WithLabelValues("test", "users").Observe(duration.Seconds())
    }()
    
    // ... 处理逻辑
    
    if err := h.redis.Del(ctx, keys...).Err(); err != nil {
        metrics.CanalCacheDeleteErrors.WithLabelValues("test", "users").Inc()
        return err
    }
    
    metrics.CanalCacheDeletes.WithLabelValues("test", "users").Add(float64(len(keys)))
    return nil
}
```

## 告警规则

### 1. Prometheus 告警规则

**alerts/canal.yml**：
```yaml
groups:
  - name: canal
    rules:
      # Canal Server 连接断开
      - alert: CanalServerDisconnected
        expr: canal_server_connected{instance="example"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Canal Server disconnected"
          description: "Canal Server instance {{ $labels.instance }} is disconnected"
      
      # Binlog 延迟过高
      - alert: CanalBinlogDelayHigh
        expr: canal_binlog_delay_seconds > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Canal binlog delay is high"
          description: "Canal binlog delay is {{ $value }}s for {{ $labels.database }}.{{ $labels.table }}"
      
      # 事件处理错误率过高
      - alert: CanalEventErrorRateHigh
        expr: rate(canal_event_errors_total[5m]) / rate(canal_events_total[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Canal event error rate is high"
          description: "Canal event error rate is {{ $value | humanizePercentage }}"
      
      # 缓存删除失败率过高
      - alert: CanalCacheDeleteErrorRateHigh
        expr: rate(canal_cache_delete_errors_total[5m]) / rate(canal_cache_deletes_total[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Canal cache delete error rate is high"
          description: "Canal cache delete error rate is {{ $value | humanizePercentage }}"
```

### 2. 告警通知

**配置告警通知渠道**：
- 邮件通知
- 钉钉/企业微信通知
- PagerDuty 集成

## 日志管理

### 1. 日志级别

**建议配置**：
- Canal Server：INFO 级别
- Canal Client：INFO 级别
- 缓存同步：DEBUG 级别（生产环境可调整为 INFO）

### 2. 日志格式

**结构化日志**：
```go
// 记录缓存删除日志
h.logger.WithContext(ctx).Infof(
    "cache deleted: database=%s, table=%s, event_type=%s, keys=%v, count=%d",
    event.Database,
    event.Table,
    event.EventType,
    keys,
    len(keys),
)
```

### 3. 日志聚合

**使用 ELK 或 Loki**：
- 集中收集 Canal 相关日志
- 设置日志保留策略（建议 30 天）
- 配置日志查询和告警

## 性能优化

### 1. 批量处理优化

**调整批量大小**：
```yaml
cache_sync:
  batch_size: 1000      # 根据实际情况调整
  batch_timeout: 5000   # 根据实际情况调整
```

**监控批量处理效果**：
- 批量大小分布
- 批量处理延迟
- 批量处理成功率

### 2. 异步处理优化

**调整异步工作线程数**：
```yaml
cache_sync:
  async: true
  async_workers: 50     # 根据 CPU 和 Redis 性能调整
```

**监控异步处理效果**：
- 异步队列长度
- 异步处理延迟
- 异步处理错误率

### 3. 缓存删除优化

**使用 Pipeline**：
```go
// 批量删除缓存（使用 Pipeline）
pipe := h.redis.Pipeline()
for _, key := range keys {
    pipe.Del(ctx, key)
}
_, err := pipe.Exec(ctx)
```

**监控 Pipeline 效果**：
- Pipeline 批量大小
- Pipeline 执行时间
- Pipeline 成功率

## 故障处理

### 1. Canal Server 故障

**症状**：
- Canal Client 连接失败
- 事件处理停止
- 缓存不再同步

**处理步骤**：
1. 检查 Canal Server 运行状态
2. 检查 MySQL 连接
3. 检查 Binlog 配置
4. 重启 Canal Server
5. 验证 Canal Client 重连

### 2. Binlog 延迟

**症状**：
- Binlog 延迟持续增长
- 缓存同步延迟

**处理步骤**：
1. 检查 MySQL 主从复制状态
2. 检查 Canal Server 处理速度
3. 检查网络延迟
4. 优化 Canal Server 配置
5. 考虑增加 Canal Server 实例

### 3. 缓存删除失败

**症状**：
- 缓存删除错误率上升
- 缓存不一致

**处理步骤**：
1. 检查 Redis 连接
2. 检查 Redis 内存使用
3. 检查网络延迟
4. 实现重试机制
5. 实现补偿删除

## 运维脚本

### 1. 健康检查脚本

**scripts/canal-health-check.sh**：
```bash
#!/bin/bash

# 检查 Canal Server 健康状态
check_canal_server() {
    curl -s http://127.0.0.1:11111/health > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✓ Canal Server is healthy"
        return 0
    else
        echo "✗ Canal Server is unhealthy"
        return 1
    fi
}

# 检查 Canal Client 连接
check_canal_client() {
    # 检查应用日志中是否有 Canal Client 错误
    tail -n 100 /var/log/app.log | grep -i "canal.*error" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✗ Canal Client has errors"
        return 1
    else
        echo "✓ Canal Client is healthy"
        return 0
    fi
}

# 执行检查
check_canal_server
check_canal_client
```

### 2. 数据一致性检查脚本

**scripts/check-cache-consistency.sh**：
```bash
#!/bin/bash

# 检查缓存和数据库一致性
check_consistency() {
    USER_ID=$1
    
    # 从数据库查询
    DB_DATA=$(mysql -h 127.0.0.1 -u root -p${MYSQL_PASSWORD} -D test -e "SELECT id, username, email FROM users WHERE id=$USER_ID" -N)
    
    # 从缓存查询
    CACHE_DATA=$(redis-cli GET "user:$USER_ID")
    
    # 比较数据
    if [ "$DB_DATA" == "$CACHE_DATA" ]; then
        echo "✓ User $USER_ID: Cache and DB are consistent"
    else
        echo "✗ User $USER_ID: Cache and DB are inconsistent"
        echo "  DB: $DB_DATA"
        echo "  Cache: $CACHE_DATA"
    fi
}

# 检查多个用户
for user_id in 1 2 3 4 5; do
    check_consistency $user_id
done
```

## 备份和恢复

### 1. Canal Server 配置备份

**备份脚本**：
```bash
#!/bin/bash

# 备份 Canal Server 配置
BACKUP_DIR="/backup/canal/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# 备份配置文件
cp -r /opt/canal/conf $BACKUP_DIR/

# 备份数据目录
cp -r /opt/canal/data $BACKUP_DIR/

echo "Canal Server configuration backed up to $BACKUP_DIR"
```

### 2. 恢复步骤

**恢复 Canal Server**：
1. 停止 Canal Server
2. 恢复配置文件
3. 恢复数据目录
4. 启动 Canal Server
5. 验证 Canal Client 连接

## 相关文档

- [Canal 方案概述](./canal-cache-consistency-overview.md) - 方案概述
- [Canal 架构设计](./canal-cache-consistency-architecture.md) - 详细架构设计
- [Canal 实现指南](./canal-cache-consistency-implementation.md) - 具体实现步骤
- [Canal 配置说明](./canal-cache-consistency-config.md) - 配置参数说明

