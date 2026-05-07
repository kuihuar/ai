

## openTelemetry


### tracing 
- 手动埋点核心业务
- 半自动中间件标准化监控
- eBPF全自动不修改代码
- 一般保存到kafka,再保存到ES

### log

- log 也存入Kafka-ES

### metrics
- metrics需要存入时序数据库，采取拉取模式
- 1. 入口实现prometheus 的exporter
- 短时任务push
- http端口用pull
- metrics 会在SDK 先聚合，再导出