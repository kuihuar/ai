

[interface-struct](https://linxuyalun.github.io/2021/02/18/golang-struct-with-embeded-interface/)

- 在给一个接口类型的变量赋予实际的值之前，它的动态类型是不存在的

```
type Action interface {
	Jump()
	Run()
}

type Person struct {
	name string
}

func (p *Person) Jump() {
	fmt.Println("JUMP")
}

func (p *Person) Run() {
	fmt.Println("RUN")
}

func main() {
	var person Action = &Person{name: "Keroro"}
	person.Jump()
}
```
对于一个接口类型的变量来说，例如上面的变量 person，赋给它的值可以被叫做它的实际值（也称动态值），而该值的类型可以被叫做这个变量的实际类型（也称动态类型）。

比如，把取址表达式 &Person 的结果值赋给了变量 person，值的类型 *Person 就是该变量的动态类型。

动态类型这个叫法是相对于静态类型而言的。对于变量 person 来讲，它的静态类型就是 Action，并且永远是 Action，但是它的动态类型却会随着赋给它的动态值而变化。比如，只有我把一个 *Person 类型的值赋给变量 person 之后，该变量的动态类型才会是 *Person 。如果还有一个 Action 接口的实现类型 *ET，并且我又把一个此类型的值赋给了 person，那么它的动态类型就会变为 *ET。

- 一个实现了某个 interface 的 struct 需要实现该 interface 下的所有方法
```
type Action interface {
	Jump()
	Run()
}

type Person struct {
	name string
}

func (p *Person) Jump() {
	fmt.Println("JUMP")
}

// 编译失败：missing method Run
func main() {
  var p Action = &Person{"Keroro"}
}
```

- 这个合法的是啥意思
```
type Action interface {
	Jump()
	Run()
}

type Person struct {
	Action
}

func main() {
	var person Action = &Person{}
	// 打印结果：person:size=[16],value=[&{<nil>}]
	fmt.Printf("person:size=[%d],value=[%v]\n", unsafe.Sizeof(person), person)
}
// 该结构体就隐式地拥有了该接口的所有方法,但需要注意的是，此时 Person 类型本身并没有实现 Action 接口的具体方法，它只是“继承”了接口的方法签名
// 但需要注意的是，此时 Person 类型本身并没有实现 Action 接口的具体方法，它只是“继承”了接口的方法签名.
func (p Person) Jump() {
  p.Action.Jump()
}
// 编译器自动为类型 *Person 实现 Jump() 方法
// Run 同理
```


- 一个 interface 占据的大小是 16 个字节，其中前 8 字节用来存储指向具体数据的指针，后 8 字节用来存储指向具体类型的指针。

Use Case
回到最开始的问题，什么时候会用上它们？

Example: Interface Wrapper
从一个角度来说，它可以做到 DRY(Don’t Repeat Yourself) 原则，另一方面，它同样可以帮助开发者遵守 SOLID() 中的依赖倒置原则。
DRY: 该原则旨在强调在编程过程中，避免代码重复，减少冗余，提高代码的可维护性、可扩展性和可读性
- 提高可维护性：如果代码中存在大量重复的部分，当需求发生变化时，需要在多个地方进行修改，容易出现遗漏或修改不一致的情况。遵循 DRY 原则，只需要在一处修改代码，就可以影响到所有使用该功能的地方，降低了维护成本。
- 增强可扩展性：当需要对某个功能进行扩展时，由于代码没有重复，只需要在封装的代码部分进行修改和扩展，而不会影响到其他部分的代码，提高了代码的灵活性和可扩展性。
- 提升可读性：去除重复代码后，代码结构更加简洁清晰，其他开发者更容易理解代码的功能和逻辑，提高了代码的可读性。
如下示例：

SOLID:
SOLID 是面向对象编程和设计中的五条重要原则，由美国软件工程师罗伯特·C·马丁（Robert C. Martin）提出，这五条原则首字母组合起来恰好是 “SOLID”，旨在帮助开发者设计出更加易维护、可扩展和灵活的软件系统。以下是这五条原则的详细介绍：

单一职责原则（Single Responsibility Principle，SRP）
单一职责原则强调一个类或模块应该只有一个引起它变化的原因。在 Go 中，我们可以将不同的功能封装到不同的结构体或函数中。

```go
package main

import "fmt"

// UserValidator 负责用户信息验证
type UserValidator struct{}

func (uv *UserValidator) ValidateUsername(username string) bool {
    return len(username) > 3
}

func (uv *UserValidator) ValidatePassword(password string) bool {
    return len(password) > 6
}

// UserRepository 负责用户信息存储
type UserRepository struct{}

func (ur *UserRepository) SaveUser(user string) {
    fmt.Printf("Saving user %s to database\n", user)
}
```
开闭原则（Open/Closed Principle，OCP）
开闭原则要求软件实体（类、模块、函数等）对扩展开放，对修改关闭。在 Go 中，我们可以通过接口和组合来实现。

```go
package main

import "fmt"

// Shape 图形接口
type Shape interface {
    Draw()
}

// Circle 圆形结构体
type Circle struct{}

func (c *Circle) Draw() {
    fmt.Println("Drawing a circle")
}

// Rectangle 矩形结构体
type Rectangle struct{}

func (r *Rectangle) Draw() {
    fmt.Println("Drawing a rectangle")
}

// DrawingProgram 绘图程序结构体
type DrawingProgram struct{}

func (dp *DrawingProgram) DrawShape(shape Shape) {
    shape.Draw()
}
```
里氏替换原则（Liskov Substitution Principle，LSP）
里氏替换原则指出所有引用基类的地方必须能透明地使用其子类的对象。在 Go 中，通过接口实现多态来体现该原则。
简单来说，子类对象能够替换父类对象，而程序的逻辑和功能不受影响。

```go
package main

import "fmt"

// Bird 鸟接口
type Bird interface {
    Move()
}

// FlyingBird 会飞的鸟结构体
type FlyingBird struct{}

func (fb *FlyingBird) Move() {
    fmt.Println("Flying")
}

// Penguin 企鹅结构体
type Penguin struct{}

func (p *Penguin) Move() {
    fmt.Println("Walking")
}
```
接口隔离原则（Interface Segregation Principle，ISP）
接口隔离原则要求客户端不应该依赖它不需要的接口。在 Go 中，将大接口拆分成多个小接口。

```go
package main

import "fmt"

// Printable 打印接口
type Printable interface {
    PrintDocument()
}

// Scannable 扫描接口
type Scannable interface {
    ScanDocument()
}

// Printer 打印机结构体
type Printer struct{}

func (p *Printer) PrintDocument() {
    fmt.Println("Printing document")
}

// Scanner 扫描仪结构体
type Scanner struct{}

func (s *Scanner) ScanDocument() {
    fmt.Println("Scanning document")
}
```
依赖倒置原则（Dependency Inversion Principle，DIP）
依赖倒置原则要求高层模块不应该依赖低层模块，两者都应该依赖抽象；抽象不应该依赖细节，细节应该依赖抽象。在 Go 中，通过接口实现依赖注入。

```go
package main

import "fmt"

// OrderRepository 订单存储接口
type OrderRepository interface {
    SaveOrder(order string)
}

// DatabaseOrderRepository 数据库订单存储结构体
type DatabaseOrderRepository struct{}

func (dor *DatabaseOrderRepository) SaveOrder(order string) {
    fmt.Printf("Saving order %s to database\n", order)
}

// OrderProcessor 订单处理结构体
type OrderProcessor struct {
    repository OrderRepository
}

func NewOrderProcessor(repository OrderRepository) *OrderProcessor {
    return &OrderProcessor{
       repository: repository,
    }
}

func (op *OrderProcessor) ProcessOrder(order string) {
    // 处理订单的逻辑
    op.repository.SaveOrder(order)
}
```
### 调用示例
```go
func main() {
    // 单一职责原则调用示例
    validator := &UserValidator{}
    fmt.Println(validator.ValidateUsername("abc"))
    fmt.Println(validator.ValidatePassword("1234567"))

    repository := &UserRepository{}
    repository.SaveUser("testuser")

    // 开闭原则调用示例
    dp := &DrawingProgram{}
    circle := &Circle{}
    dp.DrawShape(circle)
    rectangle := &Rectangle{}
    dp.DrawShape(rectangle)

    // 里氏替换原则调用示例
    flyingBird := &FlyingBird{}
    var bird Bird = flyingBird
    bird.Move()
    penguin := &Penguin{}
    bird = penguin
    bird.Move()

    // 接口隔离原则调用示例
    printer := &Printer{}
    printer.PrintDocument()
    scanner := &Scanner{}
    scanner.ScanDocument()

    // 依赖倒置原则调用示例
    dbRepo := &DatabaseOrderRepository{}
    orderProcessor := NewOrderProcessor(dbRepo)
    orderProcessor.ProcessOrder("order123")
}
```
这些示例展示了如何在 Go 语言中应用 SOLID 原则，帮助你编写更具可维护性、可扩展性和灵活性的代码。

### 开闭原则详细解说
开闭原则（Open/Closed Principle，OCP）是面向对象设计中的一个重要原则，由 Bertrand Meyer 提出，该原则强调软件实体（如类、模块、函数等）应该对扩展开放，对修改关闭。下面来详细解释这个原则，并通过具体的例子说明。

原则理解
对扩展开放：意味着软件实体应该具备良好的扩展性，当有新的需求出现时，可以通过新增代码的方式来满足需求，而不是去修改现有的核心逻辑。这样可以在不影响现有系统稳定性的前提下，灵活地添加新功能。
对修改关闭：在设计良好的系统中，一旦某个模块或类实现并测试通过后，就不应该轻易地去修改它的内部代码。因为修改现有代码可能会引入新的错误，影响到其他依赖该模块的部分。
举例说明
假设我们正在开发一个简单的图形绘制程序，最初只支持绘制圆形，后续可能需要扩展支持更多的图形，如矩形、三角形等。

不遵循开闭原则的实现
```go
package main

import "fmt"

// 图形类型的枚举
const (
    CIRCLE = iota
    RECTANGLE
)

// Shape 结构体，包含图形的基本信息
type Shape struct {
    shapeType int
    radius    float64
    width     float64
    height    float64
}

// DrawShape 函数，根据不同的图形类型进行绘制
func DrawShape(shape Shape) {
    switch shape.shapeType {
    case CIRCLE:
        fmt.Printf("Drawing a circle with radius %.2f\n", shape.radius)
    case RECTANGLE:
        fmt.Printf("Drawing a rectangle with width %.2f and height %.2f\n", shape.width, shape.height)
    default:
        fmt.Println("Unsupported shape type")
    }
}

func main() {
    circle := Shape{shapeType: CIRCLE, radius: 5.0}
    DrawShape(circle)

    rectangle := Shape{shapeType: RECTANGLE, width: 4.0, height: 6.0}
    DrawShape(rectangle)
}
```
问题分析：当我们需要新增一种图形（例如三角形）时，就需要修改 DrawShape 函数，添加新的 case 分支。这违反了开闭原则，因为每次新增图形都要修改现有的代码，增加了出错的风险，并且可能影响到其他依赖 DrawShape 函数的部分。

遵循开闭原则的实现
```go
package main

import "fmt"

// Shape 接口，定义了绘制图形的方法
type Shape interface {
    Draw()
}

// Circle 结构体，表示圆形
type Circle struct {
    radius float64
}

// Draw 方法，实现了 Shape 接口，用于绘制圆形
func (c Circle) Draw() {
    fmt.Printf("Drawing a circle with radius %.2f\n", c.radius)
}

// Rectangle 结构体，表示矩形
type Rectangle struct {
    width  float64
    height float64
}

// Draw 方法，实现了 Shape 接口，用于绘制矩形
func (r Rectangle) Draw() {
    fmt.Printf("Drawing a rectangle with width %.2f and height %.2f\n", r.width, r.height)
}

// DrawShape 函数，接受一个 Shape 接口类型的参数，调用其 Draw 方法进行绘制
func DrawShape(shape Shape) {
    shape.Draw()
}

func main() {
    circle := Circle{radius: 5.0}
    DrawShape(circle)

    rectangle := Rectangle{width: 4.0, height: 6.0}
    DrawShape(rectangle)
}
```
优点分析：在这个实现中，我们定义了一个 Shape 接口，所有具体的图形类（如 Circle 和 Rectangle）都实现了这个接口。当我们需要新增一种图形（例如三角形）时，只需要创建一个新的结构体并实现 Shape 接口的 Draw 方法，而不需要修改 DrawShape 函数。这样就实现了对扩展开放，对修改关闭的原则，提高了代码的可维护性和可扩展性。

```go
// Triangle 结构体，表示三角形
type Triangle struct {
    base   float64
    height float64
}

// Draw 方法，实现了 Shape 接口，用于绘制三角形
func (t Triangle) Draw() {
    fmt.Printf("Drawing a triangle with base %.2f and height %.2f\n", t.base, t.height)
}

func main() {
    circle := Circle{radius: 5.0}
    DrawShape(circle)

    rectangle := Rectangle{width: 4.0, height: 6.0}
    DrawShape(rectangle)

    triangle := Triangle{base: 3.0, height: 4.0}
    DrawShape(triangle)
}
```
通过上面的例子可以看到，新增三角形图形时，我们只需要添加新的结构体和实现接口方法，而不会对原有的代码造成影响，很好地遵循了开闭原则。

### 里氏替换原则详解

示例二：交通工具的继承关系
违反里氏替换原则的情况
假设有一个交通工具类 Vehicle，包含 StartEngine 和 Move 方法，还有一个子类 Bicycle 继承自 Vehicle。
```go
package main

import "fmt"

// Vehicle 交通工具结构体
type Vehicle struct{}

// StartEngine 启动引擎
func (v *Vehicle) StartEngine() {
	fmt.Println("Engine started")
}

// Move 移动
func (v *Vehicle) Move() {
	fmt.Println("Vehicle is moving")
}

// Bicycle 自行车结构体，继承自 Vehicle
type Bicycle struct {
	Vehicle
}

// StartEngine 自行车没有引擎，抛出错误
func (b *Bicycle) StartEngine() {
	fmt.Println("Error: Bicycles don't have engines")
}

// startAndMove 启动引擎并移动
func startAndMove(vehicle *Vehicle) {
	vehicle.StartEngine()
	vehicle.Move()
}

func main() {
	// 创建车辆对象
	car := &Vehicle{}
	startAndMove(car) // 正常输出

	// 创建自行车对象
	bicycle := &Bicycle{}
	startAndMove(bicycle) // 输出错误信息
}
```
问题分析：在 startAndMove 函数中，原本预期传入的交通工具对象可以正常启动引擎并移动。但当传入自行车对象时，由于自行车没有引擎，调用 StartEngine 方法会输出错误信息，导致程序行为发生变化，违反了里氏替换原则。

遵循里氏替换原则的改进
可以将 Vehicle 类拆分为有引擎的交通工具类 MotorVehicle 和无引擎的交通工具类 NonMotorVehicle。
```go
package main

import "fmt"

// Vehicle 交通工具接口
type Vehicle interface {
	Move()
}

// MotorVehicle 有引擎的交通工具结构体
type MotorVehicle struct{}

// StartEngine 启动引擎
func (m *MotorVehicle) StartEngine() {
	fmt.Println("Engine started")
}

// Move 移动
func (m *MotorVehicle) Move() {
	fmt.Println("Motor vehicle is moving")
}

// NonMotorVehicle 无引擎的交通工具结构体
type NonMotorVehicle struct{}

// Move 移动
func (n *NonMotorVehicle) Move() {
	fmt.Println("Non-motor vehicle is moving")
}

// Car 汽车结构体，继承自 MotorVehicle
type Car struct {
	MotorVehicle
}

// Bicycle 自行车结构体，继承自 NonMotorVehicle
type Bicycle struct {
	NonMotorVehicle
}

// moveVehicle 移动交通工具
func moveVehicle(vehicle Vehicle) {
	vehicle.Move()
}

func main() {
	// 创建汽车对象
	car := &Car{}
	moveVehicle(car) // 正常输出

	// 创建自行车对象
	bicycle := &Bicycle{}
	moveVehicle(bicycle) // 正常输出
}
```
优点分析：通过合理的接口设计，将不同类型的交通工具进行分类。在 moveVehicle 函数中，无论是传入有引擎的汽车对象还是无引擎的自行车对象，都能正常调用 Move 方法，程序的行为不会因为对象的替换而改变，遵循了里氏替换原则。

### 接口隔离原则
接口隔离原则（Interface Segregation Principle，ISP）是面向对象设计原则之一，它的核心思想是：客户端不应该依赖它不需要的接口，一个类对另一个类的依赖应该建立在最小的接口上。也就是说，要把大而全的接口拆分成更小、更具体的接口，这样可以让接口的使用者只需要知道他们真正需要的方法。


遵循接口隔离原则的改进
我们将大的 OfficeEquipment 接口拆分成多个小接口，每个接口只包含一个特定的功能。
```go
package main

import "fmt"

// Printable 可打印接口
type Printable interface {
	Print()
}

// Copyable 可复印接口
type Copyable interface {
	Copy()
}

// Scannable 可扫描接口
type Scannable interface {
	Scan()
}

// Faxable 可传真接口
type Faxable interface {
	Fax()
}

// Printer 打印机结构体
type Printer struct{}

// Print 实现打印功能
func (p *Printer) Print() {
	fmt.Println("Printing...")
}

// MultifunctionDevice 多功能设备结构体
type MultifunctionDevice struct{}

// Print 实现打印功能
func (m *MultifunctionDevice) Print() {
	fmt.Println("Printing...")
}

// Copy 实现复印功能
func (m *MultifunctionDevice) Copy() {
	fmt.Println("Copying...")
}

// Scan 实现扫描功能
func (m *MultifunctionDevice) Scan() {
	fmt.Println("Scanning...")
}

// Fax 实现传真功能
func (m *MultifunctionDevice) Fax() {
	fmt.Println("Faxing...")
}

func main() {
	printer := &Printer{}
	printer.Print()

	mfd := &MultifunctionDevice{}
	mfd.Print()
	mfd.Copy()
	mfd.Scan()
	mfd.Fax()
}
```
优点分析：通过将大接口拆分成多个小接口，Printer 类只需要实现它真正需要的 Printable 接口，而多功能设备 MultifunctionDevice 类可以实现多个接口。这样每个类只依赖于它真正需要的接口方法，遵循了接口隔离原则。

示例二：用户服务功能拆分
违反接口隔离原则的情况
假设有一个用户服务接口 UserService，包含用户注册、登录、修改密码和删除用户等功能。有一个具体的普通用户服务类 NormalUserService 实现这个接口。
```go
package main

import "fmt"

// UserService 用户服务接口，包含多个功能
type UserService interface {
	Register()
	Login()
	ChangePassword()
	DeleteUser()
}

// NormalUserService 普通用户服务结构体
type NormalUserService struct{}

// Register 实现用户注册功能
func (n *NormalUserService) Register() {
	fmt.Println("User registered.")
}

// Login 实现用户登录功能
func (n *NormalUserService) Login() {
	fmt.Println("User logged in.")
}

// ChangePassword 实现修改密码功能
func (n *NormalUserService) ChangePassword() {
	fmt.Println("Password changed.")
}

// DeleteUser 实现删除用户功能，普通用户可能没有权限，只是占位
func (n *NormalUserService) DeleteUser() {
	fmt.Println("Normal users cannot delete users.")
}

func main() {
	userService := &NormalUserService{}
	userService.Register()
	userService.Login()
	userService.ChangePassword()
	userService.DeleteUser()
}
```
问题分析：NormalUserService 类实现了 UserService 接口，但普通用户可能没有删除用户的权限，所以 DeleteUser 方法只是简单地占位输出提示信息，这就导致 NormalUserService 类依赖了它不需要的接口方法，违反了接口隔离原则。

遵循接口隔离原则的改进
我们将大的 UserService 接口拆分成多个小接口，每个接口只包含一个特定的功能。
```go
package main

import "fmt"

// Registerable 可注册接口
type Registerable interface {
	Register()
}

// Loginable 可登录接口
type Loginable interface {
	Login()
}

// PasswordChangable 可修改密码接口
type PasswordChangable interface {
	ChangePassword()
}

// UserDeletable 可删除用户接口
type UserDeletable interface {
	DeleteUser()
}

// NormalUserService 普通用户服务结构体
type NormalUserService struct{}

// Register 实现用户注册功能
func (n *NormalUserService) Register() {
	fmt.Println("User registered.")
}

// Login 实现用户登录功能
func (n *NormalUserService) Login() {
	fmt.Println("User logged in.")
}

// ChangePassword 实现修改密码功能
func (n *NormalUserService) ChangePassword() {
	fmt.Println("Password changed.")
}

// AdminUserService 管理员用户服务结构体
type AdminUserService struct{}

// Register 实现用户注册功能
func (a *AdminUserService) Register() {
	fmt.Println("User registered.")
}

// Login 实现用户登录功能
func (a *AdminUserService) Login() {
	fmt.Println("Admin logged in.")
}

// ChangePassword 实现修改密码功能
func (a *AdminUserService) ChangePassword() {
	fmt.Println("Admin password changed.")
}

// DeleteUser 实现删除用户功能
func (a *AdminUserService) DeleteUser() {
	fmt.Println("User deleted.")
}

func main() {
	normalUserService := &NormalUserService{}
	normalUserService.Register()
	normalUserService.Login()
	normalUserService.ChangePassword()

	adminUserService := &AdminUserService{}
	adminUserService.Register()
	adminUserService.Login()
	adminUserService.ChangePassword()
	adminUserService.DeleteUser()
}
```
优点分析：通过将大接口拆分成多个小接口，NormalUserService 类只需要实现它真正需要的接口，而管理员用户服务 AdminUserService 类可以实现所有接口。这样每个类只依赖于它真正需要的接口方法，遵循了接口隔离原则。


#### 依赖倒置原则
依赖倒置原则（Dependency Inversion Principle，DIP）是面向对象设计中的一个重要原则，它包含两个核心要点：

高层模块不应该依赖低层模块，两者都应该依赖抽象。
抽象不应该依赖细节，细节应该依赖抽象。
简单来说，依赖倒置原则强调的是依赖于抽象接口而不是具体实现类，这样可以降低模块之间的耦合度，提高代码的可维护性和可扩展性。

下面通过一个详细的 Go 语言示例来讲解依赖倒置原则。

示例场景
假设我们正在开发一个简单的电商系统，该系统中有订单服务和支付服务。订单服务在完成订单处理后需要调用支付服务进行支付操作。我们将分别展示不遵循依赖倒置原则和遵循依赖倒置原则的实现方式。

不遵循依赖倒置原则的实现
```go
package main

import "fmt"

// Alipay 是支付宝支付服务的具体实现
type Alipay struct{}

func (a *Alipay) Pay(amount float64) {
	fmt.Printf("使用支付宝支付 %.2f 元\n", amount)
}

// OrderService 是订单服务，直接依赖于 Alipay 具体实现
type OrderService struct {
	alipay *Alipay
}

func NewOrderService() *OrderService {
	return &OrderService{
		alipay: &Alipay{},
	}
}

func (os *OrderService) CreateOrder(amount float64) {
	fmt.Println("订单创建成功")
	os.alipay.Pay(amount)
}

func main() {
	orderService := NewOrderService()
	orderService.CreateOrder(100.0)
}
```
问题分析：在这个实现中，OrderService 直接依赖于 Alipay 这个具体的支付服务实现类。如果后续需要添加新的支付方式（如微信支付），就需要修改 OrderService 的代码，这违反了开闭原则，并且增加了模块之间的耦合度。

遵循依赖倒置原则的实现
```go
package main

import "fmt"

// Payment 定义支付服务的抽象接口
type Payment interface {
	Pay(amount float64)
}

// Alipay 是支付宝支付服务的具体实现
type Alipay struct{}

func (a *Alipay) Pay(amount float64) {
	fmt.Printf("使用支付宝支付 %.2f 元\n", amount)
}

// WechatPay 是微信支付服务的具体实现
type WechatPay struct{}

func (w *WechatPay) Pay(amount float64) {
	fmt.Printf("使用微信支付 %.2f 元\n", amount)
}

// OrderService 是订单服务，依赖于 Payment 抽象接口
type OrderService struct {
	payment Payment
}

func NewOrderService(payment Payment) *OrderService {
	return &OrderService{
		payment: payment,
	}
}

func (os *OrderService) CreateOrder(amount float64) {
	fmt.Println("订单创建成功")
	os.payment.Pay(amount)
}

func main() {
	// 使用支付宝支付
	alipay := &Alipay{}
	orderServiceWithAlipay := NewOrderService(alipay)
	orderServiceWithAlipay.CreateOrder(100.0)

	// 使用微信支付
	wechatPay := &WechatPay{}
	orderServiceWithWechatPay := NewOrderService(wechatPay)
	orderServiceWithWechatPay.CreateOrder(200.0)
}
```
代码解释：

定义抽象接口：首先定义了一个 Payment 接口，该接口包含一个 Pay 方法，这就是我们的抽象。
具体实现类：Alipay 和 WechatPay 结构体都实现了 Payment 接口的 Pay 方法，它们是具体的实现细节。
订单服务依赖抽象：OrderService 结构体依赖于 Payment 接口，而不是具体的支付服务实现类。在创建 OrderService 实例时，通过构造函数注入具体的支付服务实现。
优点：通过这种方式，OrderService 与具体的支付服务实现解耦。如果后续需要添加新的支付方式，只需要实现 Payment 接口，然后在创建 OrderService 实例时注入新的支付服务实现即可，无需修改 OrderService 的代码，提高了代码的可维护性和可扩展性。