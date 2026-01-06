# PostgreSQL 索引与查询优化

## 目录
- [索引基础](#索引基础)
- [索引类型](#索引类型)
- [索引创建与管理](#索引创建与管理)
- [查询计划分析](#查询计划分析)
- [查询优化技巧](#查询优化技巧)
- [性能监控](#性能监控)

---

## 索引基础

### 什么是索引

索引是数据库中的数据结构，用于快速定位数据，类似于书籍的目录。

### 索引的作用

- **加速查询**：特别是 WHERE、JOIN、ORDER BY 操作
- **唯一性约束**：确保数据唯一性
- **加速排序**：ORDER BY 和 GROUP BY 操作

### 索引的代价

- **存储空间**：索引需要额外的磁盘空间
- **写入性能**：INSERT、UPDATE、DELETE 需要维护索引
- **维护成本**：需要定期维护（VACUUM、REINDEX）

---

## 索引类型

### B-tree 索引（默认）

最常用的索引类型，适用于大多数查询场景。

```sql
-- 创建 B-tree 索引
CREATE INDEX idx_users_email ON users(email);

-- 自动创建（主键和唯一约束）
CREATE TABLE users (
    id SERIAL PRIMARY KEY,        -- 自动创建 B-tree 索引
    email VARCHAR(100) UNIQUE     -- 自动创建 B-tree 索引
);
```

**适用场景**：
- 等值查询（=）
- 范围查询（<, >, <=, >=, BETWEEN）
- 排序（ORDER BY）
- LIKE 'prefix%'（前缀匹配）

### Hash 索引

适用于等值查询，不支持范围查询和排序。

```sql
-- 创建 Hash 索引
CREATE INDEX idx_users_username_hash ON users USING HASH (username);
```

**适用场景**：
- 频繁的等值查询
- 不支持范围查询和排序

**限制**：
- 不支持多列索引
- 不支持唯一索引
- 不支持部分索引

### GiST 索引（通用搜索树）

适用于复杂数据类型和全文搜索。

```sql
-- 创建 GiST 索引（需要扩展）
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE INDEX idx_products_price_gist ON products USING GIST (price);
```

**适用场景**：
- 几何数据类型
- 全文搜索
- 数组类型
- 范围类型

### GIN 索引（通用倒排索引）

适用于包含多个值的列，如数组、JSONB、全文搜索。

```sql
-- JSONB 索引
CREATE INDEX idx_documents_data_gin ON documents USING GIN (data);

-- 数组索引
CREATE INDEX idx_products_tags_gin ON products USING GIN (tags);

-- 全文搜索索引
CREATE INDEX idx_articles_content_gin ON articles USING GIN (to_tsvector('english', content));
```

**适用场景**：
- JSONB 查询
- 数组查询
- 全文搜索
- 多值列查询

### BRIN 索引（块范围索引）

适用于按顺序存储的大表，占用空间小。

```sql
-- 创建 BRIN 索引
CREATE INDEX idx_orders_created_at_brin ON orders USING BRIN (created_at);
```

**适用场景**：
- 按时间顺序插入的大表
- 范围查询
- 空间占用要求低

---

## 索引创建与管理

### 创建索引

```sql
-- 基本创建
CREATE INDEX idx_users_email ON users(email);

-- 唯一索引
CREATE UNIQUE INDEX idx_users_username_unique ON users(username);

-- 多列索引（复合索引）
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at);

-- 部分索引（条件索引）
CREATE INDEX idx_active_users ON users(email) WHERE status = 'active';

-- 表达式索引
CREATE INDEX idx_users_lower_email ON users(LOWER(email));

-- 包含列索引（PostgreSQL 11+）
CREATE INDEX idx_orders_user_date_inc 
ON orders(user_id, created_at) 
INCLUDE (amount, status);
```

### 索引命名规范

```sql
-- 推荐命名格式：idx_表名_列名_类型
CREATE INDEX idx_users_email_btree ON users(email);
CREATE INDEX idx_products_tags_gin ON products USING GIN (tags);
```

### 查看索引

```sql
-- 查看表的所有索引
\d users

-- 或使用 SQL
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'users';

-- 查看索引大小
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;
```

### 删除索引

```sql
-- 删除索引
DROP INDEX idx_users_email;

-- 如果不存在则忽略错误
DROP INDEX IF EXISTS idx_users_email;

-- 级联删除（删除依赖对象）
DROP INDEX idx_users_email CASCADE;
```

### 重建索引

```sql
-- 重建索引（释放空间，提高性能）
REINDEX INDEX idx_users_email;

-- 重建表的所有索引
REINDEX TABLE users;

-- 重建数据库的所有索引
REINDEX DATABASE mydb;
```

### 索引维护

```sql
-- 分析表（更新统计信息）
ANALYZE users;

-- 分析所有表
ANALYZE;

-- 查看索引使用情况
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,           -- 索引扫描次数
    idx_tup_read,       -- 读取的元组数
    idx_tup_fetch       -- 获取的元组数
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

---

## 查询计划分析

### EXPLAIN 命令

```sql
-- 基本 EXPLAIN
EXPLAIN SELECT * FROM users WHERE email = 'john@example.com';

-- 显示实际执行时间
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'john@example.com';

-- 详细输出
EXPLAIN (ANALYZE, BUFFERS, VERBOSE) 
SELECT * FROM users WHERE email = 'john@example.com';

-- 格式化输出
EXPLAIN (FORMAT JSON) SELECT * FROM users WHERE email = 'john@example.com';
```

### 执行计划解读

```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'john@example.com';

-- 输出示例：
-- Seq Scan on users  (cost=0.00..25.00 rows=1 width=64) (actual time=0.123..0.456 rows=1 loops=1)
--   Filter: (email = 'john@example.com'::text)
-- Planning Time: 0.123 ms
-- Execution Time: 0.567 ms
```

**关键指标**：
- **cost**: 预估成本（启动成本..总成本）
- **rows**: 预估返回行数
- **width**: 预估行宽度（字节）
- **actual time**: 实际执行时间
- **loops**: 循环次数

### 扫描类型

1. **Seq Scan（顺序扫描）**
   - 全表扫描
   - 适用于小表或大部分数据需要读取

2. **Index Scan（索引扫描）**
   - 使用索引查找
   - 适用于等值查询

3. **Index Only Scan（仅索引扫描）**
   - 只从索引读取数据
   - 性能最好

4. **Bitmap Index Scan + Bitmap Heap Scan**
   - 使用位图索引
   - 适用于多条件查询

### 连接类型

1. **Nested Loop（嵌套循环）**
   - 小表驱动大表
   - 适用于小数据集

2. **Hash Join（哈希连接）**
   - 构建哈希表
   - 适用于等值连接

3. **Merge Join（归并连接）**
   - 两个排序后的数据集合并
   - 适用于排序后的连接

---

## 查询优化技巧

### 1. 使用索引

```sql
-- ❌ 慢：全表扫描
SELECT * FROM users WHERE LOWER(email) = 'john@example.com';

-- ✅ 快：使用表达式索引
CREATE INDEX idx_users_lower_email ON users(LOWER(email));
SELECT * FROM users WHERE LOWER(email) = 'john@example.com';
```

### 2. 避免函数调用

```sql
-- ❌ 慢：无法使用索引
SELECT * FROM users WHERE EXTRACT(YEAR FROM created_at) = 2024;

-- ✅ 快：使用范围查询
SELECT * FROM users 
WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';
```

### 3. 使用 LIMIT

```sql
-- ✅ 使用 LIMIT 减少处理数据量
SELECT * FROM users ORDER BY created_at DESC LIMIT 10;

-- 如果只需要计数，使用 COUNT(*)
SELECT COUNT(*) FROM users WHERE status = 'active';
```

### 4. 避免 SELECT *

```sql
-- ❌ 慢：读取所有列
SELECT * FROM users WHERE id = 1;

-- ✅ 快：只选择需要的列
SELECT id, username, email FROM users WHERE id = 1;
```

### 5. 使用覆盖索引

```sql
-- 包含列索引（PostgreSQL 11+）
CREATE INDEX idx_orders_user_date_inc 
ON orders(user_id, created_at) 
INCLUDE (amount, status);

-- 查询可以直接从索引获取数据
SELECT user_id, created_at, amount, status 
FROM orders 
WHERE user_id = 1;
```

### 6. 优化 JOIN

```sql
-- ✅ 确保 JOIN 条件有索引
CREATE INDEX idx_orders_user_id ON orders(user_id);

SELECT u.username, o.amount
FROM users u
INNER JOIN orders o ON u.id = o.user_id
WHERE u.id = 1;
```

### 7. 使用 EXISTS 而不是 COUNT

```sql
-- ❌ 慢：需要计算所有匹配行
SELECT * FROM users 
WHERE (SELECT COUNT(*) FROM orders WHERE user_id = users.id) > 0;

-- ✅ 快：找到第一个匹配就返回
SELECT * FROM users 
WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id);
```

### 8. 优化子查询

```sql
-- ❌ 慢：相关子查询
SELECT * FROM users u
WHERE (SELECT COUNT(*) FROM orders o WHERE o.user_id = u.id) > 5;

-- ✅ 快：使用 JOIN
SELECT u.*
FROM users u
INNER JOIN (
    SELECT user_id, COUNT(*) as order_count
    FROM orders
    GROUP BY user_id
    HAVING COUNT(*) > 5
) o ON u.id = o.user_id;
```

### 9. 使用 UNION ALL 而不是 UNION

```sql
-- ❌ 慢：需要去重
SELECT id FROM users WHERE status = 'active'
UNION
SELECT id FROM users WHERE created_at > '2024-01-01';

-- ✅ 快：如果不需要去重
SELECT id FROM users WHERE status = 'active'
UNION ALL
SELECT id FROM users WHERE created_at > '2024-01-01';
```

### 10. 批量操作优化

```sql
-- ❌ 慢：逐条插入
INSERT INTO users (username, email) VALUES ('user1', 'user1@example.com');
INSERT INTO users (username, email) VALUES ('user2', 'user2@example.com');

-- ✅ 快：批量插入
INSERT INTO users (username, email) VALUES
    ('user1', 'user1@example.com'),
    ('user2', 'user2@example.com'),
    ('user3', 'user3@example.com');
```

---

## 性能监控

### 查看慢查询

```sql
-- 启用慢查询日志（在 postgresql.conf 中）
log_min_duration_statement = 1000  -- 记录超过1秒的查询

-- 查看当前慢查询
SELECT 
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';
```

### 使用 pg_stat_statements

```sql
-- 启用扩展
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- 查看最慢的查询
SELECT 
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- 重置统计信息
SELECT pg_stat_statements_reset();
```

### 表统计信息

```sql
-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 查看索引使用情况
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE idx_scan = 0  -- 未使用的索引
ORDER BY pg_relation_size(indexrelid) DESC;
```

### 连接和锁监控

```sql
-- 查看当前连接
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query
FROM pg_stat_activity
WHERE datname = 'mydb';

-- 查看锁
SELECT 
    locktype,
    relation::regclass,
    mode,
    granted
FROM pg_locks
WHERE relation = 'users'::regclass;

-- 查看阻塞的查询
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocking_locks.pid AS blocking_pid,
    blocked_activity.query AS blocked_query,
    blocking_activity.query AS blocking_query
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;
```

---

## 最佳实践

1. **索引策略**
   - 为主键和外键创建索引
   - 为频繁查询的列创建索引
   - 为 JOIN 条件创建索引
   - 定期检查未使用的索引

2. **查询优化**
   - 使用 EXPLAIN ANALYZE 分析查询
   - 避免在 WHERE 子句中使用函数
   - 使用合适的 JOIN 类型
   - 限制返回的数据量

3. **监控和维护**
   - 定期运行 ANALYZE
   - 监控慢查询
   - 检查索引使用情况
   - 定期重建索引

---

## 下一步学习

- [PostgreSQL 事务与并发控制](./postgresql-transactions.md)
- [PostgreSQL 管理与运维](./postgresql-admin.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)

---

*最后更新：2024年*

