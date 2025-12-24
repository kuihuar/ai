责任链模式（Chain of Responsibility）
作用：将请求沿着处理链传递，直到有对象处理它。

Go 实现：常用于中间件（如 HTTP 处理器链）。
责任链模式是一种行为设计模式，它允许你将请求沿着处理者链进行传递，直到有一个处理者能够处理该请求为止
在 Go 中实现责任链模式时，通常会定义一个接口来表示处理者的共同行为，并且每个具体的处理者都实现了这个接口

```go
package main

import (
    "fmt"
    "net/http"
)

// Handler 定义处理者接口
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// LoggerMiddleware 日志记录中间件
type LoggerMiddleware struct{}

func (l LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    fmt.Printf("Logging request: %s %s\n", r.Method, r.URL.Path)
    next(w, r)
}

// AuthMiddleware 身份验证中间件
type AuthMiddleware struct{}

func (a AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    // 简单模拟身份验证
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    next(w, r)
}

// FinalHandler 最终处理者
func FinalHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

// Chain 构建责任链
func Chain(handlers ...Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var next http.HandlerFunc = FinalHandler
        for i := len(handlers) - 1; i >= 0; i-- {
            handler := handlers[i]
            next = func(h Handler, n http.HandlerFunc) http.HandlerFunc {
                return func(w http.ResponseWriter, r *http.Request) {
                    h.ServeHTTP(w, r, n)
                }
            }(handler, next)
        }
        next(w, r)
    }
}

func main() {
    logger := LoggerMiddleware{}
    auth := AuthMiddleware{}

    http.HandleFunc("/", Chain(logger, auth))
    http.ListenAndServe(":8080", nil)
}
```

2. 工作流审批系统

在工作流审批系统中，一个申请可能需要经过多个审批环节，每个审批环节都有不同的负责人，只有当前环节审批通过后，申请才能进入下一个环节。
```go
package main

import (
    "fmt"
)

// Request 定义请求结构体
type Request struct {
    Amount float64
}

// Approver 定义审批者接口
type Approver interface {
    SetNext(approver Approver)
    Approve(request Request)
}

// Manager 经理审批者
type Manager struct {
    next Approver
}

func (m *Manager) SetNext(approver Approver) {
    m.next = approver
}

func (m *Manager) Approve(request Request) {
    if request.Amount <= 1000 {
        fmt.Println("Manager approved the request.")
    } else if m.next != nil {
        m.next.Approve(request)
    } else {
        fmt.Println("Request cannot be approved.")
    }
}

// Director 总监审批者
type Director struct {
    next Approver
}

func (d *Director) SetNext(approver Approver) {
    d.next = approver
}

func (d *Director) Approve(request Request) {
    if request.Amount <= 5000 {
        fmt.Println("Director approved the request.")
    } else if d.next != nil {
        d.next.Approve(request)
    } else {
        fmt.Println("Request cannot be approved.")
    }
}

// CEO 首席执行官审批者
type CEO struct {
    next Approver
}

func (c *CEO) SetNext(approver Approver) {
    c.next = approver
}

func (c *CEO) Approve(request Request) {
    if request.Amount > 0 {
        fmt.Println("CEO approved the request.")
    } else {
        fmt.Println("Request cannot be approved.")
    }
}

func main() {
    manager := &Manager{}
    director := &Director{}
    ceo := &CEO{}

    manager.SetNext(director)
    director.SetNext(ceo)

    request := Request{Amount: 3000}
    manager.Approve(request)
}
```
在这个例子中，Manager、Director 和 CEO 是审批者，它们构成了一个责任链。当有申请时，从 Manager 开始审批，如果 Manager 无法处理，则传递给下一个审批者。

3. 输入验证系统
在表单提交等场景中，输入数据可能需要经过多个验证步骤，如格式验证、长度验证、唯一性验证等。每个验证步骤都可以作为一个处理者，依次对输入数据进行验证。
```go
package main

import (
    "fmt"
    "strings"
)

// InputData 定义输入数据结构体
type InputData struct {
    Value string
}

// Validator 定义验证者接口
type Validator interface {
    SetNext(validator Validator)
    Validate(data InputData) bool
}

// FormatValidator 格式验证者
type FormatValidator struct {
    next Validator
}

func (f *FormatValidator) SetNext(validator Validator) {
    f.next = validator
}

func (f *FormatValidator) Validate(data InputData) bool {
    if !strings.Contains(data.Value, "@") {
        fmt.Println("Invalid format: missing @ symbol.")
        return false
    }
    if f.next != nil {
        return f.next.Validate(data)
    }
    return true
}

// LengthValidator 长度验证者
type LengthValidator struct {
    next Validator
}

func (l *LengthValidator) SetNext(validator Validator) {
    l.next = validator
}

func (l *LengthValidator) Validate(data InputData) bool {
    if len(data.Value) < 5 {
        fmt.Println("Invalid length: too short.")
        return false
    }
    if l.next != nil {
        return l.next.Validate(data)
    }
    return true
}

func main() {
    formatValidator := &FormatValidator{}
    lengthValidator := &LengthValidator{}

    formatValidator.SetNext(lengthValidator)

    input := InputData{Value: "test"}
    if formatValidator.Validate(input) {
        fmt.Println("Input data is valid.")
    }
}
```
在这个例子中，FormatValidator 和 LengthValidator 是验证者，它们构成了一个责任链。输入数据先经过 FormatValidator 验证，如果通过则传递给 LengthValidator 继续验证。

责任链模式的优点是可以灵活地组织处理流程，增加或修改处理者不会影响其他处理者。但它也可能导致调试困难，因为请求的处理路径可能不那么直观。