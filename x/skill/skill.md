
## 单机吞吐优化

### 数据处理

1. 切片尽量指定容量，避免扩容迁移和频繁扩容
2. 字典集合使用空struct
3. 字典尽量指定容量，避免扩容迁移和频繁扩容
4. 字符串切片相互转换可以考虑切unsafe包（零内存拷贝）
    ```go
    func Str2BytesZeroCopy(s string) []byte {
        return unsafe.Slice(unsafe.StringData(s), len(s))
    }
    func Str2BytesSafe(s string) []byte {
        return []byte(s)
    }

    func Bytes2StrZeroCopy(b []byte) string {
        return *(*string)(unsafe.Pointer(&b))
    }
    func Bytes2StrSafew(b []byte) string {
        return string(b)
    }
    ```
5. 字符串转数字，整型转字字串由fmt改为strconv
6. 字符串接拼，字会串拼接改为builder

### 资源复用
1. 使用对象池，避免频繁分配相同类型的临时对象开销（内存分配和临时对象垃圾回收）
2. 协程池，避免创建协程创建（内存分配）和协程调度（CPU开销）
  - CPU密集，一般设置为cpu核心数的倍数
  - IO密集，一般要设置大很多
  - 混合型则介于CPU密和集IO密集之间，
  - 设置后，要通过压测进一步优化，直到系统吞吐量无法显著提升。


### 测试命令帮助
    ```shell
    go help testflag

    go tool pprof -help

    ```



## 并发等待

### waitgroup

1. 阻塞等待多个并发任务执行完成
### errorgroup


## 并发锁


### 并发访问
- 写操作多，读不频繁，使用互斥锁保证并发访问安全
- 读操作远远大于写操作，使用读写锁提高并发读取的性能
- map 分布均匀，使用分段锁，降低锁粒度
- 需要共享对像进行原子操作时，可以使用atomic无锁编程。
- 无锁性能高于锁

### sync.Map
- 读多写少的使用模式
- 多个协程频繁读写一个map
### 实现思想
- 读写分离，无锁读取，延迟删除
- 频繁写入或者键冲突多时，性能可能劣于map+RWMutex
```go
type Map struct {
        _ noCopy

        mu Mutex
        // read 包含 map 内容中可以安全并发访问的部分（无论是否持有 mu）。
        //
        // read 字段本身始终可以安全加载，但必须仅使用 mu 进行存储。
        //
        // read 中存储的条目可以在没有 mu 的情况下并发更新，但更新先前已删除的条目需要将条目复制到 dirty
        // map 中并在持有 mu 的情况下取消删除。
        // read contains the portion of the map's contents that are safe for
        // concurrent access (with or without mu held).
        //
        // The read field itself is always safe to load, but must only be stored with
        // mu held.
        //
        // Entries stored in read may be updated concurrently without mu, but updating
        // a previously-expunged entry requires that the entry be copied to the dirty
        // map and unexpunged with mu held.
        read atomic.Pointer[readOnly]
        // dirty 包含需要 mu 保存的映射内容部分。为了确保可以快速将 dirty 映射提升为 read 映射，它还包括 read 映射中所有未删除的条目。
        //
        // 已删除的条目不存储在 dirty 映射中。必须先取消删除 clean 映射中已删除的条目并将其添加到 dirty 映射中，然后才能将新值存储到该映射中。
        //
        // 如果 dirty 映射为 nil，则下次写入映射时将通过对 clean 映射进行浅拷贝来初始化它，从而省略陈旧的条目。
        // dirty contains the portion of the map's contents that require mu to be
        // held. To ensure that the dirty map can be promoted to the read map quickly,
        // it also includes all of the non-expunged entries in the read map.
        //
        // Expunged entries are not stored in the dirty map. An expunged entry in the
        // clean map must be unexpunged and added to the dirty map before a new value
        // can be stored to it.
        //
        // If the dirty map is nil, the next write to the map will initialize it by
        // making a shallow copy of the clean map, omitting stale entries.
        dirty map[any]*entry

        // misses counts the number of loads since the read map was last updated that
        // needed to lock mu to determine whether the key was present.
        //
        // Once enough misses have occurred to cover the cost of copying the dirty
        // map, the dirty map will be promoted to the read map (in the unamended
        // state) and the next store to the map will make a new dirty copy.
        misses int
}

type readOnly struct {
        m       map[any]*entry
        amended bool // true if the dirty map contains some key not in m.
}
type entry struct {
        // p points to the interface{} value stored for the entry.
        //
        // If p == nil, the entry has been deleted, and either m.dirty == nil or
        // m.dirty[key] is e.
        //
        // If p == expunged, the entry has been deleted, m.dirty != nil, and the entry
        // is missing from m.dirty.
        //
        // Otherwise, the entry is valid and recorded in m.read.m[key] and, if m.dirty
        // != nil, in m.dirty[key].
        //
        // An entry can be deleted by atomic replacement with nil: when m.dirty is
        // next created, it will atomically replace nil with expunged and leave
        // m.dirty[key] unset.
        //
        // An entry's associated value can be updated by atomic replacement, provided
        // p != expunged. If p == expunged, an entry's associated value can be updated
        // only after first setting m.dirty[key] = e so that lookups using the dirty
        // map find the entry.
        p atomic.Pointer[any]
}
```
