package mianshiskill

import (
	"fmt"
	"sync"
)

type TaskWithId struct {
	ID int
}

func TaskWithWorkerPool() {
	maxWorkers := 5
	totalTasks := 100

	var wg sync.WaitGroup
	taskQueue := make(chan TaskWithId, totalTasks)
	// 启动worker池
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, taskQueue)
	}
	// 发送任务
	for i := 0; i < totalTasks; i++ {
		taskQueue <- TaskWithId{ID: i}
	}
	close(taskQueue)
	wg.Wait()
}

func worker(workerID int, wg *sync.WaitGroup, tasks <-chan TaskWithId) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", workerID, task.ID)
	}
}
