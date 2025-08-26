package skill

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"time"
)

// ==================== select 多路复用详解 ====================

// 1. select 基本概念

/*
select 是 Go 语言中用于多路复用的关键字，它允许在多个 channel 操作之间进行非阻塞的选择。
select 语句会随机选择一个可执行的 case 执行，如果没有 case 可执行，则执行 default（如果存在）。

核心特性：
1. 非阻塞 - 不会因为某个 channel 阻塞而影响其他 channel
2. 随机选择 - 当多个 case 同时可执行时，随机选择一个
3. 超时控制 - 结合 time.After 实现超时机制
4. 取消控制 - 结合 context 实现取消机制
*/

// 2. select 基本语法

// 2.1 基本 select 语句
func BasicSelectExample() {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// 发送数据到 channel
	ch1 <- "Hello from ch1"
	ch2 <- "Hello from ch2"

	// 使用 select 接收数据
	select {
	case msg1 := <-ch1:
		fmt.Printf("Received from ch1: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("Received from ch2: %s\n", msg2)
	default:
		fmt.Println("No message available")
	}
}

// 2.2 select 与 default
func SelectWithDefaultExample() {
	ch := make(chan string)

	// 尝试从空的 channel 接收数据
	select {
	case msg := <-ch:
		fmt.Printf("Received: %s\n", msg)
	default:
		fmt.Println("Channel is empty, using default")
	}

	// 尝试向无缓冲的 channel 发送数据（没有接收者）
	select {
	case ch <- "message":
		fmt.Println("Message sent successfully")
	default:
		fmt.Println("No receiver available, using default")
	}
}

// 2.3 select 阻塞行为
func SelectBlockingExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动 goroutine 发送数据
	go func() {
		time.Sleep(time.Millisecond * 100)
		ch1 <- "Data from ch1"
	}()

	go func() {
		time.Sleep(time.Millisecond * 200)
		ch2 <- "Data from ch2"
	}()

	// select 会阻塞直到有数据可接收
	select {
	case msg1 := <-ch1:
		fmt.Printf("Received from ch1: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("Received from ch2: %s\n", msg2)
	}
}

// 3. select 高级用法

// 3.1 超时控制
func SelectTimeoutExample() {
	ch := make(chan string)

	// 启动 goroutine 模拟长时间操作
	go func() {
		time.Sleep(time.Second * 2)
		ch <- "Operation completed"
	}()

	// 使用 select 实现超时控制
	select {
	case result := <-ch:
		fmt.Printf("Operation result: %s\n", result)
	case <-time.After(time.Second):
		fmt.Println("Operation timed out")
	}
}

// 3.2 超时控制（可重复使用）
func SelectTimeoutReusableExample() {
	ch := make(chan string)

	// 创建定时器
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	go func() {
		time.Sleep(time.Millisecond * 500)
		ch <- "Operation completed"
	}()

	select {
	case result := <-ch:
		fmt.Printf("Operation result: %s\n", result)
	case <-timer.C:
		fmt.Println("Operation timed out")
	}
}

// 3.3 结合 Context 的取消控制
// 核心价值：解决 “异步操作超时” 问题
func SelectWithContextExample() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := make(chan string)

	go func() {
		time.Sleep(time.Second * 2)
		select {
		case ch <- "Operation completed":
		case <-ctx.Done():
			return
		}
	}()

	select {
	case result := <-ch:
		fmt.Printf("Operation result: %s\n", result)
	case <-ctx.Done():
		fmt.Printf("Operation cancelled: %v\n", ctx.Err())
	}
}

// 4. select 在并发模式中的应用

// 4.1 扇入模式（Fan-in）
func FanInExample() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// 启动数据源
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
			time.Sleep(time.Millisecond * 100)
		}
		close(ch1)
	}()

	go func() {
		for i := 10; i < 15; i++ {
			ch2 <- i
			time.Sleep(time.Millisecond * 150)
		}
		close(ch2)
	}()

	// 使用 select 合并数据流
	for {
		select {
		case val, ok := <-ch1:
			if !ok {
				ch1 = nil // 标记 channel 已关闭
			} else {
				fmt.Printf("Received from ch1: %d\n", val)
			}
		case val, ok := <-ch2:
			if !ok {
				ch2 = nil // 标记 channel 已关闭
			} else {
				fmt.Printf("Received from ch2: %d\n", val)
			}
		}

		// 如果两个 channel 都关闭了，退出循环
		if ch1 == nil && ch2 == nil {
			break
		}
	}
	fmt.Println("Fan-in completed")
}

// 4.2 扇出模式（Fan-out）
func FanOutExample() {
	input := make(chan int)

	// 启动多个工作 goroutine
	for i := 0; i < 3; i++ {
		go func(workerID int) {
			for val := range input {
				fmt.Printf("Worker %d processing: %d\n", workerID, val)
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	// 发送数据
	for i := 0; i < 10; i++ {
		input <- i
	}
	close(input)

	// 等待工作完成
	time.Sleep(time.Second)
	fmt.Println("Fan-out completed")
}

// 4.3 工作池模式
type WorkerPool struct {
	workers int
	input   chan interface{}
	output  chan interface{}
	done    chan struct{}
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		input:   make(chan interface{}, 100),
		output:  make(chan interface{}, 100),
		done:    make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// 启动输出处理
	go wp.outputHandler()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.input:
			if !ok {
				return
			}
			// 处理任务
			result := wp.processTask(task)
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

func (wp *WorkerPool) processTask(task interface{}) interface{} {
	// 模拟任务处理
	time.Sleep(time.Millisecond * 50)
	return fmt.Sprintf("Processed: %v", task)
}

func (wp *WorkerPool) outputHandler() {
	for {
		select {
		case result, ok := <-wp.output:
			if !ok {
				return
			}
			fmt.Printf("Output: %v\n", result)
		case <-wp.done:
			return
		}
	}
}

func (wp *WorkerPool) Submit(task interface{}) {
	select {
	case wp.input <- task:
	case <-wp.done:
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.done)
	close(wp.input)
	wp.wg.Wait()
	close(wp.output)
}

// 5. select 与错误处理

// 5.1 错误通道模式
func ErrorChannelExample() {
	resultCh := make(chan string)
	errorCh := make(chan error)

	// 启动工作 goroutine
	go func() {
		// 模拟可能出错的操作
		if rand.Float32() < 0.5 {
			errorCh <- fmt.Errorf("random error occurred")
			return
		}
		time.Sleep(time.Millisecond * 100)
		resultCh <- "Operation successful"
	}()

	// 使用 select 处理结果和错误
	select {
	case result := <-resultCh:
		fmt.Printf("Success: %s\n", result)
	case err := <-errorCh:
		fmt.Printf("Error: %v\n", err)
	case <-time.After(time.Second):
		fmt.Println("Operation timed out")
	}
}

// 5.2 优雅关闭模式
func GracefulShutdownExample() {
	done := make(chan struct{})
	signalCh := make(chan os.Signal, 1)

	// 启动主工作循环
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Shutting down gracefully...")
				return
			default:
				// 执行正常工作
				fmt.Println("Working...")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	// 模拟信号接收
	go func() {
		time.Sleep(time.Second * 2)
		close(done)
	}()

	// 等待关闭信号
	select {
	case <-signalCh:
		fmt.Println("Received shutdown signal")
	case <-done:
		fmt.Println("Shutdown initiated")
	}

	fmt.Println("Application stopped")
}

// 6. select 性能优化

// 6.1 避免 select 中的重复计算
func SelectPerformanceExample() {
	ch := make(chan int)

	// 错误示例：在 select 中重复计算
	select {
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	case val := <-ch:
		fmt.Printf("Value: %d\n", val)
	}

	// 正确示例：预先计算
	timeout := time.After(time.Second)
	select {
	case <-timeout:
		fmt.Println("Timeout")
	case val := <-ch:
		fmt.Printf("Value: %d\n", val)
	}
}

// 6.2 select 与 channel 缓冲
func SelectBufferExample() {
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

// 7. select 高级模式

// 7.1 优先级选择
func PrioritySelectExample() {
	highPriority := make(chan string)
	lowPriority := make(chan string)

	// 启动数据源
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Millisecond * 200)
			lowPriority <- fmt.Sprintf("Low priority %d", i)
		}
	}()

	go func() {
		time.Sleep(time.Millisecond * 500)
		highPriority <- "High priority message"
	}()

	// 优先处理高优先级消息
	for i := 0; i < 4; i++ {
		select {
		case msg := <-highPriority:
			fmt.Printf("Processing: %s\n", msg)
		case msg := <-lowPriority:
			fmt.Printf("Processing: %s\n", msg)
		}
	}
}

// 7.2 条件选择
func ConditionalSelectExample() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	condition := true

	// 根据条件选择不同的 channel
	if condition {
		select {
		case val := <-ch1:
			fmt.Printf("Received from ch1: %d\n", val)
		default:
			fmt.Println("No data from ch1")
		}
	} else {
		select {
		case val := <-ch2:
			fmt.Printf("Received from ch2: %d\n", val)
		default:
			fmt.Println("No data from ch2")
		}
	}
}

// 7.3 动态 select
func DynamicSelectExample() {
	channels := []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
	}

	// 启动数据源
	for i, ch := range channels {
		go func(id int, c chan int) {
			time.Sleep(time.Duration(id+1) * time.Millisecond * 100)
			c <- id
		}(i, ch)
	}

	// 动态构建 select cases
	cases := make([]reflect.SelectCase, len(channels))
	for i, ch := range channels {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	// 执行动态 select
	chosen, value, ok := reflect.Select(cases)
	if ok {
		fmt.Printf("Received %d from channel %d\n", value.Int(), chosen)
	}
}

// 8. select 常见陷阱和解决方案

// 8.1 空 select 语句
func EmptySelectExample() {
	// 空 select 会永远阻塞
	// select {}

	// 解决方案：添加 default 或可执行的 case
	select {
	default:
		fmt.Println("Empty select with default")
	}
}

// 8.2 select 中的 nil channel
func NilChannelSelectExample() {
	var ch chan int // nil channel

	// nil channel 在 select 中永远不会被选中
	select {
	case <-ch:
		fmt.Println("This will never execute")
	default:
		fmt.Println("Nil channel case")
	}
}

// 8.3 select 中的重复 case
func DuplicateCaseExample() {
	ch := make(chan int)

	// 编译错误：重复的 case
	// select {
	// case <-ch:
	// case <-ch:
	// }

	// 正确做法：使用不同的 channel 或条件
	ch2 := make(chan int)
	select {
	case <-ch:
		fmt.Println("Received from ch")
	case <-ch2:
		fmt.Println("Received from ch2")
	}
}

// 9. select 测试和调试

// 9.1 select 行为测试
func SelectBehaviorTest() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 同时发送数据
	go func() {
		ch1 <- "ch1"
	}()

	go func() {
		ch2 <- "ch2"
	}()

	// 测试随机选择行为
	counts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		select {
		case msg := <-ch1:
			counts[msg]++
		case msg := <-ch2:
			counts[msg]++
		}
	}

	fmt.Printf("Selection counts: %v\n", counts)
}

// 9.2 select 性能基准测试
func SelectBenchmarkExample() {
	ch := make(chan int, 1000)

	// 填充 channel
	for i := 0; i < 1000; i++ {
		ch <- i
	}

	// 基准测试 select 性能
	start := time.Now()
	for i := 0; i < 1000; i++ {
		select {
		case <-ch:
		default:
		}
	}
	duration := time.Since(start)

	fmt.Printf("Select benchmark: %v\n", duration)
}

// 10. select 最佳实践总结

func SelectBestPractices() {
	// 1. 总是处理 default 情况
	ch := make(chan int)
	select {
	case val := <-ch:
		fmt.Printf("Received: %d\n", val)
	default:
		fmt.Println("No data available")
	}

	// 2. 使用超时控制长时间操作
	select {
	case result := <-ch:
		fmt.Printf("Result: %v\n", result)
	case <-time.After(time.Second):
		fmt.Println("Operation timed out")
	}

	// 3. 结合 Context 进行取消控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case result := <-ch:
		fmt.Printf("Result: %v\n", result)
	case <-ctx.Done():
		fmt.Printf("Cancelled: %v\n", ctx.Err())
	}

	// 4. 避免在 select 中重复计算
	timeout := time.After(time.Second)
	select {
	case <-timeout:
		fmt.Println("Timeout")
	case val := <-ch:
		fmt.Printf("Value: %d\n", val)
	}

	// 5. 正确处理 channel 关闭
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			fmt.Printf("Received: %d\n", val)
		}
	}
}
