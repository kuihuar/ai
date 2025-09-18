package main

import (
	"fmt"
	"sort"
	"strconv"
)

// Stringer: 约束类型需实现 String() 方法
type Stringer interface{ String() string }

// ToStrings: 将实现 Stringer 的元素切片映射为字符串切片
func ToStrings[T Stringer](xs []T) []string {
	out := make([]string, 0, len(xs))
	for _, v := range xs {
		out = append(out, v.String())
	}
	return out
}

// User 与 Product 实现 String() 以用于 ToStrings
type User struct {
	ID   int
	Name string
}

func (u User) String() string { return fmt.Sprintf("%d:%s", u.ID, u.Name) }

type Product struct {
	SKU   string
	Price int
}

func (p Product) String() string { return fmt.Sprintf("%s:%d", p.SKU, p.Price) }

// DemoToStrings: 展示使用泛型接口约束进行映射
func DemoToStrings() {
	us := []User{{1, "A"}, {2, "B"}}
	ps := []Product{{"X", 100}, {"Y", 200}}
	fmt.Println(ToStrings(us)) // [1:A 2:B]
	fmt.Println(ToStrings(ps)) // [X:100 Y:200]
}

// Less 接口：定义可比较顺序的约束
type Less interface{ Less(other any) bool }

// Sortable: 基于 Less 的泛型可排序切片
type Sortable[T Less] []T

// Sort: 使用 sort.Slice 根据 Less 排序
func (s Sortable[T]) Sort() {
	sort.Slice(s, func(i, j int) bool { return s[i].Less(s[j]) })
}

// Score: 示例类型，按 V 字段升序
type Score struct {
	Name string
	V    int
}

func (a Score) Less(b any) bool { return a.V < b.(Score).V }

// DemoSortable: 展示基于接口约束的通用排序
func DemoSortable() {
	arr := Sortable[Score]{{"a", 3}, {"b", 1}, {"c", 2}}
	arr.Sort()
	fmt.Println(arr) // [{b 1} {c 2} {a 3}]
}

// Bag: 简单的泛型容器，支持添加和读取所有元素
type Bag[T any] struct{ items []T }

// NewBag: 构造函数；Add/All：添加与快照读取
func NewBag[T any](xs ...T) *Bag[T] { return &Bag[T]{items: append([]T(nil), xs...)} }
func (b *Bag[T]) Add(v T)           { b.items = append(b.items, v) }
func (b *Bag[T]) All() []T          { return append([]T(nil), b.items...) }

// MapBag: 对 Bag 内部元素做映射，返回新的 Bag
func MapBag[T any, R any](b *Bag[T], f func(T) R) *Bag[R] {
	out := make([]R, len(b.items))
	for i, v := range b.items {
		out[i] = f(v)
	}
	return &Bag[R]{items: out}
}

// DemoBag: 演示 Bag 的使用以及映射为其他类型
func DemoBag() {
	b := NewBag(1, 2, 3)
	b.Add(4)
	fmt.Println(b.All()) // [1 2 3 4]
	b2 := MapBag(b, func(x int) string { return strconv.Itoa(x) })
	fmt.Println(b2.All()) // ["1" "2" "3" "4"]
}
