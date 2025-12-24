1. 什么是鸭子类型？
鸭子类型（Duck Typing）是一种编程概念，常见于动态类型语言，核心思想是，如果一个对象的行为像鸭子（拥有鸭子的方法或属性），那么我们就可以把它当做鸭子来使用，而无需关心它的实际是什么类型。来源于一句谚语，如果它走起路来像鸭子，叫起来也像鸭子，那么它就是鸭子。

2. 鸭子类型的核心特点是？
- 关注行为，而非类型。对象的类型不重要，重要的是它是否实现了所需的方法或属性
- 动态绑定。方法的调用在运行时检查，而非编译时。若对象没有对应在的方法，运行时会报错。

3. go语言和鸭子类型的关型？
鸭子类型作为动态类型的的风格，在这种风格中，一个对象有效的语义，不是由继承自特定的类或实现特定的接口，而是由它当前的方法和属性的集合决定。go做为一种静态类型语言，通过接口实现了鸭子类型，实际上是go编译器在其中做了隐藏的工作。

4. 方法和函数的区别？
- 方法是一个与特定类型关联的函数，通过接收者（Receiver）来绑定到类型上
- 函数独立的代码块，不绑定到任何类型
- 调用方式不同
- 访问权限不同，方法可以访问接收者的字段段和方法，函数只能访问传入的参数
- 方法隐式实现接口，函数无法直接实现接口
- 方法定义指针接收者以修改原值，函数老板娘显式传递指针参数
- 方法可以重名，只要接收者类型不同
- 总结： 方法是面向对象编程的核用于封装类型的行为，函数是独立的代码单元适用于通用逻辑，当逻辑与类型相关时使用方法，当逻辑独立或通用时使用函数

5. 值接收者和指针接收者的区别？
- 从操作对象看，值接收者操作结构体的副本，指针接收者操作结构体的原实例
- 从内存开销看，每次调用方法时复制结构体，大对象性能差；指针接收者仅传递指针，无拷贝性能更优
- 从修改字段的功能看，值接收者无法修改，指针接收者可以看修改
- 从接口实现看，值接收者，值类型和指针类型都可以调用；仅指针类型可以调用指针接收者
- 值接收者不能为nil;指针接收者可为nil

6. 值接收者的使用场景
- 不需要修改原实例，只读取数据或者操作副本
- 结构僶较小，拷贝开销可忽略
- 确保并发安全，每个方法调用操盘独立副本，避免态条件
7. 指针接收者使用场景
- 需要修改原实例
- 结构体较大，如大数组或嵌套结构
- 实现接口方法， 若接口方法需要修改接收者，必须使用指针接收者
- 处理nil接收者，允许方法处理空指针，需要显式检查

8. 方法调用的隐式转找
- 值类型调用指针接收者方法时，会自动取地址
- 指钍类型调用值接收者方法时，会自动解引用

9. 方法的接口实现
- 若接口方法定义为指针接收者，只有指针类型实现了该接口
- 若接口方法定义为值接收者，值类型和指针类型均实实现了该接口

10. 何时使用值接收者和指针接收者
- 使用哪种接收者，不是由该方法是否修改了调用者，而是基于该类型的本质
- 方法能修改接收者指向的值，在值的类型是大型结构体时，指针接收者在每次调用时能避免复制该值。
- 如果类型具备原始的本质，也就是成员都来自内置原始类型，那就定义值接收者，如果是内置的引用类型，比如slice,map,interface,channel，声明的时候只是创建了一个header，对于他们也是直接定义值接收者类型，这样调用的时候，是直接复制了这些类型的header，header本身就是为复制设计
- 如果类型具备非原始的本质，不能被安全的复制，这种类型总是应该被共享，那就定义指针接收者，比如 struct File,就不应该被复制，应该只有一份实体

12. iface and eface的区别
iface和eface都是描述Go中接口的底层结构体，区别在于iface描述的接口包含方法， 而eface则是不包含任何方法的空接口
```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
type itab struct {
    inter *interfacetype
    _type *_type
    hash  uint32 // copy of _type.hash. Used for type switches.
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type InterfaceType struct {
        Type
        PkgPath Name      // import path
        Methods []Imethod // sorted by hash
}


type Type struct {
    Size_       uintptr
    PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
    Hash        uint32  // hash of type; avoids computation in hash tables
    TFlag       TFlag   // extra type information flags
    Align_      uint8   // alignment of variable with this type
    FieldAlign_ uint8   // alignment of struct field with this type
    Kind_       uint8   // enumeration for C
    // function for comparing objects of this type
    // (ptr to object A, ptr to object B) -> ==?
    Equal func(unsafe.Pointer, unsafe.Pointer) bool
    // GCData stores the GC type data for the garbage collector.
    // If the KindGCProg bit is set in kind, GCData is a GC program.
    // Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
    GCData    *byte
    Str       NameOff // string form
    PtrToThis TypeOff // type for pointer to this type, may be zero
}

type Imethod struct {
	Name NameOff // name of method
	Typ  TypeOff // .(*FuncType) underneath
}
```

- 接口变量由两个指针组成，一个是类型信息，一个是值， tab对应的是类型信息，data对应的是值
- itab对应的是接口表存储了接口的类型信息和相关方法，当接口变量被赋值时，runtime需要检查是否实现了接口中的所有方法，如果实现了就会填充tab和data
- data的值可能是值也可能是指针，取决于值的大小是否超过一个字的大小
- 假设有一个接口Animal和一个结构体Dog，当将Dog实例赋值给Animal接口变量时，Go会创建itab，其中包含Animal接口的类型信息和Dog的类型信息，以及Dog实现Animal接口的方法地址。data则指向Dog实例的指针或值

#### 示例 接口赋值的过程示例,将具体类型赋值给接口

```go 
type Animal interface {
    Speak()
}

type Dog struct {
    Name string
}

func (d Dog) Speak() {
    fmt.Println(d.Name, "says woof!")
}

func main() {
    var a Animal
    d := Dog{Name: "Buddy"}
    a = d // 此处触发接口赋值
}
```
- 类型检查：编译器验证Dog是否实现了Animal接口的所有方法（此处Speak方法）。

- 创建itab：若首次将Dog赋值给Animal接口，Go会生成或查找缓存的itab：

    - inter字段指向Animal接口的类型信息。

    - _type字段指向Dog的类型元数据。

    - fun数组填充Dog.Speak的方法地址。

- 设置data：

    - 由于Dog是结构体，data指向d的拷贝（若Dog较大，则存储指针）。

#### 示例方法调用与动态分派
当通过接口调用方法时，Go通过itab.fun找到具体方法实现：

```go
a.Speak() // 实际调用流程：
// 1. 通过a.tab找到itab
// 2. 通过itab.fun[0]获取Dog.Speak的地址
// 3. 传递a.data指向的值作为接收者（即方法的接收者参数）
```
#### 类型断言的工作原理
```go
if dog, ok := a.(Dog); ok {
    dog.Speak()
}
```
- Go通过比较a.tab._type和目标类型（Dog）的类型元数据，判断是否匹配。

- 若匹配，data字段的值会被转换为Dog类型并返回。

####  实际应用场景
- 性能优化：itab缓存了类型与方法的关系，避免每次方法调用都进行类型检查。

- 多态实现：允许不同具体类型通过同一接口调用，实现运行时多态。

- 反射基础：反射机制（如reflect.TypeOf）依赖_type元数据获取类型信息。

- 内存布局示意图
复制
iface
+--------+       +------------------+
| tab    | ----> | itab             |
|        |       | inter: *Animal   |
| data   |       | _type: *Dog      |
+--------+       | fun: [Speak()]   |
                 +------------------+
                 
data指针指向Dog实例的内存：
+------------------+
| Name: "Buddy"    |
+------------------+
#### 总结
- tab（*itab）：桥梁作用，连接接口与具体类型，存储类型元数据和方法表。

- data（unsafe.Pointer）：指向具体值的内存，保证接口操作真实数据。

- 协同工作：两者共同实现接口的动态特性，使得Go的接口既灵活又高效。

6. 建议
- 默认使用指针接收者，除明确不需要修改原数据或结构体极小
- 同一类型的方法尽量统一使用值或指针接收者，避免混淆

编译器如何自动检查类型是否实现接口
- var _ io.Wrter= (*myWriter)(nil)

[golang interface(上)](https://www.infoq.cn/article/wDURRBz1Nv3IbIeviIJF)
[golang interface(下)](https://www.infoq.cn/article/sP3pe06aFuGut2cl3Txt)



类型断言是对接口变量进行操作

类型转换


如果接收者类型为指针类型，只有通过指钍类型调用，
如果是结构体类型，则可以通过结体体和指针调用，直接调用指针仅仅是语法糖
接口转换的原理

接口类型interface,实体类型_type
如果类型的方法集完全包含接口的方法集，则认为实现了该接口

