# PostgreSQL 常见问题与解决方案

## 目录
- [连接问题](#连接问题)
- [性能问题](#性能问题)
- [锁问题](#锁问题)
- [空间问题](#空间问题)
- [数据损坏](#数据损坏)
- [复制问题](#复制问题)
- [常见错误](#常见错误)

---

## 连接问题

### 问题 1：无法连接到数据库

**错误信息**：
```
FATAL: no pg_hba.conf entry for host
```

**解决方案**：

```conf
# 检查 pg_hba.conf
# 添加允许连接的规则
host    all    all    192.168.1.0/24    md5

# 重启 PostgreSQL
sudo systemctl restart postgresql
```

### 问题 2：连接数过多

**错误信息**：
```
FATAL: too many connections
```

**解决方案**：

```sql
-- 查看当前连接数
SELECT count(*) FROM pg_stat_activity;

-- 查看最大连接数
SHOW max_connections;

-- 终止空闲连接
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle'
    AND state_change < now() - interval '5 minutes';

-- 增加最大连接数（需要重启）
-- postgresql.conf
max_connections = 200
```

### 问题 3：密码认证失败

**错误信息**：
```
FATAL: password authentication failed
```

**解决方案**：

```sql
-- 重置密码
ALTER USER username WITH PASSWORD 'new_password';

-- 检查 pg_hba.conf 配置
-- 确保使用正确的认证方法
```

---

## 性能问题

### 问题 1：查询很慢

**诊断步骤**：

```sql
-- 1. 使用 EXPLAIN ANALYZE
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'john@example.com';

-- 2. 检查是否有索引
\d users

-- 3. 查看统计信息
SELECT 
    schemaname,
    tablename,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables
WHERE tablename = 'users';

-- 4. 更新统计信息
ANALYZE users;
```

**解决方案**：

```sql
-- 创建索引
CREATE INDEX idx_users_email ON users(email);

-- 使用 pg_stat_statements 找出慢查询
SELECT 
    query,
    calls,
    mean_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 问题 2：表膨胀

**诊断**：

```sql
-- 查看表大小和死元组
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    n_dead_tup,
    n_live_tup,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
```

**解决方案**：

```sql
-- 执行 VACUUM
VACUUM ANALYZE users;

-- 如果表很大，使用 VACUUM FULL（需要锁表）
VACUUM FULL users;

-- 配置自动 VACUUM
-- postgresql.conf
autovacuum = on
autovacuum_vacuum_scale_factor = 0.1
```

### 问题 3：索引未使用

**诊断**：

```sql
-- 查看未使用的索引
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;
```

**解决方案**：

```sql
-- 删除未使用的索引
DROP INDEX idx_unused_index;

-- 重建索引
REINDEX INDEX idx_users_email;
```

---

## 锁问题

### 问题 1：查询被阻塞

**诊断**：

```sql
-- 查看阻塞的查询
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocking_locks.pid AS blocking_pid,
    blocked_activity.query AS blocked_query,
    blocking_activity.query AS blocking_query,
    blocked_activity.application_name AS blocked_app,
    blocking_activity.application_name AS blocking_app
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted
    AND blocking_locks.pid != blocked_locks.pid;
```

**解决方案**：

```sql
-- 终止阻塞的查询
SELECT pg_terminate_backend(blocking_pid)
FROM (
    SELECT blocking_locks.pid AS blocking_pid
    FROM pg_catalog.pg_locks blocked_locks
    JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
    WHERE NOT blocked_locks.granted
        AND blocking_locks.pid != blocked_locks.pid
) blocking;
```

### 问题 2：死锁

**错误信息**：
```
ERROR: deadlock detected
```

**解决方案**：

```sql
-- 查看死锁日志
-- postgresql.conf
log_lock_waits = on
deadlock_timeout = 1s

-- 应用层需要重试事务
-- 确保按相同顺序访问资源
```

---

## 空间问题

### 问题 1：磁盘空间不足

**诊断**：

```sql
-- 查看数据库大小
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
ORDER BY pg_database_size(pg_database.datname) DESC;

-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

**解决方案**：

```sql
-- 清理未使用的空间
VACUUM FULL users;

-- 删除旧数据
DELETE FROM logs WHERE created_at < '2020-01-01';

-- 清理 WAL 文件
SELECT pg_switch_wal();  -- 切换 WAL
-- 删除旧的 WAL 文件（确保已归档）
```

### 问题 2：WAL 文件占用空间

**诊断**：

```bash
# 查看 WAL 目录大小
du -sh /var/lib/postgresql/14/data/pg_wal
```

**解决方案**：

```sql
-- 检查 WAL 归档
SELECT * FROM pg_stat_archiver;

-- 配置 WAL 保留
-- postgresql.conf
max_wal_size = 4GB
min_wal_size = 1GB

-- 使用复制槽防止 WAL 被删除
SELECT * FROM pg_replication_slots;
```

---

## 数据损坏

### 问题 1：数据页损坏

**错误信息**：
```
ERROR: invalid page in block
```

**解决方案**：

```sql
-- 检查数据完整性
SELECT * FROM users WHERE ctid = '(0,1)';

-- 如果可能，从备份恢复
-- 如果无法恢复，尝试导出数据
pg_dump -t users mydb > users_backup.sql
```

### 问题 2：索引损坏

**错误信息**：
```
ERROR: index is corrupted
```

**解决方案**：

```sql
-- 重建索引
REINDEX INDEX idx_users_email;

-- 重建表的所有索引
REINDEX TABLE users;
```

---

## 复制问题

### 问题 1：复制延迟

**诊断**：

```sql
-- 在主库查看复制状态
SELECT 
    application_name,
    client_addr,
    state,
    pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) AS replication_lag_bytes
FROM pg_stat_replication;

-- 在从库查看延迟
SELECT 
    EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) AS replication_lag_seconds;
```

**解决方案**：

```sql
-- 检查网络连接
-- 检查从库性能
-- 增加 WAL 保留
-- postgresql.conf
wal_keep_size = 1GB
```

### 问题 2：复制中断

**诊断**：

```sql
-- 查看复制状态
SELECT * FROM pg_stat_replication;

-- 查看错误日志
-- /var/lib/postgresql/14/data/log/postgresql-*.log
```

**解决方案**：

```bash
# 重新同步从库
# 1. 停止从库
sudo systemctl stop postgresql

# 2. 清空数据目录
sudo rm -rf /var/lib/postgresql/14/data/*

# 3. 重新创建基础备份
sudo -u postgres pg_basebackup \
    -h master_host \
    -D /var/lib/postgresql/14/data \
    -U replicator \
    -P \
    -v \
    -R \
    -W

# 4. 启动从库
sudo systemctl start postgresql
```

---

## 常见错误

### 错误 1：关系不存在

**错误信息**：
```
ERROR: relation "table_name" does not exist
```

**解决方案**：

```sql
-- 检查表是否存在
SELECT * FROM information_schema.tables WHERE table_name = 'table_name';

-- 检查模式
SELECT current_schema();

-- 使用完整名称
SELECT * FROM schema_name.table_name;
```

### 错误 2：违反唯一约束

**错误信息**：
```
ERROR: duplicate key value violates unique constraint
```

**解决方案**：

```sql
-- 查找重复值
SELECT email, COUNT(*) 
FROM users 
GROUP BY email 
HAVING COUNT(*) > 1;

-- 删除重复数据
DELETE FROM users 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM users 
    GROUP BY email
);
```

### 错误 3：违反外键约束

**错误信息**：
```
ERROR: insert or update on table violates foreign key constraint
```

**解决方案**：

```sql
-- 检查外键约束
SELECT 
    tc.constraint_name,
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY';

-- 确保引用的数据存在
```

### 错误 4：事务中无法执行某些操作

**错误信息**：
```
ERROR: cannot execute ... in a read-only transaction
```

**解决方案**：

```sql
-- 某些操作（如 VACUUM、CREATE INDEX CONCURRENTLY）不能在事务中执行
-- 确保在事务外执行
VACUUM users;  -- 不在事务中
```

---

## 诊断工具

### pg_stat_statements

```sql
-- 查看慢查询
SELECT 
    query,
    calls,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 日志分析

```bash
# 使用 pgBadger 分析日志
pgbadger /var/lib/postgresql/14/data/log/postgresql-*.log -o report.html
```

### 系统监控

```sql
-- 查看系统资源使用
SELECT * FROM pg_stat_database;
SELECT * FROM pg_stat_bgwriter;
```

---

## 预防措施

1. **定期备份**
   - 每日全量备份
   - 启用 WAL 归档
   - 测试恢复流程

2. **监控**
   - 监控连接数
   - 监控慢查询
   - 监控磁盘空间
   - 监控复制延迟

3. **维护**
   - 定期 VACUUM
   - 定期 ANALYZE
   - 定期检查日志

4. **文档**
   - 记录配置变更
   - 记录故障处理过程
   - 维护运行手册

---

## 下一步学习

- [PostgreSQL 管理与运维](./postgresql-admin.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)
- [PostgreSQL 复制与高可用](./postgresql-replication.md)

---

*最后更新：2024年*

