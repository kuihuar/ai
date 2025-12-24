package main

import (
	"math/rand"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("profile.pprof")
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()
	expensiveCPU()
}

// go:noline
func expensiveCPU() {
	var sum float32
	for i := 0; i < 10000000; i++ {
		sum += rand.Float32()
	}
	anotherExpensiveCPU()
}

// go:noline
func anotherExpensiveCPU() {
	var sum int
	for i := 0; i < 10000000; i++ {
		sum += rand.Intn(10)
	}
}
