适配器模式（Adapter Pattern） 的核心思想是将一个接口转换为另一个客户端期望的接口，使原本不兼容的类（或类型）能够协同工作

场景描述
假设系统中存在两个不兼容的日志组件：

旧版文件日志：写入本地文件，方法签名为 LogToFile(msg string)。

新版云日志：上传日志到云服务，接口定义为 CloudLogger，需要实现 UploadLog(msg string) error。

目标：将旧版文件日志的 LogToFile 方法适配到 CloudLogger 接口。

1. 定义目标接口（新版云日志接口）

```go
/ CloudLogger 是新的日志接口标准
type CloudLogger interface {
    UploadLog(msg string) error
}
```
2. 现有旧版文件日志（需适配的类）
```go
// FileLogger 旧版文件日志实现
type FileLogger struct{}

// LogToFile 旧版方法，不符合 CloudLogger 接口
func (f *FileLogger) LogToFile(msg string) {
    fmt.Printf("[旧版日志] 写入文件: %s\n", msg)
    // 实际文件操作逻辑...
}

```

3. 创建适配器（关键步骤）

```go
// FileLoggerAdapter 是适配器，将 FileLogger 适配到 CloudLogger 接口
type FileLoggerAdapter struct {
    fileLogger *FileLogger
}

// 实现 CloudLogger 接口的 UploadLog 方法
func (a *FileLoggerAdapter) UploadLog(msg string) error {
    // 调用旧版方法，并添加适配逻辑（如错误处理转换）
    a.fileLogger.LogToFile(msg)
    return nil // 假设旧版方法无错误返回，适配时可忽略或转换
}

// 创建适配器的构造函数
func NewFileLoggerAdapter(logger *FileLogger) CloudLogger {
    return &FileLoggerAdapter{fileLogger: logger}
}

```
4. 客户端使用（无缝对接新版接口）

```go
func main() {
    // 旧版日志实例
    oldFileLogger := &FileLogger{}

    // 创建适配器，将旧版日志转换为 CloudLogger 接口
    cloudLogger := NewFileLoggerAdapter(oldFileLogger)

    // 客户端代码统一调用新版接口
    err := cloudLogger.UploadLog("用户登录成功")
    if err != nil {
        fmt.Println("日志上传失败:", err)
    }
}

```
适配器模式 vs 装饰器模式
模式	目的	特点
适配器	转换接口，解决兼容性问题	接口不同，功能不变
装饰器	动态添加功能，不改变接口	接口相同，增强功能
--------------------------------------------------------------------------

1. 新旧接口适配

```go
package main

import "fmt"

// 旧的日志记录器接口
type OldLogger interface {
    LogOld(message string)
}

// 旧的日志记录器实现
type OldLoggerImpl struct{}

func (ol *OldLoggerImpl) LogOld(message string) {
    fmt.Println("Old Logger: ", message)
}

// 新的日志记录器接口
type NewLogger interface {
    LogNew(message string)
}

// 新的日志记录器实现
type NewLoggerImpl struct{}

func (nl *NewLoggerImpl) LogNew(message string) {
    fmt.Println("New Logger: ", message)
}

// 适配器，将新日志记录器适配为旧日志记录器接口
type LoggerAdapter struct {
    newLogger NewLogger
}

func (la *LoggerAdapter) LogOld(message string) {
    la.newLogger.LogNew(message)
}

func main() {
    // 创建新日志记录器实例
    newLogger := &NewLoggerImpl{}
    // 创建适配器实例
    adapter := &LoggerAdapter{newLogger: newLogger}

    // 使用适配器，调用旧接口方法
    adapter.LogOld("This is a log message")
}

```

2. 第三方库接口适配
假设你使用一个第三方库提供的 Square 结构体和计算面积的方法，同时你有自己的 Shape 接口，你想让 Square 结构体适配到 Shape 接口上。

package main

import "fmt"

// 第三方库的 Square 结构体
type Square struct {
    Side float64
}

func (s *Square) CalculateSquareArea() float64 {
    return s.Side * s.Side
}

// 自己定义的 Shape 接口
type Shape interface {
    Area() float64
}

// 适配器，将 Square 适配到 Shape 接口
type SquareAdapter struct {
    square *Square
}

func (sa *SquareAdapter) Area() float64 {
    return sa.square.CalculateSquareArea()
}

func main() {
    // 创建 Square 实例
    square := &Square{Side: 5}
    // 创建适配器实例
    adapter := &SquareAdapter{square: square}

    // 通过适配器调用 Shape 接口的 Area 方法
    fmt.Println("Area of the square:", adapter.Area())
}
代码解释：

Square 是第三方库提供的结构体，有自己的计算面积的方法 CalculateSquareArea。
Shape 是你自己定义的接口，有 Area 方法。
SquareAdapter 是适配器，它持有一个 Square 指针，并实现了 Shape 接口的 Area 方法，在该方法中调用 Square 的 CalculateSquareArea 方法。
在 main 函数中，创建 Square 和适配器实例，然后通过适配器调用 Shape 接口的 Area 方法，实现了第三方库接口到自定义接口的适配。
