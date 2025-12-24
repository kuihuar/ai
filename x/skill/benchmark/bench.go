package main

import "unsafe"

func Bytes2strRaw(b []byte) string {
	return string(b)
}
func Bytes2strUnsafe(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func main() {

}
