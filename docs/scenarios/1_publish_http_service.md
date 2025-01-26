# Ingress NGINX Controller による HTTP サービスの公開

このシナリオでは、Ingress NGINX Controller を使用して、HTTP サービスを公開します。

## 事前条件

- [AKS クラスターのセットアップ](./0_setup_aks_cluster.md) を完了していること

### Ingress-Nginx Controller

- [Ingress NGINX Controller](https://github.com/kubernetes/ingress-nginx)
- [Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/#quick-start)
- [Create an ingress controller](https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/load-bal-ingress-c/create-unmanaged-ingress-controller?tabs=azure-cli)

```shell
NAMESPACE=ingress-nginx

# use a YAML manifest to deploy the Ingress-Nginx Controller
# kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.0/deploy/static/provider/cloud/deploy.yaml

# Use Helm to deploy the Ingress-Nginx Controller
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# Deploy the Ingress-Nginx Controller
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --create-namespace \
  --namespace $NAMESPACE \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-load-balancer-health-probe-request-path"=/healthz \
  --set controller.service.externalTrafficPolicy=Local

# Verify the deployment
kubectl get svc -n $NAMESPACE

# Uninstall the Ingress-Nginx Controller
helm uninstall ingress-nginx -n $NAMESPACE
```

- [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)

```shell
NAMESPACE=develop

# Create a namespace
kubectl create namespace $NAMESPACE

kubens $NAMESPACE

# HTTP サーバーのデプロイ
kubectl apply -f k8s/http-server/

# Deploy Ingress resource
kubectl apply -f k8s/ingress.yaml

kubectl get svc -A

# Verify the deployment
EXTERNAL_IP=xxx.xxx.xxx.xxx
curl http://$EXTERNAL_IP/http-server --verbose
```
