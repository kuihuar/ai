# 题库维护方案：从 0 到 1 全面展开

---

## 一、总体蓝图

题库维护不只是"存题目"，而是一个从**来源 → 清洗 → 入库 → 质量保障 → 检索 → 生命周期管理 → 下线**的完整系统工程。

```
                        ┌─────────────────────┐
                        │    题 目 来 源       │
                        │ 录入/抓取/交换/AI生成 │
                        └──────────┬──────────┘
                                   ↓
                        ┌─────────────────────┐
                        │    清 洗 管 道       │
                        │ 格式化/去重/分类/标签 │
                        └──────────┬──────────┘
                                   ↓
                        ┌─────────────────────┐
                        │    存 储 层          │
                        │ PG / ES / 向量库 / OSS│
                        └──────────┬──────────┘
                                   ↓
                        ┌─────────────────────┐
                        │    质 量 保 障       │
                        │ 自动验收/人工审核/反馈 │
                        └──────────┬──────────┘
                                   ↓
                        ┌─────────────────────┐
                        │    检 索 服 务       │
                        │ 精确/BM25/语义/Rerank │
                        └──────────┬──────────┘
                                   ↓
                        ┌─────────────────────┐
                        │    生 命 周 期       │
                        │ 发布/更新/降级/归档   │
                        └─────────────────────┘
```

---

## 二、题目来源：如何把题目"搞进来"

### 2.1 来源全景

| 来源 | 数量 | 质量 | 成本 | 速度 | 适合阶段 |
|------|------|------|------|------|----------|
| **人工录入** | 低（一天几十道） | 最高 | 高（人力） | 慢 | 初期打底 |
| **公开题库抓取** | 高（一天几万道） | 参差不齐 | 低 | 快 | 快速扩量 |
| **厂商合作/交换** | 中（批量） | 中-高 | 中（分成/买断） | 中 | 中期 |
| **版权采购** | 中 | 高（出版社审核过） | 高（按题/按书付费） | 中 | 建立壁垒 |
| **AI 生成（裸题）** | 极高 | 需验证 | 低 | 极快 | 补充偏难怪 |
| **用户上传/纠错** | 低但精准 | 中 | 极低 | 慢 | 数据闭环 |
| **教材/教辅 OCR** | 中（一本教辅几百道） | 中（需 OCR 校对） | 中 | 中 | 特定版本覆盖 |
| **竞品题目交叉比对** | - | - | - | - | 补缺分析（不直接抓） |

### 2.2 人工录入——种子题库

**做法**：
- 初期投入 2-3 名教师/教研员，每个学科录 5000-10000 道核心题
- 录入标准：选题覆盖所有考点，确保每个知识点都有题目
- 工具：Web 后台 + 富文本编辑器 + LaTeX 插件 + 图片上传

**录入后台核心功能**：

```
┌──────────────────────────────────────────┐
│  题目录入                                 │
├──────────────────────────────────────────┤
│  学科：[下拉] 数学                        │
│  题型：[下拉] 解答题                      │
│  难度：[下拉] 中等             年级：[下拉] 高一 │
│                                          │
│  知识点标签：[多选] 二次函数  配方法  最值  │
│  教材版本：[多选] 人教A版  北师大版        │
│                                          │
│  题目：                                   │
│  ┌────────────────────────────────────┐  │
│  │ 已知函数 $f(x)=x^2+2x+1$，          │  │
│  │ 求 $f(x)$ 在 $[-3,0]$ 上的最小值。   │  │
│  │                            [LaTeX] │  │
│  └────────────────────────────────────┘  │
│                                          │
│  (+ 添加图片)                            │
│                                          │
│  答案：$-1$ 时取得最小值 $0$              │
│                                          │
│  解析：[高级编辑器]                       │
│  ┌────────────────────────────────────┐  │
│  │ ### 考点分析                        │  │
│  │ 二次函数在闭区间上的最值。            │  │
│  │ ### 解题步骤                        │  │
│  │ ...                                │  │
│  └────────────────────────────────────┘  │
│                                          │
│  [保存草稿] [提交审核]                    │
└──────────────────────────────────────────┘
```

**效率目标**：熟练后 5-10 分钟/题（含解析），每人每天 50-100 题。

### 2.3 公开题库抓取——快速扩量

**抓取目标**：各类教育网站、作业问答平台、公开试卷库。

```python
# 抓取管道架构
class ScrapingPipeline:
    def __init__(self):
        self.sources = [
            ZhihuEducation(),    # 知乎教育板块
            ZuoyebangPublic(),   # 作业帮公开页面
            ExamPaperSites(),    # 试卷下载站
            ForumQA(),           # 问答论坛
            BaiduZhidao(),       # 百度知道
            WenkuFree(),         # 文库免费部分
        ]
        self.normalizer = QuestionNormalizer()
        self.deduper = QuestionDeduplicator()
        self.classifier = QuestionClassifier()

    def run(self, source_url):
        raw = self.fetch(source_url)           # 1. 抓取
        clean = self.normalizer.clean(raw)     # 2. 清洗
        if self.deduper.is_duplicate(clean):   # 3. 去重
            return
        structured = self.classifier.classify(clean)  # 4. 分类
        self.save(structured)                  # 5. 入库
```

**抓取频率控制**：
- IP 轮换代理池，避免被封
- 每个源间隔 1-5 秒，模拟人类行为
- 遵守 robots.txt
- 增量抓取：只抓新增/更新的页面，不重复爬

**合规风险**：
- 抓取公开数据用于内部题库建设，存在版权风险
- 不要直接展示抓取的原文（版权方的文字），而是用自己的格式重新组织
- 标注"用户贡献"或"整理自公开资料"作为来源标注
- 对于明确禁止转载的平台，不抓

### 2.4 厂商合作/交换

- 与教育出版社、教辅机构合作：用流量换题目版权，或按调用次数分成
- 题目交换：与互补产品（如：你做数学，别人做英语）互换题库
- 标准：统一格式（双方约定 JSON/Excel schema），导入前用同样的清洗管道处理

### 2.5 AI 生成裸题——补充长尾

对于搜索量大但题库没有的题目，用 LLM + 知识点体系自动"造题"：

```
知识点树 → 选考点 → 指定难度 → LLM 生成题目 + 答案 → 自动验收 → 人工抽检
```

```python
def generate_question(knowledge_point: str, difficulty: str) -> dict:
    prompt = f"""请生成一道 {difficulty} 难度的 {knowledge_point} 题目。

要求：
1. 题目有明确的已知条件和求解目标
2. 不能与常见教材题雷同（数字、情景要有变化）
3. 必须确保有唯一正确答案
4. 格式：题目 + 答案 + 完整解析

返回 JSON。
"""
    return llm.generate(prompt)
```

**关键控制**：
- 每道 AI 生成的题标记 `source=ai_generated`，不直接推送给用户作为"推荐解析"
- 需经过自动验收 + 人工抽检后，改为 `source=ai_verified` 才能作为正式内容
- 限制每天 AI 生成量（如每天 5000 道），防止低质量内容泛滥

---

## 三、题目清洗：把"脏数据"洗干净

### 3.1 清洗管线

```
原始题目 → 格式标准化 → 题干/答案/解析分离 → 公式规范化 → 去HTML/广告 → 完整性检查
```

### 3.2 格式标准化

```python
def normalize_question(raw_text: str) -> str:
    """统一格式"""
    text = raw_text

    # 全角转半角（英文、数字、标点）
    text = fullwidth_to_halfwidth(text)

    # LaTeX 分隔符统一：$$ → $，\[ → $
    text = re.sub(r'\$\$(.+?)\$\$', r'$\1$', text)
    text = re.sub(r'\\\[(.+?)\\\]', r'$\1$', text)

    # 多余空白、换行清理
    text = re.sub(r'\n{3,}', '\n\n', text)
    text = re.sub(r' {2,}', ' ', text)

    # 题号统一：1. / (1) / 一、 / 题1： → 1.
    text = re.sub(r'^[\(（]\d+[\)）]', lambda m: m.group()[1:-1]+'.', text)

    return text.strip()
```

### 3.3 题干/答案/解析分离

抓取的题目很多是"题干答案解析混在一起"的。用 LLM 做分离是最可靠的方式。

```python
def separate_parts(raw: str) -> dict:
    prompt = f"""请将以下文本拆分为：题目、答案、解析 三部分。

如果原文没有解析，solution 字段留空。
如果文本包含多道题，只返回主标题对应的那一道。

原始文本：
{raw}

返回 JSON：
{{"question": "...", "answer": "...", "solution": "..."}}
"""
    return llm.generate(prompt, temperature=0)
```

### 3.4 图片处理

```python
class ImageProcessor:
    def process(self, img_url: str) -> str:
        # 1. 下载图片
        img = download(img_url)

        # 2. 校验
        if not is_valid_image(img):
            return None  # 死链、非图片文件

        # 3. 压缩 + 转格式（WebP 比 JPEG 小 30%）
        img = compress(img, max_width=1200, quality=85)
        img = convert_to_webp(img)

        # 4. 上传 OSS
        oss_url = upload_to_oss(img, path=f"questions/{date}/{uuid}.webp")

        # 5. 如果图里有公式文字，OCR 提取
        if has_formula(img):
            latex = ocr_formula(img)
            return {"url": oss_url, "latex": latex}

        return {"url": oss_url}
```

### 3.5 题目去重（多维度）

**五层去重策略**：

```
第一层：内容哈希（SHA256）
    → 完全相同文本去重（100% 准确）

第二层：题目指纹（结构化特征 MD5）
    → 同题不同标点/排版的去重

第三层：SimHash（局部敏感哈希）
    → 高度相似的题目（相似度 > 95%）

第四层：Embedding 聚类
    → 同题不同表述的聚类（"求解 x²+1=0" vs "方程 x²+1=0 的解"）

第五层：人工合并
    → Embedding 聚类后的簇，人工确认是否为同一道题
```

```python
class QuestionDeduplicator:
    def __init__(self):
        self.fingerprinter = QuestionFingerprinter()
        self.simhasher = SimHash()

    def is_duplicate(self, question: dict) -> tuple[bool, str]:
        # 第一层：精确哈希
        content_hash = sha256(question["text"])
        if db.exists("content_hash", content_hash):
            return True, "exact_match"

        # 第二层：题目指纹
        fp = self.fingerprinter.hash(question)
        existing = db.query("SELECT id FROM questions WHERE fingerprint = ?", fp)
        if existing:
            return True, f"fingerprint_match: {existing.id}"

        # 第三层：SimHash 海明距离
        simhash = self.simhasher.hash(question["text"])
        candidates = db.query("""
            SELECT id, text, simhash
            FROM questions
            WHERE subject = ? AND hamming_distance(simhash, ?) < 3
            LIMIT 20
        """, question["subject"], simhash)

        for c in candidates:
            if text_similarity(question["text"], c.text) > 0.95:
                return True, f"near_dup: {c.id}"

        return False, None

    def cluster_duplicates(self):
        """离线：用 Embedding 发现隐藏的重复题"""
        # 1. 取最近 N 天的新题 Embedding
        embeddings = db.query("SELECT id, embedding FROM questions WHERE embedding IS NOT NULL")

        # 2. DBSCAN 聚类
        clusters = DBSCAN(eps=0.05, min_samples=2).fit(embeddings)

        # 3. 每个簇内的题 → 人工审核合并
        for cluster_id, ids in clusters:
            db.insert("duplicate_clusters", {
                "cluster_id": cluster_id,
                "question_ids": ids,
                "status": "pending_review",
            })

        # 4. 审核后台展示
        # "以下 5 道题系统判定为同一道题，请确认是否合并："
```

---

## 四、题目分类与标签体系

### 4.1 知识点树设计

知识点树不能太深（不好用），也不能太浅（太粗粒度）。推荐 **3 级结构**：

```
数学
├── 代数
│   ├── 集合与常用逻辑用语      # ← 知识点（Leaf，挂题目）
│   ├── 一元二次函数、方程和不等式
│   ├── 函数概念与性质
│   └── 指数函数与对数函数
├── 几何
│   ├── 平面向量
│   ├── 立体几何
│   └── 解析几何
├── 概率与统计
│   └── ...
└── 微积分
    └── ...
```

```sql
CREATE TABLE knowledge_tree (
    id          BIGINT PRIMARY KEY,
    parent_id   BIGINT,
    name        VARCHAR(128),
    level       SMALLINT,       -- 1=学科 2=模块 3=知识点
    sort_order  INT,
    INDEX idx_parent (parent_id)
);

CREATE TABLE question_knowledge (
    question_id  BIGINT,
    knowledge_id BIGINT,
    PRIMARY KEY (question_id, knowledge_id)
);
```

**题目应该挂多少个知识点标签？**
- 典型：1-2 个知识点（综合题最多 3 个）
- 标签太多 = 检索精度下降
- 规则：只标题目主要考察的知识点

### 4.2 自动分类

```python
def auto_tag(question: dict) -> list:
    """自动给题目打知识点标签"""
    prompt = f"""这是我们的知识点树：
{knowledge_tree_json}

题目：{question['text']}

请从知识点树中选择 1-2 个最匹配的知识点 ID，返回格式：
{{"knowledge_ids": [3, 7], "confidence": [0.95, 0.82]}}

如果 confidence < 0.7，不要勉强标注。
"""
    result = llm.generate(prompt, temperature=0)
    return result["knowledge_ids"]
```

**人工审核环节**：
- 新入库题目先由 LLM 自动分类
- 每天随机抽 20% 人工核对分类是否正确
- LLM 置信度 < 0.7 的题目，全部转人工分类

### 4.3 教材版本映射

一道题可能适用于多个教材版本。需要维护教材-知识点映射表：

```sql
-- 教材版本
CREATE TABLE textbook_editions (
    id      BIGINT PRIMARY KEY,
    name    VARCHAR(64),      -- "人教A版"
    subject VARCHAR(32),
    grade   VARCHAR(16),      -- "高一上"
    year    INT,              -- 出版年份（2023版 vs 2019版）
);

-- 教材目录 → 知识点 映射
CREATE TABLE textbook_knowledge_map (
    textbook_id  BIGINT,
    chapter_id   VARCHAR(16),      -- "必修1-第2章-第3节"
    chapter_name VARCHAR(128),     -- "二次函数与一元二次方程"
    knowledge_id BIGINT,           -- 对应知识点树的 ID
    is_core      BOOLEAN,          -- 是否为核心考点
);
```

有了这张映射表，一道题目标注了知识点后，自动知道它适用于哪些教材版本的哪个章节。

---

## 五、存储方案

### 5.1 存储矩阵

| 数据类型 | 存储方案 | 索引方式 | 分库分表 |
|----------|----------|----------|----------|
| 题目结构化数据 | PostgreSQL | B+Tree / GIN (全文) | 按学科 hash 分表 |
| 题目文本全文搜索 | Elasticsearch | 倒排索引 | ES 天然分片 |
| 题目 Embedding 向量 | pgvector / Milvus | IVFFlat / HNSW | 按学科分区 |
| 题目图片 | 对象存储 (OSS/S3) | CDN 加速 | 按日期/学科分目录 |
| 解析版本历史 | PostgreSQL (大字段 JSONB) | question_id + version | 同题目表 |
| 用户行为日志 | ClickHouse / 时序库 | 时间分区 | 按天分区，保留 90 天 |
| 缓存 (热点题目) | Redis | key=question:{id} | 按 type hash slot |

### 5.2 PostgreSQL 表设计

```sql
-- 核心题目表
CREATE TABLE questions (
    id              BIGSERIAL PRIMARY KEY,
    fingerprint     VARCHAR(64) NOT NULL,      -- 题目指纹
    source          VARCHAR(32) NOT NULL,      -- manual/scraped/ai_generated/exchange
    source_url      TEXT,                      -- 抓取来源 URL（可选）
    subject         VARCHAR(32) NOT NULL,      -- math/physics/chemistry/english/chinese/biology
    type            VARCHAR(32) NOT NULL,      -- choice/fill/solution/proof
    difficulty      VARCHAR(16),               -- easy/medium/hard
    grade_range     INT4RANGE,                 -- 适用年级范围 [7,12)

    -- 内容
    content         TEXT NOT NULL,             -- 题目正文（纯文本 + LaTeX）
    content_hash    VARCHAR(64) NOT NULL,      -- SHA256 唯一索引
    options         JSONB,                     -- 选择题选项 ["A.xx","B.xx",...]
    images          JSONB,                     -- [{"url":"...","alt":"...","latex":"..."}]

    -- 检索
    embedding       VECTOR(1536),              -- pgvector，用于语义检索
    search_text     TSVECTOR,                  -- PostgreSQL 全文检索
    simhash         BIGINT,                    -- SimHash 用于近似去重

    -- 状态与统计
    status          VARCHAR(16) DEFAULT 'draft', -- draft/reviewing/published/deprecated
    review_version  BIGINT,                    -- 当前展示的解析版本 ID
    search_count    BIGINT DEFAULT 0,          -- 总搜索次数（用于热度排序）
    like_rate       DECIMAL(4,3) DEFAULT 0,    -- 综合好评率

    -- 时间
    published_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),

    -- 索引
    UNIQUE (content_hash),
    INDEX idx_subject_type (subject, type),
    INDEX idx_subject_difficulty (subject, difficulty),
    INDEX idx_search_count (search_count DESC),
    INDEX idx_embedding USING ivfflat (embedding vector_cosine_ops) WITH (lists = 200),
    INDEX idx_search_text USING GIN (search_text),
);

-- 解析版本表（前面定义过，这里补充）
CREATE TABLE solution_versions (
    id              BIGSERIAL PRIMARY KEY,
    question_id     BIGINT NOT NULL,
    version         INT NOT NULL,
    source          VARCHAR(32),         -- manual/llm_gpt4/llm_deepseek/scraped/user_fix
    method          VARCHAR(32),         -- standard/alternative_1/alternative_2

    analysis        TEXT,                -- 考点分析
    steps           JSONB NOT NULL,      -- 解题步骤
    answer          TEXT NOT NULL,       -- 最终答案
    tips            TEXT,                -- 易错提醒
    summary         TEXT,                -- 方法总结

    quality_score   DECIMAL(3,1),
    auto_verified   BOOLEAN DEFAULT FALSE,
    review_status   VARCHAR(16) DEFAULT 'pending',

    -- 效果
    impression_count BIGINT DEFAULT 0,
    like_count       BIGINT DEFAULT 0,
    dislike_count    BIGINT DEFAULT 0,
    report_count     BIGINT DEFAULT 0,
    avg_view_time_ms BIGINT DEFAULT 0,

    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE (question_id, version),
    INDEX idx_question_quality (question_id, quality_score DESC),
);
```

### 5.3 分库分表策略

**什么时候需要分？**
- 单表 > 5000 万行 → 查询变慢
- 题库存量通常在百万到千万级，单表够用

**真要分的话，按学科 hash**：

```sql
-- 8 张分表
questions_0, questions_1, ..., questions_7

-- 路由
table_suffix = hash(subject + id) % 8
```

### 5.4 Elasticsearch 索引设计

PG 负责精确查询和关联查询，ES 负责模糊搜索和聚合分析。

```json
{
  "mappings": {
    "properties": {
      "content": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "subject": { "type": "keyword" },
      "type": { "type": "keyword" },
      "difficulty": { "type": "keyword" },
      "grade": { "type": "integer" },
      "knowledge_ids": { "type": "keyword" },
      "search_count": { "type": "long" },
      "quality_score": { "type": "float" },
      "created_at": { "type": "date" }
    }
  }
}
```

---

## 六、质量保障体系

### 6.1 质量控制的完整链路

```
新题入库
    ↓
自动验收（格式/答案验证/相似题检测）
    ├── 不通过 → 标记待处理
    └── 通过 → 进入待审核池
                ↓
        人工审核（优先级排序：热门题 > 搜索多 > 新入库）
                ↓
        ┌── 通过 → 发布（status=published）
        ├── 修改后通过 → 修改后发布
        └── 驳回 → 标记删除/重做
                ↓
        线上持续监控（用户反馈 + LLM 抽检）
                ↓
        低分/报错 → 触发重审
```

### 6.2 人工审核工作台

```
待审列表（按优先级排序）：
┌──────┬──────────┬────────┬─────────────────────┬──────────┐
│ 优先级│  学科    │  来源   │  预览               │  操作    │
├──────┼──────────┼────────┼─────────────────────┼──────────┤
│ 高   │ 数学     │ 抓取   │ 已知 f(x)=x²+...    │ [审核]   │
│ 中   │ 英语     │ AI生成 │ Choose the corre... │ [审核]   │
│ 低   │ 历史     │ 录入   │ 鸦片战争发生于...    │ [审核]   │
└──────┴──────────┴────────┴─────────────────────┴──────────┘

审核详情页：
┌──────────────────────────────────────────────┐
│ 题目预览（渲染后的效果）                       │
│ ┌──────────────────────────────────────────┐ │
│ │ ### 已知 f(x)=x²+2x+1，求 f(x) 在 [-3,0] │ │
│ │ 的最小值                                  │ │
│ └──────────────────────────────────────────┘ │
│                                              │
│ 解析预览：                                    │
│ ┌──────────────────────────────────────────┐ │
│ │ ### 考点分析：二次函数闭区间最值           │ │
│ │ ### 解题步骤：                            │ │
│ │ 1. 配方法：f(x)=(x+1)²                    │ │
│ │ 2. 对称轴 x=-1 在区间 [-3,0] 内          │ │
│ │ 3. f(-1)=0 为最小值                      │ │
│ │ ...                                      │ │
│ └──────────────────────────────────────────┘ │
│                                              │
│ AI 标注的潜在问题：                           │
│ ⚠️ 未比较端点值 f(-3)=4, f(0)=1              │
│ ✅ 答案正确                                  │
│ ✅ 公式格式正确                              │
│                                              │
│ 审核操作：                                    │
│ [通过] [驳回] [修改] [标记为优质]              │
│                                              │
│ 修改编辑区：                                  │
│ ┌──────────────────────────────────────────┐ │
│ │ [可编辑的解析内容...]                      │ │
│ └──────────────────────────────────────────┘ │
└──────────────────────────────────────────────┘
```

### 6.3 质量控制指标

| 指标 | 定义 | 目标 | 监控频率 |
|------|------|------|----------|
| 入库一次性通过率 | 自动验收通过 / 总入库 | > 80% | 每日 |
| 人工审核驳回率 | 驳回 / 审核总量 | < 15% | 每日 |
| 待审积压量 | pending_review 状态的题目数 | < 500 | 实时 |
| 审核效率 | 每人每天审核题数 | > 200 | 每周 |
| 用户报错率 | 报错 / 展示次数 | < 2% | 每日 |
| 低分题比例 | quality < 70 的题 / 总量 | < 5% | 每周 |
| 热门题覆盖率 | Top 1000 搜索有解析的比例 | > 99% | 每日 |

---

## 七、题目生命周期管理（容易忽略）

### 7.1 状态流转

```
draft（草稿）
   ↓
reviewing（审核中）
   ↓
published（已发布）
   ↓
  ├── updated（有新版本）
  ├── deprecated（内容过时但仍可搜到）
  └── archived（归档：不再展示但保留数据）
```

### 7.2 何时 deprecate（降级）一道题？

| 触发条件 | 处理 |
|----------|------|
| 教材改版，这道题对应的考点被删了 | status → deprecated，搜索结果中标记"旧版大纲" |
| 用户报错率 > 10%，多次修正仍高 | deprecate，提示用户"此答案可能有争议" |
| 有更高质量的同题版本入库 | 旧版 deprecate，引导到新版 |
| 题目本身有歧义/错误 | deprecate + 备注原因 |

### 7.3 教材改版影响面分析（一个容易被忽略的大问题）

教材 3-5 年一改版。一次改版可能影响上万道题。

**应对方案**：

```sql
-- 题目与教材版本的关联
CREATE TABLE question_textbook_links (
    question_id   BIGINT,
    textbook_id   BIGINT,
    chapter_id    VARCHAR(32),
    status        VARCHAR(16),  -- active/deprecated
    PRIMARY KEY (question_id, textbook_id)
);

-- 教材改版时：
-- 1. 标记旧版本的题目关联
UPDATE question_textbook_links
SET status = 'deprecated'
WHERE textbook_id = ?;

-- 2. 跑新旧教材差异分析
-- 新教材删了哪些知识点 → 对应的题目可以归档
-- 新教材加了哪些知识点 → 题库需要补新题
```

### 7.4 题目"新鲜度"衰减

一道题目被搜索和使用的频率，随时间自然衰减。这影响搜索排序。

```python
def freshness_decay(question: dict) -> float:
    """计算题目新鲜度衰减因子"""
    days_since_published = (now() - question["published_at"]).days
    days_since_last_search = (now() - question["last_searched_at"]).days

    # 发布时间衰减（90 天半衰期）
    publish_decay = 0.5 ** (days_since_published / 90)

    # 搜索活跃度衰减（30 天半衰期）
    search_decay = 0.5 ** (days_since_last_search / 30)

    return publish_decay * 0.3 + search_decay * 0.7
```

---

## 八、你可能没想到的问题

### 8.1 题目安全性（防泄露）

如果题库是产品的核心壁垒，需要防止被竞争对手"反向扒库"：

| 措施 | 做法 |
|------|------|
| **API 限流** | 同一用户/IP 每天最多搜 N 题（针对爬虫行为） |
| **内容水印** | 展示的解析中加入不可见的追踪标识（零宽字符/隐写） |
| **题目内容分片** | API 不返回完整解析，鼓励在 App 内查看 |
| **异常检测** | 检测"遍历式搜索"（按 ID 递增请求），标记并封禁 |
| **CDN 防盗链** | 题目图片设置 Referer 白名单 |

### 8.2 不同年级的同一知识点

"(a+b)² 展开"这个知识点：
- 初中版：用具体数字举例 → $(3+4)^2 = 49$
- 高中版：用字母代数化 → $(a+b)^2 = a^2+2ab+b^2$

同一道题可能对不同年级展示不同风格的解析。在 `solution_versions` 表中用 `difficulty_level` 字段区分：

```sql
solution_versions
  ├── question_id=42, difficulty_level=basic     (初中生看)
  ├── question_id=42, difficulty_level=intermediate (高中生看)
  └── question_id=42, difficulty_level=advanced   (竞赛生看)
```

### 8.3 题目依赖关系

一道题可能是另一道题的"前置"或"变式"。比如：

- 题 A：基础题 → 题 B：题 A 的变式
- 题 C：题 D 需要用到的公式

```sql
CREATE TABLE question_relations (
    from_question_id BIGINT,
    to_question_id   BIGINT,
    relation_type    VARCHAR(32),  -- prerequisite/variant/follow_up/uses_formula
    PRIMARY KEY (from_question_id, to_question_id)
);
```

**用处**：
- 错题本推荐："你错的是这题的变式，看看这道原题"
- 知识点前置提醒："做这题需要先掌握 XXX，点击学习"

### 8.4 题目难度漂移

一道题原本标为"中等"，随着时间推移，可能因为：
- 教学大纲调整，这题变成"基础题"
- 用户正确率普遍很高（其实偏简单）

```python
def detect_difficulty_drift(question_id: int):
    """检测实际难度是否与标注不符"""
    stats = db.query("""
        SELECT
            AVG(CASE WHEN is_correct THEN 1 ELSE 0 END) as accuracy_rate,
            COUNT(*) as total_attempts
        FROM user_answer_logs
        WHERE question_id = ? AND created_at > NOW() - INTERVAL 30 DAY
    """, question_id)

    if stats.total_attempts < 100:
        return None  # 样本不足不调整

    if stats.accuracy_rate > 0.9 and question.difficulty != "easy":
        alert("题目 {} 正确率 {:.0%}，建议降为 easy".format(
            question_id, stats.accuracy_rate))
    elif stats.accuracy_rate < 0.3 and question.difficulty != "hard":
        alert("题目 {} 正确率 {:.0%}，建议升为 hard".format(
            question_id, stats.accuracy_rate))
```

### 8.5 图文匹配问题

题目中有图，但图片丢失了（抓取的 URL 失效），这道题就废了。

**方案**：
- 所有外链图片在入库时**下载到自有 OSS**，不依赖源站
- 定期巡检图片可用性（HEAD 请求），死链告警
- 对重要的图（如函数图像、几何图），尝试用 AI 重新生成（DALL-E/根据 LaTeX 用 Python 画图）

### 8.6 题目"一题多问"的处理

有的题目是多问结构：
> (1) 求 f(x) 的导数
> (2) 求 f(x) 的单调区间
> (3) 求 f(x) 在 [-2,2] 上的最大值

**不要拆成 3 道题**——它们共享上下文。在 `questions` 表加一个字段：

```sql
ALTER TABLE questions ADD COLUMN parent_question_id BIGINT;
-- NULL = 独立题，非 NULL = 子问题
```

### 8.7 解析的"教学风格"适配

不同用户偏好的解析风格不同：
- "给我最快的方法" vs "给我最容易理解的方法"
- "详细到每一步" vs "要点提示即可"

可以将解析版本按风格打标：

```sql
ALTER TABLE solution_versions ADD COLUMN style VARCHAR(32);
-- 'fastest' / 'easiest' / 'detailed' / 'concise' / 'visual'(多图)
```

用户画像中记录偏好：

```sql
ALTER TABLE user_profiles ADD COLUMN solution_style_preference VARCHAR(32);
```

用户反馈（👍/👎 / 切换其他解法）自动修正偏好。

### 8.8 法律与版权

| 风险 | 措施 |
|------|------|
| 抓取内容侵权 | 做"转换性使用"——用自己的格式重新组织，不是直接复制 |
| 用户上传的侵权内容 | DMCA 式通知-删除机制 |
| 教材截图侵权 | 不直接用教材截图，改为文字描述 + 自己画的图 |
| 竞品诉不正当竞争 | 不使用反编译、逆向手段获取题库；只抓公开内容 |

### 8.9 多语言题目

如果要支持港澳台/海外华人，涉及：
- 繁体中文（OCR 要支持繁体）
- 英文数学题（terminology 不同："一元二次方程" = "quadratic equation"）
- 题目本身多语言，但解析和标签体系需要支持跨语言映射

```sql
CREATE TABLE question_translations (
    question_id BIGINT,
    language    VARCHAR(8),    -- zh-CN/zh-TW/en
    content     TEXT,
    PRIMARY KEY (question_id, language)
);
```

### 8.10 灾备与回滚

```bash
# 每日全量备份（PG dump）
pg_dump -Fc question_bank > backup_$(date +%Y%m%d).dump

# 增量备份（WAL 归档）
archive_command = 'rsync %p backup-server:/wal/%f'

# 题目数据变更日志（审计）
CREATE TABLE question_audit_log (
    id           BIGSERIAL PRIMARY KEY,
    question_id  BIGINT,
    action       VARCHAR(16),   -- create/update/delete/deprecate
    changed_by   BIGINT,        -- 操作人
    old_data     JSONB,
    new_data     JSONB,
    created_at   TIMESTAMP DEFAULT NOW()
);
```

---

## 九、运维与监控

### 9.1 关键监控

| 监控项 | 指标 | 告警阈值 |
|--------|------|----------|
| 新题入库速率 | 每小时入库量 | < 100/h （抓取挂了） |
| 待审积压池 | pending_review 状态的题目数 | > 1000 |
| 低分题比例 | quality < 70 | > 10% |
| 解析覆盖率 | 有解析的题 / 总题 | < 90% |
| OCR 成功率 | OCR 成功 / 总图片 | < 80% |
| 搜索命中率 | 搜索结果 Top-1 相关 | < 65% |
| 图片可用率 | 图片 200 OK | < 99% |
| 抓取成功率 | 抓取 200 OK | < 70%（封 IP 了） |

### 9.2 告警响应

```
[监控告警] "待审积压池 > 1000，当前 1200"
  ↓
  1. 检查是否抓取管道爆发（新上了数据源？）
  2. 临时调高"自动通过"阈值，减少人工审核量
  3. 暂时关闭非核心数据源的抓取
  4. 增加审核人力（兼职老师）
```

---

## 十、分阶段实施路线

### Phase 1：种子库（0-1 个月）

- [x] 2-3 名教师录入核心科目，每科 5000 题
- [x] 建好知识点树 + 教材版本映射
- [x] PG 主表 + 图片 OSS
- [x] 人工审核流程 + Web 后台

### Phase 2：自动扩量（1-3 个月）

- [ ] 上线抓取管道 + 自动清洗
- [ ] ES 全文检索 + pgvector 语义检索
- [ ] LLM 辅助分类（置信度 > 0.85 自动分类）
- [ ] 自动验收（格式 + 答案验证）

### Phase 3：精细化运营（3-6 个月）

- [ ] 用户反馈闭环（👍/👎 驱动质量排序）
- [ ] 难点漂移检测 + 自动调整
- [ ] 多版本解析共存 + AB 测试
- [ ] 主动监测热门缺题 → AI 自动补全
- [ ] 教材改版影响分析工具

### Phase 4：智能化（6-12 个月）

- [ ] AI 自动生成裸题补齐长尾
- [ ] 个性化推荐（薄弱点 → 针对性题目）
- [ ] 解析风格个性化（最快/最详细/最多图）
- [ ] 跨学科关联（数学知识点 → 对应的物理题）
- [ ] 题目质量预测模型（入库前预判质量，决定审核优先级）
