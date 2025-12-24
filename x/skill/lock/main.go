package main

import (
	"sync"
	"sync/atomic"
)

type MutexCache struct {
	mu   sync.Mutex
	data map[string]interface{}
}

func NewMutexCache() *MutexCache {
	return &MutexCache{data: make(map[string]interface{})}
}

func (c *MutexCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
func (c *MutexCache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[key]

}

type RWMutexCache struct {
	rwmu sync.RWMutex
	data map[string]interface{}
}

func NewRWMutexCache() *RWMutexCache {
	return &RWMutexCache{data: make(map[string]interface{})}
}

func (c *RWMutexCache) Set(key string, value interface{}) {
	c.rwmu.Lock()
	defer c.rwmu.Unlock()
	c.data[key] = value
}

func (c *RWMutexCache) Get(key string) interface{} {
	c.rwmu.RLock()
	defer c.rwmu.RUnlock()
	return c.data[key]

}

// 分段锁
type ConcurrentMap struct {
	segments    []*segment
	numSegments int
}
type segment struct {
	lock sync.RWMutex
	data map[string]interface{}
}

func NewConcurrentMap(numSegments int) {
	cm := &ConcurrentMap{
		segments:    make([]*segment, numSegments),
		numSegments: numSegments,
	}
	for i := range cm.segments {
		cm.segments[i] = &segment{
			data: make(map[string]interface{}),
		}
	}
}

func (cm *ConcurrentMap) getSegmentIndex(key string) int {
	hash := 0
	for _, char := range key {
		hash += int(char)
	}
	return hash % cm.numSegments
}

func (cm *ConcurrentMap) Set(key string, value interface{}) {
	segmentIndex := cm.getSegmentIndex(key)
	cm.segments[segmentIndex].lock.RLock()
	defer cm.segments[segmentIndex].lock.Unlock()
	cm.segments[segmentIndex].data[key] = value
}
func (cm *ConcurrentMap) Get(key string) interface{} {
	segmentIndex := cm.getSegmentIndex(key)
	cm.segments[segmentIndex].lock.RLock()
	defer cm.segments[segmentIndex].lock.RUnlock()

	return cm.segments[segmentIndex].data[key]
}

type LockFreeCache struct {
	cacheMap atomic.Value
}

func NewLockFreeCache() *LockFreeCache {
	c := &LockFreeCache{}
	initialMap := make(map[string]interface{})
	c.cacheMap.Store(initialMap)
	return c
}

func (c *LockFreeCache) Update(newMap map[string]interface{}) {
	newMapPtr := &newMap
	c.cacheMap.Store(newMapPtr)
}
func (c *LockFreeCache) Get(key string) (interface{}, bool) {
	cacheMap := c.cacheMap.Load().(*map[string]interface{})
	value, ok := (*cacheMap)[key]
	return value, ok
}

type CacheMap struct {
	cacheMap sync.Map
}

func NewCacheMap() *CacheMap {
	return &CacheMap{}
}

func (c *CacheMap) Get(key string) (interface{}, bool) {
	return c.cacheMap.Load(key)
}

func (c *CacheMap) Update(newMap map[string]interface{}) {
	c.cacheMap.Range(func(key, value interface{}) bool {
		c.cacheMap.Delete(key)
		return true
	})

	for k, v := range newMap {
		c.cacheMap.Store(k, v)
	}
}

type Node struct {
	value interface{}
	next  *Node
}

// 只能操作栈项，不能操作栈底，后进先出
// 任何不在栈顶的数据都无法访问
// 只能在一端操作的线性数据表
// 具有记忆作用
// 占据固定大小空间
// 函数调用， 表达式求值， 括号匹配
// 表达式求值用两个栈，一个保存操作数，下保存运算符
type LockFreeStack struct {
	top atomic.Pointer[Node]
}

func NewLockFreeStack() *LockFreeStack {
	return &LockFreeStack{}
}

func (s *LockFreeStack) Push(value interface{}) {
	node := &Node{value: value}

	for {
		oldTop := s.top.Load()
		node.next = oldTop
		if s.top.CompareAndSwap(oldTop, node) {
			return
		}
	}
}

func (s *LockFreeStack) Pop() (interface{}, bool) {
	for {
		oldTop := s.top.Load()
		if oldTop == nil {
			return nil, false
		}
		newTop := oldTop.next
		if s.top.CompareAndSwap(oldTop, newTop) {
			return oldTop.value, true
		}
	}
}
func main() {}
