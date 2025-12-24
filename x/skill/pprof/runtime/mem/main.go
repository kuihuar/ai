package main

import (
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("heap.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	runtime.GC()

	expensiveMem()

	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}
}

func expensiveMem() {
	m := make([]int, 10000000)

	for i := range m {
		m[i] = rand.Intn(127)
	}
	anotherExpensiveMem()
}

func anotherExpensiveMem() {
	m := make(map[int]float32, 10000000)
	for key := range m {
		m[key] = rand.Float32()
	}
}
