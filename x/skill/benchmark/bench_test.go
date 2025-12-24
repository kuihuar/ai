package main

import "testing"

func BenchmarkBytes2strRaw(b *testing.B) {
	aa := []byte("abcdefg")
	for n := 0; n < b.N; n++ {
		Bytes2strRaw(aa)
	}
}
func BenchmarkBytes2strUnsafe(b *testing.B) {
	aa := []byte("hijklmn")
	for n := 0; n < b.N; n++ {
		Bytes2strRaw(aa)
	}
}
