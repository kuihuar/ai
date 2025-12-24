### Channel 知识体系整理

#### 基础
1. 本质：类型安全的管道，用于 Goroutine 间通信（FIFO队列）
2. 声明：`var ch chan int`（无缓冲）或 `ch := make(chan int, 10)`（有缓冲）
3. 操作：发送 `ch <- 1`，接收 `x := <-ch`，关闭 `close(ch)`
4. 特性：
   - 无缓冲：发送和接收会阻塞，直到配对操作
   - 有缓冲：发送和接收不会阻塞，直到队列满或空
5. 阻塞与非阻塞
   - 阻塞：无配对操作，会导致 Goroutine 阻塞
   - 非阻塞：有配对操作，不会阻塞 Goroutine
6. 关闭 Channel
   - 发送方：`close(ch)`
   - 接收方：`x, ok := <-ch`，`ok` 为 false 表示 Channel 已关闭
7. 单向 Channel
   - `chan<- int`：只写 Channel，只能发送数据
   - `<-chan int`：只读 Channel，只能接收数据
8. 选择器（select）
   - 阻塞等待多个 Channel 操作
   - 随机选择一个执行，多个就绪时会随机选择
   - 示例：
     ```go
     select {
     case x := <-ch1:
         fmt.Println("ch1:", x)
     case ch2 <- 2:
         fmt.Println("ch2 sent 2")
     default:
         fmt.Println("no operation")
     }
     ```
9.  无缓冲 vs 缓冲通道
| 特性 | 无缓冲通道 | 缓冲通道 |
| --- | --- | --- |
| 同步性 | 发送和接收必须同时就绪 | 发送方不阻塞直到缓冲区满 |
| 使用场景 | 强同步控制（如信号通知） | 解耦生产消费速率 |
| 示例 | `make(chan struct{})` | `make(chan int, 10)` |

#### 应用场景
1. 生产者-消费者模型
2. 工作池（Worker Pool）
   - 多个工作 Goroutine 并发执行任务
   - 任务队列缓冲多个任务，避免过载
3. 任务分发与结果收集
   - 多个生产者 Goroutine 分发任务
   - 多个消费者 Goroutine 收集结果
4. 事件通知与响应
   - 一个事件 Goroutine 通知多个等待者
   - 多个接收者 Goroutine 响应事件
5. 超时控制
   - 避免 Goroutine 阻塞，设置超时时间
   - 使用 `select` 结合 `time.After` 实现
#### 高级特性
1. 关闭通道的规则
- 关闭后：
    - 仍可接收剩余数据，收到零值后通过val, ok := <-ch的ok判断是否关闭。
    - 向已关闭通道发送数据会触发panic。

- 最佳实践：
    ```go
    defer close(ch) // 确保通道关闭   
    ```
2. 单向通道（类型约束）
- 用途：限制通道的发送或接收权限，增强代码安全性。

```go
func producer(ch chan<- int) { // 只允许发送
    ch <- 1
}
func consumer(ch <-chan int) { // 只允许接收
    val := <-ch
}   
```
3. Select 多路复用
- 作用：监听多个通道操作，避免阻塞。
- 语法：
```go
select {
case msg := <-ch1:
    fmt.Println(msg)
case ch2 <- 42:
    fmt.Println("sent")
default:
    fmt.Println("no activity")
}
```
- 典型应用：
    - 超时控制（time.After）。
    - 非阻塞检查（default分支）。
    - 多个通道选择（随机）。
    - 事件通知与响应（select + case）。
4. Channel 底层结构（面试深挖）
- 数据结构：
    - 循环队列（有缓冲）或单向链表（无缓冲）
    - 包含缓冲数据、接收指针、发送指针、锁等
- 操作：
    - 发送：阻塞直到队列有空位
    - 接收：阻塞直到队列有数据
- 注意事项：
    - 关闭无缓冲通道会触发所有阻塞的接收操作
- 关键点：
    - 通道通过lock实现线程安全。
    - 缓冲区为环形队列，避免频繁内存分配。

#### Channel 面试高频问题
1. 基础问题
- Q1：Channel 和 Mutex 如何选择？
    - 答：Channel 用于 Goroutine 间通信/同步，Mutex 用于保护共享资源。
    - 场景：
    - 数据传递 → Channel
    - 状态保护 → Mutex

- Q2：向已关闭的 Channel 发送数据会发生什么？
    - 答：触发 panic。应在发送端调用 close，且接收端通过 val, ok := <-ch 检测关闭。
2. 设计模式问题
- Q3：如何用 Channel 实现 Worker Pool？
    - 答：使用有缓冲通道和 Goroutine 池模式。
    - 示例：
```go
    func worker(id int, tasks <-chan int, results chan<- int) {
        for task := range tasks {
            results <- task * 2
        }
    }
    func main() {
    tasks := make(chan int, 10)
        results := make(chan int, 10)
        // 启动3个Worker
        for i := 0; i < 3; i++ {
            go worker(tasks, results)
        }
        // 发送任务
        for i := 0; i < 5; i++ {
            tasks <- i
        }
        close(tasks)
        // 收集结果
        for i := 0; i < 5; i++ {
            fmt.Println(<-results)
        }
    }
    ```
- Q4: 如何用 Channel 实现 Pub-Sub 模式？
    - 答：使用 map 存储订阅关系，Goroutine 池处理并发发布。
    - 示例:
```go
    type PubSub struct {
        subs map[string][]chan string
        mu   sync.RWMutex
    }
    func (ps *PubSub) Publish(topic string, msg string) {
        ps.mu.RLock()
        defer ps.mu.RUnlock()
        for _, ch := range ps.subs[topic] {
            ch <- msg
        }
    }
    func (ps *PubSub) Subscribe(topic string) <-chan string {
        ps.mu.Lock()
        defer ps.mu.Unlock()
        ch := make(chan string)
        ps.subs[topic] = append(ps.subs[topic], ch)
        return ch
    }
```
- Q5: 如何用 Channel 实现带超时的操作？
    - 答：使用 select 结合 time.After 实现。
    - 示例：
```go
    ch := make(chan string)
    timeout := time.After(100 * time.Millisecond)
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    case <-timeout:
        fmt.Println("Timeout")
    }  
```
- Q6:   缓冲的队列？
    - 答：使用有缓冲通道。
    - 示例：
    ```go
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    // 阻塞，队列满
    ch <- 4
    ```
    - 接收端：
    ```go
    // 阻塞，队列空
    val := <-ch
    ```
- Q7:  用 Channel 实现一个并发安全的计数器
    - 答：使用原子操作（sync/atomic）。
    - 示例：
    ```go
    type Counter struct {
        value int64
    }
    func (c *Counter) Inc() {
        atomic.AddInt64(&c.value, 1)
    }
    func (c *Counter) Get() int64 {
        return atomic.LoadInt64(&c.value)
    }


    type Counter struct {
    ch  chan int
    val int
    }
    func NewCounter() *Counter {
        c := &Counter{ch: make(chan int)}
        go c.run() // 后台Goroutine处理更新
        return c
    }
    func (c *Counter) run() {
        for delta := range c.ch {
            c.val += delta
        }
    }
    func (c *Counter) Add(delta int) {
        c.ch <- delta
    }
    func (c *Counter) Value() int {
        return c.val
    }
    ```

3. 陷阱与调试
- Q5: Channel 导致 Goroutine 泄漏的场景？
    - 答：Goroutine 阻塞在无接收者的 Channel 上。
    - 解决：使用 context.Context 或带超时的 select。
- Q6: 如何诊断 Channel 死锁？
    - 答：go run -race 检测数据竞争，pprof 分析阻塞的 Goroutine。

四、性能优化技巧
1. 缓冲大小选择：
    - 计算生产消费速率差，避免频繁阻塞。
    - 合理设置缓冲大小，避免频繁切换上下文。
2. 结构体传递：
    - 传递指针而非大结构体（chan *Task）。
3. 替代方案：
    - 高吞吐场景考虑 sync.Pool 或原子操作。
    - 避免频繁内存分配，使用 sync.Pool 复用对象。
    - 避免锁，使用原子操作代替 Mutex。
