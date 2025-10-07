# Go 面向对象编程详解

## 📚 目录

- [结构体和方法](#结构体和方法)
- [接口](#接口)
- [组合与继承](#组合与继承)
- [多态](#多态)
- [封装](#封装)
- [设计模式](#设计模式)
- [最佳实践](#最佳实践)

## 结构体和方法

### 基本结构体定义

```go
package main

import "fmt"

// 1. 基本结构体
type Person struct {
    Name string
    Age  int
    City string
}

// 2. 结构体方法 - 值接收者
func (p Person) Introduce() string {
    return fmt.Sprintf("Hi, I'm %s, %d years old, from %s", p.Name, p.Age, p.City)
}

// 3. 结构体方法 - 指针接收者
func (p *Person) HaveBirthday() {
    p.Age++
}

func (p *Person) MoveTo(city string) {
    p.City = city
}

// 4. 构造函数模式
func NewPerson(name string, age int, city string) *Person {
    return &Person{
        Name: name,
        Age:  age,
        City: city,
    }
}

func main() {
    // 创建结构体实例
    person1 := Person{Name: "Alice", Age: 30, City: "Beijing"}
    person2 := NewPerson("Bob", 25, "Shanghai")
    
    // 调用方法
    fmt.Println(person1.Introduce())
    fmt.Println(person2.Introduce())
    
    // 修改状态
    person1.HaveBirthday()
    person1.MoveTo("Shanghai")
    fmt.Println("After changes:", person1.Introduce())
}
```

### 方法集规则

```go
package main

import "fmt"

type Counter int

// 值接收者方法
func (c Counter) Value() int {
    return int(c)
}

func (c Counter) String() string {
    return fmt.Sprintf("Counter: %d", c)
}

// 指针接收者方法
func (c *Counter) Increment() {
    *c++
}

func (c *Counter) Decrement() {
    *c--
}

func main() {
    // 值类型调用
    c1 := Counter(10)
    fmt.Println(c1.Value())    // 值接收者
    fmt.Println(c1.String())   // 值接收者
    
    // 指针类型调用
    c2 := &Counter(20)
    fmt.Println(c2.Value())    // 值接收者（通过指针）
    fmt.Println(c2.String())   // 值接收者（通过指针）
    
    // 只有指针可以调用指针接收者方法
    c2.Increment()             // 指针接收者
    c2.Decrement()             // 指针接收者
    
    // 值类型不能直接调用指针接收者方法
    // c1.Increment()  // 编译错误
    
    // 但可以通过取地址调用
    (&c1).Increment()          // 可以
    fmt.Println("c1 after increment:", c1)
}
```

## 接口

### 基本接口定义

```go
package main

import "fmt"

// 1. 基本接口
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

// 3. 接口实现
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

func main() {
    // 接口使用
    var w Writer = &File{name: "test.txt"}
    w.Write([]byte("Hello, World!"))
    
    // 类型断言
    if f, ok := w.(*File); ok {
        fmt.Printf("File content: %s\n", string(f.data))
    }
}
```

### 空接口和类型断言

```go
package main

import "fmt"

func main() {
    // 空接口可以存储任何类型
    var any interface{}
    
    any = 42
    processValue(any)
    
    any = "Hello"
    processValue(any)
    
    any = []int{1, 2, 3}
    processValue(any)
    
    any = map[string]int{"a": 1, "b": 2}
    processValue(any)
}

func processValue(v interface{}) {
    // 类型断言
    switch val := v.(type) {
    case int:
        fmt.Printf("Integer: %d\n", val)
    case string:
        fmt.Printf("String: %s\n", val)
    case []int:
        fmt.Printf("Slice: %v\n", val)
    case map[string]int:
        fmt.Printf("Map: %v\n", val)
    default:
        fmt.Printf("Unknown type: %T\n", val)
    }
}
```

## 组合与继承

### 结构体嵌入

```go
package main

import "fmt"

// 基础结构体
type Animal struct {
    Name string
    Age  int
}

func (a Animal) Speak() {
    fmt.Printf("%s makes a sound\n", a.Name)
}

// 嵌入结构体
type Dog struct {
    Animal  // 嵌入
    Breed   string
}

func (d Dog) Speak() {
    fmt.Printf("%s barks: Woof!\n", d.Name)
}

func (d Dog) Fetch() {
    fmt.Printf("%s fetches the ball\n", d.Name)
}

// 多层嵌入
type WorkingDog struct {
    Dog
    Job string
}

func (w WorkingDog) Work() {
    fmt.Printf("%s works as a %s\n", w.Name, w.Job)
}

func main() {
    // 创建实例
    dog := Dog{
        Animal: Animal{Name: "Buddy", Age: 3},
        Breed:  "Golden Retriever",
    }
    
    // 调用方法
    dog.Speak()    // 重写的方法
    dog.Fetch()    // 自己的方法
    
    // 访问嵌入字段
    fmt.Printf("Dog name: %s, age: %d\n", dog.Name, dog.Age)
    
    // 多层嵌入
    workingDog := WorkingDog{
        Dog: Dog{
            Animal: Animal{Name: "Rex", Age: 5},
            Breed:  "German Shepherd",
        },
        Job: "Police Dog",
    }
    
    workingDog.Speak()  // 继承的方法
    workingDog.Fetch()  // 继承的方法
    workingDog.Work()   // 自己的方法
}
```

### 接口组合

```go
package main

import "fmt"

// 基础接口
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type Closer interface {
    Close() error
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// 实现
type File struct {
    name string
    data []byte
}

func (f *File) Read(b []byte) (int, error) {
    copy(b, f.data)
    return len(f.data), nil
}

func (f *File) Write(b []byte) (int, error) {
    f.data = append(f.data, b...)
    return len(b), nil
}

func (f *File) Close() error {
    fmt.Printf("File %s closed\n", f.name)
    return nil
}

func main() {
    file := &File{name: "test.txt"}
    
    // 可以赋值给任何组合接口
    var r Reader = file
    var w Writer = file
    var rw ReadWriter = file
    var rwc ReadWriteCloser = file
    
    // 使用接口
    r.Read(make([]byte, 10))
    w.Write([]byte("Hello"))
    rw.Read(make([]byte, 10))
    rw.Write([]byte("World"))
    rwc.Close()
}
```

## 多态

### 接口多态

```go
package main

import "fmt"

// 形状接口
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 矩形
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

// 圆形
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// 三角形
type Triangle struct {
    A, B, C float64
}

func (t Triangle) Area() float64 {
    // 海伦公式
    s := (t.A + t.B + t.C) / 2
    return s * (s - t.A) * (s - t.B) * (s - t.C)
}

func (t Triangle) Perimeter() float64 {
    return t.A + t.B + t.C
}

// 多态函数
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    // 创建不同形状
    shapes := []Shape{
        Rectangle{Width: 5, Height: 3},
        Circle{Radius: 4},
        Triangle{A: 3, B: 4, C: 5},
    }
    
    // 多态调用
    for i, shape := range shapes {
        fmt.Printf("Shape %d: ", i+1)
        printShapeInfo(shape)
    }
}
```

### 工厂模式

```go
package main

import "fmt"

// 产品接口
type Product interface {
    GetName() string
    GetPrice() float64
}

// 具体产品
type Book struct {
    name  string
    price float64
}

func (b Book) GetName() string {
    return b.name
}

func (b Book) GetPrice() float64 {
    return b.price
}

type Laptop struct {
    name  string
    price float64
}

func (l Laptop) GetName() string {
    return l.name
}

func (l Laptop) GetPrice() float64 {
    return l.price
}

// 工厂
type ProductFactory struct{}

func (f ProductFactory) CreateProduct(productType string) Product {
    switch productType {
    case "book":
        return Book{name: "Go Programming", price: 29.99}
    case "laptop":
        return Laptop{name: "MacBook Pro", price: 1999.99}
    default:
        return nil
    }
}

func main() {
    factory := ProductFactory{}
    
    // 创建产品
    book := factory.CreateProduct("book")
    laptop := factory.CreateProduct("laptop")
    
    // 使用产品
    if book != nil {
        fmt.Printf("Book: %s, Price: $%.2f\n", book.GetName(), book.GetPrice())
    }
    
    if laptop != nil {
        fmt.Printf("Laptop: %s, Price: $%.2f\n", laptop.GetName(), laptop.GetPrice())
    }
}
```

## 封装

### 访问控制

```go
package main

import "fmt"

// 包级别的私有类型（小写开头）
type person struct {
    name string  // 私有字段
    age  int     // 私有字段
}

// 公开的构造函数
func NewPerson(name string, age int) *person {
    return &person{
        name: name,
        age:  age,
    }
}

// 公开的getter方法
func (p *person) GetName() string {
    return p.name
}

func (p *person) GetAge() int {
    return p.age
}

// 公开的setter方法
func (p *person) SetName(name string) {
    if len(name) > 0 {
        p.name = name
    }
}

func (p *person) SetAge(age int) {
    if age >= 0 {
        p.age = age
    }
}

// 公开的方法
func (p *person) Introduce() string {
    return fmt.Sprintf("Hi, I'm %s, %d years old", p.name, p.age)
}

func main() {
    // 只能通过构造函数创建
    p := NewPerson("Alice", 30)
    
    // 只能通过方法访问和修改
    fmt.Println(p.Introduce())
    
    p.SetName("Bob")
    p.SetAge(25)
    fmt.Println(p.Introduce())
    
    // 不能直接访问私有字段
    // fmt.Println(p.name)  // 编译错误
}
```

### 接口封装

```go
package main

import "fmt"

// 公开接口
type Database interface {
    Save(data string) error
    Load(id string) (string, error)
    Delete(id string) error
}

// 私有实现
type mysqlDB struct {
    host     string
    port     int
    username string
    password string
}

func (m *mysqlDB) Save(data string) error {
    fmt.Printf("Saving to MySQL: %s\n", data)
    return nil
}

func (m *mysqlDB) Load(id string) (string, error) {
    fmt.Printf("Loading from MySQL: %s\n", id)
    return "data", nil
}

func (m *mysqlDB) Delete(id string) error {
    fmt.Printf("Deleting from MySQL: %s\n", id)
    return nil
}

// 公开构造函数
func NewMySQLDatabase(host string, port int, username, password string) Database {
    return &mysqlDB{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }
}

func main() {
    // 通过接口使用，隐藏实现细节
    db := NewMySQLDatabase("localhost", 3306, "user", "pass")
    
    db.Save("test data")
    data, _ := db.Load("123")
    fmt.Printf("Loaded: %s\n", data)
    db.Delete("123")
}
```

## 设计模式

### 单例模式

```go
package main

import (
    "fmt"
    "sync"
)

type Singleton struct {
    name string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{name: "Singleton Instance"}
    })
    return instance
}

func (s *Singleton) GetName() string {
    return s.name
}

func main() {
    // 多次调用返回同一个实例
    s1 := GetInstance()
    s2 := GetInstance()
    
    fmt.Printf("s1 == s2: %t\n", s1 == s2)
    fmt.Printf("s1 name: %s\n", s1.GetName())
    fmt.Printf("s2 name: %s\n", s2.GetName())
}
```

### 观察者模式

```go
package main

import "fmt"

// 观察者接口
type Observer interface {
    Update(message string)
}

// 主题接口
type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify(message string)
}

// 具体主题
type NewsAgency struct {
    observers []Observer
}

func (n *NewsAgency) Attach(observer Observer) {
    n.observers = append(n.observers, observer)
}

func (n *NewsAgency) Detach(observer Observer) {
    for i, obs := range n.observers {
        if obs == observer {
            n.observers = append(n.observers[:i], n.observers[i+1:]...)
            break
        }
    }
}

func (n *NewsAgency) Notify(message string) {
    for _, observer := range n.observers {
        observer.Update(message)
    }
}

// 具体观察者
type NewsChannel struct {
    name string
}

func (nc *NewsChannel) Update(message string) {
    fmt.Printf("%s received: %s\n", nc.name, message)
}

func main() {
    agency := &NewsAgency{}
    
    channel1 := &NewsChannel{name: "CNN"}
    channel2 := &NewsChannel{name: "BBC"}
    
    agency.Attach(channel1)
    agency.Attach(channel2)
    
    agency.Notify("Breaking news: Go 1.20 released!")
    
    agency.Detach(channel1)
    agency.Notify("Update: New features added")
}
```

## 最佳实践

### 1. 接口设计原则

```go
package main

import "fmt"

// 好的接口设计：小而专注
type Writer interface {
    Write([]byte) (int, error)
}

type Closer interface {
    Close() error
}

// 组合小接口
type WriteCloser interface {
    Writer
    Closer
}

// 避免大而全的接口
// type BadInterface interface {
//     Read([]byte) (int, error)
//     Write([]byte) (int, error)
//     Close() error
//     Open() error
//     Flush() error
//     Seek(int64, int) (int64, error)
// }
```

### 2. 结构体设计

```go
package main

import "fmt"

// 好的结构体设计
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// 提供构造函数
func NewUser(username, email string) *User {
    return &User{
        Username: username,
        Email:    email,
    }
}

// 提供验证方法
func (u *User) IsValid() bool {
    return u.Username != "" && u.Email != ""
}

// 提供业务方法
func (u *User) GetDisplayName() string {
    if u.Username != "" {
        return u.Username
    }
    return u.Email
}
```

### 3. 错误处理

```go
package main

import (
    "errors"
    "fmt"
)

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// 业务逻辑
func ValidateUser(user *User) error {
    if user.Username == "" {
        return ValidationError{
            Field:   "username",
            Message: "username is required",
        }
    }
    
    if user.Email == "" {
        return ValidationError{
            Field:   "email",
            Message: "email is required",
        }
    }
    
    return nil
}

func main() {
    user := &User{Username: "", Email: "test@example.com"}
    
    if err := ValidateUser(user); err != nil {
        fmt.Printf("Error: %v\n", err)
        
        // 类型断言处理特定错误
        if validationErr, ok := err.(ValidationError); ok {
            fmt.Printf("Field: %s, Message: %s\n", 
                      validationErr.Field, validationErr.Message)
        }
    }
}
```

## 总结

Go 的面向对象编程具有以下特点：

1. **组合优于继承**: 通过嵌入实现代码复用
2. **接口驱动**: 通过接口实现多态和松耦合
3. **简洁明了**: 没有复杂的继承层次
4. **类型安全**: 编译时检查接口实现
5. **性能优秀**: 零成本抽象

**核心概念**:
- 结构体和方法
- 接口和类型断言
- 组合和嵌入
- 多态和工厂模式
- 封装和访问控制

**设计原则**:
- 接口要小而专注
- 优先使用组合而非继承
- 通过接口实现多态
- 合理使用封装
- 遵循 SOLID 原则

掌握这些面向对象特性，可以编写出更加灵活、可维护的 Go 代码。
