# 云服务提供商服务发现

## 📖 概述

云服务提供商（如 AWS、Azure、GCP）提供了原生的服务发现解决方案，与云基础设施深度集成，提供高可用性和可扩展性。

## 🎯 核心特性

### 1. 云原生集成
- 与云基础设施深度集成
- 自动扩缩容
- 高可用性

### 2. 多区域支持
- 跨区域服务发现
- 地理分布
- 故障转移

### 3. 安全集成
- IAM 集成
- VPC 网络隔离
- 加密传输

## 🚀 快速开始

### 1. AWS ECS 服务发现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/servicediscovery"
    "google.golang.org/grpc"
    pb "your-project/proto"
)

// AWSServiceDiscovery AWS 服务发现
type AWSServiceDiscovery struct {
    client *servicediscovery.ServiceDiscovery
    namespace string
    serviceName string
}

// NewAWSServiceDiscovery 创建 AWS 服务发现
func NewAWSServiceDiscovery(region, namespace, serviceName string) (*AWSServiceDiscovery, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
    })
    if err != nil {
        return nil, err
    }
    
    client := servicediscovery.New(sess)
    
    return &AWSServiceDiscovery{
        client: client,
        namespace: namespace,
        serviceName: serviceName,
    }, nil
}

// Discover 发现服务
func (asd *AWSServiceDiscovery) Discover() ([]string, error) {
    input := &servicediscovery.DiscoverInstancesInput{
        NamespaceName: aws.String(asd.namespace),
        ServiceName:   aws.String(asd.serviceName),
    }
    
    result, err := asd.client.DiscoverInstances(input)
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, instance := range result.Instances {
        if instance.Attributes["AWS_INSTANCE_IPV4"] != nil {
            ip := *instance.Attributes["AWS_INSTANCE_IPV4"]
            port := "8080"
            if instance.Attributes["AWS_INSTANCE_PORT"] != nil {
                port = *instance.Attributes["AWS_INSTANCE_PORT"]
            }
            addresses = append(addresses, fmt.Sprintf("%s:%s", ip, port))
        }
    }
    
    return addresses, nil
}

// Register 注册服务
func (asd *AWSServiceDiscovery) Register(instanceID, ip, port string, attributes map[string]string) error {
    if attributes == nil {
        attributes = make(map[string]string)
    }
    
    attributes["AWS_INSTANCE_IPV4"] = ip
    attributes["AWS_INSTANCE_PORT"] = port
    
    input := &servicediscovery.RegisterInstanceInput{
        InstanceId:    aws.String(instanceID),
        ServiceId:     aws.String(asd.serviceName),
        Attributes:    aws.StringMap(attributes),
    }
    
    _, err := asd.client.RegisterInstance(input)
    return err
}

// Deregister 注销服务
func (asd *AWSServiceDiscovery) Deregister(instanceID string) error {
    input := &servicediscovery.DeregisterInstanceInput{
        InstanceId: aws.String(instanceID),
        ServiceId:  aws.String(asd.serviceName),
    }
    
    _, err := asd.client.DeregisterInstance(input)
    return err
}
```

### 2. Azure Service Fabric 服务发现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicefabric/armservicefabric"
    "google.golang.org/grpc"
    pb "your-project/proto"
)

// AzureServiceDiscovery Azure 服务发现
type AzureServiceDiscovery struct {
    client *armservicefabric.ClustersClient
    clusterName string
    resourceGroup string
    subscriptionID string
}

// NewAzureServiceDiscovery 创建 Azure 服务发现
func NewAzureServiceDiscovery(subscriptionID, resourceGroup, clusterName string) (*AzureServiceDiscovery, error) {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        return nil, err
    }
    
    client, err := armservicefabric.NewClustersClient(subscriptionID, cred, nil)
    if err != nil {
        return nil, err
    }
    
    return &AzureServiceDiscovery{
        client: client,
        clusterName: clusterName,
        resourceGroup: resourceGroup,
        subscriptionID: subscriptionID,
    }, nil
}

// Discover 发现服务
func (asd *AzureServiceDiscovery) Discover(serviceName string) ([]string, error) {
    // 这里需要根据具体的 Azure Service Fabric 实现
    // 通常通过 REST API 或 SDK 查询服务实例
    return []string{}, nil
}
```

### 3. GCP Cloud Run 服务发现

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "cloud.google.com/go/run/apiv2"
    "cloud.google.com/go/run/apiv2/runpb"
    "google.golang.org/api/option"
    "google.golang.org/grpc"
    pb "your-project/proto"
)

// GCPServiceDiscovery GCP 服务发现
type GCPServiceDiscovery struct {
    client *run.ServicesClient
    projectID string
    region string
}

// NewGCPServiceDiscovery 创建 GCP 服务发现
func NewGCPServiceDiscovery(projectID, region string) (*GCPServiceDiscovery, error) {
    ctx := context.Background()
    client, err := run.NewServicesClient(ctx, option.WithCredentialsFile("path/to/service-account.json"))
    if err != nil {
        return nil, err
    }
    
    return &GCPServiceDiscovery{
        client: client,
        projectID: projectID,
        region: region,
    }, nil
}

// Discover 发现服务
func (gsd *GCPServiceDiscovery) Discover(serviceName string) ([]string, error) {
    ctx := context.Background()
    
    req := &runpb.GetServiceRequest{
        Name: fmt.Sprintf("projects/%s/locations/%s/services/%s", gsd.projectID, gsd.region, serviceName),
    }
    
    service, err := gsd.client.GetService(ctx, req)
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    if service.Status.Url != "" {
        addresses = append(addresses, service.Status.Url)
    }
    
    return addresses, nil
}
```

## 📝 使用示例

### 1. AWS ECS 部署

```yaml
# aws-ecs-task-definition.json
{
  "family": "your-service",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::123456789012:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "your-service",
      "image": "your-service:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/your-service",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### 2. Azure Service Fabric 部署

```yaml
# azure-service-fabric-application.yaml
apiVersion: servicefabric.azure.com/v1beta1
kind: Application
metadata:
  name: your-service
  namespace: default
spec:
  applicationTypeName: YourServiceType
  applicationTypeVersion: "1.0.0"
  parameters:
    - name: InstanceCount
      value: "3"
    - name: Port
      value: "8080"
```

### 3. GCP Cloud Run 部署

```yaml
# gcp-cloud-run-service.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: your-service
  namespace: default
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "10"
        autoscaling.knative.dev/minScale: "1"
    spec:
      containers:
      - image: gcr.io/your-project/your-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        resources:
          limits:
            cpu: "1000m"
            memory: "512Mi"
```

## 🔧 高级配置

### 1. AWS ECS 服务发现配置

```yaml
# aws-ecs-service-discovery.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-ecs-config
data:
  service-discovery.yaml: |
    namespace: your-namespace
    service: your-service
    region: us-west-2
    cluster: your-cluster
    task-definition: your-task-definition
    desired-count: 3
    launch-type: FARGATE
    network-configuration:
      awsvpc-configuration:
        subnets:
          - subnet-12345678
          - subnet-87654321
        security-groups:
          - sg-12345678
        assign-public-ip: ENABLED
```

### 2. Azure Service Fabric 配置

```yaml
# azure-service-fabric-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: azure-sf-config
data:
  service-fabric.yaml: |
    cluster-name: your-cluster
    resource-group: your-resource-group
    subscription-id: your-subscription-id
    location: eastus
    node-count: 3
    vm-size: Standard_D2s_v3
    os-type: Linux
```

### 3. GCP Cloud Run 配置

```yaml
# gcp-cloud-run-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gcp-cloud-run-config
data:
  cloud-run.yaml: |
    project-id: your-project-id
    region: us-central1
    service-name: your-service
    image: gcr.io/your-project/your-service:latest
    port: 8080
    cpu: "1000m"
    memory: "512Mi"
    min-instances: 1
    max-instances: 10
```

## 📊 性能优化

### 1. 连接池配置

```go
// 配置 gRPC 连接池
func createGRPCConnectionWithPool(serviceAddr string) (*grpc.ClientConn, error) {
    return grpc.Dial(serviceAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
        Time:                10 * time.Second,
        Timeout:             3 * time.Second,
        PermitWithoutStream: true,
    }))
}
```

### 2. 负载均衡

```go
// 使用云提供商的负载均衡
func createGRPCConnectionWithLB(serviceAddr string) (*grpc.ClientConn, error) {
    return grpc.Dial(serviceAddr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
}
```

### 3. 缓存优化

```go
// 使用缓存减少 API 调用
type CachedCloudServiceDiscovery struct {
    *AWSServiceDiscovery
    cache map[string][]string
    mutex sync.RWMutex
}

func (ccsd *CachedCloudServiceDiscovery) Discover() ([]string, error) {
    ccsd.mutex.RLock()
    if addresses, ok := ccsd.cache["instances"]; ok {
        ccsd.mutex.RUnlock()
        return addresses, nil
    }
    ccsd.mutex.RUnlock()
    
    addresses, err := ccsd.AWSServiceDiscovery.Discover()
    if err != nil {
        return nil, err
    }
    
    ccsd.mutex.Lock()
    ccsd.cache["instances"] = addresses
    ccsd.mutex.Unlock()
    
    return addresses, nil
}
```

## ��️ 安全配置

### 1. IAM 配置

```yaml
# aws-iam-policy.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-iam-policy
data:
  iam-policy.json: |
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "servicediscovery:DiscoverInstances",
            "servicediscovery:RegisterInstance",
            "servicediscovery:DeregisterInstance"
          ],
          "Resource": "*"
        }
      ]
    }
```

### 2. VPC 配置

```yaml
# aws-vpc-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-vpc-config
data:
  vpc-config.yaml: |
    vpc-id: vpc-12345678
    subnets:
      - subnet-12345678
      - subnet-87654321
    security-groups:
      - sg-12345678
    route-tables:
      - rt-12345678
```

## 🔍 监控和调试

### 1. CloudWatch 监控

```yaml
# aws-cloudwatch-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-cloudwatch-config
data:
  cloudwatch.yaml: |
    namespace: your-service
    metrics:
      - name: RequestCount
        unit: Count
        value: 1
      - name: ResponseTime
        unit: Milliseconds
        value: 100
    alarms:
      - name: HighErrorRate
        metric: ErrorRate
        threshold: 5
        comparison: GreaterThanThreshold
        period: 300
        evaluation-periods: 2
```

### 2. Azure Monitor 配置

```yaml
# azure-monitor-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: azure-monitor-config
data:
  monitor.yaml: |
    workspace-id: your-workspace-id
    workspace-key: your-workspace-key
    metrics:
      - name: RequestCount
        unit: Count
        value: 1
      - name: ResponseTime
        unit: Milliseconds
        value: 100
```

## 📚 最佳实践

1. **云原生设计**: 充分利用云提供商的特性
2. **安全配置**: 实施严格的安全策略
3. **监控告警**: 监控服务状态和性能
4. **成本优化**: 优化资源使用和成本
5. **故障恢复**: 配置合适的故障恢复策略

## 🔗 相关资源

- [AWS ECS 服务发现](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-discovery.html)
- [Azure Service Fabric](https://docs.microsoft.com/en-us/azure/service-fabric/)
- [GCP Cloud Run](https://cloud.google.com/run/docs)
- [云原生服务发现最佳实践](https://cloud.google.com/architecture/best-practices-for-operating-containers)
