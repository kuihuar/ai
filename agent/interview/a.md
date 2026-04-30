# AI开发相关面试核心知识点（代码\+实践\+区别）

# 一、Skill 实现（代码层面深入）

## 1\. 核心结构

每个Skill本质是一个独立目录，核心包含2类关键内容，代码/目录结构如下：

```plain text
# Skill目录结构（核心）
skill-demo/          # 单个Skill根目录
├─ SKILL.md          # 核心配置+规则文件（必选）
└─ reference/        # 参考文件目录（按需，可选）
    └─ demo-doc.md   # 参考文档，按需加载
```

## 2\. 核心文件 SKILL\.md 实现

用 YAML Frontmatter 定义元信息，Markdown 编写执行规则，示例代码：

```yaml
---
# YAML Frontmatter 元信息（第一级加载内容）
name: "PDF预览渲染修复Skill"  # Skill名称，供LLM决策选择
description: "用于修复PDF预览时首屏白页、拖拽渲染失败的问题，适配Vue组件渲染时序"  # 功能描述
---

# Markdown 执行规则（第二级加载内容）
## 执行前提
1.  需获取前端Vue组件（pdf-document-viewer.vue）源码
2.  需开启Chrome Devtools MCP权限，用于调试DOM渲染时序

## 执行步骤
1.  检查loadDocument()方法中forceRender()调用时机，确认documentLoading状态
2.  在documentLoading切换为false后，添加nextTick()等待DOM挂载
3.  重构渲染调度逻辑，改为串行调度，避免resize时任务打断
4.  增加宽度变化校验，仅在真实宽度变化时触发重绘

## 参考依赖
- 依赖Vue3的nextTick API
- 依赖Chrome Devtools MCP的页面操控、日志获取能力
```

## 3\. 加载机制（三级渐进加载，核心设计）

核心目的：减少LLM的context占用，节省token，代码层面核心实现（基于SkillRegistry）：

```java
// 核心类：SkillRegistry（技能注册器）
public class SkillRegistry {
    // 缓存Skill元信息（仅name+description，第一级加载）
    private Map<String, SkillMeta&gt; skillMetaCache = new ConcurrentHashMap<>();
    // 缓存完整Skill内容（SKILL.md全量，第二级加载）
    private Map<String, String> fullSkillCache = new ConcurrentHashMap<>();
    // 缓存reference参考文件（第三级，按需加载）
    private Map<String, List<String>> referenceCache = new ConcurrentHashMap<>();

    // 项目启动时扫描所有Skill目录，初始化元信息缓存（第一级加载）
    @PostConstruct
    public void init() {
        // 扫描指定目录下所有Skill目录
        File[] skillDirs = new File("skills").listFiles(File::isDirectory);
        for (File skillDir : skillDirs) {
            // 读取SKILL.md的YAML Frontmatter，提取name和description
            SkillMeta meta = parseSkillMeta(skillDir);
            skillMetaCache.put(meta.getName(), meta);
        }
    }

    // 第二级加载：获取完整SKILL.md内容
    public String getFullSkillContent(String skillName) {
        if (fullSkillCache.containsKey(skillName)) {
            return fullSkillCache.get(skillName);
        }
        // 读取SKILL.md全量内容，缓存后返回
        String fullContent = readSkillMd(skillName);
        fullSkillCache.put(skillName, fullContent);
        return fullContent;
    }

    // 第三级加载：按需加载reference目录下的参考文件
    public List<String> getReferenceFiles(String skillName) {
        if (referenceCache.containsKey(skillName)) {
            return referenceCache.get(skillName);
        }
        // 读取reference目录下所有文件，缓存后返回
        List<String> references = readReferenceFiles(skillName);
        referenceCache.put(skillName, references);
        return references;
    }

    // 工作流执行时，注入Skill到LLM提示词
    public String injectSkillToPrompt(String skillName) {
        // 拼接Skill内容，生成系统提示词
        String fullSkill = getFullSkillContent(skillName);
        return String.format("请按照以下Skill规则执行任务：\n%s", fullSkill);
    }
}

// Skill元信息封装类
@Data
class SkillMeta {
    private String name;       // Skill名称
    private String description;// Skill描述
}
```

## 4\. 工作流中调用逻辑

工作流执行到LLM节点时，若配置了skillName，从SkillRegistry获取对应Skill，拼接成系统提示词注入LLM：

```java
// 工作流LLM节点执行器
public class LLMNodeExecutor {
    @Autowired
    private SkillRegistry skillRegistry;
    @Autowired
    private ChatClient chatClient;

    public String execute(LLMNodeConfig config, String userPrompt) {
        // 1. 若配置了skillName，获取并注入Skill
        String systemPrompt = "你是一个AI开发助手，负责完成代码开发、调试任务";
        if (StringUtils.isNotBlank(config.getSkillName())) {
            systemPrompt = skillRegistry.injectSkillToPrompt(config.getSkillName());
        }

        // 2. 拼接完整提示词，调用LLM
        String fullPrompt = String.format("%s\n用户需求：%s", systemPrompt, userPrompt);
        return chatClient.generate(fullPrompt);
    }
}
```

# 二、AI辅助开发实践经验（面试重点表述）

近一个月所有项目（技术派、派聪明RAG、PaiFlow Agent等）均采用AI辅助开发，核心实践围绕2类工具展开，落地性极强：

## 1\. OpenAI Codex（核心工具，代码层面主力）

- IDE集成：在IntelliJ IDEA安装Codex插件，直接对接ChatGPT账号或OpenAI API密钥，无需额外配置，适配Java后端项目（解决VSCode无法配置JDK、Maven的痛点）。

- 核心应用场景：
            

    - 后端代码生成：快速生成CRUD、接口、中间件代码，适配Spring Boot、Gin等框架，减少重复编码。

    - 脚本开发：例如为派聪明RAG项目编写infra\.sh脚本，一键启动MinIO、ElasticSearch、Kafka，自动处理PID管理、日志输出、端口校验，且能自主修复脚本中的进程识别问题。

    - 可视化调试：配合Chrome Devtools MCP，自主操控浏览器，复现PDF预览等动态渲染Bug（如首屏白页、拖拽报错），自主抓取DOM状态、控制台日志，定位时序问题并修复，无需手动描述复杂Bug现象。

## 2\. Claude 系列（辅助补充，复杂场景优势）

- Claude（文本理解）：解决线上难以描述的Bug，文本理解能力优于Codex，能快速get真实需求，避免盲目调试。

- Claude Code（最强Agent工具）：端到端完成开发闭环，从需求拆解、代码编写、测试执行到Bug修复一条龙，适合复杂任务、多文件重构，搭配国内模型（如GLM\-5\.1）可提升战斗力，曾用其完成简历Agent项目。

## 3\. 核心实践心得

AI辅助开发的核心是“人主导、AI干活”：用AI替代80%的重复劳动（编码、脚本、调试），人聚焦架构设计、业务逻辑校验，同时需注意AI生成内容的Review，避免逻辑漏洞和安全问题。

# 三、Agent 现状与预期

## 1\. 现状：未达预期

当前主流Agent（包括Claude Code、Qoder专家团模式）仍需人工介入，无法完全自主完成完整功能开发，需人全程盯着，处理异常情况、校验结果。

## 2\. 理想预期

Agent能自主完成全流程开发：理解需求 → 设计方案 → 编写代码 → 运行测试 → 修复Bug → 提交PR，全程无需人工干预，成为真正的“AI开发助手”。

## 3\. 趋势

Agent进步速度极快，预计未来1年（如OpenClaw、爱马仕Agent）将实现突破性进展，逐步接近理想预期。

# 四、多模态知识检索（派聪明RAG项目实操）

传统RAG仅支持文本，多模态检索核心是“将非文本内容转为可向量化表示”，再沿用RAG流程，核心实现分3类场景：

## 1\. 核心通用流程（代码示例）

```java
// 多模态文档处理核心方法（派聪明项目实际代码）
public void vectorizeFile(MultipartFile file, String userId, String orgTag, boolean isPublic) {
    // 1. 文档解析和分块（适配多类型文件：文本、PDF、音频、视频）
    List<Chunk> chunks = fileParsingService.parseAndChunk(file);
    // 2. 批量向量化（根据文件类型选择对应Embedding模型）
    List<Vector> vectors = embeddingClient.batchEmbedding(chunks);
    // 3. 构建文档对象（关联元信息：文件类型、时间戳、用户ID等）
    List<Document> documents = buildDocuments(chunks, vectors, buildMetadata(file, userId, orgTag));
    // 4. 批量存储到向量库（Elasticsearch）
    elasticsearchService.bulkIndex(documents);
}

// 元信息构建
private Map<String, Object> buildMetadata(MultipartFile file, String userId, String orgTag) {
    Map<String, Object> metadata = new HashMap<>();
    metadata.put("userId", userId);
    metadata.put("orgTag", orgTag);
    metadata.put("fileType", file.getContentType());
    metadata.put("uploadTime", LocalDateTime.now());
    return metadata;
}
```

## 2\. 具体场景实现

- 场景1：音频/无画面视频 → 转文本检索


    - 流程：音频 → ASR（自动语音识别）转文字稿 → 按时间戳分块 → 文本Embedding → 存入向量库。

    - 查询：用户提问 → 文本向量化 → 检索相关文字片段 → 返回回答\+对应时间戳。

- 场景2：带画面视频 → 多模态检索
            

    - 流程：视频抽帧（每秒1帧/关键帧） → 多模态Embedding模型（如CLIP、Chinese\-CLIP）将图像转向量 → 与文本分块向量一同存入向量库。

    - 查询：用户提问（文本） → 多模态模型将问题向量化 → 检索相关图像/文本片段 → 定位视频时间戳。

- 场景3：PDF/Word中的图表 → 多模态Embedding
            

    - 用多模态模型解析图表内容，转为向量，与文本内容关联存储，实现“文本\+图表”联合检索。

## 3\. 核心模型选择

多模态Embedding模型优先选择：OpenAI CLIP（通用场景）、Chinese\-CLIP（中文场景）、BGE\-M3（多模态版本，适配中文\+多类型文件）。

# 五、A2A 与 MCP 核心区别（面试高频）

两者定位不同，互补而非竞争，核心区别用通俗表述\+实例说明，便于理解：

|对比维度|MCP（Model Control Protocol）|A2A（Agent to Agent）|
|---|---|---|
|核心定位|解决“单个Agent如何调用工具”，相当于“AI的USB接口”|解决“多个Agent如何协作”，相当于“Agent之间的通信协议”|
|核心功能|Agent通过MCP Client连接MCP Server，调用Server暴露的工具、资源、提示词能力|每个Agent发布Agent Card（JSON格式）描述自身能力，其他Agent通过Card发现并调用它|
|实例说明|库存Agent通过MCP调用数据库工具，查询当前库存数量|库存Agent通过A2A调用供应商的订货Agent，当库存不足时自动触发下单|
|核心总结|管“我怎么用工具”|管“我怎么找别人帮忙”|

# 六、RAG 长上下文解决方案（实操重点）

当检索回几十个chunk，导致上下文过长时，采用“三层递进”解决方案，落地性强，已在派聪明项目中实践：

## 1\. 第一层：检索前置过滤（减少无效chunk）

- 小粒度chunk召回：用256 token的小粒度chunk做向量检索，初步召回30条候选chunk。

- 分层召回\+重排：
           

    - 第一轮：向量检索快速缩小范围；

    - 第二轮：关键词匹配精确定位；

    - 第三轮：用Reranker（重排模型）打分，仅保留Top5\~Top8最相关的chunk，减少上下文占用。

## 2\. 第二层：Context Window 预算管理（合理分配token）

按比例分配token额度（以8K context为例），避免超出限制：

- 系统提示词：20%（1\.6K token）；

- 检索结果：50%（4K token）；

- 对话历史：30%（2\.4K token）。

对话历史处理：最近3轮保留原文，更早的对话压缩为200 token的渐进式摘要，避免信息丢失。

## 3\. 第三层：长文档特殊处理（单条超8K token场景）

- 方案一：Map\-Reduce：将长文档切分多个chunk，每个chunk单独查询LLM，最后合并回答（缺点：丢失跨chunk关联信息）。

- 方案二：RAG \+ 递归摘要：先对长文档分层摘要（5个chunk生成1个中间摘要，5个中间摘要生成1个高层摘要），检索时先匹配高层摘要定位章节，再深入匹配原始chunk。

- 方案三：使用长上下文模型：如Claude Code（支持1M上下文），直接处理长文档，无需切分。

# 七、核心项目经验（面试精简版）

## 1\. PaiAgent/PaiFlow（2026\-03 至今）

企业级AI工作流平台，基于LangGraph4j \+ Spring AI，支持可视化拖拽编排大模型和工具节点。

- 技术栈：Java 21、Spring Boot 3\.4、Spring AI 1\.0、LangGraph4j 1\.8、React 18、ReactFlow。

- 核心实现：
            

    - 基于LangGraph4j StateGraph构建工作流引擎，实现节点注册、状态传递。

    - 设计ChatClientFactory动态工厂，实现多厂商LLM（OpenAI、DeepSeek等）无缝切换。

    - 用模板方法模式重构LLM节点执行器，代码量从800\+行精简至每个10行。

    - 实现Skill预置知识包机制，支持自动加载、缓存、渐进式注入。

## 2\. 派聪明 RAG 知识库（2026\-01 ～ 2026\-02）

企业级智能对话平台，支持上传文档构建私有知识库，通过自然语言交互查询。

- 技术栈：SpringBoot、MySQL、Redis、Apache Tika、Ollama、Elasticsearch、MinIO、Kafka。

- 核心实现：
            

    - 基于Elasticsearch实现“关键词\+语义”双引擎检索，集成阿里Embedding模型。

    - 编写Shell脚本，一键启动Kafka KRaft模式，处理集群ID冲突。

    - 基于WebSocket \+ LLM Stream API，实现“打字机”式流式响应。

    - 引入MCP协议，实现Agent与工具生态解耦。