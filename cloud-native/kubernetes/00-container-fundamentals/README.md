# 容器基础概念

## 📚 学习目标

通过本模块学习，您将掌握：
- 容器技术的核心概念和原理
- Docker 基础操作和最佳实践
- 容器镜像构建和优化
- 容器网络和存储
- 容器安全基础

## 🎯 核心概念

### 1. 什么是容器？

容器是一种轻量级的虚拟化技术，它将应用程序及其依赖项打包在一起，提供一致的运行环境。

**关键特性：**
- **隔离性**: 进程、网络、文件系统隔离
- **可移植性**: 一次构建，到处运行
- **轻量级**: 共享宿主机内核，资源占用少
- **快速启动**: 秒级启动时间

### 2. 容器 vs 虚拟机

| 特性 | 容器 | 虚拟机 |
|------|------|--------|
| 资源占用 | 轻量级 | 重量级 |
| 启动时间 | 秒级 | 分钟级 |
| 隔离性 | 进程级 | 硬件级 |
| 性能 | 接近原生 | 有损耗 |
| 可移植性 | 高 | 中等 |

### 3. 容器核心技术

#### Namespace（命名空间）
- **PID Namespace**: 进程隔离
- **Network Namespace**: 网络隔离
- **Mount Namespace**: 文件系统隔离
- **UTS Namespace**: 主机名隔离
- **IPC Namespace**: 进程间通信隔离
- **User Namespace**: 用户隔离

#### Cgroups（控制组）
- **CPU 限制**: 限制 CPU 使用率
- **内存限制**: 限制内存使用量
- **I/O 限制**: 限制磁盘 I/O
- **网络限制**: 限制网络带宽

## 🐳 Docker 基础

### 1. Docker 架构

```
┌─────────────────────────────────────┐
│           Docker Client             │
├─────────────────────────────────────┤
│         Docker Daemon               │
│  ┌─────────────┐ ┌─────────────┐   │
│  │  Container  │ │  Container  │   │
│  └─────────────┘ └─────────────┘   │
│  ┌─────────────────────────────┐   │
│  │      Docker Images          │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
```

### 2. 核心组件

- **Docker Client**: 命令行工具
- **Docker Daemon**: 后台服务
- **Docker Registry**: 镜像仓库
- **Docker Images**: 只读模板
- **Docker Containers**: 运行实例

### 3. 基本命令

```bash
# 镜像操作
docker pull nginx:latest
docker images
docker rmi nginx:latest

# 容器操作
docker run -d --name web nginx:latest
docker ps
docker stop web
docker rm web

# 进入容器
docker exec -it web /bin/bash
```

## 📦 容器镜像

### 1. 镜像分层结构

```
┌─────────────────────────────────────┐
│         Application Layer           │
├─────────────────────────────────────┤
│         Runtime Layer               │
├─────────────────────────────────────┤
│         OS Libraries                │
├─────────────────────────────────────┤
│         Base OS Image               │
└─────────────────────────────────────┘
```

### 2. Dockerfile 最佳实践

```dockerfile
# 使用官方基础镜像
FROM node:18-alpine

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制应用代码
COPY . .

# 创建非root用户
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs

# 暴露端口
EXPOSE 3000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

# 启动命令
CMD ["npm", "start"]
```

### 3. 镜像优化技巧

```dockerfile
# 多阶段构建
FROM node:18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine AS production
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package*.json ./
USER node
EXPOSE 3000
CMD ["npm", "start"]
```

## 🌐 容器网络

### 1. Docker 网络模式

```bash
# 查看网络
docker network ls

# 创建自定义网络
docker network create my-network

# 运行容器并指定网络
docker run -d --name web --network my-network nginx
```

### 2. 网络类型

- **Bridge**: 默认网络模式
- **Host**: 使用宿主机网络
- **None**: 无网络连接
- **Overlay**: 跨主机网络

## 💾 容器存储

### 1. 存储类型

```bash
# 数据卷
docker volume create my-volume
docker run -v my-volume:/data nginx

# 绑定挂载
docker run -v /host/path:/container/path nginx

# 临时文件系统
docker run --tmpfs /tmp nginx
```

### 2. 存储驱动

- **overlay2**: 推荐，性能好
- **aufs**: 兼容性好
- **devicemapper**: 企业级
- **btrfs**: 高级特性

## 🔒 容器安全

### 1. 基础安全实践

```dockerfile
# 使用非root用户
USER 1001

# 最小化镜像
FROM alpine:latest

# 扫描漏洞
docker scan nginx:latest

# 只读文件系统
docker run --read-only nginx
```

### 2. 安全扫描

```bash
# 使用 Trivy 扫描
trivy image nginx:latest

# 使用 Docker Scout
docker scout quickview nginx:latest
```

## 🛠️ 实践练习

### 练习1: 构建 Web 应用镜像

```dockerfile
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 练习2: 多容器应用

```yaml
# docker-compose.yml
version: '3.8'
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html
  
  db:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 📚 相关资源

### 官方文档
- [Docker 官方文档](https://docs.docker.com/)
- [容器运行时规范](https://github.com/opencontainers/runtime-spec)

### 学习资源
- [Docker 最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [容器安全指南](https://kubernetes.io/docs/concepts/security/)

### 工具推荐
- **Docker Desktop**: 本地开发环境
- **Docker Compose**: 多容器编排
- **Trivy**: 安全扫描工具
- **Dive**: 镜像分析工具

---

**掌握容器基础，为 Kubernetes 学习打下坚实基础！** 🚀
