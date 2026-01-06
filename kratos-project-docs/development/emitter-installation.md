# Emitter 安装和部署

## 概述

本文档介绍如何安装和部署 Emitter 消息网关，包括单机部署、集群部署和 Docker 部署等方案。

## 系统要求

### 最低要求

- **CPU**：2 核
- **内存**：2GB
- **磁盘**：10GB（用于日志和消息存储）
- **操作系统**：Linux、macOS、Windows

### 推荐配置

- **CPU**：4+ 核
- **内存**：8GB+
- **磁盘**：SSD，100GB+
- **网络**：千兆网络

## 安装方式

### 1. Docker 部署（推荐）

#### 单机部署

```bash
# 拉取镜像
docker pull emitter/server

# 运行容器
docker run -d \
  --name emitter \
  -p 8080:8080 \
  -p 443:443 \
  -p 8443:8443 \
  emitter/server
```

#### 使用 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  emitter:
    image: emitter/server:latest
    container_name: emitter
    ports:
      - "8080:8080"   # HTTP API
      - "443:443"      # MQTT over TLS
      - "8443:8443"    # MQTT over WebSocket
    environment:
      - EMITTER_LICENSE=your-license-key
      - EMITTER_CLUSTER_NAME=emitter-cluster
      - EMITTER_CLUSTER_NODES=emitter:8080
    volumes:
      - emitter-data:/var/lib/emitter
    restart: unless-stopped

volumes:
  emitter-data:
```

启动服务：

```bash
docker-compose up -d
```

### 2. 二进制部署

#### 下载二进制文件

```bash
# Linux
wget https://github.com/emitter-io/emitter/releases/latest/download/emitter-linux-amd64
chmod +x emitter-linux-amd64
sudo mv emitter-linux-amd64 /usr/local/bin/emitter

# macOS
wget https://github.com/emitter-io/emitter/releases/latest/download/emitter-darwin-amd64
chmod +x emitter-darwin-amd64
sudo mv emitter-darwin-amd64 /usr/local/bin/emitter
```

#### 运行服务

```bash
# 基本运行
emitter

# 指定配置文件
emitter -config /etc/emitter/config.yaml

# 后台运行
nohup emitter > /var/log/emitter.log 2>&1 &
```

### 3. 源码编译

```bash
# 克隆仓库
git clone https://github.com/emitter-io/emitter.git
cd emitter

# 编译
go build -o emitter ./cmd/emitter

# 运行
./emitter
```

## 配置文件

### 基本配置

```yaml
# config.yaml
# Emitter 服务器配置

# 服务器配置
server:
  # HTTP API 端口
  http_port: 8080
  
  # MQTT 端口
  mqtt_port: 1883
  
  # MQTT over TLS 端口
  mqtt_tls_port: 8883
  
  # WebSocket 端口
  ws_port: 8080
  
  # WebSocket over TLS 端口
  wss_port: 8443

# 集群配置
cluster:
  # 集群名称
  name: emitter-cluster
  
  # 节点列表
  nodes:
    - emitter1:8080
    - emitter2:8080
    - emitter3:8080
  
  # 监听地址
  listen: 0.0.0.0:8080

# 存储配置
storage:
  # 存储类型：memory, badger, bolt
  type: badger
  
  # 数据目录
  path: /var/lib/emitter/data

# 日志配置
logging:
  # 日志级别：debug, info, warn, error
  level: info
  
  # 日志文件
  file: /var/log/emitter.log

# 安全配置
security:
  # TLS 证书文件
  cert_file: /etc/emitter/cert.pem
  key_file: /etc/emitter/key.pem
```

### 环境变量配置

```bash
# 许可证密钥
export EMITTER_LICENSE=your-license-key

# 集群配置
export EMITTER_CLUSTER_NAME=emitter-cluster
export EMITTER_CLUSTER_NODES=emitter1:8080,emitter2:8080

# 存储路径
export EMITTER_STORAGE_PATH=/var/lib/emitter/data

# 日志级别
export EMITTER_LOG_LEVEL=info
```

## 集群部署

### 1. 准备节点

准备 3 个或更多节点（建议奇数个，用于选举）：

```
Node 1: emitter1.example.com
Node 2: emitter2.example.com
Node 3: emitter3.example.com
```

### 2. 配置每个节点

#### Node 1 配置

```yaml
# emitter1/config.yaml
cluster:
  name: emitter-cluster
  nodes:
    - emitter1:8080
    - emitter2:8080
    - emitter3:8080
  listen: 0.0.0.0:8080
```

#### Node 2 配置

```yaml
# emitter2/config.yaml
cluster:
  name: emitter-cluster
  nodes:
    - emitter1:8080
    - emitter2:8080
    - emitter3:8080
  listen: 0.0.0.0:8080
```

#### Node 3 配置

```yaml
# emitter3/config.yaml
cluster:
  name: emitter-cluster
  nodes:
    - emitter1:8080
    - emitter2:8080
    - emitter3:8080
  listen: 0.0.0.0:8080
```

### 3. 启动集群

在每个节点上启动服务：

```bash
# Node 1
emitter -config /etc/emitter/emitter1/config.yaml

# Node 2
emitter -config /etc/emitter/emitter2/config.yaml

# Node 3
emitter -config /etc/emitter/emitter3/config.yaml
```

### 4. 验证集群

```bash
# 检查集群状态
curl http://emitter1:8080/api/cluster/status

# 应该返回所有节点的状态
```

## TLS/SSL 配置

### 1. 生成证书

```bash
# 生成私钥
openssl genrsa -out key.pem 2048

# 生成证书签名请求
openssl req -new -key key.pem -out csr.pem

# 生成自签名证书（开发环境）
openssl x509 -req -days 365 -in csr.pem -signkey key.pem -out cert.pem
```

### 2. 配置 TLS

```yaml
# config.yaml
security:
  cert_file: /etc/emitter/cert.pem
  key_file: /etc/emitter/key.pem
```

### 3. 使用 Let's Encrypt（生产环境）

```bash
# 安装 certbot
sudo apt-get install certbot

# 获取证书
sudo certbot certonly --standalone -d emitter.example.com

# 配置证书路径
security:
  cert_file: /etc/letsencrypt/live/emitter.example.com/fullchain.pem
  key_file: /etc/letsencrypt/live/emitter.example.com/privkey.pem
```

## 监控和日志

### 1. 健康检查

```bash
# HTTP API 健康检查
curl http://localhost:8080/health

# 应该返回：{"status":"ok"}
```

### 2. 指标监控

Emitter 提供 Prometheus 格式的指标：

```bash
# 获取指标
curl http://localhost:8080/metrics
```

### 3. 日志管理

```yaml
# config.yaml
logging:
  level: info
  file: /var/log/emitter.log
  max_size: 100  # MB
  max_backups: 10
  max_age: 30    # days
```

使用 logrotate 管理日志：

```bash
# /etc/logrotate.d/emitter
/var/log/emitter.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0644 emitter emitter
}
```

## 性能调优

### 1. 系统参数调优

```bash
# 增加文件描述符限制
ulimit -n 65535

# 增加网络缓冲区
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
```

### 2. 存储优化

```yaml
# config.yaml
storage:
  type: badger
  path: /var/lib/emitter/data
  # Badger 选项
  badger:
    sync_writes: false  # 异步写入，提高性能
    value_log_file_size: 1073741824  # 1GB
```

### 3. 连接池配置

```yaml
# config.yaml
server:
  max_connections: 100000
  connection_timeout: 60s
  keepalive: 300s
```

## 故障排查

### 1. 检查服务状态

```bash
# 检查进程
ps aux | grep emitter

# 检查端口
netstat -tlnp | grep 8080
```

### 2. 查看日志

```bash
# 查看日志
tail -f /var/log/emitter.log

# 查看错误日志
grep ERROR /var/log/emitter.log
```

### 3. 常见问题

#### 问题：无法启动

**原因**：端口被占用

**解决**：
```bash
# 检查端口占用
lsof -i :8080

# 修改配置文件中的端口
```

#### 问题：集群无法连接

**原因**：网络不通或防火墙阻止

**解决**：
```bash
# 检查网络连通性
ping emitter2

# 检查防火墙
sudo ufw status
sudo ufw allow 8080/tcp
```

#### 问题：性能问题

**原因**：资源不足或配置不当

**解决**：
- 增加内存和 CPU
- 优化存储配置
- 调整连接池大小

## 备份和恢复

### 1. 数据备份

```bash
# 备份数据目录
tar -czf emitter-backup-$(date +%Y%m%d).tar.gz /var/lib/emitter/data

# 定期备份（使用 cron）
0 2 * * * tar -czf /backup/emitter-$(date +\%Y\%m\%d).tar.gz /var/lib/emitter/data
```

### 2. 数据恢复

```bash
# 停止服务
systemctl stop emitter

# 恢复数据
tar -xzf emitter-backup-20240101.tar.gz -C /

# 启动服务
systemctl start emitter
```

## 参考资源

- [Emitter 官方文档](https://emitter.io/docs/)
- [Emitter GitHub](https://github.com/emitter-io/emitter)
- [Docker Hub](https://hub.docker.com/r/emitter/server)

## 总结

Emitter 的安装和部署相对简单，支持多种部署方式。推荐使用 Docker 部署，便于管理和扩展。对于生产环境，建议使用集群部署以提高可用性和性能。

