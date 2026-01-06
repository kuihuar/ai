# PostgreSQL 数据类型

## 目录
- [数值类型](#数值类型)
- [字符类型](#字符类型)
- [日期时间类型](#日期时间类型)
- [布尔类型](#布尔类型)
- [数组类型](#数组类型)
- [JSON 类型](#json-类型)
- [网络地址类型](#网络地址类型)
- [几何类型](#几何类型)
- [UUID 类型](#uuid-类型)
- [自定义类型](#自定义类型)

---

## 数值类型

### 整数类型

| 类型 | 存储大小 | 范围 |
|------|---------|------|
| `SMALLINT` | 2 字节 | -32,768 到 32,767 |
| `INTEGER` 或 `INT` | 4 字节 | -2,147,483,648 到 2,147,483,647 |
| `BIGINT` | 8 字节 | -9,223,372,036,854,775,808 到 9,223,372,036,854,775,807 |
| `SERIAL` | 4 字节 | 1 到 2,147,483,647 (自增) |
| `BIGSERIAL` | 8 字节 | 1 到 9,223,372,036,854,775,807 (自增) |

```sql
-- 使用示例
CREATE TABLE numbers (
    id SERIAL PRIMARY KEY,
    small_num SMALLINT,
    normal_num INTEGER,
    big_num BIGINT
);

INSERT INTO numbers (small_num, normal_num, big_num) VALUES
    (100, 1000000, 1000000000000);
```

### 精确数值类型

```sql
-- DECIMAL 和 NUMERIC 是等价的
-- DECIMAL(precision, scale)
-- precision: 总位数
-- scale: 小数位数

CREATE TABLE prices (
    id SERIAL PRIMARY KEY,
    price DECIMAL(10, 2),      -- 总共10位，2位小数
    discount NUMERIC(5, 2)     -- 总共5位，2位小数
);

INSERT INTO prices (price, discount) VALUES
    (999.99, 0.15),
    (1234.56, 0.20);
```

### 浮点类型

```sql
-- REAL: 单精度浮点数 (4字节)
-- DOUBLE PRECISION: 双精度浮点数 (8字节)

CREATE TABLE measurements (
    id SERIAL PRIMARY KEY,
    temperature REAL,
    precision_value DOUBLE PRECISION
);

INSERT INTO measurements (temperature, precision_value) VALUES
    (36.5, 3.141592653589793);
```

---

## 字符类型

### 基本字符类型

```sql
-- CHAR(n): 固定长度，不足补空格
-- VARCHAR(n): 可变长度，最大n字符
-- TEXT: 无长度限制

CREATE TABLE text_examples (
    id SERIAL PRIMARY KEY,
    fixed_char CHAR(10),        -- 固定10字符
    var_char VARCHAR(100),     -- 最多100字符
    long_text TEXT              -- 无限制
);

INSERT INTO text_examples (fixed_char, var_char, long_text) VALUES
    ('hello', 'world', 'This is a very long text...');
```

### 字符类型选择建议

- **CHAR(n)**: 固定长度数据（如状态码、代码）
- **VARCHAR(n)**: 可变长度，有最大限制（如用户名、邮箱）
- **TEXT**: 长文本（如文章内容、描述）

### 字符串函数

```sql
-- 长度
SELECT LENGTH('hello');                    -- 5
SELECT CHAR_LENGTH('hello');               -- 5

-- 连接
SELECT 'hello' || ' ' || 'world';          -- 'hello world'
SELECT CONCAT('hello', ' ', 'world');      -- 'hello world'

-- 大小写转换
SELECT UPPER('hello');                     -- 'HELLO'
SELECT LOWER('HELLO');                     -- 'hello'

-- 子字符串
SELECT SUBSTRING('hello world', 1, 5);     -- 'hello'
SELECT SUBSTR('hello world', 7);          -- 'world'

-- 替换
SELECT REPLACE('hello world', 'world', 'PostgreSQL');  -- 'hello PostgreSQL'

-- 去除空格
SELECT TRIM('  hello  ');                  -- 'hello'
SELECT LTRIM('  hello');                   -- 'hello'
SELECT RTRIM('hello  ');                   -- 'hello'

-- 模式匹配
SELECT 'hello' LIKE 'he%';                 -- true
SELECT 'hello' SIMILAR TO 'he%';           -- true
SELECT 'hello' ~ '^he';                    -- true (正则表达式)
```

---

## 日期时间类型

### 日期时间类型

| 类型 | 存储大小 | 描述 |
|------|---------|------|
| `DATE` | 4 字节 | 日期（年月日） |
| `TIME` | 8 字节 | 时间（时分秒） |
| `TIMESTAMP` | 8 字节 | 日期和时间（无时区） |
| `TIMESTAMPTZ` | 8 字节 | 日期和时间（带时区） |
| `INTERVAL` | 16 字节 | 时间间隔 |

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_date DATE,
    event_time TIME,
    created_at TIMESTAMP,
    updated_at TIMESTAMPTZ,
    duration INTERVAL
);

INSERT INTO events (event_date, event_time, created_at, updated_at, duration) VALUES
    ('2024-01-15', '14:30:00', 
     '2024-01-15 14:30:00', 
     '2024-01-15 14:30:00+08',
     '2 hours 30 minutes');
```

### 日期时间函数

```sql
-- 当前时间
SELECT CURRENT_DATE;           -- 当前日期
SELECT CURRENT_TIME;           -- 当前时间
SELECT CURRENT_TIMESTAMP;      -- 当前时间戳
SELECT NOW();                  -- 当前时间戳（带时区）

-- 提取部分
SELECT EXTRACT(YEAR FROM NOW());           -- 2024
SELECT EXTRACT(MONTH FROM NOW());          -- 1
SELECT EXTRACT(DAY FROM NOW());            -- 15
SELECT EXTRACT(HOUR FROM NOW());           -- 14
SELECT EXTRACT(DOW FROM NOW());            -- 星期几 (0=周日)

-- 日期运算
SELECT NOW() + INTERVAL '1 day';          -- 加1天
SELECT NOW() - INTERVAL '1 week';        -- 减1周
SELECT NOW() + INTERVAL '2 hours';        -- 加2小时

-- 日期格式化
SELECT TO_CHAR(NOW(), 'YYYY-MM-DD HH24:MI:SS');
SELECT TO_CHAR(NOW(), 'Day, Month DD, YYYY');

-- 日期解析
SELECT TO_DATE('2024-01-15', 'YYYY-MM-DD');
SELECT TO_TIMESTAMP('2024-01-15 14:30:00', 'YYYY-MM-DD HH24:MI:SS');

-- 年龄计算
SELECT AGE('2000-01-01'::DATE);           -- 计算年龄
SELECT AGE('2024-01-15'::DATE, '2000-01-01'::DATE);
```

---

## 布尔类型

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100),
    completed BOOLEAN DEFAULT FALSE
);

INSERT INTO tasks (title, completed) VALUES
    ('Task 1', TRUE),
    ('Task 2', FALSE),
    ('Task 3', NULL);

-- 布尔运算
SELECT * FROM tasks WHERE completed = TRUE;
SELECT * FROM tasks WHERE completed IS TRUE;
SELECT * FROM tasks WHERE NOT completed;
```

布尔值可以表示为：
- `TRUE` / `true` / `t` / `yes` / `y` / `1` / `on`
- `FALSE` / `false` / `f` / `no` / `n` / `0` / `off`
- `NULL`

---

## 数组类型

PostgreSQL 支持数组类型，任何数据类型都可以创建数组。

```sql
-- 创建数组列
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    tags TEXT[],                    -- 文本数组
    prices INTEGER[],               -- 整数数组
    dimensions INTEGER[3]           -- 固定长度数组
);

-- 插入数组数据
INSERT INTO products (name, tags, prices, dimensions) VALUES
    ('Laptop', 
     ARRAY['electronics', 'computer', 'portable'],
     ARRAY[999, 899, 799],
     ARRAY[30, 20, 5]);

-- 或使用字面量语法
INSERT INTO products (name, tags) VALUES
    ('Mouse', '{"wireless", "bluetooth", "ergonomic"}');

-- 查询数组
SELECT * FROM products WHERE 'electronics' = ANY(tags);
SELECT * FROM products WHERE tags @> ARRAY['electronics'];
SELECT * FROM products WHERE tags && ARRAY['computer', 'phone'];

-- 数组函数
SELECT array_length(tags, 1) FROM products;        -- 数组长度
SELECT unnest(tags) FROM products;                 -- 展开数组
SELECT array_append(tags, 'newtag') FROM products; -- 追加元素
SELECT array_prepend('newtag', tags) FROM products; -- 前置元素
SELECT array_remove(tags, 'electronics') FROM products; -- 移除元素

-- 多维数组
CREATE TABLE matrix (
    id SERIAL PRIMARY KEY,
    data INTEGER[][]
);

INSERT INTO matrix (data) VALUES
    (ARRAY[[1, 2, 3], [4, 5, 6], [7, 8, 9]]);
```

---

## JSON 类型

PostgreSQL 提供两种 JSON 类型：`JSON` 和 `JSONB`。

### JSON vs JSONB

- **JSON**: 存储原始文本，保留空格和键顺序，查询较慢
- **JSONB**: 二进制格式，不保留空格和键顺序，查询更快，支持索引

```sql
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    metadata JSON,          -- JSON 类型
    data JSONB              -- JSONB 类型（推荐）
);

INSERT INTO documents (metadata, data) VALUES
    ('{"name": "John", "age": 30, "city": "New York"}',
     '{"name": "John", "age": 30, "city": "New York"}');
```

### JSON 操作符

```sql
-- 访问操作符
SELECT data->'name' FROM documents;              -- 获取 JSON 对象字段（返回 JSON）
SELECT data->>'name' FROM documents;             -- 获取 JSON 对象字段（返回文本）
SELECT data->'tags'->0 FROM documents;           -- 获取数组元素
SELECT data->'tags'->>0 FROM documents;          -- 获取数组元素（文本）

-- 路径操作符
SELECT data#>'{tags,0}' FROM documents;          -- 路径访问
SELECT data#>>'{tags,0}' FROM documents;         -- 路径访问（文本）

-- 包含操作符
SELECT * FROM documents WHERE data @> '{"name": "John"}';
SELECT * FROM documents WHERE data ? 'name';      -- 键是否存在
SELECT * FROM documents WHERE data ?| ARRAY['name', 'age']; -- 任一键存在
SELECT * FROM documents WHERE data ?& ARRAY['name', 'age']; -- 所有键存在
```

### JSON 函数

```sql
-- 构建 JSON
SELECT json_build_object('name', 'John', 'age', 30);
SELECT json_build_array(1, 2, 3);

-- 转换
SELECT to_jsonb('hello'::text);
SELECT to_jsonb(123);

-- 提取
SELECT jsonb_extract_path(data, 'name') FROM documents;
SELECT jsonb_extract_path_text(data, 'name') FROM documents;

-- 键和值
SELECT jsonb_object_keys(data) FROM documents;
SELECT jsonb_each(data) FROM documents;

-- 类型检查
SELECT jsonb_typeof(data->'age') FROM documents;  -- 'number'
SELECT jsonb_typeof(data->'tags') FROM documents; -- 'array'

-- 更新（PostgreSQL 9.5+）
UPDATE documents 
SET data = jsonb_set(data, '{age}', '31') 
WHERE id = 1;

-- 删除键
UPDATE documents 
SET data = data - 'age' 
WHERE id = 1;
```

### JSON 索引

```sql
-- GIN 索引（推荐）
CREATE INDEX idx_documents_data ON documents USING GIN (data);

-- 表达式索引
CREATE INDEX idx_documents_name ON documents ((data->>'name'));

-- 查询使用索引
SELECT * FROM documents WHERE data @> '{"name": "John"}';
```

---

## 网络地址类型

```sql
CREATE TABLE network_info (
    id SERIAL PRIMARY KEY,
    ip_address INET,        -- IPv4 或 IPv6
    mac_address MACADDR,    -- MAC 地址
    cidr_block CIDR         -- 网络地址块
);

INSERT INTO network_info (ip_address, mac_address, cidr_block) VALUES
    ('192.168.1.100', '08:00:2b:01:02:03', '192.168.1.0/24');

-- 网络操作
SELECT * FROM network_info WHERE ip_address << '192.168.1.0/24';  -- 包含在
SELECT * FROM network_info WHERE ip_address = '192.168.1.100';
```

---

## 几何类型

需要启用 PostGIS 扩展（见扩展章节）。

```sql
-- 基本几何类型
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    point POINT,            -- 点
    line LINE,              -- 线
    lseg LSEG,              -- 线段
    box BOX,                -- 矩形
    path PATH,              -- 路径
    polygon POLYGON,        -- 多边形
    circle CIRCLE           -- 圆
);

INSERT INTO locations (name, point) VALUES
    ('Location 1', '(10, 20)');
```

---

## UUID 类型

```sql
-- 启用 uuid-ossp 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50)
);

-- 或使用 gen_random_uuid() (PostgreSQL 13+)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50)
);

-- 插入
INSERT INTO users (username) VALUES ('john');
-- id 会自动生成

-- 手动指定 UUID
INSERT INTO users (id, username) VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'alice');
```

---

## 自定义类型

### 创建枚举类型

```sql
-- 创建枚举类型
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'shipped', 'delivered', 'cancelled');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    status order_status DEFAULT 'pending'
);

INSERT INTO orders (status) VALUES ('pending');
INSERT INTO orders (status) VALUES ('shipped');

-- 查询
SELECT * FROM orders WHERE status = 'pending';
```

### 创建复合类型

```sql
-- 创建复合类型
CREATE TYPE address AS (
    street VARCHAR(100),
    city VARCHAR(50),
    zip_code VARCHAR(10)
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    home_address address
);

INSERT INTO customers (name, home_address) VALUES
    ('John', ROW('123 Main St', 'New York', '10001'));

-- 或使用语法
INSERT INTO customers (name, home_address) VALUES
    ('John', ('123 Main St', 'New York', '10001')::address);

-- 访问字段
SELECT name, (home_address).street FROM customers;
SELECT name, home_address.street FROM customers;
```

### 修改和删除类型

```sql
-- 添加枚举值
ALTER TYPE order_status ADD VALUE 'refunded' AFTER 'cancelled';

-- 重命名类型
ALTER TYPE order_status RENAME TO order_status_type;

-- 删除类型（需要先删除依赖）
DROP TYPE order_status CASCADE;
```

---

## 类型转换

```sql
-- 使用 CAST
SELECT CAST('123' AS INTEGER);
SELECT '123'::INTEGER;

-- 使用 :: 操作符
SELECT '2024-01-15'::DATE;
SELECT '14:30:00'::TIME;
SELECT '{"key": "value"}'::JSONB;

-- 使用函数
SELECT to_number('123.45', '999.99');
SELECT to_timestamp('2024-01-15 14:30:00', 'YYYY-MM-DD HH24:MI:SS');

-- 隐式转换
SELECT '123' + 456;  -- 自动转换为数字
```

---

## 最佳实践

1. **选择合适的类型**
   - 使用 `TEXT` 而不是 `VARCHAR(n)` 如果没有长度限制
   - 使用 `JSONB` 而不是 `JSON` 用于频繁查询
   - 使用 `TIMESTAMPTZ` 而不是 `TIMESTAMP` 处理时区

2. **性能考虑**
   - `JSONB` 支持索引，查询更快
   - 数组类型适合固定结构的数据
   - UUID 适合分布式系统

3. **数据完整性**
   - 使用枚举类型限制值范围
   - 使用 CHECK 约束验证数据
   - 使用外键维护引用完整性

---

## MySQL 对比提示

如果你是 MySQL 用户，以下是对比要点：

- **数组类型**：PostgreSQL 原生支持数组，MySQL 需要使用 JSON 或关联表
- **JSON 类型**：PostgreSQL 的 `JSONB` 性能更好，MySQL 5.7+ 支持 `JSON`
- **UUID 类型**：PostgreSQL 原生支持，MySQL 需要使用 `CHAR(36)` 或 `BINARY(16)`
- **自增主键**：PostgreSQL 使用 `SERIAL`，MySQL 使用 `AUTO_INCREMENT`

更多详细对比请参考 [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md)。

## 下一步学习

- [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md)
- [PostgreSQL 存储过程与函数](./postgresql-functions.md)
- [PostgreSQL 扩展与插件](./postgresql-extensions.md)
- [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md)

---

*最后更新：2024年*

