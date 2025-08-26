package skill

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// ==================== Go语言信号量 semaphore.Weighted 详解 ====================

// 1. semaphore.Weighted 基本介绍

/*
semaphore.Weighted 是 Go 官方提供的信号量实现
- 基于权重机制，支持不同权重的资源获取
- 线程安全，支持并发访问
- 支持超时和上下文取消
- 性能优异，适合高并发场景
- 需要导入 "golang.org/x/sync/semaphore"
*/

// 2. 基本用法示例

// 2.1 简单信号量使用
func BasicSemaphoreExample() {
	// 创建一个容量为3的信号量
	sem := semaphore.NewWeighted(3)

	var wg sync.WaitGroup

	// 启动5个goroutine，但最多只有3个能同时运行
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 获取信号量
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Printf("Goroutine %d failed to acquire semaphore: %v", id, err)
				return
			}
			defer sem.Release(1) // 释放信号量

			log.Printf("Goroutine %d is running", id)
			time.Sleep(time.Second) // 模拟工作
			log.Printf("Goroutine %d finished", id)
		}(i)
	}

	wg.Wait()
	log.Println("All goroutines completed")
}

// 2.2 带权重的信号量使用
func WeightedSemaphoreExample() {
	// 创建一个总容量为10的信号量
	sem := semaphore.NewWeighted(10)

	var wg sync.WaitGroup

	// 不同类型的任务需要不同数量的资源
	tasks := []struct {
		name   string
		weight int64
	}{
		{"light-task", 1},
		{"medium-task", 3},
		{"heavy-task", 5},
		{"light-task", 1},
		{"medium-task", 3},
	}

	for i, task := range tasks {
		wg.Add(1)
		go func(id int, taskName string, weight int64) {
			defer wg.Done()

			log.Printf("Task %d (%s) trying to acquire %d resources", id, taskName, weight)

			// 获取指定权重的资源
			if err := sem.Acquire(context.Background(), weight); err != nil {
				log.Printf("Task %d failed to acquire semaphore: %v", id, err)
				return
			}
			defer sem.Release(weight) // 释放相同权重的资源

			log.Printf("Task %d (%s) acquired %d resources, running...", id, taskName, weight)
			time.Sleep(time.Duration(weight) * time.Second) // 模拟工作
			log.Printf("Task %d (%s) finished", id, taskName)
		}(i, task.name, task.weight)
	}

	wg.Wait()
	log.Println("All tasks completed")
}

// 3. 高级用法示例

// 3.1 带超时的信号量
func TimeoutSemaphoreExample() {
	sem := semaphore.NewWeighted(2)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 创建带超时的上下文
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			log.Printf("Goroutine %d trying to acquire semaphore with timeout", id)

			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("Goroutine %d failed to acquire semaphore: %v", id, err)
				return
			}
			defer sem.Release(1)

			log.Printf("Goroutine %d acquired semaphore, working...", id)
			time.Sleep(3 * time.Second) // 工作时间超过超时时间
			log.Printf("Goroutine %d finished", id)
		}(i)
	}

	wg.Wait()
}

// 3.2 可取消的信号量
func CancellableSemaphoreExample() {
	sem := semaphore.NewWeighted(3)

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// 启动工作goroutine
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("Goroutine %d cancelled: %v", id, err)
				return
			}
			defer sem.Release(1)

			log.Printf("Goroutine %d working...", id)
			time.Sleep(time.Second)
			log.Printf("Goroutine %d finished", id)
		}(i)
	}

	// 2秒后取消所有操作
	time.Sleep(2 * time.Second)
	log.Println("Cancelling all operations...")
	cancel()

	wg.Wait()
	log.Println("All goroutines finished")
}

// 4. 实际应用场景

// 4.1 连接池管理
type ConnectionPool struct {
	sem    *semaphore.Weighted
	conns  chan *Connection
	mu     sync.Mutex
	closed bool
}

type Connection struct {
	ID   int
	Used bool
}

func NewConnectionPool(maxConns int) *ConnectionPool {
	pool := &ConnectionPool{
		sem:   semaphore.NewWeighted(int64(maxConns)),
		conns: make(chan *Connection, maxConns),
	}

	// 初始化连接
	for i := 0; i < maxConns; i++ {
		pool.conns <- &Connection{ID: i}
	}

	return pool
}

func (cp *ConnectionPool) GetConnection(ctx context.Context) (*Connection, error) {
	// 获取信号量许可
	if err := cp.sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}

	// 获取连接
	select {
	case conn := <-cp.conns:
		conn.Used = true
		return conn, nil
	case <-ctx.Done():
		cp.sem.Release(1)
		return nil, ctx.Err()
	}
}

func (cp *ConnectionPool) ReleaseConnection(conn *Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.closed {
		return
	}

	conn.Used = false
	cp.conns <- conn
	cp.sem.Release(1)
}

func (cp *ConnectionPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.closed = true
	close(cp.conns)
}

// 4.2 限流器实现
type RateLimiter struct {
	sem    *semaphore.Weighted
	rate   int64
	burst  int64
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRateLimiter(rate int64, burst int64) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	limiter := &RateLimiter{
		sem:    semaphore.NewWeighted(burst),
		rate:   rate,
		burst:  burst,
		ctx:    ctx,
		cancel: cancel,
	}

	// 启动令牌桶填充
	limiter.ticker = time.NewTicker(time.Second / time.Duration(rate))
	go limiter.fillTokens()

	return limiter
}

func (rl *RateLimiter) fillTokens() {
	for {
		select {
		case <-rl.ticker.C:
			// 释放一个令牌
			rl.sem.Release(1)
		case <-rl.ctx.Done():
			return
		}
	}
}

func (rl *RateLimiter) Allow(ctx context.Context) error {
	return rl.sem.Acquire(ctx, 1)
}

func (rl *RateLimiter) Close() {
	rl.cancel()
	rl.ticker.Stop()
}

// 4.3 资源管理器
type ResourceManager struct {
	sem       *semaphore.Weighted
	resources map[string]interface{}
	mu        sync.RWMutex
}

func NewResourceManager(maxResources int64) *ResourceManager {
	return &ResourceManager{
		sem:       semaphore.NewWeighted(maxResources),
		resources: make(map[string]interface{}),
	}
}

func (rm *ResourceManager) AcquireResource(ctx context.Context, resourceID string, weight int64) error {
	// 获取资源许可
	if err := rm.sem.Acquire(ctx, weight); err != nil {
		return fmt.Errorf("failed to acquire resource permit: %w", err)
	}

	// 检查资源是否可用
	rm.mu.RLock()
	if _, exists := rm.resources[resourceID]; exists {
		rm.mu.RUnlock()
		rm.sem.Release(weight) // 释放许可
		return fmt.Errorf("resource %s already in use", resourceID)
	}
	rm.mu.RUnlock()

	// 标记资源为使用中
	rm.mu.Lock()
	rm.resources[resourceID] = struct{}{}
	rm.mu.Unlock()

	return nil
}

func (rm *ResourceManager) ReleaseResource(resourceID string, weight int64) {
	rm.mu.Lock()
	delete(rm.resources, resourceID)
	rm.mu.Unlock()

	rm.sem.Release(weight)
}

// 5. 性能测试和对比

// 5.1 信号量性能测试
func SemaphorePerformanceTest() {
	const (
		numGoroutines = 1000
		semLimit      = 100
	)

	sem := semaphore.NewWeighted(semLimit)

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if err := sem.Acquire(context.Background(), 1); err != nil {
				return
			}
			defer sem.Release(1)

			// 模拟工作
			time.Sleep(time.Millisecond)
		}(i)
	}

	wg.Wait()

	duration := time.Since(start)
	log.Printf("Processed %d goroutines in %v", numGoroutines, duration)
}

// 5.2 不同权重场景测试
func WeightedSemaphoreTest() {
	sem := semaphore.NewWeighted(10)

	var wg sync.WaitGroup

	// 混合权重任务
	tasks := []int64{1, 2, 3, 1, 2, 3, 1, 2, 3, 1}

	start := time.Now()

	for i, weight := range tasks {
		wg.Add(1)
		go func(id int, w int64) {
			defer wg.Done()

			if err := sem.Acquire(context.Background(), w); err != nil {
				log.Printf("Task %d failed: %v", id, err)
				return
			}
			defer sem.Release(w)

			log.Printf("Task %d (weight: %d) running", id, w)
			time.Sleep(time.Duration(w) * time.Millisecond)
			log.Printf("Task %d completed", id)
		}(i, weight)
	}

	wg.Wait()

	duration := time.Since(start)
	log.Printf("All weighted tasks completed in %v", duration)
}

// 6. 错误处理和最佳实践

// 6.1 错误处理示例
func SemaphoreErrorHandling() {
	sem := semaphore.NewWeighted(2)

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 尝试获取信号量
	if err := sem.Acquire(ctx, 1); err != nil {
		switch err {
		case context.DeadlineExceeded:
			log.Println("Semaphore acquisition timed out")
		case context.Canceled:
			log.Println("Semaphore acquisition was cancelled")
		default:
			log.Printf("Unexpected error: %v", err)
		}
		return
	}

	defer sem.Release(1)
	log.Println("Successfully acquired semaphore")
}

// 6.2 最佳实践示例
func SemaphoreBestPractices() {
	// 1. 总是使用defer释放资源
	sem := semaphore.NewWeighted(5)

	ctx := context.Background()
	if err := sem.Acquire(ctx, 1); err != nil {
		log.Printf("Failed to acquire: %v", err)
		return
	}
	defer sem.Release(1) // 确保资源被释放

	// 2. 使用适当的超时
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sem.Acquire(timeoutCtx, 1); err != nil {
		log.Printf("Acquisition failed: %v", err)
		return
	}
	defer sem.Release(1)

	// 3. 合理设置权重
	// 根据资源消耗设置权重，而不是固定值
	resourceWeight := int64(rand.Intn(3) + 1) // 1-3的随机权重

	if err := sem.Acquire(ctx, resourceWeight); err != nil {
		log.Printf("Failed to acquire with weight %d: %v", resourceWeight, err)
		return
	}
	defer sem.Release(resourceWeight)
}

// 7. 常见陷阱和注意事项

// 7.1 避免死锁
func AvoidDeadlockExample() {
	sem := semaphore.NewWeighted(1)

	// 错误示例：在同一个goroutine中重复获取
	// sem.Acquire(ctx, 1)
	// sem.Acquire(ctx, 1) // 这会导致死锁

	// 正确示例：使用权重
	if err := sem.Acquire(context.Background(), 2); err != nil {
		log.Printf("Failed to acquire: %v", err)
		return
	}
	defer sem.Release(2)

	log.Println("Successfully acquired with weight 2")
}

// 7.2 避免资源泄漏
func AvoidResourceLeakExample() {
	sem := semaphore.NewWeighted(5)

	// 错误示例：忘记释放资源
	// sem.Acquire(ctx, 1)
	// // 忘记调用 sem.Release(1)

	// 正确示例：使用defer确保释放
	if err := sem.Acquire(context.Background(), 1); err != nil {
		return
	}
	defer sem.Release(1) // 确保资源被释放

	// 执行工作...
	log.Println("Working with acquired resource")
}

// 8. 信号量的其他实现方式

// 8.1 基于channel的简单信号量
type SimpleSemaphore struct {
	sem chan struct{}
}

func NewSimpleSemaphore(limit int) *SimpleSemaphore {
	return &SimpleSemaphore{
		sem: make(chan struct{}, limit),
	}
}

func (ss *SimpleSemaphore) Acquire() {
	ss.sem <- struct{}{}
}

func (ss *SimpleSemaphore) Release() {
	<-ss.sem
}

func (ss *SimpleSemaphore) TryAcquire() bool {
	select {
	case ss.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// 8.2 带超时的channel信号量
func (ss *SimpleSemaphore) AcquireWithTimeout(timeout time.Duration) bool {
	select {
	case ss.sem <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// 9. 运行示例的主函数
func RunSemaphoreExamples() {
	log.Println("=== Basic Semaphore Example ===")
	BasicSemaphoreExample()

	log.Println("\n=== Weighted Semaphore Example ===")
	WeightedSemaphoreExample()

	log.Println("\n=== Timeout Semaphore Example ===")
	TimeoutSemaphoreExample()

	log.Println("\n=== Cancellable Semaphore Example ===")
	CancellableSemaphoreExample()

	log.Println("\n=== Performance Test ===")
	SemaphorePerformanceTest()

	log.Println("\n=== Weighted Test ===")
	WeightedSemaphoreTest()

	log.Println("\n=== Best Practices ===")
	SemaphoreBestPractices()
}

// 10. 安装依赖说明

/*
要使用 semaphore.Weighted，需要先安装依赖：

go get golang.org/x/sync

或者使用 Go modules：

go mod init your-project
go get golang.org/x/sync

然后在代码中导入：
import "golang.org/x/sync/semaphore"
*/
