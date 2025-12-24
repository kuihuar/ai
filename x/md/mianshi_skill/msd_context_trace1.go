package mianshiskill

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type TraceInfo struct {
	TraceId      string
	SpanId       string
	ParentSpanId string
	ServiceName  string
}

type TraceIdKey string

type SpanIdKey string
type ParentSpanIdKey string
type contextKey string

const (
	TraceId      TraceIdKey      = "traceId"
	SpanId       SpanIdKey       = "spanId"
	ParentSpanId ParentSpanIdKey = "parentSpanId"
)

func GenerateTraceID() string {
	return uuid.New().String()
}

// GenerateSpanID 生成唯一的 Span ID
func GenerateSpanID() string {
	return uuid.New().String()
}

const (
	TraceKey contextKey = "trace"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 生成 Trace ID 和 Span ID
		traceId := GenerateTraceID()
		spanId := GenerateSpanID()
		// 将 Trace ID 和 Span ID 存储到上下文中
		ctx := context.WithValue(r.Context(), TraceKey, TraceInfo{
			TraceId: traceId,
			SpanId:  spanId,
		})
		r = r.WithContext(ctx)
		// 调用下一个处理程序
		next.ServeHTTP(w, r)
	})
}
func CallAnotherService(url string, ctx context.Context) {
	// 从上下文中获取 TraceInfo
	traceInfo, ok := ctx.Value(TraceKey).(TraceInfo)

	if !ok {
		// 如果上下文中没有 TraceInfo，创建一个新的 TraceInfo
		traceInfo = TraceInfo{
			TraceId:      GenerateTraceID(),
			SpanId:       GenerateSpanID(),
			ParentSpanId: "",
		}
	} else {
		traceInfo.ParentSpanId = traceInfo.SpanId
		traceInfo.SpanId = GenerateSpanID()
	}
	// 创建一个新的请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// 处理错误
		return
	}
	req.Header.Set("X-Trace-Id", traceInfo.TraceId)
	req.Header.Set("X-Span-Id", traceInfo.SpanId)
	req.Header.Set("X-Parent-Span-Id", traceInfo.ParentSpanId)
	// 将 TraceInfo 存储到上下文中
	req = req.WithContext(context.WithValue(req.Context(), TraceKey, traceInfo))
	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// 处理错误
		return
	}
	defer resp.Body.Close()
}

func SomeHandler(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取 TraceInfo
	traceInfo, _ := r.Context().Value(TraceKey).(TraceInfo)
	log.Printf("TraceId: %s, SpanId: %s, ParentSpanId: %s", traceInfo.TraceId, traceInfo.SpanId, traceInfo.ParentSpanId)
	go CallAnotherService("aa", r.Context())
}

func UseSomeHander() {
	mux := http.NewServeMux()
	mux.Handle("/", TraceMiddleware(http.HandlerFunc(SomeHandler)))

	log.Printf("Starting server at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

var (
	logFile *os.File
	spans   map[string]TraceInfo
)

func initLogFile(filename string) {
	var err error
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}
func closeLogFile() {
	if err := logFile.Close(); err != nil {
		log.Fatalf("Error closing log file: %v", err)
	}
}
func AddSpanIdToLog(span TraceInfo) {
	spans[span.SpanId] = span
	log.Printf("TraceId: %s, SpanId: %s, ParentSpanId: %s, Servicename",
		span.TraceId, span.SpanId, span.ParentSpanId, span.ServiceName)

}
