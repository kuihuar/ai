# LLM 选型与解析版本匹配方案

## 第一部分：题目解答场景的 LLM 选型

### 一、国内外可选 LLM 全景

#### 1.1 国外模型

| 模型 | 提供方 | 推理能力 | 数学 | 代码 | 成本(输入/输出 $/1M token) | 中文 | 最佳场景 |
|------|--------|----------|------|------|---------------------------|------|----------|
| **GPT-4o** | OpenAI | ★★★★★ | ★★★★★ | ★★★★★ | 2.5/10 | ★★★★ | 全科最强，复杂推理首选 |
| **GPT-4o-mini** | OpenAI | ★★★ | ★★★ | ★★★ | 0.15/0.6 | ★★★★ | 简单题、POI 分类、文本润色 |
| **Claude Opus 4** | Anthropic | ★★★★★ | ★★★★ | ★★★★★ | 15/75 | ★★★★ | 长文本推理、证明题 |
| **Claude Sonnet 4** | Anthropic | ★★★★ | ★★★★ | ★★★★★ | 3/15 | ★★★★ | 性价比最好的推理模型 |
| **Claude Haiku 4** | Anthropic | ★★★ | ★★★ | ★★★ | 0.8/4 | ★★★★ | 简单任务、快速分类 |
| **Gemini 2.5 Pro** | Google | ★★★★★ | ★★★★★ | ★★★★ | 1.25/10 | ★★★ | 数学极强，中文略弱 |
| **Gemini 2.5 Flash** | Google | ★★★★ | ★★★★ | ★★★★ | 0.15/0.6 | ★★★ | 高性价比 |

#### 1.2 国内模型

| 模型 | 提供方 | 推理能力 | 数学 | 中文 | 成本(元/百万 token) | 最佳场景 |
|------|--------|----------|------|------|---------------------|----------|
| **DeepSeek-V3** | DeepSeek | ★★★★ | ★★★★★ | ★★★★★ | ~2/8 | 数学推理极强，成本低 |
| **DeepSeek-R1** | DeepSeek | ★★★★★ | ★★★★★ | ★★★★★ | ~4/16 | 最强调推理，慢但准确 |
| **Qwen3-235B** | 阿里 | ★★★★ | ★★★★ | ★★★★★ | ~4/8 | 中文综合最强 |
| **Qwen3-Max** | 阿里 | ★★★★★ | ★★★★★ | ★★★★★ | ~10/40 | 旗舰，效果接近 GPT-4o |
| **Kimi K2** | 月之暗面 | ★★★★★ | ★★★★ | ★★★★★ | ~2/8 | 超长上下文，适合大题 |
| **GLM-4** | 智谱 | ★★★★ | ★★★★ | ★★★★ | ~1/2 | 性价比高 |
| **ERNIE 4.5** | 百度 | ★★★★ | ★★★★ | ★★★★★ | ~8/24 | 百科知识丰富 |
| **Doubao-Pro** | 字节 | ★★★★ | ★★★ | ★★★★★ | ~1/2 | 语文/英语类题目 |
| **Step-3** | 阶跃星辰 | ★★★★ | ★★★★ | ★★★★ | ~2/8 | 多模态、图片理解 |

#### 1.3 专用模型

| 模型 | 用途 | 适用 |
|------|------|------|
| **Math-Σ** (DeepSeek数学专用) | 数学题推导 | 复杂数学计算 |
| **Qwen2.5-Math** | 数学推理 | 数学竞赛题 |
| **Pix2Text** | OCR + 公式识别 | 图片→文本+LaTeX |
| **BGE-M3** | Embedding | 题目检索 |
| **BGE-Reranker** | 重排序 | 检索结果精排 |

---

### 二、按题型/科目的模型路由规划

#### 2.1 路由矩阵

```
用户拍照 → 题型分类
              │
              ├── 数学计算题（方程、求导、积分）
              │       ├── 简单 → DeepSeek-V3（成本低，数学强）
              │       ├── 中等 → GPT-4o-mini + SymPy 验证
              │       └── 难题 → DeepSeek-R1 深度推理（时间长但准）
              │
              ├── 数学证明题
              │       ├── 首选 → DeepSeek-R1（深度推理链）
              │       └── 备选 → Claude Opus（逻辑严谨）
              │
              ├── 物理解答题
              │       ├── 首选 → GPT-4o / Qwen3-Max
              │       └── 备选 → Claude Sonnet
              │
              ├── 化学题
              │       ├── 配平计算 → 符号引擎（不用 LLM）
              │       ├── 概念解释 → GLM-4（成本低）
              │       └── 实验推断 → GPT-4o-mini
              │
              ├── 英语题
              │       ├── 语法选择 → GLM-4 + 模板（成本最低）
              │       ├── 阅读理解 → Kimi K2（长文本好）
              │       └── 写作 → Claude Sonnet（文字质量高）
              │
              ├── 语文题
              │       ├── 文言文 → Qwen3-Max（中文底子好）
              │       ├── 阅读 → Kimi K2（长上下文）
              │       └── 作文 → Claude Sonnet
              │
              └── 生物/历史/地理
                      └── RAG 检索 + GPT-4o-mini / GLM-4（总结型任务）
```

#### 2.2 路由决策引擎

```go
type ModelRouter struct {
    models  map[string]Model
    rules   []RouteRule
}

type RouteRule struct {
    // 匹配条件
    Subject     string   // 学科
    Difficulty  string   // 难度：easy/medium/hard
    QuestionType string  // 题型：calculation/proof/concept/reading
    MaxCost     float64  // 成本上限（元/次）

    // 路由目标
    PrimaryModel   string  // 首选模型
    FallbackModel  string  // 主模型不可用时降级
    RetryModel     string  // 答案验证失败时重试
}

func (r *ModelRouter) Route(q Question) (*RouteDecision, error) {
    for _, rule := range r.rules {
        if rule.match(q) {
            return &RouteDecision{
                Primary:  r.models[rule.PrimaryModel],
                Fallback: r.models[rule.FallbackModel],
                Retry:    r.models[rule.RetryModel],
            }, nil
        }
    }
    // 兜底：用最便宜的模型
    return r.defaultRoute(), nil
}
```

#### 2.3 分级路由的成本对比

假设日均 100 万题，各难度占比及模型分配：

| 难度 | 占比 | 模型选择 | 单次成本 | 日成本 |
|------|------|----------|----------|--------|
| 简单 | 40% | GLM-4 / GPT-4o-mini | ~0.002 元 | 800 元 |
| 中等 | 35% | DeepSeek-V3 / Qwen3 | ~0.01 元 | 3,500 元 |
| 困难 | 15% | DeepSeek-R1 / GPT-4o | ~0.05 元 | 7,500 元 |
| 题库命中(无LLM) | 10% | 零成本 | 0 元 | 0 元 |
| **合计** | | | | **~11,800 元/天** |

如果全用 GPT-4o：100 万 × 0.03 元 = **30,000 元/天**。分级路由节省 60%。

---

### 三、LLM 选型的核心权衡

```go
// 选型五要素
type ModelSelectionCriteria struct {
    Accuracy    float64 // 正确率权重
    Cost        float64 // 成本权重
    Latency     float64 // 延迟权重
    Stability   float64 // 服务稳定性（SLA）
    Compliance  float64 // 合规（数据出境、内容安全）
}

// 场景化权重
var weights = map[string]ModelSelectionCriteria{
    "user_facing": {  // 用户实时请求
        Accuracy:    0.40,
        Cost:        0.20,
        Latency:     0.25,  // 延迟很敏感
        Stability:   0.10,
        Compliance:  0.05,
    },
    "batch_gen":   {  // 批量生成解析
        Accuracy:    0.45,
        Cost:        0.30,  // 成本更敏感
        Latency:     0.05,  // 延迟不敏感
        Stability:   0.10,
        Compliance:  0.10,
    },
    "verification":{  // 答案验证
        Accuracy:    0.50,  // 正确率最重要
        Cost:        0.15,
        Latency:     0.10,
        Stability:   0.15,
        Compliance:  0.10,
    },
}
```

**场景影响选型的关键例子**：

- **用户实时请求**：必须低延迟 → 用 DeepSeek-V3（快+便宜+数学好）而不是 DeepSeek-R1（好但慢 10 倍）
- **题库批量生成**：延迟不敏感 → 用 DeepSeek-R1 + 多路投票，追求最高质量
- **教育行业合规**：涉及学生数据的国内产品 → 必须用国内模型（DeepSeek/Qwen/GLM），不能调 OpenAI

---

## 第二部分：解析版本管理——同一道题多个解析，哪个更好？

### 一、为什么一道题会有多个解析版本

| 来源 | 说明 |
|------|------|
| **自有录入 v1** | 最初录入时老师写的解析 |
| **LLM 生成 v2** | AI 批量生成的标准化解析 |
| **用户纠错 v3** | 用户反馈后修正的版本 |
| **抓取源 A** | 从网站 A 抓的，质量未知 |
| **抓取源 B** | 从网站 B 抓的，可能解法不同 |
| **一题多解 v4** | 另一种解题方法（代数法 vs 几何法） |
| **分层版本** | 基础版/进阶版，针对不同年级 |

### 二、解析版本的数据模型

```sql
-- 题目主表
CREATE TABLE questions (
    id          BIGINT PRIMARY KEY,
    fingerprint VARCHAR(64) NOT NULL,  -- 题目指纹（用于去重）
    subject     VARCHAR(32),
    type        VARCHAR(32),           -- choice/fill/solution/proof
    difficulty  VARCHAR(16),           -- easy/medium/hard
    content     TEXT,                  -- 题目文本
    content_hash VARCHAR(64),          -- 内容哈希
    created_at  TIMESTAMP,
    INDEX idx_fingerprint (fingerprint)
);

-- 解析版本表（核心）
CREATE TABLE solution_versions (
    id            BIGINT PRIMARY KEY,
    question_id   BIGINT NOT NULL REFERENCES questions(id),
    version       INT NOT NULL,          -- 版本号 1,2,3...
    source        VARCHAR(32),           -- manual/llm_gpt4/llm_deepseek/scraped/user_fix
    method        VARCHAR(32),           -- 解法：standard/alternative_1/alternative_2
    difficulty_level VARCHAR(16),        -- basic/intermediate/advanced（针对不同年级）

    -- 内容
    analysis      TEXT,                  -- 考点分析
    steps         JSONB,                 -- [{"step":1,"title":"...","content":"...","reasoning":"..."}]
    answer        TEXT,                  -- 最终答案
    tips          TEXT,                  -- 易错提醒
    summary       TEXT,                  -- 方法总结

    -- 质量元数据
    quality_score DECIMAL(3,1),          -- 综合质量分 0-100
    auto_verified BOOLEAN DEFAULT FALSE, -- 是否通过自动验收
    reviewed_by   BIGINT,                -- 审核人 ID
    review_status VARCHAR(16),           -- pending/approved/rejected

    -- 效果数据
    impression_count BIGINT DEFAULT 0,   -- 展示次数
    click_count      BIGINT DEFAULT 0,   -- 用户展开查看次数
    like_count       BIGINT DEFAULT 0,   -- 👍
    dislike_count    BIGINT DEFAULT 0,   -- 👎
    report_count     BIGINT DEFAULT 0,   -- 报错次数
    avg_view_time_ms BIGINT DEFAULT 0,   -- 平均查看时长（能看完=质量好）

    -- 模型元数据
    llm_model     VARCHAR(64),           -- 使用的 LLM 模型
    llm_prompt_version VARCHAR(32),      -- Prompt 版本
    llm_tokens    INT,                   -- Token 消耗
    generation_cost DECIMAL(10,6),

    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP,

    INDEX idx_question_id (question_id),
    INDEX idx_quality (question_id, quality_score DESC)
);

-- 解析效果日志（用于 AB 对比）
CREATE TABLE solution_ab_log (
    id            BIGINT PRIMARY KEY,
    question_id   BIGINT,
    user_id       BIGINT,
    solution_id   BIGINT,               -- 展示给用户的版本
    action        VARCHAR(16),          -- shown/expanded/liked/disliked/reported
    dwell_time_ms BIGINT,               -- 停留时长
    created_at    TIMESTAMP,
    INDEX idx_question_user (question_id, user_id)
);
```

### 三、解析质量评分系统

#### 3.1 静态质量分（入库时计算）

```python
def calculate_static_quality(solution: dict, question: dict) -> float:
    """入库时计算的静态质量分"""
    scores = {}

    # 1. 答案正确性（30 分）
    if question["type"] == "equation":
        scores["correctness"] = 30 if verify_with_sympy(solution) else 0
    elif question["type"] == "proof":
        scores["correctness"] = 15  # 证明题无法自动验证，给基准分
    else:
        scores["correctness"] = 20  # 其他题型默认分，人工审核后调整

    # 2. 步骤完整性（25 分）
    step_count = len(solution["steps"])
    has_final_answer = bool(solution.get("answer"))
    scores["completeness"] = min(25,
        (step_count >= 3) * 15 +      # 至少 3 步
        has_final_answer * 10         # 有最终答案
    )

    # 3. 教学价值（20 分）
    scores["teaching"] = (
        bool(solution.get("analysis")) * 5 +    # 考点分析
        bool(solution.get("tips")) * 5 +        # 易错提醒
        bool(solution.get("summary")) * 5 +     # 方法总结
        (step_count >= 5) * 5                   # 步骤详细
    )

    # 4. 格式规范（15 分）
    scores["formatting"] = (
        latex_well_formed(solution) * 5 +
        step_numbering_consistent(solution) * 5 +
        has_clear_answer_label(solution) * 5
    )

    # 5. LLM-as-Judge 评分（10 分）
    scores["llm_judge"] = llm_score_solution(question, solution) * 0.1

    return sum(scores.values())
```

#### 3.2 动态效果分（运行时更新）

```python
def calculate_dynamic_score(solution_id: int) -> float:
    """基于用户行为计算动态效果分"""
    stats = db.query("""
        SELECT
            impression_count,
            click_count,
            like_count, dislike_count, report_count,
            avg_view_time_ms
        FROM solution_versions WHERE id = ?
    """, solution_id)

    # 最少展示量阈值（数据量够才有统计意义）
    if stats["impression_count"] < 100:
        return 0.5  # 数据不足，给中性分

    # 点击率（展示了是否愿意展开看）
    ctr = stats["click_count"] / stats["impression_count"]

    # 好评率
    total_feedback = stats["like_count"] + stats["dislike_count"]
    like_rate = stats["like_count"] / total_feedback if total_feedback > 0 else 0.5

    # 报错率（反向指标）
    report_rate = stats["report_count"] / stats["impression_count"]

    # 完读率（查看时长 > 解析预期时长的比例）
    expected_read_ms = 30000  # 假设读完需要 30 秒
    completion_rate = stats["avg_view_time_ms"] / expected_read_ms

    # 综合动态分
    dynamic_score = (
        ctr * 0.2 +
        like_rate * 0.35 +
        (1 - report_rate * 10) * 0.30 +  # 报错权重很大
        min(completion_rate, 1.0) * 0.15
    )

    return dynamic_score
```

#### 3.3 综合排名分

```python
def calculate_final_rank(solution_id: int) -> float:
    """综合静态质量 + 动态效果分"""
    static = get_static_quality(solution_id)
    dynamic = get_dynamic_score(solution_id)
    freshness = get_freshness_factor(solution_id)  # 越新的版本有微小加成
    manual_bonus = 1.1 if is_manually_reviewed(solution_id) else 1.0

    # 权重：静态 0.4 + 动态 0.5 + 新鲜度 0.1
    base = static * 0.4 + dynamic * 0.5 + freshness * 0.1
    return base * manual_bonus
```

### 四、如何选出"最适合"的解析

#### 4.1 策略：默认最优 + 用户匹配 + AB 测试

```go
type SolutionSelector struct {
    db    *sql.DB
    cache *redis.Client
}

func (s *SolutionSelector) SelectBest(
    ctx context.Context,
    questionID int64,
    user *UserProfile,
) (*Solution, string) {

    // 1. 获取该题所有候选版本（按综合分降序）
    candidates := s.getCandidates(questionID) // ORDER BY quality_score DESC

    if len(candidates) == 0 {
        return nil, "no_solution"
    }

    // 2. 如果有审核通过的版本，优先展示
    for _, c := range candidates {
        if c.ReviewStatus == "approved" {
            return c, "manually_approved"
        }
    }

    // 3. 如果多解，根据用户画像选
    if len(candidates) >= 2 && user != nil {
        // 基础生 → 显示"标准解法"
        if user.Level == "basic" {
            return s.filterByLevel(candidates, "basic"), "user_level_match"
        }
        // 尖子生 → 显示"进阶解法"或"一题多解"
        if user.Level == "advanced" {
            return s.filterByLevel(candidates, "advanced"), "user_level_match"
        }
    }

    // 4. 如果存在 AB 测试，随机分配
    if s.isABTestActive(questionID) {
        return s.assignABTest(questionID, user.ID), "ab_test"
    }

    // 5. 默认：返回综合排名最高、且不是 LLM 初次生成的版本
    for _, c := range candidates {
        if c.Source != "llm_unverified" {
            return c, "highest_ranked"
        }
    }

    // 6. 兜底：返回最高分（即使未审核）
    return candidates[0], "fallback"
}
```

#### 4.2 多解展示：不只选一个

同一道题，用户可能想看到不同解法。策略不是单选，而是**分层展示**：

```
┌─────────────────────────────────────────┐
│  📝 题目：求函数 f(x)=x²+2x+1 的最小值    │
├─────────────────────────────────────────┤
│  ⭐ 推荐解析（综合评分 92，987人👍）      │
│  配方法 → f(x)=(x+1)² ≥ 0 → 最小值为 0    │
│  标签：通法 · 适合所有人                   │
├─────────────────────────────────────────┤
│  📐 其他解法                             │
│  ├─ 求导法（评分 88）- 适合学过导数的     │
│  ├─ 图像法（评分 85）- 数形结合更直观     │
│  └─ 判别式法（评分 82）- 竞赛常用         │
├─────────────────────────────────────────┤
│  这道题有帮助吗？ 👍 987 · 👎 23          │
└─────────────────────────────────────────┘
```

```sql
-- 一题多解的去重和排序
SELECT method, MAX(quality_score) as best_score
FROM solution_versions
WHERE question_id = ?
  AND quality_score > 80
GROUP BY method
ORDER BY best_score DESC;
```

#### 4.3 不同场景的选取策略

| 场景 | 选取策略 | 原因 |
|------|----------|------|
| **拍照搜题结果** | 最高分 single | 用户要的是正确答案，不要选择题 |
| **题目详情页** | 最高分为主 + 备选折叠 | 用户可以探索多种解法 |
| **错题本回顾** | 优先选"用户上一次看到的版本" | 保持一致性 |
| **打印导出** | 最高分 + 最详细（步骤最多的） | 离线阅读需要详尽 |
| **API 输出给第三方** | 最高分 + 必须人工审核 | 对外的必须经过审核 |
| **新题目冷启动** | 展示 LLM 版本 + 标记"AI生成" | 没有更好的，但要诚实标注 |

---

## 第三部分：LLM 参与解析版本匹配的方式

### 一、LLM 在版本选择中的作用

LLM 本身也可以参与"判断哪个版本更好"，实现方式：

```
同一道题有 3 个候选解析：
  解析 A：配方法，来自自有题库，人工审核
  解析 B：求导法，来自 GPT-4o 生成，自动验收通过
  解析 C：抓取自网站 X，格式粗糙但有图解

LLM-as-Judge 评估：
  输入：题目 + 3 个候选解析
  输出：排序 + 每份的优缺点 + 是否可以合并
```

```python
def llm_judge_solutions(question: str, solutions: list) -> dict:
    prompt = f"""你是一位资深教研员。请评估以下题目的 {len(solutions)} 个解析版本。

题目：{question}

{solutions_text}

请从以下维度评估并排序：
1. 正确性：答案和推导是否正确
2. 教学性：是否有助于学生理解
3. 完整性：步骤是否完整无跳跃
4. 可读性：格式、公式、排版

返回 JSON：
{{
  "ranking": [
    {{"rank": 1, "solution_id": "B", "score": 92, "strength": "...", "weakness": "..."}},
    ...
  ],
  "best_for_display": "A",  // 最适合展示的
  "can_merge": true,         // 是否可以合并优点
  "merged_suggestion": "..." // 如果能合并，给出建议
}}
"""
    return llm.generate(prompt)
```

**使用时机**：只在有新版本入库、或用户大量 👎 时触发 LLM 重新评估。不是每次请求都调。

### 二、LLM 合并多个解析版本（取各自最好的部分）

```python
def llm_merge_solutions(question: str, solutions: list) -> str:
    """多版本解析 → LLM 合并为最优版本"""
    prompt = f"""以下是一道题的多个解析版本，请你取各自的优点，合并为一个最佳解析。

题目：{question}

{solutions_text}

合并不等于简单拼接。请遵循：
1. 如果多个版本的解法相同，选表达最好的那份
2. 如果解法不同，保留"推荐解析"，"其他解法"简要提要点
3. 取每个版本中最好的部分（谁的考点分析好？谁的步骤更清晰？谁的易错提醒好？）
4. 修正全部版本中的任何错误
5. 统一格式和 LaTeX

输出合并后的标准解析：
"""
    return llm.generate(prompt)
```

**触发条件**：
- 该题有 ≥ 2 个版本且质量分都在 70-85 之间（各有长处，但都不完美）
- 用户报错率 > 5%
- 人工审核员触发"智能合并"操作

---

## 第四部分：整体架构设计

### 4.1 系统全貌

```
                         ┌──────────────┐
                         │  题目入口     │
                         │ (拍照/文本)   │
                         └──────┬───────┘
                                ↓
                    ┌───────────────────────┐
                    │    题目指纹 + 去重     │
                    └───────────┬───────────┘
                                ↓
                    ┌───────────────────────┐
                    │   已有解析版本？       │
                    └────┬────────────┬─────┘
                       Yes            No
                        ↓              ↓
              ┌─────────────┐   ┌──────────────┐
              │  版本选择器  │   │  解析生成器   │
              │              │   │              │
              │ 静态评分     │   │ 模型路由器    │
              │ 动态效果分   │   │ 4 种生成方案  │
              │ AB 测试      │   │ 自动验收      │
              │ 用户匹配     │   └──────┬───────┘
              └──────┬──────┘          ↓
                     │         ┌──────────────┐
                     │         │  入库 + 评分  │
                     │         └──────┬───────┘
                     │                │
                     └────────┬───────┘
                              ↓
                    ┌───────────────────┐
                    │   解析质量评分     │
                    │  (静态 + 动态)    │
                    └────────┬──────────┘
                             ↓
                    ┌───────────────────┐
                    │   返回给用户       │
                    │   + 记录行为日志   │
                    └────────┬──────────┘
                             ↓
                    ┌───────────────────┐
                    │   数据闭环         │
                    │  👍/👎 → 重新评分 │
                    │  低分 → 重新生成  │
                    │  报错 → 人工审核  │
                    └───────────────────┘
```

### 4.2 关键设计原则

**1. 版本不可变（Immutable）**
- 每次修改生成**新版本**，旧版本保留不删除
- 方便追溯、回滚、AB 测试对比
- `version` 字段递增

**2. 冷热分离**
- 热门题目（日搜索 > 100）：多版本共存，AB 测试持续优化
- 冷门题目（日搜索 < 10）：只保留最高分版本，节省存储

**3. 用户行为权重 > 静态评分**
- 入库时的静态评分只是初筛
- 跑了两周后的动态效果分（👍/👎/完读率/报错率）才是最终排名依据
- 一个"LLM 生成+无人审核"但 👎 率 2% 的版本 > "人工审核"但 👎 率 8% 的版本

**4. LLM 是工具，不是裁判**
- LLM 负责**生成**解析和**辅助判断**质量（LLM-as-Judge）
- **最终排名**用真实用户行为数据，不用 LLM 排名（避免 LLM 偏好偏差）

---

## 五、总结清单

### LLM 选型速查

| 场景 | 首选 | 理由 |
|------|------|------|
| 拍照搜题-数学 | DeepSeek-V3 | 数学强+成本低+延迟低 |
| 拍照搜题-语文/英语 | Qwen3-Max / Claude Sonnet | 中文/文字质量好 |
| 批量生成解析 | DeepSeek-R1 + 多路投票 | 质量第一，延迟不敏感 |
| 解析质量评分 | GPT-4o-mini / GLM-4 | 简单打分，最便宜的够用 |
| 涉及学生数据 | 国内模型（DeepSeek/Qwen） | 合规要求 |
| 长文本（物理大题） | Kimi K2 | 超长上下文 |
| 多模态（图里含图表） | GPT-4V / Qwen-VL | 需要视觉理解 |

### 版本匹配速查

| 场景 | 策略 |
|------|------|
| 该题首次被搜 | LLM 生成 → 自动验收 → 标记"AI生成" |
| 该题已有审核版 | 直接返回审核版 |
| 该题有多版本 | 按综合分排序，默认展示最高分 |
| 用户是个性化场景 | 基础生→标准版，尖子生→进阶版 |
| 多版本都差不多（70-85） | LLM 合并各版本优点 |
| 用户大量 👎 | 自动触发重生成 |
