# 用户同步 API 调用示例

## API 端点

- **全量同步**: `POST /api/v1/users/sync/full`
- **增量同步**: `POST /api/v1/users/sync/incremental`

服务地址：`http://localhost:8000`（根据配置文件 `configs/config.yaml` 中的 `server.http.addr` 配置）

## 1. 全量同步用户

### 请求示例

#### 同步所有部门的用户（不指定部门ID）

```bash
curl -X POST http://localhost:8000/api/v1/users/sync/full \
  -H "Content-Type: application/json" \
  -d '{}'
```

或者指定空的 `dept_id`：

```bash
curl -X POST http://localhost:8000/api/v1/users/sync/full \
  -H "Content-Type: application/json" \
  -d '{"dept_id": 0}'
```

#### 同步指定部门的用户

```bash
curl -X POST http://localhost:8000/api/v1/users/sync/full \
  -H "Content-Type: application/json" \
  -d '{"dept_id": 123456}'
```

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `dept_id` | int64 | 否 | 部门ID，不传或传 0 表示同步所有部门 |

### 响应示例

```json
{
  "created_count": 10,
  "updated_count": 5,
  "total_count": 15,
  "message": "同步完成：创建 10 个用户，更新 5 个用户，共处理 15 个用户"
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `created_count` | int32 | 创建的用户数量 |
| `updated_count` | int32 | 更新的用户数量 |
| `total_count` | int32 | 同步的总用户数量 |
| `message` | string | 同步结果消息 |

## 2. 增量同步用户

### 请求示例

#### 同步最近 24 小时的所有用户（默认）

```bash
curl -X POST http://localhost:8000/api/v1/users/sync/incremental \
  -H "Content-Type: application/json" \
  -d '{}'
```

#### 同步指定时间之后的所有用户

```bash
# 同步 2024-01-01 00:00:00 之后的所有用户
curl -X POST http://localhost:8000/api/v1/users/sync/incremental \
  -H "Content-Type: application/json" \
  -d '{"since": 1704067200}'
```

#### 同步指定部门的增量用户

```bash
curl -X POST http://localhost:8000/api/v1/users/sync/incremental \
  -H "Content-Type: application/json" \
  -d '{"dept_id": 123456, "since": 1704067200}'
```

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `dept_id` | int64 | 否 | 部门ID，不传或传 0 表示同步所有部门 |
| `since` | int64 | 否 | 同步起始时间（Unix 时间戳），不传则默认同步最近 24 小时的数据 |

### 响应示例

```json
{
  "created_count": 3,
  "updated_count": 2,
  "total_count": 5,
  "message": "增量同步完成：创建 3 个用户，更新 2 个用户，共处理 5 个用户"
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `created_count` | int32 | 创建的用户数量 |
| `updated_count` | int32 | 更新的用户数量 |
| `total_count` | int32 | 同步的总用户数量 |
| `message` | string | 同步结果消息 |

## 3. 使用其他工具调用

### 使用 HTTPie

```bash
# 全量同步
http POST localhost:8000/api/v1/users/sync/full dept_id:=0

# 增量同步
http POST localhost:8000/api/v1/users/sync/incremental dept_id:=0 since:=1704067200
```

### 使用 Postman

1. 方法：`POST`
2. URL：`http://localhost:8000/api/v1/users/sync/full` 或 `/api/v1/users/sync/incremental`
3. Headers：
   - `Content-Type: application/json`
4. Body（raw JSON）：
   ```json
   {
     "dept_id": 0,
     "since": 1704067200  // 仅增量同步需要
   }
   ```

### 使用 Go 代码调用

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "sre/api/user/v1"
    "github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
    // 创建 HTTP 客户端
    conn, err := http.NewClient(
        context.Background(),
        http.WithEndpoint("http://localhost:8000"),
    )
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    client := v1.NewUserHTTPClient(conn)
    ctx := context.Background()
    
    // 全量同步
    fullResp, err := client.SyncUsersFull(ctx, &v1.SyncUsersFullRequest{
        DeptId: 0, // 0 表示同步所有部门
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("全量同步结果: %+v\n", fullResp)
    
    // 增量同步（最近 24 小时）
    incResp, err := client.SyncUsersIncremental(ctx, &v1.SyncUsersIncrementalRequest{
        DeptId: 0,
        Since:  time.Now().Add(-24 * time.Hour).Unix(),
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("增量同步结果: %+v\n", incResp)
}
```

### 使用 Python 调用

```python
import requests
import time

BASE_URL = "http://localhost:8000"

# 全量同步
def sync_full(dept_id=0):
    url = f"{BASE_URL}/api/v1/users/sync/full"
    data = {"dept_id": dept_id}
    response = requests.post(url, json=data)
    return response.json()

# 增量同步
def sync_incremental(dept_id=0, since=None):
    url = f"{BASE_URL}/api/v1/users/sync/incremental"
    data = {"dept_id": dept_id}
    if since:
        data["since"] = int(since)
    response = requests.post(url, json=data)
    return response.json()

# 使用示例
if __name__ == "__main__":
    # 全量同步所有部门
    result = sync_full()
    print("全量同步结果:", result)
    
    # 增量同步最近 24 小时
    yesterday = time.time() - 24 * 3600
    result = sync_incremental(since=yesterday)
    print("增量同步结果:", result)
```

## 4. 注意事项

1. **部门ID**: 
   - 钉钉的根部门ID通常是 `1`
   - 传 `0` 或不传 `dept_id` 表示同步所有部门

2. **时间戳格式**:
   - `since` 参数使用 Unix 时间戳（秒级）
   - 例如：`1704067200` 表示 `2024-01-01 00:00:00 UTC`

3. **增量同步默认时间**:
   - 如果不传 `since` 参数，默认同步最近 24 小时的数据

4. **错误处理**:
   - 如果钉钉客户端未配置，会返回错误
   - 如果外部服务不可用，会返回相应的错误信息

5. **性能考虑**:
   - 全量同步可能耗时较长，建议在低峰期执行
   - 增量同步通常较快，适合定期执行

## 5. 定时任务示例

### 使用 cron 定期执行增量同步

```bash
# 编辑 crontab
crontab -e

# 每天凌晨 2 点执行增量同步
0 2 * * * curl -X POST http://localhost:8000/api/v1/users/sync/incremental -H "Content-Type: application/json" -d '{}'
```

### 使用 systemd timer

创建 `/etc/systemd/system/user-sync.timer`:

```ini
[Unit]
Description=User Incremental Sync Timer

[Timer]
OnCalendar=daily
OnCalendar=02:00
Persistent=true

[Install]
WantedBy=timers.target
```

创建 `/etc/systemd/system/user-sync.service`:

```ini
[Unit]
Description=User Incremental Sync Service

[Service]
Type=oneshot
ExecStart=/usr/bin/curl -X POST http://localhost:8000/api/v1/users/sync/incremental -H "Content-Type: application/json" -d '{}'
```

启用定时器：

```bash
sudo systemctl enable user-sync.timer
sudo systemctl start user-sync.timer
```

