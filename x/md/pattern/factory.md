创建型模式
作用：封装对象的创建逻辑，通过统一接口创建不同类型对象。
Go 实现：通过函数返回接口实例，隐藏具体类型。

```go
type Database interface {
    Connect() string
}

type MySQL struct{}
func (m MySQL) Connect() string { return "MySQL connected" }

type MongoDB struct{}
func (m MongoDB) Connect() string { return "MongoDB connected" }

// 工厂函数
func NewDatabase(dbType string) Database {
    switch dbType {
    case "mysql":
        return MySQL{}
    case "mongodb":
        return MongoDB{}
    default:
        panic("unknown database type")
    }
}

// 使用
db := NewDatabase("mysql")
fmt.Println(db.Connect()) // 输出: MySQL connected
```
