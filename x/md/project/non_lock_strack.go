package main

import (
	"sync/atomic"
	"unsafe"
)

type LockFreeStack struct {
	top unsafe.Pointer
}
type Node struct {
	next unsafe.Pointer
	val  interface{}
}

func NewLockFreeStrack() *LockFreeStack {
	return &LockFreeStack{}
}
func (s *LockFreeStack) Push(val interface{}) {
	item := &Node{val: val}
	for {
		top := atomic.LoadPointer(&s.top)
		item.next = top
		if atomic.CompareAndSwapPointer(&s.top, top, unsafe.Pointer(item)) {
			return
		}
	}
}
func (s *LockFreeStack) Pop() (interface{}, bool) {
	for {
		top := atomic.LoadPointer(&s.top)
		if top == nil {
			return nil, false
		}
		item := (*Node)(top)
		next := atomic.LoadPointer(&item.next)
		if atomic.CompareAndSwapPointer(&s.top, top, next) {
			return item.val, true
		}
	}

}
