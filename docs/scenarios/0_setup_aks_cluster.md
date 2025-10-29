# AKS クラスターのセットアップ

このシナリオでは、Azure Kubernetes Service (AKS) クラスターをセットアップします。

## チュートリアル

### Azure リソースの作成

#### Azure CLI

AKS クラスタを作成するために、Azure CLI を使用します。
[scripts/create_kubernetes_cluster.sh](../scripts/create_kubernetes_cluster.sh) を使用して、AKS クラスターを作成します。

[Quickstart: Deploy an Azure Kubernetes Service (AKS) cluster using Azure CLI](https://learn.microsoft.com/azure/aks/learn/quick-kubernetes-deploy-cli) を参考にしてください。

```shell
# AKS クラスターの作成
sh scripts/create_kubernetes_cluster.sh

# INFO: リソースの削除
az group delete \
  --name $RESOURCE_GROUP_NAME \
  --yes --no-wait --verbose
```

### AKS クラスター接続

- [Stop and start an Azure Kubernetes Service (AKS) cluster](https://learn.microsoft.com/azure/aks/start-stop-cluster?tabs=azure-cli)

```shell
RESOURCE_GROUP_NAME=your-resource-group-name
CLUSTER_NAME=your-aks-cluster-name

# AKS クラスターの情報の取得
az aks get-credentials \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $CLUSTER_NAME \
  --verbose

# kubectl のバージョン確認
kubectl get nodes

# kubectl の設定確認
kubectl config -h
kubectl config get-contexts

# (Optional) AKS クラスターの情報の確認
az aks show \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME

# (Optional) AKS クラスターを起動
az aks start \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --no-wait

# (Optional) AKS クラスターを停止
az aks stop \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --no-wait
```

### AKS でのデバッグ

```shell
# ノードの確認
kubectl get pods --all-namespaces
# kubectl get pods -A

# ノードの詳細情報の確認
kubectl exec -n monitoring -it $POD_NAME -- sh

# ノードのログの確認
kubectl run -it --rm debug --image=ubuntu:latest --restart=Never -- sh
```

コンテナ内で、次のコマンドを実行できます。

```shell
apt update && apt install -y curl

# Name resolution test
# curl http://<service-name>.<namespace>.svc.cluster.local:<port>
curl http://http-server.develop.svc.cluster.local:8080
curl http://my-prometheus-stack-prometheus-node-exporter.monitoring.svc.cluster.local:9100/metrics
```
