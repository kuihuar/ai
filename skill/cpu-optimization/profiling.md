# Go CPU 性能分析

## 1. pprof 工具使用

### 1.1 启用 pprof

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "time"
)

func main() {
    // 启用 pprof HTTP 服务
    go func() {
        log.Println("pprof server started on :6060")
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 你的程序逻辑
    runApplication()
}

// 手动触发 CPU 分析
func startCPUProfile() {
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    if err := pprof.StartCPUProfile(f); err != nil {
        log.Fatal(err)
    }
    defer pprof.StopCPUProfile()
    
    // 运行需要分析的代码
    runCPUIntensiveTask()
}
```

### 1.2 命令行分析

```bash
# 启动 CPU 分析（30秒）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 生成火焰图
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile

# 查看 CPU 使用率最高的函数
go tool pprof -top http://localhost:6060/debug/pprof/profile

# 查看特定函数的调用栈
go tool pprof -list=functionName http://localhost:6060/debug/pprof/profile

# 生成调用图
go tool pprof -web http://localhost:6060/debug/pprof/profile

# 生成文本报告
go tool pprof -text http://localhost:6060/debug/pprof/profile

# 生成树形图
go tool pprof -tree http://localhost:6060/debug/pprof/profile
```

### 1.3 交互式分析

```bash
# 进入交互模式
go tool pprof http://localhost:6060/debug/pprof/profile

# 常用命令
(pprof) top10          # 显示前10个最耗CPU的函数
(pprof) list function  # 显示特定函数的代码
(pprof) web            # 在浏览器中打开调用图
(pprof) svg            # 生成SVG格式的调用图
(pprof) png            # 生成PNG格式的调用图
(pprof) pdf            # 生成PDF格式的调用图
(pprof) peek function  # 显示函数的调用者和被调用者
(pprof) disasm function # 显示函数的汇编代码
(pprof) weblist function # 在浏览器中显示函数代码
```

## 2. 性能分析示例

### 2.1 计算密集型任务分析

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

// 计算密集型任务
func cpuIntensiveTask() {
    sum := 0
    for i := 0; i < 1000000; i++ {
        sum += i * i
    }
    fmt.Printf("Sum: %d\n", sum)
}

// 递归计算斐波那契数列
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// 矩阵乘法
func matrixMultiply(a, b [][]int) [][]int {
    n := len(a)
    result := make([][]int, n)
    for i := range result {
        result[i] = make([]int, n)
    }
    
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            for k := 0; k < n; k++ {
                result[i][j] += a[i][k] * b[k][j]
            }
        }
    }
    
    return result
}

func main() {
    // 启用 pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 运行计算密集型任务
    go func() {
        for {
            cpuIntensiveTask()
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // 运行递归任务
    go func() {
        for {
            result := fibonacci(30)
            fmt.Printf("Fibonacci(30): %d\n", result)
            time.Sleep(200 * time.Millisecond)
        }
    }()
    
    // 运行矩阵乘法
    go func() {
        size := 100
        a := make([][]int, size)
        b := make([][]int, size)
        
        for i := 0; i < size; i++ {
            a[i] = make([]int, size)
            b[i] = make([]int, size)
            for j := 0; j < size; j++ {
                a[i][j] = i + j
                b[i][j] = i - j
            }
        }
        
        for {
            result := matrixMultiply(a, b)
            _ = result
            time.Sleep(500 * time.Millisecond)
        }
    }()
    
    select {}
}
```

### 2.2 I/O 密集型任务分析

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "time"
)

// 文件 I/O 密集型任务
func fileIOTask() {
    // 创建临时文件
    file, err := os.CreateTemp("", "cpu_profile_test")
    if err != nil {
        log.Fatal(err)
    }
    defer os.Remove(file.Name())
    defer file.Close()
    
    // 写入大量数据
    data := make([]byte, 1024*1024) // 1MB
    for i := range data {
        data[i] = byte(i % 256)
    }
    
    for i := 0; i < 100; i++ {
        _, err := file.Write(data)
        if err != nil {
            log.Fatal(err)
        }
    }
    
    // 读取数据
    file.Seek(0, 0)
    buffer := make([]byte, 1024)
    for {
        _, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
    }
}

// 网络 I/O 密集型任务
func networkIOTask() {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    
    resp, err := client.Get("https://httpbin.org/delay/1")
    if err != nil {
        log.Printf("网络请求失败: %v", err)
        return
    }
    defer resp.Body.Close()
    
    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("读取响应失败: %v", err)
        return
    }
    
    fmt.Printf("响应大小: %d bytes\n", len(body))
}

func main() {
    // 启用 pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 运行文件 I/O 任务
    go func() {
        for {
            fileIOTask()
            time.Sleep(1 * time.Second)
        }
    }()
    
    // 运行网络 I/O 任务
    go func() {
        for {
            networkIOTask()
            time.Sleep(2 * time.Second)
        }
    }()
    
    select {}
}
```

## 3. 性能分析技巧

### 3.1 识别热点函数

```go
// 使用 pprof 标签
func processData(data []int) {
    defer pprof.SetGoroutineLabels(pprof.Labels("function", "processData"))
    
    for i, v := range data {
        if v > 100 {
            // 热点代码
            result := expensiveCalculation(v)
            data[i] = result
        }
    }
}

func expensiveCalculation(n int) int {
    defer pprof.SetGoroutineLabels(pprof.Labels("function", "expensiveCalculation"))
    
    // 模拟复杂计算
    result := 0
    for i := 0; i < n; i++ {
        result += i * i
    }
    return result
}
```

### 3.2 分析特定时间段

```go
func analyzeTimeRange() {
    // 开始分析
    pprof.StartCPUProfile(os.Stdout)
    defer pprof.StopCPUProfile()
    
    // 运行需要分析的代码
    start := time.Now()
    for time.Since(start) < 10*time.Second {
        cpuIntensiveTask()
    }
}
```

### 3.3 比较不同版本

```bash
# 生成基线版本的分析文件
go run -cpuprofile=baseline.prof main.go

# 生成优化版本的分析文件
go run -cpuprofile=optimized.prof main.go

# 比较两个版本
go tool pprof -base=baseline.prof optimized.prof
```

## 4. 高级分析技巧

### 4.1 使用 trace 工具

```go
package main

import (
    "context"
    "log"
    "os"
    "runtime/trace"
    "time"
)

func main() {
    // 创建 trace 文件
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    // 开始 trace
    if err := trace.Start(f); err != nil {
        log.Fatal(err)
    }
    defer trace.Stop()
    
    // 运行需要分析的代码
    ctx := context.Background()
    runWithTrace(ctx)
}

func runWithTrace(ctx context.Context) {
    // 创建多个 goroutine
    for i := 0; i < 10; i++ {
        go func(id int) {
            for j := 0; j < 1000; j++ {
                // 模拟工作
                time.Sleep(1 * time.Millisecond)
            }
        }(i)
    }
    
    time.Sleep(5 * time.Second)
}
```

### 4.2 分析 trace 文件

```bash
# 启动 trace 分析服务器
go tool trace trace.out

# 在浏览器中打开 http://localhost:8000
```

### 4.3 自定义分析器

```go
package main

import (
    "fmt"
    "log"
    "runtime"
    "time"
)

type CPUProfiler struct {
    samples []CPUSample
    enabled bool
}

type CPUSample struct {
    Timestamp time.Time
    CPUUsage  float64
    Goroutines int
}

func NewCPUProfiler() *CPUProfiler {
    return &CPUProfiler{
        samples: make([]CPUSample, 0),
        enabled: true,
    }
}

func (p *CPUProfiler) Start() {
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for range ticker.C {
            if p.enabled {
                p.sample()
            }
        }
    }()
}

func (p *CPUProfiler) sample() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    sample := CPUSample{
        Timestamp:   time.Now(),
        CPUUsage:    getCPUUsage(),
        Goroutines:  runtime.NumGoroutine(),
    }
    
    p.samples = append(p.samples, sample)
    
    // 保持最近1000个样本
    if len(p.samples) > 1000 {
        p.samples = p.samples[1:]
    }
}

func (p *CPUProfiler) GetStats() []CPUSample {
    return p.samples
}

func (p *CPUProfiler) Stop() {
    p.enabled = false
}

func getCPUUsage() float64 {
    // 简化的 CPU 使用率计算
    // 实际实现需要更复杂的逻辑
    return float64(runtime.NumGoroutine()) * 0.1
}

func main() {
    profiler := NewCPUProfiler()
    profiler.Start()
    defer profiler.Stop()
    
    // 运行你的程序
    runApplication()
    
    // 分析结果
    stats := profiler.GetStats()
    for _, sample := range stats {
        fmt.Printf("时间: %v, CPU使用率: %.2f%%, Goroutines: %d\n",
            sample.Timestamp, sample.CPUUsage, sample.Goroutines)
    }
}
```

## 5. 性能分析最佳实践

### 5.1 分析前准备

1. **确定分析目标**: 明确要优化的性能指标
2. **选择合适的时间**: 在系统负载稳定时进行分析
3. **准备测试数据**: 使用真实或接近真实的数据
4. **设置基准**: 记录优化前的性能指标

### 5.2 分析过程

1. **多次采样**: 进行多次分析确保结果稳定
2. **不同负载**: 在不同负载下进行分析
3. **对比分析**: 与历史数据或基准进行对比
4. **深入分析**: 不仅看热点函数，还要分析调用关系

### 5.3 分析后优化

1. **优先优化热点**: 先优化占用CPU最多的函数
2. **验证优化效果**: 优化后重新分析验证效果
3. **持续监控**: 建立持续的性能监控机制
4. **文档记录**: 记录优化过程和结果

## 6. 常见问题排查

### 6.1 分析结果不准确

```go
// 确保分析时间足够长
func accurateProfile() {
    // 至少运行30秒
    pprof.StartCPUProfile(os.Stdout)
    defer pprof.StopCPUProfile()
    
    // 运行足够长的时间
    time.Sleep(30 * time.Second)
}
```

### 6.2 分析文件过大

```go
// 限制分析时间
func limitedProfile() {
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // 只分析10秒
    time.Sleep(10 * time.Second)
}
```

### 6.3 分析特定函数

```go
// 使用标签分析特定函数
func analyzeSpecificFunction() {
    defer pprof.SetGoroutineLabels(pprof.Labels("function", "target"))
    
    // 目标函数
    targetFunction()
}
```

通过合理使用这些分析工具和技巧，可以有效地识别和解决 Go 程序中的 CPU 性能问题。
