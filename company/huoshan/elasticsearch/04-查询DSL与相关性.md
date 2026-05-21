# 查询 DSL 与相关性

---

## 1. `match` 和 `term` 的区别？

**知识点：**
- `match`：对查询文本分词 → 各 term 算分（BM25）→ 可指定 `operator: and/or`。
- `term` / `terms`：不分词，精确匹配 keyword 或数值。
- `match_phrase`：短语，要求 term 顺序且邻近（`slop` 允许间隔）。

**口述：**
> match 用于 text 字段的全文搜，会先 analyzer 再打分；term 用于 keyword、状态、ID 这种精确值。典型错误是对 text 用 term 搜「北京」，因为索引里没有整个「北京」一个 term，只有分词后的词项，所以经常零结果。精确搜 text 要用 `title.keyword` 子字段。

---

## 2. bool 查询：`must` / `should` / `filter` / `must_not`？

**知识点：**
- `filter`：布尔过滤，**不算分**，结果可缓存（bitset），CPU 更省。
- `must`：必须匹配且算分；`should`：或关系，可设 `minimum_should_match`。
- 最佳实践：不影响排序的条件一律 `filter`。

**口述：**
> 我会把「要不要出现在结果里」和「怎么排序」分开。状态、时间范围、租户 ID 这种放 filter，不参与打分还能走缓存；用户输入的关键词放 must 里做相关性。should 用于「加分项」比如 VIP 加权，并设置 minimum_should_match 避免 should 变成可选。

---

## 3. `query` 和 `filter` 的本质差异？

**知识点：**
- Query context：计算 `_score`。
- Filter context：yes/no，可缓存，适合高频相同条件。
- `constant_score`：filter 包一层固定分。

**口述：**
> 简单说 query 要算相关性，filter 只判断在不在集合里。高 QPS 列表页如果每次都把 category 写在 must 里，就浪费算分了。改成 filter 后延迟和 CPU 通常明显下降，这是生产优化里最容易落地的一条。

---

## 4. BM25 是什么？和 TF-IDF 比？

**知识点：**
- 默认相似度（5.x+ 默认 BM25，可配 classic TF-IDF）。
- 词频 **饱和**：同一词出现多次收益递减。
- 文档长度归一化：长文档不被动抬高分。
- 参数 `k1`（词频饱和度）、`b`（长度归一强度），一般先用默认。

**口述：**
> BM25 是 ES 默认打分模型，比传统 TF-IDF 更稳，避免「堆关键词」刷分。理解面试即可：词越稀有 IDF 越高，词频有上限，文档越长单 term 贡献会被摊薄。调优相关性我会先动 boost、同义词和业务特征，而不是先改 k1、b。

---

## 5. 如何提升搜索相关性？（召回 vs 排序）

**知识点：**
- **召回**：分词、同义词、拼写纠错（`fuzziness`）、拼音、ngram。
- **排序**：`function_score`、字段 boost、`script_score`、学习排序（业务特征）。
- **评估**：CTR、NDCG、人工标注集、A/B。

**口述：**
> 相关性分两段：先保证该出来的能出来，再保证好的排前面。召回侧检查分词、同义词、是否用了 wrong 字段 type；排序侧用 function_score 叠加销量、点击率、距离等。最后一定要离线评测集加线上 A/B，不能凭感觉调 boost。

---

## 6. 深分页为什么慢？怎么优化？

**知识点：**
- `from + size`：协调节点每分片取 `from+size` 条再全局排序 → 深度越大越慢，默认 `index.max_result_window` 10000。
- **search_after**：上一页 sort 游标，适合实时翻页。
- **scroll / point in time (PIT)**：导出、重处理，不适合用户翻页。
- **slice**：并行 scroll 分片。

**口述：**
> 深分页慢是因为要全局排序并丢掉前面所有结果。用户翻页我们用 search_after 加唯一 tiebreaker 字段（如 `_id`）。批量导出用 PIT+search_after 或 scroll，但 scroll 会占资源，导完要关闭。绝对禁止无限制 from/size 扫全库。

---

## 7. `wildcard` / `regexp` / `prefix` 为什么慎用？

**知识点：**
- 前缀查询可能走 index term，但通配符 leading `*abc` 往往 **全索引扫描**。
- 高成本、易拖垮集群；替代：ngram、edge_ngram、completion suggester。

**口述：**
> 生产尽量避免前导通配符和复杂正则，因为 Lucene 很难用倒排剪枝，会变成扫很多 term。要做「包含」搜索，更稳妥是建 ngram 子字段或接入专门的建议器索引，而不是线上 live 查 `*keyword*`。

---

## 8. 聚合和查询能一起省资源吗？

**知识点：**
- `size: 0` 只要聚合不要 hits。
- `track_total_hits: false` 或设上限，避免精确 count 太贵。
- `terminate_after`：早停（慎用，结果不完整）。

**口述：**
> 看板类请求只关心聚合时设 size 0，并关闭精确 total hits。列表页如果不需要「共 1000 万条」这种精确数，用 track_total_hits 上限到 10000 就够，能明显降协调节点压力。

---

## 9. 慢查询怎么排查？

**知识点：**
- `slowlog`：index.search / index.indexing 阈值。
- `profile: true`：看各阶段耗时（query/fetch/collector）。
- `_explain`：单文档打分解释。
- Hot threads、节点 CPU/heap、分片是否倾斜。

**口述：**
> 我会先看慢日志定位索引和 query DSL，再 profile 看是 query 阶段慢还是 fetch 慢。query 慢常见原因：wildcard、深分页、错误 fielddata、高基数聚合；fetch 慢是 `_source` 太大或 hit 太多。同时看是不是个别分片热点，而不是全集群都慢。

---

## 10. 别名（alias）、索引模板、滚动索引？

**知识点：**
- **alias**：逻辑名指向物理索引，零停机切换（reindex 后 flip alias）。
- **index template / component template**：匹配 `logs-*` 自动应用 mapping/settings。
- **rollover**：按大小/文档数/时间滚动新索引（`logs-000001`）。

**口述：**
> 线上查询只打 alias `orders-search`，底层索引按天或按月滚动，ILM 负责删旧。发版改 mapping 时建新索引 reindex，切 alias 一秒完成，应用无感。这是 ES 运维里最标准的发布姿势。
