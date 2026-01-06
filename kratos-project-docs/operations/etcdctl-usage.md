# etcdctl 使用指南

本文档介绍如何使用 `etcdctl` 命令行工具查看和管理 etcd 中的服务注册信息。

## 安装 etcdctl

### macOS

使用 Homebrew 安装：

```bash
brew install etcd
```

### Linux

```bash
# Ubuntu/Debian
sudo apt-get install etcd-client

# CentOS/RHEL
sudo yum install etcd
```

### 从源码安装

```bash
# 下载 etcd 二进制文件
ETCD_VER=v3.5.0
wget https://github.com/etcd-io/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf etcd-${ETCD_VER}-linux-amd64.tar.gz
sudo mv etcd-${ETCD_VER}-linux-amd64/etcdctl /usr/local/bin/
```

## 基本配置

### 设置 API 版本

etcd v3 需要设置环境变量：

```bash
export ETCDCTL_API=3
```

或者在使用命令时指定：

```bash
ETCDCTL_API=3 etcdctl <command>
```

### 指定 etcd 端点

如果 etcd 不在默认位置，需要指定端点：

```bash
etcdctl --endpoints=http://127.0.0.1:2379 <command>
```

## 查看服务注册信息

### 1. 查看所有注册的服务

Kratos 框架将服务注册到 `/kratos` 前缀下：

```bash
# 设置 API 版本
export ETCDCTL_API=3

# 查看所有注册的服务（只显示键）
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos --prefix --keys-only
```

### 2. 查看特定服务的详细信息

```bash
# 查看所有服务实例的完整信息（键和值）
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos --prefix
```

### 3. 查看特定服务名称的实例

如果服务名称为 `sre`：

```bash
# 查看 sre 服务的所有实例
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos/sre --prefix
```

### 4. 查看特定协议类型的服务

Kratos 会为每个端点（gRPC 和 HTTP）创建独立的注册项：

```bash
# 查看所有 gRPC 服务
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos --prefix | grep grpc

# 查看所有 HTTP 服务
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos --prefix | grep http
```

### 5. 格式化输出 JSON

etcd 中存储的值是 JSON 格式，可以使用 `jq` 工具美化输出：

```bash
# 安装 jq（如果未安装）
# macOS: brew install jq
# Linux: sudo apt-get install jq

# 查看并格式化 JSON
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos --prefix | grep -v "^/kratos" | jq .
```

## 键的格式

Kratos etcd 注册的键格式为：

```
/kratos/{Name}/{scheme}/{ID}
```

- `Name`: 服务名称（通过 `kratos.Name()` 设置）
- `scheme`: 协议类型（`grpc` 或 `http`）
- `ID`: 服务实例 ID（通常是主机名）

示例：
- `/kratos/sre/grpc/JianfendeMacBook-Pro.local`
- `/kratos/sre/http/JianfendeMacBook-Pro.local`

如果服务名称为空，键格式可能是：
- `/kratos//grpc/{ID}`
- `/kratos//http/{ID}`

## 值的格式

每个键对应的值是一个 JSON 对象，包含以下字段：

```json
{
  "id": "JianfendeMacBook-Pro.local",
  "name": "sre",
  "version": "v1.0.0",
  "endpoints": [
    "grpc://192.168.1.100:9000",
    "http://192.168.1.100:8000"
  ],
  "metadata": {
    "env": "development",
    "region": "local"
  }
}
```

## 常用命令

### 列出所有键

```bash
etcdctl --endpoints=http://127.0.0.1:2379 get "" --prefix --keys-only
```

### 查看特定键的值

```bash
etcdctl --endpoints=http://127.0.0.1:2379 get /kratos/sre/grpc/JianfendeMacBook-Pro.local
```

### 删除服务注册（手动取消注册）

```bash
# 删除特定服务实例
etcdctl --endpoints=http://127.0.0.1:2379 del /kratos/sre/grpc/JianfendeMacBook-Pro.local

# 删除所有服务实例
etcdctl --endpoints=http://127.0.0.1:2379 del /kratos --prefix
```

### 监控服务变化（Watch）

```bash
# 监控 /kratos 前缀下的所有变化
etcdctl --endpoints=http://127.0.0.1:2379 watch /kratos --prefix
```

## 使用项目提供的检查工具

项目提供了一个 Go 工具来查看注册信息：

```bash
# 使用默认 etcd 地址 (127.0.0.1:2379)
go run cmd/check-registry/main.go

# 指定 etcd 地址
go run cmd/check-registry/main.go 192.168.1.100:2379
```

## 故障排查

### 1. 连接失败

如果出现连接错误，检查：
- etcd 是否在运行：`curl http://127.0.0.1:2379/health`
- 端口是否正确
- 防火墙设置

### 2. 找不到服务

如果找不到注册的服务：
- 确认服务已启动并成功注册
- 检查服务名称是否正确
- 如果服务名称为空，尝试查找 `/kratos//grpc/` 或 `/kratos//http/`
- 查看服务日志确认注册是否成功

### 3. 权限问题

如果使用 TLS，需要提供证书：

```bash
etcdctl --endpoints=https://127.0.0.1:2379 \
        --cacert=/path/to/ca.crt \
        --cert=/path/to/client.crt \
        --key=/path/to/client.key \
        get /kratos --prefix
```

## 参考

- [etcd 官方文档](https://etcd.io/docs/)
- [etcdctl 命令参考](https://etcd.io/docs/latest/dev-guide/interacting_v3/)
- [Kratos 服务注册文档](../architecture/service-registry-discovery.md)

