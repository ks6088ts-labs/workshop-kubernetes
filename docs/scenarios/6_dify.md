# Dify 利用環境の構築

このシナリオでは、Dify を利用するための環境を Kubernetes 上に構築します。

## Dify

[Advanced Setup](https://github.com/langgenius/dify?tab=readme-ov-file#advanced-setup) に従って、Dify をデプロイします。

構築方法が複数あるが、ここでは [Winson-030/dify-kubernetes](https://github.com/Winson-030/dify-kubernetes) を利用します。

```shell
NAMESPACE=dify
MANIFEST=https://raw.githubusercontent.com/Winson-030/dify-kubernetes/main/dify-deployment.yaml

# alias kubectl=k
k create namespace $NAMESPACE

# deploy
k -n $NAMESPACE apply -f $MANIFEST

# reset password (何故かこれをやらないと動かない)
k -n $NAMESPACE exec svc/dify-api -- flask reset-password

# connect from localhost
k -n $NAMESPACE port-forward svc/dify-nginx 8080:80

# open http://localhost:8080
# initial password = `password`

# cleanup
k -n $NAMESPACE delete -f $MANIFEST
```

## Troubleshooting

### AKS クラスターのスケールアップ

CPU やメモリが不足している場合は、AKS クラスターをスケールアップする。
(これが適切であるかは自信がない)

```shell
RESOURCE_GROUP_NAME=YOUR_RESOURCE_GROUP_NAME
CLUSTER_NAME=YOUR_CLUSTER_NAME

# Start the AKS cluster
az aks start \
  --name $CLUSTER_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --no-wait

# Verify the current node count
az aks show \
    --resource-group $RESOURCE_GROUP_NAME \
    --name $CLUSTER_NAME \
    --query agentPoolProfiles

# (Optional) Scale the node pool to 2 nodes
az aks scale \
    --resource-group $RESOURCE_GROUP_NAME \
    --name $CLUSTER_NAME \
    --node-count 2 \
    --nodepool-name default
```

### 初期パスワードの設定

[4. How to reset the password of the admin account?](https://docs.dify.ai/getting-started/install-self-hosted/faqs#id-4.-how-to-reset-the-password-of-the-admin-account)

```shell
k -n $NAMESPACE exec svc/dify-api -- flask reset-password
```
