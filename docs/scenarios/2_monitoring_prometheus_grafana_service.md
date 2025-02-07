# Prometheus と Grafana を使用してメトリクスを収集する

このシナリオでは、HTTP サービスをデプロイし、Prometheus と Grafana を使用してメトリクスを収集します。

## 事前条件

- [AKS クラスターのセットアップ](./0_setup_aks_cluster.md) を完了していること

## チュートリアル

### Grafana のデプロイ

[Deploy Grafana on Kubernetes](https://grafana.com/docs/grafana/latest/setup-grafana/installation/kubernetes/)

```shell
NAMESPACE=grafana

# Create the namespace
kubectl create namespace $NAMESPACE

# Deploy Grafana
kubectl apply -f k8s/grafana.yaml --namespace $NAMESPACE

# Verify the deployment
kubectl get pvc --namespace=$NAMESPACE
kubectl get deployments --namespace=$NAMESPACE
kubectl get svc --namespace=$NAMESPACE
kubectl get all --namespace=$NAMESPACE

# Forward the port
kubectl port-forward service/grafana 3000:3000 --namespace=$NAMESPACE
# Access the Grafana dashboard at http://localhost:3000

# Delete the deployment
kubectl delete -f k8s/grafana.yaml --namespace $NAMESPACE

# Delete the namespace
kubectl delete namespace $NAMESPACE
```

### Prometheus

```shell
# 変数の設定
CHART_REPOSITORY_NAME=prometheus-community
CHART_REPOSITORY_URL=https://prometheus-community.github.io/helm-charts
RELEASE_NAME=my-prometheus-stack
CHART_NAME=kube-prometheus-stack
NAMESPACE=monitoring
VERSION=68.2.1
VALUES_FILE=k8s/kube-prometheus-stack-values.yaml

# Helm リポジトリの追加
helm repo add $CHART_REPOSITORY_NAME $CHART_REPOSITORY_URL

# チャートのバージョンを確認
helm search repo $CHART_REPOSITORY_NAME/$CHART_NAME --versions

# Chart のデプロイ
helm install $RELEASE_NAME $CHART_REPOSITORY_NAME/$CHART_NAME \
  --namespace $NAMESPACE \
  --create-namespace \
  --version $VERSION \
  --values $VALUES_FILE

# デプロイの確認
helm list -n $NAMESPACE

# Prometheus ダッシュボードへのアクセス
kubectl port-forward -n $NAMESPACE services/$RELEASE_NAME-kube-p-prometheus 9090:9090
# ブラウザで http://localhost:9090 にアクセス (username: admin, password: prom-operator)
kubectl port-forward -n $NAMESPACE services/$RELEASE_NAME-grafana 3000:80
# Query `go_gc_duration_seconds{job="http-server"}` in the Prometheus dashboard
```

Collect metrics from the HTTP server.

```shell
kubectl create namespace develop
kubens develop

# HTTP サーバーのデプロイ
kubectl apply -f k8s/http-server/

# 設定の変更
NAMESPACE=monitoring
helm upgrade $RELEASE_NAME \
  -f k8s/kube-prometheus-stack-values.yaml \
  -n $NAMESPACE \
  prometheus-community/kube-prometheus-stack
```

### Ingress

#### Grafana

```shell
# Grafana の Ingress リソースの作成
k -n monitoring apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: monitoring-ingress
spec:
  ingressClassName: nginx
  rules:
    - http:
        paths:
          - pathType: Prefix
            path: /
            backend:
              service:
                name: $RELEASE_NAME-grafana
                port:
                  number: 80
EOF

# Ingress リソースの削除
k -n monitoring delete ingress monitoring-ingress
```

#### Prometheus

```shell
# Grafana の Ingress リソースの作成
k -n monitoring apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: monitoring-ingress
spec:
  ingressClassName: nginx
  rules:
    - http:
        paths:
          - pathType: Prefix
            path: /
            backend:
              service:
                name: $RELEASE_NAME-kube-p-prometheus
                port:
                  number: 9090
EOF

# Ingress リソースの削除
k -n monitoring delete ingress monitoring-ingress
```
