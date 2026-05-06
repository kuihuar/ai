package dancebyte

import (
	"bytes"
	"context"
	"fmt"
	"hash/fnv"
	"sync"
)

func a() {
	ctx := context.Background()
	taskCh := make(chan int)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-taskCh:
				process(task)
			}
		}
	}(ctx)
}

func process(task int) {}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func log(data string) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString(data)
	fmt.Println(buf.String())
	bufPool.Put(buf)
}

type SegmentLock struct {
	segmentCnt int
	locks      []sync.Mutex
}

func NewSegmentLock(segmentCnt int) *SegmentLock {
	return &SegmentLock{
		segmentCnt: segmentCnt,
		locks:      make([]sync.Mutex, segmentCnt),
	}
}

func (s *SegmentLock) Lock(key string) {
	h := fnv.New64a()
	h.Write([]byte(key))
	index := int(h.Sum64() % uint64(s.segmentCnt))
	s.locks[index].Lock()
}

func (s *SegmentLock) Unlock(key string) {
	h := fnv.New64()
	h.Write([]byte(key))
	idx := int(h.Sum64() % uint64(s.segmentCnt))
	s.locks[idx].Unlock()
}
