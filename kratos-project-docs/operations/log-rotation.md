# 日志轮转和清理策略

## 概述

项目已实现日志轮转和清理策略功能，使用 `lumberjack` 库自动管理日志文件，防止日志文件过大和磁盘空间不足。

## 功能特性

- ✅ **自动日志轮转**：当日志文件达到指定大小时自动轮转
- ✅ **备份文件管理**：自动保留指定数量的备份文件
- ✅ **自动清理**：自动删除超过保留天数的旧日志文件
- ✅ **压缩支持**：自动压缩旧日志文件，节省磁盘空间
- ✅ **配置化**：所有参数都可通过配置文件调整

## 配置说明

### 配置文件位置

**位置**：`configs/config.yaml`

```yaml
log:
  level: info
  format: json
  output_paths:
    - stdout
    - ./logs/app.log
  rotation:                # 日志轮转配置
    enable: true           # 是否启用日志轮转
    max_size: 100          # 每个日志文件最大大小（MB）
    max_backups: 10        # 保留的备份文件数量
    max_age: 30            # 保留天数
    compress: true         # 是否压缩旧日志文件
    local_time: true       # 使用本地时间而非 UTC
```

### 配置参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enable` | bool | false | 是否启用日志轮转 |
| `max_size` | int32 | 100 | 每个日志文件最大大小（MB），达到此大小会轮转 |
| `max_backups` | int32 | 10 | 保留的备份文件数量 |
| `max_age` | int32 | 30 | 保留天数，超过此天数的日志文件会被删除 |
| `compress` | bool | true | 是否压缩旧日志文件（.gz 格式） |
| `local_time` | bool | true | 使用本地时间而非 UTC 时间 |

## 工作原理

### 1. 日志轮转机制

当日志文件达到 `max_size` 时，会自动进行轮转：

```
logs/app.log          # 当前日志文件
logs/app.log.2025-12-16-001  # 轮转后的备份文件（如果 compress=false）
logs/app.log.2025-12-16-001.gz  # 压缩后的备份文件（如果 compress=true）
logs/app.log.2025-12-16-002.gz
...
```

**轮转规则**：
- 当前文件 `app.log` 达到 `max_size` 时，重命名为 `app.log.YYYY-MM-DD-NNN`
- 如果 `compress=true`，自动压缩为 `app.log.YYYY-MM-DD-NNN.gz`
- 创建新的 `app.log` 文件继续写入

### 2. 备份文件管理

**保留策略**：
- 最多保留 `max_backups` 个备份文件
- 超过数量的旧备份文件会被自动删除
- 删除顺序：最旧的文件优先删除

**示例**：
- `max_backups = 10`
- 当前有 15 个备份文件
- 自动删除最旧的 5 个文件，保留最新的 10 个

### 3. 自动清理机制

**清理规则**：
- 备份文件超过 `max_age` 天会被自动删除
- 清理在每次轮转时触发
- 使用文件的修改时间判断文件年龄

**示例**：
- `max_age = 30`
- 如果备份文件修改时间超过 30 天，会被自动删除

### 4. 压缩机制

**压缩规则**：
- 如果 `compress = true`，备份文件会自动压缩为 `.gz` 格式
- 压缩在轮转时进行
- 压缩可以显著减少磁盘占用（通常减少 70-90%）

**压缩效果**：
- 未压缩：100MB 日志文件
- 压缩后：约 10-30MB（取决于日志内容）

## 配置示例

### 开发环境配置

```yaml
log:
  rotation:
    enable: true
    max_size: 10          # 10MB，便于快速测试轮转
    max_backups: 5        # 保留 5 个备份
    max_age: 7            # 保留 7 天
    compress: false       # 开发环境不压缩，便于查看
    local_time: true
```

### 测试环境配置

```yaml
log:
  rotation:
    enable: true
    max_size: 50          # 50MB
    max_backups: 10       # 保留 10 个备份
    max_age: 14           # 保留 14 天
    compress: true        # 压缩旧文件
    local_time: true
```

### 生产环境配置

```yaml
log:
  rotation:
    enable: true
    max_size: 100         # 100MB
    max_backups: 20       # 保留 20 个备份（约 2GB 日志）
    max_age: 30           # 保留 30 天
    compress: true        # 压缩旧文件，节省空间
    local_time: true
```

### 高日志量环境配置

```yaml
log:
  rotation:
    enable: true
    max_size: 200         # 200MB，减少轮转频率
    max_backups: 30       # 保留 30 个备份
    max_age: 60           # 保留 60 天
    compress: true
    local_time: true
```

## 实现细节

### 1. 依赖库

使用 `gopkg.in/natefinch/lumberjack.v2` 库实现日志轮转：

```go
import "gopkg.in/natefinch/lumberjack.v2"
```

### 2. 代码实现

**文件**：`internal/logger/zap.go`

```go
// RotationConfig 日志轮转配置
type RotationConfig struct {
    Enable     bool  // 是否启用日志轮转
    MaxSize    int   // 每个日志文件最大大小（MB）
    MaxBackups int   // 保留的备份文件数量
    MaxAge     int   // 保留天数
    Compress   bool  // 是否压缩旧日志文件
    LocalTime  bool  // 使用本地时间而非 UTC
}

// getWriteSyncer 根据路径创建 WriteSyncer
func getWriteSyncer(paths []string, rotation *RotationConfig) []zapcore.WriteSyncer {
    // ...
    if rotation != nil && rotation.Enable {
        writer := &lumberjack.Logger{
            Filename:   path,
            MaxSize:    rotation.MaxSize,    // MB
            MaxBackups: rotation.MaxBackups,
            MaxAge:     rotation.MaxAge,     // 天
            Compress:   rotation.Compress,
            LocalTime:  rotation.LocalTime,
        }
        writers = append(writers, zapcore.AddSync(writer))
    } else {
        // 未启用轮转，使用普通文件
        file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        // ...
    }
}
```

### 3. 配置加载

**文件**：`internal/config/kratos.go`

```go
// 日志轮转配置
if v.IsSet("log.rotation") {
    rotation := &conf.Rotation{}
    if v.IsSet("log.rotation.enable") {
        rotation.Enable = v.GetBool("log.rotation.enable")
    }
    if v.IsSet("log.rotation.max_size") {
        rotation.MaxSize = int32(v.GetInt("log.rotation.max_size"))
    }
    // ... 其他字段 ...
    logConfig.Rotation = rotation
}
```

## 日志文件命名规则

### 未压缩模式

```
logs/app.log                    # 当前日志文件
logs/app.log.2025-12-16-001    # 备份文件 1
logs/app.log.2025-12-16-002    # 备份文件 2
logs/app.log.2025-12-16-003    # 备份文件 3
```

### 压缩模式

```
logs/app.log                    # 当前日志文件
logs/app.log.2025-12-16-001.gz # 压缩备份文件 1
logs/app.log.2025-12-16-002.gz # 压缩备份文件 2
logs/app.log.2025-12-16-003.gz # 压缩备份文件 3
```

**命名格式**：`{原文件名}.{日期}-{序号}[.gz]`

- 日期格式：`YYYY-MM-DD`
- 序号：从 001 开始递增
- 如果启用压缩，添加 `.gz` 后缀

## 磁盘空间估算

### 示例计算

假设配置：
- `max_size = 100` MB
- `max_backups = 10`
- `compress = true`

**未压缩情况**：
- 总大小 = 100MB × 10 = 1GB

**压缩后（假设压缩率 80%）**：
- 总大小 ≈ 100MB + (100MB × 10 × 0.2) = 300MB

### 不同配置的磁盘占用

| max_size | max_backups | 压缩 | 最大磁盘占用（估算） |
|----------|-------------|------|-------------------|
| 50MB | 10 | 否 | 500MB |
| 50MB | 10 | 是 | 150MB |
| 100MB | 10 | 否 | 1GB |
| 100MB | 10 | 是 | 300MB |
| 100MB | 20 | 是 | 600MB |
| 200MB | 30 | 是 | 1.8GB |

## 最佳实践

### 1. 根据日志量调整

**低日志量**（< 10MB/天）：
```yaml
max_size: 50
max_backups: 10
max_age: 30
```

**中日志量**（10-100MB/天）：
```yaml
max_size: 100
max_backups: 20
max_age: 30
```

**高日志量**（> 100MB/天）：
```yaml
max_size: 200
max_backups: 30
max_age: 60
```

### 2. 启用压缩

**建议**：生产环境始终启用压缩
- 节省 70-90% 磁盘空间
- 压缩和解压对性能影响很小
- 便于日志归档和传输

### 3. 合理设置保留天数

**建议**：
- **开发环境**：7-14 天（快速清理）
- **测试环境**：14-30 天（保留足够时间排查问题）
- **生产环境**：30-90 天（根据合规要求调整）

### 4. 监控磁盘空间

**建议监控指标**：
- 日志目录总大小
- 日志文件数量
- 磁盘使用率

**告警规则**：
- 日志目录 > 10GB：警告
- 日志目录 > 20GB：严重
- 磁盘使用率 > 80%：警告

## 故障排查

### 问题 1：日志文件没有轮转

**可能原因**：
1. 未启用日志轮转（`enable: false`）
2. 日志文件大小未达到 `max_size`
3. 配置未正确加载

**排查步骤**：
1. 检查配置文件中的 `rotation.enable` 是否为 `true`
2. 检查日志文件大小：`ls -lh logs/app.log`
3. 查看启动日志，确认配置是否加载

### 问题 2：备份文件过多

**可能原因**：
1. `max_backups` 设置过大
2. 清理机制未正常工作

**解决方案**：
1. 减少 `max_backups` 配置
2. 手动删除旧备份文件
3. 检查文件权限

### 问题 3：磁盘空间不足

**可能原因**：
1. `max_backups` 设置过大
2. `max_age` 设置过长
3. 未启用压缩

**解决方案**：
1. 减少 `max_backups` 和 `max_age`
2. 启用压缩（`compress: true`）
3. 手动清理旧日志文件

### 问题 4：压缩文件无法查看

**解决方案**：
```bash
# 查看压缩日志文件
zcat logs/app.log.2025-12-16-001.gz | less

# 或解压后查看
gunzip logs/app.log.2025-12-16-001.gz
cat logs/app.log.2025-12-16-001
```

## 手动管理日志

### 查看日志文件

```bash
# 查看当前日志
tail -f logs/app.log

# 查看压缩日志
zcat logs/app.log.2025-12-16-001.gz | tail -100

# 查看所有日志文件
ls -lh logs/
```

### 手动清理旧日志

```bash
# 删除超过 30 天的日志文件
find logs/ -name "app.log.*" -mtime +30 -delete

# 删除所有压缩日志
find logs/ -name "*.gz" -delete
```

### 压缩现有日志

```bash
# 压缩所有未压缩的备份文件
find logs/ -name "app.log.*" ! -name "*.gz" -exec gzip {} \;
```

## 监控和告警

### 1. 日志目录大小监控

```bash
# 查看日志目录大小
du -sh logs/

# 查看各个日志文件大小
du -h logs/*
```

### 2. Prometheus 监控（可选）

可以添加自定义指标监控日志目录大小：

```go
// 示例：监控日志目录大小
logDirSize := prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "log_directory_size_bytes",
        Help: "Size of log directory in bytes",
    },
    []string{"path"},
)
```

### 3. 告警规则

```yaml
# prometheus/alerts.yml
groups:
  - name: log_rotation
    rules:
      - alert: LogDirectoryTooLarge
        expr: log_directory_size_bytes > 10737418240  # 10GB
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Log directory size exceeds 10GB"
```

## 性能影响

### 轮转性能

- **轮转操作**：几乎无影响（文件重命名很快）
- **压缩操作**：轻微影响（在后台异步进行）
- **清理操作**：几乎无影响（在轮转时触发）

### 建议

- 对于高并发场景，建议：
  - 使用较大的 `max_size`（减少轮转频率）
  - 启用压缩（减少 I/O）
  - 使用 SSD 存储日志

## 参考资源

- [lumberjack 官方文档](https://github.com/natefinch/lumberjack)
- [Zap 日志库文档](https://github.com/uber-go/zap)
- [日志管理最佳实践](https://www.loggly.com/ultimate-guide/log-rotation/)

