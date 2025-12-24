package main

import (
	"net/http"
	_ "net/http/pprof"
)

// curl "http://localhost:8888/debug/pprof/trace?seconds=30s" > trace.out
func main() {
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}

// view trace by proc // 处理器视角
// Goroutine analysis // 协程视角
