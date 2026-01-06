# 测试策略

## 测试金字塔

```
        /\
       /  \      E2E 测试（少量）
      /____\
     /      \    集成测试（适量）
    /________\
   /          \  单元测试（大量）
  /____________\
```

### 单元测试
- **数量**：最多
- **速度**：最快
- **范围**：单个函数或方法
- **目标**：验证逻辑正确性

### 集成测试
- **数量**：适中
- **速度**：中等
- **范围**：多个组件协作
- **目标**：验证组件集成

### E2E 测试
- **数量**：最少
- **速度**：最慢
- **范围**：完整流程
- **目标**：验证用户场景

## Go 测试实践

### 单元测试示例

```go
func TestGreeterUsecase_CreateUser(t *testing.T) {
    // 准备测试数据
    repo := &MockGreeterRepo{}
    uc := biz.NewGreeterUsecase(repo)
    
    // 执行测试
    user, err := uc.CreateUser(context.Background(), &biz.User{
        Name:  "Alice",
        Email: "alice@example.com",
    })
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "Alice", user.Name)
}
```

### 表驱动测试

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Mock 使用

使用接口进行 Mock：

```go
type GreeterRepo interface {
    Save(ctx context.Context, user *User) error
}

type MockGreeterRepo struct {
    SaveFunc func(ctx context.Context, user *User) error
}

func (m *MockGreeterRepo) Save(ctx context.Context, user *User) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(ctx, user)
    }
    return nil
}
```

## 测试覆盖率

### 查看覆盖率
```bash
# 运行测试并生成覆盖率报告
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 覆盖率目标
- **单元测试**：核心业务逻辑 > 80%
- **集成测试**：关键流程 > 60%
- **整体覆盖率**：> 70%

## 测试最佳实践

1. **测试命名**：使用 `Test函数名_场景` 格式
2. **测试隔离**：每个测试独立，不依赖其他测试
3. **快速反馈**：保持测试快速执行
4. **可读性**：测试代码要清晰易读
5. **维护性**：测试代码也要保持整洁

## 测试工具

- **testing**：Go 标准测试包
- **testify**：断言和 Mock 工具
- **gomock**：接口 Mock 生成工具
- **httptest**：HTTP 测试工具

