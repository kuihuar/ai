# Biz 层文件命名分析

## 当前文件职责分析

### 1. Usecase 文件（包含完整业务逻辑）

| 文件名 | 包含内容 | 命名合理性 | 说明 |
|--------|---------|-----------|------|
| `user.go` | User 模型 + UserRepo 接口 + UserUsecase | ✅ **合理** | 完整的业务实体文件 |
| `order.go` | Order 模型 + OrderRepo 接口 + OrderUsecase | ✅ **合理** | 完整的业务实体文件 |
| `product.go` | Product 模型 + ProductRepo 接口 + ProductUsecase | ✅ **合理** | 完整的业务实体文件 |
| `dingtalk_event.go` | DingTalkEventUsecase | ⚠️ **可优化** | 只有 Usecase，没有模型和 Repo |

### 2. 接口定义文件（只包含接口）

| 文件名 | 包含内容 | 命名合理性 | 说明 |
|--------|---------|-----------|------|
| `cache.go` | CacheRepo 接口 | ⚠️ **不够明确** | 应该明确是 Repository 接口 |
| `external_user_service.go` | ExternalUser 模型 + ExternalUserService 接口 | ⚠️ **可优化** | 可以更明确是 Service 接口 |

## 命名问题分析

### 问题 1：`cache.go` 命名不够明确

**当前**：`cache.go` - 只包含 `CacheRepo` 接口

**问题**：
- 文件名没有体现这是 Repository 接口
- 容易与 Usecase 文件混淆

**建议**：`repo_cache.go` 或 `cache_repo.go`

### 问题 2：`external_user_service.go` 命名可优化

**当前**：`external_user_service.go` - 包含模型和 Service 接口

**问题**：
- 命名较长
- 可以更明确标识是 Service 接口

**建议**：`service_external_user.go` 或保持现状

### 问题 3：`dingtalk_event.go` 命名风格不一致

**当前**：`dingtalk_event.go` - 只有 Usecase

**问题**：
- 与其他 Usecase 文件（`user.go`, `order.go`）命名风格不一致
- 如果希望统一，可以考虑使用前缀

**建议**：保持现状或统一为 `usecase_dingtalk_event.go`

## 已实施的优化方案

### ✅ 方案 1：统一接口文件命名（已完成）

**目标**：让接口定义文件更清晰，便于区分

**已执行**：
1. ✅ `cache.go` → `repo_cache.go` ⭐
2. ✅ `external_user_service.go` → `service_external_user.go` ⭐

**优点**：
- ✅ 通过前缀清晰标识文件类型
- ✅ 接口文件在文件列表中会聚集
- ✅ 便于区分 Usecase 文件和接口文件

### ✅ 方案 2：统一 Usecase 文件命名（已完成）

**目标**：统一所有 Usecase 文件的命名风格

**已执行**：
1. ✅ `dingtalk_event.go` → `usecase_dingtalk_event.go` ⭐

**优点**：
- ✅ 统一的命名风格
- ✅ 所有 Usecase 文件会聚集
- ✅ 便于识别业务逻辑文件

## 命名规范建议

### 文件类型识别

| 文件内容 | 推荐命名模式 | 示例 |
|---------|------------|------|
| Repository 接口 | `repo_*.go` | `repo_cache.go` |
| Service 接口 | `service_*.go` | `service_external_user.go` |
| Usecase（完整实体） | `*.go` | `user.go`, `order.go` |
| Usecase（仅业务逻辑） | `usecase_*.go` 或 `*.go` | `dingtalk_event.go` |

### 命名原则

1. **接口文件**：使用前缀标识类型
   - `repo_*.go` - Repository 接口
   - `service_*.go` - Service 接口

2. **Usecase 文件**：简洁命名或统一前缀
   - 实体文件：`user.go`, `order.go`（简洁）
   - 业务逻辑：`dingtalk_event.go` 或 `usecase_dingtalk_event.go`

## 优化后的目录结构（已实施）

```
internal/biz/
├── biz.go
│
├── user.go                   # Usecase（保持）✅
├── order.go                  # Usecase（保持）✅
├── product.go                # Usecase（保持）✅
├── usecase_dingtalk_event.go # Usecase（已重命名）⭐
│
├── repo_cache.go             # Repository 接口 ⭐（已重命名）
├── service_external_user.go  # Service 接口 ⭐（已重命名）
│
├── cache_example.go          # 示例文件
└── ...
```

## 已实施的优化步骤

### ✅ 步骤 1：重命名接口文件（已完成）

1. ✅ `cache.go` → `repo_cache.go`
   - 原因：只包含接口定义，命名应该明确
   - 状态：已完成

2. ✅ `external_user_service.go` → `service_external_user.go`
   - 原因：统一接口文件命名风格
   - 状态：已完成

### ✅ 步骤 2：统一 Usecase 文件命名（已完成）

3. ✅ `dingtalk_event.go` → `usecase_dingtalk_event.go`
   - 原因：统一 Usecase 文件命名风格
   - 状态：已完成

### ✅ 步骤 3：更新文档引用（已完成）

- ✅ 更新 `internal/data/cache_README.md` 中的引用
- ✅ 更新 `internal/data/tracing_example.go` 中的注释
- ✅ 更新 `internal/data/docs/dingtalk-architecture.md` 中的引用

### ✅ 步骤 4：验证编译（已完成）

- ✅ 代码编译通过
- ✅ 无 linter 错误
- ✅ 所有引用已更新

## 总结

**已解决的命名问题**：
1. ✅ `cache.go` → `repo_cache.go` - 已重命名，明确是 Repository 接口
2. ✅ `external_user_service.go` → `service_external_user.go` - 已重命名，统一接口文件命名
3. ✅ `dingtalk_event.go` → `usecase_dingtalk_event.go` - 已重命名，统一 Usecase 文件命名

**最终命名规范**：
- ✅ **接口文件**：使用前缀（`repo_*`, `service_*`）
  - `repo_cache.go` - Repository 接口
  - `service_external_user.go` - Service 接口
- ✅ **Usecase 文件**：统一使用前缀（`usecase_*.go`）或简洁命名（`*.go`）
  - `user.go`, `order.go`, `product.go` - 完整业务实体（保持简洁）
  - `usecase_dingtalk_event.go` - 业务逻辑 Usecase（使用前缀）

**优化收益**：
- ✅ 文件类型清晰可辨
- ✅ 接口文件和 Usecase 文件自动分组
- ✅ 符合命名规范和最佳实践

