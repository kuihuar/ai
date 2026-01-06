# PostgreSQL 基础

## 目录
- [安装与配置](#安装与配置)
- [基本概念](#基本概念)
- [数据库操作](#数据库操作)
- [表操作](#表操作)
- [数据操作](#数据操作)
- [SQL 基础](#sql-基础)
- [常用命令](#常用命令)

---

## 安装与配置

### Linux 安装

#### Ubuntu/Debian
```bash
# 更新包列表
sudo apt update

# 安装 PostgreSQL
sudo apt install postgresql postgresql-contrib

# 启动服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 检查状态
sudo systemctl status postgresql
```

#### CentOS/RHEL
```bash
# 安装 PostgreSQL 仓库
sudo yum install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm

# 安装 PostgreSQL
sudo yum install -y postgresql14-server postgresql14

# 初始化数据库
sudo /usr/pgsql-14/bin/postgresql-14-setup initdb

# 启动服务
sudo systemctl start postgresql-14
sudo systemctl enable postgresql-14
```

#### 使用 Docker
```bash
# 拉取镜像
docker pull postgres:14

# 运行容器
docker run --name postgres \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -d postgres:14

# 连接
docker exec -it postgres psql -U myuser -d mydb
```

### 配置文件

#### postgresql.conf
主要配置项：
```conf
# 连接设置
listen_addresses = 'localhost'          # 监听地址
port = 5432                             # 端口
max_connections = 100                   # 最大连接数

# 内存设置
shared_buffers = 128MB                  # 共享缓冲区
effective_cache_size = 4GB              # 有效缓存大小
work_mem = 4MB                          # 工作内存

# 日志设置
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d.log'
log_statement = 'all'                   # 记录所有SQL
```

#### pg_hba.conf
客户端认证配置：
```
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     peer
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5
```

### 初始设置

```bash
# 切换到 postgres 用户
sudo -u postgres psql

# 创建新用户
CREATE USER myuser WITH PASSWORD 'mypassword';

# 创建数据库
CREATE DATABASE mydb OWNER myuser;

# 授予权限
GRANT ALL PRIVILEGES ON DATABASE mydb TO myuser;

# 退出
\q
```

---

## 基本概念

### 数据库架构

```
PostgreSQL 实例
├── 数据库 (Database)
│   ├── 模式 (Schema)
│   │   ├── 表 (Table)
│   │   ├── 视图 (View)
│   │   ├── 函数 (Function)
│   │   ├── 序列 (Sequence)
│   │   └── 索引 (Index)
│   └── 其他对象
└── 系统目录
```

### 核心概念

1. **数据库 (Database)**
   - 独立的命名空间
   - 包含多个模式

2. **模式 (Schema)**
   - 数据库内的命名空间
   - 默认模式：`public`
   - 用于组织数据库对象

3. **表 (Table)**
   - 存储数据的二维结构
   - 由行和列组成

4. **视图 (View)**
   - 虚拟表
   - 基于查询结果

5. **索引 (Index)**
   - 提高查询性能
   - 不影响数据逻辑

---

## 数据库操作

### 创建数据库

```sql
-- 基本创建
CREATE DATABASE mydb;

-- 指定所有者
CREATE DATABASE mydb OWNER myuser;

-- 指定编码
CREATE DATABASE mydb 
  WITH ENCODING 'UTF8'
  LC_COLLATE='en_US.UTF-8'
  LC_CTYPE='en_US.UTF-8'
  TEMPLATE=template0;
```

### 查看数据库

```sql
-- 列出所有数据库
\l

-- 或使用 SQL
SELECT datname FROM pg_database;

-- 查看数据库详情
\l+ mydb
```

### 连接数据库

```sql
-- 切换数据库
\c mydb

-- 或使用 SQL
\connect mydb
```

### 删除数据库

```sql
-- 删除数据库（需要断开所有连接）
DROP DATABASE mydb;

-- 强制删除（PostgreSQL 9.2+）
DROP DATABASE mydb WITH (FORCE);
```

---

## 表操作

### 创建表

```sql
-- 基本创建
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 带约束
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 使用 IF NOT EXISTS
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2)
);
```

### 数据类型

#### 数值类型
- `SMALLINT` - 2字节整数 (-32768 到 32767)
- `INTEGER` 或 `INT` - 4字节整数
- `BIGINT` - 8字节整数
- `DECIMAL(p,s)` 或 `NUMERIC(p,s)` - 精确数值
- `REAL` - 单精度浮点数
- `DOUBLE PRECISION` - 双精度浮点数
- `SERIAL` - 自增整数（相当于 INTEGER + 序列）

#### 字符类型
- `CHAR(n)` - 固定长度字符串
- `VARCHAR(n)` - 可变长度字符串
- `TEXT` - 无长度限制文本

#### 日期时间类型
- `DATE` - 日期
- `TIME` - 时间
- `TIMESTAMP` - 日期和时间
- `TIMESTAMPTZ` - 带时区的时间戳
- `INTERVAL` - 时间间隔

#### 布尔类型
- `BOOLEAN` 或 `BOOL` - true/false

### 修改表

```sql
-- 添加列
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- 删除列
ALTER TABLE users DROP COLUMN phone;

-- 修改列类型
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(200);

-- 添加约束
ALTER TABLE users ADD CONSTRAINT email_unique UNIQUE (email);

-- 删除约束
ALTER TABLE users DROP CONSTRAINT email_unique;

-- 重命名表
ALTER TABLE users RENAME TO customers;

-- 重命名列
ALTER TABLE users RENAME COLUMN username TO name;
```

### 查看表结构

```sql
-- 查看表结构
\d users

-- 查看详细信息
\d+ users

-- 列出所有表
\dt

-- 列出所有表（包括系统表）
\dt+

-- 使用 SQL 查询
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_name = 'users';
```

### 删除表

```sql
-- 删除表
DROP TABLE users;

-- 级联删除（删除依赖对象）
DROP TABLE users CASCADE;

-- 使用 IF EXISTS
DROP TABLE IF EXISTS users;
```

---

## 数据操作

### 插入数据

```sql
-- 基本插入
INSERT INTO users (username, email) 
VALUES ('john', 'john@example.com');

-- 插入多行
INSERT INTO users (username, email) VALUES
    ('alice', 'alice@example.com'),
    ('bob', 'bob@example.com'),
    ('charlie', 'charlie@example.com');

-- 使用默认值
INSERT INTO users (username, email) 
VALUES ('david', 'david@example.com');
-- created_at 会自动使用默认值

-- 从查询插入
INSERT INTO users (username, email)
SELECT name, email FROM old_users;
```

### 查询数据

```sql
-- 基本查询
SELECT * FROM users;

-- 选择特定列
SELECT id, username, email FROM users;

-- 条件查询
SELECT * FROM users WHERE id = 1;
SELECT * FROM users WHERE email LIKE '%@example.com';

-- 排序
SELECT * FROM users ORDER BY created_at DESC;

-- 限制结果
SELECT * FROM users LIMIT 10;
SELECT * FROM users LIMIT 10 OFFSET 20;

-- 去重
SELECT DISTINCT email FROM users;

-- 聚合函数
SELECT COUNT(*) FROM users;
SELECT AVG(amount) FROM orders;
SELECT MAX(created_at) FROM users;
SELECT MIN(price) FROM products;
```

### 更新数据

```sql
-- 基本更新
UPDATE users SET email = 'newemail@example.com' WHERE id = 1;

-- 更新多列
UPDATE users 
SET email = 'newemail@example.com', 
    username = 'newname' 
WHERE id = 1;

-- 使用表达式
UPDATE orders SET amount = amount * 1.1 WHERE status = 'pending';

-- 更新所有行（谨慎使用）
UPDATE users SET status = 'active';
```

### 删除数据

```sql
-- 删除特定行
DELETE FROM users WHERE id = 1;

-- 删除所有行（谨慎使用）
DELETE FROM users;

-- 使用 TRUNCATE（更快，但不可回滚）
TRUNCATE TABLE users;

-- 重置序列
TRUNCATE TABLE users RESTART IDENTITY;
```

---

## SQL 基础

### 连接查询

```sql
-- 内连接
SELECT u.username, o.amount
FROM users u
INNER JOIN orders o ON u.id = o.user_id;

-- 左连接
SELECT u.username, o.amount
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;

-- 右连接
SELECT u.username, o.amount
FROM users u
RIGHT JOIN orders o ON u.id = o.user_id;

-- 全外连接
SELECT u.username, o.amount
FROM users u
FULL OUTER JOIN orders o ON u.id = o.user_id;

-- 交叉连接
SELECT * FROM users CROSS JOIN orders;
```

### 子查询

```sql
-- 标量子查询
SELECT username, 
       (SELECT COUNT(*) FROM orders WHERE user_id = users.id) as order_count
FROM users;

-- EXISTS 子查询
SELECT * FROM users
WHERE EXISTS (
    SELECT 1 FROM orders WHERE orders.user_id = users.id
);

-- IN 子查询
SELECT * FROM users
WHERE id IN (SELECT user_id FROM orders WHERE amount > 100);

-- 相关子查询
SELECT * FROM users u
WHERE (SELECT COUNT(*) FROM orders o WHERE o.user_id = u.id) > 5;
```

### 分组和聚合

```sql
-- GROUP BY
SELECT user_id, COUNT(*) as order_count, SUM(amount) as total_amount
FROM orders
GROUP BY user_id;

-- HAVING（过滤分组）
SELECT user_id, COUNT(*) as order_count
FROM orders
GROUP BY user_id
HAVING COUNT(*) > 5;

-- 多列分组
SELECT user_id, status, COUNT(*) as count
FROM orders
GROUP BY user_id, status;
```

### 窗口函数

```sql
-- ROW_NUMBER
SELECT id, username, 
       ROW_NUMBER() OVER (ORDER BY created_at) as row_num
FROM users;

-- RANK
SELECT id, username, amount,
       RANK() OVER (ORDER BY amount DESC) as rank
FROM orders;

-- 分区窗口
SELECT user_id, amount,
       SUM(amount) OVER (PARTITION BY user_id) as user_total
FROM orders;

-- 移动平均
SELECT date, amount,
       AVG(amount) OVER (ORDER BY date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as moving_avg
FROM daily_sales;
```

---

## 常用命令

### psql 命令

```bash
# 连接数据库
psql -U username -d database_name -h hostname -p 5432

# 执行 SQL 文件
psql -U username -d database_name -f script.sql

# 执行 SQL 命令
psql -U username -d database_name -c "SELECT * FROM users;"
```

### psql 内部命令

```sql
-- 帮助
\?          -- 所有命令帮助
\h          -- SQL 命令帮助
\h SELECT   -- 特定命令帮助

-- 数据库操作
\l          -- 列出数据库
\c dbname   -- 连接数据库
\dt         -- 列出表
\d table    -- 描述表结构
\du         -- 列出用户
\dn         -- 列出模式

-- 执行控制
\q          -- 退出
\i file.sql -- 执行文件
\o file.txt -- 输出到文件
\timing     -- 显示执行时间
\echo text  -- 输出文本

-- 格式化
\x          -- 切换扩展显示
\pset       -- 设置输出格式
```

### 系统信息查询

```sql
-- 查看版本
SELECT version();

-- 查看当前数据库
SELECT current_database();

-- 查看当前用户
SELECT current_user;

-- 查看所有表
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public';

-- 查看表大小
SELECT pg_size_pretty(pg_total_relation_size('users'));

-- 查看数据库大小
SELECT pg_size_pretty(pg_database_size('mydb'));
```

---

## 实践练习

### 练习 1：创建电商数据库

```sql
-- 创建用户表
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建商品表
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    stock INTEGER DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建订单表
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建订单项表
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 练习 2：插入和查询数据

```sql
-- 插入测试数据
INSERT INTO users (username, email, password_hash) VALUES
    ('alice', 'alice@example.com', 'hash1'),
    ('bob', 'bob@example.com', 'hash2');

INSERT INTO products (name, description, price, stock) VALUES
    ('Laptop', 'High performance laptop', 999.99, 10),
    ('Mouse', 'Wireless mouse', 29.99, 50);

-- 查询练习
-- 1. 查找所有用户
SELECT * FROM users;

-- 2. 查找价格大于 100 的商品
SELECT * FROM products WHERE price > 100;

-- 3. 统计每个用户的订单数
SELECT u.username, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.username;
```

---

## 下一步学习

- [PostgreSQL 数据类型](./postgresql-data-types.md) - 深入学习各种数据类型
- [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md) - 性能优化
- [PostgreSQL 事务与并发控制](./postgresql-transactions.md) - 事务管理
- [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md) - 与 MySQL 的详细对比

---

## MySQL 对比提示

如果你是 MySQL 用户，以下是对比要点：

- **自增主键**：PostgreSQL 使用 `SERIAL`，MySQL 使用 `AUTO_INCREMENT`
- **字符串连接**：PostgreSQL 使用 `||`，MySQL 使用 `CONCAT()`
- **LIMIT 语法**：PostgreSQL 使用 `LIMIT count OFFSET offset`，MySQL 支持 `LIMIT offset, count`
- **模式概念**：PostgreSQL 有 Database + Schema，MySQL 只有 Database

更多详细对比请参考 [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md)。

---

*最后更新：2024年*

