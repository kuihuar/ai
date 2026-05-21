# 面试口述题索引（电信 / AI Agent 后端）

本目录面试题依据 [`readme.md`](./readme.md) 中的岗位要求整理，按主题拆成多个文件，便于分块准备。

| 文件 | 内容侧重 |
|------|-----------|
| [interview-go.md](./interview-go.md) | Go 核心、并发、内存模型、GC、errgroup、singleflight、context、pprof / trace / race |
| [interview-web.md](./interview-web.md) | Hertz / Gin、REST、WebSocket、gRPC / Protobuf、网关与流式 |
| [interview-data.md](./interview-data.md) | MySQL、MongoDB、GORM（Hook、预加载、事务）、慢 SQL |
| [mongodb-gorm-guide.md](./mongodb-gorm-guide.md) | **MongoDB + GORM 详解**：基本操作、聚合管道、Hook / Preload / 事务与分工 |
| [interview-cache-mq.md](./interview-cache-mq.md) | Redis、Kafka、Asynq、一致性与幂等 |
| [interview-infra-observability.md](./interview-infra-observability.md) | Docker、K8s、对象存储与预签名、Prometheus、OpenTelemetry、日志 |
| [interview-engineering-security.md](./interview-engineering-security.md) | Git、JWT、API Key、限流、Code Review、协作与交付 |
| [interview-ai-agent-bonus.md](./interview-ai-agent-bonus.md) | Agent 平台、RAG、向量库、加分项中间件与诚实边界 |

## 使用建议

- 每题回答结构：**结论一句 + 项目例子 + 指标或取舍**。
- 高级岗多追问「为什么这样设计、失败时怎么办、如何观测与回滚」。
- 与 JD 不完全重合的技术（如 Eino、WASM），无经验则说明理解边界与学习计划即可。
