# 配置系统迁移指南：从 Kratos 到 Viper

本文档详细说明如何从 Kratos 配置系统迁移到 Viper 配置系统。

## 快速对比

### 代码变化

**之前（Kratos）：**
```go
c := config.New(
    config.WithSource(
        file.NewSource(flagconf),
    ),
)
defer c.Close()
c.Load()
var bc conf.Bootstrap
c.Scan(&bc)
```

**之后（Viper）：**
```go
loader := config.NewLoader()
loader.LoadFromFile(flagconf)
bootstrap, _ := loader.LoadBootstrap()
```

### 主要改进

1. ✅ **环境变量自动支持**：无需手动处理
2. ✅ **多配置文件支持**：一行代码加载多个配置
3. ✅ **配置热更新**：内置支持，无需额外实现
4. ✅ **获取单个值**：无需加载整个配置结构
5. ✅ **更多配置格式**：支持 JSON、TOML 等

## 迁移步骤

### 步骤 1: 更新导入

**之前：**
```go
import (
    "sre/internal/conf"
    "github.com/go-kratos/kratos/v2/config"
    "github.com/go-kratos/kratos/v2/config/file"
)
```

**之后：**
```go
import (
    "sre/internal/config"
    "sre/internal/conf"  // 仍然需要，用于 Bootstrap 结构
)
```

### 步骤 2: 替换配置加载代码

**之前：**
```go
func main() {
    flag.Parse()
    
    c := config.New(
        config.WithSource(
            file.NewSource(flagconf),
        ),
    )
    defer c.Close()

    if err := c.Load(); err != nil {
        panic(err)
    }

    var bc conf.Bootstrap
    if err := c.Scan(&bc); err != nil {
        panic(err)
    }

    app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
    // ...
}
```

**之后：**
```go
func main() {
    flag.Parse()
    
    loader := config.NewLoader()
    if err := loader.LoadFromFile(flagconf); err != nil {
        panic(err)
    }

    bootstrap, err := loader.LoadBootstrap()
    if err != nil {
        panic(err)
    }

    app, cleanup, err := wireApp(bootstrap.Server, bootstrap.Data, logger)
    // ...
}
```

### 步骤 3: 配置文件无需修改

✅ **配置文件格式完全兼容**，无需修改任何 YAML 文件。

```yaml
# configs/config.yaml - 保持不变
server:
  http:
    network: tcp
    addr: 0.0.0.0:8000
    timeout: 1s
```

### 步骤 4: 利用新功能（可选）

迁移后，你可以使用 Viper 的新功能：

#### 使用环境变量

```bash
# 设置环境变量，自动覆盖配置文件
export SRE_SERVER_HTTP_ADDR=0.0.0.0:8080
export SRE_DATA_DATABASE_SOURCE="user:pass@tcp(localhost:3306)/db"
```

#### 多配置文件

```go
loader := config.NewLoader()
loader.LoadFromPaths(
    "configs/base.yaml",
    "configs/config.prod.yaml",
)
bootstrap, _ := loader.LoadBootstrap()
```

#### 配置热更新

```go
loader := config.NewLoader()
loader.LoadFromFile(flagconf)

// 监听配置变化
loader.WatchConfig(func() {
    log.Info("Configuration reloaded")
    bootstrap, _ := loader.LoadBootstrap()
    // 更新应用配置
})

bootstrap, _ := loader.LoadBootstrap()
```

#### 获取单个配置值

```go
// 无需加载整个配置结构
addr := loader.GetString("server.http.addr")
timeout := loader.GetInt("server.http.timeout")
```

## 完整迁移示例

### main.go 完整对比

**之前（Kratos）：**
```go
package main

import (
    "flag"
    "os"
    "sre/internal/conf"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/config"
    "github.com/go-kratos/kratos/v2/config/file"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/middleware/tracing"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
    _ "go.uber.org/automaxprocs"
)

var (
    Name     string
    Version  string
    flagconf string
    id, _    = os.Hostname()
)

func init() {
    flag.StringVar(&flagconf, "conf", "../../configs", "config path")
}

func main() {
    flag.Parse()
    logger := log.With(log.NewStdLogger(os.Stdout),
        "ts", log.DefaultTimestamp,
        "caller", log.DefaultCaller,
        "service.id", id,
        "service.name", Name,
        "service.version", Version,
        "trace.id", tracing.TraceID(),
        "span.id", tracing.SpanID(),
    )
    
    c := config.New(
        config.WithSource(
            file.NewSource(flagconf),
        ),
    )
    defer c.Close()

    if err := c.Load(); err != nil {
        panic(err)
    }

    var bc conf.Bootstrap
    if err := c.Scan(&bc); err != nil {
        panic(err)
    }

    app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
    if err != nil {
        panic(err)
    }
    defer cleanup()

    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

**之后（Viper）：**
```go
package main

import (
    "flag"
    "os"
    "sre/internal/config"
    "sre/internal/conf"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/middleware/tracing"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
    _ "go.uber.org/automaxprocs"
)

var (
    Name     string
    Version  string
    flagconf string
    id, _    = os.Hostname()
)

func init() {
    flag.StringVar(&flagconf, "conf", "../../configs", "config path")
}

func main() {
    flag.Parse()
    logger := log.With(log.NewStdLogger(os.Stdout),
        "ts", log.DefaultTimestamp,
        "caller", log.DefaultCaller,
        "service.id", id,
        "service.name", Name,
        "service.version", Version,
        "trace.id", tracing.TraceID(),
        "span.id", tracing.SpanID(),
    )
    
    // 使用 Viper 加载配置
    loader := config.NewLoader()
    if err := loader.LoadFromFile(flagconf); err != nil {
        panic(err)
    }

    bootstrap, err := loader.LoadBootstrap()
    if err != nil {
        panic(err)
    }

    app, cleanup, err := wireApp(bootstrap.Server, bootstrap.Data, logger)
    if err != nil {
        panic(err)
    }
    defer cleanup()

    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

## 常见问题

### Q1: 配置文件需要修改吗？

**A:** 不需要。YAML 配置文件格式完全兼容，无需任何修改。

### Q2: 配置结构会变化吗？

**A:** 不会。仍然使用 `conf.Bootstrap` 结构，输出完全相同。

### Q3: 环境变量如何设置？

**A:** 使用 `SRE_` 前缀，配置路径中的 `.` 替换为 `_`：
- `server.http.addr` → `SRE_SERVER_HTTP_ADDR`
- `data.database.source` → `SRE_DATA_DATABASE_SOURCE`

### Q4: 可以同时使用两种配置系统吗？

**A:** 可以，但不推荐。建议统一使用一种配置系统。

### Q5: 迁移后性能有影响吗？

**A:** 几乎没有影响。Viper 性能优秀，且配置加载通常在启动时进行。

### Q6: 如何回退到 Kratos 配置系统？

**A:** 只需恢复原来的代码即可，配置文件无需修改。

## 迁移检查清单

- [ ] 更新导入语句
- [ ] 替换配置加载代码
- [ ] 测试配置加载是否正常
- [ ] 验证环境变量支持（如需要）
- [ ] 测试多配置文件（如需要）
- [ ] 测试配置热更新（如需要）
- [ ] 更新相关文档

## 总结

从 Kratos 配置系统迁移到 Viper 配置系统非常简单：

1. ✅ **配置文件无需修改**
2. ✅ **配置结构保持不变**
3. ✅ **只需修改加载代码**
4. ✅ **获得更多功能**

迁移后，你可以享受环境变量、多配置文件、配置热更新等强大功能，同时保持与现有代码的完全兼容。
