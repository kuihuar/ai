# Data 目录结构优化建议

## 当前目录结构

```
internal/data/
├── cache_README.md
├── cache_tracing_example.go
├── cache.go
├── data.go                    # 数据层初始化和 Data 结构体
├── dingtalk_client.go
├── docs/                       # 文档
├── ent/                        # Ent ORM（自动生成）
├── external/                   # 第三方服务客户端
│   ├── dingtalk/
│   └── wps/
├── external_user_service.go
├── kafka.go
├── locker.go
├── order_repo.go              # Repository 实现
├── product_repo.go           # Repository 实现
├── README.md
├── redis.go
├── sql/                       # SQL 脚本
├── user_repo.go              # Repository 实现
└── utils.go                   # 工具函数
```

## 优化方案（已实施）

> **注意**：由于 Go 语言要求同一个包的所有文件必须在同一个目录下，我们采用**文件命名约定**的方式来优化，而不是创建子目录。

### 1. Repository 文件命名优化

**优化前**：
- `user_repo.go`
- `order_repo.go`
- `product_repo.go`

**优化后**：
- `repo_user.go` ⭐
- `repo_order.go` ⭐
- `repo_product.go` ⭐

**优点**：
- ✅ 通过文件名前缀清晰标识 Repository 文件
- ✅ 在文件列表中 Repository 文件会聚集在一起
- ✅ 符合 Go 包结构要求（所有文件在同一目录）

### 2. 客户端文件命名优化

**优化前**：
- `redis.go`
- `kafka.go`
- `locker.go`

**优化后**：
- `client_redis.go` ⭐
- `client_kafka.go` ⭐
- `client_locker.go` ⭐

**优点**：
- ✅ 通过文件名前缀清晰标识客户端文件
- ✅ 在文件列表中客户端文件会聚集在一起
- ✅ 符合 Go 包结构要求

### 3. 缓存文件命名优化

**优化前**：
- `cache.go`
- `cache_README.md`
- `cache_tracing_example.go`

**优化后**：
- `cache_impl.go` ⭐（实现文件）
- `cache_README.md`（保持不变）
- `cache_tracing_example.go`（可选，可重命名为 `cache_tracing_example.go.bak`）

**优点**：
- ✅ 通过文件名前缀清晰标识缓存相关文件
- ✅ 在文件列表中缓存文件会聚集在一起
- ✅ 符合 Go 包结构要求

### 4. 第三方服务客户端组织

**当前状态**：
- `external/` 目录已存在，结构良好
- `dingtalk_client.go` 在根目录，应该移到 `external/dingtalk/` 或保持当前结构

**建议**：
- 保持 `external/` 目录结构
- 如果 `dingtalk_client.go` 只是创建函数，可以保留在根目录或移到 `clients/`

### 5. 工具函数组织

**当前状态**：
- `utils.go` 在根目录，只有一个函数

**建议**：
- 如果工具函数较少，可以保留在根目录
- 如果后续增加，可以考虑创建 `utils/` 目录

## 优化后的目录结构（已实施）

```
internal/data/
├── data.go                    # 数据层初始化和 Data 结构体
├── utils.go                   # 工具函数
│
├── repo_user.go              # Repository: User ⭐
├── repo_order.go             # Repository: Order ⭐
├── repo_product.go           # Repository: Product ⭐
│
├── client_redis.go           # 客户端: Redis ⭐
├── client_kafka.go           # 客户端: Kafka ⭐
├── client_locker.go          # 客户端: 分布式锁 ⭐
│
├── cache_impl.go             # 缓存实现 ⭐
├── cache_README.md           # 缓存文档
│
├── external/                  # 第三方服务客户端（保持原结构）
│   ├── dingtalk/
│   │   ├── client.go
│   │   └── types.go
│   └── wps/
│       ├── client.go
│       ├── signature.go
│       └── types.go
│
├── ent/                       # Ent ORM（自动生成，不修改）
│   ├── schema/               # Schema 定义（手动编写）
│   └── ...
│
├── docs/                      # 文档
│   ├── dingtalk-architecture.md
│   └── las-full-sync-usage.md
│
└── sql/                       # SQL 脚本
    └── ...
```

## 已实施的优化步骤

### ✅ 步骤 1：重命名 Repository 文件

```bash
mv internal/data/user_repo.go internal/data/repo_user.go
mv internal/data/order_repo.go internal/data/repo_order.go
mv internal/data/product_repo.go internal/data/repo_product.go
```

### ✅ 步骤 2：重命名客户端文件

```bash
mv internal/data/redis.go internal/data/client_redis.go
mv internal/data/kafka.go internal/data/client_kafka.go
mv internal/data/locker.go internal/data/client_locker.go
```

### ✅ 步骤 3：重命名缓存文件

```bash
mv internal/data/cache.go internal/data/cache_impl.go
```

### ✅ 步骤 4：验证编译

所有文件都在同一个包（`package data`）下，重命名文件不影响：
- ✅ 函数和类型的可见性
- ✅ 导入路径
- ✅ Wire 依赖注入
- ✅ 代码编译

## 注意事项

### 1. Go 包结构限制

**重要**：Go 语言要求同一个包的所有 `.go` 文件必须在同一个目录下。因此我们采用**文件命名约定**而不是创建子目录。

### 2. 包名保持不变

所有文件都在 `package data` 下，重命名文件不会影响包结构。

### 3. 导入路径

由于所有文件都在同一个包下，重命名文件不会影响导入路径。

### 4. 测试文件

如果有测试文件，建议也采用相同的命名约定：
- `repo_user_test.go`
- `repo_order_test.go`
- `client_redis_test.go`

### 5. Wire 依赖注入

`data.go` 中的 `ProviderSet` 不需要修改，因为函数名和包名都没有变化。

### 6. 文件排序

使用命名约定后，在 IDE 的文件列表中：
- Repository 文件会按 `repo_*` 前缀聚集
- 客户端文件会按 `client_*` 前缀聚集
- 缓存文件会按 `cache_*` 前缀聚集

## 优化收益

1. **更清晰的结构**：通过文件命名约定按功能分类
2. **更好的可维护性**：相关文件通过命名前缀聚集
3. **更好的可扩展性**：新增功能时遵循命名约定即可
4. **符合 Go 最佳实践**：遵循 Go 包结构要求，同时保持代码组织清晰
5. **IDE 友好**：文件列表中的文件会按前缀自动排序分组

## 可选优化

### 1. 将 `external_user_service.go` 移到 `external/` 目录

如果 `external_user_service.go` 是第三方服务相关，可以考虑：
- 移到 `external/user_service/` 目录
- 或保持当前结构（如果只是服务封装）

### 2. 将 `dingtalk_client.go` 移到 `clients/` 或 `external/dingtalk/`

根据其职责：
- 如果是客户端创建函数，可以移到 `clients/dingtalk.go`
- 如果是服务封装，可以移到 `external/dingtalk/client.go`

## 总结

主要优化点（已实施）：
1. ✅ **Repository 文件** → 重命名为 `repo_*.go` 前缀
2. ✅ **客户端文件** → 重命名为 `client_*.go` 前缀
3. ✅ **缓存文件** → 重命名为 `cache_*.go` 前缀
4. ✅ **第三方服务** → 保持 `external/` 结构（已良好）
5. ✅ **工具函数** → 保留 `utils.go`（函数少）

**优化方式**：采用文件命名约定而非目录结构，既符合 Go 包结构要求，又实现了代码组织的清晰化。

这样的结构更清晰、更易维护，符合 Clean Architecture 和 Kratos 最佳实践。

