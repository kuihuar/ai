# MySQL从Azure迁移到AWS数据一致性对比指南

## 背景

MySQL数据库从Azure Database for MySQL迁移到AWS RDS for MySQL时，确保数据一致性是迁移成功的关键。本文档总结了数据一致性验证的方法、工具和最佳实践。

## 迁移前准备

### 1. 环境评估

**Azure环境信息收集：**
```bash
# 获取数据库基本信息
SELECT VERSION(), @@hostname, @@port;
SELECT DATABASE(), USER();

# 获取数据库大小
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables 
GROUP BY table_schema;

# 获取表数量统计
SELECT 
    table_schema,
    COUNT(*) as table_count
FROM information_schema.tables 
WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
GROUP BY table_schema;
```

**AWS环境准备：**
- 创建RDS实例，确保版本兼容性
- 配置安全组和子网
- 设置参数组（如果需要特殊配置）

### 2. 一致性检查策略

#### 策略1：基于校验和的对比
```sql
-- 为每个表生成校验和
SELECT 
    table_name,
    COUNT(*) as row_count,
    CHECKSUM(*) as table_checksum
FROM your_table_name
GROUP BY table_name;
```

#### 策略2：基于行数的对比
```sql
-- 统计每个表的行数
SELECT 
    table_schema,
    table_name,
    table_rows
FROM information_schema.tables 
WHERE table_schema = 'your_database_name'
ORDER BY table_name;
```

#### 策略3：基于关键字段的对比
```sql
-- 对比关键业务表的主键范围
SELECT 
    MIN(id) as min_id,
    MAX(id) as max_id,
    COUNT(*) as total_count
FROM your_important_table;
```

## 数据一致性验证工具

### 1. 自建对比脚本

**Python脚本示例：**
```python
import pymysql
import hashlib
import json
from datetime import datetime

class DatabaseComparator:
    def __init__(self, azure_config, aws_config):
        self.azure_conn = pymysql.connect(**azure_config)
        self.aws_conn = pymysql.connect(**aws_config)
    
    def get_table_list(self, connection):
        """获取数据库表列表"""
        cursor = connection.cursor()
        cursor.execute("""
            SELECT table_name 
            FROM information_schema.tables 
            WHERE table_schema = DATABASE()
            AND table_type = 'BASE TABLE'
        """)
        return [row[0] for row in cursor.fetchall()]
    
    def calculate_table_checksum(self, connection, table_name):
        """计算表的校验和"""
        cursor = connection.cursor()
        cursor.execute(f"SELECT COUNT(*) FROM {table_name}")
        row_count = cursor.fetchone()[0]
        
        # 对于大表，可以分批计算校验和
        if row_count > 1000000:
            return self.calculate_large_table_checksum(connection, table_name)
        
        cursor.execute(f"SELECT * FROM {table_name} ORDER BY 1")
        data = cursor.fetchall()
        data_str = str(data)
        return hashlib.md5(data_str.encode()).hexdigest()
    
    def compare_databases(self):
        """对比两个数据库"""
        azure_tables = self.get_table_list(self.azure_conn)
        aws_tables = self.get_table_list(self.aws_conn)
        
        results = {
            'timestamp': datetime.now().isoformat(),
            'azure_tables': len(azure_tables),
            'aws_tables': len(aws_tables),
            'comparisons': []
        }
        
        for table in azure_tables:
            if table in aws_tables:
                azure_checksum = self.calculate_table_checksum(self.azure_conn, table)
                aws_checksum = self.calculate_table_checksum(self.aws_conn, table)
                
                results['comparisons'].append({
                    'table': table,
                    'azure_checksum': azure_checksum,
                    'aws_checksum': aws_checksum,
                    'match': azure_checksum == aws_checksum
                })
        
        return results

# 使用示例
azure_config = {
    'host': 'your-azure-mysql.mysql.database.azure.com',
    'user': 'your_user',
    'password': 'your_password',
    'database': 'your_database'
}

aws_config = {
    'host': 'your-aws-rds.region.rds.amazonaws.com',
    'user': 'your_user',
    'password': 'your_password',
    'database': 'your_database'
}

comparator = DatabaseComparator(azure_config, aws_config)
results = comparator.compare_databases()
print(json.dumps(results, indent=2))
```

### 2. 使用专业工具

**pt-table-checksum (Percona Toolkit):**
```bash
# 安装Percona Toolkit
# Ubuntu/Debian
sudo apt-get install percona-toolkit

# 对比源数据库和目标数据库
pt-table-checksum \
  --host=azure-mysql-host \
  --user=username \
  --password=password \
  --databases=your_database \
  --replicate=percona.checksums \
  --no-check-binlog-format

# 在目标数据库上检查差异
pt-table-sync \
  --host=aws-rds-host \
  --user=username \
  --password=password \
  --databases=your_database \
  --replicate=percona.checksums \
  --print
```

**mysqldbcompare (MySQL Utilities):**
```bash
# 安装MySQL Utilities
pip install mysql-utilities

# 对比数据库结构
mysqldbcompare \
  --server1=user:pass@azure-host:3306 \
  --server2=user:pass@aws-host:3306 \
  --difftype=sql \
  database1:database2

# 对比数据
mysqldbcompare \
  --server1=user:pass@azure-host:3306 \
  --server2=user:pass@aws-host:3306 \
  --difftype=sql \
  --skip-table-options \
  database1:database2
```

## 迁移过程中的一致性保证

### 1. 双写策略

在迁移期间，可以实施双写策略确保数据一致性：

```python
class DualWriteManager:
    def __init__(self, primary_conn, secondary_conn):
        self.primary = primary_conn
        self.secondary = secondary_conn
    
    def execute_with_dual_write(self, sql, params=None):
        """执行SQL并同时写入两个数据库"""
        try:
            # 写入主数据库
            cursor1 = self.primary.cursor()
            cursor1.execute(sql, params)
            self.primary.commit()
            
            # 写入次数据库
            cursor2 = self.secondary.cursor()
            cursor2.execute(sql, params)
            self.secondary.commit()
            
            return True
        except Exception as e:
            # 如果次数据库写入失败，记录日志但不影响主数据库
            print(f"Secondary write failed: {e}")
            return False
```

### 2. 实时同步监控

```python
class SyncMonitor:
    def __init__(self, source_conn, target_conn):
        self.source = source_conn
        self.target = target_conn
    
    def monitor_lag(self):
        """监控同步延迟"""
        # 检查binlog位置
        source_cursor = self.source.cursor()
        source_cursor.execute("SHOW MASTER STATUS")
        source_status = source_cursor.fetchone()
        
        target_cursor = self.target.cursor()
        target_cursor.execute("SHOW SLAVE STATUS")
        target_status = target_cursor.fetchone()
        
        if target_status:
            lag = target_status[32]  # Seconds_Behind_Master
            return lag
        return None
```

## 常见问题和解决方案

### 1. 字符集和排序规则差异

**问题：** Azure和AWS的默认字符集可能不同

**解决方案：**
```sql
-- 检查字符集
SHOW VARIABLES LIKE 'character_set%';
SHOW VARIABLES LIKE 'collation%';

-- 统一字符集
ALTER DATABASE your_database CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 时区设置差异

**问题：** 不同云平台的时区设置可能影响时间戳数据

**解决方案：**
```sql
-- 检查时区设置
SELECT @@global.time_zone, @@session.time_zone;

-- 设置统一时区
SET GLOBAL time_zone = '+00:00';
SET SESSION time_zone = '+00:00';
```

### 3. 大表迁移策略

**问题：** 大表迁移时间长，容易产生不一致

**解决方案：**
```bash
# 使用mysqldump进行大表迁移
mysqldump \
  --single-transaction \
  --routines \
  --triggers \
  --lock-tables=false \
  --host=azure-host \
  --user=username \
  --password=password \
  your_database your_large_table | \
mysql \
  --host=aws-host \
  --user=username \
  --password=password \
  your_database
```

## 验证清单

### 迁移前检查
- [ ] 数据库版本兼容性确认
- [ ] 字符集和排序规则统一
- [ ] 时区设置统一
- [ ] 网络连通性测试
- [ ] 权限配置确认

### 迁移后验证
- [ ] 表数量对比
- [ ] 行数统计对比
- [ ] 关键表数据校验和对比
- [ ] 索引完整性检查
- [ ] 外键约束验证
- [ ] 存储过程和触发器验证
- [ ] 用户权限对比

### 业务验证
- [ ] 关键业务流程测试
- [ ] 性能基准测试
- [ ] 应用程序连接测试
- [ ] 备份恢复测试

## 最佳实践总结

1. **分阶段迁移**：先迁移非关键数据，验证流程后再迁移核心数据
2. **并行验证**：在迁移过程中持续进行数据一致性检查
3. **回滚准备**：制定详细的回滚计划，确保可以快速恢复
4. **监控告警**：设置实时监控，及时发现数据不一致问题
5. **文档记录**：详细记录迁移过程和验证结果，便于后续参考

## 工具推荐

- **pt-table-checksum**: 专业的数据一致性检查工具
- **mysqldbcompare**: MySQL官方工具，支持结构和数据对比
- **Percona XtraBackup**: 物理备份工具，适合大数据库迁移
- **AWS DMS**: AWS数据迁移服务，支持实时同步
- **自建脚本**: 根据具体需求定制化验证逻辑

## 总结

MySQL从Azure迁移到AWS的数据一致性验证需要综合考虑多个方面，包括技术工具、验证策略、监控机制等。通过系统性的方法和工具，可以确保迁移过程中数据的完整性和一致性，降低业务风险。
