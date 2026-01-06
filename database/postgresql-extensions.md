# PostgreSQL 扩展与插件

## 目录
- [扩展基础](#扩展基础)
- [常用扩展](#常用扩展)
- [PostGIS 地理空间扩展](#postgis-地理空间扩展)
- [性能监控扩展](#性能监控扩展)
- [全文搜索扩展](#全文搜索扩展)
- [其他实用扩展](#其他实用扩展)

---

## 扩展基础

### 什么是扩展

扩展（Extension）是 PostgreSQL 的插件系统，可以添加新功能、数据类型、函数等。

### 扩展管理

```sql
-- 查看可用扩展
SELECT * FROM pg_available_extensions;

-- 查看已安装的扩展
SELECT * FROM pg_extension;

-- 安装扩展
CREATE EXTENSION IF NOT EXISTS extension_name;

-- 删除扩展
DROP EXTENSION IF EXISTS extension_name;

-- 查看扩展版本
SELECT extname, extversion FROM pg_extension;
```

---

## 常用扩展

### uuid-ossp（UUID 生成）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 生成 UUID
SELECT uuid_generate_v1();  -- 基于时间戳
SELECT uuid_generate_v4();  -- 随机 UUID

-- 在表中使用
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50)
);
```

### pgcrypto（加密函数）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 哈希
SELECT crypt('password', gen_salt('bf'));  -- Blowfish 加密

-- MD5
SELECT md5('text');

-- SHA256
SELECT encode(digest('text', 'sha256'), 'hex');

-- 加密/解密
SELECT encrypt('secret', 'key', 'aes');
SELECT decrypt(encrypted_data, 'key', 'aes');
```

### hstore（键值对存储）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS hstore;

-- 创建表
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    attributes HSTORE
);

-- 插入数据
INSERT INTO products (name, attributes) VALUES
    ('Laptop', 'color => "black", weight => "2kg", price => "999"');

-- 查询
SELECT * FROM products WHERE attributes->'color' = 'black';
SELECT * FROM products WHERE attributes ? 'price';

-- 更新
UPDATE products 
SET attributes = attributes || 'discount => "10%"'::hstore
WHERE id = 1;
```

---

## PostGIS 地理空间扩展

### 安装 PostGIS

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

-- 查看版本
SELECT PostGIS_Version();
```

### 基本使用

```sql
-- 创建表
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    geom GEOMETRY(POINT, 4326)  -- 点，WGS84 坐标系
);

-- 插入数据
INSERT INTO locations (name, geom) VALUES
    ('Beijing', ST_GeomFromText('POINT(116.4074 39.9042)', 4326)),
    ('Shanghai', ST_GeomFromText('POINT(121.4737 31.2304)', 4326));

-- 计算距离
SELECT 
    a.name,
    b.name,
    ST_Distance(a.geom, b.geom) AS distance
FROM locations a, locations b
WHERE a.id < b.id;

-- 查找附近的点
SELECT name
FROM locations
WHERE ST_DWithin(
    geom,
    ST_GeomFromText('POINT(116.4074 39.9042)', 4326),
    0.1  -- 距离阈值
);
```

---

## 性能监控扩展

### pg_stat_statements

```sql
-- 安装
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

-- 重置统计
SELECT pg_stat_statements_reset();
```

### pg_trgm（相似度搜索）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 创建索引
CREATE INDEX idx_users_username_trgm ON users USING GIN (username gin_trgm_ops);

-- 相似度搜索
SELECT username, similarity(username, 'john') AS sim
FROM users
WHERE username % 'john'  -- 相似度操作符
ORDER BY sim DESC;

-- 使用距离
SELECT username
FROM users
WHERE username <-> 'john' < 0.5;  -- 距离操作符
```

---

## 全文搜索扩展

### zhparser（中文分词）

```bash
# 需要先安装 zhparser（从源码编译）
# 然后创建扩展
```

```sql
-- 创建中文搜索配置
CREATE EXTENSION IF NOT EXISTS zhparser;

CREATE TEXT SEARCH CONFIGURATION chinese_parser (PARSER = zhparser);

-- 使用
SELECT to_tsvector('chinese_parser', 'PostgreSQL 是一个强大的数据库');
```

---

## 其他实用扩展

### pg_buffercache（缓冲区缓存）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS pg_buffercache;

-- 查看缓冲区使用情况
SELECT 
    c.relname,
    count(*) AS buffers
FROM pg_buffercache b
INNER JOIN pg_class c ON b.relfilenode = pg_relation_filenode(c.oid)
GROUP BY c.relname
ORDER BY count(*) DESC;
```

### pg_freespacemap（空闲空间映射）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS pg_freespacemap;

-- 查看表的空闲空间
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(avail) AS free_space
FROM pg_freespacemap_relations
ORDER BY avail DESC;
```

### tablefunc（表函数）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS tablefunc;

-- 交叉表（透视表）
SELECT * FROM crosstab(
    'SELECT user_id, status, COUNT(*) 
     FROM orders 
     GROUP BY user_id, status 
     ORDER BY 1, 2',
    'SELECT DISTINCT status FROM orders ORDER BY 1'
) AS ct(user_id INTEGER, pending BIGINT, completed BIGINT, cancelled BIGINT);
```

### dblink（数据库链接）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS dblink;

-- 连接到远程数据库
SELECT dblink_connect('myconn', 'host=remote_host dbname=mydb user=myuser password=mypass');

-- 查询远程数据库
SELECT * FROM dblink('myconn', 'SELECT * FROM users') AS t(id INTEGER, username TEXT);

-- 断开连接
SELECT dblink_disconnect('myconn');
```

### pg_prewarm（预热）

```sql
-- 安装
CREATE EXTENSION IF NOT EXISTS pg_prewarm;

-- 预热表到缓存
SELECT pg_prewarm('users');

-- 预热索引
SELECT pg_prewarm('idx_users_email');
```

---

## 扩展开发

### 创建简单扩展

```sql
-- 1. 创建扩展目录结构
-- myextension/
--   ├── myextension.control
--   ├── myextension--1.0.sql
--   └── Makefile

-- 2. myextension.control
-- comment = 'My custom extension'
-- default_version = '1.0'
-- module_pathname = '$libdir/myextension'
-- relocatable = true

-- 3. myextension--1.0.sql
-- CREATE FUNCTION my_function() RETURNS TEXT AS $$
--     SELECT 'Hello from extension';
-- $$ LANGUAGE SQL;

-- 4. 安装
-- CREATE EXTENSION myextension;
```

---

## 最佳实践

1. **扩展选择**
   - 使用官方和社区维护的扩展
   - 检查扩展的兼容性和维护状态

2. **性能考虑**
   - 某些扩展可能影响性能
   - 在生产环境前测试

3. **版本管理**
   - 记录使用的扩展版本
   - 在迁移脚本中包含扩展安装

---

## 下一步学习

- [PostgreSQL 全文搜索](./postgresql-fulltext-search.md)
- [PostgreSQL 管理与运维](./postgresql-admin.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)

---

*最后更新：2024年*

