# OpenTelemetry Demo 利用環境の構築

このシナリオでは、 [OpenTelemetry Demo](https://github.com/open-telemetry/opentelemetry-demo) を利用するための環境を構築します。

## インストール

[Kubernetes deployment](https://opentelemetry.io/docs/demo/kubernetes-deployment/) を参照して OpenTelemetry Demo をインストールします。

```shell
# Add OpenTelemetry Helm repository
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update

# Install OpenTelemetry Demo
helm install otel-demo open-telemetry/opentelemetry-demo \
  --create-namespace \
  --namespace otel

# Check the status
k9s -n otel

# Use the demo
k -n otel port-forward svc/frontend-proxy 8080:8080

# Web store: http://localhost:8080/
# Grafana: http://localhost:8080/grafana/
# Load Generator UI: http://localhost:8080/loadgen/
# Jaeger UI: http://localhost:8080/jaeger/ui/
# Flagd configurator UI: http://localhost:8080/feature

# Uninstall OpenTelemetry Demo
helm uninstall otel-demo -n otel
```
