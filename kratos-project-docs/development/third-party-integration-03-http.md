# 第三方服务集成指南 - 第三步：HTTP REST API 集成

## 概述

本文档介绍如何集成 HTTP REST API 第三方服务，包括定义请求/响应类型、创建 HTTP 客户端、以及在 Repository 中使用。

## 步骤 1: 定义请求和响应类型

### 1.1 创建目录结构

```bash
mkdir -p internal/data/external/{service-name}
```

**示例：**
```bash
mkdir -p internal/data/external/payment
```

### 1.2 定义类型

在 `internal/data/external/{service-name}/types.go` 中定义请求和响应类型。

**示例：`internal/data/external/payment/types.go`**

```go
package payment

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description,omitempty"`
	CustomerID  string  `json:"customer_id,omitempty"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentID   string `json:"payment_id"`
	Status      string `json:"status"`
	Amount      float64 `json:"amount"`
	Currency    string `json:"currency"`
	CreatedAt   string `json:"created_at"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// GetPaymentRequest 获取支付信息请求
type GetPaymentRequest struct {
	PaymentID string `json:"payment_id"`
}

// GetPaymentResponse 获取支付信息响应
type GetPaymentResponse struct {
	PaymentID string  `json:"payment_id"`
	OrderID   string  `json:"order_id"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
```

## 步骤 2: 创建 HTTP 客户端

### 2.1 创建客户端接口

在 `internal/data/external/{service-name}/client.go` 中创建 HTTP 客户端。

**示例：`internal/data/external/payment/client.go`**

```go
package payment

import (
	"context"
	"fmt"
	"time"

	"sre/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-resty/resty/v2"
)

// Client 支付服务 HTTP 客户端
type Client struct {
	client  *resty.Client
	baseURL string
	log     *log.Helper
}

// NewClient 创建支付服务客户端
func NewClient(cfg *conf.Data_HTTP_ClientConfig, logger log.Logger) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("payment client config is nil")
	}

	client := resty.New().
		SetBaseURL(cfg.Endpoint).
		SetTimeout(cfg.Timeout.AsDuration()).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	// 设置默认请求头
	if cfg.Headers != nil {
		for k, v := range cfg.Headers {
			client.SetHeader(k, v)
		}
	}

	// 设置重试策略
	client.SetRetryCount(3).
		SetRetryWaitTime(100 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second)

	// 设置请求日志（可选）
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.NewHelper(logger).Debugf("HTTP Request: %s %s", req.Method, req.URL)
		return nil
	})

	// 设置响应日志（可选）
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.NewHelper(logger).Debugf("HTTP Response: %d %s", resp.StatusCode(), resp.String())
		return nil
	})

	return &Client{
		client:  client,
		baseURL: cfg.Endpoint,
		log:     log.NewHelper(logger),
	}, nil
}

// CreatePayment 创建支付
func (c *Client) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	var resp CreatePaymentResponse
	var errResp ErrorResponse

	httpResp, err := c.client.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&resp).
		SetError(&errResp).
		Post("/api/v1/payments")

	if err != nil {
		c.log.Errorf("failed to create payment: %v", err)
		return nil, fmt.Errorf("payment service error: %w", err)
	}

	if !httpResp.IsSuccess() {
		c.log.Errorf("payment service returned error: %d, %s", httpResp.StatusCode(), errResp.Message)
		return nil, fmt.Errorf("payment service error: %s", errResp.Message)
	}

	return &resp, nil
}

// GetPayment 获取支付信息
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*GetPaymentResponse, error) {
	var resp GetPaymentResponse
	var errResp ErrorResponse

	httpResp, err := c.client.R().
		SetContext(ctx).
		SetPathParam("payment_id", paymentID).
		SetResult(&resp).
		SetError(&errResp).
		Get("/api/v1/payments/{payment_id}")

	if err != nil {
		c.log.Errorf("failed to get payment %s: %v", paymentID, err)
		return nil, fmt.Errorf("payment service error: %w", err)
	}

	if !httpResp.IsSuccess() {
		c.log.Errorf("payment service returned error: %d, %s", httpResp.StatusCode(), errResp.Message)
		return nil, fmt.Errorf("payment service error: %s", errResp.Message)
	}

	return &resp, nil
}
```

### 2.2 创建客户端管理器（可选）

如果需要管理多个 HTTP 客户端，可以创建管理器：

**示例：`internal/data/clients/http.go`**

```go
package clients

import (
	"fmt"
	"sync"

	"sre/internal/conf"
	"sre/internal/data/external/payment"

	"github.com/go-kratos/kratos/v2/log"
)

// HTTPClients 管理所有 HTTP 客户端
type HTTPClients struct {
	clients map[string]interface{}
	mu      sync.RWMutex
	log     *log.Helper
}

// NewHTTPClients 创建 HTTP 客户端管理器
func NewHTTPClients(c *conf.Data, logger log.Logger) (*HTTPClients, error) {
	if c.Http == nil || len(c.Http.Clients) == 0 {
		return &HTTPClients{
			clients: make(map[string]interface{}),
			log:     log.NewHelper(logger),
		}, nil
	}

	clients := &HTTPClients{
		clients: make(map[string]interface{}),
		log:     log.NewHelper(logger),
	}

	// 初始化支付客户端
	if cfg, ok := c.Http.Clients["payment-api"]; ok {
		paymentClient, err := payment.NewClient(cfg, logger)
		if err != nil {
			clients.log.Warnf("failed to create payment client: %v", err)
		} else {
			clients.clients["payment-api"] = paymentClient
			clients.log.Info("payment client initialized")
		}
	}

	return clients, nil
}

// GetPaymentClient 获取支付客户端
func (c *HTTPClients) GetPaymentClient() (*payment.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	client, ok := c.clients["payment-api"]
	if !ok {
		return nil, fmt.Errorf("payment client not found")
	}

	paymentClient, ok := client.(*payment.Client)
	if !ok {
		return nil, fmt.Errorf("invalid payment client type")
	}

	return paymentClient, nil
}
```

## 步骤 3: 创建业务封装（可选）

在 `internal/data/external/{service-name}/{service}.go` 中创建业务封装，提供更高级的 API。

**示例：`internal/data/external/payment/payment.go`**

```go
package payment

import (
	"context"
	"sre/internal/biz"
)

// Service 支付服务业务封装
type Service struct {
	client *Client
}

// NewService 创建支付服务
func NewService(client *Client) *Service {
	return &Service{
		client: client,
	}
}

// CreatePayment 创建支付（业务层接口）
func (s *Service) CreatePayment(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	req := &CreatePaymentRequest{
		OrderID:     order.ID,
		Amount:      order.Amount,
		Currency:    order.Currency,
		Description: order.Description,
		CustomerID:  order.CustomerID,
	}

	resp, err := s.client.CreatePayment(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换为业务类型
	return &biz.Payment{
		ID:        resp.PaymentID,
		OrderID:   order.ID,
		Status:    resp.Status,
		Amount:    resp.Amount,
		Currency:  resp.Currency,
		CreatedAt: resp.CreatedAt,
	}, nil
}
```

## 步骤 4: 在 Repository 中使用

在 Repository 中使用 HTTP 客户端：

```go
package data

import (
	"context"
	"sre/internal/biz"
	"sre/internal/data/external/payment"

	"github.com/go-kratos/kratos/v2/log"
)

type orderRepo struct {
	data          *Data
	paymentClient *payment.Client
	log           *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	// 从配置或 Data 中获取客户端
	paymentClient, _ := data.httpClients.GetPaymentClient()
	
	return &orderRepo{
		data:          data,
		paymentClient: paymentClient,
		log:           log.NewHelper(logger),
	}
}

func (r *orderRepo) CreatePayment(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	req := &payment.CreatePaymentRequest{
		OrderID:  order.ID,
		Amount:   order.Amount,
		Currency: order.Currency,
	}

	resp, err := r.paymentClient.CreatePayment(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换为业务类型
	return &biz.Payment{
		ID:        resp.PaymentID,
		OrderID:   order.ID,
		Status:    resp.Status,
		Amount:    resp.Amount,
		Currency:  resp.Currency,
		CreatedAt: resp.CreatedAt,
	}, nil
}
```

## 步骤 5: 更新 Data 结构体

在 `internal/data/data.go` 中添加 HTTP 客户端管理器：

```go
// Data 统一管理所有数据访问依赖
type Data struct {
	db          *gorm.DB
	httpClients *clients.HTTPClients  // 添加 HTTP 客户端
}

// NewData 创建 Data 实例
func NewData(
	db *gorm.DB,
	httpClients *clients.HTTPClients,
	logger log.Logger,
) (*Data, func(), error) {
	// ... 实现 ...
}
```

更新 `ProviderSet`：

```go
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	clients.NewHTTPClients,  // 添加 HTTP 客户端提供者
	NewUserRepo,
)
```

## 步骤 6: 配置服务地址

在 `configs/config.yaml` 中配置服务地址：

```yaml
data:
  http:
    clients:
      payment-api:
        endpoint: https://api.payment.com
        timeout: 5s
        headers:
          X-API-Key: your-api-key
          Authorization: Bearer your-token
```

## 高级配置

### 使用中间件

可以在客户端中添加中间件，如认证、日志、监控等：

```go
client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
	// 添加认证头
	req.SetHeader("Authorization", "Bearer "+token)
	return nil
})
```

### 错误处理和重试

```go
client.SetRetryCount(3).
	SetRetryWaitTime(100 * time.Millisecond).
	SetRetryMaxWaitTime(2 * time.Second).
	AddRetryCondition(func(r *resty.Response, err error) bool {
		// 只在 5xx 错误时重试
		return r.StatusCode() >= 500
	})
```

### 请求/响应拦截器

```go
// 请求拦截器
client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
	// 记录请求日志
	// 添加追踪 ID
	req.SetHeader("X-Trace-ID", generateTraceID())
	return nil
})

// 响应拦截器
client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
	// 记录响应日志
	// 处理错误
	if resp.StatusCode() >= 400 {
		// 记录错误指标
	}
	return nil
})
```

## 测试

创建单元测试时，可以使用 mock 来模拟 HTTP 客户端：

```go
//go:build !integration

package data_test

import (
	"testing"
	"github.com/jarcoal/httpmock"
)
```

## 下一步

完成 HTTP REST API 集成后，可以参考：
- [第四步：在 Repository 中使用第三方服务](./third-party-integration-04-usage.md)

