# Mapping、分词与字段设计

---

## 1. `text` 和 `keyword` 的区别？

**知识点：**
- `text`：分词，用于全文 match；默认不用于排序/聚合（需 `.keyword` 子字段）。
- `keyword`：不分词，整串一个 term，用于精确匹配、sort、terms 聚合。
- **multi-fields**：`title` (text) + `title.keyword` (keyword)。

**口述：**
> text 会走 analyzer 分词，适合「苹果手机」这种模糊搜；keyword 原样存储，适合状态码、订单号、标签。实际 mapping 里标题、描述用 text，同时加 keyword 子字段做排序和按标题精确筛选。对 text 直接 term 查询或聚合是常见踩坑，查不到或极慢。

---

## 2. 为什么 ES 7+ 移除了 type？

**知识点：**
- 6.x 起 deprecated，7.x 默认 `_doc`，8.x 移除多 type。
- 原因：不同 type 同名字段 mapping 冲突；Lucene 层本就单索引单类型更高效。
- 替代：按业务拆 **index**（`orders-2024`）或 **data stream**。

**口述：**
> 以前一个索引多个 type 像 MySQL 多张表，但底层 Lucene 是一个索引，不同 type 字段 mapping 冲突很难维护。现在推荐一个索引一种文档结构，要区分业务就拆索引或用别名 alias 统一查询。

---

## 3. 动态 mapping（dynamic mapping）的风险与治理？

**知识点：**
- `dynamic: true/false/strict`；`strict` 遇未知字段直接拒绝写入。
- 风险：日志 JSON 字段爆炸 → mapping 条目百万级 → 集群状态变大、内存涨。
- 治理：index template 预定义、`dynamic: false` + 白名单、`copy_to` 收敛、写入前清洗。

**口述：**
> 动态 mapping 适合原型，生产最怕上游多传字段，ES 自动加 mapping，几个月后集群元数据巨大。我们会用 index template 固定核心字段，未知字段 strict 拒绝或 dynamic false 忽略，日志场景用 ECS 规范或 ingest pipeline 删无用字段。

---

## 4. analyzer 的组成？索引分词器和搜索分词器为何可不同？

**知识点：**
- analyzer = char_filter（可选）+ **tokenizer** + token_filter（lowercase、stop、synonym…）。
- `index_analyzer`：建索引时；`search_analyzer`：查询时（可更粗或更细）。
- 例：索引 `ik_max_word` 提高召回，搜索 `ik_smart` 提高精度。

**口述：**
> 分词器就是字符流切成 term 再过滤的管道。索引和搜索可以用不同 analyzer，比如索引用最细粒度把词都拆开提高召回，搜索用智能分词减少噪声。同义词一般放 search_analyzer 或单独 synonym_graph filter，改同义词通常要 **reindex** 才能生效到旧数据。

---

## 5. 中文分词 IK：`ik_max_word` vs `ik_smart`？

**知识点：**
- `ik_max_word`：最细拆分，召回高，索引体积大。
- `ik_smart`：粗粒度，适合搜索侧。
- 自定义词典：热词、行业词（需 reload 或重启，视部署方式）。

**口述：**
> 中文没有天然空格，必须靠 IK 等分词器。我们索引侧常用 max_word 保证「北京大学」和「北京」「大学」都能匹配，查询侧用 smart 避免过度拆分。业务新词通过扩展词典维护，上线前要评估对已有索引是否需要 reindex。

---

## 6. `nested` 和 `object` 的区别？

**知识点：**
- `object`：数组里多个对象字段会被 **打平**，丢失对象边界（错误关联）。
- `nested`：每个数组元素是独立 hidden 文档，查询需 `nested` query。
- 代价：nested 写入和查询更贵，能拍平就拍平。

**口述：**
> 如果订单里有多行商品，用 object 存 comments 数组会把不同行的字段交叉匹配，结果错乱。要保证数组内对象原子性，用 nested 类型，查询时包一层 nested query。nested 有额外存储和查询成本，设计时会先想能不能把结构改成父子文档或宽表。

---

## 7. `join` 父子文档适合什么场景？

**知识点：**
- 同索引内 `join` 类型：`parent` / `child` relation。
- 适合：子文档远多于父、更新子文档频繁、查询常带父过滤。
- 缺点：查询复杂、性能不如 routing 同分片设计；8.x 仍可用但不如独立索引主流。

**口述：**
> 父子文档适合一对多且子文档量特别大、又需要和父一起查的场景，比如问答里的 question-answer。但现在更常见的是按 routing 把父子放同一分片，或干脆拆成两个索引用 terms lookup。面试我会说知道 join，但优先 simpler 的数据模型。

---

## 8. 哪些字段类型容易踩坑？

**知识点：**
| 类型 | 注意 |
|------|------|
| `date` | 格式多样，统一 `strict_date_optional_time` |
| `scaled_float` | 省空间，精度有限 |
| `geo_point` / `geo_shape` | 地理查询 |
| `dense_vector` | kNN 语义搜索（8.x+） |
| `flattened` | 整对象 keyword 化，适合动态 key 少量 |

**口述：**
> date 字段格式不统一会导致解析失败或静默错误。金额如果不需要高精度可以用 scaled_float 省空间。现在 RAG 场景会用 dense_vector 做向量检索，但要配 HNSW 参数和内存评估，不是简单加个字段就完事。

---

## 9. 为什么不建议对高基数字段做 terms 聚合？

**知识点：**
- 每个分片构建 global ordinals，高基数 → 内存、GC、慢查询。
- 替代：预聚合表、Rollup/Transform、ClickHouse/Druid 离线、采样。

**口述：**
> 对 user_id、device_id 这种上亿取值的字段做 terms 聚合，相当于让 ES 在内存里维护巨大字典，非常容易 OOM 或打满断路器。正确做法是把维度降下来，或在写入时用 transform 做小时级汇总，查询打汇总索引而不是明细索引。

---

## 10. copy_to、normalizer、ignore_above 实用技巧？

**知识点：**
- `copy_to`：多字段合成一个 `all_text` 统一搜。
- `normalizer`：keyword 的低配 analyzer（如 lowercase），不影响分词。
- `ignore_above`：超长 keyword 不索引，防巨型 term。

**口述：**
> copy_to 适合「标题+摘要+标签」合并成一个搜索字段，简化查询 DSL。keyword 需要大小写不敏感时用 normalizer 而不是 text。URL、堆栈类字段设 ignore_above 防止异常长字符串拖垮索引。
