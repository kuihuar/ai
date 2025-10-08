# 多数据库一致性验证工具 - Cobra + Viper 优化版本

这是一个用Go语言编写的多数据库一致性验证工具，使用Cobra + Viper框架进行了优化重构，用于验证MySQL数据库从Azure迁移到AWS后的数据一致性。

## 🚀 主要特性

### 框架优化
- **Cobra CLI框架**: 提供强大的命令行接口，支持子命令、标志和自动补全
- **Viper配置管理**: 支持多种配置文件格式（JSON、YAML、TOML），环境变量和命令行参数
- **模块化架构**: 清晰的包结构，易于维护和扩展

### 功能特性
- 支持Azure和AWS多个数据库实例的对比验证
- 支持JSON、YAML、TOML等多种配置文件格式
- 支持环境变量配置
- 支持命令行参数覆盖
- 并行验证多个数据库对比对
- 大表分批处理，避免内存溢出
- 详细的验证报告和日志记录
- 自动配置文件生成

## 📁 项目结构

```
go-validator-optimization/
├── bin/                   # 编译的二进制文件目录
│   └── validator-optimization
├── cmd/                   # Cobra命令定义
│   ├── root.go           # 根命令
│   ├── init.go           # init命令 - 创建配置文件
│   └── validate.go       # validate命令 - 执行验证
├── configs/              # 配置文件目录
│   ├── config.yaml       # 默认配置文件
│   ├── dev.yaml          # 开发环境配置
│   ├── prod.yaml         # 生产环境配置
│   └── test.yaml         # 测试环境配置
├── examples/             # 配置示例
│   └── config.example.yaml
├── internal/             # 内部包
│   ├── config/          # 配置管理包
│   │   └── config.go
│   ├── types/           # 类型定义包
│   │   └── types.go
│   └── validator/       # 验证器核心逻辑包
│       └── validator.go
├── output/              # 输出目录
│   ├── logs/           # 日志文件
│   ├── reports/        # 报告文件
│   └── temp/           # 临时文件
├── scripts/             # 脚本目录
│   ├── dev.sh          # 开发环境脚本
│   ├── run_example.sh  # 快速启动脚本
│   ├── setup.sh        # 环境设置脚本
│   └── test.sh         # 测试脚本
├── main.go              # 程序入口
├── Makefile             # 构建管理
├── go.mod               # Go模块定义
└── README.md            # 项目文档
```

## 🛠️ 安装和构建

### 前置要求
- Go 1.19+
- MySQL数据库访问权限

### 快速开始
```bash
# 克隆或下载项目
cd go-validator-optimization

# 运行环境设置脚本（推荐）
./scripts/setup.sh

# 或者手动构建
make build
```

### 使用脚本快速启动
```bash
# 环境设置（首次使用）
./scripts/setup.sh

# 开发环境快速启动
./scripts/dev.sh

# 运行测试
./scripts/test.sh

# 快速启动示例
./scripts/run_example.sh
```

### 构建项目
```bash
# 使用Makefile构建（推荐）
make build

# 或者直接使用go build
go build -o bin/validator-optimization

# 跨平台构建
make build-all

# 清理构建产物
make clean
```

## 📖 使用方法

### 1. 创建配置文件

```bash
# 创建YAML格式配置文件（默认）
./bin/validator-optimization init

# 创建JSON格式配置文件
./bin/validator-optimization init --format json

# 创建TOML格式配置文件
./bin/validator-optimization init --format toml

# 指定输出文件名
./bin/validator-optimization init --output configs/my-config.yaml
```

### 2. 编辑配置文件

生成的配置文件示例（YAML格式）：

```yaml
azure:
  - name: azure-db1
    host: your-azure-mysql1.mysql.database.azure.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4
  - name: azure-db2
    host: your-azure-mysql2.mysql.database.azure.com
    user: your_username
    password: your_password
    database: db2
    charset: utf8mb4

aws:
  - name: aws-db1
    host: your-aws-rds1.region.rds.amazonaws.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4
  - name: aws-db2
    host: your-aws-rds2.region.rds.amazonaws.com
    user: your_username
    password: your_password
    database: db2
    charset: utf8mb4

max_workers: 3
```

### 3. 执行验证

```bash
# 使用默认配置文件验证
./bin/validator-optimization validate

# 指定配置文件
./bin/validator-optimization validate --config configs/my-config.yaml

# 设置并发数
./bin/validator-optimization validate --workers 5

# 试运行模式
./bin/validator-optimization validate --dry-run

# 指定输出文件
./bin/validator-optimization validate --output my-report.json

# 详细输出
./bin/validator-optimization validate --verbose
```

### 4. 命令行参数覆盖

```bash
# 使用命令行参数覆盖配置文件
./validator-optimization validate \
  --azure-host azure.example.com \
  --azure-user myuser \
  --azure-password mypass \
  --azure-database mydb \
  --aws-host aws.example.com \
  --aws-user myuser \
  --aws-password mypass \
  --aws-database mydb
```

## 🔧 配置方式

### 配置优先级
1. **命令行参数** (最高优先级)
2. **环境变量**
3. **配置文件**
4. **默认值** (最低优先级)

### 环境变量

```bash
# Azure配置
export MDV_AZURE_HOST="azure.example.com"
export MDV_AZURE_USER="myuser"
export MDV_AZURE_PASSWORD="mypass"
export MDV_AZURE_DATABASE="mydb"

# AWS配置
export MDV_AWS_HOST="aws.example.com"
export MDV_AWS_USER="myuser"
export MDV_AWS_PASSWORD="mypass"
export MDV_AWS_DATABASE="mydb"

# 其他配置
export MDV_MAX_WORKERS="5"
export MDV_OUTPUT="my-report.json"
```

## 📊 验证报告

验证完成后会生成详细的JSON格式报告，包含：

- 验证时间戳
- 总数据库数量
- 验证成功数量
- 数据不一致数量
- 验证错误数量
- 成功率统计
- 每个数据库的详细验证结果

### 报告示例

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "total_databases": 2,
  "successful_validations": 1,
  "inconsistent_databases": 1,
  "error_databases": 0,
  "success_rate": "50.00%",
  "results": {
    "db1": {
      "database": "db1",
      "azure_instance": "azure-db1",
      "aws_instance": "aws-db1",
      "azure_tables": 10,
      "aws_tables": 10,
      "status": "SUCCESS",
      "table_comparisons": [...]
    }
  }
}
```

## 🔧 脚本工具

### 开发脚本
项目提供了多个便捷脚本来简化开发和部署过程：

#### `scripts/setup.sh` - 环境设置脚本
```bash
./scripts/setup.sh
```
- 检查Go环境
- 安装项目依赖
- 构建项目
- 创建必要目录
- 生成默认配置文件

#### `scripts/dev.sh` - 开发环境脚本
```bash
./scripts/dev.sh
```
- 构建项目
- 运行测试
- 格式化代码
- 使用开发配置验证

#### `scripts/test.sh` - 测试脚本
```bash
./scripts/test.sh
```
- 运行单元测试
- 运行基准测试
- 生成覆盖率报告
- 竞态检测

#### `scripts/run_example.sh` - 快速启动脚本
```bash
./scripts/run_example.sh
```
- 检查环境
- 构建项目
- 创建配置文件
- 提供使用示例

### Makefile 命令
```bash
# 基本构建
make build              # 构建二进制文件
make clean              # 清理构建产物
make run                # 构建并运行

# 开发工具
make test               # 运行测试
make fmt                # 格式化代码
make deps               # 安装依赖

# 跨平台构建
make build-all          # 构建所有平台版本

# 配置管理
make init-config        # 创建默认配置
make validate-dry       # 试运行验证
make validate-dev       # 使用开发配置验证

# 信息查看
make help               # 显示帮助
make version            # 显示版本信息
```

## 🔍 命令参考

### 全局标志
- `--config string`: 配置文件路径 (默认: config.yaml)
- `--log-level string`: 日志级别 (debug, info, warn, error) (默认: info)
- `-v, --verbose`: 详细输出
- `--version`: 显示版本信息

### init 命令
- `-f, --format string`: 配置文件格式 (json, yaml, toml) (默认: yaml)
- `-o, --output string`: 输出文件名 (默认: config.yaml)

### validate 命令
- `-w, --workers int`: 最大并发数 (默认: 3)
- `-o, --output string`: 输出报告文件 (默认: consistency_report.json)
- `--dry-run`: 试运行模式，不执行实际验证
- `--azure-host string`: Azure数据库主机
- `--azure-user string`: Azure数据库用户名
- `--azure-password string`: Azure数据库密码
- `--azure-database string`: Azure数据库名称
- `--aws-host string`: AWS数据库主机
- `--aws-user string`: AWS数据库用户名
- `--aws-password string`: AWS数据库密码
- `--aws-database string`: AWS数据库名称

## 🆚 与原版本的区别

### 架构优化
- **模块化设计**: 使用internal包组织代码，提高可维护性
- **类型安全**: 统一的类型定义，避免重复代码
- **依赖注入**: 清晰的依赖关系，便于测试

### 用户体验
- **命令行界面**: 使用Cobra提供专业的CLI体验
- **配置管理**: Viper支持多种配置方式，更加灵活
- **自动补全**: 支持shell自动补全功能
- **帮助系统**: 详细的帮助信息和示例

### 功能增强
- **多格式支持**: 支持JSON、YAML、TOML配置文件
- **环境变量**: 支持环境变量配置
- **参数覆盖**: 命令行参数可以覆盖配置文件
- **配置文件生成**: 自动生成默认配置文件

## 🐛 故障排除

### 常见问题

1. **配置文件格式错误**
   ```bash
   # 检查配置文件语法
   ./validator-optimization validate --dry-run
   ```

2. **数据库连接失败**
   - 检查网络连接
   - 验证数据库凭据
   - 确认防火墙设置

3. **权限不足**
   - 确保数据库用户有SELECT权限
   - 检查information_schema访问权限

### 调试模式

```bash
# 启用详细日志
./validator-optimization validate --log-level debug --verbose

# 试运行模式
./validator-optimization validate --dry-run
```

## 📝 开发说明

### 添加新命令
1. 在`cmd/`目录下创建新的命令文件
2. 在`root.go`中注册新命令
3. 实现命令逻辑

### 扩展配置
1. 在`internal/types/types.go`中添加新的配置字段
2. 在`config.go`中添加Viper绑定
3. 在命令中添加相应的标志

### 测试
```bash
# 运行测试
go test ./...

# 测试特定包
go test ./internal/validator
```

## 📄 许可证

本项目采用MIT许可证。

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个工具。

## 📞 支持

如有问题或建议，请通过以下方式联系：
- 提交GitHub Issue
- 发送邮件至项目维护者
