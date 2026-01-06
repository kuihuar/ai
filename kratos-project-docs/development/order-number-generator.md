# 唯一编号生成器实现说明

## 概述

实现了独立的、公用的、防重复的唯一编号生成器，支持不同业务使用不同前缀生成唯一编号。确保在高并发场景下生成的编号唯一。

## 特性

- ✅ **支持多业务前缀**: 不同业务可以使用不同的前缀（如：ORD-订单，PAY-支付单，REF-退款单等）
- ✅ **防重复**: 使用数据库事务 + 行锁保证唯一性
- ✅ **高并发安全**: 支持高并发场景下的原子递增
- ✅ **通用设计**: 可复用于任何需要唯一编号的业务场景

## 架构设计

### 1. 分层架构

```
┌─────────────────────────────────────┐
│   internal/pkg/number               │
│   - Generator 接口                   │
│   - DBGenerator 实现                 │
└─────────────────────────────────────┘
              ▲
              │ 依赖
              │
┌─────────────────────────────────────┐
│   internal/data                      │
│   - NumberRepo 接口                   │
│   - numberRepo 实现                   │
│   - NewNumberGenerator               │
└─────────────────────────────────────┘
              ▲
              │ 使用
              │
┌─────────────────────────────────────┐
│   internal/biz                       │
│   - OrderUsecase (使用 "ORD" 前缀)    │
│   - OrderCreateSaga                  │
│   - 其他业务 (可使用不同前缀)          │
└─────────────────────────────────────┘
```

### 2. 核心组件

#### 2.1 唯一编号生成器接口 (`internal/pkg/number/generator.go`)

```go
type Generator interface {
    // Generate 生成唯一的编号
    // prefix: 业务前缀（如 "ORD" 表示订单，"PAY" 表示支付单等）
    // 格式: {prefix}{YYYYMMDD}{6位序列号}
    // 例如: Generate(ctx, "ORD") -> ORD20241216000001
    //      Generate(ctx, "PAY") -> PAY20241216000001
    Generate(ctx context.Context, prefix string) (string, error)
}
```

#### 2.2 数据库序列号表 (`internal/data/ent/schema/number_sequence.go`)

- **表名**: `number_sequences`
- **字段**:
  - `prefix` (主键): 业务前缀（如：ORD、PAY、REF等）
  - `date` (主键): 日期（YYYYMMDD格式）
  - `sequence`: 当前序列号
  - `created_at`: 创建时间
  - `updated_at`: 更新时间

- **唯一性保证**: 通过 `(prefix, date)` 联合唯一索引保证每个业务前缀每天只有一条记录

#### 2.3 序列号仓储 (`internal/data/number_repo.go`)

```go
type NumberRepo interface {
    // GetAndIncrement 获取并递增序列号（原子操作）
    // prefix: 业务前缀
    // date: 日期（YYYYMMDD格式）
    // 返回新的序列号
    GetAndIncrement(ctx context.Context, prefix, date string) (int64, error)
}
```

**实现要点**:
- 使用数据库事务 + `SELECT FOR UPDATE` 锁定记录
- 支持并发场景下的原子递增
- 自动创建当天的序列号记录
- 支持不同前缀的独立序列号

## 编号格式

**格式**: `{prefix}{YYYYMMDD}{6位序列号}`

**示例**:
- `ORD20241216000001` - 2024年12月16日的第1个订单（前缀：ORD）
- `ORD20241216000002` - 2024年12月16日的第2个订单
- `PAY20241216000001` - 2024年12月16日的第1个支付单（前缀：PAY）
- `REF20241216000001` - 2024年12月16日的第1个退款单（前缀：REF）
- `ORD20241217000001` - 2024年12月17日的第1个订单

**特点**:
- 支持不同业务使用不同前缀
- 每个前缀每天从1开始独立计数
- 6位序列号，支持每个前缀每天最多999,999个编号
- 包含日期信息，便于查询和归档

## 常用前缀建议

| 前缀 | 业务类型 | 说明 |
|------|---------|------|
| ORD  | 订单 | Order |
| PAY  | 支付单 | Payment |
| REF  | 退款单 | Refund |
| INV  | 发票 | Invoice |
| SHI  | 发货单 | Shipment |
| RET  | 退货单 | Return |
| COU  | 优惠券 | Coupon |
| ACT  | 活动 | Activity |

## 使用方式

### 1. 在 OrderUsecase 中使用（订单号，前缀：ORD）

```go
// 生成订单号（内部使用 "ORD" 前缀）
orderNo, err := uc.GenerateOrderNo(ctx)
if err != nil {
    return nil, err
}
```

### 2. 在其他业务中使用（自定义前缀）

```go
// 直接使用生成器，传入业务前缀
paymentNo, err := generator.Generate(ctx, "PAY")
if err != nil {
    return nil, err
}

refundNo, err := generator.Generate(ctx, "REF")
if err != nil {
    return nil, err
}
```

### 3. 在 Saga 中使用

```go
// OrderCreateSaga.Run() 中
orderNo, err := s.uc.GenerateOrderNo(ctx)  // 内部使用 "ORD" 前缀
if err != nil {
    return nil, fmt.Errorf("failed to generate order number: %w", err)
}
// 使用 orderNo 作为 Saga ID
```

### 4. 创建业务特定的生成方法

```go
// 在 PaymentUsecase 中
func (uc *PaymentUsecase) GeneratePaymentNo(ctx context.Context) (string, error) {
    return uc.numberGenerator.Generate(ctx, "PAY")
}

// 在 RefundUsecase 中
func (uc *RefundUsecase) GenerateRefundNo(ctx context.Context) (string, error) {
    return uc.numberGenerator.Generate(ctx, "REF")
}
```

## 依赖注入

### Wire 配置

**数据层** (`internal/data/data.go`):
```go
var ProviderSet = wire.NewSet(
    // ...
	NewNumberRepo,
	NewNumberGenerator,
    // ...
)
```

**业务层** (`internal/biz/order.go`):
```go
type OrderUsecase struct {
    // ...
    orderNoGenerator orderno.Generator
    // ...
}

func NewOrderUsecase(
    repo OrderRepo,
    orderItemRepo OrderItemRepo,
    productRepo ProductRepo,
	orderNoGenerator number.Generator,  // 注入生成器
    logger log.Logger,
) *OrderUsecase {
    // ...
}
```

## 并发安全性

### 1. 数据库事务 + 行锁

使用 `SELECT FOR UPDATE` 在事务中锁定记录，确保并发安全：

```sql
SELECT prefix, date, sequence 
FROM number_sequences 
WHERE prefix = ? AND date = ? 
FOR UPDATE
```

### 2. 原子递增

在事务中完成：
1. 锁定记录（按 prefix + date）
2. 读取当前序列号
3. 递增序列号
4. 更新数据库
5. 提交事务

### 3. 自动创建记录

如果当天的记录不存在，自动创建：

```go
if result.Error == gorm.ErrRecordNotFound {
    // 创建新记录
		seq = NumberSequence{
        Prefix:    prefix,
        Date:      datePrefix,
        Sequence:  0,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    tx.Create(&seq)
}
```

### 4. 前缀隔离

不同前缀的序列号完全独立，互不影响：
- `ORD20241216000001` 和 `PAY20241216000001` 可以同时存在
- 每个前缀每天从1开始独立计数

## 性能考虑

### 1. 数据库索引

- `(prefix, date)` 联合唯一索引，查询速度快
- 每个前缀每天只有一条记录，数据量小
- `prefix` 字段有单独索引，便于查询某个业务的所有记录

### 2. 事务开销

- 每次生成编号需要一次数据库事务
- 对于高并发场景，可以考虑：
  - 使用 Redis 原子操作（后续实现）
  - 批量预生成编号（后续优化）

### 3. 备选实现

`GetAndIncrementRawSQL` 方法提供了使用 `INSERT ... ON DUPLICATE KEY UPDATE` 的实现，性能更好，但需要数据库支持。

## 测试

### 1. 单元测试

```go
// 测试订单号生成
func TestDBGenerator_Generate(t *testing.T) {
    // ...
}
```

### 2. 并发测试

```go
// 测试并发场景下的唯一性
func TestDBGenerator_Concurrent(t *testing.T) {
    // 启动多个 goroutine 并发生成订单号
    // 验证所有订单号都是唯一的
}
```

## 迁移数据库

运行 Ent 迁移以创建 `number_sequences` 表：

```bash
# 生成迁移文件
go run entgo.io/ent/cmd/ent generate ./internal/data/ent/schema

# 应用迁移
./sre-client migrate
```

## 注意事项

1. **前缀规范**:
   - 前缀不能为空
   - 前缀长度不能超过10个字符
   - 建议使用大写字母，如：ORD、PAY、REF等
   - 前缀应该具有业务语义，便于识别

2. **日期格式**: 使用 `YYYYMMDD` 格式（8位数字），例如 `20241216`

3. **序列号范围**: 6位数字，支持每个前缀每天最多 999,999 个编号

4. **时区**: 使用服务器本地时区，确保日期计算正确

5. **错误处理**: 生成失败时返回错误，调用方需要处理

6. **向后兼容**: 保留了 `generateOrderNo()` 方法，但新代码应使用 `GenerateOrderNo(ctx)`

7. **前缀隔离**: 不同前缀的序列号完全独立，互不影响

## 后续优化

1. **Redis 实现**: 使用 Redis 的原子操作提高性能
2. **批量生成**: 预生成一批编号，减少数据库访问
3. **分布式ID**: 考虑使用雪花算法等分布式ID生成方案
4. **监控**: 添加编号生成性能监控和告警
5. **前缀管理**: 可以考虑添加前缀配置表，统一管理所有业务前缀

