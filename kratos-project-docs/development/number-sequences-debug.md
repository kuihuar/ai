# number_sequences 表无数据问题排查

## 当前订单生成逻辑

### 1. 生成流程

```
OrderUsecase.CreateOrder()
  └─> uc.GenerateOrderNo(ctx)
      └─> uc.orderNoGenerator.Generate(ctx, "ORD")
          └─> g.repo.GetAndIncrement(ctx, "ORD", "20241216")
              └─> 数据库事务中：
                  1. SELECT ... FROM number_sequences WHERE prefix='ORD' AND date='20241216' FOR UPDATE
                  2. 如果不存在，创建新记录 (sequence=0)
                  3. sequence = sequence + 1
                  4. UPDATE number_sequences SET sequence=?, ...
```

### 2. 依赖注入检查

从 `cmd/sre/wire_gen.go` 可以看到：

```go
numberRepo := data.NewNumberRepo(dataData, logger)
generator := data.NewNumberGenerator(numberRepo, logger)
orderUsecase := biz.NewOrderUsecase(orderRepo, orderItemRepo, productRepo, generator, logger)
```

✅ **生成器已正确注入**

### 3. 可能的问题

#### 问题 1: 数据库表未创建

**检查方法**：
```sql
SHOW TABLES LIKE 'number_sequences';
```

**解决方法**：
```bash
# 运行数据库迁移
./sre-client migrate
```

#### 问题 2: 生成器为 nil

**检查方法**：
- 如果生成器为 nil，调用 `GenerateOrderNo()` 会 panic
- 检查 Wire 依赖注入是否正确

**解决方法**：
- 确保 Wire 配置中包含 `data.NewNumberGenerator`
- 检查服务是否重启（需要重新生成 Wire 代码）
- 确保 `orderNoGenerator` 字段已正确注入

#### 问题 3: 数据库连接问题

**检查方法**：
- 查看日志中是否有数据库错误
- 检查 `number_sequences` 表是否存在

#### 问题 4: 事务回滚

**检查方法**：
- 查看日志中是否有 `Failed to get and increment sequence` 错误
- 检查数据库事务日志

## 排查步骤

### 步骤 1: 检查表是否存在

```sql
-- 检查表是否存在
SHOW TABLES LIKE 'number_sequences';

-- 如果不存在，检查表结构
DESC number_sequences;
```

### 步骤 2: 检查日志

查看应用日志，搜索：
- `OrderNoGenerator not injected` - 如果出现，说明生成器为 nil
- `Failed to get sequence` - 如果出现，说明数据库操作失败
- `Generated order number` - 如果出现，说明生成成功

### 步骤 3: 手动测试生成器

```go
// 在测试代码中
func TestNumberGenerator(t *testing.T) {
    // 创建生成器
    repo := data.NewNumberRepo(dataData, logger)
    generator := data.NewNumberGenerator(repo, logger)
    
    // 生成编号
    no, err := generator.Generate(ctx, "ORD")
    if err != nil {
        t.Fatalf("Failed to generate: %v", err)
    }
    
    // 检查数据库
    var seq NumberSequence
    db.Where("prefix = ? AND date = ?", "ORD", time.Now().Format("20060102")).First(&seq)
    t.Logf("Sequence: %+v", seq)
}
```

### 步骤 4: 检查数据库迁移

```bash
# 检查迁移状态
./sre-client migrate status

# 运行迁移
./sre-client migrate
```

## 验证方法

### 方法 1: 创建订单后检查表

```sql
-- 创建订单后，检查 number_sequences 表
SELECT * FROM number_sequences WHERE prefix = 'ORD';
```

### 方法 2: 查看订单号格式

订单号格式应该是 `ORD{日期}{序列号}`（如 `ORD20241216000001`），表示使用了数据库生成器。

如果格式不正确，说明生成器可能有问题。

## 常见问题

### Q1: 为什么表没有数据？

**A**: 可能原因：
1. 表未创建（需要运行迁移）
2. 生成器为 nil（会导致 panic）
3. 数据库连接失败
4. 事务回滚

### Q2: 如何确认使用了数据库生成器？

**A**: 
1. 检查订单号格式：应该是 `ORD20241216000001` 格式
2. 检查日志：应该有 `Generated order number: ORD...` 日志
3. 检查数据库：`number_sequences` 表应该有数据

### Q3: 如何强制使用数据库生成器？

**A**: 
1. 确保 Wire 依赖注入正确
2. 确保数据库表已创建
3. 重启服务（重新加载 Wire 生成的代码）

## 修复建议

1. **运行数据库迁移**：
   ```bash
   ./sre-client migrate
   ```

2. **检查 Wire 代码**：
   ```bash
   go generate ./cmd/sre
   ```

3. **重启服务**：
   确保使用最新的 Wire 生成的代码

4. **检查日志**：
   查看是否有错误或警告信息

