# MongoDB 与 GORM 详解

> 对应 `readme.md` 中：**MongoDB**（基本操作与聚合）与 **GORM**（Hook、预加载、事务）。  
> 本文偏「原理 + 使用要点 + 踩坑」，可与 [interview-data.md](./interview-data.md) 口述题配合使用。

---

## 第一部分：MongoDB

### 1. 数据模型与核心概念

- **Database**：逻辑库，一般按业务或环境划分（如 `agent_prod`）。
- **Collection**：集合，约等于「无固定 schema 的表」；同一集合内文档结构可以不同，但生产环境通常会约定字段规范。
- **Document**：BSON 文档（JSON 的超集），支持嵌套文档与数组；有 `_id` 主键（默认 ObjectId）。

**和关系型数据库的直觉差异：**  
没有强制的「行结构一致」；多表关系靠**嵌套**或**引用**（存 `ObjectId`）表达；跨文档事务在 4.0+ 才成熟，设计时仍宜减少跨文档强一致依赖。

---

### 2. 基本操作（CRUD）要点

| 操作 | 含义 | 常见注意点 |
|------|------|------------|
| **insert / insertMany** | 插入 | 批量插入减少往返；考虑 **ordered: false** 时部分失败语义 |
| **find / findOne** | 查询 | 投影（projection）只取需要的字段，减少网络与反序列化 |
| **updateOne / updateMany** | 更新 | 区分 **$set**（改部分字段）与替换整文档；默认不会 upsert，需显式 `upsert: true` |
| **replaceOne** | 替换 | 整文档替换，易误删未出现在新文档里的字段 |
| **deleteOne / deleteMany** | 删除 | 一定先用带条件的 find 在预发验证影响行数 |

**查询运算符（面试常提）：**  
`$eq` `$ne` `$in` `$nin` `$gt` `$gte` `$lt` `$lte` `$regex` `$exists` `$elemMatch`（数组内多条件）等。

**索引与查询：**  
没有合适索引时，大集合上的 sort / 大范围 scan 会拖垮延迟；`explain("executionStats")` 看是否 COLLSCAN、是否用到 index。

---

### 3. 聚合管道（Aggregation Pipeline）

**心智模型：**  
数据像水流一样经过多个 **Stage**，前一个 stage 的输出是下一个 stage 的输入。尽量让 **$match**、**$project（早期裁剪字段）** 靠前，减少后续 stage 处理的数据量。

**常用 Stage 简述：**

| Stage | 作用 |
|--------|------|
| **$match** | 过滤，等价于 WHERE；尽量早用，可利用索引（与后续 stage 有关，复杂管道需看执行计划） |
| **$project** | 选字段、计算新字段、重命名 |
| **$group** | 分组聚合，类似 GROUP BY；需理解 `_id` 为分组键 |
| **$sort** | 排序；大结果集 + 内存排序有内存限制，需留意 `allowDiskUse` |
| **$limit / $skip** | 截断与分页；深 **$skip** 成本高，大数据量宜用范围分页 |
| **$lookup** | 左外连接另一集合；数据量大时易慢，要想索引与「是否该反范式设计」 |
| **$unwind** | 展开数组，一行变多行 |
| **$facet** | 单管道内多子管道，适合一页多统计 |
| **$bucket / $bucketAuto** | 分桶统计 |

**与 MySQL 的对比（口述用）：**  
聚合管道适合在库内完成「多步变换 + 分组统计」；极复杂的报表或要 join 多张大表时，可能更适合 **同步到 ClickHouse / 离线数仓** 再查。

**性能与实践：**

- 管道深度和中间结果集大小决定延迟；要对 **最大文档数、数组长度** 做产品或技术上的上限。
- **$lookup** 注意外键是否有索引、右表是否过大。
- 需要时可 **`explain`** 聚合管道，观察各 stage 耗时与索引使用。

---

### 4. MongoDB 事务（与「基本操作」并列的工程常识）

- **单文档** 原子性一直成立；**多文档** 事务在副本集 / 分片集群上可用，但有开销与限制（如事务时长、条数建议保守）。
- **写关注（writeConcern）**、**读关注（readConcern）** 影响持久性与读写一致性，与业务容忍度绑定。
- 若强依赖长事务跨很多文档，往往说明模型或边界可再切一刀（例如用 **幂等 + 补偿** 替代大事务）。

---

## 第二部分：GORM

以下以 GORM v2 常见用法为准；具体 API 以项目版本文档为准。

### 1. GORM 在架构中的位置

- 适合 **单服务内** 的 CRUD、简单关联加载、迁移辅助（AutoMigrate 慎用生产大表）。
- **复杂报表、超复杂动态 SQL** 往往手写 SQL 或专用查询层更清晰、可控。

---

### 2. Hook（回调）

**是什么：**  
在创建、更新、删除、查询等生命周期节点插入逻辑，例如统一维护 `CreatedAt` / `UpdatedAt`、审计字段、默认值、软删除前的校验等。

**常见 Hook（示例命名，以官方文档为准）：**

- 创建链：`BeforeCreate` → DB → `AfterCreate`
- 更新链：`BeforeUpdate` → DB → `AfterUpdate`
- 删除链：`BeforeDelete` / `AfterDelete`（含软删除场景）
- 查询链：`AfterFind` 等

**使用原则（面试可展开）：**

1. **Hook 里不要做远程 RPC、长耗时 IO**，否则事务边界拉长、锁持有时间变长、排障困难。  
2. **错误返回**：在 Before* 里返回 `error` 可中止后续写入，适合校验；要区分「业务错误」与「系统错误」日志级别。  
3. **避免递归或隐式再次触发同一 Hook**：例如在 `AfterUpdate` 里又 `Save` 同模型，可能死循环或重复触发，需用标志位或拆层。  
4. **与事务关系**：Hook 仍在当前 `Session` / 事务内执行时，失败应能让整个事务回滚；跨资源（发消息）仍建议 **Outbox** 模式，不要只在 Hook 里发 MQ。

---

### 3. 预加载（Preload）与关联

**是什么：**  
加载主模型时，一并加载关联表数据，避免手写多次查询或漏查。

**典型用法语义：**

- `Preload("Orders")`：按外键分批加载关联，一般生成 **额外 SQL**（不是一条大 JOIN）。
- `Joins`：生成 JOIN，适合需要 **WHERE 条件打在关联表** 或列表页过滤的场景。

**Preload vs Joins（口述对比）：**

| 维度 | Preload | Joins |
|------|---------|--------|
| SQL 形态 | 常为多段查询 | 单条 JOIN |
| 列表 + 条件在子表 | 可能不如 Joins 直观 | 更自然 |
| N+1 | 配置得当可避免裸循环查 | 一条 SQL 但仍要注意是否重复行、是否需要 Distinct |

**嵌套预加载：**  
`Preload("Orders.Items")` 适合树状关联；层级过深时 SQL 次数与数据量上升，要评估是否改 **手写批量查询 + 内存组装**。

**N+1 典型坑：**  
在循环里对每条记录再 `Find`，应改为 **批量 Preload** 或 **Where IN**。开发阶段打开 **Debug / Logger** 数 SQL 条数是基本功。

---

### 4. 事务（Transaction）

**标准写法（语义）：**  
`db.Transaction(func(tx *gorm.DB) error { ... })`  
在闭包内用 **`tx`** 做所有读写，返回 `nil` 提交，返回 `error` 回滚。

**要点：**

1. **只用 tx，不要混用全局 db**，否则不在同一事务里。  
2. **事务尽量短**：不要包 HTTP 调用、RPC、sleep；只包必须原子的 DB 操作。  
3. **嵌套事务**：部分驱动支持 SavePoint，语义要团队统一；否则用「一个大闭包」或业务层编排。  
4. **只读事务**：只读场景可显式只读优化（视驱动与配置而定）。  
5. **与 Hook**：Hook 仍会触发；要清楚「哪一步失败会回滚整体」。

**与 MongoDB 同项目并存时：**  
GORM 只管 SQL 库这一侧；跨 MongoDB + MySQL 的分布式事务一般不用 XA，而用 **Saga / 最终一致 / Outbox**。

---

### 5. GORM 与 MongoDB 在同一 JD 下的分工（总结句）

- **MongoDB**：文档形态、灵活 schema、聚合分析、海量写入或嵌套结构。  
- **GORM + MySQL/PG**：强约束、事务、复杂关系与报表 SQL。  
- **Hook / Preload / 事务** 解决的是「ORM 层一致性、可维护性、性能可控」；**聚合管道** 解决的是「在文档库内完成多步计算」。二者搭配时，边界画在「谁为主存储、谁为查询形态」即可在面试里说清楚。

---

## 延伸阅读（自行查官方文档）

- MongoDB：Aggregation、Indexes、Transactions、Read/Write Concern  
- GORM：Hooks、Associations、Session、Transactions、Performance

若你希望再拆一份 **「仅面试口述版」**（每节 3～5 个问答、无长文），可以说明我单独生成 `mongodb-gorm-interview.md`。
