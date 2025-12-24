package main

import (
	"runtime"
)

// "fmt"
// "sync"

type SmallNoPoointer struct {
	data [16]int64
}
type LargeNoPointer struct {
	data [17]int64
}

type Data struct {
	Value int
}

func main() {

	// var m sync.Map

	// m.Store("key1", "value1")
	// m.Store("key1", 1)
	// m.Store("key2", 2)

	// value, ok := m.Load("key1")
	// if ok {
	// 	fmt.Println("value:", value)
	// }

	// m.Delete("key1")
	// if ok {
	// 	fmt.Println("value:", value)
	// } else {
	// 	fmt.Println("key not found")
	// }
	// sm := make(map[int]SmallNoPoointer)
	// lm := make(map[int]LargeNoPointer)

	// m := make(map[int]*int)
	m := make(map[int]Data)
	for i := 0; i < 1e6; i++ {
		// sm[i] = SmallNoPoointer{}
		// lm[i] = LargeNoPointer{}
		// m[i] = new(int)
		m[i] = Data{Value: i}
	}
	runtime.GC()
	// time.Sleep(time.Second)
}
