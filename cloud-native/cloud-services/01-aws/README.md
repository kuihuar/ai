# AWS 云原生服务详解

## 📚 学习目标

通过本模块学习，您将掌握：
- AWS 核心云原生服务架构
- 容器化服务 EKS、ECS、Fargate
- 无服务器服务 Lambda、API Gateway
- 存储和数据库服务
- 监控和日志服务
- 安全服务和最佳实践

## 🎯 AWS 云原生架构

### 1. AWS 服务生态

```
AWS 云原生服务
├── 计算服务
│   ├── EC2 (Elastic Compute Cloud)
│   ├── ECS (Elastic Container Service)
│   ├── EKS (Elastic Kubernetes Service)
│   ├── Fargate (Serverless Containers)
│   └── Lambda (Serverless Functions)
├── 存储服务
│   ├── S3 (Simple Storage Service)
│   ├── EBS (Elastic Block Store)
│   ├── EFS (Elastic File System)
│   └── FSx (Managed File Systems)
├── 数据库服务
│   ├── RDS (Relational Database Service)
│   ├── DynamoDB (NoSQL Database)
│   ├── ElastiCache (In-Memory Cache)
│   └── DocumentDB (MongoDB Compatible)
├── 网络服务
│   ├── VPC (Virtual Private Cloud)
│   ├── ALB/NLB (Load Balancers)
│   ├── CloudFront (CDN)
│   └── Route 53 (DNS)
└── 监控服务
    ├── CloudWatch (Monitoring)
    ├── X-Ray (Distributed Tracing)
    ├── CloudTrail (Audit Logs)
    └── Config (Configuration Management)
```

### 2. 服务对比

| 服务 | 类型 | 用途 | 优势 |
|------|------|------|------|
| ECS | 容器编排 | 容器管理 | 简单易用，AWS 原生 |
| EKS | Kubernetes | 容器编排 | 标准 K8s，生态丰富 |
| Fargate | 无服务器容器 | 容器运行 | 无需管理基础设施 |
| Lambda | 无服务器函数 | 事件驱动 | 按需付费，自动扩展 |

## 🐳 容器化服务

### 1. Amazon EKS (Elastic Kubernetes Service)

#### EKS 集群创建
```bash
# 创建 EKS 集群
eksctl create cluster \
  --name my-cluster \
  --version 1.24 \
  --region us-west-2 \
  --nodegroup-name workers \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 1 \
  --nodes-max 4 \
  --managed

# 更新 kubeconfig
aws eks update-kubeconfig --region us-west-2 --name my-cluster
```

#### EKS 配置示例
```yaml
# eksctl-config.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-west-2
  version: "1.24"

nodeGroups:
  - name: workers
    instanceType: t3.medium
    desiredCapacity: 3
    minSize: 1
    maxSize: 4
    ssh:
      allow: true
      publicKeyName: my-key

addons:
  - name: vpc-cni
    version: latest
  - name: coredns
    version: latest
  - name: kube-proxy
    version: latest
  - name: aws-ebs-csi-driver
    version: latest
```

### 2. Amazon ECS (Elastic Container Service)

#### ECS 任务定义
```json
{
  "family": "my-app",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::123456789012:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "my-app",
      "image": "123456789012.dkr.ecr.us-west-2.amazonaws.com/my-app:latest",
      "portMappings": [
        {
          "containerPort": 80,
          "protocol": "tcp"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/my-app",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### ECS 服务配置
```yaml
# ecs-service.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ecs-service-config
data:
  cluster: my-cluster
  service: my-app-service
  task-definition: my-app:1
  desired-count: "3"
  launch-type: FARGATE
  network-configuration: |
    {
      "awsvpcConfiguration": {
        "subnets": ["subnet-12345", "subnet-67890"],
        "securityGroups": ["sg-12345"],
        "assignPublicIp": "ENABLED"
      }
    }
```

### 3. AWS Fargate

#### Fargate 任务配置
```yaml
# fargate-task.yaml
apiVersion: v1
kind: Pod
metadata:
  name: fargate-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    ports:
    - containerPort: 80
  nodeSelector:
    eks.amazonaws.com/compute-type: fargate
```

## ⚡ 无服务器服务

### 1. AWS Lambda

#### Lambda 函数示例
```python
# lambda_function.py
import json
import boto3

def lambda_handler(event, context):
    # 处理 API Gateway 事件
    if 'httpMethod' in event:
        return {
            'statusCode': 200,
            'headers': {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': '*'
            },
            'body': json.dumps({
                'message': 'Hello from Lambda!',
                'event': event
            })
        }
    
    # 处理 S3 事件
    if 'Records' in event:
        for record in event['Records']:
            bucket = record['s3']['bucket']['name']
            key = record['s3']['object']['key']
            print(f'Processing {key} from {bucket}')
    
    return {'statusCode': 200}
```

#### Lambda 部署配置
```yaml
# serverless.yml
service: my-lambda-app

provider:
  name: aws
  runtime: python3.9
  region: us-west-2
  environment:
    STAGE: ${opt:stage, 'dev'}
  iam:
    role:
      statements:
        - Effect: Allow
          Action:
            - s3:GetObject
            - s3:PutObject
          Resource: "arn:aws:s3:::my-bucket/*"

functions:
  hello:
    handler: lambda_function.lambda_handler
    events:
      - http:
          path: hello
          method: get
      - s3:
          bucket: my-bucket
          event: s3:ObjectCreated:*

plugins:
  - serverless-python-requirements
```

### 2. API Gateway

#### API Gateway 配置
```yaml
# api-gateway.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-gateway-config
data:
  api-spec: |
    {
      "openapi": "3.0.0",
      "info": {
        "title": "My API",
        "version": "1.0.0"
      },
      "paths": {
        "/hello": {
          "get": {
            "x-amazon-apigateway-integration": {
              "type": "aws_proxy",
              "httpMethod": "POST",
              "uri": "arn:aws:lambda:us-west-2:123456789012:function:hello"
            }
          }
        }
      }
    }
```

## 🗄️ 存储和数据库

### 1. Amazon S3

#### S3 配置示例
```python
# s3_operations.py
import boto3
from botocore.exceptions import ClientError

class S3Manager:
    def __init__(self, bucket_name, region='us-west-2'):
        self.s3_client = boto3.client('s3', region_name=region)
        self.bucket_name = bucket_name
    
    def upload_file(self, file_path, object_key):
        try:
            self.s3_client.upload_file(file_path, self.bucket_name, object_key)
            print(f"File {file_path} uploaded to s3://{self.bucket_name}/{object_key}")
        except ClientError as e:
            print(f"Error uploading file: {e}")
    
    def download_file(self, object_key, file_path):
        try:
            self.s3_client.download_file(self.bucket_name, object_key, file_path)
            print(f"File s3://{self.bucket_name}/{object_key} downloaded to {file_path}")
        except ClientError as e:
            print(f"Error downloading file: {e}")
    
    def list_objects(self, prefix=''):
        try:
            response = self.s3_client.list_objects_v2(
                Bucket=self.bucket_name,
                Prefix=prefix
            )
            return response.get('Contents', [])
        except ClientError as e:
            print(f"Error listing objects: {e}")
            return []
```

### 2. Amazon RDS

#### RDS 配置示例
```yaml
# rds-instance.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: rds-config
data:
  instance-identifier: my-db-instance
  engine: postgres
  engine-version: "13.7"
  instance-class: db.t3.micro
  allocated-storage: 20
  storage-type: gp2
  master-username: admin
  master-password: mypassword
  vpc-security-group-ids: sg-12345
  db-subnet-group-name: my-db-subnet-group
  backup-retention-period: 7
  multi-az: false
  publicly-accessible: false
```

## 📊 监控和日志

### 1. CloudWatch

#### CloudWatch 配置
```yaml
# cloudwatch-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cloudwatch-config
data:
  log-group: "/aws/eks/my-cluster"
  log-stream: "my-app"
  region: "us-west-2"
  retention-days: "7"
```

#### CloudWatch 指标
```python
# cloudwatch_metrics.py
import boto3
from datetime import datetime

class CloudWatchMetrics:
    def __init__(self, region='us-west-2'):
        self.cloudwatch = boto3.client('cloudwatch', region_name=region)
    
    def put_metric(self, namespace, metric_name, value, unit='Count'):
        try:
            self.cloudwatch.put_metric_data(
                Namespace=namespace,
                MetricData=[
                    {
                        'MetricName': metric_name,
                        'Value': value,
                        'Unit': unit,
                        'Timestamp': datetime.utcnow()
                    }
                ]
            )
            print(f"Metric {metric_name} sent successfully")
        except Exception as e:
            print(f"Error sending metric: {e}")
    
    def get_metric_statistics(self, namespace, metric_name, start_time, end_time):
        try:
            response = self.cloudwatch.get_metric_statistics(
                Namespace=namespace,
                MetricName=metric_name,
                StartTime=start_time,
                EndTime=end_time,
                Period=300,
                Statistics=['Average', 'Sum', 'Maximum']
            )
            return response['Datapoints']
        except Exception as e:
            print(f"Error getting metrics: {e}")
            return []
```

### 2. X-Ray 分布式追踪

#### X-Ray 配置
```python
# xray_tracing.py
from aws_xray_sdk.core import xray_recorder
from aws_xray_sdk.core import patch_all
import boto3

# 启用 X-Ray 追踪
patch_all()

@xray_recorder.capture('my_function')
def my_function():
    # 创建子段
    with xray_recorder.capture('database_query') as subsegment:
        # 模拟数据库查询
        subsegment.put_metadata('query', 'SELECT * FROM users')
        subsegment.put_annotation('table', 'users')
        # 执行查询...
    
    # 创建另一个子段
    with xray_recorder.capture('external_api_call') as subsegment:
        # 模拟外部 API 调用
        subsegment.put_metadata('url', 'https://api.example.com')
        subsegment.put_annotation('service', 'external-api')
        # 执行 API 调用...

# 使用示例
if __name__ == "__main__":
    my_function()
```

## 🔒 安全服务

### 1. IAM 角色配置

#### IAM 策略示例
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::my-bucket/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
```

#### EKS 服务账户配置
```yaml
# iam-service-account.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-west-2

iam:
  withOIDC: true
  serviceAccounts:
  - metadata:
      name: my-app-sa
      namespace: default
    attachPolicyARNs:
    - arn:aws:iam::123456789012:policy/MyAppPolicy
    - arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
```

### 2. VPC 安全配置

#### VPC 配置示例
```yaml
# vpc-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vpc-config
data:
  vpc-cidr: "10.0.0.0/16"
  public-subnets:
    - "10.0.1.0/24"
    - "10.0.2.0/24"
  private-subnets:
    - "10.0.10.0/24"
    - "10.0.20.0/24"
  database-subnets:
    - "10.0.100.0/24"
    - "10.0.200.0/24"
```

## 🛠️ 实践练习

### 练习1: 部署 EKS 集群

```bash
# 创建 EKS 集群
eksctl create cluster \
  --name my-eks-cluster \
  --version 1.24 \
  --region us-west-2 \
  --nodegroup-name workers \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 1 \
  --nodes-max 4

# 部署应用
kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/application/nginx-app.yaml
```

### 练习2: 创建 Lambda 函数

```python
# lambda_function.py
import json
import boto3

def lambda_handler(event, context):
    # 处理不同的触发事件
    if 'Records' in event:
        # S3 事件
        for record in event['Records']:
            bucket = record['s3']['bucket']['name']
            key = record['s3']['object']['key']
            print(f'Processing {key} from {bucket}')
    
    elif 'httpMethod' in event:
        # API Gateway 事件
        return {
            'statusCode': 200,
            'body': json.dumps({'message': 'Hello from Lambda!'})
        }
    
    return {'statusCode': 200}
```

## 📚 相关资源

### 官方文档
- [AWS 官方文档](https://docs.aws.amazon.com/)
- [EKS 用户指南](https://docs.aws.amazon.com/eks/)
- [Lambda 开发者指南](https://docs.aws.amazon.com/lambda/)

### 学习资源
- [AWS 架构中心](https://aws.amazon.com/architecture/)
- [AWS 最佳实践](https://aws.amazon.com/architecture/well-architected/)

### 工具推荐
- **AWS CLI**: 命令行工具
- **eksctl**: EKS 集群管理
- **AWS CDK**: 基础设施即代码
- **Terraform**: 多云基础设施管理

---

**掌握 AWS 云原生服务，构建可扩展的云应用！** 🚀
