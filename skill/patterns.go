package skill

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

//////////////////////////////////////////////////
// 1) 单例 Singleton（sync.Once）
// 意图：全局仅有一个实例，惰性初始化且并发安全。
// 适用：配置中心、连接池、全局注册表等高复用的无状态或轻状态对象。
// 关键点：使用 sync.Once 保证初始化逻辑仅执行一次；避免在 init 中做复杂/可失败的初始化。
//////////////////////////////////////////////////

type Config struct {
	Endpoint string
	Timeout  time.Duration
}

var (
	cfg     *Config
	cfgOnce sync.Once
)

func GetConfig() *Config {
	cfgOnce.Do(func() {
		cfg = &Config{Endpoint: "https://api.example.com", Timeout: 3 * time.Second}
	})
	return cfg
}

//////////////////////////////////////////////////
// 2) 简单工厂 Simple Factory
// 意图：由一个工厂函数基于入参创建并返回抽象类型的不同实现。
// 适用：创建逻辑简单、分支较少，且调用方不关心具体实现类型。
// 关键点：返回接口类型；分支新增需要改动工厂（违背开闭原则，规模大时考虑工厂方法/注册表）。
//////////////////////////////////////////////////

type Shape interface {
	Draw() string
}

type circle struct{}

func (circle) Draw() string { return "draw circle" }

type rect struct{}

func (rect) Draw() string { return "draw rect" }

func NewShape(kind string) (Shape, error) {
	switch kind {
	case "circle":
		return circle{}, nil
	case "rect":
		return rect{}, nil
	default:
		return nil, fmt.Errorf("unknown shape: %s", kind)
	}
}

//////////////////////////////////////////////////
// 3) 工厂方法 Factory Method
// 意图：将“如何创建产品”的决策下放到子类/具体工厂，以便扩展新产品时无需改动调用方。
// 适用：同一产品族的多实现需要可插拔扩展；希望遵循开闭原则。
// 关键点：面向接口编程；新增实现时新增对应工厂类型，不修改现有代码。
//////////////////////////////////////////////////

type Notifier interface {
	Send(msg string) string
}

type NotifierFactory interface {
	New() Notifier
}

type Email struct{}

func (Email) Send(msg string) string { return "email: " + msg }

type EmailFactory struct{}

func (EmailFactory) New() Notifier { return Email{} }

type SMS struct{}

func (SMS) Send(msg string) string { return "sms: " + msg }

type SMSFactory struct{}

func (SMSFactory) New() Notifier { return SMS{} }

//////////////////////////////////////////////////
// 4) 抽象工厂 Abstract Factory
// 意图：为一组相关或相互依赖的产品提供统一的创建接口，且无需指定具体类型。
// 适用：需要创建“产品族”（如同一主题的 Button/Checkbox）；保证同族产品搭配使用的一致性。
// 关键点：工厂返回多个相关接口；新增产品族通过新增工厂实现，避免调用方分支判断。
//////////////////////////////////////////////////

type Button interface{ Render() string }
type Checkbox interface{ Check() string }

type WinButton struct{}

func (WinButton) Render() string { return "render win button" }

type WinCheckbox struct{}

func (WinCheckbox) Check() string { return "check win checkbox" }

type MacButton struct{}

func (MacButton) Render() string { return "render mac button" }

type MacCheckbox struct{}

func (MacCheckbox) Check() string { return "check mac checkbox" }

type UIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

type WinFactory struct{}

func (WinFactory) CreateButton() Button     { return WinButton{} }
func (WinFactory) CreateCheckbox() Checkbox { return WinCheckbox{} }

type MacFactory struct{}

func (MacFactory) CreateButton() Button     { return MacButton{} }
func (MacFactory) CreateCheckbox() Checkbox { return MacCheckbox{} }

//////////////////////////////////////////////////
// 5) 建造者 Builder
// 意图：分步骤构造复杂对象，屏蔽构造细节，提升可读性与可维护性。
// 适用：可选参数多、字段间存在默认值/依赖关系的构造场景。
// 关键点：方法链返回构建器自身；最终提供 Build 输出不可变快照（或副本）。
//////////////////////////////////////////////////

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

type RequestBuilder struct {
	r Request
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{r: Request{Method: "GET", Headers: map[string]string{}}}
}
func (b *RequestBuilder) Method(m string) *RequestBuilder { b.r.Method = m; return b }
func (b *RequestBuilder) URL(u string) *RequestBuilder    { b.r.URL = u; return b }
func (b *RequestBuilder) Header(k, v string) *RequestBuilder {
	b.r.Headers[k] = v
	return b
}
func (b *RequestBuilder) Body(p []byte) *RequestBuilder { b.r.Body = p; return b }
func (b *RequestBuilder) Build() Request                { return b.r }

//////////////////////////////////////////////////
// 6) 原型 Prototype（Clone）
// 意图：通过复制原型来创建对象，避免昂贵/复杂的构造。
// 适用：对象创建成本高、或需要复制包含复杂内部结构的实例。
// 关键点：注意深浅拷贝边界；对切片/映射/指针字段做必要的深拷贝以避免共享可变状态。
//////////////////////////////////////////////////

type Clonable interface {
	Clone() Clonable
}

type Doc struct {
	Title string
	Tags  []string
}

func (d *Doc) Clone() Clonable {
	cp := *d
	cp.Tags = append([]string(nil), d.Tags...) // 深拷贝切片
	return &cp
}

//////////////////////////////////////////////////
// 7) 适配器 Adapter
// 意图：在不修改原有类型的前提下，使其满足新的目标接口。
// 适用：新老接口不兼容；需要复用既有实现并适配到新接口。
// 关键点：适配器持有被适配者，将调用转换为目标接口期望的形式。
//////////////////////////////////////////////////

// 目标接口
type JSONLogger interface {
	LogJSON(s string) string
}

// 被适配者（已有接口）
type PlainLogger struct{}

func (PlainLogger) Log(s string) string { return "plain: " + s }

// 适配器：将 PlainLogger 适配成 JSONLogger
type PlainToJSONAdapter struct{ p PlainLogger }

func (a PlainToJSONAdapter) LogJSON(s string) string {
	return `{"msg":"` + a.p.Log(s) + `"}`
}

//////////////////////////////////////////////////
// 8) 桥接 Bridge（抽象/实现分离）
// 意图：将抽象与实现解耦，使二者可以独立扩展。
// 适用：抽象层次与实现层次都可能独立变化（如遥控器与设备）。
// 关键点：抽象持有实现接口；新增抽象或实现互不影响，减少组合爆炸。
//////////////////////////////////////////////////

type Device interface {
	On() string
	Off() string
	SetChannel(int) string
}

type TV struct{}

func (TV) On() string              { return "tv on" }
func (TV) Off() string             { return "tv off" }
func (TV) SetChannel(c int) string { return fmt.Sprintf("tv ch %d", c) }

type Remote struct{ dev Device }

func NewRemote(d Device) *Remote  { return &Remote{dev: d} }
func (r *Remote) On() string      { return r.dev.On() }
func (r *Remote) Off() string     { return r.dev.Off() }
func (r *Remote) Ch(c int) string { return r.dev.SetChannel(c) }

//////////////////////////////////////////////////
// 9) 装饰器 Decorator
// 意图：在不改变原对象的前提下，按层次为其动态地添加职责（如日志、缓存、指标）。
// 适用：横切关注点（日志/重试/熔断/监控）与业务解耦。
// 关键点：保持与被装饰者相同接口；包装调用前后扩展行为；可多层叠加组合。
//////////////////////////////////////////////////

type Service interface {
	Do(ctx context.Context, in string) (string, error)
}

type baseService struct{}

func (baseService) Do(_ context.Context, in string) (string, error) { return "base:" + in, nil }

// 装饰器：增加日志
type logDecorator struct{ next Service }

func (d logDecorator) Do(ctx context.Context, in string) (string, error) {
	// 省略实际日志
	out, err := d.next.Do(ctx, in)
	return "log->" + out, err
}

// 装饰器：增加缓存（简化）
type cacheDecorator struct {
	next Service
	mu   sync.Mutex
	data map[string]string
}

func (d *cacheDecorator) Do(ctx context.Context, in string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.data == nil {
		d.data = map[string]string{}
	}
	if v, ok := d.data[in]; ok {
		return "cache->" + v, nil
	}
	v, err := d.next.Do(ctx, in)
	if err == nil {
		d.data[in] = v
	}
	return v, err
}

//////////////////////////////////////////////////
// 10) 代理 Proxy（权限/延迟/远程）
// 意图：以代理对象控制对真实对象的访问，可用于鉴权、缓存、懒加载、远程访问等。
// 适用：需要在访问前后增加控制逻辑，且不希望侵入真实对象。
// 关键点：代理实现与被代理者相同接口；在代理中做权限校验、熔断、限流等。
//////////////////////////////////////////////////

type Repo interface {
	Get(id int) (string, error)
}

type realRepo struct{}

func (realRepo) Get(id int) (string, error) { return fmt.Sprintf("row:%d", id), nil }

// 代理：鉴权后访问真实 Repo
type authRepo struct {
	next Repo
	role string
}

func (a authRepo) Get(id int) (string, error) {
	if a.role != "admin" {
		return "", errors.New("forbidden")
	}
	return a.next.Get(id)
}

//////////////////////////////////////////////////
// 11) 外观 Facade
// 意图：为复杂子系统提供一个统一的高层接口，简化使用。
// 适用：屏蔽多组件协作的复杂性，向外暴露简洁 API。
// 关键点：外观内部编排多个子系统；不阻碍直接使用底层组件（按需绕过）。
//////////////////////////////////////////////////

type SubA struct{}

func (SubA) Step() string { return "A" }

type SubB struct{}

func (SubB) Step() string { return "B" }

type Facade struct {
	a SubA
	b SubB
}

func (f Facade) Do() string { return f.a.Step() + "->" + f.b.Step() }

//////////////////////////////////////////////////
// 12) 组合 Composite
// 意图：将对象组合成树形结构，使客户端对单个对象和组合对象的使用具有一致性。
// 适用：层级结构，如目录树、组织架构、UI 组件树。
// 关键点：统一接口（叶子与分支均实现）；在组合节点中递归地聚合行为或属性。
//////////////////////////////////////////////////

type Node interface {
	Name() string
	Total() int // 子树计数
}

type Leaf struct{ n string }

func (l Leaf) Name() string { return l.n }
func (l Leaf) Total() int   { return 1 }

type Branch struct {
	n        string
	children []Node
}

func (b Branch) Name() string { return b.n }
func (b Branch) Total() int {
	sum := 1
	for _, c := range b.children {
		sum += c.Total()
	}
	return sum
}

//////////////////////////////////////////////////
// 13) 享元 Flyweight（简化：对象池）
// 意图：共享细粒度对象，减少内存占用；此处用对象池表达可复用对象。
// 适用：对象创建/销毁成本高、数量大且可复用的场景。
// 关键点：区分内蕴状态（可共享）与外蕴状态（调用时传入）；注意并发安全。
//////////////////////////////////////////////////

type Obj struct{ kind string }

type ObjPool struct {
	mu   sync.Mutex
	free []*Obj
}

func NewObjPool() *ObjPool { return &ObjPool{} }

func (p *ObjPool) Get(kind string) *Obj {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := len(p.free)
	if n > 0 {
		o := p.free[n-1]
		p.free = p.free[:n-1]
		o.kind = kind
		return o
	}
	return &Obj{kind: kind}
}
func (p *ObjPool) Put(o *Obj) {
	p.mu.Lock()
	p.free = append(p.free, o)
	p.mu.Unlock()
}

//////////////////////////////////////////////////
// 14) 策略 Strategy（也可用函数类型）
// 意图：将一类算法封装为可互换的策略，运行时自由切换。
// 适用：排序/路由/打分等可替换算法；避免大量条件分支。
// 关键点：上下文持有策略接口；新增策略不影响调用方；在 Go 中可直接用函数类型简化实现。
//////////////////////////////////////////////////

type SortStrategy interface {
	Sort([]int) []int
}

type AscSort struct{}

func (AscSort) Sort(xs []int) []int {
	ys := append([]int(nil), xs...) // 省略排序
	return ys
}

type DescSort struct{}

func (DescSort) Sort(xs []int) []int {
	ys := append([]int(nil), xs...) // 省略排序
	return ys
}

type Sorter struct{ s SortStrategy }

func (s *Sorter) SetStrategy(ss SortStrategy) { s.s = ss }
func (s *Sorter) Do(xs []int) []int           { return s.s.Sort(xs) }

//////////////////////////////////////////////////
// 15) 观察者 Observer
// 意图：定义对象间的一对多依赖，当被观察者状态变化时通知所有观察者。
// 适用：事件分发、订阅通知、UI 事件、领域事件。
// 关键点：订阅/取消订阅并发安全；通知时通常复制快照避免长时间持锁。
//////////////////////////////////////////////////

type Event struct{ Data string }

type Observer interface{ On(Event) }
type Subject struct {
	mu sync.Mutex
	os []Observer
}

func (s *Subject) Subscribe(o Observer) {
	s.mu.Lock()
	s.os = append(s.os, o)
	s.mu.Unlock()
}
func (s *Subject) Notify(e Event) {
	s.mu.Lock()
	os := append([]Observer(nil), s.os...)
	s.mu.Unlock()
	for _, o := range os {
		o.On(e)
	}
}

type PrintObserver struct{ id string }

func (p PrintObserver) On(e Event) { _ = p.id /* fmt.Println(p.id, e.Data) */ }

//////////////////////////////////////////////////
// 16) 命令 Command
// 意图：将请求封装为对象，以便参数化、排队、记录和回放。
// 适用：任务队列、操作日志、撤销/重做。
// 关键点：命令对象实现统一接口；调用者（Invoker）仅负责调度队列。
//////////////////////////////////////////////////

type Command interface{ Exec() string }

type PrintCmd struct{ s string }

func (p PrintCmd) Exec() string { return "print:" + p.s }

type Invoker struct {
	queue []Command
}

func (i *Invoker) Add(c Command) { i.queue = append(i.queue, c) }
func (i *Invoker) Run() []string {
	var out []string
	for _, c := range i.queue {
		out = append(out, c.Exec())
	}
	return out
}

//////////////////////////////////////////////////
// 17) 责任链 Chain of Responsibility
// 意图：将请求沿链传递，直到有处理者处理或到达链尾。
// 适用：校验/鉴权/路由等步骤化、可组合的处理流程。
// 关键点：每个处理者决定处理或传递给下一个；链的顺序影响结果。
//////////////////////////////////////////////////

type Handler interface {
	SetNext(Handler) Handler
	Handle(int) string
}

type baseHandler struct{ next Handler }

func (b *baseHandler) SetNext(h Handler) Handler { b.next = h; return h }

type EvenHandler struct{ baseHandler }

func (h *EvenHandler) Handle(v int) string {
	if v%2 == 0 {
		return "even"
	}
	if h.next != nil {
		return h.next.Handle(v)
	}
	return "none"
}

type PositiveHandler struct{ baseHandler }

func (h *PositiveHandler) Handle(v int) string {
	if v > 0 {
		return "positive"
	}
	if h.next != nil {
		return h.next.Handle(v)
	}
	return "none"
}

//////////////////////////////////////////////////
// 18) 状态 State
// 意图：将对象在不同状态下的行为封装到独立状态对象，使状态切换更清晰。
// 适用：状态驱动的流程机（订单、任务、连接生命周期）。
// 关键点：上下文持有当前状态；由状态对象决定下一个状态与行为。
//////////////////////////////////////////////////

type State interface {
	Next(*ContextState) string
}
type ContextState struct{ s State }

func (c *ContextState) Set(s State) { c.s = s }
func (c *ContextState) Do() string  { return c.s.Next(c) }

type Idle struct{}

func (Idle) Next(c *ContextState) string { c.Set(Running{}); return "idle->running" }

type Running struct{}

func (Running) Next(c *ContextState) string { c.Set(Stopped{}); return "running->stopped" }

type Stopped struct{}

func (Stopped) Next(c *ContextState) string { return "stopped" }

//////////////////////////////////////////////////
// 19) 模板方法 Template Method
// 意图：在基类中定义算法骨架，将若干步骤延迟到子类实现。
// 适用：流程相同但部分步骤不同；复用不可变的算法结构。
// 关键点：通过组合嵌入基类的 Run，子类仅实现差异步骤。
//////////////////////////////////////////////////

type Template interface {
	step1() string
	step2() string
}

type BaseTemplate struct{}

func (BaseTemplate) Run(t Template) string { return t.step1() + "|" + t.step2() }

type Impl struct{ BaseTemplate }

func (Impl) step1() string { return "a" }
func (Impl) step2() string { return "b" }

//////////////////////////////////////////////////
// 20) 迭代器 Iterator（简化）
// 意图：顺序访问聚合对象的元素而不暴露其内部表示。
// 适用：自定义遍历逻辑、惰性遍历等。
// 关键点：提供 HasNext/Next；注意并发访问与修改时的语义。
//////////////////////////////////////////////////

type IntIterator struct {
	data []int
	i    int
}

func NewIntIterator(xs []int) *IntIterator { return &IntIterator{data: xs} }
func (it *IntIterator) HasNext() bool      { return it.i < len(it.data) }
func (it *IntIterator) Next() (int, bool) {
	if !it.HasNext() {
		return 0, false
	}
	v := it.data[it.i]
	it.i++
	return v, true
}

//////////////////////////////////////////////////
// 21) 中介者 Mediator（简化）
// 意图：用中介对象封装一组对象的交互，使各对象不需要显式相互引用。
// 适用：网状交互复杂度高；希望解耦同事对象之间的依赖。
// 关键点：集中路由转发；防止对象间形成强耦合。
//////////////////////////////////////////////////

type Mediator interface {
	Send(from, to string, msg string)
}
type ChatRoom struct {
	mu  sync.Mutex
	buf []string
}

func (c *ChatRoom) Send(from, to, msg string) {
	c.mu.Lock()
	c.buf = append(c.buf, fmt.Sprintf("[%s->%s]%s", from, to, msg))
	c.mu.Unlock()
}

//////////////////////////////////////////////////
// 22) 备忘录 Memento（快照）
// 意图：在不破坏封装的前提下，捕获并保存对象内部状态，以便日后恢复。
// 适用：撤销/恢复、版本切换、临时试验变更。
// 关键点：快照为值语义；注意敏感数据与存储成本。
//////////////////////////////////////////////////

type Editor struct{ content string }
type Snapshot struct{ content string }

func (e *Editor) Type(s string)      { e.content += s }
func (e *Editor) Save() Snapshot     { return Snapshot{content: e.content} }
func (e *Editor) Restore(s Snapshot) { e.content = s.content }
func (e *Editor) Content() string    { return e.content }
