# API 设计规范

## RESTful API 设计

### URL 设计原则

1. **使用名词，避免动词**
   - ✅ `/api/v1/users`
   - ❌ `/api/v1/getUsers`

2. **使用复数形式**
   - ✅ `/api/v1/users`
   - ❌ `/api/v1/user`

3. **层次化资源**
   - ✅ `/api/v1/users/123/posts`
   - ❌ `/api/v1/user-posts?userId=123`

### HTTP 方法

- `GET`：获取资源
- `POST`：创建资源
- `PUT`：完整更新资源
- `PATCH`：部分更新资源
- `DELETE`：删除资源

### 状态码

- `200 OK`：成功
- `201 Created`：创建成功
- `204 No Content`：删除成功
- `400 Bad Request`：客户端错误
- `401 Unauthorized`：未认证
- `403 Forbidden`：无权限
- `404 Not Found`：资源不存在
- `500 Internal Server Error`：服务器错误

## gRPC API 设计

### 服务定义

```protobuf
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);
}
```

### 消息定义

```protobuf
message GetUserRequest {
  int64 id = 1;
}

message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
}
```

### 设计原则

1. **使用有意义的服务名和方法名**
2. **请求和响应消息分离**
3. **使用标准类型（如 `google.protobuf.Empty`）**
4. **字段编号从 1 开始，避免使用已废弃的编号**

## 版本管理

### URL 版本控制
```
/api/v1/users
/api/v2/users
```

### 向后兼容
- 添加新字段时使用可选字段
- 不删除已存在的字段
- 不改变字段类型

## 错误处理

### 统一错误格式

```json
{
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User with ID 123 not found",
    "details": []
  }
}
```

### 错误码设计
- 使用有意义的错误码
- 错误码应该稳定，不轻易改变
- 提供详细的错误信息

## 最佳实践

1. **幂等性**：确保 PUT、DELETE 操作是幂等的
2. **分页**：列表接口支持分页
3. **过滤和排序**：提供灵活的查询参数
4. **限流**：实现 API 限流保护
5. **文档**：提供完整的 API 文档

