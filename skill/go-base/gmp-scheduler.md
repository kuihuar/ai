# Go GMP 调度器详解

## 📚 目录

- [GMP 模型概述](#gmp-模型概述)
- [Goroutine (G)](#goroutine-g)
- [Machine (M)](#machine-m)
- [Processor (P)](#processor-p)
- [调度器工作原理](#调度器工作原理)
- [调度策略](#调度策略)
- [工作窃取算法](#工作窃取算法)
- [调度器状态](#调度器状态)
- [性能优化](#性能优化)
- [调试和监控](#调试和监控)

## GMP 模型概述

Go 的调度器采用 GMP 模型：
- **G (Goroutine)**: 轻量级线程，用户态协程
- **M (Machine)**: 系统线程，与操作系统线程一一对应
- **P (Processor)**: 逻辑处理器，管理 G 的执行

### 基本概念

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main() {
    // 获取系统信息
    fmt.Printf("CPU cores: %d\n", runtime.NumCPU())
    fmt.Printf("Max procs: %d\n", runtime.GOMAXPROCS(0))
    
    // 设置最大处理器数
    runtime.GOMAXPROCS(4)
    fmt.Printf("Set max procs to: %d\n", runtime.GOMAXPROCS(0))
    
    // 启动多个 goroutine
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d running on M\n", id)
            time.Sleep(100 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Printf("All goroutines completed\n")
}
```

## Goroutine (G)

### Goroutine 生命周期

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func goroutineLifecycle() {
    fmt.Println("=== Goroutine Lifecycle ===")
    
    // 1. 创建 goroutine
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // 2. 运行状态
            fmt.Printf("Goroutine %d: Running\n", id)
            
            // 3. 阻塞状态 (I/O)
            time.Sleep(100 * time.Millisecond)
            
            // 4. 运行状态
            fmt.Printf("Goroutine %d: Completed\n", id)
        }(i)
    }
    
    // 5. 等待完成
    wg.Wait()
    fmt.Println("All goroutines completed")
}

func main() {
    goroutineLifecycle()
    
    // 检查当前 goroutine 数量
    fmt.Printf("Current goroutines: %d\n", runtime.NumGoroutine())
}
```

### Goroutine 状态转换

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func stateTransitions() {
    fmt.Println("=== Goroutine State Transitions ===")
    
    var wg sync.WaitGroup
    ch := make(chan int, 1)
    
    // 1. 创建 -> 运行
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("State: Created -> Running")
        
        // 2. 运行 -> 阻塞 (等待通道)
        fmt.Println("State: Running -> Blocked (waiting for channel)")
        ch <- 1
        
        // 3. 阻塞 -> 运行
        fmt.Println("State: Blocked -> Running (received from channel)")
        <-ch
        
        // 4. 运行 -> 完成
        fmt.Println("State: Running -> Completed")
    }()
    
    wg.Wait()
    fmt.Println("Goroutine completed")
}

func main() {
    stateTransitions()
}
```

## Machine (M)

### M 的创建和管理

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func machineManagement() {
    fmt.Println("=== Machine Management ===")
    
    // 获取当前 M 数量
    fmt.Printf("Current goroutines: %d\n", runtime.NumGoroutine())
    
    // 创建大量 goroutine 来触发 M 创建
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            time.Sleep(10 * time.Millisecond)
        }(i)
    }
    
    // 等待一段时间
    time.Sleep(100 * time.Millisecond)
    
    wg.Wait()
    fmt.Printf("After goroutines: %d\n", runtime.NumGoroutine())
}

func main() {
    machineManagement()
}
```

### M 的阻塞和唤醒

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func machineBlocking() {
    fmt.Println("=== Machine Blocking ===")
    
    var wg sync.WaitGroup
    ch := make(chan int)
    
    // 创建阻塞的 goroutine
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d: Waiting for channel\n", id)
            <-ch
            fmt.Printf("Goroutine %d: Received from channel\n", id)
        }(i)
    }
    
    // 等待一段时间
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Goroutines waiting: %d\n", runtime.NumGoroutine())
    
    // 唤醒所有 goroutine
    for i := 0; i < 5; i++ {
        ch <- i
    }
    
    wg.Wait()
    fmt.Printf("All goroutines completed: %d\n", runtime.NumGoroutine())
}

func main() {
    machineBlocking()
}
```

## Processor (P)

### P 的作用和管理

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func processorManagement() {
    fmt.Println("=== Processor Management ===")
    
    // 设置 P 数量
    cpus := runtime.NumCPU()
    runtime.GOMAXPROCS(cpus)
    fmt.Printf("Set processors to: %d\n", runtime.GOMAXPROCS(0))
    
    // 创建大量 goroutine 来测试 P 的调度
    var wg sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d: Running on P\n", id)
            time.Sleep(50 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All goroutines completed")
}

func main() {
    processorManagement()
}
```

### P 的工作队列

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func processorWorkQueue() {
    fmt.Println("=== Processor Work Queue ===")
    
    // 创建大量短任务
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // 短任务
            _ = id * id
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All tasks completed")
}

func main() {
    processorWorkQueue()
}
```

## 调度器工作原理

### 调度循环

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func schedulingLoop() {
    fmt.Println("=== Scheduling Loop ===")
    
    // 创建不同优先级的任务
    var wg sync.WaitGroup
    
    // 高优先级任务
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("High priority task %d\n", id)
            time.Sleep(10 * time.Millisecond)
        }(i)
    }
    
    // 低优先级任务
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Low priority task %d\n", id)
            time.Sleep(5 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All tasks completed")
}

func main() {
    schedulingLoop()
}
```

### 抢占式调度

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func preemptiveScheduling() {
    fmt.Println("=== Preemptive Scheduling ===")
    
    var wg sync.WaitGroup
    
    // 创建长时间运行的任务
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Long task %d started\n", id)
            
            // 长时间运行
            start := time.Now()
            for time.Since(start) < 100*time.Millisecond {
                // 主动让出 CPU
                runtime.Gosched()
            }
            
            fmt.Printf("Long task %d completed\n", id)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All long tasks completed")
}

func main() {
    preemptiveScheduling()
}
```

## 调度策略

### 本地队列调度

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func localQueueScheduling() {
    fmt.Println("=== Local Queue Scheduling ===")
    
    // 创建大量短任务
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Task %d\n", id)
            time.Sleep(1 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All tasks completed")
}

func main() {
    localQueueScheduling()
}
```

### 全局队列调度

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func globalQueueScheduling() {
    fmt.Println("=== Global Queue Scheduling ===")
    
    // 创建大量任务来填满本地队列
    var wg sync.WaitGroup
    for i := 0; i < 200; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Global task %d\n", id)
            time.Sleep(1 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All global tasks completed")
}

func main() {
    globalQueueScheduling()
}
```

## 工作窃取算法

### 工作窃取实现

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func workStealing() {
    fmt.Println("=== Work Stealing ===")
    
    // 创建不平衡的工作负载
    var wg sync.WaitGroup
    
    // P1: 大量任务
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Heavy task %d\n", id)
            time.Sleep(1 * time.Millisecond)
        }(i)
    }
    
    // P2: 少量任务
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Light task %d\n", id)
            time.Sleep(1 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Work stealing completed")
}

func main() {
    workStealing()
}
```

### 负载均衡

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func loadBalancing() {
    fmt.Println("=== Load Balancing ===")
    
    // 创建不同负载的任务
    var wg sync.WaitGroup
    
    // 高负载任务
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("High load task %d\n", id)
            time.Sleep(10 * time.Millisecond)
        }(i)
    }
    
    // 低负载任务
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Low load task %d\n", id)
            time.Sleep(1 * time.Millisecond)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Load balancing completed")
}

func main() {
    loadBalancing()
}
```

## 调度器状态

### 状态监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func schedulerMonitoring() {
    fmt.Println("=== Scheduler Monitoring ===")
    
    // 监控调度器状态
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            fmt.Printf("Goroutines: %d, CPUs: %d\n", 
                runtime.NumGoroutine(), 
                runtime.NumCPU())
        }
    }()
    
    // 创建一些工作
    for i := 0; i < 10; i++ {
        go func(id int) {
            time.Sleep(200 * time.Millisecond)
            fmt.Printf("Task %d completed\n", id)
        }(i)
    }
    
    time.Sleep(1 * time.Second)
}

func main() {
    schedulerMonitoring()
}
```

### 调度器统计

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func schedulerStats() {
    fmt.Println("=== Scheduler Statistics ===")
    
    // 获取调度器统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Printf("CPU cores: %d\n", runtime.NumCPU())
    fmt.Printf("Max procs: %d\n", runtime.GOMAXPROCS(0))
    fmt.Printf("Memory: %d KB\n", m.Alloc/1024)
    
    // 创建一些工作
    for i := 0; i < 100; i++ {
        go func(id int) {
            time.Sleep(10 * time.Millisecond)
        }(i)
    }
    
    time.Sleep(100 * time.Millisecond)
    
    // 再次获取统计
    runtime.ReadMemStats(&m)
    fmt.Printf("After work - Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Printf("After work - Memory: %d KB\n", m.Alloc/1024)
}

func main() {
    schedulerStats()
}
```

## 性能优化

### 减少 Goroutine 创建

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// 使用 worker pool 模式
func workerPool() {
    fmt.Println("=== Worker Pool ===")
    
    const numWorkers = 4
    const numTasks = 100
    
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
                // 处理任务
                result := task * task
                results <- result
                fmt.Printf("Worker %d processed task %d\n", workerID, task)
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
        _ = result
    }
    
    fmt.Println("Worker pool completed")
}

func main() {
    workerPool()
}
```

### 避免 Goroutine 泄漏

```go
package main

import (
    "context"
    "fmt"
    "runtime"
    "time"
)

func avoidGoroutineLeak() {
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
    
    // 检查 goroutine 数量
    fmt.Printf("Goroutines after timeout: %d\n", runtime.NumGoroutine())
}

func main() {
    avoidGoroutineLeak()
}
```

### 优化调度性能

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func optimizeScheduling() {
    fmt.Println("=== Optimize Scheduling ===")
    
    // 设置合适的 P 数量
    cpus := runtime.NumCPU()
    runtime.GOMAXPROCS(cpus)
    
    // 使用批量处理
    const batchSize = 100
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i += batchSize {
        wg.Add(1)
        go func(start int) {
            defer wg.Done()
            
            // 批量处理
            for j := start; j < start+batchSize && j < 1000; j++ {
                _ = j * j
            }
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Batch processing completed")
}

func main() {
    optimizeScheduling()
}
```

## 调试和监控

### 调度器调试

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func schedulerDebugging() {
    fmt.Println("=== Scheduler Debugging ===")
    
    // 获取当前 goroutine 信息
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    fmt.Printf("Current goroutine:\n%s\n", string(buf[:n]))
    
    // 获取所有 goroutine 信息
    buf = make([]byte, 1024*1024)
    n = runtime.Stack(buf, true)
    fmt.Printf("All goroutines:\n%s\n", string(buf[:n]))
}

func main() {
    schedulerDebugging()
}
```

### 性能分析

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/pprof"
    "os"
    "time"
)

func performanceAnalysis() {
    fmt.Println("=== Performance Analysis ===")
    
    // CPU 分析
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        fmt.Printf("Error creating CPU profile: %v\n", err)
        return
    }
    defer cpuFile.Close()
    
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // 运行一些工作
    for i := 0; i < 1000000; i++ {
        _ = i * i
    }
    
    // 内存分析
    memFile, err := os.Create("mem.prof")
    if err != nil {
        fmt.Printf("Error creating memory profile: %v\n", err)
        return
    }
    defer memFile.Close()
    
    pprof.WriteHeapProfile(memFile)
    
    fmt.Println("Profiles created: cpu.prof, mem.prof")
}

func main() {
    performanceAnalysis()
}
```

## 最佳实践

### 1. 合理设置 GOMAXPROCS

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    // 设置 P 数量为 CPU 核心数
    cpus := runtime.NumCPU()
    runtime.GOMAXPROCS(cpus)
    fmt.Printf("Set processors to: %d\n", cpus)
}
```

### 2. 使用 Worker Pool 模式

```go
package main

import (
    "fmt"
    "sync"
)

func workerPool() {
    const numWorkers = 4
    const numTasks = 100
    
    tasks := make(chan int, numTasks)
    results := make(chan int, numTasks)
    
    // 启动 worker
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range tasks {
                results <- task * task
            }
        }()
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
        _ = result
    }
}

func main() {
    workerPool()
}
```

### 3. 避免 Goroutine 泄漏

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 使用 context 控制生命周期
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                // 工作
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
    
    <-ctx.Done()
}
```

### 4. 监控调度器性能

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 定期监控
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
        }
    }()
    
    // 主程序逻辑
    time.Sleep(10 * time.Second)
}
```

## 总结

Go 的 GMP 调度器是一个高效的多线程调度系统：

1. **G (Goroutine)**: 轻量级协程，用户态线程
2. **M (Machine)**: 系统线程，与操作系统线程对应
3. **P (Processor)**: 逻辑处理器，管理 G 的执行

**核心特性**:
- 工作窃取算法实现负载均衡
- 抢占式调度避免饥饿
- 本地队列提高缓存效率
- 全局队列处理溢出任务

**优化建议**:
- 合理设置 GOMAXPROCS
- 使用 Worker Pool 模式
- 避免 Goroutine 泄漏
- 监控调度器性能

理解 GMP 调度器的工作原理对于编写高效的 Go 并发程序至关重要。
