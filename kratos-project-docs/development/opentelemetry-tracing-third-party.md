# 在第三方服务调用中创建 OpenTelemetry Span

## 概述

在调用第三方服务（HTTP API、gRPC 服务等）时，创建新的 OpenTelemetry span 可以帮助我们：

- **追踪外部调用**：了解每个第三方服务调用的耗时和状态
- **分布式追踪**：将外部服务调用纳入完整的请求追踪链路
- **问题排查**：快速定位是哪个第三方服务调用出现问题
- **性能分析**：分析外部服务调用的性能瓶颈

## 基本用法

### 1. 导入必要的包

```go
import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)
```

### 2. 创建 Tracer

在客户端结构体中添加 tracer，或者在方法中获取全局 tracer：

```go
// 方式 1: 在结构体中保存 tracer（推荐，性能更好）
type Client struct {
    tracer trace.Tracer
    // ... 其他字段
}

func NewClient(cfg *Config, logger log.Logger) (*Client, error) {
    return &Client{
        tracer: otel.Tracer("dingtalk-client"), // 使用服务名作为 tracer 名称
        // ... 其他初始化
    }, nil
}

// 方式 2: 在方法中获取全局 tracer（简单，但每次调用都有开销）
func (c *Client) SomeMethod(ctx context.Context) {
    tracer := otel.Tracer("dingtalk-client")
    // ...
}
```

### 3. 创建 Span

在调用第三方服务的方法中创建 span：

```go
func (c *Client) GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
    // 创建新的 span，span 名称应该描述操作
    ctx, span := c.tracer.Start(ctx, "dingtalk.GetUserInfo",
        trace.WithAttributes(
            attribute.String("dingtalk.user_id", userID),
            attribute.String("dingtalk.operation", "get_user_info"),
        ),
    )
    defer span.End() // 确保 span 被结束

    // 执行实际的 HTTP 调用
    var resp GetUserInfoResponse
    httpResp, err := c.client.R().
        SetContext(ctx). // 重要：将包含 span 的 context 传递给 HTTP 客户端
        SetQueryParam("access_token", token).
        SetQueryParam("userid", userID).
        SetResult(&resp).
        Get("/topapi/v2/user/get")

    // 处理错误和结果
    if err != nil {
        // 记录错误到 span
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }

    if !httpResp.IsSuccess() || resp.ErrCode != 0 {
        err := fmt.Errorf("get user info failed: %s", resp.ErrMsg)
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int("dingtalk.error_code", resp.ErrCode),
            attribute.String("dingtalk.error_msg", resp.ErrMsg),
        )
        return nil, err
    }

    // 成功时设置状态和属性
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.String("dingtalk.user_name", resp.Result.Name),
        attribute.String("http.status_code", fmt.Sprintf("%d", httpResp.StatusCode())),
    )

    return resp.Result, nil
}
```

## 最佳实践

### 1. Span 命名规范

- **格式**：`{service-name}.{operation-name}`
- **示例**：
  - `dingtalk.GetUserInfo`
  - `payment.CreateOrder`
  - `user-service.GetUser`

### 2. 属性命名规范

使用统一的属性命名前缀，便于在追踪系统中过滤和查询：

- **服务相关**：`{service}.{field}`，如 `dingtalk.user_id`、`dingtalk.error_code`
- **HTTP 相关**：`http.method`、`http.url`、`http.status_code`
- **gRPC 相关**：`rpc.method`、`rpc.service`、`rpc.status_code`

### 3. 错误处理

始终记录错误信息：

```go
if err != nil {
    span.RecordError(err)           // 记录错误堆栈
    span.SetStatus(codes.Error, err.Error()) // 设置错误状态
    return nil, err
}
```

### 4. 性能指标

记录关键性能指标：

```go
start := time.Now()
// ... 执行操作
duration := time.Since(start)
span.SetAttributes(
    attribute.Int64("duration_ms", duration.Milliseconds()),
)
```

### 5. Context 传播

**重要**：确保将包含 span 的 context 传递给 HTTP/gRPC 客户端，这样：

- HTTP 客户端可以自动注入 trace headers（如果配置了）
- 子操作可以创建子 span
- 日志可以自动关联 trace 信息

```go
// ✅ 正确：传递 context
httpResp, err := c.client.R().
    SetContext(ctx). // 包含 span 的 context
    Get("/api/users")

// ❌ 错误：使用新的 context
httpResp, err := c.client.R().
    SetContext(context.Background()). // 丢失了 trace 信息
    Get("/api/users")
```

## 完整示例

### HTTP 客户端示例（钉钉服务）

```go
package dingtalk

import (
    "context"
    "fmt"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-resty/resty/v2"
)

type Client struct {
    client    *resty.Client
    baseURL   string
    appKey    string
    appSecret string
    log       *log.Helper
    tracer    trace.Tracer // 添加 tracer
    // ... 其他字段
}

func NewClient(cfg *Config, logger log.Logger) (*Client, error) {
    // ... 初始化 resty client
    
    return &Client{
        client:    client,
        baseURL:   cfg.BaseURL,
        appKey:    cfg.AppKey,
        appSecret: cfg.AppSecret,
        log:       log.NewHelper(logger),
        tracer:    otel.Tracer("dingtalk-client"), // 初始化 tracer
    }, nil
}

func (c *Client) GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
    // 创建 span
    ctx, span := c.tracer.Start(ctx, "dingtalk.GetUserInfo",
        trace.WithAttributes(
            attribute.String("dingtalk.user_id", userID),
            attribute.String("http.method", "GET"),
            attribute.String("http.url", fmt.Sprintf("%s/topapi/v2/user/get", c.baseURL)),
        ),
    )
    defer span.End()

    start := time.Now()

    // 获取 access_token
    token, err := c.getAccessToken(ctx)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "failed to get access token")
        return nil, err
    }

    // 执行 HTTP 请求
    var resp GetUserInfoResponse
    httpResp, err := c.client.R().
        SetContext(ctx). // 传递包含 span 的 context
        SetQueryParam("access_token", token).
        SetQueryParam("userid", userID).
        SetResult(&resp).
        Get("/topapi/v2/user/get")

    duration := time.Since(start)

    // 处理错误
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int64("duration_ms", duration.Milliseconds()),
        )
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }

    // 处理业务错误
    if !httpResp.IsSuccess() || resp.ErrCode != 0 {
        err := fmt.Errorf("get user info failed: %s", resp.ErrMsg)
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        span.SetAttributes(
            attribute.Int("http.status_code", httpResp.StatusCode()),
            attribute.Int("dingtalk.error_code", resp.ErrCode),
            attribute.String("dingtalk.error_msg", resp.ErrMsg),
            attribute.Int64("duration_ms", duration.Milliseconds()),
        )
        return nil, err
    }

    // 成功
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int("http.status_code", httpResp.StatusCode()),
        attribute.String("dingtalk.user_name", resp.Result.Name),
        attribute.Int64("duration_ms", duration.Milliseconds()),
    )

    return resp.Result, nil
}
```

### gRPC 客户端示例

```go
package payment

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    
    "google.golang.org/grpc"
    pb "api/external/payment/v1"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.PaymentServiceClient
    tracer trace.Tracer
}

func (c *Client) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    // 创建 span
    ctx, span := c.tracer.Start(ctx, "payment.CreateOrder",
        trace.WithAttributes(
            attribute.String("rpc.service", "payment.PaymentService"),
            attribute.String("rpc.method", "CreateOrder"),
            attribute.String("payment.user_id", req.UserId),
            attribute.Float64("payment.amount", float64(req.Amount)),
        ),
    )
    defer span.End()

    // 调用 gRPC 服务（context 会自动传播 trace 信息）
    resp, err := c.client.CreateOrder(ctx, req)

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.String("payment.order_id", resp.OrderId),
        attribute.String("rpc.status_code", "OK"),
    )

    return resp, nil
}
```

## 嵌套 Span

对于复杂的操作，可以创建嵌套的 span：

```go
func (c *Client) GetAllUsersByDeptID(ctx context.Context, deptID int64) ([]*UserInfo, error) {
    // 父 span：整个操作
    ctx, span := c.tracer.Start(ctx, "dingtalk.GetAllUsersByDeptID",
        trace.WithAttributes(
            attribute.Int64("dingtalk.dept_id", deptID),
        ),
    )
    defer span.End()

    var allUsers []*UserInfo
    cursor := int64(0)
    size := 100

    for {
        // 子 span：每次分页请求
        ctx, pageSpan := c.tracer.Start(ctx, "dingtalk.GetUserList",
            trace.WithAttributes(
                attribute.Int64("dingtalk.cursor", cursor),
                attribute.Int("dingtalk.page_size", size),
            ),
        )

        result, err := c.GetUserList(ctx, deptID, cursor, size)
        
        if err != nil {
            pageSpan.RecordError(err)
            pageSpan.SetStatus(codes.Error, err.Error())
            pageSpan.End()
            span.RecordError(err)
            return nil, err
        }

        if result != nil && result.List != nil {
            allUsers = append(allUsers, result.List...)
        }

        pageSpan.SetAttributes(
            attribute.Int("dingtalk.users_count", len(result.List)),
            attribute.Bool("dingtalk.has_more", result.HasMore),
        )
        pageSpan.SetStatus(codes.Ok, "success")
        pageSpan.End()

        if !result.HasMore {
            break
        }

        cursor = result.NextCursor
    }

    span.SetAttributes(
        attribute.Int("dingtalk.total_users", len(allUsers)),
    )
    span.SetStatus(codes.Ok, "success")

    return allUsers, nil
}
```

## 注意事项

1. **总是使用 defer span.End()**：确保 span 在函数返回时被结束
2. **传递 context**：确保将包含 span 的 context 传递给子操作
3. **记录错误**：使用 `span.RecordError()` 和 `span.SetStatus()` 记录错误
4. **添加属性**：添加有助于调试和监控的属性
5. **避免过度追踪**：不要为每个小操作都创建 span，只在关键操作上使用

## 相关文档

- [OpenTelemetry Go 文档](https://opentelemetry.io/docs/instrumentation/go/)
- [Kratos Tracing 中间件](https://github.com/go-kratos/kratos/tree/main/middleware/tracing)
- [项目日志集成 OpenTelemetry Trace](./internal/logger/ADD_OPENTELEMETRY_TRACE.md)

