package skill

import (
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== 1. 基础锁类型 ====================

// 1.1 互斥锁 (Mutex) - 最基本的锁

type MutexExample struct {
	mu    sync.Mutex
	count int
}

func (m *MutexExample) Increment() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.count++
}
func (m *MutexExample) GetCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.count
}

// 1.2 读写锁 (RWMutex) - 允许多个读操作，但只允许一个写操作

type RWMutexExample struct {
	mu   sync.RWMutex
	data map[string]string
}

func (r *RWMutexExample) Read(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, exists := r.data[key]
	return value, exists
}
func (r *RWMutexExample) Write(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

// ==================== 2. 高级锁模式 ====================

type FairLock struct {
	mu      sync.Mutex
	waiting []chan struct{}
	locked  bool
}

func NewFairLock() *FairLock {
	return &FairLock{
		waiting: make([]chan struct{}, 0),
	}
}

func (f *FairLock) Lock() {
	f.mu.Lock()
	if !f.locked {
		f.locked = true
		f.mu.Unlock()
		return
	}
	wait := make(chan struct{})
	f.waiting = append(f.waiting, wait)
	f.mu.Unlock()
	<-wait
}

func (f *FairLock) Unlock() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.waiting) > 0 {
		wait := f.waiting[0]
		f.waiting = f.waiting[1:]
		close(wait)
	} else {
		f.locked = false
	}
}

// 2.2 自旋锁 (SpinLock) - 避免线程切换开销
// 适用于短时间锁定，避免频繁的上下文切换

type SpinLock struct {
	locked int32
}

func (s *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&s.locked, 0, 1) {
		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32(&s.locked, 0)
}

// 2.3 可重入锁 - 基于原子计数实现
// 可重入 的核心价值是: 允许同一持有者（Goroutine）在持有锁的情况下，再次获取同一把锁而不阻塞，解决了嵌套调用、递归调用时的死锁问题。

// 普通锁（不可重入）做不到这一点，同一持有者再次获取时会被自己阻塞，导致死锁

type ReentrantLock struct {
	mu    sync.Mutex
	owner int
	count int
}

func (r *ReentrantLock) Lock() {
	goroutineID := runtime.NumGoroutine()
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.owner == goroutineID {
		r.count++
	} else {
		for r.count > 0 {
			r.mu.Unlock()
			time.Sleep(time.Microsecond)
			r.mu.Lock()
		}
		r.owner = goroutineID
		r.count = 1
	}

}

func (r *ReentrantLock) Unlock() {
	goroutineID := runtime.NumGoroutine()
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.owner == goroutineID {
		r.count--
		if r.count == 0 {
			r.owner = 0
		}
	}
}

// 2.4 升降级锁 - 允许从写锁降级为读锁
// 升降级锁的作用：解决 “先读后写” 的原子性问题。普通读写锁做不到这一点，因为读写锁是互斥的
// 升级过程是原子的，避免了 “释放读锁后、获取写锁前” 的间隙被其他线程干扰，保证了数据一致性

type UpgradeLock struct {
	mu        sync.Mutex
	readers   int
	writers   int
	upgrading int
	owner     int
}

func (u *UpgradeLock) Lock() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.writers > 0 || u.upgrading > 0 {
		u.mu.Unlock()
		time.Sleep(time.Microsecond)
		u.mu.Lock()
	}
	u.readers++
}

func (u *UpgradeLock) Unlock() {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.readers > 0 {
		u.readers--
	}
}

func (u *UpgradeLock) Upgrade() bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.readers == 1 && u.writers == 0 && u.upgrading == 0 {
		u.upgrading = 1
		u.readers = 0
		u.writers = 1
		u.owner = runtime.NumGoroutine()
		return true
	}
	return false
}

func (u *UpgradeLock) Downgrade() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.writers > 0 && u.owner == runtime.NumGoroutine() {
		u.upgrading = 0
		u.readers = 1
		u.upgrading = 0
		u.owner = 0
	}
}

// 推荐：使用RWMutex
type Cache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// 推荐：使用原子操作
type Counter struct {
	count int64
}

// 推荐：使用channel
type Pool struct {
	conns chan net.Conn
}
