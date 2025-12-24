package mianshiskill

import (
	"fmt"
	"sync"
	"time"
)

func TaskWithWorkerSem() {
	const (
		maxConcurrent = 5  // 最大并发数
		totalTasks    = 20 // 总任务数
	)

	sem := make(chan struct{}, maxConcurrent) // 信号量通道
	var wg sync.WaitGroup

	for i := 0; i < totalTasks; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			sem <- struct{}{}        // 获取信号量（阻塞直到有空位）
			defer func() { <-sem }() // 释放信号量

			fmt.Printf("Task %d started\n", taskID)
			time.Sleep(1 * time.Second) // 模拟耗时操作
			fmt.Printf("Task %d done\n", taskID)
		}(i)
	}

	wg.Wait()
	fmt.Println("All tasks completed")
}
