# WebRTC 技术详解

WebRTC (Web Real-Time Communication) 实时通信技术详解。

## 目录

- [WebRTC 概述](#webrtc-概述)
- [核心组件](#核心组件)
- [信令流程](#信令流程)
- [ICE 机制](#ice-机制)
- [媒体处理](#媒体处理)
- [实际应用](#实际应用)

---

## WebRTC 概述

### 什么是 WebRTC

WebRTC 是一个开源项目，提供实时音视频通信能力。

**特点：**
- 浏览器原生支持
- P2P 通信
- 低延迟
- 加密传输

### 应用场景

- 视频会议
- 在线教育
- 远程医疗
- 游戏直播
- 在线客服

---

## 核心组件

### 1. MediaStream (getUserMedia)

**功能：**
- 获取音视频流
- 设备管理
- 流控制

**示例：**
```javascript
navigator.mediaDevices.getUserMedia({
    video: true,
    audio: true
}).then(stream => {
    // 使用 stream
});
```

### 2. RTCPeerConnection

**功能：**
- 建立 P2P 连接
- 媒体传输
- 连接管理

**关键方法：**
- `createOffer()`: 创建 offer
- `createAnswer()`: 创建 answer
- `setLocalDescription()`: 设置本地描述
- `setRemoteDescription()`: 设置远程描述
- `addIceCandidate()`: 添加 ICE 候选

### 3. RTCDataChannel

**功能：**
- 数据传输
- 低延迟
- 可靠或不可靠传输

---

## 信令流程

### 什么是信令服务器？

**信令服务器（Signaling Server）** 是 WebRTC 中用于**协调和协商**的服务器，它**不传输媒体数据**，只负责：

1. **交换连接信息**：帮助两个客户端互相发现和连接
2. **交换媒体描述（SDP）**：告诉对方自己支持什么编码格式、分辨率等
3. **交换网络信息（ICE 候选）**：告诉对方自己的网络地址
4. **会话管理**：管理通话的建立、结束等

**为什么需要信令服务器？**

WebRTC 是 **P2P（点对点）** 通信，但两个客户端在建立连接前：
- ❌ **不知道对方的 IP 地址**
- ❌ **不知道对方支持什么编码格式**
- ❌ **不知道如何建立连接**

**信令服务器的作用：**
```
客户端A ──信令服务器── 客户端B
   │                      │
   │  "我想和B通话"        │
   │  ←───────────────────┘
   │                      │
   │  "A想和你通话"        │
   │  ───────────────────→│
   │                      │
   │  交换SDP和ICE信息     │
   │  ←──────────────────→│
   │                      │
   │  建立P2P连接         │
   │  ←══════════════════→│
   │  (直接传输媒体数据)    │
```

**关键点：**
- ✅ **信令服务器只负责协商**，不传输音视频数据
- ✅ **媒体数据是 P2P 直连**，不经过信令服务器
- ✅ **信令可以用任何协议**：WebSocket、HTTP、SIP 等

### 信令服务器的工作流程

#### 1. 用户发现和连接请求

```javascript
// 客户端A：发起通话请求
signalingServer.send({
    type: 'offer',
    target: 'userB',
    sdp: offerSDP,  // 媒体描述
    iceCandidates: []  // 网络候选
});
```

#### 2. 信令服务器转发

```javascript
// 信令服务器：转发给客户端B
// 服务器代码（Node.js 示例）
io.on('connection', (socket) => {
    socket.on('offer', (data) => {
        // 转发给目标用户
        io.to(data.target).emit('offer', {
            from: socket.id,
            sdp: data.sdp,
            iceCandidates: data.iceCandidates
        });
    });
});
```

#### 3. 客户端B响应

```javascript
// 客户端B：接收并响应
signalingServer.on('offer', (data) => {
    // 创建 Answer
    peerConnection.setRemoteDescription(data.sdp);
    const answer = await peerConnection.createAnswer();
    await peerConnection.setLocalDescription(answer);
    
    // 发送 Answer 回客户端A
    signalingServer.send({
        type: 'answer',
        target: data.from,
        sdp: answer
    });
});
```

#### 4. 交换 ICE 候选

```javascript
// 客户端A：收集到新的 ICE 候选
peerConnection.onicecandidate = (event) => {
    if (event.candidate) {
        signalingServer.send({
            type: 'ice-candidate',
            target: 'userB',
            candidate: event.candidate
        });
    }
};
```

### 信令服务器的实现方式

#### 1. WebSocket 实现（推荐）

**优点：**
- ✅ 双向通信
- ✅ 实时性好
- ✅ 浏览器原生支持

**示例：**
```javascript
// 客户端
const ws = new WebSocket('wss://signaling.example.com');

ws.onopen = () => {
    // 发送信令
    ws.send(JSON.stringify({
        type: 'offer',
        target: 'userB',
        sdp: offerSDP
    }));
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === 'answer') {
        peerConnection.setRemoteDescription(data.sdp);
    }
};
```

```go
// 服务端（Go 示例）
func handleWebSocket(conn *websocket.Conn) {
    for {
        var msg SignalingMessage
        err := conn.ReadJSON(&msg)
        if err != nil {
            break
        }
        
        // 转发给目标用户
        targetConn := getUserConnection(msg.Target)
        if targetConn != nil {
            targetConn.WriteJSON(msg)
        }
    }
}
```

#### 2. HTTP 轮询实现

**优点：**
- ✅ 简单易实现
- ✅ 兼容性好

**缺点：**
- ❌ 延迟较高
- ❌ 服务器压力大

**示例：**
```javascript
// 客户端：轮询获取信令
setInterval(async () => {
    const response = await fetch('/api/signaling/poll?userId=userA');
    const messages = await response.json();
    messages.forEach(handleSignalingMessage);
}, 1000);  // 每秒轮询一次
```

#### 3. SIP 协议实现

**优点：**
- ✅ 标准协议
- ✅ 兼容传统电话系统

**适用场景：**
- 企业 VoIP 系统
- 与传统电话系统集成

### 信令服务器 vs 媒体服务器

| 特性 | 信令服务器 | 媒体服务器 |
|------|-----------|-----------|
| **作用** | 协商和协调 | 转发/处理媒体数据 |
| **数据量** | 小（文本消息） | 大（音视频流） |
| **延迟要求** | 较低（< 1秒） | 极低（< 200ms） |
| **协议** | WebSocket, HTTP, SIP | RTP/RTCP, WebRTC |
| **是否必需** | ✅ 必需（建立连接） | ⚠️ 可选（P2P 时不需要） |

### 信令服务器的功能

#### 1. 会话管理

```javascript
// 建立会话
{
    type: 'call-request',
    from: 'userA',
    to: 'userB',
    sessionId: 'session-123'
}

// 接受/拒绝
{
    type: 'call-accept',
    sessionId: 'session-123'
}

// 结束会话
{
    type: 'call-end',
    sessionId: 'session-123'
}
```

#### 2. 用户状态管理

```javascript
// 用户上线
{
    type: 'user-online',
    userId: 'userA'
}

// 用户离线
{
    type: 'user-offline',
    userId: 'userA'
}
```

#### 3. 房间管理（多人会议）

```javascript
// 加入房间
{
    type: 'join-room',
    roomId: 'room-123',
    userId: 'userA'
}

// 房间内用户列表
{
    type: 'room-users',
    roomId: 'room-123',
    users: ['userA', 'userB', 'userC']
}
```

### 信令服务器的架构设计

#### 简单架构（小规模）

```
┌─────────────┐
│  客户端A     │
└──────┬──────┘
       │ WebSocket
       │
┌──────▼──────────────┐
│   信令服务器          │
│  (单机)              │
└──────┬──────────────┘
       │ WebSocket
       │
┌──────▼──────┐
│  客户端B     │
└─────────────┘
```

#### 分布式架构（大规模）

```
┌─────────────┐
│  客户端A     │
└──────┬──────┘
       │
┌──────▼──────────────┐
│   信令服务器集群      │
│  ┌────┐  ┌────┐    │
│  │ S1 │  │ S2 │    │
│  └────┘  └────┘    │
│     │        │     │
│  ┌──▼────────▼──┐  │
│  │  Redis Pub/Sub │ │
│  └───────────────┘  │
└──────┬──────────────┘
       │
┌──────▼──────┐
│  客户端B     │
└─────────────┘
```

**实现要点：**
- 使用 **Redis Pub/Sub** 或 **消息队列** 实现服务器间通信
- 客户端可以连接到任意信令服务器
- 服务器间转发信令消息

### 信令服务器的安全考虑

#### 1. 身份验证

```javascript
// 连接时验证身份
ws.on('connection', (socket) => {
    socket.on('authenticate', (token) => {
        const user = verifyToken(token);
        if (user) {
            socket.userId = user.id;
        } else {
            socket.close();
        }
    });
});
```

#### 2. 权限控制

```javascript
// 检查用户是否有权限呼叫目标用户
function canCall(from, to) {
    // 检查黑名单
    if (isBlocked(from, to)) {
        return false;
    }
    
    // 检查好友关系
    if (!isFriend(from, to)) {
        return false;
    }
    
    return true;
}
```

#### 3. 消息加密

```javascript
// 使用 TLS/SSL 加密 WebSocket 连接
const ws = new WebSocket('wss://signaling.example.com');

// 或者对消息内容加密
const encrypted = encrypt(JSON.stringify(message));
ws.send(encrypted);
```

### 基本流程

```
1. 创建 RTCPeerConnection
   ↓
2. 添加本地流
   ↓
3. 创建 Offer
   ↓
4. 设置本地描述
   ↓
5. 发送 Offer 给对端（通过信令服务器）
   ↓
6. 对端接收 Offer，设置远程描述
   ↓
7. 对端创建 Answer
   ↓
8. 对端设置本地描述（Answer）
   ↓
9. 对端发送 Answer 给本端
   ↓
10. 本端设置远程描述（Answer）
   ↓
11. 交换 ICE 候选（通过信令服务器）
   ↓
12. 建立 P2P 连接，开始传输媒体（不经过信令服务器）
```

### SDP 交换

**SDP (Session Description Protocol)** 包含：
- 媒体类型（音频/视频）
- 编码格式
- 网络信息
- 传输协议

**示例：**
```
v=0
o=- 123456 2 IN IP4 127.0.0.1
s=-
t=0 0
m=audio 9 UDP/TLS/RTP/SAVPF 111
a=rtpmap:111 opus/48000/2
m=video 9 UDP/TLS/RTP/SAVPF 96
a=rtpmap:96 VP8/90000
```

**SDP 通过信令服务器交换：**
```
客户端A → 信令服务器 → 客户端B
  (SDP Offer)          (接收 Offer)
  
客户端B → 信令服务器 → 客户端A
  (SDP Answer)         (接收 Answer)
```

---

## ICE 机制

### ICE 候选类型

1. **Host Candidate**
   - 本地 IP 地址
   - 优先级最高

2. **Server Reflexive Candidate**
   - 通过 STUN 获取的公网 IP
   - 处理简单 NAT

3. **Relay Candidate**
   - 通过 TURN 获取的中继地址
   - 处理复杂 NAT

### STUN/TURN 服务器

**STUN (Session Traversal Utilities for NAT)**
- 获取公网 IP
- 检测 NAT 类型
- 免费使用

**TURN (Traversal Using Relays around NAT)**
- 中继服务器
- 处理复杂 NAT
- 需要服务器资源

### ICE 连接检查

```
1. 收集 ICE 候选
   ↓
2. 交换候选（通过信令）
   ↓
3. 连接检查（按优先级）
   ↓
4. 选择最佳连接
   ↓
5. 建立连接
```

---

## 媒体处理

### 音频处理流程

```
麦克风 → 采集 → 编码（Opus） → RTP 打包 → 网络传输
                                              ↓
扬声器 ← 播放 ← 解码（Opus） ← RTP 解包 ← 网络接收
```

### 视频处理流程

```
摄像头 → 采集 → 编码（VP8/H.264） → RTP 打包 → 网络传输
                                              ↓
显示器 ← 渲染 ← 解码（VP8/H.264） ← RTP 解包 ← 网络接收
```

### 编码参数

**音频：**
- 编码器：Opus（推荐）
- 采样率：48kHz
- 声道：立体声

**视频：**
- 编码器：VP8, VP9, H.264
- 分辨率：自适应
- 帧率：30fps

---

## 实际应用

### 1. 一对一通话

**架构：**
- 直连 P2P
- 信令服务器
- STUN/TURN 服务器

### 2. 多人会议

**架构选择：**
- **SFU (Selective Forwarding Unit)**
  - 服务器转发
  - 低延迟
  - 适合大规模

- **MCU (Multipoint Control Unit)**
  - 服务器混流
  - 客户端简单
  - 适合小规模

### 3. 直播推流

**流程：**
```
主播 → WebRTC 推流 → 媒体服务器 → 转码 → CDN → 观众
```

---

## 性能优化

### 1. 延迟优化

- 使用低延迟编码器
- 减少缓冲
- 优化网络路径

### 2. 带宽优化

- 自适应码率
- 动态分辨率
- 帧率控制

### 3. 质量优化

- 网络自适应
- 丢包恢复
- 错误隐藏

---

## 常见问题

### Q: WebRTC 连接失败怎么办？

**排查步骤：**
1. 检查 STUN/TURN 服务器
2. 检查防火墙设置
3. 检查 ICE 候选交换
4. 查看浏览器控制台日志

### Q: 如何降低延迟？

**方法：**
1. 使用低延迟编码器
2. 减少缓冲时间
3. 优化网络路径
4. 使用硬件加速

---

## 参考资料

- [WebRTC 官方文档](https://webrtc.org/)
- [MDN WebRTC 指南](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API)
- [WebRTC Samples](https://webrtc.github.io/samples/)

