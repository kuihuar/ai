# Go 网络性能优化策略

## 1. HTTP 服务优化

### 1.1 使用 HTTP/2

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
)

func main() {
    // 创建 HTTP/2 服务器
    server := &http.Server{
        Addr:    ":443",
        Handler: setupRoutes(),
        TLSConfig: &tls.Config{
            NextProtos: []string{"h2", "http/1.1"}, // 支持 HTTP/2
        },
    }
    
    // 启动 HTTPS 服务器（HTTP/2 需要 TLS）
    log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}

func setupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    
    // 设置路由
    mux.HandleFunc("/api/data", handleData)
    mux.HandleFunc("/api/users", handleUsers)
    
    return mux
}
```

### 1.2 启用压缩

```go
import (
    "compress/gzip"
    "net/http"
)

func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检查客户端是否支持 gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        // 设置 gzip 响应头
        w.Header().Set("Content-Encoding", "gzip")
        w.Header().Set("Vary", "Accept-Encoding")
        
        // 创建 gzip 写入器
        gz := gzip.NewWriter(w)
        defer gz.Close()
        
        // 包装 ResponseWriter
        gzw := &gzipResponseWriter{
            ResponseWriter: w,
            Writer:         gz,
        }
        
        next.ServeHTTP(gzw, r)
    })
}

type gzipResponseWriter struct {
    http.ResponseWriter
    *gzip.Writer
}

func (grw *gzipResponseWriter) Write(data []byte) (int, error) {
    return grw.Writer.Write(data)
}
```

### 1.3 优化响应头

```go
func optimizeHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 设置缓存头
        w.Header().Set("Cache-Control", "public, max-age=3600")
        w.Header().Set("ETag", generateETag(r.URL.Path))
        
        // 设置安全头
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // 设置性能头
        w.Header().Set("X-Response-Time", time.Since(start).String())
        
        next.ServeHTTP(w, r)
    })
}
```

## 2. 连接池优化

### 2.1 HTTP 客户端连接池

```go
import (
    "net/http"
    "time"
)

func createOptimizedHTTPClient() *http.Client {
    transport := &http.Transport{
        // 连接池配置
        MaxIdleConns:        100,              // 最大空闲连接数
        MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
        IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
        
        // 连接超时配置
        DialTimeout:         5 * time.Second,  // 连接超时
        TLSHandshakeTimeout: 5 * time.Second,  // TLS 握手超时
        
        // 请求超时配置
        ResponseHeaderTimeout: 10 * time.Second, // 响应头超时
        ExpectContinueTimeout: 1 * time.Second,  // Expect: 100-continue 超时
        
        // 其他优化
        DisableKeepAlives: false, // 启用 Keep-Alive
        DisableCompression: false, // 启用压缩
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second, // 总超时时间
    }
}
```

### 2.2 自定义连接池

```go
type ConnectionPool struct {
    pool    chan net.Conn
    factory func() (net.Conn, error)
    maxSize int
    mu      sync.RWMutex
    stats   PoolStats
}

type PoolStats struct {
    TotalConns    int
    ActiveConns   int
    IdleConns     int
    CreatedConns  int
    DestroyedConns int
}

func NewConnectionPool(factory func() (net.Conn, error), maxSize int) *ConnectionPool {
    return &ConnectionPool{
        pool:    make(chan net.Conn, maxSize),
        factory: factory,
        maxSize: maxSize,
    }
}

func (cp *ConnectionPool) Get() (net.Conn, error) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    
    select {
    case conn := <-cp.pool:
        cp.stats.ActiveConns++
        cp.stats.IdleConns--
        return conn, nil
    default:
        if cp.stats.TotalConns < cp.maxSize {
            conn, err := cp.factory()
            if err != nil {
                return nil, err
            }
            cp.stats.TotalConns++
            cp.stats.CreatedConns++
            cp.stats.ActiveConns++
            return conn, nil
        }
        
        // 等待可用连接
        select {
        case conn := <-cp.pool:
            cp.stats.ActiveConns++
            cp.stats.IdleConns--
            return conn, nil
        case <-time.After(5 * time.Second):
            return nil, fmt.Errorf("连接池超时")
        }
    }
}

func (cp *ConnectionPool) Put(conn net.Conn) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    
    select {
    case cp.pool <- conn:
        cp.stats.ActiveConns--
        cp.stats.IdleConns++
    default:
        // 池已满，关闭连接
        conn.Close()
        cp.stats.TotalConns--
        cp.stats.DestroyedConns++
    }
}
```

## 3. 并发优化

### 3.1 使用 worker pool 模式

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
    // 处理网络请求
    time.Sleep(100 * time.Millisecond) // 模拟网络延迟
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

### 3.2 实现请求限流

```go
type RateLimiter struct {
    tokens   int
    capacity int
    rate     int
    mu       sync.Mutex
    lastTime time.Time
}

func NewRateLimiter(rate, capacity int) *RateLimiter {
    return &RateLimiter{
        tokens:   capacity,
        capacity: capacity,
        rate:     rate,
        lastTime: time.Now(),
    }
}

func (rl *RateLimiter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(rl.lastTime)
    
    // 添加令牌
    tokensToAdd := int(elapsed.Seconds()) * rl.rate
    rl.tokens = min(rl.tokens+tokensToAdd, rl.capacity)
    rl.lastTime = now
    
    if rl.tokens > 0 {
        rl.tokens--
        return true
    }
    
    return false
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// 限流中间件
func rateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 4. 缓存优化

### 4.1 HTTP 缓存

```go
type HTTPCache struct {
    cache map[string]CacheEntry
    mu    sync.RWMutex
    ttl   time.Duration
}

type CacheEntry struct {
    Data      []byte
    Headers   map[string]string
    Timestamp time.Time
}

func NewHTTPCache(ttl time.Duration) *HTTPCache {
    return &HTTPCache{
        cache: make(map[string]CacheEntry),
        ttl:   ttl,
    }
}

func (hc *HTTPCache) Get(key string) ([]byte, map[string]string, bool) {
    hc.mu.RLock()
    defer hc.mu.RUnlock()
    
    entry, exists := hc.cache[key]
    if !exists {
        return nil, nil, false
    }
    
    if time.Since(entry.Timestamp) > hc.ttl {
        return nil, nil, false
    }
    
    return entry.Data, entry.Headers, true
}

func (hc *HTTPCache) Set(key string, data []byte, headers map[string]string) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    
    hc.cache[key] = CacheEntry{
        Data:      data,
        Headers:   headers,
        Timestamp: time.Now(),
    }
}

// 缓存中间件
func cacheMiddleware(cache *HTTPCache) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 只缓存 GET 请求
            if r.Method != http.MethodGet {
                next.ServeHTTP(w, r)
                return
            }
            
            // 生成缓存键
            cacheKey := r.URL.String()
            
            // 检查缓存
            if data, headers, found := cache.Get(cacheKey); found {
                // 设置响应头
                for key, value := range headers {
                    w.Header().Set(key, value)
                }
                w.Write(data)
                return
            }
            
            // 包装 ResponseWriter 以捕获响应
            cw := &cacheResponseWriter{
                ResponseWriter: w,
                cache:          cache,
                cacheKey:       cacheKey,
            }
            
            next.ServeHTTP(cw, r)
        })
    }
}

type cacheResponseWriter struct {
    http.ResponseWriter
    cache    *HTTPCache
    cacheKey string
    data     []byte
    headers  map[string]string
}

func (cw *cacheResponseWriter) Write(data []byte) (int, error) {
    cw.data = append(cw.data, data...)
    return cw.ResponseWriter.Write(data)
}

func (cw *cacheResponseWriter) WriteHeader(statusCode int) {
    cw.headers = make(map[string]string)
    for key, values := range cw.ResponseWriter.Header() {
        if len(values) > 0 {
            cw.headers[key] = values[0]
        }
    }
    cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *cacheResponseWriter) Close() {
    if len(cw.data) > 0 {
        cw.cache.Set(cw.cacheKey, cw.data, cw.headers)
    }
}
```

### 4.2 应用层缓存

```go
type AppCache struct {
    cache map[string]interface{}
    mu    sync.RWMutex
    ttl   time.Duration
    times map[string]time.Time
}

func NewAppCache(ttl time.Duration) *AppCache {
    return &AppCache{
        cache: make(map[string]interface{}),
        ttl:   ttl,
        times: make(map[string]time.Time),
    }
}

func (ac *AppCache) Get(key string) (interface{}, bool) {
    ac.mu.RLock()
    defer ac.mu.RUnlock()
    
    if value, exists := ac.cache[key]; exists {
        if time.Since(ac.times[key]) < ac.ttl {
            return value, true
        }
        // 过期，删除
        delete(ac.cache, key)
        delete(ac.times, key)
    }
    
    return nil, false
}

func (ac *AppCache) Set(key string, value interface{}) {
    ac.mu.Lock()
    defer ac.mu.Unlock()
    
    ac.cache[key] = value
    ac.times[key] = time.Now()
}
```

## 5. 负载均衡

### 5.1 轮询负载均衡

```go
type RoundRobinBalancer struct {
    servers []string
    current int
    mu      sync.Mutex
}

func NewRoundRobinBalancer(servers []string) *RoundRobinBalancer {
    return &RoundRobinBalancer{
        servers: servers,
    }
}

func (rrb *RoundRobinBalancer) GetServer() string {
    rrb.mu.Lock()
    defer rrb.mu.Unlock()
    
    if len(rrb.servers) == 0 {
        return ""
    }
    
    server := rrb.servers[rrb.current]
    rrb.current = (rrb.current + 1) % len(rrb.servers)
    return server
}
```

### 5.2 加权轮询负载均衡

```go
type WeightedServer struct {
    Server string
    Weight int
}

type WeightedRoundRobinBalancer struct {
    servers []WeightedServer
    current int
    mu      sync.Mutex
}

func NewWeightedRoundRobinBalancer(servers []WeightedServer) *WeightedRoundRobinBalancer {
    return &WeightedRoundRobinBalancer{
        servers: servers,
    }
}

func (wrrb *WeightedRoundRobinBalancer) GetServer() string {
    wrrb.mu.Lock()
    defer wrrb.mu.Unlock()
    
    if len(wrrb.servers) == 0 {
        return ""
    }
    
    // 找到权重最大的服务器
    maxWeight := 0
    selectedIndex := 0
    
    for i, server := range wrrb.servers {
        if server.Weight > maxWeight {
            maxWeight = server.Weight
            selectedIndex = i
        }
    }
    
    // 减少选中服务器的权重
    wrrb.servers[selectedIndex].Weight--
    
    // 如果所有权重都为0，重置权重
    allZero := true
    for _, server := range wrrb.servers {
        if server.Weight > 0 {
            allZero = false
            break
        }
    }
    
    if allZero {
        for i := range wrrb.servers {
            wrrb.servers[i].Weight = 1 // 重置为1
        }
    }
    
    return wrrb.servers[selectedIndex].Server
}
```

## 6. 监控和预警

### 6.1 性能监控

```go
type NetworkMonitor struct {
    requestCount    int64
    errorCount      int64
    totalLatency    time.Duration
    startTime       time.Time
    mu              sync.RWMutex
}

func NewNetworkMonitor() *NetworkMonitor {
    return &NetworkMonitor{
        startTime: time.Now(),
    }
}

func (nm *NetworkMonitor) RecordRequest(latency time.Duration) {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    
    atomic.AddInt64(&nm.requestCount, 1)
    nm.totalLatency += latency
}

func (nm *NetworkMonitor) RecordError() {
    atomic.AddInt64(&nm.errorCount, 1)
}

func (nm *NetworkMonitor) GetStats() map[string]interface{} {
    nm.mu.RLock()
    defer nm.mu.RUnlock()
    
    duration := time.Since(nm.startTime)
    avgLatency := time.Duration(0)
    
    if nm.requestCount > 0 {
        avgLatency = nm.totalLatency / time.Duration(nm.requestCount)
    }
    
    return map[string]interface{}{
        "request_count":   nm.requestCount,
        "error_count":     nm.errorCount,
        "avg_latency":     avgLatency,
        "requests_per_sec": float64(nm.requestCount) / duration.Seconds(),
        "error_rate":      float64(nm.errorCount) / float64(nm.requestCount),
    }
}
```

### 6.2 预警系统

```go
type AlertSystem struct {
    monitors []Monitor
    alerts   chan Alert
}

type Monitor interface {
    Check() []Alert
}

type Alert struct {
    Type        string
    Message     string
    Timestamp   time.Time
    Severity    string
}

func NewAlertSystem() *AlertSystem {
    return &AlertSystem{
        monitors: make([]Monitor, 0),
        alerts:   make(chan Alert, 100),
    }
}

func (as *AlertSystem) AddMonitor(monitor Monitor) {
    as.monitors = append(as.monitors, monitor)
}

func (as *AlertSystem) Start() {
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            for _, monitor := range as.monitors {
                alerts := monitor.Check()
                for _, alert := range alerts {
                    select {
                    case as.alerts <- alert:
                    default:
                        // 通道已满，丢弃警告
                    }
                }
            }
        }
    }()
}

func (as *AlertSystem) GetAlerts() <-chan Alert {
    return as.alerts
}
```

## 7. 最佳实践总结

1. **使用 HTTP/2**: 提高多路复用和头部压缩
2. **启用压缩**: 减少传输数据量
3. **优化连接池**: 合理配置连接参数
4. **实现缓存**: 减少重复请求
5. **使用负载均衡**: 分散请求压力
6. **实现限流**: 防止系统过载
7. **监控性能**: 实时监控网络指标
8. **设置预警**: 及时发现和处理问题

通过遵循这些优化策略，可以显著提高 Go 程序的网络性能。
