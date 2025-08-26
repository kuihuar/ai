package skill

import (
	"sync"
	"sync/atomic"
)

// ConcurrentMap 线程安全的 map
type ConcurrentMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewConcurrentMap 创建新的并发 map
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data: make(map[K]V),
	}
}

// Set 设置键值对
func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = value
}

// Get 获取值
func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, exists := cm.data[key]
	return value, exists
}

// Delete 删除键
func (cm *ConcurrentMap[K, V]) Delete(key K) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.data, key)
}

// Len 获取长度
func (cm *ConcurrentMap[K, V]) Len() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.data)
}

// Keys 获取所有键
func (cm *ConcurrentMap[K, V]) Keys() []K {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	keys := make([]K, 0, len(cm.data))
	for key := range cm.data {
		keys = append(keys, key)
	}
	return keys
}

// Values 获取所有值
func (cm *ConcurrentMap[K, V]) Values() []V {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	values := make([]V, 0, len(cm.data))
	for _, value := range cm.data {
		values = append(values, value)
	}
	return values
}

// ConcurrentSlice 线程安全的 slice
type ConcurrentSlice[T any] struct {
	mu    sync.RWMutex
	items []T
}

// NewConcurrentSlice 创建新的并发 slice
func NewConcurrentSlice[T any]() *ConcurrentSlice[T] {
	return &ConcurrentSlice[T]{
		items: make([]T, 0),
	}
}

// Append 添加元素
func (cs *ConcurrentSlice[T]) Append(item T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.items = append(cs.items, item)
}

// Get 获取指定索引的元素
func (cs *ConcurrentSlice[T]) Get(index int) (T, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	if index < 0 || index >= len(cs.items) {
		var zero T
		return zero, false
	}
	return cs.items[index], true
}

// Set 设置指定索引的元素
func (cs *ConcurrentSlice[T]) Set(index int, item T) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if index < 0 || index >= len(cs.items) {
		return false
	}
	cs.items[index] = item
	return true
}

// Len 获取长度
func (cs *ConcurrentSlice[T]) Len() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.items)
}

// Remove 删除指定索引的元素
func (cs *ConcurrentSlice[T]) Remove(index int) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if index < 0 || index >= len(cs.items) {
		return false
	}
	cs.items = append(cs.items[:index], cs.items[index+1:]...)
	return true
}

// ConcurrentQueue 线程安全的队列
type ConcurrentQueue[T any] struct {
	mu    sync.Mutex
	items []T
}

// NewConcurrentQueue 创建新的并发队列
func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		items: make([]T, 0),
	}
}

// Enqueue 入队
func (cq *ConcurrentQueue[T]) Enqueue(item T) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.items = append(cq.items, item)
}

// Dequeue 出队
func (cq *ConcurrentQueue[T]) Dequeue() (T, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if len(cq.items) == 0 {
		var zero T
		return zero, false
	}
	item := cq.items[0]
	cq.items = cq.items[1:]
	return item, true
}

// Peek 查看队首元素
func (cq *ConcurrentQueue[T]) Peek() (T, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if len(cq.items) == 0 {
		var zero T
		return zero, false
	}
	return cq.items[0], true
}

// Len 获取队列长度
func (cq *ConcurrentQueue[T]) Len() int {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return len(cq.items)
}

// IsEmpty 检查队列是否为空
func (cq *ConcurrentQueue[T]) IsEmpty() bool {
	return cq.Len() == 0
}

// ConcurrentCounter 线程安全的计数器
type ConcurrentCounter struct {
	value int64
}

// NewConcurrentCounter 创建新的并发计数器
func NewConcurrentCounter() *ConcurrentCounter {
	return &ConcurrentCounter{}
}

// Increment 增加计数
func (cc *ConcurrentCounter) Increment() {
	atomic.AddInt64(&cc.value, 1)
}

// Decrement 减少计数
func (cc *ConcurrentCounter) Decrement() {
	atomic.AddInt64(&cc.value, -1)
}

// Get 获取当前值
func (cc *ConcurrentCounter) Get() int64 {
	return atomic.LoadInt64(&cc.value)
}

// Set 设置值
func (cc *ConcurrentCounter) Set(value int64) {
	atomic.StoreInt64(&cc.value, value)
}

// Add 增加指定值
func (cc *ConcurrentCounter) Add(delta int64) {
	atomic.AddInt64(&cc.value, delta)
}

// ConcurrentSet 线程安全的集合
type ConcurrentSet[T comparable] struct {
	mu   sync.RWMutex
	data map[T]struct{}
}

// NewConcurrentSet 创建新的并发集合
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{
		data: make(map[T]struct{}),
	}
}

// Add 添加元素
func (cs *ConcurrentSet[T]) Add(item T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.data[item] = struct{}{}
}

// Remove 删除元素
func (cs *ConcurrentSet[T]) Remove(item T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.data, item)
}

// Contains 检查是否包含元素
func (cs *ConcurrentSet[T]) Contains(item T) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	_, exists := cs.data[item]
	return exists
}

// Len 获取集合大小
func (cs *ConcurrentSet[T]) Len() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.data)
}

// Items 获取所有元素
func (cs *ConcurrentSet[T]) Items() []T {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	items := make([]T, 0, len(cs.data))
	for item := range cs.data {
		items = append(items, item)
	}
	return items
}

// Clear 清空集合
func (cs *ConcurrentSet[T]) Clear() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.data = make(map[T]struct{})
}
