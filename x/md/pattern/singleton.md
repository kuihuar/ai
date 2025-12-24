创建型模式
作用：确保一个类只有一个实例，并提供全局访问点。
Go 实现：利用 sync.Once 保证线程安全的延迟初始化。

```go
type Singleton struct{
    Foo string
}

var instance *Singleton
var once sync.Once

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            Foo: "bar",
        }
    })
    return instance
}

// 使用
s1 := GetInstance()
s2 := GetInstance()
fmt.Println(s1 == s2) // 输出: true

```