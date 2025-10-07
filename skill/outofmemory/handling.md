# Go 内存溢出处理策略

## 🛠️ 优雅降级

### 1. 内存限制管理器

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

type MemoryManager struct {
    maxMemory    uint64
    currentUsage uint64
    mutex        sync.RWMutex
}

func NewMemoryManager(maxMemoryMB uint64) *MemoryManager {
    return &MemoryManager{
        maxMemory: maxMemoryMB * 1024 * 1024, // 转换为字节
    }
}

func (mm *MemoryManager) CheckMemory() bool {
    mm.mutex.RLock()
    defer mm.mutex.RUnlock()
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return m.Alloc < mm.maxMemory
}

func (mm *MemoryManager) AllocateMemory(size uint64) error {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    if m.Alloc+size > mm.maxMemory {
        return fmt.Errorf("内存不足，无法分配 %d 字节", size)
    }
    
    mm.currentUsage = m.Alloc
    return nil
}

func (mm *MemoryManager) ProcessWithMemoryLimit(data []byte) error {
    if !mm.CheckMemory() {
        return fmt.Errorf("内存使用率过高，拒绝处理")
    }
    
    // 模拟处理数据
    result := make([]byte, len(data)*2)
    copy(result, data)
    copy(result[len(data):], data)
    
    return nil
}
```

### 2. 内存使用监控

```go
func main() {
    mm := NewMemoryManager(100) // 100MB 限制
    
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024*1024) // 1MB
        
        err := mm.ProcessWithMemoryLimit(data)
        if err != nil {
            fmt.Printf("处理失败: %v\n", err)
            break
        }
        
        fmt.Printf("成功处理第 %d 个数据块\n", i+1)
        time.Sleep(10 * time.Millisecond)
    }
}
```

## 🏊 内存池管理

### 1. 字节池

```go
type BytePool struct {
    pool sync.Pool
    size int
}

func NewBytePool(size int) *BytePool {
    return &BytePool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        },
        size: size,
    }
}

func (bp *BytePool) Get() []byte {
    return bp.pool.Get().([]byte)
}

func (bp *BytePool) Put(b []byte) {
    if len(b) == bp.size {
        bp.pool.Put(b)
    }
}

// 使用内存池
func main() {
    pool := NewBytePool(1024 * 1024) // 1MB 池
    
    for i := 0; i < 1000; i++ {
        // 从池中获取
        data := pool.Get()
        
        // 使用数据
        for j := range data {
            data[j] = byte(i % 256)
        }
        
        // 处理数据
        fmt.Printf("处理数据块 %d\n", i+1)
        
        // 归还到池中
        pool.Put(data)
    }
}
```

### 2. 对象池

```go
type ObjectPool struct {
    pool sync.Pool
}

func NewObjectPool() *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &struct {
                    ID   int
                    Data []byte
                }{}
            },
        },
    }
}

func (op *ObjectPool) Get() *struct {
    ID   int
    Data []byte
} {
    return op.pool.Get().(*struct {
        ID   int
        Data []byte
    })
}

func (op *ObjectPool) Put(obj *struct {
    ID   int
    Data []byte
}) {
    // 重置对象
    obj.ID = 0
    obj.Data = obj.Data[:0]
    op.pool.Put(obj)
}
```

### 3. 切片池

```go
type SlicePool struct {
    pool sync.Pool
    size int
}

func NewSlicePool(size int) *SlicePool {
    return &SlicePool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]int, 0, size)
            },
        },
        size: size,
    }
}

func (sp *SlicePool) Get() []int {
    return sp.pool.Get().([]int)
}

func (sp *SlicePool) Put(s []int) {
    if cap(s) == sp.size {
        s = s[:0]
        sp.pool.Put(s)
    }
}
```

## 🌊 流式处理

### 1. 流式处理器

```go
type StreamProcessor struct {
    bufferSize int
}

func NewStreamProcessor(bufferSize int) *StreamProcessor {
    return &StreamProcessor{
        bufferSize: bufferSize,
    }
}

func (sp *StreamProcessor) ProcessStream(reader io.Reader) error {
    buffer := make([]byte, sp.bufferSize)
    
    for {
        n, err := reader.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        
        // 处理数据块
        sp.processChunk(buffer[:n])
    }
    
    return nil
}

func (sp *StreamProcessor) processChunk(data []byte) {
    // 模拟处理
    fmt.Printf("处理数据块，大小: %d 字节\n", len(data))
}

func main() {
    // 创建大量数据
    data := strings.Repeat("Hello, World! ", 1000000)
    reader := strings.NewReader(data)
    
    // 使用流式处理
    processor := NewStreamProcessor(1024 * 1024) // 1MB 缓冲区
    err := processor.ProcessStream(reader)
    if err != nil {
        fmt.Printf("处理失败: %v\n", err)
    }
}
```

### 2. 批量处理

```go
func batchProcessing() {
    const batchSize = 1000
    data := make([]int, 0, batchSize)
    
    for i := 0; i < 10000; i++ {
        data = append(data, i)
        
        // 批量处理
        if len(data) == batchSize {
            processBatch(data)
            data = data[:0] // 重置切片
        }
    }
    
    // 处理剩余数据
    if len(data) > 0 {
        processBatch(data)
    }
}

func processBatch(data []int) {
    // 模拟处理
    sum := 0
    for _, v := range data {
        sum += v
    }
    fmt.Printf("处理批次，大小: %d, 和: %d\n", len(data), sum)
}
```

## 🔄 内存回收策略

### 1. 及时释放资源

```go
func resourceManagement() {
    data := make([]byte, 1024*1024)
    defer func() {
        data = nil
        runtime.GC()
    }()
    
    // 使用 data
    for i := range data {
        data[i] = byte(i % 256)
    }
}
```

### 2. 定期清理

```go
func periodicCleanup() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for range ticker.C {
            // 清理缓存
            cleanupCache()
            
            // 强制垃圾回收
            runtime.GC()
            
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            fmt.Printf("清理后内存使用: %d KB\n", m.Alloc/1024)
        }
    }()
}

func cleanupCache() {
    // 清理缓存的逻辑
    fmt.Println("执行缓存清理")
}
```

### 3. 内存使用监控

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

## 🚨 紧急处理

### 1. 内存不足时的处理

```go
func handleMemoryShortage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 检查内存使用率
    usagePercent := float64(m.Alloc) / float64(m.Sys) * 100
    
    if usagePercent > 90 {
        fmt.Println("紧急情况: 内存使用率超过 90%")
        
        // 1. 立即释放非关键资源
        releaseNonCriticalResources()
        
        // 2. 强制垃圾回收
        runtime.GC()
        
        // 3. 如果仍然不足，拒绝新请求
        if !checkMemoryAfterCleanup() {
            fmt.Println("内存仍然不足，拒绝新请求")
            return
        }
    }
}

func releaseNonCriticalResources() {
    // 释放非关键资源的逻辑
    fmt.Println("释放非关键资源")
}

func checkMemoryAfterCleanup() bool {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    usagePercent := float64(m.Alloc) / float64(m.Sys) * 100
    return usagePercent < 80
}
```

### 2. 降级服务

```go
type ServiceDegrader struct {
    normalMode    bool
    degradedMode  bool
    memoryLimit   uint64
}

func NewServiceDegrader(memoryLimitMB uint64) *ServiceDegrader {
    return &ServiceDegrader{
        normalMode:   true,
        degradedMode: false,
        memoryLimit:  memoryLimitMB * 1024 * 1024,
    }
}

func (sd *ServiceDegrader) CheckAndDegrade() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    if m.Alloc > sd.memoryLimit {
        if sd.normalMode {
            fmt.Println("切换到降级模式")
            sd.normalMode = false
            sd.degradedMode = true
        }
    } else {
        if sd.degradedMode {
            fmt.Println("恢复正常模式")
            sd.normalMode = true
            sd.degradedMode = false
        }
    }
}

func (sd *ServiceDegrader) ProcessRequest(data []byte) error {
    if sd.degradedMode {
        // 降级模式：只处理小数据
        if len(data) > 1024 {
            return fmt.Errorf("降级模式：拒绝大数据请求")
        }
    }
    
    // 正常处理
    return nil
}
```

## 📊 内存使用报告

### 1. 内存使用统计

```go
type MemoryStats struct {
    Alloc      uint64
    TotalAlloc uint64
    NumGC      uint32
    Timestamp  time.Time
}

func (ms *MemoryStats) String() string {
    return fmt.Sprintf("内存: %d KB, 累计: %d KB, GC: %d, 时间: %s",
        ms.Alloc/1024, ms.TotalAlloc/1024, ms.NumGC, ms.Timestamp.Format("15:04:05"))
}
```

### 2. 内存使用历史记录

```go
type MemoryHistory struct {
    stats []MemoryStats
    mutex sync.RWMutex
}

func NewMemoryHistory() *MemoryHistory {
    return &MemoryHistory{
        stats: make([]MemoryStats, 0),
    }
}

func (mh *MemoryHistory) Record() {
    mh.mutex.Lock()
    defer mh.mutex.Unlock()
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    stat := MemoryStats{
        Alloc:      m.Alloc,
        TotalAlloc: m.TotalAlloc,
        NumGC:      m.NumGC,
        Timestamp:  time.Now(),
    }
    
    mh.stats = append(mh.stats, stat)
    
    // 只保留最近100条记录
    if len(mh.stats) > 100 {
        mh.stats = mh.stats[1:]
    }
}
```

### 3. 内存使用分析

```go
type MemoryAnalyzer struct {
    history *MemoryHistory
}

func NewMemoryAnalyzer() *MemoryAnalyzer {
    return &MemoryAnalyzer{
        history: NewMemoryHistory(),
    }
}

func (ma *MemoryAnalyzer) Analyze() {
    stats := ma.history.GetStats()
    
    if len(stats) < 2 {
        fmt.Println("数据不足，无法分析")
        return
    }
    
    // 计算内存使用趋势
    first := stats[0]
    last := stats[len(stats)-1]
    
    growth := last.Alloc - first.Alloc
    timeDiff := last.Timestamp.Sub(first.Timestamp)
    
    fmt.Printf("内存使用趋势分析:\n")
    fmt.Printf("  初始内存: %d KB\n", first.Alloc/1024)
    fmt.Printf("  当前内存: %d KB\n", last.Alloc/1024)
    fmt.Printf("  内存增长: %d KB\n", growth/1024)
    fmt.Printf("  时间跨度: %v\n", timeDiff)
    fmt.Printf("  平均增长: %d KB/s\n", growth/uint64(timeDiff.Seconds())/1024)
}
```

## 🔧 处理策略总结

1. **预防为主**: 设置内存限制，监控内存使用
2. **优雅降级**: 内存不足时拒绝非关键请求
3. **资源池化**: 使用对象池减少内存分配
4. **流式处理**: 处理大数据时使用流式方式
5. **及时释放**: 使用 defer 确保资源释放
6. **定期清理**: 定期清理缓存和临时数据
7. **紧急处理**: 内存不足时的紧急处理策略
8. **监控分析**: 持续监控和分析内存使用情况
