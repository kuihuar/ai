## Go 高级类型用法：interface / struct / 泛型（实用总结）

> 关注可组合性、可测试性与零成本抽象；配合简短示例便于迁移到生产代码。

---

### 1) 函数实现接口（Function as Interface）
当接口仅有一个方法时，可用函数类型适配，减少样板结构体。
```go
// 接口
type Handler interface {
	Handle(ctx context.Context, in string) (string, error)
}

// 函数类型（适配器）
type HandlerFunc func(context.Context, string) (string, error)

func (f HandlerFunc) Handle(ctx context.Context, in string) (string, error) {
	return f(ctx, in)
}

// 使用：
var h Handler = HandlerFunc(func(ctx context.Context, s string) (string, error) {
	return strings.ToUpper(s), nil
})
```
- **优点**: 简洁、易组合、便于测试与装饰。
- **装饰器模式**: 使用闭包叠加重试/限流/日志。
```go
func WithLog(next Handler) Handler {
	return HandlerFunc(func(ctx context.Context, in string) (string, error) {
		log.Println("in=", in)
		out, err := next.Handle(ctx, in)
		log.Println("out=", out, "err=", err)
		return out, err
	})
}
```
- **注意**: 若接口有多个方法，不宜用函数适配，应定义结构体并实现全部方法。

---

### 2) 方法集与接收者（Method Set; Value vs Pointer Receiver）
- 值接收者方法属于类型与其指针；指针接收者方法仅属于指针类型。
- 若方法需修改状态或避免拷贝，使用指针接收者。
```go
type Counter struct { n int }
func (c Counter) Value() int           { return c.n }
func (c *Counter) Inc()                { c.n++ }

var c Counter
_ = c.Value()   // ok
// c.Inc()      // 编译错误：方法集不含指针接收者
(&c).Inc()      // ok
```
- 在接口赋值时要注意方法集匹配：若接口需要 `Inc()`，则需要 `*Counter` 满足接口。
- 自动取址/解址：调用 `c.Inc()` 会报错，但调用者若持有 `*Counter`，可直接 `p.Inc()`。
- 接口变量不支持自动取址：`var x interface{ Inc() } = c` 不成立，需要 `&c`。

---

### 3) 组合与嵌入（Embedding）
- 结构体嵌入可“提升”字段/方法，形成组合优于继承的复用。
```go
type Logger struct{ io.Writer }
func (l Logger) Info(msg string) { fmt.Fprintln(l.Writer, "INFO:", msg) }

type Service struct{ Logger } // 嵌入

func NewService(w io.Writer) Service { return Service{Logger{w}} }

// 使用：
s := NewService(os.Stdout)
s.Info("started") // 直接访问提升的方法
```
- 嵌入接口实现接口聚合：
```go
type Closer interface{ Close() error }

type ReadCloser interface {
	io.Reader
	Closer
}
```
- **陷阱**: 多重嵌入出现同名字段/方法会产生选择歧义，需要显式限定：`s.Logger.Info(...)`。
- **建议**: 嵌入用于复用行为，不等同于继承层级；避免深层嵌入导致 API 难以理解。

---

### 4) 接口设计原则与零值友好（Nil Safety）
- 小接口优先（interface segregation）：如 `io.Reader`、`io.Writer`。
- 返回接口还是结构体？创建者返回具体类型，调用者依赖接口。
- 允许 `nil` 实现：例如日志器为 `nil` 时降级为 no-op。
```go
type LoggerI interface{ Printf(string, ...any) }

type NopLogger struct{}
func (NopLogger) Printf(string, ...any) {}

// 依赖注入
func Do(l LoggerI) { if l == nil { l = NopLogger{} }; l.Printf("ok") }
```
- **空接口 `any` 的位置**：对外 API 避免随处 `any`，应先尝试最小接口；边界层（序列化/配置）可接受 `any`。

---

### 5) 泛型基础：类型参数与约束（Constraints）
- 使用 `~` 指定底层类型，支持别名。
- 通过自定义约束表达能力边界。
```go
// 约束：可有序比较
type Ordered interface { ~int | ~int64 | ~float64 | ~string }

func Min[T Ordered](a, b T) T {
	if a < b { return a }
	return b
}

// 使用别名也可：
type MyInt int
var _ = Min(MyInt(3), MyInt(5))
```
- **类型推断**：常见函数可省略类型实参；必要时显式指定 `Min[int](...)`。
- **约束越窄越好**：避免在约束里放不需要的操作符或方法。

---

### 6) 类型集（Type Sets）与 union
- 用 `|` 组合类型，精确控制支持的类型集合；与 `~` 一起用于“底层类型相同”。
```go
type Number interface { ~int | ~int64 | ~float64 }

func Sum[T Number](xs ...T) T {
	var s T
	for _, x := range xs { s += x }
	return s
}
```
- **陷阱**：带类型集的接口不可作为普通接口值使用（它们通常仅用于约束）。

---

### 7) 泛型与接口的协作（约束含方法）
- 约束可包含方法签名；类型实参需实现这些方法。
```go
type Stringer interface{ String() string }

type ToStrings[T Stringer](xs []T) []string {
	out := make([]string, 0, len(xs))
	for _, v := range xs { out = append(out, v.String()) }
	return out
}
```
- **限制**：方法暂不支持独立类型参数（按 Go 现状，以类型级泛型组织 API）。

// 补充：更完整的可运行示例（两种类型同时实现约束方法）
```go
type User struct{ ID int; Name string }
func (u User) String() string { return fmt.Sprintf("%d:%s", u.ID, u.Name) }

type Product struct{ SKU string; Price int }
func (p Product) String() string { return fmt.Sprintf("%s:%d", p.SKU, p.Price) }

func DemoToStrings() {
	us := []User{{1,"A"},{2,"B"}}
	ps := []Product{{"X",100},{"Y",200}}
	fmt.Println(ToStrings(us)) // [1:A 2:B]
	fmt.Println(ToStrings(ps)) // [X:100 Y:200]
}
```

// 补充：方法型约束（带行为的约束），如比较接口
```go
type Less interface{ Less(other any) bool }

type Sortable[T Less] []T

func (s Sortable[T]) Sort() {
	sort.Slice(s, func(i, j int) bool { return s[i].Less(s[j]) })
}

// 实现示例
type Score struct{ Name string; V int }
func (a Score) Less(b any) bool { return a.V < b.(Score).V }

func DemoSortable() {
	arr := Sortable[Score]{{"a",3},{"b",1},{"c",2}}
	arr.Sort()
	fmt.Println(arr) // [{b 1} {c 2} {a 3}]
}
```

// 补充：以“类型级泛型”组织 API（而非每个函数都写类型参数）
// 思路：把类型参数提升到类型定义上，方法使用接收者的类型参数。
```go
// 类型级泛型的集合 API
type Bag[T any] struct{ items []T }

func NewBag[T any](xs ...T) *Bag[T] { return &Bag[T]{items: append([]T(nil), xs...)} }
func (b *Bag[T]) Add(v T)           { b.items = append(b.items, v) }
func (b *Bag[T]) All() []T          { return append([]T(nil), b.items...) }

// 方法也可以返回携带新类型参数的新容器（组织在类型上，而不是散落的函数）
func MapBag[T any, R any](b *Bag[T], f func(T) R) *Bag[R] {
	out := make([]R, len(b.items))
	for i, v := range b.items { out[i] = f(v) }
	return &Bag[R]{items: out}
}

func DemoBag() {
	b := NewBag(1,2,3)
	b.Add(4)
	fmt.Println(b.All()) // [1 2 3 4]
	b2 := MapBag(b, func(x int) string { return strconv.Itoa(x) })
	fmt.Println(b2.All()) // ["1" "2" "3" "4"]
}
```

---

### 9) 协变/逆变直觉（Go 的方式）
- Go 无显式协变/逆变，但通过“读取用接口、写入用具体”与最小接口可达成大部分需求。
- 切片与数组不变；用泛型函数抽象变型场景。
```go
func ReadAll[T any](it Iterator[T]) []T { /* ... */ return nil }
```

// 补充：可运行迭代器示例，强调“读取用接口”
```go
// 读取侧的最小接口
interface Iterator[T any] {
	Next() (T, bool)
}

// 基于切片的实现
type sliceIter[T any] struct{ i int; xs []T }
func (s *sliceIter[T]) Next() (T, bool) {
	if s.i >= len(s.xs) { var zero T; return zero, false }
	v := s.xs[s.i]
	s.i++
	return v, true
}

func FromSlice[T any](xs []T) Iterator[T] { return &sliceIter[T]{xs: xs} }

// 读取：将迭代器中的元素读完
func ReadAll[T any](it Iterator[T]) []T {
	var out []T
	for {
		v, ok := it.Next(); if !ok { break }
		out = append(out, v)
	}
	return out
}

func DemoIterator() {
	it := FromSlice([]int{1,2,3})
	fmt.Println(ReadAll(it)) // [1 2 3]
}
```

// 说明：为什么“不用协变”也能满足大部分需求？
// - 读取端用最小接口（这里只读 Next），调用方只依赖读取能力。
// - 写入端用具体类型接收者，避免复杂变型规则。

// 写入用具体：把元素写入到给定切片（或缓冲区）
```go
func CopyInto[T any](dst []T, it Iterator[T]) int {
	count := 0
	for {
		v, ok := it.Next(); if !ok { break }
		if count < len(dst) {
			dst[count] = v
			count++
		} else {
			break
		}
	}
	return count
}

func DemoCopyInto() {
	buf := make([]int, 2)
	it := FromSlice([]int{7,8,9})
	n := CopyInto(buf, it)
	fmt.Println(n, buf) // 2 [7 8]
}
```

// 类型级泛型组织流式 API（以类型承载操作，而非分散函数）
```go
type Stream[T any] struct{ it Iterator[T] }

func From[T any](xs []T) Stream[T] { return Stream[T]{it: FromSlice(xs)} }

// Map 返回新类型参数的 Stream（方法内部可使用新的类型参数）
func Map[T any, R any](s Stream[T], f func(T) R) Stream[R] {
	type mapIter struct{ src Iterator[T]; f func(T) R }
	var _ Iterator[R] = (*mapIter)(nil)
	func (m *mapIter) Next() (R, bool) {
		v, ok := m.src.Next(); if !ok { var zero R; return zero, false }
		return m.f(v), true
	}
	return Stream[R]{it: &mapIter{src: s.it, f: f}}
}

func Collect[T any](s Stream[T]) []T { return ReadAll(s.it) }

func DemoStream() {
	res := Collect(Map(From([]int{1,2,3}), func(x int) string { return strconv.Itoa(x) }))
	fmt.Println(res) // ["1" "2" "3"]
}
```

- 实践要点：
  - 读取能力抽象为极小接口，最大化适配性；
  - 写入或构建由具体类型掌控，避免协变/逆变复杂度；
  - 将类型参数提升到“容器/流”类型上，方法围绕该类型组织，API 更聚合、可组合。

---

### 10) 可选参数与 Builder（泛型 + 函数式选项）
- 通过函数式选项提升可读性与扩展性。
```go
type Server struct{ addr string; timeout time.Duration }

type Option func(*Server)

func WithAddr(a string) Option { return func(s *Server){ s.addr = a } }
func WithTimeout(d time.Duration) Option { return func(s *Server){ s.timeout = d } }

func NewServer(opts ...Option) *Server {
	s := &Server{addr: ":8080", timeout: 5*time.Second}
	for _, opt := range opts { opt(s) }
	return s
}

_ = NewServer(WithAddr(":9090"), WithTimeout(2*time.Second))
```
- **扩展**：可用泛型将 Option 复用到不同类型：`type OptionOf[T any] func(*T)`。

---

### 11) 通过接口模拟枚举行为（带方法）
- 枚举值实现统一接口，行为与数据并存，避免巨大 switch。
```go
type Op interface{ Apply(a, b int) int }

type add struct{}; func (add) Apply(a,b int) int { return a+b }

type mul struct{}; func (mul) Apply(a,b int) int { return a*b }

func Calc(op Op, a, b int) int { return op.Apply(a,b) }
```
- **扩展点**：为每个值增加额外方法或元数据，而不污染调用方。

---

### 12) 错误值与接口（多态错误）
- 通过实现 `error` 与自定义方法实现富错误；配合 `errors.As` 解包。
```go
type NotFoundError struct{ Key string }
func (e NotFoundError) Error() string { return "not found: " + e.Key }
func (e NotFoundError) NotFound() bool { return true }

// 使用
var err error = NotFoundError{"id:1"}
var nfe NotFoundError
if errors.As(err, &nfe) { /* handle */ }
```
- **建议**：错误类型要小而明确；暴露语义方法便于分支而不依赖字符串匹配。

---

### 13) 约束型工厂与接口返回（避免泄漏实现）
- 创建者返回接口，内部选择具体类型；配合泛型约束限制输入类型。
```go
type Repo interface{ Get(id int64) (string, error) }

type memRepo struct{ data map[int64]string }
func (m *memRepo) Get(id int64) (string, error) { return m.data[id], nil }

func NewRepo() Repo { return &memRepo{data: map[int64]string{1:"x"}} }
```
- **封装**：便于将来迁移到 SQL/远程服务而不影响调用方。

---

### 14) Map/Set 的泛型封装
```go
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T){ s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool { _, ok := s[v]; return ok }
func (s Set[T]) Keys() []T { ks := make([]T,0,len(s)); for k := range s { ks = append(ks, k) }; return ks }
```
- **技巧**：利用 `comparable` 约束即可直接作为 map 的键；必要时自定义哈希容器。

---

### 15) 性能与逃逸：接口与泛型的选择
- 接口有动态分发与逃逸成本；泛型在单态化后通常零成本。
- 热路径优先：能用泛型就不要接口；边界处保持接口易扩展。
- **基准**：对关键路径用 `testing.B` 做对比，避免“感觉优化”。

---

### 16) 类型断言与类型分支（Type Assertion / Type Switch）
- 将接口值还原为具体类型，注意断言失败分支。
```go
func F(x any) {
	if s, ok := x.(string); ok {
		fmt.Println("string:", s)
		return
	}
	switch v := x.(type) {
	case int:
		fmt.Println("int:", v)
	case fmt.Stringer:
		fmt.Println("stringer:", v.String())
	default:
		fmt.Println("unknown")
	}
}
```
- **建议**：公共 API 避免把 `any` 继续向上传递；尽量在边界处做解析与校验。

---

### 17) 接口值、动态类型与 nil 陷阱
- 接口值包含“动态类型”与“动态值”。当动态类型非 `nil` 但动态值为 `nil`，接口本身不为 `nil`。
```go
var e error          // e == nil
var p *os.PathError  // p == nil

e = p                // 接口 e != nil（动态类型为 *os.PathError）
fmt.Println(e == nil) // false
```
- **建议**：
  - 返回 `error` 时，若要表示无错，直接返回 `nil`，不要返回带 `nil` 值的具体错误类型。
  - 在接口字段中存放实现时，判空用“协议方法”或显式布尔标记，避免直接与 `nil` 比较。

---

### 18) 接口可比较性与 map/set 键
- 只有当动态类型可比较时，接口值才可比较/作为 map 键。
- 包含切片、map、函数的结构体不可比较；作为接口值比较会 `panic`。
- **建议**：作为键的类型设计为不可变、可比较（只含基础类型与可比较结构）。

---

### 19) 嵌入与方法覆盖的细节
- 子类型定义同名方法将“遮蔽”嵌入类型的方法。
- 指针接收者的嵌入：若想通过值接收者访问其指针方法，需要持有指针的外层类型。
```go
type A struct{}
func (A) Hello() { fmt.Println("A") }

type B struct{ *A }
func (B) Hello() { fmt.Println("B") }

b := B{A: &A{}}
b.Hello()  // 输出 B
b.A.Hello() // 显式调用嵌入者方法
```

---

### 20) 实战清单（Checklist）
- 定义最小接口；调用方依赖接口，创建方返回具体。
- 单方法接口用 `Func` 适配，拥抱函数式组合与装饰器。
- 明确方法集与接收者选择；面向指针修改状态；接口赋值考虑方法集匹配。
- 泛型优先热路径；约束表达明确边界；类型集配合底层类型；必要时显式类型实参。
- 组合/嵌入优于继承；同名遮蔽要显式分派；避免深层嵌入与接口歧义。
- 用类型断言/分支在边界处还原；避免 `any` 继续传播。
- 理解接口 `nil` 陷阱；错误返回用真正的 `nil`。
- 设计可比较的键类型；谨慎将接口值用作集合键。 