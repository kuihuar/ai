# SE-Client 实现文档

## 概述

`se-client` 是一个命令行工具，用于从 CSV 或 Excel 文件导入组织架构数据（部门、用户、部门用户关系）到数据库。该工具支持批量导入，并自动识别文件类型和数据结构。

## 项目结构

```
cmd/se-client/
├── main.go          # 应用入口，命令行参数解析
├── parser.go        # 文件解析器实现（CSV/Excel）
├── wire.go          # Wire 依赖注入定义
└── wire_gen.go      # Wire 自动生成的代码
```

## 核心组件

### 1. 文件解析器 (parser.go)

#### 位置说明

`parser.go` 当前位于 `cmd/se-client/` 目录下。如果未来需要复用解析器功能，可以考虑以下位置：

**选项 1：保持当前位置** ✅ **推荐（当前）**
- `cmd/se-client/parser.go`
- **优点**：简单直接，与命令紧密耦合
- **适用场景**：解析器仅用于 se-client 命令

**选项 2：独立包**
- `internal/parser/parser.go`
- **优点**：可被多个命令复用（如 se-client、其他导入工具）
- **适用场景**：解析器功能需要被多个应用共享

**选项 3：子包**
- `cmd/se-client/parser/parser.go`
- **优点**：模块化，保持命令目录整洁
- **适用场景**：解析器代码较多，需要拆分为多个文件

#### 接口设计

```go
// FileParser 文件解析器接口
type FileParser interface {
    Parse(filePath string, taskID, thirdCompanyID, platformID string) (*ParsedData, error)
}
```

#### 支持的格式

- **CSV**：通过 `csvParser` 实现
- **Excel**：通过 `excelParser` 实现（支持 .xlsx 和 .xls）

#### 自动识别机制

解析器会根据以下规则自动识别数据类型：

1. **文件名识别**：
   - 包含 "department" → 部门数据
   - 包含 "user" → 用户数据
   - 包含 "relation" 或 "department_user" → 部门用户关系

2. **列名识别**：
   - 包含 `did`, `name`, `pid` → 部门数据
   - 包含 `uid`, `account`, `nick_name` → 用户数据
   - 包含 `uid`, `did` → 部门用户关系

3. **Excel 工作表识别**：
   - 根据工作表名称判断数据类型

### 2. 数据模型

解析器支持三种数据类型的导入：

#### 部门数据 (TbLasDepartment)

**必填字段**：
- `did`: 部门ID
- `name`: 部门名称

**可选字段**：
- `pid`: 父部门ID
- `order`: 排序
- `source`: 来源（默认：sync）
- `type`: 类型
- `check_type`: 检查类型

#### 用户数据 (TbLasUser)

**必填字段**：
- `uid`: 用户ID
- `account`: 账号
- `nick_name` 或 `nickname`: 昵称

**可选字段**：
- `def_did`: 默认部门ID
- `def_did_order`: 默认部门排序
- `password`: 密码
- `avatar`: 头像
- `email`: 邮箱
- `gender`: 性别
- `title`: 职位
- `work_place`: 工作地点
- `leader`: 领导
- `employer`: 雇主
- `employment_status`: 雇佣状态（默认：notactive）
- `employment_type`: 雇佣类型
- `phone`: 手机号
- `telephone`: 电话
- `source`: 来源（默认：sync）
- `custom_fields`: 自定义字段
- `check_type`: 检查类型

#### 部门用户关系 (TbLasDepartmentUser)

**必填字段**：
- `uid`: 用户ID
- `did`: 部门ID

**可选字段**：
- `order`: 排序
- `main`: 是否主部门（默认：0）
- `check_type`: 检查类型

### 3. 依赖注入

使用 Wire 进行依赖注入，定义在 `wire.go` 中：

```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "sre/internal/data"
)

func wireApp(server *conf.Server, confData *conf.Data, logger log.Logger) (*AppDependencies, func(), error) {
    panic(wire.Build(
        data.ProviderSet,
        provideAppDependencies,
    ))
}
```

注入的依赖包括：
- `LasDepartmentRepo`: 部门数据仓库
- `LasUserRepo`: 用户数据仓库
- `LasDepartmentUserRepo`: 部门用户关系数据仓库

## 使用方法

### 命令行参数

```bash
./se-client \
  -conf ../../configs \
  -file data.csv \
  -task-id task_001 \
  -third-company-id company_001 \
  -platform-id dingtalk
```

**参数说明**：
- `-conf`: 配置文件路径（默认：../../configs）
- `-file`: 要导入的文件路径（CSV 或 Excel）
- `-task-id`: 任务ID（可选，未指定时自动生成）
- `-third-company-id`: 租户ID（必填）
- `-platform-id`: 平台ID（必填，如：dingtalk）

### 文件格式示例

#### CSV 格式 - 部门数据

```csv
did,name,pid,order,type
dept_001,技术部,,1,tech
dept_002,研发组,dept_001,1,tech
```

#### CSV 格式 - 用户数据

```csv
uid,account,nick_name,email,def_did
user_001,zhangsan,张三,zhangsan@example.com,dept_001
user_002,lisi,李四,lisi@example.com,dept_002
```

#### CSV 格式 - 部门用户关系

```csv
uid,did,order,main
user_001,dept_001,1,1
user_002,dept_002,1,1
```

#### Excel 格式

Excel 文件可以包含多个工作表，每个工作表对应一种数据类型：
- `Department` 或包含 "department" 的工作表 → 部门数据
- `User` 或包含 "user" 的工作表 → 用户数据
- `Relation` 或包含 "relation" 的工作表 → 部门用户关系

## 工作流程

1. **参数解析**：解析命令行参数
2. **配置加载**：使用 Viper 加载配置文件
3. **日志初始化**：初始化 Zap Logger
4. **依赖注入**：使用 Wire 注入数据仓库依赖
5. **文件解析**：
   - 根据文件扩展名选择解析器（CSV 或 Excel）
   - 解析文件内容，识别数据类型
   - 转换为数据模型对象
6. **数据保存**：
   - 批量保存部门数据
   - 批量保存用户数据
   - 批量保存部门用户关系数据
7. **完成**：输出处理结果

## 错误处理

- **文件不存在**：返回明确的错误信息
- **文件格式不支持**：提示支持的格式（.csv, .xlsx, .xls）
- **必填参数缺失**：提示缺少的参数
- **解析错误**：记录警告日志，跳过错误行，继续处理
- **保存失败**：返回详细错误信息

## 日志记录

工具使用 Zap Logger 记录以下信息：
- 文件解析进度
- 解析结果统计（部门数、用户数、关系数）
- 保存进度
- 错误和警告

## 扩展性

### 添加新的文件格式支持

1. 实现 `FileParser` 接口
2. 在 `handleFileImport` 中添加格式判断逻辑

```go
case ".json":
    parser = NewJSONParser(logger)
```

### 添加新的数据类型

1. 在 `ParsedData` 中添加新字段
2. 在解析器中添加解析逻辑
3. 在 `handleFileImport` 中添加保存逻辑

## 最佳实践

1. **文件命名**：使用描述性文件名，包含数据类型（如：`departments.csv`）
2. **数据验证**：在导入前验证数据完整性
3. **批量处理**：使用批量保存接口提高性能
4. **错误恢复**：部分数据失败不影响整体导入
5. **日志记录**：记录详细的处理日志便于排查问题

## 相关文档

- [项目结构文档](../project/structure.md)
- [数据层文档](../../internal/data/README.md)
- [LAS 全量同步使用文档](../../internal/data/docs/las-full-sync-usage.md)

