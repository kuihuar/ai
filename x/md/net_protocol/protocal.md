1. OSI七层模型和TCP/IP四层模型的区别？
OSI模型（理论标准）：

应用层、表示层、会话层、传输层、网络层、数据链路层、物理层。



应用层合并了表示层
TCP/IP模型（实际实现）：

应用层（HTTP/FTP/DNS）、传输层（TCP/UDP）、网络层（IP）、网络接口层（以太网/Wi-Fi）。

核心差异：

OSI的会话层、表示层在TCP/IP中被合并到应用层。

TCP/IP更注重实际协议实现，如IP协议直接对应网络层。

2. TCP和UDP的核心区别？
特性	TCP	UDP
连接性	面向连接（三次握手）	无连接
可靠性	可靠传输（确认、重传）	不保证可靠
顺序性	数据按序到达	可能乱序
流量控制	滑动窗口机制	无
头部开销	大（20~60字节）	小（8字节）
典型应用	HTTP、FTP、数据库	视频流、DNS、实时游戏



计算层数时从最底层物理层开始计数，依次向上递增。

第 1 层 - 物理层（Physical Layer）
功能：负责传输比特流，定义了物理设备和传输介质的电气、机械特性，如电缆的类型、接口的形状、信号的电压等。
示例：以太网电缆、光纤、无线信号等。
第 2 层 - 数据链路层（Data Link Layer）
功能：将物理层接收到的信号进行处理，形成有意义的数据帧，并负责帧的传输、差错检测和纠正。
示例：以太网协议、Wi - Fi 协议等，MAC 地址就工作在这一层。
第 3 层 - 网络层（Network Layer）
功能：负责将数据帧从源节点传输到目标节点，进行路由选择和寻址，使用逻辑地址（如 IP 地址）。
示例：IP 协议（IPv4、IPv6）。
第 4 层 - 传输层（Transport Layer）
功能：提供端到端的可靠通信，确保数据的正确传输，包括流量控制、错误恢复等。
示例：TCP（传输控制协议）提供可靠的、面向连接的传输；UDP（用户数据报协议）提供不可靠的、无连接的传输。
第 5 层 - 会话层（Session Layer）
功能：负责建立、管理和终止会话，协调不同主机上的应用程序之间的通信会话。
示例：远程登录协议（如 Telnet、SSH）在会话层建立和管理用户与远程服务器之间的会话。
第 6 层 - 表示层（Presentation Layer）
功能：处理数据的表示和转换，如数据的加密、解密、压缩、解压缩等，确保不同系统之间能够正确理解数据。
示例：JPEG 图像格式的压缩和解压缩、SSL/TLS 协议进行数据加密。
第 7 层 - 应用层（Application Layer）
功能：为用户的应用程序提供网络服务，如文件传输、电子邮件、网页浏览等。
示例：HTTP（超文本传输协议）用于网页浏览、SMTP（简单邮件传输协议）用于发送电子邮件。
TCP/IP 模型（4 层）
TCP/IP 模型是互联网实际使用的网络通信模型，它更加简洁实用，计算层数同样从最底层开始计数。


第 1 层 - 网络接口层（Network Interface Layer）
功能：对应 OSI 模型的物理层和数据链路层，负责将数据包封装成适合在物理网络上传输的帧，并进行物理传输。
示例：以太网、Wi - Fi 等网络接口设备。
第 2 层 - 网络层（Internet Layer）
功能：与 OSI 模型的网络层类似，负责数据包的路由和转发，使用 IP 地址进行寻址。
示例：IP 协议。
第 3 层 - 传输层（Transport Layer）
功能：和 OSI 模型的传输层功能相同，提供端到端的通信服务，确保数据的可靠传输或高效传输。
示例：TCP 和 UDP 协议。
第 4 层 - 应用层（Application Layer）
功能：整合了 OSI 模型的会话层、表示层和应用层的功能，为用户应用程序提供各种网络服务。
示例：HTTP、FTP（文件传输协议）、SMTP 等。


===========================

二、TCP协议深入
1. TCP三次握手的过程？为什么需要三次？
过程：

客户端发送SYN=1, seq=x。

服务端回复SYN=1, ACK=1, seq=y, ack=x+1。

客户端发送ACK=1, seq=x+1, ack=y+1。

为什么是三次？
防止已失效的连接请求报文突然传到服务端，导致资源浪费（两次可能建立无用连接）。

2. TCP四次挥手的过程？为什么需要四次？
过程：

主动方发送FIN=1, seq=u。

被动方回复ACK=1, seq=v, ack=u+1。

被动方发送FIN=1, ACK=1, seq=w, ack=u+1。

主动方回复ACK=1, seq=u+1, ack=w+1。

为什么是四次？
TCP是全双工协议，需双方分别关闭发送和接收通道。

3. TIME_WAIT状态的作用？
作用：

确保最后一个ACK能被被动关闭方接收（若丢失，被动方会重发FIN）。

等待网络中残留的旧报文过期，避免干扰新连接。

持续时间：2MSL（Maximum Segment Lifetime，通常2分钟）。

4. TCP如何保证可靠性？
机制：

校验和：检测数据是否损坏。

序列号与确认应答：确保数据有序到达。

超时重传：未收到ACK则重发。

流量控制：滑动窗口动态调整发送速率。

拥塞控制：慢启动、拥塞避免、快速重传、快速恢复。




三、HTTP与HTTPS
1. HTTP/1.1、HTTP/2、HTTP/3的区别？
HTTP/1.1：

持久连接（Keep-Alive），但队头阻塞（一个请求阻塞后续请求）。

HTTP/2：

多路复用（Multiplexing），头部压缩（HPACK），服务器推送（Server Push）。

HTTP/3：

基于QUIC协议（UDP实现），解决TCP队头阻塞，支持0-RTT握手。

2. HTTPS的握手过程？
客户端发送支持的TLS版本、加密套件、随机数。

服务器返回选择的加密套件、随机数、数字证书。

客户端验证证书，生成预主密钥（Pre-Master Secret），用服务器公钥加密发送。

双方通过随机数和预主密钥生成会话密钥（Session Key）。

后续通信使用对称加密。

3. HTTP状态码及其含义？
200 OK：请求成功。

301 Moved Permanently：永久重定向。

400 Bad Request：客户端请求错误。

401 Unauthorized：未认证。

403 Forbidden：无权限。

404 Not Found：资源不存在。

500 Internal Server Error：服务器内部错误。

502 Bad Gateway：网关错误。



四、DNS与网络层
1. DNS解析过程？
浏览器缓存 → 本地Hosts文件 → 本地DNS服务器。

本地DNS依次查询根域名服务器、顶级域（.com）、权威域名服务器。

返回IP地址并缓存。

2. CDN的工作原理？
核心：将内容缓存到离用户最近的边缘节点。

过程：

用户请求资源时，DNS解析返回最佳CDN节点IP。

CDN节点若有缓存则直接返回，否则回源站获取。


五、场景与设计题
1. 如何设计一个实时消息系统（如微信）的网络协议？
需求：低延迟、可靠、支持海量连接。

方案：

传输层：TCP保证消息可靠（如文本消息），UDP用于音视频流。

应用层：自定义协议或MQTT/WebSocket。

优化：心跳机制、消息重试、离线缓存。

2. TCP粘包问题如何解决？
原因：TCP是字节流协议，无消息边界。

方案：

定长协议：每个消息固定长度。

分隔符：如\r\n标记结束。

头部声明长度：如HTTP的Content-Length。


六、高频进阶问题
QUIC协议解决了哪些问题？

基于UDP，避免TCP队头阻塞。

内置加密（TLS 1.3），减少握手延迟（0-RTT）。

连接迁移（切换网络IP不影响连接）。

WebSocket与HTTP长轮询的区别？

WebSocket：全双工通信，服务端可主动推送。

HTTP长轮询：客户端轮询请求，服务器保持连接直到有数据。

ARP协议的作用？

将IP地址解析为MAC地址，实现局域网内通信。

ARP欺骗攻击如何防范？


实战代码示例（TCP Server）

```go
// Go语言实现简单TCP服务端
package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("Received: %s\n", string(buf[:n]))
	conn.Write([]byte("Hello from server!"))
}
```