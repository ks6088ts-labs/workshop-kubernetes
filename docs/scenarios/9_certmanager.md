AKS で独自ドメインを使った HTTPS の Web アプリケーションをデプロイするための手順です。ここでは、cert-manager を使用して Let's Encrypt から証明書を取得し、Ingress リソースを通じて HTTPS を有効にします。

<!-- @gemini AKS で独自ドメインを使ったHTTPS対応なWebサーバー公開方法を一番シンプルにわかる方法で教えて -->

```shell
# web applicationのデプロイメントを作成
kubectl apply -f k8s/cert-manager/nginx-deployment.yaml

# Ingress Nginx Controllerのデプロイメントを ingress-nginx 名前空間に作成
# https://kubernetes.github.io/ingress-nginx/deploy/#azure
# https://learn.microsoft.com/ja-jp/troubleshoot/azure/azure-kubernetes/load-bal-ingress-c/create-unmanaged-ingress-controller?tabs=azure-cli#create-an-ingress-controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.13.0/deploy/static/provider/cloud/deploy.yaml

# External IPアドレスの確認 (ingress-nginx-controller の LoadBalancer)
kubectl get svc -n ingress-nginx

# cert-manager を cert-manager 名前空間にデプロイ
# https://cert-manager.io/docs/installation/kubectl/
# Install all cert-manager components:
# kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.18.2/cert-manager.yaml
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.18.2 \
  --set crds.enabled=true

# ClusterIssuerの作成
# https://cert-manager.io/docs/configuration/acme/
kubectl apply -f k8s/cert-manager/letsencrypt-clusterissuer.yaml

# Ingressリソースの作成
# https://cert-manager.io/docs/usage/ingress/
kubectl apply -f k8s/cert-manager/nginx-ingress.yaml

# Ingressリソースの状態を確認
kubectl get ingress nginx-ingress -w

# 証明書の状態を確認
kubectl get certificaterequest -A
kubectl get certificate -A
kubectl describe certificate ks6088ts-com-tls
```

### [ks6088ts-labs/template-fastapi](https://github.com/ks6088ts-labs/template-fastapi) をデプロイして公開する

```shell
# FIXME: 環境変数の設定
# k8s/cert-manager/template-fastapi-deployment.yaml の env セクションを編集

# FastAPIのデプロイメントを作成
kubectl apply -f k8s/cert-manager/template-fastapi-deployment.yaml

# Ingressリソースの作成 (仕向けるサービスは template-fastapi-service)
# https://cert-manager.io/docs/usage/ingress/
kubectl apply -f k8s/cert-manager/nginx-ingress.yaml
```
