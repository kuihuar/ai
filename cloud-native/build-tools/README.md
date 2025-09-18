# 构建工具与容器化

## 📚 学习目标

通过本模块学习，您将掌握：
- 现代构建工具的使用和最佳实践
- 容器化构建流程和优化策略
- 多阶段构建和构建缓存优化
- 构建安全扫描和漏洞管理
- 构建流水线集成和自动化

## 🎯 核心概念

### 1. 构建工具生态

现代应用构建涉及多个层面的工具：

```
┌─────────────────────────────────────────────────┐
│                应用构建工具                      │
├─────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────┐   │
│  │ 语言构建工具 │ │ 包管理工具   │ │ 测试工具 │   │
│  │ Maven/Gradle│ │ npm/yarn    │ │ JUnit   │   │
│  │ npm/yarn    │ │ pip/poetry  │ │ Jest    │   │
│  │ pip/poetry  │ │ go mod      │ │ pytest  │   │
│  │ go build    │ │ cargo       │ │ go test │   │
│  └─────────────┘ └─────────────┘ └─────────┘   │
├─────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────┐   │
│  │ 容器构建工具 │ │ 镜像优化工具 │ │ 安全工具 │   │
│  │ Docker      │ │ BuildKit    │ │ Trivy   │   │
│  │ Buildah     │ │ Dive        │ │ Snyk    │   │
│  │ Podman      │ │ Distroless  │ │ Clair   │   │
│  └─────────────┘ └─────────────┘ └─────────┘   │
├─────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────┐   │
│  │ CI/CD工具    │ │ 包管理工具   │ │ 监控工具 │   │
│  │ Jenkins     │ │ Helm        │ │ Prometheus│   │
│  │ GitLab CI   │ │ Kustomize   │ │ Grafana  │   │
│  │ GitHub Actions│ │ Operator   │ │ ELK Stack│   │
│  └─────────────┘ └─────────────┘ └─────────┘   │
└─────────────────────────────────────────────────┘
```

### 2. 构建流程演进

```
传统构建                   现代构建
┌─────────────┐            ┌─────────────────┐
│ 本地编译    │            │ 容器化构建      │
├─────────────┤            ├─────────────────┤
│ 手动打包    │            │ 多阶段构建      │
├─────────────┤            ├─────────────────┤
│ 手动部署    │            │ 自动化流水线    │
├─────────────┤            ├─────────────────┤
│ 环境差异    │            │ 环境一致性      │
├─────────────┤            ├─────────────────┤
│ 配置复杂    │            │ 声明式配置      │
└─────────────┘            └─────────────────┘
```

## 🛠️ 语言构建工具

### 1. Java 构建工具

#### Maven
```xml
<!-- pom.xml -->
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>myapp</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>

    <properties>
        <maven.compiler.source>11</maven.compiler.source>
        <maven.compiler.target>11</maven.compiler.target>
        <spring.boot.version>2.7.0</spring.boot.version>
    </properties>

    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>${spring.boot.version}</version>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
                <version>${spring.boot.version}</version>
                <executions>
                    <execution>
                        <goals>
                            <goal>repackage</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
</project>
```

#### Gradle
```gradle
// build.gradle
plugins {
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
    id 'java'
}

group = 'com.example'
version = '1.0.0'
sourceCompatibility = '11'

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}

tasks.named('test') {
    useJUnitPlatform()
}
```

### 2. Node.js 构建工具

#### npm/yarn
```json
// package.json
{
  "name": "myapp",
  "version": "1.0.0",
  "scripts": {
    "build": "next build",
    "start": "next start",
    "dev": "next dev",
    "test": "jest",
    "lint": "eslint .",
    "type-check": "tsc --noEmit"
  },
  "dependencies": {
    "next": "12.3.0",
    "react": "18.2.0",
    "react-dom": "18.2.0"
  },
  "devDependencies": {
    "@types/react": "18.0.0",
    "@types/node": "18.0.0",
    "typescript": "4.8.0",
    "jest": "28.0.0",
    "eslint": "8.0.0"
  }
}
```

#### 构建脚本
```bash
#!/bin/bash
# build.sh

set -e

echo "Installing dependencies..."
npm ci

echo "Running type check..."
npm run type-check

echo "Running linting..."
npm run lint

echo "Running tests..."
npm run test

echo "Building application..."
npm run build

echo "Build completed successfully!"
```

### 3. Python 构建工具

#### Poetry
```toml
# pyproject.toml
[tool.poetry]
name = "myapp"
version = "1.0.0"
description = "My Python Application"
authors = ["Your Name <your.email@example.com>"]

[tool.poetry.dependencies]
python = "^3.9"
fastapi = "^0.68.0"
uvicorn = "^0.15.0"

[tool.poetry.group.dev.dependencies]
pytest = "^6.2.0"
black = "^21.0.0"
flake8 = "^3.9.0"

[build-system]
requires = ["poetry-core>=1.0.0"]
build-backend = "poetry.core.masonry.api"

[tool.poetry.scripts]
myapp = "myapp.main:main"
```

#### 构建脚本
```bash
#!/bin/bash
# build.sh

set -e

echo "Installing dependencies..."
poetry install

echo "Running linting..."
poetry run black --check .
poetry run flake8 .

echo "Running tests..."
poetry run pytest

echo "Building package..."
poetry build

echo "Build completed successfully!"
```

### 4. Go 构建工具

#### Go Modules
```go
// go.mod
module github.com/example/myapp

go 1.19

require (
    github.com/gin-gonic/gin v1.9.0
    github.com/stretchr/testify v1.8.2
)

require (
    github.com/davecgh/go-spew v1.1.1 // indirect
    github.com/pmezard/go-difflib v1.0.0 // indirect
    gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

#### 构建脚本
```bash
#!/bin/bash
# build.sh

set -e

echo "Installing dependencies..."
go mod download
go mod verify

echo "Running tests..."
go test -v ./...

echo "Running linting..."
golangci-lint run

echo "Building application..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

echo "Build completed successfully!"
```

## 🐳 容器化构建

### 1. Docker 多阶段构建

#### 基础多阶段构建
```dockerfile
# 构建阶段
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# 生产阶段
FROM node:18-alpine AS production
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY . .
USER node
EXPOSE 3000
CMD ["npm", "start"]
```

#### 高级多阶段构建
```dockerfile
# 依赖阶段
FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# 构建阶段
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# 测试阶段
FROM node:18-alpine AS tester
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run test

# 生产阶段
FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package*.json ./

USER nextjs
EXPOSE 3000
CMD ["npm", "start"]
```

### 2. BuildKit 高级特性

#### 启用 BuildKit
```bash
# 设置环境变量
export DOCKER_BUILDKIT=1

# 或使用 buildx
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .
```

#### 构建缓存优化
```dockerfile
# syntax=docker/dockerfile:1
FROM node:18-alpine

# 使用构建缓存
RUN --mount=type=cache,target=/root/.npm \
    npm install

# 并行构建
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm ci
```

#### 多平台构建
```dockerfile
# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:18-alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

FROM node:18-alpine
COPY --from=builder /app /app
```

### 3. 构建优化策略

#### 镜像大小优化
```dockerfile
# 使用 distroless 基础镜像
FROM gcr.io/distroless/nodejs18-debian11

# 使用 Alpine 基础镜像
FROM node:18-alpine

# 多阶段构建减少层数
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci && npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
```

#### 构建速度优化
```dockerfile
# 利用 Docker 层缓存
FROM node:18-alpine
WORKDIR /app

# 先复制依赖文件
COPY package*.json ./
RUN npm ci

# 再复制源代码
COPY . .
RUN npm run build
```

## 🔒 构建安全

### 1. 安全扫描工具

#### Trivy 扫描
```bash
# 安装 Trivy
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh

# 扫描镜像
trivy image nginx:latest

# 扫描文件系统
trivy fs .

# 扫描配置文件
trivy config .
```

#### Snyk 扫描
```bash
# 安装 Snyk
npm install -g snyk

# 扫描项目
snyk test

# 扫描容器镜像
snyk container test nginx:latest

# 监控项目
snyk monitor
```

### 2. 安全构建实践

#### 非 root 用户
```dockerfile
FROM node:18-alpine

# 创建非 root 用户
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# 切换到非 root 用户
USER nextjs

# 设置工作目录权限
WORKDIR /app
RUN chown nextjs:nodejs /app
```

#### 最小权限原则
```dockerfile
FROM node:18-alpine

# 只安装必要的包
RUN apk add --no-cache \
    dumb-init \
    && rm -rf /var/cache/apk/*

# 使用非特权端口
EXPOSE 3000

# 只读文件系统
RUN mkdir -p /app/tmp
VOLUME ["/app/tmp"]
```

## 🚀 CI/CD 集成

### 1. GitHub Actions

#### 构建流水线
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

### 2. GitLab CI

#### 构建流水线
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

### 3. Jenkins Pipeline

#### 构建流水线
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

## 🛠️ 实践练习

### 练习1: 多语言应用构建

```dockerfile
# 前端构建
FROM node:18-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# 后端构建
FROM golang:1.19-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 最终镜像
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend /app/backend/main .
COPY --from=frontend /app/frontend/dist ./static
CMD ["./main"]
```

### 练习2: 构建优化

```dockerfile
# 使用 BuildKit 缓存
# syntax=docker/dockerfile:1
FROM node:18-alpine

# 缓存依赖安装
RUN --mount=type=cache,target=/root/.npm \
    npm install

# 并行构建
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm ci
```

## 📚 相关资源

### 官方文档
- [Docker 官方文档](https://docs.docker.com/)
- [BuildKit 文档](https://docs.docker.com/build/buildkit/)
- [Trivy 文档](https://aquasecurity.github.io/trivy/)

### 学习资源
- [Docker 最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [容器安全指南](https://kubernetes.io/docs/concepts/security/)

### 工具推荐
- **Docker**: 容器化平台
- **BuildKit**: 高级构建特性
- **Trivy**: 安全扫描工具
- **Dive**: 镜像分析工具
- **Snyk**: 漏洞管理平台

---

**掌握现代构建工具，实现高效容器化部署！** 🚀
