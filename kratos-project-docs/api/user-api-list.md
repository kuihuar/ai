# 用户 API 列表

本文档列出了所有用户相关的 API 接口。

## API 概览

| 序号 | API 名称 | gRPC 方法 | HTTP 方法 | HTTP 路径 | 功能描述 |
|------|---------|-----------|-----------|-----------|----------|
| 1 | 创建用户 | `CreateUser` | POST | `/api/v1/users` | 创建新用户 |
| 2 | 获取用户 | `GetUser` | GET | `/api/v1/users/{id}` | 根据ID获取用户信息 |
| 3 | 更新用户 | `UpdateUser` | PUT | `/api/v1/users/{id}` | 更新用户信息 |
| 4 | 删除用户 | `DeleteUser` | DELETE | `/api/v1/users/{id}` | 删除用户 |
| 5 | 列出用户 | `ListUsers` | GET | `/api/v1/users` | 分页列出用户，支持搜索 |
| 6 | 全量同步用户 | `SyncUsersFull` | POST | `/api/v1/users/sync/full` | 全量同步用户数据 |
| 7 | 增量同步用户 | `SyncUsersIncremental` | POST | `/api/v1/users/sync/incremental` | 增量同步用户数据 |
| 8 | 获取访问令牌 | `GetAccessToken` | GET/POST | `/api/v1/oauth/access_token` | 获取钉钉 OAuth 访问令牌 |
| 9 | 获取用户信息（OAuth） | `GetUserInfo` | GET | `/api/v1/oauth/userinfo` | 根据访问令牌获取用户信息 |
| 10 | 触发 WPS 账号同步 | `SyncWpsAccounts` | POST | `/api/v1/wps/accounts/sync` | 触发 WPS 全量账号同步 |

## 详细说明

### 1. 创建用户

**接口信息**
- **gRPC**: `user.v1.User.CreateUser`
- **HTTP**: `POST /api/v1/users`

**请求参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| email | string | 是 | 邮箱 |
| password | string | 是 | 密码 |
| nickname | string | 否 | 昵称 |
| avatar | string | 否 | 头像URL |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| user | UserInfo | 用户信息对象 |

---

### 2. 获取用户

**接口信息**
- **gRPC**: `user.v1.User.GetUser`
- **HTTP**: `GET /api/v1/users/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 用户ID |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| user | UserInfo | 用户信息对象 |

---

### 3. 更新用户

**接口信息**
- **gRPC**: `user.v1.User.UpdateUser`
- **HTTP**: `PUT /api/v1/users/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 用户ID |
| nickname | string | 否 | 请求体 | 昵称 |
| avatar | string | 否 | 请求体 | 头像URL |
| email | string | 否 | 请求体 | 邮箱 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| user | UserInfo | 用户信息对象 |

---

### 4. 删除用户

**接口信息**
- **gRPC**: `user.v1.User.DeleteUser`
- **HTTP**: `DELETE /api/v1/users/{id}`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| id | int64 | 是 | 路径参数 | 用户ID |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | bool | 是否成功 |

---

### 5. 列出用户

**接口信息**
- **gRPC**: `user.v1.User.ListUsers`
- **HTTP**: `GET /api/v1/users`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 默认值 | 说明 |
|--------|------|------|------|--------|------|
| page | int32 | 否 | 查询参数 | 1 | 页码（从1开始） |
| page_size | int32 | 否 | 查询参数 | 10 | 每页数量（最大100） |
| keyword | string | 否 | 查询参数 | - | 搜索关键词（用户名或邮箱） |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| users | UserInfo[] | 用户列表 |
| total | int32 | 总数量 |
| page | int32 | 当前页码 |
| page_size | int32 | 每页数量 |

---

### 6. 全量同步用户

**接口信息**
- **gRPC**: `user.v1.User.SyncUsersFull`
- **HTTP**: `POST /api/v1/users/sync/full`

**请求参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| dept_id | int64 | 否 | 部门ID，不传则同步所有部门 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| created_count | int32 | 创建的用户数量 |
| updated_count | int32 | 更新的用户数量 |
| total_count | int32 | 同步的总用户数量 |
| dept_created_count | int32 | 创建的部门数量 |
| dept_updated_count | int32 | 更新的部门数量 |
| dept_total_count | int32 | 同步的总部门数量 |
| relation_created_count | int32 | 创建的部门用户关系数量 |
| relation_updated_count | int32 | 更新的部门用户关系数量 |
| relation_total_count | int32 | 同步的总关系数量 |
| message | string | 同步结果消息 |

---

### 7. 增量同步用户

**接口信息**
- **gRPC**: `user.v1.User.SyncUsersIncremental`
- **HTTP**: `POST /api/v1/users/sync/incremental`

**请求参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| dept_id | int64 | 否 | 部门ID，不传则同步所有部门 |
| since | int64 | 否 | 同步起始时间（Unix时间戳） |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| created_count | int32 | 创建的用户数量 |
| updated_count | int32 | 更新的用户数量 |
| total_count | int32 | 同步的总用户数量 |
| message | string | 同步结果消息 |

---

### 8. 获取访问令牌（钉钉 OAuth）

**接口信息**
- **gRPC**: `user.v1.User.GetAccessToken`
- **HTTP**: 
  - `GET /api/v1/oauth/access_token`
  - `POST /api/v1/oauth/access_token`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| code | string | 是 | 查询参数/请求体 | 授权码，从钉钉 OAuth 回调中获取 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| access_token | string | 访问令牌 |
| refresh_token | string | 刷新令牌（可选） |
| expires_in | int64 | 过期时间（秒） |
| scope | string | 授权范围（可选） |

---

### 9. 获取用户信息（钉钉 OAuth）

**接口信息**
- **gRPC**: `user.v1.User.GetUserInfo`
- **HTTP**: `GET /api/v1/oauth/userinfo`

**请求参数**

| 参数名 | 类型 | 必填 | 位置 | 说明 |
|--------|------|------|------|------|
| access_token | string | 是 | 查询参数/Header | 访问令牌 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| union_id | string | 用户 union_id |
| user_id | string | 用户 user_id |
| name | string | 用户姓名 |
| email | string | 用户邮箱（可选） |
| avatar | string | 用户头像 URL（可选） |
| mobile | string | 用户手机号（可选） |
| position | string | 职位（可选） |

---

### 10. 触发 WPS 全量账号同步

**接口信息**
- **gRPC**: `user.v1.User.SyncWpsAccounts`
- **HTTP**: `POST /api/v1/wps/accounts/sync`

**请求参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| task_id | string | 否 | 任务ID，不传则自动生成 |
| third_company_id | string | 否 | 第三方公司ID，不传则使用配置中的值 |

**响应数据**

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | bool | 是否成功 |
| message | string | 响应消息 |
| task_id | string | 任务ID |
| third_company_id | string | 第三方公司ID |

---

## 数据模型

### UserInfo

用户信息对象，用于返回用户的基本信息。

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 用户ID |
| username | string | 用户名 |
| email | string | 邮箱 |
| nickname | string | 昵称 |
| avatar | string | 头像URL |
| created_at | int64 | 创建时间（Unix时间戳） |
| updated_at | int64 | 更新时间（Unix时间戳） |

---

## 错误码

所有 API 的错误码定义在 `api/user/v1/error_reason.proto` 中。

常见错误码：

| 错误码 | 说明 |
|--------|------|
| USER_NOT_FOUND | 用户不存在 |
| USER_ALREADY_EXISTS | 用户已存在 |
| INVALID_USERNAME | 无效的用户名 |
| INVALID_EMAIL | 无效的邮箱 |
| INVALID_PASSWORD | 无效的密码 |

---

## 使用示例

### HTTP 请求示例

#### 创建用户
```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "nickname": "测试用户"
  }'
```

#### 获取用户
```bash
curl http://localhost:8000/api/v1/users/1
```

#### 列出用户
```bash
curl "http://localhost:8000/api/v1/users?page=1&page_size=10&keyword=test"
```

### gRPC 请求示例

```go
import (
    "context"
    "sre/api/user/v1"
    "google.golang.org/grpc"
)

conn, _ := grpc.Dial("localhost:9000", grpc.WithInsecure())
client := v1.NewUserClient(conn)

// 创建用户
user, err := client.CreateUser(context.Background(), &v1.CreateUserRequest{
    Username: "testuser",
    Email:    "test@example.com",
    Password: "password123",
    Nickname: "测试用户",
})
```

---

## 相关文档

- [API 设计规范](../../code-standards/api-design.md)
- [用户同步示例](./user-sync-examples.md)
- [第三方 API 定义](../../architecture/third-party-api-definitions.md)
