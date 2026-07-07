# 配置管理与数据库选型

## 一、配置管理

### 1.1 方案对比

| 方案 | 复杂度 | 热更新 | 多格式 | 远程配置 | 适用 |
|------|--------|--------|--------|----------|------|
| `os.Getenv` + 默认值 | 极低 | 需要重启 | - | - | 简单服务、Docker 部署 |
| `envconfig` / `cleanenv` | 低 | 需要重启 | - | - | 仅需环境变量，自动映射 struct |
| `Viper` | 中 | 支持 | yaml/json/toml/env | 支持(etcd/consul) | 需要多来源、多格式的项目 |
| `Koanf` | 中 | 支持 | 多格式 | 支持 | Viper 替代，更轻量、无全局状态 |
| 远程配置中心 (Nacos/Apollo) | 高 | 原生支持 | 多格式 | 原生 | 大型微服务体系 |

### 1.2 os.Getenv + 默认值（最简单）

```go
func getEnv(key, defaultVal string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return defaultVal
}

dbHost := getEnv("DB_HOST", "localhost")
dbPort := getEnv("DB_PORT", "5432")
```

**优点**：零依赖、Docker/K8s 天然支持、12-Factor App 标准。
**缺点**：多个配置时代码冗长、不支持复杂嵌套结构、无类型校验。
**适用**：配置项 < 10 个的微服务。

### 1.3 envconfig / cleanenv（环境变量到 Struct）

```go
// cleanenv
type Config struct {
    ServerPort int    `env:"SERVER_PORT" env-default:"8080"`
    DBHost     string `env:"DB_HOST" env-default:"localhost"`
    DBPort     int    `env:"DB_PORT" env-default:"5432"`
    LogLevel   string `env:"LOG_LEVEL" env-default:"info"`
}

var cfg Config
cleanenv.ReadEnv(&cfg) // 自动解析 + 默认值
```

**优点**：保留环境变量的简洁性，同时类型安全、自动映射到 struct。
**缺点**：仍需要重启或重新读取环境变量才能变更。
**适用**：Docker/K8s 部署的中型服务，配置简单不想引入文件。

### 1.4 Viper（功能最全）

```go
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath("./configs")
viper.AutomaticEnv()               // 环境变量可覆盖文件中的值
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

if err := viper.ReadInConfig(); err != nil {
    panic(err)
}

var cfg Config
viper.Unmarshal(&cfg)

// 热更新回调
viper.OnConfigChange(func(e fsnotify.Event) {
    viper.Unmarshal(&cfg)
    log.Println("config updated:", e.Name)
})
viper.WatchConfig()
```

**Viper vs Koanf**：

| 维度 | Viper | Koanf |
|------|-------|-------|
| 全局状态 | 有（viper 包级函数） | 无（New() 返回实例） |
| 性能 | 中等 | 优于 Viper |
| 维护 | 活跃但 issue 积压多 | 活跃 |
| 生态 | 更广泛 | 较新 |
| JSON 数字精度 | 有坑，float64 转换问题 | 无此问题 |

**选择**：如果新项目二选一，**Koanf** 更干净（无全局状态）；如果团队已有 Viper 经验，Viper 也够用。

### 1.5 配置优先级（推荐）

```
命令行参数 > 环境变量 > 配置文件 > 默认值
```

这符合 12-Factor App 的原则，Docker/K8s 通过环境变量覆盖配置文件，本地开发直接用配置文件。

---

## 二、数据库与 ORM 选型

### 2.1 ORM vs 原生 SQL：思维框架

Go 社区没有像 Java Hibernate / Python SQLAlchemy 那样"一统天下"的 ORM，选型核心问题是：**你愿意放弃多少 SQL 控制，换取多少开发效率。**

### 2.2 方案对比

| 方案 | 类型 | 学习成本 | SQL 可控性 | 开发速度 | 类型安全 |
|------|------|----------|------------|----------|----------|
| `database/sql` + `sqlx` | SQL Builder | 低 | 极高 | 慢 | 运行时 |
| **GORM** | 全功能 ORM | 中 | 低（复杂查询难控） | 快 | 运行时 |
| **Ent** (Facebook) | 代码生成 ORM | 高 | 中 | 快（生成后） | 编译时 |
| **Bun** | 混合 ORM | 中 | 高 | 中 | 运行时 |
| **sqlc** | 代码生成，SQL-first | 低 | 极高 | 快 | 编译时 |
| **Bob** | 代码生成 ORM | 中 | 高 | 中 | 编译时 |

### 2.3 原生 SQL + sqlx（完全控制）

```go
import "github.com/jmoiron/sqlx"

type User struct {
    ID   int64  `db:"id"`
    Name string `db:"name"`
}

var users []User
err := db.Select(&users, "SELECT * FROM users WHERE age > $1", 18)

// 命名参数
query := `SELECT * FROM users WHERE name = :name`
rows, err := db.NamedQuery(query, map[string]interface{}{"name": "Alice"})
```

**优点**：
- 完全掌控 SQL，不会生成意外的查询
- 零学习成本（就是写 SQL）
- 轻量，编译快
- 复杂 JOIN、子查询、窗口函数毫无限制

**缺点**：
- 大量 CRUD 时重复代码多
- 结构体与字段需要手动映射（sqlx 的 `db` tag 帮忙了但仍需手写）
- 没有 Migration 集成

**适用**：SQL 是强项、查询复杂（多表 join、报表）、或对 ORM 生成的 SQL 不放心的团队。

### 2.4 GORM（最流行，争议也最大）

```go
type User struct {
    gorm.Model
    Name  string
    Email string `gorm:"uniqueIndex"`
    Age   int
    Orders []Order // Has Many
}

// CRUD 极简
db.Create(&user)
db.First(&user, 1)
db.Where("age > ?", 18).Find(&users)
db.Model(&user).Update("name", "NewName")
db.Delete(&user)

// 关联查询
db.Preload("Orders").Find(&users)

// Auto Migration（仅开发环境）
db.AutoMigrate(&User{}, &Order{})
```

**优点**：
- CRUD 极快，中文社区资料最多
- AutoMigrate（开发阶段省心）
- 关联查询（Preload）方便
- Hook 机制（BeforeCreate/AfterUpdate）灵活

**缺点**：
- 复杂查询的 SQL 可能低效，需要查日志确认
- **软删除**（`gorm.DeletedAt`）默认行为让很多新手踩坑——`DELETE` 变成 `UPDATE`
- 零值（`""`, `0`, `false`）更新问题——Updates 默认忽略零值
- N+1 查询问题——Preload 忘了加就会出现
- 生成 SQL 的调试需要额外步骤

**适用**：常规 CRUD 为主、团队新手多、追求开发速度 > SQL 精细控制。

### 2.5 Ent（类型安全的代码生成 ORM）

```go
// 1. 定义 Schema（ent/schema/user.go）
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id"),
        field.String("name"),
        field.String("email").Unique(),
    }
}
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("orders", Order.Type),
    }
}

// 2. 运行 go generate ./ent 生成类型安全的客户端

// 3. 使用（编译时类型检查）
users, err := client.User.Query().
    Where(user.AgeGT(18)).
    WithOrders().
    All(ctx)
```

**优点**：
- 编译时类型安全——如果改了 Schema 但没更新查询代码，编译就报错
- 代码生成——不需要手动字符串映射
- 天然的 GraphQL 集成（生成 GraphQL schema）
- 自动 Migration（基于 Schema diff）

**缺点**：
- 学习成本高——Schema DSL 需要时间熟悉
- 生成代码量大，编译变慢
- 复杂 SQL（CTE、窗口函数）需要回退到原生 SQL
- 文件多，项目仓库体积增大

**适用**：数据模型复杂且频繁变动、需要 GraphQL、追求类型安全的团队。

### 2.6 sqlc（SQL-first 代码生成）

```sql
-- query.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users WHERE age > $1 ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;
```

```bash
sqlc generate  # 生成类型安全的 Go 函数
```

```go
// 生成的代码（编译时类型安全，SQL 完全可控）
user, err := queries.GetUser(ctx, 1)
users, err := queries.ListUsers(ctx, 18)
```

**优点**：
- SQL 完全掌控（你写的就是最终执行的 SQL）
- 编译时类型安全（返回 struct 字段在编译时检查）
- 生成代码干净、高性能，无反射

**缺点**：
- 动态查询（多条件筛选、排序）需要处理（可以用条件子句或备选方案）
- 没有 Migration 集成
- 只支持 PostgreSQL / MySQL / SQLite

**适用**：SQL 优先、追求类型安全和性能、不喜欢传统 ORM 魔法。

我推荐项目中**sqlc 做查询 + golang-migrate 做迁移**的组合——SQL 完全可控，编译时类型安全，零运行时开销。

### 2.7 最终选型决策

```
你主要做 CRUD？
    ├── 是 → GORM（最快上手）
    └── 否 → 你的 SQL 控制需求有多高？
             ├── 极高（复杂报表、多表 JOIN 为主）→ sqlx / sqlc
             ├── 高（混合场景，想保留 SQL 能力）→ Bun
             └── 中（类型安全 > SQL 控制）→ Ent

团队新手多、需要快速交付？
    └── 是 → GORM（资料最多，上手最快）

需要 GraphQL / 数据模型复杂多变？
    └── 是 → Ent（原生支持，编译时安全）

追求极简 + 完全控制？
    └── sqlc + golang-migrate（最佳组合）
```

---

## 三、数据库迁移

### 3.1 方案对比

| 工具 | 方式 | 适用 | 特点 |
|------|------|------|------|
| **golang-migrate** | SQL 文件 | 通用，最流行 | 纯 SQL，简单可靠 |
| **Atlas** (Ariga) | SQL / HCL DSL | 想要声明式管理 | 自动 diff、声明式迁移 |
| **goose** | SQL / Go 函数 | 需要 Go 代码逻辑的迁移 | 迁移中可写 Go 代码 |
| **GORM AutoMigrate** | 代码驱动 | 开发阶段 | 只建议开发用，生产不用 |
| **Ent Migrate** | Schema diff | Ent 项目 | 基于 Schema DSL 自动迁移 |

### 3.2 golang-migrate（推荐，最通用）

```bash
# 创建迁移
migrate create -ext sql -dir migrations -seq create_users_table
# → migrations/000001_create_users_table.up.sql
#   migrations/000001_create_users_table.down.sql

# 执行
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" up 1   # 只执行一个版本
migrate -path migrations -database "$DATABASE_URL" down   # 回退一个版本

# 在 Go 代码中嵌入迁移
import "github.com/golang-migrate/migrate/v4"
m, _ := migrate.New("file://migrations", databaseURL)
m.Up()
```

**优点**：纯 SQL、语言不受限、支持 Docker 和 K8s init container、社区最大。

### 3.3 迁移最佳实践

- **Up 和 Down 必须成对**——确保每个版本都可以回滚
- **幂等性**——Up 用 `IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS`
- **破坏性变更分两步**——先加字段（Up），确认没问题后再建 Down 写删除
- **不要改已执行的迁移**——已部署的迁移视为不可变，有问题就建新迁移
- **CI/CD 中自动执行**——作为部署的第一步

---

## 四、缓存策略

### 4.1 缓存层选型

| 方案 | 适用 | 性能 | 持久化 |
|------|------|------|--------|
| **Redis** | 通用分布式缓存 | 高 | 可选（RDB/AOF） |
| **go-cache / freecache** | 单实例进程内缓存 | 极高 | 无 |
| **ristretto** (Dgraph) | 单实例，对命中率要求高 | 极高 | 无 |
| **Memcached** | 简单 KV、不需要数据结构 | 高 | 无 |

### 4.2 策略对比

| 策略 | 流程 | 适用 |
|------|------|------|
| **Cache-Aside**（最常用）| 查缓存→miss→查 DB→写缓存 | 读多写少 |
| **Write-Through** | 先写缓存→再写 DB→返回 | 写后需要立即读到 |
| **Write-Behind** | 先写缓存→异步批量写 DB | 写密集、允许短暂不一致 |

### 4.3 常见问题与避坑

| 问题 | 方案 |
|------|------|
| **缓存穿透**（查不存在的数据每次穿透到 DB） | 缓存空值（短 TTL）或布隆过滤器 |
| **缓存击穿**（热点 key 过期瞬间大量请求到 DB） | 互斥锁（singleflight）或永不过期 + 异步更新 |
| **缓存雪崩**（大量 key 同时过期） | TTL 加随机值（`baseTTL + rand.Intn(60)`） |
| **数据不一致** | 先删缓存再写 DB + 延迟双删 |

```go
// singleflight 防止击穿
import "golang.org/x/sync/singleflight"

var g singleflight.Group

func GetUser(ctx context.Context, id int64) (*User, error) {
    // 先从缓存取
    if user := getFromCache(id); user != nil {
        return user, nil
    }
    // 缓存未命中，用 singleflight 确保同一 key 只查一次 DB
    v, err, _ := g.Do(fmt.Sprintf("user:%d", id), func() (interface{}, error) {
        return db.GetUser(ctx, id)
    })
    if err != nil {
        return nil, err
    }
    user := v.(*User)
    setCache(id, user, 5*time.Minute)
    return user, nil
}
```
