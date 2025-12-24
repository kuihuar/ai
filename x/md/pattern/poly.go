package pattern

import "fmt"

type Number interface {
	int | int64 | float32 | float64
}

func Add[T Number](a, b T) T {
	return a + b
}

// 接口定义
type SoundMaker interface {
	MakeSound() string
}

type Cat struct {
	Name string
}

// 结构体实现接口
func (c Cat) MakeSound() string {
	return fmt.Sprintf("%s miaoMiao", c.Name)
}

type Dog struct {
	Name string
}

func (d Dog) MakeSound() string {
	return fmt.Sprintf("%s wangwang", d.Name)
}

// 多态函数
func makeSound(s SoundMaker) {
	fmt.Println(s.MakeSound())
}

func ExamplePoly() {
	cat := Cat{Name: "huarhua"}
	dog := Dog{Name: "dahuang"}
	//函数内部调用MakeSound方法， 由于cat 实现了该接口，所以它们可以作为实例作为参数传入go
	makeSound(cat)
	makeSound(dog)
}
