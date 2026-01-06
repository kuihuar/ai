## 容器访问宿主机 MySQL / Redis 说明

### 场景概述

当前项目已经使用 Docker 将 `sre` 服务打包成容器运行，但 MySQL 和 Redis 仍运行在宿主机上（本机）。  
容器内的服务需要访问宿主机的 MySQL / Redis，这里说明推荐做法和配置示例。

---

### 关键原则

- **容器里的 `127.0.0.1` 不是宿主机**，而是容器自己。
- 在 macOS / Windows 的 Docker Desktop 中，可以使用 **`host.docker.internal`** 这个固定主机名访问宿主机。
- 在 Linux 服务器上，可以使用 **宿主机实际 IP** 或 `--network host` 等方式。

---

### 1. 在 macOS / Windows 上的推荐做法

#### 1.1 修改配置：使用 `host.docker.internal`

当前 `configs/config.yaml` 中的关键部分：

```yaml
data:
  database:
    enable: true
    driver: mysql
    source: root:密码@tcp(127.0.0.1:3306)/test?timeout=15s&charset=utf8mb4&parseTime=True
  redis:
    enable: true
    addr: 127.0.0.1:6379
```

在容器内，`127.0.0.1` 指向容器本身，而不是宿主机，需要改为：

```yaml
data:
  database:
    enable: true
    driver: mysql
    # 关键：将 127.0.0.1 换成 host.docker.internal
    source: root:密码@tcp(host.docker.internal:3306)/test?timeout=15s&charset=utf8mb4&parseTime=True
  redis:
    enable: true
    addr: host.docker.internal:6379
```

> 说明：`host.docker.internal` 是 Docker Desktop 内置的 DNS 名称，指向宿主机。

#### 1.2 使用环境变量覆盖（推荐）

为了避免在配置文件里写死宿主机地址，可以使用环境变量覆盖：

```bash
docker run -d \
  --name sre \
  -p 8000:8000 \
  -p 8989:8989 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e REDIS_ADDR=host.docker.internal:6379 \
  sre:latest
```

对应的配置可以改为读取环境变量（示意）：

```yaml
data:
  database:
    source: root:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/test?timeout=15s&charset=utf8mb4&parseTime=True
  redis:
    addr: ${REDIS_ADDR}
```

> 实际读取环境变量由 `viper` + `LoadBootstrapFromViper` 负责，配置中写 `${VAR}` 表达意图，具体实现可以参考 `internal/config/kratos.go`。

---

### 2. 在 Linux 服务器上的几种方式

#### 2.1 方式一：使用宿主机实际 IP

假设宿主机 IP 为 `192.168.1.100`，MySQL / Redis 在宿主机上监听 `0.0.0.0`：

```yaml
data:
  database:
    source: root:密码@tcp(192.168.1.100:3306)/test?timeout=15s&charset=utf8mb4&parseTime=True
  redis:
    addr: 192.168.1.100:6379
```

启动容器时：

```bash
docker run -d \
  --name sre \
  -p 8000:8000 \
  -p 8989:8989 \
  sre:latest
```

#### 2.2 方式二：使用 `--network host`（仅 Linux）

在 Linux 上可以让容器和宿主机共用网络命名空间：

```bash
docker run -d \
  --name sre \
  --network host \
  sre:latest
```

此时容器内访问 `127.0.0.1:3306` 就等同于访问宿主机的 `127.0.0.1:3306`，配置文件可以保持为本地开发时的写法：

```yaml
data:
  database:
    source: root:密码@tcp(127.0.0.1:3306)/test?...
  redis:
    addr: 127.0.0.1:6379
```

> 注意：`--network host` 在 macOS / Windows 上基本无效，只在 Linux 上可用。

---

### 3. 宿主机 MySQL / Redis 的必要配置

无论采用哪种方式，宿主机上的 MySQL / Redis 必须允许来自容器的连接：

#### 3.1 MySQL

1. `my.cnf` 中的 `bind-address` 不能只绑定 `127.0.0.1`，需要允许外部访问，例如：
   ```ini
   bind-address = 0.0.0.0
   ```
2. 创建允许从容器网段访问的用户，例如：
   ```sql
   CREATE USER 'root'@'%' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON test.* TO 'root'@'%';
   FLUSH PRIVILEGES;
   ```
3. 防火墙开放 `3306` 端口（根据实际情况配置）。

#### 3.2 Redis

1. `redis.conf` 中配置：
   ```conf
   bind 0.0.0.0
   requirepass your_password   # 建议开启密码
   ```
2. 防火墙开放 `6379` 端口。

---

### 4. 连通性测试

#### 4.1 宿主机上测试

```bash
# MySQL
mysql -h 127.0.0.1 -P 3306 -u root -p

# Redis
redis-cli -h 127.0.0.1 -p 6379
```

#### 4.2 容器内测试

进入容器：

```bash
docker exec -it sre sh
```

在容器内测试端口连通性：

```bash
# 测试 MySQL 端口
nc -vz host.docker.internal 3306

# 测试 Redis 端口
nc -vz host.docker.internal 6379
```

（Linux 上如果使用宿主机 IP，则替换为对应 IP）

---

### 5. 推荐实践

1. **本地开发（Docker Desktop）**
   - 配置中将 `127.0.0.1` 换成 `host.docker.internal`。
   - 或使用环境变量控制 DB/Redis 地址。

2. **测试 / 预发环境**
   - MySQL / Redis 跑在独立实例（VM 或容器）上，用实际 IP + 端口访问。

3. **生产环境**
   - 一般不会直接访问宿主机 DB，而是访问集群内的 DB 服务（例如 RDS、K8s Service）。
   - 此文档方案主要用于本地开发和简单部署。

---

### 6. 配置示例对比

#### 本地直接跑（不使用 Docker）

```yaml
data:
  database:
    source: root:密码@tcp(127.0.0.1:3306)/test?...
  redis:
    addr: 127.0.0.1:6379
```

#### Docker + Docker Desktop（推荐）

```yaml
data:
  database:
    source: root:密码@tcp(host.docker.internal:3306)/test?...
  redis:
    addr: host.docker.internal:6379
```

#### Docker + Linux 服务器（使用宿主机 IP）

```yaml
data:
  database:
    source: root:密码@tcp(192.168.1.100:3306)/test?...
  redis:
    addr: 192.168.1.100:6379
```

---

### 7. 与项目现有配置的关系

- 当前 `configs/config.yaml` 中数据库和 Redis 默认指向 `127.0.0.1`：
  ```yaml
  data:
    database:
      source: root:...@tcp(127.0.0.1:3306)/test?...
    redis:
      addr: 127.0.0.1:6379
  ```
- 本地直接运行（`go run` 或 `./bin/sre`）时，这样配置是正确的。
- **容器化运行时，必须根据运行环境调整为 `host.docker.internal` 或宿主机 IP**，否则容器内无法访问宿主机的 DB/Redis。


