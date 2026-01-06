# Vitess 分库分表方案

## 概述

Vitess 是一个用于 MySQL 分库分表的数据库中间件，由 YouTube 开发并开源。它提供了透明的分库分表能力，让应用程序可以像使用单个数据库一样使用分片数据库集群。

## 为什么需要分库分表

### 1. 单库单表的瓶颈

当数据量增长到一定程度时，单库单表会遇到以下问题：

- **存储瓶颈**：单表数据量过大，查询性能下降
- **写入瓶颈**：单库写入能力有限，无法满足高并发写入
- **连接数瓶颈**：单库连接数有限，无法支持大量并发连接
- **备份恢复困难**：大表备份和恢复时间过长

### 2. 分库分表的优势

- **水平扩展**：通过增加分片数量，线性提升存储和性能
- **高可用**：单个分片故障不影响其他分片
- **负载均衡**：将请求分散到多个分片，提升整体吞吐量
- **灵活扩容**：可以根据业务增长动态增加分片

## Vitess 架构

### 1. 核心组件

```
┌─────────────┐
│ Application │
└──────┬──────┘
       │
       │ gRPC/MySQL Protocol
       │
┌──────▼─────────────────────────────────────┐
│          Vitess Gateway (VTGate)            │
│  - 路由查询到正确的分片                     │
│  - 聚合跨分片查询结果                       │
│  - 管理事务                                 │
└──────┬─────────────────────────────────────┘
       │
       ├──────────┬──────────┬──────────┐
       │          │          │          │
┌──────▼──┐  ┌────▼───┐  ┌───▼────┐  ┌──▼─────┐
│ VTTablet│  │VTTablet│  │VTTablet│  │VTTablet│
│ Shard 1 │  │Shard 2 │  │Shard 3 │  │Shard N │
└────┬────┘  └───┬────┘  └───┬────┘  └───┬────┘
     │          │          │          │
┌────▼───┐  ┌───▼────┐  ┌───▼────┐  ┌──▼─────┐
│ MySQL  │  │ MySQL  │  │ MySQL  │  │ MySQL  │
│ Shard 1│  │Shard 2 │  │Shard 3 │  │Shard N │
└────────┘  └────────┘  └────────┘  └────────┘
```

#### VTGate（网关）

- **功能**：应用程序的入口点，负责路由和查询聚合
- **协议**：支持 MySQL 协议和 gRPC 协议
- **路由**：根据分片键（Sharding Key）将查询路由到正确的分片

#### VTTablet（分片代理）

- **功能**：每个分片的代理，管理 MySQL 连接池
- **职责**：
  - 连接池管理
  - 查询重写（如跨分片查询）
  - 健康检查
  - 主从切换

#### Topology Service（拓扑服务）

- **功能**：存储集群元数据（分片信息、路由规则等）
- **实现**：通常使用 etcd 或 Consul

### 2. 分片策略

#### 范围分片（Range Sharding）

按数据范围划分分片：

```
Shard 1: user_id 0-1000000
Shard 2: user_id 1000001-2000000
Shard 3: user_id 2000001-3000000
...
```

**优点**：
- 实现简单
- 范围查询效率高

**缺点**：
- 数据分布可能不均匀
- 热点数据问题

#### 哈希分片（Hash Sharding）

按哈希值划分分片：

```
Shard = hash(user_id) % num_shards
```

**优点**：
- 数据分布均匀
- 避免热点问题

**缺点**：
- 范围查询需要跨分片
- 扩容需要数据迁移

#### 目录分片（Directory Sharding）

使用查找表（Lookup Table）存储分片映射：

```
user_id -> shard_id 映射表
```

**优点**：
- 灵活，可以动态调整
- 支持复杂的分片规则

**缺点**：
- 需要额外的查找表
- 性能开销

## 分片键（Sharding Key）选择

### 1. 选择原则

- **高基数**：分片键的值应该分布均匀，避免热点
- **业务相关**：选择业务中经常用于查询的字段
- **不可变**：分片键一旦确定，不应该修改

### 2. 常见分片键

#### 用户ID（User ID）

```sql
-- 按用户ID分片
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    user_id BIGINT,  -- 分片键
    order_no VARCHAR(64),
    amount DECIMAL(10,2),
    ...
) ENGINE=InnoDB;
```

**适用场景**：
- 用户相关数据（订单、支付等）
- 查询通常按用户ID过滤

#### 订单号（Order No）

```sql
-- 按订单号分片
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    order_no VARCHAR(64),  -- 分片键
    user_id BIGINT,
    amount DECIMAL(10,2),
    ...
) ENGINE=InnoDB;
```

**适用场景**：
- 订单系统
- 查询通常按订单号查询

#### 时间（Time-based）

```sql
-- 按时间分片（按月）
CREATE TABLE logs (
    id BIGINT PRIMARY KEY,
    created_at TIMESTAMP,  -- 分片键
    level VARCHAR(16),
    message TEXT,
    ...
) ENGINE=InnoDB;
```

**适用场景**：
- 日志系统
- 时序数据
- 历史数据归档

### 3. 复合分片键

对于复杂场景，可以使用复合分片键：

```sql
-- 按 (user_id, order_date) 分片
-- 先按 user_id 分片，再按 order_date 分片
```

## 查询路由

### 1. 单分片查询（Single Shard Query）

查询条件包含分片键，可以直接路由到单个分片：

```sql
-- ✅ 包含分片键，路由到单个分片
SELECT * FROM orders WHERE user_id = 12345;

-- ✅ 包含分片键，路由到单个分片
SELECT * FROM orders WHERE user_id IN (12345, 67890);
```

### 2. 跨分片查询（Cross Shard Query）

查询条件不包含分片键，需要查询所有分片并聚合：

```sql
-- ❌ 不包含分片键，需要查询所有分片
SELECT * FROM orders WHERE status = 'pending';

-- ❌ 不包含分片键，需要查询所有分片
SELECT COUNT(*) FROM orders WHERE created_at > '2024-01-01';
```

**性能影响**：
- 需要查询所有分片
- 需要聚合结果
- 性能较差，应尽量避免

### 3. 范围查询（Range Query）

按分片键范围查询，可能涉及多个分片：

```sql
-- 可能涉及多个分片
SELECT * FROM orders 
WHERE user_id BETWEEN 1000000 AND 2000000;
```

## 事务处理

### 1. 单分片事务

事务只涉及单个分片，性能最好：

```sql
BEGIN;
INSERT INTO orders (user_id, order_no, amount) VALUES (12345, 'ORD001', 100.00);
UPDATE users SET balance = balance - 100.00 WHERE id = 12345;
COMMIT;
```

### 2. 分布式事务

事务涉及多个分片，需要使用两阶段提交（2PC）：

```sql
-- 涉及多个分片的事务
BEGIN;
-- 分片1：创建订单
INSERT INTO orders (user_id, order_no, amount) VALUES (12345, 'ORD001', 100.00);
-- 分片2：扣减库存
UPDATE inventory SET stock = stock - 1 WHERE product_id = 999;
COMMIT;  -- 需要2PC
```

**性能影响**：
- 需要协调多个分片
- 性能较差，应尽量避免
- 建议使用 Saga 模式替代

## 数据迁移和扩容

### 1. 垂直分片（Vertical Sharding）

按表拆分，将不同表放到不同分片：

```
Shard 1: users, user_profiles
Shard 2: orders, order_items
Shard 3: products, inventory
```

### 2. 水平分片（Horizontal Sharding）

按行拆分，将同一表的数据分布到多个分片：

```
Shard 1: orders (user_id 0-1000000)
Shard 2: orders (user_id 1000001-2000000)
Shard 3: orders (user_id 2000001-3000000)
```

### 3. 扩容流程

#### 准备阶段

1. 创建新的分片（Shard N+1）
2. 配置路由规则
3. 准备数据迁移工具

#### 迁移阶段

1. **双写阶段**：同时写入旧分片和新分片
2. **数据迁移**：将历史数据迁移到新分片
3. **数据校验**：验证数据一致性
4. **切换阶段**：切换到新分片，停止写入旧分片
5. **清理阶段**：清理旧分片数据

#### 使用 Vitess 的 Resharding

Vitess 提供了自动化的 Resharding 工具：

```bash
# 1. 创建新的分片
vtctlclient Reshard -source_shards 0 -target_shards -80,80- commerce

# 2. 开始迁移
vtctlclient MigrateServedTypes commerce/0 rdonly
vtctlclient MigrateServedTypes commerce/0 replica
vtctlclient MigrateServedTypes commerce/0 master
```

## 在项目中的集成方案

### 1. 连接配置

#### 使用 MySQL 协议

```go
// configs/config.yaml
data:
  database:
    driver: "mysql"
    source: "user:password@tcp(vtgate-host:15306)/commerce?parseTime=true"
    # Vitess VTGate 默认端口：15306 (MySQL 协议)
```

#### 使用 gRPC 协议（推荐）

```go
// configs/config.yaml
data:
  database:
    driver: "vitess"
    source: "vtgate-host:15999"  # gRPC 端口
    keyspace: "commerce"
    sharding_key: "user_id"  # 分片键
```

### 2. 代码示例

#### 查询时指定分片键

```go
// internal/data/repo_order.go
func (r *orderRepo) FindByUserID(ctx context.Context, userID int64) ([]*biz.Order, error) {
    // Vitess 会根据 user_id 自动路由到正确的分片
    orders, err := r.data.ent.Order.
        Query().
        Where(order.UserID(userID)).
        All(ctx)
    
    return orders, err
}
```

#### 插入时包含分片键

```go
func (r *orderRepo) Save(ctx context.Context, order *biz.Order) (*biz.Order, error) {
    // 确保包含分片键 user_id
    create := r.data.ent.Order.Create().
        SetUserID(order.UserID).  // 分片键
        SetOrderNo(order.OrderNo).
        SetAmount(order.Amount)
    
    return create.Save(ctx)
}
```

### 3. 跨分片查询的处理

#### 避免跨分片查询

```go
// ❌ 避免：不包含分片键的查询
func (r *orderRepo) FindByStatus(ctx context.Context, status string) ([]*biz.Order, error) {
    // 这会查询所有分片，性能很差
    orders, err := r.data.ent.Order.
        Query().
        Where(order.Status(status)).
        All(ctx)
    return orders, err
}

// ✅ 推荐：使用分片键 + 其他条件
func (r *orderRepo) FindByUserIDAndStatus(ctx context.Context, userID int64, status string) ([]*biz.Order, error) {
    // 只查询单个分片
    orders, err := r.data.ent.Order.
        Query().
        Where(
            order.UserID(userID),  // 分片键
            order.Status(status),
        ).
        All(ctx)
    return orders, err
}
```

#### 使用聚合查询

```go
// 如果必须跨分片查询，使用聚合函数
func (r *orderRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
    // Vitess 会在各个分片执行 COUNT，然后聚合结果
    count, err := r.data.ent.Order.
        Query().
        Where(order.Status(status)).
        Count(ctx)
    return int64(count), err
}
```

## 最佳实践

### 1. 分片键设计

- ✅ **选择高基数字段**：如 user_id、order_id
- ✅ **避免热点**：不要使用单调递增的 ID（如自增主键）
- ✅ **考虑查询模式**：选择经常用于查询的字段
- ❌ **避免频繁修改**：分片键一旦确定，不应修改

### 2. 查询优化

- ✅ **优先使用分片键**：查询条件包含分片键
- ✅ **避免跨分片查询**：尽量减少跨分片查询
- ✅ **使用索引**：在分片键上建立索引
- ❌ **避免全表扫描**：不使用分片键的查询会扫描所有分片

### 3. 事务设计

- ✅ **单分片事务**：尽量让事务只涉及单个分片
- ✅ **使用 Saga 模式**：跨分片事务使用 Saga 模式
- ❌ **避免分布式事务**：2PC 性能差，应尽量避免

### 4. 数据一致性

- ✅ **最终一致性**：接受最终一致性，使用补偿机制
- ✅ **幂等性**：确保操作幂等，支持重试
- ✅ **数据校验**：定期校验数据一致性

### 5. 监控和运维

- ✅ **分片监控**：监控每个分片的性能
- ✅ **查询分析**：分析慢查询，优化路由
- ✅ **容量规划**：根据数据增长规划分片数量
- ✅ **备份策略**：每个分片独立备份

## 与项目现有架构的集成

### 1. 与 Ent ORM 集成

Ent 支持 Vitess，可以通过配置连接 Vitess：

```go
// internal/data/data.go
func NewEntClient(db *gorm.DB, c *conf.Data, logger log.Logger) (*ent.Client, error) {
    drv, err := sql.Open("mysql", c.Database.Source)
    if err != nil {
        return nil, err
    }
    
    // Ent 会自动识别 Vitess 连接
    return ent.NewClient(ent.Driver(drv)), nil
}
```

### 2. 与 Saga 模式结合

Vitess 的分布式事务性能较差，建议使用 Saga 模式：

```go
// internal/biz/order_saga.go
type OrderCreateSaga struct {
    // 订单创建涉及多个分片时，使用 Saga 模式
    // 而不是 Vitess 的分布式事务
}
```

### 3. 与 Outbox 模式结合

使用 Outbox 模式确保跨分片事件的一致性：

```go
// internal/data/repo_order.go
func (r *orderRepo) SaveWithEvent(ctx context.Context, order *biz.Order, event biz.DomainEvent) (*biz.Order, error) {
    // 在同一个分片的事务中保存订单和事件
    // 确保原子性
}
```

## 性能对比

### 单库 vs 分片

| 指标 | 单库 | 分片（4个） |
|------|------|------------|
| 写入QPS | 10,000 | 40,000 |
| 查询QPS | 20,000 | 80,000 |
| 存储容量 | 1TB | 4TB |
| 单分片查询延迟 | 10ms | 10ms |
| 跨分片查询延迟 | 10ms | 40ms+ |

### 查询性能

- **单分片查询**：性能与单库相当
- **跨分片查询**：性能下降，延迟增加
- **聚合查询**：需要聚合多个分片结果，性能较差

## 适用场景

### ✅ 适合使用 Vitess

- 数据量巨大（TB 级别）
- 高并发写入（10万+ QPS）
- 读多写少场景
- 可以接受最终一致性

### ❌ 不适合使用 Vitess

- 数据量小（GB 级别）
- 需要强一致性
- 频繁的跨分片查询
- 复杂的分布式事务

## 参考资源

- [Vitess 官方文档](https://vitess.io/docs/)
- [Vitess GitHub](https://github.com/vitessio/vitess)
- [Vitess 架构设计](https://vitess.io/docs/overview/architecture/)
- [分片键选择指南](https://vitess.io/docs/user-guides/sharding/)

## 总结

Vitess 是一个强大的分库分表解决方案，适合大规模数据场景。在使用时需要注意：

1. **合理选择分片键**：选择高基数、业务相关的字段
2. **优化查询模式**：尽量使用分片键，避免跨分片查询
3. **事务设计**：单分片事务优先，跨分片使用 Saga 模式
4. **监控运维**：建立完善的监控和运维体系

通过合理的设计和使用，Vitess 可以帮助系统实现水平扩展，支撑大规模业务。

