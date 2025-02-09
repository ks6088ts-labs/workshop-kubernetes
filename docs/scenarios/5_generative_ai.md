# 生成 AI の利用環境の構築

このシナリオでは、生成 AI の利用環境を構築します。

## Open WebUI

- [Helm Setup for Kubernetes](https://docs.openwebui.com/getting-started/quick-start/)
- [Open WebUI > workshop-llm-agents/docs/references.md](https://github.com/ks6088ts-labs/workshop-llm-agents/blob/main/docs/references.md#open-webui)

```shell
helm repo add open-webui https://open-webui.github.io/helm-charts
helm repo update

NAMESPACE=genai
RELEASE_NAME=openwebui

# Install Open WebUI
helm install \
  --namespace $NAMESPACE \
  --create-namespace \
  $RELEASE_NAME \
  open-webui/open-webui

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
                name: open-webui
                port:
                  number: 80
EOF

# Ingress リソースの削除
k -n $NAMESPACE delete ingress $NAMESPACE-ingress
```
