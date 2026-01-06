# PostgreSQL 管理与运维

## 目录
- [用户与权限管理](#用户与权限管理)
- [备份与恢复](#备份与恢复)
- [性能监控](#性能监控)
- [日志管理](#日志管理)
- [维护任务](#维护任务)
- [配置优化](#配置优化)
- [故障排查](#故障排查)

---

## 用户与权限管理

### 用户管理

```sql
-- 创建用户
CREATE USER username WITH PASSWORD 'password';

-- 创建超级用户
CREATE USER admin WITH PASSWORD 'password' SUPERUSER;

-- 创建用户并设置属性
CREATE USER developer WITH
    PASSWORD 'password'
    CREATEDB
    CREATEROLE
    LOGIN
    VALID UNTIL '2025-12-31';

-- 修改用户
ALTER USER username WITH PASSWORD 'newpassword';
ALTER USER username CREATEDB;
ALTER USER username VALID UNTIL '2025-12-31';

-- 删除用户
DROP USER username;

-- 查看用户
\du                    -- psql 命令
SELECT * FROM pg_user; -- SQL 查询
```

### 角色管理

```sql
-- 创建角色
CREATE ROLE role_name;

-- 角色和用户（PostgreSQL 中用户是带 LOGIN 权限的角色）
CREATE ROLE app_user WITH LOGIN PASSWORD 'password';

-- 授予角色权限
GRANT role_name TO username;

-- 撤销角色权限
REVOKE role_name FROM username;

-- 查看角色
\dg                    -- psql 命令
SELECT * FROM pg_roles; -- SQL 查询
```

### 权限管理

#### 数据库权限

```sql
-- 授予连接权限
GRANT CONNECT ON DATABASE mydb TO username;

-- 授予创建权限
GRANT CREATE ON DATABASE mydb TO username;

-- 撤销权限
REVOKE CONNECT ON DATABASE mydb FROM username;
```

#### 模式权限

```sql
-- 授予使用权限
GRANT USAGE ON SCHEMA public TO username;

-- 授予创建权限
GRANT CREATE ON SCHEMA public TO username;

-- 授予所有权限
GRANT ALL ON SCHEMA public TO username;
```

#### 表权限

```sql
-- 授予 SELECT 权限
GRANT SELECT ON TABLE users TO username;

-- 授予多个权限
GRANT SELECT, INSERT, UPDATE ON TABLE users TO username;

-- 授予所有权限
GRANT ALL ON TABLE users TO username;

-- 授予所有表的权限
GRANT ALL ON ALL TABLES IN SCHEMA public TO username;

-- 撤销权限
REVOKE SELECT ON TABLE users FROM username;
```

#### 列权限

```sql
-- 授予特定列的权限
GRANT SELECT (id, username) ON TABLE users TO username;
GRANT UPDATE (email) ON TABLE users TO username;
```

#### 序列权限

```sql
-- 授予序列权限
GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO username;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO username;
```

### 权限查看

```sql
-- 查看表权限
SELECT 
    grantee,
    privilege_type
FROM information_schema.table_privileges
WHERE table_name = 'users';

-- 查看用户权限
SELECT 
    grantee,
    table_schema,
    table_name,
    privilege_type
FROM information_schema.table_privileges
WHERE grantee = 'username';
```

---

## 备份与恢复

### pg_dump（逻辑备份）

#### 备份单个数据库

```bash
# 基本备份
pg_dump -U username -d database_name > backup.sql

# 压缩备份
pg_dump -U username -d database_name | gzip > backup.sql.gz

# 自定义格式（可并行恢复）
pg_dump -U username -d database_name -F c -f backup.dump

# 只备份数据
pg_dump -U username -d database_name -a > data_only.sql

# 只备份结构
pg_dump -U username -d database_name -s > schema_only.sql

# 备份特定表
pg_dump -U username -d database_name -t table_name > table_backup.sql

# 排除特定表
pg_dump -U username -d database_name -T excluded_table > backup.sql
```

#### 备份所有数据库

```bash
# 备份所有数据库
pg_dumpall -U username > all_databases.sql

# 只备份全局对象（用户、角色等）
pg_dumpall -U username -g > globals.sql
```

### pg_restore（恢复）

```bash
# 恢复自定义格式备份
pg_restore -U username -d database_name backup.dump

# 恢复并创建数据库
pg_restore -U username -d new_database backup.dump

# 只恢复数据
pg_restore -U username -d database_name -a backup.dump

# 只恢复结构
pg_restore -U username -d database_name -s backup.dump

# 并行恢复（加快速度）
pg_restore -U username -d database_name -j 4 backup.dump
```

### psql 恢复

```bash
# 恢复 SQL 文件
psql -U username -d database_name < backup.sql

# 恢复并显示进度
psql -U username -d database_name -f backup.sql
```

### 物理备份（文件系统备份）

```bash
# 停止 PostgreSQL
sudo systemctl stop postgresql

# 备份数据目录
sudo tar -czf pgdata_backup.tar.gz /var/lib/postgresql/14/data

# 启动 PostgreSQL
sudo systemctl start postgresql
```

### 连续归档（WAL 归档）

#### 配置归档

```conf
# postgresql.conf
wal_level = replica
archive_mode = on
archive_command = 'cp %p /path/to/archive/%f'
```

#### 基础备份

```bash
# 创建基础备份
pg_basebackup -D /path/to/backup -Ft -z -P

# 流式备份
pg_basebackup -D /path/to/backup -Ft -z -P -h localhost -U replicator
```

---

## 性能监控

### 查看连接

```sql
-- 查看当前连接
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query_start,
    state_change,
    wait_event_type,
    wait_event,
    query
FROM pg_stat_activity
WHERE datname = 'mydb';

-- 查看连接数
SELECT 
    count(*) as total_connections,
    count(*) FILTER (WHERE state = 'active') as active_connections,
    count(*) FILTER (WHERE state = 'idle') as idle_connections
FROM pg_stat_activity;
```

### 查看慢查询

```sql
-- 查看长时间运行的查询
SELECT 
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query,
    state
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes'
    AND state != 'idle';
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
    max_exec_time,
    rows
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- 重置统计
SELECT pg_stat_statements_reset();
```

### 表统计信息

```sql
-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_indexes_size(schemaname||'.'||tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 查看数据库大小
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
ORDER BY pg_database_size(pg_database.datname) DESC;
```

### 索引使用情况

```sql
-- 查看索引使用统计
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan;

-- 查找未使用的索引
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;
```

---

## 日志管理

### 日志配置

```conf
# postgresql.conf

# 日志记录
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB

# 日志级别
log_min_messages = warning        # 记录警告及以上
log_min_error_statement = error   # 记录错误语句
log_min_duration_statement = 1000 # 记录超过1秒的查询

# 日志内容
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_timezone = 'Asia/Shanghai'
log_statement = 'all'             # 记录所有SQL
log_duration = on                 # 记录执行时间
log_connections = on              # 记录连接
log_disconnections = on           # 记录断开
log_lock_waits = on               # 记录锁等待
```

### 查看日志

```bash
# 查看日志文件
tail -f /var/lib/postgresql/14/data/log/postgresql-*.log

# 使用 less 查看
less /var/lib/postgresql/14/data/log/postgresql-*.log

# 搜索错误
grep ERROR /var/lib/postgresql/14/data/log/postgresql-*.log
```

### 使用 pgBadger 分析日志

```bash
# 安装 pgBadger
sudo apt install pgbadger

# 分析日志
pgbadger /var/lib/postgresql/14/data/log/postgresql-*.log -o report.html

# 查看报告
open report.html
```

---

## 维护任务

### VACUUM

```sql
-- 手动 VACUUM
VACUUM;

-- VACUUM 特定表
VACUUM users;

-- VACUUM ANALYZE（同时更新统计信息）
VACUUM ANALYZE;

-- VACUUM FULL（回收空间，需要锁表）
VACUUM FULL users;

-- 查看需要 VACUUM 的表
SELECT 
    schemaname,
    tablename,
    n_dead_tup,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
WHERE n_dead_tup > 0
ORDER BY n_dead_tup DESC;
```

### ANALYZE

```sql
-- 更新统计信息
ANALYZE;

-- 分析特定表
ANALYZE users;

-- 查看统计信息最后更新时间
SELECT 
    schemaname,
    tablename,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables;
```

### REINDEX

```sql
-- 重建索引
REINDEX INDEX idx_users_email;

-- 重建表的所有索引
REINDEX TABLE users;

-- 重建数据库的所有索引
REINDEX DATABASE mydb;
```

### 自动维护

```conf
# postgresql.conf

# 自动 VACUUM
autovacuum = on
autovacuum_max_workers = 3
autovacuum_naptime = 1min

# 自动 ANALYZE
autovacuum_analyze_scale_factor = 0.1
```

---

## 配置优化

### 内存配置

```conf
# postgresql.conf

# 共享内存
shared_buffers = 4GB              # 25% 的 RAM

# 工作内存
work_mem = 16MB                   # 每个操作的内存

# 维护工作内存
maintenance_work_mem = 1GB        # 维护操作的内存

# 有效缓存大小
effective_cache_size = 12GB      # 50-75% 的 RAM
```

### 连接配置

```conf
# 最大连接数
max_connections = 100

# 连接超时
statement_timeout = 0             # 0 = 无限制
lock_timeout = 0
idle_in_transaction_session_timeout = 0
```

### 查询优化

```conf
# 查询规划器
random_page_cost = 1.1            # SSD: 1.1, HDD: 4.0
effective_io_concurrency = 200   # SSD: 200, HDD: 2

# 并行查询
max_parallel_workers_per_gather = 4
max_parallel_workers = 8
max_worker_processes = 8
```

### WAL 配置

```conf
# WAL 设置
wal_level = replica
max_wal_size = 4GB
min_wal_size = 1GB
checkpoint_timeout = 15min
```

---

## 故障排查

### 连接问题

```sql
-- 查看连接限制
SHOW max_connections;

-- 查看当前连接数
SELECT count(*) FROM pg_stat_activity;

-- 终止连接
SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE pid <> pg_backend_pid();
```

### 锁问题

```sql
-- 查看锁
SELECT 
    locktype,
    relation::regclass,
    mode,
    granted,
    pid
FROM pg_locks
WHERE NOT granted;

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

### 性能问题

```sql
-- 查看慢查询
SELECT 
    pid,
    now() - query_start AS duration,
    query
FROM pg_stat_activity
WHERE state = 'active'
    AND now() - query_start > interval '1 minute';

-- 查看等待事件
SELECT 
    wait_event_type,
    wait_event,
    count(*)
FROM pg_stat_activity
WHERE wait_event IS NOT NULL
GROUP BY wait_event_type, wait_event;
```

### 磁盘空间

```sql
-- 查看数据库大小
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database;

-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 最佳实践

1. **定期备份**
   - 每日全量备份
   - 启用 WAL 归档
   - 测试恢复流程

2. **监控**
   - 监控连接数
   - 监控慢查询
   - 监控磁盘空间

3. **维护**
   - 定期 VACUUM
   - 定期 ANALYZE
   - 定期检查日志

4. **安全**
   - 使用强密码
   - 限制连接
   - 定期更新

---

## 下一步学习

- [PostgreSQL 复制与高可用](./postgresql-replication.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)
- [PostgreSQL 常见问题与解决方案](./postgresql-troubleshooting.md)

---

*最后更新：2024年*

