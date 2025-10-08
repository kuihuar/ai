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
import argparse

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
    1. 连接Azure和AWS的多个MySQL数据库实例
    2. 获取每个数据库实例的表列表
    3. 计算每个表的数据校验和
    4. 对比两个环境的数据一致性
    5. 生成详细的验证报告
    """
    
    def __init__(self, database_pairs):
        """
        初始化验证器
        
        Args:
            database_pairs (list): 数据库对比对列表，每个元素包含azure_instance和aws_instance
        """
        self.database_pairs = database_pairs
        self.results = {}  # 存储验证结果
        
        # 验证配置参数
        self._validate_config()
    
    def _validate_config(self):
        """验证配置参数的有效性"""
        if not self.database_pairs:
            raise ValueError("数据库对比对列表不能为空")
        
        required_keys = ['name', 'host', 'user', 'password', 'database']
        
        for i, pair in enumerate(self.database_pairs):
            # 检查Azure实例配置
            if 'azure_instance' not in pair or 'aws_instance' not in pair:
                raise ValueError(f"第{i+1}个对比对缺少azure_instance或aws_instance")
            
            azure_instance = pair['azure_instance']
            aws_instance = pair['aws_instance']
            
            for key in required_keys:
                if key not in azure_instance:
                    raise ValueError(f"第{i+1}个对比对的Azure实例配置缺少必要参数: {key}")
                if key not in aws_instance:
                    raise ValueError(f"第{i+1}个对比对的AWS实例配置缺少必要参数: {key}")
    
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
    
    def validate_database(self, pair):
        """
        验证单个数据库对比对的一致性
        
        验证步骤：
        1. 连接Azure和AWS数据库实例
        2. 获取表列表
        3. 对比表数量
        4. 逐个验证表的数据一致性
        5. 记录验证结果
        
        Args:
            pair (dict): 数据库对比对，包含azure_instance和aws_instance
            
        Returns:
            dict: 验证结果
        """
        azure_instance = pair['azure_instance']
        aws_instance = pair['aws_instance']
        
        logging.info(f"开始验证数据库对比: {azure_instance['database']} (Azure: {azure_instance['name']}) vs {aws_instance['database']} (AWS: {aws_instance['name']})")
        
        azure_conn = None
        aws_conn = None
        
        try:
            # 连接数据库
            logging.debug(f"连接Azure数据库: {azure_instance['host']}")
            azure_conn = pymysql.connect(
                host=azure_instance['host'],
                user=azure_instance['user'],
                password=azure_instance['password'],
                database=azure_instance['database'],
                charset=azure_instance.get('charset', 'utf8mb4'),
                connect_timeout=30,
                read_timeout=30,
                write_timeout=30
            )
            
            logging.debug(f"连接AWS数据库: {aws_instance['host']}")
            aws_conn = pymysql.connect(
                host=aws_instance['host'],
                user=aws_instance['user'],
                password=aws_instance['password'],
                database=aws_instance['database'],
                charset=aws_instance.get('charset', 'utf8mb4'),
                connect_timeout=30,
                read_timeout=30,
                write_timeout=30
            )
            
            # 获取表列表
            azure_tables = self.get_table_list(azure_conn, azure_instance['database'])
            aws_tables = self.get_table_list(aws_conn, aws_instance['database'])
            
            # 初始化验证结果
            result = {
                'database': azure_instance['database'],
                'azure_instance': azure_instance['name'],
                'aws_instance': aws_instance['name'],
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
                logging.warning(f"数据库 {azure_instance['database']}: {error_msg}")
            
            # 对比每个表的数据一致性
            logging.info(f"开始验证数据库 {azure_instance['database']} 中的 {len(azure_tables)} 个表")
            
            for i, table in enumerate(azure_tables, 1):
                logging.info(f"验证表 {i}/{len(azure_tables)}: {table}")
                
                if table in aws_tables:
                    try:
                        # 计算Azure表的校验和
                        azure_checksum = self.calculate_table_checksum(azure_conn, azure_instance['database'], table)
                        
                        # 计算AWS表的校验和
                        aws_checksum = self.calculate_table_checksum(aws_conn, aws_instance['database'], table)
                        
                        # 记录对比结果
                        table_result = {
                            'table': table,
                            'azure_checksum': azure_checksum,
                            'aws_checksum': aws_checksum,
                            'match': azure_checksum == aws_checksum,
                            'azure_instance': azure_instance['name'],
                            'aws_instance': aws_instance['name'],
                            'azure_database': azure_instance['database'],
                            'aws_database': aws_instance['database']
                        }
                        
                        result['table_comparisons'].append(table_result)
                        
                        # 检查是否一致
                        if azure_checksum != aws_checksum:
                            result['status'] = 'INCONSISTENT'
                            logging.warning(f"数据不一致 - Azure实例: {azure_instance['name']} 数据库: {azure_instance['database']} 表: {table} vs AWS实例: {aws_instance['name']} 数据库: {aws_instance['database']} 表: {table}")
                        else:
                            logging.debug(f"数据一致 - Azure实例: {azure_instance['name']} 数据库: {azure_instance['database']} 表: {table} vs AWS实例: {aws_instance['name']} 数据库: {aws_instance['database']} 表: {table}")
                            
                    except Exception as e:
                        error_msg = f"表 {table} 验证失败: {str(e)}"
                        result['errors'].append(error_msg)
                        result['status'] = 'ERROR'
                        logging.error(f"数据库 {azure_instance['database']}: {error_msg}")
                else:
                    error_msg = f"表 {table} 在AWS中不存在"
                    result['errors'].append(error_msg)
                    result['status'] = 'INCONSISTENT'
                    logging.warning(f"数据库 {azure_instance['database']}: {error_msg}")
            
            # 记录结束时间
            result['end_time'] = datetime.now().isoformat()
            
            # 统计验证结果
            consistent_tables = sum(1 for t in result['table_comparisons'] if t['match'])
            total_tables = len(result['table_comparisons'])
            
            logging.info(f"数据库 {azure_instance['database']} 验证完成:")
            logging.info(f"  状态: {result['status']}")
            logging.info(f"  表数量: Azure({result['azure_tables']}) vs AWS({result['aws_tables']})")
            logging.info(f"  数据一致性: {consistent_tables}/{total_tables}")
            logging.info(f"  错误数量: {len(result['errors'])}")
            
            return result
            
        except Exception as e:
            error_msg = f"数据库 {azure_instance['database']} 验证失败: {str(e)}"
            logging.error(error_msg)
            return {
                'database': azure_instance['database'],
                'azure_instance': azure_instance['name'],
                'aws_instance': aws_instance['name'],
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
        并行验证所有数据库对比对
        
        使用线程池并行处理多个数据库对比对，提高验证效率
        
        Args:
            max_workers (int): 最大并发线程数，默认3个
        """
        logging.info(f"开始验证 {len(self.database_pairs)} 个数据库对比对，最大并发数: {max_workers}")
        
        # 使用线程池并行处理
        with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
            # 提交所有验证任务
            future_to_pair = {
                executor.submit(self.validate_database, pair): pair 
                for pair in self.database_pairs
            }
            
            # 收集验证结果
            for future in concurrent.futures.as_completed(future_to_pair):
                pair = future_to_pair[future]
                try:
                    result = future.result()
                    db_name = result['database']
                    self.results[db_name] = result
                    logging.info(f"数据库对比 {result['azure_instance']} vs {result['aws_instance']} 验证完成，状态: {result['status']}")
                except Exception as e:
                    error_msg = f"数据库对比验证异常: {str(e)}"
                    logging.error(error_msg)
                    db_name = pair['azure_instance']['database']
                    self.results[db_name] = {
                        'database': db_name,
                        'azure_instance': pair['azure_instance']['name'],
                        'aws_instance': pair['aws_instance']['name'],
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
        total_databases = len(self.database_pairs)
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

def load_config(config_file):
    """
    从配置文件加载配置
    
    Args:
        config_file (str): 配置文件路径
        
    Returns:
        list: 数据库对比对列表
    """
    try:
        with open(config_file, 'r', encoding='utf-8') as f:
            config = json.load(f)
        
        # 验证配置结构
        if 'azure' not in config or 'aws' not in config:
            raise ValueError("配置文件必须包含azure和aws字段")
        
        if len(config['azure']) != len(config['aws']):
            raise ValueError("Azure和AWS实例数量必须相同")
        
        # 构建数据库对比对
        database_pairs = []
        for i in range(len(config['azure'])):
            pair = {
                'azure_instance': config['azure'][i],
                'aws_instance': config['aws'][i]
            }
            database_pairs.append(pair)
        
        return database_pairs, config.get('max_workers', 3)
        
    except FileNotFoundError:
        raise FileNotFoundError(f"配置文件不存在: {config_file}")
    except json.JSONDecodeError as e:
        raise ValueError(f"配置文件格式错误: {e}")
    except Exception as e:
        raise ValueError(f"加载配置文件失败: {e}")

def create_default_config(config_file):
    """
    创建默认配置文件
    
    Args:
        config_file (str): 配置文件路径
    """
    default_config = {
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
    
    with open(config_file, 'w', encoding='utf-8') as f:
        json.dump(default_config, f, indent=2, ensure_ascii=False)
    
    print(f"默认配置文件已创建: {config_file}")

def main():
    """主函数 - 程序入口点"""
    parser = argparse.ArgumentParser(description='多数据库一致性验证工具')
    parser.add_argument('--config', '-c', default='config.json', 
                       help='配置文件路径 (默认: config.json)')
    parser.add_argument('--init', action='store_true',
                       help='创建默认配置文件')
    parser.add_argument('--max-workers', type=int, default=3,
                       help='最大并发数 (默认: 3)')
    
    args = parser.parse_args()
    
    # 处理init命令
    if args.init:
        config_file = args.config
        create_default_config(config_file)
        return 0
    
    try:
        # 尝试从配置文件加载
        if os.path.exists(args.config):
            print(f"从配置文件加载配置: {args.config}")
            database_pairs, max_workers = load_config(args.config)
            # 使用命令行参数覆盖配置文件中的max_workers
            if args.max_workers != 3:
                max_workers = args.max_workers
        else:
            print("配置文件不存在，使用默认配置")
            # 使用默认配置
            database_pairs = [
                {
                    'azure_instance': {
                        'name': 'azure-db1',
                        'host': 'your-azure-mysql1.mysql.database.azure.com',
                        'user': 'your_username',
                        'password': 'your_password',
                        'database': 'db1',
                        'charset': 'utf8mb4'
                    },
                    'aws_instance': {
                        'name': 'aws-db1',
                        'host': 'your-aws-rds1.region.rds.amazonaws.com',
                        'user': 'your_username',
                        'password': 'your_password',
                        'database': 'db1',
                        'charset': 'utf8mb4'
                    }
                }
            ]
            max_workers = args.max_workers
        
        # 创建验证器实例
        validator = MultiDatabaseValidator(database_pairs)
        
        # 执行验证
        validator.validate_all_databases(max_workers=max_workers)
        
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
                    print(f"Azure实例: {result['azure_instance']}")
                    print(f"AWS实例: {result['aws_instance']}")
                    print(f"状态: {result['status']}")
                    
                    # 显示不一致的表信息
                    if result['status'] == 'INCONSISTENT' and 'table_comparisons' in result:
                        inconsistent_tables = [t for t in result['table_comparisons'] if not t['match']]
                        if inconsistent_tables:
                            print("不一致的表:")
                            for table in inconsistent_tables:
                                print(f"  - 表名: {table['table']}")
                                print(f"    Azure实例: {table['azure_instance']} 数据库: {table['azure_database']}")
                                print(f"    AWS实例: {table['aws_instance']} 数据库: {table['aws_database']}")
                                print(f"    Azure校验和: {table['azure_checksum']}")
                                print(f"    AWS校验和: {table['aws_checksum']}")
                    
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
