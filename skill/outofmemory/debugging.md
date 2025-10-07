# Go 内存溢出排查方法

## 🔍 基本内存监控

### 1. 使用 runtime 包监控内存

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func basicMemoryMonitoring() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("当前内存使用: %d KB\n", m.Alloc/1024)
    fmt.Printf("累计分配: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("系统内存: %d KB\n", m.Sys/1024)
    fmt.Printf("堆内存: %d KB\n", m.HeapAlloc/1024)
    fmt.Printf("GC 次数: %d\n", m.NumGC)
    fmt.Printf("GC 暂停时间: %d ns\n", m.PauseTotalNs)
}
```

### 2. 内存使用趋势监控

```go
func memoryTrendMonitoring() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            fmt.Printf("时间: %s, 内存: %d KB, GC: %d\n", 
                time.Now().Format("15:04:05"), m.Alloc/1024, m.NumGC)
        }
    }()
}
```

## 🐛 内存泄漏检测

### 1. 基本内存泄漏检测

```go
func memoryLeakDetection() {
    var data [][]byte
    
    for i := 0; i < 100; i++ {
        chunk := make([]byte, 1024*1024) // 1MB
        data = append(data, chunk)
        
        if i%10 == 0 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            fmt.Printf("分配 %d 个块后，内存使用: %d KB\n", i+1, m.Alloc/1024)
        }
    }
    
    // 这里应该释放 data，但故意不释放来模拟内存泄漏
    // data = nil
}
```

### 2. 内存使用率检查

```go
func memoryUsageCheck() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 计算内存使用率
    usagePercent := float64(m.Alloc) / float64(m.Sys) * 100
    
    fmt.Printf("内存使用率: %.2f%%\n", usagePercent)
    
    if usagePercent > 80 {
        fmt.Println("警告: 内存使用率过高!")
        runtime.GC() // 强制垃圾回收
    }
}
```

## 📊 使用 pprof 进行内存分析

### 1. 启动 pprof 服务器

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "time"
)

func startPprofServer() {
    go func() {
        log.Println("pprof 服务器启动在 :6060")
        log.Println("访问 http://localhost:6060/debug/pprof/ 查看内存信息")
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}

func main() {
    // 启动 pprof
    startPprofServer()
    
    // 模拟内存使用
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024*1024) // 1MB
        _ = data
        
        if i%100 == 0 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            log.Printf("内存使用: %d KB", m.Alloc/1024)
        }
        
        time.Sleep(100 * time.Millisecond)
    }
}
```

### 2. 命令行分析工具

```bash
# 启动程序
go run main.go

# 在另一个终端中分析内存
go tool pprof http://localhost:6060/debug/pprof/heap

# 分析内存分配
go tool pprof http://localhost:6060/debug/pprof/allocs

# 分析内存使用趋势
go tool pprof http://localhost:6060/debug/pprof/heap?seconds=30

# 生成内存使用图
go tool pprof -png http://localhost:6060/debug/pprof/heap > heap.png
```

## 🔧 内存分配统计

### 1. 详细内存统计

```go
func detailedMemoryUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("=== 堆内存 ===\n")
    fmt.Printf("堆分配: %d KB\n", m.HeapAlloc/1024)
    fmt.Printf("堆系统: %d KB\n", m.HeapSys/1024)
    fmt.Printf("堆空闲: %d KB\n", m.HeapIdle/1024)
    fmt.Printf("堆使用: %d KB\n", m.HeapInuse/1024)
    fmt.Printf("堆释放: %d KB\n", m.HeapReleased/1024)
    fmt.Printf("堆对象数: %d\n", m.HeapObjects)
    
    fmt.Printf("\n=== 栈内存 ===\n")
    fmt.Printf("栈使用: %d KB\n", m.StackInuse/1024)
    fmt.Printf("栈系统: %d KB\n", m.StackSys/1024)
    
    fmt.Printf("\n=== 其他内存 ===\n")
    fmt.Printf("MSpan 使用: %d KB\n", m.MSpanInuse/1024)
    fmt.Printf("MSpan 系统: %d KB\n", m.MSpanSys/1024)
    fmt.Printf("MCache 使用: %d KB\n", m.MCacheInuse/1024)
    fmt.Printf("MCache 系统: %d KB\n", m.MCacheSys/1024)
    fmt.Printf("哈希表: %d KB\n", m.BuckHashSys/1024)
    fmt.Printf("GC 系统: %d KB\n", m.GCSys/1024)
    fmt.Printf("其他系统: %d KB\n", m.OtherSys/1024)
}
```

### 2. 内存分配统计

```go
func memoryAllocationStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("分配次数: %d\n", m.Mallocs)
    fmt.Printf("释放次数: %d\n", m.Frees)
    fmt.Printf("净分配: %d\n", m.Mallocs-m.Frees)
    fmt.Printf("平均分配大小: %d 字节\n", m.TotalAlloc/uint64(m.Mallocs))
}
```

## ⚡ GC 性能分析

### 1. GC 性能监控

```go
func gcPerformanceAnalysis() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("GC 次数: %d\n", m.NumGC)
    fmt.Printf("强制 GC 次数: %d\n", m.NumForcedGC)
    fmt.Printf("GC 暂停总时间: %d ns\n", m.PauseTotalNs)
    fmt.Printf("平均 GC 暂停时间: %d ns\n", m.PauseTotalNs/uint64(m.NumGC))
    fmt.Printf("GC CPU 使用率: %.2f%%\n", m.GCCPUFraction*100)
}
```

### 2. 内存使用历史记录

```go
type MemorySnapshot struct {
    Timestamp time.Time
    Alloc     uint64
    TotalAlloc uint64
    NumGC     uint32
}

func memoryHistoryTracking() {
    var history []MemorySnapshot
    
    // 记录内存使用历史
    for i := 0; i < 10; i++ {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        snapshot := MemorySnapshot{
            Timestamp:  time.Now(),
            Alloc:      m.Alloc,
            TotalAlloc: m.TotalAlloc,
            NumGC:      m.NumGC,
        }
        
        history = append(history, snapshot)
        
        // 模拟内存使用
        data := make([]byte, 1024*1024) // 1MB
        _ = data
        
        time.Sleep(100 * time.Millisecond)
    }
    
    // 打印历史记录
    for i, snapshot := range history {
        fmt.Printf("快照 %d: 时间=%s, 内存=%d KB, GC=%d\n", 
            i+1, 
            snapshot.Timestamp.Format("15:04:05"), 
            snapshot.Alloc/1024, 
            snapshot.NumGC)
    }
}
```

## 🚨 内存使用警告

### 1. 内存使用率警告

```go
func memoryUsageWarning() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 检查内存使用率
    usagePercent := float64(m.Alloc) / float64(m.Sys) * 100
    
    if usagePercent > 90 {
        fmt.Println("严重警告: 内存使用率超过 90%!")
        runtime.GC()
    } else if usagePercent > 80 {
        fmt.Println("警告: 内存使用率超过 80%")
        runtime.GC()
    } else if usagePercent > 70 {
        fmt.Println("注意: 内存使用率超过 70%")
    } else {
        fmt.Println("内存使用正常")
    }
}
```

### 2. 内存使用预测

```go
func memoryUsagePrediction() {
    var history []uint64
    
    // 收集历史数据
    for i := 0; i < 10; i++ {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        history = append(history, m.Alloc)
        
        // 模拟内存使用
        data := make([]byte, 1024*1024) // 1MB
        _ = data
        
        time.Sleep(100 * time.Millisecond)
    }
    
    // 简单线性预测
    if len(history) >= 2 {
        growth := history[len(history)-1] - history[0]
        avgGrowth := growth / uint64(len(history)-1)
        
        predicted := history[len(history)-1] + avgGrowth*5 // 预测5步后
        
        fmt.Printf("当前内存: %d KB\n", history[len(history)-1]/1024)
        fmt.Printf("预测5步后内存: %d KB\n", predicted/1024)
    }
}
```

## 🛠️ 内存限制设置

### 1. 设置内存限制

```go
func setMemoryLimits() {
    // 设置 GC 目标百分比
    debug.SetGCPercent(100)
    fmt.Println("GC 目标百分比设置为 100%")
    
    // 设置内存限制 (Go 1.19+)
    debug.SetMemoryLimit(100 * 1024 * 1024) // 100MB
    fmt.Println("内存限制设置为 100MB")
    
    // 设置最大栈大小
    debug.SetMaxStack(64 * 1024 * 1024) // 64MB
    fmt.Println("最大栈大小设置为 64MB")
}
```

## 📈 内存使用分析工具

### 1. 使用 go tool trace

```go
package main

import (
    "os"
    "runtime/trace"
    "time"
)

func main() {
    // 开始跟踪
    f, err := os.Create("trace.out")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    
    err = trace.Start(f)
    if err != nil {
        panic(err)
    }
    defer trace.Stop()
    
    // 模拟工作
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024*1024) // 1MB
        _ = data
        time.Sleep(1 * time.Millisecond)
    }
}
```

```bash
# 分析跟踪文件
go tool trace trace.out
```

### 2. 内存使用监控器

```go
type MemoryMonitor struct {
    maxMemory     uint64
    checkInterval time.Duration
    stopCh        chan struct{}
    onWarning     func(uint64)
}

func NewMemoryMonitor(maxMemoryMB uint64, checkInterval time.Duration) *MemoryMonitor {
    return &MemoryMonitor{
        maxMemory:     maxMemoryMB * 1024 * 1024,
        checkInterval: checkInterval,
        stopCh:        make(chan struct{}),
    }
}

func (mm *MemoryMonitor) Start() {
    go func() {
        ticker := time.NewTicker(mm.checkInterval)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                mm.checkMemory()
            case <-mm.stopCh:
                return
            }
        }
    }()
}

func (mm *MemoryMonitor) checkMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    if m.Alloc > mm.maxMemory {
        if mm.onWarning != nil {
            mm.onWarning(m.Alloc)
        }
        
        // 触发垃圾回收
        runtime.GC()
    }
}
```

## 🔍 排查步骤总结

1. **基本监控**: 使用 `runtime.ReadMemStats()` 监控内存使用
2. **趋势分析**: 定期记录内存使用情况，分析增长趋势
3. **泄漏检测**: 检查内存是否持续增长而不释放
4. **pprof 分析**: 使用 pprof 工具深入分析内存使用
5. **GC 分析**: 监控垃圾回收性能
6. **设置限制**: 设置内存使用限制和警告
7. **历史记录**: 记录内存使用历史，便于分析
8. **预测分析**: 基于历史数据预测内存使用趋势

## 📚 常用命令

```bash
# 查看内存使用
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看内存分配
go tool pprof http://localhost:6060/debug/pprof/allocs

# 查看内存使用趋势
go tool pprof http://localhost:6060/debug/pprof/heap?seconds=30

# 生成内存使用图
go tool pprof -png http://localhost:6060/debug/pprof/heap > heap.png

# 分析跟踪文件
go tool trace trace.out
```
