# GoFrame 与 Kratos 深入对比分析

## 一、核心架构设计深入分析

### 1.1 实例管理机制

#### GoFrame: 全局单例 + 懒加载

**实现原理**:
```go
// frame/gins/gins.go
var (
    instances = gins.NewInstances()  // 全局实例管理器
)

// frame/g/g_object.go
func Server(name ...any) *ghttp.Server {
    return gins.Server(name...)  // 委托给 gins
}

// net/ghttp/ghttp_server.go
func GetServer(name ...any) *Server {
    serverName := DefaultServerName
    if len(name) > 0 && name[0] != "" {
        serverName = gconv.String(name[0])
    }
    // 线程安全的单例获取
    v := serverMapping.GetOrSetFuncLock(serverName, func() any {
        s := &Server{
            instance:         serverName,
            plugins:          make([]Plugin, 0),
            servers:          make([]*graceful.Server, 0),
            // ... 初始化
        }
        // 使用默认配置初始化
        if err := s.SetConfig(NewConfig()); err != nil {
            panic(gerror.WrapCode(gcode.CodeInvalidConfiguration, err, ""))
        }
        // 默认启用 OpenTelemetry
        s.Use(internalMiddlewareServerTracing)
        return s
    })
    return v.(*Server)
}
```

**设计特点**:
1. **线程安全**: 使用 `GetOrSetFuncLock` 确保并发安全
2. **懒加载**: 首次访问时才初始化，节省资源
3. **多实例支持**: 通过名称区分不同实例
4. **配置驱动**: 自动从配置文件加载
5. **默认行为**: 自动启用常用功能（如追踪）

**优势**:
- ✅ 使用简单：`g.Server()` 即可获取
- ✅ 零配置：自动加载配置
- ✅ 资源高效：按需初始化

**劣势**:
- ❌ 隐式依赖：难以测试和替换
- ❌ 全局状态：可能造成并发问题（虽然已处理）
- ❌ 难以扩展：单例模式限制了灵活性

#### Kratos: 显式创建 + 生命周期管理

**实现原理**:
```go
// app.go
type App struct {
    opts     options
    ctx      context.Context
    cancel   context.CancelFunc
    mu       sync.Mutex
    instance *registry.ServiceInstance
}

func New(opts ...Option) *App {
    o := options{
        ctx:              context.Background(),
        sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
        registrarTimeout: 10 * time.Second,
    }
    // 生成唯一ID
    if id, err := uuid.NewUUID(); err == nil {
        o.id = id.String()
    }
    // 应用选项
    for _, opt := range opts {
        opt(&o)
    }
    ctx, cancel := context.WithCancel(o.ctx)
    return &App{
        ctx:    ctx,
        cancel: cancel,
        opts:   o,
    }
}

func (a *App) Run() error {
    // 1. 构建服务实例
    instance, err := a.buildInstance()
    
    // 2. 启动前钩子
    for _, fn := range a.opts.beforeStart {
        if err = fn(sctx); err != nil {
            return err
        }
    }
    
    // 3. 并发启动所有服务器
    eg, ctx := errgroup.WithContext(sctx)
    for _, srv := range a.opts.servers {
        server := srv
        eg.Go(func() error {
            <-ctx.Done()  // 等待停止信号
            return server.Stop(stopCtx)
        })
        eg.Go(func() error {
            return server.Start(octx)  // 启动服务器
        })
    }
    
    // 4. 服务注册
    if a.opts.registrar != nil {
        a.opts.registrar.Register(rctx, instance)
    }
    
    // 5. 启动后钩子
    for _, fn := range a.opts.afterStart {
        if err = fn(sctx); err != nil {
            return err
        }
    }
    
    // 6. 等待停止信号
    signal.Notify(c, a.opts.sigs...)
    eg.Go(func() error {
        select {
        case <-ctx.Done():
            return nil
        case <-c:
            return a.Stop()
        }
    })
    
    return eg.Wait()
}
```

**设计特点**:
1. **显式创建**: 必须显式调用 `kratos.New()`
2. **生命周期管理**: 完整的启动/停止流程
3. **并发启动**: 使用 `errgroup` 并发启动多个服务
4. **优雅停止**: 信号处理 + 超时控制
5. **钩子函数**: beforeStart/afterStart/beforeStop/afterStop

**优势**:
- ✅ 易于测试：可以注入 mock 对象
- ✅ 灵活扩展：Option 模式支持各种配置
- ✅ 生命周期清晰：明确的启动/停止流程
- ✅ 并发安全：使用 errgroup 管理并发

**劣势**:
- ❌ 使用复杂：需要显式创建和配置
- ❌ 代码量大：需要写更多代码

### 1.2 配置管理深入对比

#### GoFrame: 适配器模式

**实现原理**:
```go
// os/gcfg/gcfg.go
type Config struct {
    adapter Adapter  // 适配器接口
}

// 适配器接口
type Adapter interface {
    Available(ctx context.Context, resource ...string) bool
    Get(ctx context.Context, pattern string) (*gvar.Var, error)
    Data(ctx context.Context) (*gvar.Var, error)
}

// 文件适配器
type AdapterFile struct {
    // 文件监听、热更新
}

// Apollo 适配器
type AdapterApollo struct {
    // Apollo 客户端
}

// 使用
g.Cfg().MustGet(ctx, "database.default.link")
```

**特点**:
- ✅ 统一的接口，多种实现
- ✅ 支持热更新（文件监听）
- ✅ 懒加载配置
- ✅ 支持多种格式（YAML、TOML、JSON、INI）

**实现细节**:
```go
// 配置读取（带缓存）
func (c *Config) Get(ctx context.Context, pattern string, def ...any) (*gvar.Var, error) {
    // 1. 从适配器获取
    value, err := c.adapter.Get(ctx, pattern)
    if err != nil {
        // 2. 返回默认值
        if len(def) > 0 {
            return gvar.New(def[0]), nil
        }
        return nil, err
    }
    return value, nil
}
```

#### Kratos: 多源配置 + 观察者模式

**实现原理**:
```go
// config/config.go
type Config interface {
    Load() error
    Scan(v any) error
    Value(key string) Value
    Watch(key string, o Observer) error
    Close() error
}

type config struct {
    opts      options
    reader    Reader
    cached    sync.Map      // 缓存配置值
    observers sync.Map      // 观察者列表
    watchers  []Watcher     // 配置监听器
}

// 多源配置加载
func (c *config) Load() error {
    for _, src := range c.opts.sources {
        kvs, err := src.Load()  // 从各个源加载
        if err != nil {
            return err
        }
        // 合并配置（使用 mergo）
        if err := c.reader.Merge(kvs...); err != nil {
            return err
        }
    }
    // 解析环境变量
    return c.reader.Resolve()
}

// 配置监听
func (c *config) watch(w Watcher) {
    for {
        kvs, err := w.Next()  // 监听配置变更
        if err != nil {
            continue
        }
        // 合并新配置
        c.reader.Merge(kvs...)
        // 通知观察者
        c.cached.Range(func(key, value any) bool {
            k := key.(string)
            v := value.(Value)
            if n, ok := c.reader.Value(k); ok {
                if !reflect.DeepEqual(n.Load(), v.Load()) {
                    v.Store(n.Load())  // 更新缓存值
                    // 通知观察者
                    if o, ok := c.observers.Load(k); ok {
                        o.(Observer)(k, v)
                    }
                }
            }
            return true
        })
    }
}
```

**特点**:
- ✅ 多源配置合并（mergo）
- ✅ 观察者模式（配置变更通知）
- ✅ 类型安全（Scan 方法）
- ✅ 环境变量解析（Resolver）

**使用示例**:
```go
// 多源配置
c := config.New(
    config.WithSource(file.NewSource("config.yaml")),
    config.WithSource(etcd.NewSource(etcdClient)),
    config.WithSource(consul.NewSource(consulClient)),
)

// 观察配置变更
c.Watch("database.host", func(key string, value Value) {
    // 配置变更时自动更新连接
    db.UpdateConnection(value.String())
})
```

### 1.3 传输层抽象深入对比

#### GoFrame: HTTP 专用实现

**实现原理**:
```go
// net/ghttp/ghttp_server.go
type Server struct {
    instance         string
    config           *ServerConfig
    plugins          []Plugin
    servers          []*graceful.Server
    routesMap        map[string][]*HandlerItem
    middleware       []HandlerFunc
    // ...
}

// 路由注册
func (s *Server) Bind(controller interface{}) {
    // 通过反射分析 controller 的方法
    // 自动注册路由
}

// 中间件
func (s *Server) Use(middleware ...HandlerFunc) {
    s.middleware = append(s.middleware, middleware...)
}
```

**特点**:
- ✅ HTTP 专用，深度优化
- ✅ 自动路由注册（反射）
- ✅ 内置 OpenAPI 生成
- ✅ 支持 WebSocket

#### Kratos: 统一传输层抽象

**实现原理**:
```go
// transport/transport.go
type Server interface {
    Start(context.Context) error
    Stop(context.Context) error
}

type Transporter interface {
    Kind() Kind  // http/grpc
    Endpoint() string
    Operation() string
    RequestHeader() Header
    ReplyHeader() Header
}

// HTTP 实现
// transport/http/server.go
type Server struct {
    *http.Server
    middleware  matcher.Matcher
    router      *mux.Router
    // ...
}

func (s *Server) Start(ctx context.Context) error {
    // HTTP 服务器启动
}

// gRPC 实现
// transport/grpc/server.go
type Server struct {
    *grpc.Server
    middleware  matcher.Matcher
    // ...
}

func (s *Server) Start(ctx context.Context) error {
    // gRPC 服务器启动
}
```

**特点**:
- ✅ 统一的接口，多种实现
- ✅ HTTP 和 gRPC 通用中间件
- ✅ 统一的元数据传输
- ✅ 服务发现集成

**中间件统一**:
```go
// middleware/middleware.go
type Middleware func(Handler) Handler
type Handler func(ctx context.Context, req any) (any, error)

// HTTP 和 gRPC 都使用这个接口
func Recovery() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req any) (reply any, err error) {
            defer func() {
                if rerr := recover(); rerr != nil {
                    // 恢复逻辑
                    err = handlePanic(ctx, req, rerr)
                }
            }()
            return handler(ctx, req)
        }
    }
}
```

### 1.4 中间件设计深入对比

#### GoFrame: HTTP 中间件

**实现原理**:
```go
// net/ghttp/ghttp_request.go
type Request struct {
    Server     *Server
    Request    *http.Request
    Response   *ResponseWriter
    // ...
}

type HandlerFunc func(r *Request)

// 中间件执行
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. 创建 Request 对象
    request := s.newRequest(w, r)
    
    // 2. 执行中间件
    for _, middleware := range s.middleware {
        middleware(request)
        if request.IsExited() {
            return
        }
    }
    
    // 3. 执行路由处理器
    s.handleRequest(request)
}
```

**特点**:
- ✅ HTTP 专用
- ✅ 请求/响应拦截
- ✅ 支持分组中间件

#### Kratos: 函数式通用中间件

**实现原理**:
```go
// middleware/middleware.go
type Middleware func(Handler) Handler
type Handler func(ctx context.Context, req any) (any, error)

// 链式组合
func Chain(m ...Middleware) Middleware {
    return func(next Handler) Handler {
        for i := len(m) - 1; i >= 0; i-- {
            next = m[i](next)  // 从后往前包装
        }
        return next
    }
}

// 使用示例
middleware.Chain(
    logging.Server(),
    metrics.Server(),
    tracing.Server(),
    recovery.Recovery(),
)
```

**Recovery 中间件实现**:
```go
// middleware/recovery/recovery.go
func Recovery(opts ...Option) middleware.Middleware {
    op := options{
        handler: func(context.Context, any, any) error {
            return ErrUnknownRequest
        },
    }
    for _, o := range opts {
        o(&op)
    }
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req any) (reply any, err error) {
            startTime := time.Now()
            defer func() {
                if rerr := recover(); rerr != nil {
                    buf := make([]byte, 64<<10)
                    n := runtime.Stack(buf, false)
                    buf = buf[:n]
                    log.Context(ctx).Errorf("%v: %+v\n%s\n", rerr, req, buf)
                    ctx = context.WithValue(ctx, Latency{}, time.Since(startTime).Seconds())
                    err = op.handler(ctx, req, rerr)
                }
            }()
            return handler(ctx, req)
        }
    }
}
```

**特点**:
- ✅ 函数式设计，易于组合
- ✅ HTTP 和 gRPC 通用
- ✅ 类型安全（泛型支持）
- ✅ 易于测试

## 二、数据库/ORM 设计对比

### 2.1 GoFrame: 内置 ORM（gdb）

**实现原理**:
```go
// database/gdb/gdb_core.go
type Core struct {
    db     DB
    ctx    context.Context
    cache  *gcache.Cache
    links  *gmap.StrAnyMap  // 连接池
    // ...
}

// 链式调用
func (c *Core) Model(tableNameOrStruct ...any) *Model {
    return NewModel(c.db, tableNameOrStruct...)
}

// 上下文传递
func (c *Core) Ctx(ctx context.Context) DB {
    // 创建新的 DB 对象（浅拷贝）
    newCore := &Core{}
    *newCore = *c
    newCore.ctx = WithDB(ctx, newCore.db)
    return newCore.db
}

// 超时控制
func (c *Core) GetCtxTimeout(ctx context.Context, timeoutType ctxTimeoutType) (context.Context, context.CancelFunc) {
    var config = c.db.GetConfig()
    switch timeoutType {
    case ctxTimeoutTypeExec:
        if config.ExecTimeout > 0 {
            return context.WithTimeout(ctx, config.ExecTimeout)
        }
    case ctxTimeoutTypeQuery:
        if config.QueryTimeout > 0 {
            return context.WithTimeout(ctx, config.QueryTimeout)
        }
    // ...
    }
    return ctx, func() {}
}
```

**特点**:
- ✅ 内置 ORM，无需额外依赖
- ✅ 链式调用，API 简洁
- ✅ 上下文传递，支持超时控制
- ✅ 主从分离支持
- ✅ 连接池管理

### 2.2 Kratos: 无内置 ORM

**设计理念**:
- ❌ 不内置 ORM，推荐使用第三方库
- ✅ 推荐：GORM、Ent、sqlx
- ✅ 框架只关注传输层，不关注数据层

**原因**:
1. **关注点分离**: 框架只负责传输层
2. **灵活性**: 可以选择最适合项目的 ORM
3. **轻量级**: 框架本身更轻量

## 三、服务发现与负载均衡

### 3.1 GoFrame: 通过 contrib 支持

**实现方式**:
```go
// 通过 gsvc 包支持服务发现
import "github.com/gogf/gf/v2/net/gsvc"

// 服务注册
registrar := gsvc.GetRegistry()
registrar.Register(ctx, service)

// 服务发现
discovery := gsvc.GetRegistry()
instances, err := discovery.Search(ctx, gsvc.SearchInput{
    Name: "user-service",
})
```

**特点**:
- ✅ 通过 contrib 支持多种注册中心
- ✅ 统一的接口
- ✅ 可选功能

### 3.2 Kratos: 内置支持

**实现原理**:
```go
// registry/registry.go
type Registrar interface {
    Register(ctx context.Context, service *ServiceInstance) error
    Deregister(ctx context.Context, service *ServiceInstance) error
}

type Discovery interface {
    GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
    Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// selector/selector.go
type Selector interface {
    Rebalancer
    Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error)
}
```

**负载均衡实现**:
```go
// selector/p2c/p2c.go - Power of Two Choices
type p2cPicker struct {
    r     *rand.Rand
    mu    sync.Mutex
    nodes []*Node
}

func (p *p2cPicker) Pick(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error) {
    // 1. 随机选择两个节点
    // 2. 选择负载较低的节点
    // 3. 更新节点统计信息
}
```

**特点**:
- ✅ 内置服务发现
- ✅ 多种负载均衡算法（P2C、WRR、Random）
- ✅ 节点健康检查
- ✅ 自动故障转移

## 四、缓存设计对比

### 4.1 GoFrame: 内置缓存（gcache）

**实现原理**:
```go
// os/gcache/gcache.go
type Cache struct {
    adapter Adapter
}

// 适配器接口
type Adapter interface {
    Get(ctx context.Context, key any) (*gvar.Var, error)
    Set(ctx context.Context, key, value any, duration time.Duration) error
    // ...
}

// 内存适配器
type AdapterMemory struct {
    // LRU 缓存实现
}

// Redis 适配器
type AdapterRedis struct {
    redis *gredis.Redis
}

// 高级功能
func GetOrSetFunc(ctx context.Context, key any, f Func, duration time.Duration) (*gvar.Var, error) {
    // 1. 尝试获取
    value, err := cache.Get(ctx, key)
    if err == nil {
        return value, nil
    }
    // 2. 执行函数
    value, err = f(ctx)
    if err != nil {
        return nil, err
    }
    // 3. 设置缓存
    cache.Set(ctx, key, value, duration)
    return value, nil
}

// 带锁的版本（防止缓存击穿）
func GetOrSetFuncLock(ctx context.Context, key any, f Func, duration time.Duration) (*gvar.Var, error) {
    // 使用互斥锁确保只有一个 goroutine 执行函数
}
```

**特点**:
- ✅ 统一的缓存接口
- ✅ 支持内存和 Redis
- ✅ 防止缓存击穿（GetOrSetFuncLock）
- ✅ 自动序列化

### 4.2 Kratos: 无内置缓存

**设计理念**:
- ❌ 不内置缓存
- ✅ 推荐使用 go-redis、go-cache
- ✅ 框架只关注传输层

## 五、错误处理对比

### 5.1 GoFrame: 错误码系统

**实现原理**:
```go
// errors/gcode/gcode.go
type Code int

const (
    CodeNil Code = iota
    CodeInternalError
    CodeInvalidParameter
    CodeInvalidOperation
    // ...
)

// errors/gerror/gerror.go
type Error struct {
    code    Code
    message string
    stack   string
    cause   error
}

func NewCode(code Code, message string) error {
    return &Error{
        code:    code,
        message: message,
        stack:   gdebug.Stack(),
    }
}
```

**特点**:
- ✅ 统一的错误码
- ✅ 错误堆栈跟踪
- ✅ 错误链（cause）

### 5.2 Kratos: Protobuf 错误定义

**实现原理**:
```go
// errors/errors.proto
message Error {
    int32 code = 1;
    string reason = 2;
    string message = 3;
    map<string, string> metadata = 4;
}

// errors/errors.go
type Error struct {
    Code     int32
    Reason   string
    Message  string
    Metadata map[string]string
}

// 错误生成工具
// cmd/protoc-gen-go-errors
// 从 proto 文件生成错误代码
```

**特点**:
- ✅ Protobuf 定义错误
- ✅ 代码生成工具
- ✅ 跨语言错误定义
- ✅ 元数据支持

## 六、性能优化策略

### 6.1 GoFrame 性能优化

1. **连接池管理**:
```go
// 数据库连接池
type Core struct {
    links *gmap.StrAnyMap  // 连接池缓存
}

func (c *Core) getSqlDb(master bool, schema string) (*sql.DB, error) {
    // 1. 从连接池获取
    // 2. 如果不存在，创建新连接
    // 3. 缓存连接
}
```

2. **路由缓存**:
```go
// 路由匹配结果缓存
type Server struct {
    serveCache *gcache.Cache  // 路由缓存
}

func (s *Server) getHandlerByRequest(r *Request) {
    // 1. 检查缓存
    // 2. 如果未命中，匹配路由
    // 3. 缓存结果
}
```

3. **懒加载**:
```go
// 全局单例懒加载
var defaultCache = sync.OnceValue(func() *Cache {
    return New()
})
```

### 6.2 Kratos 性能优化

1. **连接池**:
```go
// HTTP 客户端连接池
type Client struct {
    *http.Client
    // 使用标准库的连接池
}
```

2. **负载均衡算法**:
```go
// P2C 算法（Power of Two Choices）
// 时间复杂度 O(1)，性能优于 WRR
func (p *p2cPicker) Pick(ctx context.Context) (Node, error) {
    // 随机选择两个节点，选择负载较低的
}
```

3. **并发控制**:
```go
// 使用 errgroup 管理并发
eg, ctx := errgroup.WithContext(ctx)
for _, srv := range servers {
    eg.Go(func() error {
        return srv.Start(ctx)
    })
}
```

## 七、扩展性设计

### 7.1 GoFrame: 插件系统

**实现原理**:
```go
// net/ghttp/ghttp_plugin.go
type Plugin interface {
    Install(s *Server) error
}

// 插件安装
func (s *Server) Use(plugin Plugin) {
    s.plugins = append(s.plugins, plugin)
}

// 启动时安装
for _, p := range s.plugins {
    if err := p.Install(s); err != nil {
        s.Logger().Fatalf(ctx, `%+v`, err)
    }
}
```

### 7.2 Kratos: Option 模式

**实现原理**:
```go
// 通过 Option 函数扩展
type Option func(*options)

func WithServer(srv ...transport.Server) Option {
    return func(o *options) {
        o.servers = append(o.servers, srv...)
    }
}

func WithRegistrar(reg registry.Registrar) Option {
    return func(o *options) {
        o.registrar = reg
    }
}
```

## 八、最佳实践总结

### 8.1 选择建议

**选择 GoFrame 如果**:
- ✅ 开发单体 Web 应用
- ✅ 需要快速开发，开箱即用
- ✅ 需要内置 ORM、缓存等工具
- ✅ 团队喜欢约定式开发

**选择 Kratos 如果**:
- ✅ 开发微服务架构
- ✅ 需要服务发现、负载均衡
- ✅ 需要可观测性（追踪、指标）
- ✅ 需要同时支持 HTTP 和 gRPC

### 8.2 混合使用

可以混合使用两个框架的优势：

```go
// 使用 Kratos 做微服务框架
app := kratos.New(
    kratos.Server(httpSrv, grpcSrv),
    kratos.Registrar(registry),
)

// 使用 GoFrame 的工具类
import "github.com/gogf/gf/v2/util/gvalid"
validator := gvalid.New()

// 使用 GoFrame 的 ORM
import "github.com/gogf/gf/v2/database/gdb"
db := gdb.New(...)
```

### 8.3 设计模式总结

| 设计模式 | GoFrame | Kratos | 学习价值 |
|---------|---------|--------|---------|
| **单例模式** | ✅ 全局单例 | ❌ 显式创建 | 懒加载、线程安全 |
| **适配器模式** | ✅ 配置/缓存适配器 | ✅ 配置源/注册中心适配器 | 可插拔设计 |
| **观察者模式** | ✅ 配置热更新 | ✅ 配置观察者 | 事件驱动 |
| **工厂模式** | ✅ 实例工厂 | ✅ Option 模式 | 对象创建 |
| **策略模式** | ✅ 负载均衡策略 | ✅ Selector 策略 | 算法封装 |
| **装饰器模式** | ✅ 中间件 | ✅ 中间件链 | 功能增强 |
| **模板方法** | ✅ 服务器启动流程 | ✅ App 生命周期 | 流程抽象 |

## 九、深入技术细节

### 9.1 上下文传递

#### GoFrame:
```go
// 上下文传递到数据库操作
func (c *Core) Ctx(ctx context.Context) DB {
    newCore := &Core{}
    *newCore = *c  // 浅拷贝
    newCore.ctx = WithDB(ctx, newCore.db)
    return newCore.db
}
```

#### Kratos:
```go
// 上下文传递到中间件
type Handler func(ctx context.Context, req any) (any, error)

// 中间件可以访问和修改上下文
func Logging() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req any) (any, error) {
            // 从上下文获取信息
            tr, _ := transport.FromServerContext(ctx)
            // 记录日志
            log.Context(ctx).Info("request", ...)
            return handler(ctx, req)
        }
    }
}
```

### 9.2 并发安全

#### GoFrame:
```go
// 使用 sync.Map 和 gmap
type Server struct {
    routesMap map[string][]*HandlerItem  // 需要加锁
    serveCache *gcache.Cache  // 线程安全
}

// 使用 GetOrSetFuncLock 确保线程安全
v := serverMapping.GetOrSetFuncLock(serverName, func() any {
    // 初始化逻辑
})
```

#### Kratos:
```go
// 使用 sync.Map
type config struct {
    cached    sync.Map      // 线程安全
    observers sync.Map      // 线程安全
}

// 使用 errgroup 管理并发
eg, ctx := errgroup.WithContext(ctx)
for _, srv := range servers {
    eg.Go(func() error {
        return srv.Start(ctx)
    })
}
```

## 十、总结

### GoFrame 的核心优势
1. **开箱即用**: 提供完整的工具链
2. **约定优于配置**: 减少配置，提高效率
3. **工具类丰富**: 字符串、数组、时间等工具
4. **简单易用**: 全局单例，使用方便
5. **内置功能**: ORM、缓存、日志等

### Kratos 的核心优势
1. **微服务导向**: 专为微服务设计
2. **可观测性**: 内置追踪、指标、日志
3. **传输层抽象**: 统一 HTTP/gRPC
4. **生命周期管理**: 优雅启动/停止
5. **服务治理**: 服务发现、负载均衡、熔断

### 可以学到的设计思想
1. **全局单例 + 懒加载**: 提高性能，简化使用
2. **适配器模式**: 可插拔设计，易于扩展
3. **应用生命周期管理**: 优雅启动/停止
4. **函数式中间件**: 易于组合和测试
5. **传输层抽象**: 统一不同协议
6. **观察者模式**: 配置热更新
7. **Option 模式**: 灵活的配置方式
8. **上下文传递**: 请求链路追踪

两个框架各有优势，可以根据项目需求选择，也可以混合使用，取长补短！

