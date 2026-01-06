# 第三方服务集成指南 - 第四步：在 Repository 中使用第三方服务

## 概述

本文档介绍如何在 Repository 中正确使用已集成的第三方服务，包括类型转换、错误处理、事务管理等最佳实践。

## 步骤 1: 在 Repository 中获取客户端

### 1.1 通过 Data 结构体获取

在 Repository 中通过 `Data` 结构体获取第三方服务客户端：

```go
package data

import (
	"context"
	"sre/internal/biz"
	"sre/internal/data/external/payment"
	"sre/internal/data/clients"

	"github.com/go-kratos/kratos/v2/log"
)

type orderRepo struct {
	data          *Data
	paymentClient *payment.Client
	log           *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	// 从 HTTP 客户端管理器获取
	paymentClient, err := data.httpClients.GetPaymentClient()
	if err != nil {
		log.NewHelper(logger).Warnf("payment client not available: %v", err)
	}

	return &orderRepo{
		data:          data,
		paymentClient: paymentClient,
		log:           log.NewHelper(logger),
	}
}
```

### 1.2 通过 gRPC 客户端管理器获取

```go
func (r *orderRepo) GetUserFromService(ctx context.Context, userID int64) (*biz.User, error) {
	// 获取 gRPC 客户端连接
	conn, err := r.data.grpcClients.GetClient("user-service")
	if err != nil {
		return nil, err
	}

	// 创建服务客户端
	client := userservice.NewClient(conn, r.log)
	
	// 调用服务
	user, err := client.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为业务类型
	return &biz.User{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
```

## 步骤 2: 类型转换

### 2.1 第三方类型转换为业务类型

在 Data 层进行类型转换，业务层不依赖第三方服务的具体类型：

```go
// 将 gRPC 类型转换为业务类型
func (r *orderRepo) convertUserFromGRPC(user *v1.User) *biz.User {
	return &biz.User{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Unix(user.CreatedAt, 0),
	}
}

// 将 HTTP 响应类型转换为业务类型
func (r *orderRepo) convertPaymentFromHTTP(resp *payment.CreatePaymentResponse) *biz.Payment {
	return &biz.Payment{
		ID:        resp.PaymentID,
		Status:    resp.Status,
		Amount:    resp.Amount,
		Currency:  resp.Currency,
		CreatedAt: parseTime(resp.CreatedAt),
	}
}
```

### 2.2 业务类型转换为第三方类型

```go
// 将业务类型转换为 gRPC 请求类型
func (r *orderRepo) convertCreateUserRequest(user *biz.User) *v1.CreateUserRequest {
	return &v1.CreateUserRequest{
		Name:  user.Name,
		Email: user.Email,
	}
}

// 将业务类型转换为 HTTP 请求类型
func (r *orderRepo) convertPaymentRequest(order *biz.Order) *payment.CreatePaymentRequest {
	return &payment.CreatePaymentRequest{
		OrderID:     order.ID,
		Amount:      order.Amount,
		Currency:    order.Currency,
		Description: order.Description,
	}
}
```

## 步骤 3: 错误处理

### 3.1 统一错误处理

将第三方服务的错误转换为业务层错误：

```go
func (r *orderRepo) CreatePayment(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	req := &payment.CreatePaymentRequest{
		OrderID:  order.ID,
		Amount:   order.Amount,
		Currency: order.Currency,
	}

	resp, err := r.paymentClient.CreatePayment(ctx, req)
	if err != nil {
		// 记录错误日志
		r.log.Errorf("failed to create payment for order %s: %v", order.ID, err)
		
		// 转换为业务错误
		if isTimeoutError(err) {
			return nil, errors.New(500, "PAYMENT_TIMEOUT", "支付服务超时，请稍后重试")
		}
		if isNetworkError(err) {
			return nil, errors.New(500, "PAYMENT_NETWORK_ERROR", "支付服务网络错误")
		}
		return nil, errors.New(500, "PAYMENT_SERVICE_ERROR", "支付服务错误")
	}

	return r.convertPaymentFromHTTP(resp), nil
}
```

### 3.2 使用 Kratos 错误定义

参考 `api/helloworld/v1/error_reason.proto` 定义错误原因：

```protobuf
enum ErrorReason {
  // 支付相关错误
  PAYMENT_SERVICE_ERROR = 1000;
  PAYMENT_TIMEOUT = 1001;
  PAYMENT_NETWORK_ERROR = 1002;
  PAYMENT_INVALID_REQUEST = 1003;
}
```

## 步骤 4: 组合多个第三方服务

在 Repository 中可以组合使用多个第三方服务：

```go
func (r *orderRepo) CreateOrderWithPayment(ctx context.Context, order *biz.Order) (*biz.Order, error) {
	// 1. 保存订单到数据库
	if err := r.data.db.Create(order).Error; err != nil {
		return nil, err
	}

	// 2. 调用支付服务创建支付
	paymentReq := &payment.CreatePaymentRequest{
		OrderID:  order.ID,
		Amount:   order.Amount,
		Currency: order.Currency,
	}
	paymentResp, err := r.paymentClient.CreatePayment(ctx, paymentReq)
	if err != nil {
		// 支付失败，回滚订单（可选）
		r.log.Errorf("payment failed for order %s, rolling back", order.ID)
		return nil, err
	}

	// 3. 更新订单支付信息
	order.PaymentID = paymentResp.PaymentID
	order.PaymentStatus = paymentResp.Status
	if err := r.data.db.Save(order).Error; err != nil {
		return nil, err
	}

	// 4. 发送通知（可选）
	if conn, err := r.data.grpcClients.GetClient("notification-service"); err == nil {
		notificationClient := notificationservice.NewClient(conn, r.log)
		notificationClient.SendNotification(ctx, &v1.SendNotificationRequest{
			UserID:  order.UserID,
			Title:   "订单创建成功",
			Content: fmt.Sprintf("订单 %s 已创建，支付金额 %.2f", order.ID, order.Amount),
		})
	}

	return order, nil
}
```

## 步骤 5: 缓存策略

对于频繁调用的第三方服务，可以添加缓存：

```go
func (r *orderRepo) GetUserFromService(ctx context.Context, userID int64) (*biz.User, error) {
	// 1. 先查缓存
	cacheKey := fmt.Sprintf("user:%d", userID)
	if r.data.redis != nil {
		val, err := r.data.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var user biz.User
			if json.Unmarshal([]byte(val), &user) == nil {
				return &user, nil
			}
		}
	}

	// 2. 调用第三方服务
	conn, err := r.data.grpcClients.GetClient("user-service")
	if err != nil {
		return nil, err
	}
	client := userservice.NewClient(conn, r.log)
	user, err := client.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. 转换为业务类型
	bizUser := r.convertUserFromGRPC(user)

	// 4. 写入缓存
	if r.data.redis != nil {
		data, _ := json.Marshal(bizUser)
		r.data.redis.Set(ctx, cacheKey, data, time.Hour)
	}

	return bizUser, nil
}
```

## 步骤 6: 超时和重试

### 6.1 设置超时

```go
func (r *orderRepo) CreatePaymentWithTimeout(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &payment.CreatePaymentRequest{
		OrderID:  order.ID,
		Amount:   order.Amount,
		Currency: order.Currency,
	}

	resp, err := r.paymentClient.CreatePayment(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New(500, "PAYMENT_TIMEOUT", "支付服务超时")
		}
		return nil, err
	}

	return r.convertPaymentFromHTTP(resp), nil
}
```

### 6.2 实现重试逻辑

```go
func (r *orderRepo) CreatePaymentWithRetry(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	var lastErr error
	maxRetries := 3
	retryDelay := 100 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		req := &payment.CreatePaymentRequest{
			OrderID:  order.ID,
			Amount:   order.Amount,
			Currency: order.Currency,
		}

		resp, err := r.paymentClient.CreatePayment(ctx, req)
		if err == nil {
			return r.convertPaymentFromHTTP(resp), nil
		}

		lastErr = err
		
		// 如果是不可重试的错误，直接返回
		if !isRetryableError(err) {
			return nil, err
		}

		// 等待后重试
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // 指数退避
		}
	}

	return nil, lastErr
}

func isRetryableError(err error) bool {
	// 判断是否为可重试的错误（如网络错误、超时、5xx 错误）
	return true // 根据实际情况实现
}
```

## 步骤 7: 监控和日志

### 7.1 记录调用日志

```go
func (r *orderRepo) CreatePayment(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	startTime := time.Now()
	
	req := &payment.CreatePaymentRequest{
		OrderID:  order.ID,
		Amount:   order.Amount,
		Currency: order.Currency,
	}

	resp, err := r.paymentClient.CreatePayment(ctx, req)
	
	duration := time.Since(startTime)
	
	if err != nil {
		r.log.Errorf("payment creation failed: order_id=%s, duration=%v, error=%v", 
			order.ID, duration, err)
		return nil, err
	}

	r.log.Infof("payment created: order_id=%s, payment_id=%s, duration=%v", 
		order.ID, resp.PaymentID, duration)

	return r.convertPaymentFromHTTP(resp), nil
}
```

### 7.2 添加指标收集

```go
func (r *orderRepo) CreatePayment(ctx context.Context, order *biz.Order) (*biz.Payment, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		// 记录指标
		metrics.RecordPaymentDuration(duration)
		metrics.IncrementPaymentRequests()
	}()

	// ... 调用支付服务 ...
}
```

## 最佳实践总结

1. **类型转换**：在 Data 层进行类型转换，业务层不依赖第三方类型
2. **错误处理**：统一错误处理，转换为业务层错误
3. **超时控制**：为第三方服务调用设置合理的超时时间
4. **重试策略**：对可重试的错误实现重试逻辑
5. **缓存策略**：对频繁调用的服务添加缓存
6. **日志记录**：记录关键操作的日志和指标
7. **优雅降级**：当第三方服务不可用时，提供降级方案
8. **资源管理**：确保正确关闭客户端连接

## 相关文档

- [第一步：准备工作](./third-party-integration-01-preparation.md)
- [第二步：gRPC 服务集成](./third-party-integration-02-grpc.md)
- [第三步：HTTP REST API 集成](./third-party-integration-03-http.md)
- [第三方服务接口定义最佳实践](../architecture/third-party-api-definitions.md)
- [数据访问层外部依赖管理](../architecture/data-layer-dependencies.md)

