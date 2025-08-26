package skill

// 6. 系统级同步原语
// 6.1 系统调用
// runtime.Gosched() - 让出CPU时间片
// runtime.LockOSThread() - 锁定当前goroutine到系统线程
// runtime.UnlockOSThread() - 解锁当前goroutine
// 6.2 内存屏障
// runtime.GC() - 强制垃圾回收
// runtime.ReadMemStats() - 读取内存统计信息
// 7. 网络和I/O同步原语
// 7.1 网络同步
// net.Listener - 网络监听器
// net.Conn - 网络连接
// http.Server - HTTP服务器
// 7.2 I/O同步
// os.File - 文件操作
// bufio.Reader/Writer - 缓冲I/O
// io.Pipe - 管道
