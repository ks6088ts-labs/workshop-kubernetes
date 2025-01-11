# Workshop

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

### Deploy Kubernetes Dashboard

[github.com/kubernetes/dashboard > aio/deploy/recommended.yaml](https://github.com/kubernetes/dashboard/blob/v2.7.0/aio/deploy/recommended.yaml)

```shell
# Deploy the Kubernetes dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/refs/tags/v2.7.0/aio/deploy/recommended.yaml

# Verify the deployment
kubectl get pods --namespace kubernetes-dashboard
```

Create [dashboard-user.yaml](../manifests/dashboard-user.yaml) file with the following content.

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
kubectl apply -f manifests/dashboard-user.yaml

# Create a token for the dashboard user
TOKEN=$(kubectl -n kubernetes-dashboard create token admin-user)
```

Run the following command to start a proxy.

```shell
# Run the proxy
kubectl proxy
```

Access the Kubernetes dashboard at http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ after running the proxy.

### Hands on Docker

```shell
# Run the nginx container
docker run --rm --detach --publish 8080:80 --name web nginx:latest

# Verify the container
curl http://localhost:8080 -v

# Stop the container
docker stop web
```

### Run the Workshop CLI

```shell
# Build CLI
make build

# Help
./dist/workshop-kubernetes --help

# Run HTTP server from the CLI
./dist/workshop-kubernetes sandbox http --port 8888
```

### Run apps in Kubernetes

#### nginx

```shell
# Deploy the nginx app
kubectl apply -f manifests/nginx.yaml

# Forward the port
kubectl port-forward nginx 8080:80

# Verify the app
curl http://localhost:8080 -v

# Delete the app
kubectl delete -f manifests/nginx.yaml
```

#### workshop-kubernetes

```shell
# Deploy the workshop-kubernetes app
kubectl apply -f manifests/workshop-kubernetes.yaml
kubectl apply -f manifests/cronjob.yaml

# Forward the port
kubectl port-forward workshop-kubernetes 8080:8080

# Verify the app
curl http://localhost:8080 -v

# Watch for changes in the pods
kubectl get pods --watch

# Delete the app
kubectl delete -f manifests/workshop-kubernetes.yaml
kubectl delete -f manifests/cronjob.yaml
```

#### Grafana

[Deploy Grafana on Kubernetes](https://grafana.com/docs/grafana/latest/setup-grafana/installation/kubernetes/)

```shell
NAMESPACE=grafana

# Create the namespace
kubectl create namespace $NAMESPACE

# Deploy Grafana
kubectl apply -f manifests/grafana.yaml --namespace $NAMESPACE

# Verify the deployment
kubectl get pvc --namespace=$NAMESPACE
kubectl get deployments --namespace=$NAMESPACE
kubectl get svc --namespace=$NAMESPACE
kubectl get all --namespace=$NAMESPACE

# Forward the port
kubectl port-forward service/grafana 3000:3000 --namespace=$NAMESPACE
# Access the Grafana dashboard at http://localhost:3000

# Delete the deployment
kubectl delete -f manifests/grafana.yaml --namespace $NAMESPACE

# Delete the namespace
kubectl delete namespace $NAMESPACE
```

### Use helm to deploy apps

- [Helm](https://helm.sh/docs/intro/quickstart/)

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

# Delete the deployment
helm uninstall kube-prometheus-stack -n $NAMESPACE
```

# References

- [Docker/Kubernetes 実践コンテナ開発入門 改訂新版](https://gihyo.jp/book/2024/978-4-297-14017-5)
  - https://github.com/gihyodocker
- [つくって、壊して、直して学ぶ Kubernetes 入門](https://www.shoeisha.co.jp/book/detail/9784798183961)
  - https://github.com/aoi1/bbf-kubernetes
