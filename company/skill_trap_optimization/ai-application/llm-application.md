# AI 项目应用落地（细化版）

## 一、RAG（检索增强生成）深入

### 1.1 完整 Pipeline

```
文档 → Chunk 切分 → Embedding → 向量数据库
                                    ↓
用户提问 → Embedding → 相似检索 → Rerank → Prompt拼接 → LLM生成 → 返回
```

### 1.2 Chunk 切分策略

| 策略 | 做法 | 适用 |
|------|------|------|
| **固定长度** | 按 token 数切分，重叠 N token | 通用场景，简单高效 |
| **语义切分** | 按段落/章节，或用小模型判断语义边界 | 对上下文连续性要求高 |
| **递归切分** | 先用大分隔符（\\n\\n），不够再逐级细化 | LangChain 默认策略 |
| **句子级切分** | 按句号/换行切分 | FAQ 类短文本 |

```go
// 固定长度 + 重叠切分
func chunkText(text string, chunkSize, overlap int) []Chunk {
    var chunks []Chunk
    runes := []rune(text)
    for i := 0; i < len(runes); i += chunkSize - overlap {
        end := min(i+chunkSize, len(runes))
        chunks = append(chunks, Chunk{
            Content: string(runes[i:end]),
            Index:   len(chunks),
        })
    }
    return chunks
}
```

### 1.3 Embedding 模型选型

| 模型 | 维度 | 中文效果 | 成本 | 部署 |
|------|------|----------|------|------|
| OpenAI text-embedding-3-small | 512/1536 | 好 | 低（$0.02/1M token） | API |
| OpenAI text-embedding-3-large | 256/1024/3072 | 很好 | 高（$0.13/1M token） | API |
| BGE-M3 (BAAI) | 1024 | 很好 | 免费 | 本地 GPU |
| M3E (moka-ai) | 768 | 好 | 免费 | 本地 CPU 可跑 |
| Cohere Embed v3 | 1024 | 中 | 中 | API |
| Jina Embeddings v2 | 768 | 好 | 中 | API |

**选型建议**：
- 预算足 + 追求效果 → OpenAI text-embedding-3-large
- 预算有限 + 中文为主 → BGE-M3 本地部署
- 轻量 + CPU 部署 → M3E

### 1.4 向量数据库选型

| 数据库 | 类型 | 适用 | 核心优势 |
|--------|------|------|----------|
| **pgvector** | PostgreSQL 扩展 | 已有 PG 的项目，小规模 | 零额外运维，SQL 即可 |
| **Milvus** | 专用向量库 | 大规模，百万级以上 | 性能最高，分布式架构 |
| **Qdrant** | 专用向量库 | 中等规模，性能好 | Rust 编写，过滤强大 |
| **Weaviate** | 专用向量库 | 需要混合搜索 | 内建 GraphQL，生态丰富 |
| **Chroma** | 嵌入式 | 原型开发 | 最轻量，Python 优先 |
| **Elasticsearch + dense_vector** | 搜索引擎扩展 | 已有 ES 设施 | 混合全文+向量搜索 |

```go
// pgvector 示例
import "github.com/pgvector/pgvector-go"

type Document struct {
    ID        int64
    Content   string
    Embedding pgvector.Vector
}

// 插入向量
_, err := db.Exec("INSERT INTO docs (content, embedding) VALUES ($1, $2)",
    chunk.Content, pgvector.NewVector(embedding))

// 相似检索
rows, err := db.Query(`
    SELECT content, 1 - (embedding <=> $1) AS similarity
    FROM docs ORDER BY embedding <=> $1 LIMIT 10`,
    pgvector.NewVector(queryEmbedding))
```

### 1.5 Rerank（重排序）

初检索返回的 Top-K 并不一定最相关，通过 Rerank 模型二次排序，提高精度。

| 方案 | 模型 | 特点 |
|------|------|------|
| Cohere Rerank | command-rerank | API 调用，效果好 |
| BGE-Reranker | BAAI/bge-reranker | 本地部署，中文好 |
| Cross-Encoder | sentence-transformers | 精度高但慢 |

```
初检索 Top-100 → Rerank → 取 Top-5 → 拼入 Prompt → LLM
```

### 1.6 RAG 完整代码示例（Go 实现检索端）

```go
func (s *RAGService) Query(ctx context.Context, userQuery string) (string, error) {
    // 1. 查询向量化
    queryVec, err := s.embedder.Embed(ctx, userQuery)
    if err != nil {
        return "", fmt.Errorf("embed query: %w", err)
    }

    // 2. 向量检索
    docs, err := s.vectorDB.Search(ctx, queryVec, 10)
    if err != nil {
        return "", fmt.Errorf("search: %w", err)
    }

    // 3. Rerank (可选，提升精度)
    docs = s.reranker.Rerank(ctx, userQuery, docs, 5)

    // 4. 拼接 Prompt
    prompt := s.buildRAGPrompt(userQuery, docs)

    // 5. LLM 生成
    answer, err := s.llm.Generate(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("generate: %w", err)
    }
    return answer, nil
}

func (s *RAGService) buildRAGPrompt(query string, docs []Document) string {
    var ctx strings.Builder
    for i, doc := range docs {
        fmt.Fprintf(&ctx, "[%d] %s\n\n", i+1, doc.Content)
    }
    return fmt.Sprintf(`基于以下参考资料回答问题，如果资料中没有答案就说不知道。

参考资料：
%s

问题：%s

回答：`, ctx.String(), query)
}
```

---

## 二、Agent 模式深入

### 2.1 ReAct 模式（Reasoning + Acting）

Agent 最核心的模式，LLM 在"思考"和"行动"之间交替循环：

```
用户输入 → [思考] → [调用工具] → [观察结果] → [思考] → [返回最终答案]
```

```go
type Agent struct {
    llm    LLM
    tools  map[string]Tool
    maxIter int // 最大循环次数，防止死循环
}

func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
    messages := []Message{{Role: "user", Content: userInput}}

    for i := 0; i < a.maxIter; i++ {
        resp, err := a.llm.Chat(ctx, messages, a.toolDefs())
        if err != nil {
            return "", err
        }

        if resp.ToolCalls == nil {
            return resp.Content, nil // 最终回复
        }

        // 执行工具调用
        for _, tc := range resp.ToolCalls {
            result := a.executeTool(ctx, tc)
            messages = append(messages, Message{
                Role: "tool", ToolCallID: tc.ID, Content: result,
            })
        }
    }
    return "", fmt.Errorf("exceeded max iterations")
}
```

### 2.2 Function Calling 工具定义

```go
type Tool struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Parameters  json.RawMessage `json:"parameters"`
    Execute     func(ctx context.Context, args map[string]any) (string, error)
}

// 示例：查询天气工具
weatherTool := Tool{
    Name:        "get_weather",
    Description: "获取指定城市的当前天气",
    Parameters: toJSON(map[string]any{
        "type": "object",
        "properties": map[string]any{
            "city": map[string]any{
                "type": "string", "description": "城市名称，如 Beijing",
            },
        },
        "required": []string{"city"},
    }),
    Execute: func(ctx context.Context, args map[string]any) (string, error) {
        city := args["city"].(string)
        // 调用天气 API
        return fetchWeather(city), nil
    },
}
```

### 2.3 Agent 记忆管理

| 记忆类型 | 作用 | 实现 |
|----------|------|------|
| **短期记忆** | 当前对话上下文 | messages 数组（窗口截断） |
| **长期记忆** | 跨会话信息 | 向量库存储 + 检索 |
| **摘要记忆** | 压缩历史对话 | 用 LLM 摘要旧消息 |

```go
// 对话窗口管理：总 token 数超限时，摘要最旧的消息
func (m *Memory) Truncate(messages []Message, maxTokens int) []Message {
    totalTokens := countTokens(messages)
    if totalTokens <= maxTokens {
        return messages
    }
    // 保留最近 N 轮对话 + 摘要旧消息
    cutoff := len(messages) * maxTokens / totalTokens
    oldSummarized := m.summarize(messages[:cutoff])
    recent := messages[cutoff:]
    return append([]Message{{Role: "system", Content: "对话历史摘要：" + oldSummarized}}, recent...)
}
```

---

## 三、MCP Server 深入

### 3.1 协议概念

MCP（Model Context Protocol）是一种标准协议，让 LLM 以统一方式访问外部工具和数据源。

```
LLM Client ──JSON-RPC──→ MCP Server ──→ 工具/资源

三个核心原语：
- Resource：暴露数据（文件、数据库记录、API 数据）
- Tool：暴露操作（执行 SQL、发邮件、查询天气）
- Prompt：预设 Prompt 模板
```

### 3.2 Go 实现 MCP Server

```go
import "github.com/mark3labs/mcp-go"

func main() {
    server := mcp.NewServer("my-service", "1.0.0")

    // 注册工具
    server.AddTool(mcp.Tool{
        Name: "query_database",
        Description: "执行 SQL 查询（仅 SELECT）",
        InputSchema: mcp.ToolInputSchema{
            Type: "object",
            Properties: map[string]any{
                "sql": map[string]any{"type": "string"},
            },
        },
    }, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        sql := req.Params.Arguments["sql"].(string)
        result, err := executeQuery(sql)
        if err != nil {
            return nil, err
        }
        return &mcp.CallToolResult{
            Content: []mcp.Content{mcp.NewTextContent(result)},
        }, nil
    })

    // 注册资源
    server.AddResourceTemplate(mcp.ResourceTemplate{
        URITemplate: "docs://{doc_id}",
        Name: "document",
    }, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
        docID := req.Params.Arguments["doc_id"].(string)
        content := loadDocument(docID)
        return []mcp.ResourceContents{mcp.TextResourceContents{Text: content}}, nil
    })

    // 启动 SSE 服务
    sseServer := mcp.NewSSEServer(server)
    sseServer.Start(":8080")
}
```

---

## 四、流式输出实现

### 4.1 SSE (Server-Sent Events)

```go
func StreamChat(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    flusher := w.(http.Flusher)

    stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
    if err != nil {
        return
    }
    defer stream.Close()

    for stream.Next() {
        delta := stream.Current().Choices[0].Delta.Content
        fmt.Fprintf(w, "data: %s\n\n", delta)
        flusher.Flush()
    }
    fmt.Fprintf(w, "data: [DONE]\n\n")
    flusher.Flush()
}
```

### 4.2 WebSocket（双向通信）

```go
var upgrader = websocket.Upgrader{}

func WsChat(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        // 流式回复
        stream, _ := openaiClient.CreateChatCompletionStream(ctx, req)
        for stream.Next() {
            delta := stream.Current().Choices[0].Delta.Content
            conn.WriteJSON(map[string]string{"delta": delta})
        }
        conn.WriteJSON(map[string]string{"done": "true"})
    }
}
```

---

## 五、Prompt 工程

### 5.1 技巧清单

| 技巧 | 做法 | 效果 |
|------|------|------|
| **Few-Shot** | 提供 2-5 个示例 | 输出格式可控，准确率提升 |
| **Chain-of-Thought** | 加 "Let's think step by step" | 推理题准确率显著提升 |
| **角色设定** | "你是资深 Go 工程师..." | 回答风格符合角色 |
| **结构化输出** | "以 JSON 格式返回" | 程序可解析 |
| **思维树** (ToT) | 多路径推理 + 投票 | 复杂决策题 |
| **Self-Consistency** | 多次采样 + 投票 | CoT 增强，降低单次错误 |

### 5.2 Prompt 模板管理

```go
// 模板与代码分离
const ragPromptTpl = `基于以下参考资料回答问题，如果资料中没有答案就说不知道。

参考资料：
{{range .Docs}}
[{{.Index}}] {{.Content}}
{{end}}

问题：{{.Query}}

回答：`

func buildPrompt(query string, docs []Doc) string {
    tmpl := template.Must(template.New("rag").Parse(ragPromptTpl))
    var buf strings.Builder
    tmpl.Execute(&buf, map[string]any{
        "Query": query,
        "Docs":  docs,
    })
    return buf.String()
}
```

### 5.3 A/B 测试 Prompt

```
API 请求时附带 prompt_version 参数
→ 记录到日志/数据库
→ 对比不同版本的成功率、用户满意度
→ 逐步优化
```

---

## 六、模型路由与网关

### 6.1 模型网关架构

```
客户端 → 统一 API Gateway
            ├── 路由层：根据请求类型/成本/延迟 → 选择模型
            ├── 限流层：用户级/模型级配额
            ├── 缓存层：相似问题返回缓存结果
            ├── 降级层：主模型不可用时切换备用模型
            └── 计费层：Token 用量统计
              ├── GPT-4（复杂推理）
              ├── GPT-4o-mini（常规对话）
              ├── Claude Opus（代码生成）
              └── 本地模型（敏感数据不出域）
```

```go
type ModelRouter struct {
    models map[string]Model
    rules  []RoutingRule
}

type RoutingRule struct {
    Condition func(req ChatRequest) bool
    Model     string
}

func (r *ModelRouter) Route(req ChatRequest) Model {
    for _, rule := range r.rules {
        if rule.Condition(req) {
            return r.models[rule.Model]
        }
    }
    return r.models["default"]
}
```

### 6.2 Sematic Cache（语义缓存）

相似问题不重复调 LLM，返回缓存结果：

```
用户问题 → Embedding → 向量库查相似问题 → 相似度 > 阈值 → 返回缓存答案
                                                     └── < 阈值 → 调 LLM → 缓存结果
```

```go
func (s *SemanticCache) Get(query string) (string, bool) {
    queryVec := s.embedder.Embed(query)
    results := s.store.Search(queryVec, 1)
    if len(results) > 0 && results[0].Score > s.threshold {
        return results[0].Answer, true
    }
    return "", false
}
```

---

## 七、成本控制

| 措施 | 效果 |
|------|------|
| 语义缓存 | 相似问题命中率可达 30%，直接节省 30% 成本 |
| 分级模型路由 | 简单问题用廉价模型（gpt-4o-mini），复杂问题才用强模型 |
| 限制输出长度 | `max_tokens` 设合理上限 |
| Prompt 压缩 | 去掉冗余的 system prompt，用 LLMLingua 压缩上下文 |
| 批处理 | Batch API（OpenAI 50% 折扣，24h 内返回） |
| 监控告警 | 单用户 Token 消耗异常告警，防止滥用 |

---

## 八、安全防护

### 8.1 Prompt 注入防护

```go
func sanitizeInput(input string) string {
    // 1. 检测注入特征
    patterns := []string{
        "ignore previous instructions",
        "忽略之前的指令",
        "system prompt",
        "you are now",
    }
    for _, p := range patterns {
        if strings.Contains(strings.ToLower(input), p) {
            return "[检测到注入尝试，已拒绝]"
        }
    }
    // 2. 用户输入用特殊标记包裹，与 system prompt 区分
    return fmt.Sprintf("<user_input>\n%s\n</user_input>", input)
}
```

### 8.2 内容安全

- 输入过滤：敏感词、PII（身份证、手机号）脱敏
- 输出审核：调用审核 API（OpenAI Moderation / 本地模型）检测有害内容
- 越狱检测：检测提示越狱的攻击模式

---

## 九、评估体系

| 评估维度 | 方法 | 指标 |
|----------|------|------|
| **检索质量** | 人工标注 + 自动评估 | Recall@K, MRR, NDCG |
| **生成质量** | LLM-as-Judge / 人工评分 | Faithfulness（忠实度）、Relevance |
| **端到端效果** | 用户反馈（👍/👎）、A/B Test | 转化率、留存率、满意度 |
| **性能** | 基准测试 | 首 token 延迟、生成速度 (tokens/s) |

```go
// RAG 评估示例
func evaluateRAG(testCases []TestCase) Metrics {
    var totalRecall, totalMRR float64
    for _, tc := range testCases {
        results := rag.Search(tc.Query, 10)
        // Recall@5: 前5个结果中命中的比例
        hits := intersect(results[:5], tc.RelevantDocIDs)
        totalRecall += float64(len(hits)) / float64(len(tc.RelevantDocIDs))
        // MRR: 第一个命中的排名的倒数
        for i, r := range results {
            if contains(tc.RelevantDocIDs, r.ID) {
                totalMRR += 1.0 / float64(i+1)
                break
            }
        }
    }
    return Metrics{
        RecallAt5: totalRecall / float64(len(testCases)),
        MRR:       totalMRR / float64(len(testCases)),
    }
}
```
