# newApp 函数参数设计指南

## 概述

`newApp` 函数是 Kratos 应用的核心构造函数，负责组装所有组件并创建 `kratos.App` 实例。本文档说明 `newApp` 函数应该接收哪些类型的参数，以及如何组织这些参数。

## 参数分类

### 1. 核心必需参数

这些参数是大多数 Kratos 应用都需要的：

#### 1.1 日志记录器
```go
logger log.Logger
```
- **类型**: `log.Logger`
- **用途**: 应用日志记录
- **必需性**: ✅ 必需
- **说明**: 所有应用都需要日志功能

#### 1.2 传输层服务器
```go
gs *grpc.Server  // gRPC 服务器
hs *http.Server  // HTTP 服务器
```
- **类型**: `*grpc.Server`, `*http.Server`
- **用途**: 提供 gRPC 和/或 HTTP 服务
- **必需性**: ⚠️ 至少需要一个
- **说明**: 
  - 可以只提供 gRPC 服务器
  - 可以只提供 HTTP 服务器
  - 也可以同时提供两者
  - 使用可变参数或指针（nil 表示未启用）

### 2. 基础设施参数（可选）

这些参数用于增强应用功能，但不是必需的：

#### 2.1 服务注册中心
```go
registrar registry.Registrar
```
- **类型**: `registry.Registrar` (接口)
- **用途**: 服务注册与发现
- **必需性**: ❌ 可选
- **说明**: 
  - 微服务架构中需要
  - 单体应用通常不需要
  - 可以为 `nil`，需要做空值检查

#### 2.2 服务发现
```go
discovery registry.Discovery
```
- **类型**: `registry.Discovery` (接口)
- **用途**: 客户端服务发现
- **必需性**: ❌ 可选
- **说明**: 
  - 通常用于客户端应用
  - 服务端应用通常不需要
  - 如果应用需要调用其他服务，则需要

### 3. 配置参数（可选）

如果需要动态配置或运行时调整：

#### 3.1 服务器配置
```go
serverConf *conf.Server
```
- **类型**: `*conf.Server`
- **用途**: 服务器配置信息
- **必需性**: ❌ 可选
- **说明**: 
  - 通常配置在创建 Server 时已使用
  - 如果需要在 `newApp` 中使用配置，可以传入
  - 大多数情况下不需要

### 4. 自定义组件（可选）

根据项目需求添加的自定义组件：

#### 4.1 生命周期钩子
```go
lifecycle *LifecycleManager
```
- **类型**: 自定义类型
- **用途**: 管理应用生命周期事件
- **必需性**: ❌ 可选

#### 4.2 监控组件
```go
metrics *MetricsCollector
```
- **类型**: 自定义类型
- **用途**: 收集应用指标
- **必需性**: ❌ 可选

## 参数组织原则

### 原则 1: 必需参数在前，可选参数在后

```go
// ✅ 好的设计
func newApp(
    logger log.Logger,           // 必需
    gs *grpc.Server,             // 必需（至少一个）
    hs *http.Server,             // 必需（至少一个）
    registrar registry.Registrar, // 可选
) *kratos.App

// ❌ 不好的设计
func newApp(
    registrar registry.Registrar, // 可选参数不应该在前面
    logger log.Logger,
    gs *grpc.Server,
    hs *http.Server,
) *kratos.App
```

### 原则 2: 相关参数分组

```go
// ✅ 好的设计：传输层参数分组
func newApp(
    logger log.Logger,
    // 传输层
    gs *grpc.Server,
    hs *http.Server,
    // 服务注册与发现
    registrar registry.Registrar,
    discovery registry.Discovery,
) *kratos.App
```

### 原则 3: 使用接口类型，而非具体实现

```go
// ✅ 好的设计：使用接口
func newApp(
    registrar registry.Registrar,  // 接口类型
) *kratos.App

// ❌ 不好的设计：使用具体类型
func newApp(
    registrar *consul.Registry,  // 具体实现类型
) *kratos.App
```

### 原则 4: 可选参数使用指针或接口，支持 nil

```go
// ✅ 好的设计：支持 nil
func newApp(
    logger log.Logger,
    gs *grpc.Server,
    hs *http.Server,
    registrar registry.Registrar,  // 可以为 nil
) *kratos.App {
    opts := []kratos.Option{
        kratos.Logger(logger),
        kratos.Server(gs, hs),
    }
    if registrar != nil {  // 空值检查
        opts = append(opts, kratos.Registrar(registrar))
    }
    return kratos.New(opts...)
}
```

## 常见模式

### 模式 1: 最小化参数（推荐）

适用于大多数标准应用：

```go
func newApp(
    logger log.Logger,
    gs *grpc.Server,
    hs *http.Server,
) *kratos.App {
    return kratos.New(
        kratos.ID(id),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Logger(logger),
        kratos.Server(gs, hs),
    )
}
```

### 模式 2: 带服务注册

适用于微服务架构：

```go
func newApp(
    logger log.Logger,
    gs *grpc.Server,
    hs *http.Server,
    registrar registry.Registrar,
) *kratos.App {
    opts := []kratos.Option{
        kratos.ID(id),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Logger(logger),
        kratos.Server(gs, hs),
    }
    if registrar != nil {
        opts = append(opts, kratos.Registrar(registrar))
    }
    return kratos.New(opts...)
}
```

### 模式 3: 仅 gRPC 服务

适用于纯 gRPC 服务：

```go
func newApp(
    logger log.Logger,
    gs *grpc.Server,
    registrar registry.Registrar,
) *kratos.App {
    opts := []kratos.Option{
        kratos.ID(id),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Logger(logger),
        kratos.Server(gs),  // 只传 gRPC
    }
    if registrar != nil {
        opts = append(opts, kratos.Registrar(registrar))
    }
    return kratos.New(opts...)
}
```

### 模式 4: 仅 HTTP 服务

适用于纯 HTTP/REST API 服务：

```go
func newApp(
    logger log.Logger,
    hs *http.Server,
    registrar registry.Registrar,
) *kratos.App {
    opts := []kratos.Option{
        kratos.ID(id),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Logger(logger),
        kratos.Server(hs),  // 只传 HTTP
    }
    if registrar != nil {
        opts = append(opts, kratos.Registrar(registrar))
    }
    return kratos.New(opts...)
}
```

### 模式 5: 完整配置（高级）

适用于需要更多控制的应用：

```go
func newApp(
    logger log.Logger,
    gs *grpc.Server,
    hs *http.Server,
    registrar registry.Registrar,
    discovery registry.Discovery,
    serverConf *conf.Server,
) *kratos.App {
    opts := []kratos.Option{
        kratos.ID(id),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Logger(logger),
        kratos.Server(gs, hs),
    }
    
    // 添加元数据
    metadata := map[string]string{
        "env": "production",
    }
    if serverConf != nil {
        // 从配置中提取元数据
    }
    opts = append(opts, kratos.Metadata(metadata))
    
    // 注册服务
    if registrar != nil {
        opts = append(opts, kratos.Registrar(registrar))
    }
    
    return kratos.New(opts...)
}
```

## 参数类型总结表

| 参数类型 | 类型定义 | 必需性 | 用途 | 示例 |
|---------|---------|--------|------|------|
| Logger | `log.Logger` | ✅ 必需 | 日志记录 | `logger log.Logger` |
| gRPC Server | `*grpc.Server` | ⚠️ 至少一个 | gRPC 服务 | `gs *grpc.Server` |
| HTTP Server | `*http.Server` | ⚠️ 至少一个 | HTTP 服务 | `hs *http.Server` |
| Registrar | `registry.Registrar` | ❌ 可选 | 服务注册 | `registrar registry.Registrar` |
| Discovery | `registry.Discovery` | ❌ 可选 | 服务发现 | `discovery registry.Discovery` |
| Server Config | `*conf.Server` | ❌ 可选 | 服务器配置 | `serverConf *conf.Server` |

## 最佳实践

### ✅ 推荐做法

1. **保持参数列表简洁**
   - 只包含真正需要的参数
   - 避免传递不必要的配置对象

2. **使用接口类型**
   - 提高灵活性
   - 便于测试和替换实现

3. **可选参数支持 nil**
   - 使用空值检查
   - 条件性添加 Option

4. **参数顺序合理**
   - 必需参数在前
   - 可选参数在后
   - 相关参数分组

5. **使用 Wire 自动注入**
   - 通过 Wire 管理依赖
   - 减少手动组装代码

### ❌ 避免的做法

1. **传递过多参数**
   ```go
   // ❌ 不好：参数过多
   func newApp(
       logger log.Logger,
       gs *grpc.Server,
       hs *http.Server,
       registrar registry.Registrar,
       discovery registry.Discovery,
       config *conf.Bootstrap,
       metrics *Metrics,
       tracer *Tracer,
       // ... 太多参数
   ) *kratos.App
   ```

2. **传递具体实现类型**
   ```go
   // ❌ 不好：使用具体类型
   func newApp(
       registrar *consul.Registry,  // 应该使用接口
   ) *kratos.App
   ```

3. **忽略空值检查**
   ```go
   // ❌ 不好：没有空值检查
   opts = append(opts, kratos.Registrar(registrar))  // registrar 可能为 nil
   ```

4. **在 newApp 中创建组件**
   ```go
   // ❌ 不好：在 newApp 中创建组件
   func newApp(logger log.Logger) *kratos.App {
       gs := grpc.NewServer(...)  // 应该在外部创建
       // ...
   }
   ```

## 当前项目实现

当前项目的 `newApp` 实现：

```go
func newApp(
    logger log.Logger,           // ✅ 必需参数
    gs *grpc.Server,             // ✅ 传输层
    hs *http.Server,             // ✅ 传输层
    registrar registry.Registrar, // ✅ 可选基础设施
) *kratos.App {
    opts := []kratos.Option{
        kratos.ID(id),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Metadata(map[string]string{}),
        kratos.Logger(logger),
        kratos.Server(gs, hs),
    }
    // ✅ 空值检查
    if registrar != nil {
        opts = append(opts, kratos.Registrar(registrar))
    }
    return kratos.New(opts...)
}
```

**评价**: ✅ 符合最佳实践
- 参数顺序合理
- 使用接口类型
- 支持可选参数
- 有适当的空值检查

## 总结

`newApp` 函数的参数设计应该遵循以下原则：

1. **必需参数**: `log.Logger` + 至少一个传输层服务器
2. **可选参数**: 根据项目需求添加（如 `registry.Registrar`）
3. **参数顺序**: 必需在前，可选在后
4. **类型选择**: 优先使用接口类型
5. **空值处理**: 可选参数需要空值检查

保持参数列表简洁，只包含真正需要的组件，通过 Wire 自动管理依赖注入。

