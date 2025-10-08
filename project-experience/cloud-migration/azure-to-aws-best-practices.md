# Azure到AWS迁移最佳实践

## 概述

从Azure迁移到AWS是一个复杂的项目，涉及多个层面的考虑。本文档总结了迁移过程中的最佳实践和经验教训。

## 迁移策略

### 1. 迁移方法选择

#### 重新部署 (Rehost)
- **适用场景**: 简单的lift-and-shift迁移
- **优点**: 快速、风险低
- **缺点**: 无法充分利用云原生特性

#### 重构 (Refactor)
- **适用场景**: 需要优化架构的应用
- **优点**: 提升性能、降低成本
- **缺点**: 开发工作量大、风险较高

#### 重新架构 (Re-architect)
- **适用场景**: 现代化应用架构
- **优点**: 充分利用云原生服务
- **缺点**: 需要大量重构工作

### 2. 迁移阶段规划

```
阶段1: 评估和规划 (2-4周)
├── 应用清单和依赖分析
├── 成本分析
├── 风险评估
└── 迁移计划制定

阶段2: 试点迁移 (4-6周)
├── 选择非关键应用
├── 建立迁移流程
├── 验证工具和方法
└── 团队培训

阶段3: 批量迁移 (8-12周)
├── 按优先级迁移应用
├── 数据迁移
├── 网络和安全配置
└── 监控和优化

阶段4: 优化和清理 (2-4周)
├── 性能优化
├── 成本优化
├── 清理Azure资源
└── 文档更新
```

## 技术迁移指南

### 1. 计算服务迁移

#### Azure VM → AWS EC2
```bash
# 使用AWS Server Migration Service (SMS)
# 1. 在Azure VM上安装SMS Agent
# 2. 配置迁移任务
# 3. 执行迁移

# 或者使用第三方工具
# 如CloudEndure, Carbonite等
```

#### Azure App Service → AWS Elastic Beanstalk/ECS
```yaml
# Elastic Beanstalk配置示例
version: 1
application:
  name: my-app
  platform: node.js
  version: 16.x
environment:
  name: production
  solution_stack: 64bit Amazon Linux 2 v3.4.0 running Node.js 16
  option_settings:
    aws:elasticbeanstalk:container:nodejs:
      NodeCommand: "npm start"
      NodeVersion: 16.18.0
```

### 2. 数据库迁移

#### Azure SQL Database → AWS RDS
```sql
-- 使用AWS DMS进行迁移
-- 1. 创建DMS实例
-- 2. 配置源和目标端点
-- 3. 创建迁移任务
-- 4. 执行全量+增量迁移

-- 数据一致性验证
SELECT 
    table_name,
    COUNT(*) as row_count
FROM information_schema.tables 
WHERE table_schema = 'your_database'
GROUP BY table_name;
```

#### Azure Cosmos DB → AWS DynamoDB
```python
# 使用AWS Database Migration Service
# 或自定义迁移脚本

import boto3
from azure.cosmos import CosmosClient

def migrate_cosmos_to_dynamodb():
    # Azure Cosmos DB连接
    cosmos_client = CosmosClient(cosmos_endpoint, cosmos_key)
    database = cosmos_client.get_database_client(database_name)
    container = database.get_container_client(container_name)
    
    # AWS DynamoDB连接
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(table_name)
    
    # 迁移数据
    for item in container.read_all_items():
        table.put_item(Item=item)
```

### 3. 存储服务迁移

#### Azure Blob Storage → AWS S3
```python
import boto3
from azure.storage.blob import BlobServiceClient

def migrate_blob_to_s3():
    # Azure Blob Storage
    blob_service = BlobServiceClient(
        account_url=azure_account_url,
        credential=azure_credential
    )
    
    # AWS S3
    s3_client = boto3.client('s3')
    
    # 迁移文件
    container_client = blob_service.get_container_client(container_name)
    for blob in container_client.list_blobs():
        blob_data = container_client.download_blob(blob.name).readall()
        s3_client.put_object(
            Bucket=bucket_name,
            Key=blob.name,
            Body=blob_data
        )
```

### 4. 网络和安全迁移

#### 网络架构映射
```
Azure                    AWS
├── Virtual Network     → VPC
├── Subnet             → Subnet
├── Network Security Group → Security Group
├── Load Balancer      → Application/Network Load Balancer
├── VPN Gateway        → VPN Gateway
└── ExpressRoute       → Direct Connect
```

#### 安全配置迁移
```bash
# Azure Key Vault → AWS Secrets Manager
aws secretsmanager create-secret \
    --name "my-secret" \
    --description "Migrated from Azure Key Vault" \
    --secret-string '{"username":"admin","password":"secret"}'

# Azure Active Directory → AWS IAM
# 需要重新配置用户和权限
```

## 成本优化策略

### 1. 资源优化

#### 实例类型选择
```bash
# 使用AWS Compute Optimizer分析实例类型
aws compute-optimizer get-recommendation-summaries \
    --recommendation-sources EC2Instance

# 使用Reserved Instances降低成本
aws ec2 describe-reserved-instances-offerings \
    --instance-type t3.medium \
    --offering-type All Upfront
```

#### 存储优化
```bash
# 使用S3 Intelligent Tiering
aws s3api put-bucket-intelligent-tiering-configuration \
    --bucket my-bucket \
    --id my-config \
    --intelligent-tiering-configuration '{
        "Id": "my-config",
        "Status": "Enabled",
        "Tierings": [
            {
                "Days": 30,
                "AccessTier": "ARCHIVE_ACCESS"
            }
        ]
    }'
```

### 2. 监控和告警

```python
import boto3

def setup_cost_monitoring():
    cloudwatch = boto3.client('cloudwatch')
    
    # 设置成本告警
    cloudwatch.put_metric_alarm(
        AlarmName='High-Cost-Alert',
        ComparisonOperator='GreaterThanThreshold',
        EvaluationPeriods=1,
        MetricName='EstimatedCharges',
        Namespace='AWS/Billing',
        Period=86400,
        Statistic='Maximum',
        Threshold=1000.0,
        ActionsEnabled=True,
        AlarmActions=['arn:aws:sns:region:account:topic']
    )
```

## 迁移工具和脚本

### 1. 自动化迁移脚本

```python
#!/usr/bin/env python3
"""
Azure到AWS迁移自动化脚本
"""

import boto3
import json
from datetime import datetime

class AzureToAWSMigrator:
    def __init__(self, aws_region='us-east-1'):
        self.aws_region = aws_region
        self.ec2 = boto3.client('ec2', region_name=aws_region)
        self.rds = boto3.client('rds', region_name=aws_region)
        self.s3 = boto3.client('s3', region_name=aws_region)
    
    def create_vpc(self, vpc_config):
        """创建VPC"""
        response = self.ec2.create_vpc(
            CidrBlock=vpc_config['cidr_block'],
            TagSpecifications=[{
                'ResourceType': 'vpc',
                'Tags': [
                    {'Key': 'Name', 'Value': vpc_config['name']},
                    {'Key': 'Environment', 'Value': vpc_config['environment']}
                ]
            }]
        )
        return response['Vpc']['VpcId']
    
    def create_subnets(self, vpc_id, subnet_configs):
        """创建子网"""
        subnet_ids = []
        for config in subnet_configs:
            response = self.ec2.create_subnet(
                VpcId=vpc_id,
                CidrBlock=config['cidr_block'],
                AvailabilityZone=config['az'],
                TagSpecifications=[{
                    'ResourceType': 'subnet',
                    'Tags': [
                        {'Key': 'Name', 'Value': config['name']}
                    ]
                }]
            )
            subnet_ids.append(response['Subnet']['SubnetId'])
        return subnet_ids
    
    def create_security_groups(self, vpc_id, sg_configs):
        """创建安全组"""
        sg_ids = []
        for config in sg_configs:
            response = self.ec2.create_security_group(
                GroupName=config['name'],
                Description=config['description'],
                VpcId=vpc_id,
                TagSpecifications=[{
                    'ResourceType': 'security-group',
                    'Tags': [
                        {'Key': 'Name', 'Value': config['name']}
                    ]
                }]
            )
            sg_id = response['GroupId']
            
            # 添加入站规则
            if 'ingress_rules' in config:
                self.ec2.authorize_security_group_ingress(
                    GroupId=sg_id,
                    IpPermissions=config['ingress_rules']
                )
            
            sg_ids.append(sg_id)
        return sg_ids
    
    def migrate_database(self, db_config):
        """迁移数据库"""
        response = self.rds.create_db_instance(
            DBInstanceIdentifier=db_config['identifier'],
            DBInstanceClass=db_config['instance_class'],
            Engine=db_config['engine'],
            MasterUsername=db_config['username'],
            MasterUserPassword=db_config['password'],
            AllocatedStorage=db_config['allocated_storage'],
            VpcSecurityGroupIds=db_config['security_group_ids'],
            DBSubnetGroupName=db_config['subnet_group_name'],
            BackupRetentionPeriod=db_config['backup_retention'],
            MultiAZ=db_config['multi_az'],
            Tags=[
                {'Key': 'Name', 'Value': db_config['name']},
                {'Key': 'Environment', 'Value': db_config['environment']}
            ]
        )
        return response['DBInstance']['DBInstanceIdentifier']

# 使用示例
if __name__ == "__main__":
    migrator = AzureToAWSMigrator()
    
    # VPC配置
    vpc_config = {
        'name': 'migration-vpc',
        'cidr_block': '10.0.0.0/16',
        'environment': 'production'
    }
    
    # 创建VPC
    vpc_id = migrator.create_vpc(vpc_config)
    print(f"Created VPC: {vpc_id}")
```

### 2. 数据迁移脚本

```bash
#!/bin/bash
# 数据库迁移脚本

# 配置变量
AZURE_DB_HOST="your-azure-db.mysql.database.azure.com"
AWS_DB_HOST="your-aws-rds.region.rds.amazonaws.com"
DB_NAME="your_database"
DB_USER="your_user"
DB_PASSWORD="your_password"

# 创建备份
echo "Creating backup from Azure..."
mysqldump \
  --host=$AZURE_DB_HOST \
  --user=$DB_USER \
  --password=$DB_PASSWORD \
  --single-transaction \
  --routines \
  --triggers \
  --lock-tables=false \
  $DB_NAME > azure_backup.sql

# 恢复备份到AWS
echo "Restoring backup to AWS..."
mysql \
  --host=$AWS_DB_HOST \
  --user=$DB_USER \
  --password=$DB_PASSWORD \
  $DB_NAME < azure_backup.sql

# 验证数据
echo "Verifying data consistency..."
mysql \
  --host=$AZURE_DB_HOST \
  --user=$DB_USER \
  --password=$DB_PASSWORD \
  -e "SELECT COUNT(*) as azure_count FROM $DB_NAME.your_table;" > azure_count.txt

mysql \
  --host=$AWS_DB_HOST \
  --user=$DB_USER \
  --password=$DB_PASSWORD \
  -e "SELECT COUNT(*) as aws_count FROM $DB_NAME.your_table;" > aws_count.txt

# 比较结果
if diff azure_count.txt aws_count.txt; then
    echo "Data migration successful!"
else
    echo "Data migration failed - counts don't match!"
    exit 1
fi

# 清理临时文件
rm azure_backup.sql azure_count.txt aws_count.txt
```

## 监控和运维

### 1. 迁移监控仪表板

```python
import boto3
import json
from datetime import datetime, timedelta

class MigrationMonitor:
    def __init__(self):
        self.cloudwatch = boto3.client('cloudwatch')
        self.ec2 = boto3.client('ec2')
        self.rds = boto3.client('rds')
    
    def get_migration_metrics(self):
        """获取迁移相关指标"""
        metrics = {}
        
        # EC2实例状态
        instances = self.ec2.describe_instances()
        running_instances = sum(
            1 for reservation in instances['Reservations']
            for instance in reservation['Instances']
            if instance['State']['Name'] == 'running'
        )
        metrics['running_instances'] = running_instances
        
        # RDS实例状态
        db_instances = self.rds.describe_db_instances()
        available_dbs = sum(
            1 for db in db_instances['DBInstances']
            if db['DBInstanceStatus'] == 'available'
        )
        metrics['available_databases'] = available_dbs
        
        return metrics
    
    def create_migration_dashboard(self):
        """创建迁移监控仪表板"""
        dashboard_body = {
            "widgets": [
                {
                    "type": "metric",
                    "properties": {
                        "metrics": [
                            ["AWS/EC2", "CPUUtilization"],
                            [".", "NetworkIn"],
                            [".", "NetworkOut"]
                        ],
                        "period": 300,
                        "stat": "Average",
                        "region": "us-east-1",
                        "title": "EC2 Metrics"
                    }
                },
                {
                    "type": "metric",
                    "properties": {
                        "metrics": [
                            ["AWS/RDS", "CPUUtilization"],
                            [".", "DatabaseConnections"],
                            [".", "FreeableMemory"]
                        ],
                        "period": 300,
                        "stat": "Average",
                        "region": "us-east-1",
                        "title": "RDS Metrics"
                    }
                }
            ]
        }
        
        self.cloudwatch.put_dashboard(
            DashboardName='Migration-Dashboard',
            DashboardBody=json.dumps(dashboard_body)
        )

# 使用示例
monitor = MigrationMonitor()
metrics = monitor.get_migration_metrics()
print(f"Migration Metrics: {metrics}")
monitor.create_migration_dashboard()
```

### 2. 告警配置

```python
def setup_migration_alerts():
    """设置迁移告警"""
    cloudwatch = boto3.client('cloudwatch')
    
    # 高CPU使用率告警
    cloudwatch.put_metric_alarm(
        AlarmName='High-CPU-Migration',
        ComparisonOperator='GreaterThanThreshold',
        EvaluationPeriods=2,
        MetricName='CPUUtilization',
        Namespace='AWS/EC2',
        Period=300,
        Statistic='Average',
        Threshold=80.0,
        ActionsEnabled=True,
        AlarmActions=['arn:aws:sns:region:account:migration-alerts']
    )
    
    # 数据库连接数告警
    cloudwatch.put_metric_alarm(
        AlarmName='High-DB-Connections',
        ComparisonOperator='GreaterThanThreshold',
        EvaluationPeriods=1,
        MetricName='DatabaseConnections',
        Namespace='AWS/RDS',
        Period=300,
        Statistic='Average',
        Threshold=100.0,
        ActionsEnabled=True,
        AlarmActions=['arn:aws:sns:region:account:migration-alerts']
    )
```

## 风险管理和回滚策略

### 1. 风险评估矩阵

| 风险类型 | 概率 | 影响 | 缓解措施 |
|---------|------|------|----------|
| 数据丢失 | 低 | 高 | 多重备份、实时同步 |
| 服务中断 | 中 | 高 | 蓝绿部署、渐进式迁移 |
| 性能下降 | 中 | 中 | 性能测试、容量规划 |
| 成本超支 | 高 | 中 | 成本监控、资源优化 |
| 安全漏洞 | 低 | 高 | 安全审计、权限最小化 |

### 2. 回滚计划

```python
class RollbackManager:
    def __init__(self):
        self.azure_client = self.setup_azure_client()
        self.aws_client = boto3.client('ec2')
    
    def rollback_application(self, app_name):
        """回滚应用程序"""
        # 1. 停止AWS服务
        self.stop_aws_services(app_name)
        
        # 2. 恢复Azure服务
        self.restore_azure_services(app_name)
        
        # 3. 更新DNS记录
        self.update_dns_records(app_name, 'azure')
        
        # 4. 验证服务状态
        self.verify_service_health(app_name)
    
    def rollback_database(self, db_name):
        """回滚数据库"""
        # 1. 停止AWS RDS
        self.stop_aws_rds(db_name)
        
        # 2. 恢复Azure数据库
        self.restore_azure_database(db_name)
        
        # 3. 验证数据一致性
        self.verify_data_consistency(db_name)
```

## 总结

Azure到AWS的迁移是一个复杂的项目，需要：

1. **充分的规划**: 详细的迁移计划和风险评估
2. **分阶段执行**: 降低风险，便于问题定位
3. **自动化工具**: 提高效率，减少人为错误
4. **实时监控**: 及时发现和处理问题
5. **回滚准备**: 确保可以快速恢复

通过遵循这些最佳实践，可以大大提高迁移的成功率和效率。
