# PostgreSQL 事务与并发控制

## 目录
- [事务基础](#事务基础)
- [ACID 特性](#acid-特性)
- [事务隔离级别](#事务隔离级别)
- [并发控制](#并发控制)
- [锁机制](#锁机制)
- [MVCC（多版本并发控制）](#mvcc多版本并发控制)
- [死锁处理](#死锁处理)

---

## 事务基础

### 什么是事务

事务是一组数据库操作，要么全部成功，要么全部失败，保证数据一致性。

### 事务的基本操作

```sql
-- 开始事务
BEGIN;
-- 或
START TRANSACTION;

-- 提交事务
COMMIT;

-- 回滚事务
ROLLBACK;

-- 保存点
SAVEPOINT savepoint_name;
ROLLBACK TO savepoint_name;
RELEASE SAVEPOINT savepoint_name;
```

### 事务示例

```sql
-- 示例：转账操作
BEGIN;

-- 从账户 A 扣款
UPDATE accounts SET balance = balance - 100 WHERE id = 1;

-- 向账户 B 存款
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- 如果一切正常，提交
COMMIT;

-- 如果出错，回滚
-- ROLLBACK;
```

### 自动提交

PostgreSQL 默认启用自动提交模式：

```sql
-- 查看自动提交状态
SHOW autocommit;

-- 关闭自动提交（在 psql 中）
\set AUTOCOMMIT off

-- 每个语句都需要显式提交
BEGIN;
INSERT INTO users (username) VALUES ('john');
COMMIT;
```

---

## ACID 特性

### Atomicity（原子性）

事务中的所有操作要么全部成功，要么全部失败。

```sql
BEGIN;
INSERT INTO orders (user_id, amount) VALUES (1, 100);
INSERT INTO order_items (order_id, product_id, quantity) VALUES (1, 1, 2);
-- 如果任何一步失败，整个事务回滚
COMMIT;
```

### Consistency（一致性）

事务执行前后，数据库保持一致状态。

```sql
-- 示例：确保账户余额不为负
BEGIN;
UPDATE accounts SET balance = balance - 200 WHERE id = 1;
-- 如果余额变为负数，触发约束错误，事务回滚
COMMIT;
```

### Isolation（隔离性）

并发事务之间相互隔离，互不干扰。

```sql
-- 事务 A
BEGIN;
SELECT balance FROM accounts WHERE id = 1;  -- 读取 1000
UPDATE accounts SET balance = 900 WHERE id = 1;
-- 事务 B 此时看不到这个更新
COMMIT;

-- 事务 B
BEGIN;
SELECT balance FROM accounts WHERE id = 1;  -- 仍然读取 1000（取决于隔离级别）
COMMIT;
```

### Durability（持久性）

已提交的事务对数据库的修改是永久的。

```sql
BEGIN;
INSERT INTO users (username) VALUES ('john');
COMMIT;  -- 提交后，即使系统崩溃，数据也不会丢失
```

---

## 事务隔离级别

PostgreSQL 支持四种隔离级别（基于 SQL 标准）：

### 1. READ UNCOMMITTED（读未提交）

**实际上 PostgreSQL 不支持此级别**，会自动升级为 READ COMMITTED。

```sql
-- 设置隔离级别
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
BEGIN;
SELECT * FROM users;
COMMIT;
```

### 2. READ COMMITTED（读已提交）- 默认级别

只能读取已提交的数据，避免脏读。

```sql
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
BEGIN;

-- 事务 A
SELECT balance FROM accounts WHERE id = 1;  -- 读取 1000

-- 事务 B 提交了更新
-- UPDATE accounts SET balance = 900 WHERE id = 1; COMMIT;

-- 事务 A 再次读取
SELECT balance FROM accounts WHERE id = 1;  -- 读取 900（可能不同）

COMMIT;
```

**特点**：
- ✅ 避免脏读
- ❌ 可能出现不可重复读
- ❌ 可能出现幻读

### 3. REPEATABLE READ（可重复读）

在同一事务中，多次读取同一数据结果一致。

```sql
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
BEGIN;

-- 第一次读取
SELECT balance FROM accounts WHERE id = 1;  -- 读取 1000

-- 即使其他事务提交了更新，这里仍然读取 1000
SELECT balance FROM accounts WHERE id = 1;  -- 仍然读取 1000

COMMIT;
```

**特点**：
- ✅ 避免脏读
- ✅ 避免不可重复读
- ❌ 可能出现幻读（PostgreSQL 中实际上可以避免）

### 4. SERIALIZABLE（可串行化）

最严格的隔离级别，完全避免并发问题。

```sql
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
BEGIN;

-- 所有事务按顺序执行，就像串行执行一样
SELECT * FROM accounts WHERE balance > 1000;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;

COMMIT;
```

**特点**：
- ✅ 避免脏读
- ✅ 避免不可重复读
- ✅ 避免幻读
- ⚠️ 可能出现序列化错误，需要重试

### 隔离级别对比

| 隔离级别 | 脏读 | 不可重复读 | 幻读 | 性能 |
|---------|------|-----------|------|------|
| READ UNCOMMITTED | ❌ | ❌ | ❌ | 最快 |
| READ COMMITTED | ✅ | ❌ | ❌ | 快 |
| REPEATABLE READ | ✅ | ✅ | ⚠️ | 中等 |
| SERIALIZABLE | ✅ | ✅ | ✅ | 最慢 |

---

## 并发控制

### 并发问题

#### 1. 脏读（Dirty Read）

读取到未提交的数据。

```sql
-- 事务 A
BEGIN;
UPDATE accounts SET balance = 900 WHERE id = 1;  -- 未提交

-- 事务 B（READ UNCOMMITTED）
BEGIN;
SELECT balance FROM accounts WHERE id = 1;  -- 读取 900（脏读）
COMMIT;

-- 事务 A 回滚
ROLLBACK;  -- 事务 B 读取的数据无效
```

#### 2. 不可重复读（Non-repeatable Read）

同一事务中，多次读取同一数据结果不同。

```sql
-- 事务 A
BEGIN;
SELECT balance FROM accounts WHERE id = 1;  -- 读取 1000

-- 事务 B
BEGIN;
UPDATE accounts SET balance = 900 WHERE id = 1;
COMMIT;

-- 事务 A 再次读取
SELECT balance FROM accounts WHERE id = 1;  -- 读取 900（不同）
COMMIT;
```

#### 3. 幻读（Phantom Read）

同一事务中，多次查询返回的行数不同。

```sql
-- 事务 A
BEGIN;
SELECT COUNT(*) FROM accounts WHERE balance > 1000;  -- 返回 10

-- 事务 B
BEGIN;
INSERT INTO accounts (balance) VALUES (2000);
COMMIT;

-- 事务 A 再次查询
SELECT COUNT(*) FROM accounts WHERE balance > 1000;  -- 返回 11（幻读）
COMMIT;
```

---

## 锁机制

### 锁的类型

#### 1. 表级锁

```sql
-- 共享锁（SELECT）
LOCK TABLE users IN SHARE MODE;

-- 排他锁（UPDATE, DELETE, INSERT）
LOCK TABLE users IN EXCLUSIVE MODE;

-- 访问排他锁（ALTER TABLE, DROP TABLE）
LOCK TABLE users IN ACCESS EXCLUSIVE MODE;
```

#### 2. 行级锁

PostgreSQL 自动在行级别加锁：

```sql
-- FOR UPDATE：排他锁
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;

-- FOR SHARE：共享锁
SELECT * FROM accounts WHERE id = 1 FOR SHARE;

-- FOR NO KEY UPDATE：非键排他锁
SELECT * FROM accounts WHERE id = 1 FOR NO KEY UPDATE;

-- FOR KEY SHARE：键共享锁
SELECT * FROM accounts WHERE id = 1 FOR KEY SHARE;
```

### 锁的兼容性

| 锁类型 | ACCESS SHARE | ROW SHARE | ROW EXCLUSIVE | SHARE UPDATE EXCLUSIVE | SHARE | SHARE ROW EXCLUSIVE | EXCLUSIVE | ACCESS EXCLUSIVE |
|--------|--------------|-----------|---------------|------------------------|-------|---------------------|-----------|------------------|
| ACCESS SHARE | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| ROW SHARE | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| ROW EXCLUSIVE | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| SHARE UPDATE EXCLUSIVE | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| SHARE | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| SHARE ROW EXCLUSIVE | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| EXCLUSIVE | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| ACCESS EXCLUSIVE | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |

### 查看锁

```sql
-- 查看当前锁
SELECT 
    locktype,
    relation::regclass,
    mode,
    granted,
    pid
FROM pg_locks
WHERE relation = 'accounts'::regclass;

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

## MVCC（多版本并发控制）

### MVCC 原理

PostgreSQL 使用 MVCC 实现并发控制，每个事务看到数据库的一个快照。

### 系统列

每个表都有隐藏的系统列：

- `xmin`: 插入该行的事务 ID
- `xmax`: 删除该行的事务 ID（如果未删除则为 0）
- `ctid`: 行的物理位置

```sql
-- 查看系统列
SELECT xmin, xmax, ctid, * FROM users WHERE id = 1;
```

### 事务 ID

```sql
-- 查看当前事务 ID
SELECT txid_current();

-- 查看快照
SELECT txid_current_snapshot();
```

### 可见性规则

一行数据对事务可见，当且仅当：
1. `xmin` < 当前事务 ID
2. `xmax` = 0 或 `xmax` > 当前事务 ID
3. 插入该行的事务已提交
4. 删除该行的事务未提交或不存在

---

## 死锁处理

### 什么是死锁

两个或多个事务相互等待对方释放锁，导致无法继续执行。

```sql
-- 事务 A
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- 等待事务 B 释放 id = 2 的锁
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- 事务 B
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 2;
-- 等待事务 A 释放 id = 1 的锁
UPDATE accounts SET balance = balance + 100 WHERE id = 1;
```

### 死锁检测

PostgreSQL 自动检测死锁，并回滚其中一个事务。

```sql
-- 查看死锁日志（在 postgresql.conf 中）
log_lock_waits = on
deadlock_timeout = 1s
```

### 避免死锁

1. **按相同顺序访问资源**

```sql
-- ✅ 好：两个事务都按 id 顺序访问
-- 事务 A
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- 事务 B
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
```

2. **使用超时**

```sql
-- 设置锁超时
SET lock_timeout = '5s';

BEGIN;
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
-- 如果 5 秒内无法获取锁，抛出错误
COMMIT;
```

3. **减少事务时间**

```sql
-- ✅ 好：事务尽可能短
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
COMMIT;  -- 立即提交

-- ❌ 坏：长时间持有锁
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- ... 执行其他耗时操作 ...
COMMIT;
```

4. **使用较低的隔离级别**

```sql
-- 如果可能，使用 READ COMMITTED 而不是 SERIALIZABLE
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
```

---

## 实践示例

### 示例 1：银行转账

```sql
-- 安全的转账实现
BEGIN;

-- 检查余额
SELECT balance FROM accounts WHERE id = 1 FOR UPDATE;

-- 扣款
UPDATE accounts SET balance = balance - 100 WHERE id = 1;

-- 检查余额是否足够
DO $$
DECLARE
    current_balance DECIMAL;
BEGIN
    SELECT balance INTO current_balance FROM accounts WHERE id = 1;
    IF current_balance < 0 THEN
        RAISE EXCEPTION 'Insufficient balance';
    END IF;
END $$;

-- 存款
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

COMMIT;
```

### 示例 2：库存管理

```sql
-- 使用 SELECT FOR UPDATE 防止超卖
BEGIN;

-- 锁定库存行
SELECT stock FROM products WHERE id = 1 FOR UPDATE;

-- 检查库存
DO $$
DECLARE
    current_stock INTEGER;
BEGIN
    SELECT stock INTO current_stock FROM products WHERE id = 1;
    IF current_stock < 1 THEN
        RAISE EXCEPTION 'Out of stock';
    END IF;
END $$;

-- 减少库存
UPDATE products SET stock = stock - 1 WHERE id = 1;

COMMIT;
```

### 示例 3：序列化事务

```sql
-- 使用 SERIALIZABLE 隔离级别
BEGIN ISOLATION LEVEL SERIALIZABLE;

-- 检查条件
SELECT COUNT(*) FROM orders WHERE user_id = 1 AND created_at > CURRENT_DATE;

-- 如果条件满足，执行操作
INSERT INTO orders (user_id, amount) VALUES (1, 100);

COMMIT;
-- 如果出现序列化错误，需要重试
```

---

## 最佳实践

1. **事务设计**
   - 保持事务尽可能短
   - 避免在事务中执行长时间操作
   - 按相同顺序访问资源

2. **隔离级别选择**
   - 默认使用 READ COMMITTED
   - 需要可重复读时使用 REPEATABLE READ
   - 仅在必要时使用 SERIALIZABLE

3. **锁的使用**
   - 尽量使用行级锁而不是表级锁
   - 使用 SELECT FOR UPDATE 明确锁定
   - 设置合理的锁超时

4. **错误处理**
   - 处理序列化错误和死锁
   - 实现重试机制
   - 记录错误日志

---

## MySQL 对比提示

如果你是 MySQL 用户，以下是对比要点：

- **默认隔离级别**：PostgreSQL 默认 `READ COMMITTED`，MySQL (InnoDB) 默认 `REPEATABLE READ`
- **READ UNCOMMITTED**：PostgreSQL 不支持（自动升级为 READ COMMITTED），MySQL 支持
- **锁语法**：PostgreSQL 使用 `FOR UPDATE` / `FOR SHARE`，MySQL 使用 `FOR UPDATE` / `LOCK IN SHARE MODE`（8.0+ 支持 `FOR SHARE`）

更多详细对比请参考 [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md)。

## 下一步学习

- [PostgreSQL 存储过程与函数](./postgresql-functions.md)
- [PostgreSQL 管理与运维](./postgresql-admin.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)
- [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md)

---

*最后更新：2024年*

