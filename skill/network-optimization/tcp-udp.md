# Go TCP/UDP 网络优化

## 1. TCP 连接优化

### 1.1 TCP 连接池

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

type TCPConnectionPool struct {
    factory    func() (net.Conn, error)
    pool       chan net.Conn
    maxSize    int
    minSize    int
    mu         sync.RWMutex
    stats      PoolStats
}

type PoolStats struct {
    TotalConns    int
    ActiveConns   int
    IdleConns     int
    CreatedConns  int
    DestroyedConns int
}

func NewTCPConnectionPool(factory func() (net.Conn, error), minSize, maxSize int) *TCPConnectionPool {
    pool := &TCPConnectionPool{
        factory: factory,
        pool:    make(chan net.Conn, maxSize),
        maxSize: maxSize,
        minSize: minSize,
    }
    
    // 预创建最小连接数
    for i := 0; i < minSize; i++ {
        conn, err := factory()
        if err == nil {
            pool.pool <- conn
            pool.stats.TotalConns++
            pool.stats.IdleConns++
        }
    }
    
    return pool
}

func (p *TCPConnectionPool) Get() (net.Conn, error) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    select {
    case conn := <-p.pool:
        p.stats.ActiveConns++
        p.stats.IdleConns--
        return conn, nil
    default:
        if p.stats.TotalConns < p.maxSize {
            conn, err := p.factory()
            if err != nil {
                return nil, err
            }
            p.stats.TotalConns++
            p.stats.CreatedConns++
            p.stats.ActiveConns++
            return conn, nil
        }
        
        // 等待可用连接
        select {
        case conn := <-p.pool:
            p.stats.ActiveConns++
            p.stats.IdleConns--
            return conn, nil
        case <-time.After(5 * time.Second):
            return nil, fmt.Errorf("连接池超时")
        }
    }
}

func (p *TCPConnectionPool) Put(conn net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    select {
    case p.pool <- conn:
        p.stats.ActiveConns--
        p.stats.IdleConns++
    default:
        // 池已满，关闭连接
        conn.Close()
        p.stats.TotalConns--
        p.stats.DestroyedConns++
    }
}

func (p *TCPConnectionPool) GetStats() PoolStats {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.stats
}
```

### 1.2 TCP 连接复用

```go
type TCPClient struct {
    pool   *TCPConnectionPool
    addr   string
    timeout time.Duration
}

func NewTCPClient(addr string, timeout time.Duration) *TCPClient {
    factory := func() (net.Conn, error) {
        return net.DialTimeout("tcp", addr, timeout)
    }
    
    return &TCPClient{
        pool:    NewTCPConnectionPool(factory, 5, 20),
        addr:    addr,
        timeout: timeout,
    }
}

func (c *TCPClient) Send(data []byte) ([]byte, error) {
    conn, err := c.pool.Get()
    if err != nil {
        return nil, err
    }
    defer c.pool.Put(conn)
    
    // 设置写入超时
    conn.SetWriteDeadline(time.Now().Add(c.timeout))
    
    // 发送数据
    _, err = conn.Write(data)
    if err != nil {
        return nil, err
    }
    
    // 设置读取超时
    conn.SetReadDeadline(time.Now().Add(c.timeout))
    
    // 读取响应
    buffer := make([]byte, 4096)
    n, err := conn.Read(buffer)
    if err != nil {
        return nil, err
    }
    
    return buffer[:n], nil
}

func (c *TCPClient) Close() {
    // 关闭连接池中的所有连接
    for {
        select {
        case conn := <-c.pool.pool:
            conn.Close()
        default:
            return
        }
    }
}
```

### 1.3 TCP 服务器优化

```go
type TCPServer struct {
    addr     string
    handler  func(net.Conn)
    listener net.Listener
    wg       sync.WaitGroup
    mu       sync.RWMutex
    clients  map[net.Conn]bool
}

func NewTCPServer(addr string, handler func(net.Conn)) *TCPServer {
    return &TCPServer{
        addr:    addr,
        handler: handler,
        clients: make(map[net.Conn]bool),
    }
}

func (s *TCPServer) Start() error {
    listener, err := net.Listen("tcp", s.addr)
    if err != nil {
        return err
    }
    
    s.listener = listener
    
    go s.acceptConnections()
    
    return nil
}

func (s *TCPServer) acceptConnections() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            continue
        }
        
        s.mu.Lock()
        s.clients[conn] = true
        s.mu.Unlock()
        
        s.wg.Add(1)
        go s.handleConnection(conn)
    }
}

func (s *TCPServer) handleConnection(conn net.Conn) {
    defer func() {
        conn.Close()
        s.mu.Lock()
        delete(s.clients, conn)
        s.mu.Unlock()
        s.wg.Done()
    }()
    
    s.handler(conn)
}

func (s *TCPServer) GetClientCount() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.clients)
}

func (s *TCPServer) Stop() error {
    if s.listener != nil {
        s.listener.Close()
    }
    
    s.wg.Wait()
    return nil
}
```

## 2. UDP 优化

### 2.1 UDP 连接池

```go
type UDPConnectionPool struct {
    addr    string
    conn    *net.UDPConn
    mu      sync.RWMutex
    stats   UDPStats
}

type UDPStats struct {
    PacketsSent     int64
    PacketsReceived int64
    BytesSent       int64
    BytesReceived   int64
    Errors          int64
}

func NewUDPConnectionPool(addr string) (*UDPConnectionPool, error) {
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return nil, err
    }
    
    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        return nil, err
    }
    
    return &UDPConnectionPool{
        addr: addr,
        conn: conn,
    }, nil
}

func (p *UDPConnectionPool) Send(data []byte) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    _, err := p.conn.Write(data)
    if err != nil {
        p.stats.Errors++
        return err
    }
    
    p.stats.PacketsSent++
    p.stats.BytesSent += int64(len(data))
    
    return nil
}

func (p *UDPConnectionPool) Receive(buffer []byte) (int, error) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    n, err := p.conn.Read(buffer)
    if err != nil {
        p.stats.Errors++
        return 0, err
    }
    
    p.stats.PacketsReceived++
    p.stats.BytesReceived += int64(n)
    
    return n, nil
}

func (p *UDPConnectionPool) GetStats() UDPStats {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.stats
}

func (p *UDPConnectionPool) Close() error {
    return p.conn.Close()
}
```

### 2.2 UDP 服务器优化

```go
type UDPServer struct {
    addr     string
    handler  func([]byte, *net.UDPAddr)
    conn     *net.UDPConn
    wg       sync.WaitGroup
    mu       sync.RWMutex
    stats    UDPStats
}

func NewUDPServer(addr string, handler func([]byte, *net.UDPAddr)) *UDPServer {
    return &UDPServer{
        addr:    addr,
        handler: handler,
    }
}

func (s *UDPServer) Start() error {
    udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
    if err != nil {
        return err
    }
    
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return err
    }
    
    s.conn = conn
    
    // 启动多个 goroutine 处理 UDP 包
    for i := 0; i < runtime.NumCPU(); i++ {
        s.wg.Add(1)
        go s.handlePackets()
    }
    
    return nil
}

func (s *UDPServer) handlePackets() {
    defer s.wg.Done()
    
    buffer := make([]byte, 4096)
    
    for {
        n, addr, err := s.conn.ReadFromUDP(buffer)
        if err != nil {
            s.mu.Lock()
            s.stats.Errors++
            s.mu.Unlock()
            continue
        }
        
        s.mu.Lock()
        s.stats.PacketsReceived++
        s.stats.BytesReceived += int64(n)
        s.mu.Unlock()
        
        // 复制数据以避免竞态条件
        data := make([]byte, n)
        copy(data, buffer[:n])
        
        s.handler(data, addr)
    }
}

func (s *UDPServer) Send(data []byte, addr *net.UDPAddr) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    _, err := s.conn.WriteToUDP(data, addr)
    if err != nil {
        s.stats.Errors++
        return err
    }
    
    s.stats.PacketsSent++
    s.stats.BytesSent += int64(len(data))
    
    return nil
}

func (s *UDPServer) GetStats() UDPStats {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.stats
}

func (s *UDPServer) Stop() error {
    if s.conn != nil {
        s.conn.Close()
    }
    
    s.wg.Wait()
    return nil
}
```

## 3. 网络缓冲区优化

### 3.1 缓冲区池

```go
type BufferPool struct {
    pool sync.Pool
    size int
}

func NewBufferPool(size int) *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        },
        size: size,
    }
}

func (bp *BufferPool) Get() []byte {
    return bp.pool.Get().([]byte)
}

func (bp *BufferPool) Put(buf []byte) {
    if len(buf) == bp.size {
        bp.pool.Put(buf)
    }
}

// 使用缓冲区池的 TCP 客户端
type OptimizedTCPClient struct {
    pool       *TCPConnectionPool
    bufferPool *BufferPool
    addr       string
    timeout    time.Duration
}

func NewOptimizedTCPClient(addr string, timeout time.Duration) *OptimizedTCPClient {
    factory := func() (net.Conn, error) {
        return net.DialTimeout("tcp", addr, timeout)
    }
    
    return &OptimizedTCPClient{
        pool:       NewTCPConnectionPool(factory, 5, 20),
        bufferPool: NewBufferPool(4096),
        addr:       addr,
        timeout:    timeout,
    }
}

func (c *OptimizedTCPClient) Send(data []byte) ([]byte, error) {
    conn, err := c.pool.Get()
    if err != nil {
        return nil, err
    }
    defer c.pool.Put(conn)
    
    // 获取缓冲区
    buffer := c.bufferPool.Get()
    defer c.bufferPool.Put(buffer)
    
    // 设置写入超时
    conn.SetWriteDeadline(time.Now().Add(c.timeout))
    
    // 发送数据
    _, err = conn.Write(data)
    if err != nil {
        return nil, err
    }
    
    // 设置读取超时
    conn.SetReadDeadline(time.Now().Add(c.timeout))
    
    // 读取响应
    n, err := conn.Read(buffer)
    if err != nil {
        return nil, err
    }
    
    // 返回数据的副本
    result := make([]byte, n)
    copy(result, buffer[:n])
    
    return result, nil
}
```

### 3.2 零拷贝优化

```go
// 使用 io.CopyBuffer 进行零拷贝传输
func zeroCopyTransfer(dst net.Conn, src net.Conn, buffer []byte) (int64, error) {
    return io.CopyBuffer(dst, src, buffer)
}

// 使用 sendfile 系统调用（在支持的系统上）
func sendFile(dst net.Conn, src *os.File, buffer []byte) (int64, error) {
    return io.CopyBuffer(dst, src, buffer)
}
```

## 4. 网络协议优化

### 4.1 自定义协议

```go
type Protocol struct {
    headerSize int
    maxBodySize int
}

func NewProtocol(headerSize, maxBodySize int) *Protocol {
    return &Protocol{
        headerSize:  headerSize,
        maxBodySize: maxBodySize,
    }
}

func (p *Protocol) Encode(data []byte) []byte {
    if len(data) > p.maxBodySize {
        return nil
    }
    
    // 创建消息：头部 + 数据
    message := make([]byte, p.headerSize+len(data))
    
    // 写入头部（包含数据长度）
    binary.BigEndian.PutUint32(message[:4], uint32(len(data)))
    
    // 写入数据
    copy(message[p.headerSize:], data)
    
    return message
}

func (p *Protocol) Decode(data []byte) ([]byte, error) {
    if len(data) < p.headerSize {
        return nil, fmt.Errorf("数据太短")
    }
    
    // 读取头部
    bodySize := binary.BigEndian.Uint32(data[:4])
    
    if bodySize > uint32(p.maxBodySize) {
        return nil, fmt.Errorf("数据太大")
    }
    
    if len(data) < p.headerSize+int(bodySize) {
        return nil, fmt.Errorf("数据不完整")
    }
    
    // 返回数据部分
    return data[p.headerSize : p.headerSize+int(bodySize)], nil
}
```

### 4.2 消息分片

```go
type MessageFragmenter struct {
    maxFragmentSize int
    fragmentID      uint32
    mu              sync.Mutex
}

func NewMessageFragmenter(maxFragmentSize int) *MessageFragmenter {
    return &MessageFragmenter{
        maxFragmentSize: maxFragmentSize,
    }
}

func (mf *MessageFragmenter) Fragment(data []byte) [][]byte {
    mf.mu.Lock()
    defer mf.mu.Unlock()
    
    mf.fragmentID++
    
    if len(data) <= mf.maxFragmentSize {
        // 数据不需要分片
        header := make([]byte, 12)
        binary.BigEndian.PutUint32(header[0:4], mf.fragmentID)
        binary.BigEndian.PutUint32(header[4:8], 1) // 总片段数
        binary.BigEndian.PutUint32(header[8:12], 0) // 当前片段索引
        
        fragment := make([]byte, 12+len(data))
        copy(fragment[:12], header)
        copy(fragment[12:], data)
        
        return [][]byte{fragment}
    }
    
    // 数据需要分片
    totalFragments := (len(data) + mf.maxFragmentSize - 1) / mf.maxFragmentSize
    fragments := make([][]byte, totalFragments)
    
    for i := 0; i < totalFragments; i++ {
        start := i * mf.maxFragmentSize
        end := start + mf.maxFragmentSize
        if end > len(data) {
            end = len(data)
        }
        
        header := make([]byte, 12)
        binary.BigEndian.PutUint32(header[0:4], mf.fragmentID)
        binary.BigEndian.PutUint32(header[4:8], uint32(totalFragments))
        binary.BigEndian.PutUint32(header[8:12], uint32(i))
        
        fragment := make([]byte, 12+end-start)
        copy(fragment[:12], header)
        copy(fragment[12:], data[start:end])
        
        fragments[i] = fragment
    }
    
    return fragments
}
```

## 5. 网络监控

### 5.1 连接监控

```go
type NetworkMonitor struct {
    connections map[net.Conn]ConnectionInfo
    mu          sync.RWMutex
    stats       NetworkStats
}

type ConnectionInfo struct {
    StartTime    time.Time
    BytesSent    int64
    BytesReceived int64
    LastActivity time.Time
}

type NetworkStats struct {
    TotalConnections int
    ActiveConnections int
    TotalBytesSent    int64
    TotalBytesReceived int64
    AverageLatency    time.Duration
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

func (nm *NetworkMonitor) GetStats() NetworkStats {
    nm.mu.RLock()
    defer nm.mu.RUnlock()
    return nm.stats
}
```

### 5.2 性能指标收集

```go
type PerformanceCollector struct {
    metrics map[string]interface{}
    mu      sync.RWMutex
}

func NewPerformanceCollector() *PerformanceCollector {
    return &PerformanceCollector{
        metrics: make(map[string]interface{}),
    }
}

func (pc *PerformanceCollector) SetMetric(key string, value interface{}) {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    
    pc.metrics[key] = value
}

func (pc *PerformanceCollector) GetMetrics() map[string]interface{} {
    pc.mu.RLock()
    defer pc.mu.RUnlock()
    
    result := make(map[string]interface{})
    for k, v := range pc.metrics {
        result[k] = v
    }
    return result
}

func (pc *PerformanceCollector) StartCollection() {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            // 收集系统指标
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            pc.SetMetric("memory_alloc", m.Alloc)
            pc.SetMetric("memory_total", m.TotalAlloc)
            pc.SetMetric("goroutines", runtime.NumGoroutine())
        }
    }()
}
```

## 6. 最佳实践总结

1. **连接池**: 使用连接池减少连接创建开销
2. **缓冲区优化**: 使用缓冲区池减少内存分配
3. **零拷贝**: 使用零拷贝技术提高传输效率
4. **协议优化**: 设计高效的自定义协议
5. **消息分片**: 处理大数据包的分片传输
6. **监控**: 实时监控网络性能指标
7. **错误处理**: 完善的错误处理和重试机制
8. **资源管理**: 及时释放网络资源

通过遵循这些优化策略，可以显著提高 Go 程序的 TCP/UDP 网络性能。
