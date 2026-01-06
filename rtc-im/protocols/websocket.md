# WebSocket 协议详解

WebSocket 协议原理、实现和应用。

## 目录

- [协议概述](#协议概述)
- [握手过程](#握手过程)
- [数据帧格式](#数据帧格式)
- [实现示例](#实现示例)
- [最佳实践](#最佳实践)

---

## 协议概述

### 什么是 WebSocket

WebSocket 是一种在单个 TCP 连接上进行全双工通信的协议。

**特点：**
- 全双工通信
- 低延迟
- 持久连接
- 支持扩展

### 与 HTTP 对比

| 特性 | HTTP | WebSocket |
|------|------|-----------|
| **连接** | 请求-响应 | 持久连接 |
| **通信** | 单向 | 全双工 |
| **开销** | 每次请求头 | 首次握手后无头 |
| **适用** | 请求-响应 | 实时通信 |

---

## 握手过程

### 客户端请求

```
GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
```

### 服务端响应

```
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

### 密钥验证

```go
func computeAcceptKey(key string) string {
    const magic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
    h := sha1.New()
    h.Write([]byte(key + magic))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
```

---

## 数据帧格式

### 帧结构

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
```

### 操作码 (Opcode)

- `0x0`: 连续帧
- `0x1`: 文本帧
- `0x2`: 二进制帧
- `0x8`: 关闭帧
- `0x9`: Ping 帧
- `0xA`: Pong 帧

---

## 实现示例

### Go 实现

```go
package main

import (
    "github.com/gorilla/websocket"
    "net/http"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()
    
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            break
        }
        
        // 处理消息
        response := processMessage(message)
        
        // 发送响应
        if err := conn.WriteMessage(messageType, response); err != nil {
            break
        }
    }
}
```

### 心跳机制

```go
func heartbeat(conn *websocket.Conn) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

---

## 最佳实践

### 1. 连接管理

- 连接池管理
- 心跳保活
- 超时处理
- 重连机制

### 2. 消息处理

- 消息队列
- 异步处理
- 错误处理
- 限流控制

### 3. 安全考虑

- WSS (WebSocket Secure)
- 身份验证
- 权限控制
- 防攻击

---

## 参考资料

- [RFC 6455 - WebSocket Protocol](https://tools.ietf.org/html/rfc6455)
- [MDN WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

