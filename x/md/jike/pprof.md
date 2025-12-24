1. profile 采样CPU
2. allocs侧重于频繁进行内存分配的函数
3. heap 查看存活对象的内存分配情况，侧重于定位内存泄漏问题。
4. trace 工具查看是什么导致了接口口响应延时高

trace 会记录以下事件：
1. 协程的创建过程，开始运行时刻及结束运行的时间点
2. 协程由于系统调用、通道操作、锁的使用等情况而出现被阻塞的现象
3. 网络IO相关的操作情况
4. 垃圾回收的相关情况
```shell
curl "http://localhost:8888/debug/pprof/trace?seconds=30s" > trace.out
```

