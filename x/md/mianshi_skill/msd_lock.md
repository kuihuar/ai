### 简述高并发中锁的适用场景和效率对比，


|特性| sync.Mutex| sync.RWMutex|sync.Map|
|---|---|---|---|
|锁类型|互斥锁|读写锁(读共享，写独占)|并发安全的map，无显式锁操作|
|适用场景| 高竞争写操作| 读多写少| 读多写少，且键相对稳定|
|内存占用| 低| 中| 较高|
|性能| 写性能高| 读操作并发必好| 读操作快，写操作慢|
|复杂度| 手动管理锁| 需要区分读/写锁| 自动处理|
|死锁产生| 容易产生死锁| 容易产生死锁| 自动检测|




### 死锁的产生的排查方式

1. 死锁产生的场景
- channel发送/接收不匹配（解决方案为使用缓冲channel）
```go
func main(){
    ch := make(chan int)
    ch<-10 //发送阻塞（没有接收方）
    fmt.Println(<-ch)
}
```
- 互斥锁未释放, 重复加锁(解决方案为 defer mu.Unlock()) 
```go
func main(){
    var mu sync.Mutex
    mu.Lock()
    mu.Lock() //死锁
}
```
- waitgroup 等待组未释放（解决方案为 defer wg.Done()）
```go
func main(){
    var wg sync.WaitGroup
    wg.Add(2)
    go func(){
        defer wg.Done()
        fmt.Println("goroutine 1")
    }()
    wg.Wait() //死锁
}
```
### 排查方式（压测验证和静态分析）
>优先利用运行时自检机制，结合可视化工具分析并发关系，最后通过代码规范预防死锁发生
1. go test -race -count=100 -run TestDeadlock 
2. go vet -atomic -copylocks ./... 
3. go install github.com/divan/gobenchui@latest && go test -bench . | gobenchui


