### 一、容器的底层实现：Linux 内核机制

#### 1. Namespaces（命名空间）
作用：隔离进程的“视图”，让容器内的进程只能看到受限的系统资源。
关键命名空间类型：
- PID Namespace
  - 隔离进程树：容器内进程的 PID 独立，从 1 开始编号（类似独立操作系统）。
  - 实现方式：

```bash
# 查看进程的 PID Namespace
ls -l /proc/<PID>/ns/pid
```
- Mount Namespace
  - 隔离文件系统挂载点：容器拥有独立的文件系统视图（如 / 根目录）。
```bash
# 创建一个新的 Mount Namespace
unshare --mount
mount -t tmpfs tmpfs /mnt  # 容器内挂载不影响宿主机
```
- Network Namespace
  - 隔离网络设备：容器拥有独立的 IP、端口、路由表等。
  - 实现方式：
```bash
# 创建新的 Network Namespace
ip netns add mynet
ip netns exec mynet ip a  # 查看容器内网络设备
```

- UTS Namespace

  - 隔离主机名和域名：容器可以设置独立的 hostname。
  - 示例：
``` bash
unshare --uts
hostname mycontainer  # 修改容器主机名
```
- IPC Namespace

  - 隔离进程间通信：如信号量、共享内存等。

- User Namespace

  - 隔离用户和用户组：容器内可使用虚拟 UID/GID，映射到宿主机的非特权用户。

  - 安全优势：容器进程以非 root 用户运行（即使容器内显示为 root）。

#### 2. cgroups（Control Groups）
作用：限制、统计和隔离进程组的资源使用。
关键子系统：

- CPU 子系统：限制 CPU 使用份额（如 cpu.shares）或硬性配额（如 cpu.cfs_period_us）。
- Memory 子系统：限制内存使用量，触发 OOM（Out-Of-Memory）时终止容器。
- Blkio 子系统：限制磁盘 I/O 带宽。
- Devices 子系统：控制容器能否访问特定设备（如 /dev/sda）。

实现方式：

- cgroups 的配置通过虚拟文件系统暴露在 /sys/fs/cgroup/ 目录下。
示例：限制容器内存为 512MB：
```bash
# 创建 cgroup
mkdir /sys/fs/cgroup/memory/mycontainer
echo 536870912 > /sys/fs/cgroup/memory/mycontainer/memory.limit_in_bytes
# 将进程 PID 加入 cgroup
echo <PID> > /sys/fs/cgroup/memory/mycontainer/cgroup.procs
```

#### 3. Union File System（联合文件系统）
作用：为容器提供分层文件系统，支持镜像的复用和快速启动。
- 常见实现：OverlayFS、AUFS、devicemapper。
示例（OverlayFS）：
```bash
mount -t overlay overlay \
  -o lowerdir=/lower,upperdir=/upper,workdir=/work \
  /merged
```
  - lowerdir：只读的基础镜像层（如 Ubuntu 基础镜像）。
  - upperdir：可写的容器层（存放运行时修改）。
  - merged：合并后的视图，供容器进程使用。

4. Capabilities（权能）
作用：精细化控制进程的权限，替代传统的 root 全权模式。

示例：禁止容器修改系统时间：

```bash
# 移除 CAP_SYS_TIME 权能
capsh --drop=cap_sys_time -- -c "date -s '2021-01-01'"
```
- 常见权能：
  - CAP_NET_BIND_SERVICE：允许绑定低端口（如 80）。
  - CAP_SYS_ADMIN：允许执行特权操作（如挂载文件系统）。

### 二、容器运行时的工作流程（以 Docker 为例）
- 镜像准备：
  - 基于 OverlayFS 分层下载镜像（如 docker pull nginx）。
- 创建容器：
  - 调用 runc 创建容器进程，为其分配独立的 Namespaces。
- 资源限制：
  - 通过 cgroups 设置 CPU、内存等限制（对应 docker run --cpus=2 --memory=1g）。
- 文件系统挂载：
  - 使用 OverlayFS 挂载容器文件系统，合并镜像层和可写层。
- 安全隔离：
  - 启用 AppArmor/SELinux 配置文件，限制容器访问敏感路径。
  - 通过 Seccomp 过滤系统调用（如禁止 mount 系统调用）。
- 网络配置：
  - 创建 veth pair，将容器加入独立的 Network Namespace，并通过 iptables 配置 NAT 规则。
### 三、容器与虚拟机的区别
|特性|容器|虚拟机|
|:---|---|---|
|隔离级别|进程级隔离（共享宿主机内核）|硬件级隔离（独立内核）|
|启动速度|秒级启动|分钟级启动|
|资源开销|低（仅需运行进程）|高（需运行完整操作系统）|
|安全性|依赖内核隔离，存在逃逸风险|强隔离，安全性更高|
|适用场景|微服务、快速伸缩、CI/CD|需要强隔离的环境（如多租户）|

### 四、调试与验证容器的 Linux 机制
####　1. 查看容器的 Namespaces
``` bash
# 查看容器进程的 PID（例如 1234）
docker inspect --format '{{.State.Pid}}' <容器ID>

# 查看该进程的 Namespace 信息
ls -l /proc/1234/ns/
```
#### 2. 检查 cgroups 配置
```bash
# 查看容器的 cgroup 路径
cat /proc/1234/cgroup

# 查看内存限制
cat /sys/fs/cgroup/memory/docker/<容器ID>/memory.limit_in_bytes
```
#### 3. 验证 Capabilities
```bash
# 查看容器进程的权能
cat /proc/1234/status | grep Cap
```
### 五、总结
容器本质是 通过 Linux 内核的 Namespaces 和 cgroups 实现的进程级隔离：
- 隔离性：Namespaces 提供资源视图隔离，cgroups 提供资源使用限制。
- 轻量化：共享宿主机内核，无需虚拟化硬件。
- 安全性：依赖权能（Capabilities）、Seccomp、AppArmor 等机制补强隔离。