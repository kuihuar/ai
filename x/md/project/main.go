package main

import (
	"context"
	"sort"

	// "errors"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

type User struct {
	Age int
}

func (u User) ReflectCallFunc(name string) {
	fmt.Printf("age is %d, name %+v\n", u.Age, name)
}
func main() {
	age := 18
	pointerValue := reflect.ValueOf(&age)
	newValue := pointerValue.Elem()
	newValue.SetInt(20)
	fmt.Println(age)

	// pointerValue = reflect.ValueOf(age)
	// newValue = pointerValue.Elem()

	user := User{Age: 18}

	getValue := reflect.ValueOf(user)
	methodValue := getValue.MethodByName("ReflectCallFunc")
	args := []reflect.Value{reflect.ValueOf("zhangsan")}
	methodValue.Call(args)
	ret := Max[int](3, 4)
	ret1 := Max[float64](3.0, 4.0)
	fmt.Printf("%v\n", ret)
	fmt.Printf("%v\n", ret1)

}

func Max[T int | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func fn() *Obj {
	err := somefun()
	if err != nil {
		wrapErr := errors.Wrap(err, "fn error")
		fmt.Printf("%+v\n", wrapErr)
	}
	return &Obj{}
}

type Obj struct {
}

func handle(timeout time.Duration) *Obj {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan *Obj, 1)

	go func() {
		result := fn()
		select {
		case ch <- result:
		case <-ctx.Done():
		}
	}()
	select {
	case res := <-ch:
		return res
	case <-time.After(timeout):
		return nil
	}

}

func somefun() error {
	return errors.New("an error occurred")

}

type SliceFn[T any] struct {
	s    []T
	less func(T, T) bool
}

func (s SliceFn[T]) Len() int {
	return len(s.s)
}
func (s SliceFn[T]) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
}
func (s SliceFn[T]) Less(i, j int) bool {
	return s.less(s.s[i], s.s[j])
}
func SortFn[T any](s []T, less func(T, T) bool) {
	sort.Sort(SliceFn[T]{s, less})
}
