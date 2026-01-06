# 第三方服务集成指南索引

## 概述

本系列文档提供分步骤的第三方服务集成指南，帮助开发者按照最佳实践集成 gRPC 和 HTTP REST API 第三方服务。

## 文档列表

### 第一步：准备工作
**文档路径：** [third-party-integration-01-preparation.md](./third-party-integration-01-preparation.md)

**内容：**
- 扩展配置定义（conf.proto）
- 创建目录结构
- 确定服务类型
- 安装依赖

**适用场景：** 开始集成任何第三方服务前，都需要完成准备工作。

---

### 第二步：gRPC 服务集成
**文档路径：** [third-party-integration-02-grpc.md](./third-party-integration-02-grpc.md)

**内容：**
- 定义 Proto 文件
- 创建 gRPC 客户端管理器
- 在 Repository 中使用 gRPC 客户端
- 高级配置（TLS、服务发现）

**适用场景：** 集成内部微服务或外部 gRPC 服务。

---

### 第三步：HTTP REST API 集成
**文档路径：** [third-party-integration-03-http.md](./third-party-integration-03-http.md)

**内容：**
- 定义请求/响应类型
- 创建 HTTP 客户端
- 创建业务封装
- 在 Repository 中使用

**适用场景：** 集成 HTTP REST API 第三方服务。

---

### 第四步：在 Repository 中使用第三方服务
**文档路径：** [third-party-integration-04-usage.md](./third-party-integration-04-usage.md)

**内容：**
- 获取客户端
- 类型转换
- 错误处理
- 组合多个服务
- 缓存策略
- 超时和重试
- 监控和日志

**适用场景：** 在业务代码中正确使用已集成的第三方服务。

---

## 快速开始

### 集成 gRPC 服务

1. 阅读 [第一步：准备工作](./third-party-integration-01-preparation.md)
2. 按照 [第二步：gRPC 服务集成](./third-party-integration-02-grpc.md) 完成集成
3. 参考 [第四步：在 Repository 中使用第三方服务](./third-party-integration-04-usage.md) 使用服务

### 集成 HTTP REST API

1. 阅读 [第一步：准备工作](./third-party-integration-01-preparation.md)
2. 按照 [第三步：HTTP REST API 集成](./third-party-integration-03-http.md) 完成集成
3. 参考 [第四步：在 Repository 中使用第三方服务](./third-party-integration-04-usage.md) 使用服务

## 相关文档

- [第三方服务接口定义最佳实践](../architecture/third-party-api-definitions.md) - 了解如何组织和管理接口定义
- [数据访问层外部依赖管理](../architecture/data-layer-dependencies.md) - 了解 Data 层的依赖管理原则
- [依赖注入](../architecture/dependency-injection.md) - 了解 Wire 依赖注入的使用

## 常见问题

### Q: 我应该先看哪个文档？

A: 如果是第一次集成第三方服务，建议按顺序阅读：
1. 第一步：准备工作
2. 根据服务类型选择第二步（gRPC）或第三步（HTTP）
3. 第四步：在 Repository 中使用

### Q: 如何选择 gRPC 还是 HTTP？

A: 
- **gRPC**：适用于内部微服务、需要高性能的场景
- **HTTP REST API**：适用于外部第三方服务、需要简单集成的场景

### Q: 配置应该放在哪里？

A: 第三方服务的连接信息应该放在 `configs/config.yaml` 中，通过 `internal/conf/conf.proto` 定义配置结构。

### Q: 如何测试第三方服务集成？

A: 使用 mock 来模拟第三方服务，参考各步骤文档中的测试部分。

## 贡献

如果发现文档有误或需要补充，请提交 Issue 或 Pull Request。

