# 订单号生成流程说明

## 核心原则

**如果订单号生成失败，直接返回错误，不继续创建 Saga。**

## 完整流程

```
HTTP Request: POST /api/v1/orders/saga
  └─> OrderService.CreateOrderSaga()
      └─> OrderCreateSaga.Run()
          └─> OrderUsecase.GenerateOrderNo(ctx)  ⚠️ 关键检查点
              └─> number.Generator.Generate(ctx, "ORD")
                  └─> NumberRepo.GetAndIncrement(ctx, "ORD", "20251217")
                      └─> 数据库事务：
                          1. SELECT ... FROM number_sequences WHERE prefix='ORD' AND date='20251217' FOR UPDATE
                          2. 如果不存在，创建新记录 (sequence=0)
                          3. sequence = sequence + 1
                          4. UPDATE number_sequences SET sequence=?, ...
                          5. 提交事务
                      └─> 返回序列号（如果失败，返回错误）
                  └─> 生成订单号: "ORD20251217000001"（如果失败，返回错误）
              └─> 如果失败，直接返回错误 ❌，不继续
          └─> 如果订单号生成成功，继续创建 Saga
              └─> saga.Run(ctx, sagaCtx, steps)
                  └─> 创建 Saga 实例（使用订单号作为 saga_id）
```

## 错误处理

### 1. 订单号生成失败

**位置**: `OrderCreateSaga.Run()` 第 55-59 行

```go
orderNo, err := s.uc.GenerateOrderNo(ctx)
if err != nil {
    s.log.WithContext(ctx).Errorf("Failed to generate order number: %v", err)
    return nil, fmt.Errorf("failed to generate order number: %w", err)  // ✅ 直接返回错误
}
```

**行为**:
- ✅ 如果 `GenerateOrderNo()` 返回错误，直接返回错误
- ✅ 不会继续创建 Saga
- ✅ 不会创建订单

### 2. 序列号获取失败

**位置**: `number.Generator.Generate()` 第 64-68 行

```go
sequence, err := g.repo.GetAndIncrement(ctx, prefix, datePrefix)
if err != nil {
    g.logger.WithContext(ctx).Errorf("Failed to get sequence: prefix=%s, date=%s, error=%v", prefix, datePrefix, err)
    return "", errors.InternalServer("GENERATE_NO_FAILED", "failed to generate number").WithCause(err)  // ✅ 返回错误
}
```

**行为**:
- ✅ 如果 `GetAndIncrement()` 返回错误，直接返回错误
- ✅ 不会生成订单号
- ✅ 错误会向上传播到 `OrderCreateSaga.Run()`

### 3. 数据库事务失败

**位置**: `NumberRepo.GetAndIncrement()` 第 52-104 行

```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // ... 数据库操作 ...
    return nil  // 如果返回错误，事务会回滚
})

if err != nil {
    r.log.WithContext(ctx).Errorf("Failed to get and increment sequence: prefix=%s, date=%s, error=%v", prefix, date, err)
    return 0, errors.InternalServer("GET_SEQUENCE_FAILED", "failed to get sequence").WithCause(err)  // ✅ 返回错误
}
```

**行为**:
- ✅ 如果事务失败，`Transaction()` 会自动回滚
- ✅ 返回错误，不会返回序列号
- ✅ `number_sequences` 表不会有数据（因为事务回滚）

## 问题排查

### 问题：`number_sequences` 表没有数据，但订单号却被生成了

**可能的原因**:

1. **数据库表不存在**
   - `GetAndIncrement()` 会失败
   - 应该返回错误，不会生成订单号
   - **检查**: 查看日志中是否有 `Failed to get sequence` 错误

2. **数据库连接失败**
   - `GetAndIncrement()` 会失败
   - 应该返回错误，不会生成订单号
   - **检查**: 查看日志中是否有数据库连接错误

3. **事务回滚**
   - 如果事务回滚，`number_sequences` 表不会有数据
   - 但 `GetAndIncrement()` 会返回错误
   - 应该不会生成订单号
   - **检查**: 查看日志中是否有事务错误

4. **订单号在其他地方生成**
   - 如果订单号不是通过 `GenerateOrderNo()` 生成的
   - 可能绕过了 `number_sequences` 表
   - **检查**: 搜索代码中是否有其他地方生成订单号

### 验证方法

1. **检查日志**:
   ```bash
   # 搜索订单号生成相关的日志
   grep "Failed to generate order number" logs/
   grep "Failed to get sequence" logs/
   grep "Generated number" logs/
   ```

2. **检查数据库**:
   ```sql
   -- 检查表是否存在
   SHOW TABLES LIKE 'number_sequences';
   
   -- 检查数据
   SELECT * FROM number_sequences WHERE prefix = 'ORD';
   
   -- 检查重复的 saga_id
   SELECT saga_id, COUNT(*) as count 
   FROM saga_instances 
   GROUP BY saga_id 
   HAVING count > 1;
   ```

3. **检查代码**:
   ```bash
   # 搜索所有生成订单号的地方
   grep -r "GenerateOrderNo\|ORD.*Format\|ORD.*Unix" internal/
   ```

## 修复建议

如果发现 `number_sequences` 表没有数据，但订单号却被生成了：

1. **检查日志**: 查看是否有错误被忽略
2. **检查数据库**: 确认表是否存在，连接是否正常
3. **检查代码**: 确认没有其他地方在生成订单号
4. **添加监控**: 监控订单号生成的成功率和失败率

## 关键检查点

✅ **订单号生成失败时，必须直接返回错误，不继续创建 Saga**

当前的实现已经满足这个要求：
- `OrderCreateSaga.Run()` 在订单号生成失败时直接返回错误
- `GenerateOrderNo()` 在生成失败时返回错误
- `GetAndIncrement()` 在数据库操作失败时返回错误

如果仍然出现 `number_sequences` 表没有数据但订单号被生成的情况，需要进一步排查：
1. 是否有其他地方在生成订单号
2. 是否有错误被忽略或捕获
3. 是否有事务回滚但错误未正确传播

