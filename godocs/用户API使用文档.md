# 用户API使用文档

## 概述

本文档说明用户管理相关的API接口，包括用户的增删改查功能。

## API列表

### 1. 创建用户

**接口地址**: `POST /user`

**请求参数**:
```json
{
  "name": "张三",
  "email": "zhangsan@example.com",
  "phone": "13800138000"
}
```

**参数说明**:
- `name` (必填): 用户名，长度2-50字符
- `email` (必填): 邮箱地址，需符合邮箱格式
- `phone` (可选): 手机号，需符合手机号格式

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

**cURL示例**:
```bash
curl -X POST http://localhost:8000/user \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000"
  }'
```

---

### 2. 根据ID获取用户

**接口地址**: `GET /user/:id`

**路径参数**:
- `id` (必填): 用户ID，必须大于0

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "createdAt": "2024-01-01 10:00:00",
    "updatedAt": "2024-01-01 10:00:00"
  }
}
```

**cURL示例**:
```bash
curl http://localhost:8000/user/1
```

---

### 3. 获取用户列表

**接口地址**: `GET /users`

**查询参数**:
- `page` (可选): 页码，默认1，必须大于0
- `pageSize` (可选): 每页数量，默认10，范围1-100

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "张三",
        "email": "zhangsan@example.com",
        "phone": "13800138000",
        "createdAt": "2024-01-01 10:00:00",
        "updatedAt": "2024-01-01 10:00:00"
      },
      {
        "id": 2,
        "name": "李四",
        "email": "lisi@example.com",
        "phone": "13900139000",
        "createdAt": "2024-01-02 10:00:00",
        "updatedAt": "2024-01-02 10:00:00"
      }
    ],
    "total": 2,
    "page": 1,
    "size": 10
  }
}
```

**cURL示例**:
```bash
# 获取第一页，每页10条
curl http://localhost:8000/users?page=1&pageSize=10

# 获取第二页，每页20条
curl http://localhost:8000/users?page=2&pageSize=20
```

---

### 4. 更新用户

**接口地址**: `PUT /user/:id`

**路径参数**:
- `id` (必填): 用户ID，必须大于0

**请求参数** (所有字段可选，只更新提供的字段):
```json
{
  "name": "张三（已更新）",
  "email": "zhangsan_new@example.com",
  "phone": "13800138001"
}
```

**参数说明**:
- `name` (可选): 用户名，长度2-50字符
- `email` (可选): 邮箱地址，需符合邮箱格式
- `phone` (可选): 手机号，需符合手机号格式

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true
  }
}
```

**cURL示例**:
```bash
curl -X PUT http://localhost:8000/user/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三（已更新）",
    "email": "zhangsan_new@example.com"
  }'
```

---

### 5. 删除用户

**接口地址**: `DELETE /user/:id`

**路径参数**:
- `id` (必填): 用户ID，必须大于0

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true
  }
}
```

**cURL示例**:
```bash
curl -X DELETE http://localhost:8000/user/1
```

---

## 错误响应格式

当请求失败时，响应格式如下：

```json
{
  "code": 1,
  "message": "错误信息描述",
  "data": null
}
```

### 常见错误码

- `code: 1`: 参数验证失败
- `code: 1`: 数据库操作失败
- `code: 1`: 用户不存在（查询/更新/删除时）

---

## 数据表结构

假设 `users` 表结构如下：

```sql
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `name` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

**注意**: 如果您的表结构不同，请相应调整 `internal/model/entity/user.go` 和 `internal/model/do/user.go` 中的字段定义。

---

## 测试步骤

1. **启动服务**:
   ```bash
   go run main.go
   ```

2. **创建用户**:
   ```bash
   curl -X POST http://localhost:8000/user \
     -H "Content-Type: application/json" \
     -d '{"name":"测试用户","email":"test@example.com","phone":"13800138000"}'
   ```

3. **查询用户列表**:
   ```bash
   curl http://localhost:8000/users
   ```

4. **根据ID查询用户**:
   ```bash
   curl http://localhost:8000/user/1
   ```

5. **更新用户**:
   ```bash
   curl -X PUT http://localhost:8000/user/1 \
     -H "Content-Type: application/json" \
     -d '{"name":"更新后的名字"}'
   ```

6. **删除用户**:
   ```bash
   curl -X DELETE http://localhost:8000/user/1
   ```

---

## Swagger文档

启动服务后，可以访问 Swagger UI 查看完整的API文档：

- Swagger UI: http://localhost:8000/swagger
- OpenAPI规范: http://localhost:8000/api.json

---

## 注意事项

1. 所有接口都需要数据库连接正常
2. 确保 `users` 表已创建
3. 邮箱字段建议在数据库层面设置唯一索引
4. 分页查询默认按ID倒序排列
5. 更新操作只更新提供的非空字段

