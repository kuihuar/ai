# 安全机制

## 📖 安全概述

Kubernetes 提供了多层次的安全机制来保护集群和应用程序。从认证授权到网络策略，从 Pod 安全到运行时安全，Kubernetes 构建了完整的安全防护体系。

## 🔐 认证 (Authentication)

### 1. 证书认证
基于 TLS 证书的认证方式，适用于集群内部组件。

```bash
# 生成证书
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=client"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365
```

### 2. Token 认证
基于 Bearer Token 的认证方式。

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-sa
  namespace: default
---
apiVersion: v1
kind: Secret
metadata:
  name: app-token
  namespace: default
  annotations:
    kubernetes.io/service-account.name: app-sa
type: kubernetes.io/service-account-token
```

### 3. OpenID Connect
集成 OAuth2 和 OpenID Connect 进行身份认证。

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-apiserver-config
data:
  oidc-issuer-url: "https://accounts.google.com"
  oidc-client-id: "your-client-id"
  oidc-username-claim: "email"
  oidc-groups-claim: "groups"
```

## 🔑 授权 (Authorization)

### 1. RBAC (Role-Based Access Control)
基于角色的访问控制，是 Kubernetes 的主要授权机制。

#### Role 和 ClusterRole
```yaml
# 命名空间级别的角色
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]

---
# 集群级别的角色
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
```

#### RoleBinding 和 ClusterRoleBinding
```yaml
# 命名空间级别的角色绑定
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
subjects:
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io

---
# 集群级别的角色绑定
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: read-secrets-global
subjects:
- kind: Group
  name: manager
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

### 2. ServiceAccount
为 Pod 提供身份认证。

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-sa
  namespace: default
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  template:
    spec:
      serviceAccountName: app-sa
      containers:
      - name: app
        image: myapp:latest
```

## 🛡️ Pod 安全

### 1. Pod Security Standards
Kubernetes 定义了三个 Pod 安全级别。

#### Privileged (特权)
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: privileged-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    securityContext:
      privileged: true
      runAsUser: 0
      capabilities:
        add: ["ALL"]
```

#### Baseline (基线)
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: baseline-pod
spec:
  containers:
  - name: app
    image: nginx:latest
    securityContext:
      runAsNonRoot: true
      runAsUser: 1000
      allowPrivilegeEscalation: false
      capabilities:
        drop: ["ALL"]
```

#### Restricted (限制)
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: restricted-pod
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 2000
  containers:
  - name: app
    image: nginx:latest
    securityContext:
      runAsNonRoot: true
      runAsUser: 1000
      allowPrivilegeEscalation: false
      capabilities:
        drop: ["ALL"]
      readOnlyRootFilesystem: true
```

### 2. Pod Security Admission
在 Pod 创建时进行安全策略检查。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-pod
  labels:
    pod-security.kubernetes.io/enforce: baseline
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
spec:
  containers:
  - name: app
    image: nginx:latest
```

## 🌐 网络策略

### 1. NetworkPolicy
控制 Pod 间的网络通信。

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

### 2. 允许特定流量
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-web-traffic
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 80
    - protocol: TCP
      port: 443
```

### 3. 出站策略
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Egress
  egress:
  - to: []
    ports:
    - protocol: UDP
      port: 53
    - protocol: TCP
      port: 53
```

## 🔒 准入控制

### 1. ValidatingAdmissionWebhook
验证请求的准入控制器。

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pod-policy.example.com
webhooks:
- name: pod-policy.example.com
  rules:
  - apiGroups: [""]
    apiVersions: ["v1"]
    operations: ["CREATE"]
    resources: ["pods"]
    scope: "Namespaced"
  clientConfig:
    service:
      namespace: "example-system"
      name: "pod-policy-webhook"
      path: "/validate"
  admissionReviewVersions: ["v1"]
  sideEffects: None
```

### 2. MutatingAdmissionWebhook
修改请求的准入控制器。

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: pod-mutator.example.com
webhooks:
- name: pod-mutator.example.com
  rules:
  - apiGroups: [""]
    apiVersions: ["v1"]
    operations: ["CREATE"]
    resources: ["pods"]
    scope: "Namespaced"
  clientConfig:
    service:
      namespace: "example-system"
      name: "pod-mutator-webhook"
      path: "/mutate"
  admissionReviewVersions: ["v1"]
  sideEffects: None
```

## 🛠️ 安全工具

### 1. kube-bench
检查 Kubernetes 集群安全配置。

```bash
# 运行 kube-bench
kube-bench --benchmark cis-1.6

# 生成报告
kube-bench --benchmark cis-1.6 --json > report.json
```

### 2. Falco
运行时安全监控。

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: falco-config
data:
  falco.yaml: |
    rules_file:
      - /etc/falco/falco_rules.yaml
      - /etc/falco/k8s_audit_rules.yaml
    webserver:
      enabled: true
      listen_port: 9376
```

### 3. OPA Gatekeeper
策略执行引擎。

```yaml
apiVersion: config.gatekeeper.sh/v1alpha1
kind: Config
metadata:
  name: config
  namespace: gatekeeper-system
spec:
  sync:
    syncOnly:
    - group: ""
      version: "v1"
      kind: "Pod"
```

## 🎯 安全最佳实践

### 1. 身份认证
- 使用强密码和证书
- 定期轮换凭据
- 启用多因素认证

### 2. 访问控制
- 遵循最小权限原则
- 定期审查权限
- 使用 ServiceAccount

### 3. 网络安全
- 配置网络策略
- 使用 TLS 加密
- 监控网络流量

### 4. 运行时安全
- 使用非特权容器
- 扫描容器镜像
- 监控异常行为

## 🛠️ 实践练习

### 练习 1：RBAC 配置
1. 创建角色和角色绑定
2. 测试权限控制
3. 审计访问日志

### 练习 2：网络策略
1. 创建网络策略
2. 测试网络隔离
3. 配置允许规则

### 练习 3：Pod 安全
1. 配置 Pod 安全策略
2. 测试安全限制
3. 监控安全事件

## 📚 扩展阅读

- [Kubernetes 安全官方文档](https://kubernetes.io/docs/concepts/security/)
- [RBAC 官方文档](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
- [网络策略官方文档](https://kubernetes.io/docs/concepts/services-networking/network-policies/)

## 🎯 下一步

掌握安全机制后，继续学习：
- [监控与日志](./10-monitoring/README.md)
- [Helm包管理](./11-helm/README.md) 