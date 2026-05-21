# 基础设施与可观测（口述问答）

> 对应 `readme.md` 中「基础设施」与监控日志：Docker、K8s、对象存储、Prometheus、OpenTelemetry、zap / zerolog。

## 1) Docker 多阶段构建解决什么问题？

**口述参考：**  
构建阶段带编译器与依赖，运行阶段只留二进制和最小运行时镜像，减小攻击面和拉取体积。CI 里 reproducible build，固定基础镜像 digest，分层缓存加速构建。

## 2) Kubernetes 你用过哪些概念？发布要注意什么？

**口述参考：**  
Deployment、Service、Ingress、ConfigMap/Secret、HPA、探针 liveness/readiness。发布要灰度或滚动更新， readiness 失败不摘流量会雪崩；资源 request/limit 要合理，避免 OOM 与 CPU throttle。

## 3) 对象存储预签名 URL 的使用场景与风险？

**口述参考：**  
适合直传大文件、降低服务端带宽，短期授权下载上传。风险是 URL 泄露与时间窗过长，要用最短 TTL、HTTPS、必要时绑定 IP 或一次性 token。服务端仍要做业务鉴权与配额。

## 4) 策略模式在存储访问里你怎么理解？

**口述参考：**  
对 S3、MinIO、OBS、OSS 抽象统一接口，按配置选择实现，便于多云切换与单测 mock。注意各厂商签名与 header 差异，错误码与重试策略要封装。

## 5) Prometheus 指标设计有什么原则？

**口述参考：**  
指标名稳定、label 高基数要克制；RED/USE 方法看请求率、错误、延迟与资源饱和度。告警规则可执行，避免「狼来了」。业务指标与基础设施指标分层。

## 6) OpenTelemetry 你用来解决什么？

**口述参考：**  
统一 traces、metrics、logs 的关联，跨服务传播 trace context。接入后能在一次请求里看到各 span 耗时，定位下游慢与并行度问题。采样策略要平衡成本与可观测性。

## 7) zap 和 zerolog 选型差异？

**口述参考：**  
都偏向结构化高性能日志。zap 字段 API 成熟、生态多；zerolog API 链式、零分配取向。团队统一一种，约定日志级别、采样与敏感字段脱敏。

## 8) 线上故障「止血 → 定位 → 修复」你怎么描述一次经历？

**口述参考：**  
先限流降级或回滚恢复 SLO，再靠监控/trace 定位根因，修复后做复盘：缺失的告警、单测、容量评估。回答时带时间线和指标前后对比更有说服力。

## 9) Agent 平台长连接与批任务，监控上要额外看什么？

**口述参考：**  
连接数、每连接内存、goroutine、消息队列 lag、异步任务失败率、模型调用超时与重试。把「一次对话」作为 trace 根，便于关联 LLM 与工具调用。
