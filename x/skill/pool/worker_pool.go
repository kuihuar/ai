package main

import (
	"sync"
)

func Foo() {

}

var benchmarkTimes = 10000

func Handle(a, b int) int {
	if b < 100 {
		return Handle(0, b+1)
	}
	return 0
}

type handleFunc func()

type WorkPool struct {
	WorkNum int
	tasks   chan handleFunc
}

func NewWorkPool(workNum int) *WorkPool {
	w := &WorkPool{WorkNum: workNum, tasks: make(chan handleFunc)}
	w.Start()
	return w
}
func (w *WorkPool) Start() {
	for i := 0; i < w.WorkNum; i++ {
		go func() {
			for task := range w.tasks {
				task()
			}
		}()
	}
}

func (w *WorkPool) addTask(task handleFunc) {
	w.tasks <- task
}

func useWorkPool() {
	pool := NewWorkPool(5)
	var wg sync.WaitGroup
	for j := 0; j < benchmarkTimes; j++ {

		wg.Add(1)
		pool.addTask(func() {
			Handle(0, 0)
			wg.Done()
		})
	}
	wg.Wait()
}
