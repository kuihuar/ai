# 容器化构建详解

## 📚 学习目标

通过本模块学习，您将掌握：
- 容器化构建的核心概念和最佳实践
- 多阶段构建和构建缓存优化
- 容器镜像安全扫描和漏洞管理
- 构建工具选择和性能优化
- 企业级容器构建流水线

## 🎯 容器构建概览

### 1. 容器构建工具生态

```
容器构建工具生态
├── 容器运行时
│   ├── Docker
│   ├── Podman
│   ├── Containerd
│   └── CRI-O
├── 构建工具
│   ├── Docker Build
│   ├── BuildKit
│   ├── Buildah
│   ├── Kaniko
│   └── Jib
├── 镜像优化
│   ├── Dive
│   ├── Distroless
│   ├── Alpine
│   └── Scratch
└── 安全扫描
    ├── Trivy
    ├── Snyk
    ├── Clair
    └── Anchore
```

### 2. 构建流程对比

| 特性 | Docker | Buildah | Kaniko | Jib |
|------|--------|---------|--------|-----|
| 无守护进程 | ❌ | ✅ | ✅ | ✅ |
| 多阶段构建 | ✅ | ✅ | ✅ | ✅ |
| 缓存支持 | ✅ | ✅ | ✅ | ✅ |
| 安全扫描 | ✅ | ✅ | ✅ | ✅ |
| 学习曲线 | 简单 | 中等 | 简单 | 简单 |

## 🐳 Docker 构建详解

### 1. 基础 Dockerfile

```dockerfile
# 多阶段构建示例
FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

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

### 2. 高级构建技巧

#### 构建缓存优化
```dockerfile
# 利用层缓存
FROM node:18-alpine
WORKDIR /app

# 先复制依赖文件，利用缓存
COPY package*.json ./
RUN npm ci

# 再复制源代码
COPY . .
RUN npm run build
```

#### 多平台构建
```dockerfile
# 多平台构建
FROM --platform=$BUILDPLATFORM node:18-alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "Building on $BUILDPLATFORM for $TARGETPLATFORM"

FROM node:18-alpine
COPY --from=builder /app /app
```

## 🚀 BuildKit 高级特性

### 1. BuildKit 配置

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

### 2. 构建优化

```bash
# 启用 BuildKit
export DOCKER_BUILDKIT=1

# 多平台构建
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .

# 构建缓存
docker buildx build --cache-from=type=local,src=/tmp/.buildx-cache .
```

## 🔒 容器安全

### 1. 安全扫描

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
```

### 2. 安全最佳实践

```dockerfile
# 使用非 root 用户
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

**掌握容器化构建，实现高效部署！** 🚀
