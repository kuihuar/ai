# Go 闭包 (Closure) 详解

## 📖 概述

闭包是 Go 语言中的一个重要概念，它允许函数访问其外部作用域中的变量，即使外部函数已经返回。闭包在函数式编程、回调函数、状态管理等方面有广泛应用。

## 🎯 什么是闭包

闭包是一个函数值，它引用了其函数体之外的变量。该函数可以访问并赋予引用的变量的值，换句话说，该函数被"绑定"在了这些变量上。

## 🔧 基础语法

### 1. 基本闭包示例

```go
package main

import "fmt"

func main() {
    // 基本闭包
    add := func(x int) func(int) int {
        return func(y int) int {
            return x + y
        }
    }
    
    add5 := add(5)
    fmt.Println(add5(3)) // 输出: 8
    
    // 直接调用
    fmt.Println(add(10)(20)) // 输出: 30
}
```

### 2. 闭包捕获外部变量

```go
package main

import "fmt"

func main() {
    x := 10
    
    // 闭包捕获外部变量 x
    closure := func() int {
        return x * 2
    }
    
    fmt.Println(closure()) // 输出: 20
    
    // 修改外部变量
    x = 20
    fmt.Println(closure()) // 输出: 40
}
```

### 3. 闭包修改外部变量

```go
package main

import "fmt"

func main() {
    counter := 0
    
    // 闭包可以修改外部变量
    increment := func() int {
        counter++
        return counter
    }
    
    fmt.Println(increment()) // 输出: 1
    fmt.Println(increment()) // 输出: 2
    fmt.Println(increment()) // 输出: 3
}
```

## 🚀 高级用法

### 1. 函数工厂模式

```go
package main

import "fmt"

// 创建加法器
func createAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

// 创建乘法器
func createMultiplier(x int) func(int) int {
    return func(y int) int {
        return x * y
    }
}

func main() {
    add10 := createAdder(10)
    multiply5 := createMultiplier(5)
    
    fmt.Println(add10(5))      // 输出: 15
    fmt.Println(multiply5(3))  // 输出: 15
}
```

### 2. 状态管理

```go
package main

import "fmt"

// 计数器闭包
func createCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// 累加器闭包
func createAccumulator(initial int) func(int) int {
    sum := initial
    return func(x int) int {
        sum += x
        return sum
    }
}

func main() {
    // 计数器
    counter := createCounter()
    fmt.Println(counter()) // 输出: 1
    fmt.Println(counter()) // 输出: 2
    fmt.Println(counter()) // 输出: 3
    
    // 累加器
    acc := createAccumulator(10)
    fmt.Println(acc(5))  // 输出: 15
    fmt.Println(acc(3))  // 输出: 18
    fmt.Println(acc(2))  // 输出: 20
}
```

### 3. 配置函数

```go
package main

import "fmt"

// 配置结构体
type Config struct {
    Host string
    Port int
    Timeout int
}

// 配置函数类型
type ConfigFunc func(*Config)

// 设置主机
func WithHost(host string) ConfigFunc {
    return func(c *Config) {
        c.Host = host
    }
}

// 设置端口
func WithPort(port int) ConfigFunc {
    return func(c *Config) {
        c.Port = port
    }
}

// 设置超时
func WithTimeout(timeout int) ConfigFunc {
    return func(c *Config) {
        c.Timeout = timeout
    }
}

// 应用配置
func applyConfig(config *Config, funcs ...ConfigFunc) {
    for _, f := range funcs {
        f(config)
    }
}

func main() {
    config := &Config{}
    
    applyConfig(config,
        WithHost("localhost"),
        WithPort(8080),
        WithTimeout(30),
    )
    
    fmt.Printf("Config: %+v\n", config)
    // 输出: Config: {Host:localhost Port:8080 Timeout:30}
}
```

## 🔄 循环中的闭包

### 1. 常见陷阱

```go
package main

import "fmt"

func main() {
    // 错误示例 - 所有闭包都引用同一个变量
    var funcs []func() int
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() int {
            return i // 所有闭包都引用同一个 i
        })
    }
    
    for _, f := range funcs {
        fmt.Println(f()) // 输出: 3, 3, 3
    }
}
```

### 2. 正确做法

```go
package main

import "fmt"

func main() {
    // 方法1: 通过参数传递
    var funcs []func() int
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func(val int) func() int {
            return func() int {
                return val
            }
        }(i))
    }
    
    for _, f := range funcs {
        fmt.Println(f()) // 输出: 0, 1, 2
    }
    
    // 方法2: 在循环内创建局部变量
    var funcs2 []func() int
    for i := 0; i < 3; i++ {
        val := i // 创建局部变量
        funcs2 = append(funcs2, func() int {
            return val
        })
    }
    
    for _, f := range funcs2 {
        fmt.Println(f()) // 输出: 0, 1, 2
    }
}
```

## �� 实际应用场景

### 1. 中间件模式

```go
package main

import (
    "fmt"
    "log"
    "time"
)

// 中间件函数类型
type Middleware func(http.Handler) http.Handler

// 日志中间件
func LoggingMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
        })
    }
}

// 认证中间件
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
```

### 2. 回调函数

```go
package main

import "fmt"

// 处理函数类型
type ProcessFunc func(int) int

// 处理数据
func processData(data []int, processor ProcessFunc) []int {
    result := make([]int, len(data))
    for i, v := range data {
        result[i] = processor(v)
    }
    return result
}

func main() {
    data := []int{1, 2, 3, 4, 5}
    
    // 使用闭包创建处理函数
    double := func(x int) int {
        return x * 2
    }
    
    square := func(x int) int {
        return x * x
    }
    
    fmt.Println(processData(data, double))  // 输出: [2 4 6 8 10]
    fmt.Println(processData(data, square))  // 输出: [1 4 9 16 25]
}
```

### 3. 延迟执行

```go
package main

import (
    "fmt"
    "time"
)

// 延迟执行函数
func delayedExecution(delay time.Duration, fn func()) func() {
    return func() {
        time.Sleep(delay)
        fn()
    }
}

func main() {
    // 创建延迟执行函数
    delayedPrint := delayedExecution(2*time.Second, func() {
        fmt.Println("延迟2秒执行")
    })
    
    fmt.Println("开始执行...")
    delayedPrint() // 2秒后执行
    fmt.Println("执行完成")
}
```

## ⚠️ 注意事项

### 1. 内存泄漏

```go
package main

import "fmt"

func main() {
    // 可能导致内存泄漏的示例
    var bigData []int
    for i := 0; i < 1000000; i++ {
        bigData = append(bigData, i)
    }
    
    // 闭包持有 bigData 的引用，即使不再使用
    closure := func() int {
        return len(bigData) // 持有 bigData 的引用
    }
    
    fmt.Println(closure())
    
    // 解决方案：在不需要时显式释放
    bigData = nil
}
```

### 2. 变量捕获时机

```go
package main

import "fmt"

func main() {
    var funcs []func() int
    
    for i := 0; i < 3; i++ {
        // 注意：这里捕获的是 i 的地址，不是值
        funcs = append(funcs, func() int {
            return i
        })
    }
    
    // 当执行闭包时，i 的值已经是 3
    for _, f := range funcs {
        fmt.Println(f()) // 输出: 3, 3, 3
    }
}
```

## 📚 最佳实践

1. **明确闭包的生命周期**: 确保闭包不会持有不必要的引用
2. **避免在循环中直接使用闭包**: 使用参数传递或局部变量
3. **合理使用闭包进行状态管理**: 闭包适合简单的状态管理
4. **注意内存使用**: 闭包会持有外部变量的引用，可能导致内存泄漏
5. **使用闭包实现函数式编程**: 闭包是实现高阶函数的基础

## 🔗 相关资源

- [Go 官方文档 - 函数值](https://golang.org/ref/spec#Function_types)
- [Go 官方博客 - 函数式编程](https://blog.golang.org/function-values)
- [Go 闭包详解](https://golang.org/doc/effective_go.html#closures)
