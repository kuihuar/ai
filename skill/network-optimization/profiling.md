# Go 网络性能分析

## 1. 网络性能分析工具

### 1.1 使用 pprof 分析网络性能

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "time"
)

func main() {
    // 启用 pprof
    go func() {
        log.Println("pprof server started on :6060")
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 启动网络服务
    startNetworkService()
}

func startNetworkService() {
    // HTTP 服务
    http.HandleFunc("/api/data", handleData)
    http.HandleFunc("/api/users", handleUsers)
    http.HandleFunc("/api/health", handleHealth)
    
    // 启动服务器
    server := &http.Server{
        Addr:         ":8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    log.Println("Server started on :8080")
    log.Fatal(server.ListenAndServe())
}

func handleData(w http.ResponseWriter, r *http.Request) {
    // 模拟数据处理
    time.Sleep(10 * time.Millisecond)
    
    data := map[string]interface{}{
        "timestamp": time.Now().Unix(),
        "data":      "sample data",
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    // 模拟用户数据查询
    time.Sleep(50 * time.Millisecond)
    
    users := []map[string]interface{}{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
        {"id": 3, "name": "Charlie"},
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(users)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

### 1.2 网络性能分析命令

```bash
# 分析 CPU 使用情况（包含网络处理）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 分析内存使用情况
go tool pprof http://localhost:6060/debug/pprof/heap

# 分析 goroutine 使用情况
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 分析阻塞情况
go tool pprof http://localhost:6060/debug/pprof/block

# 分析互斥锁使用情况
go tool pprof http://localhost:6060/debug/pprof/mutex

# 生成火焰图
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile

# 查看网络相关的函数
go tool pprof -focus=net http://localhost:6060/debug/pprof/profile
```

## 2. 网络连接分析

### 2.1 连接数监控

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "sync"
    "time"
)

type ConnectionMonitor struct {
    connections map[net.Conn]time.Time
    mu          sync.RWMutex
    maxConns    int
}

func NewConnectionMonitor(maxConns int) *ConnectionMonitor {
    return &ConnectionMonitor{
        connections: make(map[net.Conn]time.Time),
        maxConns:    maxConns,
    }
}

func (cm *ConnectionMonitor) AddConnection(conn net.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.connections[conn] = time.Now()
    
    if len(cm.connections) > cm.maxConns {
        log.Printf("警告: 连接数超过限制 %d", cm.maxConns)
    }
}

func (cm *ConnectionMonitor) RemoveConnection(conn net.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    delete(cm.connections, conn)
}

func (cm *ConnectionMonitor) GetConnectionCount() int {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    return len(cm.connections)
}

func (cm *ConnectionMonitor) GetConnectionStats() map[string]interface{} {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    now := time.Now()
    stats := map[string]interface{}{
        "total_connections": len(cm.connections),
        "max_connections":   cm.maxConns,
        "connection_ages":   make([]time.Duration, 0),
    }
    
    for _, startTime := range cm.connections {
        age := now.Sub(startTime)
        stats["connection_ages"] = append(stats["connection_ages"].([]time.Duration), age)
    }
    
    return stats
}

// 自定义 HTTP 服务器，支持连接监控
type MonitoredServer struct {
    *http.Server
    monitor *ConnectionMonitor
}

func NewMonitoredServer(addr string, monitor *ConnectionMonitor) *MonitoredServer {
    server := &http.Server{
        Addr: addr,
    }
    
    return &MonitoredServer{
        Server:  server,
        monitor: monitor,
    }
}

func (ms *MonitoredServer) ListenAndServe() error {
    listener, err := net.Listen("tcp", ms.Addr)
    if err != nil {
        return err
    }
    
    return ms.Serve(ms.monitorConnection(listener))
}

func (ms *MonitoredServer) monitorConnection(listener net.Listener) net.Listener {
    return &monitoredListener{
        Listener: listener,
        monitor:  ms.monitor,
    }
}

type monitoredListener struct {
    net.Listener
    monitor *ConnectionMonitor
}

func (ml *monitoredListener) Accept() (net.Conn, error) {
    conn, err := ml.Listener.Accept()
    if err != nil {
        return nil, err
    }
    
    ml.monitor.AddConnection(conn)
    
    return &monitoredConn{
        Conn:    conn,
        monitor: ml.monitor,
    }, nil
}

type monitoredConn struct {
    net.Conn
    monitor *ConnectionMonitor
}

func (mc *monitoredConn) Close() error {
    mc.monitor.RemoveConnection(mc.Conn)
    return mc.Conn.Close()
}
```

### 2.2 连接池分析

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

func (cp *ConnectionPool) GetStats() PoolStats {
    cp.mu.RLock()
    defer cp.mu.RUnlock()
    
    return cp.stats
}
```

## 3. 网络延迟分析

### 3.1 延迟监控

```go
type LatencyMonitor struct {
    samples []time.Duration
    mu      sync.RWMutex
    maxSamples int
}

func NewLatencyMonitor(maxSamples int) *LatencyMonitor {
    return &LatencyMonitor{
        samples:    make([]time.Duration, 0, maxSamples),
        maxSamples: maxSamples,
    }
}

func (lm *LatencyMonitor) RecordLatency(duration time.Duration) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    
    lm.samples = append(lm.samples, duration)
    if len(lm.samples) > lm.maxSamples {
        lm.samples = lm.samples[1:]
    }
}

func (lm *LatencyMonitor) GetStats() map[string]interface{} {
    lm.mu.RLock()
    defer lm.mu.RUnlock()
    
    if len(lm.samples) == 0 {
        return map[string]interface{}{
            "count": 0,
        }
    }
    
    var total time.Duration
    var min, max time.Duration = lm.samples[0], lm.samples[0]
    
    for _, sample := range lm.samples {
        total += sample
        if sample < min {
            min = sample
        }
        if sample > max {
            max = sample
        }
    }
    
    avg := total / time.Duration(len(lm.samples))
    
    return map[string]interface{}{
        "count":    len(lm.samples),
        "min":      min,
        "max":      max,
        "avg":      avg,
        "total":    total,
    }
}

// 中间件：记录请求延迟
func latencyMiddleware(monitor *LatencyMonitor) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // 包装 ResponseWriter 以捕获状态码
            ww := &responseWriter{ResponseWriter: w, statusCode: 200}
            
            next.ServeHTTP(ww, r)
            
            duration := time.Since(start)
            monitor.RecordLatency(duration)
            
            log.Printf("请求 %s %s 耗时 %v, 状态码 %d", 
                r.Method, r.URL.Path, duration, ww.statusCode)
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 3.2 延迟分析工具

```go
func analyzeLatency(monitor *LatencyMonitor) {
    stats := monitor.GetStats()
    
    fmt.Printf("延迟统计:\n")
    fmt.Printf("  请求数量: %d\n", stats["count"])
    fmt.Printf("  最小延迟: %v\n", stats["min"])
    fmt.Printf("  最大延迟: %v\n", stats["max"])
    fmt.Printf("  平均延迟: %v\n", stats["avg"])
    fmt.Printf("  总延迟: %v\n", stats["total"])
    
    // 分析延迟分布
    analyzeLatencyDistribution(monitor)
}

func analyzeLatencyDistribution(monitor *LatencyMonitor) {
    monitor.mu.RLock()
    defer monitor.mu.RUnlock()
    
    if len(monitor.samples) == 0 {
        return
    }
    
    // 计算百分位数
    percentiles := []float64{50, 90, 95, 99}
    
    for _, p := range percentiles {
        idx := int(float64(len(monitor.samples)) * p / 100)
        if idx >= len(monitor.samples) {
            idx = len(monitor.samples) - 1
        }
        
        // 简单的排序（实际应用中应该使用更高效的排序算法）
        samples := make([]time.Duration, len(monitor.samples))
        copy(samples, monitor.samples)
        sort.Slice(samples, func(i, j int) bool {
            return samples[i] < samples[j]
        })
        
        fmt.Printf("  P%.0f: %v\n", p, samples[idx])
    }
}
```

## 4. 网络吞吐量分析

### 4.1 吞吐量监控

```go
type ThroughputMonitor struct {
    bytesTransferred int64
    requestsProcessed int64
    startTime        time.Time
    mu               sync.RWMutex
}

func NewThroughputMonitor() *ThroughputMonitor {
    return &ThroughputMonitor{
        startTime: time.Now(),
    }
}

func (tm *ThroughputMonitor) RecordBytes(bytes int64) {
    tm.mu.Lock()
    defer tm.mu.Unlock()
    
    atomic.AddInt64(&tm.bytesTransferred, bytes)
}

func (tm *ThroughputMonitor) RecordRequest() {
    tm.mu.Lock()
    defer tm.mu.Unlock()
    
    atomic.AddInt64(&tm.requestsProcessed, 1)
}

func (tm *ThroughputMonitor) GetThroughput() map[string]interface{} {
    tm.mu.RLock()
    defer tm.mu.RUnlock()
    
    duration := time.Since(tm.startTime)
    seconds := duration.Seconds()
    
    bytesPerSecond := float64(tm.bytesTransferred) / seconds
    requestsPerSecond := float64(tm.requestsProcessed) / seconds
    
    return map[string]interface{}{
        "duration_seconds":    seconds,
        "bytes_transferred":   tm.bytesTransferred,
        "requests_processed":  tm.requestsProcessed,
        "bytes_per_second":    bytesPerSecond,
        "requests_per_second": requestsPerSecond,
    }
}

// 中间件：记录吞吐量
func throughputMiddleware(monitor *ThroughputMonitor) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 记录请求
            monitor.RecordRequest()
            
            // 包装 ResponseWriter 以捕获响应大小
            ww := &throughputResponseWriter{
                ResponseWriter: w,
                monitor:        monitor,
            }
            
            next.ServeHTTP(ww, r)
        })
    }
}

type throughputResponseWriter struct {
    http.ResponseWriter
    monitor *ThroughputMonitor
}

func (trw *throughputResponseWriter) Write(data []byte) (int, error) {
    n, err := trw.ResponseWriter.Write(data)
    trw.monitor.RecordBytes(int64(n))
    return n, err
}
```

### 4.2 吞吐量分析

```go
func analyzeThroughput(monitor *ThroughputMonitor) {
    stats := monitor.GetThroughput()
    
    fmt.Printf("吞吐量统计:\n")
    fmt.Printf("  运行时间: %.2f 秒\n", stats["duration_seconds"])
    fmt.Printf("  传输字节数: %d\n", stats["bytes_transferred"])
    fmt.Printf("  处理请求数: %d\n", stats["requests_processed"])
    fmt.Printf("  字节/秒: %.2f\n", stats["bytes_per_second"])
    fmt.Printf("  请求/秒: %.2f\n", stats["requests_per_second"])
    
    // 分析吞吐量趋势
    analyzeThroughputTrend(monitor)
}

func analyzeThroughputTrend(monitor *ThroughputMonitor) {
    // 这里可以实现更复杂的趋势分析
    // 比如计算吞吐量的变化率、预测未来趋势等
    fmt.Printf("吞吐量趋势分析: 需要实现更复杂的算法\n")
}
```

## 5. 网络错误分析

### 5.1 错误监控

```go
type ErrorMonitor struct {
    errors map[string]int
    mu     sync.RWMutex
}

func NewErrorMonitor() *ErrorMonitor {
    return &ErrorMonitor{
        errors: make(map[string]int),
    }
}

func (em *ErrorMonitor) RecordError(errorType string) {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    em.errors[errorType]++
}

func (em *ErrorMonitor) GetErrorStats() map[string]interface{} {
    em.mu.RLock()
    defer em.mu.RUnlock()
    
    totalErrors := 0
    for _, count := range em.errors {
        totalErrors += count
    }
    
    return map[string]interface{}{
        "total_errors": totalErrors,
        "error_types":  em.errors,
    }
}

// 中间件：记录错误
func errorMiddleware(monitor *ErrorMonitor) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 包装 ResponseWriter 以捕获错误
            ww := &errorResponseWriter{
                ResponseWriter: w,
                monitor:        monitor,
            }
            
            next.ServeHTTP(ww, r)
        })
    }
}

type errorResponseWriter struct {
    http.ResponseWriter
    monitor *ErrorMonitor
}

func (erw *errorResponseWriter) WriteHeader(code int) {
    if code >= 400 {
        errorType := fmt.Sprintf("http_%d", code)
        erw.monitor.RecordError(errorType)
    }
    erw.ResponseWriter.WriteHeader(code)
}
```

## 6. 综合网络性能分析

### 6.1 性能分析器

```go
type NetworkProfiler struct {
    connectionMonitor *ConnectionMonitor
    latencyMonitor    *LatencyMonitor
    throughputMonitor *ThroughputMonitor
    errorMonitor      *ErrorMonitor
}

func NewNetworkProfiler(maxConns int) *NetworkProfiler {
    return &NetworkProfiler{
        connectionMonitor: NewConnectionMonitor(maxConns),
        latencyMonitor:    NewLatencyMonitor(1000),
        throughputMonitor: NewThroughputMonitor(),
        errorMonitor:      NewErrorMonitor(),
    }
}

func (np *NetworkProfiler) GetComprehensiveStats() map[string]interface{} {
    return map[string]interface{}{
        "connections": np.connectionMonitor.GetConnectionStats(),
        "latency":     np.latencyMonitor.GetStats(),
        "throughput":  np.throughputMonitor.GetThroughput(),
        "errors":      np.errorMonitor.GetErrorStats(),
    }
}

func (np *NetworkProfiler) StartMonitoring() {
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            stats := np.GetComprehensiveStats()
            log.Printf("网络性能统计: %+v", stats)
        }
    }()
}
```

### 6.2 性能报告生成

```go
func generatePerformanceReport(profiler *NetworkProfiler) {
    stats := profiler.GetComprehensiveStats()
    
    fmt.Println("=== 网络性能报告 ===")
    fmt.Printf("生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Println()
    
    // 连接统计
    if connStats, ok := stats["connections"].(map[string]interface{}); ok {
        fmt.Println("连接统计:")
        fmt.Printf("  总连接数: %d\n", connStats["total_connections"])
        fmt.Printf("  最大连接数: %d\n", connStats["max_connections"])
        fmt.Println()
    }
    
    // 延迟统计
    if latencyStats, ok := stats["latency"].(map[string]interface{}); ok {
        fmt.Println("延迟统计:")
        fmt.Printf("  请求数量: %d\n", latencyStats["count"])
        fmt.Printf("  平均延迟: %v\n", latencyStats["avg"])
        fmt.Printf("  最小延迟: %v\n", latencyStats["min"])
        fmt.Printf("  最大延迟: %v\n", latencyStats["max"])
        fmt.Println()
    }
    
    // 吞吐量统计
    if throughputStats, ok := stats["throughput"].(map[string]interface{}); ok {
        fmt.Println("吞吐量统计:")
        fmt.Printf("  运行时间: %.2f 秒\n", throughputStats["duration_seconds"])
        fmt.Printf("  传输字节数: %d\n", throughputStats["bytes_transferred"])
        fmt.Printf("  处理请求数: %d\n", throughputStats["requests_processed"])
        fmt.Printf("  字节/秒: %.2f\n", throughputStats["bytes_per_second"])
        fmt.Printf("  请求/秒: %.2f\n", throughputStats["requests_per_second"])
        fmt.Println()
    }
    
    // 错误统计
    if errorStats, ok := stats["errors"].(map[string]interface{}); ok {
        fmt.Println("错误统计:")
        fmt.Printf("  总错误数: %d\n", errorStats["total_errors"])
        if errorTypes, ok := errorStats["error_types"].(map[string]int); ok {
            for errorType, count := range errorTypes {
                fmt.Printf("  %s: %d\n", errorType, count)
            }
        }
        fmt.Println()
    }
}
```

通过使用这些网络性能分析工具和技巧，可以有效地识别和解决 Go 程序中的网络性能问题。
