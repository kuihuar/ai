## Go 设计模式速查（配套 skill/patterns.go）

> 精炼版：每个模式包含「意图 / 何时使用 / 要点 / 极简示例」。示例与 `skill/patterns.go` 中实现对应，便于对照阅读与练习。

---

### 1) 单例 Singleton
- **意图**: 全局仅一个实例，并发安全、惰性初始化
- **何时使用**: 配置、注册表、共享资源
- **要点**: `sync.Once` 保证仅初始化一次
```go
cfg := GetConfig() // 全局复用
```

### 2) 简单工厂 Simple Factory
- **意图**: 一个工厂函数按入参创建不同实现
- **何时使用**: 分支少、创建简单
- **要点**: 返回接口；新增分支需改工厂
```go
s, _ := NewShape("circle")
_ = s.Draw()
```

### 3) 工厂方法 Factory Method
- **意图**: 创建交给具体工厂，新增实现无需改调用方
- **何时使用**: 同族多实现、可插拔
- **要点**: 面向接口；新增实现=新增工厂
```go
var nf NotifierFactory = EmailFactory{}
_ = nf.New().Send("hi")
```

### 4) 抽象工厂 Abstract Factory
- **意图**: 生产同一产品族的多组件
- **何时使用**: Button/Checkbox 等主题一致的家族
- **要点**: 工厂产出多种相关产品
```go
ui := WinFactory{}
_ = ui.CreateButton().Render()
_ = ui.CreateCheckbox().Check()
```

### 5) 建造者 Builder
- **意图**: 分步骤构造复杂对象
- **何时使用**: 可选参数多、默认值多
- **要点**: 链式 API；最终 Build 快照
```go
req := NewRequestBuilder().Method("POST").URL("/x").Header("K","V").Build()
_ = req
```

### 6) 原型 Prototype
- **意图**: 通过克隆创建对象
- **何时使用**: 创建成本高或深结构复制
- **要点**: 深浅拷贝边界清晰
```go
d := &Doc{Title:"t", Tags:[]string{"a"}}
_ = d.Clone()
```

### 7) 适配器 Adapter
- **意图**: 让已有类型满足新接口
- **何时使用**: 新老接口不兼容
- **要点**: 适配器持有被适配者
```go
adapter := PlainToJSONAdapter{p: PlainLogger{}}
_ = adapter.LogJSON("hi")
```

### 8) 桥接 Bridge
- **意图**: 抽象与实现分离、可独立扩展
- **何时使用**: 抽象和实现都可能变化
- **要点**: 抽象持有实现接口
```go
r := NewRemote(TV{})
_ = r.On(); _ = r.Ch(3)
```

### 9) 装饰器 Decorator
- **意图**: 运行时分层叠加新职责
- **何时使用**: 日志/缓存/指标等横切逻辑
- **要点**: 与被装饰者接口一致，可组合
```go
svc := &cacheDecorator{next: logDecorator{next: baseService{}}}
_, _ = svc.Do(context.Background(), "job")
```

### 10) 代理 Proxy
- **意图**: 控制对真实对象的访问
- **何时使用**: 鉴权、缓存、懒加载、远程
- **要点**: 代理与目标接口一致
```go
repo := authRepo{next: realRepo{}, role: "admin"}
_, _ = repo.Get(1)
```

### 11) 外观 Facade
- **意图**: 为复杂子系统提供统一接口
- **何时使用**: 简化对多组件的调用
- **要点**: 外观内部编排子系统
```go
f := Facade{}
_ = f.Do()
```

### 12) 组合 Composite
- **意图**: 树结构下统一对待叶子与组合
- **何时使用**: 目录树、组织架构、UI 树
- **要点**: 叶子和分支实现相同接口
```go
b := Branch{n:"root", children:[]Node{Leaf{"a"}, Leaf{"b"}}}
_ = b.Total()
```

### 13) 享元 Flyweight（对象池）
- **意图**: 共享细粒度对象减少内存
- **何时使用**: 大量可复用对象
- **要点**: 区分内/外蕴状态，并发安全
```go
p := NewObjPool()
o := p.Get("A")
p.Put(o)
```

### 14) 策略 Strategy
- **意图**: 可替换的算法族
- **何时使用**: 排序、路由、打分等
- **要点**: 上下文持有策略接口
```go
s := &Sorter{}
s.SetStrategy(AscSort{})
_ = s.Do([]int{3,1,2})
```

### 15) 观察者 Observer
- **意图**: 发布-订阅通知
- **何时使用**: 事件驱动、UI/领域事件
- **要点**: 并发安全，通知用快照
```go
var subj Subject
subj.Subscribe(PrintObserver{"o1"})
subj.Notify(Event{Data:"changed"})
```

### 16) 命令 Command
- **意图**: 将请求封装为对象
- **何时使用**: 队列、日志、撤销/重做
- **要点**: 调用者仅调度命令
```go
inv := &Invoker{}
inv.Add(PrintCmd{"x"})
_ = inv.Run()
```

### 17) 责任链 Chain of Responsibility
- **意图**: 链式传递直到被处理
- **何时使用**: 校验、鉴权、路由
- **要点**: 顺序影响结果
```go
e := &EvenHandler{}
p := &PositiveHandler{}
e.SetNext(p)
_ = e.Handle(3)
```

### 18) 状态 State
- **意图**: 将状态与行为封装到状态对象
- **何时使用**: 流程机/生命周期
- **要点**: 状态决定迁移与动作
```go
c := &ContextState{}
c.Set(Idle{})
_ = c.Do() // idle->running
```

### 19) 模板方法 Template Method
- **意图**: 固化算法骨架，步骤延迟到子类
- **何时使用**: 过程相同、细节不同
- **要点**: 组合嵌入基类 Run
```go
impl := Impl{}
_ = impl.Run(impl)
```

### 20) 迭代器 Iterator
- **意图**: 不暴露内部表示的遍历
- **何时使用**: 自定义/惰性遍历
- **要点**: HasNext/Next 语义清晰
```go
it := NewIntIterator([]int{1,2,3})
for it.HasNext(){ v, _ := it.Next(); _ = v }
```

### 21) 中介者 Mediator
- **意图**: 由中介管理同事对象交互
- **何时使用**: 网状依赖复杂
- **要点**: 集中路由，解耦对象
```go
room := &ChatRoom{}
room.Send("a","b","hi")
```

### 22) 备忘录 Memento
- **意图**: 捕获并恢复对象内部状态
- **何时使用**: 撤销/恢复、版本切换
- **要点**: 快照值语义，注意敏感数据
```go
ed := &Editor{}
ed.Type("hello")
snap := ed.Save()
ed.Type(" world")
ed.Restore(snap)
```

---

### 使用建议
- 将本文件与 `skill/patterns.go` 对照阅读，结合注释与代码快速掌握要点。
- 在真实项目中优先考虑：可读性、简洁性、测试性，再选择合适模式落地。 