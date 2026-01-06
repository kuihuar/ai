# 唯一编号生成器改进说明

## 问题背景

之前的订单号格式 `ORD20251217000001`（前缀+日期+序列号）存在以下问题：
1. 同一天内，序列号从1开始，如果并发高，可能会有重复的风险
2. 虽然使用了数据库事务和 `FOR UPDATE`，但如果有问题，还是可能重复

旧的订单号格式 `ORD1765937919577000`（前缀+时间戳+纳秒）的优势：
1. 时间戳是秒级的，同一秒内生成多个订单的概率很低
2. 即使在同一秒内，纳秒部分也会不同
3. 不需要数据库，性能好

## 改进方案

### 新的编号格式

**格式**: `{prefix}{Unix时间戳秒}{6位序列号}`

**示例**:
- `ORD1765937919000001` - 时间戳 1765937919，序列号 1
- `ORD1765937919000002` - 时间戳 1765937919，序列号 2
- `ORD1765937920000001` - 时间戳 1765937920，序列号 1（新的一秒）

### 核心改进

1. **时间戳 + 序列号组合**
   - 使用 Unix 时间戳（秒）作为时间标识
   - 结合数据库序列号保证唯一性
   - 既保证唯一性，又避免同一天内序列号重复的风险

2. **时间区间优化**
   - 使用时间戳向下取整到10秒作为 `date` 键
   - 例如：`1765937919` -> `1765937910`（每10秒一个区间）
   - 减少序列号重置频率，同时保持足够的唯一性

3. **数据库序列号表**
   - 仍然使用 `number_sequences` 表存储序列号
   - `date` 字段存储时间戳（每10秒一个区间）
   - 通过数据库事务和 `FOR UPDATE` 保证原子性

## 实现细节

### 1. 生成器接口 (`internal/pkg/number/generator.go`)

```go
// Generate 生成唯一的编号
// 格式: {prefix}{Unix时间戳秒}{6位序列号}
// 例如: Generate(ctx, "ORD") -> ORD1765937919000001
func (g *DBGenerator) Generate(ctx context.Context, prefix string) (string, error) {
    // 获取当前 Unix 时间戳（秒）
    now := time.Now()
    unixTimestamp := now.Unix()
    
    // 使用时间戳向下取整到10秒作为 date 键
    timestampKey := (unixTimestamp / 10) * 10
    dateKey := fmt.Sprintf("%d", timestampKey)
    
    // 获取并递增序列号（原子操作）
    sequence, err := g.repo.GetAndIncrement(ctx, prefix, dateKey)
    if err != nil {
        return "", err
    }
    
    // 生成编号: {prefix}{Unix时间戳秒}{6位序列号}
    no := fmt.Sprintf("%s%d%06d", prefix, unixTimestamp, sequence)
    return no, nil
}
```

### 2. 数据库表结构

**表名**: `number_sequences`

**字段**:
- `prefix` (主键): 业务前缀（如：ORD、PAY、REF等）
- `date` (主键): 时间标识（Unix时间戳，每10秒一个区间）
- `sequence`: 当前序列号
- `created_at`: 创建时间
- `updated_at`: 更新时间

**唯一性保证**: 通过 `(prefix, date)` 联合唯一索引保证每个业务前缀每个时间区间只有一条记录

### 3. 序列号仓储 (`internal/data/number_repo.go`)

```go
// GetAndIncrement 获取并递增序列号（原子操作）
// prefix: 业务前缀
// date: 时间标识（Unix时间戳，每10秒一个区间）
// 返回新的序列号
func (r *numberRepo) GetAndIncrement(ctx context.Context, prefix, date string) (int64, error) {
    // 使用数据库事务 + SELECT FOR UPDATE 锁定记录
    // 支持并发场景下的原子递增
    // 自动创建该时间区间的序列号记录
}
```

## 优势对比

### 旧格式 `ORD20251217000001`（日期+序列号）
- ❌ 同一天内，序列号从1开始，并发高时可能重复
- ✅ 可读性好，包含日期信息
- ✅ 使用数据库序列号表保证唯一性

### 旧格式 `ORD1765937919577000`（时间戳+纳秒）
- ✅ 不需要数据库，性能好
- ✅ 同一秒内生成多个订单的概率很低
- ❌ 没有数据库保证，理论上可能重复（虽然概率极低）

### 新格式 `ORD1765937919000001`（时间戳+序列号）
- ✅ 结合时间戳和序列号，既保证唯一性，又避免重复风险
- ✅ 使用数据库序列号表保证唯一性
- ✅ 时间戳提供时间信息，序列号提供唯一性
- ✅ 每10秒一个区间，减少序列号重置频率
- ⚠️ 需要数据库，但性能影响很小（事务很快）

## 迁移说明

### 数据库迁移

1. **更新表结构**:
   ```sql
   ALTER TABLE number_sequences 
   MODIFY COLUMN date VARCHAR(20) COMMENT '时间标识(Unix时间戳，每10秒一个区间)';
   ```

2. **清理旧数据**（可选）:
   ```sql
   -- 如果需要清理旧的日期格式数据
   DELETE FROM number_sequences WHERE LENGTH(date) = 8;
   ```

### 代码变更

1. ✅ 更新 `Generator` 接口注释
2. ✅ 更新 `DBGenerator.Generate()` 实现
3. ✅ 更新 `NumberSequence` schema 注释
4. ✅ 更新 `number_repo.go` 中的结构体注释

## 测试建议

1. **并发测试**: 同时生成多个订单号，验证唯一性
2. **时间边界测试**: 测试跨10秒边界时的序列号生成
3. **数据库故障测试**: 测试数据库连接失败时的错误处理
4. **性能测试**: 测试高并发场景下的性能表现

## 总结

新的编号格式 `ORD1765937919000001` 结合了时间戳和序列号的优点：
- 时间戳提供时间信息和唯一性基础
- 序列号提供额外的唯一性保证
- 数据库事务保证原子性
- 每10秒一个区间，减少序列号重置频率

这样既保证了唯一性，又避免了同一天内序列号重复的风险。

