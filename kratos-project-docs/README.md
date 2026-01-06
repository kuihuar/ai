# 软件工程最佳实践文档

本项目是一个基于 Kratos 框架的研究项目，专注于探索和记录软件工程领域的最佳实践。

## 文档结构

### 架构设计
- [Kratos 架构实践](./architecture/kratos-architecture.md) - Kratos 框架的架构设计模式
- [分层架构设计](./architecture/layered-architecture.md) - 业务、数据、服务层的设计原则
- [依赖注入实践](./architecture/dependency-injection.md) - Wire 依赖注入的使用模式
- [多应用支持](./architecture/multi-app.md) - Kratos 多应用架构实践
- [服务注册与发现](./architecture/service-registry-discovery.md) - 微服务注册与发现的实现指南
- [数据访问层依赖管理](./architecture/data-layer-dependencies.md) - Redis、Kafka、MQ、RPC/HTTP 等外部依赖的组织和管理
- [第三方服务接口定义](./architecture/third-party-api-definitions.md) - 第三方服务 request/response 类型定义的最佳实践

### 代码规范
- [Go 代码规范](./code-standards/go-standards.md) - Go 语言编码规范和最佳实践
- [API 设计规范](./code-standards/api-design.md) - RESTful 和 gRPC API 设计原则
- [错误处理](./code-standards/error-handling.md) - 错误处理和异常管理

### 开发流程
- [开发工作流](./development/workflow.md) - 从需求到部署的完整流程
- [代码审查](./development/code-review.md) - 代码审查的最佳实践
- [测试策略](./development/testing-strategy.md) - 单元测试、集成测试策略
- [Cursor IDE 配置](./development/cursor-config.md) - Cursor IDE 配置和多设备同步
- [Worker 集成最佳实践](./development/worker-integration.md) - 在主应用中集成 daemon-worker 或 cron-worker 的最佳实践
- [Table Consumer Daemon Ants 架构设计](./development/daemon-ants-architecture.md) - table_consumer_ants.go 和 table_consumer_ants_biz.go 的分工和设计原则

### 运维实践
- [配置管理](./operations/config-management.md) - 配置文件的组织和管理
- [日志规范](./operations/logging.md) - 日志记录的最佳实践
- [监控与可观测性](./operations/observability.md) - 指标、追踪、日志的实践

### 项目管理
- [项目结构](./project/structure.md) - 项目目录组织规范
- [文档编写](./project/documentation.md) - 技术文档编写指南
- [版本管理](./project/versioning.md) - 版本控制和发布流程

## 贡献指南

欢迎添加和更新最佳实践文档。请遵循以下原则：

1. **实践导向**：基于实际项目经验，而非理论空谈
2. **示例丰富**：提供具体的代码示例和场景说明
3. **持续更新**：随着项目演进，不断更新和完善文档

