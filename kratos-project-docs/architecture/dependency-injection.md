# 依赖注入实践

## Wire 依赖注入

Kratos 使用 Google Wire 进行编译时依赖注入，相比运行时注入，具有以下优势：

- **编译时检查**：依赖关系在编译时确定，避免运行时错误
- **性能更好**：无需反射，直接生成代码
- **类型安全**：编译期保证类型正确

## Wire 使用模式

### Provider 函数

定义提供依赖的函数：

```go
// Provider 函数命名规范：New + 类型名
func NewGreeterUsecase(repo biz.GreeterRepo) *biz.GreeterUsecase {
    return biz.NewGreeterUsecase(repo)
}

func NewGreeterRepo(data *Data) biz.GreeterRepo {
    return &greeterRepo{data: data}
}
```

### Wire Set

将相关的 Provider 组织成 Set：

```go
// ProviderSet 命名规范：类型名 + ProviderSet
var ProviderSet = wire.NewSet(NewGreeterRepo)
```

### Wire 初始化

在 `wire.go` 中定义初始化函数：

```go
//go:build wireinject
// +build wireinject

func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        newApp,
    ))
}
```

## 最佳实践

### 1. Provider 函数设计
- 函数名清晰表达创建的对象
- 参数只包含必要的依赖
- 返回值和错误处理明确

### 2. ProviderSet 组织
- 按模块组织 ProviderSet
- 避免循环依赖
- 保持 Set 的粒度适中

### 3. 接口注入
- Biz 层定义接口
- Data 层实现接口
- 通过接口注入，实现解耦

### 4. 生命周期管理
- 使用 cleanup 函数管理资源
- 确保资源正确释放
- 处理初始化失败的情况

## 常见问题

### 1. 如何处理可选依赖？
使用 `wire.Value` 或 `wire.InterfaceValue` 提供默认值。

### 2. 如何注入配置？
将配置作为参数传入 Provider 函数。

### 3. 如何处理循环依赖？
重新设计依赖关系，通常通过引入接口或中间层解决。

