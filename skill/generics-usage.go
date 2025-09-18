package skill

import "fmt"

// 2. 泛型函数
func Add[T int32 | int64 | float32 | float64](a, b T) T {
	return a + b
}

// 3. 泛型结构体
type Stack[T any] struct {
	data []T
}

// 4. 泛型方法
func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

// 5. 泛型接口
type Stackable[T any] interface {
	Push(v T)
}

// 6. 泛型约束
type Number interface {
	int32 | int64 | float32 | float64
}

// 7. 泛型约束的使用
func AddNumbers[T Number](a, b T) T {
	return a + b
}

// 8. 泛型接口的实现
func (s *Stack[T]) Pop() T {
	if len(s.data) == 0 {
		var zero T
		return zero
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v
}

// 9. 泛型结构体的使用
func DemoStack() {
	s := &Stack[int]{
		data: []int{1, 2, 3},
	}
	s.Push(4)
	fmt.Printf("s: %v\n", s)
}
