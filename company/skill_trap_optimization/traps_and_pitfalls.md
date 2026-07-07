# Go 陷井踩坑

## Slice 陷阱

### 1. append 后原 slice 被覆盖
```go
a := []int{1, 2, 3}
b := a[:2]
b = append(b, 99) // b 和 a 共享底层数组，a[2] 变成 99
```
**避坑**：不确定是否共享底层数组时，用 `copy` 或 `append([]int{}, a...)` 创建独立 slice。

### 2. slice 作为参数
```go
func modify(s []int) {
    s = append(s, 4) // 可能扩容后指向新数组，不影响原 slice
}
```
**避坑**：如需修改 slice 本身（非元素），传指针 `*[]int` 或返回新 slice。

### 3. for-range 擦除 slice 元素
删除元素后索引不回溯会跳过下一个元素。应使用递减循环或收集后批量删除。

---

## Map 陷阱

### 1. 并发读写 panic
```go
// fatal error: concurrent map writes
go func() { m["a"] = 1 }()
go func() { m["a"] = 2 }()
```
**解决**：`sync.Mutex`、`sync.RWMutex` 或读多写少场景用 `sync.Map`。

### 2. nil map 写入 panic
```go
var m map[string]int
m["key"] = 1 // panic: assignment to entry in nil map
```
**避坑**：声明时用 `make(map[string]int)` 或 `map[string]int{}`。

### 3. map 遍历顺序不确定
Go 故意随机化 map 遍历顺序，不要依赖顺序。

---

## Channel 陷阱

### 1. 未初始化 nil channel 永久阻塞
```go
var ch chan int
ch <- 1  // 永久阻塞
<-ch     // 永久阻塞
```
**避坑**：总是用 `make(chan T)` 初始化 channel。

### 2. 向已关闭 channel 发送数据
```go
close(ch)
ch <- 1 // panic: send on closed channel
```
**避坑**：由发送方关闭 channel，用 select 或 sync.Once 保护 close。

### 3. 无缓冲 channel 死锁
```go
ch := make(chan int)
ch <- 1 // 死锁：没有接收方
<-ch    // 仅在函数内单个 goroutine 时
```

### 4. goroutine 泄漏
goroutine 阻塞在 channel 上永远不退出。用 context 或 done channel 控制退出。

---

## Interface 陷阱

### 1. nil 接口 != nil 具体类型
```go
var p *MyStruct = nil
var i interface{} = p
fmt.Println(i == nil) // false! 接口的 type 不为 nil
```
**避坑**：返回具体 nil 指针给接口时，直接返回 `nil` 而不是类型化的 nil。

---

## defer 陷阱

### 1. defer 参数预求值
```go
func f() {
    x := 1
    defer fmt.Println(x) // 注册时 x=1，输出 1
    x = 2
}
```

### 2. defer 与命名返回值
```go
func f() (result int) {
    defer func() { result++ }()
    return 0 // result = 0; defer 执行后 result = 1; 返回 1
}
```

### 3. defer 在循环中
```go
for _, f := range files {
    f, _ := os.Open(f) // defer 在循环结束后才执行
    defer f.Close()     // 大量文件描述符积压
}
```
**避坑**：将循环体提取为函数，或在循环内不依赖 defer 直接关闭。

---

## for-range 陷阱

### 1. 循环变量捕获（Go 1.21 已修复）
```go
// Go 1.21 之前
for _, v := range items {
    go func() { fmt.Println(v) }() // 所有 goroutine 可能打印最后一个值
}
// 修复方式：v := v 或 Go 1.22+ 自动修复
```

### 2. 值拷贝
```go
for _, item := range items { // item 是拷贝
    item.Field = "new"       // 不影响 items 中的原始数据
}
```
**避坑**：使用索引 `items[i]` 或 `for i := range` 修改元素。

---

## 错误处理陷阱

### 1. error 接口判断
```go
func getError() error { return (*MyError)(nil) }
err := getError()
if err != nil { // true! err 的 type 不为 nil
}
```
**避坑**：返回 error 时直接 return nil，不要返回类型化的 nil 指针。

### 2. 忽略 defer 中的 error
```go
defer f.Close() // Close 的错误被忽略
```
**避坑**：将 defer 包装为匿名函数处理 error，或使用命名返回值捕获。

---

## 其他陷阱

### 闭包与循环变量
```go
for i := 0; i < 3; i++ {
    defer func() { fmt.Print(i) }() // 输出 3, 3, 3（defer 匿名函数捕获的是 i 的引用）
}
// vs
for i := 0; i < 3; i++ {
    defer fmt.Print(i) // 输出 2, 1, 0（defer 语句注册时 i 值已求值）
}
```

### time.After 内存泄漏
```go
for {
    select {
    case <-time.After(time.Second): // 每次创建新 timer，不会被 GC 直到超时
    }
}
```
**避坑**：用 `time.NewTimer` + `timer.Reset` 复用。

### json.Unmarshal 后的 slice 扩容
目标 slice 可能被复用而非重新分配，需注意已有内容的处理。
