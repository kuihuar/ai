# 多数据库一致性验证方案

## 背景

针对MySQL从Azure迁移到AWS的多个数据库一致性验证需求，本文档提供了几种可行的验证方案，适用于DBA团队执行。

## 方案对比

| 方案 | 适用场景 | 优点 | 缺点 | 推荐度 |
|------|----------|------|------|--------|
| pt-table-checksum | 生产环境，大数据库 | 专业、稳定、支持大表 | 需要安装工具 | ⭐⭐⭐⭐⭐ |
| mysqldbcompare | 中小型数据库 | 官方工具、易用 | 性能有限 | ⭐⭐⭐⭐ |
| 自建脚本 | 定制化需求 | 灵活、可控 | 开发工作量大 | ⭐⭐⭐ |
| AWS DMS验证 | 实时验证 | 自动化、实时 | 成本较高 | ⭐⭐⭐⭐ |

## 方案一：pt-table-checksum（推荐）

### 适用场景
- 生产环境
- 大数据库（TB级别）
- 需要高精度验证
- 有Percona Toolkit使用经验

### 实施步骤

#### 1. 环境准备
```bash
# 在验证服务器上安装Percona Toolkit
# CentOS/RHEL
sudo yum install percona-toolkit

# Ubuntu/Debian
sudo apt-get install percona-toolkit

# 验证安装
pt-table-checksum --version
```

#### 2. 批量验证脚本
```bash
#!/bin/bash
# multi-db-consistency-check.sh

# 配置变量
AZURE_HOST="your-azure-mysql.mysql.database.azure.com"
AWS_HOST="your-aws-rds.region.rds.amazonaws.com"
USERNAME="your_username"
PASSWORD="your_password"
DATABASES=("db1" "db2" "db3" "db4" "db5")  # 需要验证的数据库列表

# 创建结果目录
mkdir -p consistency_results/$(date +%Y%m%d_%H%M%S)
RESULT_DIR="consistency_results/$(date +%Y%m%d_%H%M%S)"

echo "开始多数据库一致性验证..."
echo "验证时间: $(date)"
echo "结果目录: $RESULT_DIR"

# 循环验证每个数据库
for db in "${DATABASES[@]}"; do
    echo "正在验证数据库: $db"
    
    # 执行pt-table-checksum
    pt-table-checksum \
        --host=$AZURE_HOST \
        --user=$USERNAME \
        --password=$PASSWORD \
        --databases=$db \
        --replicate=percona.checksums \
        --no-check-binlog-format \
        --chunk-size=1000 \
        --chunk-time=0.5 \
        --max-lag=1s \
        --check-slave-lag \
        --recursion-method=hosts \
        --replicate-check-only \
        --output=$RESULT_DIR/${db}_azure_checksum.txt
    
    # 在AWS端检查差异
    pt-table-sync \
        --host=$AWS_HOST \
        --user=$USERNAME \
        --password=$PASSWORD \
        --databases=$db \
        --replicate=percona.checksums \
        --print \
        --output=$RESULT_DIR/${db}_aws_diff.txt
    
    # 生成验证报告
    echo "数据库 $db 验证完成" >> $RESULT_DIR/summary.txt
    echo "Azure校验和文件: ${db}_azure_checksum.txt" >> $RESULT_DIR/summary.txt
    echo "AWS差异文件: ${db}_aws_diff.txt" >> $RESULT_DIR/summary.txt
    echo "---" >> $RESULT_DIR/summary.txt
done

echo "所有数据库验证完成，结果保存在: $RESULT_DIR"
```

#### 3. 结果分析脚本
```python
#!/usr/bin/env python3
# analyze_consistency_results.py

import os
import re
from datetime import datetime

def analyze_checksum_results(result_dir):
    """分析pt-table-checksum结果"""
    summary = {
        'total_databases': 0,
        'consistent_databases': 0,
        'inconsistent_databases': 0,
        'details': []
    }
    
    for filename in os.listdir(result_dir):
        if filename.endswith('_aws_diff.txt'):
            db_name = filename.replace('_aws_diff.txt', '')
            summary['total_databases'] += 1
            
            filepath = os.path.join(result_dir, filename)
            with open(filepath, 'r') as f:
                content = f.read()
                
            # 检查是否有差异
            if 'DIFFS' in content or 'REPLACE' in content:
                summary['inconsistent_databases'] += 1
                summary['details'].append({
                    'database': db_name,
                    'status': 'INCONSISTENT',
                    'issues': extract_issues(content)
                })
            else:
                summary['consistent_databases'] += 1
                summary['details'].append({
                    'database': db_name,
                    'status': 'CONSISTENT',
                    'issues': []
                })
    
    return summary

def extract_issues(content):
    """提取具体的不一致问题"""
    issues = []
    lines = content.split('\n')
    
    for line in lines:
        if 'REPLACE' in line or 'DIFFS' in line:
            issues.append(line.strip())
    
    return issues

def generate_report(summary, result_dir):
    """生成验证报告"""
    report_path = os.path.join(result_dir, 'consistency_report.html')
    
    html_content = f"""
    <!DOCTYPE html>
    <html>
    <head>
        <title>数据库一致性验证报告</title>
        <style>
            body {{ font-family: Arial, sans-serif; margin: 20px; }}
            .header {{ background-color: #f0f0f0; padding: 20px; border-radius: 5px; }}
            .summary {{ background-color: #e8f4f8; padding: 15px; margin: 20px 0; border-radius: 5px; }}
            .consistent {{ color: green; }}
            .inconsistent {{ color: red; }}
            .details {{ margin: 20px 0; }}
            table {{ border-collapse: collapse; width: 100%; }}
            th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
            th {{ background-color: #f2f2f2; }}
        </style>
    </head>
    <body>
        <div class="header">
            <h1>数据库一致性验证报告</h1>
            <p>生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
        </div>
        
        <div class="summary">
            <h2>验证摘要</h2>
            <p>总数据库数: {summary['total_databases']}</p>
            <p class="consistent">一致数据库数: {summary['consistent_databases']}</p>
            <p class="inconsistent">不一致数据库数: {summary['inconsistent_databases']}</p>
        </div>
        
        <div class="details">
            <h2>详细结果</h2>
            <table>
                <tr>
                    <th>数据库</th>
                    <th>状态</th>
                    <th>问题描述</th>
                </tr>
    """
    
    for detail in summary['details']:
        status_class = 'consistent' if detail['status'] == 'CONSISTENT' else 'inconsistent'
        issues_text = '<br>'.join(detail['issues']) if detail['issues'] else '无'
        
        html_content += f"""
                <tr>
                    <td>{detail['database']}</td>
                    <td class="{status_class}">{detail['status']}</td>
                    <td>{issues_text}</td>
                </tr>
        """
    
    html_content += """
            </table>
        </div>
    </body>
    </html>
    """
    
    with open(report_path, 'w', encoding='utf-8') as f:
        f.write(html_content)
    
    print(f"验证报告已生成: {report_path}")

if __name__ == "__main__":
    import sys
    if len(sys.argv) != 2:
        print("用法: python analyze_consistency_results.py <结果目录>")
        sys.exit(1)
    
    result_dir = sys.argv[1]
    summary = analyze_checksum_results(result_dir)
    generate_report(summary, result_dir)
    
    print(f"验证完成:")
    print(f"总数据库数: {summary['total_databases']}")
    print(f"一致数据库数: {summary['consistent_databases']}")
    print(f"不一致数据库数: {summary['inconsistent_databases']}")
```

## 方案二：mysqldbcompare（适合中小型数据库）

### 适用场景
- 中小型数据库（GB级别）
- 需要结构对比
- 团队对MySQL Utilities熟悉

### 批量验证脚本
```bash
#!/bin/bash
# mysql-utilities-batch-check.sh

AZURE_HOST="your-azure-mysql.mysql.database.azure.com"
AWS_HOST="your-aws-rds.region.rds.amazonaws.com"
USERNAME="your_username"
PASSWORD="your_password"
DATABASES=("db1" "db2" "db3" "db4" "db5")

RESULT_DIR="mysql_utilities_results/$(date +%Y%m%d_%H%M%S)"
mkdir -p $RESULT_DIR

echo "开始使用MySQL Utilities验证..."

for db in "${DATABASES[@]}"; do
    echo "验证数据库: $db"
    
    # 结构对比
    mysqldbcompare \
        --server1=$USERNAME:$PASSWORD@$AZURE_HOST:3306 \
        --server2=$USERNAME:$PASSWORD@$AWS_HOST:3306 \
        --difftype=sql \
        --skip-table-options \
        $db:$db \
        > $RESULT_DIR/${db}_structure_diff.txt
    
    # 数据对比
    mysqldbcompare \
        --server1=$USERNAME:$PASSWORD@$AZURE_HOST:3306 \
        --server2=$USERNAME:$PASSWORD@$AWS_HOST:3306 \
        --difftype=sql \
        --skip-table-options \
        --skip-data-check=false \
        $db:$db \
        > $RESULT_DIR/${db}_data_diff.txt
    
    echo "数据库 $db 验证完成"
done

echo "所有验证完成，结果保存在: $RESULT_DIR"
```

## 方案三：自建Python验证脚本

### 适用场景
- 需要定制化验证逻辑
- 有Python开发能力
- 需要集成到现有系统

### 多数据库验证脚本（Python版本）

```python
#!/usr/bin/env python3
# multi_database_validator.py
# 多数据库一致性验证工具 - Python版本
# 功能：验证MySQL数据库从Azure迁移到AWS后的数据一致性

import pymysql
import hashlib
import json
import concurrent.futures
from datetime import datetime
import logging
import sys
import os

# 配置日志系统
# 同时输出到文件和控制台，便于调试和监控
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('consistency_check.log', encoding='utf-8'),
        logging.StreamHandler(sys.stdout)
    ]
)

class MultiDatabaseValidator:
    """
    多数据库一致性验证器
    
    主要功能：
    1. 连接Azure和AWS的MySQL数据库
    2. 获取数据库表列表
    3. 计算每个表的数据校验和
    4. 对比两个环境的数据一致性
    5. 生成详细的验证报告
    """
    
    def __init__(self, azure_config, aws_config, databases):
        """
        初始化验证器
        
        Args:
            azure_config (dict): Azure数据库连接配置
            aws_config (dict): AWS数据库连接配置
            databases (list): 需要验证的数据库名称列表
        """
        self.azure_config = azure_config
        self.aws_config = aws_config
        self.databases = databases
        self.results = {}  # 存储验证结果
        
        # 验证配置参数
        self._validate_config()
    
    def _validate_config(self):
        """验证配置参数的有效性"""
        required_keys = ['host', 'user', 'password']
        
        for config_name, config in [('Azure', self.azure_config), ('AWS', self.aws_config)]:
            for key in required_keys:
                if key not in config:
                    raise ValueError(f"{config_name}配置缺少必要参数: {key}")
        
        if not self.databases:
            raise ValueError("数据库列表不能为空")
    
    def get_table_list(self, connection, database):
        """
        获取指定数据库的表列表
        
        Args:
            connection: 数据库连接对象
            database (str): 数据库名称
            
        Returns:
            list: 表名称列表
        """
        cursor = connection.cursor()
        try:
            # 查询指定数据库中的所有用户表（排除系统表）
            cursor.execute("""
                SELECT table_name 
                FROM information_schema.tables 
                WHERE table_schema = %s
                AND table_type = 'BASE TABLE'
                ORDER BY table_name
            """, (database,))
            
            tables = [row[0] for row in cursor.fetchall()]
            logging.debug(f"数据库 {database} 包含 {len(tables)} 个表")
            return tables
            
        except Exception as e:
            logging.error(f"获取数据库 {database} 表列表失败: {str(e)}")
            raise
        finally:
            cursor.close()
    
    def calculate_table_checksum(self, connection, database, table_name):
        """
        计算表的校验和
        
        策略：
        1. 空表返回特殊标识
        2. 小表（<10万行）直接计算
        3. 大表（>=10万行）分批计算
        
        Args:
            connection: 数据库连接对象
            database (str): 数据库名称
            table_name (str): 表名称
            
        Returns:
            str: 表的MD5校验和
        """
        cursor = connection.cursor()
        try:
            # 获取表的行数
            cursor.execute(f"SELECT COUNT(*) FROM `{database}`.`{table_name}`")
            row_count = cursor.fetchone()[0]
            
            # 空表处理
            if row_count == 0:
                logging.debug(f"表 {database}.{table_name} 为空表")
                return "empty_table"
            
            # 大表分批处理，避免内存溢出
            if row_count > 100000:
                logging.info(f"表 {database}.{table_name} 行数较多({row_count})，使用分批计算")
                return self.calculate_large_table_checksum(connection, database, table_name)
            
            # 小表直接计算
            logging.debug(f"计算表 {database}.{table_name} 的校验和，行数: {row_count}")
            cursor.execute(f"SELECT * FROM `{database}`.`{table_name}` ORDER BY 1")
            data = cursor.fetchall()
            data_str = str(data)
            return hashlib.md5(data_str.encode('utf-8')).hexdigest()
            
        except Exception as e:
            logging.error(f"计算表 {database}.{table_name} 校验和失败: {str(e)}")
            raise
        finally:
            cursor.close()
    
    def calculate_large_table_checksum(self, connection, database, table_name):
        """
        大表分批计算校验和
        
        对于大表，为了避免内存溢出，采用分批读取的方式：
        1. 每次读取10000行数据
        2. 计算每批数据的校验和
        3. 最后合并所有批次的校验和
        
        Args:
            connection: 数据库连接对象
            database (str): 数据库名称
            table_name (str): 表名称
            
        Returns:
            str: 合并后的MD5校验和
        """
        cursor = connection.cursor()
        try:
            # 获取总行数
            cursor.execute(f"SELECT COUNT(*) FROM `{database}`.`{table_name}`")
            total_rows = cursor.fetchone()[0]
            
            batch_size = 10000  # 每批处理10000行
            checksums = []
            
            logging.info(f"开始分批计算表 {database}.{table_name}，总行数: {total_rows}")
            
            # 分批读取数据
            for offset in range(0, total_rows, batch_size):
                logging.debug(f"处理批次: {offset}-{min(offset + batch_size, total_rows)}")
                
                cursor.execute(f"""
                    SELECT * FROM `{database}`.`{table_name}` 
                    ORDER BY 1 
                    LIMIT {batch_size} OFFSET {offset}
                """)
                
                batch_data = cursor.fetchall()
                batch_str = str(batch_data)
                batch_checksum = hashlib.md5(batch_str.encode('utf-8')).hexdigest()
                checksums.append(batch_checksum)
            
            # 合并所有批次的校验和
            combined_checksum = hashlib.md5(''.join(checksums).encode('utf-8')).hexdigest()
            logging.info(f"表 {database}.{table_name} 分批计算完成，共 {len(checksums)} 个批次")
            
            return combined_checksum
            
        except Exception as e:
            logging.error(f"分批计算表 {database}.{table_name} 校验和失败: {str(e)}")
            raise
        finally:
            cursor.close()
    
    def validate_database(self, database):
        """
        验证单个数据库的一致性
        
        验证步骤：
        1. 连接Azure和AWS数据库
        2. 获取表列表
        3. 对比表数量
        4. 逐个验证表的数据一致性
        5. 记录验证结果
        
        Args:
            database (str): 数据库名称
            
        Returns:
            dict: 验证结果
        """
        logging.info(f"开始验证数据库: {database}")
        
        azure_conn = None
        aws_conn = None
        
        try:
            # 连接数据库
            logging.debug(f"连接Azure数据库: {self.azure_config['host']}")
            azure_conn = pymysql.connect(**self.azure_config)
            
            logging.debug(f"连接AWS数据库: {self.aws_config['host']}")
            aws_conn = pymysql.connect(**self.aws_config)
            
            # 获取表列表
            azure_tables = self.get_table_list(azure_conn, database)
            aws_tables = self.get_table_list(aws_conn, database)
            
            # 初始化验证结果
            result = {
                'database': database,
                'azure_tables': len(azure_tables),
                'aws_tables': len(aws_tables),
                'table_comparisons': [],
                'status': 'SUCCESS',
                'errors': [],
                'start_time': datetime.now().isoformat()
            }
            
            # 检查表数量一致性
            if len(azure_tables) != len(aws_tables):
                result['status'] = 'WARNING'
                error_msg = f"表数量不一致: Azure({len(azure_tables)}) vs AWS({len(aws_tables)})"
                result['errors'].append(error_msg)
                logging.warning(f"数据库 {database}: {error_msg}")
            
            # 对比每个表的数据一致性
            logging.info(f"开始验证数据库 {database} 中的 {len(azure_tables)} 个表")
            
            for i, table in enumerate(azure_tables, 1):
                logging.info(f"验证表 {i}/{len(azure_tables)}: {table}")
                
                if table in aws_tables:
                    try:
                        # 计算Azure表的校验和
                        azure_checksum = self.calculate_table_checksum(azure_conn, database, table)
                        
                        # 计算AWS表的校验和
                        aws_checksum = self.calculate_table_checksum(aws_conn, database, table)
                        
                        # 记录对比结果
                        table_result = {
                            'table': table,
                            'azure_checksum': azure_checksum,
                            'aws_checksum': aws_checksum,
                            'match': azure_checksum == aws_checksum
                        }
                        
                        result['table_comparisons'].append(table_result)
                        
                        # 检查是否一致
                        if azure_checksum != aws_checksum:
                            result['status'] = 'INCONSISTENT'
                            logging.warning(f"表 {database}.{table} 数据不一致")
                        else:
                            logging.debug(f"表 {database}.{table} 数据一致")
                            
                    except Exception as e:
                        error_msg = f"表 {table} 验证失败: {str(e)}"
                        result['errors'].append(error_msg)
                        result['status'] = 'ERROR'
                        logging.error(f"数据库 {database}: {error_msg}")
                else:
                    error_msg = f"表 {table} 在AWS中不存在"
                    result['errors'].append(error_msg)
                    result['status'] = 'INCONSISTENT'
                    logging.warning(f"数据库 {database}: {error_msg}")
            
            # 记录结束时间
            result['end_time'] = datetime.now().isoformat()
            
            # 统计验证结果
            consistent_tables = sum(1 for t in result['table_comparisons'] if t['match'])
            total_tables = len(result['table_comparisons'])
            
            logging.info(f"数据库 {database} 验证完成:")
            logging.info(f"  状态: {result['status']}")
            logging.info(f"  表数量: Azure({result['azure_tables']}) vs AWS({result['aws_tables']})")
            logging.info(f"  数据一致性: {consistent_tables}/{total_tables}")
            logging.info(f"  错误数量: {len(result['errors'])}")
            
            return result
            
        except Exception as e:
            error_msg = f"数据库 {database} 验证失败: {str(e)}"
            logging.error(error_msg)
            return {
                'database': database,
                'status': 'ERROR',
                'errors': [error_msg],
                'start_time': datetime.now().isoformat(),
                'end_time': datetime.now().isoformat()
            }
        finally:
            # 确保连接被正确关闭
            if azure_conn:
                azure_conn.close()
            if aws_conn:
                aws_conn.close()
    
    def validate_all_databases(self, max_workers=3):
        """
        并行验证所有数据库
        
        使用线程池并行处理多个数据库，提高验证效率
        
        Args:
            max_workers (int): 最大并发线程数，默认3个
        """
        logging.info(f"开始验证 {len(self.databases)} 个数据库，最大并发数: {max_workers}")
        
        # 使用线程池并行处理
        with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
            # 提交所有验证任务
            future_to_db = {
                executor.submit(self.validate_database, db): db 
                for db in self.databases
            }
            
            # 收集验证结果
            for future in concurrent.futures.as_completed(future_to_db):
                db = future_to_db[future]
                try:
                    result = future.result()
                    self.results[db] = result
                    logging.info(f"数据库 {db} 验证完成，状态: {result['status']}")
                except Exception as e:
                    error_msg = f"数据库 {db} 验证异常: {str(e)}"
                    logging.error(error_msg)
                    self.results[db] = {
                        'database': db,
                        'status': 'ERROR',
                        'errors': [error_msg],
                        'start_time': datetime.now().isoformat(),
                        'end_time': datetime.now().isoformat()
                    }
        
        logging.info("所有数据库验证完成")
    
    def generate_report(self, output_file='consistency_report.json'):
        """
        生成验证报告
        
        生成包含详细验证结果的JSON报告，便于后续分析和存档
        
        Args:
            output_file (str): 输出文件名
            
        Returns:
            dict: 验证摘要
        """
        # 统计验证结果
        total_databases = len(self.databases)
        successful_validations = sum(1 for r in self.results.values() if r['status'] == 'SUCCESS')
        inconsistent_databases = sum(1 for r in self.results.values() if r['status'] == 'INCONSISTENT')
        error_databases = sum(1 for r in self.results.values() if r['status'] == 'ERROR')
        
        # 生成报告摘要
        summary = {
            'timestamp': datetime.now().isoformat(),
            'total_databases': total_databases,
            'successful_validations': successful_validations,
            'inconsistent_databases': inconsistent_databases,
            'error_databases': error_databases,
            'success_rate': f"{(successful_validations/total_databases*100):.2f}%" if total_databases > 0 else "0%",
            'results': self.results
        }
        
        # 保存报告到文件
        try:
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(summary, f, indent=2, ensure_ascii=False)
            logging.info(f"验证报告已生成: {output_file}")
        except Exception as e:
            logging.error(f"生成报告失败: {str(e)}")
            raise
        
        return summary

def main():
    """主函数 - 程序入口点"""
    try:
        # 配置数据库连接参数
        azure_config = {
            'host': 'your-azure-mysql.mysql.database.azure.com',
            'user': 'your_username',
            'password': 'your_password',
            'charset': 'utf8mb4',
            'connect_timeout': 30,
            'read_timeout': 30,
            'write_timeout': 30
        }
        
        aws_config = {
            'host': 'your-aws-rds.region.rds.amazonaws.com',
            'user': 'your_username',
            'password': 'your_password',
            'charset': 'utf8mb4',
            'connect_timeout': 30,
            'read_timeout': 30,
            'write_timeout': 30
        }
        
        # 需要验证的数据库列表
        databases = ['db1', 'db2', 'db3', 'db4', 'db5']
        
        # 创建验证器实例
        validator = MultiDatabaseValidator(azure_config, aws_config, databases)
        
        # 执行验证
        validator.validate_all_databases(max_workers=3)
        
        # 生成报告
        summary = validator.generate_report()
        
        # 输出验证摘要
        print("\n" + "="*60)
        print("数据库一致性验证完成")
        print("="*60)
        print(f"验证时间: {summary['timestamp']}")
        print(f"总数据库数: {summary['total_databases']}")
        print(f"验证成功: {summary['successful_validations']}")
        print(f"数据不一致: {summary['inconsistent_databases']}")
        print(f"验证错误: {summary['error_databases']}")
        print(f"成功率: {summary['success_rate']}")
        print("="*60)
        
        # 如果有不一致或错误，显示详细信息
        if summary['inconsistent_databases'] > 0 or summary['error_databases'] > 0:
            print("\n详细信息:")
            for db_name, result in summary['results'].items():
                if result['status'] in ['INCONSISTENT', 'ERROR']:
                    print(f"\n数据库: {db_name}")
                    print(f"状态: {result['status']}")
                    if result['errors']:
                        print("错误信息:")
                        for error in result['errors']:
                            print(f"  - {error}")
        
        return 0 if summary['error_databases'] == 0 else 1
        
    except Exception as e:
        logging.error(f"程序执行失败: {str(e)}")
        return 1

if __name__ == "__main__":
    exit_code = main()
    sys.exit(exit_code)
```

### 多数据库验证脚本（Go语言版本）

```go
// multi_database_validator.go
// 多数据库一致性验证工具 - Go语言版本
// 功能：验证MySQL数据库从Azure迁移到AWS后的数据一致性

package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

// TableComparison 表对比结果
type TableComparison struct {
	Table           string `json:"table"`
	AzureChecksum   string `json:"azure_checksum"`
	AWSChecksum     string `json:"aws_checksum"`
	Match           bool   `json:"match"`
}

// DatabaseResult 数据库验证结果
type DatabaseResult struct {
	Database         string             `json:"database"`
	AzureTables      int                `json:"azure_tables"`
	AWSTables        int                `json:"aws_tables"`
	TableComparisons []TableComparison  `json:"table_comparisons"`
	Status           string             `json:"status"`
	Errors           []string           `json:"errors"`
	StartTime        string             `json:"start_time"`
	EndTime          string             `json:"end_time"`
}

// ValidationSummary 验证摘要
type ValidationSummary struct {
	Timestamp              string                    `json:"timestamp"`
	TotalDatabases         int                       `json:"total_databases"`
	SuccessfulValidations  int                       `json:"successful_validations"`
	InconsistentDatabases  int                       `json:"inconsistent_databases"`
	ErrorDatabases         int                       `json:"error_databases"`
	SuccessRate            string                    `json:"success_rate"`
	Results                map[string]DatabaseResult `json:"results"`
}

// MultiDatabaseValidator 多数据库验证器
type MultiDatabaseValidator struct {
	azureConfig DatabaseConfig
	awsConfig   DatabaseConfig
	databases   []string
	results     map[string]DatabaseResult
	mutex       sync.RWMutex
}

// NewMultiDatabaseValidator 创建新的验证器实例
func NewMultiDatabaseValidator(azureConfig, awsConfig DatabaseConfig, databases []string) *MultiDatabaseValidator {
	// 验证配置参数
	if err := validateConfig(azureConfig, awsConfig, databases); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	return &MultiDatabaseValidator{
		azureConfig: azureConfig,
		awsConfig:   awsConfig,
		databases:   databases,
		results:     make(map[string]DatabaseResult),
	}
}

// validateConfig 验证配置参数的有效性
func validateConfig(azureConfig, awsConfig DatabaseConfig, databases []string) error {
	// 检查Azure配置
	if azureConfig.Host == "" || azureConfig.User == "" || azureConfig.Password == "" {
		return fmt.Errorf("Azure配置缺少必要参数")
	}

	// 检查AWS配置
	if awsConfig.Host == "" || awsConfig.User == "" || awsConfig.Password == "" {
		return fmt.Errorf("AWS配置缺少必要参数")
	}

	// 检查数据库列表
	if len(databases) == 0 {
		return fmt.Errorf("数据库列表不能为空")
	}

	return nil
}

// getConnectionString 生成数据库连接字符串
func (v *MultiDatabaseValidator) getConnectionString(config DatabaseConfig) string {
	charset := config.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/?charset=%s&parseTime=true&loc=Local",
		config.User, config.Password, config.Host, charset)
}

// getTableList 获取指定数据库的表列表
func (v *MultiDatabaseValidator) getTableList(db *sql.DB, database string) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = ? 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.Query(query, database)
	if err != nil {
		return nil, fmt.Errorf("查询表列表失败: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("扫描表名失败: %v", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历表列表失败: %v", err)
	}

	log.Printf("数据库 %s 包含 %d 个表", database, len(tables))
	return tables, nil
}

// calculateTableChecksum 计算表的校验和
func (v *MultiDatabaseValidator) calculateTableChecksum(db *sql.DB, database, tableName string) (string, error) {
	// 获取表的行数
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", database, tableName)
	if err := db.QueryRow(countQuery).Scan(&rowCount); err != nil {
		return "", fmt.Errorf("获取表行数失败: %v", err)
	}

	// 空表处理
	if rowCount == 0 {
		log.Printf("表 %s.%s 为空表", database, tableName)
		return "empty_table", nil
	}

	// 大表分批处理，避免内存溢出
	if rowCount > 100000 {
		log.Printf("表 %s.%s 行数较多(%d)，使用分批计算", database, tableName, rowCount)
		return v.calculateLargeTableChecksum(db, database, tableName)
	}

	// 小表直接计算
	log.Printf("计算表 %s.%s 的校验和，行数: %d", database, tableName, rowCount)
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s` ORDER BY 1", database, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("查询表数据失败: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("获取列信息失败: %v", err)
	}

	// 读取所有数据
	var allData []string
	for rows.Next() {
		// 创建扫描目标
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", fmt.Errorf("扫描行数据失败: %v", err)
		}

		// 将数据转换为字符串
		rowData := make([]string, len(columns))
		for i, val := range values {
			if val != nil {
				rowData[i] = fmt.Sprintf("%v", val)
			} else {
				rowData[i] = "NULL"
			}
		}

		// 将行数据添加到总数据中
		rowStr := fmt.Sprintf("%v", rowData)
		allData = append(allData, rowStr)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("遍历表数据失败: %v", err)
	}

	// 计算MD5校验和
	dataStr := fmt.Sprintf("%v", allData)
	hash := md5.Sum([]byte(dataStr))
	return fmt.Sprintf("%x", hash), nil
}

// calculateLargeTableChecksum 大表分批计算校验和
func (v *MultiDatabaseValidator) calculateLargeTableChecksum(db *sql.DB, database, tableName string) (string, error) {
	// 获取总行数
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", database, tableName)
	if err := db.QueryRow(countQuery).Scan(&totalRows); err != nil {
		return "", fmt.Errorf("获取表总行数失败: %v", err)
	}

	batchSize := 10000 // 每批处理10000行
	var checksums []string

	log.Printf("开始分批计算表 %s.%s，总行数: %d", database, tableName, totalRows)

	// 分批读取数据
	for offset := 0; offset < totalRows; offset += batchSize {
		log.Printf("处理批次: %d-%d", offset, min(offset+batchSize, totalRows))

		query := fmt.Sprintf(`
			SELECT * FROM `%s`.`%s` 
			ORDER BY 1 
			LIMIT %d OFFSET %d
		`, database, tableName, batchSize, offset)

		rows, err := db.Query(query)
		if err != nil {
			return "", fmt.Errorf("查询批次数据失败: %v", err)
		}

		// 获取列信息
		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			return "", fmt.Errorf("获取列信息失败: %v", err)
		}

		// 读取批次数据
		var batchData []string
		for rows.Next() {
			// 创建扫描目标
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			// 扫描行数据
			if err := rows.Scan(valuePtrs...); err != nil {
				rows.Close()
				return "", fmt.Errorf("扫描批次数据失败: %v", err)
			}

			// 将数据转换为字符串
			rowData := make([]string, len(columns))
			for i, val := range values {
				if val != nil {
					rowData[i] = fmt.Sprintf("%v", val)
				} else {
					rowData[i] = "NULL"
				}
			}

			// 将行数据添加到批次数据中
			rowStr := fmt.Sprintf("%v", rowData)
			batchData = append(batchData, rowStr)
		}

		rows.Close()

		if err := rows.Err(); err != nil {
			return "", fmt.Errorf("遍历批次数据失败: %v", err)
		}

		// 计算批次数据的校验和
		batchStr := fmt.Sprintf("%v", batchData)
		hash := md5.Sum([]byte(batchStr))
		checksums = append(checksums, fmt.Sprintf("%x", hash))
	}

	// 合并所有批次的校验和
	combinedStr := ""
	for _, checksum := range checksums {
		combinedStr += checksum
	}
	
	combinedHash := md5.Sum([]byte(combinedStr))
	log.Printf("表 %s.%s 分批计算完成，共 %d 个批次", database, tableName, len(checksums))
	
	return fmt.Sprintf("%x", combinedHash), nil
}

// validateDatabase 验证单个数据库的一致性
func (v *MultiDatabaseValidator) validateDatabase(database string) DatabaseResult {
	log.Printf("开始验证数据库: %s", database)

	startTime := time.Now()
	result := DatabaseResult{
		Database:         database,
		Status:           "SUCCESS",
		Errors:           []string{},
		TableComparisons: []TableComparison{},
		StartTime:        startTime.Format(time.RFC3339),
	}

	// 连接数据库
	azureConnStr := v.getConnectionString(v.azureConfig)
	awsConnStr := v.getConnectionString(v.awsConfig)

	azureDB, err := sql.Open("mysql", azureConnStr)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("连接Azure数据库失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}
	defer azureDB.Close()

	awsDB, err := sql.Open("mysql", awsConnStr)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("连接AWS数据库失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}
	defer awsDB.Close()

	// 测试连接
	if err := azureDB.Ping(); err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("Azure数据库连接测试失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	if err := awsDB.Ping(); err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("AWS数据库连接测试失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	// 获取表列表
	azureTables, err := v.getTableList(azureDB, database)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("获取Azure表列表失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	awsTables, err := v.getTableList(awsDB, database)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("获取AWS表列表失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	result.AzureTables = len(azureTables)
	result.AWSTables = len(awsTables)

	// 检查表数量一致性
	if len(azureTables) != len(awsTables) {
		result.Status = "WARNING"
		errorMsg := fmt.Sprintf("表数量不一致: Azure(%d) vs AWS(%d)", len(azureTables), len(awsTables))
		result.Errors = append(result.Errors, errorMsg)
		log.Printf("数据库 %s: %s", database, errorMsg)
	}

	// 对比每个表的数据一致性
	log.Printf("开始验证数据库 %s 中的 %d 个表", database, len(azureTables))

	for i, table := range azureTables {
		log.Printf("验证表 %d/%d: %s", i+1, len(azureTables), table)

		// 检查表是否在AWS中存在
		tableExists := false
		for _, awsTable := range awsTables {
			if awsTable == table {
				tableExists = true
				break
			}
		}

		if !tableExists {
			errorMsg := fmt.Sprintf("表 %s 在AWS中不存在", table)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "INCONSISTENT"
			log.Printf("数据库 %s: %s", database, errorMsg)
			continue
		}

		// 计算Azure表的校验和
		azureChecksum, err := v.calculateTableChecksum(azureDB, database, table)
		if err != nil {
			errorMsg := fmt.Sprintf("计算Azure表 %s 校验和失败: %v", table, err)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "ERROR"
			log.Printf("数据库 %s: %s", database, errorMsg)
			continue
		}

		// 计算AWS表的校验和
		awsChecksum, err := v.calculateTableChecksum(awsDB, database, table)
		if err != nil {
			errorMsg := fmt.Sprintf("计算AWS表 %s 校验和失败: %v", table, err)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "ERROR"
			log.Printf("数据库 %s: %s", database, errorMsg)
			continue
		}

		// 记录对比结果
		tableComparison := TableComparison{
			Table:         table,
			AzureChecksum: azureChecksum,
			AWSChecksum:   awsChecksum,
			Match:         azureChecksum == awsChecksum,
		}

		result.TableComparisons = append(result.TableComparisons, tableComparison)

		// 检查是否一致
		if azureChecksum != awsChecksum {
			result.Status = "INCONSISTENT"
			log.Printf("表 %s.%s 数据不一致", database, table)
		} else {
			log.Printf("表 %s.%s 数据一致", database, table)
		}
	}

	// 记录结束时间
	result.EndTime = time.Now().Format(time.RFC3339)

	// 统计验证结果
	consistentTables := 0
	for _, comparison := range result.TableComparisons {
		if comparison.Match {
			consistentTables++
		}
	}
	totalTables := len(result.TableComparisons)

	log.Printf("数据库 %s 验证完成:", database)
	log.Printf("  状态: %s", result.Status)
	log.Printf("  表数量: Azure(%d) vs AWS(%d)", result.AzureTables, result.AWSTables)
	log.Printf("  数据一致性: %d/%d", consistentTables, totalTables)
	log.Printf("  错误数量: %d", len(result.Errors))

	return result
}

// validateAllDatabases 并行验证所有数据库
func (v *MultiDatabaseValidator) validateAllDatabases(maxWorkers int) {
	log.Printf("开始验证 %d 个数据库，最大并发数: %d", len(v.databases), maxWorkers)

	// 使用goroutine和channel实现并发控制
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, database := range v.databases {
		wg.Add(1)
		go func(db string) {
			defer wg.Done()
			
			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行验证
			result := v.validateDatabase(db)
			
			// 保存结果
			v.mutex.Lock()
			v.results[db] = result
			v.mutex.Unlock()

			log.Printf("数据库 %s 验证完成，状态: %s", db, result.Status)
		}(database)
	}

	// 等待所有验证完成
	wg.Wait()
	log.Printf("所有数据库验证完成")
}

// generateReport 生成验证报告
func (v *MultiDatabaseValidator) generateReport(outputFile string) (*ValidationSummary, error) {
	// 统计验证结果
	totalDatabases := len(v.databases)
	successfulValidations := 0
	inconsistentDatabases := 0
	errorDatabases := 0

	v.mutex.RLock()
	for _, result := range v.results {
		switch result.Status {
		case "SUCCESS":
			successfulValidations++
		case "INCONSISTENT":
			inconsistentDatabases++
		case "ERROR":
			errorDatabases++
		}
	}
	v.mutex.RUnlock()

	// 计算成功率
	successRate := "0%"
	if totalDatabases > 0 {
		successRate = fmt.Sprintf("%.2f%%", float64(successfulValidations)/float64(totalDatabases)*100)
	}

	// 生成报告摘要
	summary := &ValidationSummary{
		Timestamp:              time.Now().Format(time.RFC3339),
		TotalDatabases:         totalDatabases,
		SuccessfulValidations:  successfulValidations,
		InconsistentDatabases:  inconsistentDatabases,
		ErrorDatabases:         errorDatabases,
		SuccessRate:            successRate,
		Results:                make(map[string]DatabaseResult),
	}

	// 复制结果
	v.mutex.RLock()
	for db, result := range v.results {
		summary.Results[db] = result
	}
	v.mutex.RUnlock()

	// 保存报告到文件
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("创建报告文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		return nil, fmt.Errorf("写入报告文件失败: %v", err)
	}

	log.Printf("验证报告已生成: %s", outputFile)
	return summary, nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// main 主函数 - 程序入口点
func main() {
	// 配置数据库连接参数
	azureConfig := DatabaseConfig{
		Host:     "your-azure-mysql.mysql.database.azure.com",
		User:     "your_username",
		Password: "your_password",
		Charset:  "utf8mb4",
	}

	awsConfig := DatabaseConfig{
		Host:     "your-aws-rds.region.rds.amazonaws.com",
		User:     "your_username",
		Password: "your_password",
		Charset:  "utf8mb4",
	}

	// 需要验证的数据库列表
	databases := []string{"db1", "db2", "db3", "db4", "db5"}

	// 创建验证器实例
	validator := NewMultiDatabaseValidator(azureConfig, awsConfig, databases)

	// 执行验证
	validator.validateAllDatabases(3) // 最大并发数为3

	// 生成报告
	summary, err := validator.generateReport("consistency_report.json")
	if err != nil {
		log.Fatalf("生成报告失败: %v", err)
	}

	// 输出验证摘要
	fmt.Println("\n" + "="*60)
	fmt.Println("数据库一致性验证完成")
	fmt.Println("="*60)
	fmt.Printf("验证时间: %s\n", summary.Timestamp)
	fmt.Printf("总数据库数: %d\n", summary.TotalDatabases)
	fmt.Printf("验证成功: %d\n", summary.SuccessfulValidations)
	fmt.Printf("数据不一致: %d\n", summary.InconsistentDatabases)
	fmt.Printf("验证错误: %d\n", summary.ErrorDatabases)
	fmt.Printf("成功率: %s\n", summary.SuccessRate)
	fmt.Println("="*60)

	// 如果有不一致或错误，显示详细信息
	if summary.InconsistentDatabases > 0 || summary.ErrorDatabases > 0 {
		fmt.Println("\n详细信息:")
		for dbName, result := range summary.Results {
			if result.Status == "INCONSISTENT" || result.Status == "ERROR" {
				fmt.Printf("\n数据库: %s\n", dbName)
				fmt.Printf("状态: %s\n", result.Status)
				if len(result.Errors) > 0 {
					fmt.Println("错误信息:")
					for _, error := range result.Errors {
						fmt.Printf("  - %s\n", error)
					}
				}
			}
		}
	}

	// 根据验证结果设置退出码
	if summary.ErrorDatabases > 0 {
		os.Exit(1)
	}
}
```

### Go版本使用说明

Go语言版本已经重构为模块化结构，包含以下文件：

#### 项目结构
```
go-validator/
├── go.mod              # Go模块文件
├── main.go             # 主程序入口
├── types.go            # 数据结构定义
├── validator.go        # 验证器核心逻辑
├── config.go           # 配置文件处理
├── config.json.example # 配置文件示例
├── run_example.sh      # 使用示例脚本
└── README.md           # 详细说明文档
```

#### 1. 快速开始
```bash
# 进入项目目录
cd go-validator

# 运行示例脚本
./run_example.sh

# 或者手动执行
go mod tidy
go run . init config.json
```

#### 2. 配置文件
编辑 `config.json` 文件：
```json
{
  "azure": {
    "host": "your-azure-mysql.mysql.database.azure.com",
    "user": "your_username",
    "password": "your_password",
    "charset": "utf8mb4"
  },
  "aws": {
    "host": "your-aws-rds.region.rds.amazonaws.com",
    "user": "your_username",
    "password": "your_password",
    "charset": "utf8mb4"
  },
  "databases": ["db1", "db2", "db3", "db4", "db5"],
  "max_workers": 3
}
```

#### 3. 运行验证
```bash
# 使用配置文件运行
go run .

# 或者编译后运行
go build -o validator
./validator
```

#### 4. 命令行选项
```bash
# 显示帮助信息
go run . help

# 创建默认配置文件
go run . init [filename]

# 运行验证
go run .
```

#### 5. 输出结果
- **控制台输出**: 实时验证进度和摘要
- **consistency_report.json**: 详细的验证报告
- **日志文件**: 完整的验证过程日志

#### 6. 主要特性
- **模块化设计**: 代码分离，易于维护
- **配置文件支持**: 灵活的配置管理
- **并发控制**: 可配置的最大并发数
- **大表处理**: 自动分批处理大表
- **详细报告**: JSON格式的验证报告
- **错误处理**: 完善的错误处理和日志记录

## 方案四：AWS DMS验证（实时验证）

### 适用场景
- 需要实时验证
- 预算充足
- 需要自动化监控

### DMS验证配置
```json
{
  "TaskSettings": {
    "TargetMetadata": {
      "TargetSchema": "",
      "SupportLobs": true,
      "FullLobMode": false,
      "LobChunkSize": 0,
      "LimitedSizeLobMode": true,
      "LobMaxSize": 32,
      "InlineLobMaxSize": 0,
      "LoadMaxFileSize": 0,
      "ParallelLoadThreads": 0,
      "ParallelLoadBufferSize": 0,
      "BatchApplyEnabled": false,
      "TaskRecoveryTableEnabled": false,
      "ParallelApplyThreads": 0,
      "ParallelApplyBufferSize": 0,
      "ParallelApplyQueuesPerThread": 0
    },
    "FullLoadSettings": {
      "TargetTablePrepMode": "DO_NOTHING",
      "CreatePkAfterFullLoad": false,
      "StopTaskCachedChangesApplied": false,
      "StopTaskCachedChangesNotApplied": false,
      "MaxFullLoadSubTasks": 8,
      "TransactionConsistencyTimeout": 600,
      "CommitRate": 10000
    },
    "Logging": {
      "EnableLogging": true,
      "LogComponents": [
        {
          "Id": "SOURCE_UNLOAD",
          "Severity": "LOGGER_SEVERITY_DEFAULT"
        },
        {
          "Id": "TARGET_LOAD",
          "Severity": "LOGGER_SEVERITY_DEFAULT"
        }
      ]
    },
    "ControlTablesSettings": {
      "historyTimeslotInMinutes": 5,
      "ControlSchema": "",
      "HistoryTimeslotInMinutes": 5,
      "HistoryTableEnabled": false,
      "SuspendedTablesTableEnabled": false,
      "StatusTableEnabled": false
    },
    "StreamBufferSettings": {
      "StreamBufferCount": 3,
      "StreamBufferSizeInMB": 8,
      "CtrlStreamBufferSizeInMB": 5
    },
    "ChangeProcessingDdlHandlingPolicy": {
      "HandleSourceTableDropped": true,
      "HandleSourceTableTruncated": true,
      "HandleSourceTableAltered": true
    },
    "ErrorBehavior": {
      "DataErrorPolicy": "LOG_ERROR",
      "DataTruncationErrorPolicy": "LOG_ERROR",
      "DataErrorEscalationPolicy": "SUSPEND_TABLE",
      "DataErrorEscalationCount": 0,
      "TableErrorPolicy": "SUSPEND_TABLE",
      "TableErrorEscalationPolicy": "STOP_TASK",
      "TableErrorEscalationCount": 0,
      "RecoverableErrorCount": -1,
      "RecoverableErrorInterval": 5,
      "RecoverableErrorThrottling": true,
      "RecoverableErrorThrottlingMax": 1800,
      "ApplyErrorDeletePolicy": "IGNORE_RECORD",
      "ApplyErrorInsertPolicy": "LOG_ERROR",
      "ApplyErrorUpdatePolicy": "LOG_ERROR",
      "ApplyErrorEscalationPolicy": "LOG_ERROR",
      "ApplyErrorEscalationCount": 0,
      "ApplyErrorFailOnTruncationDdl": false,
      "FullLoadIgnoreConflicts": true,
      "FailOnTransactionLogOnly": false,
      "FailOnNoTablesCaptured": true
    },
    "ChangeProcessingTuning": {
      "BatchApplyPreserveTransaction": true,
      "BatchApplyTimeoutMin": 1,
      "BatchApplyTimeoutMax": 30,
      "BatchApplyMemoryLimit": 500,
      "BatchSplitSize": 0,
      "MinTransactionSize": 1000,
      "CommitTimeoutMs": 1,
      "MemoryLimitTotal": 1024,
      "MemoryKeepTime": 60,
      "StatementCacheSize": 50
    }
  }
}
```

## 推荐方案选择

### 对于你的场景（多数据库，DBA负责），推荐：

1. **首选：pt-table-checksum**
   - 专业、稳定、适合生产环境
   - 支持大数据库
   - 有详细的文档和社区支持

2. **备选：自建Python脚本**
   - 灵活、可控
   - 可以集成到现有系统
   - 支持并行验证

3. **简单场景：mysqldbcompare**
   - 官方工具，易用
   - 适合中小型数据库

### 实施建议：

1. **分阶段验证**：先验证1-2个数据库，确认方案可行
2. **并行执行**：使用多线程/多进程提高效率
3. **结果记录**：详细记录验证结果，便于问题定位
4. **监控告警**：设置验证失败时的告警机制
5. **定期验证**：建立定期验证机制，确保长期一致性

需要我详细解释某个方案的具体实施步骤吗？
