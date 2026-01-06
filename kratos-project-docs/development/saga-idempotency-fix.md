# Saga ID 重复问题修复说明

## 问题描述

**错误信息**：
```
Duplicate entry 'ORD20251217000001' for key 'saga_instances.saga_id'
```

**问题现象**：
- 创建订单时，`number_sequences` 表中没有数据
- 但是订单号 `ORD20251217000001` 却被生成了
- 这个订单号被用作 `saga_id`，导致重复键冲突

## 问题分析

### 1. 订单号生成流程

```
OrderCreateSaga.Run()
  └─> s.uc.GenerateOrderNo(ctx)
      └─> orderNoGenerator.Generate(ctx, "ORD")
          └─> repo.GetAndIncrement(ctx, "ORD", "20251217")
              └─> 数据库事务：
                  1. SELECT ... FROM number_sequences WHERE prefix='ORD' AND date='20251217' FOR UPDATE
                  2. 如果不存在，创建新记录
                  3. sequence = sequence + 1
                  4. UPDATE number_sequences SET sequence=?, ...
                  5. 返回序列号
          └─> 生成订单号: "ORD20251217000001"
  └─> 使用订单号作为 saga_id
      └─> 创建 Saga 实例
```

### 2. 可能的原因

#### 原因 1: 数据库表不存在
- 如果 `number_sequences` 表不存在，`GetAndIncrement` 会失败
- 但是，如果错误被忽略或处理不当，可能会继续执行

#### 原因 2: 事务回滚
- 如果 `GetAndIncrement` 的事务回滚，`number_sequences` 表没有数据
- 但是序列号已经被读取并用于生成订单号
- 订单号生成成功，但数据库没有记录

#### 原因 3: 并发问题
- 多个请求同时生成订单号
- 第一个请求成功写入 `number_sequences`
- 第二个请求也成功生成订单号（因为序列号已递增）
- 但是两个请求都使用相同的订单号创建 Saga 实例，导致重复

#### 原因 4: 幂等性问题
- 同一个请求被调用两次（网络重试、客户端重试）
- 第一次调用成功创建了 Saga 实例
- 第二次调用尝试创建相同的 Saga 实例，导致重复

## 修复方案

### 1. 添加幂等性检查

在创建 Saga 实例前，先检查是否已存在：

```go
// 先检查是否已存在（幂等性检查）
existing, err := r.data.ent.SagaInstance.Query().
    Where(sagainstance.SagaIDEQ(instance.SagaID)).
    Only(ctx)
if err == nil && existing != nil {
    // 已存在，返回已存在的实例
    return r.sagaInstanceToBiz(existing), nil
}
```

### 2. 处理并发创建

如果创建时遇到唯一约束冲突，重新查询已存在的实例：

```go
entInstance, err := create.Save(ctx)
if err != nil {
    // 检查是否是唯一约束冲突（并发创建）
    if ent.IsConstraintError(err) {
        // 并发创建，重新查询已存在的实例
        existing, queryErr := r.data.ent.SagaInstance.Query().
            Where(sagainstance.SagaIDEQ(instance.SagaID)).
            Only(ctx)
        if queryErr == nil && existing != nil {
            return r.sagaInstanceToBiz(existing), nil
        }
    }
    return nil, err
}
```

### 3. 确保订单号生成失败时不会继续

在 `OrderCreateSaga.Run()` 中，如果订单号生成失败，应该直接返回错误：

```go
orderNo, err := s.uc.GenerateOrderNo(ctx)
if err != nil {
    s.log.WithContext(ctx).Errorf("Failed to generate order number: %v", err)
    return nil, fmt.Errorf("failed to generate order number: %w", err)
}
```

## 验证步骤

### 1. 检查 number_sequences 表

```sql
-- 检查表是否存在
SHOW TABLES LIKE 'number_sequences';

-- 检查表结构
DESC number_sequences;

-- 检查数据
SELECT * FROM number_sequences WHERE prefix = 'ORD';
```

### 2. 检查订单号生成日志

查看日志中是否有：
- `Failed to get sequence` - 序列号获取失败
- `Generated number: ORD...` - 订单号生成成功
- `Failed to generate order number` - 订单号生成失败

### 3. 检查 Saga 实例

```sql
-- 检查重复的 saga_id
SELECT saga_id, COUNT(*) as count 
FROM saga_instances 
GROUP BY saga_id 
HAVING count > 1;

-- 检查特定订单号的 Saga 实例
SELECT * FROM saga_instances WHERE saga_id = 'ORD20251217000001';
```

## 根本原因

**核心问题**：订单号生成和 Saga 实例创建不是原子操作。

**流程问题**：
1. 生成订单号（可能成功，但 `number_sequences` 表没有数据）
2. 使用订单号创建 Saga 实例
3. 如果 Saga 实例创建失败，订单号已经生成，无法回滚

**解决方案**：
1. ✅ 添加幂等性检查（已实现）
2. ✅ 处理并发创建（已实现）
3. ⚠️ 确保订单号生成失败时不会继续（需要验证）

## 后续优化建议

1. **事务一致性**：考虑将订单号生成和 Saga 实例创建放在同一个事务中
2. **重试机制**：对于幂等性操作，支持自动重试
3. **监控告警**：监控 `number_sequences` 表的数据增长，确保正常写入

