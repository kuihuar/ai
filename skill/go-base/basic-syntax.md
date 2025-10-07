# Go 基础语法详解

## 📚 目录

- [变量和常量](#变量和常量)
- [数据类型](#数据类型)
- [控制结构](#控制结构)
- [函数](#函数)
- [数组和切片](#数组和切片)
- [映射](#映射)
- [结构体](#结构体)
- [指针](#指针)
- [接口](#接口)
- [类型系统](#类型系统)

## 变量和常量

### 变量声明

```go
package main

import "fmt"

func main() {
    // 1. 基本声明
    var name string
    name = "Go"
    
    // 2. 声明并初始化
    var age int = 25
    
    // 3. 类型推导
    var city = "Beijing"
    
    // 4. 简短声明（仅在函数内使用）
    count := 10
    
    // 5. 多变量声明
    var (
        firstName string = "John"
        lastName  string = "Doe"
        age2      int    = 30
    )
    
    // 6. 多变量简短声明
    x, y, z := 1, 2, 3
    
    fmt.Printf("Name: %s, Age: %d, City: %s\n", name, age, city)
    fmt.Printf("Count: %d\n", count)
    fmt.Printf("Person: %s %s, Age: %d\n", firstName, lastName, age2)
    fmt.Printf("Coordinates: x=%d, y=%d, z=%d\n", x, y, z)
}
```

### 常量

```go
package main

import "fmt"

func main() {
    // 1. 基本常量
    const pi = 3.14159
    const e = 2.71828
    
    // 2. 类型化常量
    const maxInt int = 9223372036854775807
    
    // 3. 常量组
    const (
        StatusOK    = 200
        StatusNotFound = 404
        StatusError = 500
    )
    
    // 4. iota 生成器
    const (
        Sunday = iota    // 0
        Monday           // 1
        Tuesday          // 2
        Wednesday        // 3
        Thursday         // 4
        Friday           // 5
        Saturday         // 6
    )
    
    // 5. 表达式常量
    const (
        KB = 1024
        MB = KB * 1024
        GB = MB * 1024
    )
    
    fmt.Printf("PI: %.5f\n", pi)
    fmt.Printf("Status codes: %d, %d, %d\n", StatusOK, StatusNotFound, StatusError)
    fmt.Printf("Days: Sunday=%d, Monday=%d\n", Sunday, Monday)
    fmt.Printf("Sizes: KB=%d, MB=%d, GB=%d\n", KB, MB, GB)
}
```

## 数据类型

### 基本类型

```go
package main

import "fmt"

func main() {
    // 整数类型
    var i8 int8 = 127
    var i16 int16 = 32767
    var i32 int32 = 2147483647
    var i64 int64 = 9223372036854775807
    
    // 无符号整数类型
    var u8 uint8 = 255
    var u16 uint16 = 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615
    
    // 浮点类型
    var f32 float32 = 3.14
    var f64 float64 = 3.141592653589793
    
    // 复数类型
    var c64 complex64 = 1 + 2i
    var c128 complex128 = 1 + 2i
    
    // 布尔类型
    var flag bool = true
    
    // 字符串类型
    var str string = "Hello, Go!"
    
    // 字节类型
    var b byte = 'A'  // byte 是 uint8 的别名
    
    // 符文类型
    var r rune = '好'  // rune 是 int32 的别名
    
    fmt.Printf("Integers: %d, %d, %d, %d\n", i8, i16, i32, i64)
    fmt.Printf("Unsigned: %d, %d, %d, %d\n", u8, u16, u32, u64)
    fmt.Printf("Floats: %.2f, %.15f\n", f32, f64)
    fmt.Printf("Complex: %v, %v\n", c64, c128)
    fmt.Printf("Bool: %t\n", flag)
    fmt.Printf("String: %s\n", str)
    fmt.Printf("Byte: %c (%d)\n", b, b)
    fmt.Printf("Rune: %c (%d)\n", r, r)
}
```

### 类型转换

```go
package main

import "fmt"
import "strconv"

func main() {
    // 数值类型转换
    var i int = 42
    var f float64 = float64(i)
    var u uint = uint(f)
    
    // 字符串转换
    var s1 string = strconv.Itoa(i)     // int to string
    var s2 string = strconv.FormatFloat(f, 'f', 2, 64)  // float to string
    var s3 string = string(rune(i))     // int to string (character)
    
    // 字符串解析
    if num, err := strconv.Atoi("123"); err == nil {
        fmt.Printf("Parsed int: %d\n", num)
    }
    
    if val, err := strconv.ParseFloat("3.14", 64); err == nil {
        fmt.Printf("Parsed float: %.2f\n", val)
    }
    
    fmt.Printf("Type conversions: %d -> %.2f -> %d\n", i, f, u)
    fmt.Printf("String conversions: %s, %s, %s\n", s1, s2, s3)
}
```

## 控制结构

### 条件语句

```go
package main

import "fmt"

func main() {
    // 1. 基本 if 语句
    age := 18
    if age >= 18 {
        fmt.Println("Adult")
    } else {
        fmt.Println("Minor")
    }
    
    // 2. if 语句中的变量声明
    if score := 85; score >= 90 {
        fmt.Println("Grade: A")
    } else if score >= 80 {
        fmt.Println("Grade: B")
    } else if score >= 70 {
        fmt.Println("Grade: C")
    } else {
        fmt.Println("Grade: F")
    }
    
    // 3. switch 语句
    day := "Monday"
    switch day {
    case "Monday":
        fmt.Println("Start of work week")
    case "Friday":
        fmt.Println("TGIF!")
    case "Saturday", "Sunday":
        fmt.Println("Weekend!")
    default:
        fmt.Println("Regular day")
    }
    
    // 4. 类型 switch
    var value interface{} = 42
    switch v := value.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

### 循环语句

```go
package main

import "fmt"

func main() {
    // 1. 基本 for 循环
    fmt.Println("Basic for loop:")
    for i := 0; i < 5; i++ {
        fmt.Printf("%d ", i)
    }
    fmt.Println()
    
    // 2. while 风格循环
    fmt.Println("While-style loop:")
    j := 0
    for j < 5 {
        fmt.Printf("%d ", j)
        j++
    }
    fmt.Println()
    
    // 3. 无限循环
    fmt.Println("Infinite loop (with break):")
    k := 0
    for {
        if k >= 5 {
            break
        }
        fmt.Printf("%d ", k)
        k++
    }
    fmt.Println()
    
    // 4. continue 语句
    fmt.Println("Loop with continue:")
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            continue
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()
    
    // 5. range 循环
    fmt.Println("Range loop over slice:")
    numbers := []int{10, 20, 30, 40, 50}
    for index, value := range numbers {
        fmt.Printf("Index: %d, Value: %d\n", index, value)
    }
    
    // 6. range 循环只获取值
    fmt.Println("Range loop (values only):")
    for _, value := range numbers {
        fmt.Printf("%d ", value)
    }
    fmt.Println()
    
    // 7. range 循环只获取索引
    fmt.Println("Range loop (indices only):")
    for index := range numbers {
        fmt.Printf("%d ", index)
    }
    fmt.Println()
}
```

## 函数

### 基本函数定义

```go
package main

import "fmt"

// 1. 基本函数
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 2. 函数返回值
func add(a, b int) int {
    return a + b
}

// 3. 多返回值
func divide(a, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

// 4. 命名返回值
func calculate(a, b int) (sum, product int) {
    sum = a + b
    product = a * b
    return  // 返回 sum 和 product
}

// 5. 变参函数
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

// 6. 函数作为参数
func apply(fn func(int) int, x int) int {
    return fn(x)
}

func square(x int) int {
    return x * x
}

func main() {
    // 调用基本函数
    greet("Go")
    
    // 调用返回函数
    result := add(5, 3)
    fmt.Printf("5 + 3 = %d\n", result)
    
    // 调用多返回函数
    q, r := divide(10, 3)
    fmt.Printf("10 / 3 = %d remainder %d\n", q, r)
    
    // 调用命名返回函数
    s, p := calculate(4, 5)
    fmt.Printf("4 + 5 = %d, 4 * 5 = %d\n", s, p)
    
    // 调用变参函数
    total1 := sum(1, 2, 3, 4, 5)
    total2 := sum()
    total3 := sum(10, 20)
    fmt.Printf("Sum: %d, %d, %d\n", total1, total2, total3)
    
    // 调用高阶函数
    squared := apply(square, 6)
    fmt.Printf("Square of 6 = %d\n", squared)
}
```

### 闭包和匿名函数

```go
package main

import "fmt"

func main() {
    // 1. 匿名函数
    func() {
        fmt.Println("Anonymous function")
    }()
    
    // 2. 匿名函数赋值给变量
    greet := func(name string) {
        fmt.Printf("Hello, %s!\n", name)
    }
    greet("Anonymous")
    
    // 3. 闭包
    counter := makeCounter()
    fmt.Printf("Counter: %d\n", counter())
    fmt.Printf("Counter: %d\n", counter())
    fmt.Printf("Counter: %d\n", counter())
    
    // 4. 闭包的实际应用
    multiplier := createMultiplier(5)
    fmt.Printf("5 * 3 = %d\n", multiplier(3))
    fmt.Printf("5 * 7 = %d\n", multiplier(7))
    
    // 5. 闭包捕获循环变量
    var funcs []func() int
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() int {
            return i  // 捕获的是循环结束后的 i 值
        })
    }
    
    fmt.Println("Closure capture issue:")
    for _, f := range funcs {
        fmt.Printf("%d ", f())
    }
    fmt.Println()
    
    // 6. 正确的闭包捕获
    var correctFuncs []func() int
    for i := 0; i < 3; i++ {
        i := i  // 创建局部变量
        correctFuncs = append(correctFuncs, func() int {
            return i
        })
    }
    
    fmt.Println("Correct closure capture:")
    for _, f := range correctFuncs {
        fmt.Printf("%d ", f())
    }
    fmt.Println()
}

// 创建计数器闭包
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// 创建乘数闭包
func createMultiplier(factor int) func(int) int {
    return func(x int) int {
        return factor * x
    }
}
```

## 数组和切片

### 数组

```go
package main

import "fmt"

func main() {
    // 1. 数组声明
    var arr1 [5]int                    // 零值数组
    var arr2 [5]int = [5]int{1, 2, 3, 4, 5}  // 初始化数组
    var arr3 = [...]int{1, 2, 3}      // 自动长度推导
    arr4 := [3]string{"Go", "Rust", "Python"}
    
    // 2. 数组操作
    fmt.Printf("Array 1: %v\n", arr1)
    fmt.Printf("Array 2: %v\n", arr2)
    fmt.Printf("Array 3: %v (length: %d)\n", arr3, len(arr3))
    fmt.Printf("Array 4: %v\n", arr4)
    
    // 3. 数组索引
    fmt.Printf("arr2[0] = %d\n", arr2[0])
    fmt.Printf("arr2[2] = %d\n", arr2[2])
    
    // 4. 数组长度和容量
    fmt.Printf("Length of arr2: %d\n", len(arr2))
    fmt.Printf("Length of arr4: %d\n", len(arr4))
    
    // 5. 数组遍历
    fmt.Println("Array iteration:")
    for i := 0; i < len(arr2); i++ {
        fmt.Printf("arr2[%d] = %d\n", i, arr2[i])
    }
    
    // 6. range 遍历数组
    fmt.Println("Range iteration:")
    for index, value := range arr4 {
        fmt.Printf("arr4[%d] = %s\n", index, value)
    }
}
```

### 切片

```go
package main

import "fmt"

func main() {
    // 1. 切片声明
    var s1 []int                       // nil 切片
    s2 := []int{1, 2, 3, 4, 5}        // 初始化切片
    s3 := make([]int, 5)               // 使用 make 创建
    s4 := make([]int, 5, 10)           // 长度为5，容量为10
    
    fmt.Printf("Slice 1: %v (len: %d, cap: %d)\n", s1, len(s1), cap(s1))
    fmt.Printf("Slice 2: %v (len: %d, cap: %d)\n", s2, len(s2), cap(s2))
    fmt.Printf("Slice 3: %v (len: %d, cap: %d)\n", s3, len(s3), cap(s3))
    fmt.Printf("Slice 4: %v (len: %d, cap: %d)\n", s4, len(s4), cap(s4))
    
    // 2. 切片操作
    fmt.Println("\n=== Slice Operations ===")
    
    // 追加元素
    s2 = append(s2, 6, 7, 8)
    fmt.Printf("After append: %v (len: %d, cap: %d)\n", s2, len(s2), cap(s2))
    
    // 切片切片
    subSlice := s2[1:4]
    fmt.Printf("Sub slice s2[1:4]: %v\n", subSlice)
    
    // 完整切片表达式
    fullSlice := s2[1:4:4]
    fmt.Printf("Full slice s2[1:4:4]: %v (len: %d, cap: %d)\n", 
               fullSlice, len(fullSlice), cap(fullSlice))
    
    // 3. 切片修改
    fmt.Println("\n=== Slice Modification ===")
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Printf("Original: %v\n", numbers)
    
    numbers[0] = 10
    fmt.Printf("After modify numbers[0]: %v\n", numbers)
    
    // 修改会影响原始数组
    originalArray := [5]int{1, 2, 3, 4, 5}
    sliceFromArray := originalArray[1:4]
    sliceFromArray[0] = 99
    fmt.Printf("Original array after slice modification: %v\n", originalArray)
    
    // 4. 切片复制
    fmt.Println("\n=== Slice Copy ===")
    source := []int{1, 2, 3, 4, 5}
    dest := make([]int, len(source))
    copy(dest, source)
    
    dest[0] = 99
    fmt.Printf("Source: %v\n", source)
    fmt.Printf("Dest: %v\n", dest)
    
    // 5. 切片删除元素
    fmt.Println("\n=== Slice Deletion ===")
    data := []int{1, 2, 3, 4, 5}
    fmt.Printf("Original: %v\n", data)
    
    // 删除索引为2的元素
    index := 2
    data = append(data[:index], data[index+1:]...)
    fmt.Printf("After deleting index 2: %v\n", data)
    
    // 6. 切片作为函数参数
    fmt.Println("\n=== Slice as Function Parameter ===")
    testSlice := []int{1, 2, 3, 4, 5}
    modifySlice(testSlice)
    fmt.Printf("Modified slice: %v\n", testSlice)
}

func modifySlice(s []int) {
    s[0] = 999
    s = append(s, 6, 7, 8) // 这个追加不会影响原始切片
}
```

## 映射

```go
package main

import "fmt"

func main() {
    // 1. 映射声明
    var m1 map[string]int                    // nil 映射
    m2 := make(map[string]int)               // 空映射
    m3 := map[string]int{                    // 初始化映射
        "apple":  5,
        "banana": 3,
        "orange": 8,
    }
    
    fmt.Printf("Map 1: %v (nil: %t)\n", m1, m1 == nil)
    fmt.Printf("Map 2: %v (nil: %t)\n", m2, m2 == nil)
    fmt.Printf("Map 3: %v\n", m3)
    
    // 2. 映射操作
    fmt.Println("\n=== Map Operations ===")
    
    // 添加/修改元素
    m2["key1"] = 100
    m2["key2"] = 200
    fmt.Printf("After adding elements: %v\n", m2)
    
    // 获取元素
    val1, ok1 := m2["key1"]
    val2, ok2 := m2["key3"]
    fmt.Printf("key1: %d (exists: %t)\n", val1, ok1)
    fmt.Printf("key3: %d (exists: %t)\n", val2, ok2)
    
    // 删除元素
    delete(m2, "key1")
    fmt.Printf("After deleting key1: %v\n", m2)
    
    // 3. 映射遍历
    fmt.Println("\n=== Map Iteration ===")
    
    for key, value := range m3 {
        fmt.Printf("%s: %d\n", key, value)
    }
    
    // 只遍历键
    fmt.Println("Keys only:")
    for key := range m3 {
        fmt.Printf("%s ", key)
    }
    fmt.Println()
    
    // 只遍历值
    fmt.Println("Values only:")
    for _, value := range m3 {
        fmt.Printf("%d ", value)
    }
    fmt.Println()
    
    // 4. 映射作为函数参数
    fmt.Println("\n=== Map as Function Parameter ===")
    modifyMap(m3)
    fmt.Printf("Modified map: %v\n", m3)
    
    // 5. 嵌套映射
    fmt.Println("\n=== Nested Maps ===")
    nested := map[string]map[string]int{
        "fruits": {
            "apple":  5,
            "banana": 3,
        },
        "vegetables": {
            "carrot": 10,
            "potato": 8,
        },
    }
    
    for category, items := range nested {
        fmt.Printf("%s:\n", category)
        for item, count := range items {
            fmt.Printf("  %s: %d\n", item, count)
        }
    }
}

func modifyMap(m map[string]int) {
    m["modified"] = 999
}
```

## 结构体

```go
package main

import "fmt"

// 1. 基本结构体定义
type Person struct {
    Name string
    Age  int
    City string
}

// 2. 结构体方法
func (p Person) Introduce() string {
    return fmt.Sprintf("Hi, I'm %s, %d years old, from %s", p.Name, p.Age, p.City)
}

// 3. 指针接收者方法
func (p *Person) HaveBirthday() {
    p.Age++
}

// 4. 嵌套结构体
type Address struct {
    Street string
    City   string
    State  string
    Zip    string
}

type Employee struct {
    Person  // 嵌入结构体
    Address // 嵌入结构体
    ID      int
    Salary  float64
}

// 5. 结构体标签
type User struct {
    ID       int    `json:"id" db:"user_id"`
    Username string `json:"username" db:"username"`
    Email    string `json:"email" db:"email"`
}

func main() {
    // 1. 创建结构体实例
    person1 := Person{
        Name: "Alice",
        Age:  30,
        City: "Beijing",
    }
    
    person2 := Person{"Bob", 25, "Shanghai"}
    
    person3 := Person{
        Name: "Charlie",
        Age:  35,
    } // City 为零值
    
    fmt.Printf("Person 1: %+v\n", person1)
    fmt.Printf("Person 2: %+v\n", person2)
    fmt.Printf("Person 3: %+v\n", person3)
    
    // 2. 访问结构体字段
    fmt.Printf("Person 1 name: %s\n", person1.Name)
    fmt.Printf("Person 1 age: %d\n", person1.Age)
    
    // 3. 调用结构体方法
    fmt.Println(person1.Introduce())
    
    // 4. 修改结构体字段
    person1.Age = 31
    fmt.Printf("After changing age: %+v\n", person1)
    
    // 5. 指针方法
    person1.HaveBirthday()
    fmt.Printf("After birthday: %+v\n", person1)
    
    // 6. 嵌套结构体
    employee := Employee{
        Person: Person{
            Name: "David",
            Age:  28,
        },
        Address: Address{
            Street: "123 Main St",
            City:   "New York",
            State:  "NY",
            Zip:    "10001",
        },
        ID:     1001,
        Salary: 75000.0,
    }
    
    fmt.Printf("Employee: %+v\n", employee)
    
    // 7. 访问嵌套字段
    fmt.Printf("Employee name: %s\n", employee.Name)        // 嵌入字段
    fmt.Printf("Employee city: %s\n", employee.City)        // 嵌入字段
    fmt.Printf("Employee street: %s\n", employee.Street)    // 嵌套字段
    
    // 8. 结构体指针
    employeePtr := &employee
    fmt.Printf("Employee via pointer: %+v\n", *employeePtr)
    
    // 9. 匿名结构体
    anonymous := struct {
        ID   int
        Name string
    }{
        ID:   999,
        Name: "Anonymous",
    }
    fmt.Printf("Anonymous struct: %+v\n", anonymous)
}
```

## 指针

```go
package main

import "fmt"

func main() {
    // 1. 基本指针操作
    x := 42
    p := &x  // p 是指向 x 的指针
    
    fmt.Printf("x = %d\n", x)
    fmt.Printf("p = %p\n", p)      // 指针的地址
    fmt.Printf("*p = %d\n", *p)    // 指针指向的值
    
    // 2. 通过指针修改值
    *p = 100
    fmt.Printf("After *p = 100, x = %d\n", x)
    
    // 3. 指针和函数
    fmt.Println("\n=== Pointers and Functions ===")
    
    a, b := 10, 20
    fmt.Printf("Before swap: a=%d, b=%d\n", a, b)
    swap(&a, &b)
    fmt.Printf("After swap: a=%d, b=%d\n", a, b)
    
    // 4. 指针和切片
    fmt.Println("\n=== Pointers and Slices ===")
    
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Printf("Original slice: %v\n", numbers)
    
    modifySliceWithPointer(&numbers)
    fmt.Printf("Modified slice: %v\n", numbers)
    
    // 5. 指针和数组
    fmt.Println("\n=== Pointers and Arrays ===")
    
    arr := [3]int{1, 2, 3}
    fmt.Printf("Original array: %v\n", arr)
    
    modifyArrayWithPointer(&arr)
    fmt.Printf("Modified array: %v\n", arr)
    
    // 6. 指针和结构体
    fmt.Println("\n=== Pointers and Structs ===")
    
    person := Person{Name: "Alice", Age: 30}
    fmt.Printf("Original person: %+v\n", person)
    
    modifyPersonWithPointer(&person)
    fmt.Printf("Modified person: %+v\n", person)
}

// 交换两个整数的值
func swap(x, y *int) {
    *x, *y = *y, *x
}

// 修改切片
func modifySliceWithPointer(s *[]int) {
    *s = append(*s, 6, 7, 8)
}

// 修改数组
func modifyArrayWithPointer(arr *[3]int) {
    for i := 0; i < len(arr); i++ {
        (*arr)[i] *= 2
    }
}

type Person struct {
    Name string
    Age  int
}

// 修改结构体
func modifyPersonWithPointer(p *Person) {
    p.Name = "Bob"
    p.Age = 25
}
```

## 接口

```go
package main

import "fmt"

// 1. 基本接口定义
type Writer interface {
    Write([]byte) (int, error)
}

type Reader interface {
    Read([]byte) (int, error)
}

// 2. 组合接口
type ReadWriter interface {
    Reader
    Writer
}

// 3. 空接口
type Any interface{}

// 4. 接口实现
type File struct {
    name string
    data []byte
}

func (f *File) Write(b []byte) (int, error) {
    f.data = append(f.data, b...)
    return len(b), nil
}

func (f *File) Read(b []byte) (int, error) {
    copy(b, f.data)
    return len(f.data), nil
}

// 5. 接口实现检查
func main() {
    // 1. 基本接口使用
    var w Writer = &File{name: "test.txt"}
    w.Write([]byte("Hello, World!"))
    fmt.Println("Written to file")
    
    // 2. 类型断言
    file := w.(*File)
    fmt.Printf("File content: %s\n", string(file.data))
    
    // 3. 安全的类型断言
    if f, ok := w.(*File); ok {
        fmt.Printf("File name: %s\n", f.name)
    }
    
    // 4. 类型开关
    var value interface{} = "Hello"
    
    switch v := value.(type) {
    case string:
        fmt.Printf("String: %s\n", v)
    case int:
        fmt.Printf("Integer: %d\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
    
    // 5. 空接口使用
    var any Any = 42
    fmt.Printf("Any value: %v\n", any)
    
    any = "Hello"
    fmt.Printf("Any value: %v\n", any)
    
    any = []int{1, 2, 3}
    fmt.Printf("Any value: %v\n", any)
    
    // 6. 接口作为函数参数
    processWriter(w)
    
    // 7. 返回接口
    rw := createReadWriter("example.txt")
    rw.Write([]byte("Test data"))
    
    data := make([]byte, 20)
    n, _ := rw.Read(data)
    fmt.Printf("Read data: %s\n", string(data[:n]))
}

func processWriter(w Writer) {
    w.Write([]byte("Processed data"))
}

func createReadWriter(name string) ReadWriter {
    return &File{name: name}
}
```

## 类型系统

### 类型定义

```go
package main

import "fmt"

// 1. 类型别名
type MyInt int
type MyString string

// 2. 自定义类型
type Temperature float64

func (t Temperature) Celsius() float64 {
    return float64(t)
}

func (t Temperature) Fahrenheit() float64 {
    return float64(t)*9/5 + 32
}

// 3. 方法集
type Counter int

func (c Counter) Value() int {
    return int(c)
}

func (c *Counter) Increment() {
    *c++
}

func (c *Counter) Decrement() {
    *c--
}

// 4. 接口组合
type Stringer interface {
    String() string
}

type Valuer interface {
    Value() int
}

type CounterStringer interface {
    Stringer
    Valuer
}

func (c Counter) String() string {
    return fmt.Sprintf("Counter: %d", c)
}

func main() {
    // 1. 类型别名使用
    var mi MyInt = 42
    var ms MyString = "Hello"
    
    fmt.Printf("MyInt: %d\n", mi)
    fmt.Printf("MyString: %s\n", ms)
    
    // 2. 自定义类型使用
    temp := Temperature(25.0)
    fmt.Printf("Temperature: %.1f°C, %.1f°F\n", 
               temp.Celsius(), temp.Fahrenheit())
    
    // 3. 方法集使用
    counter := Counter(10)
    fmt.Printf("Counter: %s\n", counter)
    fmt.Printf("Value: %d\n", counter.Value())
    
    counter.Increment()
    fmt.Printf("After increment: %s\n", counter)
    
    // 4. 接口组合使用
    var cs CounterStringer = counter
    fmt.Printf("CounterStringer: %s, Value: %d\n", cs.String(), cs.Value())
    
    // 5. 类型转换
    var i int = 100
    var c Counter = Counter(i)  // 显式转换
    
    fmt.Printf("Int to Counter: %s\n", c)
}
```

## 总结

Go 的基础语法具有以下特点：

1. **简洁明了**: 语法简单，易于学习和使用
2. **类型安全**: 静态类型系统，编译时检查错误
3. **高效性能**: 编译型语言，运行速度快
4. **并发友好**: 内置并发支持
5. **内存安全**: 自动垃圾回收

**核心概念**:
- 变量和常量声明简单直接
- 丰富的数据类型系统
- 强大的函数支持（包括闭包）
- 灵活的控制结构
- 动态数组（切片）和映射
- 结构体和方法
- 指针操作
- 接口系统

**最佳实践**:
- 使用简短声明在函数内部
- 优先使用切片而非数组
- 利用接口实现多态
- 合理使用指针避免复制
- 善用闭包和匿名函数

掌握这些基础语法是学习 Go 的重要一步，为后续的高级特性学习打下坚实基础。
