package mianshiskill

import (
	"fmt"
	"sync"
)

// 任务调度器, 工作池模式
// Worker Pool 模式（任务分发）
// 预先启动固定数量的 Worker Goroutine，通过 Channel 分发任务，避免频繁创建 Goroutine。

type Task interface {
	Execute(workerId int) error
	PrintTask() string
}
type DemoTask struct {
	Id   int
	Name string
}

func (t DemoTask) Execute(workerId int) error {
	fmt.Printf("taskId: %d, workerId: %d\n", t.Id, workerId)
	return nil
}
func (t DemoTask) PrintTask() string {
	return fmt.Sprintf("taskId: %d\n", t.Id)
}

type SchedulerWorkerPool struct {
	tasks        chan Task
	wg           sync.WaitGroup
	maxWorkerNum int
}

func NewSchedulerWorkerPool(maxWorkerNum, maxTaskNum int) *SchedulerWorkerPool {
	return &SchedulerWorkerPool{
		tasks:        make(chan Task, maxTaskNum),
		maxWorkerNum: maxWorkerNum,
	}
}

func (swp *SchedulerWorkerPool) Run() {
	for i := 0; i < swp.maxWorkerNum; i++ {
		swp.wg.Add(1)
		workerId := i
		go func(int) {
			defer swp.wg.Done()
			for task := range swp.tasks {
				task.Execute(workerId)
			}
		}(workerId)
	}
}
func (swp *SchedulerWorkerPool) AddTask(task Task) {
	swp.tasks <- task
}
func (swp *SchedulerWorkerPool) WaitAndClose() {
	close(swp.tasks)
	swp.wg.Wait()
}
func UseTaskWithWorkerPool() {
	const numWorkers = 5
	const numTasks = 1000
	swp := NewSchedulerWorkerPool(numWorkers, numTasks)
	swp.Run()
	for i := 0; i < 100; i++ {
		swp.AddTask(DemoTask{Id: i})
	}
	swp.WaitAndClose()
	fmt.Println("所有任务已完成")

}
