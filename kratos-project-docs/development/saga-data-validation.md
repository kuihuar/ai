# Saga 数据验证指南

本文档说明如何验证 `saga_instances` 和 `saga_steps` 表中的数据是否正确。

## 快速验证

### 方法 1: 使用验证脚本

```bash
# 使用默认配置（test 数据库）
./scripts/validate-saga-data.sh

# 指定数据库名
./scripts/validate-saga-data.sh your_database_name

# 指定数据库连接参数
DB_USER=root DB_PASS=password DB_HOST=localhost DB_PORT=3306 \
  ./scripts/validate-saga-data.sh your_database_name
```

### 方法 2: 使用 SQL 脚本

```bash
# 连接到数据库
mysql -u root -p test < scripts/check-saga-data.sql

# 或者交互式执行
mysql -u root -p test
source scripts/check-saga-data.sql
```

## 验证标准

### 1. 表结构验证

**检查项**:
- ✅ `saga_instances` 表存在
- ✅ `saga_steps` 表存在
- ✅ 表结构符合 Ent Schema 定义

**验证 SQL**:
```sql
SHOW TABLES LIKE 'saga_%';
DESCRIBE saga_instances;
DESCRIBE saga_steps;
```

### 2. 数据完整性验证

#### 2.1 每个 Saga 实例应该有 3 个步骤

**规则**: 
- `step_order = 0`: `create-order`
- `step_order = 1`: `reserve-inventory`
- `step_order = 2`: `freeze-payment`

**验证 SQL**:
```sql
SELECT 
    si.saga_id,
    COUNT(ss.id) AS step_count,
    CASE 
        WHEN COUNT(ss.id) = 3 THEN '✅ 正确'
        ELSE '❌ 错误：应该有 3 个步骤'
    END AS validation
FROM saga_instances si
LEFT JOIN saga_steps ss ON si.saga_id = ss.saga_id
GROUP BY si.id
HAVING COUNT(ss.id) != 3;
```

**预期结果**: 应该返回 0 行（所有实例都有 3 个步骤）

#### 2.2 步骤顺序和名称必须匹配

**验证 SQL**:
```sql
SELECT 
    saga_id,
    step_name,
    step_order,
    CASE 
        WHEN step_order = 0 AND step_name = 'create-order' THEN '✅'
        WHEN step_order = 1 AND step_name = 'reserve-inventory' THEN '✅'
        WHEN step_order = 2 AND step_name = 'freeze-payment' THEN '✅'
        ELSE '❌ 不匹配'
    END AS validation
FROM saga_steps
WHERE NOT (
    (step_order = 0 AND step_name = 'create-order') OR
    (step_order = 1 AND step_name = 'reserve-inventory') OR
    (step_order = 2 AND step_name = 'freeze-payment')
);
```

**预期结果**: 应该返回 0 行

### 3. 状态一致性验证

#### 3.1 COMPLETED 状态的实例

**规则**: 
- Saga 实例状态 = `COMPLETED (3)`
- 所有步骤状态 = `EXECUTED (2)`

**验证 SQL**:
```sql
SELECT 
    si.saga_id,
    si.status AS instance_status,
    COUNT(CASE WHEN ss.status = 2 THEN 1 END) AS executed_steps,
    COUNT(CASE WHEN ss.status != 2 THEN 1 END) AS other_steps
FROM saga_instances si
LEFT JOIN saga_steps ss ON si.saga_id = ss.saga_id
WHERE si.status = 3
GROUP BY si.id
HAVING COUNT(CASE WHEN ss.status != 2 THEN 1 END) > 0;
```

**预期结果**: 应该返回 0 行

#### 3.2 COMPENSED 状态的实例

**规则**: 
- Saga 实例状态 = `COMPENSED (5)`
- 至少有一个步骤状态 = `EXECUTE_FAILED (3)`
- 至少有一个步骤状态 = `COMPENSATED (5)`

**验证 SQL**:
```sql
SELECT 
    si.saga_id,
    COUNT(CASE WHEN ss.status = 3 THEN 1 END) AS failed_steps,
    COUNT(CASE WHEN ss.status = 5 THEN 1 END) AS compensated_steps
FROM saga_instances si
LEFT JOIN saga_steps ss ON si.saga_id = ss.saga_id
WHERE si.status = 5
GROUP BY si.id
HAVING COUNT(CASE WHEN ss.status = 3 THEN 1 END) = 0 
    OR COUNT(CASE WHEN ss.status = 5 THEN 1 END) = 0;
```

**预期结果**: 应该返回 0 行

### 4. 时间逻辑验证

#### 4.1 执行时间应该在补偿时间之前

**规则**: 
- 如果 `compensated_at` 不为空，则 `executed_at` 应该早于 `compensated_at`

**验证 SQL**:
```sql
SELECT 
    saga_id,
    step_name,
    executed_at,
    compensated_at
FROM saga_steps
WHERE compensated_at IS NOT NULL
  AND executed_at IS NOT NULL
  AND executed_at >= compensated_at;
```

**预期结果**: 应该返回 0 行

#### 4.2 完成时间应该在开始时间之后

**验证 SQL**:
```sql
SELECT 
    saga_id,
    started_at,
    completed_at
FROM saga_instances
WHERE completed_at IS NOT NULL
  AND completed_at < started_at;
```

**预期结果**: 应该返回 0 行

### 5. 数据关联验证

#### 5.1 所有步骤都应该有对应的实例

**验证 SQL**:
```sql
SELECT ss.*
FROM saga_steps ss
LEFT JOIN saga_instances si ON ss.saga_id = si.saga_id
WHERE si.saga_id IS NULL;
```

**预期结果**: 应该返回 0 行

#### 5.2 实例和步骤的 saga_id 应该一致

**验证 SQL**:
```sql
SELECT 
    si.saga_id AS instance_saga_id,
    ss.saga_id AS step_saga_id
FROM saga_instances si
INNER JOIN saga_steps ss ON si.saga_id = ss.saga_id
WHERE si.saga_id != ss.saga_id;
```

**预期结果**: 应该返回 0 行（这个查询理论上不应该有结果）

## 状态值说明

### Saga 实例状态 (saga_instances.status)

| 值 | 名称 | 说明 |
|---|------|------|
| 1 | PENDING | 待开始 |
| 2 | RUNNING | 执行中 |
| 3 | COMPLETED | 已完成（所有步骤成功） |
| 4 | FAILED | 失败待补偿 |
| 5 | COMPENSED | 已补偿完成 |

### Saga 步骤状态 (saga_steps.status)

| 值 | 名称 | 说明 |
|---|------|------|
| 0 | NOT_EXECUTED | 未执行 |
| 1 | EXECUTING | 执行中 |
| 2 | EXECUTED | 执行成功 |
| 3 | EXECUTE_FAILED | 执行失败 |
| 4 | COMPENSATING | 补偿中 |
| 5 | COMPENSATED | 补偿成功 |
| 6 | COMPENSATE_FAILED | 补偿失败 |

## 常见问题

### 问题 1: 步骤数量不是 3

**可能原因**:
- Saga 执行过程中断
- 数据库迁移不完整
- 手动插入了数据

**解决方法**:
- 检查是否有未完成的 Saga 实例
- 重新运行数据库迁移
- 清理不完整的数据

### 问题 2: 状态不一致

**可能原因**:
- Saga 执行过程中服务崩溃
- 补偿逻辑有 bug
- 手动修改了数据

**解决方法**:
- 检查日志，查看 Saga 执行历史
- 验证补偿逻辑是否正确
- 修复数据或重新执行 Saga

### 问题 3: 时间顺序错误

**可能原因**:
- 系统时间不同步
- 数据库时区设置错误
- 手动修改了时间戳

**解决方法**:
- 检查系统时间
- 检查数据库时区设置
- 修复时间戳

## 验证清单

使用以下清单快速验证数据：

- [ ] `saga_instances` 表存在
- [ ] `saga_steps` 表存在
- [ ] 每个 Saga 实例都有 3 个步骤
- [ ] 步骤顺序和名称正确
- [ ] COMPLETED 状态的实例所有步骤都是 EXECUTED
- [ ] COMPENSED 状态的实例有失败的步骤和补偿的步骤
- [ ] 执行时间早于补偿时间
- [ ] 所有步骤都有对应的实例
- [ ] 没有孤立的数据

## 示例：正确的数据

### 成功场景的数据

```sql
-- Saga 实例
saga_id: SAGA-1234567890-1
status: 3 (COMPLETED)
started_at: 2024-01-01 10:00:00
completed_at: 2024-01-01 10:00:05

-- 步骤
step_order=0, step_name='create-order', status=2 (EXECUTED), executed_at=2024-01-01 10:00:01
step_order=1, step_name='reserve-inventory', status=2 (EXECUTED), executed_at=2024-01-01 10:00:03
step_order=2, step_name='freeze-payment', status=2 (EXECUTED), executed_at=2024-01-01 10:00:05
```

### 失败补偿场景的数据

```sql
-- Saga 实例
saga_id: SAGA-1234567890-2
status: 5 (COMPENSED)
started_at: 2024-01-01 10:01:00
completed_at: 2024-01-01 10:01:06

-- 步骤
step_order=0, step_name='create-order', status=5 (COMPENSATED), 
  executed_at=2024-01-01 10:01:01, compensated_at=2024-01-01 10:01:05
step_order=1, step_name='reserve-inventory', status=5 (COMPENSATED),
  executed_at=2024-01-01 10:01:03, compensated_at=2024-01-01 10:01:04
step_order=2, step_name='freeze-payment', status=3 (EXECUTE_FAILED),
  executed_at=NULL, compensated_at=NULL
```

## 总结

数据验证的关键点：

1. **完整性**: 每个实例都有 3 个步骤，步骤顺序和名称正确
2. **一致性**: 实例状态和步骤状态匹配
3. **逻辑性**: 时间顺序正确，关联关系正确
4. **可追溯性**: 可以追踪每个 Saga 的完整执行历史

如果所有验证都通过，说明数据是正确的！

