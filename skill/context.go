package skill

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ==================== context.Context 上下文控制详解 ====================

// 1. Context 基本概念

/*
Context 是 Go 语言中用于跨 API 边界和进程间传递截止时间、取消信号和其他请求范围值的标准方式。
它不是同步原语，而是一种协调机制，用于在多个 goroutine 之间传递控制信息。

核心价值：
1. 取消控制 - 支持优雅取消和超时控制
2. 值传递 - 在调用链中传递请求范围的值
3. 截止时间 - 设置操作的截止时间
4. 请求追踪 - 支持分布式追踪和日志记录
*/

// 2. 基础 Context 类型

// 2.1 空 Context
func EmptyContextExample() {
	// context.Background() - 根 Context，永不取消
	ctx := context.Background()

	// context.TODO() - 当不确定使用哪个 Context 时使用
	todoCtx := context.TODO()

	// 检查是否已取消
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled")
	default:
		fmt.Println("Context is active")
	}

	// 获取截止时间
	if deadline, ok := todoCtx.Deadline(); ok {
		fmt.Printf("Deadline: %v\n", deadline)
	} else {
		fmt.Println("No deadline set")
	}
}

// 2.2 带取消的 Context
func CancellableContextExample() {
	// 创建可取消的 Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保资源被释放

	// 启动工作 goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Work cancelled")
				return
			default:
				fmt.Println("Working...")
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	// 模拟工作一段时间后取消
	time.Sleep(time.Second)
	fmt.Println("Cancelling work...")
	cancel()

	time.Sleep(time.Millisecond * 100)
}

// 2.3 带超时的 Context
func TimeoutContextExample() {
	// 创建带超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 启动工作 goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Work timeout: %v\n", ctx.Err())
				return
			default:
				fmt.Println("Working with timeout...")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	// 等待超时
	<-ctx.Done()
	fmt.Println("Timeout reached")
}

// 2.4 带截止时间的 Context
func DeadlineContextExample() {
	// 创建带截止时间的 Context
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 启动工作 goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Work deadline reached: %v\n", ctx.Err())
				return
			default:
				fmt.Println("Working with deadline...")
				time.Sleep(time.Millisecond * 800)
			}
		}
	}()

	// 等待截止时间
	<-ctx.Done()
	fmt.Println("Deadline reached")
}

// 2.5 带值的 Context
func ValueContextExample() {
	// 创建带值的 Context
	ctx := context.WithValue(context.Background(), "user_id", "12345")
	ctx = context.WithValue(ctx, "request_id", "req_67890")

	// 获取值
	if userID, ok := ctx.Value("user_id").(string); ok {
		fmt.Printf("User ID: %s\n", userID)
	}

	if requestID, ok := ctx.Value("request_id").(string); ok {
		fmt.Printf("Request ID: %s\n", requestID)
	}

	// 不存在的键返回 nil
	if value := ctx.Value("non_existent"); value == nil {
		fmt.Println("Key does not exist")
	}
}

// 3. Context 链式传递

// 3.1 多层 Context 传递
func ContextChainExample() {
	// 创建根 Context
	rootCtx := context.Background()

	// 第一层：添加超时
	timeoutCtx, cancel1 := context.WithTimeout(rootCtx, 5*time.Second)
	defer cancel1()

	// 第二层：添加值
	valueCtx := context.WithValue(timeoutCtx, "layer", "second")

	// 第三层：添加取消
	cancelCtx, cancel2 := context.WithCancel(valueCtx)
	defer cancel2()

	// 启动工作 goroutine
	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				fmt.Printf("Context cancelled: %v\n", cancelCtx.Err())
				return
			default:
				if layer, ok := cancelCtx.Value("layer").(string); ok {
					fmt.Printf("Working in %s layer...\n", layer)
				}
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	// 模拟工作一段时间后取消
	time.Sleep(time.Second * 2)
	fmt.Println("Cancelling context chain...")
	cancel2()

	time.Sleep(time.Millisecond * 100)
}

// 4. Context 在实际应用中的使用

// 4.1 HTTP 请求处理
type HTTPHandler struct {
	timeout time.Duration
}

func NewHTTPHandler(timeout time.Duration) *HTTPHandler {
	return &HTTPHandler{timeout: timeout}
}

func (h *HTTPHandler) HandleRequest(userID string) error {
	// 创建带超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	// 添加用户信息
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "request_time", time.Now())

	// 执行请求处理
	return h.processRequest(ctx)
}

func (h *HTTPHandler) processRequest(ctx context.Context) error {
	// 模拟数据库查询
	if err := h.queryDatabase(ctx); err != nil {
		return err
	}

	// 模拟外部 API 调用
	if err := h.callExternalAPI(ctx); err != nil {
		return err
	}

	return nil
}

func (h *HTTPHandler) queryDatabase(ctx context.Context) error {
	// 检查 Context 是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// 模拟数据库查询
		time.Sleep(time.Millisecond * 100)
		fmt.Printf("Database query for user: %s\n", ctx.Value("user_id"))
		return nil
	}
}

func (h *HTTPHandler) callExternalAPI(ctx context.Context) error {
	// 创建带超时的子 Context
	apiCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// 模拟 API 调用
	select {
	case <-apiCtx.Done():
		return apiCtx.Err()
	case <-time.After(time.Millisecond * 200):
		fmt.Println("External API call completed")
		return nil
	}
}

// 4.2 并发任务控制
type TaskManager struct {
	wg sync.WaitGroup
}

func NewTaskManager() *TaskManager {
	return &TaskManager{}
}

func (tm *TaskManager) RunTasks(ctx context.Context, taskCount int) error {
	// 创建错误通道
	errChan := make(chan error, taskCount)

	// 启动任务
	for i := 0; i < taskCount; i++ {
		tm.wg.Add(1)
		go func(taskID int) {
			defer tm.wg.Done()
			if err := tm.runTask(ctx, taskID); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(i)
	}

	// 等待所有任务完成或 Context 取消
	done := make(chan struct{})
	go func() {
		tm.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		// 检查是否有错误
		close(errChan)
		for err := range errChan {
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (tm *TaskManager) runTask(ctx context.Context, taskID int) error {
	// 模拟任务执行
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Printf("Task %d, step %d\n", taskID, i+1)
			time.Sleep(time.Millisecond * 100)
		}
	}
	return nil
}

// 4.3 资源池管理
type ResourcePoolWithContext struct {
	resources chan *ResourceWithContext
	ctx       context.Context
	cancel    context.CancelFunc
}

type ResourceWithContext struct {
	ID   int
	Name string
}

func NewResourcePoolWithContext(poolSize int) *ResourcePoolWithContext {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &ResourcePoolWithContext{
		resources: make(chan *ResourceWithContext, poolSize),
		ctx:       ctx,
		cancel:    cancel,
	}

	// 初始化资源
	for i := 0; i < poolSize; i++ {
		pool.resources <- &ResourceWithContext{
			ID:   i + 1,
			Name: fmt.Sprintf("Resource-%d", i+1),
		}
	}

	return pool
}

func (rp *ResourcePoolWithContext) GetResource(ctx context.Context) (*ResourceWithContext, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rp.ctx.Done():
		return nil, rp.ctx.Err()
	case resource := <-rp.resources:
		return &ResourceWithContext{
			ID:   resource.ID,
			Name: resource.Name,
		}, nil
	}
}

func (rp *ResourcePoolWithContext) ReturnResource(resource *ResourceWithContext) {
	select {
	case <-rp.ctx.Done():
		// 池已关闭，丢弃资源
		return
	case rp.resources <- &ResourceWithContext{
		ID:   resource.ID,
		Name: resource.Name,
	}:
		// 成功返回资源
	}
}

func (rp *ResourcePoolWithContext) Close() {
	rp.cancel()
	close(rp.resources)
}

// 5. Context 最佳实践

// 5.1 Context 传递原则
func ContextBestPractices() {
	// 1. 总是将 Context 作为第一个参数传递
	ctx := context.Background()

	// 2. 不要将 Context 存储在结构体中
	// 错误示例：
	// type BadStruct struct {
	//     ctx context.Context
	// }

	// 3. 使用 WithCancel、WithTimeout、WithDeadline 创建子 Context
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// 4. 总是调用 cancel 函数释放资源
	// defer cancel() // 已在上面调用

	// 5. 在长时间运行的操作中定期检查 Context
	go func() {
		for {
			select {
			case <-timeoutCtx.Done():
				fmt.Println("Context cancelled")
				return
			default:
				// 执行工作
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}

// 5.2 Context 值传递最佳实践
func ContextValueBestPractices() {
	// 1. 使用类型安全的键
	type contextKey string

	const (
		UserIDKey    contextKey = "user_id"
		RequestIDKey contextKey = "request_id"
	)

	// 2. 创建带类型安全键的 Context
	ctx := context.WithValue(context.Background(), UserIDKey, "12345")
	ctx = context.WithValue(ctx, RequestIDKey, "req_67890")

	// 3. 安全地获取值
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		fmt.Printf("User ID: %s\n", userID)
	}

	// 4. 避免传递可变的共享数据
	// 错误示例：
	// ctx = context.WithValue(ctx, "shared_map", make(map[string]interface{}))
}

// 5.3 Context 取消最佳实践
func ContextCancellationBestPractices() {
	// 1. 使用 WithCancel 进行优雅取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. 在适当的时候取消 Context
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Cancelling context...")
		cancel()
	}()

	// 3. 检查取消原因
	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.Canceled:
			fmt.Println("Context was cancelled")
		case context.DeadlineExceeded:
			fmt.Println("Context deadline exceeded")
		}
	}
}

// 6. Context 与同步原语的结合

// 6.1 Context 与 Channel 结合
func ContextWithChannel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataChan := make(chan int, 10)

	// 生产者
	go func() {
		defer close(dataChan)
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return
			case dataChan <- i:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	// 消费者
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer cancelled")
			return
		case data, ok := <-dataChan:
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			fmt.Printf("Received: %d\n", data)
		}
	}
}

// 6.2 Context 与 WaitGroup 结合
func ContextWithWaitGroup() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// 启动工作 goroutine
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d cancelled\n", id)
					return
				default:
					fmt.Printf("Worker %d working...\n", id)
					time.Sleep(time.Millisecond * 200)
				}
			}
		}(i)
	}

	// 等待一段时间后取消
	time.Sleep(time.Second)
	cancel()

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("All workers completed")
}

// 7. Context 性能考虑

// 7.1 Context 链的性能影响
func ContextPerformanceExample() {
	// 创建多层 Context 链
	ctx := context.Background()

	// 添加多层值
	for i := 0; i < 100; i++ {
		ctx = context.WithValue(ctx, fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
	}

	// 获取值（需要遍历整个链）
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_ = ctx.Value("key_50")
	}
	duration := time.Since(start)

	fmt.Printf("Context value lookup took: %v\n", duration)
}

// 7.2 避免 Context 值滥用
func ContextValueAbuseExample() {
	// 错误示例：在 Context 中存储大量数据
	ctx := context.Background()

	// 不要这样做
	largeData := make([]int, 10000)
	for i := range largeData {
		largeData[i] = i
	}
	ctx = context.WithValue(ctx, "large_data", largeData)

	// 正确做法：只存储必要的标识符
	ctx = context.WithValue(ctx, "data_id", "large_data_123")
}
