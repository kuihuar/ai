# 雪花算法编号生成器

## 概述

实现了雪花算法（Snowflake）作为唯一编号生成器的备选方案。雪花算法可以在不依赖数据库的情况下生成全局唯一的64位ID。

## 雪花算法结构

64位ID的组成：
- **1位符号位**（固定为0）
- **41位时间戳**（毫秒，从2020-01-01开始，可用约69年）
- **10位机器ID**（5位数据中心ID + 5位机器ID，最多支持1024个节点）
- **12位序列号**（同一毫秒内最多4096个ID）

## 配置方式

### 1. 更新配置文件

在 `configs/config.yaml` 中添加：

```yaml
data:
  number_generator:
    type: "snowflake"      # "db" 或 "snowflake"
    data_center_id: 1      # 数据中心ID (0-31)
    machine_id: 1          # 机器ID (0-31)
```

### 2. 使用数据库生成器（默认）

```yaml
data:
  number_generator:
    type: "db"  # 或不配置，默认使用数据库生成器
```

## 代码实现

### 1. 雪花算法生成器 (`internal/pkg/number/snowflake.go`)

```go
type SnowflakeGenerator struct {
    mu            sync.Mutex
    epoch         int64 // 起始时间戳（毫秒）
    dataCenterID  int64 // 数据中心ID (0-31)
    machineID     int64 // 机器ID (0-31)
    sequence      int64 // 序列号 (0-4095)
    lastTimestamp int64 // 上次生成ID的时间戳（毫秒）
}

func NewSnowflakeGenerator(dataCenterID, machineID int64, logger log.Logger) Generator
```

### 2. 生成器工厂 (`internal/data/number_generator.go`)

```go
// NewNumberGenerator 默认使用数据库生成器
func NewNumberGenerator(repo NumberRepo, logger log.Logger) number.Generator

// NewNumberGeneratorWithConfig 根据配置选择生成器
func NewNumberGeneratorWithConfig(repo NumberRepo, c *conf.Data, logger log.Logger) number.Generator
```

## 使用方式

### 方式1：通过 Wire 注入（推荐）

修改 `cmd/sre/wire.go`：

```go
// 将 NewNumberGenerator 替换为 NewNumberGeneratorWithConfig
data.NewNumberGeneratorWithConfig,
```

然后在 `wireApp` 函数中传入 `*conf.Data`。

### 方式2：手动创建

```go
// 使用数据库生成器
generator := data.NewNumberGenerator(repo, logger)

// 使用雪花算法生成器
generator := number.NewSnowflakeGenerator(1, 1, logger)
```

## 编号格式

### 数据库生成器
- 格式: `{prefix}{Unix时间戳秒}{6位序列号}`
- 示例: `ORD1765946794000001`

### 雪花算法生成器
- 格式: `{prefix}{64位雪花算法ID}`
- 示例: `ORD1234567890123456789`

## 优势对比

### 数据库生成器
- ✅ 需要数据库，但保证全局唯一性
- ✅ 支持不同业务前缀的独立序列号
- ✅ 可以查询和追踪序列号使用情况
- ❌ 依赖数据库性能

### 雪花算法生成器
- ✅ 不需要数据库，性能高
- ✅ 全局唯一，支持分布式环境
- ✅ 包含时间信息，可以反推生成时间
- ❌ 需要配置数据中心ID和机器ID
- ❌ 时钟回拨会导致ID重复（已处理）

## 时钟回拨处理

雪花算法生成器已实现时钟回拨检测：

```go
if now < g.lastTimestamp {
    return 0, fmt.Errorf("clock moved backwards, refusing to generate id")
}
```

如果检测到时钟回拨，会返回错误，避免生成重复ID。

## 性能特点

- **数据库生成器**: 每次生成需要数据库事务，性能取决于数据库
- **雪花算法生成器**: 纯内存操作，性能极高（每秒可生成数百万个ID）

## 注意事项

1. **数据中心ID和机器ID**: 必须在集群中唯一，否则可能生成重复ID
2. **时钟同步**: 建议使用NTP同步服务器时间，避免时钟回拨
3. **序列号溢出**: 同一毫秒内超过4096个ID时，会等待下一毫秒
4. **起始时间**: 默认从2020-01-01开始，可用到2089年

## 迁移建议

1. **新项目**: 可以直接使用雪花算法生成器
2. **现有项目**: 建议继续使用数据库生成器，保证兼容性
3. **混合使用**: 可以为不同业务使用不同的生成器

