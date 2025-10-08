# Go版本YAML配置文件支持更新总结

## 更新概述

为Go语言版本的多数据库一致性验证工具添加了YAML配置文件支持，现在支持JSON和YAML两种配置文件格式。

## 更新内容

### 1. 依赖添加

添加了YAML解析库依赖：
```bash
go get gopkg.in/yaml.v3
```

### 2. 代码更新

#### 2.1 config.go 更新
- 添加了YAML解析支持
- 根据文件扩展名自动选择解析方式
- 支持创建YAML格式的配置文件

```go
// 根据文件扩展名选择解析方式
ext := strings.ToLower(filepath.Ext(filename))
switch ext {
case ".yaml", ".yml":
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析YAML配置文件失败: %v", err)
    }
case ".json":
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析JSON配置文件失败: %v", err)
    }
default:
    // 默认尝试JSON格式
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析配置文件失败，请确保文件格式为JSON或YAML: %v", err)
    }
}
```

#### 2.2 types.go 更新
为所有结构体添加了YAML标签：

```go
// Config 配置文件结构
type Config struct {
    Azure      []DatabaseInstance `json:"azure" yaml:"azure"`             // Azure实例列表
    AWS        []DatabaseInstance `json:"aws" yaml:"aws"`                 // AWS实例列表
    MaxWorkers int                `json:"max_workers" yaml:"max_workers"` // 最大并发数
}

// DatabaseInstance 数据库实例配置
type DatabaseInstance struct {
    Name     string `json:"name" yaml:"name"`         // 实例名称
    Host     string `json:"host" yaml:"host"`         // 实例主机地址
    User     string `json:"user" yaml:"user"`         // 用户名
    Password string `json:"password" yaml:"password"` // 密码
    Database string `json:"database" yaml:"database"` // 数据库名称
    Charset  string `json:"charset" yaml:"charset"`   // 字符集
}
```

#### 2.3 main.go 更新
- 更新了帮助信息，显示YAML支持
- 修改了配置文件加载逻辑，支持多种格式
- 配置文件优先级：config.yaml > config.yml > config.json

```go
// 尝试从配置文件加载
configFiles := []string{"config.yaml", "config.yml", "config.json"}
var configFile string

for _, file := range configFiles {
    if _, err := os.Stat(file); err == nil {
        configFile = file
        break
    }
}
```

### 3. 文件结构更新

新增文件：
- `config.yaml.example` - YAML格式配置文件示例

更新文件：
- `README.md` - 添加YAML配置说明
- `run_example.sh` - 支持创建YAML配置文件

## 支持的配置文件格式

### JSON格式 (config.json)
```json
{
  "azure": [
    {
      "name": "azure-db1",
      "host": "your-azure-mysql1.mysql.database.azure.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db1",
      "charset": "utf8mb4"
    }
  ],
  "aws": [
    {
      "name": "aws-db1",
      "host": "your-aws-rds1.region.rds.amazonaws.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db1",
      "charset": "utf8mb4"
    }
  ],
  "max_workers": 3
}
```

### YAML格式 (config.yaml)
```yaml
# Azure数据库实例列表
azure:
  - name: azure-db1
    host: your-azure-mysql1.mysql.database.azure.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4

# AWS数据库实例列表
aws:
  - name: aws-db1
    host: your-aws-rds1.region.rds.amazonaws.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4

# 最大并发数
max_workers: 3
```

## 使用方法

### 1. 创建配置文件

```bash
# 创建JSON格式配置文件
./validator init config.json

# 创建YAML格式配置文件
./validator init config.yaml
```

### 2. 运行验证

```bash
# 自动检测配置文件（优先级：config.yaml > config.yml > config.json）
./validator

# 指定配置文件
./validator --config config.yaml
./validator --config config.json
```

### 3. 查看帮助

```bash
./validator help
```

## 配置文件优先级

程序会按以下顺序查找配置文件：
1. `config.yaml`
2. `config.yml`
3. `config.json`

如果找到配置文件，会使用该文件；如果都没有找到，会使用默认配置。

## 优势

1. **格式灵活性**: 支持JSON和YAML两种格式，用户可以根据喜好选择
2. **自动检测**: 根据文件扩展名自动选择解析方式
3. **向后兼容**: 完全兼容现有的JSON配置文件
4. **易于阅读**: YAML格式更易读，特别适合复杂配置
5. **注释支持**: YAML格式支持注释，便于配置说明

## 测试验证

已通过以下测试：
- ✅ YAML配置文件创建
- ✅ YAML配置文件加载
- ✅ JSON配置文件兼容性
- ✅ 配置文件优先级
- ✅ 帮助信息显示
- ✅ 运行示例脚本

## 兼容性

- **Go版本**: 1.16+
- **依赖**: gopkg.in/yaml.v3 v3.0.1
- **向后兼容**: 完全兼容现有JSON配置文件
- **跨平台**: 支持Linux、macOS、Windows

## 总结

成功为Go版本添加了YAML配置文件支持，提供了更灵活的配置方式。用户现在可以根据需要选择JSON或YAML格式的配置文件，程序会自动检测和解析。这大大提升了工具的易用性和灵活性。
