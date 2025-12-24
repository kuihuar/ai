组合模式（Composite Pattern）
作用：将对象组合成树形结构以表示“部分-整体”层次结构。
组合模式是一种结构型设计模式，它允许你将对象组合成树形结构以表示“部分 - 整体”的层次结构，使得用户对单个对象和组合对象的使用具有一致性
Go 实现：通过接口统一叶子节点和组合节点的行为。

```go
//抽象组件
// Component 是一个接口，它定义了所有组件（包括叶子节点和组合节点）都必须实现的方法 Operation()。这个接口是组合模式的核心抽象，它确保了客户端可以统一地对待单个对象（叶子节点）和组合对象（包含子节点的节点）。
type Component interface {
    Operation() string
}
// 结构体是叶子节点，它实现了 Component 接口的 Operation() 方法。叶子节点是树形结构中的最底层节点，没有子节点，代表了树形结构中的基本元素
type Leaf struct{}
func (l Leaf) Operation() string { return "Leaf" }

// 组合节点（Composite）它包含一个 children 切片，用于存储子组件（可以是叶子节点或其他组合节点）。
type Composite struct {
    children []Component
}
// 方法用于向组合节点中添加子组件
func (c *Composite) Add(child Component) {
    c.children = append(c.children, child)
}
// Operation 方法遍历所有子组件，并调用它们的 Operation 方法，将结果拼接起来，最后返回一个表示整个组合节点操作结果的字符串。
func (c *Composite) Operation() string {
    var result string
    for _, child := range c.children {
        result += child.Operation() + " "
    }
    return "Composite: [" + result + "]"
}

// 使用
composite := &Composite{}
composite.Add(Leaf{})
composite.Add(Leaf{})
fmt.Println(composite.Operation()) // 输出: Composite: [Leaf Leaf ]

```

组合模式的优点
一致性：客户端可以统一地对待单个对象和组合对象，无需关心处理的是单个对象还是对象组合，简化了客户端代码。
灵活性：可以方便地添加或删除子组件，动态地构建和修改树形结构。
可扩展性：可以很容易地添加新的组件类型，而不需要修改现有代码。

适用场景
当你需要表示对象的部分 - 整体层次结构时，例如文件系统、菜单系统、组织架构等。
当你希望客户端能够统一处理单个对象和组合对象时。
通过组合模式，你可以构建出灵活、可扩展的树形结构，并且能够以一致的方式处理这些结构中的对象。
