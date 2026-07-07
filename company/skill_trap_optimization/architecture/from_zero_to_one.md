# Go 项目从 0 到 1 架构设计

> 各模块详细拆解见同目录下的专题文件，本文档为总览索引。

## 专题文件

| 文件 | 覆盖内容 |
|------|----------|
| [project-layout-and-di.md](project-layout-and-di.md) | 目录结构、分层架构、依赖注入（手动/Wire/Fx 对比） |
| [middleware-selection.md](middleware-selection.md) | Web 框架选型（Gin/Chi/Fiber/Echo）、中间件逐项对比与选型 |
| [database-and-config.md](database-and-config.md) | 配置管理方案对比、ORM 选型、数据库迁移、缓存策略 |
| [observability-and-deploy.md](observability-and-deploy.md) | 日志/指标/追踪、API 设计（REST/gRPC/GraphQL）、测试策略、Docker/K8s、微服务演进、安全 |

## 核心设计决策速查

### 分层架构
```
Handler → Service → Repository → Model
```
上层依赖下层，下层定义接口（依赖反转），方便测试 mock 和替换实现。

### 依赖注入
| 方案 | 适用 | 关键取舍 |
|------|------|----------|
| 手动构造函数 | 小项目 < 20 依赖 | 零依赖、最直接，但依赖多时 main.go 膨胀 |
| Google Wire | 中大型项目 | 编译时生成、无运行时开销、报错清晰 |
| Uber Fx | 大型/微服务 | 运行时 DI、生命周期管理、功能最全但引入复杂度 |

### Web 框架
| 框架 | 适用 | 性能 | 生态 |
|------|------|------|------|
| Gin | 最通用，社区最大 | 中 | 成熟 |
| Chi | 偏好标准库风格 | 中 | 轻量 |
| Fiber | 追求极致性能 | 高（fasthttp） | 不兼容 net/http |
| Echo | 功能全面 | 中 | 成熟 |

### 数据库
| 方案 | 适用 | 关键取舍 |
|------|------|----------|
| 原生 sql + sqlx | 复杂 SQL、多表 join | 完全控制，手写 SQL 多 |
| GORM | 常规 CRUD 为主 | 开发快，复杂查询难控 |
| Ent (Facebook) | 数据模型复杂、GraphQL 集成 | 代码生成、类型安全、学习成本高 |
| Bun | 追求性能和灵活性平衡 | 介于 ORM 和原生之间 |

### 配置管理
| 方案 | 适用 |
|------|------|
| os.Getenv + 默认值 | 简单服务、Docker 部署 |
| Viper | 需多格式/多来源/热更新的项目 |
| envconfig / cleanenv | 仅需环境变量映射到 struct |

### 日志库
| 方案 | 适用 | 特点 |
|------|------|------|
| slog (Go 1.21+) | 标准库首选，零依赖 | 结构化、plugin 扩展 |
| zap (Uber) | 极致性能 | 零分配、高性能 |
| zerolog | 同样极致性能 | JSON 原生，API 链式调用 |
| logrus | 老项目 | 功能全但慢，不再推荐新项目 |

### API 风格
| 风格 | 适用 |
|------|------|
| RESTful | 对外 API、跨语言调用、前后端分离 |
| gRPC | 内部微服务间高性能调用 |
| GraphQL | 前端需要灵活数据查询 |

### 中间件链（洋葱模型）
```
Recovery → Logger → CORS → Auth → RateLimit → Tracing → Handler
```
