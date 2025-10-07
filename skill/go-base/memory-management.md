# Go 内存管理详解

## 📚 目录

- [内存管理概述](#内存管理概述)
- [内存分配器](#内存分配器)
- [垃圾回收器](#垃圾回收器)
- [内存布局](#内存布局)
- [内存池](#内存池)
- [内存泄漏检测](#内存泄漏检测)
- [性能优化](#性能优化)
- [调试和监控](#调试和监控)

## 内存管理概述

Go 的内存管理采用自动内存管理，包括内存分配、垃圾回收和内存优化。

### 内存管理组件

```go
package main

import (
    "fmt"
    "runtime"
    "unsafe"
)

func main() {
    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("=== Memory Management Overview ===\n")
    fmt.Printf("Alloc: %d KB (当前分配的内存)\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc: %d KB (累计分配的内存)\n", m.TotalAlloc/1024)
    fmt.Printf("Sys: %d KB (从系统获得的内存)\n", m.Sys/1024)
    fmt.Printf("HeapAlloc: %d KB (堆内存)\n", m.HeapAlloc/1024)
    fmt.Printf("HeapSys: %d KB (堆系统内存)\n", m.HeapSys/1024)
    fmt.Printf("HeapIdle: %d KB (空闲堆内存)\n", m.HeapIdle/1024)
    fmt.Printf("HeapInuse: %d KB (使用中的堆内存)\n", m.HeapInuse/1024)
    fmt.Printf("HeapReleased: %d KB (释放给系统的堆内存)\n", m.HeapReleased/1024)
    fmt.Printf("HeapObjects: %d (堆对象数量)\n", m.HeapObjects)
    fmt.Printf("StackInuse: %d KB (栈内存)\n", m.StackInuse/1024)
    fmt.Printf("StackSys: %d KB (栈系统内存)\n", m.StackSys/1024)
    fmt.Printf("MSpanInuse: %d KB (MSpan使用内存)\n", m.MSpanInuse/1024)
    fmt.Printf("MSpanSys: %d KB (MSpan系统内存)\n", m.MSpanSys/1024)
    fmt.Printf("MCacheInuse: %d KB (MCache使用内存)\n", m.MCacheInuse/1024)
    fmt.Printf("MCacheSys: %d KB (MCache系统内存)\n", m.MCacheSys/1024)
    fmt.Printf("BuckHashSys: %d KB (哈希表内存)\n", m.BuckHashSys/1024)
    fmt.Printf("GCSys: %d KB (GC系统内存)\n", m.GCSys/1024)
    fmt.Printf("OtherSys: %d KB (其他系统内存)\n", m.OtherSys/1024)
    fmt.Printf("NextGC: %d KB (下次GC阈值)\n", m.NextGC/1024)
    fmt.Printf("LastGC: %v (上次GC时间)\n", m.LastGC)
    fmt.Printf("PauseTotalNs: %v (GC暂停总时间)\n", m.PauseTotalNs)
    fmt.Printf("NumGC: %d (GC次数)\n", m.NumGC)
    fmt.Printf("NumForcedGC: %d (强制GC次数)\n", m.NumForcedGC)
    fmt.Printf("GCCPUFraction: %.6f (GC占用CPU比例)\n", m.GCCPUFraction)
    fmt.Printf("EnableGC: %t (是否启用GC)\n", m.EnableGC)
    fmt.Printf("DebugGC: %t (是否调试GC)\n", m.DebugGC)
}
```

## 内存分配器

### 内存分配策略

```go
package main

import (
    "fmt"
    "runtime"
    "unsafe"
)

func memoryAllocation() {
    fmt.Println("=== Memory Allocation ===")
    
    // 小对象分配 (<= 32KB)
    smallObj := make([]byte, 1024)
    fmt.Printf("Small object size: %d bytes\n", len(smallObj))
    
    // 大对象分配 (> 32KB)
    largeObj := make([]byte, 64*1024)
    fmt.Printf("Large object size: %d bytes\n", len(largeObj))
    
    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("Mallocs: %d\n", m.Mallocs)
    fmt.Printf("Frees: %d\n", m.Frees)
    fmt.Printf("HeapObjects: %d\n", m.HeapObjects)
}

func main() {
    memoryAllocation()
}
```

### 内存分配器类型

```go
package main

import (
    "fmt"
    "runtime"
    "unsafe"
)

// 小对象分配器
func smallObjectAllocator() {
    fmt.Println("=== Small Object Allocator ===")
    
    // 分配小对象
    objects := make([]*[1024]byte, 1000)
    for i := 0; i < 1000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 获取统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Small objects allocated: %d\n", len(objects))
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("Mallocs: %d\n", m.Mallocs)
}

// 大对象分配器
func largeObjectAllocator() {
    fmt.Println("=== Large Object Allocator ===")
    
    // 分配大对象
    objects := make([]*[64*1024]byte, 10)
    for i := 0; i < 10; i++ {
        objects[i] = &[64*1024]byte{}
    }
    
    // 获取统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Large objects allocated: %d\n", len(objects))
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("Mallocs: %d\n", m.Mallocs)
}

func main() {
    smallObjectAllocator()
    largeObjectAllocator()
}
```

### 内存对齐

```go
package main

import (
    "fmt"
    "unsafe"
)

// 内存对齐示例
type AlignedStruct struct {
    a bool    // 1 byte
    b int32   // 4 bytes
    c int64   // 8 bytes
    d string  // 16 bytes
}

type UnalignedStruct struct {
    a int64   // 8 bytes
    b bool    // 1 byte
    c int32   // 4 bytes
    d string  // 16 bytes
}

func main() {
    fmt.Println("=== Memory Alignment ===")
    
    // 对齐的结构体
    aligned := AlignedStruct{}
    fmt.Printf("Aligned struct size: %d bytes\n", unsafe.Sizeof(aligned))
    fmt.Printf("Aligned struct alignment: %d bytes\n", unsafe.Alignof(aligned))
    
    // 未对齐的结构体
    unaligned := UnalignedStruct{}
    fmt.Printf("Unaligned struct size: %d bytes\n", unsafe.Sizeof(unaligned))
    fmt.Printf("Unaligned struct alignment: %d bytes\n", unsafe.Alignof(unaligned))
    
    // 字段偏移
    fmt.Printf("Aligned struct field offsets:\n")
    fmt.Printf("  a: %d\n", unsafe.Offsetof(aligned.a))
    fmt.Printf("  b: %d\n", unsafe.Offsetof(aligned.b))
    fmt.Printf("  c: %d\n", unsafe.Offsetof(aligned.c))
    fmt.Printf("  d: %d\n", unsafe.Offsetof(aligned.d))
    
    fmt.Printf("Unaligned struct field offsets:\n")
    fmt.Printf("  a: %d\n", unsafe.Offsetof(unaligned.a))
    fmt.Printf("  b: %d\n", unsafe.Offsetof(unaligned.b))
    fmt.Printf("  c: %d\n", unsafe.Offsetof(unaligned.c))
    fmt.Printf("  d: %d\n", unsafe.Offsetof(unaligned.d))
}
```

## 垃圾回收器

### GC 工作原理

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func garbageCollection() {
    fmt.Println("=== Garbage Collection ===")
    
    // 获取GC前统计
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    fmt.Printf("Before GC - Alloc: %d KB, NumGC: %d\n", m1.Alloc/1024, m1.NumGC)
    
    // 创建大量对象
    objects := make([]*[1024]byte, 10000)
    for i := 0; i < 10000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 获取分配后统计
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("After allocation - Alloc: %d KB, NumGC: %d\n", m2.Alloc/1024, m2.NumGC)
    
    // 释放对象
    objects = nil
    
    // 手动触发GC
    runtime.GC()
    
    // 获取GC后统计
    var m3 runtime.MemStats
    runtime.ReadMemStats(&m3)
    fmt.Printf("After GC - Alloc: %d KB, NumGC: %d\n", m3.Alloc/1024, m3.NumGC)
    fmt.Printf("GC pause: %v\n", time.Duration(m3.PauseNs[(m3.NumGC+255)%256]))
}

func main() {
    garbageCollection()
}
```

### GC 调优

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

## 内存布局

### 堆内存布局

```go
package main

import (
    "fmt"
    "runtime"
    "unsafe"
)

func heapLayout() {
    fmt.Println("=== Heap Memory Layout ===")
    
    // 获取堆统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Heap memory:\n")
    fmt.Printf("  Alloc: %d KB (当前分配)\n", m.HeapAlloc/1024)
    fmt.Printf("  Sys: %d KB (系统内存)\n", m.HeapSys/1024)
    fmt.Printf("  Idle: %d KB (空闲)\n", m.HeapIdle/1024)
    fmt.Printf("  Inuse: %d KB (使用中)\n", m.HeapInuse/1024)
    fmt.Printf("  Released: %d KB (释放)\n", m.HeapReleased/1024)
    fmt.Printf("  Objects: %d (对象数量)\n", m.HeapObjects)
    
    // 创建一些对象
    objects := make([]*[1024]byte, 1000)
    for i := 0; i < 1000; i++ {
        objects[i] = &[1024]byte{}
    }
    
    // 再次获取统计
    runtime.ReadMemStats(&m)
    fmt.Printf("\nAfter allocation:\n")
    fmt.Printf("  Alloc: %d KB\n", m.HeapAlloc/1024)
    fmt.Printf("  Objects: %d\n", m.HeapObjects)
}

func main() {
    heapLayout()
}
```

### 栈内存布局

```go
package main

import (
    "fmt"
    "runtime"
    "unsafe"
)

func stackLayout() {
    fmt.Println("=== Stack Memory Layout ===")
    
    // 获取栈统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Stack memory:\n")
    fmt.Printf("  Inuse: %d KB (使用中)\n", m.StackInuse/1024)
    fmt.Printf("  Sys: %d KB (系统内存)\n", m.StackSys/1024)
    
    // 递归函数测试栈
    testStack(0)
}

func testStack(depth int) {
    if depth >= 10 {
        return
    }
    
    // 在栈上分配一些数据
    data := [1024]byte{}
    data[0] = byte(depth)
    
    // 递归调用
    testStack(depth + 1)
    
    // 防止优化
    _ = data
}

func main() {
    stackLayout()
}
```

### 内存段布局

```go
package main

import (
    "fmt"
    "runtime"
    "unsafe"
)

func memorySegmentLayout() {
    fmt.Println("=== Memory Segment Layout ===")
    
    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Memory segments:\n")
    fmt.Printf("  Heap: %d KB\n", m.HeapSys/1024)
    fmt.Printf("  Stack: %d KB\n", m.StackSys/1024)
    fmt.Printf("  MSpan: %d KB\n", m.MSpanSys/1024)
    fmt.Printf("  MCache: %d KB\n", m.MCacheSys/1024)
    fmt.Printf("  BuckHash: %d KB\n", m.BuckHashSys/1024)
    fmt.Printf("  GC: %d KB\n", m.GCSys/1024)
    fmt.Printf("  Other: %d KB\n", m.OtherSys/1024)
    fmt.Printf("  Total: %d KB\n", m.Sys/1024)
}

func main() {
    memorySegmentLayout()
}
```

## 内存池

### 对象池实现

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
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
    fmt.Println("=== Object Pool ===")
    
    // 创建对象池
    pool := NewObjectPool()
    
    // 获取内存统计
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    // 使用对象池
    objects := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        obj := pool.Get()
        obj[0] = byte(i % 256)
        objects[i] = obj
    }
    
    // 归还对象
    for _, obj := range objects {
        pool.Put(obj)
    }
    
    // 获取最终统计
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    fmt.Printf("Memory before pool: %d KB\n", m1.Alloc/1024)
    fmt.Printf("Memory after pool: %d KB\n", m2.Alloc/1024)
    fmt.Printf("Memory difference: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
}
```

### 内存池优化

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// 优化的对象池
type OptimizedPool struct {
    pools []sync.Pool
    sizes []int
}

func NewOptimizedPool() *OptimizedPool {
    sizes := []int{64, 256, 1024, 4096, 16384}
    pools := make([]sync.Pool, len(sizes))
    
    for i, size := range sizes {
        size := size // 捕获循环变量
        pools[i] = sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        }
    }
    
    return &OptimizedPool{
        pools: pools,
        sizes: sizes,
    }
}

func (p *OptimizedPool) Get(size int) []byte {
    for i, poolSize := range p.sizes {
        if size <= poolSize {
            return p.pools[i].Get().([]byte)
        }
    }
    // 如果大小超过所有池，直接分配
    return make([]byte, size)
}

func (p *OptimizedPool) Put(obj []byte) {
    size := len(obj)
    for i, poolSize := range p.sizes {
        if size == poolSize {
            // 清空对象
            for j := range obj {
                obj[j] = 0
            }
            p.pools[i].Put(obj)
            return
        }
    }
    // 如果大小不匹配，丢弃对象
}

func main() {
    fmt.Println("=== Optimized Pool ===")
    
    pool := NewOptimizedPool()
    
    // 测试不同大小的对象
    sizes := []int{64, 256, 1024, 4096, 16384}
    
    for _, size := range sizes {
        start := time.Now()
        
        // 分配和释放对象
        for i := 0; i < 1000; i++ {
            obj := pool.Get(size)
            obj[0] = byte(i % 256)
            pool.Put(obj)
        }
        
        duration := time.Since(start)
        fmt.Printf("Size %d: %v\n", size, duration)
    }
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
        
        time.Sleep(1 * time.Millisecond)
    }
    
    // 运行一段时间
    time.Sleep(10 * time.Second)
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

## 性能优化

### 内存分配优化

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func memoryAllocationOptimization() {
    fmt.Println("=== Memory Allocation Optimization ===")
    
    // 1. 预分配切片
    start := time.Now()
    
    // 不好的做法
    var badSlice []int
    for i := 0; i < 100000; i++ {
        badSlice = append(badSlice, i)
    }
    badTime := time.Since(start)
    
    // 好的做法
    start = time.Now()
    goodSlice := make([]int, 0, 100000)
    for i := 0; i < 100000; i++ {
        goodSlice = append(goodSlice, i)
    }
    goodTime := time.Since(start)
    
    fmt.Printf("Bad approach: %v\n", badTime)
    fmt.Printf("Good approach: %v\n", goodTime)
    fmt.Printf("Improvement: %.2fx\n", float64(badTime)/float64(goodTime))
}

func main() {
    memoryAllocationOptimization()
}
```

### 内存复用优化

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func memoryReuseOptimization() {
    fmt.Println("=== Memory Reuse Optimization ===")
    
    // 使用对象池
    pool := &sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }
    
    start := time.Now()
    
    // 使用对象池
    for i := 0; i < 100000; i++ {
        obj := pool.Get().([]byte)
        obj[0] = byte(i % 256)
        pool.Put(obj)
    }
    
    poolTime := time.Since(start)
    
    // 不使用对象池
    start = time.Now()
    
    for i := 0; i < 100000; i++ {
        obj := make([]byte, 1024)
        obj[0] = byte(i % 256)
    }
    
    noPoolTime := time.Since(start)
    
    fmt.Printf("With pool: %v\n", poolTime)
    fmt.Printf("Without pool: %v\n", noPoolTime)
    fmt.Printf("Improvement: %.2fx\n", float64(noPoolTime)/float64(poolTime))
}

func main() {
    memoryReuseOptimization()
}
```

### 内存对齐优化

```go
package main

import (
    "fmt"
    "unsafe"
)

// 未对齐的结构体
type UnalignedStruct struct {
    a bool    // 1 byte
    b int64   // 8 bytes
    c bool    // 1 byte
    d int32   // 4 bytes
    e bool    // 1 byte
}

// 对齐的结构体
type AlignedStruct struct {
    b int64   // 8 bytes
    d int32   // 4 bytes
    a bool    // 1 byte
    c bool    // 1 byte
    e bool    // 1 byte
}

func main() {
    fmt.Println("=== Memory Alignment Optimization ===")
    
    unaligned := UnalignedStruct{}
    aligned := AlignedStruct{}
    
    fmt.Printf("Unaligned struct size: %d bytes\n", unsafe.Sizeof(unaligned))
    fmt.Printf("Aligned struct size: %d bytes\n", unsafe.Sizeof(aligned))
    fmt.Printf("Memory saved: %d bytes\n", unsafe.Sizeof(unaligned)-unsafe.Sizeof(aligned))
}
```

## 调试和监控

### 内存分析

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/pprof"
    "os"
)

func memoryAnalysis() {
    fmt.Println("=== Memory Analysis ===")
    
    // 创建内存分析文件
    f, err := os.Create("mem.prof")
    if err != nil {
        fmt.Printf("Error creating profile: %v\n", err)
        return
    }
    defer f.Close()
    
    // 写入内存分析
    pprof.WriteHeapProfile(f)
    
    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Memory profile written to mem.prof\n")
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("Sys: %d KB\n", m.Sys/1024)
    fmt.Printf("HeapObjects: %d\n", m.HeapObjects)
}

func main() {
    memoryAnalysis()
}
```

### 内存监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func memoryMonitoring() {
    fmt.Println("=== Memory Monitoring ===")
    
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
    memoryMonitoring()
}
```

## 最佳实践

### 1. 合理使用内存池

```go
package main

import (
    "fmt"
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

### 3. 监控内存使用

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 定期监控内存
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            if m.Alloc > 100*1024*1024 { // 100MB
                fmt.Printf("High memory usage: %d MB\n", m.Alloc/1024/1024)
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

import "fmt"

func main() {
    // 预分配切片
    slice := make([]int, 0, 1000)
    
    // 使用切片
    for i := 0; i < 1000; i++ {
        slice = append(slice, i)
    }
    
    fmt.Printf("Slice length: %d\n", len(slice))
}
```

## 总结

Go 的内存管理是一个复杂而高效的系统：

1. **内存分配器**: 分层分配策略，小对象和大对象分别处理
2. **垃圾回收器**: 三色标记算法，并发回收
3. **内存布局**: 堆、栈、元数据分离管理
4. **内存池**: 对象复用，减少分配开销
5. **性能优化**: 对齐、预分配、复用等策略

**关键要点**:
- 理解内存分配和回收机制
- 使用对象池减少分配开销
- 监控内存使用，避免泄漏
- 优化内存布局，提高缓存效率
- 合理设置GC参数

掌握 Go 的内存管理对于编写高性能的 Go 程序至关重要。
