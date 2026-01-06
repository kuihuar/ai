# Biz 层文件命名规范

## 当前文件分析

### Usecase 文件（包含业务逻辑）

| 文件名 | 内容 | 命名合理性 | 建议 |
|--------|------|-----------|------|
| `user.go` | User 模型 + UserRepo 接口 + UserUsecase | ✅ 合理 | 保持现状 |
| `order.go` | Order 模型 + OrderRepo 接口 + OrderUsecase | ✅ 合理 | 保持现状 |
| `product.go` | Product 模型 + ProductRepo 接口 + ProductUsecase | ✅ 合理 | 保持现状 |
| `dingtalk_event.go` | DingTalkEventUsecase | ⚠️ 可优化 | 考虑统一命名 |

### 接口定义文件（只包含接口）

| 文件名 | 内容 | 命名合理性 | 建议 |
|--------|------|-----------|------|
| `cache.go` | CacheRepo 接口 | ⚠️ 不够明确 | 重命名为 `repo_cache.go` |
| `external_user_service.go` | ExternalUserService 接口 | ⚠️ 不够明确 | 重命名为 `service_external_user.go` |

## 命名规范建议

### 原则

1. **Usecase 文件**：使用简洁的实体名称
   - ✅ `user.go` - User 业务逻辑
   - ✅ `order.go` - Order 业务逻辑
   - ✅ `product.go` - Product 业务逻辑

2. **接口定义文件**：使用前缀标识类型
   - ✅ `repo_cache.go` - 缓存 Repository 接口
   - ✅ `service_external_user.go` - 外部用户服务接口

3. **Usecase 文件（可选）**：如果实体增多，可以使用前缀
   - ⚠️ `usecase_user.go` - 如果实体 > 10 个时考虑
   - ⚠️ `usecase_dingtalk_event.go` - 统一命名风格

## 推荐优化方案

### 方案 1：统一接口文件命名（推荐）

**目标**：让接口定义文件更清晰

**操作**：
1. `cache.go` → `repo_cache.go`（Repository 接口）
2. `external_user_service.go` → `service_external_user.go`（Service 接口）

**优点**：
- ✅ 通过前缀清晰标识文件类型
- ✅ 接口文件在文件列表中会聚集
- ✅ 便于区分 Usecase 文件和接口文件

### 方案 2：统一 Usecase 文件命名（可选）

**目标**：统一所有 Usecase 文件的命名风格

**操作**：
1. `user.go` → `usecase_user.go`
2. `order.go` → `usecase_order.go`
3. `product.go` → `usecase_product.go`
4. `dingtalk_event.go` → `usecase_dingtalk_event.go`

**优点**：
- ✅ 统一的命名风格
- ✅ 所有 Usecase 文件会聚集

**缺点**：
- ❌ 文件名变长
- ❌ 如果实体不多，可能过度设计

**建议**：
- ⚠️ 如果实体数量 < 10 个，保持简洁命名
- ✅ 如果实体数量 > 10 个，考虑统一前缀

## 最终推荐

### 高优先级：接口文件重命名

```
internal/biz/
├── cache.go                  → repo_cache.go ⭐
├── external_user_service.go  → service_external_user.go ⭐
```

### 低优先级：Usecase 文件命名

**当前状态**：文件不多，命名简洁清晰

**建议**：
- ✅ 保持现状：`user.go`, `order.go`, `product.go`
- ⚠️ `dingtalk_event.go` 可以考虑重命名为 `usecase_dingtalk_event.go`（如果希望统一）

## 文件类型识别

### 通过文件名前缀识别

| 前缀 | 类型 | 示例 |
|------|------|------|
| `repo_*.go` | Repository 接口 | `repo_cache.go` |
| `service_*.go` | Service 接口 | `service_external_user.go` |
| `usecase_*.go` | Usecase 实现 | `usecase_user.go`（可选） |
| `*.go` | Usecase 实现 | `user.go`（当前方式） |

## 实施建议

### 立即执行（高优先级）

1. ✅ `cache.go` → `repo_cache.go`
2. ✅ `external_user_service.go` → `service_external_user.go`

### 后续考虑（低优先级）

1. ⚠️ 如果实体增多，考虑统一 Usecase 文件命名
2. ⚠️ `dingtalk_event.go` → `usecase_dingtalk_event.go`（可选）

## 命名对比

### 优化前

```
internal/biz/
├── user.go                   # Usecase（清晰）
├── order.go                  # Usecase（清晰）
├── product.go                # Usecase（清晰）
├── cache.go                  # 接口（不够明确）
├── dingtalk_event.go         # Usecase（清晰）
└── external_user_service.go  # 接口（不够明确）
```

### 优化后（推荐）

```
internal/biz/
├── user.go                   # Usecase（保持）
├── order.go                  # Usecase（保持）
├── product.go                # Usecase（保持）
├── repo_cache.go             # Repository 接口 ⭐
├── dingtalk_event.go         # Usecase（保持或重命名）
└── service_external_user.go  # Service 接口 ⭐
```

## 总结

**当前问题**：
- `cache.go` 和 `external_user_service.go` 只包含接口定义，命名不够明确

**推荐优化**：
1. ✅ 重命名接口文件，使用前缀标识类型
2. ⚠️ Usecase 文件保持现状（文件不多时）

**命名原则**：
- 接口文件：使用前缀（`repo_*`, `service_*`）
- Usecase 文件：简洁命名（`user.go`）或统一前缀（`usecase_*.go`）

