# 中间件组织方式最佳实践

## 当前状态

项目中的中间件目前采用混合组织方式：

```
internal/
├── middleware/              # 通用中间件
│   ├── auth.go             # 认证中间件
│   ├── ratelimit.go        # 限流中间件
│   ├── chain.go            # 中间件链
│   └── router.go           # 路由管理器
├── tracing/
│   └── middleware.go       # Tracing 中间件
└── metrics/
    └── middleware.go       # Metrics 中间件
```

## 组织原则

### 原则 1：功能相关中间件放在功能包中

**适用于**：与特定功能模块紧密相关的中间件

**示例**：
- ✅ `internal/tracing/middleware.go` - Tracing 中间件属于可观测性功能
- ✅ `internal/metrics/middleware.go` - Metrics 中间件属于可观测性功能

**理由**：
1. **高内聚**：功能相关的代码聚合在一起，便于维护和理解
2. **低耦合**：功能模块可以独立开发和测试
3. **自包含**：功能包包含完整的实现（初始化、配置、中间件）

### 原则 2：通用中间件放在 `internal/middleware/`

**适用于**：跨功能、通用的中间件

**示例**：
- ✅ `internal/middleware/auth.go` - 认证中间件，多个功能模块都会使用
- ✅ `internal/middleware/ratelimit.go` - 限流中间件，通用功能
- ✅ `internal/middleware/chain.go` - 中间件链工具，通用基础设施

**理由**：
1. **通用性**：这些中间件不依赖特定功能模块
2. **复用性**：多个服务或功能都会使用
3. **基础设施**：属于框架层面的基础设施

## 决策树

使用以下决策树来确定中间件的放置位置：

```
中间件是否需要特定功能包的依赖？
├─ 是 → 放在功能包中（如 tracing/middleware.go）
└─ 否 → 放在 internal/middleware/（如 auth.go）
```

### 具体判断标准

#### 放在功能包中（如 `internal/tracing/`）

如果中间件满足以下任一条件：
- ✅ 依赖特定功能包的配置或初始化（如需要 TracerProvider）
- ✅ 是功能包的核心组成部分（如 tracing 是 OpenTelemetry 的一部分）
- ✅ 与功能包的其他代码紧密耦合（如使用相同的配置结构）

**示例**：
```go
// internal/tracing/middleware.go
// 依赖 tracing 包的配置和初始化
func ServerWithConfig(config TracingConfig) middleware.Middleware {
    // 使用 tracing 包的配置
}
```

#### 放在 `internal/middleware/`

如果中间件满足以下条件：
- ✅ 不依赖特定功能包
- ✅ 可以被多个功能模块使用
- ✅ 是通用的基础设施（如认证、限流、日志）

**示例**：
```go
// internal/middleware/auth.go
// 通用认证中间件，不依赖特定功能
func Auth(validator TokenValidator, logger log.Logger) middleware.Middleware {
    // 通用实现
}
```

## 最佳实践

### 1. 保持一致性

一旦确定了组织方式，应该在整个项目中保持一致：
- 所有可观测性相关的中间件放在各自的功能包中
- 所有通用中间件放在 `internal/middleware/`

### 2. 文档说明

在功能包的 README 中说明中间件的位置和使用方式：

```markdown
# Tracing 包

## 中间件

本包提供了 tracing 中间件，位于 `middleware.go`：
- `Server()` - 默认配置的 tracing 中间件
- `ServerWithConfig()` - 可配置的 tracing 中间件
```

### 3. 导入路径

使用清晰的导入路径：

```go
// 通用中间件
import "sre/internal/middleware"

// 功能特定中间件
import "sre/internal/tracing"
import "sre/internal/metrics"
```

### 4. 避免循环依赖

确保中间件不会造成循环依赖：
- 功能包的中间件不应该依赖 `internal/middleware/`
- `internal/middleware/` 不应该依赖业务功能包

## 迁移建议

如果发现中间件放错了位置，可以按以下步骤迁移：

### 从 `internal/middleware/` 迁移到功能包

1. 在功能包中创建 `middleware.go`
2. 移动中间件代码
3. 更新导入路径
4. 更新文档

### 从功能包迁移到 `internal/middleware/`

1. 确保中间件不依赖功能包特定代码
2. 移动到 `internal/middleware/`
3. 更新导入路径
4. 更新文档

## 当前项目的组织方式（推荐）

基于当前项目的架构原则，推荐采用以下组织方式：

### ✅ 推荐的组织方式

```
internal/
├── middleware/              # 通用中间件
│   ├── auth.go             # 认证中间件（通用）
│   ├── ratelimit.go        # 限流中间件（通用）
│   ├── chain.go            # 中间件链（基础设施）
│   └── router.go           # 路由管理器（基础设施）
├── tracing/
│   ├── provider.go         # Tracing 提供者
│   ├── helpers.go          # Tracing 辅助函数
│   └── middleware.go       # Tracing 中间件（功能特定）
└── metrics/
    ├── provider.go         # Metrics 提供者
    └── middleware.go       # Metrics 中间件（功能特定）
```

### 理由

1. **符合 Clean Architecture**：功能模块自包含
2. **便于维护**：功能相关的代码聚合在一起
3. **清晰职责**：通用中间件和功能特定中间件分离
4. **易于扩展**：新增功能时，中间件自然放在功能包中

## 总结

**推荐做法**：
- ✅ **功能特定中间件** → 放在功能包中（如 `tracing/middleware.go`）
- ✅ **通用中间件** → 放在 `internal/middleware/`（如 `auth.go`）

**当前项目的组织方式是正确的**，符合 Clean Architecture 和关注点分离原则。

