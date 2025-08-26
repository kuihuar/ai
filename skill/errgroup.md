# Go errgroup 错误组详解

## 概述

`errgroup` 是 Go 语言中用于管理一组 goroutine 的同步原语，它提供了以下功能：

1. **等待所有 goroutine 完成** - 类似 sync.WaitGroup
2. **收集第一个发生的错误** - 当任何 goroutine 返回错误时，立即返回
3. **当发生错误时取消所有 goroutine** - 通过 context 实现取消机制
4. **支持 context 取消** - 提供更好的控制机制

errgroup 是 `sync.WaitGroup` 的增强版本，专门用于处理可能出错的并发任务。

## 基本用法

### 1. 基本 errgroup 使用

```go
import "golang.org/x/sync/errgroup"

func basicExample() {
    var g errgroup.Group

    // 启动多个 goroutine
    g.Go(func() error {
        time.Sleep(time.Millisecond * 100)
        fmt.Println("Task 1 completed")
        return nil
    })

    g.Go(func() error {
        time.Sleep(time.Millisecond * 200)
        fmt.Println("Task 2 completed")
        return nil
    })

    g.Go(func() error {
        time.Sleep(time.Millisecond * 150)
        fmt.Println("Task 3 completed")
        return nil
    })

    // 等待所有 goroutine 完成并检查错误
    if err := g.Wait(); err != nil {
        fmt.Printf("Error occurred: %v\n", err)
    } else {
        fmt.Println("All tasks completed successfully")
    }
}
```

### 2. errgroup 错误处理

当任何一个 goroutine 返回错误时，`Wait()` 会立即返回第一个错误：

```go
func errorHandlingExample() {
    var g errgroup.Group

    g.Go(func() error {
        time.Sleep(time.Millisecond * 100)
        fmt.Println("Task 1 completed")
        return nil
    })

    g.Go(func() error {
        time.Sleep(time.Millisecond * 50)
        fmt.Println("Task 2 failed")
        return errors.New("task 2 failed")
    })

    g.Go(func() error {
        time.Sleep(time.Millisecond * 200)
        fmt.Println("Task 3 completed")
        return nil
    })

    // 当任何一个 goroutine 返回错误时，Wait 会立即返回第一个错误
    if err := g.Wait(); err != nil {
        fmt.Printf("First error: %v\n", err)
    }
}
```

### 3. errgroup 与 context 结合

```go
func withContextExample() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
    defer cancel()

    g, ctx := errgroup.WithContext(ctx)

    g.Go(func() error {
        select {
        case <-time.After(time.Millisecond * 500):
            fmt.Println("Task 1 completed")
            return nil
        case <-ctx.Done():
            fmt.Println("Task 1 cancelled")
            return ctx.Err()
        }
    })

    g.Go(func() error {
        select {
        case <-time.After(time.Millisecond * 100):
            fmt.Println("Task 2 completed")
            return nil
        case <-ctx.Done():
            fmt.Println("Task 2 cancelled")
            return ctx.Err()
        }
    })

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 高级用法

### 1. 限制并发数量

#### 使用信号量限制并发

```go
func withLimitExample() {
    g := new(errgroup.Group)
    sem := make(chan struct{}, 2) // 限制最多同时运行 2 个 goroutine

    tasks := []string{"task1", "task2", "task3", "task4", "task5"}

    for _, task := range tasks {
        task := task // 创建副本避免闭包问题
        g.Go(func() error {
            sem <- struct{}{} // 获取信号量
            defer func() {
                <-sem // 释放信号量
            }()

            time.Sleep(time.Millisecond * 100)
            fmt.Printf("%s completed\n", task)
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

#### 使用 SetLimit 方法（Go 1.19+）

```go
func setLimitExample() {
    g := new(errgroup.Group)
    g.SetLimit(2) // 限制最多同时运行 2 个 goroutine

    tasks := []string{"task1", "task2", "task3", "task4", "task5"}

    for _, task := range tasks {
        task := task
        g.Go(func() error {
            time.Sleep(time.Millisecond * 100)
            fmt.Printf("%s completed\n", task)
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### 2. 错误传播和取消

```go
func errorPropagationExample() {
    g, ctx := errgroup.WithContext(context.Background())

    g.Go(func() error {
        select {
        case <-time.After(time.Millisecond * 100):
            fmt.Println("Task 1 completed")
            return nil
        case <-ctx.Done():
            fmt.Println("Task 1 cancelled due to error in other task")
            return ctx.Err()
        }
    })

    g.Go(func() error {
        time.Sleep(time.Millisecond * 50)
        fmt.Println("Task 2 failed")
        return errors.New("task 2 failed")
    })

    g.Go(func() error {
        select {
        case <-time.After(time.Millisecond * 200):
            fmt.Println("Task 3 completed")
            return nil
        case <-ctx.Done():
            fmt.Println("Task 3 cancelled due to error in other task")
            return ctx.Err()
        }
    })

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 实际应用场景

### 1. 并行文件处理

```go
type FileProcessor struct {
    concurrency int
}

func NewFileProcessor(concurrency int) *FileProcessor {
    return &FileProcessor{concurrency: concurrency}
}

func (fp *FileProcessor) ProcessFiles(files []string) error {
    g := new(errgroup.Group)
    g.SetLimit(fp.concurrency)

    for _, file := range files {
        file := file
        g.Go(func() error {
            return fp.processFile(file)
        })
    }

    return g.Wait()
}

func (fp *FileProcessor) processFile(filename string) error {
    // 模拟文件处理
    time.Sleep(time.Millisecond * 100)
    
    // 模拟某些文件处理失败
    if filename == "error.txt" {
        return fmt.Errorf("failed to process file: %s", filename)
    }
    
    fmt.Printf("Processed file: %s\n", filename)
    return nil
}
```

### 2. 并行 HTTP 请求

```go
type HTTPClient struct {
    timeout time.Duration
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
    return &HTTPClient{timeout: timeout}
}

func (hc *HTTPClient) FetchURLs(urls []string) ([]string, error) {
    g, ctx := errgroup.WithContext(context.Background())
    g.SetLimit(5) // 限制并发请求数

    results := make([]string, len(urls))
    
    for i, url := range urls {
        i, url := i, url // 创建副本
        g.Go(func() error {
            result, err := hc.fetchURL(ctx, url)
            if err != nil {
                return err
            }
            results[i] = result
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }

    return results, nil
}

func (hc *HTTPClient) fetchURL(ctx context.Context, url string) (string, error) {
    // 模拟 HTTP 请求
    select {
    case <-time.After(hc.timeout):
        return fmt.Sprintf("Response from %s", url), nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

### 3. 数据库批量操作

```go
type DatabaseProcessor struct {
    batchSize int
}

func NewDatabaseProcessor(batchSize int) *DatabaseProcessor {
    return &DatabaseProcessor{batchSize: batchSize}
}

func (dp *DatabaseProcessor) ProcessBatch(records []Record) error {
    g := new(errgroup.Group)
    g.SetLimit(3) // 限制并发数据库连接数

    // 分批处理
    for i := 0; i < len(records); i += dp.batchSize {
        end := i + dp.batchSize
        if end > len(records) {
            end = len(records)
        }
        
        batch := records[i:end]
        g.Go(func() error {
            return dp.processBatch(batch)
        })
    }

    return g.Wait()
}

type Record struct {
    ID   int
    Data string
}

func (dp *DatabaseProcessor) processBatch(records []Record) error {
    // 模拟数据库批量操作
    time.Sleep(time.Millisecond * 50)
    
    // 模拟某些批次处理失败
    if len(records) > 0 && records[0].ID == 999 {
        return fmt.Errorf("failed to process batch with ID: %d", records[0].ID)
    }
    
    fmt.Printf("Processed batch with %d records\n", len(records))
    return nil
}
```

## 与其他同步原语结合

### 1. errgroup 与 sync.Mutex

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (sc *SafeCounter) Increment() {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.count++
}

func (sc *SafeCounter) GetCount() int {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    return sc.count
}

func withMutexExample() {
    counter := &SafeCounter{}
    var g errgroup.Group

    // 启动多个 goroutine 并发增加计数器
    for i := 0; i < 100; i++ {
        g.Go(func() error {
            counter.Increment()
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }

    fmt.Printf("Final count: %d\n", counter.GetCount())
}
```

### 2. errgroup 与 sync.Map

```go
func withSyncMapExample() {
    var m sync.Map
    var g errgroup.Group

    // 并发写入 sync.Map
    for i := 0; i < 10; i++ {
        i := i
        g.Go(func() error {
            m.Store(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }

    // 读取所有值
    for i := 0; i < 10; i++ {
        if value, ok := m.Load(fmt.Sprintf("key%d", i)); ok {
            fmt.Printf("key%d: %s\n", i, value)
        }
    }
}
```

## 错误处理策略

### 1. 错误分类处理

```go
type TaskError struct {
    TaskID string
    Err    error
}

func (te TaskError) Error() string {
    return fmt.Sprintf("task %s failed: %v", te.TaskID, te.Err)
}

func errorClassificationExample() {
    g := new(errgroup.Group)
    errors := make(chan TaskError, 10)

    tasks := []string{"task1", "task2", "task3", "task4"}

    for _, taskID := range tasks {
        taskID := taskID
        g.Go(func() error {
            if err := processTask(taskID); err != nil {
                taskErr := TaskError{TaskID: taskID, Err: err}
                select {
                case errors <- taskErr:
                default:
                    // 错误通道已满，记录日志
                    log.Printf("Error channel full, dropping error: %v", taskErr)
                }
                return taskErr
            }
            return nil
        })
    }

    // 等待所有任务完成
    if err := g.Wait(); err != nil {
        fmt.Printf("Group error: %v\n", err)
    }

    // 处理收集到的错误
    close(errors)
    for taskErr := range errors {
        fmt.Printf("Task error: %v\n", taskErr)
    }
}

func processTask(taskID string) error {
    time.Sleep(time.Millisecond * 100)
    
    // 模拟某些任务失败
    if taskID == "task2" {
        return errors.New("task2 failed")
    }
    
    fmt.Printf("Task %s completed\n", taskID)
    return nil
}
```

### 2. 错误重试机制

```go
func withRetryExample() {
    g := new(errgroup.Group)
    maxRetries := 3

    tasks := []string{"task1", "task2", "task3"}

    for _, taskID := range tasks {
        taskID := taskID
        g.Go(func() error {
            return retryTask(taskID, maxRetries)
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error after retries: %v\n", err)
    }
}

func retryTask(taskID string, maxRetries int) error {
    var lastErr error
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        if err := processTask(taskID); err != nil {
            lastErr = err
            fmt.Printf("Task %s failed (attempt %d/%d): %v\n", taskID, attempt, maxRetries, err)
            
            if attempt < maxRetries {
                time.Sleep(time.Millisecond * time.Duration(attempt*100))
                continue
            }
        } else {
            return nil
        }
    }
    
    return fmt.Errorf("task %s failed after %d attempts: %v", taskID, maxRetries, lastErr)
}
```

## 性能优化

### 1. 批量任务处理

```go
func batchProcessingExample() {
    g := new(errgroup.Group)
    g.SetLimit(5) // 限制并发数

    tasks := make([]string, 100)
    for i := range tasks {
        tasks[i] = fmt.Sprintf("task%d", i)
    }

    // 分批处理任务
    batchSize := 20
    for i := 0; i < len(tasks); i += batchSize {
        end := i + batchSize
        if end > len(tasks) {
            end = len(tasks)
        }
        
        batch := tasks[i:end]
        g.Go(func() error {
            return processBatch(batch)
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}

func processBatch(tasks []string) error {
    for _, task := range tasks {
        if err := processTask(task); err != nil {
            return err
        }
    }
    return nil
}
```

### 2. 内存优化

```go
func memoryOptimizationExample() {
    g := new(errgroup.Group)
    g.SetLimit(10)

    // 使用对象池减少内存分配
    var pool sync.Pool
    pool.New = func() interface{} {
        return make([]byte, 1024)
    }

    for i := 0; i < 100; i++ {
        g.Go(func() error {
            // 从池中获取对象
            buf := pool.Get().([]byte)
            defer pool.Put(buf) // 归还到池中

            // 使用缓冲区处理任务
            time.Sleep(time.Millisecond * 10)
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 最佳实践

### 1. 资源清理

```go
func resourceCleanupExample() {
    g, ctx := errgroup.WithContext(context.Background())
    
    // 模拟资源
    resources := make([]string, 5)
    for i := range resources {
        resources[i] = fmt.Sprintf("resource%d", i)
    }

    // 启动任务
    for i, resource := range resources {
        i, resource := i, resource
        g.Go(func() error {
            defer func() {
                // 确保资源被清理
                fmt.Printf("Cleaning up resource: %s\n", resource)
            }()

            select {
            case <-time.After(time.Millisecond * 100):
                fmt.Printf("Processing resource: %s\n", resource)
                return nil
            case <-ctx.Done():
                return ctx.Err()
            }
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### 2. 优雅关闭

```go
func gracefulShutdownExample() {
    g, ctx := errgroup.WithContext(context.Background())
    shutdown := make(chan struct{})

    // 启动工作 goroutine
    g.Go(func() error {
        for {
            select {
            case <-shutdown:
                fmt.Println("Worker shutting down")
                return nil
            case <-ctx.Done():
                fmt.Println("Worker cancelled")
                return ctx.Err()
            default:
                // 执行工作
                time.Sleep(time.Millisecond * 100)
            }
        }
    })

    // 模拟关闭信号
    go func() {
        time.Sleep(time.Second)
        close(shutdown)
    }()

    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### 3. 错误传播策略

```go
func errorPropagationStrategyExample() {
    g := new(errgroup.Group)
    g.SetLimit(3)

    // 定义错误处理策略
    errorHandler := func(taskID string, err error) {
        fmt.Printf("Handling error for task %s: %v\n", taskID, err)
        // 可以记录日志、发送告警等
    }

    tasks := []string{"task1", "task2", "task3", "task4"}

    for _, taskID := range tasks {
        taskID := taskID
        g.Go(func() error {
            if err := processTask(taskID); err != nil {
                errorHandler(taskID, err)
                return err // 继续传播错误
            }
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        fmt.Printf("Group completed with error: %v\n", err)
    } else {
        fmt.Println("All tasks completed successfully")
    }
}
```

## 与 sync.WaitGroup 的对比

| 特性 | sync.WaitGroup | errgroup |
|------|----------------|----------|
| 等待所有 goroutine 完成 | ✅ | ✅ |
| 错误处理 | ❌ | ✅ |
| Context 取消 | ❌ | ✅ |
| 并发数量限制 | ❌ | ✅ |
| 错误传播 | ❌ | ✅ |
| 资源管理 | 手动 | 自动 |

## 总结

errgroup 是 Go 并发编程中非常重要的工具，它提供了比 sync.WaitGroup 更强大的功能：

### 主要优势

1. **错误处理** - 自动收集和传播错误
2. **Context 支持** - 提供取消和超时机制
3. **并发控制** - 支持限制并发数量
4. **资源管理** - 自动处理资源清理
5. **错误传播** - 当发生错误时自动取消其他任务

### 适用场景

1. **并行文件处理** - 处理大量文件时
2. **并行 HTTP 请求** - 批量请求 API
3. **数据库批量操作** - 并发处理数据库记录
4. **任务队列处理** - 处理任务队列
5. **资源密集型操作** - 需要限制并发数的场景

### 最佳实践

1. **总是使用 context** - 提供取消和超时机制
2. **设置合理的并发限制** - 避免资源耗尽
3. **正确处理错误** - 实现错误分类和重试机制
4. **资源清理** - 确保资源被正确释放
5. **优雅关闭** - 实现优雅的关闭机制

errgroup 是编写高质量 Go 并发程序的必备工具，特别适合需要错误处理的并发场景。
