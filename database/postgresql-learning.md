# PostgreSQL 学习指南

## 目录

### 基础入门
- [PostgreSQL 基础](./postgresql-basics.md) - 安装、配置、基本操作、SQL基础
- [PostgreSQL 数据类型](./postgresql-data-types.md) - 基础类型、高级类型、JSON、数组、自定义类型

### 核心功能
- [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md) - 索引类型、查询计划、性能调优
- [PostgreSQL 事务与并发控制](./postgresql-transactions.md) - ACID、隔离级别、锁机制、MVCC

### 高级特性
- [PostgreSQL 存储过程与函数](./postgresql-functions.md) - PL/pgSQL、触发器、自定义函数
- [PostgreSQL 全文搜索](./postgresql-fulltext-search.md) - 全文索引、搜索配置、多语言支持
- [PostgreSQL 扩展与插件](./postgresql-extensions.md) - 常用扩展、PostGIS、pg_stat_statements

### 管理与运维
- [PostgreSQL 管理与运维](./postgresql-admin.md) - 用户权限、备份恢复、监控、日志管理
- [PostgreSQL 复制与高可用](./postgresql-replication.md) - 主从复制、流复制、高可用方案

### 实践应用
- [PostgreSQL 最佳实践](./postgresql-best-practices.md) - 设计原则、性能优化、安全实践
- [PostgreSQL 常见问题与解决方案](./postgresql-troubleshooting.md) - 常见错误、性能问题、故障排查
- [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md) - 与 MySQL 的详细对比，便于对比学习

---

## 学习路径建议

### 初学者路径
1. **第一步**：阅读 [PostgreSQL 基础](./postgresql-basics.md)
2. **第二步**：学习 [PostgreSQL 数据类型](./postgresql-data-types.md)
3. **第三步**：掌握 [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md)
4. **第四步**：了解 [PostgreSQL 事务与并发控制](./postgresql-transactions.md)

### 进阶路径
1. **存储过程开发**：[PostgreSQL 存储过程与函数](./postgresql-functions.md)
2. **高级特性**：[PostgreSQL 全文搜索](./postgresql-fulltext-search.md)
3. **扩展使用**：[PostgreSQL 扩展与插件](./postgresql-extensions.md)

### 运维路径
1. **日常管理**：[PostgreSQL 管理与运维](./postgresql-admin.md)
2. **高可用**：[PostgreSQL 复制与高可用](./postgresql-replication.md)
3. **问题排查**：[PostgreSQL 常见问题与解决方案](./postgresql-troubleshooting.md)

---

## 快速参考

### 常用命令
```bash
# 连接数据库
psql -U username -d database_name

# 查看版本
psql --version

# 列出所有数据库
psql -l

# 备份数据库
pg_dump -U username database_name > backup.sql

# 恢复数据库
psql -U username database_name < backup.sql
```

### 官方资源
- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [PostgreSQL 中文社区](https://www.postgresql.org/)
- [PostgreSQL Wiki](https://wiki.postgresql.org/)

### 推荐工具
- **客户端工具**：pgAdmin、DBeaver、DataGrip
- **命令行工具**：psql
- **监控工具**：pg_stat_statements、pgBadger
- **迁移工具**：pg_dump、pg_restore、pg_upgrade

---

## 版本说明

本文档主要基于 **PostgreSQL 14+** 版本编写，大部分内容适用于 PostgreSQL 12+ 版本。

---

*最后更新：2024年*

