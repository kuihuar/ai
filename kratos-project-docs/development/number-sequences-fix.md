# number_sequences 表无数据问题修复

## 问题描述

订单号已生成，saga_instances 和 saga_steps 表有数据，但 `number_sequences` 表没有数据。说明订单号不是从 `number_sequences` 表生成的。

## 问题分析

### 调用链

```
HTTP Request: POST /api/v1/orders/saga
  └─> OrderService.CreateOrderSaga()
      └─> OrderCreateSaga.Run()
          └─> OrderUsecase.GenerateOrderNo(ctx)
              └─> number.Generator.Generate(ctx, "ORD")
                  └─> NumberRepo.GetAndIncrement(ctx, "ORD", timestampKey)
                      └─> 数据库事务：SELECT ... FROM number_sequences ... FOR UPDATE
                      └─> INSERT/UPDATE number_sequences
```

### 可能的原因

1. **生成器未正确注入**: `orderNoGenerator` 为 nil，但应该有错误日志
2. **错误被忽略**: `GetAndIncrement` 返回错误，但被忽略或处理不当
3. **事务回滚**: 事务执行了，但被回滚了
4. **有其他生成方式**: 有代码直接生成订单号，绕过了生成器

## 修复措施

### 1. 确保所有错误都被正确处理

已添加详细的错误日志和检查：

```go
// internal/data/number_repo.go
func (r *numberRepo) GetAndIncrement(ctx context.Context, prefix, date string) (int64, error) {
    // 检查数据库连接
    if r.db == nil {
        r.log.WithContext(ctx).Errorf("Database connection is nil, cannot get sequence")
        return 0, errors.InternalServer("DB_CONNECTION_NIL", "database connection is not available")
    }

    r.log.WithContext(ctx).Infof("GetAndIncrement called: prefix=%s, date=%s", prefix, date)
    
    // ... 事务逻辑 ...
    
    // 检查更新结果
    if updateResult.RowsAffected == 0 {
        return fmt.Errorf("failed to update sequence: no rows affected")
    }
    
    r.log.WithContext(ctx).Infof("GetAndIncrement completed: prefix=%s, date=%s, sequence=%d", prefix, date, newSequence)
    return newSequence, nil
}
```

### 2. 确保生成器必须调用 GetAndIncrement

```go
// internal/pkg/number/generator.go
func (g *DBGenerator) Generate(ctx context.Context, prefix string) (string, error) {
    // 检查 repo 是否为 nil
    if g.repo == nil {
        g.logger.WithContext(ctx).Errorf("NumberRepo is nil, cannot generate number")
        return "", errors.InternalServer("REPO_NIL", "number repository is not initialized")
    }
    
    // 必须调用 GetAndIncrement，没有 fallback
    sequence, err := g.repo.GetAndIncrement(ctx, prefix, dateKey)
    if err != nil {
        g.logger.WithContext(ctx).Errorf("Failed to get sequence: prefix=%s, dateKey=%s, error=%v", prefix, dateKey, err)
        return "", errors.InternalServer("GENERATE_NO_FAILED", "failed to generate number").WithCause(err)
    }
    
    // 生成编号
    no := fmt.Sprintf("%s%d%06d", prefix, unixTimestamp, sequence)
    return no, nil
}
```

### 3. 确保 GenerateOrderNo 必须使用生成器

```go
// internal/biz/order.go
func (uc *OrderUsecase) GenerateOrderNo(ctx context.Context) (string, error) {
    if uc.orderNoGenerator == nil {
        uc.log.WithContext(ctx).Errorf("OrderNoGenerator is nil, cannot generate order number")
        return "", fmt.Errorf("order number generator is not initialized")
    }
    const orderPrefix = "ORD"
    return uc.orderNoGenerator.Generate(ctx, orderPrefix)
}
```

## 验证步骤

### 1. 检查日志

运行服务并调用 `CreateOrderSaga` API，查看日志中是否有：

1. `GetAndIncrement called: prefix=ORD, date=...` → 说明方法被调用
2. `Transaction started: prefix=ORD, date=...` → 说明事务开始
3. `Sequence updated: prefix=ORD, date=..., oldSequence=..., newSequence=...` → 说明更新成功
4. `GetAndIncrement completed: prefix=ORD, date=..., sequence=...` → 说明方法完成
5. `Generated number: ORD...` → 说明订单号生成成功

### 2. 检查数据库

```sql
-- 检查 number_sequences 表是否有数据
SELECT * FROM number_sequences WHERE prefix = 'ORD' ORDER BY updated_at DESC LIMIT 10;

-- 检查订单表中的订单号格式
SELECT order_no, created_at FROM orders ORDER BY created_at DESC LIMIT 10;
```

### 3. 检查错误

如果 `number_sequences` 表仍然没有数据，检查日志中是否有错误：

- `Database connection is nil` → 数据库连接问题
- `Failed to query sequence` → 查询失败
- `Failed to update sequence` → 更新失败
- `Failed to get sequence` → 整体失败

## 关键点

1. **没有 fallback 逻辑**: 如果 `GetAndIncrement` 失败，直接返回错误，不会继续生成订单号
2. **所有错误都被记录**: 每个关键步骤都有日志记录
3. **事务保证原子性**: 使用数据库事务和 `FOR UPDATE` 保证序列号的唯一性
4. **必须通过 number_sequences**: 所有订单号生成都必须通过 `number_sequences` 表，没有其他路径

## 如果问题仍然存在

如果修复后 `number_sequences` 表仍然没有数据，可能的原因：

1. **数据库表不存在**: 运行 `./sre-client migrate` 创建表
2. **数据库连接问题**: 检查数据库连接配置
3. **事务隔离级别问题**: 检查数据库事务隔离级别
4. **权限问题**: 检查数据库用户是否有 INSERT/UPDATE 权限

