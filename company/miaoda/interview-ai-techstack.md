# 大模型应用架构、Agent、RAG、MCP 面试问答专题

---

## 一、大模型应用架构基础

**Q1: 一个典型的 LLM 应用平台由哪些核心模块组成？请画出架构并解释各模块职责。**

核心模块：

```
┌─────────────────────────────────────────────────────┐
│                     接入层                           │
│  API Gateway（限流/鉴权/路由/流式代理）               │
└────────────┬──────────────────────────┬─────────────┘
             │                          │
┌────────────▼──────────┐  ┌────────────▼────────────┐
│      控制面            │  │        推理面            │
│                        │  │                          │
│ ┌──────────────────┐  │  │ ┌────────────────────┐  │
│ │ Agent 编排引擎    │  │  │ │ 模型推理网关       │  │
│ │ (ReAct/Plan-Exec)│  │  │ │ (路由/负载均衡/降级)│  │
│ └────────┬─────────┘  │  │ └────────┬───────────┘  │
│          │             │  │          │              │
│ ┌────────▼─────────┐  │  │ ┌────────▼───────────┐  │
│ │ RAG 检索管线     │  │  │ │ 推理引擎集群       │  │
│ │ (Embedding/检索  │  │  │ │ vLLM/TGI/         │  │
│ │  /重排序/生成)   │  │  │ │ TensorRT-LLM      │  │
│ └────────┬─────────┘  │  │ └────────────────────┘  │
│          │             │  │                          │
│ ┌────────▼─────────┐  │  │                          │
│ │ MCP Server 管理  │  │  │                          │
│ │ (工具注册/发现   │  │  │                          │
│ │  /生命周期)      │  │  │                          │
│ └──────────────────┘  │  │                          │
│                        │  │                          │
│ ┌──────────────────┐  │  │                          │
│ │ Prompt 管理      │  │  │                          │
│ │ (模板/版本/AB)   │  │  │                          │
│ └──────────────────┘  │  │                          │
└────────────────────────┘  └──────────────────────────┘
             │                          │
┌────────────▼──────────────────────────▼─────────────┐
│                   基础设施层                          │
│  PostgreSQL │ Redis │ Milvus/Qdrant │ Kafka │ MinIO  │
│  GPU Cluster (A100/A10/H100) │ K8s │ Prometheus     │
└─────────────────────────────────────────────────────┘
```

各模块职责：

- **接入层**：统一 API 入口，兼容 OpenAI API 格式；按租户/API Key 做鉴权和配额控制；流式推理的 SSE 透传。
- **Agent 编排引擎**：解析用户意图，按 ReAct 等模式自主调用工具链，管理多轮推理循环和终止条件。
- **RAG 检索管线**：文档入库→分段→向量化→索引；Query→改写→检索→重排序→上下文拼接。
- **MCP Server 管理**：工具的注册、发现、生命周期管理，支持 stdio/HTTP 多种传输协议。
- **Prompt 管理**：模板库存、版本管理、变量注入、AB 测试、效果追踪。
- **模型推理网关**：按模型/场景/成本多维路由，负载均衡，故障降级，Token 计数与计费。
- **推理引擎**：vLLM/TGI/TensorRT-LLM 等高效推理引擎，管理 GPU 资源和 KV Cache。
- **基础设施层**：关系库存业务数据、Redis 做缓存和会话、向量库存 Embedding、消息队列做异步解耦、对象存储存模型权重和文档。

---

**Q2: 如何理解大模型应用中的"上下文窗口"？长上下文的挑战和解决方案是什么？**

上下文窗口（Context Window）是 LLM 单次推理能接受的 token 上限（如 GPT-4 128K、Claude 200K）。所有 System Prompt + 历史对话 + RAG 检索内容 + 用户 Query 都计入上下文。

挑战：
- **注意力复杂度 O(n²)**：上下文越长，注意力计算量和显存消耗呈平方增长。
- **"Lost in the Middle"**：长上下文中，模型对开头和结尾的内容关注度高，中间部分信息容易被忽略。
- **KV Cache 膨胀**：长序列的 KV Cache 占据大量显存，限制并发能力。
- **推理延迟**：长上下文导致 Prefill 阶段耗时显著增加。

解决方案：
- **上下文压缩**：对长对话历史做摘要（Summarization），只保留关键信息而非原始全文。
- **滑动窗口注意力**（Sliding Window Attention）：每个 token 只关注前后 w 个 token，将复杂度降到 O(n×w)，如 Mistral 的方案。
- **RAG 替代长上下文**：不把所有文档塞进 Context，而是检索最相关片段后再拼入。
- **分块处理**：将长文档拆成多段，逐段推理后再综合结果（如 Map-Reduce、Refine 策略）。
- **Prompt 压缩技术**：LLMLingua 等框架对 Prompt 做无损/有损压缩。

---

**Q3: 什么是 Prompt Engineering？有哪些核心技术和实践原则？**

Prompt Engineering 是通过设计输入指令来引导 LLM 产生期望输出的系统化方法，不是简单的一句话提示，而是工程化实践。

核心技术：

**Few-Shot Prompting（少样本提示）**
- 在 Prompt 中给出 2-5 个输入→输出示例，让模型理解任务模式。
- 示例选择很关键：覆盖边界情况、格式一致性、与当前 query 的相似度。

**Chain-of-Thought（思维链，CoT）**
- 要求模型"一步步思考"（Let's think step by step），展示推理过程。
- 显著提升数学、逻辑和多步推理任务准确率。
- 进阶：Auto-CoT（自动生成推理链）、Complexity-based CoT（选复杂样本做示例）。

**角色扮演（Persona）**
- 设定角色身份、专业背景、语气风格、输出格式等约束。

**结构化输出**
- 明确要求 JSON/Markdown/表格 等输出格式，给出 Schema 示例。
- JSON Mode / Function Calling / Structured Outputs（OpenAI）确保格式可靠。

**实践原则**
- 清晰 > 简短：说清楚比说得少重要。
- 正面指令 > 负面指令：告诉模型"做什么"比"不要做什么"效果好。
- 分隔符：用 `###`、`"""`、`<context>` 等标记区分不同段落。
- 迭代优化：Prompt 版本化、AB 测试、效果评估（准确性/一致性/延迟/成本）。

---

**Q4: 如何评估一个大模型应用的效果？有哪些核心指标？**

评估维度分层：

**模型能力层**
- 通用基准：MMLU（多任务语言理解）、HumanEval（代码生成）、GSM8K（数学推理）。
- 特定领域评测：自建领域测试集，覆盖业务场景。

**检索质量层（RAG 专用）**
- Ragas 框架指标：
  - **Context Precision**：检索结果中相关文档的排序精度。
  - **Context Recall**：答案所需信息在检索结果中的覆盖程度。
  - **Faithfulness**：答案中的每个断言是否可从检索上下文推导（幻觉检测）。
  - **Answer Relevancy**：答案是否紧贴问题。

**应用效果层**
- 端到端准确率：人工标注评估，答案满足用户需求的比例。
- 首 Token 延迟（TTFT）+ Token 生成速度（TPS）。
- 任务成功率：Agent 任务的自主完成率（不需人工介入的比例）。

**运营指标层**
- 用户满意度（点赞/点踩比例）、复问率（同一个问题问多次的比例）。
- 日均调用量、Token 消耗趋势、成本核算。

**线上监控层**
- 幻觉率监控：采样线上对话，统计事实性错误比例。
- 安全合规：敏感内容拦截率、越狱攻击检测率。
- 模型漂移监控：同一 Prompt 在不同时间的输出一致性（分布偏移检测）。

---

## 二、RAG（检索增强生成）深度

**Q5: RAG 的完整技术链路是什么？每一步有哪些关键设计选型？**

完整链路分为离线和在线两条 Pipeline。

**离线索引 Pipeline**

```
文档源（PDF/Word/HTML/DB）
  → 文档解析（Unstructured/pypdf/LangChain Document Loader）
    → 文本分段 Chunking
      → Embedding 向量化
        → 写入向量数据库 + 构建索引
```

每一步选型：
- **文档解析**：非结构化文档用 Unstructured 库或 Azure Document Intelligence；PDF 复杂排版考虑保留阅读顺序和表格结构。
- **分段策略**：
  - 固定大小（如 512 tokens）+ 重叠（如 50 tokens）——简单可靠。
  - 语义分段（按段落/章节自然边界）——保留语义完整性。
  - 递归字符分割（`RecursiveCharacterTextSplitter`）——面向代码等特殊格式。
  - **父子文档（Parent-Child）**：检索用小粒度 chunk，返回大粒度 parent，兼顾精确召回和上下文完整性。
- **Embedding 模型**：中文场景推荐 bge-large-zh-v1.5 / GTE-Qwen2-7B / text2vec-large-chinese；多语言用 multilingual-e5-large。
- **向量数据库**：
  - Milvus：分布式、高性能、混合检索，适合大规模。
  - Qdrant：Rust 实现、性能好、过滤能力强，适合中等规模。
  - PGVector：基于 PostgreSQL，轻量化、事务支持好，适合小规模或已有 PG 的场景。

**在线 Query Pipeline**

```
用户 Query
  → Query Rewriting（改写/扩展/分解/HyDE）
    → Embedding 向量化
      → 混合检索（Dense + Sparse）
        → Reranker 重排序
          → 上下文拼接（Context Window 分配）
            → LLM 生成
```

关键选型：
- **Query 改写**：多轮对话指代消解 + 意图补充；复杂查询拆成子问题分步检索。
- **混合检索**：向量检索（语义相似）+ BM25 关键词检索（精确匹配），互补召回。权重通常 0.7:0.3 或动态调整。
- **Reranker**：BGE-Reranker-v2-m3 / Cohere Rerank API / ColBERT 等 Cross-Encoder 模型对 Top-K（如 100）精排，取 Top-N（如 5）送入 LLM。这是提升回答质量性价比最高的手段之一。
- **上下文拼接**：相关片段按相关度排序 + 关键信息前置 + 标注来源。

---

**Q6: RAG 中常见的"检索失败"模式有哪些？分别如何优化？**

| 失败模式 | 现象 | 根因 | 优化方案 |
|---------|------|------|---------|
| **检索不相关** | 返回的文档与问题无关 | Embedding 模型与领域不匹配；Query 表达不清 | 微调 Embedding 模型（使用领域数据）；Query Rewriting 改写增强 |
| **检索不完整** | 答案需要的信息分散在多处，只召回部分 | Chunk 过小导致信息碎片化；仅依赖向量检索遗漏关键词匹配 | 增大 Chunk 或使用父子文档；引入 BM25 混合检索；多路召回融合 |
| **检索到过时信息** | 文档内容已更新，但向量库中的旧版未被替换 | 索引未及时同步数据源变更 | 增量索引更新；文档版本管理；时效性过滤 |
| **Lost in the Middle** | 正确答案在第 5 个文档中，但 LLM 只关注了前 2 个 | 长上下文中段信息被忽略 | Reranker 精排使最相关结果前置；对关键信息做高亮标记 |
| **检索到矛盾信息** | 多个文档包含相互矛盾的说法，LLM 无法分辨 | 数据源未做权威性管控 | 引入来源权威性权重；冲突检测后用 Prompt 指示 LLM 以最新/最权威来源为准 |

---

**Q7: 混合检索（Hybrid Search）的融合策略有哪些？如何实现？**

混合检索 = 稠密检索（Dense, 语义向量）+ 稀疏检索（Sparse, BM25/TF-IDF 关键词）。

融合策略：

1. **RRF（Reciprocal Rank Fusion，倒数排名融合）**
   ```
   RRF_score(doc) = Σ 1/(k + rank_i(doc))
   ```
   - k 默认 60，用于平滑。
   - 无需归一化分数，对排名敏感而非绝对分值，简单有效。推荐首选。

2. **线性加权融合**
   ```
   final_score = α × normalized_dense_score + (1-α) × normalized_sparse_score
   ```
   - 需要先对两种分数做归一化（Min-Max 或 Z-score），否则量纲差异导致一侧权重被吞掉。
   - α 通常取 0.6-0.8（偏好语义）。

3. **两阶段融合**
   - 第一阶段：用 Dense 检索 Top-K1（如 200）。
   - 第二阶段：在 K1 内用 Sparse 做重排或布尔过滤，筛掉不满足关键词约束的文档。
   - 适合强关键词约束场景（如必须包含特定实体名）。

实现（Python 伪代码）：
```python
def hybrid_search(query, top_k=10, alpha=0.7):
    # Dense search
    query_vec = embedding_model.encode(query)
    dense_results = vector_db.search(query_vec, top_k=top_k * 2)
    
    # Sparse search
    sparse_results = bm25_index.search(query, top_k=top_k * 2)
    
    # RRF fusion
    rrf_scores = {}
    for rank, doc in enumerate(dense_results):
        rrf_scores[doc.id] = rrf_scores.get(doc.id, 0) + 1 / (60 + rank)
    for rank, doc in enumerate(sparse_results):
        rrf_scores[doc.id] = rrf_scores.get(doc.id, 0) + 1 / (60 + rank)
    
    # Sort and return top_k
    sorted_docs = sorted(rrf_scores.items(), key=lambda x: x[1], reverse=True)
    return [doc_id for doc_id, _ in sorted_docs[:top_k]]
```

---

**Q8: 如何对 RAG 系统进行效果评估和持续优化？**

评估体系搭建：

**离线评估（Offline）**
- 构建黄金测试集：收集 200-500 个真实用户 Query + 标准答案 + 标注相关文档。
- 用 Ragas 框架自动化评估（见 Q4），定期跑一轮生成评估报告。
- 关注趋势：Context Recall < 0.6 说明遗漏了信息，需要调 Embedding 或 Chunk 策略；Faithfulness < 0.7 说明幻觉严重。

**在线评估（Online）**
- 收集用户反馈：点赞/点踩、答案复制率、追问率。
- 采样人工标注：每两周抽 100 条会话做人工评测。
- AB 测试：新检索策略先在 10% 流量上验证，对比核心指标。

**持续优化 Pipeline**
```
评估发现问题 → 假设根因 → 实验验证 → 上线全量 → 再次评估
```
- 举个迭代例子：Context Recall 低 → 假设 Chunk 太小吃掉了上下文 → 改为父子文档检索 → AB 测试观察 Recall 提升 → 确认后全量上线。

---

**Q9: Graph RAG 是什么？和传统 RAG 有什么区别和优势？**

**传统 RAG 的局限**：依赖向量相似度搜索，擅长语义匹配但缺乏对**实体关系、多跳推理、全局总结**的支持。例如："总结这家公司过去三年的技术战略变化"——需要理解跨越多个文档的时间线和概念演进，单纯靠 chunk 相似度检索很难覆盖。

**Graph RAG（微软提出）**：
在传统 RAG 基础上增加**知识图谱层**。核心流程：

1. **图构建**：从文档中提取实体（人物/公司/技术/事件）和关系（收购/发布/合作），构建知识图谱（存于 Neo4j/NetworkX）。
2. **社区检测**：用 Leiden 算法将图中紧密关联的节点划分为"社区"（Community）。
3. **社区摘要生成**：对每个社区用 LLM 生成结构化摘要（描述了社区内的实体和关系）。
4. **检索时**：
   - 向量检索获取相关 Chunk（传统 RAG）。
   - 同时从图谱中匹配相关实体和社区摘要（图谱检索）。
   - 融合两种信息源，LLM 获得更全局和结构化的上下文。

**优势**：
- 多跳问答能力显著提升（"A 公司的 CEO 之前在哪家公司任职"）。
- 全局性、总结性问题回答质量高。
- 实体消歧和关系追溯更精确。

**代价**：
- 图构建成本高（每百万 token 需额外 LLM 调用来提取实体关系）。
- 索引构建慢，数据更新后需重建图谱。
- 适合文档规模大、需要全局理解的知识库场景（如企业知识库、法律文档库）。

---

**Q10: RAG 系统中如何处理多模态（图片/表格）内容的检索？**

三种策略：

**策略一：文本描述替代（简单场景）**
- 图片用多模态 LLM（GPT-4V/Llama 3.2 Vision/Qwen-VL）生成文字描述。
- 表格抽取为 CSV/JSON 文本。
- 将描述文本与原始文本一起 Embedding 入库。
- 检索时命中文本描述，回答时引回原图/原表。适用于图片不多、精度要求不高的场景。

**策略二：多模态 Embedding（中等场景）**
- 使用 CLIP/BLIP/SigLIP 等多模态 Embedding 模型，将图片直接编码为向量。
- 文本和图片在同一个向量空间中可互检索。
- 表格转图片后用多模态 Embedding 或专门的 Table Embedding 模型。

**策略三：独立管线 + 融合（复杂场景）**
- 文本检索和图片检索各走独立管线。
- 两路结果融合（RRF 或单独拼入），LLM 最终生成时同时参考文本和图片/表格。
- 需要支持多模态输入的 LLM（如 GPT-4o / Claude 3.5 / Gemini）才能直接"看懂"图片。

实践建议：大多数企业场景用策略一即可覆盖 80% 需求。如果表格或图片是信息核心载体，优先文本化（表格→Markdown 表格格式）；图片数量多且信息密度高，再升级到策略二/三。

---

## 三、Agent 架构深度

**Q11: Agent 的核心架构是什么？标准 Agent 循环是怎样的？**

Agent 的标准抽象模型：

```
┌─────────────────────────────────────────┐
│                Agent 核心               │
│                                         │
│  ┌───────────┐  ┌───────────────────┐  │
│  │  Planner   │  │    Memory         │  │
│  │ (任务规划) │  │ (短期/长期/工作)  │  │
│  └─────┬─────┘  └───────────────────┘  │
│        │                                │
│  ┌─────▼─────────────────────────────┐ │
│  │         Reasoning Engine          │ │
│  │   (LLM 推理，ReAct/CoT/ToT)       │ │
│  └─────┬─────────────────────────────┘ │
│        │                                │
│  ┌─────▼──────┐  ┌──────────────────┐  │
│  │ Tool Use   │  │   Reflection     │  │
│  │ (工具调用) │  │   (自我反思)     │  │
│  └────────────┘  └──────────────────┘  │
└─────────────────────────────────────────┘
```

标准 Agent 循环（ReAct）：

```
1. 用户输入任务 Task
2. Agent 进入循环：
   a. Think（思考）: 分析当前状态，规划下一步
   b. Act（行动）: 根据规划选择并调用工具
   c. Observe（观察）: 获取工具返回结果
   d. Reflect（反思）: 判断是否需要调整计划
   e. 判断是否满足终止条件（任务完成 / 达到最大步数 / 死胡同）
3. 输出最终结果
```

伪代码：
```python
class Agent:
    def run(self, task: str, max_steps: int = 10) -> str:
        messages = [{"role": "system", "content": SYSTEM_PROMPT},
                     {"role": "user", "content": task}]
        
        for step in range(max_steps):
            response = llm.chat(messages, tools=available_tools)
            
            if response.has_final_answer():
                return response.content
            
            if response.has_tool_calls():
                for tool_call in response.tool_calls:
                    result = execute_tool(tool_call.name, tool_call.arguments)
                    messages.append({"role": "tool", "content": result, 
                                     "tool_call_id": tool_call.id})
            
            # 更新短期记忆
            self.short_term_memory.append({"step": step, "thought": response})
        
        raise MaxStepsExceeded(f"任务在 {max_steps} 步内未完成")
```

---

**Q12: ReAct、Plan-and-Execute、ReWOO 等 Agent 模式的对比和选型？**

| 模式 | 核心思想 | 执行方式 | 优点 | 缺点 | 适用场景 |
|------|---------|---------|------|------|---------|
| **ReAct** | 每步交替推理和行动：Think → Act → Observe → Think ... | 单步串行 | 灵活，能根据观察动态调整；实现简单 | Token 消耗大；复杂任务步骤多，容易跑偏 | 通用场景，中等复杂度任务 |
| **Plan-and-Execute** | 先制定完整计划，再逐步执行（计划可动态修正） | 计划阶段 + 执行阶段分离 | 全局观好；复杂任务步骤清晰；可跟踪进度 | 初始计划可能不完美；Token 消耗大（生成计划 + 执行） | 多步复杂任务（如代码项目、研究分析） |
| **ReWOO** | Reasoning WithOut Observation：一次性规划所有工具调用，批量执行 | 全规划 → 全执行 → 汇总 | Token 效率极高（不需每次观察后重新推理）；速度快 | 依赖初始规划质量；无法根据中间结果动态调整 | 确定性高、工具调用结果可预期的任务 |
| **Reflexion** | 执行失败后自我反思，修正策略后重试 | ReAct + 反思循环 | 能从失败中学习；复杂任务成功率高 | 执行时间长；Token 消耗更大 | 需要多次尝试的困难任务（代码调试、策略优化） |
| **Tree/Graph of Thoughts** | 多条推理路径并行探索，选择最佳路径 | 树/图搜索 | 探索空间大；找到更优解 | Token 消耗极大；延迟高 | 需要多方案对比的决策场景 |

选型建议：
- 任务步骤少、环境反馈重要 → ReAct。
- 任务可预先规划清楚 → Plan-and-Execute。
- 成本和延迟敏感 → ReWOO（略过中间思考）。
- 任务难度大、首次成功率低 → Reflexion。

---

**Q13: Agent 的工具调用（Function Calling/Tool Use）机制是什么？如何确保调用可靠性？**

机制原理（以 OpenAI 兼容 API 为例）：

1. **工具注册**：在请求中传入 `tools` 参数，每个 tool 包含：
   - `type: "function"`
   - `function.name`：工具唯一标识
   - `function.description`：工具功能描述（LLM 据此判断是否调用）
   - `function.parameters`：JSON Schema 定义输入参数

2. **模型决策**：LLM 收到用户消息 + 工具列表，自主判断：
   - 不需要工具 → 直接回复文本（`finish_reason: "stop"`）
   - 需要工具 → 返回 tool_calls 列表（`finish_reason: "tool_calls"`），包含 `function.name` 和 `function.arguments`(JSON)

3. **工具执行**：开发者解析 tool_calls → 执行对应函数 → 以 `role: "tool"` 追加结果到消息列表。

4. **模型消化结果**：再次调用 LLM，模型解析工具返回结果后生成最终回答或继续调用工具。

保障可靠性：

- **严格的 JSON Schema**：`parameters` 声明类型(type)、枚举(enum)、必填(required)、描述(description)，减少模型生成非法参数。
- **Tool Description 写好"何时用"**：不光写功能，还要说明"当用户问 X 或需要 Y 时调用此工具"。
- **错误返回结构化**：工具执行失败时不抛异常，返回 `{"error": "原因"}`，让 LLM 自行决定重试还是降级回答。
- **参数校验 + 兜底**：收到 tool_calls 后在代码层校验参数合规性，不合法时直接返回结构化错误消息。
- **max_steps 限制**：限制最大工具调用轮次，防止无限循环。

---

**Q14: Agent 的记忆系统（Memory）如何设计？**

记忆分层设计：

```
┌────────────────────────────────────────────┐
│              短期记忆（Short-term）          │
│  当前会话的对话历史，window_size=20 轮      │
│  直接拼入 LLM Context                      │
│  存储：Redis, TTL=会话超时时间              │
└────────────────┬───────────────────────────┘
                 │
┌────────────────▼───────────────────────────┐
│              长期记忆（Long-term）           │
│  跨会话的持久化记忆：                        │
│  - 事实记忆（用户姓名/偏好/重要决策）        │
│  - 经验记忆（某任务的处理经验）              │
│  存储：向量数据库 + 结构化 DB               │
│  召回：基于当前 Query 语义检索 + 重要性过滤  │
└────────────────┬───────────────────────────┘
                 │
┌────────────────▼───────────────────────────┐
│              工作记忆（Working）             │
│  当前任务的临时草稿：                        │
│  - 中间计算结果                              │
│  - 当前计划的执行进度                        │
│  - 已收集但尚未处理的信息                    │
│  存储：Redis, TTL=任务执行期间               │
└────────────────────────────────────────────┘
```

**长期记忆的工程实现**：

```python
class LongTermMemory:
    def __init__(self, vector_db, llm):
        self.vector_db = vector_db
        self.llm = llm
    
    def extract_memories(self, conversation: list[Message]) -> list[Memory]:
        """从对话中提取值得记住的事实/偏好/经验"""
        prompt = """从此对话中提取值得记住的关键信息，JSON 格式返回：
        [{
            "type": "fact|preference|decision|experience",
            "content": "记忆内容",
            "importance": 0.0-1.0  // 重要性评分
        }]
        规则：仅提取对后续对话可能有用的信息；琐碎的寒暄跳过。"""
        
        response = self.llm.chat(prompt + str(conversation), response_format="json")
        return parse_memories(response)
    
    def store(self, memory: Memory):
        """存储记忆：写入向量库（用于语义检索）+ PG（用于精确查询）"""
        embedding = embed(memory.content)
        self.vector_db.insert(
            id=memory.id,
            vector=embedding,
            metadata={"type": memory.type, "importance": memory.importance}
        )
    
    def recall(self, query: str, top_k: int = 5) -> list[Memory]:
        """检索相关记忆"""
        query_vec = embed(query)
        results = self.vector_db.search(query_vec, top_k=top_k * 2)
        # 按重要性过滤和排序
        results = [r for r in results if r.metadata["importance"] > 0.3]
        return results[:top_k]
```

**记忆管理策略**：
- **记忆衰减**：旧记忆按时间衰减（指数衰减或线性衰减），检索时乘衰减因子。
- **去重与合并**：新旧记忆冲突时（"我住在北京"→"我搬到上海了"），新记忆覆盖旧记忆或标记旧记忆失效。
- **遗忘机制**：低重要性 + 长期未被召回的记忆自动清理。

---

**Q15: 多 Agent 协作架构有哪些模式？如何实现 Agent 间的消息传递和任务编排？**

协作模式：

**1. 顺序流水线（Sequential Pipeline）**
```
Agent A 输出 → Agent B 输入 → Agent C 输入 → 最终输出
```
- 适合作业流程固定的场景（如 文档解析 → 摘要 → 翻译）。
- 实现简单，但缺乏灵活性和错误恢复。

**2. 主从模式（Orchestrator-Worker）**
```
         ┌───────┐
         │Orch-   │
         │estrator│  ← 拆解任务，分配，汇总
         └──┬──┬──┘
      ┌─────┘  └─────┐
   ┌──▼──┐  ┌────┐  ┌──▼──┐
   │Worker│  │... │  │Worker│
   │  A   │  │    │  │  N   │
   └──────┘  └────┘  └──────┘
```
- 一个主 Agent（Orchestrator）拆解任务、分配给专门的 Worker Agent、汇总结果。
- 最通用的模式，适合大多数复杂任务。

**3. 去中心化协作（Peer-to-Peer / Debate）**
```
Agent A ←→ Agent B ←→ Agent C
```
- 平等角色，通过讨论/辩论达成共识。适合需要多视角审视的决策或创作任务。
- 实现复杂，消息路由规则需精心设计。

**4. 分层架构（Hierarchical）**
```
         Top Agent
        /    |    \
   Middle  Middle  Middle
    /  \    /  \    /  \
  Leaf Leaf Leaf Leaf Leaf
```
- 组织化结构，上层做决策和协调，下层做执行。适合极复杂任务（如大型软件项目开发模拟）。

**消息传递实现**：
- 通过共享消息总线（如 Redis Pub/Sub、Kafka）或 Agent 框架内置机制（LangGraph State、CrewAI Task）。
- 消息格式标准化：`{"from": "agent_a", "to": "agent_b", "type": "task|result|inquiry", "payload": {...}}`

**任务编排实现**：
- 用 DAG（有向无环图）描述 Agent 间的依赖关系：Task B 依赖 Task A 的结果。
- 框架支持：LangGraph 构建 Agent 状态图、CrewAI 定义 Task 依赖链、AutoGen 的 GroupChat 模式。

---

**Q16: Agent 开发中常见的问题有哪些？如何解决？**

| 问题 | 现象 | 根因 | 解决方案 |
|------|------|------|---------|
| **无限循环** | Agent 反复调用同一工具，无法终止 | 工具描述不清导致模型"迷路"；缺少明确的终止条件 | 设 `max_steps`；在 Prompt 中明确终止条件；检测到循环模式时强制中断并降级回答 |
| **幻觉工具调用** | 调用了不存在的工具，或参数明显错误 | 工具描述不清晰或与用户意图不匹配 | 严格的 JSON Schema 约束；代码层参数校验 + 结构化错误返回；Tool Description 细化"何时用" |
| **Token 消耗过大** | 单个任务消耗数万 Token，成本高 | 每步 Think 过于啰嗦；历史信息未压缩 | 要求 LLM 精简思考；工具返回信息只给关键字段；历史对话超过一定轮次做摘要压缩 |
| **上下文窗口溢出** | 多轮工具调用后 Context 超限 | 工具返回内容过多 + 历史积累 | 滑动窗口裁剪历史；重要信息做摘要保留；工具返回限制长度 |
| **任务漂移** | Agent 逐渐偏离原始目标，做起了其他事情 | 多步执行中逐步被新信息分散注意 | 在 System Prompt 中强调原始目标；每 N 步自我检查是否偏离；Orchestrator 角色监督 |
| **工具执行超时** | 某个工具调用长时间无响应，Agent 挂起 | 外部 API 慢或不可达 | 工具执行加超时（如 30s）；超时后返回结构化错误告知 Agent；重要工具增加重试机制 |
| **安全风险** | Agent 执行了危险操作（删除数据/执行任意代码） | 工具权限过大，没有安全边界 | 工具分级（只读/读写/管理员）；危险操作需人工确认（Human-in-the-Loop）；沙盒隔离执行 |

---

## 四、MCP 协议深度

**Q17: MCP 协议的设计理念和核心架构是什么？解决了什么问题？**

**背景**：在 MCP 出现之前，每个 LLM 应用都要自己实现工具连接——查数据库要写一个 Connector，调 API 要写另一个 Connector。形成 M×N 集成问题（M 个 LLM 应用 × N 个工具/数据源），每个组合都要定制开发。

**MCP 的设计理念**：
- 类比 USB-C 协议之于硬件：定义了一套 LLM 与外部工具/数据源交互的标准接口。只要实现了 MCP，任何 LLM 应用都能使用任何 MCP Server。
- 将 M×N 问题降维为 M+N：LLM 应用只要支持 MCP Client 就能调用所有 MCP Server；工具开发者只要实现 MCP Server 就能被所有 LLM 应用调用。

**核心架构**：

```
┌──────────┐     MCP Protocol     ┌──────────┐
│  MCP     │◄───────────────────►│  MCP     │
│  Client  │   (JSON-RPC 2.0)    │  Server  │
│          │                      │          │
│ Claude   │                      │ DB Server│
│ VS Code  │                      │ File Sys │
│ 自研APP  │                      │ API Wrap │
└──────────┘                      └──────────┘

Transport: stdio | Streamable HTTP
```

**三大原语（Primitives）**：

| 原语 | 用途 | 类比 | 控制方 |
|------|------|------|--------|
| **Tools** | 模型可调用的函数（有副作用），如查数据库、发邮件、创建文件 | REST POST/PUT/DELETE | 模型主动调用 |
| **Resources** | 模型可读取的数据（无副作用），如文件内容、API 数据、数据库查询结果 | REST GET | 模型/用户均可读取 |
| **Prompts** | 预定义的提示词模板，标准化常用交互模式 | 快捷指令/模板 | 用户选择使用 |

**解决了什么问题**：
- 标准化：不受限于特定 LLM 厂商（OpenAI Function Calling、Anthropic Tool Use 等格式被 MCP 统一抽象）。
- 生态复用：一次开发 MCP Server，所有 MCP Client 都能使用。
- 安全性：通过 Transport 层做进程隔离（stdio）或网络鉴权（HTTP），工具权限可精确控制。
- 动态发现：Client 可自动发现 Server 暴露的所有 Tools/Resources/Prompts，无需手动配置。

---

**Q18: MCP 的 Tools、Resources、Prompts 三个原语分别如何设计？各自的边界在哪？**

**Tools（工具）**

设计场景：模型需要**执行操作**，有副作用。

```python
# 示例：天气查询工具
{
    "name": "get_weather",
    "description": "查询指定城市的当前天气。当用户询问天气、气温、是否下雨等问题时调用此工具。",
    "inputSchema": {
        "type": "object",
        "properties": {
            "city": {
                "type": "string",
                "description": "城市名称，如 'Beijing' 或 '北京'"
            },
            "unit": {
                "type": "string",
                "enum": ["celsius", "fahrenheit"],
                "description": "温度单位，默认 celsius"
            }
        },
        "required": ["city"]
    }
}
```

Tool 设计要点：
- `description` 是给 LLM 看的，写清楚"什么时候用"和"返回什么"，比参数 Schema 更重要。
- Tool 应该做**一件事**，粒度适中（不要一个 Tool 查天气+查股票+发邮件）。

**Resources（资源）**

设计场景：暴露**只读数据**，让 LLM 了解上下文。

```python
# 示例：数据库表结构资源
{
    "uri": "db://schema/users",
    "name": "用户表结构",
    "description": "users 表的 DDL 和字段说明",
    "mimeType": "text/plain"
}
# 还可支持 Resource Template：db://table/{table_name}/schema
```

Resource 设计要点：
- Resource 是**被动提供数据**，Client 读取后放入 Context。
- 适合文件内容、数据库 Schema、系统状态、配置信息等只读上下文。

**Prompts（提示模板）**

设计场景：标准化**常用交互模式**，提供可选的参数化 Prompt。

```python
{
    "name": "code_review",
    "description": "对代码进行 Review",
    "arguments": [
        {"name": "language", "description": "编程语言", "required": True},
        {"name": "code", "description": "待审查代码", "required": True}
    ]
}
```

**三者的边界判断**：

| 问题 | 选 Tool | 选 Resource | 选 Prompt |
|------|---------|------------|-----------|
| 会改变外部系统状态吗？（写入/删除/发送） | 是 | 否 | 否 |
| 只是提供数据让 LLM 阅读理解吗？ | 否 | 是 | 否 |
| 是用户可选的标准交互模板吗？ | 否 | 否 | 是 |
| 需要严格的参数校验和执行逻辑吗？ | 是 | 否 | 否 |

典型误区：把数据库 SELECT 查询做成 Tool 而不是 Resource。简单查询适合 Resource，但带复杂参数/RAG 检索的查询适合 Tool（因为需要参数校验和动态处理）。

---

**Q19: MCP 的 Transport 层有哪些选择？各有什么适用场景？**

| Transport | 通信方式 | 生命周期 | 优点 | 缺点 | 适用场景 |
|-----------|---------|---------|------|------|---------|
| **stdio** | 标准输入/输出流，JSON-RPC over stdin/stdout | 随 Client 进程启停 | 零网络配置；天然进程隔离安全；低延迟 | Client 与 Server 必须在同一台机器；只能一对一 | 本地工具（文件操作、本地数据库、CLI 工具、代码执行沙箱） |
| **Streamable HTTP** | HTTP POST + SSE（2025 新规范，替代旧 SSE transport） | Server 独立运行，多 Client 共享 | 远程访问；多 Client 共享；可独立扩缩容 | 需网络配置、鉴权、HTTPS；有网络延迟 | 远程服务（企业内部 API、云服务、共享工具平台） |

**stdio 模式的工作流程**：
1. Client 以子进程方式启动 Server（`subprocess.Popen`）。
2. Client 通过 stdin 发送 JSON-RPC 请求。
3. Server 通过 stdout 返回 JSON-RPC 响应。
4. Client 退出时关闭 stdin，Server 收到 EOF 后退出。

**Streamable HTTP 模式的工作流程**：
1. Server 作为独立 HTTP 服务运行，暴露 `/mcp` 端点。2. Client 发送 POST 请求（JSON-RPC），Server 立即返回 200 + SSE Stream（对于需要流式响应的场景，如工具执行进度）。
3. 支持 Session 管理，多 Client 可同时连接。
4. 需配套鉴权层（OAuth2/API Key）。

**选型指南**：
- 秒哒平台的内置工具（文件处理、沙箱执行）→ stdio，不暴露网络攻击面。
- 企业共享的工具服务（CRM 查询、数据仓库访问）→ Streamable HTTP，支持多租户和多语言 SDK。

---

**Q20: 如何开发一个 MCP Server？以 Python SDK 为例。**

完整示例：数据库查询 MCP Server。

```python
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent
import asyncpg

# 1. 初始化 Server
server = Server("postgres-mcp-server")

# 2. 注册 Tools
@server.list_tools()
async def list_tools() -> list[Tool]:
    return [
        Tool(
            name="query_database",
            description="执行 SQL 查询并返回结果。用于查询数据库中的业务数据。"
                        "当用户需要查找、统计、分析数据时调用此工具。"
                        "注意：仅支持 SELECT 查询，不支持 INSERT/UPDATE/DELETE。",
            inputSchema={
                "type": "object",
                "properties": {
                    "sql": {
                        "type": "string",
                        "description": "要执行的 SELECT SQL 语句"
                    }
                },
                "required": ["sql"]
            }
        ),
        Tool(
            name="list_tables",
            description="列出数据库中所有的表和视图，帮助了解有哪些数据可用。",
            inputSchema={
                "type": "object",
                "properties": {}
            }
        )
    ]

# 3. 实现 Tool Handler
@server.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    conn = await asyncpg.connect(DSN)
    try:
        if name == "query_database":
            sql = arguments["sql"]
            # 安全检查：仅允许 SELECT
            if not sql.strip().upper().startswith("SELECT"):
                return [TextContent(
                    type="text",
                    text="错误：仅允许 SELECT 查询。"
                )]
            rows = await conn.fetch(sql)
            # 格式化结果为 Markdown 表格
            result = format_table(rows)
            return [TextContent(type="text", text=result)]
        
        elif name == "list_tables":
            rows = await conn.fetch(
                "SELECT table_name FROM information_schema.tables "
                "WHERE table_schema='public'"
            )
            result = "\n".join(r["table_name"] for r in rows)
            return [TextContent(type="text", text=result)]
    finally:
        await conn.close()

# 4. 注册 Resources
@server.list_resources()
async def list_resources():
    return [
        Resource(
            uri="db://schema",
            name="数据库 Schema",
            description="所有表的 DDL 结构信息",
            mimeType="text/plain"
        )
    ]

@server.read_resource()
async def read_resource(uri: str):
    if uri == "db://schema":
        conn = await asyncpg.connect(DSN)
        # 读取所有表的 DDL
        schemas = await get_all_table_schemas(conn)
        await conn.close()
        return schemas

# 5. 启动 Server
async def main():
    async with stdio_server() as (read, write):
        await server.run(read, write)

if __name__ == "__main__":
    import asyncio
    asyncio.run(main())
```

**关键设计原则**：
- **安全性**：SQL 查询仅允许 SELECT，防止删除/修改数据库。
- **错误处理**：返回结构化错误消息而非直接崩溃，让 LLM 自行决定如何处理。
- **格式化输出**：返回 Markdown 表格等 LLM 易消化的格式。
- **资源发现**：同时注册 `list_tables` Tool 和 Schema Resource，让 LLM 有"探索"能力。

---

**Q21: MCP 与 OpenAI Function Calling 的关系和区别？已经用了 Function Calling 还需要 MCP 吗？**

| 维度 | OpenAI Function Calling | MCP |
|------|------------------------|-----|
| **定位** | LLM API 的一个功能：模型决定何时调用函数并生成参数 JSON | 完整的 Client-Server 协议：定义工具的发现、调用、生命周期管理 |
| **绑定** | 与 OpenAI API 格式绑定（其他厂商有类似实现但格式各异） | 厂商无关的开放协议 |
| **工具注册** | 每轮 API 调用时传入 tools 列表 | Client 启动时一次性发现并注册所有 Tools |
| **工具实现** | 开发者自行实现函数的调用、结果格式化、错误处理 | Server 封装了工具的实现细节，Client 只需调用标准接口 |
| **生态复用** | 每个应用自己写工具集成代码 | 一次开发 MCP Server，所有 MCP Client 复用 |
| **生命周期** | 随 API 请求的一次性声明 | Server 独立运行，有自己的生命周期 |

**有 Function Calling 还需要 MCP 吗？**

需要。它们是不同层次的抽象：
- Function Calling 解决的是"LLM 怎么决定调用哪个工具、传什么参数"的问题。
- MCP 解决的是"工具怎么定义、怎么部署、怎么被不同 LLM 应用复用"的问题。

在工程实践中，两者是互补关系：

```
用户请求 → LLM (Function Calling 决策调用哪个工具)
                ↓ Tool Call (name + arguments)
          MCP Client → MCP Server (执行工具逻辑)
                ↓ Tool Result
          LLM (消化结果，生成回答)
```

MCP 替代的是工具接入层（过去自己写一堆 Adapter），Function Calling 替代的是 LLM 的意图识别层。即使用 MCP，LLM 内部仍然用 Function Calling（或 Anthropic Tool Use）机制来生成 tool_calls。

---

**Q22: 在生产环境中部署 MCP Server 需要注意哪些问题？**

**1. 鉴权与安全**
- stdio 模式天然隔离，权限 = 启动进程的用户权限。
- Streamable HTTP 模式必须有鉴权层：OAuth2 Bearer Token / mTLS / IP 白名单。
- 工具按能力分级：只读工具 / 读写工具 / 管理员工具。根据调用方身份授予最低权限。

**2. 性能与并发**
- `list_tools` 和 `call_tool` 是高频接口，需保证低延迟（< 100ms 建立连接，< 5s 工具执行超时）。
- 单个 Server 的资源限制：连接池大小、并发工具调用数。
- 对于耗时长的工具（如大文件处理），返回进度通知（通过 Notifications）。

**3. 可观测性**
- 记录每次工具调用的元数据：调用方、工具名、参数（脱敏）、耗时、结果状态。
- 异常告警：工具调用失败率突增、响应时间突增。
- 与 OpenTelemetry 集成，串起 LLM API Trace + MCP Tool Trace 的端到端链路。

**4. 错误处理与重试**
- 区分可重试错误（网络超时、临时不可用）和不可重试错误（参数非法、权限不足）。
- 工具超时后的熔断：短时间内失败超过阈值，暂停暴露该工具，避免级联故障。
- 错误消息对 LLM 友好：返回自然语言描述的错误，而不是堆栈 trace。

**5. 版本管理**
- MCP Server 的 Tools/Resources 定义变更（改名、改参数、删工具）需要版本化管理。
- Client 启动时通过 `list_tools` 动态发现，但建议 Client 端有版本兼容性检查。

**6. 多租户**
- 不同租户的 MCP Server 实例隔离（Namespace/独立安装），防止数据泄露。
- 工具调用结果的缓存策略按租户隔离。

---

## 五、推理引擎与模型服务化

**Q23: vLLM 的核心技术创新是什么？为什么它比传统推理方式快？**

vLLM 两大核心创新：

**1. PagedAttention（分页注意力）**

传统 KV Cache 为每个请求分配整块连续显存，用完释放后形成空洞 → 碎片化严重 → 显存利用仅 ~30%。

PagedAttention 将 KV Cache 切分成固定大小的页（Block，如 16 token/page），页内连续、页间链表连接：
- 新请求按需申请页，释放后页回池，无碎片。
- 多个请求可共享同一页（如相同的 System Prompt 共享 KV Cache）。
- 显存利用率提升到 70-80%+。

对比：

```
传统：  [AAAAAAAAA........] [BBBBB..............] [CCCCCCCCCCCC.......]
        (大量碎片/空洞，无法复用)

PagedAttention: [A][B][A][C][B][C][A][ ][C][ ][ ][ ][ ][ ][ ][ ]
                (页粒度分配，页回收后立即可被新请求使用)
```

**2. Continuous Batching（连续批处理）**

传统静态 Batch：等 N 个请求一起处理，最长的请求做完才释放整个 batch → 短请求被长请求拖累。

Continuous Batching：以每个推理 step（生成一个 token）为调度粒度：
- Step 1：推理 [Req1, Req2, Req3] 各生成一个 token。
- Step 2：Req1 已完成（EOS），移出 batch；Req4 刚到达，加入 batch；推理 [Req2, Req3, Req4]。
- 随到随处理、生成完即退出，不等其他请求。

效果：GPU 计算单元持续保持高负载，吞吐提升 **5-10x**。

---

**Q24: 大模型推理中的 Prefill 和 Decode 阶段分别是什么？为什么 Prefill 更慢？**

两阶段的本质：

**Prefill（预填充/编码）阶段**
- 输入：整个 Prompt（input tokens），一次性计算。
- 做的事：并行处理所有 input tokens，生成第一个 output token + 填充整个序列的 KV Cache。
- 计算特点：**Compute-bound**（计算瓶颈），并行度 = input_token 数量 × 模型层数。
- 慢的原因：需要承载大量的矩阵乘法（尤其长 Prompt，如 10K input tokens）。

**Decode（解码/生成）阶段**
- 输入：每步仅一个新 token。
- 做的事：自回归生成，每步处理一个 token，用已有的 KV Cache 生成下一个 token。
- 计算特点：**Memory-bound**（内存/显存带宽瓶颈），每次只计算 1 个 token 的 Attention（并行度低），受限于显存带宽读写 KV Cache。
- 相对 Prefill 快：单步计算量小（1 token vs 整个输入序列）。

首 Token 延迟（TTFT）= Prefill 时间。
每 Token 生成延迟（TPOT）= Decode 单步时间。

优化方向：
- Prefill 慢 → 用 Prefix Caching（相同前缀复用 KV Cache）、Chunked Prefill（将长 Prefill 拆成多块与 Decode 混跑）。
- Decode 显存带宽瓶颈 → 量化（INT4/INT8 减少 KV Cache 体积）、Speculative Decoding（用小型草稿模型批量生成候选 token，大模型一次验证多个）。

---

**Q25: 如何选择合适的模型量化方法？GPTQ、AWQ、GGUF 各自的适用场景？**

| 方法 | 原理 | 量化对象 | 精度保持 | 推理需要 | 最佳场景 |
|------|------|---------|---------|---------|---------|
| **GPTQ** | 逐层权重量化，用校准数据补偿误差（OBS 最优脑手术） | 权重 | INT4 精度损失 ~1-3% | GPU 推理引擎（vLLM/TGI/TensorRT-LLM） | GPU 服务端推理，显存优先 |
| **AWQ** | 基于激活值重要性做量化，保留对输出影响大的通道精度 | 权重 | INT4 精度与 GPTQ 相当或更好 | vLLM/TGI 原生支持 | GPU 服务端推理（推荐优先选 AWQ） |
| **GGUF** | llama.cpp 原生格式，支持混合精度（不同层不同 bit）、CPU 推理优化 | 权重 | 灵活，可在准确率与体积间平衡 | llama.cpp / Ollama | CPU 推理、边缘设备、本地部署 |
| **BitsAndBytes** | HuggingFace 深度集成，4bit NormalFloat + 双重量化 | 权重 | INT4 训练/微调场景优化 | Transformers 加载 + QLoRA | 微调训练（QLoRA）、实验探索 |
| **FP8** | 原生 FP8 Tensor Core（H100/L40S），硬件加速 | 权重+激活值 | 精度损失极小（<0.5%） | vLLM FP8、TGI FP8 | H100 等新 GPU，追求极致吞吐 |

**选型决策树**：
```
有 H100/FP8 硬件？ → FP8（速度最快，精度最好）
   ↓ 否
GPU 服务端推理？ → AWQ（vLLM/TGI 原生支持，推荐首选）
   ↓ 否
CPU 推理/边缘？ → GGUF（llama.cpp/Ollama）
   ↓ 否
微调训练？ → BitsAndBytes 4-bit + QLoRA
```

---

**Q26: 如何设计一个模型推理网关？核心功能和技术难点是什么？**

核心功能模块：

**1. 统一 API 适配**
- 对外暴露 OpenAI 兼容 API（`/v1/chat/completions`），降低业务方接入成本。
- 内部 Adapter 模式：每个推理引擎（vLLM/TGI/Triton/云端 API）封装为独立 Adapter。
- 模型与引擎解耦：同一模型可配置多种推理后端，按需切换。

**2. 智能路由**
- 路由维度：按模型名称、按场景（对话/代码/嵌入）、按成本优先/延迟优先。
- 基于后端实时负载的加权路由（从 Prometheus metrics 采集队列深度、GPU 利用率、TTFT）。
- 故障转移：某模型实例不健康，自动切到备选实例或备选模型。

**3. 流式代理**
- SSE/WebSocket 透传，逐 Token 返回到客户端。
- 流中断处理：客户端断开 → 通知后端取消推理（释放 KV Cache）。
- 过程中做增量 Token 计数，用于实时配额校验。

**4. 限流与配额**
- 多维度限流：全局网关 QPS、API Key 级别 Token 配额、用户级别并发限制。
- 排队机制：高负载时请求入优先级队列，超时降级或拒绝。
- 软硬限：软限降级到备用模型，硬限直接返回 429。

**5. 计费与用量统计**
- 流式中每一段 SSE 报文落地时增量统计 Token（输入+输出），上报 Kafka。
- 按租户/应用/模型/时间窗口聚合消费。

**6. 安全**
- 敏感词检测（Prompt 注入过滤 + Response 内容合规）。
- API Key 泄露检测（异常调用模式识别）。

**技术难点**：
- **流式代理的 Token 精确计数**：没有接入推理引擎内部计数器的话，只能按 `choices[0].delta.content` 逐段估算，精度取决于 Tokenizer。建议推理引擎暴露 metrics，由网关采集。
- **后端负载感知**：单纯轮询/加权轮询不感知推理实例的实际负载（队列深度 × 平均 token 生成时间），需要基于推理引擎 Prometheus metrics 做动态路由。
- **取消传播**：客户端断开 SSE 连接后，网关需立即向推理引擎发取消请求（`POST /v1/chat/completions cancel` 或 HTTP/2 RST_STREAM），否则浪费 GPU 算力。

---

## 六、实战综合场景

**Q27: 设计一个企业知识库 RAG 问答系统，从零开始你的技术方案是怎样的？**

**第一阶段：需求明确**
- 知识范围：哪些类型的文档（制度文档/技术文档/工单记录）？数量多大？
- 问答类型：事实查询 vs 总结分析 vs 多文档对比？
- 用户规模：并发量？延迟要求？
- 权限要求：不同人员能查的文档范围是否不同？

**第二阶段：数据接入与预处理**
- 文档接入：支持哪些来源（本地文件/SharePoint/Confluence/数据库）？增量同步还是全量？
- 文档解析：不同类型用不同解析器（PDF→Unstructured，Word→python-docx，网页→HTML 提取）。
- 分块策略：技术文档 512 tokens + 语义分段，制度文档按条款自然边界分段。
- 元数据保留：文档标题、来源、更新时间、权限标签（用于检索后过滤）。

**第三阶段：索引构建**
- Embedding 模型选型：中文 → bge-large-zh-v1.5（768维），综合性价比。
- 向量数据库：规模 < 10 万文档且已有 PG → PGVector；大规模 → Milvus。
- BM25 索引：用 Elasticsearch 或 Milvus 内置的 Sparse Vector。
- 构建 Pipeline 的并发和失败重试。

**第四阶段：检索与生成**
- Query Rewriting：多轮改写 + 同义词扩展。
- 混合检索：Dense (bge) + Sparse (BM25)，RRF 融合。
- Reranker：BGE-Reranker-v2-m3，Top-100 → Top-5。
- Prompt 模板：将检索上下文 + 用户问题 + 来源引用格式注入。
- LLM 选型：平衡成本与效果（GPT-4o-mini/DeepSeek-V3/Qwen3）。

**第五阶段：评估与上线**
- 构建测试集（至少 200 条真实场景 Query + 标注答案）。
- Ragas 自动化评估 + 人工抽检。
- 灰度上线（10% 流量）→ 观察指标 → 全量。
- 建立反馈闭环：用户点踩 → 人工复查 → 更新索引/调整策略。

**第六阶段：持续运营**
- 文档更新同步机制（监听变更 → 重新索引受影响文档）。
- 效果监控：Context Recall、Faithfulness、Answer Relevancy 趋势 Dashboard。
- 定期 Bad Case 分析 → 策略迭代。

---

**Q28: 客户要求私有化部署的 Agent 平台中，Agent 可以调用客户自建的 HTTP API，如何设计这个能力？**

设计目标：让客户在平台上注册自己的 API 作为 Agent 可调用的 Tool。

**方案设计**：

**1. Tool 注册接口**
```json
POST /api/v1/tools/register
{
    "name": "query_customer_crm",
    "description": "查询客户 CRM 系统中的客户信息。当需要查找客户联系方式、订单记录时调用。",
    "method": "POST",
    "url": "http://customer-internal.corp.com/api/crm/query",
    "headers": {
        "Authorization": "Bearer {{secret:crm_api_key}}",
        "Content-Type": "application/json"
    },
    "request_template": {
        "customer_name": "{{customer_name}}",
        "query_type": "{{query_type}}"
    },
    "parameters": {
        "type": "object",
        "properties": {
            "customer_name": {"type": "string", "description": "客户名称"},
            "query_type": {"type": "string", "enum": ["contact", "order", "all"], 
                           "description": "查询类型：contact=联系方式，order=订单，all=全部"}
        },
        "required": ["customer_name"]
    },
    "timeout_seconds": 30,
    "retry": {"max_retries": 2, "backoff": "exponential"}
}
```

**2. 变量模板引擎**
- `{{variable}}` 从 LLM 生成的参数中取值。
- `{{secret:key}}` 从平台密钥管理模块取值（不暴露给 LLM）。
- 支持对变量做转换：`{{customer_name | urlencode}}`。

**3. 执行流程**
```
1. LLM 决定调用 query_customer_crm(customer_name="张三")
2. MCP Server 收到 tool call → 匹配注册的 Tool 定义
3. 取出 request_template，用传入参数渲染模板
4. 从密钥库取出 crm_api_key，拼到 Authorization Header
5. 发送 HTTP POST 到客户 API
6. 返回结果格式化 → 转 TextContent → 返回 LLM
```

**4. 安全边界**
- **网络白名单**：只允许访问客户预先配置的 API 域名/IP 段，拒绝访问内网敏感地址（如 10.x/192.168.x）。
- **超时控制**：每个 Tool 设置最大执行时间（默认 30s），超时后返回结构化错误。
- **返回值截断**：API 返回的 Response 如果过大（如 > 64KB），截断并附加提示"返回值过长，已截断"。
- **密钥安全管理**：`{{secret:xxx}}` 的值由平台加密存储，不出现在日志/LLM Context 中。
- **审计记录**：记录每次 Tool 调用的完整信息（调用方、时间、API URL、参数、返回状态码），用于事后审计。

**5. 错误处理**
- HTTP 4xx → 返回错误码+消息，由 LLM 判断是否可修正参数后重试。
- HTTP 5xx → 根据 retry 策略自动重试，超 max_retries 后返回失败。
- DNS 解析失败/连接超时 → 返回明确的网络错误，供 LLM 或用户判断。

---

**Q29: 某客户私有化环境中有 8 张 A100-80G，需要同时服务 3 个不同的模型（70B 通用对话、13B 代码生成、7B Embedding），如何规划 GPU 资源分配和推理部署？**

**资源评估**：

| 模型 | 参数量 | 量化方案 | 单卡/多卡 | 单实例显存 | 副本数 | 总 GPU 卡数 |
|------|--------|---------|----------|-----------|--------|------------|
| 70B 通用对话 | 70B | AWQ INT4 (~40GB) | 2 卡 Tensor Parallel | ~25GB/卡 | 2 副本 | 4 卡 |
| 13B 代码生成 | 13B | AWQ INT4 (~8GB) | 单卡 | ~12GB(含 KV Cache) | 2 副本 | 2 卡 |
| 7B Embedding | 7B | FP16 (~14GB) | 单卡 | ~15GB(含 batch) | 1 副本 | 1 卡 |
| **合计** | | | | | | **7 卡** |

预留 1 卡做冗余和实验。总计 8 卡刚好。

**部署策略**：

**推理引擎与管理**
- 70B：vLLM，`--tensor-parallel-size 2` + `--quantization awq`，2 副本跨 2 台物理机/跨 NUMA 节点（避免双卡带宽瓶颈）。
- 13B：vLLM + AWQ，单卡即可，2 副本分布在不同 GPU 卡，做负载均衡。
- Embedding：用 TEI（Text Embeddings Inference）或 vLLM Embeddings API，FP16 单卡运行，吞吐优先。

**K8s GPU 调度**
- 使用 NVIDIA GPU Operator + Device Plugin，为不同模型打 nodeSelector 或 affinity。
- 70B 的 Pod 使用 `nvidia.com/gpu: 2` + `topologySpreadConstraints` 保证 2 副本不在同一节点。
- 将 70B 的 4 张卡设为高优先级（`priorityClassName: high-priority`），保证核心模型资源不被抢占。

**模型存储与加载**
- 模型权重存于共享高速存储（节点本地 NVMe SSD 或分布式文件系统），镜像不包含模型文件。
- vLLM 启动时从本地 NVMe 加载（速度快），如需跨节点用 MinIO 分布式模式。
- 权重校验和检查：启动前校验模型文件完整性。

**推理网关路由**
- API Gateway 层按 `model` 参数路由到对应推理引擎的 Service。
- 70B（高成本/高延迟）限制并发数，13B 宽松，Embedding 不限制。
- 按用户等级做优先级调度：VIP 用户优先使用 70B。

**弹性策略**
- 对话模型（70B/13B）：按队列深度 HPA，正常时段 70B 保持 2 副本，低峰缩为 1 副本释放 GPU 给离线批处理。
- Embedding 模型：通常固定副本数，因为批处理需求稳定。

---

**Q30: 如何理解"大模型应用的私有化部署"与传统软件的私有化部署在技术上的核心差异？**

| 维度 | 传统软件私有化 | 大模型应用私有化 |
|------|--------------|----------------|
| **核心依赖** | 依赖关系型数据库 + 中间件 | 新增**GPU 集群** + 模型权重（GB~TB 级）+ 推理引擎 |
| **硬件要求** | CPU + 内存 + 普通磁盘 | 必须有 **NVIDIA GPU**（A100/A10/H100），显存/算力需精确评估 |
| **部署体积** | 镜像几百 MB + 少量数据 | 模型文件动辄 10GB-140GB/个，多模型场景总物料 > 1TB |
| **环境适配** | Linux 发行版适配（CentOS/Ubuntu/麒麟） | 额外需要 **CUDA/cuDNN 驱动版本适配** + GPU 型号兼容矩阵 |
| **升级复杂度** | 代码升级 + 数据库迁移 | 新增**模型权重升级**（换新模型文件可能涉及量化格式变更、显存需求变化） |
| **性能调优** | 数据库调优 + 缓存 + JVM 参数 | 新增**推理引擎参数调优**（TP size、max_num_seqs、KV Cache block size、量化精度） |
| **监控运维** | CPU/内存/磁盘/网络/QPS | 新增**GPU 指标**（利用率/显存/Tensor Core/FLOPS）+ 推理特有指标（TTFT/TPS/队列深度） |
| **高可用** | 服务多副本 + 数据库主从 | 推理服务的高可用成本极高（GPU 昂贵，难以做到传统意义的 N+1 冗余；模型冷启动需数分钟） |
| **容量规划** | 按 QPS 规划 CPU 核数 | 按**并发推理请求数 × 模型大小**规划 GPU 卡数和显存 |

**核心差异总结**：
1. **GPU 是稀缺资源**：无法像 CPU 那样弹性伸缩，容量规划必须精确，且涉及巨额硬件成本。
2. **模型即"巨型依赖"**：模型权重的管理（版本/分发/校验）成为部署的核心挑战，类似于"带了一个小数据库的只读副本"。
3. **运维门槛大幅提高**：需要同时懂 K8s、GPU 驱动、CUDA 生态和推理引擎的复合型工程师。
4. **"慢启动"是常态**：推理服务启动需要加载几十 GB 模型到显存，预热需数分钟，传统快速扩缩容模式不适用。

---

*本文档为秒哒产品大模型应用架构面试专题，共 30 题，覆盖 LLM 应用架构基础、RAG 深度解析、Agent 架构设计、MCP 协议实战、推理引擎与模型服务化、实战综合场景六大领域。*
