package skill

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ==================== Channel 通道详解 ====================

// 1. Channel 基本概念

/*
Channel 是 Go 语言中用于 goroutine 间通信的核心机制，它提供了以下特性：
1. 类型安全 - 编译时类型检查
2. 线程安全 - 内置同步机制
3. 阻塞机制 - 自动处理同步
4. 内存安全 - 避免数据竞争

Channel 遵循 "Don't communicate by sharing memory; share memory by communicating" 的设计哲学。
*/

// 2. Channel 基本用法

// 2.1 创建和使用 channel
func BasicChannelExample() {
	// 创建无缓冲 channel
	ch := make(chan int)

	// 启动 goroutine 发送数据
	go func() {
		ch <- 42
		fmt.Println("Sent: 42")
	}()

	// 主 goroutine 接收数据
	value := <-ch
	fmt.Printf("Received: %d\n", value)
}

// 2.2 有缓冲 channel
func BufferedChannelExample() {
	// 创建有缓冲 channel，容量为 3
	ch := make(chan int, 3)

	// 发送数据（不会阻塞，因为缓冲区有空间）
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("Sent 3 values to buffered channel")

	// 接收数据
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n", <-ch)
}

// 2.3 channel 方向性
func ChannelDirectionExample() {
	// 双向 channel
	ch := make(chan int)

	// 发送专用 channel
	sendCh := (chan<- int)(ch)

	// 接收专用 channel
	recvCh := (<-chan int)(ch)

	// 启动发送者
	go func() {
		sendCh <- 100
		fmt.Println("Sent: 100")
	}()

	// 启动接收者
	go func() {
		value := <-recvCh
		fmt.Printf("Received: %d\n", value)
	}()

	time.Sleep(time.Millisecond * 100)
}

// 3. Channel 高级用法

// 3.1 select 语句
func SelectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动数据源
	go func() {
		time.Sleep(time.Millisecond * 100)
		ch1 <- "Data from ch1"
	}()

	go func() {
		time.Sleep(time.Millisecond * 200)
		ch2 <- "Data from ch2"
	}()

	// 使用 select 处理多个 channel
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received from ch1: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received from ch2: %s\n", msg2)
		case <-time.After(time.Second):
			fmt.Println("Timeout")
		}
	}
}

// 3.2 超时控制
func TimeoutExample() {
	ch := make(chan string)

	// 启动长时间操作
	go func() {
		time.Sleep(time.Second * 2)
		ch <- "Operation completed"
	}()

	// 使用 select 实现超时控制
	select {
	case result := <-ch:
		fmt.Printf("Result: %s\n", result)
	case <-time.After(time.Second):
		fmt.Println("Operation timed out")
	}
}

// 3.3 非阻塞操作
func NonBlockingExample() {
	ch := make(chan int)

	// 非阻塞发送
	select {
	case ch <- 42:
		fmt.Println("Sent successfully")
	default:
		fmt.Println("Send would block")
	}

	// 非阻塞接收
	select {
	case value := <-ch:
		fmt.Printf("Received: %d\n", value)
	default:
		fmt.Println("Receive would block")
	}
}

// 4. Channel 模式应用

// 4.1 工作池模式
type WorkerPoolWithChannel struct {
	workers int
	input   chan interface{}
	output  chan interface{}
	done    chan struct{}
	wg      sync.WaitGroup
}

func NewWorkerPoolWithChannel(workers int) *WorkerPoolWithChannel {
	return &WorkerPoolWithChannel{
		workers: workers,
		input:   make(chan interface{}, 100),
		output:  make(chan interface{}, 100),
		done:    make(chan struct{}),
	}
}

func (wp *WorkerPoolWithChannel) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPoolWithChannel) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.input:
			if !ok {
				return
			}
			// 处理任务
			result := wp.processTask(task, id)
			select {
			case wp.output <- result:
			case <-wp.done:
				return
			}
		case <-wp.done:
			return
		}
	}
}

func (wp *WorkerPoolWithChannel) processTask(task interface{}, workerID int) interface{} {
	// 模拟任务处理
	time.Sleep(time.Millisecond * 100)
	return fmt.Sprintf("Task %v processed by worker %d", task, workerID)
}

func (wp *WorkerPoolWithChannel) Submit(task interface{}) {
	select {
	case wp.input <- task:
	case <-wp.done:
	}
}

func (wp *WorkerPoolWithChannel) Results() <-chan interface{} {
	return wp.output
}

func (wp *WorkerPoolWithChannel) Stop() {
	close(wp.done)
	close(wp.input)
	wp.wg.Wait()
	close(wp.output)
}

// 4.2 扇入模式（Fan-in）
func FanInExampleWithChannel() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// 启动数据源
	go func() {
		defer close(ch1)
		for i := 0; i < 5; i++ {
			ch1 <- i
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		defer close(ch2)
		for i := 5; i < 10; i++ {
			ch2 <- i
			time.Sleep(time.Millisecond * 150)
		}
	}()

	// 扇入：合并多个输入 channel
	out := fanInWithChannel(ch1, ch2)

	// 接收合并后的数据
	for value := range out {
		fmt.Printf("Received: %d\n", value)
	}
}

func fanInWithChannel(ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case val, ok := <-ch1:
				if !ok {
					ch1 = nil
				} else {
					out <- val
				}
			case val, ok := <-ch2:
				if !ok {
					ch2 = nil
				} else {
					out <- val
				}
			}

			// 如果两个 channel 都关闭了，退出循环
			if ch1 == nil && ch2 == nil {
				break
			}
		}
	}()

	return out
}

// 4.3 扇出模式（Fan-out）
func FanOutExampleWithChannel() {
	input := make(chan int)
	workers := 3

	// 启动数据源
	go func() {
		defer close(input)
		for i := 0; i < 10; i++ {
			input <- i
			time.Sleep(time.Millisecond * 100)
		}
	}()

	// 扇出：分发到多个工作 goroutine
	fanOutWithChannel(input, workers)
}

func fanOutWithChannel(input <-chan int, workers int) {
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for value := range input {
				fmt.Printf("Worker %d processing: %d\n", workerID, value)
				time.Sleep(time.Millisecond * 50)
			}
		}(i)
	}

	wg.Wait()
}

// 4.4 管道模式（Pipeline）
func PipelineExample() {
	// 创建管道：生成 -> 平方 -> 过滤 -> 输出
	numbers := generate(10)
	squares := square(numbers)
	filtered := filter(squares, func(n int) bool { return n%2 == 0 })
	output := output(filtered)

	// 等待管道完成
	<-output
}

// <-chan int 只读 channel（只能接收数据）
// ch2 chan<- int 只写 channel（只能发送数据）
func generate(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			out <- i
			time.Sleep(time.Millisecond * 50)
		}
	}()
	return out
}
// input <-chan int 只读 channel（只能接收数据）
func square(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for value := range input {
			out <- value * value
		}
	}()
	return out
}

func filter(input <-chan int, predicate func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for value := range input {
			if predicate(value) {
				out <- value
			}
		}
	}()
	return out
}

func output(input <-chan int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for value := range input {
			fmt.Printf("Output: %d\n", value)
		}
	}()
	return done
}

// 5. Channel 与 Context 结合

// 5.1 Context 取消控制
func ContextCancelExampleWithChannel() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	ch := make(chan string)

	go func() {
		select {
		case <-time.After(time.Second):
			ch <- "Operation completed"
		case <-ctx.Done():
			ch <- "Operation cancelled"
		}
	}()

	select {
	case result := <-ch:
		fmt.Printf("Result: %s\n", result)
	case <-ctx.Done():
		fmt.Printf("Context cancelled: %v\n", ctx.Err())
	}
}

// 5.2 优雅关闭
func GracefulShutdownExampleWithChannel() {
	done := make(chan struct{})
	shutdown := make(chan struct{})

	// 启动工作 goroutine
	go func() {
		for {
			select {
			case <-shutdown:
				fmt.Println("Shutting down gracefully...")
				close(done)
				return
			default:
				// 执行正常工作
				fmt.Println("Working...")
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()

	// 模拟信号接收
	go func() {
		time.Sleep(time.Second)
		close(shutdown)
	}()

	// 等待关闭完成
	<-done
	fmt.Println("Application stopped")
}

// 6. Channel 错误处理

// 6.1 错误通道模式
func ErrorChannelExampleWithChannel() {
	resultCh := make(chan string)
	errorCh := make(chan error)

	// 启动可能出错的操作
	go func() {
		if rand.Float32() < 0.5 {
			errorCh <- fmt.Errorf("random error occurred")
			return
		}
		time.Sleep(time.Millisecond * 100)
		resultCh <- "Operation successful"
	}()

	// 处理结果和错误
	select {
	case result := <-resultCh:
		fmt.Printf("Success: %s\n", result)
	case err := <-errorCh:
		fmt.Printf("Error: %v\n", err)
	case <-time.After(time.Second):
		fmt.Println("Operation timed out")
	}
}

// 6.2 带错误的结果类型
type Result struct {
	Value interface{}
	Error error
}

func ResultChannelExample() {
	resultCh := make(chan Result)

	// 启动任务
	go func() {
		if rand.Float32() < 0.3 {
			resultCh <- Result{Error: fmt.Errorf("task failed")}
			return
		}
		time.Sleep(time.Millisecond * 100)
		resultCh <- Result{Value: "Task completed successfully"}
	}()

	// 处理结果
	result := <-resultCh
	if result.Error != nil {
		fmt.Printf("Error: %v\n", result.Error)
	} else {
		fmt.Printf("Success: %v\n", result.Value)
	}
}

// 7. Channel 性能优化

// 7.1 对象池模式
type ObjectPool struct {
	pool chan interface{}
}

func NewObjectPool(size int, factory func() interface{}) *ObjectPool {
	pool := &ObjectPool{
		pool: make(chan interface{}, size),
	}

	// 预填充对象池
	for i := 0; i < size; i++ {
		pool.pool <- factory()
	}

	return pool
}

func (op *ObjectPool) Get() interface{} {
	select {
	case obj := <-op.pool:
		return obj
	default:
		// 池为空，创建新对象
		return make([]byte, 1024)
	}
}

func (op *ObjectPool) Put(obj interface{}) {
	select {
	case op.pool <- obj:
		// 成功放回池中
	default:
		// 池已满，丢弃对象
	}
}

// 7.2 批量处理
func BatchProcessingExample() {
	input := make(chan int, 100)
	batchSize := 10

	// 启动数据源
	go func() {
		defer close(input)
		for i := 0; i < 50; i++ {
			input <- i
		}
	}()

	// 批量处理
	go func() {
		batch := make([]int, 0, batchSize)
		for value := range input {
			batch = append(batch, value)
			if len(batch) >= batchSize {
				processBatch(batch)
				batch = batch[:0] // 清空切片，保持容量
			}
		}
		// 处理剩余数据
		if len(batch) > 0 {
			processBatch(batch)
		}
	}()

	time.Sleep(time.Second)
}

func processBatch(batch []int) {
	fmt.Printf("Processing batch of %d items: %v\n", len(batch), batch)
	time.Sleep(time.Millisecond * 100)
}

// 8. Channel 测试和调试

// 8.1 Channel 行为测试
func ChannelBehaviorTest() {
	ch := make(chan int, 3)

	// 测试发送
	fmt.Printf("Channel capacity: %d\n", cap(ch))
	fmt.Printf("Channel length: %d\n", len(ch))

	// 发送数据
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("After sending 3 items, length: %d\n", len(ch))

	// 接收数据
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("After receiving 1 item, length: %d\n", len(ch))
}

// 8.2 Channel 性能基准测试
func ChannelBenchmarkExample() {
	ch := make(chan int, 1000)

	// 填充 channel
	for i := 0; i < 1000; i++ {
		ch <- i
	}

	// 基准测试
	start := time.Now()
	for i := 0; i < 1000; i++ {
		<-ch
	}
	duration := time.Since(start)

	fmt.Printf("Processed 1000 items in %v\n", duration)
}

// 9. Channel 常见陷阱和解决方案

// 9.1 死锁问题
func DeadlockExample() {
	ch := make(chan int)

	// 这会导致死锁，因为没有接收者
	// ch <- 42

	// 解决方案：在另一个 goroutine 中发送
	go func() {
		ch <- 42
	}()

	value := <-ch
	fmt.Printf("Received: %d\n", value)
}

// 9.2 内存泄漏
func MemoryLeakExample() {
	ch := make(chan int)

	// 启动发送者
	go func() {
		for i := 0; i < 1000; i++ {
			ch <- i
		}
		close(ch) // 重要：关闭 channel
	}()

	// 接收所有数据
	for value := range ch {
		fmt.Printf("Received: %d\n", value)
	}
}

// 9.3 竞态条件
func RaceConditionExample() {
	var counter int
	ch := make(chan int)

	// 启动多个 goroutine 增加计数器
	for i := 0; i < 100; i++ {
		go func() {
			counter++
			ch <- counter
		}()
	}

	// 接收结果
	for i := 0; i < 100; i++ {
		<-ch
	}

	fmt.Printf("Final counter: %d\n", counter)
}

// 10. Channel 最佳实践

// 10.1 关闭 channel 的规则
func ChannelCloseRules() {
	ch := make(chan int)

	// 规则1：只有发送者应该关闭 channel
	go func() {
		defer close(ch) // 在发送者中关闭
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()

	// 规则2：接收者不应该关闭 channel
	for value := range ch {
		fmt.Printf("Received: %d\n", value)
	}
}

// 10.2 使用 select 避免阻塞
func NonBlockingSelectExample() {
	ch := make(chan int)

	// 非阻塞发送
	select {
	case ch <- 42:
		fmt.Println("Sent successfully")
	default:
		fmt.Println("Send would block")
	}

	// 非阻塞接收
	select {
	case value := <-ch:
		fmt.Printf("Received: %d\n", value)
	default:
		fmt.Println("Receive would block")
	}
}

// 10.3 使用缓冲 channel 提高性能
func BufferedChannelPerformanceExample() {
	// 无缓冲 channel
	unbuffered := make(chan int)

	// 有缓冲 channel
	buffered := make(chan int, 10)

	// 测试无缓冲 channel
	select {
	case unbuffered <- 1:
		fmt.Println("Sent to unbuffered channel")
	default:
		fmt.Println("Cannot send to unbuffered channel (no receiver)")
	}

	// 测试有缓冲 channel
	select {
	case buffered <- 1:
		fmt.Println("Sent to buffered channel")
	default:
		fmt.Println("Cannot send to buffered channel")
	}
}

// 11. 总结

func ChannelSummary() {
	fmt.Println("=== Channel 总结 ===")
	fmt.Println("1. Channel 是 Go 并发编程的核心")
	fmt.Println("2. 提供了类型安全和线程安全的通信机制")
	fmt.Println("3. 支持阻塞和非阻塞操作")
	fmt.Println("4. 可以与 select 和 context 结合使用")
	fmt.Println("5. 支持多种并发模式：工作池、扇入、扇出、管道等")
	fmt.Println("6. 需要正确关闭 channel 避免内存泄漏")
	fmt.Println("7. 使用缓冲 channel 可以提高性能")
	fmt.Println("8. 遵循 'Don't communicate by sharing memory' 的设计哲学")
}
