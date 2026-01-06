# Conductor 流程引擎

## 概述

Conductor 是 Netflix 开源的工作流编排引擎，基于 JSON 定义工作流，提供了任务调度、执行、重试等功能。Conductor 主要使用 Java 实现，但也提供了 Go 客户端支持（社区维护）。

## 核心概念

### 1. Workflow（工作流）

Workflow 是用 JSON 定义的业务流程，描述了任务的执行顺序和依赖关系。

### 2. Task（任务）

Task 是工作流中的执行单元，可以是系统任务（System Task）或用户任务（User Task）。

### 3. Worker（工作器）

Worker 是执行任务的进程。

## 架构

```
┌─────────────┐
│  Client     │
└──────┬──────┘
       │
       │ HTTP/REST
       │
┌──────▼─────────────────────────────────┐
│        Conductor Server                │
│  - Workflow Definition Store           │
│  - Task Queue                          │
│  - Execution Store                     │
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

### 1. 安装 Conductor Server

#### 使用 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_PASSWORD: conductor
      POSTGRES_DB: conductor

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.14.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"

  conductor-server:
    image: conductor/conductor-server:latest
    ports:
      - "8080:8080"
    environment:
      - CONFIG_PROP=config.properties
    depends_on:
      - postgres
      - elasticsearch
```

### 2. 安装 Go 客户端

```bash
go get github.com/conductor-sdk/conductor-go
```

## 基本使用

### 1. 定义 Workflow（JSON）

```json
{
  "name": "order_process",
  "description": "订单处理流程",
  "version": 1,
  "tasks": [
    {
      "name": "create_order",
      "taskReferenceName": "create_order_ref",
      "type": "SIMPLE",
      "inputParameters": {
        "userID": "${workflow.input.userID}",
        "items": "${workflow.input.items}",
        "amount": "${workflow.input.amount}"
      }
    },
    {
      "name": "reserve_inventory",
      "taskReferenceName": "reserve_inventory_ref",
      "type": "SIMPLE",
      "inputParameters": {
        "orderID": "${create_order_ref.output.orderID}"
      }
    },
    {
      "name": "freeze_payment",
      "taskReferenceName": "freeze_payment_ref",
      "type": "SIMPLE",
      "inputParameters": {
        "orderID": "${create_order_ref.output.orderID}",
        "amount": "${workflow.input.amount}"
      }
    }
  ],
  "outputParameters": {
    "orderID": "${create_order_ref.output.orderID}"
  }
}
```

### 2. 注册 Workflow

```go
// internal/workflow/conductor_register.go
package workflow

import (
    "context"
    "encoding/json"
    "io/ioutil"
    
    conductor "github.com/conductor-sdk/conductor-go/sdk/client"
)

func RegisterWorkflow(client *conductor.APIClient, workflowPath string) error {
    data, err := ioutil.ReadFile(workflowPath)
    if err != nil {
        return err
    }
    
    var workflowDef map[string]interface{}
    err = json.Unmarshal(data, &workflowDef)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    _, err = client.WorkflowResourceApi.CreateOrUpdateWorkflowDef(ctx, workflowDef)
    
    return err
}
```

### 3. 创建 Worker

```go
// cmd/conductor-worker/main.go
package main

import (
    "context"
    "log"
    
    conductor "github.com/conductor-sdk/conductor-go/sdk/client"
    "github.com/conductor-sdk/conductor-go/sdk/worker"
)

func main() {
    // 创建 Conductor Client
    apiClient := conductor.NewAPIClient(&conductor.Configuration{
        BasePath: "http://localhost:8080/api",
    })
    
    // 创建 Task Runner
    taskRunner := worker.NewTaskRunner(apiClient)
    
    // 注册任务处理器
    taskRunner.Start("create_order", handleCreateOrder, 1)  // 1 个并发
    taskRunner.Start("reserve_inventory", handleReserveInventory, 1)
    taskRunner.Start("freeze_payment", handleFreezePayment, 1)
    
    // 保持运行
    select {}
}

func handleCreateOrder(task *conductor.Task) (*conductor.TaskResult, error) {
    // 获取输入参数
    userID := task.InputData["userID"].(float64)
    amount := task.InputData["amount"].(float64)
    
    // 执行业务逻辑
    orderID := createOrder(int64(userID), int64(amount))
    
    // 返回结果
    return &conductor.TaskResult{
        TaskId: task.TaskId,
        Status: conductor.COMPLETED,
        OutputData: map[string]interface{}{
            "orderID": orderID,
        },
    }, nil
}

func handleReserveInventory(task *conductor.Task) (*conductor.TaskResult, error) {
    // 类似的处理逻辑
}

func handleFreezePayment(task *conductor.Task) (*conductor.TaskResult, error) {
    // 类似的处理逻辑
}
```

### 4. 启动 Workflow

```go
// internal/service/order_service.go
package service

import (
    "context"
    
    conductor "github.com/conductor-sdk/conductor-go/sdk/client"
)

type OrderService struct {
    conductorClient *conductor.APIClient
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
    // 启动工作流
    startRequest := conductor.StartWorkflowRequest{
        Name:        "order_process",
        Version:     1,
        Input: map[string]interface{}{
            "userID": req.UserID,
            "items":  req.Items,
            "amount": req.Amount,
        },
    }
    
    workflowId, err := s.conductorClient.WorkflowResourceApi.StartWorkflow(ctx, startRequest)
    if err != nil {
        return nil, err
    }
    
    return &CreateOrderResponse{
        WorkflowID: workflowId,
    }, nil
}
```

## 高级特性

### 1. 条件分支（Switch Task）

```json
{
  "name": "check_inventory",
  "type": "SWITCH",
  "decisionCases": {
    "available": [
      {
        "name": "reserve_inventory",
        "type": "SIMPLE"
      }
    ],
    "unavailable": [
      {
        "name": "notify_user",
        "type": "SIMPLE"
      }
    ]
  },
  "defaultCase": [
    {
      "name": "wait_for_inventory",
      "type": "SIMPLE"
    }
  ]
}
```

### 2. 并行执行（Fork/Join）

```json
{
  "name": "fork_tasks",
  "type": "FORK_JOIN",
  "forkTasks": [
    [
      {
        "name": "task1",
        "type": "SIMPLE"
      }
    ],
    [
      {
        "name": "task2",
        "type": "SIMPLE"
      }
    ]
  ],
  "joinOn": ["task1", "task2"]
}
```

### 3. 循环（Do-While Task）

```json
{
  "name": "retry_payment",
  "type": "DO_WHILE",
  "loopCondition": "${retry_payment_ref.output.retryCount < 3}",
  "loopOver": [
    {
      "name": "attempt_payment",
      "type": "SIMPLE"
    }
  ]
}
```

## 与项目集成

### 1. 配置

```yaml
# configs/config.yaml
conductor:
  server_url: http://localhost:8080/api  # Conductor Server URL
  workflow_name: order_process           # 工作流名称
```

### 2. 创建 Conductor Client

```go
// internal/data/client_conductor.go
package data

import (
    conductor "github.com/conductor-sdk/conductor-go/sdk/client"
)

func NewConductorClient(c *conf.Data) (*conductor.APIClient, error) {
    return conductor.NewAPIClient(&conductor.Configuration{
        BasePath: c.Conductor.ServerUrl,
    }), nil
}
```

## 最佳实践

### 1. Workflow 设计

- 保持工作流简洁，避免过于复杂的嵌套
- 合理使用系统任务（System Task）和用户任务（User Task）
- 使用输入/输出参数传递数据

### 2. 任务处理

- 实现幂等性
- 合理设置超时和重试
- 处理错误情况

### 3. 性能优化

- 使用多个 Worker 实例
- 合理设置任务并发数
- 使用并行执行提高效率

## 与其他引擎的对比

| 特性 | Conductor | Temporal | Zeebe |
|------|-----------|----------|-------|
| 工作流定义 | JSON | 代码 | BPMN |
| 学习曲线 | 低 | 中 | 中 |
| 社区支持 | 中 | 高 | 中 |
| Netflix 使用 | ✅ | ❌ | ❌ |
| 适用场景 | 任务编排 | 通用 | BPM |

## 参考资源

- [Conductor GitHub](https://github.com/Netflix/conductor)
- [Conductor 文档](https://conductor.netflix.com/)
- [Conductor Go SDK](https://github.com/conductor-sdk/conductor-go)

## 总结

Conductor 是一个基于 JSON 的工作流编排引擎，适合 Netflix 风格的微服务编排场景。它的主要特点：

- ✅ JSON 配置驱动，易于理解
- ✅ 支持复杂的任务编排
- ✅ Netflix 内部使用，经过大规模验证
- ⚠️ Go 客户端支持由社区维护

如果需要 JSON 配置驱动的工作流编排，Conductor 是一个可以考虑的选择。

