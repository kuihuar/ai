# Go CPU 性能优化策略

## 1. 算法优化

### 1.1 选择合适的数据结构

```go
// ❌ 错误：使用低效的数据结构
func inefficientSearch(data []int, target int) bool {
    for _, v := range data {
        if v == target {
            return true
        }
    }
    return false
}

// ✅ 正确：使用 map 进行 O(1) 查找
func efficientSearch(data []int, target int) bool {
    lookup := make(map[int]bool)
    for _, v := range data {
        lookup[v] = true
    }
    return lookup[target]
}

// ✅ 更好：如果数据已排序，使用二分查找
func binarySearch(data []int, target int) bool {
    left, right := 0, len(data)-1
    for left <= right {
        mid := (left + right) / 2
        if data[mid] == target {
            return true
        } else if data[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return false
}
```

### 1.2 避免重复计算

```go
// ❌ 错误：重复计算
func inefficientCalculation(data []int) []int {
    result := make([]int, len(data))
    for i, v := range data {
        // 每次都重新计算
        result[i] = v*v + 2*v + 1
    }
    return result
}

// ✅ 正确：使用缓存避免重复计算
var calculationCache = make(map[int]int)

func efficientCalculation(data []int) []int {
    result := make([]int, len(data))
    for i, v := range data {
        if cached, exists := calculationCache[v]; exists {
            result[i] = cached
        } else {
            calculated := v*v + 2*v + 1
            calculationCache[v] = calculated
            result[i] = calculated
        }
    }
    return result
}

// ✅ 更好：预计算常用值
func precomputedCalculation(data []int) []int {
    result := make([]int, len(data))
    for i, v := range data {
        result[i] = v*v + 2*v + 1
    }
    return result
}
```

### 1.3 优化循环

```go
// ❌ 错误：嵌套循环
func inefficientNestedLoop(matrix [][]int) int {
    sum := 0
    for i := 0; i < len(matrix); i++ {
        for j := 0; j < len(matrix[i]); j++ {
            sum += matrix[i][j]
        }
    }
    return sum
}

// ✅ 正确：优化循环顺序
func optimizedNestedLoop(matrix [][]int) int {
    sum := 0
    for i := 0; i < len(matrix); i++ {
        row := matrix[i]
        for j := 0; j < len(row); j++ {
            sum += row[j]
        }
    }
    return sum
}

// ✅ 更好：使用 range 循环
func rangeLoop(matrix [][]int) int {
    sum := 0
    for _, row := range matrix {
        for _, val := range row {
            sum += val
        }
    }
    return sum
}
```

## 2. 并发优化

### 2.1 合理使用 goroutine

```go
// ❌ 错误：过度使用 goroutine
func inefficientConcurrency(data []int) []int {
    results := make([]int, len(data))
    var wg sync.WaitGroup
    
    for i, v := range data {
        wg.Add(1)
        go func(i, v int) {
            defer wg.Done()
            results[i] = v * v // 简单计算不需要并发
        }(i, v)
    }
    
    wg.Wait()
    return results
}

// ✅ 正确：合理使用并发
func efficientConcurrency(data []int) []int {
    const numWorkers = runtime.NumCPU()
    results := make([]int, len(data))
    
    var wg sync.WaitGroup
    chunkSize := len(data) / numWorkers
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            for j := start; j < end; j++ {
                results[j] = data[j] * data[j]
            }
        }(i*chunkSize, (i+1)*chunkSize)
    }
    
    wg.Wait()
    return results
}
```

### 2.2 使用 worker pool 模式

```go
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultChan chan Result
    wg         sync.WaitGroup
}

type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID int
    Data  interface{}
    Error error
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan Job, workers*2),
        resultChan: make(chan Result, workers*2),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    
    for job := range wp.jobQueue {
        result := wp.processJob(job)
        wp.resultChan <- result
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // 处理任务
    time.Sleep(100 * time.Millisecond) // 模拟处理时间
    return Result{
        JobID: job.ID,
        Data:  job.Data,
        Error: nil,
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Close() {
    close(wp.jobQueue)
    wp.wg.Wait()
    close(wp.resultChan)
}
```

### 2.3 避免锁竞争

```go
// ❌ 错误：使用全局锁
type InefficientCounter struct {
    mu    sync.Mutex
    count int
}

func (c *InefficientCounter) Increment() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

func (c *InefficientCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

// ✅ 正确：使用原子操作
type EfficientCounter struct {
    count int64
}

func (c *EfficientCounter) Increment() {
    atomic.AddInt64(&c.count, 1)
}

func (c *EfficientCounter) Value() int64 {
    return atomic.LoadInt64(&c.count)
}

// ✅ 更好：使用分片减少竞争
type ShardedCounter struct {
    shards []*CounterShard
    mask   uint64
}

type CounterShard struct {
    count int64
}

func NewShardedCounter(shards int) *ShardedCounter {
    sc := &ShardedCounter{
        shards: make([]*CounterShard, shards),
        mask:   uint64(shards - 1),
    }
    
    for i := 0; i < shards; i++ {
        sc.shards[i] = &CounterShard{}
    }
    
    return sc
}

func (sc *ShardedCounter) Increment() {
    shard := sc.shards[rand.Uint64()&sc.mask]
    atomic.AddInt64(&shard.count, 1)
}

func (sc *ShardedCounter) Value() int64 {
    var total int64
    for _, shard := range sc.shards {
        total += atomic.LoadInt64(&shard.count)
    }
    return total
}
```

## 3. 内存优化

### 3.1 减少内存分配

```go
// ❌ 错误：频繁分配内存
func inefficientMemory(data []int) []int {
    result := []int{}
    for _, v := range data {
        if v > 0 {
            result = append(result, v*2)
        }
    }
    return result
}

// ✅ 正确：预分配容量
func efficientMemory(data []int) []int {
    result := make([]int, 0, len(data))
    for _, v := range data {
        if v > 0 {
            result = append(result, v*2)
        }
    }
    return result
}

// ✅ 更好：使用对象池
var resultPool = sync.Pool{
    New: func() interface{} {
        return make([]int, 0, 1000)
    },
}

func pooledMemory(data []int) []int {
    result := resultPool.Get().([]int)
    defer resultPool.Put(result[:0])
    
    for _, v := range data {
        if v > 0 {
            result = append(result, v*2)
        }
    }
    
    // 返回副本
    return append([]int(nil), result...)
}
```

### 3.2 优化字符串操作

```go
// ❌ 错误：频繁字符串拼接
func inefficientString(parts []string) string {
    result := ""
    for _, part := range parts {
        result += part
    }
    return result
}

// ✅ 正确：使用 strings.Builder
func efficientString(parts []string) string {
    var builder strings.Builder
    builder.Grow(len(parts) * 10) // 预分配容量
    
    for _, part := range parts {
        builder.WriteString(part)
    }
    return builder.String()
}

// ✅ 更好：使用 bytes.Buffer
func bufferString(parts []string) string {
    var buffer bytes.Buffer
    buffer.Grow(len(parts) * 10)
    
    for _, part := range parts {
        buffer.WriteString(part)
    }
    return buffer.String()
}
```

## 4. I/O 优化

### 4.1 批量处理

```go
// ❌ 错误：逐个处理
func inefficientIO(data []string) error {
    for _, item := range data {
        err := processItem(item)
        if err != nil {
            return err
        }
    }
    return nil
}

// ✅ 正确：批量处理
func efficientIO(data []string) error {
    const batchSize = 100
    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }
        
        batch := data[i:end]
        err := processBatch(batch)
        if err != nil {
            return err
        }
    }
    return nil
}

func processBatch(batch []string) error {
    // 批量处理逻辑
    for _, item := range batch {
        // 处理单个项目
        _ = item
    }
    return nil
}
```

### 4.2 使用连接池

```go
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

func (cp *ConnectionPool) Get() (*sql.DB, error) {
    select {
    case conn := <-cp.pool:
        return conn, nil
    default:
        // 创建新连接
        return sql.Open("mysql", "user:password@/dbname")
    }
}

func (cp *ConnectionPool) Put(conn *sql.DB) {
    select {
    case cp.pool <- conn:
    default:
        // 池已满，关闭连接
        conn.Close()
    }
}
```

## 5. 编译器优化

### 5.1 使用编译器标志

```bash
# 启用优化
go build -ldflags="-s -w" main.go

# 禁用调试信息
go build -ldflags="-s -w" -trimpath main.go

# 启用内联优化
go build -gcflags="-l=4" main.go

# 启用逃逸分析
go build -gcflags="-m" main.go
```

### 5.2 避免不必要的函数调用

```go
// ❌ 错误：不必要的函数调用
func inefficientCalls(data []int) int {
    sum := 0
    for _, v := range data {
        sum += getValue(v) // 函数调用开销
    }
    return sum
}

func getValue(v int) int {
    return v * 2
}

// ✅ 正确：内联计算
func efficientCalls(data []int) int {
    sum := 0
    for _, v := range data {
        sum += v * 2 // 直接计算
    }
    return sum
}
```

## 6. 缓存优化

### 6.1 使用本地缓存

```go
type LocalCache struct {
    mu    sync.RWMutex
    data  map[string]interface{}
    ttl   time.Duration
    times map[string]time.Time
}

func NewLocalCache(ttl time.Duration) *LocalCache {
    return &LocalCache{
        data:  make(map[string]interface{}),
        ttl:   ttl,
        times: make(map[string]time.Time),
    }
}

func (c *LocalCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if value, exists := c.data[key]; exists {
        if time.Since(c.times[key]) < c.ttl {
            return value, true
        }
        // 过期，删除
        delete(c.data, key)
        delete(c.times, key)
    }
    
    return nil, false
}

func (c *LocalCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data[key] = value
    c.times[key] = time.Now()
}
```

### 6.2 使用 LRU 缓存

```go
type LRUCache struct {
    capacity int
    cache    map[string]*Node
    head     *Node
    tail     *Node
    mu       sync.Mutex
}

type Node struct {
    key   string
    value interface{}
    prev  *Node
    next  *Node
}

func NewLRUCache(capacity int) *LRUCache {
    head := &Node{}
    tail := &Node{}
    head.next = tail
    tail.prev = head
    
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*Node),
        head:     head,
        tail:     tail,
    }
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if node, exists := lru.cache[key]; exists {
        lru.moveToHead(node)
        return node.value, true
    }
    
    return nil, false
}

func (lru *LRUCache) Set(key string, value interface{}) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if node, exists := lru.cache[key]; exists {
        node.value = value
        lru.moveToHead(node)
    } else {
        if len(lru.cache) >= lru.capacity {
            lru.removeTail()
        }
        
        newNode := &Node{
            key:   key,
            value: value,
        }
        lru.cache[key] = newNode
        lru.addToHead(newNode)
    }
}

func (lru *LRUCache) moveToHead(node *Node) {
    lru.removeNode(node)
    lru.addToHead(node)
}

func (lru *LRUCache) addToHead(node *Node) {
    node.prev = lru.head
    node.next = lru.head.next
    lru.head.next.prev = node
    lru.head.next = node
}

func (lru *LRUCache) removeNode(node *Node) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (lru *LRUCache) removeTail() {
    last := lru.tail.prev
    lru.removeNode(last)
    delete(lru.cache, last.key)
}
```

## 7. 性能测试和验证

### 7.1 基准测试

```go
func BenchmarkInefficient(b *testing.B) {
    data := make([]int, 1000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        inefficientSearch(data, 500)
    }
}

func BenchmarkEfficient(b *testing.B) {
    data := make([]int, 1000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        efficientSearch(data, 500)
    }
}

func BenchmarkBinarySearch(b *testing.B) {
    data := make([]int, 1000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        binarySearch(data, 500)
    }
}
```

### 7.2 性能对比

```go
func comparePerformance() {
    data := make([]int, 10000)
    for i := range data {
        data[i] = i
    }
    
    // 测试不同方法的性能
    methods := []struct {
        name string
        fn   func([]int, int) bool
    }{
        {"Inefficient", inefficientSearch},
        {"Efficient", efficientSearch},
        {"BinarySearch", binarySearch},
    }
    
    for _, method := range methods {
        start := time.Now()
        for i := 0; i < 1000; i++ {
            method.fn(data, 5000)
        }
        duration := time.Since(start)
        fmt.Printf("%s: %v\n", method.name, duration)
    }
}
```

## 8. 最佳实践总结

1. **算法优化**: 选择合适的数据结构和算法
2. **并发优化**: 合理使用 goroutine 和避免锁竞争
3. **内存优化**: 减少内存分配和使用对象池
4. **I/O 优化**: 批量处理和连接池
5. **编译器优化**: 使用合适的编译标志
6. **缓存优化**: 使用本地缓存和 LRU 缓存
7. **性能测试**: 使用基准测试验证优化效果
8. **持续监控**: 建立性能监控和预警机制

通过遵循这些优化策略，可以显著提高 Go 程序的 CPU 性能。
