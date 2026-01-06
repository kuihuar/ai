# 分层架构设计

## 设计原则

### 单一职责原则
每一层只负责一个明确的职责，避免职责混乱。

### 依赖倒置原则
上层定义接口，下层实现接口，实现依赖倒置。

### 关注点分离
将不同的关注点（业务逻辑、数据访问、协议处理）分离到不同层。

## 层间通信

### 数据流向

```
请求流程：
Client → API → Service → Biz → Data → Database

响应流程：
Database → Data → Biz → Service → API → Client
```

### 接口定义

Biz 层定义接口，Data 层实现：

```go
// internal/biz/greeter.go
type GreeterRepo interface {
    Save(context.Context, *Greeter) error
    FindByID(context.Context, int64) (*Greeter, error)
}

// internal/data/greeter.go
type greeterRepo struct {
    data *Data
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) error {
    // 实现数据持久化
}
```

## 各层详细说明

### Service 层
- **输入**：接收来自 API 层的请求
- **处理**：参数校验、协议转换
- **输出**：调用 Biz 层，返回结果

### Biz 层
- **输入**：接收来自 Service 层的业务请求
- **处理**：业务逻辑、规则验证、流程编排
- **输出**：调用 Data 层接口，返回业务对象

### Data 层
- **输入**：接收来自 Biz 层的数据操作请求
- **处理**：数据库操作、外部服务调用、缓存操作
- **输出**：返回数据实体或错误

## 常见问题

### 1. 业务逻辑应该放在哪一层？
**答案**：Biz 层。Service 层只做协议转换，Data 层只做数据访问。

### 2. 如何避免循环依赖？
**答案**：严格遵循依赖方向，上层定义接口，下层实现。

### 3. 跨层调用是否允许？
**答案**：不允许。必须通过相邻层进行调用。

