# number_sequences 表 date 字段长度修复

## 问题描述

错误信息：
```
Error 1406 (22001): Data too long for column 'date' at row 1
INSERT INTO `number_sequences` (`prefix`,`date`,`sequence`,`created_at`,`updated_at`) 
VALUES ('ORD','1765947160',0,'2025-12-17 04:52:48.86','2025-12-17 04:52:48.86')
```

## 问题分析

1. **插入的值**: `'1765947160'` - 这是一个10位的 Unix 时间戳
2. **字段长度**: `date` 字段的长度可能不够（可能是 `varchar(8)` 或其他较小的值）
3. **需要长度**: Unix 时间戳（秒）最大是 10 位（2147483647，对应 2038-01-19），所以至少需要 `varchar(10)`
4. **当前设置**: 我们使用每10秒一个区间，所以 `date` 字段存储的是时间戳向下取整到10秒的值，仍然是10位数字

## 修复方案

### 1. 更新 Ent Schema

```go
// internal/data/ent/schema/number_sequence.go
field.String("date").
    MaxLen(20).
    Comment("时间标识（Unix时间戳，每10秒一个区间）").
    SchemaType(map[string]string{
        "mysql": "varchar(20)",
    }),
```

### 2. 更新数据库表结构

需要手动执行 SQL 修改表结构：

```sql
ALTER TABLE number_sequences 
MODIFY COLUMN date VARCHAR(20) NOT NULL COMMENT '时间标识(Unix时间戳，每10秒一个区间)';
```

### 3. 验证修复

```sql
-- 检查表结构
DESC number_sequences;

-- 应该看到 date 字段类型为 varchar(20)
```

## 执行步骤

1. **运行数据库迁移**（如果使用 Ent 迁移）:
   ```bash
   ./sre-client migrate
   ```

2. **或者手动执行 SQL**:
   ```sql
   ALTER TABLE number_sequences 
   MODIFY COLUMN date VARCHAR(20) NOT NULL COMMENT '时间标识(Unix时间戳，每10秒一个区间)';
   ```

3. **重新生成 Ent 代码**（如果需要）:
   ```bash
   make ent-generate
   ```

4. **测试订单号生成**:
   - 调用 `CreateOrderSaga` API
   - 检查 `number_sequences` 表是否有数据
   - 检查日志是否有错误

## 注意事项

- Unix 时间戳（秒）最大是 10 位（2147483647）
- 我们使用 `varchar(20)` 是为了留有余地，未来如果需要更长的标识符也可以支持
- 如果表已经存在数据，修改字段长度不会丢失数据（从短改长是安全的）

