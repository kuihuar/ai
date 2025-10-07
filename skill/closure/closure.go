package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 1. 基本闭包示例
func basicClosure() {
	fmt.Println("=== 基本闭包示例 ===")

	// 基本闭包
	add := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}

	add5 := add(5)
	fmt.Println("add5(3) =", add5(3)) // 输出: 8

	// 直接调用
	fmt.Println("add(10)(20) =", add(10)(20)) // 输出: 30
}

// 2. 闭包捕获外部变量
func captureVariable() {
	fmt.Println("\n=== 闭包捕获外部变量 ===")

	x := 10

	// 闭包捕获外部变量 x
	closure := func() int {
		return x * 2
	}

	fmt.Println("初始 x = 10, closure() =", closure()) // 输出: 20

	// 修改外部变量
	x = 20
	fmt.Println("修改 x = 20, closure() =", closure()) // 输出: 40
}

// 3. 闭包修改外部变量
func modifyVariable() {
	fmt.Println("\n=== 闭包修改外部变量 ===")

	counter := 0

	// 闭包可以修改外部变量
	increment := func() int {
		counter++
		return counter
	}

	fmt.Println("increment() =", increment()) // 输出: 1
	fmt.Println("increment() =", increment()) // 输出: 2
	fmt.Println("increment() =", increment()) // 输出: 3
}

// 4. 函数工厂模式
func functionFactory() {
	fmt.Println("\n=== 函数工厂模式 ===")

	// 创建加法器
	createAdder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}

	// 创建乘法器
	createMultiplier := func(x int) func(int) int {
		return func(y int) int {
			return x * y
		}
	}

	add10 := createAdder(10)
	multiply5 := createMultiplier(5)

	fmt.Println("add10(5) =", add10(5))         // 输出: 15
	fmt.Println("multiply5(3) =", multiply5(3)) // 输出: 15
}

// 5. 状态管理
func stateManagement() {
	fmt.Println("\n=== 状态管理 ===")

	// 计数器闭包
	createCounter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}

	// 累加器闭包
	createAccumulator := func(initial int) func(int) int {
		sum := initial
		return func(x int) int {
			sum += x
			return sum
		}
	}

	// 计数器
	counter := createCounter()
	fmt.Println("counter() =", counter()) // 输出: 1
	fmt.Println("counter() =", counter()) // 输出: 2
	fmt.Println("counter() =", counter()) // 输出: 3

	// 累加器
	acc := createAccumulator(10)
	fmt.Println("acc(5) =", acc(5)) // 输出: 15
	fmt.Println("acc(3) =", acc(3)) // 输出: 18
	fmt.Println("acc(2) =", acc(2)) // 输出: 20
}

// 6. 配置函数模式
type Config struct {
	Host    string
	Port    int
	Timeout int
}

type ConfigFunc func(*Config)

func WithHost(host string) ConfigFunc {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port int) ConfigFunc {
	return func(c *Config) {
		c.Port = port
	}
}

func WithTimeout(timeout int) ConfigFunc {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func applyConfig(config *Config, funcs ...ConfigFunc) {
	for _, f := range funcs {
		f(config)
	}
}

func configPattern() {
	fmt.Println("\n=== 配置函数模式 ===")

	config := &Config{}

	applyConfig(config,
		WithHost("localhost"),
		WithPort(8080),
		WithTimeout(30),
	)

	fmt.Printf("Config: %+v\n", config)
	// 输出: Config: {Host:localhost Port:8080 Timeout:30}
}

// 7. 循环中的闭包陷阱
func loopClosureTrap() {
	fmt.Println("\n=== 循环中的闭包陷阱 ===")

	// 错误示例 - 所有闭包都引用同一个变量
	var funcs []func() int
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() int {
			return i // 所有闭包都引用同一个 i
		})
	}

	fmt.Println("错误示例:")
	for i, f := range funcs {
		fmt.Printf("funcs[%d]() = %d\n", i, f()) // 输出: 3, 3, 3
	}
}

// 8. 循环中闭包的正确做法
func loopClosureCorrect() {
	fmt.Println("\n=== 循环中闭包的正确做法 ===")

	// 方法1: 通过参数传递
	var funcs []func() int
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func(val int) func() int {
			return func() int {
				return val
			}
		}(i))
	}

	fmt.Println("方法1 - 参数传递:")
	for i, f := range funcs {
		fmt.Printf("funcs[%d]() = %d\n", i, f()) // 输出: 0, 1, 2
	}

	// 方法2: 在循环内创建局部变量
	var funcs2 []func() int
	for i := 0; i < 3; i++ {
		val := i // 创建局部变量
		funcs2 = append(funcs2, func() int {
			return val
		})
	}

	fmt.Println("方法2 - 局部变量:")
	for i, f := range funcs2 {
		fmt.Printf("funcs2[%d]() = %d\n", i, f()) // 输出: 0, 1, 2
	}
}

// 9. 中间件模式
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

func AuthMiddleware(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+token {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func middlewarePattern() {
	fmt.Println("\n=== 中间件模式 ===")

	// 创建处理器
	// handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, World!")
	// })

	// 应用中间件
	//handler = LoggingMiddleware()(handler)
	//handler = AuthMiddleware("secret-token")(handler)

	fmt.Println("中间件已应用，处理器已准备就绪")
}

// 10. 回调函数
type ProcessFunc func(int) int

func processData(data []int, processor ProcessFunc) []int {
	result := make([]int, len(data))
	for i, v := range data {
		result[i] = processor(v)
	}
	return result
}

func callbackPattern() {
	fmt.Println("\n=== 回调函数模式 ===")

	data := []int{1, 2, 3, 4, 5}

	// 使用闭包创建处理函数
	double := func(x int) int {
		return x * 2
	}

	square := func(x int) int {
		return x * x
	}

	fmt.Println("原始数据:", data)
	fmt.Println("双倍处理:", processData(data, double)) // 输出: [2 4 6 8 10]
	fmt.Println("平方处理:", processData(data, square)) // 输出: [1 4 9 16 25]
}

// 11. 延迟执行
func delayedExecution(delay time.Duration, fn func()) func() {
	return func() {
		time.Sleep(delay)
		fn()
	}
}

func delayedExecutionPattern() {
	fmt.Println("\n=== 延迟执行模式 ===")

	// 创建延迟执行函数
	delayedPrint := delayedExecution(1*time.Second, func() {
		fmt.Println("延迟1秒执行")
	})

	fmt.Println("开始执行...")
	delayedPrint() // 1秒后执行
	fmt.Println("执行完成")
}

// 12. 内存泄漏示例
func memoryLeakExample() {
	fmt.Println("\n=== 内存泄漏示例 ===")

	// 可能导致内存泄漏的示例
	var bigData []int
	for i := 0; i < 1000; i++ {
		bigData = append(bigData, i)
	}

	// 闭包持有 bigData 的引用，即使不再使用
	closure := func() int {
		return len(bigData) // 持有 bigData 的引用
	}

	fmt.Println("闭包结果:", closure())

	// 解决方案：在不需要时显式释放
	bigData = nil
	fmt.Println("已释放 bigData")
}

func main() {
	basicClosure()
	captureVariable()
	modifyVariable()
	functionFactory()
	stateManagement()
	configPattern()
	loopClosureTrap()
	loopClosureCorrect()
	middlewarePattern()
	callbackPattern()
	delayedExecutionPattern()
	memoryLeakExample()
}
