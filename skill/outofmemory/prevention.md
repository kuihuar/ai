# Go 内存溢出预防措施

## 1. 内存泄漏预防

### 1.1 及时关闭资源

```go
// ❌ 错误：忘记关闭文件
func readFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    // 忘记关闭文件
    return nil
}

// ✅ 正确：使用 defer 确保资源关闭
func readFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // 确保文件被关闭
    
    // 处理文件内容
    return nil
}

// ✅ 更好：使用匿名函数处理错误
func readFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer func() {
        if closeErr := file.Close(); closeErr != nil {
            log.Printf("关闭文件失败: %v", closeErr)
        }
    }()
    
    return nil
}
```

### 1.2 避免循环引用

```go
// ❌ 错误：循环引用导致内存泄漏
type Node struct {
    Value    string
    Parent   *Node
    Children []*Node
}

func createCircularReference() {
    parent := &Node{Value: "parent"}
    child := &Node{Value: "child", Parent: parent}
    parent.Children = append(parent.Children, child)
    // parent 和 child 相互引用，无法被GC回收
}

// ✅ 正确：使用弱引用或接口
type Node struct {
    Value    string
    Parent   interface{} // 使用接口避免强引用
    Children []*Node
}

// ✅ 或者使用 context 管理生命周期
type Node struct {
    Value    string
    Parent   *Node
    Children []*Node
    ctx      context.Context
    cancel   context.CancelFunc
}

func NewNode(value string, parent *Node) *Node {
    ctx, cancel := context.WithCancel(context.Background())
    return &Node{
        Value:  value,
        Parent: parent,
        ctx:    ctx,
        cancel: cancel,
    }
}

func (n *Node) Close() {
    n.cancel() // 取消context，帮助GC
}
```

### 1.3 及时清理定时器和协程

```go
// ❌ 错误：定时器没有停止
func startTimer() {
    ticker := time.NewTicker(1 * time.Second)
    go func() {
        for range ticker.C {
            // 处理定时任务
        }
    }()
    // 忘记停止ticker
}

// ✅ 正确：管理定时器生命周期
type Service struct {
    ticker *time.Ticker
    done   chan bool
}

func NewService() *Service {
    return &Service{
        ticker: time.NewTicker(1 * time.Second),
        done:   make(chan bool),
    }
}

func (s *Service) Start() {
    go func() {
        for {
            select {
            case <-s.ticker.C:
                // 处理定时任务
            case <-s.done:
                return
            }
        }
    }()
}

func (s *Service) Stop() {
    s.ticker.Stop()
    close(s.done)
}
```

## 2. 内存使用优化

### 2.1 对象池模式

```go
// 使用 sync.Pool 重用对象
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024) // 预分配容量
    },
}

func processData(data []byte) []byte {
    // 从池中获取buffer
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf[:0]) // 重置并放回池中
    
    // 使用buffer处理数据
    buf = append(buf, "processed: "...)
    buf = append(buf, data...)
    
    result := make([]byte, len(buf))
    copy(result, buf)
    return result
}

// 自定义对象池
type ObjectPool struct {
    pool chan interface{}
    new  func() interface{}
}

func NewObjectPool(size int, newFunc func() interface{}) *ObjectPool {
    return &ObjectPool{
        pool: make(chan interface{}, size),
        new:  newFunc,
    }
}

func (p *ObjectPool) Get() interface{} {
    select {
    case obj := <-p.pool:
        return obj
    default:
        return p.new()
    }
}

func (p *ObjectPool) Put(obj interface{}) {
    select {
    case p.pool <- obj:
    default:
        // 池已满，丢弃对象
    }
}
```

### 2.2 字符串优化

```go
// ❌ 错误：频繁字符串拼接
func buildString(parts []string) string {
    result := ""
    for _, part := range parts {
        result += part // 每次拼接都创建新字符串
    }
    return result
}

// ✅ 正确：使用 strings.Builder
func buildString(parts []string) string {
    var builder strings.Builder
    builder.Grow(len(parts) * 10) // 预分配容量
    
    for _, part := range parts {
        builder.WriteString(part)
    }
    return builder.String()
}

// ✅ 或者使用 bytes.Buffer
func buildString(parts []string) string {
    var buffer bytes.Buffer
    buffer.Grow(len(parts) * 10)
    
    for _, part := range parts {
        buffer.WriteString(part)
    }
    return buffer.String()
}
```

### 2.3 切片优化

```go
// ❌ 错误：频繁扩容
func processItems(items []int) []int {
    result := []int{} // 零值切片，每次append都可能扩容
    for _, item := range items {
        if item > 0 {
            result = append(result, item*2)
        }
    }
    return result
}

// ✅ 正确：预分配容量
func processItems(items []int) []int {
    result := make([]int, 0, len(items)) // 预分配容量
    for _, item := range items {
        if item > 0 {
            result = append(result, item*2)
        }
    }
    return result
}

// ✅ 更好：如果知道确切大小，直接分配
func processItems(items []int) []int {
    count := 0
    for _, item := range items {
        if item > 0 {
            count++
        }
    }
    
    result := make([]int, count)
    idx := 0
    for _, item := range items {
        if item > 0 {
            result[idx] = item * 2
            idx++
        }
    }
    return result
}
```

## 3. 内存分配策略

### 3.1 减少小对象分配

```go
// ❌ 错误：创建大量小对象
type Point struct {
    X, Y float64
}

func createPoints(count int) []Point {
    points := make([]Point, count)
    for i := 0; i < count; i++ {
        points[i] = Point{X: float64(i), Y: float64(i * 2)}
    }
    return points
}

// ✅ 正确：使用结构体切片减少分配
func createPoints(count int) []Point {
    points := make([]Point, count)
    for i := 0; i < count; i++ {
        points[i] = Point{X: float64(i), Y: float64(i * 2)}
    }
    return points
}

// ✅ 更好：如果不需要随机访问，使用数组
func createPointsArray(count int) [1000]Point {
    var points [1000]Point
    for i := 0; i < count && i < 1000; i++ {
        points[i] = Point{X: float64(i), Y: float64(i * 2)}
    }
    return points
}
```

### 3.2 避免不必要的内存拷贝

```go
// ❌ 错误：不必要的拷贝
func processData(data []byte) []byte {
    result := make([]byte, len(data))
    copy(result, data) // 不必要的拷贝
    
    // 处理数据
    for i := range result {
        result[i] = result[i] * 2
    }
    return result
}

// ✅ 正确：直接修改原数据
func processData(data []byte) []byte {
    // 直接修改原数据
    for i := range data {
        data[i] = data[i] * 2
    }
    return data
}

// ✅ 或者使用指针避免拷贝
func processData(data *[]byte) {
    for i := range *data {
        (*data)[i] = (*data)[i] * 2
    }
}
```

## 4. 并发安全的内存管理

### 4.1 使用原子操作

```go
// ❌ 错误：使用互斥锁保护简单操作
type Counter struct {
    mu    sync.Mutex
    value int64
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

func (c *Counter) Value() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// ✅ 正确：使用原子操作
type Counter struct {
    value int64
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Value() int64 {
    return atomic.LoadInt64(&c.value)
}
```

### 4.2 避免竞态条件

```go
// ❌ 错误：竞态条件
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, ok := c.data[key]
    return value, ok
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

// ✅ 正确：使用 sync.Map
type Cache struct {
    data sync.Map
}

func (c *Cache) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}

func (c *Cache) Set(key string, value interface{}) {
    c.data.Store(key, value)
}
```

## 5. 内存监控和预警

### 5.1 定期检查内存使用

```go
func monitorMemory() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        // 检查内存使用率
        if m.Alloc > 100*1024*1024 { // 100MB
            log.Printf("内存使用过高: %d MB", m.Alloc/1024/1024)
        }
        
        // 检查GC频率
        if m.NumGC > 0 {
            avgGC := float64(m.PauseTotalNs) / float64(m.NumGC) / 1e6
            if avgGC > 10 { // 10ms
                log.Printf("GC暂停时间过长: %.2f ms", avgGC)
            }
        }
    }
}
```

### 5.2 设置内存限制

```go
func setMemoryLimit() {
    // 设置最大内存使用量
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    maxMemory := 200 * 1024 * 1024 // 200MB
    if m.Alloc > maxMemory {
        log.Printf("内存使用超过限制: %d MB", m.Alloc/1024/1024)
        // 触发GC
        runtime.GC()
        // 或者退出程序
        os.Exit(1)
    }
}
```

## 6. 最佳实践总结

1. **及时释放资源**：使用 `defer` 确保资源被正确释放
2. **避免循环引用**：使用弱引用或 context 管理生命周期
3. **预分配容量**：为切片、字符串等预分配合适的容量
4. **使用对象池**：重用对象减少GC压力
5. **减少小对象分配**：批量处理数据
6. **使用原子操作**：避免不必要的锁竞争
7. **定期监控内存**：设置内存使用预警
8. **合理设置GC参数**：根据应用特点调整GC参数

通过遵循这些预防措施，可以大大降低Go程序出现内存溢出的风险。
