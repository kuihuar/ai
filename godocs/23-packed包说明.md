# packed 包说明

## 什么是 packed 包？

`internal/packed/packed.go` 是 GoFrame 框架中的**资源打包**功能，用于将静态资源（配置文件、模板文件等）打包到二进制文件中。

## 功能说明

### 1. 资源嵌入

GoFrame 支持将资源文件打包到编译后的二进制文件中，这样：
- ✅ 不需要单独部署配置文件
- ✅ 可以生成单个可执行文件
- ✅ 便于分发和部署

### 2. 工作原理

```go
// main.go
import (
    _ "hz/internal/packed"  // 空导入，执行 init 函数
)
```

**执行流程**:
1. 导入 `packed` 包时，会执行 `init()` 函数
2. `init()` 函数会加载打包的资源文件
3. 资源文件被嵌入到二进制文件中（通过 `go:embed` 或 `gres`）

### 3. 生成方式

使用 GoFrame 的 `gf` 工具链生成：

```bash
# 打包资源文件
gf pack public,template,config packed

# 或者使用 Makefile
make pack
```

**打包后的效果**:
- 资源文件被编译到 `internal/packed/packed.go` 中
- 通过 `gres` 包访问打包的资源

### 4. 使用示例

```go
// 访问打包的资源
import "github.com/gogf/gf/v2/os/gres"

// 读取打包的配置文件
configContent := gres.GetContent("config/config.yaml")

// 读取打包的模板文件
templateContent := gres.GetContent("template/index.html")

// 检查资源是否存在
if gres.Contains("config/config.yaml") {
    // 资源存在
}
```

## 当前项目中的使用

### 当前状态

```go
// internal/packed/packed.go
package packed
```

**说明**:
- 当前文件是空的，表示**还没有打包资源**
- 如果需要打包资源，需要：
  1. 使用 `gf pack` 命令打包资源
  2. 或者手动实现资源嵌入

### 是否需要打包？

**不需要打包的情况**:
- ✅ 配置文件放在外部（如 `manifest/config/`）
- ✅ 模板文件放在外部
- ✅ 需要动态修改配置文件

**需要打包的情况**:
- ✅ 生成单个可执行文件
- ✅ 配置文件不需要修改
- ✅ 简化部署（不需要额外的文件）

## 实际应用场景

### 场景1: 单文件部署

```bash
# 打包资源
gf pack config,template packed

# 编译
go build -o app

# 部署（只需要一个文件）
./app
```

### 场景2: 配置热更新

```go
// 不使用打包，配置文件放在外部
// manifest/config/config.yaml
// 可以动态修改，应用自动热更新
```

## 与 Kratos 的对比

### GoFrame (packed)
```go
// 资源打包到二进制文件
import _ "hz/internal/packed"
gres.GetContent("config.yaml")
```

### Kratos
```go
// 配置文件放在外部，通过 Source 加载
config.New(
    config.WithSource(file.NewSource("config.yaml")),
)
```

**区别**:
- **GoFrame**: 支持资源打包，单文件部署
- **Kratos**: 配置文件外部化，支持动态更新

## 总结

`packed` 包是 GoFrame 的资源打包功能，用于：
1. **资源嵌入**: 将静态资源打包到二进制文件
2. **简化部署**: 生成单个可执行文件
3. **可选功能**: 根据需求决定是否使用

**当前项目**: 文件为空，表示未使用资源打包功能，配置文件放在 `manifest/config/` 目录中，这是更灵活的方式。

