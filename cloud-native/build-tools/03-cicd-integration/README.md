# CI/CD 集成构建

## 📚 学习目标

通过本模块学习，您将掌握：
- CI/CD 流水线中的构建集成策略
- 构建缓存和依赖管理优化
- 多环境构建和部署策略
- 构建监控和告警机制
- 企业级构建流水线设计

## 🎯 CI/CD 构建概览

### 1. CI/CD 构建流程

```
CI/CD 构建流程
├── 代码提交
│   ├── Git Hook
│   ├── Webhook
│   └── 定时触发
├── 构建阶段
│   ├── 代码检出
│   ├── 依赖安装
│   ├── 代码编译
│   ├── 测试执行
│   └── 构建产物
├── 容器化
│   ├── 镜像构建
│   ├── 安全扫描
│   ├── 镜像推送
│   └── 镜像标签
└── 部署
    ├── 环境部署
    ├── 健康检查
    ├── 回滚机制
    └── 监控告警
```

### 2. 构建工具集成

| 工具 | 构建支持 | 容器化 | 缓存 | 并行 | 安全扫描 |
|------|----------|--------|------|------|----------|
| Jenkins | ✅ | ✅ | ✅ | ✅ | ✅ |
| GitLab CI | ✅ | ✅ | ✅ | ✅ | ✅ |
| GitHub Actions | ✅ | ✅ | ✅ | ✅ | ✅ |
| Azure DevOps | ✅ | ✅ | ✅ | ✅ | ✅ |
| CircleCI | ✅ | ✅ | ✅ | ✅ | ✅ |

## 🚀 GitHub Actions 集成

### 1. 基础构建流水线

```yaml
# .github/workflows/build.yml
name: Build and Deploy

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - name: Checkout repository
      uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=sha

    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
        format: 'sarif'
        output: 'trivy-results.sarif'

    - name: Upload Trivy scan results
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'trivy-results.sarif'
```

### 2. 多环境构建

```yaml
# .github/workflows/multi-env.yml
name: Multi-Environment Build

on:
  push:
    branches: [ main, develop, staging ]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        environment: [dev, staging, prod]
        include:
          - environment: dev
            branch: develop
            registry: dev-registry.example.com
          - environment: staging
            branch: staging
            registry: staging-registry.example.com
          - environment: prod
            branch: main
            registry: prod-registry.example.com

    steps:
    - name: Checkout
      uses: actions/checkout@v3

    - name: Build for ${{ matrix.environment }}
      run: |
        docker build -t ${{ matrix.registry }}/myapp:${{ github.sha }} .
        docker push ${{ matrix.registry }}/myapp:${{ github.sha }}
```

## 🔧 GitLab CI 集成

### 1. 基础构建流水线

```yaml
# .gitlab-ci.yml
stages:
  - build
  - test
  - security
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"

build:
  stage: build
  image: docker:20.10.16
  services:
    - docker:20.10.16-dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main
    - develop

test:
  stage: test
  image: node:18-alpine
  script:
    - npm ci
    - npm run test
    - npm run lint
  coverage: '/Lines\s*:\s*(\d+\.\d+)%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage/cobertura-coverage.xml

security:
  stage: security
  image: aquasec/trivy:latest
  script:
    - trivy image --exit-code 0 --no-progress $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  allow_failure: true
```

### 2. 高级构建流水线

```yaml
# .gitlab-ci.yml
stages:
  - build
  - test
  - security
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"
  DOCKER_BUILDKIT: 1

build:
  stage: build
  image: docker:20.10.16
  services:
    - docker:20.10.16-dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker buildx create --use
    - docker buildx build --platform linux/amd64,linux/arm64 -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA --push .
  only:
    - main
    - develop

test:
  stage: test
  image: node:18-alpine
  script:
    - npm ci
    - npm run test
    - npm run lint
  coverage: '/Lines\s*:\s*(\d+\.\d+)%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage/cobertura-coverage.xml

security:
  stage: security
  image: aquasec/trivy:latest
  script:
    - trivy image --exit-code 0 --no-progress $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  allow_failure: true

deploy:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main
  when: manual
```

## 🏗️ Jenkins 集成

### 1. Pipeline 构建

```groovy
// Jenkinsfile
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
                sh 'npm ci'
                sh 'npm run build'
            }
        }
        
        stage('Test') {
            steps {
                sh 'npm run test'
                sh 'npm run lint'
            }
            post {
                always {
                    junit 'test-results.xml'
                }
            }
        }
        
        stage('Security Scan') {
            steps {
                sh 'trivy image --exit-code 0 --no-progress ${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}'
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
                    }
                }
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

### 2. 多分支构建

```groovy
// Jenkinsfile
pipeline {
    agent any
    
    parameters {
        choice(
            name: 'ENVIRONMENT',
            choices: ['dev', 'staging', 'prod'],
            description: 'Target environment'
        )
        string(
            name: 'VERSION',
            defaultValue: 'latest',
            description: 'Image version'
        )
    }
    
    environment {
        DOCKER_IMAGE = 'myapp'
        DOCKER_TAG = "${params.VERSION}"
        REGISTRY = "registry-${params.ENVIRONMENT}.example.com"
    }
    
    stages {
        stage('Build') {
            steps {
                sh 'npm ci'
                sh 'npm run build'
            }
        }
        
        stage('Test') {
            steps {
                sh 'npm run test'
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
                    }
                }
            }
        }
        
        stage('Deploy') {
            steps {
                script {
                    sh "kubectl set image deployment/myapp myapp=${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG} -n ${params.ENVIRONMENT}"
                }
            }
        }
    }
}
```

## 🛠️ 实践练习

### 练习1: 多语言项目构建

```yaml
# .github/workflows/multi-lang.yml
name: Multi-Language Build

on:
  push:
    branches: [ main ]

jobs:
  build-frontend:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
    - name: Install dependencies
      run: cd frontend && npm ci
    - name: Build frontend
      run: cd frontend && npm run build

  build-backend:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Setup Java
      uses: actions/setup-java@v3
      with:
        java-version: '11'
        distribution: 'temurin'
    - name: Build backend
      run: cd backend && ./mvnw clean package

  build-docker:
    needs: [build-frontend, build-backend]
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Build Docker image
      run: docker build -t myapp:latest .
```

### 练习2: 构建缓存优化

```yaml
# .github/workflows/cache-optimized.yml
name: Cache Optimized Build

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Cache dependencies
      uses: actions/cache@v3
      with:
        path: |
          ~/.npm
          ~/.m2
        key: ${{ runner.os }}-deps-${{ hashFiles('**/package-lock.json', '**/pom.xml') }}
        restore-keys: |
          ${{ runner.os }}-deps-
    
    - name: Install dependencies
      run: |
        npm ci
        cd backend && ./mvnw dependency:resolve
    
    - name: Build
      run: |
        npm run build
        cd backend && ./mvnw clean package
```

## 📚 相关资源

### 官方文档
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [GitLab CI 文档](https://docs.gitlab.com/ee/ci/)
- [Jenkins 文档](https://www.jenkins.io/doc/)

### 学习资源
- [CI/CD 最佳实践](https://martinfowler.com/articles/ci.html)
- [构建性能优化](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows)

### 工具推荐
- **GitHub Actions**: GitHub 集成 CI/CD
- **GitLab CI**: GitLab 集成 CI/CD
- **Jenkins**: 企业级 CI/CD 平台
- **Azure DevOps**: 微软云 CI/CD
- **CircleCI**: 云原生 CI/CD

---

**掌握 CI/CD 集成构建，实现自动化部署！** 🚀
