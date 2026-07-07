# Golang 面试基础

## 语言基础

### Slice 底层结构
- Slice 由指针、长度(len)、容量(cap)组成
- 多个 slice 可能共享同一个底层数组
- append 触发扩容时，会分配新数组并拷贝，不再共享底层数组
- 扩容策略：Go 1.18+ 按增长因子平滑递增

### Map 底层原理
- 哈希表实现，bucket 大小为 8 个 key-value 对
- 溢出桶(overflow bucket)链表解决哈希冲突
- 扩容机制：渐进式搬迁（rehash 或 overflow 过多时触发）
- 读取、写入、删除均为 O(1) 平均复杂度
- **非并发安全**：并发写/读写会 panic，需用 sync.Mutex 或 sync.Map

### Channel 机制
- 无缓冲 channel：同步通信，发送方阻塞直到接收方就绪
- 有缓冲 channel：缓冲满时才阻塞发送方
- 底层结构：hchan 包含环形队列 buf、sendq/recvq 等待队列
- 关闭 channel：close 后仍可读，不可写（panic），不可重复关闭（panic）
- select 多路复用，随机选择一个就绪的 case
- nil channel 的发送和接收永久阻塞

### Interface 底层
- 空接口 `interface{}`：eface 结构（type, data 指针）
- 非空接口：iface 结构（tab 指向 itable, data 指针）
- **nil 判断陷阱**：接口值为 nil 需要 type 和 data 都为 nil
- 类型断言：`v, ok := x.(T)` 安全断言

### defer 执行顺序
- 后进先出（LIFO）顺序执行
- 注册时参数即求值，执行时使用注册时的参数值
- Go 1.14+ 改为内联优化，性能大幅提升
- recover 必须放在 defer 函数中才能捕获 panic

### Goroutine 调度
- GMP 模型：G(Goroutine)、M(Machine-OS线程)、P(Processor-逻辑处理器)
- P 的数量由 `GOMAXPROCS` 决定，默认等于 CPU 核数
- 工作窃取(work-stealing)：空闲 P 从其他 P 偷取 G
- 抢占式调度：Go 1.14+ 基于信号的异步抢占

### 内存管理
- TCMalloc 思想：多级缓存（mcache → mcentral → mheap）
- 逃逸分析：编译器决定变量分配在栈还是堆
- GC：三色标记 + 并发回收 + 混合写屏障

---

## 并发编程

### sync 包
- Mutex：互斥锁，Lock/Unlock 成对使用
- RWMutex：读写锁，读共享，写独占
- WaitGroup：等待一组 goroutine 完成
- Once：确保函数只执行一次
- Cond：条件变量，Broadcast/Signal/Wait

### Context 使用
- 传递取消信号、超时控制、请求域数据
- 父 context 取消会传播到所有子 context
- 不要将 context 存到 struct 中长期持有
- 不要传 nil context，不确定时用 `context.TODO()`

### 常见并发模式
- Fan-in/Fan-out
- Pipeline
- Worker Pool
- Or-Done Channel
- Tee Channel

---

## 常见面试题

### 1. Slice 扩容机制是怎样的？append 后原 slice 会变吗？

**扩容机制**：当 append 超过 cap 时，Go 会分配新的底层数组，将原数据拷贝过去，然后追加新元素。新容量在 Go 1.18+ 的扩容策略为：
- 原容量 < 256：扩容为 2 倍
- 原容量 >= 256：按 `(oldcap + 3*256) / 4` 公式平滑递增

**append 后原 slice 会变吗？** 分两种情况：
- **未触发扩容**（len + 新元素 <= cap）：在原底层数组上追加，如果多个 slice 共享该底层数组，其他 slice 的值也会受影响。
- **触发扩容**（len + 新元素 > cap）：分配新数组，返回的 slice 指向新数组，与原 slice 不再共享底层数组，不再互相影响。

```go
a := []int{1, 2, 3, 4, 5} // len=5, cap=5
b := a[:3]                 // len=3, cap=5，共享底层数组
b = append(b, 99)          // 未扩容，b[3]=99，影响了 a[3]
fmt.Println(a)             // [1 2 3 99 5]

c := a[:]
c = append(c, 6, 7, 8, 9, 10, 11) // 超过 cap，触发扩容，新数组
fmt.Println(a)                      // [1 2 3 99 5]，不受影响
```

---

### 2. Map 并发读写会怎样？如何安全并发访问？

**并发读写/写写会触发运行时的 fatal error**，直接 crash 进程：
> `fatal error: concurrent map writes` 或 `concurrent map read and map write`

这是 Go runtime 检测到并发访问后主动抛出的，无法通过 recover 捕获。

**解决方案（按场景选择）**：

| 方案 | 适用场景 |
|------|----------|
| `sync.Mutex` / `sync.RWMutex` | 通用方案，读写都需要保护 |
| `sync.Map` | 读多写少，key 相对固定 |
| 将 map 限制在单个 goroutine 内 | channel 传递消息的模式 |

```go
// 方案一：RWMutex（读多写少）
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}
func (s *SafeMap) Get(key string) int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.m[key]
}
func (s *SafeMap) Set(key string, val int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[key] = val
}

// 方案二：sync.Map
var m sync.Map
m.Store("key", "value")
v, ok := m.Load("key")
m.Range(func(k, v interface{}) bool { return true })
```

---

### 3. Channel 关闭后还能读吗？向已关闭 channel 写会怎样？

| 操作 | 结果 |
|------|------|
| 向未初始化(nil) channel 读/写 | 永久阻塞 |
| 向已关闭 channel 写 | **panic**: `send on closed channel` |
| 从已关闭 channel 读（有缓冲数据） | 读到缓冲数据，ok=true |
| 从已关闭 channel 读（无缓冲数据） | 读到零值，ok=false |
| 重复关闭 channel | **panic**: `close of closed channel` |
| 关闭 nil channel | **panic**: `close of nil channel` |

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)
v, ok := <-ch // v=1, ok=true
v, ok = <-ch  // v=2, ok=true
v, ok = <-ch  // v=0, ok=false，可无限读，一直返回零值

// 使用 for-range 读取直到 channel 关闭
for v := range ch { }

// 检测 channel 是否关闭
v, ok := <-ch
if !ok { /* 已关闭 */ }
```

**最佳实践**：
- 由发送方关闭 channel
- 多个发送方时，用 `sync.WaitGroup` + 专门的 closer goroutine 关闭
- 不确定是否已关闭时，加 `recover` 或使用 select

---

### 4. 如何判断 interface 是否为 nil？

interface 底层结构包含两部分：**type**（动态类型）和 **data**（动态值指针）。只有两者都为 nil，接口才等于 nil。

```go
var p *int = nil          // p 是 nil 指针
var i interface{} = p     // i 不为 nil！因为 type = *int, data = nil
fmt.Println(i == nil)     // false

var i2 interface{} = nil  // type = nil, data = nil
fmt.Println(i2 == nil)    // true

// 实际项目中常见坑
func getError() error {
    var err *MyError = nil
    return err // 返回的 error != nil！
}
err := getError()
if err != nil { // 会进入此分支，即使 err 的值是 nil
}
```

**正确做法**：返回 error 时直接 `return nil`，不要返回类型化的 nil 指针。

---

### 5. defer 执行顺序？defer 中修改返回值有效吗？

**执行顺序**：后进先出（LIFO，即栈序），最后注册的 defer 最先执行。

**defer 修改返回值**——只在**命名返回值**的情况下有效：

```go
// 命名返回值：defer 可以修改
func f1() (result int) {
    defer func() { result++ }() // 有效：返回 1
    return 0
}

// 匿名返回值：defer 无法修改
func f2() int {
    result := 0
    defer func() { result++ }() // 无效：仍返回 0
    return result
}

// defer 参数预求值
func f3() {
    x := 1
    defer fmt.Println(x) // 打印 1，参数在注册时已求值
    x = 2
}
```

**关键**：命名返回值时，return 语句分两步：1) 将返回值赋给命名变量；2) 执行 defer 链；3) 真正的返回。所以 defer 可以修改命名返回值。

---

### 6. GPM 调度模型是什么？goroutine 阻塞会影响其他吗？

**GMP 模型**：

| 组件 | 含义 | 说明 |
|------|------|------|
| **G** (Goroutine) | 协程 | 轻量级线程，初始栈 2KB |
| **P** (Processor) | 逻辑处理器 | 本地任务队列，数量 = GOMAXPROCS |
| **M** (Machine) | OS 线程 | 实际执行单元，与 P 一对一绑定 |

**调度流程**：M 必须绑定 P 才能执行 G。每个 P 有本地 G 队列，M 从绑定的 P 的队列中取 G 执行。P 的本地队列空了，就从其他 P 偷取一半 G（work-stealing）。

**goroutine 阻塞是否影响其他 goroutine？** 分两种情况：

- **用户态阻塞**（channel、锁、select 等）：G 进入等待队列，M 与 P 解绑，P 被另一个 M 或新建 M 接管，继续执行其他 G。不影响。
- **系统调用阻塞**（read/write 等）：M 与 P 解绑，P 被另一个 M 接管。阻塞的 M 和 G 一起挂起，系统调用返回后 G 重新排队等待 P。

所以 goroutine 阻塞不会阻塞其他 goroutine，这是 GMP 模型的核心优势。

---

### 7. Go GC 如何工作？如何优化 GC？

**三色标记清除 + 并发回收**：

1. **初始标记**（STW 短）：扫描栈和全局变量，标记根对象
2. **并发标记**（并发）：从根对象出发，三色标记遍历对象图
3. **混合写屏障**（并发）：记录 GC 期间指针变化，保证不漏标
4. **并发清除**（并发）：回收白色对象

**优化 GC 的方法**：

| 方向 | 具体措施 |
|------|----------|
| 减少堆分配 | 用值类型代替指针、预分配 slice 容量、strings.Builder |
| 减少指针 | 指针越少，标记越快；可考虑用索引替代指针 |
| 对象复用 | `sync.Pool` 复用高频临时对象 |
| 调整 GOGC | 默认 100（堆增长 100% 触发 GC），调大可减少 GC 频率 |
| 设置 GOMEMLIMIT | Go 1.19+，软内存上限，避免无限增长 |
| 减少分配频率 | 避免在热路径中创建大量临时对象 |

```bash
# 查看 GC 日志
GODEBUG=gctrace=1 ./app

# 示例输出解读
gc 1 @0.012s 2%: 0.12+0.45+0.013 ms clock, 0.12+0.23/0.35/0.15+0.013 ms cpu
#           GC 1  扫描STW  并发标记    STW标记结束   GC 占 CPU 时间
```

---

### 8. 逃逸分析是什么？如何查看变量是否逃逸？

**逃逸分析**是编译器的一项优化，决定变量分配到**栈**还是**堆**：
- **栈分配**：函数返回时自动回收，无 GC 开销，速度快
- **堆分配**：需要 GC 回收，有分配锁竞争和 GC 扫描开销

**常见逃逸场景**：
```go
// 1. 返回局部变量的指针
func f() *int { x := 1; return &x } // x 逃逸到堆

// 2. 变量被 interface{} 包裹
func f() { fmt.Println(1) } // 1 逃逸（fmt.Println 接收 interface{}）

// 3. 变量大小不确定或过大
func f() { make([]int, 100000) } // 大对象逃逸

// 4. slice/map 中存指针
s := []*int{&x} // 可能导致底层数组整体逃逸

// 5. 闭包引用
func f() func() { x := 1; return func() { x++ } } // x 逃逸
```

**查看逃逸分析结果**：
```bash
go build -gcflags="-m" main.go          # 查看逃逸分析
go build -gcflags="-m -m" main.go       # 更详细的信息
go tool compile -m main.go              # 单个文件
```

---

### 9. new 和 make 的区别？

| 特性 | `new(T)` | `make(T, args)` |
|------|----------|-----------------|
| 适用范围 | 任意类型 | **仅** slice, map, channel |
| 返回值 | `*T`（指向零值 T 的指针） | `T`（初始化后的 T 本身） |
| 是否初始化 | 只分配内存，填零值 | 初始化内部数据结构（如 slice 的底层数组） |
| 可直接使用 | 可以（零值可用） | 可以，make 后即可使用 |

```go
// new：返回指针
p := new(int)       // p 是 *int，*p = 0
s := new([]int)     // s 是 *[]int，*s = nil，不能直接 append

// make：返回初始化后的值
m := make([]int, 0, 10)  // 可立即 append
c := make(chan int, 5)   // 可立即 send/receive
mp := make(map[string]int) // 可立即写

// 关键区别
var m1 map[string]int     // nil map，写会 panic
m2 := make(map[string]int) // 初始化后，可以安全使用

var ch1 chan int           // nil channel，读写永久阻塞
ch2 := make(chan int)      // 可用 channel
```

---

### 10. 数组和 slice 的区别？函数传参是值传递还是引用传递？

| 特性 | 数组 `[n]T` | Slice `[]T` |
|------|------------|------------|
| 长度 | 固定，编译时确定，是类型的一部分 | 动态，运行时可变 |
| 值/引用 | 值类型，赋值/传参会拷贝整个数组 | 引用语义（header 拷贝，共享底层数组） |
| 内存 | 整个数组连续分配 | header(24B) + 底层数组 |
| 比较 | 可用 `==` 比较（元素可比较时） | 只能和 nil 比较 `== nil` |
| 零值 | 所有元素为零值 | nil（但 `len(nil slice)==0`） |

**函数传参：Go 只有值传递，没有引用传递。**

```go
// 数组：整个数组被拷贝，函数内修改不影响原数组
func modifyArray(arr [3]int) {
    arr[0] = 999 // 不影响调用方的数组
}

// Slice：header 被拷贝（指针、len、cap），但底层数组共享
func modifySlice(s []int) {
    s[0] = 999 // 修改底层数组，会影响调用方
    s = append(s, 1) // header 的修改不影响调用方（重新赋值）
}

// Map 和 Channel：传递的是底层指针的拷贝，函数内修改会影响调用方

// 如果想修改 slice 本身（如 append），有两种方式
// 方式一：返回新 slice
func add(s []int, v int) []int { return append(s, v) }
// 方式二：传指针
func add(s *[]int, v int) { *s = append(*s, v) }
```
