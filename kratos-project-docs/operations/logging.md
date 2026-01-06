# 日志规范

## Kratos 日志系统

Kratos 使用结构化日志，支持多种日志后端。

### 日志使用

```go
import "github.com/go-kratos/kratos/v2/log"

// 创建 logger
logger := log.NewStdLogger(os.Stdout)

// 使用 logger
log.NewHelper(logger).Infow("key", "value", "message", "user created")
log.NewHelper(logger).Errorw("key", "value", "err", err)
```

### 日志级别

- `DEBUG`：调试信息，开发时使用
- `INFO`：一般信息，记录重要操作
- `WARN`：警告信息，潜在问题
- `ERROR`：错误信息，需要关注
- `FATAL`：致命错误，程序无法继续

## 日志规范

### 1. 结构化日志
使用键值对形式，便于日志分析：

```go
// ✅ 正确：结构化日志
logger.Infow("user_created", 
    "user_id", userID,
    "username", username,
    "timestamp", time.Now(),
)

// ❌ 错误：非结构化日志
logger.Info(fmt.Sprintf("User %s created with ID %d", username, userID))
```

### 2. 日志上下文
在请求处理中传递日志上下文：

```go
func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    log := log.NewHelper(log.FromContext(ctx))
    log.Infow("get_user", "user_id", req.Id)
    // ...
}
```

### 3. 敏感信息
不要在日志中记录敏感信息：
- 密码
- 令牌
- 信用卡号
- 个人隐私信息

### 4. 错误日志
记录错误时包含足够的上下文：

```go
if err != nil {
    log.Errorw("failed_to_save_user",
        "user_id", userID,
        "error", err,
        "stack", fmt.Sprintf("%+v", err),
    )
    return err
}
```

## 日志格式

### 推荐格式
```
时间戳 [级别] 消息 key1=value1 key2=value2
```

### 示例
```
2024-01-01T10:00:00Z [INFO] user_created user_id=123 username=alice
2024-01-01T10:00:01Z [ERROR] failed_to_save_user user_id=123 error="connection timeout"
```

## 最佳实践

1. **日志级别合理**：根据重要性选择合适的日志级别
2. **日志量控制**：避免过多日志影响性能
3. **日志聚合**：使用日志聚合系统（如 ELK）统一管理
4. **日志监控**：设置日志告警，及时发现问题

