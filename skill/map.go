package skill

import (
	"fmt"
	"sync"
	"time"
)

// ==================== sync.Map 同步原语 ====================

// 1. 基本用法示例

// 1.1 简单的并发安全映射
type SimpleConcurrentMap struct {
	data sync.Map
}

func NewSimpleConcurrentMap() *SimpleConcurrentMap {
	return &SimpleConcurrentMap{}
}

func (scm *SimpleConcurrentMap) Set(key, value interface{}) {
	scm.data.Store(key, value)
}

func (scm *SimpleConcurrentMap) Get(key interface{}) (interface{}, bool) {
	return scm.data.Load(key)
}

func (scm *SimpleConcurrentMap) Delete(key interface{}) {
	scm.data.Delete(key)
}

func (scm *SimpleConcurrentMap) Range(f func(key, value interface{}) bool) {
	scm.data.Range(f)
}

// 1.2 类型安全的映射
type TypedConcurrentMap struct {
	data sync.Map
}

func NewTypedConcurrentMap() *TypedConcurrentMap {
	return &TypedConcurrentMap{}
}

func (tcm *TypedConcurrentMap) SetString(key string, value string) {
	tcm.data.Store(key, value)
}

func (tcm *TypedConcurrentMap) GetString(key string) (string, bool) {
	if value, ok := tcm.data.Load(key); ok {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

func (tcm *TypedConcurrentMap) SetInt(key string, value int) {
	tcm.data.Store(key, value)
}

func (tcm *TypedConcurrentMap) GetInt(key string) (int, bool) {
	if value, ok := tcm.data.Load(key); ok {
		if num, ok := value.(int); ok {
			return num, true
		}
	}
	return 0, false
}

// 2. 高级用法

// 2.1 计数器映射
type CounterMap struct {
	counters sync.Map
}

func NewCounterMap() *CounterMap {
	return &CounterMap{}
}

func (cm *CounterMap) Increment(key string) int {
	for {
		if value, loaded := cm.counters.Load(key); loaded {
			if count, ok := value.(int); ok {
				if cm.counters.CompareAndSwap(key, count, count+1) {
					return count + 1
				}
			}
		} else {
			if cm.counters.CompareAndSwap(key, nil, 1) {
				return 1
			}
		}
	}
}

func (cm *CounterMap) GetCount(key string) int {
	if value, ok := cm.counters.Load(key); ok {
		if count, ok := value.(int); ok {
			return count
		}
	}
	return 0
}

// 2.2 缓存映射
type CacheMap struct {
	cache sync.Map
}

type CacheEntry struct {
	Value      interface{}
	ExpireTime time.Time
}

func NewCacheMap() *CacheMap {
	cm := &CacheMap{}

	// 启动清理协程
	go cm.cleanupExpired()

	return cm
}

func (cm *CacheMap) Set(key string, value interface{}, ttl time.Duration) {
	entry := &CacheEntry{
		Value:      value,
		ExpireTime: time.Now().Add(ttl),
	}
	cm.cache.Store(key, entry)
}

func (cm *CacheMap) Get(key string) (interface{}, bool) {
	if value, ok := cm.cache.Load(key); ok {
		if entry, ok := value.(*CacheEntry); ok {
			if time.Now().Before(entry.ExpireTime) {
				return entry.Value, true
			} else {
				cm.cache.Delete(key)
			}
		}
	}
	return nil, false
}

func (cm *CacheMap) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cm.cache.Range(func(key, value interface{}) bool {
			if entry, ok := value.(*CacheEntry); ok {
				if time.Now().After(entry.ExpireTime) {
					cm.cache.Delete(key)
				}
			}
			return true
		})
	}
}

// 2.3 读写分离映射
type ReadWriteMap struct {
	readMap  sync.Map
	writeMap sync.Map
	mu       sync.RWMutex
}

func NewReadWriteMap() *ReadWriteMap {
	return &ReadWriteMap{}
}

func (rwm *ReadWriteMap) Set(key, value interface{}) {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	rwm.writeMap.Store(key, value)
}

func (rwm *ReadWriteMap) Get(key interface{}) (interface{}, bool) {
	rwm.mu.RLock()
	defer rwm.mu.RUnlock()

	// 先从写映射读取
	if value, ok := rwm.writeMap.Load(key); ok {
		return value, true
	}

	// 再从读映射读取
	return rwm.readMap.Load(key)
}

func (rwm *ReadWriteMap) Flush() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	// 将写映射的内容合并到读映射
	rwm.writeMap.Range(func(key, value interface{}) bool {
		rwm.readMap.Store(key, value)
		return true
	})

	// 清空写映射
	rwm.writeMap = sync.Map{}
}

// 3. 实际应用场景

// 3.1 配置管理器
type ConfigManager struct {
	config sync.Map
}

func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{}

	// 加载默认配置
	cm.config.Store("timeout", 30)
	cm.config.Store("max_connections", 100)
	cm.config.Store("debug", false)

	return cm
}

func (cm *ConfigManager) SetConfig(key string, value interface{}) {
	cm.config.Store(key, value)
}

func (cm *ConfigManager) GetConfig(key string) (interface{}, bool) {
	return cm.config.Load(key)
}

func (cm *ConfigManager) GetConfigWithDefault(key string, defaultValue interface{}) interface{} {
	if value, ok := cm.config.Load(key); ok {
		return value
	}
	return defaultValue
}

// 3.2 会话管理器
type SessionManager struct {
	sessions sync.Map
}

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	LastSeen  time.Time
	Data      map[string]interface{}
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{}

	// 启动会话清理协程
	go sm.cleanupExpiredSessions()

	return sm
}

func (sm *SessionManager) CreateSession(userID string) *Session {
	session := &Session{
		ID:        fmt.Sprintf("session-%d", time.Now().UnixNano()),
		UserID:    userID,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
		Data:      make(map[string]interface{}),
	}

	sm.sessions.Store(session.ID, session)
	return session
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	if value, ok := sm.sessions.Load(sessionID); ok {
		if session, ok := value.(*Session); ok {
			session.LastSeen = time.Now()
			return session, true
		}
	}
	return nil, false
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.sessions.Delete(sessionID)
}

func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		expireTime := time.Now().Add(-time.Hour * 24) // 24小时过期

		sm.sessions.Range(func(key, value interface{}) bool {
			if session, ok := value.(*Session); ok {
				if session.LastSeen.Before(expireTime) {
					sm.sessions.Delete(key)
				}
			}
			return true
		})
	}
}

// 4. 性能优化技巧

// 4.1 批量操作
type BatchMap struct {
	data sync.Map
}

func (bm *BatchMap) BatchSet(items map[interface{}]interface{}) {
	for key, value := range items {
		bm.data.Store(key, value)
	}
}

func (bm *BatchMap) BatchGet(keys []interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})

	for _, key := range keys {
		if value, ok := bm.data.Load(key); ok {
			result[key] = value
		}
	}

	return result
}

// 4.2 条件更新
func (bm *BatchMap) UpdateIf(key interface{}, condition func(interface{}) bool, newValue interface{}) bool {
	for {
		if oldValue, loaded := bm.data.Load(key); loaded {
			if condition(oldValue) {
				if bm.data.CompareAndSwap(key, oldValue, newValue) {
					return true
				}
			} else {
				return false
			}
		} else {
			if bm.data.CompareAndSwap(key, nil, newValue) {
				return true
			}
		}
	}
}

// 5. 注意事项和最佳实践

// 5.1 避免存储大对象
type SafeMap struct {
	data sync.Map
}

func (sm *SafeMap) SafeStore(key, value interface{}) {
	// 检查值的大小
	if size := estimateSize(value); size > 1024*1024 { // 1MB
		fmt.Printf("Warning: storing large object of size %d bytes\n", size)
	}
	sm.data.Store(key, value)
}

func estimateSize(value interface{}) int {
	// 简单的大小估算
	switch v := value.(type) {
	case string:
		return len(v)
	case []byte:
		return len(v)
	case map[string]interface{}:
		return len(v) * 100 // 粗略估算
	default:
		return 100 // 默认估算
	}
}

// 5.2 正确处理类型断言
func (sm *SafeMap) SafeLoad(key interface{}) (interface{}, bool) {
	value, ok := sm.data.Load(key)
	if !ok {
		return nil, false
	}
	return value, true
}

// 5.3 内存泄漏防护
type LeakProtectedMap struct {
	data sync.Map
}

func (lpm *LeakProtectedMap) StoreWithCleanup(key, value interface{}, cleanup func()) {
	lpm.data.Store(key, value)

	// 设置清理函数
	go func() {
		time.Sleep(time.Hour) // 1小时后清理
		if _, loaded := lpm.data.LoadAndDelete(key); loaded {
			cleanup()
		}
	}()
}
