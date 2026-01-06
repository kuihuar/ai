# 产品 API 列表

本文档列出了所有产品相关的 API 接口。

## API 概览

| 序号 | API 名称 | gRPC 方法 | HTTP 方法 | HTTP 路径 | 功能描述 |
|------|---------|-----------|-----------|-----------|----------|
| 1 | 创建产品 | `CreateProduct` | POST | `/api/v1/products` | 创建新产品 |
| 2 | 获取产品 | `GetProduct` | GET | `/api/v1/products/{id}` | 根据ID获取产品信息 |
| 3 | 更新产品 | `UpdateProduct` | PUT | `/api/v1/products/{id}` | 更新产品信息 |
| 4 | 删除产品 | `DeleteProduct` | DELETE | `/api/v1/products/{id}` | 删除产品 |
| 5 | 列出产品 | `ListProducts` | GET | `/api/v1/products` | 分页列出产品，支持搜索和价格筛选 |
| 6 | 根据SKU获取产品 | `GetProductBySku` | GET | `/api/v1/products/sku/{sku}` | 根据SKU编码获取产品信息 |

## 详细说明

### 1. 创建产品

**接口信息**
- **gRPC**: `product.v1.Product.CreateProduct`
- **HTTP**: `POST /api/v1/products`

**请求参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 产品名称（最大长度100） |
| sku | string | 是 | SKU编码（最大长度64，唯一） |
| price | int64 | 是 | 价格（分，必须 >= 0） |
| stock | int32 | 否 | 库存数量（默认0，必须 >= 0） |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| product | ProductInfo | 产品信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| PRODUCT_ALREADY_EXISTS | SKU 已存在 |
| PRODUCT_INVALID_NAME | 产品名称无效（为空或超过100字符） |
| PRODUCT_INVALID_SKU | SKU 无效（为空或超过64字符） |
| PRODUCT_INVALID_PRICE | 价格无效（小于0） |

---

### 2. 获取产品

**接口信息**
- **gRPC**: `product.v1.Product.GetProduct`
- **HTTP**: `GET /api/v1/products/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 产品ID |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| product | ProductInfo | 产品信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| PRODUCT_NOT_FOUND | 产品不存在 |

---

### 3. 更新产品

**接口信息**
- **gRPC**: `product.v1.Product.UpdateProduct`
- **HTTP**: `PUT /api/v1/products/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 产品ID |
| name | string | 否 | 请求体 | 产品名称（可选，不为空时更新） |
| price | int64 | 否 | 请求体 | 价格（分，可选，> 0 时更新） |
| stock | int32 | 否 | 请求体 | 库存数量（可选，>= 0 时更新） |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| product | ProductInfo | 产品信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| PRODUCT_NOT_FOUND | 产品不存在 |
| PRODUCT_INVALID_NAME | 产品名称无效 |
| PRODUCT_INVALID_PRICE | 价格无效 |

---

### 4. 删除产品

**接口信息**
- **gRPC**: `product.v1.Product.DeleteProduct`
- **HTTP**: `DELETE /api/v1/products/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 产品ID |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | bool | 是否成功 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| PRODUCT_NOT_FOUND | 产品不存在 |

---

### 5. 列出产品

**接口信息**
- **gRPC**: `product.v1.Product.ListProducts`
- **HTTP**: `GET /api/v1/products`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 默认值 | 说明 |
|--------|------|------|------|--------|------|
| page | int32 | 否 | 查询参数 | 1 | 页码（从1开始） |
| page_size | int32 | 否 | 查询参数 | 10 | 每页数量（最大100） |
| keyword | string | 否 | 查询参数 | - | 搜索关键词（可选，用于搜索产品名称或SKU） |
| min_price | int64 | 否 | 查询参数 | - | 最低价格（可选，> 0 时生效） |
| max_price | int64 | 否 | 查询参数 | - | 最高价格（可选，> 0 时生效） |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| products | ProductInfo[] | 产品列表 |
| total | int32 | 总数量 |
| page | int32 | 当前页码 |
| page_size | int32 | 每页数量 |

**示例**

```bash
# 基本查询
GET /api/v1/products?page=1&page_size=10

# 关键词搜索
GET /api/v1/products?keyword=iPhone

# 价格范围筛选
GET /api/v1/products?min_price=10000&max_price=100000

# 组合查询
GET /api/v1/products?page=1&page_size=20&keyword=iPhone&min_price=50000&max_price=100000
```

---

### 6. 根据SKU获取产品

**接口信息**
- **gRPC**: `product.v1.Product.GetProductBySku`
- **HTTP**: `GET /api/v1/products/sku/{sku}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| sku | string | 是 | 路径参数 | SKU编码 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| product | ProductInfo | 产品信息对象 |

**错误码**

| 错误码 | 说明 |
|--------|------|
| PRODUCT_NOT_FOUND | 产品不存在 |

---

## 数据模型

### ProductInfo

产品信息对象，用于返回产品的基本信息。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 产品ID |
| name | string | 产品名称 |
| sku | string | SKU编码（唯一） |
| price | int64 | 价格（分） |
| stock | int32 | 库存数量 |
| created_at | int64 | 创建时间（Unix时间戳） |
| updated_at | int64 | 更新时间（Unix时间戳） |

---

## 错误码

所有 API 的错误码定义在 `api/product/v1/error_reason.proto` 中。

常见错误码：

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| PRODUCT_NOT_FOUND | 404 | 产品不存在 |
| PRODUCT_ALREADY_EXISTS | 409 | 产品已存在（SKU重复） |
| PRODUCT_INVALID_NAME | 400 | 无效的产品名称 |
| PRODUCT_INVALID_SKU | 400 | 无效的SKU编码 |
| PRODUCT_INVALID_PRICE | 400 | 无效的产品价格 |
| PRODUCT_SAVE_FAILED | 500 | 保存产品失败 |
| PRODUCT_UPDATE_FAILED | 500 | 更新产品失败 |
| PRODUCT_QUERY_FAILED | 500 | 查询产品失败 |
| PRODUCT_DELETE_FAILED | 500 | 删除产品失败 |

---

## 使用示例

### HTTP 请求示例

#### 创建产品
```bash
curl -X POST http://localhost:8000/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "sku": "IPHONE15-PRO-001",
    "price": 799900,
    "stock": 50
  }'
```

**响应示例：**
```json
{
  "product": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "sku": "IPHONE15-PRO-001",
    "price": 799900,
    "stock": 50,
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```

#### 获取产品
```bash
curl http://localhost:8000/api/v1/products/1
```

#### 更新产品
```bash
curl -X PUT http://localhost:8000/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro Max",
    "price": 899900,
    "stock": 30
  }'
```

#### 删除产品
```bash
curl -X DELETE http://localhost:8000/api/v1/products/1
```

#### 列出产品
```bash
# 基本列表
curl "http://localhost:8000/api/v1/products?page=1&page_size=10"

# 关键词搜索
curl "http://localhost:8000/api/v1/products?keyword=iPhone"

# 价格范围筛选
curl "http://localhost:8000/api/v1/products?min_price=50000&max_price=100000"

# 组合查询
curl "http://localhost:8000/api/v1/products?page=1&page_size=20&keyword=iPhone&min_price=50000&max_price=100000"
```

#### 根据SKU获取产品
```bash
curl http://localhost:8000/api/v1/products/sku/IPHONE15-PRO-001
```

### gRPC 请求示例

```go
import (
    "context"
    "sre/api/product/v1"
    "google.golang.org/grpc"
)

conn, _ := grpc.Dial("localhost:8989", grpc.WithInsecure())
client := v1.NewProductClient(conn)

// 创建产品
product, err := client.CreateProduct(context.Background(), &v1.CreateProductRequest{
    Name:  "iPhone 15 Pro",
    Sku:   "IPHONE15-PRO-001",
    Price: 799900,
    Stock: 50,
})

// 获取产品
product, err := client.GetProduct(context.Background(), &v1.GetProductRequest{
    Id: 1,
})

// 列出产品
products, err := client.ListProducts(context.Background(), &v1.ListProductsRequest{
    Page:     1,
    PageSize: 10,
    Keyword:  "iPhone",
    MinPrice: 50000,
    MaxPrice: 100000,
})
```

---

## 业务规则

### 产品创建规则

1. **SKU 唯一性**：每个产品的 SKU 必须唯一，创建时会检查 SKU 是否已存在
2. **价格验证**：价格必须 >= 0（单位为分）
3. **库存验证**：库存数量必须 >= 0
4. **名称验证**：产品名称不能为空，最大长度 100 字符
5. **SKU 验证**：SKU 不能为空，最大长度 64 字符

### 产品更新规则

1. **部分更新**：只更新提供的字段，未提供的字段保持不变
2. **价格验证**：如果提供价格，必须 > 0
3. **库存验证**：如果提供库存，必须 >= 0
4. **名称验证**：如果提供名称，不能为空且最大长度 100 字符

### 产品查询规则

1. **分页**：默认每页 10 条，最大 100 条
2. **关键词搜索**：支持搜索产品名称和 SKU
3. **价格筛选**：支持按价格范围筛选（min_price <= price <= max_price）
4. **排序**：按 ID 降序排列（最新创建的在前）

---

## 认证和授权

### 当前配置

- **认证要求**：`/api/v1/products` 及其子路径**不需要认证**（公开接口）
- **限流策略**：公开接口限流（100 请求/分钟）

### 修改认证配置

如需为产品接口添加认证，请修改 `internal/server/http.go` 中的 `isAuthenticatedRoute` 函数和 `createAuthMiddlewareWithSkip` 函数。

---

## 相关文档

- [API 设计规范](../../code-standards/api-design.md)
- [用户 API 列表](./user-api-list.md)
- [订单 API 列表](./order-api-list.md)
- [产品业务逻辑](../../internal/biz/product.go)
