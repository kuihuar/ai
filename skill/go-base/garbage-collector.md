# Go 垃圾回收器详解

## 📚 目录

- [垃圾回收器概述](#垃圾回收器概述)
- [三色标记算法](#三色标记算法)
- [并发垃圾回收](#并发垃圾回收)
- [GC 触发条件](#gc-触发条件)
- [GC 调优](#gc-调优)
- [GC 性能监控](#gc-性能监控)
- [内存泄漏检测](#内存泄漏检测)
- [最佳实践](#最佳实践)

## 垃圾回收器概述

Go 的垃圾回收器采用三色标记算法，实现了并发、低延迟的垃圾回收。

### GC 基本概念

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    fmt.Println("=== Garbage Collector Overview ===")
    
    // 获取GC统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("GC Statistics:\n")
    fmt.Printf("  NumGC: %d (GC次数)\n", m.NumGC)
    fmt.Printf("  PauseTotalNs: %v (GC暂停总时间)\n", time.Duration(m.PauseTotalNs))
    fmt.Printf("  PauseNs: %v (最近一次GC暂停时间)\n", time.Duration(m.PauseNs[(m.NumGC+255)%256]))
    fmt.Printf("  LastGC: %v (上次GC时间)\n", time.Unix(0, int64(m.LastGC)))
    fmt.Printf("  NextGC: %d KB (下次GC阈值)\n", m.NextGC/1024)
    fmt.Printf("  GCCPUFraction: %.6f (GC占用CPU比例)\n", m.GCCPUFraction)
    fmt.Printf("  EnableGC: %t (是否启用GC)\n", m.EnableGC)
    fmt.Printf("  DebugGC: %t (是否调试GC)\n", m.DebugGC)
    
    // 内存统计
    fmt.Printf("\nMemory Statistics:\n")
    fmt.Printf("  Alloc: %d KB (当前分配的内存)\n", m.Alloc/1024)
    fmt.Printf("  TotalAlloc: %d KB (累计分配的内存)\n", m.TotalAlloc/1024)
    fmt.Printf("  Sys: %d KB (从系统获得的内存)\n", m.Sys/1024)
    fmt.Printf("  HeapAlloc: %d KB (堆内存)\n", m.HeapAlloc/1024)
    fmt.Printf("  HeapSys: %d KB (堆系统内存)\n", m.HeapSys/1024)
    fmt.Printf("  HeapObjects: %d (堆对象数量)\n", m.HeapObjects)
}
```

### GC 工作流程

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func gcWorkflow() {
    fmt.Println("=== GC Workflow ===")
    
    // 1. 创建大量对象
    fmt.Println("1. Creating objects...")
    objects := make([]*[1024]byte, 10000)
    for i := 0; i < 10000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 获取分配后统计
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    fmt.Printf("   After allocation - Alloc: %d KB, Objects: %d\n", 
        m1.Alloc/1024, m1.HeapObjects)
    
    // 2. 释放部分对象
    fmt.Println("2. Releasing some objects...")
    for i := 0; i < 5000; i++ {
        objects[i] = nil
    }
    
    // 3. 手动触发GC
    fmt.Println("3. Triggering GC...")
    start := time.Now()
    runtime.GC()
    gcTime := time.Since(start)
    
    // 获取GC后统计
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("   GC time: %v\n", gcTime)
    fmt.Printf("   After GC - Alloc: %d KB, Objects: %d\n", 
        m2.Alloc/1024, m2.HeapObjects)
    fmt.Printf("   GC count: %d\n", m2.NumGC)
}
```

## 三色标记算法

### 三色标记原理

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

// 模拟三色标记算法
type Object struct {
    ID       int
    Children []*Object
    Marked   bool
}

func createObjectGraph() *Object {
    // 创建对象图
    root := &Object{ID: 1}
    
    // 创建子对象
    child1 := &Object{ID: 2}
    child2 := &Object{ID: 3}
    child3 := &Object{ID: 4}
    
    // 建立引用关系
    root.Children = []*Object{child1, child2}
    child1.Children = []*Object{child3}
    
    return root
}

func markObjects(root *Object) {
    // 三色标记算法
    // 1. 白色：未访问
    // 2. 灰色：已访问但子对象未访问
    // 3. 黑色：已访问且子对象已访问
    
    // 从根对象开始标记
    markObject(root)
}

func markObject(obj *Object) {
    if obj == nil || obj.Marked {
        return
    }
    
    // 标记为灰色
    obj.Marked = true
    
    // 递归标记子对象
    for _, child := range obj.Children {
        markObject(child)
    }
    
    // 标记为黑色
    fmt.Printf("Object %d marked as black\n", obj.ID)
}

func main() {
    fmt.Println("=== Three-Color Marking Algorithm ===")
    
    // 创建对象图
    root := createObjectGraph()
    
    // 执行标记
    markObjects(root)
    
    // 模拟GC统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("GC count: %d\n", m.NumGC)
}
```

### 标记阶段实现

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func markingPhase() {
    fmt.Println("=== Marking Phase ===")
    
    // 创建对象图
    objects := make([]*Object, 1000)
    for i := 0; i < 1000; i++ {
        objects[i] = &Object{ID: i}
    }
    
    // 建立引用关系
    for i := 0; i < 999; i++ {
        objects[i].Children = []*Object{objects[i+1]}
    }
    
    // 执行标记
    start := time.Now()
    markObjects(objects[0])
    markTime := time.Since(start)
    
    fmt.Printf("Marking time: %v\n", markTime)
    fmt.Printf("Objects marked: %d\n", countMarkedObjects(objects))
}

type Object struct {
    ID       int
    Children []*Object
    Marked   bool
}

func markObjects(root *Object) {
    if root == nil {
        return
    }
    
    // 使用栈进行迭代标记
    stack := []*Object{root}
    
    for len(stack) > 0 {
        obj := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        if obj.Marked {
            continue
        }
        
        obj.Marked = true
        
        // 将子对象加入栈
        for _, child := range obj.Children {
            if !child.Marked {
                stack = append(stack, child)
            }
        }
    }
}

func countMarkedObjects(objects []*Object) int {
    count := 0
    for _, obj := range objects {
        if obj.Marked {
            count++
        }
    }
    return count
}

func main() {
    markingPhase()
}
```

## 并发垃圾回收

### 并发标记

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func concurrentMarking() {
    fmt.Println("=== Concurrent Marking ===")
    
    // 创建大量对象
    objects := make([]*Object, 10000)
    for i := 0; i < 10000; i++ {
        objects[i] = &Object{ID: i}
    }
    
    // 建立引用关系
    for i := 0; i < 9999; i++ {
        objects[i].Children = []*Object{objects[i+1]}
    }
    
    // 并发标记
    start := time.Now()
    concurrentMarkObjects(objects[0])
    markTime := time.Since(start)
    
    fmt.Printf("Concurrent marking time: %v\n", markTime)
    fmt.Printf("Objects marked: %d\n", countMarkedObjects(objects))
}

type Object struct {
    ID       int
    Children []*Object
    Marked   bool
    mu       sync.Mutex
}

func concurrentMarkObjects(root *Object) {
    if root == nil {
        return
    }
    
    // 使用工作池进行并发标记
    const numWorkers = 4
    workChan := make(chan *Object, 1000)
    var wg sync.WaitGroup
    
    // 启动工作协程
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for obj := range workChan {
                markObject(obj)
            }
        }()
    }
    
    // 发送工作
    go func() {
        defer close(workChan)
        workChan <- root
    }()
    
    wg.Wait()
}

func markObject(obj *Object) {
    obj.mu.Lock()
    if obj.Marked {
        obj.mu.Unlock()
        return
    }
    obj.Marked = true
    obj.mu.Unlock()
    
    // 标记子对象
    for _, child := range obj.Children {
        if !child.Marked {
            markObject(child)
        }
    }
}

func countMarkedObjects(objects []*Object) int {
    count := 0
    for _, obj := range objects {
        if obj.Marked {
            count++
        }
    }
    return count
}

func main() {
    concurrentMarking()
}
```

### 并发清理

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func concurrentSweeping() {
    fmt.Println("=== Concurrent Sweeping ===")
    
    // 创建对象
    objects := make([]*Object, 1000)
    for i := 0; i < 1000; i++ {
        objects[i] = &Object{ID: i}
    }
    
    // 标记部分对象
    for i := 0; i < 500; i++ {
        objects[i].Marked = true
    }
    
    // 并发清理
    start := time.Now()
    concurrentSweepObjects(objects)
    sweepTime := time.Since(start)
    
    fmt.Printf("Concurrent sweeping time: %v\n", sweepTime)
    fmt.Printf("Objects remaining: %d\n", countMarkedObjects(objects))
}

type Object struct {
    ID       int
    Children []*Object
    Marked   bool
    mu       sync.Mutex
}

func concurrentSweepObjects(objects []*Object) {
    const numWorkers = 4
    workChan := make(chan *Object, 1000)
    var wg sync.WaitGroup
    
    // 启动工作协程
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for obj := range workChan {
                sweepObject(obj)
            }
        }()
    }
    
    // 发送工作
    go func() {
        defer close(workChan)
        for _, obj := range objects {
            workChan <- obj
        }
    }()
    
    wg.Wait()
}

func sweepObject(obj *Object) {
    obj.mu.Lock()
    if !obj.Marked {
        // 清理未标记的对象
        obj.Children = nil
    }
    obj.mu.Unlock()
}

func countMarkedObjects(objects []*Object) int {
    count := 0
    for _, obj := range objects {
        if obj.Marked {
            count++
        }
    }
    return count
}

func main() {
    concurrentSweeping()
}
```

## GC 触发条件

### 自动触发

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func gcTriggerConditions() {
    fmt.Println("=== GC Trigger Conditions ===")
    
    // 获取初始统计
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    fmt.Printf("Initial - Alloc: %d KB, NextGC: %d KB\n", 
        m1.Alloc/1024, m1.NextGC/1024)
    
    // 创建对象直到触发GC
    objects := make([]*[1024]byte, 0)
    
    for i := 0; i < 1000; i++ {
        obj := &[1024]byte{}
        objects = append(objects, obj)
        
        // 检查是否触发GC
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        if m.NumGC > m1.NumGC {
            fmt.Printf("GC triggered at iteration %d\n", i)
            fmt.Printf("  Alloc: %d KB, NextGC: %d KB\n", 
                m.Alloc/1024, m.NextGC/1024)
            break
        }
        
        if i%100 == 0 {
            fmt.Printf("Iteration %d - Alloc: %d KB\n", i, m.Alloc/1024)
        }
    }
}

func main() {
    gcTriggerConditions()
}
```

### 手动触发

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func manualGCTrigger() {
    fmt.Println("=== Manual GC Trigger ===")
    
    // 创建对象
    objects := make([]*[1024]byte, 1000)
    for i := 0; i < 1000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 获取分配后统计
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    fmt.Printf("Before GC - Alloc: %d KB, NumGC: %d\n", 
        m1.Alloc/1024, m1.NumGC)
    
    // 手动触发GC
    start := time.Now()
    runtime.GC()
    gcTime := time.Since(start)
    
    // 获取GC后统计
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("After GC - Alloc: %d KB, NumGC: %d\n", 
        m2.Alloc/1024, m2.NumGC)
    fmt.Printf("GC time: %v\n", gcTime)
}

func main() {
    manualGCTrigger()
}
```

## GC 调优

### GC 参数调优

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

func gcTuning() {
    fmt.Println("=== GC Tuning ===")
    
    // 设置GC目标百分比
    debug.SetGCPercent(100)
    
    // 获取GC统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("GC target percentage: 100%%\n")
    fmt.Printf("Current GC count: %d\n", m.NumGC)
    fmt.Printf("GC pause total: %v\n", time.Duration(m.PauseTotalNs))
    
    // 创建大量对象
    objects := make([]*[1024]byte, 50000)
    for i := 0; i < 50000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 等待GC
    time.Sleep(100 * time.Millisecond)
    
    // 获取GC后统计
    runtime.ReadMemStats(&m)
    fmt.Printf("After work - GC count: %d\n", m.NumGC)
    fmt.Printf("GC pause total: %v\n", time.Duration(m.PauseTotalNs))
    fmt.Printf("GC CPU fraction: %.6f\n", m.GCCPUFraction)
}

func main() {
    gcTuning()
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

func gcPerformanceTest() {
    fmt.Println("=== GC Performance Test ===")
    
    // 记录开始时间
    start := time.Now()
    
    // 创建大量对象
    objects := make([]*[1024]byte, 100000)
    for i := 0; i < 100000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 记录分配时间
    allocTime := time.Since(start)
    
    // 释放对象
    objects = nil
    
    // 记录GC时间
    gcStart := time.Now()
    runtime.GC()
    gcTime := time.Since(gcStart)
    
    // 获取统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Allocation time: %v\n", allocTime)
    fmt.Printf("GC time: %v\n", gcTime)
    fmt.Printf("Total time: %v\n", time.Since(start))
    fmt.Printf("GC count: %d\n", m.NumGC)
    fmt.Printf("GC pause: %v\n", time.Duration(m.PauseNs[(m.NumGC+255)%256]))
}

func main() {
    gcPerformanceTest()
}
```

## GC 性能监控

### GC 统计监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func gcMonitoring() {
    fmt.Println("=== GC Monitoring ===")
    
    // 创建监控器
    monitor := &GCMonitor{}
    monitor.Start()
    
    // 模拟工作负载
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024)
        data[0] = byte(i % 256)
        
        time.Sleep(1 * time.Millisecond)
    }
    
    // 运行一段时间
    time.Sleep(5 * time.Second)
}

type GCMonitor struct {
    startTime time.Time
    lastGC    time.Time
}

func (m *GCMonitor) Start() {
    m.startTime = time.Now()
    
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for range ticker.C {
            m.printGCStats()
        }
    }()
}

func (m *GCMonitor) printGCStats() {
    var mem runtime.MemStats
    runtime.ReadMemStats(&mem)
    
    if mem.NumGC > 0 {
        lastGC := time.Unix(0, int64(mem.LastGC))
        if lastGC.After(m.lastGC) {
            m.lastGC = lastGC
            fmt.Printf("GC #%d: %v, Pause: %v, Alloc: %d KB\n",
                mem.NumGC,
                lastGC,
                time.Duration(mem.PauseNs[(mem.NumGC+255)%256]),
                mem.Alloc/1024)
        }
    }
}

func main() {
    gcMonitoring()
}
```

### GC 性能分析

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/pprof"
    "os"
    "time"
)

func gcPerformanceAnalysis() {
    fmt.Println("=== GC Performance Analysis ===")
    
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
        data := make([]byte, 1024)
        data[0] = byte(i % 256)
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
    gcPerformanceAnalysis()
}
```

## 内存泄漏检测

### 内存泄漏检测

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func memoryLeakDetection() {
    fmt.Println("=== Memory Leak Detection ===")
    
    // 监控内存使用
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            fmt.Printf("Memory: %d KB, Goroutines: %d\n", 
                m.Alloc/1024, runtime.NumGoroutine())
        }
    }()
    
    // 模拟内存泄漏
    var leakedData [][]byte
    
    for i := 0; i < 100; i++ {
        // 分配内存但不释放
        data := make([]byte, 1024*1024) // 1MB
        leakedData = append(leakedData, data)
        
        time.Sleep(100 * time.Millisecond)
    }
    
    // 运行一段时间
    time.Sleep(5 * time.Second)
}

func main() {
    memoryLeakDetection()
}
```

### 内存使用监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func memoryUsageMonitoring() {
    fmt.Println("=== Memory Usage Monitoring ===")
    
    // 创建监控器
    monitor := &MemoryMonitor{}
    monitor.Start()
    
    // 模拟工作负载
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024)
        data[0] = byte(i % 256)
        
        time.Sleep(10 * time.Millisecond)
    }
    
    // 运行一段时间
    time.Sleep(5 * time.Second)
}

type MemoryMonitor struct {
    startTime time.Time
    maxMemory uint64
}

func (m *MemoryMonitor) Start() {
    m.startTime = time.Now()
    
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for range ticker.C {
            var mem runtime.MemStats
            runtime.ReadMemStats(&mem)
            
            if mem.Alloc > m.maxMemory {
                m.maxMemory = mem.Alloc
            }
            
            fmt.Printf("Memory: %d KB, Max: %d KB, Goroutines: %d\n",
                mem.Alloc/1024, m.maxMemory/1024, runtime.NumGoroutine())
        }
    }()
}

func main() {
    memoryUsageMonitoring()
}
```

## 最佳实践

### 1. 合理设置GC参数

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
)

func main() {
    // 设置GC目标百分比
    debug.SetGCPercent(100)
    
    // 获取GC统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("GC target percentage: 100%%\n")
    fmt.Printf("Current GC count: %d\n", m.NumGC)
}
```

### 2. 避免内存泄漏

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

### 3. 监控GC性能

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 定期监控GC
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            if m.NumGC > 0 {
                fmt.Printf("GC count: %d, Pause: %v\n", 
                    m.NumGC, 
                    time.Duration(m.PauseNs[(m.NumGC+255)%256]))
            }
        }
    }()
    
    // 主程序逻辑
    time.Sleep(10 * time.Second)
}
```

### 4. 优化内存分配

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func main() {
    // 使用对象池
    pool := &sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }
    
    // 获取对象
    obj := pool.Get().([]byte)
    
    // 使用对象
    obj[0] = 1
    
    // 归还对象
    pool.Put(obj)
}
```

## 总结

Go 的垃圾回收器是一个高效的内存管理系统：

1. **三色标记算法**: 白色、灰色、黑色标记对象状态
2. **并发回收**: 与程序并发运行，减少暂停时间
3. **自动触发**: 基于内存使用情况自动触发
4. **性能监控**: 提供详细的GC统计信息
5. **调优参数**: 可配置的GC参数

**关键特性**:
- 低延迟：GC暂停时间短
- 高并发：与程序并发运行
- 自适应：根据内存使用情况调整
- 可监控：提供详细的性能统计

**优化建议**:
- 合理设置GC参数
- 避免内存泄漏
- 监控GC性能
- 使用对象池减少分配

理解 Go 的垃圾回收机制对于编写高效的 Go 程序至关重要。
