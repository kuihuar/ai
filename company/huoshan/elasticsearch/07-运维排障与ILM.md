# 运维排障与 ILM

---

## 1. 集群 yellow / red 排查步骤？

**知识点：**
1. `GET _cluster/health` — 看重分片、未分配数。
2. `GET _cat/shards?v&h=index,shard,prirep,state,unassigned.reason`。
3. `GET _cluster/allocation/explain` — 对 unassigned 分片给原因。
4. 查节点：`_cat/nodes`、磁盘、日志、近期重启/滚动发布。

**口述：**
> yellow 先看是不是副本没配上，单节点、磁盘水位、节点数不够都常见。red 说明 primary 丢了，立即看 explain API 和该索引最近是否有关闭节点、强制分配错误。我会先恢复可用性（扩容、释放磁盘、修正 allocation 规则），再追查根因，避免手动 `allocate_empty_primary` 除非确认数据可丢。

---

## 2. 磁盘水位（watermark）三阶段？

**知识点：**
- **low**：不再往该节点分配新分片。
- **high**：尝试把分片迁走。
- **flood_stage**：索引 **read-only-allow-delete**，阻断写入，需释放空间后 `PUT index/_settings` 解除。

**口述：**
> 磁盘用到 85%、90%、95% 这类阈值会逐级收紧。flood 阶段索引变只读，写入全挂，这是保护集群不整体崩溃。处理就是删旧索引、扩盘、迁冷数据，然后取消 read-only。平时用 ILM 删过期索引，不要等 flood 才救火。

---

## 3. ILM（Index Lifecycle Management）怎么设计？

**知识点：**
- 阶段：**hot** → **warm** → **cold** → **delete**（可插 `frozen`）。
- 动作：rollover、shrink、forcemerge、allocate、delete、searchable snapshot。
- 绑定：`index.lifecycle.name` 或通过 index template 默认 policy。

**口述：**
> 日志类索引我会配 ILM：hot 阶段 rollover（50GB 或 1 天），warm 缩副本+迁 warm 节点，30 天后 cold，90 天 delete。这样磁盘可预测，运维不用 cron 手工删索引。policy 变更只影响新索引，老索引要 reapply 或 rollover 后的新 backing index。

---

## 4. data stream 和传统按天索引的区别？

**知识点：**
- **data stream**：写入抽象层，背后多个 backing index，`@timestamp` 必填，适合日志/指标。
- 查询打 stream 名，自动查所有 backing；rollover 由 ILM 触发。
- 对比 `logs-2024.05.20` 手工建索引：stream 统一命名、权限、生命周期。

**口述：**
> 8.x 日志场景推荐 data stream，应用只写 `logs-app`，不用每天换索引名。ILM 绑在 stream 的 template 上自动 rollover 和删除。和传统按天索引比，运维更简单，但要注意 mapping 升级要走 reindex 新 template。

---

## 5. reindex 和零停机切换流程？

**知识点：**
- `POST _reindex`：源 → 目标（可跨集群 `source.remote`）。
- 流程：建新索引/新 mapping → reindex → 校验 count/抽样 → **alias 切换** → 删旧索引。
- 增量：双写 + 定时 reindex 差量，或 CDC 持续同步。

**口述：**
> 改 mapping 不能原地改字段类型，必须新索引 reindex。我们会先 alias 双写或短暂停写，全量 reindex 后对比文档数，然后把 alias 从 old 切到 new，应用无感。大索引用 sliced reindex 并行，并限流避免打满集群。

---

## 6. 快照备份与恢复？

**知识点：**
- 仓库：S3、共享文件系统（需注册 `repository`）。
- `PUT _snapshot/repo/snap-1` 定期 SLM（Snapshot Lifecycle Management）。
- 恢复：`POST _snapshot/repo/snap-1/_restore`，可改索引名、部分索引。

**口述：**
> 生产必须有 snapshot，配合 SLM 每天增量。恢复演练季度做一次，否则真宕机才发现 S3 权限或路径错了。恢复时注意集群版本兼容和索引名冲突，生产恢复常恢复到新集群再切流量。

---

## 7. 滚动升级（rolling restart）注意点？

**知识点：**
- 先禁分配 `cluster.routing.allocation.enable=primaries` 或逐节点 drain。
- 版本跳跃需看官方 matrix；6→7→8 大版本不能一步跨。
- 升级后 cluster state 变大时检查 `_cluster/settings` 持久化。

**口述：**
> 滚动重启时一次只停一个 data 节点，等分片 reallocate 变 green 再下一个。升级大版本前先看 breaking changes，比如 type 移除、Java 版本、安全默认开启。我们有次升级失败是因为插件版本和 ES 版本不匹配，所以插件也要纳入 checklist。

---

## 8. 常见异常怎么口述？

| 现象 | 可能原因 | 处理方向 |
|------|----------|----------|
| `rejected execution` | 线程池队列满 | 降并发、bulk 限流、扩容 |
| `circuit_breaking_exception` | 聚合/fielddata 过大 | 优化 DSL、预聚合 |
| `too_many_buckets` | 聚合桶超限 | 提高 limit 或改查询 |
| `version_conflict` | 并发更新同文档 | 重试、乐观锁 |
| `index read-only` | flood watermark | 清磁盘、解除只读 |

**口述：**
> 线程池拒绝说明写入或查询超过节点处理能力，先限流再扩容，不是单纯调大队列。断路器是内存保护，根本是查询太大。版本冲突说明没有做幂等或用了外部版本号管理。

---

## 9. 监控指标该看哪些？

**知识点：**
- 集群：health、unassigned shards、pending tasks。
- 节点：CPU、heap、GC、disk、indexing/search rate、rejected。
- 索引：segment count、merge 耗时、refresh/flush 时间。
- 业务：P95 查询延迟、bulk 失败率、同步滞后（若 ES 从 DB 同步）。

**口述：**
> 除了 green 还要看 pending cluster tasks 堆积，说明 master 压力大。数据节点重点看 search/indexing 线程池 rejected 和 old GC。业务侧如果 ES 是 MySQL 同步来的，要监控 lag 和补偿队列深度，这比只看 ES 健康更有用。
