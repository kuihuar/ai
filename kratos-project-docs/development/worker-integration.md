# 主应用集成 Worker 最佳实践

## 概述

本文档说明如何在主应用 `cmd/sre` 中集成 `daemon-worker` 或 `cron-worker` 的功能，提供配置驱动、优雅启动/停止、资源隔离等最佳实践。

## 目录

- [设计原则](#设计原则)
- [方案对比](#方案对比)
- [推荐方案：配置驱动集成](#推荐方案配置驱动集成)
- [实现步骤](#实现步骤)
- [配置示例](#配置示例)
- [部署场景](#部署场景)
- [注意事项](#注意事项)

---

## 设计原则

### 1. 配置驱动
- ✅ 通过配置文件控制是否启用 worker
- ✅ 支持独立部署和集成部署两种模式
- ✅ 配置变更无需修改代码

### 2. 生命周期管理
- ✅ 使用 Kratos 生命周期钩子统一管理
- ✅ 优雅启动和停止
- ✅ 错误处理和恢复

### 3. 资源隔离
- ✅ Worker 错误不影响主服务
- ✅ 独立的 goroutine 运行
- ✅ 独立的日志标识

### 4. 可维护性
- ✅ 代码结构清晰，易于扩展
- ✅ 遵循项目架构规范
- ✅ 最小化代码重复

---

## 方案对比

### 方案 A：独立部署（当前方式）

**优点**：
- ✅ 完全隔离，worker 崩溃不影响主服务
- ✅ 独立扩缩容
- ✅ 独立监控和日志
- ✅ 资源隔离（CPU、内存）

**缺点**：
- ❌ 需要部署多个进程
- ❌ 配置需要同步
- ❌ 资源占用更多

**适用场景**：
- 生产环境
- 需要独立扩缩容
- 对稳定性要求高

### 方案 B：集成部署（本文档方案）

**优点**：
- ✅ 单进程部署，简化运维
- ✅ 共享配置和连接池
- ✅ 资源占用更少
- ✅ 适合小规模应用

**缺点**：
- ❌ Worker 错误可能影响主服务（需要良好隔离）
- ❌ 无法独立扩缩容
- ❌ 日志混合（需要良好标识）

**适用场景**：
- 开发/测试环境
- 小规模应用
- 资源受限环境

### 方案 C：混合模式（推荐）

**优点**：
- ✅ 通过配置选择部署模式
- ✅ 灵活性最高
- ✅ 可以逐步迁移

**实现**：
- 配置中增加 `worker.mode` 字段
- `mode: "standalone"` → 独立部署
- `mode: "embedded"` → 集成部署

---

## 推荐方案：配置驱动集成

### 架构设计

```
┌─────────────────────────────────────────────────┐
│              cmd/sre (主应用)                    │
│  ┌──────────────────────────────────────────┐   │
│  │  gRPC Server + HTTP Server               │   │
│  └──────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────┐   │
│  │  Worker Manager (可选)                    │   │
│  │  ├── Daemon Worker Set (可选)            │   │
│  │  └── Cron Manager (可选)                 │   │
│  └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

### 配置结构

在 `internal/conf/conf.proto` 中添加 Worker 配置：

```protobuf
message Bootstrap {
  Server server = 1;
  Data data = 2;
  Registry registry = 3;
  ConfigCenter config_center = 4;
  Log log = 5;
  Service service = 6;
  Worker worker = 7;  // 新增：Worker 配置
}

message Worker {
  bool enable = 1;              // 是否启用 Worker（总开关）
  string mode = 2;              // 模式: "standalone" | "embedded"
  DaemonWorker daemon = 3;      // Daemon Worker 配置
  CronWorker cron = 4;          // Cron Worker 配置
}

message DaemonWorker {
  bool enable = 1;              // 是否启用 Daemon Worker
  repeated string jobs = 2;     // 启用的 job 列表，如: ["table-consumer-ants"]
}

message CronWorker {
  bool enable = 1;              // 是否启用 Cron Worker
  string timezone = 2;          // 时区，默认: "Asia/Shanghai"
  repeated string jobs = 3;     // 启用的 job 列表，如: ["sync_user"]
}
```

### 配置示例

```yaml
# configs/config.yaml
worker:
  enable: true                    # 总开关
  mode: "embedded"                # 模式: standalone | embedded
  daemon:
    enable: true                   # 启用 Daemon Worker
    jobs:                          # 启用的 daemon jobs
      - "table-consumer-ants"
  cron:
    enable: true                   # 启用 Cron Worker
    timezone: "Asia/Shanghai"
    jobs:                          # 启用的 cron jobs
      - "sync_user"
```

---

## 实现步骤

### 步骤 1：更新配置定义

#### 1.1 更新 `internal/conf/conf.proto`

```protobuf
message Bootstrap {
  // ... 现有字段 ...
  Worker worker = 7;
}

message Worker {
  bool enable = 1;
  string mode = 2;
  DaemonWorker daemon = 3;
  CronWorker cron = 4;
}

message DaemonWorker {
  bool enable = 1;
  repeated string jobs = 2;
}

message CronWorker {
  bool enable = 1;
  string timezone = 2;
  repeated string jobs = 3;
}
```

#### 1.2 重新生成配置代码

```bash
make config
# 或
protoc --go_out=. --go_opt=paths=source_relative internal/conf/conf.proto
```

### 步骤 2：创建 Worker Manager

#### 2.1 创建 `internal/app/worker/manager.go`

```go
package worker

import (
	"context"
	"sre/internal/biz/cron"
	"sre/internal/biz/daemon"
	"sre/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// Manager 管理所有 Worker（Daemon 和 Cron）
type Manager struct {
	logger        log.Logger
	daemonJobSet  *DaemonJobSet
	cronManager   *cron.Manager
	config        *conf.Worker
}

// DaemonJobSet 管理多个 daemon jobs
type DaemonJobSet struct {
	Jobs []daemon.DaemonJob
}

// NewManager 创建 Worker Manager
func NewManager(
	logger log.Logger,
	config *conf.Worker,
	daemonJobSet *DaemonJobSet,
	cronManager *cron.Manager,
) *Manager {
	return &Manager{
		logger:       logger,
		daemonJobSet: daemonJobSet,
		cronManager:  cronManager,
		config:       config,
	}
}

// Start 启动所有启用的 Worker
func (m *Manager) Start(ctx context.Context) error {
	logHelper := log.NewHelper(m.logger)

	if m.config == nil || !m.config.Enable {
		logHelper.Info("worker is disabled, skipping worker startup")
		return nil
	}

	// 启动 Daemon Worker
	if m.config.Daemon != nil && m.config.Daemon.Enable {
		if err := m.startDaemonWorkers(ctx, logHelper); err != nil {
			return err
		}
	}

	// 启动 Cron Worker
	if m.config.Cron != nil && m.config.Cron.Enable {
		if err := m.startCronWorkers(ctx, logHelper); err != nil {
			return err
		}
	}

	return nil
}

// Stop 停止所有 Worker
func (m *Manager) Stop(ctx context.Context) error {
	logHelper := log.NewHelper(m.logger)

	if m.config == nil || !m.config.Enable {
		return nil
	}

	// 停止 Daemon Worker
	if m.config.Daemon != nil && m.config.Daemon.Enable && m.daemonJobSet != nil {
		logHelper.Info("stopping daemon workers")
		for _, job := range m.daemonJobSet.Jobs {
			if err := job.Stop(); err != nil {
				logHelper.Errorf("failed to stop daemon job %s: %v", job.Name(), err)
			}
		}
	}

	// 停止 Cron Worker
	if m.config.Cron != nil && m.config.Cron.Enable && m.cronManager != nil {
		logHelper.Info("stopping cron workers")
		// Cron Manager 的停止逻辑（如果有）
	}

	return nil
}

// startDaemonWorkers 启动 Daemon Workers
func (m *Manager) startDaemonWorkers(ctx context.Context, logHelper *log.Helper) error {
	if m.daemonJobSet == nil || len(m.daemonJobSet.Jobs) == 0 {
		logHelper.Warn("no daemon jobs registered")
		return nil
	}

	// 根据配置过滤启用的 jobs
	enabledJobs := m.filterEnabledDaemonJobs(m.config.Daemon.Jobs)

	logHelper.Infof("starting %d daemon job(s)", len(enabledJobs))

	for _, job := range enabledJobs {
		job := job // 避免闭包问题
		go func() {
			logHelper.Infof("starting daemon job: %s", job.Name())
			if err := job.Run(ctx); err != nil {
				logHelper.Errorf("daemon job %s stopped with error: %v", job.Name(), err)
			} else {
				logHelper.Infof("daemon job %s stopped gracefully", job.Name())
			}
		}()
	}

	return nil
}

// startCronWorkers 启动 Cron Workers
func (m *Manager) startCronWorkers(ctx context.Context, logHelper *log.Helper) error {
	if m.cronManager == nil {
		logHelper.Warn("cron manager is nil")
		return nil
	}

	logHelper.Info("starting cron workers")

	// 在独立的 goroutine 中启动 Cron Manager
	go func() {
		if err := m.cronManager.Start(ctx); err != nil {
			logHelper.Errorf("cron manager stopped with error: %v", err)
		}
	}()

	return nil
}

// filterEnabledDaemonJobs 根据配置过滤启用的 daemon jobs
func (m *Manager) filterEnabledDaemonJobs(enabledJobNames []string) []daemon.DaemonJob {
	if len(enabledJobNames) == 0 {
		// 如果没有配置，启用所有 jobs
		return m.daemonJobSet.Jobs
	}

	enabledMap := make(map[string]bool)
	for _, name := range enabledJobNames {
		enabledMap[name] = true
	}

	var filtered []daemon.DaemonJob
	for _, job := range m.daemonJobSet.Jobs {
		if enabledMap[job.Name()] {
			filtered = append(filtered, job)
		}
	}

	return filtered
}
```

### 步骤 3：更新 `cmd/sre/app.go`

#### 3.1 修改 `newApp` 函数

```go
package main

import (
	"context"
	"sre/internal/app/worker"
	"sre/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func newApp(
	logger log.Logger,
	gs *grpc.Server,
	hs *http.Server,
	registrar kratosRegistry.Registrar,
	dingTalkEventService *service.DingTalkEventService,
	workerManager *worker.Manager,  // 新增
	workerConfig *conf.Worker,       // 新增
) *kratos.App {
	logHelper := log.NewHelper(logger)

	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	}

	// 注册中心
	if registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	// 钉钉事件服务
	if dingTalkEventService != nil {
		opts = append(opts,
			kratos.BeforeStart(func(ctx context.Context) error {
				dingTalkEventService.Start()
				return nil
			}),
			kratos.BeforeStop(func(ctx context.Context) error {
				dingTalkEventService.Stop()
				return nil
			}),
		)
	}

	// Worker Manager（如果启用）
	if workerConfig != nil && workerConfig.Enable && workerConfig.Mode == "embedded" {
		opts = append(opts,
			kratos.BeforeStart(func(ctx context.Context) error {
				logHelper.Info("starting embedded workers")
				return workerManager.Start(ctx)
			}),
			kratos.BeforeStop(func(ctx context.Context) error {
				logHelper.Info("stopping embedded workers")
				return workerManager.Stop(ctx)
			}),
		)
	}

	return kratos.New(opts...)
}
```

### 步骤 4：更新 Wire 配置

#### 4.1 更新 `cmd/sre/wire.go`

```go
//go:build wireinject
// +build wireinject

package main

import (
	"sre/internal/app/worker"
	"sre/internal/biz"
	"sre/internal/biz/cron"
	"sre/internal/biz/cron/jobs"
	"sre/internal/biz/daemon"
	"sre/internal/conf"
	"sre/internal/data"
	"sre/internal/registry"
	"sre/internal/server"
	"sre/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Server,
	*conf.Data,
	*conf.Registry,
	*conf.Service,
	*conf.Worker,  // 新增
	log.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		// 基础设施
		data.NewDB,
		data.NewData,
		registry.NewRegistry,
		// 数据层
		data.NewUserRepo,
		data.NewDingTalkClient,
		data.NewExternalUserService,
		// 业务层
		biz.ProvideBusinessConfig,
		biz.NewUserUsecase,
		biz.NewDingTalkEventUsecase,
		// 服务层
		service.NewUserService,
		service.NewDingTalkEventService,
		// 服务器
		server.NewGRPCServer,
		server.NewHTTPServer,
		// Worker 相关（可选，根据配置决定是否创建）
		provideDaemonJobSet,      // 新增
		provideCronManager,       // 新增
		worker.NewManager,        // 新增
		// 应用
		newApp,
	))
}

// provideDaemonJobSet 提供 Daemon Job Set（可选）
func provideDaemonJobSet(
	tableConsumerDaemon *daemon.TableConsumerDaemonAnts,
) *worker.DaemonJobSet {
	jobs := []daemon.DaemonJob{}
	
	if tableConsumerDaemon != nil {
		jobs = append(jobs, tableConsumerDaemon)
	}
	
	return &worker.DaemonJobSet{
		Jobs: jobs,
	}
}

// provideCronManager 提供 Cron Manager（可选）
func provideCronManager(
	config *conf.Worker,
	logger log.Logger,
	syncUserJob *jobs.SyncUserJob,
) (*cron.Manager, error) {
	if config == nil || config.Cron == nil || !config.Cron.Enable {
		return nil, nil
	}

	timezone := config.Cron.Timezone
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}

	manager, err := cron.NewManager(logger, timezone)
	if err != nil {
		return nil, err
	}

	// 注册启用的 jobs
	if syncUserJob != nil {
		if err := manager.RegisterJob(syncUserJob); err != nil {
			return nil, err
		}
	}

	return manager, nil
}
```

#### 4.2 更新 `cmd/sre/wire.go` 中的 Daemon Job 创建

需要确保 Daemon Job 的创建逻辑在 Wire 中：

```go
// 在 wire.Build 中添加
daemon.NewTableConsumerDaemonAntsForBiz,  // 或相应的构造函数
```

### 步骤 5：更新配置文件

在 `configs/config.yaml` 中添加 Worker 配置：

```yaml
worker:
  enable: true
  mode: "embedded"  # 或 "standalone"
  daemon:
    enable: true
    jobs:
      - "table-consumer-ants"
  cron:
    enable: true
    timezone: "Asia/Shanghai"
    jobs:
      - "sync_user"
```

### 步骤 6：生成 Wire 代码

```bash
cd cmd/sre
wire
```

---

## 配置示例

### 场景 1：仅启用 Daemon Worker

```yaml
worker:
  enable: true
  mode: "embedded"
  daemon:
    enable: true
    jobs:
      - "table-consumer-ants"
  cron:
    enable: false
```

### 场景 2：仅启用 Cron Worker

```yaml
worker:
  enable: true
  mode: "embedded"
  daemon:
    enable: false
  cron:
    enable: true
    timezone: "Asia/Shanghai"
    jobs:
      - "sync_user"
```

### 场景 3：同时启用两者

```yaml
worker:
  enable: true
  mode: "embedded"
  daemon:
    enable: true
    jobs:
      - "table-consumer-ants"
  cron:
    enable: true
    timezone: "Asia/Shanghai"
    jobs:
      - "sync_user"
```

### 场景 4：完全禁用 Worker

```yaml
worker:
  enable: false
```

### 场景 5：独立部署模式

```yaml
worker:
  enable: false  # 主应用不启用
  mode: "standalone"  # 通过独立的 daemon-worker 和 cron-worker 部署
```

---

## 部署场景

### 场景 A：开发/测试环境（集成部署）

**配置**：
```yaml
worker:
  enable: true
  mode: "embedded"
  daemon:
    enable: true
  cron:
    enable: true
```

**部署**：
```bash
# 单进程运行
go run cmd/sre/main.go cmd/sre/wire_gen.go -conf configs
```

**优点**：
- 简化部署
- 快速开发测试
- 资源占用少

### 场景 B：生产环境（独立部署）

**配置**：
```yaml
# configs/config.yaml (主应用)
worker:
  enable: false

# configs/daemon-worker.yaml
worker:
  enable: true
  mode: "standalone"
  daemon:
    enable: true

# configs/cron-worker.yaml
worker:
  enable: true
  mode: "standalone"
  cron:
    enable: true
```

**部署**：
```bash
# 主应用
./bin/sre -conf configs

# Daemon Worker
./bin/daemon-worker -conf configs

# Cron Worker
./bin/cron-worker -conf configs
```

**优点**：
- 完全隔离
- 独立扩缩容
- 高可用

### 场景 C：混合部署（推荐）

**配置**：
```yaml
# 主应用：只运行 API 服务
worker:
  enable: false

# 部分 Worker 集成在主应用
# 部分 Worker 独立部署（通过环境变量或配置中心控制）
```

**部署**：
- 根据负载和重要性选择集成或独立部署
- 关键 Worker 独立部署
- 次要 Worker 集成部署

---

## 注意事项

### 1. 错误隔离

**问题**：Worker 错误可能影响主服务

**解决方案**：
- ✅ 使用独立的 goroutine 运行 Worker
- ✅ Worker 错误只记录日志，不 panic
- ✅ 使用 recover 捕获 panic

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            logHelper.Errorf("daemon job %s panicked: %v", job.Name(), r)
        }
    }()
    // ... worker logic
}()
```

### 2. 资源竞争

**问题**：Worker 和主服务可能竞争资源（数据库连接池、CPU 等）

**解决方案**：
- ✅ 合理配置连接池大小
- ✅ 使用限流控制 Worker 处理速度
- ✅ 监控资源使用情况

### 3. 日志标识

**问题**：Worker 日志和主服务日志混合

**解决方案**：
- ✅ 在日志中添加 Worker 标识
- ✅ 使用不同的日志文件（如果配置了文件输出）
- ✅ 使用结构化日志字段区分

```go
logHelper.WithFields(map[string]interface{}{
    "worker": "daemon",
    "job": job.Name(),
}).Info("starting daemon job")
```

### 4. 配置同步

**问题**：集成部署时，Worker 配置和主服务配置需要同步

**解决方案**：
- ✅ 使用统一的配置源（配置中心）
- ✅ 配置验证和默认值
- ✅ 配置变更时优雅重启

### 5. 优雅停止

**问题**：主服务停止时，Worker 需要优雅停止

**解决方案**：
- ✅ 使用 Kratos `BeforeStop` 钩子
- ✅ Worker 实现 `Stop()` 方法
- ✅ 设置合理的超时时间

### 6. 监控和可观测性

**问题**：需要监控 Worker 的运行状态

**解决方案**：
- ✅ 添加健康检查接口
- ✅ 暴露 Prometheus 指标
- ✅ 使用 OpenTelemetry 追踪

---

## 最佳实践总结

### ✅ 推荐做法

1. **配置驱动**：通过配置文件控制 Worker 的启用/禁用
2. **生命周期管理**：使用 Kratos 生命周期钩子统一管理
3. **错误隔离**：Worker 错误不影响主服务
4. **优雅停止**：确保 Worker 能够优雅停止
5. **日志标识**：清晰标识 Worker 日志
6. **资源监控**：监控 Worker 的资源使用

### ❌ 避免做法

1. **硬编码**：不要在代码中硬编码 Worker 的启用/禁用
2. **阻塞主服务**：不要让 Worker 阻塞主服务的启动
3. **忽略错误**：不要忽略 Worker 的错误，至少记录日志
4. **资源泄漏**：确保 Worker 停止时释放所有资源
5. **配置混乱**：不要在不同地方重复配置 Worker

---

## 相关文档

- [应用启动流程](../../architecture/application-startup-flow.md)
- [Cron 定时任务集成指南](./cron-integration.md)
- [Daemon Worker 实现指南](../internal/biz/daemon/README.md)
- [配置管理](../operations/config-management.md)

---

## 示例代码

完整的示例代码请参考：
- `internal/app/worker/manager.go` - Worker Manager 实现
- `cmd/sre/app.go` - 主应用集成示例
- `cmd/sre/wire.go` - Wire 依赖注入配置

