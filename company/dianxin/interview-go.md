# Go 与并发、性能调优（口述问答）

> 对应 `readme.md` 中「Go 语言」与性能工具相关条目。

## 1) Goroutine 和 OS 线程区别？什么场景容易出问题？

**口述参考：**  
Goroutine 由运行时调度，创建和切换成本远低于大量线程，适合高并发。风险主要是泄漏和失控并发，比如没有退出条件、阻塞在 channel、下游调用无超时。我会用 `context` 控制生命周期，配合 `WaitGroup` 收口，并用 pprof 看 goroutine 数量是否异常增长。

## 2) Channel 用在什么场景？如何避免误用？

**口述参考：**  
我主要用在任务编排、生产者消费者、协程间信号传递。要避免无缓冲导致的死锁、重复 close、发送方和接收方生命周期不一致。习惯上明确谁负责 close，退出统一走 `context`，复杂场景用 `select` 处理取消和超时。

## 3) 你怎么理解 Go 的 interface 与多态？

**口述参考：**  
interface 是隐式实现，利于依赖倒置和单测 mock。我会控制 interface 粒度，避免「大而全」；对外暴露小接口，对内组合。注意 nil interface 与具体类型 nil 的坑，以及值接收者与指针接收者对满足接口的影响。

## 4) 谈谈内存模型里 happens-before 与你在并发里怎么用？

**口述参考：**  
happens-before 描述可见性顺序，比如 channel 发送 happens-before 接收、`sync` 解锁 happens-before 加锁。写代码时我会用 channel 或 `sync` 明确传递「完成信号」，避免靠「碰巧可见」的写法；共享数据尽量用消息传递或受控锁保护。

## 5) GC 带来什么问题？你怎么优化？

**口述参考：**  
分配过频会推高 GC CPU 和延迟抖动。我会用 pprof 的 heap、alloc_space 找热点，减少临时对象、复用 buffer、避免不必要逃逸；关键路径避免在循环里大量分配。优化后看 P99 和 GC pause 是否改善。

## 6) errgroup 和裸开 goroutine 怎么选？

**口述参考：**  
多个并行子任务且要统一错误和取消时，用 errgroup 更合适，可以和一个 `context` 绑定，任一失败取消其它任务。简单 fire-and-forget 且生命周期清晰时才会裸开协程，并保证可退出、可观测。

## 7) singleflight 解决什么问题？有什么副作用？

**口述参考：**  
同一 key 的并发请求合并成一次执行，适合缓存击穿、热点配置拉取。副作用是可能拉长个别请求的等待时间，以及要处理好 panic 和超时，避免拖垮共享调用。我会设超时和限流，并监控合并命中率。

## 8) context 的适用边界是什么？

**口述参考：**  
适合传递取消、超时、请求级元数据（如 trace id），不适合塞业务大对象。取消要一路传到下游 HTTP/RPC/DB。注意 `context.Background` 只在根节点用，子调用用派生 context，避免泄漏取消链。

## 9) pprof、trace、race 分别什么时候用？

**口述参考：**  
CPU profile 看热点函数；heap 看分配和泄漏嫌疑；trace 看调度、阻塞和网络等待的时序；race 在开发阶段跑并发测试抓数据竞争。线上一般用连续 profile 或采样，结合监控定位再本地复现。

## 10) 如何避免 Goroutine 泄漏？（收口清单）

**口述参考：**  
每条协程都要有退出条件；channel 配对清晰；外部 IO 必须带超时；后台任务用独立 supervisor 或 worker pool 限制并发；上线后看 goroutine 指标和 pprof goroutine 堆栈是否持续增长。
