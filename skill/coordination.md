# Go语言协调机制全面学习指南

## 概述

Go语言提供了多种协调机制来支持goroutine之间的同步和通信，包括条件变量、等待组、屏障等。这些机制不是锁，而是基于锁实现的协调工具。

## 1. 条件变量 (sync.Cond)

### 1.1 基本概念

**特点：**
- 基于互斥锁实现的等待-通知机制
- 允许goroutine等待特定条件
- 支持单个或多个goroutine的唤醒

**核心价值：**
- 解耦等待和通知
- 避免无效的轮询消耗资源
- 实现复杂的同步逻辑

### 1.2 典型应用场景

**生产者消费者模式：**
- 协调生产和消费速度
- 管理缓冲区状态
- 实现异步处理

**资源池管理：**
- 管理有限资源的分配和回收
- 控制并发数量
- 实现连接池

**任务队列：**
- 实现线程池和任务调度
- 管理任务状态
- 协调工作流程

### 1.3 使用示例

**生产者消费者：**
```go
type Buffer struct {
    mu      sync.Mutex
    cond    *sync.Cond
    items   []interface{}
    size    int
    maxSize int
}

func (b *Buffer) Produce(item interface{}) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    for b.size >= b.maxSize {
        b.cond.Wait() // 等待缓冲区有空间
    }
    
    b.items = append(b.items, item)
    b.size++
    b.cond.Signal() // 通知消费者
}

func (b *Buffer) Consume() interface{} {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    for b.size == 0 {
        b.cond.Wait() // 等待缓冲区有数据
    }
    
    item := b.items[0]
    b.items = b.items[1:]
    b.size--
    b.cond.Signal() // 通知生产者
    return item
}
```

## 2. 等待组 (sync.WaitGroup)

### 2.1 基本概念

**特点：**
- 等待多个goroutine完成任务
- 基于原子操作实现
- 开销很低

**使用场景：**
- 主协程等待子协程完成
- 任务协调
- 资源清理

### 2.2 使用示例

**基本用法：**
```go
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            work(id)
        }(i)
    }
    
    wg.Wait() // 等待所有goroutine完成
    fmt.Println("所有任务完成")
}
```

**错误处理：**
```go
func processWithErrorHandling() error {
    var wg sync.WaitGroup
    errChan := make(chan error, 1)
    
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            if err := work(id); err != nil {
                select {
                case errChan <- err:
                default:
                }
            }
        }(i)
    }
    
    wg.Wait()
    close(errChan)
    
    if err := <-errChan; err != nil {
        return err
    }
    return nil
}
```

## 3. 屏障 (Barrier)

### 3.1 基本概念

**特点：**
- 多个goroutine在特定点汇合
- 同步后所有goroutine继续执行
- 支持分阶段同步

**使用场景：**
- 分阶段处理
- 循环同步
- 并行算法

### 3.2 实现示例

**基本屏障：**
```go
type Barrier struct {
    mu       sync.Mutex
    cond     *sync.Cond
    count    int
    parties  int
    phase    int
}

func NewBarrier(parties int) *Barrier {
    b := &Barrier{parties: parties}
    b.cond = sync.NewCond(&b.mu)
    return b
}

func (b *Barrier) Await() {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    phase := b.phase
    b.count++
    
    if b.count == b.parties {
        // 最后一个到达的goroutine
        b.count = 0
        b.phase++
        b.cond.Broadcast()
    } else {
        // 等待其他goroutine
        for phase == b.phase {
            b.cond.Wait()
        }
    }
}
```

**使用示例：**
```go
func parallelSort(data []int) {
    barrier := NewBarrier(2)
    
    go func() {
        sortFirstHalf(data[:len(data)/2])
        barrier.Await() // 等待另一半排序完成
        merge(data)     // 合并结果
    }()
    
    go func() {
        sortSecondHalf(data[len(data)/2:])
        barrier.Await() // 等待前一半排序完成
        merge(data)     // 合并结果
    }()
}
```

## 4. 选择指导原则

### 4.1 选择条件变量当：
- 需要等待特定条件
- 实现生产者消费者模式
- 复杂的等待-通知逻辑

### 4.2 选择等待组当：
- 等待多个goroutine完成
- 简单的任务协调
- 主协程等待子协程

### 4.3 选择屏障当：
- 需要分阶段同步
- 多个goroutine在特定点汇合
- 循环的同步需求

## 5. 性能考虑

| 协调机制 | 性能特点 | 适用场景 |
|----------|----------|----------|
| 条件变量 | 较高开销，但避免忙等待 | 复杂同步 |
| 等待组 | 很低开销，基于原子操作 | 任务协调 |
| 屏障 | 中等开销，适合批量同步 | 分阶段处理 |

## 6. 总结

协调机制的核心价值在于：

1. **条件变量**：解决"等待"问题，不关心"同步"
2. **等待组**：解决"单向等待"问题，不关心"互相等待"
3. **屏障**：解决"协同节奏"问题，不关心"单向等待"

通过合理选择和组合使用这些协调机制，可以构建出高效、安全、可维护的并发程序。
