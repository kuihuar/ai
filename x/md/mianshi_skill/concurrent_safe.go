package mianshiskill

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	//mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.value++
}
func (c *Counter) Value() int {
	return c.value
}

type CounterWithMutex struct {
	mu    sync.Mutex
	value int
}

func (c *CounterWithMutex) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// 读多写少时可以改成读锁
func (c *CounterWithMutex) Value() int {
	return c.value
}

type CounterWithRWMutex struct {
	mu    sync.RWMutex
	value int
}

// 多个读操作并行
func (c *CounterWithRWMutex) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// 性能对比
// 方案	吞吐量（ops/ns）	适用场景
// 互斥锁	中	复杂操作或需强一致性
// 原子操作	高	简单计数器、高性能场景
// Channel	低	需要与其他异步逻辑配合
// 读写锁	高	读多写少的场景，如缓存

type CounterWithAtomic struct {
	value int32
}

func (c *CounterWithAtomic) Increment() {
	atomic.AddInt32(&c.value, 1)
}
func (c *CounterWithAtomic) Value() int {
	res := atomic.LoadInt32(&c.value)
	return int(res)
}

type Integer interface {
	~int32 | ~int64
}

type CounerWithGeneric[T Integer] struct {
	value T
}

func (c *CounerWithGeneric[T]) Increment() {

	switch ptr := any(&c.value).(type) {
	case *int32:
		atomic.AddInt32(ptr, 1)

	case *int64:
		atomic.AddInt64(ptr, 1)

	}
}

func (c *CounerWithGeneric[T]) Value() T {
	switch ptr := any(&c.value).(type) {
	case *int32:
		return T(atomic.LoadInt32(ptr))
	case *int64:
		return T(atomic.LoadInt64(ptr))
	default:
		return 0
	}
}

func UseCounter() {
	c32 := CounerWithGeneric[int32]{value: 0}
	c32.Increment()
	fmt.Printf("c32: %v \n", c32.Value())
	c32.Increment()
	fmt.Printf("c32: %v \n", c32.Value())

}
