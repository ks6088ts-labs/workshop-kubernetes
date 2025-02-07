# Argo CD による GitOps

このシナリオでは、Argo CD を使用して GitOps を実装します。

## 事前条件

- [AKS クラスターのセットアップ](./0_setup_aks_cluster.md) を完了していること

## チュートリアル

### Argo CD のデプロイ

#### kubectl での Argo CD のデプロイ

[Argo CD > Getting Started](https://argo-cd.readthedocs.io/en/stable/getting_started/)

```shell
NAMESPACE=argocd

# Namespace の作成
kubectl create namespace $NAMESPACE

# Namespace の設定
kubectl config set-context --current --namespace=$NAMESPACE

# Argo CD のデプロイ
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Argo CD API サーバーへのアクセス: https://argo-cd.readthedocs.io/en/stable/getting_started/#3-access-the-argo-cd-api-server
kubectl port-forward service/argocd-server 8080:443
```

#### Helm チャートを利用した Argo CD のデプロイ

- [クラウドネイティブで実現する マイクロサービス開発・運用 実践ガイド > 7.2 　 Argo CD による GitOps の実装](https://gihyo.jp/book/2023/978-4-297-13783-0)

```shell
REPO_NAME=argo
NAMESPACE=argo-cd
RELEASE_NAME=my-argo-cd

helm repo add $REPO_NAME https://argoproj.github.io/argo-helm

helm install \
  --namespace $NAMESPACE \
  --create-namespace \
  $RELEASE_NAME \
  $REPO_NAME/argo-cd

# Set the namespace
kubectl config set-context --current --namespace=$NAMESPACE

helm list -A

# Access The Argo CD API Server
kubectl port-forward service/$RELEASE_NAME-argocd-server 8080:443

# Uninstall Argo CD
helm uninstall $RELEASE_NAME -n $NAMESPACE
```

### アプリケーションのデプロイ

アプリケーションを Git リポジトリから作成するには、[Create An Application From A Git Repository](https://argo-cd.readthedocs.io/en/stable/getting_started/#6-create-an-application-from-a-git-repository) を参照してください。

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
kubectl port-forward -n guestbook service/guestbook-ui 8888:80

# Delete the application
argocd app delete guestbook --yes
```

### Argo Workflows のデプロイ

- [Argo Workflows - The workflow engine for Kubernetes > Quick Start](https://argo-workflows.readthedocs.io/en/latest/quick-start/)

```shell
REPO_NAME=argo
NAMESPACE=argo-workflows
RELEASE_NAME=my-argo-workflows

helm install \
  --namespace $NAMESPACE \
  --create-namespace \
  $RELEASE_NAME \
  $REPO_NAME/argo-workflows

kubens $NAMESPACE

# Get token
POD_NAME=$(kubectl get pods -l app=server -o jsonpath='{.items[0].metadata.name}')
k exec -it $POD_NAME -- argo auth token

# Access the Argo Workflows UI
kubectl port-forward svc/$RELEASE_NAME-server 2746:2746 --address='0.0.0.0'

argo list -A

# Create cluster role
k apply -f k8s/argo-workflows/

# Submit a workflow
argo submit --watch https://raw.githubusercontent.com/argoproj/argo-workflows/main/examples/hello-world.yaml

# kubectl apply -f - <<EOF
# apiVersion: rbac.authorization.k8s.io/v1
# kind: ClusterRole
# metadata:
#   name: argo-workflow-argo-workflows-workflow-controller
# rules:
# - apiGroups: ["argoproj.io"]
#   resources: ["workflowtaskresults"]
#   verbs: ["create", "get", "list", "watch", "update", "patch", "delete"]
# EOF
```

### HTTP サーバーのデプロイ

```shell
# Create HTTP Services
argocd app create http-server \
  --repo https://github.com/ks6088ts-labs/workshop-kubernetes.git \
  --path k8s/http-server \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace develop

argocd app sync http-server
```

### Ingress リソースの作成

```shell
NAMESPACE=argo-workflows
RELEASE_NAME=my-argo-workflows
SERVICE=$RELEASE_NAME-server
PORT=2746

# Ingress リソースの作成
k -n $NAMESPACE apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: $NAMESPACE-ingress
spec:
  ingressClassName: nginx
  rules:
    - http:
        paths:
          - pathType: Prefix
            path: /
            backend:
              service:
                name: $SERVICE
                port:
                  number: $PORT
EOF

# Ingress リソースの削除
k -n $NAMESPACE delete ingress $NAMESPACE-ingress
```
