# Go 内存监控和工具使用

## 1. 内置内存监控

### 1.1 runtime.MemStats 使用

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
    
    fmt.Printf("内存统计:\n")
    fmt.Printf("  已分配内存: %d KB\n", m.Alloc/1024)
    fmt.Printf("  总分配内存: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("  系统内存: %d KB\n", m.Sys/1024)
    fmt.Printf("  堆内存: %d KB\n", m.HeapAlloc/1024)
    fmt.Printf("  堆系统内存: %d KB\n", m.HeapSys/1024)
    fmt.Printf("  堆空闲内存: %d KB\n", m.HeapIdle/1024)
    fmt.Printf("  堆使用内存: %d KB\n", m.HeapInuse/1024)
    fmt.Printf("  堆释放内存: %d KB\n", m.HeapReleased/1024)
    fmt.Printf("  GC次数: %d\n", m.NumGC)
    fmt.Printf("  上次GC时间: %v\n", time.Unix(0, int64(m.LastGC)))
    fmt.Printf("  下次GC目标: %d KB\n", m.NextGC/1024)
    fmt.Printf("  GC暂停时间: %d ns\n", m.PauseTotalNs)
    fmt.Printf("  平均GC暂停: %d ns\n", m.PauseTotalNs/uint64(m.NumGC))
    fmt.Printf("  对象数量: %d\n", m.Mallocs)
    fmt.Printf("  释放对象数量: %d\n", m.Frees)
    fmt.Printf("  存活对象数量: %d\n", m.Mallocs-m.Frees)
}

func main() {
    printMemStats()
    
    // 分配一些内存
    data := make([]byte, 1024*1024) // 1MB
    _ = data
    
    printMemStats()
    
    // 触发GC
    runtime.GC()
    printMemStats()
}
```

### 1.2 内存使用趋势监控

```go
type MemoryMonitor struct {
    samples []MemorySample
    maxSamples int
}

type MemorySample struct {
    Timestamp time.Time
    Alloc     uint64
    Sys       uint64
    NumGC     uint32
}

func NewMemoryMonitor(maxSamples int) *MemoryMonitor {
    return &MemoryMonitor{
        samples: make([]MemorySample, 0, maxSamples),
        maxSamples: maxSamples,
    }
}

func (m *MemoryMonitor) Sample() {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)
    
    sample := MemorySample{
        Timestamp: time.Now(),
        Alloc:     stats.Alloc,
        Sys:       stats.Sys,
        NumGC:     stats.NumGC,
    }
    
    m.samples = append(m.samples, sample)
    if len(m.samples) > m.maxSamples {
        m.samples = m.samples[1:]
    }
}

func (m *MemoryMonitor) GetTrend() (float64, float64) {
    if len(m.samples) < 2 {
        return 0, 0
    }
    
    first := m.samples[0]
    last := m.samples[len(m.samples)-1]
    
    allocTrend := float64(last.Alloc-first.Alloc) / float64(len(m.samples))
    gcTrend := float64(last.NumGC-first.NumGC) / float64(len(m.samples))
    
    return allocTrend, gcTrend
}

func (m *MemoryMonitor) Start(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            m.Sample()
        }
    }()
}
```

## 2. pprof 内存分析

### 2.1 启用 pprof

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "runtime"
)

func main() {
    // 启用 pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 你的程序逻辑
    runApplication()
}

// 在代码中手动触发内存分析
func triggerMemoryProfile() {
    f, err := os.Create("mem.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    runtime.GC() // 触发GC
    if err := pprof.WriteHeapProfile(f); err != nil {
        log.Fatal(err)
    }
}
```

### 2.2 使用 pprof 命令行工具

```bash
# 查看内存使用情况
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看内存分配情况
go tool pprof http://localhost:6060/debug/pprof/allocs

# 查看内存使用趋势
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# 比较两个内存快照
go tool pprof -base=old.prof new.prof

# 查看内存使用最多的函数
go tool pprof -top http://localhost:6060/debug/pprof/heap

# 查看内存使用最多的代码行
go tool pprof -list=main http://localhost:6060/debug/pprof/heap
```

### 2.3 内存分析示例

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "time"
)

func main() {
    // 启用 pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 模拟内存使用
    go func() {
        for {
            data := make([]byte, 1024*1024) // 1MB
            _ = data
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // 模拟内存泄漏
    go func() {
        var leak []byte
        for {
            leak = append(leak, make([]byte, 1024)...)
            time.Sleep(50 * time.Millisecond)
        }
    }()
    
    select {}
}
```

## 3. 内存泄漏检测

### 3.1 使用 goleak 检测协程泄漏

```go
package main

import (
    "go.uber.org/goleak"
    "testing"
)

func TestNoGoroutineLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    
    // 你的测试代码
    runTest()
}

// 检测特定协程泄漏
func TestSpecificGoroutineLeak(t *testing.T) {
    defer goleak.VerifyNone(t, 
        goleak.IgnoreTopFunction("runtime.gopark"),
        goleak.IgnoreTopFunction("runtime.goparkunlock"),
    )
    
    runTest()
}
```

### 3.2 内存泄漏检测工具

```go
type LeakDetector struct {
    initialAlloc uint64
    threshold    uint64
}

func NewLeakDetector(threshold uint64) *LeakDetector {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return &LeakDetector{
        initialAlloc: m.Alloc,
        threshold:    threshold,
    }
}

func (ld *LeakDetector) Check() bool {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return m.Alloc > ld.initialAlloc+ld.threshold
}

func (ld *LeakDetector) GetLeakSize() uint64 {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    if m.Alloc > ld.initialAlloc {
        return m.Alloc - ld.initialAlloc
    }
    return 0
}

// 使用示例
func main() {
    detector := NewLeakDetector(1024 * 1024) // 1MB阈值
    
    // 运行程序
    runApplication()
    
    // 检查内存泄漏
    if detector.Check() {
        log.Printf("检测到内存泄漏: %d KB", detector.GetLeakSize()/1024)
    }
}
```

## 4. 内存使用预警

### 4.1 内存使用率监控

```go
type MemoryAlert struct {
    maxMemory    uint64
    alertChannel chan string
}

func NewMemoryAlert(maxMemory uint64) *MemoryAlert {
    return &MemoryAlert{
        maxMemory:    maxMemory,
        alertChannel: make(chan string, 10),
    }
}

func (ma *MemoryAlert) Start() {
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            ma.checkMemory()
        }
    }()
}

func (ma *MemoryAlert) checkMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    if m.Alloc > ma.maxMemory {
        alert := fmt.Sprintf("内存使用超过限制: %d MB (限制: %d MB)", 
            m.Alloc/1024/1024, ma.maxMemory/1024/1024)
        
        select {
        case ma.alertChannel <- alert:
        default:
            // 通道已满，丢弃警告
        }
    }
}

func (ma *MemoryAlert) GetAlert() <-chan string {
    return ma.alertChannel
}
```

### 4.2 内存使用趋势分析

```go
type MemoryTrendAnalyzer struct {
    samples []MemorySample
    window  int
}

func NewMemoryTrendAnalyzer(window int) *MemoryTrendAnalyzer {
    return &MemoryTrendAnalyzer{
        samples: make([]MemorySample, 0, window),
        window:  window,
    }
}

func (mta *MemoryTrendAnalyzer) AddSample() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    sample := MemorySample{
        Timestamp: time.Now(),
        Alloc:     m.Alloc,
        Sys:       m.Sys,
        NumGC:     m.NumGC,
    }
    
    mta.samples = append(mta.samples, sample)
    if len(mta.samples) > mta.window {
        mta.samples = mta.samples[1:]
    }
}

func (mta *MemoryTrendAnalyzer) GetTrend() (float64, bool) {
    if len(mta.samples) < 2 {
        return 0, false
    }
    
    first := mta.samples[0]
    last := mta.samples[len(mta.samples)-1]
    
    // 计算内存增长趋势
    trend := float64(last.Alloc-first.Alloc) / float64(len(mta.samples))
    
    // 判断是否持续增长
    growing := true
    for i := 1; i < len(mta.samples); i++ {
        if mta.samples[i].Alloc <= mta.samples[i-1].Alloc {
            growing = false
            break
        }
    }
    
    return trend, growing
}
```

## 5. 内存优化建议

### 5.1 自动内存优化

```go
type MemoryOptimizer struct {
    maxMemory    uint64
    gcThreshold  uint64
    lastGC       time.Time
}

func NewMemoryOptimizer(maxMemory uint64) *MemoryOptimizer {
    return &MemoryOptimizer{
        maxMemory:   maxMemory,
        gcThreshold: maxMemory * 80 / 100, // 80%阈值
    }
}

func (mo *MemoryOptimizer) CheckAndOptimize() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 检查是否需要GC
    if m.Alloc > mo.gcThreshold {
        runtime.GC()
        mo.lastGC = time.Now()
    }
    
    // 检查是否需要强制GC
    if m.Alloc > mo.maxMemory {
        runtime.GC()
        runtime.GC() // 连续两次GC
        
        // 如果仍然超限，记录警告
        runtime.ReadMemStats(&m)
        if m.Alloc > mo.maxMemory {
            log.Printf("内存使用仍然超限: %d MB", m.Alloc/1024/1024)
        }
    }
}
```

### 5.2 内存使用报告

```go
func generateMemoryReport() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    report := fmt.Sprintf(`
内存使用报告:
=============
当前内存使用: %d MB
总分配内存: %d MB
系统内存: %d MB
堆内存: %d MB
堆空闲内存: %d MB
堆使用内存: %d MB
GC次数: %d
上次GC: %v
GC暂停时间: %d ns
平均GC暂停: %d ns
对象数量: %d
释放对象数量: %d
存活对象数量: %d
`,
        m.Alloc/1024/1024,
        m.TotalAlloc/1024/1024,
        m.Sys/1024/1024,
        m.HeapAlloc/1024/1024,
        m.HeapIdle/1024/1024,
        m.HeapInuse/1024/1024,
        m.NumGC,
        time.Unix(0, int64(m.LastGC)),
        m.PauseTotalNs,
        m.PauseTotalNs/uint64(m.NumGC),
        m.Mallocs,
        m.Frees,
        m.Mallocs-m.Frees,
    )
    
    log.Println(report)
}
```

## 6. 生产环境监控

### 6.1 集成监控系统

```go
type PrometheusExporter struct {
    memoryGauge prometheus.Gauge
    gcCounter   prometheus.Counter
}

func NewPrometheusExporter() *PrometheusExporter {
    return &PrometheusExporter{
        memoryGauge: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "go_memory_alloc_bytes",
            Help: "当前分配的内存字节数",
        }),
        gcCounter: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "go_gc_total",
            Help: "GC总次数",
        }),
    }
}

func (pe *PrometheusExporter) Update() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    pe.memoryGauge.Set(float64(m.Alloc))
    pe.gcCounter.Add(float64(m.NumGC))
}

func (pe *PrometheusExporter) Register() {
    prometheus.MustRegister(pe.memoryGauge)
    prometheus.MustRegister(pe.gcCounter)
}
```

### 6.2 日志记录

```go
func logMemoryUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    log.WithFields(log.Fields{
        "alloc_mb":      m.Alloc / 1024 / 1024,
        "sys_mb":        m.Sys / 1024 / 1024,
        "heap_mb":       m.HeapAlloc / 1024 / 1024,
        "gc_count":      m.NumGC,
        "last_gc":       time.Unix(0, int64(m.LastGC)),
        "pause_total_ns": m.PauseTotalNs,
    }).Info("内存使用情况")
}
```

## 7. 最佳实践总结

1. **定期监控**：使用 runtime.MemStats 定期检查内存使用
2. **使用 pprof**：使用 pprof 工具分析内存使用和泄漏
3. **设置预警**：设置内存使用阈值和预警机制
4. **趋势分析**：分析内存使用趋势，提前发现问题
5. **自动优化**：根据内存使用情况自动触发优化
6. **生产监控**：集成到监控系统中，实时监控
7. **日志记录**：记录内存使用情况，便于问题排查
8. **性能测试**：使用基准测试验证内存优化效果

通过合理使用这些监控工具和技巧，可以及时发现和解决Go程序中的内存问题。
