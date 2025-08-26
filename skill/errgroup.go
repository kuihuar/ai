package skill

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// ==================== errgroup 错误组详解 ====================

// 1. errgroup 基本概念

/*
errgroup 是 Go 语言中用于管理一组 goroutine 的同步原语，它提供了以下功能：
1. 等待所有 goroutine 完成
2. 收集第一个发生的错误
3. 当发生错误时取消所有 goroutine
4. 支持 context 取消

errgroup 是 sync.WaitGroup 的增强版本，专门用于处理可能出错的并发任务。
*/

// 2. errgroup 基本用法

// 2.1 基本 errgroup 使用
func BasicErrGroupExample() {
	var g errgroup.Group

	// 启动多个 goroutine
	g.Go(func() error {
		time.Sleep(time.Millisecond * 100)
		fmt.Println("Task 1 completed")
		return nil
	})

	g.Go(func() error {
		time.Sleep(time.Millisecond * 200)
		fmt.Println("Task 2 completed")
		return nil
	})

	g.Go(func() error {
		time.Sleep(time.Millisecond * 150)
		fmt.Println("Task 3 completed")
		return nil
	})

	// 等待所有 goroutine 完成并检查错误
	if err := g.Wait(); err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}
}

// 2.2 errgroup 错误处理
func ErrGroupErrorHandlingExample() {
	var g errgroup.Group

	g.Go(func() error {
		time.Sleep(time.Millisecond * 100)
		fmt.Println("Task 1 completed")
		return nil
	})

	g.Go(func() error {
		time.Sleep(time.Millisecond * 50)
		fmt.Println("Task 2 failed")
		return errors.New("task 2 failed")
	})

	g.Go(func() error {
		time.Sleep(time.Millisecond * 200)
		fmt.Println("Task 3 completed")
		return nil
	})

	// 当任何一个 goroutine 返回错误时，Wait 会立即返回第一个错误
	if err := g.Wait(); err != nil {
		fmt.Printf("First error: %v\n", err)
	}
}

// 2.3 errgroup 与 context 结合
func ErrGroupWithContextExample() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		select {
		case <-time.After(time.Millisecond * 500):
			fmt.Println("Task 1 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 1 cancelled")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		select {
		case <-time.After(time.Millisecond * 100):
			fmt.Println("Task 2 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 2 cancelled")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		select {
		case <-time.After(time.Millisecond * 200):
			fmt.Println("Task 3 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 3 cancelled")
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 3. errgroup 高级用法

// 3.1 限制并发数量
func ErrGroupWithLimitExample() {
	// 限制最多同时运行 2 个 goroutine
	g := new(errgroup.Group)
	sem := make(chan struct{}, 2)

	tasks := []string{"task1", "task2", "task3", "task4", "task5"}

	for _, task := range tasks {
		task := task // 创建副本避免闭包问题
		g.Go(func() error {
			sem <- struct{}{} // 获取信号量
			defer func() {
				<-sem // 释放信号量
			}()

			time.Sleep(time.Millisecond * 100)
			fmt.Printf("%s completed\n", task)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 3.2 使用 errgroup.Group 的 SetLimit 方法（Go 1.19+）
func ErrGroupSetLimitExample() {
	g := new(errgroup.Group)
	g.SetLimit(2) // 限制最多同时运行 2 个 goroutine

	tasks := []string{"task1", "task2", "task3", "task4", "task5"}

	for _, task := range tasks {
		task := task
		g.Go(func() error {
			time.Sleep(time.Millisecond * 100)
			fmt.Printf("%s completed\n", task)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 3.3 错误传播和取消
func ErrGroupErrorPropagationExample() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		select {
		case <-time.After(time.Millisecond * 100):
			fmt.Println("Task 1 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 1 cancelled due to error in other task")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		time.Sleep(time.Millisecond * 50)
		fmt.Println("Task 2 failed")
		return errors.New("task 2 failed")
	})

	g.Go(func() error {
		select {
		case <-time.After(time.Millisecond * 200):
			fmt.Println("Task 3 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 3 cancelled due to error in other task")
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 4. errgroup 实际应用场景

// 4.1 并行文件处理
type FileProcessorWithErrGroup struct {
	concurrency int
}

func NewFileProcessorWithErrGroup(concurrency int) *FileProcessorWithErrGroup {
	return &FileProcessorWithErrGroup{concurrency: concurrency}
}

func (fp *FileProcessorWithErrGroup) ProcessFiles(files []string) error {
	g := new(errgroup.Group)
	g.SetLimit(fp.concurrency)

	for _, file := range files {
		file := file
		g.Go(func() error {
			return fp.processFile(file)
		})
	}

	return g.Wait()
}

func (fp *FileProcessorWithErrGroup) processFile(filename string) error {
	// 模拟文件处理
	time.Sleep(time.Millisecond * 100)

	// 模拟某些文件处理失败
	if filename == "error.txt" {
		return fmt.Errorf("failed to process file: %s", filename)
	}

	fmt.Printf("Processed file: %s\n", filename)
	return nil
}

// 4.2 并行 HTTP 请求
type HTTPClient struct {
	timeout time.Duration
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{timeout: timeout}
}

func (hc *HTTPClient) FetchURLs(urls []string) ([]string, error) {
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(5) // 限制并发请求数

	results := make([]string, len(urls))

	for i, url := range urls {
		i, url := i, url // 创建副本
		g.Go(func() error {
			result, err := hc.fetchURL(ctx, url)
			if err != nil {
				return err
			}
			results[i] = result
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func (hc *HTTPClient) fetchURL(ctx context.Context, url string) (string, error) {
	// 模拟 HTTP 请求
	select {
	case <-time.After(hc.timeout):
		return fmt.Sprintf("Response from %s", url), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// 4.3 数据库批量操作
type DatabaseProcessor struct {
	batchSize int
}

func NewDatabaseProcessor(batchSize int) *DatabaseProcessor {
	return &DatabaseProcessor{batchSize: batchSize}
}

func (dp *DatabaseProcessor) ProcessBatch(records []Record) error {
	g := new(errgroup.Group)
	g.SetLimit(3) // 限制并发数据库连接数

	// 分批处理
	for i := 0; i < len(records); i += dp.batchSize {
		end := i + dp.batchSize
		if end > len(records) {
			end = len(records)
		}

		batch := records[i:end]
		g.Go(func() error {
			return dp.processBatch(batch)
		})
	}

	return g.Wait()
}

type Record struct {
	ID   int
	Data string
}

func (dp *DatabaseProcessor) processBatch(records []Record) error {
	// 模拟数据库批量操作
	time.Sleep(time.Millisecond * 50)

	// 模拟某些批次处理失败
	if len(records) > 0 && records[0].ID == 999 {
		return fmt.Errorf("failed to process batch with ID: %d", records[0].ID)
	}

	fmt.Printf("Processed batch with %d records\n", len(records))
	return nil
}

// 5. errgroup 与其他同步原语结合

// 5.1 errgroup 与 sync.Mutex
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *SafeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count++
}

func (sc *SafeCounter) GetCount() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

func ErrGroupWithMutexExample() {
	counter := &SafeCounter{}
	var g errgroup.Group

	// 启动多个 goroutine 并发增加计数器
	for i := 0; i < 100; i++ {
		g.Go(func() error {
			counter.Increment()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Final count: %d\n", counter.GetCount())
}

// 5.2 errgroup 与 sync.Map
func ErrGroupWithSyncMapExample() {
	var m sync.Map
	var g errgroup.Group

	// 并发写入 sync.Map
	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			m.Store(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// 读取所有值
	for i := 0; i < 10; i++ {
		if value, ok := m.Load(fmt.Sprintf("key%d", i)); ok {
			fmt.Printf("key%d: %s\n", i, value)
		}
	}
}

// 6. errgroup 错误处理策略

// 6.1 错误分类处理
type TaskError struct {
	TaskID string
	Err    error
}

func (te TaskError) Error() string {
	return fmt.Sprintf("task %s failed: %v", te.TaskID, te.Err)
}

func ErrGroupErrorClassificationExample() {
	g := new(errgroup.Group)
	errors := make(chan TaskError, 10)

	tasks := []string{"task1", "task2", "task3", "task4"}

	for _, taskID := range tasks {
		taskID := taskID
		g.Go(func() error {
			if err := processTask(taskID); err != nil {
				taskErr := TaskError{TaskID: taskID, Err: err}
				select {
				case errors <- taskErr:
				default:
					// 错误通道已满，记录日志
					log.Printf("Error channel full, dropping error: %v", taskErr)
				}
				return taskErr
			}
			return nil
		})
	}

	// 等待所有任务完成
	if err := g.Wait(); err != nil {
		fmt.Printf("Group error: %v\n", err)
	}

	// 处理收集到的错误
	close(errors)
	for taskErr := range errors {
		fmt.Printf("Task error: %v\n", taskErr)
	}
}

func processTask(taskID string) error {
	time.Sleep(time.Millisecond * 100)

	// 模拟某些任务失败
	if taskID == "task2" {
		return errors.New("task2 failed")
	}

	fmt.Printf("Task %s completed\n", taskID)
	return nil
}

// 6.2 错误重试机制
func ErrGroupWithRetryExample() {
	g := new(errgroup.Group)
	maxRetries := 3

	tasks := []string{"task1", "task2", "task3"}

	for _, taskID := range tasks {
		taskID := taskID
		g.Go(func() error {
			return retryTask(taskID, maxRetries)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error after retries: %v\n", err)
	}
}

func retryTask(taskID string, maxRetries int) error {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := processTask(taskID); err != nil {
			lastErr = err
			fmt.Printf("Task %s failed (attempt %d/%d): %v\n", taskID, attempt, maxRetries, err)

			if attempt < maxRetries {
				time.Sleep(time.Millisecond * time.Duration(attempt*100))
				continue
			}
		} else {
			return nil
		}
	}

	return fmt.Errorf("task %s failed after %d attempts: %v", taskID, maxRetries, lastErr)
}

// 7. errgroup 性能优化

// 7.1 批量任务处理
func ErrGroupBatchProcessingExample() {
	g := new(errgroup.Group)
	g.SetLimit(5) // 限制并发数

	tasks := make([]string, 100)
	for i := range tasks {
		tasks[i] = fmt.Sprintf("task%d", i)
	}

	// 分批处理任务
	batchSize := 20
	for i := 0; i < len(tasks); i += batchSize {
		end := i + batchSize
		if end > len(tasks) {
			end = len(tasks)
		}

		batch := tasks[i:end]
		g.Go(func() error {
			return processBatchWithErrGroup(batch)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func processBatchWithErrGroup(tasks []string) error {
	for _, task := range tasks {
		if err := processTask(task); err != nil {
			return err
		}
	}
	return nil
}

// 7.2 内存优化
func ErrGroupMemoryOptimizationExample() {
	g := new(errgroup.Group)
	g.SetLimit(10)

	// 使用对象池减少内存分配
	var pool sync.Pool
	pool.New = func() interface{} {
		return make([]byte, 1024)
	}

	for i := 0; i < 100; i++ {
		g.Go(func() error {
			// 从池中获取对象
			buf := pool.Get().([]byte)
			defer pool.Put(buf) // 归还到池中

			// 使用缓冲区处理任务
			time.Sleep(time.Millisecond * 10)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 8. errgroup 测试和调试

// 8.1 性能基准测试
func ErrGroupBenchmarkExample() {
	g := new(errgroup.Group)
	g.SetLimit(10)

	start := time.Now()

	for i := 0; i < 1000; i++ {
		g.Go(func() error {
			time.Sleep(time.Microsecond * 100)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	duration := time.Since(start)
	fmt.Printf("Processed 1000 tasks in %v\n", duration)
}

// 8.2 错误统计
func ErrGroupErrorStatisticsExample() {
	g := new(errgroup.Group)
	g.SetLimit(5)

	var errorCount int
	var mu sync.Mutex

	tasks := make([]string, 50)
	for i := range tasks {
		tasks[i] = fmt.Sprintf("task%d", i)
	}

	for _, taskID := range tasks {
		taskID := taskID
		g.Go(func() error {
			if err := processTask(taskID); err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Group error: %v\n", err)
	}

	fmt.Printf("Total errors: %d\n", errorCount)
}

// 9. errgroup 最佳实践

// 9.1 资源清理
func ErrGroupResourceCleanupExample() {
	g, ctx := errgroup.WithContext(context.Background())

	// 模拟资源
	resources := make([]string, 5)
	for i := range resources {
		resources[i] = fmt.Sprintf("resource%d", i)
	}

	// 启动任务
	for _, resource := range resources {
		//	i, resource := i, resource
		g.Go(func() error {
			defer func() {
				// 确保资源被清理
				fmt.Printf("Cleaning up resource: %s\n", resource)
			}()

			select {
			case <-time.After(time.Millisecond * 100):
				fmt.Printf("Processing resource: %s\n", resource)
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 9.2 优雅关闭
func ErrGroupGracefulShutdownExample() {
	g, ctx := errgroup.WithContext(context.Background())
	shutdown := make(chan struct{})

	// 启动工作 goroutine
	g.Go(func() error {
		for {
			select {
			case <-shutdown:
				fmt.Println("Worker shutting down")
				return nil
			case <-ctx.Done():
				fmt.Println("Worker cancelled")
				return ctx.Err()
			default:
				// 执行工作
				time.Sleep(time.Millisecond * 100)
			}
		}
	})

	// 模拟关闭信号
	go func() {
		time.Sleep(time.Second)
		close(shutdown)
	}()

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// 9.3 错误传播策略
func ErrGroupErrorPropagationStrategyExample() {
	g := new(errgroup.Group)
	g.SetLimit(3)

	// 定义错误处理策略
	errorHandler := func(taskID string, err error) {
		fmt.Printf("Handling error for task %s: %v\n", taskID, err)
		// 可以记录日志、发送告警等
	}

	tasks := []string{"task1", "task2", "task3", "task4"}

	for _, taskID := range tasks {
		taskID := taskID
		g.Go(func() error {
			if err := processTask(taskID); err != nil {
				errorHandler(taskID, err)
				return err // 继续传播错误
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Group completed with error: %v\n", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}
}

// 10. 总结

func ErrGroupSummary() {
	fmt.Println("=== errgroup 总结 ===")
	fmt.Println("1. errgroup 是 sync.WaitGroup 的增强版本")
	fmt.Println("2. 支持错误收集和传播")
	fmt.Println("3. 支持 context 取消")
	fmt.Println("4. 支持并发数量限制")
	fmt.Println("5. 适用于需要错误处理的并发任务")
	fmt.Println("6. 提供了优雅的资源管理和错误处理机制")
}
