# PostgreSQL 复制与高可用

## 目录
- [复制基础](#复制基础)
- [流复制（Streaming Replication）](#流复制streaming-replication)
- [逻辑复制](#逻辑复制)
- [高可用方案](#高可用方案)
- [故障转移](#故障转移)
- [监控与维护](#监控与维护)

---

## 复制基础

### 复制类型

1. **物理复制（流复制）**
   - 复制整个数据库集群
   - 字节级复制
   - 主从必须相同版本

2. **逻辑复制**
   - 复制特定表或数据库
   - SQL 语句级复制
   - 可以跨版本

### 复制模式

- **同步复制**：主库等待从库确认
- **异步复制**：主库不等待从库确认

---

## 流复制（Streaming Replication）

### 主库配置

#### 1. 创建复制用户

```sql
-- 在主库创建复制用户
CREATE USER replicator WITH REPLICATION PASSWORD 'replicator_password';
```

#### 2. 配置 postgresql.conf

```conf
# 启用 WAL 归档
wal_level = replica

# 设置最大 WAL 发送进程数
max_wal_senders = 3

# 设置 WAL 保留数量
wal_keep_segments = 32  # PostgreSQL 13+ 使用 wal_keep_size
wal_keep_size = 512MB   # PostgreSQL 13+

# 启用热备
hot_standby = on
```

#### 3. 配置 pg_hba.conf

```
# 允许复制连接
host    replication    replicator    192.168.1.0/24    md5
```

#### 4. 重启主库

```bash
sudo systemctl restart postgresql
```

### 从库配置

#### 1. 基础备份

```bash
# 停止从库（如果已运行）
sudo systemctl stop postgresql

# 清空数据目录
sudo rm -rf /var/lib/postgresql/14/data/*

# 创建基础备份
sudo -u postgres pg_basebackup \
    -h master_host \
    -D /var/lib/postgresql/14/data \
    -U replicator \
    -P \
    -v \
    -R \
    -W

# -R 选项会自动创建 standby.signal 和 postgresql.auto.conf
```

#### 2. 配置 recovery.conf（PostgreSQL 12+ 使用 postgresql.auto.conf）

```conf
# postgresql.auto.conf（由 pg_basebackup -R 自动创建）
primary_conninfo = 'host=master_host port=5432 user=replicator password=replicator_password'
primary_slot_name = 'standby_slot'
```

#### 3. 创建恢复信号文件

```bash
# PostgreSQL 12+
touch /var/lib/postgresql/14/data/standby.signal
```

#### 4. 启动从库

```bash
sudo systemctl start postgresql
```

### 验证复制

```sql
-- 在主库查看复制状态
SELECT 
    application_name,
    client_addr,
    state,
    sync_state,
    sync_priority
FROM pg_stat_replication;

-- 在从库查看复制延迟
SELECT 
    pg_last_wal_receive_lsn(),
    pg_last_wal_replay_lsn(),
    pg_last_wal_receive_lsn() - pg_last_wal_replay_lsn() AS replication_lag;
```

---

## 逻辑复制

### 发布（Publication）

```sql
-- 在主库创建发布
CREATE PUBLICATION my_publication FOR TABLE users, orders;

-- 发布所有表
CREATE PUBLICATION all_tables FOR ALL TABLES;

-- 发布特定操作
CREATE PUBLICATION insert_only FOR TABLE users WITH (publish = 'insert');
```

### 订阅（Subscription）

```sql
-- 在从库创建订阅
CREATE SUBSCRIPTION my_subscription
    CONNECTION 'host=master_host port=5432 user=replicator password=replicator_password dbname=mydb'
    PUBLICATION my_publication;

-- 查看订阅状态
SELECT * FROM pg_subscription;
SELECT * FROM pg_stat_subscription;
```

---

## 高可用方案

### 1. 主从复制 + 手动故障转移

**架构**：
```
主库 (Primary) → 从库 (Standby)
```

**故障转移步骤**：
```bash
# 1. 提升从库为主库
sudo -u postgres pg_ctl promote -D /var/lib/postgresql/14/data

# 2. 更新应用连接配置
# 3. 原主库恢复后，重新配置为从库
```

### 2. 使用 Patroni

Patroni 是 PostgreSQL 高可用解决方案。

```yaml
# patroni.yml
scope: postgres
namespace: /db/
name: node1

restapi:
  listen: 0.0.0.0:8008
  connect_address: 192.168.1.10:8008

etcd:
  hosts: 192.168.1.20:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 30
    maximum_lag_on_failover: 1048576
  initdb:
  - encoding: UTF8
  - locale: en_US.UTF-8
  pg_hba:
  - host replication replicator 0.0.0.0/0 md5
  - host all all 0.0.0.0/0 md5
  users:
    admin:
      password: admin
      options:
        - createrole
        - createdb

postgresql:
  listen: 0.0.0.0:5432
  connect_address: 192.168.1.10:5432
  data_dir: /var/lib/postgresql/14/data
  pgpass: /tmp/pgpass
  authentication:
    replication:
      username: replicator
      password: replicator
    superuser:
      username: postgres
      password: postgres
  parameters:
    wal_level: replica
    hot_standby: "on"
    max_connections: 100
    max_wal_senders: 10
    wal_keep_size: 512MB
```

### 3. 使用 pgpool-II

pgpool-II 提供连接池和负载均衡。

```conf
# pgpool.conf
listen_addresses = '*'
port = 5432
socket_dir = '/var/run/postgresql'

backend_hostname0 = 'master_host'
backend_port0 = 5432
backend_weight0 = 1
backend_flag0 = 'ALLOW_TO_FAILOVER'

backend_hostname1 = 'standby_host'
backend_port1 = 5432
backend_weight1 = 1
backend_flag1 = 'ALLOW_TO_FAILOVER'

failover_on_backend_error = on
```

---

## 故障转移

### 自动故障转移

使用 Patroni 或类似工具实现自动故障转移。

### 手动故障转移

```bash
# 1. 检查从库状态
sudo -u postgres psql -c "SELECT pg_is_in_recovery();"

# 2. 提升从库
sudo -u postgres pg_ctl promote -D /var/lib/postgresql/14/data

# 3. 验证
sudo -u postgres psql -c "SELECT pg_is_in_recovery();"  # 应返回 false
```

### 重新配置原主库为从库

```bash
# 1. 停止原主库
sudo systemctl stop postgresql

# 2. 清空数据目录
sudo rm -rf /var/lib/postgresql/14/data/*

# 3. 从新主库创建基础备份
sudo -u postgres pg_basebackup \
    -h new_primary_host \
    -D /var/lib/postgresql/14/data \
    -U replicator \
    -P \
    -v \
    -R \
    -W

# 4. 启动
sudo systemctl start postgresql
```

---

## 监控与维护

### 监控复制延迟

```sql
-- 主库：查看复制状态
SELECT 
    application_name,
    client_addr,
    state,
    sync_state,
    pg_wal_lsn_diff(pg_current_wal_lsn(), sent_lsn) AS sending_lag,
    pg_wal_lsn_diff(sent_lsn, write_lsn) AS write_lag,
    pg_wal_lsn_diff(write_lsn, flush_lsn) AS flush_lag,
    pg_wal_lsn_diff(flush_lsn, replay_lsn) AS replay_lag
FROM pg_stat_replication;

-- 从库：查看延迟
SELECT 
    EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) AS replication_lag_seconds;
```

### 监控工具

```bash
# 使用 pg_stat_replication 视图
# 使用 Prometheus + postgres_exporter
# 使用 Zabbix 模板
```

### 维护任务

```sql
-- 检查 WAL 归档
SELECT * FROM pg_stat_archiver;

-- 检查复制槽
SELECT * FROM pg_replication_slots;

-- 创建复制槽（防止 WAL 被删除）
SELECT pg_create_physical_replication_slot('standby_slot');

-- 删除复制槽
SELECT pg_drop_replication_slot('standby_slot');
```

---

## 最佳实践

1. **复制配置**
   - 使用同步复制关键数据
   - 配置多个从库提高可用性
   - 使用复制槽防止 WAL 丢失

2. **监控**
   - 监控复制延迟
   - 监控主从库状态
   - 设置告警

3. **故障转移**
   - 定期测试故障转移流程
   - 准备故障转移脚本
   - 文档化恢复步骤

4. **备份**
   - 即使有复制，也要定期备份
   - 测试备份恢复流程

---

## 下一步学习

- [PostgreSQL 管理与运维](./postgresql-admin.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)
- [PostgreSQL 常见问题与解决方案](./postgresql-troubleshooting.md)

---

*最后更新：2024年*

