# 不一致输出格式更新总结

## 更新概述

根据用户需求，对Python和Go版本的多数据库一致性验证工具进行了更新，让不一致的结果输出更详细的数据库实例和表名信息。

## 更新内容

### 1. 数据结构更新

#### Python版本
在`TableComparison`结果中添加了以下字段：
```python
table_result = {
    'table': table,
    'azure_checksum': azure_checksum,
    'aws_checksum': aws_checksum,
    'match': azure_checksum == aws_checksum,
    'azure_instance': azure_instance['name'],      # 新增
    'aws_instance': aws_instance['name'],          # 新增
    'azure_database': azure_instance['database'],  # 新增
    'aws_database': aws_instance['database']       # 新增
}
```

#### Go版本
在`TableComparison`结构体中添加了以下字段：
```go
type TableComparison struct {
    Table         string `json:"table"`
    AzureChecksum string `json:"azure_checksum"`
    AWSChecksum   string `json:"aws_checksum"`
    Match         bool   `json:"match"`
    AzureInstance string `json:"azure_instance"`  // 新增
    AWSInstance   string `json:"aws_instance"`    // 新增
    AzureDatabase string `json:"azure_database"`  // 新增
    AWSDatabase   string `json:"aws_database"`    // 新增
}
```

### 2. 日志输出更新

#### Python版本
```python
# 不一致时的日志输出
logging.warning(f"数据不一致 - Azure实例: {azure_instance['name']} 数据库: {azure_instance['database']} 表: {table} vs AWS实例: {aws_instance['name']} 数据库: {aws_instance['database']} 表: {table}")

# 一致时的日志输出
logging.debug(f"数据一致 - Azure实例: {azure_instance['name']} 数据库: {azure_instance['database']} 表: {table} vs AWS实例: {aws_instance['name']} 数据库: {aws_instance['database']} 表: {table}")
```

#### Go版本
```go
// 不一致时的日志输出
log.Printf("数据不一致 - Azure实例: %s 数据库: %s 表: %s vs AWS实例: %s 数据库: %s 表: %s", 
    pair.AzureInstance.Name, pair.AzureInstance.Database, table,
    pair.AWSInstance.Name, pair.AWSInstance.Database, table)

// 一致时的日志输出
log.Printf("数据一致 - Azure实例: %s 数据库: %s 表: %s vs AWS实例: %s 数据库: %s 表: %s", 
    pair.AzureInstance.Name, pair.AzureInstance.Database, table,
    pair.AWSInstance.Name, pair.AWSInstance.Database, table)
```

### 3. 最终输出格式更新

#### Python版本输出格式
```
详细信息:

数据库: db1
Azure实例: azure-db1
AWS实例: aws-db1
状态: INCONSISTENT
不一致的表:
  - 表名: orders
    Azure实例: azure-db1 数据库: db1
    AWS实例: aws-db1 数据库: db1
    Azure校验和: xyz789uvw012
    AWS校验和: different_checksum_here
```

#### Go版本输出格式
```
详细信息:

数据库: db1
Azure实例: azure-db1
AWS实例: aws-db1
状态: INCONSISTENT
表对比结果:
  ✓ users
  ✗ orders
    Azure实例: azure-db1 数据库: db1
    AWS实例: aws-db1 数据库: db1
    Azure校验和: xyz789uvw012
    AWS校验和: different_checksum_here
  ✓ products
```

## 输出格式对比

| 特性 | Python版本 | Go版本 |
|------|------------|--------|
| 不一致表显示 | 单独列出不一致的表 | 显示所有表，用✓/✗标记 |
| 实例信息 | 包含实例名和数据库名 | 包含实例名和数据库名 |
| 校验和显示 | 显示Azure和AWS校验和 | 显示Azure和AWS校验和 |
| 格式风格 | 简洁列表式 | 树状结构式 |

## 使用示例

### 实际运行输出示例

当发现数据不一致时，两个版本都会输出详细的实例和表信息：

**Python版本**：
```
2025-10-08 14:11:21,535 - WARNING - 数据不一致 - Azure实例: azure-db1 数据库: db1 表: orders vs AWS实例: aws-db1 数据库: db1 表: orders

详细信息:
数据库: db1
Azure实例: azure-db1
AWS实例: aws-db1
状态: INCONSISTENT
不一致的表:
  - 表名: orders
    Azure实例: azure-db1 数据库: db1
    AWS实例: aws-db1 数据库: db1
    Azure校验和: xyz789uvw012
    AWS校验和: different_checksum_here
```

**Go版本**：
```
2025/10/08 14:11:21 数据不一致 - Azure实例: azure-db1 数据库: db1 表: orders vs AWS实例: aws-db1 数据库: db1 表: orders

详细信息:
数据库: db1
Azure实例: azure-db1
AWS实例: aws-db1
状态: INCONSISTENT
表对比结果:
  ✓ users
  ✗ orders
    Azure实例: azure-db1 数据库: db1
    AWS实例: aws-db1 数据库: db1
    Azure校验和: xyz789uvw012
    AWS校验和: different_checksum_here
  ✓ products
```

## 优势

1. **精确定位**: 可以准确知道哪个实例的哪个数据库的哪个表出现了不一致
2. **便于排查**: 提供了完整的上下文信息，便于DBA快速定位问题
3. **详细对比**: 显示具体的校验和值，便于进一步分析差异
4. **统一格式**: 两个版本都提供了相同级别的详细信息

## 兼容性

- 所有现有功能保持不变
- 只是增强了输出信息的详细程度
- 向后兼容，不影响现有的报告格式
- JSON报告文件中也包含了新的字段信息

## 总结

这次更新显著提升了数据不一致时的可观测性，让用户能够快速准确地定位到具体的问题表，大大提高了数据库迁移验证的效率。两个版本都提供了相同级别的详细信息，用户可以根据需要选择合适的版本使用。
