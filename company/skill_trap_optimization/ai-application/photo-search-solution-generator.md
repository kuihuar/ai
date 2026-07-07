# 精美解析生成方案

## 问题场景

题库有两个来源：
- **自有题库**：自己录入的题目，可能只有题干和答案，缺解析
- **抓取题库**：从公开源抓的题目，解析质量参差不齐（有的缺步骤、有的格式乱、有的是错的）

目标：将这些"裸题"变成**步骤完整、格式统一、有教学价值**的精美解析。

---

## 一、精美解析的标准定义

一个"精美解析"应该包含：

```
┌─────────────────────────────────┐
│  📝 题目                         │
│  已知 f(x) = x² + 2x + 1，求     │
│  f(x) 在区间 [-3, 0] 上的最小值。 │
├─────────────────────────────────┤
│  🔍 考点分析                     │
│  二次函数最值、配方法、区间讨论    │
├─────────────────────────────────┤
│  📐 解题步骤                     │
│  第 1 步：配方                   │
│    f(x) = (x+1)²                 │
│  第 2 步：分析对称轴              │
│    对称轴 x = -1，在区间内        │
│  第 3 步：计算最小值              │
│    f(-1) = 0                     │
│  第 4 步：验证                   │
│    端点值 f(-3)=4, f(0)=1       │
│    最小值为 0 ✓                  │
├─────────────────────────────────┤
│  ✨ 答案：最小值为 0              │
├─────────────────────────────────┤
│  💡 易错点提醒                    │
│  忽略区间端点 vs 顶点的大小比较    │
├─────────────────────────────────┤
│  📚 知识点链接                    │
│  二次函数图像、配方法、函数最值    │
└─────────────────────────────────┘
```

**关键维度**：

| 维度 | 标准 |
|------|------|
| 步骤完整性 | 每一步都有推导和依据 |
| 格式统一 | LaTeX 公式、步骤编号、层级一致 |
| 教学价值 | 有考点标注、易错提醒、方法总结 |
| 可读性 | 步骤不可太长，关键转换处有说明 |
| 正确性 | 答案正确，推导逻辑无漏洞 |

---

## 二、整体流程

```
原始题目（自有录入 / 抓取）
        ↓
┌──────────────────────┐
│ 阶段 1：题目清洗      │  格式化、去噪、归一化
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ 阶段 2：难度/类型分类 │  判断用什么方式生成解析
└──────┬───────────────┘
       ↓
       ├── 计算型（数学计算、化学方程式）
       │       → 符号计算引擎自动推导 + LLM 润色
       │
       ├── 推理型（证明题、物理推导、语文阅读）
       │       → 多轮 LLM 生成 + 自动验证
       │
       ├── 知识型（历史、生物、英语语法）
       │       → RAG 检索教科书 + LLM 总结
       │
       └── 客观题（选择题、填空题）
               → 先确定答案 → 逆推解析步骤
       ↓
┌──────────────────────┐
│ 阶段 3：解析生成      │  核心生成逻辑
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ 阶段 4：自动验收      │  正确性验证 + 质量评分
└──────┬───────────────┘
       ↓
    ┌── 通过 ──→ 入库发布
    │
    └── 不通过 ──→ 阶段 5：人工精修
                          ↓
                      入库发布
```

---

## 三、阶段详解

### 3.1 阶段 1：题目清洗

**自有题库的清洗**：
- 题目文本格式统一（全角/半角标点、多余空白）
- LaTeX 公式规范化（`x^2` → `x^{2}`、`\frac12` → `\frac{1}{2}`）
- 题干完整性检查（是否缺图、缺条件）

**抓取题库的清洗**（难点在这里）：

| 常见问题 | 处理方式 |
|----------|----------|
| 题目和答案混在一起 | 用 LLM 分离题干/答案/解析 |
| 包含网页噪声（HTML 标签、广告文字） | 正则清洗 + LLM 提取正文 |
| 图片 URL 已失效 | 检测死链 + 尝试文字描述替代 |
| 多道题挤在一段 | LLM 按题号切分 |
| 公式是图片而非 LaTeX | OCR 识别图片公式 → LaTeX（Pix2Text） |
| 不同来源格式不统一 | 统一转为 LaTeX + Markdown |

```python
def clean_scraped_question(raw_text: str) -> dict:
    """将抓取的原始文本清洗为结构化题目"""
    prompt = f"""从以下原始文本中提取结构化题目信息，以 JSON 返回。

原始文本：
{raw_text}

返回格式：
{{
  "subject": "学科",
  "question": "只提取题目正文，去掉无关内容",
  "answer": "最终答案",
  "solution": "如果有解析就提取，没有就留空",
  "difficulty": "基础/中等/提高",
  "knowledge_points": ["知识点1", "知识点2"]
}}
"""
    result = llm.generate(prompt, temperature=0)
    return parse_json(result)
```

### 3.2 阶段 2：难度/类型分类

```python
def classify_question(question: dict) -> str:
    """判断题目类型，决定用哪种生成方式"""
    subject = question["subject"]
    text = question["question"]

    # 数学/物理——大部分可用符号计算
    if subject in ["数学", "物理"]:
        if is_computational(text):     # 含方程、函数、数值计算
            return "symbolic"           # → 符号计算引擎
        elif is_proof(text):            # 含"证明"、"求证"
            return "llm_reasoning"      # → LLM 推理 + 多路验证
        else:
            return "llm_hybrid"         # → LLM 生成 + 符号验证

    # 化学方程式
    if subject == "化学":
        if has_equation(text):
            return "symbolic"           # 配平、计算用引擎
        return "llm_knowledge"          # 概念解释用 LLM

    # 英语——语法解析模板 + LLM
    if subject == "英语":
        return "template + llm"         # 语法点有模板，翻译/作文用 LLM

    # 语文、历史——RAG + LLM
    return "rag + llm"

def is_computational(text: str) -> bool:
    """判断是否为计算型题目"""
    patterns = [
        r"求解?", r"计算", r"求.*的值", r"解方程",
        r"化简", r"求导", r"积分",
        r"=\s*\?",                  # f(x)=?
    ]
    return any(re.search(p, text) for p in patterns)
```

### 3.3 阶段 3：解析生成（核心）

#### 方案 A：符号计算引擎（计算型题目）

**适用**：数学计算、解方程、求导积分、化学方程式配平

```
题目 → 结构化查询 → SymPy / Wolfram → 推导步骤 → LLM 润色成教学语言
```

```python
import sympy as sp

def generate_solution_symbolic(question: str, answer: str) -> dict:
    # 1. 用 LLM 提取数学表达式
    expr_data = llm_extract_math(question)
    # {"type": "solve_quadratic", "equation": "x**2 - 5*x + 6 = 0"}

    # 2. SymPy 给出分步推导
    x = sp.Symbol('x')
    eq = sp.Eq(x**2 - 5*x + 6, 0)

    steps = [
        {"step": 1, "title": "因式分解",
         "content": f"${sp.latex(sp.factor(x**2 - 5*x + 6))} = 0$",
         "reasoning": "使用十字相乘法，找两个数积为 6 和为 -5"},
        {"step": 2, "title": "写出两个一次方程",
         "content": f"${sp.latex(sp.Eq(x-2, 0))}$ 或 ${sp.latex(sp.Eq(x-3, 0))}$",
         "reasoning": "积为零则至少一个因子为零"},
        {"step": 3, "title": "求解",
         "content": "$x = 2$ 或 $x = 3$",
         "reasoning": "移项得到解"},
    ]

    # 3. LLM 润色：把 SymPy 的机械步骤转为教学语言
    polished = llm.polish(steps, style="高中教师",
                          addon=["易错提醒", "方法总结"])
    return polished
```

**能覆盖的题型**：

| 题型 | 计算引擎能力 | LLM 辅助做什么 |
|------|-------------|---------------|
| 一元/多元方程 | 100% 可解 | 步骤解释、多种解法对比 |
| 因式分解 | 100% | 解释分解思路 |
| 求导/积分 | 95% | 解释每一步的规则 |
| 代数化简 | 100% | 讲解化简技巧 |
| 化学配平 | 100% | 解释配平思路 |
| 几何计算 | 50%（解析几何可，纯几何不行） | 添加辅助线思路 |
| 概率统计 | 30% | LLM 主导，引擎计算概率值 |

#### 方案 B：LLM 多轮生成（推理型题目）

**适用**：证明题、物理推导、几何题、阅读理解

**核心思路**：LLM 做推理，但通过多轮生成 + 交叉验证保证质量。

```
题目 → Prompt 生成
          ↓
    3 路并行生成：
    路线 A：GPT-4 + 详细 CoT Prompt
    路线 B：Claude + 详细 CoT Prompt
    路线 C：GPT-4 + 不同的 Few-shot 示例
          ↓
    最终答案比对：
    ┌── 3 条一致 → 合并 + 选最好的解析
    ├── 2 条一致 → 多数答案 + 人工抽检
    └── 都不一致 → 全部转人工
```

```python
def generate_solution_multi_path(question: dict) -> dict:
    """多路生成 + 投票 + 选最优"""
    # 1. 三路并行生成
    paths = parallel_generate([
        ("gpt4", build_cot_prompt(question)),
        ("claude", build_cot_prompt(question)),
        ("gpt4", build_fewshot_prompt(question)),
    ])

    # 2. 最终答案比对
    final_answers = [extract_final_answer(p) for p in paths]

    # 3. 投票
    answer_counts = Counter(final_answers)
    winner, count = answer_counts.most_common(1)[0]

    if count >= 2:
        # 多数一致，选质量最好的那份解析
        best = score_and_pick_best([p for p in paths
                if extract_final_answer(p) == winner])
        return {"status": "auto", "solution": best}
    else:
        # 全不一致，转人工
        return {"status": "manual_review", "candidates": paths}
```

**分科目 Prompt 设计**（以物理为例）：

```markdown
## 物理题解析生成 Prompt

你是高中物理教师，请为下面这道题生成一份精美解析。

题目：{question}

## 解析结构要求

### 1. 审题分析（2-3 句）
- 识别物理情景（什么运动/什么过程）
- 明确已知量和待求量
- 判断涉及的物理规律

### 2. 解题步骤（逐步推导）

每一步必须包含：
- 使用的物理公式（LaTeX 格式）
- 代入数值的过程
- 公式的物理意义解释
- 中间结果的单位

示例格式：
> **第 1 步：受力分析**
> 物体受重力 G=mg，支持力 N。
> 沿斜面方向：$mg\sin\theta - f = ma$
> 代入：$2\times10\times0.5 - f = 2\times a$
>
> **第 2 步：...**

### 3. 答案
明确标出最终结果，带单位

### 4. 易错提醒
- 最常见的解题错误（1-2 点）
- 容易混淆的概念

### 5. 方法总结
- 这类题的通法是什么
- 有没有更简便的解法（一题多解）

## 输出要求
- 所有公式用 $...$（行内）或 $$...$$（独立行）包裹
- 步骤编号：第 1 步、第 2 步...
- 关键转换处加简短解释
```

#### 方案 C：RAG + LLM（知识型题目）

**适用**：历史、生物、英语语法等需要知识储备的题

```
题目 → 识别知识点 → 检索教材/教辅/百科
                         ↓
                   LLM 基于检索内容生成解析
                         ↓
                   标注引用来源
```

```python
def generate_solution_rag(question: dict) -> dict:
    # 1. 识别涉及知识点
    knowledge_points = llm.extract_knowledge_points(question["text"])

    # 2. 从教材库/教辅库/百科检索相关内容
    references = []
    for kp in knowledge_points:
        refs = knowledge_base.search(kp, top_k=3)
        references.extend(refs)

    # 3. LLM 基于权威文献生成解析
    prompt = f"""基于以下参考资料，为题目生成解析：

参考资料：
{format_references(references)}

题目：{question["text"]}

要求：
1. 解析必须基于参考资料的权威说法，不要臆造
2. 在引用处标注来源
3. 按"考点→解析→答案→知识点扩展"组织
"""
    return llm.generate(prompt)
```

#### 方案 D：模板 + 变量填充（可模板化的题型）

**适用**：英语语法选择题、化学方程式配平、数学基础运算

```python
TEMPLATES = {
    "english_grammar": """## 考点分析
本题考察{grammar_point}。

## 解题思路
{grammar_point}的规则是：{rule}。

- 选项 A: {option_a} — {analysis_a}
- 选项 B: {option_b} — {analysis_b}
- 选项 C: {option_c} — {analysis_c}
- 选项 D: {option_d} — {analysis_d}

## 答案：{answer}

## 易错提醒
{common_mistake}
""",
    "chem_balance": """## 配平过程
化学方程式：{raw_equation}

**第 1 步：标注化合价**
{oxidation_states}

**第 2 步：列电子守恒方程**
{electron_balance}

**第 3 步：确定系数**
{coefficients}

**配平结果**
$${balanced_equation}$$

## 反应类型
{reaction_type}
""",
}
```

---

### 3.4 阶段 4：自动验收

**验收不是一次性的，是分层级的**：

```
生成解析
    ↓
第一层：格式检查（100% 自动化）
    → LaTeX 语法是否合法
    → 步骤结构是否完整（考点 + 步骤 + 答案 + 易错）
    → 是否有明显截断
    ↓ 通过
第二层：答案正确性检查（视题型自动化）  
    → 方程解代入验证（SymPy）：自动
    → 选择题选项匹配：自动
    → 证明题：无法自动，抽样人工
    → 阅读题：无法自动，抽样人工
    ↓ 通过
第三层：质量评分（LLM-as-Judge）
    → 用另一个 LLM 评分
    → 不合格（< 70 分）→ 人工复核
    ↓ 通过
第四层：人工抽检
    → 每天随机抽 5% 检查
    → 发现问题回滚该批次的自动生成
```

```python
def auto_verify(question: dict, solution: dict) -> VerifyResult:
    """自动验收生成的解析"""

    # 第一层：格式
    if not latex_valid(solution["content"]):
        return VerifyResult.FAIL("LaTeX 语法错误")
    if not structure_complete(solution):
        return VerifyResult.FAIL("缺少必要结构（步骤/答案/考点）")

    # 第二层：答案验证
    if question["type"] == "equation":
        if not verify_equation_solution(question["text"], solution["answer"]):
            return VerifyResult.FAIL("答案代入方程不成立")
    if question["type"] == "choice":
        if solution["answer"] not in question["options"]:
            return VerifyResult.FAIL("答案不在选项中")

    # 第三层：LLM 质量评分
    score = llm_score_solution(question, solution)
    if score < 70:
        return VerifyResult.MANUAL_REVIEW(f"质量评分过低: {score}")

    return VerifyResult.PASS()

def llm_score_solution(question: dict, solution: dict) -> int:
    """用另一个 LLM 给解析打分"""
    prompt = f"""请对以下题目的解析质量打分（0-100）：

题目：{question["text"]}
答案：{question["answer"]}
解析：{solution["content"]}

评分维度（各 25 分）：
1. 正确性：答案正确，推导无误（25分）
2. 完整性：步骤完整，无跳跃（25分）
3. 教学性：有考点分析、有方法总结（25分）
4. 可读性：格式规范，公式正确（25分）

只返回分数，不要其他内容。"""
    return int(llm.generate(prompt, temperature=0))
```

---

### 3.5 人工精修工单

自动验收不通过或低分解析，进入人工精修流程。

```
审核工作台：
┌────────────────────────────────────────┐
│  题目：求函数 f(x)=x²+2x+1 在 [-3,0]   │
│        上的最小值。                     │
│                                        │
│  标准答案：最小值为 0                   │
│                                        │
│  ┌── AI 生成的解析 ──────────────────┐  │
│  │ 第 1 步：求导                       │  │
│  │ f'(x) = 2x + 2                      │  │
│  │ 令 f'(x) = 0，得 x = -1             │  │
│  │                                     │  │
│  │ 第 2 步：代回原函数                  │  │
│  │ f(-1) = 0    ← 缺少端点值比较       │  │
│  │                                     │  │
│  │ ⚠️ AI 质量评分：72，扣分原因：       │  │
│  │   - 未比较端点值 [-3, 0]            │  │
│  └────────────────────────────────────┘  │
│                                        │
│  ✏️ 人工修正：                          │
│  ┌────────────────────────────────────┐  │
│  │ 在"第 2 步"后面加上：               │  │
│  │ 检查端点：f(-3)=4, f(0)=1          │  │
│  │ ∵ f(-1)=0 < min(f(-3), f(0))       │  │
│  │ ∴ 最小值为 0                        │  │
│  └────────────────────────────────────┘  │
│                                        │
│  [通过] [驳回重生成] [修改后通过]       │
└────────────────────────────────────────┘
```

**人工精修 UI 关键特性**：
- 左侧：题目 + 标准答案
- 右侧：AI 生成的解析 + AI 评分 + 扣分原因
- 一键定位 AI 标注的问题点
- 修正后可立即预览
- 通过后的解析自动入库，作为将来 Few-shot 的示例

**效率**：一个人一天可审核 200-300 题（大部分只需确认，少部分修改）。

---

## 四、批量生成策略

### 4.1 优先级排序

不同题目的生成优先级不同：

| 优先级 | 条件 | 处理方式 |
|--------|------|----------|
| P0（最高） | 用户搜索但无解析的题目 | 准实时触发生成 |
| P1 | 用户高频搜索的题目 | 每日批量补全 |
| P2 | 有解析但用户 👎 多的 | 优先安排重生成 |
| P3 | 新入库的"裸题" | 空闲时段批量处理 |

### 4.2 批量生成流水线

```
消息队列（题目 ID 列表）
     ↓
消费者池（N 个 Worker）
     ↓
1. 取题目 → 查是否已有解析（有且👎少就跳过）
2. 题目清洗 → 类型分类
3. 生成解析（走对应的方案 A/B/C/D）
4. 自动验收
5. 通过 → 入库 + 更新状态
   不通过 → 写入"待人工精修"队列
```

---

## 五、抓取题库的特殊处理

抓取的题目有些**已经有解析**，但质量参差不齐。不是所有都需要重新生成。

### 抓取解析的质量分级

```python
def classify_scraped_solution(raw_solution: str) -> str:
    """对抓取到的解析做质量分级"""
    score = rate_solution_quality(raw_solution)

    if score >= 85:
        return "good"         # 格式清洗后直接入库
    elif score >= 60:
        return "fixable"      # LLM 润色 + 补全缺口，人工抽检
    else:
        return "regenerate"   # 丢弃原始解析，重新生成
```

**评分维度**：

| 检查项 | 权重 | 检测方式 |
|--------|------|----------|
| 是否有完整步骤推导 | 30% | LLM 判断 + 步骤数量 |
| 答案是否正确 | 25% | 代入验证 / 与标准答案比对 |
| 格式是否规范 | 20% | LaTeX 检测、结构检测 |
| 是否有教学价值（考点+易错） | 15% | LLM 判断 |
| 是否有垃圾内容（广告等） | 10% | 规则 + LLM 检测 |

### "fixable" 类型的 LLM 润色

```python
def polish_scraped_solution(question: str, raw_solution: str) -> str:
    """将抓取的低质量解析润色成标准格式"""
    prompt = f"""请将下面的原始解析润色成标准格式，保留正确的部分，修正错误，补充缺失。

题目：{question}

原始解析（可能有格式问题、错误或缺失）：
{raw_solution}

标准格式：
### 考点分析
[一句话点出考察的知识点]

### 解题步骤
第 1 步：[步骤名]
[推导过程 + 公式 + 理由]
...

### 答案
[明确标出最终答案]

### 易错提醒
[指出常见错误]

### 方法总结
[总结这类题的通法]

重要：
- 如果原始解析中有错误，请修正它
- 如果步骤有跳跃，请补充中间步骤
- 保持原始解析中正确的部分不变
- 所有公式使用 LaTeX 格式
"""
    return llm.generate(prompt)
```

---

## 六、不同科目的差异化策略

| 科目 | 生成策略 | 验证方式 | 特殊要求 |
|------|----------|----------|----------|
| **数学** | 计算型→符号引擎；证明型→LLM多路 | 代入验证 + SymPy | 多解展示（代数/几何） |
| **物理** | LLM 为主 + 公式计算验证 | 量纲检查 + 数值验证 | 受力分析图、过程示意图 |
| **化学** | 方程式配平→符号引擎；概念→RAG+LLM | 原子守恒验证 | 结构式、反应流程图 |
| **英语** | 语法→模板填充；翻译/阅读→LLM | 选项匹配 | 语法规则引用 |
| **语文** | RAG + LLM（教材/权威参考） | 人工抽检为主 | 文言文注释、背诵提示 |
| **生物** | RAG + LLM | 人工抽检 | 过程图解（如有丝分裂） |

---

## 七、质量衡量指标

| 指标 | 定义 | 目标 |
|------|------|------|
| 自动化率 | 无需人工修改直接入库的比例 | > 80% |
| 验收通过率 | 自动验收通过的比例 | > 85% |
| 人工修正率 | 需要人工修改才入库的比例 | < 15% |
| 用户好评率 | 答案页面 👎 的比例 | < 5% |
| 解析完备率 | 有解析的题 / 总题目数 | > 95% |
| 入库速度 | 每天可处理多少道"裸题" | > 5000 题/天 |

---

## 八、总结

```
自有“裸题” ──→ 题型分类 ──→ 选生成方案 ──→ 自动验收 ──→ 人工抽检 ──→ 入库
                                       │
抓取题目 ──→ 质量分级 ──→ good: 清洗入库
                         fixable: LLM 润色 → 验收
                         bad: 丢弃 → 重新生成
```

核心原则：
- **计算归引擎，推理归 LLM，知识归检索**——不同题型不该用同一套方案
- **多路验证比单路调优更可靠**——两条路线的错误不重叠
- **人工审核做裁判，不做苦力**——LLM 生成 + 自动验 + 人工抽检
- **抓取的解析要分级处理**——好的直接洗、中的 LLM 润色、差的丢掉重做
