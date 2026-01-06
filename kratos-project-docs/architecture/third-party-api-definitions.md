# 第三方服务接口定义最佳实践

## 概述

在微服务架构中，调用第三方服务（包括其他微服务、外部 API）时，需要定义请求和响应的数据结构。本文档介绍如何组织和管理这些接口定义，保持代码清晰、可维护。

## 分类和组织原则

### 1. 按服务类型分类

- **gRPC 服务**：使用 Protobuf 定义
- **HTTP REST API**：使用 Go 结构体定义
- **GraphQL API**：使用 GraphQL Schema 或 Go 结构体

### 2. 按服务来源分类

- **内部微服务**：同一组织内的其他服务
- **外部第三方服务**：外部供应商提供的服务
- **公共 API**：公开的 API 服务

## 目录结构

### 推荐结构

```
sre/
├── api/                          # API 定义层
│   ├── helloworld/              # 本服务的 API 定义
│   │   └── v1/
│   └── external/                # 外部服务的 gRPC 定义
│       ├── user-service/        # 用户服务
│       │   └── v1/
│       │       └── user.proto
│       └── order-service/       # 订单服务
│           └── v1/
│               └── order.proto
├── internal/
│   └── data/
│       ├── clients/             # 客户端封装
│       │   ├── grpc.go
│       │   └── http.go
│       └── external/            # 第三方服务接口定义
│           ├── payment/        # 支付服务
│           │   ├── client.go   # 客户端实现
│           │   ├── types.go    # 请求/响应类型定义
│           │   └── payment.go  # 业务封装
│           ├── notification/   # 通知服务
│           │   ├── client.go
│           │   ├── types.go
│           │   └── notification.go
│           └── sms/            # 短信服务
│               ├── client.go
│               ├── types.go
│               └── sms.go
└── third_party/                 # 第三方 Protobuf 依赖
    ├── google/
    └── openapi/
```

## gRPC 服务接口定义

### 内部微服务（推荐方式）

对于同一组织内的微服务，将 proto 定义放在 `api/external/` 目录下：

```
api/
└── external/
    └── user-service/
        └── v1/
            ├── user.proto
            ├── user.pb.go
            └── user_grpc.pb.go
```

**示例：user.proto**

```protobuf
syntax = "proto3";

package user.service.v1;

option go_package = "sre/api/external/user-service/v1;v1";

// 用户服务定义
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

// 请求消息
message GetUserRequest {
  int64 id = 1;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

// 响应消息
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int64 created_at = 4;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}
```

**在 Data 层使用**

```go
// internal/data/external/user/user.go
package external

import (
	"context"
	
	"sre/api/external/user-service/v1"
	"sre/internal/biz"
	
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	client v1.UserServiceClient
	logger log.Logger
}

func NewUserServiceClient(conn *grpc.ClientConn, logger log.Logger) *UserServiceClient {
	return &UserServiceClient{
		client: v1.NewUserServiceClient(conn),
		logger: log.NewHelper(logger),
	}
}

func (c *UserServiceClient) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	resp, err := c.client.GetUser(ctx, &v1.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	
	// 转换为业务对象
	return &biz.User{
		ID:    resp.Id,
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}
```

### 外部第三方 gRPC 服务

如果第三方服务提供了 proto 定义，可以：

1. **直接使用第三方 proto**（推荐）
   - 将 proto 文件放在 `third_party/` 目录
   - 保持原始包名和结构

2. **重新定义适配层**（如果第三方 proto 不符合规范）
   - 在 `api/external/` 下定义适配的 proto
   - 在 Data 层做转换

## HTTP REST API 接口定义

### 组织方式

对于 HTTP REST API，将请求/响应类型定义放在 `internal/data/external/{service}/types.go`：

```
internal/data/external/
└── payment/
    ├── types.go        # 请求/响应类型定义
    ├── client.go       # HTTP 客户端实现
    └── payment.go      # 业务封装
```

### 示例：支付服务

**types.go - 请求/响应类型定义**

```go
// internal/data/external/payment/types.go
package payment

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	ReturnURL   string  `json:"return_url"`
	NotifyURL   string  `json:"notify_url"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentID   string `json:"payment_id"`
	PaymentURL  string `json:"payment_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// QueryPaymentRequest 查询支付请求
type QueryPaymentRequest struct {
	PaymentID string `json:"payment_id"`
}

// QueryPaymentResponse 查询支付响应
type QueryPaymentResponse struct {
	PaymentID   string  `json:"payment_id"`
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	PaidAt      string  `json:"paid_at,omitempty"`
	FailedAt    string  `json:"failed_at,omitempty"`
	FailureCode string  `json:"failure_code,omitempty"`
}

// PaymentError 支付服务错误响应
type PaymentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *PaymentError) Error() string {
	return e.Message
}
```

**client.go - HTTP 客户端实现**

```go
// internal/data/external/payment/client.go
package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
)

// Client 支付服务客户端
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     log.Logger
}

// NewClient 创建支付服务客户端
func NewClient(c *conf.Data, logger log.Logger) (*Client, error) {
	if c.Payment == nil {
		return nil, fmt.Errorf("payment config not found")
	}
	
	return &Client{
		baseURL: c.Payment.BaseURL,
		apiKey:  c.Payment.ApiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: log.NewHelper(logger),
	}, nil
}

// CreatePayment 创建支付
func (c *Client) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments", c.baseURL)
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		var paymentErr PaymentError
		if err := json.Unmarshal(respBody, &paymentErr); err == nil {
			return nil, &paymentErr
		}
		return nil, fmt.Errorf("payment API error: %s", string(respBody))
	}
	
	var result CreatePaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// QueryPayment 查询支付状态
func (c *Client) QueryPayment(ctx context.Context, req *QueryPaymentRequest) (*QueryPaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s", c.baseURL, req.PaymentID)
	
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		var paymentErr PaymentError
		if err := json.Unmarshal(respBody, &paymentErr); err == nil {
			return nil, &paymentErr
		}
		return nil, fmt.Errorf("payment API error: %s", string(respBody))
	}
	
	var result QueryPaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}
```

**payment.go - 业务封装**

```go
// internal/data/external/payment/payment.go
package payment

import (
	"context"
	
	"sre/internal/biz"
	
	"github.com/go-kratos/kratos/v2/log"
)

// PaymentService 支付服务封装
type PaymentService struct {
	client *Client
	logger log.Logger
}

// NewPaymentService 创建支付服务
func NewPaymentService(client *Client, logger log.Logger) *PaymentService {
	return &PaymentService{
		client: client,
		logger: log.NewHelper(logger),
	}
}

// CreatePayment 创建支付（业务层接口）
func (s *PaymentService) CreatePayment(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	req := &CreatePaymentRequest{
		OrderID:     order.ID,
		Amount:      order.Amount,
		Currency:    "CNY",
		Description: order.Description,
		ReturnURL:   order.ReturnURL,
		NotifyURL:   order.NotifyURL,
	}
	
	resp, err := s.client.CreatePayment(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 转换为业务对象
	return &biz.Payment{
		PaymentID:  resp.PaymentID,
		PaymentURL: resp.PaymentURL,
		Status:      resp.Status,
	}, nil
}
```

## 命名规范

### 1. 文件命名

- **types.go**：请求/响应类型定义
- **client.go**：客户端实现
- **{service}.go**：业务封装（可选）

### 2. 类型命名

- **请求类型**：`{Action}{Resource}Request`
  - 例如：`CreatePaymentRequest`、`QueryPaymentRequest`
- **响应类型**：`{Action}{Resource}Response`
  - 例如：`CreatePaymentResponse`、`QueryPaymentResponse`
- **错误类型**：`{Service}Error`
  - 例如：`PaymentError`、`SMSError`

### 3. 包命名

- 使用服务名称作为包名（小写）
- 例如：`payment`、`sms`、`notification`

## 最佳实践

### 1. 类型定义原则

#### ✅ 推荐做法

```go
// 为每个 API 端点定义独立的请求/响应类型
type CreatePaymentRequest struct {
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type CreatePaymentResponse struct {
	PaymentID  string `json:"payment_id"`
	PaymentURL string `json:"payment_url"`
	Status     string `json:"status"`
}
```

#### ❌ 不推荐做法

```go
// 不要使用通用的 map 或 interface{}
func CreatePayment(req map[string]interface{}) (map[string]interface{}, error)

// 不要复用不相关的类型
type PaymentRequest struct {
	// 混合了创建和查询的字段
	PaymentID string  // 查询用
	Amount    float64 // 创建用
}
```

### 2. 错误处理

```go
// 定义服务特定的错误类型
type PaymentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *PaymentError) Error() string {
	return fmt.Sprintf("payment error [%s]: %s", e.Code, e.Message)
}

// 在客户端中处理错误
if resp.StatusCode != http.StatusOK {
	var paymentErr PaymentError
	if err := json.Unmarshal(respBody, &paymentErr); err == nil {
		return nil, &paymentErr
	}
	return nil, fmt.Errorf("payment API error: %s", string(respBody))
}
```

### 3. 类型转换

在 Data 层进行第三方类型和业务类型的转换：

```go
// 第三方类型 -> 业务类型
func toBizPayment(resp *CreatePaymentResponse) *biz.Payment {
	return &biz.Payment{
		PaymentID:  resp.PaymentID,
		PaymentURL: resp.PaymentURL,
		Status:     resp.Status,
	}
}

// 业务类型 -> 第三方类型
func fromBizOrder(order *biz.Order) *CreatePaymentRequest {
	return &CreatePaymentRequest{
		OrderID:   order.ID,
		Amount:    order.Amount,
		Currency:  "CNY",
		// ...
	}
}
```

### 4. 配置管理

在配置中定义第三方服务的连接信息：

```protobuf
// internal/conf/conf.proto
message Data {
  message Payment {
    string base_url = 1;
    string api_key = 2;
    google.protobuf.Duration timeout = 3;
  }
  
  Payment payment = 1;
}
```

```yaml
# configs/config.yaml
data:
  payment:
    base_url: https://api.payment.com
    api_key: your-api-key
    timeout: 10s
```

### 5. 版本管理

对于可能变更的第三方 API，使用版本号：

```
internal/data/external/
└── payment/
    ├── v1/
    │   ├── types.go
    │   └── client.go
    └── v2/
        ├── types.go
        └── client.go
```

### 6. 文档注释

为所有公开的类型和方法添加文档注释：

```go
// CreatePaymentRequest 创建支付请求
// 
// 字段说明：
//   - OrderID: 订单ID，必填
//   - Amount: 支付金额，单位：元，必填
//   - Currency: 货币类型，默认：CNY
type CreatePaymentRequest struct {
	OrderID  string  `json:"order_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency"`
}
```

### 7. 测试支持

为测试提供 Mock 实现：

```go
// internal/data/external/payment/mock.go
package payment

type MockClient struct {
	CreatePaymentFunc func(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	QueryPaymentFunc  func(ctx context.Context, req *QueryPaymentRequest) (*QueryPaymentResponse, error)
}

func (m *MockClient) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	if m.CreatePaymentFunc != nil {
		return m.CreatePaymentFunc(ctx, req)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *MockClient) QueryPayment(ctx context.Context, req *QueryPaymentRequest) (*QueryPaymentResponse, error) {
	if m.QueryPaymentFunc != nil {
		return m.QueryPaymentFunc(ctx, req)
	}
	return nil, fmt.Errorf("not implemented")
}
```

## 目录组织决策树

```
调用第三方服务？
├── gRPC 服务？
│   ├── 内部微服务？
│   │   └── api/external/{service}/v1/*.proto
│   └── 外部服务？
│       └── third_party/{vendor}/*.proto 或
│           └── api/external/{service}/v1/*.proto (适配层)
└── HTTP REST API？
    └── internal/data/external/{service}/
        ├── types.go      (请求/响应类型)
        ├── client.go     (HTTP 客户端)
        └── {service}.go  (业务封装，可选)
```

## 总结

### 核心原则

1. **按服务组织**：每个第三方服务有独立的目录
2. **类型明确**：为每个 API 端点定义独立的请求/响应类型
3. **职责分离**：类型定义、客户端实现、业务封装分离
4. **版本管理**：支持 API 版本演进
5. **错误处理**：定义服务特定的错误类型
6. **类型转换**：在 Data 层进行第三方类型和业务类型的转换

### 目录位置总结

| 服务类型 | 定义位置 | 说明 |
|---------|---------|------|
| 内部 gRPC 服务 | `api/external/{service}/v1/` | 同一组织的微服务 |
| 外部 gRPC 服务 | `third_party/{vendor}/` 或 `api/external/{service}/v1/` | 第三方 proto 或适配层 |
| HTTP REST API | `internal/data/external/{service}/types.go` | 请求/响应类型 |
| HTTP 客户端 | `internal/data/external/{service}/client.go` | 客户端实现 |

通过以上组织方式，可以清晰地管理所有第三方服务的接口定义，保持代码的可维护性和可测试性。

