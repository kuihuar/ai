package main

import (
	"sync"
	"testing"
	"time"
)

const (
	cost     = 1 * time.Millisecond
	readCnt  = 10000
	writeCnt = 1
)

func BenchmarkMutexReadMore(b *testing.B) {
	c := NewMutexCache()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for j := 0; j < 10000; j++ {

			for k := 0; k < readCnt; k++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					c.Get("key")
					// time.Sleep(cost)
				}()
			}

			for k := 0; k < writeCnt; k++ {
				wg.Add(1)
				go func(k int) {
					defer wg.Done()
					c.Set("key", "value")
					// time.Sleep(cost)
				}(k)
			}

		}
		wg.Wait()
	}
}
func BenchmarkRWMutexReadMore(b *testing.B) {
	c := NewRWMutexCache()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for j := 0; j < 1000; j++ {

			for k := 0; k < readCnt; k++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					c.Get("key")
					// time.Sleep(cost)
				}()
			}

			for k := 0; k < writeCnt; k++ {
				wg.Add(1)
				go func(k int) {
					defer wg.Done()
					c.Set("key", "value")
					// time.Sleep(cost)
				}(k)
			}

		}
		wg.Wait()
	}
}
