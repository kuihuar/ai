# 配置管理

## 配置系统选择

项目支持两种配置管理方式：

1. **Kratos 配置系统**（默认）：使用 Protobuf 定义配置结构，支持 YAML 格式
2. **Viper 配置系统**（推荐）：功能更强大，支持多种配置格式、环境变量、配置热更新等

## 两种配置系统对比

### 代码使用对比

#### Kratos 配置系统（之前的方式）

```go
import (
    "sre/internal/conf"
    "github.com/go-kratos/kratos/v2/config"
    "github.com/go-kratos/kratos/v2/config/file"
)

func main() {
    // 创建配置源
    c := config.New(
        config.WithSource(
            file.NewSource(flagconf),
        ),
    )
    defer c.Close()

    // 加载配置
    if err := c.Load(); err != nil {
        panic(err)
    }

    // 扫描到结构体
    var bc conf.Bootstrap
    if err := c.Scan(&bc); err != nil {
        panic(err)
    }

    // 使用配置
    app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
}
```

#### Viper 配置系统（新的方式）

```go
import "sre/internal/config"

func main() {
    // 创建配置加载器
    loader := config.NewLoader()
    
    // 加载配置文件
    if err := loader.LoadFromFile(flagconf); err != nil {
        panic(err)
    }

    // 加载并转换为 Bootstrap
    bootstrap, err := loader.LoadBootstrap()
    if err != nil {
        panic(err)
    }

    // 使用配置（与 Kratos 方式完全兼容）
    app, cleanup, err := wireApp(bootstrap.Server, bootstrap.Data, logger)
}
```

### 功能对比表

| 功能特性 | Kratos 配置系统 | Viper 配置系统 |
|---------|----------------|---------------|
| **配置文件格式** | 主要支持 YAML | 支持 YAML、JSON、TOML、HCL、INI 等 |
| **环境变量支持** | ❌ 需要手动处理 | ✅ 自动支持，前缀 `SRE_` |
| **多配置文件** | ❌ 需要手动合并 | ✅ 支持多文件，自动合并（后面的覆盖前面的） |
| **配置热更新** | ❌ 需要自己实现 | ✅ 内置支持 `WatchConfig()` |
| **获取单个配置值** | ❌ 需要先 Scan 整个结构 | ✅ 直接 `GetString()`、`GetInt()` 等 |
| **配置验证** | ✅ Protobuf 类型安全 | ✅ 支持自定义验证 |
| **默认值** | ✅ 支持 | ✅ 支持 `SetDefault()` |
| **配置路径查找** | 需要指定完整路径 | ✅ 支持目录自动查找 |
| **与 Kratos 集成** | ✅ 原生支持 | ✅ 完全兼容，输出相同结构 |
| **代码复杂度** | 简单 | 简单（封装后） |

### 主要区别说明

#### 1. API 使用方式

**Kratos 方式：**
- 使用 `config.New()` 创建配置对象
- 使用 `c.Scan()` 将配置扫描到结构体
- 需要手动管理配置源

**Viper 方式：**
- 使用 `config.NewLoader()` 创建加载器
- 使用 `loader.LoadBootstrap()` 直接获取配置
- 封装了配置加载逻辑，使用更简单

#### 2. 环境变量支持

**Kratos 方式：**
```go
// 需要手动读取环境变量并设置
addr := os.Getenv("SERVER_HTTP_ADDR")
if addr != "" {
    // 手动设置到配置中
}
```

**Viper 方式：**
```bash
# 自动支持，无需代码修改
export SRE_SERVER_HTTP_ADDR=0.0.0.0:8080
# 配置文件中对应的值会被自动覆盖
```

#### 3. 多配置文件支持

**Kratos 方式：**
```go
// 需要手动合并多个配置源
c := config.New(
    config.WithSource(
        file.NewSource("configs/base.yaml"),
        file.NewSource("configs/config.prod.yaml"),
    ),
)
```

**Viper 方式：**
```go
// 一行代码支持多文件
loader.LoadFromPaths("configs/base.yaml", "configs/config.prod.yaml")
```

#### 4. 获取单个配置值

**Kratos 方式：**
```go
// 必须先 Scan 整个配置结构
var bc conf.Bootstrap
c.Scan(&bc)
addr := bc.Server.Http.Addr
```

**Viper 方式：**
```go
// 直接获取，无需加载整个结构
addr := loader.GetString("server.http.addr")
timeout := loader.GetInt("server.http.timeout")
```

#### 5. 配置热更新

**Kratos 方式：**
```go
// 需要自己实现文件监听和重新加载逻辑
// 通常需要额外的库（如 fsnotify）
```

**Viper 方式：**
```go
// 内置支持，一行代码
loader.WatchConfig(func() {
    log.Info("Configuration reloaded")
    bootstrap, _ := loader.LoadBootstrap()
    // 更新应用配置
})
```

### 迁移建议

如果你正在使用 Kratos 配置系统，迁移到 Viper 非常简单：

1. **配置文件格式不变**：YAML 格式完全兼容
2. **配置结构不变**：仍然使用 `conf.Bootstrap` 结构
3. **只需修改加载代码**：将 `config.New()` + `c.Scan()` 替换为 `config.NewLoader()` + `LoadBootstrap()`

**迁移示例：**

```go
// 之前（Kratos）
c := config.New(config.WithSource(file.NewSource(flagconf)))
defer c.Close()
c.Load()
var bc conf.Bootstrap
c.Scan(&bc)

// 之后（Viper）
loader := config.NewLoader()
loader.LoadFromFile(flagconf)
bootstrap, _ := loader.LoadBootstrap()
// bc 和 bootstrap 是相同的结构，可以直接替换使用
```

### 选择建议

- **使用 Kratos 配置系统**：如果项目简单，只需要基本的 YAML 配置加载
- **使用 Viper 配置系统**：如果需要环境变量、多配置文件、配置热更新等高级功能（推荐）

> 📖 **迁移指南**：如果你正在使用 Kratos 配置系统，想迁移到 Viper，请参考 [配置系统迁移指南](./config-migration-guide.md)

## Kratos 配置系统

Kratos 使用 Protobuf 定义配置结构，支持 YAML 格式配置文件。

### 配置定义

在 `internal/conf/conf.proto` 中定义配置结构：

```protobuf
message Bootstrap {
  Server server = 1;
  Data data = 2;
}

message Server {
  message HTTP {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration timeout = 3;
  }
  message GRPC {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration timeout = 3;
  }
  HTTP http = 1;
  GRPC grpc = 2;
}

message Data {
  message Database {
    string driver = 1;
    string source = 2;
  }
  Database database = 1;
}
```

### 配置文件

`configs/config.yaml`：

```yaml
server:
  http:
    network: tcp
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    network: tcp
    addr: 0.0.0.0:9000
    timeout: 1s
data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

## 配置管理原则

### 1. 环境分离
不同环境使用不同的配置文件：
- `configs/config.yaml` - 开发环境
- `configs/config.prod.yaml` - 生产环境
- `configs/config.test.yaml` - 测试环境

### 2. 敏感信息保护
- 敏感信息（密码、密钥）不提交到代码库
- 使用环境变量或密钥管理服务
- 配置文件模板化

### 3. 配置验证
- 启动时验证配置完整性
- 提供默认值
- 清晰的错误提示

### 4. 配置热更新
- 支持配置热更新（可选）
- 更新时验证配置有效性
- 记录配置变更日志

## Viper 配置系统

Viper 是一个功能强大的配置管理库，支持多种配置格式、环境变量、配置热更新等功能。

### 基本使用

#### 方式 1: 使用配置加载器（推荐）

```go
import "sre/internal/config"

loader := config.NewLoader()
if err := loader.LoadFromFile("configs/config.yaml"); err != nil {
    panic(err)
}

bootstrap, err := loader.LoadBootstrap()
if err != nil {
    panic(err)
}
```

#### 方式 2: 从文件直接加载（便捷方法）

```go
bootstrap, err := config.LoadBootstrapFromFile("configs/config.yaml")
if err != nil {
    panic(err)
}
```

### 环境变量支持

Viper 自动支持环境变量覆盖配置文件中的值。环境变量命名规则：

- 前缀：`SRE_`（可通过 `SetEnvPrefix` 修改）
- 分隔符：`.` 会被替换为 `_`
- 示例：
  - `server.http.addr` → `SRE_SERVER_HTTP_ADDR`
  - `data.database.source` → `SRE_DATA_DATABASE_SOURCE`

```bash
# 使用环境变量覆盖配置
export SRE_SERVER_HTTP_ADDR=0.0.0.0:8080
export SRE_DATA_DATABASE_SOURCE="user:pass@tcp(localhost:3306)/db"
```

### 多配置文件支持

支持加载多个配置文件，后面的配置会覆盖前面的：

```go
loader := config.NewLoader()
loader.LoadFromPaths(
    "configs/base.yaml",        // 基础配置
    "configs/config.prod.yaml", // 环境特定配置
)
bootstrap, err := loader.LoadBootstrap()
```

### 配置热更新

支持监听配置文件变化并自动重新加载：

```go
loader := config.NewLoader()
loader.LoadFromFile("configs/config.yaml")

// 监听配置变化
loader.WatchConfig(func() {
    log.Info("Configuration reloaded")
    // 重新加载配置
    bootstrap, _ := loader.LoadBootstrap()
    // 更新应用配置
})

bootstrap, err := loader.LoadBootstrap()
```

### 获取单个配置值

```go
loader := config.NewLoader()
loader.LoadFromFile("configs/config.yaml")

// 获取配置值
addr := loader.GetString("server.http.addr")
timeout := loader.GetInt("server.http.timeout")
enabled := loader.GetBool("feature.enabled")
```

### 解析到自定义结构体

```go
type CustomConfig struct {
    ServerAddr string `mapstructure:"server_addr"`
    Timeout    int    `mapstructure:"timeout"`
}

var customConfig CustomConfig
loader.UnmarshalKey("custom", &customConfig)
```

### 支持的配置格式

Viper 支持多种配置格式：
- YAML（默认）
- JSON
- TOML
- HCL
- INI
- 环境变量
- 命令行参数

### 在 main.go 中使用 Viper

可以替换或补充 Kratos 配置系统：

```go
import (
    "sre/internal/config"
    "sre/internal/conf"
)

func main() {
    // 使用 Viper 加载配置
    loader := config.NewLoader()
    if err := loader.LoadFromFile(flagconf); err != nil {
        panic(err)
    }
    
    bootstrap, err := loader.LoadBootstrap()
    if err != nil {
        panic(err)
    }
    
    // 后续使用 bootstrap 配置
    app, cleanup, err := wireApp(bootstrap.Server, bootstrap.Data, logger)
    // ...
}
```

## 最佳实践

1. **配置分层**：按功能模块组织配置
2. **类型安全**：使用 Protobuf 定义，保证类型安全
3. **文档完善**：为每个配置项添加注释说明
4. **版本管理**：配置变更要有版本记录
5. **环境变量优先**：生产环境优先使用环境变量，避免敏感信息泄露
6. **配置验证**：启动时验证配置完整性和有效性
7. **配置热更新**：对于需要动态调整的配置，使用配置热更新功能

