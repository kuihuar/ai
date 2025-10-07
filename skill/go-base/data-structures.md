# Go 基础数据结构详解

## 📚 目录

- [基本数据类型](#基本数据类型)
- [复合数据类型](#复合数据类型)
- [引用类型](#引用类型)
- [接口类型](#接口类型)
- [类型转换](#类型转换)
- [零值概念](#零值概念)
- [内存布局](#内存布局)
- [性能对比](#性能对比)

## 基本数据类型

### 数值类型

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    // 整数类型
    var i8 int8 = 127
    var i16 int16 = 32767
    var i32 int32 = 2147483647
    var i64 int64 = 9223372036854775807
    var i int = 42 // 平台相关，32位或64位
    
    // 无符号整数类型
    var u8 uint8 = 255
    var u16 uint16 = 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615
    var u uint = 42
    
    // 浮点类型
    var f32 float32 = 3.14
    var f64 float64 = 3.141592653589793
    
    // 复数类型
    var c64 complex64 = 1 + 2i
    var c128 complex128 = 1 + 2i
    
    // 字节类型 (uint8 的别名)
    var b byte = 'A'
    
    // Unicode 码点类型 (int32 的别名)
    var r rune = '中'
    
    fmt.Printf("int8: %d, size: %d bytes\n", i8, unsafe.Sizeof(i8))
    fmt.Printf("int16: %d, size: %d bytes\n", i16, unsafe.Sizeof(i16))
    fmt.Printf("int32: %d, size: %d bytes\n", i32, unsafe.Sizeof(i32))
    fmt.Printf("int64: %d, size: %d bytes\n", i64, unsafe.Sizeof(i64))
    fmt.Printf("int: %d, size: %d bytes\n", i, unsafe.Sizeof(i))
    fmt.Printf("float32: %f, size: %d bytes\n", f32, unsafe.Sizeof(f32))
    fmt.Printf("float64: %f, size: %d bytes\n", f64, unsafe.Sizeof(f64))
    fmt.Printf("complex64: %v, size: %d bytes\n", c64, unsafe.Sizeof(c64))
    fmt.Printf("complex128: %v, size: %d bytes\n", c128, unsafe.Sizeof(c128))
    fmt.Printf("byte: %c, size: %d bytes\n", b, unsafe.Sizeof(b))
    fmt.Printf("rune: %c, size: %d bytes\n", r, unsafe.Sizeof(r))
}
```

### 布尔类型

```go
package main

import "fmt"

func main() {
    var b1 bool = true
    var b2 bool = false
    var b3 bool // 零值为 false
    
    fmt.Printf("b1: %t\n", b1)
    fmt.Printf("b2: %t\n", b2)
    fmt.Printf("b3: %t\n", b3)
    
    // 布尔运算
    fmt.Printf("b1 && b2: %t\n", b1 && b2)
    fmt.Printf("b1 || b2: %t\n", b1 || b2)
    fmt.Printf("!b1: %t\n", !b1)
}
```

### 字符串类型

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    // 字符串声明和初始化
    var s1 string = "Hello, World!"
    var s2 string = `这是一个
多行字符串`
    
    // 字符串长度
    fmt.Printf("s1 length: %d\n", len(s1))
    fmt.Printf("s1 rune count: %d\n", utf8.RuneCountInString(s1))
    
    // 字符串是不可变的
    // s1[0] = 'h' // 编译错误
    
    // 字符串拼接
    s3 := s1 + " " + s2
    fmt.Printf("s3: %s\n", s3)
    
    // 使用 strings 包
    fmt.Printf("Contains: %t\n", strings.Contains(s1, "World"))
    fmt.Printf("Index: %d\n", strings.Index(s1, "World"))
    fmt.Printf("Replace: %s\n", strings.Replace(s1, "World", "Go", 1))
    fmt.Printf("ToUpper: %s\n", strings.ToUpper(s1))
    fmt.Printf("ToLower: %s\n", strings.ToLower(s1))
    
    // 字符串遍历
    fmt.Println("遍历字节:")
    for i := 0; i < len(s1); i++ {
        fmt.Printf("%c ", s1[i])
    }
    fmt.Println()
    
    fmt.Println("遍历字符:")
    for _, r := range s1 {
        fmt.Printf("%c ", r)
    }
    fmt.Println()
}
```

## 复合数据类型

### 数组

```go
package main

import "fmt"

func main() {
    // 数组声明
    var arr1 [5]int
    var arr2 [5]int = [5]int{1, 2, 3, 4, 5}
    var arr3 = [...]int{1, 2, 3, 4, 5} // 编译器推导长度
    
    // 数组初始化
    arr1[0] = 10
    arr1[1] = 20
    
    fmt.Printf("arr1: %v\n", arr1)
    fmt.Printf("arr2: %v\n", arr2)
    fmt.Printf("arr3: %v\n", arr3)
    
    // 数组长度
    fmt.Printf("arr1 length: %d\n", len(arr1))
    
    // 数组遍历
    fmt.Println("遍历数组:")
    for i, v := range arr2 {
        fmt.Printf("arr2[%d] = %d\n", i, v)
    }
    
    // 数组比较
    var arr4 [5]int = [5]int{1, 2, 3, 4, 5}
    fmt.Printf("arr2 == arr4: %t\n", arr2 == arr4)
    
    // 多维数组
    var matrix [3][3]int
    matrix[0] = [3]int{1, 2, 3}
    matrix[1] = [3]int{4, 5, 6}
    matrix[2] = [3]int{7, 8, 9}
    
    fmt.Println("矩阵:")
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            fmt.Printf("%d ", matrix[i][j])
        }
        fmt.Println()
    }
}
```

### 切片 (Slice)

```go
package main

import "fmt"

func main() {
    // 切片声明和初始化
    var s1 []int
    var s2 []int = []int{1, 2, 3, 4, 5}
    var s3 = make([]int, 5)        // 长度为5，容量为5
    var s4 = make([]int, 5, 10)    // 长度为5，容量为10
    
    fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))
    fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
    fmt.Printf("s3: %v, len: %d, cap: %d\n", s3, len(s3), cap(s3))
    fmt.Printf("s4: %v, len: %d, cap: %d\n", s4, len(s4), cap(s4))
    
    // 切片操作
    s5 := s2[1:3]  // 切片操作
    fmt.Printf("s5: %v, len: %d, cap: %d\n", s5, len(s5), cap(s5))
    
    // 追加元素
    s6 := append(s2, 6, 7, 8)
    fmt.Printf("s6: %v, len: %d, cap: %d\n", s6, len(s6), cap(s6))
    
    // 复制切片
    s7 := make([]int, len(s2))
    copy(s7, s2)
    fmt.Printf("s7: %v\n", s7)
    
    // 切片遍历
    fmt.Println("遍历切片:")
    for i, v := range s2 {
        fmt.Printf("s2[%d] = %d\n", i, v)
    }
    
    // 切片作为函数参数
    modifySlice(s2)
    fmt.Printf("修改后的s2: %v\n", s2)
}

func modifySlice(s []int) {
    if len(s) > 0 {
        s[0] = 999
    }
}
```

### 映射 (Map)

```go
package main

import "fmt"

func main() {
    // 映射声明和初始化
    var m1 map[string]int
    var m2 map[string]int = make(map[string]int)
    var m3 = map[string]int{
        "apple":  5,
        "banana": 3,
        "orange": 8,
    }
    
    // 映射操作
    m1 = make(map[string]int)
    m1["key1"] = 100
    m1["key2"] = 200
    
    fmt.Printf("m1: %v\n", m1)
    fmt.Printf("m2: %v\n", m2)
    fmt.Printf("m3: %v\n", m3)
    
    // 访问元素
    value, exists := m1["key1"]
    fmt.Printf("m1[\"key1\"]: %d, exists: %t\n", value, exists)
    
    // 删除元素
    delete(m1, "key1")
    fmt.Printf("删除key1后的m1: %v\n", m1)
    
    // 映射遍历
    fmt.Println("遍历映射:")
    for k, v := range m3 {
        fmt.Printf("m3[%s] = %d\n", k, v)
    }
    
    // 映射长度
    fmt.Printf("m3 length: %d\n", len(m3))
    
    // 映射作为函数参数
    modifyMap(m3)
    fmt.Printf("修改后的m3: %v\n", m3)
}

func modifyMap(m map[string]int) {
    m["grape"] = 10
}
```

### 结构体 (Struct)

```go
package main

import "fmt"

// 结构体定义
type Person struct {
    Name string
    Age  int
    City string
}

// 嵌套结构体
type Address struct {
    Street string
    City   string
    State  string
    Zip    string
}

type Employee struct {
    Person
    Address
    ID     int
    Salary float64
}

// 方法定义
func (p Person) String() string {
    return fmt.Sprintf("Person{Name: %s, Age: %d, City: %s}", p.Name, p.Age, p.City)
}

func (p *Person) SetAge(age int) {
    p.Age = age
}

func (e Employee) GetFullAddress() string {
    return fmt.Sprintf("%s, %s, %s %s", e.Street, e.City, e.State, e.Zip)
}

func main() {
    // 结构体初始化
    p1 := Person{"Alice", 30, "New York"}
    p2 := Person{Name: "Bob", Age: 25, City: "Los Angeles"}
    p3 := Person{Name: "Charlie"}
    
    fmt.Printf("p1: %v\n", p1)
    fmt.Printf("p2: %v\n", p2)
    fmt.Printf("p3: %v\n", p3)
    
    // 访问字段
    fmt.Printf("p1.Name: %s\n", p1.Name)
    fmt.Printf("p1.Age: %d\n", p1.Age)
    
    // 修改字段
    p1.SetAge(31)
    fmt.Printf("修改年龄后的p1: %v\n", p1)
    
    // 嵌套结构体
    emp := Employee{
        Person: Person{Name: "David", Age: 35, City: "Chicago"},
        Address: Address{
            Street: "123 Main St",
            City:   "Chicago",
            State:  "IL",
            Zip:    "60601",
        },
        ID:     1001,
        Salary: 75000.0,
    }
    
    fmt.Printf("emp: %+v\n", emp)
    fmt.Printf("emp.Name: %s\n", emp.Name) // 直接访问嵌套字段
    fmt.Printf("emp.FullAddress: %s\n", emp.GetFullAddress())
    
    // 结构体比较
    p4 := Person{"Alice", 31, "New York"}
    fmt.Printf("p1 == p4: %t\n", p1 == p4)
}
```

## 引用类型

### 指针

```go
package main

import "fmt"

func main() {
    // 指针声明
    var p *int
    var i int = 42
    
    // 取地址
    p = &i
    
    fmt.Printf("i: %d\n", i)
    fmt.Printf("p: %p\n", p)
    fmt.Printf("*p: %d\n", *p)
    
    // 通过指针修改值
    *p = 100
    fmt.Printf("修改后的i: %d\n", i)
    
    // 指针的指针
    var pp **int = &p
    fmt.Printf("pp: %p\n", pp)
    fmt.Printf("*pp: %p\n", *pp)
    fmt.Printf("**pp: %d\n", **pp)
    
    // 指针作为函数参数
    modifyValue(&i)
    fmt.Printf("函数修改后的i: %d\n", i)
    
    // 指针和切片
    arr := []int{1, 2, 3, 4, 5}
    modifySlice(arr)
    fmt.Printf("修改后的arr: %v\n", arr)
}

func modifyValue(p *int) {
    *p = 999
}

func modifySlice(s []int) {
    if len(s) > 0 {
        s[0] = 999
    }
}
```

### 函数类型

```go
package main

import "fmt"

// 函数类型定义
type Calculator func(int, int) int

// 函数作为参数
func calculate(a, b int, op Calculator) int {
    return op(a, b)
}

// 函数作为返回值
func getOperation(op string) Calculator {
    switch op {
    case "add":
        return func(a, b int) int { return a + b }
    case "subtract":
        return func(a, b int) int { return a - b }
    case "multiply":
        return func(a, b int) int { return a * b }
    case "divide":
        return func(a, b int) int { return a / b }
    default:
        return func(a, b int) int { return 0 }
    }
}

func main() {
    // 函数变量
    var add Calculator = func(a, b int) int { return a + b }
    var subtract Calculator = func(a, b int) int { return a - b }
    
    fmt.Printf("add(5, 3): %d\n", add(5, 3))
    fmt.Printf("subtract(5, 3): %d\n", subtract(5, 3))
    
    // 函数作为参数
    result1 := calculate(10, 5, add)
    result2 := calculate(10, 5, subtract)
    fmt.Printf("calculate(10, 5, add): %d\n", result1)
    fmt.Printf("calculate(10, 5, subtract): %d\n", result2)
    
    // 函数作为返回值
    addOp := getOperation("add")
    multiplyOp := getOperation("multiply")
    
    fmt.Printf("addOp(10, 5): %d\n", addOp(10, 5))
    fmt.Printf("multiplyOp(10, 5): %d\n", multiplyOp(10, 5))
}
```

## 接口类型

```go
package main

import "fmt"

// 接口定义
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 矩形结构体
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// 圆形结构体
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// 接口使用
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    // 接口实现
    var s Shape
    
    r := Rectangle{Width: 5, Height: 3}
    s = r
    fmt.Printf("Rectangle: ")
    printShapeInfo(s)
    
    c := Circle{Radius: 4}
    s = c
    fmt.Printf("Circle: ")
    printShapeInfo(s)
    
    // 类型断言
    if rect, ok := s.(Rectangle); ok {
        fmt.Printf("It's a rectangle with width %.2f and height %.2f\n", rect.Width, rect.Height)
    } else if circle, ok := s.(Circle); ok {
        fmt.Printf("It's a circle with radius %.2f\n", circle.Radius)
    }
    
    // 类型开关
    switch shape := s.(type) {
    case Rectangle:
        fmt.Printf("Rectangle: width=%.2f, height=%.2f\n", shape.Width, shape.Height)
    case Circle:
        fmt.Printf("Circle: radius=%.2f\n", shape.Radius)
    default:
        fmt.Println("Unknown shape")
    }
}
```

## 类型转换

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // 数值类型转换
    var i int = 42
    var f float64 = float64(i)
    var u uint = uint(i)
    
    fmt.Printf("int: %d, float64: %.2f, uint: %d\n", i, f, u)
    
    // 字符串转换
    str := strconv.Itoa(i)
    fmt.Printf("int to string: %s\n", str)
    
    num, err := strconv.Atoi("123")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("string to int: %d\n", num)
    }
    
    // 类型断言
    var value interface{} = "Hello, World!"
    
    if str, ok := value.(string); ok {
        fmt.Printf("Value is string: %s\n", str)
    }
    
    // 类型开关
    switch v := value.(type) {
    case string:
        fmt.Printf("String: %s\n", v)
    case int:
        fmt.Printf("Int: %d\n", v)
    case float64:
        fmt.Printf("Float64: %.2f\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

## 零值概念

```go
package main

import "fmt"

func main() {
    // 各种类型的零值
    var i int
    var f float64
    var b bool
    var s string
    var p *int
    var sl []int
    var m map[string]int
    var ch chan int
    var iface interface{}
    
    fmt.Printf("int zero value: %d\n", i)
    fmt.Printf("float64 zero value: %.2f\n", f)
    fmt.Printf("bool zero value: %t\n", b)
    fmt.Printf("string zero value: '%s'\n", s)
    fmt.Printf("pointer zero value: %v\n", p)
    fmt.Printf("slice zero value: %v\n", sl)
    fmt.Printf("map zero value: %v\n", m)
    fmt.Printf("channel zero value: %v\n", ch)
    fmt.Printf("interface zero value: %v\n", iface)
    
    // 零值检查
    if sl == nil {
        fmt.Println("slice is nil")
    }
    
    if m == nil {
        fmt.Println("map is nil")
    }
    
    if p == nil {
        fmt.Println("pointer is nil")
    }
}
```

## 内存布局

```go
package main

import (
    "fmt"
    "unsafe"
)

type Example struct {
    a bool    // 1 byte
    b int32   // 4 bytes
    c int64   // 8 bytes
    d string  // 16 bytes (8 + 8)
}

func main() {
    var e Example
    fmt.Printf("Example size: %d bytes\n", unsafe.Sizeof(e))
    fmt.Printf("Example alignment: %d bytes\n", unsafe.Alignof(e))
    
    // 字段偏移
    fmt.Printf("a offset: %d\n", unsafe.Offsetof(e.a))
    fmt.Printf("b offset: %d\n", unsafe.Offsetof(e.b))
    fmt.Printf("c offset: %d\n", unsafe.Offsetof(e.c))
    fmt.Printf("d offset: %d\n", unsafe.Offsetof(e.d))
    
    // 字段大小
    fmt.Printf("a size: %d\n", unsafe.Sizeof(e.a))
    fmt.Printf("b size: %d\n", unsafe.Sizeof(e.b))
    fmt.Printf("c size: %d\n", unsafe.Sizeof(e.c))
    fmt.Printf("d size: %d\n", unsafe.Sizeof(e.d))
}
```

## 性能对比

```go
package main

import (
    "fmt"
    "testing"
    "time"
)

// 数组性能测试
func BenchmarkArray(b *testing.B) {
    arr := [1000]int{}
    for i := 0; i < b.N; i++ {
        for j := 0; j < 1000; j++ {
            arr[j] = j
        }
    }
}

// 切片性能测试
func BenchmarkSlice(b *testing.B) {
    sl := make([]int, 1000)
    for i := 0; i < b.N; i++ {
        for j := 0; j < 1000; j++ {
            sl[j] = j
        }
    }
}

// 映射性能测试
func BenchmarkMap(b *testing.B) {
    m := make(map[int]int)
    for i := 0; i < b.N; i++ {
        for j := 0; j < 1000; j++ {
            m[j] = j
        }
    }
}

func main() {
    // 运行基准测试
    fmt.Println("运行基准测试...")
    
    // 数组测试
    start := time.Now()
    arr := [1000]int{}
    for i := 0; i < 1000; i++ {
        for j := 0; j < 1000; j++ {
            arr[j] = j
        }
    }
    arrayTime := time.Since(start)
    
    // 切片测试
    start = time.Now()
    sl := make([]int, 1000)
    for i := 0; i < 1000; i++ {
        for j := 0; j < 1000; j++ {
            sl[j] = j
        }
    }
    sliceTime := time.Since(start)
    
    // 映射测试
    start = time.Now()
    m := make(map[int]int)
    for i := 0; i < 1000; i++ {
        for j := 0; j < 1000; j++ {
            m[j] = j
        }
    }
    mapTime := time.Since(start)
    
    fmt.Printf("数组时间: %v\n", arrayTime)
    fmt.Printf("切片时间: %v\n", sliceTime)
    fmt.Printf("映射时间: %v\n", mapTime)
}
```

## 最佳实践

### 1. 选择合适的数据类型

```go
// 好的做法
var count uint32 = 0  // 明确表示非负数
var price float64 = 99.99  // 使用 float64 进行货币计算

// 避免的做法
var count int = 0  // 可能为负数
var price float32 = 99.99  // 精度可能不够
```

### 2. 使用切片而不是数组

```go
// 好的做法
func processItems(items []int) {
    // 处理切片
}

// 避免的做法
func processItems(items [100]int) {
    // 固定长度数组不够灵活
}
```

### 3. 合理使用映射

```go
// 好的做法
func getUserByID(users map[int]User, id int) (User, bool) {
    user, exists := users[id]
    return user, exists
}

// 避免的做法
func getUserByID(users map[int]User, id int) User {
    return users[id]  // 可能返回零值
}
```

### 4. 使用结构体方法

```go
// 好的做法
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++
}

func (c *Counter) Value() int {
    return c.value
}

// 避免的做法
func incrementCounter(c *Counter) {
    c.value++
}
```

### 5. 接口设计原则

```go
// 好的做法 - 小接口
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

// 避免的做法 - 大接口
type File interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
    Seek(int64, int) (int64, error)
    // ... 更多方法
}
```

## 总结

Go 的数据类型系统设计简洁而强大：

1. **基本类型**: 提供常用的数值、布尔、字符串类型
2. **复合类型**: 数组、切片、映射、结构体满足不同需求
3. **引用类型**: 指针、函数、接口提供高级功能
4. **类型安全**: 编译时类型检查，运行时类型断言
5. **零值概念**: 每个类型都有明确的零值
6. **内存效率**: 合理的内存布局和垃圾回收

选择合适的数据类型是编写高效 Go 程序的基础，理解各种类型的特点和使用场景对于 Go 开发者至关重要。
