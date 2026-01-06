# number_sequences 表无数据问题详细排查

## 调用链分析

### 完整调用链

```
HTTP Request: POST /api/v1/orders/saga
  └─> OrderService.CreateOrderSaga(ctx, in)
      └─> OrderCreateSaga.Run(ctx, userID, amount, currency, description, items)
          └─> OrderUsecase.GenerateOrderNo(ctx)  ⚠️ 关键点
              └─> number.Generator.Generate(ctx, "ORD")
                  └─> NumberRepo.GetAndIncrement(ctx, "ORD", "20251217")
                      └─> r.db.WithContext(ctx).Transaction(...)  ⚠️ 事务开始
                          └─> SELECT ... FROM number_sequences ... FOR UPDATE
                          └─> INSERT INTO number_sequences ... (如果不存在)
                          └─> UPDATE number_sequences SET sequence=?, ...
                      └─> 事务提交（如果成功）
```

## 可能的问题

### 问题 1: 数据库连接为 nil

**检查点**: `NewNumberRepo` 时，`data.db` 可能为 nil

**验证**:
```go
// internal/data/number_repo.go
func NewNumberRepo(data *Data, logger log.Logger) NumberRepo {
    return &numberRepo{
        db:  data.db,  // 如果 data.db 为 nil，这里就是 nil
        log: log.NewHelper(logger),
    }
}
```

**修复**: 已添加 nil 检查

### 问题 2: 数据库表不存在

**检查点**: `number_sequences` 表可能未创建

**验证**:
```sql
SHOW TABLES LIKE 'number_sequences';
```

**修复**: 运行数据库迁移
```bash
./sre-client migrate
```

### 问题 3: 事务嵌套问题

**检查点**: 如果 `GenerateOrderNo` 在一个事务中被调用，而 `GetAndIncrement` 也开启事务，可能会有问题

**GORM 行为**:
- `WithContext(ctx).Transaction()` 会检查 context 中是否已有事务
- 如果已有事务，会复用现有事务（不会开启新事务）
- 如果现有事务回滚，嵌套的事务也会回滚

**当前实现**:
- `OrderCreateSaga.Run()` 不在事务中
- `GenerateOrderNo()` 不在事务中
- `GetAndIncrement()` 开启独立事务

**结论**: 不应该有事务嵌套问题

### 问题 4: 数据库连接失败但错误被忽略

**检查点**: 如果数据库连接失败，`GetAndIncrement` 应该返回错误

**验证**: 查看日志中是否有错误

## 调试步骤

### 步骤 1: 添加详细日志

已在以下位置添加日志：
1. `GenerateOrderNo()`: 检查生成器是否为 nil
2. `Generate()`: 记录生成开始和结果
3. `GetAndIncrement()`: 检查数据库连接是否为 nil

### 步骤 2: 检查数据库连接

```go
// 在 NewNumberRepo 中添加日志
func NewNumberRepo(data *Data, logger log.Logger) NumberRepo {
    logHelper := log.NewHelper(logger)
    if data == nil {
        logHelper.Error("Data is nil")
    } else if data.db == nil {
        logHelper.Error("Database connection is nil")
    } else {
        logHelper.Info("NumberRepo created successfully")
    }
    return &numberRepo{
        db:  data.db,
        log: logHelper,
    }
}
```

### 步骤 3: 检查表是否存在

```sql
-- 检查表是否存在
SHOW TABLES LIKE 'number_sequences';

-- 如果不存在，检查迁移状态
-- 运行迁移
./sre-client migrate
```

### 步骤 4: 手动测试生成器

```go
// 在测试代码中
func TestNumberGenerator(t *testing.T) {
    // 创建生成器
    data := &Data{db: testDB}
    repo := data.NewNumberRepo(data, logger)
    generator := data.NewNumberGenerator(repo, logger)
    
    // 生成编号
    no, err := generator.Generate(ctx, "ORD")
    if err != nil {
        t.Fatalf("Failed to generate: %v", err)
    }
    
    // 检查数据库
    var seq NumberSequence
    testDB.Where("prefix = ? AND date = ?", "ORD", time.Now().Format("20060102")).First(&seq)
    t.Logf("Sequence: %+v", seq)
    assert.NotEmpty(t, seq.Prefix)
}
```

## 关键检查点

1. ✅ **生成器是否为 nil**: 已添加检查
2. ✅ **数据库连接是否为 nil**: 已添加检查
3. ✅ **Repo 是否为 nil**: 已添加检查
4. ⚠️ **表是否存在**: 需要手动检查
5. ⚠️ **数据库连接是否正常**: 需要检查日志

## 日志检查清单

运行服务后，查看日志中是否有：

1. `OrderNoGenerator is nil` → 生成器未注入
2. `NumberRepo is nil` → Repo 未初始化
3. `Database connection is nil` → 数据库连接失败
4. `Generating number: prefix=ORD, date=...` → 开始生成
5. `Failed to get sequence` → 序列号获取失败
6. `Generated number: ORD...` → 生成成功

如果以上日志都没有，说明 `GenerateOrderNo()` 可能根本没有被调用。

