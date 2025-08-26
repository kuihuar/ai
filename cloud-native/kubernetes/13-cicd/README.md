# CI/CD 流水线

## 📖 CI/CD 概述

CI/CD（持续集成/持续部署）是现代软件开发的核心实践。在 Kubernetes 环境中，CI/CD 流水线自动化了代码构建、测试、部署和发布的全过程。

## 🎯 CI/CD 流程

### 1. 持续集成 (CI)
- 代码提交触发构建
- 自动运行测试
- 代码质量检查
- 构建容器镜像

### 2. 持续部署 (CD)
- 自动部署到测试环境
- 自动化测试验证
- 生产环境部署
- 回滚机制

## 🏗️ CI/CD 架构

### 1. 典型架构
```
┌─────────────────────────────────────────────────────────────┐
│                    Git Repository                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   Jenkins   │ │   GitLab    │ │   GitHub    │            │
│  │   CI/CD     │ │   Actions   │ │   Actions   │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │  Container  │ │   Helm      │ │   ArgoCD    │            │
│  │  Registry   │ │   Charts    │ │   GitOps    │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   Dev       │ │   Staging   │ │ Production  │            │
│  │  Cluster    │ │  Cluster    │ │  Cluster    │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### 2. GitOps 模式
- 以 Git 为单一真实源
- 声明式基础设施
- 自动化同步
- 审计和回滚

## 🛠️ Jenkins CI/CD

### 1. Jenkins 部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jenkins
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jenkins
  template:
    metadata:
      labels:
        app: jenkins
    spec:
      serviceAccount: jenkins
      containers:
      - name: jenkins
        image: jenkins/jenkins:lts
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 50000
          name: jnlp
        volumeMounts:
        - name: jenkins-home
          mountPath: /var/jenkins_home
        - name: docker-sock
          mountPath: /var/run/docker.sock
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
      volumes:
      - name: jenkins-home
        persistentVolumeClaim:
          claimName: jenkins-pvc
      - name: docker-sock
        hostPath:
          path: /var/run/docker.sock
---
apiVersion: v1
kind: Service
metadata:
  name: jenkins
spec:
  selector:
    app: jenkins
  ports:
  - port: 8080
    targetPort: 8080
    name: http
  - port: 50000
    targetPort: 50000
    name: jnlp
  type: LoadBalancer
```

### 2. Jenkinsfile
```groovy
pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = 'myapp'
        DOCKER_TAG = "${env.BUILD_NUMBER}"
        REGISTRY = 'registry.example.com'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Build') {
            steps {
                sh 'mvn clean package'
            }
        }
        
        stage('Test') {
            steps {
                sh 'mvn test'
            }
            post {
                always {
                    junit '**/target/surefire-reports/*.xml'
                }
            }
        }
        
        stage('SonarQube') {
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh 'mvn sonar:sonar'
                }
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    docker.build("${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}")
                }
            }
        }
        
        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry("https://${REGISTRY}", 'registry-credentials') {
                        docker.image("${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}").push()
                        docker.image("${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}").push('latest')
                    }
                }
            }
        }
        
        stage('Deploy to Dev') {
            when {
                branch 'develop'
            }
            steps {
                sh "kubectl set image deployment/myapp myapp=${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG} -n dev"
            }
        }
        
        stage('Deploy to Staging') {
            when {
                branch 'main'
            }
            steps {
                sh "kubectl set image deployment/myapp myapp=${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG} -n staging"
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            input {
                message "Deploy to production?"
                ok "Deploy"
            }
            steps {
                sh "kubectl set image deployment/myapp myapp=${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG} -n production"
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}
```

## 🚀 GitHub Actions

### 1. 工作流配置
```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: registry.example.com
  IMAGE_NAME: myapp

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up JDK 11
      uses: actions/setup-java@v3
      with:
        java-version: '11'
        distribution: 'temurin'
    
    - name: Run tests
      run: mvn test
    
    - name: Run SonarQube
      run: mvn sonar:sonar
      env:
        SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/develop'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ secrets.REGISTRY_USERNAME }}
        password: ${{ secrets.REGISTRY_PASSWORD }}
    
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}
          type=sha
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}

  deploy-dev:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    environment: development
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up kubectl
      uses: azure/setup-kubectl@v3
      with:
        version: 'latest'
    
    - name: Configure kubectl
      run: |
        echo "${{ secrets.KUBE_CONFIG_DEV }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
    
    - name: Deploy to development
      run: |
        kubectl set image deployment/myapp myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} -n dev
        kubectl rollout status deployment/myapp -n dev

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: staging
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up kubectl
      uses: azure/setup-kubectl@v3
      with:
        version: 'latest'
    
    - name: Configure kubectl
      run: |
        echo "${{ secrets.KUBE_CONFIG_STAGING }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
    
    - name: Deploy to staging
      run: |
        kubectl set image deployment/myapp myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} -n staging
        kubectl rollout status deployment/myapp -n staging

  deploy-production:
    needs: [build, deploy-staging]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up kubectl
      uses: azure/setup-kubectl@v3
      with:
        version: 'latest'
    
    - name: Configure kubectl
      run: |
        echo "${{ secrets.KUBE_CONFIG_PROD }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
    
    - name: Deploy to production
      run: |
        kubectl set image deployment/myapp myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} -n production
        kubectl rollout status deployment/myapp -n production
```

## 🔄 ArgoCD GitOps

### 1. ArgoCD 部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: argocd-server
  template:
    metadata:
      labels:
        app.kubernetes.io/name: argocd-server
    spec:
      containers:
      - name: argocd-server
        image: quay.io/argoproj/argocd:latest
        ports:
        - containerPort: 8080
        - containerPort: 8083
        env:
        - name: ARGOCD_API_SERVER_REPLICAS
          value: "1"
        - name: ARGOCD_APPLICATION_CONTROLLER_REPLICAS
          value: "1"
        - name: ARGOCD_REPO_SERVER_REPLICAS
          value: "1"
        - name: ARGOCD_REDIS_REPLICAS
          value: "1"
        - name: ARGOCD_APPLICATION_SET_CONTROLLER_REPLICAS
          value: "1"
---
apiVersion: v1
kind: Service
metadata:
  name: argocd-server
spec:
  selector:
    app.kubernetes.io/name: argocd-server
  ports:
  - port: 80
    targetPort: 8080
    name: http
  - port: 443
    targetPort: 8083
    name: https
  type: LoadBalancer
```

### 2. Application 定义
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/example/myapp
    targetRevision: HEAD
    path: k8s
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
  revisionHistoryLimit: 10
```

## 🏗️ Helm Charts

### 1. Chart 结构
```
myapp/
├── Chart.yaml
├── values.yaml
├── values-dev.yaml
├── values-staging.yaml
├── values-prod.yaml
└── templates/
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    └── configmap.yaml
```

### 2. 多环境配置
```yaml
# values-dev.yaml
replicaCount: 1
image:
  repository: myapp
  tag: latest
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"
ingress:
  enabled: false
```

```yaml
# values-prod.yaml
replicaCount: 3
image:
  repository: myapp
  tag: stable
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: nginx
  hosts:
    - host: myapp.example.com
      paths:
        - path: /
          pathType: Prefix
```

## 🎯 最佳实践

### 1. 流水线设计
- 快速反馈
- 自动化测试
- 渐进式部署
- 回滚机制

### 2. 安全考虑
- 镜像扫描
- 密钥管理
- 访问控制
- 审计日志

### 3. 监控告警
- 部署状态监控
- 应用性能监控
- 错误率告警
- 回滚告警

### 4. 环境管理
- 环境隔离
- 配置管理
- 数据管理
- 资源管理

## 🛠️ 实践练习

### 练习 1：Jenkins 流水线
1. 部署 Jenkins
2. 配置流水线
3. 测试自动化部署

### 练习 2：GitHub Actions
1. 配置工作流
2. 集成容器注册表
3. 自动化部署

### 练习 3：ArgoCD GitOps
1. 部署 ArgoCD
2. 配置应用
3. 测试 GitOps 流程

## 📚 扩展阅读

- [Jenkins 官方文档](https://www.jenkins.io/doc/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [ArgoCD 官方文档](https://argo-cd.readthedocs.io/)

## 🎯 下一步

掌握 CI/CD 后，继续学习：
- [故障排查](./14-troubleshooting/README.md)
- [性能优化](./15-optimization/README.md) 