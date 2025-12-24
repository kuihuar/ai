package main

import "fmt"

type Number interface {
	int | float64
}

func Add[T Number](nums []T) T {
	var sum T

	for _, v := range nums {
		sum += v
	}
	return sum
}

func main() {
	ints := []int{1, 2, 3}
	floats := []float64{1.1, 2.2, 3.3}

	fmt.Println(Add(ints))
	fmt.Println(Add(floats))

	intStack := &Stack[int]{}

	intStack.Push(1)
	fmt.Println(intStack.Pop())

	var candyBox Box[string] = &CandyBox{}
	candyBox.Put("candy")
	fmt.Println(candyBox.Get())

}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

type Box[T any] interface {
	Put(item T)
	Get() T
}
type CandyBox struct {
	item string
}

func (c *CandyBox) Put(item string) {
	c.item = item
}
func (c *CandyBox) Get() string {
	return c.item
}

type Calculator[T int | float64] interface {
	Add(a, b T) T
}

var intCalc Calculator[int]
intCalc.Add(1, 2)