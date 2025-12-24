package main

import (
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/bytedance/gopkg/util/gopool"
)

func TestFoo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Foo()
		})
	}
}

const BenchmarkTimes = 10000

func BenchmarkWorkerNoPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(BenchmarkTimes)
		for j := 0; j < BenchmarkTimes; j++ {
			go func() {
				Handle(0, 0)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	workNum := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool := NewWorkPool(workNum)
		for j := 0; j < benchmarkTimes; j++ {
			wg.Add(1)
			pool.addTask(func() {
				Handle(0, 0)
				wg.Done()
			})
		}
		wg.Wait()
	}
}
func BenchmarkWorkerPoolBD(b *testing.B) {
	workNum := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool := gopool.NewPool("test"+strconv.Itoa(i), int32(workNum), &gopool.Config{})
		for j := 0; j < benchmarkTimes; j++ {
			wg.Add(1)
			pool.Go(func() {
				Handle(0, 0)
				wg.Done()
			})
		}
		wg.Wait()
	}
}
