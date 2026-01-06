# Kratos Option 模式详解

## 一、什么是 Option 模式？

Option 模式（也叫函数式选项模式）是一种**灵活的配置方式**，通过函数参数来设置对象的配置，避免了：
- 构造函数参数过多
- 需要创建多个构造函数
- 配置结构体过于复杂

## 二、传统配置方式的问题

### 2.1 构造函数参数过多

```go
// ❌ 不好的设计：参数过多
func NewApp(id string, name string, version string, logger Logger, 
            registrar Registrar, servers []Server, ...) *App {
    // ...
}
```

**问题**:
- 参数太多，难以记忆
- 必须按顺序传递所有参数
- 可选参数需要传 `nil` 或零值

### 2.2 配置结构体

```go
// ❌ 不好的设计：配置结构体
type AppConfig struct {
    ID        string
    Name      string
    Version   string
    Logger    Logger
    Registrar Registrar
    Servers   []Server
    // ... 很多字段
}

func NewApp(config AppConfig) *App {
    // ...
}

// 使用
app := NewApp(AppConfig{
    ID:      "app-1",
    Name:    "myapp",
    Version: "v1.0.0",
    // 必须设置所有字段，即使不需要
})
```

**问题**:
- 必须设置所有字段
- 无法区分零值和未设置
- 配置结构体可能很大

## 三、Kratos Option 模式实现

### 3.1 核心设计

```go
// options.go
// Option 是一个函数类型，接收 options 指针
type Option func(o *options)

// options 是内部配置结构体
type options struct {
    id        string
    name      string
    version   string
    metadata  map[string]string
    endpoints []*url.URL
    
    ctx  context.Context
    sigs []os.Signal
    
    logger           log.Logger
    registrar        registry.Registrar
    registrarTimeout time.Duration
    stopTimeout      time.Duration
    servers          []transport.Server
    
    // 生命周期钩子
    beforeStart []func(context.Context) error
    beforeStop  []func(context.Context) error
    afterStart  []func(context.Context) error
    afterStop   []func(context.Context) error
}

// New 创建 App，接收多个 Option
func New(opts ...Option) *App {
    o := options{
        // 设置默认值
        ctx:              context.Background(),
        sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
        registrarTimeout: 10 * time.Second,
    }
    
    // 应用所有 Option
    for _, opt := range opts {
        opt(&o)  // 每个 Option 函数修改 options
    }
    
    // 生成唯一ID（如果没有提供）
    if o.id == "" {
        if id, err := uuid.NewUUID(); err == nil {
            o.id = id.String()
        }
    }
    
    ctx, cancel := context.WithCancel(o.ctx)
    return &App{
        ctx:    ctx,
        cancel: cancel,
        opts:   o,
    }
}
```

### 3.2 Option 函数定义

```go
// ID 设置服务ID
func ID(id string) Option {
    return func(o *options) {
        o.id = id
    }
}

// Name 设置服务名称
func Name(name string) Option {
    return func(o *options) {
        o.name = name
    }
}

// Version 设置版本
func Version(version string) Option {
    return func(o *options) {
        o.version = version
    }
}

// Server 设置传输服务器
func Server(srv ...transport.Server) Option {
    return func(o *options) {
        o.servers = srv  // 可以设置多个服务器
    }
}

// Registrar 设置服务注册中心
func Registrar(r registry.Registrar) Option {
    return func(o *options) {
        o.registrar = r
    }
}

// Logger 设置日志
func Logger(logger log.Logger) Option {
    return func(o *options) {
        o.logger = logger
    }
}

// BeforeStart 设置启动前钩子
func BeforeStart(fn func(context.Context) error) Option {
    return func(o *options) {
        o.beforeStart = append(o.beforeStart, fn)
    }
}

// AfterStart 设置启动后钩子
func AfterStart(fn func(context.Context) error) Option {
    return func(o *options) {
        o.afterStart = append(o.afterStart, fn)
    }
}
```

## 四、Option 模式使用示例

### 4.1 基础使用

```go
package main

import (
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/go-kratos/kratos/v2/transport/grpc"
)

func main() {
    // 创建 HTTP 服务器
    httpSrv := http.NewServer(
        http.Address(":8000"),
    )
    
    // 创建 gRPC 服务器
    grpcSrv := grpc.NewServer(
        grpc.Address(":9000"),
    )
    
    // 使用 Option 模式创建 App
    app := kratos.New(
        kratos.ID("myapp-001"),
        kratos.Name("myapp"),
        kratos.Version("v1.0.0"),
        kratos.Server(httpSrv, grpcSrv),  // 可以传多个服务器
    )
    
    app.Run()
}
```

### 4.2 完整配置示例

```go
package main

import (
    "context"
    "time"
    
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/registry"
    "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/go-kratos/kratos/v2/transport/grpc"
)

func main() {
    // 创建日志
    logger := log.NewStdLogger(os.Stdout)
    
    // 创建服务注册中心
    registrar := consul.New(consulClient)
    
    // 创建 HTTP 服务器
    httpSrv := http.NewServer(
        http.Address(":8000"),
        http.Middleware(
            logging.Server(),
            recovery.Recovery(),
        ),
    )
    
    // 创建 gRPC 服务器
    grpcSrv := grpc.NewServer(
        grpc.Address(":9000"),
        grpc.Middleware(
            logging.Server(),
            recovery.Recovery(),
        ),
    )
    
    // 使用 Option 模式创建 App
    app := kratos.New(
        // 基本信息
        kratos.ID("myapp-001"),
        kratos.Name("myapp"),
        kratos.Version("v1.0.0"),
        kratos.Metadata(map[string]string{
            "env": "production",
        }),
        
        // 服务器
        kratos.Server(httpSrv, grpcSrv),
        
        // 日志
        kratos.Logger(logger),
        
        // 服务注册
        kratos.Registrar(registrar),
        kratos.RegistrarTimeout(5 * time.Second),
        
        // 生命周期钩子
        kratos.BeforeStart(func(ctx context.Context) error {
            log.Info("应用启动前...")
            // 初始化数据库连接等
            return nil
        }),
        kratos.AfterStart(func(ctx context.Context) error {
            log.Info("应用启动后...")
            // 启动后台任务等
            return nil
        }),
        kratos.BeforeStop(func(ctx context.Context) error {
            log.Info("应用停止前...")
            // 清理资源
            return nil
        }),
        kratos.AfterStop(func(ctx context.Context) error {
            log.Info("应用停止后...")
            return nil
        }),
    )
    
    app.Run()
}
```

### 4.3 自定义 Option

```go
package main

import (
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/transport/http"
)

// CustomOption 自定义 Option
func CustomOption(customValue string) kratos.Option {
    return func(o *kratos.Options) {
        // 可以访问内部 options 结构
        // 注意：这里需要 kratos 暴露 Options 类型，或者通过其他方式
        // 实际使用中，可以通过扩展 options 结构来支持
    }
}

// 或者通过包装其他 Option
func WithDatabase(dsn string) kratos.Option {
    return func(o *kratos.Options) {
        // 在启动前初始化数据库
        kratos.BeforeStart(func(ctx context.Context) error {
            // 初始化数据库
            db, err := sql.Open("mysql", dsn)
            if err != nil {
                return err
            }
            // 存储到 context 或其他地方
            return nil
        })(o)
    }
}
```

## 五、HTTP 服务器的 Option 模式

Kratos 的 HTTP 服务器也使用 Option 模式：

```go
// transport/http/server.go
type ServerOption func(*Server)

// Network 设置网络类型
func Network(network string) ServerOption {
    return func(s *Server) {
        s.network = network
    }
}

// Address 设置地址
func Address(addr string) ServerOption {
    return func(s *Server) {
        s.address = addr
    }
}

// Middleware 设置中间件
func Middleware(m ...middleware.Middleware) ServerOption {
    return func(o *Server) {
        o.middleware.Use(m...)
    }
}

// NewServer 创建服务器
func NewServer(opts ...ServerOption) *Server {
    srv := &Server{
        network:     "tcp",
        address:     ":0",
        timeout:     1 * time.Second,
        middleware:  matcher.New(),
        // ... 默认值
    }
    
    // 应用所有 Option
    for _, o := range opts {
        o(srv)
    }
    
    return srv
}
```

**使用示例**:
```go
httpSrv := http.NewServer(
    http.Address(":8000"),
    http.Timeout(5 * time.Second),
    http.Middleware(
        logging.Server(),
        recovery.Recovery(),
        metrics.Server(),
    ),
    http.RequestDecoder(customDecoder),
    http.ResponseEncoder(customEncoder),
)
```

## 六、Option 模式的优势

### 6.1 灵活性

- ✅ **可选参数**: 只设置需要的参数
- ✅ **任意顺序**: Option 可以任意顺序传递
- ✅ **组合使用**: 可以组合多个 Option

### 6.2 可扩展性

- ✅ **易于扩展**: 添加新配置只需添加新的 Option 函数
- ✅ **向后兼容**: 新 Option 不影响现有代码
- ✅ **类型安全**: 编译时检查类型

### 6.3 可读性

- ✅ **自文档化**: Option 函数名清晰表达意图
- ✅ **链式调用**: 可以链式调用多个 Option
- ✅ **IDE 支持**: IDE 可以自动补全

### 6.4 可测试性

- ✅ **易于 Mock**: 可以轻松创建测试用的 Option
- ✅ **依赖注入**: 通过 Option 注入依赖
- ✅ **配置隔离**: 不同测试可以使用不同配置

## 七、Option 模式的实现细节

### 7.1 默认值处理

```go
func New(opts ...Option) *App {
    o := options{
        // 设置合理的默认值
        ctx:              context.Background(),
        sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
        registrarTimeout: 10 * time.Second,
    }
    
    // 应用 Option（会覆盖默认值）
    for _, opt := range opts {
        opt(&o)
    }
    
    // 处理特殊情况
    if o.id == "" {
        // 自动生成ID
        if id, err := uuid.NewUUID(); err == nil {
            o.id = id.String()
        }
    }
    
    return &App{opts: o}
}
```

### 7.2 累积性 Option

某些 Option 可以累积（如中间件、钩子函数）：

```go
// BeforeStart 可以多次调用，累积多个函数
func BeforeStart(fn func(context.Context) error) Option {
    return func(o *options) {
        o.beforeStart = append(o.beforeStart, fn)  // 追加，不是覆盖
    }
}

// 使用
app := kratos.New(
    kratos.BeforeStart(func(ctx context.Context) error {
        // 第一个钩子
        return initDatabase(ctx)
    }),
    kratos.BeforeStart(func(ctx context.Context) error {
        // 第二个钩子
        return initCache(ctx)
    }),
)
```

### 7.3 条件性 Option

```go
// 根据环境变量决定是否启用某些功能
func NewApp() *kratos.App {
    opts := []kratos.Option{
        kratos.Name("myapp"),
    }
    
    // 只在生产环境启用服务注册
    if os.Getenv("ENV") == "production" {
        opts = append(opts, kratos.Registrar(consul.New(...)))
    }
    
    return kratos.New(opts...)
}
```

## 八、Option 模式的最佳实践

### 8.1 命名规范

- ✅ 使用动词开头：`With*`、`Set*`、`Enable*`
- ✅ 清晰表达意图：`Address`、`Timeout`、`Middleware`
- ✅ 保持一致性：所有 Option 函数命名风格一致

### 8.2 文档注释

```go
// Address sets the server address.
// If not set, defaults to ":0" (random port).
func Address(addr string) ServerOption {
    return func(s *Server) {
        s.address = addr
    }
}
```

### 8.3 类型安全

```go
// ✅ 好的设计：类型安全
func Timeout(timeout time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

// ❌ 不好的设计：类型不安全
func Timeout(timeout interface{}) ServerOption {
    return func(s *Server) {
        s.timeout = timeout.(time.Duration)  // 可能 panic
    }
}
```

### 8.4 验证和错误处理

```go
func Address(addr string) ServerOption {
    return func(s *Server) {
        // 验证地址格式
        if _, err := net.ResolveTCPAddr("tcp", addr); err != nil {
            panic(fmt.Sprintf("invalid address: %v", err))
        }
        s.address = addr
    }
}
```

## 九、总结

Kratos 的 Option 模式通过**函数式配置**实现了框架的扩展性：

1. **函数类型**: `Option func(o *options)`
2. **默认值**: 在 `New()` 中设置默认值
3. **应用 Option**: 遍历所有 Option 函数并应用
4. **灵活配置**: 只设置需要的参数，任意顺序

**适用场景**:
- ✅ 配置项较多的情况
- ✅ 需要可选参数
- ✅ 需要向后兼容
- ✅ 需要类型安全

**优势**:
- ✅ 灵活性高
- ✅ 可读性好
- ✅ 易于扩展
- ✅ 类型安全

Option 模式是 Go 语言中非常流行的配置模式，Kratos 将其应用得淋漓尽致，让框架的使用变得非常灵活和优雅。

