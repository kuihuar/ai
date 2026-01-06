# Cron 定时任务集成指南

## 概述

本文档介绍如何在 Kratos 项目中集成定时任务（Cron Job）功能。定时任务通常用于执行周期性任务，如数据同步、清理、统计等。

## 架构设计

### 推荐方案：独立应用

在 Kratos 多应用架构中，**推荐将定时任务作为独立应用**，放在 `cmd/` 目录下，例如 `cmd/cron-worker`。

#### 优点
- **独立部署**：定时任务可以独立部署和扩展
- **资源隔离**：不影响主服务的性能和稳定性
- **灵活调度**：可以为不同类型的任务创建不同的 worker 应用
- **易于维护**：任务代码独立，便于管理和监控

#### 项目结构

```
sre/
├── cmd/
│   ├── sre/              # 主服务应用
│   ├── user-consumer/    # Kafka 消费者应用
│   └── cron-worker/      # 定时任务应用（新增）
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── internal/
│   ├── biz/
│   │   └── cron/         # 定时任务业务逻辑（新增）
│   │       ├── job.go    # 任务接口定义
│   │       └── jobs/     # 具体任务实现
│   │           ├── sync_user.go
│   │           └── cleanup_log.go
│   ├── data/
│   └── service/
└── configs/
    └── cron-worker.yaml  # 定时任务应用配置
```

## 技术选型

### Cron 库选择

推荐使用 `github.com/robfig/cron/v3`，这是 Go 生态中最成熟和广泛使用的 cron 库。

#### 特点
- 支持标准 cron 表达式（5 位和 6 位）
- 支持秒级精度（6 位表达式）
- 支持时区配置
- 支持任务链和依赖
- 活跃维护，社区支持好

#### 安装

```bash
go get github.com/robfig/cron/v3
```

## 实现步骤

### 步骤 1: 定义任务接口

在 `internal/biz/cron/` 目录下定义任务接口：

```go
// internal/biz/cron/job.go
package cron

import (
	"context"
	
	"github.com/go-kratos/kratos/v2/log"
)

// Job 定时任务接口
type Job interface {
	// Name 返回任务名称
	Name() string
	
	// Spec 返回 cron 表达式
	// 标准格式：秒 分 时 日 月 星期
	// 例如："0 0 2 * * *" 表示每天凌晨 2 点执行
	Spec() string
	
	// Run 执行任务
	Run(ctx context.Context) error
}

// JobRunner 任务运行器
type JobRunner struct {
	logger log.Logger
}

// NewJobRunner 创建任务运行器
func NewJobRunner(logger log.Logger) *JobRunner {
	return &JobRunner{
		logger: logger,
	}
}

// Run 运行任务，包含错误处理和日志记录
func (r *JobRunner) Run(job Job) func() {
	return func() {
		ctx := context.Background()
		logHelper := log.NewHelper(r.logger)
		
		logHelper.Infof("starting cron job: %s", job.Name())
		
		if err := job.Run(ctx); err != nil {
			logHelper.Errorf("cron job %s failed: %v", job.Name(), err)
			return
		}
		
		logHelper.Infof("cron job %s completed successfully", job.Name())
	}
}
```

### 步骤 2: 实现具体任务

在 `internal/biz/cron/jobs/` 目录下实现具体任务：

```go
// internal/biz/cron/jobs/sync_user.go
package jobs

import (
	"context"
	
	"sre/internal/biz/cron"
	"sre/internal/biz/user"
	
	"github.com/go-kratos/kratos/v2/log"
)

// SyncUserJob 同步用户任务
type SyncUserJob struct {
	userUsecase *user.UserUsecase
	logger      log.Logger
}

// NewSyncUserJob 创建同步用户任务
func NewSyncUserJob(userUsecase *user.UserUsecase, logger log.Logger) *SyncUserJob {
	return &SyncUserJob{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

// Name 返回任务名称
func (j *SyncUserJob) Name() string {
	return "sync-user"
}

// Spec 返回 cron 表达式：每天凌晨 2 点执行
func (j *SyncUserJob) Spec() string {
	return "0 0 2 * * *"
}

// Run 执行任务
func (j *SyncUserJob) Run(ctx context.Context) error {
	logHelper := log.NewHelper(j.logger)
	logHelper.Info("syncing users from external service...")
	
	// 调用业务逻辑
	// if err := j.userUsecase.SyncUsers(ctx); err != nil {
	//     return err
	// }
	
	logHelper.Info("user sync completed")
	return nil
}

// 确保实现了 cron.Job 接口
var _ cron.Job = (*SyncUserJob)(nil)
```

```go
// internal/biz/cron/jobs/cleanup_log.go
package jobs

import (
	"context"
	"time"
	
	"sre/internal/biz/cron"
	
	"github.com/go-kratos/kratos/v2/log"
)

// CleanupLogJob 清理日志任务
type CleanupLogJob struct {
	logger log.Logger
}

// NewCleanupLogJob 创建清理日志任务
func NewCleanupLogJob(logger log.Logger) *CleanupLogJob {
	return &CleanupLogJob{
		logger: logger,
	}
}

// Name 返回任务名称
func (j *CleanupLogJob) Name() string {
	return "cleanup-log"
}

// Spec 返回 cron 表达式：每天凌晨 3 点执行
func (j *CleanupLogJob) Spec() string {
	return "0 0 3 * * *"
}

// Run 执行任务
func (j *CleanupLogJob) Run(ctx context.Context) error {
	logHelper := log.NewHelper(j.logger)
	
	// 清理 30 天前的日志
	cutoffTime := time.Now().AddDate(0, 0, -30)
	logHelper.Infof("cleaning up logs before %v", cutoffTime)
	
	// 实现清理逻辑
	// ...
	
	logHelper.Info("log cleanup completed")
	return nil
}

// 确保实现了 cron.Job 接口
var _ cron.Job = (*CleanupLogJob)(nil)
```

### 步骤 3: 创建 Cron 管理器

在 `internal/biz/cron/` 目录下创建 Cron 管理器：

```go
// internal/biz/cron/manager.go
package cron

import (
	"context"
	"fmt"
	"time"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

// Manager Cron 任务管理器
type Manager struct {
	cron    *cron.Cron
	jobs    []Job
	runner  *JobRunner
	logger  log.Logger
}

// NewManager 创建 Cron 管理器
func NewManager(logger log.Logger, timezone string) (*Manager, error) {
	// 解析时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}
	
	// 创建 cron 实例，支持秒级精度
	c := cron.New(
		cron.WithSeconds(),           // 支持秒级精度（6 位表达式）
		cron.WithLocation(loc),       // 设置时区
		cron.WithChain(               // 添加恢复链，防止 panic
			cron.Recover(cron.DefaultLogger),
		),
	)
	
	return &Manager{
		cron:   c,
		jobs:   make([]Job, 0),
		runner: NewJobRunner(logger),
		logger: logger,
	}, nil
}

// RegisterJob 注册任务
func (m *Manager) RegisterJob(job Job) error {
	logHelper := log.NewHelper(m.logger)
	
	// 注册任务（AddFunc 会自动验证 cron 表达式）
	_, err := m.cron.AddFunc(job.Spec(), m.runner.Run(job))
	if err != nil {
		return fmt.Errorf("failed to register job %s: %w", job.Name(), err)
	}
	
	m.jobs = append(m.jobs, job)
	logHelper.Infof("registered cron job: %s, spec: %s", job.Name(), job.Spec())
	
	return nil
}

// Start 启动 Cron 管理器
func (m *Manager) Start(ctx context.Context) error {
	logHelper := log.NewHelper(m.logger)
	
	if len(m.jobs) == 0 {
		logHelper.Warn("no cron jobs registered")
		return nil
	}
	
	logHelper.Infof("starting cron manager with %d jobs", len(m.jobs))
	m.cron.Start()
	
	// 等待上下文取消
	<-ctx.Done()
	
	logHelper.Info("stopping cron manager...")
	m.cron.Stop()
	
	// 等待所有运行中的任务完成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		logHelper.Warn("timeout waiting for cron jobs to finish")
	case <-m.cron.Stop().Done():
		logHelper.Info("all cron jobs stopped")
	}
	
	return nil
}

// GetJobs 获取所有注册的任务
func (m *Manager) GetJobs() []Job {
	return m.jobs
}
```

### 步骤 4: 创建应用入口

创建 `cmd/cron-worker/main.go`：

```go
// cmd/cron-worker/main.go
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"sre/internal/biz/cron"
	"sre/internal/biz/cron/jobs"
	"sre/internal/config"
	"sre/internal/logger"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	_ "go.uber.org/automaxprocs"
)

var (
	Name    string
	Version string
	flagconf string
	id, _   = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	
	// ============================================
	// 加载配置
	// ============================================
	bootstrap, err := config.LoadBootstrapWithViper(flagconf)
	if err != nil {
		panic(err)
	}
	
	// ============================================
	// 初始化 logger
	// ============================================
	zapLogger := logger.NewZapLoggerWithConfig(
		bootstrap.Log.Level,
		bootstrap.Log.Format,
		bootstrap.Log.OutputPaths,
		Name,
		id,
		Version,
	)
	
	logger := log.With(zapLogger,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	
	// ============================================
	// 初始化依赖
	// ============================================
	cleanup, err := wireApp(bootstrap.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	
	// ============================================
	// 创建 Cron 管理器
	// ============================================
	timezone := "Asia/Shanghai" // 可以从配置读取
	if bootstrap.Cron != nil && bootstrap.Cron.Timezone != "" {
		timezone = bootstrap.Cron.Timezone
	}
	
	manager, err := cron.NewManager(logger, timezone)
	if err != nil {
		panic(err)
	}
	
	// ============================================
	// 注册任务
	// ============================================
	// 这里需要根据实际业务注入依赖
	// userUsecase := ... // 从 wire 注入
	// syncUserJob := jobs.NewSyncUserJob(userUsecase, logger)
	// if err := manager.RegisterJob(syncUserJob); err != nil {
	//     panic(err)
	// }
	
	cleanupLogJob := jobs.NewCleanupLogJob(logger)
	if err := manager.RegisterJob(cleanupLogJob); err != nil {
		panic(err)
	}
	
	// ============================================
	// 启动 Cron 管理器
	// ============================================
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// 处理优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.NewHelper(logger).Info("received shutdown signal, stopping cron worker...")
		cancel()
	}()
	
	log.NewHelper(logger).Info("starting cron worker...")
	if err := manager.Start(ctx); err != nil {
		log.NewHelper(logger).Errorf("cron manager error: %v", err)
	}
	
	log.NewHelper(logger).Info("cron worker stopped")
}
```

### 步骤 5: 配置 Wire 依赖注入

创建 `cmd/cron-worker/wire.go`：

```go
// cmd/cron-worker/wire.go
//go:build wireinject
// +build wireinject

package main

import (
	"sre/internal/biz"
	"sre/internal/conf"
	"sre/internal/data"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp 初始化应用依赖
func wireApp(c *conf.Data, logger log.Logger) (func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		// 可以添加 cron 相关的 ProviderSet
	))
}
```

运行 Wire 生成代码：

```bash
cd cmd/cron-worker
go generate
```

### 步骤 6: 扩展配置定义

如果需要配置 Cron 相关参数，在 `internal/conf/conf.proto` 中添加：

```protobuf
// internal/conf/conf.proto
message Bootstrap {
  // ... 现有配置 ...
  
  Cron cron = 7;  // Cron 配置
}

message Cron {
  string timezone = 1;  // 时区，默认 Asia/Shanghai
}
```

重新生成配置代码：

```bash
make config
```

### 步骤 7: 配置文件

创建 `configs/cron-worker.yaml`：

```yaml
# configs/cron-worker.yaml
data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/test?parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s

cron:
  timezone: Asia/Shanghai  # 时区配置

log:
  level: info
  format: json
  output_paths:
    - stdout
```

## Cron 表达式说明

### 标准格式（6 位，支持秒级）

```
秒 分 时 日 月 星期
*  *  *  *  *  *
```

### 常用表达式示例

| 表达式 | 说明 |
|--------|------|
| `0 0 2 * * *` | 每天凌晨 2 点执行 |
| `0 0 */2 * * *` | 每 2 小时执行一次 |
| `0 */5 * * * *` | 每 5 分钟执行一次 |
| `0 0 0 * * 0` | 每周日凌晨执行 |
| `0 0 0 1 * *` | 每月 1 号凌晨执行 |
| `0 0 9-17 * * 1-5` | 工作日上午 9 点到下午 5 点，每小时执行 |
| `0 0 0,12 * * *` | 每天 0 点和 12 点执行 |

### 特殊字符

- `*`：匹配所有值
- `,`：指定多个值，如 `0,30` 表示 0 分和 30 分
- `-`：指定范围，如 `9-17` 表示 9 到 17
- `/`：指定步长，如 `*/5` 表示每 5 个单位

## 最佳实践

### 1. 任务设计原则

- **单一职责**：每个任务只做一件事
- **幂等性**：任务可以安全地重复执行
- **错误处理**：任务内部要妥善处理错误，避免影响其他任务
- **超时控制**：长时间运行的任务要设置超时

```go
// 示例：带超时控制的任务
func (j *SyncUserJob) Run(ctx context.Context) error {
	// 设置任务超时时间为 10 分钟
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	
	// 执行任务逻辑
	return j.userUsecase.SyncUsers(ctx)
}
```

### 2. 日志记录

- 任务开始和结束都要记录日志
- 记录任务执行时间
- 错误要记录详细上下文

```go
func (r *JobRunner) Run(job Job) func() {
	return func() {
		startTime := time.Now()
		ctx := context.Background()
		logHelper := log.NewHelper(r.logger)
		
		logHelper.Infof("starting cron job: %s", job.Name())
		
		if err := job.Run(ctx); err != nil {
			logHelper.Errorf("cron job %s failed after %v: %v", 
				job.Name(), time.Since(startTime), err)
			return
		}
		
		logHelper.Infof("cron job %s completed successfully in %v", 
			job.Name(), time.Since(startTime))
	}
}
```

### 3. 监控和告警

- 记录任务执行次数和成功率
- 任务失败时发送告警
- 监控任务执行时间

```go
// 可以集成 Prometheus metrics
func (r *JobRunner) Run(job Job) func() {
	return func() {
		startTime := time.Now()
		// ... 执行任务 ...
		
		// 记录 metrics
		// jobDuration.WithLabelValues(job.Name()).Observe(time.Since(startTime).Seconds())
		// jobCounter.WithLabelValues(job.Name(), "success").Inc()
	}
}
```

### 4. 分布式锁（可选）

如果多个实例运行，可以使用分布式锁确保任务只在一个实例执行：

```go
// 使用 Redis 分布式锁
func (j *SyncUserJob) Run(ctx context.Context) error {
	lockKey := fmt.Sprintf("cron:lock:%s", j.Name())
	
	// 尝试获取锁，超时时间 5 分钟
	lock, err := j.redisClient.SetNX(ctx, lockKey, "locked", 5*time.Minute).Result()
	if err != nil || !lock {
		return fmt.Errorf("failed to acquire lock for job %s", j.Name())
	}
	defer j.redisClient.Del(ctx, lockKey)
	
	// 执行任务逻辑
	return j.userUsecase.SyncUsers(ctx)
}
```

### 5. 任务依赖

如果任务之间有依赖关系，可以使用任务链：

```go
// 先执行任务 A，成功后再执行任务 B
func (m *Manager) RegisterJobChain(jobs ...Job) error {
	var prevJob Job
	for _, job := range jobs {
		if prevJob != nil {
			// 创建依赖任务
			dependentJob := &DependentJob{
				Job:      job,
				DependsOn: prevJob,
			}
			if err := m.RegisterJob(dependentJob); err != nil {
				return err
			}
		} else {
			if err := m.RegisterJob(job); err != nil {
				return err
			}
		}
		prevJob = job
	}
	return nil
}
```

## 测试

### 单元测试

```go
// internal/biz/cron/jobs/sync_user_test.go
package jobs

import (
	"context"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestSyncUserJob_Name(t *testing.T) {
	job := NewSyncUserJob(nil, nil)
	assert.Equal(t, "sync-user", job.Name())
}

func TestSyncUserJob_Spec(t *testing.T) {
	job := NewSyncUserJob(nil, nil)
	assert.Equal(t, "0 0 2 * * *", job.Spec())
}

func TestSyncUserJob_Run(t *testing.T) {
	// Mock userUsecase
	// job := NewSyncUserJob(mockUserUsecase, logger)
	// err := job.Run(context.Background())
	// assert.NoError(t, err)
}
```

### 集成测试

```go
// cmd/cron-worker/integration_test.go
package main

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
)

func TestCronManager_Start(t *testing.T) {
	// 创建测试用的 manager
	manager, err := cron.NewManager(logger, "UTC")
	assert.NoError(t, err)
	
	// 注册测试任务
	testJob := &TestJob{}
	err = manager.RegisterJob(testJob)
	assert.NoError(t, err)
	
	// 启动 manager
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	go func() {
		manager.Start(ctx)
	}()
	
	// 等待任务执行
	time.Sleep(2 * time.Second)
	cancel()
}
```

## 部署

### 构建

```bash
# 构建 cron-worker
go build -ldflags "-X main.Version=1.0.0 -X main.Name=cron-worker" \
  -o bin/cron-worker ./cmd/cron-worker
```

### Dockerfile

```dockerfile
# cmd/cron-worker/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o cron-worker ./cmd/cron-worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/cron-worker .
COPY --from=builder /app/configs/cron-worker.yaml ./configs/

CMD ["./cron-worker", "-conf", "./configs"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  cron-worker:
    build:
      context: .
      dockerfile: cmd/cron-worker/Dockerfile
    volumes:
      - ./configs/cron-worker.yaml:/app/configs/cron-worker.yaml
    restart: unless-stopped
    depends_on:
      - mysql
      - redis
```

## 替代方案：集成到现有应用

如果不想创建独立应用，也可以将 Cron 集成到现有应用中：

```go
// cmd/sre/main.go
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, 
            cronManager *cron.Manager, registrar kratosRegistry.Registrar) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	}
	
	// 启动 Cron 管理器
	if cronManager != nil {
		ctx := context.Background()
		go func() {
			if err := cronManager.Start(ctx); err != nil {
				log.NewHelper(logger).Errorf("cron manager error: %v", err)
			}
		}()
	}
	
	if registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}
	return kratos.New(opts...)
}
```

**注意**：这种方式不推荐用于生产环境，因为：
- 任务执行可能影响主服务性能
- 任务失败可能影响主服务稳定性
- 难以独立扩展和监控

## 常见问题

### Q: 如何确保任务只在一个实例执行？

A: 使用分布式锁（Redis 或 etcd），在任务执行前获取锁，执行完成后释放锁。

### Q: 任务执行时间过长怎么办？

A: 在任务内部设置超时控制，使用 `context.WithTimeout`。

### Q: 如何动态添加或删除任务？

A: Cron 管理器支持运行时添加任务，但删除任务需要重启应用。如果需要动态管理，可以考虑使用任务调度系统（如 Airflow、Temporal）。

### Q: 任务失败后如何重试？

A: 可以在任务内部实现重试逻辑，或使用支持重试的 cron 库。

### Q: 如何监控任务执行情况？

A: 集成 Prometheus 或 OpenTelemetry，记录任务执行时间、成功/失败次数等指标。

## 相关文档

- [多应用架构](../architecture/multi-app.md) - 了解 Kratos 多应用架构
- [依赖注入](../architecture/dependency-injection.md) - 了解 Wire 依赖注入
- [日志管理](../operations/logging.md) - 了解日志配置和管理

## 总结

1. **推荐使用独立应用**：将定时任务作为独立应用部署，实现资源隔离
2. **使用标准库**：使用 `github.com/robfig/cron/v3` 作为 cron 库
3. **遵循架构原则**：任务逻辑放在 `internal/biz/cron/` 目录
4. **完善错误处理**：任务内部要妥善处理错误，避免影响其他任务
5. **监控和日志**：记录任务执行情况，便于排查问题

