# Canal 缓存一致性方案文档索引

## 文档概览

本文档索引提供了基于 Canal 架构的缓存一致性方案的完整文档导航。

## 文档列表

### 1. [方案概述](./canal-cache-consistency-overview.md)

**内容**：
- Canal 简介
- 为什么使用 Canal
- 架构对比
- 整体架构
- 核心流程
- 适用场景
- 实施建议

**适合人群**：
- 架构师
- 技术负责人
- 需要了解整体方案的开发人员

### 2. [架构设计](./canal-cache-consistency-architecture.md)

**内容**：
- 架构组件详解
- 数据流设计
- 缓存同步策略
- 错误处理
- 性能优化
- 一致性保证
- 部署架构
- 安全考虑

**适合人群**：
- 架构师
- 高级开发人员
- 需要深入理解架构的设计人员

### 3. [实现指南](./canal-cache-consistency-implementation.md)

**内容**：
- 前置条件
- 代码实现
- 配置管理
- 测试验证
- 部署步骤
- 故障排查

**适合人群**：
- 开发人员
- 需要具体实现代码的工程师

### 4. [配置说明](./canal-cache-consistency-config.md)

**内容**：
- Canal Server 配置
- Canal Client 配置
- 缓存键映射配置
- 环境变量配置
- 配置验证
- 配置最佳实践

**适合人群**：
- 运维人员
- 开发人员
- 需要配置系统的工程师

### 5. [监控运维](./canal-cache-consistency-monitoring.md)

**内容**：
- 监控指标
- 告警规则
- 日志管理
- 性能优化
- 故障处理
- 运维脚本
- 备份和恢复

**适合人群**：
- 运维人员
- SRE 工程师
- 需要监控和运维系统的工程师

## 快速开始

### 新手入门路径

1. **第一步**：阅读 [方案概述](./canal-cache-consistency-overview.md)
   - 了解 Canal 是什么
   - 理解为什么使用 Canal
   - 了解整体架构

2. **第二步**：阅读 [架构设计](./canal-cache-consistency-architecture.md)
   - 深入理解架构设计
   - 了解缓存同步策略
   - 理解一致性保证机制

3. **第三步**：阅读 [实现指南](./canal-cache-consistency-implementation.md)
   - 了解实现步骤
   - 参考代码示例
   - 进行实际开发

4. **第四步**：阅读 [配置说明](./canal-cache-consistency-config.md)
   - 配置 Canal Server
   - 配置 Canal Client
   - 验证配置正确性

5. **第五步**：阅读 [监控运维](./canal-cache-consistency-monitoring.md)
   - 设置监控指标
   - 配置告警规则
   - 建立运维流程

## 文档关系图

```
方案概述
  │
  ├─→ 架构设计 ──→ 实现指南
  │                    │
  │                    ├─→ 配置说明
  │                    │
  │                    └─→ 监控运维
  │
  └─→ 当前方案（Cache-Aside）
```

## 相关文档

### 当前方案文档

- [缓存一致性策略](./cache-consistency-strategy.md) - Cache-Aside 模式实现

### 架构文档

- [分布式事件一致性架构设计](../architecture/distributed-event-consistency.md) - Outbox + Saga 模式

## 常见问题

### Q1: Canal 和 Cache-Aside 有什么区别？

**A**: 
- **Canal**：基于 Binlog 的自动同步，业务代码无需关心缓存删除
- **Cache-Aside**：业务代码显式删除缓存，实现简单但需要业务代码配合

详见 [方案概述](./canal-cache-consistency-overview.md#架构对比)

### Q2: 什么时候使用 Canal，什么时候使用 Cache-Aside？

**A**: 
- **Canal**：适合高频热点数据、跨服务数据、实时性要求高的场景
- **Cache-Aside**：适合低频数据、计算型数据、临时数据

详见 [方案概述](./canal-cache-consistency-overview.md#适用场景)

### Q3: Canal 如何保证一致性？

**A**: 
- 基于 MySQL Binlog，保证所有数据变更都会被捕获
- 最终一致性，不保证强一致性
- 通常延迟 < 1 秒

详见 [架构设计](./canal-cache-consistency-architecture.md#一致性保证)

### Q4: Canal 的性能如何？

**A**: 
- 异步处理，对业务性能影响极小
- 支持批量处理，提高效率
- 支持过滤规则，只处理需要的表

详见 [架构设计](./canal-cache-consistency-architecture.md#性能优化)

### Q5: 如何监控 Canal？

**A**: 
- 监控 Canal Server 运行状态
- 监控 Binlog 延迟
- 监控事件处理速度
- 监控缓存删除成功率

详见 [监控运维](./canal-cache-consistency-monitoring.md#监控指标)

## 更新日志

- **2024-01-XX**: 初始版本，创建文档索引和所有子文档

## 贡献指南

如果发现文档问题或有改进建议，请：
1. 提交 Issue
2. 或直接提交 Pull Request

## 参考资源

- [Canal 官方文档](https://github.com/alibaba/canal)
- [Canal 快速开始](https://github.com/alibaba/canal/wiki/QuickStart)
- [MySQL Binlog 格式](https://dev.mysql.com/doc/refman/8.0/en/binary-log.html)

