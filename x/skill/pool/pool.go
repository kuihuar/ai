package main

import (
	"bytes"
	"sync"
)

var data = make([]byte, 1000)

func WriteBufferNoPool() {
	var buf bytes.Buffer
	buf.Write(data)
}

var objPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func WriteBufferWithPool() {
	buf := objPool.Get().(*bytes.Buffer)
	buf.Write(data)
	buf.Reset()
	objPool.Put(buf)
}
