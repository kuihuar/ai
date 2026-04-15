# Kubeflow Pipeline 实战示例

## 目标

构建一个简单的训练流水线：`数据预处理 -> 训练 -> 评估 -> 模型发布`。

## 流程拆分建议

1. **preprocess**：读取数据并输出清洗后的数据集
2. **train**：训练模型并输出模型文件
3. **evaluate**：计算指标并判断是否达标
4. **deploy**：达标则触发 KServe 发布

## 示例 DSL（简化）

```python
from kfp import dsl

@dsl.pipeline(name="demo-train-pipeline")
def train_pipeline():
    p = preprocess_op()
    t = train_op(p.outputs["dataset"])
    e = evaluate_op(t.outputs["model"])
    deploy_op(t.outputs["model"]).after(e)
```

## 运行与观测

- 在 Kubeflow UI 提交 Pipeline Run
- 观察每个步骤的日志与产物
- 通过指标面板查看准确率、召回率等指标

## 常见问题

- **任务 Pending**：检查节点资源和调度约束
- **镜像拉取失败**：检查镜像地址、凭据、网络
- **产物丢失**：检查对象存储配置和访问权限

## 下一步

- 接入 Katib 做自动化参数搜索
- 给 Pipeline 增加数据质量校验步骤
- 引入多环境发布（dev/staging/prod）

