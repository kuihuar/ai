### 结构体嵌套结构体

场景：基础模型扩展

```go
type Base struct { ID int }
type User struct {
    Base  // 嵌入：字段提升，可以直接访问 user.ID
    Name string
}
```

### 结构体嵌入接口
场景：行为可动态替换、策略模式、依赖注入 → 用「嵌入接口」
结构体里面组合了一个接口字段，是「我持有一个接口，借用它的能力」。

```go
type Reader interface {
    Read() (int, error)
}
type User struct {
    ID int
    Name string
    Reader

}
```

### 结构体实现接口
场景：固定行为、抽象规范 → 用「实现接口」
结构体里面实现了接口的所有方法，是「我是一个接口，我有所有方法」。

```go
type Reader interface {
    Read() (int, error)
}
type User struct {
    ID int
    Name string
}
func (u *User) Read() (int, error) {
    return 0, nil
}
```


### 结构体包含函数
场景：行为可动态替换、策略模式、依赖注入 → 用「包含函数」
调函数
钩子函数（中间件、事件处理）
简单策略（只需要一个行为）
动态替换逻辑（运行时改行为）
结构体里面包含了一个函数字段，是「我持有一个函数，调用它的能力」。

```go
type User struct {
    ID int
    Name string
    Read func() (int, error)
}
``` 


### 结构体包含接口（指针）
场景：行为可动态替换、策略模式、依赖注入 → 用「包含接口（指针）」
结构体里面包含了一个接口字段，是「我持有一个接口，调用它的能力」。

```go
type User struct {
    ID int
    Name string
    Reader *Reader
}
```

### 结构体包含接口（值）
场景：行为可动态替换、策略模式、依赖注入 → 用「包含接口（值）」
结构体里面包含了一个接口字段，是「我持有一个接口，调用它的能力」。

```go
type User struct {
    ID int
    Name string
    Reader Reader
}
```

### 结构体嵌入函数别名类型
场景：行为可动态替换、策略模式、依赖注入 → 用「嵌入函数别名类型」                           
结构体嵌入函数别名类型，是「我持有一个函数，调用它的能力」。


### 接口嵌入接口
场景：接口嵌套接口，合并方法集，最常用在标准库。

### 普通函数 实现 接口


### 递归组合（自己嵌套自己）


## 3. 函数 (Function) 的高级组合

### A. 函数实现接口（适配器模式）
正如我们之前讨论的 `http.HandlerFunc`。
*   **场景**：让普通函数满足接口要求。
*   **用法**：定义一个函数类型，并给它写一个满足接口的方法，方法内部调用函数自己。

### B. 函数返回函数（闭包/工厂）
*   **场景**：中间件、装饰器、延迟执行。
*   **用法**：

## 2. 接口 (Interface) 的组合与嵌套

### A. 接口嵌套接口（组合）
Go 不鼓励大接口，鼓励通过小接口组合成大接口。
*   **场景**：标准库的 `ReadWriter`。
*   **用法**：
    
```go
    type Reader interface { Read() }
    type Writer interface { Write() }
    type ReadWriter interface {
        Reader // 嵌入
        Writer // 嵌入
    }
    ```
### C. 结构体包含函数字段（策略模式）
*   **场景**：运行时动态改变对象的行为。
*   **用法**：
    
```go
    type Calculator struct {
        Operate func(int, int) int // 函数作为字段
    }
    add := Calculator{Operate: func(a, b int) int { return a + b }}
    ```
    



1. 有字段名（包含 / 嵌套，无方法提升）
结构体包含结构体
结构体包含接口
结构体包含函数
2. 无字段名（嵌入，有方法提升）
结构体嵌入结构体
结构体嵌入接口
结构体嵌入函数别名
接口嵌入接口
3. 类型实现接口（is-a 关系）
结构体实现接口
自定义函数类型实现接口（你要的）
自定义基础类型 (int/string 别名) 实现接口
4. 函数特殊形态
方法值、方法表达式


### 总结：如何选择？
1.  如果你要**存数据**，基础一定是 **Struct**。
2.  如果你要**定协议**，基础一定是 **Interface**。
3.  如果你要**传逻辑**，基础一定是 **Function**。
4.  **嵌套**是为了“偷懒”（自动获得能力），**组合**是为了“解耦”（不依赖具体实现）。