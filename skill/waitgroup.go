package skill

import (
	"fmt"
	"sync"
	"time"
)

// ==================== sync.WaitGroup 同步原语 ====================

// 1. 基本用法示例

// 1.1 简单的任务协调
type TaskCoordinator struct {
	wg sync.WaitGroup
}

func NewTaskCoordinator() *TaskCoordinator {
	return &TaskCoordinator{}
}

func (tc *TaskCoordinator) RunTasks(taskCount int) {
	for i := 0; i < taskCount; i++ {
		tc.wg.Add(1)
		go func(taskID int) {
			defer tc.wg.Done()
			fmt.Printf("Task %d starting\n", taskID)
			time.Sleep(time.Millisecond * 100)
			fmt.Printf("Task %d completed\n", taskID)
		}(i)
	}

	tc.wg.Wait()
	fmt.Println("All tasks completed")
}

// 1.2 错误处理版本
type ErrorHandlingCoordinator struct {
	wg      sync.WaitGroup
	errChan chan error
	mu      sync.Mutex
	errors  []error
}

func NewErrorHandlingCoordinator() *ErrorHandlingCoordinator {
	return &ErrorHandlingCoordinator{
		errChan: make(chan error, 10),
		errors:  make([]error, 0),
	}
}

func (ehc *ErrorHandlingCoordinator) RunTasksWithErrorHandling(taskCount int) error {
	// 启动错误收集协程
	ehc.wg.Add(1)
	go ehc.collectErrors()

	// 启动任务
	for i := 0; i < taskCount; i++ {
		ehc.wg.Add(1)
		go func(taskID int) {
			defer ehc.wg.Done()
			if err := ehc.runTask(taskID); err != nil {
				ehc.errChan <- err
			}
		}(i)
	}

	// 等待所有任务完成
	ehc.wg.Wait()
	close(ehc.errChan)

	// 检查是否有错误
	if len(ehc.errors) > 0 {
		return fmt.Errorf("completed with %d errors", len(ehc.errors))
	}

	return nil
}

func (ehc *ErrorHandlingCoordinator) collectErrors() {
	defer ehc.wg.Done()

	for err := range ehc.errChan {
		ehc.mu.Lock()
		ehc.errors = append(ehc.errors, err)
		ehc.mu.Unlock()
	}
}

func (ehc *ErrorHandlingCoordinator) runTask(taskID int) error {
	fmt.Printf("Task %d starting\n", taskID)
	time.Sleep(time.Millisecond * 100)

	// 模拟某些任务失败
	if taskID%3 == 0 {
		return fmt.Errorf("task %d failed", taskID)
	}

	fmt.Printf("Task %d completed\n", taskID)
	return nil
}

// 2. 高级用法

// 2.1 带超时的等待组

type TimeoutWaitGorup struct {
	wg      sync.WaitGroup
	timeout time.Duration
}

func NewTimeoutWaitGorup(timeout time.Duration) *TimeoutWaitGorup {
	return &TimeoutWaitGorup{
		timeout: timeout,
	}
}
func (twg *TimeoutWaitGorup) WaitWithTimeout() error {
	done := make(chan struct{})

	go func() {
		twg.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(twg.timeout):
		return fmt.Errorf("wait group timeout after %v", twg.timeout)
	}
}

// 2.2 进度跟踪等待组
type ProgressWaitGroup struct {
	wg        sync.WaitGroup
	mu        sync.Mutex
	total     int
	completed int
}

func NewProgressWaitGroup(total int) *ProgressWaitGroup {
	return &ProgressWaitGroup{
		total: total,
	}
}

func (pwg *ProgressWaitGroup) AddTask() {
	pwg.wg.Add(1)
}

func (pwg *ProgressWaitGroup) CompleteTask() {
	pwg.mu.Lock()
	pwg.completed++
	progress := float64(pwg.completed) / float64(pwg.total) * 100
	pwg.mu.Unlock()

	fmt.Printf("Progress: %.1f%% (%d/%d)\n", progress, pwg.completed, pwg.total)
	pwg.wg.Done()
}

func (pwg *ProgressWaitGroup) Wait() {
	pwg.wg.Wait()
	fmt.Println("All tasks completed!")
}

// 2.3 分层等待组
type HierarchicalWaitGroup struct {
	parentWg sync.WaitGroup
	childWgs map[string]*sync.WaitGroup
	mu       sync.RWMutex
}

func NewHierarchicalWaitGroup() *HierarchicalWaitGroup {
	return &HierarchicalWaitGroup{
		childWgs: make(map[string]*sync.WaitGroup),
	}
}

func (hwg *HierarchicalWaitGroup) AddChildGroup(name string) {
	hwg.mu.Lock()
	defer hwg.mu.Unlock()

	hwg.childWgs[name] = &sync.WaitGroup{}
}

func (hwg *HierarchicalWaitGroup) AddTask(groupName string) {
	hwg.mu.RLock()
	childWg, exists := hwg.childWgs[groupName]
	hwg.mu.RUnlock()

	if !exists {
		hwg.AddChildGroup(groupName)
		hwg.mu.RLock()
		childWg = hwg.childWgs[groupName]
		hwg.mu.RUnlock()
	}

	hwg.parentWg.Add(1)
	childWg.Add(1)
}

func (hwg *HierarchicalWaitGroup) CompleteTask(groupName string) {
	hwg.mu.RLock()
	childWg := hwg.childWgs[groupName]
	hwg.mu.RUnlock()

	childWg.Done()
	hwg.parentWg.Done()
}

func (hwg *HierarchicalWaitGroup) WaitForGroup(groupName string) {
	hwg.mu.RLock()
	childWg := hwg.childWgs[groupName]
	hwg.mu.RUnlock()

	childWg.Wait()
}

func (hwg *HierarchicalWaitGroup) WaitForAll() {
	hwg.parentWg.Wait()
}

// 3. 实际应用场景

// 3.1 文件处理器
type FileProcessor struct {
	wg sync.WaitGroup
}

func (fp *FileProcessor) ProcessFiles(files []string) {
	for _, file := range files {
		fp.wg.Add(1)
		go func(filename string) {
			defer fp.wg.Done()
			fp.processFile(filename)
		}(file)
	}

	fp.wg.Wait()
	fmt.Println("All files processed")
}

func (fp *FileProcessor) processFile(filename string) {
	fmt.Printf("Processing file: %s\n", filename)
	time.Sleep(time.Millisecond * 200)
	fmt.Printf("Completed file: %s\n", filename)
}

// 3.2 网络请求处理器
type NetworkProcessor struct {
	wg sync.WaitGroup
}

func (np *NetworkProcessor) ProcessRequests(urls []string) {
	results := make(chan string, len(urls))

	for _, url := range urls {
		np.wg.Add(1)
		go func(u string) {
			defer np.wg.Done()
			result := np.fetchURL(u)
			results <- result
		}(url)
	}

	// 等待所有请求完成
	go func() {
		np.wg.Wait()
		close(results)
	}()

	// 收集结果
	for result := range results {
		fmt.Println(result)
	}
}

func (np *NetworkProcessor) fetchURL(url string) string {
	fmt.Printf("Fetching: %s\n", url)
	time.Sleep(time.Millisecond * 100)
	return fmt.Sprintf("Result from %s", url)
}

// 4. 性能优化技巧

// 4.1 批量任务处理
type BatchProcessorWithWaitGroup struct {
	wg sync.WaitGroup
}

func (bp *BatchProcessorWithWaitGroup) ProcessBatch(items []interface{}, batchSize int) {
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		bp.wg.Add(1)
		go func(batchItems []interface{}) {
			defer bp.wg.Done()
			bp.processBatch(batchItems)
		}(batch)
	}

	bp.wg.Wait()
}

func (bp *BatchProcessorWithWaitGroup) processBatch(items []interface{}) {
	fmt.Printf("Processing batch of %d items\n", len(items))
	time.Sleep(time.Millisecond * 50)
}

// 4.2 资源池协调
type ResourcePoolCoordinator struct {
	wg sync.WaitGroup
}

func (rpc *ResourcePoolCoordinator) CoordinateWorkers(workerCount int, taskCount int) {
	tasks := make(chan int, taskCount)

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		rpc.wg.Add(1)
		go func(workerID int) {
			defer rpc.wg.Done()
			rpc.worker(workerID, tasks)
		}(i)
	}

	// 分发任务
	for i := 0; i < taskCount; i++ {
		tasks <- i
	}
	close(tasks)

	rpc.wg.Wait()
	fmt.Println("All workers completed")
}

func (rpc *ResourcePoolCoordinator) worker(workerID int, tasks <-chan int) {
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", workerID, task)
		time.Sleep(time.Millisecond * 50)
	}
}

// 5. 注意事项和最佳实践

// 5.1 避免死锁
type SafeWaitGroup struct {
	wg sync.WaitGroup
}

func (swg *SafeWaitGroup) SafeAdd(delta int) {
	if delta < 0 {
		panic("negative delta")
	}
	swg.wg.Add(delta)
}

func (swg *SafeWaitGroup) SafeDone() {
	swg.wg.Done()
}

// 5.2 正确处理 panic
func (swg *SafeWaitGroup) SafeWait() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	swg.wg.Wait()
}

// 5.3 避免重复调用
type OnceWaitGroup struct {
	wg   sync.WaitGroup
	once sync.Once
}

func (owg *OnceWaitGroup) WaitOnce() {
	owg.once.Do(func() {
		owg.wg.Wait()
	})
}
