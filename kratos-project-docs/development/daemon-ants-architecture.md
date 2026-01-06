# Table Consumer Daemon Ants 架构设计

本文档说明 `table_consumer_ants.go` 和 `table_consumer_ants_biz.go` 两个文件的分工和设计原则。

## 概述

`TableConsumerDaemonAnts` 是一个基于 [ants](https://github.com/panjf2000/ants) goroutine pool 的数据库表轮询守护进程实现。它采用**框架与业务分离**的设计模式，将通用框架逻辑和具体业务逻辑分开，提高代码的可复用性和可维护性。

## 文件分工

### `table_consumer_ants.go` - 通用框架层

**职责**：提供通用的守护进程框架，不包含任何业务逻辑。

**包含内容**：

1. **核心结构体**
   - `TableConsumerDaemonAnts` - 守护进程主体结构
   - `PoolStats` - Pool 统计信息结构

2. **通用逻辑实现**
   - `Run()` - 主轮询循环，负责从数据库获取记录并提交到 ants pool
   - `Stop()` - 优雅停止守护进程
   - `waitTasks()` - 等待所有任务完成
   - `GetPoolStats()` - 获取 pool 统计信息（用于监控）

3. **配置选项（Option Pattern）**
   - `WithAntsPoolSize()` - 设置 goroutine pool 大小
   - `WithAntsBatchSize()` - 设置每批查询的记录数
   - `WithAntsPollInterval()` - 设置轮询间隔
   - `WithAntsRecordFetcher()` - 设置自定义记录获取函数

4. **构造函数**
   - `NewTableConsumerDaemonAnts()` - 创建守护进程实例
     - 接受 `handler` 和 `fetcher` 作为函数参数
     - 不依赖具体业务类型

**特点**：
- ✅ **可复用**：可用于不同的业务场景
- ✅ **无业务依赖**：不包含任何业务逻辑
- ✅ **职责单一**：只负责调度和任务管理

### `table_consumer_ants_biz.go` - 业务层实现

**职责**：提供业务层的实现模板和示例代码。

**包含内容**：

1. **业务构造函数**
   - `NewTableConsumerDaemonAntsForBiz()` - 业务层入口函数
     - 接受业务相关的依赖（如 `db`、`logger`、`usecase`、`repo` 等）
     - 内部实现具体的 `handler` 和 `fetcher`
     - 调用框架层的 `NewTableConsumerDaemonAnts()` 创建实例

2. **Handler 模板**
   - 处理单条记录的业务逻辑
   - 包含 TODO 注释，指导如何实现具体业务
   - 示例场景：
     - 调用 usecase 处理业务逻辑
     - 调用外部服务（HTTP/gRPC）
     - 更新数据库状态
     - 发送消息到消息队列

3. **Fetcher 模板**
   - 从数据库获取待处理记录的查询逻辑
   - 提供两种方式：
     - **方式 1**：使用 `BuildDefaultFetcher`（推荐，适用于简单查询）
     - **方式 2**：自定义 fetcher（适用于复杂查询）
   - 包含 TODO 注释，指导如何实现具体查询

4. **配置示例**
   - Pool 大小、批次大小、轮询间隔等建议值
   - 根据实际业务需求调整的指导

**特点**：
- ✅ **业务导向**：包含具体的业务逻辑模板
- ✅ **易于定制**：提供清晰的 TODO 指导
- ✅ **示例丰富**：包含多种使用场景的示例代码

## 架构关系图

```
┌─────────────────────────────────────────────────────────────┐
│              table_consumer_ants_biz.go                      │
│                    (业务层实现)                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  NewTableConsumerDaemonAntsForBiz()                         │
│    ├─ 定义 handler (业务处理逻辑)                            │
│    ├─ 定义 fetcher (数据获取逻辑)                            │
│    └─ 调用 NewTableConsumerDaemonAnts()                     │
│         └─ 返回 TableConsumerDaemonAnts 实例                │
│                                                              │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              table_consumer_ants.go                         │
│                    (通用框架层)                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  NewTableConsumerDaemonAnts()                               │
│    ├─ 创建 ants.Pool                                        │
│    ├─ 初始化配置选项                                         │
│    └─ 返回 TableConsumerDaemonAnts 实例                     │
│                                                              │
│  TableConsumerDaemonAnts.Run()                              │
│    ├─ 启动轮询循环                                           │
│    ├─ 调用 fetcher 获取记录                                  │
│    ├─ 提交任务到 ants pool                                   │
│    └─ 在 pool 中执行 handler                                 │
│                                                              │
│  TableConsumerDaemonAnts.Stop()                             │
│    ├─ 停止轮询循环                                           │
│    └─ 等待所有任务完成                                       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## 执行流程详解

### Task 与 Handler 的关系

**重要**：提交给 ants 协程池的是 **task（闭包函数）**，而不是 handler 本身。

执行流程：

```
┌─────────────────────────────────────────────────────────────┐
│  Run() 主循环（单 goroutine）                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 调用 fetcher 获取一批记录                               │
│     records := fetcher(ctx, batchSize)                     │
│                                                             │
│  2. 遍历每条记录，创建 task（闭包）                          │
│     for _, record := range records {                       │
│         recordCopy := record  // 避免闭包问题               │
│         task := func() {                                    │
│             defer wg.Done()                                 │
│             handler(ctx, recordCopy)  // ← 调用 handler     │
│         }                                                   │
│                                                             │
│  3. 提交 task 到 ants pool                                  │
│     pool.Submit(task)  // ← 提交的是 task，不是 handler    │
│     }                                                       │
│                                                             │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼ 多个 task 并发提交
┌─────────────────────────────────────────────────────────────┐
│  Ants Pool（多个 goroutine，数量 = poolSize）               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Goroutine 1: 执行 task1 → 调用 handler(record1)          │
│  Goroutine 2: 执行 task2 → 调用 handler(record2)           │
│  Goroutine 3: 执行 task3 → 调用 handler(record3)           │
│  ...                                                        │
│  Goroutine N: 执行 taskN → 调用 handler(recordN)           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**代码实现**（来自 `table_consumer_ants.go`）：

```165:178:internal/biz/daemon/table_consumer_ants.go
// 创建任务闭包
recordCopy := record // 避免闭包问题
task := func() {
    defer d.wg.Done()
    
    // 处理记录
    if err := d.handler(ctx, recordCopy); err != nil {
        logHelper.Errorf("failed to process record: %v", err)
        // 根据业务需求决定是否重试或更新状态为失败
        return
    }

    logHelper.Debugf("processed record successfully")
}

// 提交任务到 pool
d.wg.Add(1)
if err := d.pool.Submit(task); err != nil {
```

**总结**：
- ✅ 提交给 ants pool 的是 **task（闭包函数）**
- ✅ 每个 task 捕获一条记录（`recordCopy`）
- ✅ task 内部调用 **handler** 处理记录
- ✅ 多个 goroutine 并发执行不同的 task，实现并发处理

### 1. 初始化阶段

```go
// 业务层：定义 handler 和 fetcher
daemon, err := NewTableConsumerDaemonAntsForBiz(db, logger)
// ↓
// 框架层：创建守护进程实例
daemon := NewTableConsumerDaemonAnts(db, handler, logger, opts...)
```

### 2. 运行阶段

```go
// Worker Manager 启动守护进程
daemon.Run(ctx)
// ↓
// 框架层执行：
// 1. 启动轮询循环（ticker）
// 2. 每轮询间隔调用 fetcher 获取记录
// 3. 对每条记录创建一个 task（闭包函数）
// 4. 将 task 提交到 ants pool（pool.Submit(task)）
// 5. Pool 中的 goroutine 执行 task，task 内部调用 handler
```

**关键点**：
- 提交给 ants pool 的是 **task（闭包函数）**，不是 handler 本身
- 每个 task 是一个闭包，捕获了单条记录（`recordCopy`）
- task 内部调用 `handler(ctx, recordCopy)` 处理记录
- 多个 goroutine 并发执行不同的 task，每个 task 处理不同的记录

### 3. 停止阶段

```go
// Worker Manager 停止守护进程
daemon.Stop()
// ↓
// 框架层执行：
// 1. 停止轮询循环
// 2. 关闭 ants pool（不再接受新任务）
// 3. 等待所有进行中的任务完成
```

## 设计优势

### 1. 关注点分离（Separation of Concerns）

- **框架层**：专注于调度、并发控制、生命周期管理
- **业务层**：专注于业务逻辑实现

### 2. 可复用性（Reusability）

- `table_consumer_ants.go` 可以用于不同的业务场景
- 只需实现不同的 `handler` 和 `fetcher` 即可

### 3. 易于维护（Maintainability）

- 业务变更只需修改 `table_consumer_ants_biz.go`
- 框架优化只需修改 `table_consumer_ants.go`
- 两者互不影响

### 4. 清晰的职责划分（Clear Responsibilities）

- **框架层**：负责"如何调度"
- **业务层**：负责"如何处理"

## 使用示例

### 基本使用

```go
// 1. 在业务层实现 handler 和 fetcher
daemon, err := NewTableConsumerDaemonAntsForBiz(
    db,
    logger,
    // 可以传入 usecase 或 repo
    // userUsecase,
    // orderRepo,
)

// 2. 注册到 Worker Manager
workerManager.RegisterDaemon(daemon)

// 3. 启动 Worker Manager
workerManager.Start(ctx)
```

### 自定义配置

在 `NewTableConsumerDaemonAntsForBiz` 中修改配置选项：

```go
daemon, err := NewTableConsumerDaemonAnts(
    db,
    handler,
    logger,
    WithAntsPoolSize(20),                // 自定义 pool 大小
    WithAntsBatchSize(200),               // 自定义批次大小
    WithAntsPollInterval(2*time.Second),  // 自定义轮询间隔
    WithAntsRecordFetcher(customFetcher),  // 使用自定义 fetcher
)
```

## 代码位置

- **框架层**：`internal/biz/daemon/table_consumer_ants.go`
- **业务层**：`internal/biz/daemon/table_consumer_ants_biz.go`

## 相关文档

- [Worker 集成最佳实践](./worker-integration.md) - 如何在主应用中集成 daemon-worker
- [Job 实现指南](./job-implementation.md) - Daemon Job 的实现模式和最佳实践
- [分层架构设计](../architecture/layered-architecture.md) - Clean Architecture 分层原则

## 总结

这种**框架与业务分离**的设计模式符合 Clean Architecture 的分层原则：

1. **框架层**（`table_consumer_ants.go`）提供通用的、可复用的基础设施
2. **业务层**（`table_consumer_ants_biz.go`）提供具体的业务实现模板

通过这种设计，我们实现了：
- ✅ 代码复用
- ✅ 职责清晰
- ✅ 易于维护
- ✅ 易于测试

