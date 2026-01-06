# OpenTelemetry 追踪数据存储到数据库

## 概述

OpenTelemetry 本身不直接支持将追踪数据存储到数据库，但可以通过以下方案实现：

1. **Jaeger + 数据库存储**（推荐）
2. **OpenTelemetry Collector + 自定义 Exporter**
3. **自定义 Exporter**

## 方案 1: Jaeger + MySQL/PostgreSQL（推荐）

### 优点
- ✅ 成熟稳定，生产环境广泛使用
- ✅ 支持 MySQL、PostgreSQL、Cassandra、Elasticsearch
- ✅ 提供完整的 UI 界面
- ✅ 无需修改应用代码，只需配置 Jaeger

### 配置步骤

#### 1. 创建数据库和表结构

Jaeger 需要特定的表结构。可以使用 Jaeger 官方提供的 SQL 脚本：

**MySQL:**
```sql
-- 创建数据库
CREATE DATABASE jaeger CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE jaeger;

-- 下载并执行 Jaeger 的 MySQL schema
-- https://github.com/jaegertracing/jaeger/tree/main/plugin/storage/mysql
```

**PostgreSQL:**
```sql
-- 创建数据库
CREATE DATABASE jaeger;

-- 下载并执行 Jaeger 的 PostgreSQL schema
-- https://github.com/jaegertracing/jaeger/tree/main/plugin/storage/postgresql
```

#### 2. 启动 Jaeger with MySQL/PostgreSQL

**使用 Docker Compose（推荐）:**

```yaml
# docker-compose-jaeger-mysql.yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: jaeger
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./jaeger-mysql-schema.sql:/docker-entrypoint-initdb.d/init.sql

  jaeger:
    image: jaegertracing/all-in-one:latest
    environment:
      - SPAN_STORAGE_TYPE=mysql
      - MYSQL_HOST=mysql
      - MYSQL_PORT=3306
      - MYSQL_DB=jaeger
      - MYSQL_USER=root
      - MYSQL_PASSWORD=root
    ports:
      - "16686:16686"  # UI
      - "14268:14268"  # HTTP Collector
      - "4317:4317"    # OTLP gRPC
      - "4318:4318"    # OTLP HTTP
    depends_on:
      - mysql

volumes:
  mysql_data:
```

**使用 Docker 命令:**

```bash
# 启动 MySQL
docker run -d --name jaeger-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=jaeger \
  -p 3306:3306 \
  mysql:8.0

# 初始化 Jaeger 表结构（需要先下载 schema 文件）
# mysql -h 127.0.0.1 -u root -p jaeger < jaeger-mysql-schema.sql

# 启动 Jaeger with MySQL
docker run -d --name jaeger \
  -e SPAN_STORAGE_TYPE=mysql \
  -e MYSQL_HOST=127.0.0.1 \
  -e MYSQL_PORT=3306 \
  -e MYSQL_DB=jaeger \
  -e MYSQL_USER=root \
  -e MYSQL_PASSWORD=root \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 4317:4317 \
  jaegertracing/all-in-one:latest
```

**PostgreSQL 版本:**

```bash
# 启动 PostgreSQL
docker run -d --name jaeger-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=jaeger \
  -p 5432:5432 \
  postgres:15

# 启动 Jaeger with PostgreSQL
docker run -d --name jaeger \
  -e SPAN_STORAGE_TYPE=postgresql \
  -e POSTGRES_HOST=127.0.0.1 \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_DB=jaeger \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 4317:4317 \
  jaegertracing/all-in-one:latest
```

#### 3. 配置应用使用 Jaeger

在 `configs/config.yaml` 中配置：

```yaml
tracing:
  service_name: "sre"
  service_version: "v1.0.0"
  environment: "dev"
  jaeger_endpoint: "http://localhost:14268/api/traces"  # Jaeger HTTP Collector
  # 或者使用 OTLP（如果 Jaeger 支持）
  # otlp_endpoint: "localhost:4317"
  sampling_ratio: 1.0
```

#### 4. 验证

1. 访问 Jaeger UI: http://localhost:16686
2. 执行一些请求，生成追踪数据
3. 在 Jaeger UI 中查看追踪数据
4. 检查数据库中的表是否有数据：

```sql
-- MySQL
USE jaeger;
SHOW TABLES;
SELECT COUNT(*) FROM traces;
SELECT COUNT(*) FROM spans;
```

### 表结构说明

Jaeger 在数据库中创建的主要表：

- `traces`: 追踪信息
- `spans`: Span 信息
- `operations`: 操作信息
- `dependencies`: 服务依赖关系

## 方案 2: OpenTelemetry Collector + 自定义 Exporter

### 架构

```
应用 → OTLP → OpenTelemetry Collector → 自定义 Exporter → 数据库
```

### 实现步骤

#### 1. 创建自定义数据库 Exporter

```go
// internal/tracing/exporter/database_exporter.go
package exporter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type DatabaseExporter struct {
	db *sql.DB
}

func NewDatabaseExporter(db *sql.DB) *DatabaseExporter {
	return &DatabaseExporter{db: db}
}

func (e *DatabaseExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	// 实现将 spans 写入数据库的逻辑
	for _, span := range spans {
		if err := e.saveSpan(ctx, span); err != nil {
			return err
		}
	}
	return nil
}

func (e *DatabaseExporter) Shutdown(ctx context.Context) error {
	return e.db.Close()
}

func (e *DatabaseExporter) saveSpan(ctx context.Context, span trace.ReadOnlySpan) error {
	// 提取 span 信息
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	parentSpanID := span.Parent().SpanID().String()
	
	// 获取服务名
	serviceName := "unknown"
	if resource := span.Resource(); resource != nil {
		if attr := resource.Attributes(); attr != nil {
			for _, a := range attr {
				if a.Key == semconv.ServiceNameKey {
					serviceName = a.Value.AsString()
					break
				}
			}
		}
	}

	// 插入数据库
	query := `
		INSERT INTO spans (
			trace_id, span_id, parent_span_id, 
			service_name, operation_name, 
			start_time, duration_ms, 
			tags, logs, status_code, status_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := e.db.ExecContext(ctx, query,
		traceID, spanID, parentSpanID,
		serviceName, span.Name(),
		span.StartTime(), span.EndTime().Sub(span.StartTime()).Milliseconds(),
		// tags, logs 需要序列化为 JSON
		span.Status().Code.String(), span.Status().Description,
	)
	
	return err
}
```

#### 2. 创建数据库表结构

```sql
CREATE TABLE spans (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trace_id VARCHAR(32) NOT NULL,
    span_id VARCHAR(16) NOT NULL,
    parent_span_id VARCHAR(16),
    service_name VARCHAR(255) NOT NULL,
    operation_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    duration_ms BIGINT NOT NULL,
    tags JSON,
    logs JSON,
    status_code VARCHAR(50),
    status_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_trace_id (trace_id),
    INDEX idx_service_name (service_name),
    INDEX idx_start_time (start_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 3. 在 provider.go 中添加数据库导出器支持

需要修改 `internal/tracing/provider.go` 以支持数据库导出器。

## 方案 3: 使用现有的数据库存储方案

### SigNoz

[SigNoz](https://signoz.io/) 是一个开源的 APM 工具，支持将数据存储到 ClickHouse 或 PostgreSQL。

### Tempo (Grafana)

[Grafana Tempo](https://grafana.com/docs/tempo/latest/) 支持多种后端存储，包括对象存储和数据库。

## 推荐方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **Jaeger + MySQL/PostgreSQL** | 成熟稳定、有 UI、无需改代码 | 需要额外部署 Jaeger | 生产环境推荐 |
| **自定义 Exporter** | 完全控制、灵活 | 需要自己实现和维护 | 特殊需求 |
| **SigNoz/Tempo** | 功能完整、现代化 | 需要额外部署 | 需要完整 APM 功能 |

## 最佳实践

1. **生产环境推荐使用 Jaeger + MySQL/PostgreSQL**
   - 成熟稳定
   - 有完整的 UI 和查询功能
   - 社区支持好

2. **数据库选择**
   - **MySQL**: 适合中小规模，查询性能好
   - **PostgreSQL**: 适合大规模，JSON 支持更好
   - **Cassandra**: 适合超大规模分布式场景

3. **数据保留策略**
   - 设置数据保留时间（如 7 天、30 天）
   - 定期清理旧数据
   - 考虑使用分区表

4. **性能优化**
   - 为常用查询字段创建索引
   - 使用批量插入
   - 考虑异步写入

## 相关资源

- [Jaeger 存储后端文档](https://www.jaegertracing.io/docs/latest/deployment/#storage-backends)
- [Jaeger MySQL Schema](https://github.com/jaegertracing/jaeger/tree/main/plugin/storage/mysql)
- [Jaeger PostgreSQL Schema](https://github.com/jaegertracing/jaeger/tree/main/plugin/storage/postgresql)
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)

