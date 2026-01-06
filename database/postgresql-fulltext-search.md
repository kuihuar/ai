# PostgreSQL 全文搜索

## 目录
- [全文搜索基础](#全文搜索基础)
- [文本搜索类型](#文本搜索类型)
- [搜索配置](#搜索配置)
- [索引优化](#索引优化)
- [高级功能](#高级功能)
- [多语言支持](#多语言支持)

---

## 全文搜索基础

### 什么是全文搜索

全文搜索是在文本数据中查找关键词的功能，比简单的 LIKE 查询更强大和高效。

### 基本概念

- **文档（Document）**：要搜索的文本
- **词位（Lexeme）**：经过词干提取和规范化的词
- **tsvector**：文档的向量表示
- **tsquery**：搜索查询

---

## 文本搜索类型

### tsvector

`tsvector` 是文档的向量表示，包含词位和位置信息。

```sql
-- 创建 tsvector
SELECT 'The quick brown fox jumps over the lazy dog'::tsvector;

-- 输出：'The' 'brown' 'dog' 'fox' 'jumps' 'lazy' 'over' 'quick' 'the'

-- 使用 to_tsvector 函数（会进行词干提取和规范化）
SELECT to_tsvector('english', 'The quick brown fox jumps over the lazy dog');

-- 输出：'brown':3 'dog':9 'fox':4 'jump':5 'lazi':8 'quick':2
```

### tsquery

`tsquery` 是搜索查询，支持布尔运算符。

```sql
-- 创建 tsquery
SELECT 'fox & dog'::tsquery;

-- 使用 to_tsquery 函数
SELECT to_tsquery('english', 'fox & dog');

-- 布尔运算符
SELECT 'fox & dog'::tsquery;        -- AND
SELECT 'fox | dog'::tsquery;        -- OR
SELECT '!cat'::tsquery;              -- NOT
SELECT '(fox | dog) & !cat'::tsquery; -- 组合
```

### 执行搜索

```sql
-- 基本搜索
SELECT 
    title,
    content
FROM articles
WHERE to_tsvector('english', content) @@ to_tsquery('english', 'postgresql');

-- 使用列存储 tsvector（推荐）
ALTER TABLE articles ADD COLUMN content_vector tsvector;

UPDATE articles 
SET content_vector = to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(content, ''));

-- 搜索
SELECT title, content
FROM articles
WHERE content_vector @@ to_tsquery('english', 'postgresql');
```

---

## 搜索配置

### 内置配置

PostgreSQL 提供多种语言配置：

```sql
-- 查看可用配置
SELECT cfgname FROM pg_ts_config;

-- 常用配置
SELECT to_tsvector('english', 'running runs ran');
SELECT to_tsvector('simple', 'running runs ran');  -- 不进行词干提取
```

### 自定义配置

```sql
-- 创建文本搜索配置
CREATE TEXT SEARCH CONFIGURATION my_config (COPY = english);

-- 添加同义词
CREATE TEXT SEARCH DICTIONARY my_synonym (
    TEMPLATE = synonym,
    SYNONYMS = my_synonyms
);

-- 使用同义词
ALTER TEXT SEARCH CONFIGURATION my_config
    ALTER MAPPING FOR asciiword WITH my_synonym, english_stem;
```

---

## 索引优化

### GIN 索引

```sql
-- 创建 GIN 索引（推荐）
CREATE INDEX idx_articles_content_gin 
ON articles USING GIN (to_tsvector('english', content));

-- 使用表达式索引
CREATE INDEX idx_articles_content_gin 
ON articles USING GIN (to_tsvector('english', title || ' ' || content));

-- 使用存储的 tsvector 列
CREATE INDEX idx_articles_content_vector_gin 
ON articles USING GIN (content_vector);
```

### 触发器自动更新

```sql
-- 创建触发器函数
CREATE OR REPLACE FUNCTION update_content_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.content_vector := to_tsvector('english', 
        COALESCE(NEW.title, '') || ' ' || COALESCE(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER update_articles_content_vector
    BEFORE INSERT OR UPDATE ON articles
    FOR EACH ROW
    EXECUTE FUNCTION update_content_vector();
```

---

## 高级功能

### 排名（Ranking）

```sql
-- ts_rank：基于词频排名
SELECT 
    title,
    ts_rank(content_vector, query) AS rank
FROM articles, to_tsquery('english', 'postgresql') query
WHERE content_vector @@ query
ORDER BY rank DESC;

-- ts_rank_cd：考虑词距离的排名
SELECT 
    title,
    ts_rank_cd(content_vector, query) AS rank
FROM articles, to_tsquery('english', 'postgresql') query
WHERE content_vector @@ query
ORDER BY rank DESC;

-- 加权排名
SELECT 
    title,
    ts_rank_cd(
        setweight(to_tsvector('english', title), 'A') ||
        setweight(to_tsvector('english', content), 'B'),
        query
    ) AS rank
FROM articles, to_tsquery('english', 'postgresql') query
WHERE content_vector @@ query
ORDER BY rank DESC;
```

### 高亮显示

```sql
-- ts_headline：高亮匹配的文本
SELECT 
    title,
    ts_headline('english', content, query) AS headline
FROM articles, to_tsquery('english', 'postgresql') query
WHERE content_vector @@ query;

-- 自定义高亮选项
SELECT 
    title,
    ts_headline(
        'english', 
        content, 
        query,
        'StartSel=<mark>, StopSel=</mark>, MaxWords=35, MinWords=15'
    ) AS headline
FROM articles, to_tsquery('english', 'postgresql') query
WHERE content_vector @@ query;
```

### 短语搜索

```sql
-- 使用 <-> 操作符搜索短语
SELECT title, content
FROM articles
WHERE content_vector @@ to_tsquery('english', 'postgresql <-> database');

-- 使用 <N> 指定词距离
SELECT title, content
FROM articles
WHERE content_vector @@ to_tsquery('english', 'postgresql <2> database');
```

### 前缀搜索

```sql
-- 使用 :* 进行前缀搜索
SELECT title, content
FROM articles
WHERE content_vector @@ to_tsquery('english', 'post:*');
```

---

## 多语言支持

### 不同语言的配置

```sql
-- 中文搜索（需要安装 zhparser 扩展）
CREATE EXTENSION IF NOT EXISTS zhparser;

CREATE TEXT SEARCH CONFIGURATION chinese_parser (PARSER = zhparser);

-- 英文搜索
SELECT to_tsvector('english', 'The quick brown fox');

-- 中文搜索
SELECT to_tsvector('chinese_parser', '快速棕色狐狸');
```

### 混合语言

```sql
-- 处理混合语言内容
SELECT to_tsvector('english', 'PostgreSQL 是一个强大的数据库系统');
```

---

## 实践示例

### 示例 1：文章搜索系统

```sql
-- 创建表
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    content_vector tsvector,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_articles_content_vector_gin 
ON articles USING GIN (content_vector);

-- 创建触发器
CREATE OR REPLACE FUNCTION update_articles_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.content_vector := to_tsvector('english', 
        COALESCE(NEW.title, '') || ' ' || COALESCE(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_articles_vector_trigger
    BEFORE INSERT OR UPDATE ON articles
    FOR EACH ROW
    EXECUTE FUNCTION update_articles_vector();

-- 插入数据
INSERT INTO articles (title, content) VALUES
    ('PostgreSQL Full-Text Search', 'PostgreSQL provides powerful full-text search capabilities.'),
    ('Database Performance', 'Optimizing database queries is essential for performance.');

-- 搜索
SELECT 
    id,
    title,
    ts_rank_cd(content_vector, query) AS rank
FROM articles, to_tsquery('english', 'postgresql & search') query
WHERE content_vector @@ query
ORDER BY rank DESC;
```

### 示例 2：多字段搜索

```sql
-- 搜索多个字段
SELECT 
    id,
    title,
    content,
    ts_rank_cd(
        setweight(to_tsvector('english', title), 'A') ||
        setweight(to_tsvector('english', content), 'B'),
        query
    ) AS rank
FROM articles, to_tsquery('english', 'postgresql | database') query
WHERE 
    to_tsvector('english', title) @@ query OR
    to_tsvector('english', content) @@ query
ORDER BY rank DESC;
```

### 示例 3：搜索函数

```sql
-- 创建搜索函数
CREATE OR REPLACE FUNCTION search_articles(search_term TEXT)
RETURNS TABLE (
    id INTEGER,
    title VARCHAR,
    content TEXT,
    rank REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id,
        a.title,
        a.content,
        ts_rank_cd(a.content_vector, query) AS rank
    FROM articles a, to_tsquery('english', search_term) query
    WHERE a.content_vector @@ query
    ORDER BY rank DESC
    LIMIT 20;
END;
$$ LANGUAGE plpgsql;

-- 使用
SELECT * FROM search_articles('postgresql & search');
```

---

## 性能优化

### 索引策略

```sql
-- 使用存储的 tsvector 列（推荐）
ALTER TABLE articles ADD COLUMN content_vector tsvector;
CREATE INDEX idx_articles_content_vector_gin ON articles USING GIN (content_vector);

-- 使用触发器自动更新
CREATE TRIGGER update_content_vector
    BEFORE INSERT OR UPDATE ON articles
    FOR EACH ROW
    EXECUTE FUNCTION update_content_vector();
```

### 查询优化

```sql
-- ✅ 好：使用索引列
SELECT * FROM articles 
WHERE content_vector @@ to_tsquery('english', 'postgresql');

-- ❌ 差：每次计算 tsvector
SELECT * FROM articles 
WHERE to_tsvector('english', content) @@ to_tsquery('english', 'postgresql');
```

---

## 最佳实践

1. **使用存储的 tsvector 列**
   - 预先计算并存储 tsvector
   - 使用触发器自动更新

2. **创建 GIN 索引**
   - 对 tsvector 列创建 GIN 索引
   - 提高搜索性能

3. **使用排名**
   - 使用 ts_rank 或 ts_rank_cd 排序结果
   - 考虑使用加权排名

4. **优化查询**
   - 避免在 WHERE 子句中计算 tsvector
   - 使用存储的 tsvector 列

---

## 下一步学习

- [PostgreSQL 扩展与插件](./postgresql-extensions.md)
- [PostgreSQL 索引与查询优化](./postgresql-indexes-optimization.md)
- [PostgreSQL 最佳实践](./postgresql-best-practices.md)

---

*最后更新：2024年*

