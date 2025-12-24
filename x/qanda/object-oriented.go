package qanda

import (
	"fmt"
	"io"
)

type User struct {
	Name string
	Age  int
}

func (u *User) UpdateNameByPointer(name string) {
	u.Name = name
}
func (u User) UpdateAgeByValue(age int) {
	u.Age = age
}

type UpdaterPointer interface {
	UpdateNameByPointer(string)
}

type UpdaterValue interface {
	UpdateAgeByValue(int)
}

func ExampleObjectriented() {
	var updaterPointer UpdaterPointer
	var updaterValue UpdaterValue
	user := User{
		Name: "zhangesan",
		Age:  18,
	}
	user.UpdateNameByPointer("lisi")
	updaterPointer = &user
	fmt.Println(updaterPointer)

	updaterValue = &user
	fmt.Println(updaterValue)
	updaterValue = user
	fmt.Println(updaterValue)
}

type Coder interface{ code() }
type Gopher struct{ name string }

func (g Gopher) code() { fmt.Printf("%s is coding\n", g.name) }
func ExampleInterIsEqual() {
	var c Coder
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)
	var g *Gopher
	fmt.Println(g == nil)
	c = g
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)

}

type MyError struct{}

func (i MyError) Error() string {
	return "MyError"
}
func Process() error {
	var err *MyError = nil
	return err
}
func ExampleInterIsEqual1() {

	err := Process()
	fmt.Println(err)
	fmt.Println(err == nil)

}

type myWriter struct{}

func (w *myWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func ExampleInterIsEqual2() {
	// 变量是 _, 类型是 o.Writer
	// 将nil 转化为 *myWriter类型的指针，即nil指针
	// 如果*myWriter未实现io.Writer接口，此赋值会失败
	var _ io.Writer = (*myWriter)(nil)

}

//类型断言和类型转换

// 类型转换，需要转换前后的两个类型要相互妆容才行
// 断言， go中的 所有类型都实现了 interface{} 空接口, 当一个函数的形参是interface{}, 在函数中就要断言，
// 目标类型，布尔参数 ：= 表达式.(目标类型)
// 安全断言
func ExampleTypeTransfor() {
	var i int = 9
	var f float64

	f = float64(i)
	fmt.Println(f)
}

func juge(v interface{}) {
	switch t := v.(type) {
	// case nil:
	case User:
		fmt.Printf("%p, %v", &t, v)
	case nil:
		fmt.Println("nil")
	default:
		fmt.Println("unknown")
	}
}
