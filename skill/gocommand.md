# Go 命令知识概览

## 基础命令

### 1. 构建和运行
- **`go build`** - 编译包和依赖，但不安装结果
- **`go run`** - 编译并运行 Go 程序
- **`go install`** - 编译并安装包和依赖
- **`go clean`** - 删除对象文件和缓存文件

### 2. 模块管理
- **`go mod init`** - 初始化新模块
- **`go mod tidy`** - 添加缺失的模块，删除未使用的模块
- **`go mod download`** - 下载模块到本地缓存
- **`go mod vendor`** - 创建依赖的 vendored 副本
- **`go mod verify`** - 验证依赖具有预期内容
- **`go mod graph`** - 打印模块需求图
- **`go mod edit`** - 从工具或脚本编辑 go.mod
- **`go mod why`** - 解释为什么需要包或模块

### 3. 依赖管理
- **`go get`** - 添加依赖到当前模块并安装它们
- **`go list`** - 列出包或模块

### 4. 测试
- **`go test`** - 测试包
- **`go vet`** - 报告包中可能的错误

### 5. 代码质量
- **`go fmt`** - 格式化包源代码 (gofmt)
- **`go doc`** - 显示包或符号的文档
- **`go fix`** - 更新包以使用新的 API

### 6. 工具链
- **`go tool`** - 运行指定的 go 工具
- **`go generate`** - 通过处理源代码生成 Go 文件

### 7. 环境信息
- **`go env`** - 打印 Go 环境信息
- **`go version`** - 打印 Go 版本

### 8. 工作区管理
- **`go work`** - 工作区维护

### 9. 调试和报告
- **`go bug`** - 开始错误报告
- **`go telemetry`** - 管理遥测数据和设置

## 详细命令说明

### go bug
**功能**: 启动一个错误报告，自动打开浏览器访问 Go 项目的错误报告页面

**语法**: `go bug`

**详细说明**:
- 自动检测 Go 版本、操作系统、架构等信息
- 打开默认浏览器访问 https://github.com/golang/go/issues/new
- 预填充错误报告模板，包含系统信息
- 帮助开发者快速提交 bug 报告

**使用场景**:
- 发现 Go 工具链或标准库的 bug
- 需要向 Go 团队报告问题
- 获取标准化的错误报告格式

**示例**:
```bash
# 启动错误报告
go bug
# 会自动打开浏览器并预填充系统信息
```

### go build
**功能**: 编译包和依赖，生成可执行文件或库文件

**语法**: `go build [build flags] [packages]`

**详细说明**:
- 编译指定的包及其依赖
- 默认在当前目录编译，生成与目录同名的可执行文件
- 支持交叉编译（通过 GOOS 和 GOARCH 环境变量）
- 可以编译多个包或整个模块
- 不会安装编译结果到 $GOPATH/bin 或 $GOBIN

**常用标志**:
- `-o file` - 指定输出文件名
- `-v` - 显示编译的包名
- `-race` - 启用竞态检测
- `-ldflags` - 传递给链接器的标志
- `-tags` - 构建标签
- `-trimpath` - 移除文件路径信息

**使用场景**:
- 编译 Go 程序生成可执行文件
- 测试编译是否成功
- 交叉编译到不同平台
- 生成静态链接的二进制文件

**示例**:
```bash
# 编译当前目录的包
go build

# 编译指定文件
go build main.go

# 指定输出文件名
go build -o myapp

# 编译整个模块
go build ./...

# 交叉编译到 Linux
GOOS=linux GOARCH=amd64 go build

# 启用竞态检测
go build -race

# 传递链接器标志
go build -ldflags="-s -w" -o myapp
```

### go clean
**功能**: 删除对象文件和缓存文件，清理构建产物

**语法**: `go clean [clean flags] [packages]`

**详细说明**:
- 删除由 go build 生成的对象文件和可执行文件
- 清理 Go 模块缓存和构建缓存
- 可以清理特定包或整个模块
- 不会删除源代码文件
- 有助于解决构建问题和节省磁盘空间

**常用标志**:
- `-cache` - 删除整个构建缓存
- `-modcache` - 删除模块下载缓存
- `-testcache` - 删除测试结果缓存
- `-fuzzcache` - 删除模糊测试缓存
- `-i` - 删除 go install 安装的包
- `-r` - 递归清理导入的包
- `-n` - 显示将要执行的命令但不执行
- `-x` - 显示执行的命令

**使用场景**:
- 解决构建问题（如缓存损坏）
- 释放磁盘空间
- 强制重新编译所有依赖
- 清理测试缓存
- 重置模块缓存

**示例**:
```bash
# 清理当前目录的构建文件
go clean

# 清理指定包
go clean ./pkg1 ./pkg2

# 清理整个模块
go clean ./...

# 删除构建缓存
go clean -cache

# 删除模块缓存
go clean -modcache

# 删除测试缓存
go clean -testcache

# 递归清理所有依赖
go clean -r

# 显示将要执行的清理命令
go clean -n

# 清理并显示执行的命令
go clean -x
```

### go doc
**功能**: 显示包或符号的文档注释

**语法**: `go doc [doc flags] [package|[package.]symbol[.methodOrField]]`

**详细说明**:
- 显示 Go 包、函数、类型、变量等的文档注释
- 支持显示标准库和第三方包的文档
- 可以显示包级别的文档或特定符号的文档
- 支持显示方法、字段等详细信息
- 是 Go 代码文档的标准查看工具

**常用标志**:
- `-all` - 显示所有导出的符号
- `-short` - 显示简短格式（不显示示例）
- `-src` - 显示源代码
- `-u` - 显示未导出的符号
- `-c` - 匹配时区分大小写
- `-cmd` - 将命令视为包（用于 main 包）

**使用场景**:
- 查看标准库 API 文档
- 了解第三方包的使用方法
- 查看自己代码的文档注释
- 学习 Go 标准库的使用
- 快速查找函数签名和用法

**示例**:
```bash
# 查看 fmt 包的文档
go doc fmt

# 查看 fmt.Printf 函数的文档
go doc fmt.Printf

# 查看 strings 包的所有导出符号
go doc -all strings

# 查看 time.Time 类型的文档
go doc time.Time

# 查看 time.Time.Format 方法的文档
go doc time.Time.Format

# 显示源代码
go doc -src fmt.Printf

# 查看 main 包的文档
go doc -cmd .

# 查看当前包的文档
go doc .

# 查看特定符号的简短文档
go doc -short fmt.Printf
```

### go env
**功能**: 打印 Go 环境信息，显示或设置环境变量

**语法**: `go env [env flags] [var ...]`

**详细说明**:
- 显示 Go 工具链的环境变量和配置信息
- 可以查看特定环境变量的值
- 支持设置环境变量（使用 -w 标志）
- 显示 Go 安装路径、版本、架构等信息
- 帮助诊断 Go 环境配置问题

**常用标志**:
- `-json` - 以 JSON 格式输出
- `-u` - 取消设置环境变量
- `-w` - 写入环境变量到配置文件
- `-v` - 显示详细信息

**重要环境变量**:
- `GOROOT` - Go 安装根目录
- `GOPATH` - Go 工作空间路径
- `GOPROXY` - 模块代理设置
- `GOSUMDB` - 校验和数据库
- `GOOS` - 目标操作系统
- `GOARCH` - 目标架构
- `CGO_ENABLED` - 是否启用 CGO
- `GO111MODULE` - 模块模式设置

**使用场景**:
- 检查 Go 环境配置
- 诊断构建问题
- 设置开发环境
- 查看 Go 版本和路径信息
- 配置模块代理和校验和数据库

**示例**:
```bash
# 显示所有环境变量
go env

# 显示特定环境变量
go env GOPATH GOROOT

# 以 JSON 格式显示
go env -json

# 设置模块代理
go env -w GOPROXY=https://goproxy.cn,direct

# 设置校验和数据库
go env -w GOSUMDB=sum.golang.org

# 启用模块模式
go env -w GO111MODULE=on

# 取消设置环境变量
go env -u GOPROXY

# 显示详细信息
go env -v

# 查看 Go 版本信息
go env GOVERSION

# 查看目标平台
go env GOOS GOARCH
```

### go fix
**功能**: 运行 go tool fix 来修复包，将旧版本的 Go 代码更新到新版本

**语法**: `go fix [packages]`

**详细说明**:
- 自动修复 Go 代码中的过时语法和 API
- 将旧版本 Go 代码迁移到新版本
- 应用一系列预定义的修复规则
- 支持批量修复多个包
- 是 Go 版本升级的重要工具

**修复内容**:
- 过时的包导入路径
- 废弃的 API 调用
- 语法变更
- 标准库 API 变化
- 编译器要求的变更

**使用场景**:
- 升级 Go 版本时迁移代码
- 修复过时的 API 使用
- 批量更新项目代码
- 确保代码兼容新版本 Go
- 自动化代码迁移

**注意事项**:
- 修复前建议备份代码
- 某些修复可能需要手动调整
- 建议在修复后运行测试
- 某些第三方库可能需要单独处理

**示例**:
```bash
# 修复当前目录的包
go fix

# 修复指定包
go fix ./pkg1 ./pkg2

# 修复整个模块
go fix ./...

# 修复特定文件
go fix file.go

# 修复并显示详细信息
go fix -v

# 修复前查看将要执行的修复
go fix -n

# 修复并显示执行的命令
go fix -x
```

### go fmt
**功能**: 格式化 Go 包中的源代码，统一代码风格

**语法**: `go fmt [packages]`

**详细说明**:
- 自动格式化 Go 源代码文件
- 统一代码缩进、空格、换行等格式
- 应用 Go 官方的代码格式化规则
- 支持批量格式化多个包
- 是 Go 代码风格的标准工具

**格式化规则**:
- 统一缩进（使用 tab）
- 调整空格和换行
- 对齐注释
- 规范化导入语句
- 调整运算符周围的空格
- 统一括号和逗号的使用

**常用标志**:
- `-n` - 显示将要执行的命令但不执行
- `-x` - 显示执行的命令
- `-w` - 将结果写入源文件（默认只输出到标准输出）
- `-l` - 列出需要格式化的文件

**使用场景**:
- 统一团队代码风格
- 提交代码前格式化
- 批量整理代码格式
- 确保代码符合 Go 标准
- 自动化代码美化

**最佳实践**:
- 在提交代码前运行 go fmt
- 配置编辑器自动格式化
- 在 CI/CD 流程中包含格式化检查
- 团队统一使用相同的格式化规则

**示例**:
```bash
# 格式化当前目录的包
go fmt

# 格式化指定包
go fmt ./pkg1 ./pkg2

# 格式化整个模块
go fmt ./...

# 格式化特定文件
go fmt file.go

# 显示将要执行的命令
go fmt -n

# 显示执行的命令
go fmt -x

# 将结果写入源文件
go fmt -w

# 列出需要格式化的文件
go fmt -l

# 格式化并写入文件
go fmt -w ./...
```

### go generate
**功能**: 通过扫描源代码中的特殊注释来生成代码

**语法**: `go generate [build flags] [file.go... | packages]`

**详细说明**:
- 扫描源代码中的 `//go:generate` 注释
- 执行注释中指定的命令来生成代码
- 支持多种代码生成工具和模板
- 常用于生成重复性代码、接口实现等
- 是 Go 代码生成的标准方式

**特殊注释格式**:
```go
//go:generate command argument...
```

**常用生成工具**:
- `stringer` - 为枚举类型生成 String() 方法
- `mockgen` - 生成接口的 mock 实现
- `protoc` - 从 Protocol Buffers 生成 Go 代码
- `swag` - 生成 Swagger 文档
- `wire` - 生成依赖注入代码
- `sqlc` - 从 SQL 生成类型安全的 Go 代码

**常用标志**:
- `-n` - 显示将要执行的命令但不执行
- `-x` - 显示执行的命令
- `-v` - 显示详细信息
- `-run regexp` - 只运行匹配正则表达式的命令

**使用场景**:
- 生成重复性代码
- 创建接口的 mock 实现
- 从协议文件生成代码
- 生成 API 文档
- 创建依赖注入代码
- 生成数据库访问层

**最佳实践**:
- 在生成的文件中添加 "DO NOT EDIT" 注释
- 将生成的文件加入版本控制
- 在构建流程中包含代码生成步骤
- 使用有意义的生成命令注释

**示例**:
```bash
# 生成当前包的所有代码
go generate

# 生成指定包
go generate ./pkg1 ./pkg2

# 生成特定文件
go generate file.go

# 显示将要执行的命令
go generate -n

# 显示执行的命令
go generate -x

# 显示详细信息
go generate -v

# 只运行匹配的命令
go generate -run "stringer"

# 生成整个模块
go generate ./...
```

**源代码示例**:
```go
//go:generate stringer -type=Pill
type Pill int

const (
    Placebo Pill = iota
    Aspirin
    Ibuprofen
    Paracetamol
)

//go:generate mockgen -destination=mock_mock.go -package=mock github.com/example/interface
type MyInterface interface {
    DoSomething() error
}
```

### go get
**功能**: 下载并安装包和依赖项

**语法**: `go get [build flags] [packages]`

**详细说明**:
- 下载指定的包及其依赖项
- 将包安装到 GOPATH 或模块缓存中
- 更新 go.mod 文件中的依赖版本
- 支持从 Git、SVN、Mercurial 等版本控制系统下载
- 是 Go 包管理的核心命令

**常用标志**:
- `-d` - 只下载包，不安装
- `-u` - 更新包及其依赖项到最新版本
- `-u=patch` - 只更新补丁版本
- `-t` - 同时下载测试依赖
- `-insecure` - 允许不安全的 HTTP 连接
- `-v` - 显示详细信息
- `-x` - 显示执行的命令

**版本控制**:
- 支持 Git、SVN、Mercurial、Bazaar
- 自动识别版本控制系统
- 支持分支、标签、提交哈希
- 支持私有仓库（需要认证）

**使用场景**:
- 安装第三方包
- 更新依赖版本
- 添加新的依赖项
- 下载特定版本的包
- 安装开发工具
- 更新项目依赖

**模块模式**:
- 在模块模式下，go get 会更新 go.mod 文件
- 支持语义化版本控制
- 自动解析依赖关系
- 支持 replace 和 exclude 指令

**最佳实践**:
- 使用 go.mod 管理依赖
- 定期更新依赖版本
- 使用 -u 标志更新依赖
- 检查依赖的安全性
- 使用 go.sum 验证依赖完整性

**示例**:
```bash
# 安装包
go get github.com/gin-gonic/gin

# 安装特定版本
go get github.com/gin-gonic/gin@v1.9.0

# 安装最新版本
go get -u github.com/gin-gonic/gin

# 只下载不安装
go get -d github.com/gin-gonic/gin

# 更新所有依赖
go get -u ./...

# 更新补丁版本
go get -u=patch ./...

# 安装开发工具
go get golang.org/x/tools/cmd/goimports

# 显示详细信息
go get -v github.com/gin-gonic/gin

# 显示执行的命令
go get -x github.com/gin-gonic/gin

# 安装测试依赖
go get -t ./...

# 从特定分支安装
go get github.com/user/repo@branch-name

# 从特定提交安装
go get github.com/user/repo@commit-hash
```

### go install
**功能**: 编译并安装包和可执行文件

**语法**: `go install [build flags] [packages]`

**详细说明**:
- 编译指定的包并安装到 GOBIN 目录
- 生成可执行文件或库文件
- 在模块模式下，安装到 $GOPATH/bin 或 $GOBIN
- 支持交叉编译和条件编译
- 是 Go 程序部署的标准方式

**安装位置**:
- 可执行文件: $GOBIN 或 $GOPATH/bin
- 库文件: $GOPATH/pkg
- 模块模式下: $GOPATH/bin 或 $GOBIN

**常用标志**:
- `-v` - 显示详细信息
- `-x` - 显示执行的命令
- `-race` - 启用竞态检测
- `-msan` - 启用内存清理器
- `-asan` - 启用地址清理器
- `-buildmode` - 指定构建模式
- `-ldflags` - 传递链接器标志
- `-tags` - 指定构建标签

**构建模式**:
- `exe` - 可执行文件（默认）
- `shared` - 共享库
- `c-archive` - C 归档文件
- `c-shared` - C 共享库
- `plugin` - Go 插件

**使用场景**:
- 安装开发工具
- 构建可执行程序
- 安装第三方工具
- 部署 Go 应用
- 创建系统工具
- 构建库文件

**与 go build 的区别**:
- go install 会安装到标准位置
- go build 只编译到当前目录
- go install 适合工具安装
- go build 适合开发调试

**最佳实践**:
- 使用 go install 安装工具
- 设置合适的 GOBIN 路径
- 使用版本标签管理工具版本
- 在 CI/CD 中使用 go install
- 定期清理不需要的工具

**示例**:
```bash
# 安装当前包
go install

# 安装指定包
go install ./cmd/server

# 安装多个包
go install ./cmd/... ./pkg/...

# 安装远程包
go install github.com/gin-gonic/gin

# 安装特定版本
go install github.com/gin-gonic/gin@v1.9.0

# 显示详细信息
go install -v ./cmd/server

# 显示执行的命令
go install -x ./cmd/server

# 启用竞态检测
go install -race ./cmd/server

# 使用构建标签
go install -tags=debug ./cmd/server

# 传递链接器标志
go install -ldflags="-s -w" ./cmd/server

# 构建共享库
go install -buildmode=shared ./pkg/mylib

# 构建 C 归档文件
go install -buildmode=c-archive ./pkg/mylib

# 安装开发工具
go install golang.org/x/tools/cmd/goimports@latest

# 安装整个模块
go install ./...
```

### go list
**功能**: 列出包或模块信息

**语法**: `go list [build flags] [-f format] [-json] [packages]`

**详细说明**:
- 显示包或模块的详细信息
- 支持自定义输出格式
- 可以输出 JSON 格式的数据
- 用于查询包的结构和依赖关系
- 是 Go 包信息查询的核心工具

**输出格式**:
- 默认格式: 包路径
- JSON 格式: 使用 -json 标志
- 自定义格式: 使用 -f 标志和模板语法

**常用标志**:
- `-f format` - 指定输出格式
- `-json` - 输出 JSON 格式
- `-m` - 列出模块而不是包
- `-u` - 显示可更新的依赖
- `-versions` - 显示可用版本
- `-deps` - 显示依赖关系
- `-test` - 包含测试文件
- `-e` - 显示错误信息

**模板变量**:
- `.Name` - 包名
- `.ImportPath` - 导入路径
- `.Dir` - 包目录
- `.GoFiles` - Go 源文件
- `.TestGoFiles` - 测试文件
- `.Imports` - 导入的包
- `.Deps` - 依赖的包
- `.Module` - 模块信息

**使用场景**:
- 查询包信息
- 分析依赖关系
- 检查模块状态
- 查找可更新的依赖
- 生成包列表
- 调试包问题

**最佳实践**:
- 使用 -json 获取完整信息
- 使用模板自定义输出
- 结合其他工具分析依赖
- 定期检查依赖更新
- 使用 -e 调试包错误

**示例**:
```bash
# 列出当前包
go list

# 列出指定包
go list ./pkg1 ./pkg2

# 列出所有包
go list ./...

# 输出 JSON 格式
go list -json ./pkg1

# 自定义输出格式
go list -f '{{.Name}} -> {{.ImportPath}}' ./pkg1

# 显示包目录
go list -f '{{.Dir}}' ./pkg1

# 显示源文件
go list -f '{{.GoFiles}}' ./pkg1

# 显示导入的包
go list -f '{{.Imports}}' ./pkg1

# 显示依赖关系
go list -deps ./pkg1

# 列出模块
go list -m

# 列出所有模块
go list -m all

# 显示可更新的依赖
go list -m -u all

# 显示模块版本
go list -m -versions github.com/gin-gonic/gin

# 显示测试文件
go list -f '{{.TestGoFiles}}' ./pkg1

# 显示错误信息
go list -e ./pkg1

# 显示包和测试包
go list -test ./pkg1

# 自定义格式显示多个字段
go list -f 'Package: {{.Name}}\nPath: {{.ImportPath}}\nDir: {{.Dir}}' ./pkg1

# 显示模块路径
go list -f '{{.Module.Path}}' ./pkg1

# 显示模块版本
go list -f '{{.Module.Version}}' ./pkg1

# 显示模块替换
go list -f '{{.Module.Replace}}' ./pkg1
```

### go mod
**功能**: 管理 Go 模块

**语法**: `go mod [command] [arguments]`

**详细说明**:
- 管理 Go 模块的依赖关系
- 操作 go.mod 和 go.sum 文件
- 支持模块的初始化、编辑、下载等操作
- 是 Go 模块系统的核心管理工具
- 提供多种子命令处理不同任务

**子命令**:
- `init` - 初始化新模块
- `download` - 下载模块依赖
- `edit` - 编辑 go.mod 文件
- `graph` - 打印模块依赖图
- `verify` - 验证依赖
- `why` - 解释为什么需要包或模块
- `tidy` - 整理和清理模块依赖
- `vendor` - 创建 vendor 目录

**go mod init**:
- 初始化新的 Go 模块
- 创建 go.mod 文件
- 设置模块路径和 Go 版本
- 支持自定义模块路径

**go mod download**:
- 下载模块到本地缓存
- 验证模块完整性
- 支持特定模块下载
- 更新 go.sum 文件

**go mod edit**:
- 编辑 go.mod 文件
- 添加、删除、替换模块
- 修改 Go 版本要求
- 支持批量操作

**go mod graph**:
- 显示模块依赖关系图
- 用于分析依赖结构
- 支持过滤特定模块
- 输出格式便于处理

**go mod verify**:
- 验证模块的完整性
- 检查 go.sum 文件
- 确保依赖未被篡改
- 报告验证结果

**go mod why**:
- 解释包或模块的依赖原因
- 显示依赖路径
- 帮助理解依赖关系
- 支持递归分析

**go mod tidy**:
- 添加缺失的依赖
- 删除未使用的依赖
- 更新 go.mod 和 go.sum
- 确保依赖一致性

**go mod vendor**:
- 创建 vendor 目录
- 复制依赖到本地
- 支持离线构建
- 确保构建一致性

**使用场景**:
- 初始化新项目
- 管理项目依赖
- 分析依赖关系
- 验证模块完整性
- 创建离线构建环境
- 调试依赖问题

**最佳实践**:
- 使用语义化版本
- 定期运行 go mod tidy
- 验证依赖完整性
- 使用 vendor 目录确保一致性
- 记录依赖变更
- 定期更新依赖

**示例**:
```bash
# 初始化模块
go mod init myproject

# 初始化带路径的模块
go mod init github.com/user/myproject

# 下载依赖
go mod download

# 下载特定模块
go mod download github.com/gin-gonic/gin

# 编辑 go.mod
go mod edit -require=github.com/gin-gonic/gin@v1.9.0

# 添加替换
go mod edit -replace=old/path=new/path

# 显示依赖图
go mod graph

# 验证依赖
go mod verify

# 解释依赖
go mod why github.com/gin-gonic/gin

# 整理依赖
go mod tidy

# 创建 vendor 目录
go mod vendor

# 显示模块信息
go mod edit -print

# 设置 Go 版本
go mod edit -go=1.21

# 添加排除
go mod edit -exclude=github.com/old/module@v1.0.0

# 批量编辑
go mod edit -droprequire=github.com/old/module

# 显示依赖路径
go mod why -m github.com/gin-gonic/gin

# 验证特定模块
go mod verify github.com/gin-gonic/gin

# 下载并验证
go mod download -x

# 显示模块图
go mod graph | grep gin
```

### go work
**功能**: 管理 Go 工作区

**语法**: `go work [command] [arguments]`

**详细说明**:
- 管理多模块 Go 工作区
- 操作 go.work 文件
- 支持工作区的初始化、编辑、同步等操作
- 是 Go 1.18+ 引入的工作区功能
- 用于管理包含多个模块的项目

**子命令**:
- `init` - 初始化工作区
- `use` - 添加模块到工作区
- `edit` - 编辑 go.work 文件
- `sync` - 同步工作区模块
- `vendor` - 为工作区创建 vendor 目录

**go work init**:
- 初始化新的 Go 工作区
- 创建 go.work 文件
- 设置工作区配置
- 支持指定初始模块

**go work use**:
- 添加模块到工作区
- 支持相对路径和绝对路径
- 自动更新 go.work 文件
- 支持批量添加模块

**go work edit**:
- 编辑 go.work 文件
- 添加、删除、替换模块
- 修改工作区配置
- 支持批量操作

**go work sync**:
- 同步工作区中的模块
- 更新模块依赖
- 确保工作区一致性
- 处理模块间依赖关系

**go work vendor**:
- 为工作区创建 vendor 目录
- 复制所有模块的依赖
- 支持离线构建
- 确保构建一致性

**工作区结构**:
- go.work 文件定义工作区
- 包含多个模块路径
- 支持模块替换和排除
- 与 go.mod 文件配合使用

**使用场景**:
- 管理多模块项目
- 开发大型应用
- 处理模块间依赖
- 统一管理相关模块
- 简化开发工作流
- 支持微服务架构

**最佳实践**:
- 合理组织工作区结构
- 定期同步工作区
- 使用相对路径引用模块
- 保持模块独立性
- 记录工作区变更
- 测试工作区配置

**示例**:
```bash
# 初始化工作区
go work init

# 初始化带模块的工作区
go work init ./module1 ./module2

# 添加模块到工作区
go work use ./module1

# 添加多个模块
go work use ./module1 ./module2 ./module3

# 编辑工作区
go work edit -use=./newmodule

# 替换模块
go work edit -replace=old/path=new/path

# 同步工作区
go work sync

# 创建 vendor 目录
go work vendor

# 显示工作区信息
go work edit -print

# 删除模块
go work edit -dropuse=./oldmodule

# 添加排除
go work edit -exclude=module@v1.0.0

# 批量编辑
go work edit -dropreplace=old/path

# 使用绝对路径
go work use /absolute/path/to/module

# 同步特定模块
go work sync ./module1

# 显示工作区图
go work graph

# 验证工作区
go work verify

# 显示工作区依赖
go work why ./module1

# 下载工作区依赖
go work download

# 整理工作区
go work tidy
```

### go run
**功能**: 编译并运行 Go 程序

**语法**: `go run [build flags] [-exec xprog] package [arguments...]`

**详细说明**:
- 编译指定的包并立即运行
- 支持传递命令行参数给程序
- 自动处理依赖关系
- 临时编译，不生成可执行文件
- 适合快速测试和开发

**编译过程**:
- 解析包依赖
- 编译源代码
- 链接生成可执行文件
- 运行程序
- 清理临时文件

**支持的文件类型**:
- .go 源文件
- 包目录
- 多个源文件
- 工作区模块

**常用标志**:
- `-v` - 显示编译过程
- `-x` - 显示执行的命令
- `-race` - 启用竞态检测
- `-msan` - 启用内存清理器
- `-asan` - 启用地址清理器
- `-exec xprog` - 使用指定的执行器
- `-ldflags` - 传递链接器标志
- `-tags` - 指定构建标签

**性能优化**:
- 支持增量编译
- 缓存编译结果
- 并行编译
- 智能依赖分析

**错误处理**:
- 显示编译错误
- 显示运行时错误
- 提供详细错误信息
- 支持错误定位

**使用场景**:
- 快速测试代码
- 开发调试
- 运行示例程序
- 执行脚本
- 测试新功能
- 演示程序

**最佳实践**:
- 使用 -v 查看编译过程
- 使用 -race 检测竞态条件
- 传递合适的参数
- 处理程序退出码
- 使用构建标签控制编译
- 合理组织代码结构

**示例**:
```bash
# 运行单个文件
go run main.go

# 运行包
go run ./cmd/server

# 运行多个文件
go run file1.go file2.go

# 传递参数
go run main.go arg1 arg2

# 显示编译过程
go run -v main.go

# 显示执行的命令
go run -x main.go

# 启用竞态检测
go run -race main.go

# 使用构建标签
go run -tags=debug main.go

# 传递链接器标志
go run -ldflags="-s -w" main.go

# 使用自定义执行器
go run -exec=time main.go

# 运行工作区模块
go run ./module1/cmd/app

# 运行带环境变量
GOOS=linux go run main.go

# 运行测试文件
go run *_test.go

# 运行示例
go run example_test.go

# 运行基准测试
go run -bench=. *_test.go

# 运行特定测试
go run -run=TestFunction *_test.go

# 运行带覆盖率的测试
go run -cover *_test.go

# 运行带性能分析的测试
go run -cpuprofile=cpu.prof *_test.go

# 运行带内存分析的测试
go run -memprofile=mem.prof *_test.go

# 运行带阻塞分析的测试
go run -blockprofile=block.prof *_test.go

# 运行带互斥分析的测试
go run -mutexprofile=mutex.prof *_test.go

# 运行带跟踪的测试
go run -trace=trace.out *_test.go

# 运行带调试信息的程序
go run -gcflags="-N -l" main.go

# 运行优化版本
go run -gcflags="-O2" main.go

# 运行带内联信息的程序
go run -gcflags="-m" main.go

# 运行带逃逸分析的程序
go run -gcflags="-m -m" main.go
```

### go build
**功能**: 编译 Go 包和依赖

**语法**: `go build [build flags] [packages]`

**详细说明**:
- 编译指定的包和依赖
- 生成可执行文件或库文件
- 支持交叉编译
- 自动处理依赖关系
- 支持多种输出格式

**编译目标**:
- 可执行文件（main 包）
- 静态库（非 main 包）
- 动态库（使用 -buildmode）
- 插件（使用 -buildmode=plugin）

**输出控制**:
- `-o` - 指定输出文件名
- `-buildmode` - 指定构建模式
- `-ldflags` - 传递链接器标志
- `-tags` - 指定构建标签
- `-trimpath` - 移除文件路径信息

**构建模式**:
- `exe` - 可执行文件（默认）
- `archive` - 静态库
- `c-archive` - C 静态库
- `c-shared` - C 动态库
- `shared` - Go 动态库
- `plugin` - Go 插件

**优化选项**:
- `-gcflags` - 传递编译器标志
- `-ldflags` - 传递链接器标志
- `-asmflags` - 传递汇编器标志
- `-trimpath` - 移除调试路径
- `-race` - 启用竞态检测
- `-msan` - 启用内存清理器

**交叉编译**:
- `GOOS` - 目标操作系统
- `GOARCH` - 目标架构
- `CGO_ENABLED` - 启用 CGO
- `CC` - C 编译器
- `CXX` - C++ 编译器

**调试支持**:
- `-gcflags="-N -l"` - 禁用优化和内联
- `-ldflags="-s -w"` - 移除调试信息
- `-trimpath` - 移除路径信息
- `-race` - 竞态检测

**性能分析**:
- `-cpuprofile` - CPU 性能分析
- `-memprofile` - 内存性能分析
- `-blockprofile` - 阻塞性能分析
- `-mutexprofile` - 互斥性能分析

**使用场景**:
- 构建生产程序
- 交叉编译
- 创建库文件
- 构建插件
- 性能优化
- 调试程序

**最佳实践**:
- 使用语义化版本
- 设置合适的构建标签
- 优化二进制大小
- 处理交叉编译
- 使用构建缓存
- 合理组织代码结构

**示例**:
```bash
# 构建当前包
go build

# 构建指定包
go build ./cmd/server

# 构建多个包
go build ./cmd/server ./cmd/client

# 指定输出文件
go build -o myapp main.go

# 构建静态库
go build -buildmode=archive -o lib.a ./pkg

# 构建 C 静态库
go build -buildmode=c-archive -o lib.a ./pkg

# 构建 C 动态库
go build -buildmode=c-shared -o lib.so ./pkg

# 构建插件
go build -buildmode=plugin -o plugin.so ./plugin

# 交叉编译 Linux
GOOS=linux GOARCH=amd64 go build -o app main.go

# 交叉编译 Windows
GOOS=windows GOARCH=amd64 go build -o app.exe main.go

# 交叉编译 ARM
GOOS=linux GOARCH=arm go build -o app main.go

# 使用构建标签
go build -tags=debug main.go

# 传递链接器标志
go build -ldflags="-s -w" main.go

# 传递编译器标志
go build -gcflags="-N -l" main.go

# 启用竞态检测
go build -race main.go

# 移除路径信息
go build -trimpath main.go

# 构建带版本信息
go build -ldflags="-X main.Version=1.0.0" main.go

# 构建带调试信息
go build -gcflags="-N -l" -ldflags="-s -w" main.go

# 构建优化版本
go build -gcflags="-O2" main.go

# 构建带内联信息
go build -gcflags="-m" main.go

# 构建带逃逸分析
go build -gcflags="-m -m" main.go

# 构建带性能分析
go build -cpuprofile=cpu.prof main.go

# 构建带内存分析
go build -memprofile=mem.prof main.go

# 构建带阻塞分析
go build -blockprofile=block.prof main.go

# 构建带互斥分析
go build -mutexprofile=mutex.prof main.go

# 构建带跟踪
go build -trace=trace.out main.go

# 构建带覆盖率的测试
go build -cover main.go

# 构建带基准测试
go build -bench=. main.go

# 构建带特定测试
go build -run=TestFunction main.go

# 构建带覆盖率分析的测试
go build -coverprofile=coverage.out main.go

# 构建带覆盖率模式的测试
go build -covermode=atomic main.go

# 构建带覆盖率函数的测试
go build -coverpkg=./... main.go

# 构建带覆盖率输出的测试
go build -coverdir=coverage main.go

# 构建带覆盖率格式的测试
go build -coverformat=html main.go

# 构建带覆盖率阈值的测试
go build -coverthreshold=80 main.go

# 构建带覆盖率排除的测试
go build -coverexclude=testdata main.go

# 构建带覆盖率包含的测试
go build -coverinclude=*.go main.go

# 构建带覆盖率排除模式的测试
go build -coverexclude=*_test.go main.go

# 构建带覆盖率包含模式的测试
go build -coverinclude=main.go main.go
```

### go install
**功能**: 编译并安装 Go 包

**语法**: `go install [build flags] [packages]`

**详细说明**:
- 编译指定的包和依赖
- 将可执行文件安装到 GOBIN 目录
- 将库文件安装到 GOPATH/pkg 目录
- 自动处理依赖关系
- 支持版本管理

**安装位置**:
- 可执行文件：`$GOBIN` 或 `$GOPATH/bin`
- 库文件：`$GOPATH/pkg`
- 缓存文件：`$GOCACHE`
- 模块缓存：`$GOMODCACHE`

**版本管理**:
- 支持语义化版本
- 自动解析依赖版本
- 支持版本约束
- 处理版本冲突
- 支持模块替换

**安装类型**:
- 可执行文件（main 包）
- 库文件（非 main 包）
- 工具和命令
- 插件和扩展

**常用标志**:
- `-v` - 显示编译过程
- `-x` - 显示执行的命令
- `-race` - 启用竞态检测
- `-ldflags` - 传递链接器标志
- `-tags` - 指定构建标签
- `-trimpath` - 移除路径信息

**环境变量**:
- `GOBIN` - 可执行文件安装目录
- `GOPATH` - Go 工作空间
- `GOCACHE` - 编译缓存目录
- `GOMODCACHE` - 模块缓存目录
- `GOOS` - 目标操作系统
- `GOARCH` - 目标架构

**依赖管理**:
- 自动下载依赖
- 解析版本约束
- 处理间接依赖
- 支持模块替换
- 处理版本冲突

**缓存机制**:
- 编译结果缓存
- 依赖下载缓存
- 增量编译支持
- 并行编译优化
- 智能缓存失效

**使用场景**:
- 安装工具和命令
- 部署应用程序
- 安装开发工具
- 管理项目依赖
- 构建发布版本
- 安装第三方工具

**最佳实践**:
- 使用语义化版本
- 设置合适的 GOBIN
- 管理依赖版本
- 使用构建标签
- 优化安装性能
- 处理版本冲突

**示例**:
```bash
# 安装当前包
go install

# 安装指定包
go install ./cmd/server

# 安装多个包
go install ./cmd/server ./cmd/client

# 安装到指定目录
GOBIN=/usr/local/bin go install ./cmd/app

# 安装带版本的工具
go install golang.org/x/tools/cmd/goimports@latest

# 安装特定版本
go install golang.org/x/tools/cmd/goimports@v0.1.0

# 安装带构建标签
go install -tags=debug ./cmd/app

# 安装带链接器标志
go install -ldflags="-s -w" ./cmd/app

# 安装带编译器标志
go install -gcflags="-N -l" ./cmd/app

# 安装启用竞态检测
go install -race ./cmd/app

# 安装移除路径信息
go install -trimpath ./cmd/app

# 安装带版本信息
go install -ldflags="-X main.Version=1.0.0" ./cmd/app

# 安装优化版本
go install -gcflags="-O2" ./cmd/app

# 安装带调试信息
go install -gcflags="-N -l" -ldflags="-s -w" ./cmd/app

# 安装带内联信息
go install -gcflags="-m" ./cmd/app

# 安装带逃逸分析
go install -gcflags="-m -m" ./cmd/app

# 安装带性能分析
go install -cpuprofile=cpu.prof ./cmd/app

# 安装带内存分析
go install -memprofile=mem.prof ./cmd/app

# 安装带阻塞分析
go install -blockprofile=block.prof ./cmd/app

# 安装带互斥分析
go install -mutexprofile=mutex.prof ./cmd/app

# 安装带跟踪
go install -trace=trace.out ./cmd/app

# 安装带覆盖率的测试
go install -cover ./cmd/app

# 安装带基准测试
go install -bench=. ./cmd/app

# 安装带特定测试
go install -run=TestFunction ./cmd/app

# 安装带覆盖率分析的测试
go install -coverprofile=coverage.out ./cmd/app

# 安装带覆盖率模式的测试
go install -covermode=atomic ./cmd/app

# 安装带覆盖率函数的测试
go install -coverpkg=./... ./cmd/app

# 安装带覆盖率输出的测试
go install -coverdir=coverage ./cmd/app

# 安装带覆盖率格式的测试
go install -coverformat=html ./cmd/app

# 安装带覆盖率阈值的测试
go install -coverthreshold=80 ./cmd/app

# 安装带覆盖率排除的测试
go install -coverexclude=testdata ./cmd/app

# 安装带覆盖率包含的测试
go install -coverinclude=*.go ./cmd/app

# 安装带覆盖率排除模式的测试
go install -coverexclude=*_test.go ./cmd/app

# 安装带覆盖率包含模式的测试
go install -coverinclude=main.go ./cmd/app

# 安装常用开发工具
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/tools/cmd/godoc@latest
go install golang.org/x/tools/cmd/gorename@latest
go install golang.org/x/tools/cmd/guru@latest
go install golang.org/x/tools/cmd/godoc@latest
go install golang.org/x/tools/cmd/godoc@latest

# 安装性能分析工具
go install github.com/google/pprof@latest
go install github.com/uber/go-torch@latest
go install github.com/rakyll/hey@latest

# 安装代码质量工具
go install golang.org/x/lint/golint@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install github.com/dominikh/go-tools/cmd/staticcheck@latest

# 安装测试工具
go install github.com/axw/gocov/gocov@latest
go install github.com/axw/gocov/cmd/gocov@latest
go install github.com/axw/gocov/cmd/gocov-xml@latest
go install github.com/axw/gocov/cmd/gocov-html@latest

# 安装文档工具
go install golang.org/x/tools/cmd/godoc@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装调试工具
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/ramya-rao-a/go-outline@latest
go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest

# 安装代码生成工具
go install github.com/golang/mock/mockgen@latest
go install github.com/vektra/mockery/v2@latest
go install github.com/99designs/gqlgen@latest

# 安装数据库工具
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/jackc/pgx/v4/cmd/pgx@latest
go install github.com/go-sql-driver/mysql@latest

# 安装 Web 开发工具
go install github.com/gin-gonic/gin@latest
go install github.com/gorilla/mux@latest
go install github.com/labstack/echo/v4@latest

# 安装微服务工具
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### go test
**功能**: 测试 Go 包

**语法**: `go test [build/test flags] [packages] [build/test flags & test binary flags]`

**详细说明**:
- 编译并运行指定包的测试
- 支持单元测试、基准测试和示例测试
- 自动发现测试文件
- 并行执行测试
- 生成测试报告

**测试类型**:
- 单元测试（TestXxx）
- 基准测试（BenchmarkXxx）
- 示例测试（ExampleXxx）
- 模糊测试（FuzzXxx）
- 集成测试
- 表驱动测试

**测试文件命名**:
- `*_test.go` - 测试文件
- `*_bench_test.go` - 基准测试文件
- `*_example_test.go` - 示例测试文件
- `*_fuzz_test.go` - 模糊测试文件

**测试函数命名**:
- `TestXxx(t *testing.T)` - 单元测试
- `BenchmarkXxx(b *testing.B)` - 基准测试
- `ExampleXxx()` - 示例测试
- `FuzzXxx(f *testing.F)` - 模糊测试

**常用标志**:
- `-v` - 详细输出
- `-run` - 运行匹配的测试
- `-bench` - 运行基准测试
- `-cover` - 启用覆盖率
- `-race` - 启用竞态检测
- `-short` - 运行短测试
- `-timeout` - 设置超时时间
- `-parallel` - 设置并行度

**覆盖率选项**:
- `-coverprofile` - 生成覆盖率文件
- `-covermode` - 设置覆盖率模式
- `-coverpkg` - 指定覆盖率包
- `-coverdir` - 设置覆盖率目录
- `-coverformat` - 设置覆盖率格式
- `-coverthreshold` - 设置覆盖率阈值

**性能分析**:
- `-cpuprofile` - CPU 性能分析
- `-memprofile` - 内存性能分析
- `-blockprofile` - 阻塞性能分析
- `-mutexprofile` - 互斥性能分析
- `-trace` - 执行跟踪

**测试执行**:
- 自动发现测试文件
- 并行执行测试
- 支持测试超时
- 支持测试清理
- 支持测试钩子

**测试报告**:
- 测试结果摘要
- 失败测试详情
- 覆盖率报告
- 性能分析报告
- 基准测试结果

**使用场景**:
- 单元测试
- 集成测试
- 性能测试
- 覆盖率分析
- 代码质量检查
- 回归测试

**最佳实践**:
- 编写全面的测试
- 使用表驱动测试
- 设置合适的超时
- 使用测试钩子
- 处理测试清理
- 优化测试性能

**示例**:
```bash
# 运行当前包的所有测试
go test

# 运行指定包的测试
go test ./pkg

# 运行多个包的测试
go test ./pkg1 ./pkg2

# 详细输出
go test -v

# 运行匹配的测试
go test -run=TestFunction

# 运行基准测试
go test -bench=.

# 运行特定基准测试
go test -bench=BenchmarkFunction

# 启用覆盖率
go test -cover

# 生成覆盖率文件
go test -coverprofile=coverage.out

# 设置覆盖率模式
go test -covermode=atomic

# 指定覆盖率包
go test -coverpkg=./...

# 设置覆盖率目录
go test -coverdir=coverage

# 设置覆盖率格式
go test -coverformat=html

# 设置覆盖率阈值
go test -coverthreshold=80

# 排除覆盖率文件
go test -coverexclude=testdata

# 包含覆盖率文件
go test -coverinclude=*.go

# 排除覆盖率模式
go test -coverexclude=*_test.go

# 包含覆盖率模式
go test -coverinclude=main.go

# 启用竞态检测
go test -race

# 运行短测试
go test -short

# 设置超时时间
go test -timeout=30s

# 设置并行度
go test -parallel=4

# CPU 性能分析
go test -cpuprofile=cpu.prof

# 内存性能分析
go test -memprofile=mem.prof

# 阻塞性能分析
go test -blockprofile=block.prof

# 互斥性能分析
go test -mutexprofile=mutex.prof

# 执行跟踪
go test -trace=trace.out

# 显示执行的命令
go test -x

# 显示编译过程
go test -v -x

# 运行测试并生成报告
go test -v -cover -coverprofile=coverage.out

# 运行基准测试并生成报告
go test -bench=. -benchmem -benchtime=1s

# 运行特定测试并生成报告
go test -run=TestFunction -v -cover

# 运行所有测试并生成报告
go test ./... -v -cover -coverprofile=coverage.out

# 运行测试并检查覆盖率
go test -cover -coverprofile=coverage.out && go tool cover -func=coverage.out

# 运行测试并生成 HTML 报告
go test -cover -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# 运行测试并检查覆盖率阈值
go test -cover -coverprofile=coverage.out && go tool cover -func=coverage.out | grep total | awk '{if ($3 < 80) exit 1}'

# 运行测试并生成性能报告
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof

# 运行测试并分析性能
go test -bench=. -benchmem -cpuprofile=cpu.prof && go tool pprof cpu.prof

# 运行测试并分析内存
go test -bench=. -benchmem -memprofile=mem.prof && go tool pprof mem.prof

# 运行测试并生成火焰图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -http=:8080 cpu.prof

# 运行测试并生成调用图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -web cpu.prof

# 运行测试并生成文本报告
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -text cpu.prof

# 运行测试并生成列表报告
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -list=. cpu.prof

# 运行测试并生成汇编报告
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -disasm=. cpu.prof

# 运行测试并生成调用树
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -tree cpu.prof

# 运行测试并生成调用图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -callgrind cpu.prof

# 运行测试并生成 Graphviz 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -dot cpu.prof

# 运行测试并生成 SVG 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -svg cpu.prof

# 运行测试并生成 PDF 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -pdf cpu.prof

# 运行测试并生成 PNG 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -png cpu.prof

# 运行测试并生成 JPEG 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -jpeg cpu.prof

# 运行测试并生成 GIF 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -gif cpu.prof

# 运行测试并生成 BMP 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -bmp cpu.prof

# 运行测试并生成 TIFF 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -tiff cpu.prof

# 运行测试并生成 RAW 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -raw cpu.prof

# 运行测试并生成 XML 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -xml cpu.prof

# 运行测试并生成 JSON 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -json cpu.prof

# 运行测试并生成 YAML 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -yaml cpu.prof

# 运行测试并生成 CSV 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -csv cpu.prof

# 运行测试并生成 TSV 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -tsv cpu.prof

# 运行测试并生成 HTML 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -html cpu.prof

# 运行测试并生成 Markdown 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -markdown cpu.prof

# 运行测试并生成 LaTeX 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -latex cpu.prof

# 运行测试并生成 RTF 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -rtf cpu.prof

# 运行测试并生成 ODT 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -odt cpu.prof

# 运行测试并生成 DOCX 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -docx cpu.prof

# 运行测试并生成 XLSX 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -xlsx cpu.prof

# 运行测试并生成 PPTX 图
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -pptx cpu.prof
```

### go vet
**功能**: 检查 Go 源代码中常见的错误

**语法**: `go vet [vet flags] [packages]`

**详细说明**:
- 静态分析 Go 代码中的潜在问题
- 检查常见的编程错误和可疑构造
- 不编译代码，只进行语法和语义分析
- 提供详细的错误信息和修复建议
- 支持自定义检查器

**检查类型**:
- 未使用的变量和函数
- 错误的函数调用
- 格式字符串错误
- 结构体标签错误
- 通道使用错误
- 接口实现错误
- 类型断言错误
- 并发安全问题

**常用标志**:
- `-all` - 运行所有检查器
- `-v` - 详细输出
- `-n` - 打印命令但不执行
- `-x` - 显示执行的命令
- `-tags` - 构建标签

**检查器列表**:
- `assign` - 检查赋值语句
- `atomic` - 检查 atomic 包使用
- `bools` - 检查布尔表达式
- `buildtag` - 检查构建标签
- `cgocall` - 检查 CGO 调用
- `composite` - 检查复合字面量
- `copylocks` - 检查锁的复制
- `directive` - 检查指令
- `errorsas` - 检查 errors.As 使用
- `framepointer` - 检查帧指针
- `httpresponse` - 检查 HTTP 响应
- `ifaceassert` - 检查接口断言
- `loopclosure` - 检查循环闭包
- `lostcancel` - 检查丢失的取消
- `nilfunc` - 检查 nil 函数调用
- `printf` - 检查 printf 函数
- `shift` - 检查位移操作
- `sigchanyzer` - 检查信号通道
- `slog` - 检查 slog 包使用
- `stdmethods` - 检查标准方法
- `stringintconv` - 检查字符串整数转换
- `structtag` - 检查结构体标签
- `testinggoroutine` - 检查测试 goroutine
- `tests` - 检查测试函数
- `timeformat` - 检查时间格式
- `unmarshal` - 检查 unmarshal 函数
- `unreachable` - 检查不可达代码
- `unsafeptr` - 检查 unsafe.Pointer
- `unusedresult` - 检查未使用的结果

**使用场景**:
- 代码质量检查
- 静态分析
- 错误预防
- 代码审查
- CI/CD 集成
- 开发工具集成

**最佳实践**:
- 在提交前运行
- 集成到 CI/CD 流程
- 使用 IDE 插件
- 定期检查代码
- 关注新的检查器

**示例**:
```bash
# 检查当前包
go vet

# 检查指定包
go vet ./pkg

# 检查多个包
go vet ./pkg1 ./pkg2

# 检查所有包
go vet ./...

# 运行所有检查器
go vet -all

# 详细输出
go vet -v

# 打印命令但不执行
go vet -n

# 显示执行的命令
go vet -x

# 使用构建标签
go vet -tags=debug

# 检查特定检查器
go vet -printf ./pkg

# 检查结构体标签
go vet -structtag ./pkg

# 检查 HTTP 响应
go vet -httpresponse ./pkg

# 检查接口断言
go vet -ifaceassert ./pkg

# 检查循环闭包
go vet -loopclosure ./pkg

# 检查 printf 函数
go vet -printf ./pkg

# 检查时间格式
go vet -timeformat ./pkg

# 检查 unmarshal 函数
go vet -unmarshal ./pkg

# 检查 unsafe.Pointer
go vet -unsafeptr ./pkg

# 检查未使用的结果
go vet -unusedresult ./pkg

# 检查测试函数
go vet -tests ./pkg

# 检查测试 goroutine
go vet -testinggoroutine ./pkg

# 检查原子操作
go vet -atomic ./pkg

# 检查布尔表达式
go vet -bools ./pkg

# 检查复合字面量
go vet -composite ./pkg

# 检查锁的复制
go vet -copylocks ./pkg

# 检查 CGO 调用
go vet -cgocall ./pkg

# 检查错误处理
go vet -errorsas ./pkg

# 检查帧指针
go vet -framepointer ./pkg

# 检查指令
go vet -directive ./pkg

# 检查位移操作
go vet -shift ./pkg

# 检查信号通道
go vet -sigchanyzer ./pkg

# 检查 slog 包
go vet -slog ./pkg

# 检查标准方法
go vet -stdmethods ./pkg

# 检查字符串整数转换
go vet -stringintconv ./pkg

# 检查不可达代码
go vet -unreachable ./pkg

# 检查 nil 函数调用
go vet -nilfunc ./pkg

# 检查丢失的取消
go vet -lostcancel ./pkg

# 检查赋值语句
go vet -assign ./pkg

# 检查构建标签
go vet -buildtag ./pkg

# 集成到 Makefile
vet:
	go vet ./...

# 集成到 CI/CD
- name: Run go vet
  run: go vet ./...

# 集成到 pre-commit hook
#!/bin/sh
go vet ./...

# 集成到 IDE
# 在 VS Code 中配置 go vet 作为 linter

# 自定义检查器
go vet -vettool=$(which custom-vet) ./pkg

# 忽略特定检查器
go vet -printf=false ./pkg

# 检查特定文件
go vet file.go

# 检查特定函数
go vet -printf -printf.funcs=log.Printf ./pkg

# 检查特定包
go vet -printf -printf.packages=log ./pkg

# 检查特定函数和包
go vet -printf -printf.funcs=log.Printf -printf.packages=log ./pkg
```

### go version
**功能**: 显示 Go 版本信息

**语法**: `go version [-m] [-v] [file ...]`

**详细说明**:
- 显示当前安装的 Go 版本
- 显示编译信息和构建标签
- 可以检查二进制文件的 Go 版本
- 支持模块版本信息显示
- 提供详细的版本元数据

**输出信息**:
- Go 版本号
- 构建时间
- 构建平台
- 构建标签
- 编译器版本
- 运行时版本

**常用标志**:
- `-m` - 显示模块版本信息
- `-v` - 详细输出

**版本格式**:
- 语义化版本 (semver)
- 预发布版本
- 构建元数据
- 开发版本

**使用场景**:
- 检查 Go 版本
- 验证环境配置
- 调试版本问题
- 脚本自动化
- CI/CD 集成

**示例**:
```bash
# 显示当前 Go 版本
go version

# 显示详细版本信息
go version -v

# 检查二进制文件的 Go 版本
go version ./app

# 检查多个文件的版本
go version file1 file2

# 显示模块版本信息
go version -m ./app

# 检查远程二进制文件
go version https://example.com/app

# 在脚本中使用
VERSION=$(go version)
echo "Go version: $VERSION"

# 检查版本是否满足要求
go version | grep -q "go1.21" && echo "Version OK" || echo "Version too old"

# 提取版本号
VERSION=$(go version | awk '{print $3}')
echo "Version: $VERSION"

# 检查特定版本
if go version | grep -q "go1.21"; then
    echo "Go 1.21 or higher"
else
    echo "Need Go 1.21 or higher"
fi

# 在 Makefile 中使用
check-version:
	@go version | grep -q "go1.21" || (echo "Need Go 1.21 or higher" && exit 1)

# 在 CI/CD 中使用
- name: Check Go version
  run: |
    go version
    go version | grep -q "go1.21" || exit 1

# 检查模块依赖版本
go version -m ./cmd/app

# 显示所有模块版本
go version -m ./... | grep "mod"

# 检查特定模块版本
go version -m ./app | grep "github.com/example/module"

# 验证二进制文件兼容性
go version ./app | grep -q "go1.21" && echo "Compatible" || echo "Incompatible"

# 批量检查版本
for file in ./bin/*; do
    echo "$file: $(go version $file)"
done

# 检查 Docker 镜像中的 Go 版本
docker run --rm golang:1.21 go version

# 检查不同架构的版本
GOOS=linux GOARCH=amd64 go version
GOOS=windows GOARCH=386 go version

# 检查开发版本
go version | grep -q "devel" && echo "Development version" || echo "Release version"

# 检查预发布版本
go version | grep -q "beta\|rc\|alpha" && echo "Pre-release" || echo "Stable release"

# 提取构建信息
go version | sed 's/.*buildinfo=\([^ ]*\).*/\1/'

# 检查构建标签
go version | grep -o 'buildinfo=[^ ]*' | cut -d= -f2

# 显示编译器信息
go version | grep -o 'compiler=[^ ]*' | cut -d= -f2

# 显示平台信息
go version | grep -o 'GOOS=[^ ]*' | cut -d= -f2
go version | grep -o 'GOARCH=[^ ]*' | cut -d= -f2

# 检查是否支持特定功能
go version | grep -q "go1.21" && echo "Supports generics" || echo "No generics support"

# 版本比较脚本
#!/bin/bash
REQUIRED_VERSION="go1.21"
CURRENT_VERSION=$(go version | awk '{print $3}')

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$CURRENT_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
    echo "Go version $CURRENT_VERSION is sufficient"
else
    echo "Go version $CURRENT_VERSION is too old, need $REQUIRED_VERSION or higher"
    exit 1
fi
```

### go tool
**功能**: 运行指定的 Go 工具链工具

**语法**: `go tool [tool] [args]`

**详细说明**:
- 提供对 Go 工具链中各种工具的访问
- 包括编译器、链接器、汇编器、分析工具等
- 每个工具都有特定的功能和参数
- 支持跨平台工具链
- 提供调试和性能分析能力

**可用工具**:
- `compile` - Go 编译器
- `link` - 链接器
- `asm` - 汇编器
- `cgo` - CGO 工具
- `cover` - 代码覆盖率工具
- `pprof` - 性能分析工具
- `trace` - 执行跟踪工具
- `vet` - 代码检查工具
- `addr2line` - 地址到行号转换
- `nm` - 符号表查看器
- `objdump` - 目标文件反汇编器
- `pack` - 包管理器
- `test2json` - 测试输出转换器

**使用场景**:
- 编译和链接 Go 程序
- 性能分析和优化
- 调试和故障排除
- 代码覆盖率分析
- 内存和 CPU 分析
- 执行跟踪分析

**示例**:
```bash
# 查看所有可用工具
go tool

# 编译 Go 文件
go tool compile main.go

# 链接目标文件
go tool link -o main main.o

# 汇编汇编文件
go tool asm -o output.o input.s

# 生成代码覆盖率报告
go tool cover -html=coverage.out -o coverage.html

# 分析 CPU profile
go tool pprof cpu.prof

# 分析内存 profile
go tool pprof mem.prof

# 查看执行跟踪
go tool trace trace.out

# 反汇编二进制文件
go tool objdump -s main main

# 查看符号表
go tool nm main

# 地址到行号转换
go tool addr2line main 0x123456

# 运行 CGO 工具
go tool cgo -godefs types.go

# 分析 goroutine profile
go tool pprof goroutine.prof

# 分析阻塞 profile
go tool pprof block.prof

# 分析互斥锁 profile
go tool pprof mutex.prof

# 生成火焰图
go tool pprof -http=:8080 cpu.prof

# 查看函数调用图
go tool pprof -list=main cpu.prof

# 分析特定函数
go tool pprof -focus=MyFunction cpu.prof

# 排除特定函数
go tool pprof -ignore=MyFunction cpu.prof

# 生成 SVG 图
go tool pprof -svg cpu.prof > profile.svg

# 生成 PDF 图
go tool pprof -pdf cpu.prof > profile.pdf

# 分析堆内存分配
go tool pprof -alloc_space mem.prof

# 分析堆内存对象数量
go tool pprof -alloc_objects mem.prof

# 分析内存使用峰值
go tool pprof -inuse_space mem.prof

# 分析内存对象数量峰值
go tool pprof -inuse_objects mem.prof

# 查看汇编代码
go tool compile -S main.go

# 生成调试信息
go tool compile -N -l main.go

# 优化编译
go tool compile -O main.go

# 内联函数
go tool compile -l main.go

# 禁用内联
go tool compile -l=false main.go

# 生成位置信息
go tool compile -trimpath=false main.go

# 设置构建标签
go tool compile -tags=debug main.go

# 交叉编译
GOOS=linux GOARCH=amd64 go tool compile main.go

# 设置编译器标志
go tool compile -gcflags="-N -l" main.go

# 设置链接器标志
go tool link -ldflags="-s -w" -o main main.o

# 生成静态链接
go tool link -ldflags="-extldflags=-static" -o main main.o

# 设置入口点
go tool link -ldflags="-X main.Version=1.0.0" -o main main.o

# 生成调试信息
go tool link -ldflags="-w=false" -o main main.o

# 设置库路径
go tool link -L /path/to/libs -o main main.o

# 生成共享库
go tool link -buildmode=c-shared -o libmain.so main.o

# 生成静态库
go tool link -buildmode=c-archive -o libmain.a main.o

# 生成插件
go tool link -buildmode=plugin -o plugin.so main.o

# 分析包依赖
go tool nm -sort=size main

# 查看未定义符号
go tool nm -u main

# 查看外部符号
go tool nm -g main

# 查看动态符号
go tool nm -D main

# 反汇编特定函数
go tool objdump -s main.main main

# 反汇编特定地址范围
go tool objdump -s main -start-addr=0x1000 -stop-addr=0x2000 main

# 查看文件头信息
go tool objdump -h main

# 查看段信息
go tool objdump -s .text main

# 生成测试 JSON 输出
go test -json ./... | go tool test2json

# 分析测试结果
go test -v ./... 2>&1 | go tool test2json -t

# 生成覆盖率函数报告
go tool cover -func=coverage.out

# 生成覆盖率行报告
go tool cover -mode=count coverage.out

# 生成覆盖率原子报告
go tool cover -mode=atomic coverage.out

# 设置覆盖率模式
go test -covermode=atomic -coverprofile=coverage.out ./...

# 分析特定包的覆盖率
go tool cover -func=coverage.out | grep "mypackage"

# 生成覆盖率差异报告
go tool cover -func=coverage.out | diff -u expected.txt -

# 分析 goroutine 泄漏
go tool pprof -top goroutine.prof

# 分析内存泄漏
go tool pprof -top -alloc_space mem.prof

# 分析 CPU 热点
go tool pprof -top cpu.prof

# 分析阻塞热点
go tool pprof -top block.prof

# 分析互斥锁热点
go tool pprof -top mutex.prof

# 生成调用图
go tool pprof -weblist=main cpu.prof

# 分析特定 goroutine
go tool pprof -focus=goroutine1 cpu.prof

# 分析特定时间范围
go tool pprof -seconds=10 cpu.prof

# 生成采样报告
go tool pprof -raw cpu.prof

# 分析网络 profile
go tool pprof network.prof

# 分析系统调用 profile
go tool pprof syscall.prof

# 分析文件 I/O profile
go tool pprof fileio.prof

# 生成对比报告
go tool pprof -base=old.prof new.prof

# 分析特定函数调用
go tool pprof -focus=MyFunction -ignore=OtherFunction cpu.prof

# 生成树形图
go tool pprof -tree cpu.prof

# 生成调用图
go tool pprof -callgrind cpu.prof > callgrind.out

# 生成 Graphviz 图
go tool pprof -dot cpu.prof > profile.dot

# 分析特定包
go tool pprof -focus=mypackage cpu.prof

# 排除标准库
go tool pprof -ignore=std cpu.prof

# 分析特定方法
go tool pprof -focus=MyStruct.MyMethod cpu.prof

# 生成内存分配图
go tool pprof -weblist=main -alloc_space mem.prof

# 分析内存分配模式
go tool pprof -list=alloc mem.prof

# 生成内存使用图
go tool pprof -weblist=main -inuse_space mem.prof

# 分析内存使用模式
go tool pprof -list=inuse mem.prof

# 分析特定大小的分配
go tool pprof -focus=size:1024 mem.prof

# 生成阻塞分析图
go tool pprof -weblist=main block.prof

# 分析阻塞模式
go tool pprof -list=block block.prof

# 分析特定阻塞时间
go tool pprof -focus=duration:100ms block.prof

# 生成互斥锁分析图
go tool pprof -weblist=main mutex.prof

# 分析互斥锁模式
go tool pprof -list=mutex mutex.prof

# 分析特定锁等待时间
go tool pprof -focus=duration:1s mutex.prof

# 生成 goroutine 分析图
go tool pprof -weblist=main goroutine.prof

# 分析 goroutine 状态
go tool pprof -list=goroutine goroutine.prof

# 分析特定 goroutine 类型
go tool pprof -focus=state:running goroutine.prof

# 生成线程分析图
go tool pprof -weblist=main thread.prof

# 分析线程状态
go tool pprof -list=thread thread.prof

# 分析特定线程状态
go tool pprof -focus=state:syscall thread.prof

# 生成堆栈跟踪
go tool pprof -traces cpu.prof

# 分析特定堆栈
go tool pprof -focus=stack:main.main cpu.prof

# 生成调用路径
go tool pprof -callgrind cpu.prof

# 分析特定调用路径
go tool pprof -focus=path:main.main->helper.func cpu.prof

# 生成内存分配路径
go tool pprof -callgrind -alloc_space mem.prof

# 分析内存分配路径
go tool pprof -focus=path:main.main->alloc.func -alloc_space mem.prof

# 生成阻塞路径
go tool pprof -callgrind block.prof

# 分析阻塞路径
go tool pprof -focus=path:main.main->block.func block.prof

# 生成互斥锁路径
go tool pprof -callgrind mutex.prof

# 分析互斥锁路径
go tool pprof -focus=path:main.main->lock.func mutex.prof

# 生成 goroutine 路径
go tool pprof -callgrind goroutine.prof

# 分析 goroutine 路径
go tool pprof -focus=path:main.main->goroutine.func goroutine.prof

# 生成线程路径
go tool pprof -callgrind thread.prof

# 分析线程路径
go tool pprof -focus=path:main.main->thread.func thread.prof

# 生成网络分析图
go tool pprof -weblist=main network.prof

# 分析网络模式
go tool pprof -list=network network.prof

# 分析特定网络操作
go tool pprof -focus=op:read network.prof

# 生成系统调用分析图
go tool pprof -weblist=main syscall.prof

# 分析系统调用模式
go tool pprof -list=syscall syscall.prof

# 分析特定系统调用
go tool pprof -focus=call:read syscall.prof

# 生成文件 I/O 分析图
go tool pprof -weblist=main fileio.prof

# 分析文件 I/O 模式
go tool pprof -list=fileio fileio.prof

# 分析特定文件操作
go tool pprof -focus=op:read fileio.prof

# 生成对比分析
go tool pprof -base=baseline.prof current.prof

# 分析性能回归
go tool pprof -diff_base=regression.prof current.prof

# 生成增量分析
go tool pprof -inuse_space -base=old.prof new.prof

# 分析内存增长
go tool pprof -alloc_space -base=old.prof new.prof

# 生成 CPU 对比
go tool pprof -base=cpu_old.prof cpu_new.prof

# 分析 CPU 变化
go tool pprof -diff_base=cpu_old.prof cpu_new.prof

# 生成阻塞对比
go tool pprof -base=block_old.prof block_new.prof

# 分析阻塞变化
go tool pprof -diff_base=block_old.prof block_new.prof

# 生成互斥锁对比
go tool pprof -base=mutex_old.prof mutex_new.prof

# 分析互斥锁变化
go tool pprof -diff_base=mutex_old.prof mutex_new.prof

# 生成 goroutine 对比
go tool pprof -base=goroutine_old.prof goroutine_new.prof

# 分析 goroutine 变化
go tool pprof -diff_base=goroutine_old.prof goroutine_new.prof

# 生成线程对比
go tool pprof -base=thread_old.prof thread_new.prof

# 分析线程变化
go tool pprof -diff_base=thread_old.prof thread_new.prof

# 生成网络对比
go tool pprof -base=network_old.prof network_new.prof

# 分析网络变化
go tool pprof -diff_base=network_old.prof network_new.prof

# 生成系统调用对比
go tool pprof -base=syscall_old.prof syscall_new.prof

# 分析系统调用变化
go tool pprof -diff_base=syscall_old.prof syscall_new.prof

# 生成文件 I/O 对比
go tool pprof -base=fileio_old.prof fileio_new.prof

# 分析文件 I/O 变化
go tool pprof -diff_base=fileio_old.prof fileio_new.prof

# 生成内存分配对比
go tool pprof -alloc_space -base=mem_old.prof mem_new.prof

# 分析内存分配变化
go tool pprof -alloc_space -diff_base=mem_old.prof mem_new.prof

# 生成内存使用对比
go tool pprof -inuse_space -base=mem_old.prof mem_new.prof

# 分析内存使用变化
go tool pprof -inuse_space -diff_base=mem_old.prof mem_new.prof

# 生成内存对象对比
go tool pprof -alloc_objects -base=mem_old.prof mem_new.prof

# 分析内存对象变化
go tool pprof -alloc_objects -diff_base=mem_old.prof mem_new.prof

# 生成内存对象使用对比
go tool pprof -inuse_objects -base=mem_old.prof mem_new.prof

# 分析内存对象使用变化
go tool pprof -inuse_objects -diff_base=mem_old.prof mem_new.prof

# 生成堆栈对比
go tool pprof -traces -base=old.prof new.prof

# 分析堆栈变化
go tool pprof -traces -diff_base=old.prof new.prof

# 生成调用路径对比
go tool pprof -callgrind -base=old.prof new.prof

# 分析调用路径变化
go tool pprof -callgrind -diff_base=old.prof new.prof

# 生成树形对比
go tool pprof -tree -base=old.prof new.prof

# 分析树形变化
go tool pprof -tree -diff_base=old.prof new.prof

# 生成原始对比
go tool pprof -raw -base=old.prof new.prof

# 分析原始变化
go tool pprof -raw -diff_base=old.prof new.prof

# 生成 Graphviz 对比
go tool pprof -dot -base=old.prof new.prof

# 分析 Graphviz 变化
go tool pprof -dot -diff_base=old.prof new.prof

# 生成 SVG 对比
go tool pprof -svg -base=old.prof new.prof

# 分析 SVG 变化
go tool pprof -svg -diff_base=old.prof new.prof

# 生成 PDF 对比
go tool pprof -pdf -base=old.prof new.prof

# 分析 PDF 变化
go tool pprof -pdf -diff_base=old.prof new.prof

# 生成火焰图对比
go tool pprof -flamegraph -base=old.prof new.prof

# 分析火焰图变化
go tool pprof -flamegraph -diff_base=old.prof new.prof

# 生成 Web 界面对比
go tool pprof -http=:8080 -base=old.prof new.prof

# 分析 Web 界面变化
go tool pprof -http=:8080 -diff_base=old.prof new.prof

# 生成列表对比
go tool pprof -list=main -base=old.prof new.prof

# 分析列表变化
go tool pprof -list=main -diff_base=old.prof new.prof

# 生成 Web 列表对比
go tool pprof -weblist=main -base=old.prof new.prof

# 分析 Web 列表变化
go tool pprof -weblist=main -diff_base=old.prof new.prof

# 生成顶部对比
go tool pprof -top -base=old.prof new.prof

# 分析顶部变化
go tool pprof -top -diff_base=old.prof new.prof

# 生成树形顶部对比
go tool pprof -top -tree -base=old.prof new.prof

# 分析树形顶部变化
go tool pprof -top -tree -diff_base=old.prof new.prof

# 生成原始顶部对比
go tool pprof -top -raw -base=old.prof new.prof

# 分析原始顶部变化
go tool pprof -top -raw -diff_base=old.prof new.prof

# 生成调用图顶部对比
go tool pprof -top -callgrind -base=old.prof new.prof

# 分析调用图顶部变化
go tool pprof -top -callgrind -diff_base=old.prof new.prof

# 生成堆栈顶部对比
go tool pprof -top -traces -base=old.prof new.prof

# 分析堆栈顶部变化
go tool pprof -top -traces -diff_base=old.prof new.prof

# 生成 Graphviz 顶部对比
go tool pprof -top -dot -base=old.prof new.prof

# 分析 Graphviz 顶部变化
go tool pprof -top -dot -diff_base=old.prof new.prof

# 生成 SVG 顶部对比
go tool pprof -top -svg -base=old.prof new.prof

# 分析 SVG 顶部变化
go tool pprof -top -svg -diff_base=old.prof new.prof

# 生成 PDF 顶部对比
go tool pprof -top -pdf -base=old.prof new.prof

# 分析 PDF 顶部变化
go tool pprof -top -pdf -diff_base=old.prof new.prof

# 生成火焰图顶部对比
go tool pprof -top -flamegraph -base=old.prof new.prof

# 分析火焰图顶部变化
go tool pprof -top -flamegraph -diff_base=old.prof new.prof

# 生成 Web 界面顶部对比
go tool pprof -top -http=:8080 -base=old.prof new.prof

# 分析 Web 界面顶部变化
go tool pprof -top -http=:8080 -diff_base=old.prof new.prof

# 生成列表顶部对比
go tool pprof -top -list=main -base=old.prof new.prof

# 分析列表顶部变化
go tool pprof -top -list=main -diff_base=old.prof new.prof

# 生成 Web 列表顶部对比
go tool pprof -top -weblist=main -base=old.prof new.prof

# 分析 Web 列表顶部变化
go tool pprof -top -weblist=main -diff_base=old.prof new.prof

# 生成树形列表顶部对比
go tool pprof -top -tree -list=main -base=old.prof new.prof

# 分析树形列表顶部变化
go tool pprof -top -tree -list=main -diff_base=old.prof new.prof

# 生成原始列表顶部对比
go tool pprof -top -raw -list=main -base=old.prof new.prof

# 分析原始列表顶部变化
go tool pprof -top -raw -list=main -diff_base=old.prof new.prof

# 生成调用图列表顶部对比
go tool pprof -top -callgrind -list=main -base=old.prof new.prof

# 分析调用图列表顶部变化
go tool pprof -top -callgrind -list=main -diff_base=old.prof new.prof

# 生成堆栈列表顶部对比
go tool pprof -top -traces -list=main -base=old.prof new.prof

# 分析堆栈列表顶部变化
go tool pprof -top -traces -list=main -diff_base=old.prof new.prof

# 生成 Graphviz 列表顶部对比
go tool pprof -top -dot -list=main -base=old.prof new.prof

# 分析 Graphviz 列表顶部变化
go tool pprof -top -dot -list=main -diff_base=old.prof new.prof

# 生成 SVG 列表顶部对比
go tool pprof -top -svg -list=main -base=old.prof new.prof

# 分析 SVG 列表顶部变化
go tool pprof -top -svg -list=main -diff_base=old.prof new.prof

# 生成 PDF 列表顶部对比
go tool pprof -top -pdf -list=main -base=old.prof new.prof

# 分析 PDF 列表顶部变化
go tool pprof -top -pdf -list=main -diff_base=old.prof new.prof

# 生成火焰图列表顶部对比
go tool pprof -top -flamegraph -list=main -base=old.prof new.prof

# 分析火焰图列表顶部变化
go tool pprof -top -flamegraph -list=main -diff_base=old.prof new.prof

# 生成 Web 界面列表顶部对比
go tool pprof -top -http=:8080 -list=main -base=old.prof new.prof

# 分析 Web 界面列表顶部变化
go tool pprof -top -http=:8080 -list=main -diff_base=old.prof new.prof

# 生成 Web 列表树形顶部对比
go tool pprof -top -weblist=main -tree -base=old.prof new.prof

# 分析 Web 列表树形顶部变化
go tool pprof -top -weblist=main -tree -diff_base=old.prof new.prof

# 生成 Web 列表原始顶部对比
go tool pprof -top -weblist=main -raw -base=old.prof new.prof

# 分析 Web 列表原始顶部变化
go tool pprof -top -weblist=main -raw -diff_base=old.prof new.prof

# 生成 Web 列表调用图顶部对比
go tool pprof -top -weblist=main -callgrind -base=old.prof new.prof

# 分析 Web 列表调用图顶部变化
go tool pprof -top -weblist=main -callgrind -diff_base=old.prof new.prof

# 生成 Web 列表堆栈顶部对比
go tool pprof -top -weblist=main -traces -base=old.prof new.prof

# 分析 Web 列表堆栈顶部变化
go tool pprof -top -weblist=main -traces -diff_base=old.prof new.prof

# 生成 Web 列表 Graphviz 顶部对比
go tool pprof -top -weblist=main -dot -base=old.prof new.prof

# 分析 Web 列表 Graphviz 顶部变化
go tool pprof -top -weblist=main -dot -diff_base=old.prof new.prof

# 生成 Web 列表 SVG 顶部对比
go tool pprof -top -weblist=main -svg -base=old.prof new.prof

# 分析 Web 列表 SVG 顶部变化
go tool pprof -top -weblist=main -svg -diff_base=old.prof new.prof

# 生成 Web 列表 PDF 顶部对比
go tool pprof -top -weblist=main -pdf -base=old.prof new.prof

# 分析 Web 列表 PDF 顶部变化
go tool pprof -top -weblist=main -pdf -diff_base=old.prof new.prof

# 生成 Web 列表火焰图顶部对比
go tool pprof -top -weblist=main -flamegraph -base=old.prof new.prof

# 分析 Web 列表火焰图顶部变化
go tool pprof -top -weblist=main -flamegraph -diff_base=old.prof new.prof

# 生成 Web 界面 Web 列表顶部对比
go tool pprof -top -weblist=main -http=:8080 -base=old.prof new.prof

# 分析 Web 界面 Web 列表顶部变化
go tool pprof -top -weblist=main -http=:8080 -diff_base=old.prof new.prof
```

## 常用标志

### 构建标志 (适用于 build, clean, get, install, list, run, test)
- `-a` - 强制重新构建已经是最新的包
- `-n` - 打印命令但不执行
- `-p n` - 可以并行运行的程序数量
- `-race` - 启用数据竞争检测
- `-msan` - 启用内存清理器互操作
- `-asan` - 启用地址清理器互操作
- `-cover` - 启用代码覆盖率检测
- `-v` - 打印正在编译的包名
- `-work` - 打印临时工作目录名且退出时不删除
- `-x` - 打印命令
- `-C dir` - 在运行命令前切换到 dir
- `-mod mode` - 模块下载模式 (readonly, vendor, mod)
- `-tags tag,list` - 构建标签列表
- `-trimpath` - 从结果可执行文件中删除所有文件系统路径

### 编译器标志
- `-gcflags` - 传递给 go tool compile 的参数
- `-ldflags` - 传递给 go tool link 的参数
- `-asmflags` - 传递给 go tool asm 的参数
- `-gccgoflags` - 传递给 gccgo 编译器/链接器的参数

## 常用模式

### 1. 开发流程
```bash
# 初始化项目
go mod init myproject

# 添加依赖
go get github.com/gin-gonic/gin

# 运行程序
go run main.go

# 构建程序
go build -o myapp

# 测试
go test ./...
```

### 2. 模块管理
```bash
# 整理依赖
go mod tidy

# 查看依赖图
go mod graph

# 验证依赖
go mod verify

# 创建 vendor 目录
go mod vendor
```

### 3. 代码质量
```bash
# 格式化代码
go fmt ./...

# 检查代码问题
go vet ./...

# 运行测试并生成覆盖率报告
go test -cover ./...
```

## 环境变量

### 重要的 Go 环境变量
- `GOROOT` - Go 安装目录
- `GOPATH` - Go 工作空间目录
- `GOPROXY` - 模块代理
- `GOSUMDB` - 校验和数据库
- `GOPRIVATE` - 私有模块
- `GONOPROXY` - 不使用代理的模块
- `GONOSUMDB` - 不查询校验和数据库的模块

## 最佳实践

1. **使用 Go Modules** - 现代 Go 项目应该使用 Go modules 而不是 GOPATH
2. **定期运行 `go mod tidy`** - 保持依赖整洁
3. **使用 `go vet`** - 在提交前检查代码问题
4. **使用 `go fmt`** - 保持代码格式一致
5. **使用 `go test`** - 编写和运行测试
6. **使用 `go build` 进行构建** - 而不是直接使用编译器工具
7. **使用 `go install` 安装工具** - 而不是手动复制二进制文件

## 版本兼容性

- Go 1.11+ 支持 Go modules
- Go 1.16+ 默认启用 Go modules
- Go 1.17+ 支持工作区
- Go 1.18+ 支持泛型
- Go 1.19+ 支持软内存限制
- Go 1.20+ 支持 profile-guided optimization
- Go 1.21+ 支持结构化日志
- Go 1.22+ 支持循环变量捕获修复

## go telemetry

`go telemetry` 命令用于管理 Go 的遥测数据收集，帮助 Go 团队了解工具使用情况并改进产品。

### 常用命令
```bash
go telemetry                    # 查看当前遥测状态
go telemetry on                 # 启用遥测
go telemetry off                # 禁用遥测
go telemetry status             # 查看遥测数据
go telemetry status -d          # 查看遥测数据目录
go telemetry status -f          # 查看遥测数据文件
go telemetry status -c          # 查看遥测配置
go telemetry reset              # 重置遥测设置
go telemetry help               # 查看遥测帮助
```

### 遥测数据说明
- **收集内容**: Go 工具链使用情况、错误报告、性能数据、构建信息
- **隐私保护**: 不收集源代码、个人信息、项目内容
- **用途**: 改进 Go 工具链、识别问题、优化性能、指导开发方向
- **控制**: 用户可以随时启用/禁用，数据本地存储

### 环境变量
- `GOTELEMETRY` - 设置遥测模式 (on/off/query)
- `GOTELEMETRYDIR` - 遥测数据存储目录
- `GOTELEMETRYUPLOAD` - 控制数据上传 (on/off)

### 最佳实践
- 开发环境可以启用遥测以帮助改进 Go
- 生产环境或敏感项目建议禁用
- 定期检查遥测状态和设置
- 了解收集的数据类型和用途

