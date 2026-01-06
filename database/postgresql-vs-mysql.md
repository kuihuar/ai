# PostgreSQL vs MySQL 全面对比

## 目录
- [概述](#概述)
- [架构对比](#架构对比)
- [数据类型对比](#数据类型对比)
- [SQL 语法对比](#sql-语法对比)
- [索引对比](#索引对比)
- [事务与锁对比](#事务与锁对比)
- [存储引擎对比](#存储引擎对比)
- [性能特性对比](#性能特性对比)
- [高级功能对比](#高级功能对比)
- [使用场景建议](#使用场景建议)

---

## 概述

### PostgreSQL
- **类型**：对象关系型数据库（ORDBMS）
- **许可证**：PostgreSQL License（类似 BSD）
- **开发**：社区驱动
- **特点**：标准兼容性强、功能丰富、扩展性强

### MySQL
- **类型**：关系型数据库（RDBMS）
- **许可证**：GPL（社区版）或商业许可
- **开发**：Oracle 公司主导
- **特点**：易用、性能好、生态丰富

---

## 架构对比

### PostgreSQL

```sql
-- 架构层次
PostgreSQL 实例
├── 数据库 (Database)
│   ├── 模式 (Schema) - 命名空间
│   │   ├── 表 (Table)
│   │   ├── 视图 (View)
│   │   ├── 函数 (Function)
│   │   └── 其他对象
│   └── 系统目录
```

**特点**：
- 支持多模式（Schema）
- 单一存储引擎（但功能强大）
- 进程模型（每个连接一个进程）

### MySQL

```sql
-- 架构层次
MySQL 实例
├── 数据库 (Database) - 等同于 PostgreSQL 的 Schema
│   ├── 表 (Table)
│   ├── 视图 (View)
│   └── 其他对象
```

**特点**：
- 数据库即模式
- 多存储引擎（InnoDB、MyISAM、Memory 等）
- 线程模型（每个连接一个线程）

### 对比表

| 特性 | PostgreSQL | MySQL |
|------|-----------|-------|
| 架构模型 | 进程模型 | 线程模型 |
| 存储引擎 | 单一引擎 | 多引擎（InnoDB、MyISAM等） |
| 命名空间 | Database + Schema | Database |
| 连接开销 | 较高（进程） | 较低（线程） |
| 并发控制 | MVCC | MVCC（InnoDB）或表锁（MyISAM） |

---

## 数据类型对比

### 数值类型

| PostgreSQL | MySQL | 说明 |
|-----------|-------|------|
| `SMALLINT` | `SMALLINT` | 2字节整数 |
| `INTEGER` / `INT` | `INT` / `INTEGER` | 4字节整数 |
| `BIGINT` | `BIGINT` | 8字节整数 |
| `DECIMAL(p,s)` | `DECIMAL(p,s)` | 精确数值 |
| `NUMERIC(p,s)` | `NUMERIC(p,s)` | 精确数值 |
| `REAL` | `FLOAT` | 单精度浮点 |
| `DOUBLE PRECISION` | `DOUBLE` | 双精度浮点 |
| `SERIAL` | `AUTO_INCREMENT` | 自增整数 |

**示例**：

```sql
-- PostgreSQL
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    price DECIMAL(10, 2)
);

-- MySQL
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    price DECIMAL(10, 2)
);
```

### 字符类型

| PostgreSQL | MySQL | 说明 |
|-----------|-------|------|
| `CHAR(n)` | `CHAR(n)` | 固定长度 |
| `VARCHAR(n)` | `VARCHAR(n)` | 可变长度 |
| `TEXT` | `TEXT` | 长文本 |
| - | `TINYTEXT` | 小文本（MySQL特有） |
| - | `MEDIUMTEXT` | 中等文本（MySQL特有） |
| - | `LONGTEXT` | 超长文本（MySQL特有） |

**差异**：
- PostgreSQL 的 `TEXT` 无长度限制
- MySQL 有多种 TEXT 类型，有长度限制

### 日期时间类型

| PostgreSQL | MySQL | 说明 |
|-----------|-------|------|
| `DATE` | `DATE` | 日期 |
| `TIME` | `TIME` | 时间 |
| `TIMESTAMP` | `TIMESTAMP` | 日期时间（无时区） |
| `TIMESTAMPTZ` | `DATETIME` | 日期时间（PostgreSQL带时区） |
| `INTERVAL` | - | 时间间隔（PostgreSQL特有） |

**示例**：

```sql
-- PostgreSQL
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    duration INTERVAL
);

-- MySQL
CREATE TABLE events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    duration INT  -- 需要手动计算
);
```

### 高级类型

| 类型 | PostgreSQL | MySQL |
|------|-----------|-------|
| 数组 | ✅ `INTEGER[]` | ❌ 不支持 |
| JSON | ✅ `JSON` / `JSONB` | ✅ `JSON` (5.7+) |
| UUID | ✅ `UUID` | ❌ 使用 `CHAR(36)` |
| 网络地址 | ✅ `INET`, `CIDR` | ❌ 不支持 |
| 几何类型 | ✅ `POINT`, `POLYGON` | ✅ `POINT`, `POLYGON` (5.7+) |
| 全文搜索 | ✅ 内置 | ✅ 支持 |

**示例**：

```sql
-- PostgreSQL: 数组
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    tags TEXT[]
);
INSERT INTO products (tags) VALUES (ARRAY['electronics', 'computer']);

-- MySQL: 需要额外表或 JSON
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tags JSON
);
INSERT INTO products (tags) VALUES ('["electronics", "computer"]');
```

---

## SQL 语法对比

### 创建表

```sql
-- PostgreSQL
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- MySQL
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 自增主键

```sql
-- PostgreSQL: SERIAL
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50)
);

-- MySQL: AUTO_INCREMENT
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50)
);
```

### 字符串连接

```sql
-- PostgreSQL
SELECT 'Hello' || ' ' || 'World';  -- 使用 ||
SELECT CONCAT('Hello', ' ', 'World');  -- 也支持 CONCAT

-- MySQL
SELECT CONCAT('Hello', ' ', 'World');  -- 使用 CONCAT
SELECT 'Hello' || 'World';  -- 在 MySQL 中 || 是逻辑或
```

### 限制查询结果

```sql
-- PostgreSQL
SELECT * FROM users LIMIT 10;
SELECT * FROM users LIMIT 10 OFFSET 20;

-- MySQL
SELECT * FROM users LIMIT 10;
SELECT * FROM users LIMIT 20, 10;  -- MySQL 支持 LIMIT offset, count
```

### 日期函数

```sql
-- PostgreSQL
SELECT NOW();
SELECT CURRENT_DATE;
SELECT EXTRACT(YEAR FROM NOW());
SELECT NOW() + INTERVAL '1 day';

-- MySQL
SELECT NOW();
SELECT CURDATE();
SELECT YEAR(NOW());
SELECT DATE_ADD(NOW(), INTERVAL 1 DAY);
```

### 条件表达式

```sql
-- PostgreSQL: CASE 表达式
SELECT 
    CASE 
        WHEN age < 18 THEN 'Minor'
        WHEN age < 65 THEN 'Adult'
        ELSE 'Senior'
    END AS age_group
FROM users;

-- MySQL: 同样支持 CASE，还支持 IF
SELECT 
    IF(age < 18, 'Minor', 'Adult') AS age_group
FROM users;
```

### 窗口函数

```sql
-- PostgreSQL: 完整支持
SELECT 
    id,
    name,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) AS rank
FROM employees;

-- MySQL: 8.0+ 支持
SELECT 
    id,
    name,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) AS rank
FROM employees;
```

---

## 索引对比

### 索引类型

| 索引类型 | PostgreSQL | MySQL |
|---------|-----------|-------|
| B-tree | ✅ 默认 | ✅ 默认（InnoDB） |
| Hash | ✅ 支持 | ✅ Memory 引擎 |
| GiST | ✅ 支持 | ❌ 不支持 |
| GIN | ✅ 支持 | ❌ 不支持 |
| 全文索引 | ✅ GIN/GiST | ✅ FULLTEXT |
| 空间索引 | ✅ GiST | ✅ SPATIAL |

### 创建索引

```sql
-- PostgreSQL
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_name_gin ON users USING GIN (to_tsvector('english', name));
CREATE UNIQUE INDEX idx_users_username_unique ON users(username);

-- MySQL
CREATE INDEX idx_users_email ON users(email);
CREATE FULLTEXT INDEX idx_users_name_fulltext ON users(name);
CREATE UNIQUE INDEX idx_users_username_unique ON users(username);
```

### 覆盖索引

```sql
-- PostgreSQL 11+
CREATE INDEX idx_orders_user_date_inc 
ON orders(user_id, created_at) 
INCLUDE (amount, status);

-- MySQL 8.0+ (InnoDB)
-- 使用覆盖索引（索引包含所有查询列）
CREATE INDEX idx_orders_covering 
ON orders(user_id, created_at, amount, status);
```

---

## 事务与锁对比

### 隔离级别

| 隔离级别 | PostgreSQL | MySQL (InnoDB) |
|---------|-----------|----------------|
| READ UNCOMMITTED | ❌ 不支持（升级为 READ COMMITTED） | ✅ 支持 |
| READ COMMITTED | ✅ 默认 | ✅ 支持 |
| REPEATABLE READ | ✅ 支持 | ✅ 默认 |
| SERIALIZABLE | ✅ 支持 | ✅ 支持 |

### 锁机制

```sql
-- PostgreSQL: 行级锁
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
SELECT * FROM accounts WHERE id = 1 FOR SHARE;

-- MySQL (InnoDB): 行级锁
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
SELECT * FROM accounts WHERE id = 1 LOCK IN SHARE MODE;  -- 8.0+ 使用 FOR SHARE
```

### MVCC 实现

| 特性 | PostgreSQL | MySQL (InnoDB) |
|------|-----------|----------------|
| MVCC | ✅ 完整支持 | ✅ 支持 |
| 回滚段 | ✅ 使用 | ✅ 使用 |
| 可见性判断 | ✅ 基于事务ID | ✅ 基于 Read View |

---

## 存储引擎对比

### PostgreSQL

- **单一存储引擎**：功能强大，支持所有特性
- **可扩展性**：通过扩展添加功能

### MySQL

- **InnoDB**（默认）：
  - 支持事务、外键、行级锁
  - MVCC
  - 聚簇索引

- **MyISAM**：
  - 不支持事务
  - 表级锁
  - 全文索引（5.6-）
  - 适合读多写少

- **Memory**：
  - 内存存储
  - 表级锁
  - 适合临时数据

**对比**：

| 特性 | PostgreSQL | MySQL InnoDB | MySQL MyISAM |
|------|-----------|--------------|--------------|
| 事务 | ✅ | ✅ | ❌ |
| 外键 | ✅ | ✅ | ❌ |
| 行级锁 | ✅ | ✅ | ❌ |
| 崩溃恢复 | ✅ | ✅ | ⚠️ 有限 |
| 全文搜索 | ✅ | ✅ (5.6+) | ✅ |

---

## 性能特性对比

### 查询优化器

| 特性 | PostgreSQL | MySQL |
|------|-----------|-------|
| 优化器 | ✅ 基于成本的优化器 | ✅ 基于成本的优化器 |
| 执行计划 | ✅ EXPLAIN ANALYZE | ✅ EXPLAIN FORMAT=JSON |
| 统计信息 | ✅ 自动收集 | ✅ 自动收集 |
| 并行查询 | ✅ 支持 | ✅ 8.0+ 支持 |

### 连接池

```sql
-- PostgreSQL: 使用 pgBouncer 或 pgpool-II
-- 配置连接池参数
max_connections = 100

-- MySQL: 内置连接池
max_connections = 200
thread_cache_size = 8
```

### 分区

```sql
-- PostgreSQL: 原生分区
CREATE TABLE orders (
    id SERIAL,
    order_date DATE,
    amount DECIMAL(10, 2)
) PARTITION BY RANGE (order_date);

CREATE TABLE orders_2024 PARTITION OF orders
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

-- MySQL: 5.7+ 支持分区
CREATE TABLE orders (
    id INT AUTO_INCREMENT,
    order_date DATE,
    amount DECIMAL(10, 2),
    PRIMARY KEY (id, order_date)
) PARTITION BY RANGE (YEAR(order_date)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026)
);
```

---

## 高级功能对比

### JSON 支持

```sql
-- PostgreSQL: JSONB（二进制，性能更好）
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    data JSONB
);

SELECT * FROM documents WHERE data @> '{"status": "active"}';
CREATE INDEX idx_documents_data_gin ON documents USING GIN (data);

-- MySQL: JSON（5.7+）
CREATE TABLE documents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    data JSON
);

SELECT * FROM documents WHERE JSON_EXTRACT(data, '$.status') = 'active';
CREATE INDEX idx_documents_status ON documents ((CAST(data->>'$.status' AS CHAR(20))));
```

### 全文搜索

```sql
-- PostgreSQL: 内置全文搜索
CREATE INDEX idx_articles_content_gin 
ON articles USING GIN (to_tsvector('english', content));

SELECT * FROM articles 
WHERE to_tsvector('english', content) @@ to_tsquery('english', 'postgresql');

-- MySQL: FULLTEXT 索引
CREATE FULLTEXT INDEX idx_articles_content_fulltext ON articles(content);

SELECT * FROM articles 
WHERE MATCH(content) AGAINST('postgresql' IN NATURAL LANGUAGE MODE);
```

### 存储过程

```sql
-- PostgreSQL: PL/pgSQL
CREATE OR REPLACE FUNCTION get_user_count()
RETURNS INTEGER AS $$
DECLARE
    count INTEGER;
BEGIN
    SELECT COUNT(*) INTO count FROM users;
    RETURN count;
END;
$$ LANGUAGE plpgsql;

-- MySQL: 存储过程
DELIMITER //
CREATE PROCEDURE get_user_count()
BEGIN
    SELECT COUNT(*) FROM users;
END //
DELIMITER ;
```

### 触发器

```sql
-- PostgreSQL
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- MySQL
DELIMITER //
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
BEGIN
    SET NEW.updated_at = NOW();
END //
DELIMITER ;
```

### 复制

| 特性 | PostgreSQL | MySQL |
|------|-----------|-------|
| 流复制 | ✅ 物理复制 | ✅ 二进制日志复制 |
| 逻辑复制 | ✅ 9.4+ | ✅ 5.7+ |
| 主从复制 | ✅ | ✅ |
| 主主复制 | ⚠️ 需要第三方 | ✅ |
| 组复制 | ❌ | ✅ 5.7+ |

---

## 使用场景建议

### 选择 PostgreSQL 的场景

1. **复杂查询和数据分析**
   - 强大的 SQL 功能
   - 窗口函数、CTE、递归查询
   - 适合数据仓库和分析

2. **需要高级数据类型**
   - 数组、JSONB、UUID
   - 自定义类型
   - 地理空间数据（PostGIS）

3. **标准兼容性要求高**
   - 遵循 SQL 标准
   - 易于迁移

4. **需要扩展功能**
   - 丰富的扩展生态
   - 自定义扩展开发

5. **数据完整性要求高**
   - 严格的约束检查
   - 外键支持完善

### 选择 MySQL 的场景

1. **Web 应用**
   - 简单易用
   - 生态丰富（ORM、框架支持）
   - 文档完善

2. **高并发读写**
   - 线程模型性能好
   - 连接开销小
   - 适合 OLTP 场景

3. **需要多种存储引擎**
   - 根据场景选择引擎
   - MyISAM 适合读多写少

4. **已有 MySQL 生态**
   - 团队熟悉 MySQL
   - 现有工具和脚本

5. **云服务支持**
   - 各大云平台支持好
   - 托管服务丰富

### 性能对比总结

| 场景 | PostgreSQL | MySQL |
|------|-----------|-------|
| 复杂查询 | ✅ 更优 | ⚠️ 一般 |
| 简单查询 | ✅ 优秀 | ✅ 优秀 |
| 高并发写入 | ✅ 优秀 | ✅ 优秀 |
| 全文搜索 | ✅ 优秀 | ✅ 良好 |
| JSON 操作 | ✅ 优秀（JSONB） | ✅ 良好 |
| 地理空间 | ✅ 优秀（PostGIS） | ⚠️ 基础支持 |
| 扩展性 | ✅ 优秀 | ⚠️ 有限 |

---

## 迁移指南

### 从 MySQL 迁移到 PostgreSQL

1. **数据类型映射**
   ```sql
   -- AUTO_INCREMENT → SERIAL
   -- DATETIME → TIMESTAMP
   -- TEXT → TEXT（注意长度限制）
   -- ENUM → CHECK 约束或单独表
   ```

2. **语法差异**
   ```sql
   -- LIMIT offset, count → LIMIT count OFFSET offset
   -- CONCAT() → || 或 CONCAT()
   -- IF() → CASE WHEN
   ```

3. **存储引擎**
   - InnoDB → PostgreSQL 默认
   - MyISAM → 需要评估（无表锁）

### 从 PostgreSQL 迁移到 MySQL

1. **数据类型映射**
   ```sql
   -- SERIAL → AUTO_INCREMENT
   -- TIMESTAMPTZ → DATETIME
   -- ARRAY → JSON 或关联表
   -- UUID → CHAR(36) 或 BINARY(16)
   ```

2. **功能限制**
   - 无 Schema 支持
   - 数组需要转换
   - 某些高级类型不支持

---

## 总结

### PostgreSQL 优势
- ✅ SQL 标准兼容性强
- ✅ 功能丰富（数组、JSONB、全文搜索等）
- ✅ 扩展性强
- ✅ 数据完整性好
- ✅ 适合复杂查询和分析

### MySQL 优势
- ✅ 简单易用
- ✅ 性能优秀（简单查询）
- ✅ 生态丰富
- ✅ 文档完善
- ✅ 云服务支持好

### 选择建议

- **新项目**：根据团队熟悉度和需求选择
- **复杂查询**：优先考虑 PostgreSQL
- **简单 Web 应用**：MySQL 可能更合适
- **数据分析**：PostgreSQL 更适合
- **高并发 OLTP**：两者都很好，根据具体情况选择

---

## 相关资源

- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [MySQL 官方文档](https://dev.mysql.com/doc/)
- [PostgreSQL vs MySQL 性能对比](https://www.postgresql.org/about/featurematrix/)
- [MySQL vs PostgreSQL 选择指南](https://www.mysql.com/why-mysql/)

---

*最后更新：2024年*

