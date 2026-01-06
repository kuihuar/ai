# RTC & IM 知识体系

实时通信（RTC）和即时通讯（IM）技术知识整理，包括基础知识、面试问答、工程实践等。

## 📚 目录结构

```
rtc-im/
├── README.md           # 本文件，总目录
├── basics/             # 基础知识
│   ├── concepts.md     # 核心概念
│   ├── architecture.md # 架构设计
│   └── protocols.md    # 协议基础
├── interview/          # 面试问答
│   ├── rtc-qa.md      # RTC 面试题
│   ├── im-qa.md       # IM 面试题
│   └── system-design.md # 系统设计题
├── engineering/       # 工程实践
│   ├── best-practices.md # 最佳实践
│   ├── performance.md    # 性能优化
│   ├── troubleshooting.md # 问题排查
│   └── case-studies.md   # 案例分析
├── rtc/               # 实时音视频
│   ├── webrtc.md      # WebRTC 技术
│   ├── streaming.md   # 流媒体技术
│   ├── codec.md       # 编解码
│   └── network.md     # 网络传输
├── im/                # 即时通讯
│   ├── messaging.md   # 消息系统
│   ├── presence.md    # 在线状态
│   ├── push.md        # 推送技术
│   └── storage.md     # 消息存储
├── protocols/         # 协议详解
│   ├── websocket.md   # WebSocket
│   ├── xmpp.md        # XMPP
│   ├── mqtt.md        # MQTT
│   └── sip.md         # SIP
└── tools/             # 工具和框架
    ├── frameworks.md  # 开源框架
    ├── sdk.md         # SDK 使用
    └── monitoring.md  # 监控工具
```

## 🎯 知识体系

### 1. 基础知识 (basics/)

- [ ] 核心概念
  - [ ] RTC vs IM 的区别
  - [ ] 实时通信的技术挑战
  - [ ] 延迟、丢包、带宽的关系
  - [ ] QoS 和 QoE

- [ ] 架构设计
  - [ ] 客户端-服务器架构
  - [ ] P2P 架构
  - [ ] 混合架构
  - [ ] 分布式系统设计

- [ ] 协议基础
  - [ ] TCP vs UDP
  - [ ] 应用层协议选择
  - [ ] 信令协议
  - [ ] 媒体传输协议

### 2. 面试问答 (interview/)

- [ ] RTC 面试题
  - [ ] WebRTC 工作原理
  - [ ] 音视频编解码
  - [ ] 网络自适应
  - [ ] 延迟优化

- [ ] IM 面试题
  - [ ] 消息可靠性保证
  - [ ] 消息顺序性
  - [ ] 离线消息处理
  - [ ] 群聊实现

- [ ] 系统设计题
  - [ ] 设计一个 IM 系统
  - [ ] 设计一个视频会议系统
  - [ ] 设计一个直播系统
  - [ ] 设计一个实时游戏系统

### 3. 工程实践 (engineering/)

- [ ] 最佳实践
  - [ ] 架构设计原则
  - [ ] 代码组织
  - [ ] 错误处理
  - [ ] 日志和监控

- [ ] 性能优化
  - [ ] 延迟优化
  - [ ] 带宽优化
  - [ ] CPU/内存优化
  - [ ] 网络优化

- [ ] 问题排查
  - [ ] 常见问题
  - [ ] 调试技巧
  - [ ] 性能分析
  - [ ] 故障处理

- [ ] 案例分析
  - [ ] 微信 IM 架构
  - [ ] 钉钉音视频
  - [ ] Zoom 技术栈
  - [ ] Discord 架构

### 4. 实时音视频 (rtc/)

- [ ] WebRTC 技术
  - [ ] 核心概念
  - [ ] 信令流程
  - [ ] ICE 候选
  - [ ] SDP 协商

- [ ] 流媒体技术
  - [ ] 推流和拉流
  - [ ] 转码
  - [ ] 混流
  - [ ] 录制

- [ ] 编解码
  - [ ] 音频编解码（Opus, AAC）
  - [ ] 视频编解码（H.264, VP8, VP9, AV1）
  - [ ] 编码参数调优
  - [ ] 硬件加速

- [ ] 网络传输
  - [ ] RTP/RTCP
  - [ ] 网络自适应
  - [ ] 拥塞控制
  - [ ] 丢包恢复

### 5. 即时通讯 (im/)

- [ ] 消息系统
  - [ ] 消息类型（文本、图片、文件、语音、视频）
  - [ ] 消息存储
  - [ ] 消息同步
  - [ ] 消息搜索

- [ ] 在线状态
  - [ ] 状态管理
  - [ ] 心跳机制
  - [ ] 状态同步
  - [ ] 离线检测

- [ ] 推送技术
  - [ ] APNs (iOS)
  - [ ] FCM (Android)
  - [ ] 第三方推送
  - [ ] 推送策略

- [ ] 消息存储
  - [ ] 数据库设计
  - [ ] 消息索引
  - [ ] 冷热数据分离
  - [ ] 数据迁移

### 6. 协议详解 (protocols/)

- [ ] WebSocket
  - [ ] 协议原理
  - [ ] 握手过程
  - [ ] 帧格式
  - [ ] 应用场景

- [ ] XMPP
  - [ ] 协议架构
  - [ ] 消息格式
  - [ ] 扩展协议
  - [ ] 优缺点

- [ ] MQTT
  - [ ] 协议特点
  - [ ] QoS 级别
  - [ ] 主题订阅
  - [ ] 应用场景

- [ ] SIP
  - [ ] 协议基础
  - [ ] 会话建立
  - [ ] 媒体协商
  - [ ] 应用场景

### 7. 工具和框架 (tools/)

- [ ] 开源框架
  - [ ] WebRTC (Google)
  - [ ] Janus Gateway
  - [ ] Kurento
  - [ ] MediaSoup
  - [ ] Pion (Go)

- [ ] SDK 使用
  - [ ] 客户端 SDK
  - [ ] 服务端 SDK
  - [ ] 集成指南
  - [ ] 常见问题

- [ ] 监控工具
  - [ ] 性能监控
  - [ ] 质量监控
  - [ ] 日志分析
  - [ ] 告警系统

## 📖 学习路径

### 初级
1. 理解 RTC 和 IM 的基本概念
2. 学习 WebSocket 协议
3. 了解 WebRTC 基础
4. 实现简单的聊天应用

### 中级
1. 深入理解 WebRTC 协议栈
2. 学习音视频编解码
3. 掌握网络自适应技术
4. 实现完整的 IM 系统

### 高级
1. 优化延迟和带宽
2. 设计大规模分布式系统
3. 处理复杂网络环境
4. 性能调优和问题排查

## 🔗 相关资源

### 官方文档
- [WebRTC 官方文档](https://webrtc.org/)
- [RFC 文档](https://www.rfc-editor.org/)

### 开源项目
- [Pion WebRTC](https://github.com/pion/webrtc) - Go 实现的 WebRTC
- [Janus Gateway](https://github.com/meetecho/janus-gateway) - WebRTC 服务器
- [MediaSoup](https://github.com/versatica/mediasoup) - SFU 服务器

### 技术博客
- WebRTC 技术博客
- 各大公司的技术分享

## 📝 已创建文档

### 基础知识
- ✅ [核心概念](./basics/concepts.md) - RTC vs IM、技术挑战、关键指标、架构模式

### 面试问答
- ✅ [RTC 面试问答](./interview/rtc-qa.md) - 17个问题，涵盖 WebRTC、编解码、网络传输、性能优化
- ✅ [IM 面试问答](./interview/im-qa.md) - 16个问题，涵盖消息系统、在线状态、推送技术、存储设计
- ✅ [系统设计题](./interview/system-design.md) - 4个系统设计（IM系统、视频会议、直播、实时游戏）

### 工程实践
- ✅ [最佳实践](./engineering/best-practices.md) - 架构设计、代码组织、错误处理、日志监控、性能优化、安全实践

### 实时音视频 (rtc/)
- ✅ [WebRTC 技术](./rtc/webrtc.md) - WebRTC 概述、核心组件、信令流程、ICE 机制、媒体处理
- ✅ [编解码技术](./rtc/codec.md) - 音视频编解码器详解（Opus、AAC、VP8/VP9、H.264等）
- ✅ [流媒体技术](./rtc/streaming.md) - 推流、拉流、转码、CDN 分发
- ✅ [网络传输技术](./rtc/network.md) - RTP/RTCP、网络自适应、拥塞控制、丢包恢复

### 即时通讯 (im/)
- ✅ [消息系统设计](./im/messaging.md) - 消息模型、消息流程、可靠性保证、顺序性保证、存储设计
- ✅ [推送技术](./im/push.md) - iOS推送(APNs)、Android推送(FCM)、Web推送、推送策略
- ✅ [在线状态管理](./im/presence.md) - 状态类型、状态管理、心跳机制、状态同步

### 协议详解 (protocols/)
- ✅ [WebSocket 协议](./protocols/websocket.md) - 协议概述、握手过程、数据帧格式、实现示例
- ✅ [MQTT 协议](./protocols/mqtt.md) - 协议概述、QoS级别、主题订阅、应用场景
- ✅ [XMPP 协议](./protocols/xmpp.md) - 协议概述、协议架构、消息格式、扩展协议

### 工具和框架 (tools/)
- ✅ [开源框架](./tools/frameworks.md) - WebRTC框架、IM框架、媒体服务器、SDK对比

## 📝 更新日志

- 2024-01-XX: 创建知识体系目录结构
- 2024-01-XX: 完成核心文档创建（基础知识、面试问答、工程实践）
- 2024-01-XX: 完善所有方向文档（RTC、IM、协议、工具框架）

## 📊 文档统计

- **基础知识**: 1 篇
- **面试问答**: 3 篇（共 37+ 个问题）
- **工程实践**: 1 篇
- **实时音视频**: 4 篇
- **即时通讯**: 3 篇
- **协议详解**: 3 篇
- **工具框架**: 1 篇

**总计**: 16 篇核心文档

