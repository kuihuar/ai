# Biz 目录结构优化建议

## 当前目录结构

```
internal/biz/
├── biz.go                    # ProviderSet 和配置提供
├── cache_example.go          # 缓存使用示例
├── cache.go                  # 缓存接口定义
├── cron/                     # Cron 任务相关
│   ├── job.go
│   ├── jobs/
│   │   └── sync_user.go
│   └── manager.go
├── daemon/                   # Daemon 任务相关
│   ├── 01-daemonjob-technologies.md
│   ├── 02-comparison-with-ants.md
│   ├── 04-implementation-comparison.md
│   ├── 05-goroutine-pool-explained.md
│   ├── 06-faq-daemonjob.md
│   ├── 07-daemon-vs-cron-comparison.md
│   ├── 08-open-source-alternatives.md
│   ├── 09-engineering-optimizations.md
│   ├── 10-architecture-placement.md
│   ├── example_ants_with_lock.go
│   ├── example_ants.go
│   ├── example_pool_concept.go
│   ├── example.go
│   ├── job.go
│   ├── README.md
│   ├── table_consumer_ants_biz.go
│   ├── table_consumer_ants.go
│   └── table_consumer.go
├── dingtalk_event.go         # 钉钉事件业务逻辑
├── external_user_service.go  # 外部用户服务接口
├── order.go                  # Order 业务逻辑
├── product.go                # Product 业务逻辑
├── README.md                 # 文档
└── user.go                   # User 业务逻辑
```

## 优化建议

### 1. 业务实体文件命名优化（可选）

**当前状态**：
- `user.go` - User 业务逻辑
- `order.go` - Order 业务逻辑
- `product.go` - Product 业务逻辑

**优化方案**（可选）：
采用统一的命名约定，便于识别：
- `usecase_user.go` 或保持 `user.go`
- `usecase_order.go` 或保持 `order.go`
- `usecase_product.go` 或保持 `product.go`

**建议**：
- ⚠️ **保持现状**：如果文件不多，保持简洁的命名更好
- ✅ **如果实体增多**：可以考虑使用 `usecase_*.go` 前缀

### 2. 缓存相关文件组织

**当前问题**：
- `cache.go` - 接口定义
- `cache_example.go` - 使用示例

**优化方案**：
```
internal/biz/
├── cache.go                  # 缓存接口定义（保持不变）
└── cache_example.go          # 缓存示例（可选：重命名为 cache_example.go.bak 或移到 docs）
```

**建议**：
- ✅ `cache.go` 保持现状（接口定义）
- ⚠️ `cache_example.go` 可以：
  - 保留（如果作为参考示例）
  - 重命名为 `cache_example.go.bak`（如果不再使用）
  - 移到 `docs/development/` 目录（如果只是文档）

### 3. 外部服务接口文件命名

**当前状态**：
- `external_user_service.go` - 外部用户服务接口

**优化方案**：
- 保持现状（如果只有一个）
- 如果有多个外部服务，可以考虑：
  - `external_user_service.go` → `external_user.go`
  - 或创建 `external/` 目录（如果外部服务接口增多）

**建议**：
- ✅ **保持现状**：文件少时，保持简洁命名

### 4. Daemon 目录文档整理

**当前问题**：
- `daemon/` 目录下有很多文档文件（10 个 .md 文件）
- 文档和代码混在一起

**优化方案**：
```
internal/biz/daemon/
├── job.go                    # 核心接口和实现
├── table_consumer.go         # 核心实现
├── table_consumer_ants.go    # Ants 实现
├── table_consumer_ants_biz.go # Ants 业务逻辑
├── README.md                 # 主要文档
└── docs/                     # 文档目录（新建）
    ├── 01-daemonjob-technologies.md
    ├── 02-comparison-with-ants.md
    ├── 04-implementation-comparison.md
    ├── 05-goroutine-pool-explained.md
    ├── 06-faq-daemonjob.md
    ├── 07-daemon-vs-cron-comparison.md
    ├── 08-open-source-alternatives.md
    ├── 09-engineering-optimizations.md
    └── 10-architecture-placement.md
```

**优点**：
- ✅ 代码和文档分离
- ✅ 目录更清晰
- ✅ 便于查找代码文件

### 5. Daemon 示例文件整理

**当前问题**：
- `daemon/` 目录下有多个示例文件：
  - `example.go`
  - `example_ants.go`
  - `example_ants_with_lock.go`
  - `example_pool_concept.go`

**优化方案**（已实施）：
> **注意**：由于 Go 语言要求同一个包的所有文件必须在同一个目录下，示例文件不能移到子目录。但示例文件已经使用了 `example*.go` 命名模式，这是合理的。

**保持现状**：
```
internal/biz/daemon/
├── example.go                # 示例文件（保持现状）
├── example_ants.go
├── example_ants_with_lock.go
└── example_pool_concept.go
```

**优点**：
- ✅ 通过 `example*` 前缀清晰标识示例文件
- ✅ 符合 Go 包结构要求
- ✅ 在文件列表中示例文件会聚集在一起

### 6. 钉钉事件文件命名

**当前状态**：
- `dingtalk_event.go` - 钉钉事件业务逻辑

**优化方案**：
- 保持现状（如果只有一个钉钉相关业务）
- 如果有多个钉钉相关业务，可以考虑：
  - `usecase_dingtalk_event.go`
  - 或创建 `dingtalk/` 目录

**建议**：
- ✅ **保持现状**：文件少时，保持简洁命名

## 推荐优化方案（优先级排序）

### 高优先级：Daemon 目录整理

**理由**：文档和示例文件较多，影响代码查找

**操作**：
1. 创建 `internal/biz/daemon/docs/` 目录
2. 移动所有 `.md` 文档文件到 `docs/` 目录
3. 创建 `internal/biz/daemon/examples/` 目录
4. 移动所有 `example_*.go` 文件到 `examples/` 目录

### 中优先级：缓存示例文件处理

**理由**：示例文件可能不再需要

**操作**：
- 如果不再使用，重命名为 `cache_example.go.bak`
- 或移到 `docs/development/` 目录

### 低优先级：业务实体文件命名

**理由**：当前文件不多，命名简洁清晰

**操作**：
- 保持现状
- 如果后续实体增多（>10 个），再考虑统一命名约定

## 优化后的目录结构（已实施）

```
internal/biz/
├── biz.go                    # ProviderSet
├── cache.go                  # 缓存接口定义
├── cache_example.go          # 缓存示例（保持现状）
│
├── user.go                   # User 业务逻辑（保持现状）
├── order.go                  # Order 业务逻辑（保持现状）
├── product.go                # Product 业务逻辑（保持现状）
│
├── dingtalk_event.go        # 钉钉事件业务逻辑
├── external_user_service.go # 外部用户服务接口
│
├── cron/                     # Cron 任务（保持现状）
│   ├── job.go
│   ├── jobs/
│   └── manager.go
│
├── daemon/                   # Daemon 任务（已优化）⭐
│   ├── job.go                # 核心接口
│   ├── table_consumer.go     # 核心实现
│   ├── table_consumer_ants.go
│   ├── table_consumer_ants_biz.go
│   ├── example.go            # 示例文件（保持现状）
│   ├── example_ants.go
│   ├── example_ants_with_lock.go
│   ├── example_pool_concept.go
│   ├── README.md             # 主要文档
│   └── docs/                 # 文档目录 ⭐（已创建）
│       ├── 01-daemonjob-technologies.md
│       ├── 02-comparison-with-ants.md
│       ├── 04-implementation-comparison.md
│       ├── 05-goroutine-pool-explained.md
│       ├── 06-faq-daemonjob.md
│       ├── 07-daemon-vs-cron-comparison.md
│       ├── 08-open-source-alternatives.md
│       ├── 09-engineering-optimizations.md
│       └── 10-architecture-placement.md
│
└── README.md                 # 文档
```

## 已实施的优化步骤

### ✅ 步骤 1：整理 Daemon 目录文档（已完成）

```bash
mkdir -p internal/biz/daemon/docs
mv internal/biz/daemon/01-*.md internal/biz/daemon/docs/
mv internal/biz/daemon/02-*.md internal/biz/daemon/docs/
# ... 其他文档文件
```

**结果**：
- ✅ 10 个文档文件已移动到 `daemon/docs/` 目录
- ✅ `README.md` 保留在 `daemon/` 根目录

### ⚠️ 步骤 2：Daemon 示例文件（保持现状）

**原因**：由于 Go 语言要求同一个包的所有文件必须在同一个目录下，示例文件不能移到子目录。

**现状**：
- ✅ 示例文件已使用 `example*.go` 命名模式
- ✅ 通过文件名前缀清晰标识
- ✅ 在文件列表中会自动聚集

### ⚠️ 步骤 3：缓存示例文件（保持现状）

**原因**：
- `cache_example.go` 在文档中被引用（`internal/data/cache_README.md`）
- 作为参考示例保留

**现状**：
- ✅ 保持 `cache_example.go` 在根目录
- ✅ 作为缓存使用示例参考

### ✅ 步骤 4：验证编译（已完成）

```bash
go build ./internal/biz/...
# 编译通过 ✅
```

## 注意事项

### 1. 文档链接更新

如果文档之间有相互引用，需要更新路径：
- `daemon/01-*.md` → `daemon/docs/01-*.md`

### 2. 导入路径

所有文件都在 `package biz` 下，移动文件不会影响导入路径。

### 3. 示例代码

示例文件中的导入路径可能需要更新（如果移动到不同目录）。

## 优化收益

1. **更清晰的结构**：文档和代码分离
2. **更好的可维护性**：相关文件集中管理
3. **更好的可查找性**：核心代码文件更容易找到
4. **符合最佳实践**：遵循常见的项目组织方式

## 总结

主要优化点（已实施）：
1. ✅ **Daemon 文档整理** → `daemon/docs/` 目录（已完成）
2. ✅ **Daemon 示例文件** → 保持现状，使用 `example*.go` 命名（符合 Go 规范）
3. ✅ **缓存示例文件** → 保持现状（作为参考示例）
4. ✅ **业务实体文件** → 保持现状（命名简洁清晰）

**优化方式**：
- 采用目录结构优化（文档整理到子目录）
- 采用文件命名约定（示例文件使用 `example*` 前缀）
- 符合 Go 包结构要求（同一包的文件在同一目录）

**优化收益**：
- ✅ 文档和代码分离，目录更清晰
- ✅ 核心代码文件更容易找到
- ✅ 符合 Go 语言规范和最佳实践

这样的结构更清晰、更易维护，符合 Clean Architecture 和 Kratos 最佳实践。

