# Kubecost 利用環境の構築

このシナリオでは、Kubecost を利用するための環境を構築します。

## Kubecost のインストール

[Installing Kubecost](https://docs.kubecost.com/install-and-configure/install) を参照して Kubecost をインストールします。

```shell
helm repo add kubecost https://kubecost.github.io/cost-analyzer/
helm repo update
helm search repo kubecost

# Deploy the Kubecost Controller
helm install kubecost kubecost/cost-analyzer \
  --create-namespace \
  --namespace kubecost \
  --version 2.6.2

# Access the Kubecost UI
k -n kubecost port-forward deployments/kubecost-cost-analyzer 9090:9090
```
