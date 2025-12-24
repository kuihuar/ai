观察者模式（Observer Pattern）
作用：定义对象间的一对多依赖关系，当一个对象状态改变时，所有依赖者自动收到通知。

Go 实现：利用 Channel 或回调函数实现事件驱动。


```go
// Event 是一个结构体，用于表示主题发生变化时传递给观察者的事件信息。这里的 Data 字段是一个字符串，用于存储事件的具体数据
type Event struct {
    Data string
}

// Observer 是一个接口，它定义了一个 Notify 方法。所有实现该接口的类型都可以作为观察者，当主题发生变化并发出事件时，主题会调用观察者的 Notify 方法，将事件信息传递给观察者。
type Observer interface {
    Notify(event Event)
}

// Subject 是主题结构体，它包含一个 observers 切片，用于存储所有注册到该主题的观察者。
type Subject struct {
    observers []Observer
}
//Register 方法：该方法用于将一个观察者注册到主题中。它将传入的观察者添加到 observers 切片中，这样当主题发出事件时，这个观察者就会收到通知。
//Emit 方法：该方法用于发出事件。它遍历 observers 切片，对每个观察者调用 Notify 方法，并将事件信息传递给它们，从而实现通知所有观察者的功能。
func (s *Subject) Register(observer Observer) {
    s.observers = append(s.observers, observer)
}
func (s *Subject) Emit(event Event) {
    for _, observer := range s.observers {
        observer.Notify(event)
    }
}
// LogObserver 是一个具体的观察者类型，它实现了 Observer 接口的 Notify 方法。在 Notify 方法中，它打印出事件的 Data 字段，用于记录事件信息
type LogObserver struct{}
func (l LogObserver) Notify(event Event) {
    fmt.Println("LogObserver:", event.Data)
}

// 使用
subject := &Subject{}
subject.Register(LogObserver{})
subject.Emit(Event{Data: "User logged in"}) // 输出: LogObserver: User logged in
```
通过观察者模式，我们可以实现主题和观察者之间的解耦。主题只负责管理观察者和发出事件，而观察者只负责处理接收到的事件。这样，当需要添加新的观察者或修改观察者的行为时，不需要修改主题的代码，提高了代码的可维护性和可扩展性。如果要添加新的观察者，只需实现 Observer 接口，并将其注册到主题中即可。