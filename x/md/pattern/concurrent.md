Worker Pool（协程池）
作用：限制并发协程数量，避免资源耗尽。

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        results <- job * 2
    }
}

// 创建任务和结果通道
jobs := make(chan int, 100)
results := make(chan int, 100)

// 启动 3 个 Worker
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}

// 发送任务
for j := 1; j <= 5; j++ {
    jobs <- j
}
close(jobs)

// 收集结果
for a := 1; a <= 5; a++ {
    <-results
}


```