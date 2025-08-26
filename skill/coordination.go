package skill

import (
	"fmt"
	"sync"
	"time"
)

// 1. 条件变量 (Cond) - 线程同步
// 条件变量是一种同步原语，用于在多个goroutine之间进行通信和协调。它不是锁，而是基于锁实现的等待和通知机制。
// 在 Go 语言中，条件变量（sync.Cond）是一种用于协调多个 Goroutine 同步的机制，它通常与互斥锁（sync.Mutex 或 sync.RWMutex）配合使用，
// 用于等待某个 "条件" 的发生，
// 或通知其他 Goroutine 条件已满足。其核心价值在于：让 Goroutine 在特定条件不满足时阻塞等待，直到条件满足后被唤醒，避免无效的轮询消耗资源。
// 典型应用场景
// 条件变量的核心场景是 "多个 Goroutine 依赖某个共享状态的变化"，以下是最常见的应用场景：

// 1. 生产者消费者 - 协调生产和消费速度
// 2.资源池 - 管理有限资源的分配和回收
// 3.任务队列 - 实现线程池和任务调度
// 4.读写锁 - 实现复杂的读写同步
// 5.信号量 - 控制并发数量
// 6.屏障 - 同步多个goroutine的执行

// 条件变量的核心价值在于：
// 解耦等待和通知：让 Goroutine 在特定条件不满足时阻塞等待，直到条件满足后被唤醒，避免无效的轮询消耗资源。

// 1.1 生产者消费者

type BufferWithCoordination struct {
	mutex    sync.Mutex
	cond     sync.Cond
	data     []int
	capacity int
	size     int
}

func NewBufferWithCoordination(capacity int) *BufferWithCoordination {
	buffer := &BufferWithCoordination{
		data:     make([]int, 0, capacity),
		capacity: capacity,
	}
	buffer.cond = sync.Cond{L: &buffer.mutex}
	return buffer
}

func (b *BufferWithCoordination) Produce(item int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for b.size >= b.capacity {
		fmt.Println("Buffer is full, waiting for consumer to consume")
		b.cond.Wait() // 释放锁，等待条件满足，等待
	}
	b.data = append(b.data, item)
	b.size++
	fmt.Printf("Produced item: %d, Buffer size: %d\n", item, b.size)

	//b.cond.Broadcast()
	// 所有等待的 Goroutine	全局状态变更，需所有等待者响应（如初始化完成）	配置加载完成后唤醒所有依赖 Goroutine
	b.cond.Signal() //生产者唤醒一个消费者
}

func (b *BufferWithCoordination) Consume() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for b.size == 0 {
		fmt.Println("Buffer is empty, waiting for producer to produce")
		b.cond.Wait()
	}
	item := b.data[0]
	b.data = b.data[1:]
	b.size--
	fmt.Printf("Consumed item: %d, Buffer size: %d\n", item, b.size)
	b.cond.Signal() // 通知生产者

	return item
}

// 1.2 资源池
type ResourcePoolWithCoordination struct {
	mu        sync.Mutex
	cond      *sync.Cond
	resources []*ResourceWithCoordination
	maxSize   int
	created   int
}

type ResourceWithCoordination struct {
	ID   int
	Name string
}

func NewResourcePoolWithCoordination(maxSize int) *ResourcePoolWithCoordination {
	pool := ResourcePoolWithCoordination{
		resources: make([]*ResourceWithCoordination, 0, maxSize),
		maxSize:   maxSize,
	}
	pool.cond = sync.NewCond(&pool.mu)
	return &pool
}

func (rp *ResourcePoolWithCoordination) GetResource() *ResourceWithCoordination {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	// 等待资源可用
	for len(rp.resources) == 0 && rp.created >= rp.maxSize {
		fmt.Println("No resources available, waiting...")
		rp.cond.Wait()
	}

	if len(rp.resources) > 0 {
		// 从池中获取资源
		resource := rp.resources[len(rp.resources)-1]
		rp.resources = rp.resources[:len(rp.resources)-1]
		fmt.Printf("Got resource from pool: %d\n", resource.ID)
		return resource
	}

	// 创建新资源
	rp.created++
	resource := &ResourceWithCoordination{ID: rp.created, Name: fmt.Sprintf("Resource-%d", rp.created)}
	fmt.Printf("Created new resource: %d\n", resource.ID)
	return resource
}

func (rp *ResourcePoolWithCoordination) ReturnResource(resource *ResourceWithCoordination) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	if len(rp.resources) < rp.maxSize {
		rp.resources = append(rp.resources, resource)
		fmt.Printf("Returned resource to pool: %d\n", resource.ID)
		rp.cond.Signal() // 通知等待的goroutine
	} else {
		fmt.Printf("Pool is full, discarding resource: %d\n", resource.ID)
	}
}

// 2. 等待组 (WaitGroup)

// 2.1 基本用法
func WaitGroupExample() {
	var wg sync.WaitGroup

	// 启动多个goroutine
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Second)
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	// 等待所有goroutine完成
	wg.Wait()
	fmt.Println("All workers completed")
}

// 2.2 错误处理
func WaitGroupWithErrorHandling() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := workWithError(id); err != nil {
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

func workWithError(id int) error {
	time.Sleep(time.Millisecond * 100)
	if id == 1 {
		return fmt.Errorf("worker %d failed", id)
	}
	return nil
}

// 3. 屏障 (Barrier)

// 3.1 基本屏障实现
type Barrier struct {
	mu      sync.Mutex
	cond    *sync.Cond
	count   int
	parties int
	phase   int
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

// 3.2 使用示例
func BarrierExample() {
	barrier := NewBarrier(2)

	go func() {
		fmt.Println("Goroutine 1: Phase 1")
		barrier.Await()
		fmt.Println("Goroutine 1: Phase 2")
		barrier.Await()
		fmt.Println("Goroutine 1: Phase 3")
	}()

	go func() {
		fmt.Println("Goroutine 2: Phase 1")
		barrier.Await()
		fmt.Println("Goroutine 2: Phase 2")
		barrier.Await()
		fmt.Println("Goroutine 2: Phase 3")
	}()

	time.Sleep(time.Second * 2)
}

// 4. 组合使用示例

// 4.1 工作池
type WorkPool struct {
	mu      sync.Mutex
	cond    *sync.Cond
	wg      sync.WaitGroup
	barrier *Barrier
	jobs    chan int
	workers int
	active  int
}

func NewWorkPool(workers int) *WorkPool {
	wp := &WorkPool{
		jobs:    make(chan int, 100),
		workers: workers,
	}
	wp.cond = sync.NewCond(&wp.mu)
	wp.barrier = NewBarrier(workers)
	return wp
}

func (wp *WorkPool) Start() {
	// 启动工作协程
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// 等待所有工作协程准备就绪
	wp.barrier.Await()
}

func (wp *WorkPool) worker(id int) {
	defer wp.wg.Done()

	// 等待开始信号
	wp.barrier.Await()

	for job := range wp.jobs {
		wp.mu.Lock()
		wp.active++
		wp.mu.Unlock()

		// 处理工作
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Millisecond * 100)

		wp.mu.Lock()
		wp.active--
		if wp.active == 0 {
			wp.cond.Signal() // 通知所有工作完成
		}
		wp.mu.Unlock()
	}
}

func (wp *WorkPool) AddJob(job int) {
	wp.jobs <- job
}

func (wp *WorkPool) WaitForCompletion() {
	wp.mu.Lock()
	for wp.active > 0 {
		wp.cond.Wait()
	}
	wp.mu.Unlock()

	close(wp.jobs)
	wp.wg.Wait()
}

// 5. 门闩 (Latch) - 一次性同步点

// 5.1 基本门闩实现
// 门闩是一种一次性的同步原语，类似于倒计时器
// 当计数器归零时，所有等待的goroutine都会被唤醒
// 门闩只能被触发一次，之后就不能再重置

type Latch struct {
	mu        sync.Mutex
	cond      *sync.Cond
	count     int
	triggered bool
}

func NewLatch(count int) *Latch {
	l := &Latch{count: count}
	l.cond = sync.NewCond(&l.mu)
	return l
}

// CountDown 减少计数器，当计数器归零时触发门闩
func (l *Latch) CountDown() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.triggered {
		return // 已经触发过了
	}

	l.count--
	if l.count <= 0 {
		l.triggered = true
		l.cond.Broadcast() // 唤醒所有等待的goroutine
	}
}

// Await 等待门闩被触发
func (l *Latch) Await() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for !l.triggered {
		l.cond.Wait()
	}
}

// AwaitWithTimeout 带超时的等待
func (l *Latch) AwaitWithTimeout(timeout time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.triggered {
		return true
	}

	// 创建超时通道
	timeoutChan := time.After(timeout)

	// 等待条件满足或超时
	for !l.triggered {
		select {
		case <-timeoutChan:
			return false // 超时
		default:
			l.cond.Wait()
		}
	}

	return true
}

// IsTriggered 检查门闩是否已经触发
func (l *Latch) IsTriggered() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.triggered
}

// GetCount 获取当前计数器值
func (l *Latch) GetCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.count
}

// 5.2 门闩使用示例
func LatchExample() {
	// 创建一个计数为3的门闩
	latch := NewLatch(3)

	// 启动等待的goroutine
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d waiting for latch\n", id)
			latch.Await()
			fmt.Printf("Goroutine %d released from latch\n", id)
		}(i)
	}

	// 模拟一些工作
	time.Sleep(time.Millisecond * 100)

	// 逐步减少计数器
	fmt.Println("Counting down...")
	latch.CountDown()
	time.Sleep(time.Millisecond * 50)

	latch.CountDown()
	time.Sleep(time.Millisecond * 50)

	latch.CountDown() // 最后一个，触发门闩
	fmt.Println("Latch triggered!")

	time.Sleep(time.Millisecond * 100)
}

// 5.3 带超时的门闩示例
func LatchWithTimeoutExample() {
	latch := NewLatch(2)

	// 启动一个带超时的等待goroutine
	go func() {
		fmt.Println("Waiting for latch with timeout...")
		if latch.AwaitWithTimeout(2 * time.Second) {
			fmt.Println("Latch triggered successfully")
		} else {
			fmt.Println("Latch wait timed out")
		}
	}()

	// 只减少一次计数器
	time.Sleep(time.Millisecond * 100)
	latch.CountDown()

	// 等待超时
	time.Sleep(3 * time.Second)
}

// 5.4 门闩在初始化场景中的应用
type InitializationManager struct {
	latch *Latch
	data  map[string]interface{}
	mu    sync.RWMutex
}

func NewInitializationManager() *InitializationManager {
	return &InitializationManager{
		latch: NewLatch(1), // 只需要一次初始化
		data:  make(map[string]interface{}),
	}
}

func (im *InitializationManager) Initialize() {
	// 模拟初始化工作
	time.Sleep(time.Millisecond * 100)

	im.mu.Lock()
	im.data["initialized"] = true
	im.data["timestamp"] = time.Now()
	im.mu.Unlock()

	fmt.Println("Initialization completed")
	im.latch.CountDown() // 触发门闩
}

func (im *InitializationManager) GetData(key string) (interface{}, bool) {
	// 等待初始化完成
	im.latch.Await()

	im.mu.RLock()
	defer im.mu.RUnlock()

	value, exists := im.data[key]
	return value, exists
}

// 5.5 门闩在资源准备场景中的应用
type ResourcePreparer struct {
	latch     *Latch
	resources []string
	ready     bool
	mu        sync.Mutex
}

func NewResourcePreparer(resourceCount int) *ResourcePreparer {
	return &ResourcePreparer{
		latch:     NewLatch(resourceCount),
		resources: make([]string, 0, resourceCount),
	}
}

func (rp *ResourcePreparer) PrepareResource(resourceID int) {
	// 模拟资源准备
	time.Sleep(time.Millisecond * time.Duration(50+resourceID*10))

	rp.mu.Lock()
	rp.resources = append(rp.resources, fmt.Sprintf("Resource-%d", resourceID))
	rp.mu.Unlock()

	fmt.Printf("Resource %d prepared\n", resourceID)
	rp.latch.CountDown()
}

func (rp *ResourcePreparer) WaitForAllResources() []string {
	rp.latch.Await()

	rp.mu.Lock()
	defer rp.mu.Unlock()

	rp.ready = true
	result := make([]string, len(rp.resources))
	copy(result, rp.resources)
	return result
}

// 5.6 门闩与屏障的对比
/*
门闩 (Latch) vs 屏障 (Barrier):

门闩特点：
- 一次性触发，不能重置
- 等待者被动等待，由外部控制触发
- 适用于初始化、资源准备等场景
- 计数器归零时自动触发

屏障特点：
- 可以重复使用
- 等待者主动参与，到达指定数量时触发
- 适用于阶段同步、并行计算等场景
- 所有参与者到达时自动触发
*/

// 5.7 门闩的最佳实践
func LatchBestPractices() {
	// 1. 总是检查是否已经触发
	latch := NewLatch(1)

	if latch.IsTriggered() {
		fmt.Println("Latch already triggered")
		return
	}

	// 2. 使用超时避免无限等待
	if !latch.AwaitWithTimeout(5 * time.Second) {
		fmt.Println("Latch wait timed out, taking fallback action")
		// 执行备用逻辑
	}

	// 3. 在适当的时候触发门闩
	go func() {
		// 执行初始化工作
		time.Sleep(time.Millisecond * 100)
		latch.CountDown()
	}()

	// 4. 避免重复触发
	latch.CountDown() // 第一次触发
	latch.CountDown() // 后续调用会被忽略
}
