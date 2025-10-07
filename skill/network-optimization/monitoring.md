# Go 网络性能监控工具

## 1. 内置监控工具

### 1.1 使用 pprof 进行网络监控

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

## 2. 自定义监控系统

### 2.1 网络连接监控

```go
type NetworkMonitor struct {
    connections map[net.Conn]ConnectionInfo
    mu          sync.RWMutex
    stats       NetworkStats
}

type ConnectionInfo struct {
    StartTime      time.Time
    BytesSent      int64
    BytesReceived  int64
    LastActivity   time.Time
    RequestCount   int64
    ErrorCount     int64
}

type NetworkStats struct {
    TotalConnections   int
    ActiveConnections  int
    TotalBytesSent     int64
    TotalBytesReceived int64
    TotalRequests      int64
    TotalErrors        int64
    AverageLatency     time.Duration
}

func NewNetworkMonitor() *NetworkMonitor {
    return &NetworkMonitor{
        connections: make(map[net.Conn]ConnectionInfo),
    }
}

func (nm *NetworkMonitor) AddConnection(conn net.Conn) {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    
    nm.connections[conn] = ConnectionInfo{
        StartTime:    time.Now(),
        LastActivity: time.Now(),
    }
    
    nm.stats.TotalConnections++
    nm.stats.ActiveConnections++
}

func (nm *NetworkMonitor) RemoveConnection(conn net.Conn) {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    
    if info, exists := nm.connections[conn]; exists {
        nm.stats.TotalBytesSent += info.BytesSent
        nm.stats.TotalBytesReceived += info.BytesReceived
        nm.stats.TotalRequests += info.RequestCount
        nm.stats.TotalErrors += info.ErrorCount
        delete(nm.connections, conn)
        nm.stats.ActiveConnections--
    }
}

func (nm *NetworkMonitor) RecordBytes(conn net.Conn, sent, received int64) {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    
    if info, exists := nm.connections[conn]; exists {
        info.BytesSent += sent
        info.BytesReceived += received
        info.LastActivity = time.Now()
        nm.connections[conn] = info
    }
}

func (nm *NetworkMonitor) RecordRequest(conn net.Conn) {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    
    if info, exists := nm.connections[conn]; exists {
        info.RequestCount++
        info.LastActivity = time.Now()
        nm.connections[conn] = info
    }
}

func (nm *NetworkMonitor) RecordError(conn net.Conn) {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    
    if info, exists := nm.connections[conn]; exists {
        info.ErrorCount++
        info.LastActivity = time.Now()
        nm.connections[conn] = info
    }
}

func (nm *NetworkMonitor) GetStats() NetworkStats {
    nm.mu.RLock()
    defer nm.mu.RUnlock()
    return nm.stats
}

func (nm *NetworkMonitor) GetConnectionStats() map[string]interface{} {
    nm.mu.RLock()
    defer nm.mu.RUnlock()
    
    stats := make(map[string]interface{})
    stats["total_connections"] = nm.stats.TotalConnections
    stats["active_connections"] = nm.stats.ActiveConnections
    stats["total_bytes_sent"] = nm.stats.TotalBytesSent
    stats["total_bytes_received"] = nm.stats.TotalBytesReceived
    stats["total_requests"] = nm.stats.TotalRequests
    stats["total_errors"] = nm.stats.TotalErrors
    
    if nm.stats.TotalRequests > 0 {
        stats["error_rate"] = float64(nm.stats.TotalErrors) / float64(nm.stats.TotalRequests)
    } else {
        stats["error_rate"] = 0.0
    }
    
    return stats
}
```

### 2.2 延迟监控

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

func (lm *LatencyMonitor) GetPercentiles() map[string]time.Duration {
    lm.mu.RLock()
    defer lm.mu.RUnlock()
    
    if len(lm.samples) == 0 {
        return nil
    }
    
    // 复制并排序样本
    sortedSamples := make([]time.Duration, len(lm.samples))
    copy(sortedSamples, lm.samples)
    sort.Slice(sortedSamples, func(i, j int) bool {
        return sortedSamples[i] < sortedSamples[j]
    })
    
    percentiles := map[string]time.Duration{
        "p50": sortedSamples[int(float64(len(sortedSamples))*0.5)],
        "p90": sortedSamples[int(float64(len(sortedSamples))*0.9)],
        "p95": sortedSamples[int(float64(len(sortedSamples))*0.95)],
        "p99": sortedSamples[int(float64(len(sortedSamples))*0.99)],
    }
    
    return percentiles
}
```

### 2.3 吞吐量监控

```go
type ThroughputMonitor struct {
    bytesTransferred  int64
    requestsProcessed int64
    startTime         time.Time
    mu                sync.RWMutex
}

func NewThroughputMonitor() *ThroughputMonitor {
    return &ThroughputMonitor{
        startTime: time.Now(),
    }
}

func (tm *ThroughputMonitor) RecordBytes(bytes int64) {
    atomic.AddInt64(&tm.bytesTransferred, bytes)
}

func (tm *ThroughputMonitor) RecordRequest() {
    atomic.AddInt64(&tm.requestsProcessed, 1)
}

func (tm *ThroughputMonitor) GetThroughput() map[string]interface{} {
    duration := time.Since(tm.startTime)
    seconds := duration.Seconds()
    
    bytesPerSecond := float64(atomic.LoadInt64(&tm.bytesTransferred)) / seconds
    requestsPerSecond := float64(atomic.LoadInt64(&tm.requestsProcessed)) / seconds
    
    return map[string]interface{}{
        "duration_seconds":    seconds,
        "bytes_transferred":   atomic.LoadInt64(&tm.bytesTransferred),
        "requests_processed":  atomic.LoadInt64(&tm.requestsProcessed),
        "bytes_per_second":    bytesPerSecond,
        "requests_per_second": requestsPerSecond,
    }
}
```

## 3. 监控中间件

### 3.1 HTTP 监控中间件

```go
func monitoringMiddleware(monitor *NetworkMonitor, latencyMonitor *LatencyMonitor, throughputMonitor *ThroughputMonitor) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // 记录请求
            throughputMonitor.RecordRequest()
            
            // 包装 ResponseWriter 以捕获响应大小
            mw := &monitoringResponseWriter{
                ResponseWriter: w,
                monitor:        monitor,
                throughputMonitor: throughputMonitor,
            }
            
            next.ServeHTTP(mw, r)
            
            // 记录延迟
            latency := time.Since(start)
            latencyMonitor.RecordLatency(latency)
            
            // 记录错误
            if mw.statusCode >= 400 {
                monitor.RecordError(nil) // 这里需要连接信息
            }
        })
    }
}

type monitoringResponseWriter struct {
    http.ResponseWriter
    monitor           *NetworkMonitor
    throughputMonitor *ThroughputMonitor
    statusCode        int
}

func (mw *monitoringResponseWriter) Write(data []byte) (int, error) {
    n, err := mw.ResponseWriter.Write(data)
    mw.throughputMonitor.RecordBytes(int64(n))
    return n, err
}

func (mw *monitoringResponseWriter) WriteHeader(code int) {
    mw.statusCode = code
    mw.ResponseWriter.WriteHeader(code)
}
```

### 3.2 连接监控中间件

```go
func connectionMonitoringMiddleware(monitor *NetworkMonitor) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 获取连接信息
            if conn, ok := getConnectionFromRequest(r); ok {
                monitor.RecordRequest(conn)
                
                // 包装 ResponseWriter 以监控字节传输
                cw := &connectionMonitoringResponseWriter{
                    ResponseWriter: w,
                    monitor:        monitor,
                    conn:           conn,
                }
                
                next.ServeHTTP(cw, r)
            } else {
                next.ServeHTTP(w, r)
            }
        })
    }
}

type connectionMonitoringResponseWriter struct {
    http.ResponseWriter
    monitor *NetworkMonitor
    conn    net.Conn
}

func (cw *connectionMonitoringResponseWriter) Write(data []byte) (int, error) {
    n, err := cw.ResponseWriter.Write(data)
    cw.monitor.RecordBytes(cw.conn, 0, int64(n))
    return n, err
}

func getConnectionFromRequest(r *http.Request) (net.Conn, bool) {
    // 这里需要根据具体的实现来获取连接
    // 通常需要自定义的 ResponseWriter 来捕获连接信息
    return nil, false
}
```

## 4. 监控数据收集

### 4.1 指标收集器

```go
type MetricsCollector struct {
    metrics map[string]interface{}
    mu      sync.RWMutex
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        metrics: make(map[string]interface{}),
    }
}

func (mc *MetricsCollector) SetMetric(key string, value interface{}) {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    mc.metrics[key] = value
}

func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
    mc.mu.RLock()
    defer mc.mu.RUnlock()
    
    result := make(map[string]interface{})
    for k, v := range mc.metrics {
        result[k] = v
    }
    return result
}

func (mc *MetricsCollector) StartCollection(monitor *NetworkMonitor, latencyMonitor *LatencyMonitor, throughputMonitor *ThroughputMonitor) {
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            // 收集网络指标
            networkStats := monitor.GetConnectionStats()
            for k, v := range networkStats {
                mc.SetMetric("network_"+k, v)
            }
            
            // 收集延迟指标
            latencyStats := latencyMonitor.GetStats()
            for k, v := range latencyStats {
                mc.SetMetric("latency_"+k, v)
            }
            
            // 收集吞吐量指标
            throughputStats := throughputMonitor.GetThroughput()
            for k, v := range throughputStats {
                mc.SetMetric("throughput_"+k, v)
            }
            
            // 收集系统指标
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            mc.SetMetric("memory_alloc", m.Alloc)
            mc.SetMetric("memory_total", m.TotalAlloc)
            mc.SetMetric("goroutines", runtime.NumGoroutine())
        }
    }()
}
```

### 4.2 监控数据导出

```go
type MetricsExporter struct {
    collector *MetricsCollector
    interval  time.Duration
}

func NewMetricsExporter(collector *MetricsCollector, interval time.Duration) *MetricsExporter {
    return &MetricsExporter{
        collector: collector,
        interval:  interval,
    }
}

func (me *MetricsExporter) StartExport() {
    go func() {
        ticker := time.NewTicker(me.interval)
        defer ticker.Stop()
        
        for range ticker.C {
            metrics := me.collector.GetMetrics()
            me.exportMetrics(metrics)
        }
    }()
}

func (me *MetricsExporter) exportMetrics(metrics map[string]interface{}) {
    // 导出到日志
    log.Printf("Metrics: %+v", metrics)
    
    // 导出到文件
    me.exportToFile(metrics)
    
    // 导出到外部监控系统
    me.exportToExternalSystem(metrics)
}

func (me *MetricsExporter) exportToFile(metrics map[string]interface{}) {
    filename := fmt.Sprintf("metrics_%s.json", time.Now().Format("2006-01-02_15-04-05"))
    
    data, err := json.MarshalIndent(metrics, "", "  ")
    if err != nil {
        log.Printf("Error marshaling metrics: %v", err)
        return
    }
    
    err = ioutil.WriteFile(filename, data, 0644)
    if err != nil {
        log.Printf("Error writing metrics file: %v", err)
    }
}

func (me *MetricsExporter) exportToExternalSystem(metrics map[string]interface{}) {
    // 这里可以实现导出到 Prometheus、InfluxDB 等外部监控系统
    // 示例：导出到 HTTP 端点
    data, err := json.Marshal(metrics)
    if err != nil {
        return
    }
    
    resp, err := http.Post("http://monitoring-system:8080/metrics", "application/json", bytes.NewBuffer(data))
    if err != nil {
        log.Printf("Error exporting to external system: %v", err)
        return
    }
    resp.Body.Close()
}
```

## 5. 预警系统

### 5.1 预警规则

```go
type AlertRule struct {
    Name        string
    Condition   func(map[string]interface{}) bool
    Severity    string
    Message     string
}

type AlertSystem struct {
    rules   []AlertRule
    alerts  chan Alert
    mu      sync.RWMutex
}

type Alert struct {
    Rule      string
    Message   string
    Severity  string
    Timestamp time.Time
    Metrics   map[string]interface{}
}

func NewAlertSystem() *AlertSystem {
    return &AlertSystem{
        alerts: make(chan Alert, 100),
    }
}

func (as *AlertSystem) AddRule(rule AlertRule) {
    as.mu.Lock()
    defer as.mu.Unlock()
    
    as.rules = append(as.rules, rule)
}

func (as *AlertSystem) CheckRules(metrics map[string]interface{}) {
    as.mu.RLock()
    defer as.mu.RUnlock()
    
    for _, rule := range as.rules {
        if rule.Condition(metrics) {
            alert := Alert{
                Rule:      rule.Name,
                Message:   rule.Message,
                Severity:  rule.Severity,
                Timestamp: time.Now(),
                Metrics:   metrics,
            }
            
            select {
            case as.alerts <- alert:
            default:
                // 通道已满，丢弃警告
            }
        }
    }
}

func (as *AlertSystem) GetAlerts() <-chan Alert {
    return as.alerts
}
```

### 5.2 预警规则示例

```go
func setupAlertRules(alertSystem *AlertSystem) {
    // 高错误率预警
    alertSystem.AddRule(AlertRule{
        Name: "high_error_rate",
        Condition: func(metrics map[string]interface{}) bool {
            if errorRate, ok := metrics["network_error_rate"].(float64); ok {
                return errorRate > 0.1 // 错误率超过 10%
            }
            return false
        },
        Severity: "warning",
        Message:  "Error rate is too high",
    })
    
    // 高延迟预警
    alertSystem.AddRule(AlertRule{
        Name: "high_latency",
        Condition: func(metrics map[string]interface{}) bool {
            if avgLatency, ok := metrics["latency_avg"].(time.Duration); ok {
                return avgLatency > 100*time.Millisecond
            }
            return false
        },
        Severity: "warning",
        Message:  "Average latency is too high",
    })
    
    // 高内存使用预警
    alertSystem.AddRule(AlertRule{
        Name: "high_memory_usage",
        Condition: func(metrics map[string]interface{}) bool {
            if memoryAlloc, ok := metrics["memory_alloc"].(uint64); ok {
                return memoryAlloc > 100*1024*1024 // 超过 100MB
            }
            return false
        },
        Severity: "critical",
        Message:  "Memory usage is too high",
    })
}
```

## 6. 监控面板

### 6.1 简单的监控面板

```go
func setupMonitoringDashboard(monitor *NetworkMonitor, latencyMonitor *LatencyMonitor, throughputMonitor *ThroughputMonitor) {
    http.HandleFunc("/monitoring/dashboard", func(w http.ResponseWriter, r *http.Request) {
        // 获取监控数据
        networkStats := monitor.GetConnectionStats()
        latencyStats := latencyMonitor.GetStats()
        throughputStats := throughputMonitor.GetThroughput()
        
        // 生成 HTML 面板
        html := generateDashboardHTML(networkStats, latencyStats, throughputStats)
        
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(html))
    })
}

func generateDashboardHTML(networkStats, latencyStats map[string]interface{}, throughputStats map[string]interface{}) string {
    return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Network Monitoring Dashboard</title>
    <meta http-equiv="refresh" content="5">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .metric { margin: 10px 0; padding: 10px; border: 1px solid #ccc; }
        .warning { background-color: #fff3cd; }
        .critical { background-color: #f8d7da; }
    </style>
</head>
<body>
    <h1>Network Monitoring Dashboard</h1>
    
    <h2>Network Statistics</h2>
    <div class="metric">
        <strong>Active Connections:</strong> %v<br>
        <strong>Total Requests:</strong> %v<br>
        <strong>Error Rate:</strong> %.2f%%<br>
    </div>
    
    <h2>Latency Statistics</h2>
    <div class="metric">
        <strong>Average Latency:</strong> %v<br>
        <strong>Min Latency:</strong> %v<br>
        <strong>Max Latency:</strong> %v<br>
    </div>
    
    <h2>Throughput Statistics</h2>
    <div class="metric">
        <strong>Requests per Second:</strong> %.2f<br>
        <strong>Bytes per Second:</strong> %.2f<br>
    </div>
</body>
</html>
    `, 
    networkStats["active_connections"],
    networkStats["total_requests"],
    networkStats["error_rate"].(float64)*100,
    latencyStats["avg"],
    latencyStats["min"],
    latencyStats["max"],
    throughputStats["requests_per_second"],
    throughputStats["bytes_per_second"])
}
```

## 7. 最佳实践总结

1. **使用内置工具**: 充分利用 pprof 等内置监控工具
2. **自定义监控**: 根据业务需求实现自定义监控指标
3. **实时监控**: 建立实时监控和预警系统
4. **数据收集**: 定期收集和导出监控数据
5. **预警机制**: 设置合理的预警规则和阈值
6. **可视化**: 提供直观的监控面板
7. **性能优化**: 确保监控系统本身不影响性能
8. **数据存储**: 合理存储和管理监控数据

通过使用这些监控工具和技巧，可以有效地监控和优化 Go 程序的网络性能。
