package skill

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// 整数类型
var int32Val int32
var int64Val int64
var uint32Val uint32
var uint64Val uint64

// 指针类型
var ptrVal unsafe.Pointer

// 布尔类型（通过int32实现）
var boolVal int32

// 原子地加载值
// val := atomic.LoadInt32(&int32Val)
// ptr := atomic.LoadPointer(&ptrVal)

// // 原子地存储值
// atomic.StoreInt32(&int32Val, 100)
// atomic.StorePointer(&ptrVal, unsafe.Pointer(&someStruct))
// // 原子地交换值，返回旧值
// oldVal := atomic.SwapInt32(&int32Val, 200)

// // 原子地比较并交换，成功返回true
// swapped := atomic.CompareAndSwapInt32(&int32Val, 100, 200)

// 1. 原子计数器

type AtomicCounter struct {
	value int64
}

func (ac *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&ac.value, 1)
}

func (ac *AtomicCounter) Decrement() int64 {
	return atomic.AddInt64(&ac.value, -1)
}

func (ac *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&ac.value)
}

func (ac *AtomicCounter) Set(val int64) {
	atomic.StoreInt64(&ac.value, val)
}

// 2. 原子布尔标志
type AtomicFlag struct {
	flag int32
}

func (af *AtomicFlag) Set() {
	atomic.StoreInt32(&af.flag, 1)
}

func (af *AtomicFlag) Clear() {
	atomic.StoreInt32(&af.flag, 0)
}

func (af *AtomicFlag) IsSet() bool {
	return atomic.LoadInt32(&af.flag) == 1
}

func (af *AtomicFlag) TrySet() bool {
	return atomic.CompareAndSwapInt32(&af.flag, 0, 1)
}

// 3. 原子指针

type AtomicReference struct {
	ptr unsafe.Pointer
}

func (ar *AtomicReference) Set(value interface{}) {
	atomic.StorePointer(&ar.ptr, unsafe.Pointer(&value))
}
func (ar *AtomicReference) Get() interface{} {
	// 1. ptr - 原始指针
	ptr := atomic.LoadPointer(&ar.ptr)
	if ptr == nil {
		return nil
	}
	// 2. (*interface{})(ptr) - 类型转换
	// (*interface{})：这是一个指向 interface{} 的指针类型
	// (ptr)：将 unsafe.Pointer 转换为 *interface{} 类型
	// 结果：得到一个指向 interface{} 的指针

	// 3.*(*interface{})(ptr) - 解引用
	// 最外层的 *：对指针进行解引用操作
	// 结果：获取指针指向的实际 interface{} 值

	return *(*interface{})(ptr)
}

// 4. 原子高级
// 4.1 原子自旋锁
type AtomicSpinLock struct {
	locked int32
}

func (asl *AtomicSpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&asl.locked, 0, 1) {
		// 自旋等待，避免上下文切换
		runtime.Gosched()
	}
}

func (asl *AtomicSpinLock) Unlock() {
	atomic.StoreInt32(&asl.locked, 0)
}

func (asl *AtomicSpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&asl.locked, 0, 1)
}

// 4.2 原子读写锁
type AtomicRWLock struct {
	readers int32
	writers int32
}

func (arl *AtomicRWLock) RLock() {
	for {
		// 等待写锁释放
		for atomic.LoadInt32(&arl.writers) > 0 {
			runtime.Gosched()
		}

		// 尝试增加读者计数
		atomic.AddInt32(&arl.readers, 1)
		if atomic.LoadInt32(&arl.writers) == 0 {
			return // 成功获取读锁
		}

		// 写锁被获取，回退读者计数
		atomic.AddInt32(&arl.readers, -1)
	}
}

func (arl *AtomicRWLock) RUnlock() {
	atomic.AddInt32(&arl.readers, -1)
}

func (arl *AtomicRWLock) Lock() {
	// 等待所有读者和写者完成
	for !atomic.CompareAndSwapInt32(&arl.writers, 0, 1) {
		runtime.Gosched()
	}

	for atomic.LoadInt32(&arl.readers) > 0 {
		runtime.Gosched()
	}
}

func (arl *AtomicRWLock) Unlock() {
	atomic.StoreInt32(&arl.writers, 0)
}

// 4.3 原子对像池
type AtomicPool struct {
	head unsafe.Pointer
}

type poolNode struct {
	value interface{}
	next  unsafe.Pointer
}

func (ap *AtomicPool) Push(value interface{}) {
	node := &poolNode{value: value}
	for {
		head := atomic.LoadPointer(&ap.head)
		node.next = head
		if atomic.CompareAndSwapPointer(&ap.head, head, unsafe.Pointer(node)) {
			return
		}
	}
}

func (ap *AtomicPool) Pop() (interface{}, bool) {
	for {
		head := atomic.LoadPointer(&ap.head)
		if head == nil {
			return nil, false
		}

		node := (*poolNode)(head)
		next := node.next

		if atomic.CompareAndSwapPointer(&ap.head, head, next) {
			return node.value, true
		}
	}
}

// 5. 原子操作最佳实践

// 5.1 内存对齐
// 确保原子操作的内存对齐
// 性能：对齐的数据访问更快
// 正确性：某些CPU架构要求原子操作必须对齐
// 缓存效率：避免伪共享（false sharing）

type AlignedCounter struct {
	_     [0]int64 // 填充，确保对齐
	value int64
	_     [0]int64 // 填充，避免伪共享
}

// 5.2 避免ABA问题

// 使用版本号避免ABA问题
// 可能导致逻辑错误
// 在复杂数据结构中特别危险
// 难以调试和发现
// ABA问题示例
// var value int64 = 100

// 线程1：想要将100改为200
// old := atomic.LoadInt64(&value) // 读取到100
// 此时线程2将100改为300，然后又改回100
// 线程1继续执行
// atomic.CompareAndSwapInt64(&value, old, 200) // 成功！但这是错误的
// 解决方案是:
// 1. 使用版本号, 同时比较值和版本号
// 2. 使用指针

type VersionedValue struct {
	value   int64
	version int64
}

func (vv *VersionedValue) CompareAndSwap(expected, new int64, expectedVersion int64) bool {
	// 需要同时比较值和版本号
	// 这里简化实现，实际需要更复杂的逻辑
	return atomic.CompareAndSwapInt64(&vv.value, expected, new)
}

// 5.3 原子操作的性能优化

// 批量操作比单个操作更高效
type BatchCounter struct {
	counters [8]int64 // 避免伪共享
}

func (bc *BatchCounter) AddBatch(values []int64) {
	for i, v := range values {
		atomic.AddInt64(&bc.counters[i%8], v)
	}
}

// type Problem struct {
//     counter1 int64 // 可能在同一缓存行
//     counter2 int64 // 导致缓存行冲突
// }

// // 解决：使用填充分离
// type Solution struct {
//     counter1 int64
//     _        [56]byte // 填充，确保counter2在下一个缓存行
//     counter2 int64
// }
// 原子操作是Go并发编程中最基础也是最高效的同步机制，掌握好原子操作对于编写高性能的并发程序非常重要。你想深入了解哪个特定的原子操作模式？

// 高性能计数器

// 64+56+8 = 128 每个实例独占一行
// golang 头部 8字节
// [56]byte 填充的目的是：
// 避免伪共享：确保每个计数器实例独占缓存行
// 考虑Go运行时开销：补偿对象头部的8字节占用
// 内存对齐：满足Go语言的内存对齐要求
// 总大小128字节：是缓存行大小(64字节)的整数倍
type HighPerformanceCounter struct {
	counters [8]int64
	_        [56]byte
}

func (hpc *HighPerformanceCounter) Increment(goroutineID int) {
	index := goroutineID % 8
	atomic.AddInt64(&hpc.counters[index], 1)
}
