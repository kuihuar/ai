# 流媒体技术

实时音视频流媒体技术详解。

## 目录

- [流媒体概述](#流媒体概述)
- [推流技术](#推流技术)
- [拉流技术](#拉流技术)
- [转码技术](#转码技术)
- [CDN 分发](#cdn-分发)

---

## 流媒体概述

### 什么是流媒体

流媒体是指通过网络实时传输音视频数据的技术。

**特点：**
- 实时传输
- 边下载边播放
- 不需要完整下载

### 应用场景

- 直播
- 点播
- 视频会议
- 在线教育

---

## 推流技术

### RTMP 推流

**RTMP (Real-Time Messaging Protocol)**
- Adobe 开发
- 基于 TCP
- 广泛支持

**流程：**
```
采集 → 编码 → RTMP 打包 → 推流服务器
```

**实现：**
```go
// 使用 FFmpeg 推流
ffmpeg -i input.mp4 -c copy -f flv rtmp://server/live/stream
```

### WebRTC 推流

**特点：**
- 低延迟
- 浏览器原生
- P2P 或服务器

**流程：**
```
浏览器 → WebRTC → 媒体服务器 → 转码 → CDN
```

### SRT 推流

**SRT (Secure Reliable Transport)**
- 低延迟
- 可靠传输
- 加密

**适用场景：**
- 专业直播
- 低延迟需求

---

## 拉流技术

### HLS (HTTP Live Streaming)

**特点：**
- Apple 开发
- 基于 HTTP
- 自适应码率

**格式：**
```
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXTINF:10.0,
segment1.ts
#EXTINF:10.0,
segment2.ts
```

**优点：**
- 广泛支持
- CDN 友好
- 自适应

**缺点：**
- 延迟较高（> 10秒）

### HTTP-FLV

**特点：**
- 低延迟（< 3秒）
- 基于 HTTP
- Flash 格式

**适用场景：**
- 实时直播
- 低延迟需求

### DASH

**DASH (Dynamic Adaptive Streaming over HTTP)**
- 自适应码率
- 标准化
- 广泛支持

---

## 转码技术

### 转码目的

1. **多码率适配**
   - 不同网络环境
   - 不同设备
   - 自适应播放

2. **格式转换**
   - 推流格式 → 拉流格式
   - 编码格式转换

3. **质量优化**
   - 降噪
   - 增强
   - 水印

### 转码流程

```
输入流 → 解码 → 处理 → 编码 → 输出流
```

### 实现方案

**FFmpeg：**
```bash
# 转码为多码率
ffmpeg -i input.mp4 \
  -c:v libx264 -b:v 1000k -s 640x360 output_360p.mp4 \
  -c:v libx264 -b:v 2500k -s 1280x720 output_720p.mp4 \
  -c:v libx264 -b:v 5000k -s 1920x1080 output_1080p.mp4
```

**硬件加速：**
```bash
# 使用 GPU 加速
ffmpeg -hwaccel cuda -i input.mp4 -c:v h264_nvenc output.mp4
```

---

## CDN 分发

### CDN 作用

- 就近分发
- 减少延迟
- 降低带宽成本
- 提高可用性

### CDN 架构

```
源站 → CDN 边缘节点 → 用户
```

### 缓存策略

**缓存内容：**
- 静态资源（HLS 切片）
- 图片、文件
- 配置信息

**缓存时间：**
- 直播：不缓存或短缓存
- 点播：长缓存
- 静态资源：永久缓存

---

## 性能优化

### 1. 延迟优化

- 使用低延迟协议（SRT, WebRTC）
- 减少缓冲
- 优化转码流程
- 就近部署

### 2. 带宽优化

- 自适应码率
- 多码率转码
- 压缩优化
- CDN 加速

### 3. 质量优化

- 编码参数调优
- 预处理优化
- 后处理增强
- 硬件加速

---

## 参考资料

- [FFmpeg 文档](https://ffmpeg.org/documentation.html)
- [HLS 规范](https://tools.ietf.org/html/rfc8216)
- [流媒体技术](https://en.wikipedia.org/wiki/Streaming_media)

