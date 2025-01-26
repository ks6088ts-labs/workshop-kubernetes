# Workshop

## Fundamentals

### Hands on Docker

```shell
# Run the nginx container
docker run --rm --detach --publish 8080:80 --name web nginx:latest

# Verify the container
curl http://localhost:8080 -v

# Stop the container
docker stop web
```

### Hands on Go

Run the Workshop CLI

```shell
# Build CLI
make build

# Help
./dist/workshop-kubernetes --help

# Run HTTP server from the CLI
./dist/workshop-kubernetes sandbox http --port 8888
```

## Deploy Applications

### Kubernetes Dashboard

[github.com/kubernetes/dashboard > aio/deploy/recommended.yaml](https://github.com/kubernetes/dashboard/blob/v2.7.0/aio/deploy/recommended.yaml)

```shell
# Deploy the Kubernetes dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/refs/tags/v2.7.0/aio/deploy/recommended.yaml

# Verify the deployment
kubectl get pods --namespace kubernetes-dashboard
```

Create [dashboard-user.yaml](../k8s/dashboard-user.yaml) file with the following content.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: admin-user
    namespace: kubernetes-dashboard
```

Apply the manifest file.

```shell
# Apply the dashboard user
kubectl apply -f k8s/dashboard-user.yaml

# Create a token for the dashboard user
TOKEN=$(kubectl -n kubernetes-dashboard create token admin-user)
```

Run the following command to start a proxy.

```shell
# Run the proxy
kubectl proxy
```

Access the Kubernetes dashboard at http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ after running the proxy.

### nginx

```shell
# Deploy the nginx app
kubectl apply -f k8s/nginx.yaml

# Forward the port
kubectl port-forward nginx 8080:80

# Verify the app
curl http://localhost:8080 -v

# Delete the app
kubectl delete -f k8s/nginx.yaml
```

### workshop-kubernetes CLI

```shell
# Deploy the workshop-kubernetes app
kubectl apply -f k8s/workshop-kubernetes.yaml
kubectl apply -f k8s/cronjob.yaml

# Forward the port
kubectl port-forward workshop-kubernetes 8080:8080

# Verify the app
curl http://localhost:8080 -v

# Watch for changes in the pods
kubectl get pods --watch

# Delete the app
kubectl delete -f k8s/workshop-kubernetes.yaml
kubectl delete -f k8s/cronjob.yaml
```

## Use helm to deploy apps

- [Helm > Quickstart Guide](https://helm.sh/docs/intro/quickstart/)
- [Artifact Hub](https://artifacthub.io/)

```shell
# Helm チャートの作成
CHART_NAME=charts/my-chart
mkdir -p charts

helm create $CHART_NAME
helm template $CHART_NAME
helm package $CHART_NAME \
  --version v0.0.1 \
  --app-version v0.0.1
```

### Elastic Stack

- [Elastic Stack Helm Chart](https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-stack-helm-chart.html)

```shell
NAMESPACE=elastic-stack

# Add the Elastic Helm repository
helm repo add elastic https://helm.elastic.co
helm repo update

# Search for the Elastic Stack chart
helm search repo elastic/

# Deploy the Elastic Stack
helm install es-quickstart elastic/eck-stack \
  -n $NAMESPACE \
  --create-namespace

# Verify the deployment
helm list --all-namespaces
```

## Tools

### k9s

- [k9s](https://github.com/derailed/k9s)
- [k9s 使い方まとめ](https://qiita.com/t0m0ya/items/15ee4d43dcda4701946c)

```shell
# Install k9s
go install github.com/derailed/k9s@latest

# Run k9s
k9s
```

# References

- [Docker/Kubernetes 実践コンテナ開発入門 改訂新版](https://gihyo.jp/book/2024/978-4-297-14017-5)
  - https://github.com/gihyodocker
