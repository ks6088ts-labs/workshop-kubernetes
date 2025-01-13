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

## Infrastructure

### Create Azure Kubernetes Service (AKS) Cluster

To create a Kubernetes cluster in Azure, you need to run [scripts/create_kubernetes_cluster.sh](../scripts/create_kubernetes_cluster.sh) script.

Refer to [Quickstart: Deploy an Azure Kubernetes Service (AKS) cluster using Azure CLI](https://learn.microsoft.com/azure/aks/learn/quick-kubernetes-deploy-cli) for more information.

```shell
sh scripts/create_kubernetes_cluster.sh
```

### Delete Azure Kubernetes Service (AKS) Cluster

To delete the Kubernetes cluster, simply delete the resource group.

```shell
RESOURCE_GROUP_NAME=rg-workshop-kubernetes-DEADBEEF

az group delete \
  --name $RESOURCE_GROUP_NAME \
  --yes --no-wait --verbose
```

### Connect to Azure Kubernetes Service (AKS) Cluster

```shell
RANDOM_SUFFIX=DEADBEEF
RESOURCE_GROUP_NAME=rg-workshop-kubernetes-$RANDOM_SUFFIX
CLUSTER_NAME=workshop-kubernetes-$RANDOM_SUFFIX

# Connect to the Kubernetes cluster
az aks get-credentials \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $CLUSTER_NAME \
  --verbose

# Verify the connection
kubectl get nodes

# Contexts
kubectl config -h
kubectl config get-contexts
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

### Grafana

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

### Argo CD

[Argo CD > Getting Started](https://argo-cd.readthedocs.io/en/stable/getting_started/)

```shell
NAMESPACE=argocd

# Deploy Argo CD
kubectl create namespace $NAMESPACE

# Set the namespace
kubectl config set-context --current --namespace=$NAMESPACE

# Install Argo CD
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Access The Argo CD API Server: https://argo-cd.readthedocs.io/en/stable/getting_started/#3-access-the-argo-cd-api-server
kubectl port-forward service/argocd-server 8080:443
```

To create an application from a Git repository, refer to [Create An Application From A Git Repository](https://argo-cd.readthedocs.io/en/stable/getting_started/#6-create-an-application-from-a-git-repository).

```shell
# Retrieve the initial password
argocd admin initial-password

# Set the password
PASSWORD=YOUR_PASSWORD

# Login to Argo CD
argocd login localhost:8080 \
  --username admin \
  --password $PASSWORD \
  --insecure

# Update the password
argocd account update-password

# Delete the initial password
kubectl delete secret argocd-initial-admin-secret -n $NAMESPACE

# Logout
argocd logout localhost:8080

# Login to Argo CD
argocd login localhost:8080 \
  --username admin \
  --password $PASSWORD \
  --insecure

# Create an application
argocd app create guestbook \
  --repo https://github.com/argoproj/argocd-example-apps.git \
  --path guestbook \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace default

# List the applications
argocd app list

# Sync the application
argocd app sync guestbook

# Verify the deployment
kubectl get svc
kubectl port-forward service/guestbook-ui 8888:80

# Delete the application
argocd app delete guestbook --yes
```

## Use helm to deploy apps

- [Helm > Quickstart Guide](https://helm.sh/docs/intro/quickstart/)

### Prometheus

- [つくって、壊して、直して学ぶ Kubernetes 入門 > Chapter 11 　オブザーバビリティとモニタリングに触れてみよう](https://www.shoeisha.co.jp/book/detail/9784798183961)

```shell
# Add the Prometheus Helm repository
helm repo list
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Create the namespace
NAMESPACE=monitoring
kubectl create ns $NAMESPACE

# Deploy Prometheus
helm install kube-prometheus-stack \
  --namespace monitoring \
  prometheus-community/kube-prometheus-stack

# Verify the deployment
helm list -n $NAMESPACE
kubectl get pod -n $NAMESPACE --watch

# Access to the Prometheus dashboard
kubectl port-forward service/kube-prometheus-stack-prometheus 9090:9090 -n $NAMESPACE

# Access to the Grafana dashboard (username: admin, password: prom-operator)
kubectl port-forward service/kube-prometheus-stack-grafana 3000:80 -n $NAMESPACE

# Delete the deployment
helm uninstall kube-prometheus-stack -n $NAMESPACE
```

Collect metrics from the HTTP server.

```shell
# Launch HTTP server on develop namespace
kubectl apply -f k8s/collect-metrics/namespace.yaml
kubectl apply -f k8s/collect-metrics/http-server.yaml

# Update the Prometheus configuration
helm upgrade kube-prometheus-stack \
  -f k8s/collect-metrics/kube-prometheus-stack-values.yaml \
  -n $NAMESPACE \
  prometheus-community/kube-prometheus-stack

# Verify the deployment from Prometheus dashboard
kubectl port-forward service/kube-prometheus-stack-prometheus 9090:9090 -n $NAMESPACE
# Query `go_gc_duration_seconds{job="http-server"}` in the Prometheus dashboard
```

### Argo CD

- [クラウドネイティブで実現する マイクロサービス開発・運用 実践ガイド > 7.2 　 Argo CD による GitOps の実装](https://gihyo.jp/book/2023/978-4-297-13783-0)

```shell
REPO_NAME=argo-cd
helm repo add $REPO_NAME https://argoproj.github.io/argo-helm

NAMESPACE=gitops

# Deploy Argo CD
kubectl create namespace $NAMESPACE
# Set the namespace
kubectl config set-context --current --namespace=$NAMESPACE

helm install -n $NAMESPACE mygitops $REPO_NAME/argo-cd

helm list --all-namespaces

# Access The Argo CD API Server
kubectl port-forward service/mygitops-argocd-server 8080:443

# Uninstall Argo CD
helm uninstall mygitops -n $NAMESPACE
```

# References

- [Docker/Kubernetes 実践コンテナ開発入門 改訂新版](https://gihyo.jp/book/2024/978-4-297-14017-5)
  - https://github.com/gihyodocker
