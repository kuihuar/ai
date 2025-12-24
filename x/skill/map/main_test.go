package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func main_test(t *testing.T) {

}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test1kGcDuration(t *testing.T) {
	size := 1000
	m := GenerateStringMap(size)
	runtime.GC()
	gcCost := timeGC()
	t.Logf("size Td GC duration: %v\n", size, gcCost)
	runtime.KeepAlive(m)

}
func Test500wGcDuration(t *testing.T) {
	size := 5000000
	m := GenerateStringMap(size)
	runtime.GC()
	gcCost := timeGC()
	t.Logf("size Td GC duration: %v\n", size, gcCost)
	runtime.KeepAlive(m)
}

func GenerateStringMap(size int) map[string]string {
	m := make(map[string]string, size)
	for i := 0; i < size; i++ {
		m[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("val_%d", i)
	}
	return m
}
func timeGC() time.Duration {
	gcStartTime := time.Now()
	runtime.GC()
	gcCost := time.Since(gcStartTime)
	return gcCost
}

func Test500w2GcDuration(t *testing.T) {
	size := 5000000
	func() {
		m := GenerateStringMap(size)
		runtime.KeepAlive(m)
	}()

	runtime.GC()
	//var mem runtime.MemStats
	start := time.Now()
	runtime.GC()
	gcCost := time.Since(start)
	t.Logf("size Td GC duration: %v\n", size, gcCost)
}
