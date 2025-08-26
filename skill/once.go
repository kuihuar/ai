package skill

import (
	"fmt"
	"sync"
	"time"
)

// 基础用法
// 1. 单例

type Singleton struct {
	data string
	once sync.Once
}

func (s *Singleton) Initialize() {
	s.once.Do(
		func() {
			s.data = "initialize data"
		})
}

func (s *Singleton) GetData() string {
	return s.data
}

// 2. 延迟初始化

type lazyInitializer struct {
	config map[string]interface{}
	once   sync.Once
}

func (li *lazyInitializer) GetConfig() map[string]interface{} {
	li.once.Do(func() {
		li.config = loadConfiguration()
	})
	return li.config
}
func loadConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
}

// 1.3 资源初始化
type ResourceManagerWithOnce struct {
	connection *ConnectionWithOnce
	once       sync.Once
}

type ConnectionWithOnce struct {
	ID   int
	Open bool
}

func (rm *ResourceManagerWithOnce) GetConnection() *ConnectionWithOnce {
	rm.once.Do(func() {
		rm.connection = &ConnectionWithOnce{
			ID:   int(time.Now().UnixNano()),
			Open: true,
		}
		fmt.Printf("Connection established: %s\n", rm.connection.ID)
	})
	return rm.connection
}

// 2.1 错误处理版本

// 无论多少个 Goroutine 同时调用 Initialize()，once.Do 内部的初始化函数（performInitialization）都只会执行一次。
// 后续的调用会直接返回第一次初始化的结果（data 或 err）。

type SafeInitializer struct {
	data interface{}
	err  error
	once sync.Once
}

func (si *SafeInitializer) Initialize() (interface{}, error) {
	si.once.Do(func() {
		si.data, si.err = performInitialization()
	})
	return si.data, si.err
}

func performInitialization() (interface{}, error) {
	// 模拟可能失败的初始化
	if time.Now().UnixNano()%2 == 0 {
		return "success", nil
	}
	return nil, fmt.Errorf("initialization failed")
}

// 2.2 带参数的初始化
type ParameterizedInitializer struct {
	initializers map[string]interface{}
	mu           sync.RWMutex
}

func NewParameterizedInitializer() *ParameterizedInitializer {
	return &ParameterizedInitializer{
		initializers: make(map[string]interface{}),
	}
}

// 带参数的初始化（ParameterizedInitializer）：
// 多 key 场景的初始化（按需初始化）

func (pi *ParameterizedInitializer) GetOrInitialize(key string, initFunc func() interface{}) interface{} {
	pi.mu.RLock()
	if value, exists := pi.initializers[key]; exists {
		pi.mu.RUnlock()
		return value
	}
	pi.mu.RUnlock()

	pi.mu.Lock()
	defer pi.mu.Unlock()

	// 双重检查
	if value, exists := pi.initializers[key]; exists {
		return value
	}

	value := initFunc()
	pi.initializers[key] = value
	return value
}

// 3. 实际应用场景

// 3.1 数据库连接池初始化
type DatabasePool struct {
	connections []*ConnectionWithOnce
	once        sync.Once
}

func (dp *DatabasePool) InitializePool() {
	dp.once.Do(func() {
		dp.connections = make([]*ConnectionWithOnce, 10)
		for i := 0; i < 10; i++ {
			dp.connections[i] = &ConnectionWithOnce{
				ID:   int(time.Now().UnixNano()),
				Open: true,
			}
		}
		fmt.Println("Database pool initialized with 10 connections")
	})
}

// 3.2 日志系统初始化
type Logger struct {
	level string
	once  sync.Once
}

func (l *Logger) SetLevel(level string) {
	l.once.Do(func() {
		l.level = level
		fmt.Printf("Logger level set to: %s\n", level)
	})
}

// 4. 性能测试和对比

func BenchmarkOnceVsMutex() {
	// sync.Once 版本
	var once sync.Once
	// 核心逻辑：sync.Once 是 Go 标准库专门设计的同步原语，其 Do 方法接收的函数无论被多少个 Goroutine 调用，都只会被执行一次
	// 底层保证：sync.Once 内部通过原子操作（atomic）和互斥锁（mutex）的组合实现
	once.Do(func() {
		time.Sleep(time.Millisecond * 10)
	})

	// 互斥锁版本
	// 虽正确但有锁竞争开销
	var mu sync.Mutex
	var initialized bool
	mu.Lock()
	if !initialized {
		time.Sleep(time.Millisecond * 10)
		initialized = true
	}
	mu.Unlock()
}
