# Temporal 流程引擎

## 概述

Temporal 是由 Uber 开源的分布式工作流引擎，是目前 Go 语言中最流行和功能最完善的流程引擎。它提供了可靠的工作流执行、自动故障恢复、可观测性等企业级特性。

## 核心概念

### 1. Workflow（工作流）

Workflow 是业务逻辑的执行单元，定义了业务流程的步骤和顺序。

**特点**：
- 必须确定性（Deterministic）：相同的输入总是产生相同的输出
- 持久化执行：状态被持久化，可以恢复
- 长时间运行：可以运行数天或数月

### 2. Activity（活动）

Activity 是工作流中的实际执行单元，执行具体的业务逻辑（如调用外部服务、数据库操作等）。

**特点**：
- 可以执行非确定性的操作
- 支持重试和超时
- 可以访问外部资源

### 3. Task Queue（任务队列）

Task Queue 用于将任务分发给 Worker。

### 4. Worker（工作器）

Worker 是执行 Workflow 和 Activity 的进程。

## 架构

```
┌─────────────┐
│  Client     │
└──────┬──────┘
       │
       │ gRPC
       │
┌──────▼─────────────────────────────────┐
│        Temporal Server                 │
│  - Frontend (API Gateway)              │
│  - Matching (Task Queue)               │
│  - History (Event Store)               │
│  - Worker (Workflow Engine)            │
└──────┬─────────────────────────────────┘
       │
       ├─────────────┬─────────────┐
       │             │             │
┌──────▼─────┐  ┌────▼──────┐  ┌──▼──────┐
│  Worker 1  │  │  Worker 2 │  │ Worker N│
│ (Go App)   │  │ (Go App)  │  │ (Go App)│
└────────────┘  └───────────┘  └─────────┘
```

## 安装和配置

### 1. 安装 Temporal Server

#### 使用 Docker Compose（开发环境）

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgresql:
    image: postgres:14
    environment:
      POSTGRES_PWD: temporal
      POSTGRES_SEEDS: postgresql
    volumes:
      - postgresql-data:/var/lib/postgresql/data

  temporal:
    image: temporalio/auto-setup:latest
    ports:
      - "7233:7233"
    environment:
      - DB=postgresql
      - DB_PORT=5432
      - POSTGRES_USER=temporal
      - POSTGRES_PWD=temporal
      - POSTGRES_SEEDS=postgresql
    depends_on:
      - postgresql

  temporal-ui:
    image: temporalio/ui:latest
    ports:
      - "8088:8088"
    environment:
      - TEMPORAL_ADDRESS=temporal:7233
    depends_on:
      - temporal

volumes:
  postgresql-data:
```

启动服务：

```bash
docker-compose up -d
```

访问 Web UI：http://localhost:8088

### 2. 安装 Go SDK

```bash
go get go.temporal.io/sdk
```

## 基本使用

### 1. 定义 Activity

```go
// internal/workflow/activities.go
package workflow

import (
    "context"
    "time"
)

// CreateOrderActivity 创建订单的 Activity
func CreateOrderActivity(ctx context.Context, req CreateOrderRequest) (CreateOrderResponse, error) {
    // 执行实际的订单创建逻辑
    // 可以调用外部服务、数据库操作等
    
    // 模拟耗时操作
    time.Sleep(100 * time.Millisecond)
    
    return CreateOrderResponse{
        OrderID: "ORD123",
        Status:  "created",
    }, nil
}

// ReserveInventoryActivity 预留库存的 Activity
func ReserveInventoryActivity(ctx context.Context, req ReserveInventoryRequest) (ReserveInventoryResponse, error) {
    // 执行库存预留逻辑
    time.Sleep(100 * time.Millisecond)
    
    return ReserveInventoryResponse{
        ReserveID: "RESERVE123",
    }, nil
}

// FreezePaymentActivity 冻结支付的 Activity
func FreezePaymentActivity(ctx context.Context, req FreezePaymentRequest) (FreezePaymentResponse, error) {
    // 执行支付冻结逻辑
    time.Sleep(100 * time.Millisecond)
    
    return FreezePaymentResponse{
        FreezeID: "FREEZE123",
    }, nil
}
```

### 2. 定义 Workflow

```go
// internal/workflow/order_workflow.go
package workflow

import (
    "context"
    "time"
    
    "go.temporal.io/sdk/workflow"
)

// OrderWorkflowInput 工作流输入
type OrderWorkflowInput struct {
    UserID  int64
    Items   []OrderItem
    Amount  int64
}

// OrderWorkflowOutput 工作流输出
type OrderWorkflowOutput struct {
    OrderID   string
    ReserveID string
    FreezeID  string
}

// OrderWorkflow 订单创建工作流
func OrderWorkflow(ctx workflow.Context, input OrderWorkflowInput) (OrderWorkflowOutput, error) {
    // 配置选项
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,    // Activity 超时时间
        RetryPolicy: &workflow.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)
    
    var result OrderWorkflowOutput
    
    // Step 1: 创建订单
    createOrderReq := CreateOrderRequest{
        UserID: input.UserID,
        Items:  input.Items,
        Amount: input.Amount,
    }
    var createOrderResp CreateOrderResponse
    err := workflow.ExecuteActivity(ctx, CreateOrderActivity, createOrderReq).Get(ctx, &createOrderResp)
    if err != nil {
        return result, err
    }
    result.OrderID = createOrderResp.OrderID
    
    // Step 2: 预留库存
    reserveReq := ReserveInventoryRequest{
        OrderID: createOrderResp.OrderID,
        Items:   input.Items,
    }
    var reserveResp ReserveInventoryResponse
    err = workflow.ExecuteActivity(ctx, ReserveInventoryActivity, reserveReq).Get(ctx, &reserveResp)
    if err != nil {
        // 如果失败，需要补偿（取消订单）
        // 这里可以调用补偿 Activity
        return result, err
    }
    result.ReserveID = reserveResp.ReserveID
    
    // Step 3: 冻结支付
    freezeReq := FreezePaymentRequest{
        OrderID: createOrderResp.OrderID,
        Amount:  input.Amount,
    }
    var freezeResp FreezePaymentResponse
    err = workflow.ExecuteActivity(ctx, FreezePaymentActivity, freezeReq).Get(ctx, &freezeResp)
    if err != nil {
        // 如果失败，需要补偿（释放库存、取消订单）
        return result, err
    }
    result.FreezeID = freezeResp.FreezeID
    
    return result, nil
}
```

### 3. 创建 Worker

```go
// cmd/temporal-worker/main.go
package main

import (
    "log"
    
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    
    "sre/internal/workflow"
)

func main() {
    // 创建 Temporal Client
    c, err := client.Dial(client.Options{
        HostPort: "localhost:7233",
    })
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()
    
    // 创建 Worker
    w := worker.New(c, "order-task-queue", worker.Options{})
    
    // 注册 Workflow
    w.RegisterWorkflow(workflow.OrderWorkflow)
    
    // 注册 Activity
    w.RegisterActivity(workflow.CreateOrderActivity)
    w.RegisterActivity(workflow.ReserveInventoryActivity)
    w.RegisterActivity(workflow.FreezePaymentActivity)
    
    // 启动 Worker
    err = w.Run(worker.InterruptCh())
    if err != nil {
        log.Fatalln("Unable to start worker", err)
    }
}
```

### 4. 启动 Workflow

```go
// internal/service/order_service.go
package service

import (
    "context"
    
    "go.temporal.io/sdk/client"
    
    "sre/internal/workflow"
)

type OrderService struct {
    temporalClient client.Client
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
    // 工作流输入
    workflowInput := workflow.OrderWorkflowInput{
        UserID: req.UserID,
        Items:  req.Items,
        Amount: req.Amount,
    }
    
    // 启动工作流
    options := client.StartWorkflowOptions{
        ID:        "order-" + generateOrderID(),
        TaskQueue: "order-task-queue",
    }
    
    we, err := s.temporalClient.ExecuteWorkflow(ctx, options, workflow.OrderWorkflow, workflowInput)
    if err != nil {
        return nil, err
    }
    
    // 等待工作流完成
    var result workflow.OrderWorkflowOutput
    err = we.Get(ctx, &result)
    if err != nil {
        return nil, err
    }
    
    return &CreateOrderResponse{
        OrderID: result.OrderID,
    }, nil
}
```

## 高级特性

### 1. 补偿事务（Saga 模式）

```go
func OrderWorkflowWithCompensation(ctx workflow.Context, input OrderWorkflowInput) (OrderWorkflowOutput, error) {
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
    }
    ctx = workflow.WithActivityOptions(ctx, ao)
    
    var compensations []func() error
    
    // Step 1: 创建订单
    var orderResp CreateOrderResponse
    err := workflow.ExecuteActivity(ctx, CreateOrderActivity, createOrderReq).Get(ctx, &orderResp)
    if err != nil {
        return result, err
    }
    // 记录补偿操作
    compensations = append(compensations, func() error {
        return workflow.ExecuteActivity(ctx, CancelOrderActivity, orderResp.OrderID).Get(ctx, nil)
    })
    
    // Step 2: 预留库存
    var reserveResp ReserveInventoryResponse
    err = workflow.ExecuteActivity(ctx, ReserveInventoryActivity, reserveReq).Get(ctx, &reserveResp)
    if err != nil {
        // 执行补偿
        executeCompensations(ctx, compensations)
        return result, err
    }
    compensations = append(compensations, func() error {
        return workflow.ExecuteActivity(ctx, ReleaseInventoryActivity, reserveResp.ReserveID).Get(ctx, nil)
    })
    
    // Step 3: 冻结支付
    var freezeResp FreezePaymentResponse
    err = workflow.ExecuteActivity(ctx, FreezePaymentActivity, freezeReq).Get(ctx, &freezeResp)
    if err != nil {
        // 执行补偿
        executeCompensations(ctx, compensations)
        return result, err
    }
    
    return result, nil
}

func executeCompensations(ctx workflow.Context, compensations []func() error) {
    // 按相反顺序执行补偿
    for i := len(compensations) - 1; i >= 0; i-- {
        compensations[i]()
    }
}
```

### 2. 信号（Signal）

用于从外部向工作流发送事件：

```go
// 在工作流中接收信号
func OrderWorkflowWithSignal(ctx workflow.Context, input OrderWorkflowInput) (OrderWorkflowOutput, error) {
    signalChan := workflow.GetSignalChannel(ctx, "cancel-order")
    selector := workflow.NewSelector(ctx)
    
    selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
        var signal string
        c.Receive(ctx, &signal)
        if signal == "cancel" {
            // 处理取消逻辑
        }
    })
    
    // 等待信号或工作流完成
    selector.Select(ctx)
    
    return result, nil
}

// 从外部发送信号
func (s *OrderService) CancelOrder(ctx context.Context, workflowID string) error {
    err := s.temporalClient.SignalWorkflow(ctx, workflowID, "", "cancel-order", "cancel")
    return err
}
```

### 3. 查询（Query）

用于查询工作流的当前状态：

```go
// 在工作流中定义查询
func OrderWorkflowWithQuery(ctx workflow.Context, input OrderWorkflowInput) (OrderWorkflowOutput, error) {
    var currentStatus string
    
    // 注册查询处理器
    err := workflow.SetQueryHandler(ctx, "status", func() (string, error) {
        return currentStatus, nil
    })
    if err != nil {
        return result, err
    }
    
    // 更新状态
    currentStatus = "creating-order"
    // ...
    
    return result, nil
}

// 从外部查询状态
func (s *OrderService) GetOrderStatus(ctx context.Context, workflowID string) (string, error) {
    resp, err := s.temporalClient.QueryWorkflow(ctx, workflowID, "", "status")
    if err != nil {
        return "", err
    }
    
    var status string
    err = resp.Get(&status)
    return status, err
}
```

## 与项目集成

### 1. 配置

```yaml
# configs/config.yaml
temporal:
  address: localhost:7233  # Temporal Server 地址
  namespace: default       # 命名空间
  task_queue: order-task-queue  # 任务队列
```

### 2. 创建 Temporal Client

```go
// internal/data/client_temporal.go
package data

import (
    "go.temporal.io/sdk/client"
)

func NewTemporalClient(c *conf.Data) (client.Client, error) {
    opts := client.Options{
        HostPort: c.Temporal.Address,
        Namespace: c.Temporal.Namespace,
    }
    return client.Dial(opts)
}
```

### 3. 替换现有的 Saga 实现

可以考虑使用 Temporal 替换现有的 Saga 实现：

```go
// internal/biz/order_temporal.go
type OrderTemporalUseCase struct {
    temporalClient client.Client
}

func (uc *OrderTemporalUseCase) CreateOrder(ctx context.Context, req *CreateOrderRequest) error {
    workflowInput := workflow.OrderWorkflowInput{
        UserID: req.UserID,
        Items:  req.Items,
        Amount: req.Amount,
    }
    
    options := client.StartWorkflowOptions{
        ID:        "order-" + req.OrderNo,
        TaskQueue: "order-task-queue",
    }
    
    we, err := uc.temporalClient.ExecuteWorkflow(ctx, options, workflow.OrderWorkflow, workflowInput)
    if err != nil {
        return err
    }
    
    var result workflow.OrderWorkflowOutput
    return we.Get(ctx, &result)
}
```

## 最佳实践

### 1. Workflow 确定性

Workflow 代码必须是确定性的：

```go
// ❌ 错误：使用非确定性函数
func BadWorkflow(ctx workflow.Context) error {
    id := uuid.New()  // 每次执行结果不同
    return nil
}

// ✅ 正确：使用 workflow 提供的 API
func GoodWorkflow(ctx workflow.Context) error {
    id := workflow.GetInfo(ctx).WorkflowExecution.ID  // 确定性
    return nil
}
```

### 2. Activity 超时和重试

合理配置 Activity 的超时和重试策略：

```go
ao := workflow.ActivityOptions{
    StartToCloseTimeout: time.Minute,
    RetryPolicy: &workflow.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    time.Minute,
        MaximumAttempts:    3,
    },
}
```

### 3. 长时间运行的工作流

对于长时间运行的工作流，使用 `workflow.Sleep` 而不是 `time.Sleep`：

```go
// ❌ 错误
time.Sleep(24 * time.Hour)

// ✅ 正确
workflow.Sleep(ctx, 24*time.Hour)
```

## 参考资源

- [Temporal 官方文档](https://docs.temporal.io/)
- [Temporal Go SDK](https://github.com/temporalio/sdk-go)
- [Temporal 示例代码](https://github.com/temporalio/samples-go)

## 总结

Temporal 是一个功能强大、可靠的流程引擎，适合复杂业务流程的编排。它的主要优势：

- ✅ 功能完善，支持长时间运行的工作流
- ✅ 自动故障恢复和重试
- ✅ 完整的可观测性
- ✅ 活跃的社区和丰富的文档

对于需要高可靠性和复杂业务流程的场景，Temporal 是一个很好的选择。

