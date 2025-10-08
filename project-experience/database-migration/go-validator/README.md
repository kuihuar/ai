# Go语言多数据库一致性验证工具

这是一个用Go语言编写的多数据库一致性验证工具，用于验证MySQL数据库从Azure迁移到AWS后的数据一致性。

## 功能特性

- **多实例数据库验证**: 支持Azure和AWS多个数据库实例的对比验证
- **实例级配置**: 每个数据库实例独立配置，支持不同的连接参数
- **并行验证**: 支持同时验证多个数据库对比对
- **大表分批处理**: 自动处理大表，避免内存溢出
- **详细日志记录**: 完整的验证过程日志
- **JSON报告生成**: 生成详细的验证报告
- **配置文件支持**: 支持JSON和YAML格式的配置文件
- **并发控制**: 可配置的最大并发数

## 项目结构

```
go-validator/
├── go.mod              # Go模块文件
├── main.go             # 主程序入口
├── types.go            # 数据结构定义
├── validator.go        # 验证器核心逻辑
├── config.go           # 配置文件处理
└── README.md           # 说明文档
```

## 安装和运行

### 1. 安装依赖

```bash
# 进入项目目录
cd go-validator

# 初始化Go模块（如果还没有）
go mod init multi-database-validator

# 安装MySQL驱动
go get github.com/go-sql-driver/mysql
```

### 2. 创建配置文件

```bash
# 创建JSON格式的默认配置文件
go run . init config.json

# 或创建YAML格式的默认配置文件
go run . init config.yaml
```

### 3. 编辑配置文件

编辑配置文件（`config.json` 或 `config.yaml`），设置正确的数据库连接信息：

#### JSON格式配置文件示例 (config.json)：

```json
{
  "azure": [
    {
      "name": "azure-db1",
      "host": "your-azure-mysql1.mysql.database.azure.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db1",
      "charset": "utf8mb4"
    },
    {
      "name": "azure-db2",
      "host": "your-azure-mysql2.mysql.database.azure.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db2",
      "charset": "utf8mb4"
    }
  ],
  "aws": [
    {
      "name": "aws-db1",
      "host": "your-aws-rds1.region.rds.amazonaws.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db1",
      "charset": "utf8mb4"
    },
    {
      "name": "aws-db2",
      "host": "your-aws-rds2.region.rds.amazonaws.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db2",
      "charset": "utf8mb4"
    }
  ],
  "max_workers": 3
}
```

#### YAML格式配置文件示例 (config.yaml)：

```yaml
# Azure数据库实例列表
azure:
  - name: azure-db1
    host: your-azure-mysql1.mysql.database.azure.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4
  
  - name: azure-db2
    host: your-azure-mysql2.mysql.database.azure.com
    user: your_username
    password: your_password
    database: db2
    charset: utf8mb4

# AWS数据库实例列表
aws:
  - name: aws-db1
    host: your-aws-rds1.region.rds.amazonaws.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4
  
  - name: aws-db2
    host: your-aws-rds2.region.rds.amazonaws.com
    user: your_username
    password: your_password
    database: db2
    charset: utf8mb4

# 最大并发数
max_workers: 3
```

### 4. 运行验证

```bash
# 使用配置文件运行
go run .

# 或者编译后运行
go build -o validator
./validator
```

## 使用方法

### 命令行选项

```bash
# 显示帮助信息
go run . help

# 创建默认配置文件
go run . init [filename]

# 运行验证（使用默认配置或config.json）
go run .
```

### 配置说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| azure[].name | Azure实例名称 | - |
| azure[].host | Azure数据库主机地址 | - |
| azure[].user | Azure数据库用户名 | - |
| azure[].password | Azure数据库密码 | - |
| azure[].database | Azure数据库名称 | - |
| azure[].charset | Azure数据库字符集 | utf8mb4 |
| aws[].name | AWS实例名称 | - |
| aws[].host | AWS数据库主机地址 | - |
| aws[].user | AWS数据库用户名 | - |
| aws[].password | AWS数据库密码 | - |
| aws[].database | AWS数据库名称 | - |
| aws[].charset | AWS数据库字符集 | utf8mb4 |
| max_workers | 最大并发数 | 3 |

**注意**: Azure和AWS实例数组的长度必须相同，工具会按索引顺序进行配对验证。

## 验证策略

### 1. 实例配对验证
- 按索引顺序配对Azure和AWS实例
- 验证每个实例对的数据一致性

### 2. 表数量对比
- 比较Azure和AWS环境中每个数据库的表数量
- 检查是否存在缺失的表

### 3. 数据一致性验证
- **小表（<10万行）**: 直接计算整个表的MD5校验和
- **大表（>=10万行）**: 分批读取数据，计算每批的校验和，最后合并

### 4. 空表处理
- 空表返回特殊标识 "empty_table"
- 确保空表的一致性

## 输出结果

### 1. 控制台输出
- 实时显示验证进度
- 显示验证摘要
- 显示错误和不一致信息

### 2. 日志文件
- 详细的验证过程日志
- 错误信息和调试信息

### 3. JSON报告
生成 `consistency_report.json` 文件，包含：
- 验证摘要统计
- 每个数据库对比对的详细结果
- 实例信息（Azure实例名称、AWS实例名称）
- 表级别的对比结果
- 错误信息列表

## 示例输出

### 控制台输出
```
数据库一致性验证完成
============================================================
验证时间: 2024-01-15T10:30:00Z
总耗时: 2m30s
总数据库数: 5
验证成功: 3
数据不一致: 1
验证错误: 1
成功率: 60.00%
============================================================
```

### JSON报告示例
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "total_databases": 2,
  "successful_validations": 1,
  "inconsistent_databases": 1,
  "error_databases": 0,
  "success_rate": "50.00%",
  "results": {
    "db1": {
      "database": "db1",
      "azure_instance": "azure-db1",
      "aws_instance": "aws-db1",
      "azure_tables": 10,
      "aws_tables": 10,
      "table_comparisons": [
        {
          "table": "users",
          "azure_checksum": "abc123...",
          "aws_checksum": "abc123...",
          "match": true
        }
      ],
      "status": "SUCCESS",
      "errors": [],
      "start_time": "2024-01-15T10:28:00Z",
      "end_time": "2024-01-15T10:29:30Z"
    }
  }
}
```

## 性能优化

### 1. 并发控制
- 使用goroutine并行验证多个数据库
- 可配置的最大并发数，避免资源耗尽

### 2. 大表处理
- 自动检测大表（>10万行）
- 分批读取数据，每批10000行
- 避免内存溢出

### 3. 连接管理
- 每个数据库使用独立的连接
- 自动关闭连接，避免连接泄漏

## 错误处理

### 1. 连接错误
- 数据库连接失败
- 连接超时
- 认证失败

### 2. 查询错误
- SQL语法错误
- 权限不足
- 表不存在

### 3. 数据错误
- 数据格式不一致
- 字符集问题
- 数据类型不匹配

## 注意事项

1. **网络连接**: 确保能够同时连接到Azure和AWS数据库
2. **权限要求**: 需要对所有数据库有SELECT权限
3. **资源消耗**: 大表验证会消耗较多内存和网络带宽
4. **时间考虑**: 大数据库验证可能需要较长时间
5. **字符集**: 确保两个环境的字符集设置一致

## 故障排除

### 常见问题

1. **连接失败**
   - 检查网络连通性
   - 验证数据库地址和端口
   - 确认用户名和密码

2. **权限错误**
   - 确保用户有SELECT权限
   - 检查数据库访问权限

3. **内存不足**
   - 减少max_workers数量
   - 分批验证数据库

4. **验证超时**
   - 增加数据库连接超时时间
   - 减少并发数

## 扩展功能

可以考虑添加的功能：
- 支持其他数据库类型（PostgreSQL、Oracle等）
- 增量验证（只验证变更的数据）
- 实时监控和告警
- Web界面
- 更多验证策略（行数对比、关键字段对比等）
