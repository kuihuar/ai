装饰器模式是一种结构型设计模式，它允许向一个现有的对象添加新的功能，同时又不改变其结构。在 Go 语言中，由于没有像 Python 那样的装饰器语法糖，但可以通过接口和组合来实现类似的功能。下面通过一个简单的日志记录装饰器示例来展示 Go 语言中装饰器模式的实现。

场景描述
假设我们有一个服务接口，它提供了一个处理请求的方法。现在我们希望在处理请求前后添加日志记录功能，并且不修改原服务的代码。

```go
package main

import (
	"fmt"
	"time"
)

// Service 定义服务接口
type Service interface {
	HandleRequest() string
}

// ConcreteService 实现 Service 接口
type ConcreteService struct{}

func (cs *ConcreteService) HandleRequest() string {
	return "Request handled"
}

// LoggingDecorator 日志记录装饰器
type LoggingDecorator struct {
	service Service
}

func NewLoggingDecorator(service Service) *LoggingDecorator {
	return &LoggingDecorator{
		service: service,
	}
}

func (ld *LoggingDecorator) HandleRequest() string {
	// 记录开始时间
	start := time.Now()
	fmt.Printf("Request started at: %s\n", start.Format(time.RFC3339))

	// 调用被装饰的服务的方法
	result := ld.service.HandleRequest()

	// 记录结束时间和耗时
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Request ended at: %s, elapsed time: %s\n", end.Format(time.RFC3339), elapsed)

	return result
}

func main() {
	// 创建具体服务实例
	concreteService := &ConcreteService{}

	// 使用日志记录装饰器包装具体服务
	loggingService := NewLoggingDecorator(concreteService)

	// 调用装饰后的服务方法
	response := loggingService.HandleRequest()
	fmt.Println("Response:", response)
}
```