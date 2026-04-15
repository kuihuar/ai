# AI 组件对比：推理服务与网关

## 对比范围

- KServe
- vLLM
- Triton Inference Server
- Ray Serve

## 一句话定位

- **KServe**：模型服务治理平台（部署/版本/流量）
- **vLLM**：LLM 推理性能引擎
- **Triton**：多框架高性能推理服务器
- **Ray Serve**：分布式服务编排框架

## 核心能力对比

| 维度 | KServe | vLLM | Triton | Ray Serve |
| :--- | :--- | :--- | :--- | :--- |
| 模型生命周期治理 | 强 | 弱 | 中 | 中 |
| LLM 推理吞吐优化 | 中 | 强 | 中高 | 中 |
| 多框架支持 | 高 | LLM 为主 | 高 | 中 |
| 灰度/金丝雀 | 强 | 需外部实现 | 需外部实现 | 可编排实现 |
| 上手复杂度 | 中高 | 中 | 中高 | 中 |

## 选型建议

- **企业推理平台首选**：KServe + (vLLM/Triton 作为 runtime)
- **LLM 单场景高性能**：vLLM
- **多框架统一推理**：Triton
- **复杂在线编排**：Ray Serve

## 落地建议

1. 先定“控制平面”（KServe 或 Ray Serve）
2. 再定“运行时引擎”（vLLM / Triton）
3. 最后补灰度、限流、观测与回滚策略
