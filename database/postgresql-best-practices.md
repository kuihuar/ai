# PostgreSQL 最佳实践

## 目录
- [数据库设计](#数据库设计)
- [性能优化](#性能优化)
- [安全实践](#安全实践)
- [开发规范](#开发规范)
- [运维实践](#运维实践)

---

## 数据库设计

### 命名规范

```sql
-- ✅ 表名：小写，使用下划线，复数形式
CREATE TABLE users (...);
CREATE TABLE order_items (...);

-- ✅ 列名：小写，使用下划线
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(50),
    created_at TIMESTAMP
);

-- ✅ 索引名：idx_表名_列名
CREATE INDEX idx_users_email ON users(email);

-- ✅ 约束名：表名_列名_约束类型
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
```

### 数据类型选择

```sql
-- ✅ 使用合适的数据类型
CREATE TABLE products (
    id SERIAL PRIMARY KEY,           -- 自增ID
    name VARCHAR(100) NOT NULL,       -- 有长度限制
    description TEXT,                 -- 无长度限制
    price DECIMAL(10, 2),            -- 精确数值
    stock INTEGER CHECK (stock >= 0), -- 非负整数
    is_active BOOLEAN DEFAULT TRUE,   -- 布尔值
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  -- 带时区时间戳
);

-- ❌ 避免使用过大的数据类型
-- SMALLINT 足够时不要用 INTEGER
-- VARCHAR(50) 足够时不要用 TEXT
```

### 约束使用

```sql
-- ✅ 使用约束保证数据完整性
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 索引策略

```sql
-- ✅ 为主键和外键创建索引（自动创建）
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,        -- 自动创建索引
    user_id INTEGER REFERENCES users(id)  -- 需要手动创建索引
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

-- ✅ 为频繁查询的列创建索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- ✅ 使用复合索引
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at);

-- ✅ 使用部分索引
CREATE INDEX idx_active_users ON users(email) WHERE status = 'active';
```

---

## 性能优化

### 查询优化

```sql
-- ✅ 使用 EXPLAIN ANALYZE 分析查询
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'john@example.com';

-- ✅ 避免 SELECT *
SELECT id, username, email FROM users WHERE id = 1;

-- ✅ 使用 LIMIT
SELECT * FROM users ORDER BY created_at DESC LIMIT 10;

-- ✅ 使用索引列查询
SELECT * FROM users WHERE email = 'john@example.com';  -- email 有索引

-- ❌ 避免在 WHERE 子句中使用函数
-- SELECT * FROM users WHERE LOWER(email) = 'john@example.com';
-- ✅ 使用表达式索引
CREATE INDEX idx_users_lower_email ON users(LOWER(email));
```

### 连接优化

```sql
-- ✅ 确保 JOIN 条件有索引
CREATE INDEX idx_orders_user_id ON orders(user_id);

SELECT u.username, o.amount
FROM users u
INNER JOIN orders o ON u.id = o.user_id
WHERE u.id = 1;

-- ✅ 使用 EXISTS 而不是 COUNT
SELECT * FROM users u
WHERE EXISTS (SELECT 1 FROM orders o WHERE o.user_id = u.id);

-- ❌ 避免相关子查询
-- SELECT * FROM users u
-- WHERE (SELECT COUNT(*) FROM orders o WHERE o.user_id = u.id) > 5;
```

### 批量操作

```sql
-- ✅ 批量插入
INSERT INTO users (username, email) VALUES
    ('user1', 'user1@example.com'),
    ('user2', 'user2@example.com'),
    ('user3', 'user3@example.com');

-- ✅ 使用 COPY 导入大量数据
COPY users (username, email) FROM '/path/to/data.csv' WITH CSV HEADER;

-- ✅ 批量更新
UPDATE users 
SET status = 'active'
WHERE id IN (1, 2, 3, 4, 5);
```

### 事务优化

```sql
-- ✅ 保持事务尽可能短
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;  -- 立即提交

-- ❌ 避免长时间事务
-- BEGIN;
-- UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- -- ... 执行其他耗时操作 ...
-- COMMIT;
```

---

## 安全实践

### 用户和权限

```sql
-- ✅ 创建专用用户，不要使用 postgres 用户
CREATE USER app_user WITH PASSWORD 'strong_password';

-- ✅ 授予最小权限
GRANT CONNECT ON DATABASE mydb TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE ON TABLE users TO app_user;
-- 不授予 DELETE 权限，除非必要

-- ✅ 使用角色管理权限
CREATE ROLE read_only;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO read_only;
GRANT read_only TO app_user;
```

### 密码安全

```sql
-- ✅ 使用强密码
CREATE USER app_user WITH PASSWORD 'ComplexP@ssw0rd!123';

-- ✅ 定期更换密码
ALTER USER app_user WITH PASSWORD 'NewComplexP@ssw0rd!123';

-- ✅ 设置密码过期
ALTER USER app_user VALID UNTIL '2025-12-31';
```

### 连接安全

```conf
# pg_hba.conf
# ✅ 使用 SSL 连接
hostssl    all    all    0.0.0.0/0    md5

# ✅ 限制 IP 访问
host    all    all    192.168.1.0/24    md5

# ❌ 避免使用 trust
# local    all    all    trust
```

### SQL 注入防护

```sql
-- ✅ 使用参数化查询（在应用层）
-- Python 示例
-- cursor.execute("SELECT * FROM users WHERE email = %s", (email,))

-- ❌ 避免字符串拼接
-- cursor.execute("SELECT * FROM users WHERE email = '" + email + "'")
```

---

## 开发规范

### 函数和存储过程

```sql
-- ✅ 使用有意义的函数名
CREATE OR REPLACE FUNCTION get_user_by_email(user_email TEXT)
RETURNS users AS $$
DECLARE
    user_record users%ROWTYPE;
BEGIN
    SELECT * INTO user_record
    FROM users
    WHERE email = user_email;
    
    RETURN user_record;
END;
$$ LANGUAGE plpgsql;

-- ✅ 添加注释
COMMENT ON FUNCTION get_user_by_email IS '根据邮箱获取用户信息';
```

### 错误处理

```sql
-- ✅ 处理异常
CREATE OR REPLACE FUNCTION safe_divide(a NUMERIC, b NUMERIC)
RETURNS NUMERIC AS $$
BEGIN
    BEGIN
        RETURN a / b;
    EXCEPTION
        WHEN division_by_zero THEN
            RAISE NOTICE 'Division by zero';
            RETURN NULL;
        WHEN OTHERS THEN
            RAISE NOTICE 'Error: %', SQLERRM;
            RETURN NULL;
    END;
END;
$$ LANGUAGE plpgsql;
```

### 版本控制

```sql
-- ✅ 使用迁移脚本管理 Schema
-- migrations/001_create_users_table.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/002_add_phone_to_users.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

---

## 运维实践

### 备份策略

```bash
# ✅ 每日全量备份
0 2 * * * pg_dump -U postgres mydb > /backup/mydb_$(date +\%Y\%m\%d).sql

# ✅ 启用 WAL 归档
# postgresql.conf
wal_level = replica
archive_mode = on
archive_command = 'cp %p /archive/%f'

# ✅ 定期测试恢复
```

### 监控

```sql
-- ✅ 监控慢查询
-- postgresql.conf
log_min_duration_statement = 1000

-- ✅ 使用 pg_stat_statements
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- ✅ 监控连接数
SELECT count(*) FROM pg_stat_activity;

-- ✅ 监控表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### 维护任务

```sql
-- ✅ 定期 VACUUM
VACUUM ANALYZE;

-- ✅ 定期更新统计信息
ANALYZE;

-- ✅ 检查未使用的索引
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;
```

### 配置优化

```conf
# postgresql.conf
# ✅ 根据硬件调整内存设置
shared_buffers = 4GB              # 25% RAM
effective_cache_size = 12GB       # 50-75% RAM
work_mem = 16MB                   # 根据并发连接数调整
maintenance_work_mem = 1GB

# ✅ 调整连接数
max_connections = 100

# ✅ 启用自动 VACUUM
autovacuum = on
autovacuum_max_workers = 3
```

---

## 代码审查清单

### 数据库设计
- [ ] 表名和列名符合命名规范
- [ ] 使用合适的数据类型
- [ ] 添加必要的约束
- [ ] 创建必要的索引

### 查询优化
- [ ] 使用 EXPLAIN ANALYZE 分析
- [ ] 避免 SELECT *
- [ ] 使用索引列查询
- [ ] 避免在 WHERE 子句中使用函数

### 安全
- [ ] 使用参数化查询
- [ ] 授予最小权限
- [ ] 使用强密码
- [ ] 启用 SSL

### 维护
- [ ] 添加注释
- [ ] 处理异常
- [ ] 记录日志
- [ ] 版本控制

---

## 下一步学习

- [PostgreSQL 常见问题与解决方案](./postgresql-troubleshooting.md)
- [PostgreSQL 管理与运维](./postgresql-admin.md)
- [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md)

---

*最后更新：2024年*

