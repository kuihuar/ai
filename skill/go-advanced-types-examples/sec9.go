package main

import (
	"fmt"
	"strconv"
)

// Iterator: 抽象出“按次读取”元素的接口
type Iterator[T any] interface {
	Next() (T, bool)
}

// sliceIter: 基于切片的 Iterator 实现
type sliceIter[T any] struct{ i int; xs []T }
func (s *sliceIter[T]) Next() (T, bool) {
	if s.i >= len(s.xs) { var zero T; return zero, false }
	v := s.xs[s.i]
	s.i++
	return v, true
}

// FromSlice: 将切片适配为 Iterator
func FromSlice[T any](xs []T) Iterator[T] { return &sliceIter[T]{ xs: xs } }

// ReadAll: 读取迭代器的所有元素到切片
func ReadAll[T any](it Iterator[T]) []T {
	var out []T
	for {
		v, ok := it.Next(); if !ok { break }
		out = append(out, v)
	}
	return out
}

// DemoIterator: 展示 Iterator + ReadAll
func DemoIterator() {
	it := FromSlice([]int{1, 2, 3})
	fmt.Println(ReadAll(it)) // [1 2 3]
}

// CopyInto: 将迭代器的元素复制到固定大小缓冲区，返回写入数量
func CopyInto[T any](dst []T, it Iterator[T]) int {
	count := 0
	for {
		v, ok := it.Next(); if !ok { break }
		if count < len(dst) {
			dst[count] = v
			count++
		} else {
			break
		}
	}
	return count
}

// DemoCopyInto: 展示如何向固定缓冲区复制若干元素
func DemoCopyInto() {
	buf := make([]int, 2)
	it := FromSlice([]int{7, 8, 9})
	n := CopyInto(buf, it)
	fmt.Println(n, buf) // 2 [7 8]
}

// mapIter: Stream 映射时的迭代器实现（包级，便于泛型实例化）
type mapIter[T any, R any] struct{ src Iterator[T]; f func(T) R }

func (m *mapIter[T, R]) Next() (R, bool) {
	v, ok := m.src.Next(); if !ok { var zero R; return zero, false }
	return m.f(v), true
}

// Stream: 轻量级“流”，内部由 Iterator 驱动
type Stream[T any] struct{ it Iterator[T] }

// From: 将切片转换为 Stream
func From[T any](xs []T) Stream[T] { return Stream[T]{ it: FromSlice(xs) } }

// Map: 对流中元素进行函数式映射
func Map[T any, R any](s Stream[T], f func(T) R) Stream[R] {
	return Stream[R]{ it: &mapIter[T, R]{ src: s.it, f: f } }
}

// Collect: 将流收集为切片
func Collect[T any](s Stream[T]) []T { return ReadAll(s.it) }

// DemoStream: 展示 From + Map + Collect 的组合
func DemoStream() {
	res := Collect(Map(From([]int{1,2,3}), func(x int) string { return strconv.Itoa(x) }))
	fmt.Println(res) // ["1" "2" "3"]
}
