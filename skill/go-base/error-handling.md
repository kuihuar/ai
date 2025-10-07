# Go 错误处理详解

## 📚 目录

- [错误处理基础](#错误处理基础)
- [自定义错误类型](#自定义错误类型)
- [错误包装和展开](#错误包装和展开)
- [错误处理模式](#错误处理模式)
- [panic 和 recover](#panic-和-recover)
- [最佳实践](#最佳实践)
- [常见错误处理场景](#常见错误处理场景)

## 错误处理基础

### 基本错误处理

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

// 1. 基本错误返回
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 2. 使用 fmt.Errorf 创建错误
func validateAge(age int) error {
    if age < 0 {
        return fmt.Errorf("age cannot be negative: %d", age)
    }
    if age > 150 {
        return fmt.Errorf("age cannot be greater than 150: %d", age)
    }
    return nil
}

// 3. 错误检查和处理
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file %s: %w", filename, err)
    }
    defer file.Close()
    
    // 处理文件...
    return nil
}

func main() {
    // 基本错误处理
    result, err := divide(10, 2)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %d\n", result)
    }
    
    // 错误处理示例
    result, err = divide(10, 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // 验证错误
    if err := validateAge(25); err != nil {
        fmt.Printf("Validation error: %v\n", err)
    }
    
    if err := validateAge(-5); err != nil {
        fmt.Printf("Validation error: %v\n", err)
    }
}
```

### 错误类型检查

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func (e ValidationError) Is(target error) bool {
    var v ValidationError
    if errors.As(target, &v) {
        return e.Field == v.Field
    }
    return false
}

func main() {
    // 1. 基本错误检查
    err := os.Open("nonexistent.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // 2. 错误类型断言
    err = os.Open("nonexistent.txt")
    if err != nil {
        if pathErr, ok := err.(*os.PathError); ok {
            fmt.Printf("Path error: %s, op: %s, path: %s\n", 
                      pathErr.Err, pathErr.Op, pathErr.Path)
        }
    }
    
    // 3. 使用 errors.Is 检查特定错误
    err = os.Open("nonexistent.txt")
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("File does not exist")
        }
    }
    
    // 4. 使用 errors.As 类型断言
    err = ValidationError{Field: "email", Message: "invalid format"}
    var validationErr ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation error: field=%s, message=%s\n", 
                  validationErr.Field, validationErr.Message)
    }
    
    // 5. 自定义错误的 Is 方法
    err1 := ValidationError{Field: "email", Message: "invalid"}
    err2 := ValidationError{Field: "email", Message: "required"}
    
    if errors.Is(err1, err2) {
        fmt.Println("Same field validation error")
    } else {
        fmt.Println("Different validation errors")
    }
}
```

## 自定义错误类型

### 结构化错误

```go
package main

import (
    "fmt"
    "time"
)

// 1. 基本自定义错误
type APIError struct {
    Code    int
    Message string
    Time    time.Time
}

func (e APIError) Error() string {
    return fmt.Sprintf("API error %d: %s at %v", e.Code, e.Message, e.Time)
}

// 2. 带上下文的错误
type DatabaseError struct {
    Operation string
    Table     string
    Err       error
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("database error during %s on table %s: %v", 
                      e.Operation, e.Table, e.Err)
}

func (e DatabaseError) Unwrap() error {
    return e.Err
}

// 3. 错误链
type BusinessError struct {
    Code    string
    Message string
    Cause   error
}

func (e BusinessError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e BusinessError) Unwrap() error {
    return e.Cause
}

func main() {
    // 使用自定义错误
    err := APIError{
        Code:    404,
        Message: "Resource not found",
        Time:    time.Now(),
    }
    fmt.Printf("Error: %v\n", err)
    
    // 使用带上下文的错误
    dbErr := DatabaseError{
        Operation: "INSERT",
        Table:     "users",
        Err:       fmt.Errorf("duplicate key"),
    }
    fmt.Printf("Database error: %v\n", dbErr)
    
    // 使用错误链
    businessErr := BusinessError{
        Code:    "USER_NOT_FOUND",
        Message: "User does not exist",
        Cause:   dbErr,
    }
    fmt.Printf("Business error: %v\n", businessErr)
}
```

### 错误分类

```go
package main

import (
    "errors"
    "fmt"
)

// 错误类型定义
type ErrorType int

const (
    ErrorTypeValidation ErrorType = iota
    ErrorTypeNotFound
    ErrorTypePermission
    ErrorTypeInternal
)

// 分类错误
type CategorizedError struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (e CategorizedError) Error() string {
    typeName := []string{"validation", "not_found", "permission", "internal"}[e.Type]
    if e.Cause != nil {
        return fmt.Sprintf("%s error: %s (caused by: %v)", typeName, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s error: %s", typeName, e.Message)
}

func (e CategorizedError) Unwrap() error {
    return e.Cause
}

// 错误类型检查
func (e CategorizedError) Is(target error) bool {
    var t CategorizedError
    if errors.As(target, &t) {
        return e.Type == t.Type
    }
    return false
}

// 错误创建函数
func NewValidationError(message string, cause error) error {
    return CategorizedError{
        Type:    ErrorTypeValidation,
        Message: message,
        Cause:   cause,
    }
}

func NewNotFoundError(message string, cause error) error {
    return CategorizedError{
        Type:    ErrorTypeNotFound,
        Message: message,
        Cause:   cause,
    }
}

func main() {
    // 创建不同类型的错误
    validationErr := NewValidationError("email is required", nil)
    notFoundErr := NewNotFoundError("user not found", nil)
    
    fmt.Printf("Validation error: %v\n", validationErr)
    fmt.Printf("Not found error: %v\n", notFoundErr)
    
    // 错误类型检查
    if errors.Is(validationErr, CategorizedError{Type: ErrorTypeValidation}) {
        fmt.Println("This is a validation error")
    }
    
    if errors.Is(notFoundErr, CategorizedError{Type: ErrorTypeNotFound}) {
        fmt.Println("This is a not found error")
    }
}
```

## 错误包装和展开

### 错误包装

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

// 1. 使用 fmt.Errorf 和 %w 包装错误
func readConfig(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open config file %s: %w", filename, err)
    }
    defer file.Close()
    
    // 模拟读取配置
    return fmt.Errorf("invalid config format: %w", errors.New("missing required field"))
}

// 2. 使用 errors.Wrap 包装错误
func processData(data string) error {
    if data == "" {
        return fmt.Errorf("processing data: %w", errors.New("empty data"))
    }
    
    // 模拟处理
    return fmt.Errorf("processing data: %w", errors.New("validation failed"))
}

// 3. 多层错误包装
func handleRequest() error {
    err := readConfig("config.json")
    if err != nil {
        return fmt.Errorf("handling request: %w", err)
    }
    
    err = processData("test data")
    if err != nil {
        return fmt.Errorf("handling request: %w", err)
    }
    
    return nil
}

func main() {
    // 测试错误包装
    err := handleRequest()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        
        // 展开错误链
        fmt.Println("\nError chain:")
        for {
            fmt.Printf("  %v\n", err)
            err = errors.Unwrap(err)
            if err == nil {
                break
            }
        }
    }
}
```

### 错误展开和检查

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func (e ValidationError) Is(target error) bool {
    var v ValidationError
    if errors.As(target, &v) {
        return e.Field == v.Field
    }
    return false
}

func main() {
    // 创建错误链
    err := fmt.Errorf("api call failed: %w", 
        fmt.Errorf("database error: %w", 
            ValidationError{Field: "email", Message: "invalid format"}))
    
    fmt.Printf("Original error: %v\n", err)
    
    // 1. 使用 errors.Is 检查错误链
    var validationErr ValidationError
    if errors.Is(err, validationErr) {
        fmt.Println("Found validation error in chain")
    }
    
    // 2. 使用 errors.As 提取特定错误类型
    var extracted ValidationError
    if errors.As(err, &extracted) {
        fmt.Printf("Extracted validation error: field=%s, message=%s\n", 
                  extracted.Field, extracted.Message)
    }
    
    // 3. 手动展开错误链
    fmt.Println("\nManual error unwrapping:")
    current := err
    for i := 0; current != nil && i < 10; i++ {
        fmt.Printf("Level %d: %v\n", i, current)
        current = errors.Unwrap(current)
    }
    
    // 4. 检查特定错误类型
    err = fmt.Errorf("file operation failed: %w", os.ErrNotExist)
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File does not exist")
    }
    
    // 5. 检查多个错误类型
    err = fmt.Errorf("io error: %w", io.EOF)
    if errors.Is(err, io.EOF) {
        fmt.Println("End of file reached")
    }
}
```

## 错误处理模式

### 错误处理策略

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "os"
)

// 1. 快速失败模式
func quickFail(data string) error {
    if data == "" {
        return errors.New("data cannot be empty")
    }
    
    if len(data) < 3 {
        return errors.New("data too short")
    }
    
    return nil
}

// 2. 重试模式
func retryOperation(operation func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := operation(); err != nil {
            lastErr = err
            fmt.Printf("Attempt %d failed: %v\n", i+1, err)
            continue
        }
        return nil
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, lastErr)
}

// 3. 优雅降级模式
func gracefulDegradation(data string) (string, error) {
    // 尝试主要处理
    result, err := processData(data)
    if err == nil {
        return result, nil
    }
    
    // 降级到备用处理
    fmt.Printf("Primary processing failed: %v, using fallback\n", err)
    return fallbackProcess(data), nil
}

func processData(data string) (string, error) {
    if data == "error" {
        return "", errors.New("processing failed")
    }
    return "processed: " + data, nil
}

func fallbackProcess(data string) string {
    return "fallback: " + data
}

// 4. 错误聚合模式
type MultiError struct {
    Errors []error
}

func (m MultiError) Error() string {
    return fmt.Sprintf("multiple errors occurred: %v", m.Errors)
}

func (m MultiError) Add(err error) {
    if err != nil {
        m.Errors = append(m.Errors, err)
    }
}

func (m MultiError) HasErrors() bool {
    return len(m.Errors) > 0
}

func main() {
    // 快速失败
    if err := quickFail(""); err != nil {
        fmt.Printf("Quick fail: %v\n", err)
    }
    
    // 重试模式
    attemptCount := 0
    err := retryOperation(func() error {
        attemptCount++
        if attemptCount < 3 {
            return errors.New("temporary failure")
        }
        return nil
    }, 5)
    
    if err != nil {
        fmt.Printf("Retry failed: %v\n", err)
    } else {
        fmt.Println("Retry succeeded")
    }
    
    // 优雅降级
    result, err := gracefulDegradation("test")
    if err != nil {
        fmt.Printf("Graceful degradation failed: %v\n", err)
    } else {
        fmt.Printf("Result: %s\n", result)
    }
    
    // 错误聚合
    var multiErr MultiError
    multiErr.Add(errors.New("error 1"))
    multiErr.Add(errors.New("error 2"))
    multiErr.Add(nil) // 不会添加
    
    if multiErr.HasErrors() {
        fmt.Printf("Multiple errors: %v\n", multiErr)
    }
}
```

### 错误处理中间件

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "time"
)

// 错误处理中间件
type ErrorHandler struct {
    logger *log.Logger
}

func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{
        logger: log.New(os.Stdout, "ERROR: ", log.LstdFlags),
    }
}

// 1. 记录错误
func (h *ErrorHandler) LogError(err error, context string) {
    h.logger.Printf("%s: %v", context, err)
}

// 2. 错误恢复
func (h *ErrorHandler) Recover(operation func() error) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
            h.LogError(err, "panic recovery")
        }
    }()
    
    return operation()
}

// 3. 错误重试
func (h *ErrorHandler) Retry(operation func() error, maxRetries int, delay time.Duration) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := operation(); err != nil {
            lastErr = err
            h.LogError(err, fmt.Sprintf("attempt %d", i+1))
            
            if i < maxRetries-1 {
                time.Sleep(delay)
            }
            continue
        }
        return nil
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, lastErr)
}

// 4. 错误转换
func (h *ErrorHandler) TransformError(err error, transform func(error) error) error {
    if err == nil {
        return nil
    }
    
    transformed := transform(err)
    if transformed != err {
        h.LogError(transformed, "error transformed")
    }
    
    return transformed
}

func main() {
    handler := NewErrorHandler()
    
    // 测试错误恢复
    err := handler.Recover(func() error {
        panic("something went wrong")
    })
    if err != nil {
        fmt.Printf("Recovered error: %v\n", err)
    }
    
    // 测试错误重试
    attemptCount := 0
    err = handler.Retry(func() error {
        attemptCount++
        if attemptCount < 3 {
            return errors.New("temporary failure")
        }
        return nil
    }, 5, 100*time.Millisecond)
    
    if err != nil {
        fmt.Printf("Retry failed: %v\n", err)
    } else {
        fmt.Println("Retry succeeded")
    }
    
    // 测试错误转换
    originalErr := errors.New("original error")
    transformedErr := handler.TransformError(originalErr, func(err error) error {
        return fmt.Errorf("transformed: %w", err)
    })
    fmt.Printf("Transformed error: %v\n", transformedErr)
}
```

## panic 和 recover

### panic 使用场景

```go
package main

import (
    "fmt"
    "log"
)

// 1. 不可恢复的错误
func mustNotBeZero(n int) {
    if n == 0 {
        panic("number cannot be zero")
    }
}

// 2. 程序逻辑错误
func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

// 3. 初始化失败
func initialize() {
    // 模拟初始化失败
    panic("initialization failed")
}

func main() {
    // 基本 panic 使用
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
        }
    }()
    
    mustNotBeZero(5)
    mustNotBeZero(0) // 这会触发 panic
    
    // 不会执行到这里
    fmt.Println("This won't be printed")
}
```

### recover 使用模式

```go
package main

import (
    "fmt"
    "log"
    "os"
)

// 1. 基本 recover 使用
func safeOperation() (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("operation panicked: %v", r)
        }
    }()
    
    // 可能 panic 的操作
    result = divide(10, 0)
    return result, nil
}

// 2. 带日志的 recover
func safeOperationWithLog() (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered: %v", r)
            err = fmt.Errorf("operation panicked: %v", r)
        }
    }()
    
    result = divide(10, 0)
    return result, nil
}

// 3. 恢复后继续执行
func resilientOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v, continuing...\n", r)
        }
    }()
    
    // 可能 panic 的操作
    panic("something went wrong")
    
    // 这行不会执行
    fmt.Println("This won't be printed")
}

// 4. 多层 recover
func outerFunction() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Outer recover: %v\n", r)
        }
    }()
    
    innerFunction()
}

func innerFunction() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Inner recover: %v\n", r)
            panic("re-panic") // 重新 panic
        }
    }()
    
    panic("inner panic")
}

func main() {
    // 测试基本 recover
    result, err := safeOperation()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %d\n", result)
    }
    
    // 测试带日志的 recover
    result, err = safeOperationWithLog()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // 测试恢复后继续执行
    resilientOperation()
    fmt.Println("Program continues after panic recovery")
    
    // 测试多层 recover
    outerFunction()
}
```

## 最佳实践

### 1. 错误处理原则

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// 1. 总是检查错误
func badExample() {
    file, _ := os.Open("file.txt") // 错误：忽略了错误
    defer file.Close()
    // 处理文件...
}

func goodExample() error {
    file, err := os.Open("file.txt")
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    // 处理文件...
    return nil
}

// 2. 提供有意义的错误信息
func validateUser(name, email string) error {
    if name == "" {
        return errors.New("name is required")
    }
    
    if email == "" {
        return errors.New("email is required")
    }
    
    if len(email) < 5 {
        return fmt.Errorf("email too short: %s", email)
    }
    
    return nil
}

// 3. 使用错误包装保持上下文
func processUserData(userID string) error {
    err := validateUser("", "test@example.com")
    if err != nil {
        return fmt.Errorf("processing user %s: %w", userID, err)
    }
    
    return nil
}

// 4. 区分错误类型
func handleError(err error) {
    if err == nil {
        return
    }
    
    switch {
    case errors.Is(err, os.ErrNotExist):
        fmt.Println("File not found")
    case errors.Is(err, io.EOF):
        fmt.Println("End of file")
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}

func main() {
    // 测试错误处理
    if err := goodExample(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    if err := processUserData("123"); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // 测试错误类型处理
    handleError(os.ErrNotExist)
    handleError(io.EOF)
    handleError(errors.New("custom error"))
}
```

### 2. 错误处理工具函数

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// 错误处理工具函数
type ErrorUtils struct{}

// 1. 忽略错误（谨慎使用）
func (ErrorUtils) Ignore(err error) {
    if err != nil {
        // 记录日志但不返回错误
        fmt.Printf("Ignoring error: %v\n", err)
    }
}

// 2. 错误转换
func (ErrorUtils) Convert(err error, message string) error {
    if err == nil {
        return nil
    }
    return fmt.Errorf("%s: %w", message, err)
}

// 3. 错误聚合
func (ErrorUtils) Combine(errs ...error) error {
    var nonNilErrs []error
    for _, err := range errs {
        if err != nil {
            nonNilErrs = append(nonNilErrs, err)
        }
    }
    
    if len(nonNilErrs) == 0 {
        return nil
    }
    
    if len(nonNilErrs) == 1 {
        return nonNilErrs[0]
    }
    
    return fmt.Errorf("multiple errors: %v", nonNilErrs)
}

// 4. 错误重试
func (ErrorUtils) Retry(operation func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := operation(); err != nil {
            lastErr = err
            continue
        }
        return nil
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, lastErr)
}

func main() {
    utils := ErrorUtils{}
    
    // 测试错误转换
    originalErr := errors.New("original error")
    convertedErr := utils.Convert(originalErr, "conversion failed")
    fmt.Printf("Converted error: %v\n", convertedErr)
    
    // 测试错误聚合
    err1 := errors.New("error 1")
    err2 := errors.New("error 2")
    err3 := error(nil)
    
    combinedErr := utils.Combine(err1, err2, err3)
    fmt.Printf("Combined error: %v\n", combinedErr)
    
    // 测试错误重试
    attemptCount := 0
    retryErr := utils.Retry(func() error {
        attemptCount++
        if attemptCount < 3 {
            return errors.New("temporary failure")
        }
        return nil
    }, 5)
    
    if retryErr != nil {
        fmt.Printf("Retry failed: %v\n", retryErr)
    } else {
        fmt.Println("Retry succeeded")
    }
}
```

## 常见错误处理场景

### 1. 文件操作错误处理

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func copyFile(src, dst string) error {
    // 打开源文件
    srcFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("failed to open source file %s: %w", src, err)
    }
    defer srcFile.Close()
    
    // 创建目标文件
    dstFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("failed to create destination file %s: %w", dst, err)
    }
    defer dstFile.Close()
    
    // 复制文件内容
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return fmt.Errorf("failed to copy file content: %w", err)
    }
    
    // 确保数据写入磁盘
    err = dstFile.Sync()
    if err != nil {
        return fmt.Errorf("failed to sync file: %w", err)
    }
    
    return nil
}

func main() {
    err := copyFile("source.txt", "destination.txt")
    if err != nil {
        fmt.Printf("Copy failed: %v\n", err)
    } else {
        fmt.Println("File copied successfully")
    }
}
```

### 2. 网络请求错误处理

```go
package main

import (
    "errors"
    "fmt"
    "net"
    "net/http"
    "time"
)

// 网络错误类型
type NetworkError struct {
    Operation string
    URL       string
    Err       error
}

func (e NetworkError) Error() string {
    return fmt.Sprintf("network error during %s to %s: %v", e.Operation, e.URL, e.Err)
}

func (e NetworkError) Unwrap() error {
    return e.Err
}

// 网络请求函数
func makeRequest(url string) error {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    resp, err := client.Get(url)
    if err != nil {
        // 检查网络错误类型
        if netErr, ok := err.(net.Error); ok {
            if netErr.Timeout() {
                return NetworkError{
                    Operation: "GET",
                    URL:       url,
                    Err:       fmt.Errorf("request timeout: %w", err),
                }
            }
        }
        
        return NetworkError{
            Operation: "GET",
            URL:       url,
            Err:       err,
        }
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return NetworkError{
            Operation: "GET",
            URL:       url,
            Err:       fmt.Errorf("HTTP error: %d", resp.StatusCode),
        }
    }
    
    return nil
}

func main() {
    // 测试网络请求
    err := makeRequest("https://httpbin.org/status/200")
    if err != nil {
        fmt.Printf("Request failed: %v\n", err)
    } else {
        fmt.Println("Request successful")
    }
    
    // 测试错误请求
    err = makeRequest("https://httpbin.org/status/404")
    if err != nil {
        fmt.Printf("Request failed: %v\n", err)
    }
}
```

### 3. 数据库操作错误处理

```go
package main

import (
    "errors"
    "fmt"
)

// 数据库错误类型
type DatabaseError struct {
    Operation string
    Table     string
    Err       error
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("database error during %s on table %s: %v", e.Operation, e.Table, e.Err)
}

func (e DatabaseError) Unwrap() error {
    return e.Err
}

// 模拟数据库操作
func insertUser(userID, name string) error {
    // 模拟数据库错误
    if userID == "" {
        return DatabaseError{
            Operation: "INSERT",
            Table:     "users",
            Err:       errors.New("user ID cannot be empty"),
        }
    }
    
    if name == "" {
        return DatabaseError{
            Operation: "INSERT",
            Table:     "users",
            Err:       errors.New("name cannot be empty"),
        }
    }
    
    // 模拟唯一约束错误
    if userID == "duplicate" {
        return DatabaseError{
            Operation: "INSERT",
            Table:     "users",
            Err:       errors.New("duplicate key: user already exists"),
        }
    }
    
    return nil
}

func main() {
    // 测试数据库操作
    testCases := []struct {
        userID string
        name   string
    }{
        {"", "Alice"},
        {"123", ""},
        {"duplicate", "Bob"},
        {"456", "Charlie"},
    }
    
    for _, tc := range testCases {
        err := insertUser(tc.userID, tc.name)
        if err != nil {
            fmt.Printf("Insert failed for userID=%s, name=%s: %v\n", tc.userID, tc.name, err)
        } else {
            fmt.Printf("Insert successful for userID=%s, name=%s\n", tc.userID, tc.name)
        }
    }
}
```

## 总结

Go 的错误处理具有以下特点：

1. **显式错误处理**: 错误作为返回值，必须显式处理
2. **错误包装**: 使用 `fmt.Errorf` 和 `%w` 保持错误上下文
3. **类型安全**: 通过类型断言和 `errors.As` 检查错误类型
4. **错误链**: 使用 `errors.Unwrap` 和 `errors.Is` 处理错误链
5. **panic/recover**: 用于不可恢复的错误和程序逻辑错误

**核心概念**:
- 错误作为值返回
- 错误包装和展开
- 自定义错误类型
- 错误处理模式
- panic 和 recover

**最佳实践**:
- 总是检查错误
- 提供有意义的错误信息
- 使用错误包装保持上下文
- 区分不同类型的错误
- 谨慎使用 panic 和 recover

掌握这些错误处理技巧，可以编写出更加健壮和可维护的 Go 代码。
