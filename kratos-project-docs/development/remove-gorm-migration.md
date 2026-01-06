# 移除 GORM 依赖迁移方案

## 概述

当前项目同时使用了 GORM 和 Ent，但 GORM 实际上只作为"连接管理器"使用，所有业务逻辑都通过 Ent 实现。本文档分析是否可以移除 GORM 依赖，并提供迁移方案。

## 当前 GORM 使用情况分析

### 1. 使用位置统计

| 位置 | 用途 | 是否可替换 |
|------|------|-----------|
| `internal/data/data.go` | 创建数据库连接 | ✅ 可替换为 `database/sql` |
| `internal/data/data.go` | 从 GORM 提取 `*sql.DB` 给 Ent | ✅ 可直接创建 `*sql.DB` |
| `internal/data/health.go` | 健康检查（获取 `*sql.DB`） | ✅ 可直接使用 `*sql.DB` |
| `internal/biz/daemon/table_consumer.go` | `BuildDefaultFetcher` 使用 `*gorm.DB` | ✅ 可改为 `*sql.DB` |
| `internal/biz/daemon/table_consumer_ants.go` | `TableConsumerDaemonAnts` 使用 `*gorm.DB` | ✅ 可改为 `*sql.DB` |
| `internal/logger/gorm.go` | GORM logger 适配器 | ⚠️ 需要移除或保留（如果其他地方需要） |

### 2. 关键发现

**✅ 所有 Repository 层都只使用 Ent**
- `repo_user.go` - 只使用 `ent.Client`
- `repo_order.go` - 只使用 `ent.Client`
- `repo_product.go` - 只使用 `ent.Client`
- `repo_order_item.go` - 只使用 `ent.Client`

**✅ GORM 只用于连接管理**
- `NewDB` 创建 GORM 连接
- `NewEntClient` 从 GORM 提取 `*sql.DB` 给 Ent
- Daemon 中通过 `db.DB()` 获取 `*sql.DB` 执行原始 SQL

**✅ 没有使用 GORM 的 ORM 功能**
- 没有使用 GORM 的查询构建器
- 没有使用 GORM 的关联查询
- 没有使用 GORM 的迁移功能

## 迁移方案

### 方案 1：完全移除 GORM（推荐）

**优点**：
- ✅ 减少依赖，简化项目
- ✅ 减少二进制体积
- ✅ 统一使用 Ent 作为 ORM
- ✅ 连接管理更直接

**缺点**：
- ⚠️ 需要修改多个文件
- ⚠️ 需要移除 GORM logger 适配器（如果不需要）

### 方案 2：保留 GORM 作为连接管理器（不推荐）

**优点**：
- ✅ 改动最小
- ✅ 保留 GORM logger 功能

**缺点**：
- ❌ 增加不必要的依赖
- ❌ 增加二进制体积
- ❌ 维护成本高

## 详细迁移步骤

### 步骤 1：修改 `internal/data/data.go`

#### 1.1 移除 GORM 导入，添加 `database/sql`

```go
package data

import (
	"database/sql"
	"os"
	"sre/internal/conf"
	"sre/internal/data/ent"
	loggerpkg "sre/internal/logger"
	"sre/internal/pkg/cipherutil"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)
```

#### 1.2 修改 `Data` 结构体

```go
// Data .
type Data struct {
	db     *sql.DB  // 改为 *sql.DB
	ent    *ent.Client
	redis  *redis.Client
	locker DistributedLocker
}

// DB 返回数据库连接（用于健康检查）
func (d *Data) DB() *sql.DB {
	return d.db
}
```

#### 1.3 修改 `NewDB` 函数

```go
// NewDB creates a new database connection.
// Returns nil if database is not configured, disabled, or connection fails (non-blocking).
func NewDB(c *conf.Data, logger log.Logger) (*sql.DB, error) {
	logHelper := log.NewHelper(logger)

	if c.Database == nil {
		logHelper.Warn("database config is nil, database features will be disabled")
		return nil, nil
	}

	// 检查是否启用数据库
	if !c.Database.Enable {
		logHelper.Info("database is disabled, database features will be disabled")
		return nil, nil
	}

	// 根据驱动类型创建数据库连接
	var db *sql.DB
	var err error

	switch c.Database.Driver {
	case "mysql":
		var dsn string
		encryptedDsn := os.Getenv("ECIS_ECISACCOUNTSYNC_DB")
		logHelper.Infof("encryptedDsn: %s", encryptedDsn)
		if encryptedDsn != "" {
			logHelper.Infof("encryptedDsn is not empty, decrypting database dsn")
			// 从配置获取解密密钥
			decryptKey := c.Database.DecryptKey

			if decryptKey == "" {
				logHelper.Warn("decrypt_key is not configured, cannot decrypt database dsn, database features will be disabled")
				return nil, nil
			}
			decryptedDsn, err := cipherutil.DecryptByAes(encryptedDsn, decryptKey)
			logHelper.Infof("decryptedDsn: %s", decryptedDsn)
			if err != nil {
				logHelper.Warnf("failed to decrypt database dsn: %v, database features will be disabled", err)
				return nil, nil
			}
			dsn = decryptedDsn
		} else {
			logHelper.Debugf("encryptedDsn is empty, using database dsn from config")
			dsn = c.Database.Source
		}
		logHelper.Infof("dsn: %s", dsn)

		// 使用 database/sql 创建连接
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logHelper.Warnf("failed to open database: %v, database features will be disabled", err)
			return nil, nil
		}

		// 配置连接池
		db.SetMaxOpenConns(25)        // 最大打开连接数
		db.SetMaxIdleConns(10)        // 最大空闲连接数
		db.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生存时间
		db.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大生存时间

		// 测试连接
		if err := db.Ping(); err != nil {
			logHelper.Warnf("failed to ping database: %v, database features will be disabled", err)
			return nil, nil
		}

	default:
		logHelper.Warnf("unsupported database driver: %s, database features will be disabled", c.Database.Driver)
		return nil, nil
	}

	logHelper.Info("database connection established")
	return db, nil
}
```

#### 1.4 修改 `NewEntClient` 函数

```go
// NewEntClient creates a new ent client from *sql.DB.
func NewEntClient(db *sql.DB, logger log.Logger) (*ent.Client, error) {
	if db == nil {
		return nil, nil
	}

	// 从 *sql.DB 创建 ent driver
	drv := entsql.OpenDB("mysql", db)
	
	// 创建 ent client
	client := ent.NewClient(ent.Driver(drv))
	return client, nil
}
```

#### 1.5 修改 `NewData` 函数

```go
// NewData .
func NewData(db *sql.DB, entClient *ent.Client, redisClient *redis.Client, locker DistributedLocker, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logHelper := log.NewHelper(logger)
		logHelper.Info("closing the data resources")
		if entClient != nil {
			if err := entClient.Close(); err != nil {
				logHelper.Errorf("failed to close ent client: %v", err)
			} else {
				logHelper.Info("ent client closed")
			}
		}
		if db != nil {
			if err := db.Close(); err != nil {
				logHelper.Errorf("failed to close database: %v", err)
			} else {
				logHelper.Info("database closed")
			}
		}
		if redisClient != nil {
			if err := redisClient.Close(); err != nil {
				logHelper.Errorf("failed to close Redis client: %v", err)
			} else {
				logHelper.Info("Redis client closed")
			}
		}
	}
	return &Data{
		db:     db,
		ent:    entClient,
		redis:  redisClient,
		locker: locker,
	}, cleanup, nil
}
```

### 步骤 2：修改 `internal/data/health.go`

```go
// HealthCheck 执行健康检查
func (d *Data) HealthCheck(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Healthy:   true,
		Message:   "healthy",
		Details:   make(map[string]string),
		Timestamp: time.Now(),
	}

	// 检查数据库连接
	if d.db != nil {
		// 设置超时上下文
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		if err := d.db.PingContext(pingCtx); err != nil {
			status.Healthy = false
			status.Message = "database ping failed"
			status.Details["database"] = "error: " + err.Error()
			return status
		}

		// 获取连接池状态
		stats := d.db.Stats()
		status.Details["database"] = "ok"
		// 连接池统计信息（可选，用于调试）
		status.Details["database_open_conns"] = fmt.Sprintf("%d", stats.OpenConnections)
		status.Details["database_idle_conns"] = fmt.Sprintf("%d", stats.Idle)
	} else {
		status.Details["database"] = "not configured"
	}

	// ... Redis 检查保持不变 ...
}

// ReadinessCheck 执行就绪检查
func (d *Data) ReadinessCheck(ctx context.Context, logger log.Logger) HealthStatus {
	// ... 类似修改，直接使用 d.db.PingContext ...
}
```

### 步骤 3：修改 Daemon 相关代码

#### 3.1 修改 `internal/biz/daemon/table_consumer.go`

```go
// BuildDefaultFetcher 构建默认的记录获取函数
// 参数：
//   - db: 数据库连接（*sql.DB）
//   - tableName: 表名
//   - statusField: 状态字段名
//   - pendingStatus: 待处理状态值
//   - orderBy: 排序字段（如 "created_at ASC"）
//   - scanFunc: 扫描函数，用于将数据库行转换为业务对象
func BuildDefaultFetcher(
	db *sql.DB,  // 改为 *sql.DB
	tableName string,
	statusField string,
	pendingStatus string,
	orderBy string,
	scanFunc func(*sql.Rows) (interface{}, error),
) RecordFetcher {
	return func(ctx context.Context, limit int) ([]interface{}, error) {
		// 设置查询超时
		queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// 构建 SQL 查询
		query := fmt.Sprintf(
			"SELECT * FROM %s WHERE %s = ? ORDER BY %s LIMIT ? FOR UPDATE SKIP LOCKED",
			tableName, statusField, orderBy,
		)

		// 直接使用 *sql.DB 执行查询
		rows, err := db.QueryContext(queryCtx, query, pendingStatus, limit)
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		var records []interface{}
		for rows.Next() {
			record, err := scanFunc(rows)
			if err != nil {
				return nil, fmt.Errorf("failed to scan record: %w", err)
			}
			records = append(records, record)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		return records, nil
	}
}
```

#### 3.2 修改 `internal/biz/daemon/table_consumer_ants.go`

```go
// TableConsumerDaemonAnts 基于 ants 的数据库表轮询守护进程
type TableConsumerDaemonAnts struct {
	*BaseDaemon
	db              *sql.DB  // 改为 *sql.DB
	handler         RecordHandler
	fetcher         RecordFetcher
	// ... 其他字段保持不变 ...
}

// NewTableConsumerDaemonAnts 创建基于 ants 的表轮询守护进程
func NewTableConsumerDaemonAnts(
	db *sql.DB,  // 改为 *sql.DB
	handler RecordHandler,
	logger log.Logger,
	opts ...TableConsumerAntsOption,
) (*TableConsumerDaemonAnts, error) {
	// ... 实现保持不变 ...
}
```

#### 3.3 修改 Daemon 使用示例

```go
// internal/biz/daemon/example.go
import (
	"database/sql"
	// 移除 gorm.io/gorm
)

func ExampleTableConsumer(db *sql.DB, logger log.Logger) (*TableConsumerDaemon, error) {
	// ... 使用 *sql.DB 而不是 *gorm.DB ...
}
```

### 步骤 4：更新 Wire 依赖注入

`cmd/sre/wire.go` 中的 `NewDB` 函数签名会自动更新，无需手动修改。

### 步骤 5：更新 go.mod

移除 GORM 相关依赖：

```bash
go mod edit -droprequire gorm.io/gorm
go mod edit -droprequire gorm.io/driver/mysql
go mod tidy
```

### 步骤 6：处理 GORM Logger（可选）

如果 `internal/logger/gorm.go` 只用于 GORM，可以移除。如果其他地方需要，可以保留但标记为废弃。

## 迁移检查清单

### 代码修改

- [ ] 修改 `internal/data/data.go`
  - [ ] 移除 GORM 导入
  - [ ] 添加 `database/sql` 导入
  - [ ] 修改 `Data` 结构体
  - [ ] 修改 `NewDB` 函数
  - [ ] 修改 `NewEntClient` 函数
  - [ ] 修改 `NewData` 函数

- [ ] 修改 `internal/data/health.go`
  - [ ] 更新健康检查逻辑使用 `*sql.DB`

- [ ] 修改 `internal/biz/daemon/table_consumer.go`
  - [ ] 修改 `BuildDefaultFetcher` 函数签名

- [ ] 修改 `internal/biz/daemon/table_consumer_ants.go`
  - [ ] 修改 `TableConsumerDaemonAnts` 结构体
  - [ ] 修改 `NewTableConsumerDaemonAnts` 函数

- [ ] 修改所有使用 Daemon 的代码
  - [ ] 更新函数签名和调用

### 依赖清理

- [ ] 从 `go.mod` 移除 GORM 依赖
- [ ] 运行 `go mod tidy`
- [ ] 检查是否有其他文件引用 GORM

### 测试验证

- [ ] 编译通过
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 健康检查功能正常
- [ ] Daemon 功能正常

## 迁移后的优势

1. **减少依赖**
   - 移除 `gorm.io/gorm`
   - 移除 `gorm.io/driver/mysql`
   - 只保留 `github.com/go-sql-driver/mysql`（Ent 需要）

2. **简化代码**
   - 连接管理更直接
   - 减少抽象层

3. **统一 ORM**
   - 所有数据库操作都通过 Ent
   - 代码风格更一致

4. **性能提升**
   - 减少一层抽象
   - 二进制体积更小

## 风险评估

### 低风险

- ✅ Repository 层已经只使用 Ent
- ✅ GORM 只用于连接管理
- ✅ 没有使用 GORM 的高级功能

### 需要注意

- ⚠️ Daemon 代码需要修改
- ⚠️ 需要测试健康检查功能
- ⚠️ 需要测试连接池配置

## 建议

**推荐进行迁移**，因为：
1. GORM 在当前项目中只是"连接管理器"
2. 所有业务逻辑都使用 Ent
3. 迁移成本低，收益明显
4. 可以统一技术栈

## 参考资源

- [database/sql 官方文档](https://pkg.go.dev/database/sql)
- [Ent 官方文档](https://entgo.io/)
- [MySQL 驱动文档](https://github.com/go-sql-driver/mysql)

