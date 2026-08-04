package generic

import (
	"cmp"
	"fmt"
)

// T 类型参数，any 等价 interface{}，允许任意类型
func Print[T any](v T) {
	fmt.Println(v)
}

type Stack[T any] struct{}

func (s *Stack[T]) Push(v T) {}

func GetMax[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Number 限制为几种基础数字类型
type Number interface {
	int | int64 | float32 | float64
}

func Sum[T Number](a, b T) T {
	return a + b
}

// 支持 int 以及任何底层是 int 的自定义类型
type FlexibleInt interface {
	~int | ~int64
}

type MyInt int // 自定义衍生类型

func Add[T FlexibleInt](a, b T) T {
	return a + b
}

//泛型使用注意事项（踩坑点）

// 1. 多个不同类型的泛型参数
func MapKeys[K comparable, V any](m map[K]V) []K { return make([]K, 0, len(m)) }

// 2. 泛型切片类型定义
type Vector[T any] []T

// 3. 泛型 Map 类型定义
type HashMap[K comparable, V any] map[K]V

// 4. 带有方法的接口约束
type Stringer interface {
	~string | ~[]byte
	String() string
}

// 泛型（Type-centric）：关注的是数据的类型本身。它解决“同一种数据结构或算法，如何作用于不同的数据类型”的问题（编译期确定类型）。

// 接口（Behavior-centric）：关注的是对象的行为/能力。它解决“不同的组件，如何暴露相同的行为进行解耦”的问题（运行期动态派发）。
