# JOB 任务调度实现方案

## 目录

- [概述](#概述)
- [JOB 分类](#job-分类)
- [工程实现方式](#工程实现方式)
- [开源方案对比](#开源方案对比)
- [项目实现方案](#项目实现方案)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

## 概述

JOB（任务调度）是软件系统中用于执行定时任务、异步任务、批量任务等的重要组件。本文档介绍 JOB 的工程实现方式和开源方案，帮助开发者选择合适的实现方案。

## JOB 分类

### 1. 按执行方式分类

#### 定时任务（Scheduled Jobs）
- **特点**：按固定时间间隔或特定时间点执行
- **示例**：每天凌晨数据同步、每小时统计报表、每周数据清理
- **实现**：Cron 表达式、定时器

#### 异步任务（Async Jobs）
- **特点**：立即触发，但异步执行
- **示例**：发送邮件、生成报表、图片处理
- **实现**：消息队列、任务队列

#### 批量任务（Batch Jobs）
- **特点**：处理大量数据，通常需要分批处理
- **示例**：数据迁移、批量导入、批量计算
- **实现**：任务队列 + 分批处理

#### 延迟任务（Delayed Jobs）
- **特点**：延迟一定时间后执行
- **示例**：订单超时取消、优惠券过期提醒
- **实现**：延迟队列、定时器

#### 守护进程任务（Daemon Jobs）
- **特点**：长期运行的后台任务，持续监听和处理
- **示例**：Kafka 消费者、WebSocket 连接处理、事件监听器
- **实现**：goroutine + 循环、消息队列消费者、长连接处理
- **关键点**：需要优雅关闭、错误恢复、资源管理

### 2. 按执行环境分类

#### 单机任务
- **特点**：在单个进程/容器中执行
- **适用场景**：轻量级任务、开发环境
- **限制**：无法水平扩展、单点故障

#### 分布式任务
- **特点**：可在多个节点执行，支持负载均衡
- **适用场景**：生产环境、高可用要求
- **优势**：可扩展、高可用、负载均衡

## 工程实现方式

### 方案一：基于 Cron 库（轻量级）

#### 适用场景
- 简单的定时任务
- 单机或少量实例部署
- 任务执行时间短（秒级到分钟级）
- 不需要复杂的任务管理功能

#### 实现方式

**1. 单机 Cron**

```go
// 使用 github.com/robfig/cron/v3
package main

import (
    "github.com/robfig/cron/v3"
)

func main() {
    c := cron.New(cron.WithSeconds())
    
    // 每天凌晨 2 点执行
    c.AddFunc("0 0 2 * * *", func() {
        // 执行任务
    })
    
    c.Start()
    defer c.Stop()
    
    // 保持运行
    select {}
}
```

**优点**：
- 实现简单，代码量少
- 资源占用低
- 无需额外依赖

**缺点**：
- 不支持分布式
- 任务失败无持久化
- 无法动态管理任务
- 多实例会重复执行

**2. 分布式 Cron（加锁）**

```go
// 使用 Redis 分布式锁
func (j *Job) Run(ctx context.Context) error {
    lockKey := fmt.Sprintf("cron:lock:%s", j.Name())
    
    // 尝试获取锁，超时时间 5 分钟
    lock, err := redisClient.SetNX(ctx, lockKey, "locked", 5*time.Minute).Result()
    if err != nil || !lock {
        return fmt.Errorf("failed to acquire lock")
    }
    defer redisClient.Del(ctx, lockKey)
    
    // 执行任务逻辑
    return j.execute(ctx)
}
```

**优点**：
- 支持多实例部署
- 确保任务只执行一次
- 实现相对简单

**缺点**：
- 需要额外的 Redis 依赖
- 锁超时时间需要估算
- 任务执行时间过长可能导致锁过期

### 方案二：基于消息队列（中量级）

#### 适用场景
- 异步任务处理
- 需要任务持久化
- 需要任务重试机制
- 需要任务优先级

#### 实现方式

**1. 基于 Kafka**

```go
// 生产者：定时发送任务消息
func scheduleJob() {
    c := cron.New()
    c.AddFunc("0 0 2 * * *", func() {
        message := &JobMessage{
            JobName: "sync-user",
            Payload: map[string]interface{}{},
        }
        kafkaProducer.Send("job-queue", message)
    })
    c.Start()
}

// 消费者：处理任务
func consumeJob() {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "job-queue",
        GroupID: "job-workers",
    })
    
    for {
        msg, err := reader.ReadMessage(context.Background())
        if err != nil {
            continue
        }
        
        var jobMsg JobMessage
        json.Unmarshal(msg.Value, &jobMsg)
        
        // 处理任务
        processJob(jobMsg)
    }
}
```

**优点**：
- 任务持久化
- 支持多消费者负载均衡
- 高吞吐量
- 支持消息重试

**缺点**：
- 需要维护 Kafka 集群
- 配置相对复杂
- 不适合精确的定时任务

**2. 基于 Redis Stream**

```go
// 生产者
func scheduleJob() {
    c := cron.New()
    c.AddFunc("0 0 2 * * *", func() {
        redisClient.XAdd(context.Background(), &redis.XAddArgs{
            Stream: "job-queue",
            Values: map[string]interface{}{
                "job": "sync-user",
                "data": "{}",
            },
        })
    })
    c.Start()
}

// 消费者
func consumeJob() {
    for {
        streams, err := redisClient.XRead(context.Background(), &redis.XReadArgs{
            Streams: []string{"job-queue", "0"},
            Count:   10,
            Block:   time.Second,
        }).Result()
        
        for _, stream := range streams {
            for _, msg := range stream.Messages {
                // 处理任务
                processJob(msg.Values)
                // 确认消息
                redisClient.XAck(context.Background(), "job-queue", "job-group", msg.ID)
            }
        }
    }
}
```

**优点**：
- 实现简单
- 支持消息确认和重试
- 支持消费者组
- 轻量级

**缺点**：
- 消息持久化依赖 Redis 持久化配置
- 不适合大规模任务

### 方案三：Daemon Job（守护进程任务）

#### 适用场景
- 长期运行的后台任务
- 消息队列消费者（Kafka、RabbitMQ、Redis Stream）
- WebSocket 连接处理
- 事件监听器
- 持续监控和处理任务

#### Go 语言最佳实践

**1. 基础 Daemon Job 模式**

```go
// Daemon 任务接口
type DaemonJob interface {
    Name() string
    Run(ctx context.Context) error
    Stop() error
}

// 基础实现
type BaseDaemon struct {
    name   string
    logger log.Logger
    done   chan struct{}
}

func (d *BaseDaemon) Name() string {
    return d.name
}

func (d *BaseDaemon) Stop() error {
    close(d.done)
    return nil
}

// 使用示例
type KafkaConsumerDaemon struct {
    *BaseDaemon
    reader *kafka.Reader
}

func (d *KafkaConsumerDaemon) Run(ctx context.Context) error {
    logHelper := log.NewHelper(d.logger)
    logHelper.Infof("starting daemon job: %s", d.name)
    
    for {
        select {
        case <-ctx.Done():
            logHelper.Info("context cancelled, stopping daemon")
            return ctx.Err()
        case <-d.done:
            logHelper.Info("daemon stopped")
            return nil
        default:
            // 处理消息
            msg, err := d.reader.ReadMessage(ctx)
            if err != nil {
                logHelper.Errorf("failed to read message: %v", err)
                time.Sleep(time.Second) // 避免快速重试
                continue
            }
            
            // 处理消息
            if err := d.processMessage(ctx, msg); err != nil {
                logHelper.Errorf("failed to process message: %v", err)
                // 根据业务需求决定是否重试或丢弃
            }
        }
    }
}
```

**2. Worker Pool 模式（推荐）**

```go
// 使用 worker pool 处理任务，提高并发性能
type WorkerPoolDaemon struct {
    *BaseDaemon
    workerCount int
    jobQueue     chan Job
    workers      []*Worker
}

func (d *WorkerPoolDaemon) Run(ctx context.Context) error {
    logHelper := log.NewHelper(d.logger)
    
    // 启动 worker pool
    for i := 0; i < d.workerCount; i++ {
        worker := NewWorker(i, d.jobQueue, d.logger)
        d.workers = append(d.workers, worker)
        go worker.Start(ctx)
    }
    
    // 主循环：接收任务并分发
    for {
        select {
        case <-ctx.Done():
            logHelper.Info("stopping worker pool")
            d.stopWorkers()
            return ctx.Err()
        case <-d.done:
            logHelper.Info("daemon stopped")
            d.stopWorkers()
            return nil
        case job := <-d.jobQueue:
            // 任务已分发到 worker pool
            _ = job
        }
    }
}

type Worker struct {
    id      int
    queue   chan Job
    logger  log.Logger
}

func (w *Worker) Start(ctx context.Context) {
    logHelper := log.NewHelper(w.logger)
    logHelper.Infof("worker %d started", w.id)
    
    for {
        select {
        case <-ctx.Done():
            logHelper.Infof("worker %d stopped", w.id)
            return
        case job := <-w.queue:
            if err := job.Process(ctx); err != nil {
                logHelper.Errorf("worker %d failed to process job: %v", w.id, err)
            }
        }
    }
}
```

**3. Kafka Consumer Daemon（项目实际应用）**

```go
// internal/biz/daemon/kafka_consumer.go
package daemon

import (
    "context"
    "time"
    
    "github.com/go-kratos/kratos/v2/log"
    "github.com/segmentio/kafka-go"
)

type KafkaConsumerDaemon struct {
    reader  *kafka.Reader
    logger  log.Logger
    handler MessageHandler
}

type MessageHandler func(ctx context.Context, msg kafka.Message) error

func NewKafkaConsumerDaemon(
    reader *kafka.Reader,
    handler MessageHandler,
    logger log.Logger,
) *KafkaConsumerDaemon {
    return &KafkaConsumerDaemon{
        reader:  reader,
        logger:  logger,
        handler: handler,
    }
}

func (d *KafkaConsumerDaemon) Name() string {
    return "kafka-consumer"
}

func (d *KafkaConsumerDaemon) Run(ctx context.Context) error {
    logHelper := log.NewHelper(d.logger)
    logHelper.Info("starting kafka consumer daemon")
    
    // 使用 goroutine 池处理消息
    const workerCount = 10
    jobQueue := make(chan kafka.Message, 100)
    
    // 启动 worker pool
    for i := 0; i < workerCount; i++ {
        go d.worker(ctx, i, jobQueue)
    }
    
    // 主循环：读取消息
    for {
        select {
        case <-ctx.Done():
            logHelper.Info("context cancelled, stopping kafka consumer")
            return ctx.Err()
        default:
            // 设置读取超时，避免阻塞
            readCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
            msg, err := d.reader.ReadMessage(readCtx)
            cancel()
            
            if err != nil {
                if err == context.DeadlineExceeded {
                    // 超时是正常的，继续循环
                    continue
                }
                if err == context.Canceled {
                    return nil
                }
                logHelper.Errorf("failed to read message: %v", err)
                time.Sleep(time.Second)
                continue
            }
            
            // 非阻塞发送到队列
            select {
            case jobQueue <- msg:
            case <-ctx.Done():
                return ctx.Err()
            default:
                // 队列满时，记录警告但继续处理
                logHelper.Warn("job queue is full, dropping message")
            }
        }
    }
}

func (d *KafkaConsumerDaemon) worker(ctx context.Context, id int, queue chan kafka.Message) {
    logHelper := log.NewHelper(d.logger)
    logHelper.Infof("worker %d started", id)
    
    for {
        select {
        case <-ctx.Done():
            logHelper.Infof("worker %d stopped", id)
            return
        case msg := <-queue:
            // 处理消息
            if err := d.handler(ctx, msg); err != nil {
                logHelper.Errorf("worker %d failed to process message: %v", id, err)
                // 根据业务需求决定是否重试
                continue
            }
            
            // 处理成功，可以提交 offset（如果使用手动提交）
            // d.reader.CommitMessages(ctx, msg)
        }
    }
}

func (d *KafkaConsumerDaemon) Stop() error {
    logHelper := log.NewHelper(d.logger)
    logHelper.Info("stopping kafka consumer daemon")
    return d.reader.Close()
}
```

**4. 数据库表轮询 Daemon（消费数据表数据）**

适用于需要持续轮询数据库表，处理待处理数据的场景。

**模式选择建议**：
- **数据量小（< 100条/批次）**：使用基础 Daemon Job 模式
- **数据量大（> 100条/批次）**：使用 Worker Pool 模式（推荐）

**实现示例（Worker Pool 模式 - 推荐）**：

```go
// internal/biz/daemon/table_consumer.go
package daemon

import (
    "context"
    "database/sql"
    "fmt"
    "sync"
    "time"
    
    "github.com/go-kratos/kratos/v2/log"
)

// TableConsumerDaemon 数据库表轮询守护进程
type TableConsumerDaemon struct {
    db          *sql.DB
    logger      log.Logger
    handler     RecordHandler
    workerCount int
    batchSize   int
    pollInterval time.Duration
    wg          sync.WaitGroup
}

type RecordHandler func(ctx context.Context, record interface{}) error

// NewTableConsumerDaemon 创建表轮询守护进程
func NewTableConsumerDaemon(
    db *sql.DB,
    handler RecordHandler,
    logger log.Logger,
    opts ...TableConsumerOption,
) *TableConsumerDaemon {
    daemon := &TableConsumerDaemon{
        db:          db,
        logger:      logger,
        handler:     handler,
        workerCount: 10,        // 默认 10 个 worker
        batchSize:   100,       // 默认每批 100 条
        pollInterval: 1 * time.Second, // 默认 1 秒轮询一次
    }
    
    for _, opt := range opts {
        opt(daemon)
    }
    
    return daemon
}

type TableConsumerOption func(*TableConsumerDaemon)

func WithWorkerCount(count int) TableConsumerOption {
    return func(d *TableConsumerDaemon) {
        d.workerCount = count
    }
}

func WithBatchSize(size int) TableConsumerOption {
    return func(d *TableConsumerDaemon) {
        d.batchSize = size
    }
}

func WithPollInterval(interval time.Duration) TableConsumerOption {
    return func(d *TableConsumerDaemon) {
        d.pollInterval = interval
    }
}

func (d *TableConsumerDaemon) Name() string {
    return "table-consumer"
}

func (d *TableConsumerDaemon) Run(ctx context.Context) error {
    logHelper := log.NewHelper(d.logger)
    logHelper.Info("starting table consumer daemon")
    
    // 创建任务队列
    jobQueue := make(chan interface{}, d.batchSize*2)
    
    // 启动 worker pool
    for i := 0; i < d.workerCount; i++ {
        d.wg.Add(1)
        go d.worker(ctx, i, jobQueue)
    }
    
    // 主循环：轮询数据库表
    ticker := time.NewTicker(d.pollInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            logHelper.Info("context cancelled, stopping table consumer")
            d.waitWorkers()
            return ctx.Err()
        case <-ticker.C:
            // 从数据库读取待处理数据
            records, err := d.fetchPendingRecords(ctx)
            if err != nil {
                logHelper.Errorf("failed to fetch records: %v", err)
                continue
            }
            
            if len(records) == 0 {
                // 没有数据，继续等待
                continue
            }
            
            logHelper.Infof("fetched %d records to process", len(records))
            
            // 将记录分发到 worker pool
            for _, record := range records {
                select {
                case jobQueue <- record:
                case <-ctx.Done():
                    d.waitWorkers()
                    return ctx.Err()
                default:
                    // 队列满时，记录警告
                    logHelper.Warn("job queue is full, waiting...")
                    select {
                    case jobQueue <- record:
                    case <-ctx.Done():
                        d.waitWorkers()
                        return ctx.Err()
                    }
                }
            }
        }
    }
}

// fetchPendingRecords 从数据库获取待处理记录
// 示例：查询待处理的订单
func (d *TableConsumerDaemon) fetchPendingRecords(ctx context.Context) ([]interface{}, error) {
    // 设置查询超时
    queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // 示例 SQL：查询待处理的订单
    // 实际使用时需要根据业务表结构调整
    query := `
        SELECT id, order_no, status, created_at 
        FROM orders 
        WHERE status = 'pending' 
        ORDER BY created_at ASC 
        LIMIT ?
        FOR UPDATE SKIP LOCKED
    `
    
    rows, err := d.db.QueryContext(queryCtx, query, d.batchSize)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    var records []interface{}
    for rows.Next() {
        var record struct {
            ID        int64
            OrderNo   string
            Status    string
            CreatedAt time.Time
        }
        
        if err := rows.Scan(&record.ID, &record.OrderNo, &record.Status, &record.CreatedAt); err != nil {
            logHelper := log.NewHelper(d.logger)
            logHelper.Errorf("failed to scan record: %v", err)
            continue
        }
        
        records = append(records, record)
    }
    
    return records, rows.Err()
}

// worker 处理单个记录
func (d *TableConsumerDaemon) worker(ctx context.Context, id int, queue chan interface{}) {
    defer d.wg.Done()
    logHelper := log.NewHelper(d.logger)
    logHelper.Infof("worker %d started", id)
    
    for {
        select {
        case <-ctx.Done():
            logHelper.Infof("worker %d stopped", id)
            return
        case record := <-queue:
            // 处理记录
            if err := d.handler(ctx, record); err != nil {
                logHelper.Errorf("worker %d failed to process record: %v", id, err)
                // 根据业务需求决定是否重试或更新状态为失败
                continue
            }
            
            // 处理成功，可以更新记录状态
            // d.updateRecordStatus(ctx, record)
        }
    }
}

func (d *TableConsumerDaemon) waitWorkers() {
    logHelper := log.NewHelper(d.logger)
    logHelper.Info("waiting for workers to finish...")
    
    done := make(chan struct{})
    go func() {
        d.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        logHelper.Info("all workers finished")
    case <-time.After(30 * time.Second):
        logHelper.Warn("timeout waiting for workers")
    }
}

func (d *TableConsumerDaemon) Stop() error {
    logHelper := log.NewHelper(d.logger)
    logHelper.Info("stopping table consumer daemon")
    d.waitWorkers()
    return nil
}
```

**使用示例**：

```go
// 创建 daemon
daemon := daemon.NewTableConsumerDaemon(
    db,
    func(ctx context.Context, record interface{}) error {
        // 处理订单
        order := record.(struct {
            ID        int64
            OrderNo   string
            Status    string
            CreatedAt time.Time
        })
        
        logHelper.Infof("processing order: %s", order.OrderNo)
        
        // 业务处理逻辑
        // ...
        
        // 更新订单状态
        _, err := db.ExecContext(ctx, 
            "UPDATE orders SET status = 'processed' WHERE id = ?", 
            order.ID,
        )
        return err
    },
    logger,
    daemon.WithWorkerCount(20),        // 20 个 worker
    daemon.WithBatchSize(200),        // 每批 200 条
    daemon.WithPollInterval(2*time.Second), // 2 秒轮询一次
)

// 启动 daemon
go daemon.Run(ctx)
```

**关键要点**：

1. **使用 `FOR UPDATE SKIP LOCKED`**：避免多个实例同时处理同一条记录
2. **批量查询**：使用 `LIMIT` 限制每次查询数量，避免内存溢出
3. **Worker Pool**：并发处理提高吞吐量
4. **轮询间隔**：根据数据产生频率调整，避免数据库压力过大
5. **状态更新**：处理成功后及时更新记录状态，避免重复处理
6. **错误处理**：记录失败记录，支持重试或死信队列

**5. 优雅关闭模式**

```go
// Daemon Manager：管理多个 daemon 任务
type DaemonManager struct {
    daemons []DaemonJob
    logger  log.Logger
}

func (m *DaemonManager) Start(ctx context.Context) error {
    logHelper := log.NewHelper(m.logger)
    
    // 启动所有 daemon
    for _, daemon := range m.daemons {
        go func(d DaemonJob) {
            if err := d.Run(ctx); err != nil {
                logHelper.Errorf("daemon %s error: %v", d.Name(), err)
            }
        }(daemon)
    }
    
    // 等待上下文取消
    <-ctx.Done()
    
    // 优雅关闭
    logHelper.Info("shutting down daemons...")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan struct{})
    go func() {
        for _, daemon := range m.daemons {
            if err := daemon.Stop(); err != nil {
                logHelper.Errorf("failed to stop daemon %s: %v", daemon.Name(), err)
            }
        }
        close(done)
    }()
    
    select {
    case <-done:
        logHelper.Info("all daemons stopped")
    case <-shutdownCtx.Done():
        logHelper.Warn("timeout waiting for daemons to stop")
    }
    
    return nil
}
```

**5. 错误恢复和重试机制**

```go
// 带指数退避的重试机制
func (d *KafkaConsumerDaemon) processWithRetry(
    ctx context.Context,
    msg kafka.Message,
    maxRetries int,
) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := d.handler(ctx, msg)
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // 指数退避：1s, 2s, 4s, 8s...
        backoff := time.Duration(1<<uint(i)) * time.Second
        if backoff > 30*time.Second {
            backoff = 30 * time.Second
        }
        
        logHelper := log.NewHelper(d.logger)
        logHelper.Warnf("retry %d/%d after %v: %v", i+1, maxRetries, backoff, err)
        
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff):
            // 继续重试
        }
    }
    
    return fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

**6. 健康检查和监控**

```go
// Daemon 健康检查
type HealthChecker interface {
    HealthCheck(ctx context.Context) error
}

func (d *KafkaConsumerDaemon) HealthCheck(ctx context.Context) error {
    // 检查 Kafka 连接
    if d.reader == nil {
        return fmt.Errorf("kafka reader is nil")
    }
    
    // 可以添加更多健康检查逻辑
    // 例如：检查最近一次消息处理时间
    return nil
}

// 定期健康检查
func (m *DaemonManager) StartHealthCheck(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            for _, daemon := range m.daemons {
                if checker, ok := daemon.(HealthChecker); ok {
                    if err := checker.HealthCheck(ctx); err != nil {
                        logHelper := log.NewHelper(m.logger)
                        logHelper.Errorf("daemon %s health check failed: %v", daemon.Name(), err)
                    }
                }
            }
        }
    }
}
```

**优点**：
- 支持长期运行
- 可以处理持续的数据流
- 支持并发处理（Worker Pool）
- 资源利用率高
- 支持优雅关闭

**缺点**：
- 需要处理优雅关闭
- 错误恢复机制复杂
- 资源管理要求高
- 监控和调试相对困难

**关键最佳实践**：

1. **使用 Context 控制生命周期**
   ```go
   // ✅ 正确：使用 context 控制
   for {
       select {
       case <-ctx.Done():
           return ctx.Err()
       default:
           // 处理逻辑
       }
   }
   ```

2. **优雅关闭**
   ```go
   // ✅ 正确：等待正在处理的任务完成
   shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

3. **错误恢复**
   ```go
   // ✅ 正确：错误时等待后重试，避免快速循环
   if err != nil {
       time.Sleep(time.Second)
       continue
   }
   ```

4. **资源限制**
   ```go
   // ✅ 正确：限制 worker 数量和队列大小
   const workerCount = 10
   jobQueue := make(chan Job, 100)
   ```

5. **超时控制**
   ```go
   // ✅ 正确：设置读取超时，避免永久阻塞
   readCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
   defer cancel()
   ```

### 方案四：基于任务调度系统（重量级）

#### 适用场景
- 复杂的任务调度需求
- 需要任务依赖管理
- 需要任务监控和可视化
- 需要任务历史记录
- 大规模任务调度

#### 开源方案

**1. Apache Airflow**

- **语言**：Python
- **特点**：工作流编排、DAG 支持、丰富的 UI
- **适用**：数据管道、ETL 任务
- **部署**：需要 Python 环境，相对复杂

**2. Temporal**

- **语言**：Go/Java/Python/TypeScript
- **特点**：工作流引擎、任务编排、状态管理
- **适用**：复杂业务流程、长期运行任务
- **部署**：需要 Temporal Server

**3. Apache DolphinScheduler**

- **语言**：Java
- **特点**：分布式调度、可视化、多租户
- **适用**：大数据任务调度
- **部署**：需要 Java 环境

**4. XXL-JOB**

- **语言**：Java
- **特点**：分布式任务调度、Web UI、任务监控
- **适用**：Java 生态、中小规模任务
- **部署**：需要 Java 环境

**5. Asynq (Go)**

- **语言**：Go
- **特点**：Redis 后端、任务队列、重试机制
- **适用**：Go 项目、异步任务
- **部署**：需要 Redis

## 开源方案对比

### Go 生态 Cron 库对比

| 库名 | 特点 | 适用场景 | 维护状态 |
|------|------|----------|----------|
| `github.com/robfig/cron/v3` | 成熟稳定、功能完整、支持秒级精度 | 标准定时任务 | ⭐⭐⭐⭐⭐ 活跃 |
| `github.com/go-co-op/gocron` | 简单易用、链式调用 | 简单定时任务 | ⭐⭐⭐⭐ 活跃 |
| `github.com/rakanalh/scheduler` | 轻量级、支持时区 | 轻量级任务 | ⭐⭐⭐ 一般 |

### 任务队列库对比

| 库名 | 后端 | 特点 | 适用场景 |
|------|------|------|----------|
| `github.com/hibiken/asynq` | Redis | 功能完整、重试、优先级 | Go 异步任务 |
| `github.com/adjust/rmq` | Redis | 简单易用 | 简单任务队列 |
| `github.com/segmentio/kafka-go` | Kafka | 高吞吐、分布式 | 大规模消息处理 |

### 任务调度系统对比

| 系统 | 语言 | 复杂度 | 功能 | 适用规模 |
|------|------|--------|------|----------|
| Airflow | Python | 高 | ⭐⭐⭐⭐⭐ | 大规模 |
| Temporal | 多语言 | 中 | ⭐⭐⭐⭐⭐ | 大规模 |
| XXL-JOB | Java | 中 | ⭐⭐⭐⭐ | 中大规模 |
| Asynq | Go | 低 | ⭐⭐⭐ | 中小规模 |

## 项目实现方案

### 当前实现

本项目采用 **基于 Cron 库 + 分布式锁** 的方案：

```12:96:internal/biz/cron/manager.go
// Manager Cron 任务管理器
type Manager struct {
	cron   *cron.Cron
	jobs   []Job
	runner *JobRunner
	logger log.Logger
}

// NewManager 创建 Cron 管理器
func NewManager(logger log.Logger, timezone string) (*Manager, error) {
	// 解析时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}

	// 创建 cron 实例，支持秒级精度
	c := cron.New(
		cron.WithSeconds(), // 支持秒级精度（6 位表达式）
		cron.WithLocation(loc), // 设置时区
		cron.WithChain( // 添加恢复链，防止 panic
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &Manager{
		cron:   c,
		jobs:   make([]Job, 0),
		runner: NewJobRunner(logger),
		logger: logger,
	}, nil
}

// RegisterJob 注册任务
func (m *Manager) RegisterJob(job Job) error {
	logHelper := log.NewHelper(m.logger)

	// 注册任务（AddFunc 会自动验证 cron 表达式）
	_, err := m.cron.AddFunc(job.Spec(), m.runner.Run(job))
	if err != nil {
		return fmt.Errorf("failed to register job %s: %w", job.Name(), err)
	}

	m.jobs = append(m.jobs, job)
	logHelper.Infof("registered cron job: %s, spec: %s", job.Name(), job.Spec())

	return nil
}

// Start 启动 Cron 管理器
func (m *Manager) Start(ctx context.Context) error {
	logHelper := log.NewHelper(m.logger)

	if len(m.jobs) == 0 {
		logHelper.Warn("no cron jobs registered")
		return nil
	}

	logHelper.Infof("starting cron manager with %d jobs", len(m.jobs))
	m.cron.Start()

	// 等待上下文取消
	<-ctx.Done()

	logHelper.Info("stopping cron manager...")
	stopCtx := m.cron.Stop()

	// 等待所有运行中的任务完成
	waitCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	select {
	case <-waitCtx.Done():
		logHelper.Warn("timeout waiting for cron jobs to finish")
	case <-stopCtx.Done():
		logHelper.Info("all cron jobs stopped")
	}

	return nil
}

// GetJobs 获取所有注册的任务
func (m *Manager) GetJobs() []Job {
	return m.jobs
}
```

### 架构设计

```
┌─────────────────────────────────────────┐
│         Cron Worker Application         │
│  ┌───────────────────────────────────┐  │
│  │      Cron Manager                 │  │
│  │  ┌─────────────────────────────┐   │  │
│  │  │  robfig/cron/v3            │   │  │
│  │  └─────────────────────────────┘   │  │
│  │  ┌─────────────────────────────┐   │  │
│  │  │  Job Runner                │   │  │
│  │  │  - 错误处理                │   │  │
│  │  │  - 日志记录                │   │  │
│  │  └─────────────────────────────┘   │  │
│  └───────────────────────────────────┘  │
│  ┌───────────────────────────────────┐  │
│  │  Jobs                             │  │
│  │  - SyncUserJob                    │  │
│  │  - CleanupLogJob                  │  │
│  │  - ...                            │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
           │
           │ 执行任务
           ▼
┌─────────────────────────────────────────┐
│      Business Logic Layer               │
│  ┌───────────────────────────────────┐  │
│  │  UserUsecase                      │  │
│  │  LogUsecase                       │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
           │
           │ 数据访问
           ▼
┌─────────────────────────────────────────┐
│      Data Layer                          │
│  - Database                             │
│  - Redis (分布式锁)                      │
│  - External Services                     │
└─────────────────────────────────────────┘
```

### 扩展方案

#### 1. 添加分布式锁支持

```go
// internal/biz/cron/job.go
type Job interface {
    Name() string
    Spec() string
    Run(ctx context.Context) error
    // 新增：是否需要分布式锁
    RequireLock() bool
    // 新增：锁超时时间
    LockTimeout() time.Duration
}

// internal/biz/cron/manager.go
type Manager struct {
    // ... 现有字段
    redisClient *redis.Client // 新增 Redis 客户端
}

// 在 JobRunner 中添加锁逻辑
func (r *JobRunner) Run(job Job) func() {
    return func() {
        ctx := context.Background()
        
        // 如果需要分布式锁
        if job.RequireLock() {
            lockKey := fmt.Sprintf("cron:lock:%s", job.Name())
            timeout := job.LockTimeout()
            if timeout == 0 {
                timeout = 5 * time.Minute
            }
            
            lock, err := r.redisClient.SetNX(ctx, lockKey, "locked", timeout).Result()
            if err != nil || !lock {
                logHelper.Warnf("failed to acquire lock for job %s", job.Name())
                return
            }
            defer r.redisClient.Del(ctx, lockKey)
        }
        
        // 执行任务
        // ...
    }
}
```

#### 2. 添加任务监控

```go
// 集成 Prometheus metrics
var (
    jobDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "cron_job_duration_seconds",
            Help: "Cron job execution duration",
        },
        []string{"job_name"},
    )
    
    jobCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cron_job_total",
            Help: "Total number of cron jobs",
        },
        []string{"job_name", "status"},
    )
)

func (r *JobRunner) Run(job Job) func() {
    return func() {
        startTime := time.Now()
        // ... 执行任务
        
        duration := time.Since(startTime).Seconds()
        jobDuration.WithLabelValues(job.Name()).Observe(duration)
        
        if err != nil {
            jobCounter.WithLabelValues(job.Name(), "failed").Inc()
        } else {
            jobCounter.WithLabelValues(job.Name(), "success").Inc()
        }
    }
}
```

#### 3. 迁移到 Asynq（如果需要）

如果未来需要更强大的任务队列功能，可以考虑迁移到 Asynq：

```go
// 使用 Asynq
import "github.com/hibiken/asynq"

// 创建客户端
client := asynq.NewClient(asynq.RedisClientOpt{
    Addr: "localhost:6379",
})

// 定时任务：使用 Cron 触发，发送到 Asynq
c := cron.New()
c.AddFunc("0 0 2 * * *", func() {
    task := asynq.NewTask("sync-user", []byte("{}"))
    client.Enqueue(task)
})
c.Start()

// Worker：处理任务
srv := asynq.NewServer(
    asynq.RedisClientOpt{Addr: "localhost:6379"},
    asynq.Config{
        Concurrency: 10,
    },
)

mux := asynq.NewServeMux()
mux.HandleFunc("sync-user", func(ctx context.Context, t *asynq.Task) error {
    // 处理任务
    return syncUser(ctx)
})

srv.Run(mux)
```

## 最佳实践

### 1. 任务设计原则

#### 单一职责
每个任务只做一件事，便于测试和维护。

```go
// ❌ 不好：一个任务做多件事
func (j *Job) Run(ctx context.Context) error {
    j.syncUsers(ctx)
    j.cleanupLogs(ctx)
    j.sendReports(ctx)
    return nil
}

// ✅ 好：拆分为多个任务
type SyncUserJob struct {}
type CleanupLogJob struct {}
type SendReportJob struct {}
```

#### 幂等性
任务可以安全地重复执行，不会产生副作用。

```go
// ✅ 幂等性设计
func (j *SyncUserJob) Run(ctx context.Context) error {
    users := j.fetchUsers(ctx)
    
    for _, user := range users {
        // 使用 Upsert 而不是 Insert
        j.repo.UpsertUser(ctx, user)
    }
    
    return nil
}
```

#### 超时控制
长时间运行的任务要设置超时。

```go
func (j *Job) Run(ctx context.Context) error {
    // 设置任务超时时间为 10 分钟
    ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
    defer cancel()
    
    return j.execute(ctx)
}
```

#### 错误处理
妥善处理错误，避免影响其他任务。

```go
func (r *JobRunner) Run(job Job) func() {
    return func() {
        defer func() {
            if r := recover(); r != nil {
                logHelper.Errorf("panic in job %s: %v", job.Name(), r)
            }
        }()
        
        if err := job.Run(ctx); err != nil {
            logHelper.Errorf("job %s failed: %v", job.Name(), err)
            // 可以发送告警
            // j.sendAlert(err)
        }
    }
}
```

### 2. 日志记录

#### 结构化日志
使用结构化日志，便于查询和分析。

```go
logHelper.WithFields(log.Fields{
    "job_name": job.Name(),
    "start_time": startTime,
    "duration": time.Since(startTime),
}).Info("job completed")
```

#### 关键节点记录
在任务的关键节点记录日志。

```go
func (j *SyncUserJob) Run(ctx context.Context) error {
    logHelper.Info("starting user sync")
    
    users, err := j.fetchUsers(ctx)
    if err != nil {
        return err
    }
    logHelper.Infof("fetched %d users", len(users))
    
    synced, err := j.syncUsers(ctx, users)
    if err != nil {
        return err
    }
    logHelper.Infof("synced %d users", synced)
    
    logHelper.Info("user sync completed")
    return nil
}
```

### 3. 监控和告警

#### Metrics 指标
记录任务执行的关键指标。

```go
// 任务执行时间
jobDuration.WithLabelValues(job.Name()).Observe(duration)

// 任务执行次数
jobCounter.WithLabelValues(job.Name(), "success").Inc()

// 任务失败次数
jobCounter.WithLabelValues(job.Name(), "failed").Inc()
```

#### 告警规则
设置合理的告警规则。

```yaml
# Prometheus 告警规则
groups:
  - name: cron_jobs
    rules:
      - alert: CronJobFailed
        expr: rate(cron_job_total{status="failed"}[5m]) > 0.1
        for: 5m
        annotations:
          summary: "Cron job {{ $labels.job_name }} is failing"
      
      - alert: CronJobSlow
        expr: histogram_quantile(0.95, cron_job_duration_seconds) > 300
        for: 5m
        annotations:
          summary: "Cron job {{ $labels.job_name }} is slow"
```

### 4. 分布式部署

#### 分布式锁
多实例部署时使用分布式锁。

```go
// 使用 Redis 分布式锁
func (j *Job) acquireLock(ctx context.Context) (bool, error) {
    lockKey := fmt.Sprintf("cron:lock:%s", j.Name())
    timeout := j.LockTimeout()
    
    // 使用 SET NX EX 原子操作
    result, err := redisClient.SetNX(ctx, lockKey, "locked", timeout).Result()
    return result, err
}
```

#### 锁超时时间
合理设置锁超时时间，避免任务执行时间过长导致锁过期。

```go
// 根据任务类型设置超时时间
func (j *SyncUserJob) LockTimeout() time.Duration {
    // 用户同步任务预计需要 5 分钟
    return 10 * time.Minute // 设置 10 分钟，留有余量
}
```

### 5. 任务配置化

#### 从配置读取 Cron 表达式
支持从配置文件读取 Cron 表达式，便于动态调整。

```go
// config.yaml
cron:
  jobs:
    - name: sync-user
      spec: "0 0 2 * * *"
      enabled: true
    - name: cleanup-log
      spec: "0 0 3 * * *"
      enabled: true

// 代码中读取配置
for _, jobConfig := range config.Cron.Jobs {
    if !jobConfig.Enabled {
        continue
    }
    
    job := jobs.NewJob(jobConfig.Name, jobConfig.Spec)
    manager.RegisterJob(job)
}
```

### 6. Daemon Job 最佳实践

#### Context 控制生命周期
始终使用 `context.Context` 控制 daemon 的生命周期，支持优雅关闭。

```go
// ✅ 正确：使用 context 控制
func (d *Daemon) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 处理逻辑
        }
    }
}

// ❌ 错误：使用全局变量或 channel 控制
var stopFlag bool
func (d *Daemon) Run() error {
    for !stopFlag {
        // 处理逻辑
    }
}
```

#### 优雅关闭
实现优雅关闭，等待正在处理的任务完成。

```go
// ✅ 正确：等待任务完成
func (d *Daemon) Stop() error {
    // 1. 停止接收新任务
    close(d.jobQueue)
    
    // 2. 等待正在处理的任务完成
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan struct{})
    go func() {
        d.wg.Wait() // 等待所有 worker 完成
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return fmt.Errorf("timeout waiting for tasks to complete")
    }
}
```

#### 错误恢复和重试
实现合理的错误恢复机制，避免快速重试导致资源浪费。

```go
// ✅ 正确：指数退避重试
func (d *Daemon) processWithRetry(ctx context.Context, task Task) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := d.process(ctx, task)
        if err == nil {
            return nil
        }
        
        // 指数退避：1s, 2s, 4s
        backoff := time.Duration(1<<uint(i)) * time.Second
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff):
            // 继续重试
        }
    }
    return fmt.Errorf("max retries exceeded")
}

// ❌ 错误：快速重试，可能导致资源浪费
func (d *Daemon) process(ctx context.Context, task Task) error {
    for {
        err := d.doProcess(ctx, task)
        if err == nil {
            return nil
        }
        // 没有等待，立即重试
    }
}
```

#### 资源限制
合理设置 worker 数量和队列大小，避免资源耗尽。

```go
// ✅ 正确：限制资源使用
const (
    maxWorkers = 10
    queueSize  = 100
)

func NewDaemon() *Daemon {
    return &Daemon{
        workerCount: maxWorkers,
        jobQueue:    make(chan Job, queueSize),
    }
}

// ❌ 错误：无限制创建 goroutine
func (d *Daemon) processAll(tasks []Task) {
    for _, task := range tasks {
        go d.process(task) // 可能创建大量 goroutine
    }
}
```

#### 超时控制
为所有阻塞操作设置超时，避免永久阻塞。

```go
// ✅ 正确：设置超时
func (d *Daemon) readMessage(ctx context.Context) (*Message, error) {
    readCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    return d.reader.ReadMessage(readCtx)
}

// ❌ 错误：可能永久阻塞
func (d *Daemon) readMessage(ctx context.Context) (*Message, error) {
    return d.reader.ReadMessage(ctx) // 没有超时
}
```

#### 健康检查
实现健康检查机制，便于监控和故障排查。

```go
// ✅ 正确：实现健康检查
type HealthChecker interface {
    HealthCheck(ctx context.Context) error
}

func (d *Daemon) HealthCheck(ctx context.Context) error {
    // 检查连接状态
    if !d.isConnected() {
        return fmt.Errorf("connection lost")
    }
    
    // 检查最近处理时间
    if time.Since(d.lastProcessTime) > 5*time.Minute {
        return fmt.Errorf("no message processed for 5 minutes")
    }
    
    return nil
}
```

#### 监控和指标
记录关键指标，便于监控和性能分析。

```go
// ✅ 正确：记录指标
var (
    messagesProcessed = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "daemon_messages_processed_total",
            Help: "Total number of messages processed",
        },
        []string{"daemon_name", "status"},
    )
    
    processingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "daemon_processing_duration_seconds",
            Help: "Message processing duration",
        },
        []string{"daemon_name"},
    )
)

func (d *Daemon) processMessage(ctx context.Context, msg Message) error {
    startTime := time.Now()
    defer func() {
        duration := time.Since(startTime).Seconds()
        processingDuration.WithLabelValues(d.name).Observe(duration)
    }()
    
    err := d.handler(ctx, msg)
    if err != nil {
        messagesProcessed.WithLabelValues(d.name, "failed").Inc()
        return err
    }
    
    messagesProcessed.WithLabelValues(d.name, "success").Inc()
    return nil
}
```

### 7. 测试

#### 单元测试
为任务编写单元测试。

```go
func TestSyncUserJob_Run(t *testing.T) {
    // Mock 依赖
    mockRepo := &MockUserRepo{}
    mockLogger := log.NewNopLogger()
    
    job := NewSyncUserJob(mockRepo, mockLogger)
    
    ctx := context.Background()
    err := job.Run(ctx)
    
    assert.NoError(t, err)
    assert.True(t, mockRepo.UpsertCalled)
}
```

#### 集成测试
测试任务在真实环境中的执行。

```go
func TestCronManager_Integration(t *testing.T) {
    // 使用测试数据库
    db := setupTestDB(t)
    defer db.Close()
    
    manager, err := cron.NewManager(logger, "UTC")
    require.NoError(t, err)
    
    job := NewSyncUserJob(userRepo, logger)
    err = manager.RegisterJob(job)
    require.NoError(t, err)
    
    // 启动 manager，等待任务执行
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    go manager.Start(ctx)
    time.Sleep(1 * time.Second)
    cancel()
}
```

## 常见问题

### Q1: 如何确保任务只在一个实例执行？

**A**: 使用分布式锁（Redis 或 etcd）。在任务执行前获取锁，执行完成后释放锁。

```go
func (j *Job) Run(ctx context.Context) error {
    lockKey := fmt.Sprintf("cron:lock:%s", j.Name())
    
    // 尝试获取锁
    lock, err := redisClient.SetNX(ctx, lockKey, "locked", 5*time.Minute).Result()
    if err != nil || !lock {
        return fmt.Errorf("failed to acquire lock")
    }
    defer redisClient.Del(ctx, lockKey)
    
    // 执行任务
    return j.execute(ctx)
}
```

### Q2: 任务执行时间过长怎么办？

**A**: 
1. 在任务内部设置超时控制
2. 将大任务拆分为多个小任务
3. 使用任务队列异步处理

```go
// 设置超时
ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
defer cancel()

// 分批处理
const batchSize = 100
for i := 0; i < len(items); i += batchSize {
    batch := items[i:min(i+batchSize, len(items))]
    if err := processBatch(ctx, batch); err != nil {
        return err
    }
}
```

### Q3: 如何动态添加或删除任务？

**A**: 
- **简单方案**：重启应用（当前实现）
- **高级方案**：使用任务调度系统（Airflow、Temporal）
- **折中方案**：通过 API 动态注册任务（需要扩展 Manager）

```go
// 扩展 Manager 支持动态注册
func (m *Manager) RegisterJobDynamic(job Job) error {
    entryID, err := m.cron.AddFunc(job.Spec(), m.runner.Run(job))
    if err != nil {
        return err
    }
    
    m.jobEntries[job.Name()] = entryID
    m.jobs = append(m.jobs, job)
    return nil
}

func (m *Manager) RemoveJob(jobName string) error {
    entryID, ok := m.jobEntries[jobName]
    if !ok {
        return fmt.Errorf("job not found: %s", jobName)
    }
    
    m.cron.Remove(entryID)
    delete(m.jobEntries, jobName)
    return nil
}
```

### Q4: 任务失败后如何重试？

**A**: 
1. **任务内部重试**：在任务逻辑中实现重试
2. **使用支持重试的库**：如 Asynq
3. **手动重试机制**：记录失败任务，定时重试

```go
// 任务内部重试
func (j *Job) Run(ctx context.Context) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := j.execute(ctx)
        if err == nil {
            return nil
        }
        
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second)
        }
    }
    return fmt.Errorf("max retries exceeded")
}
```

### Q5: 如何监控任务执行情况？

**A**: 
1. **集成 Prometheus**：记录任务执行时间、成功/失败次数
2. **日志聚合**：使用 ELK、Loki 等工具
3. **任务调度系统**：使用 Airflow、XXL-JOB 等自带监控

```go
// Prometheus metrics
var (
    jobDuration = prometheus.NewHistogramVec(...)
    jobCounter = prometheus.NewCounterVec(...)
)

func (r *JobRunner) Run(job Job) func() {
    return func() {
        startTime := time.Now()
        err := job.Run(ctx)
        duration := time.Since(startTime)
        
        jobDuration.WithLabelValues(job.Name()).Observe(duration.Seconds())
        if err != nil {
            jobCounter.WithLabelValues(job.Name(), "failed").Inc()
        } else {
            jobCounter.WithLabelValues(job.Name(), "success").Inc()
        }
    }
}
```

### Q6: 如何处理任务依赖？

**A**: 
1. **简单依赖**：在任务内部调用依赖任务
2. **复杂依赖**：使用任务调度系统（Airflow DAG、Temporal Workflow）

```go
// 简单依赖：在任务内部处理
func (j *JobB) Run(ctx context.Context) error {
    // 先执行依赖任务
    if err := j.depJob.Run(ctx); err != nil {
        return err
    }
    
    // 再执行当前任务
    return j.execute(ctx)
}

// 复杂依赖：使用任务链
type JobChain struct {
    jobs []Job
}

func (c *JobChain) Run(ctx context.Context) error {
    for _, job := range c.jobs {
        if err := job.Run(ctx); err != nil {
            return fmt.Errorf("job %s failed: %w", job.Name(), err)
        }
    }
    return nil
}
```

### Q7: Daemon Job 如何实现优雅关闭？

**A**: 使用 `context.Context` 和 `sync.WaitGroup` 实现优雅关闭。

```go
type Daemon struct {
    workers []*Worker
    wg      sync.WaitGroup
    done    chan struct{}
}

func (d *Daemon) Start(ctx context.Context) error {
    // 启动 workers
    for _, worker := range d.workers {
        d.wg.Add(1)
        go func(w *Worker) {
            defer d.wg.Done()
            w.Run(ctx)
        }(worker)
    }
    
    // 等待上下文取消
    <-ctx.Done()
    
    // 停止接收新任务
    close(d.done)
    
    // 等待所有 workers 完成
    done := make(chan struct{})
    go func() {
        d.wg.Wait()
        close(done)
    }()
    
    // 设置超时
    select {
    case <-done:
        return nil
    case <-time.After(30 * time.Second):
        return fmt.Errorf("timeout waiting for workers")
    }
}
```

### Q8: Daemon Job 如何处理消息积压？

**A**: 使用 Worker Pool 和背压机制。

```go
// 1. 使用 Worker Pool 提高处理能力
const (
    workerCount = 10
    queueSize   = 1000
)

// 2. 监控队列长度，实现背压
func (d *Daemon) processWithBackpressure(ctx context.Context) {
    for {
        // 检查队列是否满
        if len(d.jobQueue) > queueSize*0.8 {
            logHelper.Warn("queue is nearly full, slowing down")
            time.Sleep(100 * time.Millisecond)
            continue
        }
        
        // 读取消息
        msg, err := d.reader.ReadMessage(ctx)
        if err != nil {
            continue
        }
        
        // 非阻塞发送
        select {
        case d.jobQueue <- msg:
        default:
            logHelper.Warn("queue is full, dropping message")
        }
    }
}

// 3. 动态调整 worker 数量
func (d *Daemon) adjustWorkers() {
    queueLen := len(d.jobQueue)
    
    if queueLen > queueSize*0.8 && len(d.workers) < maxWorkers {
        // 增加 worker
        d.addWorker()
    } else if queueLen < queueSize*0.2 && len(d.workers) > minWorkers {
        // 减少 worker
        d.removeWorker()
    }
}
```

### Q9: 消费数据表数据适合哪种 Daemon 模式？

**A**: 根据数据量和处理复杂度选择：

**1. 基础 Daemon Job 模式**（数据量小）
- 每批数据 < 100 条
- 处理逻辑简单
- 单线程处理即可

```go
func (d *SimpleTableDaemon) Run(ctx context.Context) error {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            records, _ := d.fetchPendingRecords(ctx)
            for _, record := range records {
                d.handler(ctx, record)
            }
        }
    }
}
```

**2. Worker Pool 模式**（推荐，数据量大）
- 每批数据 > 100 条
- 需要并发处理
- 提高吞吐量

```go
// 使用 Worker Pool 模式，参考文档中的 TableConsumerDaemon 实现
daemon := daemon.NewTableConsumerDaemon(
    db, handler, logger,
    daemon.WithWorkerCount(20),
    daemon.WithBatchSize(200),
)
```

**选择建议**：
- **小规模**（< 100条/批次）：基础模式
- **中大规模**（> 100条/批次）：Worker Pool 模式
- **超大规模**（> 1000条/批次）：考虑使用消息队列（Kafka）替代轮询

**关键优化点**：
1. 使用 `FOR UPDATE SKIP LOCKED` 避免并发冲突
2. 批量查询，避免单条查询
3. 合理设置轮询间隔，平衡实时性和数据库压力
4. 处理成功后及时更新状态，避免重复处理

### Q10: Daemon Job 如何实现错误恢复？

**A**: 实现指数退避重试和死信队列。

```go
// 1. 指数退避重试
func (d *Daemon) processWithRetry(ctx context.Context, msg Message) error {
    maxRetries := 3
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := d.handler(ctx, msg)
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // 指数退避
        backoff := time.Duration(1<<uint(i)) * time.Second
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff):
            // 继续重试
        }
    }
    
    // 2. 重试失败后，发送到死信队列
    return d.sendToDeadLetterQueue(ctx, msg, lastErr)
}

// 3. 死信队列处理
func (d *Daemon) sendToDeadLetterQueue(ctx context.Context, msg Message, err error) error {
    dlqMsg := DeadLetterMessage{
        Original: msg,
        Error:    err.Error(),
        Time:     time.Now(),
    }
    
    return d.dlqProducer.Send(ctx, dlqMsg)
}
```

### Q11: 如何选择实现方案？

**A**: 根据以下因素选择：

| 因素 | 轻量级 Cron | 消息队列 | 任务调度系统 |
|------|------------|----------|--------------|
| **任务数量** | < 10 | 10-100 | > 100 |
| **任务复杂度** | 简单 | 中等 | 复杂 |
| **执行时间** | 秒级-分钟级 | 分钟级-小时级 | 任意 |
| **依赖关系** | 无 | 简单 | 复杂 |
| **监控需求** | 基础 | 中等 | 完整 |
| **运维成本** | 低 | 中 | 高 |

**推荐**：
- **小规模**：Cron + 分布式锁（当前方案）
- **中规模**：Asynq + Redis
- **大规模**：Airflow 或 Temporal

## 总结

1. **选择合适的方案**：根据任务规模、复杂度、运维成本选择
   - **定时任务**：Cron + 分布式锁
   - **异步任务**：消息队列（Kafka、Redis Stream）
   - **守护进程任务**：Daemon Job + Worker Pool
   - **复杂工作流**：任务调度系统（Airflow、Temporal）

2. **遵循最佳实践**：
   - **任务设计**：单一职责、幂等性、超时控制、错误处理
   - **Daemon Job**：Context 控制、优雅关闭、错误恢复、资源限制
   - **监控告警**：记录关键指标，设置合理告警
   - **分布式部署**：使用分布式锁确保任务只执行一次

3. **Daemon Job 关键要点**：
   - 使用 `context.Context` 控制生命周期
   - 实现优雅关闭，等待任务完成
   - 使用 Worker Pool 提高并发性能
   - 实现错误恢复和重试机制
   - 设置合理的资源限制和超时控制
   - 实现健康检查和监控指标

4. **持续优化**：根据实际需求逐步优化和扩展

## 相关文档

- [Cron 集成指南](./cron-integration.md) - 详细的集成步骤
- [多应用架构](../architecture/multi-app.md) - 了解多应用架构
- [依赖注入](../architecture/dependency-injection.md) - 了解 Wire 依赖注入
- [日志管理](../operations/logging.md) - 了解日志配置和管理
- [可观测性](../operations/observability.md) - 了解监控和追踪

