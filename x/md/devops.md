Operator是基于Kubernetes的API扩展机制，用于管理复杂有状态应用的自定义控制器。它通过自定义资源定义（CustomResourceDefinition, CRD）来表示应用的配置和状态，并通过控制器来监控这些资源的变化，执行相应的操作，如创建、更新和删除应用实例。
Operator：侧重于管理复杂有状态应用的生命周期。它可以根据应用的状态和变化自动执行复杂的操作，如数据库的备份恢复、应用的滚动升级等。Operator通过自定义资源和控制器来实现对应用的精细化管理。
Helm Chart：主要用于应用的部署和配置管理。它是一种模板化的包管理工具，通过将Kubernetes资源文件打包成Chart，可以方便地在不同环境中部署和升级应用。Helm Chart更关注应用的初始部署和配置。
Helm也可以用于部署Operator。可以将Operator的相关资源（如CRD、Deployment等）打包成Helm Chart，通过Helm命令进行部署和管理。这样可以简化Operator的部署过程，提高部署效率。
互补使用：在实际应用中，Operator和Helm Chart可以互补使用。可以使用Helm Chart来进行应用的初始部署，将应用的基本配置和资源部署到Kubernetes集群中。然后，使用Operator来管理应用的生命周期，确保应用在运行过程中始终处于健康状态。
升级方案：
复杂的operator
无状态的helm
其它gitlab CICD

fluent, es,kibana
Vector + ClickHouse + Grafana
Loki + Promtail + Grafana
devops
Jenkins + Helm + Kubernetes
GitLab CI/CD + Kubernetes

JuiceFS
JuiceFS
https://github.com/apache/apisix-ingress-controller
https://github.com/apache/apisix-ingress-controller
https://github.com/k8snetworkplumbingwg/sriov-network-device-plugin
kube-sriov-device-plugin  （PCIe设备虚拟化）
https://github.com/NVIDIA/k8s-device-plugin
https://volcano.sh/zh/docs/plugins/
keda      （弹性伸缩）	https://github.com/kedacore/keda
gpu-manager  （GPU虚拟化） 	https://github.com/tkestack/gpu-manager
https://vector.dev/docs/reference/configuration/sources/kubernetes_logs/
