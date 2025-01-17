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

### Manage AKS Cluster

- [Stop and start an Azure Kubernetes Service (AKS) cluster](https://learn.microsoft.com/azure/aks/start-stop-cluster?tabs=azure-cli)

```shell
# Start the AKS cluster
az aks start \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME

# Verify the AKS cluster status
az aks show \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME

# Stop the AKS cluster
az aks stop \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --no-wait
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

# Create namespace for guestbook
kubectl create namespace guestbook
kubectl config set-context --current --namespace=guestbook

# Create an application
argocd app create guestbook \
  --repo https://github.com/argoproj/argocd-example-apps.git \
  --path guestbook \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace guestbook

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
RELEASE_NAME=kube-prometheus-stack

# Deploy Prometheus
helm install $RELEASE_NAME \
  --namespace $NAMESPACE \
  --create-namespace \
  prometheus-community/kube-prometheus-stack

# Verify the deployment
helm list -n $NAMESPACE
kubectl get pod -n $NAMESPACE --watch

# Access to the Prometheus dashboard
kubectl port-forward service/kube-prometheus-stack-prometheus 9090:9090 -n $NAMESPACE

# Access to the Grafana dashboard (username: admin, password: prom-operator)
kubectl port-forward service/kube-prometheus-stack-grafana 3000:80 -n $NAMESPACE

# Delete the deployment
helm uninstall $RELEASE_NAME -n $NAMESPACE
```

Collect metrics from the HTTP server.

```shell
# Launch HTTP server on develop namespace
kubectl apply -f k8s/collect-metrics/namespace.yaml
kubectl apply -f k8s/collect-metrics/http-server.yaml

# Update the Prometheus configuration
helm upgrade $RELEASE_NAME \
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
NAMESPACE=gitops
RELEASE_NAME=mygitops

helm repo add $REPO_NAME https://argoproj.github.io/argo-helm

helm install \
  --namespace $NAMESPACE \
  --create-namespace \
  $RELEASE_NAME \
  $REPO_NAME/argo-cd

# Set the namespace
kubectl config set-context --current --namespace=$NAMESPACE

helm list --all-namespaces

# Access The Argo CD API Server
kubectl port-forward service/mygitops-argocd-server 8080:443

# Uninstall Argo CD
helm uninstall mygitops -n $NAMESPACE
```

### Keycloak

- [Bitnami package for Keycloak](https://bitnami.com/stack/keycloak/helm)
- [Installing the Chart](https://github.com/bitnami/charts/tree/main/bitnami/keycloak/#installing-the-chart)

```shell
NAMESPACE=iam
RELEASE_NAME=myiam

# Install Keycloak
helm install \
  --namespace $NAMESPACE \
  --create-namespace \
  $RELEASE_NAME \
  oci://registry-1.docker.io/bitnamicharts/keycloak

# Verify the deployment
helm list --all-namespaces

# Set namespace
kubens $NAMESPACE

# Port forward
kubectl port-forward service/myiam-keycloak 8080:80

# Retrieve the admin password
# https://github.com/bitnami/charts/issues/17522#issuecomment-1626921535
kubectl get secret myiam-keycloak \
  --output jsonpath='{.data.admin-password}' | base64 -d && echo

# Login to Keycloak at http://localhost:8080 with username `user` and the password
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
