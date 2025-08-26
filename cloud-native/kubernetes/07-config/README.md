# ConfigMap 与 Secret

## 📖 什么是 ConfigMap？

ConfigMap 是 Kubernetes 中用于存储非敏感配置数据的资源对象。它可以将配置数据与容器镜像分离，使应用程序更加灵活和可移植。

## 🎯 ConfigMap 特点

### 1. 配置分离
- 将配置从容器镜像中分离
- 支持不同环境的配置管理
- 便于配置的版本控制

### 2. 多种数据格式
- 支持键值对、文件、目录
- 支持 YAML、JSON、纯文本
- 支持二进制数据

### 3. 动态更新
- 支持配置的热更新
- 无需重启容器即可更新配置
- 支持配置的回滚

## 📝 ConfigMap 配置

### 1. 从命令行创建
```bash
# 从字面量创建
kubectl create configmap app-config --from-literal=APP_ENV=production --from-literal=LOG_LEVEL=info

# 从文件创建
kubectl create configmap nginx-config --from-file=nginx.conf

# 从目录创建
kubectl create configmap app-config --from-file=config/
```

### 2. 从 YAML 文件创建
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  # 键值对
  APP_ENV: "production"
  LOG_LEVEL: "info"
  DATABASE_URL: "mysql://localhost:3306/app"
  
  # 配置文件
  nginx.conf: |
    server {
        listen 80;
        server_name localhost;
        
        location / {
            root /usr/share/nginx/html;
            index index.html;
        }
    }
  
  # JSON 配置
  app.json: |
    {
      "database": {
        "host": "localhost",
        "port": 3306,
        "name": "app"
      },
      "redis": {
        "host": "localhost",
        "port": 6379
      }
    }
```

## 🔐 什么是 Secret？

Secret 是 Kubernetes 中用于存储敏感数据的资源对象，如密码、令牌、密钥等。Secret 数据以 base64 编码存储，提供了一定程度的安全性。

## 🎯 Secret 特点

### 1. 敏感数据管理
- 存储密码、令牌、密钥等敏感信息
- 支持多种类型的敏感数据
- 提供访问控制机制

### 2. 数据编码
- 数据以 base64 编码存储
- 支持二进制数据
- 提供数据加密选项

### 3. 类型支持
- **Opaque**：通用类型
- **kubernetes.io/service-account-token**：服务账户令牌
- **kubernetes.io/dockercfg**：Docker 配置
- **kubernetes.io/tls**：TLS 证书

## 📝 Secret 配置

### 1. 从命令行创建
```bash
# 从字面量创建
kubectl create secret generic db-secret --from-literal=username=admin --from-literal=password=secret123

# 从文件创建
kubectl create secret generic tls-secret --from-file=tls.crt --from-file=tls.key

# 从 Docker 配置创建
kubectl create secret docker-registry regcred --docker-server=<your-registry-server> --docker-username=<your-username> --docker-password=<your-password>
```

### 2. 从 YAML 文件创建
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
  namespace: default
type: Opaque
data:
  # base64 编码的数据
  username: YWRtaW4=  # admin
  password: c2VjcmV0MTIz  # secret123
  database-url: bXlzcWw6Ly9sb2NhbGhvc3Q6MzMwNi9hcHA=
```

### 3. TLS Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tls-secret
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-certificate>
  tls.key: <base64-encoded-private-key>
```

## 🔧 在 Pod 中使用

### 1. 环境变量方式
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  containers:
  - name: app
    image: myapp:latest
    env:
    # 从 ConfigMap 获取
    - name: APP_ENV
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: APP_ENV
    - name: LOG_LEVEL
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: LOG_LEVEL
    # 从 Secret 获取
    - name: DB_USERNAME
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: password
```

### 2. 文件挂载方式
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
    volumeMounts:
    # 挂载 ConfigMap
    - name: nginx-config
      mountPath: /etc/nginx/conf.d
      readOnly: true
    # 挂载 Secret
    - name: tls-secret
      mountPath: /etc/nginx/ssl
      readOnly: true
  volumes:
  # ConfigMap 卷
  - name: nginx-config
    configMap:
      name: nginx-config
  # Secret 卷
  - name: tls-secret
    secret:
      secretName: tls-secret
```

### 3. 子路径挂载
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  containers:
  - name: app
    image: myapp:latest
    volumeMounts:
    - name: config-volume
      mountPath: /app/config/app.json
      subPath: app.json
  volumes:
  - name: config-volume
    configMap:
      name: app-config
```

## 🛠️ 常用操作

### 1. 创建和查看
```bash
# 创建 ConfigMap
kubectl create configmap app-config --from-literal=APP_ENV=production

# 创建 Secret
kubectl create secret generic db-secret --from-literal=username=admin --from-literal=password=secret123

# 查看 ConfigMap
kubectl get configmaps
kubectl describe configmap app-config

# 查看 Secret
kubectl get secrets
kubectl describe secret db-secret
```

### 2. 更新配置
```bash
# 更新 ConfigMap
kubectl patch configmap app-config -p '{"data":{"APP_ENV":"staging"}}'

# 更新 Secret
kubectl patch secret db-secret -p '{"data":{"password":"bmV3cGFzc3dvcmQ="}}'
```

### 3. 删除配置
```bash
# 删除 ConfigMap
kubectl delete configmap app-config

# 删除 Secret
kubectl delete secret db-secret
```

## 🔄 配置更新策略

### 1. 自动更新
- ConfigMap 和 Secret 更新后，挂载的卷会自动更新
- 应用程序需要支持配置热重载
- 某些情况下可能需要重启 Pod

### 2. 滚动更新
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  template:
    metadata:
      annotations:
        checksum/config: "${CONFIG_CHECKSUM}"
    spec:
      containers:
      - name: app
        image: myapp:latest
```

## 🎯 最佳实践

### 1. 命名规范
- 使用有意义的名称
- 遵循命名空间约定
- 使用标签进行分类

### 2. 数据管理
- 避免在 ConfigMap 中存储敏感数据
- 使用 Secret 存储敏感信息
- 定期轮换敏感数据

### 3. 访问控制
- 使用 RBAC 控制访问权限
- 限制 Secret 的访问范围
- 监控配置访问日志

### 4. 版本管理
- 使用版本控制管理配置
- 支持配置回滚
- 记录配置变更历史

## 🛠️ 实践练习

### 练习 1：基础配置管理
1. 创建 ConfigMap 存储应用配置
2. 在 Pod 中使用环境变量
3. 测试配置更新

### 练习 2：文件配置
1. 创建包含配置文件的 ConfigMap
2. 挂载到 Pod 中
3. 测试配置热更新

### 练习 3：敏感数据管理
1. 创建 Secret 存储数据库凭据
2. 在应用中使用 Secret
3. 测试安全访问

## 📚 扩展阅读

- [Kubernetes ConfigMap 官方文档](https://kubernetes.io/docs/concepts/configuration/configmap/)
- [Kubernetes Secret 官方文档](https://kubernetes.io/docs/concepts/configuration/secret/)
- [配置管理最佳实践](https://kubernetes.io/docs/concepts/configuration/overview/)

## 🎯 下一步

掌握配置管理后，继续学习：
- [存储管理](./08-storage/README.md)
- [安全机制](./09-security/README.md) 