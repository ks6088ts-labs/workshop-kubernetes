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
