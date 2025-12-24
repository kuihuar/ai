### 实现并发任务调度器，限制运行时的goroutine为5个，处理100个耗时任务
>**问题分析：**
1. 最大并发数限制
2. 耗时任务处理
3. 100个任务不算大

>**实现方式：**
1. 信号量模式，控制并发数（任务间独立性高，无状态）
2. 工作池模式，处理任务（任务处理逻辑统一，资源复用）

>**两种方法的对比：**
| 特性 | 信号量模式 | Worker Pool 模式 |
| --- | --- | --- |
| Goroutine 数量 | 创建 100 个 Goroutine（轻量） | 固定 5 个 Worker Goroutine |
| 任务分配 | 每个任务独立 Goroutine | Worker 复用，循环处理多个任务 |
| 资源开销 | 略高（大量 Goroutine 等待信号量） | 更低（复用固定数量的 Goroutine） |
| 适用场景 | 任务间独立性高、无状态 | 任务处理逻辑统一、需资源复用 |


>**实现代码：**
```go
package main
import (
	"fmt"
	"sync"
	"time"
)
// 任务处理函数
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
	}    
}
func main() {
	// 任务队列
	jobs := make(chan int, 100)
	// 等待组
	var wg sync.WaitGroup
	// 并发限制
	concurrency := 5
	// 启动工作池
	for w := 1; w <= concurrency; w++ {
		wg.Add(1) 
		go worker(w, jobs, &wg)
	}
	// 发送任务
	for j := 1; j <= 100; j++ {
		jobs <- j
	}
	close(jobs)
	// 等待所有任务完成
	wg.Wait()

	fmt.Println("All jobs completed")

}
```
