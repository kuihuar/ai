## master node
- kube-api-server (systemd service)（Static Pod）
- kube-scheduler (systemd service)（Static Pod）
- kube-controller-manager(systemd service)（Static Pod）
- etcd(systemd service)（Static Pod）
- CNI 插件
  1. kube-proxy
  2. Calico(calico-node DaemonSet)Flannel
- Addon Components

  - Metrics Server
  - DasBoard

## worker node
- kubelet(systemd service) 节点代理
- Container Runtime（containerd， CRI-O）
- CNI 插件
  1. kube-proxy
  2. Calico( DaemonSet)Flannel
  3. calico-node
- CoreDNS


CNI插件（DaemonSet）：CNI插件不是Master节点的系统组件，它是集群的网络插件，以工作负载（Pod）的形式运行在所有节点（包括Master）上