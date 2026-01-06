# 错误处理

## Kratos 错误处理

Kratos 提供了统一的错误处理机制，基于 Protobuf 定义错误原因。

### 错误定义

在 `api/helloworld/v1/error_reason.proto` 中定义错误：

```protobuf
enum ErrorReason {
  USER_NOT_FOUND = 0;
  INVALID_PARAMS = 1;
  INTERNAL_ERROR = 2;
}
```

### 错误使用

```go
import "sre/api/helloworld/v1"

// 返回错误
return nil, errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")

// 检查错误
if errors.Is(err, v1.ErrorReason_USER_NOT_FOUND) {
    // 处理用户不存在的情况
}
```

## 错误处理原则

### 1. 及时处理错误
不要忽略错误，总是检查并处理。

```go
// ❌ 错误：忽略错误
result, _ := someFunction()

// ✅ 正确：处理错误
result, err := someFunction()
if err != nil {
    return nil, err
}
```

### 2. 添加上下文
使用 `fmt.Errorf` 和 `%w` 包装错误，添加上下文信息。

```go
if err != nil {
    return nil, fmt.Errorf("failed to save user %d: %w", userID, err)
}
```

### 3. 错误分类
区分可恢复错误和不可恢复错误。

- **可恢复错误**：网络超时、临时性错误，可以重试
- **不可恢复错误**：参数错误、权限错误，不应该重试

### 4. 错误传播
在适当的层级处理错误，不要过度包装。

```go
// 在 Data 层
func (r *repo) Save(ctx context.Context, user *User) error {
    if err := r.db.Save(user).Error; err != nil {
        return fmt.Errorf("database save failed: %w", err)
    }
    return nil
}

// 在 Biz 层
func (uc *Usecase) CreateUser(ctx context.Context, user *User) error {
    if err := uc.repo.Save(ctx, user); err != nil {
        return errors.Internal("CREATE_USER_FAILED", "failed to create user")
    }
    return nil
}
```

## 错误类型

### 业务错误
使用 Kratos 错误类型表示业务错误：

```go
errors.NotFound(reason, message)    // 404
errors.BadRequest(reason, message)  // 400
errors.Unauthorized(reason, message) // 401
errors.Forbidden(reason, message)   // 403
errors.Internal(reason, message)    // 500
```

### 系统错误
系统级错误（如数据库连接失败）应该记录日志并返回通用错误。

## 最佳实践

1. **统一错误格式**：使用框架提供的错误类型
2. **错误日志**：记录详细的错误信息用于调试
3. **用户友好**：返回给用户的错误信息要友好
4. **错误监控**：集成错误监控系统，及时发现问题

