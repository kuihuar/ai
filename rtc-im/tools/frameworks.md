# RTC & IM 开源框架

实时通信和即时通讯相关的开源框架和工具。

## 目录

- [WebRTC 框架](#webrtc-框架)
- [IM 框架](#im-框架)
- [媒体服务器](#媒体服务器)
- [SDK](#sdk)

---

## WebRTC 框架

### Pion WebRTC (Go)

**特点：**
- Go 语言实现
- 纯 Go，无 C 依赖
- 跨平台
- 活跃维护

**GitHub:** https://github.com/pion/webrtc

**适用场景：**
- Go 后端服务
- 服务器端 WebRTC
- 自定义实现

### Janus Gateway

**特点：**
- C 语言实现
- 高性能
- 插件架构
- 支持多种协议

**GitHub:** https://github.com/meetecho/janus-gateway

**适用场景：**
- 视频会议
- 直播
- 录制

### MediaSoup

**特点：**
- Node.js 实现
- SFU 架构
- 高性能
- 易扩展

**GitHub:** https://github.com/versatica/mediasoup

**适用场景：**
- 大规模会议
- 实时通信
- 流媒体

### Kurento

**特点：**
- Java 实现
- MCU 架构
- 媒体处理
- 录制转码

**GitHub:** https://github.com/Kurento/kurento-media-server

**适用场景：**
- 小规模会议
- 媒体处理
- 录制服务

---

## IM 框架

### OpenIM

**特点：**
- Go 语言实现
- 开源 IM 框架
- 完整功能
- 高性能

**GitHub:** https://github.com/openimsdk/open-im-server

**功能：**
- 单聊、群聊
- 消息推送
- 在线状态
- 文件传输

### Rocket.Chat

**特点：**
- Node.js 实现
- 企业级 IM
- 丰富功能
- 可定制

**GitHub:** https://github.com/RocketChat/Rocket.Chat

**适用场景：**
- 企业内部通讯
- 团队协作
- 客服系统

### Mattermost

**特点：**
- Go 语言实现
- 企业级
- 开源
- 可自托管

**GitHub:** https://github.com/mattermost/mattermost

**适用场景：**
- 企业协作
- 团队沟通
- 项目管理

---

## 媒体服务器

### SRS (Simple Realtime Server)

**特点：**
- C++ 实现
- 高性能
- 支持多种协议
- 易部署

**GitHub:** https://github.com/ossrs/srs

**功能：**
- RTMP 推流
- HLS 拉流
- WebRTC
- 录制转码

### Livego

**特点：**
- Go 语言实现
- 简单易用
- 轻量级
- 快速部署

**GitHub:** https://github.com/gwuhaolin/livego

**适用场景：**
- 简单直播
- 快速原型
- 学习研究

---

## SDK

### WebRTC SDK

**客户端 SDK：**
- **WebRTC JS SDK**: 浏览器原生
- **libwebrtc**: C++ 库
- **Pion WebRTC**: Go SDK

**服务端 SDK：**
- **Pion**: Go
- **aiortc**: Python
- **Kurento Client**: Java/JS

### IM SDK

**开源 SDK：**
- **OpenIM SDK**: Go/JS/Java
- **Rocket.Chat SDK**: JS/Java/Swift
- **Mattermost SDK**: JS/Go

**商业 SDK：**
- 环信 SDK
- 融云 SDK
- 腾讯云 IM SDK

---

## 框架选择指南

### WebRTC 场景

**Go 后端：**
- 推荐：Pion WebRTC
- 原因：纯 Go，易集成

**Node.js 后端：**
- 推荐：MediaSoup
- 原因：高性能，易扩展

**Java 后端：**
- 推荐：Kurento
- 原因：成熟稳定

### IM 场景

**自建 IM：**
- 推荐：OpenIM
- 原因：功能完整，性能好

**企业协作：**
- 推荐：Mattermost
- 原因：企业级功能

**快速开发：**
- 推荐：Rocket.Chat
- 原因：功能丰富，易定制

---

## 参考资料

- [WebRTC 官方](https://webrtc.org/)
- [开源 IM 对比](https://github.com/)
- [媒体服务器对比](https://github.com/)

