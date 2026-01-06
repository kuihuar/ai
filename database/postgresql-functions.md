# PostgreSQL 存储过程与函数

## 目录
- [函数基础](#函数基础)
- [PL/pgSQL 语言](#plpgsql-语言)
- [函数类型](#函数类型)
- [触发器](#触发器)
- [游标](#游标)
- [异常处理](#异常处理)
- [最佳实践](#最佳实践)

---

## 函数基础

### 创建函数

```sql
-- 基本语法
CREATE OR REPLACE FUNCTION function_name(parameters)
RETURNS return_type AS $$
    -- 函数体
$$ LANGUAGE language_name;
```

### 简单示例

```sql
-- SQL 函数
CREATE OR REPLACE FUNCTION add_numbers(a INTEGER, b INTEGER)
RETURNS INTEGER AS $$
    SELECT a + b;
$$ LANGUAGE SQL;

-- 调用
SELECT add_numbers(10, 20);  -- 返回 30
```

### 函数特性

```sql
-- 带默认参数
CREATE OR REPLACE FUNCTION greet(name TEXT, greeting TEXT DEFAULT 'Hello')
RETURNS TEXT AS $$
    SELECT greeting || ', ' || name || '!';
$$ LANGUAGE SQL;

SELECT greet('John');                    -- 'Hello, John!'
SELECT greet('John', 'Hi');              -- 'Hi, John!'

-- 命名参数
SELECT greet(greeting => 'Hi', name => 'John');

-- 可变参数
CREATE OR REPLACE FUNCTION sum_numbers(VARIADIC numbers INTEGER[])
RETURNS INTEGER AS $$
    SELECT SUM(n) FROM unnest(numbers) AS n;
$$ LANGUAGE SQL;

SELECT sum_numbers(1, 2, 3, 4, 5);  -- 返回 15
```

---

## PL/pgSQL 语言

### PL/pgSQL 基础

PL/pgSQL 是 PostgreSQL 的过程化语言，类似于 Oracle 的 PL/SQL。

```sql
-- 启用 PL/pgSQL（通常已默认启用）
CREATE EXTENSION IF NOT EXISTS plpgsql;

-- 基本函数结构
CREATE OR REPLACE FUNCTION get_user_name(user_id INTEGER)
RETURNS TEXT AS $$
DECLARE
    user_name TEXT;
BEGIN
    SELECT username INTO user_name
    FROM users
    WHERE id = user_id;
    
    RETURN user_name;
END;
$$ LANGUAGE plpgsql;
```

### 变量声明

```sql
CREATE OR REPLACE FUNCTION example()
RETURNS VOID AS $$
DECLARE
    -- 基本变量
    counter INTEGER := 0;
    user_name TEXT;
    total_amount DECIMAL(10, 2);
    
    -- 使用 %TYPE（推荐）
    user_id users.id%TYPE;
    user_email users.email%TYPE;
    
    -- 使用 %ROWTYPE
    user_record users%ROWTYPE;
    
    -- 常量
    MAX_COUNT CONSTANT INTEGER := 100;
BEGIN
    -- 函数体
END;
$$ LANGUAGE plpgsql;
```

### 控制结构

#### IF 语句

```sql
CREATE OR REPLACE FUNCTION check_balance(account_id INTEGER, amount DECIMAL)
RETURNS BOOLEAN AS $$
DECLARE
    current_balance DECIMAL;
BEGIN
    SELECT balance INTO current_balance
    FROM accounts
    WHERE id = account_id;
    
    IF current_balance IS NULL THEN
        RAISE EXCEPTION 'Account not found';
    ELSIF current_balance < amount THEN
        RETURN FALSE;
    ELSE
        RETURN TRUE;
    END IF;
END;
$$ LANGUAGE plpgsql;
```

#### CASE 语句

```sql
CREATE OR REPLACE FUNCTION get_status_text(status_code INTEGER)
RETURNS TEXT AS $$
BEGIN
    CASE status_code
        WHEN 1 THEN RETURN 'Pending';
        WHEN 2 THEN RETURN 'Processing';
        WHEN 3 THEN RETURN 'Completed';
        WHEN 4 THEN RETURN 'Cancelled';
        ELSE RETURN 'Unknown';
    END CASE;
END;
$$ LANGUAGE plpgsql;
```

#### 循环

```sql
-- LOOP
CREATE OR REPLACE FUNCTION factorial(n INTEGER)
RETURNS BIGINT AS $$
DECLARE
    result BIGINT := 1;
    i INTEGER := 1;
BEGIN
    LOOP
        EXIT WHEN i > n;
        result := result * i;
        i := i + 1;
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- WHILE
CREATE OR REPLACE FUNCTION countdown(n INTEGER)
RETURNS VOID AS $$
DECLARE
    i INTEGER := n;
BEGIN
    WHILE i > 0 LOOP
        RAISE NOTICE 'Count: %', i;
        i := i - 1;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- FOR
CREATE OR REPLACE FUNCTION sum_range(start_num INTEGER, end_num INTEGER)
RETURNS INTEGER AS $$
DECLARE
    total INTEGER := 0;
BEGIN
    FOR i IN start_num..end_num LOOP
        total := total + i;
    END LOOP;
    RETURN total;
END;
$$ LANGUAGE plpgsql;

-- FOR IN SELECT
CREATE OR REPLACE FUNCTION process_users()
RETURNS VOID AS $$
DECLARE
    user_record users%ROWTYPE;
BEGIN
    FOR user_record IN SELECT * FROM users LOOP
        RAISE NOTICE 'Processing user: %', user_record.username;
        -- 处理逻辑
    END LOOP;
END;
$$ LANGUAGE plpgsql;
```

---

## 函数类型

### 返回标量值

```sql
CREATE OR REPLACE FUNCTION get_user_count()
RETURNS INTEGER AS $$
DECLARE
    count INTEGER;
BEGIN
    SELECT COUNT(*) INTO count FROM users;
    RETURN count;
END;
$$ LANGUAGE plpgsql;
```

### 返回表

```sql
-- 返回表（RETURNS TABLE）
CREATE OR REPLACE FUNCTION get_active_users()
RETURNS TABLE (
    id INTEGER,
    username TEXT,
    email TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.username, u.email
    FROM users u
    WHERE u.status = 'active';
END;
$$ LANGUAGE plpgsql;

-- 使用
SELECT * FROM get_active_users();
```

### 返回集合

```sql
-- 返回 SETOF
CREATE OR REPLACE FUNCTION get_user_emails()
RETURNS SETOF TEXT AS $$
BEGIN
    RETURN QUERY
    SELECT email FROM users;
END;
$$ LANGUAGE plpgsql;

-- 使用
SELECT * FROM get_user_emails();
```

### 返回记录

```sql
CREATE OR REPLACE FUNCTION get_user_by_id(user_id INTEGER)
RETURNS users AS $$
DECLARE
    user_record users%ROWTYPE;
BEGIN
    SELECT * INTO user_record
    FROM users
    WHERE id = user_id;
    
    RETURN user_record;
END;
$$ LANGUAGE plpgsql;

-- 使用
SELECT * FROM get_user_by_id(1);
```

### 返回 JSON

```sql
CREATE OR REPLACE FUNCTION get_user_json(user_id INTEGER)
RETURNS JSONB AS $$
DECLARE
    result JSONB;
BEGIN
    SELECT to_jsonb(u.*) INTO result
    FROM users u
    WHERE u.id = user_id;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;
```

---

## 触发器

### 什么是触发器

触发器是在特定事件（INSERT、UPDATE、DELETE）发生时自动执行的函数。

### 创建触发器函数

```sql
-- 触发器函数示例：自动更新时间戳
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
```

### 触发器类型

#### BEFORE 触发器

```sql
-- 在操作执行前触发
CREATE OR REPLACE FUNCTION validate_user_email()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$' THEN
        RAISE EXCEPTION 'Invalid email format';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_email_trigger
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION validate_user_email();
```

#### AFTER 触发器

```sql
-- 在操作执行后触发
CREATE OR REPLACE FUNCTION log_user_changes()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_audit_log (user_id, action, changed_at)
    VALUES (NEW.id, TG_OP, CURRENT_TIMESTAMP);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_user_changes_trigger
    AFTER INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION log_user_changes();
```

#### INSTEAD OF 触发器

```sql
-- 用于视图，替代实际操作
CREATE VIEW active_users_view AS
SELECT id, username, email FROM users WHERE status = 'active';

CREATE OR REPLACE FUNCTION insert_into_active_users()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO users (username, email, status)
    VALUES (NEW.username, NEW.email, 'active');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_active_users_trigger
    INSTEAD OF INSERT ON active_users_view
    FOR EACH ROW
    EXECUTE FUNCTION insert_into_active_users();
```

### 触发器变量

```sql
CREATE OR REPLACE FUNCTION example_trigger()
RETURNS TRIGGER AS $$
BEGIN
    -- TG_OP: 操作类型 (INSERT, UPDATE, DELETE)
    RAISE NOTICE 'Operation: %', TG_OP;
    
    -- TG_TABLE_NAME: 表名
    RAISE NOTICE 'Table: %', TG_TABLE_NAME;
    
    -- TG_WHEN: 触发时机 (BEFORE, AFTER)
    RAISE NOTICE 'When: %', TG_WHEN;
    
    -- OLD: 旧行（UPDATE, DELETE）
    IF TG_OP = 'UPDATE' THEN
        RAISE NOTICE 'Old value: %', OLD.username;
    END IF;
    
    -- NEW: 新行（INSERT, UPDATE）
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        RAISE NOTICE 'New value: %', NEW.username;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### 删除触发器

```sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
```

---

## 游标

### 显式游标

```sql
CREATE OR REPLACE FUNCTION process_large_dataset()
RETURNS VOID AS $$
DECLARE
    user_cursor CURSOR FOR
        SELECT id, username, email FROM users;
    user_record RECORD;
BEGIN
    OPEN user_cursor;
    
    LOOP
        FETCH user_cursor INTO user_record;
        EXIT WHEN NOT FOUND;
        
        -- 处理每一行
        RAISE NOTICE 'Processing user: %', user_record.username;
    END LOOP;
    
    CLOSE user_cursor;
END;
$$ LANGUAGE plpgsql;
```

### 游标参数

```sql
CREATE OR REPLACE FUNCTION process_users_by_status(user_status TEXT)
RETURNS VOID AS $$
DECLARE
    user_cursor CURSOR(status_filter TEXT) FOR
        SELECT * FROM users WHERE status = status_filter;
    user_record users%ROWTYPE;
BEGIN
    OPEN user_cursor(user_status);
    
    LOOP
        FETCH user_cursor INTO user_record;
        EXIT WHEN NOT FOUND;
        
        -- 处理逻辑
        RAISE NOTICE 'User: %', user_record.username;
    END LOOP;
    
    CLOSE user_cursor;
END;
$$ LANGUAGE plpgsql;
```

---

## 异常处理

### 异常处理语法

```sql
CREATE OR REPLACE FUNCTION safe_divide(a NUMERIC, b NUMERIC)
RETURNS NUMERIC AS $$
DECLARE
    result NUMERIC;
BEGIN
    BEGIN
        result := a / b;
        RETURN result;
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

### 异常类型

```sql
CREATE OR REPLACE FUNCTION handle_exceptions()
RETURNS VOID AS $$
BEGIN
    BEGIN
        -- 可能抛出异常的代码
        INSERT INTO users (username) VALUES (NULL);
    EXCEPTION
        WHEN not_null_violation THEN
            RAISE NOTICE 'Not null violation';
        WHEN unique_violation THEN
            RAISE NOTICE 'Unique violation';
        WHEN foreign_key_violation THEN
            RAISE NOTICE 'Foreign key violation';
        WHEN check_violation THEN
            RAISE NOTICE 'Check violation';
        WHEN OTHERS THEN
            RAISE NOTICE 'Unexpected error: %', SQLERRM;
    END;
END;
$$ LANGUAGE plpgsql;
```

### 获取错误信息

```sql
CREATE OR REPLACE FUNCTION get_error_info()
RETURNS TEXT AS $$
DECLARE
    error_message TEXT;
    error_detail TEXT;
    error_hint TEXT;
BEGIN
    BEGIN
        -- 可能出错的代码
        INSERT INTO users (id) VALUES (1);
    EXCEPTION
        WHEN OTHERS THEN
            error_message := SQLERRM;
            error_detail := SQLSTATE;
            error_hint := 'Check the error details';
            
            RAISE NOTICE 'Error: %', error_message;
            RAISE NOTICE 'State: %', error_detail;
            RAISE NOTICE 'Hint: %', error_hint;
            
            RETURN error_message;
    END;
END;
$$ LANGUAGE plpgsql;
```

### 抛出异常

```sql
CREATE OR REPLACE FUNCTION validate_age(age INTEGER)
RETURNS VOID AS $$
BEGIN
    IF age < 0 THEN
        RAISE EXCEPTION 'Age cannot be negative: %', age;
    ELSIF age > 150 THEN
        RAISE EXCEPTION USING
            ERRCODE = 'P0001',
            MESSAGE = 'Age is too large',
            HINT = 'Please check the age value';
    END IF;
END;
$$ LANGUAGE plpgsql;
```

---

## 高级特性

### 函数重载

```sql
-- 同名函数，不同参数
CREATE OR REPLACE FUNCTION add_numbers(a INTEGER, b INTEGER)
RETURNS INTEGER AS $$
    SELECT a + b;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION add_numbers(a NUMERIC, b NUMERIC)
RETURNS NUMERIC AS $$
    SELECT a + b;
$$ LANGUAGE SQL;

-- 调用时会根据参数类型选择
SELECT add_numbers(1, 2);           -- 调用 INTEGER 版本
SELECT add_numbers(1.5, 2.5);       -- 调用 NUMERIC 版本
```

### 函数稳定性

```sql
-- IMMUTABLE: 相同输入总是返回相同结果，不访问数据库
CREATE OR REPLACE FUNCTION calculate_area(radius NUMERIC)
RETURNS NUMERIC AS $$
    SELECT 3.14159 * radius * radius;
$$ LANGUAGE SQL IMMUTABLE;

-- STABLE: 相同输入在同一事务中返回相同结果
CREATE OR REPLACE FUNCTION get_current_user_count()
RETURNS INTEGER AS $$
    SELECT COUNT(*) FROM users;
$$ LANGUAGE SQL STABLE;

-- VOLATILE: 默认，每次调用可能返回不同结果
CREATE OR REPLACE FUNCTION get_random_number()
RETURNS INTEGER AS $$
    SELECT floor(random() * 100)::INTEGER;
$$ LANGUAGE SQL VOLATILE;
```

### 函数安全

```sql
-- SECURITY DEFINER: 以函数所有者权限执行
CREATE OR REPLACE FUNCTION admin_delete_user(user_id INTEGER)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users WHERE id = user_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- SECURITY INVOKER: 以调用者权限执行（默认）
CREATE OR REPLACE FUNCTION user_delete_own_account(user_id INTEGER)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users WHERE id = user_id AND id = current_user_id();
END;
$$ LANGUAGE plpgsql SECURITY INVOKER;
```

---

## 最佳实践

1. **函数设计**
   - 保持函数简洁，单一职责
   - 使用有意义的函数名
   - 添加注释说明

2. **性能优化**
   - 使用 IMMUTABLE 标记纯函数
   - 避免在循环中执行查询
   - 使用批量操作

3. **错误处理**
   - 总是处理可能的异常
   - 提供有意义的错误消息
   - 记录错误日志

4. **安全性**
   - 使用 SECURITY INVOKER 除非必要
   - 验证输入参数
   - 防止 SQL 注入

---

## 下一步学习

- [PostgreSQL 全文搜索](./postgresql-fulltext-search.md)
- [PostgreSQL 扩展与插件](./postgresql-extensions.md)
- [PostgreSQL 管理与运维](./postgresql-admin.md)

---

*最后更新：2024年*

