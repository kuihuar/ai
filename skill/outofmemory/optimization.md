# Go 内存优化最佳实践

## 1. GC 调优

### 1.1 设置合适的 GC 目标百分比

```go
// 设置 GC 目标百分比为 50%
func main() {
    debug.SetGCPercent(50)
    
    // 你的程序逻辑
    runApplication()
}

// 动态调整 GC 目标百分比
func adjustGCPercent() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 根据内存使用情况动态调整
    if m.Alloc > 100*1024*1024 { // 100MB
        debug.SetGCPercent(30) // 更频繁的GC
    } else {
        debug.SetGCPercent(50) // 正常GC
    }
}
```

### 1.2 手动触发 GC

```go
func triggerGC() {
    // 在关键时刻手动触发GC
    runtime.GC()
    
    // 等待GC完成
    runtime.Gosched()
}

// 在批量处理完成后触发GC
func processBatch(items []Item) {
    for _, item := range items {
        processItem(item)
    }
    
    // 处理完一批后触发GC
    runtime.GC()
}
```

### 1.3 使用 GOGC 环境变量

```bash
# 设置 GC 目标百分比
export GOGC=50

# 禁用 GC（仅用于测试）
export GOGC=off

# 运行程序
go run main.go
```

## 2. 内存池化技术

### 2.1 字节切片池

```go
// 字节切片池
type BytePool struct {
    pool sync.Pool
    size int
}

func NewBytePool(size int) *BytePool {
    return &BytePool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, size)
            },
        },
        size: size,
    }
}

func (p *BytePool) Get() []byte {
    return p.pool.Get().([]byte)
}

func (p *BytePool) Put(buf []byte) {
    if cap(buf) >= p.size {
        buf = buf[:0] // 重置长度
        p.pool.Put(buf)
    }
}

// 使用示例
var bytePool = NewBytePool(1024)

func processData(data []byte) []byte {
    buf := bytePool.Get()
    defer bytePool.Put(buf)
    
    // 使用 buf 处理数据
    buf = append(buf, "processed: "...)
    buf = append(buf, data...)
    
    result := make([]byte, len(buf))
    copy(result, buf)
    return result
}
```

### 2.2 结构体池

```go
// 结构体池
type Person struct {
    Name string
    Age  int
    Data []byte
}

type PersonPool struct {
    pool sync.Pool
}

func NewPersonPool() *PersonPool {
    return &PersonPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &Person{
                    Data: make([]byte, 0, 1024),
                }
            },
        },
    }
}

func (p *PersonPool) Get() *Person {
    person := p.pool.Get().(*Person)
    person.Name = ""
    person.Age = 0
    person.Data = person.Data[:0] // 重置切片
    return person
}

func (p *PersonPool) Put(person *Person) {
    p.pool.Put(person)
}

// 使用示例
var personPool = NewPersonPool()

func createPerson(name string, age int) *Person {
    person := personPool.Get()
    person.Name = name
    person.Age = age
    return person
}

func releasePerson(person *Person) {
    personPool.Put(person)
}
```

### 2.3 连接池

```go
// 数据库连接池
type ConnectionPool struct {
    pool chan *sql.DB
    max  int
}

func NewConnectionPool(max int) *ConnectionPool {
    return &ConnectionPool{
        pool: make(chan *sql.DB, max),
        max:  max,
    }
}

func (p *ConnectionPool) Get() (*sql.DB, error) {
    select {
    case conn := <-p.pool:
        return conn, nil
    default:
        // 创建新连接
        return sql.Open("mysql", "user:password@/dbname")
    }
}

func (p *ConnectionPool) Put(conn *sql.DB) {
    select {
    case p.pool <- conn:
    default:
        // 池已满，关闭连接
        conn.Close()
    }
}
```

## 3. 内存映射文件

### 3.1 使用 mmap 处理大文件

```go
import (
    "os"
    "syscall"
    "unsafe"
)

// 内存映射文件
type MappedFile struct {
    data []byte
    file *os.File
}

func OpenMappedFile(filename string) (*MappedFile, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }
    
    size := stat.Size()
    if size == 0 {
        return &MappedFile{data: nil, file: file}, nil
    }
    
    // 映射文件到内存
    data, err := syscall.Mmap(int(file.Fd()), 0, int(size), 
        syscall.PROT_READ, syscall.MAP_SHARED)
    if err != nil {
        file.Close()
        return nil, err
    }
    
    return &MappedFile{data: data, file: file}, nil
}

func (mf *MappedFile) Close() error {
    if mf.data != nil {
        syscall.Munmap(mf.data)
    }
    return mf.file.Close()
}

func (mf *MappedFile) Data() []byte {
    return mf.data
}

// 使用示例
func processLargeFile(filename string) error {
    mf, err := OpenMappedFile(filename)
    if err != nil {
        return err
    }
    defer mf.Close()
    
    data := mf.Data()
    // 处理数据，不需要加载到内存
    for i := 0; i < len(data); i += 1024 {
        chunk := data[i:i+1024]
        processChunk(chunk)
    }
    
    return nil
}
```

## 4. 流式处理

### 4.1 流式读取大文件

```go
func processLargeFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    // 设置缓冲区大小
    buf := make([]byte, 0, 64*1024)
    scanner.Buffer(buf, 1024*1024)
    
    for scanner.Scan() {
        line := scanner.Text()
        processLine(line)
    }
    
    return scanner.Err()
}

// 分块处理
func processFileInChunks(filename string, chunkSize int) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    buffer := make([]byte, chunkSize)
    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        
        // 处理这一块数据
        processChunk(buffer[:n])
    }
    
    return nil
}
```

### 4.2 流式处理网络数据

```go
func processHTTPResponse(resp *http.Response) error {
    defer resp.Body.Close()
    
    // 使用流式读取
    reader := bufio.NewReader(resp.Body)
    buffer := make([]byte, 4096)
    
    for {
        n, err := reader.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        
        // 处理这一块数据
        processChunk(buffer[:n])
    }
    
    return nil
}
```

## 5. 内存压缩和优化

### 5.1 使用更紧凑的数据结构

```go
// ❌ 浪费内存的结构
type User struct {
    ID       int64
    Name     string
    Email    string
    Age      int
    IsActive bool
    Created  time.Time
}

// ✅ 更紧凑的结构
type User struct {
    ID       int32  // 使用 int32 而不是 int64
    Name     string
    Email    string
    Age      uint8  // 使用 uint8 而不是 int
    IsActive bool
    Created  int64  // 使用时间戳而不是 time.Time
}

// ✅ 使用位字段
type UserFlags struct {
    IsActive   bool
    IsVerified bool
    IsPremium  bool
    IsAdmin    bool
}

// 使用位操作
func (f *UserFlags) SetActive(active bool) {
    if active {
        f.flags |= 1 << 0
    } else {
        f.flags &^= 1 << 0
    }
}

func (f *UserFlags) IsActive() bool {
    return f.flags&(1<<0) != 0
}
```

### 5.2 字符串优化

```go
// 使用字符串常量池
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
    StatusPending  = "pending"
)

// 使用枚举而不是字符串
type Status int

const (
    StatusActive Status = iota
    StatusInactive
    StatusPending
)

func (s Status) String() string {
    switch s {
    case StatusActive:
        return "active"
    case StatusInactive:
        return "inactive"
    case StatusPending:
        return "pending"
    default:
        return "unknown"
    }
}

// 使用字符串构建器
func buildQuery(conditions []string) string {
    var builder strings.Builder
    builder.Grow(len(conditions) * 20) // 预分配容量
    
    builder.WriteString("SELECT * FROM table WHERE ")
    for i, condition := range conditions {
        if i > 0 {
            builder.WriteString(" AND ")
        }
        builder.WriteString(condition)
    }
    
    return builder.String()
}
```

## 6. 内存监控和调试

### 6.1 内存使用监控

```go
type MemoryMonitor struct {
    maxMemory    uint64
    checkInterval time.Duration
    stopChan     chan bool
}

func NewMemoryMonitor(maxMemory uint64, interval time.Duration) *MemoryMonitor {
    return &MemoryMonitor{
        maxMemory:    maxMemory,
        checkInterval: interval,
        stopChan:     make(chan bool),
    }
}

func (m *MemoryMonitor) Start() {
    ticker := time.NewTicker(m.checkInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            m.checkMemory()
        case <-m.stopChan:
            return
        }
    }
}

func (m *MemoryMonitor) Stop() {
    close(m.stopChan)
}

func (m *MemoryMonitor) checkMemory() {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)
    
    if stats.Alloc > m.maxMemory {
        log.Printf("内存使用超过限制: %d MB (限制: %d MB)", 
            stats.Alloc/1024/1024, m.maxMemory/1024/1024)
        
        // 触发GC
        runtime.GC()
        
        // 再次检查
        runtime.ReadMemStats(&stats)
        if stats.Alloc > m.maxMemory {
            log.Printf("GC后内存仍超限: %d MB", stats.Alloc/1024/1024)
        }
    }
}
```

### 6.2 内存分析

```go
func analyzeMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    log.Printf("内存统计:")
    log.Printf("  已分配: %d MB", m.Alloc/1024/1024)
    log.Printf("  总分配: %d MB", m.TotalAlloc/1024/1024)
    log.Printf("  系统内存: %d MB", m.Sys/1024/1024)
    log.Printf("  GC次数: %d", m.NumGC)
    log.Printf("  上次GC: %v", time.Unix(0, int64(m.LastGC)))
    
    if m.NumGC > 0 {
        avgGC := float64(m.PauseTotalNs) / float64(m.NumGC) / 1e6
        log.Printf("  平均GC暂停: %.2f ms", avgGC)
    }
}
```

## 7. 性能测试和基准测试

### 7.1 内存使用基准测试

```go
func BenchmarkMemoryUsage(b *testing.B) {
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        // 测试代码
        data := make([]byte, 1024)
        processData(data)
    }
}

func BenchmarkPoolUsage(b *testing.B) {
    pool := NewBytePool(1024)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        buf := pool.Get()
        processData(buf)
        pool.Put(buf)
    }
}
```

### 7.2 内存泄漏检测

```go
func TestMemoryLeak(t *testing.T) {
    var m1, m2 runtime.MemStats
    
    // 记录初始内存
    runtime.ReadMemStats(&m1)
    
    // 执行可能泄漏的操作
    for i := 0; i < 1000; i++ {
        data := make([]byte, 1024)
        processData(data)
    }
    
    // 强制GC
    runtime.GC()
    runtime.GC()
    
    // 记录GC后内存
    runtime.ReadMemStats(&m2)
    
    // 检查内存是否显著增加
    if m2.Alloc > m1.Alloc+1024*1024 { // 1MB
        t.Errorf("可能存在内存泄漏: 初始 %d, 现在 %d", 
            m1.Alloc, m2.Alloc)
    }
}
```

## 8. 最佳实践总结

1. **合理设置GC参数**：根据应用特点调整GOGC
2. **使用对象池**：重用对象减少GC压力
3. **流式处理**：避免一次性加载大量数据
4. **内存映射**：处理大文件时使用mmap
5. **紧凑数据结构**：使用合适的数据类型
6. **定期监控**：设置内存使用预警
7. **性能测试**：使用基准测试验证优化效果
8. **内存分析**：使用pprof等工具分析内存使用

通过遵循这些最佳实践，可以显著提高Go程序的内存使用效率。
