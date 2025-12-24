package main

import "testing"

func TestWriteBufferNoPool(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteBufferNoPool()
		})
	}
}

func BenchmarkWriteBufferNoPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//其它请求处理
		WriteBufferNoPool()
	}
}

func BenchmarkWriteBufferWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//其它请求处理
		WriteBufferWithPool()
	}
}
