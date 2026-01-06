# 网络传输技术

实时音视频网络传输技术详解。

## 目录

- [RTP/RTCP](#rtprtcp)
- [网络自适应](#网络自适应)
- [拥塞控制](#拥塞控制)
- [丢包恢复](#丢包恢复)

---

## RTP/RTCP

### RTP (Real-time Transport Protocol)

**特点：**
- 传输音视频数据
- 基于 UDP
- 提供时间戳、序列号
- 不保证可靠性

**RTP 头格式：**
```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X|  CC   |M|     PT      |       sequence number         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           timestamp                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
|            contributing source (CSRC) identifiers             |
|                             ....                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

### RTCP (RTP Control Protocol)

**功能：**
- 传输控制信息
- 质量反馈
- 同步信息

**RTCP 包类型：**
- SR (Sender Report): 发送者报告
- RR (Receiver Report): 接收者报告
- SDES: 源描述
- BYE: 离开

---

## 网络自适应

### 带宽估计

**方法：**
1. **基于延迟**
   - 测量 RTT
   - 延迟增加 → 带宽减少

2. **基于丢包**
   - 统计丢包率
   - 丢包率高 → 带宽减少

3. **基于吞吐量**
   - 测量实际吞吐量
   - 动态调整

### 码率调整

```go
type BandwidthEstimator struct {
    currentBandwidth int
    lossRate         float64
    rtt              time.Duration
}

func (b *BandwidthEstimator) Estimate() int {
    // 基于丢包率调整
    if b.lossRate > 0.1 {
        b.currentBandwidth = int(float64(b.currentBandwidth) * 0.8)
    } else if b.lossRate < 0.01 {
        b.currentBandwidth = int(float64(b.currentBandwidth) * 1.1)
    }
    
    return b.currentBandwidth
}
```

---

## 拥塞控制

### GCC (Google Congestion Control)

**算法：**
1. 基于延迟的控制器
2. 基于丢包的控制器
3. 选择较小的码率

### 实现策略

```go
func adjustBitrate(lossRate float64, rtt time.Duration, currentBitrate int) int {
    // 基于丢包
    if lossRate > 0.02 {
        return int(float64(currentBitrate) * 0.85)
    }
    
    // 基于延迟
    if rtt > 200*time.Millisecond {
        return int(float64(currentBitrate) * 0.9)
    }
    
    // 网络良好，可以增加
    if lossRate < 0.001 && rtt < 50*time.Millisecond {
        return int(float64(currentBitrate) * 1.05)
    }
    
    return currentBitrate
}
```

---

## 丢包恢复

### 前向纠错 (FEC)

**原理：**
- 发送冗余数据
- 可以恢复部分丢包
- 增加带宽开销

**实现：**
```go
func encodeFEC(data []byte) ([]byte, []byte) {
    // 原始数据
    original := data
    
    // 生成冗余数据
    redundant := generateRedundant(original)
    
    return original, redundant
}
```

### 重传机制

**选择性重传：**
- 只重传关键帧
- 减少带宽开销
- 保证质量

**实现：**
```go
func handlePacketLoss(seqNum int) {
    // 检查是否是关键帧
    if isKeyFrame(seqNum) {
        // 请求重传
        requestRetransmission(seqNum)
    } else {
        // 使用错误隐藏
        errorConcealment(seqNum)
    }
}
```

### 错误隐藏

**方法：**
1. **帧复制**
   - 使用前一帧
   - 简单有效

2. **插值**
   - 前后帧插值
   - 更平滑

3. **运动补偿**
   - 基于运动向量
   - 更准确

---

## 参考资料

- [RTP 规范 RFC 3550](https://tools.ietf.org/html/rfc3550)
- [WebRTC 拥塞控制](https://webrtc.org/)
- [网络自适应算法](https://github.com/)

