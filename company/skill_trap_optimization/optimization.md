# Go 性能优化

## CPU 优化

### 逃逸分析优化
- 使用 `go build -gcflags="-m"` 查看逃逸分析结果
- 避免返回局部变量指针导致堆分配
- 避免将变量赋值给 interface{} 导致逃逸
- slice/map 存指针导致底层数组逃逸到堆

### 字符串与 []byte 转换
- 避免频繁的 string/[]byte 转换，使用 `unsafe` 零拷贝转换（确定安全时）
- 使用 `strings.Builder` 替代 `+` 拼接，预分配容量
- `bytes.Buffer` 用于 []byte 拼接

### 循环优化
- 预分配 slice/map 容量，减少扩容开销
- 避免在热循环中使用 defer
- for-range 遍历大数组用指针或索引，避免值拷贝
- 将不变的条件判断提取到循环外部

### 反射优化
- 反射性能差，热路径中避免使用
- 缓存 `reflect.TypeOf` 和 `reflect.ValueOf` 的结果
- 考虑代码生成替代反射

---

## 内存优化

### 结构体内存对齐
- 字段按类型大小降序排列，减少 padding
- 使用 `unsafe.Sizeof` 和 `unsafe.Alignof` 检查
- 使用 `go tool compile -S` 或 `structlayout` 工具分析

### 对象复用
- `sync.Pool` 复用高频创建的对象，减少 GC 压力
- 对象池适用于无状态、可复用的临时对象
- 注意 Pool 中的对象可能被 GC 自动清理

### 减少内存分配
- 使用 `[]byte` 替代 `string` 做读写缓冲区
- 小对象优先值类型（栈分配），大对象才用指针
- `make` 时给出准确容量，减少 slice 扩容

### 栈 vs 堆
- 栈分配：快，自动回收，无 GC 开销
- 堆分配：触发 GC，有分配锁竞争
- 逃逸分析决定分配位置

---

## GC 优化

### 减少 GC 压力
- 减少堆上的对象数量和大小
- 减少指针数量（指针越多，标记越慢）
- `GOGC` 参数调整 GC 触发阈值，默认 100
- `GOMEMLIMIT` 设置软内存上限（Go 1.19+）

### GC 监控
```bash
GODEBUG=gctrace=1 ./app
```
- 使用 `runtime.ReadMemStats` 监控内存指标
- `go tool pprof` 分析 heap profile

---

## 并发优化

### 锁优化
- 减小锁粒度：分段锁、读写锁
- 避免在持有锁时做 IO 操作
- `sync.Map` 适用于读多写少的场景
- 原子操作 `sync/atomic` 替代互斥锁

### Goroutine 管理
- 控制 goroutine 数量，避免 goroutine 泄漏
- 使用 Worker Pool 限制并发数
- 用 `runtime.NumGoroutine()` 监控泄漏

### Channel 优化
- 选择合适的缓冲大小
- 避免在 select 中有永不就绪的 case
- struct{} 做信号 channel，零内存占用

---

## 工具链

| 工具 | 用途 |
|------|------|
| `go test -bench` | 基准测试 |
| `go test -cpuprofile` | CPU profiling |
| `go test -memprofile` | 内存 profiling |
| `go tool pprof` | 分析 profile 数据 |
| `go tool trace` | 运行时追踪 |
| `benchstat` | 基准测试结果对比 |

### pprof 常用命令
```bash
# CPU profiling
go test -bench=. -cpuprofile=cpu.prof
go tool pprof -http=:8080 cpu.prof

# 内存 profiling
go test -bench=. -memprofile=mem.prof
go tool pprof -http=:8080 mem.prof
```
