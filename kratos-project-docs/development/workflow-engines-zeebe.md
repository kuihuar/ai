# Zeebe (Camunda) 流程引擎

## 概述

Zeebe 是 Camunda 开源的分布式工作流引擎，基于 BPMN 2.0 标准，提供了高性能、可扩展的流程编排能力。Zeebe 支持 Go 客户端，适合需要 BPMN 标准支持和图形化建模的场景。

## 核心概念

### 1. BPMN 2.0

BPMN（Business Process Model and Notation）是业务流程建模和标注的标准，Zeebe 支持 BPMN 2.0 标准。

### 2. Process（流程）

Process 是 BPMN 定义的业务流程，可以通过图形化工具（Camunda Modeler）设计。

### 3. Job（任务）

Job 是流程中的执行单元，由 Worker 处理。

### 4. Worker（工作器）

Worker 是执行 Job 的进程。

## 架构

```
┌─────────────┐
│  Client     │
└──────┬──────┘
       │
       │ gRPC
       │
┌──────▼─────────────────────────────────┐
│        Zeebe Gateway                   │
│  - gRPC API                            │
│  - REST API                            │
└──────┬─────────────────────────────────┘
       │
┌──────▼─────────────────────────────────┐
│        Zeebe Broker                    │
│  - Process Engine                      │
│  - Job Queue                           │
│  - Event Store                         │
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

### 1. 安装 Zeebe

#### 使用 Docker（开发环境）

```bash
docker run -p 26500:26500 camunda/zeebe:latest
```

#### 使用 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  zeebe:
    image: camunda/zeebe:latest
    ports:
      - "26500:26500"
      - "9600:9600"
    environment:
      - ZEEBE_BROKER_CLUSTER_NODEID=0
      - ZEEBE_BROKER_CLUSTER_CLUSTERNAME=zeebe-cluster
      - ZEEBE_BROKER_CLUSTER_REPLICATIONFACTOR=1
      - ZEEBE_BROKER_CLUSTER_PARTITIONSCOUNT=1
```

### 2. 安装 Go 客户端

```bash
go get github.com/camunda/zeebe/clients/go/v8/pkg/zbc
```

## 基本使用

### 1. 创建 BPMN 流程（使用 Camunda Modeler）

使用 Camunda Modeler（图形化工具）创建 BPMN 流程文件（`.bpmn`），或者手动编写 XML。

示例 BPMN 流程（`order-process.bpmn`）：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL"
                  xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI"
                  xmlns:zeebe="http://camunda.org/schema/zeebe/1.0">
  <bpmn:process id="order-process" isExecutable="true">
    <bpmn:startEvent id="start"/>
    <bpmn:serviceTask id="create-order" name="创建订单">
      <bpmn:extensionElements>
        <zeebe:taskDefinition type="create-order" />
      </bpmn:extensionElements>
    </bpmn:serviceTask>
    <bpmn:serviceTask id="reserve-inventory" name="预留库存">
      <bpmn:extensionElements>
        <zeebe:taskDefinition type="reserve-inventory" />
      </bpmn:extensionElements>
    </bpmn:serviceTask>
    <bpmn:serviceTask id="freeze-payment" name="冻结支付">
      <bpmn:extensionElements>
        <zeebe:taskDefinition type="freeze-payment" />
      </bpmn:extensionElements>
    </bpmn:serviceTask>
    <bpmn:endEvent id="end"/>
    <bpmn:sequenceFlow id="flow1" sourceRef="start" targetRef="create-order"/>
    <bpmn:sequenceFlow id="flow2" sourceRef="create-order" targetRef="reserve-inventory"/>
    <bpmn:sequenceFlow id="flow3" sourceRef="reserve-inventory" targetRef="freeze-payment"/>
    <bpmn:sequenceFlow id="flow4" sourceRef="freeze-payment" targetRef="end"/>
  </bpmn:process>
</bpmn:definitions>
```

### 2. 部署流程

```go
// internal/workflow/zeebe_deploy.go
package workflow

import (
    "context"
    "os"
    
    "github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func DeployProcess(client zbc.Client, bpmnPath string) error {
    file, err := os.Open(bpmnPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    ctx := context.Background()
    response, err := client.NewDeployResourceCommand().
        AddResourceFile(bpmnPath).
        Send(ctx)
    
    if err != nil {
        return err
    }
    
    // 获取部署的流程定义
    process := response.GetProcesses()[0]
    fmt.Printf("Deployed process: %s (version: %d)\n", process.BpmnProcessId, process.Version)
    
    return nil
}
```

### 3. 创建 Worker

```go
// cmd/zeebe-worker/main.go
package main

import (
    "context"
    "log"
    
    "github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
    "github.com/camunda/zeebe/clients/go/v8/pkg/entities"
    "github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

func main() {
    // 创建 Zeebe Client
    client, err := zbc.NewClient(&zbc.ClientConfig{
        GatewayAddress: "localhost:26500",
    })
    if err != nil {
        log.Fatalln("Failed to create Zeebe client", err)
    }
    defer client.Close()
    
    // 创建 Job Worker
    jobWorker := client.NewJobWorker().
        JobType("create-order").
        Handler(handleCreateOrder).
        Open()
    defer jobWorker.Close()
    
    // 创建其他 Job Worker
    jobWorker2 := client.NewJobWorker().
        JobType("reserve-inventory").
        Handler(handleReserveInventory).
        Open()
    defer jobWorker2.Close()
    
    jobWorker3 := client.NewJobWorker().
        JobType("freeze-payment").
        Handler(handleFreezePayment).
        Open()
    defer jobWorker3.Close()
    
    // 保持运行
    select {}
}

func handleCreateOrder(client worker.JobClient, job entities.Job) {
    ctx := context.Background()
    
    // 获取变量
    variables, err := job.GetVariablesAsMap()
    if err != nil {
        failJob(client, job, err)
        return
    }
    
    // 执行业务逻辑
    userID := int64(variables["userID"].(float64))
    amount := int64(variables["amount"].(float64))
    
    // 创建订单
    orderID := createOrder(userID, amount)
    
    // 完成 Job
    request, err := client.NewCompleteJobCommand().JobKey(job.GetKey()).VariablesFromMap(map[string]interface{}{
        "orderID": orderID,
    })
    if err != nil {
        failJob(client, job, err)
        return
    }
    
    ctx = context.Background()
    _, err = request.Send(ctx)
    if err != nil {
        log.Printf("Failed to complete job %d: %v", job.GetKey(), err)
        return
    }
    
    log.Printf("Completed job %d with orderID: %s", job.GetKey(), orderID)
}

func handleReserveInventory(client worker.JobClient, job entities.Job) {
    // 类似的处理逻辑
}

func handleFreezePayment(client worker.JobClient, job entities.Job) {
    // 类似的处理逻辑
}

func failJob(client worker.JobClient, job entities.Job, err error) {
    ctx := context.Background()
    _, err = client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
    if err != nil {
        log.Printf("Failed to fail job %d: %v", job.GetKey(), err)
    }
}
```

### 4. 启动流程实例

```go
// internal/service/order_service.go
package service

import (
    "context"
    
    "github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

type OrderService struct {
    zeebeClient zbc.Client
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
    // 启动流程实例
    processInstanceResponse, err := s.zeebeClient.NewCreateInstanceCommand().
        BPMNProcessId("order-process").
        Version(1).
        VariablesFromMap(map[string]interface{}{
            "userID": req.UserID,
            "items":  req.Items,
            "amount": req.Amount,
        }).
        Send(ctx)
    
    if err != nil {
        return nil, err
    }
    
    processInstanceKey := processInstanceResponse.GetProcessInstanceKey()
    
    // 等待流程完成（可选）
    // 可以通过订阅事件或查询状态来获取结果
    
    return &CreateOrderResponse{
        ProcessInstanceKey: processInstanceKey,
    }, nil
}
```

## 高级特性

### 1. 补偿（Compensation）

Zeebe 支持 BPMN 的补偿事件：

```xml
<bpmn:boundaryEvent id="compensation-event" attachedToRef="freeze-payment">
  <bpmn:compensateEventDefinition />
</bpmn:boundaryEvent>
```

### 2. 消息事件（Message Events）

支持消息驱动的流程：

```xml
<bpmn:intermediateCatchEvent id="payment-confirmed">
  <bpmn:messageEventDefinition messageRef="PaymentConfirmedMessage" />
</bpmn:intermediateCatchEvent>
```

### 3. 定时器（Timer Events）

支持基于时间的流程：

```xml
<bpmn:boundaryEvent id="timeout" attachedToRef="freeze-payment">
  <bpmn:timerEventDefinition>
    <bpmn:timeDuration>PT1H</bpmn:timeDuration>
  </bpmn:timerEventDefinition>
</bpmn:boundaryEvent>
```

## 与项目集成

### 1. 配置

```yaml
# configs/config.yaml
zeebe:
  address: localhost:26500  # Zeebe Gateway 地址
  process_id: order-process # 流程 ID
```

### 2. 创建 Zeebe Client

```go
// internal/data/client_zeebe.go
package data

import (
    "github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func NewZeebeClient(c *conf.Data) (zbc.Client, error) {
    return zbc.NewClient(&zbc.ClientConfig{
        GatewayAddress: c.Zeebe.Address,
    })
}
```

## 最佳实践

### 1. BPMN 设计

- 使用 Camunda Modeler 进行可视化设计
- 保持流程简洁，避免过于复杂的嵌套
- 合理使用网关（Gateway）控制流程分支

### 2. Job 处理

- 实现幂等性：相同输入产生相同输出
- 合理设置重试次数和超时时间
- 处理错误情况，使用 Fail Job 命令

### 3. 性能优化

- 使用多个 Worker 实例提高并发处理能力
- 合理设置 Job Worker 的并发数
- 使用分区（Partition）分散负载

## 与 Temporal 的对比

| 特性 | Zeebe | Temporal |
|------|-------|----------|
| BPMN 支持 | ✅ 原生支持 | ❌ 不支持 |
| 图形化建模 | ✅ Camunda Modeler | ❌ 代码定义 |
| 性能 | 高 | 高 |
| 学习曲线 | 中（需要了解 BPMN） | 中 |
| 社区 | 中 | 高 |
| 适用场景 | BPM 场景 | 通用场景 |

## 参考资源

- [Zeebe 官方文档](https://docs.camunda.io/)
- [Zeebe Go 客户端](https://github.com/camunda/zeebe/clients/go)
- [Camunda Modeler](https://camunda.com/download/modeler/)
- [BPMN 2.0 规范](https://www.omg.org/spec/BPMN/2.0/)

## 总结

Zeebe 是一个基于 BPMN 标准的流程引擎，适合需要图形化建模和 BPMN 标准支持的场景。它的主要优势：

- ✅ 基于 BPMN 2.0 标准
- ✅ 提供图形化建模工具
- ✅ 高性能和可扩展性
- ✅ 适合企业级 BPM 需求

如果需要 BPMN 标准支持或图形化工作流设计，Zeebe 是一个很好的选择。

