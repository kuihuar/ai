# XMPP 协议

XMPP (Extensible Messaging and Presence Protocol) 协议详解。

## 目录

- [协议概述](#协议概述)
- [协议架构](#协议架构)
- [消息格式](#消息格式)
- [扩展协议](#扩展协议)

---

## 协议概述

### 什么是 XMPP

XMPP 是一个基于 XML 的即时通讯协议。

**特点：**
- 开放标准
- 可扩展
- 分布式
- 实时通信

### 应用场景

- 企业 IM
- 聊天应用
- IoT 设备通信
- 游戏通信

---

## 协议架构

### 客户端-服务器架构

```
客户端1 ←→ XMPP 服务器 ←→ 客户端2
```

### 服务器间通信

```
服务器1 ←→ 服务器2
```

---

## 消息格式

### 基本消息

```xml
<message
    from='alice@example.com'
    to='bob@example.com'
    type='chat'>
    <body>Hello, Bob!</body>
</message>
```

### 在线状态

```xml
<presence
    from='alice@example.com'
    to='bob@example.com'>
    <show>away</show>
    <status>In a meeting</status>
</presence>
```

---

## 扩展协议

### Jingle (音视频)

XMPP 的音视频扩展协议。

### File Transfer

文件传输扩展。

---

## 参考资料

- [XMPP 官方](https://xmpp.org/)
- [XMPP 规范](https://xmpp.org/rfcs/)

