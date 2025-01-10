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

# References

- [つくって、壊して、直して学ぶ Kubernetes 入門](https://www.shoeisha.co.jp/book/detail/9784798183961)
  - https://github.com/aoi1/bbf-kubernetes
