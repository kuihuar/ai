## master node
- kube-api-server (systemd service)
- kube-scheduler (systemd service)
- kube-controller-manager(systemd service)
- etcd(systemd service)
- CNI 插件
  1. kube-proxy
  2. Calico(calico-node DaemonSet)
- Addon Components
  - CoreDNS
  - Metrics Server
  - DasBoard

## worker node
- kubelet
- CNI 插件
  1. kube-proxy
  2. Calico(calico-node DaemonSet)
