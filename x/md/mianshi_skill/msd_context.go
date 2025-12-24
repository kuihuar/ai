package mianshiskill

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// OpenTelemetry 是一个开放的、跨语言的分布式追踪系统，用于收集、存储和分析应用程序的性能数据。
// 它提供了一个统一的接口，用于收集和传播分布式追踪数据，支持多种后端存储和分析工具。
// OpenTelemetry 的目标是成为一个可扩展、可配置的分布式追踪系统，支持多种语言和框架，
// 并提供了丰富的功能和工具，用于开发者和运维人员进行性能分析和故障排查。
type traceKey struct{}
type Trace struct {
	TraceID  string
	SpanID   string
	ParentID string
}

const (
	TraceIDHeader = "X-Trace-ID"
	SpanIDHeader  = "X-Span-ID"
)

// 中间件注入traceID
// traceID 是一个用于标识整个请求链的唯一标识符。
func generateTraceID() string {
	return "1234567890"
}

// spanID 是一个用于标识请求中的一个具体操作或步骤的标识符。
func generateSpanID() string {
	return "1234567890"
}
func InjectTraceIDWithMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 优先尝试获取上游 TraceID
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = generateTraceID() // 生成唯一ID (如UUIDv4)
		}

		// 创建新 SpanID
		spanID := generateSpanID()

		// 注入 context
		ctx := context.WithValue(r.Context(), traceKey{}, Trace{
			TraceID:  traceID,
			SpanID:   spanID,
			ParentID: r.Header.Get(SpanIDHeader), // 记录父级Span
		})

		// 传播到下游（设置响应头）
		w.Header().Set(TraceIDHeader, traceID)
		w.Header().Set(SpanIDHeader, spanID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 上下文传递traceID
// HTTP 客户端调用前设置头部
// func SendRequest(ctx context.Context, url string) {
// 	trace, ok := ctx.Value(traceKey{}).(Trace)
// 	if !ok {
// 		trace = Trace{TraceID: "unknown"}
// 	}

// 	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
// 	req.Header.Set(TraceIDHeader, trace.TraceID)
// 	req.Header.Set(SpanIDHeader, generateSpanID()) // 生成新Span

// 	// 执行请求...
// }

// gRPC 元数据传递
// func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
// 	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

// 	if trace, ok := ctx.Value(traceKey{}).(Trace); ok {
// 		md, _ := metadata.FromOutgoingContext(ctx)
// 		md = md.Copy()
// 		md.Set(TraceIDHeader, trace.TraceID)
// 		md.Set(SpanIDHeader, generateSpanID())
// 		ctx = metadata.NewOutgoingContext(ctx, md)
// 	}
// 	return invoker(ctx, method, req, reply, cc, opts...)
// }

// gRPC 传递traceID

// 日志集成

// 从 context 提取追踪信息的日志器
// func LogWithTrace(ctx context.Context) *log.Entry {
// 	fields := log.Fields{}

// 	if trace, ok := ctx.Value(traceKey{}).(Trace); ok {
// 		fields["trace_id"] = trace.TraceID
// 		fields["span_id"] = trace.SpanID
// 		fields["parent_span"] = trace.ParentID
// 	}

// 	return log.WithFields(fields)
// }

// 使用示例
//
//	func HandleRequest(ctx context.Context) {
//		LogWithTrace(ctx).Info("Start processing request")
//		// ...业务逻辑
//	}
func TestContext1() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("Request timed out")
	case <-time.After(1 * time.Second):
		fmt.Println("Request completed")
	}
	fmt.Println("end")

}
func TestContext() {
	// WithTimeout 相对时间
	// WithDeadline 绝对时间
	// WithCancel 手动取消
	// WithValue 传递值
	// 所以派生自该context的子content，都会被取消
	// 当父context被取消时，所有派生的子context都会被取消·
	timeout := 4 * time.Second // 设置4秒时，dowork 不会被取消
	timeout = 3 * time.Second  // 设置3秒时，dowork 会被取消

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // 确保在函数结束时调用取消

	for i := 1; i <= 3; i++ {
		go dowork(ctx, i) // 启动子 goroutine
	}

	// 等待一段时间以确保 goroutine 有时间运行
	time.Sleep(5 * time.Second)
	fmt.Println("主程序结束")
}

func dowork(ctx context.Context, id int) {
	fmt.Printf("Worker %d 开始工作\n", id)
	select {
	case <-time.After(3 * time.Second): // 模拟长时间工作
		fmt.Printf("Worker %d 完成工作\n", id)
	case <-ctx.Done(): // 监听上下文取消
		fmt.Printf("Worker %d 被取消\n", id)
	}
}

// TraceID/SpanID
func TraceIDSpanID() {
	// 创建一个根上下文
	rootCtx := context.Background()
	// 创建一个带有超时的子上下文
	_, cancel := context.WithCancel(rootCtx)
	defer cancel()
}
