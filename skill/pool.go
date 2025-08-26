package skill

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type SimplePool struct {
	pool sync.Pool
}

func NewSimplePool() *SimplePool {
	return &SimplePool{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (p *SimplePool) GetBuffer() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func (p *SimplePool) PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}

// 1.2 连接池
type ConnectionPoolWithPool struct {
	pool sync.Pool
}

type PoolConnection struct {
	ID        string
	CreatedAt time.Time
	InUse     bool
}

func NewConnectionPoolWithPool() *ConnectionPoolWithPool {
	return &ConnectionPoolWithPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &PoolConnection{
					ID:        fmt.Sprintf("conn-%d", time.Now().UnixNano()),
					CreatedAt: time.Now(),
					InUse:     false,
				}
			},
		},
	}
}

func (cp *ConnectionPoolWithPool) GetConnection() *PoolConnection {
	conn := cp.pool.Get().(*PoolConnection)
	conn.InUse = true
	return conn
}

func (cp *ConnectionPoolWithPool) PutConnection(conn *PoolConnection) {
	conn.InUse = false
	cp.pool.Put(conn)
}

// 2. 高级用法

// 2.1 带清理的对象池
type CleanablePool struct {
	pool sync.Pool
}

type CleanableObject struct {
	Data      []byte
	Timestamp time.Time
	Clean     bool
}

func NewCleanablePool() *CleanablePool {
	return &CleanablePool{
		pool: sync.Pool{
			New: func() interface{} {
				return &CleanableObject{
					Data:      make([]byte, 0, 1024),
					Timestamp: time.Now(),
					Clean:     true,
				}
			},
		},
	}
}

func (cp *CleanablePool) GetObject() *CleanableObject {
	obj := cp.pool.Get().(*CleanableObject)

	// 检查对象是否需要清理
	if time.Since(obj.Timestamp) > time.Minute*5 {
		obj.Data = obj.Data[:0] // 清空数据
		obj.Timestamp = time.Now()
		obj.Clean = true
	}

	return obj
}

func (cp *CleanablePool) PutObject(obj *CleanableObject) {
	if len(obj.Data) > 1024*1024 { // 如果数据太大，不回收
		return
	}
	cp.pool.Put(obj)
}

// 2.2 多级对象池
type MultiLevelPool struct {
	smallPool  sync.Pool
	mediumPool sync.Pool
	largePool  sync.Pool
}

type SizedObject struct {
	Data []byte
	Size int
}

func NewMultiLevelPool() *MultiLevelPool {
	return &MultiLevelPool{
		smallPool: sync.Pool{
			New: func() interface{} {
				return &SizedObject{
					Data: make([]byte, 0, 1024),
					Size: 1024,
				}
			},
		},
		mediumPool: sync.Pool{
			New: func() interface{} {
				return &SizedObject{
					Data: make([]byte, 0, 1024*1024),
					Size: 1024 * 1024,
				}
			},
		},
		largePool: sync.Pool{
			New: func() interface{} {
				return &SizedObject{
					Data: make([]byte, 0, 10*1024*1024),
					Size: 10 * 1024 * 1024,
				}
			},
		},
	}
}

func (mlp *MultiLevelPool) GetObject(size int) *SizedObject {
	switch {
	case size <= 1024:
		return mlp.smallPool.Get().(*SizedObject)
	case size <= 1024*1024:
		return mlp.mediumPool.Get().(*SizedObject)
	default:
		return mlp.largePool.Get().(*SizedObject)
	}
}

func (mlp *MultiLevelPool) PutObject(obj *SizedObject) {
	obj.Data = obj.Data[:0] // 重置数据

	switch obj.Size {
	case 1024:
		mlp.smallPool.Put(obj)
	case 1024 * 1024:
		mlp.mediumPool.Put(obj)
	case 10 * 1024 * 1024:
		mlp.largePool.Put(obj)
	}
}

// 3. 实际应用场景

// 3.1 JSON 编码器池
// 在 Go 中，使用 encoding/json 包进行 JSON 编码时，通常需要创建 json.Encoder 或临时的 bytes.Buffer（用于存储编码后的字节数据）
// 问题：如果在高并发场景（如 Web 服务器处理大量请求）中频繁调用这类函数，会导致大量临时对象（bytes.Buffer、json.Encoder）被创建和销毁。
// 后果：这些对象会触发频繁的垃圾回收（GC），增加系统开销，降低程序性能（尤其是在每秒处理数万次请求的场景中）
// JSONEncoderPool 的核心价值：复用对象，减少开销
// JSONEncoderPool 通过 sync.Pool 缓存 JSONEncoder 实例，实现 “创建一次，多次复用”

type JSONEncoderPool struct {
	pool sync.Pool
}

type JSONEncoder struct {
	buffer  *bytes.Buffer
	encoder interface{} // 实际的 JSON 编码器
}

func NewJSONEncoderPool() *JSONEncoderPool {
	return &JSONEncoderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &JSONEncoder{
					buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
				}
			},
		},
	}
}

func (jep *JSONEncoderPool) GetEncoder() *JSONEncoder {
	encoder := jep.pool.Get().(*JSONEncoder)
	encoder.buffer.Reset()
	return encoder
}

func (jep *JSONEncoderPool) PutEncoder(encoder *JSONEncoder) {
	jep.pool.Put(encoder)
}

type TempBufferPool struct {
	pool sync.Pool
}

type TempBuffer struct {
	Data   []byte
	Length int
}

func NewTempBufferPool() *TempBufferPool {
	return &TempBufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &TempBuffer{
					Data:   make([]byte, 4096),
					Length: 0,
				}
			},
		},
	}
}

func (tbp *TempBufferPool) GetBuffer() *TempBuffer {
	buffer := tbp.pool.Get().(*TempBuffer)
	buffer.Length = 0
	return buffer
}
func (tbp *TempBufferPool) PutBuffer(buffer *TempBuffer) {
	if buffer.Length > len(buffer.Data)/2 {
		// 如果使用了超过一半的缓冲区，不回收
		return
	}
	tbp.pool.Put(buffer)
}

// 4. 性能优化技巧

// 4.1 预热池
func (sp *SimplePool) WarmUp(count int) {
	for i := 0; i < count; i++ {
		buf := sp.pool.Get().(*bytes.Buffer)
		sp.pool.Put(buf)
	}
}

// 4.2 批量操作
type BatchProcessor struct {
	pool sync.Pool
}

func (bp *BatchProcessor) ProcessBatch(items []string) []string {
	results := make([]string, 0, len(items))

	for _, item := range items {
		buffer := bp.pool.Get().(*bytes.Buffer)
		buffer.Reset()

		// 处理项目
		buffer.WriteString("processed: ")
		buffer.WriteString(item)

		results = append(results, buffer.String())
		bp.pool.Put(buffer)
	}

	return results
}

// 5. 注意事项和最佳实践

// 5.1 不要存储大对象
type BadPoolExample struct {
	pool sync.Pool
}

func (bpe *BadPoolExample) BadPut(obj *SizedObject) {
	if len(obj.Data) > 1024*1024 { // 1MB
		// 大对象不放入池中
		return
	}
	bpe.pool.Put(obj)
}

// 5.2 正确处理 nil 值
type SafePool struct {
	pool sync.Pool
}

func (sp *SafePool) SafeGet() *bytes.Buffer {
	obj := sp.pool.Get()
	if obj == nil {
		return &bytes.Buffer{}
	}
	return obj.(*bytes.Buffer)
}

// 5.3 池的清理策略
type CleanupPool struct {
	pool      sync.Pool
	lastClean time.Time
	mu        sync.Mutex
}

func (cp *CleanupPool) GetWithCleanup() interface{} {
	cp.mu.Lock()
	if time.Since(cp.lastClean) > time.Minute*10 {
		// 每10分钟清理一次池
		cp.pool = sync.Pool{
			New: cp.pool.New,
		}
		cp.lastClean = time.Now()
	}
	cp.mu.Unlock()

	return cp.pool.Get()
}
