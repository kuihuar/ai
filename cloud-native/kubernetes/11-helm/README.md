# Helm 包管理

## 📦 什么是 Helm？

Helm 是 Kubernetes 的包管理工具，类似于 Linux 的 apt 或 yum。它简化了 Kubernetes 应用的部署和管理，通过 Chart 来定义、安装和升级复杂的 Kubernetes 应用。

## 🎯 Helm 特点

### 1. 包管理
- 简化应用部署
- 版本管理
- 依赖管理

### 2. 模板化
- 使用 Go 模板语言
- 支持变量和条件
- 可重用组件

### 3. 生命周期管理
- 安装、升级、回滚
- 卸载和清理
- 状态管理

## 🏗️ Helm 架构

### 1. 核心组件
- **Helm Client**: 命令行工具
- **Tiller Server**: 服务端组件（Helm 3 中已移除）
- **Chart Repository**: 包存储库

### 2. Helm 3 变化
- 移除了 Tiller
- 使用 Kubernetes API 直接操作
- 改进了安全性

## 📁 Chart 结构

### 1. 标准 Chart 目录
```
mychart/
├── Chart.yaml          # Chart 元数据
├── values.yaml         # 默认配置值
├── charts/             # 依赖的 Charts
├── templates/          # 模板文件
│   ├── _helpers.tpl    # 辅助函数
│   ├── deployment.yaml
│   ├── service.yaml
│   └── NOTES.txt       # 安装说明
└── .helmignore         # 忽略文件
```

### 2. Chart.yaml
```yaml
apiVersion: v2
name: myapp
description: A Helm chart for My Application
type: application
version: 0.1.0
appVersion: "1.0.0"
keywords:
  - web
  - application
home: https://github.com/example/myapp
sources:
  - https://github.com/example/myapp
maintainers:
  - name: John Doe
    email: john@example.com
dependencies:
  - name: mysql
    version: 8.0.0
    repository: https://charts.bitnami.com/bitnami
```

### 3. values.yaml
```yaml
# 默认配置值
replicaCount: 1

image:
  repository: nginx
  pullPolicy: IfNotPresent
  tag: ""

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  create: true
  annotations: {}
  name: ""

podAnnotations: {}

podSecurityContext: {}

securityContext: {}

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: ""
  annotations: {}
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls: []

resources: {}

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}
```

## 📝 模板语法

### 1. 基本语法
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "mychart.fullname" . }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "mychart.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "mychart.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
```

### 2. 条件语句
```yaml
{{- if .Values.ingress.enabled -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "mychart.fullname" . }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if .Values.ingress.className }}
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
{{- end }}
```

### 3. 循环语句
```yaml
spec:
  rules:
  {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
        {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType }}
            backend:
              service:
                name: {{ include "mychart.fullname" $ }}
                port:
                  number: {{ $.Values.service.port }}
        {{- end }}
  {{- end }}
```

### 4. 辅助函数
```yaml
{{/*
Expand the name of the chart.
*/}}
{{- define "mychart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "mychart.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mychart.labels" -}}
helm.sh/chart: {{ include "mychart.chart" . }}
{{ include "mychart.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
```

## 🛠️ Helm 命令

### 1. 基础命令
```bash
# 安装 Helm
curl https://get.helm.sh/helm-v3.10.0-linux-amd64.tar.gz | tar xz
sudo mv linux-amd64/helm /usr/local/bin/

# 查看版本
helm version

# 查看帮助
helm help
```

### 2. Repository 管理
```bash
# 添加仓库
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add stable https://charts.helm.sh/stable

# 更新仓库
helm repo update

# 查看仓库
helm repo list

# 搜索 Chart
helm search repo nginx
helm search repo bitnami/nginx
```

### 3. Chart 操作
```bash
# 创建新 Chart
helm create mychart

# 安装 Chart
helm install my-release ./mychart
helm install my-release bitnami/nginx

# 升级 Chart
helm upgrade my-release ./mychart
helm upgrade my-release bitnami/nginx

# 回滚 Chart
helm rollback my-release 1

# 卸载 Chart
helm uninstall my-release

# 查看 Release
helm list
helm status my-release
```

### 4. 配置管理
```bash
# 使用自定义 values
helm install my-release ./mychart -f values.yaml
helm install my-release ./mychart --set replicaCount=3

# 查看生成的 YAML
helm template my-release ./mychart
helm template my-release ./mychart -f values.yaml

# 验证 Chart
helm lint ./mychart
```

## 📦 Chart 开发

### 1. 创建 Chart
```bash
# 创建新 Chart
helm create myapp

# 修改 Chart.yaml
vim myapp/Chart.yaml

# 修改 values.yaml
vim myapp/values.yaml
```

### 2. 开发模板
```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}
  labels:
    {{- include "myapp.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "myapp.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "myapp.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          ports:
            - name: http
              containerPort: {{ .Values.service.port }}
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

### 3. 测试 Chart
```bash
# 本地测试
helm install --dry-run --debug my-release ./myapp

# 验证语法
helm lint ./myapp

# 打包 Chart
helm package ./myapp
```

## 🔄 依赖管理

### 1. 添加依赖
```yaml
# Chart.yaml
dependencies:
  - name: mysql
    version: 8.0.0
    repository: https://charts.bitnami.com/bitnami
  - name: redis
    version: 16.0.0
    repository: https://charts.bitnami.com/bitnami
```

### 2. 管理依赖
```bash
# 更新依赖
helm dependency update ./myapp

# 构建依赖
helm dependency build ./myapp

# 查看依赖
helm dependency list ./myapp
```

## 🎯 最佳实践

### 1. Chart 设计
- 使用有意义的名称
- 提供合理的默认值
- 支持自定义配置

### 2. 模板开发
- 使用辅助函数
- 添加注释说明
- 处理边界情况

### 3. 版本管理
- 遵循语义化版本
- 记录变更日志
- 测试兼容性

### 4. 安全考虑
- 验证用户输入
- 限制资源使用
- 保护敏感信息

## 🛠️ 实践练习

### 练习 1：基础 Chart
1. 创建简单的 Web 应用 Chart
2. 配置部署和服务
3. 测试安装和升级

### 练习 2：复杂应用
1. 创建多组件应用 Chart
2. 管理依赖关系
3. 配置资源限制

### 练习 3：自定义配置
1. 创建灵活的配置选项
2. 支持不同环境
3. 实现条件部署

## 📚 扩展阅读

- [Helm 官方文档](https://helm.sh/docs/)
- [Chart 开发指南](https://helm.sh/docs/chart_template_guide/)
- [最佳实践](https://helm.sh/docs/chart_best_practices/)

## 🎯 下一步

掌握 Helm 后，继续学习：
- [微服务部署](./12-microservices/README.md)
- [CI/CD流水线](./13-cicd/README.md) 