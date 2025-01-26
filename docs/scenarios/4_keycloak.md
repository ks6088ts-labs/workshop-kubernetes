# Keycloak を使用した認証

このシナリオでは、Keycloak を使用してアプリケーションに認証を追加する方法を学びます。

## 事前条件

- [AKS クラスターのセットアップ](./0_setup_aks_cluster.md) を完了していること

## チュートリアル

### Keycloak のデプロイ

- [Bitnami package for Keycloak](https://bitnami.com/stack/keycloak/helm)
- [Installing the Chart](https://github.com/bitnami/charts/tree/main/bitnami/keycloak/#installing-the-chart)

```shell
NAMESPACE=iam
RELEASE_NAME=myiam

# Login to Docker Hub
docker login -u $DOCKER_USERNAME

# Install Keycloak
helm install \
  --namespace $NAMESPACE \
  --create-namespace \
  $RELEASE_NAME \
  oci://registry-1.docker.io/bitnamicharts/keycloak

# Verify the deployment
helm list -A

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
