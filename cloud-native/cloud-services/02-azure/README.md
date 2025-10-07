# Azure 云原生服务详解

## 📚 学习目标

通过本模块学习，您将掌握：
- Azure 核心云原生服务架构
- 容器化服务 AKS、Container Instances
- 无服务器服务 Functions、Logic Apps
- 存储和数据库服务
- 监控和日志服务
- 安全服务和最佳实践

## 🎯 Azure 云原生架构

### 1. Azure 服务生态

```
Azure 云原生服务
├── 计算服务
│   ├── Virtual Machines (VM)
│   ├── AKS (Azure Kubernetes Service)
│   ├── Container Instances (ACI)
│   ├── App Service
│   └── Functions (Serverless)
├── 存储服务
│   ├── Blob Storage
│   ├── Managed Disks
│   ├── Files
│   └── Data Lake Storage
├── 数据库服务
│   ├── SQL Database
│   ├── Cosmos DB
│   ├── Redis Cache
│   └── Database for PostgreSQL
├── 网络服务
│   ├── Virtual Network (VNet)
│   ├── Load Balancer
│   ├── CDN
│   └── DNS
└── 监控服务
    ├── Monitor
    ├── Application Insights
    ├── Log Analytics
    └── Security Center
```

### 2. 服务对比

| 服务 | 类型 | 用途 | 优势 |
|------|------|------|------|
| AKS | Kubernetes | 容器编排 | 完全托管，标准 K8s |
| Container Instances | 容器运行 | 快速部署 | 按需付费，无基础设施管理 |
| App Service | 应用托管 | Web 应用 | 自动扩展，多语言支持 |
| Functions | 无服务器函数 | 事件驱动 | 按需付费，自动扩展 |

## 🐳 容器化服务

### 1. Azure Kubernetes Service (AKS)

#### AKS 集群创建
```bash
# 创建资源组
az group create --name myResourceGroup --location eastus

# 创建 AKS 集群
az aks create \
  --resource-group myResourceGroup \
  --name myAKSCluster \
  --node-count 3 \
  --node-vm-size Standard_B2s \
  --enable-addons monitoring \
  --generate-ssh-keys

# 获取凭据
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster
```

#### AKS 配置示例
```yaml
# aks-cluster.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aks-config
data:
  cluster-name: myAKSCluster
  resource-group: myResourceGroup
  location: eastus
  node-count: "3"
  node-vm-size: Standard_B2s
  kubernetes-version: "1.24.0"
  enable-monitoring: "true"
  enable-rbac: "true"
```

### 2. Azure Container Instances (ACI)

#### ACI 配置示例
```yaml
# aci-deployment.yaml
apiVersion: v1
kind: Pod
metadata:
  name: aci-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    ports:
    - containerPort: 80
  nodeSelector:
    kubernetes.io/arch: amd64
  tolerations:
  - key: "azure.com/aci"
    operator: "Equal"
    value: "true"
    effect: "NoSchedule"
```

#### ACI 部署脚本
```bash
# 创建 ACI 实例
az container create \
  --resource-group myResourceGroup \
  --name mycontainer \
  --image nginx:latest \
  --cpu 1 \
  --memory 1 \
  --ports 80 \
  --dns-name-label myapp \
  --location eastus

# 查看状态
az container show --resource-group myResourceGroup --name mycontainer
```

### 3. Azure App Service

#### App Service 配置
```yaml
# app-service.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-service-config
data:
  app-name: my-web-app
  resource-group: myResourceGroup
  location: eastus
  sku: F1
  runtime: "DOTNETCORE|3.1"
  always-on: "true"
  https-only: "true"
```

#### App Service 部署
```bash
# 创建 App Service 计划
az appservice plan create \
  --name myAppServicePlan \
  --resource-group myResourceGroup \
  --sku F1

# 创建 Web 应用
az webapp create \
  --resource-group myResourceGroup \
  --plan myAppServicePlan \
  --name myWebApp \
  --runtime "DOTNETCORE|3.1"

# 部署代码
az webapp deployment source config \
  --resource-group myResourceGroup \
  --name myWebApp \
  --repo-url https://github.com/Azure-Samples/dotnetcore-api \
  --branch main \
  --manual-integration
```

## ⚡ 无服务器服务

### 1. Azure Functions

#### Functions 项目结构
```
MyFunctionApp/
├── host.json
├── local.settings.json
├── requirements.txt
└── HttpTrigger/
    ├── function.json
    └── __init__.py
```

#### Functions 代码示例
```python
# HttpTrigger/__init__.py
import logging
import azure.functions as func

def main(req: func.HttpRequest) -> func.HttpResponse:
    logging.info('Python HTTP trigger function processed a request.')

    name = req.params.get('name')
    if not name:
        try:
            req_body = req.get_json()
        except ValueError:
            pass
        else:
            name = req_body.get('name')

    if name:
        return func.HttpResponse(f"Hello, {name}. This HTTP triggered function executed successfully.")
    else:
        return func.HttpResponse(
             "This HTTP triggered function executed successfully. Pass a name in the query string or in the request body for a personalized response.",
             status_code=200
        )
```

#### Functions 配置
```json
// host.json
{
  "version": "2.0",
  "logging": {
    "applicationInsights": {
      "samplingSettings": {
        "isEnabled": true,
        "excludedTypes": "Request"
      }
    }
  },
  "extensionBundle": {
    "id": "Microsoft.Azure.Functions.ExtensionBundle",
    "version": "[2.*, 3.0.0)"
  }
}
```

```json
// HttpTrigger/function.json
{
  "scriptFile": "__init__.py",
  "bindings": [
    {
      "authLevel": "function",
      "type": "httpTrigger",
      "direction": "in",
      "name": "req",
      "methods": [
        "get",
        "post"
      ]
    },
    {
      "type": "http",
      "direction": "out",
      "name": "$return"
    }
  ]
}
```

### 2. Logic Apps

#### Logic App 配置
```json
{
  "definition": {
    "$schema": "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
    "contentVersion": "1.0.0.0",
    "parameters": {},
    "triggers": {
      "When_a_HTTP_request_is_received": {
        "type": "Request",
        "kind": "Http",
        "inputs": {
          "schema": {}
        }
      }
    },
    "actions": {
      "Response": {
        "type": "Response",
        "kind": "Http",
        "inputs": {
          "statusCode": 200,
          "body": "Hello from Logic App!"
        }
      }
    },
    "outputs": {}
  }
}
```

## 🗄️ 存储和数据库

### 1. Azure Blob Storage

#### Blob Storage 操作
```python
# blob_storage.py
from azure.storage.blob import BlobServiceClient, BlobClient, ContainerClient
import os

class BlobStorageManager:
    def __init__(self, connection_string):
        self.blob_service_client = BlobServiceClient.from_connection_string(connection_string)
    
    def create_container(self, container_name):
        try:
            self.blob_service_client.create_container(container_name)
            print(f"Container {container_name} created successfully")
        except Exception as e:
            print(f"Error creating container: {e}")
    
    def upload_blob(self, container_name, blob_name, file_path):
        try:
            blob_client = self.blob_service_client.get_blob_client(
                container=container_name, 
                blob=blob_name
            )
            with open(file_path, "rb") as data:
                blob_client.upload_blob(data)
            print(f"File {file_path} uploaded to {blob_name}")
        except Exception as e:
            print(f"Error uploading blob: {e}")
    
    def download_blob(self, container_name, blob_name, download_path):
        try:
            blob_client = self.blob_service_client.get_blob_client(
                container=container_name, 
                blob=blob_name
            )
            with open(download_path, "wb") as download_file:
                download_file.write(blob_client.download_blob().readall())
            print(f"Blob {blob_name} downloaded to {download_path}")
        except Exception as e:
            print(f"Error downloading blob: {e}")
    
    def list_blobs(self, container_name):
        try:
            container_client = self.blob_service_client.get_container_client(container_name)
            blobs = container_client.list_blobs()
            return list(blobs)
        except Exception as e:
            print(f"Error listing blobs: {e}")
            return []
```

### 2. Azure Cosmos DB

#### Cosmos DB 配置
```python
# cosmos_db.py
from azure.cosmos import CosmosClient, PartitionKey
import json

class CosmosDBManager:
    def __init__(self, endpoint, key, database_name):
        self.client = CosmosClient(endpoint, key)
        self.database = self.client.get_database_client(database_name)
    
    def create_container(self, container_name, partition_key):
        try:
            container = self.database.create_container(
                id=container_name,
                partition_key=PartitionKey(path=partition_key)
            )
            print(f"Container {container_name} created successfully")
            return container
        except Exception as e:
            print(f"Error creating container: {e}")
            return None
    
    def create_item(self, container_name, item):
        try:
            container = self.database.get_container_client(container_name)
            container.create_item(item)
            print(f"Item created successfully")
        except Exception as e:
            print(f"Error creating item: {e}")
    
    def read_item(self, container_name, item_id, partition_key):
        try:
            container = self.database.get_container_client(container_name)
            item = container.read_item(item_id, partition_key)
            return item
        except Exception as e:
            print(f"Error reading item: {e}")
            return None
    
    def query_items(self, container_name, query):
        try:
            container = self.database.get_container_client(container_name)
            items = list(container.query_items(query, enable_cross_partition_query=True))
            return items
        except Exception as e:
            print(f"Error querying items: {e}")
            return []
```

## 📊 监控和日志

### 1. Azure Monitor

#### Monitor 配置
```yaml
# monitor-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: monitor-config
data:
  workspace-id: "your-workspace-id"
  workspace-key: "your-workspace-key"
  log-analytics-workspace: "your-workspace"
  region: "eastus"
```

#### 自定义指标
```python
# custom_metrics.py
from azure.monitor.query import MetricsQueryClient
from azure.identity import DefaultAzureCredential
from datetime import datetime, timedelta

class CustomMetrics:
    def __init__(self, workspace_id):
        self.credential = DefaultAzureCredential()
        self.client = MetricsQueryClient(self.credential)
        self.workspace_id = workspace_id
    
    def query_metrics(self, metric_name, start_time, end_time):
        try:
            response = self.client.query_resource(
                resource_uri=f"/subscriptions/{self.workspace_id}",
                metric_names=[metric_name],
                start_time=start_time,
                end_time=end_time,
                granularity=timedelta(minutes=5)
            )
            return response
        except Exception as e:
            print(f"Error querying metrics: {e}")
            return None
```

### 2. Application Insights

#### Application Insights 配置
```python
# app_insights.py
from opencensus.ext.azure.log_exporter import AzureLogHandler
from opencensus.ext.azure.trace_exporter import AzureExporter
from opencensus.trace import config_integration
from opencensus.trace.samplers import ProbabilitySampler
from opencensus.trace.tracer import Tracer
import logging

class AppInsightsManager:
    def __init__(self, connection_string):
        self.connection_string = connection_string
        self.setup_logging()
        self.setup_tracing()
    
    def setup_logging(self):
        logger = logging.getLogger(__name__)
        logger.addHandler(AzureLogHandler(connection_string=self.connection_string))
        logger.setLevel(logging.INFO)
        self.logger = logger
    
    def setup_tracing(self):
        config_integration.trace_integrations(['logging'])
        exporter = AzureExporter(connection_string=self.connection_string)
        tracer = Tracer(exporter=exporter, sampler=ProbabilitySampler(1.0))
        self.tracer = tracer
    
    def log_info(self, message):
        self.logger.info(message)
    
    def log_error(self, message):
        self.logger.error(message)
    
    def trace_function(self, func_name):
        with self.tracer.span(name=func_name):
            # 函数逻辑
            pass
```

## 🔒 安全服务

### 1. Azure Active Directory (AAD)

#### AAD 应用注册
```yaml
# aad-app.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aad-app-config
data:
  client-id: "your-client-id"
  client-secret: "your-client-secret"
  tenant-id: "your-tenant-id"
  redirect-uri: "https://your-app.com/auth/callback"
  scope: "https://graph.microsoft.com/.default"
```

#### AAD 认证示例
```python
# aad_auth.py
from msal import ConfidentialClientApplication
import requests

class AADAuth:
    def __init__(self, client_id, client_secret, tenant_id):
        self.client_id = client_id
        self.client_secret = client_secret
        self.tenant_id = tenant_id
        self.authority = f"https://login.microsoftonline.com/{tenant_id}"
        self.scope = ["https://graph.microsoft.com/.default"]
        
        self.app = ConfidentialClientApplication(
            client_id=client_id,
            client_credential=client_secret,
            authority=self.authority
        )
    
    def get_access_token(self):
        try:
            result = self.app.acquire_token_silent(self.scope, account=None)
            if not result:
                result = self.app.acquire_token_for_client(scopes=self.scope)
            
            if "access_token" in result:
                return result["access_token"]
            else:
                print(f"Error: {result.get('error')}")
                return None
        except Exception as e:
            print(f"Error getting access token: {e}")
            return None
```

### 2. Key Vault

#### Key Vault 配置
```python
# key_vault.py
from azure.keyvault.secrets import SecretClient
from azure.identity import DefaultAzureCredential

class KeyVaultManager:
    def __init__(self, vault_url):
        self.credential = DefaultAzureCredential()
        self.client = SecretClient(vault_url=vault_url, credential=self.credential)
    
    def set_secret(self, secret_name, secret_value):
        try:
            self.client.set_secret(secret_name, secret_value)
            print(f"Secret {secret_name} set successfully")
        except Exception as e:
            print(f"Error setting secret: {e}")
    
    def get_secret(self, secret_name):
        try:
            secret = self.client.get_secret(secret_name)
            return secret.value
        except Exception as e:
            print(f"Error getting secret: {e}")
            return None
    
    def delete_secret(self, secret_name):
        try:
            self.client.begin_delete_secret(secret_name)
            print(f"Secret {secret_name} deleted successfully")
        except Exception as e:
            print(f"Error deleting secret: {e}")
```

## 🛠️ 实践练习

### 练习1: 部署 AKS 集群

```bash
# 创建资源组
az group create --name myResourceGroup --location eastus

# 创建 AKS 集群
az aks create \
  --resource-group myResourceGroup \
  --name myAKSCluster \
  --node-count 3 \
  --node-vm-size Standard_B2s \
  --enable-addons monitoring

# 获取凭据
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster

# 部署应用
kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/application/nginx-app.yaml
```

### 练习2: 创建 Azure Function

```python
# function_app.py
import azure.functions as func
import logging

app = func.FunctionApp()

@app.function_name(name="HttpTrigger1")
@app.route(route="hello")
def test_function(req: func.HttpRequest) -> func.HttpResponse:
    logging.info('Python HTTP trigger function processed a request.')

    name = req.params.get('name')
    if not name:
        try:
            req_body = req.get_json()
        except ValueError:
            pass
        else:
            name = req_body.get('name')

    if name:
        return func.HttpResponse(f"Hello, {name}. This HTTP triggered function executed successfully.")
    else:
        return func.HttpResponse(
             "This HTTP triggered function executed successfully. Pass a name in the query string or in the request body for a personalized response.",
             status_code=200
        )
```

## 📚 相关资源

### 官方文档
- [Azure 官方文档](https://docs.microsoft.com/azure/)
- [AKS 用户指南](https://docs.microsoft.com/azure/aks/)
- [Functions 开发者指南](https://docs.microsoft.com/azure/azure-functions/)

### 学习资源
- [Azure 架构中心](https://docs.microsoft.com/azure/architecture/)
- [Azure 最佳实践](https://docs.microsoft.com/azure/architecture/framework/)

### 工具推荐
- **Azure CLI**: 命令行工具
- **Azure PowerShell**: PowerShell 模块
- **Azure Portal**: Web 管理界面
- **Azure Data Studio**: 数据库管理工具

---

**掌握 Azure 云原生服务，构建现代化的云应用！** 🚀
