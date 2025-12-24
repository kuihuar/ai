package mianshiskill

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	serviceName    = "trace-demo"
	jaegerEndpoint = "http://lcoalhost:12345/api/traces"
)

func initTracerProvider() (*sdktrace.TracerProvider, error) {
	// 创建 Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, err
	}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("environment", "dev"),
		attribute.Int64("version", 1),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),                     //批处理Span
		sdktrace.WithResource(res),                    //关联资源
		sdktrace.WithSampler(sdktrace.AlwaysSample()), //全采样
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp, nil
}
func orderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tracer := otel.Tracer("order-handler")
	ctx, span := tracer.Start(ctx, "process-order", trace.WithAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
	))
	defer span.End()

	processPayment(ctx)
	updateInventory(ctx)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order processed successfully"))

}
func processPayment(ctx context.Context) {

}
func updateInventory(ctx context.Context) {

}
func UseTrace() {
	// 初始化OpenTelemetry
	tp, err := initTracerProvider()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	http.HandleFunc("/orders", orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
