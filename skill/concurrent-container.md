# Go 并发容器 (Concurrent Containers)

Go 并发容器提供线程安全的数据结构，适用于多 goroutine 环境下的数据共享和操作。

## 概述

并发容器通过互斥锁和原子操作保证线程安全，支持泛型，可以存储任意类型的数据。

## 容器类型

### 1. ConcurrentMap - 并发 Map

线程安全的键值对存储容器。

```go
// 创建并发 map
cm := NewConcurrentMap[string, int]()

// 基本操作
cm.Set("key1", 100)
value, exists := cm.Get("key1")  // 100, true
cm.Delete("key1")

// 获取信息
length := cm.Len()
keys := cm.Keys()
values := cm.Values()
```

**特点：**
- 使用读写锁，读操作并发性能好
- 支持任意 comparable 类型的键
- 支持任意类型的值

### 2. ConcurrentSlice - 并发 Slice

线程安全的动态数组。

```go
// 创建并发 slice
cs := NewConcurrentSlice[int]()

// 基本操作
cs.Append(1)
cs.Append(2)
value, ok := cs.Get(0)  // 1, true
cs.Set(1, 100)
cs.Remove(0)

// 获取长度
length := cs.Len()
```

**特点：**
- 支持动态扩容
- 索引访问和修改
- 边界检查保护

### 3. ConcurrentQueue - 并发队列

线程安全的 FIFO 队列。

```go
// 创建并发队列
cq := NewConcurrentQueue[string]()

// 队列操作
cq.Enqueue("task1")
cq.Enqueue("task2")
item, ok := cq.Dequeue()  // "task1", true
peekItem, ok := cq.Peek() // "task2", true

// 状态检查
isEmpty := cq.IsEmpty()
length := cq.Len()
```

**特点：**
- 先进先出 (FIFO) 顺序
- 支持查看队首元素
- 空队列安全处理

### 4. ConcurrentCounter - 并发计数器

高性能的原子计数器。

```go
// 创建计数器
counter := NewConcurrentCounter()

// 计数操作
counter.Increment()
counter.Decrement()
counter.Add(10)
counter.Set(100)

// 获取值
value := counter.Get()
```

**特点：**
- 使用原子操作，性能极高
- 无锁设计
- 适合高频计数场景

### 5. ConcurrentSet - 并发集合

线程安全的无序集合。

```go
// 创建并发集合
set := NewConcurrentSet[int]()

// 集合操作
set.Add(1)
set.Add(2)
exists := set.Contains(1)  // true
set.Remove(1)
set.Clear()

// 获取信息
size := set.Len()
items := set.Items()
```

**特点：**
- 基于 map 实现，查找效率高
- 自动去重
- 支持清空操作

## 使用场景

### 1. 缓存系统
```go
cache := NewConcurrentMap[string, []byte]()
cache.Set("user:123", userData)
data, exists := cache.Get("user:123")
```

### 2. 任务队列
```go
taskQueue := NewConcurrentQueue[Task]()
taskQueue.Enqueue(newTask)

// 工作协程
go func() {
    for {
        if task, ok := taskQueue.Dequeue(); ok {
            processTask(task)
        }
    }
}()
```

### 3. 计数器统计
```go
requestCounter := NewConcurrentCounter()
go func() {
    for {
        requestCounter.Increment()
        // 处理请求
    }
}()
```

### 4. 去重处理
```go
processedIDs := NewConcurrentSet[string]()
if !processedIDs.Contains(id) {
    processedIDs.Add(id)
    // 处理逻辑
}
```

## 性能考虑

### 1. 锁粒度
- ConcurrentMap 使用读写锁，读多写少场景性能好
- ConcurrentCounter 使用原子操作，性能最佳
- 其他容器使用互斥锁，适合中等并发

### 2. 内存使用
- 所有容器都有额外的锁开销
- 预分配容量可以减少扩容开销

### 3. 并发度
- 高并发场景优先使用原子操作
- 中等并发可以使用读写锁
- 低并发场景互斥锁足够

## 最佳实践

### 1. 选择合适的容器
```go
// 高频计数 - 使用原子计数器
counter := NewConcurrentCounter()

// 读多写少 - 使用并发 map
config := NewConcurrentMap[string, string]()

// 任务队列 - 使用并发队列
queue := NewConcurrentQueue[Task]()
```

### 2. 避免长时间持有锁
```go
// 好的做法
func processData(cm *ConcurrentMap[string, int]) {
    data, exists := cm.Get("key")
    if exists {
        // 在锁外处理数据
        process(data)
    }
}

// 避免的做法
func processData(cm *ConcurrentMap[string, int]) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    // 长时间处理数据 - 阻塞其他操作
    time.Sleep(1 * time.Second)
}
```

### 3. 批量操作优化
```go
// 批量添加数据
func batchAdd(cs *ConcurrentSlice[int], items []int) {
    cs.mu.Lock()
    defer cs.mu.Unlock()
    cs.items = append(cs.items, items...)
}
```

### 4. 错误处理
```go
// 检查操作是否成功
if value, ok := cs.Get(index); ok {
    // 使用 value
} else {
    // 处理索引越界
}

if item, ok := cq.Dequeue(); ok {
    // 处理队列项
} else {
    // 处理空队列
}
```

## 注意事项

1. **死锁预防**：避免在持有锁时调用可能获取同一锁的方法
2. **性能监控**：在高并发场景下监控锁竞争情况
3. **内存泄漏**：及时清理不再使用的数据
4. **类型安全**：使用泛型确保类型安全

## 扩展功能

可以根据需要扩展更多功能：

- **ConcurrentStack** - 并发栈
- **ConcurrentPriorityQueue** - 并发优先队列
- **ConcurrentLRUCache** - 并发 LRU 缓存
- **ConcurrentRingBuffer** - 并发环形缓冲区

这些并发容器为 Go 并发编程提供了安全可靠的数据结构基础。 