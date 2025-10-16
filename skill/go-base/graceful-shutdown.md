# Go程序优雅退出详解

## 1. 概述

优雅退出（Graceful Shutdown）是指程序在收到退出信号时，能够：
- 停止接收新的请求
- 完成正在处理的请求
- 释放资源（关闭连接、文件等）
- 保存必要的数据
- 然后安全退出

## 2. 信号处理

### 2.1 常见的退出信号

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // 创建信号通道
    sigChan := make(chan os.Signal, 1)
    
    // 注册要捕获的信号
    signal.Notify(sigChan, 
        syscall.SIGINT,  // Ctrl+C
        syscall.SIGTERM, // 终止信号
        syscall.SIGHUP,  // 挂起信号（常用于重载配置）
    )
    
    // 等待信号
    sig := <-sigChan
    fmt.Printf("收到信号: %v\n", sig)
    
    // 执行清理工作
    cleanup()
}

func cleanup() {
    fmt.Println("执行清理工作...")
    // 关闭数据库连接
    // 保存数据
    // 释放资源
}
```

### 2.2 信号类型说明

| 信号 | 值 | 说明 | 用途 |
|------|----|----|----|
| SIGINT | 2 | 中断信号 | Ctrl+C，用户主动终止 |
| SIGTERM | 15 | 终止信号 | 系统或进程管理器发送 |
| SIGHUP | 1 | 挂起信号 | 终端关闭或配置重载 |
| SIGQUIT | 3 | 退出信号 | 生成core dump后退出 |
| SIGKILL | 9 | 强制杀死 | 无法捕获，立即终止 |

## 3. HTTP服务器优雅退出

### 3.1 基本实现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // 创建HTTP服务器
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.HandlerFunc(handler),
    }
    
    // 启动服务器（非阻塞）
    go func() {
        fmt.Println("服务器启动在 :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("服务器启动失败: %v", err)
        }
    }()
    
    // 等待退出信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("开始优雅关闭服务器...")
    
    // 创建超时上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 优雅关闭服务器
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("服务器强制关闭: %v", err)
    }
    
    fmt.Println("服务器已优雅关闭")
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 模拟长时间处理
    time.Sleep(2 * time.Second)
    w.Write([]byte("Hello, World!"))
}
```

### 3.2 带健康检查的优雅退出

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type Server struct {
    server     *http.Server
    isShutdown bool
    mu         sync.RWMutex
}

func NewServer() *Server {
    s := &Server{}
    
    mux := http.NewServeMux()
    mux.HandleFunc("/", s.handler)
    mux.HandleFunc("/health", s.healthHandler)
    
    s.server = &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    return s
}

func (s *Server) Start() {
    go func() {
        fmt.Println("服务器启动在 :8080")
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("服务器启动失败: %v", err)
        }
    }()
}

func (s *Server) Shutdown(ctx context.Context) error {
    s.mu.Lock()
    s.isShutdown = true
    s.mu.Unlock()
    
    return s.server.Shutdown(ctx)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
    // 检查是否正在关闭
    s.mu.RLock()
    if s.isShutdown {
        s.mu.RUnlock()
        http.Error(w, "服务器正在关闭", http.StatusServiceUnavailable)
        return
    }
    s.mu.RUnlock()
    
    // 模拟处理时间
    time.Sleep(2 * time.Second)
    w.Write([]byte("Hello, World!"))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
    s.mu.RLock()
    isShutdown := s.isShutdown
    s.mu.RUnlock()
    
    if isShutdown {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("服务器正在关闭"))
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("健康"))
}

func main() {
    server := NewServer()
    server.Start()
    
    // 等待退出信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("开始优雅关闭服务器...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("服务器强制关闭: %v", err)
    }
    
    fmt.Println("服务器已优雅关闭")
}
```

## 4. 数据库连接优雅关闭

### 4.1 基本实现

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
)

type App struct {
    db *sql.DB
}

func NewApp() (*App, error) {
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        return nil, err
    }
    
    // 设置连接池参数
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return &App{db: db}, nil
}

func (a *App) Start() {
    // 启动业务逻辑
    go a.processData()
}

func (a *App) processData() {
    for {
        // 模拟数据处理
        time.Sleep(1 * time.Second)
        fmt.Println("处理数据中...")
    }
}

func (a *App) Shutdown(ctx context.Context) error {
    fmt.Println("关闭数据库连接...")
    return a.db.Close()
}

func main() {
    app, err := NewApp()
    if err != nil {
        log.Fatal(err)
    }
    
    app.Start()
    
    // 等待退出信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("开始优雅关闭应用...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := app.Shutdown(ctx); err != nil {
        log.Printf("关闭应用失败: %v", err)
    }
    
    fmt.Println("应用已优雅关闭")
}
```

### 4.2 带事务处理的优雅关闭

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
)

type App struct {
    db           *sql.DB
    activeTx     map[string]*sql.Tx
    mu           sync.RWMutex
    isShutdown   bool
}

func NewApp() (*App, error) {
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        return nil, err
    }
    
    return &App{
        db:       db,
        activeTx: make(map[string]*sql.Tx),
    }, nil
}

func (a *App) StartTransaction(id string) (*sql.Tx, error) {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if a.isShutdown {
        return nil, fmt.Errorf("应用正在关闭")
    }
    
    tx, err := a.db.Begin()
    if err != nil {
        return nil, err
    }
    
    a.activeTx[id] = tx
    return tx, nil
}

func (a *App) CommitTransaction(id string) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    tx, exists := a.activeTx[id]
    if !exists {
        return fmt.Errorf("事务不存在")
    }
    
    delete(a.activeTx, id)
    return tx.Commit()
}

func (a *App) RollbackTransaction(id string) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    tx, exists := a.activeTx[id]
    if !exists {
        return fmt.Errorf("事务不存在")
    }
    
    delete(a.activeTx, id)
    return tx.Rollback()
}

func (a *App) Shutdown(ctx context.Context) error {
    a.mu.Lock()
    a.isShutdown = true
    a.mu.Unlock()
    
    // 等待所有事务完成或超时
    done := make(chan struct{})
    go func() {
        for {
            a.mu.RLock()
            if len(a.activeTx) == 0 {
                a.mu.RUnlock()
                close(done)
                return
            }
            a.mu.RUnlock()
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    select {
    case <-done:
        fmt.Println("所有事务已完成")
    case <-ctx.Done():
        fmt.Println("等待事务超时，强制回滚")
        a.mu.Lock()
        for id, tx := range a.activeTx {
            tx.Rollback()
            delete(a.activeTx, id)
        }
        a.mu.Unlock()
    }
    
    return a.db.Close()
}
```

## 5. 消息队列优雅关闭

### 5.1 RabbitMQ示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    
    "github.com/streadway/amqp"
)

type MessageConsumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   amqp.Queue
    wg      sync.WaitGroup
    quit    chan struct{}
}

func NewMessageConsumer() (*MessageConsumer, error) {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return nil, err
    }
    
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    q, err := ch.QueueDeclare(
        "task_queue", // name
        true,         // durable
        false,        // delete when unused
        false,        // exclusive
        false,        // no-wait
        nil,          // arguments
    )
    if err != nil {
        return nil, err
    }
    
    return &MessageConsumer{
        conn:    conn,
        channel: ch,
        queue:   q,
        quit:    make(chan struct{}),
    }, nil
}

func (mc *MessageConsumer) Start() {
    msgs, err := mc.channel.Consume(
        mc.queue.Name, // queue
        "",            // consumer
        false,         // auto-ack
        false,         // exclusive
        false,         // no-local
        false,         // no-wait
        nil,           // args
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 启动多个消费者
    for i := 0; i < 3; i++ {
        mc.wg.Add(1)
        go mc.worker(i, msgs)
    }
}

func (mc *MessageConsumer) worker(id int, msgs <-chan amqp.Delivery) {
    defer mc.wg.Done()
    
    for {
        select {
        case msg := <-msgs:
            if msg.Body == nil {
                continue
            }
            
            fmt.Printf("Worker %d 处理消息: %s\n", id, msg.Body)
            
            // 模拟处理时间
            time.Sleep(2 * time.Second)
            
            // 确认消息
            msg.Ack(false)
            
        case <-mc.quit:
            fmt.Printf("Worker %d 停止\n", id)
            return
        }
    }
}

func (mc *MessageConsumer) Shutdown(ctx context.Context) error {
    fmt.Println("停止消息消费者...")
    
    // 发送停止信号
    close(mc.quit)
    
    // 等待所有worker完成
    done := make(chan struct{})
    go func() {
        mc.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        fmt.Println("所有消息处理完成")
    case <-ctx.Done():
        fmt.Println("等待消息处理超时")
    }
    
    // 关闭连接
    mc.channel.Close()
    return mc.conn.Close()
}

func main() {
    consumer, err := NewMessageConsumer()
    if err != nil {
        log.Fatal(err)
    }
    
    consumer.Start()
    
    // 等待退出信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("开始优雅关闭消息消费者...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := consumer.Shutdown(ctx); err != nil {
        log.Printf("关闭消息消费者失败: %v", err)
    }
    
    fmt.Println("消息消费者已优雅关闭")
}
```

## 6. 完整的应用优雅退出框架

### 6.1 优雅退出管理器

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// Shutdowner 定义可关闭的组件接口
type Shutdowner interface {
    Shutdown(ctx context.Context) error
}

// GracefulShutdown 优雅关闭管理器
type GracefulShutdown struct {
    components []Shutdowner
    timeout    time.Duration
    mu         sync.Mutex
}

func NewGracefulShutdown(timeout time.Duration) *GracefulShutdown {
    return &GracefulShutdown{
        components: make([]Shutdowner, 0),
        timeout:    timeout,
    }
}

func (gs *GracefulShutdown) AddComponent(component Shutdowner) {
    gs.mu.Lock()
    defer gs.mu.Unlock()
    gs.components = append(gs.components, component)
}

func (gs *GracefulShutdown) WaitForSignal() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    sig := <-quit
    fmt.Printf("收到信号: %v，开始优雅关闭...\n", sig)
    
    gs.Shutdown()
}

func (gs *GracefulShutdown) Shutdown() {
    ctx, cancel := context.WithTimeout(context.Background(), gs.timeout)
    defer cancel()
    
    gs.mu.Lock()
    components := make([]Shutdowner, len(gs.components))
    copy(components, gs.components)
    gs.mu.Unlock()
    
    // 并发关闭所有组件
    var wg sync.WaitGroup
    for i, component := range components {
        wg.Add(1)
        go func(idx int, comp Shutdowner) {
            defer wg.Done()
            
            if err := comp.Shutdown(ctx); err != nil {
                log.Printf("组件 %d 关闭失败: %v", idx, err)
            } else {
                log.Printf("组件 %d 已成功关闭", idx)
            }
        }(i, component)
    }
    
    // 等待所有组件关闭完成
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        fmt.Println("所有组件已优雅关闭")
    case <-ctx.Done():
        fmt.Println("关闭超时，强制退出")
    }
}

// 示例组件
type HTTPServer struct {
    server *http.Server
}

func (h *HTTPServer) Shutdown(ctx context.Context) error {
    return h.server.Shutdown(ctx)
}

type Database struct {
    db *sql.DB
}

func (d *Database) Shutdown(ctx context.Context) error {
    return d.db.Close()
}

type MessageQueue struct {
    conn *amqp.Connection
}

func (mq *MessageQueue) Shutdown(ctx context.Context) error {
    return mq.conn.Close()
}

func main() {
    // 创建优雅关闭管理器
    shutdown := NewGracefulShutdown(30 * time.Second)
    
    // 添加组件
    shutdown.AddComponent(&HTTPServer{})
    shutdown.AddComponent(&Database{})
    shutdown.AddComponent(&MessageQueue{})
    
    // 启动应用
    fmt.Println("应用启动...")
    
    // 等待退出信号
    shutdown.WaitForSignal()
}
```

## 7. 最佳实践

### 7.1 设计原则

1. **信号处理**：正确处理SIGINT和SIGTERM信号
2. **超时控制**：设置合理的关闭超时时间
3. **资源清理**：确保所有资源都被正确释放
4. **状态管理**：维护应用状态，拒绝新请求
5. **错误处理**：妥善处理关闭过程中的错误

### 7.2 常见陷阱

1. **忘记处理信号**：程序无法响应退出信号
2. **超时设置不当**：过短导致强制退出，过长导致响应慢
3. **资源泄漏**：忘记关闭连接、文件等资源
4. **并发问题**：关闭过程中的数据竞争
5. **状态不一致**：部分组件关闭，部分仍在运行

### 7.3 监控和日志

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Logger struct {
    *log.Logger
}

func NewLogger() *Logger {
    return &Logger{
        Logger: log.New(os.Stdout, "[GRACEFUL] ", log.LstdFlags),
    }
}

func (l *Logger) LogShutdownStart() {
    l.Println("开始优雅关闭流程")
}

func (l *Logger) LogShutdownComponent(name string, err error) {
    if err != nil {
        l.Printf("组件 %s 关闭失败: %v", name, err)
    } else {
        l.Printf("组件 %s 关闭成功", name)
    }
}

func (l *Logger) LogShutdownComplete() {
    l.Println("优雅关闭完成")
}

func (l *Logger) LogShutdownTimeout() {
    l.Println("优雅关闭超时，强制退出")
}

// 使用示例
func main() {
    logger := NewLogger()
    
    // 创建上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    logger.LogShutdownStart()
    
    // 模拟组件关闭
    components := []string{"HTTP服务器", "数据库", "消息队列"}
    for _, component := range components {
        // 模拟关闭过程
        time.Sleep(1 * time.Second)
        
        // 模拟可能的错误
        var err error
        if component == "数据库" {
            err = fmt.Errorf("连接超时")
        }
        
        logger.LogShutdownComponent(component, err)
    }
    
    logger.LogShutdownComplete()
}
```

## 8. 总结

优雅退出是生产环境应用的基本要求，主要涉及：

1. **信号处理**：捕获和处理退出信号
2. **资源管理**：正确释放所有资源
3. **状态控制**：拒绝新请求，完成现有请求
4. **超时控制**：设置合理的关闭超时
5. **错误处理**：妥善处理关闭过程中的错误
6. **监控日志**：记录关闭过程，便于调试

通过合理的优雅退出机制，可以确保应用在重启、部署或故障时能够安全地保存数据、释放资源，避免数据丢失和资源泄漏。
