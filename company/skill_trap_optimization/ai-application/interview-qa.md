# AI/LLM 面试问答

## 一、基础理论

### Q1：Transformer 的核心结构是什么？

Transformer 由 Encoder 和 Decoder 组成，核心组件：

```
Encoder: Input → Embedding + Positional Encoding → Multi-Head Self-Attention → FFN → ... ×N 层
Decoder: Output → Embedding + Positional Encoding → Masked Self-Attention → Cross-Attention → FFN → ... ×N 层
```

**Attention 计算公式**：

$$\text{Attention}(Q, K, V) = \text{softmax}\left(\frac{QK^T}{\sqrt{d_k}}\right)V$$

**关键组件**：
- **Self-Attention**：每个 token 与序列中所有 token 计算注意力权重，捕获全局依赖
- **Multi-Head**：多组 Q/K/V，不同"头"关注不同的语义子空间
- **Positional Encoding**：注入位置信息（正弦/余弦函数或可学习参数）
- **FFN**：两层全连接，中间激活 ReLU/GELU
- **LayerNorm + 残差连接**：稳定训练，防止梯度消失

**为什么除以 √dk？**
点积 QK^T 的方差会随 d_k 增大而增大，导致 softmax 输出的梯度进入极小区。除以 √d_k 将方差控制在 1，保持梯度稳定。

---

### Q2：GPT 和 BERT 的核心区别？

| 维度 | GPT | BERT |
|------|-----|------|
| 架构 | Decoder-only | Encoder-only |
| Attention Mask | Causal（因果，只看左侧） | Bidirectional（双向，看全部） |
| 预训练目标 | Next Token Prediction | MLM（Masked Language Model）+ NSP |
| 擅长 | 生成（文本生成、对话） | 理解（分类、NER、QA） |
| 生成方式 | 自回归（逐 token 生成） | 不支持直接生成 |
| 代表 | GPT-3/4, LLaMA | BERT, RoBERTa |

**追问：为什么现在 Decoder-only 架构（GPT 系列）更主流？**
- 统一生成和理解任务（万物皆可序列生成）
- 上下文学习（In-Context Learning）只在 Decoder-only 模型中涌现
- 训练效率：Causal Attention 可通过 KV Cache 高效推理
- Scaling Law 更好：Decoder-only 模型增大后能力增长更可预期

---

### Q3：大模型训练流程是什么样的？

```
1. 预训练 (Pre-training)
   万亿级 token，Next Token Prediction
   → 基座模型（Base Model）

2. 监督微调 (SFT, Supervised Fine-Tuning)
   数万条高质量指令数据
   → 能理解指令

3. 对齐 (Alignment, RLHF/DPO)
   - RLHF: 奖励模型 → PPO 强化学习
   - DPO: 直接在偏好对上优化（更简单，替代 RLHF 的趋势）
   → 符合人类偏好

4. (可选) 领域微调
   特定领域数据继续训练
   → 医疗/法律/代码等垂直模型
```

**追问：RLHF vs DPO 的优劣？**

| 维度 | RLHF | DPO |
|------|------|-----|
| 复杂度 | 高（4 个模型：policy, reference, reward, critic） | 低（2 个模型） |
| 训练稳定性 | 差，需要精细调参 | 好，直接优化 |
| 效果 | SOTA（但差距在缩小） | 接近 RLHF |
| 工程代价 | 大 | 小 |
| 趋势 | - | 越来越用 DPO 替代 RLHF |

---

### Q4：什么是注意力机制中的 KV Cache？为什么能加速推理？

自回归生成时，每生成一个新 token，需要重新计算所有历史 token 的 K 和 V。但其实历史 token 的 K、V 是不变的，没必要重复算。

**KV Cache**：把每个 token 计算过的 K、V 缓存起来，下一轮只算新 token 的 Q/K/V，计算量从 O(n²) 降到 O(n)。

```
无 KV Cache：每步重新算全部 n 个 token → O(n²d)
有 KV Cache：每步只算新 token，矩阵拼接缓存 → O(nd)
```

**代价**：显存占用随序列长度线性增长。长序列（>128K）时 KV Cache 成为瓶颈，需要量化（KV8/int4）或稀疏化。

---

### Q5：LoRA 微调的原理是什么？

**LoRA**（Low-Rank Adaptation）：在预训练权重旁加一个低秩分解的可训练矩阵，冻结原权重不动。

```
原计算：  h = Wx
LoRA:    h = Wx + BAx
其中 W 冻结，只训练 A (d×r) 和 B (r×k)，r << d
```

**为什么有效**：
- 微调时权重变化落在低秩空间内（Intrinsic Dimensionality假说）
- 假设 $\Delta W$ 的秩很低，可以用两个小矩阵 $B \cdot A$ 近似

**优点**：
- 训练参数极少（原模型的 0.1%~1%）
- 不增加推理延迟（$BA$ 可合并回 $W$）
- 多个 LoRA 可快速切换（一个基座 + 多个 LoRA 适配多任务）

---

## 二、RAG（检索增强生成）

### Q6：RAG 的完整流程是什么？有哪些核心设计选择？

```
离线阶段：文档 → 切分 → Embedding → 向量数据库
在线阶段：用户问题 → Embedding → 检索 Top-K → (Rerank) → Prompt 拼接 → LLM 生成
```

**核心设计选择**：

| 环节 | 关键决策 | 影响 |
|------|----------|------|
| 切分 (Chunking) | 大小、重叠、策略（固定/语义） | 太小缺上下文，太大降低检索精度 |
| Embedding | 模型选型、维度 | 影响检索效果天花板 |
| 检索策略 | 稠密/稀疏/混合检索 | 关键词 vs 语义匹配各有优劣 |
| Top-K | 返回数量 | K 大覆盖全但 cost 高，K 小精确但可能遗漏 |
| Rerank | 是否需要、模型选型 | 用计算换精度，通常显著提升效果 |
| Prompt 模板 | 引用方式、拒绝策略 | 影响幻觉率 |

### Q7：检索阶段纯向量检索有什么问题？混合检索（Hybrid Search）是什么？

**纯向量检索的问题**：
- 对精确关键词匹配不敏感（如产品型号、人名）
- 长文本的 Embedding 可能丢失细节

**混合检索 = 稠密检索 + 稀疏检索**

```
稠密检索 (Dense)：Embedding 向量相似度 → 语义匹配
稀疏检索 (Sparse)：BM25/TF-IDF 关键词匹配 → 精确匹配
                              ↓
                   加权融合（Reciprocal Rank Fusion）
                              ↓
                         Rerank → Top-K
```

```python
# RRF (Reciprocal Rank Fusion) 融合公式
score(doc) = Σ 1/(k + rank_i(doc))
# k=60 是常用经验值
```

### Q8：如何评估 RAG 系统的好坏？

**自动评估指标**：

| 指标 | 测量什么 | 计算方式 |
|------|----------|----------|
| Recall@K | Top K 中命中几个正确答案 | 命中数 / 正确答案总数 |
| MRR (Mean Reciprocal Rank) | 第一个正确结果排第几 | avg(1/rank) |
| NDCG | 排序质量（考虑排名位置） | 位置权重 + 相关性加权 |
| Faithfulness | 回答是否忠实于检索文档 | LLM-as-Judge 逐个陈述检查 |
| Answer Relevance | 回答是否切题 | LLM-as-Judge |

**人工评估维度**：准确性、完整性、可读性、引用准确性。

---

## 三、Agent

### Q9：Agent 的核心循环是什么样的？

```
Loop:
    1. LLM 接收任务 + 工具列表 + 历史上下文
    2. LLM 决策：直接回答 OR 调用工具
    3. 如果调用工具：执行 → 观察结果 → 回到步骤 1
    4. 如果直接回答：输出，结束
```

**追问：Agent 容易出什么问题？**
- **死循环**：LLM 反复调用同一工具拿不到有效信息。措施：设置 max_iter
- **幻觉工具调用**：调用不存在的工具或参数错误
- **遗忘目标**：上下文变长后忘记原始任务
- **工具误用**：参数错误导致工具执行失败

### Q10：Multi-Agent 架构如何设计？

常见模式：

| 模式 | 描述 | 例子 |
|------|------|------|
| **主从模式** | 一个主 Agent 分配任务给子 Agent | 主 Agent 分析需求 → 分发 |
| **辩论模式** | 多个 Agent 生成方案并辩论 | 代码审查 Agent + 生成 Agent |
| **流水线** | 各 Agent 串行执行业务 Pipeline | 检索 Agent → 分析 Agent → 写作 Agent |
| **分层模式** | Agent 按领域能力层级化 | 通用 Agent → 专业 Agent（代码/数学等） |

```go
// 主从模式
type Orchestrator struct {
    specialistAgents map[string]*Agent
}

func (o *Orchestrator) Execute(ctx context.Context, task string) (string, error) {
    // 1. 由 Master Agent 拆解任务
    subtasks := o.masterAgent.Decompose(task)

    // 2. 路由到对应的 Specialist Agent
    results := make([]string, len(subtasks))
    for i, st := range subtasks {
        agent := o.route(st.Domain)
        results[i], _ = agent.Execute(ctx, st.Description)
    }

    // 3. 汇总
    return o.masterAgent.Synthesize(results)
}
```

---

## 四、模型部署与推理优化

### Q11：大模型推理加速有哪些技术？

| 技术 | 原理 | 加速比 | 精度损失 |
|------|------|--------|----------|
| **量化 (Quantization)** | FP16→INT8/INT4 | 2-4x | 低（INT8 几乎无损） |
| **KV Cache** | 缓存历史 K/V | 每步 O(n) 变 O(1) | 无损 |
| **Flash Attention** | 分块计算 + IO 优化 | 2-4x | 无损 |
| **Speculative Decoding** | 小模型草稿 + 大模型验证 | 1.5-3x | 无损 |
| **Continuous Batching** | 动态管理多请求的 KV Cache | 吞吐量 10x+ | 无损 |
| **Tensor/Pipeline Parallelism** | 多 GPU 并行 | 近线性 | 无损 |
| **vLLM / TensorRT-LLM** | 综合推理框架 | 显著 | 可忽略 |

**追问：INT8 量化为什么几乎无损？**
大模型权重分布通常集中在狭窄范围，INT8 的表示精度（-128~127）足够——尤其在使用动态量化（per-tensor/per-channel/per-token）时。

### Q12：如何为应用选择部署方案？

| 需求 | 推荐方案 |
|------|----------|
| 快速验证、低 QPS | Ollama + 消费级 GPU |
| 生产环境 API 服务 | vLLM / TensorRT-LLM + A100/H100 |
| 低延迟需求 | TensorRT-LLM（优化最深） |
| 低成本、非实时 | llama.cpp / Ollama + CPU |
| 多租户 SaaS | vLLM（原生 multi-LoRA 支持） |

---

## 五、系统设计

### Q13：设计一个支持百万用户的 AI 客服系统

**整体架构**：

```
CDN / WAF
    ↓
API Gateway（限流、鉴权、路由）
    ↓
会话管理服务（WebSocket 长连接 / SSE / HTTP polling）
    ↓
意图识别 → 简单 FAQ → 直接返回
         → 复杂问题 → RAG Pipeline → LLM → 返回
         → 转人工 → 排队 → 人工客服
    ↓
缓存层（语义缓存 + Redis 热点缓存）
    ↓
监控层（Token 消耗、延迟、满意度追踪）
```

**关键设计要点**：

1. **会话管理**：用 Redis 存储 session，保存最近 N 轮对话用于多轮理解
2. **分级响应**：FAQ 直接匹配（<10ms）→ 语义缓存（<50ms）→ RAG+LLM（2-5s）
3. **削峰填谷**：高峰期 influx 用消息队列缓冲请求
4. **成本控制**：语义缓存（目标命中率 30%）、模型分级路由（简单用 fast model）
5. **可观测性**：每次对话记录 trace——retrieval recall、生成延迟、用户反馈

### Q14：设计一个拍照搜题系统

> 详见 [photo-search.md](photo-search.md) — 完整技术链路拆解

---

## 六、Prompt Engineering 与 LLM 能力

### Q15：什么是 In-Context Learning（ICL）？为什么有效？

ICL 是通过 few-shot 示例影响模型行为，不需要更新参数。

**为什么有效**：
- 大模型的 Attention 机制可根据上下文中的示例动态建立"隐式任务表示"
- Transformer 等价于元学习的梯度下降步骤（理论假说）

**提升 ICL 效果的技巧**：
- 示例多样化 + 分布与测试集一致
- 示例顺序有影响（可尝试不同顺序取平均）
- 加入 CoT 推理链的示例显著提升复杂任务效果

### Q16：Chain-of-Thought (CoT) 和 Zero-shot-CoT 的区别？

- **CoT**：用带推理步骤的 Few-shot 示例引导模型推理
- **Zero-shot-CoT**：只加 `"Let's think step by step"`，不需要示例

```markdown
# Few-shot CoT（需要示例）
Q: 小明有 5 个苹果，给了小红 2 个，又买了 3 个，现在有几个？
A: 小明开始有 5 个苹果。给了小红 2 个后剩 3 个。又买了 3 个后一共有 6 个。答案是 6。

Q: [新问题]
A:

# Zero-shot CoT（不需要示例）
Q: [问题]
A: Let's think step by step.
```

**追问：什么时候 CoT 不适用？**
- 简单事实性问题（"法国的首都是哪里"）
- 需要想象力的创意写作
- Token 成本/延迟敏感场景（CoT 大幅增加 token 消耗）

### Q17：幻觉（Hallucination）的成因和缓解策略

**成因**：
- 训练数据噪声和偏差
- 模型是概率模型，本质上是"猜"下一个 token
- 最大似然目标不区分"不知道"和"编造"
- 训练数据和用户 prompt 的不匹配

**缓解策略**：

| 策略 | 原理 | 效果 |
|------|------|------|
| RAG | 用外部知识约束生成 | 最有效，大幅降低幻觉 |
| 明确告知"不知道" | SFT 阶段训练模型说不知道 | 中等 |
| 约束解码 | 限制输出空间（如只能从候选列表选） | 特定场景有效 |
| Self-reflection | 让模型检查自己输出的一致性 | 中等 |
| 降低 temperature | 减少输出随机性 | 副作用：降低创造性 |
