# Go 并发编程详解

## 📚 目录

- [并发编程概述](#并发编程概述)
- [Goroutine 基础](#goroutine-基础)
- [Channel 通信](#channel-通信)
- [同步原语](#同步原语)
- [并发模式](#并发模式)
- [性能优化](#性能优化)
- [最佳实践](#最佳实践)

## 并发编程概述

Go 的并发模型基于 CSP (Communicating Sequential Processes) 理论，通过 goroutine 和 channel 实现。

### 基本概念

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    fmt.Println("=== Concurrency Overview ===")
    
    // 获取系统信息
    fmt.Printf("CPU cores: %d\n", runtime.NumCPU())
    fmt.Printf("Max procs: %d\n", runtime.GOMAXPROCS(0))
    fmt.Printf("Current goroutines: %d\n", runtime.NumGoroutine())
    
    // 启动 goroutine
    go func() {
        fmt.Println("Hello from goroutine!")
    }()
    
    // 等待 goroutine 完成
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())
}
```

## Goroutine 基础

### 创建和管理

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Goroutine Basics ===")
    
    // 1. 基本创建
    go func() {
        fmt.Println("Goroutine 1")
    }()
    
    // 2. 带参数的 goroutine
    go func(id int) {
        fmt.Printf("Goroutine %d\n", id)
    }(2)
    
    // 3. 使用 WaitGroup 同步
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d\n", id)
            time.Sleep(100 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All goroutines completed")
}
```

### 生命周期管理

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Goroutine Lifecycle ===")
    
    // 使用 context 控制生命周期
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // 启动受控的 goroutine
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine stopped by context")
                return
            case <-ticker.C:
                fmt.Println("Goroutine running...")
            }
        }
    }()
    
    // 等待 context 超时
    <-ctx.Done()
    fmt.Println("Main goroutine completed")
}
```

## Channel 通信

### 基本使用

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Channel Basics ===")
    
    // 1. 无缓冲 channel
    ch := make(chan int)
    
    go func() {
        ch <- 42
    }()
    
    value := <-ch
    fmt.Printf("Received: %d\n", value)
    
    // 2. 有缓冲 channel
    bufferedCh := make(chan int, 3)
    bufferedCh <- 1
    bufferedCh <- 2
    bufferedCh <- 3
    
    fmt.Printf("Buffered channel length: %d\n", len(bufferedCh))
    
    // 3. 关闭 channel
    close(bufferedCh)
    
    for value := range bufferedCh {
        fmt.Printf("Received: %d\n", value)
    }
}
```

### Channel 模式

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Channel Patterns ===")
    
    // 1. 生产者-消费者模式
    producer := make(chan int, 10)
    consumer := make(chan int, 10)
    
    // 生产者
    go func() {
        for i := 0; i < 10; i++ {
            producer <- i
        }
        close(producer)
    }()
    
    // 消费者
    go func() {
        for value := range producer {
            consumer <- value * 2
        }
        close(consumer)
    }()
    
    // 处理结果
    for result := range consumer {
        fmt.Printf("Result: %d\n", result)
    }
}
```

## 同步原语

### Mutex 互斥锁

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Mutex ===")
    
    var mu sync.Mutex
    var counter int
    
    // 启动多个 goroutine 修改共享变量
    for i := 0; i < 10; i++ {
        go func(id int) {
            for j := 0; j < 1000; j++ {
                mu.Lock()
                counter++
                mu.Unlock()
            }
            fmt.Printf("Goroutine %d completed\n", id)
        }(i)
    }
    
    time.Sleep(2 * time.Second)
    fmt.Printf("Final counter: %d\n", counter)
}
```

### RWMutex 读写锁

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== RWMutex ===")
    
    var rwmu sync.RWMutex
    var data map[string]int = make(map[string]int)
    
    // 写入数据
    go func() {
        for i := 0; i < 10; i++ {
            rwmu.Lock()
            data[fmt.Sprintf("key%d", i)] = i
            rwmu.Unlock()
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // 读取数据
    for i := 0; i < 5; i++ {
        go func(id int) {
            for j := 0; j < 10; j++ {
                rwmu.RLock()
                value := data[fmt.Sprintf("key%d", j)]
                rwmu.RUnlock()
                fmt.Printf("Reader %d: key%d = %d\n", id, j, value)
                time.Sleep(50 * time.Millisecond)
            }
        }(i)
    }
    
    time.Sleep(2 * time.Second)
}
```

### WaitGroup 等待组

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== WaitGroup ===")
    
    var wg sync.WaitGroup
    
    // 启动多个 goroutine
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d starting\n", id)
            time.Sleep(time.Duration(id) * 100 * time.Millisecond)
            fmt.Printf("Goroutine %d completed\n", id)
        }(i)
    }
    
    // 等待所有 goroutine 完成
    wg.Wait()
    fmt.Println("All goroutines completed")
}
```

### Once 单次执行

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    fmt.Println("=== Once ===")
    
    var once sync.Once
    var initialized bool
    
    // 多次调用，但只执行一次
    for i := 0; i < 5; i++ {
        go func(id int) {
            once.Do(func() {
                initialized = true
                fmt.Printf("Initialized by goroutine %d\n", id)
            })
            fmt.Printf("Goroutine %d: initialized = %t\n", id, initialized)
        }(i)
    }
    
    time.Sleep(100 * time.Millisecond)
}
```

## 并发模式

### Worker Pool 模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Worker Pool Pattern ===")
    
    const numWorkers = 3
    const numTasks = 10
    
    // 创建任务通道
    tasks := make(chan int, numTasks)
    results := make(chan int, numTasks)
    
    // 启动 worker
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for task := range tasks {
                fmt.Printf("Worker %d processing task %d\n", workerID, task)
                time.Sleep(100 * time.Millisecond)
                results <- task * 2
            }
        }(i)
    }
    
    // 发送任务
    go func() {
        for i := 0; i < numTasks; i++ {
            tasks <- i
        }
        close(tasks)
    }()
    
    // 收集结果
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 处理结果
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}
```

### Pipeline 模式

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    fmt.Println("=== Pipeline Pattern ===")
    
    // 创建管道
    input := make(chan int, 10)
    stage1 := make(chan int, 10)
    stage2 := make(chan int, 10)
    output := make(chan int, 10)
    
    // 输入阶段
    go func() {
        for i := 0; i < 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // 阶段1：乘以2
    go func() {
        for value := range input {
            stage1 <- value * 2
        }
        close(stage1)
    }()
    
    // 阶段2：加1
    go func() {
        for value := range stage1 {
            stage2 <- value + 1
        }
        close(stage2)
    }()
    
    // 输出阶段
    go func() {
        for value := range stage2 {
            output <- value
        }
        close(output)
    }()
    
    // 处理输出
    for result := range output {
        fmt.Printf("Result: %d\n", result)
    }
}
```

### Fan-out/Fan-in 模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Fan-out/Fan-in Pattern ===")
    
    // 输入通道
    input := make(chan int, 10)
    
    // 输出通道
    output := make(chan int, 10)
    
    // 启动输入
    go func() {
        for i := 0; i < 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // Fan-out：分发到多个 worker
    const numWorkers = 3
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for value := range input {
                fmt.Printf("Worker %d processing %d\n", workerID, value)
                time.Sleep(100 * time.Millisecond)
                output <- value * 2
            }
        }(i)
    }
    
    // Fan-in：收集结果
    go func() {
        wg.Wait()
        close(output)
    }()
    
    // 处理输出
    for result := range output {
        fmt.Printf("Result: %d\n", result)
    }
}
```

## 性能优化

### 减少 Goroutine 创建

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== Reduce Goroutine Creation ===")
    
    // 使用 worker pool 而不是为每个任务创建 goroutine
    const numWorkers = 4
    const numTasks = 1000
    
    tasks := make(chan int, numTasks)
    results := make(chan int, numTasks)
    
    // 启动固定数量的 worker
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for task := range tasks {
                results <- task * task
            }
        }(i)
    }
    
    // 发送任务
    go func() {
        for i := 0; i < numTasks; i++ {
            tasks <- i
        }
        close(tasks)
    }()
    
    // 收集结果
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 处理结果
    count := 0
    for range results {
        count++
    }
    
    fmt.Printf("Processed %d tasks\n", count)
}
```

### 避免 Goroutine 泄漏

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Avoid Goroutine Leak ===")
    
    // 使用 context 控制 goroutine 生命周期
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    // 启动受控的 goroutine
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine stopped by context")
                return
            case <-ticker.C:
                fmt.Println("Goroutine running...")
            }
        }
    }()
    
    // 等待 context 超时
    <-ctx.Done()
    fmt.Println("Main goroutine completed")
}
```

## 最佳实践

### 1. 合理使用 Goroutine

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    // 使用 WaitGroup 同步
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d\n", id)
        }(i)
    }
    
    wg.Wait()
}
```

### 2. 正确使用 Channel

```go
package main

import "fmt"

func main() {
    // 使用 channel 进行通信
    ch := make(chan int, 1)
    
    go func() {
        ch <- 42
    }()
    
    value := <-ch
    fmt.Printf("Received: %d\n", value)
}
```

### 3. 避免竞态条件

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var mu sync.Mutex
    var counter int
    
    // 使用互斥锁保护共享变量
    for i := 0; i < 10; i++ {
        go func() {
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    
    fmt.Printf("Counter: %d\n", counter)
}
```

### 4. 监控并发性能

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 监控 goroutine 数量
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
        }
    }()
    
    // 主程序逻辑
    time.Sleep(5 * time.Second)
}
```

## 总结

Go 的并发编程基于以下核心概念：

1. **Goroutine**: 轻量级协程，用户态线程
2. **Channel**: 用于 goroutine 间通信
3. **同步原语**: Mutex、RWMutex、WaitGroup、Once
4. **并发模式**: Worker Pool、Pipeline、Fan-out/Fan-in

**关键特性**:
- 简单易用：语法简洁，易于理解
- 高效性能：轻量级协程，低开销
- 安全并发：通过 channel 避免竞态条件
- 灵活模式：支持多种并发模式

**最佳实践**:
- 合理使用 goroutine，避免过度创建
- 使用 channel 进行通信，避免共享内存
- 正确使用同步原语，避免竞态条件
- 监控并发性能，及时发现问题

掌握 Go 的并发编程对于编写高效的 Go 程序至关重要。
