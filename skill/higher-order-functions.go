package skill

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// 函数式编程三大核心操作：
// Map：一对一转换（[a, b, c] → [f(a), f(b), f(c)]）
// Filter：筛选元素（[a, b, c] → [a, c]）
// Reduce：多对一归约（[a, b, c] → f(f(initial, a), b), c)）

// 1. 高阶函数作为参数 - Map函数
// 将切片中的每个元素通过转换函数进行一对一映射
func MapWithHiger[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v) // 对每个元素应用转换函数
	}
	return result
}

// 2. 高阶函数作为返回值 - 工厂函数模式
// 返回一个乘法器函数，实现函数工厂
func CreateMultipliterWithHiger(factor int) func(int) int {
	return func(x int) int {
		return x * factor // 闭包捕获了外部的factor变量
	}
}

// 3. 闭包 - 函数内部引用了外部变量 - 计数器
// 返回一个闭包函数，维护内部状态
func CreateCounterWithHiger() func() int {
	count := 0 // 外部变量，被内部函数引用
	return func() int {
		count++ // 修改外部变量
		return count
	}
}

// 4. 作为参数 - 过滤器
// 根据谓词函数筛选切片中的元素
func FilterWithHiger[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) { // 使用谓词函数判断是否保留元素
			result = append(result, v)
		}
	}
	return result
}

// 5. 作为参数 - 归约函数
// 将切片归约为单个值，常用于求和、求积等操作
func ReduceWithHiger[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, v := range slice {
		result = reducer(result, v) // 累积操作
	}
	return result
}

// 6. 作为参数 - 排序函数
// 使用自定义比较函数对切片进行排序
func SortWithHiger[T any](slice []T, less func(T, T) bool) []T {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j]) // 使用自定义比较函数
	})
	return slice
}

// 7. 作为返回值 - 装饰器模式
// 为函数添加额外功能而不修改原函数
func DecorateWithHiger[T any](fn func(T) T, decorator func(func(T) T) func(T) T) func(T) T {
	return decorator(fn) // 返回装饰后的函数
}

// 装饰器示例：添加日志功能
func WithLoggingWithHiger(fn func(int) int) func(int) int {
	return func(x int) int {
		fmt.Println("Calling function with input:", x)
		result := fn(x)
		fmt.Println("Function returned:", result)
		return result
	}
}

// 8. 作为返回值 - 函数组合
// 将两个函数组合成一个新函数：f(g(x))
func ComposeWithHiger[T any](f func(T) T, g func(T) T) func(T) T {
	return func(x T) T {
		return f(g(x)) // 先执行g，再执行f
	}
}

// 9. 作为返回值 - 函数柯里化
// 函数柯里化（Currying）是一种将接受多个参数的函数转换为一系列接受单个参数的函数的技术。
func CurryWithHiger[T any, R any](fn func(T) R) func(T) R {
	return func(x T) R {
		return fn(x)
	}
}

// 10. 作为返回值 - 函数记忆化
// 缓存函数结果，避免重复计算
func MemoizeWithHiger[T comparable, R any](fn func(T) R) func(T) R {
	cache := make(map[T]R)

	return func(x T) R {
		if val, ok := cache[x]; ok {
			return val // 返回缓存的结果
		}
		result := fn(x)
		cache[x] = result // 缓存新结果
		return result
	}
}

// 11. 作为返回值 - 防抖函数
// 延迟执行函数，如果在延迟期间再次调用则重新计时
func DebounceWithHiger[T any](fn func(T), delay time.Duration) func(T) {
	var timer *time.Timer
	return func(x T) {
		if timer != nil {
			timer.Stop() // 取消之前的定时器
		}
		timer = time.AfterFunc(delay, func() {
			fn(x) // 延迟执行
		})
	}
}

// 12. 作为返回值 - 节流函数
// 限制函数执行频率，确保在指定时间内最多执行一次
func ThrottleWithHiger[T any](fn func(T), interval time.Duration) func(T) {
	var lastCall time.Time
	return func(x T) {
		now := time.Now()
		if now.Sub(lastCall) >= interval {
			fn(x)
			lastCall = now
		}
	}
}

// 13. 作为参数 - 异步执行器
// 异步执行函数并返回结果通道
func AsyncExecuteWithHiger[T any, R any](fn func(T) R) func(T) <-chan R {
	return func(x T) <-chan R {
		result := make(chan R, 1)
		go func() {
			result <- fn(x) // 在goroutine中执行
			close(result)
		}()
		return result
	}
}

// 14. 作为返回值 - 重试函数
// 为函数添加重试机制
func WithRetryWithHiger[T any](fn func(T) error, maxRetries int) func(T) error {
	return func(x T) error {
		var err error
		for i := 0; i <= maxRetries; i++ {
			if err = fn(x); err == nil {
				return nil // 成功则返回
			}
			if i < maxRetries {
				time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			}
		}
		return err
	}
}

// 15. 作为返回值 - 超时控制
// 为函数添加超时控制
func WithTimeoutWithHiger[T any, R any](fn func(T) R, timeout time.Duration) func(T) (R, error) {
	return func(x T) (R, error) {
		resultChan := make(chan R, 1)
		errorChan := make(chan error, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					errorChan <- fmt.Errorf("panic: %v", r)
				}
			}()
			resultChan <- fn(x)
		}()

		select {
		case result := <-resultChan:
			return result, nil
		case err := <-errorChan:
			var zero R
			return zero, err
		case <-time.After(timeout):
			var zero R
			return zero, fmt.Errorf("function execution timed out")
		}
	}
}

// 16. 作为参数 - 管道操作
// 将多个函数串联成管道
func PipelineWithHiger[T any](functions ...func(T) T) func(T) T {
	return func(x T) T {
		result := x
		for _, fn := range functions {
			result = fn(result) // 依次执行每个函数
		}
		return result
	}
}

// 17. 作为返回值 - 部分应用
// 固定函数的部分参数，返回接受剩余参数的函数
func PartialWithHiger[T1, T2, R any](fn func(T1, T2) R, arg1 T1) func(T2) R {
	return func(arg2 T2) R {
		return fn(arg1, arg2) // 固定第一个参数
	}
}

// 18. 作为返回值 - 惰性求值
// 延迟函数执行，只在需要时才计算
func LazyWithHiger[T any, R any](fn func(T) R) func(T) func() R {
	return func(x T) func() R {
		var result R
		var computed bool
		return func() R {
			if !computed {
				result = fn(x)
				computed = true
			}
			return result
		}
	}
}

// 19. 作为参数 - 错误处理装饰器
// 为函数添加统一的错误处理
func WithErrorHandlerWithHiger[T any, R any](fn func(T) (R, error), handler func(error)) func(T) R {
	return func(x T) R {
		result, err := fn(x)
		if err != nil {
			handler(err) // 统一错误处理
		}
		return result
	}
}

// 20. 作为返回值 - 中间件模式
// 为函数添加中间件处理链
func WithMiddlewareWithHiger[T any, R any](fn func(T) R, middlewares ...func(func(T) R) func(T) R) func(T) R {
	result := fn
	for i := len(middlewares) - 1; i >= 0; i-- {
		result = middlewares[i](result) // 从右到左应用中间件
	}
	return result
}

// 示例函数
func squareWithHiger(x int) int {
	return x * x
}

func isEven(x int) bool {
	return x%2 == 0
}

func add(a, b int) int {
	return a + b
}

// 演示高阶函数的使用
func DemonstrateHigherOrderFunctions() {
	fmt.Println("=== 高阶函数演示 ===")

	// 1. 作为参数使用
	numbers := []int{1, 2, 3, 4, 5}
	squared := MapWithHiger(numbers, squareWithHiger)
	fmt.Printf("原始数组: %v\n", numbers)
	fmt.Printf("平方后: %v\n", squared)

	// 使用匿名函数作为参数
	uppercase := MapWithHiger([]string{"hello", "world"}, strings.ToUpper)
	fmt.Printf("转大写: %v\n", uppercase)

	// 2. 作为返回值使用
	double := CreateMultipliterWithHiger(2)
	triple := CreateMultipliterWithHiger(3)
	fmt.Printf("2的3倍: %d\n", double(3))
	fmt.Printf("3的4倍: %d\n", triple(4))

	// 3. 闭包使用
	counter := CreateCounterWithHiger()
	fmt.Printf("计数器: %d\n", counter()) // 1
	fmt.Printf("计数器: %d\n", counter()) // 2
	fmt.Printf("计数器: %d\n", counter()) // 3

	// 4. 过滤器使用
	evenNumbers := FilterWithHiger(numbers, isEven)
	fmt.Printf("偶数: %v\n", evenNumbers)

	// 5. 归约使用
	sum := ReduceWithHiger(numbers, 0, add)
	fmt.Printf("求和: %d\n", sum)

	// 6. 装饰器使用
	loggedSquare := WithLoggingWithHiger(squareWithHiger)
	result := loggedSquare(5)
	fmt.Printf("最终结果: %d\n", result)

	// 7. 记忆化使用
	memoizedSquare := MemoizeWithHiger(squareWithHiger)
	fmt.Printf("记忆化平方(5): %d\n", memoizedSquare(5))
	fmt.Printf("记忆化平方(5): %d\n", memoizedSquare(5)) // 从缓存获取

	// 8. 管道使用
	pipeline := PipelineWithHiger(squareWithHiger, func(x int) int { return x + 1 })
	fmt.Printf("管道结果: %d\n", pipeline(3)) // (3^2) + 1 = 10

	// 9. 部分应用使用
	addFive := PartialWithHiger(add, 5)
	fmt.Printf("部分应用: %d\n", addFive(3)) // 5 + 3 = 8
}
