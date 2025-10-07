# Go Runtime 详解

## 📚 目录

- [Runtime 概述](#runtime-概述)
- [Goroutine 调度器](#goroutine-调度器)
- [内存管理](#内存管理)
- [垃圾回收器](#垃圾回收器)
- [网络轮询器](#网络轮询器)
- [系统调用](#系统调用)
- [运行时统计](#运行时统计)
- [调试和监控](#调试和监控)

## Runtime 概述

Go Runtime 是 Go 程序运行时的核心组件，负责管理内存、调度 goroutine、垃圾回收等关键功能。

### 核心组件

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 获取 Go 版本
    fmt.Printf("Go version: %s\n", runtime.Version())
    
    // 获取操作系统信息
    fmt.Printf("OS: %s\n", runtime.GOOS)
    fmt.Printf("Architecture: %s\n", runtime.GOARCH)
    
    // 获取 CPU 核心数
    fmt.Printf("CPU cores: %d\n", runtime.NumCPU())
    
    // 设置最大 CPU 使用数
    runtime.GOMAXPROCS(runtime.NumCPU())
    fmt.Printf("Max procs: %d\n", runtime.GOMAXPROCS(0))
    
    // 获取当前 goroutine 数量
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    
    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("Memory: %d KB\n", m.Alloc/1024)
}
```

### Runtime 包核心功能

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 1. Goroutine 管理
    fmt.Println("=== Goroutine 管理 ===")
    
    // 启动多个 goroutine
    for i := 0; i < 5; i++ {
        go func(id int) {
            fmt.Printf("Goroutine %d running\n", id)
            time.Sleep(100 * time.Millisecond)
        }(i)
    }
    
    // 等待所有 goroutine 完成
    time.Sleep(200 * time.Millisecond)
    fmt.Printf("Current goroutines: %d\n", runtime.NumGoroutine())
    
    // 2. 内存管理
    fmt.Println("\n=== 内存管理 ===")
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("Sys: %d KB\n", m.Sys/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
    
    // 3. 垃圾回收
    fmt.Println("\n=== 垃圾回收 ===")
    
    // 手动触发垃圾回收
    runtime.GC()
    
    // 设置垃圾回收目标百分比
    runtime.GC()
    fmt.Printf("GC completed\n")
    
    // 4. 栈管理
    fmt.Println("\n=== 栈管理 ===")
    
    // 获取当前栈大小
    stackSize := runtime.Stack(nil, false)
    fmt.Printf("Stack size: %d bytes\n", len(stackSize))
    
    // 5. 系统调用
    fmt.Println("\n=== 系统调用 ===")
    
    // 获取调用栈
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    fmt.Printf("Stack trace:\n%s\n", string(buf[:n]))
}
```

## Goroutine 调度器

### 调度器原理

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// 工作函数
func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d started\n", id)
    
    // 模拟工作
    for i := 0; i < 3; i++ {
        fmt.Printf("Worker %d: step %d\n", id, i+1)
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Printf("Worker %d finished\n", id)
}

func main() {
    // 设置最大 CPU 使用数
    runtime.GOMAXPROCS(2)
    
    var wg sync.WaitGroup
    
    // 启动多个 goroutine
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    // 等待所有 goroutine 完成
    wg.Wait()
    
    fmt.Printf("All workers completed\n")
    fmt.Printf("Final goroutine count: %d\n", runtime.NumGoroutine())
}
```

### 调度器状态监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func monitorScheduler() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        fmt.Printf("Goroutines: %d, CPUs: %d, MaxProcs: %d\n",
            runtime.NumGoroutine(),
            runtime.NumCPU(),
            runtime.GOMAXPROCS(0))
    }
}

func main() {
    // 启动监控
    go monitorScheduler()
    
    // 启动一些工作
    for i := 0; i < 10; i++ {
        go func(id int) {
            time.Sleep(2 * time.Second)
            fmt.Printf("Task %d completed\n", id)
        }(i)
    }
    
    // 运行一段时间
    time.Sleep(5 * time.Second)
}
```

## 内存管理

### 内存分配器

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 获取内存统计
    var m1, m2 runtime.MemStats
    
    runtime.ReadMemStats(&m1)
    
    // 分配一些内存
    data := make([]byte, 1024*1024) // 1MB
    for i := range data {
        data[i] = byte(i % 256)
    }
    
    runtime.ReadMemStats(&m2)
    
    // 计算内存使用
    fmt.Printf("Before allocation: %d KB\n", m1.Alloc/1024)
    fmt.Printf("After allocation: %d KB\n", m2.Alloc/1024)
    fmt.Printf("Memory used: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
    
    // 释放内存
    data = nil
    runtime.GC()
    
    var m3 runtime.MemStats
    runtime.ReadMemStats(&m3)
    fmt.Printf("After GC: %d KB\n", m3.Alloc/1024)
}
```

### 内存池

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// 对象池
type ObjectPool struct {
    pool sync.Pool
}

func NewObjectPool() *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 1024)
            },
        },
    }
}

func (p *ObjectPool) Get() []byte {
    return p.pool.Get().([]byte)
}

func (p *ObjectPool) Put(obj []byte) {
    // 清空对象
    for i := range obj {
        obj[i] = 0
    }
    p.pool.Put(obj)
}

func main() {
    pool := NewObjectPool()
    
    // 使用对象池
    for i := 0; i < 1000; i++ {
        obj := pool.Get()
        // 使用对象
        obj[0] = byte(i % 256)
        pool.Put(obj)
    }
    
    // 检查内存使用
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("Memory after pool usage: %d KB\n", m.Alloc/1024)
}
```

## 垃圾回收器

### GC 配置和监控

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

func main() {
    // 设置 GC 目标百分比
    debug.SetGCPercent(100)
    
    // 获取 GC 统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("GC Stats:\n")
    fmt.Printf("  NumGC: %d\n", m.NumGC)
    fmt.Printf("  PauseTotal: %v\n", time.Duration(m.PauseTotalNs))
    fmt.Printf("  PauseNs: %v\n", time.Duration(m.PauseNs[(m.NumGC+255)%256]))
    
    // 手动触发 GC
    fmt.Println("Triggering GC...")
    runtime.GC()
    
    // 再次获取统计
    runtime.ReadMemStats(&m)
    fmt.Printf("After GC:\n")
    fmt.Printf("  NumGC: %d\n", m.NumGC)
    fmt.Printf("  Alloc: %d KB\n", m.Alloc/1024)
}
```

### GC 性能测试

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func createGarbage() {
    // 创建大量垃圾
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024)
        data[0] = byte(i % 256)
    }
}

func main() {
    // 记录开始时间
    start := time.Now()
    
    // 创建垃圾
    createGarbage()
    
    // 记录创建时间
    createTime := time.Since(start)
    
    // 触发 GC
    gcStart := time.Now()
    runtime.GC()
    gcTime := time.Since(gcStart)
    
    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Create time: %v\n", createTime)
    fmt.Printf("GC time: %v\n", gcTime)
    fmt.Printf("Memory after GC: %d KB\n", m.Alloc/1024)
    fmt.Printf("GC count: %d\n", m.NumGC)
}
```

## 网络轮询器

### 网络 I/O 处理

```go
package main

import (
    "fmt"
    "net"
    "runtime"
    "time"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // 模拟处理
    time.Sleep(100 * time.Millisecond)
    
    // 发送响应
    conn.Write([]byte("Hello from server\n"))
}

func main() {
    // 监听端口
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Printf("Error listening: %v\n", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Server listening on :8080")
    
    // 处理连接
    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                fmt.Printf("Error accepting: %v\n", err)
                continue
            }
            
            // 每个连接一个 goroutine
            go handleConnection(conn)
        }
    }()
    
    // 监控 goroutine 数量
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    }
}
```

## 系统调用

### 系统调用监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 获取系统调用统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("System memory: %d KB\n", m.Sys/1024)
    fmt.Printf("Heap memory: %d KB\n", m.HeapSys/1024)
    fmt.Printf("Stack memory: %d KB\n", m.StackSys/1024)
    
    // 模拟系统调用
    for i := 0; i < 1000; i++ {
        // 分配内存（可能触发系统调用）
        data := make([]byte, 1024)
        data[0] = byte(i % 256)
        
        // 每100次检查一次内存
        if i%100 == 0 {
            runtime.ReadMemStats(&m)
            fmt.Printf("Iteration %d: Alloc=%d KB\n", i, m.Alloc/1024)
        }
    }
    
    // 最终统计
    runtime.ReadMemStats(&m)
    fmt.Printf("Final memory: %d KB\n", m.Alloc/1024)
}
```

## 运行时统计

### 详细统计信息

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("=== Memory Statistics ===\n")
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("Sys: %d KB\n", m.Sys/1024)
    fmt.Printf("Lookups: %d\n", m.Lookups)
    fmt.Printf("Mallocs: %d\n", m.Mallocs)
    fmt.Printf("Frees: %d\n", m.Frees)
    fmt.Printf("HeapAlloc: %d KB\n", m.HeapAlloc/1024)
    fmt.Printf("HeapSys: %d KB\n", m.HeapSys/1024)
    fmt.Printf("HeapIdle: %d KB\n", m.HeapIdle/1024)
    fmt.Printf("HeapInuse: %d KB\n", m.HeapInuse/1024)
    fmt.Printf("HeapReleased: %d KB\n", m.HeapReleased/1024)
    fmt.Printf("HeapObjects: %d\n", m.HeapObjects)
    fmt.Printf("StackInuse: %d KB\n", m.StackInuse/1024)
    fmt.Printf("StackSys: %d KB\n", m.StackSys/1024)
    fmt.Printf("MSpanInuse: %d KB\n", m.MSpanInuse/1024)
    fmt.Printf("MSpanSys: %d KB\n", m.MSpanSys/1024)
    fmt.Printf("MCacheInuse: %d KB\n", m.MCacheInuse/1024)
    fmt.Printf("MCacheSys: %d KB\n", m.MCacheSys/1024)
    fmt.Printf("BuckHashSys: %d KB\n", m.BuckHashSys/1024)
    fmt.Printf("GCSys: %d KB\n", m.GCSys/1024)
    fmt.Printf("OtherSys: %d KB\n", m.OtherSys/1024)
    fmt.Printf("NextGC: %d KB\n", m.NextGC/1024)
    fmt.Printf("LastGC: %v\n", time.Unix(0, int64(m.LastGC)))
    fmt.Printf("PauseTotalNs: %v\n", time.Duration(m.PauseTotalNs))
    fmt.Printf("NumGC: %d\n", m.NumGC)
    fmt.Printf("NumForcedGC: %d\n", m.NumForcedGC)
    fmt.Printf("GCCPUFraction: %.6f\n", m.GCCPUFraction)
    fmt.Printf("EnableGC: %t\n", m.EnableGC)
    fmt.Printf("DebugGC: %t\n", m.DebugGC)
    fmt.Printf("BySize: %v\n", m.BySize)
    fmt.Printf("========================\n")
}

func main() {
    // 初始统计
    printMemStats()
    
    // 分配一些内存
    data := make([]byte, 1024*1024) // 1MB
    for i := range data {
        data[i] = byte(i % 256)
    }
    
    fmt.Println("\nAfter allocation:")
    printMemStats()
    
    // 触发 GC
    runtime.GC()
    
    fmt.Println("\nAfter GC:")
    printMemStats()
}
```

## 调试和监控

### 运行时调试

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/pprof"
    "os"
)

func main() {
    // 1. 获取调用栈
    fmt.Println("=== Call Stack ===")
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    fmt.Printf("Stack trace:\n%s\n", string(buf[:n]))
    
    // 2. 获取所有 goroutine 的调用栈
    fmt.Println("\n=== All Goroutines ===")
    buf = make([]byte, 1024*1024)
    n = runtime.Stack(buf, true)
    fmt.Printf("All goroutines:\n%s\n", string(buf[:n]))
    
    // 3. 获取 goroutine ID
    fmt.Println("\n=== Goroutine ID ===")
    buf = make([]byte, 64)
    n = runtime.Stack(buf, false)
    fmt.Printf("Current goroutine: %s\n", string(buf[:n]))
    
    // 4. 内存分析
    fmt.Println("\n=== Memory Profile ===")
    f, err := os.Create("mem.prof")
    if err != nil {
        fmt.Printf("Error creating profile: %v\n", err)
        return
    }
    defer f.Close()
    
    pprof.WriteHeapProfile(f)
    fmt.Println("Memory profile written to mem.prof")
    
    // 5. CPU 分析
    fmt.Println("\n=== CPU Profile ===")
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        fmt.Printf("Error creating CPU profile: %v\n", err)
        return
    }
    defer cpuFile.Close()
    
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // 模拟一些工作
    for i := 0; i < 1000000; i++ {
        _ = i * i
    }
    
    fmt.Println("CPU profile written to cpu.prof")
}
```

### 性能监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

type Monitor struct {
    startTime time.Time
    lastGC    time.Time
}

func NewMonitor() *Monitor {
    return &Monitor{
        startTime: time.Now(),
    }
}

func (m *Monitor) Start() {
    go m.monitor()
}

func (m *Monitor) monitor() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        m.printStats()
    }
}

func (m *Monitor) printStats() {
    var mem runtime.MemStats
    runtime.ReadMemStats(&mem)
    
    uptime := time.Since(m.startTime)
    
    fmt.Printf("=== Runtime Stats ===\n")
    fmt.Printf("Uptime: %v\n", uptime)
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Printf("Memory: %d KB\n", mem.Alloc/1024)
    fmt.Printf("GC Count: %d\n", mem.NumGC)
    fmt.Printf("GC Time: %v\n", time.Duration(mem.PauseTotalNs))
    fmt.Printf("====================\n")
}

func main() {
    monitor := NewMonitor()
    monitor.Start()
    
    // 模拟一些工作
    for i := 0; i < 100; i++ {
        go func(id int) {
            time.Sleep(time.Duration(id) * 10 * time.Millisecond)
            fmt.Printf("Task %d completed\n", id)
        }(i)
    }
    
    // 运行一段时间
    time.Sleep(10 * time.Second)
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
    // 设置最大 CPU 使用数
    cpus := runtime.NumCPU()
    runtime.GOMAXPROCS(cpus)
    
    fmt.Printf("CPU cores: %d\n", cpus)
    fmt.Printf("Max procs: %d\n", runtime.GOMAXPROCS(0))
}
```

### 2. 监控内存使用

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func monitorMemory() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        if m.Alloc > 100*1024*1024 { // 100MB
            fmt.Printf("High memory usage: %d MB\n", m.Alloc/1024/1024)
        }
    }
}

func main() {
    go monitorMemory()
    
    // 主程序逻辑
    time.Sleep(30 * time.Second)
}
```

### 3. 合理使用垃圾回收

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
)

func main() {
    // 设置 GC 目标百分比
    debug.SetGCPercent(100)
    
    // 在关键时刻手动触发 GC
    defer func() {
        runtime.GC()
        fmt.Println("GC triggered on exit")
    }()
    
    // 程序逻辑
    fmt.Println("Program running...")
}
```

### 4. 避免内存泄漏

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 定期检查内存使用
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            if m.Alloc > 50*1024*1024 { // 50MB
                fmt.Printf("Memory usage: %d MB\n", m.Alloc/1024/1024)
            }
        }
    }()
    
    // 程序逻辑
    time.Sleep(10 * time.Second)
}
```

## 总结

Go Runtime 是 Go 程序运行的核心，提供了：

1. **Goroutine 调度**: 高效的并发调度器
2. **内存管理**: 自动内存分配和垃圾回收
3. **网络轮询**: 高效的网络 I/O 处理
4. **系统调用**: 与操作系统的接口
5. **运行时统计**: 详细的性能监控信息

理解 Go Runtime 的工作原理对于编写高效的 Go 程序至关重要。通过合理使用 Runtime 提供的功能，可以显著提高程序的性能和稳定性。
