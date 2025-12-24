1. 事件处理系统
在事件处理系统中，当某个事件发生时，需要通知多个不同的组件进行相应的处理。例如，在一个 Web 应用程序中，当用户注册成功时，可能需要同时执行发送欢迎邮件、记录日志、更新统计信息等操作。
```go
package main

import (
    "fmt"
)

// Event 定义事件结构体
type Event struct {
    UserID string
}

// Observer 定义观察者接口
type Observer interface {
    Notify(event Event)
}

// Subject 定义主题结构体
type Subject struct {
    observers []Observer
}

// Register 注册观察者
func (s *Subject) Register(observer Observer) {
    s.observers = append(s.observers, observer)
}

// Emit 发出事件
func (s *Subject) Emit(event Event) {
    for _, observer := range s.observers {
        observer.Notify(event)
    }
}

// EmailSender 邮件发送器，实现 Observer 接口
type EmailSender struct{}

func (e EmailSender) Notify(event Event) {
    fmt.Printf("Sending welcome email to user %s\n", event.UserID)
}

// Logger 日志记录器，实现 Observer 接口
type Logger struct{}

func (l Logger) Notify(event Event) {
    fmt.Printf("Logging user registration: %s\n", event.UserID)
}

func main() {
    subject := &Subject{}
    subject.Register(EmailSender{})
    subject.Register(Logger{})

    event := Event{UserID: "123"}
    subject.Emit(event)
}
```


2. 状态监控系统

在系统状态监控场景中，当系统的某个状态发生变化时，需要通知多个监控组件进行处理。例如，监控服务器的 CPU 使用率、内存使用率等，当这些指标超过阈值时，通知管理员、记录日志、触发报警等
```go
package main

import (
    "fmt"
)

// StateChangeEvent 定义状态变化事件结构体
type StateChangeEvent struct {
    Metric string
    Value  float64
}

// Observer 定义观察者接口
type Observer interface {
    Notify(event StateChangeEvent)
}

// Subject 定义主题结构体
type Subject struct {
    observers []Observer
}

// Register 注册观察者
func (s *Subject) Register(observer Observer) {
    s.observers = append(s.observers, observer)
}

// Emit 发出事件
func (s *Subject) Emit(event StateChangeEvent) {
    for _, observer := range s.observers {
        observer.Notify(event)
    }
}

// AdminNotifier 管理员通知器，实现 Observer 接口
type AdminNotifier struct{}

func (a AdminNotifier) Notify(event StateChangeEvent) {
    fmt.Printf("Notifying admin: %s value %.2f exceeded threshold\n", event.Metric, event.Value)
}

// LogRecorder 日志记录器，实现 Observer 接口
type LogRecorder struct{}

func (l LogRecorder) Notify(event StateChangeEvent) {
    fmt.Printf("Logging state change: %s value %.2f\n", event.Metric, event.Value)
}

func main() {
    subject := &Subject{}
    subject.Register(AdminNotifier{})
    subject.Register(LogRecorder{})

    event := StateChangeEvent{Metric: "CPU Usage", Value: 90.0}
    subject.Emit(event)
}
```