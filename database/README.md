# 数据库学习目录

## 数据库基础
- [ ] 数据库概念
  - [ ] 数据库类型
  - [ ] ACID特性
  - [ ] 事务管理
  - [ ] 并发控制
- [ ] 数据建模
  - [ ] ER模型
  - [ ] 关系模型
  - [ ] 范式理论
  - [ ] 反范式化
- [ ] SQL基础
  - [ ] DDL (数据定义语言)
  - [ ] DML (数据操作语言)
  - [ ] DQL (数据查询语言)
  - [ ] DCL (数据控制语言)

## 关系型数据库
- [ ] MySQL
  - [ ] 安装和配置
  - [ ] 存储引擎 (InnoDB, MyISAM)
  - [ ] 索引优化
  - [ ] 查询优化
  - [ ] 主从复制
  - [ ] 分库分表
- [x] PostgreSQL
  - [x] [PostgreSQL 学习指南](./postgresql-learning.md) - 完整的 PostgreSQL 学习资料
  - [x] [PostgreSQL 基础](./postgresql-basics.md) - 安装、配置、基本操作
  - [x] [PostgreSQL 数据类型](./postgresql-data-types.md) - 基础类型、高级类型、JSON、数组
  - [x] [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md) - 索引类型、查询计划、性能调优
  - [x] [PostgreSQL 事务与并发控制](./postgresql-transactions.md) - ACID、隔离级别、锁机制、MVCC
  - [x] [PostgreSQL 存储过程与函数](./postgresql-functions.md) - PL/pgSQL、触发器、自定义函数
  - [x] [PostgreSQL 全文搜索](./postgresql-fulltext-search.md) - 全文索引、搜索配置、多语言支持
  - [x] [PostgreSQL 扩展与插件](./postgresql-extensions.md) - 常用扩展、PostGIS、pg_stat_statements
  - [x] [PostgreSQL 管理与运维](./postgresql-admin.md) - 用户权限、备份恢复、监控、日志管理
  - [x] [PostgreSQL 复制与高可用](./postgresql-replication.md) - 主从复制、流复制、高可用方案
  - [x] [PostgreSQL 最佳实践](./postgresql-best-practices.md) - 设计原则、性能优化、安全实践
  - [x] [PostgreSQL 常见问题与解决方案](./postgresql-troubleshooting.md) - 常见错误、性能问题、故障排查
  - [x] [PostgreSQL vs MySQL 全面对比](./postgresql-vs-mysql.md) - 与 MySQL 的详细对比，便于对比学习
- [ ] Oracle
  - [ ] 体系结构
  - [ ] PL/SQL
  - [ ] 性能调优
  - [ ] RAC集群
- [ ] SQL Server
  - [ ] T-SQL
  - [ ] 集成服务
  - [ ] 分析服务
  - [ ] 报表服务

## NoSQL数据库
- [ ] 键值存储
  - [ ] Redis
    - [ ] 数据类型
    - [ ] 持久化
    - [ ] 集群模式
    - [ ] 缓存策略
  - [ ] Memcached
- [ ] 文档数据库
  - [ ] MongoDB
    - [ ] 文档模型
    - [ ] 聚合管道
    - [ ] 索引策略
    - [ ] 分片集群
  - [ ] CouchDB
- [ ] 列族数据库
  - [ ] Cassandra
    - [ ] 数据模型
    - [ ] 一致性级别
    - [ ] 分区策略
  - [ ] HBase
- [ ] 图数据库
  - [ ] Neo4j
    - [ ] Cypher查询语言
    - [ ] 图算法
    - [ ] 性能优化
  - [ ] ArangoDB

## 数据仓库
- [ ] 数据仓库概念
  - [ ] OLTP vs OLAP
  - [ ] 星型模式
  - [ ] 雪花模式
  - [ ] 事实表和维度表
- [ ] ETL/ELT
  - [ ] 数据抽取
  - [ ] 数据转换
  - [ ] 数据加载
  - [ ] 数据质量
- [ ] 数据湖
  - [ ] 存储格式 (Parquet, ORC)
  - [ ] 元数据管理
  - [ ] 数据目录

## 大数据技术
- [ ] Hadoop生态
  - [ ] HDFS
  - [ ] MapReduce
  - [ ] YARN
  - [ ] Hive
  - [ ] HBase
- [ ] Spark
  - [ ] RDD
  - [ ] DataFrame
  - [ ] Spark SQL
  - [ ] Spark Streaming
- [ ] Flink
  - [ ] 流处理
  - [ ] 批处理
  - [ ] 状态管理
  - [ ] 容错机制

## 数据库设计
- [ ] 需求分析
  - [ ] 业务需求收集
  - [ ] 数据需求分析
  - [ ] 性能需求
- [ ] 概念设计
  - [ ] 实体识别
  - [ ] 关系建模
  - [ ] 属性定义
- [ ] 逻辑设计
  - [ ] 表结构设计
  - [ ] 索引设计
  - [ ] 约束设计
- [ ] 物理设计
  - [ ] 存储优化
  - [ ] 分区策略
  - [ ] 备份恢复

## 性能优化
- [ ] 查询优化
  - [ ] 执行计划分析
  - [ ] 索引优化
  - [ ] SQL重写
  - [ ] 统计信息
- [ ] 系统优化
  - [ ] 内存配置
  - [ ] 磁盘I/O优化
  - [ ] 并发控制
  - [ ] 连接池管理
- [ ] 监控和诊断
  - [ ] 性能监控
  - [ ] 慢查询分析
  - [ ] 资源使用分析

## 高可用和容灾
- [ ] 主从复制
  - [ ] 同步复制
  - [ ] 异步复制
  - [ ] 半同步复制
- [ ] 集群技术
  - [ ] 读写分离
  - [ ] 负载均衡
  - [ ] 故障转移
- [ ] 备份恢复
  - [ ] 全量备份
  - [ ] 增量备份
  - [ ] 时间点恢复
  - [ ] 灾难恢复

## 数据安全
- [ ] 访问控制
  - [ ] 用户管理
  - [ ] 权限管理
  - [ ] 角色管理
- [ ] 数据加密
  - [ ] 传输加密
  - [ ] 存储加密
  - [ ] 字段级加密
- [ ] 审计和合规
  - [ ] 操作审计
  - [ ] 数据脱敏
  - [ ] 合规要求

## 云数据库
- [ ] AWS数据库服务
  - [ ] RDS
  - [ ] DynamoDB
  - [ ] Aurora
  - [ ] Redshift
- [ ] 阿里云数据库
  - [ ] RDS
  - [ ] PolarDB
  - [ ] AnalyticDB
  - [ ] MaxCompute
- [ ] 腾讯云数据库
  - [ ] TencentDB
  - [ ] TDSQL
  - [ ] 云数据仓库

## 数据库工具和框架
- [ ] 管理工具
  - [ ] phpMyAdmin
  - [ ] Navicat
  - [ ] DBeaver
  - [ ] DataGrip
- [ ] 监控工具
  - [ ] Prometheus
  - [ ] Grafana
  - [ ] Zabbix
- [ ] 迁移工具
  - [ ] [数据库迁移管理工具与实践](./migration-tools-practices.md) - 全面的迁移工具对比和实践指南
  - [ ] Flyway
  - [ ] Liquibase
  - [ ] Alembic (Python/SQLAlchemy)
  - [ ] TypeORM Migrations
  - [ ] Ent Migrations (Go)
  - [ ] Prisma Migrate
  - [ ] 数据同步工具
- [ ] ORM 和框架
  - [ ] [Ent、Schema 和 GraphQL](./ent-graphql-schema.md) - Ent 框架与 GraphQL 集成指南

## 实际项目
- [ ] 电商数据库设计
- [ ] 用户行为分析
- [ ] 实时数据处理
- [ ] 数据迁移项目
- [ ] 性能优化项目 