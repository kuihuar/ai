# Ent 表操作指南 - 第一步：定义 Schema

本文档是 Ent 表操作指南的第一部分，介绍如何定义 Ent Schema。

## 概述

Ent Schema 是 Ent ORM 的核心，它定义了数据库表的结构、字段、索引、约束等。在修改或新增表之前，首先需要在 `internal/data/ent/schema/` 目录下定义或修改 Schema 文件。

## 准备工作

### 1. 安装 Ent CLI

确保已安装 Ent 代码生成工具：

```bash
go install entgo.io/ent/cmd/ent@latest
```

### 2. 了解 Schema 文件位置

所有 Schema 定义文件位于：
```
internal/data/ent/schema/
```

## 新增表：定义 Schema

### 步骤 1：创建 Schema 文件

在 `internal/data/ent/schema/` 目录下创建新的 Schema 文件，例如 `product.go`：

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Comment("主键ID"),
		field.String("name").
			MaxLen(100).
			Comment("产品名称"),
		field.String("sku").
			MaxLen(64).
			Unique().
			Comment("SKU编码"),
		field.Int64("price").
			Comment("价格（分）"),
		field.Int32("stock").
			Default(0).
			Comment("库存数量"),
		field.Int64("created_at").
			Comment("创建时间（Unix时间戳）"),
		field.Int64("updated_at").
			Comment("更新时间（Unix时间戳）"),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return nil
}

// Indexes of the Product.
func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sku"),
		index.Fields("name"),
	}
}
```

### 步骤 2：字段类型说明

Ent 支持以下常用字段类型：

| Ent 类型 | Go 类型 | 说明 |
|---------|---------|------|
| `field.Int64()` | `int64` | 64位整数，常用于 ID、时间戳 |
| `field.Int32()` | `int32` | 32位整数 |
| `field.String()` | `string` | 字符串 |
| `field.Bool()` | `bool` | 布尔值 |
| `field.Float64()` | `float64` | 浮点数 |
| `field.Time()` | `time.Time` | 时间类型 |
| `field.JSON()` | `[]byte` | JSON 数据 |

### 步骤 3：常用字段选项

#### 基础选项

```go
field.String("name").
	MaxLen(100).              // 最大长度
	MinLen(1).                // 最小长度
	Default("").              // 默认值
	Optional().               // 可选字段（允许 NULL）
	Comment("产品名称").       // 字段注释
```

#### 唯一约束

```go
field.String("sku").
	Unique().                 // 唯一约束
	Comment("SKU编码")
```

#### 索引

```go
// 单字段索引
index.Fields("user_id")

// 复合索引
index.Fields("user_id", "status")

// 唯一索引
index.Fields("order_no").Unique()
```

### 步骤 4：参考示例

可以参考现有的 `Order` Schema 定义：

```12:64:internal/data/ent/schema/order.go
// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Comment("主键ID"),
		field.Int64("user_id").
			Comment("用户ID"),
		field.String("order_no").
			MaxLen(64).
			Unique().
			Comment("订单号"),
		field.Int32("status").
			Default(1).
			Comment("订单状态: 1-待支付, 2-已支付, 3-已发货, 4-已完成, 5-已取消, 6-已退款"),
		field.Int64("amount").
			Comment("订单金额（分）"),
		field.String("currency").
			MaxLen(10).
			Default("CNY").
			Comment("货币类型"),
		field.String("description").
			MaxLen(500).
			Optional().
			Comment("订单描述"),
		field.Int64("created_at").
			Comment("创建时间（Unix时间戳）"),
		field.Int64("updated_at").
			Comment("更新时间（Unix时间戳）"),
		field.Int64("paid_at").
			Default(0).
			Comment("支付时间（Unix时间戳，0表示未支付）"),
		field.Int64("cancelled_at").
			Default(0).
			Comment("取消时间（Unix时间戳，0表示未取消）"),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return nil
}

// Indexes of the Order.
func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("order_no"),
		index.Fields("status"),
	}
}
```

## 修改表：更新 Schema

### 步骤 1：修改字段

在现有的 Schema 文件中修改字段定义：

```go
// 添加新字段
field.String("new_field").
	MaxLen(100).
	Optional().
	Comment("新字段")

// 修改字段属性
field.String("existing_field").
	MaxLen(200).              // 修改最大长度
	Optional().                // 改为可选
	Comment("更新后的注释")
```

### 步骤 2：添加索引

在 `Indexes()` 方法中添加新索引：

```go
func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("order_no"),
		index.Fields("status"),
		index.Fields("created_at"),  // 新增索引
	}
}
```

### 步骤 3：删除字段

⚠️ **注意**：删除字段需要谨慎处理：

1. **不要直接删除字段**：直接删除会导致数据丢失
2. **先标记为废弃**：在字段注释中标记为废弃
3. **后续版本删除**：在确认无影响后再删除

```go
// 标记为废弃（不推荐直接删除）
field.String("old_field").
	Optional().
	Comment("已废弃，将在下个版本删除")
```

## 最佳实践

### 1. 命名规范

- Schema 名称使用**大驼峰**（PascalCase），如 `Product`、`OrderItem`
- 字段名称使用**蛇形命名**（snake_case），如 `user_id`、`created_at`
- 文件名称与 Schema 名称一致，使用小写，如 `product.go`、`order_item.go`

### 2. 时间字段

推荐使用 Unix 时间戳（`int64`）存储时间：

```go
field.Int64("created_at").
	Comment("创建时间（Unix时间戳）"),
field.Int64("updated_at").
	Comment("更新时间（Unix时间戳）"),
```

### 3. 状态字段

使用 `int32` 存储状态，并在注释中说明状态值：

```go
field.Int32("status").
	Default(1).
	Comment("订单状态: 1-待支付, 2-已支付, 3-已发货, 4-已完成, 5-已取消, 6-已退款"),
```

### 4. 金额字段

使用 `int64` 存储金额（以分为单位），避免浮点数精度问题：

```go
field.Int64("amount").
	Comment("订单金额（分）"),
```

### 5. 必填字段

- 业务必填字段：不使用 `Optional()`
- 可选字段：使用 `Optional()`
- 有默认值的字段：使用 `Default()`

### 6. 索引设计

- 为经常用于查询条件的字段添加索引
- 为外键字段添加索引
- 为唯一字段添加唯一索引
- 避免过度索引（影响写入性能）

## 常见问题

### Q1: 如何定义外键关系？

使用 `Edges()` 方法定义实体间的关系，详见 Ent 官方文档的 Edges 部分。

### Q2: 如何添加字段验证？

使用 `Validate()` 方法：

```go
field.String("email").
	Validate(func(s string) error {
		if !strings.Contains(s, "@") {
			return fmt.Errorf("invalid email")
		}
		return nil
	})
```

### Q3: 如何定义枚举类型？

使用 `Values()` 方法：

```go
field.Enum("status").
	Values("pending", "paid", "cancelled").
	Default("pending")
```

## 下一步

完成 Schema 定义后，请继续阅读：

- [第二步：生成 Ent 代码](./ent-table-operations-02-code-generation.md) - 生成 Ent 代码和迁移文件

